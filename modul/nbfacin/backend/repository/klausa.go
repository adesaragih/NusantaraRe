package repository

// Tab Clauses kasus FIRE (tiket 47): ClauseList kasus (T_CLAUSELIST, migrasi 197 - satu baris per argumen, K47-4),
// pencarian klausa (popup Choose Clause), dan argumen master klausa.
//
// Pencarian `[terverifikasi]` `DDL\SearchClauseFireSQL_PreAct.xml` (Data-Clause; dikirim work owner 05-10-2026): Page-New
// SearchClauseListInput, `.Type = "FIRE"`, `.Language = ClauseLanguageID`, `.pyNote` = kata kunci, lalu RDB-List
// `DDL\RetrieveClauseSQL.xml` (Int-CLAUSE), PERSIS:
//
//	select * from ( select ID, decode({Language},'1',titleeng,'2',titledual,'0',titleina) as "Title",
//	  decode({Language},'1',texteng,'2',textdual,'0',textina) as "Text", Info as "Info"
//	  from clause where type = {Type}) bca where "Title" is not null and "Info" like '%' || {pyNote} || '%'
//
// atas view `DDL\CLAUSE.txt` (JSON M_CLAUSE; INFO FIRE = Description). Tanpa ESCAPE - seperti asal.
//
// K47-7 - MENYIMPANG dari RetrieveClauseSQL atas permintaan pengguna 05-10-2026: (1) "ini keyword buatin bisa huruf besar
// atau kecil dong" -> `UPPER("Info") LIKE '%' || UPPER(kata) || '%'` (asal peka huruf); (2) "kalau pilihnya bahasa
// inggris, ... ambil yg bahasa inggris juga dong, kn di tabel nya ada, kolom TEXTENG, kalau kolom itu kosong baru ambil yg
// bahasa indonesia ( kolom TEXTINA)" -> isi bahasa '1' = TEXTENG, atau TEXTINA bila TEXTENG kosong / spasi saja (asal:
// TEXTENG apa adanya). Judul dan bahasa '0' / '2' tetap seperti asal.
// Keputusan K47-5: berhalaman 10 (C-3) dengan urut ID (SQL asal tanpa urutan); argumentCount per hasil = jumlah
// `RetrieveArgumentNumberClauseSQL` (`select count(*) from m_argclausefire a where oldid = (select oldid from m_clause where
// id = {CARI1})`) - yang PostAct pakai untuk ArgumentCount, BUKAN SUMOFARGUMENT view.
//
// Argumen `[terverifikasi]` `NB FacIn\RDBList\SearchClauseArgFireSQL.xml` (Int-ARGCLAUSEFIRE), PERSIS: `select
// a.jsondata.ArgumentNumber, a.jsondata.Description, a.jsondata.DefaultValue from m_argclausefire a where oldid = (select
// oldid from m_clause where id = {...}) order by a.jsondata.ArgumentNumber`. ⚠️ Urutan = teks JSON apa adanya ("10"
// sebelum "2"). DDL M_CLAUSE (ID / OLDID VARCHAR2(7), JSONDATA) `[terverifikasi]` 05-10-2026; M_ARGCLAUSEFIRE hanya dari SQL.
//
// K47-6 (05-10-2026, bukti query work owner di DEV: ALL_OBJECTS POOLDATA hanya CLAUSE (VIEW) dan M_CLAUSE (TABLE) -
// M_ARGCLAUSEFIRE TIDAK ADA): keberadaan M_ARGCLAUSEFIRE ditanya ke katalog SEKALI per permintaan. Tidak ada -> pencarian
// tetap jalan dengan argumentCount kosong (BUKAN "0") + satu peringatan log; argumen -> daftar KOSONG + peringatan log
// (pengguna 05-10-2026: "M_ARGCLAUSEFIRE memang tidak ada, saat di klik See Clause and Argument, yg muncul itu adalah
// TEXTINA atau TEXTENG" - popup Isi Klasula hanya menampilkan isi klausa, tanpa grid argumen). Galat lain tidak ditangkap.

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelClauseList - ClauseList kasus beserta argumennya (migrasi 197).
	TabelClauseList    = "T_CLAUSELIST"
	sequenceClauseList = "SEQ_T_CLAUSELIST"
	// TabelMasterKlausa / TabelArgumenKlausa / TabelViewKlausa - master klausa warisan (dibaca saja).
	TabelMasterKlausa  = "M_CLAUSE"
	TabelArgumenKlausa = "M_ARGCLAUSEFIRE"
	TabelViewKlausa    = "CLAUSE"
	// TipeKlausaFire - `.Type = "FIRE"` SearchClauseFireSQL_PreAct.
	TipeKlausaFire = "FIRE"
)

