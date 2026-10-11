package repository

// Uji teks SQL tanpa Oracle: bentuk WHERE setia rule asal, penampung unik, nol COMMIT, nama berkualifikasi.

import (
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/endorsmentfacin/backend/models"
)

var polaPenampung = regexp.MustCompile(`(?:^|[^A-Za-z0-9_':]):(\d+)\b`)

// penampung - nomor penampung menurut urutan kemunculan.
func penampung(q string) []string {
	var out []string
	for _, m := range polaPenampung.FindAllStringSubmatch(q, -1) {
		out = append(out, m[1])
	}
	return out
}

func rata(q string) string { return strings.Join(strings.Fields(q), " ") }

func TestSQLGerbang(t *testing.T) {
	for _, k := range []struct {
		nama, q   string
		penampung []string
		memuat    []string
	}{
		{"StatusEDM", SqlNilaiPolisTerakhir("S.JSON_POLIS", JalurEdmType), []string{"1", "2"}, []string{
			"JSON_VALUE(b.DATA_JSON, '$.QuotationData.EdmType' NULL ON ERROR)", "FROM S.JSON_POLIS b",
			"b.NOPOLIS = :1 AND b.PRODKE = (SELECT COUNT(c.NOPOLIS) - 1 FROM S.JSON_POLIS c WHERE c.NOPOLIS = :2)",
			"FETCH FIRST 1 ROWS ONLY"}},
		{"JenisBisnis", SqlNilaiPolisTerakhir("S.JSON_POLIS", JalurBusinessType), []string{"1", "2"}, []string{
			"'$.QuotationData.BusinessType'"}},
		{"Klaim", SqlKlaimTerakhir("S.OS_AKSEPTASI_KLAIM"), []string{"1"}, []string{
			"SELECT TO_CHAR(k.STS_REJECT) FROM S.OS_AKSEPTASI_KLAIM k WHERE k.NOPOLIS = :1",
			"ORDER BY k.TANGGAL DESC FETCH FIRST 1 ROWS ONLY"}},
		{"Renewal", SqlAdaRenewal("S.JSON_POLIS"), []string{"1"}, []string{"FROM S.JSON_POLIS j WHERE j.OLDNOPOLIS = :1"}},
		{"Facout", SqlAdaFacout("S.FACOUTPRODUCTION"), []string{"1"}, []string{"FROM S.FACOUTPRODUCTION f WHERE f.POLICYNO = :1"}},
		{"Pembayaran", SqlAdaPembayaran(), []string{"1", "2"}, []string{
			"FROM ARASAPAS.DETAIL_INVOICE WHERE INV_INV_NO = :1 AND IVD_JR_ID = :2"}},
		{"Periode", SqlPeriodeProduksi("S.FACINPRODUCTION"), []string{"1"}, []string{
			"SELECT DISTINCT TO_CHAR(a.BEGINDATE, 'YYYYMMDDHH24MISS'), TO_CHAR(a.ENDDATE, 'YYYYMMDDHH24MISS')",
			"FROM S.FACINPRODUCTION a WHERE a.NOPOLIS = :1",
			"b.SPREAD_DATE = (SELECT MAX(c.SPREAD_DATE) FROM S.FACINPRODUCTION c WHERE c.NOPOLIS = a.NOPOLIS)",
			"FETCH FIRST 1 ROWS ONLY"}},
		{"Proteksi", SqlAdaProteksi("S.OPENPROTEKSI_EDM"), []string{"1", "2", "3"}, []string{
			"FROM S.OPENPROTEKSI_EDM o WHERE o.POLICY_NO = :1 AND o.TYPE = :2",
			"AND (o.STS_AKSEP <> :3 OR o.STS_AKSEP IS NULL) FETCH FIRST 1 ROWS ONLY"}},
		{"EDMBelumSelesai", SqlEDMBelumSelesai("S.T_WORK_POLIS", "S.T_EDM_FACIN_KASUS"), []string{"1", "2", "3"}, []string{
			"SELECT w.ID FROM S.T_WORK_POLIS w JOIN S.T_EDM_FACIN_KASUS k ON k.ID = w.ID",
			"WHERE w.LINI = :1 AND k.OLD_POLICY_NO = :2 AND (w.STATUS_WORK IS NULL OR w.STATUS_WORK NOT LIKE :3)",
			"ORDER BY w.ID FETCH FIRST 1 ROWS ONLY"}},
		{"SisipWorkPolis", SqlSisipWorkPolis("S.T_WORK_POLIS"), []string{"1", "2", "3", "4", "5", "6", "7"}, []string{
			"INSERT INTO S.T_WORK_POLIS (ID, LINI, POSITION, FLAG_ONGOING_POLICY, STATUS_WORK, CREATE_OP, CREATE_OP_NAME, TGL_CREATE, TGL_UPDATE)",
			"VALUES (:1, :2, :3, :4, :5, :6, :7, SYSDATE, SYSDATE)"}},
		{"SisipKasus", SqlSisipKasus("S.T_EDM_FACIN_KASUS"),
			[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15"}, []string{
				"INSERT INTO S.T_EDM_FACIN_KASUS (" + strings.Join(KolomKasus, ", ") + ")",
				"TO_DATE(:6, 'YYYY-MM-DD')"}},
		{"Pengecualian", SqlPengecualian("S.M_EDM_PENGECUALIAN_POLIS"), []string{"1"}, []string{
			"SELECT p.NOPOLIS FROM S.M_EDM_PENGECUALIAN_POLIS p WHERE p.JENIS = :1 ORDER BY p.NOPOLIS"}},
	} {
		q := rata(k.q)
		if err := db.PeriksaSQL(q); err != nil {
			t.Errorf("%s: %v", k.nama, err)
		}
		if got := penampung(q); strings.Join(got, ",") != strings.Join(k.penampung, ",") {
			t.Errorf("%s: penampung %v, ingin %v", k.nama, got, k.penampung)
		}
		for _, s := range k.memuat {
			if !strings.Contains(q, s) {
				t.Errorf("%s: tidak memuat %q\n%s", k.nama, s, q)
			}
		}
		if strings.Contains(q, "SELECT *") {
			t.Errorf("%s: SELECT *", k.nama)
		}
	}
}

