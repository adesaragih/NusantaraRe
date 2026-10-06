package repository

// Baca `POOLDATA.M_ATTACHMENTTREATY_2` — panel Attachment modul Adjustment.
//
// ⛔ BACA SAJA. Nol `INSERT`/`UPDATE`/`DELETE`/`MERGE`, nol DDL, nol migrasi
// menyentuhnya. Tabel WARISAN yang hidup, 43 baris terukur 4 Oktober 2026.
//
// ⛔ NOL TABEL BARU dan nol migrasi di modul ini untuk lampiran.
//
// =====================================================================
// KENAPA KODENYA ADA DUA KALI, DI DUA MODUL
// =====================================================================
//   `Section/ShowAttachmentTreaty.xml`, `WorkAttachments.xml`, dan
//   `TreatyInActionButtons.xml` ada di KEDUA ekspor — Treaty In dan Treaty
//   In Adjustment — dan keduanya membaca tabel yang sama.
//
//   Yang TIDAK boleh dilakukan: satu modul mengimpor repository modul lain.
//   `TestModulTidakMengimporModulLain` menolaknya, dan alasannya berdiri
//   sendiri — dua modul yang berbagi kode repository berbagi juga jadwal
//   rilisnya, dan itu persis yang pemisahan modul ada untuk mencegahnya.
//
//   Jadi pembacaannya ditulis dua kali, dan kesamaannya DINYATAKAN di sini
//   alih-alih disembunyikan. Bila salah satu berubah, yang lain tidak ikut
//   berubah dengan sendirinya — dan itu memang yang dikehendaki: modul
//   Adjustment boleh menampilkan lampiran dengan cara yang berbeda.
//
//   ⚠️ Yang TIDAK boleh menyimpang adalah nama tabel dan nama kolomnya.
//   Keduanya milik sistem lama, bukan milik modul mana pun.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyinadjustment/backend/models"
)

// TabelWarisanLampiran - tabel lampiran sistem lama.
const TabelWarisanLampiran = "M_ATTACHMENTTREATY_2"

// BacaLampiranKontrak membaca seluruh lampiran satu kontrak.
//
// ⛔ `ORDER BY` berlapis, sama dengan modul Treaty In: kategori lebih dulu,
// lalu `DATEINPUT` MENURUN - yang terbaru di atas - lalu `ID` supaya dua
// unggahan pada detik yang sama punya urutan tetap.
func (g *Gudang) BacaLampiranKontrak(ctx context.Context, masterID string) ([]models.BarisLampiranWarisan, error) {
	nama, err := g.db.Qualify(TabelWarisanLampiran)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT ID, CATEGORY_ID, CATEGORY, FILENAME, FILEMIMETYPE,
		T_STORAGE_ID, TO_CHAR(DATEINPUT, 'YYYYMMDD'), USERNAME
		FROM %s WHERE TREATYID = :1
		ORDER BY CATEGORY_ID, DATEINPUT DESC, ID`, nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}

	rows, err := g.db.QueryContext(ctx, q, masterID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s kontrak %s: %w", TabelWarisanLampiran, masterID, err)
	}
	defer func() { _ = rows.Close() }()

	out := []models.BarisLampiranWarisan{}
	for rows.Next() {
		var id, kode, nk, berkas, mime, simpanan, tanggal, pengguna sql.NullString
		if err := rows.Scan(&id, &kode, &nk, &berkas, &mime, &simpanan, &tanggal, &pengguna); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s kontrak %s: %w", TabelWarisanLampiran, masterID, err)
		}
		out = append(out, models.BarisLampiranWarisan{
			ID:           id.String,
			KodeKategori: kode.String,
			NamaKategori: nk.String,
			NamaBerkas:   berkas.String,
			JenisMime:    mime.String,
			IDSimpanan:   simpanan.String,
			Diunggah:     tanggal.String,
			Pengunggah:   pengguna.String,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca %s kontrak %s: %w", TabelWarisanLampiran, masterID, err)
	}
	return out, nil
}

// BacaKatalogKategoriLampiran membaca pasangan kode↔nama DARI DATA.
//
// ⛔ Dibaca, tidak dihafal. `M_ATTACHMENTTREATY_2` menyimpan `CATEGORY_ID`
// dan `CATEGORY` berdampingan, jadi ketujuh pasangan yang terpakai adalah
// fakta yang diambil - bukan daftar di dalam kode yang akan membeku pada
// hari seseorang mengganti sebuah nama di sistem lama.
func (g *Gudang) BacaKatalogKategoriLampiran(ctx context.Context) (map[string]string, error) {
	nama, err := g.db.Qualify(TabelWarisanLampiran)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT CATEGORY_ID, MAX(CATEGORY) FROM %s
		WHERE CATEGORY_ID IS NOT NULL GROUP BY CATEGORY_ID ORDER BY CATEGORY_ID`, nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca katalog %s: %w", TabelWarisanLampiran, err)
	}
	defer func() { _ = rows.Close() }()

	out := map[string]string{}
	for rows.Next() {
		var kode, nk sql.NullString
		if err := rows.Scan(&kode, &nk); err != nil {
			return nil, fmt.Errorf("repository: membaca katalog %s: %w", TabelWarisanLampiran, err)
		}
		out[kode.String] = nk.String
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca katalog %s: %w", TabelWarisanLampiran, err)
	}
	return out, nil
}
