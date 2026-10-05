package handlers_test

// Uji seam 1 - jalur NB NonProporsional / XOL UJUNG KE UJUNG lewat HTTP
// (`[keputusan work owner]` K8): pilih bisnis kontrak NonProportional ->
// master XOL dibaca (tiruan `PembacaMasterTreaty`) -> InputPolicyTreatyInDetail_NonProp
// -> hitung -> simpan -> baca kembali -> tangga sampai selesai. Nilai harapan
// dihitung tangan dari langkah XML (lihat models/nonprop_test.go). Fixture UJI-.

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/services"
)

// masterNP - dua layer, angsuran IDR dan USD, ringkasan layer untuk preACT 18.
func masterNP() models.MasterXOL {
	return models.MasterXOL{
		Nilai: map[string]string{"RNMShare": "10", "FacultativeShare": "0", "ProportionType": "NonProportional"},
		Daftar: map[string][]models.Baris{
			"Share": {
				{"Layer": "1", "LayerType": "UJI-LT", "SpreadingTypeXOL": "UJI-SPR", "SpreadingTypeIDXOL": "UJI-SPR-ID"},
				{"Layer": "2", "LayerType": "UJI-LT"},
			},
			"Share(1).GrossPremiumList":   {{"Currency": "IDR", "Value": "1000"}, {"Currency": "USD", "Value": "50"}},
			"Share(1).NetPremiumList":     {{"Currency": "IDR", "Value": "900"}, {"Currency": "USD", "Value": "45"}},
			"Share(1).DeductionList":      {{"Currency": "IDR", "Deduction": "102.2"}, {"Currency": "USD", "Deduction": "10.22"}},
			"Share(1).DeductionTotalList": {{"Currency": "IDR", "Value": "306.6"}},
			"Share(1).RnmLimitList":       {{"Currency": "IDR", "Value": "5000"}},
			"Share(2).GrossPremiumList":   {{"Currency": "IDR", "Value": "2000"}},
			"Share(2).NetPremiumList":     {{"Currency": "IDR", "Value": "1800"}},
			"Share(2).DeductionList":      {{"Currency": "IDR", "Deduction": "204.4"}},
			"Installment": {
				{"Currency": "IDR", "AmountTotal": "2700", "PctTotal": "100"},
				{"Currency": "USD", "AmountTotal": "45", "PctTotal": "100"},
			},
			"Installment(1).InstallmentList": {
				{"Installment": "1", "PaymentDate": "2026-11-01", "InstallmentPct": "40", "Amount": "1080", "Currency": "IDR"},
				{"Installment": "2", "PaymentDate": "2027-02-01", "InstallmentPct": "60", "Amount": "1620", "Currency": "IDR"},
			},
			"Installment(2).InstallmentList": {
				{"Installment": "1", "PaymentDate": "2026-11-01", "InstallmentPct": "100", "Amount": "45", "Currency": "USD"},
			},
			"LimitShareSummaryList": {
				{"Note": "UJI", "Limit": "1000", "Limit2": "0", "Deductible": "102.2", "Deductible2": "0", "NetPremi": "900", "NetPremi2": "0"},
			},
			"TotalShareDeductionNP": {{"Currency": "IDR", "Value": "306.6"}},
		},
	}
}

func kontrakNP(u *uji) {
	u.g.Kontrak["UJI-D-NP"] = models.BarisKontrak{"ID": "UJI-D-NP", "TREATYID": "UJI-T-NP", "LIMITCURRENCY": "IDR",
		"CLASSOFBUSINESS": "UJI BISNIS", "TREATYTYPE": "XOL", "PROPORTIONTYPE": "NonProportional",
		"COMMENCEMENT": "2026-10-01 00:00:00", "TERMINATION": "2027-09-30 00:00:00"}
	u.g.Bisnis["UJI BISNIS"] = models.BarisBisnis{OldID: "01", GroupPanel: "006", ID: "UJI-B1"}
	u.g.Master["UJI-T-NP"] = masterNP()
}

// layar membaca jawaban JSON `services.Layar`.
func (u *uji) layar(isi string) services.Layar {
	u.t.Helper()
	var ly services.Layar
	if err := json.Unmarshal([]byte(isi), &ly); err != nil {
		u.t.Fatalf("%v: %s", err, isi)
	}
	return ly
}

