package models

// Untuk apa berkas ini: PENGGOLONG JENIS USAHA - port `DecisionTable/
// BusinessType_DeT` (tiket 06, spec §5.5, AC 19-22) dan `When/IsLife`
// (AC 67).
//
// `[terverifikasi]` Sumber: `docs/INVENTARIS-XML.md` bab 10. Baris tabel di
// bawah DIBANGKITKAN dari korpus oleh `docs/alat` (bukan disalin tangan):
// 36 baris, dua kolom penguji `.Quotation.GroupPanel` dan
// `.Quotation.BusinessOldId` (tipe kolom `text`, operator `=`), 128 kode
// BusinessOldId seluruhnya unik.
//
// Semantik tabel keputusan Pega yang ditiru:
//
//	sel kosong                  -> kolom itu TIDAK diuji (apa saja cocok)
//	daftar-ATAU (pyOrConditions) -> nilai cocok salah satu anggota; pyRowNum
//	                               berbasis NOL dipetakan ke baris pyRowNum+1
//	pyEvaluateAllRows = no      -> BERHENTI di baris pertama yang cocok (AC 19)
//	pyDefaultResult = "UNKNOWN" -> tak ada yang cocok (AC 20)
//
// ⚠️ Sel `"26"   ` di baris 30 memuat spasi DI LUAR tanda kutip. Sebagai
// ekspresi Pega nilainya `"26"` - spasi di luar literal hanyalah spasi
// ekspresi. Kode 26 cocok di baris 30, bukan jatuh ke baris 32.
//
// ⛔ Keterangan rule `When` diabaikan seluruhnya (P23, AC 22): yang ditiru
// adalah syarat yang DIJALANKAN.

// BarisPenggolong adalah satu baris `BusinessType_DeT`.
type BarisPenggolong struct {
	Baris            int
	GroupPanel       string
	BusinessOldID    string
	UjiBusinessOldID bool
	Atau             []string
	Hasil            string
}

// JenisUsahaTakDikenal adalah `pyDefaultResult` tabel - VERBATIM.
const JenisUsahaTakDikenal = "UNKNOWN"

