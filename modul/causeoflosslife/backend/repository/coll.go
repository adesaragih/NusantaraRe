package repository

// SQL modul Cause Of Loss Life. Dibaca DAN ditulis di kolom `CAUSEOFLOSS_LIFE` (keputusan work owner 08-10-2026 K1,
// migrasi modul 090-092) - kolom bernama, nol JSON. Nol DELETE (XML `InboxCauseofLossLife` tanpa Delete:
// `pyGridDeleteActivityExists` false b2847).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/causeoflosslife/backend/models"
)

var (
	// ErrTidakAda - ID tidak ada.
	ErrTidakAda = errors.New("repository: Cause of Loss tidak ada")
	// ErrBelumAda - tabel, kolom, atau sequence tidak ada di skema ini (ORA-00942 / ORA-00904 / ORA-02289).
	ErrBelumAda = errors.New("repository: tabel atau sequence Cause of Loss tidak ada di skema ini")
	// ErrKembar - ORA-00001 (PK `SYS_C008825`): ID sudah dipakai.
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

// Urutan - klausa ORDER BY grid: ID ANGKA MENAIK lalu ID (`pySortType` ASC b3923, `pySortOrder` 1 b3929). Tidak ada
// pilihan pengguna: ketiga kolom `pyColumnSorting` false (b3925 / b3947 / b3969).
const Urutan = "TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) ASC NULLS LAST, ID ASC"

// SqlDaftar - satu halaman grid `BrowseCauseofLossLife_RD` (b3981). Tanpa saring (`pyGridFiltering` false b4036).
func SqlDaftar(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s ORDER BY %s OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY`, kolomBaca, t, Urutan)
}

// SqlJumlah - jumlah baris.
func SqlJumlah(t string) string { return fmt.Sprintf(`SELECT COUNT(*) FROM %s`, t) }

// SqlAmbil - satu baris.
func SqlAmbil(t string) string { return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, kolomBaca, t) }

// SqlPemakaiNama - baris lain bernama sama (tanpa beda huruf dan spasi tepi), selain kecualiID (K4, di luar XML).
// Baris bernama NULL (100001 DEV) tidak pernah cocok.
func SqlPemakaiNama(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE UPPER(TRIM(CAUSEOFLOSS)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0)) ORDER BY ID`,
		kolomBaca, t)
}

// SqlNomorBaru - satu nomor sequence sebagai teks.
func SqlNomorBaru(seq string) string { return fmt.Sprintf(`SELECT TO_CHAR(%s.NEXTVAL) FROM DUAL`, seq) }

// SqlAdaID - ID sudah terpakai?
func SqlAdaID(t string) string { return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ID = :1`, t) }

// SqlSisip - Save Add (`AddToList_Act` b5385).
func SqlSisip(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, CAUSEOFLOSS) VALUES (:1, :2)`, t)
}

// SqlUbah - Edit (`EditList_DT` b3770) lalu Save.
func SqlUbah(t string) string { return fmt.Sprintf(`UPDATE %s SET CAUSEOFLOSS = :1 WHERE ID = :2`, t) }

func pindai(rows *sql.Rows) ([]models.CauseOfLoss, error) {
	var out []models.CauseOfLoss
	for rows.Next() {
		var id, nama sql.NullString
		if err := rows.Scan(&id, &nama); err != nil {
			return nil, err
		}
		out = append(out, models.CauseOfLoss{ID: id.String, CauseOfLoss: nama.String})
	}
	return out, rows.Err()
}

func (g *Gudang) baca(ctx context.Context, tx *db.Tx, q string, args ...any) ([]models.CauseOfLoss, error) {
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

// Daftar - satu halaman dan total barisnya.
func (g *Gudang) Daftar(ctx context.Context, s models.Saringan) ([]models.CauseOfLoss, int, error) {
	t, _, err := g.nama()
	if err != nil {
		return nil, 0, err
	}
	total, err := g.satuNilai(ctx, nil, Tabel, SqlJumlah(t))
	if err != nil {
		return nil, 0, err
	}
	d, err := g.baca(ctx, nil, SqlDaftar(t), (s.Halaman-1)*models.UkuranHalaman, models.UkuranHalaman)
	return d, angka(total), err
}

// Ambil - satu baris; ErrTidakAda bila ID tidak ada.
func (g *Gudang) Ambil(ctx context.Context, tx *db.Tx, id string) (models.CauseOfLoss, error) {
	t, _, err := g.nama()
	if err != nil {
		return models.CauseOfLoss{}, err
	}
	d, err := g.baca(ctx, tx, SqlAmbil(t), id)
	if err != nil {
		return models.CauseOfLoss{}, err
	}
	if len(d) == 0 {
		return models.CauseOfLoss{}, ErrTidakAda
	}
	return d[0], nil
}

// PemakaiNama - baris bernama sama selain kecualiID ("" = semua).
func (g *Gudang) PemakaiNama(ctx context.Context, tx *db.Tx, nama, kecualiID string) ([]models.CauseOfLoss, error) {
	t, _, err := g.nama()
	if err != nil {
		return nil, err
	}
	return g.baca(ctx, tx, SqlPemakaiNama(t), nama, db.KosongJadiNil(kecualiID))
}

// NomorBaru - satu nomor `M_CAUSEOFLOSS_LIFE_SEQ` (teks).
func (g *Gudang) NomorBaru(ctx context.Context, tx *db.Tx) (string, error) {
	_, seq, err := g.nama()
	if err != nil {
		return "", err
	}
	return g.satuNilai(ctx, tx, "DUAL", SqlNomorBaru(seq))
}

// AdaID - ID sudah terpakai di CAUSEOFLOSS_LIFE?
func (g *Gudang) AdaID(ctx context.Context, tx *db.Tx, id string) (bool, error) {
	t, _, err := g.nama()
	if err != nil {
		return false, err
	}
	s, err := g.satuNilai(ctx, tx, Tabel, SqlAdaID(t), id)
	return angka(s) > 0, err
}

// Sisip - baris baru; ErrKembar bila ID terpakai (ORA-00001).
func (g *Gudang) Sisip(ctx context.Context, tx *db.Tx, c models.CauseOfLoss) error {
	t, _, err := g.nama()
	if err != nil {
		return err
	}
	_, err = g.tulis(ctx, tx, SqlSisip(t), "menyimpan", c.ID, db.KosongJadiNil(c.CauseOfLoss))
	return err
}

// Ubah - CAUSEOFLOSS satu baris; ErrTidakAda bila ID tidak ada.
func (g *Gudang) Ubah(ctx context.Context, tx *db.Tx, c models.CauseOfLoss) error {
	t, _, err := g.nama()
	if err != nil {
		return err
	}
	n, err := g.tulis(ctx, tx, SqlUbah(t), "mengubah", db.KosongJadiNil(c.CauseOfLoss), c.ID)
	if err == nil && n == 0 {
		return ErrTidakAda
	}
	return err
}
