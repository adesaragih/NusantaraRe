package repository

// Tulisan ke `T_CLAIMLF_DOCUMENT` dan `T_CLAIMLF_STORAGE` - butir be.
//
// Untuk apa berkas ini: yang MENULIS dokumen. Yang MEMBACA ada di
// `klaimlife.go` (`AmbilDokumen`), dan sengaja tidak dipindah: pembacaan
// dokumen ikut pembacaan klaim, penulisannya berdiri sendiri.
//
// ⛔ Kedua tabel ditulis DI DALAM transaksi pemanggil. Baris dokumen tanpa
// kartu penyimpanannya, atau sebaliknya, adalah keadaan yang layar tidak
// dapat jelaskan.
//
// Dibaca sesudah: klaimlife.go dan efekkeluar.go (outbox-nya).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimlife/models"
)

// ErrDokumenTidakAda - baris dokumen yang diminta bukan milik klaim itu.
var ErrDokumenTidakAda = errors.New("repository: dokumen tidak ada pada klaim ini")

// sqlSisipDokumen merakit penulisan satu baris dokumen.
//
// ⛔ Urutan kolomnya mengikuti `InsertDocument_Act.xml` langkah 3 b627:
// `.ID` b647 (cap waktu), `.TANGGAL` b701 (`@CurrentDateTime()`), `.MIME`
// b781 (`@toLowerCase`), `KATEGORI_1` = kunci kelompok `DL-` b595,
// `KATEGORI_2` = kategori yang dipilih pemakai.
//
// ⚠️ `T_STORAGE_ID` ditulis KOSONG lebih dulu - lihat `SisipDokumen`.
func sqlSisipDokumen(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s
		(ID, PREMIUM_LIST_DETAIL_ID, NAMA_FILE, MIME, KATEGORI_1, KATEGORI_2,
		 TANGGAL, T_STORAGE_ID)
		VALUES (:1,:2,:3,:4,:5,:6,:7,:8)`, tabel)
}

// SisipDokumen menulis satu baris dokumen.
//
// ⛔ PENYIMPANGAN SADAR ATAS URUTAN PEGA, dan arahnya dinyatakan.
// `InsertDocument_Act` langkah 5 b1185 berprasyarat b1283
// `NewDocument.T_STORAGE_ID==""` dengan `WhenTrue=3` - artinya baris TIDAK
// disimpan selama penyimpanan belum menjawab; `models.BolehSimpanBarisDokumen`
// menirukan gerbang itu. Di Pega hal itu MUNGKIN sebab langkah 4 b1023
// memanggil penyimpanan SEREMPAK, sehingga saat langkah 5 tiba jawabannya
// sudah ada.
//
// Dengan outbox jawabannya BELUM ada: efeknya dijalankan sesudah transaksi
// ini. Maka dua pilihan, dan keduanya berbiaya:
//
//	a. tulis baris hanya sesudah efek selesai - setia pada b1283, tetapi
//	   efek yang gagal membuat dokumennya HILANG dan tidak ada apa pun untuk
//	   dicoba ulang;
//	b. tulis baris sekarang dengan `T_STORAGE_ID` kosong - baris itulah yang
//	   outbox rujuk, dan efeknya mengisi kolom itu.
//
// Dipilih (b), sesuai butir **be**. Yang hilang: gerbang b1283 tidak berlaku
// di jalur ini. Yang didapat: berkas yang sudah mendarat tidak pernah tanpa
// catatan. Layar sudah lebih dulu menandai baris ber-`T_STORAGE_ID` kosong
// sebagai "belum terunggah" - penyimpangan yang dicatat sejak `AmbilDokumen`.
func (r *KlaimLife) SisipDokumen(ctx context.Context, tx *db.Tx,
	d models.Dokumen, saat time.Time) error {

	tabel, err := r.db.Qualify("T_CLAIMLF_DOCUMENT")
	if err != nil {
		return err
	}
	q := sqlSisipDokumen(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, d.ID, d.PesertaID,
		db.KosongJadiNil(d.NamaFile), db.KosongJadiNil(d.Mime),
		db.KosongJadiNil(d.Kategori1), db.KosongJadiNil(d.Kategori2),
		saat, db.KosongJadiNil(d.TStorageID))
	if err != nil {
		return fmt.Errorf("repository: menyisip dokumen: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "penyisipan dokumen")
}

// SatuDokumen membaca satu baris dokumen milik sebuah klaim.
//
// ⛔ Dibatasi KLAIM-nya lewat pesertanya, bukan dibaca dengan `ID` saja.
// Pengenal dokumen datang dari jalur URL; tanpa batas itu, `GET
// /api/dokumen/9/isi` menyerahkan berkas milik klaim siapa pun.
func (r *KlaimLife) SatuDokumen(ctx context.Context, klaimID string, id int64) (
	models.Dokumen, error) {

	dok, err := r.db.Qualify("T_CLAIMLF_DOCUMENT")
	if err != nil {
		return models.Dokumen{}, err
	}
	pes, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return models.Dokumen{}, err
	}
	q := fmt.Sprintf(
		`SELECT d.PREMIUM_LIST_DETAIL_ID, d.ID, d.NAMA_FILE, d.MIME,
		        d.KATEGORI_1, d.KATEGORI_2, d.T_STORAGE_ID
		   FROM %s d JOIN %s p ON p.ID = d.PREMIUM_LIST_DETAIL_ID
		  WHERE p.CLAIM_ID = :1 AND d.ID = :2 AND p.STS_HAPUS IS NULL`, dok, pes)
	if err := db.PeriksaSQL(q); err != nil {
		return models.Dokumen{}, err
	}
	var (
		pesertaID, nama, mime, kat1, kat2, storage sql.NullString
		baca                                       sql.NullInt64
	)
	err = r.db.QueryRowContext(ctx, q, klaimID, id).Scan(&pesertaID, &baca,
		&nama, &mime, &kat1, &kat2, &storage)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Dokumen{}, ErrDokumenTidakAda
	}
	if err != nil {
		return models.Dokumen{}, fmt.Errorf("repository: membaca dokumen: %w", err)
	}
	return models.Dokumen{
		ID: baca.Int64, PesertaID: pesertaID.String, NamaFile: nama.String,
		Mime: mime.String, Kategori1: kat1.String, Kategori2: kat2.String,
		TStorageID: storage.String,
	}, nil
}

