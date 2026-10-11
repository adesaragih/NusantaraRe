package models_test

import (
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimnonprop/backend/models"
)

var saatUji = time.Date(2026, 10, 8, 10, 0, 0, 0, models.Jakarta)

// kasusUji - dua penyetuju UJI-, tingkat berjalan `count`.
func kasusUji(count int, putusan ...string) models.Kasus {
	k := models.Kasus{ID: "KMTNP-UJI001", KlaimID: "CLMNP-UJI001", AdjustmentID: "UJI-ADJ-1", Loop: 2, Count: count,
		UsulTutup: models.UsulTidak, UsulCadang: models.UsulTidak, Tahap: models.TahapKomite,
		Tangga: []models.Anggota{
			{ID: "L1", Urut: 1, OperatorID: "UJI-K1", Jabatan: "UJI-JABATAN-1", Keputusan: models.KeputusanMenunggu},
			{ID: "L2", Urut: 2, OperatorID: "UJI-K2", Jabatan: "UJI-JABATAN-2", Keputusan: models.KeputusanMenunggu},
		}}
	for i, p := range putusan {
		k.Tangga[i].Keputusan = p
	}
	return k
}

func klaimUji(acceptedNo string) kontrak.KlaimTreaty {
	return kontrak.KlaimTreaty{Adjustment: 1, Daftar: map[string][]map[string]string{
		"ClaimData.AdjustmentList": {{"ID": "UJI-ADJ-1", "AcceptedNo": acceptedNo}},
	}}
}

func TestGiliranKomiteRouter(t *testing.T) {
	k := kasusUji(1)
	if a, ada := k.Giliran(); !ada || a.OperatorID != "UJI-K1" {
		t.Fatalf("giliran tingkat 1: %+v %v", a, ada)
	}
	if !k.Pemegang("UJI-K1", nil) || k.Pemegang("UJI-K2", nil) {
		t.Fatal("pemegang tingkat 1 = UJI-K1 saja")
	}
	if !k.Pemegang("UJI-LAIN", []string{"UJI-K1"}) {
		t.Fatal("KomiteID workbasket: pemegang workbasket aktif ikut memegang")
	}
	k = kasusUji(2, models.KeputusanSetuju)
	if !k.Pemegang("UJI-K2", nil) || k.Pemegang("UJI-K1", nil) {
		t.Fatal("sesudah UJI-K1 setuju, giliran UJI-K2")
	}
	k = kasusUji(3, models.KeputusanSetuju, models.KeputusanSetuju)
	if _, ada := k.Giliran(); ada || k.Pemegang("UJI-K2", nil) {
		t.Fatal("tanpa baris menunggu tidak ada pemegang")
	}
	k = kasusUji(1)
	k.StatusWork = models.StatusSelesai
	if k.Pemegang("UJI-K1", nil) {
		t.Fatal("kasus tertutup tidak dipegang siapa pun")
	}
}

func TestPeriksaIsianShowTransfer(t *testing.T) {
	p := models.PeriksaIsian(kasusUji(1), models.Keputusan{})
	if len(p) != 2 || !strings.HasPrefix(p[0], models.LabelAcceptStatus) || !strings.HasPrefix(p[1], "Note:") {
		t.Fatalf("AcceptStatus dan Note wajib: %v", p)
	}
	k := kasusUji(1)
	k.Loop = 1
	k.Tangga = k.Tangga[:1]
	p = models.PeriksaIsian(k, models.Keputusan{AcceptStatus: "1", Comment: "UJI", IsSubjectivity: true})
	if len(p) != 1 || !strings.HasPrefix(p[0], models.LabelSubjectivityNote) {
		t.Fatalf("Subjectivity Note wajib bila Subjectivity: %v", p)
	}
	if p := models.PeriksaIsian(k, models.Keputusan{AcceptStatus: "1", Comment: "UJI", IsSubjectivity: true,
		SubjectivityNote: "UJI bebas"}); len(p) != 1 || p[0] != "Subjectivity Note: Value is not in the list" {
		t.Fatalf("Subjectivity Note di luar daftar pilihan: %v", p)
	}
	// tingkat 2: isian tingkat 1 nonaktif - diabaikan, tidak memicu pesan
	if p := models.PeriksaIsian(kasusUji(2, "1"), models.Keputusan{AcceptStatus: "1", Comment: "UJI",
		IsSubjectivity: true}); len(p) != 0 {
		t.Fatalf("tingkat 2: %v", p)
	}
}

