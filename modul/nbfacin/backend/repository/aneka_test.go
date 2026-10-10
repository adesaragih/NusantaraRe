package repository

import (
	"context"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
)

var polaBind = regexp.MustCompile(`:([0-9]+)`)

// bindTerbesar - nomor bind terbesar di teks SQL (go-ora: bind menurut urutan kemunculan; nomor wajib 1..n tanpa celah).
func bindTerbesar(t *testing.T, q string) int {
	t.Helper()
	lihat := map[string]bool{}
	for _, m := range polaBind.FindAllStringSubmatch(q, -1) {
		if lihat[m[1]] {
			t.Errorf("bind :%s muncul dua kali: %s", m[1], q)
		}
		lihat[m[1]] = true
	}
	for i := 1; i <= len(lihat); i++ {
		if !lihat[strconv.Itoa(i)] {
			t.Errorf("bind :%d tidak ada (celah): %s", i, q)
		}
	}
	return len(lihat)
}

func tabelAnekaUji() tabelAneka {
	return tabelAneka{tabelObjek: tabelObjek{work: "S.W", general: "S.G", loc: "S.LOC", prop: "S.PROP", risk: "S.RISK",
		okupasi: "S.OKU", cov: "S.COV", ded: "S.DED", spread: "S.SPR", item: "S.ITEM", bang: "S.BANG", sekitar: "S.SEK",
		tol: "S.TOL", fea: "S.FEA", rugi: "S.RUGI", koas: "S.KOAS"}, aneka: "S.ANEKA", marine: "S.MAR", aviasi: "S.AVI",
		he: "S.HE", memo: "S.MEMO"}
}

// kolomMigrasi214 - kolom CREATE TABLE `tabel` di migrasi 214 (urai teks berkas).
func kolomMigrasi214(t *testing.T, tabel string) map[string]bool {
	t.Helper()
	b, err := os.ReadFile("../migrations/214_objek_aneka.sql")
	if err != nil {
		t.Fatal(err)
	}
	m := regexp.MustCompile(`(?s)CREATE TABLE \{skema\}\.` + tabel + ` \((.*?)\n\)`).FindStringSubmatch(strings.ReplaceAll(string(b), "\r", ""))
	if m == nil {
		t.Fatalf("%s tidak ada di 214", tabel)
	}
	k := map[string]bool{}
	for _, baris := range strings.Split(m[1], "\n") {
		f := strings.Fields(baris)
		if len(f) >= 2 && f[0] != "CONSTRAINT" {
			k[f[0]] = true
		}
	}
	return k
}

// TestKolomAnekaAdaDiMigrasi - setiap kolom yang ditulis / dibaca ada di 214 (dua arah untuk empat anak: semua kolom data anak
// dipakai). Instrumen: T_ANEKALIST 50 kolom, T_MEMORANDUMACCEPT 20.
func TestKolomAnekaAdaDiMigrasi(t *testing.T) {
	a := kolomMigrasi214(t, "T_ANEKALIST")
	if len(a) != 50 || len(kolomMigrasi214(t, "T_MEMORANDUMACCEPT")) != 20 {
		t.Fatalf("instrumen: %d kolom T_ANEKALIST", len(a))
	}
	for _, k := range kolomAnekaData {
		if !a[k.nama] {
			t.Errorf("T_ANEKALIST.%s tidak ada di 214", k.nama)
		}
	}
	for _, c := range []struct {
		tabel string
		nama  []string
	}{
		{"T_MARINEHULL", namaKolom(kolomMarineHull)}, {"T_AVIATIONHULL", namaKolom(kolomAviationHull)},
		{"T_VEHICLEHE", namaKolom(kolomVehicleHE)}, {"T_MEMORANDUMACCEPT", namaKolom(kolomMemorandum)},
	} {
		m := kolomMigrasi214(t, c.tabel)
		for _, n := range c.nama {
			if !m[n] {
				t.Errorf("%s.%s tidak ada di 214", c.tabel, n)
			}
			delete(m, n)
		}
		for _, sistem := range []string{"ID", "IDPEGA", "COB_GROUP", "PARENT_ID"} {
			delete(m, sistem)
		}
		if len(m) != 0 {
			t.Errorf("%s: kolom 214 tidak dipakai %v", c.tabel, m)
		}
	}
}

