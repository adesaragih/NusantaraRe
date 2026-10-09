package services

// Untuk apa berkas ini: EFEK KELUAR Claim Fac In - antrean ke T_LOG_SERVICE_RNM (hanya produksi, `When/IsPEGAPROD`)
// dan PELAKSANA baris outbox `MODUL = claimfacin` (pola Claim Non Prop `PelaksanaClaimNonProp`, disalin bukan diimpor).
//
// ⛔ Non-produksi: pengirim STUB - nol panggilan keluar, baris berhenti `gagal-permanen` dengan
// `ErrPengirimStubNonProduksi`. Produksi: alamat M_LINK_SERVICE di-resolve (kunci Kategori_1/Kategori_2 VERBATIM)
// lalu berhenti terang (`…BelumDisetujui`) sampai manusia menyetujui panggilan nyata (prompt §6).

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
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/repository"
)

// Jenis efek keluar Claim Fac In.
const (
	JenisEfekKonversi = "konversi-klaim" // KonversiKlaim_Act -> KonversiKlaimNonLife (Klaim / insertClaimAccept)
	JenisEfekKasir    = "kasir"          // HitServiceToKasir_Act -> SendAcceptationToKasir (Kasir / insertAllPaymentKasir)
	JenisEfekEmail    = "email-klaim"    // SendEmailDLA_ACT / SendEmailKlaim
	JenisEfekDLA      = "dla-klaim"      // ChooseDla_Act 9-10 -> HitDLAClaimFacin (Klaim / insertClaimDLA)
)

// MuatanOutbox - isi baris T_LOG_SERVICE_RNM.MUATAN (teks JSON milik outbox inti). Hanya pengenal - nol data pribadi.
type MuatanOutbox struct {
	Kategori1 string `json:"kategori1,omitempty"`
	Kategori2 string `json:"kategori2,omitempty"`
	KlaimID   string `json:"klaimId"`
	Isi       any    `json:"isi,omitempty"`
}

// kunciEfek - kunci M_LINK_SERVICE per jenis (Kategori_1/Kategori_2 VERBATIM dari activity).
var kunciEfek = map[string][2]string{
	JenisEfekKonversi: {"Klaim", "insertClaimAccept"},
	JenisEfekKasir:    {"Kasir", "insertAllPaymentKasir"},
	JenisEfekDLA:      {"Klaim", "insertClaimDLA"},
}

// ErrEfekTakDikenal - jenis efek tanpa penanganan.
var ErrEfekTakDikenal = errors.New("services: jenis efek outbox tidak dikenal")

// antre - efek keluar hanya di produksi.
func (j *jalanAksi) antre(jenis, rujukan string, isi any) error {
	if !j.l.produksi {
		return nil
	}
	k := kunciEfek[jenis]
	teks, err := json.Marshal(MuatanOutbox{Kategori1: k[0], Kategori2: k[1], KlaimID: j.kasus.ID, Isi: isi})
	if err != nil {
		return err
	}
	_, err = j.l.g.AntreEfek(j.ctx, j.tx, jenis, rujukan, string(teks), j.k.Sekarang)
	return err
}

// antreKonversi = KonversiKlaim_Act (INSKEY, NOPOLIS, STSREJECT).
func (j *jalanAksi) antreKonversi(sts string) error {
	return j.antre(JenisEfekKonversi, j.kasus.ID, map[string]string{"CASEID": models.KunciInstans(j.kasus.ID),
		"NOPOLIS": j.h.Ambil(models.JalurNoPolis), "STS_REJECT": sts})
}

// PelaksanaClaimFacIn menjalankan baris outbox Claim Fac In.
type PelaksanaClaimFacIn struct {
	Lingkungan inti.Lingkungan
	Resolver   layanan.ResolverEndpoint
}

// Laksanakan memenuhi `outbox.PelaksanaEfek`.
func (p PelaksanaClaimFacIn) Laksanakan(ctx context.Context, _ *db.Tx, b outbox.BarisEfekKeluar) error {
	if b.Modul != repository.ModulOutbox {
		return fmt.Errorf("%w: pelaksana Claim Fac In menerima baris modul %q", galat.ErrPermintaanTidakSah, b.Modul)
	}
	var m MuatanOutbox
	if err := json.Unmarshal([]byte(b.Muatan), &m); err != nil {
		return fmt.Errorf("%w: muatan outbox Claim Fac In tak terbaca: %v", galat.ErrPermintaanTidakSah, err)
	}
	switch b.Jenis {
	case JenisEfekKonversi, JenisEfekKasir, JenisEfekEmail, JenisEfekDLA:
	default:
		return fmt.Errorf("%w: %w %q", galat.ErrPermintaanTidakSah, ErrEfekTakDikenal, b.Jenis)
	}
	if !p.Lingkungan.AdalahProduksi() {
		return outbox.ErrPengirimStubNonProduksi
	}
	if b.Jenis == JenisEfekEmail {
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

// PekerjaClaimFacInOracle menyusun pekerja outbox Claim Fac In (tidak dijalankan modul - lihat `modul.go`).
func PekerjaClaimFacInOracle(svc inti.Akar, l inti.Lingkungan) *outbox.PekerjaEfek {
	return outbox.NewPekerjaEfekModul(svc, PelaksanaClaimFacIn{Lingkungan: l,
		Resolver: layanan.ResolverLinkServiceOracle(svc)}, repository.ModulOutbox)
}
