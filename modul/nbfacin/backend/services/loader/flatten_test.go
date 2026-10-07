package loader

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// Kasus di berkas ini UJI-rekaan: dokumen disusun dari NAMA medan rancangan dengan
// nilai buatan; tidak satu pun nilai berasal dari kasus nyata. Kasus nyata ada di
// kasus_test.go (fixture de-identifikasi tiket 15).

const idUji = "UJI-KELAS NB-1"

type dok = map[string]any

func ratakan(t *testing.T, d dok) *Hasil {
	t.Helper()
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	h, err := Flatten(Masukan{IDPega: idUji, DataJSON: b})
	if err != nil {
		t.Fatalf("Flatten: %v", err)
	}
	return h
}

func galatRata(t *testing.T, d any) error {
	t.Helper()
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	h, err := Flatten(Masukan{IDPega: idUji, DataJSON: b})
	if err == nil {
		t.Fatalf("Flatten tidak berhenti, %d tabel", len(h.Baris))
	}
	if h != nil {
		t.Error("berhenti keras tetapi mengembalikan Hasil sebagian")
	}
	return err
}

func satu(t *testing.T, h *Hasil, tabel string) Baris {
	t.Helper()
	if n := len(h.Baris[tabel]); n != 1 {
		t.Fatalf("%s: %d baris, mau 1", tabel, n)
	}
	return h.Baris[tabel][0]
}

func teks(b Baris, kol string) string { return b.Kolom[kol].Teks }

func angka(b Baris, kol string) string {
	if b.Kolom[kol].Angka == nil {
		return "<nil>"
	}
	return b.Kolom[kol].Angka.String()
}

var polaUUID5 = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// Spec 11 uji 1 - satu penawaran sederhana: bentuk dasar baris dan kolom sistem.
func TestFlattenBentukDasar(t *testing.T) {
	h := ratakan(t, dok{
		"pxObjClass":    "UJI",
		"CurrentYear":   "2099",
		"QuotationData": dok{"BusinessType": "Life", "BusinessCode": "UJI-01"},
		"PersonList": []any{
			dok{"Age": "30", "CoverageList": []any{dok{"Coverage": "UJI-A"}, dok{"Coverage": "UJI-B"}}},
			dok{"Age": "40", "CoverageList": []any{dok{"Coverage": "UJI-C"}}},
		},
	})
	kerja, umum := satu(t, h, "T_WORK_POLIS"), satu(t, h, "T_GENERAL_POLIS")
	if kerja.Kunci != 1 || kerja.Induk != 0 || umum.Kunci != 2 || umum.Induk != 1 {
		t.Errorf("kunci kerja %d/%d, umum %d/%d", kerja.Kunci, kerja.Induk, umum.Kunci, umum.Induk)
	}
	if teks(kerja, "JENIS_WORK") != "NB" || teks(kerja, "NO_WORK") != "NB-1" || teks(kerja, "IDPEGA") != idUji {
		t.Errorf("T_WORK_POLIS %v", kerja.Kolom)
	}
	if _, ada := kerja.Kolom["COB_GROUP"]; ada {
		t.Error("T_WORK_POLIS tidak punya COB_GROUP di rancangan")
	}
	if teks(umum, "COB_GROUP") != "Life" || teks(umum, "CURRENT_YEAR") != "2099" {
		t.Errorf("T_GENERAL_POLIS %v", umum.Kolom)
	}
	q := satu(t, h, "T_QUOTATIONDATA")
	if q.Induk != umum.Kunci || teks(q, "BUSINESS_CODE") != "UJI-01" {
		t.Errorf("T_QUOTATIONDATA %+v", q)
	}
	orang, jam := h.Baris["T_PERSONLIST"], h.Baris["T_COVERAGELIST"]
	if len(orang) != 2 || len(jam) != 3 {
		t.Fatalf("%d orang, %d jaminan", len(orang), len(jam))
	}
	mauInduk := []int{orang[0].Kunci, orang[0].Kunci, orang[1].Kunci}
	mauSeq := []string{"1", "2", "1"}
	uid := map[string]bool{}
	for i, c := range jam {
		if c.Induk != mauInduk[i] || angka(c, "SEQ_NO") != mauSeq[i] {
			t.Errorf("jaminan %d: induk %d seq %s", i, c.Induk, angka(c, "SEQ_NO"))
		}
		if teks(c, "PARENT_TABLE") != "T_PERSONLIST" || teks(c, "SRC_PATH") != "PersonList/CoverageList" {
			t.Errorf("jaminan %d: PARENT_TABLE %q SRC_PATH %q", i, teks(c, "PARENT_TABLE"), teks(c, "SRC_PATH"))
		}
		if !polaUUID5.MatchString(teks(c, "ROW_UID")) || uid[teks(c, "ROW_UID")] {
			t.Errorf("jaminan %d: ROW_UID %q tidak berbentuk UUID v5 unik", i, teks(c, "ROW_UID"))
		}
		uid[teks(c, "ROW_UID")] = true
		if teks(c, "IDPEGA") != idUji || teks(c, "COB_GROUP") != "Life" {
			t.Errorf("jaminan %d: IDPEGA/COB_GROUP %v", i, c.Kolom)
		}
		if _, ada := c.Kolom["ID"]; ada {
			t.Error("Flatten tidak mengisi ID - itu milik repository")
		}
	}
	// Kunci = urutan lahir: induk selalu lebih kecil dari anaknya.
	for tb, bs := range h.Baris {
		for _, b := range bs {
			if b.Induk >= b.Kunci {
				t.Errorf("%s kunci %d induk %d", tb, b.Kunci, b.Induk)
			}
		}
	}
	// CurrentYear 1 + BusinessType 1 + BusinessCode 1 + Age 2 + Coverage 3 = 8.
	if h.Diagnostik.Terpetakan != 8 {
		t.Errorf("Terpetakan %d, mau 8", h.Diagnostik.Terpetakan)
	}
}

