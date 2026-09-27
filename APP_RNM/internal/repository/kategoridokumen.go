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
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	baris, err := r.db.sql.QueryContext(ctx, q, kodeBisnis)
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
