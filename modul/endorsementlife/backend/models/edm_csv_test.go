package models

import (
	"errors"
	"fmt"
	"html"
	"os"
	"regexp"
	"slices"
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

// barisLengkapQR - baris CSV QR yang memenuhi seluruh kolom wajib (10-10-2026), ditimpa `tambah`.
func barisLengkapQR(tambah map[string]string) map[string]string {
	b := map[string]string{
		"NAME_OF_INSURED": "UJI-NAMA", "CERTIFICATE_NO": "UJI-N1", "PLAN": "UJI-PLAN", "POLICY_HOLDER": "UJI-PH", "SEX": "M",
		"DOB": "02/01/1990", "CURRENCY": "IDR", "BEGIN_DATE": "01/01/2026", "EXPIRED_DATE": "31/12/2026", "SUM_INSURED": "1000",
		"AGE": "36", "GROSS_PREMIUM": "10", "DEDUCTION": "0", "RI_ADMIN_FEE": "0", "BROKERAGE_FEE": "0", "TAX": "0",
		"PROF_COMM": "0", "CLAIM": "0", "NET_PREMIUM": "0",
	}
	for k, v := range tambah {
		b[k] = v
	}
	return b
}

func TestRapikanBarisCSV(t *testing.T) {
	acuan := AcuanCSV{Plan: "UJI-PLAN", PolicyHolder: "UJI-PH", Tipe: TypeQR, MataUang: []string{"IDR"}}
	nilai, pesan := RapikanBarisCSV(1, barisLengkapQR(map[string]string{
		"GROSS_PREMIUM": "10,5", "STNC": "03/04/2026", "FACTOR": "1.25", "UW_STATUS": "UJI", "RETROCESSION_VALUATION_BEGIN_DATE": "01/01/2026",
	}), acuan)
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
	// PLAN / POLICY_HOLDER dari polis (10-10-2026): nilai CSV beda tidak ditolak.
	if _, ada := ditemukan["PLAN"]; ada {
		t.Errorf("PLAN CSV masih divalidasi: %q", ditemukan["PLAN"])
	}
	for _, k := range []string{"GROSS_PREMIUM", "DOB", "AGE", "FACTOR", "CERTIFICATE_NO"} {
		if ditemukan[k] == "" {
			t.Errorf("kolom %s tidak dilaporkan (%v)", k, pesan)
		}
	}
	// Nilai CSV berbeda (bahkan > 255 byte) diabaikan; yang disimpan selalu nilai polis.
	nilai, pesan = RapikanBarisCSV(2, barisLengkapQR(map[string]string{"PLAN": strings.Repeat("P", 300), "POLICY_HOLDER": "UJI-LAIN"}), acuan)
	if len(pesan) != 0 || nilai["PLAN"] != "UJI-PLAN" || nilai["POLICY_HOLDER"] != "UJI-PH" {
		t.Errorf("PLAN/POLICY_HOLDER bukan nilai polis: %v %q %q", pesan, nilai["PLAN"], nilai["POLICY_HOLDER"])
	}
	// Kolom tidak ada di baris sama sekali: tetap terisi nilai polis.
	b := barisLengkapQR(nil)
	delete(b, "PLAN")
	delete(b, "POLICY_HOLDER")
	if nilai, pesan = RapikanBarisCSV(2, b, acuan); len(pesan) != 0 || nilai["PLAN"] != "UJI-PLAN" || nilai["POLICY_HOLDER"] != "UJI-PH" {
		t.Errorf("tanpa kolom PLAN/POLICY_HOLDER: %v %v", pesan, nilai)
	}
}

// TestNilaiPolisCSV - PLAN / POLICY_HOLDER dari kepala kasus, fallback peserta Old, keduanya kosong ditolak
// (keputusan work owner 10-10-2026, pola PremiumList Life DataPolisPeserta).
func TestNilaiPolisCSV(t *testing.T) {
	lama := AcuanCSV{Plan: "UJI-PLAN-LAMA", PolicyHolder: "UJI-PH-LAMA"}
	if a, err := NilaiPolisCSV(lama, " GROUP TERM LIFE ", "UJI-PEMEGANG"); err != nil || a.Plan != "GROUP TERM LIFE" || a.PolicyHolder != "UJI-PEMEGANG" {
		t.Errorf("kepala: %+v %v", a, err)
	}
	if a, err := NilaiPolisCSV(lama, "", " "); err != nil || a.Plan != "UJI-PLAN-LAMA" || a.PolicyHolder != "UJI-PH-LAMA" {
		t.Errorf("fallback peserta Old: %+v %v", a, err)
	}
	if a, err := NilaiPolisCSV(AcuanCSV{PolicyHolder: "UJI-PH-LAMA"}, "UJI-COB", ""); err != nil || a.Plan != "UJI-COB" || a.PolicyHolder != "UJI-PH-LAMA" {
		t.Errorf("per kolom: %+v %v", a, err)
	}
	_, err := NilaiPolisCSV(AcuanCSV{}, "", "")
	if !errors.Is(err, ErrDataPolisCSVKosong) || !strings.Contains(err.Error(), "POLICY_HOLDER and PLAN are empty on the policy") {
		t.Errorf("keduanya kosong: %v", err)
	}
	if _, err := NilaiPolisCSV(AcuanCSV{PolicyHolder: "UJI-PH"}, "", ""); !errors.Is(err, ErrDataPolisCSVKosong) || !strings.Contains(err.Error(), "PLAN is empty") {
		t.Errorf("PLAN kosong: %v", err)
	}
}

func TestJudulCSV(t *testing.T) {
	wajib := KolomWajibCSV(TypeQR)
	masuk := append([]string{string(rune(0xFEFF)) + "plan", " Policy_Holder ", "OVRR_COMM", "STATUS"}, wajib...)
	j, asing, err := JudulCSV(masuk, TypeQR)
	if err != nil || strings.Join(j[:4], ",") != "PLAN,POLICY_HOLDER,OVRR_COMM,STATUS" {
		t.Fatalf("%v %v", j, err)
	}
	// PLAN / POLICY_HOLDER tidak wajib di judul (10-10-2026): judul wajib saja diterima.
	if _, _, err := JudulCSV(wajib, TypeQR); err != nil || slices.Contains(wajib, "PLAN") || slices.Contains(wajib, "POLICY_HOLDER") {
		t.Errorf("judul tanpa PLAN/POLICY_HOLDER: %v %v", wajib, err)
	}
	// Kepala unduhan `Generate Data Detail` b276 memuat dua properti di luar 4.1: diabaikan, dilaporkan.
	if strings.Join(asing, ",") != "OVRR_COMM,STATUS" {
		t.Errorf("diabaikan %v", asing)
	}
	for _, rusak := range [][]string{
		{"PLAN"},                          // POLICY_HOLDER hilang - acuan 4.1 tak dapat diperiksa
		{"PLAN", "POLICY_HOLDER", "PLAN"}, // ganda
	} {
		if _, _, err := JudulCSV(rusak, TypeQR); err == nil {
			t.Errorf("%v diterima", rusak)
		}
	}
}

// TestValidasiWajibCSV - kolom wajib per Type, usia, mata uang, sertifikat (keputusan work owner 10-10-2026).
func TestValidasiWajibCSV(t *testing.T) {
	for _, tipe := range []string{TypeQR, TypeQP, TypeTP, TypeTR, "XX"} {
		acuan := AcuanCSV{Plan: "UJI-PLAN", PolicyHolder: "UJI-PH", Tipe: tipe, MataUang: []string{"IDR"}, Sertifikat: map[string]bool{"UJI-D1": true}}
		lengkap := map[string]string{}
		for _, k := range KolomWajibCSV(tipe) {
			lengkap[k] = "0"
		}
		for k, v := range map[string]string{"NAME_OF_INSURED": "UJI-NAMA", "CERTIFICATE_NO": "UJI-N1", "PLAN": "UJI-PLAN",
			"POLICY_HOLDER": "UJI-PH", "SEX": "F", "DOB": "02/01/1990", "CURRENCY": "IDR", "BEGIN_DATE": "01/01/2026",
			"EXPIRED_DATE": "31/12/2026", "ENTRY_AGE": ""} {
			lengkap[k] = v
		}
		// Angka 0 lolos; ENTRY_AGE kosong lolos karena AGE terisi.
		if _, p := RapikanBarisCSV(2, lengkap, acuan); len(p) != 0 {
			t.Fatalf("%s: baris lengkap ditolak %v", tipe, p)
		}
		for _, k := range append(append([]string{}, KolomWajibUmumCSV...), KolomWajibTipeCSV(tipe)...) {
			b := salin(lengkap)
			b[k] = "  "
			if !adaPesan(pesanSaja(RapikanBarisCSV(2, b, acuan)), k, fmt.Sprintf(PesanCSVWajib, k)) {
				t.Errorf("%s: %s kosong tidak ditolak", tipe, k)
			}
		}
		// NET_PREMIUM wajib untuk SEMUA Type (keputusan work owner 10-10-2026): kosong ditolak, 0 lolos.
		b := salin(lengkap)
		b["NET_PREMIUM"] = ""
		if !adaPesan(pesanSaja(RapikanBarisCSV(2, b, acuan)), "NET_PREMIUM", "NET_PREMIUM is empty. Please fill it in.") {
			t.Errorf("%s: NET_PREMIUM kosong lolos", tipe)
		}
		b["NET_PREMIUM"] = "0"
		if _, p := RapikanBarisCSV(2, b, acuan); len(p) != 0 {
			t.Errorf("%s: NET_PREMIUM 0 ditolak %v", tipe, p)
		}
		if !slices.Contains(KolomWajibCSV(tipe), "NET_PREMIUM") || !slices.Contains(JudulTemplatCSV(tipe)[:len(KolomWajibCSV(tipe))], "NET_PREMIUM") {
			t.Errorf("%s: NET_PREMIUM tidak di kolom wajib / depan templat", tipe)
		}
		b = salin(lengkap)
		b["AGE"], b["ENTRY_AGE"] = "", ""
		if !adaPesan(pesanSaja(RapikanBarisCSV(2, b, acuan)), "AGE", PesanCSVUsiaWajib) {
			t.Errorf("%s: AGE dan ENTRY_AGE kosong lolos", tipe)
		}
		b["ENTRY_AGE"] = "30"
		if _, p := RapikanBarisCSV(2, b, acuan); len(p) != 0 {
			t.Errorf("%s: ENTRY_AGE saja ditolak %v", tipe, p)
		}
	}
	acuan := AcuanCSV{Plan: "UJI-PLAN", PolicyHolder: "UJI-PH", Tipe: TypeQR, MataUang: []string{"IDR"}, Sertifikat: map[string]bool{"UJI-D1": true}}
	if !adaPesan(pesanSaja(RapikanBarisCSV(2, barisLengkapQR(map[string]string{"CURRENCY": "EUR"}), acuan)), "CURRENCY", PesanCSVMataUang) {
		t.Error("EUR lolos")
	}
	_, p := RapikanBarisCSV(2, barisLengkapQR(map[string]string{"CURRENCY": "USD"}), acuan)
	if len(p) != 1 || p[0].Pesan != "CURRENCY USD is not used in this policy (use IDR)." {
		t.Errorf("USD bukan mata uang polis: %v", p)
	}
	// Sertifikat peserta Old kasus: tidak peka huruf, spasi dibuang.
	if !adaPesan(pesanSaja(RapikanBarisCSV(2, barisLengkapQR(map[string]string{"CERTIFICATE_NO": " uji-d1 "}), acuan)), "CERTIFICATE_NO", `Certificate No "uji-d1" is already registered in this policy.`) {
		t.Error("sertifikat yang sudah ada lolos")
	}
	// Judul: SELURUH kolom wajib yang kurang disebut sekaligus.
	_, _, err := JudulCSV([]string{"PLAN", "POLICY_HOLDER"}, TypeTR)
	if err == nil || !strings.Contains(err.Error(), "Your file is missing these columns: NAME_OF_INSURED, CERTIFICATE_NO, SEX, DOB, CURRENCY, BEGIN_DATE, EXPIRED_DATE, SUM_INSURED, NET_PREMIUM, AGE or ENTRY_AGE, GROSS_PREMIUM_REFUND_RETRO") {
		t.Errorf("judul kurang: %v", err)
	}
	if TipeCSV("XX") != TypeQR || strings.Join(DaftarKolomWajibCSV(TypeTP), ",") !=
		strings.Join(KolomWajibUmumCSV, ",")+",AGE or ENTRY_AGE,GROSS_PREMIUM_RETRO,DISCOUNT_PREMIUM_RETRO,RI_ADMIN_FEE_RETRO,BROKERAGE_FEE_RETRO" {
		t.Errorf("daftar wajib TP %v", DaftarKolomWajibCSV(TypeTP))
	}
}

// TestTemplatCSV - judul 72 properti dengan kolom wajib di depan; baris contoh lolos bolak-balik.
func TestTemplatCSV(t *testing.T) {
	for _, tipe := range []string{TypeQR, TypeQP, TypeTP, TypeTR} {
		j := JudulTemplatCSV(tipe)
		wajib := KolomWajibCSV(tipe)
		if len(j) != len(PemetaanCSV) || strings.Join(j[:len(wajib)], ",") != strings.Join(wajib, ",") {
			t.Fatalf("%s: judul %v", tipe, j)
		}
		lihat := map[string]bool{}
		for _, k := range j {
			if lihat[k] {
				t.Fatalf("%s: judul ganda %s", tipe, k)
			}
			lihat[k] = true
		}
	}
	j := JudulTemplatCSV(TypeQR)
	nilai := map[string]string{"NAME_OF_INSURED": "UJI-RAHASIA", "CERTIFICATE_NO": "UJI-D1", "PLAN": "UJI-PLAN", "POLICY_HOLDER": "UJI-PH",
		"SEX": "M", "DOB": "02/01/1990", "CURRENCY": "IDR", "BEGIN_DATE": "01/01/2026", "EXPIRED_DATE": "31/12/2026",
		"SUM_INSURED": "1000.123456789", "AGE": "36", "GROSS_PREMIUM": ".5", "STNC": "2026-04-03", "WPC": "bukan", "FACTOR": "1.5"}
	b := BarisContohCSV(j, nilai)
	baris := map[string]string{}
	for i, k := range j {
		baris[k] = b[i]
	}
	if baris["CERTIFICATE_NO"] != "" || baris["NAME_OF_INSURED"] != "" || baris["GROSS_PREMIUM"] != "0.5" ||
		baris["SUM_INSURED"] != "1000.12345679" || baris["STNC"] != "03/04/2026" || baris["WPC"] != "" || baris["DEDUCTION"] != "0" ||
		baris["LAPSE_DATE"] != "" || baris["FACTOR"] != "1.5" {
		t.Errorf("baris contoh %v", baris)
	}
	baris["CERTIFICATE_NO"], baris["NAME_OF_INSURED"] = "UJI-N1", "UJI-NAMA"
	acuan := AcuanCSV{Plan: "UJI-PLAN", PolicyHolder: "UJI-PH", Tipe: TypeQR, MataUang: []string{"IDR"}, Sertifikat: map[string]bool{"UJI-D1": true}}
	if _, p := RapikanBarisCSV(2, baris, acuan); len(p) != 0 {
		t.Errorf("baris contoh tidak lolos: %v", p)
	}
}

func salin(m map[string]string) map[string]string {
	h := map[string]string{}
	for k, v := range m {
		h[k] = v
	}
	return h
}

func pesanSaja(_ map[string]string, p []PesanCSV) []PesanCSV { return p }

func adaPesan(p []PesanCSV, kolom, pesan string) bool {
	for _, x := range p {
		if x.Kolom == kolom && x.Pesan == pesan {
			return true
		}
	}
	return false
}

// TestTaxProfCommClaimTidakWajib - keputusan work owner 10-10-2026: TAX, PROF_COMM, CLAIM tidak wajib (QR dan QP):
// tanpa kolom di judul dan tanpa nilai lolos; bila diisi tetap divalidasi; kosong tersimpan 0.
func TestTaxProfCommClaimTidakWajib(t *testing.T) {
	opsional := []string{"TAX", "PROF_COMM", "CLAIM"}
	for _, tipe := range []string{TypeQR, TypeQP} {
		acuan := AcuanCSV{Plan: "UJI-PLAN", PolicyHolder: "UJI-PH", Tipe: tipe, MataUang: []string{"IDR"}}
		baris := map[string]string{"NAME_OF_INSURED": "UJI-NAMA", "CERTIFICATE_NO": "UJI-N1", "PLAN": "UJI-PLAN",
			"POLICY_HOLDER": "UJI-PH", "SEX": "M", "DOB": "02/01/1990", "CURRENCY": "IDR", "BEGIN_DATE": "01/01/2026",
			"EXPIRED_DATE": "31/12/2026", "AGE": "30"}
		for _, k := range append(append([]string{}, KolomWajibUmumCSV...), KolomWajibTipeCSV(tipe)...) {
			if _, ada := baris[k]; !ada {
				baris[k] = "0"
			}
		}
		for _, k := range opsional {
			if slices.Contains(KolomWajibCSV(tipe), k) || slices.Contains(DaftarKolomWajibCSV(tipe), k) {
				t.Errorf("%s: %s masih wajib", tipe, k)
			}
		}
		// Tanpa kolom sama sekali.
		nilai, p := RapikanBarisCSV(2, baris, acuan)
		if len(p) != 0 || nilai["TAX"] != "0" || nilai["PROF_COMM"] != "0" || nilai["CLAIM"] != "0" {
			t.Errorf("%s tanpa kolom: %v %v", tipe, p, nilai)
		}
		// Kolom ada tetapi kosong.
		for _, k := range opsional {
			baris[k] = " "
		}
		if _, p := RapikanBarisCSV(2, baris, acuan); len(p) != 0 {
			t.Errorf("%s kosong: %v", tipe, p)
		}
		// Diisi bentuk rusak: tetap divalidasi.
		baris["TAX"] = "1.234.5"
		if _, p := RapikanBarisCSV(2, baris, acuan); len(p) != 1 || p[0].Kolom != "TAX" {
			t.Errorf("%s TAX rusak: %v", tipe, p)
		}
		// Judul tanpa ketiganya diterima.
		j := KolomWajibCSV(tipe)
		if _, _, err := JudulCSV(j, tipe); err != nil {
			t.Errorf("%s judul tanpa TAX/PROF_COMM/CLAIM: %v", tipe, err)
		}
	}
}

// TestKalimatCSVSederhana - kalimat penolakan untuk user awam (keputusan work owner 10-10-2026): nama kolom disebut,
// contoh format, tanpa awalan teknis.
func TestKalimatCSVSederhana(t *testing.T) {
	acuan := AcuanCSV{Plan: "UJI-PLAN", PolicyHolder: "UJI-PH", Tipe: TypeQR, MataUang: []string{"IDR", "USD"}}
	_, p := RapikanBarisCSV(2, barisLengkapQR(map[string]string{
		"DOB": "1990-11-22", "GROSS_PREMIUM": "1.500.000,50", "DEDUCTION": "abc", "RI_ADMIN_FEE": "0.123456789",
		"AGE": "3x", "FACTOR": "0,75", "NAME_OF_INSURED": strings.Repeat("N", 256), "SEX": " ", "CURRENCY": "EUR",
	}), acuan)
	mau := map[string]string{
		"DOB":             "DOB must be a date in DD/MM/YYYY format, e.g. 22/11/1990.",
		"GROSS_PREMIUM":   "GROSS_PREMIUM must not use thousands separators, e.g. 1500000.50 or 1500000,50.",
		"DEDUCTION":       "DEDUCTION must be a number, e.g. 1500000.50.",
		"RI_ADMIN_FEE":    "RI_ADMIN_FEE has more than 8 decimal places.",
		"AGE":             "AGE must be a whole number (max 5 digits).",
		"FACTOR":          "FACTOR must be a number with a dot as decimal separator, e.g. 0.75.",
		"NAME_OF_INSURED": "NAME_OF_INSURED is too long (max 255 characters).",
		"SEX":             "SEX is empty. Please fill it in.",
		"CURRENCY":        "CURRENCY must be IDR or USD.",
	}
	dapat := map[string]string{}
	for _, x := range p {
		dapat[x.Kolom] = x.Pesan
	}
	for k, v := range mau {
		if dapat[k] != v {
			t.Errorf("%s: %q, mau %q", k, dapat[k], v)
		}
	}
	b := barisLengkapQR(nil)
	b["AGE"], b["ENTRY_AGE"] = "", ""
	if !adaPesan(pesanSaja(RapikanBarisCSV(2, b, acuan)), "AGE", "AGE or ENTRY_AGE must be filled in.") {
		t.Error("kalimat AGE/ENTRY_AGE")
	}
	// Mata uang polis lebih dari satu: "use IDR or USD".
	_, p = RapikanBarisCSV(2, barisLengkapQR(map[string]string{"CURRENCY": "USD"}), AcuanCSV{Plan: "UJI-PLAN", PolicyHolder: "UJI-PH", Tipe: TypeQR, MataUang: []string{"IDR"}})
	if len(p) != 1 || p[0].Pesan != "CURRENCY USD is not used in this policy (use IDR)." {
		t.Errorf("mata uang polis: %v", p)
	}
	// Judul: kalimat layar tanpa awalan teknis; errors.Is tetap ErrJudulCSV.
	for masuk, kalimat := range map[string]string{
		"kurang": "Your file is missing these columns:",
		"ganda":  "Column PLAN appears more than once in the header.",
		"kosong": "The header has an empty column name (column 2).",
	} {
		var judul []string
		switch masuk {
		case "kurang":
			judul = []string{"PLAN"}
		case "ganda":
			judul = append(KolomWajibCSV(TypeQR), "PLAN", "plan")
		case "kosong":
			judul = append([]string{"", ""}, KolomWajibCSV(TypeQR)...)
		}
		_, _, err := JudulCSV(judul, TypeQR)
		var g PesanGalatCSV
		if !errors.Is(err, ErrJudulCSV) || !errors.As(err, &g) || !strings.HasPrefix(g.PesanLayar(), kalimat) ||
			strings.Contains(g.PesanLayar(), "models:") || strings.Contains(g.PesanLayar(), "invalid CSV header") {
			t.Errorf("%s: %v", masuk, err)
		}
	}
}
