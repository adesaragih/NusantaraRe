package repository

// SQL modul R/I Risk. Ringkasan dibaca DAN ditulis di kolom `RIRISK_LIFE_SUMMARY`, rincian di kolom `RIRISK_LIFE`
// (keputusan work owner 08-10-2026 K1/K2, migrasi inti 935-940) - kolom bernama, nol JSON. CONTRACT, YEAR, MONTH =
// kolom TEKS warisan (VARCHAR2(10)): ditulis teks angka kanonik apa adanya. RISK = NUMBER warisan tanpa skala: ditulis
// `TO_NUMBER(:koef) / POWER(10, :skala)` (teks angka murni - TANPA bergantung NLS sesi; pola ricommlife, ADR-U-0003 -
// nol float), dibaca `fmtAngka`. AGE (kolom view lama, tidak ada di XML) tidak dibaca dan tidak ditulis.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/ririsklife/backend/models"
)

var (
	// ErrTidakAda - ID tidak ada.
	ErrTidakAda = errors.New("repository: R/I Risk tidak ada")
	// ErrBelumAda - tabel, kolom, atau sequence tidak ada di skema ini (ORA-00942 / ORA-00904 / ORA-02289).
	ErrBelumAda = errors.New("repository: tabel atau sequence R/I Risk tidak ada di skema ini")
	// ErrKembar - ORA-00001.
	ErrKembar = errors.New("repository: ID sudah dipakai")
	// ErrBacaSaja - SQL tulis diarahkan ke objek di luar DaftarTabelDitulis.
	ErrBacaSaja = errors.New("repository: objek ini dibaca saja")
	// ErrSitus - `M_SITE_DATABASE` tidak punya tepat satu situs aktif (CURRENT_SITE = '1').
	ErrSitus = errors.New("repository: M_SITE_DATABASE harus punya tepat satu baris CURRENT_SITE = '1'")
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

func siap(objek, q string) error {
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	return PeriksaTulis(objek, q)
}

// nama - nama berskema setiap objek.
type nama struct{ tabelRingkasan, tabelRincian, tabelSitus, seqRingkasan, seqRincian string }

func (g *Gudang) nama() (nama, error) {
	var n nama
	for _, p := range []struct {
		ke    *string
		objek string
	}{{&n.tabelRingkasan, TabelRingkasan}, {&n.tabelRincian, TabelRincian}, {&n.tabelSitus, TabelSitus},
		{&n.seqRingkasan, SeqRingkasan}, {&n.seqRincian, SeqRincian}} {
		q, err := g.db.Qualify(p.objek)
		if err != nil {
			return nama{}, err
		}
		*p.ke = q
	}
	return n, nil
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

// urutID - ID angka (bukan urut teks), lalu ID.
func urutID(arah string) string {
	return fmt.Sprintf("TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) %s NULLS LAST, ID %s", arah, arah)
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
	return urutID(arah)
}

const saringRingkasan = `(:1 IS NULL OR UPPER(ID) LIKE :2 ESCAPE '\') AND (:3 IS NULL OR UPPER(USEDBY) LIKE :4 ESCAPE '\')`

// SqlDaftar - satu halaman grid `BrowseRIRiskSummary`.
func SqlDaftar(v, urut string, turun bool) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE %s ORDER BY %s OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY`,
		kolomRingkasan, v, saringRingkasan, Urutan(urut, turun))
}

// SqlJumlah - jumlah baris bersaring.
func SqlJumlah(v string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, v, saringRingkasan)
}

// SqlAmbil - satu ringkasan.
func SqlAmbil(v string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, kolomRingkasan, v)
}

// SqlPemakaiNama - ringkasan lain bernama sama (tanpa beda huruf dan spasi tepi), selain kecualiID.
func SqlPemakaiNama(v string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE UPPER(TRIM(USEDBY)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0)) ORDER BY ID`,
		kolomRingkasan, v)
}

// SqlSitus - ID situs aktif (prosedur PEGA_M_RIRISK_LIFE_SUMMARY), sebagai teks; semua baris dibaca supaya "tidak tepat
// satu" terdengar.
func SqlSitus(t string) string {
	return fmt.Sprintf(`SELECT TO_CHAR(ID) FROM %s WHERE CURRENT_SITE = :1`, t)
}

// SqlNomorBaru - satu nomor sequence sebagai teks.
func SqlNomorBaru(seq string) string { return fmt.Sprintf(`SELECT TO_CHAR(%s.NEXTVAL) FROM DUAL`, seq) }

