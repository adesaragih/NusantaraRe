package models_test

import (
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimprop/backend/models"
)

var saatUji = time.Date(2026, 10, 8, 10, 0, 0, 0, models.Jakarta)

// kasusUji - dua penyetuju UJI-, tingkat berjalan `count`.
func kasusUji(count int, putusan ...string) models.Kasus {
	k := models.Kasus{ID: "TKMT-UJI001", KlaimID: "CLMP-UJI001", AdjustmentID: "UJI-ADJ-1", Loop: 2, Count: count,
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
	if !k.Pemegang("UJI-K1") || k.Pemegang("UJI-K2") {
		t.Fatal("pemegang tingkat 1 = UJI-K1 saja")
	}
	k = kasusUji(2, models.KeputusanSetuju)
	if !k.Pemegang("UJI-K2") || k.Pemegang("UJI-K1") {
		t.Fatal("sesudah UJI-K1 setuju, giliran UJI-K2")
	}
	k = kasusUji(3, models.KeputusanSetuju, models.KeputusanSetuju)
	if _, ada := k.Giliran(); ada || k.Pemegang("UJI-K2") {
		t.Fatal("tanpa baris menunggu tidak ada pemegang (KomiteRouter tanpa cadangan)")
	}
	k = kasusUji(1)
	k.StatusWork = models.StatusSelesai
	if k.Pemegang("UJI-K1") {
		t.Fatal("kasus tertutup tidak dipegang siapa pun")
	}
}

func TestMasihBerjalanIsKomiteLoop(t *testing.T) {
	cases := []struct {
		status      string
		count, loop int
		mau         bool
	}{{"1", 2, 2, true}, {"1", 3, 2, false}, {"2", 2, 2, false}, {"2", 3, 2, false}}
	for _, c := range cases {
		if got := models.MasihBerjalan(c.status, c.count, c.loop); got != c.mau {
			t.Errorf("IsKomiteLoop(%s, %d, %d) = %v, mau %v", c.status, c.count, c.loop, got, c.mau)
		}
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
		SubjectivityNote: "1"}); len(p) != 0 {
		t.Fatalf("tangga satu tingkat: %v", p)
	}
	p = models.PeriksaIsian(kasusUji(1), models.Keputusan{AcceptStatus: "1", Comment: "UJI", IsSubjectivity: true,
		SubjectivityNote: "1"})
	if len(p) != 1 || p[0] != models.PesanSubjectivityBertingkat {
		t.Fatalf("OQ-KCP-01 tangga dua tingkat: %v", p)
	}
	// tingkat 2: isian tingkat 1 nonaktif - diabaikan, tidak memicu pesan
	if p := models.PeriksaIsian(kasusUji(2, "1"), models.Keputusan{AcceptStatus: "1", Comment: "UJI",
		IsSubjectivity: true}); len(p) != 0 {
		t.Fatalf("tingkat 2: %v", p)
	}
}

func TestRencanaSetujuBukanTingkatAkhir(t *testing.T) {
	r := models.Rencanakan(kasusUji(1), klaimUji(""), models.Keputusan{AcceptStatus: "1", Comment: "UJI-OK",
		UsulTutup: true}, "UJI-K1", saatUji)
	if len(r.Tangga) != 1 || r.Tangga[0].ID != "L1" || r.Tangga[0].Keputusan != "1" || r.Tangga[0].Komentar != "UJI-OK" {
		t.Fatalf("S6 tangga: %+v", r.Tangga)
	}
	if len(r.Klaim.Riwayat) != 1 || r.Klaim.Riwayat[0].Teks != "Accepted by UJI-JABATAN-1" ||
		r.Klaim.Riwayat[0].Tingkat != "UJI-JABATAN-1" || r.Klaim.Riwayat[0].Pelaku != "UJI-K1" {
		t.Fatalf("S8-S10 riwayat: %+v", r.Klaim.Riwayat)
	}
	if r.Klaim.Header["ClaimData.IsCloseFile"] != "true" || r.Klaim.Header["ClaimData.IsReservedClaim"] != "false" ||
		r.UsulTutup != "1" || r.UsulCadang != "0" {
		t.Fatalf("S11 usul: header %+v, kepala %s/%s", r.Klaim.Header, r.UsulTutup, r.UsulCadang)
	}
	if r.Count != 2 || r.Selesai || r.TingkatAkhirSetuju || r.TerbitkanNomor || r.Kasir || r.JSONKlaim {
		t.Fatalf("tingkat 1 dari 2: %+v", r)
	}
	if _, ada := r.Klaim.Header["IsAnyAcceptation"]; ada || r.Klaim.Adjustment != nil {
		t.Fatalf("tingkat bukan akhir tidak menyentuh akseptasi: %+v / %+v", r.Klaim.Header, r.Klaim.Adjustment)
	}
	if r.Email != models.EmailPenyetujuBerikut || r.StatusRiwayat != "ACCEPT" {
		t.Fatalf("email %q, riwayat %q", r.Email, r.StatusRiwayat)
	}
}

