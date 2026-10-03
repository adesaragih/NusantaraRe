package loader

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

// init - menggabungkan amandemen ke skema bangkitan. ⛔ Bila workbook kelak sudah
// memuat tabel/kolom/jalur yang sama, penggabungan diam-diam akan menggandakan atau
// menimpanya; karena itu tabrakan = panic saat paket dimuat (amandemen ini harus
// dicabut, bukan ditumpuk).
// Duplikat di dalam amandemen sendiri juga panic.
func init() {
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
