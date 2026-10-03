package repository

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

func tabelObjekUji() tabelObjek {
	return tabelObjek{work: "UJI.W", general: "UJI.G", loc: "UJI.L", prop: "UJI.P", risk: "UJI.R", bang: "UJI.B", sekitar: "UJI.S"}
}

// TestSQLObjek - tiket 35: baca urut SEQ_NO lewat case; hapus ANAK sebelum INDUK (FK tanpa
// ON DELETE), tiap pernyataan satu bind :1; General dipastikan ada; bind sisip bernomor.
func TestSQLObjek(t *testing.T) {
	tb := tabelObjekUji()
	baca := sqlBacaObjek(tb)
	for _, harus := range []string{"FROM UJI.L l", "LEFT JOIN UJI.P p ON p.PARENT_ID = l.ID", "LEFT JOIN UJI.R r ON r.PARENT_ID = p.ID",
		"LEFT JOIN UJI.B b ON b.PARENT_ID = p.ID", "LEFT JOIN UJI.S s ON s.PARENT_ID = p.ID", "WHERE l.PARENT_ID = :1", "ORDER BY l.SEQ_NO"} {
		if !strings.Contains(baca, harus) {
			t.Errorf("baca tanpa %q", harus)
		}
	}
	if n := len(strings.Split(baca[len("SELECT "):strings.Index(baca, "\nFROM")], ",")); n != 47 {
		t.Errorf("SELECT %d kolom, mau 47 (23 tiket 35 + 24 tiket 38)", n)
	}
	hapus := sqlHapusObjek(tb)
	urut := []string{"DELETE FROM UJI.R ", "DELETE FROM UJI.B ", "DELETE FROM UJI.S ", "DELETE FROM UJI.P ", "DELETE FROM UJI.L "}
	if len(hapus) != len(urut) {
		t.Fatalf("%d DELETE, mau %d", len(hapus), len(urut))
	}
	for i, q := range hapus {
		if !strings.HasPrefix(q, urut[i]) || strings.Count(q, ":1") != 1 || strings.Contains(q, ":2") {
			t.Errorf("hapus[%d] %q", i, q)
		}
	}
	if got := sqlPastikanGeneral("UJI.G"); got != "INSERT INTO UJI.G (ID) SELECT :1 FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM UJI.G WHERE ID = :2)" {
		t.Errorf("pastikan General %q", got)
	}
	if teksBool(true) != "true" || teksBool(false) != "false" {
		t.Error("boolean Pega harus teks true/false")
	}
}

// TestSQLObjekMemakaiKolomMigrasi - setiap kolom yang ditulis/dibaca ada di CREATE TABLE /
// ALTER TABLE ADD migrasi 186-187, dan jumlah bind = jumlah kolom tiap INSERT.
func TestSQLObjekMemakaiKolomMigrasi(t *testing.T) {
	sql := bacaMigrasi(t, "186_t_objek_fire.sql") + "\n" + bacaMigrasi(t, "187_t_surroundingrisk.sql")
	kolom := map[string]map[string]bool{}
	for _, m := range regexp.MustCompile(`(?s)(?:CREATE TABLE|ALTER TABLE) \{skema\}\.(\w+) (?:ADD )?\((.*?)\n\)`).FindAllStringSubmatch(sql, -1) {
		if kolom[m[1]] == nil {
			kolom[m[1]] = map[string]bool{}
		}
		for _, b := range strings.Split(m[2], "\n") {
			if f := strings.Fields(b); len(f) > 0 && f[0] != "CONSTRAINT" {
				kolom[m[1]][f[0]] = true
			}
		}
	}
	if len(kolom["T_SURROUNDINGRISK"]) != 24 || !kolom["T_PROPERTY"]["OWNERSHIP"] {
		t.Fatalf("migrasi 187 tak terbaca: %d kolom T_SURROUNDINGRISK", len(kolom["T_SURROUNDINGRISK"]))
	}
	tb := tabelObjek{loc: "T_LOCATIONLIST", prop: "T_PROPERTY", risk: "T_RISKLOCATION", bang: "T_BUILDINGCONSTRUCTION", sekitar: "T_SURROUNDINGRISK"}
	for _, q := range []string{sqlSisipLokasi(tb.loc), sqlSisipProperty(tb.prop), sqlSisipRisk(tb.risk), sqlSisipBangunan(tb.bang),
		sqlSisipSekitar(tb.sekitar)} {
		m := regexp.MustCompile(`INSERT INTO (\w+) \(([^)]*)\) VALUES \(([^)]*)\)`).FindStringSubmatch(q)
		if m == nil {
			t.Fatalf("INSERT tak terbaca: %q", q)
		}
		cols := strings.Split(m[2], ", ")
		if len(cols) != len(strings.Split(m[3], ", ")) {
			t.Errorf("%s: %d kolom, bind berbeda", m[1], len(cols))
		}
		for _, c := range cols {
			if !kolom[m[1]][c] {
				t.Errorf("%s.%s ditulis, tidak ada di 186/187", m[1], c)
			}
		}
	}
	alias := map[string]string{"p": "T_PROPERTY", "r": "T_RISKLOCATION", "b": "T_BUILDINGCONSTRUCTION", "s": "T_SURROUNDINGRISK"}
	for _, m := range regexp.MustCompile(`\b([prbs])\.([A-Z_]+)\b`).FindAllStringSubmatch(sqlBacaObjek(tabelObjekUji()), -1) {
		if !kolom[alias[m[1]]][m[2]] {
			t.Errorf("%s.%s dibaca, tidak ada di 186/187", alias[m[1]], m[2])
		}
	}
}

