package handlers

// Untuk apa berkas ini: rute pemilih Source Of Business (XOL Retro) - popup
// Harness `SOB` / `Section/SourceHierarki` (lihat `services/sumberbisnis.go`).

import (
	"net/http"

	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbtreatyin/backend/models"
)

// badanSumberBisnis - klik satu baris TreeGrid beserta halaman layar
// (showHarness `pySubmitData=Yes`: isian layar ikut terkirim).
type badanSumberBisnis struct {
	Halaman *models.Halaman `json:"halaman"`
	// IDAgen - `.ID` baris RD `BrowseAgentHierarkiList_RD` yang diklik.
	IDAgen string `json:"idAgen"`
}

func (h *rute) sumberBisnis(w http.ResponseWriter, r *http.Request) {
	out, err := h.l.DaftarSumberBisnis(r.Context(), h.pelaku(r))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	if out == nil {
		out = []models.BarisAgen{}
	}
	galat.TulisJSON(w, out)
}

func (h *rute) pilihSumberBisnis(w http.ResponseWriter, r *http.Request) {
	var b badanSumberBisnis
	if !bacaJSON(w, r, &b) {
		return
	}
	ly, err := h.l.PilihSumberBisnis(r.Context(), h.pelaku(r), r.PathValue("id"), b.IDAgen, b.Halaman)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, ly)
}
