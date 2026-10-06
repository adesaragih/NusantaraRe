package repository

// SQL modul R/I Rate Life. Baca dari view (`RATE_LIFE_SUMMARY`, `RATE_LIFE`); tulis ke tabel fisik
// (`M_RATE_LIFE_SUMMARY`, `M_RATE_LIFE`): sisip = `JSON_OBJECT`, ubah = `JSON_MERGEPATCH` (kunci JSON lain milik Pega
// tetap). Baris `M_RATE_LIFE` milik satu ringkasan dipilih lewat view (`ID IN (SELECT ID FROM RATE_LIFE WHERE
// IDUSEDBY = :n)`) supaya maknanya SAMA dengan pembaca lain (`GetRateRetro`, `BrowseLifeRate_SQL`).
//
// ⚠️ `RATE_LIFE` view atas CLOB tanpa indeks: setiap kueri ber-IDUSEDBY mengurai seluruh JSON (puluhan detik di DEV
// untuk agregat). Karena itu kunci kembar upload dibaca SEKALI untuk seluruh ringkasan berkas (daftar IN).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/riratelife/backend/models"
)

var (
	// ErrTidakAda - ID tidak ada.
	ErrTidakAda = errors.New("repository: ringkasan R/I rate tidak ada")
	// ErrBelumAda - tabel, view, atau sequence tidak ada di skema ini (ORA-00942 / ORA-00904 / ORA-02289).
	ErrBelumAda = errors.New("repository: tabel, view, atau sequence R/I Rate Life tidak ada di skema ini")
	// ErrKembar - ORA-00001.
	ErrKembar = errors.New("repository: ID sudah dipakai")
	// ErrBacaSaja - SQL tulis diarahkan ke view atau objek di luar DaftarTabelDitulis.
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

// nama - nama berskema setiap objek.
type nama struct{ tabelRingkasan, tabelRate, viewRingkasan, viewRate, seqRingkasan, seqRate string }

func (g *Gudang) nama() (nama, error) {
	var n nama
	for _, p := range []struct {
		ke    *string
		objek string
	}{{&n.tabelRingkasan, TabelRingkasan}, {&n.tabelRate, TabelRate}, {&n.viewRingkasan, ViewRingkasan},
		{&n.viewRate, ViewRate}, {&n.seqRingkasan, SeqRingkasan}, {&n.seqRate, SeqRate}} {
		q, err := g.db.Qualify(p.objek)
		if err != nil {
			return nama{}, err
		}
		*p.ke = q
	}
	return n, nil
}

func siap(objek, q string) error {
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	return PeriksaTulis(objek, q)
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

// Urutan - klausa ORDER BY grid ringkasan (kolom dari daftar putih, bukan dari masukan).
func Urutan(urut string, turun bool) string {
	arah := "ASC"
	if turun {
		arah = "DESC"
	}
	switch urut {
	case models.UrutUsedBy:
		return fmt.Sprintf("UPPER(USEDBY) %s NULLS LAST, ID", arah)
	case models.UrutOperator:
		return fmt.Sprintf("UPPER(OPERATORID) %s NULLS LAST, ID", arah)
	case models.UrutTanggal:
		return fmt.Sprintf("MODIFIEDDATE %s NULLS LAST, ID", arah)
	}
	return fmt.Sprintf("TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) %s NULLS LAST, ID %s", arah, arah)
}

const saringRingkasan = `(:1 IS NULL OR UPPER(ID) LIKE :2 ESCAPE '\') AND (:3 IS NULL OR UPPER(USEDBY) LIKE :4 ESCAPE '\')`

// SqlDaftar - satu halaman grid `BrowseRateLifeSummary` (param `id`, `idusedby` b9308-b9314).
func SqlDaftar(v, urut string, turun bool) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE %s ORDER BY %s OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY`,
		kolomRingkasan, v, saringRingkasan, Urutan(urut, turun))
}

// SqlJumlah - jumlah baris bersaring.
func SqlJumlah(v string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, v, saringRingkasan)
}

// SqlAmbil - satu ringkasan.
func SqlAmbil(v string) string { return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, kolomRingkasan, v) }

// SqlPemakaiNama - ringkasan lain bernama sama (tanpa beda huruf dan spasi tepi), selain kecualiID.
func SqlPemakaiNama(v string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE UPPER(TRIM(USEDBY)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0)) ORDER BY ID`,
		kolomRingkasan, v)
}

