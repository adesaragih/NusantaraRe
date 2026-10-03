package rekonsiliasi

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/alat/deidentifikasi"
	"nusantarare/modul/nbfacin/backend/services/premium"
)

const (
	berkasDaftarIzin = "../premium/testdata/daftarizin/premi.json"
	folderKasus      = "../premium/testdata/kasus"
)

func baca(t *testing.T, jalur string) []byte {
	t.Helper()
	isi, err := os.ReadFile(jalur)
	if err != nil {
		t.Fatal(err)
	}
	return isi
}

// TestDaftarIzinCocokEksak - kedelapan kasus nyata PA dan MBU identik sampai
// digit terakhir, dengan asal rumusnya disebut.
func TestDaftarIzinCocokEksak(t *testing.T) {
	hasil, err := DaftarIzin(baca(t, berkasDaftarIzin))
	if err != nil {
		t.Fatal(err)
	}
	if len(hasil) != 8 {
		t.Fatalf("%d hasil, mau 8", len(hasil))
	}
	for _, h := range hasil {
		if h.Status != Cocok || h.Langkah == "" {
			t.Errorf("%+v: mau cocok dengan langkah terisi", h)
		}
	}
}

// TestSelisihDilaporkan - selisih menyebut kasus, lini, langkah, dan kedua nilai;
// satu digit terakhir pun bukan "hampir sama".
func TestSelisihDilaporkan(t *testing.T) {
	isi := `[{"kasus":"UJI-rekaan","input":{"LiniBisnis":"PA","CalculateMethod":"3","TSI":"2787500000",
	"Rate":"13.47311827957","ProRatePercent":"400.27397260273972602700"},"premi":"37556317.2044"}]`
	hasil, err := DaftarIzin([]byte(isi))
	if err != nil {
		t.Fatal(err)
	}
	mau := Hasil{Kasus: "UJI-rekaan", Lini: "PA", Langkah: "CalculatePremiPA_FacIn langkah 6 L1003",
		Status: Selisih, SistemLama: "37556317.2044", SistemBaru: "37556317.2043"}
	if len(hasil) != 1 || hasil[0] != mau {
		t.Fatalf("dapat %+v\nmau   %+v", hasil, mau)
	}
}

// TestLiniBelumDiportBelumTercakup - lini tanpa rumus dilaporkan belum tercakup,
// bukan lulus dan bukan galat.
func TestLiniBelumDiportBelumTercakup(t *testing.T) {
	isi := `[{"kasus":"UJI-rekaan","input":{"LiniBisnis":"BONDING","TSI":"1","Rate":"1"},"premi":"1"}]`
	hasil, err := DaftarIzin([]byte(isi))
	if err != nil || len(hasil) != 1 || hasil[0].Status != BelumTercakup {
		t.Fatalf("dapat %+v (%v)", hasil, err)
	}
}

