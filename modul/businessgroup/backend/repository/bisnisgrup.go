// Package repository membaca dan menulis `POOLDATA.BUSINESSGROUP`; `TREATYGROUP` dan `M_SITE_DATABASE` dibaca saja.
// Setiap nilai diikat (`:n`, go-ora mengikat menurut URUTAN kemunculan); nama objek lewat `Qualify`. Nol COMMIT -
// transaksi milik services.
//
// ⛔ Prosedur `UPSERT_BUSINESSGROUP` TIDAK dipanggil; `M_BUSINESSGROUP` (JSON Pega) tidak disentuh. ID baru dibentuk
// seperti prosedur itu: situs aktif + `BUSINESSGROUP_SEQ` 4 digit (DEV: `1` + 4 digit).
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/businessgroup/backend/models"
)

// Nama objek warisan.
const (
	Tabel        = "BUSINESSGROUP"
	TabelTreaty  = "TREATYGROUP"
	TabelSitus   = "M_SITE_DATABASE"
	Sequence     = "BUSINESSGROUP_SEQ"
	situsAktif   = "1" // nilai `M_SITE_DATABASE.CURRENT_SITE` situs yang sedang dipakai; diikat, bukan literal SQL
	bukanSyariah = `NVL(UPPER(TRIM(NOTE)), ' ') NOT LIKE '%` + models.AkhiranSyariah + `'`
)

// DaftarTabelDitulis dan DaftarTabelDibacaSaja - penjaga modul.
var (
	DaftarTabelDitulis    = []string{Tabel}
	DaftarTabelDibacaSaja = []string{TabelTreaty, TabelSitus}
)