func TestRencanaSetujuBukanTingkatAkhir(t *testing.T) {
	r := models.Rencanakan(kasusUji(1), klaimUji(""), models.Keputusan{AcceptStatus: "1", Comment: "UJI-OK",
		UsulTutup: true}, "UJI-K1", "UJI Nama Satu", saatUji)
	if len(r.Tangga) != 1 || r.Tangga[0].ID != "L1" || r.Tangga[0].Keputusan != "1" || r.Tangga[0].Komentar != "UJI-OK" ||
		r.Tangga[0].Pemutus != "UJI-K1" {
		t.Fatalf("S6 tangga: %+v", r.Tangga)
	}
	if len(r.Klaim.Riwayat) != 1 || r.Klaim.Riwayat[0].Teks != "Accepted by UJI Nama Satu" {
		t.Fatalf("S8 / S10 riwayat: %+v", r.Klaim.Riwayat)
	}
	if len(r.Klaim.Header) != 0 || r.Klaim.Adjustment != nil || r.UsulTutup != "1" || r.UsulCadang != "0" {
		t.Fatalf("tingkat bukan akhir tidak menyentuh klaim: %+v / %+v, usul %s/%s", r.Klaim.Header, r.Klaim.Adjustment,
			r.UsulTutup, r.UsulCadang)
	}
	if r.Count != 2 || r.Selesai || r.TingkatAkhirSetuju || r.TerbitkanNomor || r.SimpanOS || r.Kasir ||
		r.SimpanOSSubjectivity || r.Posisi != "UJI-K2" {
		t.Fatalf("tingkat 1 dari 2: %+v", r)
	}
	if r.Email != models.EmailPenyetujuBerikut || r.StatusRiwayat != "ACCEPT" {
		t.Fatalf("email %q, riwayat %q", r.Email, r.StatusRiwayat)
	}
}

func TestRencanaSetujuTingkatAkhir(t *testing.T) {
	k := kasusUji(2, models.KeputusanSetuju)
	k.UsulTutup = models.UsulYa // diisi tingkat 1
	r := models.Rencanakan(k, klaimUji(""), models.Keputusan{AcceptStatus: "1", Comment: "UJI-AKHIR"}, "UJI-K2",
		"UJI Nama Dua", saatUji)
	h, a := r.Klaim.Header, r.Klaim.Adjustment
	if h["IsCloseFile"] != "1" || h["CNPStatusCase"] != "CLAIM ACCEPTED" || len(h) != 2 {
		t.Fatalf("S14.12 / S14.17 header: %+v", h)
	}
	if a["IsSubjectivity"] != "false" || a["SubjectivityNote"] != "" || len(a) != 2 {
		t.Fatalf("S18 akseptasi: %+v", a)
	}
	if !r.TingkatAkhirSetuju || !r.TerbitkanNomor || !r.SimpanOS || !r.Konversi || !r.Kasir || !r.SimpanOSSubjectivity {
		t.Fatalf("langkah tingkat akhir: %+v", r)
	}
	if r.Count != 3 || !r.Selesai || r.Posisi != "" || r.Email != models.EmailPembuatSetuju {
		t.Fatalf("selesai: count %d selesai %v posisi %q email %q", r.Count, r.Selesai, r.Posisi, r.Email)
	}
	r.AdjustmentDiterima("UJI-NOMOR", saatUji)
	if a := r.Klaim.Adjustment; a["AcceptedNo"] != "UJI-NOMOR" || a["AcceptanceStatus"] != "1" ||
		a["AcceptedDate"] != "2026-10-08 10:00:00" {
		t.Fatalf("S14.13: %+v", a)
	}
	// nomor sudah ada -> S14.9-S14.11 dilewati, OS tetap
	if r := models.Rencanakan(k, klaimUji("UJI-SUDAH"), models.Keputusan{AcceptStatus: "1", Comment: "x"}, "UJI-K2", "x",
		saatUji); r.TerbitkanNomor || !r.SimpanOS {
		t.Fatalf("AcceptedNo terisi: nomor %v OS %v", r.TerbitkanNomor, r.SimpanOS)
	}
}

