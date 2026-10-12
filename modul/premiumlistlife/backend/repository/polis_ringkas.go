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

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
)

// RingkasPolisLife membaca ringkasan polis.
type RingkasPolisLife struct{ db *db.DB }

// NewRingkasPolisLife menyusunnya.
func NewRingkasPolisLife(db *db.DB) *RingkasPolisLife { return &RingkasPolisLife{db: db} }

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
	return fmt.Sprintf(`SELECT `+kolomPolisRingkas("p.STATUSS")+`
	   FROM %s p
	  WHERE p.NO_POLIS = :1
	  ORDER BY NVL(p.PROD_KE, 0) DESC
	  FETCH FIRST 1 ROWS ONLY`, polis)
}

// kolomPolisRingkas - daftar pilih SATU sumber untuk sqlPolisRingkas dan sqlPolisRingkasIdPega (urutan = bacaRingkas);
// `status` = ekspresi kolom status (Ringkas lama: `p.STATUSS` apa adanya).
func kolomPolisRingkas(status string) string {
	return `p.NO_POLIS, p.TYPE, p.MARKETING_NAME,
	        p.CEDING_CO_NAME, p.POLICY_HOLDER_NAME, p.BUSINESS_NAME,
	        TO_CHAR(p.DATE_RECEIVED, 'YYYY-MM-DD'),
	        ` + status + `, p.STATUS_UPDATE, p.PRODUCT_NAME_ID, p.PRODUCT_NAME,
	        NVL(p.PROD_KE, 0),
	        p.BUSINESS_CODE, p.CEDING_CO, p.RETRO_ID, p.SECURITY_REINSURER_ID,
	        p.TYPE_CEDING, p.TYPE_CEDING_NAME, p.PRO_RATE_TYPE,
	        TO_CHAR(p.WPC, 'YYYY-MM-DD'), p.RETRO_NAME, p.SECURITY_REINSURER, p.SOB_NAME`
}

// sqlPolisRingkasIdPega - header polis dari IDPEGA peserta (keputusan work owner 10-10-2026, perbaikan Claim Life:
// "ambil ke tabel T_PREMIUM_LIST berdasarkan idpega peserta nya saja").
//
// ⛔ `M_LIFE_PREMIUM_DETAIL.IDPEGA` = pengenal kasus (`pzInsKey` -> `NBLF-`/`EDML-`, polis_warisan.go b79; EDM Confirm
// menulis seluruh peserta versi itu ber-IDPEGA kasus EDM, endorsementlife edm_putusan.go), dan baris UTAMA
// T_PREMIUM_LIST kasus itu ber-`ID = ID_PEGA = pengenal kasus` (polis_barisstatus.go). Baris status penawaran juga
// ber-`ID_PEGA` = kasus tetapi `ID`-nya lain - karena itu `p.ID = :1 AND p.ID_PEGA = p.ID`: tepat baris utama, lewat
// PK. NO_POLIS / PL tidak dipakai, jadi NB yang belum pernah di-endorse (NO_POLIS NULL) tetap terbaca dan kasus EDM
// terbuka tidak mungkin terbawa (pesertanya belum pernah ditulis ke M_LIFE). Urutan PROD_KE hanya pengaman.
//
// STATUS EFEKTIF (keputusan work owner 10-10-2026 "polis harus sudah resolve completed"): kasus NB menyimpan status
// kerjanya di `T_WORK_POLIS.STATUS_WORK` (models.StatusPolisSelesai "Resolved-Completed"; `T_PREMIUM_LIST.STATUSS` tidak
// ditulis PremiumList), sedangkan kasus EDM (tanpa baris T_WORK_POLIS) menulis `T_PREMIUM_LIST.STATUSS` saat Confirm
// (endorsementlife `edm_putusan.go` b115, "Resolved-Completed"). Kolom status jawaban = NVL keduanya. Ringkas(NO_POLIS)
// lama tetap `p.STATUSS` apa adanya.
func sqlPolisRingkasIdPega(polis, kerja string) string {
	status := `NVL((SELECT w.STATUS_WORK FROM ` + kerja + ` w WHERE w.ID = p.ID), p.STATUSS)`
	return fmt.Sprintf(`SELECT `+kolomPolisRingkas(status)+`
	   FROM %s p
	  WHERE p.ID = :1 AND p.ID_PEGA = p.ID
	  ORDER BY NVL(p.PROD_KE, 0) DESC
	  FETCH FIRST 1 ROWS ONLY`, polis)
}

