package models

// Uji port XOL endorsemen (edm_xol.go, edm_xol_keluar.go, edm_xol_selisih.go). Angka karangan sederhana; nilai
// harapan dihitung tangan di komentar setiap uji. Data UJI-.

import (
	"errors"
	"testing"

	"github.com/cockroachdb/apd/v3"
)

// xeuSama - dua teks desimal bernilai sama (format boleh berbeda).
func xeuSama(t *testing.T, label, dapat, harap string) {
	t.Helper()
	d, _, err1 := apd.NewFromString(dapat)
	w, _, err2 := apd.NewFromString(harap)
	if err1 != nil || err2 != nil || d.Cmp(w) != 0 {
		t.Errorf("%s = %q, harap %s", label, dapat, harap)
	}
}

// xeuBaris membandingkan medan angka baris dengan harapan.
func xeuBaris(t *testing.T, label string, b Baris, harap map[string]string) {
	t.Helper()
	for m, v := range harap {
		xeuSama(t, label+"."+m, b[m], v)
	}
}

// xeuLapisan menulis satu layer master (gross/net/potongan per mata uang) di `daftar(si)`.
func xeuLapisan(h *Halaman, daftar string, si int, gross, net, ded []Baris) {
	h.SetelDaftar(JalurAnak(daftar, si, "GrossPremiumList"), gross)
	h.SetelDaftar(JalurAnak(daftar, si, "NetPremiumList"), net)
	h.SetelDaftar(JalurAnak(daftar, si, "DeductionList"), ded)
}

func xeuMV(mu, v string) Baris  { return Baris{"Currency": mu, "Value": v} }
func xeuPot(mu, v string) Baris { return Baris{"Currency": mu, "Deduction": v} }

// masterDuaLayer: layer 1 IDR gross 1000 / net 800 / potongan 200 (+ gross USD 7); layer 2 IDR 500 / 400 / 100.
func xeuMasterDuaLayer(h *Halaman, daftar string) {
	h.SetelDaftar(daftar, []Baris{{"LayerType": "UJI-LT", "Layer": "1", "LayerPartType": "UJI-PT", "LayerPart": "1"}, {"Layer": "2"}})
	xeuLapisan(h, daftar, 1, []Baris{xeuMV("IDR", "1000"), xeuMV("USD", "7")}, []Baris{xeuMV("IDR", "800")}, []Baris{xeuPot("IDR", "200")})
	xeuLapisan(h, daftar, 2, []Baris{xeuMV("IDR", "500")}, []Baris{xeuMV("IDR", "400")}, []Baris{xeuPot("IDR", "100")})
}

