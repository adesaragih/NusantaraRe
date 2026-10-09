package models

// Uji hitung kolom peserta Type QR - Calculate CSV (keputusan work owner 05-10-2026).
// TANPA Oracle; nol nama orang (UJI-).

import (
	"errors"
	"strings"
	"testing"

	"nusantarare/inti/backend/utils"
)

// paramUji - Ceding Retention 25%, Ceding's Limit 100.000.000, RNM Share 50%,
// Deduction 20%, Brokerage 5%; Class of Business UJI-GTL.
func paramUji() ParamProdukQR {
	return ParamProdukQR{ProductName: "UJI-PRODUK", CedingRetentionNum: "25", CedingLimit: "100000000",
		RNMShare: "50", RIComm: "20", Brokerage: "5", RIRiskID: "RISK-1", BusinessName: "UJI-GTL", ProRateType: "1"}
}

var planUji = []PlanProdukQR{{Name: "UJI-LAIN", RIRateID: "RATE-9"}, {Name: " uji-gtl ", RIRateID: "RATE-1"}}

// rateUnisex - AGE 40 CONTRACT 1 → 3,2 (koma desimal, seperti view).
var rateUnisex = []BarisRateQR{
	{ID: "1", Gender: "U", Age: "40", Contract: "1", Rate: "3,2"},
	{ID: "2", Gender: "U", Age: "40", Contract: "2", Rate: "6.4"},
	{ID: "3", Gender: "U", Age: "41", Contract: "1", Rate: "3.3"},
}

var riskUji = []BarisRiskQR{
	{ID: "1", Year: "1", Contract: "1", Risk: "620"},
	{ID: "2", Year: "2", Contract: "1", Risk: "410"},
}

// barisQR - SI 500.000.000, ENTRY_AGE 40, PERIOD_MM 12, tahun pertama, FACTOR
// kosong, EM_PERCENT 0.5.
func barisQR(ubah map[string]string) []BarisUnggah {
	n := map[string]string{
		"SUM_INSURED": "500000000", "ENTRY_AGE": "40", "PERIOD_MM": "12", "SEX": "",
		"BEGIN_DATE": "01/01/2026", "GROSS_VALUATION_BEGIN_DATE": "01/01/2026",
		"FACTOR": "", "EM_PERCENT": "0.5", "RI_ADMIN_FEE": "",
		// Kolom hasil di CSV diabaikan dan ditimpa.
		"GROSS_PREMIUM": "999", "RATE": "999",
	}
	for k, v := range ubah {
		n[k] = v
	}
	return []BarisUnggah{{Nomor: 1, Nilai: n}}
}

