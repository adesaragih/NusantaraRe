package repository

// JSON_OFFER_LIFE - satu baris per kasus, ditulis setiap Save Offer (keputusan
// work owner 06-10-2026).
//
// Untuk apa berkas ini: meniru `InputOfferLife_ACT` langkah 3 ("Pega to
// json_offer_life") + 4 (`SaveOfferJsonLife_SQL` → `POOLDATA.INSERTJSONOFFERLIFE`)
// dengan SQL LANGSUNG di transaksi Save Offer milik pemanggil, SESUDAH isian
// tersimpan ke baris utama `T_PREMIUM_LIST`. Prosedurnya TIDAK dipanggil.
//
// `[terverifikasi]` nama kolom = ALL_TAB_COLUMNS POOLDATA (diberikan user
// 06-10-2026). Tabel WARISAN POOLDATA - nol DDL, nol sequence baru.
//
// Asal nilai (`TempInputDataLife.CARIn` langkah 3 → `p_*` SaveOfferJsonLife_SQL):
//
//	IDPEGA ← pzInsKey = pengenal kasus (SAMA dengan IDPEGA LIFEINPRODUCTION, polis_produksi.go)
//	BUSINESSCODE/NAME ← CARI1/2 BusinessCode/Name; CEDINGCO/NAME ← CARI3/4; POLICYHOLDER/NAME ← CARI6/7
//	REINSURANCETYPE ← CARI5 JenisAsuransi; TYPECEDING ← CARI8 @if(TypeCeding="1","QS",...)
//	RECEIVED_DATE..BINDING_DATE ← CARI9..16 (DateReceived, TanggalPenawaran, TanggalRespon,
//	    TanggalKonfirmasi, TanggalTBC, TanggalKonfirmasiBalik, TanggalRealisasi, TanggalBind) - tanpa jam
//	STATUS ← CARI17 Status; STATUS_UPDATE/FINAL ← CARI18/19; MARKETING_NOTE ← CARI20 KeteranganMarketing
//	OPERATORNAME ← CARI21 pxCreateOpName = CREATE_OP_NAME; UW_NOTE ← CARI23 KetentuanUnderwriting
//	TGL_INPUT ← SYSDATE
//
// ⛔ JSONDATA (CARI2 halaman JSON) dan OLDID (p_OFFERID ← OfferID) TIDAK diisi (NULL).
// ⛔ ID: ID baris lama kasus yang sama DIPERTAHANKAN; kasus baru memakai sequence
// warisan `JSON_OFFER_SEQ` (PK prosedur, `[data DBA]` 15-09-2026 - tiket 09) yang
// SUDAH ada. `GetIdOffer_SQL` membacanya kembali sebagai `'OFR-'||ID`.
// ⛔ Sekali per kasus: DELETE lalu INSERT - Save Offer berulang tidak menggandakan.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// ErrNilaiJSONOfferTerlaluPanjang - isian melebihi lebar kolom JSON_OFFER_LIFE
// (mis. STATUS / BUSINESSCODE VARCHAR2(10)); ditolak sebelum ORA-12899.
var ErrNilaiJSONOfferTerlaluPanjang = errors.New("repository: isian melebihi lebar kolom JSON_OFFER_LIFE")

// TabelJSONOfferLife dan urutan ID-nya (keduanya warisan POOLDATA).
const (
	TabelJSONOfferLife  = "JSON_OFFER_LIFE"
	UrutanJSONOfferLife = "JSON_OFFER_SEQ"
)

