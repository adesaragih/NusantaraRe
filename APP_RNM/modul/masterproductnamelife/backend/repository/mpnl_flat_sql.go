package repository

// Penulis dan pembaca tabel flat (`mpnl_flat.go` = pemetaannya). Dipakai gudang (baca/tulis produk) dan alat pindah.
//
// ⛔ Angka ditulis TANPA bergantung NLS sesi: desimal dipecah menjadi koefisien bulat (teks) dan skala, dirakit Oracle
// `TO_NUMBER(:koef) / POWER(10, :skala)` (pola `mastercontractretrolife`, ADR-U-0003 - nol float); dibaca lewat
// `db.FmtDesimal` lalu diseragamkan AngkaKanonik (`TM9` menulis `.5`, bukan `0.5`). Tanggal `TO_DATE`/`TO_CHAR`
// `YYYY-MM-DD`; stempel komentar `TO_TIMESTAMP`/`TO_CHAR` `…FF3` (GMT).
// ⛔ go-ora mengikat menurut URUTAN KEMUNCULAN penampung - nomornya dirakit berurutan, nilainya disusun urutan yang sama.
// ⛔ Seluruh anak satu produk ditulis ULANG di transaksi pemanggil (hapus lalu sisip); nol COMMIT di teks SQL.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
)

// ErrNilaiTidakMuat - nilai produk tidak muat kolom flatnya (services menolaknya lebih dulu, berkalimat).
var ErrNilaiTidakMuat = errors.New("repository: a product value does not fit its column")

// bentukStempelOracle - stempel komentar sebagai teks bind `TO_TIMESTAMP(…, 'YYYYMMDDHH24MISS.FF3')`.
const bentukStempelOracle = "20060102150405.000"

// penampung merakit `:1`, `:2`, … menurut urutan kemunculan.
type penampung struct{ n int }

func (p *penampung) berikut() string {
	p.n++
	return ":" + strconv.Itoa(p.n)
}

// ekspresiTulis - ekspresi SQL satu nilai kolom (satu atau dua penampung).
func ekspresiTulis(jenis JenisKolomFlat, ph *penampung) string {
	switch jenis {
	case FlatDesimal, FlatBulat:
		koef := ph.berikut()
		return fmt.Sprintf("(TO_NUMBER(%s) / POWER(10, %s))", koef, ph.berikut())
	case FlatTanggal:
		return fmt.Sprintf("TO_DATE(%s, 'YYYY-MM-DD')", ph.berikut())
	case FlatStempel:
		return fmt.Sprintf("TO_TIMESTAMP(%s, 'YYYYMMDDHH24MISS.FF3')", ph.berikut())
	default:
		return ph.berikut()
	}
}

// pecahAngka - desimal → koefisien bulat (teks) dan skala: 1500000000.10 → ("150000000010", 2).
func pecahAngka(d *apd.Decimal) (string, int64) {
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

// argTulis - nilai bind satu kolom; `v` SUDAH kanonik (NormalkanFlat). Kosong = NULL.
func argTulis(jenis JenisKolomFlat, v string) []any {
	switch jenis {
	case FlatDesimal, FlatBulat:
		if v == "" {
			return []any{nil, int64(0)}
		}
		_, d, _ := AngkaKanonik(v)
		koef, skala := pecahAngka(d)
		return []any{koef, skala}
	case FlatStempel:
		if v == "" {
			return []any{nil}
		}
		t, _ := time.Parse(BentukStempelPega, v)
		return []any{t.UTC().Format(bentukStempelOracle)}
	default:
		return []any{db.KosongJadiNil(v)}
	}
}

// ekspresiBaca - kolom sebagai teks.
func ekspresiBaca(nama string, jenis JenisKolomFlat) string {
	switch jenis {
	case FlatDesimal, FlatBulat:
		return fmt.Sprintf(db.FmtDesimal, nama)
	case FlatTanggal:
		return fmt.Sprintf("TO_CHAR(%s, 'YYYY-MM-DD')", nama)
	case FlatStempel:
		return fmt.Sprintf(`TO_CHAR(%s, 'YYYYMMDD"T"HH24MISS.FF3')`, nama)
	default:
		return nama
	}
}

// nilaiBaca - teks Oracle → teks model kanonik.
func nilaiBaca(jenis JenisKolomFlat, v sql.NullString) string {
	if !v.Valid {
		return ""
	}
	switch jenis {
	case FlatDesimal, FlatBulat:
		if k, _, ok := AngkaKanonik(v.String); ok {
			return k
		}
		return v.String
	case FlatStempel:
		return v.String + " GMT"
	default:
		return v.String
	}
}

func bendera01(v bool) int {
	if v {
		return 1
	}
	return 0
}

// --- induk --------------------------------------------------------------------------

func sqlSisipInduk(tabel string) string {
	ph := &penampung{}
	kolom, nilai := []string{KolomIDFlat}, []string{ph.berikut()}
	for _, k := range KolomFlatInduk {
		kolom = append(kolom, k.Nama)
		nilai = append(nilai, ekspresiTulis(k.Jenis, ph))
	}
	kolom, nilai = append(kolom, KolomIsORS), append(nilai, ph.berikut())
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, tabel, strings.Join(kolom, ", "), strings.Join(nilai, ", "))
}

