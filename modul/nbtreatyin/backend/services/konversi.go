package services

// Untuk apa berkas ini: EFEK KELUAR konversi Arasapas sesudah realisasi
// selesai (Utility2 flow; models/konversi.go).
//
// ⛔ Dijalankan SESUDAH transaksi submit selesai dan kegagalannya TIDAK
// menggagalkan submit (KEPUTUSAN-RONDE-12 butir 7 langkah 3; pola ADR-U-0008
// repo, `inti/backend/outbox`). Di luar produksi efek ini DILEWATI.
//
// ⚠️ `[terbuka]` Sambungan nyatanya belum ada: kunci `M_LINK_SERVICE` rule
// treaty tidak terekspor (salinan kelas `Data-PolicyTreatyIn` stub), dan
// menghubungkan layanan luar menuntut persetujuan manusia (pola
// `outbox.ErrArasapasBelumDisetujui`). Pengirim bawaan karena itu GAGAL TERANG;
// pengirim sungguhan dipasang lewat `DenganKonversi`. Email kegagalan
// (`outbox.EfekEmail`) belum disetujui repo - tidak dikirim.

import (
	"context"
	"errors"
	"log"

	"nusantarare/modul/nbtreatyin/backend/models"
)

// PengirimKonversi mengirim muatan konversi ke Arasapas.
type PengirimKonversi interface {
	Kirim(ctx context.Context, m models.MuatanKonversi) error
}

// ErrKonversiBelumDisambung - pengirim bawaan.
var ErrKonversiBelumDisambung = errors.New(
	"services: sambungan konversi Arasapas treaty belum disetujui (kunci M_LINK_SERVICE tidak terekspor)")

type konversiBelumDisambung struct{}

func (konversiBelumDisambung) Kirim(context.Context, models.MuatanKonversi) error {
	return ErrKonversiBelumDisambung
}

// DenganKonversi memasang pengirim konversi dan menyatakan lingkungan
// produksi (di luar produksi efek dilewati).
func (l *Layanan) DenganKonversi(p PengirimKonversi, produksi bool) *Layanan {
	salin := *l
	salin.konversi, salin.produksi = p, produksi
	return &salin
}

// HasilKirim adalah akibat satu submit.
type HasilKirim struct {
	Kasus models.Kasus `json:"kasus"`
	// PesanKonversi - `FlagErrorKonversi` bila konversi gagal; kosong bila
	// berhasil, dilewati, atau kasus belum selesai.
	PesanKonversi string `json:"pesanKonversi,omitempty"`
}

// konversikan menjalankan efek keluar sesudah realisasi selesai.
func (l *Layanan) konversikan(ctx context.Context, id string, h *models.Halaman) string {
	if !l.produksi {
		return ""
	}
	p := l.konversi
	if p == nil {
		p = konversiBelumDisambung{}
	}
	if err := p.Kirim(ctx, models.RakitMuatanKonversi(id, h)); err != nil {
		log.Printf("nbtreatyin: konversi Arasapas %s gagal: %v", id, err)
		return models.PesanGagalKonversi
	}
	return ""
}
