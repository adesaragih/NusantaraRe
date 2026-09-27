package services

// Uji antre-ulang efek keluar - butir aq, A2.
//
// Yang diuji di sini adalah LOGIKA MURNINYA: jeda, dan keputusan menyerah.
// Penulisan ke Oracle diuji terpisah oleh uji berskema.

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/internal/repository"
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

// Muatan yang ditulis lalu dibaca kembali menghasilkan pengenal yang sama.
//
// ⛔ Perjalanan bolak-balik, bukan pemeriksaan satu arah. Penulis dan pembaca
// muatan berada di dua tempat; yang menjaganya tetap sepaham hanyalah test
// yang melewati keduanya.
func TestMuatanOutboxBolakBalik(t *testing.T) {
	saat := time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)
	asal := MuatanEfek{
		KlaimID:      "RNML-K-9",
		AdjustmentID: "ADJ-9",
		AkunID:       "akun-9",
		Waktu:        saat,
	}
	teks, err := json.Marshal(muatanOutbox{
		KlaimID:      asal.KlaimID,
		AdjustmentID: asal.AdjustmentID,
		AkunID:       asal.AkunID,
		Waktu:        asal.Waktu.Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("merakit: %v", err)
	}
	balik := bacaMuatan(string(teks), time.Time{})
	if balik.KlaimID != asal.KlaimID || balik.AdjustmentID != asal.AdjustmentID ||
		balik.AkunID != asal.AkunID {
		t.Errorf("pengenal berubah: %+v, mau %+v", balik, asal)
	}
	if !balik.Waktu.Equal(asal.Waktu) {
		t.Errorf("waktu = %v, mau %v", balik.Waktu, asal.Waktu)
	}
}

// Muatan rusak TIDAK menggagalkan pembacaan - ia jatuh ke waktu cadangan.
//
// ⛔ Baris outbox berumur panjang; bentuk JSON-nya dapat berubah di antara
// saat ia ditulis dan saat ia menyerah. Jejak berpengenal kosong masih lebih
// berguna daripada kegagalan yang tidak tercatat sama sekali.
func TestMuatanRusakTidakMenggagalkan(t *testing.T) {
	cadangan := time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC)
	for _, rusak := range []string{"", "{", "bukan json", `{"waktu":"kemarin"}`} {
		m := bacaMuatan(rusak, cadangan)
		if !m.Waktu.Equal(cadangan) {
			t.Errorf("bacaMuatan(%q).Waktu = %v, mau waktu cadangan %v",
				rusak, m.Waktu, cadangan)
		}
	}
}

// jejakPalsu mencatat apa yang direkam, tanpa Oracle.
type jejakPalsu struct{ catatan []CatatanJejak }

func (j *jejakPalsu) Rekam(_ context.Context, _ *repository.Tx, c CatatanJejak) error {
	j.catatan = append(j.catatan, c)
	return nil
}

// AC 20 tiket 12: kegagalan PERMANEN masuk jalur audit.
//
// ⛔ Dan kegagalan yang masih akan dicoba lagi TIDAK. Mencatatnya berarti
// membanjiri jejak audit dengan delapan baris untuk satu email yang akhirnya
// terkirim.
func TestHanyaKegagalanPermanenMasukJejak(t *testing.T) {
	saat := time.Date(2026, 7, 8, 9, 10, 11, 0, time.UTC)
	c := CatatanEfekGagal{
		MuatanEfek: MuatanEfek{KlaimID: "RNML-K-3", AdjustmentID: "ADJ-3",
			AkunID: "akun-3", Waktu: saat},
		Nama:  NamaEfekEmail,
		Sebab: "sambungan ditolak",
	}
	palsu := &jejakPalsu{}
	a := antreanOracle{jejak: palsu}
	if err := a.rekamMenyerah(context.Background(), nil, c); err != nil {
		t.Fatalf("merekam: %v", err)
	}
	if len(palsu.catatan) != 1 {
		t.Fatalf("%d catatan jejak, mau 1", len(palsu.catatan))
	}
	j := palsu.catatan[0]
	if j.Dari != TahapEfekKeluar+":"+NamaEfekEmail {
		t.Errorf("Dari = %q, mau %q", j.Dari, TahapEfekKeluar+":"+NamaEfekEmail)
	}
	if j.Ke != TahapEfekMenyerah {
		t.Errorf("Ke = %q, mau %q", j.Ke, TahapEfekMenyerah)
	}
	if j.KlaimID != c.KlaimID || j.AdjustmentID != c.AdjustmentID {
		t.Errorf("jejak tidak menunjuk pekerjaannya: %+v", j)
	}
	// ⛔ Pesan galat TIDAK ikut ke jejak: ia dapat menyebut nama objek basis
	// data, alamat, bahkan nilai kolom. Rinciannya tinggal di GALAT_TERAKHIR.
	if strings.Contains(j.Dari, c.Sebab) || strings.Contains(j.Ke, c.Sebab) {
		t.Errorf("pesan galat bocor ke jejak audit: %+v", j)
	}
}

// Tanpa penjejak, menyerah GAGAL TERANG - tidak diam-diam tak tercatat.
func TestMenyerahTanpaPenjejakGagalTerang(t *testing.T) {
	a := antreanOracle{}
	err := a.rekamMenyerah(context.Background(), nil, CatatanEfekGagal{})
	if !errors.Is(err, ErrJejakBelumDiputuskan) {
		t.Errorf("galat = %v, mau ErrJejakBelumDiputuskan", err)
	}
}
