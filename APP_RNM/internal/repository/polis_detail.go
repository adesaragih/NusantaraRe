package repository

// Pembacaan grid Premium List Detail - tiket 03 PremiumList Life.
//
// Untuk apa berkas ini: membaca baris peserta satu polis, kolomnya diturunkan
// dari `models.KolomGridPeserta` - daftar yang sama yang dipakai layar.
//
// ⛔ SQL-NYA DIRAKIT DARI DAFTAR KOLOM, tidak diketik ulang. Query yang
// diketik terpisah dari daftar kolomnya akan menyimpang darinya, dan
// penyimpangannya muncul di layar sebagai angka di bawah judul yang salah -
// kekeliruan yang tidak satu pun galat tunjukkan.
//
// ⛔ SELURUH ANGKA KELUAR SEBAGAI TEKS lewat `TO_CHAR(..., 'TM9', ...)`
// (ADR-U-0003, ADR-U-0016). Nol `float64` di jalur ini.
//
// ⛔ Setiap query menyebut skemanya lewat `Qualify` (ADR-U-0033).
//
// Dibaca sesudah: polis_nomor.go.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/internal/models"
	"nusantarare/inti/db"
)

// ⛔ PEMBUNGKUSNYA DIPINJAM, tidak diketik ulang: `fmtDesimal`
// (klaimlife.go) dan `fmtTanggalOracle` (barislamakolom.go), keduanya di
// paket yang SAMA. Ronde pertama berkas ini menyalin keduanya dengan alasan
// "kebetulan sama, bukan satu aturan yang dibagi" - dan itu keliru:
// ADR-U-0003 dan ADR-U-0016 menjadikannya SATU aturan. Dua salinan satu
// aturan adalah satu aturan yang suatu hari hanya diperbaiki di satu tempat.

// ekspresiKolomPeserta merakit daftar kolom SELECT dari daftar kolom grid.
func ekspresiKolomPeserta(kolom []models.KolomPeserta) string {
	bagian := make([]string, 0, len(kolom)+1)
	bagian = append(bagian, "d.ID")
	for _, k := range kolom {
		medan := "d." + k.Nama
		switch k.Jenis {
		case models.KolomPesertaAngka:
			bagian = append(bagian, fmt.Sprintf(db.FmtDesimal, medan))
		case models.KolomPesertaTanggal:
			bagian = append(bagian, fmt.Sprintf(db.FmtTanggalOracle, medan))
		default:
			bagian = append(bagian, medan)
		}
	}
	return strings.Join(bagian, ", ")
}

// sqlGridPeserta merakit pembacaan satu halaman grid.
//
// ⚠️ Urutannya `CERTIFICATE_NO` lalu `ID` - nomor sertifikat yang dilihat
// orang, dengan pemutus seri supaya baris tidak berpindah sendiri di antara
// dua halaman.
func sqlGridPeserta(detail string, kolom []models.KolomPeserta) string {
	return fmt.Sprintf(`SELECT %s FROM %s d
	  WHERE d.PREMIUM_LIST_ID = :1
	  ORDER BY d.CERTIFICATE_NO, d.ID
	  OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY`,
		ekspresiKolomPeserta(kolom), detail)
}

// sqlCacahPeserta merakit pencacahnya - query tersendiri, bukan jendela.
func sqlCacahPeserta(detail string) string {
	return fmt.Sprintf(
		`SELECT COUNT(*) FROM %s d WHERE d.PREMIUM_LIST_ID = :1`, detail)
}

// HalamanPeserta adalah satu halaman grid beserta cacah totalnya.
type HalamanPeserta struct {
	Baris []models.BarisPeserta
	Total int
}

// GridPeserta membaca baris peserta satu polis.
type GridPeserta struct{ db *db.DB }

// NewGridPeserta menyusunnya.
func NewGridPeserta(db *db.DB) *GridPeserta { return &GridPeserta{db: db} }

