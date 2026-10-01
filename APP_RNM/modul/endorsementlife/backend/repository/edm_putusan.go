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
