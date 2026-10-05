// Package repository membaca dan menulis tabel Bordereaux di Oracle. SQL dirakit dari potongan tetap; setiap nilai
// diikat (`:n`, go-ora mengikat menurut URUTAN kemunculan). Nama tabel dan kolom berasal dari `models` (bukan dari
// masukan pengguna) dan tetap lewat `Qualify`.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/bordereaux/backend/models"
)

// Nama objek warisan dan baru.
const (
	TabelHeader  = "BORDEREAUX"
	TabelJSON    = "M_BORDEREAUX"
	TabelRiwayat = "BORDEREAUX_HISTORY"
	SeqRiwayat   = "SEQ_BORDEREAUX_HISTORY"
	TabelTreaty  = "TREATY_IN"
	TabelAgent   = "AGENT"
)

var (
	// ErrTidakAda - berkas tidak ada.
	ErrTidakAda = errors.New("repository: berkas bordereaux tidak ada")
	// ErrBelumDimigrasi - BORDEREAUX_HISTORY atau sequence-nya belum ada (migrasi 890 belum dijalankan).
	ErrBelumDimigrasi = errors.New("repository: migrasi bordereaux 890 belum dijalankan (-migrate)")
	// ErrIDTerpakai - ORA-00001 saat menyisipkan (ID bentrok).
	ErrIDTerpakai = errors.New("repository: ID sudah terpakai")
)

// Gudang - akses Oracle modul Bordereaux.
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

func (g *Gudang) nama(objek string) (string, error) { return g.db.Qualify(objek) }

func bungkus(err error, apa string) error {
	if err == nil {
		return nil
	}
	s := err.Error()
	switch {
	case strings.Contains(s, "ORA-00942"), strings.Contains(s, "ORA-02289"):
		return fmt.Errorf("%w: %v", ErrBelumDimigrasi, err)
	case strings.Contains(s, "ORA-00001"):
		return fmt.Errorf("%w: %v", ErrIDTerpakai, err)
	}
	return fmt.Errorf("repository: %s: %w", apa, err)
}

type pemindai interface{ Scan(...any) error }

func pindaiTeks(p pemindai, n int) ([]string, error) {
	v := make([]sql.NullString, n)
	tujuan := make([]any, n)
	for i := range v {
		tujuan[i] = &v[i]
	}
	if err := p.Scan(tujuan...); err != nil {
		return nil, err
	}
	out := make([]string, n)
	for i := range v {
		out[i] = v[i].String
	}
	return out, nil
}

func daftar[T any](ctx context.Context, j penjalan, q string, pindai func(pemindai) (T, error), args ...any) ([]T, error) {
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
	if err := rows.Err(); err != nil {
		return nil, bungkus(err, "membaca")
	}
	return out, nil
}

func jalankan(ctx context.Context, j penjalan, q, apa string, args ...any) (sql.Result, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	h, err := j.ExecContext(ctx, q, args...)
	return h, bungkus(err, apa)
}

// PolaCari - pola LIKE ber-ESCAPE '\' untuk "Contains, case-insensitive" Pega; kosong = nil.
func PolaCari(kata string) any {
	kata = strings.ToUpper(strings.TrimSpace(kata))
	if kata == "" {
		return nil
	}
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return "%" + r.Replace(kata) + "%"
}

// ---------------------------------------------------------------------------- header

const kolomHeader = `BDX_ID, TO_CHAR(TANGGAL, 'DD-MM-YYYY HH24:MI'), USER_INPUT, TYPE, TYPE_BUSINESS, MASTERID, CEDINGID,
  CEDINGNAME, SOBID, SOBNAME, TREATYNAME, TO_CHAR(BDXREPORT_START, 'DD-MM-YYYY'), TO_CHAR(BDXREPORT_END, 'DD-MM-YYYY'),
  REFFNO_OF_SOA, REFFNO_OF_BDX, POSITION, STATUSAKSEP`