// TestBerkasKasusP5 - patokan tiket 18: setiap premi coverage kelima kasus nyata P-5
// cocok sampai digit terakhir dengan rumus lininya. Nilai dasar akseptasi tetap belum
// tercakup (nilai lama tidak terekam).
func TestBerkasKasusP5(t *testing.T) {
	berkas, err := filepath.Glob(filepath.Join(folderKasus, "*.json"))
	if err != nil || len(berkas) != 5 {
		t.Fatalf("%d fixture (%v)", len(berkas), err)
	}
	// Jumlah coverage per kasus dan lini: dihitung dua cara (tes ini dan urai Python
	// independen), jendela + perintah: `../premium/testdata/kasus/README.md`.
	mau := map[string]map[string]int{"edm-fire-1": {"FIRE": 45}, "nb-fire-1": {"FIRE": 5}, "rnw-fire-1": {"FIRE": 4},
		"nb-kredit-1": {"ANEKA": 1}, "nb-marinecargo-1": {"MARINE CARGO": 104}}
	// Tiket 07: total per item (TotalGrossPremi) dan per lokasi per mata uang
	// (TotalTSIPremiGrossList). Dihitung dua cara - tes Go ini dan urai Python atas
	// fixture (02-10-2026): item cocok 9 (EDM) + 1 (NB) + 1 (RNW); lokasi NB dan RNW (1
	// item) cocok; lokasi EDM (9 item USD, cabang 1.3.4 lewat double, A48) belum tercakup.
	total := map[string]int{}
	mauTotal := map[string]int{"nb-fire-1 cocok": 2, "rnw-fire-1 cocok": 2, "edm-fire-1 cocok": 9, "edm-fire-1 belum tercakup": 1}
	for _, b := range berkas {
		nama := strings.TrimSuffix(filepath.Base(b), ".json")
		hasil, err := BerkasKasus(nama, baca(t, b))
		if err != nil {
			t.Fatalf("%s: %v", nama, err)
		}
		cocok := map[string]int{}
		var bayar []Hasil
		for _, h := range hasil[:len(hasil)-1] {
			if h.Langkah == langkahPembayaran {
				bayar = append(bayar, h)
				continue
			}
			if strings.HasPrefix(h.Langkah, langkahTotal) {
				total[nama+" "+string(h.Status)]++
				if h.Status == Selisih || h.Status == Galat {
					t.Logf("%s: %+v", nama, h)
				}
				continue
			}
			// Kasus EDM dibandingkan dengan rumus NB; catatannya wajib menyebutnya.
			if edm := strings.Contains(h.Catatan, "StatusBusiness 3"); edm != strings.HasPrefix(nama, "edm-") {
				t.Errorf("%s: catatan EDM %v: %q", nama, edm, h.Catatan)
			}
			if h.Status != Cocok {
				t.Errorf("%s: %+v", nama, h)
				continue
			}
			cocok[h.Lini]++
		}
		// Tiket 19: pembayaran MARINE dihitung, tetapi `pyWorkPage.Policy.Payment` lama
		// tidak terekam (akar fixture = OfferFacIn). Jumlah 26395202 dihitung dua cara
		// (Python decimal dan fractions atas CoverageList, 01-10-2026).
		if mauBayar := nama == "nb-marinecargo-1"; mauBayar != (len(bayar) == 1) ||
			(mauBayar && (bayar[0].Status != BelumTercakup || bayar[0].SistemBaru != "26395202.0000")) {
			t.Errorf("%s: baris pembayaran %+v", nama, bayar)
		}
		if !reflect.DeepEqual(cocok, mau[nama]) {
			t.Errorf("%s: cocok per lini %v, mau %v", nama, cocok, mau[nama])
		}
		if akhir := hasil[len(hasil)-1]; akhir.Langkah != langkahNilaiDasar || akhir.Status != BelumTercakup {
			t.Errorf("%s: baris nilai dasar %+v", nama, akhir)
		}
	}
	if !reflect.DeepEqual(total, mauTotal) {
		t.Errorf("baris total %v, mau %v", total, mauTotal)
	}
}

// TestBerkasKasusSelisihDilaporkan - premi lama satu coverage ANEKA diubah satu
// digit terakhir: selisih dengan jalur coverage-nya, bukan disesuaikan.
func TestBerkasKasusSelisihDilaporkan(t *testing.T) {
	var akar map[string]any
	if err := json.Unmarshal(baca(t, filepath.Join(folderKasus, "nb-kredit-1.json")), &akar); err != nil {
		t.Fatal(err)
	}
	cov := akar["LocationList"].([]any)[0].(map[string]any)["Property"].(map[string]any)["RiskLocation"].(map[string]any)["OccupationList"].([]any)[0].(map[string]any)["AnekaList"].([]any)[0].(map[string]any)["CoverageList"].([]any)[0].(map[string]any)
	lama := cov["Premium"].(string)
	ubah := lama[:len(lama)-1] + map[byte]string{'9': "8"}[lama[len(lama)-1]]
	if ubah == lama[:len(lama)-1] {
		ubah = lama[:len(lama)-1] + "9"
	}
	cov["Premium"] = ubah
	isi, err := json.Marshal(akar)
	if err != nil {
		t.Fatal(err)
	}
	hasil, err := BerkasKasus("UJI-rekaan", isi)
	if err != nil {
		t.Fatal(err)
	}
	if r := Ringkas(hasil); r[Selisih] != 1 {
		t.Fatalf("ringkasan %v", r)
	}
	for _, h := range hasil {
		if h.Status == Selisih && (h.Lini != "ANEKA" || !strings.Contains(h.Langkah, "CoverageList[0]") || h.SistemLama != ubah) {
			t.Errorf("selisih tanpa rincian: %+v", h)
		}
	}
}

