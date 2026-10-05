// Package repository membaca dan menulis `POOLDATA.TREATYEXCHANGEYEARLY`; view `CURRENCY` dan `M_SITE_DATABASE` dibaca
// saja. Setiap nilai diikat (`:n`, go-ora mengikat menurut URUTAN kemunculan); nama objek lewat `Qualify`. Nol COMMIT -
// transaksi milik services.
//
// ⛔ Prosedur `PEGA_TREATYEXCHANGE` TIDAK dipanggil; `M_TREATYEXCHANGE` (JSON Pega) tidak disentuh. ID baru dibentuk
// seperti prosedur itu: situs aktif + `TREATYEXCHANGE_SEQ` 4 digit. ID warisan tidak unik, jadi baris dikenali lewat
// ROWID (`Kunci`).
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyexchangeyearly/backend/models"
)

// Nama objek warisan.
const (
	Tabel         = "TREATYEXCHANGEYEARLY"
	ViewMataUang  = "CURRENCY"
	TabelSitus    = "M_SITE_DATABASE"
	Sequence      = "TREATYEXCHANGE_SEQ"
	situsAktif    = "1" // nilai `M_SITE_DATABASE.CURRENT_SITE` situs yang sedang dipakai; diikat, bukan literal SQL
	kolomKursBaku = `ROWIDTOCHAR(t.ROWID), t.ID, t.TREATYYEAR, t.IDCURRENCY, t.CURRENCY, c.NOTE, t.STARTDATE, t.ENDDATE,
  t.TOIDR, t.TOUSD, t.QUARTER, t.USERID, t.DATEIU, t.DATEIN`
)

// DaftarTabelDitulis dan DaftarTabelDibacaSaja - penjaga modul.
var (
	DaftarTabelDitulis    = []string{Tabel}
	DaftarTabelDibacaSaja = []string{ViewMataUang, TabelSitus}
)

var (
	// ErrTidakAda - baris tidak ada (juga ROWID yang tidak sah).
	ErrTidakAda = errors.New("repository: kurs tidak ada")
	// ErrBelumAda - tabel, view, atau sequence warisan tidak ada di skema ini.
	ErrBelumAda = errors.New("repository: tabel, view, atau sequence Treaty Exchange Yearly tidak ada di skema ini")
	// ErrTanpaSitus - `M_SITE_DATABASE` tidak punya situs aktif (awalan ID).
	ErrTanpaSitus = errors.New("repository: M_SITE_DATABASE tanpa situs aktif (CURRENT_SITE)")
)

// Gudang - akses Oracle modul ini.
type Gudang struct{ db *db.DB }

// Baru membuat gudang.
func Baru(d *db.DB) *Gudang { return &Gudang{db: d} }

type penjalan interface {
	QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, q string, args ...any) *sql.Row
	ExecContext(ctx context.Context, q string, args ...any) (sql.Result, error)
}

func (g *Gudang) dari(tx *db.Tx) penjalan {
	if tx.Terisi() {
		return tx
	}
	return g.db
}

func bungkus(err error, apa string) error {
	if err == nil {
		return nil
	}
	s := err.Error()
	switch {
	case strings.Contains(s, "ORA-01410"):
		return ErrTidakAda
	case strings.Contains(s, "ORA-00942") || strings.Contains(s, "ORA-00904") || strings.Contains(s, "ORA-02289"):
		return fmt.Errorf("%w: %v", ErrBelumAda, err)
	}
	return fmt.Errorf("repository: %s: %w", apa, err)
}

// PolaCari - pola LIKE ber-ESCAPE '\' untuk cari "memuat", tanpa beda huruf; kosong = nil.
func PolaCari(kata string) any {
	kata = strings.ToUpper(strings.TrimSpace(kata))
	if kata == "" {
		return nil
	}
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return "%" + r.Replace(kata) + "%"
}

func (g *Gudang) nama(objek string) (string, error) { return g.db.Qualify(objek) }

type pemindai interface{ Scan(...any) error }

func baca[T any](ctx context.Context, j penjalan, q string, pindai func(pemindai) (T, error), args ...any) ([]T, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := j.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, bungkus(err, "membaca")
	}
	defer func() { _ = rows.Close() }()
	out := []T{}
	for rows.Next() {
		v, err := pindai(rows)
		if err != nil {
			return nil, bungkus(err, "memindai")
		}
		out = append(out, v)
	}
	return out, bungkus(rows.Err(), "membaca")
}

func teks(p pemindai, n int) ([]string, error) {
	v := make([]sql.NullString, n)
	tuju := make([]any, n)
	for i := range v {
		tuju[i] = &v[i]
	}
	if err := p.Scan(tuju...); err != nil {
		return nil, err
	}
	out := make([]string, n)
	for i := range v {
		out[i] = v[i].String
	}
	return out, nil
}

