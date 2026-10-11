package repository

// Lampiran case NB FacIn (tiket 67) - tabel Pega POOLDATA (keputusan work owner 08-10-2026 K67-1 "disimpan ke tabel
// DOCUMENT_POLIS, T_STORAGE_IMAGE"), SQL persis rule-nya `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\`:
//
//	CategoryAttach_SQL           `select * from CATEGORY_ATTACH_REAS order by note` (hanya NOTE dibaca)
//	InputParamUploadReas_act 5   Obj-Browse Int-DOCUMENT_POLIS `.IDPEGA = pzInsKey AND .KATEGORI_2 = kategori` (Count)
//	GetAllDocument (RD)          ID, T_STORAGE_ID, INSKEY_LINK, MIME, NAMAFILE, KATEGORI_2, TANGGAL; KATEGORI_2 + IDPEGA; maks 500
//	InsertDocument_Act 3 / 5     Obj-Save DOCUMENT_POLIS (ID cap waktu, KATEGORI_1 "Reas", INSKEY_* "")
//	DeleteDocumentPolis_Act 3    Obj-Delete (Immediate, Remove) = hapus FISIK (K-A: "ikutin pega aja", pengecualian ADR-0031)
//	Insert_T_Storage_SQL / Update_T_Storage_SQL / GetLinkStorage_SQL / DeleteStorage_SQL   T_STORAGE_IMAGE
//	GetAppName_SQL               `SELECT APPNAME … FROM POOLDATA.T_FOLDER_IMAGE`
//
// Status unggah dibaca dari outbox bersama T_LOG_SERVICE_RNM (MODUL "NBFACIN"). Nol COMMIT (ADR-U-0029), nol CALL
// (ADR-0043): blok PL/SQL `begin insert … commit; end;` Insert_T_Storage_SQL menjadi INSERT biasa di transaksi pemanggil.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/nbfacin/backend/models"
)

// Tabel lampiran (POOLDATA, warisan Pega).
const (
	TabelDokumenPolis       = "DOCUMENT_POLIS"
	TabelStorageImage       = "T_STORAGE_IMAGE"
	TabelKategoriLampiran   = "CATEGORY_ATTACH_REAS"
	TabelFolderImage        = "T_FOLDER_IMAGE"
	TabelTokenStorage       = "GCP_IMAGE"
	TabelOutboxBersama      = "T_LOG_SERVICE_RNM"
	batasBarisDokumenPolis  = 500 // GetAllDocument pyMaxRecords
	batasGeserIDDokumen     = 1000
	kodeOracleUnikDilanggar = "ORA-00001"
)

var (
	// ErrLampiranTidakAda - dokumen tidak ada, atau bukan milik case itu.
	ErrLampiranTidakAda = errors.New("repository: dokumen lampiran tidak ada di case ini")
	// ErrAppStorageKosong - T_FOLDER_IMAGE tidak memberi APPNAME.
	ErrAppStorageKosong = errors.New("repository: APPNAME penyimpanan tidak ada di T_FOLDER_IMAGE")
)

