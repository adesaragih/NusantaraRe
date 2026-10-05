// Package repository adalah SATU-SATUNYA lapisan modul Aggregate yang berbicara ke Oracle.
//
// `AGGREGATE` ditulis (sisip, hapus) dan dibaca; view `TREATYINDETAILJOINEDM` (Master ID), `ASSESSMENT_ZONE`,
// `TREATYYEAR`, dan `TREATYEXCHANGEYEARLY` hanya dibaca - semuanya tabel warisan Pega yang strukturnya tidak diubah.
// Nomor ID baru dari `SEQ_AGGREGATE` (migrasi 880).
//
// Setiap tabel lewat `Qualify` (ADR-U-0033), nilai lewat bind, nol COMMIT di SQL (ADR-U-0029): transaksinya milik
// services. Angka dibaca sebagai teks (`db.FmtDesimal`) dan ditulis `TO_NUMBER(:koef) / POWER(10, :skala)` - nol
// float, kebal NLS sesi.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/aggregate/backend/models"
)

// Objek Oracle modul ini.
const (
	TabelAggregate = "AGGREGATE"
	ViewTreaty     = "TREATYINDETAILJOINEDM"
	TabelZona      = "ASSESSMENT_ZONE"
	TabelTahun     = "TREATYYEAR"
	TabelKurs      = "TREATYEXCHANGEYEARLY"
	SeqAggregate   = "SEQ_AGGREGATE"
)

// Saringan Master ID (`GetMasterIDAgg_Act` langkah 3, Logic `A AND ((B AND C) OR D)`).
const (
	ProporsionalPropertyJenis = "Proportional"
	ProporsionalPropertyGrup  = "PROPERTY"
	NonProporsional           = "NonProportional"
)

var (
	// ErrTidakAda - baris yang dicari tidak ada.
	ErrTidakAda = errors.New("repository: row not found")
	// ErrBelumDimigrasi - SEQ_AGGREGATE (migrasi 880) belum ada.
	ErrBelumDimigrasi = errors.New("repository: SEQ_AGGREGATE is not migrated yet - run -migrate (migration 880)")
)

// Gudang membaca dan menulis Oracle.
type Gudang struct{ db *db.DB }

// Baru menyusun Gudang.
func Baru(d *db.DB) *Gudang { return &Gudang{db: d} }

// penjalan - *db.DB dan *db.Tx.
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

// belumDimigrasi - ORA-02289: sequence tidak ada.
func belumDimigrasi(err error) bool { return err != nil && strings.Contains(err.Error(), "ORA-02289") }

func bungkus(err error, apa string) error {
	if belumDimigrasi(err) {
		return ErrBelumDimigrasi
	}
	return fmt.Errorf("repository: %s: %w", apa, err)
}

// PolaCari membentuk pola LIKE huruf besar yang aman: `\`, `%`, dan `_` ketikan pemakai bukan wildcard.
func PolaCari(kueri string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return "%" + r.Replace(strings.ToUpper(strings.TrimSpace(kueri))) + "%"
}

// --- SQL ---------------------------------------------------------------------------------------------------------

// ekspresiBaca - ekspresi SELECT satu kolom grid.
func ekspresiBaca(k models.Kolom) string {
	switch k.Jenis {
	case models.Angka:
		return fmt.Sprintf(db.FmtDesimal, k.Nama)
	case models.Tanggal:
		return fmt.Sprintf(`TO_CHAR(%s, 'DD-MM-YYYY')`, k.Nama)
	}
	return k.Nama
}

// kolomBaca - ID, TANGGAL_INPUT, USER_INPUT, lalu kolom grid berurutan (`pindaiBaris`).
func kolomBaca() string {
	bagian := []string{models.KolomID, `TO_CHAR(TANGGAL_INPUT, 'DD-MM-YYYY HH24:MI:SS')`, models.KolomUserInput}
	for _, k := range models.KolomGrid {
		bagian = append(bagian, ekspresiBaca(k))
	}
	return strings.Join(bagian, ", ")
}