// InsertToTreatyXOLListEDM: 1 menyalin TreatyIn.Share ke ActualValue.Share; induk = jumlah layer; pajak bukan
// Inclusive dari potongan kumulatif 300: PPH 300x0,02 = 6, PPN 300x0,022 = 6,6; layer 1 CARI45 = .Deduction 200.
// USD: CARI31 tidak dinolkan per layer -> Currency induk USD, gross hanya 7 (CARI32 dinolkan 3.3.3.1).
func TestInsertToTreatyXOLListEDMSalinDanJumlah(t *testing.T) {
	h := HalamanBaru()
	h.Setel(pt+"FlagPPH", "true")
	h.Setel(pt+"TypeTax", "Exclusive")
	xeuMasterDuaLayer(h, jMaster+"Share")
	h.SetelDaftar(jMaster+"Installment", []Baris{{"Currency": "IDR"}, {"Currency": "USD"}})
	if err := InsertToTreatyXOLListEDM(h, IDMataUang{"IDR": "UJI-ID-IDR"}); err != nil {
		t.Fatal(err)
	}
	if n := len(h.AmbilDaftar(jMaster + "ActualValue.Share")); n != 2 {
		t.Fatalf("ActualValue.Share %d baris, harap 2 (salinan langkah 1)", n)
	}
	if g := h.AmbilDaftar(JalurAnak(jMaster+"ActualValue.Share", 2, "GrossPremiumList")); len(g) != 1 || g[0]["Value"] != "500" {
		t.Fatalf("daftar bersarang ikut tersalin: %v", g)
	}
	induk := h.AmbilDaftar(DaftarXOL)
	if len(induk) != 2 || induk[0]["Currency"] != "IDR" || induk[0]["IDCurrency"] != "UJI-ID-IDR" || induk[0]["DueTo"] != "" {
		t.Fatalf("induk %v", induk)
	}
	xeuBaris(t, "induk IDR", induk[0], map[string]string{"GrossPremi": "1500", "NetPremi": "1200", "DueToValue": "1200",
		"Deduction": "300", "BrokerageFeeSebenarnya": "300", "PPHValue": "6", "PPNValue": "6.6",
		"NetPremiAfterPPH": "1206", "NetPremiAfterPPN": "1206.6", "NetPremiAfterTax": "1212.6"})
	if induk[1]["Currency"] != "USD" || induk[1]["IDCurrency"] != "" {
		t.Errorf("induk USD %v", induk[1])
	}
	xeuBaris(t, "induk USD", induk[1], map[string]string{"GrossPremi": "7", "NetPremi": "0", "Deduction": "0"})
	l := h.AmbilDaftar(JalurAnak(DaftarXOL, 1, AnakLayerXOL))
	if len(l) != 2 || l[0]["Layer"] != "1" || l[0]["LayerType"] != "UJI-LT" || l[1]["Layer"] != "2" || l[0]["DueTo"] != "" {
		t.Fatalf("layer %v", l)
	}
	if l[0]["BrokerageFeeSebenarnya"] != "200" { // CARI45 = .Deduction mentah (bukan Inclusive)
		t.Errorf("BFS layer 1 %q", l[0]["BrokerageFeeSebenarnya"])
	}
	xeuBaris(t, "layer 1", l[0], map[string]string{"GrossPremi": "1000", "NetPremi": "800", "DueToValue": "800",
		"Deduction": "200", "PPHValue": "4", "PPNValue": "4.4", "NetPremiAfterPPH": "804", "NetPremiAfterPPN": "804.4",
		"NetPremiAfterTax": "808.4"})
	xeuBaris(t, "layer 2", l[1], map[string]string{"GrossPremi": "500", "PPHValue": "2", "NetPremiAfterTax": "404.2"})
}

// Inclusive: BFS induk = @divide(300, @divide(102.2,100,8)=1.022, 8) = 293.54207436; PPH = x0,02 = 5.8708414872;
// layer 1 CARI45 = @divide(200,1.022,8) = 195.69471624.
func TestInsertToTreatyXOLListEDMInclusive(t *testing.T) {
	h := HalamanBaru()
	h.Setel(pt+"FlagPPH", "true")
	h.Setel(pt+"TypeTax", TypeTaxInclusive)
	xeuMasterDuaLayer(h, jMaster+"ActualValue.Share")
	h.SetelDaftar(jMaster+"Share", []Baris{{"Layer": "UJI-TIDAK-DISALIN"}})
	h.SetelDaftar(jMaster+"Installment", []Baris{{"Currency": "IDR"}})
	if err := InsertToTreatyXOLListEDM(h, IDMataUang{}); err != nil {
		t.Fatal(err)
	}
	if h.AmbilDaftar(jMaster + "ActualValue.Share")[0]["Layer"] != "1" {
		t.Fatal("ActualValue.Share berisi - langkah 1 tidak menyalin Share")
	}
	xeuBaris(t, "induk", h.AmbilDaftar(DaftarXOL)[0], map[string]string{"BrokerageFeeSebenarnya": "293.54207436",
		"PPHValue": "5.8708414872"})
	xeuSama(t, "CARI45 layer 1", h.AmbilDaftar(JalurAnak(DaftarXOL, 1, AnakLayerXOL))[0]["BrokerageFeeSebenarnya"], "195.69471624")
}

