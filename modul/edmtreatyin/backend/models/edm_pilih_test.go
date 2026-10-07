package models

// Uji rantai pilih bisnis endorsemen (edm_pilih.go): tiap aktivitas, lalu `PilihBisnisEDM` berlapis untuk
// Proporsional, NonProporsional (XOL), XOL Retro, AdjPremi (EDMType 3), dan Cancel (EDMType 4). Pembaca master
// tiruan; data UJI- dan angka karangan sederhana.

import (
	"strings"
	"testing"
	"time"
)

// pbuPembaca - tiruan `PembacaMasterEDM` yang mencatat setiap panggilan (urutan RDB-List).
type pbuPembaca struct {
	oldID, id, edm, out map[string][]MasterXOL
	mu                  map[string]string
	panggil             []string
}

func (f *pbuPembaca) catat(s string) { f.panggil = append(f.panggil, s) }
func (f *pbuPembaca) MasterMenurutOldID(x string) ([]MasterXOL, error) {
	f.catat("OLDID:" + x)
	return f.oldID[x], nil
}
func (f *pbuPembaca) MasterMenurutID(x string) ([]MasterXOL, error) {
	f.catat("ID:" + x)
	return f.id[x], nil
}
func (f *pbuPembaca) MasterEDMMenurutID(x string) ([]MasterXOL, error) {
	f.catat("EDM:" + x)
	return f.edm[x], nil
}
func (f *pbuPembaca) MasterOutMenurutID(x string) ([]MasterXOL, error) {
	f.catat("OUT:" + x)
	return f.out[x], nil
}
func (f *pbuPembaca) IDMataUang(n string) (string, error) { return f.mu[n], nil }

// pbuMaster - satu baris master: skalar + satu layer `Share` (IDR gross/net/potongan, DeductionTotalList = potongan),
// satu angsuran IDR dengan `rinci` baris InstallmentList.
func pbuMaster(id, gross, net, ded string, rinci int) MasterXOL {
	m := MasterXOL{Nilai: map[string]string{"ID": id, "Commencement": "2026-01-01", "Termination": "2026-12-31",
		"FacultativeShare": "0", "InstallmentNo": "UJI-NO"}, Daftar: map[string][]Baris{}}
	m.Daftar["Share"] = []Baris{{"Layer": "1", "LayerType": "UJI-LT"}}
	m.Daftar["Share(1).GrossPremiumList"] = []Baris{xeuMV("IDR", gross)}
	m.Daftar["Share(1).NetPremiumList"] = []Baris{xeuMV("IDR", net)}
	m.Daftar["Share(1).DeductionTotalList"] = []Baris{xeuMV("IDR", ded)}
	m.Daftar["Share(1).DeductionList"] = []Baris{xeuPot("IDR", ded)}
	m.Daftar["Installment"] = []Baris{{"Currency": "IDR"}}
	var r []Baris
	for i := 0; i < rinci; i++ {
		r = append(r, Baris{"Installment": "UJI"})
	}
	m.Daftar["Installment(1).InstallmentList"] = r
	return m
}

// pbuHalaman - kasus EDM baru (CreateEDMT): OldData = generasi lama, New hanya EDMType/Currency/ClaimType.
func pbuHalaman(edm, claim, nonProp string) *Halaman {
	h := HalamanBaru()
	h.Setel(pt+"EDMType", edm)
	h.Setel(pt+"Currency", "IDR")
	h.Setel(pt+"ClaimType", claim)
	for m, v := range map[string]string{"NoOffer": "UJI-M1", "PolicyNo": "UJI-POL-1", "IsNewPolicyNonProp": nonProp,
		"Currency": "IDR", "SOBName": "UJI-SOB", "TreatyGroupName": "UJI-GRUP-LAMA", "StartDate": "2025-01-01",
		"NetPremium": "800", "PremiOgp": "1000", "Deduction1": "200", "BalanceDueTo": "800", "Installment": "1"} {
		h.Setel(od+m, v)
	}
	h.SetelDaftar(od+"SpreadingRiskList", []Baris{{"TreatyType": "UJI-TT", "SharePercentage": "100", "ClaimPercentage": "100"}})
	h.Setel(pt+"QuotationData.BusinessCode", "UJI-BIZ")
	return h
}

func pbuSekarang() time.Time { return time.Date(2026, 10, 6, 9, 0, 0, 0, time.UTC) }

// SetValueEDM_Act pilihan pertama (NoOffer kosong, bukan Retro): RDB 5 (OLDID) lalu 6 (ID) - `[dugaan]` hasil 6
// MENGGANTI hasil 5, jadi kunci yang hanya ada di master ber-OLDID tidak pernah di-adopt.
func TestSetValueEDMPilihanPertamaHasilTerakhirMenang(t *testing.T) {
	h := pbuHalaman("1", "", "1")
	edm := pbuMaster("UJI-M1-E1", "1", "1", "1", 0)
	edm.Nilai["IsProRate"] = "true"
	f := &pbuPembaca{oldID: map[string][]MasterXOL{"UJI-M1": {edm}},
		id: map[string][]MasterXOL{"UJI-M1": {pbuMaster("UJI-M1", "1000", "800", "200", 1)}}}
	atas := HalamanBaru()
	atas.Setel(jMaster+"Sisa", "UJI-LAMA")
	if err := SetValueEDM(h, atas, f); err != nil {
		t.Fatal(err)
	}
	if strings.Join(f.panggil, ",") != "OLDID:UJI-M1,ID:UJI-M1" {
		t.Fatalf("urutan RDB %v", f.panggil)
	}
	if h.Ambil(jMaster+"ID") != "UJI-M1" || h.Ambil(jMaster+"IsProRate") != "" {
		t.Fatalf("adopt: ID %q IsProRate %q", h.Ambil(jMaster+"ID"), h.Ambil(jMaster+"IsProRate"))
	}
	if atas.Ambil(jMaster+"ID") != "UJI-M1" || atas.Ambil(jMaster+"Sisa") != "" {
		t.Fatalf("TreatyIn tingkat atas (2-4): %v", atas.Nilai)
	}
	// NP (EDMType 1): DueTo 1, tanggal master, jumlah Share ber-mata uang polis
	xeuSama(t, "PremiOgp", h.Ambil(pt+"PremiOgp"), "1000")
	xeuSama(t, "NetPremium", h.Ambil(pt+"NetPremium"), "800")
	xeuSama(t, "Deduction1", h.Ambil(pt+"Deduction1"), "200")
	xeuSama(t, "ShareValue", h.Ambil(pt+"ShareValue"), "0")
	if h.Ambil(pt+"DueTo") != "1" || h.Ambil(pt+"StartDate") != "2026-01-01" || h.Ambil(pt+"EndDate") != "2026-12-31" {
		t.Errorf("langkah 2: %q %q %q", h.Ambil(pt+"DueTo"), h.Ambil(pt+"StartDate"), h.Ambil(pt+"EndDate"))
	}
	if len(h.AmbilDaftar(DaftarXOL)) != 1 { // 8 InsertToTreatyXOLList (NB)
		t.Errorf("TreatyXOLList %v", h.AmbilDaftar(DaftarXOL))
	}
}

