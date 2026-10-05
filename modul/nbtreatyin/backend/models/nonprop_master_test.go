package models_test

// Uji seam 3 - pemulihan limit master (`TreatySetReinstatement`,
// `SetReinstatementPct`), pra-proses `TreatyRealizationCheckXOLList`, dan
// tampilan master saat berkas dibuka ulang. Nilai harapan dari langkah XML.

import (
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestTreatySetReinstatement(t *testing.T) {
	h := models.HalamanBaru()
	m := masterUji()
	m.Daftar["Limits"] = []models.Baris{
		{"ReinstatementValue": "2", "ReinstatementNote": "UJI-N", "Limit": "1000", "Limit2": "50"},
		{"ReinstatementValue": "3"},
		{"ReinstatementValue": ""},
	}
	m.Daftar["Limits(1).MDPList"] = []models.Baris{{"Currency": "IDR", "Value": "10"}, {"Currency": "USD", "Value": "2"}}
	m.Daftar["Limits(2).Reinstatement_List"] = []models.Baris{{"ReinstatementValue": "1"}}
	models.TerapkanMasterXOL(h, m)
	if err := models.TreatySetReinstatement(h); err != nil {
		t.Fatal(err)
	}
	r := h.AmbilDaftar("TreatyIn.Limits(1).Reinstatement_List")
	if len(r) != 2 {
		t.Fatalf("SetReinstatementPct 3: ReinstatementValue 2 -> 2 baris, dapat %d", len(r))
	}
	// 3.1 + 3.3 (MDPList dua baris: (1) IDR, (2) USD)
	harap := map[string]string{"ReinstatementPct": "100", "ReinstatementValue": "2", "ReinstatementNote": "UJI-N",
		"ReinstatementAmount1": "1000", "ReinstatementAmount2": "50", "ID": "1", "AdditionalPct": "100",
		"AdditionalAmount1": "10", "AdditionalAmount2": "2"}
	for k, v := range harap {
		if r[1][k] != v {
			t.Errorf("baris 2 %s = %q, harap %q", k, r[1][k], v)
		}
	}
	if r[0]["ReinstatementValue"] != "1" {
		t.Errorf("baris 1 ReinstatementValue %q", r[0]["ReinstatementValue"])
	}
	// TreatySetReinstatement 1.1: Reinstatement_List(1).ReinstatementValue terisi -> dilewati.
	if r2 := h.AmbilDaftar("TreatyIn.Limits(2).Reinstatement_List"); len(r2) != 1 {
		t.Errorf("limit 2 tidak disentuh: %d baris", len(r2))
	}
	if r3 := h.AmbilDaftar("TreatyIn.Limits(3).Reinstatement_List"); len(r3) != 0 {
		t.Errorf("ReinstatementValue kosong -> nol pemulihan: %d baris", len(r3))
	}
}

func TestSetReinstatementPctSatuMDP(t *testing.T) {
	h := models.HalamanBaru()
	models.TerapkanMasterXOL(h, models.MasterXOL{Daftar: map[string][]models.Baris{
		"Limits":            {{"ReinstatementValue": "1"}},
		"Limits(1).MDPList": {{"Currency": "USD", "Value": "7"}},
	}})
	if err := models.SetReinstatementPct(h, 1); err != nil {
		t.Fatal(err)
	}
	// 3.2: satu baris MDP ber-USD -> Amount1 0, Amount2 7.
	b := h.AmbilDaftar("TreatyIn.Limits(1).Reinstatement_List")[0]
	if b["AdditionalAmount1"] != "0" || b["AdditionalAmount2"] != "7" {
		t.Fatalf("baris %+v", b)
	}
}

func TestPerluCekDaftarXOL(t *testing.T) { // InputPolicyTreatyInPre_Act 10
	kasus := []struct {
		isNew, edm, jenis string
		harap             bool
	}{
		{"1", "", "NonProportional", true},
		{"1", "3", "NonProportional", false},
		{"0", "", "NonProportional", false},
		{"1", "", "Proportional", false}, // penjaga AC 33
	}
	for _, c := range kasus {
		h := models.HalamanBaru()
		h.Setel("PolicyTreatyIn.IsNewPolicyNonProp", c.isNew)
		h.Setel("PolicyTreatyIn.EDMType", c.edm)
		h.Setel("PolicyTreatyIn.QuotationData.ProportionalType", c.jenis)
		if got := models.PerluCekDaftarXOL(h); got != c.harap {
			t.Errorf("%+v -> %v", c, got)
		}
	}
}

func TestCekDaftarXOL(t *testing.T) { // TreatyRealizationCheckXOLList
	h := halamanXOL("", "")
	h.TambahPesan("", "UJI-pesan")
	if !models.AwalCekDaftarXOL(h) || len(h.SemuaPesan()) != 0 {
		t.Fatal("langkah 1-2: pesan dibersihkan, daftar kosong -> lanjut")
	}
	if err := models.LengkapiDaftarXOL(h, idUji); err != nil {
		t.Fatal(err)
	}
	if n := len(h.AmbilDaftar("PolicyTreatyIn.TreatyXOLList")); n != 2 || len(h.SemuaPesan()) != 0 {
		t.Fatalf("langkah 6: %d baris, pesan %v", n, h.SemuaPesan())
	}
	if models.AwalCekDaftarXOL(h) {
		t.Fatal("langkah 2: daftar terisi -> keluar")
	}
	kosong := models.HalamanBaru()
	if err := models.LengkapiDaftarXOL(kosong, idUji); err != nil {
		t.Fatal(err)
	}
	if p := kosong.SemuaPesan(); len(p) != 1 || p[0] != "Error fetching XolList" {
		t.Fatalf("langkah 8: %v", p)
	}
}

func TestTampilanMasterNonPropTidakMenyentuhPolis(t *testing.T) {
	h := halamanXOL("true", "Inclusive")
	m := masterUji()
	m.Daftar["LimitShareSummaryList"] = []models.Baris{{"Limit": "1", "Deductible": "102.2", "NetPremi": "900"}}
	m.Daftar["TotalShareDeductionNP"] = []models.Baris{{"Currency": "IDR"}}
	models.TerapkanMasterXOL(h, m)
	h.SetelDaftar("PolicyTreatyIn.ListInstallment", []models.Baris{{"Currency": "IDR", "PaymentTotal": "10"}})
	if err := models.TampilanMasterNonProp(h); err != nil {
		t.Fatal(err)
	}
	sama(t, "PPNValue layer", h.AmbilDaftar("TreatyIn.LimitShareSummaryList")[0]["PPNValue"], "2.2")
	if _, ada := h.AmbilDaftar("PolicyTreatyIn.ListInstallment")[0]["PPN"]; ada {
		t.Fatal("tampilan master tidak menulis medan polis tersimpan")
	}
}

func TestPilihBisnisUlangMembuangRincianAngsuran(t *testing.T) {
	// preACT langkah 11 `Property-Remove ListInstallment`: PageList dibuang beserta
	// halaman bersarangnya - rincian NonProp lama tidak tertinggal saat kontrak
	// dipilih ulang (penyimpanan polis Proportional menolak rincian bersarang, AC 31).
	h := models.HalamanBaru()
	h.SetelDaftar("PolicyTreatyIn.ListInstallment", []models.Baris{{"Currency": "IDR"}})
	h.SetelDaftar("PolicyTreatyIn.ListInstallment(1).InstallmentList", []models.Baris{{"Premium": "1"}})
	models.TerapkanDetailKontrak(h, models.BarisKontrak{"PROPORTIONTYPE": "Proportional"})
	if len(h.AmbilDaftar("PolicyTreatyIn.ListInstallment(1).InstallmentList")) != 0 {
		t.Fatal("rincian bersarang ikut terhapus")
	}
}