// Spec 11 uji 2, 3, 4 - SETIAP jalur rancangan (148, kecuali dua ruas sintetis
// V-30 yang diuji TestFlattenSkoring) menjadi satu dokumen minimal: larik untuk
// ruas bertabel berulang, objek selainnya, satu medan teks di ujungnya. Yang
// ditagih: rantai baris persis sesuai Jalur Sumber, PARENT_TABLE = induk tertulis
// (termasuk kelima induk T_COVERAGELIST), SRC_PATH = jalur rancangan, kedalaman
// sampai 8, dan tabel wadah murni tetap lahir berisi kolom sistem saja.
func TestFlattenSetiapJalur(t *testing.T) {
	tabelJalur := map[string]string{}
	for _, j := range jalurSumber {
		tabelJalur[j.jalur] = j.tabel
	}
	diuji, wadah := 0, map[string]bool{}
	indukTercapai := map[string]map[string]bool{}
	for _, j := range jalurSumber {
		if j.jalur == "" || strings.Contains(j.jalur, ruasFaktor) {
			continue
		}
		medan, nilai := medanUji(j.tabel)
		if medan == "" {
			continue // tabel wadah: tercakup sebagai leluhur jalur lain
		}
		ruas := strings.Split(j.jalur, "/")
		var rantai []string
		var daun any = dok{medan: nilai}
		for i := len(ruas) - 1; i >= 0; i-- {
			tb := tabelJalur[strings.Join(ruas[:i+1], "/")]
			rantai = append([]string{tb}, rantai...)
			if punyaKolom(tb, kolomSeq) {
				daun = []any{daun}
			}
			daun = dok{ruas[i]: daun}
		}
		d := daun.(dok)
		if q, ada := d["QuotationData"].(dok); ada {
			q["BusinessType"] = "PA"
		} else {
			d["QuotationData"] = dok{"BusinessType": "PA"}
		}
		h := ratakan(t, d)
		diuji++
		mau := map[string]int{"T_WORK_POLIS": 1, "T_GENERAL_POLIS": 1, "T_QUOTATIONDATA": 1}
		for _, tb := range rantai {
			if tb != "T_QUOTATIONDATA" {
				mau[tb]++
			}
		}
		got := map[string]int{}
		for tb, bs := range h.Baris {
			got[tb] = len(bs)
		}
		if !reflect.DeepEqual(got, mau) {
			t.Errorf("%s: baris %v, mau %v", j.jalur, got, mau)
			continue
		}
		// Setiap kolom NOT NULL yang bukan milik repository terisi (ID/PARENT_ID dari repository).
		for tb, bs := range h.Baris {
			for _, b := range bs {
				for _, k := range skemaTabel[tb] {
					if _, ada := b.Kolom[k.nama]; k.wajib && !ada {
						if asal, _ := asalKolom(tb, k); asal != asalRepository {
							t.Errorf("%s: %s.%s NOT NULL tidak terisi", j.jalur, tb, k.nama)
						}
					}
				}
			}
		}
		ujung := h.Baris[j.tabel][len(h.Baris[j.tabel])-1]
		if punyaKolom(j.tabel, kolomTabelInduk) && teks(ujung, kolomTabelInduk) != j.induk {
			t.Errorf("%s: PARENT_TABLE %q, mau %s", j.jalur, teks(ujung, kolomTabelInduk), j.induk)
		}
		if punyaKolom(j.tabel, kolomJalurSumber) && teks(ujung, kolomJalurSumber) != j.jalur {
			t.Errorf("%s: SRC_PATH %q", j.jalur, teks(ujung, kolomJalurSumber))
		}
		if v := ujung.Kolom[medanKolom(j.tabel, medan)]; v.Teks != nilai && (v.Angka == nil || v.Angka.String() != nilai) {
			t.Errorf("%s: medan %s tidak tertulis", j.jalur, medan)
		}
		if indukTercapai[j.tabel] == nil {
			indukTercapai[j.tabel] = map[string]bool{}
		}
		indukTercapai[j.tabel][j.induk] = true
		for _, tb := range rantai {
			if m, _ := medanUji(tb); m == "" {
				wadah[tb] = true
				for k := range h.Baris[tb][0].Kolom {
					if asal, _ := asalKolom(tb, kolomBernama(tb, k)); asal == asalMedan {
						t.Errorf("wadah %s berisi kolom medan %s", tb, k)
					}
				}
			}
		}
	}
	// 148 + 2 amandemen (T_ADDITIONALSHIP; T_FEALIST tiket 41) - 1 jalur akar - 2 ruas V-30 - 4 jalur tiga
	// tabel wadah (T_FR_POLICY dua jalur).
	if diuji != 143 {
		t.Errorf("%d jalur diuji, mau 143", diuji)
	}
	if len(wadah) != 3 {
		t.Errorf("tabel wadah tercapai %v, mau T_FR_FACOFFERLIST, T_FR_OBJECT, T_FR_POLICY", wadah)
	}
	if len(indukTercapai["T_COVERAGELIST"]) != 5 {
		t.Errorf("induk T_COVERAGELIST tercapai %v, mau 5", indukTercapai["T_COVERAGELIST"])
	}
}

// medanUji - FIELD ASLI pertama (urut abjad, teks lebih dulu) yang terisi dari medan
// dokumen, beserta nilai ujinya; "" bila tabel itu wadah murni.
func medanUji(tb string) (medan, nilai string) {
	var tks, ang []string
	for _, k := range skemaTabel[tb] {
		if asal, _ := asalKolom(tb, k); asal == asalMedan && !strings.HasPrefix(k.nama, "PAY_") {
			if k.angka() {
				ang = append(ang, k.medan)
			} else {
				tks = append(tks, k.medan)
			}
		}
	}
	sort.Strings(tks)
	sort.Strings(ang)
	switch {
	case len(tks) > 0:
		return tks[0], "UJI"
	case len(ang) > 0:
		return ang[0], "1"
	}
	return "", ""
}

func medanKolom(tb, medan string) string {
	for _, k := range skemaTabel[tb] {
		if k.medan == medan && !k.turunan {
			return k.nama
		}
	}
	return ""
}

func kolomBernama(tb, nama string) kolomSkema {
	for _, k := range skemaTabel[tb] {
		if k.nama == nama {
			return k
		}
	}
	return kolomSkema{}
}

// Spec 11 uji 5 - tiap kelompok lini bisnis V-16, termasuk BusinessType ganda.
func TestFlattenLiniBisnis(t *testing.T) {
	lokasiFire := []any{dok{"Property": dok{"PropertyItemList": []any{dok{"ItemType": "UJI"}}}}}
	lokasiAneka := []any{dok{"Property": dok{"RiskLocation": dok{"OccupationList": []any{dok{"AnekaList": []any{dok{"Quantity": "1"}}}}}}}}
	for _, k := range []struct {
		nama string
		d    dok
		mau  string
	}{
		{"FIRE langkah 1", dok{"LocationList": lokasiFire}, "FIRE"},
		{"Aneka langkah 2", dok{"LocationList": lokasiAneka}, "Aneka"},
		{"FIRE mendahului Aneka", dok{"LocationList": append(append([]any{}, lokasiAneka...), lokasiFire...)}, "FIRE"},
		{"MBUCar langkah 3", dok{"VehicleList": []any{dok{"Brand": "UJI"}}}, "MBUCar"},
		{"MarineCargo langkah 4", dok{"CargoList": []any{dok{"FromRute": "UJI"}}}, "MarineCargo"},
		{"Life langkah 5", dok{"QuotationData": dok{"BusinessType": "Life"}}, "Life"},
		{"PA langkah 5", dok{"QuotationData": dok{"BusinessType": "PA"}}, "PA"},
		// V-16 diperbaiki: AnekaList di cabang FacRetroList tidak menjadikan Aneka.
		{"AnekaList Fac Retro bukan Aneka", dok{"QuotationData": dok{"BusinessType": "PA"},
			"FacRetroList": []any{dok{"LocationList": lokasiAneka}}}, "PA"},
		// LocationList "bisnis": AnekaList yang hanya ada di cabang dibuang tidak dihitung.
		{"AnekaList di cabang OldData bukan Aneka", dok{"QuotationData": dok{"BusinessType": "Life"},
			"LocationList": []any{dok{"OldData": dok{"AnekaList": []any{dok{"Quantity": "1"}}}}}}, "Life"},
		// "ada" = larik berisi; larik kosong tidak menggeser langkah.
		{"larik kosong", dok{"QuotationData": dok{"BusinessType": "Life"}, "VehicleList": []any{}, "CargoList": []any{}}, "Life"},
		// BusinessType ganda: yang mengikat medan akar; salinan OldData dibuang V-19.
		{"BusinessType ganda", dok{"QuotationData": dok{"BusinessType": "Life"},
			"OldData": dok{"QuotationData": dok{"BusinessType": "PA"}}}, "Life"},
	} {
		h := ratakan(t, k.d)
		if got := teks(satu(t, h, "T_GENERAL_POLIS"), "COB_GROUP"); got != k.mau {
			t.Errorf("%s: COB_GROUP %q, mau %s", k.nama, got, k.mau)
		}
	}
	for nama, d := range map[string]dok{
		"tanpa BusinessType":                      {},
		"BusinessType tak dikenal":                {"QuotationData": dok{"BusinessType": "UJI-TAKDIKENAL"}},
		"dikenal tapi bukan Life/PA di langkah 5": {"QuotationData": dok{"BusinessType": "Aneka"}},
	} {
		if err := galatRata(t, d); !errors.Is(err, ErrLiniBisnis) {
			t.Errorf("%s: %v", nama, err)
		}
	}
}

