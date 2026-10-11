package repository

// Untuk apa berkas ini: PENOMORAN - padanan tanpa procedure untuk PROC_GENERATE_SEQUENCE_NUMBER (GetSequenceNumber_SQL:
// CLaimFaceSheet_Act 14 klaim "K", GeneratePLATreaty_Act 7 PLA treaty "G", GenerateDLAFacin_Act / DLAFacintoTreaty_Act
// DLA "P" / "S") lewat `inti/backend/penomor`.
//
// Nomor terbit di transaksi aksi, ikut batal bila aksinya batal (nomor sequence tetap dapat berlubang, sama seperti
// Pega). PERBAIKAN prompt §5 butir 11 (PARITAS `[penyimpangan sadar]`): `ParamSeq.CARI3` (tanggal procedure) tidak
// pernah diisi pemanggil mana pun - periode dari tanggal tutup buku, pola Claim Prop / Claim Non Prop.

import (
	"context"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penomor"
	"nusantarare/modul/claimfacin/backend/models"
)

// BahanNomor - keluaran penghitung bersama (`ParamSeq.HASIL1` MM.YYYY, `ParamSeq.HASIL2` urut).
type BahanNomor struct {
	Jenis  string
	MMYYYY string
	Urut   int
}

// UrutNomor = GetKodeProdNonLife_SQL + PROC_GENERATE_SEQUENCE_NUMBER(pxObjClass, kode+huruf, NULL): jenis = kode
// produksi NONLIFE + huruf ("K" klaim, "G" PLA treaty, "P" / "S" DLA), periode dari tanggal tutup buku.
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
