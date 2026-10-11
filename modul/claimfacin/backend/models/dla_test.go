package models_test

// Uji tabel Print DLA (ChooseDla_Act, GenerateDLAFacin_Act, DLAFacintoTreaty_Act), Acceptation dan muatan kasir
// (SaveAcceptation, HitServiceToKasir_Act jalur CLM, GetStatusKasir_Act), premi (CekPremiLunas_Act). Angka `UJI`.

import (
	"testing"

	"nusantarare/modul/claimfacin/backend/models"
)

// halamanDLA - objek retro dengan satu item (TSI RNM 250 jt, spreading polis FAC 40%) dan adjustment `adj`.
func halamanDLA(adj ...models.Baris) *models.Halaman {
	h := halamanAdj(models.Baris{})
	h.SetelDaftar(models.DaftarObjek, []models.Baris{{"IsFacretro": "1", "PlaStatus": "1", "RemarksDLA": "- Refer to PLA : "}})
	h.SetelDaftar(models.DaftarDiItem(1, 1, models.AnakSpreadPolis), []models.Baris{
		{"TreatyType": "10003", "TreatyName": "UJI QS FAC", "SharePercentage": "60"},
		{"TreatyType": models.TreatyFacRetro, "TreatyName": "UJI FAC RETRO", "SharePercentage": "40"}})
	h.SetelDaftar(models.DaftarFacRetro, []models.Baris{{"ReinsurerID": "UJI-R1", "ReinsurerName": "UJI REAS",
		"PctShareAllObj": "100", "AdditionalInfo": "UJI INFO"}})
	h.SetelDaftar(models.DaftarAdj(1, 1), adj)
	return h
}

func layak(no, pt, nilai string) models.Baris {
	b := models.Baris{"AcceptanceStatus": "1", "IsPrintAccept": "1", "AcceptedNo": no, "PaymentType": pt,
		"IsFacRetro": "1"}
	switch pt {
	case models.BayarSalvage:
		b["SalvageValue"] = nilai
	case models.BayarFee, models.BayarExpense:
		b["AdjusterFeeValue"] = nilai
	default:
		b["AdjustmentValue"] = nilai
	}
	return b
}

func TestChooseDLA(t *testing.T) {
	h := halamanDLA()
	r, err := models.ChooseDLA(h, 1)
	if err != nil {
		t.Fatal(err)
	}
	ob, _ := models.Objek(h, 1)
	if r.Keluar || len(r.Urutan) != 1 || r.Urutan[0] != models.JenisNomorDLAFac || ob["DLAStatus"] != "1" ||
		h.Ambil(models.JalurIsTreatyOut) != "1" {
		t.Fatalf("rencana %+v objek %v", r, ob)
	}
	// gerbang treaty CekLimit.CARI1 (OQ-CFI-24): hanya bila 1
	h = halamanDLA()
	h.Setel(models.JalurCekLimit, "1")
	if r, _ = models.ChooseDLA(h, 1); len(r.Urutan) != 2 || r.Urutan[0] != models.JenisNomorDLATreaty {
		t.Fatalf("dengan CekLimit 1: %+v", r)
	}
	// 1: retro tanpa PLA keluar
	h = halamanDLA()
	h.SetelDaftar(models.DaftarObjek, []models.Baris{{"IsFacretro": "1"}})
	if r, _ = models.ChooseDLA(h, 1); !r.Keluar {
		t.Fatalf("retro tanpa PLA tidak keluar: %+v", r)
	}
}

func TestDLASekaliPerJenisBukanPerBaris(t *testing.T) {
	// [penyimpangan sadar]: Pega memanggil GenerateDLAFacin_Act untuk SETIAP baris 10015 setiap item.
	h := halamanDLA()
	h.SetelDaftar(models.DaftarItem(1), []models.Baris{{"TSINusare": "1"}, {"TSINusare": "1"}})
	for i := 1; i <= 2; i++ {
		h.SetelDaftar(models.DaftarDiItem(1, i, models.AnakSpreadPolis), []models.Baris{
			{"TreatyType": models.TreatyFacRetro, "SharePercentage": "40"}})
	}
	r, _ := models.ChooseDLA(h, 1)
	if len(r.Urutan) != 1 {
		t.Fatalf("dua item retro: %+v", r)
	}
}

