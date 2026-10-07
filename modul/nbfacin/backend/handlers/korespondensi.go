package handlers

// Tab Correspondence (tiket 56): GET / POST /api/nbfacin/kasus/{caseId}/korespondensi, PUT / DELETE …/korespondensi/{no},
// POST …/korespondensi/{no}/kirim, GET /api/nbfacin/korespondensi/penerima/saran. Kontrak frontend api.ts `Korespondensi`.

import (
	"encoding/json"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// batasBadanKorespondensi - badan mentah (pesan bersih ≤ 4000 bita diperiksa layanan; HTML mentah bisa lebih panjang).
const batasBadanKorespondensi = 1 << 20

type penerimaKorespondensiKabel struct {
	Recipient     string `json:"recipient"`
	RecipientMail string `json:"recipientMail"`
}

type korespondensiKabel struct {
	No         string                       `json:"no"`
	Tanggal    string                       `json:"tanggal"`
	Sender     string                       `json:"sender"`
	SenderMail string                       `json:"senderMail"`
	Subject    string                       `json:"subject"`
	Message    string                       `json:"message"`
	Penerima   []penerimaKorespondensiKabel `json:"penerima"`
	Terkirim   bool                         `json:"terkirim"`
}

func keKorespondensiKabel(k models.Korespondensi) korespondensiKabel {
	out := korespondensiKabel{No: k.No, Tanggal: k.Tanggal, Sender: k.Sender, SenderMail: k.SenderMail, Subject: k.Subject,
		Message: k.Message, Penerima: []penerimaKorespondensiKabel{}, Terkirim: k.Terkirim}
	for _, p := range k.Penerima {
		out.Penerima = append(out.Penerima, penerimaKorespondensiKabel(p))
	}
	return out
}

// uraiIsiKorespondensi - badan ketat {subject, message, penerima: [{recipient, recipientMail}]}.
func uraiIsiKorespondensi(w http.ResponseWriter, r *http.Request) (models.IsiKorespondensi, bool) {
	var b struct {
		Subject  string                       `json:"subject"`
		Message  string                       `json:"message"`
		Penerima []penerimaKorespondensiKabel `json:"penerima"`
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanKorespondensi))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&b); err != nil {
		galat.Tulis(w, http.StatusBadRequest, `badan permintaan harus {"subject", "message", "penerima": [{"recipient", "recipientMail"}]}`)
		return models.IsiKorespondensi{}, false
	}
	isi := models.IsiKorespondensi{Subject: b.Subject, Message: b.Message, Penerima: []models.PenerimaKorespondensi{}}
	for _, p := range b.Penerima {
		isi.Penerima = append(isi.Penerima, models.PenerimaKorespondensi(p))
	}
	return isi, true
}

// daftarKorespondensi - GET /api/nbfacin/kasus/{caseId}/korespondensi -> {daftar, pesan}.
func daftarKorespondensi(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := svc.DaftarKorespondensi(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		daftar := make([]korespondensiKabel, 0, len(t.Daftar))
		for _, k := range t.Daftar {
			daftar = append(daftar, keKorespondensiKabel(k))
		}
		galat.TulisJSON(w, struct {
			Daftar []korespondensiKabel `json:"daftar"`
			Pesan  []string             `json:"pesan"`
		}{daftar, t.Pesan})
	}
}

// tambahKorespondensi - POST /api/nbfacin/kasus/{caseId}/korespondensi -> draf baru.
func tambahKorespondensi(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		k, err := svc.TambahKorespondensi(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keKorespondensiKabel(k))
	}
}

// ubahKorespondensi - PUT /api/nbfacin/kasus/{caseId}/korespondensi/{no} (Save) atau POST …/{no}/kirim (Send Email).
func ubahKorespondensi(svc *services.Service, stubPelaku, kirim bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if err := svc.WajibAksesKorespondensi(r.Context(), pelaku); err != nil {
			tulisGalat(w, err) // 401 / 403 sebelum badan diurai
			return
		}
		isi, ok := uraiIsiKorespondensi(w, r)
		if !ok {
			return
		}
		aksi := svc.SimpanKorespondensi
		if kirim {
			aksi = svc.KirimKorespondensi
		}
		k, err := aksi(r.Context(), pelaku, r.PathValue("caseId"), r.PathValue("no"), isi)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keKorespondensiKabel(k))
	}
}

// hapusKorespondensi - DELETE /api/nbfacin/kasus/{caseId}/korespondensi/{no} -> 204.
func hapusKorespondensi(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := svc.HapusKorespondensi(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), r.PathValue("no")); err != nil {
			tulisGalat(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// saranPenerimaKorespondensi - GET /api/nbfacin/korespondensi/penerima/saran?q= -> {baris: [{id, label, ekstra}]}.
func saranPenerimaKorespondensi(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.SaranPenerimaKorespondensi(r.Context(), inti.PelakuDari(r, stubPelaku), r.URL.Query().Get("q"))
		tulisSaran(w, d, err)
	}
}
