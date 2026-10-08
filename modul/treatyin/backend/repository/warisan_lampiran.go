package repository

// Baca `POOLDATA.M_ATTACHMENTTREATY_2` — panel Attachment.
//
// ⭐ Sejak 8 Oktober 2026 tombol `Upload file` MENULIS ke tabel ini
// (keputusan pemakai) — jalurnya di `lampiran_tulis.go`. Berkas ini tetap
// hanya MEMBACA. Nol DDL, nol penyebutan di migrasi mana pun: tabel ini
// WARISAN yang hidup — 43 baris terukur 4 Oktober 2026.
//
// ⛔ NOL TABEL BARU. Aturan "tabel baru hanya untuk struktur tab ber-JSON"
// tidak berlaku di sini: lampirannya sudah relasional, dengan `CATEGORY_ID`
// dan `CATEGORY` berdampingan dan `T_STORAGE_ID` menunjuk `T_STORAGE_IMAGE`.
//
// Kueri di bawah mengikuti `GetAttachment2_Sql` dari ekspor:
//
//	select a.ID, a.filename as "pyFileName", a.CATEGORY_ID as "pyCategory",
//	       a.FILEMIMETYPE as "pyFileMimeType", a.T_STORAGE_ID as "type"
//	  from M_ATTACHMENTTREATY_2 a
//	 where treatyid = {TreatyIn.ID} AND CATEGORY_ID = {StatusDoc.CARI40}
//
// ⚠️ Dengan DUA penyimpangan yang dinyatakan:
//
//  1. saringan `CATEGORY_ID` DIBUANG. Panel menampilkan SELURUH kategori
//     beserta cacahnya sekaligus; menanyakannya satu per satu berarti
//     sebelas perjalanan ke basis data untuk satu panel.
//  2. `DATEINPUT` dan `USERNAME` ikut dibaca. Keduanya ada di tabel dan
//     tidak ada di kueri lama — dan "siapa mengunggah, kapan" adalah hal
//     pertama yang ditanyakan tentang sebuah lampiran.
//
// ⛔ `DATA_JSON` TIDAK dibaca. Ia `CLOB` berisi muatan berkasnya, dan panel
// hanya menampilkan daftar; menariknya pada tiap pembukaan form berarti
// memindahkan berkas yang tidak ada yang minta.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
)

// TabelWarisanLampiran - tabel lampiran sistem lama.
const TabelWarisanLampiran = "M_ATTACHMENTTREATY_2"

// BacaLampiranKontrak membaca seluruh lampiran satu kontrak.
//
// ⛔ `ORDER BY` berlapis: kategori lebih dulu supaya panel dapat
// mengelompokkan tanpa mengurutkan ulang, lalu `DATEINPUT` MENURUN —
// lampiran terbaru di atas, sebab itu yang dicari orang — lalu `ID` supaya
// dua unggahan pada detik yang sama punya urutan tetap.
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
		var id, kode, nama, berkas, mime, simpanan, tanggal, pengguna sql.NullString
		if err := rows.Scan(&id, &kode, &nama, &berkas, &mime, &simpanan, &tanggal, &pengguna); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s kontrak %s: %w", TabelWarisanLampiran, masterID, err)
		}
		out = append(out, models.BarisLampiranWarisan{
			ID:           id.String,
			KodeKategori: kode.String,
			NamaKategori: nama.String,
			NamaBerkas:   berkas.String,
			JenisMime:    mime.String,
			IDSimpanan:   simpanan.String,
			// ⛔ `YYYYMMDD` apa adanya — bentuk yang `services.TanggalTampil`
			// sudah tahu menerjemahkan. Nol penerjemah kelima.
			Diunggah:   tanggal.String,
			Pengunggah: pengguna.String,
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
// `M_KATEGORIMASTERTREATY` (`GetMasterTreatyCategory_SQL`). Kesebelas
// pasangannya ada di sana, termasuk keempat kode yang dulu "tanpa nama"
// (`00003` Binding, signed share Email · `00004` Info Pack · `00008` LOA ·
// `00009` Claim Data) — `docs/PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md`
// TERJAWAB. Pasangan dari DATA lampiran melengkapi kode yang tidak ada di
// katalog (nol hari ini).
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

// bacaKatalogDariData - pasangan `CATEGORY_ID`/`CATEGORY` yang terpakai di
// `M_ATTACHMENTTREATY_2`.
func (g *Gudang) bacaKatalogDariData(ctx context.Context) (map[string]string, error) {
	nama, err := g.db.Qualify(TabelWarisanLampiran)
	if err != nil {
		return nil, err
	}
	// ⛔ `MAX(CATEGORY)` bukan sekadar menghindari GROUP BY: bila satu kode
	// ternyata punya dua nama di tabel, yang dipilih tetap, bukan berubah
	// tiap pembacaan. Ketidakcocokannya ditagih `TestSatuKodeSatuNama`.
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
