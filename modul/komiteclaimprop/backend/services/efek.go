package services

// Untuk apa berkas ini: PELAKSANA OUTBOX Komite Claim Prop - menjalankan satu baris T_LOG_SERVICE_RNM
// `MODUL = komiteclaimprop` (pola pelaksana Claim Prop / Komite Claim Life, DISALIN bukan diimpor).
//
// ⛔ Non-produksi: baris berhenti `gagal-permanen` dengan `ErrPengirimStubNonProduksi` (nol panggilan keluar).
// Produksi: alamat M_LINK_SERVICE di-resolve sungguhan (kunci VERBATIM dari activity) lalu berhenti terang
// (`…BelumDisetujui`) sampai manusia menyetujui panggilan nyata. Email
// (`SendEmailWithAttachments`) memakai SMTP, bukan M_LINK_SERVICE.

import (
	"context"
	"encoding/json"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/komiteclaimprop/backend/repository"
)

// PelaksanaKomiteClaimProp menjalankan baris outbox modul ini.
type PelaksanaKomiteClaimProp struct {
	Lingkungan inti.Lingkungan
	Resolver   layanan.ResolverEndpoint
}

// Laksanakan memenuhi `outbox.PelaksanaEfek`.
func (p PelaksanaKomiteClaimProp) Laksanakan(ctx context.Context, _ *db.Tx, b outbox.BarisEfekKeluar) error {
	if b.Modul != repository.ModulOutbox {
		return fmt.Errorf("%w: pelaksana Komite Claim Prop menerima baris modul %q", galat.ErrPermintaanTidakSah, b.Modul)
	}
	var m MuatanOutbox
	if err := json.Unmarshal([]byte(b.Muatan), &m); err != nil {
		return fmt.Errorf("%w: muatan outbox Komite Claim Prop tak terbaca: %v", galat.ErrPermintaanTidakSah, err)
	}
	switch b.Jenis {
	case JenisEfekKonversi, JenisEfekKasir, JenisEfekEmailKomite:
	default:
		return fmt.Errorf("%w: jenis efek Komite Claim Prop %q", galat.ErrPermintaanTidakSah, b.Jenis)
	}
	if !p.Lingkungan.AdalahProduksi() {
		return outbox.ErrPengirimStubNonProduksi
	}
	if b.Jenis == JenisEfekEmailKomite {
		return outbox.EfekEmail{}.Jalankan(ctx, outbox.MuatanEfek{KlaimID: m.KlaimID})
	}
	if p.Resolver == nil {
		return layanan.ErrResolverBelumDiputuskan
	}
	if _, err := layanan.AlamatLayanan(ctx, p.Resolver, layanan.KunciLayanan{Kategori1: m.Kategori1,
		Kategori2: m.Kategori2}); err != nil {
		return err
	}
	if b.Jenis == JenisEfekKasir {
		return outbox.ErrKasirBelumDisetujui
	}
	return outbox.ErrArasapasBelumDisetujui
}

// PekerjaKomiteClaimPropOracle menyusun pekerja outbox modul ini (tidak dijalankan modul - lihat `modul.go`).
func PekerjaKomiteClaimPropOracle(svc inti.Akar, l inti.Lingkungan) *outbox.PekerjaEfek {
	return outbox.NewPekerjaEfekModul(svc, PelaksanaKomiteClaimProp{Lingkungan: l,
		Resolver: layanan.ResolverLinkServiceOracle(svc)}, repository.ModulOutbox)
}
