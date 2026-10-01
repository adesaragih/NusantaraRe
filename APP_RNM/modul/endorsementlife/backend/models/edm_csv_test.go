package models

import (
	"html"
	"os"
	"regexp"
	"strings"
	"testing"
)

// TestPemetaanCSVDariKorpus - 72 penetapan `SaveCSVEDMLife` 4.1, urutan,
// nama properti, dan normalisasi koma diturunkan ulang dari korpus.
func TestPemetaanCSVDariKorpus(t *testing.T) {
	// Berkas MENTAH: `&lt;APPEND&gt;` dibuka-escape sesudah pasangan dipotong.
	mentah, err := os.ReadFile(akarKorpus + `\Activity\SaveCSVEDMLife.xml`)
	if err != nil {
		t.Skipf("korpus tidak terjangkau (%v)", err)
	}
	pasangan := pasanganProperti(string(mentah), "<pyStepsDescription>Set property PremiumListDetail</pyStepsDescription>",
		"<pyStepsDescription>Set message error</pyStepsDescription>")
	for i := range pasangan {
		pasangan[i] = [2]string{html.UnescapeString(pasangan[i][0]), html.UnescapeString(pasangan[i][1])}
	}
	tujuan := regexp.MustCompile(`^pyWorkPage\.PremiumListSummary\.PremiumListDetail\(<(APPEND|LAST)>\)\.(\w+)$`)
	uang := regexp.MustCompile(`^@divide\(@toDecimal\(@replaceAll\(\.(\w+),",","\."\)\),1,20\)$`)
	polos := regexp.MustCompile(`^\.(\w+)$`)
	var korpus []KolomCSV
	for _, p := range pasangan {
		m := tujuan.FindStringSubmatch(p[0])
		if m == nil {
			continue // `local.errmsg`
		}
		k := KolomCSV{Properti: m[2]}
		if u := uang.FindStringSubmatch(p[1]); u != nil && u[1] == m[2] {
			k.Jenis = CSVUang
		} else if v := polos.FindStringSubmatch(p[1]); v != nil && v[1] == m[2] {
			k.Jenis = CSVTeks
		} else {
			t.Errorf("%s <= %s: bukan salin langsung maupun uang", p[0], p[1])
		}
		korpus = append(korpus, k)
	}
	if len(korpus) != 72 || len(PemetaanCSV) != 72 {
		t.Fatalf("korpus %d, kode %d penetapan; mau 72", len(korpus), len(PemetaanCSV))
	}
	for i, k := range korpus {
		kode := PemetaanCSV[i]
		if kode.Properti != k.Properti {
			t.Errorf("#%d properti %s, korpus %s", i+1, kode.Properti, k.Properti)
		}
		if (kode.Jenis == CSVUang) != (k.Jenis == CSVUang) {
			t.Errorf("%s: normalisasi koma kode %v, korpus %v", k.Properti, kode.Jenis == CSVUang, k.Jenis == CSVUang)
		}
	}
}

// TestKolomCSVBertipeSesuaiSaveMasterLPDet - kolom bertanggal adalah yang
// `SaveMasterLPDet` bungkus `To_date(…, 'DD/MM/YYYY')`.
func TestKolomCSVBertipeSesuaiSaveMasterLPDet(t *testing.T) {
	isi := bacaKorpus(t, `RDBList\SaveMasterLPDet.xml`)
	var korpus []string
	for _, m := range regexp.MustCompile(`To_date\(\{TempValue\.(\w+)\}, 'DD/MM/YYYY'\)`).FindAllStringSubmatch(isi, -1) {
		korpus = append(korpus, m[1])
	}
	var kode []string
	for _, k := range PemetaanCSV {
		if k.Jenis == CSVTanggal || k.Jenis == CSVTeksTanggal {
			kode = append(kode, k.Properti)
		}
	}
	urut := func(s []string) string {
		c := append([]string{}, s...)
		for i := range c {
			for j := i + 1; j < len(c); j++ {
				if c[j] < c[i] {
					c[i], c[j] = c[j], c[i]
				}
			}
		}
		return strings.Join(c, ",")
	}
	if urut(korpus) != urut(kode) {
		t.Errorf("kolom tanggal\n kode   %v\n korpus %v", kode, korpus)
	}
}