// PenyimpanKlausa - ClauseList kasus, pencarian, dan argumen master.
type PenyimpanKlausa interface {
	// BacaKlausa - ClauseList kasus `id` urut SEQ_NO; ErrKasusTidakAda bila case tidak ada / bukan LINI Fac In.
	BacaKlausa(ctx context.Context, id string) ([]models.KlausaKasus, error)
	// GantiKlausa - ganti utuh ClauseList kasus `id` di transaksi pemanggil; ErrKasusTidakAda bila case tidak ada.
	GantiKlausa(ctx context.Context, tx *db.Tx, id string, baris []models.KlausaKasus) error
	// CariKlausa - satu halaman RetrieveClauseSQL (bahasa = kode "0" / "1" / "2") dan jumlah seluruhnya.
	CariKlausa(ctx context.Context, bahasa, kata string, nomor, ukuran int) ([]models.HasilKlausa, int, error)
	// ArgumenKlausa - SearchClauseArgFireSQL atas klausa ber-ID `id` (kosong bila tidak ada).
	ArgumenKlausa(ctx context.Context, id string) ([]models.ArgumenKlausa, error)
}

// KlausaOracle - PenyimpanKlausa atas Oracle.
type KlausaOracle struct{ db *db.DB }

// NewKlausaOracle merakit penyimpan klausa.
func NewKlausaOracle(d *db.DB) *KlausaOracle { return &KlausaOracle{db: d} }

// sqlCariKlausa - RetrieveClauseSQL + K47-7 (lihat kepala berkas), bind :1 / :2 bahasa (decode judul / isi), :3 tipe,
// :4 kata.
// Dasar bersama hitungan dan halaman.
func sqlCariKlausa(view string) string {
	return `SELECT * FROM (SELECT ID, DECODE(:1, '1', TITLEENG, '2', TITLEDUAL, '0', TITLEINA) AS "Title", ` +
		`DECODE(:2, '1', CASE WHEN TRIM(TEXTENG) IS NULL THEN TEXTINA ELSE TEXTENG END, '2', TEXTDUAL, '0', TEXTINA) AS "Text", ` +
		`INFO AS "Info" FROM ` + view + ` WHERE TYPE = :3) bca WHERE "Title" IS NOT NULL AND UPPER("Info") LIKE '%' || UPPER(:4) || '%'`
}

// sqlHitungKlausa - jumlah hasil pencarian.
func sqlHitungKlausa(dasar string) string { return "SELECT COUNT(*) FROM (" + dasar + ")" }

// sqlHalamanKlausa - satu halaman urut ID (:5 OFFSET, :6 FETCH) + jumlah argumen RetrieveArgumentNumberClauseSQL per baris.
func sqlHalamanKlausa(dasar, arg, master string) string {
	return `SELECT h.ID, h."Title", h."Text", h."Info", (SELECT COUNT(*) FROM ` + arg + ` a WHERE a.OLDID = (SELECT m.OLDID FROM ` +
		master + ` m WHERE m.ID = h.ID)) FROM (` + dasar + `) h ORDER BY h.ID OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY`
}

// sqlHalamanKlausaTanpaArgumen - sqlHalamanKlausa bila M_ARGCLAUSEFIRE tidak ada (K47-6): jumlah argumen NULL.
func sqlHalamanKlausaTanpaArgumen(dasar string) string {
	return `SELECT h.ID, h."Title", h."Text", h."Info", NULL FROM (` + dasar + `) h ORDER BY h.ID OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY`
}

