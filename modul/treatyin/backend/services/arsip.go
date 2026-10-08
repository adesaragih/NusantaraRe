package services

// Arsip muatan keluar — tiket 42.
//
// ⛔ SATU jalur tulis, SATU jalur baca, dan jalur bacanya TIDAK mengembalikan
// muatan. `ADR-0034`: arsipnya disimpan, jalur bacanya tidak ada. `INV-61`
// membawanya sebagai invarian.
//
// ⚠️ Yang dilarang bukan membaca TABELNYA — melainkan membaca KOLOM MUATANNYA.
// Menghitung berapa arsip yang ada tidak membuat arsip menjadi sumber kedua;
// membaca isinya membuatnya begitu.

import (
	"context"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

// Keempat tujuan hilir yang sistem lama kirimi — himpunan TERTUTUP.
//
// ⛔ Disebut satu per satu, bukan diterima apa adanya dari pemanggil: tujuan
// yang salah ketik akan membuat arsipnya tidak pernah ditemukan kembali, dan
// arsip yang tidak dapat ditemukan sama saja dengan arsip yang tidak ada.
var tujuanHilir = map[string]bool{
	"PEGA_TREATY_IN":              true,
	"PEGA_M_TREATY_IN_EDM":        true,
	"PEGA_M_TREATY_IN_DETAIL":     true,
	"PEGA_M_TREATY_IN_DETAIL_EDM": true,
}

// SimpanArsipMuatanKeluar menyimpan satu pengiriman apa adanya.
//
// ⚠️ Muatan TIDAK diperiksa bentuknya, dan itu disengaja. Arsip yang menolak
// muatan yang tidak dapat diurai tidak melestarikan apa pun yang berguna —
// justru muatan yang rusak itulah yang paling perlu ditunjukkan nanti. Yang
// ditolak hanya muatan KOSONG: arsip tanpa isi tidak mengarsipkan apa pun.
func (l *Layanan) SimpanArsipMuatanKeluar(ctx context.Context, p inti.Pelaku, a models.ArsipMuatanKeluar) error {
	if err := inti.WajibIdentitas(p); err != nil {
		return err
	}
	var salah []string
	if a.IDKontrak <= 0 {
		salah = append(salah, fmt.Sprintf("pengenal kontrak %d bukan pengenal yang sah", a.IDKontrak))
	}
	tujuan := strings.TrimSpace(a.Tujuan)
	if !tujuanHilir[tujuan] {
		salah = append(salah, fmt.Sprintf("tujuan %q bukan salah satu dari keempat tujuan hilir", tujuan))
	}
	if strings.TrimSpace(a.Muatan) == "" {
		salah = append(salah, "muatan kosong; arsip tanpa isi tidak mengarsipkan apa pun")
	}
	// ⛔ Waktu pengiriman WAJIB datang dari pemanggil, tidak diisi di sini.
	// Tiket 42 menolak arsip yang disimpan SESUDAH migrasi alih-alih PADA
	// SAAT-nya; mengisi `time.Now()` di sini menghapus bedanya.
	if strings.TrimSpace(a.DikirimPada) == "" {
		salah = append(salah, "waktu pengiriman kosong; ia waktu KIRIM, bukan waktu simpan, jadi tidak dapat diisi di sini")
	} else if _, err := time.Parse(time.RFC3339, a.DikirimPada); err != nil {
		salah = append(salah, fmt.Sprintf("waktu pengiriman %q bukan RFC 3339", a.DikirimPada))
	}
	if len(salah) > 0 {
		return fmt.Errorf("%w: %v", ErrMasukanTidakSah, salah)
	}
	a.Tujuan = tujuan
	return l.gudang.SimpanArsipMuatanKeluar(ctx, a)
}

// BuktiArsipKontrak menunjukkan BAHWA arsipnya ada — bukan isinya.
//
// ⛔ Nilai baliknya `models.BuktiArsip`, yang tidak punya ruas muatan. Itu
// satu-satunya jalur baca tabel arsip di seluruh modul ini, dan ia tidak dapat
// mengembalikan muatan tanpa seseorang menambahkan ruas ke struct-nya —
// perubahan yang terlihat di tinjauan, berbeda dengan `SELECT` yang bertambah
// satu kolom.
//
// ⚠️ Cacah nol adalah JAWABAN, bukan galat: kontrak yang belum pernah dikirim
// ke hilir memang belum punya arsip.
func (l *Layanan) BuktiArsipKontrak(ctx context.Context, p inti.Pelaku, idKontrak int64) (models.BuktiArsip, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.BuktiArsip{}, err
	}
	if idKontrak <= 0 {
		return models.BuktiArsip{}, fmt.Errorf("%w: pengenal kontrak %d bukan pengenal yang sah", ErrMasukanTidakSah, idKontrak)
	}
	b, err := l.gudang.BuktiArsipKontrak(ctx, idKontrak)
	if err != nil {
		return models.BuktiArsip{}, err
	}
	b.IDKontrak = idKontrak
	if b.Tujuan == nil {
		b.Tujuan = []string{}
	}
	return b, nil
}
