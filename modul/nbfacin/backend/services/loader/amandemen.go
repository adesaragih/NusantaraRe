package loader

import (
	"sort"
	"strings"
)

// Amandemen rancangan - keputusan work owner 02-10-2026, butir 70 (diteruskan sesi
// `nusantarare-0f`, "setuju") atas usulan `docs/USULAN-KOLOM-PENAMPUNG.md`. Workbook
// dan DDL draf (`D:\migrasi\RNM\OUTPUT\08-flat\`) tetap READ-ONLY: kolom dan tabel di
// bawah TIDAK ada di sana; skema_gen.go juga tidak disunting. Keduanya digabung saat
// paket dimuat. DDL / migrasinya menunggu tiket 23 (ditahan).
//
//	P1  T_COVERAGELIST.COVERAGE_INITIAL VARCHAR2(500)            <- CoverageList.CoverageInitial
//	P2-3 tabel T_ADDITIONALSHIP (anak T_SHIP, berulang)          <- Ship/AdditionalShip[n]
//	P5  T_CURRENCYLIST.CURRENCY_REF_ID VARCHAR2(50)              <- CurrencyList.ID (pola V-49)
//	P6  T_FR_CURRENCYLIST.POLICY_TSI                             <- FR CurrencyList/Policy.TSI
//	    tipe "mengikuti keputusan presisi tim inti": ditulis NUMBER tanpa presisi -
//	    penanda, bukan pilihan; Flatten menyimpannya desimal eksak seperti kolom uang lain
//
// P4 (IsCedingConfirm) TIDAK di sini: kolomnya di tabel lama
// POOLDATA.HISTORYAKSEPTASIPRODUCTION, menunggu DBA - tetap di penampung.

var amandemenKolom = map[string][]kolomSkema{
	"T_COVERAGELIST":    {{nama: "COVERAGE_INITIAL", tipe: "VARCHAR2(500)", medan: "CoverageInitial"}},
	"T_CURRENCYLIST":    {{nama: "CURRENCY_REF_ID", tipe: "VARCHAR2(50)", turunan: true}},
	"T_FR_CURRENCYLIST": {{nama: "POLICY_TSI", tipe: "NUMBER", turunan: true}},
	// Butir 76.2: penanda lini K-064, tipe kolom yang ada (premiumlistlife 050).
	"T_WORK_POLIS": {{nama: "LINI", tipe: "VARCHAR2(255)"}},
}

// Penyelarasan dengan T_WORK_POLIS yang ADA - butir 76 (keputusan work owner
// 03-10-2026, AskUserQuestion di sesi ini). K-064: T_WORK_POLIS dan T_GENERAL_POLIS
// adalah tabel yang SAMA dengan lini lain; T_WORK_POLIS sudah dibuat premiumlistlife
// (050, diubah 057/059/063): ID VARCHAR2(32) NOT NULL berisi pengenal work, LINI,
// POSITION, STATUS_WORK, FLAG_ONGOING_POLICY, COVER_KEY, CREATE_OP, CREATE_OP_NAME,
// TGL_CREATE, TGL_UPDATE. Fac In menyambung, tidak membuat ulang.
//
//	76.1 ID mengikuti tabel yang ada: tipeIDKasus, isinya pengenal work (pyID =
//	     NO_WORK, mis. NB-184351), bukan surrogate NUMBER. T_GENERAL_POLIS berbagi PK
//	     (K-064 relasi 52, ID = ID) jadi ikut; PARENT_ID tabel yang berinduk salah satu
//	     tabel itu juga (10 tabel, dihitung dari jalurSumber di init). IDPEGA,
//	     JENIS_WORK, NO_WORK tetap kolom tambahan.
//	76.4 Empat kolom rancangan DIGABUNG ke kolom yang ada, bertipe kolom yang ada -
//	     tidak satu pun menyempit (VARCHAR2(100)/(20)/(50) -> (255)/(255)/(64); DATE ->
//	     DATE). Arti kolom yang ada TIDAK diubah; pengisinya (repository, tiket 24)
//	     menulis kolom pasangannya.
const tipeIDKasus = "VARCHAR2(32)"

var tabelIDKasus = []string{"T_WORK_POLIS", "T_GENERAL_POLIS"}

var amandemenGabung = map[string]map[string]kolomSkema{
	"T_WORK_POLIS": {
		"POSISI":        {nama: "POSITION", tipe: "VARCHAR2(255)"},
		"STATUS_PROSES": {nama: "STATUS_WORK", tipe: "VARCHAR2(255)"},
		"TGL_INPUT":     {nama: "TGL_CREATE", tipe: "DATE"},
		"USERNAME":      {nama: "CREATE_OP", tipe: "VARCHAR2(64)"},
	},
}

// anakIDKasus - tabel yang PARENT_ID-nya diganti tipeIDKasus (diisi init; dibaca uji).
var anakIDKasus []string

