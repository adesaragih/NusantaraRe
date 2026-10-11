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
	"reflect"
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
	// awal - SEQ_NO lokasi = awal + urutan (tiket 75: kelompok berikutnya unggah CSV menambah sesudah lokasi yang sudah ada); 0 =
	// jalur lama.
	awal int
	// seq - bila tidak nil: SEQ_NO lokasi ke-i = seq[i], awal diabaikan (tiket 75 F3: lokasi diubah ditulis ulang di tempatnya).
	seq []int
}

// seqNo - SEQ_NO lokasi ke-i (0..).
func (w *tulisObjek) seqNo(i int) int {
	if w.seq != nil {
		return w.seq[i]
	}
	return w.awal + i + 1
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

// kolomArrayBertipe - kolomArray dengan SATU tipe Go per kolom (bug DEV 10-10-2026: Save Spreading Aneka ORA-01722 di
// T_SPREADINGLIST). Sebab `[terverifikasi: go-ora v2.8.19]`: INSERT array-bind (command.go Stmt._exec, cabang bulkExec)
// menyandi tiap elemen dengan SATU ParameterInfo yang dipakai ulang - encodeValue memanggil par.init() (parameter_encode.go
// L516: CharsetID / CharsetForm / ContFlag = 0) lalu setDataType menurut tipe Go elemen itu; nil -> NCHAR TANPA charset
// (L16-18), string -> NCHAR + charset server (L98-102). DataType kolom = elemen non-kosong terakhir, sedangkan CharsetID /
// CharsetForm yang dikirim di kepala bind (parameter.go write L373-374) = milik elemen TERAKHIR. Kolom teks yang baris
// terakhirnya NULL dan baris sebelumnya berisi terkirim sebagai NCHAR ber-charset 0; dampaknya di server (ORA-01722 di
// TO_NUMBER angkaMasuk) `[dugaan]` - tidak dapat dibuktikan dari kode klien.
// Perbaikan: kolom teks yang memuat NULL diikat sql.NullString di SETIAP baris (NULL = Valid false; setDataType cabang
// tyNullString = cabang string L98-102, nilai lewat driver.Valuer) - satu tipe Go per kolom, kepala bind sama di posisi mana
// pun; kolom tanpa NULL apa adanya. Kolom bercampur tipe, atau tipe bukan teks yang bercampur NULL, = galat (tidak ditebak).
// Dipakai jalur Aneka (sisipAnekaBerkelompok) dan jalur FIRE (sisipBerkelompok - pencegahan, keputusan work owner 10-10-2026,
// pengecualian sadar dari "jangan ubah yang sudah berhasil"; isi data yang disimpan tidak berubah).
func kolomArrayBertipe(baris [][]any) ([]any, error) {
	kolom := kolomArray(baris)
	for c, k := range kolom {
		v := k.([]any)
		var tipe reflect.Type
		adaNil := false
		for _, x := range v {
			if x == nil {
				adaNil = true
				continue
			}
			if t := reflect.TypeOf(x); tipe == nil {
				tipe = t
			} else if t != tipe {
				return nil, fmt.Errorf("repository: kolom %d array-bind bercampur tipe %v dan %v", c+1, tipe, t)
			}
		}
		if !adaNil {
			continue
		}
		if tipe != nil && tipe != tipeTeks {
			return nil, fmt.Errorf("repository: kolom %d array-bind bertipe %v bercampur NULL", c+1, tipe)
		}
		for r, x := range v {
			s, ok := x.(string)
			v[r] = sql.NullString{String: s, Valid: ok}
		}
	}
	return kolom, nil
}

var tipeTeks = reflect.TypeOf("")

// sisipBerkelompok - lihat kepala berkas.
func sisipBerkelompok(ctx context.Context, e eksekutorKelompok, t tabelObjek, id string, baris []models.ObjekFire,
	uid func() (string, error)) error {
	return sisipBerkelompokMulai(ctx, e, t, id, baris, uid, 0)
}

// sisipBerkelompokMulai - sisipBerkelompok dengan SEQ_NO lokasi mulai awal+1 (tiket 75).
func sisipBerkelompokMulai(ctx context.Context, e eksekutorKelompok, t tabelObjek, id string, baris []models.ObjekFire,
	uid func() (string, error), awal int) error {
	return sisipBerkelompokPosisi(ctx, e, t, id, baris, uid, awal, nil)
}

// sisipBerkelompokPosisi - sisipBerkelompok dengan SEQ_NO lokasi awal+1.. atau, bila seq tidak nil, seq[i] (tiket 75 F3).
func sisipBerkelompokPosisi(ctx context.Context, e eksekutorKelompok, t tabelObjek, id string, baris []models.ObjekFire,
	uid func() (string, error), awal int, seq []int) error {
	if seq != nil && len(seq) != len(baris) {
		return fmt.Errorf("repository: %d SEQ_NO untuk %d lokasi", len(seq), len(baris))
	}
	hitung := &penghitungID{n: map[string]int{}}
	nolUID := func() (string, error) { return "", nil }
	if err := (&tulisObjek{t: t, p: hitung, uid: nolUID}).sisipObjek(ctx, id, baris); err != nil {
		return err
	}
	kumpulan, err := ambilIDBerkelompok(ctx, e, hitung)
	if err != nil {
		return err
	}
	kel := &penulisKelompok{id: kumpulan, indeks: map[string]*kelompokSisip{}}
	if err := (&tulisObjek{t: t, p: kel, uid: uid, awal: awal, seq: seq}).sisipObjek(ctx, id, baris); err != nil {
		return err
	}
	return jalankanPenulisKelompok(ctx, e, kel)
}

// ambilIDBerkelompok - ID tiap sequence hasil lintasan hitung; sequence ber-jumlah sama -> satu kueri (urutan kemunculan pertama
// dipertahankan).
func ambilIDBerkelompok(ctx context.Context, e eksekutorKelompok, hitung *penghitungID) (map[string][]string, error) {
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
			return nil, err
		}
		for i, s := range perJumlah[n] {
			kumpulan[s] = ids[i]
		}
	}
	return kumpulan, nil
}

