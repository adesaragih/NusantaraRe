package repository

// Salinan peserta ke tabel warisan `M_LIFE_PREMIUM_DETAIL` - pl2, tiket 05a.
//
// Untuk apa berkas ini: Claim Life membaca peserta polis dari tabel warisan
// (`GET /api/peserta-life`, pesertapolis.go). Selama tabel itu yang dibaca,
// polis yang disubmit di sistem baru harus MUNCUL di sana - kalau tidak,
// klaim atas polis baru tidak menemukan pesertanya sama sekali.
//
// ⛔ BERKAS INI HANYA MENULIS. Tidak ada satu pun kueri baca di sini, dan itu
// disengaja: penjaga `TestQueryTabelPesertaSelaluBerindexDanBerbatas`
// memeriksa setiap kueri baca di berkas yang menyebut `namaTabelPeserta`, dan
// syaratnya (penyaring ber-index + batas hasil) benar untuk MEMBACA tabel
// 66,8 juta baris. Bahan salinannya dibaca dari tabel KAMI oleh
// `SummaryPolis.SumberWarisan` (polis_summary.go).
//
// ⛔ SATU TRANSAKSI dengan penomoran dan rekap (pl2). `SaveMasterLPDet` di
// Pega menulis `COMMIT;` di baris 252 - itu titik potong yang membuat nomor
// terbit tanpa peserta warisan, atau sebaliknya. Di sini nol `COMMIT`.
//
// Dibaca sesudah: polis_summary.go, pesertapolis.go (pembacanya).

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nusantarare/internal/models"
)

// jenisNilai membedakan cara sebuah nilai menyeberang ke tabel warisan.
type jenisNilai int

const (
	nilaiTeks jenisNilai = iota
	// nilaiAngka - dibaca lewat `TO_CHAR` ber-NLS, dikirim sebagai teks.
	nilaiAngka
	// nilaiTanggal - kolom DATE di tabel kami; dibaca `TO_CHAR(…,'DD/MM/YYYY')`.
	nilaiTanggal
	// nilaiTeksTanggal - kolom TEKS `dd/mm/yyyy` di tabel kami (STNC, WPC).
	nilaiTeksTanggal
)

