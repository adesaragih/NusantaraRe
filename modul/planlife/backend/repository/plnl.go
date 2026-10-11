package repository

// SQL modul Plan. Dibaca DAN ditulis di kolom `PRODUCT_TYPE_LIFE` (keputusan work owner 08-10-2026 K1, migrasi inti
// 946-948) - kolom bernama, nol JSON. Pilihan Business / Benefit dibaca dari master `BUSINESS` / `BENEFIT_LIFE` (K4).
// Nol DELETE (XML `InboxProductType` tanpa Delete: `pyGridDeleteActivityExists` false b419 / b2839).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/planlife/backend/models"
)

var (
	// ErrTidakAda - ID tidak ada.
	ErrTidakAda = errors.New("repository: Plan tidak ada")
	// ErrBelumAda - tabel, kolom, atau sequence tidak ada di skema ini (ORA-00942 / ORA-00904 / ORA-02289).
	ErrBelumAda = errors.New("repository: tabel atau sequence Plan tidak ada di skema ini")
	// ErrKembar - ORA-00001 (PK_PRODUCT_TYPE_LIFE): ID sudah dipakai.
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

// Bungkus memetakan galat Oracle ke galat paket ini.
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

func (g *Gudang) nama(objek string) (string, error) { return g.db.Qualify(objek) }

// Urutan - klausa ORDER BY grid (kolom dari daftar putih). Kosong = ID ANGKA menaik (stabil).
func Urutan(urut string, turun bool) string {
	arah := "ASC"
	if turun {
		arah = "DESC"
	}
	switch urut {
	case models.UrutCoverName:
		return fmt.Sprintf("UPPER(COVERNAME) %s NULLS LAST, ID", arah)
	case models.UrutBenefit:
		return fmt.Sprintf("UPPER(BENEFIT) %s NULLS LAST, ID", arah)
	}
	return "TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) ASC NULLS LAST, ID ASC"
}

// SqlDaftar - satu halaman grid `BrowseProductTypeLife_RD` (b4325).
func SqlDaftar(t, urut string, turun bool) string {
	return fmt.Sprintf(`SELECT %s FROM %s ORDER BY %s OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY`, kolomBaca, t, Urutan(urut, turun))
}

// SqlJumlah - jumlah baris.
func SqlJumlah(t string) string { return fmt.Sprintf(`SELECT COUNT(*) FROM %s`, t) }

// SqlAmbil - satu baris.
func SqlAmbil(t string) string { return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, kolomBaca, t) }

// SqlPemakaiNama - plan lain bernama sama (tanpa beda huruf dan spasi tepi), selain kecualiID (K5).
func SqlPemakaiNama(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE UPPER(TRIM(COVERNAME)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0)) ORDER BY ID`,
		kolomBaca, t)
}

// SqlNomorBaru - satu nomor sequence sebagai teks.
func SqlNomorBaru(seq string) string { return fmt.Sprintf(`SELECT TO_CHAR(%s.NEXTVAL) FROM DUAL`, seq) }

// SqlAdaID - ID sudah terpakai?
func SqlAdaID(t string) string { return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ID = :1`, t) }