func sqlPerbaruiInduk(tabel string) string {
	ph := &penampung{}
	set := make([]string, 0, len(KolomFlatInduk)+1)
	for _, k := range KolomFlatInduk {
		set = append(set, k.Nama+" = "+ekspresiTulis(k.Jenis, ph))
	}
	set = append(set, KolomIsORS+" = "+ph.berikut())
	return fmt.Sprintf(`UPDATE %s SET %s WHERE %s = %s`, tabel, strings.Join(set, ", "), KolomIDFlat, ph.berikut())
}

// argInduk - nilai bind kolom induk (tanpa ID), urutan KolomFlatInduk lalu IS_ORS; `n` sudah kanonik.
func argInduk(n models.Produk) []any {
	var args []any
	for _, k := range KolomFlatInduk {
		args = append(args, argTulis(k.Jenis, *k.ambil(&n))...)
	}
	return append(args, bendera01(n.Umum.IsORS))
}

func sqlBacaInduk(tabel string, kunci bool) string {
	kolom := []string{KolomIDFlat}
	for _, k := range KolomFlatInduk {
		kolom = append(kolom, ekspresiBaca(k.Nama, k.Jenis))
	}
	kolom = append(kolom, "TO_CHAR("+KolomIsORS+")")
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = :1`, strings.Join(kolom, ", "), tabel, KolomIDFlat)
	if kunci {
		q += ` FOR UPDATE`
	}
	return q
}

// sqlDaftarFlat - grid daftar (`BrowseProduct_Life`: tanpa saringan b682, urut `.ID ASC` b1096, `pyMaxRecords`
// 500 b1078) - kelima kolom grid langsung dari induk.
func sqlDaftarFlat(tabel string) string {
	return fmt.Sprintf(`SELECT ID, CEDING, TREATYNUMBER, INWARDNAME, CREATEOP, UPDATEOP FROM %s
		ORDER BY ID ASC FETCH FIRST %d ROWS ONLY`, tabel, MaksBarisGrid)
}

func sqlSemuaIDFlat(tabel string) string {
	return fmt.Sprintf(`SELECT ID FROM %s ORDER BY ID ASC`, tabel)
}

// --- anak ---------------------------------------------------------------------------

func sqlHapusAnak(tabel string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE %s = :1`, tabel, KolomProductID)
}

func sqlSisipAnak[T any](a AnakFlat[T], tabel string) string {
	ph := &penampung{}
	kolom, nilai := []string{KolomProductID, KolomUrut}, []string{ph.berikut(), ph.berikut()}
	for _, k := range a.Kolom {
		kolom = append(kolom, k.Nama)
		nilai = append(nilai, ekspresiTulis(k.Jenis, ph))
	}
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, tabel, strings.Join(kolom, ", "), strings.Join(nilai, ", "))
}

