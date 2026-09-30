package repository

// Penjaga baris adjustment baru - TANPA Oracle.
//
// Pemilik: tiket 11.

import (
	"context"
	"errors"
	"testing"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/models"
)

// TestBarisBaruTidakBolehLahirBerkeputusan - kepemilikan keputusan.
//
// ⛔ Diperiksa SAAT JALAN, di jalur tulisnya, bukan lewat pola atas teks
// sumber. Pola `KodeStatus = models.KodeAksep` dapat dielakkan oleh
// `models.BarisAdjustment{KodeStatus: "1"}` - titik dua, bukan sama dengan -
// dan oleh setiap variabel perantara. Pemeriksaan nilai tidak dapat dielakkan
// bentuk penulisan apa pun.
//
// `[keputusan work owner 2026-09-15]` Komite menulis, Claim Life membaca.
func TestBarisBaruTidakBolehLahirBerkeputusan(t *testing.T) {
	r := NewPohonKlaim(nil)
	ctx := context.Background()
	for _, kode := range []string{kontrak.KodeAksep, kontrak.KodeDitolak, "9"} {
		_, err := r.SisipkanBaris(ctx, nil, "P-1",
			models.BarisAdjustment{KodeStatus: kode})
		if !errors.Is(err, ErrBarisBaruBerkeputusan) {
			t.Errorf("kode %q: galat = %v, mau ErrBarisBaruBerkeputusan", kode, err)
		}
	}
	// ⛔ Outstanding dan kosong LOLOS - keduanya baris yang baru memulai
	// putaran. Tanpa kasus ini penjaganya dapat menolak SEGALANYA dan tetap
	// tampak benar.
	//
	// ⚠️ Diperiksa lewat `PeriksaBarisBaru`, bukan lewat `SisipkanBaris`:
	// yang lolos penjaga akan lanjut menyentuh basis data, dan test ini
	// berjalan tanpa Oracle. Penolakannya tetap dibuktikan lewat penulis
	// sungguhan di atas, sehingga aturan ini terbukti BENAR-BENAR terpasang.
	for _, kode := range []string{kontrak.KodeOutstanding, ""} {
		if err := PeriksaBarisBaru(kode); err != nil {
			t.Errorf("kode %q ditolak; ia baris pembuka putaran: %v", kode, err)
		}
	}
}