// Spec 11 uji 6 - baris tanpa mata uang: UNKNOWN terpasang dan terhitung; pewarisan
// K-063 (c) membawa kode leluhur, termasuk UNKNOWN-nya.
func TestFlattenMataUang(t *testing.T) {
	bayar := dok{"ListInstallment": []any{dok{"InstallmentNo": "1"}}, "Premium": "10.5"}
	h := ratakan(t, dok{
		"QuotationData": dok{"BusinessType": "Life"},
		"CurrencyList": []any{
			dok{"Name": "USD", "OldID": "02", "Policy": dok{"Payment": bayar}},
			dok{"OldID": "01", "Policy": dok{"Payment": bayar}},
		},
		"LocationList": []any{dok{"CurrencyID": "10001", "Property": dok{"ObjectNo": "1",
			"PropertyItemList": []any{dok{"ItemType": "UJI", "CoverageList": []any{dok{"Coverage": "UJI",
				"DeductibleList": []any{dok{"Currency": "IDR", "Amount": "1"}, dok{"Amount": "2"}}}}}}}}},
	})
	uang := h.Baris["T_CURRENCYLIST"]
	if teks(uang[0], "CURRENCY_CODE") != "USD" || teks(uang[0], "NAME") != "USD" || angka(uang[0], "PAY_PREMIUM") != "10.5" {
		t.Errorf("CurrencyList 1 %v", uang[0].Kolom)
	}
	if teks(uang[1], "CURRENCY_CODE") != kodeTakDiketahui {
		t.Errorf("CurrencyList tanpa Name: %q", teks(uang[1], "CURRENCY_CODE"))
	}
	cicil := h.Baris["T_LISTINSTALLMENT"]
	if len(cicil) != 2 || teks(cicil[0], "CURRENCY_CODE") != "USD" || teks(cicil[1], "CURRENCY_CODE") != kodeTakDiketahui {
		t.Fatalf("ListInstallment %v", cicil)
	}
	if cicil[0].Induk != uang[0].Kunci || teks(cicil[0], "SRC_PATH") != "" && teks(cicil[0], "SRC_PATH") != "CurrencyList/ListInstallment" {
		t.Errorf("V-22a ListInstallment induk %d, mau %d", cicil[0].Induk, uang[0].Kunci)
	}
	// Coverage FIRE ini tanpa halaman Currency -> UNKNOWN (dihitung di bawah).
	potong := h.Baris["T_DEDUCTIBLELIST"]
	if teks(potong[0], "CURRENCY") != "IDR" || teks(potong[1], "CURRENCY") != kodeTakDiketahui {
		t.Errorf("DeductibleList %v / %v", potong[0].Kolom["CURRENCY"], potong[1].Kolom["CURRENCY"])
	}
	// T_PROPERTY mewarisi dari T_LOCATIONLIST, yang hanya punya CURRENCY_ID angka.
	if teks(satu(t, h, "T_PROPERTY"), "CURRENCY_CODE") != kodeTakDiketahui {
		t.Error("T_PROPERTY.CURRENCY_CODE bukan UNKNOWN")
	}
	mau := map[string]int{"T_CURRENCYLIST.CURRENCY_CODE": 1, "T_LISTINSTALLMENT.CURRENCY_CODE": 1,
		"T_DEDUCTIBLELIST.CURRENCY": 1, "T_PROPERTY.CURRENCY_CODE": 1, "T_COVERAGELIST.CURRENCY_CODE": 1,
		"T_SPREADINGLIST.CURRENCY_CODE": 0}
	for k, n := range mau {
		if h.Diagnostik.MataUangTakDiketahui[k] != n {
			t.Errorf("MataUangTakDiketahui[%s] = %d, mau %d", k, h.Diagnostik.MataUangTakDiketahui[k], n)
		}
	}
}

// Butir 68.3 + 69 - CURRENCY_CODE T_COVERAGELIST / T_ANEKALIST dan kembaran T_FR_*-nya
// dari halaman Currency/Name anaknya, SEBELUM pewarisan, sehingga turunannya ikut
// berkode; kode sendiri yang berbeda dari anaknya berhenti keras (A66).
func TestFlattenKodeDariHalamanCurrency(t *testing.T) {
	jaminan := dok{"Coverage": "UJI", "SpreadingList": []any{dok{"FlagSpreading": "UJI"}},
		"Currency": dok{"ID": "10026", "Name": "IDR"}}
	h := ratakan(t, dok{"CargoList": []any{dok{"FromRute": "UJI", "CoverageList": []any{jaminan}}},
		"FacRetroList": []any{dok{"CargoList": []any{dok{"FromRute": "UJI", "CoverageList": []any{jaminan}}}}}})
	if got := teks(satu(t, h, "T_COVERAGELIST"), "CURRENCY_CODE"); got != "IDR" {
		t.Errorf("T_COVERAGELIST.CURRENCY_CODE %q", got)
	}
	if got := teks(satu(t, h, "T_SPREADINGLIST"), "CURRENCY_CODE"); got != "IDR" {
		t.Errorf("T_SPREADINGLIST mewarisi %q", got)
	}
	// Butir 69: kembaran FR ikut 68.3 (anak Currency-nya baris T_FR_CURRENCY).
	if got := teks(satu(t, h, "T_FR_COVERAGELIST"), "CURRENCY_CODE"); got != "IDR" {
		t.Errorf("T_FR_COVERAGELIST.CURRENCY_CODE %q", got)
	}
	if got := teks(satu(t, h, "T_FR_SPREADINGLIST"), "CURRENCY_CODE"); got != "IDR" {
		t.Errorf("T_FR_SPREADINGLIST mewarisi %q", got)
	}
	// Aneka dan kembaran FR-nya; T_FR_FACOUTOBJECTLIST mewarisi dari T_FR_COVERAGELIST.
	aneka := func(kode string) []any {
		return []any{dok{"Property": dok{"RiskLocation": dok{"OccupationList": []any{dok{"AnekaList": []any{dok{"Quantity": "1",
			"Currency": dok{"Name": kode}, "CoverageList": []any{dok{"Coverage": "UJI", "Currency": dok{"Name": kode},
				"FacOutObjectList": []any{dok{"Rate": "1"}}}}}}}}}}}}
	}
	h = ratakan(t, dok{"LocationList": aneka("USD"), "FacRetroList": []any{dok{"LocationList": aneka("JPY")}}})
	for tb, mau := range map[string]string{"T_ANEKALIST": "USD", "T_FR_ANEKALIST": "JPY", "T_FR_COVERAGELIST": "JPY", "T_FR_FACOUTOBJECTLIST": "JPY"} {
		if got := teks(satu(t, h, tb), "CURRENCY_CODE"); got != mau {
			t.Errorf("%s.CURRENCY_CODE %q, mau %s", tb, got, mau)
		}
	}
	err := galatRata(t, dok{"CargoList": []any{dok{"FromRute": "UJI", "CoverageList": []any{dok{"Name": "USD",
		"Currency": dok{"Name": "IDR"}}}}}})
	if !errors.Is(err, ErrKolomGanda) || strings.Contains(err.Error(), "IDR") {
		t.Errorf("kode sendiri berbeda dari Currency/Name: %v", err)
	}
	// Dua halaman Currency berkode berbeda di bawah satu coverage: tidak ditebak.
	err = galatRata(t, dok{"CargoList": []any{dok{"FromRute": "UJI", "CoverageList": []any{dok{"Coverage": "UJI",
		"Currency": []any{dok{"Name": "IDR"}, dok{"Name": "USD"}}}}}}})
	if !errors.Is(err, ErrKolomGanda) {
		t.Errorf("Currency ganda berbeda: %v", err)
	}
	// Kode sendiri SAMA dengan anaknya: bukan konflik.
	h = ratakan(t, dok{"CargoList": []any{dok{"FromRute": "UJI", "CoverageList": []any{dok{"Name": "IDR",
		"Currency": dok{"Name": "IDR"}}}}}})
	if got := teks(satu(t, h, "T_COVERAGELIST"), "CURRENCY_CODE"); got != "IDR" {
		t.Errorf("kode sama: %q", got)
	}
}

