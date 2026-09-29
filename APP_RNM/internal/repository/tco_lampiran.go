package repository

// Lampiran tahun treaty - tiket 12 Treaty Contract Out, tco4: tabel WARISAN.
//
// Untuk apa berkas ini: rekam lampiran di `M_ATTACHMENTTREATY_2`, objek berkas
// di `T_STORAGE_IMAGE`, dan master kategori `CATEGORY_ATTACH_REAS` (dibaca
// saja) - persis RDB warisannya:
//
//	GetAllAttachment2_Sql.xml b84 / GetAttachment2_Sql.xml b85   baca, `where treatyid = {TreatyIn.ID}`
//	DeleteAttachment2_Sql.xml b84    `delete M_ATTACHMENTTREATY_2 where treatyid = .. and id = ..`
//	InsertAtatchment_Sql.xml b60     `PEGA_M_ATTACHMENT(IDPEGA, DATAPEGA)` - badan [terbuka - DBA],
//	                                 kolom ditiru dari penulis langsung tabel yang SAMA:
//	                                 `Treaty In/RDBList/InsertAttachment2_Sql.xml` b84 (OQ-TCO-24)
//	Insert_T_Storage_SQL.xml b85 (Claim Fac In), DeleteStorage_SQL.xml b85   `T_STORAGE_IMAGE`
//
// ⛔ Kunci pemilik `TREATYID` = `TreatyYear + TreatyYearID` (teks disambung):
// `TreatyOutSaveAttachment` b1402 dan `DeleteAttachmentTreaty` b252 mengisi
// `TreatyIn.ID` begitu [terverifikasi]. RALAT tiket 12 ("kunci treaty inward"):
// nama halamannya `TreatyIn`, isinya milik modul ini.
//
// ⛔ Status "terkirim" = baris `T_STORAGE_IMAGE` untuk `T_STORAGE_ID` ada;
// "gagal" dibaca dari outbox bersama `T_LOG_SERVICE_RNM` (milik aplikasi) -
// disaring `MODUL` modul ini.
//
// Dibaca sesudah: tco_tahun.go, efekkeluar.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"nusantarare/internal/models"
)

// MasterKategoriLampiranTCO - master kategori lampiran, dibaca saja.
//
// `[terverifikasi]` `Treaty Contract Out/RDBList/CategoryAttach_SQL.xml`
// `select * from CATEGORY_ATTACH_REAS order by note`.
const MasterKategoriLampiranTCO = "CATEGORY_ATTACH_REAS"

// TabelOutboxBersama - outbox efek keluar lintas modul (migrasi 015).
const TabelOutboxBersama = "T_LOG_SERVICE_RNM"

// Tabel dan nilai WARISAN lampiran (tco4).
const (
	// TabelLampiranTCO - `M_ATTACHMENTTREATY_2`.
	TabelLampiranTCO = "M_ATTACHMENTTREATY_2"
	// TabelObjekPenyimpananTCO - `T_STORAGE_IMAGE`.
	TabelObjekPenyimpananTCO = "T_STORAGE_IMAGE"
	// KategoriPegaLampiranTCO - `CATEGORY` = `.pyCategory` = "File"
	// (`TreatyOutSaveAttachment` b558); kategori PILIHAN di `CATEGORY_ID`.
	KategoriPegaLampiranTCO = "File"
)

// ErrLampiranTidakAda - lampiran tidak ada, atau bukan milik tahun treaty itu.
var ErrLampiranTidakAda = errors.New("repository: lampiran tidak ditemukan pada tahun treaty ini")

// BarisLampiranTCO adalah satu lampiran beserta nasib efek unggah TERAKHIRnya.
type BarisLampiranTCO struct {
	models.LampiranTCO
	// StatusEfek - status efek `storage-unggah` terakhir; kosong bila belum ada.
	StatusEfek string
	// GalatEfek - `GALAT_TERAKHIR` efek itu.
	GalatEfek string
	// PercobaanEfek - cacah percobaan efek itu.
	PercobaanEfek int
}

// KategoriLampiran membaca master kategori lampiran.
type KategoriLampiran struct{ db *DB }

