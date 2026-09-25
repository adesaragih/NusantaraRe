package repository

// Pelari migrasi basis data.
//
// Untuk apa berkas ini: menjalankan berkas SQL di folder `migrations/` secara
// berurutan, sekali saja masing-masing, dan menyediakan jalur mundur.
//
// Dibaca sesudah: repository.go (yang memperkenalkan DB dan Qualify).
//
// Istilah yang dipakai di sini, sekali dijelaskan:
//   - migrasi    : satu langkah perubahan bentuk basis data, disimpan sebagai
//                  satu berkas .sql bernomor. Nomornya menentukan urutan.
//   - idempoten  : aman dijalankan berulang kali. Menjalankan dua kali memberi
//                  hasil yang sama dengan menjalankan sekali.
//   - embed      : menanam isi berkas ke dalam biner Go saat dibangun, sehingga
//                  program tidak perlu mencari berkas .sql di disk saat jalan.
//   - jalur mundur: berkas berpasangan (_down.sql) yang membatalkan langkahnya.
//
// Aturan yang dijaga:
//   ADR-U-0029  nol COMMIT di teks SQL. Perintah DDL Oracle memang menutup
//               transaksinya sendiri, dan itu sifat Oracle - bukan alasan untuk
//               menulis COMMIT di teks SQL. PeriksaSQL menolaknya.
//   ADR-U-0033  nama skema disebut eksplisit; berkas .sql memakai penanda
//               {skema} yang diganti saat dijalankan.
//   ADR-U-0005  menolak berjalan bila lingkungan menunjuk produksi Pega.

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"path"
	"sort"
	"strings"
	"time"
)

//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// ErrMigrasiDiProduksi menolak perubahan bentuk basis data di lingkungan
// produksi Pega. Migrasi membuat dan membuang tabel; ia tidak pernah boleh
// berjalan di sana tanpa keputusan manusia.
var ErrMigrasiDiProduksi = errors.New("repository: menolak migrasi saat IS_PEGA_PROD=true")

// namaTabelMigrasi mencatat langkah mana yang sudah dijalankan.
//
// [usulan] Nama dan bentuk tabel pencatat ini belum pernah diputuskan lewat
// ADR; ia boleh diganti lewat keputusan tertulis.
const namaTabelMigrasi = "T_MIGRASI"

// Migrasi adalah satu langkah, sudah dipecah menjadi pernyataan-pernyataan.
type Migrasi struct {
	Nama       string
	Pernyataan []string
}

// LaporanMigrasi menceritakan apa yang benar-benar terjadi.
type LaporanMigrasi struct {
	Dijalankan []string // langkah yang baru dijalankan kali ini
	Dilewati   []string // langkah yang memang sudah pernah dijalankan
	Pernyataan int      // cacah pernyataan SQL yang dieksekusi
}

// pecahPernyataan memecah isi berkas .sql menjadi pernyataan terpisah.
//
// Pemisahnya adalah baris yang HANYA berisi tanda garis miring - konvensi
// SQL*Plus. Memakai titik koma sebagai pemisah tidak aman karena titik koma
// juga muncul di dalam blok PL/SQL.
func pecahPernyataan(isi string) []string {
	var out []string
	var sekarang []string
	for _, baris := range strings.Split(isi, "\n") {
		if strings.TrimSpace(baris) == "/" {
			if p := gabung(sekarang); p != "" {
				out = append(out, p)
			}
			sekarang = nil
			continue
		}
		sekarang = append(sekarang, baris)
	}
	if p := gabung(sekarang); p != "" {
		out = append(out, p)
	}
	return out
}

// gabung menyatukan baris menjadi satu pernyataan, membuang baris komentar
// murni dan spasi di tepi. Mengembalikan teks kosong bila tidak ada isi.
func gabung(baris []string) string {
	var isi []string
	for _, b := range baris {
		// U+FEFF (byte order mark) menempel di awal berkas bila penyunting
		// menulisnya sebagai UTF-8 ber-BOM. TrimSpace TIDAK membuangnya,
		// sehingga tanpa baris ini satu baris komentar dapat lolos menjadi
		// bagian pernyataan SQL dan Oracle menolaknya. Pernah terjadi di
		// berkas 004 - karena itu dijaga di sini, bukan hanya dibersihkan
		// sekali di berkasnya.
		b = strings.TrimPrefix(b, "\ufeff")
		if strings.HasPrefix(strings.TrimSpace(b), "--") {
			continue
		}
		isi = append(isi, b)
	}
	return strings.TrimSpace(strings.Join(isi, "\n"))
}