// NomorPolisDariID membaca `NO_POLIS` satu polis dari pengenal work-nya
// (`T_PREMIUM_LIST.ID` = `T_WORK_POLIS.ID`) - parameter `noPolis`
// `convertJsonNusareToProduction` (OQ-PL-14, GILIRAN-17).
func (r *RingkasPolisLife) NomorPolisDariID(ctx context.Context, polisID string) (string, error) {
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf(`SELECT p.NO_POLIS FROM %s p WHERE p.ID = :1`, polis)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var no sql.NullString
	if err := r.db.QueryRowContext(ctx, q, polisID).Scan(&no); err != nil {
		return "", fmt.Errorf("repository: membaca nomor polis %q: %w", polisID, err)
	}
	return strings.TrimSpace(no.String), nil
}

// Ringkas membaca versi berjalan sebuah nomor polis.
func (r *RingkasPolisLife) Ringkas(ctx context.Context, nomorPolis string) (
	kontrak.PolisRingkas, error) {

	if strings.TrimSpace(nomorPolis) == "" {
		return kontrak.PolisRingkas{}, fmt.Errorf("%w: nomor polis kosong",
			kontrak.ErrPolisNomorTakDitemukan)
	}
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return kontrak.PolisRingkas{}, err
	}
	q := sqlPolisRingkas(polis)
	if err := db.PeriksaSQL(q); err != nil {
		return kontrak.PolisRingkas{}, err
	}
	return r.bacaRingkas(ctx, q, nomorPolis)
}

// RingkasDariIdPega membaca header polis dari IDPEGA peserta (sqlPolisRingkasIdPega). Tidak ada = ErrPolisNomorTakDitemukan.
func (r *RingkasPolisLife) RingkasDariIdPega(ctx context.Context, idPega string) (kontrak.PolisRingkas, error) {
	idPega = strings.TrimSpace(idPega)
	if idPega == "" {
		return kontrak.PolisRingkas{}, fmt.Errorf("%w: IDPEGA kosong", kontrak.ErrPolisNomorTakDitemukan)
	}
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return kontrak.PolisRingkas{}, err
	}
	kerja, err := r.db.Qualify("T_WORK_POLIS")
	if err != nil {
		return kontrak.PolisRingkas{}, err
	}
	q := sqlPolisRingkasIdPega(polis, kerja)
	if err := db.PeriksaSQL(q); err != nil {
		return kontrak.PolisRingkas{}, err
	}
	return r.bacaRingkas(ctx, q, idPega)
}

// bacaRingkas menjalankan q (satu pengikat) dan memindai kolomPolisRingkas.
func (r *RingkasPolisLife) bacaRingkas(ctx context.Context, q, kunci string) (kontrak.PolisRingkas, error) {
	var (
		no, tipe, marketing, ceding, pemegang, bisnis sql.NullString
		diterima, status, statusUbah, prodID, prodNm  sql.NullString
		kodeBisnis, cedingKode, retroID, secID        sql.NullString
		tipeCeding, tipeCedingNm, prorata, wpc        sql.NullString
		retroNm, secNm, sobNm                         sql.NullString
		prodKe                                        int
	)
	if err := r.db.QueryRowContext(ctx, q, kunci).Scan(
		&no, &tipe, &marketing, &ceding, &pemegang, &bisnis,
		&diterima, &status, &statusUbah, &prodID, &prodNm, &prodKe,
		&kodeBisnis, &cedingKode, &retroID, &secID,
		&tipeCeding, &tipeCedingNm, &prorata, &wpc, &retroNm, &secNm, &sobNm); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return kontrak.PolisRingkas{}, fmt.Errorf("%w: %q", kontrak.ErrPolisNomorTakDitemukan, kunci)
		}
		return kontrak.PolisRingkas{}, fmt.Errorf("repository: membaca ringkasan polis: %w", err)
	}
	hasil := kontrak.PolisRingkas{
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