func pindaiHeader(p pemindai) (models.Header, error) {
	v, err := pindaiTeks(p, 17)
	if err != nil {
		return models.Header{}, err
	}
	return models.Header{BdxID: v[0], Tanggal: v[1], UserInput: v[2], Type: v[3], TypeBusiness: v[4], MasterID: v[5],
		CedingID: v[6], CedingName: v[7], SobID: v[8], SobName: v[9], TreatyName: v[10], ReportStart: v[11],
		ReportEnd: v[12], ReffNoSOA: v[13], ReffNoBDX: v[14], Position: v[15], StatusAksep: v[16]}, nil
}

// saringan merakit WHERE daftar dari filter yang terisi saja; nomor penampung mulai `awal`.
func saringan(f models.Filter, awal int) (string, []any) {
	var (
		syarat []string
		args   []any
	)
	n := awal
	tambah := func(pola string, nilai any) {
		syarat = append(syarat, fmt.Sprintf(pola, n))
		args = append(args, nilai)
		n++
	}
	berisi := func(kolom, kata string) {
		if p := PolaCari(kata); p != nil {
			tambah("UPPER("+kolom+") LIKE :%d ESCAPE '\\'", p)
		}
	}
	sama := func(kolom, nilai string) {
		if nilai = strings.TrimSpace(nilai); nilai != "" {
			tambah("UPPER("+kolom+") = UPPER(:%d)", nilai)
		}
	}
	berisi("BDX_ID", f.BdxID)
	sama("TYPE", f.Type)
	sama("TYPE_BUSINESS", f.Business)
	berisi("REFFNO_OF_SOA", f.ReffNoSOA)
	berisi("REFFNO_OF_BDX", f.ReffNoBDX)
	berisi("CEDINGNAME", f.Ceding)
	berisi("TREATYNAME", f.Treaty)
	if s := strings.TrimSpace(f.ReportStart); s != "" {
		tambah("BDXREPORT_START >= TO_DATE(:%d, 'DD-MM-YYYY')", s)
	}
	if s := strings.TrimSpace(f.ReportEnd); s != "" {
		tambah("BDXREPORT_END <= TO_DATE(:%d, 'DD-MM-YYYY')", s)
	}
	berisi("POSITION", f.Position)
	sama("STATUSAKSEP", f.Status)
	if len(syarat) == 0 {
		return "1 = 1", nil
	}
	return strings.Join(syarat, " AND "), args
}

// Daftar - satu halaman daftar (`InboxBordereaux_RD`, urut waktu buat terbaru) dan jumlah seluruhnya.
func (g *Gudang) Daftar(ctx context.Context, f models.Filter, offset, ukuran int) ([]models.Header, int, error) {
	t, err := g.nama(TabelHeader)
	if err != nil {
		return nil, 0, err
	}
	where, args := saringan(f, 1)
	var total int
	qHitung := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, t, where)
	if err := db.PeriksaSQL(qHitung); err != nil {
		return nil, 0, err
	}
	if err := g.db.QueryRowContext(ctx, qHitung, args...).Scan(&total); err != nil {
		return nil, 0, bungkus(err, "menghitung daftar")
	}
	n := len(args) + 1
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE %s ORDER BY TANGGAL DESC NULLS LAST, BDX_ID DESC
	  OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY`, kolomHeader, t, where, n, n+1)
	hasil, err := daftar(ctx, g.db, q, pindaiHeader, append(args, offset, ukuran)...)
	return hasil, total, err
}

// AmbilHeader - satu berkas.
func (g *Gudang) AmbilHeader(ctx context.Context, tx *db.Tx, id string) (models.Header, error) {
	t, err := g.nama(TabelHeader)
	if err != nil {
		return models.Header{}, err
	}
	hs, err := daftar(ctx, g.dari(tx), fmt.Sprintf(`SELECT %s FROM %s WHERE BDX_ID = :1`, kolomHeader, t), pindaiHeader, id)
	if err != nil {
		return models.Header{}, err
	}
	if len(hs) == 0 {
		return models.Header{}, ErrTidakAda
	}
	return hs[0], nil
}

// AdaBerkas - BDX_ID sudah terpakai?
func (g *Gudang) AdaBerkas(ctx context.Context, tx *db.Tx, id string) (bool, error) {
	t, err := g.nama(TabelHeader)
	if err != nil {
		return false, err
	}
	q := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE BDX_ID = :1`, t)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var n int
	if err := g.dari(tx).QueryRowContext(ctx, q, id).Scan(&n); err != nil {
		return false, bungkus(err, "memeriksa BDX_ID")
	}
	return n > 0, nil
}