// daftarMigrasi membaca seluruh langkah maju, terurut menurut namanya.
func daftarMigrasi(mundur bool) ([]Migrasi, error) {
	entri, err := berkasMigrasi.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("repository: membaca folder migrasi: %w", err)
	}
	var nama []string
	for _, e := range entri {
		n := e.Name()
		if !strings.HasSuffix(n, ".sql") {
			continue
		}
		if strings.HasSuffix(n, "_down.sql") != mundur {
			continue
		}
		nama = append(nama, n)
	}
	sort.Strings(nama)
	if mundur {
		// Jalur mundur berjalan MENURUN supaya anak hilang sebelum induknya.
		for i, j := 0, len(nama)-1; i < j; i, j = i+1, j-1 {
			nama[i], nama[j] = nama[j], nama[i]
		}
	}

	out := make([]Migrasi, 0, len(nama))
	for _, n := range nama {
		isi, err := berkasMigrasi.ReadFile(path.Join("migrations", n))
		if err != nil {
			return nil, fmt.Errorf("repository: membaca %s: %w", n, err)
		}
		out = append(out, Migrasi{Nama: n, Pernyataan: pecahPernyataan(string(isi))})
	}
	return out, nil
}

// kunciLangkah menyamakan nama berkas maju dan mundur menjadi satu kunci,
// supaya jalur mundur tahu langkah mana yang dibatalkannya.
func kunciLangkah(nama string) string {
	return strings.TrimSuffix(strings.TrimSuffix(nama, ".sql"), "_down")
}

// siapkanTabelMigrasi membuat tabel pencatat bila belum ada.
func (d *DB) siapkanTabelMigrasi(ctx context.Context) error {
	tabel, err := d.Qualify(namaTabelMigrasi)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`CREATE TABLE %s (
		NAMA            VARCHAR2(128) NOT NULL,
		DIJALANKAN_PADA DATE,
		CONSTRAINT PK_%s PRIMARY KEY (NAMA))`, tabel, namaTabelMigrasi)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := d.sql.ExecContext(ctx, q); err != nil {
		if sudahAda(err) {
			return nil
		}
		return fmt.Errorf("repository: membuat tabel migrasi: %w", err)
	}
	return nil
}

// sudahAda mengenali galat Oracle "nama sudah dipakai objek lain" (ORA-00955).
// Itu yang membuat pembuatan tabel pencatat menjadi idempoten.
func sudahAda(err error) bool {
	return err != nil && strings.Contains(err.Error(), "ORA-00955")
}

// tidakAda mengenali galat Oracle "tabel atau view tidak ada" (ORA-00942)
// dan "sequence tidak ada" (ORA-02289). Jalur mundur memakainya supaya
// membongkar yang memang belum ada bukan kegagalan.
func tidakAda(err error) bool {
	if err == nil {
		return false
	}
	p := err.Error()
	return strings.Contains(p, "ORA-00942") || strings.Contains(p, "ORA-02289")
}

