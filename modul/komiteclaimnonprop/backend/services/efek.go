package services

// Untuk apa berkas ini: PELAKSANA OUTBOX Komite Claim Non Prop - menjalankan satu baris T_LOG_SERVICE_RNM
// `MODUL = komiteclaimnonprop` (pola pelaksana Komite Claim Prop, DISALIN bukan diimpor).
//
// ⛔ Non-produksi: baris berhenti `gagal-permanen` dengan `ErrPengirimStubNonProduksi` (nol panggilan keluar).
// Produksi: alamat M_LINK_SERVICE di-resolve sungguhan (kunci VERBATIM dari activity) lalu berhenti terang
// (`…BelumDisetujui`) sampai manusia menyetujui panggilan nyata. Email
// (`SendEmailWithAttachments`) memakai SMTP, bukan M_LINK_SERVICE. Email dirakit lebih dulu (`Penyusun`) - isi yang
// akan dikirim terbukti dapat disusun sebelum panggilan nyatanya ditahan.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/komiteclaimnonprop/backend/models"
	"nusantarare/modul/komiteclaimnonprop/backend/repository"
)

// Penyusun merakit isi efek dari pengenal MUATAN - `*Layanan`.
type Penyusun interface {
	SusunEmailKomite(ctx context.Context, komiteID string, isi map[string]string) (models.SurelKomite, error)
}

// ErrPenyusunBelumDisambung - pelaksana dirakit tanpa perakit isi (salah rakit).
var ErrPenyusunBelumDisambung = errors.New("services: perakit isi efek Komite Claim Non Prop belum disambung")

// PelaksanaKomiteClaimNonProp menjalankan baris outbox modul ini.
type PelaksanaKomiteClaimNonProp struct {
	Lingkungan inti.Lingkungan
	Resolver   layanan.ResolverEndpoint
	Penyusun   Penyusun
}

// Laksanakan memenuhi `outbox.PelaksanaEfek`.
func (p PelaksanaKomiteClaimNonProp) Laksanakan(ctx context.Context, _ *db.Tx, b outbox.BarisEfekKeluar) error {
	if b.Modul != repository.ModulOutbox {
		return fmt.Errorf("%w: pelaksana Komite Claim Non Prop menerima baris modul %q", galat.ErrPermintaanTidakSah, b.Modul)
	}
	var m MuatanOutbox
	if err := json.Unmarshal([]byte(b.Muatan), &m); err != nil {
		return fmt.Errorf("%w: muatan outbox Komite Claim Non Prop tak terbaca: %v", galat.ErrPermintaanTidakSah, err)
	}
	switch b.Jenis {
	case JenisEfekKonversi, JenisEfekKasir, JenisEfekEmailKomite:
	default:
		return fmt.Errorf("%w: jenis efek Komite Claim Non Prop %q", galat.ErrPermintaanTidakSah, b.Jenis)
	}
	if !p.Lingkungan.AdalahProduksi() {
		return outbox.ErrPengirimStubNonProduksi
	}
	if b.Jenis == JenisEfekEmailKomite {
		if p.Penyusun == nil {
			return ErrPenyusunBelumDisambung
		}
		if _, err := p.Penyusun.SusunEmailKomite(ctx, m.KomiteID, isiTeks(m.Isi)); err != nil {
			return err
		}
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

// isiTeks - `MuatanOutbox.Isi` (objek JSON) sebagai teks per kunci.
func isiTeks(v any) map[string]string {
	out := map[string]string{}
	if m, ok := v.(map[string]any); ok {
		for k, x := range m {
			if s, ok := x.(string); ok {
				out[k] = s
			}
		}
	}
	return out
}
