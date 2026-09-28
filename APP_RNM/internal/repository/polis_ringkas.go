package repository

// Ringkasan polis untuk Claim Life - butir pl4, dipakai butir av.
//
// Untuk apa berkas ini: memberi Claim Life apa yang di Pega disebut
// `pyWorkPage.PolicyDataLife.*` - data polis yang TERISI SENDIRI saat nomor
// polis dipilih, dan di layar Register maupun Outstanding berstatus
// `pyReadOnly`.
//
// ⛔ CLAIM LIFE TIDAK MEMBACA CERMIN JSON-NYA `[keputusan work owner, av]`.
// Sumbernya tabel relasional modul PremiumList Life, bukan `JSON_POLIS`.
//
// ⛔ VERSI BERJALAN ADALAH `PROD_KE` TERBESAR. Migrasi 051: satu baris = satu
// VERSI polis, dan seluruh versi hidup berdampingan (new business maupun
// endorsement). Membaca baris mana saja yang nomornya cocok akan memberi
// Claim Life data polis versi LAMA - dan klaim yang dinilai dengan data versi
// lama adalah klaim yang dinilai salah, tanpa satu pun galat.
//
// ⛔ Setiap query menyebut skemanya lewat `Qualify` (ADR-U-0033).
//
// Dibaca sesudah: polis_inbox.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// ErrPolisNomorTakDitemukan - tidak ada polis dengan nomor itu.
var ErrPolisNomorTakDitemukan = errors.New(
	"repository: nomor polis tidak ditemukan di PremiumList Life")

// PolisRingkas adalah sepuluh medan `PolicyDataLife` yang PUNYA kolom.
//
// ⚠️ TIGA MEDAN LAYAR LAMA TIDAK ADA DI SINI, dan ketiadaannya dinyatakan
// alih-alih diisi teks kosong: `TanggalRespon`, `TanggalKonfirmasi`, dan
// `TanggalRealisasi` tidak punya kolom di migrasi 050-056 mana pun. Ketiganya
// properti halaman kerja Pega; di mana nilainya tinggal `[terbuka]`.
type PolisRingkas struct {
	NomorPolis       string
	Type             string
	MarketingName    string
	CedingCoName     string
	PolicyHolderName string
	BusinessName     string
	// DateReceived adalah `Date Received Email` di layar Claim Life.
	DateReceived *string
	Status       string
	StatusUpdate string
	// ProductNameID adalah KUNCI ambang batas hari (butir ba) - ia yang
	// dipakai membaca `PRODUCTINWARD_LIFE`.
	ProductNameID string
	ProductName   string
	// ProdKe adalah versi yang terbaca, supaya pemanggil dapat menyatakannya.
	ProdKe int

	// Tujuh medan layar lagi - butir av-2 (GILIRAN-11 paket 2), urut section.
	//
	//	TypeCeding/TypeCedingName `.TypeCeding`        "System Reinsurance" b9301
	//	ProRateType               `.ProRateType`       "Premium Method"     b9694
	//	WPC                       `.WPC`               "WPC"                b10149
	//	RetroName                 `.RetroName`         "Retro Name"         b10335 (TP/TR)
	//	SecurityReinsurer         `.SecurityReinsurer` "Security Reinsurer" b10617 (TP/TR)
	//	SobName                   `.SobName`           "SOB"                b11930
	TypeCeding        string
	TypeCedingName    string
	ProRateType       string
	WPC               *string
	RetroName         string
	SecurityReinsurer string
	SobName           string

	// Empat KUNCI yang `Save to RNM` pakai - bukan medan layar, sehingga
	// TIDAK menyeberang ke JSON `PolicyDataLife` (kontraknya dikunci dua sisi).
	//
	//	BusinessCode        `PolicyDataLife.BusinessCode`  langkah 10 (ContentNote), 13-20
	//	CedingCo            `PolicyDataLife.CedingCo`      langkah 11.1 `CARI1` (klaim ganda)
	//	RetroID             `PolicyDataLife.RetroID`       langkah 27
	//	SecurityReinsurerID `PolicyDataLife.SecurityReinsurerID` langkah 27
	BusinessCode        string
	CedingCo            string
	RetroID             string
	SecurityReinsurerID string
}

// RingkasPolisLife membaca ringkasan polis.
type RingkasPolisLife struct{ db *DB }

// NewRingkasPolisLife menyusunnya.
func NewRingkasPolisLife(db *DB) *RingkasPolisLife { return &RingkasPolisLife{db: db} }

