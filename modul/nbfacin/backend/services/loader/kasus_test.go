package loader

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// Lima kasus NYATA hasil de-identifikasi (tiket 15) - dipakai bersama paket
// premium dan rekonsiliasi. Tidak ada nilai kasus yang disalin ke berkas ini: yang
// dikunci hanya hitungan baris, kelompok lini, dan nama kolom - kecuali satu nilai
// penunjuk posisi ("7", `IndexCargo` marine) yang dikutip atas keputusan butir 72 sebagai
// bukti; bukan data pribadi.
const folderKasus = "../premium/testdata/kasus"

type emasKasus struct {
	cob                              string
	terpetakan, takDikenal, penunjuk int
	baris                            map[string]int
	uang                             map[string]int
}

// Angka emas 02-10-2026 (diperbarui sesudah butir 68 dan 69). Pemeriksa silang
// T_COVERAGELIST: sensus coverage README fixture (dua cara, 01-10-2026) - 45 / 5 /
// 4 / 1 / 104, sama. Butir 70: CoverageInitial (104) dan AdditionalShip (416 medan) kini
// berkolom (terpetakan marine +520). Butir 72: ke-600 penunjuk Idx*/Index* (321 + 45 + 5 +
// 208 + 21) kini di kolom teks mentahnya - terpetakan naik tepat sebanyak itu; penampung
// tinggal IsCedingConfirm. IsCedingConfirm
// (K-069 7b, butir 69) di kelima kasus 3 + 10 + 9 + 2 + 10 = 34 = pengurai Python
// independen atas ViewSuggest (02-10-2026). Mata uang: coverage/aneka berhalaman
// Currency (marine, kredit) kini berkode; coverage FIRE tanpa halaman Currency tetap
// UNKNOWN.
var emas = map[string]emasKasus{
	"edm-fire-1.json": {cob: "FIRE", terpetakan: 2704, takDikenal: 3, penunjuk: 0,
		baris: map[string]int{"T_BUILDINGCONSTRUCTION": 1, "T_CEDINGCEDANTLIST": 1, "T_CEDING_CURRENCYLIST": 1, "T_COVERAGELIST": 45, "T_CURRENCYLIST": 1, "T_DEDUCTIBLELIST": 45, "T_GENERAL_POLIS": 1, "T_LISTCAUSEOFLOSS": 1, "T_LISTINSTALLMENT": 1, "T_LOCATIONLIST": 1, "T_OCCUPATIONLIST": 1, "T_PROPERTY": 1, "T_PROPERTYITEMLIST": 9, "T_QUOTATIONDATA": 1, "T_SPREADINGLIST": 90, "T_SURROUNDINGRISK": 1, "T_TABLEOFLIMIT": 1, "T_WORK_POLIS": 1},
		uang:  map[string]int{"T_COVERAGELIST.CURRENCY_CODE": 45, "T_DEDUCTIBLELIST.CURRENCY": 5, "T_PROPERTY.CURRENCY_CODE": 1, "T_SPREADINGLIST.CURRENCY_CODE": 90}},
	"nb-fire-1.json": {cob: "FIRE", terpetakan: 804, takDikenal: 10, penunjuk: 0,
		baris: map[string]int{"T_BUILDINGCONSTRUCTION": 1, "T_CEDINGCEDANTLIST": 1, "T_CEDINGCOLIST": 1, "T_CEDING_CURRENCYLIST": 1, "T_COVERAGELIST": 5, "T_CURRENCYLIST": 1, "T_DATASCORINGRISKLIST": 2, "T_DEDUCTIBLELIST": 7, "T_GENERAL_POLIS": 1, "T_LISTCAUSEOFLOSS": 1, "T_LISTINSTALLMENT": 4, "T_LOCATIONLIST": 1, "T_OCCUPATIONLIST": 1, "T_PROPERTY": 1, "T_PROPERTYITEMLIST": 1, "T_QUOTATIONDATA": 1, "T_SCORINGRESULT": 1, "T_SCORINGRISK": 1, "T_SCORING_FACTOR": 48, "T_SCORING_OPTION": 202, "T_SPREADINGLIST": 5, "T_SURROUNDINGRISK": 1, "T_TABLEOFLIMIT": 1, "T_WORK_POLIS": 1},
		uang:  map[string]int{"T_COVERAGELIST.CURRENCY_CODE": 5, "T_DEDUCTIBLELIST.CURRENCY": 2, "T_PROPERTY.CURRENCY_CODE": 1, "T_SPREADINGLIST.CURRENCY_CODE": 5}},
	"nb-kredit-1.json": {cob: "Aneka", terpetakan: 126, takDikenal: 9, penunjuk: 0,
		baris: map[string]int{"T_ANEKALIST": 1, "T_CEDINGCEDANTLIST": 1, "T_CEDINGCOLIST": 1, "T_CEDING_CURRENCYLIST": 1, "T_COVERAGELIST": 1, "T_CURRENCY": 2, "T_CURRENCYLIST": 1, "T_GENERAL_POLIS": 1, "T_LISTINSTALLMENT": 1, "T_LOCATIONLIST": 1, "T_OCCUPATIONLIST": 1, "T_PROPERTY": 1, "T_QUOTATIONDATA": 1, "T_RISKLOCATION": 1, "T_SPREADINGLIST": 1, "T_WORK_POLIS": 1},
		uang:  map[string]int{"T_PROPERTY.CURRENCY_CODE": 1}},
	"nb-marinecargo-1.json": {cob: "MarineCargo", terpetakan: 4842, takDikenal: 2, penunjuk: 0,
		baris: map[string]int{"T_ADDITIONALSHIP": 104, "T_CARGOLIST": 104, "T_CEDINGCEDANTLIST": 1, "T_CEDING_CURRENCYLIST": 1, "T_COVERAGELIST": 104, "T_CURRENCY": 104, "T_CURRENCYLIST": 1, "T_DEDUCTIBLELIST": 104, "T_GENERAL_POLIS": 1, "T_LISTINSTALLMENT": 1, "T_QUOTATIONDATA": 1, "T_SHIP": 104, "T_SPREADINGLIST": 104, "T_WORK_POLIS": 1},
		uang:  map[string]int{}},
	"rnw-fire-1.json": {cob: "FIRE", terpetakan: 724, takDikenal: 10, penunjuk: 0,
		baris: map[string]int{"T_BUILDINGCONSTRUCTION": 1, "T_CEDINGCEDANTLIST": 1, "T_CEDINGCOLIST": 1, "T_CEDING_CURRENCYLIST": 1, "T_COVERAGELIST": 4, "T_CURRENCYLIST": 1, "T_DATASCORINGRISKLIST": 2, "T_DEDUCTIBLELIST": 7, "T_GENERAL_POLIS": 1, "T_LISTCAUSEOFLOSS": 1, "T_LISTINSTALLMENT": 1, "T_LOCATIONLIST": 1, "T_OCCUPATIONLIST": 1, "T_PROPERTY": 1, "T_PROPERTYITEMLIST": 1, "T_QUOTATIONDATA": 1, "T_SCORINGRESULT": 1, "T_SCORINGRISK": 1, "T_SCORING_FACTOR": 48, "T_SCORING_OPTION": 202, "T_SPREADINGLIST": 4, "T_SURROUNDINGRISK": 1, "T_TABLEOFLIMIT": 1, "T_WORK_POLIS": 1},
		uang:  map[string]int{"T_COVERAGELIST.CURRENCY_CODE": 4, "T_PROPERTY.CURRENCY_CODE": 1, "T_SPREADINGLIST.CURRENCY_CODE": 4}},
}

