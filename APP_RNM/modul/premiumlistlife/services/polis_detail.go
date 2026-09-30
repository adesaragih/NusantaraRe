package services

// Layar Premium List Detail - tiket 03 PremiumList Life.
//
// Untuk apa berkas ini: merakit apa yang satu layar detail butuhkan - kepala
// polis, nomor PL, dan satu halaman grid peserta.
//
// ⛔ BENTUK JAWABANNYA MILIK LAPISAN INI, bukan lapisan repository. Handler
// TIDAK BOLEH mengimpor repository (arah ketergantungan
// handlers -> services -> repository), jadi tipe ber-tag JSON tinggal di
// sini - sejajar dengan `polis_inbox.go`.
//
// Dibaca sesudah: polis_nomor.go.

import (
	"context"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/premiumlistlife/models"
	"nusantarare/modul/premiumlistlife/repository"
)

// Galat polis yang perlu dikenali handler.
//
// ⛔ DIRUJUK ULANG DI SINI, bukan diimpor handler dari repository. Arah
// ketergantungannya handlers -> services -> repository, dan handler yang
// mengimpor repository "hanya untuk satu galat" membuka pintu yang berikutnya
// dilewati query.
var (
	// ErrPolisTakDitemukan - tidak ada polis dengan id itu.
	ErrPolisTakDitemukan = repository.ErrPolisTakDitemukan
	// ErrPolisTanpaPeserta - polis belum punya baris peserta.
	ErrPolisTanpaPeserta = repository.ErrPolisTanpaPeserta
	// ErrNomorPLBerbedaAntarPeserta - baris peserta tidak sepakat nomornya.
	ErrNomorPLBerbedaAntarPeserta = repository.ErrNomorPLBerbedaAntarPeserta
	// ErrNomorPLTerbitBersamaan - permintaan lain mendahului.
	ErrNomorPLTerbitBersamaan = repository.ErrNomorPLTerbitBersamaan
	// ErrRekapKosong - tidak ada peserta yang direkap (tiket 05a).
	ErrRekapKosong = repository.ErrRekapKosong
)

// HalamanPesertaPolis adalah satu halaman grid peserta.
type HalamanPesertaPolis struct {
	// Kolom adalah nama kolom, URUT seperti layar lama menampilkannya.
	//
	// ⚠️ IKUT DIKIRIM, tidak diketik ulang di layar. Dua daftar kolom -
	// satu di Go, satu di TypeScript - akan berselisih, dan selisihnya
	// muncul sebagai angka di bawah judul kolom yang salah.
	Kolom []string `json:"kolom"`
	// Baris adalah baris peserta halaman ini.
	Baris []models.BarisPeserta `json:"baris"`
	// Total adalah cacah SELURUH peserta polis, bukan yang di halaman ini.
	Total int `json:"total"`
	// Halaman dan Ukuran adalah halaman yang benar-benar terbaca, sesudah
	// dijepit - bukan yang diminta.
	Halaman int `json:"halaman"`
	Ukuran  int `json:"ukuran"`
}

// KepalaPolis adalah keterangan polis yang tampil di atas grid.
type KepalaPolis struct {
	PolisID    string `json:"polisId"`
	Tipe       string `json:"type"`
	KodeBisnis string `json:"businessCode"`
	Nomor      string `json:"plNumber"`
	// MedanTanpaKolom menyebut medan layar lama yang tidak kami punya.
	//
	// ⛔ IKUT DIKIRIM, bukan disimpan di komentar Go. Orang yang
	// membandingkan layar baru dengan layar lama akan menghitung kolomnya;
	// yang menemukan selisih tanpa penjelasan akan menyimpulkan datanya
	// hilang.
	MedanTanpaKolom []MedanTanpaKolomPolis `json:"medanTanpaKolom"`
}

// MedanTanpaKolomPolis menyebut medan layar lama yang tidak kami punya.
//
// ⛔ DIKIRIM KE LAYAR, bukan disimpan di komentar Go. Orang yang membandingkan
// layar baru dengan layar lama akan menghitung kolomnya; yang menemukan
// selisih tanpa penjelasan akan menyimpulkan datanya hilang.
type MedanTanpaKolomPolis struct {
	Medan  string `json:"medan"`
	Alasan string `json:"alasan"`
}