func namaKolom[T any](k []kolomData[T]) []string {
	var n []string
	for _, x := range k {
		n = append(n, x.nama)
	}
	return n
}

// TestSQLAneka - bind 1..n tanpa ganda / celah; saringan induk; kolom sisip = placeholder; mata uang opsional.
func TestSQLAneka(t *testing.T) {
	ta := tabelAnekaUji()
	if n := bindTerbesar(t, sqlAnekaKasus(ta)); n != 2 {
		t.Errorf("sqlAnekaKasus %d bind", n)
	}
	for _, q := range []string{sqlAnekaKasus(ta), sqlKunciAnekaLama(ta)} {
		for _, harus := range []string{"a.PARENT_TABLE = 'T_OCCUPATIONLIST'", "a.PARENT_TABLE = 'T_RISKLOCATION'"} {
			if !strings.Contains(q, harus) {
				t.Errorf("tanpa %s: %s", harus, q)
			}
		}
	}
	if bindTerbesar(t, sqlKunciAnekaLama(ta)) != 2 || !strings.Contains(sqlKunciAnekaLama(ta), "c.PARENT_TABLE = 'T_ANEKALIST'") {
		t.Error("sqlKunciAnekaLama")
	}
	for _, h := range sqlHapusAneka(ta) {
		if n := bindTerbesar(t, h.q); n != h.bind {
			t.Errorf("%d bind, dicatat %d: %s", n, h.bind, h.q)
		}
	}
	if h := sqlHapusAneka(ta); len(h) != 6 || !strings.HasPrefix(h[4].q, "DELETE FROM S.ANEKA") ||
		!strings.Contains(h[5].q, "o.PARENT_TABLE = 'T_RISKLOCATION'") {
		t.Errorf("urutan hapus: %v", h)
	}
	for _, q := range append(sqlHapusCoverageAneka(ta), sqlPindahCoverage(ta.cov)) {
		if !strings.Contains(q, "PARENT_TABLE = 'T_ANEKALIST'") {
			t.Errorf("coverage tanpa saringan induk Aneka: %s", q)
		}
	}
	for _, mu := range []bool{true, false} {
		q := sqlSisipAneka(ta.aneka, mu)
		kolom := strings.Split(q[strings.Index(q, "(")+1:strings.Index(q, ")")], ", ")
		if n := bindTerbesar(t, q); n != len(kolom) || strings.Contains(q, "CURRENCY_CODE") != mu {
			t.Errorf("mata uang %v: %d bind / %d kolom: %s", mu, n, len(kolom), q)
		}
		if !strings.Contains(q, "TO_NUMBER(:") {
			t.Errorf("uang tanpa angkaMasuk: %s", q)
		}
	}
	if q := sqlBacaAneka(ta); bindTerbesar(t, q) != 2 || !strings.Contains(q, angkaKeluar("a.TSI")) {
		t.Errorf("baca: %s", q)
	}
}

func anekaUji(nama, mataUang string) models.ObjekAneka {
	tsi, _ := utils.ParseDecimal("1500.5")
	return models.ObjekAneka{ObjectName: nama, Currency: mataUang, TSI: tsi}
}