// kolomJSONOfferLife - kolom yang DIISI selain ID dan IDPEGA, urut p_* SaveOfferJsonLife_SQL.
// `Lebar` = byte VARCHAR2 tujuan (0 = DATE, tidak diperiksa); `Label` = nama layar.
var kolomJSONOfferLife = []struct {
	Kolom, Nilai, Label string
	Lebar               int
}{
	{"BUSINESSCODE", "p.BUSINESS_CODE", "Business Code", 10},
	{"BUSINESSNAME", "p.BUSINESS_NAME", "Business Name", 100},
	{"CEDINGCO", "p.CEDING_CO", "Ceding Code", 100},
	{"CEDINGCONAME", "p.CEDING_CO_NAME", "Ceding Name", 100},
	{"REINSURANCETYPE", "p.JENIS_ASURANSI", "Jenis Asuransi", 50},
	{"POLICYHOLDER", "p.POLICY_HOLDER", "Policy Holder", 100},
	{"POLICYHOLDERNAME", "p.POLICY_HOLDER_NAME", "Policy Holder Name", 100},
	// CARI8 - teks kosong Pega = NULL Oracle (pola TYPECEDINGNAME LIFEINPRODUCTION).
	{"TYPECEDING", "CASE p.TYPE_CEDING WHEN '1' THEN 'QS' WHEN '2' THEN 'SURPLUS' WHEN '3' THEN 'QS + SURPLUS' WHEN '4' THEN 'XOL' END", "Type Ceding", 50},
	{"STATUS", "p.STATUS_PENAWARAN", "Status", 10},
	{"STATUS_UPDATE", "p.STATUS_UPDATE", "Status Update", 1000},
	{"STATUS_FINAL", "p.STATUS_FINAL", "Status Final", 1000},
	{"MARKETING_NOTE", "p.KETERANGAN_MARKETING", "Keterangan Marketing", 1000},
	{"OPERATORNAME", "p.CREATE_OP_NAME", "Operator", 100},
	{"RECEIVED_DATE", "TRUNC(p.DATE_RECEIVED)", "", 0},
	{"OFFERING_DATE", "TRUNC(p.TANGGAL_PENAWARAN)", "", 0},
	{"RESPONSE_DATE", "TRUNC(p.TANGGAL_RESPON)", "", 0},
	{"CONFIRMATION_DATE", "TRUNC(p.TANGGAL_KONFIRMASI)", "", 0},
	{"MAX_TBC", "TRUNC(p.TANGGAL_TBC)", "", 0},
	{"RECONFIRMATION_DATE", "TRUNC(p.TANGGAL_KONFIRMASI_BALIK)", "", 0},
	{"REALIZATION_DATE", "TRUNC(p.TANGGAL_REALISASI)", "", 0},
	{"BINDING_DATE", "TRUNC(p.TANGGAL_BIND)", "", 0},
	{"UW_NOTE", "p.KETENTUAN_UNDERWRITING", "Ketentuan Underwriting", 1000},
	{"TGL_INPUT", "SYSDATE", "", 0},
}

// sqlIDJSONOfferLama - ID baris lama kasus ini (terbaru bila warisan menyimpan lebih dari satu).
func sqlIDJSONOfferLama(tujuan string) string {
	return fmt.Sprintf(`SELECT ID FROM %s WHERE IDPEGA = :1 AND ID IS NOT NULL
	  ORDER BY TGL_INPUT DESC NULLS LAST, ID DESC FETCH FIRST 1 ROWS ONLY`, tujuan)
}

// sqlLebarJSONOffer - LENGTHB tiap nilai teks dari baris utama, urut kolom berlebar.
func sqlLebarJSONOffer(polis string) string {
	var bagian []string
	for _, k := range kolomJSONOfferLife {
		if k.Lebar > 0 {
			bagian = append(bagian, "LENGTHB("+k.Nilai+")")
		}
	}
	return fmt.Sprintf(`SELECT %s FROM %s p WHERE p.ID = :1 AND p.ID = p.ID_PEGA`, strings.Join(bagian, ", "), polis)
}