// PenyimpanLampiran - pembaca / penulis lampiran case (tiket 67).
type PenyimpanLampiran interface {
	KategoriLampiran(ctx context.Context) ([]string, error)
	JumlahLampiran(ctx context.Context, idPega string) (map[string]int, error)
	DaftarLampiran(ctx context.Context, idPega, kategori, modul, jenisUnggah string) ([]models.DokumenPolis, error)
	AmbilLampiran(ctx context.Context, idPega, id, modul, jenisUnggah string) (models.DokumenPolis, error)
	// KunciLampiran - FOR UPDATE satu dokumen (pelaksana efek); ErrLampiranTidakAda bila sudah dihapus.
	KunciLampiran(ctx context.Context, tx *db.Tx, id string) (models.DokumenPolis, error)
	SisipLampiran(ctx context.Context, tx *db.Tx, d models.DokumenPolis, saat time.Time) (string, error)
	HapusLampiran(ctx context.Context, tx *db.Tx, idPega, id string) error
	// LinkStorage - GetLinkStorage_SQL; false bila IMAGEID tidak tercatat. tx nil = pool.
	LinkStorage(ctx context.Context, tx *db.Tx, imageID string) (models.LinkStorage, bool, error)
	SimpanObjek(ctx context.Context, tx *db.Tx, o models.ObjekStorage) error
	PerbaruiObjek(ctx context.Context, tx *db.Tx, o models.ObjekStorage) error
	HapusObjek(ctx context.Context, tx *db.Tx, imageID string) error
	AppStorage(ctx context.Context) (string, error)
	TokenStorageBerlaku(ctx context.Context, tx *db.Tx, app string, saat time.Time, sisaMinimum time.Duration) (string, time.Duration, error)
	PungutEfekRujukan(ctx context.Context, tx *db.Tx, modul, rujukan string, saat time.Time) (outbox.BarisEfekKeluar, error)
}

// LampiranOracle - PenyimpanLampiran atas Oracle.
type LampiranOracle struct{ db *db.DB }

// NewLampiranOracle merakitnya.
func NewLampiranOracle(d *db.DB) *LampiranOracle { return &LampiranOracle{db: d} }

func sqlKategoriLampiran(t string) string {
	return "SELECT NOTE FROM " + t + " WHERE NOTE IS NOT NULL ORDER BY NOTE"
}

func sqlJumlahLampiran(t string) string {
	return "SELECT KATEGORI_2, COUNT(*) FROM " + t + " WHERE IDPEGA IN (:1, :2) AND KATEGORI_2 IS NOT NULL GROUP BY KATEGORI_2" // IDPEGA dua bentuk (idpegaproduksi.go)
}

// sqlDaftarLampiran - GetAllDocument + keberadaan objek (subkueri, bukan JOIN) + efek unggah TERAKHIR per dokumen.
// GetAllDocument tanpa urutan (pySortOrder 99999): ORDER BY ID `[keputusan agent]` supaya urutan stabil.
func sqlDaftarLampiran(dok, obj, ob string, satu bool) string {
	saring := "a.KATEGORI_2 = :3"
	if satu {
		saring = "a.ID = :3"
	}
	return fmt.Sprintf(`SELECT a.ID, a.NAMAFILE, a.MIME, a.KATEGORI_1, a.KATEGORI_2, a.T_STORAGE_ID, a.PXCREATEOPERATOR,
	       TO_CHAR(a.TANGGAL, 'DD-MM-YYYY HH24:MI'),
	       (SELECT COUNT(*) FROM %s s WHERE s.IMAGEID = a.T_STORAGE_ID),
	       (SELECT MAX(o.STATUS) KEEP (DENSE_RANK LAST ORDER BY o.DIBUAT, o.ID) FROM %s o
	         WHERE o.MODUL = :1 AND o.JENIS_EFEK = :2 AND o.RUJUKAN = a.ID)
	  FROM %s a
	 WHERE %s AND a.IDPEGA IN (:4, :5)
	 ORDER BY a.ID
	 FETCH FIRST %d ROWS ONLY`, obj, ob, dok, saring, batasBarisDokumenPolis)
}

func sqlKunciLampiran(dok, obj string) string {
	return fmt.Sprintf(`SELECT a.ID, a.IDPEGA, a.NAMAFILE, a.MIME, a.T_STORAGE_ID,
	       (SELECT COUNT(*) FROM %s s WHERE s.IMAGEID = a.T_STORAGE_ID)
	  FROM %s a WHERE a.ID = :1 FOR UPDATE OF a.ID`, obj, dok)
}

func sqlAdaIDLampiran(dok string) string { return "SELECT COUNT(*) FROM " + dok + " WHERE ID = :1" }

