package loader

import (
	"fmt"
	"regexp"
	"strings"
)

// Aturan khusus - HANYA yang tertulis sebagai keputusan V-xx / K-xxx di workbook
// rancangan (`Claude outputs\Tabel-Flat-per-Grup-Bisnis.xlsx` lembar BACA-INI) dan
// `08-flat\BAHAN-SPEC-PEMUATAN.md`. Selebihnya mesin generik membaca skema_gen.go.

// cabangDibuang - subpohon yang tidak menghasilkan baris maupun kolom, dengan
// dasarnya. Dicocokkan pada NAMA kunci, di kedalaman mana pun - termasuk kembaran
// cabang Fac Retro-nya (V-28: "beserta kembaran Fac Retro-nya").
var cabangDibuang = map[string]string{
	"OldData":                "V-19 salinan kerja (bukan versi sebelumnya, V-43a)",
	"FacRetro":               "V-25 cabang FacRetro tanpa List",
	"Parameters":             "V-34 parameter layar",
	"OutGoList":              "V-34 parameter layar",
	"TotalTSIList":           "V-28 tabel penjumlahan (K-067)",
	"TotalTSIPremiGrossList": "V-28 tabel penjumlahan (K-067)",
	"TotalTSIPremiSpreadRNM": "V-28 tabel penjumlahan (K-067)",
	"ViewSuggest":            "V-31 riwayat akseptasi tetap di POOLDATA.HISTORYAKSEPTASIPRODUCTION",
	"OldOfferedPayment":      "V-44 salinan sebelum endorsement",
}

// Alasan medan daun yang dibuang (kunci Diagnostik.Dibuang).
const (
	alasanMeta      = "px*/pz* metadata ekspor Pega"
	alasanPy        = "py* - rancangan tanpa satu pun kolom py* (diukur dari lembar Kolom)"
	alasanSufiksOld = "V-28b medan berakhiran Old (nilai sebelum endorsement)"
	alasanEDMOld    = "V-45 medan berpola EDMOld*"
)

// metadata - medan px/pz/py: tidak dihitung sebagai medan data di mana pun.
func metadata(medan string) bool {
	al := alasanDaun(medan)
	return al == alasanMeta || al == alasanPy
}

// alasanDaun - alasan medan daun dibuang, atau "" bila medan itu dipetakan.
// ⚠️ Urutan cek: EDMOld* lebih dulu (EDMOldPremi tidak berakhiran Old).
func alasanDaun(medan string) string {
	switch {
	case strings.HasPrefix(medan, "px"), strings.HasPrefix(medan, "pz"):
		return alasanMeta
	case strings.HasPrefix(medan, "py"):
		return alasanPy
	case strings.HasPrefix(medan, "EDMOld"):
		return alasanEDMOld
	case strings.HasSuffix(medan, "Old"):
		return alasanSufiksOld
	}
	return ""
}

// lipatan - halaman yang DILIPAT ke baris induknya (bukan tabel sendiri). Kunci:
// tabel baris induk + nama halaman.
//
//	V-39  PolicyData di akar dilipat ke T_GENERAL_POLIS ("seluruh isi PolicyData")
//	V-22b PolicyData/Payment di akar -> awalan Pay (PAY_INSTALLMENT …)
//	V-22  CurrencyList/Policy/Payment -> awalan Pay; V-22a ListInstallment-nya menjadi
//	      anak T_CURRENCYLIST. Kembaran T_FR_CURRENCYLIST: rancangannya memuat sembilan
//	      kolom PAY_ yang sama, jadi Payment-nya dilipat juga.
//	V-24b CargoList/PolicyData dilipat ke T_CARGOLIST (SailDate; V-26 BLNumber,
//	      InvoiceNumber); Ship dan LC menjadi anak T_CARGOLIST
type lipatan struct{ tabel, halaman string }

// aturanLipat - dalam: halaman di DALAM lipatan yang ikut dilipat, dengan awalan
// medannya. medanSendiri: medan daun halaman itu sendiri ikut dilipat.
//
// ⛔ Policy: medanSendiri false. V-22 melipatnya karena "nol kolom terisi" - tidak
// pernah memutuskan ke mana medannya pergi. [terverifikasi] 02-10-2026 di 115 contoh:
// FacOfferList/CurrencyList/Policy membawa TSI sendiri (15 unsur) di samping TSI
// CurrencyList-nya; melipatnya menimpa satu dengan yang lain. Medan itu kini
// terhitung TakTerpetakan, bukan menimpa.
type aturanLipat struct {
	dalam        map[string]string
	medanSendiri bool
	// khusus - medan halaman itu yang punya kolom sendiri di baris induk walau
	// medanSendiri false (amandemen P6, butir 70).
	khusus map[string]string
}