func sqlHapusJSONOffer(tujuan string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE IDPEGA = :1`, tujuan)
}

// sqlSisipJSONOffer - satu baris dari baris utama kasus.
// ID lama: :1 ID, :2 IDPEGA, :3 kasus. Kasus baru (urutan ≠ ""): ID = sequence, :1 IDPEGA, :2 kasus.
func sqlSisipJSONOffer(tujuan, polis, urutan string) string {
	kolom := []string{"ID", "IDPEGA"}
	nilai := []string{":1", ":2"}
	kunci := ":3"
	if urutan != "" {
		nilai = []string{"TO_CHAR(" + urutan + ".NEXTVAL)", ":1"}
		kunci = ":2"
	}
	for _, k := range kolomJSONOfferLife {
		kolom = append(kolom, k.Kolom)
		nilai = append(nilai, k.Nilai)
	}
	return fmt.Sprintf(`INSERT INTO %s (%s) SELECT %s FROM %s p WHERE p.ID = %s AND p.ID = p.ID_PEGA`,
		tujuan, strings.Join(kolom, ", "), strings.Join(nilai, ", "), polis, kunci)
}

// periksaLebarJSONOffer - kalimat untuk setiap nilai yang tidak muat; kosong = muat semua.
func periksaLebarJSONOffer(panjang []sql.NullInt64) []string {
	var kurang []string
	i := 0
	for _, k := range kolomJSONOfferLife {
		if k.Lebar == 0 {
			continue
		}
		if n := panjang[i]; n.Valid && n.Int64 > int64(k.Lebar) {
			kurang = append(kurang, fmt.Sprintf("%s max %d characters (now %d)", k.Label, k.Lebar, n.Int64))
		}
		i++
	}
	return kurang
}

// JSONOfferLife menulis JSON_OFFER_LIFE.
type JSONOfferLife struct{ db *db.DB }

// NewJSONOfferLife menyusunnya.
func NewJSONOfferLife(db *db.DB) *JSONOfferLife { return &JSONOfferLife{db: db} }

// Ganti memeriksa lebar, membaca ID lama, menghapus baris kasus ini, lalu
// menyisipkan satu baris dari isian penawaran TERAKHIR - di transaksi tx.
// Mengembalikan ID yang ditulis ("" bila diberi sequence).
func (r *JSONOfferLife) Ganti(ctx context.Context, tx *db.Tx, kasusID string) (string, error) {
	tujuan, err := r.db.Qualify(TabelJSONOfferLife)
	if err != nil {
		return "", err
	}
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return "", err
	}
	urutan, err := r.db.Qualify(UrutanJSONOfferLife)
	if err != nil {
		return "", err
	}
	// 0. Lebar kolom - ditolak dengan kalimat jelas, bukan ORA-12899.
	ql := sqlLebarJSONOffer(polis)
	if err := db.PeriksaSQL(ql); err != nil {
		return "", err
	}
	var panjang []sql.NullInt64
	var tuju []any
	for _, k := range kolomJSONOfferLife {
		if k.Lebar > 0 {
			panjang = append(panjang, sql.NullInt64{})
		}
	}
	for i := range panjang {
		tuju = append(tuju, &panjang[i])
	}
	if err := tx.QueryRowContext(ctx, ql, kasusID).Scan(tuju...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrHeaderPolisTidakAda
		}
		return "", fmt.Errorf("repository: memeriksa lebar %s kasus %q: %w", TabelJSONOfferLife, kasusID, err)
	}
	if kurang := periksaLebarJSONOffer(panjang); len(kurang) > 0 {
		return "", fmt.Errorf("%w: %s", ErrNilaiJSONOfferTerlaluPanjang, strings.Join(kurang, "; "))
	}
	// 1. ID lama (bisa nol baris).
	qi := sqlIDJSONOfferLama(tujuan)
	if err := db.PeriksaSQL(qi); err != nil {
		return "", err
	}
	var idLama sql.NullString
	if err := tx.QueryRowContext(ctx, qi, kasusID).Scan(&idLama); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("repository: membaca ID %s kasus %q: %w", TabelJSONOfferLife, kasusID, err)
	}
	// 2-3. Hapus lalu sisip.
	qh := sqlHapusJSONOffer(tujuan)
	qs, arg := sqlSisipJSONOffer(tujuan, polis, urutan), []any{kasusID, kasusID}
	if idLama.Valid && idLama.String != "" {
		qs, arg = sqlSisipJSONOffer(tujuan, polis, ""), []any{idLama.String, kasusID, kasusID}
	}
	for _, q := range []string{qh, qs} {
		if err := db.PeriksaSQL(q); err != nil {
			return "", err
		}
	}
	if _, err := tx.ExecContext(ctx, qh, kasusID); err != nil {
		return "", fmt.Errorf("repository: menghapus %s kasus %q: %w", TabelJSONOfferLife, kasusID, err)
	}
	h, err := tx.ExecContext(ctx, qs, arg...)
	if err != nil {
		return "", fmt.Errorf("repository: menulis %s kasus %q: %w", TabelJSONOfferLife, kasusID, err)
	}
	if n, err := h.RowsAffected(); err == nil && n != 1 {
		return "", fmt.Errorf("repository: %s kasus %q menulis %d baris, mau 1", TabelJSONOfferLife, kasusID, n)
	}
	return idLama.String, nil
}
