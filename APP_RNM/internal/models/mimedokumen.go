package models

// Jenis MIME dokumen — tabel keputusan `DecisionTable/GetMimeType.xml`.
//
// Untuk apa berkas ini: ketika pemanggil tidak menyebutkan jenis berkas,
// Pega menurunkannya dari EKSTENSI nama berkas lewat tabel keputusan itu.
//
// Sumber, dibaca 27-09-2026:
//
//	`Activity/InsertDocument_Act.xml`
//	  langkah 2 b520 `Java`, prasyarat b586 `Param.MIME==""`
//	                  `WhenTrue=2` LANJUT / `WhenFalse=3` LEWATI
//	  langkah 3 b781-782 `.MIME` = `@toLowerCase(Param.MIME)`
//
// ⛔ Prasyaratnya berarti tabel ini dipakai HANYA ketika pemanggil tidak
// menyebutkan jenisnya. `SaveAttachLife` b1467 meneruskan `.pyFileType`;
// bila peramban mengisinya, nilai ITU yang dipakai - bukan tebakan dari
// ekstensi. Membalik urutannya akan menimpa jenis yang sudah benar dengan
// tebakan yang lebih miskin.
//
// ⚠️ Keempat puluh delapan baris di bawah DISALIN DARI EKSPOR, tidak
// dikarang dan tidak dilengkapi. Beberapa di antaranya tampak ganjil
// (`et`, `oxps`, `lnk`), dan justru itu sebabnya ia disalin apa adanya:
// daftar MIME "yang masuk akal" akan berbeda dari daftar yang sistem lama
// pakai, dan bedanya baru terlihat ketika sebuah berkas ditolak.

import "strings"

// MimeBawaan adalah hasil untuk ekstensi yang tidak ada di tabel.
//
// ⛔ VERBATIM `DecisionTable/GetMimeType.xml` b89 - nilai `otherwise` tabel
// itu. Berkas tak dikenal TIDAK ditolak; ia diberi jenis umum, dan itu
// keputusan sistem lama, bukan kelalaian.
const MimeBawaan = "application/octet-stream"

// petaMime memetakan ekstensi (huruf kecil, tanpa titik) ke jenis MIME.
//
// Nomor baris menunjuk pasangan kondisi/hasil di berkas tabelnya.
var petaMime = map[string]string{
	"pdf":  "application/pdf",                                                                   // b290 / b417
	"jpg":  "image/jpeg",                                                                        // b291 / b418
	"jpeg": "image/jpeg",                                                                        // b292 / b419
	"png":  "image/png",                                                                         // b293 / b420
	"jfif": "image/jpeg",                                                                        // b294 / b421
	"gif":  "image/gif",                                                                         // b295 / b422
	"heic": "image/heic",                                                                        // b296 / b423
	"heif": "image/heif",                                                                        // b297 / b424
	"tif":  "image/tiff",                                                                        // b298 / b425
	"tiff": "image/tiff",                                                                        // b299 / b426
	"bmp":  "image/bmp",                                                                         // b300 / b427
	"webp": "image/webp",                                                                        // b301 / b428
	"doc":  "application/msword",                                                                // b302 / b429
	"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",           // b303 / b430
	"xls":  "application/vnd.ms-excel",                                                          // b304 / b431
	"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",                 // b305 / b432
	"xlsb": "application/vnd.ms-excel.sheet.binary.macroEnabled.12",                             // b306 / b433
	"xlt":  "application/vnd.ms-excel",                                                          // b307 / b434
	"xltx": "application/vnd.openxmlformats-officedocument.spreadsheetml.template",              // b308 / b435
	"xlsm": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.macroEnabled.12", // b309 / b436
	"ppt":  "application/vnd.ms-powerpoint",                                                     // b310 / b437
	"pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",         // b311 / b438
	"pps":  "application/vnd.ms-powerpoint",                                                     // b312 / b439
	"ppsx": "application/vnd.openxmlformats-officedocument.presentationml.slideshow",            // b313 / b440
	"rtf":  "application/rtf",                                                                   // b314 / b441
	"txt":  "text/plain",                                                                        // b315 / b442
	"eml":  "message/rfc822",                                                                    // b316 / b443
	"zip":  "application/zip",                                                                   // b317 / b444
	"7z":   "application/x-7z-compressed",                                                       // b318 / b445
	"rar":  "application/x-rar-compressed",                                                      // b319 / b446
	"msg":  "application/vnd.ms-outlook",                                                        // b320 / b447
	"csv":  "text/csv",                                                                          // b321 / b448
	"odt":  "application/vnd.oasis.opendocument.text",                                           // b322 / b449
	"ods":  "application/vnd.oasis.opendocument.spreadsheet",                                    // b323 / b450
	"odp":  "application/vnd.oasis.opendocument.presentation",                                   // b324 / b451
	"url":  "application/internet-shortcut",                                                     // b325 / b452
	"lnk":  "application/x-ms-shortcut",                                                         // b326 / b453
	"html": "text/html",                                                                         // b327 / b454
	"htm":  "text/html",                                                                         // b328 / b455
	"wps":  "application/vnd.ms-works",                                                          // b329 / b456
	"svg":  "image/svg+xml",                                                                     // b330 / b457
	"oxps": "application/oxps",                                                                  // b331 / b458
	"mp4":  "video/mp4",                                                                         // b332 / b459
	"zipx": "application/zip",                                                                   // b333 / b460
	"wav":  "audio/wav",                                                                         // b334 / b461
	"et":   "application/et",                                                                    // b335 / b462
	"mht":  "message/rfc822",                                                                    // b336 / b463
	"avi":  "video/x-msvideo",                                                                   // b337 / b464
}