func kosongNil(s string) any { return db.KosongJadiNil(strings.TrimSpace(s)) }

// SisipHeader - berkas baru; TANGGAL = SYSDATE, kecuali h.Tanggal terisi `DD-MM-YYYY HH24:MI:SS` (Copy Old Data:
// waktu buat kasus Pega).
func (g *Gudang) SisipHeader(ctx context.Context, tx *db.Tx, h models.Header) error {
	t, err := g.nama(TabelHeader)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`INSERT INTO %s (BDX_ID, TANGGAL, USER_INPUT, TYPE, TYPE_BUSINESS, MASTERID, CEDINGID, CEDINGNAME,
	  SOBID, SOBNAME, TREATYNAME, BDXREPORT_START, BDXREPORT_END, REFFNO_OF_SOA, REFFNO_OF_BDX, POSITION, STATUSAKSEP)
	  VALUES (:1, NVL(TO_DATE(:2, 'DD-MM-YYYY HH24:MI:SS'), SYSDATE), :3, :4, :5, :6, :7, :8, :9, :10, :11,
	  TO_DATE(:12, 'DD-MM-YYYY'), TO_DATE(:13, 'DD-MM-YYYY'), :14, :15, :16, :17)`, t)
	_, err = jalankan(ctx, g.dari(tx), q, "menyisipkan berkas", h.BdxID, kosongNil(h.Tanggal), kosongNil(h.UserInput), kosongNil(h.Type),
		kosongNil(h.TypeBusiness), kosongNil(h.MasterID), kosongNil(h.CedingID), kosongNil(h.CedingName), kosongNil(h.SobID),
		kosongNil(h.SobName), kosongNil(h.TreatyName), kosongNil(h.ReportStart), kosongNil(h.ReportEnd),
		kosongNil(h.ReffNoSOA), kosongNil(h.ReffNoBDX), kosongNil(h.Position), kosongNil(h.StatusAksep))
	return err
}

// UbahHeader - isian form; TANGGAL, USER_INPUT, POSITION, dan STATUSAKSEP tidak disentuh.
func (g *Gudang) UbahHeader(ctx context.Context, tx *db.Tx, h models.Header) error {
	t, err := g.nama(TabelHeader)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`UPDATE %s SET TYPE = :1, TYPE_BUSINESS = :2, MASTERID = :3, CEDINGID = :4, CEDINGNAME = :5,
	  SOBID = :6, SOBNAME = :7, TREATYNAME = :8, BDXREPORT_START = TO_DATE(:9, 'DD-MM-YYYY'),
	  BDXREPORT_END = TO_DATE(:10, 'DD-MM-YYYY'), REFFNO_OF_SOA = :11, REFFNO_OF_BDX = :12 WHERE BDX_ID = :13`, t)
	hasil, err := jalankan(ctx, g.dari(tx), q, "mengubah berkas", kosongNil(h.Type), kosongNil(h.TypeBusiness),
		kosongNil(h.MasterID), kosongNil(h.CedingID), kosongNil(h.CedingName), kosongNil(h.SobID), kosongNil(h.SobName),
		kosongNil(h.TreatyName), kosongNil(h.ReportStart), kosongNil(h.ReportEnd), kosongNil(h.ReffNoSOA),
		kosongNil(h.ReffNoBDX), h.BdxID)
	if err != nil {
		return err
	}
	if n, _ := hasil.RowsAffected(); n == 0 {
		return ErrTidakAda
	}
	return nil
}

