// Package repository membaca dan menulis `POOLDATA.TREATYGROUP`; `TREATYGROUPOJK`, `BUSINESSGROUP`, dan
// `M_SITE_DATABASE` dibaca saja. Setiap nilai diikat (`:n`, go-ora mengikat menurut URUTAN kemunculan); nama objek lewat
// `Qualify`. Nol COMMIT - transaksi milik services.
//
// ⛔ Prosedur `PEGA_TREATYGROUP` dan `PEGA_M_TREATYGROUP` TIDAK dipanggil; `M_TREATYGROUP` (JSON Pega) tidak disentuh.
// ID baru dibentuk seperti keduanya: situs aktif + `TREATYGROUP_SEQ` 4 digit (DEV: `1` + 4 digit).
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatygroup/backend/models"
)

// Nama objek warisan.
const (
	Tabel           = "TREATYGROUP"
	TabelOjk        = "TREATYGROUPOJK"
	TabelBisnisGrup = "BUSINESSGROUP"
	TabelSitus      = "M_SITE_DATABASE"
	Sequence        = "TREATYGROUP_SEQ"
)

// DaftarTabelDitulis dan DaftarTabelDibacaSaja - penjaga modul.
var (
	DaftarTabelDitulis    = []string{Tabel}
	DaftarTabelDibacaSaja = []string{TabelOjk, TabelBisnisGrup, TabelSitus}
)

// situsAktif - nilai `M_SITE_DATABASE.CURRENT_SITE` situs yang sedang dipakai; diikat, bukan literal SQL.
const situsAktif = "1"

// bukanSyariah - grup bisnis berakhiran SYARIAH tidak ditampilkan (saringan Pega `BrowseBusinessGroup_RD`
// `.Note NotEndsWith "SYARIAH"`; keputusan work owner 05-10-2026: "tidak ada syariah").
const bukanSyariah = `NVL(UPPER(TRIM(NOTE)), ' ') NOT LIKE '%SYARIAH'`

var (
	// ErrTidakAda - ID tidak ada.
	ErrTidakAda = errors.New("repository: treaty group tidak ada")
	// ErrBelumAda - tabel atau sequence warisan tidak ada di skema ini.
	ErrBelumAda = errors.New("repository: tabel atau sequence Treaty Group tidak ada di skema ini")
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

// baca menjalankan SELECT dan memindai setiap baris dengan pindai.
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

// kolomGrup - baris TREATYGROUP `t` dan nama grup bisnis COA-nya `b`.
const kolomGrup = `t.ID, t.OLDID, t.OJKBUSINESSID, t.OJKBUSINESSNAME, t.OJKBUSINESSNAMEIDN, t.ORDERNO, t.TREATYGROUPNAME,
  t.TREATYGROUPSOANAME, t.TGLUPDATE, t.USERID, t.COAID, b.NOTE`

func pindaiGrup(p pemindai) (models.Grup, error) {
	v, err := teks(p, 12)
	if err != nil {
		return models.Grup{}, err
	}
	return models.Grup{ID: v[0], OldID: v[1], OjkID: v[2], OjkName: v[3], OjkNameIDN: v[4], OrderNo: v[5], Name: v[6],
		SoaName: v[7], TglUpdate: v[8], UserID: v[9], CoaID: v[10], CoaName: v[11]}, nil
}

// SqlDaftar - daftar bersaring kata (ID, nama, nama SOA, nama OJK) dan OJK (nil = semua), urut Order No lalu nama.
func SqlDaftar(t, b string) string {
	return fmt.Sprintf(`SELECT %s FROM %s t LEFT JOIN %s b ON b.ID = t.COAID
	  WHERE (:1 IS NULL OR UPPER(t.ID) LIKE :2 ESCAPE '\' OR UPPER(t.TREATYGROUPNAME) LIKE :3 ESCAPE '\'
	    OR UPPER(t.TREATYGROUPSOANAME) LIKE :4 ESCAPE '\' OR UPPER(t.OJKBUSINESSNAME) LIKE :5 ESCAPE '\')
	    AND (:6 IS NULL OR t.OJKBUSINESSID = :7)
	  ORDER BY LPAD(t.ORDERNO, 10, '0'), UPPER(t.TREATYGROUPNAME), t.ID`, kolomGrup, t, b)
}

// Daftar - baris bersaring; ojk "" = semua OJK.
func (g *Gudang) Daftar(ctx context.Context, kata, ojk string) ([]models.Grup, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return nil, err
	}
	b, err := g.nama(TabelBisnisGrup)
	if err != nil {
		return nil, err
	}
	pola, o := PolaCari(kata), db.KosongJadiNil(strings.TrimSpace(ojk))
	return baca(ctx, g.db, SqlDaftar(t, b), pindaiGrup, pola, pola, pola, pola, pola, o, o)
}

