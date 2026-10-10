package repository

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// urutTiruan - sequence tiruan: nilai ke-k tiap sequence = "<seq>-k" (urutan per sequence, seperti NEXTVAL).
type urutTiruan struct{ n map[string]int }

func (u *urutTiruan) berikut(seq string) string {
	u.n[seq]++
	return fmt.Sprintf("%s-%d", seq, u.n[seq])
}

type barisRekam struct {
	q   string
	arg []any
}

// rekamLangsung - jalur lama (penulisLangsung) tanpa basis data: satu round-trip per NEXTVAL dan per INSERT.
type rekamLangsung struct {
	u     *urutTiruan
	baris []barisRekam
	rt    int
}

func (r *rekamLangsung) nomor(_ context.Context, seq string) (string, error) {
	r.rt++
	return r.u.berikut(seq), nil
}

func (r *rekamLangsung) sisip(_ context.Context, q, _ string, arg []any) error {
	r.rt++
	r.baris = append(r.baris, barisRekam{q, append([]any(nil), arg...)})
	return nil
}

// eksekutorTiruan - eksekutorKelompok tanpa basis data: array dikembalikan menjadi baris; round-trip dihitung (pengambilan
// nomor ditambah fetch prefetch 25 baris bawaan go-ora).
type eksekutorTiruan struct {
	u       *urutTiruan
	baris   []barisRekam
	rt      int
	urutan  []string // teks SQL tiap eksekusi
	selisih int      // ID dikembalikan per sequence = n + selisih (uji pengaman)
}

func (e *eksekutorTiruan) nomorBerkelompok(_ context.Context, seq []string, n int) ([][]string, error) {
	e.rt += (n + 24) / 25
	hasil := make([][]string, len(seq))
	for r := 0; r < n+e.selisih; r++ {
		for i, s := range seq {
			hasil[i] = append(hasil[i], e.u.berikut(s))
		}
	}
	return hasil, nil
}

func (e *eksekutorTiruan) sisipArray(_ context.Context, q, _ string, kolom []any, n int) error {
	e.rt++
	e.urutan = append(e.urutan, q)
	for r := 0; r < n; r++ {
		arg := make([]any, len(kolom))
		for c := range kolom {
			v := kolom[c].([]any)
			if len(v) != n {
				return fmt.Errorf("kolom %d: %d nilai, seharusnya %d", c, len(v), n)
			}
			arg[c] = v[r]
		}
		e.baris = append(e.baris, barisRekam{q, arg})
	}
	return nil
}

func uidUrut() func() (string, error) {
	n := 0
	return func() (string, error) {
		n++
		return fmt.Sprintf("uid-%d", n), nil
	}
}

var tbObjekUji = tabelObjek{work: "W", general: "G", loc: "LOC", prop: "PROP", risk: "RISK", bang: "BANG", sekitar: "SKT", item: "ITEM",
	okupasi: "OKU", tol: "TOL", fea: "FEA", rugi: "RUGI", koas: "KOAS", cov: "COV", ded: "DED", spread: "SPR"}

// objekKecilUji - dua lokasi: item ber-coverage dan deductible (dengan / tanpa mata uang = dua teks SQL), dua okupasi, FEA, dua
// catatan kerugian; lokasi kedua tanpa item.
func objekKecilUji() []models.ObjekFire {
	d := desimalUji
	cov := func(old string) models.CoverageObjek {
		return models.CoverageObjek{Coverage: "C" + old, OldID: old, CoverageBasis: "1", Rate: d("0.25"), Premium: d("12.5"),
			Deductibles: []models.Deductible{{TypeDeductible: "1", Currency: "IDR", Amount: d("100")}, {TypeDeductible: "2"}}}
	}
	return []models.ObjekFire{
		{ObjectNo: "1", ObjectType: "Office", RoadName: "UJI RAYA", ZipCode: "12345", IsTopRisk: true,
			Items: []models.ItemObjek{{ItemType: "BUILDING", Currency: "IDR", TSI: d("1000"), Coverages: []models.CoverageObjek{cov("FLEXAS"), cov("EQ")}},
				{ItemType: "MESIN", Currency: "USD", TSI: d("5")}},
			Occupations: []models.OkupasiObjek{{OccupationID: "2001"}, {OccupationID: "2002", ConstructionClass: "KELAS 1"}},
			FEA:         []models.BarisFEA{{}},
			LossRecords: []models.CatatanKerugian{{Detail: "-"}, {Detail: "UJI", CoinsName: "UJI KOAS", Currency: "IDR", Amount: d("7")}}},
		{ObjectNo: "2", ObjectType: "Others", LossRecords: []models.CatatanKerugian{{Detail: "-"}}},
	}
}

// objekCsvUji - n lokasi bentuk Save Data CSV: satu item, satu coverage, satu okupasi, satu catatan kerugian "-".
func objekCsvUji(n int) []models.ObjekFire {
	out := make([]models.ObjekFire, n)
	for i := range out {
		out[i] = models.ObjekFire{ObjectNo: fmt.Sprint(i + 1), ObjectType: "Dwelling House", RoadName: fmt.Sprintf("UJI %d", i),
			Items: []models.ItemObjek{{ItemType: "BUILDING", Currency: "IDR", TSI: desimalUji("1000"),
				Coverages: []models.CoverageObjek{{Coverage: "100815", OldID: "FLEXAS", CoverageBasis: "1", Rate: desimalUji("0.25")}}}},
			Occupations: []models.OkupasiObjek{{OccupationID: "2001"}}, LossRecords: []models.CatatanKerugian{{Detail: "-"}}}
	}
	return out
}

