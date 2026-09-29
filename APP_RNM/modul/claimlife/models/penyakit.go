package models

// Pencarian diagnosa — kelompok Medis.
//
// Untuk apa berkas ini: `DISEASE_LIFE` berisi **97.586 baris**. Setiap
// pencarian atasnya harus BERBATAS, dan batasnya bukan karangan kami: rule
// Pega menyebutkannya angka demi angka.
//
// Sumber, dibaca 27-09-2026:
//
//	`ReportDefinition/BrowseDiseaseLife_RD.xml`
//	  b659 `pyMaxRecords`        500
//	  b514 / b747 `pyPageSize`   50
//	  b535 / b754 `pyLogic`      `A AND B`
//	       A  `.ICD_Code` `Contains` `Param.ICD_Code`   b543-552
//	       B  `.Disease`  `Contains` `Param.Disease`    b556-569
//	  b598 urut `.Number` ASC (`pySortOrder` 1)
//	  medan: `.Number` (label `ID`), `.Disease`, `.ICD_Code`
//
//	`Activity/SearchDiagnose_act.xml`
//	  b254-255 `SearchDiagnose.CARI1` = `@toUpperCase(SearchDiagnose.CARI1)`
//	  b301-302 `SearchDiagnose.CARI2` = `@toUpperCase(SearchDiagnose.CARI2)`
//
//	`Section/Diagnose_Section.xml`
//	  b1645 `ICD_Code` <- `SearchDiagnose.CARI1`
//	  b1651 `Disease`  <- `SearchDiagnose.CARI2`
//
// ⛔ CARI1 adalah KODE ICD dan CARI2 adalah NAMA penyakit - bukan
// sebaliknya. Pemetaannya dibaca dari `pyRDParams` b1657-1658, bukan
// ditebak dari namanya; "CARI1/CARI2" tidak menyebutkan apa pun.

import "strings"

// BatasBarisPenyakit adalah `pyMaxRecords` b659.
//
// ⛔ 500, dan angka itu dari rule - bukan dari selera kami. Ia satu-satunya
// yang berdiri antara layar dan 97.586 baris.
const BatasBarisPenyakit = 500

// UkuranHalamanPenyakit adalah `pyPageSize` b514.
const UkuranHalamanPenyakit = 50

// Penyakit adalah satu baris `DISEASE_LIFE`.
//
// ⚠️ Ketiga medannya persis medan yang report definition itu pilih. Nomor
// dibawa sebagai TEKS: ia pengenal, bukan bilangan yang dihitung
// (ADR-U-0022), dan "007" bukan "7".
type Penyakit struct {
	Nomor   string `json:"nomor"`
	Nama    string `json:"nama"`
	KodeICD string `json:"kodeIcd"`
}

// KriteriaPenyakit adalah dua kata kunci pencarian, sesudah dinormalkan.
type KriteriaPenyakit struct {
	// KodeICD adalah `Param.ICD_Code` <- `SearchDiagnose.CARI1`.
	KodeICD string
	// Nama adalah `Param.Disease` <- `SearchDiagnose.CARI2`.
	Nama string
}

// NormalkanKriteriaPenyakit meniru kedua `@toUpperCase` b255 dan b302.
//
// ⛔ HURUF BESAR di kedua sisi, dan itu bukan kosmetik: pencarian `Contains`
// di Oracle peka huruf, sehingga tanpa ini "diabetes" tidak menemukan
// "DIABETES" - dan pemakai akan menyimpulkan penyakitnya tidak ada.
//
// ⚠️ Spasi di tepi dibuang. Rule aslinya TIDAK membuangnya, dan itu
// penyimpangan sadar: `Contains " A00"` di Pega mencari spasi pun, sehingga
// satu spasi tak sengaja membuat pencarian gagal tanpa sebab yang terlihat.
// Arah selisihnya aman - ia hanya MELEBARKAN hasil, tidak pernah
// menyembunyikan baris yang Pega tampilkan.
func NormalkanKriteriaPenyakit(kodeICD, nama string) KriteriaPenyakit {
	return KriteriaPenyakit{
		KodeICD: strings.ToUpper(strings.TrimSpace(kodeICD)),
		Nama:    strings.ToUpper(strings.TrimSpace(nama)),
	}
}

// Kosong menyatakan tidak satu pun kata kunci diisi.
//
// ⚠️ Dengan keduanya kosong, `Contains ""` di Pega cocok dengan SEMUA baris -
// dan `pyMaxRecords` 500 yang menahannya. Ditiru apa adanya: pencarian
// kosong SAH, dan batasnya yang bekerja. Melarangnya akan menutup jalan yang
// di sistem lama terbuka *(menelusuri daftar tanpa tahu kata kuncinya)*.
func (k KriteriaPenyakit) Kosong() bool { return k.KodeICD == "" && k.Nama == "" }

// BatasPenyakit menjepit batas yang diminta pemanggil ke rentang yang sah.
//
// ⛔ Nol atau negatif menjadi UkuranHalamanPenyakit, bukan "tanpa batas".
// Permintaan tanpa batas atas tabel 97.586 baris adalah permintaan yang
// memuat seluruh tabel ke dalam memori satu proses.
func BatasPenyakit(diminta int) int {
	if diminta <= 0 {
		return UkuranHalamanPenyakit
	}
	if diminta > BatasBarisPenyakit {
		return BatasBarisPenyakit
	}
	return diminta
}
