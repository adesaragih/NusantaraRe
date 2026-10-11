package repository

// Baca `POOLDATA.T_VIEW_COMMENT` — panel History modul Adjustment.
//
// ⛔ BACA SAJA. Nol `INSERT`/`UPDATE`/`DELETE`/`MERGE`, nol DDL, nol migrasi
// modul ini menyentuhnya.
//
// =====================================================================
// KENAPA MODUL INI MEMBACANYA SENDIRI
// =====================================================================
//   Ronde sebelumnya meninggalkan panel History kosong dengan alasan
//   *"sumbernya milik modul sebelah, dan impor silang ditolak
//   `TestModulTidakMengimporModulLain`"*. Alasan itu SETENGAH benar, dan
//   setengah yang keliru membuat panelnya kosong tanpa perlu.
//
//   Yang penjaga itu tolak adalah **impor paket Go lintas modul** — bukan
//   pembacaan tabel. `T_VIEW_COMMENT` berdiri di POOLDATA, skema yang
//   sama dengan `M_ATTACHMENTTREATY_2` yang modul ini sudah baca sendiri.
//   Jadi pembacanya ditulis di sini, dan NOL baris mengimpor `treatyin`.
//
//   ⚠️ Tabelnya dibuat migrasi `430` modul `treatyin` dan dimuat pemuat
//   modul itu — 11.365 baris. Modul ini bukan pemiliknya dan tidak pernah
//   menulisnya; bila `treatyin` kelak mengubah bentuknya, pembacaan di
//   sini ikut berubah, dan kuerinya menyebut kolomnya satu per satu justru
//   supaya perubahan itu berbunyi sebagai galat, bukan sebagai kolom yang
//   diam-diam kosong.
//
// Keempat kolomnya berpadanan SATU-KE-SATU dengan kolom layar lama:
//
//	TANGGAL       -> Date
//	OPERATORNAME  -> PIC
//	ISAPPROVED    -> Approval
//	SUGGEST       -> Comment

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyinadjustment/backend/models"
)

// TabelWarisanRiwayat - tabel komentar/persetujuan kontrak.
const TabelWarisanRiwayat = "T_VIEW_COMMENT"

// BacaRiwayatKontrak membaca seluruh baris riwayat satu kontrak.
//
// ⛔ `ORDER BY URUTAN`. Larik komentar di sistem lama BERURUT, dan `URUTAN`
// adalah posisi aslinya di dalam dokumen. Tanpa klausa ini Oracle bebas
// mengembalikan baris dalam urutan apa pun, dan riwayat persetujuan yang
// urutannya tertukar TERBACA BENAR.
func (g *Gudang) BacaRiwayatKontrak(ctx context.Context, masterID string) ([]models.BarisRiwayatWarisan, error) {
	nama, err := g.db.Qualify(TabelWarisanRiwayat)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT TANGGAL, OPERATORNAME, ISAPPROVED, SUGGEST
		FROM %s WHERE MASTERID = :1 ORDER BY URUTAN`, nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}

	rows, err := g.db.QueryContext(ctx, q, masterID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s kontrak %s: %w", TabelWarisanRiwayat, masterID, err)
	}
	defer func() { _ = rows.Close() }()

	out := []models.BarisRiwayatWarisan{}
	for rows.Next() {
		var tgl, operator, disetujui, catatan sql.NullString
		if err := rows.Scan(&tgl, &operator, &disetujui, &catatan); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s kontrak %s: %w", TabelWarisanRiwayat, masterID, err)
		}
		out = append(out, models.BarisRiwayatWarisan{
			// ⛔ Apa adanya. Nilainya TEKS di tabel pendaratan — nol tafsir
			// di repository; yang menerjemahkan untuk layar `services`.
			Tanggal: tgl.String, Operator: operator.String,
			Disetujui: disetujui.String, Catatan: catatan.String,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca %s kontrak %s: %w", TabelWarisanRiwayat, masterID, err)
	}
	return out, nil
}
