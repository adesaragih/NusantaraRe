package handlers_test

// Uji Choose Polis atas polis NB Fac In sistem baru (perintah work owner 10-10-2026: "pertama cek json_polis, jika json
// nya kosong ambil dari flat nya"): JSON_POLIS berisi -> dokumen itu; JSON_POLIS penunjuk kasus -> dokumen tabel datar;
// keduanya tidak ada -> 409 OQ-CFI-11. Fixture `UJI-`.

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/services"
	"nusantarare/modul/claimfacin/backend/tiruan"
)

// pilihPolisUji - CopyNB polis uji; mengembalikan kode dan nilai halaman jawaban.
func (u *uji) pilihPolisUji(id string) (int, map[string]any, map[string]any) {
	u.t.Helper()
	kode, out := u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "CopyNB", Konteks: services.ModalPilihPolis,
		Param: polisUji + "|1"})
	hal, _ := out["halaman"].(map[string]any)
	nl, _ := hal["nilai"].(map[string]any)
	return kode, out, nl
}

func TestChoosePolisJSONBerisiTidakMembacaDatar(t *testing.T) {
	u := baruUji(t)
	u.a.DokumenDatar[tiruan.KunciPolis(polisUji, "1")] = []byte(strings.Replace(polisFire, "UJI FIRE", "UJI DATAR", 1))
	kode, out, nl := u.pilihPolisUji(u.buat())
	u.wajib(kode, http.StatusOK, out, "pilih polis")
	if nl[models.OQ+"BusinessName"] != "UJI FIRE" {
		t.Fatalf("BusinessName %v (JSON_POLIS berisi harus dipakai)", nl[models.OQ+"BusinessName"])
	}
}

func TestChoosePolisPenunjukMemakaiTabelDatar(t *testing.T) {
	u := baruUji(t)
	// DEV 10-10-2026: DATA_JSON polis NB sistem baru = penunjuk kasus 27 karakter
	u.a.Dokumen[tiruan.KunciPolis(polisUji, "1")] = []byte(`{"CaseIDNewNB":"UJI-NB-1"}`)
	u.a.DokumenDatar[tiruan.KunciPolis(polisUji, "1")] = []byte(strings.Replace(polisFire, "UJI FIRE", "UJI DATAR", 1))
	id := u.buat()
	kode, out, nl := u.pilihPolisUji(id)
	u.wajib(kode, http.StatusOK, out, "pilih polis")
	if nl[models.OQ+"BusinessName"] != "UJI DATAR" || nl[models.JalurNoPolis] != polisUji ||
		nl[models.OQ+"InsuredName"] != "UJI TERTANGGUNG" {
		t.Fatalf("halaman polis datar: BusinessName %v PolicyNo %v Insured %v", nl[models.OQ+"BusinessName"],
			nl[models.JalurNoPolis], nl[models.OQ+"InsuredName"])
	}
	// muat ulang (muatPolis) memakai jalur yang sama
	kode, out = u.minta(http.MethodGet, "/api/claim-fac-in/kasus/"+id, admin, "", nil)
	u.wajib(kode, http.StatusOK, out, "buka kasus")
	if nl, _ := out["halaman"].(map[string]any)["nilai"].(map[string]any); nl[models.OQ+"BusinessName"] != "UJI DATAR" {
		t.Fatalf("muat ulang: BusinessName %v", nl[models.OQ+"BusinessName"])
	}
}

func TestChoosePolisTanpaJSONDanDatar(t *testing.T) {
	u := baruUji(t)
	u.a.Dokumen[tiruan.KunciPolis(polisUji, "1")] = []byte(`{"CaseIDNewNB":"UJI-NB-1"}`)
	kode, out, _ := u.pilihPolisUji(u.buat())
	u.wajib(kode, http.StatusConflict, out, "polis tanpa isi")
	if g, _ := out["galat"].(string); !strings.Contains(g, "OQ-CFI-11") {
		t.Fatalf("galat %v", out["galat"])
	}
}

func TestMuatKasusTetapTerbukaBilaTabelDatarGagal(t *testing.T) {
	// muatPolis gagal-lunak: kasus berpolis penunjuk tetap terbuka tanpa halaman polis; Choose Polis menolak (500).
	u := baruUji(t)
	u.a.Dokumen[tiruan.KunciPolis(polisUji, "1")] = []byte(`{"CaseIDNewNB":"UJI-NB-1"}`)
	u.a.DokumenDatar[tiruan.KunciPolis(polisUji, "1")] = []byte(polisFire)
	id := u.buat()
	kode, out, _ := u.pilihPolisUji(id)
	u.wajib(kode, http.StatusOK, out, "pilih polis")
	u.a.GagalDatar = errors.New("UJI tabel datar tak terbaca")
	kode, out = u.minta(http.MethodGet, "/api/claim-fac-in/kasus/"+id, admin, "", nil)
	u.wajib(kode, http.StatusOK, out, "buka kasus saat tabel datar gagal")
	kode, out, _ = u.pilihPolisUji(id)
	u.wajib(kode, http.StatusInternalServerError, out, "Choose Polis saat tabel datar gagal")
}