// NewKategoriLampiran menyusun pembacanya.
func NewKategoriLampiran(db *DB) *KategoriLampiran { return &KategoriLampiran{db: db} }

// sqlKategoriLampiranTCO - hanya `NOTE`, satu-satunya kolom yang terbukti
// dipakai (`SetCategoryAttachTreatyin.xml` b500 `.NOTE`).
//
// ⚠️ `DISTINCT` adalah tambahan kami: nilai pilihan layar adalah teksnya, dan
// dua baris master senama akan menjadi dua pilihan yang tidak dapat dibedakan.
func sqlKategoriLampiranTCO(tabel string) string {
	return fmt.Sprintf(`SELECT DISTINCT NOTE FROM %s WHERE NOTE IS NOT NULL ORDER BY NOTE`, tabel)
}

// Daftar membaca seluruh kategori, urut `NOTE`.
func (k *KategoriLampiran) Daftar(ctx context.Context) ([]string, error) {
	tabel, err := k.db.Qualify(MasterKategoriLampiranTCO)
	if err != nil {
		return nil, err
	}
	q := sqlKategoriLampiranTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := k.db.sql.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca master %s: %w", MasterKategoriLampiranTCO, err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []string
	for rows.Next() {
		var n sql.NullString
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		if strings.TrimSpace(n.String) != "" {
			hasil = append(hasil, n.String)
		}
	}
	return hasil, rows.Err()
}

// MasterLampiranTCO membaca dan menulis `M_ATTACHMENTTREATY_2` + `T_STORAGE_IMAGE`.
type MasterLampiranTCO struct{ db *DB }

// NewMasterLampiranTCO menyusun gudangnya.
func NewMasterLampiranTCO(db *DB) *MasterLampiranTCO { return &MasterLampiranTCO{db: db} }

// IDLampiranTCO - `TO_CHAR(SYSTIMESTAMP, 'YYYYMMDDHH24MISSFF3')` (Treaty In
// `InsertAttachment2_Sql` b84), jam WIB. Ia sekaligus cap waktu unggahnya.
func IDLampiranTCO(saat time.Time) string {
	w := saat.In(zonaJakartaTCO)
	return w.Format("20060102150405") + fmt.Sprintf("%03d", w.Nanosecond()/int(time.Millisecond))
}

// waktuDariIDLampiranTCO - kebalikan `IDLampiranTCO`; ID berbentuk lain -> nol.
func waktuDariIDLampiranTCO(id string) time.Time {
	if len(id) != 17 {
		return time.Time{}
	}
	w, err := time.ParseInLocation("20060102150405", id[:14], zonaJakartaTCO)
	if err != nil {
		return time.Time{}
	}
	ms, err := strconv.Atoi(id[14:])
	if err != nil {
		return time.Time{}
	}
	return w.Add(time.Duration(ms) * time.Millisecond)
}

// pilihLampiranTCO - kolom GetAllAttachment2_Sql + keberadaan objeknya.
//
// ⚠️ Subkueri skalar, bukan JOIN: `T_STORAGE_IMAGE` tanpa PK `[terbuka - DBA]`,
// dan JOIN atas baris objek kembar menggandakan lampiran.
func pilihLampiranTCO(objek string) string {
	return fmt.Sprintf(`a.ID, a.FILENAME, a.FILEMIMETYPE, a.CATEGORY_ID, a.T_STORAGE_ID,
	       (SELECT MAX(s.IMAGEID) FROM %s s WHERE s.IMAGEID = a.T_STORAGE_ID), a.USERNAME`, objek)
}