// FlagRetroTreaty: induk = CARI44 layer TERAKHIR (100 + retro 20 = 120), Deduction 0; layer: CARI19 = potongan
// retro + CARI44 (10+200 = 210; 20+100 = 120), Deduction = CARI44 (200 / 100).
func TestInsertToTreatyXOLListEDMRetro(t *testing.T) {
	h := HalamanBaru()
	h.Setel(pt+"FlagRetroTreaty", "true")
	xeuMasterDuaLayer(h, jMaster+"ActualValue.Share")
	h.SetelDaftar(jMaster+"ActualValue.FacultativeShareList", []Baris{{}, {}})
	h.SetelDaftar(JalurAnak(jMaster+"ActualValue.FacultativeShareList", 1, "DeductionList"), []Baris{xeuPot("IDR", "10")})
	h.SetelDaftar(JalurAnak(jMaster+"ActualValue.FacultativeShareList", 2, "DeductionList"), []Baris{xeuPot("IDR", "20")})
	h.SetelDaftar(jMaster+"Installment", []Baris{{"Currency": "IDR"}})
	if err := InsertToTreatyXOLListEDM(h, IDMataUang{}); err != nil {
		t.Fatal(err)
	}
	xeuBaris(t, "induk", h.AmbilDaftar(DaftarXOL)[0], map[string]string{"GrossPremi": "120", "NetPremi": "120",
		"DueToValue": "120", "Deduction": "0"})
	l := h.AmbilDaftar(JalurAnak(DaftarXOL, 1, AnakLayerXOL))
	xeuBaris(t, "layer 1", l[0], map[string]string{"GrossPremi": "210", "NetPremi": "210", "DueToValue": "210", "Deduction": "200"})
	xeuBaris(t, "layer 2", l[1], map[string]string{"GrossPremi": "120", "Deduction": "100"})
	if _, ada := h.AmbilDaftar(DaftarXOL)[0]["PPHValue"]; ada {
		t.Error("FlagPPH kosong - langkah pajak 3.3.7 tidak jalan")
	}
}

// OldData: sumber TreatyIn TINGKAT ATAS, tujuan OldData.TreatyXOLList, TANPA nol per layer: layer 2 tanpa baris IDR
// -> CARI layer 1 terbawa dan dijumlah lagi: gross 1000+1000 = 2000, net 1600, potongan 400. Pajak (OldData.FlagPPH)
// induk memakai TypeTax BARU (Inclusive): @divide(400,1.022,8) = 391.38943249; CARI45 layer memakai OldData.TypeTax
// (kosong -> @toDecimal(CARI44) = 200).
func TestInsertToTreatyXOLListEDMOldDataTanpaNolPerLayer(t *testing.T) {
	h := HalamanBaru()
	atas := HalamanBaru()
	h.Setel(pt+"TypeTax", TypeTaxInclusive)
	h.Setel(od+"FlagPPH", "true")
	h.SetelDaftar(DaftarXOL, []Baris{{"GrossPremi": "UJI-TETAP"}})
	atas.SetelDaftar(jMaster+"Share", []Baris{{"Layer": "1"}, {"Layer": "2"}})
	xeuLapisan(atas, jMaster+"Share", 1, []Baris{xeuMV("IDR", "1000")}, []Baris{xeuMV("IDR", "800")}, []Baris{xeuPot("IDR", "200")})
	xeuLapisan(atas, jMaster+"Share", 2, []Baris{xeuMV("USD", "9")}, []Baris{xeuMV("USD", "9")}, []Baris{xeuPot("USD", "9")})
	atas.SetelDaftar(jMaster+"Installment", []Baris{{"Currency": "IDR"}})
	if err := InsertToTreatyXOLListEDMOldData(h, atas, IDMataUang{"IDR": "UJI-ID-IDR"}); err != nil {
		t.Fatal(err)
	}
	if h.AmbilDaftar(DaftarXOL)[0]["GrossPremi"] != "UJI-TETAP" {
		t.Fatal("TreatyXOLList baru tidak boleh disentuh")
	}
	induk := h.AmbilDaftar(od + "TreatyXOLList")
	if len(induk) != 1 || induk[0]["IDCurrency"] != "UJI-ID-IDR" || induk[0]["Currency"] != "IDR" {
		t.Fatalf("induk lama %v", induk)
	}
	xeuBaris(t, "induk lama", induk[0], map[string]string{"GrossPremi": "2000", "NetPremi": "1600", "Deduction": "400",
		"BrokerageFeeSebenarnya": "391.38943249"})
	l := h.AmbilDaftar(JalurAnak(od+"TreatyXOLList", 1, AnakLayerXOL))
	xeuSama(t, "CARI45 layer 1", l[0]["BrokerageFeeSebenarnya"], "200")
	if l[1]["Currency"] != "" || l[1]["Layer"] != "2" { // Page-New per layer: layer 2 tanpa baris IDR
		t.Errorf("layer 2 %v", l[1])
	}
}

