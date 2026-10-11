package repository

// Simpan perubahan saja tab Object per halaman (tiket 75 F3, PATCH …/objek). Satu transaksi milik services:
//  1. sentuh case (kunci + 404) + pastikan General (sama dengan GantiObjek);
//  2. lokasi `ubah` / `hapus` wajib milik case (ErrLokasiBukanKasus) - SEQ_NO lokasi `ubah` diingat;
//  3. lokasi `ubah` yang coverage-nya ber-spreading DITOLAK (ErrLokasiBerSpreading) - K75-4 (keputusan work owner 10-10-2026, sementara sampai F4): spreading per lokasi
//     belum dapat dipertahankan tanpa membaca case utuh (tab Spreading case besar = F4);
//  4. pohon lokasi `ubah` + `hapus` dihapus (urutan anak -> induk sqlHapusObjek, dibatasi ke ID lokasi);
//  5. `ubah` disisip ulang di SEQ_NO-nya (ID baru);
//  6. `tambah` disisip sesudah SEQ_NO terbesar.
//
// K75-5 (keputusan work owner 10-10-2026): SEQ_NO TIDAK dirapatkan sesudah `hapus`: urutan tetap, celah boleh (pembaca FIRE memakai ORDER BY
// SEQ_NO, bukan nilainya). Merapatkan 1 juta lokasi tanpa MERGE berbiaya kuadrat; PUT …/objek dan Save Data CSV menulis ulang
// 1..n seperti semula. Satu-satunya pembaca yang memakai nilai SEQ_NO lokasi sebagai posisi adalah kunci objek Aneka
// (sqlKunciAnekaLama) - tab Object per halaman hanya untuk lini FIRE.

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// ErrLokasiBukanKasus - ID lokasi PATCH tidak ada di case (sudah dihapus / diganti pengguna lain, atau bukan milik case).
var ErrLokasiBukanKasus = errors.New("repository: lokasi bukan milik case")

// ErrLokasiBerSpreading - lokasi yang diubah sudah memiliki baris spreading.
var ErrLokasiBerSpreading = errors.New("repository: lokasi yang diubah memiliki spreading")

// BatasLokasiPerubahan - ID lokasi per kueri IN (≤ 1000 ekspresi Oracle, ORA-01795).
const BatasLokasiPerubahan = 500

// PenyimpanHalamanObjek - tiket 75 F3: baca per halaman + simpan perubahan saja.
type PenyimpanHalamanObjek interface {
	PembacaHalamanObjek
	// PerbaruiObjek - lihat kepala berkas; mengembalikan jumlah lokasi case sesudahnya.
	PerbaruiObjek(ctx context.Context, tx *db.Tx, id string, ubah []models.LokasiUbah, tambah []models.ObjekFire, hapus []string) (int, error)
}

// bindIN - ":awal, :awal+1, …" sebanyak n.
func bindIN(awal, n int) string {
	b := make([]string, n)
	for i := range b {
		b[i] = ":" + strconv.Itoa(awal+i)
	}
	return strings.Join(b, ", ")
}

// sqlSeqLokasi - ID + SEQ_NO lokasi case :1 yang ID-nya di daftar (:2..).
func sqlSeqLokasi(loc string, n int) string {
	return "SELECT TO_CHAR(ID), SEQ_NO FROM " + loc + " WHERE PARENT_ID = :1 AND ID IN (" + bindIN(2, n) + ")"
}

// sqlAdaSpreadingLokasi - ada baris spreading di coverage lokasi terpilih (lingkupBaca menyisipkan ID lokasi).
func sqlAdaSpreadingLokasi(t tabelObjek) string {
	return "SELECT COUNT(*) FROM " + t.spread + " WHERE PARENT_ID IN (" +
		sqlCoverageKasus(tabelSpreading{cov: t.cov, item: t.item, prop: t.prop, loc: t.loc}) + ") AND ROWNUM = 1"
}

// sqlHapusLokasi - sqlHapusObjek dibatasi ke ID lokasi (:2..:n+1) case :1. Setiap pernyataan sqlHapusObjek memuat tepat satu
// syarat case: `l.PARENT_ID = :1` (beralias) atau `PARENT_ID = :1` atas T_LOCATIONLIST tanpa alias.
func sqlHapusLokasi(t tabelObjek, n int) ([]string, error) {
	var hasil []string
	for _, q := range sqlHapusObjek(t) {
		if strings.Count(q, "PARENT_ID = :1") != 1 {
			return nil, fmt.Errorf("repository: pernyataan hapus objek tanpa syarat case tunggal: %s", q)
		}
		saring := " AND ID IN (" + bindIN(2, n) + ")"
		if strings.Contains(q, pembatasLokasi) {
			saring = " AND l.ID IN (" + bindIN(2, n) + ")"
			q = strings.Replace(q, pembatasLokasi, pembatasLokasi+saring, 1)
		} else {
			lokasi := "FROM " + t.loc + " WHERE PARENT_ID = :1"
			if strings.Count(q, lokasi) != 1 {
				return nil, fmt.Errorf("repository: pernyataan hapus objek tidak dikenali: %s", q)
			}
			q = strings.Replace(q, lokasi, lokasi+saring, 1)
		}
		hasil = append(hasil, q)
	}
	return hasil, nil
}