func pindaiKurs(p pemindai) (models.Kurs, error) {
	v, err := teks(p, 14)
	if err != nil {
		return models.Kurs{}, err
	}
	return models.Kurs{Kunci: v[0], ID: v[1], TreatyYear: v[2], IDCurrency: v[3], Currency: v[4], CurrencyName: v[5],
		StartDate: v[6], EndDate: v[7], ToIDR: v[8], ToUSD: v[9], Quarter: v[10], UserID: v[11], DateIU: v[12], DateIn: v[13]}, nil
}

func (g *Gudang) sumber() (string, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return "", err
	}
	c, err := g.nama(ViewMataUang)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`%s t LEFT JOIN %s c ON c.ID = t.IDCURRENCY`, t, c), nil
}

// SqlDaftar - daftar bersaring kata (ID, kode dan nama mata uang) dan tahun treaty (nil = semua); urut tahun terbaru.
func SqlDaftar(sumber string) string {
	return fmt.Sprintf(`SELECT %s FROM %s
	  WHERE (:1 IS NULL OR UPPER(t.ID) LIKE :2 ESCAPE '\' OR UPPER(t.CURRENCY) LIKE :3 ESCAPE '\' OR UPPER(c.NOTE) LIKE :4 ESCAPE '\')
	    AND (:5 IS NULL OR t.TREATYYEAR = :6)
	  ORDER BY t.TREATYYEAR DESC, t.CURRENCY, t.QUARTER, t.ID`, kolomKursBaku, sumber)
}

// Daftar - baris bersaring; tahun "" = semua tahun.
func (g *Gudang) Daftar(ctx context.Context, kata, tahun string) ([]models.Kurs, error) {
	s, err := g.sumber()
	if err != nil {
		return nil, err
	}
	pola, th := PolaCari(kata), db.KosongJadiNil(strings.TrimSpace(tahun))
	return baca(ctx, g.db, SqlDaftar(s), pindaiKurs, pola, pola, pola, pola, th, th)
}

func (g *Gudang) satu(ctx context.Context, tx *db.Tx, where string, arg any) (models.Kurs, error) {
	s, err := g.sumber()
	if err != nil {
		return models.Kurs{}, err
	}
	d, err := baca(ctx, g.dari(tx), fmt.Sprintf(`SELECT %s FROM %s WHERE %s`, kolomKursBaku, s, where), pindaiKurs, arg)
	if err != nil {
		return models.Kurs{}, err
	}
	if len(d) == 0 {
		return models.Kurs{}, ErrTidakAda
	}
	return d[0], nil
}

// AmbilKunci - satu baris menurut ROWID-nya.
func (g *Gudang) AmbilKunci(ctx context.Context, tx *db.Tx, kunci string) (models.Kurs, error) {
	return g.satu(ctx, tx, `t.ROWID = CHARTOROWID(:1)`, kunci)
}

// AmbilID - baris pertama ber-ID ini (dipakai sesudah Add: ID baru unik).
func (g *Gudang) AmbilID(ctx context.Context, tx *db.Tx, id string) (models.Kurs, error) {
	return g.satu(ctx, tx, `t.ID = :1`, id)
}

// SqlKembar - baris lain dengan tahun, mata uang, dan quarter yang sama; :4/:5 = ROWID baris sendiri (nil = Add).
func SqlKembar(t string) string {
	return fmt.Sprintf(`SELECT ID FROM %s WHERE TREATYYEAR = :1 AND IDCURRENCY = :2 AND QUARTER = :3
	  AND (:4 IS NULL OR ROWID <> CHARTOROWID(:5)) ORDER BY ID`, t)
}

// Kembar - ID baris lain dengan Treaty Year, Currency, dan Quarter yang sama, selain baris kecualiKunci.
func (g *Gudang) Kembar(ctx context.Context, tx *db.Tx, tahun, idCurrency, quarter, kecualiKunci string) ([]string, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return nil, err
	}
	k := db.KosongJadiNil(kecualiKunci)
	return baca(ctx, g.dari(tx), SqlKembar(t), func(p pemindai) (string, error) {
		v, err := teks(p, 1)
		if err != nil {
			return "", err
		}
		return v[0], nil
	}, tahun, idCurrency, quarter, k, k)
}

func pindaiMataUang(p pemindai) (models.MataUang, error) {
	v, err := teks(p, 3)
	if err != nil {
		return models.MataUang{}, err
	}
	return models.MataUang{ID: v[0], Kode: v[1], Nama: v[2]}, nil
}

// DaftarMataUang - pilihan Currency (view CURRENCY), urut kode.
func (g *Gudang) DaftarMataUang(ctx context.Context) ([]models.MataUang, error) {
	c, err := g.nama(ViewMataUang)
	if err != nil {
		return nil, err
	}
	return baca(ctx, g.db, fmt.Sprintf(`SELECT ID, CURRENCY, NOTE FROM %s ORDER BY CURRENCY, ID`, c), pindaiMataUang)
}

