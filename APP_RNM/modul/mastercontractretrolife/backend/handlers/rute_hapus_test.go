package handlers_test

// Rute hapus berjenjang (paket 7): popup lebih dulu, DELETE wajib membawa
// cacahan popup, konfirmasi basi = 409.

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/mastercontractretrolife/backend/handlers"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

func serverPohon(t *testing.T) *uji {
	u := server(t, true)
	u.g.Reinsurer["UJI-R1"] = models.Reinsurer{ID: "UJI-R1", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1"}
	u.g.Security["UJI-S1"] = models.SecurityReinsurer{ID: "UJI-S1", TreatyReinsurerID: "UJI-R1"}
	u.g.Business["UJI-B1"] = models.Business{ID: "UJI-B1", TreatyContractID: "UJI-K1"}
	return u
}

func TestHapusKontrakPopupLaluKaskade(t *testing.T) {
	u := serverPohon(t)
	kode, badan := u.minta(t, "GET", handlers.Prefix+"/kontrak/UJI-K1/dampak-hapus", true)
	if kode != http.StatusOK || !strings.Contains(badan, `"dampak":{"security":1,"reinsurer":1,"business":1}`) {
		t.Fatalf("popup: %d %s", kode, badan)
	}
	kode, _ = u.kirim(t, "DELETE", handlers.Prefix+"/kontrak/UJI-K1", `{"dampak":{"security":0,"reinsurer":1,"business":1}}`)
	if kode != http.StatusConflict || len(u.g.Security) != 1 {
		t.Errorf("konfirmasi basi: %d, security %d", kode, len(u.g.Security))
	}
	kode, badan = u.kirim(t, "DELETE", handlers.Prefix+"/kontrak/UJI-K1", `{"dampak":{"security":1,"reinsurer":1,"business":1}}`)
	if kode != http.StatusOK || !strings.Contains(badan, `"pesan":"Data Berhasil di Hapus"`) ||
		len(u.g.Kontrak) != 0 || len(u.g.Reinsurer) != 0 || len(u.g.Security) != 0 || len(u.g.Business) != 0 {
		t.Errorf("kaskade: %d %s", kode, badan)
	}
}

func TestDeleteTanpaKonfirmasiDitolak(t *testing.T) {
	u := serverPohon(t)
	for _, badan := range []string{"", `{}`} {
		kode, _ := u.kirim(t, "DELETE", handlers.Prefix+"/business/UJI-B1", badan)
		if kode != http.StatusBadRequest || len(u.g.Business) != 1 {
			t.Errorf("DELETE %q tanpa dampak: %d", badan, kode)
		}
	}
	kode, badan := u.kirim(t, "DELETE", handlers.Prefix+"/business/UJI-B1", `{"dampak":{}}`)
	if kode != http.StatusOK || !strings.Contains(badan, "Data Dengan ID UJI-B1 Berhasil di Hapus") {
		t.Errorf("hapus business: %d %s", kode, badan)
	}
}