// UbahStatus memindah berkas HANYA bila posisi dan statusnya masih seperti yang dibaca (penjaga dua pemroses
// serentak); jawab false = sudah berubah.
func (g *Gudang) UbahStatus(ctx context.Context, tx *db.Tx, id, posisiLama, statusLama, posisiBaru, statusBaru string) (bool, error) {
	t, err := g.nama(TabelHeader)
	if err != nil {
		return false, err
	}
	q := fmt.Sprintf(`UPDATE %s SET POSITION = :1, STATUSAKSEP = :2 WHERE BDX_ID = :3
	  AND NVL(POSITION, ' ') = NVL(:4, ' ') AND NVL(STATUSAKSEP, ' ') = NVL(:5, ' ')`, t)
	hasil, err := jalankan(ctx, g.dari(tx), q, "mengubah status", kosongNil(posisiBaru), kosongNil(statusBaru), id,
		kosongNil(posisiLama), kosongNil(statusLama))
	if err != nil {
		return false, err
	}
	n, _ := hasil.RowsAffected()
	return n == 1, nil
}

// Chart - jumlah berkas per Business, Type, dan Ceding (chart daftar, keputusan work owner 04-10-2026).
func (g *Gudang) Chart(ctx context.Context) ([]models.IrisanChart, error) {
	t, err := g.nama(TabelHeader)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT TYPE_BUSINESS, TYPE, CEDINGID, CEDINGNAME, COUNT(*) FROM %s
	  GROUP BY TYPE_BUSINESS, TYPE, CEDINGID, CEDINGNAME ORDER BY 1, 2, 4`, t)
	return daftar(ctx, g.db, q, func(p pemindai) (models.IrisanChart, error) {
		var (
			v [4]sql.NullString
			n int
		)
		if err := p.Scan(&v[0], &v[1], &v[2], &v[3], &n); err != nil {
			return models.IrisanChart{}, err
		}
		return models.IrisanChart{Business: v[0].String, Type: v[1].String, CedingID: v[2].String, CedingName: v[3].String, Jumlah: n}, nil
	})
}

// ---------------------------------------------------------------------------- detail

func ekspresiBaca(k models.KolomCSV) string {
	switch k.Jenis {
	case models.Angka:
		return fmt.Sprintf(db.FmtDesimal, k.Kolom)
	case models.Tanggal:
		return fmt.Sprintf("TO_CHAR(%s, 'DD-MM-YYYY')", k.Kolom)
	}
	return k.Kolom
}

// Detail - baris detail satu berkas di tabel kombinasinya, urut ID.
func (g *Gudang) Detail(ctx context.Context, tx *db.Tx, k models.KombinasiBdx, id string) ([]models.Baris, error) {
	t, err := g.nama(k.Tabel)
	if err != nil {
		return nil, err
	}
	kolom := []string{"ID"}
	for _, c := range k.Kolom {
		kolom = append(kolom, ekspresiBaca(c))
	}
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = :1 ORDER BY ID`, strings.Join(kolom, ", "), t, k.KolomInduk)
	return daftar(ctx, g.dari(tx), q, func(p pemindai) (models.Baris, error) {
		v, err := pindaiTeks(p, len(kolom))
		if err != nil {
			return nil, err
		}
		b := models.Baris{models.KolomID: v[0]}
		for i, c := range k.Kolom {
			b[c.Kolom] = v[i+1]
		}
		return b, nil
	}, id)
}

// HapusDetail - seluruh baris detail satu berkas di tabel satu kombinasi.
func (g *Gudang) HapusDetail(ctx context.Context, tx *db.Tx, k models.KombinasiBdx, id string) (int64, error) {
	t, err := g.nama(k.Tabel)
	if err != nil {
		return 0, err
	}
	hasil, err := jalankan(ctx, g.dari(tx), fmt.Sprintf(`DELETE FROM %s WHERE %s = :1`, t, k.KolomInduk), "menghapus detail", id)
	if err != nil {
		return 0, err
	}
	n, _ := hasil.RowsAffected()
	return n, nil
}