// DetailPolis melayani layar Premium List Detail.
type DetailPolis struct{ svc *Service }

// DetailPolis menyusun layanannya.
func (s *Service) DetailPolis() *DetailPolis { return &DetailPolis{svc: s} }

// Kepala membaca keterangan polis beserta nomornya.
func (d *DetailPolis) Kepala(ctx context.Context, pelaku inti.Pelaku, polisID string) (
	KepalaPolis, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return KepalaPolis{}, err
	}
	if d == nil || d.svc == nil || !d.svc.PunyaDatabase() {
		return KepalaPolis{}, db.ErrTanpaOracle
	}
	if polisID == "" {
		return KepalaPolis{}, fmt.Errorf("%w: id polis kosong", galat.ErrPermintaanTidakSah)
	}
	// ⛔ TANPA TRANSAKSI, dan itu bukan kelalaian. `Ringkas` membaca kepala
	// polis beserta nomornya dalam SATU query, jadi tidak ada dua pembacaan
	// yang dapat berselisih - dan layanan pembaca yang membuka transaksi
	// membuat penjaga butir bb menghitungnya sebagai pengubah, yaitu
	// mengaburkan mana yang benar-benar menulis.
	ringkas, err := repository.NewNomorPolis(d.svc.DB()).Ringkas(ctx, polisID)
	if err != nil {
		return KepalaPolis{}, err
	}
	return KepalaPolis{
		PolisID:         polisID,
		Tipe:            ringkas.Identitas.Tipe,
		KodeBisnis:      ringkas.Identitas.KodeBisnis,
		Nomor:           ringkas.Nomor.Nomor,
		MedanTanpaKolom: d.MedanTanpaKolom(),
	}, nil
}

// Peserta membaca satu halaman grid peserta.
func (d *DetailPolis) Peserta(ctx context.Context, pelaku inti.Pelaku, polisID string,
	halaman, ukuran int) (HalamanPesertaPolis, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HalamanPesertaPolis{}, err
	}
	if d == nil || d.svc == nil || !d.svc.PunyaDatabase() {
		return HalamanPesertaPolis{}, db.ErrTanpaOracle
	}
	if polisID == "" {
		return HalamanPesertaPolis{}, fmt.Errorf("%w: id polis kosong", galat.ErrPermintaanTidakSah)
	}
	if halaman < 1 {
		halaman = 1
	}
	ukuran = repository.BatasUkuranHalamanPeserta(ukuran)

	hal, err := repository.NewGridPeserta(d.svc.DB()).Ambil(ctx, polisID, halaman, ukuran)
	if err != nil {
		return HalamanPesertaPolis{}, err
	}
	return HalamanPesertaPolis{
		Kolom:   models.NamaKolomGridPeserta(),
		Baris:   hal.Baris,
		Total:   hal.Total,
		Halaman: halaman,
		Ukuran:  ukuran,
	}, nil
}

// MedanTanpaKolom mendaftar medan layar lama yang tidak punya kolom.
//
// ⚠️ Dipakai `Kepala`, dan hanya lewat sana ia sampai ke layar.
func (d *DetailPolis) MedanTanpaKolom() []MedanTanpaKolomPolis {
	// ⚠️ Urutannya dibuat tetap supaya jawabannya tidak berubah-ubah antar
	// permintaan - peta Go tidak berurut, dan jawaban yang berubah urutannya
	// membuat uji dan cache sama-sama gelisah.
	urut := []string{"REINSTYPENAME", "RetrocadedShare"}
	keluar := make([]MedanTanpaKolomPolis, 0, len(urut))
	for _, m := range urut {
		if alasan, ada := models.MedanGridTanpaKolom[m]; ada {
			keluar = append(keluar, MedanTanpaKolomPolis{Medan: m, Alasan: alasan})
		}
	}
	return keluar
}
