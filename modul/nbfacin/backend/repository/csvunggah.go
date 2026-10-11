package repository

// Unggah CSV objek tanpa batas (tiket 75, migrasi 257): pekerjaan T_NB_CSV_UNGGAH, penampung baris T_NB_CSV_BARIS (satu kolom
// C_<kolom> per kolom template), galat / pesan T_NB_CSV_PESAN. Tulisan di transaksi PEMANGGIL (layanan membuka transaksi pendek
// per langkah supaya kemajuan terlihat dan tahan restart, ADR-U-0029). Baris penampung disisip array-bind per kelompok lewat
// jalankanKelompok (periksaBindBaris + kolomArrayBertipe).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// Tabel migrasi 257.
	TabelUnggahCsv      = "T_NB_CSV_UNGGAH"
	TabelBarisUnggahCsv = "T_NB_CSV_BARIS"
	TabelPesanUnggahCsv = "T_NB_CSV_PESAN"
	sequenceUnggahCsv   = "SEQ_T_NB_CSV_UNGGAH"
	// TabelBarisUnggahAneka - penampung pekerjaan berjenis aneka (migrasi 250, tiket 73 P6, K73-52 opsi 1): tabel sendiri tanpa FK
	// karena 250 berjalan sebelum 257 yang membuat T_NB_CSV_UNGGAH / T_NB_CSV_BARIS.
	TabelBarisUnggahAneka = "T_NB_CSV_BARIS_ANEKA"
)

// ErrUnggahTidakAda - pekerjaan unggah tidak ada (atau milik case lain).
var ErrUnggahTidakAda = errors.New("repository: unggahan CSV tidak ada")

// PenyimpanUnggahCsv - pekerjaan unggah CSV objek (tiket 75).
type PenyimpanUnggahCsv interface {
	// BuatUnggah - pekerjaan baru (DIBUAT = DIUBAH = SYSTIMESTAMP); jawab ID.
	BuatUnggah(ctx context.Context, tx *db.Tx, u models.UnggahCsv) (string, error)
	// UbahUnggah - semua medan yang berubah + DIUBAH = SYSTIMESTAMP (detak).
	UbahUnggah(ctx context.Context, tx *db.Tx, u models.UnggahCsv) error
	BacaUnggah(ctx context.Context, id string) (models.UnggahCsv, error)
	// UnggahTerakhir - pekerjaan terbaru case (ada = false bila tidak ada).
	UnggahTerakhir(ctx context.Context, caseID string) (models.UnggahCsv, bool, error)
	// SisipBarisUnggah - satu kelompok baris penampung pekerjaan berjenis `jenis` (Isi urut models.KolomCsvJenis).
	SisipBarisUnggah(ctx context.Context, tx *db.Tx, id, jenis string, baris []models.BarisUnggahCsv) error
	// BacaBarisUnggah - baris NO_BARIS dari..sampai (inklusif), urut; Isi urut models.KolomCsvJenis(jenis).
	BacaBarisUnggah(ctx context.Context, id, jenis string, dari, sampai int) ([]models.BarisUnggahCsv, error)
	// SisipPesanUnggah - galat / pesan jenis `jenis`, URUT mulai `urut`.
	SisipPesanUnggah(ctx context.Context, tx *db.Tx, id, jenis string, urut int, pesan []models.PesanUnggahCsv) error
	// BacaPesanUnggah - paling banyak `batas` galat / pesan pertama jenis `jenis`.
	BacaPesanUnggah(ctx context.Context, id, jenis string, batas int) ([]models.PesanUnggahCsv, error)
	// HapusPesanUnggah - seluruh galat / pesan pekerjaan (validasi ulang).
	HapusPesanUnggah(ctx context.Context, tx *db.Tx, id string) error
	// AdaObjek - case sudah memuat sekurangnya satu objek jenis itu (fire: lokasi T_LOCATIONLIST; cargo: T_CARGOLIST).
	AdaObjek(ctx context.Context, caseID, jenis string) (bool, error)
}

// UnggahCsvOracle - PenyimpanUnggahCsv atas Oracle.
type UnggahCsvOracle struct{ db *db.DB }