// perSQL - baris per teks SQL, urutan dipertahankan; plus urutan kemunculan pertama teks SQL.
func perSQL(b []barisRekam) (map[string][][]any, []string) {
	m := map[string][][]any{}
	var urut []string
	for _, x := range b {
		if _, ada := m[x.q]; !ada {
			urut = append(urut, x.q)
		}
		m[x.q] = append(m[x.q], x.arg)
	}
	return m, urut
}

// bandingkanJalur - jalur lama vs berkelompok atas objek yang sama: baris per tabel identik (ID sequence, PARENT_ID, SEQ_NO,
// uid, nilai kolom) dan kelompok dijalankan induk sebelum anak; kembalikan round-trip keduanya.
func bandingkanJalur(t *testing.T, objek []models.ObjekFire) (int, int) {
	t.Helper()
	ctx := context.Background()
	lama := &rekamLangsung{u: &urutTiruan{n: map[string]int{}}}
	if err := (&tulisObjek{t: tbObjekUji, p: lama, uid: uidUrut()}).sisipObjek(ctx, "UJI-NB-1", objek); err != nil {
		t.Fatal(err)
	}
	baru := &eksekutorTiruan{u: &urutTiruan{n: map[string]int{}}}
	if err := sisipBerkelompok(ctx, baru, tbObjekUji, "UJI-NB-1", objek, uidUrut()); err != nil {
		t.Fatal(err)
	}
	ml, ul := perSQL(lama.baris)
	mb, ub := perSQL(baru.baris)
	if !reflect.DeepEqual(ml, mb) || !reflect.DeepEqual(ul, ub) || len(lama.baris) != len(baru.baris) {
		t.Fatalf("baris berbeda: lama %d baris %d SQL, baru %d baris %d SQL", len(lama.baris), len(ul), len(baru.baris), len(ub))
	}
	return lama.rt, baru.rt
}

// TestSisipBerkelompokSetara - K66-B13: objek kecil (semua jenis anak, dua varian SQL deductible) dan CSV 1128 baris sintetis
// menghasilkan baris IDENTIK dengan jalur lama; induk ditulis sebelum anak; round-trip turun.
func TestSisipBerkelompokSetara(t *testing.T) {
	kecilLama, kecilBaru := bandingkanJalur(t, objekKecilUji())
	csvLama, csvBaru := bandingkanJalur(t, objekCsvUji(1128))
	if csvBaru*10 > csvLama || kecilBaru > kecilLama {
		t.Errorf("round-trip: kecil %d -> %d, CSV %d -> %d", kecilLama, kecilBaru, csvLama, csvBaru)
	}
	// Urutan eksekusi: tabel induk sebelum anak.
	baru := &eksekutorTiruan{u: &urutTiruan{n: map[string]int{}}}
	if err := sisipBerkelompok(context.Background(), baru, tbObjekUji, "UJI-NB-1", objekKecilUji(), uidUrut()); err != nil {
		t.Fatal(err)
	}
	pos := func(tabel string) int {
		for i, q := range baru.urutan {
			if strings.HasPrefix(q, "INSERT INTO "+tabel+" ") {
				return i
			}
		}
		t.Fatalf("tanpa %s", tabel)
		return -1
	}
	for _, p := range [][2]string{{"LOC", "PROP"}, {"PROP", "RISK"}, {"PROP", "ITEM"}, {"ITEM", "COV"}, {"COV", "DED"},
		{"PROP", "OKU"}, {"OKU", "TOL"}, {"LOC", "FEA"}, {"PROP", "RUGI"}, {"RUGI", "KOAS"}} {
		if pos(p[0]) >= pos(p[1]) {
			t.Errorf("%s sesudah %s", p[0], p[1])
		}
	}
}

// TestSisipBerkelompokRoundTrip - jumlah round-trip penyisipan objek (tanpa sentuh / hapus yang sama di kedua jalur) untuk
// 1128 dan 3000 baris CSV (log uji); potongan ukuranKelompok.
func TestSisipBerkelompokRoundTrip(t *testing.T) {
	for _, n := range []int{1128, 3000} {
		lama, baru := bandingkanJalur(t, objekCsvUji(n))
		t.Logf("%d baris: round-trip penyisipan jalur lama %d, berkelompok %d (%.0fx)", n, lama, baru, float64(lama)/float64(baru))
	}
	if kolomArray(nil) != nil || !reflect.DeepEqual(kolomArray([][]any{{1, "a"}, {2, nil}}), []any{[]any{1, 2}, []any{"a", nil}}) {
		t.Error("kolomArray")
	}
	if !reflect.DeepEqual(argEksekusi([]any{[]any{"a"}, []any{nil}}, 1), []any{"a", nil}) ||
		!reflect.DeepEqual(argEksekusi([]any{[]any{"a", "b"}}, 2), []any{[]any{"a", "b"}}) {
		t.Error("argEksekusi")
	}
	if q := sqlNomorBerkelompok([]string{"S.A", "S.B"}); q != "SELECT S.A.NEXTVAL, S.B.NEXTVAL FROM DUAL CONNECT BY LEVEL <= :1" {
		t.Errorf("%q", q)
	}
}

// TestSisipBerkelompokPengaman - kumpulan ID yang tidak cocok dengan lintasan tulis (kurang / lebih) -> galat, tanpa eksekusi.
func TestSisipBerkelompokPengaman(t *testing.T) {
	for _, s := range []int{-1, 1} {
		e := &eksekutorTiruan{u: &urutTiruan{n: map[string]int{}}, selisih: s}
		if err := sisipBerkelompok(context.Background(), e, tbObjekUji, "UJI-NB-1", objekCsvUji(3), uidUrut()); err == nil || len(e.urutan) != 0 {
			t.Errorf("selisih %d: %v, %d eksekusi", s, err, len(e.urutan))
		}
	}
}
