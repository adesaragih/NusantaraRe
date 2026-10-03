// Package repository adalah SATU-SATUNYA lapisan modul Master Product Name
// Life yang berbicara ke Oracle, dan satu-satunya yang mengenal bentuk JSON
// Pega. Seluruh SQL modul ini ada di sini.
//
// Untuk apa berkas ini: nama objek Oracle kedua tabel JSON warisan dan Gudang - ditulis SEKALI.
//
// ⭐ Sejak 02-10-2026 (K5, OQ-MPNL-01 flat; tiket 01 bab bertanggal): produk disimpan di tabel FLAT
// `M_PRODUCTNAME_LIFE` + tujuh anak (`mpnl_flat.go`). Kedua tabel warisan `M_PRODUCT_LIFE` dan
// `M_PRODUCTINWARD_LIFE` (P1 lama: JSON seperti Pega) TIDAK ditulis dan TIDAK dibaca aplikasi lagi - cadangan dan
// sumber alat pindah saja (`mpnl_pindah.go`). Ketiga view DEV tidak dipakai dan tidak dibangun ulang (K7). Dijaga
// `TestMPNLAplikasiHanyaTabelFlat` (keputusan work owner 02-10-2026).
// ⛔ Setiap query menyebut skemanya lewat `db.Qualify` dan diperiksa
// `db.PeriksaSQL` (nol COMMIT). Prosedur `PEGA_M_PRODUCT_LIFE` dan
// `PEGA_M_PRODUCT_INWARD_LIFE` TIDAK dipanggil (logikanya ditiru).
//
// Dibaca sesudah: models/mpnl_produk.go.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// Dua tabel warisan produk (`dba-procedures-and-ddl.md` DDL) - sejak 02-10-2026 hanya DIBACA alat pindah.
const (
	TabelProduk = "M_PRODUCT_LIFE"
	TabelInward = "M_PRODUCTINWARD_LIFE"
)

// Kolom fisik kedua tabel - katalog DEV `ALL_TAB_COLUMNS` 01-10-2026
// (`testdata/katalog-dev.json`, uji `TestKolomDibacaAdaDiKatalogDEV`). `M_PRODUCT_LIFE`:
// `JSONDATA` CLOB (`IS JSON`) + DUA kolom datar `RIRISKID`, `RIRISK` - penulisnya di
// korpus hanya `SaveProductNameLIfeFlat` b84. ⛔ `PRODUCTNAME` dan `BEGIN_DATE` TIDAK ADA
// di DEV (ralat lanjutan 1 L1; `dba-procedures-and-ddl.md` keliru menyebut empat).
// `M_PRODUCTINWARD_LIFE`: `ID` + `JSONDATA`.
var (
	KolomProduk = []string{"ID", "JSONDATA", "RIRISKID", "RIRISK"}
	KolomInward = []string{"ID", "JSONDATA"}
)

// kuerier - baca lewat transaksi bila ada (anti-basi di dalam penyimpanan),
// selain itu lewat koneksi.
type kuerier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Gudang membaca dan menulis tabel flat produk, membaca master, dan membaca kedua tabel JSON warisan (alat pindah).
type Gudang struct{ db *db.DB }

// Baru menyusun Gudang di atas koneksi bersama.
func Baru(d *db.DB) *Gudang { return &Gudang{db: d} }

func (g *Gudang) kueri(tx *db.Tx) kuerier {
	if tx.Terisi() {
		return tx
	}
	return g.db
}

// siapkan mengualifikasi nama objek dan memeriksa teks SQL-nya.
func (g *Gudang) siapkan(objek string, susun func(tabel string) string) (string, error) {
	tabel, err := g.db.Qualify(objek)
	if err != nil {
		return "", err
	}
	q := susun(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	if err := periksaBacaSaja(objek, q); err != nil {
		return "", err
	}
	return q, nil
}

// ErrMasterBacaSaja - SQL yang bukan SELECT diarahkan ke objek master/view yang dibaca saja.
var ErrMasterBacaSaja = errors.New("repository: reference master is read-only")

// periksaBacaSaja - lapis kedua penjaga modul `TestMPNLMasterDibacaSaja`: objek di
// DaftarMasterDibacaSaja (termasuk kedua view rate, K1 01-10-2026) hanya menerima SELECT, sebelum
// teks SQL apa pun sampai ke Oracle.
func periksaBacaSaja(objek, q string) error {
	for _, m := range DaftarMasterDibacaSaja {
		if m == objek && !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(q)), "SELECT ") {
			return fmt.Errorf("%w: %s", ErrMasterBacaSaja, objek)
		}
	}
	return nil
}
