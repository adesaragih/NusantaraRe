package repository

// Daftar Ceding Co - popup Change Ceding Co (tiket 34). Penyimpanan = tabel RANCANGAN
// T_CEDINGCOLIST (jalur `QuotationData/CedingCoList`, induk T_QUOTATIONDATA; migrasi 185) +
// gabungan `;` kode dan nama di T_QUOTATIONDATA.CEDING_CO / CEDING_CO_NAME (Pega
// `Quotation.CedingCo` / `CedingCoName`, `SetCedingCo_Act`: `@If(x=="", .CedingCo, x+";"+.CedingCo)`).
// Kode dicek ke AGENT dengan syarat tiket 33 (CedingCoHierarki memakai RD yang sama,
// BrowseAgentNonLife_RD); nama dari AGENT, bukan dari klien.

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelCedingCoList - daftar ceding per case (migrasi 185).
	TabelCedingCoList = "T_CEDINGCOLIST"
	// SequenceCedingCoList - pembangkit T_CEDINGCOLIST.ID (migrasi 185).
	SequenceCedingCoList = "SEQ_T_CEDINGCOLIST"
	// pemisahCeding - pemisah gabungan kode/nama (E-8, tanpa spasi).
	pemisahCeding = ";"
	// Lebar kolom (migrasi 185 / rancangan, butir 80).
	lebarNamaCeding   = 500  // T_CEDINGCOLIST.CEDING_CO_NAME
	lebarGabunganKode = 1000 // T_QUOTATIONDATA.CEDING_CO
	lebarGabunganNama = 4000 // T_QUOTATIONDATA.CEDING_CO_NAME
)

// ErrCedingTidakSah - ada kode ceding yang tidak ada di AGENT atau tidak lolos syarat tiket 33.
var ErrCedingTidakSah = errors.New("repository: kode ceding tidak ditemukan di AGENT atau tidak lolos syarat")

// ErrCedingTerlaluPanjang - nama dari AGENT atau gabungan `;` melebihi lebar kolomnya.
var ErrCedingTerlaluPanjang = errors.New("repository: daftar ceding melebihi lebar kolom")

func sqlBacaCeding(ceding, quo string) string {
	return "SELECT c.CEDING_CO, c.CEDING_CO_NAME FROM " + ceding + " c JOIN " + quo +
		" q ON q.ID = c.PARENT_ID WHERE q.PARENT_ID = :1 ORDER BY c.SEQ_NO"
}

func sqlIDQuotation(quo string) string { return "SELECT ID FROM " + quo + " WHERE PARENT_ID = :1" }

func sqlHapusCeding(ceding string) string { return "DELETE FROM " + ceding + " WHERE PARENT_ID = :1" }

func sqlSisipCeding(ceding string) string {
	return "INSERT INTO " + ceding + " (ID, PARENT_ID, SEQ_NO, ROW_UID, CEDING_CO, CEDING_CO_NAME) VALUES (:1, :2, :3, :4, :5, :6)"
}

func sqlGabunganCeding(quo string) string {
	return "UPDATE " + quo + " SET CEDING_CO = :1, CEDING_CO_NAME = :2 WHERE ID = :3"
}

// bacaCeding - daftar ceding case `id` urut SEQ_NO (kosong = larik kosong).
func (r *KasusOracle) bacaCeding(ctx context.Context, ceding, quo, id string) ([]models.Ceding, error) {
	baris, err := r.db.QueryContext(ctx, sqlBacaCeding(ceding, quo), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelCedingCoList, err)
	}
	defer baris.Close()
	hasil := []models.Ceding{}
	for baris.Next() {
		var kode, nama sql.NullString
		if err := baris.Scan(&kode, &nama); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelCedingCoList, err)
		}
		hasil = append(hasil, models.Ceding{ID: kode.String, Name: nama.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelCedingCoList, err)
	}
	return hasil, nil
}

