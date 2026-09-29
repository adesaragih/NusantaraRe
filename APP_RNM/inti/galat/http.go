// Package galat memuat jawaban HTTP bersama - amplop galat `{"galat": ...}`
// dan jawaban JSON - serta galat permintaan yang dipakai setiap modul.
//
// Refactor bentuk B (30-09-2026): dulu `galat()` dan `tulisJSONPolis()` di
// paket handlers, dan `ErrPermintaanTidakSah` di services.
package galat

import (
	"encoding/json"
	"net/http"
)

// batasFormulir membatasi bagian NON-berkas sebuah multipart.
//
// ⚠️ Ini BUKAN batas ukuran berkas - itu `services.BatasUkuranUnggahan`, dan
// ia ditegakkan saat menyalin. Yang di sini membatasi berapa banyak formulir
// yang ditahan di MEMORI sebelum bagian berkasnya dialirkan ke disk.
const BatasFormulir = 1 << 20

func Tulis(w http.ResponseWriter, kode int, pesan string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(kode)
	_ = json.NewEncoder(w).Encode(struct {
		Galat string `json:"galat"`
	}{pesan})
}

// tulisJSONPolis menulis satu jawaban JSON 200.
func TulisJSON(w http.ResponseWriter, isi any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(isi)
}
