package unggah

// Pengenal berkas di penyimpanan - `ImageID`, butir be.
//
// ⛔ RALAT 28-09-2026 ATAS BACAAN KAMI SENDIRI. Ronde pertama memakai
// pengenal DOKUMEN (cap waktu `@CurrentDate("yyyyMMddhhmmssSSS")`,
// `InsertDocument_Act.xml` b648) sebagai `IMAGEID` kartu penyimpanan, dengan
// alasan "keduanya lahir di saat yang sama". Keliru: `IMAGEID` punya rule
// pembangkitnya SENDIRI, dan rule itu ada di korpus sejak awal:
//
//	RDBList/GenerateImageID_SQL.xml  b85-b88
//	  SELECT STANDARD_HASH(
//	           'ASMPP' || TO_CHAR(SYSTIMESTAMP, 'YYYYMMDDHH24MISSFF9') || SYS_GUID(),
//	           'MD5') AS "InsertDoc.ImageID"
//	  FROM DUAL
//
// dipanggil `Activity/InsertGoogleStorage_Act.xml` b2226, hasilnya mengalir
// ke `Param.ImageID` b2784-2785 lalu ke `Insert_T_Storage_SQL.xml` b94
// sebagai kolom `IMAGEID`.
//
// ⚠️ Bedanya BUKAN kosmetik. Pengenal dokumen dapat ditebak - ia cap waktu -
// sedangkan `IMAGEID` memuat `SYS_GUID()`. Memakai yang pertama sebagai kunci
// penyimpanan berarti siapa pun yang tahu KAPAN sebuah berkas diunggah dapat
// menyusun kunci penyimpanannya.
//
// Dibaca sesudah: dokumenbaru.go.

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// AwalanImageID adalah literal pembuka masukan hash.
//
// ⛔ VERBATIM `GenerateImageID_SQL.xml` b86: `'ASMPP'` - LIMA huruf.
// Bukan `ASMAPP`, yang merupakan awalan TOKEN penyimpanan
// (`services.awalanToken`, butir an) dan hal yang berbeda. Keduanya berdiri
// berdampingan di modul ini dan berbeda satu huruf; konstanta bernama
// dengan kutipan barisnya adalah satu-satunya cara itu tidak tertukar.
const AwalanImageID = "ASMPP"

// PanjangImageID adalah panjang keluarannya - MD5 sebagai heksa.
//
// `STANDARD_HASH(..., 'MD5')` mengembalikan RAW(16); Oracle menampilkannya
// sebagai 32 karakter heksa HURUF BESAR.
const PanjangImageID = 32

// bentukStempelImageID meniru `TO_CHAR(SYSTIMESTAMP,'YYYYMMDDHH24MISSFF9')`.
//
// ⚠️ `FF9` = sembilan angka pecahan detik. Go menulis nanodetik dengan
// `.000000000`, yang membawa titik - titik itu DIBUANG, sebab Oracle tidak
// menuliskannya untuk `FF9`.
const bentukStempelImageID = "20060102150405.000000000"

// StempelImageID menyusun bagian waktu masukan hash.
func StempelImageID(saat time.Time) string {
	return strings.Replace(saat.Format(bentukStempelImageID), ".", "", 1)
}

// MasukanImageID merakit teks yang di-hash.
//
// ⛔ Dipisah dari perhitungannya supaya dapat dibandingkan dengan Oracle
// apa adanya: uji `db` menghitung `STANDARD_HASH` atas teks INI, sehingga
// yang dibandingkan rumusnya - bukan kebetulan dua jam yang sama.
//
// `guid` adalah padanan `SYS_GUID()`: 16 byte, ditulis 32 heksa huruf besar.
// Oracle mengubah RAW menjadi heksa secara implisit saat disambung ke teks.
func MasukanImageID(saat time.Time, guid string) string {
	return AwalanImageID + StempelImageID(saat) + strings.ToUpper(guid)
}

// ImageIDDari menghitung `IMAGEID` dari masukan yang sudah dirakit.
//
// ⛔ HURUF BESAR, sebab itulah yang Oracle tuliskan. `IMAGEID` menjadi kunci
// baris; heksa huruf kecil adalah kunci yang BERBEDA, dan pencarian atasnya
// akan mengembalikan nol baris tanpa satu pun galat.
func ImageIDDari(masukan string) string {
	jumlah := md5.Sum([]byte(masukan))
	return strings.ToUpper(hex.EncodeToString(jumlah[:]))
}

// ImageIDBaru menerbitkan satu `IMAGEID` baru.
//
// ⚠️ `SYS_GUID()` diganti 16 byte acak kriptografis. Oracle tidak menjanjikan
// GUID-nya acak - ia hanya menjanjikan unik - jadi yang ditiru sifat yang
// DIPAKAI, yaitu ketidakterdugaan kunci penyimpanan.
func ImageIDBaru(saat time.Time) (string, error) {
	guid := make([]byte, 16)
	if _, err := rand.Read(guid); err != nil {
		// ⛔ Gagal TERANG. Pengenal penyimpanan yang jatuh kembali ke cap
		// waktu ketika acaknya gagal adalah pengenal yang dapat ditebak -
		// dan tidak ada satu pun yang akan tahu bahwa ia jatuh kembali.
		return "", fmt.Errorf("models: gagal menerbitkan ImageID: %w", err)
	}
	return ImageIDDari(MasukanImageID(saat, hex.EncodeToString(guid))), nil
}
