package repository

// Jalur baca isi pemilih "Choose Ceding" dan "Choose Source of Business".
//
// ⛔ BACA SAJA atas `TREATY_IN`, sama seperti `warisan_daftar.go`, dan
// tunduk pada larangan yang sama: nol `INSERT`, nol `UPDATE`, nol `DELETE`,
// nol DDL. Dijaga `TestWarisanHanyaDibaca`.
//
// ⛔ Keduanya memakai SATU fungsi, bukan dua yang disalin. Perbedaannya hanya
// sepasang nama kolom, dan dua salinan berarti dua tempat untuk lupa
// menyaring `NULL`.
//
// ⭐ KAPAN BERKAS INI BOLEH DICABUT: ketika tabel acuan cedant dan asal
// bisnis berdiri dan terisi. Sampai itu, menolak menampilkan pemilih berarti
// menyembunyikan 94 nilai yang sebenarnya ada.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/modul/treatyin/backend/models"
)

// BacaDaftarCedant mengembalikan seluruh pasangan pengenal+nama cedant.
func (g *Gudang) BacaDaftarCedant(ctx context.Context) ([]models.PilihanWarisan, error) {
	return g.bacaPilihan(ctx, "CEDINGID", "CEDING")
}

// BacaDaftarAsalBisnis mengembalikan seluruh pasangan pengenal+nama sumber.
func (g *Gudang) BacaDaftarAsalBisnis(ctx context.Context) ([]models.PilihanWarisan, error) {
	return g.bacaPilihan(ctx, "LEADINGREINSSOURCEID", "LEADINGREINSSOURCE")
}

// bacaPilihan menarik pasangan UNIK pengenal+nama, urut nama lalu pengenal.
//
// ⛔ `DISTINCT` atas KEDUA kolom, bukan atas namanya saja. Menyatukan per
// nama akan membuang pengenal kedua tanpa jejak — persis yang `PilihanWarisan`
// larang.
//
// ⚠️ Barisnya disaring `nama IS NOT NULL`: satu dari 1.854 kontrak tidak
// punya cedant maupun asal bisnis, dan baris tanpa nama di dalam pemilih
// hanya dapat dipilih secara tidak sengaja.
func (g *Gudang) bacaPilihan(ctx context.Context, kolomID, kolomNama string) ([]models.PilihanWarisan, error) {
	nama, err := g.db.Qualify(TabelWarisanKontrak)
	if err != nil {
		return nil, err
	}
	// ⛔ Nama kolom DITANAM di kode, tidak pernah datang dari permintaan —
	// keduanya konstanta di berkas ini. Tidak ada jalur dari masukan pemakai
	// ke dalam teks kueri.
	q := fmt.Sprintf(`SELECT DISTINCT %s, %s FROM %s
		WHERE %s IS NOT NULL
		ORDER BY %s, %s`, kolomID, kolomNama, nama, kolomNama, kolomNama, kolomID)

	baris, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca pilihan %s dari %s: %w", kolomNama, TabelWarisanKontrak, err)
	}
	defer func() { _ = baris.Close() }()

	keluar := []models.PilihanWarisan{}
	for baris.Next() {
		// ⚠️ Pengenalnya NULLABLE walau namanya tidak — dan itu bukan
		// kemungkinan teoretis: pengenal kosong di sini berarti kontrak lama
		// yang namanya tercatat tanpa kodenya.
		var id, nm sql.NullString
		if err := baris.Scan(&id, &nm); err != nil {
			return nil, fmt.Errorf("repository: membaca baris pilihan %s: %w", kolomNama, err)
		}
		keluar = append(keluar, models.PilihanWarisan{ID: id.String, Nama: nm.String})
	}
	return keluar, baris.Err()
}
