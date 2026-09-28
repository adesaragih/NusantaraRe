package repository

// Kotak masuk PremiumList Life - tiket 01 bagian 2.
//
// Meniru `ReportDefinition/InboxPremiumList.xml`, dibaca sebagai pohon
// 28-09-2026. Tiga belas kolom terdaftar di sana; sebelas dapat kami baca,
// dan dua tidak - lihat catatan di bawah, dan tiket 01.
//
//	b721 `A.pyID`                              -> Case ID
//	b735 `A.pxCreateDateTime`                  -> Create Date/Time
//	b750 `A.pxCreateOpName`                    -> Create Operator Name
//	b764 `A.pyStatusWork`                      -> Work Status
//	b778 `A.CedingCoName`                      -> CedingCoName
//	b794 `A.PolicyHolderName`                  -> PolicyHolderName
//	b807 `A.PremiumListSummary.PL_NUMBER`      -> PL_NUMBER       (*)
//	b821 `A.PremiumListSummary.RISLIPRNM`      -> RISLIPRNM
//	b836 `A.Type`                              -> Type
//	b849 `A.MarketingName`                     -> MarketingName
//	b863 `A.SobName`                           -> SobName
//	b878 `A.DateReceived`                      -> DateReceived
//	b891 `A.KetentuanUnderwriting`             -> Ketentuan Underwriting  (**)
//
// ⛔ (*) `PL_NUMBER` ADA DI TINGKAT YANG BERBEDA. Jalur propertinya menyebut
// `PremiumListSummary`, tetapi di skema kami kolomnya hidup di
// `T_PREMIUM_LIST_DETAIL` (migrasi 052) - bukan di `_SUMMARY` (055) maupun
// di `T_PREMIUM_LIST` (051). Dibaca dari sana dengan `MAX`, sebab satu polis
// bernomor PL satu: bila baris detailnya ternyata bernomor berbeda, itu data
// yang rusak dan `PL_NUMBER` bukan tempat menyembunyikannya.
//
// ⛔ (**) `KetentuanUnderwriting` TIDAK PUNYA KOLOM di migrasi 050-056 mana
// pun. Kolomnya karena itu TIDAK ditampilkan, dan ketiadaannya dicatat di
// tiket 01 - bukan diisi teks kosong yang terbaca "memang kosong".
//
// ⛔ Setiap query menyebut skemanya lewat `Qualify` (ADR-U-0033).
//
// Dibaca sesudah: polis_work.go.

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

// BarisInboxPolis adalah satu baris kotak masuk PremiumList.
//
// ⚠️ Nol nama orang. `PolicyHolderName` adalah nama PEMEGANG POLIS, yang di
// modul ini adalah badan usaha - bukan tertanggung. Bila kelak ternyata ia
// dapat berisi nama perorangan, kolom ini yang pertama harus ditinjau.
type BarisInboxPolis struct {
	CaseID           string
	TglCreate        time.Time
	CreateOpName     string
	StatusWork       string
	Position         string
	CedingCoName     string
	PolicyHolderName string
	PLNumber         string
	RISlipRNM        string
	Type             string
	MarketingName    string
	SobName          string
	DateReceived     *time.Time
}

// HalamanInboxPolis adalah satu halaman kotak masuk beserta cacah totalnya.
type HalamanInboxPolis struct {
	Baris []BarisInboxPolis
	// Total adalah cacah SELURUH baris yang cocok - bukan yang di halaman
	// ini. Lencana tab memakainya, dan lencana yang mencacah halaman akan
	// berkata "20" untuk antrean seribu.
	Total int
}

// InboxPolis membaca kotak masuk PremiumList.
type InboxPolis struct{ db *DB }

// NewInboxPolis menyusunnya.
func NewInboxPolis(db *DB) *InboxPolis { return &InboxPolis{db: db} }

// sqlInboxPolis merakit pembacaan satu halaman.
//
// ⛔ `LEFT JOIN` ke detail, bukan `JOIN`: polis yang CSV-nya belum diunggah
// belum punya baris detail sama sekali, dan ia justru yang paling perlu
// tampil di kotak masuk - pekerjaannya belum selesai.
//
// ⚠️ Urutannya `TGL_INPUT DESC` lalu `ID` - terbaru dahulu, dengan pemutus
// seri. Daftar tanpa pemutus seri menampilkan baris yang berpindah sendiri
// di antara dua halaman, dan baris yang berpindah dapat terlewat.
func sqlInboxPolis(work, polis, detail string) string {
	return fmt.Sprintf(`SELECT w.ID, w.STATUS, w.POSITION,
	        p.CREATE_OP_NAME, p.CEDING_CO_NAME, p.POLICY_HOLDER_NAME,
	        p.RI_SLIP_RNM, p.TYPE, p.MARKETING_NAME, p.SOB_NAME,
	        p.DATE_RECEIVED, p.TGL_INPUT,
	        (SELECT MAX(d.PL_NUMBER) FROM %s d
	          WHERE d.PREMIUM_LIST_ID = p.ID) AS PL_NUMBER
	   FROM %s w LEFT JOIN %s p ON p.ID_PEGA = w.ID
	  WHERE (:1 IS NULL OR w.POSITION = :1)
	  ORDER BY p.TGL_INPUT DESC, w.ID
	  OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY`, detail, work, polis)
}

