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

	"nusantarare/internal/repository"
)

// ErrPengirimStubNonProduksi - non-produksi TIDAK mengirim, dan TIDAK mengaku terkirim.
//
// ⛔ Temuan /code-review giliran 10: stub yang menuntaskan baris `selesai`
// tidak terbedakan dari kiriman nyata - anti-dobel lalu menganggap Kasir
// SUDAH dibayar bila basis data non-produksi kelak dipromosikan atau flag
// lingkungannya berubah. Galat permanen ini membuat barisnya berhenti di
// `gagal-permanen` dengan sebab yang menyebut dirinya, bukan `selesai`.
var ErrPengirimStubNonProduksi = errors.New(
	"services: lingkungan bukan produksi - efek Komite TIDAK dikirim (pengirim stub); " +
		"baris ini bukan kiriman yang berhasil")

// ErrKasirBelumDisetujui - pemanggilan Kasir nyata menuntut persetujuan.
var ErrKasirBelumDisetujui = errors.New(
	"services: pemanggilan Kasir belum disetujui; menghubungkan layanan pembayaran " +
		"menuntut persetujuan manusia (OQ-002: kontrak Kasir tidak ada di korpus)")

// RiwayatEfek menjawab apakah efek yang sama sudah pernah tuntas.
type RiwayatEfek interface {
	SudahSelesai(ctx context.Context, tx *repository.Tx, modul, jenis, rujukan, kecualiID string) (bool, error)
}

// PelaksanaKomite menjalankan baris outbox Komite.
type PelaksanaKomite struct {
	Lingkungan Lingkungan
	Resolver   ResolverEndpoint
	Riwayat    RiwayatEfek
}

// perluCekDobel - Email dan Kasir (AC 22 spec). Arasapas tidak disebut spec.
func perluCekDobel(jenis string) bool {
	return jenis == JenisEfekKomiteEmail || jenis == JenisEfekKomiteKasir
}

// Laksanakan memenuhi `PelaksanaEfek`.
func (p PelaksanaKomite) Laksanakan(ctx context.Context, tx *repository.Tx,
	b repository.BarisEfekKeluar) error {

	if b.Modul != ModulKomiteLife {
		return fmt.Errorf("%w: pelaksana Komite menerima baris modul %q", ErrPermintaanTidakSah, b.Modul)
	}
	var m muatanOutboxKomite
	if err := json.Unmarshal([]byte(b.Muatan), &m); err != nil {
		return fmt.Errorf("%w: muatan outbox komite tak terbaca: %v", ErrPermintaanTidakSah, err)
	}
	switch b.Jenis {
	case JenisEfekKomiteArasapas, JenisEfekKomiteEmail, JenisEfekKomiteKasir:
	default:
		return fmt.Errorf("%w: jenis efek komite %q", ErrPermintaanTidakSah, b.Jenis)
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
		return ErrPengirimStubNonProduksi
	}
	kunci := KunciLayanan{Kategori1: m.Kategori1, Kategori2: m.Kategori2}
	switch b.Jenis {
	case JenisEfekKomiteArasapas:
		return EfekArasapas{Resolver: p.Resolver}.Jalankan(ctx, MuatanEfek{KlaimID: m.KlaimID})
	case JenisEfekKomiteEmail:
		return EfekEmail{}.Jalankan(ctx, MuatanEfek{KlaimID: m.KlaimID})
	default: // kasir
		if p.Resolver == nil {
			return ErrResolverBelumDiputuskan
		}
		if _, err := AlamatLayanan(ctx, p.Resolver, kunci); err != nil {
			return err
		}
		return ErrKasirBelumDisetujui
	}
}

// riwayatOracle membaca outbox.
type riwayatOracle struct{ pohon *repository.PohonKlaim }

func (r riwayatOracle) SudahSelesai(ctx context.Context, tx *repository.Tx,
	modul, jenis, rujukan, kecualiID string) (bool, error) {
	return r.pohon.EfekSudahSelesai(ctx, tx, modul, jenis, rujukan, kecualiID)
}

// PekerjaKomiteOracle menyusun pekerja outbox Komite.
func PekerjaKomiteOracle(svc *Service) *PekerjaEfek {
	return NewPekerjaEfekModul(svc, PelaksanaKomite{
		Lingkungan: svc.lingkungan,
		Resolver:   ResolverLinkServiceOracle(svc),
		Riwayat:    riwayatOracle{pohon: repository.NewPohonKlaim(svc.db)},
	}, ModulKomiteLife)
}