// amandemenLebar - butir 80 (keputusan work owner 03-10-2026, tiket 34): kolom gabungan `;`
// daftar Ceding Co di T_QUOTATIONDATA dilebarkan - kode = DDL Pega FACINOFFER/FACINPRODUCTION
// CEDINGCO VARCHAR2(1000), nama VARCHAR2(4000). Migrasi 185.
// Tiket 40 (A138, pola butir 80): kode/nama okupasi = lebar sumbernya OCCUPATION.OLDID / NAME VARCHAR2(1000)
// (DDL OCCUPATION.txt); rancangan 50 / 500. Migrasi 189.
var amandemenLebar = map[string]string{
	"T_QUOTATIONDATA.CEDING_CO":        "VARCHAR2(1000)",
	"T_QUOTATIONDATA.CEDING_CO_NAME":   "VARCHAR2(4000)",
	"T_OCCUPATIONLIST.OCCUPATION_ID":   "VARCHAR2(1000)",
	"T_OCCUPATIONLIST.OCCUPATION_NAME": "VARCHAR2(1000)",
	// Tiket 42 (A146): Loss Detail pxTextArea - rancangan 50 (diturunkan dari contoh kosong), pola V-6 catatan 500. Migrasi 191.
	"T_LISTCAUSEOFLOSS.DETAIL": "VARCHAR2(500)",
	// Butir 87/88 (perintah work owner 03-10-2026, bug DEV): Risk Location / Address = lebar sumber RISKADDRESS
	// VARCHAR2(4000) (pola butir 80); delapan kolom alamat lain VARCHAR2(100). Migrasi 192.
	"T_RISKLOCATION.ASM_ADDRESS":  "VARCHAR2(4000)",
	"T_RISKLOCATION.ASM_CITY":     "VARCHAR2(100)",
	"T_RISKLOCATION.ASM_DISTRICT": "VARCHAR2(100)",
	"T_RISKLOCATION.ASMRW":        "VARCHAR2(100)",
	"T_RISKLOCATION.ASM_ZIP_CODE": "VARCHAR2(100)",
	"T_PROPERTY.ROAD_NAME":        "VARCHAR2(4000)",
	"T_PROPERTY.ROAD_TYPE":        "VARCHAR2(100)",
	"T_PROPERTY.PROVINCE":         "VARCHAR2(100)",
	"T_PROPERTY.COUNTRY":          "VARCHAR2(100)",
	"T_PROPERTY.ALM_RISK_ID":      "VARCHAR2(100)",
	// Tiket 52 (pola butir 80): Zone = ZONES.DESCRIPTION / ZONEOJK.DESCRIPTION VARCHAR2(1000) (DDL ZONEOJK.txt; RD
	// BrowseZones_RD menulis .Description ke .Zone); rancangan 50, sampel `DDL\CONTOH` sampai 111 byte. Migrasi 203.
	// Dilebarkan VARCHAR2(4000) - bug DEV 06-10-2026 (Zone > 1000 byte; ZONES.DESCRIPTION view atas JSON), perintah work
	// owner "buat 4000". Migrasi 204.
	"T_COVERAGELIST.ZONE": "VARCHAR2(4000)",
	// Bug DEV 11-10-2026 (CSV Marine Hull ORA-12899, pola butir 80): Coverage Note = lebar sumbernya POOLDATA.COVERAGE.NAME
	// VARCHAR2(4000) (`DDL\COVERAGE.txt`); rancangan 500. Migrasi 251. (T_ADDITIONALCOVERAGE.COVERAGE_NOTE tetap 500: tabelnya lahir di
	// 256, sesudah 251 menurut urutan nama berkas.)
	"T_COVERAGELIST.COVERAGE_NOTE": "VARCHAR2(4000)",
	// Tiket 54 (pola butir 80): kolom Scoring Risk = lebar sumbernya - TopRiskLocation = .Property.RiskLocation.ASMAddress
	// (T_RISKLOCATION.ASM_ADDRESS 4000, 192), OccupationCode / OccupationNote = OccupationList(1).OccupationId / Name
	// (T_OCCUPATIONLIST 1000, 189). Migrasi 205.
	"T_DATASCORINGRISKLIST.TOP_RISK_LOCATION": "VARCHAR2(4000)",
	"T_DATASCORINGRISKLIST.OCCUPATION_CODE":   "VARCHAR2(1000)",
	"T_DATASCORINGRISKLIST.OCCUPATION_NOTE":   "VARCHAR2(1000)",
	// Tiket 73 P1 (pola butir 80, keputusan agent): lebar sumber - OCCUPATION_ID = 189; FLAG / kolom MarineHull = POOLDATA.SHIP
	// FLAGNAME / CLASSIFICATION / CONSTRUCTION / DWT / GRT / NRT VARCHAR2(1000) (`DDL\SHIP.txt`, di-SET BrowseShip_RD);
	// OBJECT_NAME 4000 (SHIP.NAME 1000; SetObjectDescHE_DT menyambung tiga medan 500). Migrasi 214.
	"T_ANEKALIST.OCCUPATION_ID":   "VARCHAR2(1000)",
	"T_ANEKALIST.FLAG":            "VARCHAR2(1000)",
	"T_ANEKALIST.OBJECT_NAME":     "VARCHAR2(4000)",
	"T_MARINEHULL.CLASSIFICATION": "VARCHAR2(1000)",
	"T_MARINEHULL.CONSTRUCTION":   "VARCHAR2(1000)",
	"T_MARINEHULL.DWTRT":          "VARCHAR2(1000)",
	"T_MARINEHULL.GRTRT":          "VARCHAR2(1000)",
	"T_MARINEHULL.NRTRT":          "VARCHAR2(1000)",
	// Tiket 73 P2 (migrasi 214): PctShortPeriod angka (K73-6) - rancangan VARCHAR2(30). TYPE_OF_DISCOUNT sudah teks
	// (kolomTeksMenyimpang, butir 68.1).
	"T_COVERAGELIST.PCT_SHORT_PERIOD": "NUMBER",
}

// T_ADDITIONALSHIP - kolom sistem sepola tabel berulang berjalur tunggal di DDL draf
// (ID, IDPEGA, COB_GROUP, PARENT_ID, SEQ_NO, ROW_UID; tanpa PARENT_TABLE/SRC_PATH,
// V-23); DWT/GRT/NRT VARCHAR2(50) sama dengan kolom bernama sama di T_SHIP.
var amandemenTabel = map[string][]kolomSkema{
	"T_ADDITIONALSHIP": {
		{nama: "ID", tipe: "NUMBER", wajib: true},
		{nama: "IDPEGA", tipe: "VARCHAR2(50)"},
		{nama: "COB_GROUP", tipe: "VARCHAR2(20)"},
		{nama: "PARENT_ID", tipe: "NUMBER", wajib: true},
		{nama: "SEQ_NO", tipe: "NUMBER(5)", wajib: true},
		{nama: "ROW_UID", tipe: "VARCHAR2(36)", wajib: true},
		{nama: "DWT", tipe: "VARCHAR2(50)", medan: "DWT"},
		{nama: "GRT", tipe: "VARCHAR2(50)", medan: "GRT"},
		{nama: "NRT", tipe: "VARCHAR2(50)", medan: "NRT"},
		{nama: "ADDITIONAL_SHIP_REF_ID", tipe: "VARCHAR2(50)", turunan: true},
	},
}