// JumlahBarisMime adalah cacah baris tabel keputusan itu.
//
// ⚠️ Dikunci uji. Baris yang HILANG dari tabel tidak berbunyi sendiri: berkas
// yang jenisnya hilang akan diam-diam menjadi `application/octet-stream`, dan
// tidak ada yang tahu ia pernah punya jenis sendiri.
const JumlahBarisMime = 48

// MimeDariNamaFile menurunkan jenis MIME dari nama berkas.
//
// ⛔ Ekstensi diambil dari titik TERAKHIR. "laporan.2026.pdf" berekstensi
// "pdf", bukan "2026.pdf"; nama berkas bertitik banyak lumrah pada lampiran
// klaim.
//
// ⛔ Nama berkas TANPA titik, atau yang berakhir dengan titik, mendapat
// MimeBawaan - bukan galat. Tabel keputusannya punya `otherwise`, dan
// menolak berkas yang sistem lama terima berarti menutup pintu yang terbuka.
func MimeDariNamaFile(nama string) string {
	titik := strings.LastIndex(nama, ".")
	if titik < 0 || titik == len(nama)-1 {
		return MimeBawaan
	}
	// ⚠️ Huruf kecil di KEDUA sisi. Tabelnya berkunci huruf kecil, dan
	// "LAPORAN.PDF" adalah nama berkas yang sangat biasa.
	if m, ada := petaMime[strings.ToLower(nama[titik+1:])]; ada {
		return m
	}
	return MimeBawaan
}

// MimeDokumen memilih antara jenis yang DISEBUT pemanggil dan tebakan tabel.
//
// ⛔ Ini padanan prasyarat b586, dan urutannya penting: yang disebut
// pemanggil MENANG. Tabel hanya mengisi yang kosong.
//
// ⚠️ Hasilnya selalu huruf kecil - padanan `@toLowerCase(Param.MIME)` b782.
// Tanpa itu "APPLICATION/PDF" dan "application/pdf" menjadi dua jenis
// berbeda di kolom yang sama.
func MimeDokumen(disebut, namaFile string) string {
	if strings.TrimSpace(disebut) != "" {
		return strings.ToLower(strings.TrimSpace(disebut))
	}
	return MimeDariNamaFile(namaFile)
}