// Pilihan berikut (NoOffer terisi): hanya RDB 5 (OLDID); TreatyIn tingkat atas tanpa ID. XOL Retro: hanya RDB 7
// (M_TREATY_OUT) - dan 6 bila NoOffer kosong, yang hasilnya diganti 7.
func TestSetValueEDMPilihanBerikutDanRetro(t *testing.T) {
	h := pbuHalaman("1", "", "1")
	h.Setel(pt+"NoOffer", "UJI-M1-E1")
	f := &pbuPembaca{}
	atas := HalamanBaru()
	if err := SetValueEDM(h, atas, f); err != nil {
		t.Fatal(err)
	}
	if strings.Join(f.panggil, ",") != "OLDID:UJI-M1" || atas.Ambil(jMaster+"ID") != "" {
		t.Fatalf("pilihan berikut: %v, ID atas %q", f.panggil, atas.Ambil(jMaster+"ID"))
	}
	h = pbuHalaman("1", KlaimXOLRetro, "1")
	out := MasterXOL{Nilai: map[string]string{"ID": "UJI-OUT"}, Daftar: map[string][]Baris{}}
	f = &pbuPembaca{id: map[string][]MasterXOL{"UJI-M1": {pbuMaster("UJI-M1", "1", "1", "1", 0)}},
		out: map[string][]MasterXOL{"UJI-M1": {out}}}
	if err := SetValueEDM(h, HalamanBaru(), f); err != nil {
		t.Fatal(err)
	}
	if strings.Join(f.panggil, ",") != "ID:UJI-M1,OUT:UJI-M1" || h.Ambil(jMaster+"ID") != "UJI-OUT" {
		t.Fatalf("Retro: %v ID %q", f.panggil, h.Ambil(jMaster+"ID"))
	}
	if _, ada := h.Daftar[jMaster+"Share"]; ada {
		t.Error("hasil RDB 6 diganti RDB 7 - Share master lama tidak boleh ter-adopt")
	}
}

// adoptJSONObject berurutan: baris kedua menimpa kunci tingkat atas dan mengganti daftar utuh (beserta anaknya);
// kunci yang tidak disebut tetap.
func TestAdopsiMasterEDMBerurutan(t *testing.T) {
	h := HalamanBaru()
	h.Setel(jMaster+"Lain", "UJI-TETAP")
	a := pbuMaster("UJI-A", "1", "1", "1", 2)
	b := MasterXOL{Nilai: map[string]string{"ID": "UJI-B"}, Daftar: map[string][]Baris{"Installment": {{"Currency": "USD"}}}}
	for _, m := range []MasterXOL{a, b} {
		AdopsiMasterEDM(h, HalamanMaster, m)
	}
	if h.Ambil(jMaster+"ID") != "UJI-B" || h.Ambil(jMaster+"Commencement") != "2026-01-01" || h.Ambil(jMaster+"Lain") != "UJI-TETAP" {
		t.Fatalf("skalar %v", h.Nilai)
	}
	if d := h.AmbilDaftar(jMaster + "Installment"); len(d) != 1 || d[0]["Currency"] != "USD" {
		t.Fatalf("Installment %v", d)
	}
	if _, ada := h.Daftar[JalurAnak(jMaster+"Installment", 1, "InstallmentList")]; ada {
		t.Fatal("anak daftar lama ikut terganti")
	}
	if len(h.AmbilDaftar(jMaster+"Share")) != 1 {
		t.Fatal("daftar yang tidak disebut baris kedua tetap")
	}
}

// AdjPremi (EDMType 3): jumlah dari ValueDifference.Share; IsFacRetro bila FacultativeShare > 0; XOL dari
// InsertToTreatyXOLListEDM (ActualValue.Share disalin dari Share).
func TestInputPolicyTreatyEDMDetailNPAdjPremi(t *testing.T) {
	h := pbuHalaman("3", "", "1")
	m := pbuMaster("UJI-M1-E1", "1000", "800", "200", 0)
	m.Nilai["FacultativeShare"] = "5"
	m.Daftar["ValueDifference.Share"] = []Baris{{}, {}}
	m.Daftar["ValueDifference.Share(1).GrossPremiumList"] = []Baris{xeuMV("IDR", "60"), xeuMV("USD", "9")}
	m.Daftar["ValueDifference.Share(2).GrossPremiumList"] = []Baris{xeuMV("IDR", "40")}
	m.Daftar["ValueDifference.Share(1).NetPremiumList"] = []Baris{xeuMV("IDR", "50")}
	m.Daftar["ValueDifference.Share(2).DeductionTotalList"] = []Baris{xeuMV("IDR", "10")}
	AdopsiMasterEDM(h, HalamanMaster, m)
	h.SetelDaftar(DaftarAngsuran, []Baris{{"Premium": "UJI-HAPUS"}})
	if err := InputPolicyTreatyEDMDetailNPAdjPremi(h, &pbuPembaca{mu: map[string]string{"IDR": "UJI-ID-IDR"}}); err != nil {
		t.Fatal(err)
	}
	xeuSama(t, "PremiOgp", h.Ambil(pt+"PremiOgp"), "100")
	xeuSama(t, "NetPremium", h.Ambil(pt+"NetPremium"), "50")
	xeuSama(t, "BalanceDueTo", h.Ambil(pt+"BalanceDueTo"), "50")
	xeuSama(t, "Deduction1", h.Ambil(pt+"Deduction1"), "10")
	xeuSama(t, "Deduction2", h.Ambil(pt+"Deduction2"), "0")
	if h.Ambil("OfferFacIn.IsFacRetro") != "1" || len(h.AmbilDaftar(DaftarAngsuran)) != 0 {
		t.Errorf("IsFacRetro %q, ListInstallment %v", h.Ambil("OfferFacIn.IsFacRetro"), h.AmbilDaftar(DaftarAngsuran))
	}
	x := h.AmbilDaftar(DaftarXOL)
	if len(x) != 1 || x[0]["IDCurrency"] != "UJI-ID-IDR" || len(h.AmbilDaftar(jMaster+"ActualValue.Share")) != 1 {
		t.Fatalf("XOL EDM %v", x)
	}
	xeuSama(t, "XOL GrossPremi (ActualValue = Share)", x[0]["GrossPremi"], "1000")
}

