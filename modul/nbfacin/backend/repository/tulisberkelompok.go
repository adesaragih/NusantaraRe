package repository

// Penulisan objek berkelompok (tiket 66 K66-B13, keputusan work owner 08-10-2026 "Percepat penyimpanan"): baris objek yang SAMA
// dengan jalur lama (satu NEXTVAL + satu INSERT per baris) ditulis dengan round-trip jauh lebih sedikit.
//
// Cara aman (anak butuh ID induk SEBELUM induk ditulis, sehingga INSERT … RETURNING ber-array tidak dipakai):
//  1. lintasan hitung - traversal yang sama (tulisObjek.sisipObjek) dengan penghitungID: berapa ID per sequence;
//  2. ID diambil berkelompok dari sequence-nya (BUKAN ditebak): sequence ber-jumlah sama dalam SATU kueri
//     `SELECT s1.NEXTVAL, s2.NEXTVAL, … FROM DUAL CONNECT BY LEVEL <= :1`;
//  3. lintasan tulis - traversal yang sama dengan penulisKelompok: ID diambil berurutan dari kumpulan, baris disangga per teks
//     SQL (urutan kemunculan pertama = induk sebelum anak), uid per baris seperti jalur lama;
//  4. tiap kelompok dijalankan array-bind go-ora (setiap argumen berupa slice satu nilai per baris), potongan ukuranKelompok
//     baris; jumlah baris terpengaruh wajib = jumlah baris (pengganti db.PastikanSatuBaris per baris).
// Jalur lama (satu NEXTVAL + satu INSERT per baris) tidak dipakai lagi; padanannya di uji (rekamLangsung).
// Kesetaraan dengan jalur lama diuji tanpa basis data (tulisberkelompok_test.go): baris per tabel (urutan, SEQ_NO, PARENT_ID,
// nilai kolom) identik.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// ukuranKelompok - baris per eksekusi array-bind.
const ukuranKelompok = 500

// penulisBaris - tujuan traversal tulisObjek: nomor (NEXTVAL) dan sisip (satu baris INSERT).
type penulisBaris interface {
	nomor(ctx context.Context, sequence string) (string, error)
	sisip(ctx context.Context, q, tabel string, arg []any) error
}

// tulisObjek - traversal baris objek (objek.go / item.go / coverage.go / deductible.go / okupasi.go / fea.go / kerugian.go).
type tulisObjek struct {
	t   tabelObjek
	p   penulisBaris
	uid func() (string, error)
}

// penghitungID - lintasan hitung: jumlah ID per sequence (urutan kemunculan pertama dicatat); tidak menulis apa pun.
type penghitungID struct {
	n    map[string]int
	urut []string
}

func (p *penghitungID) nomor(_ context.Context, sequence string) (string, error) {
	if p.n[sequence] == 0 {
		p.urut = append(p.urut, sequence)
	}
	p.n[sequence]++
	return "0", nil
}

func (p *penghitungID) sisip(context.Context, string, string, []any) error { return nil }

// kelompokSisip - baris satu teks SQL, urutan traversal.
type kelompokSisip struct {
	q, tabel string
	baris    [][]any
}

// penulisKelompok - lintasan tulis: ID dari kumpulan (berurutan), baris disangga per SQL.
type penulisKelompok struct {
	id       map[string][]string
	kelompok []*kelompokSisip
	indeks   map[string]*kelompokSisip
}

func (p *penulisKelompok) nomor(_ context.Context, sequence string) (string, error) {
	sisa := p.id[sequence]
	if len(sisa) == 0 {
		return "", fmt.Errorf("repository: kumpulan ID %s habis (lintasan hitung dan tulis berbeda)", sequence)
	}
	p.id[sequence] = sisa[1:]
	return sisa[0], nil
}

func (p *penulisKelompok) sisip(_ context.Context, q, tabel string, arg []any) error {
	k, ada := p.indeks[q]
	if !ada {
		k = &kelompokSisip{q: q, tabel: tabel}
		p.indeks[q] = k
		p.kelompok = append(p.kelompok, k)
	}
	k.baris = append(k.baris, append([]any(nil), arg...))
	return nil
}

// eksekutorKelompok - dua operasi basis data jalur berkelompok (Oracle: oracleKelompok; uji: tiruan penghitung).
type eksekutorKelompok interface {
	// nomorBerkelompok - n ID tiap sequence (satu kueri untuk semua sequence yang diminta).
	nomorBerkelompok(ctx context.Context, sequence []string, n int) ([][]string, error)
	// sisipArray - satu INSERT array-bind: kolom[i] = nilai kolom i untuk n baris; wajib n baris terpengaruh.
	sisipArray(ctx context.Context, q, tabel string, kolom []any, n int) error
}

// kolomArray - baris (n × c) -> c kolom, masing-masing []any n nilai (bentuk argumen array-bind go-ora).
func kolomArray(baris [][]any) []any {
	if len(baris) == 0 {
		return nil
	}
	kolom := make([]any, len(baris[0]))
	for c := range kolom {
		v := make([]any, len(baris))
		for r := range baris {
			v[r] = baris[r][c]
		}
		kolom[c] = v
	}
	return kolom
}