// NewUnggahCsvOracle merakit penyimpan unggah CSV.
func NewUnggahCsvOracle(d *db.DB) *UnggahCsvOracle { return &UnggahCsvOracle{db: d} }

type tabelUnggah struct{ unggah, baris, barisAneka, pesan, loc, cargo string }

// penampung - tabel baris pekerjaan berjenis `jenis` (teks berkualifikasi, nama untuk pesan).
func (t tabelUnggah) penampung(jenis string) (string, string) {
	if jenis == models.JenisCsvAneka {
		return t.barisAneka, TabelBarisUnggahAneka
	}
	return t.baris, TabelBarisUnggahCsv
}

func (r *UnggahCsvOracle) tabel() (tabelUnggah, error) {
	var t tabelUnggah
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelUnggahCsv, &t.unggah}, {TabelBarisUnggahCsv, &t.baris}, {TabelPesanUnggahCsv, &t.pesan}, {TabelLocationList, &t.loc},
		{TabelCargoList, &t.cargo}, {TabelBarisUnggahAneka, &t.barisAneka}} {
		n, err := r.db.Qualify(x.nama)
		if err != nil {
			return t, err
		}
		*x.ke = n
	}
	return t, nil
}

// kolomUnggahUbah - kolom yang ditulis UbahUnggah (urutan bind :1..:n, lalu ID).
var kolomUnggahUbah = []string{"STATUS", "FASE", "JUMLAH_BARIS", "BARIS_DIPROSES", "BARIS_TERSIMPAN", "KELOMPOK_TERAKHIR",
	"JUMLAH_GALAT", "JUMLAH_PESAN", "NO_OFFER_SLIP", "BEGIN_DATE", "END_DATE", "GALAT_PROSES", "GALAT_DI_BARIS", "UKURAN"}

func argUnggahUbah(u models.UnggahCsv) []any {
	var diBaris any
	if u.GalatDiBaris > 0 {
		diBaris = u.GalatDiBaris
	}
	return []any{u.Status, u.Fase, u.TotalBaris, u.BarisDiproses, u.BarisTersimpan, u.KelompokTerakhir, u.JumlahGalat,
		u.JumlahPesan, db.KosongJadiNil(u.NoOfferSlip), db.KosongJadiNil(u.BeginDate), db.KosongJadiNil(u.EndDate),
		db.KosongJadiNil(potong4000(u.GalatProses)), diBaris, u.Ukuran}
}

// potong4000 - teks VARCHAR2(4000) BYTE: dipotong pada batas rune.
func potong4000(s string) string {
	if len(s) <= 4000 {
		return s
	}
	n := 4000
	for n > 0 && (s[n]&0xC0) == 0x80 {
		n--
	}
	return s[:n]
}

func sqlBuatUnggah(t string) string {
	return "INSERT INTO " + t + " (ID, CASE_ID, STATUS, FASE, NAMA_BERKAS, UKURAN, DIBUAT_OLEH, JENIS, KOLOM, DIBUAT, DIUBAH) VALUES" +
		" (:1, :2, :3, :4, :5, :6, :7, :8, :9, SYSTIMESTAMP, SYSTIMESTAMP)"
}

// pemisahKolomUnggah - pemisah judul T_NB_CSV_UNGGAH.KOLOM.
const pemisahKolomUnggah = ";"

// jenisUnggah - jenis tersimpan (kosong = fire).
func jenisUnggah(j string) string {
	if j == "" {
		return models.JenisCsvFire
	}
	return j
}

func sqlUbahUnggah(t string) string {
	set := make([]string, len(kolomUnggahUbah))
	for i, k := range kolomUnggahUbah {
		set[i] = k + " = :" + strconv.Itoa(i+1)
	}
	return "UPDATE " + t + " SET " + strings.Join(set, ", ") + ", DIUBAH = SYSTIMESTAMP WHERE ID = :" + strconv.Itoa(len(set)+1)
}