// SetValueOldTax: hanya baris ber-PPNValue kosong (induk dan layer, terpisah).
func TestSetValueOldTax(t *testing.T) {
	h := HalamanBaru()
	h.SetelDaftar(od+"TreatyXOLList", []Baris{{"NetPremi": "80"}, {"NetPremi": "90", "PPNValue": "1"}})
	h.SetelDaftar(JalurAnak(od+"TreatyXOLList", 2, AnakLayerXOL), []Baris{{"NetPremi": "70"}})
	SetValueOldTax(h)
	d := h.AmbilDaftar(od + "TreatyXOLList")
	if d[0]["PPHValue"] != "0" || d[0]["PPNValue"] != "0" || d[0]["NetPremiAfterTax"] != "80" || d[0]["NetPremiAfterPPN"] != "80" {
		t.Errorf("induk 1 %v", d[0])
	}
	if d[1]["PPNValue"] != "1" || d[1]["NetPremiAfterTax"] != "" {
		t.Errorf("induk 2 (PPN terisi) %v", d[1])
	}
	if l := h.AmbilDaftar(JalurAnak(od+"TreatyXOLList", 2, AnakLayerXOL))[0]; l["NetPremiAfterPPH"] != "70" {
		t.Errorf("layer %v", l)
	}
}

// CopyGeneralDataEDM: IDCurrency = NAMA mata uang; NoOffer = TreatyIn.ID; TreatyGroupName akhir = OldData
// (Limits(1).TreatyGroupList(1) ditimpa); langkah 2 hanya bila OldData bukan NonProp; BalanceBefore* kosong/0 ->
// OldData.NetPremium.
func TestCopyGeneralDataEDM(t *testing.T) {
	for _, nonProp := range []string{"1", "0"} {
		h := pbuHalaman("1", "", nonProp)
		h.Setel(od+"BalanceBeforePPH", "0")
		h.Setel(od+"BalanceBeforeTax", "777")
		h.Setel(jMaster+"ID", "UJI-M1-E1")
		h.SetelDaftar(JalurAnak(jMaster+"Limits", 1, "TreatyGroupList"), []Baris{{"TreatyGroup": "UJI-GRUP-MASTER"}})
		CopyGeneralDataEDM(h, pbuSekarang())
		if h.Ambil(pt+"IDCurrency") != "IDR" || h.Ambil(pt+"NoOffer") != "UJI-M1-E1" || h.Ambil(pt+"TreatyGroupName") != "UJI-GRUP-LAMA" {
			t.Errorf("langkah 1: %q %q %q", h.Ambil(pt+"IDCurrency"), h.Ambil(pt+"NoOffer"), h.Ambil(pt+"TreatyGroupName"))
		}
		if !strings.HasPrefix(h.Ambil(pt+"StatementDate"), "2026-10-06") || h.Ambil(pt+"StartDate") != "2025-01-01" {
			t.Errorf("StatementDate %q StartDate %q", h.Ambil(pt+"StatementDate"), h.Ambil(pt+"StartDate"))
		}
		salin := h.Ambil(pt+"PremiOgp") == "1000"
		if salin != (nonProp == "0") {
			t.Errorf("IsNewPolicyNonProp %s: langkah 2 jalan = %v", nonProp, salin)
		}
		if nonProp == "0" && (h.Ambil(od+"BalanceBeforePPH") != "800" || h.Ambil(od+"BalanceBeforeTax") != "777") {
			t.Errorf("BalanceBefore %q %q", h.Ambil(od+"BalanceBeforePPH"), h.Ambil(od+"BalanceBeforeTax"))
		}
	}
}

// SetEDMTCancel Proporsional: uang "0", daftar disalin dari OldData beserta anaknya lalu dinolkan; GrossPremium tidak.
func TestSetEDMTCancelProporsional(t *testing.T) {
	h := pbuHalaman("4", "", "0")
	h.Setel(pt+"GrossPremium", "UJI-TETAP")
	h.Setel(pt+"PremiOgp", "1000")
	h.SetelDaftar(od+"ListInstallment", []Baris{{"Premium": "800", "PaymentTotal": "800", "InstallmentNo": "1"}})
	h.SetelDaftar(JalurAnak(od+"ListInstallment", 1, AnakRinciAngsuran), []Baris{{"Premium": "800"}})
	SetEDMTCancel(h)
	if h.Ambil(pt+"PremiOgp") != "0" || h.Ambil(pt+"OutstandingClaim") != "0" || h.Ambil(pt+"GrossPremium") != "UJI-TETAP" {
		t.Errorf("1.1: %q %q %q", h.Ambil(pt+"PremiOgp"), h.Ambil(pt+"OutstandingClaim"), h.Ambil(pt+"GrossPremium"))
	}
	a := h.AmbilDaftar(DaftarAngsuran)
	if len(a) != 1 || a[0]["Premium"] != "0" || a[0]["PaymentTotal"] != "0" || a[0]["InstallmentNo"] != "1" {
		t.Errorf("ListInstallment %v", a)
	}
	if r := h.AmbilDaftar(JalurAnak(DaftarAngsuran, 1, AnakRinciAngsuran)); len(r) != 1 || r[0]["Premium"] != "800" {
		t.Errorf("rincian ikut tersalin (tidak dinolkan 1.3): %v", r)
	}
	if s := h.AmbilDaftar(DaftarSpreading); len(s) != 1 || s[0]["SharePercentage"] != "0" || s[0]["ClaimPercentage"] != "100" {
		t.Errorf("SpreadingRiskList %v", s)
	}
	if h.AmbilDaftar(od + "ListInstallment")[0]["Premium"] != "800" {
		t.Error("OldData tidak boleh ikut berubah (salinan)")
	}
	if h.Ambil(pt+"Installment") != "1" {
		t.Errorf("Installment %q", h.Ambil(pt+"Installment"))
	}
}