func angkaSamaTeks(t *testing.T, apa, dapat, harap string) {
	t.Helper()
	a, err1 := models.AngkaTeks(apa, dapat)
	b, err2 := models.AngkaTeks(apa, harap)
	if err1 != nil || err2 != nil || dapat == "" || a.Cmp(b) != 0 {
		t.Errorf("%s = %q, harap %s", apa, dapat, harap)
	}
}

func TestNonPropPilihBisnisHitungSimpanBacaKembali(t *testing.T) {
	u := baru(t)
	kontrakNP(u)
	id := u.buat()
	// FlagPPH/TypeTax diisi layar SEBELUM pilih bisnis (centang FlagPPH hanya
	// menjalankan RemoveTypeTax_ACT - nilai XOL dihitung saat pilih bisnis).
	awal := models.HalamanBaru()
	awal.Setel("PolicyTreatyIn.FlagPPH", "true")
	awal.Setel("PolicyTreatyIn.TypeTax", "Inclusive")
	kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D-NP", "halaman": awal})
	if kode != http.StatusOK {
		t.Fatalf("pilih bisnis NonProp: %d %s", kode, isi)
	}
	s := u.g.Halaman[id] // yang TERSIMPAN
	// preACT 16 -> NonProp 13, 16-19 (hitung tangan: models/nonprop_test.go)
	for j, harap := range map[string]string{
		"PolicyTreatyIn.PremiOgp": "3000", "PolicyTreatyIn.NetPremium": "2700", "PolicyTreatyIn.Deduction1": "306.6",
		"PolicyTreatyIn.PPNValue": "6.6", "PolicyTreatyIn.PPHValue": "6", "PolicyTreatyIn.BalanceDueTo": "2712.6",
		"PolicyTreatyIn.ShareValue": "5000",
	} {
		angkaSamaTeks(t, j, s.Ambil(j), harap)
	}
	if s.Ambil("PolicyTreatyIn.DueTo") != "1" || s.Ambil("PolicyTreatyIn.StartDate") != "2026-10-01 00:00:00" {
		t.Errorf("NonProp 13: DueTo %q StartDate %q", s.Ambil("PolicyTreatyIn.DueTo"), s.Ambil("PolicyTreatyIn.StartDate"))
	}
	// T_POLIS_XOL / T_POLIS_XOL_LAYER (katalog TabelXOL, TabelLayerXOL)
	xol := s.AmbilDaftar(models.TabelXOL.Daftar)
	if len(xol) != 2 {
		t.Fatalf("TreatyXOLList tersimpan %d baris, harap 2", len(xol))
	}
	angkaSamaTeks(t, "XOL(1).GrossPremi", xol[0]["GrossPremi"], "3000")
	angkaSamaTeks(t, "XOL(1).NetPremiAfterTax", xol[0]["NetPremiAfterTax"], "2712.6")
	if xol[0]["IDCurrency"] != "UJI-ID-IDR" || xol[1]["IDCurrency"] != "UJI-ID-USD" {
		t.Errorf("IDCurrency dari GetDataCurrencyByName_SQL: %q %q", xol[0]["IDCurrency"], xol[1]["IDCurrency"])
	}
	lap := s.AmbilDaftar(models.JalurAnak(models.TabelXOL.Daftar, 1, models.TabelLayerXOL.Daftar))
	if len(lap) != 2 || lap[0]["Layer"] != "1" || lap[1]["Layer"] != "2" {
		t.Fatalf("ValueList %v", lap)
	}
	// T_POLIS_INSTALMENT / T_POLIS_INSTALMENT_DETAIL + preACT 18 (PPN per angsuran)
	ang := s.AmbilDaftar(models.DaftarAngsuran)
	if len(ang) != 2 {
		t.Fatalf("ListInstallment %d baris", len(ang))
	}
	angkaSamaTeks(t, "Angsuran(IDR).PPN", ang[0]["PPN"], "2.2")
	rinci := s.AmbilDaftar(models.JalurAnak(models.DaftarAngsuran, 1, models.TabelAngsuranRinci.Daftar))
	if len(rinci) != 2 || rinci[0]["DueDate"] != "2026-11-01" {
		t.Fatalf("InstallmentList %v", rinci)
	}
	// preACT 18.3.4.2.1: PremiumAfterPPN = PaymentTotalAfterPPN (2700+2.2) x 40%.
	// `.PPN`/`.PPh` rincian tidak berkolom (diagram R58, PERBANDINGAN-KOLOM-DIAGRAM bab 5).
	angkaSamaTeks(t, "Rinci(1).PremiumAfterPPN", rinci[0]["PremiumAfterPPN"], "1080.88")
	// T_POLIS_SPREADING (TreatyNonPropSetSpreading)
	if sp := s.AmbilDaftar(models.DaftarSpreading); len(sp) != 1 || sp[0]["TreatyType"] != "UJI-SPR-ID" {
		t.Fatalf("SpreadingRiskList %v", sp)
	}
	// Halaman master TIDAK tersimpan (nol penulisan master; hanya TREATY_IN_ID).
	if s.Ambil("TreatyIn.RNMShare") != "" || len(s.AmbilDaftar("TreatyIn.Share")) != 0 || s.Ambil("TreatyIn.ID") != "UJI-D-NP" {
		t.Fatal("halaman master tidak disimpan; TreatyIn.ID tetap")
	}

	// Baca kembali: master dibaca ulang, IsNewPolicyNonProp 1 (preDT), daftar XOL dari simpanan.
	_, isi = u.panggil("GET", "/kasus/"+id, admin, nil)
	ly := u.layar(isi)
	if ly.Halaman.Ambil("PolicyTreatyIn.IsNewPolicyNonProp") != "1" {
		t.Fatal("InputPolicyTreatyIn_preDT 4: ProportionalType NonProportional -> IsNewPolicyNonProp 1")
	}
	if len(ly.Halaman.AmbilDaftar("TreatyIn.Share")) != 2 || ly.Halaman.Ambil("TreatyIn.RNMShare") != "10" {
		t.Fatal("subsection DetailPoliciesNonProportional membaca master lagi saat dibuka")
	}
	angkaSamaTeks(t, "tampilan LimitShareSummaryList.PPNValue", ly.Halaman.AmbilDaftar("TreatyIn.LimitShareSummaryList")[0]["PPNValue"], "2.2")
	for _, m := range ly.MedanWajib {
		if m == "PolicyTreatyIn.PremiOgp" || m == "PolicyTreatyIn.Deduction1" {
			t.Errorf("medan uang tersembunyi (kontainer .IsNewPolicyNonProp != 1) tidak wajib: %s", m)
		}
	}

	// Hitung: CountSpreading baris 1 (SpreadingRiskList di DetailPolicyTreatyInNonProportional).
	h := ly.Halaman
	h.Setel("PolicyTreatyIn.QuotationData.MOID", "UJI-MO-1")
	h.Setel("PolicyTreatyIn.IsApproved", "1")
	h.Setel("PolicyTreatyIn.Suggest", "UJI-catatan")
	h.Setel("PolicyTreatyIn.NetPremium", "999999") // turunan: layar tidak pernah menulisnya
	h.SetelDaftar(models.TabelXOL.Daftar, []models.Baris{{"GrossPremi": "UJI-KARANGAN"}})
	kode, isi = u.panggil("POST", "/kasus/"+id+"/hitung", admin, map[string]any{"urutan": []map[string]string{{"aksi": "CountSpreading"}}, "indeks": 1, "halaman": h})
	if kode != http.StatusOK {
		t.Fatalf("hitung: %d %s", kode, isi)
	}
	angkaSamaTeks(t, "PremiumSpreaded", u.layar(isi).Halaman.AmbilDaftar(models.DaftarSpreading)[0]["PremiumSpreaded"], "2700")

	// Simpan: CountNetPremi TIDAK dijalankan untuk NonProp (medan uang tersembunyi,
	// tak satu pun refresh memicunya) - NetPremium tetap 2700, bukan 3000-306.6.
	kode, isi = u.panggil("PUT", "/kasus/"+id, admin, map[string]any{"halaman": h})
	if kode != http.StatusOK {
		t.Fatalf("simpan: %d %s", kode, isi)
	}
	s = u.g.Halaman[id]
	angkaSamaTeks(t, "NetPremium sesudah simpan", s.Ambil("PolicyTreatyIn.NetPremium"), "2700")
	if x := s.AmbilDaftar(models.TabelXOL.Daftar); len(x) != 2 || x[0]["GrossPremi"] == "UJI-KARANGAN" {
		t.Fatalf("TreatyXOLList tidak tampil di section mana pun - tidak diterima dari layar: %v", x)
	}
	if len(s.AmbilDaftar(models.JalurAnak(models.DaftarAngsuran, 1, models.TabelAngsuranRinci.Daftar))) != 2 {
		t.Fatal("rincian angsuran tetap tersimpan")
	}

	// Tangga sampai selesai: SetValidateInstallment (2 x 100%) tidak berlaku -
	// rantai layar uang tidak terpicu untuk NonProp.
	if kode, isi := u.kirim(id, admin, h); kode != http.StatusOK {
		t.Fatalf("admin submit NonProp: %d %s", kode, isi)
	}
	if kode, isi := u.kirim(id, secHead, putusan("1")); kode != http.StatusOK {
		t.Fatalf("Sec Head NonProp (medan uang tersembunyi tidak wajib): %d %s", kode, isi)
	}
	kode, isi = u.panggil("POST", "/kasus/"+id+"/nomor-polis", deptHead, map[string]any{"halaman": putusan("1")})
	if kode != http.StatusOK || !strings.Contains(isi, "UJI-QR.") {
		t.Fatalf("nomor polis (DueTo 1 -> QR): %d %s", kode, isi)
	}
	if kode, isi := u.kirim(id, deptHead, putusan("1")); kode != http.StatusOK || u.g.Kasus[id].StatusWork != models.StatusSelesai {
		t.Fatalf("Dept Head: %d %s", kode, isi)
	}
	if n := len(u.g.Halaman[id].AmbilDaftar(models.TabelXOL.Daftar)); n != 2 {
		t.Fatalf("XOL tetap tersimpan sampai selesai: %d", n)
	}
}