// Spec 11 uji 7 - format angka (BAHAN §7): titik dan bulat apa adanya, koma desimal
// pengecualian yang terhitung; format lain berhenti keras tanpa menyebut nilainya.
func TestFlattenFormatAngka(t *testing.T) {
	// T_GENERAL_POLIS.CEDING_RETENTION (NUMBER). PCT_LIMIT tidak lagi dipakai di sini:
	// sejak butir 68.1 ia teks apa adanya (TestFlattenTeksMenyimpang).
	batas := func(v string) dok {
		return dok{"QuotationData": dok{"BusinessType": "Life"}, "CedingRetention": v}
	}
	for v, mau := range map[string]string{"70,5": "70.5", "70.5": "70.5", "70": "70", "-0,25": "-0.25",
		"0.10000000000000000000": "0.10000000000000000000"} {
		h := ratakan(t, batas(v))
		if got := angka(satu(t, h, "T_GENERAL_POLIS"), "CEDING_RETENTION"); got != mau {
			t.Errorf("%q -> %s, mau %s", v, got, mau)
		}
		koma := 0
		if strings.Contains(v, ",") {
			koma = 1
		}
		if h.Diagnostik.KomaDesimal["T_GENERAL_POLIS.CEDING_RETENTION"] != koma {
			t.Errorf("%q: KomaDesimal %d", v, h.Diagnostik.KomaDesimal["T_GENERAL_POLIS.CEDING_RETENTION"])
		}
	}
	for _, v := range []string{"70,5 ", "1.234,5", "1,234.5", "1e5", "true", "5%", ".5"} {
		err := galatRata(t, batas(v))
		if !errors.Is(err, ErrBukanAngka) || strings.Contains(err.Error(), strings.TrimSpace(v)) {
			t.Errorf("%q: %v", v, err)
		}
	}
	h := ratakan(t, batas("1234567890123456789012345678901234567890.5"))
	if h.Diagnostik.LebihDari38Digit["T_GENERAL_POLIS.CEDING_RETENTION"] != 1 {
		t.Error("angka 41 digit tidak terhitung LebihDari38Digit")
	}
	// Nol di ujung bukan digit bermakna (41 digit tertulis, 1 bermakna).
	h = ratakan(t, batas("1.0000000000000000000000000000000000000000"))
	if n := h.Diagnostik.LebihDari38Digit["T_GENERAL_POLIS.CEDING_RETENTION"]; n != 0 {
		t.Errorf("nol ujung terhitung > 38 digit: %d", n)
	}
	if got := angka(satu(t, h, "T_GENERAL_POLIS"), "CEDING_RETENTION"); got != "1.0000000000000000000000000000000000000000" {
		t.Errorf("nilai dibulatkan: %s", got)
	}
}

// Butir 68.1 - delapan kolom NUMBER berisi teks disimpan TEKS APA ADANYA: "%" tidak
// ditafsirkan, koma desimal tidak dinormalkan, spasi di ujung PctLimit tidak
// dipangkas ([terverifikasi] 57 nilai PctLimit di 115 contoh berakhiran spasi).
func TestFlattenTeksMenyimpang(t *testing.T) {
	// T_FR_COVERAGELIST tidak punya TYPE_OF_DISCOUNT - cabang FR tanpa CoverageList.
	batas := func(v string, jaminan []any) []any {
		return []any{dok{"Property": dok{"PropertyItemList": []any{dok{"ItemType": "UJI", "CoverageList": jaminan}},
			"OccupationList": []any{dok{"Category": "UJI", "TableOfLimit": []any{dok{"PctLimit": v}}}}}}}
	}
	h := ratakan(t, dok{"PPnCheck": "true", "QuotationData": dok{"BusinessType": "Life", "ShareOfCeding": "12.5%", "EdmChargeFee": "UJI"},
		"LocationList": batas("70,000 ", []any{dok{"TypeOfDiscount": "true"}}), "PersonList": []any{dok{"MasterRateCoverage": "UJI-A"}},
		"FacRetroList": []any{dok{"PrintRISlip": dok{"WarrPayment": "As Agreed"}, "LocationList": batas("12,5", nil)}}})
	// Kedelapan kolom, termasuk koma desimal TANPA spasi ("12,5") - koma tidak dinormalkan.
	for _, k := range []struct{ tabel, kolom, mau string }{
		{"T_GENERAL_POLIS", "PPN_CHECK", "true"}, {"T_QUOTATIONDATA", "SHARE_OF_CEDING", "12.5%"},
		{"T_QUOTATIONDATA", "EDM_CHARGE_FEE", "UJI"}, {"T_TABLEOFLIMIT", "PCT_LIMIT", "70,000 "},
		{"T_COVERAGELIST", "TYPE_OF_DISCOUNT", "true"}, {"T_PERSONLIST", "MASTER_RATE_COVERAGE", "UJI-A"},
		{"T_FR_PRINTRISLIP", "WARR_PAYMENT", "As Agreed"}, {"T_FR_TABLEOFLIMIT", "PCT_LIMIT", "12,5"},
	} {
		v := satu(t, h, k.tabel).Kolom[k.kolom]
		if v.Angka != nil || v.Teks != k.mau {
			t.Errorf("%s.%s = %+v, mau teks %q", k.tabel, k.kolom, v, k.mau)
		}
	}
	if n := len(h.Penampung); n != 0 {
		t.Errorf("penampung %d - medan uji tidak terpetakan: %+v", n, h.Penampung)
	}
	if len(h.Diagnostik.KomaDesimal) != 0 {
		t.Errorf("koma desimal ditafsirkan di kolom teks: %v", h.Diagnostik.KomaDesimal)
	}
	if len(kolomTeksMenyimpang) != 8 {
		t.Errorf("%d kolom menyimpang, mau 8", len(kolomTeksMenyimpang))
	}
	for tk := range kolomTeksMenyimpang {
		tb, kol, _ := strings.Cut(tk, ".")
		if k := kolomBernama(tb, kol); !k.angka() {
			t.Errorf("%s bukan NUMBER di DDL draf - penyimpangannya tidak perlu", tk)
		}
	}
}

