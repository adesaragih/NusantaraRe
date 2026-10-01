package repository

// Bentuk SQL `Save` (`SetPremi_EDM`) - tanpa Oracle; rumus rekap diturunkan
// ulang dari korpus.

import (
	"html"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"

	"nusantarare/modul/endorsementlife/backend/models"
)

const uRekap = "UJISKEMA.T_PREMIUM_LIST_SUMMARY"

// akarKorpus - korpus READ-ONLY; uji yang membacanya DILEWATI bila tak terjangkau.
const akarKorpus = `D:\XML\RNM_BRD\Endorsement Life`

func TestSQLTandaiMembalikHanyaPesertaOld(t *testing.T) {
	q := sqlTandai(uPeserta, 0, 0)
	if n := penampungUnik(t, "sqlTandai", q); n != 3 {
		t.Fatalf("%d penampung, mau 3", n)
	}
	if !strings.Contains(q, "SET d.EDM_STATUS = :1, ") || !strings.HasSuffix(q, "WHERE d.PREMIUM_LIST_ID = :2 AND d.EDM_STATUS = :3") {
		t.Errorf("status dan penyaring Old tidak di satu pernyataan: %s", q)
	}
	for _, k := range models.KolomJurnalBalik {
		if !strings.Contains(q, "d."+k+" = -d."+k) {
			t.Errorf("kolom %s tidak dibalik", k)
		}
	}
	if c := strings.Count(q, " = -d."); c != 32 {
		t.Errorf("%d kolom dibalik, mau 32 (R28: hanya 2.1)", c)
	}
	p := sqlTandai(uPeserta, 2, 0)
	if penampungUnik(t, "sqlTandai pilih", p) != 5 || !strings.HasSuffix(p, "AND d.ID IN (:4, :5)") {
		t.Errorf("pilihan: %s", p)
	}
	k := sqlTandai(uPeserta, 0, 1)
	if penampungUnik(t, "sqlTandai kecuali", k) != 4 || !strings.HasSuffix(k, "AND d.ID NOT IN (:4)") {
		t.Errorf("DELETE ALL berpengecualian: %s", k)
	}
}

func TestSQLRekapBerkunciKasus(t *testing.T) {
	if q := sqlHapusRekap(uRekap); penampungUnik(t, "sqlHapusRekap", q) != 1 || !strings.Contains(q, "WHERE r.PREMIUM_LIST_ID = :1") {
		t.Errorf("hapus rekap: %s", q)
	}
	for _, tipe := range []string{models.TypeQR, models.TypeQP, models.TypeTP, models.TypeTR, ""} {
		q := sqlSisipRekap(uRekap, uPeserta, tipe)
		if n := penampungUnik(t, "sqlSisipRekap "+tipe, q); n != 3 {
			t.Fatalf("%s: %d penampung", tipe, n)
		}
		if !strings.HasSuffix(q, "WHERE d.PREMIUM_LIST_ID = :3 GROUP BY d.CURRENCY") {
			t.Errorf("%s: rekap tidak per mata uang kasus: %s", tipe, q)
		}
	}
	if p, b := rumusPremiBalance("XX"); p != "0" || b != "0" {
		t.Errorf("tipe lain: %s / %s, mau nilai awal 0 (2.2)", p, b)
	}
	if q := sqlRekapKasus(uRekap); penampungUnik(t, "sqlRekapKasus", q) != 1 || !strings.HasSuffix(q, "ORDER BY r.CURRENCY") {
		t.Errorf("baca rekap: %s", q)
	}
}

// TestKolomRekapSamaDenganDDL - setiap kolom rekap ada di DDL 055; kolom
// DDL yang tidak dijumlah hanya identitas dan empat kolom SUM (R08).
func TestKolomRekapSamaDenganDDL(t *testing.T) {
	ddl := map[string]bool{}
	for _, k := range kolomDDL(t, "055_t_premium_list_summary.sql") {
		ddl[k] = true
	}
	ditulis := map[string]bool{"ID": true, "PREMIUM_LIST_ID": true, "CURRENCY": true, "PREMIUM": true, "BALANCE": true}
	for _, k := range kolomJumlahRekap {
		ditulis[k[0]] = true
	}
	for _, k := range models.KolomRekapKasus {
		if !ditulis[k] {
			t.Errorf("kolom layar %s tidak pernah dihitung", k)
		}
	}
	var sisa []string
	for k := range ddl {
		if !ditulis[k] {
			sisa = append(sisa, k)
		}
	}
	for k := range ditulis {
		if !ddl[k] {
			t.Errorf("kolom %s bukan kolom DDL 055", k)
		}
	}
	sort.Strings(sisa)
	if strings.Join(sisa, ",") != "CEDING_RETENTION,SHARE_NUSANTARA_RE_GROSS,SUM_AT_RISK_GROSS,SUM_REASURED" {
		t.Errorf("kolom DDL tak dijumlah: %v", sisa)
	}
	if len(models.KolomRekapKasus) != len(ditulis)-2 {
		t.Errorf("layar membaca %d kolom, rekap menulis %d", len(models.KolomRekapKasus), len(ditulis)-2)
	}
}

