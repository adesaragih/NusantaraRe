package handlers_test

// Medan pengendali tampil (work owner 11-10-2026 "CETANGAN Deductible Info TIDAK BERFUNGSI!"): XML `change:postValue` /
// `change:refresh` tanpa activity - nilai baru dikirim ke server dan section dievaluasi ulang, sehingga medan yang
// tampilnya bergantung padanya (`VIS=OTHER` / `VIS=CONDITION`) muncul seketika. Dahulu medan itu tanpa aksi: nilai hanya
// berubah di layar, server tidak pernah mengevaluasi ulang.

import (
	"net/http"
	"testing"

	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/services"
)

// adaMedan - medan berjalur `jalur` tampil di salah satu tata `out` (utama, panel, modal).
func adaMedan(out map[string]any, jalur string) bool {
	var cari func(any) bool
	cari = func(x any) bool {
		switch v := x.(type) {
		case map[string]any:
			if v["jenis"] == "medan" && v["jalur"] == jalur {
				return true
			}
			for _, w := range v {
				if cari(w) {
					return true
				}
			}
		case []any:
			for _, w := range v {
				if cari(w) {
					return true
				}
			}
		}
		return false
	}
	for _, k := range []string{"tata", "panel", "modal"} {
		if cari(out[k]) {
			return true
		}
	}
	return false
}

func TestDeductibleInfoMemunculkanFormat(t *testing.T) {
	u := baruUji(t)
	id := u.sampaiEstimasi()
	est := models.KunciPanel(models.PanelObjekEst, models.DaftarObjek, 1)
	item := models.KunciPanel(models.PanelItemEst, models.DaftarItem(1), 1)
	kode, out := u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "SetObjectItem", Konteks: est, Indeks: 1, Param: "1"})
	u.wajib(kode, http.StatusOK, out, "pilih property item")
	j := func(p string) string { return models.JalurItem(1, 1) + "." + p }
	if adaMedan(out, j("FormType")) {
		t.Fatal("Format tampil sebelum Deductible Info dicentang")
	}
	// centang Deductible Info -> Format tampil (Estimasi LS4 `VIS=OTHER:.DeductibleType=='true'`)
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "PostValue", Konteks: item,
		Masukan: map[string]string{j("DeductibleType"): "true"}})
	u.wajib(kode, http.StatusOK, out, "centang Deductible Info")
	if !adaMedan(out, j("FormType")) {
		t.Fatal("Format tidak tampil sesudah Deductible Info dicentang")
	}
	// Format 2 -> baris % / of / Minimum (LS11 `VIS .FormType==2`)
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "PostValue", Konteks: item,
		Masukan: map[string]string{j("FormType"): "2"}})
	u.wajib(kode, http.StatusOK, out, "pilih Format")
	if !adaMedan(out, j("Amount")) || !adaMedan(out, j("DeductibleValue")) {
		t.Fatal("isian Format 2 (% / Minimum) tidak tampil")
	}
	if got := models.AmbilJalur(u.g.Halaman(id), j("DeductibleType")); got != "true" {
		t.Fatalf("DeductibleType tersimpan %q", got)
	}
	// hilangkan centang -> Format tersembunyi lagi
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "PostValue", Konteks: item,
		Masukan: map[string]string{j("DeductibleType"): "false"}})
	u.wajib(kode, http.StatusOK, out, "hapus centang Deductible Info")
	if adaMedan(out, j("FormType")) {
		t.Fatal("Format tetap tampil sesudah centang dihapus")
	}
}