// Langkah 3: baris ReinsuranceListTONP ber-ReinsName lain dibuang. `[dugaan]` perulangan menurut subskrip: dua baris
// berbeda berurutan -> baris kedua bergeser dan TIDAK diperiksa (bertahan); daftar bersarang ikut bergeser.
func TestBuangReinsBerbedaMenurutSubskrip(t *testing.T) {
	h := HalamanBaru()
	h.SetelDaftar(jMaster+"ShareReins", []Baris{{}})
	j := xeTONP(1)
	h.SetelDaftar(j, []Baris{{"ReinsName": "UJI-X"}, {"ReinsName": "UJI-Y"}, {"ReinsName": "UJI-SOB"}})
	h.SetelDaftar(JalurAnak(j, 2, "GrossPremiumList"), []Baris{xeuMV("IDR", "2")})
	h.SetelDaftar(JalurAnak(j, 3, "GrossPremiumList"), []Baris{xeuMV("IDR", "3")})
	xeBuangReinsBerbeda(h, "UJI-SOB")
	d := h.AmbilDaftar(j)
	if len(d) != 2 || d[0]["ReinsName"] != "UJI-Y" || d[1]["ReinsName"] != "UJI-SOB" {
		t.Fatalf("TONP %v", d)
	}
	if g := h.AmbilDaftar(JalurAnak(j, 1, "GrossPremiumList")); len(g) != 1 || g[0]["Value"] != "2" {
		t.Errorf("anak baris 2 bergeser ke 1: %v", g)
	}
	if g := h.AmbilDaftar(JalurAnak(j, 2, "GrossPremiumList")); len(g) != 1 || g[0]["Value"] != "3" {
		t.Errorf("anak baris 3 bergeser ke 2: %v", g)
	}
	if _, ada := h.Daftar[JalurAnak(j, 3, "GrossPremiumList")]; ada {
		t.Error("kunci subskrip 3 tertinggal")
	}
}

// xeuMasterKeluar: ShareReins(1).TONP = [A (SOB): gross 100 net 80 pot 20; B (lain) dibuang; C (SOB): 50/40/10].
func xeuMasterKeluar(h *Halaman) {
	h.Setel(pt+"SOBName", "UJI-SOB")
	h.SetelDaftar(jMaster+"ShareReins", []Baris{{}})
	j := xeTONP(1)
	h.SetelDaftar(j, []Baris{{"ReinsName": "UJI-SOB", "Layer": "1"}, {"ReinsName": "UJI-LAIN", "Layer": "9"},
		{"ReinsName": "UJI-SOB", "Layer": "2"}})
	xeuLapisan(h, j, 1, []Baris{xeuMV("IDR", "100")}, []Baris{xeuMV("IDR", "80")}, []Baris{xeuPot("IDR", "20")})
	xeuLapisan(h, j, 2, []Baris{xeuMV("IDR", "999")}, []Baris{xeuMV("IDR", "999")}, []Baris{xeuPot("IDR", "999")})
	xeuLapisan(h, j, 3, []Baris{xeuMV("IDR", "50")}, []Baris{xeuMV("IDR", "40")}, []Baris{xeuPot("IDR", "10")})
	h.SetelDaftar(jMaster+"Installment", []Baris{{"Currency": "IDR"}})
}

// InsertToTreatyOutXOLList: induk gross 150, net DUA KALI (80+40)x2 = 240, DueToValue 120, potongan 30; tanpa pajak.
func TestInsertToTreatyOutXOLList(t *testing.T) {
	h := HalamanBaru()
	h.Setel(pt+"FlagPPH", "true")
	xeuMasterKeluar(h)
	if err := InsertToTreatyOutXOLList(h, IDMataUang{"IDR": "UJI-ID-IDR"}); err != nil {
		t.Fatal(err)
	}
	induk := h.AmbilDaftar(DaftarXOL)
	xeuBaris(t, "induk", induk[0], map[string]string{"GrossPremi": "150", "NetPremi": "240", "DueToValue": "120", "Deduction": "30"})
	if induk[0]["DueTo"] != "" || induk[0]["IDCurrency"] != "UJI-ID-IDR" || induk[0]["PPHValue"] != "" {
		t.Errorf("induk %v", induk[0])
	}
	l := h.AmbilDaftar(JalurAnak(DaftarXOL, 1, AnakLayerXOL))
	if len(l) != 2 || l[0]["Layer"] != "1" || l[1]["Layer"] != "2" || l[0]["BrokerageFeeSebenarnya"] != "" {
		t.Fatalf("layer %v", l)
	}
	xeuBaris(t, "layer 2", l[1], map[string]string{"GrossPremi": "50", "NetPremi": "40", "DueToValue": "40", "Deduction": "10"})
}

