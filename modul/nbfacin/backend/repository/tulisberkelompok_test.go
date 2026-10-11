package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hash/fnv"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// urutTiruan - sequence tiruan: nilai ke-k tiap sequence = dasar(seq) + k, teks digit seperti NEXTVAL (kolom NUMBER(19),
// periksaBindBaris); dasar dari nama sequence sehingga sama di kedua jalur apa pun urutan pengambilannya.
type urutTiruan struct{ n map[string]int }

func (u *urutTiruan) berikut(seq string) string {
	h := fnv.New32a()
	h.Write([]byte(seq))
	u.n[seq]++
	return fmt.Sprint(int64(h.Sum32()%9000+1000)*1000000 + int64(u.n[seq]))
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
	seragam bool     // jalur Aneka: setiap kolom array wajib lolos periksaKolomSeragam
	gagal   string   // tabel yang eksekusinya ditolak (uji pembungkus galat)
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

func (e *eksekutorTiruan) sisipArray(_ context.Context, q, tabel string, kolom []any, n int) error {
	if tabel == e.gagal {
		return errors.New("ORA-01722: invalid number")
	}
	e.rt++
	e.urutan = append(e.urutan, q)
	if e.seragam {
		if err := periksaKolomSeragam(q, kolom); err != nil {
			return err
		}
	}
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

// bindAngkaMasuk - nomor bind (:n) yang diterima angkaMasuk (TO_NUMBER ber-mask) dalam teks SQL.
var bindAngkaMasuk = regexp.MustCompile(`TO_NUMBER\(:(\d+), 'FM`)

// periksaKolomSeragam - instrumen bug DEV 10-10-2026 (ORA-01722 T_SPREADINGLIST, lihat kolomArrayBertipe): setiap kolom
// array-bind SATU tipe Go untuk semua baris dan tanpa nil telanjang; nilai teks di bind angkaMasuk wajib cocok PolaAngkaMasuk.
func periksaKolomSeragam(q string, kolom []any) error {
	angka := map[int]bool{}
	for _, m := range bindAngkaMasuk.FindAllStringSubmatch(q, -1) {
		n, _ := strconv.Atoi(m[1])
		angka[n-1] = true
	}
	for c, k := range kolom {
		v := k.([]any)
		for r, x := range v {
			if x == nil && len(v) == 1 {
				continue // n == 1: bind skalar (argEksekusi), NULL sah seperti jalur per baris
			}
			if x == nil {
				return fmt.Errorf("kolom %d baris %d: nil telanjang", c+1, r+1)
			}
			if len(v) > 1 && reflect.TypeOf(x) != reflect.TypeOf(v[0]) {
				return fmt.Errorf("kolom %d: tipe %T baris 1, %T baris %d", c+1, v[0], x, r+1)
			}
			s, ok := x.(string)
			if n, adaN := x.(sql.NullString); adaN {
				s, ok = n.String, n.Valid
			}
			if ok && angka[c] && !PolaAngkaMasuk.MatchString(s) {
				return fmt.Errorf("kolom %d baris %d: %q di luar mask angkaMasuk", c+1, r+1, s)
			}
		}
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
	baru := &eksekutorTiruan{u: &urutTiruan{n: map[string]int{}}, seragam: true}
	if err := sisipBerkelompok(ctx, baru, tbObjekUji, "UJI-NB-1", objek, uidUrut()); err != nil {
		t.Fatal(err)
	}
	ml, ul := perSQL(lama.baris)
	mb, ub := perSQL(nilaiAsli(baru.baris))
	if !reflect.DeepEqual(ml, mb) || !reflect.DeepEqual(ul, ub) || len(lama.baris) != len(baru.baris) {
		t.Fatalf("baris berbeda: lama %d baris %d SQL, baru %d baris %d SQL", len(lama.baris), len(ul), len(baru.baris), len(ub))
	}
	return lama.rt, baru.rt
}

// nilaiAsli - baris berkelompok dengan sql.NullString (kolomArrayBertipe) dikembalikan ke nilai jalur lama (Valid -> teks,
// tidak Valid -> nil): isi data yang disimpan tidak berubah.
func nilaiAsli(baris []barisRekam) []barisRekam {
	hasil := make([]barisRekam, len(baris))
	for i, b := range baris {
		arg := make([]any, len(b.arg))
		for c, x := range b.arg {
			arg[c] = x
			if n, ok := x.(sql.NullString); ok {
				arg[c] = nil
				if n.Valid {
					arg[c] = n.String
				}
			}
		}
		hasil[i] = barisRekam{b.q, arg}
	}
	return hasil
}

// TestSisipBerkelompokFireNullTerakhir - bug DEV 10-10-2026 (pencegahan FIRE, work owner 10-10-2026): fixture objekKecilUji
// memuat kolom T_LOCATIONLIST bercampur NULL; jalur FIRE lolos periksaKolomSeragam, dan kolom yang NULL di baris TERAKHIR
// setelah baris berisi (bentuk pemicu ORA-01722) memang ada di fixture - instrumen diuji atas bentuk kolomArray lama.
func TestSisipBerkelompokFireNullTerakhir(t *testing.T) {
	ctx := context.Background()
	baru := &eksekutorTiruan{u: &urutTiruan{n: map[string]int{}}, seragam: true}
	if err := sisipBerkelompok(ctx, baru, tbObjekUji, "UJI-NB-1", objekKecilUji(), uidUrut()); err != nil {
		t.Fatal(err)
	}
	lama := &rekamLangsung{u: &urutTiruan{n: map[string]int{}}}
	if err := (&tulisObjek{t: tbObjekUji, p: lama, uid: uidUrut()}).sisipObjek(ctx, "UJI-NB-1", objekKecilUji()); err != nil {
		t.Fatal(err)
	}
	ml, ul := perSQL(lama.baris)
	nullTerakhir := 0
	for _, q := range ul {
		baris := ml[q]
		if len(baris) < 2 {
			continue
		}
		kolom := kolomArray(baris)
		for _, k := range kolom {
			v := k.([]any)
			adaIsi := false
			for _, x := range v[:len(v)-1] {
				adaIsi = adaIsi || x != nil
			}
			if adaIsi && v[len(v)-1] == nil {
				nullTerakhir++
			}
		}
		if strings.HasPrefix(q, "INSERT INTO LOC ") && periksaKolomSeragam(q, kolom) == nil {
			t.Error("instrumen: kolom T_LOCATIONLIST bercampur NULL bentuk lama lolos")
		}
	}
	if nullTerakhir == 0 {
		t.Error("fixture tanpa kolom NULL di baris terakhir setelah baris berisi")
	}
	t.Logf("kolom NULL di baris terakhir setelah baris berisi (jalur lama): %d", nullTerakhir)
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
	baru := &eksekutorTiruan{u: &urutTiruan{n: map[string]int{}}, seragam: true}
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

// TestKolomArrayBertipe - bug DEV 10-10-2026: kolom teks yang memuat nil -> sql.NullString di setiap baris (NULL di posisi mana
// pun, termasuk baris terakhir); kolom tanpa nil apa adanya;
// kolom bercampur tipe / tipe bukan teks bercampur NULL -> galat. Instrumen periksaKolomSeragam diuji dulu atas jawaban yang
// diketahui: bentuk kolomArray lama (nil di baris terakhir) DITOLAK, bentuk kolomArrayBertipe diterima.
func TestKolomArrayBertipe(t *testing.T) {
	q := "INSERT INTO S.SPR (A, B, C) VALUES (:1, :2, " + angkaMasuk(":3") + ")"
	baris := [][]any{{"1", 1, "40"}, {"2", 2, nil}} // bentuk DEV: desimal baris terakhir NULL
	if err := periksaKolomSeragam(q, kolomArray(baris)); err == nil {
		t.Error("instrumen: nil di baris terakhir lolos")
	}
	for _, v := range []any{"0.123456789", sql.NullString{String: "1E+3", Valid: true}} {
		if err := periksaKolomSeragam(q, []any{[]any{"1"}, []any{1}, []any{v}}); err == nil {
			t.Errorf("instrumen: %v di angkaMasuk lolos", v)
		}
	}
	if err := periksaKolomSeragam(q, []any{[]any{"1", "2"}, []any{1, "2"}, []any{"1", "2"}}); err == nil {
		t.Error("instrumen: kolom bercampur tipe lolos")
	}
	kolom, err := kolomArrayBertipe(baris)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(kolom, []any{[]any{"1", "2"}, []any{1, 2}, []any{sql.NullString{String: "40", Valid: true}, sql.NullString{}}}) {
		t.Errorf("kolom: %v", kolom)
	}
	if err := periksaKolomSeragam(q, kolom); err != nil {
		t.Errorf("instrumen menolak bentuk benar: %v", err)
	}
	if k, err := kolomArrayBertipe([][]any{{nil}, {nil}}); err != nil || !reflect.DeepEqual(k, []any{[]any{sql.NullString{}, sql.NullString{}}}) {
		t.Errorf("kolom NULL semua: %v %v", k, err)
	}
	for _, salah := range [][][]any{{{1}, {nil}}, {{"a"}, {1}}, {{nil}, {2}}} {
		if _, err := kolomArrayBertipe(salah); err == nil {
			t.Errorf("%v: tanpa galat", salah)
		}
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