// grupBusinessType - lembar `Grup Bisnis` workbook asal (`Claude outputs\Tabel-Flat-
// per-Grup-Bisnis.xlsx`, kolom BUSINESSTYPE YANG MASUK): cara KEDUA menetapkan
// kelompok, dari NILAI - V-16 langkah 1-4 membaca BENTUK. Keduanya harus sepakat.
var grupBusinessType = map[string]string{
	"Aneka": "Aneka", "AviationHull": "Aneka", "Bonding": "Aneka", "BondingKBG": "Aneka", "CustomBond": "Aneka",
	"ElectronicEquipment": "Aneka", "GolfInsurance": "Aneka", "HE": "Aneka", "LandRig": "Aneka", "Liability": "Aneka",
	"MBD": "Aneka", "MarineHull": "Aneka", "FireStyle1": "FIRE", "FireStyle2": "FIRE", "Life": "Life",
	"MBUCar": "MBUCar", "MarineCargo": "MarineCargo", "PA": "PA",
}

func bacaKasus(t *testing.T) map[string][]byte {
	t.Helper()
	berkas, err := filepath.Glob(filepath.Join(folderKasus, "*.json"))
	if err != nil || len(berkas) != len(emas) {
		t.Fatalf("%d fixture (%v), mau %d", len(berkas), err, len(emas))
	}
	isi := map[string][]byte{}
	for _, f := range berkas {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		isi[filepath.Base(f)] = b
	}
	return isi
}

