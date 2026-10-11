package repository

// Uji teks SQL polis tabel datar (perintah work owner 10-10-2026, `polis_datar.go`): baca saja, bind urut kemunculan
// (Oracle), hanya kolom yang ada di katalog, kolom NUMBER lewat TO_CHAR 'TM9', tabel ber-SRC_PATH disaring PARENT_TABLE.

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/claimfacin/backend/models"
)

// katalogUji - katalog tiruan: setiap kolom rancangan + kolom sistem, NUMBER untuk `angka`.
func katalogUji(t string, angka ...string) map[string]string {
	spec, _ := models.TabelPolisDatarDari(t)
	m := map[string]string{"ID": "NUMBER", "PARENT_ID": "NUMBER"}
	if spec.Daftar {
		m["SEQ_NO"] = "NUMBER"
	}
	if spec.Src {
		m["SRC_PATH"], m["PARENT_TABLE"] = "VARCHAR2", "VARCHAR2"
	}
	for _, k := range spec.Kolom {
		m[k.Nama] = "VARCHAR2"
	}
	for _, a := range angka {
		m[a] = "NUMBER"
	}
	return m
}

func TestSQLPolisDatarBacaSaja(t *testing.T) {
	akar, _ := models.TabelPolisDatarDari(models.TabelAkarPolisDatar)
	kat := katalogUji(models.TabelAkarPolisDatar)
	semua := []string{sqlKatalogDatar(3), sqlKasusPolisDatar("JP"), sqlKasusProduksiDatar("FP"),
		sqlAkarPolisDatar("G", kolomDatar(akar, kat), kat), sqlOpportunityDatar("O", "A"), sqlBisnisProduksiDatar("B", "FP"),
		sqlBisnisNamaDatar("B"), sqlMarketingDatar("M"), sqlIDMataUangDatar("C")}
	for _, s := range models.SisiPolisDatarSemua() {
		spec, _ := models.TabelPolisDatarDari(s.Tabel)
		k := katalogUji(s.Tabel)
		semua = append(semua, sqlAnakPolisDatar(s.Tabel, kolomDatar(spec, k), k, 3))
	}
	for _, q := range semua {
		u := strings.ToUpper(strings.TrimSpace(q))
		if !strings.HasPrefix(u, "SELECT ") || strings.Contains(u, "COMMIT") || strings.Contains(q, "%[") ||
			strings.Contains(u, "INSERT ") || strings.Contains(u, "UPDATE ") || strings.Contains(u, "DELETE ") {
			t.Errorf("SQL bukan baca saja: %s", q)
		}
	}
}

func TestSQLAnakPolisDatarMenurutKatalog(t *testing.T) {
	cov, _ := models.TabelPolisDatarDari("T_COVERAGELIST")
	// katalog hanya tiga kolom bermedan: kolom rancangan lain tidak dipilih
	kat := map[string]string{"ID": "NUMBER", "PARENT_ID": "NUMBER", "SEQ_NO": "NUMBER", "SRC_PATH": "VARCHAR2",
		"PARENT_TABLE": "VARCHAR2", "COVERAGE": "VARCHAR2", "OLDID": "VARCHAR2", "TSI": "NUMBER"}
	kat[indeksIndukDatar] = "UJI_IX_PARENT"
	q := sqlAnakPolisDatar("UJI.T_COVERAGELIST", kolomDatar(cov, kat), kat, 2)
	mau := "SELECT /*+ INDEX(x UJI_IX_PARENT) */ ID, PARENT_ID, SRC_PATH, COVERAGE, OLDID, " +
		"TO_CHAR(TSI, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''') " +
		"FROM UJI.T_COVERAGELIST x WHERE PARENT_ID IN (:1, :2) AND PARENT_TABLE = :3 ORDER BY SEQ_NO, ID"
	if q != mau {
		t.Fatalf("SQL coverage:\n%s\nmau\n%s", q, mau)
	}
	// halaman tunggal tanpa SEQ_NO / SRC_PATH di katalog
	prop, _ := models.TabelPolisDatarDari("T_PROPERTY")
	kp := map[string]string{"ID": "NUMBER", "PARENT_ID": "NUMBER", "OBJECT_NAME": "VARCHAR2"}
	if q := sqlAnakPolisDatar("UJI.T_PROPERTY", kolomDatar(prop, kp), kp, 1); q !=
		"SELECT ID, PARENT_ID, OBJECT_NAME FROM UJI.T_PROPERTY x WHERE PARENT_ID IN (:1) ORDER BY ID" {
		t.Fatalf("SQL property: %s", q)
	}
}

func TestNilaiDatarSentinelMataUang(t *testing.T) {
	ks := []models.KolomPolisDatar{{Nama: "CURRENCY", Medan: "Currency"}, {Nama: "CURRENCY_CODE", Medan: "Currency.Name"},
		{Nama: "ITEM_TYPE", Medan: "ItemType"}}
	m := nilaiDatar(ks, []string{"UNKNOWN", "UNKNOWN", "UNKNOWN"})
	if len(m) != 1 || m["ItemType"] != "UNKNOWN" {
		t.Fatalf("sentinel UNKNOWN hanya untuk kolom mata uang: %v", m)
	}
}

