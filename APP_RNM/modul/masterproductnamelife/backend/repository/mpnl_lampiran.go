package repository

// Lampiran produk (paket 8, tiket 08–09) - tabel warisan seperti XML (P5):
//
//	InsertAttachProdName_Sql b84  INSERT M_ATTACHMENTPRODUCTNAME (ID = TO_CHAR(SYSTIMESTAMP,'YYYYMMDDHH24MISSFF3'),
//	                              TREATYID, CATEGORY, FILENAME, FILEMIMETYPE, DATA_JSON, USERNAME, T_STORAGE_ID)
//	GetAttachmentProdName_Sql b84 SELECT … WHERE treatyid = ProductName.ID
//	DeleteAttachProdName_Sql b85  DELETE … WHERE treatyid = … AND id = …
//	Insert_T_Storage_SQL b85      INSERT T_STORAGE_IMAGE (IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE='standard')
//	DeleteStorage_SQL b85         DELETE T_STORAGE_IMAGE WHERE imageid = …
//	GetAppName_SQL b58            SELECT APPNAME FROM POOLDATA.T_FOLDER_IMAGE
//
// ⛔ Pengiriman berkas = EFEK KELUAR lewat outbox bersama `T_LOG_SERVICE_RNM`
// (preseden Treaty Contract Out): rekam lampiran + antrean di SATU transaksi;
// pelaksana stub menulis `T_STORAGE_IMAGE` saat berhasil. `URLPUBLIC` NULL -
// alamat nyata tidak dipanggil dan tidak ditulis (OQ-MPNL-10).
// ⛔ Status "terunggah" = baris `T_STORAGE_IMAGE` untuk `T_STORAGE_ID` ada.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/masterproductnamelife/backend/models"
)

// Tabel warisan lampiran dan outbox bersama.
const (
	TabelLampiran = "M_ATTACHMENTPRODUCTNAME"
	TabelObjek    = "T_STORAGE_IMAGE"
	TabelFolder   = "T_FOLDER_IMAGE"
	TabelOutbox   = "T_LOG_SERVICE_RNM"
)

// Penanda outbox modul ini.
const (
	ModulOutbox     = "MASTERPRODUCTNAMELIFE"
	JenisEfekUnggah = "unggah-lampiran"
)

// Kolom yang disebut SQL lampiran (penjaga kata cadangan).
var (
	KolomLampiran = []string{"ID", "TREATYID", "CATEGORY", "FILENAME", "FILEMIMETYPE", "DATA_JSON", "USERNAME", "T_STORAGE_ID"}
	KolomObjek    = []string{"IMAGEID", "URLPUBLIC", "APPFOLDER", "EXPDATE", "FILENAME", "APPNAME", "STORAGE"}
	KolomOutbox   = []string{"ID", "MODUL", "RUJUKAN", "STATUS", "PERCOBAAN", "GALAT_TERAKHIR", "DIBUAT"}
)

var (
	// ErrLampiranTidakAda - bukan lampiran produk itu (404).
	ErrLampiranTidakAda = errors.New("repository: attachment not found")
	// ErrIDLampiranBentrok - stempel waktu milidetik dipakai baris lain berulang kali.
	ErrIDLampiranBentrok = errors.New("repository: could not obtain a free attachment ID")
)

func sqlWaktuID() string { return `SELECT TO_CHAR(SYSTIMESTAMP, 'YYYYMMDDHH24MISSFF3') FROM DUAL` }

func sqlSisipLampiran(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, TREATYID, CATEGORY, FILENAME, FILEMIMETYPE, DATA_JSON, USERNAME, T_STORAGE_ID)
		VALUES (:1, :2, :3, :4, :5, NULL, :6, :7)`, tabel)
}

func sqlDaftarLampiran(tabel string) string {
	return fmt.Sprintf(`SELECT ID, TREATYID, CATEGORY, FILENAME, FILEMIMETYPE, USERNAME, T_STORAGE_ID
		FROM %s WHERE TREATYID = :1 ORDER BY ID ASC`, tabel)
}

func sqlAmbilLampiran(tabel string) string {
	return fmt.Sprintf(`SELECT ID, TREATYID, CATEGORY, FILENAME, FILEMIMETYPE, USERNAME, T_STORAGE_ID
		FROM %s WHERE TREATYID = :1 AND ID = :2`, tabel)
}

func sqlHapusLampiran(tabel string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE TREATYID = :1 AND ID = :2`, tabel)
}