// Jalur rancangan (sesudah V-24b: CargoList/PolicyData/Ship -> CargoList/Ship).
// Kembaran Fac Retro (FacRetroList/.../Ship/AdditionalShip) tidak diputuskan - bila ada,
// isinya tetap ke penampung.
var amandemenJalur = []jalurSkema{
	{jalur: "CargoList/Ship/AdditionalShip", tabel: "T_ADDITIONALSHIP", induk: "T_SHIP"},
}

// Penunjuk Idx*/Index* - butir 72 (keputusan work owner 02-10-2026, diteruskan sesi
// `nusantarare-0f`, "setuju"): SATU aturan untuk semua penunjuk - disimpan APA ADANYA
// sebagai kolom teks mentah di tabel barisnya, tidak ditafsirkan. Amandemen V-47:
// penunjuk induk langsung TIDAK dibuang (premis "PARENT_ID menyatakan hal yang sama"
// dibantah data, A68: fixture 103 beda, korpus 8). ⛔ Kolom ini BUKAN pengganti kolom
// FK V-47 (LOCATION_ID, ...) dan tidak dipakai mengisinya; FK tetap kosong (68.4).
//
// 48 kolom = 48 kunci yang teramati di penampung, [terverifikasi] 02-10-2026 dihitung
// atas 5 fixture + 115 contoh korpus; seluruh nilainya bilangan bulat, terpanjang 3
// karakter. Nama: SNAKE_CASE dari FIELD ASLI, pola kolom kunci rancangan yang sudah
// ada (IdxLocation -> IDX_LOCATION, IndexPropertyItem -> INDEX_PROPERTY_ITEM); tipe
// VARCHAR2(50), pola V-6 kolom kode ("teks mentah", bukan NUMBER seperti kolom kunci
// lama). Kelompok menurut lembar Kandidat Hapus; "induk langsung" vs "leluhur" hanya
// pengelompokan laporan (`[dugaan]` dari nama medan), tidak memengaruhi perlakuan.
type kolomPenunjuk struct{ tabel, nama, medan string }

var amandemenPenunjuk = []kolomPenunjuk{
	// R1 indeks-diri (13)
	{tabel: "T_ANEKALIST", nama: "IDX_ANEKA", medan: "IdxAneka"},
	{tabel: "T_CARGOLIST", nama: "IDX_CARGO", medan: "IdxCargo"},
	{tabel: "T_COVERAGELIST", nama: "INDEX_COVERAGE", medan: "IndexCoverage"},
	{tabel: "T_DEDUCTIBLELIST", nama: "INDEX_DEDUCTIBLE", medan: "IndexDeductible"},
	{tabel: "T_FR_ANEKALIST", nama: "IDX_ANEKA", medan: "IdxAneka"},
	{tabel: "T_FR_CARGOLIST", nama: "IDX_CARGO", medan: "IdxCargo"},
	{tabel: "T_FR_COVERAGELIST", nama: "INDEX_COVERAGE", medan: "IndexCoverage"},
	{tabel: "T_FR_DEDUCTIBLELIST", nama: "INDEX_DEDUCTIBLE", medan: "IndexDeductible"},
	{tabel: "T_FR_OCCUPATIONLIST", nama: "IDX_OCCUPATION", medan: "IdxOccupation"},
	{tabel: "T_FR_PROPERTYITEMLIST", nama: "INDEX_PROPERTY_ITEM", medan: "IndexPropertyItem"},
	{tabel: "T_OCCUPATIONLIST", nama: "IDX_OCCUPATION", medan: "IdxOccupation"},
	{tabel: "T_PROPERTYITEMLIST", nama: "INDEX_PROPERTY_ITEM", medan: "IndexPropertyItem"},
	{tabel: "T_VEHICLELIST", nama: "IDX_VEHICLE", medan: "IdxVehicle"},
	// R3 sasaran induk langsung (16)
	{tabel: "T_ADDITIONALCOVERAGE", nama: "INDEX_COVERAGE", medan: "IndexCoverage"},
	{tabel: "T_ANEKALIST", nama: "IDX_OCCUPATION", medan: "IdxOccupation"},
	{tabel: "T_COVERAGELIST", nama: "INDEX_ANEKA", medan: "IndexAneka"},
	{tabel: "T_COVERAGELIST", nama: "INDEX_CARGO", medan: "IndexCargo"},
	{tabel: "T_COVERAGELIST", nama: "INDEX_PROPERTY_ITEM", medan: "IndexPropertyItem"},
	{tabel: "T_COVERAGELIST", nama: "INDEX_VEHICLE_LIST", medan: "IndexVehicleList"},
	{tabel: "T_DEDUCTIBLELIST", nama: "INDEX_COVERAGE", medan: "IndexCoverage"},
	{tabel: "T_FR_CARGOLIST", nama: "IDX_FAC_RETRO", medan: "IdxFacRetro"},
	{tabel: "T_FR_COVERAGELIST", nama: "INDEX_ANEKA", medan: "IndexAneka"},
	{tabel: "T_FR_COVERAGELIST", nama: "INDEX_CARGO", medan: "IndexCargo"},
	{tabel: "T_FR_COVERAGELIST", nama: "INDEX_PROPERTY_ITEM", medan: "IndexPropertyItem"},
	{tabel: "T_FR_DEDUCTIBLELIST", nama: "INDEX_COVERAGE", medan: "IndexCoverage"},
	{tabel: "T_FR_PERSONLIST", nama: "IDX_FAC_RETRO", medan: "IdxFacRetro"},
	{tabel: "T_FR_PROPERTYITEMLIST", nama: "INDEX_PROPERTY", medan: "IndexProperty"},
	{tabel: "T_PROPERTYITEMLIST", nama: "INDEX_PROPERTY", medan: "IndexProperty"},
	{tabel: "T_SPREADINGLIST", nama: "INDEX_COVERAGE", medan: "IndexCoverage"},
	// R3 sasaran leluhur (16)
	{tabel: "T_ADDITIONALCOVERAGE", nama: "INDEX_VEHICLE_LIST", medan: "IndexVehicleList"},
	{tabel: "T_ANEKALIST", nama: "IDX_LOCATION", medan: "IdxLocation"},
	{tabel: "T_COVERAGELIST", nama: "INDEX_LOCATION", medan: "IndexLocation"},
	{tabel: "T_COVERAGELIST", nama: "INDEX_OCCUPATION", medan: "IndexOccupation"},
	{tabel: "T_COVERAGELIST", nama: "INDEX_PROPERTY", medan: "IndexProperty"},
	{tabel: "T_DEDUCTIBLELIST", nama: "INDEX_ANEKA", medan: "IndexAneka"},
	{tabel: "T_DEDUCTIBLELIST", nama: "INDEX_OCCUPATION", medan: "IndexOccupation"},
	{tabel: "T_DEDUCTIBLELIST", nama: "INDEX_PROPERTY", medan: "IndexProperty"},
	{tabel: "T_DEDUCTIBLELIST", nama: "INDEX_PROPERTY_ITEM", medan: "IndexPropertyItem"},
	{tabel: "T_FR_ANEKALIST", nama: "IDX_FAC_RETRO", medan: "IdxFacRetro"},
	{tabel: "T_FR_COVERAGELIST", nama: "INDEX_OCCUPATION", medan: "IndexOccupation"},
	{tabel: "T_FR_COVERAGELIST", nama: "INDEX_PROPERTY", medan: "IndexProperty"},
	{tabel: "T_FR_DEDUCTIBLELIST", nama: "INDEX_PROPERTY_ITEM", medan: "IndexPropertyItem"},
	{tabel: "T_OCCUPATIONLIST", nama: "IDX_LOCATION", medan: "IdxLocation"},
	{tabel: "T_SPREADINGLIST", nama: "INDEX_ANEKA", medan: "IndexAneka"},
	{tabel: "T_SPREADINGLIST", nama: "INDEX_LOCATION", medan: "IndexLocation"},
	// status lembar berselisih (2)
	{tabel: "T_FR_ANEKALIST", nama: "IDX_OCCUPATION", medan: "IdxOccupation"},
	{tabel: "T_FR_DEDUCTIBLELIST", nama: "INDEX_PROPERTY", medan: "IndexProperty"},
	// tidak tercantum di Kandidat Hapus (1)
	{tabel: "T_FR_PERSONLIST", nama: "IDX_PERSON", medan: "IdxPerson"},
}

