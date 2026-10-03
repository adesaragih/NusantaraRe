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

// init - menggabungkan amandemen ke skema bangkitan. ⛔ Bila workbook kelak sudah
// memuat tabel/kolom/jalur yang sama, penggabungan diam-diam akan menggandakan atau
// menimpanya; karena itu tabrakan = panic saat paket dimuat (amandemen ini harus
// dicabut, bukan ditumpuk).
// Duplikat di dalam amandemen sendiri juga panic.
func init() {
	amandemenKolom["T_BUILDINGCONSTRUCTION"] = append(amandemenKolom["T_BUILDINGCONSTRUCTION"], amandemenBangunan...)
	amandemenKolom["T_SURROUNDINGRISK"] = append(amandemenKolom["T_SURROUNDINGRISK"], amandemenSekitar...)
	amandemenKolom["T_PROPERTYITEMLIST"] = append(amandemenKolom["T_PROPERTYITEMLIST"], amandemenItem...)
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