// sqlAdaObjek - objek bernama itu ada di skema sasaran (katalog; pola migrasi.objekAda).
const sqlAdaObjek = "SELECT COUNT(*) FROM SYS.ALL_OBJECTS WHERE UPPER(OWNER) = UPPER(:1) AND OBJECT_NAME = :2"

// adaArgumen - M_ARGCLAUSEFIRE ada di skema (K47-6).
func (r *KlausaOracle) adaArgumen(ctx context.Context) (bool, error) {
	var n int
	if err := r.db.QueryRowContext(ctx, sqlAdaObjek, r.db.Skema(), TabelArgumenKlausa).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: katalog %s: %w", TabelArgumenKlausa, err)
	}
	return n > 0, nil
}

// sqlArgumenKlausa - SearchClauseArgFireSQL (lihat kepala berkas), placeholder ber-bind; alias kolom dibuang.
func sqlArgumenKlausa(arg, klausa string) string {
	return "SELECT a.jsondata.ArgumentNumber, a.jsondata.Description, a.jsondata.DefaultValue FROM " + arg +
		" a WHERE a.OLDID = (SELECT c.OLDID FROM " + klausa + " c WHERE c.ID = :1) ORDER BY a.jsondata.ArgumentNumber"
}

// sqlBacaKlausa - seluruh baris klausa kasus (satu per argumen), urut klausa lalu argumen; baris tanpa argumen lebih dulu.
func sqlBacaKlausa(cl string) string {
	return "SELECT SEQ_NO, ARGUMENT_SEQ_NO, CLAUSE_CODE, CLAUSE_TITLE, CLAUSE_DESCRIPTION, CLAUSE_LANGUAGE, CLAUSE_LANGUAGE_ID," +
		" CLAUSE_CONTENT, CLAUSE_CONTENT_TEMP, ARGUMENT_COUNT, ARGUMENT_NUMBER, ARGUMENT_DESCRIPTION, ARGUMENT_VALUE FROM " + cl +
		" WHERE PARENT_ID = :1 ORDER BY SEQ_NO, ARGUMENT_SEQ_NO NULLS FIRST"
}

func sqlHapusKlausa(cl string) string { return "DELETE FROM " + cl + " WHERE PARENT_ID = :1" }

func sqlSisipKlausa(cl string) string {
	return "INSERT INTO " + cl + " (ID, PARENT_ID, SEQ_NO, ARGUMENT_SEQ_NO, ROW_UID, CLAUSE_CODE, CLAUSE_TITLE, CLAUSE_DESCRIPTION," +
		" CLAUSE_LANGUAGE, CLAUSE_LANGUAGE_ID, CLAUSE_CONTENT, CLAUSE_CONTENT_TEMP, ARGUMENT_COUNT, ARGUMENT_NUMBER," +
		" ARGUMENT_DESCRIPTION, ARGUMENT_VALUE) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16)"
}

// barisKlausa - satu baris T_CLAUSELIST (satu argumen, atau klausa tanpa argumen bila UrutArg = 0).
type barisKlausa struct {
	Urut, UrutArg int
	Klausa        models.KlausaKasus
	Argumen       models.ArgumenKlausa
}

// pecahKlausa - ClauseList -> baris tabel: satu per argumen (UrutArg 1..n), klausa tanpa argumen = satu baris UrutArg 0
// (ARGUMENT_SEQ_NO kosong); Urut = urutan klausa mulai 1.
func pecahKlausa(daftar []models.KlausaKasus) []barisKlausa {
	var out []barisKlausa
	for i, k := range daftar {
		if len(k.Arguments) == 0 {
			out = append(out, barisKlausa{Urut: i + 1, Klausa: k})
			continue
		}
		for j, a := range k.Arguments {
			out = append(out, barisKlausa{Urut: i + 1, UrutArg: j + 1, Klausa: k, Argumen: a})
		}
	}
	return out
}

