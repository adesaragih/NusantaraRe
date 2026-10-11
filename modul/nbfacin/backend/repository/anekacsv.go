package repository

// Tulisan objek Aneka hasil unggah CSV (tiket 73 P6, keputusan work owner 11-10-2026; infrastruktur tiket 75 / 74 P6): pohon tab
// Object Aneka (tulisAneka, SATU lintasan traversal yang sama dengan PUT objek-aneka) dijalankan berkelompok dengan SEQ_NO lokasi
// sesudah lokasi terbesar yang ada. Transaksi milik pemanggil.

import (
	"context"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// PenulisAnekaCsv - tiket 73 P6.
type PenulisAnekaCsv interface {
	// TulisAnekaCsv - di transaksi tx: sentuh case (kunci + 404), pastikan General; ganti = seluruh objek case dihapus dulu
	// (InsertUploadAneka_act langkah 1 Property-Remove LocationList): coverage Aneka objek lama (beserta deductible / spreading),
	// halaman anak + objek + okupasi Aneka, lalu seluruh lokasi (sqlHapusObjek FIRE). Lokasi disisip berkelompok dengan SEQ_NO
	// sesudah SEQ_NO lokasi terbesar yang ada.
	TulisAnekaCsv(ctx context.Context, tx *db.Tx, id string, lokasi []models.LokasiAneka, ganti bool) error
}

func sqlSeqLokasiTerakhir(loc string) string {
	return "SELECT NVL(MAX(SEQ_NO), 0) FROM " + loc + " WHERE PARENT_ID = :1"
}

// TulisAnekaCsv - lihat PenulisAnekaCsv.
func (r *ObjekOracle) TulisAnekaCsv(ctx context.Context, tx *db.Tx, id string, lokasi []models.LokasiAneka, ganti bool) error {
	t, err := r.tabelAneka()
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, sqlSentuhCase(t.work), "menyentuh "+TabelWorkPolis, id, LiniFacIn)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrKasusTidakAda
	}
	if _, err := jalankan(ctx, tx, sqlPastikanGeneral(t.general), "memastikan "+TabelGeneralPolis, id, id); err != nil {
		return err
	}
	if ganti {
		lama, err := kunciAnekaLama(ctx, tx, t, id)
		if err != nil {
			return err
		}
		_, hapus := rencanaCoverageAneka(lama, nil, nil)
		for _, idLama := range hapus {
			for _, q := range sqlHapusCoverageAneka(t) {
				if _, err := jalankan(ctx, tx, q, "menghapus coverage Aneka", idLama); err != nil {
					return err
				}
			}
		}
		for _, h := range sqlHapusAneka(t) {
			arg := make([]any, h.bind)
			for i := range arg {
				arg[i] = id
			}
			if _, err := jalankan(ctx, tx, h.q, "menghapus objek Aneka", arg...); err != nil {
				return err
			}
		}
		for _, q := range sqlHapusObjek(t.tabelObjek) {
			if _, err := jalankan(ctx, tx, q, "menghapus objek", id); err != nil {
				return err
			}
		}
	}
	q := sqlSeqLokasiTerakhir(t.loc)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	var awal int
	if err := tx.QueryRowContext(ctx, q, id).Scan(&awal); err != nil {
		return fmt.Errorf("repository: SEQ_NO lokasi terakhir: %w", err)
	}
	_, _, err = sisipAnekaBerkelompokDari(ctx, oracleKelompok{d: r.db, tx: tx}, t, id, lokasi, uidAcak, awal)
	return err
}
