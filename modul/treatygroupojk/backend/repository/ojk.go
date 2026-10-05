// Package repository membaca dan menulis `POOLDATA.TREATYGROUPOJK`. Setiap nilai diikat (`:n`, go-ora mengikat menurut
// URUTAN kemunculan); nama objek lewat `Qualify`. Nol COMMIT - transaksi milik services.
//
// Tabel warisan tanpa PK, tanpa sequence: ID baru = nomor ID tertinggi + 1 (dua digit) dan Order No baru = Order No
// tertinggi + 1, dihitung SESUDAH seluruh baris dikunci `FOR UPDATE` supaya dua Add bersamaan tidak mendapat nomor yang
// sama.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatygroupojk/backend/models"
)

// Tabel - nama objek warisan.
const Tabel = "TREATYGROUPOJK"

var (
	// ErrTidakAda - ID tidak ada.
	ErrTidakAda = errors.New("repository: treaty group OJK tidak ada")
	// ErrBelumAda - tabel TREATYGROUPOJK tidak ada di skema ini.
	ErrBelumAda = errors.New("repository: tabel TREATYGROUPOJK tidak ada di skema ini")
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
	if s := err.Error(); strings.Contains(s, "ORA-00942") || strings.Contains(s, "ORA-00904") {
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

const kolom = `ID, NAME, NAMEIDN, ORDERNO`

// urutan - Order No sebagai angka (teks warisan; LPAD supaya `9` sebelum `11`), lalu ID.
const urutan = `ORDER BY LPAD(ORDERNO, 10, '0'), ID`

type pemindai interface{ Scan(...any) error }

func pindai(p pemindai) (models.Ojk, error) {
	var v [4]sql.NullString
	if err := p.Scan(&v[0], &v[1], &v[2], &v[3]); err != nil {
		return models.Ojk{}, err
	}
	return models.Ojk{ID: v[0].String, Name: v[1].String, NameIDN: v[2].String, OrderNo: v[3].String}, nil
}

func (g *Gudang) daftar(ctx context.Context, tx *db.Tx, q string, args ...any) ([]models.Ojk, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.dari(tx).QueryContext(ctx, q, args...)
	if err != nil {
		return nil, bungkus(err, "membaca")
	}
	defer func() { _ = rows.Close() }()
	out := []models.Ojk{}
	for rows.Next() {
		o, err := pindai(rows)
		if err != nil {
			return nil, bungkus(err, "memindai")
		}
		out = append(out, o)
	}
	return out, bungkus(rows.Err(), "membaca")
}

// SqlDaftar - daftar bersaring kata (ID, Name, Name IDN), urut Order No.
func SqlDaftar(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s
	  WHERE (:1 IS NULL OR UPPER(ID) LIKE :2 ESCAPE '\' OR UPPER(NAME) LIKE :3 ESCAPE '\' OR UPPER(NAMEIDN) LIKE :4 ESCAPE '\')
	  %s`, kolom, t, urutan)
}

// Daftar - baris bersaring, urut Order No.
func (g *Gudang) Daftar(ctx context.Context, kata string) ([]models.Ojk, error) {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return nil, err
	}
	pola := PolaCari(kata)
	return g.daftar(ctx, nil, SqlDaftar(t), pola, pola, pola, pola)
}

// Ambil - satu baris.
func (g *Gudang) Ambil(ctx context.Context, tx *db.Tx, id string) (models.Ojk, error) {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return models.Ojk{}, err
	}
	d, err := g.daftar(ctx, tx, fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, kolom, t), id)
	if err != nil {
		return models.Ojk{}, err
	}
	if len(d) == 0 {
		return models.Ojk{}, ErrTidakAda
	}
	return d[0], nil
}

// SqlPemakai - baris lain yang kolomnya sama (tanpa beda huruf dan spasi tepi), selain kecualiID. `kolomNilai` hanya
// salah satu nama kolom tetap di bawah - bukan masukan pengguna.
func SqlPemakai(t, kolomNilai string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE UPPER(TRIM(%s)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0)) ORDER BY ID`,
		kolom, t, kolomNilai)
}

func (g *Gudang) pemakai(ctx context.Context, tx *db.Tx, kolomNilai, nilai, kecualiID string) ([]models.Ojk, error) {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return nil, err
	}
	return g.daftar(ctx, tx, SqlPemakai(t, kolomNilai), nilai, db.KosongJadiNil(kecualiID))
}

// PemakaiNama - baris lain yang NAME-nya sama.
func (g *Gudang) PemakaiNama(ctx context.Context, tx *db.Tx, nama, kecualiID string) ([]models.Ojk, error) {
	return g.pemakai(ctx, tx, "NAME", nama, kecualiID)
}

// SqlKunci - mengunci seluruh baris sampai transaksi selesai (Add bersamaan menunggu giliran).
func SqlKunci(t string) string { return fmt.Sprintf(`SELECT ID FROM %s FOR UPDATE`, t) }

// SqlNomorTertinggi - nomor ID tertinggi (hanya ID yang seluruhnya angka); tabel kosong = 0.
func SqlNomorTertinggi(t string) string {
	return fmt.Sprintf(`SELECT NVL(MAX(TO_NUMBER(ID)), 0) FROM %s WHERE REGEXP_LIKE(ID, '^[0-9]+$')`, t)
}

// SqlOrderTertinggi - Order No tertinggi (hanya ORDERNO yang seluruhnya angka); tabel kosong = 0.
func SqlOrderTertinggi(t string) string {
	return fmt.Sprintf(`SELECT NVL(MAX(TO_NUMBER(ORDERNO)), 0) FROM %s WHERE REGEXP_LIKE(ORDERNO, '^[0-9]+$')`, t)
}

// KunciNomorTertinggi - kunci seluruh baris, lalu nomor ID tertinggi dan Order No tertinggi. Dipanggil di DALAM
// transaksi Add.
func (g *Gudang) KunciNomorTertinggi(ctx context.Context, tx *db.Tx) (id, order int, err error) {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return 0, 0, err
	}
	kunci := SqlKunci(t)
	if err := db.PeriksaSQL(kunci); err != nil {
		return 0, 0, err
	}
	rows, err := g.dari(tx).QueryContext(ctx, kunci)
	if err != nil {
		return 0, 0, bungkus(err, "mengunci")
	}
	for rows.Next() {
		// Isinya tidak dipakai - barisnya cukup terkunci.
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return 0, 0, bungkus(err, "mengunci")
	}
	if err := rows.Close(); err != nil {
		return 0, 0, bungkus(err, "mengunci")
	}
	if err := g.dari(tx).QueryRowContext(ctx, SqlNomorTertinggi(t)).Scan(&id); err != nil {
		return 0, 0, bungkus(err, "membaca nomor tertinggi")
	}
	if err := g.dari(tx).QueryRowContext(ctx, SqlOrderTertinggi(t)).Scan(&order); err != nil {
		return 0, 0, bungkus(err, "membaca Order No tertinggi")
	}
	return id, order, nil
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

// Sisip - baris baru.
func (g *Gudang) Sisip(ctx context.Context, tx *db.Tx, o models.Ojk) error {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return err
	}
	_, err = g.tulis(ctx, tx, fmt.Sprintf(`INSERT INTO %s (ID, NAME, NAMEIDN, ORDERNO) VALUES (:1, :2, :3, :4)`, t),
		"menyimpan", o.ID, db.KosongJadiNil(o.Name), db.KosongJadiNil(o.NameIDN), db.KosongJadiNil(o.OrderNo))
	return err
}

// SqlUbah - isian form; ID dan ORDERNO tidak disentuh.
func SqlUbah(t string) string {
	return fmt.Sprintf(`UPDATE %s SET NAME = :1, NAMEIDN = :2 WHERE ID = :3`, t)
}

// Ubah - isian form; ID dan Order No tetap.
func (g *Gudang) Ubah(ctx context.Context, tx *db.Tx, o models.Ojk) error {
	t, err := g.db.Qualify(Tabel)
	if err != nil {
		return err
	}
	n, err := g.tulis(ctx, tx, SqlUbah(t), "mengubah", db.KosongJadiNil(o.Name), db.KosongJadiNil(o.NameIDN), o.ID)
	if err == nil && n == 0 {
		return ErrTidakAda
	}
	return err
}
