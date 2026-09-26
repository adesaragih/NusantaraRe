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
	"regexp"
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
	// Pernyataan mencacah yang benar-benar DIEKSEKUSI. Pernyataan yang
	// dilewati karena objeknya sudah ada tidak ikut dihitung; ia muncul
	// di ObjekSudahAda, dan memang tidak dikirim ke Oracle.
	Pernyataan int
	// ObjekSudahAda mencatat pernyataan CREATE yang dilewati karena objeknya
	// ternyata sudah berdiri. Ia TIDAK kosong hanya pada keadaan tidak normal -
	// lihat JalankanMigrasi - sehingga pemanggil yang menemukannya berisi tahu
	// bahwa ada langkah yang dulu gagal separuh jalan.
	ObjekSudahAda []string
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

// PernyataanLangkah mengembalikan pernyataan SQL satu langkah migrasi, apa
// adanya, dengan penanda {skema} MASIH utuh.
//
// Dipakai test db supaya ia dapat menjalankan pernyataan yang BENAR-BENAR ada
// di langkah itu - bukan pernyataan tiruan yang bentuknya berbeda. Menirukan
// kegagalan separuh jalan dengan objek berbentuk lain akan menguji hal lain.
func PernyataanLangkah(nama string) ([]string, error) {
	isi, err := berkasMigrasi.ReadFile(path.Join("migrations", nama))
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s: %w", nama, err)
	}
	return pecahPernyataan(string(isi)), nil
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

// sudahAda mengenali SATU galat Oracle: nama yang hendak dipakai sudah dipakai
// objek lain - tabel, index, atau sequence (ORA-00955). Itu yang membuat
// pembuatan objek menjadi idempoten.
//
// ⛔ ORA-02264 sengaja TIDAK termasuk, dan ini ralat atas ronde 2.
// ORA-02264 berarti "nama itu sudah dipakai constraint lain", dan Oracle baru
// memeriksanya ketika TABELNYA BELUM ADA - bila tabelnya sudah ada, ia menjawab
// ORA-00955 lebih dulu. Jadi ORA-02264 pada sebuah CREATE TABLE berarti tabel
// itu JUSTRU TIDAK terbuat. Menelannya berarti mencatat langkah sebagai sukses
// atas tabel yang tidak pernah ada - kegagalan yang baru ketahuan jauh di hilir.
func sudahAda(err error) bool {
	return err != nil && strings.Contains(err.Error(), "ORA-00955")
}

// polaObjekDibuat menangkap nama objek yang dibuat sebuah pernyataan CREATE.
var polaObjekDibuat = regexp.MustCompile(
	`(?is)CREATE\s+(?:UNIQUE\s+)?(?:TABLE|INDEX|SEQUENCE)\s+\{skema\}\.(\w+)`)

// namaObjekDibuat membaca nama objek dari pernyataan CREATE, atau teks kosong
// bila pernyataannya bukan CREATE yang dikenali.
func namaObjekDibuat(pernyataan string) string {
	m := polaObjekDibuat.FindStringSubmatch(pernyataan)
	if m == nil {
		return ""
	}
	return strings.ToUpper(m[1])
}

