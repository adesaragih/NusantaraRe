package models

// ⛔ BERKAS INI BELUM PUNYA SATU PUN PEMANGGIL DI LUAR UJI, DAN ITU
// DISENGAJA - TETAPI HARUS DIBACA SEBAGAI UTANG, BUKAN SEBAGAI SELESAI.
//
// KunciKelompokDokumen / IDDokumenBaru / BolehSimpanBarisDokumen / PerluHapusDiPenyimpanan belum dipanggil kode produksi mana pun.
// Jalur unggah dan hapus dokumen menuntut penyambungan ke Google Storage
// (`InsertGoogleStorage_Act`, `DeleteGoogleStorage_Act`), dan itu menuntut
// persetujuan manusia. Aturannya ditiru lebih dulu supaya yang DAPAT
// diputuskan tidak menunggu yang tidak dapat.
//
// ⚠️ Bentuk ini PERSIS "rute tanpa pemanggil" yang tiga kali menjadi cacat
// di modul ini - backend hijau, layar hijau, fiturnya tidak ada. Bedanya
// satu dan hanya satu: di sini ketiadaan pemanggil DINYATAKAN, di sana ia
// tidak. Bila Anda membaca ini dan penyambungan penyimpanan sudah ada,
// maka pekerjaan yang tersisa adalah MEMANGGILNYA - bukan menulis ulang
// aturannya.

// Aturan lahirnya satu baris dokumen — kelompok Dokumen.
//
// Untuk apa berkas ini: ketika sebuah berkas dilampirkan, Pega menurunkan
// TIGA hal dan menggerbangi SATU. Ketiganya murni; gerbangnya yang menentukan
// urutan, dan urutan itu yang paling mudah dibalik tanpa sadar.
//
// Sumber, dibaca sebagai pohon 27-09-2026:
//
//	`Activity/SaveAttachLife.xml`
//	  langkah 1 b294 ULANG(EMBEDDED b1677) atas `dragDropFileUpload.pxResults`
//	    1.2 b483  `Primary.DOCUMENT` =
//	              `@if(Primary.DOCUMENT == "","DL-"+@pxReplaceAllViaRegex(
//	               @CurrentDateTime(),"[^0-9]",""),Primary.DOCUMENT)`  b595-596
//	              `Local.IDDoc` = `Primary.DOCUMENT`                    b615-616
//	    1.7 b1467 `Call InsertDocument_Act` dengan
//	              `IDPEGA=pyWorkPage.pzInsKey` · `NAMAFILE=.pyFileName`
//	              `MIME=.pyFileType` · `KATEGORI_1=Local.IDDoc`
//	              `KATEGORI_2=.pyCategory` · `BASE64=.pyFileSource`
//	              `NOAKSEP=""` · `NOPREKAS=""` · `PAYMENTDATE=""`
//	  langkah 3 b1830 `Commit`
//
//	`Activity/InsertDocument_Act.xml`
//	  2 b520  `Java` prasyarat b586 `Param.MIME==""` WhenTrue=2 LANJUT
//	  3 b627  `.ID` = `@CurrentDate("yyyyMMddhhmmssSSS","Asia/Jakarta")` b647-648
//	          `.TANGGAL` = `@CurrentDateTime()` b701-702
//	          `.MIME` = `@toLowerCase(Param.MIME)` b781-782
//	          `.pxCreateOperator` = `pxRequestor.pyUserIdentifier` b939-940
//	  4 b1023 `Call InsertGoogleStorage_Act`
//	  5 b1185 `Obj-Save` prasyarat b1283 `NewDocument.T_STORAGE_ID==""`
//	          WhenTrue=3 LEWATI
//
//	`Activity/DeleteDocument_Act.xml`
//	  2 b385  `Call DeleteGoogleStorage_Act` prasyarat b472
//	          `DeleteDocument.T_STORAGE_ID==""` WhenTrue=3 LEWATI
//	  3 b513  `Obj-Delete` — TANPA prasyarat
//
// Dibaca sesudah: mimedokumen.go.

import (
	"regexp"
	"strings"
	"time"
)