// langkahDT - `pyPropertyStepId` → (aksi, nama, nilai) satu DataTransform.
func langkahDT(t *testing.T, jalur string) map[string][3]string {
	t.Helper()
	b, err := os.ReadFile(akarKorpus + `\` + jalur)
	if err != nil {
		t.Skipf("korpus tidak terjangkau (%v)", err)
	}
	ambil := func(potong, tag string) string {
		m := regexp.MustCompile(`<` + tag + `>([^<]*)</` + tag + `>`).FindStringSubmatch(potong)
		if m == nil {
			return ""
		}
		return html.UnescapeString(m[1])
	}
	hasil := map[string][3]string{}
	for _, potong := range strings.Split(string(b), "<rowdata")[1:] {
		if id := ambil(potong, "pyPropertyStepId"); id != "" {
			hasil[id] = [3]string{ambil(potong, "pyActionName"), ambil(potong, "pyPropertiesName"), ambil(potong, "pyPropertiesValue")}
		}
	}
	return hasil
}

// TestRekapDariKorpus - pasangan kolom rekap ← kolom peserta, serta rumus
// `PREMIUM`/`BALANCE` per Type, diturunkan ulang dari `AppendCurrencySummary_DT`.
func TestRekapDariKorpus(t *testing.T) {
	l := langkahDT(t, `DataTransform\AppendCurrencySummary_DT.xml`)
	akum := regexp.MustCompile(`^(?i)param\.(\w+) \+ \.(\w+)$`)
	tujuan := regexp.MustCompile(`^@divide\((?i:param)\.(\w+),1,20\)$`)
	param := map[string]string{} // param (huruf kecil) → kolom peserta
	tipeDi := map[string]string{}
	premi, balance := map[string]string{}, map[string]string{}
	for id, s := range l {
		switch {
		case regexp.MustCompile(`^2\.34\.1\.\d+$`).MatchString(id) && s[0] == "SET":
			m := akum.FindStringSubmatch(s[2])
			if m == nil || !strings.EqualFold(s[1], "Param."+m[1]) {
				t.Errorf("%s %s <= %s bukan akumulasi", id, s[1], s[2])
				continue
			}
			param[strings.ToLower(m[1])] = m[2]
		case regexp.MustCompile(`^2\.34\.1\.\d+$`).MatchString(id) && s[0] == "WHEN":
			tipeDi[id] = regexp.MustCompile(`^pyWorkPage\.Type=="(QR|QP|TP|TR)"$`).FindStringSubmatch(s[1])[1]
		}
	}
	for id, s := range l {
		if induk, ada := tipeDi[strings.TrimSuffix(id, ".1")]; ada && strings.HasSuffix(id, ".1") {
			switch s[1] {
			case "Param.Premium":
				premi[induk] = strings.TrimSpace(strings.TrimPrefix(s[2], "Param.Premium +"))
			case "Param.Balance":
				balance[induk] = strings.TrimSuffix(strings.TrimPrefix(s[2], "Param.Balance+("), ")")
			}
		}
	}
	dariKorpus := map[string]string{}
	for id, s := range l {
		if !regexp.MustCompile(`^2\.\d+$`).MatchString(id) || !strings.HasPrefix(s[1], ".") {
			continue
		}
		m := tujuan.FindStringSubmatch(s[2])
		if m == nil {
			continue
		}
		kolom := strings.TrimPrefix(s[1], ".")
		if kolom == "PREMIUM" || kolom == "BALANCE" {
			continue
		}
		sumber, ada := param[strings.ToLower(m[1])]
		if !ada {
			t.Errorf("%s ← %s: param tanpa akumulasi", kolom, m[1])
		}
		dariKorpus[kolom] = sumber
	}
	if len(dariKorpus) != 30 || len(kolomJumlahRekap) != 30 {
		t.Fatalf("korpus %d pasangan, kode %d; mau 30", len(dariKorpus), len(kolomJumlahRekap))
	}
	for _, k := range kolomJumlahRekap {
		if dariKorpus[k[0]] != k[1] {
			t.Errorf("%s ← %s di kode, korpus %q", k[0], k[1], dariKorpus[k[0]])
		}
	}
	// `NVL(d.X, 0)` → `.X`, tanpa spasi, tanpa `SUM(…)` luar = ekspresi Pega.
	pega := func(sqlnya string) string {
		s := regexp.MustCompile(`NVL\(d\.(\w+), 0\)`).ReplaceAllString(sqlnya, ".$1")
		return strings.TrimSuffix(strings.TrimPrefix(strings.ReplaceAll(s, " ", ""), "SUM("), ")")
	}
	if len(premi) != 4 || len(balance) != 4 {
		t.Fatalf("korpus memberi %d rumus premi, %d balance; mau 4", len(premi), len(balance))
	}
	for tipe := range premi {
		p, b := rumusPremiBalance(tipe)
		if pega(p) != premi[tipe] {
			t.Errorf("%s PREMIUM %q, korpus %q", tipe, pega(p), premi[tipe])
		}
		if pega(b) != strings.ReplaceAll(balance[tipe], " ", "") {
			t.Errorf("%s BALANCE %q, korpus %q", tipe, pega(b), balance[tipe])
		}
	}
}
