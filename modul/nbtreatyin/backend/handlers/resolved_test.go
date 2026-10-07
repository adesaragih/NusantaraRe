package handlers_test

// Uji seam 1 - berkas Resolved TIDAK dapat diubah lagi (perintah work owner 06-10-2026: "yang namanya resolve tidak
// bisa di edit lagi"): layar dibuka hanya-baca tanpa tombol, setiap jalur tulis ditolak 409 dan tidak mengubah apa
// pun - baik oleh pembuatnya maupun oleh atasan penutupnya. Fixture UJI-.

import (
	"encoding/json"
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestResolvedTidakBisaDiubah(t *testing.T) {
	u := baru(t)
	id, _ := u.sampaiDeptHead()
	if kode, isi := u.kirim(id, deptHead, putusan("1")); kode != http.StatusOK || u.g.Kasus[id].StatusWork != models.StatusSelesai {
		t.Fatalf("Dept Head menyetujui: %d %s", kode, isi)
	}
	ditolak := u.buat() // Resolved-Rejected
	if kode, isi := u.kirim(ditolak, admin, halamanLengkap("0")); kode != http.StatusOK || u.g.Kasus[ditolak].StatusWork != models.StatusDitolak {
		t.Fatalf("admin menolak: %d %s", kode, isi)
	}
	riwayat, usulan, produksi := len(u.g.Riwayat), len(u.g.Usulan), len(u.g.Produksi)

	for _, k := range []string{id, ditolak} {
		suggest := u.g.Halaman[k].Ambil("PolicyTreatyIn.Suggest")
		for _, p := range []pelakuUji{admin, secHead, deptHead} {
			kode, isi := u.panggil("GET", "/kasus/"+k, p, nil)
			var ly struct {
				BolehKerja bool   `json:"bolehKerja"`
				Tombol     string `json:"tombol"`
			}
			if kode != http.StatusOK || json.Unmarshal([]byte(isi), &ly) != nil || ly.BolehKerja || ly.Tombol != "" {
				t.Fatalf("%s dibuka %s: %d bolehKerja=%v tombol=%q", k, p.akun, kode, ly.BolehKerja, ly.Tombol)
			}
			// badan lengkap: penolakan karena berkas tertutup, bukan karena permintaan tak sah
			badan := map[string]any{"halaman": putusan("1"), "idDetail": "UJI-D", "idAgen": "UJI-AGEN",
				"urutan": []map[string]string{{"aksi": "SetDueTo"}}}
			for _, c := range []struct{ metode, jalur string }{
				{"PUT", "/kasus/" + k},
				{"POST", "/kasus/" + k + "/hitung"},
				{"POST", "/kasus/" + k + "/bisnis"},
				{"POST", "/kasus/" + k + "/pilih-bisnis"},
				{"POST", "/kasus/" + k + "/pilih-sumber-bisnis"},
				{"POST", "/kasus/" + k + "/nomor-polis"},
				{"POST", "/kasus/" + k + "/kirim"},
			} {
				if kode, isi := u.panggil(c.metode, c.jalur, p, badan); kode != http.StatusConflict {
					t.Errorf("%s %s oleh %s: %d, harap 409 (%s)", c.metode, c.jalur, p.akun, kode, isi)
				}
			}
		}
		if u.g.Halaman[k].Ambil("PolicyTreatyIn.Suggest") != suggest {
			t.Errorf("%s: halaman berubah sesudah ditolak", k)
		}
	}
	if len(u.g.Riwayat) != riwayat || len(u.g.Usulan) != usulan || len(u.g.Produksi) != produksi {
		t.Fatalf("baris tertulis untuk berkas Resolved: riwayat %d->%d usulan %d->%d produksi %d->%d",
			riwayat, len(u.g.Riwayat), usulan, len(u.g.Usulan), produksi, len(u.g.Produksi))
	}
}