var (
	// ErrTidakAda - ID tidak ada.
	ErrTidakAda = errors.New("repository: business group tidak ada")
	// ErrBelumAda - tabel atau sequence warisan tidak ada di skema ini.
	ErrBelumAda = errors.New("repository: tabel atau sequence Business Group tidak ada di skema ini")
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
	if s := err.Error(); strings.Contains(s, "ORA-00942") || strings.Contains(s, "ORA-00904") || strings.Contains(s, "ORA-02289") {
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

const kolom = `ID, NOTE, ALIASNAME, TOPID, TREATYNAME`

func pindai(p pemindai) (models.BisnisGrup, error) {
	v, err := teks(p, 5)
	if err != nil {
		return models.BisnisGrup{}, err
	}
	return models.BisnisGrup{ID: v[0], Name: v[1], Alias: v[2], TopID: v[3], TreatyName: v[4]}, nil
}

// SqlDaftar - daftar tanpa SYARIAH, bersaring kata (ID, nama, alias, nama treaty) dan Treaty Group (nil = semua).
func SqlDaftar(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s
	  WHERE %s
	    AND (:1 IS NULL OR UPPER(ID) LIKE :2 ESCAPE '\' OR UPPER(NOTE) LIKE :3 ESCAPE '\'
	      OR UPPER(ALIASNAME) LIKE :4 ESCAPE '\' OR UPPER(TREATYNAME) LIKE :5 ESCAPE '\')
	    AND (:6 IS NULL OR TOPID = :7)
	  ORDER BY UPPER(NOTE), ID`, kolom, t, bukanSyariah)
}

// Daftar - baris bersaring; topID "" = semua Treaty Group.
func (g *Gudang) Daftar(ctx context.Context, kata, topID string) ([]models.BisnisGrup, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return nil, err
	}
	pola, top := PolaCari(kata), db.KosongJadiNil(strings.TrimSpace(topID))
	return baca(ctx, g.db, SqlDaftar(t), pindai, pola, pola, pola, pola, pola, top, top)
}

// Ambil - satu baris (juga yang berakhiran SYARIAH; services yang menolaknya).
func (g *Gudang) Ambil(ctx context.Context, tx *db.Tx, id string) (models.BisnisGrup, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return models.BisnisGrup{}, err
	}
	d, err := baca(ctx, g.dari(tx), fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, kolom, t), pindai, id)
	if err != nil {
		return models.BisnisGrup{}, err
	}
	if len(d) == 0 {
		return models.BisnisGrup{}, ErrTidakAda
	}
	return d[0], nil
}

// PemakaiNama - ID baris lain yang NOTE-nya sama (tanpa beda huruf dan spasi tepi), selain kecualiID.
func (g *Gudang) PemakaiNama(ctx context.Context, tx *db.Tx, nama, kecualiID string) ([]string, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT ID FROM %s WHERE UPPER(TRIM(NOTE)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0)) ORDER BY ID`, t)
	return baca(ctx, g.dari(tx), q, func(p pemindai) (string, error) {
		v, err := teks(p, 1)
		if err != nil {
			return "", err
		}
		return v[0], nil
	}, nama, db.KosongJadiNil(kecualiID))
}

func pindaiTreaty(p pemindai) (models.TreatyGroup, error) {
	v, err := teks(p, 2)
	if err != nil {
		return models.TreatyGroup{}, err
	}
	return models.TreatyGroup{ID: v[0], Name: v[1]}, nil
}

// SqlDaftarTreaty - pilihan Treaty Group, urut nama.
func SqlDaftarTreaty(t string) string {
	return fmt.Sprintf(`SELECT ID, TREATYGROUPNAME FROM %s ORDER BY UPPER(TREATYGROUPNAME), ID`, t)
}

// DaftarTreaty - seluruh baris TREATYGROUP.
func (g *Gudang) DaftarTreaty(ctx context.Context) ([]models.TreatyGroup, error) {
	t, err := g.nama(TabelTreaty)
	if err != nil {
		return nil, err
	}
	return baca(ctx, g.db, SqlDaftarTreaty(t), pindaiTreaty)
}

// AmbilTreaty - satu baris TREATYGROUP (namanya disalin ke TREATYNAME saat disimpan).
func (g *Gudang) AmbilTreaty(ctx context.Context, tx *db.Tx, id string) (models.TreatyGroup, error) {
	t, err := g.nama(TabelTreaty)
	if err != nil {
		return models.TreatyGroup{}, err
	}
	d, err := baca(ctx, g.dari(tx), fmt.Sprintf(`SELECT ID, TREATYGROUPNAME FROM %s WHERE ID = :1`, t), pindaiTreaty, id)
	if err != nil {
		return models.TreatyGroup{}, err
	}
	if len(d) == 0 {
		return models.TreatyGroup{}, ErrTidakAda
	}
	return d[0], nil
}

// SqlIDBaru - ID baru seperti `UPSERT_BUSINESSGROUP`: ID situs aktif `M_SITE_DATABASE` + nomor sequence 4 digit.
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
	_, err := g.Ambil(ctx, tx, id)
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

// Sisip - baris baru; kolom sama dengan `UPSERT_BUSINESSGROUP`.
func (g *Gudang) Sisip(ctx context.Context, tx *db.Tx, b models.BisnisGrup) error {
	t, err := g.nama(Tabel)
	if err != nil {
		return err
	}
	n := db.KosongJadiNil
	_, err = g.tulis(ctx, tx, fmt.Sprintf(`INSERT INTO %s (ID, NOTE, TOPID, ALIASNAME, TREATYNAME) VALUES (:1, :2, :3, :4, :5)`, t),
		"menyimpan", b.ID, n(b.Name), n(b.TopID), n(b.Alias), n(b.TreatyName))
	return err
}

// Ubah - isian form; ID tetap.
func (g *Gudang) Ubah(ctx context.Context, tx *db.Tx, b models.BisnisGrup) error {
	t, err := g.nama(Tabel)
	if err != nil {
		return err
	}
	n := db.KosongJadiNil
	jumlah, err := g.tulis(ctx, tx, fmt.Sprintf(`UPDATE %s SET NOTE = :1, TOPID = :2, ALIASNAME = :3, TREATYNAME = :4 WHERE ID = :5`, t),
		"mengubah", n(b.Name), n(b.TopID), n(b.Alias), n(b.TreatyName), b.ID)
	if err == nil && jumlah == 0 {
		return ErrTidakAda
	}
	return err
}
