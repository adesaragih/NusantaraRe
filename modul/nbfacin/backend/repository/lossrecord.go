package repository

// Tab utama Loss Record kasus FIRE (tiket 55) - tabel BARU migrasi 206 (tidak ada di rancangan loader: jalur tingkat case
// `.OfferFacIn.ListCauseOfLoss`, Premi PAR / MB / PL / Other, loss ratio 3 / 5 tahun dan net tidak ada di workbook):
// T_LOSSRECORD (induk T_GENERAL_POLIS, satu per kasus) -> T_LOSSRECORDLIST (baris kerugian, urut SEQ_NO). Simpan = ganti
// utuh di satu transaksi. Jabatan akun (aksesLossRecord) dibaca dari M_LOGIN_GO milik inti login (baca saja).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelLossRecord, TabelLossRecordList - tabel migrasi 206.
	TabelLossRecord     = "T_LOSSRECORD"
	TabelLossRecordList = "T_LOSSRECORDLIST"
	// TabelLoginGo - akun login (inti, migrasi 902); hanya JOB_POSITION yang dibaca.
	TabelLoginGo = "M_LOGIN_GO"

	sequenceLossRecord     = "SEQ_T_LOSSRECORD"
	sequenceLossRecordList = "SEQ_T_LOSSRECORDLIST"
)

// PenyimpanLossRecord - baca / ganti Loss Record satu kasus.
type PenyimpanLossRecord interface {
	// BacaLossRecord - isi tersimpan (kosong bila belum pernah disimpan); ErrKasusTidakAda bila case FAC tidak ada.
	BacaLossRecord(ctx context.Context, id string) (models.LossRecord, error)
	// TulisLossRecord - ganti utuh (di transaksi pemanggil); ErrKasusTidakAda bila case FAC tidak ada.
	TulisLossRecord(ctx context.Context, tx *db.Tx, id string, l models.LossRecord) error
}

// PembacaJabatan - JOB_POSITION akun login.
type PembacaJabatan interface {
	// JabatanAkun - M_LOGIN_GO.JOB_POSITION apa adanya; "" bila akun tidak ada / kolom kosong.
	JabatanAkun(ctx context.Context, akunID string) (string, error)
}

// LossRecordOracle - PenyimpanLossRecord dan PembacaJabatan atas Oracle.
type LossRecordOracle struct{ db *db.DB }

// NewLossRecordOracle merakit penyimpan Loss Record.
func NewLossRecordOracle(d *db.DB) *LossRecordOracle { return &LossRecordOracle{db: d} }

type tabelLossRecord struct{ work, general, kepala, baris string }

func (r *LossRecordOracle) tabel() (tabelLossRecord, error) {
	var t tabelLossRecord
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelWorkPolis, &t.work}, {TabelGeneralPolis, &t.general}, {TabelLossRecord, &t.kepala}, {TabelLossRecordList, &t.baris}} {
		q, err := r.db.Qualify(x.nama)
		if err != nil {
			return t, err
		}
		*x.ke = q
	}
	return t, nil
}

// kolomKepalaLossRecord - kolom uang / rasio T_LOSSRECORD = urutan argumen tulis dan baca.
var kolomKepalaLossRecord = []string{"PREMI_PAR", "PREMI_MB", "PREMI_PL", "PREMI_OTHER", "LOSS_RATIO1_YEAR_PERCENT",
	"LOSS_RATIO3_YEAR_PERCENT", "LOSS_RATIO5_YEAR_PERCENT", "NET_LOSS_RATIO1_YEAR", "NET_LOSS_RATIO3_YEAR", "NET_LOSS_RATIO5_YEAR"}

func penunjukKepala(l *models.LossRecord) []**apd.Decimal {
	return []**apd.Decimal{&l.PremiPar, &l.PremiMb, &l.PremiPl, &l.PremiOther, &l.Rasio.Gross1, &l.Rasio.Gross3, &l.Rasio.Gross5,
		&l.Rasio.Net1, &l.Rasio.Net3, &l.Rasio.Net5}
}

