package services

// Untuk apa berkas ini: EFEK KELUAR Claim Non Prop - antrean ke T_LOG_SERVICE_RNM (hanya produksi, `When/IsPEGAPROD`)
// dan PELAKSANA baris outbox `MODUL = claimnonprop` (pola Claim Prop `PelaksanaClaimProp`, disalin bukan diimpor).
//
// ⛔ Non-produksi: pengirim STUB - nol panggilan keluar, baris berhenti `gagal-permanen` dengan
// `ErrPengirimStubNonProduksi`. Produksi: alamat M_LINK_SERVICE di-resolve (kunci Kategori_1/Kategori_2 VERBATIM)
// lalu berhenti terang (`…BelumDisetujui`) sampai manusia menyetujui panggilan nyata (prompt §6 butir 6). Connector
// bertautan mati (InsertClaimOutstanding_NP, insertClaimFinalOrClosed_NP) tanpa kunci M_LINK_SERVICE = OQ-CNP-31.

import (
	"context"
	"encoding/json"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/claimnonprop/backend/models"
	"nusantarare/modul/claimnonprop/backend/repository"
)

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
		"NOPOLIS": j.h.Ambil(models.CD + "PolicyData.PolicyNo"), "STS_REJECT": sts})
}

// PelaksanaClaimNonProp menjalankan baris outbox Claim Non Prop.
type PelaksanaClaimNonProp struct {
	Lingkungan inti.Lingkungan
	Resolver   layanan.ResolverEndpoint
}

// Laksanakan memenuhi `outbox.PelaksanaEfek`.
func (p PelaksanaClaimNonProp) Laksanakan(ctx context.Context, _ *db.Tx, b outbox.BarisEfekKeluar) error {
	if b.Modul != repository.ModulOutbox {
		return fmt.Errorf("%w: pelaksana Claim Non Prop menerima baris modul %q", galat.ErrPermintaanTidakSah, b.Modul)
	}
	var m MuatanOutbox
	if err := json.Unmarshal([]byte(b.Muatan), &m); err != nil {
		return fmt.Errorf("%w: muatan outbox Claim Non Prop tak terbaca: %v", galat.ErrPermintaanTidakSah, err)
	}
	switch b.Jenis {
	case JenisEfekOutstanding, JenisEfekKonversi, JenisEfekTutup, JenisEfekKasir, JenisEfekEmailKomite:
	default:
		return fmt.Errorf("%w: %w %q", galat.ErrPermintaanTidakSah, ErrEfekTakDikenal, b.Jenis)
	}
	if !p.Lingkungan.AdalahProduksi() {
		return outbox.ErrPengirimStubNonProduksi
	}
	if b.Jenis == JenisEfekEmailKomite {
		return outbox.EfekEmail{}.Jalankan(ctx, outbox.MuatanEfek{KlaimID: m.KlaimID})
	}
	if m.Kategori1 == "" { // connector bertautan mati tanpa kunci M_LINK_SERVICE (OQ-CNP-31)
		return outbox.ErrArasapasBelumDisetujui
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

// PekerjaClaimNonPropOracle menyusun pekerja outbox Claim Non Prop (tidak dijalankan modul - lihat `modul.go`).
func PekerjaClaimNonPropOracle(svc inti.Akar, l inti.Lingkungan) *outbox.PekerjaEfek {
	return outbox.NewPekerjaEfekModul(svc, PelaksanaClaimNonProp{Lingkungan: l,
		Resolver: layanan.ResolverLinkServiceOracle(svc)}, repository.ModulOutbox)
}