func TestRencanaSetujuTingkatAkhir(t *testing.T) {
	k := kasusUji(2, models.KeputusanSetuju)
	k.UsulCadang = models.UsulYa // diisi tingkat 1
	r := models.Rencanakan(k, klaimUji(""), models.Keputusan{AcceptStatus: "1", Comment: "UJI-AKHIR", UsulTutup: true},
		"UJI-K2", saatUji)
	if r.UsulTutup != "0" || r.UsulCadang != "1" || r.Klaim.Header["ClaimData.IsReservedClaim"] != "true" {
		t.Fatalf("isian tingkat 1 dipertahankan di tingkat 2: %s/%s %+v", r.UsulTutup, r.UsulCadang, r.Klaim.Header)
	}
	h, a := r.Klaim.Header, r.Klaim.Adjustment
	if h["IsAnyAcceptation"] != "1" || h["AktifButton"] != "0" || h["IsOutstanding"] != "1" || h["IsCFS"] != "" ||
		h["ClaimData.IsSubjectivity"] != "false" {
		t.Fatalf("header tingkat akhir: %+v", h)
	}
	if a["IsApproved"] != "1" || a["IsSubjectivity"] != "false" || a["SubjectivityNote"] != "" || a["Notes"] != "UJI-AKHIR" {
		t.Fatalf("adjustment tingkat akhir: %+v", a)
	}
	if !r.TingkatAkhirSetuju || !r.TerbitkanNomor || !r.SimpanAkseptasi || !r.SimpanRetro || !r.JSONKlaim ||
		!r.Konversi || !r.LogAkseptasi || !r.Kasir {
		t.Fatalf("langkah tingkat akhir: %+v", r)
	}
	if r.Count != 3 || !r.Selesai || r.Email != models.EmailPembuatSetuju {
		t.Fatalf("selesai: count %d selesai %v email %q", r.Count, r.Selesai, r.Email)
	}
	r.AdjustmentDiterima("UJI-NOMOR", saatUji)
	if a := r.Klaim.Adjustment; a["AcceptedNo"] != "UJI-NOMOR" || a["AcceptanceStatus"] != "1" ||
		a["AcceptedDate"] != "2026-10-08 10:00:00" {
		t.Fatalf("S16.9: %+v", a)
	}
	// nomor sudah ada -> S16 dilewati
	if r := models.Rencanakan(k, klaimUji("UJI-SUDAH"), models.Keputusan{AcceptStatus: "1", Comment: "x"}, "UJI-K2",
		saatUji); r.TerbitkanNomor || !r.SimpanAkseptasi {
		t.Fatalf("AcceptedNo terisi: nomor %v akseptasi %v", r.TerbitkanNomor, r.SimpanAkseptasi)
	}
}