// BatasUkuranHalamanPeserta menjepit ukuran yang diminta klien.
//
// ⛔ Satu polis grup dapat memuat RIBUAN peserta (migrasi 052). Permintaan
// tanpa batas atas memuat seluruhnya ke memori satu proses.
func BatasUkuranHalamanPeserta(diminta int) int {
	const (
		bawaan   = 50
		maksimum = 500
	)
	switch {
	case diminta <= 0:
		return bawaan
	case diminta > maksimum:
		return maksimum
	}
	return diminta
}

// Ambil membaca satu halaman peserta.
func (r *GridPeserta) Ambil(ctx context.Context, polisID string, halaman, ukuran int) (
	HalamanPeserta, error) {

	detail, err := r.db.Qualify("T_PREMIUM_LIST_DETAIL")
	if err != nil {
		return HalamanPeserta{}, err
	}
	if halaman < 1 {
		halaman = 1
	}
	ukuran = BatasUkuranHalamanPeserta(ukuran)

	var hasil HalamanPeserta
	qCacah := sqlCacahPeserta(detail)
	if err := db.PeriksaSQL(qCacah); err != nil {
		return HalamanPeserta{}, err
	}
	if err := r.db.QueryRowContext(ctx, qCacah, polisID).Scan(&hasil.Total); err != nil {
		return HalamanPeserta{}, fmt.Errorf("repository: mencacah peserta polis: %w", err)
	}

	kolom := models.KolomGridPeserta
	q := sqlGridPeserta(detail, kolom)
	if err := db.PeriksaSQL(q); err != nil {
		return HalamanPeserta{}, err
	}
	baris, err := r.db.QueryContext(ctx, q, polisID, (halaman-1)*ukuran, ukuran)
	if err != nil {
		return HalamanPeserta{}, fmt.Errorf("repository: membaca peserta polis: %w", err)
	}
	defer baris.Close()

	hasil.Baris = []models.BarisPeserta{}
	for baris.Next() {
		// ⛔ Sasaran pindai dirakit dari daftar kolom yang SAMA dengan yang
		// merakit SELECT-nya. Dua daftar yang panjangnya berbeda adalah galat
		// `Scan` saat berjalan; satu daftar membuatnya mustahil.
		sel := make([]sql.NullString, len(kolom)+1)
		sasaran := make([]any, len(sel))
		for i := range sel {
			sasaran[i] = &sel[i]
		}
		if err := baris.Scan(sasaran...); err != nil {
			return HalamanPeserta{}, fmt.Errorf("repository: memindai peserta polis: %w", err)
		}
		b := models.BarisPeserta{
			ID:    sel[0].String,
			Nilai: make(map[string]string, len(kolom)),
		}
		for i, k := range kolom {
			// ⛔ NULL menjadi teks KOSONG, dan bukan "0".
			//
			// ⚠️ RUJUKAN ADR-U-0027 DICABUT dari sini 28-09-2026: ADR itu
			// tentang kolom yang dideklarasi nullable dan kewajiban isi yang
			// ditegakkan di kode - ia tidak berkata apa pun tentang pemetaan
			// NULL menjadi teks. Aturannya berdiri sendiri, dan sebabnya
			// cukup: premi yang BELUM DIISI dan premi yang MEMANG NOL adalah
			// dua keadaan berbeda, dan hanya satu di antaranya perlu
			// dikerjakan orang. Memetakan keduanya ke "0" membuang perbedaan
			// itu di tempat yang tidak dapat dipulihkan.
			if sel[i+1].Valid {
				b.Nilai[k.Nama] = sel[i+1].String
			} else {
				b.Nilai[k.Nama] = ""
			}
		}
		hasil.Baris = append(hasil.Baris, b)
	}
	if err := baris.Err(); err != nil {
		return HalamanPeserta{}, fmt.Errorf("repository: membaca peserta polis: %w", err)
	}
	return hasil, nil
}
