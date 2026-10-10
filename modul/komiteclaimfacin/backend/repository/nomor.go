package repository

// Untuk apa berkas ini: NOMOR AKSEPTASI - KomitePost_Adjustment S7.2.1.2.4-S7.2.1.2.7 tanpa procedure (pola Komite
// Claim Prop / Non Prop, Claim Fac In tahap 1):
//
//	S7.2.1.2.4  GetKodeProdNonLife_SQL     KODE_PRODUKSI TYPE 'NONLIFE'        -> penomor.AwalanProduksi
//	S7.2.1.2.5  CARI1 = pyWorkCover.pxObjClass, CARI2 = kode + "A"
//	S7.2.1.2.6  GetSequenceNumber_SQL      PROC_GENERATE_SEQUENCE_NUMBER(CARI1, CARI2, TO_DATE(CARI3))
//	                                       -> penomor.HariClosing + HitungPeriodeNomor + UrutNomorBerikut
//
// Celah XML: `ParamSeq.CARI3` (tanggal) tidak pernah diisi - S7.2.1.2.1-S7.2.1.2.3 ter-remark. Ditangani sama dengan
// Claim Fac In tahap 1 (prompt tahap 1 §5 butir 11): periode dari tanggal tutup buku. COMMIT procedure tidak ditiru:
// nomor terbit di transaksi Submit.

import (
	"context"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penomor"
	"nusantarare/modul/komiteclaimfacin/backend/models"
)

// tipeKodeProduksi - `KODE_PRODUKSI.TYPE` (GetKodeProdNonLife_SQL).
const tipeKodeProduksi = "NONLIFE"

// UrutNomorAkseptasi = S7.2.1.2.4-S7.2.1.2.6.
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
