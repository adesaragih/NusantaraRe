package repository

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// bindBerurut - nomor bind :n dalam urutan kemunculan = 1..n (go-ora mengikat menurut urutan); kembalikan jumlahnya.
func bindBerurut(t *testing.T, q string) int {
	t.Helper()
	b := regexp.MustCompile(`:(\d+)`).FindAllStringSubmatch(q, -1)
	for i, m := range b {
		if m[1] != strconv.Itoa(i+1) {
			t.Errorf("bind ke-%d = :%s\n%s", i+1, m[1], q)
		}
	}
	return len(b)
}

// tabelHapus - tabel sasaran tiap DELETE.
func tabelHapus(q []string) []string {
	var hasil []string
	for _, x := range q {
		hasil = append(hasil, strings.Fields(strings.TrimPrefix(x, "DELETE FROM "))[0])
	}
	return hasil
}

// TestSqlHapusLokasi - tiket 75 F3: hapus pohon lokasi terpilih = sqlHapusObjek (tabel dan urutan anak -> induk SAMA) dengan
// syarat case :1 DITAMBAH `ID IN (:2, :3)`; bind berurutan, jumlah bind = 1 + jumlah lokasi.
func TestSqlHapusLokasi(t *testing.T) {
	tb := tabelObjekUji()
	q, err := sqlHapusLokasi(tb, 2)
	if err != nil {
		t.Fatal(err)
	}
	if a, b := tabelHapus(sqlHapusObjek(tb)), tabelHapus(q); strings.Join(a, ",") != strings.Join(b, ",") {
		t.Errorf("tabel / urutan beda:\n%v\n%v", a, b)
	}
	for _, x := range q {
		if n := bindBerurut(t, x); n != 3 || strings.Count(x, "ID IN (:2, :3)") != 1 || strings.Count(x, "PARENT_ID = :1") != 1 {
			t.Errorf("%d bind:\n%s", n, x)
		}
		if err := db.PeriksaSQL(x); err != nil {
			t.Error(err)
		}
	}
	if akhir := q[len(q)-1]; akhir != "DELETE FROM "+tb.loc+" WHERE PARENT_ID = :1 AND ID IN (:2, :3)" {
		t.Error(akhir)
	}
}

// TestSqlPerubahanObjek - kueri pendukung PATCH: bind berurutan; spreading hanya coverage lokasi terpilih; SEQ_NO terakhir
// hanya baris case :1.
func TestSqlPerubahanObjek(t *testing.T) {
	tb := tabelObjekUji()
	if q := sqlSeqLokasi(tb.loc, 3); bindBerurut(t, q) != 4 || !strings.HasSuffix(q, "WHERE PARENT_ID = :1 AND ID IN (:2, :3, :4)") {
		t.Error(q)
	}
	q, arg, err := lingkupBaca{id: "K", lokasi: []string{"91"}}.kueri(sqlAdaSpreadingLokasi(tb))
	if err != nil || bindBerurut(t, q) != len(arg) || !strings.Contains(q, "l.PARENT_ID = :1 AND l.ID IN (:2)") ||
		!strings.HasPrefix(q, "SELECT COUNT(*) FROM "+tb.spread+" WHERE PARENT_ID IN (SELECT v.ID FROM "+tb.cov) {
		t.Errorf("%v %v\n%s", err, arg, q)
	}
	r := sqlSeqTerakhir(tb.loc)
	if bindBerurut(t, r) != 1 || r != "SELECT NVL(MAX(SEQ_NO), 0), COUNT(*) FROM "+tb.loc+" WHERE PARENT_ID = :1" {
		t.Error(r)
	}
	for _, x := range []string{r, q} {
		if err := db.PeriksaSQL(x); err != nil {
			t.Error(err)
		}
	}
}

// TestSisipBerkelompokPosisi - lokasi diubah ditulis ulang di SEQ_NO-nya (seq per baris); tanpa seq tetap awal+1..; jumlah seq
// harus sama dengan jumlah lokasi.
func TestSisipBerkelompokPosisi(t *testing.T) {
	ctx := context.Background()
	objek := []models.ObjekFire{{ObjectNo: "1", ObjectType: "UJI"}, {ObjectNo: "2", ObjectType: "UJI"}}
	seqLokasi := func(seq []int, awal int) []any {
		e := &eksekutorTiruan{u: &urutTiruan{n: map[string]int{}}, seragam: true}
		if err := sisipBerkelompokPosisi(ctx, e, tbObjekUji, "UJI-NB-1", objek, uidUrut(), awal, seq); err != nil {
			t.Fatal(err)
		}
		var hasil []any
		for _, b := range e.baris {
			if strings.HasPrefix(b.q, "INSERT INTO "+tbObjekUji.loc+" ") {
				hasil = append(hasil, b.arg[2])
			}
		}
		return hasil
	}
	if s := seqLokasi([]int{7, 3}, 0); len(s) != 2 || s[0] != 7 || s[1] != 3 {
		t.Errorf("seq: %v", s)
	}
	if s := seqLokasi(nil, 40); len(s) != 2 || s[0] != 41 || s[1] != 42 {
		t.Errorf("awal: %v", s)
	}
	e := &eksekutorTiruan{u: &urutTiruan{n: map[string]int{}}}
	if err := sisipBerkelompokPosisi(ctx, e, tbObjekUji, "UJI-NB-1", objek, uidUrut(), 0, []int{1}); err == nil {
		t.Error("seq kurang diterima")
	}
}