// SetEDMTCancel NonProp: induk/layer XOL dinolkan, Currency/IDCurrency/DueTo dari OldData per subskrip; layer
// Currency = Currency induk, IDCurrency tetap, Layer* dari OldData; BrokerageFeeSebenarnya tidak dinolkan.
func TestSetEDMTCancelNonProp(t *testing.T) {
	h := pbuHalaman("4", "", "1")
	h.Setel(pt+"IsNewPolicyNonProp", "1")
	h.SetelDaftar(DaftarXOL, []Baris{{"Currency": "UJI-BARU", "GrossPremi": "1200", "BrokerageFeeSebenarnya": "9"}})
	h.SetelDaftar(JalurAnak(DaftarXOL, 1, AnakLayerXOL), []Baris{{"IDCurrency": "UJI-ID-L", "NetPremi": "5", "Layer": "UJI-BARU"}})
	h.SetelDaftar(od+"TreatyXOLList", []Baris{{"Currency": "IDR", "IDCurrency": "UJI-ID-IDR", "DueTo": "UJI-DT"}})
	h.SetelDaftar(JalurAnak(od+"TreatyXOLList", 1, AnakLayerXOL), []Baris{{"Layer": "1", "LayerType": "UJI-LT"}})
	SetEDMTCancel(h)
	x := h.AmbilDaftar(DaftarXOL)[0]
	if x["Currency"] != "IDR" || x["IDCurrency"] != "UJI-ID-IDR" || x["DueTo"] != "UJI-DT" || x["GrossPremi"] != "0" ||
		x["NetPremiAfterTax"] != "0" || x["BrokerageFeeSebenarnya"] != "9" {
		t.Errorf("induk %v", x)
	}
	l := h.AmbilDaftar(JalurAnak(DaftarXOL, 1, AnakLayerXOL))[0]
	if l["Currency"] != "IDR" || l["IDCurrency"] != "UJI-ID-L" || l["DueTo"] != "UJI-DT" || l["Layer"] != "1" ||
		l["LayerType"] != "UJI-LT" || l["NetPremi"] != "0" {
		t.Errorf("layer %v", l)
	}
	if h.Ambil(pt+"PremiOgp") != "" {
		t.Error("blok Proporsional tidak boleh jalan")
	}
}

// FillMasterInstallment + DT SetInstallmentValue: M_TREATY_IN_EDM ber-ID NoOffer ke TreatyIn tingkat atas;
// Installment = pyWorkPage.TreatyIn.InstallmentNo; ListInstallment/InstallmentList per subskrip dari
// ValueDifference.Installment; TreatyIn tingkat atas dibuang (5). Retro: M_TREATY_OUT menggantikan.
func TestFillMasterInstallment(t *testing.T) {
	h := pbuHalaman("1", "", "1")
	h.Setel(pt+"NoOffer", "UJI-M1-E1")
	h.Setel(jMaster+"InstallmentNo", "2")
	m := MasterXOL{Nilai: map[string]string{"ID": "UJI-M1-E1"}, Daftar: map[string][]Baris{
		"ValueDifference.Installment": {{"Currency": "IDR", "AmountTotal": "300"}},
		"ValueDifference.Installment(1).InstallmentList": {
			{"Currency": "IDR", "Amount": "100", "DueDate": "2026-11-01", "Installment": "1", "InstallmentPct": "33.3333", "PaymentDate": "2026-11-05"},
			{"Currency": "IDR", "Amount": "200", "DueDate": "2027-01-01", "Installment": "2", "InstallmentPct": "66.6667"}},
	}}
	f := &pbuPembaca{edm: map[string][]MasterXOL{"UJI-M1-E1": {m}}}
	atas := HalamanBaru()
	atas.Setel(jMaster+"Sisa", "UJI-LAMA")
	h.SetelDaftar(DaftarAngsuran, []Baris{{"InstallmentNo": "UJI-TETAP"}})
	if err := FillMasterInstallment(h, atas, f); err != nil {
		t.Fatal(err)
	}
	if strings.Join(f.panggil, ",") != "EDM:UJI-M1-E1" || len(atas.Nilai) != 0 || len(atas.Daftar) != 0 {
		t.Fatalf("panggilan %v, atas %v %v", f.panggil, atas.Nilai, atas.Daftar)
	}
	if h.Ambil(pt+"Installment") != "2" || h.Ambil(sd+"TotalPremium") != "300" {
		t.Errorf("Installment %q TotalPremium %q", h.Ambil(pt+"Installment"), h.Ambil(sd+"TotalPremium"))
	}
	a := h.AmbilDaftar(DaftarAngsuran)
	if len(a) != 1 || a[0]["Currency"] != "IDR" || a[0]["Premium"] != "300" || a[0]["InstallmentNo"] != "UJI-TETAP" {
		t.Errorf("ListInstallment %v", a)
	}
	r := h.AmbilDaftar(JalurAnak(DaftarAngsuran, 1, AnakRinciAngsuran))
	if len(r) != 2 || r[0]["Premium"] != "100" || r[0]["InstallmentNo"] != "1" || r[0]["InstallmentPercentage"] != "33.3333" ||
		r[0]["DueDate"] != "2026-11-01" || r[0]["PaymentDate"] != "2026-11-05" || r[1]["Premium"] != "200" {
		t.Errorf("InstallmentList %v", r)
	}
	h = pbuHalaman("1", KlaimXOLRetro, "1")
	h.Setel(pt+"NoOffer", "UJI-X")
	f = &pbuPembaca{}
	if err := FillMasterInstallment(h, HalamanBaru(), f); err != nil {
		t.Fatal(err)
	}
	if strings.Join(f.panggil, ",") != "EDM:UJI-X,OUT:UJI-M1" {
		t.Fatalf("Retro %v", f.panggil)
	}
}

