package services_test

import (
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/bordereaux/backend/models"
	"nusantarare/modul/bordereaux/backend/services"
	"nusantarare/modul/bordereaux/backend/tiruan"
)

// Nilai JSON Pega: DateTime GMT dibaca di WIB (tanggal Pega = 17:00 GMT hari sebelumnya), yyyyMMdd apa adanya, epoch =
// kosong; angka bertitik banyak atau berkoma = format Indonesia, selain itu desimal titik.
func TestNilaiJSONLama(t *testing.T) {
	for masuk, mau := range map[string]string{
		"20260630T170000.000 GMT": "01-07-2026", "20260701": "01-07-2026", "19700101T000000.000 GMT": "", "": "", "01/07/2026": "01-07-2026",
	} {
		if got, ok := services.TanggalLama(masuk); !ok || got != mau {
			t.Errorf("TanggalLama(%q) = %q %v, mau %q", masuk, got, ok, mau)
		}
	}
	if _, ok := services.TanggalLama("bukan tanggal"); ok {
		t.Error("teks bukan tanggal harus ditolak")
	}
	if got := services.WaktuLama("20261003T023005.123 GMT"); got != "03-10-2026 09:30:05" {
		t.Errorf("WaktuLama = %q", got)
	}
	for masuk, mau := range map[string]string{
		"1.234.567": "1234567", "1.234,5": "1234.5", "1234.50": "1234.5", "12,5%": "12.5", "-": "", "": "", "0": "0", "-7.25": "-7.25",
	} {
		if got, ok := services.AngkaLama(masuk); !ok || got != mau {
			t.Errorf("AngkaLama(%q) = %q %v, mau %q", masuk, got, ok, mau)
		}
	}
	if _, ok := services.AngkaLama("abc"); ok {
		t.Error("teks bukan angka harus ditolak")
	}
}

// Kolom dibaca dari nama persis, nama tanpa garis bawah (Engineering), atau nama lama (aliasLama).
func TestBarisLamaMemetakanKunciPega(t *testing.T) {
	k := kombinasi(t, models.TypePremium, "ENGINEERING")
	b, err := services.BarisLama(k, []map[string]any{{
		"POLICYNUMBER": "UJI-POL-1", "START_POI": "20260630T170000.000 GMT", "TSI": "1.500.000", "PREMIUM": "2500.5",
		"PREMIUM_CEDED100": "1.000,25", "ZIPCODEE": "12345", "AMOUNT": "10", "pxObjClass": "UJI",
	}})
	if err != nil {
		t.Fatal(err)
	}
	mau := map[string]string{"POLICY_NUMBER": "UJI-POL-1", "START_POI": "01-07-2026", "TSI": "1500000", "PREMIUM_100": "2500.5",
		"PREMIUM_REINSURER": "1000.25", "ZIP_CODE": "12345", "LOL_PML_EML": "10", "COB": ""}
	for kolom, v := range mau {
		if b[0][kolom] != v {
			t.Errorf("%s = %q, mau %q", kolom, b[0][kolom], v)
		}
	}
	_, err = services.BarisLama(k, []map[string]any{{"TSI": "satu juta"}})
	var g services.GalatCSV
	if !errors.As(err, &g) || !strings.Contains(g.Error(), `PremiumEngineeringList row 1, column`) {
		t.Errorf("nilai salah harus ditolak: %v", err)
	}
}

const jsonTanpaHeader = `{"BDX_ID":"BDX-UJI.2","TYPE":"PREMIUM","TYPE_BUSINESS":"FIRE","MASTERID":"UJI-TI-1","CEDINGID":"UJI-AG-1",
"CEDINGNAME":"UJI CEDANT SATU","TREATYNAME":"UJI KONTRAK SATU","BDXREPORT_START":"20260630T170000.000 GMT",
"BDXREPORT_END":"20260929T170000.000 GMT","POSITION":"","STATUSAKSEP":"Resolve-Complete","pxCreateOperator":"UJI-MAKER-LAMA",
"pxCreateDateTime":"20260901T030000.000 GMT",
"PremiumFireList":[{"COB":"UJI","TSI100":"1.500.000","POI_START":"20260630T170000.000 GMT"},{"COB":"UJI-2","TSI100":"7.5"}],
"PremiumEngineeringList":[{"COB":"BUKAN KOMBINASINYA"}],
"CommentList":[{"Date":"20261003T023005.123 GMT","OperatorName":"UJI OPERATOR","IsApproved":"1","Suggest":"UJI OK"},
{"Date":"20261003T040000.000 GMT","OperatorName":"","IsApproved":"0","Suggest":""}]}`