// kolomPesertaWarisan adalah pemetaan `SaveMasterLPDet`, VERBATIM - pl2.
//
// `[terverifikasi]` Kolom: daftar `INSERT INTO POOLDATA.M_LIFE_PREMIUM_DETAIL`
// di `RDBList/SaveMasterLPDet.xml` (baris 87-167), urut dokumen - 80 kolom
// termasuk `ID`, yang diisi sequence dan tidak ada di daftar ini. Sumber:
// `VALUES` rule itu dan penetapan `TempInputDetail.CARIn` di
// `Activity/InsertLifePremiumDetail_act.xml`:
//
//	grep -A1 "<PropertiesName>TempInputDetail.CARI[0-9]*<" \
//	  Activity/InsertLifePremiumDetail_act.xml
//
// `d.` = baris peserta (`TempValue.X`, atau `CARIn ← @toDecimal(.X)`), `p.` =
// header polis (`TempWorkPage.X`, `pyWorkPage.pzInsKey`). Sumber kosong =
// nilai yang dirakit `nilaiSalinWarisan` (`PL_NUMBER`, `IDPEGA`, `STATUSOLD`,
// `STATUS`).
//
// ⚠️ Lima yang TIDAK bernama sama, dan sebabnya:
//
//	CEDING_CO    ← CARI32 = TempWorkPage.CedingCo    -> p.CEDING_CO
//	PRORATETYPE  ← CARI34 = TempWorkPage.ProRateType -> p.PRO_RATE_TYPE
//	IDPEGA       ← pyWorkPage.pzInsKey               -> pengenal work (polisID)
//	EDMSTATUS    ← TempValue.EDMStatus               -> d.EDM_STATUS
//	RETRO_VALUATION_*_DATE ← TempValue.RETROCESSION_VALUATION_*_DATE
//	                                   -> d.RETRO_VALUATION_*_DATE (052)
//
// ⛔ DUA NILAI TERTANAM, disalin apa adanya (lihat `nilaiSalinWarisan`):
//
//	STATUSOLD ← CARI47 = 0
//	STATUS    ← CARI48 = @if(Type=="QR","0",@if(Type=="QP","0","1"))
//
// ⚠️ `CARI46 = @toDecimal(.KTP)` ditetapkan activity tetapi TIDAK dipakai
// `SaveMasterLPDet`. Tidak disalin.
//
// ⚠️ `RISK` ditulis `@toDecimal(.RISK)` di rule aslinya, padahal kolomnya
// teks di migrasi 052 dan tipenya di tabel warisan tidak ada di korpus. Ia
// dikirim sebagai teks apa adanya - `[belum terverifikasi]` tipe kolom
// warisannya (OQ-001 keluarga DDL).
var kolomPesertaWarisan = []struct {
	Kolom  string
	Sumber string
	Jenis  jenisNilai
}{
	{"SHARE_NUSANTARA_RE", "d.SHARE_NUSANTARA_RE", nilaiAngka},
	{"SEX", "d.SEX", nilaiTeks},
	{"COMM", "d.COMM", nilaiAngka},
	{"FLEET_DISCOUNT", "d.FLEET_DISCOUNT", nilaiAngka},
	{"POLICY_HOLDER", "d.POLICY_HOLDER", nilaiTeks},
	{"PLAN", "d.PLAN", nilaiTeks},
	{"PERIOD_YY", "d.PERIOD_YY", nilaiAngka},
	{"PERIOD_MM", "d.PERIOD_MM", nilaiAngka},
	{"POLICY_NO", "d.POLICY_NO", nilaiTeks},
	{"NET_PREMIUM", "d.NET_PREMIUM", nilaiAngka},
	{"NAME_OF_INSURED", "d.NAME_OF_INSURED", nilaiTeks},
	{"DESCRIPTION", "d.DESCRIPTION", nilaiTeks},
	{"GROSS_PREMIUM", "d.GROSS_PREMIUM", nilaiAngka},
	{"EXPIRED_DATE", "d.EXPIRED_DATE", nilaiTanggal},
	{"DOB", "d.DOB", nilaiTanggal},
	{"CERTIFICATE_NO", "d.CERTIFICATE_NO", nilaiTeks},
	{"BEGIN_DATE", "d.BEGIN_DATE", nilaiTanggal},
	{"EFFECTIVE_DATE", "d.EFFECTIVE_DATE", nilaiTanggal},
	{"STNC", "d.STNC", nilaiTeksTanggal},
	{"LAPSE_DATE", "d.LAPSE_DATE", nilaiTanggal},
	{"PASSED_PERIOD", "d.PASSED_PERIOD", nilaiAngka},
	{"AGE", "d.AGE", nilaiAngka},
	{"CLAIM", "d.CLAIM", nilaiAngka},
	{"TAX", "d.TAX", nilaiAngka},
	{"BROKERAGE_FEE", "d.BROKERAGE_FEE", nilaiAngka},
	{"OVR_COMM", "d.OVR_COMM", nilaiAngka},
	{"CURRENCY", "d.CURRENCY", nilaiTeks},
	{"MEDICAL_STATUS", "d.MEDICAL_STATUS", nilaiTeks},
	{"SUM_INSURED", "d.SUM_INSURED", nilaiAngka},
	{"CEDING_RETENTION", "d.CEDING_RETENTION", nilaiAngka},
	{"SUM_REASURED", "d.SUM_REASURED", nilaiAngka},
	{"PROF_COMM", "d.PROF_COMM", nilaiAngka},
	{"GROSS_PREMIUM_REFUND", "d.GROSS_PREMIUM_REFUND", nilaiAngka},
	{"NET_PREMIUM_REFUND", "d.NET_PREMIUM_REFUND", nilaiAngka},
	{"COMM_REFUND", "d.COMM_REFUND", nilaiAngka},
	{"BROKERAGE_FEE_REFUND", "d.BROKERAGE_FEE_REFUND", nilaiAngka},
	{"OVR_COMM_REFUND", "d.OVR_COMM_REFUND", nilaiAngka},
	{"TAX_REFUND", "d.TAX_REFUND", nilaiAngka},
	{"SHARE_RETRO", "d.SHARE_RETRO", nilaiAngka},
	{"GROSS_PREMIUM_RETRO", "d.GROSS_PREMIUM_RETRO", nilaiAngka},
	{"DISCOUNT_PREMIUM_RETRO", "d.DISCOUNT_PREMIUM_RETRO", nilaiAngka},
	{"OVR_COMM_RETRO", "d.OVR_COMM_RETRO", nilaiAngka},
	{"BROKERAGE_FEE_RETRO", "d.BROKERAGE_FEE_RETRO", nilaiAngka},
	{"NET_PREMIUM_RETRO", "d.NET_PREMIUM_RETRO", nilaiAngka},
	{"GROSS_PREMIUM_REFUND_RETRO", "d.GROSS_PREMIUM_REFUND_RETRO", nilaiAngka},
	{"DISCOUNT_PREMIUM_REFUND_RETRO", "d.DISCOUNT_PREMIUM_REFUND_RETRO", nilaiAngka},
	{"OVR_COMM_REFUND_RETRO", "d.OVR_COMM_REFUND_RETRO", nilaiAngka},
	{"BROKERAGE_FEE_REFUND_RETRO", "d.BROKERAGE_FEE_REFUND_RETRO", nilaiAngka},
	{"NET_PREMIUM_REFUND_RETRO", "d.NET_PREMIUM_REFUND_RETRO", nilaiAngka},
	{"CLAIM_AMOUNT", "d.CLAIM_AMOUNT", nilaiAngka},
	{"PL_NUMBER", "", nilaiTeks},
	{"PL_NUMBER_EDM", "p.PL_NUMBER_EDM", nilaiTeks},
	{"CEDING_CO", "p.CEDING_CO", nilaiTeks},
	{"RATE", "d.RATE", nilaiAngka},
	{"PRORATETYPE", "p.PRO_RATE_TYPE", nilaiTeks},
	{"SUM_AT_RISK_GROSS", "d.SUM_AT_RISK_GROSS", nilaiAngka},
	{"SUM_AT_RISK_RETRO", "d.SUM_AT_RISK_RETRO", nilaiAngka},
	{"RETROCEDED_SHARE", "d.RETROCEDED_SHARE", nilaiAngka},
	{"SHARE_NUSANTARA_RE_GROSS", "d.SHARE_NUSANTARA_RE_GROSS", nilaiAngka},
	{"GROSS_VALUATION_BEGIN_DATE", "d.GROSS_VALUATION_BEGIN_DATE", nilaiTanggal},
	{"GROSS_VALUATION_EXPIRED_DATE", "d.GROSS_VALUATION_EXPIRED_DATE", nilaiTanggal},
	{"RETRO_VALUATION_BEGIN_DATE", "d.RETRO_VALUATION_BEGIN_DATE", nilaiTanggal},
	{"RETRO_VALUATION_EXPIRED_DATE", "d.RETRO_VALUATION_EXPIRED_DATE", nilaiTanggal},
	{"WPC", "d.WPC", nilaiTeksTanggal},
	{"ENTRY_AGE", "d.ENTRY_AGE", nilaiAngka},
	{"CURRENT_AGE", "d.CURRENT_AGE", nilaiAngka},
	{"DEDUCTION", "d.DEDUCTION", nilaiAngka},
	{"FACTOR", "d.FACTOR", nilaiAngka},
	{"DEDUCTION_REFUND", "d.DEDUCTION_REFUND", nilaiAngka},
	{"RI_ADMIN_FEE_REFUND_RETRO", "d.RI_ADMIN_FEE_REFUND_RETRO", nilaiAngka},
	{"RI_ADMIN_FEE_RETRO", "d.RI_ADMIN_FEE_RETRO", nilaiAngka},
	{"RI_ADMIN_FEE_REFUND", "d.RI_ADMIN_FEE_REFUND", nilaiAngka},
	{"RI_ADMIN_FEE", "d.RI_ADMIN_FEE", nilaiAngka},
	{"IDPEGA", "", nilaiTeks},
	{"EDMSTATUS", "d.EDM_STATUS", nilaiTeks},
	{"STATUSOLD", "", nilaiTeks},
	{"STATUS", "", nilaiTeks},
	{"EM_PERCENT", "d.EM_PERCENT", nilaiAngka},
	{"RISK", "d.RISK", nilaiTeks},
}