func masterUji(t *testing.T, p ParamProdukQR, rate []BarisRateQR) MasterQR {
	t.Helper()
	m, err := SiapkanMasterQR(p, planUji, rate, riskUji)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// samaAngka - teks desimal sama nilainya (4 desimal tersimpan, mis. 960000.0000).
func samaAngka(t *testing.T, kolom, got, mau string) {
	t.Helper()
	g, err1 := utils.ParseDecimal(got)
	w, err2 := utils.ParseDecimal(mau)
	if err1 != nil || err2 != nil || g.Cmp(w) != 0 {
		t.Errorf("%s = %q, mau %s", kolom, got, mau)
	}
}

func hitungSatu(t *testing.T, m MasterQR, ubah map[string]string) (map[string]string, []Penolakan) {
	t.Helper()
	b := barisQR(ubah)
	tolak := HitungPesertaQR(m, b, nil)
	return b[0].Nilai, tolak
}

func TestHitungQRVektorTahunPertama(t *testing.T) {
	n, tolak := hitungSatu(t, masterUji(t, paramUji(), rateUnisex), nil)
	if len(tolak) != 0 {
		t.Fatalf("ditolak: %+v", tolak)
	}
	for k, mau := range map[string]string{
		"CEDING_RETENTION": "100000000", "SUM_REASURED": "400000000", "SHARE_NUSANTARA_RE": "200000000",
		"SUM_AT_RISK_GROSS": "200000000", "RATE": "3.2", "GROSS_PREMIUM": "960000", "DEDUCTION": "192000",
		"BROKERAGE_FEE": "48000", "NET_PREMIUM": "720000", "RI_ADMIN_FEE": "0",
	} {
		samaAngka(t, k, n[k], mau)
	}
	if n["RISK"] != "" {
		t.Errorf("tahun pertama: RISK = %q, mau kosong", n["RISK"])
	}
	if n["GROSS_PREMIUM"] != "960000.0000" {
		t.Errorf("dibulatkan 4 desimal: %q", n["GROSS_PREMIUM"])
	}
}

func TestHitungQRBukanTahunPertama(t *testing.T) {
	// YEAR = 2027 - 2026 = 1 (tanpa +1), CONTRACT 1 → RISK master 620.
	n, tolak := hitungSatu(t, masterUji(t, paramUji(), rateUnisex), map[string]string{"GROSS_VALUATION_BEGIN_DATE": "01/01/2027"})
	if len(tolak) != 0 {
		t.Fatalf("ditolak: %+v", tolak)
	}
	samaAngka(t, "SUM_AT_RISK_GROSS", n["SUM_AT_RISK_GROSS"], "124000000")
	if n["RISK"] != "0.6200000" {
		t.Errorf("RISK = %q, mau 0.6200000 (7 desimal)", n["RISK"])
	}
	// Premi tetap dari SHARE, bukan dari sum at risk.
	samaAngka(t, "GROSS_PREMIUM", n["GROSS_PREMIUM"], "960000")
}

func TestHitungQRLimitTidakTerlampaui(t *testing.T) {
	// SI 200.000.000 × 25% = 50.000.000 < limit.
	n, _ := hitungSatu(t, masterUji(t, paramUji(), rateUnisex), map[string]string{"SUM_INSURED": "200000000"})
	samaAngka(t, "CEDING_RETENTION", n["CEDING_RETENTION"], "50000000")
	samaAngka(t, "SUM_REASURED", n["SUM_REASURED"], "150000000")
	samaAngka(t, "SHARE_NUSANTARA_RE", n["SHARE_NUSANTARA_RE"], "75000000")
}

func TestHitungQRFaktorNolDanEMKosong(t *testing.T) {
	m := masterUji(t, paramUji(), rateUnisex)
	// FACTOR 0 → 1, EM kosong → 0: 3.2/1000 × 200.000.000 = 640.000.
	n, _ := hitungSatu(t, m, map[string]string{"FACTOR": "0", "EM_PERCENT": ""})
	samaAngka(t, "GROSS_PREMIUM", n["GROSS_PREMIUM"], "640000")
	// FACTOR 2, EM 0 → 1.280.000.
	n, _ = hitungSatu(t, m, map[string]string{"FACTOR": "2", "EM_PERCENT": "0"})
	samaAngka(t, "GROSS_PREMIUM", n["GROSS_PREMIUM"], "1280000")
	// RI_ADMIN_FEE CSV dipertahankan, tidak mengurangi NET_PREMIUM.
	n, _ = hitungSatu(t, m, map[string]string{"RI_ADMIN_FEE": "1500"})
	if n["RI_ADMIN_FEE"] != "1500" {
		t.Errorf("RI_ADMIN_FEE = %q", n["RI_ADMIN_FEE"])
	}
	samaAngka(t, "NET_PREMIUM", n["NET_PREMIUM"], "720000")
}

func TestHitungQRKontrak(t *testing.T) {
	for periode, mau := range map[int64]string{1: "1", 11: "1", 12: "1", 17: "1", 18: "2", 24: "2", 30: "3"} {
		if got := KontrakQR(periode).Text('f'); got != mau {
			t.Errorf("PERIOD_MM %d → CONTRACT %s, mau %s", periode, got, mau)
		}
	}
	m := masterUji(t, paramUji(), rateUnisex)
	// PERIOD_MM < 12 → CONTRACT 1 → rate 3.2; 18 → CONTRACT 2 → rate 6.4.
	n, _ := hitungSatu(t, m, map[string]string{"PERIOD_MM": "6"})
	samaAngka(t, "RATE", n["RATE"], "3.2")
	n, _ = hitungSatu(t, m, map[string]string{"PERIOD_MM": "18"})
	samaAngka(t, "RATE", n["RATE"], "6.4")
	for _, v := range []string{"", "0", "1.5", "x"} {
		if _, tolak := hitungSatu(t, m, map[string]string{"PERIOD_MM": v}); len(tolak) != 1 || tolak[0].Kolom != "PERIOD_MM" {
			t.Errorf("PERIOD_MM %q: %+v", v, tolak)
		}
	}
}

var rateGender = []BarisRateQR{
	{ID: "1", Gender: "F", Age: "40", Contract: "1", Rate: "2.5"},
	{ID: "2", Gender: "M", Age: "40", Contract: "1", Rate: "3.5"},
	{ID: "3", Gender: "M", Age: "40", Contract: "2", Rate: "7"},
}

func TestHitungQRMasterGender(t *testing.T) {
	m := masterUji(t, paramUji(), rateGender)
	n, tolak := hitungSatu(t, m, map[string]string{"SEX": "f"})
	if len(tolak) != 0 {
		t.Fatalf("ditolak: %+v", tolak)
	}
	samaAngka(t, "RATE", n["RATE"], "2.5")
	// Pro Rate Type 1: CONTRACT ikut dicocokkan - M, contract 2 (PERIOD_MM 24).
	n, _ = hitungSatu(t, m, map[string]string{"SEX": "M", "PERIOD_MM": "24"})
	samaAngka(t, "RATE", n["RATE"], "7")
	// SEX kosong / asing ditolak.
	for _, v := range []string{"", "X"} {
		if _, tolak := hitungSatu(t, m, map[string]string{"SEX": v}); len(tolak) != 1 || tolak[0].Kolom != "SEX" {
			t.Errorf("SEX %q: %+v", v, tolak)
		}
	}
	// Master U: SEX CSV tidak dibaca.
	n, tolak = hitungSatu(t, masterUji(t, paramUji(), rateUnisex), map[string]string{"SEX": "X"})
	if len(tolak) != 0 || n["RATE"] == "" {
		t.Errorf("master U membaca SEX: %+v", tolak)
	}
}

func TestHitungQRProRate23TanpaKontrak(t *testing.T) {
	for _, pr := range []string{"2", "3"} {
		p := paramUji()
		p.ProRateType = pr
		m := masterUji(t, p, rateGender)
		// F: satu baris (CONTRACT diabaikan) → 2.5, walau PERIOD_MM 24.
		n, tolak := hitungSatu(t, m, map[string]string{"SEX": "F", "PERIOD_MM": "24"})
		if len(tolak) != 0 {
			t.Fatalf("Pro Rate %s: %+v", pr, tolak)
		}
		samaAngka(t, "RATE", n["RATE"], "2.5")
		// M: dua baris AGE 40 → ganda, ditolak.
		_, tolak = hitungSatu(t, m, map[string]string{"SEX": "M"})
		if len(tolak) != 1 || !strings.Contains(tolak[0].Pesan, "More than one rate found for age 40, sex M in R/I Rate RATE-1") {
			t.Errorf("Pro Rate %s ganda: %+v", pr, tolak)
		}
	}
}

func TestHitungQRRateDanRiskTidakKetemu(t *testing.T) {
	m := masterUji(t, paramUji(), rateUnisex)
	_, tolak := hitungSatu(t, m, map[string]string{"ENTRY_AGE": "70"})
	if len(tolak) != 1 || tolak[0].Pesan != "No rate found for age 70, contract 1 in R/I Rate RATE-1" {
		t.Errorf("rate hilang: %+v", tolak)
	}
	// YEAR 5 tidak ada di master risk.
	_, tolak = hitungSatu(t, m, map[string]string{"GROSS_VALUATION_BEGIN_DATE": "01/01/2031"})
	if len(tolak) != 1 || tolak[0].Pesan != "No risk found for year 5, contract 1 in R/I Risk RISK-1" {
		t.Errorf("risk hilang: %+v", tolak)
	}
	// Rate ganda pada master U.
	ganda := append(append([]BarisRateQR{}, rateUnisex...), BarisRateQR{ID: "9", Gender: "U", Age: "40", Contract: "1", Rate: "1"})
	_, tolak = hitungSatu(t, masterUji(t, paramUji(), ganda), nil)
	if len(tolak) != 1 || !strings.HasPrefix(tolak[0].Pesan, "More than one rate found") {
		t.Errorf("rate ganda: %+v", tolak)
	}
	// Risk ganda.
	riskGanda := append(append([]BarisRiskQR{}, riskUji...), BarisRiskQR{ID: "7", Year: "1", Contract: "1", Risk: "1"})
	m2, err := SiapkanMasterQR(paramUji(), planUji, rateUnisex, riskGanda)
	if err != nil {
		t.Fatal(err)
	}
	_, tolak = hitungSatu(t, m2, map[string]string{"GROSS_VALUATION_BEGIN_DATE": "01/01/2027"})
	if len(tolak) != 1 || !strings.HasPrefix(tolak[0].Pesan, "More than one risk found for year 1, contract 1") {
		t.Errorf("risk ganda: %+v", tolak)
	}
	// Baris yang sudah ditolak validasi tidak dihitung.
	b := barisQR(map[string]string{"ENTRY_AGE": "70"})
	if got := HitungPesertaQR(m, b, map[int]bool{1: true}); len(got) != 0 {
		t.Errorf("baris terlewati tetap dihitung: %+v", got)
	}
}

func TestSyaratPolisQR(t *testing.T) {
	cek := func(nama string, p ParamProdukQR, plan []PlanProdukQR, rate []BarisRateQR, risk []BarisRiskQR, pesan string) {
		t.Helper()
		_, err := SiapkanMasterQR(p, plan, rate, risk)
		if !errors.Is(err, ErrSyaratHitungQR) || !strings.Contains(err.Error(), pesan) {
			t.Errorf("%s: %v, mau memuat %q", nama, err, pesan)
		}
	}
	p := paramUji()
	p.CedingLimit = ""
	cek("P1 kosong", p, planUji, rateUnisex, riskUji, "Ceding's Limit is not set on the product")
	p.CedingLimit = "0"
	cek("P1 nol", p, planUji, rateUnisex, riskUji, "Ceding's Limit is not set on the product")
	p = paramUji()
	p.BusinessName = "UJI-TIDAK-ADA"
	cek("P2 nol", p, planUji, rateUnisex, riskUji, "No plan in product UJI-PRODUK matches Class of Business UJI-TIDAK-ADA")
	cek("P2 ganda", paramUji(), append(append([]PlanProdukQR{}, planUji...), PlanProdukQR{Name: "UJI-GTL", RIRateID: "R2"}),
		rateUnisex, riskUji, "More than one plan in product UJI-PRODUK matches Class of Business UJI-GTL")
	campur := append(append([]BarisRateQR{}, rateUnisex...), rateGender[0])
	cek("P3 campuran", paramUji(), planUji, campur, riskUji, "R/I Rate RATE-1 mixes unisex and gendered rows")
	cek("P4 rate", paramUji(), []PlanProdukQR{{Name: "UJI-GTL"}}, rateUnisex, riskUji, "has no R/I Rate")
	p = paramUji()
	p.RIRiskID = " "
	cek("P4 risk", p, planUji, rateUnisex, riskUji, "has no R/I Risk")
	rusak := []BarisRateQR{{ID: "77", Gender: "U", Age: "40", Contract: "1", Rate: "1.234,5"}}
	cek("rate tak terurai", paramUji(), planUji, rusak, riskUji, "row ID 77: RATE")
}

func TestValidasiQRAbaikanKolomHasil(t *testing.T) {
	b := barisSah(1)
	b.Nilai["GROSS_PREMIUM"] = "1,5" // bentuk salah - untuk QR diabaikan (ditimpa hitungan)
	if h, err := ValidasiUnggah("QR", []BarisUnggah{b}); err != nil || !h.Lolos() {
		t.Errorf("QR menolak kolom hasil: %v %+v", err, h.Ditolak)
	}
	if h, _ := ValidasiUnggah("QP", []BarisUnggah{b}); h.Lolos() {
		t.Error("QP tidak lagi memeriksa bentuk GROSS_PREMIUM")
	}
}

// RISK master dibaca sebagai teks: dulu view JSON berkoma ("921,9"), sejak migrasi inti 940 (ririsklife, keputusan work
// owner 08-10-2026 K2) NUMBER dibaca TM9 bertitik ("921.9", pecahan tanpa nol depan ".9"). Ketiga bentuk memberi hasil
// hitung QR yang SAMA - SUM_AT_RISK_GROSS, RISK, dan seluruh kolom hasil.
func TestHitungQRRiskKomaDanTitikSama(t *testing.T) {
	bentuk := map[string][]BarisRiskQR{
		"koma (view lama)":   {{ID: "1", Year: "1", Contract: "1", Risk: "921,9"}, {ID: "2", Year: "2", Contract: "1", Risk: "0,9"}},
		"titik (TM9 NUMBER)": {{ID: "1", Year: "1", Contract: "1", Risk: "921.9"}, {ID: "2", Year: "2", Contract: "1", Risk: ".9"}},
	}
	var acuan map[string]string
	for nama, risk := range bentuk {
		m, err := SiapkanMasterQR(paramUji(), planUji, rateUnisex, risk)
		if err != nil {
			t.Fatalf("%s: %v", nama, err)
		}
		for _, tgl := range []string{"01/01/2027", "01/01/2028"} {
			n, tolak := hitungSatu(t, m, map[string]string{"GROSS_VALUATION_BEGIN_DATE": tgl})
			if len(tolak) != 0 {
				t.Fatalf("%s %s ditolak: %+v", nama, tgl, tolak)
			}
			if acuan == nil {
				acuan = map[string]string{}
			}
			for k, v := range n {
				kunci := tgl + "|" + k
				if lama, ada := acuan[kunci]; ada && lama != v {
					t.Errorf("%s %s: %s = %q, bentuk lain %q", nama, tgl, k, v, lama)
				}
				acuan[kunci] = v
			}
		}
	}
	if acuan["01/01/2027|RISK"] == "" || acuan["01/01/2028|RISK"] == "" {
		t.Errorf("RISK tahun ke-2 / ke-3 kosong: %v", acuan)
	}
	samaAngka(t, "RISK 921,9 / 1000", acuan["01/01/2027|RISK"], "0.9219")
}