// TestSisipAnekaBerkelompok - traversal: induk, PARENT_TABLE / SRC_PATH, SEQ_NO, anak hanya bila terisi, mata uang kosong tanpa
// CURRENCY_CODE, peta posisi -> ID objek baru, jumlah bind = placeholder.
func TestSisipAnekaBerkelompok(t *testing.T) {
	ta := tabelAnekaUji()
	kapal := anekaUji("UJI KAPAL", "IDR")
	kapal.MarineHull.TypeOfVessel = "TUG"
	kapal.Memorandum.SuretyType = "SB-BB"
	pesawat := anekaUji("UJI PESAWAT", "")
	pesawat.AviationHull.Pilots = "2"
	lokasi := []models.LokasiAneka{
		{ObjectNo: "1", IsTopRisk: true, Okupasi: []models.OkupasiAneka{
			{OccupationID: "OK1", Objek: []models.ObjekAneka{kapal, anekaUji("UJI DUA", "USD")}},
			{OccupationID: "OK2"}}, ObjekGolf: []models.ObjekAneka{pesawat}},
		{ObjectNo: "2"},
	}
	e := &eksekutorTiruan{u: &urutTiruan{n: map[string]int{}}}
	baru, ganti, err := sisipAnekaBerkelompok(context.Background(), e, ta, "K1", lokasi, uidUrut())
	if err != nil {
		t.Fatal(err)
	}
	perTabel := map[string][]barisRekam{}
	for _, b := range e.baris {
		if n := bindTerbesar(t, b.q); n != len(b.arg) {
			t.Errorf("%d bind, %d argumen: %s", n, len(b.arg), b.q)
		}
		tabel := strings.Fields(b.q)[2]
		perTabel[tabel] = append(perTabel[tabel], b)
	}
	jumlah := map[string]int{"S.LOC": 2, "S.PROP": 2, "S.RISK": 2, "S.OKU": 2, "S.ANEKA": 3, "S.MAR": 1, "S.AVI": 1, "S.MEMO": 1}
	for tabel, n := range jumlah {
		if len(perTabel[tabel]) != n {
			t.Errorf("%s: %d baris, mau %d", tabel, len(perTabel[tabel]), n)
		}
	}
	if len(perTabel["S.HE"]) != 0 || len(perTabel["S.BANG"]) != 0 || len(perTabel["S.ITEM"]) != 0 {
		t.Error("tabel yang tidak diisi tertulis")
	}
	oku := perTabel["S.OKU"]
	if oku[0].arg[1] != perTabel["S.RISK"][0].arg[0] || oku[0].arg[2] != "T_RISKLOCATION" ||
		oku[0].arg[3] != "LocationList/Property/RiskLocation/OccupationList" || oku[1].arg[4] != 2 {
		t.Errorf("okupasi Aneka: %v", oku)
	}
	an := perTabel["S.ANEKA"]
	if an[0].arg[1] != oku[0].arg[0] || an[0].arg[2] != "T_OCCUPATIONLIST" || an[1].arg[4] != 2 || an[0].arg[6] != "IDR" {
		t.Errorf("objek okupasi: %v", an[:2])
	}
	if an[2].arg[1] != perTabel["S.RISK"][0].arg[0] || an[2].arg[2] != "T_RISKLOCATION" ||
		an[2].arg[3] != "LocationList/Property/RiskLocation/AnekaList" || strings.Contains(an[2].q, "CURRENCY_CODE") {
		t.Errorf("objek Golf: %v", an[2])
	}
	if perTabel["S.MAR"][0].arg[1] != an[0].arg[0] || perTabel["S.AVI"][0].arg[1] != an[2].arg[0] ||
		perTabel["S.MEMO"][0].arg[2] != "SB-BB" {
		t.Error("anak objek")
	}
	if prop := perTabel["S.PROP"][0]; prop.arg[6] != "true" || prop.arg[2] != "1" {
		t.Errorf("property: %v", prop.arg)
	}
	mau := map[string]any{"1/1/1": an[0].arg[0], "1/1/2": an[1].arg[0], "1/0/1": an[2].arg[0]}
	if len(ganti) != 0 {
		t.Errorf("coverage tidak dikirim tetapi ditandai ganti: %v", ganti)
	}
	if len(baru) != 3 {
		t.Errorf("posisi baru %v", baru)
	}
	for k, v := range mau {
		if baru[k] != v {
			t.Errorf("posisi %s = %v, mau %v", k, baru[k], v)
		}
	}
	// TSI uang lewat ikatDesimal (teks bertitik).
	if !strings.Contains(strings.Join(func() []string {
		var s []string
		for _, a := range an[0].arg {
			if x, ok := a.(string); ok {
				s = append(s, x)
			}
		}
		return s
	}(), "|"), "1500.5") {
		t.Errorf("TSI: %v", an[0].arg)
	}
}