// kolomBacaUnggah - kolom pekerjaan yang dibaca (urutan Scan bacaUnggah).
var kolomBacaUnggah = "TO_CHAR(ID), CASE_ID, STATUS, FASE, NAMA_BERKAS, UKURAN, JUMLAH_BARIS, BARIS_DIPROSES, BARIS_TERSIMPAN," +
	" KELOMPOK_TERAKHIR, JUMLAH_GALAT, JUMLAH_PESAN, NO_OFFER_SLIP, BEGIN_DATE, END_DATE, GALAT_PROSES, GALAT_DI_BARIS, DIBUAT_OLEH," +
	" CAST(DIBUAT AS DATE), CAST(DIUBAH AS DATE), CASE WHEN DIUBAH < SYSTIMESTAMP - NUMTODSINTERVAL(" +
	strconv.Itoa(models.BatasDetakUnggahMenit) + ", 'MINUTE') THEN 1 ELSE 0 END, JENIS, KOLOM"

func sqlBacaUnggah(t string) string {
	return "SELECT " + kolomBacaUnggah + " FROM " + t + " WHERE ID = :1"
}

func sqlUnggahTerakhir(t string) string {
	return "SELECT " + kolomBacaUnggah + " FROM " + t + " WHERE ID = (SELECT MAX(ID) FROM " + t + " WHERE CASE_ID = :1)"
}

// kolomBarisUnggah - C_<kolom> penampung pekerjaan berjenis `jenis` (urutan models.KolomCsvJenis).
func kolomBarisUnggah(jenis string) []string {
	sumber := models.KolomCsvJenis(jenis)
	k := make([]string, len(sumber))
	for i, x := range sumber {
		k[i] = "C_" + x
	}
	return k
}

func sqlSisipBarisUnggah(t, jenis string) string {
	kolom := append([]string{"UNGGAH_ID", "NO_BARIS"}, kolomBarisUnggah(jenis)...)
	nilai := make([]string, len(kolom))
	for i := range nilai {
		nilai[i] = ":" + strconv.Itoa(i+1)
	}
	return "INSERT INTO " + t + " (" + strings.Join(kolom, ", ") + ") VALUES (" + strings.Join(nilai, ", ") + ")"
}

func sqlBacaBarisUnggah(t, jenis string) string {
	return "SELECT NO_BARIS, " + strings.Join(kolomBarisUnggah(jenis), ", ") + " FROM " + t +
		" WHERE UNGGAH_ID = :1 AND NO_BARIS BETWEEN :2 AND :3 ORDER BY NO_BARIS"
}

func sqlSisipPesanUnggah(t string) string {
	return "INSERT INTO " + t + " (UNGGAH_ID, URUT, JENIS, NO_BARIS, TEKS) VALUES (:1, :2, :3, :4, :5)"
}

func sqlBacaPesanUnggah(t string) string {
	return "SELECT NO_BARIS, TEKS FROM " + t + " WHERE UNGGAH_ID = :1 AND JENIS = :2 ORDER BY URUT FETCH FIRST :3 ROWS ONLY"
}

// BuatUnggah - lihat PenyimpanUnggahCsv.
func (r *UnggahCsvOracle) BuatUnggah(ctx context.Context, tx *db.Tx, u models.UnggahCsv) (string, error) {
	t, err := r.tabel()
	if err != nil {
		return "", err
	}
	id, err := r.db.NomorBerikut(ctx, tx, sequenceUnggahCsv)
	if err != nil {
		return "", err
	}
	_, err = jalankan(ctx, tx, sqlBuatUnggah(t.unggah), "menyisipkan "+TabelUnggahCsv, id, u.CaseID, u.Status, u.Fase,
		db.KosongJadiNil(potong4000(u.NamaBerkas)), u.Ukuran, db.KosongJadiNil(u.DibuatOleh), jenisUnggah(u.Jenis),
		db.KosongJadiNil(potong4000(strings.Join(u.Kolom, pemisahKolomUnggah))))
	return id, err
}

// UbahUnggah - lihat PenyimpanUnggahCsv.
func (r *UnggahCsvOracle) UbahUnggah(ctx context.Context, tx *db.Tx, u models.UnggahCsv) error {
	t, err := r.tabel()
	if err != nil {
		return err
	}
	h, err := jalankan(ctx, tx, sqlUbahUnggah(t.unggah), "mengubah "+TabelUnggahCsv, append(argUnggahUbah(u), u.ID)...)
	if err != nil {
		return err
	}
	if n, err := h.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrUnggahTidakAda
	}
	return nil
}

