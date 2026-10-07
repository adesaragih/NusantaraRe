package repository

// LIFEINPRODUCTION - satu baris per polis yang di-Confirm (Resolved-Completed),
// keputusan work owner 06-10-2026.
//
// Untuk apa berkas ini: meniru `InsertJsonPolisLife_Act` langkah 9 ("Pega to
// json_offer & lifeinproduction") + 11 (`SaveLifeinProduction_SQL`) dengan SQL
// LANGSUNG, di transaksi Confirm milik pemanggil, SESUDAH PL Number dan WPC final
// tertulis ke T_PREMIUM_LIST.
//
// `[terverifikasi]` kolom = daftar INSERT `PremiumList Life/RDBList/
// SaveLifeinProduction_SQL.xml` (36 kolom; 27 di antaranya juga ditulis modul
// endorsementlife, `KolomProduksiWarisan`). Tabel WARISAN POOLDATA - nol DDL.
//
// Asal nilai (langkah 9 → VALUES), seluruhnya dari kepala kasus T_PREMIUM_LIST:
//
//	IDPEGA ← InputDataLife.CARI1 (pzInsKey) = pengenal kasus (pola rekap warisan dan endorsementlife)
//	NOPOLIS ← PremiumListSummary.PL_NUMBER; DATERECEIVED / WPC ← To_date(..,'DD/MM/YYYY') = tanggal tanpa jam
//	TYPECEDINGNAME ← CARI24 @if(TypeCeding="1","QS",...); CREATEOPNAME ← pxCreateOpName = CREATE_OP_NAME
//	RISLIPRNM ← PremiumListSummary.RISLIPRNM → RI_SLIP_RNM data polis
//	LAYER_1..LAYER_4 TIDAK diisi (NULL) - keputusan work owner 06-10-2026 "tidak perlu digunakan lagi"
//	TGL_INPUT ← sysdate
//
// ⛔ Sekali per kasus: baris lama kasus yang sama DIHAPUS lebih dulu (Confirm yang
// diulang tidak menggandakan). Nol COMMIT (ADR-U-0029).