// sqlSisipLampiran - InsertDocument_Act langkah 3: TANGGAL = SYSDATE (DDL DEFAULT sysdate; Pega @CurrentDateTime),
// INSKEY_LINK / INSKEY_DATA "" = NULL.
func sqlSisipLampiran(dok string) string {
	return "INSERT INTO " + dok + " (ID, IDPEGA, TANGGAL, NAMAFILE, MIME, KATEGORI_1, KATEGORI_2, INSKEY_LINK, INSKEY_DATA, " +
		"T_STORAGE_ID, PXCREATEOPERATOR) VALUES (:1, :2, SYSDATE, :3, :4, :5, :6, NULL, NULL, :7, :8)"
}

func sqlHapusLampiran(dok string) string {
	return "DELETE FROM " + dok + " WHERE ID = :1 AND IDPEGA IN (:2, :3)"
}

func sqlLinkStorage(obj string) string {
	return "SELECT URLPUBLIC, APPFOLDER, TO_CHAR(EXPDATE, 'YYYYMMDDHH24MISS'), APPNAME, FILENAME FROM " + obj +
		" WHERE IMAGEID = :1"
}

func sqlSimpanObjekStorage(obj string) string {
	return "INSERT INTO " + obj + " (IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE) " +
		"VALUES (:1, :2, :3, TO_DATE(:4, 'DD/MM/YYYY HH24:MI:SS'), :5, :6, 'standard')"
}

func sqlPerbaruiObjekStorage(obj string) string {
	return "UPDATE " + obj + " SET URLPUBLIC = :1, APPFOLDER = :2, EXPDATE = TO_DATE(:3, 'DD/MM/YYYY HH24:MI:SS'), " +
		"TANGGAL_UPLOAD = TO_DATE(:4, 'MM/DD/YYYY HH24:MI:SS') WHERE IMAGEID = :5"
}

func sqlHapusObjekStorage(obj string) string { return "DELETE FROM " + obj + " WHERE IMAGEID = :1" }

func sqlAppStorage(t string) string {
	return "SELECT APPNAME FROM " + t + " WHERE APPNAME IS NOT NULL FETCH FIRST 1 ROWS ONLY"
}

// sqlTokenStorageBerlaku - sisa umur dihitung DI ORACLE terhadap waktu yang di-bind (pola Treaty Contract Out).
func sqlTokenStorageBerlaku(t string) string {
	return "SELECT KODEAKSES, ROUND((CAST(INPUTDATE AS DATE) - CAST(:1 AS DATE)) * 86400) FROM " + t +
		" WHERE APPNAME = :2 AND INPUTDATE > :3 ORDER BY INPUTDATE DESC FETCH FIRST 1 ROWS ONLY"
}

// sqlPungutEfekRujukan - pungutan outbox bersama + `AND RUJUKAN = :4` (efek MILIK aksi pemakai itu saja). ⛔ Batasan Oracle:
// `FETCH FIRST` tidak boleh bersama `FOR UPDATE` (ORA-02014, bug DEV 10-10-2026, K67-8) - kandidat lewat subkueri ROWNUM, kunci
// di tabel dasar; bentuk yang sama dengan outbox.sqlPungutEfek (keputusan work owner 10-10-2026, cara 1).
func sqlPungutEfekRujukan(t string) string {
	return "SELECT ID, LINI, MODUL, JENIS_EFEK, RUJUKAN, MUATAN, PERCOBAAN FROM " + t +
		" WHERE ID = (SELECT ID FROM (SELECT ID FROM " + t +
		" WHERE STATUS = :1 AND JADWAL_BERIKUT <= :2 AND MODUL = :3 AND RUJUKAN = :4" +
		" ORDER BY JADWAL_BERIKUT ASC) WHERE ROWNUM = 1) FOR UPDATE SKIP LOCKED"
}

func (r *LampiranOracle) tabel(nama ...string) ([]string, error) {
	out := make([]string, len(nama))
	for i, n := range nama {
		t, err := r.db.Qualify(n)
		if err != nil {
			return nil, err
		}
		out[i] = t
	}
	return out, nil
}

