package handlers_test

// Uji seam 1 - T_GENERAL_POLIS BERSAMA FacIn + Treaty In (keputusan work
// owner 04-10-2026). Tabel dasar `T_GENERAL_POLIS` (nbfacin 182) dan akar
// `T_WORK_POLIS` memuat baris lini lain: kasus FacIn `T_WORK_POLIS.LINI =
// 'FAC'`. Modul ini hanya membaca dan menulis kasus lininya sendiri
// (`models.LiniKasus` = 'NONLIFE'): baris FAC tidak tampil di daftar portal,
// tidak dapat dibuka, disimpan, dikirim, maupun dibaca riwayatnya - persis
// seperti kasus yang tidak ada (404). Padanan Oracle:
// `repository/lini_db_test.go` (tag db).

import (
	"net/http"
	"reflect"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

// barisFAC menaruh satu kasus FacIn di gudang tiruan: baris T_WORK_POLIS
// ber-LINI 'FAC' + generasi T_GENERAL_POLIS-nya, di antrean admin dengan
// tahap yang sama persis dengan kasus Treaty baru - satu-satunya beda LINI.
func (u *uji) barisFAC(id string) {
	u.t.Helper()
	u.g.Kasus[id] = models.Kasus{ID: id, Position: models.PositionAdmin, StatusWork: models.AssignmentAdmin,
		PositionNote: models.PosisiAdmin, CreateOp: "UJI-FAC", TglCreate: "2026-10-03 09:00:00"}
	u.g.Halaman[id] = models.HalamanBaru()
	u.g.Lini[id] = "FAC"
}

func TestBarisLiniFACTidakTampilDanTidakTerjangkau(t *testing.T) {
	u := baru(t)
	nb := u.buat()
	const fac = "NB-990001" // berawalan sama dengan kasus Treaty: penyaring LINI, bukan awalan ID
	u.barisFAC(fac)
	sebelum := u.g.Kasus[fac]

	// Daftar portal: hanya kasus NONLIFE.
	for _, kueri := range []string{"", "?cari=NB-99", "?posisi=" + models.PosisiAdmin} {
		kode, id, isi := u.daftar(admin, kueri)
		if kode != http.StatusOK {
			t.Fatalf("daftar %q: %d %s", kueri, kode, isi)
		}
		for _, x := range id {
			if x == fac {
				t.Errorf("daftar %q menampilkan baris LINI FAC %s: %v", kueri, fac, id)
			}
		}
		if kueri == "" && !reflect.DeepEqual(id, []string{nb}) {
			t.Errorf("daftar tanpa saringan = %v, harap hanya kasus Treaty %s", id, nb)
		}
	}

	// Setiap pembaca dan penulis kasus: 404 seperti kasus yang tidak ada.
	for _, c := range []struct {
		metode, jalur string
		badan         any
	}{
		{"GET", "/kasus/" + fac, nil},
		{"PUT", "/kasus/" + fac, map[string]any{"halaman": halamanLengkap("1")}},
		{"POST", "/kasus/" + fac + "/kirim", map[string]any{"halaman": halamanLengkap("1")}},
		{"POST", "/kasus/" + fac + "/hitung", map[string]any{"urutan": []map[string]string{{"aksi": "SetDueTo"}}, "halaman": halamanLengkap("1")}},
		{"POST", "/kasus/" + fac + "/nomor-polis", map[string]any{"halaman": halamanLengkap("1")}},
		{"GET", "/kasus/" + fac + "/riwayat", nil},
	} {
		if kode, isi := u.panggil(c.metode, c.jalur, admin, c.badan); kode != http.StatusNotFound {
			t.Errorf("%s %s atas baris LINI FAC: %d %s, harap 404", c.metode, c.jalur, kode, isi)
		}
	}
	if u.g.Kasus[fac] != sebelum {
		t.Errorf("baris FAC berubah: %+v -> %+v", sebelum, u.g.Kasus[fac])
	}
	for _, r := range u.g.Riwayat {
		if r.IDPega == models.KunciInstans(fac) {
			t.Errorf("riwayat ditulis untuk baris FAC: %+v", r)
		}
	}

	// Kasus Treaty tetap terbuka seperti biasa.
	if kode, isi := u.panggil("GET", "/kasus/"+nb, admin, nil); kode != http.StatusOK {
		t.Fatalf("kasus Treaty %s: %d %s", nb, kode, isi)
	}
}