var lipat = map[lipatan]aturanLipat{
	{"T_GENERAL_POLIS", "PolicyData"}: {dalam: map[string]string{"Payment": "Pay"}, medanSendiri: true},
	{"T_CURRENCYLIST", "Policy"}:      {dalam: map[string]string{"Payment": "Pay"}},
	{"T_FR_CURRENCYLIST", "Policy"}:   {dalam: map[string]string{"Payment": "Pay"}, khusus: map[string]string{"TSI": "POLICY_TSI"}},
	{"T_CARGOLIST", "PolicyData"}:     {medanSendiri: true},
	{"T_FR_CARGOLIST", "PolicyData"}:  {medanSendiri: true},
}

// gantiNama - ruas jalur mentah yang bernama lain di rancangan.
//
//	V-17b VehicleList/Occupation -> T_OCCUPATIONLIST (jalur rancangan VehicleList/OccupationList)
//	V-33  ASMCoverage dilebur ke CoverageList (PA); lokasinya di bawah PersonList `[dugaan]`
//	      dari jalur T_COVERAGELIST "PersonList/CoverageList | Life PA" - belum ada fixture PA
var gantiNama = map[[2]string]string{
	{"VehicleList", "Occupation"}: "OccupationList",
	{"PersonList", "ASMCoverage"}: "CoverageList",
}

// medanID - V-49: medan Pega bernama `ID` pada tabel ini diberi nama baru supaya
// tidak bentrok dengan kolom sistem ID.
var medanID = map[string]string{
	"T_CURRENCY":      "CURRENCY_REF_ID",
	"T_FR_CURRENCY":   "CURRENCY_REF_ID",
	"T_SHIP":          "SHIP_REF_ID",
	"T_FR_SHIP":       "SHIP_REF_ID",
	"T_RETROLIST":     "RETRO_REF_ID",
	"T_GENERAL_POLIS": "SOURCE_ID",
	// Amandemen butir 70 (amandemen.go), pola V-49 yang sama.
	"T_CURRENCYLIST":   "CURRENCY_REF_ID",
	"T_ADDITIONALSHIP": "ADDITIONAL_SHIP_REF_ID",
}

// medanGanda - BAHAN §0: di JSON kode mata uang CurrencyList hanya ada di medan
// `Name` (di XML juga nama elemennya). CURRENCY_CODE (FIELD ASLI rekaan
// `CurrencyCode`) diisi dari Name; NAME tetap diisi dari Name.
//
// ⚠️ Bunyi V-22 "Kolom ID pada CurrencyList berisi kode mata uang" TIDAK cocok dengan
// data: [terverifikasi] 02-10-2026, 116 unsur CurrencyList akar di 115 contoh - Name
// kode tiga huruf 116/116, ID terisi 5/116 dan seluruhnya angka lima digit (rujukan
// master, sebangun Currency/ID V-49). ID karena itu TIDAK dipetakan ke CURRENCY_CODE;
// ia terhitung TakTerpetakan.
var medanGanda = map[string]map[string]string{
	"T_CURRENCYLIST": {"Name": "CURRENCY_CODE"},
}

// V-30 ScoringRisk: tiap faktor di bawah DataScoringRiskList (langsung, atau di
// bawah satu halaman kelompok) menjadi satu baris T_SCORING_FACTOR; pasangan
// ChechBoxN / ScoreN menjadi baris T_SCORING_OPTION dengan OPTION_NO = N.
// `[dugaan]` FACTOR_GROUP dan FACTOR_NAME diisi NAMA halaman apa adanya - peta "23
// faktor" yang V-30 sebut (lembar ScoringRisk) tidak ada di kedua workbook.
const (
	tabelSkoring  = "T_DATASCORINGRISKLIST"
	ruasFaktor    = "ScoringFactor"
	ruasOpsi      = "ScoringOption"
	tabelFaktor   = "T_SCORING_FACTOR"
	tabelOpsi     = "T_SCORING_OPTION"
	medanTandaCek = "Checked"
)

var polaOpsi = regexp.MustCompile(`^(ChechBox|Score)([0-9]+)$`)

// V-16: 18 nilai BusinessType yang terbaca (BAHAN §4). Langkah 5 hanya menerima
// Life dan PA; enam belas sisanya di langkah 5 berarti bentuk dan nilai tidak sepakat.
var businessTypeDikenal = map[string]bool{
	"Aneka": true, "AviationHull": true, "Bonding": true, "BondingKBG": true, "CustomBond": true,
	"ElectronicEquipment": true, "FireStyle1": true, "FireStyle2": true, "GolfInsurance": true,
	"HE": true, "LandRig": true, "Liability": true, "Life": true, "MarineCargo": true,
	"MarineHull": true, "MBD": true, "MBUCar": true, "PA": true,
}