func TestNonPropMasterTidakAdaMenghentikanPilihBisnis(t *testing.T) { // AC 36-38
	u := baru(t)
	kontrakNP(u)
	delete(u.g.Master, "UJI-T-NP")
	id := u.buat()
	kode, _ := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D-NP"})
	if kode != http.StatusUnprocessableEntity {
		t.Fatalf("master tidak ada: %d", kode)
	}
	if u.g.Halaman[id].Ambil("PolicyTreatyIn.NoOffer") != "" {
		t.Fatal("pembacaan gagal tidak menyimpan apa pun")
	}
	u.g.Master["UJI-T-NP"] = masterNP()
	u.g.MasterRusak = true
	if kode, _ := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D-NP"}); kode != http.StatusUnprocessableEntity {
		t.Fatalf("dokumen master rusak: %d", kode)
	}
}

func TestNonPropCekDaftarXOLSaatDibuka(t *testing.T) { // InputPolicyTreatyInPre_Act 10
	u := baru(t)
	kontrakNP(u)
	m := masterNP()
	m.Nilai["FacultativeShare"] = "5" // -> RetroShare; IsNewPolicyNonProp masih 0 saat pilih bisnis
	u.g.Master["UJI-T-NP"] = m
	id := u.buat()
	if kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D-NP"}); kode != http.StatusOK {
		t.Fatalf("%d %s", kode, isi)
	}
	if n := len(u.g.Halaman[id].AmbilDaftar(models.TabelXOL.Daftar)); n != 0 {
		t.Fatalf("RetroShare langkah 2 bersyarat IsNewPolicyNonProp==1: tersimpan %d baris", n)
	}
	// Dibuka: preDT -> IsNewPolicyNonProp 1 -> TreatyRealizationCheckXOLList ->
	// master dibaca (BrowseTreatyIn) -> InsertToTreatyXOLList.
	_, isi := u.panggil("GET", "/kasus/"+id, admin, nil)
	if n := len(u.layar(isi).Halaman.AmbilDaftar(models.TabelXOL.Daftar)); n != 2 {
		t.Fatalf("CheckXOLList mengisi %d baris, harap 2", n)
	}
	// Master hilang: pesan VERBATIM langkah 8, layar tetap terbuka.
	delete(u.g.Master, "UJI-T-NP")
	kode, isi := u.panggil("GET", "/kasus/"+id, admin, nil)
	if kode != http.StatusOK || !strings.Contains(isi, "Error fetching XolList") {
		t.Fatalf("%d %s", kode, isi)
	}
}