// sisipBerkelompok - lihat kepala berkas.
func sisipBerkelompok(ctx context.Context, e eksekutorKelompok, t tabelObjek, id string, baris []models.ObjekFire,
	uid func() (string, error)) error {
	hitung := &penghitungID{n: map[string]int{}}
	nolUID := func() (string, error) { return "", nil }
	if err := (&tulisObjek{t: t, p: hitung, uid: nolUID}).sisipObjek(ctx, id, baris); err != nil {
		return err
	}
	// Sequence ber-jumlah sama -> satu kueri (urutan kemunculan pertama dipertahankan).
	kumpulan := map[string][]string{}
	var jumlah []int
	perJumlah := map[int][]string{}
	for _, s := range hitung.urut {
		n := hitung.n[s]
		if _, ada := perJumlah[n]; !ada {
			jumlah = append(jumlah, n)
		}
		perJumlah[n] = append(perJumlah[n], s)
	}
	for _, n := range jumlah {
		ids, err := e.nomorBerkelompok(ctx, perJumlah[n], n)
		if err != nil {
			return err
		}
		for i, s := range perJumlah[n] {
			kumpulan[s] = ids[i]
		}
	}
	kel := &penulisKelompok{id: kumpulan, indeks: map[string]*kelompokSisip{}}
	if err := (&tulisObjek{t: t, p: kel, uid: uid}).sisipObjek(ctx, id, baris); err != nil {
		return err
	}
	for s, sisa := range kel.id {
		if len(sisa) != 0 {
			return fmt.Errorf("repository: %d ID %s tidak terpakai (lintasan hitung dan tulis berbeda)", len(sisa), s)
		}
	}
	for _, k := range kel.kelompok {
		for awal := 0; awal < len(k.baris); awal += ukuranKelompok {
			akhir := min(awal+ukuranKelompok, len(k.baris))
			if err := e.sisipArray(ctx, k.q, k.tabel, kolomArray(k.baris[awal:akhir]), akhir-awal); err != nil {
				return err
			}
		}
	}
	return nil
}

// argEksekusi - argumen ExecContext: n > 1 kolom array apa adanya; n == 1 nilai skalar tiap kolom.
func argEksekusi(kolom []any, n int) []any {
	if n != 1 {
		return kolom
	}
	arg := make([]any, len(kolom))
	for c := range kolom {
		arg[c] = kolom[c].([]any)[0]
	}
	return arg
}

// oracleKelompok - eksekutorKelompok atas Oracle (go-ora array-bind).
type oracleKelompok struct {
	d  *db.DB
	tx *db.Tx
}

// sqlNomorBerkelompok - n baris, satu kolom NEXTVAL per sequence (sequence sudah berkualifikasi skema).
func sqlNomorBerkelompok(sequence []string) string {
	kol := make([]string, len(sequence))
	for i, s := range sequence {
		kol[i] = s + ".NEXTVAL"
	}
	return "SELECT " + strings.Join(kol, ", ") + " FROM DUAL CONNECT BY LEVEL <= :1"
}

func (o oracleKelompok) nomorBerkelompok(ctx context.Context, sequence []string, n int) ([][]string, error) {
	nama := make([]string, len(sequence))
	for i, s := range sequence {
		q, err := o.d.Qualify(s)
		if err != nil {
			return nil, err
		}
		nama[i] = q
	}
	q := sqlNomorBerkelompok(nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	baris, err := o.tx.QueryContext(ctx, q, n)
	if err != nil {
		return nil, fmt.Errorf("repository: mengambil nomor %v: %w", sequence, err)
	}
	defer baris.Close()
	hasil := make([][]string, len(sequence))
	nilai := make([]sql.NullInt64, len(sequence))
	tuju := make([]any, len(sequence))
	for i := range nilai {
		tuju[i] = &nilai[i]
	}
	for baris.Next() {
		if err := baris.Scan(tuju...); err != nil {
			return nil, fmt.Errorf("repository: mengambil nomor %v: %w", sequence, err)
		}
		for i := range nilai {
			hasil[i] = append(hasil[i], fmt.Sprint(nilai[i].Int64))
		}
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: mengambil nomor %v: %w", sequence, err)
	}
	if len(hasil[0]) != n {
		return nil, fmt.Errorf("repository: mengambil nomor %v: %d dari %d", sequence, len(hasil[0]), n)
	}
	return hasil, nil
}

// sisipArray - n > 1: array-bind; n == 1: bind skalar biasa (sama dengan jalur lama; tidak bergantung pada array-bind satu
// elemen go-ora).
func (o oracleKelompok) sisipArray(ctx context.Context, q, tabel string, kolom []any, n int) error {
	h, err := jalankan(ctx, o.tx, q, "menyisipkan "+tabel, argEksekusi(kolom, n)...)
	if err != nil {
		return err
	}
	m, err := h.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: %s: %w", tabel, err)
	}
	if m != int64(n) {
		return fmt.Errorf("repository: %s: %d baris tersisip, seharusnya %d", tabel, m, n)
	}
	return nil
}

// sisipObjekBerkelompok - GantiObjek: jalur berkelompok atas transaksi tx.
func (r *ObjekOracle) sisipObjekBerkelompok(ctx context.Context, tx *db.Tx, t tabelObjek, id string, baris []models.ObjekFire) error {
	return sisipBerkelompok(ctx, oracleKelompok{d: r.db, tx: tx}, t, id, baris, uidAcak)
}
