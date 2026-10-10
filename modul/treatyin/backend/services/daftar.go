package services

// Layar daftar kontrak — ronde layar 1.
//
// ⛔ Satu jalur baca, dan ia TIDAK menyaring apa pun. Layar lama punya ikon
// saring per kolom dan penomoran halaman (`Section/InputTreatyInOffer.xml`),
// dan keduanya BELUM dibangun di sini: penyaringan yang dibangun sebelum ada
// satu baris pun untuk disaring tidak dapat diuji terhadap apa pun.
//
// ⚠️ Tabelnya KOSONG hari ini — tiket 59 yang memindahkan kepala kontrak
// warisan belum jalan. Daftar kosong karena itu jawaban yang benar, bukan
// kegagalan, dan layar menyatakannya lewat `Kosong` yang menyebut tiket 59.

import (
	"context"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

// DaftarKontrak membaca seluruh kontrak untuk layar daftar.
func (l *Layanan) DaftarKontrak(ctx context.Context, p inti.Pelaku) ([]models.BarisDaftarKontrak, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	b, err := l.gudang.DaftarKontrak(ctx)
	if err != nil {
		return nil, err
	}
	if b == nil {
		// Daftar kosong, bukan nil: pemanggil JSON tidak perlu membedakan
		// `null` dari `[]` untuk pertanyaan yang jawabannya "belum ada".
		return []models.BarisDaftarKontrak{}, nil
	}
	return b, nil
}
