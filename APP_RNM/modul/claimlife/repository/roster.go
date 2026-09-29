package repository

// Pembaca roster Komite - butir af, A2.
//
// Untuk apa berkas ini: membaca anggota komite yang berhak memutuskan sebuah
// nilai klaim, dari `POOLDATA.EMAILKOMITE`.
//
// ⛔ DATA ORANG. Tabel ini memuat nama, email, dan jabatan. Aturannya keras:
// dibaca SAAT JALAN saja; nol baris disalin ke fixture, tiket, maupun log;
// dan hanya kolom yang benar-benar dipakai yang diambil - kolom yang dibaca
// cenderung ikut tercatat di suatu tempat pada akhirnya.
//
// `[terverifikasi]` filternya persis `Claim Life/ReportDefinition/
// FilterEmailKomiteWithLimit.xml`: pecahan baris 662 `pyFilterLogic =
// A AND C AND B`; A (670-680) `.LIMIT_BOTTOM <= Param.LIMIT_BOTTOM`;
// C (683-697) `.STS_KLAIM = Param.STS_KLAIM`; B (701-715) `.STS_AKTIF = "1"`.
// Urutannya `.DEGREE` `ASC` (baris 824-831, `pySortOrder` 1).
//
// Dibaca sesudah: nomorakseptasi.go.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/db"
)

// BarisRoster adalah satu anggota komite yang berhak.
//
// ⚠️ `Email` DATA ORANG. Ia menyeberang ke pengirim email dan ke
// `T_KOMITE_KOMITELIST`, dan tidak ke mana pun lagi.
type BarisRoster struct {
	// OperatorID mengisi `KomiteList(...).KomiteID`.
	//
	// ⛔ NAMA KOLOMNYA MENIPU, dan ini terbukti dari XML - bukan tebakan:
	// `[terverifikasi]` `GetListKomiteLife.xml` pecahan baris 1151 (langkah
	// 6.1) menyetel `.KomiteID = .OPERATOR_ID`. (Salinannya di
	// `CreateKMTLife_Act` 3.1 ter-remark - ralat sensus 28-09-2026.) Yang
	// mengisi "KomiteID" adalah pengenal OPERATOR, bukan kolom `ID` roster.
	OperatorID string
	// Jabatan mengisi `KomiteList(...).IDKomite`.
	//
	// ⛔ Menipu pula, dan ke arah sebaliknya: baris 1237 menyetel
	// `.IDKomite = .JABATAN`. Yang mengisi "IDKomite" adalah JABATAN.
	Jabatan     string
	Email       string
	Degree      int
	LimitBottom string
}

// sqlRosterKomite merakit pernyataannya.
//
// Dipisah supaya bentuknya - ketiga filter dan urutannya - dapat diuji tanpa
// Oracle.
//
// ⛔ `LIMIT_BOTTOM` dibandingkan sebagai ANGKA lewat `TO_NUMBER`, sebab
// kolomnya `INTEGER` `[data DBA]` sedangkan ambang kita teks desimal. Tanpa
// itu Oracle membandingkan dua tipe berbeda dan hasilnya bergantung konversi
// implisit.
func sqlRosterKomite(tabel string) string {
	return fmt.Sprintf(`SELECT OPERATOR_ID, JABATAN, EMAIL, DEGREE, LIMIT_BOTTOM
		  FROM %s
		 WHERE LIMIT_BOTTOM <= TO_NUMBER(:1)
		   AND STS_KLAIM = :2
		   AND STS_AKTIF = :3
		 ORDER BY DEGREE ASC`, tabel)
}

// stsAktifRoster adalah nilai `.STS_AKTIF` yang XML tuntut.
//
// `[terverifikasi]` `FilterEmailKomiteWithLimit.xml` baris 703:
// `<pyFilterValue>"1"</pyFilterValue>` - TEKS, bukan bilangan (ADR-U-0022).
const stsAktifRoster = "1"

// AmbilRosterKomite membaca anggota yang menutup sebuah ambang.
//
// ⚠️ `ambang` diterima sebagai TEKS desimal, bukan float: ia nilai klaim, dan
// uang tidak pernah menjadi float (ADR-U-0003).
func (r *PohonKlaim) AmbilRosterKomite(ctx context.Context,
	ambang, stsKlaim string) ([]BarisRoster, error) {

	tabel, err := r.db.Qualify("EMAILKOMITE")
	if err != nil {
		return nil, err
	}
	q := sqlRosterKomite(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, q, ambang, stsKlaim, stsAktifRoster)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca roster komite: %w", err)
	}
	defer func() { _ = baris.Close() }()

	var hasil []BarisRoster
	for baris.Next() {
		var (
			operator, jabatan, email, limit sql.NullString
			degree                          sql.NullInt64
		)
		if err := baris.Scan(&operator, &jabatan, &email, &degree, &limit); err != nil {
			return nil, fmt.Errorf("repository: membaca baris roster: %w", err)
		}
		hasil = append(hasil, BarisRoster{
			OperatorID:  operator.String,
			Jabatan:     jabatan.String,
			Email:       email.String,
			Degree:      int(degree.Int64),
			LimitBottom: limit.String,
		})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: menelusuri roster: %w", err)
	}
	return hasil, nil
}