// TestKasusLolosKetat - temuan 02-10-2026: Flatten ketat berhenti di kelima kasus
// karena DDL draf memberi PPN_CHECK dan SHARE_OF_CEDING tipe NUMBER. Butir 68.1: kedua
// kolom itu (dan enam lainnya) kini TEKS APA ADANYA - penyimpangan sadar dari DDL.
// Uji ini menagih: kelima kasus lolos ketat, dan PPN_CHECK tersimpan sebagai teks,
// bukan angka.
func TestKasusLolosKetat(t *testing.T) {
	for nama, b := range bacaKasus(t) {
		h, err := Flatten(Masukan{IDPega: idUji, DataJSON: b})
		if err != nil {
			t.Errorf("%s: %v", nama, err)
			continue
		}
		v := h.Baris["T_GENERAL_POLIS"][0].Kolom["PPN_CHECK"]
		if v.Angka != nil || v.Teks == "" {
			t.Errorf("%s: PPN_CHECK bukan teks", nama)
		}
	}
}

// TestKasusUkur - mode ukur (konflik tipe dihitung, tidak menghentikan): baris per
// tabel, kelompok lini, konflik tipe, dan mata uang tak diketahui per kasus.
func TestKasusUkur(t *testing.T) {
	for nama, b := range bacaKasus(t) {
		e := emas[nama]
		ba := map[string]int{}
		h, err := flatten(Masukan{IDPega: idUji, DataJSON: b}, ba)
		if err != nil {
			t.Fatalf("%s: %v", nama, err)
		}
		if cob := teks(h.Baris["T_GENERAL_POLIS"][0], "COB_GROUP"); cob != e.cob {
			t.Errorf("%s: COB_GROUP %s, mau %s", nama, cob, e.cob)
		}
		if !reflect.DeepEqual(h.Diagnostik.BarisPerTabel, e.baris) {
			t.Errorf("%s: baris %v\nmau %v", nama, h.Diagnostik.BarisPerTabel, e.baris)
		}
		for tb, n := range e.baris {
			if len(h.Baris[tb]) != n {
				t.Errorf("%s: Baris[%s] %d, BarisPerTabel %d", nama, tb, len(h.Baris[tb]), n)
			}
		}
		if len(ba) != 0 {
			t.Errorf("%s: masih ada konflik tipe %v", nama, ba)
		}
		tak, tunjuk := 0, 0
		for _, m := range h.Penampung {
			if m.Penunjuk {
				tunjuk++
			} else {
				tak++
			}
		}
		// ADR-0023: penampung = setiap medan tak terpetakan dan penunjuk di luar 48 kolom
		// penunjuk, satu per medan. `penunjuk` emas 0: menagih bahwa tidak ada penunjuk
		// fixture yang tertinggal di penampung sesudah butir 72.
		if tak != e.takDikenal || tunjuk != e.penunjuk || tak != jumlah(h.Diagnostik.TakTerpetakan) ||
			tunjuk != jumlah(h.Diagnostik.PenunjukBelumDikonversi) {
			t.Errorf("%s: penampung %d/%d, mau %d/%d (Diagnostik %d/%d)", nama, tak, tunjuk, e.takDikenal, e.penunjuk,
				jumlah(h.Diagnostik.TakTerpetakan), jumlah(h.Diagnostik.PenunjukBelumDikonversi))
		}
		if !reflect.DeepEqual(h.Diagnostik.MataUangTakDiketahui, e.uang) {
			t.Errorf("%s: mata uang %v, mau %v", nama, h.Diagnostik.MataUangTakDiketahui, e.uang)
		}
		if h.Diagnostik.Terpetakan != e.terpetakan {
			t.Errorf("%s: Terpetakan %d, mau %d", nama, h.Diagnostik.Terpetakan, e.terpetakan)
		}
		// Cara kedua kelompok lini: dari BusinessType akar, lewat lembar Grup Bisnis.
		var d struct {
			QuotationData struct{ BusinessType string }
		}
		if err := json.Unmarshal(b, &d); err != nil || grupBusinessType[d.QuotationData.BusinessType] != e.cob {
			t.Errorf("%s: BusinessType -> %q, bentuk -> %s (%v)", nama, grupBusinessType[d.QuotationData.BusinessType], e.cob, err)
		}
	}
}

