package handlers

// Unggah CSV objek tanpa batas (tiket 75 F1 + F2): POST /api/nbfacin/kasus/{caseId}/objek/csv/berkas (multipart/form-data medan
// `berkas`, DIBACA MENGALIR lewat r.MultipartReader - tanpa ParseMultipartForm / batas badan), GET …/objek/csv/aktif, GET
// …/objek/csv/{idUnggah}, GET …/objek/csv/{idUnggah}/baris?halaman=&ukuran=, POST …/objek/csv/{idUnggah}/simpan dan …/ulangi.
// Langkah yang memulai pekerjaan latar menjawab 202 dengan status. Endpoint lama POST …/objek/csv (JSON <= 3000 baris) tetap.

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// statusUnggahKabel - kontrak status pekerjaan (sesi 0f, 10-10-2026). galat / pesan = paling banyak 100 pertama, berawalan nomor
// baris; galatProses = sebab status 'gagal' ("gagal di baris N: …"); dibuat DD-MM-YYYY HH:mm (jam basis data).
type statusUnggahKabel struct {
	IDUnggah       string   `json:"idUnggah"`
	Jenis          string   `json:"jenis"` // importer pekerjaan: "fire" / "cargo" (tiket 74 P6)
	Status         string   `json:"status"`
	Fase           string   `json:"fase"`
	TotalBaris     int      `json:"totalBaris"`
	BarisDiproses  int      `json:"barisDiproses"`
	BarisTersimpan int      `json:"barisTersimpan"`
	Galat          []string `json:"galat"`
	JumlahGalat    int      `json:"jumlahGalat"`
	Pesan          []string `json:"pesan"`
	JumlahPesan    int      `json:"jumlahPesan"`
	GalatProses    string   `json:"galatProses"`
	NamaBerkas     string   `json:"namaBerkas"`
	Dibuat         string   `json:"dibuat"`
	AdaObjek       bool     `json:"adaObjek"`
}

// teksPesanUnggah - "baris N: teks" (teks yang sudah berawalan "baris N" dibiarkan).
func teksPesanUnggah(p []models.PesanUnggahCsv) []string {
	hasil := make([]string, 0, len(p))
	for _, x := range p {
		t := x.Teks
		if x.NoBaris > 0 && !strings.HasPrefix(t, "baris "+strconv.Itoa(x.NoBaris)) {
			t = fmt.Sprintf("baris %d: %s", x.NoBaris, t)
		}
		hasil = append(hasil, t)
	}
	return hasil
}

func keStatusUnggahKabel(s services.StatusUnggahCsv) statusUnggahKabel {
	u := s.Unggah
	k := statusUnggahKabel{IDUnggah: u.ID, Jenis: u.Jenis, Status: u.Status, Fase: u.Fase, TotalBaris: u.TotalBaris, BarisDiproses: u.BarisDiproses,
		BarisTersimpan: u.BarisTersimpan, Galat: teksPesanUnggah(s.Galat), JumlahGalat: u.JumlahGalat, Pesan: teksPesanUnggah(s.Pesan),
		JumlahPesan: u.JumlahPesan, GalatProses: u.GalatProses, NamaBerkas: u.NamaBerkas, AdaObjek: s.AdaObjek}
	if !u.Dibuat.IsZero() {
		k.Dibuat = u.Dibuat.Format("02-01-2006 15:04")
	}
	return k
}

// unggahBerkasCsv - POST …/objek/csv/berkas: bagian multipart `berkas` dibaca mengalir; bagian lain dilewati.
func unggahBerkasCsv(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if err := inti.WajibIdentitas(pelaku); err != nil {
			tulisGalat(w, err)
			return
		}
		mr, err := r.MultipartReader()
		if err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan harus multipart/form-data dengan medan berkas")
			return
		}
		for {
			bagian, err := mr.NextPart()
			if errors.Is(err, io.EOF) {
				galat.Tulis(w, http.StatusBadRequest, "medan berkas tidak ada")
				return
			}
			if err != nil {
				galat.Tulis(w, http.StatusBadRequest, "badan multipart tidak terbaca")
				return
			}
			if bagian.FormName() != "berkas" {
				_ = bagian.Close()
				continue
			}
			st, err := svc.MulaiUnggahCsv(r.Context(), pelaku, r.PathValue("caseId"), bagian.FileName(), bagian)
			_ = bagian.Close()
			if err != nil {
				tulisGalat(w, err)
				return
			}
			tulisJSONKode(w, http.StatusAccepted, keStatusUnggahKabel(st))
			return
		}
	}
}

// tulisJSONKode - JSON dengan kode status selain 200.
func tulisJSONKode(w http.ResponseWriter, kode int, isi any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(kode)
	_ = json.NewEncoder(w).Encode(isi)
}

// statusUnggahCsv - GET …/objek/csv/{idUnggah}.
func statusUnggahCsv(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		st, err := svc.StatusUnggahCsvKasus(r.Context(), r.PathValue("caseId"), r.PathValue("idUnggah"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keStatusUnggahKabel(st))
	}
}

// unggahCsvAktif - GET …/objek/csv/aktif: status pekerjaan terakhir case, `null` bila belum pernah.
func unggahCsvAktif(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		st, ada, err := svc.UnggahCsvAktif(r.Context(), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		if !ada {
			galat.TulisJSON(w, nil)
			return
		}
		galat.TulisJSON(w, keStatusUnggahKabel(st))
	}
}

// barisUnggahCsv - GET …/objek/csv/{idUnggah}/baris?halaman=&ukuran= -> {total, kolom: string[], baris: Record<kolom,string>[]};
// kolom = judul tampil urut berkas (tiket 74 P6: pekerjaan cargo dibaca per nama judul; FIRE = urutan template).
func barisUnggahCsv(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		halaman, err1 := strconv.Atoi(r.URL.Query().Get("halaman"))
		ukuran, err2 := strconv.Atoi(r.URL.Query().Get("ukuran"))
		if err1 != nil || err2 != nil {
			galat.Tulis(w, http.StatusBadRequest, "halaman dan ukuran wajib bilangan bulat")
			return
		}
		h, err := svc.BarisUnggahCsv(r.Context(), r.PathValue("caseId"), r.PathValue("idUnggah"), halaman, ukuran)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		out := struct {
			Total int                 `json:"total"`
			Kolom []string            `json:"kolom"`
			Baris []map[string]string `json:"baris"`
		}{Total: h.Total, Kolom: h.Kolom, Baris: make([]map[string]string, 0, len(h.Baris))}
		kanonik := models.KolomCsvJenis(h.Jenis)
		for _, b := range h.Baris {
			m := make(map[string]string, len(kanonik))
			for i, k := range kanonik {
				if i < len(b.Isi) {
					m[k] = b.Isi[i]
				}
			}
			out.Baris = append(out.Baris, m)
		}
		galat.TulisJSON(w, out)
	}
}

// langkahUnggahCsv - POST …/objek/csv/{idUnggah}/simpan | /ulangi (202).
func langkahUnggahCsv(svc *services.Service, stubPelaku, ulangi bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		langkah := svc.SimpanUnggahCsv
		if ulangi {
			langkah = svc.UlangiUnggahCsv
		}
		st, err := langkah(r.Context(), pelaku, r.PathValue("caseId"), r.PathValue("idUnggah"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		tulisJSONKode(w, http.StatusAccepted, keStatusUnggahKabel(st))
	}
}