// InsertToTreatyOutXOLListEDMOldData: potongan dari SpreadingTONP(1).ReinsuranceListTONP(1).LayerList(ti); net SEKALI.
func TestInsertToTreatyOutXOLListEDMOldData(t *testing.T) {
	h := HalamanBaru()
	xeuMasterKeluar(h)
	lapis := JalurAnak(JalurAnak(jMaster+"SpreadingTONP", 1, "ReinsuranceListTONP"), 1, "LayerList")
	h.SetelDaftar(lapis, []Baris{{}, {}})
	h.SetelDaftar(JalurAnak(lapis, 1, "DeductionList"), []Baris{xeuPot("IDR", "7")})
	h.SetelDaftar(JalurAnak(lapis, 2, "DeductionList"), []Baris{xeuPot("IDR", "3")})
	if err := InsertToTreatyOutXOLListEDMOldData(h, IDMataUang{}); err != nil {
		t.Fatal(err)
	}
	if len(h.AmbilDaftar(DaftarXOL)) != 0 {
		t.Fatal("TreatyXOLList baru tidak boleh ditulis")
	}
	xeuBaris(t, "induk lama", h.AmbilDaftar(od + "TreatyXOLList")[0], map[string]string{"GrossPremi": "150",
		"NetPremi": "120", "Deduction": "10"})
	l := h.AmbilDaftar(JalurAnak(od+"TreatyXOLList", 1, AnakLayerXOL))
	xeuBaris(t, "layer lama 1", l[0], map[string]string{"Deduction": "7"})
}

// xeuSelisih menyiapkan satu induk dan satu layer baru/lama.
func xeuSelisih(h *Halaman, baru, lama Baris) {
	h.SetelDaftar(DaftarXOL, []Baris{{"Currency": "IDR", "IDCurrency": "UJI-ID-IDR"}})
	h.SetelDaftar(JalurAnak(DaftarXOL, 1, AnakLayerXOL), []Baris{baru})
	h.SetelDaftar(od+"TreatyXOLList", []Baris{{}})
	h.SetelDaftar(JalurAnak(od+"TreatyXOLList", 1, AnakLayerXOL), []Baris{lama})
}

// 1.2.1: turun -> 0 (gross 100 < 150), naik -> selisih (net 300-200 = 100); DueTo dari DueToValue selisih.
// Baris lama yang tak ada = 0 (BrokerageFeeSebenarnya baru 5, lama kosong -> 5).
func TestCalculateDifferenceEDMBatasBawahNol(t *testing.T) {
	h := HalamanBaru()
	xeuSelisih(h, Baris{"Layer": "1", "Currency": "IDR", "GrossPremi": "100", "NetPremi": "300", "DueToValue": "300",
		"Deduction": "50", "BrokerageFeeSebenarnya": "5"},
		Baris{"GrossPremi": "150", "NetPremi": "200", "DueToValue": "200", "Deduction": "50"})
	if err := CalculateDifferenceEDM(h); err != nil {
		t.Fatal(err)
	}
	l := h.AmbilDaftar(JalurAnak(DaftarSelisihXOL, 1, AnakLayerXOL))[0]
	xeuBaris(t, "layer", l, map[string]string{"GrossPremi": "0", "NetPremi": "100", "DueToValue": "100", "Deduction": "0",
		"BrokerageFeeSebenarnya": "5"})
	if l["DueTo"] != "DUE TO US" || l["Layer"] != "1" || l["IDCurrency"] != "" {
		t.Errorf("layer %v", l)
	}
	induk := h.AmbilDaftar(DaftarSelisihXOL)[0]
	if induk["DueTo"] != "DUE TO YOU" || induk["IDCurrency"] != "UJI-ID-IDR" { // lokal awal 0
		t.Errorf("induk %v", induk)
	}
	xeuBaris(t, "induk (langkah 2)", induk, map[string]string{"NetPremi": "100", "GrossPremi": "0"})
}

// 1.2.3 prorata 50%: (300-200) x @Math.divide(50,100,8) = 50; turun tetap 0.
func TestCalculateDifferenceEDMProrata(t *testing.T) {
	h := HalamanBaru()
	h.Setel(jMaster+"IsProRate", "True")
	h.Setel(jMaster+"ProRatePercent", "50")
	xeuSelisih(h, Baris{"GrossPremi": "100", "NetPremi": "300"}, Baris{"GrossPremi": "150", "NetPremi": "200"})
	if err := CalculateDifferenceEDM(h); err != nil {
		t.Fatal(err)
	}
	xeuBaris(t, "layer", h.AmbilDaftar(JalurAnak(DaftarSelisihXOL, 1, AnakLayerXOL))[0],
		map[string]string{"NetPremi": "50", "GrossPremi": "0"})
}

