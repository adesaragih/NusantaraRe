package repository

// Keputusan kasus (tiket 08/04/09/11) - riwayat `AddHistorySuggest`,
// penjaga anti-dobel `(NO_POLIS, PROD_KE)`, peresmian versi, penolakan, dan
// rekap warisan `M_LIFE_PREMIUM_SUMMARY` (isi prosedur
// `PEGA_M_LIFE_PREMIUM_SUMMARY` ditiru, prosedurnya tidak dipanggil - E2).
//
// ⛔ Nol `COMMIT`: seluruhnya di transaksi `services.Putuskan`.
//
// Dibaca sesudah: edm_simpan.go.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/endorsementlife/backend/models"
)

// batasRiwayat - baris riwayat keputusan yang dibaca layar.
const batasRiwayat = 200

// sqlRiwayat - grid `ConfirmSection` b1856, `AddHistorySuggest` 2 b493 `.No` `Descending` b491.
// Kunci `IDX_VS_PL` (PREMIUM_LIST_ID).
func sqlRiwayat(riwayat string) string {
	return fmt.Sprintf(`SELECT s.NO, TO_CHAR(s.DATE_SUGGEST, 'YYYY-MM-DD HH24:MI:SS'), s.PIC_SUGGEST, s.IS_CEDING_CONFIRM, s.COMMENT_SUGGEST
	  FROM %s s WHERE s.PREMIUM_LIST_ID = :1 ORDER BY s.NO DESC FETCH FIRST %d ROWS ONLY`, riwayat, batasRiwayat)
}

// Riwayat membaca riwayat keputusan kasus.
func (g *Gudang) Riwayat(ctx context.Context, tx *db.Tx, kasusID string) ([]models.BarisRiwayat, error) {
	n, err := g.nama(tabelRiwayat)
	if err != nil {
		return nil, err
	}
	q := sqlRiwayat(n[0])
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.pakai(tx).QueryContext(ctx, q, kasusID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca riwayat kasus %q: %w", kasusID, err)
	}
	defer func() { _ = rows.Close() }()
	hasil := []models.BarisRiwayat{}
	for rows.Next() {
		var no sql.NullInt64
		var tgl, pic, status, komentar sql.NullString
		if err := rows.Scan(&no, &tgl, &pic, &status, &komentar); err != nil {
			return nil, fmt.Errorf("repository: memindai riwayat: %w", err)
		}
		hasil = append(hasil, models.BarisRiwayat{No: int(no.Int64), Tanggal: tgl.String, PIC: pic.String, Status: status.String, Komentar: komentar.String})
	}
	return hasil, rows.Err()
}

