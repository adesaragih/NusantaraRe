package models

import (
	"errors"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"
)

// AC 24, 33, 34: DUA PULUH LIMA jenis, aturan di KODE, per jenis sendiri.
func TestAturanKlausulDuaPuluhLimaJenis(t *testing.T) {
	jenis := map[string]bool{}
	for _, a := range AturanKlausulTCO {
		jenis[a.Jenis] = true
	}
	if len(jenis) != 25 {
		t.Fatalf("jenis klausul %d, mau 25 (18 induk + 7 anak)", len(jenis))
	}
	anak := 0
	for j := range jenis {
		a, _ := AturanJenisKlausul(j, "")
		if a.Anak {
			anak++
		}
	}
	if anak != 7 {
		t.Errorf("jenis anak %d, mau 7", anak)
	}
	// Setiap DescID 10001-10018 punya induk; tujuh punya anak berDescID SAMA.
	for i := 10001; i <= 10018; i++ {
		id := strconv.Itoa(i)
		if _, ok := CariAturanKlausul(id, false, subjenisBawaan(id)); !ok {
			t.Errorf("DescID %s tanpa aturan induk", id)
		}
	}
	for _, id := range []string{"10001", "10003", "10004", "10006", "10007", "10009", "10012"} {
		if _, ok := CariAturanKlausul(id, true, ""); !ok {
			t.Errorf("DescID %s tanpa aturan anak", id)
		}
	}
}

func subjenisBawaan(id string) string {
	switch id {
	case DescExclutionTreaty:
		return SubjenisOccupation
	case DescCoinsPanel:
		return SubjenisCoinsLessThan
	}
	return ""
}

// Tabel kebenaran wajib-isi per jenis - SESUDAH ralat 29-09-2026 atas tabel
// tiket (prasyarat nonaktif `TerritorialLimit` tidak dihitung).
func TestWajibIsiPerJenis(t *testing.T) {
	mau := map[string]string{
		"TreatyLimit": "ReinsTypeID,Rp,Usd", "PLA": "ReinsTypeID,Rp,Usd", "CashLossLimit": "ReinsTypeID,Rp,Usd",
		"FacIn": "ReinsTypeID,Rp,Usd", "ExGratia": "ReinsTypeID,Rp,Usd", "EPI": "ReinsTypeID,Rp,Usd",
		"ClaimCoorp":       "ReinsTypeID,Rp,Usd",
		"TreatyLimitChild": "ReinsTypeID,Pct,Rp,Usd", "PLAList": "ReinsTypeID,Pct,Rp,Usd",
		"CashLossLimitList": "ReinsTypeID,Pct,Rp,Usd", "FacInList": "ReinsTypeID,Pct,Rp,Usd",
		"ExGratiaChildList": "ReinsTypeID,Pct,Rp,Usd", "EpiList": "ReinsTypeID,Pct,Rp,Usd",
		"ClaimCoorpChild": "ReinsTypeID,Pct,Rp,Usd",
		"ProfitComm":      "Pct,ReinsTypeID,PctMe", "Ricomm": "ReinsTypeID,Pct,Method", "BordereAux": "Method",
		"TerrLimit": "TerritorialLimit", "CoinsPanel": "CoIns_Min,CoIns_Max,TreatyLimit",
		"MaxCoinsPanel": "CoIns_Max", "MinLOL": "Pct", "MinLOLMB": "Pct",
	}
	for jenis, w := range mau {
		a, ok := AturanJenisKlausul(jenis, "")
		if !ok {
			t.Errorf("%s tidak ada", jenis)
			continue
		}
		if strings.Join(a.Wajib, ",") != w {
			t.Errorf("%s wajib %v, mau %s", jenis, a.Wajib, w)
		}
	}
	occ, _ := AturanJenisKlausul("ExclutionTreaty", SubjenisOccupation)
	if strings.Join(occ.Wajib, ",") != "ID_Occupation,Occupation,Line,Usd,Rp" {
		t.Errorf("ExclutionTreaty/Occupation wajib %v", occ.Wajib)
	}
	cl, _ := AturanJenisKlausul("ExclutionTreaty", SubjenisClause)
	if strings.Join(cl.Wajib, ",") != "ID_Clause,Clause" {
		t.Errorf("ExclutionTreaty/Clause wajib %v", cl.Wajib)
	}
}