func TestKlasifikasiGalatOracle(t *testing.T) {
	for _, k := range []struct {
		pesan               string
		belumAda, skemaLain bool
	}{
		{"ORA-00942: table or view does not exist", true, true},
		{"ORA-00904: invalid identifier", true, false},
		{"ORA-01031: insufficient privileges", false, true},
		{"ORA-03113: end-of-file on communication channel", false, false},
	} {
		asal := errors.New(k.pesan)
		if got := errors.Is(bungkus(asal, "uji"), ErrBelumAda); got != k.belumAda {
			t.Errorf("bungkus %q: ErrBelumAda = %v", k.pesan, got)
		}
		if got := errors.Is(galatArasapas(asal), ErrSkemaLainTakTerbaca); got != k.skemaLain {
			t.Errorf("galatArasapas %q: ErrSkemaLainTakTerbaca = %v", k.pesan, got)
		}
		if !errors.Is(galatArasapas(asal), asal) && !k.skemaLain {
			t.Errorf("galatArasapas %q: galat asal hilang", k.pesan)
		}
	}
	if bungkus(nil, "uji") != nil {
		t.Fatal("bungkus(nil) bukan nil")
	}
}

func TestArgSisipKasus(t *testing.T) {
	k := models.KepalaKasus{ID: "EDM-00001", OldPolicyNo: "UJI.01.0001", EdmType: "4", EdmTypeNew: "1", DetailType: "2",
		EndorsementDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.FixedZone("WIB", 7*3600)), EndorsementReason: "r",
		EndorsementSource: "s", EndorsementSourceNote: "sn", EndorsementStatus: "", AdminCharge: "true", Survey: "false",
		InternalRetro: "0", NBStatus: "NEW EDM", IsCedingConfirm: "Policy"}
	a := ArgSisipKasus(k)
	if len(a) != len(KolomKasus) || len(KolomKasus) != 15 {
		t.Fatalf("%d nilai, %d kolom", len(a), len(KolomKasus))
	}
	want := []any{"EDM-00001", "UJI.01.0001", "4", "1", "2", "2025-06-01", "r", "s", "sn", nil, "true", "false", "0",
		"NEW EDM", "Policy"}
	for i := range want {
		if a[i] != want[i] {
			t.Errorf("%s = %v, ingin %v", KolomKasus[i], a[i], want[i])
		}
	}
	// Tanggal kosong → NULL.
	if ArgSisipKasus(models.KepalaKasus{ID: "EDM-00002"})[5] != nil {
		t.Fatal("tanggal nol bukan NULL")
	}
}

func TestRakitIDKasusEDM(t *testing.T) {
	for n, want := range map[int64]string{1: "EDM-00001", 42: "EDM-00042", 99999: "EDM-99999", 100000: "EDM-100000"} {
		if got := models.RakitIDKasusEDM(n); got != want {
			t.Errorf("%d → %q, ingin %q", n, got, want)
		}
	}
	if models.LiniEDMFacIn != "FACEDM" || models.PosisiAwalEDM != "ReasFacInMarketing" ||
		models.StatusAwalEDM != "Pending-Policy" || models.FlagAwalEDM != "1" || models.PolaStatusSelesai != "Resolved-%" {
		t.Fatal("keadaan awal kasus EDM berubah (keputusan work owner W-EDM-08…11)")
	}
}

func TestNomorInvoice(t *testing.T) {
	// `@replaceAll(.PolicyNo,".","")` - SEMUA titik.
	for masuk, keluar := range map[string]string{"UJI.01.0001": "UJI010001", "TANPATITIK": "TANPATITIK", "..": ""} {
		if got := NomorInvoice(masuk); got != keluar {
			t.Errorf("%q → %q, ingin %q", masuk, got, keluar)
		}
	}
	if KodeJurnalPembayaran != "5" {
		t.Fatalf("kode jurnal %q", KodeJurnalPembayaran)
	}
	// BrowseOpenProteksiEdm_RD saringan C `.STS_AKSEP != "2"`.
	if StsAksepDikecualikan != "2" || TabelProteksi != "OPENPROTEKSI_EDM" {
		t.Fatalf("saringan klep %q / tabel %q", StsAksepDikecualikan, TabelProteksi)
	}
}
