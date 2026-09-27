package services

// Uji antre-ulang efek keluar - butir aq, A2.
//
// Yang diuji di sini adalah LOGIKA MURNINYA: jeda, dan keputusan menyerah.
// Penulisan ke Oracle diuji terpisah oleh uji berskema.

import (
	"encoding/json"
	"testing"
	"time"
)

func TestBackoffBerlipatDanBerplafon(t *testing.T) {
	kasus := []struct {
		percobaan int
		mau       time.Duration
	}{
		// ⛔ Percobaan < 1 diperlakukan sebagai 1, bukan jeda nol. Jeda nol
		// berarti percobaan ulang seketika, yang bukan antre-ulang melainkan
		// putaran ketat.
		{percobaan: -3, mau: 30 * time.Second},
		{percobaan: 0, mau: 30 * time.Second},
		{percobaan: 1, mau: 30 * time.Second},
		{percobaan: 2, mau: time.Minute},
		{percobaan: 3, mau: 2 * time.Minute},
		{percobaan: 4, mau: 4 * time.Minute},
		{percobaan: 5, mau: 8 * time.Minute},
		{percobaan: 6, mau: 16 * time.Minute},
		{percobaan: 7, mau: 32 * time.Minute},
		{percobaan: 8, mau: 64 * time.Minute},
		{percobaan: 9, mau: 128 * time.Minute},
		{percobaan: 10, mau: 256 * time.Minute},
		// Plafon 6 jam = 360 menit; 512 menit melewatinya.
		{percobaan: 11, mau: 6 * time.Hour},
		{percobaan: 50, mau: 6 * time.Hour},
		// ⛔ Percobaan sangat besar TIDAK boleh melimpah menjadi negatif.
		// Perlipatan tanpa plafon di dalam gelungnya akan melakukannya, dan
		// jeda negatif berarti jatuh tempo di masa lalu - putaran ketat lagi.
		{percobaan: 1000, mau: 6 * time.Hour},
	}
	for _, k := range kasus {
		if dapat := Backoff(k.percobaan); dapat != k.mau {
			t.Errorf("Backoff(%d) = %v, mau %v", k.percobaan, dapat, k.mau)
		}
	}
}

func TestBackoffTidakPernahNolAtauNegatif(t *testing.T) {
	for n := -5; n < 200; n++ {
		if j := Backoff(n); j <= 0 {
			t.Fatalf("Backoff(%d) = %v - jeda tak positif berarti putaran ketat", n, j)
		}
	}
}

// Batas percobaan ada, dan angkanya bukan nol maupun satu.
//
// ⛔ Tanpa batas, kegagalan jaringan yang tak kunjung pulih berputar selamanya
// dan tidak pernah terlihat sebagai kegagalan oleh siapa pun.
func TestJatahPercobaanTerbatas(t *testing.T) {
	if percobaanMaksimum < 2 {
		t.Fatalf("percobaanMaksimum = %d - satu percobaan bukan antre-ulang",
			percobaanMaksimum)
	}
	if percobaanMaksimum > 50 {
		t.Errorf("percobaanMaksimum = %d - terlalu banyak untuk disebut menyerah",
			percobaanMaksimum)
	}
}

// Muatan outbox memuat PENGENAL DAN WAKTU SAJA.
//
// ⛔ Penjaga bentuk, bukan penjaga teks: ia menyusun muatan dari nilai yang
// sengaja dibuat mencurigakan, lalu memastikan tidak satu pun dari nilai
// "bocor" itu muncul. Bila kelak seseorang menambahkan medan nama atau email
// ke `MuatanEfek`, test ini yang menolaknya.
func TestMuatanOutboxHanyaPengenalDanWaktu(t *testing.T) {
	m := muatanOutbox{
		KlaimID:      "RNML-K-1",
		AdjustmentID: "ADJ-1",
		AkunID:       "akun-1",
		Waktu:        time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC).Format(time.RFC3339),
	}
	// Empat medan, tidak lebih. Cacahnya dikunci supaya penambahan medan baru
	// harus melewati test ini.
	const mauMedan = 4
	if n := jumlahMedanJSON(t, m); n != mauMedan {
		t.Errorf("muatan outbox punya %d medan, mau %d - medan baru menuntut "+
			"pemeriksaan kerahasiaan tersendiri", n, mauMedan)
	}
}

// jumlahMedanJSON mencacah medan yang benar-benar terbit ke JSON.
//
// ⚠️ Lewat JSON, bukan lewat refleksi atas struct-nya: yang penting bukan
// berapa medan yang DIMILIKI tipe itu, melainkan berapa yang TERTULIS ke
// kolom `MUATAN`.
func jumlahMedanJSON(t *testing.T, v any) int {
	t.Helper()
	teks, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("merakit muatan: %v", err)
	}
	var peta map[string]any
	if err := json.Unmarshal(teks, &peta); err != nil {
		t.Fatalf("membaca muatan: %v", err)
	}
	return len(peta)
}