// SqlIDBaru - satu nomor sequence sebagai teks.
func SqlIDBaru(seq string) string { return fmt.Sprintf(`SELECT TO_CHAR(%s.NEXTVAL) FROM DUAL`, seq) }

// SqlMaksID - nomor ID angka tertinggi tabel fisik (ID non-angka diabaikan) - bentuk sama dengan migrasi 923.
func SqlMaksID(t string) string {
	return fmt.Sprintf(`SELECT TO_CHAR(NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$'))), 0)) FROM %s`, t)
}

// SqlAdaID - ID sudah terpakai di tabel fisik?
func SqlAdaID(t string) string { return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ID = :1`, t) }

// SqlSisipRingkasan - ringkasan baru (`JSON_OBJECT`, kunci ASUMSI A1).
func SqlSisipRingkasan(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, %s) VALUES (:1, JSON_OBJECT('%s' VALUE :2, '%s' VALUE :3, '%s' VALUE :4 ABSENT ON NULL))`,
		t, KolomJSON, JSONUsedBy, JSONOperatorID, JSONModified)
}

// SqlUbahRingkasan - Edit: hanya tiga kunci yang diganti; kunci lain milik Pega tetap.
func SqlUbahRingkasan(t string) string {
	return fmt.Sprintf(`UPDATE %s SET %s = JSON_MERGEPATCH(%s, JSON_OBJECT('%s' VALUE :1, '%s' VALUE :2, '%s' VALUE :3) RETURNING CLOB)
	  WHERE ID = :4`, t, KolomJSON, KolomJSON, JSONUsedBy, JSONOperatorID, JSONModified)
}

// SqlUbahNamaRate - salinan nama di baris rate ringkasan itu ikut diganti (kolom `USEDBY` view `RATE_LIFE`).
func SqlUbahNamaRate(t, v string) string {
	return fmt.Sprintf(`UPDATE %s SET %s = JSON_MERGEPATCH(%s, JSON_OBJECT('%s' VALUE :1) RETURNING CLOB)
	  WHERE ID IN (SELECT ID FROM %s WHERE IDUSEDBY = :2)`, t, KolomJSON, KolomJSON, JSONUsedBy, v)
}

// SqlHapusRingkasan - Delete ringkasan.
func SqlHapusRingkasan(t string) string { return fmt.Sprintf(`DELETE FROM %s WHERE ID = :1`, t) }

// SqlHapusRate - baris rate milik ringkasan (`DeleteSummaryDetail`: ringkasan BESERTA rinciannya).
func SqlHapusRate(t, v string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE ID IN (SELECT ID FROM %s WHERE IDUSEDBY = :1)`, t, v)
}

// SqlJumlahRate - jumlah baris rate milik ringkasan.
func SqlJumlahRate(v string) string { return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE IDUSEDBY = :1`, v) }

// urutRate - urutan stabil grid Rate Detail: GENDER, CONTRACT, AGE (angka; kosong dulu), ID.
const urutRate = `GENDER, TO_NUMBER(REGEXP_SUBSTR(TRIM(CONTRACT), '^[0-9]+$')) NULLS FIRST,
	  TO_NUMBER(REGEXP_SUBSTR(TRIM(AGE), '^[0-9]+$')) NULLS FIRST, ID`