// kedalamanMaks - kedalaman rowdata bersarang terbesar yang pernah terukur (BAHAN §1).
const kedalamanMaks = 8

// jenisWork - JENIS_WORK yang sah (lembar Kolom T_WORK_POLIS: "NB / RNW / EDM").
var jenisWork = map[string]bool{"NB": true, "RNW": true, "EDM": true}

// Kolom yang TIDAK diisi Flatten - pengisinya repository (sumber di luar
// DATA_JSON) atau belum ada aturan tertulisnya. Dipakai TestSetiapKolomBerasal.
var (
	// kolomRepository - diisi repository (tiket 24), bukan dari dokumen.
	kolomRepository = map[string]string{
		"ID":                           "surrogate dari sekuens; Flatten memberi Baris.Kunci",
		"PARENT_ID":                    "ID baris Baris.Induk",
		"T_GENERAL_POLIS.OLD_POLIS_ID": "K-071 J-5: hanya RENEWAL; fase 1 NULL seluruhnya",
		"T_GENERAL_POLIS.PROD_KE":      "POOLDATA.JSON_POLIS.PRODKE (VARCHAR2(5) -> NUMBER, K-071 J-3)",
		"T_WORK_POLIS.NOURUT":          "V-48: baris terakhir POOLDATA.HISTORYAKSEPTASIPRODUCTION",
		"T_WORK_POLIS.PUTARAN":         "V-48: baris terakhir POOLDATA.HISTORYAKSEPTASIPRODUCTION",
		"T_WORK_POLIS.STS_KONVERSI":    "V-48: POOLDATA.JSON_POLIS",
		"T_WORK_POLIS.TGL_KONVERSI":    "V-48: POOLDATA.JSON_POLIS",
		// Butir 76 (amandemen.go): T_WORK_POLIS = tabel yang ada (K-064).
		"T_WORK_POLIS.ID":          "butir 76.1: = NO_WORK (pyID, mis. NB-184351) - pengenal work tabel yang ada, bukan sekuens",
		"T_GENERAL_POLIS.ID":       "butir 76.1: berbagi PK = T_WORK_POLIS.ID (K-064 relasi 52)",
		"T_WORK_POLIS.LINI":        "butir 76.2: 'FAC' (penanda lini K-064)",
		"T_WORK_POLIS.POSITION":    "V-48: baris terakhir POOLDATA.HISTORYAKSEPTASIPRODUCTION (POSISI rancangan, digabung butir 76.4)",
		"T_WORK_POLIS.STATUS_WORK": "V-48a: tidak punya sumber di bahan mana pun - belum terverifikasi (STATUS_PROSES rancangan, digabung butir 76.4)",
		"T_WORK_POLIS.TGL_CREATE":  "V-48: POOLDATA.JSON_POLIS (TGL_INPUT rancangan, digabung butir 76.4)",
		"T_WORK_POLIS.CREATE_OP":   "V-48: POOLDATA.JSON_POLIS (USERNAME rancangan, digabung butir 76.4)",
	}
	// kolomFKV47 - V-47 "IdxLocation jadi LOCATION_ID ke T_LOCATIONLIST, IndexProperty
	// jadi PROPERTY_ID, IndexAneka jadi ANEKA_ID, dan seterusnya". ⛔ Belum diisi:
	// medan sumber per kolom dan cakupan "nomor posisi" (di dalam induk mana) tidak
	// tertulis untuk 29 kolom ini; DDL draf memberinya VARCHAR2(50) tanpa REFERENCES
	// sedangkan ID bertipe NUMBER. Medan Idx*/Index* sumbernya disimpan apa adanya di
	// kolom penunjuk teks mentah (butir 72, amandemen.go) - bukan pengisi FK ini.
	kolomFKV47 = map[string]bool{"ANEKA_ID": true, "CARGO_ID": true, "LOCATION_ID": true, "PROPERTY_ID": true,
		"PROPERTY_ITEM_ID": true, "VEHICLE_ID": true, "COVERAGE_ID": true, "FAC_RETRO_ID": true}
	alasanFKV47 = "V-47: medan sumber dan cakupan posisi belum tertulis per kolom - belum terverifikasi"
	// kolomTanpaRumus - V-37 TSI_TOP_RISK: "tidak ada di ekspor Pega", rumusnya tidak tertulis.
	kolomTanpaRumus = map[string]string{
		"TSI_TOP_RISK": "V-37: ditambahkan work owner, tidak ada di ekspor, rumus belum tertulis - belum terverifikasi",
	}
	// alasanWarisTanpaPembawa - pewarisan K-063 (c) yang leluhurnya tidak punya kolom kode:
	// T_LOCATIONLIST hanya punya CURRENCY_ID (angka rujukan master, bukan kode ISO -
	// lembar Audit Mata Uang). Kolomnya selalu UNKNOWN dan terhitung.
	alasanWarisTanpaPembawa = "K-063 (c) diwarisi dari %s, yang tidak punya CURRENCY_CODE - selalu UNKNOWN, terhitung"
)