// 1.2.4 (FlagPPH && EDMType 3): BFS = potongan selisih 100 (bukan Inclusive); PPH 2; PPN 2.2; AfterTax = 100+2.2+2.
// EDMType 1 dengan FlagPPH: pajak TIDAK dihitung ulang (selisih mentah berbatas bawah).
func TestCalculateDifferenceEDMPajakAdjPremi(t *testing.T) {
	for _, tc := range []struct {
		edm, pph, tax string
	}{{"3", "2", "104.2"}, {"1", "0", "0"}} {
		h := HalamanBaru()
		h.Setel(pt+"FlagPPH", "true")
		h.Setel(pt+"EDMType", tc.edm)
		xeuSelisih(h, Baris{"Deduction": "150", "NetPremi": "300"}, Baris{"Deduction": "50", "NetPremi": "200"})
		if err := CalculateDifferenceEDM(h); err != nil {
			t.Fatal(err)
		}
		xeuBaris(t, "EDMType "+tc.edm, h.AmbilDaftar(JalurAnak(DaftarSelisihXOL, 1, AnakLayerXOL))[0],
			map[string]string{"PPHValue": tc.pph, "NetPremiAfterTax": tc.tax})
	}
}

// 1.2.5 EDMType 4 / 2: selisih mentah, boleh negatif (100-150 = -50); DueTo YOU.
func TestCalculateDifferenceEDMBatalMentah(t *testing.T) {
	for _, edm := range []string{"4", "2"} {
		h := HalamanBaru()
		h.Setel(pt+"EDMType", edm)
		xeuSelisih(h, Baris{"GrossPremi": "100", "DueToValue": "0"}, Baris{"GrossPremi": "150", "DueToValue": "80"})
		if err := CalculateDifferenceEDM(h); err != nil {
			t.Fatal(err)
		}
		l := h.AmbilDaftar(JalurAnak(DaftarSelisihXOL, 1, AnakLayerXOL))[0]
		xeuBaris(t, "EDMType "+edm, l, map[string]string{"GrossPremi": "-50", "DueToValue": "-80"})
		if l["DueTo"] != "DUE TO YOU" {
			t.Errorf("DueTo %q", l["DueTo"])
		}
	}
}

// DueTo induk ke-2 memakai local.duetovalue layer terakhir induk ke-1 (positif -> US); baris selisih lama ke-3
// tertinggal (tidak dihapus) dan ikut langkah 2.
func TestCalculateDifferenceEDMDueToIndukBasiDanBarisLama(t *testing.T) {
	h := HalamanBaru()
	h.SetelDaftar(DaftarXOL, []Baris{{"Currency": "IDR"}, {"Currency": "USD"}})
	h.SetelDaftar(JalurAnak(DaftarXOL, 1, AnakLayerXOL), []Baris{{"DueToValue": "10"}})
	h.SetelDaftar(JalurAnak(DaftarXOL, 2, AnakLayerXOL), []Baris{{"DueToValue": "0"}})
	h.SetelDaftar(DaftarSelisihXOL, []Baris{{}, {}, {"Currency": "UJI-LAMA"}})
	h.SetelDaftar(JalurAnak(DaftarSelisihXOL, 3, AnakLayerXOL), []Baris{{"NetPremi": "7"}})
	if err := CalculateDifferenceEDM(h); err != nil {
		t.Fatal(err)
	}
	d := h.AmbilDaftar(DaftarSelisihXOL)
	if len(d) != 3 || d[0]["DueTo"] != "DUE TO YOU" || d[1]["DueTo"] != "DUE TO US" || d[2]["Currency"] != "UJI-LAMA" {
		t.Fatalf("induk %v", d)
	}
	xeuSama(t, "induk lama NetPremi", d[2]["NetPremi"], "7")
}

// Nilai bukan angka -> galat (tidak diam-diam 0).
func TestCalculateDifferenceEDMBukanAngka(t *testing.T) {
	h := HalamanBaru()
	xeuSelisih(h, Baris{"GrossPremi": "UJI-X"}, Baris{})
	if err := CalculateDifferenceEDM(h); !errors.Is(err, ErrBukanAngka) {
		t.Fatalf("galat %v", err)
	}
}
