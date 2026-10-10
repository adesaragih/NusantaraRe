package repository

// Untuk apa berkas ini: NOMOR AKSEPTASI - KomitePostAdjustment S14.8-S14.10 tanpa procedure (pola Komite Claim Prop):
//
//	S14.8  GetKodeProdNonLife_SQL      KODE_PRODUKSI TYPE 'NONLIFE'       -> penomor.AwalanProduksi
//	S14.9  CARI1 = TempMainWork.pxObjClass, CARI2 = kode + "A"
//	S14.10 GetSequenceNumber_SQL       PROC_GENERATE_SEQUENCE_NUMBER(CARI1, CARI2, TO_DATE(CARI3))
//	                                   -> penomor.HariClosing + HitungPeriodeNomor + UrutNomorBerikut
//
// Celah XML: `ParamSeq.CARI3` (tanggal) tidak pernah diisi - S14.3-S14.7 ter-remark. ALL_SOURCE DEV
// `PROC_GENERATE_SEQUENCE_NUMBER`: tanggal NULL = SYSDATE (Jakarta) dengan pergeseran hari tutup buku - sama dengan
// `penomor.HitungPeriodeNomor`. COMMIT procedure tidak ditiru: nomor terbit di transaksi Submit.

import (
	"context"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penomor"
	"nusantarare/modul/komiteclaimnonprop/backend/models"
)

// tipeKodeProduksi - `KODE_PRODUKSI.TYPE` (GetKodeProdNonLife_SQL).
const tipeKodeProduksi = "NONLIFE"

// UrutNomorAkseptasi = S14.8-S14.10.
func (g *Gudang) UrutNomorAkseptasi(ctx context.Context, tx *db.Tx, saat time.Time) (models.BahanNomor, error) {
	p := penomor.NewPenomor(g.db)
	kode, err := p.AwalanProduksi(ctx, tx, tipeKodeProduksi)
	if err != nil {
		return models.BahanNomor{}, err
	}
	hari, err := p.HariClosing(ctx, tx)
	if err != nil {
		return models.BahanNomor{}, err
	}
	per, err := penomor.HitungPeriodeNomor(saat, hari)
	if err != nil {
		return models.BahanNomor{}, err
	}
	jenis := kode + models.HurufAkseptasi
	urut, err := p.UrutNomorBerikut(ctx, tx, models.KelasKlaim, jenis, per, saat)
	if err != nil {
		return models.BahanNomor{}, err
	}
	return models.BahanNomor{Jenis: jenis, MMYYYY: per.MMYYYY, Urut: urut}, nil
}
