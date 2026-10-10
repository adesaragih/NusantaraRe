package repository

// Untuk apa berkas ini: NOMOR AKSEPTASI - KomitePostAdjustment S16.5-S16.7 tanpa procedure:
//
//	S16.5 GetKodeProdNonLife_SQL       KODE_PRODUKSI TYPE 'NONLIFE'       -> penomor.AwalanProduksi
//	S16.6 CARI1 = pyWorkCover.pxObjClass, CARI2 = kode + "A"
//	S16.7 GetSequenceNumber_SQL        PROC_GENERATE_SEQUENCE_NUMBER(CARI1, CARI2, TO_DATE(CARI3))
//	                                   -> penomor.HariClosing + HitungPeriodeNomor + UrutNomorBerikut
//
// ⚠️ Celah XML: `ParamSeq.CARI3` (tanggal) tidak pernah diisi - S16.1-S16.4 (tanggal tutup buku,
// `GenerateNoAcceptTreaty`) ter-remark. ALL_SOURCE DEV `PROC_GENERATE_SEQUENCE_NUMBER`: tanggal NULL = SYSDATE
// (Jakarta) dengan pergeseran hari tutup buku - sama dengan `penomor.HitungPeriodeNomor`, sama dengan nomor PLA / DLA
// Claim Prop. COMMIT procedure tidak ditiru: nomor terbit di transaksi Submit.

import (
	"context"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penomor"
	"nusantarare/modul/komiteclaimprop/backend/models"
)

// tipeKodeProduksi - `KODE_PRODUKSI.TYPE` (GetKodeProdNonLife_SQL).
const tipeKodeProduksi = "NONLIFE"

// UrutNomorAkseptasi = S16.5-S16.7.
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
