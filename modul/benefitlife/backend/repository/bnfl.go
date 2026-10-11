package repository

// SQL modul Benefit. Dibaca DAN ditulis di kolom `BENEFIT_LIFE` (keputusan work owner 08-10-2026 K1, migrasi inti
// 942-944) - kolom bernama, nol JSON. Nol DELETE (XML `InboxBenefit` tanpa Delete: `pyGridDeleteActivityExists` false
// b3292).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/benefitlife/backend/models"
)

var (
	// ErrTidakAda - ID tidak ada.
	ErrTidakAda = errors.New("repository: Benefit tidak ada")
	// ErrBelumAda - tabel, kolom, atau sequence tidak ada di skema ini (ORA-00942 / ORA-00904 / ORA-02289).
	ErrBelumAda = errors.New("repository: tabel atau sequence Benefit tidak ada di skema ini")
	// ErrKembar - ORA-00001 (PK `SYS_C009031`): ID sudah dipakai.
	ErrKembar = errors.New("repository: ID sudah dipakai")
	// ErrBacaSaja - SQL tulis diarahkan ke objek di luar DaftarTabelDitulis.
	ErrBacaSaja = errors.New("repository: objek ini dibaca saja")
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

// Bungkus memetakan galat Oracle ke galat paket ini (ORA-00001 = ErrKembar; tabel / kolom / sequence tidak ada =
// ErrBelumAda).
func Bungkus(err error, apa string) error {
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

// PeriksaTulis - lapis penjaga: pernyataan bukan SELECT hanya boleh atas objek DaftarTabelDitulis.
func PeriksaTulis(objek, q string) error {
	if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(q)), "SELECT ") {
		return nil
	}
	for _, t := range DaftarTabelDitulis {
		if t == objek {
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrBacaSaja, objek)
}

func siap(objek, q string) error {
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	return PeriksaTulis(objek, q)
}

// nama - nama berskema tabel dan sequence.
func (g *Gudang) nama() (tabel, seq string, err error) {
	if tabel, err = g.db.Qualify(Tabel); err != nil {
		return "", "", err
	}
	if seq, err = g.db.Qualify(Seq); err != nil {
		return "", "", err
	}
	return tabel, seq, nil
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

// Urutan - klausa ORDER BY grid: ID ANGKA (bukan urut teks) lalu ID. Bawaan MENURUN (`pySortType` DESC b4391); kolom
// Benefit tidak dapat diurutkan (`pyColumnSorting` false b4415) - satu-satunya pilihan adalah arah.
func Urutan(naik bool) string {
	arah := "DESC"
	if naik {
		arah = "ASC"
	}
	return fmt.Sprintf("TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) %s NULLS LAST, ID %s", arah, arah)
}

const saring = `(:1 IS NULL OR UPPER(ID) LIKE :2 ESCAPE '\') AND (:3 IS NULL OR UPPER(BENEFIT) LIKE :4 ESCAPE '\')`

// SqlDaftar - satu halaman grid `BrowseBenefitLife_RD` (b4451).
func SqlDaftar(t string, naik bool) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE %s ORDER BY %s OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY`,
		kolomBaca, t, saring, Urutan(naik))
}

// SqlJumlah - jumlah baris bersaring.
func SqlJumlah(t string) string { return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, t, saring) }

// SqlAmbil - satu baris.
func SqlAmbil(t string) string { return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, kolomBaca, t) }

// SqlNomorBaru - satu nomor sequence sebagai teks.
func SqlNomorBaru(seq string) string { return fmt.Sprintf(`SELECT TO_CHAR(%s.NEXTVAL) FROM DUAL`, seq) }

// SqlAdaID - ID sudah terpakai?
func SqlAdaID(t string) string { return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ID = :1`, t) }

// SqlSisip - Save Add (`AddToList_Act` b1791).
func SqlSisip(t string) string { return fmt.Sprintf(`INSERT INTO %s (ID, BENEFIT) VALUES (:1, :2)`, t) }

// SqlUbah - Edit (`EditList_DT` b4243) lalu Save.
func SqlUbah(t string) string { return fmt.Sprintf(`UPDATE %s SET BENEFIT = :1 WHERE ID = :2`, t) }

func pindai(rows *sql.Rows) ([]models.Benefit, error) {
	var out []models.Benefit
	for rows.Next() {
		var id, benefit sql.NullString
		if err := rows.Scan(&id, &benefit); err != nil {
			return nil, err
		}
		out = append(out, models.Benefit{ID: id.String, Benefit: benefit.String})
	}
	return out, rows.Err()
}