func sqlSeqTerakhir(loc string) string {
	return "SELECT NVL(MAX(SEQ_NO), 0), COUNT(*) FROM " + loc + " WHERE PARENT_ID = :1"
}

// argLokasi - :1 case lalu ID lokasi.
func argLokasi(id string, lokasi []string) []any {
	arg := []any{id}
	for _, x := range lokasi {
		arg = append(arg, x)
	}
	return arg
}

// PerbaruiObjek - lihat PenyimpanHalamanObjek.
func (r *ObjekOracle) PerbaruiObjek(ctx context.Context, tx *db.Tx, id string, ubah []models.LokasiUbah, tambah []models.ObjekFire,
	hapus []string) (int, error) {
	t, err := r.tabel()
	if err != nil {
		return 0, err
	}
	sentuh, err := jalankan(ctx, tx, sqlSentuhCase(t.work), "menyentuh "+TabelWorkPolis, id, LiniFacIn)
	if err != nil {
		return 0, err
	}
	if n, err := sentuh.RowsAffected(); err != nil {
		return 0, err
	} else if n == 0 {
		return 0, ErrKasusTidakAda
	}
	if _, err := jalankan(ctx, tx, sqlPastikanGeneral(t.general), "memastikan "+TabelGeneralPolis, id, id); err != nil {
		return 0, err
	}
	idUbah := make([]string, len(ubah))
	for i, u := range ubah {
		idUbah[i] = u.ID
	}
	semua := append(append([]string{}, idUbah...), hapus...)
	if len(semua) > BatasLokasiPerubahan {
		return 0, fmt.Errorf("repository: %d lokasi diubah / dihapus (paling banyak %d)", len(semua), BatasLokasiPerubahan)
	}
	seq := map[string]int{}
	if len(semua) > 0 {
		if seq, err = r.seqLokasi(ctx, tx, t, id, semua); err != nil {
			return 0, err
		}
	}
	if len(idUbah) > 0 {
		q, arg, err := lingkupBaca{id: id, lokasi: idUbah}.kueri(sqlAdaSpreadingLokasi(t))
		if err != nil {
			return 0, err
		}
		var n int
		if err := tx.QueryRowContext(ctx, q, arg...).Scan(&n); err != nil {
			return 0, fmt.Errorf("repository: cek spreading lokasi: %w", err)
		}
		if n > 0 {
			return 0, ErrLokasiBerSpreading
		}
	}
	if len(semua) > 0 {
		hapusQ, err := sqlHapusLokasi(t, len(semua))
		if err != nil {
			return 0, err
		}
		for _, q := range hapusQ {
			if _, err := jalankan(ctx, tx, q, "menghapus lokasi objek", argLokasi(id, semua)...); err != nil {
				return 0, err
			}
		}
	}
	e := oracleKelompok{d: r.db, tx: tx}
	if len(ubah) > 0 {
		baris, posisi := make([]models.ObjekFire, len(ubah)), make([]int, len(ubah))
		for i, u := range ubah {
			baris[i], posisi[i] = u.Objek, seq[u.ID]
		}
		if err := sisipBerkelompokPosisi(ctx, e, t, id, baris, uidAcak, 0, posisi); err != nil {
			return 0, err
		}
	}
	var akhir, jumlah int
	if err := tx.QueryRowContext(ctx, sqlSeqTerakhir(t.loc), id).Scan(&akhir, &jumlah); err != nil {
		return 0, fmt.Errorf("repository: SEQ_NO lokasi terakhir: %w", err)
	}
	if len(tambah) > 0 {
		if err := sisipBerkelompokPosisi(ctx, e, t, id, tambah, uidAcak, akhir, nil); err != nil {
			return 0, err
		}
	}
	return jumlah + len(tambah), nil
}

// seqLokasi - SEQ_NO per ID lokasi; ErrLokasiBukanKasus bila ada ID yang tidak ada di case.
func (r *ObjekOracle) seqLokasi(ctx context.Context, tx *db.Tx, t tabelObjek, id string, lokasi []string) (map[string]int, error) {
	baris, err := tx.QueryContext(ctx, sqlSeqLokasi(t.loc, len(lokasi)), argLokasi(id, lokasi)...)
	if err != nil {
		return nil, fmt.Errorf("repository: SEQ_NO lokasi: %w", err)
	}
	defer baris.Close()
	seq := map[string]int{}
	for baris.Next() {
		var x string
		var n int
		if err := baris.Scan(&x, &n); err != nil {
			return nil, fmt.Errorf("repository: SEQ_NO lokasi: %w", err)
		}
		seq[x] = n
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: SEQ_NO lokasi: %w", err)
	}
	for _, x := range lokasi {
		if _, ada := seq[x]; !ada {
			return nil, fmt.Errorf("%w: %s", ErrLokasiBukanKasus, x)
		}
	}
	return seq, nil
}