func TestTerapkanDLAFac(t *testing.T) {
	h := halamanDLA(layak("UJI-A1", models.BayarFinal, "4500000"),
		models.Baris{"AcceptanceStatus": "", "PaymentType": models.BayarInterim, "AdjustmentValue": "9"})
	hasil, err := models.TerapkanDLAFac(konteksUji(), h, 1, "UJI-P77.03.2026.00001")
	if err != nil {
		t.Fatal(err)
	}
	adj := h.AmbilDaftar(models.DaftarAdj(1, 1))
	it, _ := models.Item(h, 1, 1)
	ob, _ := models.Objek(h, 1)
	// facout = 250 jt x 40% = 100 jt; klaim = 4,5 jt x 40% = 1,8 jt; reas 100% -> 1,8 jt
	if hasil.Keluar || len(hasil.Akseptasi) != 1 || hasil.Akseptasi[0] != "UJI-A1" ||
		adj[0]["DLA_No"] != "UJI-P77.03.2026.00001" || adj[1]["DLA_No"] != "" ||
		it["TotalEstimasiReas"] != "1800000" || ob["NoDLA"] != "UJI-P77.03.2026.00001" {
		t.Fatalf("hasil %+v adj %v item %v", hasil, adj, it)
	}
	// 15.2: RemarksDLA + FacRetroList ke SETIAP adjustment item
	for a := 1; a <= 2; a++ {
		fr := h.AmbilDaftar(models.DaftarDiAdj(1, 1, a, models.AnakFacRetro))
		if adj[a-1]["RemarksDLA"] != "- Refer to PLA : " || len(fr) != 1 || fr[0]["ReinsurerID"] != "UJI-R1" {
			t.Fatalf("adjustment %d: %v retro %v", a, adj[a-1], fr)
		}
	}
	if kr := h.AmbilDaftar(models.DaftarKronologi); len(kr) != 1 || kr[0]["pyNote"] != models.TeksPrintDLA {
		t.Fatalf("kronologi %v", kr)
	}
}

func TestTerapkanDLAFacTanpaAdjustmentLayakKeluar(t *testing.T) {
	h := halamanDLA(models.Baris{"AcceptanceStatus": "1", "IsPrintAccept": "", "PaymentType": models.BayarFinal})
	hasil, err := models.TerapkanDLAFac(konteksUji(), h, 1, "UJI-P1")
	if err != nil {
		t.Fatal(err)
	}
	ob, _ := models.Objek(h, 1)
	if !hasil.Keluar || len(hasil.Akseptasi) != 0 || ob["NoDLA"] != "UJI-P1" {
		t.Fatalf("%+v objek %v", hasil, ob)
	}
}

func TestDraftDLAHanyaRetro(t *testing.T) {
	h := halamanDLA(models.Baris{"IsFacRetro": ""})
	if err := models.DraftDLA(konteksUji(), h, 1, 1, 1); err != nil {
		t.Fatal(err)
	}
	if len(h.AmbilDaftar(models.DaftarKronologi)) != 0 {
		t.Fatal("draft DLA untuk adjustment bukan retro")
	}
	h = halamanDLA(models.Baris{"IsFacRetro": "1"})
	if err := models.DraftDLA(konteksUji(), h, 1, 1, 1); err != nil {
		t.Fatal(err)
	}
	ob, _ := models.Objek(h, 1)
	if ob["ShareRetro"] != "1" || h.AmbilDaftar(models.DaftarKronologi)[0]["pyNote"] != models.TeksDraftDLA ||
		ob["NoDLA"] != "" {
		t.Fatalf("objek %v", ob)
	}
}

func TestSaveAcceptation(t *testing.T) {
	h := halamanDLA(models.Baris{"AcceptanceStatus": "1", "PaymentType": models.BayarFinal})
	ok, err := models.SaveAcceptation(konteksUji(), h, 1, 1, 1)
	if err != nil || !ok {
		t.Fatal(ok, err)
	}
	if err := models.SelesaiAkseptasi(h, 1, 1, 1); err != nil {
		t.Fatal(err)
	}
	b, ob := adjSatu(t, h), models.Baris{}
	ob, _ = models.Objek(h, 1)
	if b["IsPrintAccept"] != "1" || b["IsFacRetro"] != "1" || ob["IsPrintAccept"] != "1" || ob["DLAStatus"] != "0" {
		t.Fatalf("adj %v objek %v", b, ob)
	}
	if ok, _ := models.SaveAcceptation(konteksUji(), h, 1, 1, 1); ok {
		t.Fatal("Acceptation kedua tidak keluar di langkah 1")
	}
}