// FillSpreading: baris 1 = salinan OldData(1), PremiumSpreaded = selisih XOL(1).NetPremi, total; pesan bila
// spreading lama kosong.
func TestFillSpreading(t *testing.T) {
	h := pbuHalaman("1", "", "1")
	h.SetelDaftar(DaftarSpreading, []Baris{{"TreatyType": "UJI-GANTI", "Lain": "UJI-HILANG"}, {"SharePercentage": "5"}})
	h.SetelDaftar(DaftarSelisihXOL, []Baris{{"NetPremi": "200"}})
	if err := FillSpreading(h); err != nil {
		t.Fatal(err)
	}
	d := h.AmbilDaftar(DaftarSpreading)
	if len(d) != 2 || d[0]["TreatyType"] != "UJI-TT" || d[0]["Lain"] != "" || d[0]["PremiumSpreaded"] != "200" {
		t.Fatalf("spreading %v", d)
	}
	xeuSama(t, "TotalSharePercentagePremium", h.Ambil(pt+"TotalSharePercentagePremium"), "105")
	xeuSama(t, "TotalPremium", h.Ambil(pt+"TotalPremium"), "200")
	xeuSama(t, "TotalClaim", h.Ambil(pt+"TotalClaim"), "0")
	if adaPesan(h) {
		t.Errorf("pesan %v", h.SemuaPesan())
	}
	h = HalamanBaru()
	if err := FillSpreading(h); err != nil {
		t.Fatal(err)
	}
	if p := h.SemuaPesan(); len(p) != 1 || p[0] != PesanSpreadingAsalKosong {
		t.Errorf("pesan %v", p)
	}
}

// TreatyRealizationCheckXOLListEDM: OldData XOL sudah ada (bukan FlagRetroTreaty) -> keluar tanpa RDB; kosong ->
// SetTreatyIn_Act (ID = OldData.NoOffer) lalu XOL lama dari master lama; FlagRetroTreaty -> selalu dibangun ulang.
func TestTreatyRealizationCheckXOLListEDM(t *testing.T) {
	h := pbuHalaman("1", "", "1")
	h.SetelDaftar(od+"TreatyXOLList", []Baris{{"GrossPremi": "UJI-TETAP"}})
	h.TambahPesan("", "UJI-PESAN")
	f := &pbuPembaca{id: map[string][]MasterXOL{"UJI-M1": {pbuMaster("UJI-M1", "1000", "800", "200", 1)}}}
	if err := TreatyRealizationCheckXOLListEDM(h, HalamanBaru(), f); err != nil {
		t.Fatal(err)
	}
	if len(f.panggil) != 0 || h.AmbilDaftar(od + "TreatyXOLList")[0]["GrossPremi"] != "UJI-TETAP" || adaPesan(h) {
		t.Fatalf("keluar langkah 3: %v %v", f.panggil, h.SemuaPesan())
	}
	h.Setel(pt+"FlagRetroTreaty", "true")
	atas := HalamanBaru()
	if err := TreatyRealizationCheckXOLListEDM(h, atas, f); err != nil {
		t.Fatal(err)
	}
	if strings.Join(f.panggil, ",") != "ID:UJI-M1" || atas.Ambil(jMaster+"ID") != "UJI-M1" {
		t.Fatalf("SetTreatyIn_Act %v", f.panggil)
	}
	xeuSama(t, "OldData XOL dibangun ulang", h.AmbilDaftar(od + "TreatyXOLList")[0]["GrossPremi"], "200")
}

// ---------------------------------------------------------------- rantai EDMChooseBusiness_Act

// pbuRantai menyiapkan master: OLDID UJI-M1 -> master EDM (Share IDR 1200/1000/200, 2 termin); ID UJI-M1 -> master
// lama (1000/800/200); EDM ID UJI-M1-E1 -> ValueDifference.Installment.
func pbuRantai() *pbuPembaca {
	baru := pbuMaster("UJI-M1-E1", "1200", "1000", "200", 2)
	baru.Daftar["ValueDifference.Share"] = []Baris{{}}
	baru.Daftar["ValueDifference.Share(1).GrossPremiumList"] = []Baris{xeuMV("IDR", "200")}
	baru.Daftar["ValueDifference.Share(1).NetPremiumList"] = []Baris{xeuMV("IDR", "200")}
	baru.Daftar["ActualValue.Share"] = []Baris{{"Layer": "1"}}
	baru.Daftar["ActualValue.Share(1).GrossPremiumList"] = []Baris{xeuMV("IDR", "1200")}
	baru.Daftar["ActualValue.Share(1).NetPremiumList"] = []Baris{xeuMV("IDR", "1000")}
	baru.Daftar["ActualValue.Share(1).DeductionList"] = []Baris{xeuPot("IDR", "200")}
	vd := MasterXOL{Nilai: map[string]string{"ID": "UJI-M1-E1"}, Daftar: map[string][]Baris{
		"ValueDifference.Installment":                    {{"Currency": "IDR", "AmountTotal": "200"}},
		"ValueDifference.Installment(1).InstallmentList": {{"Amount": "200", "Installment": "1"}}}}
	return &pbuPembaca{
		oldID: map[string][]MasterXOL{"UJI-M1": {baru}},
		id:    map[string][]MasterXOL{"UJI-M1": {pbuMaster("UJI-M1", "1000", "800", "200", 1)}},
		edm:   map[string][]MasterXOL{"UJI-M1-E1": {vd}},
		mu:    map[string]string{"IDR": "UJI-ID-IDR"},
	}
}