const tipePenunjuk = "VARCHAR2(50)"

// amandemenDeduction2 - tiket 50 (K50-2, keputusan agent menunggu konfirmasi): `.CurrencyList(n).Policy.Payment.Deduction2`
// (Total Premium Deduction 2, grid tab Premium Deduction, `ProtekDeduction2_Act`) dilipat ke PAY_ seperti saudaranya
// (PAY_COMMISION / PAY_BROKERAGE_FEE, lipatan V-22b) tetapi tidak ada di rancangan: `[terverifikasi]` 0 dari 115 contoh
// `DDL\CONTOH` memuatnya di tingkat CurrencyList.Policy.Payment - seluruh 123 tag Deduction2 ada di ListInstallment
// (101 CurrencyList, 21 OldData, 1 FacRetroList; dihitung dua cara 05-10-2026: pengurai XML = grep tag). Tipe NUMBER =
// saudaranya.
// Migrasi 200.
var amandemenDeduction2 = kolomSkema{nama: "PAY_DEDUCTION2", tipe: "NUMBER", medan: "PayDeduction2", turunan: true}

// amandemenKonstruksiGempa - tiket 52: `.ConstructionEarthquake` (4.2 Construction, `Section\CoverageItem.xml`; isian
// `SearchClassConstructionReas_SQL` = CONSTRUCTION.CONSTRUCTION) ADA di layar tetapi tidak di rancangan T_COVERAGELIST
// (sampel `DDL\CONTOH` tidak memuatnya). Lebar VARCHAR2(500) = kolom konstruksi saudaranya (T_SURROUNDINGRISK
// *_CONSTRUCTION); lebar sumber belum terverifikasi (view CONSTRUCTION atas JSON, kolom tanpa tipe tertulis). Migrasi 203.
var amandemenKonstruksiGempa = kolomSkema{nama: "CONSTRUCTION_EARTHQUAKE", tipe: "VARCHAR2(500)", medan: "ConstructionEarthquake"}

// amandemenScoring - tiket 54: dua medan Scoring Risk yang tidak di rancangan. `.DataScoringRiskList(2).IndexLocation`
// (Dominant Risk, autocomplete `Section\ScoringRisk.xml` -> ObjectNo; 0 sampel `DDL\CONTOH`) VARCHAR2(50) = OBJECT_NO;
// `.LossRatio.LossRatio.PercentLossRatio` (halaman faktor, aturan V-30 -> baris T_SCORING_FACTOR; 2 sampel, sebelumnya
// penampung) NUMBER seperti skor saudaranya. Migrasi 205.
var amandemenScoring = map[string]kolomSkema{
	"T_DATASCORINGRISKLIST": {nama: "INDEX_LOCATION", tipe: "VARCHAR2(50)", medan: "IndexLocation"},
	"T_SCORING_FACTOR":      {nama: "PERCENT_LOSS_RATIO", tipe: "NUMBER", medan: "PercentLossRatio"},
}

// amandemenNBStatus - tiket 57: `pyWorkPage.NBStatus` (teks status kasus yang ditulis Submit tahap marketing) milik halaman
// KERJA, tidak di rancangan; ditaruh di T_GENERAL_POLIS karena T_WORK_POLIS milik premiumlistlife. Medan NBStatus tidak ada
// di halaman OfferFacIn, jadi loader data lama mengisinya kosong. VARCHAR2(1000) - lebar sumber tidak tertulis. Migrasi 208.
var amandemenNBStatus = kolomSkema{nama: "NB_STATUS", tipe: "VARCHAR2(1000)", medan: "NBStatus"}