// Amandemen butir 70: P1 CoverageInitial, P2-P3 T_ADDITIONALSHIP, P6 POLICY_TSI cabang FR.
func TestFlattenAmandemen70(t *testing.T) {
	h := ratakan(t, dok{"CargoList": []any{dok{"FromRute": "UJI",
		"CoverageList": []any{dok{"Coverage": "UJI", "CoverageInitial": "UJI A (Motor)"}},
		"PolicyData": dok{"Ship": dok{"NM_SHIP": "UJI",
			"AdditionalShip": []any{dok{"ID": "123", "DWT": "1", "GRT": "2", "NRT": "3"}, dok{"DWT": "4"}}}}}},
		"FacRetroList": []any{dok{"PrintRISlip": dok{"FacOfferList": []any{dok{"CurrencyList": []any{dok{"Name": "IDR",
			"TSI": "100", "Policy": dok{"TSI": "99.5", "Payment": dok{"TSITotal": "100"}}}}}}}}}})
	if got := teks(satu(t, h, "T_COVERAGELIST"), "COVERAGE_INITIAL"); got != "UJI A (Motor)" {
		t.Errorf("P1 COVERAGE_INITIAL %q", got)
	}
	kapal := satu(t, h, "T_SHIP")
	tambahan := h.Baris["T_ADDITIONALSHIP"]
	if len(tambahan) != 2 || tambahan[0].Induk != kapal.Kunci || teks(tambahan[0], "ADDITIONAL_SHIP_REF_ID") != "123" ||
		teks(tambahan[0], "DWT") != "1" || teks(tambahan[0], "NRT") != "3" || angka(tambahan[1], "SEQ_NO") != "2" ||
		!polaUUID5.MatchString(teks(tambahan[1], "ROW_UID")) || teks(tambahan[1], "COB_GROUP") != "MarineCargo" {
		t.Errorf("P2-P3 T_ADDITIONALSHIP %+v", tambahan)
	}
	fr := satu(t, h, "T_FR_CURRENCYLIST")
	if angka(fr, "POLICY_TSI") != "99.5" || angka(fr, "TSI") != "100" || angka(fr, "PAY_TSI_TOTAL") != "100" {
		t.Errorf("P6 %v", fr.Kolom)
	}
	if len(h.Penampung) != 0 {
		t.Errorf("penampung %+v", h.Penampung)
	}
}

// Butir 72 - SATU aturan: setiap penunjuk Idx*/Index* yang termasuk 48 kunci disimpan
// apa adanya di kolom teks mentahnya (induk langsung, R1, leluhur, berselisih, IdxPerson);
// kolom FK V-47 tetap kosong; penunjuk di luar 48 kunci tetap ke penampung.
func TestFlattenPenunjukKeKolom(t *testing.T) {
	h := ratakan(t, dok{"LocationList": []any{dok{"Property": dok{"PropertyItemList": []any{dok{"ItemType": "UJI",
		"IndexPropertyItem": "01", // R1 indeks-diri -> INDEX_PROPERTY_ITEM, teks apa adanya
		"IndexProperty":     "1",  // R3 induk langsung -> INDEX_PROPERTY
		"CoverageList": []any{dok{"Coverage": "UJI",
			"IndexPropertyItem": "1", // R3 induk langsung
			"IndexProperty":     "3", // R3 leluhur
			"IndexUjiTak":       "1", // di luar 48 kunci -> penampung
		}}}}}}}})
	item, jaminan := satu(t, h, "T_PROPERTYITEMLIST"), satu(t, h, "T_COVERAGELIST")
	for _, k := range []struct {
		b        Baris
		kol, mau string
	}{{item, "INDEX_PROPERTY_ITEM", "01"}, {item, "INDEX_PROPERTY", "1"}, {jaminan, "INDEX_PROPERTY_ITEM", "1"}, {jaminan, "INDEX_PROPERTY", "3"}} {
		if v := k.b.Kolom[k.kol]; v.Teks != k.mau || v.Angka != nil {
			t.Errorf("%s = %+v, mau teks %q", k.kol, v, k.mau)
		}
	}
	// FK V-47 tidak diisi dari penunjuk (68.4).
	for _, kol := range []string{"PROPERTY_ID", "PROPERTY_ITEM_ID", "LOCATION_ID"} {
		if _, ada := jaminan.Kolom[kol]; ada {
			t.Errorf("FK %s terisi", kol)
		}
	}
	if len(h.Penampung) != 1 || h.Penampung[0].Medan != "IndexUjiTak" || len(h.Diagnostik.Dibuang) != 0 {
		t.Errorf("penampung %+v, Dibuang %v", h.Penampung, h.Diagnostik.Dibuang)
	}
	// Kunci berselisih dan IdxPerson juga berkolom.
	h = ratakan(t, dok{"QuotationData": dok{"BusinessType": "PA"}, "FacRetroList": []any{dok{
		"PersonList": []any{dok{"IdxPerson": "2"}},
		"LocationList": []any{dok{"Property": dok{"PropertyItemList": []any{dok{"ItemType": "UJI", "CoverageList": []any{dok{"Rate": "1",
			"DeductibleList": []any{dok{"Amount": "1", "IndexProperty": "4"}}}}}}}}}}}})
	if teks(satu(t, h, "T_FR_DEDUCTIBLELIST"), "INDEX_PROPERTY") != "4" || teks(satu(t, h, "T_FR_PERSONLIST"), "IDX_PERSON") != "2" || len(h.Penampung) != 0 {
		t.Errorf("berselisih/IdxPerson: penampung %+v", h.Penampung)
	}
}

