package repository

// Sumber daftar kategori dokumen wajib - butir ar1, A2.
//
// `[data DBA]` `POOLDATA.CATEGORY_ATTACH_CLAIMLIFE` - `NOTE VARCHAR2(30)`,
// `BISNIS VARCHAR2(100)`, `POSITION VARCHAR2(100)`, `NOU VARCHAR2(2)`;
// enam baris, seluruhnya `BISNIS = 'ALL'` dengan `POSITION` dan `NOU` NULL.
//
// ⚠️ `ar1` `[USULAN yang disahkan — 27 September 2026]`, BUKAN
// `[terverifikasi]`. Rule `GetCategoryLife_SQL` **tidak ada di korpus** -
// seluruh 29 berkas `Claim Life/RDBList/` sudah dicacah - sehingga tabel mana
// yang Pega benar-benar baca tidak dapat dipastikan. Kandidat kedua yang
// dicatat dan TIDAK dipakai: view `DOCUMENTCLAIM_LIFE` atas
// `M_PRODUCT_LIFE.JSONDATA`, yang jalur JSON produknya dilarang AC 38.
//
// ⛔ Isinya LABEL KONFIGURASI, bukan data orang - boleh masuk tiruan skema
// uji.
//
// Dibaca sesudah: linkservice.go.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// sqlKategoriWajib merakit pernyataannya.
//
// ⚠️ `BISNIS IN ('ALL', <kode>)`: baris `ALL` berlaku bagi seluruh produk,
// dan baris berkode hanya bagi produknya. Keenam baris DEV seluruhnya `ALL`,
// tetapi menyaringnya tetap benar - kolomnya ada justru untuk itu.
func sqlKategoriWajib(tabel string) string {
	return fmt.Sprintf(
		`SELECT NOTE FROM %s WHERE BISNIS IN ('ALL', :1) ORDER BY NOTE`, tabel)
}

// AmbilKategoriWajib membaca daftar kategori dokumen yang wajib ada.
func (r *PohonKlaim) AmbilKategoriWajib(ctx context.Context,
	kodeBisnis string) ([]string, error) {

	tabel, err := r.db.Qualify("CATEGORY_ATTACH_CLAIMLIFE")
	if err != nil {
		return nil, err
	}
	q := sqlKategoriWajib(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, q, kodeBisnis)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca kategori dokumen wajib: %w", err)
	}
	defer func() { _ = baris.Close() }()

	var hasil []string
	for baris.Next() {
		var nota sql.NullString
		if err := baris.Scan(&nota); err != nil {
			return nil, fmt.Errorf("repository: membaca baris kategori: %w", err)
		}
		if nota.Valid && nota.String != "" {
			hasil = append(hasil, nota.String)
		}
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: menelusuri kategori: %w", err)
	}
	return hasil, nil
}

// TabelDokumenKlaimProduk - daftar "DOCUMENT CLAIM / Document List" Master Product Name Life (modul masterproductnamelife,
// migrasi 142: PRODUCTID VARCHAR2(6), URUT NUMBER(5), DOCUMENT VARCHAR2(500); PK (PRODUCTID, URUT)). BACA SAJA; salinan
// bentuk baca, tanpa impor modul master.
const TabelDokumenKlaimProduk = "M_PRODUCTNAME_LIFE_DOCCLAIM"

// sqlDokumenKlaimProduk - Document Claim satu produk, URUTAN sama dengan master (URUT).
func sqlDokumenKlaimProduk(tabel string) string {
	return fmt.Sprintf(`SELECT d.DOCUMENT FROM %s d WHERE d.PRODUCTID = :1 ORDER BY d.URUT`, tabel)
}

// DokumenKlaimProduk membaca daftar Document Claim produk (keputusan work owner 11-10-2026: "category document didapat
// dari document claim yang ada di product name yang digunakan oleh claim"). Produk tanpa daftar = kosong.
func (r *PohonKlaim) DokumenKlaimProduk(ctx context.Context, produkID string) ([]string, error) {
	tabel, err := r.db.Qualify(TabelDokumenKlaimProduk)
	if err != nil {
		return nil, err
	}
	q := sqlDokumenKlaimProduk(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, q, produkID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca Document Claim produk %s: %w", produkID, err)
	}
	defer func() { _ = baris.Close() }()
	var hasil []string
	for baris.Next() {
		var d sql.NullString
		if err := baris.Scan(&d); err != nil {
			return nil, fmt.Errorf("repository: memindai Document Claim: %w", err)
		}
		if t := strings.TrimSpace(d.String); t != "" {
			hasil = append(hasil, t)
		}
	}
	return hasil, baris.Err()
}
