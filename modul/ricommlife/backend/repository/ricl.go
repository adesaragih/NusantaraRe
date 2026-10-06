package repository

// SQL modul R/I Comm Life. Ringkasan: baca dari view `RICOMM_LIFE_SUMMARY`, tulis tabel JSON `M_RICOMM_LIFE_SUMMARY`
// (sisip `JSON_OBJECT`, ubah ricl_json.go). Rincian: tabel FLAT `RICOMM_LIFE` (migrasi inti 924) - kolom bernama, bukan
// JSONDATA. Angka ditulis TANPA bergantung NLS sesi: bulat = `TO_NUMBER(:n)` atas teks angka saja, desimal =
// `TO_NUMBER(:koef) / POWER(10, :skala)` (pola masterproductnamelife, ADR-U-0003 - nol float); dibaca `fmtAngka`.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/ricommlife/backend/models"
)

var (
	// ErrTidakAda - ID tidak ada.
	ErrTidakAda = errors.New("repository: R/I Comm Life tidak ada")
	// ErrBelumAda - tabel, view, atau sequence tidak ada di skema ini (ORA-00942 / ORA-00904 / ORA-02289).
	ErrBelumAda = errors.New("repository: tabel, view, atau sequence R/I Comm Life tidak ada di skema ini")
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
type nama struct{ tabelRingkasan, tabelKomisi, viewRingkasan, tabelLama, tabelSitus, seqRingkasan, seqKomisi string }

func (g *Gudang) nama() (nama, error) {
	var n nama
	for _, p := range []struct {
		ke    *string
		objek string
	}{{&n.tabelRingkasan, TabelRingkasan}, {&n.tabelKomisi, TabelKomisi}, {&n.viewRingkasan, ViewRingkasan},
		{&n.tabelLama, TabelJSONLama}, {&n.tabelSitus, TabelSitus}, {&n.seqRingkasan, SeqRingkasan}, {&n.seqKomisi, SeqKomisi}} {
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

// SqlDaftar - satu halaman grid `BrowseRICommSummary`.
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

// SqlSitus - ID situs aktif (prosedur PEGA_M_RICOMM_LIFE b13), sebagai teks; semua baris dibaca supaya "tidak tepat
// satu" terdengar.
func SqlSitus(t string) string {
	return fmt.Sprintf(`SELECT TO_CHAR(ID) FROM %s WHERE CURRENT_SITE = :1`, t)
}

// SqlNomorBaru - satu nomor sequence sebagai teks.
func SqlNomorBaru(seq string) string { return fmt.Sprintf(`SELECT TO_CHAR(%s.NEXTVAL) FROM DUAL`, seq) }

// SqlAdaID - ID sudah terpakai?
func SqlAdaID(t string) string { return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ID = :1`, t) }

// SqlSisipRingkasan - ringkasan baru (`JSON_OBJECT`; `pxObjClass` seperti data DEV).
func SqlSisipRingkasan(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, %s) VALUES (:1, JSON_OBJECT('%s' VALUE :2, '%s' VALUE :3, '%s' VALUE :4, '%s' VALUE :5 ABSENT ON NULL))`,
		t, KolomJSON, JSONUsedBy, JSONOperatorID, JSONModified, JSONKelas)
}

// SqlHapusRingkasan - Delete ringkasan.
func SqlHapusRingkasan(t string) string { return fmt.Sprintf(`DELETE FROM %s WHERE ID = :1`, t) }

// kolomKomisi - kolom tabel flat sebagai teks.
func kolomKomisi() string {
	return fmt.Sprintf("ID, IDUSEDBY, USEDBY, %s, %s, %s", fmt.Sprintf(fmtAngka, "CONTRACT"), fmt.Sprintf(fmtAngka, "YEAR"),
		fmt.Sprintf(fmtAngka, "COMM"))
}

// SqlDaftarKomisi - satu halaman R/I COMM DETAIL (RD `BrowseRICommLife_RD`, param idusedby b7396; sort ID ASC b9515).
func SqlDaftarKomisi(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE IDUSEDBY = :1 ORDER BY %s OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY`,
		kolomKomisi(), t, urutID("ASC"))
}

// SqlJumlahKomisi - baris rincian milik ringkasan.
func SqlJumlahKomisi(t string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE IDUSEDBY = :1`, t)
}

// SqlKomisiDari - seluruh rincian untuk n ringkasan sekaligus (kunci kembar); n >= 1.
func SqlKomisiDari(t string, n int) string {
	ikat := make([]string, n)
	for i := range ikat {
		ikat[i] = fmt.Sprintf(":%d", i+1)
	}
	return fmt.Sprintf(`SELECT %s FROM %s WHERE IDUSEDBY IN (%s) ORDER BY %s`, kolomKomisi(), t, strings.Join(ikat, ", "), urutID("ASC"))
}

// SqlSemuaKomisi - seluruh tabel flat (alat pindah).
func SqlSemuaKomisi(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s ORDER BY %s`, kolomKomisi(), t, urutID("ASC"))
}

// SqlSisipKomisi - satu baris flat; COMM dirakit dari koefisien dan skala.
func SqlSisipKomisi(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, IDUSEDBY, USEDBY, CONTRACT, YEAR, COMM) VALUES (:1, :2, :3, TO_NUMBER(:4), TO_NUMBER(:5), TO_NUMBER(:6) / POWER(10, :7))`, t)
}