// NonProp (EDMType 1, pilihan berikut): XOL baru dari master EDM, XOL lama dari master lama, selisih net 1000-800
// = 200; angsuran EDMT 2 termin x 100 (DueDate = StatementDate); spreading baris 1 dari OldData, premi 200.
func TestPilihBisnisEDMNonProp(t *testing.T) {
	h := pbuHalaman("1", "", "1")
	h.Setel(pt+"NoOffer", "UJI-M1-E1")
	f := pbuRantai()
	if err := PilihBisnisEDM(h, f, pbuSekarang()); err != nil {
		t.Fatal(err)
	}
	if strings.Join(f.panggil, ",") != "OLDID:UJI-M1,ID:UJI-M1,EDM:UJI-M1-E1" {
		t.Fatalf("RDB %v", f.panggil)
	}
	if h.Ambil(pt+"NoOffer") != "UJI-M1-E1" || h.Ambil(pt+"IDCurrency") != "IDR" || h.Ambil(pt+"PremiOgp") != "1200" {
		t.Errorf("CopyGeneralDataEDM / NP: %q %q %q", h.Ambil(pt+"NoOffer"), h.Ambil(pt+"IDCurrency"), h.Ambil(pt+"PremiOgp"))
	}
	xeuSama(t, "XOL lama GrossPremi", h.AmbilDaftar(od + "TreatyXOLList")[0]["GrossPremi"], "1000")
	d := h.AmbilDaftar(DaftarSelisihXOL)
	if len(d) != 1 || d[0]["IDCurrency"] != "UJI-ID-IDR" {
		t.Fatalf("selisih %v", d)
	}
	xeuBaris(t, "selisih induk", d[0], map[string]string{"GrossPremi": "200", "NetPremi": "200", "Deduction": "0"})
	if h.Ambil(pt+"Installment") != "2" {
		t.Errorf("Installment %q (jumlah InstallmentList master)", h.Ambil(pt+"Installment"))
	}
	a := h.AmbilDaftar(DaftarAngsuran)
	r := h.AmbilDaftar(JalurAnak(DaftarAngsuran, 1, AnakRinciAngsuran))
	if len(a) != 1 || len(r) != 2 || r[0]["DueDate"] != "2026-10-06" {
		t.Fatalf("angsuran %v / %v", a, r)
	}
	xeuSama(t, "termin 1", r[0]["Premium"], "100")
	sp := h.AmbilDaftar(DaftarSpreading)
	if len(sp) != 1 || sp[0]["TreatyType"] != "UJI-TT" || sp[0]["PremiumSpreaded"] != "200" {
		t.Errorf("spreading %v", sp)
	}
	if h.Ambil("OfferFacIn.QuotationData.BusinessCode") != "UJI-BIZ" || h.Ambil(pt+"QuotationData.OldPolicyNo") != "UJI-POL-1" ||
		h.Ambil("OfferFacIn.QuotationData.OldPolicyNo") != "" {
		t.Errorf("langkah 10: %v", h.Nilai)
	}
}

// Proporsional (OldData.IsNewPolicyNonProp 0): uang New = OldData (CopyGeneralDataEDM 2 menimpa NP), XOL tetap
// dibangun (SetValueEDM 9 tanpa syarat jenis), pilihan pertama: master lama = master baru -> selisih 0.
func TestPilihBisnisEDMProporsional(t *testing.T) {
	h := pbuHalaman("1", "", "0")
	f := pbuRantai()
	if err := PilihBisnisEDM(h, f, pbuSekarang()); err != nil {
		t.Fatal(err)
	}
	if strings.Join(f.panggil, ",") != "OLDID:UJI-M1,ID:UJI-M1,ID:UJI-M1,EDM:UJI-M1" {
		t.Fatalf("RDB %v", f.panggil)
	}
	if h.Ambil(pt+"PremiOgp") != "1000" || h.Ambil(pt+"NetPremium") != "800" || h.Ambil(pt+"NoOffer") != "UJI-M1" {
		t.Errorf("uang = OldData: %q %q %q", h.Ambil(pt+"PremiOgp"), h.Ambil(pt+"NetPremium"), h.Ambil(pt+"NoOffer"))
	}
	xeuBaris(t, "selisih", h.AmbilDaftar(DaftarSelisihXOL)[0], map[string]string{"GrossPremi": "0", "NetPremi": "0"})
	if h.Ambil(pt+"Installment") != "1" {
		t.Errorf("Installment %q", h.Ambil(pt+"Installment"))
	}
}

// XOL Retro, pilihan berikut (PolicyTreatyIn.SOBName sudah terisi CopyGeneralDataEDM pilihan sebelumnya): master
// M_TREATY_OUT; XOL baru InsertToTreatyOutXOLList (net dua kali), XOL lama InsertToTreatyOutXOLListEDMOldData
// (sumber master BARU, potongan SpreadingTONP), dibangun walau FlagRetroTreaty kosong karena OldData.TreatyXOLList
// kosong. Pilihan PERTAMA: SOBName baru masih kosong saat langkah 3 -> semua baris ReinsName terbuang -> XOL nol
// (ditiru apa adanya; satu-satunya penulis PolicyTreatyIn.SOBName di korpus = CopyGeneralDataEDM_act).
func TestPilihBisnisEDMRetro(t *testing.T) {
	pertama := pbuHalaman("1", KlaimXOLRetro, "1")
	h := pbuHalaman("1", KlaimXOLRetro, "1")
	h.Setel(pt+"NoOffer", "UJI-OUT-1")
	h.Setel(pt+"SOBName", "UJI-SOB")
	out := MasterXOL{Nilai: map[string]string{"ID": "UJI-OUT-1"}, Daftar: map[string][]Baris{
		"Installment":                       {{"Currency": "IDR"}},
		"Installment(1).InstallmentList":    {{}},
		"ShareReins":                        {{}},
		"ShareReins(1).ReinsuranceListTONP": {{"ReinsName": "UJI-SOB", "Layer": "1"}},
		"ShareReins(1).ReinsuranceListTONP(1).GrossPremiumList": {xeuMV("IDR", "100")},
		"ShareReins(1).ReinsuranceListTONP(1).NetPremiumList":   {xeuMV("IDR", "80")},
		"ShareReins(1).ReinsuranceListTONP(1).DeductionList":    {xeuPot("IDR", "20")},
		"SpreadingTONP":                                                      {{}},
		"SpreadingTONP(1).ReinsuranceListTONP":                               {{}},
		"SpreadingTONP(1).ReinsuranceListTONP(1).LayerList":                  {{}},
		"SpreadingTONP(1).ReinsuranceListTONP(1).LayerList(1).DeductionList": {xeuPot("IDR", "5")},
	}}
	f := &pbuPembaca{out: map[string][]MasterXOL{"UJI-M1": {out}}, mu: map[string]string{"IDR": "UJI-ID-IDR"}}
	if err := PilihBisnisEDM(h, f, pbuSekarang()); err != nil {
		t.Fatal(err)
	}
	if strings.Join(f.panggil, ",") != "OUT:UJI-M1,ID:UJI-M1,EDM:UJI-OUT-1,OUT:UJI-M1" {
		t.Fatalf("RDB %v", f.panggil)
	}
	xeuBaris(t, "XOL baru", h.AmbilDaftar(DaftarXOL)[0], map[string]string{"NetPremi": "160", "Deduction": "20"})
	xeuBaris(t, "XOL lama", h.AmbilDaftar(od + "TreatyXOLList")[0], map[string]string{"NetPremi": "80", "Deduction": "5"})
	// selisih PER LAYER lalu dijumlah (langkah 2): net layer 80-80 = 0 (gandaan net hanya di induk TreatyXOLList),
	// potongan 20-5 = 15
	xeuBaris(t, "selisih", h.AmbilDaftar(DaftarSelisihXOL)[0], map[string]string{"GrossPremi": "0", "NetPremi": "0", "Deduction": "15"})
	if h.Ambil(pt+"Installment") != "1" {
		t.Errorf("Installment %q", h.Ambil(pt+"Installment"))
	}
	f = &pbuPembaca{out: map[string][]MasterXOL{"UJI-M1": {out}}}
	if err := PilihBisnisEDM(pertama, f, pbuSekarang()); err != nil {
		t.Fatal(err)
	}
	xeuBaris(t, "pilihan pertama XOL baru", pertama.AmbilDaftar(DaftarXOL)[0], map[string]string{"NetPremi": "0", "GrossPremi": "0"})
	if n := len(pertama.AmbilDaftar(xeTONP(1))); n != 0 || pertama.Ambil(pt+"SOBName") != "UJI-SOB" {
		t.Errorf("pilihan pertama: %d baris TONP tersisa, SOBName %q", n, pertama.Ambil(pt+"SOBName"))
	}
}

