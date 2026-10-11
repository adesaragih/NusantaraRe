package services

// Untuk apa berkas ini: LAMPIRAN KLAIM untuk layar Komite Claim Fac In - bagian `kontrak.KlaimFacInKomite` (perintah
// work owner 11-10-2026: "TAMBAHKAN ATTACHMENT KLAIM FACIN KE KOMITE KLAIM FACIN"). Hanya-baca, jalur baca yang sama
// dengan tab lampiran Claim Fac In (`DaftarLampiran`, `UnduhLampiran`, `TautanOfficeLampiran`); pelaku sudah diperiksa
// modul Komite.

import (
	"context"
	"strings"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/claimfacin/backend/models"
)

// LampiranKlaimFacIn - lihat `kontrak.KlaimFacInKomite`.
func (k KlaimUntukKomite) LampiranKlaimFacIn(ctx context.Context, klaimID string) (kontrak.DaftarLampiranKlaim, error) {
	if err := k.siap(); err != nil {
		return kontrak.DaftarLampiranKlaim{}, err
	}
	kasus, err := k.l.g.Keadaan(ctx, nil, klaimID)
	if err != nil {
		return kontrak.DaftarLampiranKlaim{}, galatKasusKomite(err)
	}
	kat, err := k.l.g.KategoriLampiran(ctx, kasus.ID)
	if err != nil {
		return kontrak.DaftarLampiranKlaim{}, err
	}
	dok, err := k.l.g.DaftarLampiran(ctx, kasus.ID)
	if err != nil {
		return kontrak.DaftarLampiranKlaim{}, err
	}
	out := kontrak.DaftarLampiranKlaim{Kategori: make([]kontrak.KategoriLampiranKlaim, 0, len(kat)),
		Lampiran: make([]kontrak.LampiranKlaim, 0, len(dok))}
	for _, x := range kat {
		out.Kategori = append(out.Kategori, kontrak.KategoriLampiranKlaim{ID: x.ID, Label: x.Label, Jumlah: x.CountAttach})
	}
	for _, a := range dok {
		out.Lampiran = append(out.Lampiran, kontrak.LampiranKlaim{ID: a.ID, NamaFile: a.NamaFile, Kategori: a.Kategori,
			MIME: a.MIME, Tanggal: a.Tanggal, Operator: a.Operator, AdaObjek: a.AdaObjek})
	}
	return out, nil
}

// lampiranUntukKomite - satu dokumen kasus klaim `klaimID` (tanpa pemeriksaan pelaku - urusan modul Komite).
func (k KlaimUntukKomite) lampiranUntukKomite(ctx context.Context, klaimID, lampiranID string) (models.Lampiran, error) {
	if err := k.siap(); err != nil {
		return models.Lampiran{}, err
	}
	kasus, err := k.l.g.Keadaan(ctx, nil, klaimID)
	if err != nil {
		return models.Lampiran{}, galatKasusKomite(err)
	}
	dok, err := k.l.g.DaftarLampiran(ctx, kasus.ID)
	if err != nil {
		return models.Lampiran{}, err
	}
	for _, a := range dok {
		if a.ID == strings.TrimSpace(lampiranID) {
			return a, nil
		}
	}
	return models.Lampiran{}, kontrak.ErrLampiranKlaimTidakAda
}

// BukaLampiranKlaimFacIn - lihat `kontrak.KlaimFacInKomite` (sama dengan `UnduhLampiran`).
func (k KlaimUntukKomite) BukaLampiranKlaimFacIn(ctx context.Context, klaimID, lampiranID, pengguna string) (
	kontrak.IsiLampiranKlaim, error) {
	a, err := k.lampiranUntukKomite(ctx, klaimID, lampiranID)
	if err != nil {
		return kontrak.IsiLampiranKlaim{}, err
	}
	if !a.AdaObjek {
		return kontrak.IsiLampiranKlaim{}, penyimpanan.ErrObjekTidakAda
	}
	if k.l.berkas == nil {
		return kontrak.IsiLampiranKlaim{}, ErrPenyimpananBelumDipasang
	}
	isi, err := k.l.berkas.Buka(ctx, a.StorageID, models.DurasiLihatLampiran, pengguna)
	if err != nil {
		return kontrak.IsiLampiranKlaim{}, err
	}
	return kontrak.IsiLampiranKlaim{Isi: isi, Nama: a.NamaFile, Mime: unggah.MimeDariNamaFile(a.NamaFile)}, nil
}

// TautanOfficeLampiranKlaimFacIn - lihat `kontrak.KlaimFacInKomite` (sama dengan `TautanOfficeLampiran`).
func (k KlaimUntukKomite) TautanOfficeLampiranKlaimFacIn(ctx context.Context, klaimID, lampiranID, pengguna string) (
	string, error) {
	a, err := k.lampiranUntukKomite(ctx, klaimID, lampiranID)
	if err != nil {
		return "", err
	}
	if !models.EkstensiOffice[strings.ToLower(a.MIME)] {
		return "", kontrak.ErrLampiranKlaimBukanOffice
	}
	if !a.AdaObjek {
		return "", penyimpanan.ErrObjekTidakAda
	}
	if k.l.berkas == nil {
		return "", ErrPenyimpananBelumDipasang
	}
	return k.l.berkas.Tautan(ctx, a.StorageID, models.DurasiLihatLampiran, pengguna)
}