// TestBerkasMentahDitolak - kerangka tidak pernah memproses berkas kasus mentah.
func TestBerkasMentahDitolak(t *testing.T) {
	_, err := BerkasKasus("mentah", []byte(`{"InsuredName":"Rekaan","QuotationData":{"BusinessType":"Fire"}}`))
	if !errors.Is(err, deidentifikasi.ErrBelumBersih) {
		t.Fatalf("galat %v, mau ErrBelumBersih", err)
	}
}

// TestDeterministik - dua kali jalan, hasil identik dan urut.
func TestDeterministik(t *testing.T) {
	a, err := DaftarIzin(baca(t, berkasDaftarIzin))
	if err != nil {
		t.Fatal(err)
	}
	b, _ := DaftarIzin(baca(t, berkasDaftarIzin))
	if !reflect.DeepEqual(a, b) {
		t.Fatal("hasil berbeda antar-jalan")
	}
	r := Ringkas(a)
	if r[Cocok] != 8 || len(r) != 1 {
		t.Fatalf("ringkasan %v", r)
	}
}

const berkasAgregatMBU = "../premium/testdata/daftarizin/mbu_mata_uang.json"

// TestAgregatMBU - CurrencyList.Premium lima kasus MBU NB nyata (tiket 07):
// empat cocok, kasus #5 selisih terbuka - dilaporkan, bukan disembunyikan.
func TestAgregatMBU(t *testing.T) {
	hasil, err := AgregatMBU(baca(t, berkasAgregatMBU))
	if err != nil {
		t.Fatal(err)
	}
	if r := Ringkas(hasil); r[Cocok] != 4 || r[Selisih] != 1 || len(hasil) != 5 {
		t.Fatalf("ringkasan %v (%d hasil)", r, len(hasil))
	}
	mau := Hasil{Kasus: "kasus MBU NB #5", Lini: "MBU", Langkah: langkahAgregatMBU + " IDR", Status: Selisih,
		SistemLama: "309703644.8", SistemBaru: "201301628.0000"}
	if hasil[4] != mau {
		t.Fatalf("dapat %+v\nmau   %+v", hasil[4], mau)
	}
}

// TestAgregatMBUNilaiEksak - nilai sama dengan skala teks berbeda = cocok (A31,
// butir 52: Decimal); satu digit berbeda = selisih.
func TestAgregatMBUNilaiEksak(t *testing.T) {
	isi := func(premi string) []byte {
		return []byte(`[{"kasus":"UJI-rekaan","urutanMasterMataUang":["IDR"],"coverage":[{"TSI":"1000","Rate":"2",` +
			`"Loading":"","ProRatePercentCoverage":"100","MataUang":"IDR","FlagDelete":""}],"currencyList":[{"Name":"IDR","Premium":"` + premi + `"}]}]`)
	}
	for premi, mau := range map[string]Status{"20": Cocok, "20.0000": Cocok, "20.0001": Selisih, "19.9999": Selisih} {
		hasil, err := AgregatMBU(isi(premi))
		if err != nil || len(hasil) != 1 || hasil[0].Status != mau {
			t.Errorf("premi lama %s: dapat %+v (%v), mau %s", premi, hasil, err, mau)
		}
	}
}

// TestAgregatMBUMataUangSepihak - mata uang yang hanya ada di CurrencyList lama
// tetap dilaporkan sebagai selisih, tidak terlewat diam-diam.
func TestAgregatMBUMataUangSepihak(t *testing.T) {
	isi := []byte(`[{"kasus":"UJI-rekaan","urutanMasterMataUang":["IDR"],"coverage":[{"TSI":"1000","Rate":"2","Loading":"",` +
		`"ProRatePercentCoverage":"100","MataUang":"IDR","FlagDelete":""}],"currencyList":[{"Name":"IDR","Premium":"20"},{"Name":"USD","Premium":"5"}]}]`)
	hasil, err := AgregatMBU(isi)
	if err != nil {
		t.Fatal(err)
	}
	if len(hasil) != 2 || hasil[1].Langkah != langkahAgregatMBU+" USD" || hasil[1].Status != Selisih || hasil[1].SistemBaru != "" {
		t.Fatalf("dapat %+v", hasil)
	}
}

