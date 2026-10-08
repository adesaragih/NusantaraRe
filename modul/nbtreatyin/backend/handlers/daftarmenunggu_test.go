package handlers_test

// Uji seam 1 - daftar berkas kotak masuk Beranda (keputusan work owner 06-10-2026: klik workbasket / jenis di
// Beranda menampilkan berkasnya tanpa masuk menu NB Treaty In): isi = berkas yang MENUNGGU akun - Admin: buatannya
// yang masih di Admin; Sec Head / Dept Head: antrean workbasket itu. Workbasket yang tidak dipegang = 403.
// Fixture UJI-.

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func (u *uji) menunggu(p pelakuUji, workbasket string) (int, []models.RingkasanKasus) {
	u.t.Helper()
	kueri := ""
	if workbasket != "" {
		kueri = "?workbasket=" + workbasket
	}
	kode, isi := u.panggil("GET", "/kotak-masuk/kasus"+kueri, p, nil)
	if kode != http.StatusOK {
		return kode, nil
	}
	var out []models.RingkasanKasus
	if err := json.Unmarshal([]byte(isi), &out); err != nil {
		u.t.Fatalf("%v: %s", err, isi)
	}
	return kode, out
}

func idDari(b []models.RingkasanKasus) string {
	id := []string{}
	for _, x := range b {
		id = append(id, x.ID)
	}
	sort.Strings(id)
	return strings.Join(id, " ")
}

func TestDaftarMenungguBeranda(t *testing.T) {
	u := baru(t)
	diAdmin := u.buatOleh(admin)
	naik := u.buatOleh(admin)
	u.buatOleh(pelakuUji{"UJI-ADMIN2", models.PosisiAdmin})
	u.g.Halaman[diAdmin].Setel(models.HalamanPolis+".CedingCoName", "UJI CEDING")
	u.g.Halaman[diAdmin].Setel(models.HalamanPolis+".StartDate", "2026-10-01")
	if kode, isi := u.kirim(naik, admin, halamanLengkap("1")); kode != http.StatusOK {
		t.Fatalf("admin submit: %d %s", kode, isi)
	}
	for _, c := range []struct {
		nama, wb string
		p        pelakuUji
		kode     int
		harap    string
	}{
		{"admin, workbasket Admin: hanya buatannya yang masih di Admin", models.PosisiAdmin, admin, http.StatusOK, diAdmin},
		{"admin, semua workbasket", "", admin, http.StatusOK, diAdmin},
		{"Sec Head: antrean Sec Head", models.PosisiSecHead, secHead, http.StatusOK, naik},
		{"admin meminta workbasket yang tidak dipegang", models.PosisiSecHead, admin, http.StatusForbidden, ""},
		{"bukan workbasket tangga", "UJI-WB-LAIN", secHead, http.StatusForbidden, ""},
	} {
		kode, b := u.menunggu(c.p, c.wb)
		if kode != c.kode || (kode == http.StatusOK && idDari(b) != c.harap) {
			t.Errorf("%s: %d %v, harap %d %q", c.nama, kode, idDari(b), c.kode, c.harap)
		}
	}
	// kolom tambahan Beranda: Ceding Company, Inception Date
	_, b := u.menunggu(admin, models.PosisiAdmin)
	if len(b) != 1 || b[0].CedingCoName != "UJI CEDING" || b[0].StartDate != "2026-10-01" {
		t.Errorf("kolom tambahan: %+v", b)
	}
}
