package services

// Layar Adjustment — daftar penyesuaian warisan dan satu penyesuaian utuh
// (Old Data · New Data · Attachment · History).
//
// ⛔ SATU jalur baca untuk kedua cabang. Proporsional dan non-proporsional
// berbeda TAB dan MEDAN, dan perbedaan itu hidup di lapis tampilan; yang
// dibaca dari Oracle sama persis untuk keduanya.
//
// ⛔ NOL terjemahan nilai di sini, dan itu disengaja. Tanggal, angka, dan
// persen diformat di layar oleh `inti/frontend/lib/format.ts` — SATU
// pemformat untuk kedua panel. Panel Old dan New menampilkan angka yang akan
// dibandingkan mata ke mata; penerjemah kedua di sini akan membuat satu
// panel membulatkan berbeda dari yang lain, dan bedanya terbaca sebagai beda
// data.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyinadjustment/backend/models"
	"nusantarare/modul/treatyinadjustment/backend/repository"
)

var (
	// ErrPenyesuaianTidakAda - pengenalnya tidak menunjuk penyesuaian mana pun.
	ErrPenyesuaianTidakAda = repository.ErrPenyesuaianTidakAda
	// ErrJSONPenyesuaianRusak - dokumennya ada tetapi tidak dapat diurai.
	ErrJSONPenyesuaianRusak = repository.ErrJSONPenyesuaianRusak
)

// DaftarPenyesuaianWarisan membaca grid daftar layar Adjustment.
//
// Daftar kosong BUKAN galat: ia jawaban yang sah ("belum ada penyesuaian"),
// berbeda dari penyesuaian tertentu yang tidak ada.
func (l *Layanan) DaftarPenyesuaianWarisan(ctx context.Context, p inti.Pelaku) ([]models.BarisPenyesuaian, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	return l.gudang.DaftarPenyesuaianWarisan(ctx)
}

// PenyesuaianWarisan membaca satu penyesuaian: kedua sisinya sekaligus.
//
// ⛔ Kedua sisi dibaca dalam SATU pembacaan, dari dokumen yang sama. Dua
// pembacaan untuk dua panel yang dibandingkan berdampingan membuka jendela
// tempat satu sisi berubah di antara keduanya.
func (l *Layanan) PenyesuaianWarisan(ctx context.Context, p inti.Pelaku, id string) (models.Penyesuaian, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.Penyesuaian{}, err
	}
	// ⛔ Pengenal KOSONG ditolak di sini: "tidak ada" berbeda dari "tidak
	// ditanyakan". Pengenalnya TEKS (`1000080/R02`) — tidak diuraikan
	// menjadi angka.
	if strings.TrimSpace(id) == "" {
		return models.Penyesuaian{}, fmt.Errorf("%w: pengenal penyesuaian kosong", ErrIDTidakSah)
	}
	return l.gudang.BacaPenyesuaianPendaratan(ctx, id)
}

// PenyesuaianTidakAda menjawab apakah galatnya "penyesuaiannya tidak ada".
//
// Dipisah supaya handler tidak perlu mengimpor `repository` - arah
// `handlers -> services -> repository` dijaga penjaga inti.
func PenyesuaianTidakAda(err error) bool { return errors.Is(err, ErrPenyesuaianTidakAda) }