func TestPilihIndeksInduk(t *testing.T) {
	// dua kolom terdepan (urut TABLE_NAME, INDEX_NAME, COLUMN_POSITION): PARENT_ID, atau PARENT_TABLE lalu PARENT_ID
	got := pilihIndeksInduk([][]string{
		{"T_A", "IX_A_LAIN", "SEQ_NO", "1"}, {"T_A", "IX_A_PARENT", "PARENT_ID", "1"}, {"T_A", "IX_A_PARENT", "SEQ_NO", "2"},
		{"T_B", "IX_B", "PARENT_TABLE", "1"}, {"T_B", "IX_B", "PARENT_ID", "2"},
		{"T_C", "IX_C", "PARENT_TABLE", "1"}, {"T_C", "IX_C", "SEQ_NO", "2"},
	})
	if len(got) != 2 || got["T_A"] != "IX_A_PARENT" || got["T_B"] != "IX_B" {
		t.Fatalf("indeks induk: %v", got)
	}
}

func TestTembolokDatar(t *testing.T) {
	var c tembolokDatar
	t0 := time.Date(2026, 10, 10, 9, 0, 0, 0, time.UTC)
	c.simpan("UJI-1", []byte("a"), t0)
	if d, ada := c.ambil("UJI-1", t0.Add(umurTembolokDatar)); !ada || string(d) != "a" {
		t.Fatal("entri dalam umur harus ada")
	}
	if _, ada := c.ambil("UJI-1", t0.Add(umurTembolokDatar+time.Second)); ada {
		t.Fatal("entri kedaluwarsa harus hilang")
	}
	for i := 0; i < batasTembolokDatar+5; i++ {
		c.simpan(fmt.Sprintf("UJI-%d", i), []byte("x"), t0.Add(time.Duration(i)*time.Millisecond))
	}
	if len(c.isi) != batasTembolokDatar {
		t.Fatalf("tembolok %d entri, batas %d", len(c.isi), batasTembolokDatar)
	}
	if _, ada := c.isi["UJI-0"]; ada {
		t.Fatal("entri tertua harus dibuang saat penuh")
	}
}

func TestTabelTampilanNBTanpaIndeksIndukTidakDibaca(t *testing.T) {
	// DEV 10-10-2026: T_LISTCAUSEOFLOSS ~1 juta baris tanpa indeks PARENT_ID = ~10 detik per kueri. Tabel khusus tampilan
	// polis NB tanpa indeks induk tidak dikueri sama sekali (Acuan tanpa basis data: kueri apa pun akan panik).
	a := &Acuan{}
	for _, tb := range []string{"T_LISTCAUSEOFLOSS", "T_FEALIST", "T_DEDUCTIBLELIST"} {
		if !models.TabelTampilanNB(tb) {
			t.Fatalf("%s bukan tabel tampilan NB", tb)
		}
		kat := map[string]map[string]string{tb: {"PARENT_ID": "VARCHAR2"}}
		rows, err := a.anakPolisDatar(context.Background(), kat, models.SisiPolisDatar{Tabel: tb, Induk: "T_PROPERTY"}, []string{"1"})
		if err != nil || rows != nil {
			t.Fatalf("%s: %v %v", tb, rows, err)
		}
	}
	for _, tb := range []string{"T_PROPERTY", "T_COVERAGELIST", "T_SPREADINGLIST"} {
		if models.TabelTampilanNB(tb) {
			t.Fatalf("%s dibaca activity klaim - tidak boleh bergerbang indeks", tb)
		}
	}
}

func TestSQLTanpaKolomKembarTakBeralias(t *testing.T) {
	// DEV 11-10-2026: `SELECT COUNTRY, COUNTRY, ID ... FETCH FIRST` = ORA-00960 (Country / Province Input Register gagal).
	// Kolom yang dipilih dua kali wajib beralias pada salinan keduanya.
	kembar := regexp.MustCompile(`SELECT ([A-Za-z_.]+), ([A-Za-z_.]+)(\s*,|\s+FROM)`)
	nama, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range nama {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range kembar.FindAllStringSubmatch(string(b), -1) {
			if m[1] == m[2] {
				t.Errorf("%s: %q tanpa alias", f, m[0])
			}
		}
	}
}

func TestSQLAgenKlienMendahulukanKlienBeralamat(t *testing.T) {
	// Work owner 11-10-2026 "TOMBOL INI GA FUNGSI" (Choose grid Ceding Co Name): DEV dua agen bernama klien sama, yang
	// terakhir tanpa CLIENT_ADDRESS -> agen yang kliennya beralamat didahulukan.
	q := sqlAgenKlien("UJI_AGENT", "UJI_CA")
	for _, mau := range []string{"FROM UJI_AGENT g WHERE g.CLIENTNAME = :1",
		"EXISTS (SELECT 1 FROM UJI_CA c WHERE c.CLIENTID = g.CLIENTID AND NVL(c.ASMADDRESSTYPE, '-') <> '7')",
		"THEN 0 ELSE 1 END, g.ROWID DESC", "FETCH FIRST 1 ROWS ONLY"} {
		if !strings.Contains(strings.Join(strings.Fields(q), " "), mau) {
			t.Errorf("sqlAgenKlien tanpa %q:\n%s", mau, q)
		}
	}
}