// TestRencanaCoverageAneka - posisi yang tetap -> pindah (lama -> baru, bukan terbalik); posisi hilang -> hapus; objek baru
// tanpa coverage lama diabaikan.
func TestRencanaCoverageAneka(t *testing.T) {
	pindah, hapus := rencanaCoverageAneka(map[string]string{"1/1/1": "10", "1/0/2": "11", "2/1/1": "12"},
		map[string]string{"1/1/1": "90", "1/0/1": "91", "2/1/1": "92", "3/1/1": "93"}, nil)
	if !reflect.DeepEqual(pindah, []pindahCoverage{{lama: "10", baru: "90"}, {lama: "12", baru: "92"}}) ||
		!reflect.DeepEqual(hapus, []string{"11"}) {
		t.Errorf("%v / %v", pindah, hapus)
	}
	// P2: coverage posisi 2/1/1 dikirim (diganti utuh, sudah disisip baru) -> coverage lama dihapus, tidak dipindah.
	pindah, hapus = rencanaCoverageAneka(map[string]string{"1/1/1": "10", "2/1/1": "12"},
		map[string]string{"1/1/1": "90", "2/1/1": "92"}, map[string]bool{"2/1/1": true, "3/1/1": true})
	if !reflect.DeepEqual(pindah, []pindahCoverage{{lama: "10", baru: "90"}}) || !reflect.DeepEqual(hapus, []string{"12"}) {
		t.Errorf("ganti: %v / %v", pindah, hapus)
	}
	if !strings.HasPrefix(sqlPindahCoverage("S.COV"), "UPDATE S.COV SET PARENT_ID = :1 ") ||
		!strings.HasSuffix(sqlPindahCoverage("S.COV"), "AND PARENT_ID = :2") {
		t.Errorf("bind pindah: %s", sqlPindahCoverage("S.COV"))
	}
	if q := sqlHapusCoverageAneka(tabelAnekaUji()); !strings.HasPrefix(q[0], "DELETE FROM S.SPR") ||
		!strings.HasPrefix(q[1], "DELETE FROM S.DED") || !strings.HasPrefix(q[2], "DELETE FROM S.COV") {
		t.Errorf("urutan hapus coverage (FK): %v", q)
	}
}

// TestMataUangBacaAneka - DEFAULT 'UNKNOWN' (objek tanpa Currency, mis. Aviation Hull) dibaca kosong; kode lain apa adanya.
func TestMataUangBacaAneka(t *testing.T) {
	for masuk, mau := range map[string]string{"UNKNOWN": "", "IDR": "IDR", "": "", "unknown": "unknown"} {
		if got := mataUangBaca(masuk); got != mau {
			t.Errorf("%q -> %q, mau %q", masuk, got, mau)
		}
	}
}

// TestSisipCoverageAneka - P2: coverage dikirim (larik) -> disisip berinduk objek baru, PARENT_TABLE T_ANEKALIST, SRC_PATH per jalur,
// SEQ_NO, mata uang opsional; deductible berinduk coverage; larik kosong tetap menandai `ganti`; nil tidak.
func TestSisipCoverageAneka(t *testing.T) {
	ta := tabelAnekaUji()
	tsi, _ := utils.ParseDecimal("100")
	cov := []models.CoverageAneka{{Coverage: "C1", Currency: "IDR", TSI: tsi, CalculateMethodFacIn: "1",
		Deductibles: []models.Deductible{{TypeDeductible: "1"}, {TypeDeductible: "2", Currency: "USD"}}}, {Coverage: "C2"}}
	lokasi := []models.LokasiAneka{{Okupasi: []models.OkupasiAneka{{Objek: []models.ObjekAneka{
		{ObjectName: "A", Coverages: cov}, {ObjectName: "B"}, {ObjectName: "C", Coverages: []models.CoverageAneka{}}}}},
		ObjekGolf: []models.ObjekAneka{{ObjectName: "G", Coverages: []models.CoverageAneka{{Coverage: "CG"}}}}}}
	e := &eksekutorTiruan{u: &urutTiruan{n: map[string]int{}}}
	baru, ganti, err := sisipAnekaBerkelompok(context.Background(), e, ta, "K1", lokasi, uidUrut())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ganti, map[string]bool{"1/1/1": true, "1/1/3": true, "1/0/1": true}) {
		t.Errorf("ganti: %v", ganti)
	}
	var covBaris, ded []barisRekam
	for _, b := range e.baris {
		if n := bindTerbesar(t, b.q); n != len(b.arg) {
			t.Errorf("%d bind, %d argumen: %s", n, len(b.arg), b.q)
		}
		switch strings.Fields(b.q)[2] {
		case "S.COV":
			covBaris = append(covBaris, b)
		case "S.DED":
			ded = append(ded, b)
		}
	}
	if len(covBaris) != 3 || len(ded) != 2 {
		t.Fatalf("coverage %d, deductible %d", len(covBaris), len(ded))
	}
	jOk, jGolf := "LocationList/Property/RiskLocation/OccupationList/AnekaList/CoverageList",
		"LocationList/Property/RiskLocation/AnekaList/CoverageList"
	if c := covBaris[0]; c.arg[1] != baru["1/1/1"] || c.arg[2] != "T_ANEKALIST" || c.arg[3] != jOk || c.arg[4] != 1 ||
		c.arg[6] != "IDR" || !strings.Contains(c.q, "CURRENCY_CODE") {
		t.Errorf("coverage 1: %v", c.arg[:7])
	}
	if c := covBaris[1]; c.arg[4] != 2 || strings.Contains(c.q, "CURRENCY_CODE") {
		t.Errorf("coverage 2 tanpa mata uang: %v / %s", c.arg[:6], c.q)
	}
	if c := covBaris[2]; c.arg[1] != baru["1/0/1"] || c.arg[3] != jGolf {
		t.Errorf("coverage Golf: %v", c.arg[:6])
	}
	if d := ded[1]; d.arg[1] != covBaris[0].arg[0] || d.arg[2] != "T_COVERAGELIST" || d.arg[3] != jOk+"/DeductibleList" ||
		d.arg[4] != 2 || !strings.Contains(d.q, "CURRENCY") {
		t.Errorf("deductible: %v", d.arg[:6])
	}
}

