package repository

// Untuk apa berkas ini: PENOMORAN - padanan tanpa procedure untuk:
//
//	PROC_GENERATE_SEQUENCE_NUMBER (GetSequenceNumber_SQL, SaveDataToOSAksep_Act 15.6)  inti/backend/penomor
//	GenerateNoPLATNP (GeneratePlaCNP_Act 8)  'RNM-M' || OLDID || '.' || MM || '.' || yyyy || '.TX' || lpad(PLATNP_SEQ, 5)
//
// Nomor terbit di transaksi aksi, ikut batal bila aksinya batal (nomor sequence tetap dapat berlubang, sama seperti
// Pega). Tanggal `TO_DATE(CARI3)` procedure tidak pernah diisi korpus (OQ-CNP-27) - periode dari tanggal tutup buku,
// pola Claim Prop.

import (
	"context"
	"strconv"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penomor"
	"nusantarare/modul/claimnonprop/backend/models"
)

// BahanNomor - keluaran penghitung bersama (`ParamSeq.HASIL1` MM.YYYY, `ParamSeq.HASIL2` urut).
type BahanNomor struct {
	Jenis  string
	MMYYYY string
	Urut   int
}

// UrutNomor = GetKodeProdNonLife_SQL + PROC_GENERATE_SEQUENCE_NUMBER(pxObjClass, kode+huruf, NULL): jenis = kode
// produksi NONLIFE + huruf ("K" klaim), periode dari tanggal tutup buku.
func (g *Gudang) UrutNomor(ctx context.Context, tx *db.Tx, huruf string, saat time.Time) (BahanNomor, error) {
	p := penomor.NewPenomor(g.db)
	kode, err := p.AwalanProduksi(ctx, tx, models.TipeKodeProduksiNonLife)
	if err != nil {
		return BahanNomor{}, err
	}
	hari, err := p.HariClosing(ctx, tx)
	if err != nil {
		return BahanNomor{}, err
	}
	per, err := penomor.HitungPeriodeNomor(saat, hari)
	if err != nil {
		return BahanNomor{}, err
	}
	jenis := kode + huruf
	urut, err := p.UrutNomorBerikut(ctx, tx, models.KelasKasus, jenis, per, saat)
	if err != nil {
		return BahanNomor{}, err
	}
	return BahanNomor{Jenis: jenis, MMYYYY: per.MMYYYY, Urut: urut}, nil
}

// NomorPLA = GenerateNoPLATNP: urut PLATNP_SEQ, nomor dirakit `models.RakitNomorPLA` (tanggal sistem).
func (g *Gudang) NomorPLA(ctx context.Context, tx *db.Tx, oldID string, saat time.Time) (string, error) {
	urut, err := g.db.NomorBerikut(ctx, tx, "PLATNP_SEQ")
	if err != nil {
		return "", err
	}
	n, err := strconv.Atoi(urut)
	if err != nil {
		return "", err
	}
	return models.RakitNomorPLA(oldID, saat, n), nil
}