// amandemenKeUW - tiket 58: `.OfferFacIn.DateToUW` / `.SubmitToUW` (`DataTransform\IsJuniorUW_B.xml`, Accept Team Leader)
// ada di rancangan TETAPI di T_WORK_POLIS (milik premiumlistlife, K-064); ditaruh di T_GENERAL_POLIS, tipe = rancangan
// VARCHAR2(30). Tanpa FIELD ASLI (V-48: medan akar tidak boleh ganda dengan T_WORK_POLIS) - asal = repository (ditulis
// aplikasi, aturan.go kolomRepository). Migrasi 210.
var amandemenKeUW = []kolomSkema{
	{nama: "DATE_TO_UW", tipe: "VARCHAR2(30)"},
	{nama: "SUBMIT_TO_UW", tipe: "VARCHAR2(30)"},
}

// amandemenBangunan - tiket 35 (A110): tiga medan BuildingConstruction yang ADA di layar
// (`Section\ObjectDetails.xml` sel 56-58) tetapi tidak di rancangan; tipe = kolom
// saudaranya (VARCHAR2(50)). Migrasi 186. Digabung ke amandemenKolom di init.
var amandemenBangunan = []kolomSkema{
	{nama: "PARTITION_TYPE", tipe: "VARCHAR2(50)", medan: "PartitionType"},
	{nama: "SUPPORT_WALL_TYPE", tipe: "VARCHAR2(50)", medan: "SupportWallType"},
	{nama: "OTHERS_TYPE", tipe: "VARCHAR2(50)", medan: "OthersType"},
}

// amandemenSekitar - tiket 38 (A130): 18 medan SurroundingRisk yang ADA di layar
// (`Section\RiskAround.xml` sel 9-14/23-28/37-42/51-56, 68, 69) tetapi tidak di rancangan:
// empat sisi x {Occupation, Construction, Distance, Note}, FloodArea, HousekeepingRemark. Kode
// VARCHAR2(50), teks VARCHAR2(500) (pola V-6); Occupation / Note VARCHAR2(1000) = OCCUPATION.OLDID / NAME
// (butir 80); Construction VARCHAR2(500) (nilai standar terpanjang
// 215 bita, `DDL\FrontConstruction.xml`); Distance teks angka VARCHAR2(50) (A129). Migrasi 187.
var amandemenSekitar = func() []kolomSkema {
	var ks []kolomSkema
	for _, s := range []struct{ kolom, medan string }{{"FRONT", "Front"}, {"LEFT", "Left"}, {"BACK", "Back"}, {"RIGHT", "Right"}} {
		ks = append(ks,
			kolomSkema{nama: s.kolom + "_OCCUPATION", tipe: "VARCHAR2(1000)", medan: s.medan + "Occupation"},
			kolomSkema{nama: s.kolom + "_CONSTRUCTION", tipe: "VARCHAR2(500)", medan: s.medan + "Construction"},
			kolomSkema{nama: s.kolom + "_DISTANCE", tipe: "VARCHAR2(50)", medan: s.medan + "Distance"},
			kolomSkema{nama: s.kolom + "_NOTE", tipe: "VARCHAR2(1000)", medan: s.medan + "Note"})
	}
	return append(ks, kolomSkema{nama: "FLOOD_AREA", tipe: "VARCHAR2(50)", medan: "FloodArea"},
		kolomSkema{nama: "HOUSEKEEPING_REMARK", tipe: "VARCHAR2(500)", medan: "HousekeepingRemark"})
}()

// amandemenItem - tiket 39 (A132): enam medan PropertyItem yang ADA di layar
// (`Section\PropertyItemFacIn_Section.xml`: .PropertyYear, .Unit, .Condition, .Year, .NoOfTree,
// .AreaHectar) tetapi tidak di rancangan T_PROPERTYITEMLIST. Nama/tipe = medan bernama sama di tabel
// rancangan lain (YEAR / UNIT VARCHAR2(50), CONDITION VARCHAR2(500)); angka disimpan teks (pola A129).
// Migrasi 188.
var amandemenItem = []kolomSkema{
	{nama: "PROPERTY_YEAR", tipe: "VARCHAR2(50)", medan: "PropertyYear"},
	{nama: "UNIT", tipe: "VARCHAR2(50)", medan: "Unit"},
	{nama: "CONDITION", tipe: "VARCHAR2(500)", medan: "Condition"},
	{nama: "YEAR", tipe: "VARCHAR2(50)", medan: "Year"},
	{nama: "NO_OF_TREE", tipe: "VARCHAR2(50)", medan: "NoOfTree"},
	{nama: "AREA_HECTAR", tipe: "VARCHAR2(50)", medan: "AreaHectar"},
}

// amandemenFEA - tiket 41 (A142): tabel BARU T_FEALIST untuk `LocationList/FEAList` (rancangan tidak punya tabel
// FEA). Kolom sistem pola tabel berulang (T_ADDITIONALSHIP); medan dari `Section\InputFEA_IsUW.xml`; empat
// medan halaman tertanam .DataFEA dilipat ke baris ini (lipatFEA). Migrasi 190.
var amandemenFEA = []kolomSkema{
	{nama: "ID", tipe: "NUMBER", wajib: true},
	{nama: "IDPEGA", tipe: "VARCHAR2(50)"},
	{nama: "COB_GROUP", tipe: "VARCHAR2(20)"},
	{nama: "PARENT_ID", tipe: "NUMBER", wajib: true},
	{nama: "SEQ_NO", tipe: "NUMBER(5)", wajib: true},
	{nama: "ROW_UID", tipe: "VARCHAR2(36)", wajib: true},
	{nama: "APAR", tipe: "VARCHAR2(50)", medan: "APAR"},
	{nama: "SPRINKLER", tipe: "VARCHAR2(50)", medan: "Sprinkler"},
	{nama: "SMOKE_DETECTOR", tipe: "VARCHAR2(50)", medan: "SmokeDetector"},
	{nama: "HYDRANT", tipe: "VARCHAR2(50)", medan: "Hydrant"},
	{nama: "PRIVATE_TRUCK_BRIGADE", tipe: "VARCHAR2(50)", medan: "PrivateTruckBrigade"},
	{nama: "PRIVATE_FIRE_BRIGADE", tipe: "VARCHAR2(50)", medan: "PrivateFireBrigade"},
	{nama: "TEAM_SOP_SAFETY", tipe: "VARCHAR2(50)", medan: "TeamSOPSafety"},
	{nama: "TEAM_SOP_RISK_MANAGEMENT", tipe: "VARCHAR2(50)", medan: "TeamSOPRiskManagement"},
	{nama: "INFO_FEA", tipe: "VARCHAR2(500)", medan: "InfoFEA"},
}