func gudangLama() *tiruan.Gudang {
	g := tiruan.Contoh()
	g.JSON["BDX-UJI.2"] = jsonTanpaHeader
	// BDX-UJI.1: header ada, detailnya masih di JSON saja.
	g.JSON["BDX-UJI.1"] = `{"BDX_ID":"BDX-UJI.1","TYPE":"PREMIUM","TYPE_BUSINESS":"FIRE","PremiumFireList":[{"COB":"UJI-LAMA"}]}`
	// BDX-UJI.3: detail dan riwayat sudah di tabel - tidak tampil, tidak ditimpa.
	g.Header["BDX-UJI.3"] = models.Header{BdxID: "BDX-UJI.3", Type: "PREMIUM", TypeBusiness: "FIRE", UserInput: "UJI-MAKER", Position: "UJI-MAKER"}
	g.Rinci["BORDEREAUX_PREMI_FIRE"] = map[string][]models.Baris{"BDX-UJI.3": {{models.KolomID: "BDX_DTL-UJI", "COB": "BARU"}}}
	g.Riw["BDX-UJI.3"] = []models.Riwayat{{PIC: "UJI-MAKER", IsApproved: true}}
	g.JSON["BDX-UJI.3"] = `{"BDX_ID":"BDX-UJI.3","TYPE":"PREMIUM","TYPE_BUSINESS":"FIRE","PremiumFireList":[{"COB":"LAMA"}],
"CommentList":[{"Date":"20261003T023005.123 GMT","OperatorName":"UJI OPERATOR","IsApproved":"1","Suggest":"LAMA"}]}`
	return g
}

func TestCopyOldHanyaSuperadmin(t *testing.T) {
	l := layanan(gudangLama())
	if _, err := l.DaftarLama(ctx, maker); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("daftar oleh bukan superadmin: %v", err)
	}
	if _, err := l.SalinLama(ctx, maker, []string{"BDX-UJI.2"}); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("salin oleh bukan superadmin: %v", err)
	}
	if services.DaftarPilihan(maker).CopyOld || !services.DaftarPilihan(super).CopyOld {
		t.Error("tombol Copy Old Data hanya untuk superadmin")
	}
}

func TestDaftarLamaHanyaYangBelumDiTabel(t *testing.T) {
	d, err := layanan(gudangLama()).DaftarLama(ctx, super)
	if err != nil {
		t.Fatal(err)
	}
	if len(d) != 2 || d[0].ID != "BDX-UJI.2" || d[1].ID != "BDX-UJI.1" {
		t.Fatalf("daftar %+v", d)
	}
	mau := models.BerkasLama{ID: "BDX-UJI.2", Type: "PREMIUM", Business: "FIRE", Ceding: "UJI CEDANT SATU", Treaty: "UJI KONTRAK SATU",
		Status: "Resolve-Complete", BarisJSON: 2, Komentar: 2, TanpaHeader: true}
	if d[0] != mau {
		t.Errorf("berkas tanpa header %+v", d[0])
	}
	if d[1].BarisJSON != 1 || d[1].BarisTabel != 0 || d[1].TanpaHeader {
		t.Errorf("berkas ber-header %+v", d[1])
	}
}