// AC 36: LimitMB dan Portfolio TIDAK dilepas sebagai "tanpa validasi".
func TestLimitMBDanPortfolioDitahan(t *testing.T) {
	for _, j := range []string{"LimitMB", "Portfolio"} {
		a, _ := AturanJenisKlausul(j, "")
		if a.Ditahan == "" {
			t.Errorf("%s dilepas tanpa aturan wajib-isi", j)
		}
		err := PeriksaKlausulTCO(a, KlausulTreaty{})
		if !errors.Is(err, ErrKlausulDitahan) || !strings.Contains(err.Error(), j) {
			t.Errorf("%s: %v", j, err)
		}
	}
}

// AC 35: pesan menyebut SELURUH medan yang kurang.
func TestPeriksaKlausulMenyebutMedan(t *testing.T) {
	a, _ := AturanJenisKlausul("EPI", "")
	err := PeriksaKlausulTCO(a, KlausulTreaty{ReinsTypeID: "10003"})
	if !errors.Is(err, ErrKlausulMedanWajib) || !strings.Contains(err.Error(), "Rp") || !strings.Contains(err.Error(), "Usd") ||
		strings.Contains(err.Error(), "ReinsTypeID") {
		t.Errorf("pesan: %v", err)
	}
	if err := PeriksaKlausulTCO(a, KlausulTreaty{ReinsTypeID: "10003", Rp: apd.New(1, 0), Usd: apd.New(0, 0)}); err != nil {
		t.Errorf("lengkap ditolak: %v", err)
	}
}

// HitungRpUsd: Usd = Pct x Usd_induk / 100, Rp = Pct x Rp_induk / 100 - persis.
func TestRpUsdAnakTCO(t *testing.T) {
	rp, usd, err := RpUsdAnakTCO(des(t, "33.33333333"), des(t, "1000000000"), des(t, "66666.67"))
	if err != nil {
		t.Fatal(err)
	}
	if rp.Text('f') != "333333333.30000000" || usd.Text('f') != "22222.22333111" {
		t.Errorf("rp %s usd %s", rp.Text('f'), usd.Text('f'))
	}
	// Nilai besar: float64 kehilangan digit di sini (dibuktikan dengan mutasi).
	rp, usd, err = RpUsdAnakTCO(des(t, "12.34567891"), des(t, "98765432109876543.21"), des(t, "1234567.89"))
	if err != nil || rp.Text('f') != "12193263122359396.42211401" || usd.Text('f') != "152415.78762536" {
		t.Errorf("nilai besar: rp %s usd %s %v", rp.Text('f'), usd.Text('f'), err)
	}
	// Induk tanpa Rp: turunannya kosong - gerbang wajib yang menyebutnya.
	rp, _, _ = RpUsdAnakTCO(des(t, "10"), nil, des(t, "1"))
	if rp != nil {
		t.Errorf("rp tanpa induk = %v", rp)
	}
}

// SaveTreatyArr*List langkah 6-7: total Pct anak BARU > 100 ditolak;
// TreatyTestChildTotal_Act: total != 100 hanya PERINGATAN sesudah simpan.
func TestTotalPctAnakTCO(t *testing.T) {
	total, err := PeriksaTotalAnakTCO([]*apd.Decimal{des(t, "60"), nil}, des(t, "40"))
	if err != nil || total.Text('f') != "100" {
		t.Errorf("tepat 100: %v %v", total, err)
	}
	if PeringatanSpreadingTCO(total) != "" {
		t.Error("100 tetap diberi peringatan")
	}
	_, err = PeriksaTotalAnakTCO([]*apd.Decimal{des(t, "60")}, des(t, "40.00000001"))
	if !errors.Is(err, ErrTotalPctAnakMelebihi100) || !strings.Contains(err.Error(), "100.00000001") {
		t.Errorf("> 100: %v", err)
	}
	total, _ = PeriksaTotalAnakTCO(nil, des(t, "50"))
	if PeringatanSpreadingTCO(total) != "Please make sure spreading is 100%" {
		t.Errorf("peringatan: %q", PeringatanSpreadingTCO(total))
	}
}