// TestMasukanCoverage - dua pemetaan yang tidak dibedakan fixture P-5 (semua
// `LostLimit` = "100", tidak ada master policy - README fixture): ejaan korpus
// `.LostLimit` (CountPremiCoverageAneka, CountPremi_ACT) dan master policy MARINE
// `PolicyType == 1 && IsMOP == "MOP"` (CountGPWMarinePAMbu_Act 1.1.1.1.2).
func TestMasukanCoverage(t *testing.T) {
	c := map[string]any{"LostLimit": "75", "LossLimit": "10"}
	for _, u := range []struct {
		qd     map[string]any
		master bool
	}{
		{map[string]any{"PolicyType": "1", "IsMOP": "MOP"}, true},
		{map[string]any{"PolicyType": "2", "IsMOP": "MOP"}, false},
		{map[string]any{"PolicyType": "1", "IsMOP": ""}, false},
	} {
		in := masukanCoverage(map[string]any{"QuotationData": u.qd}, c, premium.LiniMarineCargo, false)
		if in.MasterPolicy != u.master || in.LossLimit != "75" {
			t.Errorf("%v: MasterPolicy %v LossLimit %q", u.qd, in.MasterPolicy, in.LossLimit)
		}
	}
}

// TestBerkasKasusEDMBedaTidakDiport - CountGrossPremiEDM_Act cabang FIRE
// (`NB FacIn\Activity\CountGrossPremiEDM_Act.xml` 1.1.2.1.4-1.1.2.1.4.3) menimpa
// `.TSI` coverage dengan `TSIObjectItem` item dan menegasikan premi coverage ber-
// `FlagDelete == "1"`. Keduanya belum diport (A44): kasus EDM yang memicunya =
// galat, bukan dibandingkan dengan rumus NB. Lini EDM selain FIRE = belum tercakup.
func TestBerkasKasusEDMBedaTidakDiport(t *testing.T) {
	nCovItem0 := 0
	ubahEDM := func(ubah func(item, cov map[string]any)) []Hasil {
		t.Helper()
		var akar map[string]any
		if err := json.Unmarshal(baca(t, filepath.Join(folderKasus, "edm-fire-1.json")), &akar); err != nil {
			t.Fatal(err)
		}
		item := akar["LocationList"].([]any)[0].(map[string]any)["Property"].(map[string]any)["PropertyItemList"].([]any)[0].(map[string]any)
		nCovItem0 = len(item["CoverageList"].([]any))
		ubah(item, item["CoverageList"].([]any)[0].(map[string]any))
		isi, err := json.Marshal(akar)
		if err != nil {
			t.Fatal(err)
		}
		hasil, err := BerkasKasus("UJI-rekaan", isi)
		if err != nil {
			t.Fatal(err)
		}
		return barisPremiSaja(hasil)
	}
	// TSIObjectItem beda: `[dugaan]` 1.1.2.1.4 tidak menimpa `.TSI` (pasangan itu
	// hanya di pyParamArray langkah bermetode kosong) - tetap dibandingkan, bertanda.
	beda := ubahEDM(func(item, _ map[string]any) { item["TSIObjectItem"] = "1" })
	bertanda := 0
	for _, h := range beda {
		if strings.Contains(h.Catatan, "[dugaan] TSIObjectItem") {
			bertanda++
		}
	}
	if r := Ringkas(beda); r[Galat] != 0 || r[Cocok] != 45 || bertanda != nCovItem0 {
		t.Errorf("TSI item beda: ringkasan %v, bertanda %d, mau 45 cocok dan %d bertanda", r, bertanda, nCovItem0)
	}
	// Nilai sama, teks beda (`8400000` lawan `8400000.0`): BUKAN perbedaan, tanpa tanda.
	samaNilaiBedaTeks := ubahEDM(func(item, cov map[string]any) { item["TSIObjectItem"] = cov["TSI"].(string) + ".0" })
	for _, h := range samaNilaiBedaTeks {
		if strings.Contains(h.Catatan, "[dugaan] TSIObjectItem") {
			t.Errorf("TSI item sama nilai bertanda: %+v", h)
		}
	}
	if r := Ringkas(ubahEDM(func(_, cov map[string]any) { cov["FlagDelete"] = "1" })); r[Galat] != 1 || r[Cocok] != 44 {
		t.Errorf("FlagDelete coverage: ringkasan %v, mau 1 galat 44 cocok", r)
	}
	// FlagDelete item juga "1": syarat negasi `Local.deleteobjitem!="1"` tidak terpenuhi.
	if r := Ringkas(ubahEDM(func(item, cov map[string]any) { item["FlagDelete"] = "1"; cov["FlagDelete"] = "1" })); r[Galat] != 0 || r[Cocok] != 45 {
		t.Errorf("FlagDelete item dan coverage: ringkasan %v", r)
	}
}

