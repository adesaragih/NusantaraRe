package models

import "time"

// KolomCsvObjek - 61 kolom template CSV objek FIRE (services.KolomCsvObjek menunjuk ke sini; tiket 66 K66-3), urutan berkas;
// juga kolom penampung T_NB_CSV_BARIS (awalan C_, migrasi 257).
var KolomCsvObjek = []string{"NO", "NO_OFFER_SLIP", "CLASS_OF_BUSINESS", "BEGIN_DATE", "END_DATE", "INSURED_NAME", "SOB_NAME",
	"CEDING_CO_NAME", "MARKETING_NAME", "OCCUPATION_CODE", "OCCUPATION_NAME", "CLASS_OF_CONSTRUCTION", "TYPE_ADDRESS", "ADDRESS",
	"BUILDING_NO", "POSTAL_CODE", "KELURAHAN", "OBJECT_TYPE", "INTEREST", "TOP_RISK", "ITEM_TYPE", "ITEM_TYPE_1", "ITEM_TYPE_2",
	"ITEM_TYPE_3", "ITEM_TYPE_4", "ITEM_TYPE_5", "ITEM_TYPE_6", "TSI_OBJECT_ITEM", "TSI_OBJECT_ITEM_1", "TSI_OBJECT_ITEM_2",
	"TSI_OBJECT_ITEM_3", "TSI_OBJECT_ITEM_4", "TSI_OBJECT_ITEM_5", "TSI_OBJECT_ITEM_6", "ACCUMULATION_CODE", "ACCUMULATION_NAME",
	"ADJUSTABLE_PCT", "CURRENCY", "INDEMNITY_MONTH", "LOL", "RATE_FLEXAS", "NET_RATE_FLEXAS", "DISCOUNT_RATE_FLEXAS", "RATE_RSMDCC",
	"NET_RATE_RSMDCC", "DISCOUNT_RATE_RSMDCC", "RATE_FWTWD", "NET_RATE_FWTWD", "DISCOUNT_RATE_FWTWD", "RATE_EQ", "NET_RATE_EQ",
	"DISCOUNT_RATE_EQ", "RATE_OTHERS", "NET_RATE_OTHERS", "DISCOUNT_RATE_OTHERS", "KODE_POS", "KECAMATAN", "KABUPATEN_KOTA",
	"PROVINSI", "SUMBER_KODE_POS", "TINGKAT_KEYAKINAN"}

// Jenis importer pekerjaan unggah CSV objek (T_NB_CSV_UNGGAH.JENIS, migrasi 258) - dipilih menurut Class of Business case seperti
// Pega `ChooseActivityMappingCSV_Act`: IsFire -> fire (InsertUploadFire_act), IsMarineCargo -> cargo (InsertUploadMarineCargo_act).
const (
	JenisCsvFire  = "fire"
	JenisCsvCargo = "cargo"
	// JenisCsvAneka - tiket 73 P6: IsAneka -> InsertUploadAneka_act (ChooseActivityMappingCSV_Act langkah 5).
	JenisCsvAneka = "aneka"
)

// Template CSV objek case (GET kasus `jenisCsv`, tiket 73 P6) - Pega `NB FacIn\Activity\SetDownloadMasterCSV.xml` langkah 2-6
// (Property-Set InputData.CARI44 berurutan, yang terakhir benar menang) `[terverifikasi]`: IsAneka TemplateCSVAneka, IsMarineHull
// TemplateCSVMarineHull, IsContractorsPlantMachinery TemplateCSVAneka_CPM, IsFire TemplateCSVFire, IsMarineCargo
// TemplateCSVMarineCargo; tidak satu pun = "" (Golf tidak termasuk IsAneka). Importer ketiga template Aneka sama
// (ChooseActivityMappingCSV_Act langkah 5 IsAneka -> InsertUploadAneka_act).
const (
	TemplateCsvAneka      = "aneka"
	TemplateCsvMarineHull = "marinehull"
	TemplateCsvAnekaCPM   = "cpm"
)