// kumpulkanKlausa - kebalikan pecahKlausa atas baris urut (Urut, UrutArg): baris ber-Urut sama = satu klausa, kolom
// klausa dari baris pertamanya; argumen dari baris ber-UrutArg > 0. argumentList selalu larik.
func kumpulkanKlausa(baris []barisKlausa) []models.KlausaKasus {
	hasil := []models.KlausaKasus{}
	for i, b := range baris {
		if i == 0 || b.Urut != baris[i-1].Urut {
			k := b.Klausa
			k.Arguments = []models.ArgumenKlausa{}
			hasil = append(hasil, k)
		}
		if b.UrutArg > 0 {
			k := &hasil[len(hasil)-1]
			k.Arguments = append(k.Arguments, b.Argumen)
		}
	}
	return hasil
}

// nama - nama berskema T_CLAUSELIST, T_WORK_POLIS, T_GENERAL_POLIS.
func (r *KlausaOracle) nama() (cl, work, general string, err error) {
	for _, x := range []struct {
		t string
		p *string
	}{{TabelClauseList, &cl}, {TabelWorkPolis, &work}, {TabelGeneralPolis, &general}} {
		if *x.p, err = r.db.Qualify(x.t); err != nil {
			return
		}
	}
	return
}

// CariKlausa - lihat PenyimpanKlausa.
func (r *KlausaOracle) CariKlausa(ctx context.Context, bahasa, kata string, nomor, ukuran int) ([]models.HasilKlausa, int, error) {
	view, err := r.db.Qualify(TabelViewKlausa)
	if err != nil {
		return nil, 0, err
	}
	arg, err := r.db.Qualify(TabelArgumenKlausa)
	if err != nil {
		return nil, 0, err
	}
	master, err := r.db.Qualify(TabelMasterKlausa)
	if err != nil {
		return nil, 0, err
	}
	ada, err := r.adaArgumen(ctx)
	if err != nil {
		return nil, 0, err
	}
	dasar := sqlCariKlausa(view)
	halaman := sqlHalamanKlausa(dasar, arg, master)
	if !ada {
		log.Printf("nbfacin: peringatan: tabel %s tidak ada di skema %s - argumentCount pencarian klausa dikosongkan (K47-6)",
			TabelArgumenKlausa, r.db.Skema())
		halaman = sqlHalamanKlausaTanpaArgumen(dasar)
	}
	bind := []any{bahasa, bahasa, TipeKlausaFire, kata}
	var total int
	if err := r.db.QueryRowContext(ctx, sqlHitungKlausa(dasar), bind...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("repository: hitung %s: %w", TabelViewKlausa, err)
	}
	baris, err := r.db.QueryContext(ctx, halaman, append(bind, (nomor-1)*ukuran, ukuran)...)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: baca %s: %w", TabelViewKlausa, err)
	}
	defer baris.Close()
	hasil := []models.HasilKlausa{}
	for baris.Next() {
		var id, judul, isi, info sql.NullString
		var n sql.NullInt64
		if err := baris.Scan(&id, &judul, &isi, &info, &n); err != nil {
			return nil, 0, fmt.Errorf("repository: %s: %w", TabelViewKlausa, err)
		}
		jumlah := "" // tabel argumen tidak ada (K47-6)
		if n.Valid {
			jumlah = strconv.FormatInt(n.Int64, 10)
		}
		hasil = append(hasil, models.HasilKlausa{ID: id.String, Title: judul.String, Text: isi.String, Info: info.String,
			ArgumentCount: jumlah})
	}
	return hasil, total, baris.Err()
}

// ArgumenKlausa - lihat PenyimpanKlausa.
func (r *KlausaOracle) ArgumenKlausa(ctx context.Context, id string) ([]models.ArgumenKlausa, error) {
	if ada, err := r.adaArgumen(ctx); err != nil {
		return nil, err
	} else if !ada {
		log.Printf("nbfacin: peringatan: tabel %s tidak ada di skema %s - argumen klausa %s kosong (K47-6)",
			TabelArgumenKlausa, r.db.Skema(), id)
		return []models.ArgumenKlausa{}, nil
	}
	arg, err := r.db.Qualify(TabelArgumenKlausa)
	if err != nil {
		return nil, err
	}
	klausa, err := r.db.Qualify(TabelMasterKlausa)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlArgumenKlausa(arg, klausa), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelArgumenKlausa, err)
	}
	defer baris.Close()
	hasil := []models.ArgumenKlausa{}
	for baris.Next() {
		var n, d, v sql.NullString
		if err := baris.Scan(&n, &d, &v); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelArgumenKlausa, err)
		}
		hasil = append(hasil, models.ArgumenKlausa{Number: n.String, Description: d.String, Value: v.String})
	}
	return hasil, baris.Err()
}

