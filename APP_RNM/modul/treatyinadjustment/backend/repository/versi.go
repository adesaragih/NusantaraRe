// Package repository memegang seluruh sentuhan basis data modul Treaty In
// Adjustment.
//
// ⛔ Tabel yang dibacanya - `KONTRAK` dan `VERSI_KONTRAK` - DIBUAT modul
// `treatyin` (migrasi 401). Modul ini menambahkan satu kolom padanya (migrasi
// 440) dan membacanya; ia TIDAK mengimpor paket modul itu. Yang dibagi adalah
// tabelnya, bukan kodenya.
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyinadjustment/backend/models"
)

// Gudang adalah seluruh sentuhan basis data modul ini.
type Gudang struct{ db *db.DB }

// Baru membuat Gudang di atas koneksi.
func Baru(d *db.DB) *Gudang { return &Gudang{db: d} }

// polaTanggal - DATE Oracle dibaca sebagai teks `YYYY-MM-DD`, tidak pernah
// bergantung NLS sesi.
const polaTanggal = "YYYY-MM-DD"

// DaftarKontrak membaca kepala kontrak, berurut menurut pengenalnya.
func (g *Gudang) DaftarKontrak(ctx context.Context) ([]models.Kontrak, error) {
	nama, err := g.db.Qualify("KONTRAK")
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT ID_KONTRAK, NOMOR_KONTRAK_WARISAN, SIFAT_PROPORSI,
		TO_CHAR(TANGGAL_MULAI, '%s'), TO_CHAR(TANGGAL_BERAKHIR, '%s')
		FROM %s ORDER BY ID_KONTRAK ASC`, polaTanggal, polaTanggal, nama)
	baris, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca KONTRAK: %w", err)
	}
	defer func() { _ = baris.Close() }()

	keluar := []models.Kontrak{}
	for baris.Next() {
		var k models.Kontrak
		var warisan sql.NullString
		if err := baris.Scan(&k.ID, &warisan, &k.SifatProporsi, &k.TanggalMulai, &k.TanggalBerakhir); err != nil {
			return nil, fmt.Errorf("repository: membaca baris KONTRAK: %w", err)
		}
		k.NomorKontrakWarisan = warisan.String
		keluar = append(keluar, k)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca KONTRAK: %w", err)
	}
	return keluar, nil
}

// RantaiVersi membaca seluruh versi satu kontrak.
//
// Urutannya NOMOR_URUT_VERSI, lalu pengenal - dan `NULLS LAST` disebut
// eksplisit: tiket 05 membuat nomor urut boleh kosong, dan baris warisan yang
// belum dinomori ulang (tiket 10) harus mendarat di satu ujung yang tetap,
// bukan di tempat yang bergantung bawaan Oracle.
func (g *Gudang) RantaiVersi(ctx context.Context, idKontrak int64) ([]models.Versi, error) {
	nama, err := g.db.Qualify("VERSI_KONTRAK")
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT ID_VERSI_KONTRAK, ID_KONTRAK, NOMOR_URUT_VERSI, KEADAAN_SIKLUS_HIDUP,
		JENIS_ADDENDUM, SIFAT_MATERIAL_ADDENDUM, TO_CHAR(TANGGAL_BERLAKU_ADDENDUM, '%s'),
		ID_VERSI_KONTRAK_DASAR, NAMA_KONTRAK
		FROM %s WHERE ID_KONTRAK = :1
		ORDER BY NOMOR_URUT_VERSI ASC NULLS LAST, ID_VERSI_KONTRAK ASC`, polaTanggal, nama)
	baris, err := g.db.QueryContext(ctx, q, idKontrak)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca VERSI_KONTRAK: %w", err)
	}
	defer func() { _ = baris.Close() }()

	keluar := []models.Versi{}
	for baris.Next() {
		var v models.Versi
		var nomor, dasar sql.NullInt64
		var jenis, material, berlaku sql.NullString
		if err := baris.Scan(&v.ID, &v.IDKontrak, &nomor, &v.KeadaanSiklusHidup,
			&jenis, &material, &berlaku, &dasar, &v.NamaKontrak); err != nil {
			return nil, fmt.Errorf("repository: membaca baris VERSI_KONTRAK: %w", err)
		}
		if nomor.Valid {
			n := nomor.Int64
			v.NomorUrutVersi = &n
		}
		if dasar.Valid {
			d := dasar.Int64
			v.IDVersiDasar = &d
		}
		v.JenisAddendum, v.SifatMaterialAddendum, v.TanggalBerlakuAddendum = jenis.String, material.String, berlaku.String
		keluar = append(keluar, v)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca VERSI_KONTRAK: %w", err)
	}
	return keluar, nil
}