// sqlDaftarLampiranTCO membaca lampiran satu tahun + efek unggah terakhirnya.
//
// ⛔ Efek TERAKHIR per lampiran (`ROW_NUMBER ... DIBUAT DESC`): pengulangan
// menambah baris outbox baru, dan status yang tampil adalah nasib percobaan
// terbaru, bukan yang pertama.
func sqlDaftarLampiranTCO(tabel, objek, outbox string, satu bool) string {
	saring := ""
	if satu {
		saring = " AND a.ID = :4"
	}
	return fmt.Sprintf(`SELECT %s,
	       o.STATUS, o.GALAT_TERAKHIR, o.PERCOBAAN
	  FROM %s a
	  LEFT JOIN (SELECT RUJUKAN, STATUS, GALAT_TERAKHIR, PERCOBAAN,
	                    ROW_NUMBER() OVER (PARTITION BY RUJUKAN ORDER BY DIBUAT DESC, ID DESC) AS URUT
	               FROM %s
	              WHERE MODUL = :1 AND JENIS_EFEK = :2) o
	    ON o.RUJUKAN = a.ID AND o.URUT = 1
	 WHERE a.TREATYID = :3%s
	 ORDER BY a.ID DESC`, pilihLampiranTCO(objek), tabel, outbox, saring)
}

func sqlKunciLampiranTCO(tabel string) string {
	return fmt.Sprintf(`SELECT ID FROM %s WHERE ID = :1 FOR UPDATE`, tabel)
}

func sqlAmbilUntukKirimLampiranTCO(tabel, objek string) string {
	return fmt.Sprintf(`SELECT %s FROM %s a WHERE a.ID = :1`, pilihLampiranTCO(objek), tabel)
}

func sqlAdaIDLampiranTCO(tabel string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ID = :1`, tabel)
}

// sqlSisipLampiranTCO - kolom `InsertAttachment2_Sql` b84; `DATA_JSON` NULL
// (nol kolom dokumen - atribut lampiran berkolom bernama).
func sqlSisipLampiranTCO(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s
	       (ID, TREATYID, CATEGORY, FILENAME, FILEMIMETYPE, DATA_JSON, USERNAME, CATEGORY_ID, T_STORAGE_ID)
	VALUES (:1, :2, :3, :4, :5, NULL, :6, :7, :8)`, tabel)
}

// sqlHapusLampiranTCO - `DeleteAttachment2_Sql` b84.
func sqlHapusLampiranTCO(tabel string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE TREATYID = :1 AND ID = :2`, tabel)
}

// sqlSimpanObjekTCO - `Insert_T_Storage_SQL` (Claim Fac In) b85: `EXPDATE`
// `To_date(exp, 'DD/MM/YYYY HH24:MI:SS')`, `STORAGE` 'standard'.
func sqlSimpanObjekTCO(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s
	       (IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE)
	VALUES (:1, :2, :3, TO_DATE(:4, 'DD/MM/YYYY HH24:MI:SS'), :5, :6, 'standard')`, tabel)
}