func TestRencanaTolakMenghentikanTangga(t *testing.T) {
	r := models.Rencanakan(kasusUji(1), klaimUji(""), models.Keputusan{AcceptStatus: "2", Comment: "UJI-TOLAK"}, "UJI-K1",
		saatUji)
	if len(r.Tangga) != 2 || r.Tangga[0].ID != "L1" || r.Tangga[0].Keputusan != "2" || r.Tangga[1].ID != "L2" ||
		r.Tangga[1].Keputusan != "2" || r.Tangga[1].IsiKomentar {
		t.Fatalf("S6 + S26.1: %+v", r.Tangga)
	}
	if r.Klaim.Adjustment["AcceptanceStatus"] != "2" || r.Klaim.Header["AktifButton"] != "0" ||
		r.Klaim.Adjustment["Notes"] != "UJI-TOLAK" {
		t.Fatalf("S25 / S27: %+v %+v", r.Klaim.Adjustment, r.Klaim.Header)
	}
	if r.Klaim.Riwayat[0].Teks != "Rejected by UJI-JABATAN-1" || r.StatusRiwayat != "REJECT" {
		t.Fatalf("riwayat: %+v / %q", r.Klaim.Riwayat, r.StatusRiwayat)
	}
	if r.Count != 3 || !r.Selesai || r.Email != models.EmailPembuatTolak || r.TerbitkanNomor || r.Kasir {
		t.Fatalf("tolak: %+v", r)
	}
}

func TestRencanaSubjectivityTingkatSatu(t *testing.T) {
	k := kasusUji(1)
	k.Loop = 1
	k.Tangga = k.Tangga[:1]
	r := models.Rencanakan(k, klaimUji(""), models.Keputusan{AcceptStatus: "1", Comment: "UJI", IsSubjectivity: true,
		SubjectivityNote: "UJI-NOTE"}, "UJI-K1", saatUji)
	if !r.Subjectivity || r.TerbitkanNomor || r.SimpanAkseptasi || r.SimpanRetro || r.Kasir {
		t.Fatalf("subjectivity melewati S16 / S17 / S21 / S34: %+v", r)
	}
	if !r.JSONKlaim || !r.Konversi {
		t.Fatalf("S28 / S29 tidak bergerbang subjectivity: %+v", r)
	}
	a := r.Klaim.Adjustment
	if a["IsKomite"] != "0" || a["IsSubjectivity"] != "true" || a["SubjectivityNote"] != "UJI-NOTE" ||
		r.Klaim.Header["ClaimData.IsSubjectivity"] != "true" {
		t.Fatalf("S23 / S24: %+v %+v", a, r.Klaim.Header)
	}
	if _, ada := r.Klaim.Header["IsOutstanding"]; ada {
		t.Fatal("SaveAcceptation_Act dilewati: IsOutstanding tidak disentuh")
	}
}

func TestTotalPenyesuaianSetKomiteList(t *testing.T) {
	rows := []map[string]string{
		{"Currency": "UJA", "GrossAdjustment": "100", "AdjustmentValue": "10", "KursIDR": "2", "AcceptanceStatus": "1"},
		{"Currency": "UJB", "GrossAdjustment": "50.5", "AdjustmentValue": "5.25", "KursIDR": "3"},
		{"Currency": "UJA", "GrossAdjustment": "1000", "AdjustmentValue": "100", "KursIDR": "2", "AcceptanceStatus": "2"},
		{"Currency": "UJA", "GrossAdjustment": "0.0001", "AdjustmentValue": "0.00005", "KursIDR": "2"},
	}
	got, err := models.TotalPenyesuaian(rows)
	if err != nil {
		t.Fatal(err)
	}
	mau := []models.TotalMataUang{
		{Currency: "UJB", AdjustmentGross: "50.5", AdjustmentValue: "5.25"},
		{Currency: "UJA", AdjustmentGross: "100.0001", AdjustmentValue: "10.00005"},
		{Currency: models.LabelTotalIDR, AdjustmentGross: "351.5002", AdjustmentValue: "35.7501"},
	}
	if len(got) != len(mau) {
		t.Fatalf("baris total %+v", got)
	}
	for i := range mau {
		if got[i] != mau[i] {
			t.Errorf("baris %d: %+v, mau %+v", i, got[i], mau[i])
		}
	}
	if _, err := models.TotalPenyesuaian([]map[string]string{{"Currency": "UJA", "GrossAdjustment": "x"}}); err == nil {
		t.Fatal("angka rusak harus galat, bukan ditebak")
	}
}

