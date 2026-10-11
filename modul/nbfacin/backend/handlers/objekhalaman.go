package handlers

// Tab Object FIRE per halaman (tiket 75 F3).

import (
	"encoding/json"
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// lokasiHalamanKabel - objekKabel + T_LOCATIONLIST.ID (kunci PATCH …/objek).
type lokasiHalamanKabel struct {
	ID string `json:"id"`
	objekKabel
}

// halamanObjekKabel - `{total, halaman, ukuran, lokasi}`.
type halamanObjekKabel struct {
	Total   int                  `json:"total"`
	Halaman int                  `json:"halaman"`
	Ukuran  int                  `json:"ukuran"`
	Lokasi  []lokasiHalamanKabel `json:"lokasi"`
}

// bacaObjekHalaman - GET /api/nbfacin/kasus/{caseId}/objek/halaman?halaman=&ukuran=&cari=.
func bacaObjekHalaman(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		halaman, err1 := strconv.Atoi(r.URL.Query().Get("halaman"))
		ukuran, err2 := strconv.Atoi(r.URL.Query().Get("ukuran"))
		if err1 != nil || err2 != nil {
			galat.Tulis(w, http.StatusBadRequest, "halaman dan ukuran wajib bilangan bulat")
			return
		}
		h, err := svc.BacaObjekHalaman(r.Context(), r.PathValue("caseId"), halaman, ukuran, r.URL.Query().Get("cari"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		out := halamanObjekKabel{Total: h.Total, Halaman: halaman, Ukuran: ukuran, Lokasi: make([]lokasiHalamanKabel, 0, len(h.Lokasi))}
		for i, o := range h.Lokasi {
			out.Lokasi = append(out.Lokasi, lokasiHalamanKabel{ID: h.IDLokasi[i], objekKabel: keKabel(o)})
		}
		galat.TulisJSON(w, out)
	}
}

// batasBadanPerubahanObjek - badan PATCH …/objek (≤ 500 lokasi diubah + 500 ditambah) `[keputusan agent]`.
const batasBadanPerubahanObjek = 8 << 20

// perubahanObjekKabel - badan PATCH `{ubah: [{id, …objek}], tambah: [objek], hapus: [id]}`.
type perubahanObjekKabel struct {
	Ubah   []lokasiHalamanKabel `json:"ubah"`
	Tambah []objekKabel         `json:"tambah"`
	Hapus  []string             `json:"hapus"`
}

// perbaruiObjek - PATCH /api/nbfacin/kasus/{caseId}/objek -> `{jumlahLokasi}`; halaman dibaca ulang layar.
func perbaruiObjek(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b perubahanObjekKabel
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanPerubahanObjek)).Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON perubahan objek yang sah")
			return
		}
		p := services.PerubahanObjek{Ubah: make([]models.LokasiUbah, len(b.Ubah)), Tambah: make([]models.ObjekFire, len(b.Tambah)),
			Hapus: b.Hapus}
		var masalah []string
		for i, o := range b.Ubah {
			var m []string
			p.Ubah[i] = models.LokasiUbah{ID: o.ID, Objek: keModel(i, o.objekKabel, &m)}
			for _, x := range m {
				masalah = append(masalah, "ubah "+x)
			}
		}
		for i, o := range b.Tambah {
			var m []string
			p.Tambah[i] = keModel(i, o, &m)
			for _, x := range m {
				masalah = append(masalah, "tambah "+x)
			}
		}
		if err := services.GalatIsianObjek(masalah); err != nil {
			tulisGalat(w, err)
			return
		}
		n, err := svc.PerbaruiObjek(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), p)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, struct {
			JumlahLokasi int `json:"jumlahLokasi"`
		}{n})
	}
}