import (
	"context"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// TabelProduksiLife - tabel warisan yang ditulis.
const TabelProduksiLife = "LIFEINPRODUCTION"

// KolomProduksiLife - kolom yang DIISI (urut SaveLifeinProduction_SQL, tanpa LAYER_1..4).
var KolomProduksiLife = []string{
	"IDPEGA", "NOPOLIS", "BUSINESSCODE", "BUSINESSNAME", "CEDINGCO", "CEDINGCONAME", "DATERECEIVED",
	"MARKETINGCODE", "MARKETINGNAME", "POLICYHOLDER", "POLICYHOLDERNAME", "PRORATETYPE", "CREATEOPNAME", "SOB",
	"SOBNAME", "TYPE", "TYPECEDING", "MOID", "NOOFFER", "RISLIPRNM", "RETROID", "RETRONAME", "DESCRIPTION",
	"TYPECEDINGNAME", "WPC", "SECURITYREINSURERID", "SECURITYREINSURER", "PRODUCTNAMEID", "PRODUCTNAME",
	"ANNUITYINTEREST", "PREMIUMREFUNDFACTOR", "TGL_INPUT",
}

// nilaiProduksiLife - ekspresi SELECT tiap kolom dari kepala kasus `p`. :1 IDPEGA, :2 NOPOLIS.
var nilaiProduksiLife = map[string]string{
	"IDPEGA": ":1", "NOPOLIS": ":2",
	"BUSINESSCODE": "p.BUSINESS_CODE", "BUSINESSNAME": "p.BUSINESS_NAME", "CEDINGCO": "p.CEDING_CO",
	"CEDINGCONAME": "p.CEDING_CO_NAME", "DATERECEIVED": "TRUNC(p.DATE_RECEIVED)",
	"MARKETINGCODE": "p.MARKETING_CODE", "MARKETINGNAME": "p.MARKETING_NAME", "POLICYHOLDER": "p.POLICY_HOLDER",
	"POLICYHOLDERNAME": "p.POLICY_HOLDER_NAME", "PRORATETYPE": "p.PRO_RATE_TYPE", "CREATEOPNAME": "p.CREATE_OP_NAME",
	"SOB": "p.SOB", "SOBNAME": "p.SOB_NAME", "TYPE": "p.TYPE", "TYPECEDING": "p.TYPE_CEDING", "MOID": "p.MO_ID",
	"NOOFFER": "p.NO_OFFER", "RISLIPRNM": "p.RI_SLIP_RNM", "RETROID": "p.RETRO_ID", "RETRONAME": "p.RETRO_NAME",
	"DESCRIPTION": "p.DESCRIPTION",
	// CARI24 - teks kosong Pega = NULL Oracle.
	"TYPECEDINGNAME": "CASE p.TYPE_CEDING WHEN '1' THEN 'QS' WHEN '2' THEN 'SURPLUS' WHEN '3' THEN 'QS + SURPLUS' WHEN '4' THEN 'XOL' END",
	"WPC":            "TRUNC(p.WPC)", "SECURITYREINSURERID": "p.SECURITY_REINSURER_ID", "SECURITYREINSURER": "p.SECURITY_REINSURER",
	"PRODUCTNAMEID": "p.PRODUCT_NAME_ID", "PRODUCTNAME": "p.PRODUCT_NAME",
	"ANNUITYINTEREST": "p.ANNUITY_INTEREST", "PREMIUMREFUNDFACTOR": "p.PREMIUM_REFUND_FACTOR",
	"TGL_INPUT": "SYSDATE",
}

func sqlHapusProduksiLife(tujuan string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE IDPEGA = :1`, tujuan)
}

// sqlSisipProduksiLife - satu baris dari kepala kasus (baris utama ID = ID_PEGA).
// Penampung urut kemunculan: :1 IDPEGA, :2 NOPOLIS, :3 kasus.
func sqlSisipProduksiLife(tujuan, polis string) string {
	nilai := make([]string, len(KolomProduksiLife))
	for i, k := range KolomProduksiLife {
		nilai[i] = nilaiProduksiLife[k]
	}
	return fmt.Sprintf(`INSERT INTO %s (%s) SELECT %s FROM %s p WHERE p.ID = :3 AND p.ID = p.ID_PEGA`,
		tujuan, strings.Join(KolomProduksiLife, ", "), strings.Join(nilai, ", "), polis)
}

// ProduksiLife menulis LIFEINPRODUCTION.
type ProduksiLife struct{ db *db.DB }

// NewProduksiLife menyusunnya.
func NewProduksiLife(db *db.DB) *ProduksiLife { return &ProduksiLife{db: db} }

// Ganti menghapus baris LIFEINPRODUCTION kasus ini lalu menyisipkan satu baris baru
// dari kepala kasus, di transaksi tx. Galat bila tidak tepat satu baris tersisip.
func (r *ProduksiLife) Ganti(ctx context.Context, tx *db.Tx, kasusID, plNumber string) error {
	tujuan, err := r.db.Qualify(TabelProduksiLife)
	if err != nil {
		return err
	}
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return err
	}
	qh, qs := sqlHapusProduksiLife(tujuan), sqlSisipProduksiLife(tujuan, polis)
	for _, q := range []string{qh, qs} {
		if err := db.PeriksaSQL(q); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, qh, kasusID); err != nil {
		return fmt.Errorf("repository: menghapus %s kasus %q: %w", TabelProduksiLife, kasusID, err)
	}
	h, err := tx.ExecContext(ctx, qs, kasusID, plNumber, kasusID)
	if err != nil {
		return fmt.Errorf("repository: menulis %s kasus %q: %w", TabelProduksiLife, kasusID, err)
	}
	if n, err := h.RowsAffected(); err == nil && n != 1 {
		return fmt.Errorf("repository: %s kasus %q menulis %d baris, mau 1", TabelProduksiLife, kasusID, n)
	}
	return nil
}
