// Package jejak merekam SIAPA dan KAPAN setiap transisi (ADR-U-0007) ke
// `T_CLAIMLF_JEJAK` - dipakai Claim Life, PremiumList, dan Komite.
//
// Refactor bentuk B (30-09-2026): dulu bagian `services/statusbaris.go` dan
// `repository/jejak.go`.
package jejak

import (
	"context"
	"errors"
	"fmt"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
)

var (
	// ErrJejakBelumDiputuskan - tempat jejak audit belum ada.
	ErrJejakBelumDiputuskan = errors.New("services: tempat jejak audit belum diputuskan")
)

// CatatanJejak adalah satu baris jejak audit sebuah transisi.
type CatatanJejak struct {
	// AdjustmentID kosong pada jalur balik TAHAP: yang berpindah kasusnya,
	// bukan satu baris.
	AdjustmentID string
	// KlaimID selalu terisi.
	//
	// ⛔ RALAT A2, 27-09-2026. Sebelum ini `tahap.go` mengisi `AdjustmentID`
	// dengan pengenal KLAIM - dua hal berbeda dikonflasi, dan jejak jalur
	// balik akan tampak menunjuk baris adjustment yang tidak pernah ada.
	// Tabel `T_CLAIMLF_JEJAK` punya kedua kolom; kini modelnya pun.
	KlaimID string
	Dari    string
	Ke      string
	AkunID  string
	Waktu   time.Time
	// Komentar - `T_CLAIMLF_JEJAK.KOMENTAR` (migrasi 021, OQ-M5): alasan
	// penolakan Admin. Kosong pada transisi lain.
	Komentar string
}

// Jejak merekam SIAPA dan KAPAN untuk setiap transisi status.
//
// ⛔ Kenapa ini antarmuka dan bukan penulisan langsung: ADR-U-0007 menuntut
// setiap transisi terekam, sedangkan TABELNYA belum ada - keputusan membuatnya
// masih `[USULAN]` butir am, dan brief melarang menebaknya. Memisahkannya
// membuat yang belum ada terlihat sebagai satu galat terang, bukan sebagai
// transisi yang diam-diam tak tercatat.
//
// Pola yang sama dengan `Penomor` pada tiket 02.
type Jejak interface {
	Rekam(ctx context.Context, tx *db.Tx, c CatatanJejak) error
}

// JejakBelumDiputuskan adalah implementasi bawaan; ia selalu gagal.
type JejakBelumDiputuskan struct{}

// Rekam selalu gagal, dengan pesan yang menyebut apa yang ditunggu.
func (JejakBelumDiputuskan) Rekam(context.Context, *db.Tx, CatatanJejak) error {
	return fmt.Errorf("%w: tabel jejak audit (butir am) belum disahkan work owner, "+
		"sedangkan ADR-U-0007 menuntut setiap transisi merekam siapa dan kapan",
		ErrJejakBelumDiputuskan)
}

// perekamOracle menulis jejak ke `T_CLAIMLF_JEJAK` - butir am, A2.
//
// ⛔ Ia menggantikan `JejakBelumDiputuskan`, yang selama ini membuat KELIMA
// jalur tulis modul ini menjawab HTTP 501: menolak baris (05), memindah tahap
// (08), menyerahkan ke Komite (10), membuka putaran (11), dan mengaksep
// (audit A0).
type perekamOracle struct{ baca *Penyimpan }

// PerekamJejakOracle menyusun perekam yang menulis ke tabel jejak.
func PerekamJejakOracle(svc inti.Akar) Jejak {
	return perekamOracle{baca: NewPenyimpan(svc.DB())}
}

// Rekam menulis satu catatan, DI DALAM transaksi pemanggilnya.
func (p perekamOracle) Rekam(ctx context.Context, tx *db.Tx,
	c CatatanJejak) error {

	return p.baca.SisipJejak(ctx, tx, c.AdjustmentID, c.KlaimID,
		c.Dari, c.Ke, c.AkunID, c.Waktu, c.Komentar)
}
