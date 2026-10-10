package repository

// Untuk apa berkas ini: PENOMORAN - padanan tanpa procedure untuk:
//
//	GENERATE_NOCLMTRTYINTEMP (GenerateNoTRTInTemp)  CLMTRTINTEMP_SEQ.NEXTVAL + INSERT OUTTREATYINTEMP   "KT.<yyyy>-<urut>"
//	GENERATE_NOCLMTREATYIN   (GenerateNoCLMTreatyIn) CLMTREATYIN_SEQ.NEXTVAL + INSERT CFSTREATYIN       "RNM-K<bis>.MM.YYYY.T<urut>"
//	PROC_GENERATE_SEQUENCE_NUMBER (GetSequenceNumber_SQL)  inti/backend/penomor (kunci CLASS, JENIS, TAHUN)
//
// Isi procedure dibaca dari ALL_SOURCE DEV 07-10-2026 (`cp/sumber_proc.out`). COMMIT procedure TIDAK ditiru - nomor
// terbit di transaksi aksi, ikut batal bila aksinya batal (nomor sequence tetap dapat berlubang, sama seperti Pega).

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penomor"
	"nusantarare/modul/claimprop/backend/models"
)

func lpad5(s string) string {
	s = strings.TrimSpace(s)
	for len(s) < 5 {
		s = "0" + s
	}
	return s
}

// sqlNomorSementara - catatan nomor sementara (isi GENERATE_NOCLMTRTYINTEMP tanpa COMMIT).
func sqlNomorSementara(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID_CFS, KODE, TAHUN, COUNT) VALUES (:1, :2, :3, :4)`, tabel)
}

// sqlNomorCFS - catatan nomor CFS (isi GENERATE_NOCLMTREATYIN tanpa COMMIT).
func sqlNomorCFS(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID_CFS, KODE, KODE_BIS, BULAN, TAHUN, COUNT) VALUES (:1, :2, :3, :4, :5, :6)`, tabel)
}

// NomorSementara = GENERATE_NOCLMTRTYINTEMP('KT', tahun): `KODE || '.' || TAHUN || '-' || lpad(seq,5,'0')`.
func (g *Gudang) NomorSementara(ctx context.Context, tx *db.Tx, tahun string) (string, error) {
	urut, err := g.db.NomorBerikut(ctx, tx, "CLMTRTINTEMP_SEQ")
	if err != nil {
		return "", err
	}
	n, _ := strconv.Atoi(urut)
	nomor := models.RakitNomorSementara(tahun, n)
	tabel, err := g.db.Qualify("OUTTREATYINTEMP")
	if err != nil {
		return "", err
	}
	q := sqlNomorSementara(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	hasil, err := tx.ExecContext(ctx, q, nomor, models.KodeNomorSementara, tahun, lpad5(urut))
	if err != nil {
		return "", fmt.Errorf("repository: mencatat nomor sementara klaim: %w", err)
	}
	return nomor, db.PastikanSatuBaris(hasil, "pencatatan nomor sementara klaim")
}

// NomorCFS = GENERATE_NOCLMTREATYIN('K', kodeBis, MM, YYYY, 'OUT') (UpdateTableOS langkah 2-3).
func (g *Gudang) NomorCFS(ctx context.Context, tx *db.Tx, kodeBis string, saat time.Time) (string, error) {
	urut, err := g.db.NomorBerikut(ctx, tx, "CLMTREATYIN_SEQ")
	if err != nil {
		return "", err
	}
	n, _ := strconv.Atoi(urut)
	lokal := saat.In(models.Jakarta)
	mm, yyyy := fmt.Sprintf("%02d", int(lokal.Month())), strconv.Itoa(lokal.Year())
	nomor := models.RakitNomorCFS(kodeBis, mm, yyyy, n)
	tabel, err := g.db.Qualify("CFSTREATYIN")
	if err != nil {
		return "", err
	}
	q := sqlNomorCFS(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	hasil, err := tx.ExecContext(ctx, q, nomor, models.JenisNomorKlaim, teksAtauNil(kodeBis), mm, yyyy, lpad5(urut))
	if err != nil {
		return "", fmt.Errorf("repository: mencatat nomor CFS klaim: %w", err)
	}
	return nomor, db.PastikanSatuBaris(hasil, "pencatatan nomor CFS klaim")
}

// BahanNomor - keluaran penghitung bersama (`ParamSeq.HASIL1` MM.YYYY, `ParamSeq.HASIL2` urut).
type BahanNomor struct {
	Jenis  string
	MMYYYY string
	Urut   int
}

// UrutNomor = GetKodeProdNonLife_SQL + PROC_GENERATE_SEQUENCE_NUMBER(pxObjClass, kode+huruf, NULL): jenis = kode
// produksi NONLIFE + huruf ("K" klaim, "H" PLA, "P" DLA), periode dari tanggal tutup buku.
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

// HariClosing - TANGGAL_CLOSING (GETTanggalClosing_SQL, SaveOutstanding_Act langkah 11-12).
func (g *Gudang) HariClosing(ctx context.Context, tx *db.Tx) (int, error) {
	return penomor.NewPenomor(g.db).HariClosing(ctx, tx)
}