func TestRakitNomorAkseptasi(t *testing.T) {
	if got := models.RakitNomorAkseptasi("UJI-"+models.HurufAkseptasi, "12", "10.2026", 7); got != "UJI-A12.10.2026.TP00007" {
		t.Fatalf("nomor %q", got)
	}
}

func klaimAkhir() kontrak.KlaimTreaty {
	return kontrak.KlaimTreaty{Adjustment: 1,
		Nilai: map[string]string{"pyID": "CLMP-UJI001", "ClaimData.NoClaim": "UJI-K-0001",
			"ClaimData.PolicyData.PolicyNo": "UJI-POLIS", "ClaimData.IDMaster": "UJI-M1",
			"ClaimData.CauseOfLoss": "UJI-SEBAB", "ClaimData.CauseOfLossID": "UJI-S1", "TreatyInMaster.RNMShareP": "25",
			"ClaimData.InsuredName": "UJI-TERTANGGUNG", "ClaimData.DateOfLoss": "2026-02-07",
			"OfferFacIn.QuotationData.BusinessOldId": "12"},
		Daftar: map[string][]map[string]string{
			"ClaimData.AdjustmentList": {{"ID": "UJI-ADJ-1", "Type": "2", "Currency": "UJA", "CurrencyID": "UJI-A",
				"GrossAdjustment": "1000", "AdjustmentValue": "250", "KursIDR": "2", "PayableTo": "UJI-PENERIMA",
				"NoAccount": "12-34 56", "AcceptedNo": "UJI-A12.10.2026.TP00001", "AcceptedDate": "2026-10-28 10:00:00",
				"IndividualRiskRNM": "0", "IndividualRiskPercentage": "0", "IDOfBank": "UJI-BANK"}},
			"ClaimData.AdjustmentList(1).SpreadingAdjustment": {{"TreatyType": "UJI-T1"}, {"TreatyType": "UJI-T2"}},
		}}
}

func TestSusunOSAkseptasiSaveAcceptation(t *testing.T) {
	kl := klaimAkhir()
	kl.Daftar["ClaimData.AdjustmentList"][0]["pxCreateOperator"] = "UJI-ADMIN"
	kl.Daftar["ClaimData.AdjustmentList"][0]["ProposeAdjustmentValue"] = "1200"
	os := models.SusunOSAkseptasi(kl, "UJI-NOMOR", "TKMT-UJI001", saatUji)
	// DATA_JSON = halaman TempOSAkseptasi (SaveAcceptation_Act S1, S4): kunci urut tanpa membedakan huruf besar, LF ","
	// di depan pasangan berikutnya, nol spasi.
	pasangan := []string{`"AcceptedNo":"UJI-NOMOR"`, `"CauseOfLoss":"UJI-SEBAB"`, `"CauseOfLossID":"UJI-S1"`,
		`"Currency":"UJA"`, `"CurrencyID":"UJI-A"`, `"EstimationDate":"20261008"`, `"GrossValue":"1000"`,
		`"KomiteNo":"TKMT-UJI001"`, `"KursValue":"2"`, `"NoClaim":"UJI-K-0001"`, `"PaymentType":"2"`,
		`"PersenRNM":"25"`, `"pxCreateOperator":"UJI-ADMIN"`, `"pxObjClass":"ASM-FW-GCNMFW-Data-osAkseptasi"`,
		`"TotalGross":"1200"`, `"Type":"1"`, `"TypeID":"2"`, `"Value":"250"`}
	mau := models.BarisOSAkseptasi{CaseID: "CLMP-UJI001", NoClaim: "UJI-K-0001", NoPolis: "UJI-POLIS", StsReject: "1",
		MasterID: "UJI-M1", CauseOfLoss: "UJI-SEBAB", CauseOfLossID: "UJI-S1", CurrencyID: "UJI-A", Currency: "UJA",
		GrossValue: "1000", PersenRNM: "25", Value: "250", KursValue: "2", Type: "1", AcceptedNo: "UJI-NOMOR",
		PaymentType: "2", EstimationDate: saatUji, DataJSON: "{" + string(rune(10)) +
			strings.Join(pasangan, string(rune(10))+",") + string(rune(10)) + "}" + string(rune(10))}
	if os != mau {
		t.Fatalf("baris OS %+v\nmau %+v", os, mau)
	}
}