// sqlSisip - satu baris: ID, TANGGAL_INPUT, USER_INPUT, lalu kolom grid; angka dua bind (koefisien, skala).
// COMMENCEMENT tidak diisi (`UploadCSVAggregate_Act` tidak pernah mengisinya; keputusan work owner 04-10-2026).
func sqlSisip(t string) string {
	kolom := []string{models.KolomID, models.KolomTanggalInput, models.KolomUserInput}
	nilai := []string{":1", "TO_DATE(:2, 'YYYY-MM-DD HH24:MI:SS')", ":3"}
	n := 4
	for _, k := range models.KolomGrid {
		kolom = append(kolom, k.Nama)
		switch k.Jenis {
		case models.Angka:
			nilai = append(nilai, fmt.Sprintf("(TO_NUMBER(:%d) / POWER(10, :%d))", n, n+1))
			n += 2
		case models.Tanggal:
			nilai = append(nilai, fmt.Sprintf("TO_DATE(:%d, 'DD-MM-YYYY')", n))
			n++
		default:
			nilai = append(nilai, fmt.Sprintf(":%d", n))
			n++
		}
	}
	return fmt.Sprintf(`INSERT INTO %s (%s)
	  VALUES (%s)`, t, strings.Join(kolom, ", "), strings.Join(nilai, ", "))
}

// kunciWhere - baris AGGREGATE berkunci daftar; DECODE menyamakan NULL dengan NULL. :1-:6.
const kunciWhere = `DECODE(TRUNC(TANGGAL_INPUT), TO_DATE(:1, 'DD-MM-YYYY'), 1, 0) = 1
	    AND DECODE(CEDING_CODE, :2, 1, 0) = 1 AND DECODE(CEDING_NAME, :3, 1, 0) = 1
	    AND DECODE(TREATY_TYPE, :4, 1, 0) = 1 AND DECODE(AS_AT, TO_DATE(:5, 'DD-MM-YYYY'), 1, 0) = 1
	    AND DECODE(UW_YEAR, :6, 1, 0) = 1`

// saringDaftar - :1 kosong (NULL) = tanpa saringan; :2-:3 pola yang sama untuk CEDING_NAME dan CEDING_CODE.
const saringDaftar = `(:1 IS NULL OR UPPER(CEDING_NAME) LIKE :2 ESCAPE '\' OR UPPER(CEDING_CODE) LIKE :3 ESCAPE '\')`

// kelompokDaftar - kolom kunci daftar `GridDasbordAgg`.
const kelompokDaftar = `TRUNC(TANGGAL_INPUT), CEDING_CODE, CEDING_NAME, TREATY_TYPE, AS_AT, UW_YEAR`

func sqlDaftar(t string) string {
	return fmt.Sprintf(`SELECT TO_CHAR(TRUNC(TANGGAL_INPUT), 'DD-MM-YYYY'), CEDING_CODE, CEDING_NAME, TREATY_TYPE,
	    TO_CHAR(AS_AT, 'DD-MM-YYYY'), UW_YEAR, COUNT(*), TO_CHAR(MAX(TANGGAL_INPUT), 'DD-MM-YYYY HH24:MI')
	  FROM %s WHERE %s
	  GROUP BY %s
	  ORDER BY MAX(TANGGAL_INPUT) DESC NULLS LAST, CEDING_NAME, TREATY_TYPE, AS_AT, UW_YEAR
	  OFFSET :4 ROWS FETCH NEXT :5 ROWS ONLY`, t, saringDaftar, kelompokDaftar)
}

func sqlHitungDaftar(t string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM (SELECT 1 FROM %s WHERE %s GROUP BY %s)`, t, saringDaftar, kelompokDaftar)
}

func sqlRincian(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE %s
	  ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '[0-9]+$')) NULLS LAST, ID`, kolomBaca(), t, kunciWhere)
}