// SqlSisip - Save Add (`SaveProductTypeLife_Act` b6422).
func SqlSisip(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, COVERNAME, BUSINESS, BUSINESSID, BENEFIT, BENEFITID) VALUES (:1, :2, :3, :4, :5, :6)`, t)
}

// SqlUbah - Edit (`EditProductTypeLife_Act` b4045) lalu Save.
func SqlUbah(t string) string {
	return fmt.Sprintf(`UPDATE %s SET COVERNAME = :1, BUSINESS = :2, BUSINESSID = :3, BENEFIT = :4, BENEFITID = :5 WHERE ID = :6`, t)
}

// SqlPilihanBusiness - `BrowseBusinessLife_RD` (b1183): ID, OLDID, NOTE grup life (K4: GROUPPANEL = '009', diturunkan
// dari data), urut ID. ID / OLDID lewat TO_CHAR (tipe kolom master tidak tercatat di repo).
func SqlPilihanBusiness(t string) string {
	return fmt.Sprintf(`SELECT TO_CHAR(ID), TO_CHAR(OLDID), NOTE FROM %s WHERE GROUPPANEL = :1 ORDER BY ID`, t)
}

// SqlPilihanBenefit - `BrowseBenefitLife_RD` (b1537): ID (= `.Number`), BENEFIT, urut ID.
func SqlPilihanBenefit(t string) string {
	return fmt.Sprintf(`SELECT ID, BENEFIT FROM %s ORDER BY ID`, t)
}

func pindai(rows *sql.Rows, n int) ([][]string, error) {
	var out [][]string
	for rows.Next() {
		v := make([]sql.NullString, n)
		tuju := make([]any, n)
		for i := range v {
			tuju[i] = &v[i]
		}
		if err := rows.Scan(tuju...); err != nil {
			return nil, err
		}
		s := make([]string, n)
		for i := range v {
			s[i] = v[i].String
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (g *Gudang) baca(ctx context.Context, tx *db.Tx, objek, q string, n int, args ...any) ([][]string, error) {
	if err := siap(objek, q); err != nil {
		return nil, err
	}
	rows, err := g.dari(tx).QueryContext(ctx, q, args...)
	if err != nil {
		return nil, Bungkus(err, "membaca")
	}
	defer func() { _ = rows.Close() }()
	out, err := pindai(rows, n)
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

func kePlan(b [][]string) []models.Plan {
	out := make([]models.Plan, 0, len(b))
	for _, s := range b {
		out = append(out, models.Plan{ID: s[0], CoverName: s[1], Business: s[2], BusinessID: s[3], Benefit: s[4], BenefitID: s[5]})
	}
	return out
}

// Daftar - satu halaman dan total barisnya.
func (g *Gudang) Daftar(ctx context.Context, s models.Saringan) ([]models.Plan, int, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return nil, 0, err
	}
	total, err := g.satuNilai(ctx, nil, Tabel, SqlJumlah(t))
	if err != nil {
		return nil, 0, err
	}
	b, err := g.baca(ctx, nil, Tabel, SqlDaftar(t, s.Urut, s.Turun), 6, (s.Halaman-1)*models.UkuranHalaman, models.UkuranHalaman)
	return kePlan(b), angka(total), err
}

// Ambil - satu baris; ErrTidakAda bila ID tidak ada.
func (g *Gudang) Ambil(ctx context.Context, tx *db.Tx, id string) (models.Plan, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return models.Plan{}, err
	}
	b, err := g.baca(ctx, tx, Tabel, SqlAmbil(t), 6, id)
	if err != nil {
		return models.Plan{}, err
	}
	if len(b) == 0 {
		return models.Plan{}, ErrTidakAda
	}
	return kePlan(b)[0], nil
}

// PemakaiNama - plan bernama sama selain kecualiID ("" = semua).
func (g *Gudang) PemakaiNama(ctx context.Context, tx *db.Tx, nama, kecualiID string) ([]models.Plan, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return nil, err
	}
	b, err := g.baca(ctx, tx, Tabel, SqlPemakaiNama(t), 6, nama, db.KosongJadiNil(kecualiID))
	return kePlan(b), err
}

// NomorBaru - satu nomor `M_PRODUCT_TYPE_LIFE_SEQ` (teks).
func (g *Gudang) NomorBaru(ctx context.Context, tx *db.Tx) (string, error) {
	seq, err := g.nama(Seq)
	if err != nil {
		return "", err
	}
	return g.satuNilai(ctx, tx, "DUAL", SqlNomorBaru(seq))
}

// AdaID - ID sudah terpakai?
func (g *Gudang) AdaID(ctx context.Context, tx *db.Tx, id string) (bool, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return false, err
	}
	s, err := g.satuNilai(ctx, tx, Tabel, SqlAdaID(t), id)
	return angka(s) > 0, err
}

func argPlan(p models.Plan) []any {
	k := db.KosongJadiNil
	return []any{k(p.CoverName), k(p.Business), k(p.BusinessID), k(p.Benefit), k(p.BenefitID)}
}

// Sisip - baris baru; ErrKembar bila ID terpakai (ORA-00001).
func (g *Gudang) Sisip(ctx context.Context, tx *db.Tx, p models.Plan) error {
	t, err := g.nama(Tabel)
	if err != nil {
		return err
	}
	_, err = g.tulis(ctx, tx, SqlSisip(t), "menyimpan", append([]any{p.ID}, argPlan(p)...)...)
	return err
}

// Ubah - satu baris; ErrTidakAda bila ID tidak ada.
func (g *Gudang) Ubah(ctx context.Context, tx *db.Tx, p models.Plan) error {
	t, err := g.nama(Tabel)
	if err != nil {
		return err
	}
	n, err := g.tulis(ctx, tx, SqlUbah(t), "mengubah", append(argPlan(p), p.ID)...)
	if err == nil && n == 0 {
		return ErrTidakAda
	}
	return err
}

// PilihanBusiness - business grup `grup` (K4).
func (g *Gudang) PilihanBusiness(ctx context.Context, grup string) ([]models.PilihanBusiness, error) {
	t, err := g.nama(MasterBusiness)
	if err != nil {
		return nil, err
	}
	b, err := g.baca(ctx, nil, MasterBusiness, SqlPilihanBusiness(t), 3, grup)
	out := make([]models.PilihanBusiness, 0, len(b))
	for _, s := range b {
		out = append(out, models.PilihanBusiness{ID: s[0], OldID: s[1], Note: s[2]})
	}
	return out, err
}

// PilihanBenefit - seluruh BENEFIT_LIFE (K4).
func (g *Gudang) PilihanBenefit(ctx context.Context) ([]models.PilihanBenefit, error) {
	t, err := g.nama(MasterBenefit)
	if err != nil {
		return nil, err
	}
	b, err := g.baca(ctx, nil, MasterBenefit, SqlPilihanBenefit(t), 2)
	out := make([]models.PilihanBenefit, 0, len(b))
	for _, s := range b {
		out = append(out, models.PilihanBenefit{ID: s[0], Benefit: s[1]})
	}
	return out, err
}
