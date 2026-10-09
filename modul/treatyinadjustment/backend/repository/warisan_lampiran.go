package repository

// Baca `POOLDATA.M_ATTACHMENTTREATY_2` — panel Attachment modul Adjustment.
//
// ⛔ BACA SAJA. Nol `INSERT`/`UPDATE`/`DELETE`/`MERGE`, nol DDL, nol migrasi
// menyentuhnya. Tabel WARISAN yang hidup, 43 baris terukur 4 Oktober 2026.
//
// ⛔ NOL TABEL BARU dan nol migrasi di modul ini untuk lampiran.
//
// ⭐ 8 Oktober 2026: nama kategori dari `POOLDATA.M_KATEGORIMASTERTREATY`
// (katalog RD Pega) — juga BACA SAJA.
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
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyinadjustment/backend/models"
)

// TabelWarisanLampiran - tabel lampiran sistem lama.
const TabelWarisanLampiran = "M_ATTACHMENTTREATY_2"

// TabelKategoriLampiran - katalog kategori yang RD Pega baca
// (`GetMasterTreatyCategory_SQL`: `SELECT id, note FROM
// POOLDATA.M_KATEGORIMASTERTREATY order by note`). BACA SAJA.
const TabelKategoriLampiran = "M_KATEGORIMASTERTREATY"

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

// BacaKatalogKategoriLampiran membaca pasangan kode↔nama.
//
// ⭐ 8 Oktober 2026 — sumber UTAMANYA katalog yang RD Pega baca:
// `M_KATEGORIMASTERTREATY` (`GetMasterTreatyCategory_SQL`), sama dengan
// modul Treaty In. Kesebelas pasangannya ada di sana, termasuk keempat kode
// yang dulu "tanpa nama" (`00003` Binding, signed share Email · `00004` Info
// Pack · `00008` LOA · `00009` Claim Data) —
// `treatyin/docs/PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md` TERJAWAB.
// Pasangan dari DATA lampiran melengkapi kode yang tidak ada di katalog.
func (g *Gudang) BacaKatalogKategoriLampiran(ctx context.Context) (map[string]string, error) {
	out, err := g.BacaKategoriMaster(ctx)
	if err != nil {
		return nil, err
	}
	data, err := g.bacaKatalogDariData(ctx)
	if err != nil {
		return nil, err
	}
	for kode, nama := range data {
		if _, ada := out[kode]; !ada {
			out[kode] = nama
		}
	}
	return out, nil
}

// BacaKategoriMaster - kesebelas kategori `M_KATEGORIMASTERTREATY`, kode → nama.
func (g *Gudang) BacaKategoriMaster(ctx context.Context) (map[string]string, error) {
	nama, err := g.db.Qualify(TabelKategoriLampiran)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf("SELECT ID, NOTE FROM %s ORDER BY NOTE", nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s: %w", TabelKategoriLampiran, err)
	}
	defer func() { _ = rows.Close() }()
	out := map[string]string{}
	for rows.Next() {
		var kode, nk sql.NullString
		if err := rows.Scan(&kode, &nk); err != nil {
			return nil, fmt.Errorf("repository: membaca %s: %w", TabelKategoriLampiran, err)
		}
		// ⚠️ `TrimSpace`: NOTE `00004` tersimpan dengan ekor CR LF (terukur
		// 8 Oktober 2026) — ekor baris baru bukan bagian namanya.
		if k := strings.TrimSpace(kode.String); k != "" {
			out[k] = strings.TrimSpace(nk.String)
		}
	}
	return out, rows.Err()
}

// bacaKatalogDariData - pasangan `CATEGORY_ID`/`CATEGORY` yang terpakai di
// `M_ATTACHMENTTREATY_2`.
//
// ⛔ Dibaca, tidak dihafal: tabel menyimpan `CATEGORY_ID` dan `CATEGORY`
// berdampingan, jadi pasangan yang terpakai adalah fakta yang diambil.
func (g *Gudang) bacaKatalogDariData(ctx context.Context) (map[string]string, error) {
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