func TestUangCSVEDM(t *testing.T) {
	for masuk, mau := range map[string]string{
		"1234.5": "1234.5", "1234,5": "1234.5", "-0,00000001": "-0.00000001", "": "0", " 7 ": "7", "0": "0",
	} {
		if g, err := UangCSVEDM(masuk); err != nil || g != mau {
			t.Errorf("UangCSVEDM(%q) = %q, %v; mau %q", masuk, g, err, mau)
		}
	}
	// Pemisah ribuan: Pega menjadikannya `1.234.567.89` (rusak) - di sini ditolak terang.
	for _, rusak := range []string{"1,234,567.89", "1.234.567,89", "1.234,5", "1,234.5", "1.2.3", "abc", "1e5", "0.000000001", "1,2,3"} {
		if g, err := UangCSVEDM(rusak); err == nil {
			t.Errorf("UangCSVEDM(%q) diterima = %q", rusak, g)
		}
	}
}

func TestRapikanBarisCSV(t *testing.T) {
	acuan := AcuanCSV{Plan: "UJI-PLAN", PolicyHolder: "UJI-PH"}
	nilai, pesan := RapikanBarisCSV(1, map[string]string{
		"PLAN": "UJI-PLAN", "POLICY_HOLDER": "UJI-PH", "GROSS_PREMIUM": "10,5", "DOB": "02/01/1990", "AGE": "36",
		"STNC": "03/04/2026", "FACTOR": "1.25", "UW_STATUS": "UJI", "RETROCESSION_VALUATION_BEGIN_DATE": "01/01/2026",
	}, acuan)
	if len(pesan) != 0 {
		t.Fatalf("baris sah ditolak: %v", pesan)
	}
	for k, mau := range map[string]string{
		"GROSS_PREMIUM": "10.5", "DOB": "1990-01-02", "AGE": "36", "STNC": "03/04/2026", "FACTOR": "1.25",
		"RETRO_VALUATION_BEGIN_DATE": "2026-01-01", "NET_PREMIUM": "0", "LAPSE_DATE": "", "PLAN": "UJI-PLAN",
	} {
		if nilai[k] != mau {
			t.Errorf("%s = %q, mau %q", k, nilai[k], mau)
		}
	}
	if _, ada := nilai["UW_STATUS"]; ada {
		t.Error("UW_STATUS tidak punya kolom; tidak boleh ikut disimpan")
	}

	_, pesan = RapikanBarisCSV(7, map[string]string{
		"PLAN": "UJI-LAIN", "POLICY_HOLDER": "UJI-PH", "GROSS_PREMIUM": "1,234.5", "DOB": "1990-01-02", "AGE": "3x",
		"FACTOR": "1,5", "CERTIFICATE_NO": strings.Repeat("X", 256),
	}, acuan)
	ditemukan := map[string]string{}
	for _, p := range pesan {
		if p.Baris != 7 {
			t.Errorf("pesan tanpa nomor baris: %+v", p)
		}
		ditemukan[p.Kolom] = p.Pesan
	}
	if ditemukan["PLAN"] != PesanCSVPlan {
		t.Errorf("PLAN beda: %q, mau VERBATIM %q", ditemukan["PLAN"], PesanCSVPlan)
	}
	for _, k := range []string{"GROSS_PREMIUM", "DOB", "AGE", "FACTOR", "CERTIFICATE_NO"} {
		if ditemukan[k] == "" {
			t.Errorf("kolom %s tidak dilaporkan (%v)", k, pesan)
		}
	}
	if _, pesan = RapikanBarisCSV(2, map[string]string{"PLAN": "UJI-PLAN", "POLICY_HOLDER": "UJI-LAIN"}, acuan); len(pesan) != 1 ||
		pesan[0].Kolom != "POLICY_HOLDER" || pesan[0].Pesan != PesanCSVPlan {
		t.Errorf("POLICY_HOLDER beda: %v", pesan)
	}
}

func TestJudulCSV(t *testing.T) {
	j, asing, err := JudulCSV([]string{string(rune(0xFEFF)) + "plan", " Policy_Holder ", "GROSS_PREMIUM", "OVRR_COMM", "STATUS"})
	if err != nil || strings.Join(j, ",") != "PLAN,POLICY_HOLDER,GROSS_PREMIUM,OVRR_COMM,STATUS" {
		t.Fatalf("%v %v", j, err)
	}
	// Kepala unduhan `Generate Data Detail` b276 memuat dua properti di luar 4.1: diabaikan, dilaporkan.
	if strings.Join(asing, ",") != "OVRR_COMM,STATUS" {
		t.Errorf("diabaikan %v", asing)
	}
	for _, rusak := range [][]string{
		{"PLAN"},                          // POLICY_HOLDER hilang - acuan 4.1 tak dapat diperiksa
		{"PLAN", "POLICY_HOLDER", "PLAN"}, // ganda
	} {
		if _, _, err := JudulCSV(rusak); err == nil {
			t.Errorf("%v diterima", rusak)
		}
	}
}