// AdjPremi (EDMType 3, FlagPPH): jumlah dari ValueDifference.Share, XOL baru dari ActualValue.Share; pajak lama
// diisi SetValueOldTax; selisih pajak dihitung ulang dari potongan selisih (1.2.4).
func TestPilihBisnisEDMAdjPremi(t *testing.T) {
	h := pbuHalaman("3", "", "1")
	h.Setel(pt+"NoOffer", "UJI-M1-E1")
	h.Setel(pt+"FlagPPH", "true")
	h.Setel(od+"FlagPPH", "true")
	f := pbuRantai()
	f.id["UJI-M1"][0].Daftar["Share(1).DeductionList"] = []Baris{xeuPot("IDR", "100")}
	if err := PilihBisnisEDM(h, f, pbuSekarang()); err != nil {
		t.Fatal(err)
	}
	xeuSama(t, "PremiOgp (ValueDifference)", h.Ambil(pt+"PremiOgp"), "200")
	if len(h.AmbilDaftar(od+"TreatyXOLList")) != 1 {
		t.Fatal("XOL lama")
	}
	l := h.AmbilDaftar(JalurAnak(DaftarSelisihXOL, 1, AnakLayerXOL))[0]
	// potongan selisih 200-100 = 100 -> BFS 100, PPH 2, PPN 2.2; net 1000-800 = 200 -> AfterTax 204.2
	xeuBaris(t, "selisih layer", l, map[string]string{"Deduction": "100", "BrokerageFeeSebenarnya": "100", "PPHValue": "2",
		"PPNValue": "2.2", "NetPremiAfterTax": "204.2"})
}

// Cancel (EDMType 4) NonProp: XOL baru dinolkan (SetEDMTCancel), selisih mentah negatif (1.2.5): net 0-800 = -800.
func TestPilihBisnisEDMCancel(t *testing.T) {
	h := pbuHalaman("4", "", "1")
	h.Setel(pt+"NoOffer", "UJI-M1-E1")
	f := pbuRantai()
	if err := PilihBisnisEDM(h, f, pbuSekarang()); err != nil {
		t.Fatal(err)
	}
	if h.AmbilDaftar(DaftarXOL)[0]["NetPremi"] != "0" {
		t.Fatalf("SetEDMTCancel: %v", h.AmbilDaftar(DaftarXOL))
	}
	d := h.AmbilDaftar(DaftarSelisihXOL)[0]
	xeuBaris(t, "selisih batal", d, map[string]string{"NetPremi": "-800", "GrossPremi": "-1000", "Deduction": "-200"})
	l := h.AmbilDaftar(JalurAnak(DaftarSelisihXOL, 1, AnakLayerXOL))[0]
	if l["DueTo"] != "DUE TO YOU" {
		t.Errorf("DueTo %q", l["DueTo"])
	}
	r := h.AmbilDaftar(JalurAnak(DaftarAngsuran, 1, AnakRinciAngsuran))
	if len(r) != 2 {
		t.Fatalf("angsuran %v", r)
	}
	xeuSama(t, "termin batal", r[0]["Premium"], "-400")
}

// Keputusan work owner 07-10-2026 ("CURRENT PREMIUM ... GA ADA ANGKANYA"): With Tax / Type Tax generasi baru yang
// BELUM diisi diwarisi dari generasi lama, sehingga grid XOL Current Premium (InsertToTreatyXOLListEDM 3.3.7 /
// 3.3.8.5) berpajak tanpa admin harus mencentang With Tax sebelum Choose Business. "false" (dilepas admin) dihormati.
func TestPilihBisnisEDMNonPropWarisPajakLama(t *testing.T) {
	h := pbuHalaman("1", "", "1")
	h.Setel(pt+"NoOffer", "UJI-M1-E1")
	h.Setel(od+"FlagPPH", "true")
	h.Setel(od+"TypeTax", "Exclusive")
	if err := PilihBisnisEDM(h, pbuRantai(), pbuSekarang()); err != nil {
		t.Fatal(err)
	}
	if h.Ambil(pt+"FlagPPH") != "true" || h.Ambil(pt+"TypeTax") != "Exclusive" {
		t.Fatalf("With Tax / Type Tax tidak diwarisi: %q %q", h.Ambil(pt+"FlagPPH"), h.Ambil(pt+"TypeTax"))
	}
	lapis := h.AmbilDaftar(JalurAnak(DaftarXOL, 1, AnakLayerXOL))
	if len(lapis) == 0 || lapis[0]["PPHValue"] == "" || lapis[0]["NetPremiAfterTax"] == "" {
		t.Fatalf("layer Current Premium tanpa pajak: %+v", lapis)
	}
	if b := h.AmbilDaftar(DaftarXOL)[0]; b["PPNValue"] == "" {
		t.Fatalf("induk Current Premium tanpa PPN: %+v", b)
	}

	h = pbuHalaman("1", "", "1")
	h.Setel(pt+"NoOffer", "UJI-M1-E1")
	h.Setel(od+"FlagPPH", "true")
	h.Setel(pt+"FlagPPH", "false") // admin melepas centang
	if err := PilihBisnisEDM(h, pbuRantai(), pbuSekarang()); err != nil {
		t.Fatal(err)
	}
	if h.Ambil(pt+"FlagPPH") != "false" || h.AmbilDaftar(JalurAnak(DaftarXOL, 1, AnakLayerXOL))[0]["PPHValue"] != "" {
		t.Fatal("With Tax yang dilepas admin tidak boleh ditimpa warisan")
	}
}