// K-069 (7b), butir 69 - IsCedingConfirm tidak hilang bersama ViewSuggest (V-31):
// masuk penampung (tabel dan tipenya menunggu work owner); medan ViewSuggest lain tetap
// dibuang dan hanya mereka yang terhitung Dibuang V-31.
func TestFlattenIsCedingConfirmDiselamatkan(t *testing.T) {
	h := ratakan(t, dok{"QuotationData": dok{"BusinessType": "Life"},
		"ViewSuggest": []any{dok{"No": "1", "IsCedingConfirm": "Offer"}, dok{"No": "2", "IsCedingConfirm": ""}, dok{"No": "3", "IsCedingConfirm": "Binding"}}})
	umum := satu(t, h, "T_GENERAL_POLIS")
	mau := []MedanTakDikenal{
		{Kunci: umum.Kunci, Jalur: "ViewSuggest[1]", Medan: "IsCedingConfirm", Nilai: "Offer"},
		{Kunci: umum.Kunci, Jalur: "ViewSuggest[3]", Medan: "IsCedingConfirm", Nilai: "Binding"},
	}
	if !reflect.DeepEqual(h.Penampung, mau) {
		t.Errorf("penampung %+v", h.Penampung)
	}
	if n := h.Diagnostik.Dibuang[cabangDibuang["ViewSuggest"]]; n != 3 {
		t.Errorf("Dibuang V-31 %d, mau 3 (No x 3)", n)
	}
}

// Spec 11 uji 8, 9 - cabang yang dikecualikan tidak menghasilkan baris; medannya
// terhitung per alasan.
func TestFlattenCabangDibuang(t *testing.T) {
	jaminan := []any{dok{"Coverage": "UJI"}}
	h := ratakan(t, dok{
		"QuotationData": dok{"BusinessType": "PA", "pxObjClass": "UJI", "pyLabel": "UJI"},
		"OldData":       dok{"QuotationData": dok{"BusinessType": "Life"}, "PersonList": []any{dok{"CoverageList": jaminan}}},
		"Parameters":    dok{"CariOccupationOldID": "1"},
		"OutGoList":     []any{dok{"OutgoAmount": "0"}},
		"ViewSuggest":   []any{dok{"No": "1"}},
		"FacRetro":      dok{"PersonList": []any{dok{"CoverageList": jaminan}}},
		"PersonList": []any{dok{"Age": "1", "TSIOld": "5", "EDMOldPremi": "5",
			"CoverageList": []any{dok{"Coverage": "UJI", "PremiumOld": "1"}}}},
		"CurrencyList": []any{dok{"Name": "USD", "OldID": "02"}},
		"LocationList": []any{dok{"Property": dok{"TotalTSIList": []any{dok{"TSI": "1"}},
			"TotalTSIPremiGrossList": []any{dok{"TSI": "1"}}, "TotalTSIPremiSpreadRNM": []any{dok{"TSI": "1"}}}}},
		"FacRetroList": []any{dok{"PrintRISlip": dok{"FacOfferList": []any{dok{"OldOfferedPayment": dok{"X": "1"}}}}}},
	})
	if len(h.Baris["T_PERSONLIST"]) != 1 || len(h.Baris["T_COVERAGELIST"]) != 1 {
		t.Errorf("cabang salinan kerja ikut menjadi baris: %d orang, %d jaminan", len(h.Baris["T_PERSONLIST"]), len(h.Baris["T_COVERAGELIST"]))
	}
	for _, tb := range []string{"T_LOCATIONLIST", "T_PROPERTY", "T_FACRETROLIST", "T_FR_PRINTRISLIP", "T_FR_FACOFFERLIST"} {
		if len(h.Baris[tb]) != 0 {
			t.Errorf("%s lahir dari cabang yang seluruh isinya dibuang (V-27)", tb)
		}
	}
	mau := map[string]int{
		cabangDibuang["OldData"]: 2, cabangDibuang["Parameters"]: 1, cabangDibuang["OutGoList"]: 1,
		cabangDibuang["ViewSuggest"]: 1, cabangDibuang["FacRetro"]: 1, cabangDibuang["TotalTSIList"]: 3,
		cabangDibuang["OldOfferedPayment"]: 1, alasanSufiksOld: 2, alasanEDMOld: 1, alasanMeta: 1, alasanPy: 1,
	}
	// OutGoList dan Parameters berbagi alasan V-34.
	mau[cabangDibuang["Parameters"]] = 2
	if !reflect.DeepEqual(h.Diagnostik.Dibuang, mau) {
		t.Errorf("Dibuang %v\nmau %v", h.Diagnostik.Dibuang, mau)
	}
	// V-28b: yang BERAWALAN Old adalah pengenal dan tetap disimpan.
	if teks(satu(t, h, "T_CURRENCYLIST"), "OLD_ID") != "02" {
		t.Error("OldID berawalan Old (pengenal) ikut terbuang")
	}
	// V-27: lokasi yang isinya hanya tabel Total* tidak lahir, tetapi tercatat.
	if h.Diagnostik.CabangKosong["LocationList"] != 1 || h.Diagnostik.CabangKosong["FacRetroList"] != 1 {
		t.Errorf("CabangKosong %v", h.Diagnostik.CabangKosong)
	}
}

// Spec 11 uji 10 - setiap keadaan berhenti keras yang dapat dicapai lewat dokumen.
// Dua yang tak tercapai dengan Jalur Sumber sekarang - induk tak tetap dan kedalaman
// > 8 - dijaga invarian TestJalurIndukAdalahAwalan dan TestKedalamanRancangan.
func TestFlattenBerhentiKeras(t *testing.T) {
	for _, id := range []string{"", "NB-1", "UJI-KELAS ", "UJI-KELAS XX-1", "UJI-KELAS NB-", "UJI-KELAS NB1"} {
		if _, err := Flatten(Masukan{IDPega: id, DataJSON: []byte(`{}`)}); !errors.Is(err, ErrIDPega) {
			t.Errorf("IDPEGA %q: %v", id, err)
		}
	}
	// Pesan galat tidak menyalin isi dokumen (encoding/json mengutip karakter masukan).
	for _, s := range []string{``, `[]`, `"x"`, `{"a":1`, `{"a":"1","a":"2"}`, `{} {}`, `{"a":RAHASIAUJI}`} {
		_, err := Flatten(Masukan{IDPega: idUji, DataJSON: []byte(s)})
		if !errors.Is(err, ErrDokumen) || strings.Contains(err.Error(), "RAHASIA") || strings.Contains(err.Error(), "'R'") {
			t.Errorf("dokumen %q: %v", s, err)
		}
	}
	if err := galatRata(t, dok{"QuotationData": dok{"BusinessType": "MarineHull"}}); strings.Contains(err.Error(), "MarineHull") {
		t.Errorf("pesan galat menyalin nilai dokumen: %v", err)
	}
	life := dok{"BusinessType": "Life"}
	for nama, k := range map[string]struct {
		d   dok
		mau error
	}{
		"PolicyData menimpa medan akar": {dok{"QuotationData": life, "EndorsementNo": "1", "PolicyData": dok{"EndorsementNo": "2"}}, ErrKolomGanda},
		"PolicyData larik":              {dok{"QuotationData": life, "PolicyData": []any{dok{"EndorsementNo": "1"}}}, ErrBentuk},
		"unsur larik bukan objek":       {dok{"QuotationData": life, "PersonList": []any{"UJI"}}, ErrBentuk},
		"NUMBER bukan angka":            {dok{"QuotationData": life, "CedingRetention": "UJI"}, ErrBukanAngka},
	} {
		if err := galatRata(t, k.d); !errors.Is(err, k.mau) {
			t.Errorf("%s: %v, mau %v", nama, err, k.mau)
		}
	}
}

