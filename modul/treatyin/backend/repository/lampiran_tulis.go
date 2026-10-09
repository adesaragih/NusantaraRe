package repository

// Jalur TULIS panel Attachment — tombol `Upload file` (`Section/
// WorkAttachments.xml`) → FlowAction `TreatyAttachContent` → Activity
// `TreatySaveAttachment`.
//
// ⭐ Keputusan pemakai 8 Oktober 2026: *"untuk upload file seharusnya kesini
// SELECT * FROM M_ATTACHMENTTREATY_2"* — tabel warisan yang sama yang panel
// baca, bukan tabel baru.
//
//	TreatySaveAttachment [1.4]  InsertGoogleStorage_Act → `Insert_T_Storage_SQL`
//	                     [1.5]  `Datain.CARI50 = ""` — kolom JSON tabel
//	                            lampiran diisi KOSONG; di sini tidak disebut
//	                            sama sekali (NULL)
//	                     [1.6]  `InsertAttachment2_Sql` — hanya bila objeknya
//	                            terkirim (`Datain.CARI51 != ""`)
//
// ⚠️ Pega menjalankan kedua INSERT di dua langkah terpisah; di sini satu
// transaksi — baris lampiran tanpa objek, atau objek tanpa baris, tidak lahir.
// ⚠️ `DATEINPUT` dan `TANGGAL_UPLOAD` diisi `DEFAULT SYSDATE` tabelnya,
// persis seperti SQL Pega yang tidak menyebutnya.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
)

const (
	// TabelKategoriLampiran - katalog kategori yang RD Pega baca
	// (`GetMasterTreatyCategory_SQL`: `SELECT id, note FROM
	// POOLDATA.M_KATEGORIMASTERTREATY order by note`).
	TabelKategoriLampiran = "M_KATEGORIMASTERTREATY"
	// TabelObjekSimpanan - catatan objek di Google Storage.
	TabelObjekSimpanan = "T_STORAGE_IMAGE"
	// TabelFolderSimpanan - `GetAppName_SQL`: `SELECT APPNAME FROM T_FOLDER_IMAGE`.
	TabelFolderSimpanan = "T_FOLDER_IMAGE"
)

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

// NamaAplikasiSimpanan - `GetAppName_SQL` (`InsertGoogleStorage_Act` [6],
// prasyarat `IsPEGAPROD` nonaktif → selalu dijalankan): baris pertama.
func (g *Gudang) NamaAplikasiSimpanan(ctx context.Context) (string, error) {
	nama, err := g.db.Qualify(TabelFolderSimpanan)
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf("SELECT APPNAME FROM %s FETCH FIRST 1 ROWS ONLY", nama)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var app sql.NullString
	if err := g.db.QueryRowContext(ctx, q).Scan(&app); err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("repository: membaca %s: %w", TabelFolderSimpanan, err)
	}
	return app.String, nil
}

// CatatLampiran mencatat objek dan baris lampirannya dalam SATU transaksi;
// mengembalikan `ID` baris lampiran.
func (g *Gudang) CatatLampiran(ctx context.Context, l models.LampiranBaru) (id string, err error) {
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	if id, err = g.catatLampiranDalam(ctx, tx, l); err != nil {
		return "", err
	}
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("repository: mengikat lampiran %s: %w", l.NamaBerkas, err)
	}
	return id, nil
}

// CatatLampiranLaluBatalkanUntukUji - kedua INSERT di dalam transaksi yang
// SELALU dibatalkan; mengembalikan `ID` dan cacah baris yang sempat tertulis.
func (g *Gudang) CatatLampiranLaluBatalkanUntukUji(ctx context.Context, l models.LampiranBaru) (string, int, int, error) {
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return "", 0, 0, err
	}
	defer func() { _ = tx.Rollback() }()
	id, err := g.catatLampiranDalam(ctx, tx, l)
	if err != nil {
		return "", 0, 0, err
	}
	hitung := func(tabel, kolom, nilai string) (int, error) {
		nama, err := g.db.Qualify(tabel)
		if err != nil {
			return 0, err
		}
		var n int
		err = tx.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = :1", nama, kolom), nilai).Scan(&n)
		return n, err
	}
	nObjek, err := hitung(TabelObjekSimpanan, "IMAGEID", l.ImageID)
	if err != nil {
		return "", 0, 0, err
	}
	nLampiran, err := hitung(TabelWarisanLampiran, "ID", id)
	return id, nObjek, nLampiran, err
}