// objekAda bertanya ke katalog Oracle apakah sebuah objek benar-benar berdiri.
//
// Dipakai sesudah sebuah CREATE dilewati karena menjawab ORA-00955. Oracle
// hanya bilang "nama sudah dipakai"; ia tidak bilang objek APA. Tanpa
// pemeriksaan ini, nama yang dipakai objek jenis lain akan lolos sebagai
// "sudah ada".
//
// ⚠️ Yang diperiksa keberadaan, BUKAN bentuk. Kolom yang berbeda tidak
// tertangkap di sini; yang menjaga hal itu jalur mundur - bongkar dulu, baru
// pasang lagi.
//
// SYS.ALL_OBJECTS ditulis berawalan skema, bukan telanjang (ADR-U-0033). Ia
// dipilih di atas USER_OBJECTS karena skema sasaran belum tentu sama dengan
// pengguna sambungan, dan di atas DBA_OBJECTS karena yang terakhir menuntut hak
// istimewa yang tidak perlu. ALL_OBJECTS memperlihatkan objek yang terlihat oleh
// sesi - dan sesi ini baru saja mencoba membuat objek di skema itu, jadi ia
// memang punya akses ke sana.
func (d *DB) objekAda(ctx context.Context, nama string) (bool, error) {
	// UPPER di kedua sisi: Oracle menyimpan pengenal tanpa kutip dalam huruf
	// besar, tetapi ORACLE_SCHEMA datang dari env var dan boleh ditulis
	// bagaimana saja. Tanpa ini, skema yang ditulis huruf kecil membuat
	// pembuktian menjawab "tidak ada" dan MENGGAGALKAN migrasi yang sehat.
	q := `SELECT COUNT(*) FROM SYS.ALL_OBJECTS
	        WHERE UPPER(OWNER) = UPPER(:1) AND UPPER(OBJECT_NAME) = UPPER(:2)`
	if err := PeriksaSQL(q); err != nil {
		return false, err
	}
	var n int
	if err := d.sql.QueryRowContext(ctx, q, d.skema, nama).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: memeriksa keberadaan %s: %w", nama, err)
	}
	return n > 0, nil
}

// pernyataanBuat menjawab apakah sebuah pernyataan SQL membuat objek baru.
//
// Hanya pernyataan CREATE yang boleh dilewati saat objeknya sudah ada. ALTER,
// INSERT, atau DROP yang menjawab galat serupa BUKAN keadaan yang sama, dan
// menelannya akan menyembunyikan kerusakan sungguhan.
func pernyataanBuat(q string) bool {
	return strings.HasPrefix(strings.ToUpper(strings.TrimSpace(q)), "CREATE")
}

// ringkasPernyataan mengambil beberapa kata pertama sebuah pernyataan, supaya
// laporan menyebut objek yang dimaksud tanpa menyalin seluruh DDL-nya.
func ringkasPernyataan(q string) string {
	kata := strings.Fields(q)
	if len(kata) > 4 {
		kata = kata[:4]
	}
	return strings.Join(kata, " ")
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
				// DDL Oracle menutup transaksinya sendiri. Bila sebuah langkah
				// dulu gagal di pernyataan kedua, pernyataan pertamanya sudah
				// terlanjur jadi dan T_MIGRASI tidak sempat mencatat apa pun -
				// sehingga percobaan berikutnya mengulang langkah itu dari awal
				// dan mati di ORA-00955. Objek yang sudah berdiri dilewati,
				// dicatat, lalu langkahnya diteruskan sampai tuntas.
				// ⚠️ Risiko yang DITERIMA sadar: Oracle hanya bilang "nama itu
				// sudah dipakai", bukan "objeknya berbentuk sama". Objek lama
				// yang berbeda bentuk karena itu ikut diterima. Yang menjaga
				// hal itu bukan kode ini melainkan jalur mundur - bongkar dulu,
				// baru pasang lagi - dan itulah yang dilakukan skema uji.
				if pernyataanBuat(q) && sudahAda(err) {
					// ⛔ Dilewati TIDAK cukup - keberadaannya dibuktikan.
					// Oracle hanya bilang namanya terpakai, bukan bahwa objek
					// yang kita maksud ada. Bila ternyata tidak ada, langkah
					// ini GAGAL dan tidak dicatat di T_MIGRASI.
					nama := namaObjekDibuat(p)
					if nama == "" {
						return lap, fmt.Errorf(
							"repository: migrasi %s: dilaporkan sudah ada, nama objeknya tidak terbaca: %w",
							m.Nama, err)
					}
					ada, errPeriksa := d.objekAda(ctx, nama)
					if errPeriksa != nil {
						return lap, fmt.Errorf("repository: migrasi %s: %w",
							m.Nama, errPeriksa)
					}
					if !ada {
						return lap, fmt.Errorf(
							"repository: migrasi %s: %s dilaporkan sudah ada, tidak ditemukan di katalog: %w",
							m.Nama, nama, err)
					}
					lap.ObjekSudahAda = append(lap.ObjekSudahAda,
						fmt.Sprintf("%s: %s", m.Nama, ringkasPernyataan(q)))
					continue
				}
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