// TestKolomCoverageAnekaAdaDiMigrasi - setiap kolom data coverage Aneka ada di T_COVERAGELIST sesudah 193/195/198/203/214.
func TestKolomCoverageAnekaAdaDiMigrasi(t *testing.T) {
	ada := map[string]bool{}
	for _, f := range []string{"193_t_coveragelist.sql", "195_t_coveragelist_unit_akumulasi.sql", "198_spreading.sql",
		"203_coverage_zona.sql", "214_objek_aneka.sql"} {
		b, err := os.ReadFile("../migrations/" + f)
		if err != nil {
			t.Fatal(err)
		}
		teks := strings.ReplaceAll(string(b), "\r", "")
		for _, pola := range []string{`(?s)CREATE TABLE \{skema\}\.T_COVERAGELIST \((.*?)\n\)`,
			`(?s)ALTER TABLE \{skema\}\.T_COVERAGELIST ADD \((.*?)\n\)`} {
			for _, m := range regexp.MustCompile(pola).FindAllStringSubmatch(teks, -1) {
				for _, baris := range strings.Split(m[1], "\n") {
					if f := strings.Fields(baris); len(f) >= 2 && f[0] != "CONSTRAINT" {
						ada[f[0]] = true
					}
				}
			}
		}
	}
	if !ada["PCT_SHORT_PERIOD"] || !ada["TSI_NUSANTARA_RE"] || !ada["UNIT"] {
		t.Fatalf("instrumen: %v", ada)
	}
	for _, k := range kolomCoverageAnekaData {
		if !ada[k.nama] {
			t.Errorf("T_COVERAGELIST.%s tidak dibuat migrasi", k.nama)
		}
	}
	for _, q := range []string{sqlSisipCoverageAneka("S.COV", true), sqlSisipCoverageAneka("S.COV", false)} {
		kolom := strings.Split(q[strings.Index(q, "(")+1:strings.Index(q, ")")], ", ")
		if bindTerbesar(t, q) != len(kolom) {
			t.Errorf("bind != kolom: %s", q)
		}
	}
	ta := tabelAnekaUji()
	for _, q := range []string{sqlBacaCoverageAneka(ta), sqlBacaDeductibleAneka(ta)} {
		if bindTerbesar(t, q) != 2 || !strings.Contains(q, "c.PARENT_TABLE = 'T_ANEKALIST'") {
			t.Errorf("baca: %s", q)
		}
	}
}

