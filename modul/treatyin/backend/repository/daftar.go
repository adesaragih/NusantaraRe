package repository

// Jalur daftar kontrak - ronde layar 1.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/modul/treatyin/backend/models"
)

// DaftarKontrak membaca kepala seluruh kontrak beserta nama dan keadaan versi
// TERAKHIRNYA.
//
// ⛔ Versi terakhir diambil lewat sub-SELECT berurut `NOMOR_URUT_VERSI DESC
// NULLS LAST, ID_VERSI_KONTRAK DESC`, bukan lewat `MAX(ID)`. Sebabnya tiket 05
// membuat `NOMOR_URUT_VERSI` boleh kosong: baris warisan yang belum dinomori
// ulang (tiket 10) harus mendarat di satu ujung yang TETAP, bukan di tempat
// yang bergantung bawaan Oracle.
//
// ⚠️ `LEFT JOIN`, bukan `JOIN`. Kontrak tanpa versi tidak dapat lahir lewat
// jalur aplikasi (tiket 14 membuat keduanya dalam satu transaksi), tetapi
// baris hasil pemindahan belum tentu demikian - dan kontrak yang hilang dari
// layar karena versinya tidak ada adalah kehilangan yang tidak terlihat.
func (g *Gudang) DaftarKontrak(ctx context.Context) ([]models.BarisDaftarKontrak, error) {
	tKontrak, err := g.db.Qualify("KONTRAK")
	if err != nil {
		return nil, err
	}
	tVersi, err := g.db.Qualify("VERSI_KONTRAK")
	if err != nil {
		return nil, err
	}
	// ⛔ Kata kerja `%s` biasa, BUKAN berindeks (`%[1]s`). Penjaga
	// `TestNamaTabelTidakTelanjang` (ADR-U-0033) mengenali `%s` sebagai nama
	// yang datang dari `Qualify`, dan bentuk berindeks lolos darinya sebagai
	// "nama tabel telanjang". Pengulangan argumennya lebih murah daripada
	// penjaga yang berhenti menjaga.
	q := fmt.Sprintf(`SELECT k.ID_KONTRAK, k.ID_CEDANT, k.ID_ASAL_BISNIS, k.SIFAT_PROPORSI,
		TO_CHAR(k.TANGGAL_MULAI,'YYYY-MM-DD'), TO_CHAR(k.TANGGAL_BERAKHIR,'YYYY-MM-DD'),
		(SELECT v.NAMA_KONTRAK FROM %s v WHERE v.ID_KONTRAK = k.ID_KONTRAK
		   ORDER BY v.NOMOR_URUT_VERSI DESC NULLS LAST, v.ID_VERSI_KONTRAK DESC
		   FETCH FIRST 1 ROWS ONLY),
		(SELECT v.KEADAAN_SIKLUS_HIDUP FROM %s v WHERE v.ID_KONTRAK = k.ID_KONTRAK
		   ORDER BY v.NOMOR_URUT_VERSI DESC NULLS LAST, v.ID_VERSI_KONTRAK DESC
		   FETCH FIRST 1 ROWS ONLY)
		FROM %s k ORDER BY k.ID_KONTRAK ASC`, tVersi, tVersi, tKontrak)

	baris, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca daftar KONTRAK: %w", err)
	}
	defer func() { _ = baris.Close() }()

	keluar := []models.BarisDaftarKontrak{}
	for baris.Next() {
		var b models.BarisDaftarKontrak
		var nama, keadaan sql.NullString
		if err := baris.Scan(&b.ID, &b.IDCedant, &b.IDAsalBisnis, &b.SifatProporsi,
			&b.TanggalMulai, &b.TanggalBerakhir, &nama, &keadaan); err != nil {
			return nil, fmt.Errorf("repository: membaca baris daftar KONTRAK: %w", err)
		}
		b.NamaKontrak = nama.String
		b.KeadaanSiklusHidup = keadaan.String
		// PosisiKe sengaja dibiarkan kosong - nol kolom menyimpannya.
		keluar = append(keluar, b)
	}
	return keluar, baris.Err()
}