func wajibTx(tx *db.Tx, apa string) error {
	if tx == nil {
		return errors.New("repository: " + apa + " menuntut transaksi")
	}
	return nil
}

// KategoriLampiran - CategoryAttach_SQL (tanpa DISTINCT, seperti Pega).
func (r *LampiranOracle) KategoriLampiran(ctx context.Context) ([]string, error) {
	t, err := r.tabel(TabelKategoriLampiran)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, sqlKategoriLampiran(t[0]))
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s: %w", TabelKategoriLampiran, err)
	}
	defer func() { _ = rows.Close() }()
	var out []string
	for rows.Next() {
		var n sql.NullString
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		if strings.TrimSpace(n.String) != "" {
			out = append(out, n.String)
		}
	}
	return out, rows.Err()
}

// JumlahLampiran - Count per KATEGORI_2 dokumen case (KATEGORI_1 TIDAK disaring - dokumen buatan PDF ikut).
func (r *LampiranOracle) JumlahLampiran(ctx context.Context, idPega string) (map[string]int, error) {
	t, err := r.tabel(TabelDokumenPolis)
	if err != nil {
		return nil, err
	}
	polos, awalan := DuaBentukIDPega(idPega)
	rows, err := r.db.QueryContext(ctx, sqlJumlahLampiran(t[0]), polos, awalan)
	if err != nil {
		return nil, fmt.Errorf("repository: menghitung %s: %w", TabelDokumenPolis, err)
	}
	defer func() { _ = rows.Close() }()
	out := map[string]int{}
	for rows.Next() {
		var k sql.NullString
		var n int
		if err := rows.Scan(&k, &n); err != nil {
			return nil, err
		}
		out[k.String] += n
	}
	return out, rows.Err()
}

func (r *LampiranOracle) baca(ctx context.Context, idPega, kunci, modul, jenis string, satu bool) ([]models.DokumenPolis, error) {
	t, err := r.tabel(TabelDokumenPolis, TabelStorageImage, TabelOutboxBersama)
	if err != nil {
		return nil, err
	}
	polos, awalan := DuaBentukIDPega(idPega)
	rows, err := r.db.QueryContext(ctx, sqlDaftarLampiran(t[0], t[1], t[2], satu), modul, jenis, kunci, polos, awalan)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s: %w", TabelDokumenPolis, err)
	}
	defer func() { _ = rows.Close() }()
	var out []models.DokumenPolis
	for rows.Next() {
		var n [8]sql.NullString
		var ada int
		var status sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6], &n[7], &ada, &status); err != nil {
			return nil, err
		}
		out = append(out, models.DokumenPolis{ID: n[0].String, IDPega: idPega, NamaFile: n[1].String, Mime: n[2].String,
			Kategori1: n[3].String, Kategori2: n[4].String, TStorageID: n[5].String, Operator: n[6].String,
			Tanggal: n[7].String, Tersimpan: ada > 0 && n[5].String != "", StatusEfek: status.String})
	}
	return out, rows.Err()
}

// DaftarLampiran - GetAllDocument satu kategori.
func (r *LampiranOracle) DaftarLampiran(ctx context.Context, idPega, kategori, modul, jenis string) ([]models.DokumenPolis, error) {
	return r.baca(ctx, idPega, kategori, modul, jenis, false)
}

// AmbilLampiran - satu dokumen, DIBATASI case-nya (ID dari jalur URL).
func (r *LampiranOracle) AmbilLampiran(ctx context.Context, idPega, id, modul, jenis string) (models.DokumenPolis, error) {
	d, err := r.baca(ctx, idPega, id, modul, jenis, true)
	if err != nil {
		return models.DokumenPolis{}, err
	}
	if len(d) == 0 {
		return models.DokumenPolis{}, ErrLampiranTidakAda
	}
	return d[0], nil
}