// sudahDijalankan membaca daftar langkah yang tercatat.
func (d *DB) sudahDijalankan(ctx context.Context) (map[string]bool, error) {
	tabel, err := d.Qualify(namaTabelMigrasi)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT NAMA FROM %s`, tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := d.sql.QueryContext(ctx, q)
	if err != nil {
		if tidakAda(err) {
			return map[string]bool{}, nil
		}
		return nil, fmt.Errorf("repository: membaca catatan migrasi: %w", err)
	}
	defer func() { _ = rows.Close() }()

	out := map[string]bool{}
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		out[n] = true
	}
	return out, rows.Err()
}

// JalankanMigrasi menjalankan seluruh langkah yang belum pernah dijalankan.
//
// Idempoten: langkah yang sudah tercatat dilewati, sehingga memanggilnya dua
// kali berturut-turut aman.
func (d *DB) JalankanMigrasi(ctx context.Context) (LaporanMigrasi, error) {
	var lap LaporanMigrasi
	if d == nil || d.sql == nil {
		return lap, ErrTanpaOracle
	}
	if d.isPegaProd {
		return lap, ErrMigrasiDiProduksi
	}
	if err := d.siapkanTabelMigrasi(ctx); err != nil {
		return lap, err
	}
	selesai, err := d.sudahDijalankan(ctx)
	if err != nil {
		return lap, err
	}
	langkah, err := daftarMigrasi(false)
	if err != nil {
		return lap, err
	}
	tabel, err := d.Qualify(namaTabelMigrasi)
	if err != nil {
		return lap, err
	}

	for _, m := range langkah {
		kunci := kunciLangkah(m.Nama)
		if selesai[kunci] {
			lap.Dilewati = append(lap.Dilewati, kunci)
			continue
		}
		for _, p := range m.Pernyataan {
			q := strings.ReplaceAll(p, "{skema}", d.skema)
			if err := PeriksaSQL(q); err != nil {
				return lap, fmt.Errorf("%s: %w", m.Nama, err)
			}
			if _, err := d.sql.ExecContext(ctx, q); err != nil {
				return lap, fmt.Errorf("repository: migrasi %s: %w", m.Nama, err)
			}
			lap.Pernyataan++
		}
		qCatat := fmt.Sprintf(`INSERT INTO %s (NAMA, DIJALANKAN_PADA) VALUES (:1, :2)`, tabel)
		if err := PeriksaSQL(qCatat); err != nil {
			return lap, err
		}
		if _, err := d.sql.ExecContext(ctx, qCatat, kunci, time.Now()); err != nil {
			return lap, fmt.Errorf("repository: mencatat migrasi %s: %w", kunci, err)
		}
		lap.Dijalankan = append(lap.Dijalankan, kunci)
	}
	return lap, nil
}

// BongkarMigrasi menjalankan jalur mundur untuk langkah yang sudah dijalankan.
//
// Urutannya menurun, sehingga tabel anak dibongkar sebelum induknya. Langkah
// yang objeknya memang sudah tidak ada dilewati tanpa galat.
func (d *DB) BongkarMigrasi(ctx context.Context) (LaporanMigrasi, error) {
	var lap LaporanMigrasi
	if d == nil || d.sql == nil {
		return lap, ErrTanpaOracle
	}
	if d.isPegaProd {
		return lap, ErrMigrasiDiProduksi
	}
	langkah, err := daftarMigrasi(true)
	if err != nil {
		return lap, err
	}
	tabel, err := d.Qualify(namaTabelMigrasi)
	if err != nil {
		return lap, err
	}
	// Hanya langkah yang TERCATAT pernah dijalankan yang dibongkar. Tanpa
	// pembacaan ini, laporan akan menyebut langkah yang tidak berbuat apa-apa
	// sebagai "dijalankan".
	selesai, err := d.sudahDijalankan(ctx)
	if err != nil {
		return lap, err
	}

	for _, m := range langkah {
		kunci := kunciLangkah(m.Nama)
		if !selesai[kunci] {
			lap.Dilewati = append(lap.Dilewati, kunci)
			continue
		}
		for _, p := range m.Pernyataan {
			q := strings.ReplaceAll(p, "{skema}", d.skema)
			if err := PeriksaSQL(q); err != nil {
				return lap, fmt.Errorf("%s: %w", m.Nama, err)
			}
			if _, err := d.sql.ExecContext(ctx, q); err != nil {
				if tidakAda(err) {
					continue
				}
				return lap, fmt.Errorf("repository: bongkar %s: %w", m.Nama, err)
			}
			lap.Pernyataan++
		}
		// DELETE ini pembukuan migrasi, BUKAN jalur pengguna - ADR-U-0031
		// melarang DELETE pada jalur pengguna, bukan pada catatan langkah.
		qHapus := fmt.Sprintf(`DELETE FROM %s WHERE NAMA = :1`, tabel)
		if err := PeriksaSQL(qHapus); err != nil {
			return lap, err
		}
		if _, err := d.sql.ExecContext(ctx, qHapus, kunci); err != nil {
			if !tidakAda(err) {
				return lap, fmt.Errorf("repository: menghapus catatan %s: %w", kunci, err)
			}
		}
		lap.Dijalankan = append(lap.Dijalankan, kunci)
	}
	return lap, nil
}