func TestRencanaTolakKeExit(t *testing.T) {
	r := models.Rencanakan(kasusUji(1), klaimUji(""), models.Keputusan{AcceptStatus: "2", Comment: "UJI-TOLAK"}, "UJI-K1",
		"UJI Nama Satu", saatUji)
	if len(r.Tangga) != 2 || r.Tangga[0].ID != "L1" || r.Tangga[0].Keputusan != "2" || r.Tangga[1].ID != "L2" ||
		r.Tangga[1].Keputusan != "2" || r.Tangga[1].Komentar != "UJI-TOLAK" || !r.Tangga[1].IsiKomentar ||
		r.Tangga[1].Pemutus != "" {
		t.Fatalf("S6 + S11.3: %+v", r.Tangga)
	}
	if r.Klaim.Adjustment["AcceptanceStatus"] != "2" || r.Klaim.Header["CNPStatusCase"] != "CLAIM REJECTED" ||
		len(r.Klaim.Adjustment) != 1 || len(r.Klaim.Header) != 1 {
		t.Fatalf("S11.4 / S20: %+v %+v", r.Klaim.Adjustment, r.Klaim.Header)
	}
	if r.Klaim.Riwayat[0].Teks != "Rejected by UJI Nama Satu" || r.StatusRiwayat != "REJECT" {
		t.Fatalf("riwayat: %+v / %q", r.Klaim.Riwayat, r.StatusRiwayat)
	}
	if r.Count != 3 || !r.Selesai || r.Email != "" || r.TerbitkanNomor || r.SimpanOS || r.Kasir ||
		r.SimpanOSSubjectivity {
		t.Fatalf("S12 tolak melompat ke EXIT: %+v", r)
	}
}

func TestRencanaSubjectivityTingkatSatu(t *testing.T) {
	k := kasusUji(1)
	k.Loop = 1
	k.Tangga = k.Tangga[:1]
	r := models.Rencanakan(k, klaimUji(""), models.Keputusan{AcceptStatus: "1", Comment: "UJI", IsSubjectivity: true,
		SubjectivityNote: "3"}, "UJI-K1", "UJI", saatUji)
	if !r.Subjectivity || r.TerbitkanNomor || r.SimpanOS || r.Kasir || r.Konversi {
		t.Fatalf("subjectivity melewati S14: %+v", r)
	}
	a := r.Klaim.Adjustment
	if !r.SimpanOSSubjectivity || a["IsKomite"] != "0" || a["IsSubjectivity"] != "true" || a["SubjectivityNote"] != "3" ||
		len(r.Klaim.Header) != 0 {
		t.Fatalf("S16 / S17 / S18: %+v %+v", a, r.Klaim.Header)
	}
	if _, ada := r.Klaim.Header["CNPStatusCase"]; ada {
		t.Fatal("S14.17 dilewati: CNPStatusCase tidak disentuh")
	}
}

// Isian tingkat 1 tersimpan di header kasus komite, dibaca tingkat akhir.
func TestRencanaSubjectivityTanggaDuaTingkat(t *testing.T) {
	r1 := models.Rencanakan(kasusUji(1), klaimUji(""), models.Keputusan{AcceptStatus: "1", Comment: "UJI",
		IsSubjectivity: true, SubjectivityNote: "3"}, "UJI-K1", "UJI", saatUji)
	if r1.SubjectivitySimpan != "1" || r1.SubjectivityNoteSimpan != "3" || r1.TingkatAkhirSetuju ||
		r1.SimpanOSSubjectivity || r1.Klaim.Adjustment["IsKomite"] != "0" {
		t.Fatalf("tingkat 1 menyimpan isian, S16 berjalan: %+v", r1)
	}
	k := kasusUji(2, "1")
	k.Subjectivity, k.SubjectivityNote = "1", "3"
	// tingkat 2: isian Subjectivity nonaktif - nilai yang dikirim diabaikan, yang tersimpan dipakai.
	r2 := models.Rencanakan(k, klaimUji(""), models.Keputusan{AcceptStatus: "1", Comment: "UJI-2"}, "UJI-K2", "UJI",
		saatUji)
	if !r2.Subjectivity || r2.TerbitkanNomor || r2.SimpanOS || r2.Kasir || !r2.SimpanOSSubjectivity ||
		r2.SubjectivitySimpan != "1" {
		t.Fatalf("tingkat akhir memakai subjectivity tersimpan: %+v", r2)
	}
	if a := r2.Klaim.Adjustment; a["IsKomite"] != "0" || a["SubjectivityNote"] != "3" {
		t.Fatalf("S16 / S18 dari isian tersimpan: %+v", a)
	}
}