// SqlDaftarRate - satu halaman Rate Detail (harness `InboxRIRate` b11444).
func SqlDaftarRate(v string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE IDUSEDBY = :1 ORDER BY %s OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY`, kolomRate, v, urutRate)
}

// SqlRateDari - seluruh baris rate untuk n ringkasan sekaligus (kunci kembar upload); n >= 1.
func SqlRateDari(v string, n int) string {
	ikat := make([]string, n)
	for i := range ikat {
		ikat[i] = fmt.Sprintf(":%d", i+1)
	}
	return fmt.Sprintf(`SELECT %s FROM %s WHERE IDUSEDBY IN (%s)`, kolomRate, v, strings.Join(ikat, ", "))
}

// SqlSisipRate - baris rate baru (`JSON_OBJECT`; TYPE tidak diisi; CONTRACT kosong = kunci tidak ditulis).
func SqlSisipRate(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, %s) VALUES (:1, JSON_OBJECT('%s' VALUE :2, '%s' VALUE :3, '%s' VALUE :4,
	  '%s' VALUE :5, '%s' VALUE :6, '%s' VALUE :7 ABSENT ON NULL))`,
		t, KolomJSON, JSONIDUsedBy, JSONUsedBy, JSONGender, JSONContract, JSONAge, JSONRate)
}

type pemindai interface{ Scan(...any) error }

