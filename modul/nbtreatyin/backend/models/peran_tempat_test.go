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
	if len(got) != len(SemuaTempat) || got[TempatTanggalProduksi] {
		t.Fatalf("tanpa pemetaan setiap tempat TERTUNDA (AC 81): %v", got)
	}
}

func TestTempatTampilMenurutArah(t *testing.T) { // AC 81, 82
	muncul := []PeranTempat{{KodeTempat: TempatTanggalProduksi, Peran: "UJI-PERAN-A", Arah: ArahMuncul}}
	if !TempatTampil(muncul, punya("UJI-PERAN-A"))[TempatTanggalProduksi] || TempatTampil(muncul, punya())[TempatTanggalProduksi] {
		t.Fatal("MUNCUL: hanya pemegang peran yang melihat")
	}
	kecuali := []PeranTempat{{KodeTempat: TempatTanggalProduksi, Peran: "UJI-PERAN-A", Arah: ArahKecuali}}
	if TempatTampil(kecuali, punya("UJI-PERAN-A"))[TempatTanggalProduksi] || !TempatTampil(kecuali, punya())[TempatTanggalProduksi] {
		t.Fatal("KECUALI: semua kecuali pemegang peran")
	}
	bentrok := append(muncul, PeranTempat{KodeTempat: TempatTanggalProduksi, Peran: "UJI-PERAN-B", Arah: ArahKecuali})
	if TempatTampil(bentrok, punya("UJI-PERAN-A"))[TempatTanggalProduksi] {
		t.Fatal("dua arah di satu tempat bertentangan: tetap tertunda (AC 82)")
	}
	asing := []PeranTempat{{KodeTempat: TempatTanggalProduksi, Peran: "UJI-PERAN-A", Arah: "UJI-ARAH"}}
	if TempatTampil(asing, punya("UJI-PERAN-A"))[TempatTanggalProduksi] {
		t.Fatal("arah tak dikenal tidak ditebak (AC 82)")
	}
	lain := []PeranTempat{{KodeTempat: "UJI-TEMPAT-LAIN", Peran: "UJI-PERAN-A", Arah: ArahMuncul}}
	if _, ada := TempatTampil(lain, punya("UJI-PERAN-A"))["UJI-TEMPAT-LAIN"]; ada {
		t.Fatal("tempat yang tidak dibaca layanan tidak ikut dilaporkan")
	}
}