// TestSetiapMedanTerhitung - rekonsiliasi dua cara: pencacah daun INDEPENDEN
// (encoding/json biasa, bukan pohon mesin) menghitung medan terisi yang bukan
// px/pz/py; jumlah itu harus sama dengan Terpetakan + Dibuang (tanpa dua alasan
// metadata) + TakTerpetakan + PenunjukBelumDikonversi + konflik tipe. Tidak ada
// medan yang hilang tanpa terhitung, dan tidak ada yang terhitung dua kali.
func TestSetiapMedanTerhitung(t *testing.T) {
	for nama, b := range bacaKasus(t) {
		var v any
		if err := json.Unmarshal(b, &v); err != nil {
			t.Fatal(err)
		}
		mau := cacahDaun(v)
		ba := map[string]int{}
		h, err := flatten(Masukan{IDPega: idUji, DataJSON: b}, ba)
		if err != nil {
			t.Fatal(err)
		}
		got := h.Diagnostik.Terpetakan + jumlah(ba) + jumlah(h.Diagnostik.TakTerpetakan) + jumlah(h.Diagnostik.PenunjukBelumDikonversi)
		for alasan, n := range h.Diagnostik.Dibuang {
			if alasan != alasanMeta && alasan != alasanPy {
				got += n
			}
		}
		if got != mau {
			t.Errorf("%s: ember %d, pencacah independen %d", nama, got, mau)
		}
	}
}

func cacahDaun(v any) int {
	switch x := v.(type) {
	case map[string]any:
		n := 0
		for k, a := range x {
			if s, ok := a.(string); ok {
				if s != "" && !strings.HasPrefix(k, "px") && !strings.HasPrefix(k, "pz") && !strings.HasPrefix(k, "py") {
					n++
				}
				continue
			}
			n += cacahDaun(a)
		}
		return n
	case []any:
		n := 0
		for _, a := range x {
			if s, ok := a.(string); ok {
				if s != "" {
					n++
				}
				continue
			}
			n += cacahDaun(a)
		}
		return n
	}
	return 0
}

func jumlah(m map[string]int) int {
	n := 0
	for _, v := range m {
		n += v
	}
	return n
}

// TestKasusPenunjukBukanPosisi - bukti butir 71-72 dikunci (hitungan; satu nilai pointer
// dikutip): `IndexCargo` di kasus marine tersimpan APA ADANYA di INDEX_CARGO - teks "7"
// di ke-104 baris - dan sama dengan SEQ_NO baris induknya hanya di 1 dari 104 (A68,
// membantah premis V-47); R1 indeks-diri di kelima kasus sama dengan SEQ_NO sendiri di
// 227 dari 227 (A67 - yang berbeda hanya di korpus, 39). Dihitung dua cara 02-10-2026:
// mesin (uji ini) dan pengurai Python atas pohon mentah.
func TestKasusPenunjukBukanPosisi(t *testing.T) {
	r1 := map[string]string{} // tabel -> kolom R1
	for _, k := range amandemenPenunjuk {
		if kat := penunjukKandidat[k.tabel+"."+k.medan]; kat.kategori == "R1 indeks-diri" {
			r1[k.tabel] = k.nama
		}
	}
	samaR1, bedaR1, cargo7, cargoSama, cargoBeda := 0, 0, 0, 0, 0
	for _, b := range bacaKasus(t) {
		h, err := Flatten(Masukan{IDPega: idUji, DataJSON: b})
		if err != nil {
			t.Fatal(err)
		}
		kunci := map[int]Baris{}
		for _, bs := range h.Baris {
			for _, r := range bs {
				kunci[r.Kunci] = r
			}
		}
		for tb, bs := range h.Baris {
			for _, r := range bs {
				if kol, ada := r1[tb]; ada {
					if v, isi := r.Kolom[kol]; isi {
						if v.Teks == r.Kolom["SEQ_NO"].Angka.String() {
							samaR1++
						} else {
							bedaR1++
						}
					}
				}
				if v, isi := r.Kolom["INDEX_CARGO"]; isi && tb == "T_COVERAGELIST" {
					if v.Teks == "7" && v.Angka == nil {
						cargo7++
					}
					if v.Teks == kunci[r.Induk].Kolom["SEQ_NO"].Angka.String() {
						cargoSama++
					} else {
						cargoBeda++
					}
				}
			}
		}
	}
	if samaR1 != 227 || bedaR1 != 0 {
		t.Errorf("R1 indeks-diri sama %d beda %d, mau 227 / 0", samaR1, bedaR1)
	}
	if cargo7 != 104 || cargoSama != 1 || cargoBeda != 103 {
		t.Errorf("INDEX_CARGO teks 7: %d, sama induk %d, beda %d; mau 104 / 1 / 103", cargo7, cargoSama, cargoBeda)
	}
}