// Lipatan V-39/V-22/V-22b/V-24b/V-48 dan penamaan ulang V-49.
func TestFlattenLipatan(t *testing.T) {
	h := ratakan(t, dok{
		"ID":            "UJI-SUMBER",
		"IsFlagReject":  "false",
		"QuotationData": dok{"BusinessType": "Life"},
		"PolicyData":    dok{"EndorsementNo": "0", "Payment": dok{"Installment": "2", "RICommision": "1.5"}},
		"CurrencyList": []any{dok{"Name": "USD", "TSI": "100", "ID": "10001",
			"Policy": dok{"TSI": "999", "Payment": dok{"TSITotal": "100"}}}},
	})
	umum, kerja := satu(t, h, "T_GENERAL_POLIS"), satu(t, h, "T_WORK_POLIS")
	if teks(umum, "SOURCE_ID") != "UJI-SUMBER" || teks(umum, "ENDORSEMENT_NO") != "0" ||
		angka(umum, "PAY_INSTALLMENT") != "2" || angka(umum, "PAY_RI_COMMISION") != "1.5" {
		t.Errorf("T_GENERAL_POLIS %v", umum.Kolom)
	}
	if teks(kerja, "IS_FLAG_REJECT") != "false" {
		t.Errorf("V-48 IsFlagReject tidak ke T_WORK_POLIS: %v", kerja.Kolom)
	}
	uang := satu(t, h, "T_CURRENCYLIST")
	if angka(uang, "TSI") != "100" || angka(uang, "PAY_TSI_TOTAL") != "100" {
		t.Errorf("T_CURRENCYLIST %v", uang.Kolom)
	}
	// Policy/TSI akar tidak menimpa CurrencyList/TSI (P6 hanya untuk cabang FR) - terhitung.
	// CurrencyList/ID (angka rujukan) tidak menjadi CURRENCY_CODE; P5: ke CURRENCY_REF_ID.
	if h.Diagnostik.TakTerpetakan["CurrencyList/Policy :: TSI"] != 1 {
		t.Errorf("TakTerpetakan %v", h.Diagnostik.TakTerpetakan)
	}
	if teks(uang, "CURRENCY_REF_ID") != "10001" || teks(uang, "CURRENCY_CODE") != "USD" {
		t.Errorf("P5 CURRENCY_REF_ID %q, CURRENCY_CODE %q", teks(uang, "CURRENCY_REF_ID"), teks(uang, "CURRENCY_CODE"))
	}

	h = ratakan(t, dok{"CargoList": []any{dok{"FromRute": "UJI",
		"PolicyData": dok{"SailDate": "20990101", "Ship": dok{"ID": "UJI-KAPAL", "NM_SHIP": "UJI"}}}}})
	kargo, kapal := satu(t, h, "T_CARGOLIST"), satu(t, h, "T_SHIP")
	if teks(kargo, "SAIL_DATE") != "20990101" || kapal.Induk != kargo.Kunci || teks(kapal, "SHIP_REF_ID") != "UJI-KAPAL" {
		t.Errorf("V-24b kargo %v kapal %+v", kargo.Kolom, kapal)
	}
}

// V-30 - faktor skoring menjadi T_SCORING_FACTOR + T_SCORING_OPTION.
func TestFlattenSkoring(t *testing.T) {
	h := ratakan(t, dok{
		"LocationList": []any{dok{"Property": dok{"PropertyItemList": []any{dok{"ItemType": "UJI"}}}}},
		"ScoringRisk": dok{"FinalScore": "1", "DataScoringRiskList": []any{dok{"Score": "7",
			"FEA": dok{"pxObjClass": "UJI", "FireAlarmSystem": dok{"Checked": "true", "ScorePerFactor": "2",
				"ChechBox1": "true", "Score1": "1", "ChechBox2": "false", "Score2": "0"}},
			"Occupation": dok{"Checked": "false", "ChechBox1": "", "Score1": "3", "ChechBox10": "true"},
			// Faktor yang hanya membawa Checked (tanpa ChechBoxN) tetap faktor.
			"Others": dok{"Rate": dok{"Checked": "true"}},
		}}},
	})
	daftar := satu(t, h, "T_DATASCORINGRISKLIST")
	f := h.Baris["T_SCORING_FACTOR"]
	if len(f) != 3 || f[0].Induk != daftar.Kunci || teks(f[2], "FACTOR_GROUP") != "Others" || teks(f[2], "CHECKED") != "true" {
		t.Fatalf("faktor %v", f)
	}
	if teks(f[0], "FACTOR_GROUP") != "FEA" || teks(f[0], "FACTOR_NAME") != "FireAlarmSystem" || teks(f[0], "CHECKED") != "true" ||
		angka(f[0], "SCORE_PER_FACTOR") != "2" || angka(f[0], "SEQ_NO") != "1" {
		t.Errorf("faktor 1 %v", f[0].Kolom)
	}
	if _, ada := f[1].Kolom["FACTOR_GROUP"]; ada || teks(f[1], "FACTOR_NAME") != "Occupation" || angka(f[1], "SEQ_NO") != "2" {
		t.Errorf("faktor 2 %v", f[1].Kolom)
	}
	o := h.Baris["T_SCORING_OPTION"]
	type opsi struct {
		induk              int
		no, cek, skor, seq string
	}
	var got []opsi
	for _, b := range o {
		got = append(got, opsi{b.Induk, teks(b, "OPTION_NO"), teks(b, "CHECK_BOX"), angka(b, "SCORE"), angka(b, "SEQ_NO")})
	}
	// SEQ_NO = N (A50): opsi 2..9 Occupation tidak ada, celahnya dipertahankan.
	mau := []opsi{{f[0].Kunci, "1", "true", "1", "1"}, {f[0].Kunci, "2", "false", "0", "2"},
		{f[1].Kunci, "1", "", "3", "1"}, {f[1].Kunci, "10", "true", "<nil>", "10"}}
	if !reflect.DeepEqual(got, mau) {
		t.Errorf("opsi %v\nmau %v", got, mau)
	}
}

// V-17b / V-33 - ruas mentah yang bernama lain di rancangan.
func TestFlattenGantiNama(t *testing.T) {
	h := ratakan(t, dok{"VehicleList": []any{dok{"Brand": "UJI", "Occupation": dok{"OccupationName": "UJI"}}}})
	o := satu(t, h, "T_OCCUPATIONLIST")
	if teks(o, "SRC_PATH") != "VehicleList/OccupationList" || teks(o, "PARENT_TABLE") != "T_VEHICLELIST" {
		t.Errorf("V-17b %v", o.Kolom)
	}
	h = ratakan(t, dok{"QuotationData": dok{"BusinessType": "PA"}, "PersonList": []any{dok{"Age": "1",
		"ASMCoverage": []any{dok{"CalculateMethod": "UJI"}}}}})
	c := satu(t, h, "T_COVERAGELIST")
	if teks(c, "CALCULATE_METHOD") != "UJI" || teks(c, "SRC_PATH") != "PersonList/CoverageList" {
		t.Errorf("V-33 %v", c.Kolom)
	}
}