func (g *Gudang) baca(ctx context.Context, tx *db.Tx, q string, args ...any) ([]models.Benefit, error) {
	if err := siap(Tabel, q); err != nil {
		return nil, err
	}
	rows, err := g.dari(tx).QueryContext(ctx, q, args...)
	if err != nil {
		return nil, Bungkus(err, "membaca")
	}
	defer func() { _ = rows.Close() }()
	out, err := pindai(rows)
	return out, Bungkus(err, "membaca")
}

func (g *Gudang) satuNilai(ctx context.Context, tx *db.Tx, objek, q string, args ...any) (string, error) {
	if err := siap(objek, q); err != nil {
		return "", err
	}
	var s sql.NullString
	if err := g.dari(tx).QueryRowContext(ctx, q, args...).Scan(&s); err != nil {
		return "", Bungkus(err, "membaca")
	}
	return s.String, nil
}

func (g *Gudang) tulis(ctx context.Context, tx *db.Tx, q, apa string, args ...any) (int64, error) {
	if err := siap(Tabel, q); err != nil {
		return 0, err
	}
	h, err := g.dari(tx).ExecContext(ctx, q, args...)
	if err != nil {
		return 0, Bungkus(err, apa)
	}
	n, _ := h.RowsAffected()
	return n, nil
}

func angka(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}

// Daftar - satu halaman bersaring dan total barisnya.
func (g *Gudang) Daftar(ctx context.Context, s models.Saringan) ([]models.Benefit, int, error) {
	t, _, err := g.nama()
	if err != nil {
		return nil, 0, err
	}
	id, bn := PolaCari(s.ID), PolaCari(s.Benefit)
	total, err := g.satuNilai(ctx, nil, Tabel, SqlJumlah(t), id, id, bn, bn)
	if err != nil {
		return nil, 0, err
	}
	d, err := g.baca(ctx, nil, SqlDaftar(t, s.Naik), id, id, bn, bn, (s.Halaman-1)*models.UkuranHalaman, models.UkuranHalaman)
	return d, angka(total), err
}

// Ambil - satu baris; ErrTidakAda bila ID tidak ada.
func (g *Gudang) Ambil(ctx context.Context, tx *db.Tx, id string) (models.Benefit, error) {
	t, _, err := g.nama()
	if err != nil {
		return models.Benefit{}, err
	}
	d, err := g.baca(ctx, tx, SqlAmbil(t), id)
	if err != nil {
		return models.Benefit{}, err
	}
	if len(d) == 0 {
		return models.Benefit{}, ErrTidakAda
	}
	return d[0], nil
}

// NomorBaru - satu nomor `M_BENEFIT_LIFE_SEQ` (teks).
func (g *Gudang) NomorBaru(ctx context.Context, tx *db.Tx) (string, error) {
	_, seq, err := g.nama()
	if err != nil {
		return "", err
	}
	return g.satuNilai(ctx, tx, "DUAL", SqlNomorBaru(seq))
}

// AdaID - ID sudah terpakai di BENEFIT_LIFE?
func (g *Gudang) AdaID(ctx context.Context, tx *db.Tx, id string) (bool, error) {
	t, _, err := g.nama()
	if err != nil {
		return false, err
	}
	s, err := g.satuNilai(ctx, tx, Tabel, SqlAdaID(t), id)
	return angka(s) > 0, err
}

// Sisip - baris baru; ErrKembar bila ID terpakai (ORA-00001).
func (g *Gudang) Sisip(ctx context.Context, tx *db.Tx, b models.Benefit) error {
	t, _, err := g.nama()
	if err != nil {
		return err
	}
	_, err = g.tulis(ctx, tx, SqlSisip(t), "menyimpan", b.ID, db.KosongJadiNil(b.Benefit))
	return err
}

// Ubah - BENEFIT satu baris; ErrTidakAda bila ID tidak ada.
func (g *Gudang) Ubah(ctx context.Context, tx *db.Tx, b models.Benefit) error {
	t, _, err := g.nama()
	if err != nil {
		return err
	}
	n, err := g.tulis(ctx, tx, SqlUbah(t), "mengubah", db.KosongJadiNil(b.Benefit), b.ID)
	if err == nil && n == 0 {
		return ErrTidakAda
	}
	return err
}