// Ambil - satu baris.
func (g *Gudang) Ambil(ctx context.Context, tx *db.Tx, id string) (models.Grup, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return models.Grup{}, err
	}
	b, err := g.nama(TabelBisnisGrup)
	if err != nil {
		return models.Grup{}, err
	}
	d, err := baca(ctx, g.dari(tx), fmt.Sprintf(`SELECT %s FROM %s t LEFT JOIN %s b ON b.ID = t.COAID WHERE t.ID = :1`,
		kolomGrup, t, b), pindaiGrup, id)
	if err != nil {
		return models.Grup{}, err
	}
	if len(d) == 0 {
		return models.Grup{}, ErrTidakAda
	}
	return d[0], nil
}

// SqlAnak - grup bisnis ber-TOPID grup ini, tanpa yang berakhiran SYARIAH.
func SqlAnak(b string) string {
	return fmt.Sprintf(`SELECT ID, NOTE, ALIASNAME FROM %s WHERE TOPID = :1 AND %s ORDER BY UPPER(NOTE), ID`, b, bukanSyariah)
}

// Anak - grup bisnis anak (View).
func (g *Gudang) Anak(ctx context.Context, id string) ([]models.BisnisGrup, error) {
	b, err := g.nama(TabelBisnisGrup)
	if err != nil {
		return nil, err
	}
	return baca(ctx, g.db, SqlAnak(b), func(p pemindai) (models.BisnisGrup, error) {
		v, err := teks(p, 3)
		if err != nil {
			return models.BisnisGrup{}, err
		}
		return models.BisnisGrup{ID: v[0], Name: v[1], Alias: v[2]}, nil
	}, id)
}

func pindaiOjk(p pemindai) (models.Ojk, error) {
	v, err := teks(p, 4)
	if err != nil {
		return models.Ojk{}, err
	}
	return models.Ojk{ID: v[0], Name: v[1], NameIDN: v[2], OrderNo: v[3]}, nil
}

// SqlDaftarOjk - pilihan OJK Business, urut Order No.
func SqlDaftarOjk(o string) string {
	return fmt.Sprintf(`SELECT ID, NAME, NAMEIDN, ORDERNO FROM %s ORDER BY LPAD(ORDERNO, 10, '0'), ID`, o)
}

// SqlCoaPerOjk - COA setiap OJK: COAID yang paling sering di grup se-OJK, lalu terkecil (sama dengan SqlCoaSeOjk),
// beserta namanya.
func SqlCoaPerOjk(t, b string) string {
	return fmt.Sprintf(`SELECT OJKBUSINESSID, COAID, NOTE FROM (
	  SELECT t.OJKBUSINESSID, t.COAID, b.NOTE,
	    ROW_NUMBER() OVER (PARTITION BY t.OJKBUSINESSID ORDER BY COUNT(*) DESC, t.COAID) rn
	  FROM %s t LEFT JOIN %s b ON b.ID = t.COAID
	  WHERE t.COAID IS NOT NULL
	  GROUP BY t.OJKBUSINESSID, t.COAID, b.NOTE)
	WHERE rn = 1`, t, b)
}

// DaftarOjk - seluruh baris TREATYGROUPOJK beserta COA masing-masing.
func (g *Gudang) DaftarOjk(ctx context.Context) ([]models.Ojk, error) {
	o, err := g.nama(TabelOjk)
	if err != nil {
		return nil, err
	}
	t, err := g.nama(Tabel)
	if err != nil {
		return nil, err
	}
	b, err := g.nama(TabelBisnisGrup)
	if err != nil {
		return nil, err
	}
	daftar, err := baca(ctx, g.db, SqlDaftarOjk(o), pindaiOjk)
	if err != nil {
		return nil, err
	}
	coa, err := baca(ctx, g.db, SqlCoaPerOjk(t, b), func(p pemindai) ([]string, error) { return teks(p, 3) })
	if err != nil {
		return nil, err
	}
	per := map[string][]string{}
	for _, c := range coa {
		per[c[0]] = c
	}
	for i := range daftar {
		if c, ada := per[daftar[i].ID]; ada {
			daftar[i].CoaID, daftar[i].CoaName = c[1], c[2]
		}
	}
	return daftar, nil
}