// sqlBacaKepalaLossRecord - case FAC (LINI) LEFT JOIN kepala: tanpa baris case -> ErrKasusTidakAda; tanpa kepala -> NULL.
func sqlBacaKepalaLossRecord(t tabelLossRecord) string {
	pilih := make([]string, len(kolomKepalaLossRecord))
	for i, k := range kolomKepalaLossRecord {
		pilih[i] = angkaKeluar("k." + k)
	}
	return "SELECT TO_CHAR(k.ID), " + strings.Join(pilih, ", ") + " FROM " + t.work + " w LEFT JOIN " + t.kepala +
		" k ON k.PARENT_ID = w.ID WHERE w.ID = :1 AND w.LINI = :2"
}

func sqlBacaBarisLossRecord(t tabelLossRecord) string {
	return "SELECT b.UW_YEAR, TO_CHAR(b.TAHUN_KE), b.CLASS_OF_BUSINESS, b.DATE_OF_LOSS, b.CAUSE_OF_LOSS, b.CURRENCY, " +
		angkaKeluar("b.LOSS_AMOUNT") + ", " + angkaKeluar("b.LOSS_AMOUNT_IDR") + ", " + angkaKeluar("b.KURS_IDR") + " FROM " + t.baris +
		" b WHERE b.PARENT_ID = :1 ORDER BY b.SEQ_NO, b.ID"
}

func sqlHapusLossRecord(t tabelLossRecord) []string {
	return []string{
		"DELETE FROM " + t.baris + " WHERE PARENT_ID IN (SELECT k.ID FROM " + t.kepala + " k WHERE k.PARENT_ID = :1)",
		"DELETE FROM " + t.kepala + " WHERE PARENT_ID = :1",
	}
}

func sqlSisipKepalaLossRecord(kepala string) string {
	nilai := make([]string, len(kolomKepalaLossRecord))
	for i := range nilai {
		nilai[i] = angkaMasuk(fmt.Sprintf(":%d", i+3))
	}
	return "INSERT INTO " + kepala + " (ID, PARENT_ID, " + strings.Join(kolomKepalaLossRecord, ", ") + ") VALUES (:1, :2, " +
		strings.Join(nilai, ", ") + ")"
}

func sqlSisipBarisLossRecord(baris string) string {
	return "INSERT INTO " + baris + " (ID, PARENT_ID, SEQ_NO, ROW_UID, UW_YEAR, TAHUN_KE, CLASS_OF_BUSINESS, DATE_OF_LOSS," +
		" CAUSE_OF_LOSS, CURRENCY, LOSS_AMOUNT, LOSS_AMOUNT_IDR, KURS_IDR) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, " +
		angkaMasuk(":11") + ", " + angkaMasuk(":12") + ", " + angkaMasuk(":13") + ")"
}

func sqlJabatanAkun(login string) string {
	return "SELECT JOB_POSITION FROM " + login + " WHERE LOGIN_ID = :1"
}

// BacaLossRecord - lihat PenyimpanLossRecord.
func (r *LossRecordOracle) BacaLossRecord(ctx context.Context, id string) (models.LossRecord, error) {
	l := models.LossRecord{Baris: []models.BarisLossRecord{}}
	t, err := r.tabel()
	if err != nil {
		return l, err
	}
	teks := make([]sql.NullString, len(kolomKepalaLossRecord)+1)
	tujuan := make([]any, len(teks))
	for i := range teks {
		tujuan[i] = &teks[i]
	}
	err = r.db.QueryRowContext(ctx, sqlBacaKepalaLossRecord(t), id, LiniFacIn).Scan(tujuan...)
	if errors.Is(err, sql.ErrNoRows) {
		return l, ErrKasusTidakAda
	}
	if err != nil {
		return l, fmt.Errorf("repository: baca %s: %w", TabelLossRecord, err)
	}
	if !teks[0].Valid {
		return l, nil // belum pernah disimpan
	}
	for i, ke := range penunjukKepala(&l) {
		if *ke, err = bacaDesimal(id, kolomKepalaLossRecord[i], &teks[i+1]); err != nil {
			return l, err
		}
	}
	baris, err := r.db.QueryContext(ctx, sqlBacaBarisLossRecord(t), teks[0].String)
	if err != nil {
		return l, fmt.Errorf("repository: baca %s: %w", TabelLossRecordList, err)
	}
	defer baris.Close()
	for baris.Next() {
		var uw, ke, cob, dol, sebab, mu, jumlah, idr, kurs sql.NullString
		if err := baris.Scan(&uw, &ke, &cob, &dol, &sebab, &mu, &jumlah, &idr, &kurs); err != nil {
			return l, fmt.Errorf("repository: %s: %w", TabelLossRecordList, err)
		}
		b := models.BarisLossRecord{UwYear: uw.String, ClassOfBusiness: cob.String, DateOfLoss: dol.String, CauseOfLoss: sebab.String,
			Currency: mu.String}
		b.TahunKe, _ = strconv.Atoi(ke.String) // kosong -> 0 (di luar jendela)
		for _, d := range []struct {
			kolom string
			v     *sql.NullString
			ke    **apd.Decimal
		}{{"LOSS_AMOUNT", &jumlah, &b.LossAmount}, {"LOSS_AMOUNT_IDR", &idr, &b.LossAmountIDR}, {"KURS_IDR", &kurs, &b.KursIDR}} {
			if *d.ke, err = bacaDesimal(id, d.kolom, d.v); err != nil {
				return l, err
			}
		}
		l.Baris = append(l.Baris, b)
	}
	if err := baris.Err(); err != nil {
		return l, fmt.Errorf("repository: %s: %w", TabelLossRecordList, err)
	}
	return l, nil
}

