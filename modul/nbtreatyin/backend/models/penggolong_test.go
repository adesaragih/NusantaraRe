package models

// Uji seam 3 - penggolong jenis usaha (tiket 06: AC 19, 20, 21, 22, 67).
//
// ⭐ Nilai harapan TIDAK dihitung dengan cara yang sama seperti kodenya:
// `harapanKorpus` dibangkitkan dari baris XML `DecisionTable/BusinessType_DeT`
// oleh skrip Python terpisah (docs/alat, korpus.json), yang mengevaluasi tabel
// sendiri; kasus `kasusTangan` dipilih tangan dari bab 10 INVENTARIS-XML.md.

import "testing"

// harapanKorpus: setiap pasangan (GroupPanel, BusinessOldId) yang tertulis di
// tabel - 128 kode - beserta hasil evaluasi baris-pertama-cocok dari XML.
var harapanKorpus = [][3]string{
	{"002", "03", "PA"},
	{"003", "36", "Bonding"},
	{"003", "37", "Bonding"},
	{"003", "38", "Bonding"},
	{"003", "39", "Bonding"},
	{"003", "40", "Bonding"},
	{"003", "41", "Bonding"},
	{"003", "42", "Bonding"},
	{"003", "43", "Bonding"},
	{"003", "44", "Bonding"},
	{"003", "45", "Bonding"},
	{"003", "46", "Bonding"},
	{"003", "47", "Bonding"},
	{"003", "48", "Bonding"},
	{"003", "49", "Bonding"},
	{"003", "50", "Bonding"},
	{"003", "51", "Bonding"},
	{"003", "52", "Bonding"},
	{"003", "C2", "Bonding"},
	{"003", "C3", "Bonding"},
	{"003", "C4", "Bonding"},
	{"003", "C5", "Bonding"},
	{"003", "C6", "Bonding"},
	{"003", "C7", "Bonding"},
	{"003", "C8", "Bonding"},
	{"003", "C9", "Bonding"},
	{"003", "D1", "Bonding"},
	{"003", "D2", "Bonding"},
	{"003", "D3", "Bonding"},
	{"003", "D4", "Bonding"},
	{"003", "D5", "Bonding"},
	{"003", "D6", "Bonding"},
	{"003", "D7", "Bonding"},
	{"003", "D8", "Bonding"},
	{"003", "D9", "Bonding"},
	{"003", "F1", "Bonding"},
	{"003", "F2", "Bonding"},
	{"003", "F3", "Bonding"},
	{"003", "F4", "Bonding"},
	{"003", "F5", "Bonding"},
	{"003", "95", "BondingKBG"},
	{"003", "96", "BondingKBG"},
	{"003", "97", "BondingKBG"},
	{"003", "98", "BondingKBG"},
	{"003", "99", "BondingKBG"},
	{"003", "A1", "BondingKBG"},
	{"003", "A2", "BondingKBG"},
	{"003", "A3", "BondingKBG"},
	{"003", "A4", "BondingKBG"},
	{"003", "A5", "BondingKBG"},
	{"003", "A6", "BondingKBG"},
	{"003", "A7", "BondingKBG"},
	{"003", "A8", "BondingKBG"},
	{"003", "A9", "BondingKBG"},
	{"003", "B4", "BondingKBG"},
	{"003", "B5", "BondingKBG"},
	{"003", "B6", "BondingKBG"},
	{"003", "B8", "BondingKBG"},
	{"003", "B9", "BondingKBG"},
	{"003", "C1", "BondingKBG"},
	{"003", "20", "HE"},
	{"003", "05", "MarineHull"},
	{"003", "55", "MarineHull"},
	{"003", "21", "Glass"},
	{"003", "15", "Liability"},
	{"003", "24", "Liability"},
	{"003", "89", "Liability"},
	{"003", "90", "Liability"},
	{"003", "91", "Liability"},
	{"003", "E7", "Liability"},
	{"003", "E9", "Liability"},
	{"003", "16", "AllRisk"},
	{"003", "06", "AviationHull"},
	{"003", "14", "Burglary"},
	{"003", "07", "Car"},
	{"003", "08", "Ear"},
	{"003", "09", "ElectronicEquipment"},
	{"003", "19", "Fidelity"},
	{"003", "27", "GolfInsurance"},
	{"003", "17", "CIT"},
	{"003", "18", "CIS"},
	{"003", "11", "Boiler"},
	{"003", "93", "Workmen"},
	{"003", "10", "MBD"},
	{"003", "12", "ContractorsPM"},
	{"003", "79", "CustomBond"},
	{"003", "80", "CustomBond"},
	{"003", "81", "CustomBond"},
	{"003", "82", "CustomBond"},
	{"003", "83", "CustomBond"},
	{"003", "84", "CustomBond"},
	{"003", "85", "CustomBond"},
	{"003", "86", "CustomBond"},
	{"003", "87", "CustomBond"},
	{"003", "88", "CustomBond"},
	{"003", "28", "FireStyle1"},
	{"003", "B3", "LandRig"},
	{"006", "78", "OilGas"},
	{"006", "22", "FireStyle1"},
	{"006", "26", "FireStyle1"},
	{"006", "29", "FireStyle1"},
	{"006", "57", "FireStyle1"},
	{"006", "01", "FireStyle2"},
	{"006", "34", "FireStyle2"},
	{"006", "59", "FireStyle2"},
	{"007", "02", "MBUCar"},
	{"007", "60", "MBUMotorCycle"},
	{"009", "L1", "Life"},
	{"009", "L2", "Life"},
	{"009", "L3", "Life"},
	{"009", "L4", "Life"},
	{"009", "L5", "Life"},
	{"009", "L6", "Life"},
	{"009", "L7", "Life"},
	{"009", "L8", "Life"},
	{"009", "L9", "Life"},
	{"009", "L10", "Life"},
	{"009", "L11", "Life"},
	{"009", "L12", "Life"},
	{"009", "L13", "Life"},
	{"009", "L14", "Life"},
	{"009", "L15", "Life"},
	{"009", "L16", "Life"},
	{"009", "L17", "Life"},
	{"009", "L18", "Life"},
	{"009", "L19", "Life"},
	{"009", "L20", "Life"},
	{"009", "L21", "Life"},
}