func TestSalinLama(t *testing.T) {
	g := gudangLama()
	l := layanan(g)
	j, err := l.SalinLama(ctx, super, []string{"BDX-UJI.2", " BDX-UJI.1", "BDX-UJI.3", "BDX-UJI.2", "TIDAK-ADA"})
	if err != nil {
		t.Fatal(err)
	}
	status := []string{}
	for _, h := range j.Hasil {
		status = append(status, h.ID+"="+h.Status)
	}
	if strings.Join(status, " ") != "BDX-UJI.2=disalin BDX-UJI.1=disalin BDX-UJI.3=sudahAda TIDAK-ADA=ditolak" || j.Disalin != 2 {
		t.Fatalf("hasil %v disalin %d", status, j.Disalin)
	}
	if p := strings.Join(j.Hasil[0].Pesan, "; "); p != "header created from the old JSON; 2 detail rows; 2 history rows" {
		t.Errorf("catatan %q", p)
	}
	h := g.Header["BDX-UJI.2"]
	if h.Tanggal != "01-09-2026 10:00:00" || h.UserInput != "UJI-MAKER-LAMA" || h.ReportStart != "01-07-2026" ||
		h.ReportEnd != "30-09-2026" || h.StatusAksep != models.StatusResolveComplete || h.CedingID != "UJI-AG-1" {
		t.Errorf("header %+v", h)
	}
	b := g.Rinci["BORDEREAUX_PREMI_FIRE"]["BDX-UJI.2"]
	if len(b) != 2 || b[0]["TSI100"] != "1500000" || b[0]["POI_START"] != "01-07-2026" || b[1]["TSI100"] != "7.5" ||
		!strings.HasPrefix(b[0][models.KolomID], "BDX_DTL-") || b[0][models.KolomID] == b[1][models.KolomID] {
		t.Errorf("detail %v", b)
	}
	if len(g.Rinci["BORDEREAUX_PREMI_ENGINEERING"]["BDX-UJI.2"]) != 0 {
		t.Error("page list kombinasi lain tidak disalin")
	}
	r := g.Riw["BDX-UJI.2"]
	if len(r) != 2 || r[0] != (models.Riwayat{Tanggal: "03-10-2026 09:30", PIC: "UJI OPERATOR", IsApproved: true, Komentar: "UJI OK"}) ||
		r[1].PIC != "-" || r[1].IsApproved {
		t.Errorf("riwayat %+v", r)
	}
	if b := g.Rinci["BORDEREAUX_PREMI_FIRE"]["BDX-UJI.3"]; len(b) != 1 || b[0]["COB"] != "BARU" || len(g.Riw["BDX-UJI.3"]) != 1 {
		t.Error("isi yang sudah di tabel tidak boleh ditimpa")
	}
	if d, _ := l.DaftarLama(ctx, super); len(d) != 0 {
		t.Errorf("sesudah disalin daftar kosong: %+v", d)
	}
	if j, _ := l.SalinLama(ctx, super, []string{"BDX-UJI.2"}); j.Hasil[0].Status != models.SalinSudahAda {
		t.Errorf("salin ulang %+v", j.Hasil)
	}
	if _, err := l.SalinLama(ctx, super, []string{" "}); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("tanpa ID: %v", err)
	}
}

// Nilai yang tidak sah menolak SELURUH berkas: header, detail, dan riwayatnya tidak ditulis sama sekali.
func TestSalinLamaDitolakTanpaJejak(t *testing.T) {
	g := tiruan.Contoh()
	g.JSON["BDX-UJI.9"] = `{"BDX_ID":"BDX-UJI.9","TYPE":"PREMIUM","TYPE_BUSINESS":"FIRE","PremiumFireList":[{"TSI100":"satu juta"}],
"CommentList":[{"Date":"20261003T023005.123 GMT","OperatorName":"UJI","IsApproved":"1","Suggest":""}]}`
	g.JSON["BDX-UJI.8"] = `bukan json`
	j, err := layanan(g).SalinLama(ctx, super, []string{"BDX-UJI.9", "BDX-UJI.8"})
	if err != nil {
		t.Fatal(err)
	}
	if j.Hasil[0].Status != models.SalinDitolak || !strings.Contains(j.Hasil[0].Pesan[0], "is not a number") ||
		j.Hasil[1].Status != models.SalinDitolak || j.Disalin != 0 {
		t.Errorf("hasil %+v", j.Hasil)
	}
	if _, ada := g.Header["BDX-UJI.9"]; ada || len(g.Riw["BDX-UJI.9"]) != 0 {
		t.Error("berkas ditolak tidak boleh meninggalkan header atau riwayat")
	}
}

// Kolom spread Aviation berisi nominal (migrasi 892, NUMBER(38,8)): nilai miliaran diterima. Dulu NUMBER(10,4) - Copy
// Old Data klaim Aviation DEV ditolak "is too large".
func TestBarisLamaSpreadAviationNominal(t *testing.T) {
	k := kombinasi(t, models.TypeClaim, "AVIATION")
	b, err := services.BarisLama(k, []map[string]any{{"SPREAD_OF_CLAIM_OR": "2500000000", "SPREAD_OF_CLAIM_QS": "2500000000.5"}})
	if err != nil {
		t.Fatal(err)
	}
	if b[0]["SPREAD_OF_CLAIM_OR"] != "2500000000" || b[0]["SPREAD_OF_CLAIM_QS"] != "2500000000.5" {
		t.Errorf("spread %v", b[0])
	}
	for _, kk := range []models.KombinasiBdx{k, kombinasi(t, models.TypePremium, "AVIATION")} {
		for _, c := range kk.Kolom {
			if strings.HasPrefix(c.Kolom, "SPREAD_OF_") && (c.Presisi != 38 || c.Skala != 8) {
				t.Errorf("%s.%s = NUMBER(%d,%d), mau NUMBER(38,8) sesuai migrasi 892", kk.Tabel, c.Kolom, c.Presisi, c.Skala)
			}
		}
	}
}