// sqlHapusObjekTCO - `DeleteStorage_SQL` b85.
func sqlHapusObjekTCO(tabel string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE IMAGEID = :1`, tabel)
}

func sqlTreatyYearLampiranTCO(tabel string) string {
	return fmt.Sprintf(`SELECT TREATYYEAR FROM %s WHERE ID = :1`, tabel)
}

// kunciTreaty - `TREATYID` lampiran satu tahun treaty: `TreatyYear + ID`.
func (m *MasterLampiranTCO) kunciTreaty(ctx context.Context, q kuerierTCO, tahunID string) (string, error) {
	tabel, err := m.db.Qualify(TabelTahunTCO)
	if err != nil {
		return "", err
	}
	sq := sqlTreatyYearLampiranTCO(tabel)
	if err := PeriksaSQL(sq); err != nil {
		return "", err
	}
	var th sql.NullString
	err = q.QueryRowContext(ctx, sq, tahunID).Scan(&th)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrTahunTreatyTidakAda
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca tahun treaty lampiran %s: %w", tahunID, err)
	}
	return models.KunciTreatyLampiranTCO(th.String, tahunID), nil
}

// pindaiLampiranTCO membaca tujuh kolom `pilihLampiranTCO`.
func pindaiLampiranTCO(n [7]sql.NullString, tahunID string) models.LampiranTCO {
	return models.LampiranTCO{
		ID: n[0].String, IDTreatyYear: tahunID, FileName: n[1].String, FileMimeType: n[2].String,
		Category: n[3].String, ImageID: n[4].String, TStorageID: n[5].String, UserID: n[6].String,
		TglUpload: waktuDariIDLampiranTCO(n[0].String),
	}
}

// Daftar membaca lampiran satu tahun treaty, ID DESC.
func (m *MasterLampiranTCO) Daftar(ctx context.Context, modul, jenisUnggah, tahunID string) (
	[]BarisLampiranTCO, error) {
	return m.baca(ctx, modul, jenisUnggah, tahunID, "")
}

// Ambil membaca SATU lampiran, dibatasi tahun treaty-nya.
//
// ⛔ Dibatasi tahun (`TREATYID`), bukan dicari dengan ID saja: ID lampiran
// datang dari jalur URL, dan tanpa batas itu jalur tahun mana pun menyerahkan
// lampiran tahun lain.
func (m *MasterLampiranTCO) Ambil(ctx context.Context, modul, jenisUnggah, tahunID, id string) (
	BarisLampiranTCO, error) {
	baris, err := m.baca(ctx, modul, jenisUnggah, tahunID, id)
	if err != nil {
		return BarisLampiranTCO{}, err
	}
	if len(baris) == 0 {
		return BarisLampiranTCO{}, ErrLampiranTidakAda
	}
	return baris[0], nil
}

func (m *MasterLampiranTCO) tabelLampiran() (string, string, error) {
	tabel, err := m.db.Qualify(TabelLampiranTCO)
	if err != nil {
		return "", "", err
	}
	objek, err := m.db.Qualify(TabelObjekPenyimpananTCO)
	if err != nil {
		return "", "", err
	}
	return tabel, objek, nil
}

func (m *MasterLampiranTCO) baca(ctx context.Context, modul, jenisUnggah, tahunID, id string) (
	[]BarisLampiranTCO, error) {
	tabel, objek, err := m.tabelLampiran()
	if err != nil {
		return nil, err
	}
	outbox, err := m.db.Qualify(TabelOutboxBersama)
	if err != nil {
		return nil, err
	}
	kunci, err := m.kunciTreaty(ctx, m.db.bacaTCO(ctx), tahunID)
	if err != nil {
		return nil, err
	}
	q := sqlDaftarLampiranTCO(tabel, objek, outbox, id != "")
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	arg := []any{modul, jenisUnggah, kunci}
	if id != "" {
		arg = append(arg, id)
	}
	rows, err := m.db.bacaTCO(ctx).QueryContext(ctx, q, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca lampiran tahun treaty %s: %w", tahunID, err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []BarisLampiranTCO
	for rows.Next() {
		var n [7]sql.NullString
		var status, galat sql.NullString
		var percobaan sql.NullInt64
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6], &status, &galat, &percobaan); err != nil {
			return nil, err
		}
		hasil = append(hasil, BarisLampiranTCO{LampiranTCO: pindaiLampiranTCO(n, tahunID), StatusEfek: status.String,
			GalatEfek: galat.String, PercobaanEfek: int(percobaan.Int64)})
	}
	return hasil, rows.Err()
}

// AmbilUntukKirim mengunci satu lampiran untuk pelaksana efek.
//
// ⛔ `FOR UPDATE`: dua pelaksana atas lampiran yang sama menunggu satu sama
// lain, dan yang kedua melihat objek `T_STORAGE_IMAGE` yang sudah tercatat.
// Kunci dan bacaan dua pernyataan: subkueri skalar tidak bercampur FOR UPDATE.
func (m *MasterLampiranTCO) AmbilUntukKirim(ctx context.Context, tx *Tx, id string) (models.LampiranTCO, error) {
	if tx == nil {
		return models.LampiranTCO{}, errors.New("repository: mengunci lampiran menuntut transaksi")
	}
	tabel, objek, err := m.tabelLampiran()
	if err != nil {
		return models.LampiranTCO{}, err
	}
	qk := sqlKunciLampiranTCO(tabel)
	if err := PeriksaSQL(qk); err != nil {
		return models.LampiranTCO{}, err
	}
	var terkunci sql.NullString
	if err := tx.tx.QueryRowContext(ctx, qk, id).Scan(&terkunci); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.LampiranTCO{}, ErrLampiranTidakAda
		}
		return models.LampiranTCO{}, fmt.Errorf("repository: mengunci lampiran %s: %w", id, err)
	}
	q := sqlAmbilUntukKirimLampiranTCO(tabel, objek)
	if err := PeriksaSQL(q); err != nil {
		return models.LampiranTCO{}, err
	}
	var n [7]sql.NullString
	if err := tx.tx.QueryRowContext(ctx, q, id).Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6]); err != nil {
		return models.LampiranTCO{}, fmt.Errorf("repository: membaca lampiran %s: %w", id, err)
	}
	return pindaiLampiranTCO(n, ""), nil
}

// Sisip menulis satu lampiran; ID = stempel waktu unggahnya (`IDLampiranTCO`).
//
// ⚠️ `M_ATTACHMENTTREATY_2` tanpa PK yang diketahui: ID yang sudah terpakai
// (dua unggahan pada milidetik yang sama) digeser satu milidetik
// `[keputusan kami]` - dua lampiran tidak boleh berbagi ID.
func (m *MasterLampiranTCO) Sisip(ctx context.Context, tx *Tx, l models.LampiranTCO) (string, error) {
	if tx == nil {
		return "", errors.New("repository: menyisipkan lampiran menuntut transaksi")
	}
	tabel, _, err := m.tabelLampiran()
	if err != nil {
		return "", err
	}
	kunci, err := m.kunciTreaty(ctx, tx.tx, l.IDTreatyYear)
	if err != nil {
		return "", err
	}
	qa := sqlAdaIDLampiranTCO(tabel)
	if err := PeriksaSQL(qa); err != nil {
		return "", err
	}
	saat := l.TglUpload
	var id string
	for i := 0; ; i++ {
		id = IDLampiranTCO(saat)
		var n int
		if err := tx.tx.QueryRowContext(ctx, qa, id).Scan(&n); err != nil {
			return "", fmt.Errorf("repository: memeriksa ID lampiran: %w", err)
		}
		if n == 0 {
			break
		}
		if i >= 1000 {
			return "", fmt.Errorf("repository: ID lampiran %s dan seribu penggantinya sudah terpakai", id)
		}
		saat = saat.Add(time.Millisecond)
	}
	q := sqlSisipLampiranTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, id, kunci, KategoriPegaLampiranTCO, kosongJadiNil(l.FileName),
		kosongJadiNil(l.FileMimeType), kosongJadiNil(l.UserID), kosongJadiNil(l.Category), l.ImageID)
	if err != nil {
		return "", fmt.Errorf("repository: menyisipkan lampiran: %w", err)
	}
	return id, pastikanSatuBaris(hasil, "penyisipan lampiran")
}

// Hapus membuang satu lampiran milik tahun treaty itu (`DeleteAttachment2_Sql`).
func (m *MasterLampiranTCO) Hapus(ctx context.Context, tx *Tx, tahunID, id string) error {
	if tx == nil {
		return errors.New("repository: menghapus lampiran menuntut transaksi")
	}
	tabel, _, err := m.tabelLampiran()
	if err != nil {
		return err
	}
	kunci, err := m.kunciTreaty(ctx, tx.tx, tahunID)
	if err != nil {
		return err
	}
	q := sqlHapusLampiranTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, kunci, id)
	if err != nil {
		return fmt.Errorf("repository: menghapus lampiran %s: %w", id, err)
	}
	if n, err := hasil.RowsAffected(); err == nil && n == 0 {
		return ErrLampiranTidakAda
	}
	return pastikanSatuBaris(hasil, "penghapusan lampiran")
}

// SimpanObjek mencatat objek yang sudah ada di penyimpanan (`Insert_T_Storage_SQL`).
func (m *MasterLampiranTCO) SimpanObjek(ctx context.Context, tx *Tx, o models.ObjekPenyimpananTCO) error {
	if tx == nil {
		return errors.New("repository: mencatat objek penyimpanan menuntut transaksi")
	}
	_, objek, err := m.tabelLampiran()
	if err != nil {
		return err
	}
	q := sqlSimpanObjekTCO(objek)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, o.ImageID, kosongJadiNil(o.URLPublic), kosongJadiNil(o.AppFolder),
		kosongJadiNil(o.Exp), kosongJadiNil(o.Namafile), kosongJadiNil(o.App))
	if err != nil {
		// ⛔ URL bertanda tangan tidak disebut: ia memuat tanda tangan akses.
		return fmt.Errorf("repository: mencatat objek penyimpanan %s", o.ImageID)
	}
	return pastikanSatuBaris(hasil, "pencatatan objek penyimpanan")
}

// HapusObjek membuang catatan objek (`DeleteStorage_SQL`); nol baris bukan galat.
func (m *MasterLampiranTCO) HapusObjek(ctx context.Context, tx *Tx, imageID string) error {
	if tx == nil {
		return errors.New("repository: menghapus catatan objek menuntut transaksi")
	}
	_, objek, err := m.tabelLampiran()
	if err != nil {
		return err
	}
	q := sqlHapusObjekTCO(objek)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := tx.tx.ExecContext(ctx, q, imageID); err != nil {
		return fmt.Errorf("repository: menghapus catatan objek %s: %w", imageID, err)
	}
	return nil
}

// sqlPungutEfekRujukanTCO - `sqlPungutEfek` + `AND RUJUKAN = :4`: pemungutan
// SATU efek jatuh tempo milik satu rujukan (temuan /code-review, tiket 12).
func sqlPungutEfekRujukanTCO(tabel string) string {
	return fmt.Sprintf(`SELECT ID, LINI, MODUL, JENIS_EFEK, RUJUKAN, MUATAN,
			   PERCOBAAN
		  FROM %s
		 WHERE STATUS = :1
		   AND JADWAL_BERIKUT <= :2
		   AND MODUL = :3
		   AND RUJUKAN = :4
		 ORDER BY JADWAL_BERIKUT ASC
		 FETCH FIRST 1 ROWS ONLY
		 FOR UPDATE SKIP LOCKED`, tabel)
}

// PungutEfekRujukanTCO - seperti `PungutEfek`, dibatasi satu rujukan: aksi
// seorang pemakai menjalankan efek MILIKNYA, bukan antrean pemakai lain.
func (r *PohonKlaim) PungutEfekRujukanTCO(ctx context.Context, tx *Tx, modul, rujukan string, saat time.Time) (BarisEfekKeluar, error) {
	if tx == nil {
		return BarisEfekKeluar{}, errors.New("repository: PungutEfekRujukanTCO menuntut transaksi")
	}
	tabel, err := r.db.Qualify("T_LOG_SERVICE_RNM")
	if err != nil {
		return BarisEfekKeluar{}, err
	}
	q := sqlPungutEfekRujukanTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return BarisEfekKeluar{}, err
	}
	var (
		b                                   BarisEfekKeluar
		lini, modulB, jenis, rujukanB, muat sql.NullString
		percobaan                           sql.NullInt64
	)
	err = tx.tx.QueryRowContext(ctx, q, StatusEfekAntre, saat, modul, rujukan).Scan(
		&b.ID, &lini, &modulB, &jenis, &rujukanB, &muat, &percobaan)
	if errors.Is(err, sql.ErrNoRows) {
		return BarisEfekKeluar{}, ErrEfekTidakAda
	}
	if err != nil {
		return BarisEfekKeluar{}, fmt.Errorf("repository: memungut efek keluar %s: %w", rujukan, err)
	}
	b.Lini, b.Modul, b.Jenis = lini.String, modulB.String, jenis.String
	b.Rujukan, b.Muatan = rujukanB.String, muat.String
	b.Percobaan = int(percobaan.Int64)
	tandai := sqlTandaiJalan(tabel)
	if err := PeriksaSQL(tandai); err != nil {
		return BarisEfekKeluar{}, err
	}
	hasil, err := tx.tx.ExecContext(ctx, tandai, StatusEfekJalan, saat, b.ID, StatusEfekAntre)
	if err != nil {
		return BarisEfekKeluar{}, fmt.Errorf("repository: menandai efek jalan: %w", err)
	}
	if err := pastikanSatuBaris(hasil, "penandaan efek jalan"); err != nil {
		return BarisEfekKeluar{}, err
	}
	return b, nil
}