// amandemenKerugian - tiket 42 (A145): lima medan CauseOfLoss yang ADA di layar
// (`Section\InputCauseOfLoss_FacIn.xml`; Amount hanya di grid `Section\CauseOfLoss_FacIn.xml`) tetapi tidak di
// rancangan T_LISTCAUSEOFLOSS (contoh data tidak
// memuatnya). DateOfLoss teks Pega VARCHAR2(30) (pola START_DATE_TIME); teks 500 (pola V-6); uang NUMBER (pola
// rancangan; migrasi NUMBER(38,8), ADR-0016). Migrasi 191.
var amandemenKerugian = []kolomSkema{
	{nama: "DATE_OF_LOSS", tipe: "VARCHAR2(30)", medan: "DateOfLoss"},
	{nama: "LOSS_OBJECT", tipe: "VARCHAR2(500)", medan: "LossObject"},
	{nama: "AMOUNT", tipe: "NUMBER", medan: "Amount"},
	{nama: "PREVENTION_OF_LOSS", tipe: "NUMBER", medan: "PreventionOfLoss"},
	{nama: "CAUSE_OF_LOSS", tipe: "VARCHAR2(500)", medan: "CauseOfLoss"},
}

// amandemenAneka - tiket 73 P1 (migrasi 214): medan form objek Aneka yang ADA di layar
// (`NB FacIn\Section\InputDtlObjectAneka_FacIn.xml`; LimitOfLiability grid okupasi `InputOkupasiAneka_FacIn.xml` L5858; TSIOld
// grid Golf `InputDtlObjectLocation_FacIn.xml` L13568; AlmRiskID di-SET BrowseShip_RD sel L10329) tetapi tidak di rancangan
// T_ANEKALIST (sampel tidak mengisinya). Teks pxTextInput VARCHAR2(50), pxTextArea / nama 500 (pola V-6); angka K73-6 dan uang
// NUMBER (pola rancangan; migrasi NUMBER(38,8)); tanggal teks VARCHAR2(30) (pola START_DATE); Briguna K73-7. TSI_OLD tanpa
// FIELD ASLI: medan TSIOld terkena V-28b (berakhiran Old) - ditulis aplikasi (aturan.go kolomRepository).
var amandemenAneka = []kolomSkema{
	{nama: "NO_KONTRAK", tipe: "VARCHAR2(50)", medan: "NO_KONTRAK"},
	{nama: "NO_KTP", tipe: "VARCHAR2(50)", medan: "NO_KTP"},
	{nama: "NO_NPWP", tipe: "VARCHAR2(50)", medan: "NO_NPWP"},
	{nama: "OTHERS", tipe: "VARCHAR2(500)", medan: "Others"},
	{nama: "INVOICE_NO", tipe: "VARCHAR2(50)", medan: "InvoiceNo"},
	{nama: "MAINTENANCE_PERIOD", tipe: "NUMBER", medan: "MaintenancePeriod"},
	{nama: "NO_OF_TREE", tipe: "NUMBER", medan: "NoOfTree"},
	{nama: "AREA_HECTAR", tipe: "NUMBER", medan: "AreaHectar"},
	{nama: "LIMIT_IDEMNITY_DAY", tipe: "NUMBER", medan: "LimitIdemnityDay"},
	{nama: "LIMIT_IDEMNITY_MONTH", tipe: "NUMBER", medan: "LimitIdemnityMonth"},
	{nama: "LIMIT_OF_LIABILITY", tipe: "NUMBER", medan: "LimitOfLiability"},
	{nama: "TSI_OLD", tipe: "NUMBER"},
	{nama: "FLAG_DELETE", tipe: "VARCHAR2(10)", medan: "FlagDelete"},
	{nama: "ALM_RISK_ID", tipe: "VARCHAR2(100)", medan: "AlmRiskID"},
	{nama: "LOAN_ACCOUNT_NUMBER", tipe: "VARCHAR2(50)", medan: "LoanAccountNumber"},
	{nama: "LOAN_NAME", tipe: "VARCHAR2(500)", medan: "LoanName"},
	{nama: "JOB", tipe: "VARCHAR2(500)", medan: "Job"},
	{nama: "DEPOSIT_ACCOUNT_NO", tipe: "VARCHAR2(50)", medan: "DepositAccountNo"},
	{nama: "ID_CARD_NO", tipe: "VARCHAR2(50)", medan: "IDCardNo"},
	{nama: "INSTITUTION", tipe: "VARCHAR2(500)", medan: "Institution"},
	{nama: "LOAN_TYPE", tipe: "VARCHAR2(50)", medan: "LoanType"},
	{nama: "DATE_OF_BIRTH", tipe: "VARCHAR2(30)", medan: "DateOfBirth"},
	{nama: "SALARY", tipe: "NUMBER", medan: "Salary"},
}

// amandemenKapal / amandemenHE - tiket 73 P1 (migrasi 214): medan layar halaman `.MarineHull` (sel L14598 / L15741 / L16884)
// dan `.VehicleHE` (sel L21528 / L21767 / L26301 / L26540; ObjectIdHE di-SET BrowseObject_RD sel L24970) yang tidak di
// rancangan. Tipe = kolom saudaranya VARCHAR2(50).
var (
	amandemenKapal = []kolomSkema{
		{nama: "LENGTHM", tipe: "VARCHAR2(50)", medan: "LENGTHM"},
		{nama: "BREADTHM", tipe: "VARCHAR2(50)", medan: "BREADTHM"},
		{nama: "DEPTHM", tipe: "VARCHAR2(50)", medan: "DEPTHM"},
	}
	amandemenHE = []kolomSkema{
		{nama: "CHASSIS_NUMBER", tipe: "VARCHAR2(50)", medan: "ChassisNumber"},
		{nama: "ENGINE_NUMBER", tipe: "VARCHAR2(50)", medan: "EngineNumber"},
		{nama: "OBJECT_ID_HE", tipe: "VARCHAR2(50)", medan: "ObjectIdHE"},
	}
)