// kolomNolBilaKosongWarisan - OQ-PL-10 DITUTUP 29-09-2026 (GILIRAN-17)
// `[keputusan work owner]`: kolom yang di `SaveMasterLPDet` memakai
// `TempInputDetail.CARIn` dengan `CARIn = @toDecimal(.X)` di
// `InsertLifePremiumDetail_act` langkah 3.3.3 (b1843, hidup) - 45 kolom
// termasuk `RISK` (CARI50). `@toDecimal("")` = 0, jadi di tabel warisan kolom
// ini TIDAK PERNAH NULL. Kolom teks, tanggal, dan `PERIOD_YY/MM`,
// `PASSED_PERIOD`, `AGE`, `ENTRY_AGE`, `CURRENT_AGE` (disalin `TempValue.X`
// langsung, tanpa `@toDecimal`) tetap NULL bila kosong.
//
// ⛔ SATU tempat, di tepi repository WARISAN. ADR-U-0027 (kosong bukan nol)
// tetap berlaku untuk tabel `T_*`. Himpunannya diturunkan ulang dari korpus
// oleh `TestKolomNolBilaKosongWarisanDariKorpus`, dua arah.
var kolomNolBilaKosongWarisan = map[string]bool{
	"SHARE_NUSANTARA_RE": true, "COMM": true, "FLEET_DISCOUNT": true, "NET_PREMIUM": true,
	"GROSS_PREMIUM": true, "CLAIM": true, "TAX": true, "BROKERAGE_FEE": true,
	"OVR_COMM": true, "SUM_INSURED": true, "CEDING_RETENTION": true, "SUM_REASURED": true,
	"PROF_COMM": true, "GROSS_PREMIUM_REFUND": true, "NET_PREMIUM_REFUND": true, "COMM_REFUND": true,
	"BROKERAGE_FEE_REFUND": true, "OVR_COMM_REFUND": true, "TAX_REFUND": true, "SHARE_RETRO": true,
	"GROSS_PREMIUM_RETRO": true, "DISCOUNT_PREMIUM_RETRO": true, "OVR_COMM_RETRO": true, "BROKERAGE_FEE_RETRO": true,
	"NET_PREMIUM_RETRO": true, "GROSS_PREMIUM_REFUND_RETRO": true, "DISCOUNT_PREMIUM_REFUND_RETRO": true, "OVR_COMM_REFUND_RETRO": true,
	"BROKERAGE_FEE_REFUND_RETRO": true, "NET_PREMIUM_REFUND_RETRO": true, "CLAIM_AMOUNT": true, "RATE": true,
	"SUM_AT_RISK_GROSS": true, "SUM_AT_RISK_RETRO": true, "RETROCEDED_SHARE": true, "SHARE_NUSANTARA_RE_GROSS": true,
	"DEDUCTION": true, "FACTOR": true, "DEDUCTION_REFUND": true, "RI_ADMIN_FEE_REFUND_RETRO": true,
	"RI_ADMIN_FEE_RETRO": true, "RI_ADMIN_FEE_REFUND": true, "RI_ADMIN_FEE": true, "EM_PERCENT": true,
	"RISK": true,
}

