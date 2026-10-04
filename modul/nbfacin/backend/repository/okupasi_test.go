package repository

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// TestSQLOkupasi - tiket 40: okupasi hanya berinduk T_PROPERTY (literal tetap), dibaca lewat case urut property lalu
// SEQ_NO, TableOfLimit LEFT JOIN; hapus TableOfLimit sebelum okupasi; PARENT_TABLE / SRC_PATH = isi loader.
func TestSQLOkupasi(t *testing.T) {
	baca := sqlBacaOkupasi(tabelObjekUji())
	for _, harus := range []string{"FROM UJI.O o", "JOIN UJI.P p ON p.ID = o.PARENT_ID", "LEFT JOIN UJI.K k ON k.PARENT_ID = o.ID",
		"WHERE o.PARENT_TABLE = 'T_PROPERTY' AND l.PARENT_ID = :1", "ORDER BY o.PARENT_ID, o.SEQ_NO"} {
		if !strings.Contains(baca, harus) {
			t.Errorf("baca okupasi tanpa %q", harus)
		}
	}
	var hapusK, hapusO string
	for _, q := range sqlHapusObjek(tabelObjekUji()) {
		switch {
		case strings.HasPrefix(q, "DELETE FROM UJI.K "):
			hapusK = q
		case strings.HasPrefix(q, "DELETE FROM UJI.O "):
			hapusO = q
		}
	}
	if !strings.HasPrefix(hapusK, "DELETE FROM UJI.K WHERE PARENT_ID IN (SELECT o.ID FROM UJI.O o WHERE o.PARENT_TABLE = 'T_PROPERTY'") ||
		!strings.HasPrefix(hapusO, "DELETE FROM UJI.O o WHERE o.PARENT_TABLE = 'T_PROPERTY' AND o.PARENT_ID IN (") {
		t.Errorf("hapus okupasi: %q / %q", hapusK, hapusO)
	}
	if indukOkupasi != "T_PROPERTY" || jalurOkupasi != "LocationList/Property/OccupationList" {
		t.Error("PARENT_TABLE / SRC_PATH harus sama dengan isi loader")
	}
	// isi loader: jalur rancangan yang sama ada di skema_gen.go.
	b, err := os.ReadFile("../services/loader/skema_gen.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `{jalur: "`+jalurOkupasi+`", tabel: "T_OCCUPATIONLIST", induk: "T_PROPERTY"}`) {
		t.Error("jalur okupasi tidak ada di rancangan loader")
	}
}

// TestBacaOkupasiMemakaiKolomBernama - setiap kunci v("...") di okupasi.go ada di kolomBacaOkupasi dan sebaliknya.
func TestBacaOkupasiMemakaiKolomBernama(t *testing.T) {
	b, err := os.ReadFile("okupasi.go")
	if err != nil {
		t.Fatal(err)
	}
	ada, dibaca := map[string]bool{}, map[string]bool{}
	for _, k := range kolomBacaOkupasi {
		ada[k] = true
	}
	for _, m := range regexp.MustCompile(`v\("([^"]+)"\)`).FindAllStringSubmatch(string(b), -1) {
		dibaca[m[1]] = true
		if !ada[m[1]] {
			t.Errorf("kunci %q dibaca tetapi tidak dipilih", m[1])
		}
	}
	for _, k := range kolomBacaOkupasi {
		if !dibaca[k] {
			t.Errorf("kolom %s dipilih tetapi tidak dibaca", k)
		}
	}
}
