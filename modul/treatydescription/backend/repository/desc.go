// Package repository membaca dan menulis `POOLDATA.TREATYDESC` - satu-satunya tabel modul ini (perintah work owner
// 05-10-2026: "hanya baca dari tabel yang saya kasih"). Setiap nilai diikat (`:n`, go-ora mengikat menurut URUTAN kemunculan); nama objek lewat `Qualify`. Nol COMMIT -
// transaksi milik services.
//
// ⛔ Prosedur `PEGA_TREATYDESC` dan `PEGA_M_TREATYDESC` TIDAK dipanggil; `M_TREATYDESC` (JSON Pega) tidak disentuh.
// ID baru dibentuk persis seperti `PEGA_TREATYDESC`: `'1' || LPAD(TREATY_DESCRIPTION_SEQ.NEXTVAL, 4, '0')`.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatydescription/backend/models"
)

// Nama objek warisan (katalog DEV 05-10-2026: TREATYDESC 13 baris, tanpa PK, constraint, indeks, maupun trigger;
// TREATY_DESCRIPTION_SEQ LAST_NUMBER 19, cache 0).
const (
	Tabel    = "TREATYDESC"
	Sequence = "TREATY_DESCRIPTION_SEQ"
)

// AwalanID - awalan ID baru yang ditulis `PEGA_TREATYDESC` (literal '1', bukan situs `M_SITE_DATABASE`).
const AwalanID = "1"

// DaftarTabelDitulis dan DaftarTabelDibacaSaja - penjaga modul.
var (
	DaftarTabelDitulis    = []string{Tabel}
	DaftarTabelDibacaSaja = []string{}
)

var (
	// ErrTidakAda - ID tidak ada.
	ErrTidakAda = errors.New("repository: treaty description tidak ada")
	// ErrBelumAda - tabel atau sequence warisan tidak ada di skema ini.
	ErrBelumAda = errors.New("repository: tabel atau sequence Treaty Description tidak ada di skema ini")
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

// kolomDesc - keempat kolom TREATYDESC; ISXOL NULL dibaca '0'.
const kolomDesc = `ID, DESCNAME, NVL(ISXOL, '0'), STATUSAKTIF`

func pindaiDesc(p pemindai) (models.Desc, error) {
	v, err := teks(p, 4)
	if err != nil {
		return models.Desc{}, err
	}
	return models.Desc{ID: v[0], DescName: v[1], IsXOL: v[2], StatusAktif: v[3], Aktif: v[3] != models.Nonaktif}, nil
}

// SqlDaftar - daftar bersaring kata (ID atau nama), jenis (ISXOL; NULL = Non XOL), dan status (STATUSAKTIF; NULL =
// aktif); nil = tanpa saringan. Urut ID seperti `BrowseTreatyDesc_RD`.
func SqlDaftar(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s
	  WHERE (:1 IS NULL OR UPPER(ID) LIKE :2 ESCAPE '\' OR UPPER(DESCNAME) LIKE :3 ESCAPE '\')
	    AND (:4 IS NULL OR NVL(ISXOL, '0') = :5)
	    AND (:6 IS NULL OR DECODE(STATUSAKTIF, '0', '0', '1') = :7)
	  ORDER BY ID`, kolomDesc, t)
}

// Daftar - baris bersaring; xol dan status "" = semua.
func (g *Gudang) Daftar(ctx context.Context, kata, xol, status string) ([]models.Desc, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return nil, err
	}
	pola, x, s := PolaCari(kata), db.KosongJadiNil(xol), db.KosongJadiNil(status)
	return baca(ctx, g.db, SqlDaftar(t), pindaiDesc, pola, pola, pola, x, x, s, s)
}

// SqlAmbil - satu baris.
func SqlAmbil(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, kolomDesc, t)
}

// Ambil - satu baris.
func (g *Gudang) Ambil(ctx context.Context, tx *db.Tx, id string) (models.Desc, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return models.Desc{}, err
	}
	d, err := baca(ctx, g.dari(tx), SqlAmbil(t), pindaiDesc, id)
	if err != nil {
		return models.Desc{}, err
	}
	if len(d) == 0 {
		return models.Desc{}, ErrTidakAda
	}
	return d[0], nil
}

// SqlPemakaiNama - ID baris lain yang DESCNAME-nya sama (tanpa beda huruf dan spasi tepi), selain :2.
func SqlPemakaiNama(t string) string {
	return fmt.Sprintf(`SELECT ID FROM %s WHERE UPPER(TRIM(DESCNAME)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0)) ORDER BY ID`, t)
}

// PemakaiNama - ID baris lain yang memakai nama ini, selain kecualiID.
func (g *Gudang) PemakaiNama(ctx context.Context, tx *db.Tx, nama, kecualiID string) ([]string, error) {
	t, err := g.nama(Tabel)
	if err != nil {
		return nil, err
	}
	return baca(ctx, g.dari(tx), SqlPemakaiNama(t), func(p pemindai) (string, error) {
		v, err := teks(p, 1)
		if err != nil {
			return "", err
		}
		return v[0], nil
	}, nama, db.KosongJadiNil(kecualiID))
}

// SqlIDBaru - ID baru seperti `PEGA_TREATYDESC`: awalan '1' (diikat) + nomor TREATY_DESCRIPTION_SEQ 4 digit.
func SqlIDBaru(seq string) string {
	return fmt.Sprintf(`SELECT :1 || LPAD(TO_CHAR(%s.NEXTVAL), 4, '0') FROM DUAL`, seq)
}

// IDBaru - satu ID baru (sequence maju satu).
func (g *Gudang) IDBaru(ctx context.Context, tx *db.Tx) (string, error) {
	seq, err := g.nama(Sequence)
	if err != nil {
		return "", err
	}
	q := SqlIDBaru(seq)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var id sql.NullString
	if err := g.dari(tx).QueryRowContext(ctx, q, AwalanID).Scan(&id); err != nil {
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

// SqlSisip - keempat kolom yang ditulis `PEGA_TREATYDESC`.
func SqlSisip(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, DESCNAME, ISXOL, STATUSAKTIF) VALUES (:1, :2, :3, :4)`, t)
}

// Sisip - baris baru.
func (g *Gudang) Sisip(ctx context.Context, tx *db.Tx, r models.Desc) error {
	t, err := g.nama(Tabel)
	if err != nil {
		return err
	}
	_, err = g.tulis(ctx, tx, SqlSisip(t), "menyimpan", r.ID, r.DescName, r.IsXOL, r.StatusAktif)
	return err
}

// SqlUbah - kolom yang diubah `PEGA_TREATYDESC`; ID tidak pernah disentuh.
func SqlUbah(t string) string {
	return fmt.Sprintf(`UPDATE %s SET DESCNAME = :1, ISXOL = :2, STATUSAKTIF = :3 WHERE ID = :4`, t)
}

// Ubah - isian form; ID tetap. Nol baris = ErrTidakAda.
func (g *Gudang) Ubah(ctx context.Context, tx *db.Tx, r models.Desc) error {
	t, err := g.nama(Tabel)
	if err != nil {
		return err
	}
	jumlah, err := g.tulis(ctx, tx, SqlUbah(t), "mengubah", r.DescName, r.IsXOL, r.StatusAktif, r.ID)
	if err == nil && jumlah == 0 {
		return ErrTidakAda
	}
	return err
}