// StatusSalinWarisan adalah `CARI48`, VERBATIM.
//
// `[terverifikasi]` `@if(TempWorkPage.Type=="QR","0",@if(TempWorkPage.Type=="QP","0","1"))`
// - `InsertLifePremiumDetail_act.xml`. ⛔ Arti `0`/`1` tidak dijelaskan
// korpus (keluarga OQ-020); yang disalin hanya pemetaannya.
func StatusSalinWarisan(tipe string) string {
	switch tipe {
	case models.TipePLQuotationRealisasi, models.TipePLQuotationProposal:
		return "0"
	}
	return "1"
}

// sqlHapusPesertaWarisan merakit pembersihan salinan warisan satu work.
//
// ⛔ IDEMPOTEN PER `PL_NUMBER` (pl2), dan DIKURUNG `IDPEGA`. `PL_NUMBER` saja
// tidak cukup: baris endorsemen polis yang sama memuat `PL_NUMBER` yang SAMA
// (dengan `PL_NUMBER_EDM` terisi) dan milik work lain. Menghapus berdasarkan
// nomor saja menghapus riwayat endorsemen yang bukan milik submit ini.
//
// ⚠️ `PL_NUMBER` di depan: tabel ini hanya ber-index pada `PL_NUMBER`,
// `CERTIFICATE_NO`, dan `POLICY_NO` (DBA 26-09-2026).
func sqlHapusPesertaWarisan(tabel string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE PL_NUMBER = :1 AND IDPEGA = :2`, tabel)
}

// sqlSisipPesertaWarisan merakit penyisipan satu baris warisan.
//
// ⛔ Tanggal diubah Oracle lewat `TO_DATE` dengan pola yang DINYATAKAN -
// `To_date(…,'DD/MM/YYYY')` persis rule aslinya - bukan `NLS_DATE_FORMAT`.
func sqlSisipPesertaWarisan(tabel, urutan string) string {
	kolom := []string{"ID"}
	penanda := []string{"TO_CHAR(" + urutan + ".NEXTVAL)"}
	for i, k := range kolomPesertaWarisan {
		kolom = append(kolom, k.Kolom)
		switch k.Jenis {
		case nilaiTanggal, nilaiTeksTanggal:
			penanda = append(penanda,
				fmt.Sprintf("TO_DATE(:%d, '%s')", i+1, BentukTanggalOracle))
		default:
			penanda = append(penanda, fmt.Sprintf(":%d", i+1))
		}
	}
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`,
		tabel, strings.Join(kolom, ", "), strings.Join(penanda, ", "))
}