// SqlAdaID - ID sudah terpakai?
func SqlAdaID(t string) string { return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ID = :1`, t) }

// SqlSisipRingkasan - ringkasan baru (kolom bernama; `pxObjClass` Pega tidak lagi disimpan - bukan kolom view).
func SqlSisipRingkasan(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, USEDBY, OPERATORID, MODIFIEDDATE) VALUES (:1, :2, :3, :4)`, t)
}

// SqlUbahRingkasan - Edit / Save: USEDBY, OPERATORID, MODIFIEDDATE satu ringkasan.
func SqlUbahRingkasan(t string) string {
	return fmt.Sprintf(`UPDATE %s SET USEDBY = :1, OPERATORID = :2, MODIFIEDDATE = :3 WHERE ID = :4`, t)
}

// SqlHapusRingkasan - Delete ringkasan.
func SqlHapusRingkasan(t string) string { return fmt.Sprintf(`DELETE FROM %s WHERE ID = :1`, t) }

// kolomRincian - kolom rincian sebagai teks (urutan = models.Rincian; RISK lewat fmtAngka).
func kolomRincian() string {
	return fmt.Sprintf("ID, IDUSEDBY, USEDBY, CONTRACT, YEAR, MONTH, %s", fmt.Sprintf(fmtAngka, "RISK"))
}

// SqlDaftarRincian - satu halaman R/I RISK DETAIL (RD `BrowseRIRiskLife_RD` b7612, param idusedby b7508; sort ID ASC
// b7607/b7619).
func SqlDaftarRincian(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE IDUSEDBY = :1 ORDER BY %s OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY`,
		kolomRincian(), t, urutID("ASC"))
}

// SqlJumlahRincian - baris rincian milik ringkasan.
func SqlJumlahRincian(t string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE IDUSEDBY = :1`, t)
}

// SqlRincianDari - seluruh rincian untuk n ringkasan sekaligus (kunci kembar); n >= 1.
func SqlRincianDari(t string, n int) string {
	ikat := make([]string, n)
	for i := range ikat {
		ikat[i] = fmt.Sprintf(":%d", i+1)
	}
	return fmt.Sprintf(`SELECT %s FROM %s WHERE IDUSEDBY IN (%s) ORDER BY %s`, kolomRincian(), t, strings.Join(ikat, ", "), urutID("ASC"))
}