// BacaKlausa - lihat PenyimpanKlausa. Baris berurutan SEQ_NO dikumpulkan menjadi satu klausa; kolom klausa dari baris
// pertamanya, argumen dari baris ber-ARGUMENT_SEQ_NO.
func (r *KlausaOracle) BacaKlausa(ctx context.Context, id string) ([]models.KlausaKasus, error) {
	cl, work, _, err := r.nama()
	if err != nil {
		return nil, err
	}
	var n int
	if err := r.db.QueryRowContext(ctx, sqlAdaKasus(work), id, LiniFacIn).Scan(&n); err != nil {
		return nil, fmt.Errorf("repository: cek case NB: %w", err)
	}
	if n == 0 {
		return nil, ErrKasusTidakAda
	}
	baris, err := r.db.QueryContext(ctx, sqlBacaKlausa(cl), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelClauseList, err)
	}
	defer baris.Close()
	var dibaca []barisKlausa
	for baris.Next() {
		var urut int
		var urutArg sql.NullInt64
		var v [11]sql.NullString
		ptr := []any{&urut, &urutArg}
		for i := range v {
			ptr = append(ptr, &v[i])
		}
		if err := baris.Scan(ptr...); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelClauseList, err)
		}
		dibaca = append(dibaca, barisKlausa{Urut: urut, UrutArg: int(urutArg.Int64),
			Klausa: models.KlausaKasus{Code: v[0].String, Title: v[1].String, Description: v[2].String, Language: v[3].String,
				LanguageID: v[4].String, Content: v[5].String, ContentTemp: v[6].String, ArgumentCount: v[7].String},
			Argumen: models.ArgumenKlausa{Number: v[8].String, Description: v[9].String, Value: v[10].String}})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelClauseList, err)
	}
	return kumpulkanKlausa(dibaca), nil
}

// GantiKlausa - lihat PenyimpanKlausa. Urutan pola GantiObjek: sentuh case (kunci + 404), pastikan General, hapus, sisip
// ulang - satu baris per argumen (klausa tanpa argumen: satu baris ber-ARGUMENT_SEQ_NO kosong), SEQ_NO 1..n per klausa.
func (r *KlausaOracle) GantiKlausa(ctx context.Context, tx *db.Tx, id string, baris []models.KlausaKasus) error {
	cl, work, general, err := r.nama()
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, sqlSentuhCase(work), "menyentuh "+TabelWorkPolis, id, LiniFacIn)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrKasusTidakAda
	}
	if _, err := jalankan(ctx, tx, sqlPastikanGeneral(general), "memastikan "+TabelGeneralPolis, id, id); err != nil {
		return err
	}
	if _, err := jalankan(ctx, tx, sqlHapusKlausa(cl), "menghapus klausa", id); err != nil {
		return err
	}
	k := db.KosongJadiNil
	for _, b := range pecahKlausa(baris) {
		var urutArg any // NULL = klausa tanpa argumen
		if b.UrutArg > 0 {
			urutArg = b.UrutArg
		}
		idBaris, err := r.db.NomorBerikut(ctx, tx, sequenceClauseList)
		if err != nil {
			return err
		}
		uid, err := uidAcak()
		if err != nil {
			return err
		}
		c, a := b.Klausa, b.Argumen
		if _, err := jalankan(ctx, tx, sqlSisipKlausa(cl), "menyisip "+TabelClauseList, idBaris, id, b.Urut, urutArg, uid,
			k(c.Code), k(c.Title), k(c.Description), k(c.Language), k(c.LanguageID), k(c.Content), k(c.ContentTemp),
			k(c.ArgumentCount), k(a.Number), k(a.Description), k(a.Value)); err != nil {
			return err
		}
	}
	return nil
}