// Kirim ulang baris subjectivity - S6 dilewati, S7 menulis ComiteeClaim(<LAST>).
func TestRencanaKirimUlangSubjectivityS7(t *testing.T) {
	kl := klaimUji("")
	kl.Daftar["ClaimData.AdjustmentList"][0]["IsSubjectivity"] = "true"
	r := models.Rencanakan(kasusUji(1), kl, models.Keputusan{AcceptStatus: "1", Comment: "UJI"}, "UJI-K1", "UJI",
		saatUji)
	if len(r.Tangga) != 1 || r.Tangga[0].ID != "L2" {
		t.Fatalf("S7 menulis baris terakhir tangga: %+v", r.Tangga)
	}
}

func TestRakitNomorAkseptasi(t *testing.T) {
	if got := models.RakitNomorAkseptasi("UJI"+models.HurufAkseptasi, "123", "10.2026", 7); got != "UJIA123.10.2026.TX00007" {
		t.Fatalf("nomor %q", got)
	}
	if !models.PanjangNoAksepCNP("UJIA123.10.2026.TX00007") || models.PanjangNoAksepCNP("UJIA12.10.2026.TX00007") {
		t.Fatal("HitServiceToKasirKMT_Act 9-10: panjang nomor 23 / 24")
	}
}

// klaimAkhir - satu akseptasi: layer UR + satu layer XoL, Spreading In dua baris mata uang sama + satu lain.
func klaimAkhir() kontrak.KlaimTreaty {
	adj := func(anak string) string { return "ClaimData.AdjustmentList(1)." + anak }
	return kontrak.KlaimTreaty{Adjustment: 1,
		Nilai: map[string]string{"pyID": "CLMNP-UJI001", "ClaimData.NoClaim": "UJI-K-0001",
			"ClaimData.PolicyData.PolicyNo": "UJI-POLIS", "ClaimData.IDMaster": "UJI-M1",
			"ClaimData.CauseOfLoss": "UJI-SEBAB", "ClaimData.CauseOfLossID": "UJI-S1", "TreatyInMaster.RNMShare": "25",
			"ClaimData.InsuredName": "UJI-TERTANGGUNG", "ClaimData.DateOfLoss": "2026-02-07",
			"OfferFacIn.QuotationData.BusinessOldId": "123"},
		Daftar: map[string][]map[string]string{
			"ClaimData.AdjustmentList": {{"ID": "UJI-ADJ-1", "PaymentType": "1", "PayableTo": "UJI-PENERIMA",
				"NoAccount": "12-34 56", "AcceptedNo": "UJIA123.10.2026.TX00001", "AcceptedDate": "2026-10-28 10:00:00",
				"pxCreateOperator": "UJI-ADMIN"}},
			adj(models.AnakXOL): {
				{"TreatyName": "UR", "Currency": "UJA", "TotalClaim": "100"},
				{"TreatyName": "UJI-L1", "TreatyType": "UJI-X1", "Currency": "UJA", "CurrencyID": "UJI-A",
					"ClaimEstimation": "400", "TotalClaim": "400", "ClaimPercentage": "25", "ClaimSpreaded": "100",
					"AdjusterFee": "5", "Salvage": "0", "CNPOthersFee": "0", "CNPReinstatement": "8",
					"CNPReinstatementRNM": "2", "Kurs": "2", "KursIDR": "2", "CNPLimit": "1000", "CNPMDP": "10",
					"CNPPctReinstate": "100"},
			},
			adj(models.AnakSpreadIn): {
				{"TreatyName": "UJI-QS", "Currency": "UJA", "CurrencyID": "UJI-A", "SharePercentage": "60",
					"TotalClaim": "60", "PremiumSpreaded": "1.5", "NoAccount": "98-76", "IDOfBank": "UJI-BANK"},
				{"TreatyName": "UJI-SP", "Currency": "UJA", "CurrencyID": "UJI-A", "SharePercentage": "40",
					"TotalClaim": "40", "PremiumSpreaded": "0.5"},
				{"TreatyName": "UJI-QS", "Currency": "UJB", "CurrencyID": "UJI-B", "SharePercentage": "100",
					"TotalClaim": "7"},
			},
		}}
}