// TulisLossRecord - lihat PenyimpanLossRecord: sentuh case (FAC), pastikan T_GENERAL_POLIS, hapus, sisip ulang (SEQ_NO =
// urutan l.Baris).
func (r *LossRecordOracle) TulisLossRecord(ctx context.Context, tx *db.Tx, id string, l models.LossRecord) error {
	t, err := r.tabel()
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, sqlSentuhCase(t.work), "menyentuh "+TabelWorkPolis, id, LiniFacIn)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrKasusTidakAda
	}
	if _, err := jalankan(ctx, tx, sqlPastikanGeneral(t.general), "memastikan "+TabelGeneralPolis, id, id); err != nil {
		return err
	}
	for _, q := range sqlHapusLossRecord(t) {
		if _, err := jalankan(ctx, tx, q, "menghapus Loss Record", id); err != nil {
			return err
		}
	}
	idKepala, err := r.db.NomorBerikut(ctx, tx, sequenceLossRecord)
	if err != nil {
		return err
	}
	arg := []any{idKepala, id}
	for _, ke := range penunjukKepala(&l) {
		arg = append(arg, ikatDesimal(*ke))
	}
	if _, err := jalankan(ctx, tx, sqlSisipKepalaLossRecord(t.kepala), "menyisipkan "+TabelLossRecord, arg...); err != nil {
		return err
	}
	for i, b := range l.Baris {
		idBaris, err := r.db.NomorBerikut(ctx, tx, sequenceLossRecordList)
		if err != nil {
			return err
		}
		uid, err := uidAcak()
		if err != nil {
			return err
		}
		var tahunKe any
		if b.TahunKe > 0 {
			tahunKe = b.TahunKe
		}
		if _, err := jalankan(ctx, tx, sqlSisipBarisLossRecord(t.baris), "menyisipkan "+TabelLossRecordList, idBaris, idKepala, i+1, uid,
			db.KosongJadiNil(b.UwYear), tahunKe, db.KosongJadiNil(b.ClassOfBusiness), db.KosongJadiNil(b.DateOfLoss),
			db.KosongJadiNil(b.CauseOfLoss), db.KosongJadiNil(b.Currency), ikatDesimal(b.LossAmount), ikatDesimal(b.LossAmountIDR),
			ikatDesimal(b.KursIDR)); err != nil {
			return err
		}
	}
	return nil
}

// JabatanAkun - lihat PembacaJabatan.
func (r *LossRecordOracle) JabatanAkun(ctx context.Context, akunID string) (string, error) {
	t, err := r.db.Qualify(TabelLoginGo)
	if err != nil {
		return "", err
	}
	var jabatan sql.NullString
	err = r.db.QueryRowContext(ctx, sqlJabatanAkun(t), akunID).Scan(&jabatan)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: baca %s: %w", TabelLoginGo, err)
	}
	return jabatan.String, nil
}
