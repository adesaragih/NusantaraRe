package daftar_test

// Uji daftar modul sungguhan - struktur tim satu folder per modul, R4.
//
// ⛔ Nilai harapan di sini disalin dari `modul/daftar.go` SEBELUM refactor
// (HEAD fcc3a8d): perakit harus menghasilkan sambungan yang SAMA dengan
// sambungan tangan yang digantikannya - bukan sambungan yang kebetulan lulus.

import (
	"reflect"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/daftar"
	"nusantarare/inti/backend/migrasi"
)

func TestSambunganSamaDenganDaftarLama(t *testing.T) {
	r, err := daftar.Rakit(inti.NewDasar(nil), config.Config{}, func(string) {})
	if err != nil {
		t.Fatal(err)
	}
	// modul/daftar.go lama:
	//   svcCL := ...DenganPembacaPolis(premiumlistservices.PembacaPolis(svcPL))
	//   svcKM := ...DenganKlaim(claimlifeservices.KlaimUntukKomite(svcCL))
	mau := []inti.Sambungan{
		{Kontrak: "kontrak.KlaimKomite", Penyedia: "claimlife", Pemakai: "komiteclaimlife"},
		{Kontrak: "kontrak.PembacaPolis", Penyedia: "premiumlistlife", Pemakai: "claimlife"},
	}
	if !reflect.DeepEqual(r.Sambungan, mau) {
		t.Errorf("sambungan %+v\nmau %+v", r.Sambungan, mau)
	}
	var nama []string
	for _, m := range r.Modul {
		nama = append(nama, m.Nama())
	}
	if mau := []string{"claimlife", "komiteclaimlife", "premiumlistlife", "treatycontractout"}; !reflect.DeepEqual(nama, mau) {
		t.Errorf("modul %v, mau %v", nama, mau)
	}
}

// Peta nama lama = `modul.NamaLama` sebelum refactor; kini tiap modul
// menyatakan nama lamanya sendiri.
func TestNamaLamaSamaDenganDaftarLama(t *testing.T) {
	mau := map[string]string{
		"premiumlist": "premiumlistlife",
		"komite":      "komiteclaimlife",
		"treaty":      "treatycontractout",
	}
	if dapat := daftar.NamaLama(); !reflect.DeepEqual(dapat, mau) {
		t.Errorf("nama lama %v, mau %v", dapat, mau)
	}
}

// Sumber migrasi = `modul.SumberMigrasi` sebelum refactor: Claim Life,
// Komite, PremiumList, dan inti - Treaty Contract Out tanpa migrasi (tco4).
func TestSumberMigrasiSamaDenganDaftarLama(t *testing.T) {
	sumber := daftar.SumberMigrasi()
	if len(sumber) != 4 {
		t.Fatalf("%d sumber migrasi, mau 4 (claimlife, komiteclaimlife, premiumlistlife, inti)", len(sumber))
	}
	langkah, err := migrasi.Daftar(false, sumber...)
	if err != nil {
		t.Fatal(err)
	}
	pertama, terakhir := langkah[0].Nama, langkah[len(langkah)-1].Nama
	if pertama != "001_t_work_claim.sql" || terakhir != "900_m_nav_menu.sql" {
		t.Errorf("langkah %s ... %s, mau 001_t_work_claim.sql ... 900_m_nav_menu.sql", pertama, terakhir)
	}
}
