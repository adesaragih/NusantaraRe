package daftar_test

// Uji daftar modul sungguhan - struktur tim satu folder per modul, R4.
//
// ⛔ Nilai harapan di sini disalin dari `modul/daftar.go` SEBELUM refactor
// (HEAD fcc3a8d): perakit harus menghasilkan sambungan yang SAMA dengan
// sambungan tangan yang digantikannya - bukan sambungan yang kebetulan lulus.
// Modul yang dimulai sesudahnya boleh MENAMBAH, tidak boleh mengubah.

import (
	"slices"
	"sort"
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
	//
	// ⚠️ SUBSET, bukan sama persis: modul yang dimulai sesudahnya boleh
	// menambah modul dan sambungan, tetapi tidak boleh MENGUBAH yang lama.
	// Uji yang menuntut daftar persis akan memaksa setiap modul baru menyunting
	// berkas milik tim inti ini (temuan uji coba bab 5).
	mau := []inti.Sambungan{
		{Kontrak: "kontrak.KlaimKomite", Penyedia: "claimlife", Pemakai: "komiteclaimlife"},
		{Kontrak: "kontrak.PembacaPolis", Penyedia: "premiumlistlife", Pemakai: "claimlife"},
	}
	ada := map[inti.Sambungan]bool{}
	for _, s := range r.Sambungan {
		ada[s] = true
	}
	for _, s := range mau {
		if !ada[s] {
			t.Errorf("sambungan %+v hilang; sambungan hasil perakit %+v", s, r.Sambungan)
		}
	}
	// Satu pemakai satu penyedia per kontrak: tidak ada sambungan lama yang
	// berpindah penyedia diam-diam.
	for _, s := range r.Sambungan {
		for _, m := range mau {
			if s.Kontrak == m.Kontrak && s.Pemakai == m.Pemakai && s != m {
				t.Errorf("sambungan %+v menggantikan %+v", s, m)
			}
		}
	}
	var nama []string
	for _, m := range r.Modul {
		nama = append(nama, m.Nama())
	}
	if !sort.StringsAreSorted(nama) {
		t.Errorf("modul %v tidak berurutan menurut nama", nama)
	}
	for _, lama := range []string{"claimlife", "komiteclaimlife", "premiumlistlife", "treatycontractout"} {
		if !slices.Contains(nama, lama) {
			t.Errorf("modul %s hilang dari daftar %v", lama, nama)
		}
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
	dapat := daftar.NamaLama()
	for lama, baru := range mau {
		if dapat[lama] != baru {
			t.Errorf("nama lama %s -> %q, mau %q (peta %v)", lama, dapat[lama], baru, dapat)
		}
	}
}

// Sumber migrasi = `modul.SumberMigrasi` sebelum refactor: Claim Life,
// Komite, PremiumList, dan inti - Treaty Contract Out tanpa migrasi (tco4).
// Kini: satu sumber per modul terdaftar yang bermigrasi, ditambah inti.
func TestSumberMigrasiSamaDenganDaftarLama(t *testing.T) {
	bermigrasi := 0
	for _, p := range daftar.Terdaftar() {
		if p.Migrasi != nil {
			bermigrasi++
		}
	}
	sumber := daftar.SumberMigrasi()
	if len(sumber) != bermigrasi+1 || bermigrasi < 3 {
		t.Fatalf("%d sumber migrasi, mau %d modul bermigrasi + inti", len(sumber), bermigrasi)
	}
	langkah, err := migrasi.Daftar(false, sumber...)
	if err != nil {
		t.Fatal(err)
	}
	var nama []string
	for _, l := range langkah {
		nama = append(nama, l.Nama)
	}
	for _, lama := range []string{"001_t_work_claim.sql", "030_komite_kaskade_dan_lebar_id.sql", "050_t_work_polis.sql", "900_m_nav_menu.sql"} {
		if !slices.Contains(nama, lama) {
			t.Errorf("langkah %s hilang dari daftar pelari", lama)
		}
	}
	if nama[0] != "001_t_work_claim.sql" {
		t.Errorf("langkah pertama %s, mau 001_t_work_claim.sql", nama[0])
	}
}