func bacaUnggah(baris interface{ Scan(...any) error }) (models.UnggahCsv, error) {
	var u models.UnggahCsv
	var nama, offer, mulai, akhir, galat, oleh, jenis, kolom sql.NullString
	var ukuran sql.NullInt64
	var diBaris sql.NullInt64
	var dibuat, diubah sql.NullTime
	var basi int
	err := baris.Scan(&u.ID, &u.CaseID, &u.Status, &u.Fase, &nama, &ukuran, &u.TotalBaris, &u.BarisDiproses, &u.BarisTersimpan,
		&u.KelompokTerakhir, &u.JumlahGalat, &u.JumlahPesan, &offer, &mulai, &akhir, &galat, &diBaris, &oleh, &dibuat, &diubah, &basi,
		&jenis, &kolom)
	if err != nil {
		return u, err
	}
	u.NamaBerkas, u.Ukuran, u.NoOfferSlip, u.BeginDate, u.EndDate = nama.String, ukuran.Int64, offer.String, mulai.String, akhir.String
	u.GalatProses, u.GalatDiBaris, u.DibuatOleh = galat.String, int(diBaris.Int64), oleh.String
	u.Dibuat, u.Diubah, u.Basi = dibuat.Time, diubah.Time, basi == 1
	u.Jenis = jenisUnggah(jenis.String)
	if kolom.String != "" {
		u.Kolom = strings.Split(kolom.String, pemisahKolomUnggah)
	}
	return u, nil
}

// BacaUnggah - lihat PenyimpanUnggahCsv.
func (r *UnggahCsvOracle) BacaUnggah(ctx context.Context, id string) (models.UnggahCsv, error) {
	t, err := r.tabel()
	if err != nil {
		return models.UnggahCsv{}, err
	}
	u, err := bacaUnggah(r.db.QueryRowContext(ctx, sqlBacaUnggah(t.unggah), id))
	if errors.Is(err, sql.ErrNoRows) {
		return u, ErrUnggahTidakAda
	}
	if err != nil {
		return u, fmt.Errorf("repository: baca %s: %w", TabelUnggahCsv, err)
	}
	return u, nil
}

// UnggahTerakhir - lihat PenyimpanUnggahCsv.
func (r *UnggahCsvOracle) UnggahTerakhir(ctx context.Context, caseID string) (models.UnggahCsv, bool, error) {
	t, err := r.tabel()
	if err != nil {
		return models.UnggahCsv{}, false, err
	}
	u, err := bacaUnggah(r.db.QueryRowContext(ctx, sqlUnggahTerakhir(t.unggah), caseID))
	if errors.Is(err, sql.ErrNoRows) {
		return u, false, nil
	}
	if err != nil {
		return u, false, fmt.Errorf("repository: baca %s: %w", TabelUnggahCsv, err)
	}
	return u, true, nil
}

// SisipBarisUnggah - lihat PenyimpanUnggahCsv (array-bind, potongan ukuranKelompok, satu transaksi per kelompok).
func (r *UnggahCsvOracle) SisipBarisUnggah(ctx context.Context, tx *db.Tx, id, jenis string, baris []models.BarisUnggahCsv) error {
	if len(baris) == 0 {
		return nil
	}
	t, err := r.tabel()
	if err != nil {
		return err
	}
	tabel, nama := t.penampung(jenis)
	k := &kelompokSisip{q: sqlSisipBarisUnggah(tabel, jenis), tabel: nama, baris: argBarisUnggah(id, jenis, baris)}
	return jalankanKelompok(ctx, oracleKelompok{d: r.db, tx: tx}, k, false)
}

// argBarisUnggah - nilai bind per baris: UNGGAH_ID, NO_BARIS, sel kolom jenis (kosong -> NULL; kurang -> NULL).
func argBarisUnggah(id, jenis string, baris []models.BarisUnggahCsv) [][]any {
	kolom := models.KolomCsvJenis(jenis)
	hasil := make([][]any, len(baris))
	for i, b := range baris {
		arg := make([]any, 0, 2+len(kolom))
		arg = append(arg, id, b.No)
		for c := range kolom {
			var v any
			if c < len(b.Isi) && b.Isi[c] != "" {
				v = potong4000(b.Isi[c])
			}
			arg = append(arg, v)
		}
		hasil[i] = arg
	}
	return hasil
}

