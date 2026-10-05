// Package repository membaca dan menulis `POOLDATA.ADJUSTERCONSULTANT`. Setiap nilai diikat (`:n`, go-ora mengikat
// menurut URUTAN kemunculan); nama objek lewat `Qualify`. Nol COMMIT - transaksi milik services.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/adjusterconsultant/backend/models"
)

// Nama objek warisan.
const (
	Tabel       = "ADJUSTERCONSULTANT"
	Sequence    = "ADJUSTERCONSULTANT_SEQ"
	TabelSitus  = "M_SITE_DATABASE"
	BatasDaftar = 2000
)

// situsAktif - nilai `M_SITE_DATABASE.CURRENT_SITE` situs yang sedang dipakai; diikat, bukan literal SQL.
const situsAktif = "1"

var (
	// ErrTidakAda - ID tidak ada.
	ErrTidakAda = errors.New("repository: adjuster consultant tidak ada")
	// ErrBelumDimigrasi - kolom IS_ACTIVE belum ada (migrasi 870 belum dijalankan) atau tabelnya tidak ada.
	ErrBelumDimigrasi = errors.New("repository: migrasi adjusterconsultant 870 belum dijalankan (-migrate)")
	// ErrIDTerpakai - ORA-00001 saat menyisipkan.
	ErrIDTerpakai = errors.New("repository: ID sudah terpakai")
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
	case strings.Contains(s, "ORA-00942"), strings.Contains(s, "ORA-00904"), strings.Contains(s, "ORA-02289"):
		return fmt.Errorf("%w: %v", ErrBelumDimigrasi, err)
	case strings.Contains(s, "ORA-00001"):
		return fmt.Errorf("%w: %v", ErrIDTerpakai, err)
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

const kolom = `ID, NAME, ADDRESS, TELPNO, USERNAME, TO_CHAR(EDITDATE, 'DD-MM-YYYY HH24:MI'), IS_ACTIVE`

type pemindai interface{ Scan(...any) error }

func pindai(p pemindai) (models.Adjuster, error) {
	var v [7]sql.NullString
	if err := p.Scan(&v[0], &v[1], &v[2], &v[3], &v[4], &v[5], &v[6]); err != nil {
		return models.Adjuster{}, err
	}
	return models.Adjuster{ID: v[0].String, Name: v[1].String, Address: v[2].String, TelpNo: v[3].String,
		Username: v[4].String, EditDate: v[5].String, Active: v[6].String == models.BenderaAktif}, nil
}

func (g *Gudang) daftar(ctx context.Context, tx *db.Tx, q string, args ...any) ([]models.Adjuster, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.dari(tx).QueryContext(ctx, q, args...)
	if err != nil {
		return nil, bungkus(err, "membaca")
	}
	defer func() { _ = rows.Close() }()
	out := []models.Adjuster{}
	for rows.Next() {
		a, err := pindai(rows)
		if err != nil {
			return nil, bungkus(err, "memindai")
		}
		out = append(out, a)
	}
	return out, bungkus(rows.Err(), "membaca")
}

// SqlDaftar - daftar bersaring: kata (ID, nama, alamat, telepon) dan bendera IS_ACTIVE (nil = semua).
func SqlDaftar(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s
	  WHERE (:1 IS NULL OR UPPER(ID) LIKE :2 ESCAPE '\' OR UPPER(NAME) LIKE :3 ESCAPE '\'
	    OR UPPER(ADDRESS) LIKE :4 ESCAPE '\' OR UPPER(TELPNO) LIKE :5 ESCAPE '\')
	    AND (:6 IS NULL OR IS_ACTIVE = :7)
	  ORDER BY UPPER(NAME), ID FETCH FIRST :8 ROWS ONLY`, kolom, t)
}

// Daftar - baris bersaring, urut nama; bendera "" = semua status.
func (g *Gudang) Daftar(ctx context.Context, kata, bendera string) ([]models.Adjuster, error) {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return nil, err
	}
	pola := PolaCari(kata)
	b := db.KosongJadiNil(bendera)
	return g.daftar(ctx, nil, SqlDaftar(t), pola, pola, pola, pola, pola, b, b, BatasDaftar)
}

// Ambil - satu baris.
func (g *Gudang) Ambil(ctx context.Context, tx *db.Tx, id string) (models.Adjuster, error) {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return models.Adjuster{}, err
	}
	d, err := g.daftar(ctx, tx, fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, kolom, t), id)
	if err != nil {
		return models.Adjuster{}, err
	}
	if len(d) == 0 {
		return models.Adjuster{}, ErrTidakAda
	}
	return d[0], nil
}

// PemakaiNama - baris lain yang NAME-nya sama (tanpa beda huruf dan spasi tepi), selain kecualiID.
func (g *Gudang) PemakaiNama(ctx context.Context, tx *db.Tx, nama, kecualiID string) ([]models.Adjuster, error) {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE UPPER(TRIM(NAME)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0)) ORDER BY ID`, kolom, t)
	return g.daftar(ctx, tx, q, nama, db.KosongJadiNil(kecualiID))
}

// SqlIDBaru - ID baru seperti `GetIDConsultanAdj_SQL` Pega: ID situs aktif `M_SITE_DATABASE` + nomor sequence 4 digit.
func SqlIDBaru(situs, seq string) string {
	return fmt.Sprintf(`SELECT ID || LPAD(TO_CHAR(%s.NEXTVAL), 4, '0') FROM %s WHERE CURRENT_SITE = :1`, seq, situs)
}

// IDBaru - satu ID baru (sequence maju satu).
func (g *Gudang) IDBaru(ctx context.Context, tx *db.Tx) (string, error) {
	s, err := g.db.Qualify(TabelSitus)
	if err != nil {
		return "", err
	}
	seq, err := g.db.Qualify(Sequence)
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

// Sisip - baris baru, aktif; EDITDATE = SYSDATE.
func (g *Gudang) Sisip(ctx context.Context, tx *db.Tx, a models.Adjuster, username string) error {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return err
	}
	_, err = g.tulis(ctx, tx, fmt.Sprintf(`INSERT INTO %s (ID, NAME, ADDRESS, TELPNO, USERNAME, EDITDATE, IS_ACTIVE)
	  VALUES (:1, :2, :3, :4, :5, SYSDATE, :6)`, t), "menyimpan", a.ID, db.KosongJadiNil(a.Name),
		db.KosongJadiNil(a.Address), db.KosongJadiNil(a.TelpNo), db.KosongJadiNil(username), models.BenderaAktif)
	return err
}

// Ubah - isian form; IS_ACTIVE tidak disentuh.
func (g *Gudang) Ubah(ctx context.Context, tx *db.Tx, a models.Adjuster, username string) error {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return err
	}
	n, err := g.tulis(ctx, tx, fmt.Sprintf(`UPDATE %s SET NAME = :1, ADDRESS = :2, TELPNO = :3, USERNAME = :4,
	  EDITDATE = SYSDATE WHERE ID = :5`, t), "mengubah", db.KosongJadiNil(a.Name), db.KosongJadiNil(a.Address),
		db.KosongJadiNil(a.TelpNo), db.KosongJadiNil(username), a.ID)
	if err == nil && n == 0 {
		return ErrTidakAda
	}
	return err
}

// SetelAktif - IS_ACTIVE, pengubah, dan waktu ubah.
func (g *Gudang) SetelAktif(ctx context.Context, tx *db.Tx, id string, aktif bool, username string) error {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return err
	}
	bendera := models.BenderaNonaktif
	if aktif {
		bendera = models.BenderaAktif
	}
	n, err := g.tulis(ctx, tx, fmt.Sprintf(`UPDATE %s SET IS_ACTIVE = :1, USERNAME = :2, EDITDATE = SYSDATE WHERE ID = :3`, t),
		"mengubah status", bendera, db.KosongJadiNil(username), id)
	if err == nil && n == 0 {
		return ErrTidakAda
	}
	return err
}
