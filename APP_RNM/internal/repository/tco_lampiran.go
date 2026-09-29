package repository

// Lampiran tahun treaty - tiket 12 Treaty Contract Out (FITUR BARU).
//
// Untuk apa berkas ini: rekam `T_TREATYYEAR_LAMPIRAN` dan master kategori
// `CATEGORY_ATTACH_REAS` (dibaca saja). Isi berkas TIDAK di sini - ia di
// penyimpanan, di balik antarmuka di services.
//
// ⛔ Lampiran melekat pada TAHUN treaty (`IDTREATYYEAR`). Rujukan ke kunci
// treaty inward tidak ada di berkas ini; penjaga
// `TestTCOLampiranTanpaRujukanTreatyInward` menegakkannya.
//
// ⛔ Status "gagal" dibaca dari outbox bersama `T_LOG_SERVICE_RNM` - HANYA
// dibaca, disaring `MODUL` modul ini. Tulisan ke outbox lewat fungsi yang sudah
// ada (`PohonKlaim.AntreEfek` dan kawan-kawan), bukan SQL baru.
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

// MasterLampiranTCO membaca dan menulis `T_TREATYYEAR_LAMPIRAN`.
type MasterLampiranTCO struct{ db *DB }

// NewMasterLampiranTCO menyusun gudangnya.
func NewMasterLampiranTCO(db *DB) *MasterLampiranTCO { return &MasterLampiranTCO{db: db} }

const pilihLampiranTCO = `l.ID, l.IDTREATYYEAR, l.FILENAME, l.FILEMIMETYPE, l.CATEGORY, l.IMAGEID,
	       l.T_STORAGE_ID, TO_CHAR(l.UKURAN), l.USERID, TO_CHAR(l.TGLUPLOAD, 'YYYY-MM-DD HH24:MI:SS')`

// sqlDaftarLampiranTCO membaca lampiran satu tahun + efek unggah terakhirnya.
//
// ⛔ Efek TERAKHIR per lampiran (`ROW_NUMBER ... DIBUAT DESC`): pengulangan
// menambah baris outbox baru, dan status yang tampil adalah nasib percobaan
// terbaru, bukan yang pertama.
func sqlDaftarLampiranTCO(tabel, outbox string, satu bool) string {
	saring := ""
	if satu {
		saring = " AND l.ID = :4"
	}
	return fmt.Sprintf(`SELECT %s,
	       o.STATUS, o.GALAT_TERAKHIR, o.PERCOBAAN
	  FROM %s l
	  LEFT JOIN (SELECT RUJUKAN, STATUS, GALAT_TERAKHIR, PERCOBAAN,
	                    ROW_NUMBER() OVER (PARTITION BY RUJUKAN ORDER BY DIBUAT DESC, ID DESC) AS URUT
	               FROM %s
	              WHERE MODUL = :1 AND JENIS_EFEK = :2) o
	    ON o.RUJUKAN = l.ID AND o.URUT = 1
	 WHERE l.IDTREATYYEAR = :3%s
	 ORDER BY l.ID DESC`, pilihLampiranTCO, tabel, outbox, saring)
}

func sqlAmbilUntukKirimLampiranTCO(tabel string) string {
	return fmt.Sprintf(`SELECT %s FROM %s l WHERE l.ID = :1 FOR UPDATE`, pilihLampiranTCO, tabel)
}

func sqlSisipLampiranTCO(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s
	       (ID, IDTREATYYEAR, FILENAME, FILEMIMETYPE, CATEGORY, IMAGEID, T_STORAGE_ID,
	        UKURAN, USERID, TGLUPLOAD)
	VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10)`, tabel)
}

func sqlHapusLampiranTCO(tabel string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE ID = :1 AND IDTREATYYEAR = :2`, tabel)
}

func sqlTandaiLampiranTCO(tabel string) string {
	return fmt.Sprintf(`UPDATE %s SET T_STORAGE_ID = :1 WHERE ID = :2`, tabel)
}

// pindaiLampiranTCO membaca sepuluh kolom `pilihLampiranTCO`.
func pindaiLampiranTCO(n [10]sql.NullString) (models.LampiranTCO, error) {
	l := models.LampiranTCO{
		ID: n[0].String, IDTreatyYear: n[1].String, FileName: n[2].String,
		FileMimeType: n[3].String, Category: n[4].String, ImageID: n[5].String,
		TStorageID: n[6].String, UserID: n[8].String,
	}
	if s := strings.TrimSpace(n[7].String); s != "" {
		u, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return l, fmt.Errorf("repository: UKURAN lampiran %s bernilai %q: %w", l.ID, s, err)
		}
		l.Ukuran = u
	}
	var err error
	l.TglUpload, err = uraiTanggalTeks(n[9], "TGLUPLOAD")
	return l, err
}

// Daftar membaca lampiran satu tahun treaty, ID DESC.
func (m *MasterLampiranTCO) Daftar(ctx context.Context, modul, jenisUnggah, tahunID string) (
	[]BarisLampiranTCO, error) {
	return m.baca(ctx, modul, jenisUnggah, tahunID, "")
}