// Keputusan WO 07-10-2026: With Tax dicentang SESUDAH Choose Business (NonProp baru) -> rantai dijalankan lagi dengan
// master yang sama; pajak Current Premium terisi, isian Installment / StatementDate / MarketingOfficer tetap.
func TestHitungUlangPajakNonPropEDMSesudahChooseBusiness(t *testing.T) {
	h := pbuHalaman("1", "", "1")
	h.Setel(pt+"NoOffer", "UJI-M1-E1")
	h.Setel(pt+"FlagPPH", "false") // admin melepas centang sebelum Choose Business
	if err := PilihBisnisEDM(h, pbuRantai(), pbuSekarang()); err != nil {
		t.Fatal(err)
	}
	if h.AmbilDaftar(JalurAnak(DaftarXOL, 1, AnakLayerXOL))[0]["PPHValue"] != "" {
		t.Fatal("prasyarat: tanpa With Tax layer belum berpajak")
	}
	h.Setel(pt+"Installment", "3") // sel S12 diubah admin
	if err := FillPaymentInstallmentEDMT(h, "3"); err != nil {
		t.Fatal(err)
	}
	h.Setel(pt+"MarketingOfficer", "UJI-MO-BARU")
	statement := h.Ambil(pt + "StatementDate")
	h.Setel(pt+"FlagPPH", "true")
	h.Setel(pt+"TypeTax", TypeTaxInclusive)

	f := pbuRantai()
	if err := HitungUlangPajakNonPropEDM(h, f, pbuSekarang().AddDate(0, 0, 3)); err != nil {
		t.Fatal(err)
	}
	// jalur kedua (NoOffer = master EDM): SetValueEDM hanya langkah 5, master EDM tetap UJI-M1-E1
	if strings.Join(f.panggil, ",") != "OLDID:UJI-M1,EDM:UJI-M1-E1" {
		t.Fatalf("jalur master berubah: %v", f.panggil)
	}
	lapis := h.AmbilDaftar(JalurAnak(DaftarXOL, 1, AnakLayerXOL))
	if lapis[0]["PPHValue"] == "" || lapis[0]["NetPremiAfterTax"] == "" || h.AmbilDaftar(DaftarXOL)[0]["PPNValue"] == "" {
		t.Fatalf("Current Premium tetap tanpa pajak: %+v / %+v", h.AmbilDaftar(DaftarXOL)[0], lapis)
	}
	if h.Ambil(pt+"Installment") != "3" || h.Ambil(pt+"StatementDate") != statement || h.Ambil(pt+"MarketingOfficer") != "UJI-MO-BARU" {
		t.Fatalf("isian ditimpa: Installment %q StatementDate %q MO %q", h.Ambil(pt+"Installment"),
			h.Ambil(pt+"StatementDate"), h.Ambil(pt+"MarketingOfficer"))
	}
	if n := len(h.AmbilDaftar(JalurAnak(DaftarAngsuran, 1, "InstallmentList"))); n != 3 {
		t.Fatalf("angsuran tidak disusun ulang dari Installment admin: %d termin", n)
	}
	if b := h.AmbilDaftar(DaftarAngsuran)[0]; b["PremiumAfterTax"] == "" {
		t.Fatalf("angsuran tanpa premi sesudah pajak: %+v", b)
	}
}

// Pilihan terakhir lewat jalur pertama (NoOffer kosong -> master lama ber-ID): hitung ulang tetap jalur itu.
func TestHitungUlangPajakNonPropEDMJalurPertamaTetap(t *testing.T) {
	h := pbuHalaman("1", "", "1")
	if err := PilihBisnisEDM(h, pbuRantai(), pbuSekarang()); err != nil {
		t.Fatal(err)
	}
	if h.Ambil(pt+"NoOffer") != h.Ambil(od+"NoOffer") {
		t.Fatalf("prasyarat: NoOffer jalur pertama %q", h.Ambil(pt+"NoOffer"))
	}
	f := pbuRantai()
	if err := HitungUlangPajakNonPropEDM(h, f, pbuSekarang()); err != nil {
		t.Fatal(err)
	}
	// langkah 6 (ID) tetap jalan dan TreatyIn.ID tetap master lama UJI-M1
	if strings.Join(f.panggil, ",") != "OLDID:UJI-M1,ID:UJI-M1,EDM:UJI-M1" {
		t.Fatalf("hitung ulang pajak maju ke master EDM: %v", f.panggil)
	}
}

func TestHitungUlangPajakNonPropEDMTanpaChooseBusiness(t *testing.T) {
	for _, h := range []*Halaman{pbuHalaman("1", "", "1"), pbuHalaman("1", "", "0")} {
		h.Setel(pt+"IsNewPolicyNonProp", h.Ambil(od+"IsNewPolicyNonProp"))
		h.SetelDaftar(DaftarXOL, nil)
		f := pbuRantai()
		if err := HitungUlangPajakNonPropEDM(h, f, pbuSekarang()); err != nil || len(f.panggil) != 0 {
			t.Fatalf("rantai jalan tanpa Choose Business: %v %v", err, f.panggil)
		}
	}
	h := pbuHalaman("1", "", "0")
	h.Setel(pt+"NoOffer", "UJI-M1-E1")
	if err := PilihBisnisEDM(h, pbuRantai(), pbuSekarang()); err != nil {
		t.Fatal(err)
	}
	h.SetelDaftar(DaftarXOL, []Baris{{"Currency": "IDR"}})
	f := pbuRantai()
	if err := HitungUlangPajakNonPropEDM(h, f, pbuSekarang()); err != nil || len(f.panggil) != 0 {
		t.Fatalf("polis Proporsional ikut rantai NonProp: %v %v", err, f.panggil)
	}
}