// KunciLampiran - FOR UPDATE: dua pelaksana atas dokumen yang sama saling menunggu, yang kedua melihat objek tercatat.
func (r *LampiranOracle) KunciLampiran(ctx context.Context, tx *db.Tx, id string) (models.DokumenPolis, error) {
	if err := wajibTx(tx, "mengunci dokumen lampiran"); err != nil {
		return models.DokumenPolis{}, err
	}
	t, err := r.tabel(TabelDokumenPolis, TabelStorageImage)
	if err != nil {
		return models.DokumenPolis{}, err
	}
	var n [5]sql.NullString
	var ada int
	err = tx.QueryRowContext(ctx, sqlKunciLampiran(t[0], t[1]), id).Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &ada)
	if errors.Is(err, sql.ErrNoRows) {
		return models.DokumenPolis{}, ErrLampiranTidakAda
	}
	if err != nil {
		return models.DokumenPolis{}, fmt.Errorf("repository: mengunci dokumen %s: %w", id, err)
	}
	return models.DokumenPolis{ID: n[0].String, IDPega: n[1].String, NamaFile: n[2].String, Mime: n[3].String,
		TStorageID: n[4].String, Tersimpan: ada > 0 && n[4].String != ""}, nil
}

// SisipLampiran - InsertDocument_Act: ID cap waktu Pega (models.IDDokumenPolis). ID terpakai (dua unggahan di milidetik yang
// sama, atau jam 12-an `hh` yang membuat 01:00 = 13:00) digeser satu milidetik, sampai 1000 kali (pola Treaty Contract Out);
// ORA-00001 dari penulis lain (Pega / generator PDF) di antara cek dan INSERT juga digeser (K-D).
func (r *LampiranOracle) SisipLampiran(ctx context.Context, tx *db.Tx, d models.DokumenPolis, saat time.Time) (string, error) {
	if err := wajibTx(tx, "menulis dokumen lampiran"); err != nil {
		return "", err
	}
	t, err := r.tabel(TabelDokumenPolis)
	if err != nil {
		return "", err
	}
	for i := 0; i < batasGeserIDDokumen; i, saat = i+1, saat.Add(time.Millisecond) {
		id := models.IDDokumenPolis(saat)
		var n int
		if err := tx.QueryRowContext(ctx, sqlAdaIDLampiran(t[0]), id).Scan(&n); err != nil {
			return "", fmt.Errorf("repository: memeriksa ID dokumen: %w", err)
		}
		if n > 0 {
			continue
		}
		hasil, err := tx.ExecContext(ctx, sqlSisipLampiran(t[0]), id, d.IDPega, db.KosongJadiNil(d.NamaFile),
			db.KosongJadiNil(d.Mime), db.KosongJadiNil(d.Kategori1), db.KosongJadiNil(d.Kategori2),
			db.KosongJadiNil(d.TStorageID), db.KosongJadiNil(d.Operator))
		if err != nil && strings.Contains(err.Error(), kodeOracleUnikDilanggar) {
			continue
		}
		if err != nil {
			return "", fmt.Errorf("repository: menulis %s: %w", TabelDokumenPolis, err)
		}
		return id, db.PastikanSatuBaris(hasil, "dokumen lampiran")
	}
	return "", fmt.Errorf("repository: ID dokumen lampiran dan %d penggantinya sudah terpakai", batasGeserIDDokumen)
}

// HapusLampiran - DeleteDocumentPolis_Act langkah 3 Obj-Delete: hapus FISIK (K-A), dibatasi case.
func (r *LampiranOracle) HapusLampiran(ctx context.Context, tx *db.Tx, idPega, id string) error {
	if err := wajibTx(tx, "menghapus dokumen lampiran"); err != nil {
		return err
	}
	t, err := r.tabel(TabelDokumenPolis)
	if err != nil {
		return err
	}
	polos, awalan := DuaBentukIDPega(idPega)
	hasil, err := tx.ExecContext(ctx, sqlHapusLampiran(t[0]), id, polos, awalan)
	if err != nil {
		return fmt.Errorf("repository: menghapus dokumen %s: %w", id, err)
	}
	if n, err := hasil.RowsAffected(); err == nil && n == 0 {
		return ErrLampiranTidakAda
	}
	return db.PastikanSatuBaris(hasil, "hapus dokumen lampiran")
}