func (g *Gudang) catatLampiranDalam(ctx context.Context, tx *db.Tx, l models.LampiranBaru) (string, error) {
	objek, err := g.db.Qualify(TabelObjekSimpanan)
	if err != nil {
		return "", err
	}
	lampiran, err := g.db.Qualify(TabelWarisanLampiran)
	if err != nil {
		return "", err
	}
	// `Insert_T_Storage_SQL` — STORAGE 'standard', EXPDATE dari `exp` jawaban.
	q1 := fmt.Sprintf(`INSERT INTO %s (IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE)
		VALUES (:1, :2, :3, TO_DATE(:4, 'DD/MM/YYYY HH24:MI:SS'), :5, :6, 'standard')`, objek)
	if err := db.PeriksaSQL(q1); err != nil {
		return "", err
	}
	if _, err := tx.ExecContext(ctx, q1, l.ImageID, l.URLPublik, l.AppFolder, l.Exp, l.NamaObjek, l.App); err != nil {
		return "", fmt.Errorf("repository: mencatat objek %s: %w", TabelObjekSimpanan, err)
	}
	// `InsertAttachment2_Sql` — ID `TO_CHAR(SYSTIMESTAMP, 'YYYYMMDDHH24MISSFF3')`,
	// dibaca lebih dulu supaya dapat dikembalikan.
	var id string
	if err := tx.QueryRowContext(ctx, "SELECT TO_CHAR(SYSTIMESTAMP, 'YYYYMMDDHH24MISSFF3') FROM DUAL").Scan(&id); err != nil {
		return "", fmt.Errorf("repository: menyusun pengenal lampiran: %w", err)
	}
	q2 := fmt.Sprintf(`INSERT INTO %s (ID, TREATYID, CATEGORY, FILENAME, FILEMIMETYPE, USERNAME, CATEGORY_ID, T_STORAGE_ID)
		VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`, lampiran)
	if err := db.PeriksaSQL(q2); err != nil {
		return "", err
	}
	if _, err := tx.ExecContext(ctx, q2, id, l.IDKontrak, l.NamaKategori, l.NamaBerkas, l.Ekstensi, l.Pengguna,
		l.KodeKategori, l.ImageID); err != nil {
		return "", fmt.Errorf("repository: mencatat lampiran %s: %w", TabelWarisanLampiran, err)
	}
	return id, nil
}

// BacaObjekSimpanan - `GetLinkStorage_SQL`. `ada` palsu bila barisnya tidak ada.
func (g *Gudang) BacaObjekSimpanan(ctx context.Context, imageID string) (models.ObjekSimpanan, bool, error) {
	nama, err := g.db.Qualify(TabelObjekSimpanan)
	if err != nil {
		return models.ObjekSimpanan{}, false, err
	}
	q := fmt.Sprintf(`SELECT URLPUBLIC, APPFOLDER, TO_CHAR(EXPDATE, 'DD/MM/YYYY HH24:MI:SS'), APPNAME, FILENAME
		FROM %s WHERE IMAGEID = :1 FETCH FIRST 1 ROWS ONLY`, nama)
	if err := db.PeriksaSQL(q); err != nil {
		return models.ObjekSimpanan{}, false, err
	}
	var url, folder, exp, app, berkas sql.NullString
	err = g.db.QueryRowContext(ctx, q, imageID).Scan(&url, &folder, &exp, &app, &berkas)
	if err == sql.ErrNoRows {
		return models.ObjekSimpanan{}, false, nil
	}
	if err != nil {
		return models.ObjekSimpanan{}, false, fmt.Errorf("repository: membaca %s %s: %w", TabelObjekSimpanan, imageID, err)
	}
	return models.ObjekSimpanan{ImageID: imageID, URLPublik: url.String, AppFolder: folder.String, Exp: exp.String,
		App: app.String, NamaObjek: berkas.String}, true, nil
}