// KolomCsvCargo - 18 judul yang dibaca InsertUploadMarineCargo_act / ProtectUploadMarineCargo_act lewat pxUploadCSVResults (per
// NAMA judul, urutan berkas bebas) `[terverifikasi]` (tiket 74 P6). Urutan di sini = urutan kanonik penampung (kolom C_<judul>,
// migrasi 258; CURRENCY / TOP_RISK dipakai bersama kolom FIRE), BUKAN urutan template.
var KolomCsvCargo = []string{"TRADING", "DEPARTURE", "DESTINATION", "GOODS", "GOODS_NAME", "PACKING", "CONVEYENCE",
	"CONVEYENCE_NAME", "COVERAGE_IDX", "COVERAGE_NAME", "TSI", "CURRENCY", "AGENT", "DATASHIP", "INVOICE", "ETD", "BLNUMBER", "TOP_RISK"}

// KolomCsvAneka - 35 judul yang dibaca InsertUploadAneka_act / ProtectUploadAneka_Act lewat pxUploadCSVResults (per NAMA judul)
// `[terverifikasi]` (tiket 73 P6, sensus dua cara). Urutan: 31 judul template Marine Hull (`DDL\templatecsvmarinehull.csv`, urutan
// berkas) lalu 4 judul yang dibaca tetapi tidak ada di template itu. Kolom penampung C_<judul> (12 dipakai bersama kolom FIRE /
// cargo, 23 migrasi 250).
var KolomCsvAneka = []string{"NO", "TYPE_ADDRESS", "ADDRESS", "BUILDING_NO", "POSTAL_CODE", "COUNTRY", "OCCUPATION_NAME",
	"OBJECT_DESCRIPTION", "SECTION", "CURRENCY", "TSI", "YEAR", "PLACE_OF_MANUFACTURE", "OBJECT_TYPE", "NRT", "GRT", "DWT", "FLAGNAME",
	"CONSTRUCTION", "CLASSIFICATION", "COVERAGE_NAME", "CALCULATE_METHOD", "DAY", "PCT_LOADING", "PCT_SCALE", "PCT_RATE", "PREMIUM",
	"PCT_DISCOUNT", "DISCOUNT", "PCT_LOL", "LOL", "SERIALNO", "CHASSISNUMBER", "ENGINENUMBER", "KELURAHAN"}

// JumlahJudulMarineHull - 31 judul pertama KolomCsvAneka = template Marine Hull (judul wajib K73-60).
const JumlahJudulMarineHull = 31

// KolomCsvJenis - kolom penampung pekerjaan berjenis `jenis` (kosong / tak dikenal = fire, pekerjaan sebelum 258).
func KolomCsvJenis(jenis string) []string {
	switch jenis {
	case JenisCsvCargo:
		return KolomCsvCargo
	case JenisCsvAneka:
		return KolomCsvAneka
	}
	return KolomCsvObjek
}

// UnggahCsv - satu pekerjaan unggah CSV objek tanpa batas (tiket 75, T_NB_CSV_UNGGAH migrasi 257). Status / Fase lihat
// services/csvunggah.go; GalatDiBaris 0 = tidak ada.
type UnggahCsv struct {
	ID, CaseID, Status, Fase, NamaBerkas                        string
	Ukuran                                                      int64
	TotalBaris, BarisDiproses, BarisTersimpan, KelompokTerakhir int
	JumlahGalat, JumlahPesan, GalatDiBaris                      int
	NoOfferSlip, BeginDate, EndDate, GalatProses, DibuatOleh    string
	Dibuat, Diubah                                              time.Time
	// Jenis - JenisCsvFire / JenisCsvCargo (migrasi 258); Kolom - judul yang dikenali urut berkas (nil = urutan KolomCsvJenis).
	Jenis string
	Kolom []string
	// Basi - DIUBAH (detak) lebih tua dari BatasDetakUnggahMenit menurut jam basis data (dihitung di SQL, bebas zona waktu).
	Basi bool
}

// BatasDetakUnggahMenit - pekerjaan berjalan tanpa detak selama ini dianggap terputus (tiket 75 K75-2).
const BatasDetakUnggahMenit = 15

// BarisUnggahCsv - satu baris penampung (T_NB_CSV_BARIS): No = nomor baris data 1..n, Isi = sel urut KolomCsvJenis(jenis pekerjaan).
type BarisUnggahCsv struct {
	No  int
	Isi []string
}

// PesanUnggahCsv - satu galat / pesan per baris (T_NB_CSV_PESAN); NoBaris 0 = tingkat berkas.
type PesanUnggahCsv struct {
	NoBaris int
	Teks    string
}