// HapusDokumen membuang satu baris dokumen.
//
// `[terverifikasi]` `DeleteDocument_Act.xml` b513 `Obj-Delete` berjalan
// TANPA prasyarat - barisnya dihapus apa pun keadaan penyimpanannya. Yang
// berprasyarat hanya panggilan ke penyimpanan (b472, `WhenTrue=3` LEWATI
// bila `T_STORAGE_ID` kosong), dan itu tinggal di `models.PerluHapusDiPenyimpanan`.
func (r *KlaimLife) HapusDokumen(ctx context.Context, tx *db.Tx, id int64) error {
	tabel, err := r.db.Qualify("T_CLAIMLF_DOCUMENT")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`DELETE FROM %s WHERE ID = :1`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("repository: menghapus dokumen: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "penghapusan dokumen")
}

// sqlSisipKartuBerkas merakit penulisan kartu penyimpanan.
//
// Meniru `RDBList/Insert_T_Storage_SQL.xml` b86-100, TANPA `commit;` b102
// yang ada di sana - transaksi milik pemanggil (ADR-U-0029).
func sqlSisipKartuBerkas(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s
		(IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE,
		 TANGGAL_UPLOAD)
		VALUES (:1,:2,:3,:4,:5,:6,:7,:8)`, tabel)
}

// KartuBerkas adalah satu baris `T_CLAIMLF_STORAGE`.
type KartuBerkas struct {
	ImageID   string
	URLPublic string
	AppFolder string
	ExpDate   time.Time
	FileName  string
	AppName   string
	Storage   string
	// TanggalUpload - `Update_T_Storage_SQL.xml` b89.
	//
	// ⚠️ Sumbernya BUKAN jam kita: `GetUrlGoogleStorage_Act.xml`
	// b2273-2274 menyetelnya dari `UploadDoc.Response.DateTime`, yaitu cap
	// waktu yang DIKEMBALIKAN layanan penyimpanan. Selama pelaksananya stub,
	// layanan itu proses kita sendiri - penyimpangan yang dinyatakan.
	TanggalUpload time.Time
}

// SisipKartuBerkas menulis kartu penyimpanan satu berkas.
func (r *KlaimLife) SisipKartuBerkas(ctx context.Context, tx *db.Tx, k KartuBerkas) error {
	tabel, err := r.db.Qualify("T_CLAIMLF_STORAGE")
	if err != nil {
		return err
	}
	q := sqlSisipKartuBerkas(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, k.ImageID, db.KosongJadiNil(k.URLPublic),
		db.KosongJadiNil(k.AppFolder), waktuJadiNil(k.ExpDate),
		db.KosongJadiNil(k.FileName), db.KosongJadiNil(k.AppName),
		db.KosongJadiNil(k.Storage), waktuJadiNil(k.TanggalUpload))
	if err != nil {
		return fmt.Errorf("repository: menyisip kartu berkas: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "penyisipan kartu berkas")
}

// HapusKartuBerkas membuang kartu penyimpanan satu berkas.
//
// ⚠️ Nol baris terpengaruh SAH: dokumen yang efek unggahnya belum selesai
// belum punya kartu. `pastikanSatuBaris` karena itu tidak dipakai.
func (r *KlaimLife) HapusKartuBerkas(ctx context.Context, tx *db.Tx, imageID string) error {
	tabel, err := r.db.Qualify("T_CLAIMLF_STORAGE")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`DELETE FROM %s WHERE IMAGEID = :1`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, imageID); err != nil {
		return fmt.Errorf("repository: menghapus kartu berkas: %w", err)
	}
	return nil
}

// sqlTautanBerkas merakit pembacaan URL publik sebuah kartu.
//
// Meniru `RDBList/GetLinkStorage_SQL.xml` b85-91: ia memilih `URLPUBLIC`
// beserta empat kolom lain dengan `where imageid = {UploadDoc.ImageID}`.
func sqlTautanBerkas(tabel string) string {
	return fmt.Sprintf(`SELECT URLPUBLIC FROM %s WHERE IMAGEID = :1`, tabel)
}

// TautanBerkas membaca URL publik sebuah kartu penyimpanan.
//
// Mengembalikan teks kosong bila kartunya belum ada - keadaan NORMAL selama
// efek unggahnya belum dijalankan, bukan kegagalan.
func (r *KlaimLife) TautanBerkas(ctx context.Context, imageID string) (string, error) {
	tabel, err := r.db.Qualify("T_CLAIMLF_STORAGE")
	if err != nil {
		return "", err
	}
	q := sqlTautanBerkas(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var url sql.NullString
	err = r.db.QueryRowContext(ctx, q, imageID).Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca tautan berkas: %w", err)
	}
	return url.String, nil
}

// TandaiDokumenTerunggah mengisi `T_STORAGE_ID` sesudah efek unggah selesai.
//
// ⛔ Syarat WHERE menyertakan `T_STORAGE_ID IS NULL`: efek yang terlanjur
// dijalankan dua kali tidak menimpa kartu yang sudah tertaut. Outbox menjamin
// SETIDAKNYA sekali, bukan tepat sekali.
func (r *KlaimLife) TandaiDokumenTerunggah(ctx context.Context, tx *db.Tx,
	id int64, imageID string) error {

	tabel, err := r.db.Qualify("T_CLAIMLF_DOCUMENT")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`UPDATE %s SET T_STORAGE_ID = :1
		 WHERE ID = :2 AND T_STORAGE_ID IS NULL`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, imageID, id); err != nil {
		return fmt.Errorf("repository: menandai dokumen terunggah: %w", err)
	}
	return nil
}

// KlaimDokumen mencari klaim pemilik sebuah dokumen.
//
// ⛔ Ia ada supaya rute `GET /api/dokumen/{id}/isi` - yang jalurnya TIDAK
// menyebut klaim, sebab `URLPUBLIC` di sistem lama pun tidak - tetap dapat
// memakai pembacaan yang SAMA dengan rute lain, yaitu `SatuDokumen` yang
// berbatas klaim. Tanpa ini rute itu akan punya pembacaan sendiri, dan
// pembacaan kedua adalah tempat kedua bagi batas untuk hilang.
func (r *KlaimLife) KlaimDokumen(ctx context.Context, id int64) (string, error) {
	dok, err := r.db.Qualify("T_CLAIMLF_DOCUMENT")
	if err != nil {
		return "", err
	}
	pes, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf(
		`SELECT p.CLAIM_ID FROM %s d JOIN %s p ON p.ID = d.PREMIUM_LIST_DETAIL_ID
		  WHERE d.ID = :1 AND p.STS_HAPUS IS NULL`, dok, pes)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var klaim sql.NullString
	err = r.db.QueryRowContext(ctx, q, id).Scan(&klaim)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrDokumenTidakAda
	}
	if err != nil {
		return "", fmt.Errorf("repository: mencari klaim dokumen: %w", err)
	}
	return klaim.String, nil
}