func TestJSONHalamanPegaTanpaEscapeHTML(t *testing.T) {
	got := models.JSONHalamanPega(map[string]string{"B": `a/b <c> "d"`, "a": "1", "Kosong": ""})
	mau := "{" + string(rune(10)) + `"a":"1"` + string(rune(10)) + "," + `"B":"a/b <c> \"d\""` + string(rune(10)) + "}" +
		string(rune(10))
	if got != mau {
		t.Fatalf("JSON %q, mau %q", got, mau)
	}
}

func TestSusunRetroSaveAcceptationTreaty(t *testing.T) {
	retro := []models.BarisRetro{{ReinsurerID: "UJI-RE1", Name: "UJI-REAS", RiComm: "5", PctShare: "40"}}
	h, err := models.SusunRetro("250", "100", retro, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(h.FacRetro) != 1 || h.FacRetro[0]["PctShareAllObj"] != "40" || !h.IsFacRetro || h.Pesan != models.PesanCetakDLA {
		t.Fatalf("retro kosong + melewati batas: %+v", h)
	}
	h, _ = models.SusunRetro("50", "100", retro, true)
	if len(h.FacRetro) != 0 || h.IsFacRetro {
		t.Fatalf("retro sudah ada + di bawah batas: %+v", h)
	}
	if h, _ := models.SusunRetro("1", "", nil, true); !h.IsFacRetro {
		t.Fatal("batas tak terbaca = @toDecimal(\"\") = 0")
	}
	if tt, ada := models.SpreadingAdjustmentTerakhir(klaimAkhir()); !ada || tt != "UJI-T2" {
		t.Fatalf("putaran terakhir S4: %q", tt)
	}
}

func TestMuatanKasirDanEmail(t *testing.T) {
	kl := klaimAkhir()
	cfg := models.KonfigurasiKasir{CompanyName: "UJI-CO", LjtdID: "UJI-LJ", LdcID: "UJI-LDC", LdcIDSyariah: "UJI-LDCS"}
	m := models.SusunMuatanKasir(kl, models.AdjustmentKlaim(kl), "uji@contoh.invalid", false, cfg, "UJI-K2", saatUji)
	if m.AccountNo != "123456" || m.TglAksep != "28-10-2026" || m.TglBolehBayar != "01-12-2026" || m.LbuId != "12" ||
		m.LdcId != "UJI-LDC" || m.CompanyName != "UJI-CO" || m.UserInput != "UJI-K2" || m.Nett != "250" {
		t.Fatalf("muatan kasir %+v", m)
	}
	if !models.PanjangNoAksepCLMP("UJI-A12.10.2026.TP00001") || models.PanjangNoAksepCLMP("UJI") {
		t.Fatal("panjang nomor CLMP 23/24")
	}
	if got := models.SubjekEmail(models.EmailPembuatSetuju, kl, "TKMT-UJI001"); got !=
		"(Approval) Pengajuan Akseptasi : CLMP-UJI001/TKMT-UJI001 UJI-TERTANGGUNG DOL 07 Febuari 2026" {
		t.Fatalf("subjek %q", got)
	}
	if got := models.SubjekEmail(models.EmailPenyetujuBerikut, kl, "X"); !strings.HasPrefix(got, "Pengajuan Akseptasi : ") {
		t.Fatalf("subjek berikut %q", got)
	}
	if !models.AdaSpreadingAdjustment(kl) {
		t.Fatal("S1 SendEmailKlaim_KMT: spreading adjustment ada")
	}
}