// PerbaruiObjekSimpanan - `Update_T_Storage_SQL` (`GetUrlGoogleStorage_Act`
// [6.6]): URL baru, APPFOLDER, EXPDATE, dan TANGGAL_UPLOAD dari `DateTime`
// jawaban (`MM/DD/YYYY HH24:MI:SS`). Bentuk yang tak terbaca → NULL.
func (g *Gudang) PerbaruiObjekSimpanan(ctx context.Context, o models.ObjekSimpanan, tanggal string) error {
	nama, err := g.db.Qualify(TabelObjekSimpanan)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`UPDATE %s SET URLPUBLIC = :1, APPFOLDER = :2,
		EXPDATE = TO_DATE(:3, 'DD/MM/YYYY HH24:MI:SS'), TANGGAL_UPLOAD = TO_DATE(:4, 'MM/DD/YYYY HH24:MI:SS')
		WHERE IMAGEID = :5`, nama)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := g.db.ExecContext(ctx, q, o.URLPublik, o.AppFolder, o.Exp, tanggal, o.ImageID); err != nil {
		return fmt.Errorf("repository: memperbarui %s %s: %w", TabelObjekSimpanan, o.ImageID, err)
	}
	return nil
}

// HapusLampiran - `DeleteStorage_SQL` (`DeleteGoogleStorage_Act` [10]) lalu
// `DeleteAttachment2_Sql` (`Delete_act` [3]), SATU transaksi.
func (g *Gudang) HapusLampiran(ctx context.Context, idKontrak, idLampiran, imageID string) (err error) {
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	if err = g.hapusLampiranDalam(ctx, tx, idKontrak, idLampiran, imageID); err != nil {
		return err
	}
	return tx.Commit()
}

// HapusLampiranLaluBatalkanUntukUji - penghapusan di dalam transaksi yang
// SELALU dibatalkan; mengembalikan cacah baris yang terhapus.
func (g *Gudang) HapusLampiranLaluBatalkanUntukUji(ctx context.Context, idKontrak, idLampiran, imageID string) (int64, int64, error) {
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer func() { _ = tx.Rollback() }()
	objek, lampiran, err := g.hapusLampiranHitung(ctx, tx, idKontrak, idLampiran, imageID)
	return objek, lampiran, err
}

func (g *Gudang) hapusLampiranDalam(ctx context.Context, tx *db.Tx, idKontrak, idLampiran, imageID string) error {
	_, _, err := g.hapusLampiranHitung(ctx, tx, idKontrak, idLampiran, imageID)
	return err
}

func (g *Gudang) hapusLampiranHitung(ctx context.Context, tx *db.Tx, idKontrak, idLampiran, imageID string) (int64, int64, error) {
	objek, err := g.db.Qualify(TabelObjekSimpanan)
	if err != nil {
		return 0, 0, err
	}
	lampiran, err := g.db.Qualify(TabelWarisanLampiran)
	if err != nil {
		return 0, 0, err
	}
	var nObjek int64
	if imageID != "" {
		q1 := fmt.Sprintf("DELETE FROM %s WHERE IMAGEID = :1", objek)
		if err := db.PeriksaSQL(q1); err != nil {
			return 0, 0, err
		}
		r, err := tx.ExecContext(ctx, q1, imageID)
		if err != nil {
			return 0, 0, fmt.Errorf("repository: menghapus %s %s: %w", TabelObjekSimpanan, imageID, err)
		}
		nObjek, _ = r.RowsAffected()
	}
	q2 := fmt.Sprintf("DELETE FROM %s WHERE TREATYID = :1 AND ID = :2", lampiran)
	if err := db.PeriksaSQL(q2); err != nil {
		return 0, 0, err
	}
	r, err := tx.ExecContext(ctx, q2, idKontrak, idLampiran)
	if err != nil {
		return 0, 0, fmt.Errorf("repository: menghapus %s %s: %w", TabelWarisanLampiran, idLampiran, err)
	}
	nLampiran, _ := r.RowsAffected()
	return nObjek, nLampiran, nil
}

// UbahKategoriLampiran - `ChangeKateAttachment2_Sql` untuk setiap baris,
// SATU transaksi: `CATEGORY_ID` dan `CATEGORY` (nama dari katalog).
func (g *Gudang) UbahKategoriLampiran(ctx context.Context, idKontrak string, ubah []models.PerubahanKategori) (err error) {
	nama, err := g.db.Qualify(TabelWarisanLampiran)
	if err != nil {
		return err
	}
	q := fmt.Sprintf("UPDATE %s SET CATEGORY_ID = :1, CATEGORY = :2 WHERE TREATYID = :3 AND ID = :4", nama)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	for _, u := range ubah {
		if _, err = tx.ExecContext(ctx, q, u.KodeKategori, u.NamaKategori, idKontrak, u.IDLampiran); err != nil {
			return fmt.Errorf("repository: mengubah kategori %s %s: %w", TabelWarisanLampiran, u.IDLampiran, err)
		}
	}
	return tx.Commit()
}
