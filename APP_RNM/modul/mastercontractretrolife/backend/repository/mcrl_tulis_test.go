package repository

// Penulis (paket 2+) - identitas dan teks SQL, tanpa Oracle.

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

func TestFormatIdentitasSepertiProcedure(t *testing.T) {
	for n, mau := range map[int64]string{44: "1000044", 1: "1000001", 999999: "1999999"} {
		if got, err := FormatIdentitas(n); err != nil || got != mau {
			t.Errorf("FormatIdentitas(%d) = %q, %v; mau %q", n, got, err, mau)
		}
	}
	// ⛔ LPAD Oracle memotong diam-diam; di sini gagal terang.
	for _, n := range []int64{1000000, -1} {
		if _, err := FormatIdentitas(n); !errors.Is(err, ErrIdentitasMelampauiLebar) {
			t.Errorf("FormatIdentitas(%d): %v, mau ErrIdentitasMelampauiLebar", n, err)
		}
	}
}

func TestSQLPenulisTahun(t *testing.T) {
	kasus := map[string]struct {
		sql   string
		wajib []string
	}{
		"sisip": {sqlSisipTahun("S.TREATYYEAR_LIFE"), []string{"INSERT INTO S.TREATYYEAR_LIFE",
			"VALUES (:1, :2, :3, :4, SYSDATE, TO_DATE(:5, 'YYYY-MM-DD'), TO_DATE(:6, 'YYYY-MM-DD'))"}},
		"perbarui": {sqlPerbaruiTahun("S.TREATYYEAR_LIFE"), []string{"TGLUPDATE = SYSDATE", "WHERE ID = :6"}},
		"salin business": {sqlSalinTahunKeBusiness("S.TREATYBUSINESS_LIFE"), []string{"SET TREATYYEAR = :1",
			"WHERE TREATYYEARID = :3 AND DECODE(TREATYYEAR, :4, 0, 1) = 1"}},
		"salin kontrak": {sqlSalinTahunKeKontrak("S.TREATYCONTRACT_LIFE"), []string{"WHERE IDTREATYYEAR = :4",
			"DECODE(TREATYSTARTDATE, TO_DATE(:5, 'YYYY-MM-DD'), 0, 1) = 1"}},
		"sequence": {sqlNomorBerikut("S.TREATYYEAR_LIFE_SEQ"), []string{"SELECT S.TREATYYEAR_LIFE_SEQ.NEXTVAL FROM DUAL"}},
	}
	for nama, k := range kasus {
		rata := strings.Join(strings.Fields(k.sql), " ")
		for _, w := range k.wajib {
			if !strings.Contains(rata, w) {
				t.Errorf("%s: tanpa %q:\n%s", nama, w, rata)
			}
		}
		if err := db.PeriksaSQL(k.sql); err != nil {
			t.Errorf("%s: %v", nama, err)
		}
	}
}

func TestPecahAngkaKebalNLS(t *testing.T) {
	for masuk, mau := range map[string]struct {
		koef  string
		skala int64
	}{
		"1500000000.10": {"150000000010", 2}, "0.000000001": {"1", 9}, "-12.5": {"-125", 1},
		"0": {"0", 0}, "100": {"100", 0}, "0.00": {"0", 0},
	} {
		d, _, err := apd.NewFromString(masuk)
		if err != nil {
			t.Fatal(err)
		}
		koef, skala := PecahAngka(d)
		if koef != mau.koef || skala != mau.skala {
			t.Errorf("PecahAngka(%s) = (%v, %d), mau (%s, %d)", masuk, koef, skala, mau.koef, mau.skala)
		}
	}
	if koef, skala := PecahAngka(nil); koef != nil || skala != 0 {
		t.Errorf("nil = NULL: (%v, %d)", koef, skala)
	}
}

