package services

// Format unggah ceding berpemisah titik koma - keputusan work owner
// 02-10-2026. TANPA Oracle. Seluruh nilai uji berawalan `UJI-`; nol nama
// orang dan nol nomor polis nyata.

import (
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/premiumlistlife/backend/models"
)

// judulFormatCeding - baris judul format ceding (31 kolom), TANPA
// START_DATE, EFFECTIVE_DATE, dan dua RETROCESSION_VALUATION_*, dan dengan
// SHARE_NUSANTARA_RE_GROSS menggantikan SHARE_NUSANTARA_RE.
const judulFormatCeding = "NO;POLICY_NO;POLICY_HOLDER;CERTIFICATE_NO;NAME_OF_INSURED;SEX;DOB;" +
	"ENTRY_AGE;CURRENT_AGE;PLAN;BEGIN_DATE;EXPIRED_DATE;GROSS_VALUATION_BEGIN_DATE;" +
	"GROSS_VALUATION_EXPIRED_DATE;PERIOD_MM;UW_STATUS;EM_PERCENT;CURRENCY;SUM_INSURED;" +
	"CEDING_RETENTION;SUM_REASURED;SHARE_NUSANTARA_RE_GROSS;SUM_AT_RISK_GROSS;" +
	"GROSS_PREMIUM;DEDUCTION;RI_ADMIN_FEE;BROKERAGE_FEE;NET_PREMIUM;FACTOR;STNC;WPC"

const barisFormatCeding = `1;UJI-POL-1;UJI PEMEGANG;UJI-C1;UJI PESERTA A;F;08/03/1990;34;34;` +
	`UJI-PLAN;13/12/2024;13/06/2025;13/12/2024;13/06/2025;6;uw status;;IDR;10000;1000;10000;` +
	`10000;10000;"10,70944011";0;"0,535472006";0;"10,1739681";;31/10/2026;31/10/2026`

func TestFormatCedingTitikKomaDiterima(t *testing.T) {
	baris, err := bacaQR(strings.NewReader(judulFormatCeding + "\r\n" + barisFormatCeding + "\r\n"))
	if err != nil {
		t.Fatalf("format ceding ditolak: %v", err)
	}
	if len(baris) != 1 {
		t.Fatalf("%d baris, mau 1", len(baris))
	}
	n := baris[0].Nilai
	for k, mau := range map[string]string{
		"GROSS_PREMIUM":      "10.70944011",
		"NET_PREMIUM":        "10.1739681",
		"RI_ADMIN_FEE":       "0.535472006",
		"DEDUCTION":          "0",
		"SHARE_NUSANTARA_RE": "10000",
		"CERTIFICATE_NO":     "UJI-C1",
	} {
		if n[k] != mau {
			t.Errorf("%s = %q, mau %q", k, n[k], mau)
		}
	}
	models.IsiNolUangKosong(baris)
	if h := validasiQR(baris); !h.Lolos() {
		t.Errorf("baris format ceding ditolak validasi: %+v", h.Ditolak)
	}
}

func TestPemisahCSVDitebakDariJudul(t *testing.T) {
	for _, k := range []struct {
		awal string
		mau  rune
	}{
		{"A;B;C\n1,5;2;3", ';'},
		{"A,B,C\n1;2;3", ','},
		{`"A;B",C,D`, ','},
		{"A", ','},
		{"", ','},
	} {
		if got := PemisahCSV(k.awal); got != k.mau {
			t.Errorf("PemisahCSV(%q) = %q, mau %q", k.awal, got, k.mau)
		}
	}
}

// TestKomaTetapDitolakDiBerkasBerpemisahKoma - `1,234` di berkas
// berpemisah koma bisa berarti seribu; ia tetap DITOLAK, bukan ditebak.
func TestKomaTetapDitolakDiBerkasBerpemisahKoma(t *testing.T) {
	r := barisUji("UJI-C1")
	// SUM_INSURED, bukan GROSS_PREMIUM: untuk QR GROSS_PREMIUM kini DIHITUNG dan
	// bentuk nilainya di CSV tidak diperiksa (05-10-2026).
	r[17] = `"1000,50"`
	baris, err := bacaQR(strings.NewReader(csvUji(judulLengkap(), r)))
	if err != nil {
		t.Fatal(err)
	}
	if baris[0].Nilai["SUM_INSURED"] != "1000,50" {
		t.Errorf("koma di berkas berpemisah koma diubah: %q", baris[0].Nilai["SUM_INSURED"])
	}
	if h := validasiQR(baris); h.Lolos() {
		t.Error("koma di berkas berpemisah koma lolos validasi")
	}
}

func TestShareGrossDidahulukanLaluShare(t *testing.T) {
	judul := judulLengkap("SHARE_NUSANTARA_RE_GROSS")
	isi := barisUji("UJI-C1")
	dua := append(append([]string{}, isi...), "77.5")
	satu := append(append([]string{}, isi...), "")
	baris, err := bacaQR(strings.NewReader(csvUji(judul, dua, satu)))
	if err != nil {
		t.Fatal(err)
	}
	if got := baris[0].Nilai["SHARE_NUSANTARA_RE"]; got != "77.5" {
		t.Errorf("GROSS terisi: SHARE_NUSANTARA_RE = %q, mau 77.5", got)
	}
	if got := baris[1].Nilai["SHARE_NUSANTARA_RE"]; got != "50.5" {
		t.Errorf("GROSS kosong: SHARE_NUSANTARA_RE = %q, mau 50.5 (kolom SHARE_NUSANTARA_RE)", got)
	}
}

func TestTanpaShareTidakLagiWajib(t *testing.T) {
	// Sejak 03-10-2026 SHARE_NUSANTARA_RE bukan kolom wajib (keputusan work
	// owner): berkas tanpanya diterima, dan sel uangnya menjadi 0.
	var judul []string
	var isi []string
	for i, k := range judulLengkap() {
		if k != "SHARE_NUSANTARA_RE" {
			judul = append(judul, k)
			isi = append(isi, barisUji("UJI-C1")[i])
		}
	}
	baris, err := bacaQR(strings.NewReader(csvUji(judul, isi)))
	if err != nil {
		t.Fatalf("berkas tanpa SHARE_NUSANTARA_RE ditolak: %v", err)
	}
	models.IsiNolUangKosong(baris)
	if h := validasiQR(baris); !h.Lolos() {
		t.Errorf("baris tanpa SHARE_NUSANTARA_RE ditolak: %+v", h.Ditolak)
	}
}

// Type kosong ditolak SEBELUM berkasnya dibaca (keputusan work owner 03-10-2026).
func TestBacaCSVTanpaTipeDitolak(t *testing.T) {
	_, err := BacaCSVUnggah(strings.NewReader(csvUji(judulLengkap(), barisUji("UJI-C1"))), "")
	if !errors.Is(err, models.ErrTipeUnggahTakDikenal) {
		t.Errorf("galat %v, mau ErrTipeUnggahTakDikenal", err)
	}
}