// TestBerkasKasusEDMLiniLainBelumTercakup - StatusBusiness 3 pada lini selain FIRE.
func TestBerkasKasusEDMLiniLainBelumTercakup(t *testing.T) {
	var akar map[string]any
	if err := json.Unmarshal(baca(t, filepath.Join(folderKasus, "nb-kredit-1.json")), &akar); err != nil {
		t.Fatal(err)
	}
	akar["QuotationData"].(map[string]any)["StatusBusiness"] = "3"
	isi, err := json.Marshal(akar)
	if err != nil {
		t.Fatal(err)
	}
	hasil, err := BerkasKasus("UJI-rekaan", isi)
	if err != nil {
		t.Fatal(err)
	}
	if r := Ringkas(hasil); r[BelumTercakup] != 2 || r[Cocok] != 0 {
		t.Errorf("ringkasan %v, mau 2 belum tercakup (premi ANEKA + nilai dasar)", r)
	}
}

// kreditDiubah - fixture nb-kredit-1 dengan perubahan buatan, lalu BerkasKasus.
func kreditDiubah(t *testing.T, ubah func(akar, aneka, cov map[string]any)) []Hasil {
	t.Helper()
	var akar map[string]any
	if err := json.Unmarshal(baca(t, filepath.Join(folderKasus, "nb-kredit-1.json")), &akar); err != nil {
		t.Fatal(err)
	}
	risk := akar["LocationList"].([]any)[0].(map[string]any)["Property"].(map[string]any)["RiskLocation"].(map[string]any)
	aneka := risk["OccupationList"].([]any)[0].(map[string]any)["AnekaList"].([]any)[0].(map[string]any)
	ubah(akar, aneka, aneka["CoverageList"].([]any)[0].(map[string]any))
	isi, err := json.Marshal(akar)
	if err != nil {
		t.Fatal(err)
	}
	hasil, err := BerkasKasus("UJI-rekaan", isi)
	if err != nil {
		t.Fatal(err)
	}
	return hasil
}

// TestBerkasKasusCabangBuatan - dua cabang yang tidak dijangkau data nyata (celah
// tes mutasi, bukan mutan ekuivalen): IsMBD dan jalur GOLF. Masukan BUATAN.
func TestBerkasKasusCabangBuatan(t *testing.T) {
	// GOLF: CoverageList di RiskLocation.AnekaList (tanpa OccupationList). Rumus sama
	// dengan ANEKA, jadi premi tersimpan tetap cocok.
	golf := kreditDiubah(t, func(akar, aneka, _ map[string]any) {
		akar["QuotationData"].(map[string]any)["BusinessType"] = "GolfInsurance"
		risk := akar["LocationList"].([]any)[0].(map[string]any)["Property"].(map[string]any)["RiskLocation"].(map[string]any)
		risk["AnekaList"] = []any{aneka}
		delete(risk, "OccupationList")
	})
	n := 0
	for _, h := range golf {
		if h.Lini == "GOLF" && h.Status == Cocok {
			n++
		}
	}
	if n != 1 {
		t.Errorf("GOLF: %d coverage cocok, mau 1: %+v", n, golf)
	}
	// IsMBD + IndemnityPercentage terisi: rumus MBD (× indemnity) dipakai, sehingga
	// premi tersimpan (rumus non-MBD) menjadi selisih - bukti cabang dijangkau.
	mbd := kreditDiubah(t, func(akar, _, cov map[string]any) {
		akar["QuotationData"].(map[string]any)["BusinessType"] = "MBD"
		cov["IndemnityPercentage"] = "50"
	})
	for _, h := range mbd {
		if h.Langkah == langkahPremi || strings.HasPrefix(h.Langkah, langkahPremi+" |") {
			t.Logf("MBD: %+v", h)
		}
	}
	if r := Ringkas(mbd); r[Selisih] != 1 {
		t.Errorf("MBD: ringkasan %v, mau 1 selisih", r)
	}
}

// barisPremiSaja - hanya baris premi coverage (tanpa baris total, pembayaran, nilai dasar).
func barisPremiSaja(hasil []Hasil) []Hasil {
	var premi []Hasil
	for _, h := range hasil {
		if !strings.HasPrefix(h.Langkah, langkahTotal) && h.Langkah != langkahPembayaran && h.Langkah != langkahNilaiDasar {
			premi = append(premi, h)
		}
	}
	return premi
}