// Jumlah penampung SQL = jumlah argumen bind yang dikirim penulis.
func TestPenampungKontrakCocokArgumen(t *testing.T) {
	k := models.Kontrak{}
	if n, mau := hitungPenampung(sqlSisipKontrak("S.T")), len(argUang(make([]any, 5), k))+2; n != mau {
		t.Errorf("sisip kontrak: %d penampung, %d argumen", n, mau)
	}
	if n, mau := hitungPenampung(sqlPerbaruiKontrak("S.T")), len(argUang(make([]any, 3), k))+3; n != mau {
		t.Errorf("ubah kontrak: %d penampung, %d argumen", n, mau)
	}
	if n := hitungPenampung(sqlSalinJenisKeAnak("S.T")); n != 6 {
		t.Errorf("salin jenis: %d penampung, mau 6", n)
	}
	if n := hitungPenampung(sqlSalinTahunKeKontrak("S.T")); n != 6 {
		t.Errorf("salin tahun ke kontrak: %d penampung, mau 6", n)
	}
	r := models.Reinsurer{}
	if n, mau := hitungPenampung(sqlSisipReinsurer("S.T")), len(argPersenReinsurer(make([]any, 8), r)); n != mau {
		t.Errorf("sisip reinsurer: %d penampung, %d argumen", n, mau)
	}
	if n, mau := hitungPenampung(sqlPerbaruiReinsurer("S.T")), len(argPersenReinsurer(make([]any, 5), r))+1; n != mau {
		t.Errorf("ubah reinsurer: %d penampung, %d argumen", n, mau)
	}
	if n := hitungPenampung(sqlSisipSecurity("S.T")); n != 9 {
		t.Errorf("sisip security: %d penampung, mau 7 + 2 angka", n)
	}
	if n := hitungPenampung(sqlPerbaruiSecurity("S.T")); n != 6 {
		t.Errorf("ubah security: %d penampung, mau 3 + 2 angka + ID", n)
	}
	if n := hitungPenampung(sqlSisipBusiness("S.T")); n != 11 {
		t.Errorf("sisip business: %d penampung, mau 11", n)
	}
	if n := hitungPenampung(sqlPerbaruiBusiness("S.T")); n != 9 {
		t.Errorf("ubah business: %d penampung, mau 9", n)
	}
}

func TestSQLTotalSharePerKontrak(t *testing.T) {
	q := strings.Join(strings.Fields(sqlTotalSharePerKontrak("S.K", "S.R", "S.Y", true)), " ")
	for _, w := range []string{"LEFT JOIN S.R r ON r.TREATYCONTRACTID = k.ID AND r.TREATYYEARID = k.IDTREATYYEAR",
		"WHERE k.IDTREATYYEAR = :1", "TO_CHAR(SUM(r.PCTSHARE), 'TM9'"} {
		if !strings.Contains(q, w) {
			t.Errorf("tanpa %q: %s", w, q)
		}
	}
	if strings.Contains(sqlTotalSharePerKontrak("S.K", "S.R", "S.Y", false), ":1") {
		t.Error("tanpa saringan tahun, tidak boleh ada penampung")
	}
}

var polaPenampung = regexp.MustCompile(`:(\d+)\b`)

// hitungPenampung - penampung BERBEDA :n, dan harus rapat 1..n.
func hitungPenampung(q string) int {
	ada := map[string]bool{}
	for _, m := range polaPenampung.FindAllStringSubmatch(q, -1) {
		ada[m[1]] = true
	}
	for i := 1; i <= len(ada); i++ {
		if !ada[strconv.Itoa(i)] {
			return -1
		}
	}
	return len(ada)
}

func TestSequenceHanyaLimaMilikModul(t *testing.T) {
	g := &Gudang{}
	if _, err := g.identitasBaru(context.Background(), nil, "TREATYSECURITY_LIFE_SEQ"); !errors.Is(err, ErrSequenceTakDikenal) {
		t.Errorf("sequence DEV yang tidak dirujuk XML harus ditolak (K3): %v", err)
	}
}
