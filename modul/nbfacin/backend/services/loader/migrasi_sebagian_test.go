package loader

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// tipeAkhirMigrasi - tipe kolom tabel flat sesudah migrasi nbfacin 182..195 dijalankan
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
		"186_t_objek_fire.sql", "187_t_surroundingrisk.sql", "188_t_propertyitemlist.sql", "189_t_occupationlist.sql",
		"191_t_listcauseofloss.sql", "192_lebar_alamat_risiko.sql",
		"193_t_coveragelist.sql", "194_t_deductiblelist.sql", "195_t_coveragelist_unit_akumulasi.sql",
		"198_spreading.sql", "199_cedant.sql"} {
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
		"T_BUILDINGCONSTRUCTION.ID": "NUMBER(19)", "T_BUILDINGCONSTRUCTION.PARENT_ID": "NUMBER(19)",
		"T_SURROUNDINGRISK.ID": "NUMBER(19)", "T_SURROUNDINGRISK.PARENT_ID": "NUMBER(19)",
		"T_PROPERTYITEMLIST.ID": "NUMBER(19)", "T_PROPERTYITEMLIST.PARENT_ID": "NUMBER(19)",
		// tiket 39: uang (ADR-0016) dan persen (presisiSah penjaga) - rancangan NUMBER polos.
		"T_PROPERTYITEMLIST.TSI_OBJECT_ITEM": "NUMBER(38,8)", "T_PROPERTYITEMLIST.PCT_ADJUST2": "NUMBER(38,8)",
		"T_PROPERTYITEMLIST.PCT_ADJUST_OTHER": "NUMBER(38,8)",
		"T_OCCUPATIONLIST.ID":                 "NUMBER(19)", "T_OCCUPATIONLIST.PARENT_ID": "NUMBER(19)",
		"T_TABLEOFLIMIT.ID": "NUMBER(19)", "T_TABLEOFLIMIT.PARENT_ID": "NUMBER(19)",
		// tiket 40: butir 68.1 - rancangan NUMBER, isi teks apa adanya (lebar A139).
		"T_TABLEOFLIMIT.PCT_LIMIT": "VARCHAR2(50)",
		// tiket 42: uang (ADR-0016) - rancangan NUMBER polos.
		"T_LOCATIONLIST.LOSS_RATIO1_YEAR_AMOUNT": "NUMBER(38,8)", "T_LOCATIONLIST.LOSS_RATIO35_YEAR_AMOUNT": "NUMBER(38,8)",
		"T_LISTCAUSEOFLOSS.ID": "NUMBER(19)", "T_LISTCAUSEOFLOSS.PARENT_ID": "NUMBER(19)", "T_LISTCAUSEOFLOSS.CLAIM": "NUMBER(38,8)",
		"T_LISTCAUSEOFLOSS.AMOUNT": "NUMBER(38,8)", "T_LISTCAUSEOFLOSS.PREVENTION_OF_LOSS": "NUMBER(38,8)",
		"T_COINSDATA.ID": "NUMBER(19)", "T_COINSDATA.PARENT_ID": "NUMBER(19)",
		// tiket 43: uang / rate / persen (ADR-0016) - rancangan NUMBER polos.
		"T_PROPERTYITEMLIST.TOTAL_GROSS_PREMI": "NUMBER(38,8)", "T_PROPERTYITEMLIST.TOTAL_NET_RATE": "NUMBER(38,8)",
		"T_COVERAGELIST.ID": "NUMBER(19)", "T_COVERAGELIST.PARENT_ID": "NUMBER(19)",
		"T_COVERAGELIST.TSI":                  "NUMBER(38,8)",
		"T_COVERAGELIST.RATE":                 "NUMBER(38,8)",
		"T_COVERAGELIST.RATE_OJK":             "NUMBER(38,8)",
		"T_COVERAGELIST.DISCOUNT_PERCENTAGE":  "NUMBER(38,8)",
		"T_COVERAGELIST.TSI_LIABILITY":        "NUMBER(38,8)",
		"T_COVERAGELIST.NET_RATE":             "NUMBER(38,8)",
		"T_COVERAGELIST.LIMITOF_LIABILITY":    "NUMBER(38,8)",
		"T_COVERAGELIST.PCT_LO_L":             "NUMBER(38,8)",
		"T_COVERAGELIST.PRO_RATE_PERCENT":     "NUMBER(38,8)",
		"T_COVERAGELIST.INDEMNITY_PERCENTAGE": "NUMBER(38,8)",
		"T_COVERAGELIST.SUBLIMIT":             "NUMBER(38,8)",
		"T_COVERAGELIST.DISCOUNT":             "NUMBER(38,8)",
		"T_COVERAGELIST.PREMIUM":              "NUMBER(38,8)",
		"T_COVERAGELIST.PCT_ADJUSTMENT":       "NUMBER(38,8)",
		// tiket 45: uang / persen (ADR-0016); kode pyStandardValue 0-7 -> NUMBER(5) (A170) - rancangan NUMBER polos.
		"T_DEDUCTIBLELIST.ID": "NUMBER(19)", "T_DEDUCTIBLELIST.PARENT_ID": "NUMBER(19)",
		"T_DEDUCTIBLELIST.AMOUNT": "NUMBER(38,8)", "T_DEDUCTIBLELIST.PCT_DEDUCTIBLE": "NUMBER(38,8)",
		"T_DEDUCTIBLELIST.PCT_DEDUCTIBLE2": "NUMBER(38,8)", "T_DEDUCTIBLELIST.TYPE_DEDUCTIBLE": "NUMBER(5)",
		"T_DEDUCTIBLELIST.TYPE_DEDUCTIBLE2": "NUMBER(5)",
		// tiket 48 (198): persen / uang (ADR-0016) - rancangan NUMBER polos.
		"T_GENERAL_POLIS.PERCENT_SHARE": "NUMBER(38,8)", "T_COVERAGELIST.TSI_NUSANTARA_RE": "NUMBER(38,8)",
		"T_COVERAGELIST.PREMI_NUSANTARA_RE": "NUMBER(38,8)",
		// tiket 49 (199): kode -> NUMBER(5) (A170); ShareOfCeding teks "N%" apa adanya (pola PCT_LIMIT); PARENT_ID = PK
		// teks T_GENERAL_POLIS (butir 76.1); persen NUMBER(38,8).
		"T_GENERAL_POLIS.SHARE_CEDANT_TYPE": "NUMBER(5)", "T_QUOTATIONDATA.SHARE_OF_CEDING": "VARCHAR2(50)",
		"T_CEDINGCEDANTLIST.ID": "NUMBER(19)", "T_CEDINGCEDANTLIST.PARENT_ID": "VARCHAR2(32)",
		"T_CEDINGCEDANTLIST.SHARE_CEDING": "NUMBER(38,8)",
	}
	akhir := tipeAkhirMigrasi(t)
	diperiksa := 0
	for _, tabel := range []string{"T_GENERAL_POLIS", "T_QUOTATIONDATA", "T_CEDINGCOLIST", "T_LOCATIONLIST",
		"T_PROPERTY", "T_RISKLOCATION", "T_BUILDINGCONSTRUCTION", "T_SURROUNDINGRISK", "T_PROPERTYITEMLIST", "T_OCCUPATIONLIST", "T_TABLEOFLIMIT",
		"T_LISTCAUSEOFLOSS", "T_COINSDATA", "T_COVERAGELIST", "T_DEDUCTIBLELIST", "T_CEDINGCEDANTLIST"} {
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
	for _, tabel := range []string{"T_CEDINGCOLIST", "T_RISKLOCATION", "T_BUILDINGCONSTRUCTION", "T_SURROUNDINGRISK", "T_TABLEOFLIMIT", "T_LISTCAUSEOFLOSS",
		"T_COINSDATA", "T_CEDINGCEDANTLIST"} {
		for _, k := range skemaTabel[tabel] {
			if akhir[tabel][k.nama] == "" {
				t.Errorf("%s.%s rancangan tidak dibuat migrasi", tabel, k.nama)
			}
		}
	}
	// 7 T_GENERAL_POLIS + 15 T_QUOTATIONDATA + 8 T_CEDINGCOLIST + 6 T_LOCATIONLIST + 20 T_PROPERTY (16 + 4 di 187)
	// + 9 T_RISKLOCATION + 11 T_BUILDINGCONSTRUCTION + 24 T_SURROUNDINGRISK + 22 T_PROPERTYITEMLIST (16 rancangan
	// + 6 baru, 188) + 10 T_OCCUPATIONLIST + 7 T_TABLEOFLIMIT (189) = 139.
	// + 4 T_LOCATIONLIST (191) + 15 T_LISTCAUSEOFLOSS + 5 T_COINSDATA (191) + 2 T_PROPERTYITEMLIST + 34 T_COVERAGELIST
	// (193) = 199. + 18 T_DEDUCTIBLELIST (194) = 217. + 3 T_COVERAGELIST (195) = 220.
	// + 1 T_GENERAL_POLIS + 2 T_COVERAGELIST (198) = 223. + 1 T_GENERAL_POLIS + 1 T_QUOTATIONDATA + 9 T_CEDINGCEDANTLIST
	// (199) = 234.
	if diperiksa != 234 {
		t.Errorf("%d kolom diperiksa, mau 234", diperiksa)
	}
	// Butir 87/88: Risk Location / Address VARCHAR2(4000), delapan kolom alamat lain VARCHAR2(100) sesudah 192.
	for tk, mau := range map[string]string{"T_RISKLOCATION.ASM_ADDRESS": "VARCHAR2(4000)", "T_PROPERTY.ROAD_NAME": "VARCHAR2(4000)",
		"T_RISKLOCATION.ASM_CITY": "VARCHAR2(100)", "T_RISKLOCATION.ASM_DISTRICT": "VARCHAR2(100)", "T_RISKLOCATION.ASMRW": "VARCHAR2(100)",
		"T_RISKLOCATION.ASM_ZIP_CODE": "VARCHAR2(100)", "T_PROPERTY.ROAD_TYPE": "VARCHAR2(100)", "T_PROPERTY.PROVINCE": "VARCHAR2(100)",
		"T_PROPERTY.COUNTRY": "VARCHAR2(100)", "T_PROPERTY.ALM_RISK_ID": "VARCHAR2(100)"} {
		t1, k1, _ := strings.Cut(tk, ".")
		if akhir[t1][k1] != mau {
			t.Errorf("%s = %s, mau %s (butir 87/88)", tk, akhir[t1][k1], mau)
		}
	}
	if akhir["T_PROPERTY"]["BUILDING_NO"] != "VARCHAR2(50)" {
		t.Errorf("BUILDING_NO tidak dari RISKADDRESS, tetap 50: %s", akhir["T_PROPERTY"]["BUILDING_NO"])
	}
	if akhir["T_LISTCAUSEOFLOSS"]["DETAIL"] != "VARCHAR2(500)" {
		t.Errorf("DETAIL dilebarkan (A146): %v", akhir["T_LISTCAUSEOFLOSS"]["DETAIL"])
	}
	if akhir["T_OCCUPATIONLIST"]["OCCUPATION_ID"] != "VARCHAR2(1000)" || akhir["T_OCCUPATIONLIST"]["OCCUPATION_NAME"] != "VARCHAR2(1000)" {
		t.Errorf("lebar okupasi (A138): %v", akhir["T_OCCUPATIONLIST"])
	}
	if akhir["T_QUOTATIONDATA"]["CEDING_CO_NAME"] != "VARCHAR2(4000)" || akhir["T_QUOTATIONDATA"]["CEDING_CO"] != "VARCHAR2(1000)" {
		t.Errorf("kolom gabungan Ceding Co: %v (butir 80)", akhir["T_QUOTATIONDATA"])
	}
}
