package handlers

// Rute kontrak - tiket 14.
//
// ⛔ Dua rute, dan keduanya API. TIDAK ada layar baru: `L-4` menyatakan tidak
// ada spesifikasi layar di mana pun, dan papan melarang mengarangnya. Yang
// mendarat di sini lapisan APLIKASI tiket 14 - jalur simpan tempat INV-53 dan
// INV-29 akhirnya ditegakkan - bukan lapisan layarnya.

import (
	"encoding/json"
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/treatyin/backend/services"
)

// batasBadanJSON - badan permintaan terbesar (1 MiB). Kepala kontrak tidak
// pernah sebesar itu; yang dibatasi badan yang tak berujung.
const batasBadanJSON = 1 << 20

func daftarkanKontrak(pasang func(string, rute)) {
	pasang("POST "+Prefix+"/kontrak", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		r.Body = http.MaxBytesReader(w, r.Body, batasBadanJSON)
		var m services.MasukanKontrak
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body must be valid JSON")
			return
		}
		h, err := l.BuatKontrak(r.Context(), p, m)
		if jawabGalat(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(h)
	})
	// Layar daftar kontrak — ronde layar 1.
	//
	// ⛔ Didaftarkan SEBELUM `/kontrak/{id}`, dan itu bukan selera: ServeMux
	// Go memilih pola paling spesifik, jadi urutannya tidak menentukan —
	// tetapi yang MEMBACA berkas ini membaca dari atas, dan rute daftar yang
	// bersembunyi di bawah rute detail terbaca sebagai tidak ada.
	pasang("GET "+Prefix+"/kontrak", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		baris, err := l.DaftarKontrak(r.Context(), p)
		tulis(w, baris, err)
	})

	// Layar daftar dari tabel WARISAN — keputusan pemilik proses 3 Okt 2026.
	//
	// ⛔ Jalur TERPISAH dari `/kontrak`, dan namanya menyatakannya. Keduanya
	// menjawab pertanyaan berbeda: `/kontrak` mengembalikan yang sudah
	// DIPINDAHKAN (hari ini nol), `/kontrak-warisan` mengembalikan yang ADA
	// DI SISTEM LAMA (1.854).
	pasang("GET "+Prefix+"/kontrak-warisan", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		// Halaman yang bukan angka dibulatkan ke 1, bukan ditolak: nomor
		// halaman datang dari tombol, dan tombol yang menghasilkan 400
		// membuat pemakai menebak apa yang ia lakukan salah.
		halaman, _ := strconv.Atoi(r.URL.Query().Get("halaman"))
		hasil, err := l.DaftarKontrakWarisan(r.Context(), p, halaman)
		tulis(w, hasil, err)
	})

	// Satu kontrak WARISAN — dibuka dari baris daftar.
	//
	// ⛔ Pengenalnya TEKS dan dipakai apa adanya; nol `ParseInt`. Kolom
	// `TREATY_IN.ID` adalah `VARCHAR2`, dan `pengenal()` yang menolak
	// non-angka akan menolak baris yang sah pada hari pertama ia muncul.
	pasang("GET "+Prefix+"/kontrak-warisan/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		k, err := l.BacaKontrakWarisan(r.Context(), p, r.PathValue("id"))
		tulis(w, k, err)
	})

	pasang("GET "+Prefix+"/kontrak/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		id, ok := pengenal(w, r)
		if !ok {
			return
		}
		k, err := l.BacaKontrak(r.Context(), p, id)
		tulis(w, k, err)
	})
}