// nilaiSalinWarisan menyusun argumen satu baris, urut `kolomPesertaWarisan`.
//
// ⛔ Kosong menjadi NULL, bukan teks kosong - KECUALI kolom
// `kolomNolBilaKosongWarisan`, yang menjadi "0" seperti `@toDecimal` Pega
// (OQ-PL-10, GILIRAN-17).
func nilaiSalinWarisan(nomorPL, idPega string, b BarisWarisan) []any {
	arg := make([]any, 0, len(kolomPesertaWarisan))
	for _, k := range kolomPesertaWarisan {
		switch k.Kolom {
		case "PL_NUMBER":
			arg = append(arg, nomorPL)
			continue
		case "IDPEGA":
			arg = append(arg, idPega)
			continue
		case "STATUSOLD":
			arg = append(arg, "0")
			continue
		case "STATUS":
			arg = append(arg, StatusSalinWarisan(b.Tipe))
			continue
		}
		v := b.Nilai[k.Kolom]
		if !v.Valid || strings.TrimSpace(v.String) == "" {
			if kolomNolBilaKosongWarisan[k.Kolom] {
				arg = append(arg, "0")
				continue
			}
			arg = append(arg, nil)
			continue
		}
		arg = append(arg, v.String)
	}
	return arg
}

// PesertaWarisan menulis salinan peserta ke tabel warisan.
type PesertaWarisan struct{ db *DB }