// TestBacaObjekMemakaiKolomBernama - setiap kunci v("...") di BacaObjek ada di
// kolomBacaObjek, dan setiap kolom terpilih dibaca (salah ketik = panic saat Oracle).
func TestBacaObjekMemakaiKolomBernama(t *testing.T) {
	b, err := os.ReadFile("objek.go")
	if err != nil {
		t.Fatal(err)
	}
	ada, dibaca := map[string]bool{}, map[string]bool{}
	for _, k := range kolomBacaObjek {
		ada[k] = true
	}
	for _, m := range regexp.MustCompile(`v\("([^"]+)"\)`).FindAllStringSubmatch(string(b), -1) {
		dibaca[m[1]] = true
		if !ada[m[1]] {
			t.Errorf("kunci %q dibaca tetapi tidak dipilih", m[1])
		}
	}
	for _, k := range kolomBacaObjek {
		if !dibaca[k] {
			t.Errorf("kolom %s dipilih tetapi tidak dibaca", k)
		}
	}
	if len(kolomBacaObjek) != 47 {
		t.Errorf("%d kolom, mau 47 (23 tiket 35 + 4 T_PROPERTY + 20 T_SURROUNDINGRISK tiket 38)", len(kolomBacaObjek))
	}
}

// TestArgSekitarMenurutKolom - tiket 38: bind sqlSisipSekitar dipasangkan lewat KUNCI (nilai uji =
// nama kolomnya sendiri), bukan lewat urutan; kosong -> NULL.
func TestArgSekitarMenurutKolom(t *testing.T) {
	sisi := func(p string) models.SisiRisiko {
		return models.SisiRisiko{Occupation: p + "_OCCUPATION", Construction: p + "_CONSTRUCTION", Distance: p + "_DISTANCE", Note: p + "_NOTE"}
	}
	s := models.SurroundingRisk{Front: sisi("FRONT"), Left: sisi("LEFT"), Back: sisi("BACK"), Right: sisi("RIGHT"),
		HousekeepingStatus: "HOUSEKEEPING_STATUS", FloodAreaStatus: "FLOOD_AREA_STATUS", FloodArea: "FLOOD_AREA",
		HousekeepingRemark: "HOUSEKEEPING_REMARK"}
	m := regexp.MustCompile(`\(([^)]*)\) VALUES`).FindStringSubmatch(sqlSisipSekitar("UJI.S"))
	cols := strings.Split(m[1], ", ")[2:]
	arg := argSekitar(s)
	if len(arg) != len(cols) || len(cols) != 20 {
		t.Fatalf("%d bind, %d kolom, mau 20", len(arg), len(cols))
	}
	for i, c := range cols {
		if arg[i] != c {
			t.Errorf("bind kolom %s = %v", c, arg[i])
		}
	}
	for i, a := range argSekitar(models.SurroundingRisk{}) {
		if a != nil {
			t.Errorf("kosong[%d] = %v, mau NULL", i, a)
		}
	}
}