// AwalanKunciDokumen adalah awalan kunci kelompok dokumen satu peserta.
//
// ⛔ VERBATIM `SaveAttachLife.xml` b596. Ia BUKAN pengenal baris dokumen
// melainkan kunci KELOMPOK: satu peserta punya satu nilai, dan seluruh
// dokumennya menyimpannya di `KATEGORI_1`. Itulah yang disaring
// `LoadDocumentLife_ACT` (`Field .KATEGORI_1` = `Value .DOCUMENT`).
const AwalanKunciDokumen = "DL-"

// bukanAngka meniru `@pxReplaceAllViaRegex(…,"[^0-9]","")`.
var bukanAngka = regexp.MustCompile(`[^0-9]`)

// KunciKelompokDokumen mengembalikan kunci kelompok dokumen peserta.
//
// ⛔ Yang SUDAH ADA tidak pernah ditimpa - itu arti `@if(Primary.DOCUMENT ==
// "", …, Primary.DOCUMENT)` b596. Menimpanya pada lampiran kedua akan
// memutus lampiran pertama dari pesertanya: baris lama tetap menyimpan kunci
// lama di `KATEGORI_1`, dan saringan layar hanya mencari kunci yang baru.
// Dokumen itu tidak hilang dari basis data; ia hilang dari layar, yang lebih
// buruk karena tidak ada yang tahu.
func KunciKelompokDokumen(sudahAda string, saat time.Time) string {
	if strings.TrimSpace(sudahAda) != "" {
		return sudahAda
	}
	return AwalanKunciDokumen + bukanAngka.ReplaceAllString(
		saat.Format("2006-01-02 15:04:05.000"), "")
}

// bentukIDDokumen meniru `@CurrentDate("yyyyMMddhhmmssSSS","Asia/Jakarta")`.
//
// ⚠️ `hh` pada pola Pega/Java adalah jam 12-JAM (1..12). Ditiru apa adanya:
// pengenal yang dihasilkan pukul 14:05 dan pukul 02:05 karena itu BERTABRAKAN
// pada milidetik yang sama. Itu cacat warisan, bukan pilihan kami, dan ia
// DICATAT alih-alih diam-diam diperbaiki - memperbaikinya di sini membuat
// pengenal baris baru berbeda bentuk dari pengenal baris lama, dan keduanya
// hidup di kolom yang sama.
const bentukIDDokumen = "20060102030405.000"

// IDDokumenBaru menurunkan pengenal baris dokumen dari waktu.
//
// ⛔ Zona waktunya `Asia/Jakarta` VERBATIM b648, bukan zona server. Pengenal
// yang bergeser tujuh jam ketika server dipindah adalah pengenal yang tidak
// dapat diurutkan bersama pengenal lama.
func IDDokumenBaru(saat time.Time, jakarta *time.Location) string {
	return bukanAngka.ReplaceAllString(saat.In(jakarta).Format(bentukIDDokumen), "")
}

// BolehSimpanBarisDokumen menjawab prasyarat b1283.
//
// ⛔ URUTANNYA: unggah DULU, baris basis data KEMUDIAN. `Obj-Save` DILEWATI
// ketika `T_STORAGE_ID` kosong, sehingga di sistem lama tidak pernah ada
// baris dokumen yang berkasnya tidak ada.
//
// ⚠️ Membaliknya - menulis baris lebih dulu lalu mengunggah - menghasilkan
// baris yang menunjuk berkas yang tidak pernah naik. Layar akan
// menampilkannya sebagai dokumen yang ada, tautannya akan gagal, dan tidak
// ada yang dapat membedakannya dari gangguan jaringan sesaat.
func BolehSimpanBarisDokumen(tStorageID string) bool {
	return strings.TrimSpace(tStorageID) != ""
}

// PerluHapusDiPenyimpanan menjawab prasyarat b472.
//
// ⛔ Dan `Obj-Delete` b513 berjalan TANPA prasyarat: barisnya dihapus
// sekalipun tidak ada berkas di penyimpanan. Arah yang benar - baris yatim
// tanpa berkas tidak berguna bagi siapa pun.
func PerluHapusDiPenyimpanan(tStorageID string) bool {
	return strings.TrimSpace(tStorageID) != ""
}
