package handlers

// Rute pemulihan limit dan jalur baca warisan — tiket 32, 40, 41.
//
// ⛔ Tiket 41 menyebut `L-4` ("tidak ada spesifikasi layar di mana pun"), jadi
// yang dibangun di sini JALUR BACANYA — bukan layarnya. Rute ini ada supaya
// konsumen hilir yang belum diketahui punya pintu pada hari peralihan
// (`ADR-0051`), dan `CARA MENYALAKANNYA` tiket 41 menuntut catatan aksesnya
// dibaca mingguan: tanpa rute tersendiri, aksesnya tidak dapat dihitung.

import (
	"encoding/json"
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

func daftarkanWarisan(pasang func(string, rute)) {
	// Tiket 32 — SELURUH pemulihan sebuah layer sekaligus.
	pasang("PUT "+Prefix+"/layer/{id}/pemulihan", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		id, ok := pengenalLayer(w, r)
		if !ok {
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, batasBadanJSON)
		var baris []models.PemulihanLimit
		if err := json.NewDecoder(r.Body).Decode(&baris); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body must be a valid JSON array")
			return
		}
		if jawabGalat(w, l.CatatPemulihanLimit(r.Context(), p, id, baris)) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	// Tiket 40 — nilai versi dasar, dibaca lewat rujukannya.
	//
	// ⚠️ Versi pertama menjawab 200 dengan `"ada": false`, BUKAN 404. Tidak
	// adanya versi sebelumnya bukan sumber daya yang hilang - ia jawaban.
	pasang("GET "+Prefix+"/versi/{id}/sebelumnya", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		id, ok := pengenalVersi(w, r)
		if !ok {
			return
		}
		hasil, err := l.NilaiVersiSebelumnya(r.Context(), p, id)
		tulis(w, hasil, err)
	})

	// Tiket 41 — identitas kontrak dalam bentuk sistem lama, diturunkan.
	pasang("GET "+Prefix+"/kontrak/{id}/bentuk-lama", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		id, ok := pengenal(w, r)
		if !ok {
			return
		}
		hasil, err := l.IdentitasBentukLama(r.Context(), p, id)
		tulis(w, hasil, err)
	})
}

// pengenalLayer dan pengenalVersi membaca `{id}` dengan kalimat galat yang
// menyebut APA yang salah bentuknya.
//
// ⛔ Sengaja TIDAK memakai `pengenal` apa adanya: kalimatnya berbunyi
// "contract id must be a number", dan rute ini bukan tentang kontrak.
// Kalimat galat yang menyebut hal yang salah lebih buruk daripada kalimat
// galat yang umum.
func pengenalLayer(w http.ResponseWriter, r *http.Request) (int64, bool) {
	return pengenalBernama(w, r, "layer")
}

func pengenalVersi(w http.ResponseWriter, r *http.Request) (int64, bool) {
	return pengenalBernama(w, r, "version")
}

func pengenalBernama(w http.ResponseWriter, r *http.Request, apa string) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		galat.Tulis(w, http.StatusBadRequest, apa+" id must be a number")
		return 0, false
	}
	return id, true
}
