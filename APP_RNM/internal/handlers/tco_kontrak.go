package handlers

// Pintu HTTP kontrak treaty di dalam tahun treaty - tiket 04.
//
//	GET  /api/treaty-contract-out/tahun/{id}/kontrak                  grid BrowseTreatyContract_RD
//	POST /api/treaty-contract-out/tahun/{id}/kontrak                  Add b8528 + Save b3618 (baru)
//	PUT  /api/treaty-contract-out/tahun/{id}/kontrak/{kid}            Edit b10519 + Save b3618
//	GET  /api/treaty-contract-out/tahun/{id}/kontrak/akhir-bawaan     SetTanggalTreatyContract
//
// Nomor baris = `Section/InputTreatyContractReinsType.xml`.
//
// Dibaca sesudah: rute_treaty_contract_out.go, services/tco_kontrak.go.

import (
	"encoding/json"
	"net/http"
	"strings"

	"nusantarare/internal/services"
)

// layananKontrakTCO memasang seluruh implementasi nyata.
func layananKontrakTCO(svc *services.Service) *services.KontrakTreatyTCO {
	return svc.KontrakTreatyTCO().
		DenganGudang(services.GudangKontrakOracle(svc)).
		DenganTahun(services.GudangTahunTreatyOracle(svc)).
		DenganJenis(services.PembacaJenisReasuransiOracle(svc))
}

func daftarkanRuteKontrakTCO(mux *http.ServeMux, svc *services.Service, stub bool) {
	mux.HandleFunc("GET /api/treaty-contract-out/tahun/{id}/kontrak", daftarKontrakTCO(svc, stub))
	mux.HandleFunc("POST /api/treaty-contract-out/tahun/{id}/kontrak", simpanKontrakTCO(svc, stub, false))
	mux.HandleFunc("PUT /api/treaty-contract-out/tahun/{id}/kontrak/{kid}", simpanKontrakTCO(svc, stub, true))
	mux.HandleFunc("GET /api/treaty-contract-out/tahun/{id}/kontrak/akhir-bawaan", akhirBawaanKontrakTCO(svc, stub))
}

type jawabanDaftarKontrak struct {
	Daftar []services.KontrakTampil `json:"daftar"`
	Total  int                      `json:"total"`
}

func daftarKontrakTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		d, err := layananKontrakTCO(svc).Daftar(r.Context(), pelakuDari(r, stub), r.PathValue("id"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		tulisJSONPolis(w, jawabanDaftarKontrak{Daftar: d, Total: len(d)})
	}
}

// simpanKontrakTCO - POST menolak `id` dari klien (AC 5); PUT menolak `id`
// badan yang berbeda dari jalur.
func simpanKontrakTCO(svc *services.Service, stub, perbarui bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		var masuk services.KontrakMasuk
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&masuk); err != nil {
			galat(w, http.StatusBadRequest, "badan permintaan bukan JSON kontrak yang sah")
			return
		}
		id := strings.TrimSpace(masuk.ID)
		switch {
		case !perbarui && id != "":
			galat(w, http.StatusBadRequest, "identitas kontrak dibuat server; POST tidak boleh membawa id")
			return
		case perbarui && id != "" && id != r.PathValue("kid"):
			galat(w, http.StatusBadRequest, "id di badan berbeda dari id di jalur")
			return
		case perbarui:
			masuk.ID = r.PathValue("kid")
		}
		k, err := layananKontrakTCO(svc).Simpan(r.Context(), pelakuDari(r, stub), r.PathValue("id"), masuk)
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		tulisJSONPolis(w, k)
	}
}

type jawabanAkhirBawaan struct {
	TreatyEndDate string `json:"treatyEndDate"`
}

func akhirBawaanKontrakTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		akhir, err := layananKontrakTCO(svc).AkhirBawaan(r.Context(), pelakuDari(r, stub), r.PathValue("id"),
			r.URL.Query().Get("mulai"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		tulisJSONPolis(w, jawabanAkhirBawaan{TreatyEndDate: akhir})
	}
}
