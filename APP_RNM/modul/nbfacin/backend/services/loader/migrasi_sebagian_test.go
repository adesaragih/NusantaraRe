package loader

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// TestMigrasiFlatSebagianCocokRancangan - butir 78.4: T_GENERAL_POLIS (182) dan
// T_QUOTATIONDATA (183) dibuat SEBAGIAN; setiap kolomnya wajib ADA di rancangan
// (skemaTabel, sesudah amandemen butir 76) dengan tipe yang SAMA. Satu-satunya selisih
// yang diizinkan: T_QUOTATIONDATA.ID NUMBER -> NUMBER(19) (identitas dari sequence,
// A92 - presisi kolom ID belum diputus tim inti).
func TestMigrasiFlatSebagianCocokRancangan(t *testing.T) {
	selisih := map[string]string{"T_QUOTATIONDATA.ID": "NUMBER(19)"}
	diperiksa := 0
	for berkas, tabel := range map[string]string{"182_t_general_polis.sql": "T_GENERAL_POLIS", "183_t_quotationdata.sql": "T_QUOTATIONDATA"} {
		b, err := os.ReadFile("../../migrations/" + berkas)
		if err != nil {
			t.Fatal(err)
		}
		blok := regexp.MustCompile(`(?s)CREATE TABLE \{skema\}\.` + tabel + ` \((.*?)\n\)`).FindStringSubmatch(strings.ReplaceAll(string(b), "\r", ""))
		if blok == nil {
			t.Fatalf("%s: CREATE TABLE %s tidak terbaca", berkas, tabel)
		}
		for _, baris := range strings.Split(blok[1], "\n") {
			f := strings.Fields(strings.TrimSuffix(strings.TrimSpace(baris), ","))
			if len(f) < 2 || f[0] == "CONSTRAINT" {
				continue
			}
			diperiksa++
			i := indeksKolom(tabel, f[0])
			if i < 0 {
				t.Errorf("%s.%s tidak ada di rancangan", tabel, f[0])
				continue
			}
			tipe, mau := strings.TrimSuffix(f[1], ","), skemaTabel[tabel][i].tipe
			if s, ada := selisih[tabel+"."+f[0]]; ada {
				mau = s
			}
			if tipe != mau {
				t.Errorf("%s.%s migrasi %s, rancangan %s", tabel, f[0], tipe, mau)
			}
		}
	}
	// 7 kolom T_GENERAL_POLIS + 13 kolom T_QUOTATIONDATA.
	if diperiksa != 20 {
		t.Errorf("%d kolom diperiksa, mau 20", diperiksa)
	}
}