func sqlHapus(t string) string { return fmt.Sprintf(`DELETE FROM %s WHERE %s`, t, kunciWhere) }

// sqlRingkasan - jumlah RNM_VALUE_IN_USD per CEDING, TREATY_TYPE, COVERAGE; :1 kosong = seluruh As At.
func sqlRingkasan(t string) string {
	return fmt.Sprintf(`SELECT CEDING_CODE, CEDING_NAME, TREATY_TYPE, COVERAGE, %s FROM %s
	  WHERE (:1 IS NULL OR AS_AT = TO_DATE(:2, 'DD-MM-YYYY'))
	  GROUP BY CEDING_CODE, CEDING_NAME, TREATY_TYPE, COVERAGE
	  ORDER BY CEDING_CODE, CEDING_NAME, TREATY_TYPE, COVERAGE`, fmt.Sprintf(db.FmtDesimal, "SUM(RNM_VALUE_IN_USD)"), t)
}

func sqlNomor(seq string) string { return fmt.Sprintf(`SELECT %s.NEXTVAL FROM DUAL`, seq) }

func sqlAdaID(t string) string { return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ID = :1`, t) }

const sqlSekarang = `SELECT TO_CHAR(SYSDATE, 'YYYY-MM-DD HH24:MI:SS') FROM DUAL`

// kolomTreaty - kolom view Master ID (`pindaiTreaty`).
var kolomTreaty = `ID, TREATYID, CEDINGID, CEDING, ` + fmt.Sprintf(db.FmtDesimal, "RNM_SHARE") +
	`, TREATYYEAR, PROPORTIONTYPE, TREATYCONTRACTNAME, TREATYGROUP, SOB`

// sqlCariTreaty - `GetMasterIDAgg_Act` langkah 3: CEDING memuat kata cari (huruf besar), lalu (Proportional dan
// PROPERTY) atau NonProportional. Tanpa urutan, seperti Obj-Browse Pega.
func sqlCariTreaty(v string) string {
	return fmt.Sprintf(`SELECT %s FROM %s
	  WHERE CEDING LIKE :1 ESCAPE '\' AND ((PROPORTIONTYPE = :2 AND TREATYGROUP = :3) OR PROPORTIONTYPE = :4)
	  FETCH FIRST :5 ROWS ONLY`, kolomTreaty, v)
}

func sqlAmbilTreaty(v string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, kolomTreaty, v)
}

func sqlDaftarZona(t string) string {
	return fmt.Sprintf(`SELECT ASSESSMENT_CODE, ASSESSMENT_NOTE FROM %s`, t)
}

// sqlTahunTreaty - `UploadCSVAggregate_Act` langkah 7.9: periode TREATYYEAR yang memuat AS_AT (teks yyyyMMdd),
// baris pertama tanpa urutan seperti Obj-Browse Pega (keputusan work owner 04-10-2026: ikuti XML).
func sqlTahunTreaty(t string) string {
	return fmt.Sprintf(`SELECT TREATYYEAR FROM %s WHERE STARTDATE <= :1 AND ENDDATE >= :2 FETCH FIRST 1 ROWS ONLY`, t)
}

// sqlKurs - langkah 7.14: kurs tahun treaty dan mata uang, baris pertama tanpa urutan.
func sqlKurs(t string) string {
	return fmt.Sprintf(`SELECT TOUSD, TOIDR FROM %s WHERE TREATYYEAR = :1 AND CURRENCY = :2 FETCH FIRST 1 ROWS ONLY`, t)
}

// --- pemindai dan pelaksana -------------------------------------------------------------------------------------

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

func pindaiBaris(p pemindai) (models.Baris, error) {
	v, err := pindaiTeks(p, 3+len(models.KolomGrid))
	if err != nil {
		return nil, err
	}
	b := models.Baris{models.KolomID: v[0], models.KolomTanggalInput: v[1], models.KolomUserInput: v[2]}
	for i, k := range models.KolomGrid {
		b[k.Nama] = v[3+i]
	}
	return b, nil
}

func pindaiTreaty(p pemindai) (models.MasterTreaty, error) {
	v, err := pindaiTeks(p, 10)
	if err != nil {
		return models.MasterTreaty{}, err
	}
	return models.MasterTreaty{ID: v[0], TreatyID: v[1], CedingID: v[2], Ceding: v[3], RnmShare: v[4], TreatyYear: v[5],
		ProportionType: v[6], TreatyContractName: v[7], TreatyGroup: v[8], SOB: v[9]}, nil
}

func pindaiKelompok(p pemindai) (models.Kelompok, error) {
	var tgl, kode, nama, jenis, asAt, uw, akhir sql.NullString
	var n int
	if err := p.Scan(&tgl, &kode, &nama, &jenis, &asAt, &uw, &n, &akhir); err != nil {
		return models.Kelompok{}, err
	}
	return models.Kelompok{Kunci: models.Kunci{TanggalInput: tgl.String, CedingCode: kode.String, CedingName: nama.String,
		TreatyType: jenis.String, AsAt: asAt.String, UwYear: uw.String}, JumlahBaris: n, InputTerakhir: akhir.String}, nil
}

func pindaiBilangan(p pemindai) (int64, error) {
	var n int64
	err := p.Scan(&n)
	return n, err
}

func pindaiSatuTeks(p pemindai) (string, error) {
	var s sql.NullString
	err := p.Scan(&s)
	return s.String, err
}

// daftar menjalankan satu SELECT banyak baris.
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

// satu menjalankan satu SELECT satu baris; nol baris = ErrTidakAda.
func satu[T any](ctx context.Context, j penjalan, q string, pindai func(pemindai) (T, error), args ...any) (T, error) {
	var nol T
	if err := db.PeriksaSQL(q); err != nil {
		return nol, err
	}
	v, err := pindai(j.QueryRowContext(ctx, q, args...))
	if errors.Is(err, sql.ErrNoRows) {
		return nol, ErrTidakAda
	}
	if err != nil {
		return nol, bungkus(err, "membaca")
	}
	return v, nil
}

// jalankan menjalankan satu DML.
func jalankan(ctx context.Context, j penjalan, q, apa string, args ...any) (sql.Result, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	hasil, err := j.ExecContext(ctx, q, args...)
	if err != nil {
		return nil, bungkus(err, apa)
	}
	return hasil, nil
}

// --- nilai bind ---------------------------------------------------------------------------------------------------

// PecahAngka memecah desimal menjadi koefisien bulat (teks) dan skala; nil = NULL.
// Contoh: 1500000000.10 -> ("150000000010", 2).
func PecahAngka(d *apd.Decimal) (any, int64) {
	if d == nil {
		return nil, 0
	}
	s := d.Text('f')
	tanda := ""
	if strings.HasPrefix(s, "-") {
		tanda, s = "-", s[1:]
	}
	bulat, pecahan, _ := strings.Cut(s, ".")
	koef := strings.TrimLeft(bulat+pecahan, "0")
	if koef == "" {
		return "0", 0
	}
	return tanda + koef, int64(len(pecahan))
}

// NilaiSisip - nilai bind `sqlSisip`, berurutan. Baris sudah diperiksa services: angka terurai, tanggal sah;
// teks kosong = NULL.
func NilaiSisip(id, sekarang, pelaku string, b models.Baris) ([]any, error) {
	nilai := []any{id, sekarang, db.KosongJadiNil(pelaku)}
	for _, k := range models.KolomGrid {
		v := strings.TrimSpace(b[k.Nama])
		switch k.Jenis {
		case models.Angka:
			var d *apd.Decimal
			if v != "" {
				var err error
				if d, err = utils.ParseDecimal(v); err != nil {
					return nil, fmt.Errorf("repository: %s: %w", k.Nama, err)
				}
			}
			koef, skala := PecahAngka(d)
			nilai = append(nilai, koef, skala)
		default:
			nilai = append(nilai, db.KosongJadiNil(v))
		}
	}
	return nilai, nil
}

// NilaiKunci - enam nilai bind `kunciWhere`; kosong = NULL.
func NilaiKunci(k models.Kunci) []any {
	n := db.KosongJadiNil
	return []any{n(k.TanggalInput), n(k.CedingCode), n(k.CedingName), n(k.TreatyType), n(k.AsAt), n(k.UwYear)}
}

// NilaiCari - tiga nilai bind `saringDaftar`; kueri kosong = tanpa saringan.
func NilaiCari(kueri string) []any {
	if strings.TrimSpace(kueri) == "" {
		return []any{nil, nil, nil}
	}
	p := PolaCari(kueri)
	return []any{p, p, p}
}

// --- operasi ------------------------------------------------------------------------------------------------------

// Daftar membaca satu halaman daftar dan jumlah seluruhnya.
func (g *Gudang) Daftar(ctx context.Context, kueri string, offset, ukuran int) ([]models.Kelompok, int, error) {
	t, err := g.nama(TabelAggregate)
	if err != nil {
		return nil, 0, err
	}
	cari := NilaiCari(kueri)
	jumlah, err := satu(ctx, g.db, sqlHitungDaftar(t), pindaiBilangan, cari...)
	if err != nil {
		return nil, 0, err
	}
	baris, err := daftar(ctx, g.db, sqlDaftar(t), pindaiKelompok, append(cari, offset, ukuran)...)
	return baris, int(jumlah), err
}

// Rincian membaca seluruh baris AGGREGATE berkunci k.
func (g *Gudang) Rincian(ctx context.Context, k models.Kunci) ([]models.Baris, error) {
	t, err := g.nama(TabelAggregate)
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.db, sqlRincian(t), pindaiBaris, NilaiKunci(k)...)
}

// Hapus membuang seluruh baris AGGREGATE berkunci k; jumlah baris yang terhapus.
func (g *Gudang) Hapus(ctx context.Context, tx *db.Tx, k models.Kunci) (int64, error) {
	t, err := g.nama(TabelAggregate)
	if err != nil {
		return 0, err
	}
	hasil, err := jalankan(ctx, g.dari(tx), sqlHapus(t), "menghapus aggregate", NilaiKunci(k)...)
	if err != nil {
		return 0, err
	}
	return hasil.RowsAffected()
}

// Ringkasan - jumlah RNM Value (USD) per Ceding, Treaty Type, Coverage; asAt kosong = seluruhnya.
func (g *Gudang) Ringkasan(ctx context.Context, asAt string) ([]models.IrisanRingkasan, error) {
	t, err := g.nama(TabelAggregate)
	if err != nil {
		return nil, err
	}
	a := db.KosongJadiNil(asAt)
	return daftar(ctx, g.db, sqlRingkasan(t), func(p pemindai) (models.IrisanRingkasan, error) {
		v, err := pindaiTeks(p, 5)
		if err != nil {
			return models.IrisanRingkasan{}, err
		}
		return models.IrisanRingkasan{CedingCode: v[0], CedingName: v[1], TreatyType: v[2], Coverage: v[3], RnmValueInUSD: v[4]}, nil
	}, a, a)
}

// CariTreaty - baris view Master ID yang CEDING-nya memuat pola, paling banyak batas.
func (g *Gudang) CariTreaty(ctx context.Context, kueri string, batas int) ([]models.MasterTreaty, error) {
	v, err := g.nama(ViewTreaty)
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.db, sqlCariTreaty(v), pindaiTreaty, PolaCari(kueri), ProporsionalPropertyJenis,
		ProporsionalPropertyGrup, NonProporsional, batas)
}

// AmbilTreaty - satu baris view Master ID; tidak ada = ErrTidakAda.
func (g *Gudang) AmbilTreaty(ctx context.Context, id string) (models.MasterTreaty, error) {
	v, err := g.nama(ViewTreaty)
	if err != nil {
		return models.MasterTreaty{}, err
	}
	return satu(ctx, g.db, sqlAmbilTreaty(v), pindaiTreaty, id)
}

// DaftarZona - kode -> catatan ASSESSMENT_ZONE (kode pertama yang terbaca menang, seperti pxResults(1) Pega).
func (g *Gudang) DaftarZona(ctx context.Context) (map[string]string, error) {
	t, err := g.nama(TabelZona)
	if err != nil {
		return nil, err
	}
	baris, err := daftar(ctx, g.db, sqlDaftarZona(t), func(p pemindai) ([]string, error) { return pindaiTeks(p, 2) })
	if err != nil {
		return nil, err
	}
	out := map[string]string{}
	for _, b := range baris {
		if _, ada := out[b[0]]; !ada {
			out[b[0]] = b[1]
		}
	}
	return out, nil
}

// TahunTreaty - tahun treaty periode yang memuat asAt (`yyyyMMdd`); tidak ada = "".
func (g *Gudang) TahunTreaty(ctx context.Context, asAt string) (string, error) {
	t, err := g.nama(TabelTahun)
	if err != nil {
		return "", err
	}
	s, err := satu(ctx, g.db, sqlTahunTreaty(t), pindaiSatuTeks, asAt, asAt)
	if errors.Is(err, ErrTidakAda) {
		return "", nil
	}
	return s, err
}

// Kurs - TOUSD dan TOIDR tahun treaty dan mata uang; ada false = tidak ada baris.
func (g *Gudang) Kurs(ctx context.Context, tahun, mataUang string) (toUSD, toIDR string, ada bool, err error) {
	t, err := g.nama(TabelKurs)
	if err != nil {
		return "", "", false, err
	}
	v, err := satu(ctx, g.db, sqlKurs(t), func(p pemindai) ([]string, error) { return pindaiTeks(p, 2) }, tahun, mataUang)
	if errors.Is(err, ErrTidakAda) {
		return "", "", false, nil
	}
	if err != nil {
		return "", "", false, err
	}
	return v[0], v[1], true, nil
}

// Sekarang - waktu basis data `YYYY-MM-DD HH24:MI:SS`; satu nilai TANGGAL_INPUT untuk seluruh unggahan.
func (g *Gudang) Sekarang(ctx context.Context, tx *db.Tx) (string, error) {
	return satu(ctx, g.dari(tx), sqlSekarang, pindaiSatuTeks)
}

// NomorBerikut - `SEQ_AGGREGATE.NEXTVAL`.
func (g *Gudang) NomorBerikut(ctx context.Context, tx *db.Tx) (int64, error) {
	s, err := g.nama(SeqAggregate)
	if err != nil {
		return 0, err
	}
	return satu(ctx, g.dari(tx), sqlNomor(s), pindaiBilangan)
}

// AdaID - ID sudah terpakai (mis. dibuat Pega).
func (g *Gudang) AdaID(ctx context.Context, tx *db.Tx, id string) (bool, error) {
	t, err := g.nama(TabelAggregate)
	if err != nil {
		return false, err
	}
	n, err := satu(ctx, g.dari(tx), sqlAdaID(t), pindaiBilangan, id)
	return n > 0, err
}

// Sisip menyisipkan satu baris AGGREGATE.
func (g *Gudang) Sisip(ctx context.Context, tx *db.Tx, id, sekarang, pelaku string, b models.Baris) error {
	t, err := g.nama(TabelAggregate)
	if err != nil {
		return err
	}
	nilai, err := NilaiSisip(id, sekarang, pelaku, b)
	if err != nil {
		return err
	}
	_, err = jalankan(ctx, g.dari(tx), sqlSisip(t), "menyisipkan aggregate", nilai...)
	return err
}