func pindaiTeks(p pemindai, n int) ([]string, error) {
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

func (g *Gudang) bacaBaris(ctx context.Context, tx *db.Tx, objek, q string, kolom int, args ...any) ([][]string, error) {
	if err := siap(objek, q); err != nil {
		return nil, err
	}
	rows, err := g.dari(tx).QueryContext(ctx, q, args...)
	if err != nil {
		return nil, bungkus(err, "membaca")
	}
	defer func() { _ = rows.Close() }()
	var out [][]string
	for rows.Next() {
		s, err := pindaiTeks(rows, kolom)
		if err != nil {
			return nil, bungkus(err, "memindai")
		}
		out = append(out, s)
	}
	return out, bungkus(rows.Err(), "membaca")
}

func (g *Gudang) satuNilai(ctx context.Context, tx *db.Tx, objek, q string, args ...any) (string, error) {
	if err := siap(objek, q); err != nil {
		return "", err
	}
	var s sql.NullString
	if err := g.dari(tx).QueryRowContext(ctx, q, args...).Scan(&s); err != nil {
		return "", bungkus(err, "membaca")
	}
	return s.String, nil
}

func (g *Gudang) tulis(ctx context.Context, tx *db.Tx, objek, q, apa string, args ...any) (int64, error) {
	if err := siap(objek, q); err != nil {
		return 0, err
	}
	h, err := g.dari(tx).ExecContext(ctx, q, args...)
	if err != nil {
		return 0, bungkus(err, apa)
	}
	n, _ := h.RowsAffected()
	return n, nil
}

func keRingkasan(b [][]string) []models.Ringkasan {
	out := make([]models.Ringkasan, 0, len(b))
	for _, s := range b {
		out = append(out, models.Ringkasan{ID: s[0], UsedBy: s[1], OperatorID: s[2], ModifiedDate: s[3]})
	}
	return out
}

func keRate(b [][]string) []models.Rate {
	out := make([]models.Rate, 0, len(b))
	for _, s := range b {
		out = append(out, models.Rate{ID: s[0], IDUsedBy: s[1], UsedBy: s[2], Gender: s[3], Contract: s[4], Age: s[5], Rate: s[6]})
	}
	return out
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

// Daftar - satu halaman ringkasan bersaring dan total barisnya.
func (g *Gudang) Daftar(ctx context.Context, s models.Saringan) ([]models.Ringkasan, int, error) {
	n, err := g.nama()
	if err != nil {
		return nil, 0, err
	}
	id, nm := PolaCari(s.ID), PolaCari(s.UsedBy)
	total, err := g.satuNilai(ctx, nil, ViewRingkasan, SqlJumlah(n.viewRingkasan), id, id, nm, nm)
	if err != nil {
		return nil, 0, err
	}
	b, err := g.bacaBaris(ctx, nil, ViewRingkasan, SqlDaftar(n.viewRingkasan, s.Urut, s.Turun), 4, id, id, nm, nm,
		(s.Halaman-1)*models.UkuranHalaman, models.UkuranHalaman)
	return keRingkasan(b), angka(total), err
}

// Ambil - satu ringkasan.
func (g *Gudang) Ambil(ctx context.Context, tx *db.Tx, id string) (models.Ringkasan, error) {
	n, err := g.nama()
	if err != nil {
		return models.Ringkasan{}, err
	}
	b, err := g.bacaBaris(ctx, tx, ViewRingkasan, SqlAmbil(n.viewRingkasan), 4, id)
	if err != nil {
		return models.Ringkasan{}, err
	}
	if len(b) == 0 {
		return models.Ringkasan{}, ErrTidakAda
	}
	return keRingkasan(b)[0], nil
}

// PemakaiNama - ringkasan bernama sama selain kecualiID ("" = semua).
func (g *Gudang) PemakaiNama(ctx context.Context, tx *db.Tx, nama, kecualiID string) ([]models.Ringkasan, error) {
	n, err := g.nama()
	if err != nil {
		return nil, err
	}
	b, err := g.bacaBaris(ctx, tx, ViewRingkasan, SqlPemakaiNama(n.viewRingkasan), 4, nama, db.KosongJadiNil(kecualiID))
	return keRingkasan(b), err
}

// IDBaru - satu nomor sequence (`ringkasan` true = SEQ_M_RATE_LIFE_SUMMARY, selain itu SEQ_M_RATE_LIFE).
func (g *Gudang) IDBaru(ctx context.Context, tx *db.Tx, ringkasan bool) (string, error) {
	n, err := g.nama()
	if err != nil {
		return "", err
	}
	seq := n.seqRate
	if ringkasan {
		seq = n.seqRingkasan
	}
	return g.satuNilai(ctx, tx, "DUAL", SqlIDBaru(seq))
}

func (n nama) tabel(ringkasan bool) (string, string) {
	if ringkasan {
		return TabelRingkasan, n.tabelRingkasan
	}
	return TabelRate, n.tabelRate
}

// MaksID - nomor ID angka tertinggi tabel fisik.
func (g *Gudang) MaksID(ctx context.Context, tx *db.Tx, ringkasan bool) (int, error) {
	n, err := g.nama()
	if err != nil {
		return 0, err
	}
	objek, t := n.tabel(ringkasan)
	s, err := g.satuNilai(ctx, tx, objek, SqlMaksID(t))
	return angka(s), err
}

// AdaID - ID sudah terpakai di tabel fisik?
func (g *Gudang) AdaID(ctx context.Context, tx *db.Tx, ringkasan bool, id string) (bool, error) {
	n, err := g.nama()
	if err != nil {
		return false, err
	}
	objek, t := n.tabel(ringkasan)
	s, err := g.satuNilai(ctx, tx, objek, SqlAdaID(t), id)
	return angka(s) > 0, err
}

// SisipRingkasan - ringkasan baru.
func (g *Gudang) SisipRingkasan(ctx context.Context, tx *db.Tx, r models.Ringkasan) error {
	n, err := g.nama()
	if err != nil {
		return err
	}
	_, err = g.tulis(ctx, tx, TabelRingkasan, SqlSisipRingkasan(n.tabelRingkasan), "menyimpan ringkasan", r.ID,
		db.KosongJadiNil(r.UsedBy), db.KosongJadiNil(r.OperatorID), db.KosongJadiNil(r.ModifiedDate))
	return err
}

// UbahRingkasan - Edit; ErrTidakAda bila ID tidak ada.
func (g *Gudang) UbahRingkasan(ctx context.Context, tx *db.Tx, r models.Ringkasan) error {
	n, err := g.nama()
	if err != nil {
		return err
	}
	j, err := g.tulis(ctx, tx, TabelRingkasan, SqlUbahRingkasan(n.tabelRingkasan), "mengubah ringkasan",
		db.KosongJadiNil(r.UsedBy), db.KosongJadiNil(r.OperatorID), db.KosongJadiNil(r.ModifiedDate), r.ID)
	if err == nil && j == 0 {
		return ErrTidakAda
	}
	return err
}

// UbahNamaRate - salinan nama di baris rate ringkasan idUsedBy.
func (g *Gudang) UbahNamaRate(ctx context.Context, tx *db.Tx, idUsedBy, nama string) (int, error) {
	n, err := g.nama()
	if err != nil {
		return 0, err
	}
	j, err := g.tulis(ctx, tx, TabelRate, SqlUbahNamaRate(n.tabelRate, n.viewRate), "mengubah nama rate", nama, idUsedBy)
	return int(j), err
}

// HapusRingkasan - Delete ringkasan; ErrTidakAda bila ID tidak ada.
func (g *Gudang) HapusRingkasan(ctx context.Context, tx *db.Tx, id string) error {
	n, err := g.nama()
	if err != nil {
		return err
	}
	j, err := g.tulis(ctx, tx, TabelRingkasan, SqlHapusRingkasan(n.tabelRingkasan), "menghapus ringkasan", id)
	if err == nil && j == 0 {
		return ErrTidakAda
	}
	return err
}

// HapusRate - baris rate milik ringkasan; jumlah terhapus.
func (g *Gudang) HapusRate(ctx context.Context, tx *db.Tx, idUsedBy string) (int, error) {
	n, err := g.nama()
	if err != nil {
		return 0, err
	}
	j, err := g.tulis(ctx, tx, TabelRate, SqlHapusRate(n.tabelRate, n.viewRate), "menghapus rate", idUsedBy)
	return int(j), err
}

// JumlahRate - baris rate milik ringkasan.
func (g *Gudang) JumlahRate(ctx context.Context, tx *db.Tx, idUsedBy string) (int, error) {
	n, err := g.nama()
	if err != nil {
		return 0, err
	}
	s, err := g.satuNilai(ctx, tx, ViewRate, SqlJumlahRate(n.viewRate), idUsedBy)
	return angka(s), err
}

// DaftarRate - satu halaman Rate Detail dan totalnya.
func (g *Gudang) DaftarRate(ctx context.Context, idUsedBy string, halaman int) ([]models.Rate, int, error) {
	n, err := g.nama()
	if err != nil {
		return nil, 0, err
	}
	total, err := g.JumlahRate(ctx, nil, idUsedBy)
	if err != nil {
		return nil, 0, err
	}
	b, err := g.bacaBaris(ctx, nil, ViewRate, SqlDaftarRate(n.viewRate), 7, idUsedBy, (halaman-1)*models.UkuranHalaman,
		models.UkuranHalaman)
	return keRate(b), total, err
}

// RateDari - seluruh baris rate untuk ringkasan ids (kunci kembar upload). ids kosong = nil.
func (g *Gudang) RateDari(ctx context.Context, tx *db.Tx, ids []string) ([]models.Rate, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	n, err := g.nama()
	if err != nil {
		return nil, err
	}
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	b, err := g.bacaBaris(ctx, tx, ViewRate, SqlRateDari(n.viewRate, len(ids)), 7, args...)
	return keRate(b), err
}

// SisipRate - baris rate baru.
func (g *Gudang) SisipRate(ctx context.Context, tx *db.Tx, r models.Rate) error {
	n, err := g.nama()
	if err != nil {
		return err
	}
	k := db.KosongJadiNil
	_, err = g.tulis(ctx, tx, TabelRate, SqlSisipRate(n.tabelRate), "menyimpan rate", r.ID, k(r.IDUsedBy), k(r.UsedBy),
		k(r.Gender), k(r.Contract), k(r.Age), k(r.Rate))
	return err
}