// NewPesertaWarisan menyusunnya.
func NewPesertaWarisan(db *DB) *PesertaWarisan { return &PesertaWarisan{db: db} }

// Ganti menghapus salinan lama milik work yang sama lalu menyisipkan yang baru.
//
// `idPega` adalah pengenal work - padanan `pyWorkPage.pzInsKey`. Di repo ini
// pengenal polis yang dipakai seluruh layanan (`T_WORK_POLIS.ID` =
// `T_PREMIUM_LIST.ID`, lihat polis_nomor.go) ADALAH kunci work itu.
//
// ⚠️ SENGAJA BUKAN `T_PREMIUM_LIST.ID_PEGA`. Kolom itu belum punya satu pun
// penulis di repo ini (header polis belum dibuat kode mana pun), sehingga
// membacanya membuat setiap submit gagal - atau, bila NULL dilewatkan,
// `IDPEGA = NULL` tidak pernah benar dan salinan lama tidak pernah terhapus.
//
// Mengembalikan cacah baris terhapus dan tersisip.
func (r *PesertaWarisan) Ganti(ctx context.Context, tx *Tx, nomorPL, idPega string,
	baris []BarisWarisan) (dihapus, disisip int, err error) {

	if strings.TrimSpace(nomorPL) == "" {
		return 0, 0, errors.New("repository: menolak menyalin peserta warisan tanpa PL_NUMBER")
	}
	if strings.TrimSpace(idPega) == "" {
		return 0, 0, errors.New("repository: menolak menyalin peserta warisan tanpa pengenal work")
	}
	if len(baris) == 0 {
		return 0, 0, ErrPolisTanpaPeserta
	}
	tabel, err := r.db.Qualify(namaTabelPeserta)
	if err != nil {
		return 0, 0, err
	}
	urutan, err := r.db.Qualify(namaUrutanPesertaWarisan)
	if err != nil {
		return 0, 0, err
	}
	hapus := sqlHapusPesertaWarisan(tabel)
	if err := PeriksaSQL(hapus); err != nil {
		return 0, 0, err
	}
	h, err := tx.tx.ExecContext(ctx, hapus, nomorPL, idPega)
	if err != nil {
		return 0, 0, fmt.Errorf("repository: menghapus salinan peserta warisan: %w", err)
	}
	nh, err := h.RowsAffected()
	if err != nil {
		return 0, 0, fmt.Errorf("repository: membaca cacah salinan terhapus: %w", err)
	}

	sisip := sqlSisipPesertaWarisan(tabel, urutan)
	if err := PeriksaSQL(sisip); err != nil {
		return 0, 0, err
	}
	stmt, err := tx.tx.PrepareContext(ctx, sisip)
	if err != nil {
		return 0, 0, fmt.Errorf("repository: menyiapkan salinan peserta warisan: %w", err)
	}
	defer func() { _ = stmt.Close() }()
	for _, b := range baris {
		if _, err := stmt.ExecContext(ctx, nilaiSalinWarisan(nomorPL, idPega, b)...); err != nil {
			return 0, 0, fmt.Errorf("repository: menyalin peserta %s ke tabel warisan: %w",
				b.IDPeserta, err)
		}
	}
	return int(nh), len(baris), nil
}

// namaUrutanPesertaWarisan adalah sequence PK `M_LIFE_PREMIUM_DETAIL`.
//
// `[terverifikasi]` `to_Char(M_LIFE_PREMIUM_DETAIL_SEQ.nextval)` - butir
// pertama `VALUES` di `RDBList/SaveMasterLPDet.xml`.
const namaUrutanPesertaWarisan = "M_LIFE_PREMIUM_DETAIL_SEQ"