// SqlSisipRincian - satu baris rincian; RISK dirakit dari koefisien dan skala.
func SqlSisipRincian(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, IDUSEDBY, USEDBY, CONTRACT, YEAR, MONTH, RISK) VALUES (:1, :2, :3, :4, :5, :6, TO_NUMBER(:7) / POWER(10, :8))`, t)
}

// SqlUbahRincian - EDIT (`EditRIRiskLife_Act` b9808) lalu Save: CONTRACT, YEAR, MONTH, RISK baris milik ringkasan itu.
func SqlUbahRincian(t string) string {
	return fmt.Sprintf(`UPDATE %s SET CONTRACT = :1, YEAR = :2, MONTH = :3, RISK = TO_NUMBER(:4) / POWER(10, :5) WHERE ID = :6 AND IDUSEDBY = :7`, t)
}

// SqlUbahNamaRincian - salinan nama di setiap rincian ringkasan (Edit nama ringkasan, pola ricommlife).
func SqlUbahNamaRincian(t string) string {
	return fmt.Sprintf(`UPDATE %s SET USEDBY = :1 WHERE IDUSEDBY = :2`, t)
}

// SqlHapusRincian - rincian milik ringkasan (`DeleteSummaryDetail`: ringkasan BESERTA rinciannya).
func SqlHapusRincian(t string) string { return fmt.Sprintf(`DELETE FROM %s WHERE IDUSEDBY = :1`, t) }

// PecahDesimal - desimal kanonik -> koefisien bulat (teks) dan skala: `12.05` -> ("1205", 2); kosong = (nil, 0).
func PecahDesimal(kanonik string) (any, int64) {
	if kanonik == "" {
		return nil, 0
	}
	bulat, pecahan, _ := strings.Cut(kanonik, ".")
	koef := strings.TrimLeft(bulat+pecahan, "0")
	if koef == "" {
		return "0", 0
	}
	return koef, int64(len(pecahan))
}

// AngkaOracle - teks angka dari Oracle -> bentuk kanonik models, diurai di Go TANPA bergantung NLS sesi: titik ATAU
// koma desimal (satu pemisah, mis. bila argumen NLS `TO_CHAR` tidak berlaku), pecahan tanpa nol depan (`TM9` menulis
// `.5` / `,5`) dan tanda minus. Bentuk lain (bukan angka) apa adanya.
func AngkaOracle(s string) string {
	s = strings.TrimSpace(s)
	minus := strings.HasPrefix(s, "-")
	t := strings.TrimPrefix(s, "-")
	if strings.Count(t, ",") == 1 && !strings.Contains(t, ".") {
		t = strings.Replace(t, ",", ".", 1)
	}
	if strings.HasPrefix(t, ".") {
		t = "0" + t
	}
	k, _, _, ok := models.DesimalKanonik(t)
	if !ok {
		return s
	}
	if minus && k != "0" {
		return "-" + k
	}
	return k
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

func keRincian(b [][]string) []models.Rincian {
	out := make([]models.Rincian, 0, len(b))
	for _, s := range b {
		out = append(out, models.Rincian{ID: s[0], IDUsedBy: s[1], UsedBy: s[2], Contract: s[3], Year: s[4], Month: s[5],
			Risk: AngkaOracle(s[6])})
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
	total, err := g.satuNilai(ctx, nil, TabelRingkasan, SqlJumlah(n.tabelRingkasan), id, id, nm, nm)
	if err != nil {
		return nil, 0, err
	}
	b, err := g.bacaBaris(ctx, nil, TabelRingkasan, SqlDaftar(n.tabelRingkasan, s.Urut, s.Turun), 4, id, id, nm, nm,
		(s.Halaman-1)*models.UkuranHalaman, models.UkuranHalaman)
	return keRingkasan(b), angka(total), err
}

// Ambil - satu ringkasan.
func (g *Gudang) Ambil(ctx context.Context, tx *db.Tx, id string) (models.Ringkasan, error) {
	n, err := g.nama()
	if err != nil {
		return models.Ringkasan{}, err
	}
	b, err := g.bacaBaris(ctx, tx, TabelRingkasan, SqlAmbil(n.tabelRingkasan), 4, id)
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
	b, err := g.bacaBaris(ctx, tx, TabelRingkasan, SqlPemakaiNama(n.tabelRingkasan), 4, nama, db.KosongJadiNil(kecualiID))
	return keRingkasan(b), err
}

// Situs - ID situs aktif `M_SITE_DATABASE` (CURRENT_SITE = '1'); ErrSitus bila tidak tepat satu baris.
func (g *Gudang) Situs(ctx context.Context, tx *db.Tx) (string, error) {
	n, err := g.nama()
	if err != nil {
		return "", err
	}
	b, err := g.bacaBaris(ctx, tx, TabelSitus, SqlSitus(n.tabelSitus), 1, situsAktif)
	if err != nil {
		return "", err
	}
	if len(b) != 1 {
		return "", fmt.Errorf("%w (%d baris)", ErrSitus, len(b))
	}
	return b[0][0], nil
}

// NomorBaru - satu nomor sequence warisan (`ringkasan` true = M_RIRISK_LIFE_SUMMARY_SEQ, selain itu M_RIRISK_LIFE_SEQ).
func (g *Gudang) NomorBaru(ctx context.Context, tx *db.Tx, ringkasan bool) (string, error) {
	n, err := g.nama()
	if err != nil {
		return "", err
	}
	seq := n.seqRincian
	if ringkasan {
		seq = n.seqRingkasan
	}
	return g.satuNilai(ctx, tx, "DUAL", SqlNomorBaru(seq))
}

// AdaID - ID sudah terpakai di tabel ringkasan / rincian (satu tabel per jenis)?
func (g *Gudang) AdaID(ctx context.Context, tx *db.Tx, ringkasan bool, id string) (bool, error) {
	n, err := g.nama()
	if err != nil {
		return false, err
	}
	objek, t := TabelRincian, n.tabelRincian
	if ringkasan {
		objek, t = TabelRingkasan, n.tabelRingkasan
	}
	s, err := g.satuNilai(ctx, tx, objek, SqlAdaID(t), id)
	return angka(s) > 0, err
}

// SisipRingkasan - ringkasan baru.
func (g *Gudang) SisipRingkasan(ctx context.Context, tx *db.Tx, r models.Ringkasan) error {
	n, err := g.nama()
	if err != nil {
		return err
	}
	k := db.KosongJadiNil
	_, err = g.tulis(ctx, tx, TabelRingkasan, SqlSisipRingkasan(n.tabelRingkasan), "menyimpan ringkasan", r.ID,
		k(r.UsedBy), k(r.OperatorID), k(r.ModifiedDate))
	return err
}

// UbahRingkasan - USEDBY, OPERATORID, MODIFIEDDATE (kosong = NULL); ErrTidakAda bila ID tidak ada.
func (g *Gudang) UbahRingkasan(ctx context.Context, tx *db.Tx, r models.Ringkasan) error {
	n, err := g.nama()
	if err != nil {
		return err
	}
	k := db.KosongJadiNil
	j, err := g.tulis(ctx, tx, TabelRingkasan, SqlUbahRingkasan(n.tabelRingkasan), "mengubah ringkasan",
		k(r.UsedBy), k(r.OperatorID), k(r.ModifiedDate), r.ID)
	if err == nil && j == 0 {
		return ErrTidakAda
	}
	return err
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

// UbahNamaRincian - salinan nama di setiap rincian ringkasan idUsedBy; jumlah baris.
func (g *Gudang) UbahNamaRincian(ctx context.Context, tx *db.Tx, idUsedBy, nama string) (int, error) {
	n, err := g.nama()
	if err != nil {
		return 0, err
	}
	j, err := g.tulis(ctx, tx, TabelRincian, SqlUbahNamaRincian(n.tabelRincian), "mengganti nama rincian", db.KosongJadiNil(nama), idUsedBy)
	return int(j), err
}

// HapusRincian - rincian milik ringkasan; jumlah terhapus.
func (g *Gudang) HapusRincian(ctx context.Context, tx *db.Tx, idUsedBy string) (int, error) {
	n, err := g.nama()
	if err != nil {
		return 0, err
	}
	j, err := g.tulis(ctx, tx, TabelRincian, SqlHapusRincian(n.tabelRincian), "menghapus rincian", idUsedBy)
	return int(j), err
}

// JumlahRincian - rincian milik ringkasan.
func (g *Gudang) JumlahRincian(ctx context.Context, tx *db.Tx, idUsedBy string) (int, error) {
	n, err := g.nama()
	if err != nil {
		return 0, err
	}
	s, err := g.satuNilai(ctx, tx, TabelRincian, SqlJumlahRincian(n.tabelRincian), idUsedBy)
	return angka(s), err
}

// DaftarRincian - satu halaman R/I RISK DETAIL dan totalnya.
func (g *Gudang) DaftarRincian(ctx context.Context, idUsedBy string, halaman int) ([]models.Rincian, int, error) {
	n, err := g.nama()
	if err != nil {
		return nil, 0, err
	}
	total, err := g.JumlahRincian(ctx, nil, idUsedBy)
	if err != nil {
		return nil, 0, err
	}
	b, err := g.bacaBaris(ctx, nil, TabelRincian, SqlDaftarRincian(n.tabelRincian), 7, idUsedBy,
		(halaman-1)*models.UkuranHalamanRincian, models.UkuranHalamanRincian)
	return keRincian(b), total, err
}

// RincianDari - seluruh rincian untuk ringkasan ids (kunci kembar). ids kosong = nil.
func (g *Gudang) RincianDari(ctx context.Context, tx *db.Tx, ids []string) ([]models.Rincian, error) {
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
	b, err := g.bacaBaris(ctx, tx, TabelRincian, SqlRincianDari(n.tabelRincian, len(ids)), 7, args...)
	return keRincian(b), err
}

// argRincian - CONTRACT, YEAR, MONTH (teks angka atau NULL), RISK (koefisien, skala).
func argRincian(k models.Rincian) []any {
	koef, skala := PecahDesimal(k.Risk)
	k0 := db.KosongJadiNil
	return []any{k0(k.Contract), k0(k.Year), k0(k.Month), koef, skala}
}

// SisipRincian - baris rincian baru; ErrKembar bila ID terpakai.
func (g *Gudang) SisipRincian(ctx context.Context, tx *db.Tx, k models.Rincian) error {
	n, err := g.nama()
	if err != nil {
		return err
	}
	args := append([]any{k.ID, db.KosongJadiNil(k.IDUsedBy), db.KosongJadiNil(k.UsedBy)}, argRincian(k)...)
	_, err = g.tulis(ctx, tx, TabelRincian, SqlSisipRincian(n.tabelRincian), "menyimpan rincian", args...)
	return err
}

// UbahRincian - Edit satu rincian milik k.IDUsedBy; ErrTidakAda bila baris tidak ada di ringkasan itu.
func (g *Gudang) UbahRincian(ctx context.Context, tx *db.Tx, k models.Rincian) error {
	n, err := g.nama()
	if err != nil {
		return err
	}
	args := append(argRincian(k), k.ID, k.IDUsedBy)
	j, err := g.tulis(ctx, tx, TabelRincian, SqlUbahRincian(n.tabelRincian), "mengubah rincian", args...)
	if err == nil && j == 0 {
		return ErrTidakAda
	}
	return err
}