// halaman - teks `@ASM.GetPageJSONString()` dari pasangan dan PageList (LF; "," di depan pasangan berikutnya).
func halaman(bagian ...string) string {
	lf := string(rune(10))
	return "{" + lf + strings.Join(bagian, lf+",") + lf + "}"
}

func daftarJSON(nama string, isi ...string) string {
	lf := string(rune(10))
	return `"` + nama + `":[ ` + lf + strings.Join(isi, lf+",") + lf + "] "
}

func TestSusunOSAkseptasiInsertOSKlaimCNP(t *testing.T) {
	kl := klaimAkhir()
	os, err := models.SusunOSAkseptasi(kl, "CLMNP-UJI001", "UJI-NOMOR", saatUji)
	if err != nil {
		t.Fatal(err)
	}
	mu := halaman(`"AdjusterFeeRNM":"5"`, `"AdjusterFeeValue":"20"`, `"CNPLimit":"1000"`, `"CNPMindep":"10"`,
		`"CNPOthersFee":"0"`, `"CNPOthersFeeRNM":"0"`, `"CNPPctReinstate":"100"`, `"CNPReinstatement":"8"`,
		`"CNPReinstatementRNM":"2"`, `"Currency":"UJA"`, `"CurrencyID":"UJI-A"`, `"GrossAdjustment":"400"`,
		`"GrossValue":"100"`, `"KursIDR":"2"`, `"PersenRNM":"25"`, `"pxObjClass":"ASM-FW-GCNMFW-Data-Adjustment"`,
		`"SalvageRNM":"0"`, `"SalvageValue":"0"`, `"TotalXOL":"400"`, `"TotalXOLGross":"400"`, `"TotalXOLLayerRNM":"100"`,
		`"TotalXOLRNM":"100"`)
	layer := halaman(`"pxObjClass":"ASM-FW-GCNMFW-Data-Adjustment"`, `"XOL":"UJI-L1"`, `"XOLID":"UJI-X1"`,
		daftarJSON("CNPCurrencyList", mu))
	dj := halaman(`"AcceptedNo":"UJI-NOMOR"`, `"CauseOfLoss":"UJI-SEBAB"`, `"CauseOfLossID":"UJI-S1"`,
		`"EstimationDate":"20261008"`, `"NoClaim":"UJI-K-0001"`, `"PersenRNM":"25"`,
		`"pxObjClass":"ASM-FW-GCNMFW-Data-osAkseptasi"`, `"Type":"1"`, daftarJSON("CNPLayerList", layer)) + string(rune(10))
	mau := models.BarisOSAkseptasi{CaseID: "CLMNP-UJI001", NoClaim: "UJI-K-0001", NoPolis: "UJI-POLIS", StsReject: "4",
		MasterID: "UJI-M1", DataJSON: dj}
	if os != mau {
		t.Fatalf("baris OS\n%+v\nmau\n%+v", os, mau)
	}
	sub, err := models.SusunOSSubjectivity(kl, "CLMNP-UJI001", "UJI-NOMOR", true, saatUji)
	if err != nil {
		t.Fatal(err)
	}
	if sub.StsSubjectivity != "1" || sub.StsReject != "4" || sub.DataJSON != dj || sub.CaseID != "CLMNP-UJI001" {
		t.Fatalf("OS subjectivity: %+v", sub)
	}
	// InsertXOLKlaimCNP: layer UR tidak masuk; subjectivity -> nol baris.
	x := models.SusunXOL2(kl, "CLMNP-UJI001", false)
	if len(x) != 1 || x[0] != (models.BarisXOL2{CaseID: "CLMNP-UJI001", GrossAdjustment: "400", CNPReinstatement: "8",
		Currency: "UJA", KursIDR: "2", XOL: "UJI-L1"}) {
		t.Fatalf("CLAIMXOL2: %+v", x)
	}
	if len(models.SusunXOL2(kl, "CLMNP-UJI001", true)) != 0 {
		t.Fatal("InsertXOLKlaimCNP S1: subjectivity keluar")
	}
	if models.StsOSAkseptasi("2") != "1" {
		t.Fatal("InsertOSKlaimCNP S3: PaymentType selain 1 -> 1")
	}
}