// BacaBarisUnggah - lihat PenyimpanUnggahCsv.
func (r *UnggahCsvOracle) BacaBarisUnggah(ctx context.Context, id, jenis string, dari, sampai int) ([]models.BarisUnggahCsv, error) {
	t, err := r.tabel()
	if err != nil {
		return nil, err
	}
	tabel, nama := t.penampung(jenis)
	baris, err := r.db.QueryContext(ctx, sqlBacaBarisUnggah(tabel, jenis), id, dari, sampai)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", nama, err)
	}
	defer baris.Close()
	hasil := []models.BarisUnggahCsv{}
	n := len(models.KolomCsvJenis(jenis))
	for baris.Next() {
		var no int
		sel := make([]sql.NullString, n)
		tuju := make([]any, 0, n+1)
		tuju = append(tuju, &no)
		for i := range sel {
			tuju = append(tuju, &sel[i])
		}
		if err := baris.Scan(tuju...); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", nama, err)
		}
		isi := make([]string, n)
		for i := range sel {
			isi[i] = sel[i].String
		}
		hasil = append(hasil, models.BarisUnggahCsv{No: no, Isi: isi})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", nama, err)
	}
	return hasil, nil
}

// SisipPesanUnggah - lihat PenyimpanUnggahCsv.
func (r *UnggahCsvOracle) SisipPesanUnggah(ctx context.Context, tx *db.Tx, id, jenis string, urut int, pesan []models.PesanUnggahCsv) error {
	if len(pesan) == 0 {
		return nil
	}
	t, err := r.tabel()
	if err != nil {
		return err
	}
	baris := make([][]any, len(pesan))
	for i, p := range pesan {
		var no any
		if p.NoBaris > 0 {
			no = p.NoBaris
		}
		baris[i] = []any{id, urut + i, jenis, no, potong4000(p.Teks)}
	}
	k := &kelompokSisip{q: sqlSisipPesanUnggah(t.pesan), tabel: TabelPesanUnggahCsv, baris: baris}
	return jalankanKelompok(ctx, oracleKelompok{d: r.db, tx: tx}, k, false)
}

// BacaPesanUnggah - lihat PenyimpanUnggahCsv.
func (r *UnggahCsvOracle) BacaPesanUnggah(ctx context.Context, id, jenis string, batas int) ([]models.PesanUnggahCsv, error) {
	t, err := r.tabel()
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlBacaPesanUnggah(t.pesan), id, jenis, batas)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelPesanUnggahCsv, err)
	}
	defer baris.Close()
	hasil := []models.PesanUnggahCsv{}
	for baris.Next() {
		var no sql.NullInt64
		var teks string
		if err := baris.Scan(&no, &teks); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelPesanUnggahCsv, err)
		}
		hasil = append(hasil, models.PesanUnggahCsv{NoBaris: int(no.Int64), Teks: teks})
	}
	return hasil, baris.Err()
}

// HapusPesanUnggah - lihat PenyimpanUnggahCsv.
func (r *UnggahCsvOracle) HapusPesanUnggah(ctx context.Context, tx *db.Tx, id string) error {
	t, err := r.tabel()
	if err != nil {
		return err
	}
	_, err = jalankan(ctx, tx, "DELETE FROM "+t.pesan+" WHERE UNGGAH_ID = :1", "menghapus "+TabelPesanUnggahCsv, id)
	return err
}

// AdaObjek - lihat PenyimpanUnggahCsv.
func (r *UnggahCsvOracle) AdaObjek(ctx context.Context, caseID, jenis string) (bool, error) {
	t, err := r.tabel()
	if err != nil {
		return false, err
	}
	tabel, nama := t.loc, TabelLocationList
	if jenis == models.JenisCsvCargo {
		tabel, nama = t.cargo, TabelCargoList
	}
	var n int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+tabel+" WHERE PARENT_ID = :1 AND ROWNUM = 1", caseID).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: baca %s: %w", nama, err)
	}
	return n > 0, nil
}