// sqlSisipRiwayat - `AddHistorySuggest` 1: `No` = terbesar + 1 per kasus.
// Penampung :1 kasus (hash), :2 kasus, :3 waktu, :4 PIC, :5 status, :6 komentar, :7 kasus.
func sqlSisipRiwayat(riwayat string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, PREMIUM_LIST_ID, NO, DATE_SUGGEST, PIC_SUGGEST, IS_CEDING_CONFIRM, COMMENT_SUGGEST)
	  SELECT RAWTOHEX(STANDARD_HASH(:1 || '/H/' || TO_CHAR(NVL(MAX(s.NO), 0) + 1), 'MD5')), :2, NVL(MAX(s.NO), 0) + 1,
	         TO_DATE(:3, 'YYYY-MM-DD HH24:MI:SS'), :4, :5, :6
	    FROM %s s WHERE s.PREMIUM_LIST_ID = :7`, riwayat, riwayat)
}

// SisipRiwayat menulis satu baris riwayat keputusan.
func (g *Gudang) SisipRiwayat(ctx context.Context, tx *db.Tx, r models.RiwayatTulis) error {
	n, err := g.nama(tabelRiwayat)
	if err != nil {
		return err
	}
	q := sqlSisipRiwayat(n[0])
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, r.KasusID, r.KasusID, r.Waktu, db.KosongJadiNil(r.PIC), r.Status,
		db.KosongJadiNil(r.Komentar), r.KasusID); err != nil {
		return fmt.Errorf("repository: menulis riwayat kasus %q: %w", r.KasusID, err)
	}
	return nil
}

// sqlAdaVersiResmi - penjaga anti-dobel (tiket 09, AC 41/71): versi resmi
// ber-`(NO_POLIS, PROD_KE)` sama. Kasus terbuka/ditolak tidak ber-`NO_POLIS`.
// Kunci `IDX_PL_NOPOLIS_PRODKE` (migrasi 480).
func sqlAdaVersiResmi(polis string) string {
	return fmt.Sprintf(`SELECT 1 FROM %s p WHERE p.NO_POLIS = :1 AND p.PROD_KE = :2 FETCH FIRST 1 ROWS ONLY`, polis)
}

// AdaVersiResmi - true bila versi `(polis, prodKe)` sudah resmi.
func (g *Gudang) AdaVersiResmi(ctx context.Context, tx *db.Tx, nomorPolis string, prodKe int) (bool, error) {
	n, err := g.nama(tabelPolis)
	if err != nil {
		return false, err
	}
	q := sqlAdaVersiResmi(n[0])
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var satu int
	err = g.pakai(tx).QueryRowContext(ctx, q, nomorPolis, prodKe).Scan(&satu)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("repository: memeriksa versi %s/%d: %w", nomorPolis, prodKe, err)
	}
	return true, nil
}

// sqlResmikanKepala / sqlResmikanPeserta - `Confirm` (bab 7 PARITAS).
func sqlResmikanKepala(polis string) string {
	return fmt.Sprintf(`UPDATE %s p SET p.NO_POLIS = :1, p.PROD_KE = :2, p.NO_ENDORS = :3, p.PL_NUMBER_EDM = :4, p.STATUSS = :5
	  WHERE p.ID = :6 AND p.STATUSS IS NULL`, polis)
}

// `STATUS_OLD` ← `CARI47` (11.4 b4920 `.EDMStatus=="Old"` → b4843, 11.5 → b4982),
// `STATUS` ← `CARI48` (11.2 b4567, 11.3 b4706). Nilai kode dari `models`, diikat.
func sqlResmikanPeserta(peserta string) string {
	return fmt.Sprintf(`UPDATE %s d SET d.PL_NUMBER_EDM = :1,
	    d.STATUS_OLD = CASE WHEN d.EDM_STATUS = :2 THEN :3 ELSE :4 END, d.STATUS = :5
	  WHERE d.PREMIUM_LIST_ID = :6`, peserta)
}

// Resmikan menjadikan kasus versi resmi; mengembalikan cacah peserta.
func (g *Gudang) Resmikan(ctx context.Context, tx *db.Tx, r models.ResmiKasus) (int, error) {
	n, err := g.nama(tabelPolis, tabelPeserta)
	if err != nil {
		return 0, err
	}
	kepala, peserta := sqlResmikanKepala(n[0]), sqlResmikanPeserta(n[1])
	for _, q := range []string{kepala, peserta} {
		if err := db.PeriksaSQL(q); err != nil {
			return 0, err
		}
	}
	h, err := tx.ExecContext(ctx, kepala, r.NomorPolis, r.ProdKe, r.Nomor, r.Nomor, models.StatusKasusSelesai, r.ID)
	if err != nil {
		return 0, fmt.Errorf("repository: meresmikan kasus %q: %w", r.ID, err)
	}
	if err := db.PastikanSatuBaris(h, "peresmian kasus "+r.ID); err != nil {
		return 0, err
	}
	h, err = tx.ExecContext(ctx, peserta, r.Nomor, models.StatusOld, models.StatusLama(models.StatusOld),
		models.StatusLama(models.StatusNew), r.StatusJenis, r.ID)
	if err != nil {
		return 0, fmt.Errorf("repository: meresmikan peserta kasus %q: %w", r.ID, err)
	}
	c, err := h.RowsAffected()
	return int(c), err
}

// sqlTolak - `Decline` b1408 → `End1` `Resolved-Rejected` b642.
func sqlTolak(polis string) string {
	return fmt.Sprintf(`UPDATE %s p SET p.STATUSS = :1 WHERE p.ID = :2 AND p.STATUSS IS NULL`, polis)
}

// Tolak menutup kasus `Resolved-Rejected`.
func (g *Gudang) Tolak(ctx context.Context, tx *db.Tx, kasusID string) error {
	n, err := g.nama(tabelPolis)
	if err != nil {
		return err
	}
	q := sqlTolak(n[0])
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, q, models.StatusKasusDitolak, kasusID)
	if err != nil {
		return fmt.Errorf("repository: menolak kasus %q: %w", kasusID, err)
	}
	return db.PastikanSatuBaris(h, "penolakan kasus "+kasusID)
}

// --- rekap warisan (E2, tiket 11) -----------------------------------------------

// KolomRekapWarisan - daftar `INSERT` prosedur `PEGA_M_LIFE_PREMIUM_SUMMARY`
// (badan prosedur baris 48-85, `modul/premiumlistlife/docs/dba-procedure-PEGA_M_LIFE_PREMIUM_SUMMARY.md`),
// tanpa `ID` (sequence). `InsertPLSummary` b85 mengisi `P_COB` ←
// `pyWorkPage.BusinessName`, `P_PL_NUMBER`, `P_PL_NUMBER_EDM`, `P_CURRENCY`…`P_DEDUCTION` ←
// `CARI2`…`CARI34` = properti `CurrencyList` bernama sama (`InsertJsonPolisLife_Act`
// 12.1 b5411…), `P_IDPEGA` ← `pzInsKey` = pengenal kasus.
var KolomRekapWarisan = []string{
	"BALANCE", "BROKERAGE_FEE", "CLAIM", "OVR_COMM", "COB", "COMMISSION",
	"NET_PREMIUM_REFUND", "GROSS_PREMIUM_REFUND", "COMM_REFUND",
	"BROKERAGE_FEE_REFUND", "OVR_COMM_REFUND", "TAX_REFUND", "CURRENCY",
	"PL_NUMBER", "PL_NUMBER_EDM", "PREMIUM", "SHARE_RETRO",
	"GROSS_PREMIUM_RETRO", "DISCOUNT_PREMIUM_RETRO", "OVR_COMM_RETRO",
	"BROKERAGE_FEE_RETRO", "NET_PREMIUM_RETRO", "GROSS_PREMIUM_REFUND_RETRO",
	"DISCOUNT_PREMIUM_REFUND_RETRO", "OVR_COMM_REFUND_RETRO",
	"BROKERAGE_FEE_REFUND_RETRO", "NET_PREMIUM_REFUND_RETRO", "PROF_COMM",
	"TAX", "CLAIM_AMOUNT", "RI_ADMIN_FEE_RETRO", "RI_ADMIN_FEE_REFUND_RETRO",
	"RI_ADMIN_FEE_REFUND", "DEDUCTION_REFUND", "RI_ADMIN_FEE", "DEDUCTION",
	"IDPEGA",
}

// RekapWarisanTulis - kepala baris rekap warisan; nilai uang dari rekap kasus.
type RekapWarisanTulis struct {
	KasusID    string
	NomorPolis string
	Nomor      string
	COB        string
}

// sqlHapusRekapWarisan - hapus-sebelum-sisip `PL_NUMBER` + `IDPEGA` (pola
// PremiumList `sqlHapusSummaryWarisan`; tabelnya satu baris per mata uang per versi).
func sqlHapusRekapWarisan(warisan string) string {
	return fmt.Sprintf(`DELETE FROM %s w WHERE w.PL_NUMBER = :1 AND w.IDPEGA = :2`, warisan)
}

// sqlSisipRekapWarisan - satu baris per mata uang rekap kasus; uang kosong 0.
//
// ⛔ NUMBER langsung dari rekap kasus, BUKAN teks: parameter prosedur memang
// `VARCHAR2`, tetapi teks desimal yang masuk kolom `NUMBER` dikonversi
// menurut NLS sesi (`,` di sesi Indonesia) - `INSERT … SELECT` NUMBER ke
// NUMBER kebal NLS dan eksak. ⚠️ `[belum terverifikasi]` tipe kolom uang
// warisan: dianggap `NUMBER` seperti tiruan `uji/skemauji` dan katalog peserta.
// Penampung :1 COB, :2 PL_NUMBER, :3 PL_NUMBER_EDM, :4 IDPEGA, :5 kasus.
func sqlSisipRekapWarisan(warisan, urutan, rekap string) string {
	nilai := []string{"TO_CHAR(" + urutan + ".NEXTVAL)"}
	for _, k := range KolomRekapWarisan {
		switch k {
		case "COB":
			nilai = append(nilai, ":1")
		case "PL_NUMBER":
			nilai = append(nilai, ":2")
		case "PL_NUMBER_EDM":
			nilai = append(nilai, ":3")
		case "IDPEGA":
			nilai = append(nilai, ":4")
		case "CURRENCY":
			nilai = append(nilai, "r.CURRENCY")
		default:
			nilai = append(nilai, "NVL(r."+k+", 0)")
		}
	}
	return fmt.Sprintf(`INSERT INTO %s (ID, %s) SELECT %s FROM %s r WHERE r.PREMIUM_LIST_ID = :5`,
		warisan, strings.Join(KolomRekapWarisan, ", "), strings.Join(nilai, ", "), rekap)
}

// TulisRekapWarisan menyalin rekap kasus ke `M_LIFE_PREMIUM_SUMMARY`; mengembalikan cacah baris.
func (g *Gudang) TulisRekapWarisan(ctx context.Context, tx *db.Tx, r RekapWarisanTulis) (int, error) {
	n, err := g.nama(tabelRekapWarisan, urutanRekapWarisan, tabelRekap)
	if err != nil {
		return 0, err
	}
	hapus, sisip := sqlHapusRekapWarisan(n[0]), sqlSisipRekapWarisan(n[0], n[1], n[2])
	for _, q := range []string{hapus, sisip} {
		if err := db.PeriksaSQL(q); err != nil {
			return 0, err
		}
	}
	if _, err := tx.ExecContext(ctx, hapus, r.NomorPolis, r.KasusID); err != nil {
		return 0, fmt.Errorf("repository: menghapus rekap warisan %q: %w", r.KasusID, err)
	}
	h, err := tx.ExecContext(ctx, sisip, db.KosongJadiNil(r.COB), r.NomorPolis, r.Nomor, r.KasusID, r.KasusID)
	if err != nil {
		return 0, fmt.Errorf("repository: menulis rekap warisan %q: %w", r.KasusID, err)
	}
	c, err := h.RowsAffected()
	return int(c), err
}

// --- produksi warisan (K4 keputusan work owner 01-10-2026, OQ-EDM-010) ----------

// KolomProduksiWarisan - daftar `INSERT INTO POOLDATA.LIFEINPRODUCTION` `RDBList/SaveLifeinProduction_SQL.xml`
// b86/b87, VERBATIM urutan korpus (27 kolom; DEV 37 kolom - sepuluh sisanya ditulis jalur new business Pega,
// tidak oleh jalur EDM).
var KolomProduksiWarisan = []string{
	"IDPEGA", "NOPOLIS", "NOENDORS", "BUSINESSCODE", "BUSINESSNAME", "CEDINGCO", "CEDINGCONAME", "DATERECEIVED",
	"MARKETINGCODE", "MARKETINGNAME", "POLICYHOLDER", "POLICYHOLDERNAME", "PRORATETYPE", "CREATEOPNAME", "SOB",
	"SOBNAME", "TYPE", "TYPECEDING", "MOID", "NOOFFER", "RISLIPRNM", "RETROID", "RETRONAME", "TYPECEDINGNAME",
	"SECURITYREINSURERID", "SECURITYREINSURER", "TGL_INPUT",
}

// nilaiProduksiWarisan - VALUES `SaveLifeinProduction_SQL` ← `InsertJsonPolisLife_Act` 7 b1606 / 8 b1805
// (`InputDataLife.CARI1` b1702, `TempInputDataLife.CARI2`…`CARI27` b1878–b2382) = properti `pyWorkPage.*`,
// diambil dari kolom kepala kasus `T_PREMIUM_LIST` yang `MappingEDMLife` 9 isi dari properti yang sama
// (`models.KolomKepalaSalin`). Penampung: :1 `IDPEGA` ← `pzInsKey` = pengenal kasus (seperti rekap warisan),
// :2 `NOPOLIS` ← `PolicyNo`, :3 `NOENDORS` ← `PL_NUMBER_EDM` (CARI25 b2340).
var nilaiProduksiWarisan = map[string]string{
	"IDPEGA": ":1", "NOPOLIS": ":2", "NOENDORS": ":3",
	"BUSINESSCODE": "p.BUSINESS_CODE", "BUSINESSNAME": "p.BUSINESS_NAME", "CEDINGCO": "p.CEDING_CO",
	"CEDINGCONAME": "p.CEDING_CO_NAME",
	// `To_date({TempInputDataLife.CARI7}, 'DD/MM/YYYY')` - tanggal tanpa jam.
	"DATERECEIVED":  "TRUNC(p.DATE_RECEIVED)",
	"MARKETINGCODE": "p.MARKETING_CODE", "MARKETINGNAME": "p.MARKETING_NAME", "POLICYHOLDER": "p.POLICY_HOLDER",
	"POLICYHOLDERNAME": "p.POLICY_HOLDER_NAME", "PRORATETYPE": "p.PRO_RATE_TYPE",
	// CARI13 b2109 ← `pyWorkPage.pxCreateOpName`: pembuat kasus (`CREATE_OP_NAME`, akun pelaku - bukan nama orang).
	"CREATEOPNAME": "p.CREATE_OP_NAME",
	"SOB":          "p.SOB", "SOBNAME": "p.SOB_NAME", "TYPE": "p.TYPE", "TYPECEDING": "p.TYPE_CEDING", "MOID": "p.MO_ID",
	"NOOFFER": "p.NO_OFFER", "RISLIPRNM": "p.RI_SLIP_RNM", "RETROID": "p.RETRO_ID", "RETRONAME": "p.RETRO_NAME",
	// CARI24 b2319 `@if(TypeCeding="1","QS",@if(…="2","SURPLUS",@if(…="3","QS + SURPLUS",@if(…="4","XOL",""))))` -
	// teks kosong Pega = NULL Oracle.
	"TYPECEDINGNAME":      "CASE p.TYPE_CEDING " + kasusJenisCeding() + " END",
	"SECURITYREINSURERID": "p.SECURITY_REINSURER_ID", "SECURITYREINSURER": "p.SECURITY_REINSURER",
	"TGL_INPUT": "SYSDATE",
}

// kasusJenisCeding - cabang `WHEN` dari models.NamaJenisCeding, urut kode.
func kasusJenisCeding() string {
	var b strings.Builder
	for _, k := range []string{"1", "2", "3", "4"} {
		fmt.Fprintf(&b, "WHEN '%s' THEN '%s' ", k, models.NamaJenisCeding[k])
	}
	return strings.TrimSpace(b.String())
}

// ProduksiWarisanTulis - kepala baris `LIFEINPRODUCTION` kasus yang diresmikan.
type ProduksiWarisanTulis struct {
	KasusID    string
	NomorPolis string
	Nomor      string
}

// sqlSisipProduksiWarisan - satu baris dari kepala kasus. Penampung :1 IDPEGA, :2 NOPOLIS, :3 NOENDORS, :4 kasus
// (urut kemunculan - godror mengikat menurut urutan).
func sqlSisipProduksiWarisan(tujuan, polis string) string {
	nilai := make([]string, len(KolomProduksiWarisan))
	for i, k := range KolomProduksiWarisan {
		nilai[i] = nilaiProduksiWarisan[k]
	}
	return fmt.Sprintf(`INSERT INTO %s (%s) SELECT %s FROM %s p WHERE p.ID = :4`,
		tujuan, strings.Join(KolomProduksiWarisan, ", "), strings.Join(nilai, ", "), polis)
}

// TulisProduksiWarisan menulis baris `LIFEINPRODUCTION` kasus yang diresmikan; mengembalikan cacah baris (1).
//
// ⛔ Seperti Pega: SISIP murni, tanpa hapus-sebelum-sisip - tabelnya tanpa PK dan satu kasus diresmikan sekali
// (kasus tertutup tidak dapat diputuskan ulang, AC 29).
func (g *Gudang) TulisProduksiWarisan(ctx context.Context, tx *db.Tx, r ProduksiWarisanTulis) (int, error) {
	n, err := g.nama(tabelProduksiWarisan, tabelPolis)
	if err != nil {
		return 0, err
	}
	q := sqlSisipProduksiWarisan(n[0], n[1])
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	h, err := tx.ExecContext(ctx, q, r.KasusID, r.NomorPolis, r.Nomor, r.KasusID)
	if err != nil {
		return 0, fmt.Errorf("repository: menulis produksi warisan %q: %w", r.KasusID, err)
	}
	c, err := h.RowsAffected()
	if err == nil && c != 1 {
		err = fmt.Errorf("repository: produksi warisan kasus %q menulis %d baris, mau 1", r.KasusID, c)
	}
	return int(c), err
}

// --- peserta warisan (K5 keputusan work owner 01-10-2026, OQ-EDM-016) -----------

// KolomPesertaWarisanEDM - daftar `INSERT INTO POOLDATA.M_LIFE_PREMIUM_DETAIL` `RDBList/SaveMasterLPDet.xml`
// b86/b87, VERBATIM urutan korpus (80 kolom, `ID` pertama). Daftar yang sama dengan penulis new business
// PremiumList (`kolomPesertaWarisan` modul itu) - DITIRU, tidak diimpor; sumber nilainya berbeda (lihat
// nilaiPesertaWarisanEDM).
var KolomPesertaWarisanEDM = []string{
	"ID", "SHARE_NUSANTARA_RE", "SEX", "COMM", "FLEET_DISCOUNT", "POLICY_HOLDER", "PLAN", "PERIOD_YY", "PERIOD_MM",
	"POLICY_NO", "NET_PREMIUM", "NAME_OF_INSURED", "DESCRIPTION", "GROSS_PREMIUM", "EXPIRED_DATE", "DOB",
	"CERTIFICATE_NO", "BEGIN_DATE", "EFFECTIVE_DATE", "STNC", "LAPSE_DATE", "PASSED_PERIOD", "AGE", "CLAIM", "TAX",
	"BROKERAGE_FEE", "OVR_COMM", "CURRENCY", "MEDICAL_STATUS", "SUM_INSURED", "CEDING_RETENTION", "SUM_REASURED",
	"PROF_COMM", "GROSS_PREMIUM_REFUND", "NET_PREMIUM_REFUND", "COMM_REFUND", "BROKERAGE_FEE_REFUND",
	"OVR_COMM_REFUND", "TAX_REFUND", "SHARE_RETRO", "GROSS_PREMIUM_RETRO", "DISCOUNT_PREMIUM_RETRO", "OVR_COMM_RETRO",
	"BROKERAGE_FEE_RETRO", "NET_PREMIUM_RETRO", "GROSS_PREMIUM_REFUND_RETRO", "DISCOUNT_PREMIUM_REFUND_RETRO",
	"OVR_COMM_REFUND_RETRO", "BROKERAGE_FEE_REFUND_RETRO", "NET_PREMIUM_REFUND_RETRO", "CLAIM_AMOUNT", "PL_NUMBER",
	"PL_NUMBER_EDM", "CEDING_CO", "RATE", "PRORATETYPE", "SUM_AT_RISK_GROSS", "SUM_AT_RISK_RETRO", "RETROCEDED_SHARE",
	"SHARE_NUSANTARA_RE_GROSS", "GROSS_VALUATION_BEGIN_DATE", "GROSS_VALUATION_EXPIRED_DATE",
	"RETRO_VALUATION_BEGIN_DATE", "RETRO_VALUATION_EXPIRED_DATE", "WPC", "ENTRY_AGE", "CURRENT_AGE", "DEDUCTION",
	"FACTOR", "DEDUCTION_REFUND", "RI_ADMIN_FEE_REFUND_RETRO", "RI_ADMIN_FEE_RETRO", "RI_ADMIN_FEE_REFUND",
	"RI_ADMIN_FEE", "IDPEGA", "EDMSTATUS", "STATUSOLD", "STATUS", "EM_PERCENT", "RISK",
}

// nilaiPesertaWarisanEDM - VALUES `SaveMasterLPDet` ← `InsertJsonPolisLife_Act` 11 (ulang
// `PremiumListSummary.PremiumListDetail`, PRE=false b2861 → selalu) 11.1 b2889 (PRE=false b2909), per peserta
// kasus `d` (`T_PREMIUM_LIST_DETAIL`) dan kepala `p`; uji `TestPesertaWarisanEDMDariKorpus` menurunkannya ulang:
//
//	TempInputDetail.CARIn ← @toDecimal(.X)    NVL(d.X, 0)        `@toDecimal("")` = 0 (OQ-PL-10)
//	TempValue.X                               d.X                teks/angka apa adanya, kosong = NULL
//	To_date(TempValue.X, 'DD/MM/YYYY')        TRUNC(d.X)         kolom DATE 052; STNC/WPC teks → TO_DATE
//	CARI32 / CARI34 ← pyWorkPage.CedingCo / .ProRateType          p.CEDING_CO / p.PRO_RATE_TYPE
//	CARI47 (11.4/11.5) / CARI48 (11.2/11.3)                       d.STATUS_OLD / d.STATUS - Resmikan lebih dulu
//	:1 PL_NUMBER ← PremiumListSummary.PL_NUMBER, :2 PL_NUMBER_EDM, :3 IDPEGA ← pzInsKey = pengenal kasus
//
// ⚠️ PERSIS Pega, termasuk dua kolom yang hilang: `EM_PERCENT` ← CARI49 dan `RISK` ← CARI50 TIDAK PERNAH
// ditetapkan 11.1 jalur endorsement (CARI12 ditetapkan dua kali: `.EM_PERCENT` lalu ditimpa
// `.GROSS_PREMIUM_REFUND`) - keduanya NULL, berbeda dari jalur new business PremiumList yang mengisinya.
var nilaiPesertaWarisanEDM = petaNilaiPesertaWarisanEDM(map[string]string{
	"ID": "", "PL_NUMBER": ":1", "PL_NUMBER_EDM": ":2", "IDPEGA": ":3",
	"SEX": "d.SEX", "POLICY_HOLDER": "d.POLICY_HOLDER", "PLAN": "d.PLAN", "PERIOD_YY": "d.PERIOD_YY",
	"PERIOD_MM": "d.PERIOD_MM", "POLICY_NO": "d.POLICY_NO", "NAME_OF_INSURED": "d.NAME_OF_INSURED",
	"DESCRIPTION": "d.DESCRIPTION", "CERTIFICATE_NO": "d.CERTIFICATE_NO", "PASSED_PERIOD": "d.PASSED_PERIOD",
	"AGE": "d.AGE", "CURRENCY": "d.CURRENCY", "MEDICAL_STATUS": "d.MEDICAL_STATUS", "ENTRY_AGE": "d.ENTRY_AGE",
	"CURRENT_AGE":  "d.CURRENT_AGE",
	"EXPIRED_DATE": "TRUNC(d.EXPIRED_DATE)", "DOB": "TRUNC(d.DOB)", "BEGIN_DATE": "TRUNC(d.BEGIN_DATE)",
	"EFFECTIVE_DATE": "TRUNC(d.EFFECTIVE_DATE)", "LAPSE_DATE": "TRUNC(d.LAPSE_DATE)",
	"GROSS_VALUATION_BEGIN_DATE":   "TRUNC(d.GROSS_VALUATION_BEGIN_DATE)",
	"GROSS_VALUATION_EXPIRED_DATE": "TRUNC(d.GROSS_VALUATION_EXPIRED_DATE)",
	// `TempValue.RETROCESSION_VALUATION_*` → kolom `RETRO_VALUATION_*` (052).
	"RETRO_VALUATION_BEGIN_DATE":   "TRUNC(d.RETRO_VALUATION_BEGIN_DATE)",
	"RETRO_VALUATION_EXPIRED_DATE": "TRUNC(d.RETRO_VALUATION_EXPIRED_DATE)",
	"STNC":                         "TO_DATE(d.STNC, 'DD/MM/YYYY')", "WPC": "TO_DATE(d.WPC, 'DD/MM/YYYY')",
	"CEDING_CO": "p.CEDING_CO", "PRORATETYPE": "p.PRO_RATE_TYPE",
	"EDMSTATUS": "d.EDM_STATUS", "STATUSOLD": "d.STATUS_OLD", "STATUS": "d.STATUS",
	"EM_PERCENT": "NULL", "RISK": "NULL",
})

// desimalPesertaWarisanEDM - kolom `CARIn ← @toDecimal(.X)` 11.1, nama kolom = nama properti.
var desimalPesertaWarisanEDM = []string{
	"SHARE_NUSANTARA_RE", "COMM", "FLEET_DISCOUNT", "NET_PREMIUM", "GROSS_PREMIUM", "CLAIM", "TAX", "BROKERAGE_FEE",
	"OVR_COMM", "SUM_INSURED", "CEDING_RETENTION", "SUM_REASURED", "PROF_COMM", "GROSS_PREMIUM_REFUND",
	"NET_PREMIUM_REFUND", "COMM_REFUND", "BROKERAGE_FEE_REFUND", "OVR_COMM_REFUND", "TAX_REFUND", "SHARE_RETRO",
	"GROSS_PREMIUM_RETRO", "DISCOUNT_PREMIUM_RETRO", "OVR_COMM_RETRO", "BROKERAGE_FEE_RETRO", "NET_PREMIUM_RETRO",
	"GROSS_PREMIUM_REFUND_RETRO", "DISCOUNT_PREMIUM_REFUND_RETRO", "OVR_COMM_REFUND_RETRO",
	"BROKERAGE_FEE_REFUND_RETRO", "NET_PREMIUM_REFUND_RETRO", "CLAIM_AMOUNT", "RATE", "SUM_AT_RISK_GROSS",
	"SUM_AT_RISK_RETRO", "RETROCEDED_SHARE", "SHARE_NUSANTARA_RE_GROSS", "DEDUCTION", "FACTOR", "DEDUCTION_REFUND",
	"RI_ADMIN_FEE_REFUND_RETRO", "RI_ADMIN_FEE_RETRO", "RI_ADMIN_FEE_REFUND", "RI_ADMIN_FEE",
}

// petaNilaiPesertaWarisanEDM melengkapi peta nilai dengan kolom desimalPesertaWarisanEDM.
func petaNilaiPesertaWarisanEDM(m map[string]string) map[string]string {
	for _, k := range desimalPesertaWarisanEDM {
		m[k] = "NVL(d." + k + ", 0)"
	}
	return m
}

// PesertaWarisanTulis - kepala baris `M_LIFE_PREMIUM_DETAIL` kasus yang diresmikan.
type PesertaWarisanTulis struct {
	KasusID    string
	NomorPolis string
	Nomor      string
}

// sqlSisipPesertaWarisanEDM - setiap peserta kasus (Old, New, Delete, Batal) satu baris warisan.
//
// ⛔ E4: tabel warisan ±66,8 juta baris hanya SASARAN `INSERT` - nol baca, nol pemindaian. Sumbernya peserta
// kasus `T_PREMIUM_LIST_DETAIL` berkunci `IDX_PLD_PL` (`PREMIUM_LIST_ID`, migrasi 052) + kepalanya menurut PK.
// Penampung :1 PL_NUMBER, :2 PL_NUMBER_EDM, :3 IDPEGA, :4 kasus - urut kemunculan.
func sqlSisipPesertaWarisanEDM(tujuan, urutan, peserta, polis string) string {
	nilai := make([]string, len(KolomPesertaWarisanEDM))
	for i, k := range KolomPesertaWarisanEDM {
		if k == "ID" {
			nilai[i] = "TO_CHAR(" + urutan + ".NEXTVAL)"
			continue
		}
		nilai[i] = nilaiPesertaWarisanEDM[k]
	}
	return fmt.Sprintf(`INSERT INTO %s (%s) SELECT %s FROM %s d JOIN %s p ON p.ID = d.PREMIUM_LIST_ID WHERE d.PREMIUM_LIST_ID = :4`,
		tujuan, strings.Join(KolomPesertaWarisanEDM, ", "), strings.Join(nilai, ", "), peserta, polis)
}

// TulisPesertaWarisan menulis peserta kasus yang diresmikan ke `M_LIFE_PREMIUM_DETAIL`; mengembalikan cacah baris.
//
// ⛔ Seperti Pega: SISIP murni (`SaveMasterLPDet` tanpa hapus) - satu kasus diresmikan sekali (AC 29); dipanggil
// SESUDAH Resmikan, yang menetapkan `STATUS`/`STATUS_OLD`/`PL_NUMBER_EDM` peserta (11.2–11.5).
func (g *Gudang) TulisPesertaWarisan(ctx context.Context, tx *db.Tx, r PesertaWarisanTulis) (int, error) {
	n, err := g.nama(tabelPesertaWarisanEDM, urutanPesertaWarisan, tabelPeserta, tabelPolis)
	if err != nil {
		return 0, err
	}
	q := sqlSisipPesertaWarisanEDM(n[0], n[1], n[2], n[3])
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	h, err := tx.ExecContext(ctx, q, r.NomorPolis, r.Nomor, r.KasusID, r.KasusID)
	if err != nil {
		return 0, fmt.Errorf("repository: menulis peserta warisan %q: %w", r.KasusID, err)
	}
	c, err := h.RowsAffected()
	return int(c), err
}
