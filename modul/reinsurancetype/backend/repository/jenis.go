// Package repository membaca dan menulis `POOLDATA.REINSURANCETYPE`. Setiap nilai diikat (`:n`, go-ora mengikat menurut
// URUTAN kemunculan); nama objek lewat `Qualify`. Nol COMMIT - transaksi milik services.
//
// ⛔ Prosedur `PEGA_REINSURANCETYPE` TIDAK dipanggil; `M_REINSURANCETYPE` (JSON Pega) tidak disentuh. ID baru dibentuk
// seperti prosedur itu: `1` + `M_REINSURANCETYPE_SEQ` 4 digit (bukan `REINSURANCETYPE_SEQ`, yang tidak dipakai).
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/reinsurancetype/backend/models"
)

// Nama objek warisan.
const (
	Tabel    = "REINSURANCETYPE"
	Sequence = "M_REINSURANCETYPE_SEQ"
	// awalanID - awalan ID baru `PEGA_REINSURANCETYPE` (`vID := 1|| lpad(...)`); diikat, bukan literal SQL.
	awalanID = "1"
	kolom    = `ID, NOTE, TYPE, SOANOTE, CODE, FLAG, NOURUT, GROUPTYPE, USERID, TGLUPDATE`
)

// DaftarTabelDitulis - penjaga modul.
var DaftarTabelDitulis = []string{Tabel}

var (
	// ErrTidakAda - ID tidak ada.
	ErrTidakAda = errors.New("repository: reinsurance type tidak ada")
	// ErrBelumAda - tabel atau sequence warisan tidak ada di skema ini.
	ErrBelumAda = errors.New("repository: tabel atau sequence Reinsurance Type tidak ada di skema ini")
	// ErrKembar - ORA-00001 (indeks unik ID + NOTE).
	ErrKembar = errors.New("repository: ID dan nama sudah dipakai (indeks unik REINSURANCETYPE_INDEX1)")
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
	case strings.Contains(s, "ORA-00001"):
		return fmt.Errorf("%w: %v", ErrKembar, err)
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

type pemindai interface{ Scan(...any) error }

func pindai(p pemindai) (models.Jenis, error) {
	var v [10]sql.NullString
	tuju := make([]any, len(v))
	for i := range v {
		tuju[i] = &v[i]
	}
	if err := p.Scan(tuju...); err != nil {
		return models.Jenis{}, err
	}
	return models.Jenis{ID: v[0].String, Name: v[1].String, Type: v[2].String, SoaName: v[3].String, Code: v[4].String,
		Flag: v[5].String, NoUrut: v[6].String, GroupType: v[7].String, UserID: v[8].String, TglUpdate: v[9].String}, nil
}

func (g *Gudang) baca(ctx context.Context, tx *db.Tx, q string, args ...any) ([]models.Jenis, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.dari(tx).QueryContext(ctx, q, args...)
	if err != nil {
		return nil, bungkus(err, "membaca")
	}
	defer func() { _ = rows.Close() }()
	out := []models.Jenis{}
	for rows.Next() {
		j, err := pindai(rows)
		if err != nil {
			return nil, bungkus(err, "memindai")
		}
		out = append(out, j)
	}
	return out, bungkus(rows.Err(), "membaca")
}

// SqlDaftar - daftar bersaring kata (ID, nama, nama SOA, code), Type, dan Flag (nil = semua); urut Type lalu nama.
func SqlDaftar(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s
	  WHERE (:1 IS NULL OR UPPER(ID) LIKE :2 ESCAPE '\' OR UPPER(NOTE) LIKE :3 ESCAPE '\'
	      OR UPPER(SOANOTE) LIKE :4 ESCAPE '\' OR UPPER(CODE) LIKE :5 ESCAPE '\')
	    AND (:6 IS NULL OR TYPE = :7)
	    AND (:8 IS NULL OR FLAG = :9)
	  ORDER BY TYPE, UPPER(NOTE), ID`, kolom, t)
}

// Daftar - baris bersaring; tipe / flag "" = semua.
func (g *Gudang) Daftar(ctx context.Context, kata, tipe, flag string) ([]models.Jenis, error) {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return nil, err
	}
	pola, ty, fl := PolaCari(kata), db.KosongJadiNil(strings.TrimSpace(tipe)), db.KosongJadiNil(strings.TrimSpace(flag))
	return g.baca(ctx, nil, SqlDaftar(t), pola, pola, pola, pola, pola, ty, ty, fl, fl)
}

// Ambil - satu baris.
func (g *Gudang) Ambil(ctx context.Context, tx *db.Tx, id string) (models.Jenis, error) {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return models.Jenis{}, err
	}
	d, err := g.baca(ctx, tx, fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, kolom, t), id)
	if err != nil {
		return models.Jenis{}, err
	}
	if len(d) == 0 {
		return models.Jenis{}, ErrTidakAda
	}
	return d[0], nil
}

// SqlPemakaiNama - baris lain yang NOTE-nya sama (tanpa beda huruf dan spasi tepi), selain kecualiID.
func SqlPemakaiNama(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE UPPER(TRIM(NOTE)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0)) ORDER BY ID`, kolom, t)
}

// PemakaiNama - baris lain bernama sama.
func (g *Gudang) PemakaiNama(ctx context.Context, tx *db.Tx, nama, kecualiID string) ([]models.Jenis, error) {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return nil, err
	}
	return g.baca(ctx, tx, SqlPemakaiNama(t), nama, db.KosongJadiNil(kecualiID))
}