// Subjenis exclusion diturunkan dari medan yang terisi (tanpa kolom tambahan).
func TestSubjenisExclusionTCO(t *testing.T) {
	kasus := map[string]KlausulTreaty{
		SubjenisOccupation: {IDOccupation: "UJI-O1"},
		SubjenisClause:     {IDClause: "UJI-C1"},
		SubjenisPeriode:    {Layer: "12"},
		SubjenisObject:     {Pct: apd.New(5, 0)},
	}
	for mau, k := range kasus {
		if dapat := SubjenisExclusionTCO(k); dapat != mau {
			t.Errorf("subjenis %q, mau %q", dapat, mau)
		}
	}
}

// Medan yang bukan milik jenis ditolak - tidak ada kolom yang ditulis diam-diam.
func TestIsiMedanKlausulTCO(t *testing.T) {
	a, _ := AturanJenisKlausul("ProfitComm", "")
	var k KlausulTreaty
	if err := IsiMedanKlausulTCO(a, &k, map[string]string{"ReinsTypeID": "10003", "Pct": "12,5", "PctMe": "1",
		"Ydcf": "UJI"}); err != nil {
		t.Fatal(err)
	}
	if k.Pct.Text('f') != "12.5" || k.Ydcf != "UJI" || k.ReinsTypeID != "10003" {
		t.Errorf("isi: %+v", k)
	}
	if err := IsiMedanKlausulTCO(a, &k, map[string]string{"Rp": "1"}); !errors.Is(err, ErrMedanBukanMilikJenis) {
		t.Errorf("medan asing: %v", err)
	}
	if err := IsiMedanKlausulTCO(a, &k, map[string]string{"Pct": "1.000,5"}); !errors.Is(err, ErrPersenBukanDesimal) {
		t.Errorf("desimal salah: %v", err)
	}
}

// TestSemuaAnakMemakaiPilihanAnakTreatyLimit - ketujuh grid `Show Child`
// memilih ReinsType dari daftar anak Treaty Limit (porsi + induknya); induk
// tidak [keputusan work owner 02-10-2026].
func TestSemuaAnakMemakaiPilihanAnakTreatyLimit(t *testing.T) {
	var anak []string
	for _, a := range AturanKlausulTCO {
		mau := ""
		if a.Anak {
			mau = PilihanReinsAnakTreatyLimit
			anak = append(anak, a.Jenis)
		}
		if a.PilihanReins != mau {
			t.Errorf("%s: PilihanReins %q, mau %q", a.Jenis, a.PilihanReins, mau)
		}
	}
	sort.Strings(anak)
	if got := strings.Join(anak, ","); got != "CashLossLimitList,ClaimCoorpChild,EpiList,ExGratiaChildList,FacInList,PLAList,TreatyLimitChild" {
		t.Errorf("aturan anak %s, mau ketujuh grid Show Child", got)
	}
}

// TestSatuBarisHanyaTigaJenis - Minimum LOL, Max Coins Panel, Minimum LOL MB
// (`GridTreatyArrangementMinLOL.xml` b2232, `…MaxCoinsPanel.xml` b2212,
// `…MInLOLMB.xml` b2262).
func TestSatuBarisHanyaTigaJenis(t *testing.T) {
	var satu []string
	for _, a := range AturanKlausulTCO {
		if a.SatuBaris {
			satu = append(satu, a.Jenis+"/"+a.DescID)
		}
	}
	sort.Strings(satu)
	if got := strings.Join(satu, ","); got != "MaxCoinsPanel/10016,MinLOL/10015,MinLOLMB/10018" {
		t.Errorf("jenis satu baris: %s", got)
	}
}