func sqlAdaObjek(tabel string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE IMAGEID = :1`, tabel)
}

func sqlSisipObjek(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE)
		VALUES (:1, NULL, :2, SYSDATE + :3 / 86400, :4, :5, 'standard')`, tabel)
}

func sqlHapusObjek(tabel string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE IMAGEID = :1`, tabel)
}

func sqlNamaAplikasi(tabel string) string {
	return fmt.Sprintf(`SELECT APPNAME FROM %s FETCH FIRST 1 ROWS ONLY`, tabel)
}

// sqlPungutUnggah - efek antre satu lampiran. ⛔ Tanpa `FETCH FIRST`: Oracle
// menolak batas baris bersama `FOR UPDATE` (ORA-02014); baris pertama diambil di Go.
func sqlPungutUnggah(tabel string) string {
	return fmt.Sprintf(`SELECT ID, PERCOBAAN FROM %s
		WHERE MODUL = :1 AND RUJUKAN = :2 AND STATUS = :3
		ORDER BY DIBUAT ASC FOR UPDATE SKIP LOCKED`, tabel)
}

// sqlAdaUnggahAntre - efek antre lampiran ini ada? (tanpa mengunci, tanpa menaikkan percobaan).
func sqlAdaUnggahAntre(tabel string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE MODUL = :1 AND RUJUKAN = :2 AND STATUS = :3`, tabel)
}

func sqlStatusUnggah(tabel string) string {
	return fmt.Sprintf(`SELECT STATUS, GALAT_TERAKHIR FROM %s
		WHERE MODUL = :1 AND RUJUKAN = :2 ORDER BY DIBUAT DESC FETCH FIRST 1 ROWS ONLY`, tabel)
}

func pindaiLampiran(rows interface{ Scan(...any) error }) (models.Lampiran, error) {
	var id, produk, kat, nama, ext, user, storage sql.NullString
	if err := rows.Scan(&id, &produk, &kat, &nama, &ext, &user, &storage); err != nil {
		return models.Lampiran{}, err
	}
	return models.Lampiran{ID: id.String, ProdukID: produk.String, Category: kat.String, FileName: nama.String,
		FileMimeType: ext.String, UserName: user.String, StorageID: storage.String}, nil
}

// lengkapiStatus - terunggah bila objeknya tercatat; selain itu dari outbox.
func (g *Gudang) lengkapiStatus(ctx context.Context, tx *db.Tx, l *models.Lampiran) error {
	ada, err := g.AdaObjek(ctx, tx, l.StorageID)
	if err != nil {
		return err
	}
	if ada {
		l.Status = models.StatusTerunggah
		return nil
	}
	q, err := g.siapkan(TabelOutbox, sqlStatusUnggah)
	if err != nil {
		return err
	}
	var status, galat sql.NullString
	err = g.kueri(tx).QueryRowContext(ctx, q, ModulOutbox, l.ID).Scan(&status, &galat)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		l.Status = models.StatusBelum
	case err != nil:
		return fmt.Errorf("repository: reading attachment %s send status: %w", l.ID, err)
	case galat.String != "" || status.String == outbox.StatusEfekGagalPermanen:
		l.Status, l.Galat = models.StatusGagal, galat.String
	default:
		l.Status = models.StatusBelum
	}
	return nil
}

// DaftarLampiran - grid lampiran produk (`LoadAttachmentProdName` 2 b370), beserta status.
func (g *Gudang) DaftarLampiran(ctx context.Context, produkID string) ([]models.Lampiran, error) {
	q, err := g.siapkan(TabelLampiran, sqlDaftarLampiran)
	if err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, produkID)
	if err != nil {
		return nil, fmt.Errorf("repository: reading %s: %w", TabelLampiran, err)
	}
	hasil := []models.Lampiran{}
	for rows.Next() {
		l, err := pindaiLampiran(rows)
		if err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("repository: reading %s: %w", TabelLampiran, err)
		}
		hasil = append(hasil, l)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	for i := range hasil {
		if err := g.lengkapiStatus(ctx, nil, &hasil[i]); err != nil {
			return nil, err
		}
	}
	return hasil, nil
}