// amandemenMemorandum - tiket 73 P1 (migrasi 214): tabel BARU T_MEMORANDUMACCEPT untuk halaman `.MemorandumAccept` baris
// Aneka (`NB FacIn\Section\InputDtlMemorandumAccept.xml` 15 FIELD L1068-3633 + SuretyType dari
// `Activity\InputDtlObjectAneka_PreAct.xml`). Kolom sistem pola tabel anak satu halaman (T_MARINEHULL); pxTextInput 500,
// pxTextArea 4000 (keputusan agent). Dua jalur: non-Golf (OccupationList) dan Golf (RiskLocation/AnekaList, kontainer C1).
var (
	amandemenMemorandum = []kolomSkema{
		{nama: "ID", tipe: "NUMBER", wajib: true},
		{nama: "IDPEGA", tipe: "VARCHAR2(50)"},
		{nama: "COB_GROUP", tipe: "VARCHAR2(20)"},
		{nama: "PARENT_ID", tipe: "NUMBER", wajib: true},
		{nama: "SURETY_TYPE", tipe: "VARCHAR2(50)", medan: "SuretyType"},
		{nama: "APPLICATION_NUMBER", tipe: "VARCHAR2(500)", medan: "ApplicationNumber"},
		{nama: "ACCUMULATION", tipe: "VARCHAR2(500)", medan: "Accumulation"},
		{nama: "COLLATERAL", tipe: "VARCHAR2(500)", medan: "Collateral"},
		{nama: "INSURED_SIUP", tipe: "VARCHAR2(500)", medan: "InsuredSiup"},
		{nama: "KEY_PERSON", tipe: "VARCHAR2(500)", medan: "KeyPerson"},
		{nama: "ORIGINAL_RATE", tipe: "VARCHAR2(500)", medan: "OriginalRate"},
		{nama: "PRINCIPALS_OUT_GO", tipe: "VARCHAR2(500)", medan: "PrincipalsOutGo"},
		{nama: "TYPE_OF_WORDING", tipe: "VARCHAR2(500)", medan: "TypeOfWording"},
		{nama: "BUILDING_STATUS", tipe: "VARCHAR2(4000)", medan: "BuildingStatus"},
		{nama: "SUPPORTING_DOCUMENT", tipe: "VARCHAR2(4000)", medan: "SupportingDocument"},
		{nama: "PRINCIPALS_DATA", tipe: "VARCHAR2(4000)", medan: "PrincipalsData"},
		{nama: "CONCLUSIONS_PROPOSALS", tipe: "VARCHAR2(4000)", medan: "ConclusionsProposals"},
		{nama: "PRINCIPALS_BACKGROUND", tipe: "VARCHAR2(4000)", medan: "PrincipalsBackground"},
		{nama: "SUMMARY_FINANCIAL", tipe: "VARCHAR2(4000)", medan: "SummaryFinancial"},
		{nama: "PRINCIPALS_GUARANTEE", tipe: "VARCHAR2(4000)", medan: "PrincipalsGuarantee"},
	}
	jalurMemorandum = []jalurSkema{
		{jalur: "LocationList/Property/RiskLocation/OccupationList/AnekaList/MemorandumAccept", tabel: "T_MEMORANDUMACCEPT", induk: "T_ANEKALIST"},
		{jalur: "LocationList/Property/RiskLocation/AnekaList/MemorandumAccept", tabel: "T_MEMORANDUMACCEPT", induk: "T_ANEKALIST"},
	}
)

// jalurFEA, lipatFEA - tiket 41: jalur FEAList di bawah baris lokasi; halaman .DataFEA dilipat (medan sendiri ikut).
var (
	jalurFEA = jalurSkema{jalur: "LocationList/FEAList", tabel: "T_FEALIST", induk: "T_LOCATIONLIST"}
	lipatFEA = lipatan{"T_FEALIST", "DataFEA"}
)