// AmbilMataUang - satu mata uang (kodenya disalin ke kolom CURRENCY).
func (g *Gudang) AmbilMataUang(ctx context.Context, tx *db.Tx, id string) (models.MataUang, error) {
	c, err := g.nama(ViewMataUang)
	if err != nil {
		return models.MataUang{}, err
	}
	d, err := baca(ctx, g.dari(tx), fmt.Sprintf(`SELECT ID, CURRENCY, NOTE FROM %s WHERE ID = :1`, c), pindaiMataUang, id)
	if err != nil {
		return models.MataUang{}, err
	}
	if len(d) == 0 {
		return models.MataUang{}, ErrTidakAda
	}
	return d[0], nil
}

// DaftarTahun - tahun treaty yang ada (saringan daftar), terbaru dulu.
func (g *Gudang) DaftarTahun(ctx context.Context) ([]string, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return nil, err
	}
	return baca(ctx, g.db, fmt.Sprintf(`SELECT DISTINCT TREATYYEAR FROM %s WHERE TREATYYEAR IS NOT NULL ORDER BY TREATYYEAR DESC`, t),
		func(p pemindai) (string, error) {
			v, err := teks(p, 1)
			if err != nil {
				return "", err
			}
			return v[0], nil
		})
}

// SqlIDBaru - ID baru seperti `PEGA_TREATYEXCHANGE`: ID situs aktif `M_SITE_DATABASE` + nomor sequence 4 digit.
func SqlIDBaru(situs, seq string) string {
	return fmt.Sprintf(`SELECT ID || LPAD(TO_CHAR(%s.NEXTVAL), 4, '0') FROM %s WHERE CURRENT_SITE = :1`, seq, situs)
}

// IDBaru - satu ID baru (sequence maju satu).
func (g *Gudang) IDBaru(ctx context.Context, tx *db.Tx) (string, error) {
	s, err := g.nama(TabelSitus)
	if err != nil {
		return "", err
	}
	seq, err := g.nama(Sequence)
	if err != nil {
		return "", err
	}
	q := SqlIDBaru(s, seq)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var id sql.NullString
	err = g.dari(tx).QueryRowContext(ctx, q, situsAktif).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) || (err == nil && !id.Valid) {
		return "", ErrTanpaSitus
	}
	if err != nil {
		return "", bungkus(err, "membuat ID")
	}
	return id.String, nil
}

// AdaID - ID sudah terpakai?
func (g *Gudang) AdaID(ctx context.Context, tx *db.Tx, id string) (bool, error) {
	_, err := g.AmbilID(ctx, tx, id)
	if errors.Is(err, ErrTidakAda) {
		return false, nil
	}
	return err == nil, err
}

func (g *Gudang) tulis(ctx context.Context, tx *db.Tx, q, apa string, args ...any) (int64, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	h, err := g.dari(tx).ExecContext(ctx, q, args...)
	if err != nil {
		return 0, bungkus(err, apa)
	}
	n, _ := h.RowsAffected()
	return n, nil
}

// SqlSisip - kolom yang ditulis `PEGA_TREATYEXCHANGE`, tanpa QURRENCYID (selalu kosong di DEV).
func SqlSisip(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, TREATYYEAR, STARTDATE, ENDDATE, TOIDR, TOUSD, USERID, DATEIU, IDCURRENCY, CURRENCY,
	  QUARTER, DATEIN) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12)`, t)
}

// Sisip - baris baru.
func (g *Gudang) Sisip(ctx context.Context, tx *db.Tx, k models.Kurs) error {
	t, err := g.nama(Tabel)
	if err != nil {
		return err
	}
	n := db.KosongJadiNil
	_, err = g.tulis(ctx, tx, SqlSisip(t), "menyimpan", k.ID, n(k.TreatyYear), n(k.StartDate), n(k.EndDate), n(k.ToIDR),
		n(k.ToUSD), n(k.UserID), n(k.DateIU), n(k.IDCurrency), n(k.Currency), n(k.Quarter), n(k.DateIn))
	return err
}

// SqlUbah - satu baris menurut ROWID; ID, QURRENCYID, dan DATEIN tidak disentuh.
func SqlUbah(t string) string {
	return fmt.Sprintf(`UPDATE %s SET TREATYYEAR = :1, STARTDATE = :2, ENDDATE = :3, TOIDR = :4, TOUSD = :5, USERID = :6,
	  DATEIU = :7, IDCURRENCY = :8, CURRENCY = :9, QUARTER = :10 WHERE ROWID = CHARTOROWID(:11)`, t)
}

// Ubah - isian form atas satu baris (Kunci).
func (g *Gudang) Ubah(ctx context.Context, tx *db.Tx, k models.Kurs) error {
	t, err := g.nama(Tabel)
	if err != nil {
		return err
	}
	n := db.KosongJadiNil
	jumlah, err := g.tulis(ctx, tx, SqlUbah(t), "mengubah", n(k.TreatyYear), n(k.StartDate), n(k.EndDate), n(k.ToIDR),
		n(k.ToUSD), n(k.UserID), n(k.DateIU), n(k.IDCurrency), n(k.Currency), n(k.Quarter), k.Kunci)
	if err == nil && jumlah == 0 {
		return ErrTidakAda
	}
	return err
}