// LinkStorage - GetLinkStorage_SQL.
func (r *LampiranOracle) LinkStorage(ctx context.Context, tx *db.Tx, imageID string) (models.LinkStorage, bool, error) {
	t, err := r.tabel(TabelStorageImage)
	if err != nil {
		return models.LinkStorage{}, false, err
	}
	var n [5]sql.NullString
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, sqlLinkStorage(t[0]), imageID)
	} else {
		row = r.db.QueryRowContext(ctx, sqlLinkStorage(t[0]), imageID)
	}
	err = row.Scan(&n[0], &n[1], &n[2], &n[3], &n[4])
	if errors.Is(err, sql.ErrNoRows) {
		return models.LinkStorage{}, false, nil
	}
	if err != nil {
		// Galat Oracle tidak menggemakan nilai kolom; URL tidak ditulis di pesan ini.
		return models.LinkStorage{}, false, fmt.Errorf("repository: membaca catatan objek %s: %w", imageID, err)
	}
	return models.LinkStorage{ImageID: imageID, URLPublic: n[0].String, AppFolder: n[1].String, ExpUTC: n[2].String,
		App: n[3].String, Namafile: n[4].String}, true, nil
}

// SimpanObjek - Insert_T_Storage_SQL (tanpa blok PL/SQL / COMMIT).
func (r *LampiranOracle) SimpanObjek(ctx context.Context, tx *db.Tx, o models.ObjekStorage) error {
	if err := wajibTx(tx, "mencatat objek penyimpanan"); err != nil {
		return err
	}
	t, err := r.tabel(TabelStorageImage)
	if err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, sqlSimpanObjekStorage(t[0]), o.ImageID, db.KosongJadiNil(o.URLPublic),
		db.KosongJadiNil(o.AppFolder), db.KosongJadiNil(o.Exp), db.KosongJadiNil(o.Namafile), db.KosongJadiNil(o.App))
	if err != nil {
		return fmt.Errorf("repository: mencatat objek penyimpanan %s: %w", o.ImageID, err)
	}
	return db.PastikanSatuBaris(hasil, "catatan objek penyimpanan")
}

// PerbaruiObjek - Update_T_Storage_SQL sesudah geturl; nol baris bukan galat (seperti Pega).
func (r *LampiranOracle) PerbaruiObjek(ctx context.Context, tx *db.Tx, o models.ObjekStorage) error {
	if err := wajibTx(tx, "menyegarkan catatan objek"); err != nil {
		return err
	}
	t, err := r.tabel(TabelStorageImage)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, sqlPerbaruiObjekStorage(t[0]), db.KosongJadiNil(o.URLPublic), db.KosongJadiNil(o.AppFolder),
		db.KosongJadiNil(o.Exp), db.KosongJadiNil(o.TanggalUpload), o.ImageID); err != nil {
		return fmt.Errorf("repository: menyegarkan catatan objek %s: %w", o.ImageID, err)
	}
	return nil
}

// HapusObjek - DeleteStorage_SQL; nol baris bukan galat.
func (r *LampiranOracle) HapusObjek(ctx context.Context, tx *db.Tx, imageID string) error {
	if err := wajibTx(tx, "menghapus catatan objek"); err != nil {
		return err
	}
	t, err := r.tabel(TabelStorageImage)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, sqlHapusObjekStorage(t[0]), imageID); err != nil {
		return fmt.Errorf("repository: menghapus catatan objek %s: %w", imageID, err)
	}
	return nil
}