// TestSQLProRateKasus - migrasi 255: tulis lewat angkaMasuk (:1 nilai, :2 ID), baca lewat angkaKeluar.
func TestSQLProRateKasus(t *testing.T) {
	if q := sqlUbahProRateKasus("S.G"); q != "UPDATE S.G SET PRO_RATE_PERCENT = "+angkaMasuk(":1")+" WHERE ID = :2" || bindTerbesar(t, q) != 2 {
		t.Errorf("%q", q)
	}
	if q := sqlProRateKasus("S.G"); q != "SELECT "+angkaKeluar("PRO_RATE_PERCENT")+" FROM S.G WHERE ID = :1" {
		t.Errorf("%q", q)
	}
	if q := sqlPercentShareKasus("S.G"); q != "SELECT "+angkaKeluar("PERCENT_SHARE")+" FROM S.G WHERE ID = :1" {
		t.Errorf("%q", q)
	}
}

// TestSisipSpreadingAneka - tiket 73 P3: `.SpreadingList` coverage Aneka berinduk coverage (PARENT_TABLE T_COVERAGELIST), SRC_PATH
// jalur coverage + "/SpreadingList", SEQ_NO 1..n, CURRENCY_CODE = mata uang coverage -> objek -> "UNKNOWN"; ID dari lintasan hitung.
func TestSisipSpreadingAneka(t *testing.T) {
	ta := tabelAnekaUji()
	share, _ := utils.ParseDecimal("40")
	baris := []models.BarisSpreading{{TreatyType: "OR", TreatyName: "UJI OR", SharePercentage: share}, {TreatyType: "QS"}}
	lokasi := []models.LokasiAneka{{Okupasi: []models.OkupasiAneka{{Objek: []models.ObjekAneka{
		{Currency: "USD", Coverages: []models.CoverageAneka{{Currency: "IDR", Spreading: baris}, {Spreading: baris[:1]}}},
		{Coverages: []models.CoverageAneka{{Spreading: baris[1:]}}}}}}}}
	e := &eksekutorTiruan{u: &urutTiruan{n: map[string]int{}}}
	if _, _, err := sisipAnekaBerkelompok(context.Background(), e, ta, "K1", lokasi, uidUrut()); err != nil {
		t.Fatal(err)
	}
	var cov, spr []barisRekam
	for _, b := range e.baris {
		switch strings.Fields(b.q)[2] {
		case "S.COV":
			cov = append(cov, b)
		case "S.SPR":
			if n := bindTerbesar(t, b.q); n != len(b.arg) {
				t.Errorf("%d bind, %d argumen", n, len(b.arg))
			}
			spr = append(spr, b)
		}
	}
	if len(cov) != 3 || len(spr) != 4 {
		t.Fatalf("coverage %d, spreading %d", len(cov), len(spr))
	}
	jalur := "LocationList/Property/RiskLocation/OccupationList/AnekaList/CoverageList/SpreadingList"
	for i, mau := range []struct {
		induk    any
		seq      int
		jenis    any
		mataUang string
	}{{cov[0].arg[0], 1, "OR", "IDR"}, {cov[0].arg[0], 2, "QS", "IDR"}, {cov[1].arg[0], 1, "OR", "USD"}, {cov[2].arg[0], 1, "QS", mataUangKosong}} {
		a := spr[i].arg
		if a[1] != mau.induk || a[2] != "T_COVERAGELIST" || a[3] != jalur || a[4] != mau.seq || a[6] != mau.jenis || a[13] != mau.mataUang {
			t.Errorf("spreading %d: %v", i, a)
		}
	}
	if spr[0].arg[8] != "40" || spr[0].arg[7] != "UJI OR" || spr[1].arg[7] != nil || spr[1].arg[8] != nil {
		t.Errorf("isi baris: %v / %v", spr[0].arg[6:9], spr[1].arg[6:9])
	}
	if q := sqlBacaSpreadingAneka(ta); bindTerbesar(t, q) != 2 || !strings.Contains(q, "FROM S.SPR s WHERE s.PARENT_ID IN (") ||
		!strings.HasSuffix(q, "ORDER BY s.PARENT_ID, s.SEQ_NO") {
		t.Errorf("baca: %s", q)
	}
}
