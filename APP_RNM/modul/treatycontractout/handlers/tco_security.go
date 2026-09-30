package handlers

// Pintu HTTP security reinsurer - tiket 06.
//
//	GET    /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/reinsurer/{rid}/security         Security Reinsurer (ViewDetailTreatyReinsurerGrid1 b5277)
//	POST   /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/reinsurer/{rid}/security         Add b15459 + Save b20246
//	PUT    /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/reinsurer/{rid}/security/{sid}   Edit b17252 + Save b20246
//	DELETE /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/reinsurer/{rid}/security/{sid}   Delete b17559
//
// Nomor baris = `Section/InputTreatyContractReinsType.xml` kecuali disebut lain.
// Pemilih `Security Name` (`BrowseAgentReinsSOA_RD` b19711) = rute
// `GET /reinsurer-master` tiket 05 - RD yang sama.
//
// Dibaca sesudah: tco_reinsurer.go, services/tco_security.go.

import (
	"encoding/json"
	"net/http"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/treatycontractout/services"
)

func layananSecurityTCO(svc *services.Service) *services.SecurityTCO {
	return svc.SecurityTCO().
		DenganGudang(services.GudangSecurityOracle(svc)).
		DenganReinsurer(services.GudangReinsurerOracle(svc)).
		DenganKontrak(services.PemegangKontrakOracle(svc)).
		DenganTahun(services.GudangTahunTreatyOracle(svc)).
		DenganMaster(services.MasterReinsurerOracle(svc))
}

const jalurSecurityTCO = "/api/treaty-contract-out/tahun/{id}/kontrak/{kid}/reinsurer/{rid}/security"

func daftarkanRuteSecurityTCO(mux *http.ServeMux, svc *services.Service, stub bool) {
	mux.HandleFunc("GET "+jalurSecurityTCO, daftarSecurityTCO(svc, stub))
	mux.HandleFunc("POST "+jalurSecurityTCO, simpanSecurityTCO(svc, stub, false))
	mux.HandleFunc("PUT "+jalurSecurityTCO+"/{sid}", simpanSecurityTCO(svc, stub, true))
	mux.HandleFunc("DELETE "+jalurSecurityTCO+"/{sid}", hapusSecurityTCO(svc, stub))
}

func daftarSecurityTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		d, err := layananSecurityTCO(svc).Daftar(r.Context(), inti.PelakuDari(r, stub), r.PathValue("id"),
			r.PathValue("kid"), r.PathValue("rid"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, d)
	}
}

func simpanSecurityTCO(svc *services.Service, stub, perbarui bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		var masuk services.SecurityMasuk
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&masuk); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body is not valid security JSON")
			return
		}
		id := strings.TrimSpace(masuk.ID)
		switch {
		case !perbarui && id != "":
			galat.Tulis(w, http.StatusBadRequest, "the server assigns the security id; POST must not carry an id")
			return
		case perbarui && id != "" && id != r.PathValue("sid"):
			galat.Tulis(w, http.StatusBadRequest, "id in the body differs from id in the path")
			return
		case perbarui:
			masuk.ID = r.PathValue("sid")
		}
		s, err := layananSecurityTCO(svc).Simpan(r.Context(), inti.PelakuDari(r, stub), r.PathValue("id"),
			r.PathValue("kid"), r.PathValue("rid"), masuk)
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, s)
	}
}

type jawabanHapusSecurity struct {
	Pesan string `json:"pesan"`
}

func hapusSecurityTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		pesan, err := layananSecurityTCO(svc).Hapus(r.Context(), inti.PelakuDari(r, stub), r.PathValue("id"),
			r.PathValue("kid"), r.PathValue("rid"), r.PathValue("sid"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, jawabanHapusSecurity{Pesan: pesan})
	}
}
