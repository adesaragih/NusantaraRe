package handlers_test

// Uji seam 1 - kotak masuk Beranda (keputusan work owner 06-10-2026: "beranda menunjukkan berapa banyak case yang
// masuk di akun dia, mengikuti workbasket"): satu baris per workbasket tangga yang DIPEGANG akun, jumlah berkas yang
// MENUNGGU dia - Admin: buatannya yang masih di Admin; Sec Head / Dept Head: antrean workbasket itu. Fixture UJI-.

import (
	"encoding/json"
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func (u *uji) kotakMasuk(p pelakuUji) []models.AntreanKotakMasuk {
	u.t.Helper()
	kode, isi := u.panggil("GET", "/kotak-masuk", p, nil)
	if kode != http.StatusOK {
		u.t.Fatalf("kotak masuk %s: %d %s", p.akun, kode, isi)
	}
	var out []models.AntreanKotakMasuk
	if err := json.Unmarshal([]byte(isi), &out); err != nil {
		u.t.Fatalf("%v: %s", err, isi)
	}
	return out
}

func TestKotakMasukBerandaMengikutiWorkbasket(t *testing.T) {
	u := baru(t)
	u.g.KotakMasuk = map[string]models.PemegangKotakMasuk{models.PosisiSecHead: {NamaWorkbasket: "UJI Treaty Inward - Section Head"}}
	admin2 := pelakuUji{"UJI-ADMIN2", models.PosisiAdmin}
	u.buatOleh(admin)
	naik := u.buatOleh(admin)
	u.buatOleh(admin2)
	if kode, isi := u.kirim(naik, admin, halamanLengkap("1")); kode != http.StatusOK {
		t.Fatalf("admin submit: %d %s", kode, isi)
	}
	harap := func(p pelakuUji, want []models.AntreanKotakMasuk) {
		t.Helper()
		got := u.kotakMasuk(p)
		if len(got) != len(want) {
			t.Fatalf("%s: %+v, harap %+v", p.akun, got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("%s baris %d: %+v, harap %+v", p.akun, i, got[i], want[i])
			}
		}
	}
	// admin: buatannya yang MASIH di Admin (yang sudah naik tidak menunggu dia); nama tanpa master = ID workbasket
	harap(admin, []models.AntreanKotakMasuk{{Workbasket: models.PosisiAdmin, Nama: models.PosisiAdmin, Jumlah: 1}})
	harap(admin2, []models.AntreanKotakMasuk{{Workbasket: models.PosisiAdmin, Nama: models.PosisiAdmin, Jumlah: 1}})
	harap(secHead, []models.AntreanKotakMasuk{{Workbasket: models.PosisiSecHead, Nama: "UJI Treaty Inward - Section Head", Jumlah: 1}})
	// workbasket yang dipegang tetap tampil walau kosong; urut tangga
	harap(pelakuUji{"UJI-SHDH", models.PosisiDeptHead + "," + models.PosisiSecHead}, []models.AntreanKotakMasuk{
		{Workbasket: models.PosisiSecHead, Nama: "UJI Treaty Inward - Section Head", Jumlah: 1},
		{Workbasket: models.PosisiDeptHead, Nama: models.PosisiDeptHead, Jumlah: 0},
	})
	// bukan pemegang workbasket Treaty: daftar kosong, bukan 403 (Beranda tidak boleh gagal karenanya)
	harap(orang, []models.AntreanKotakMasuk{})
}