// SqlUbahKomisi - Edit (`EditList_DT` b9323): CONTRACT, YEAR, COMM baris milik ringkasan itu.
func SqlUbahKomisi(t string) string {
	return fmt.Sprintf(`UPDATE %s SET CONTRACT = TO_NUMBER(:1), YEAR = TO_NUMBER(:2), COMM = TO_NUMBER(:3) / POWER(10, :4) WHERE ID = :5 AND IDUSEDBY = :6`, t)
}

// SqlUbahNamaKomisi - salinan nama di setiap rincian ringkasan (Edit nama, ASUMSI seperti riratelife A6).
func SqlUbahNamaKomisi(t string) string {
	return fmt.Sprintf(`UPDATE %s SET USEDBY = :1 WHERE IDUSEDBY = :2`, t)
}

// SqlHapusKomisi - rincian milik ringkasan (`DeleteSummaryDetail`: ringkasan BESERTA rinciannya).
func SqlHapusKomisi(t string) string { return fmt.Sprintf(`DELETE FROM %s WHERE IDUSEDBY = :1`, t) }

// SqlKunciKomisi - kunci tabel flat sepanjang transaksi alat pindah.
func SqlKunciKomisi(t string) string { return fmt.Sprintf(`LOCK TABLE %s IN EXCLUSIVE MODE`, t) }

// SqlSumberJSON - seluruh baris JSON lama (alat pindah).
func SqlSumberJSON(t string) string {
	return fmt.Sprintf(`SELECT ID, %s FROM %s ORDER BY ID`, KolomJSON, t)
}

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

// angkaBaca - teks `TM9` Oracle -> bentuk kanonik models (`.5` -> `0.5`); bentuk lain apa adanya.
func angkaBaca(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, ".") {
		s = "0" + s
	}
	if k, _, _, ok := models.DesimalKanonik(s); ok {
		return k
	}
	return s
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