func TestJSONHalamanPegaTanpaEscapeHTML(t *testing.T) {
	got := models.JSONHalamanPega(models.HalamanJSON{Nilai: map[string]string{"B": `a/b <c> "d"`, "a": "1", "Kosong": ""}})
	mau := "{" + string(rune(10)) + `"a":"1"` + string(rune(10)) + "," + `"B":"a/b <c> \"d\""` + string(rune(10)) + "}" +
		string(rune(10))
	if got != mau {
		t.Fatalf("JSON %q, mau %q", got, mau)
	}
}

func TestMuatanKasirDanEmail(t *testing.T) {
	kl := klaimAkhir()
	cfg := models.KonfigurasiKasir{CompanyName: "UJI-CO", LjtdID: "UJI-LJ", LdcID: "UJI-LDC", LdcIDSyariah: "UJI-LDCS"}
	m, err := models.SusunMuatanKasir(kl, models.AdjustmentKlaim(kl), "uji@contoh.invalid", false, cfg, "UJI-K2")
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 2 {
		t.Fatalf("satu muatan per mata uang Spreading In: %+v", m)
	}
	a, b := m[0], m[1]
	if a.LkuId != "UJI-A" || a.Nett != "98" || a.AccountNo != "9876" || a.LbgID != "UJI-BANK" || a.NoTrans !=
		"UJIA123.10.2026.TX00001" || a.TglAksep != "28-10-2026" || a.TglBolehBayar != "01-12-2026" || a.LbuId != "123" ||
		a.LdcId != "UJI-LDC" || a.CompanyName != "UJI-CO" || a.UserInput != "UJI-ADMIN" || a.AcceptType != "1" ||
		a.Deductible != "0" || a.Email != "uji@contoh.invalid" {
		t.Fatalf("muatan kasir UJA %+v", a)
	}
	if b.LkuId != "UJI-B" || b.Nett != "7" || b.AccountNo != "123456" || models.RekeningAngka("12-34 56") != "123456" {
		t.Fatalf("muatan kasir UJB (rekening akseptasi angka saja, S7): %+v", b)
	}
	if got := models.SubjekEmail(models.EmailPembuatSetuju, kl, "KMTNP-UJI001"); got !=
		"(Approval) Pengajuan Akseptasi : CLMNP-UJI001/KMTNP-UJI001 UJI-TERTANGGUNG DOL 07 Febuari 2026" {
		t.Fatalf("subjek %q", got)
	}
	if got := models.SubjekEmail(models.EmailPenyetujuBerikut, kl, "X"); !strings.HasPrefix(got, "Pengajuan Akseptasi : ") {
		t.Fatalf("subjek berikut %q", got)
	}
	if !models.AdaSpreadingAdjustment(kl) {
		t.Fatal("SendEmailKlaim_KMT: spreading adjustment ada")
	}
}

func TestTanggalBolehBayarBergulirTahun(t *testing.T) {
	for ymd, mau := range map[string]string{"20261028": "01-12-2026", "20261210": "10-01-2027",
		"20261227": "01-02-2027", "20260115": "15-02-2026"} {
		if got := models.TanggalBolehBayar(ymd); got != mau {
			t.Errorf("TanggalBolehBayar(%s) = %s, mau %s", ymd, got, mau)
		}
	}
}

func TestLabelPromptValues(t *testing.T) {
	for kode, mau := range map[string]string{"0": "Transfer to committee", "1": "Approve", "2": "Reject", "": "", "9": "9"} {
		if l := models.LabelStatusBaris(kode); l != mau {
			t.Errorf("AcceptanceStatus %q: %q, mau %q", kode, l, mau)
		}
	}
	if !models.SubjectivityNoteSah("7") || models.SubjectivityNoteSah("8") || models.SubjectivityNoteSah("") {
		t.Fatal("SubjectivityNote hanya 1..7")
	}
}