// kolomTeksMenyimpang - ⛔ PENYIMPANGAN SADAR dari DDL draf (keputusan work owner
// 02-10-2026, butir 68.1): delapan kolom yang DDL beri NUMBER tetapi isinya teks
// disimpan sebagai TEKS APA ADANYA - "9%" tidak ditafsirkan, koma desimal tidak
// dinormalkan, spasi di ujung PctLimit tidak dipangkas. Nilai peta = baris DDL draf
// (`08-flat\DDL-tabel-flat-draf.sql`) yang disimpangi. Panjang VARCHAR2-nya
// keputusan tiket 23 (ditahan).
var kolomTeksMenyimpang = map[string]string{
	"T_GENERAL_POLIS.PPN_CHECK":         "L95: PPN_CHECK NUMBER",
	"T_COVERAGELIST.TYPE_OF_DISCOUNT":   "L206: TYPE_OF_DISCOUNT NUMBER",
	"T_QUOTATIONDATA.EDM_CHARGE_FEE":    "L233: EDM_CHARGE_FEE NUMBER",
	"T_QUOTATIONDATA.SHARE_OF_CEDING":   "L272: SHARE_OF_CEDING NUMBER",
	"T_PERSONLIST.MASTER_RATE_COVERAGE": "L513: MASTER_RATE_COVERAGE NUMBER",
	"T_TABLEOFLIMIT.PCT_LIMIT":          "L914: PCT_LIMIT NUMBER",
	"T_FR_PRINTRISLIP.WARR_PAYMENT":     "L1151: WARR_PAYMENT NUMBER",
	"T_FR_TABLEOFLIMIT.PCT_LIMIT":       "L1519: PCT_LIMIT NUMBER",
}

// kodeDariCurrency - keputusan work owner 02-10-2026: CURRENCY_CODE tabel (kunci)
// diambil dari NAME baris halaman anak `Currency` (nilai = tabel baris halaman itu),
// karena FIELD ASLI `Name` milik baris itu sendiri tidak ada di JSON maupun XML.
// Butir 68.3 untuk dua tabel bisnis; butir 69 menambah kembaran T_FR_* "supaya
// konsisten" (anak Currency-nya baris T_FR_CURRENCY).
//
// T_PROPERTY ikut "bila leluhurnya punya halaman Currency" - [terverifikasi]
// 02-10-2026 (ralat butir 69: rumusan lama "nol halaman Currency di bawah
// LocationList" terlalu luas - di bawah PropertyItemList/DeductibleList/AnekaList ada
// halaman Currency): nol halaman Currency sebagai ANAK LANGSUNG LocationList maupun
// LocationList/Property (0/4 fixture, 0/297 korpus); Currency di akar ada (3/5
// fixture, 41/115 korpus) tetapi Name-nya tidak pernah terisi. Syaratnya tidak
// terpenuhi, T_PROPERTY tetap mewarisi (selalu UNKNOWN).
var kodeDariCurrency = map[string]string{
	"T_COVERAGELIST": "T_CURRENCY", "T_ANEKALIST": "T_CURRENCY",
	"T_FR_COVERAGELIST": "T_FR_CURRENCY", "T_FR_ANEKALIST": "T_FR_CURRENCY",
}

// kolomNamaCurrency - kolom kode di baris halaman Currency (V-49: medan ID-nya
// menjadi CURRENCY_REF_ID).
const kolomNamaCurrency = "NAME"

// medanDiselamatkan - medan di dalam cabang yang dibuang yang TIDAK boleh ikut
// hilang: nilainya masuk penampung (ADR-0023) sampai kolomnya ada.
//
//	ViewSuggest.IsCedingConfirm - K-069 (7b) "Keputusan work owner: kolom sendiri"
//	(salinan repo `docs/00-KEPUTUSAN-WORK-OWNER.md` L3788-3798). Butir 70 P4: kolom
//	sendiri di POOLDATA.HISTORYAKSEPTASIPRODUCTION di samping POSISI - tafsiran yang
//	dipilih work owner; tabel lama, MENUNGGU DBA (nama + tipe), jadi tetap di penampung.
//	Medan ViewSuggest lain tetap dibuang (V-31: riwayat di tabel lama itu).
var medanDiselamatkan = map[string]map[string]bool{
	"ViewSuggest": {"IsCedingConfirm": true},
}