// TestCoinsPanelDuaGridCoInsScale - Co-Ins Scale (`GridTreatyArrangementCoins.xml`) = DUA grid,
// dibedakan `Param.Type` "Less Than" (b10055) / "More Than" (b14320) yang `NewTreatyArrCoins` b406
// simpan di `.SpreadingOrder`; kunci dobel per grid.
func TestCoinsPanelDuaGridCoInsScale(t *testing.T) {
	for _, sub := range []string{SubjenisCoinsLessThan, SubjenisCoinsMoreThan} {
		a, ok := CariAturanKlausul(DescCoinsPanel, false, sub)
		if !ok {
			t.Fatalf("CoinsPanel/%s tidak ada", sub)
		}
		if a.Jenis != "CoinsPanel" || strings.Join(a.Medan, ",") != "CoIns_Min,CoIns_Max,TreatyLimit" ||
			strings.Join(a.Wajib, ",") != "CoIns_Min,CoIns_Max,TreatyLimit" ||
			strings.Join(a.KunciDobel, ",") != "SpreadingOrder,CoIns_Min,CoIns_Max" {
			t.Errorf("CoinsPanel/%s: %+v", sub, a)
		}
	}
	if SubjenisCoinsLessThan != "Less Than" || SubjenisCoinsMoreThan != "More Than" {
		t.Error("subjenis Co-Ins Scale bukan nilai Param.Type VERBATIM")
	}
	if _, ok := CariAturanKlausul(DescCoinsPanel, false, ""); ok {
		t.Error("CoinsPanel tanpa subjenis diterima - grid mana?")
	}
	for k, mau := range map[string]string{"Less Than": "Less Than", "More Than": "More Than", "": "",
		"less than": "", " More Than ": "More Than", "Lainnya": ""} {
		if got := SubjenisKlausulTCO(KlausulTreaty{TreatyDescID: DescCoinsPanel, SpreadingOrder: k}); got != mau {
			t.Errorf("SpreadingOrder %q -> %q, mau %q", k, got, mau)
		}
	}
	if got := SubjenisKlausulTCO(KlausulTreaty{TreatyDescID: DescExclutionTreaty, IDClause: "UJI-1"}); got != SubjenisClause {
		t.Errorf("exclusion: %q", got)
	}
	if got := SubjenisKlausulTCO(KlausulTreaty{TreatyDescID: DescEPI, SpreadingOrder: "Less Than"}); got != "" {
		t.Errorf("jenis lain bersubjenis %q", got)
	}
	k := KlausulTreaty{SpreadingOrder: "More Than"}
	if NilaiMedanKlausul(k, MedanSpreadingOrder) != "More Than" {
		t.Error("SpreadingOrder tidak terbaca sebagai kunci dobel")
	}
}

// TestCoinsPanelKorpus - nilai Type, simpanannya, saringan grid, dan urutannya ada di korpus.
func TestCoinsPanelKorpus(t *testing.T) {
	baris := func(rel string, n int) string {
		b, err := os.ReadFile(akarKorpusTCO + `\` + rel)
		if err != nil {
			t.Skipf("korpus tidak terjangkau (%v)", err)
		}
		return strings.TrimSpace(strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n")[n-1])
	}
	grid := `Section\GridTreatyArrangementCoins.xml`
	rd := `ReportDefinition\BrowseTreatyArrangement_CoinsPanel_RD.xml`
	for _, c := range []struct {
		rel string
		n   int
		mau string
	}{
		{grid, 10055, `<pyValue>"` + SubjenisCoinsLessThan + `"</pyValue>`},
		{grid, 14320, `<pyValue>"` + SubjenisCoinsMoreThan + `"</pyValue>`},
		{`Activity\NewTreatyArrCoins.xml`, 405, `<PropertiesName>InputTreatyCoins.SpreadingOrder</PropertiesName>`},
		{`Activity\NewTreatyArrCoins.xml`, 406, `<PropertiesValue>Param.Type</PropertiesValue>`},
		{rd, 610, `<pyFilterName>.SpreadingOrder</pyFilterName>`},
		{rd, 612, `<pyFilterValue>Param.Type</pyFilterValue>`},
		{rd, 653, `<pyFieldName>.TreatyLimit</pyFieldName>`},
		{rd, 655, `<pySortType>DESC</pySortType>`},
	} {
		if got := baris(c.rel, c.n); got != c.mau {
			t.Errorf("%s b%d: %s, mau %s", c.rel, c.n, got, c.mau)
		}
	}
}

// TestUrutCoInsScaleTCO - `TreatyLimit DESC`, kosong lebih dulu (Oracle NULLS FIRST), seri stabil.
func TestUrutCoInsScaleTCO(t *testing.T) {
	d := func(s string) *apd.Decimal {
		v, _, _ := apd.NewFromString(s)
		return v
	}
	b := []KlausulTreaty{{ID: "1", TreatyLimit: d("500")}, {ID: "2", TreatyLimit: d("1000.5")}, {ID: "3"},
		{ID: "4", TreatyLimit: d("500.00")}, {ID: "5", TreatyLimit: d("99")}}
	UrutCoInsScaleTCO(b)
	var id []string
	for _, k := range b {
		id = append(id, k.ID)
	}
	if got := strings.Join(id, ","); got != "3,2,1,4,5" {
		t.Errorf("urutan %s, mau 3,2,1,4,5", got)
	}
}