// PecahAngka - teks desimal menjadi koefisien bulat dan skala, supaya angka masuk Oracle sebagai
// `TO_NUMBER(koef) / POWER(10, skala)` tanpa float dan tanpa bergantung NLS.
func PecahAngka(teks string) (any, any, error) {
	if strings.TrimSpace(teks) == "" {
		return nil, 0, nil
	}
	d, _, err := apd.NewFromString(teks)
	if err != nil {
		return nil, 0, fmt.Errorf("angka %q tidak sah", teks)
	}
	s := d.Text('f')
	tanda := ""
	if strings.HasPrefix(s, "-") {
		tanda, s = "-", s[1:]
	}
	bulat, pecahan, _ := strings.Cut(s, ".")
	koef := strings.TrimLeft(bulat+pecahan, "0")
	if koef == "" {
		return "0", 0, nil
	}
	return tanda + koef, len(pecahan), nil
}

// SisipDetail - satu baris detail; idDetail dari layanan (`BDX_DTL-...`).
func (g *Gudang) SisipDetail(ctx context.Context, tx *db.Tx, k models.KombinasiBdx, id, idDetail string, b models.Baris) error {
	t, err := g.nama(k.Tabel)
	if err != nil {
		return err
	}
	kolom := []string{"ID", k.KolomInduk}
	nilai := []string{":1", ":2"}
	args := []any{idDetail, id}
	n := 3
	for _, c := range k.Kolom {
		kolom = append(kolom, c.Kolom)
		v := b[c.Kolom]
		switch c.Jenis {
		case models.Angka:
			koef, skala, err := PecahAngka(v)
			if err != nil {
				return fmt.Errorf("kolom %s: %w", c.Kolom, err)
			}
			nilai = append(nilai, fmt.Sprintf("(TO_NUMBER(:%d) / POWER(10, :%d))", n, n+1))
			args = append(args, koef, skala)
			n += 2
		case models.Tanggal:
			nilai = append(nilai, fmt.Sprintf("TO_DATE(:%d, 'DD-MM-YYYY')", n))
			args = append(args, kosongNil(v))
			n++
		default:
			nilai = append(nilai, fmt.Sprintf(":%d", n))
			args = append(args, kosongNil(v))
			n++
		}
	}
	q := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, t, strings.Join(kolom, ", "), strings.Join(nilai, ", "))
	_, err = jalankan(ctx, g.dari(tx), q, "menyisipkan detail", args...)
	return err
}

// ---------------------------------------------------------------------------- hapus berkas

// HapusBerkas - detail di SEMUA tabel kombinasi (Type/Business pernah berubah), riwayat, baris JSON lama, lampiran,
// lalu header. Prosedur Pega `PEGA_DELETE_BORDEREAUX` hanya menghapus header dan JSON - detailnya tertinggal
// (diperbaiki).
func (g *Gudang) HapusBerkas(ctx context.Context, tx *db.Tx, id string) error {
	for _, k := range models.Kombinasi {
		if _, err := g.HapusDetail(ctx, tx, k, id); err != nil {
			return err
		}
	}
	for _, x := range []struct{ tabel, kolom string }{
		{TabelRiwayat, "BDX_ID"}, {TabelJSON, "ID"}, {"M_ATTACHMENTBORDEREAUX", "BDX_ID"},
	} {
		t, err := g.nama(x.tabel)
		if err != nil {
			return err
		}
		if _, err := jalankan(ctx, g.dari(tx), fmt.Sprintf(`DELETE FROM %s WHERE %s = :1`, t, x.kolom), "menghapus "+x.tabel, id); err != nil {
			return err
		}
	}
	t, err := g.nama(TabelHeader)
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, g.dari(tx), fmt.Sprintf(`DELETE FROM %s WHERE BDX_ID = :1`, t), "menghapus berkas", id)
	if err != nil {
		return err
	}
	if n, _ := hasil.RowsAffected(); n == 0 {
		return ErrTidakAda
	}
	return nil
}

// ---------------------------------------------------------------------------- riwayat