// gabungCeding - gabungan `;` kode dan nama, urut daftar (E-8), PERSIS `@If(x=="", .CedingCo,
// x+";"+.CedingCo)` (SetCedingCo_Act / SetDataSobCeding_Act): pemisah dilewati selama hasil
// sementara masih kosong - nilai kosong di DEPAN tidak menghasilkan `;` di awal, di tengah
// tetap `;;`.
func gabungCeding(daftar []models.Ceding) (kode, nama string) {
	tambah := func(x, v string) string {
		if x == "" {
			return v
		}
		return x + pemisahCeding + v
	}
	for _, c := range daftar {
		kode, nama = tambah(kode, c.ID), tambah(nama, c.Name)
	}
	return kode, nama
}

// periksaLebarCeding - nama per baris dan gabungan muat di kolomnya.
func periksaLebarCeding(daftar []models.Ceding) error {
	for _, c := range daftar {
		if len(c.Name) > lebarNamaCeding {
			return fmt.Errorf("%w: nama ceding %s melebihi %d bita", ErrCedingTerlaluPanjang, c.ID, lebarNamaCeding)
		}
	}
	kode, nama := gabungCeding(daftar)
	if len(kode) > lebarGabunganKode || len(nama) > lebarGabunganNama {
		return fmt.Errorf("%w: gabungan kode %d/%d bita, nama %d/%d bita", ErrCedingTerlaluPanjang,
			len(kode), lebarGabunganKode, len(nama), lebarGabunganNama)
	}
	return nil
}

// uidAcak - ROW_UID baris buatan aplikasi: UUID versi 4 acak (A106; ROW_UID loader = v5
// deterministik, K-073: tidak boleh dijadikan sandaran di luar tabel flat).
func uidAcak() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

// simpanCeding - ganti seluruh daftar ceding case `id` (urut `kode`) di dalam transaksi
// pemanggil; T_QUOTATIONDATA baris case sudah ada (dipanggil sesudah upsert quotation).
func (r *KasusOracle) simpanCeding(ctx context.Context, tx *db.Tx, quo, id string, kode []string) error {
	ceding, err := r.db.Qualify(TabelCedingCoList)
	if err != nil {
		return err
	}
	daftar := make([]models.Ceding, len(kode))
	for i, k := range kode {
		nama, err := namaAgent(ctx, r.db, tx, k)
		if errors.Is(err, ErrSOBTidakSah) {
			return fmt.Errorf("%w: cedingIds[%d]", ErrCedingTidakSah, i)
		}
		if err != nil {
			return err
		}
		daftar[i] = models.Ceding{ID: k, Name: nama}
	}
	if err := periksaLebarCeding(daftar); err != nil {
		return err
	}
	var qid int64
	qID := sqlIDQuotation(quo)
	if err := db.PeriksaSQL(qID); err != nil {
		return err
	}
	if err := tx.QueryRowContext(ctx, qID, id).Scan(&qid); err != nil {
		return fmt.Errorf("repository: ID %s case: %w", TabelQuotationData, err)
	}
	if _, err := jalankan(ctx, tx, sqlHapusCeding(ceding), "menghapus "+TabelCedingCoList, qid); err != nil {
		return err
	}
	for i, c := range daftar {
		urut, err := r.db.NomorBerikut(ctx, tx, SequenceCedingCoList)
		if err != nil {
			return err
		}
		uid, err := uidAcak()
		if err != nil {
			return err
		}
		if err := r.sisip(ctx, tx, sqlSisipCeding(ceding), TabelCedingCoList, urut, qid, i+1, uid, c.ID, db.KosongJadiNil(c.Name)); err != nil {
			return err
		}
	}
	gk, gn := gabungCeding(daftar)
	_, err = jalankan(ctx, tx, sqlGabunganCeding(quo), "mengubah gabungan ceding", db.KosongJadiNil(gk), db.KosongJadiNil(gn), qid)
	return err
}