// Ambil membaca SATU lampiran, dibatasi tahun treaty-nya.
//
// ⛔ Dibatasi tahun, bukan dicari dengan ID saja: ID lampiran datang dari
// jalur URL, dan tanpa batas itu jalur tahun mana pun menyerahkan lampiran
// tahun lain.
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

func (m *MasterLampiranTCO) baca(ctx context.Context, modul, jenisUnggah, tahunID, id string) (
	[]BarisLampiranTCO, error) {
	tabel, err := m.db.Qualify(TabelLampiranTCO)
	if err != nil {
		return nil, err
	}
	outbox, err := m.db.Qualify(TabelOutboxBersama)
	if err != nil {
		return nil, err
	}
	q := sqlDaftarLampiranTCO(tabel, outbox, id != "")
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	arg := []any{modul, jenisUnggah, tahunID}
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
		var n [10]sql.NullString
		var status, galat sql.NullString
		var percobaan sql.NullInt64
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6], &n[7], &n[8], &n[9],
			&status, &galat, &percobaan); err != nil {
			return nil, err
		}
		l, err := pindaiLampiranTCO(n)
		if err != nil {
			return nil, err
		}
		hasil = append(hasil, BarisLampiranTCO{LampiranTCO: l, StatusEfek: status.String,
			GalatEfek: galat.String, PercobaanEfek: int(percobaan.Int64)})
	}
	return hasil, rows.Err()
}

// AmbilUntukKirim mengunci satu lampiran untuk pelaksana efek.
//
// ⛔ `FOR UPDATE`: dua pelaksana atas lampiran yang sama menunggu satu sama
// lain, dan yang kedua melihat `T_STORAGE_ID` yang sudah terisi.
func (m *MasterLampiranTCO) AmbilUntukKirim(ctx context.Context, tx *Tx, id string) (models.LampiranTCO, error) {
	if tx == nil {
		return models.LampiranTCO{}, errors.New("repository: mengunci lampiran menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelLampiranTCO)
	if err != nil {
		return models.LampiranTCO{}, err
	}
	q := sqlAmbilUntukKirimLampiranTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return models.LampiranTCO{}, err
	}
	var n [10]sql.NullString
	err = tx.tx.QueryRowContext(ctx, q, id).Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5],
		&n[6], &n[7], &n[8], &n[9])
	if errors.Is(err, sql.ErrNoRows) {
		return models.LampiranTCO{}, ErrLampiranTidakAda
	}
	if err != nil {
		return models.LampiranTCO{}, fmt.Errorf("repository: mengunci lampiran %s: %w", id, err)
	}
	return pindaiLampiranTCO(n)
}

// Sisip menulis satu lampiran; ID dari `SEQ_T_TREATYYEAR_LAMPIRAN`.
func (m *MasterLampiranTCO) Sisip(ctx context.Context, tx *Tx, l models.LampiranTCO) (string, error) {
	if tx == nil {
		return "", errors.New("repository: menyisipkan lampiran menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelLampiranTCO)
	if err != nil {
		return "", err
	}
	id, err := m.db.IdentitasBerikutTCO(ctx, tx, SeqLampiranTCO)
	if err != nil {
		return "", err
	}
	q := sqlSisipLampiranTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, id, l.IDTreatyYear, kosongJadiNil(l.FileName),
		kosongJadiNil(l.FileMimeType), kosongJadiNil(l.Category), l.ImageID,
		kosongJadiNil(l.TStorageID), l.Ukuran, kosongJadiNil(l.UserID), tanggalJadiNil(l.TglUpload))
	if err != nil {
		return "", fmt.Errorf("repository: menyisipkan lampiran: %w", err)
	}
	return id, pastikanSatuBaris(hasil, "penyisipan lampiran")
}

// Hapus membuang satu lampiran milik tahun treaty itu.
func (m *MasterLampiranTCO) Hapus(ctx context.Context, tx *Tx, tahunID, id string) error {
	if tx == nil {
		return errors.New("repository: menghapus lampiran menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelLampiranTCO)
	if err != nil {
		return err
	}
	q := sqlHapusLampiranTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, id, tahunID)
	if err != nil {
		return fmt.Errorf("repository: menghapus lampiran %s: %w", id, err)
	}
	if n, err := hasil.RowsAffected(); err == nil && n == 0 {
		return ErrLampiranTidakAda
	}
	return pastikanSatuBaris(hasil, "penghapusan lampiran")
}

// TandaiTerkirim menulis `T_STORAGE_ID`; teks kosong mengosongkannya lagi
// (perbaikan rekam yang berkasnya hilang di penyimpanan).
func (m *MasterLampiranTCO) TandaiTerkirim(ctx context.Context, tx *Tx, id, storageID string) error {
	if tx == nil {
		return errors.New("repository: menandai lampiran menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelLampiranTCO)
	if err != nil {
		return err
	}
	q := sqlTandaiLampiranTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, kosongJadiNil(storageID), id)
	if err != nil {
		return fmt.Errorf("repository: menandai lampiran %s: %w", id, err)
	}
	if n, err := hasil.RowsAffected(); err == nil && n == 0 {
		return ErrLampiranTidakAda
	}
	return pastikanSatuBaris(hasil, "penandaan lampiran")
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