// Riwayat - riwayat persetujuan satu berkas, urut waktu.
func (g *Gudang) Riwayat(ctx context.Context, id string) ([]models.Riwayat, error) {
	t, err := g.nama(TabelRiwayat)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT TO_CHAR(TANGGAL, 'DD-MM-YYYY HH24:MI'), PIC, IS_APPROVED, KOMENTAR FROM %s
	  WHERE BDX_ID = :1 ORDER BY TANGGAL, ID`, t)
	return daftar(ctx, g.db, q, func(p pemindai) (models.Riwayat, error) {
		v, err := pindaiTeks(p, 4)
		if err != nil {
			return models.Riwayat{}, err
		}
		return models.Riwayat{Tanggal: v[0], PIC: v[1], IsApproved: v[2] == "1", Komentar: v[3]}, nil
	}, id)
}

// SisipRiwayat - satu baris riwayat; TANGGAL = SYSDATE.
func (g *Gudang) SisipRiwayat(ctx context.Context, tx *db.Tx, id, pic string, setuju bool, komentar string) error {
	return g.sisipRiwayat(ctx, tx, id, nil, pic, setuju, komentar)
}

// SisipRiwayatLama - satu baris riwayat ber-TANGGAL `DD-MM-YYYY HH24:MI:SS` (Copy Old Data: CommentList Pega).
func (g *Gudang) SisipRiwayatLama(ctx context.Context, tx *db.Tx, id, tanggal, pic string, setuju bool, komentar string) error {
	return g.sisipRiwayat(ctx, tx, id, kosongNil(tanggal), pic, setuju, komentar)
}

func (g *Gudang) sisipRiwayat(ctx context.Context, tx *db.Tx, id string, tanggal any, pic string, setuju bool, komentar string) error {
	t, err := g.nama(TabelRiwayat)
	if err != nil {
		return err
	}
	nomor, err := g.db.NomorBerikut(ctx, tx, SeqRiwayat)
	if err != nil {
		return bungkus(err, "mengambil ID riwayat")
	}
	bendera := "0"
	if setuju {
		bendera = "1"
	}
	_, err = jalankan(ctx, g.dari(tx), fmt.Sprintf(`INSERT INTO %s (ID, BDX_ID, TANGGAL, PIC, IS_APPROVED, KOMENTAR)
	  VALUES (:1, :2, NVL(TO_DATE(:3, 'DD-MM-YYYY HH24:MI:SS'), SYSDATE), :4, :5, :6)`, t), "menyimpan riwayat", nomor, id,
		tanggal, pic, bendera, kosongNil(komentar))
	return err
}

// ---------------------------------------------------------------------------- copy old data

// JSONLama - baris M_BORDEREAUX (JSON Pega, hanya DIBACA) beserta Type/Business header BORDEREAUX-nya bila ada;
// id kosong = semua, urut ID turun.
func (g *Gudang) JSONLama(ctx context.Context, tx *db.Tx, id string) ([]models.JSONLama, error) {
	tj, err := g.nama(TabelJSON)
	if err != nil {
		return nil, err
	}
	th, err := g.nama(TabelHeader)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT m.ID, m.DATA_JSON, CASE WHEN b.BDX_ID IS NULL THEN 0 ELSE 1 END, b.TYPE, b.TYPE_BUSINESS
	  FROM %s m LEFT JOIN %s b ON b.BDX_ID = m.ID WHERE (:1 IS NULL OR m.ID = :2) ORDER BY m.ID DESC`, tj, th)
	return daftar(ctx, g.dari(tx), q, func(p pemindai) (models.JSONLama, error) {
		var (
			v   [3]sql.NullString
			isi db.PindaiTeksPanjang
			ada int
		)
		if err := p.Scan(&v[0], &isi, &ada, &v[1], &v[2]); err != nil {
			return models.JSONLama{}, err
		}
		return models.JSONLama{ID: v[0].String, Isi: isi.Teks(), AdaHeader: ada == 1, Type: v[1].String, Business: v[2].String}, nil
	}, kosongNil(id), kosongNil(id))
}

// CacahDetail - jumlah baris detail per BDX_ID di tabel satu kombinasi.
func (g *Gudang) CacahDetail(ctx context.Context, k models.KombinasiBdx) (map[string]int, error) {
	t, err := g.nama(k.Tabel)
	if err != nil {
		return nil, err
	}
	return cacah(ctx, g.db, fmt.Sprintf(`SELECT %s, COUNT(*) FROM %s GROUP BY %s`, k.KolomInduk, t, k.KolomInduk))
}

