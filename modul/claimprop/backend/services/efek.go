package services

// Untuk apa berkas ini: PELAKSANA OUTBOX Claim Prop - menjalankan satu baris T_LOG_SERVICE_RNM `MODUL = claimprop`
// (pola Komite Claim Life `PelaksanaKomite`, DISALIN bukan diimpor). Pekerja tahu cara memungut dan mencoba ulang;
// berkas ini tahu arti "konversi-klaim", "kasir", dan "email-komite".
//
// ⛔ Non-produksi: pengirim STUB - nol panggilan keluar, baris berhenti `gagal-permanen` dengan
// `ErrPengirimStubNonProduksi`, tidak pernah `selesai`. Produksi: alamat M_LINK_SERVICE di-resolve sungguhan (kunci
// Kategori_1/Kategori_2 VERBATIM dari activity) lalu berhenti terang (`…BelumDisetujui`) sampai manusia menyetujui
// panggilan nyata (prompt §6 butir 6). Email komite (SendEmailKlaim, AddKomiteTreatyChild_ACT 34) memakai SMTP, bukan
// M_LINK_SERVICE; alamat CC/BCC yang di-hardcode XML ke konfigurasi = OQ-CP-15.

import (
	"context"
	"encoding/json"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/claimprop/backend/repository"
)

// PelaksanaClaimProp menjalankan baris outbox Claim Prop.
type PelaksanaClaimProp struct {
	Lingkungan inti.Lingkungan
	Resolver   layanan.ResolverEndpoint
}

// Laksanakan memenuhi `outbox.PelaksanaEfek`.
func (p PelaksanaClaimProp) Laksanakan(ctx context.Context, _ *db.Tx, b outbox.BarisEfekKeluar) error {
	if b.Modul != repository.ModulOutbox {
		return fmt.Errorf("%w: pelaksana Claim Prop menerima baris modul %q", galat.ErrPermintaanTidakSah, b.Modul)
	}
	var m MuatanOutbox
	if err := json.Unmarshal([]byte(b.Muatan), &m); err != nil {
		return fmt.Errorf("%w: muatan outbox Claim Prop tak terbaca: %v", galat.ErrPermintaanTidakSah, err)
	}
	switch b.Jenis {
	case JenisEfekKonversi, JenisEfekKasir, JenisEfekEmailKomite:
	default:
		return fmt.Errorf("%w: jenis efek Claim Prop %q", galat.ErrPermintaanTidakSah, b.Jenis)
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
	if _, err := layanan.AlamatLayanan(ctx, p.Resolver, layanan.KunciLayanan{Kategori1: m.Kategori1, Kategori2: m.Kategori2}); err != nil {
		return err
	}
	if b.Jenis == JenisEfekKasir {
		return outbox.ErrKasirBelumDisetujui
	}
	return outbox.ErrArasapasBelumDisetujui
}

// PekerjaClaimPropOracle menyusun pekerja outbox Claim Prop (tidak dijalankan modul - lihat `modul.go`).
func PekerjaClaimPropOracle(svc inti.Akar, l inti.Lingkungan) *outbox.PekerjaEfek {
	return outbox.NewPekerjaEfekModul(svc, PelaksanaClaimProp{Lingkungan: l,
		Resolver: layanan.ResolverLinkServiceOracle(svc)}, repository.ModulOutbox)
}
