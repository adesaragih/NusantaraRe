package services

// Pelaksana outbox Komite - tiket 07 Komite Claim Life.
//
// Untuk apa berkas ini: menjalankan satu baris outbox `MODUL = KOMITELIFE`
// yang dipungut `PekerjaEfek` (antrean.go, milik Claim Life - DIPAKAI ULANG,
// bukan disalin). Pekerja tahu cara memungut, mencoba ulang, dan menyerah;
// berkas ini tahu arti "email-komite" dan "kasir-komite".
//
// ⛔ AT-LEAST-ONCE + ANTI-DOBEL (ADR-0015, AC 22 spec). Worker menuntaskan
// baris di transaksi yang sama dengan pelaksanaannya, jadi baris `selesai`
// tidak pernah dipungut lagi. Tetapi kiriman yang BERHASIL lalu gagal
// di-commit akan dicoba lagi - karena itu, sebelum Email atau Kasir dikirim,
// dipastikan tidak ada baris LAIN kasus yang sama berjenis sama yang sudah
// `selesai`. Penerima tetap menerima ID baris sebagai kunci idempoten.
//
// ⛔ KODE LINGKUNGAN SATU TEMPAT (ADR-0005): di non-produksi pelaksana ini
// adalah PENGIRIM STUB yang mencatat (km4) - nol panggilan keluar, dan
// barisnya berhenti `gagal-permanen` dengan `ErrPengirimStubNonProduksi`,
// TIDAK PERNAH `selesai`. Di produksi, alamat di-resolve sungguhan lalu berhenti
// terang (`…BelumDisetujui`) sampai manusia menyetujui panggilan nyata.
//
// ⛔ Gerbang `EXIT JIKA RETROID "L0000141"` (langkah 9) DIBUANG `[keputusan
// work owner]` (OQ-064): semua klaim menjalankan seluruh efeknya, termasuk
// Kasir - perubahan perilaku yang menyentuh uang, diterima sadar.
//
// Dibaca sesudah: komite_outbox.go, antrean.go.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/galat"
	"nusantarare/inti/layanan"
	"nusantarare/inti/outbox"
)

// RiwayatEfek menjawab apakah efek yang sama sudah pernah tuntas.
type RiwayatEfek interface {
	SudahSelesai(ctx context.Context, tx *db.Tx, modul, jenis, rujukan, kecualiID string) (bool, error)
}

// PelaksanaKomite menjalankan baris outbox Komite.
type PelaksanaKomite struct {
	Lingkungan inti.Lingkungan
	Resolver   layanan.ResolverEndpoint
	Riwayat    RiwayatEfek
}

// perluCekDobel - Email dan Kasir (AC 22 spec). Arasapas tidak disebut spec.
func perluCekDobel(jenis string) bool {
	return jenis == JenisEfekKomiteEmail || jenis == JenisEfekKomiteKasir
}

// Laksanakan memenuhi `PelaksanaEfek`.
func (p PelaksanaKomite) Laksanakan(ctx context.Context, tx *db.Tx,
	b outbox.BarisEfekKeluar) error {

	if b.Modul != ModulKomiteLife {
		return fmt.Errorf("%w: pelaksana Komite menerima baris modul %q", galat.ErrPermintaanTidakSah, b.Modul)
	}
	var m muatanOutboxKomite
	if err := json.Unmarshal([]byte(b.Muatan), &m); err != nil {
		return fmt.Errorf("%w: muatan outbox komite tak terbaca: %v", galat.ErrPermintaanTidakSah, err)
	}
	switch b.Jenis {
	case JenisEfekKomiteArasapas, JenisEfekKomiteEmail, JenisEfekKomiteKasir:
	default:
		return fmt.Errorf("%w: jenis efek komite %q", galat.ErrPermintaanTidakSah, b.Jenis)
	}
	if perluCekDobel(b.Jenis) {
		if p.Riwayat == nil {
			return errors.New("services: pelaksana Komite tanpa riwayat efek; anti-dobel tidak dapat dijamin")
		}
		sudah, err := p.Riwayat.SudahSelesai(ctx, tx, b.Modul, b.Jenis, b.Rujukan, b.ID)
		if err != nil {
			return err
		}
		if sudah {
			// Sudah terkirim sukses lewat baris lain - JANGAN kirim lagi.
			return nil
		}
	}
	if !p.Lingkungan.AdalahProduksi() {
		// km4: pengirim stub MENCATAT - nol panggilan keluar, dan catatannya
		// jujur: gagal-permanen bersebab, bukan `selesai`.
		return outbox.ErrPengirimStubNonProduksi
	}
	kunci := layanan.KunciLayanan{Kategori1: m.Kategori1, Kategori2: m.Kategori2}
	switch b.Jenis {
	case JenisEfekKomiteArasapas:
		return outbox.EfekArasapas{Resolver: p.Resolver}.Jalankan(ctx, outbox.MuatanEfek{KlaimID: m.KlaimID})
	case JenisEfekKomiteEmail:
		return outbox.EfekEmail{}.Jalankan(ctx, outbox.MuatanEfek{KlaimID: m.KlaimID})
	default: // kasir
		if p.Resolver == nil {
			return layanan.ErrResolverBelumDiputuskan
		}
		if _, err := layanan.AlamatLayanan(ctx, p.Resolver, kunci); err != nil {
			return err
		}
		return outbox.ErrKasirBelumDisetujui
	}
}

// riwayatOracle membaca outbox.
type riwayatOracle struct{ pohon *outbox.Penyimpan }

func (r riwayatOracle) SudahSelesai(ctx context.Context, tx *db.Tx,
	modul, jenis, rujukan, kecualiID string) (bool, error) {
	return r.pohon.EfekSudahSelesai(ctx, tx, modul, jenis, rujukan, kecualiID)
}

// PekerjaKomiteOracle menyusun pekerja outbox Komite.
func PekerjaKomiteOracle(svc *Service) *outbox.PekerjaEfek {
	return outbox.NewPekerjaEfekModul(svc, PelaksanaKomite{
		Lingkungan: svc.Lingkungan(),
		Resolver:   layanan.ResolverLinkServiceOracle(svc),
		Riwayat:    riwayatOracle{pohon: outbox.NewPenyimpan(svc.DB())},
	}, ModulKomiteLife)
}