func sqlBacaAnak[T any](a AnakFlat[T], tabel string) string {
	kolom := []string{KolomUrut}
	for _, k := range a.Kolom {
		kolom = append(kolom, ekspresiBaca(k.Nama, k.Jenis))
	}
	return fmt.Sprintf(`SELECT %s FROM %s WHERE %s = :1 ORDER BY %s`, strings.Join(kolom, ", "), tabel,
		KolomProductID, KolomUrut)
}

// eksekusi - satu pernyataan tulis di transaksi pemanggil; satuBaris = tuntut tepat satu baris.
func (g *Gudang) eksekusi(ctx context.Context, tx *db.Tx, objek string, susun func(string) string, satuBaris bool,
	args ...any) error {
	if !tx.Terisi() {
		return errors.New("repository: writing a product requires a transaction")
	}
	q, err := g.siapkan(objek, susun)
	if err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("repository: writing %s: %w", objek, err)
	}
	if satuBaris {
		return db.PastikanSatuBaris(hasil, "writing "+objek)
	}
	return nil
}

func tulisAnak[T any](ctx context.Context, g *Gudang, tx *db.Tx, a AnakFlat[T], n *models.Produk) error {
	if err := g.eksekusi(ctx, tx, a.Tabel, sqlHapusAnak, false, n.ID); err != nil {
		return err
	}
	baris := *a.daftar(n)
	for i := range baris {
		args := []any{n.ID, i + 1}
		for _, k := range a.Kolom {
			args = append(args, argTulis(k.Jenis, *k.ambil(&baris[i]))...)
		}
		if err := g.eksekusi(ctx, tx, a.Tabel, func(t string) string { return sqlSisipAnak(a, t) }, true,
			args...); err != nil {
			return err
		}
	}
	return nil
}

func bacaAnak[T any](ctx context.Context, g *Gudang, tx *db.Tx, a AnakFlat[T], n *models.Produk) error {
	q, err := g.siapkan(a.Tabel, func(t string) string { return sqlBacaAnak(a, t) })
	if err != nil {
		return err
	}
	rows, err := g.kueri(tx).QueryContext(ctx, q, n.ID)
	if err != nil {
		return fmt.Errorf("repository: reading %s %s: %w", a.Tabel, n.ID, err)
	}
	defer func() { _ = rows.Close() }()
	hasil := []T{}
	for rows.Next() {
		sel := make([]sql.NullString, 1+len(a.Kolom))
		tujuan := make([]any, len(sel))
		for i := range sel {
			tujuan[i] = &sel[i]
		}
		if err := rows.Scan(tujuan...); err != nil {
			return fmt.Errorf("repository: reading %s %s: %w", a.Tabel, n.ID, err)
		}
		var b T
		for i, k := range a.Kolom {
			*k.ambil(&b) = nilaiBaca(k.Jenis, sel[i+1])
		}
		hasil = append(hasil, b)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("repository: reading %s %s: %w", a.Tabel, n.ID, err)
	}
	*a.daftar(n) = hasil
	return nil
}

// --- produk utuh -----------------------------------------------------------------------

// tulisFlat menulis produk (induk sisip/perbarui + seluruh anak ditulis ulang) di transaksi pemanggil. Nilainya
// diseragamkan NormalkanFlat lebih dulu; nilai yang tidak muat = ErrNilaiTidakMuat (nol tulisan).
func (g *Gudang) tulisFlat(ctx context.Context, tx *db.Tx, p models.Produk, baru bool) error {
	n, masalah := NormalkanFlat(p)
	if len(masalah) > 0 {
		m := masalah[0]
		return fmt.Errorf("%w: %s.%s (%s) and %d more", ErrNilaiTidakMuat, m.Tabel, m.Kolom, m.Jenis, len(masalah)-1)
	}
	if baru {
		args := append([]any{n.ID}, argInduk(n)...)
		if err := g.eksekusi(ctx, tx, TabelFlatInduk, sqlSisipInduk, true, args...); err != nil {
			return err
		}
	} else if err := g.eksekusi(ctx, tx, TabelFlatInduk, sqlPerbaruiInduk, true,
		append(argInduk(n), n.ID)...); err != nil {
		return err
	}
	for _, tulis := range []func() error{
		func() error { return tulisAnak(ctx, g, tx, AnakLien, &n) },
		func() error { return tulisAnak(ctx, g, tx, AnakDokumen, &n) },
		func() error { return tulisAnak(ctx, g, tx, AnakPlan, &n) },
		func() error { return tulisAnak(ctx, g, tx, AnakFinUW, &n) },
		func() error { return tulisAnak(ctx, g, tx, AnakUWLimit, &n) },
		func() error { return tulisAnak(ctx, g, tx, AnakOutward, &n) },
		func() error { return tulisAnak(ctx, g, tx, AnakKomentar, &n) },
	} {
		if err := tulis(); err != nil {
			return err
		}
	}
	return nil
}