// init - menggabungkan amandemen ke skema bangkitan. ⛔ Bila workbook kelak sudah
// memuat tabel/kolom/jalur yang sama, penggabungan diam-diam akan menggandakan atau
// menimpanya; karena itu tabrakan = panic saat paket dimuat (amandemen ini harus
// dicabut, bukan ditumpuk).
// Duplikat di dalam amandemen sendiri juga panic.
func init() {
	if _, ada := amandemenTabel[jalurFEA.tabel]; ada {
		panic("loader: amandemen tabel " + jalurFEA.tabel + " ganda")
	}
	amandemenTabel[jalurFEA.tabel] = amandemenFEA
	amandemenJalur = append(amandemenJalur, jalurFEA)
	if _, ada := amandemenTabel["T_MEMORANDUMACCEPT"]; ada {
		panic("loader: amandemen tabel T_MEMORANDUMACCEPT ganda")
	}
	amandemenTabel["T_MEMORANDUMACCEPT"] = amandemenMemorandum
	amandemenJalur = append(amandemenJalur, jalurMemorandum...)
	amandemenKolom["T_ANEKALIST"] = append(amandemenKolom["T_ANEKALIST"], amandemenAneka...)
	amandemenKolom["T_MARINEHULL"] = append(amandemenKolom["T_MARINEHULL"], amandemenKapal...)
	amandemenKolom["T_VEHICLEHE"] = append(amandemenKolom["T_VEHICLEHE"], amandemenHE...)
	if _, ada := lipat[lipatFEA]; ada {
		panic("loader: lipatan T_FEALIST/DataFEA sudah ada")
	}
	lipat[lipatFEA] = aturanLipat{medanSendiri: true}
	amandemenKolom["T_BUILDINGCONSTRUCTION"] = append(amandemenKolom["T_BUILDINGCONSTRUCTION"], amandemenBangunan...)
	amandemenKolom["T_SURROUNDINGRISK"] = append(amandemenKolom["T_SURROUNDINGRISK"], amandemenSekitar...)
	amandemenKolom["T_PROPERTYITEMLIST"] = append(amandemenKolom["T_PROPERTYITEMLIST"], amandemenItem...)
	amandemenKolom["T_LISTCAUSEOFLOSS"] = append(amandemenKolom["T_LISTCAUSEOFLOSS"], amandemenKerugian...)
	amandemenKolom["T_CURRENCYLIST"] = append(amandemenKolom["T_CURRENCYLIST"], amandemenDeduction2)
	amandemenKolom["T_COVERAGELIST"] = append(amandemenKolom["T_COVERAGELIST"], amandemenKonstruksiGempa)
	amandemenKolom["T_GENERAL_POLIS"] = append(amandemenKolom["T_GENERAL_POLIS"], amandemenNBStatus)
	amandemenKolom["T_GENERAL_POLIS"] = append(amandemenKolom["T_GENERAL_POLIS"], amandemenKeUW...)
	for t, k := range amandemenScoring {
		amandemenKolom[t] = append(amandemenKolom[t], k)
	}
	for _, k := range amandemenPenunjuk {
		amandemenKolom[k.tabel] = append(amandemenKolom[k.tabel], kolomSkema{nama: k.nama, tipe: tipePenunjuk, medan: k.medan})
	}
	for t, ks := range amandemenKolom {
		lihat := map[string]bool{}
		for _, k := range ks {
			if lihat[k.nama] {
				panic("loader: amandemen kolom " + t + "." + k.nama + " ganda")
			}
			lihat[k.nama] = true
			for _, ada := range skemaTabel[t] {
				if ada.nama == k.nama {
					panic("loader: amandemen kolom " + t + "." + k.nama + " sudah ada di skema bangkitan")
				}
			}
		}
		skemaTabel[t] = append(skemaTabel[t], ks...)
	}
	for t, ks := range amandemenTabel {
		if _, ada := skemaTabel[t]; ada {
			panic("loader: amandemen tabel " + t + " sudah ada di skema bangkitan")
		}
		skemaTabel[t] = ks
	}
	for _, j := range amandemenJalur {
		for _, ada := range jalurSumber {
			if ada.jalur == j.jalur {
				panic("loader: amandemen jalur " + j.jalur + " sudah ada di skema bangkitan")
			}
		}
	}
	jalurSumber = append(jalurSumber, amandemenJalur...)
	selaraskanWorkPolis()
	for tk, tipe := range amandemenLebar {
		t, k, _ := strings.Cut(tk, ".")
		i := indeksKolom(t, k)
		if i < 0 || skemaTabel[t][i].tipe == tipe {
			panic("loader: lebar " + tk + " tidak dapat diganti " + tipe)
		}
		skemaTabel[t][i].tipe = tipe
	}
}

// selaraskanWorkPolis - butir 76.1 dan 76.4. ⛔ Panic bila kolom sumber tidak ada, nama
// tujuan sudah ada, atau tipe sudah sama: workbook berubah, amandemen ini harus ditinjau.
func selaraskanWorkPolis() {
	for t, peta := range amandemenGabung {
		for lama, baru := range peta {
			i := indeksKolom(t, lama)
			if i < 0 || indeksKolom(t, baru.nama) >= 0 {
				panic("loader: gabung " + t + "." + lama + " -> " + baru.nama + " tidak dapat diterapkan")
			}
			k := skemaTabel[t][i]
			k.nama, k.tipe = baru.nama, baru.tipe
			skemaTabel[t][i] = k
		}
	}
	idKasus := map[string]bool{}
	for _, t := range tabelIDKasus {
		gantiTipe(t, "ID")
		idKasus[t] = true
	}
	induk := map[string]map[string]bool{}
	for _, j := range jalurSumber {
		if induk[j.tabel] == nil {
			induk[j.tabel] = map[string]bool{}
		}
		induk[j.tabel][j.induk] = true
	}
	for tabel, indukTabel := range induk {
		berindukKasus := 0
		for i := range indukTabel {
			if idKasus[i] {
				berindukKasus++
			}
		}
		if berindukKasus == 0 || indeksKolom(tabel, "PARENT_ID") < 0 {
			continue // T_GENERAL_POLIS: induk T_WORK_POLIS, berbagi PK tanpa PARENT_ID
		}
		if berindukKasus != len(indukTabel) {
			panic("loader: " + tabel + " berinduk campuran - PARENT_ID tidak dapat satu tipe")
		}
		gantiTipe(tabel, "PARENT_ID")
		anakIDKasus = append(anakIDKasus, tabel)
	}
	sort.Strings(anakIDKasus)
}

func indeksKolom(t, nama string) int {
	for i, k := range skemaTabel[t] {
		if k.nama == nama {
			return i
		}
	}
	return -1
}

func gantiTipe(t, nama string) {
	i := indeksKolom(t, nama)
	if i < 0 || skemaTabel[t][i].tipe == tipeIDKasus {
		panic("loader: tipe " + t + "." + nama + " tidak dapat diganti " + tipeIDKasus)
	}
	skemaTabel[t][i].tipe = tipeIDKasus
}

// Riwayat penunjuk Idx*/Index*: P7 (butir 70) membuang sebagian; butir 71 mencabutnya
// (semua ke penampung); butir 72 menyimpan semuanya di kolom teks mentah
// (amandemenPenunjuk). Dasar butir 71-72:
//
//	A68 - induk langsung: [terverifikasi] 02-10-2026, dua cara (mesin vs pengurai Python),
//	      nilai tidak selalu = posisi baris induk (fixture 103 dari 219 beda; korpus 8),
//	      membantah premis V-47 "PARENT_ID sudah menyatakan hal yang sama".
//	A67 - R1 indeks-diri: nilai tidak selalu = posisi baris sendiri (fixture 0 dari 227
//	      beda; korpus 39), jadi "sudah diwakili kolom urutan baris" tidak berlaku umum.
//
// Arti penunjuk tidak ditafsirkan. penunjukKandidat (bangkitan lembar Kandidat Hapus)
// tetap dibangkitkan sebagai catatan rancangan; mesin tidak membacanya (hanya uji).