// AppStorage - GetAppName_SQL (InsertGoogleStorage_Act langkah 6: prasyarat IsPEGAPROD NONAKTIF, selalu dibaca).
func (r *LampiranOracle) AppStorage(ctx context.Context) (string, error) {
	t, err := r.tabel(TabelFolderImage)
	if err != nil {
		return "", err
	}
	var app sql.NullString
	err = r.db.QueryRowContext(ctx, sqlAppStorage(t[0])).Scan(&app)
	if errors.Is(err, sql.ErrNoRows) || (err == nil && strings.TrimSpace(app.String) == "") {
		return "", ErrAppStorageKosong
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca APPNAME penyimpanan: %w", err)
	}
	return app.String, nil
}

// TokenStorageBerlaku - token GCP_IMAGE yang masih berlaku LEBIH DARI sisaMinimum sesudah saat, beserta sisa umurnya.
// ⛔ Galat driver tidak diteruskan: pesannya dapat memuat nilai kolom kredensial.
func (r *LampiranOracle) TokenStorageBerlaku(ctx context.Context, tx *db.Tx, app string, saat time.Time,
	sisaMinimum time.Duration) (string, time.Duration, error) {
	if err := wajibTx(tx, "token penyimpanan"); err != nil {
		return "", 0, err
	}
	t, err := r.tabel(TabelTokenStorage)
	if err != nil {
		return "", 0, err
	}
	var kode sql.NullString
	var sisa sql.NullFloat64
	err = tx.QueryRowContext(ctx, sqlTokenStorageBerlaku(t[0]), saat, app, saat.Add(sisaMinimum)).Scan(&kode, &sisa)
	if errors.Is(err, sql.ErrNoRows) {
		return "", 0, nil
	}
	if err != nil {
		return "", 0, fmt.Errorf("repository: membaca token penyimpanan untuk %q", app)
	}
	return kode.String, time.Duration(sisa.Float64) * time.Second, nil
}

// PungutEfekRujukan - satu efek jatuh tempo milik satu rujukan, ditandai jalan (pola PungutEfek outbox bersama).
func (r *LampiranOracle) PungutEfekRujukan(ctx context.Context, tx *db.Tx, modul, rujukan string,
	saat time.Time) (outbox.BarisEfekKeluar, error) {
	if err := wajibTx(tx, "memungut efek"); err != nil {
		return outbox.BarisEfekKeluar{}, err
	}
	t, err := r.tabel(TabelOutboxBersama)
	if err != nil {
		return outbox.BarisEfekKeluar{}, err
	}
	var b outbox.BarisEfekKeluar
	var n [5]sql.NullString
	var percobaan sql.NullInt64
	err = tx.QueryRowContext(ctx, sqlPungutEfekRujukan(t[0]), outbox.StatusEfekAntre, saat, modul, rujukan).Scan(
		&b.ID, &n[0], &n[1], &n[2], &n[3], &n[4], &percobaan)
	if errors.Is(err, sql.ErrNoRows) {
		return outbox.BarisEfekKeluar{}, outbox.ErrEfekTidakAda
	}
	if err != nil {
		return outbox.BarisEfekKeluar{}, fmt.Errorf("repository: memungut efek %s: %w", rujukan, err)
	}
	b.Lini, b.Modul, b.Jenis, b.Rujukan, b.Muatan, b.Percobaan = n[0].String, n[1].String, n[2].String, n[3].String,
		n[4].String, int(percobaan.Int64)
	tandai := outbox.SQLTandaiJalan(t[0])
	hasil, err := tx.ExecContext(ctx, tandai, outbox.StatusEfekJalan, saat, b.ID, outbox.StatusEfekAntre)
	if err != nil {
		return outbox.BarisEfekKeluar{}, fmt.Errorf("repository: menandai efek berjalan: %w", err)
	}
	if err := db.PastikanSatuBaris(hasil, "tanda efek berjalan"); err != nil {
		return outbox.BarisEfekKeluar{}, err
	}
	return b, nil
}