// sqlPolisRingkas merakit pembacaan versi BERJALAN sebuah nomor polis.
//
// ⛔ `ORDER BY PROD_KE DESC` lalu `FETCH FIRST 1 ROWS ONLY`. Tanpa urutan itu
// Oracle bebas memberi versi mana pun yang paling murah dibacanya - dan yang
// paling murah bukan yang paling baru.
//
// ⚠️ `NVL(PROD_KE, 0)` supaya baris new business yang belum pernah
// di-endorse - `PROD_KE` masih kosong - tetap ikut terurut, bukan terlempar
// ke ujung oleh NULL.
func sqlPolisRingkas(polis string) string {
	return fmt.Sprintf(`SELECT p.NO_POLIS, p.TYPE, p.MARKETING_NAME,
	        p.CEDING_CO_NAME, p.POLICY_HOLDER_NAME, p.BUSINESS_NAME,
	        TO_CHAR(p.DATE_RECEIVED, 'YYYY-MM-DD'),
	        p.STATUSS, p.STATUS_UPDATE, p.PRODUCT_NAME_ID, p.PRODUCT_NAME,
	        NVL(p.PROD_KE, 0),
	        p.BUSINESS_CODE, p.CEDING_CO, p.RETRO_ID, p.SECURITY_REINSURER_ID,
	        p.TYPE_CEDING, p.TYPE_CEDING_NAME, p.PRO_RATE_TYPE,
	        TO_CHAR(p.WPC, 'YYYY-MM-DD'), p.RETRO_NAME, p.SECURITY_REINSURER, p.SOB_NAME
	   FROM %s p
	  WHERE p.NO_POLIS = :1
	  ORDER BY NVL(p.PROD_KE, 0) DESC
	  FETCH FIRST 1 ROWS ONLY`, polis)
}

// Ringkas membaca versi berjalan sebuah nomor polis.
func (r *RingkasPolisLife) Ringkas(ctx context.Context, nomorPolis string) (
	PolisRingkas, error) {

	if strings.TrimSpace(nomorPolis) == "" {
		return PolisRingkas{}, fmt.Errorf("%w: nomor polis kosong",
			ErrPolisNomorTakDitemukan)
	}
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return PolisRingkas{}, err
	}
	q := sqlPolisRingkas(polis)
	if err := PeriksaSQL(q); err != nil {
		return PolisRingkas{}, err
	}
	var (
		no, tipe, marketing, ceding, pemegang, bisnis sql.NullString
		diterima, status, statusUbah, prodID, prodNm  sql.NullString
		kodeBisnis, cedingKode, retroID, secID        sql.NullString
		tipeCeding, tipeCedingNm, prorata, wpc        sql.NullString
		retroNm, secNm, sobNm                         sql.NullString
		prodKe                                        int
	)
	if err := r.db.sql.QueryRowContext(ctx, q, nomorPolis).Scan(
		&no, &tipe, &marketing, &ceding, &pemegang, &bisnis,
		&diterima, &status, &statusUbah, &prodID, &prodNm, &prodKe,
		&kodeBisnis, &cedingKode, &retroID, &secID,
		&tipeCeding, &tipeCedingNm, &prorata, &wpc, &retroNm, &secNm, &sobNm); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PolisRingkas{}, fmt.Errorf("%w: %q", ErrPolisNomorTakDitemukan, nomorPolis)
		}
		return PolisRingkas{}, fmt.Errorf("repository: membaca ringkasan polis: %w", err)
	}
	hasil := PolisRingkas{
		NomorPolis:       no.String,
		Type:             tipe.String,
		MarketingName:    marketing.String,
		CedingCoName:     ceding.String,
		PolicyHolderName: pemegang.String,
		BusinessName:     bisnis.String,
		Status:           status.String,
		StatusUpdate:     statusUbah.String,
		ProductNameID:    prodID.String,
		ProductName:      prodNm.String,
		ProdKe:           prodKe,

		BusinessCode:        strings.TrimSpace(kodeBisnis.String),
		CedingCo:            strings.TrimSpace(cedingKode.String),
		RetroID:             strings.TrimSpace(retroID.String),
		SecurityReinsurerID: strings.TrimSpace(secID.String),

		TypeCeding:        tipeCeding.String,
		TypeCedingName:    tipeCedingNm.String,
		ProRateType:       prorata.String,
		RetroName:         retroNm.String,
		SecurityReinsurer: secNm.String,
		SobName:           sobNm.String,
	}
	if wpc.Valid {
		t := wpc.String
		hasil.WPC = &t
	}
	// ⛔ NULL tetap nil, bukan tanggal kosong berbentuk teks: kolom yang belum
	// diisi dan kolom bertanggal adalah dua keadaan berbeda.
	if diterima.Valid {
		t := diterima.String
		hasil.DateReceived = &t
	}
	return hasil, nil
}
