package models

// Uji murni pemetaan TEMPAT -> PERAN -> ARAH (tiket 05; AC 81, 82, 91).
// Peran fiktif berawalan UJI- (AC 91: nol peran karangan).

import "testing"

func punya(peran ...string) func(string) bool {
	return func(p string) bool {
		for _, x := range peran {
			if x == p {
				return true
			}
		}
		return false
	}
}

func TestPemetaanPeranTempatKosongSampaiIAMMenjawab(t *testing.T) { // K12, K16
	if len(PemetaanPeranTempat) != 0 {
		t.Fatalf("pemetaan diisi IAM bersama work owner, bukan ditebak: %+v", PemetaanPeranTempat)
	}
	got := TempatTampil(PemetaanPeranTempat, punya("UJI-PERAN-A", PosisiAdmin))
	if len(got) != len(SemuaTempat) {
		t.Fatalf("setiap tempat dilaporkan: %v", got)
	}
	for kode, tampil := range got {
		if tampil {
			t.Fatalf("tanpa pemetaan setiap tempat TERTUNDA (AC 81): %s", kode)
		}
	}
}

// Kedua belas tempat tiket 05 (grilling ronde 2 P28: `pyUserIdentifier` di
// DetailDeptHeadTreatyIn_UW 3x, GeneralDeptHeadTreatyIn_UW 3x, ListSuggest 4x;
// `pxInsName` di DetailPoliciesNonProportional 2x) - terdaftar sekali, kode unik.
func TestDuaBelasTempatTerdaftar(t *testing.T) { // K16, tiket 05
	if len(SemuaTempat) != 12 {
		t.Fatalf("tiket 05 = 12 tempat, terdaftar %d", len(SemuaTempat))
	}
	per := map[string]int{}
	kode := map[string]bool{}
	for _, tp := range DaftarTempat {
		if kode[tp.Kode] || tp.Syarat == "" || tp.Gerbang == "" {
			t.Fatalf("tempat %+v: kode ganda atau tanpa syarat/gerbang", tp)
		}
		kode[tp.Kode] = true
		per[tp.Section]++
	}
	for section, n := range map[string]int{"DetailDeptHeadTreatyIn_UW": 3, "GeneralDeptHeadTreatyIn_UW": 3,
		"ListSuggest": 4, "DetailPoliciesNonProportional": 2} {
		if per[section] != n {
			t.Errorf("%s: %d tempat, harap %d", section, per[section], n)
		}
	}
}

// `ListSuggest.ProductionDate`: tampil `.IsApproved == 1 && (<id-3> || <id-4>)`,
// wajib `pyRequiredWhen` bunyi sama - masing-masing tempat sendiri. Sel yang
// tidak tampil tidak mewajibkan.
func TestTanggalProduksiTampilDanWajibMenurutTempat(t *testing.T) { // AC 81, tiket 05
	h := HalamanBaru()
	h.Setel(HalamanPolis+".IsApproved", "1")
	wajib := func(tempat map[string]bool) bool {
		for _, m := range MedanWajibBerlaku(h, PosisiAdmin, tempat) {
			if m.Jalur == HalamanPolis+".ProductionDate" {
				return true
			}
		}
		return false
	}
	if TanggalProduksiTampil(h, nil) || wajib(nil) {
		t.Fatal("pemetaan kosong: tidak tampil, tidak wajib")
	}
	tampil := map[string]bool{TempatProduksiTampilOperator3: true}
	if !TanggalProduksiTampil(h, tampil) || wajib(tampil) {
		t.Fatal("tempat tampil saja: tampil, tidak wajib")
	}
	keduanya := map[string]bool{TempatProduksiTampilOperator3: true, TempatProduksiWajibOperator4: true}
	if !wajib(keduanya) {
		t.Fatal("tampil + wajib: Production Date wajib")
	}
	if wajib(map[string]bool{TempatProduksiWajibOperator3: true}) {
		t.Fatal("sel tak tampil tidak menegakkan wajibnya")
	}
	h.Setel(HalamanPolis+".IsApproved", "0")
	if TanggalProduksiTampil(h, keduanya) || wajib(keduanya) {
		t.Fatal("IsApproved bukan 1: tidak tampil")
	}
	if MedanWajibBerlaku(h, "UJI-POSISI-LAIN", keduanya) != nil {
		t.Fatal("posisi di luar tangga: tanpa medan wajib")
	}
}

func TestTempatTampilMenurutArah(t *testing.T) { // AC 81, 82
	muncul := []PeranTempat{{KodeTempat: TempatLabelNonEDM, Peran: "UJI-PERAN-A", Arah: ArahMuncul}}
	if !TempatTampil(muncul, punya("UJI-PERAN-A"))[TempatLabelNonEDM] || TempatTampil(muncul, punya())[TempatLabelNonEDM] {
		t.Fatal("MUNCUL: hanya pemegang peran yang melihat")
	}
	kecuali := []PeranTempat{{KodeTempat: TempatLabelNonEDM, Peran: "UJI-PERAN-A", Arah: ArahKecuali}}
	if TempatTampil(kecuali, punya("UJI-PERAN-A"))[TempatLabelNonEDM] || !TempatTampil(kecuali, punya())[TempatLabelNonEDM] {
		t.Fatal("KECUALI: semua kecuali pemegang peran")
	}
	bentrok := append(muncul, PeranTempat{KodeTempat: TempatLabelNonEDM, Peran: "UJI-PERAN-B", Arah: ArahKecuali})
	if TempatTampil(bentrok, punya("UJI-PERAN-A"))[TempatLabelNonEDM] {
		t.Fatal("dua arah di satu tempat bertentangan: tetap tertunda (AC 82)")
	}
	asing := []PeranTempat{{KodeTempat: TempatLabelNonEDM, Peran: "UJI-PERAN-A", Arah: "UJI-ARAH"}}
	if TempatTampil(asing, punya("UJI-PERAN-A"))[TempatLabelNonEDM] {
		t.Fatal("arah tak dikenal tidak ditebak (AC 82)")
	}
	lain := []PeranTempat{{KodeTempat: "UJI-TEMPAT-LAIN", Peran: "UJI-PERAN-A", Arah: ArahMuncul}}
	if _, ada := TempatTampil(lain, punya("UJI-PERAN-A"))["UJI-TEMPAT-LAIN"]; ada {
		t.Fatal("tempat yang tidak dibaca layanan tidak ikut dilaporkan")
	}
}