// Medan tanpa kolom tidak hilang diam-diam - terhitung menurut jalur dan nama.
func TestFlattenTakTerpetakan(t *testing.T) {
	h := ratakan(t, dok{"QuotationData": dok{"BusinessType": "Life"}, "UjiMedan": "1",
		"InwardScale": dok{"UjiA": "1", "UjiB": dok{"UjiC": "2"}, "UjiD": []any{"3", "", "4"}},
		"PersonList":  []any{dok{"Age": "1", "IdxPerson": "1"}}})
	// InwardScale: UjiA + UjiC + dua unsur larik bernilai terisi = 4.
	mau := map[string]int{"(akar) :: UjiMedan": 1, "InwardScale :: *": 4}
	if !reflect.DeepEqual(h.Diagnostik.TakTerpetakan, mau) {
		t.Errorf("TakTerpetakan %v", h.Diagnostik.TakTerpetakan)
	}
	if h.Diagnostik.PenunjukBelumDikonversi["PersonList :: IdxPerson"] != 1 {
		t.Errorf("PenunjukBelumDikonversi %v", h.Diagnostik.PenunjukBelumDikonversi)
	}
	// ADR-0023: tiap medan masuk penampung BESERTA nilainya dan kunci barisnya.
	orang := satu(t, h, "T_PERSONLIST")
	umum := satu(t, h, "T_GENERAL_POLIS")
	mauP := []MedanTakDikenal{
		{Kunci: umum.Kunci, Jalur: "InwardScale", Medan: "UjiA", Nilai: "1"},
		{Kunci: umum.Kunci, Jalur: "InwardScale/UjiB", Medan: "UjiC", Nilai: "2"},
		{Kunci: umum.Kunci, Jalur: "InwardScale/UjiD", Medan: "[1]", Nilai: "3"},
		{Kunci: umum.Kunci, Jalur: "InwardScale/UjiD", Medan: "[3]", Nilai: "4"},
		{Kunci: orang.Kunci, Jalur: "PersonList", Medan: "IdxPerson", Nilai: "1", Penunjuk: true},
		{Kunci: umum.Kunci, Jalur: "(akar)", Medan: "UjiMedan", Nilai: "1"},
	}
	sort.Slice(h.Penampung, func(i, j int) bool {
		return h.Penampung[i].Jalur+h.Penampung[i].Medan < h.Penampung[j].Jalur+h.Penampung[j].Medan
	})
	sort.Slice(mauP, func(i, j int) bool { return mauP[i].Jalur+mauP[i].Medan < mauP[j].Jalur+mauP[j].Medan })
	if !reflect.DeepEqual(h.Penampung, mauP) {
		t.Errorf("penampung %+v\nmau %+v", h.Penampung, mauP)
	}

	// V-27 membatalkan unsur PersonList yang hanya berisi penunjuk: entrinya pindah ke
	// baris induk, bukan memegang Kunci yang lalu dipakai ulang baris berikutnya.
	h = ratakan(t, dok{"QuotationData": dok{"BusinessType": "Life"},
		"PersonList": []any{dok{"IdxPerson": "1"}, dok{"Age": "2"}}})
	umum, orang = satu(t, h, "T_GENERAL_POLIS"), satu(t, h, "T_PERSONLIST")
	if len(h.Penampung) != 1 || h.Penampung[0].Kunci != umum.Kunci || h.Penampung[0].Kunci == orang.Kunci {
		t.Errorf("penampung baris batal %+v (umum %d, orang %d)", h.Penampung, umum.Kunci, orang.Kunci)
	}
}

// Murni dan deterministik: masukan sama -> hasil sama persis (termasuk ROW_UID);
// IDPEGA lain -> ROW_UID lain.
func TestFlattenDeterministik(t *testing.T) {
	// JSON mentah, bukan peta Go: json.Marshal mengurutkan kunci menurut abjad, sehingga
	// "urutan dokumen" tidak dapat diuji lewatnya.
	b := []byte(`{"QuotationData":{"BusinessType":"Life"},"PersonList":[{"Age":"1"},{"Age":"2"}],` +
		`"UjiB":"2","UjiA":"1","InwardScale":{"UjiD":"4","UjiC":"3"}}`)
	h1, err1 := Flatten(Masukan{IDPega: idUji, DataJSON: b})
	h2, err2 := Flatten(Masukan{IDPega: idUji, DataJSON: b})
	h3, err3 := Flatten(Masukan{IDPega: "UJI-KELAS NB-2", DataJSON: b})
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal(err1, err2, err3)
	}
	if !reflect.DeepEqual(h1, h2) {
		t.Error("dua jalan atas masukan sama berbeda")
	}
	// Penampung berurutan dokumen (urutan kunci JSON, bukan urutan abjad).
	var urut []string
	for _, m := range h1.Penampung {
		urut = append(urut, m.Medan)
	}
	if strings.Join(urut, ",") != "UjiB,UjiA,UjiD,UjiC" {
		t.Errorf("urutan penampung %v", urut)
	}
	if teks(h1.Baris["T_PERSONLIST"][0], "ROW_UID") == teks(h3.Baris["T_PERSONLIST"][0], "ROW_UID") {
		t.Error("ROW_UID tidak bergantung pada IDPEGA")
	}
}

// TestFlattenFEA - tiket 41 (A142): LocationList/FEAList -> T_FEALIST (induk T_LOCATIONLIST, SEQ_NO urut), halaman
// tertanam DataFEA DILIPAT ke baris FEA; jumlah unit teks apa adanya; nol medan ke penampung.
func TestFlattenFEA(t *testing.T) {
	fea := []any{dok{"APAR": "2", "Sprinkler": "0", "InfoFEA": "UJI INFO", "DataFEA": dok{"PrivateTruckBrigade": "1", "TeamSOPSafety": "UJI"}},
		dok{"Hydrant": "3"}}
	h := ratakan(t, dok{"QuotationData": dok{"BusinessType": "Life"}, "LocationList": []any{dok{"FEAList": fea, "Property": dok{"ObjectNo": "1"}}}})
	baris := h.Baris["T_FEALIST"]
	if len(baris) != 2 {
		t.Fatalf("%d baris T_FEALIST, mau 2", len(baris))
	}
	lokasi := satu(t, h, "T_LOCATIONLIST")
	for i, mau := range []map[string]string{
		{"APAR": "2", "SPRINKLER": "0", "INFO_FEA": "UJI INFO", "PRIVATE_TRUCK_BRIGADE": "1", "TEAM_SOP_SAFETY": "UJI"},
		{"HYDRANT": "3"},
	} {
		for k, v := range mau {
			if got := teks(baris[i], k); got != v {
				t.Errorf("FEA[%d].%s = %q, mau %q", i, k, got, v)
			}
		}
		if baris[i].Induk != lokasi.Kunci {
			t.Errorf("FEA[%d] induk %v, mau baris lokasi %v", i, baris[i].Induk, lokasi.Kunci)
		}
		if got := angka(baris[i], "SEQ_NO"); got != strconv.Itoa(i+1) {
			t.Errorf("FEA[%d] SEQ_NO %s", i, got)
		}
	}
	if n := len(h.Penampung); n != 0 {
		t.Errorf("penampung %d: %+v", n, h.Penampung)
	}
}
