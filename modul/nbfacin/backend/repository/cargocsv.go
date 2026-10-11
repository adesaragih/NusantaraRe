package repository

// Tulisan cargo hasil unggah CSV Marine Cargo (tiket 74 P6, K74-17 tanpa batas baris): pohon cargo PUT objek-cargo (tulisCargo, SATU
// lintasan traversal) dijalankan berkelompok seperti objek FIRE tiket 75 - lintasan hitung ID, ID diambil berkelompok, INSERT
// array-bind per teks SQL (jalankanKelompok: periksaBindBaris + kolomArrayBertipe). Transaksi milik pemanggil.

import (
	"context"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// PenulisCargoCsv - tiket 74 P6.
type PenulisCargoCsv interface {
	// TulisCargoCsv - di transaksi tx: sentuh case (kunci + 404), pastikan General; ganti = seluruh cargo case dihapus dulu
	// (sqlHapusCargo, Property-Remove CargoList InsertUploadMarineCargo_act langkah 2); cargo disisip berkelompok dengan SEQ_NO
	// sesudah SEQ_NO cargo terbesar yang ada.
	TulisCargoCsv(ctx context.Context, tx *db.Tx, id string, cargo []models.ObjekCargo, ganti bool) error
}

func sqlSeqCargoTerakhir(cargo string) string {
	return "SELECT NVL(MAX(SEQ_NO), 0) FROM " + cargo + " WHERE PARENT_ID = :1"
}

// TulisCargoCsv - lihat PenulisCargoCsv.
func (r *ObjekOracle) TulisCargoCsv(ctx context.Context, tx *db.Tx, id string, cargo []models.ObjekCargo, ganti bool) error {
	t, err := r.tabelCargo()
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
		for _, q := range sqlHapusCargo(t) {
			if _, err := jalankan(ctx, tx, q, "menghapus cargo", id); err != nil {
				return err
			}
		}
	}
	var awal int
	if err := tx.QueryRowContext(ctx, sqlSeqCargoTerakhir(t.cargo), id).Scan(&awal); err != nil {
		return fmt.Errorf("repository: SEQ_NO cargo terakhir: %w", err)
	}
	return sisipCargoBerkelompok(ctx, oracleKelompok{d: r.db, tx: tx}, t, id, cargo, uidAcak, awal)
}

// sisipCargoBerkelompok - pohon cargo berkelompok (lihat kepala berkas), SEQ_NO cargo awal+1..
func sisipCargoBerkelompok(ctx context.Context, e eksekutorKelompok, t tabelCargo, id string, cargo []models.ObjekCargo,
	uid func() (string, error), awal int) error {
	hitung := &penghitungID{n: map[string]int{}}
	lewat := func(string, string, []any) error { return nil }
	if err := (&tulisCargo{t: t, nomor: func(s string) (string, error) { return hitung.nomor(ctx, s) },
		uid: func() (string, error) { return "", nil }, sisip: lewat, awal: awal}).sisipCargo(id, cargo); err != nil {
		return err
	}
	kumpulan, err := ambilIDBerkelompok(ctx, e, hitung)
	if err != nil {
		return err
	}
	kel := &penulisKelompok{id: kumpulan, indeks: map[string]*kelompokSisip{}}
	w := &tulisCargo{t: t, nomor: func(s string) (string, error) { return kel.nomor(ctx, s) }, uid: uid,
		sisip: func(q, tabel string, arg []any) error { return kel.sisip(ctx, q, tabel, arg) }, awal: awal}
	if err := w.sisipCargo(id, cargo); err != nil {
		return err
	}
	return jalankanPenulisKelompok(ctx, e, kel)
}
