package repository

// Uji katalog (lanjutan 1, L1/L3): kolom yang DITULIS repository ke kedua tabel produk ada di katalog DEV
// (`testdata/katalog-dev.json`, `ALL_TAB_COLUMNS` 01-10-2026); lebar kolom datar = katalog; objek master
// yang dibaca pemilih = katalog `ALL_OBJECTS`. Kolom yang tidak ada di DEV = `ORA-00904` pada setiap simpan.
//
// ⚠️ Batas cakupannya, dinyatakan: katalog yang diberikan brief hanya memuat kedua tabel produk dan NAMA
// objek master. Penulis `M_ATTACHMENTPRODUCTNAME`, `T_STORAGE_IMAGE`, `T_LOG_SERVICE_RNM` dan kolom yang
// dibaca pemilih / `BrowseReinstypeOR_SQL` belum dicocokkan katalog DEV (sumbernya dokumen DBA) - sisa
// risiko yang dilaporkan, bukan yang dijaga uji ini.

import (
	"encoding/json"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// katalog - isi `testdata/katalog-dev.json`.
type katalog struct {
	Tabel       map[string][]string `json:"tabel"`
	Tipe        map[string]string   `json:"tipe"`
	ObjekMaster map[string]string   `json:"objekMaster"`
}

func katalogDEV(t *testing.T) katalog {
	t.Helper()
	isi, err := os.ReadFile("testdata/katalog-dev.json")
	if err != nil {
		t.Fatal(err)
	}
	var k katalog
	if err := json.Unmarshal(isi, &k); err != nil {
		t.Fatal(err)
	}
	return k
}

var (
	polaSisip  = regexp.MustCompile(`(?is)INSERT\s+INTO\s+\S+\s*\(([^)]*)\)`)
	polaUbah   = regexp.MustCompile(`(?is)UPDATE\s+\S+\s+SET\s+(.*?)\s+WHERE\b`)
	polaSetKol = regexp.MustCompile(`(?:^|,)\s*("?[A-Za-z_][A-Za-z0-9_$#]*"?)\s*=`)
	polaLebar  = regexp.MustCompile(`^VARCHAR2\((\d+)\)$`)
)

// namaKolom - pengenal Oracle: tanpa kutip = huruf besar; berkutip = apa adanya.
func namaKolom(k string) string {
	k = strings.TrimSpace(k)
	if strings.HasPrefix(k, `"`) && strings.HasSuffix(k, `"`) {
		return strings.Trim(k, `"`)
	}
	return strings.ToUpper(k)
}

// kolomDitulis - kolom yang disebut INSERT (daftar kolom) atau UPDATE (klausa SET).
func kolomDitulis(q string) []string {
	var hasil []string
	if m := polaSisip.FindStringSubmatch(q); m != nil {
		for _, k := range strings.Split(m[1], ",") {
			hasil = append(hasil, namaKolom(k))
		}
	}
	if m := polaUbah.FindStringSubmatch(q); m != nil {
		for _, s := range polaSetKol.FindAllStringSubmatch(m[1], -1) {
			hasil = append(hasil, namaKolom(s[1]))
		}
	}
	return hasil
}

// kolomAsing - kolom yang ditulis q tetapi tidak ada di katalog tabelnya.
func kolomAsing(q string, katalog []string) []string {
	ada := map[string]bool{}
	for _, k := range katalog {
		ada[k] = true
	}
	var asing []string
	for _, k := range kolomDitulis(q) {
		if !ada[k] {
			asing = append(asing, k)
		}
	}
	return asing
}

func TestKolomDitulisAdaDiKatalogDEV(t *testing.T) {
	kat := katalogDEV(t).Tabel
	kasus := []struct {
		tabel string
		sql   string
	}{
		{TabelProduk, sqlSisipUmum("S.T")},
		{TabelProduk, sqlPerbaruiUmum("S.T")},
		{TabelInward, sqlSisipInward("S.T")},
		{TabelInward, sqlPerbaruiInward("S.T")},
	}
	for _, k := range kasus {
		if len(kolomDitulis(k.sql)) == 0 {
			t.Errorf("%s: pembaca kolom tidak menemukan kolom apa pun - instrumennya yang rusak:\n%s", k.tabel, rata(k.sql))
		}
		if asing := kolomAsing(k.sql, kat[k.tabel]); len(asing) > 0 {
			t.Errorf("%s: kolom %v tidak ada di katalog DEV %v (ORA-00904):\n%s", k.tabel, asing, kat[k.tabel], rata(k.sql))
		}
	}
	// Daftar kolom penjaga kata cadangan = katalog persis.
	for tabel, kolom := range map[string][]string{TabelProduk: KolomProduk, TabelInward: KolomInward} {
		a, b := append([]string{}, kolom...), append([]string{}, kat[tabel]...)
		sort.Strings(a)
		sort.Strings(b)
		if strings.Join(a, ",") != strings.Join(b, ",") {
			t.Errorf("%s: daftar kolom %v ≠ katalog DEV %v", tabel, kolom, kat[tabel])
		}
	}
}

// Uji gigit: menambah PRODUCTNAME / BEGIN_DATE ke INSERT / UPDATE = merah - juga huruf kecil dan berkutip.
func TestAturanKatalogMenggigit(t *testing.T) {
	kat := katalogDEV(t).Tabel[TabelProduk]
	for _, tambahan := range []string{"PRODUCTNAME", "productname", `"PRODUCTNAME"`} {
		sisip := strings.Replace(sqlSisipUmum("S.T"), "RIRISK)", "RIRISK, "+tambahan+")", 1)
		if asing := kolomAsing(sisip, kat); len(asing) != 1 || asing[0] != "PRODUCTNAME" {
			t.Errorf("INSERT ber-%s harus tertangkap: %v\n%s", tambahan, asing, rata(sisip))
		}
	}
	for _, tambahan := range []string{"BEGIN_DATE", "begin_date", `"BEGIN_DATE"`} {
		ubah := strings.Replace(sqlPerbaruiUmum("S.T"), " WHERE", ", "+tambahan+" = TO_DATE(:9, 'DD/MM/YYYY') WHERE", 1)
		if asing := kolomAsing(ubah, kat); len(asing) != 1 || asing[0] != "BEGIN_DATE" {
			t.Errorf("UPDATE ber-%s harus tertangkap: %v\n%s", tambahan, asing, rata(ubah))
		}
	}
	// Butir yang jawabannya diketahui - instrumen diuji lebih dulu.
	for q, mau := range map[string]string{
		`UPDATE S.T SET A = :1, B = TO_DATE(:2, 'DD/MM/YYYY') WHERE ID = :3`: "A,B",
		`update s.t set jsondata = :1, "Ririsk" = :2 where id = :3`:          "JSONDATA,Ririsk",
		"INSERT INTO S.T (ID, X)\n VALUES (:1, :2)":                          "ID,X",
		`insert into s.t (id, "x") values (:1, :2)`:                          "ID,x",
	} {
		if got := strings.Join(kolomDitulis(q), ","); got != mau {
			t.Errorf("pembaca kolom %q: %s, mau %s", q, got, mau)
		}
	}
}

// Lebar kolom datar yang diperiksa services = lebar katalog DEV.
func TestLebarKolomDatarSamaDenganKatalogDEV(t *testing.T) {
	tipe := katalogDEV(t).Tipe
	for kolom, lebar := range map[string]int{"RIRISKID": LebarRIRiskID, "RIRISK": LebarRIRisk} {
		m := polaLebar.FindStringSubmatch(tipe[TabelProduk+"."+kolom])
		if m == nil {
			t.Fatalf("%s: katalog tanpa VARCHAR2(n): %q", kolom, tipe[TabelProduk+"."+kolom])
		}
		if n, _ := strconv.Atoi(m[1]); n != lebar {
			t.Errorf("%s: lebar services %d ≠ katalog DEV %d", kolom, lebar, n)
		}
	}
	if !strings.HasPrefix(tipe[TabelProduk+".ID"], "VARCHAR2(6)") {
		t.Errorf("ID produk: %q (FormatIdentitas menulis 6 karakter)", tipe[TabelProduk+".ID"])
	}
}

// L3 (OQ-MPNL-04 ditutup data DEV): setiap pemilih membaca objek bernama PERSIS objek yang ada di DEV -
// tabel `AGENT`, `CLIENT`; view `CURRENCY`, `CAUSEOFLOSS_LIFE`, `PRODUCT_TYPE_LIFE`, `RIRISK_LIFE_SUMMARY`,
// dan sejak K1 01-10-2026 (OQ-MPNL-03) view `RATE_LIFE_SUMMARY` (`Choose R/I Rate`) dan `RATE_LIFE` (`View Rate`).
// Yang dibuktikan: NAMA objek; kolom yang dibaca belum ada di katalog yang diberikan.
func TestObjekMasterAdaDiKatalogDEV(t *testing.T) {
	obj := katalogDEV(t).ObjekMaster
	if len(obj) != 8 {
		t.Fatalf("katalog objek master: %v", obj)
	}
	dibaca := map[string]string{MasterJenisPlan: "PLAN LIST", MasterRate: "View Rate"}
	for jenis, s := range sumberMaster {
		dibaca[s.objek] = string(jenis)
	}
	for objek, oleh := range dibaca {
		if _, ada := obj[objek]; !ada {
			t.Errorf("pemilih %s membaca %s - tidak ada di katalog DEV", oleh, objek)
		}
	}
	for objek := range obj {
		if _, ada := dibaca[objek]; !ada {
			t.Errorf("objek katalog %s tidak dibaca pemilih mana pun - salah satu pemilih membaca nama lain", objek)
		}
	}
}