func TestBolehAkseptasi(t *testing.T) {
	kasus := []struct {
		b    models.Baris
		mau  bool
		nama string
	}{
		{models.Baris{"AcceptanceStatus": "1"}, true, "baru diterima"},
		{models.Baris{"AcceptanceStatus": ""}, false, "belum diterima"},
		{models.Baris{"AcceptanceStatus": "1", "IsPrintAccept": "1", "DirectToKasir": "false"}, false, "cetak, tanpa kasir"},
		{models.Baris{"AcceptanceStatus": "1", "IsPrintAccept": "1", "DirectToKasir": "true"}, true, "kirim ulang kasir"},
		{models.Baris{"AcceptanceStatus": "1", "IsPrintAccept": "1", "DirectToKasir": "true", "StatusKasir": "x"}, false,
			"kasir berstatus"},
	}
	for _, c := range kasus {
		if got := models.BolehAkseptasi(c.b); got != c.mau {
			t.Errorf("%s: %v", c.nama, got)
		}
	}
}

func TestTglBolehBayar(t *testing.T) {
	kasus := []struct {
		tgl   string
		kini  int
		hasil string
	}{
		{"20260305", 3, "05-04-2026"},
		{"20260326", 3, "01-05-2026"},
		{"20261110", 11, "10-12-2026"},
		{"20261210", 12, "10-01-2027"}, // bulan 13 -> "1" -> "01"; tahun naik (bulan berjalan Desember)
		{"20261126", 12, "01-01-2027"}, // hari > 25 bulan Nov -> Jan
		{"20261126", 11, "01-01-2026"}, // tahun hanya naik bila bulan berjalan Desember
		{"", 3, ""},
	}
	for _, c := range kasus {
		if got := models.TglBolehBayar(c.tgl, c.kini); got != c.hasil {
			t.Errorf("TglBolehBayar(%q, %d) = %q, mau %q", c.tgl, c.kini, got, c.hasil)
		}
	}
}

func TestMuatanKasirMenjumlahAkseptasiSama(t *testing.T) {
	h := halamanDLA(
		models.Baris{"AcceptedNo": "UJI-AKS.03.2026.00001", "PaymentType": models.BayarFinal, "AdjustmentValue": "1000",
			"IndividualRiskRNM": "10", "IndividualRiskPercentage": "1", "AcceptedDate": "2026-03-05 10:00:00",
			"NoAccount": "12-34", "PayableTo": "UJI SOB", "CurrencyID": idr, "IDOfBank": "UJI-B"},
		models.Baris{"AcceptedNo": "UJI-AKS.03.2026.00001", "PaymentType": models.BayarFee, "AdjusterFeeValue": "200",
			"IndividualRiskRNM": "5", "IndividualRiskPercentage": "2", "AcceptedDate": "2026-03-05 10:00:00",
			"NoAccount": "56-78", "PayableTo": "UJI SOB", "CurrencyID": idr, "IDOfBank": "UJI-B"},
		models.Baris{"AcceptedNo": "UJI-LAIN", "PaymentType": models.BayarFinal, "AdjustmentValue": "999"})
	b, _ := models.Adj(h, 1, 1, 1)
	if !models.PanjangNoAksepCLM(b["AcceptedNo"]) {
		t.Fatalf("panjang nomor akseptasi %d", len(b["AcceptedNo"]))
	}
	m, err := models.SusunMuatanKasir(konteksUji(), h, 1, b, "")
	if err != nil {
		t.Fatal(err)
	}
	if m.Nett != "1200" || m.Deductible != "15" || m.KaliDeduct != "3" || m.AccountNo != "5678" ||
		m.AcceptType != models.BayarFee || m.TglAksep != "05-03-2026" || m.TglBolehBayar != "05-04-2026" {
		t.Fatalf("%+v", m)
	}
}

func TestStatusKasirDanPremi(t *testing.T) {
	b := models.Baris{"AcceptanceStatus": "1", "DirectToKasir": "true"}
	models.TerapkanStatusKasir(b, models.KetKasirSukses, true)
	if b["StatusKasir"] != models.StatusKasirSukses {
		t.Fatalf("%v", b)
	}
	b = models.Baris{"AcceptanceStatus": "1", "DirectToKasir": "false"}
	models.TerapkanStatusKasir(b, "UJI GAGAL", true)
	if b["StatusKasir"] != "" {
		t.Fatalf("tanpa DirectToKasir: %v", b)
	}
	for saldo, mau := range map[string]bool{"": false, "0": false, "-5": false, "1,5": true, "100": true} {
		got, err := models.PremiBelumLunas(saldo)
		if err != nil || got != mau {
			t.Errorf("PremiBelumLunas(%q) = %v %v", saldo, got, err)
		}
	}
}