// AmbilLampiran - satu lampiran produk, beserta status.
func (g *Gudang) AmbilLampiran(ctx context.Context, tx *db.Tx, produkID, id string) (models.Lampiran, error) {
	q, err := g.siapkan(TabelLampiran, sqlAmbilLampiran)
	if err != nil {
		return models.Lampiran{}, err
	}
	l, err := pindaiLampiran(g.kueri(tx).QueryRowContext(ctx, q, produkID, id))
	if errors.Is(err, sql.ErrNoRows) {
		return models.Lampiran{}, fmt.Errorf("%w: %s", ErrLampiranTidakAda, id)
	}
	if err != nil {
		return models.Lampiran{}, fmt.Errorf("repository: reading %s %s: %w", TabelLampiran, id, err)
	}
	return l, g.lengkapiStatus(ctx, tx, &l)
}

// SisipLampiran merekam lampiran baru; ID = stempel waktu basis data (FF3).
func (g *Gudang) SisipLampiran(ctx context.Context, tx *db.Tx, l models.Lampiran) (string, error) {
	if !tx.Terisi() {
		return "", errors.New("repository: recording an attachment requires a transaction")
	}
	qAda, err := g.siapkan(TabelLampiran, func(t string) string {
		return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ID = :1`, t)
	})
	if err != nil {
		return "", err
	}
	for coba := 0; coba < 5; coba++ {
		var id string
		if err := tx.QueryRowContext(ctx, sqlWaktuID()).Scan(&id); err != nil {
			return "", fmt.Errorf("repository: attachment ID: %w", err)
		}
		var n int
		if err := tx.QueryRowContext(ctx, qAda, id).Scan(&n); err != nil {
			return "", fmt.Errorf("repository: attachment ID: %w", err)
		}
		if n > 0 {
			time.Sleep(2 * time.Millisecond)
			continue
		}
		err := g.exec(ctx, tx, TabelLampiran, sqlSisipLampiran, id, l.ProdukID, l.Category, l.FileName,
			l.FileMimeType, l.UserName, l.StorageID)
		return id, err
	}
	return "", ErrIDLampiranBentrok
}

// HapusLampiran - `DeleteAttachProdName_Sql` b85.
func (g *Gudang) HapusLampiran(ctx context.Context, tx *db.Tx, produkID, id string) error {
	return g.exec(ctx, tx, TabelLampiran, sqlHapusLampiran, produkID, id)
}

// AdaObjek - objek `T_STORAGE_IMAGE` untuk IMAGEID tercatat.
func (g *Gudang) AdaObjek(ctx context.Context, tx *db.Tx, imageID string) (bool, error) {
	if imageID == "" {
		return false, nil
	}
	q, err := g.siapkan(TabelObjek, sqlAdaObjek)
	if err != nil {
		return false, err
	}
	var n int
	if err := g.kueri(tx).QueryRowContext(ctx, q, imageID).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: reading %s: %w", TabelObjek, err)
	}
	return n > 0, nil
}

// CatatObjek - `Insert_T_Storage_SQL`; idempoten (objek yang sudah tercatat tidak digandakan).
func (g *Gudang) CatatObjek(ctx context.Context, tx *db.Tx, o models.ObjekPenyimpanan) error {
	ada, err := g.AdaObjek(ctx, tx, o.ImageID)
	if err != nil || ada {
		return err
	}
	return g.exec(ctx, tx, TabelObjek, sqlSisipObjek, o.ImageID, o.AppFolder, o.DurasiDetik, o.FileName,
		db.KosongJadiNil(o.AppName))
}

// HapusObjek - `DeleteStorage_SQL`; objek yang sudah tidak ada bukan galat.
func (g *Gudang) HapusObjek(ctx context.Context, tx *db.Tx, imageID string) error {
	if !tx.Terisi() {
		return errors.New("repository: deleting a storage object requires a transaction")
	}
	q, err := g.siapkan(TabelObjek, sqlHapusObjek)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, imageID); err != nil {
		return fmt.Errorf("repository: deleting %s %s: %w", TabelObjek, imageID, err)
	}
	return nil
}

// NamaAplikasi - `GetAppName_SQL` b58 (`InsertGoogleStorage_Act` 6 b958, PRE=false → selalu).
func (g *Gudang) NamaAplikasi(ctx context.Context, tx *db.Tx) (string, error) {
	q, err := g.siapkan(TabelFolder, sqlNamaAplikasi)
	if err != nil {
		return "", err
	}
	var app sql.NullString
	err = g.kueri(tx).QueryRowContext(ctx, q).Scan(&app)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: reading %s: %w", TabelFolder, err)
	}
	return app.String, nil
}

// AntreUnggah - satu efek pengiriman berkas, di transaksi pemanggil.
func (g *Gudang) AntreUnggah(ctx context.Context, tx *db.Tx, lampiranID, muatan string, saat time.Time) error {
	_, err := outbox.NewPenyimpan(g.db).AntreEfek(ctx, tx, outbox.LiniLife, ModulOutbox, JenisEfekUnggah, lampiranID,
		muatan, saat)
	return err
}

// AdaUnggahAntre - lampiran punya efek berstatus antre (kirim ulang memakainya, bukan menambah).
func (g *Gudang) AdaUnggahAntre(ctx context.Context, tx *db.Tx, lampiranID string) (bool, error) {
	q, err := g.siapkan(TabelOutbox, sqlAdaUnggahAntre)
	if err != nil {
		return false, err
	}
	var n int
	if err := g.kueri(tx).QueryRowContext(ctx, q, ModulOutbox, lampiranID, outbox.StatusEfekAntre).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: reading attachment %s send effects: %w", lampiranID, err)
	}
	return n > 0, nil
}

// PungutUnggah - efek antre milik SATU lampiran (dikunci SKIP LOCKED), ditandai jalan.
func (g *Gudang) PungutUnggah(ctx context.Context, tx *db.Tx, lampiranID string, saat time.Time) (string, int, bool, error) {
	if !tx.Terisi() {
		return "", 0, false, errors.New("repository: taking a send effect requires a transaction")
	}
	q, err := g.siapkan(TabelOutbox, sqlPungutUnggah)
	if err != nil {
		return "", 0, false, err
	}
	var id string
	var percobaan sql.NullInt64
	rows, err := tx.QueryContext(ctx, q, ModulOutbox, lampiranID, outbox.StatusEfekAntre)
	if err != nil {
		return "", 0, false, fmt.Errorf("repository: taking attachment %s send effect: %w", lampiranID, err)
	}
	ada := rows.Next()
	if ada {
		err = rows.Scan(&id, &percobaan)
	}
	if tutup := rows.Close(); err == nil {
		err = tutup
	}
	if err != nil {
		return "", 0, false, fmt.Errorf("repository: taking attachment %s send effect: %w", lampiranID, err)
	}
	if !ada {
		return "", 0, false, nil
	}
	tabel, err := g.db.Qualify(TabelOutbox)
	if err != nil {
		return "", 0, false, err
	}
	tandai := outbox.SQLTandaiJalan(tabel)
	if err := db.PeriksaSQL(tandai); err != nil {
		return "", 0, false, err
	}
	hasil, err := tx.ExecContext(ctx, tandai, outbox.StatusEfekJalan, saat, id, outbox.StatusEfekAntre)
	if err != nil {
		return "", 0, false, fmt.Errorf("repository: marking send effect %s: %w", id, err)
	}
	if err := db.PastikanSatuBaris(hasil, "marking send effect"); err != nil {
		return "", 0, false, err
	}
	return id, int(percobaan.Int64) + 1, true, nil
}

// TuntaskanUnggah - hasil satu percobaan pengiriman.
func (g *Gudang) TuntaskanUnggah(ctx context.Context, tx *db.Tx, efekID, status string, jadwal time.Time, galat string,
	saat time.Time) error {
	return outbox.NewPenyimpan(g.db).TuntaskanEfek(ctx, tx, efekID, status, jadwal, galat, saat)
}
