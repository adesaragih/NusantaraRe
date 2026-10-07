package handlers_test

// Uji seam 1 - filter A RD `GetListOpportunity` (`A.pxCreateOperator = Param.UserIdentifier`; keputusan work owner
// 06-10-2026 "isi inbox ini muncul hanya untuk akun dia saja", RALAT "hanya filter berdasarkan create operator aja"):
// setiap akun melihat berkas BUATANNYA saja (di posisi mana pun, untuk memantau status); antrean atasan tidak tampil
// di portal. Fixture UJI-.

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func (u *uji) buatOleh(p pelakuUji) string {
	u.t.Helper()
	kode, isi := u.panggil("POST", "/kasus", p, nil)
	if kode != http.StatusCreated {
		u.t.Fatalf("buat: %d %s", kode, isi)
	}
	var k models.Kasus
	if err := json.Unmarshal([]byte(isi), &k); err != nil {
		u.t.Fatal(err)
	}
	return k.ID
}

func TestPortalAdminHanyaBerkasBuatanSendiri(t *testing.T) {
	u := baru(t)
	admin2 := pelakuUji{"UJI-ADMIN2", models.PosisiAdmin}
	milikAdmin := u.buatOleh(admin)
	milikAdmin2 := u.buatOleh(admin2)
	lihat := func(p pelakuUji, harap ...string) {
		t.Helper()
		kode, id, isi := u.daftar(p, "")
		if kode != http.StatusOK || strings.Join(id, " ") != strings.Join(harap, " ") {
			t.Fatalf("%s: %d %v, harap %v (%s)", p.akun, kode, id, harap, isi)
		}
	}
	lihat(admin, milikAdmin)
	lihat(admin2, milikAdmin2)
	// sesudah Submit: pembuat tetap memantau; Sec Head TIDAK melihatnya di portal (bukan pembuat)
	if kode, isi := u.kirim(milikAdmin, admin, halamanLengkap("1")); kode != http.StatusOK {
		t.Fatalf("admin submit: %d %s", kode, isi)
	}
	lihat(admin, milikAdmin)
	lihat(secHead)
	// pemegang Admin + Sec Head yang tidak membuat apa pun: kosong
	lihat(pelakuUji{"UJI-ADMSH", models.PosisiAdmin + "," + models.PosisiSecHead})
}