// kodeTakDiketahui - sentinel K-069 usulan 10 / ADR-0006 untuk baris lama tanpa mata uang.
const kodeTakDiketahui = "UNKNOWN"

// Kolom sistem yang diisi Flatten.
const (
	kolomIDPega      = "IDPEGA"
	kolomCOB         = "COB_GROUP"
	kolomSeq         = "SEQ_NO"
	kolomRowUID      = "ROW_UID"
	kolomTabelInduk  = "PARENT_TABLE"
	kolomJalurSumber = "SRC_PATH"
	kolomJenisWork   = "JENIS_WORK"
	kolomNoWork      = "NO_WORK"
	kolomKodeUang    = "CURRENCY_CODE"
)

// Asal isi satu kolom - dipakai TestSetiapKolomBerasal untuk menagih bahwa tidak
// satu pun dari 1.329 kolom tanpa asal tertulis.
const (
	asalMedan      = "medan dokumen"
	asalFlatten    = "diisi Flatten (aturan sistem/turunan)"
	asalRepository = "diisi repository (sumber di luar DATA_JSON)"
	asalKosong     = "sengaja tidak diisi (belum terverifikasi)"
)

var kolomSistemFlatten = map[string]bool{kolomIDPega: true, kolomCOB: true, kolomSeq: true, kolomRowUID: true,
	kolomTabelInduk: true, kolomJalurSumber: true}

// kolomV30 - kolom yang diisi aturan V-30 (bukan medan bernama sama di dokumen).
var kolomV30 = map[string]bool{"T_SCORING_FACTOR.FACTOR_GROUP": true, "T_SCORING_FACTOR.FACTOR_NAME": true,
	"T_SCORING_OPTION.OPTION_NO": true, "T_SCORING_OPTION.CHECK_BOX": true, "T_SCORING_OPTION.SCORE": true}

// asalKolom - asal isi kolom k tabel t, dan alasannya bila tidak dari medan dokumen.
func asalKolom(t string, k kolomSkema) (asal, alasan string) {
	tk := t + "." + k.nama
	// Kunci tabel.kolom lebih dulu: T_WORK_POLIS.ID (butir 76.1) mengalahkan "ID" umum.
	if a, ada := kolomRepository[tk]; ada {
		return asalRepository, a
	}
	if a, ada := kolomRepository[k.nama]; ada {
		return asalRepository, a
	}
	switch {
	case kolomSistemFlatten[k.nama]:
		return asalFlatten, "kolom sistem"
	case t == "T_WORK_POLIS" && (k.nama == kolomJenisWork || k.nama == kolomNoWork):
		return asalFlatten, "dipecah dari IDPEGA (lembar Kolom T_WORK_POLIS)"
	case kolomV30[tk]:
		return asalFlatten, "V-30"
	case medanID[t] == k.nama:
		return asalFlatten, "V-49 medan ID"
	}
	for l, a := range lipat {
		for _, kol := range a.khusus {
			if l.tabel == t && kol == k.nama {
				return asalFlatten, "amandemen P6: medan " + l.halaman + " ke kolom induk"
			}
		}
	}
	for _, kol := range medanGanda[t] {
		if kol == k.nama {
			return asalFlatten, "BAHAN §0 kode mata uang dari Name"
		}
	}
	if sumber, ada := warisMataUang[t]; ada && k.nama == kolomKodeUang {
		for _, ks := range skemaTabel[sumber] {
			if ks.nama == kolomKodeUang {
				return asalFlatten, "K-063 (c) diwarisi dari " + sumber
			}
		}
		return asalKosong, fmt.Sprintf(alasanWarisTanpaPembawa, sumber)
	}
	if k.turunan {
		switch {
		case k.kolomPay():
			return asalMedan, "V-22/V-22b medan Payment berawalan Pay"
		case kolomFKV47[k.nama]:
			return asalKosong, alasanFKV47
		}
		if a, ada := kolomTanpaRumus[k.nama]; ada {
			return asalKosong, a
		}
		return "", "kolom turunan tanpa aturan"
	}
	if k.medan != "" {
		return asalMedan, ""
	}
	return "", "kolom tanpa medan dan tanpa aturan"
}