func TestKe128KodeSamaDenganSistemLama(t *testing.T) { // AC 21
	if len(harapanKorpus) != 128 {
		t.Fatalf("harapan memuat %d kode, korpus 128", len(harapanKorpus))
	}
	unik := map[string]bool{}
	for _, h := range harapanKorpus {
		if unik[h[1]] {
			t.Fatalf("kode %q muncul dua kali - korpus menyatakan 128 unik", h[1])
		}
		unik[h[1]] = true
		if got := GolongkanJenisUsaha(h[0], h[1]); got != h[2] {
			t.Errorf("GroupPanel %q BusinessOldId %q: dapat %q, sistem lama %q", h[0], h[1], got, h[2])
		}
	}
}

func TestPenggolongKasusTangan(t *testing.T) {
	kasusTangan := []struct {
		nama, gp, bo, harap string
	}{
		// contoh data produksi di rancangan-tabel-datar §4.1: 006 + 01 -> baris 31
		{"produksi 006/01", "006", "01", "FireStyle2"},
		// sel `"26"   ` berspasi DI LUAR kutip -> tetap cocok baris 30
		{"sel 26 berspasi luar kutip", "006", "26", "FireStyle1"},
		// 006 tanpa kode yang dikenal -> penampung panel baris 32
		{"penampung panel 006", "006", "ZZ", "Fire"},
		// 003 tanpa kode yang dikenal -> penampung baris 26
		{"penampung panel 003", "003", "ZZ", "Aneka"},
		// baris 1 tidak menguji BusinessOldId sama sekali
		{"panel 001 apa saja", "001", "03", "Medicare"},
		// panel 002 hanya punya baris berkode 03
		{"panel 002 kode lain", "002", "04", JenisUsahaTakDikenal},
	}
	for _, k := range kasusTangan {
		if got := GolongkanJenisUsaha(k.gp, k.bo); got != k.harap {
			t.Errorf("%s: dapat %q, harap %q", k.nama, got, k.harap)
		}
	}
}

func TestPenggolongBerhentiDiBarisPertama(t *testing.T) { // AC 19
	// Panel 003 + kode 20 cocok baris 5 (HE) DAN penampung baris 26 (Aneka).
	// Evaluasi seluruh baris akan berakhir di Aneka.
	if got := GolongkanJenisUsaha("003", "20"); got != "HE" {
		t.Fatalf("dapat %q: penggolong tidak berhenti di baris pertama yang cocok", got)
	}
}

func TestPenggolongBawaanUnknown(t *testing.T) { // AC 20
	for _, k := range [][2]string{{"", ""}, {"010", "01"}, {"6", "01"}} {
		if got := GolongkanJenisUsaha(k[0], k[1]); got != "UNKNOWN" {
			t.Errorf("%v: dapat %q, harap UNKNOWN (bukan kosong, bukan galat)", k, got)
		}
	}
}

func TestPenggolongMembandingkanTeksBukanBilangan(t *testing.T) { // ID-16
	// "6" bukan "006": kode bernol-depan yang dibaca sebagai bilangan jatuh ke bawaan.
	if got := GolongkanJenisUsaha("6", "01"); got == "FireStyle2" {
		t.Fatal(`"6" diperlakukan sama dengan "006" - pembandingan harus teks persis`)
	}
}

func TestLiniJiwaEnamBelasKode(t *testing.T) { // AC 67
	if len(KodeLiniJiwa) != 16 {
		t.Fatalf("IsLife memuat %d kode, XML 16", len(KodeLiniJiwa))
	}
	for _, k := range []string{"L1", "L9", "L16"} {
		if !AdalahLiniJiwa(k) {
			t.Errorf("%s harus lini jiwa", k)
		}
	}
	// L17..L21 ada di baris 35 BusinessType_DeT, tetapi TIDAK di IsLife.
	for _, k := range []string{"L17", "L21", "l1", ""} {
		if AdalahLiniJiwa(k) {
			t.Errorf("%q tidak boleh lini jiwa", k)
		}
	}
}
