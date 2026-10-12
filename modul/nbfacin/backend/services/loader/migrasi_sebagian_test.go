package loader

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
)

// tipeAkhirMigrasi - tipe kolom tabel flat sesudah migrasi nbfacin 182..205, 208, 210, 212, 213 dan 214 dijalankan BERURUTAN, pernyataan demi
// pernyataan: CREATE TABLE, ALTER ... ADD (...), ALTER ... MODIFY (KOLOM TIPE), dan ALTER ... DROP COLUMN (juga di dalam blok
// berpelindung katalog, `migrasi.KolomAlterBuang`) - 202 membuang lalu menambah ulang PPN_CHECK.
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
		"198_spreading.sql", "199_cedant.sql", "200_premium_deduction.sql", "201_payment.sql", "202_ppn_check_teks.sql",
		"203_coverage_zona.sql", "204_coverage_zona_4000.sql", "205_scoring_risk.sql", "208_nb_status.sql", "210_submit_team_leader.sql", "212_submit_uw.sql",
		"213_form_binding.sql", "214_objek_aneka.sql", "251_lebar_coverage_note.sql", "252_general_polis_lengkap.sql",
		"253_quotationdata_lengkap.sql", "255_pro_rate_percent.sql"} {
		b, err := os.ReadFile("../../migrations/" + berkas)
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range strings.Split(strings.ReplaceAll(string(b), "\r", ""), "\n/\n") {
			for _, m := range polaCreate.FindAllStringSubmatch(p, -1) {
				kolom(m[1], m[2])
			}
			for _, m := range polaAdd.FindAllStringSubmatch(p, -1) {
				kolom(m[1], m[2])
			}
			for _, m := range polaModify.FindAllStringSubmatch(p, -1) {
				if akhir[m[1]][m[2]] == "" {
					t.Fatalf("%s: MODIFY %s.%s sebelum kolomnya ada", berkas, m[1], m[2])
				}
				akhir[m[1]][m[2]] = m[3]
			}
			if tabel, buang := migrasi.KolomAlterBuang(p); tabel != "" {
				for _, k := range buang {
					if akhir[tabel][k] == "" {
						t.Fatalf("%s: DROP COLUMN %s.%s yang tidak ada", berkas, tabel, k)
					}
					delete(akhir[tabel], k)
				}
			}
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
		// tiket 50 (200): PPN_CHECK teks apa adanya (butir 68.1, kolomTeksMenyimpang); persen / uang NUMBER(38,8) (ADR-0016) -
		// rancangan NUMBER polos.
		"T_GENERAL_POLIS.PPN_CHECK": "VARCHAR2(10)", "T_GENERAL_POLIS.PAY_RI_COMMISION": "NUMBER(38,8)",
		"T_GENERAL_POLIS.PAY_PCT_BROKERAGE_FEE": "NUMBER(38,8)", "T_CURRENCYLIST.ID": "NUMBER(19)",
		"T_CURRENCYLIST.PAY_PREMIUM": "NUMBER(38,8)", "T_CURRENCYLIST.PAY_COMMISION": "NUMBER(38,8)",
		"T_CURRENCYLIST.PAY_DEDUCTION2": "NUMBER(38,8)", "T_CURRENCYLIST.PAY_BROKERAGE_FEE": "NUMBER(38,8)",
		"T_CURRENCYLIST.PAY_PPH": "NUMBER(38,8)", "T_CURRENCYLIST.PAY_PPN": "NUMBER(38,8)",
		// tiket 51 (201): jumlah cicilan NUMBER(5); uang NUMBER(38,8) (ADR-0016; DEDUCTION2 rancangan VARCHAR2(50));
		// persen cicilan teks desimal apa adanya (K51-3, pola PCT_LIMIT); ID / PARENT_ID NUMBER(19) (A92).
		"T_GENERAL_POLIS.PAY_INSTALLMENT": "NUMBER(5)", "T_CURRENCYLIST.SUM_TOTAL_PAYMENT": "NUMBER(38,8)",
		"T_CURRENCYLIST.PAY_NET_PREMIUM": "NUMBER(38,8)", "T_LISTINSTALLMENT.ID": "NUMBER(19)",
		"T_LISTINSTALLMENT.PARENT_ID": "NUMBER(19)", "T_LISTINSTALLMENT.INSTALLMENT_NO": "NUMBER(5)",
		"T_LISTINSTALLMENT.INSTALLMENT_PERCENTAGE": "VARCHAR2(50)", "T_LISTINSTALLMENT.PREMIUM": "NUMBER(38,8)",
		"T_LISTINSTALLMENT.RI_COMMISION": "NUMBER(38,8)", "T_LISTINSTALLMENT.DEDUCTION2": "NUMBER(38,8)",
		"T_LISTINSTALLMENT.BROKERAGE_FEE": "NUMBER(38,8)", "T_LISTINSTALLMENT.PPH": "NUMBER(38,8)",
		"T_LISTINSTALLMENT.PPN": "NUMBER(38,8)", "T_LISTINSTALLMENT.DISCOUNT": "NUMBER(38,8)",
		"T_LISTINSTALLMENT.PAYMENT_TOTAL": "NUMBER(38,8)",
		// tiket 54 (205): ID / PARENT_ID NUMBER(19) (A92); skor / persen NUMBER(38,8) - rancangan NUMBER polos.
		"T_SCORINGRISK.ID": "NUMBER(19)", "T_SCORINGRISK.FINAL_SCORE": "NUMBER(38,8)",
		"T_SCORINGRESULT.ID": "NUMBER(19)", "T_SCORINGRESULT.PARENT_ID": "NUMBER(19)",
		"T_DATASCORINGRISKLIST.ID": "NUMBER(19)", "T_DATASCORINGRISKLIST.PARENT_ID": "NUMBER(19)",
		"T_DATASCORINGRISKLIST.SCORE": "NUMBER(38,8)", "T_SCORING_FACTOR.ID": "NUMBER(19)", "T_SCORING_FACTOR.PARENT_ID": "NUMBER(19)",
		"T_SCORING_FACTOR.SCORE_PER_FACTOR": "NUMBER(38,8)", "T_SCORING_FACTOR.PERCENT_LOSS_RATIO": "NUMBER(38,8)",
		"T_SCORING_OPTION.ID": "NUMBER(19)", "T_SCORING_OPTION.PARENT_ID": "NUMBER(19)", "T_SCORING_OPTION.SCORE": "NUMBER(38,8)",
		// tiket 73 (214): ID / PARENT_ID NUMBER(19) (A92); uang / angka K73-6 NUMBER(38,8) (ADR-0016) - rancangan / amandemen NUMBER
		// polos; TYPE_OF_DISCOUNT teks (kolomTeksMenyimpang, isi sampel "true").
		"T_ANEKALIST.ID": "NUMBER(19)", "T_ANEKALIST.PARENT_ID": "NUMBER(19)", "T_ANEKALIST.TSI": "NUMBER(38,8)",
		"T_ANEKALIST.MAINTENANCE_PERIOD": "NUMBER(38,8)", "T_ANEKALIST.NO_OF_TREE": "NUMBER(38,8)",
		"T_ANEKALIST.AREA_HECTAR": "NUMBER(38,8)", "T_ANEKALIST.LIMIT_IDEMNITY_DAY": "NUMBER(38,8)",
		"T_ANEKALIST.LIMIT_IDEMNITY_MONTH": "NUMBER(38,8)", "T_ANEKALIST.LIMIT_OF_LIABILITY": "NUMBER(38,8)",
		"T_ANEKALIST.TSI_OLD": "NUMBER(38,8)", "T_ANEKALIST.SALARY": "NUMBER(38,8)",
		"T_MARINEHULL.ID": "NUMBER(19)", "T_MARINEHULL.PARENT_ID": "NUMBER(19)",
		"T_AVIATIONHULL.ID": "NUMBER(19)", "T_AVIATIONHULL.PARENT_ID": "NUMBER(19)",
		"T_VEHICLEHE.ID": "NUMBER(19)", "T_VEHICLEHE.PARENT_ID": "NUMBER(19)",
		"T_MEMORANDUMACCEPT.ID": "NUMBER(19)", "T_MEMORANDUMACCEPT.PARENT_ID": "NUMBER(19)",
		"T_COVERAGELIST.PCT_SHORT_PERIOD": "NUMBER(38,8)", "T_COVERAGELIST.LOADING": "NUMBER(38,8)",
		"T_COVERAGELIST.MIN_PREMIUM": "NUMBER(38,8)", "T_COVERAGELIST.TYPE_OF_DISCOUNT": "VARCHAR2(10)",
		// tiket 73 P2 (255): % Pro Rate case, persen NUMBER(38,8) - rancangan NUMBER polos.
		"T_GENERAL_POLIS.PRO_RATE_PERCENT": "NUMBER(38,8)",
		// 252 (keputusan work owner 12-10-2026): NUMBER polos rancangan -> presisi sah penjaga - uang / persen / share NUMBER(38,8)
		// (ADR-0016), identitas NUMBER(19) (A92), kode NUMBER(5) (A170).
		"T_GENERAL_POLIS.CEDING_RETENTION": "NUMBER(38,8)", "T_GENERAL_POLIS.CEDING_RETENTION_NOMINAL": "NUMBER(38,8)",
		"T_GENERAL_POLIS.OLD_POLIS_ID": "NUMBER(19)", "T_GENERAL_POLIS.PRO_RATE_TYPE": "NUMBER(5)",
		"T_GENERAL_POLIS.SHARE_RNML": "NUMBER(38,8)", "T_GENERAL_POLIS.SUM_TOTAL_TSI": "NUMBER(38,8)",
		"T_GENERAL_POLIS.TOTAL_PREMI_NUSA_RE": "NUMBER(38,8)", "T_GENERAL_POLIS.TOTAL_TSI_NUSA_RE": "NUMBER(38,8)",
		"T_GENERAL_POLIS.TOTAL_TSI_NUSA_RE_SPREADING": "NUMBER(38,8)", "T_GENERAL_POLIS.TOTAL_TSI_TOP_RISK": "NUMBER(38,8)",
		"T_GENERAL_POLIS.TREATY_CAPACITY": "NUMBER(38,8)",
	}
	akhir := tipeAkhirMigrasi(t)
	diperiksa := 0
	for _, tabel := range []string{"T_GENERAL_POLIS", "T_QUOTATIONDATA", "T_CEDINGCOLIST", "T_LOCATIONLIST",
		"T_PROPERTY", "T_RISKLOCATION", "T_BUILDINGCONSTRUCTION", "T_SURROUNDINGRISK", "T_PROPERTYITEMLIST", "T_OCCUPATIONLIST", "T_TABLEOFLIMIT",
		"T_LISTCAUSEOFLOSS", "T_COINSDATA", "T_COVERAGELIST", "T_DEDUCTIBLELIST", "T_CEDINGCEDANTLIST", "T_CURRENCYLIST", "T_LISTINSTALLMENT",
		"T_SCORINGRISK", "T_SCORINGRESULT", "T_DATASCORINGRISKLIST", "T_SCORING_FACTOR", "T_SCORING_OPTION",
		"T_ANEKALIST", "T_MARINEHULL", "T_AVIATIONHULL", "T_VEHICLEHE", "T_MEMORANDUMACCEPT"} {
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
		"T_COINSDATA", "T_CEDINGCEDANTLIST",
		"T_SCORINGRISK", "T_SCORINGRESULT", "T_DATASCORINGRISKLIST", "T_SCORING_FACTOR", "T_SCORING_OPTION",
		// T_ANEKALIST tidak: penunjuk IDX_ANEKA / IDX_LOCATION / IDX_OCCUPATION (butir 72) milik loader, tidak dibuat (pola 198).
		"T_MARINEHULL", "T_AVIATIONHULL", "T_VEHICLEHE", "T_MEMORANDUMACCEPT"} {
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
	// (199) = 234. + 3 T_GENERAL_POLIS + 14 T_CURRENCYLIST (200) = 251. + 1 T_GENERAL_POLIS + 2 T_CURRENCYLIST +
	// 18 T_LISTINSTALLMENT (201) = 272. + 2 T_COVERAGELIST (203) = 274. 204 melebarkan ZONE (tanpa kolom baru) = 274.
	// + 7 T_SCORINGRISK + 20 T_SCORINGRESULT + 11 T_DATASCORINGRISKLIST + 12 T_SCORING_FACTOR + 9 T_SCORING_OPTION (205) = 333.
	// 206 / 207 tabel baru di luar rancangan (tidak diputar). + 1 T_GENERAL_POLIS NB_STATUS (208) = 334.
	// 209 MODIFY NULL (tipe tetap). + 4 T_GENERAL_POLIS (210) = 338. 211 data saja. + 4 T_GENERAL_POLIS (212) = 342.
	// + 2 T_GENERAL_POLIS (213) = 344. + 50 T_ANEKALIST + 13 T_MARINEHULL + 10 T_AVIATIONHULL + 12 T_VEHICLEHE
	// + 20 T_MEMORANDUMACCEPT + 8 T_COVERAGELIST (214) = 457. + 1 T_GENERAL_POLIS PRO_RATE_PERCENT (255) = 458.
	// + 48 T_GENERAL_POLIS (252: 53 rancangan belum ada dikurangi 5 Endorsement) = 506.
	// + 40 T_QUOTATIONDATA (253: 51 rancangan belum ada dikurangi 11 Endorsement / Renewal) = 546.
	if diperiksa != 546 {
		t.Errorf("%d kolom diperiksa, mau 546", diperiksa)
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

// TestMigrasi252GeneralPolisLengkap - keputusan work owner 12-10-2026, dua arah: sesudah 252 setiap kolom rancangan T_GENERAL_POLIS
// (sesudah amandemen) ADA di migrasi KECUALI tepat lima kolom khusus Endorsement; dan setiap kolom migrasi ada di rancangan (dijaga
// TestMigrasiFlatSebagianCocokRancangan, tipe sama / selisih presisi). Jumlah akhir 25 + 48 = 73.
func TestMigrasi252GeneralPolisLengkap(t *testing.T) {
	akhir := tipeAkhirMigrasi(t)["T_GENERAL_POLIS"]
	edm := map[string]bool{"ENDORSEMENT_NO": true, "ENDORSMENT_REASON": true, "PRORATE_START_EDM": true, "PRORATE_EDM_END": true,
		"ERRMSG": true}
	var hilang []string
	for _, k := range skemaTabel["T_GENERAL_POLIS"] {
		if akhir[k.nama] == "" && !edm[k.nama] {
			hilang = append(hilang, k.nama)
		}
		if akhir[k.nama] != "" && edm[k.nama] {
			t.Errorf("%s khusus Endorsement ikut dibuat", k.nama)
		}
	}
	if len(hilang) > 0 {
		t.Errorf("kolom rancangan belum dibuat migrasi: %v", hilang)
	}
	for k := range akhir {
		if indeksKolom("T_GENERAL_POLIS", k) < 0 {
			t.Errorf("%s di migrasi, tidak di rancangan", k)
		}
	}
	if len(akhir) != 73 || len(skemaTabel["T_GENERAL_POLIS"]) != 78 {
		t.Errorf("fisik %d (mau 73), rancangan %d (mau 78)", len(akhir), len(skemaTabel["T_GENERAL_POLIS"]))
	}
	for k, mau := range map[string]string{"POLICY_NO": "VARCHAR2(50)", "PROD_DATE_TIME": "VARCHAR2(30)", "PROD_KE": "NUMBER(5)"} {
		if akhir[k] != mau {
			t.Errorf("tahap 1 %s = %s, mau %s", k, akhir[k], mau)
		}
	}
	turun, err := os.ReadFile("../../migrations/252_general_polis_lengkap_down.sql")
	if err != nil {
		t.Fatal(err)
	}
	naik, _ := os.ReadFile("../../migrations/252_general_polis_lengkap.sql")
	for _, baris := range strings.Split(strings.ReplaceAll(string(naik), "\r", ""), "\n") {
		if f := strings.Fields(strings.TrimSuffix(strings.TrimSpace(baris), ",")); len(f) == 2 && strings.HasPrefix(f[1], "N") ||
			len(f) == 2 && strings.HasPrefix(f[1], "V") {
			if !strings.Contains(string(turun), "  "+f[0]+",") && !strings.Contains(string(turun), "  "+f[0]+"\n") {
				t.Errorf("_down tidak membuang %s", f[0])
			}
		}
	}
}

// TestMigrasi253QuotationDataLengkap - tiket 76b (keputusan work owner 12-10-2026), dua arah: sesudah 253 setiap kolom rancangan
// T_QUOTATIONDATA ADA di migrasi KECUALI tepat sebelas kolom Endorsement / Renewal; setiap kolom migrasi ada di rancangan; 16 + 40 =
// 56; _down membuang ke-40 kolom 253.
func TestMigrasi253QuotationDataLengkap(t *testing.T) {
	akhir := tipeAkhirMigrasi(t)["T_QUOTATIONDATA"]
	edm := map[string]bool{"EDM_CHARGE_FEE": true, "EDM_DATE": true, "EDM_NOTE": true, "EDM_SOURCE": true, "EDM_SOURCE_NOTE": true,
		"EDM_STATUS": true, "EDM_SURVEY": true, "EDM_TYPE": true, "EDM_TYPE_NEW": true, "ENDORSEMENT_INTERNAL_RETRO": true, "RNW_DATE": true}
	var hilang []string
	for _, k := range skemaTabel["T_QUOTATIONDATA"] {
		if akhir[k.nama] == "" && !edm[k.nama] {
			hilang = append(hilang, k.nama)
		}
		if akhir[k.nama] != "" && edm[k.nama] {
			t.Errorf("%s khusus Endorsement / Renewal ikut dibuat", k.nama)
		}
	}
	if len(hilang) > 0 {
		t.Errorf("kolom rancangan belum dibuat migrasi: %v", hilang)
	}
	for k := range akhir {
		if indeksKolom("T_QUOTATIONDATA", k) < 0 {
			t.Errorf("%s di migrasi, tidak di rancangan", k)
		}
	}
	if len(akhir) != 56 || len(skemaTabel["T_QUOTATIONDATA"]) != 67 {
		t.Errorf("fisik %d (mau 56), rancangan %d (mau 67)", len(akhir), len(skemaTabel["T_QUOTATIONDATA"]))
	}
	naik, err := os.ReadFile("../../migrations/253_quotationdata_lengkap.sql")
	if err != nil {
		t.Fatal(err)
	}
	turun, err := os.ReadFile("../../migrations/253_quotationdata_lengkap_down.sql")
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, baris := range strings.Split(strings.ReplaceAll(string(naik), "\r", ""), "\n") {
		f := strings.Fields(strings.TrimSuffix(strings.TrimSpace(baris), ","))
		if len(f) == 2 && strings.HasPrefix(f[1], "VARCHAR2(") {
			n++
			if !strings.Contains(string(turun), "  "+f[0]+",") && !strings.Contains(string(turun), "  "+f[0]+"\n") {
				t.Errorf("_down tidak membuang %s", f[0])
			}
		}
	}
	if n != 40 {
		t.Errorf("253 menambah %d kolom, mau 40", n)
	}
}