// CacahRiwayat - jumlah baris BORDEREAUX_HISTORY per BDX_ID.
func (g *Gudang) CacahRiwayat(ctx context.Context) (map[string]int, error) {
	t, err := g.nama(TabelRiwayat)
	if err != nil {
		return nil, err
	}
	return cacah(ctx, g.db, fmt.Sprintf(`SELECT BDX_ID, COUNT(*) FROM %s GROUP BY BDX_ID`, t))
}

func cacah(ctx context.Context, j penjalan, q string) (map[string]int, error) {
	type pasangan struct {
		id string
		n  int
	}
	d, err := daftar(ctx, j, q, func(p pemindai) (pasangan, error) {
		var (
			id sql.NullString
			n  int
		)
		err := p.Scan(&id, &n)
		return pasangan{id.String, n}, err
	})
	if err != nil {
		return nil, err
	}
	hasil := make(map[string]int, len(d))
	for _, x := range d {
		hasil[x.id] = x.n
	}
	return hasil, nil
}

// ---------------------------------------------------------------------------- master treaty

// agentAktif - nilai AGENT.STATUSACTIVE (teks) untuk cedant aktif; diikat, bukan literal SQL.
const agentAktif = "1"

// CariCedant - AGENT aktif ber-CLIENTID yang namanya memuat kata (`BrowseAgentNonLife_RD`; syarat
// `AgentType2 != "LIFE INSURANCE"` tidak dapat ditiru - kolomnya tidak ada di AGENT DEV).
func (g *Gudang) CariCedant(ctx context.Context, kata string, batas int) ([]models.Cedant, error) {
	t, err := g.nama(TabelAgent)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT ID, CLIENTNAME FROM %s WHERE STATUSACTIVE = :1 AND CLIENTID IS NOT NULL
	  AND UPPER(CLIENTNAME) LIKE :2 ESCAPE '\' ORDER BY UPPER(CLIENTNAME), ID FETCH FIRST :3 ROWS ONLY`, t)
	pola := PolaCari(kata)
	if pola == nil {
		pola = "%"
	}
	return daftar(ctx, g.db, q, func(p pemindai) (models.Cedant, error) {
		v, err := pindaiTeks(p, 2)
		if err != nil {
			return models.Cedant{}, err
		}
		return models.Cedant{ID: v[0], Nama: v[1]}, nil
	}, agentAktif, pola, batas)
}

const kolomTreaty = `ID, TREATYCONTRACTNAME, PROPORTIONTYPE, LEADINGREINSSOURCEID, LEADINGREINSSOURCE, CEDINGID, CEDING`

func pindaiTreaty(p pemindai) (models.MasterTreaty, error) {
	v, err := pindaiTeks(p, 7)
	if err != nil {
		return models.MasterTreaty{}, err
	}
	return models.MasterTreaty{ID: v[0], ContractName: v[1], ReinsType: v[2], SobID: v[3], SobName: v[4], CedingID: v[5], CedingName: v[6]}, nil
}

// CariTreaty - TREATY_IN satu cedant (`BrowseTREATY_IN`: CedingID = cedant terpilih, urut ID turun, maks. 500).
func (g *Gudang) CariTreaty(ctx context.Context, cedingID string, batas int) ([]models.MasterTreaty, error) {
	t, err := g.nama(TabelTreaty)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE CEDINGID = :1 ORDER BY ID DESC FETCH FIRST :2 ROWS ONLY`, kolomTreaty, t)
	return daftar(ctx, g.db, q, pindaiTreaty, cedingID, batas)
}

// AmbilTreaty - satu TREATY_IN menurut ID (Save menyalin ceding, SOB, dan nama kontrak dari sini).
func (g *Gudang) AmbilTreaty(ctx context.Context, id string) (models.MasterTreaty, error) {
	t, err := g.nama(TabelTreaty)
	if err != nil {
		return models.MasterTreaty{}, err
	}
	ts, err := daftar(ctx, g.db, fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, kolomTreaty, t), pindaiTreaty, id)
	if err != nil {
		return models.MasterTreaty{}, err
	}
	if len(ts) == 0 {
		return models.MasterTreaty{}, ErrTidakAda
	}
	return ts[0], nil
}
