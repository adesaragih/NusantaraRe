package handlers

// Pintu HTTP business pada kombinasi kontrak - tiket 07.
//
//	GET    /api/treaty-contract-out/business-master                            pemilih b6294 (BrowseFilterBusiness_RD)
//	GET    /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/business          Business List b10842 (kontrak)
//	POST   /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/business          Add b2785 + Save b6966
//	PUT    /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/business/{bid}    Edit b4547 + Save b6966
//	DELETE /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/business/{bid}    Delete b4826
//
// Nomor baris = `Section/ViewDetailTreatyBusinessGrid.xml` kecuali disebut lain.
//
// Dibaca sesudah: tco_reinsurer.go, services/tco_business.go.

import (
	"encoding/json"
	"net/http"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/treatycontractout/services"
)

func layananBusinessTCO(svc *services.Service) *services.BusinessTCO {
	return svc.BusinessTCO().
		DenganGudang(services.GudangBusinessOracle(svc)).
		DenganKontrak(services.PemegangKontrakOracle(svc)).
		DenganTahun(services.GudangTahunTreatyOracle(svc)).
		DenganMaster(services.MasterBusinessOracle(svc))
}

func daftarkanRuteBusinessTCO(mux *http.ServeMux, svc *services.Service, stub bool) {
	mux.HandleFunc("GET /api/treaty-contract-out/business-master", masterBusinessTCO(svc, stub))
	mux.HandleFunc("GET /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/business", daftarBusinessTCO(svc, stub))
	mux.HandleFunc("POST /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/business", simpanBusinessTCO(svc, stub, false))
	mux.HandleFunc("PUT /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/business/{bid}", simpanBusinessTCO(svc, stub, true))
	mux.HandleFunc("DELETE /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/business/{bid}", hapusBusinessTCO(svc, stub))
}

type jawabanBusinessMaster struct {
	Daftar []services.BusinessMasterTampil `json:"daftar"`
	Total  int                             `json:"total"`
}

func masterBusinessTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		d, err := layananBusinessTCO(svc).Master(r.Context(), inti.PelakuDari(r, stub))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, jawabanBusinessMaster{Daftar: d, Total: len(d)})
	}
}

func daftarBusinessTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		d, err := layananBusinessTCO(svc).Daftar(r.Context(), inti.PelakuDari(r, stub), r.PathValue("id"), r.PathValue("kid"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, d)
	}
}

func simpanBusinessTCO(svc *services.Service, stub, perbarui bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		var masuk services.BusinessMasuk
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&masuk); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body is not valid business JSON")
			return
		}
		id := strings.TrimSpace(masuk.ID)
		switch {
		case !perbarui && id != "":
			galat.Tulis(w, http.StatusBadRequest, "the server assigns the business id; POST must not carry an id")
			return
		case perbarui && id != "" && id != r.PathValue("bid"):
			galat.Tulis(w, http.StatusBadRequest, "id in the body differs from id in the path")
			return
		case perbarui:
			masuk.ID = r.PathValue("bid")
		}
		b, err := layananBusinessTCO(svc).Simpan(r.Context(), inti.PelakuDari(r, stub), r.PathValue("id"), r.PathValue("kid"), masuk)
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, b)
	}
}

type jawabanHapusBusiness struct {
	Pesan string `json:"pesan"`
}

func hapusBusinessTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		pesan, err := layananBusinessTCO(svc).Hapus(r.Context(), inti.PelakuDari(r, stub), r.PathValue("id"),
			r.PathValue("kid"), r.PathValue("bid"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, jawabanHapusBusiness{Pesan: pesan})
	}
}