// jalankanPenulisKelompok - semua ID terpakai (lintasan hitung = tulis), lalu setiap kelompok dijalankan urut kemunculan (induk
// sebelum anak).
func jalankanPenulisKelompok(ctx context.Context, e eksekutorKelompok, kel *penulisKelompok) error {
	for s, sisa := range kel.id {
		if len(sisa) != 0 {
			return fmt.Errorf("repository: %d ID %s tidak terpakai (lintasan hitung dan tulis berbeda)", len(sisa), s)
		}
	}
	for _, k := range kel.kelompok {
		if err := jalankanKelompok(ctx, e, k, false); err != nil {
			return err
		}
	}
	return nil
}

// jalankanKelompok - satu kelompok SQL (bug DEV 10-10-2026 ORA-01722): bind diperiksa dulu (periksaBindBaris), lalu array-bind
// potongan ukuranKelompok (kolomArrayBertipe) atau - perBaris - satu eksekusi per baris dengan nilai apa adanya (n == 1 ->
// argEksekusi bind skalar, sama dengan jalur lama per baris). Galat Oracle dibungkus tabel, rentang baris, dan jumlah bind.
func jalankanKelompok(ctx context.Context, e eksekutorKelompok, k *kelompokSisip, perBaris bool) error {
	ukuran, cara := ukuranKelompok, "array-bind"
	if perBaris {
		ukuran, cara = 1, "per baris"
	}
	for awal := 0; awal < len(k.baris); awal += ukuran {
		akhir := min(awal+ukuran, len(k.baris))
		potong := k.baris[awal:akhir]
		if err := periksaBindBaris(k.q, k.tabel, potong, awal); err != nil {
			return err
		}
		kolom := kolomArray(potong)
		if !perBaris {
			var err error
			if kolom, err = kolomArrayBertipe(potong); err != nil {
				return fmt.Errorf("%w (%s)", err, k.tabel)
			}
		}
		if err := e.sisipArray(ctx, k.q, k.tabel, kolom, akhir-awal); err != nil {
			return fmt.Errorf("%w (%s baris %d-%d dari %d, %d bind, %s)", err, k.tabel, awal+1, akhir, len(k.baris), len(potong[0]), cara)
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

// PenambahObjek - tiket 75: lokasi objek DITAMBAH ke case (tanpa menghapus yang ada) - kelompok kedua dst. unggah CSV objek.
type PenambahObjek interface {
	// TambahObjek - sentuh case (kunci + 404), pastikan General, sisip baris dengan SEQ_NO lokasi mulai awal+1.
	TambahObjek(ctx context.Context, tx *db.Tx, id string, baris []models.ObjekFire, awal int) error
}

// TambahObjek - lihat PenambahObjek.
func (r *ObjekOracle) TambahObjek(ctx context.Context, tx *db.Tx, id string, baris []models.ObjekFire, awal int) error {
	t, err := r.tabel()
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, sqlSentuhCase(t.work), "menyentuh "+TabelWorkPolis, id, LiniFacIn)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrKasusTidakAda
	}
	if _, err := jalankan(ctx, tx, sqlPastikanGeneral(t.general), "memastikan "+TabelGeneralPolis, id, id); err != nil {
		return err
	}
	return sisipBerkelompokMulai(ctx, oracleKelompok{d: r.db, tx: tx}, t, id, baris, uidAcak, awal)
}

// sisipObjekBerkelompok - GantiObjek: jalur berkelompok atas transaksi tx.
func (r *ObjekOracle) sisipObjekBerkelompok(ctx context.Context, tx *db.Tx, t tabelObjek, id string, baris []models.ObjekFire) error {
	return sisipBerkelompok(ctx, oracleKelompok{d: r.db, tx: tx}, t, id, baris, uidAcak)
}
