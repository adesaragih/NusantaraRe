package loader

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// tipeAkhirMigrasi - tipe kolom tabel flat sesudah migrasi nbfacin 182..186 dijalankan
// BERURUTAN: CREATE TABLE, lalu ALTER ... ADD (...) dan ALTER ... MODIFY (KOLOM TIPE).
func tipeAkhirMigrasi(t *testing.T) map[string]map[string]string {
	t.Helper()
	polaCreate := regexp.MustCompile(`(?s)CREATE TABLE \{skema\}\.(\w+) \((.*?)\n\)`)
	polaAdd := regexp.MustCompile(`(?s)ALTER TABLE \{skema\}\.(\w+) ADD \((.*?)\n\)`)
	polaModify := regexp.MustCompile(`ALTER TABLE \{skema\}\.(\w+) MODIFY \((\w+) ([A-Z0-9_]+(?:\([^)]*\))?)\)`)
	akhir := map[string]map[string]string{}
	kolom := func(tabel, badan string) {
		if akhir[tabel] == nil {
			akhir[tabel] = map[string]string{}
		}
		for _, baris := range strings.Split(badan, "\n") {
			f := strings.Fields(strings.TrimSuffix(strings.TrimSpace(baris), ","))
			if len(f) >= 2 && f[0] != "CONSTRAINT" {
				akhir[tabel][f[0]] = strings.TrimSuffix(f[1], ",")
			}
		}
	}
	for _, berkas := range []string{"182_t_general_polis.sql", "183_t_quotationdata.sql", "184_t_quotationdata_sob.sql", "185_t_cedingcolist.sql",
		"186_t_objek_fire.sql"} {
		b, err := os.ReadFile("../../migrations/" + berkas)
		if err != nil {
			t.Fatal(err)
		}
		isi := strings.ReplaceAll(string(b), "\r", "")
		for _, m := range polaCreate.FindAllStringSubmatch(isi, -1) {
			kolom(m[1], m[2])
		}
		for _, m := range polaAdd.FindAllStringSubmatch(isi, -1) {
			kolom(m[1], m[2])
		}
		for _, m := range polaModify.FindAllStringSubmatch(isi, -1) {
			if akhir[m[1]][m[2]] == "" {
				t.Fatalf("%s: MODIFY %s.%s sebelum kolomnya ada", berkas, m[1], m[2])
			}
			akhir[m[1]][m[2]] = m[3]
		}
	}
	return akhir
}

// TestMigrasiFlatSebagianCocokRancangan - butir 78.4/80: tabel flat yang dibuat nbfacin
// (T_GENERAL_POLIS, T_QUOTATIONDATA sebagian; T_CEDINGCOLIST utuh) - setiap kolom pada
// keadaan AKHIR migrasi 182..185 wajib ADA di rancangan (skemaTabel sesudah amandemen
// butir 76/80) dengan tipe SAMA. Selisih yang diizinkan: identitas dari sequence NUMBER ->
// NUMBER(19) (A92) - T_QUOTATIONDATA.ID, T_CEDINGCOLIST.ID dan PARENT_ID.
func TestMigrasiFlatSebagianCocokRancangan(t *testing.T) {
	selisih := map[string]string{"T_QUOTATIONDATA.ID": "NUMBER(19)", "T_CEDINGCOLIST.ID": "NUMBER(19)",
		"T_CEDINGCOLIST.PARENT_ID": "NUMBER(19)", "T_LOCATIONLIST.ID": "NUMBER(19)", "T_PROPERTY.ID": "NUMBER(19)",
		"T_PROPERTY.PARENT_ID": "NUMBER(19)", "T_RISKLOCATION.ID": "NUMBER(19)", "T_RISKLOCATION.PARENT_ID": "NUMBER(19)",
		"T_BUILDINGCONSTRUCTION.ID": "NUMBER(19)", "T_BUILDINGCONSTRUCTION.PARENT_ID": "NUMBER(19)"}
	akhir := tipeAkhirMigrasi(t)
	diperiksa := 0
	for _, tabel := range []string{"T_GENERAL_POLIS", "T_QUOTATIONDATA", "T_CEDINGCOLIST", "T_LOCATIONLIST",
		"T_PROPERTY", "T_RISKLOCATION", "T_BUILDINGCONSTRUCTION"} {
		for k, tipe := range akhir[tabel] {
			// "VARCHAR2(10)" dari "VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL" - tipe saja yang dibandingkan.
			if i := indeksKolom(tabel, k); i >= 0 && strings.HasPrefix(skemaTabel[tabel][i].tipe, tipe+" ") {
				diperiksa++
				continue
			}
			diperiksa++
			i := indeksKolom(tabel, k)
			if i < 0 {
				t.Errorf("%s.%s tidak ada di rancangan", tabel, k)
				continue
			}
			mau := skemaTabel[tabel][i].tipe
			if s, ada := selisih[tabel+"."+k]; ada {
				mau = s
			}
			if tipe != mau {
				t.Errorf("%s.%s migrasi %s, rancangan %s", tabel, k, tipe, mau)
			}
		}
	}
	// Dibuat UTUH: setiap kolom rancangan (sesudah amandemen) ada di migrasi.
	for _, tabel := range []string{"T_CEDINGCOLIST", "T_RISKLOCATION", "T_BUILDINGCONSTRUCTION"} {
		for _, k := range skemaTabel[tabel] {
			if akhir[tabel][k.nama] == "" {
				t.Errorf("%s.%s rancangan tidak dibuat migrasi", tabel, k.nama)
			}
		}
	}
	// 7 T_GENERAL_POLIS + 15 T_QUOTATIONDATA + 8 T_CEDINGCOLIST + 6 T_LOCATIONLIST + 16 T_PROPERTY
	// + 9 T_RISKLOCATION + 11 T_BUILDINGCONSTRUCTION = 72.
	if diperiksa != 72 {
		t.Errorf("%d kolom diperiksa, mau 72", diperiksa)
	}
	if akhir["T_QUOTATIONDATA"]["CEDING_CO_NAME"] != "VARCHAR2(4000)" || akhir["T_QUOTATIONDATA"]["CEDING_CO"] != "VARCHAR2(1000)" {
		t.Errorf("kolom gabungan Ceding Co: %v (butir 80)", akhir["T_QUOTATIONDATA"])
	}
}