// sqlCacahInboxPolis merakit pencacahnya.
//
// ⛔ Query TERSENDIRI, bukan `COUNT(*) OVER ()` yang ditempel ke query di
// atas. Jendela itu dihitung per baris halaman; pada halaman KOSONG - yang
// terjadi setiap kali seseorang membuka halaman terakhir lalu satu baris
// hilang - ia tidak mengembalikan apa pun, dan totalnya menjadi nol untuk
// antrean yang tidak kosong.
func sqlCacahInboxPolis(work string) string {
	return fmt.Sprintf(
		`SELECT COUNT(*) FROM %s w WHERE (:1 IS NULL OR w.POSITION = :1)`, work)
}

// Ambil membaca satu halaman kotak masuk.
//
// `posisi` kosong berarti SELURUH posisi - padanan tab "semua".
func (r *InboxPolis) Ambil(ctx context.Context, posisi string, halaman, ukuran int) (
	HalamanInboxPolis, error) {

	work, err := r.db.Qualify("T_WORK_POLIS")
	if err != nil {
		return HalamanInboxPolis{}, err
	}
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return HalamanInboxPolis{}, err
	}
	detail, err := r.db.Qualify("T_PREMIUM_LIST_DETAIL")
	if err != nil {
		return HalamanInboxPolis{}, err
	}
	if halaman < 1 {
		halaman = 1
	}
	if ukuran < 1 {
		ukuran = 20
	}

	var hasil HalamanInboxPolis
	qCacah := sqlCacahInboxPolis(work)
	if err := PeriksaSQL(qCacah); err != nil {
		return HalamanInboxPolis{}, err
	}
	if err := r.db.sql.QueryRowContext(ctx, qCacah,
		kosongJadiNil(posisi)).Scan(&hasil.Total); err != nil {
		return HalamanInboxPolis{}, fmt.Errorf("repository: mencacah kotak masuk polis: %w", err)
	}

	q := sqlInboxPolis(work, polis, detail)
	if err := PeriksaSQL(q); err != nil {
		return HalamanInboxPolis{}, err
	}
	baris, err := r.db.sql.QueryContext(ctx, q, kosongJadiNil(posisi),
		(halaman-1)*ukuran, ukuran)
	if err != nil {
		return HalamanInboxPolis{}, fmt.Errorf("repository: membaca kotak masuk polis: %w", err)
	}
	defer baris.Close()

	hasil.Baris = []BarisInboxPolis{}
	for baris.Next() {
		var (
			id, status, posisiBaris, opName, ceding, pemegang sql.NullString
			slip, tipe, marketing, sob, plNomor               sql.NullString
			diterima, dibuat                                  sql.NullTime
		)
		if err := baris.Scan(&id, &status, &posisiBaris, &opName, &ceding,
			&pemegang, &slip, &tipe, &marketing, &sob, &diterima, &dibuat,
			&plNomor); err != nil {
			return HalamanInboxPolis{}, fmt.Errorf("repository: memindai baris polis: %w", err)
		}
		b := BarisInboxPolis{
			CaseID: id.String, StatusWork: status.String, Position: posisiBaris.String,
			CreateOpName: opName.String, CedingCoName: ceding.String,
			PolicyHolderName: pemegang.String, RISlipRNM: slip.String,
			Type: tipe.String, MarketingName: marketing.String, SobName: sob.String,
			PLNumber: plNomor.String,
		}
		if dibuat.Valid {
			b.TglCreate = dibuat.Time
		}
		// ⛔ NULL tetap nil, bukan tanggal nol. Tanggal nol adalah tahun 1
		// Masehi, dan pembaca hilir tidak dapat membedakannya dari kolom
		// yang memang kosong (ADR-U-0027).
		if diterima.Valid {
			t := diterima.Time
			b.DateReceived = &t
		}
		hasil.Baris = append(hasil.Baris, b)
	}
	if err := baris.Err(); err != nil {
		return HalamanInboxPolis{}, fmt.Errorf("repository: membaca kotak masuk polis: %w", err)
	}
	return hasil, nil
}

// UkuranHalamanPolisBawaan dipakai bila pemanggil tidak menyebut ukurannya.
//
// ⚠️ `[keputusan kami]` Nol rule menyebut ukuran halaman kotak masuk
// PremiumList; `InboxPremiumList.xml` menyerahkannya ke grid bawaan Pega.
// Dua puluh sejajar dengan kotak masuk Claim Life, supaya kedua layar
// berperilaku sama.
const UkuranHalamanPolisBawaan = 20

// BatasUkuranHalamanPolis menjepit ukuran yang diminta klien.
//
// ⛔ `?ukuran=1000000` atas tabel polis adalah permintaan yang memuat
// seluruh tabel ke memori satu proses - pelajaran yang sama dengan
// `models.BatasPenyakit`.
func BatasUkuranHalamanPolis(diminta int) int {
	const maksimum = 200
	switch {
	case diminta <= 0:
		return UkuranHalamanPolisBawaan
	case diminta > maksimum:
		return maksimum
	}
	return diminta
}

// NomorHalamanPolis menjepit nomor halaman.
func NomorHalamanPolis(diminta string) int {
	n, err := strconv.Atoi(diminta)
	if err != nil || n < 1 {
		return 1
	}
	return n
}