// SqlIDBaru - ID baru seperti `PEGA_REINSURANCETYPE`: awalan `1` + nomor `M_REINSURANCETYPE_SEQ` 4 digit.
func SqlIDBaru(seq string) string {
	return fmt.Sprintf(`SELECT :1 || LPAD(TO_CHAR(%s.NEXTVAL), 4, '0') FROM DUAL`, seq)
}

// IDBaru - satu ID baru (sequence maju satu).
func (g *Gudang) IDBaru(ctx context.Context, tx *db.Tx) (string, error) {
	seq, err := g.db.Qualify(Sequence)
	if err != nil {
		return "", err
	}
	q := SqlIDBaru(seq)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var id sql.NullString
	if err := g.dari(tx).QueryRowContext(ctx, q, awalanID).Scan(&id); err != nil {
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

// SqlSisip - kolom yang ditulis `PEGA_REINSURANCETYPE` + GROUPTYPE (isian form).
func SqlSisip(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, NOTE, TYPE, SOANOTE, CODE, FLAG, USERID, NOURUT, TGLUPDATE, GROUPTYPE)
	  VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10)`, t)
}

// Sisip - baris baru.
func (g *Gudang) Sisip(ctx context.Context, tx *db.Tx, j models.Jenis) error {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return err
	}
	n := db.KosongJadiNil
	_, err = g.tulis(ctx, tx, SqlSisip(t), "menyimpan", j.ID, n(j.Name), n(j.Type), n(j.SoaName), n(j.Code), n(j.Flag),
		n(j.UserID), n(j.NoUrut), n(j.TglUpdate), n(j.GroupType))
	return err
}

// SqlUbah - kolom `PEGA_REINSURANCETYPE` + GROUPTYPE; ID tidak disentuh.
func SqlUbah(t string) string {
	return fmt.Sprintf(`UPDATE %s SET NOTE = :1, TYPE = :2, SOANOTE = :3, CODE = :4, FLAG = :5, USERID = :6, NOURUT = :7,
	  TGLUPDATE = :8, GROUPTYPE = :9 WHERE ID = :10`, t)
}

// Ubah - isian form; ID tetap.
func (g *Gudang) Ubah(ctx context.Context, tx *db.Tx, j models.Jenis) error {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return err
	}
	n := db.KosongJadiNil
	jumlah, err := g.tulis(ctx, tx, SqlUbah(t), "mengubah", n(j.Name), n(j.Type), n(j.SoaName), n(j.Code), n(j.Flag),
		n(j.UserID), n(j.NoUrut), n(j.TglUpdate), n(j.GroupType), j.ID)
	if err == nil && jumlah == 0 {
		return ErrTidakAda
	}
	return err
}