// AmbilOjk - satu baris TREATYGROUPOJK (disalin ke grup saat disimpan).
func (g *Gudang) AmbilOjk(ctx context.Context, tx *db.Tx, id string) (models.Ojk, error) {
	o, err := g.nama(TabelOjk)
	if err != nil {
		return models.Ojk{}, err
	}
	d, err := baca(ctx, g.dari(tx), fmt.Sprintf(`SELECT ID, NAME, NAMEIDN, ORDERNO FROM %s WHERE ID = :1`, o), pindaiOjk, id)
	if err != nil {
		return models.Ojk{}, err
	}
	if len(d) == 0 {
		return models.Ojk{}, ErrTidakAda
	}
	return d[0], nil
}

// PemakaiNama - ID baris lain yang TREATYGROUPNAME-nya sama (tanpa beda huruf dan spasi tepi), selain kecualiID.
func (g *Gudang) PemakaiNama(ctx context.Context, tx *db.Tx, nama, kecualiID string) ([]string, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT ID FROM %s WHERE UPPER(TRIM(TREATYGROUPNAME)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0)) ORDER BY ID`, t)
	return baca(ctx, g.dari(tx), q, func(p pemindai) (string, error) {
		v, err := teks(p, 1)
		if err != nil {
			return "", err
		}
		return v[0], nil
	}, nama, db.KosongJadiNil(kecualiID))
}

// SqlCoaSeOjk - COAID grup lain se-OJK (keputusan work owner 05-10-2026, opsi A): yang paling sering, lalu terkecil.
func SqlCoaSeOjk(t string) string {
	return fmt.Sprintf(`SELECT COAID FROM %s WHERE OJKBUSINESSID = :1 AND COAID IS NOT NULL
	  GROUP BY COAID ORDER BY COUNT(*) DESC, COAID FETCH FIRST 1 ROWS ONLY`, t)
}

// CoaSeOjk - COAID untuk grup baru OJK ini; "" = belum ada grup se-OJK yang ber-COAID.
func (g *Gudang) CoaSeOjk(ctx context.Context, tx *db.Tx, ojk string) (string, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return "", err
	}
	d, err := baca(ctx, g.dari(tx), SqlCoaSeOjk(t), func(p pemindai) (string, error) {
		v, err := teks(p, 1)
		if err != nil {
			return "", err
		}
		return v[0], nil
	}, ojk)
	if err != nil || len(d) == 0 {
		return "", err
	}
	return d[0], nil
}

// SqlIDBaru - ID baru seperti `PEGA_M_TREATYGROUP`: ID situs aktif `M_SITE_DATABASE` + nomor sequence 4 digit.
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

// SqlSisip - kolom yang ditulis `PEGA_TREATYGROUP` + COAID (opsi A); OLDID tidak ditulis.
func SqlSisip(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, OJKBUSINESSID, OJKBUSINESSNAME, OJKBUSINESSNAMEIDN, TREATYGROUPNAME,
	  TREATYGROUPSOANAME, TGLUPDATE, USERID, ORDERNO, COAID) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10)`, t)
}

// Sisip - baris baru.
func (g *Gudang) Sisip(ctx context.Context, tx *db.Tx, r models.Grup) error {
	t, err := g.nama(Tabel)
	if err != nil {
		return err
	}
	n := db.KosongJadiNil
	_, err = g.tulis(ctx, tx, SqlSisip(t), "menyimpan", r.ID, n(r.OjkID), n(r.OjkName), n(r.OjkNameIDN), n(r.Name),
		n(r.SoaName), n(r.TglUpdate), n(r.UserID), n(r.OrderNo), n(r.CoaID))
	return err
}

// SqlUbah - kolom yang diubah `PEGA_TREATYGROUP` + COAID (ikut OJK; tetap bila OJK tidak diganti); ID dan OLDID tidak
// disentuh.
func SqlUbah(t string) string {
	return fmt.Sprintf(`UPDATE %s SET OJKBUSINESSID = :1, OJKBUSINESSNAME = :2, OJKBUSINESSNAMEIDN = :3,
	  TREATYGROUPNAME = :4, TREATYGROUPSOANAME = :5, TGLUPDATE = :6, USERID = :7, ORDERNO = :8, COAID = :9
	  WHERE ID = :10`, t)
}

// Ubah - isian form; ID dan OLDID tetap.
func (g *Gudang) Ubah(ctx context.Context, tx *db.Tx, r models.Grup) error {
	t, err := g.nama(Tabel)
	if err != nil {
		return err
	}
	n := db.KosongJadiNil
	jumlah, err := g.tulis(ctx, tx, SqlUbah(t), "mengubah", n(r.OjkID), n(r.OjkName), n(r.OjkNameIDN), n(r.Name),
		n(r.SoaName), n(r.TglUpdate), n(r.UserID), n(r.OrderNo), n(r.CoaID), r.ID)
	if err == nil && jumlah == 0 {
		return ErrTidakAda
	}
	return err
}