func keKomisi(b [][]string) []models.Komisi {
	out := make([]models.Komisi, 0, len(b))
	for _, s := range b {
		out = append(out, models.Komisi{ID: s[0], IDUsedBy: s[1], UsedBy: s[2], Contract: angkaBaca(s[3]),
			Year: angkaBaca(s[4]), Comm: angkaBaca(s[5])})
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

// NomorBaru - satu nomor sequence warisan (`ringkasan` true = M_RICOMM_LIFE_SUMMARY_SEQ, selain itu M_RICOMM_LIFE_SEQ).
func (g *Gudang) NomorBaru(ctx context.Context, tx *db.Tx, ringkasan bool) (string, error) {
	n, err := g.nama()
	if err != nil {
		return "", err
	}
	seq := n.seqKomisi
	if ringkasan {
		seq = n.seqRingkasan
	}
	return g.satuNilai(ctx, tx, "DUAL", SqlNomorBaru(seq))
}

// AdaID - ID sudah terpakai di tabel ringkasan / flat?
func (g *Gudang) AdaID(ctx context.Context, tx *db.Tx, ringkasan bool, id string) (bool, error) {
	n, err := g.nama()
	if err != nil {
		return false, err
	}
	objek, t := TabelKomisi, n.tabelKomisi
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
		k(r.UsedBy), k(r.OperatorID), k(r.ModifiedDate), models.KelasRingkasan)
	return err
}

// UbahRingkasan - hanya USEDBY, OPERATORID, MODIFIEDDATE yang diganti; kunci lain milik Pega tetap.
func (g *Gudang) UbahRingkasan(ctx context.Context, tx *db.Tx, r models.Ringkasan) error {
	n, err := g.nama()
	if err != nil {
		return err
	}
	return g.ubahJSON(ctx, tx, n.tabelRingkasan, r.ID, map[string]string{
		JSONUsedBy: r.UsedBy, JSONOperatorID: r.OperatorID, JSONModified: r.ModifiedDate})
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

// UbahNamaKomisi - salinan nama di setiap rincian ringkasan idUsedBy; jumlah baris.
func (g *Gudang) UbahNamaKomisi(ctx context.Context, tx *db.Tx, idUsedBy, nama string) (int, error) {
	n, err := g.nama()
	if err != nil {
		return 0, err
	}
	j, err := g.tulis(ctx, tx, TabelKomisi, SqlUbahNamaKomisi(n.tabelKomisi), "mengganti nama rincian", db.KosongJadiNil(nama), idUsedBy)
	return int(j), err
}

// HapusKomisi - rincian milik ringkasan; jumlah terhapus.
func (g *Gudang) HapusKomisi(ctx context.Context, tx *db.Tx, idUsedBy string) (int, error) {
	n, err := g.nama()
	if err != nil {
		return 0, err
	}
	j, err := g.tulis(ctx, tx, TabelKomisi, SqlHapusKomisi(n.tabelKomisi), "menghapus rincian", idUsedBy)
	return int(j), err
}

// JumlahKomisi - rincian milik ringkasan.
func (g *Gudang) JumlahKomisi(ctx context.Context, tx *db.Tx, idUsedBy string) (int, error) {
	n, err := g.nama()
	if err != nil {
		return 0, err
	}
	s, err := g.satuNilai(ctx, tx, TabelKomisi, SqlJumlahKomisi(n.tabelKomisi), idUsedBy)
	return angka(s), err
}

// DaftarKomisi - satu halaman R/I COMM DETAIL dan totalnya.
func (g *Gudang) DaftarKomisi(ctx context.Context, idUsedBy string, halaman int) ([]models.Komisi, int, error) {
	n, err := g.nama()
	if err != nil {
		return nil, 0, err
	}
	total, err := g.JumlahKomisi(ctx, nil, idUsedBy)
	if err != nil {
		return nil, 0, err
	}
	b, err := g.bacaBaris(ctx, nil, TabelKomisi, SqlDaftarKomisi(n.tabelKomisi), 6, idUsedBy,
		(halaman-1)*models.UkuranHalaman, models.UkuranHalaman)
	return keKomisi(b), total, err
}

// KomisiDari - seluruh rincian untuk ringkasan ids (kunci kembar). ids kosong = nil.
func (g *Gudang) KomisiDari(ctx context.Context, tx *db.Tx, ids []string) ([]models.Komisi, error) {
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
	b, err := g.bacaBaris(ctx, tx, TabelKomisi, SqlKomisiDari(n.tabelKomisi, len(ids)), 6, args...)
	return keKomisi(b), err
}

// argKomisi - CONTRACT, YEAR (teks angka bulat atau NULL), COMM (koefisien, skala).
func argKomisi(k models.Komisi) []any {
	koef, skala := PecahDesimal(k.Comm)
	return []any{db.KosongJadiNil(k.Contract), db.KosongJadiNil(k.Year), koef, skala}
}

// SisipKomisi - baris flat baru; ErrKembar bila ID terpakai.
func (g *Gudang) SisipKomisi(ctx context.Context, tx *db.Tx, k models.Komisi) error {
	n, err := g.nama()
	if err != nil {
		return err
	}
	args := append([]any{k.ID, db.KosongJadiNil(k.IDUsedBy), db.KosongJadiNil(k.UsedBy)}, argKomisi(k)...)
	_, err = g.tulis(ctx, tx, TabelKomisi, SqlSisipKomisi(n.tabelKomisi), "menyimpan rincian", args...)
	return err
}

// UbahKomisi - Edit satu rincian milik k.IDUsedBy; ErrTidakAda bila baris tidak ada di ringkasan itu.
func (g *Gudang) UbahKomisi(ctx context.Context, tx *db.Tx, k models.Komisi) error {
	n, err := g.nama()
	if err != nil {
		return err
	}
	args := append(argKomisi(k), k.ID, k.IDUsedBy)
	j, err := g.tulis(ctx, tx, TabelKomisi, SqlUbahKomisi(n.tabelKomisi), "mengubah rincian", args...)
	if err == nil && j == 0 {
		return ErrTidakAda
	}
	return err
}