// TabelPenggolong adalah ke-36 baris `BusinessType_DeT`, urut baris.
var TabelPenggolong = []BarisPenggolong{
	{Baris: 1, GroupPanel: "001", BusinessOldID: "", UjiBusinessOldID: false, Atau: nil, Hasil: "Medicare"},
	{Baris: 2, GroupPanel: "002", BusinessOldID: "03", UjiBusinessOldID: true, Atau: nil, Hasil: "PA"},
	{Baris: 3, GroupPanel: "003", BusinessOldID: "", UjiBusinessOldID: true, Atau: []string{"36", "37", "38", "39", "40", "41", "42", "43", "44", "45", "46", "47", "48", "49", "50", "51", "52", "C2", "C3", "C4", "C5", "C6", "C7", "C8", "C9", "D1", "D2", "D3", "D4", "D5", "D6", "D7", "D8", "D9", "F1", "F2", "F3", "F4", "F5"}, Hasil: "Bonding"},
	{Baris: 4, GroupPanel: "003", BusinessOldID: "", UjiBusinessOldID: true, Atau: []string{"95", "96", "97", "98", "99", "A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "B4", "B5", "B6", "B8", "B9", "C1"}, Hasil: "BondingKBG"},
	{Baris: 5, GroupPanel: "003", BusinessOldID: "20", UjiBusinessOldID: true, Atau: nil, Hasil: "HE"},
	{Baris: 6, GroupPanel: "003", BusinessOldID: "", UjiBusinessOldID: true, Atau: []string{"05", "55"}, Hasil: "MarineHull"},
	{Baris: 7, GroupPanel: "003", BusinessOldID: "21", UjiBusinessOldID: true, Atau: nil, Hasil: "Glass"},
	{Baris: 8, GroupPanel: "003", BusinessOldID: "", UjiBusinessOldID: true, Atau: []string{"15", "24", "89", "90", "91", "E7", "E9"}, Hasil: "Liability"},
	{Baris: 9, GroupPanel: "003", BusinessOldID: "16", UjiBusinessOldID: true, Atau: nil, Hasil: "AllRisk"},
	{Baris: 10, GroupPanel: "003", BusinessOldID: "06", UjiBusinessOldID: true, Atau: nil, Hasil: "AviationHull"},
	{Baris: 11, GroupPanel: "003", BusinessOldID: "14", UjiBusinessOldID: true, Atau: nil, Hasil: "Burglary"},
	{Baris: 12, GroupPanel: "003", BusinessOldID: "07", UjiBusinessOldID: true, Atau: nil, Hasil: "Car"},
	{Baris: 13, GroupPanel: "003", BusinessOldID: "08", UjiBusinessOldID: true, Atau: nil, Hasil: "Ear"},
	{Baris: 14, GroupPanel: "003", BusinessOldID: "09", UjiBusinessOldID: true, Atau: nil, Hasil: "ElectronicEquipment"},
	{Baris: 15, GroupPanel: "003", BusinessOldID: "19", UjiBusinessOldID: true, Atau: nil, Hasil: "Fidelity"},
	{Baris: 16, GroupPanel: "003", BusinessOldID: "27", UjiBusinessOldID: true, Atau: nil, Hasil: "GolfInsurance"},
	{Baris: 17, GroupPanel: "003", BusinessOldID: "17", UjiBusinessOldID: true, Atau: nil, Hasil: "CIT"},
	{Baris: 18, GroupPanel: "003", BusinessOldID: "18", UjiBusinessOldID: true, Atau: nil, Hasil: "CIS"},
	{Baris: 19, GroupPanel: "003", BusinessOldID: "11", UjiBusinessOldID: true, Atau: nil, Hasil: "Boiler"},
	{Baris: 20, GroupPanel: "003", BusinessOldID: "93", UjiBusinessOldID: true, Atau: nil, Hasil: "Workmen"},
	{Baris: 21, GroupPanel: "003", BusinessOldID: "10", UjiBusinessOldID: true, Atau: nil, Hasil: "MBD"},
	{Baris: 22, GroupPanel: "003", BusinessOldID: "12", UjiBusinessOldID: true, Atau: nil, Hasil: "ContractorsPM"},
	{Baris: 23, GroupPanel: "003", BusinessOldID: "", UjiBusinessOldID: true, Atau: []string{"79", "80", "81", "82", "83", "84", "85", "86", "87", "88"}, Hasil: "CustomBond"},
	{Baris: 24, GroupPanel: "003", BusinessOldID: "28", UjiBusinessOldID: true, Atau: nil, Hasil: "FireStyle1"},
	{Baris: 25, GroupPanel: "003", BusinessOldID: "B3", UjiBusinessOldID: true, Atau: nil, Hasil: "LandRig"},
	{Baris: 26, GroupPanel: "003", BusinessOldID: "", UjiBusinessOldID: false, Atau: nil, Hasil: "Aneka"},
	{Baris: 27, GroupPanel: "004", BusinessOldID: "", UjiBusinessOldID: false, Atau: nil, Hasil: "MarineCargo"},
	{Baris: 28, GroupPanel: "005", BusinessOldID: "", UjiBusinessOldID: false, Atau: nil, Hasil: "Travel"},
	{Baris: 29, GroupPanel: "006", BusinessOldID: "78", UjiBusinessOldID: true, Atau: nil, Hasil: "OilGas"},
	{Baris: 30, GroupPanel: "006", BusinessOldID: "", UjiBusinessOldID: true, Atau: []string{"22", "26", "29", "57"}, Hasil: "FireStyle1"},
	{Baris: 31, GroupPanel: "006", BusinessOldID: "", UjiBusinessOldID: true, Atau: []string{"01", "34", "59"}, Hasil: "FireStyle2"},
	{Baris: 32, GroupPanel: "006", BusinessOldID: "", UjiBusinessOldID: false, Atau: nil, Hasil: "Fire"},
	{Baris: 33, GroupPanel: "007", BusinessOldID: "02", UjiBusinessOldID: true, Atau: nil, Hasil: "MBUCar"},
	{Baris: 34, GroupPanel: "007", BusinessOldID: "60", UjiBusinessOldID: true, Atau: nil, Hasil: "MBUMotorCycle"},
	{Baris: 35, GroupPanel: "009", BusinessOldID: "", UjiBusinessOldID: true, Atau: []string{"L1", "L2", "L3", "L4", "L5", "L6", "L7", "L8", "L9", "L10", "L11", "L12", "L13", "L14", "L15", "L16", "L17", "L18", "L19", "L20", "L21"}, Hasil: "Life"},
	{Baris: 36, GroupPanel: "008", BusinessOldID: "", UjiBusinessOldID: false, Atau: nil, Hasil: "Medicare"},
}

// GolongkanJenisUsaha menjalankan `BusinessType_DeT` atas dua masukan teks.
//
// Pembandingannya TEKS PERSIS (`pyColumnDataType` = text): `"006"` dan `"6"`
// berbeda, dan itu disengaja - spec-penyimpanan ID-16.
func GolongkanJenisUsaha(groupPanel, businessOldID string) string {
	for _, b := range TabelPenggolong {
		if b.GroupPanel != groupPanel {
			continue
		}
		if b.UjiBusinessOldID && !cocokBusinessOldID(b, businessOldID) {
			continue
		}
		return b.Hasil
	}
	return JenisUsahaTakDikenal
}

func cocokBusinessOldID(b BarisPenggolong, v string) bool {
	if b.BusinessOldID != "" && b.BusinessOldID == v {
		return true
	}
	for _, a := range b.Atau {
		if a == v {
			return true
		}
	}
	return false
}

// KodeLiniJiwa adalah ke-16 nilai `When/IsLife` - syarat yang DIJALANKAN
// (`pyLogic` A OR ... OR P atas `pyWorkPage.Quotation.BusinessOldId`).
//
// ⚠️ Keterangan penampil rule itu menyebut `BusinessCode = "10164"` milik
// halaman OfferFacIn; ia TIDAK dijalankan (P23). AC 67.
var KodeLiniJiwa = []string{
	"L1", "L2", "L3", "L4", "L5", "L6", "L7", "L8",
	"L9", "L10", "L11", "L12", "L13", "L14", "L15", "L16",
}

// AdalahLiniJiwa = `When/IsLife`.
//
// ⚠️ Pemanggilnya di korpus ini hanya dua aktivitas rantai Fac In
// (`ProtectPremiPolicy_Act`, `ProtectShareCedant_Act`) yang tidak terjangkau,
// dan `Protection_Act` varian Treaty yang berkasnya tidak ada (INVENTARIS bab
// 1.2). AC 67 menuntut penggolongnya dimigrasi apa adanya; ia tersedia di sini
// tanpa pemanggil sampai varian Treaty itu diperoleh.
func AdalahLiniJiwa(businessOldID string) bool {
	for _, k := range KodeLiniJiwa {
		if k == businessOldID {
			return true
		}
	}
	return false
}