// bacaFlat membaca satu produk utuh; `kunci` = `FOR UPDATE` atas baris induk (di dalam simpan).
func (g *Gudang) bacaFlat(ctx context.Context, tx *db.Tx, id string, kunci bool) (models.Produk, error) {
	q, err := g.siapkan(TabelFlatInduk, func(t string) string { return sqlBacaInduk(t, kunci) })
	if err != nil {
		return models.Produk{}, err
	}
	rows, err := g.kueri(tx).QueryContext(ctx, q, id)
	if err != nil {
		return models.Produk{}, fmt.Errorf("repository: reading %s %s: %w", TabelFlatInduk, id, err)
	}
	var semua [][]sql.NullString
	for rows.Next() {
		sel := make([]sql.NullString, 2+len(KolomFlatInduk))
		tujuan := make([]any, len(sel))
		for i := range sel {
			tujuan[i] = &sel[i]
		}
		if err := rows.Scan(tujuan...); err != nil {
			_ = rows.Close()
			return models.Produk{}, fmt.Errorf("repository: reading %s %s: %w", TabelFlatInduk, id, err)
		}
		semua = append(semua, sel)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return models.Produk{}, fmt.Errorf("repository: reading %s %s: %w", TabelFlatInduk, id, err)
	}
	_ = rows.Close()
	switch len(semua) {
	case 0:
		return models.Produk{}, fmt.Errorf("%w: %s", ErrTidakAda, id)
	case 1:
	default:
		return models.Produk{}, fmt.Errorf("%w: %s %s (%d rows)", ErrIdentitasGanda, TabelFlatInduk, id, len(semua))
	}
	sel := semua[0]
	p := models.Produk{ID: sel[0].String}
	for i, k := range KolomFlatInduk {
		*k.ambil(&p) = nilaiBaca(k.Jenis, sel[i+1])
	}
	p.Umum.IsORS = sel[len(sel)-1].String == "1"
	p.Umum.PolicyHolder, p.Umum.PolicyHolderName = p.Inward.PolicyHolder, p.Inward.PolicyHolderName
	p.Inward.ID, p.Inward.ProductID = p.ID, p.ID
	for _, baca := range []func() error{
		func() error { return bacaAnak(ctx, g, tx, AnakLien, &p) },
		func() error { return bacaAnak(ctx, g, tx, AnakDokumen, &p) },
		func() error { return bacaAnak(ctx, g, tx, AnakPlan, &p) },
		func() error { return bacaAnak(ctx, g, tx, AnakFinUW, &p) },
		func() error { return bacaAnak(ctx, g, tx, AnakUWLimit, &p) },
		func() error { return bacaAnak(ctx, g, tx, AnakOutward, &p) },
		func() error { return bacaAnak(ctx, g, tx, AnakKomentar, &p) },
	} {
		if err := baca(); err != nil {
			return models.Produk{}, err
		}
	}
	return p, nil
}

// semuaIDFlat - seluruh ID induk, urut.
func (g *Gudang) semuaIDFlat(ctx context.Context, tx *db.Tx) ([]string, error) {
	q, err := g.siapkan(TabelFlatInduk, sqlSemuaIDFlat)
	if err != nil {
		return nil, err
	}
	rows, err := g.kueri(tx).QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: reading %s: %w", TabelFlatInduk, err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []string
	for rows.Next() {
		var id sql.NullString
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("repository: reading %s: %w", TabelFlatInduk, err)
		}
		hasil = append(hasil, id.String)
	}
	return hasil, rows.Err()
}
