package repository

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// TestLingkupBacaKueri - tiket 75 F3: setiap kueri baca pohon objek memuat TEPAT satu pembatasLokasi (syarat lingkup halaman);
// lingkup halaman menyisipkan `AND l.ID IN (:2, :3)` sesudahnya, bind berurutan :1..:n sesuai urutan argumen.
func TestLingkupBacaKueri(t *testing.T) {
	tb := tabelObjekUji()
	for nama, q := range map[string]string{"objek": sqlBacaObjek(tb), "item": sqlBacaItem(tb), "coverage": sqlBacaCoverage(tb),
		"deductible": sqlBacaDeductible(tb), "okupasi": sqlBacaOkupasi(tb), "fea": sqlBacaFEA(tb), "kerugian": sqlBacaKerugian(tb)} {
		semua, arg, err := lingkupBaca{id: "K"}.kueri(q)
		if err != nil || semua != q || !reflect.DeepEqual(arg, []any{"K"}) {
			t.Errorf("%s semua: %v %v", nama, err, arg)
		}
		h, arg, err := lingkupBaca{id: "K", lokasi: []string{"11", "12"}}.kueri(q)
		if err != nil || !strings.Contains(h, "l.PARENT_ID = :1 AND l.ID IN (:2, :3)") || !reflect.DeepEqual(arg, []any{"K", "11", "12"}) {
			t.Errorf("%s halaman: %v %v\n%s", nama, err, arg, h)
		}
		if b := regexp.MustCompile(`:(\d+)`).FindAllStringSubmatch(h, -1); len(b) != 3 {
			t.Errorf("%s bind: %v", nama, b)
		} else {
			for i, m := range b {
				if m[1] != strconv.Itoa(i+1) {
					t.Errorf("%s bind ke-%d = :%s", nama, i+1, m[1])
				}
			}
		}
	}
	if _, _, err := (lingkupBaca{id: "K", lokasi: []string{"1"}}).kueri("SELECT 1 FROM X"); err == nil {
		t.Error("kueri tanpa pola diterima")
	}
	if _, _, err := (lingkupBaca{id: "K", lokasi: []string{}}).kueri(sqlBacaFEA(tb)); err == nil {
		t.Error("lingkup kosong diterima")
	}
}

// TestSqlHalamanObjek - bind berurutan menurut kemunculan (go-ora): tanpa cari :1 case, :2 / :3 OFFSET / FETCH; dengan cari
// :2..:5 pola, :6 / :7. Jumlah argumen = jumlah bind.
func TestSqlHalamanObjek(t *testing.T) {
	tb := tabelObjekUji()
	for _, cari := range []string{"", "a"} {
		ada := cari != ""
		for _, q := range []string{sqlJumlahLokasiCari(tb, ada), sqlIDLokasiHalaman(tb, ada)} {
			b := regexp.MustCompile(`:(\d+)`).FindAllStringSubmatch(q, -1)
			for i, m := range b {
				if m[1] != strconv.Itoa(i+1) {
					t.Errorf("bind ke-%d = :%s\n%s", i+1, m[1], q)
				}
			}
			n := len(argCari("K", cari))
			if strings.Contains(q, "OFFSET") {
				n += 2
			}
			if len(b) != n {
				t.Errorf("bind %d, argumen %d\n%s", len(b), n, q)
			}
		}
	}
	if q := sqlIDLokasiHalaman(tb, false); !strings.HasSuffix(q, "ORDER BY l.SEQ_NO, l.ID OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY") {
		t.Error(q)
	}
}

// TestPolaCariObjek - memuat, huruf besar, wildcard LIKE pengguna di-escape.
func TestPolaCariObjek(t *testing.T) {
	for masuk, mau := range map[string]string{" jl. 5_a ": `%JL. 5\_A%`, `100%`: `%100\%%`, `a\b`: `%A\\B%`} {
		if got := polaCari(masuk); got != mau {
			t.Errorf("%q -> %q, mau %q", masuk, got, mau)
		}
	}
	if !reflect.DeepEqual(argCari("K", "  "), []any{"K"}) {
		t.Error("cari kosong")
	}
}
