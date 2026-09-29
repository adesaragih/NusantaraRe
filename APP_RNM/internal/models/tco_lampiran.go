package models

// Lampiran tahun treaty - tiket 12 Treaty Contract Out (FITUR BARU).
//
// Untuk apa berkas ini: bentuk satu baris `T_TREATYYEAR_LAMPIRAN` dan aturan
// murni di sekitarnya - status, kategori, nama berkas di disk dan di zip.
//
// ⚠️ Penyimpangan sadar 9 `[keputusan work owner]`: jalur lampiran Pega
// (`M_ATTACHMENTTREATY_2`) belum rampung dan berkunci ID treaty INWARD. Korpus
// di sini sumber RANTAI TEKNIS (kategori, kunci berkas, label), bukan sumber
// perilaku; perilakunya ditetapkan AC tiket 12.
//
// Dibaca sesudah: tco_tahun.go.

import (
	"path"
	"regexp"
	"strings"
	"time"
)

// LampiranTCO adalah satu baris `T_TREATYYEAR_LAMPIRAN`.
type LampiranTCO struct {
	ID           string
	IDTreatyYear string
	FileName     string
	FileMimeType string
	Category     string
	// ImageID adalah kunci berkas di penyimpanan - acak, lahir bersama
	// barisnya, tidak pernah berubah. Pengulangan menulis ke kunci yang sama.
	ImageID string
	// TStorageID terisi sesudah penyimpanan memastikan berkasnya ada.
	TStorageID string
	Ukuran     int64
	UserID     string
	TglUpload  time.Time
}

// Status lampiran yang tampil di layar.
//
// ⛔ TEKS tetap, dibaca layar dan uji; kosakata kami (tiket 12 "terkirim" /
// "tertunda"), bukan istilah Pega - fiturnya tidak ada di Pega.
const (
	StatusLampiranTerkirim = "terkirim"
	StatusLampiranTertunda = "tertunda"
	StatusLampiranGagal    = "gagal"
)

// StatusLampiranTCO menurunkan status dari kunci penyimpanan dan outbox.
//
// ⛔ Terkirim MENANG. Berkas yang sudah dipastikan ada di penyimpanan tidak
// menjadi "gagal" hanya karena efek pengulangan sesudahnya menyerah.
// `efekMenyerah` = efek unggah TERAKHIR lampiran ini berstatus gagal-permanen.
func StatusLampiranTCO(tStorageID string, efekMenyerah bool) string {
	switch {
	case strings.TrimSpace(tStorageID) != "":
		return StatusLampiranTerkirim
	case efekMenyerah:
		return StatusLampiranGagal
	default:
		return StatusLampiranTertunda
	}
}

// KategoriLampiranSah mencocokkan pilihan pemakai dengan master kategori.
//
// Mengembalikan teks master APA ADANYA - itulah yang disimpan di `CATEGORY`,
// supaya satu kategori tidak tersimpan dalam tiga ejaan.
func KategoriLampiranSah(master []string, pilihan string) (string, bool) {
	p := strings.TrimSpace(pilihan)
	if p == "" {
		return "", false
	}
	for _, m := range master {
		if strings.EqualFold(strings.TrimSpace(m), p) {
			return m, true
		}
	}
	return "", false
}

// akhiranSah membatasi akhiran berkas di disk: huruf kecil dan angka, 1-10.
var akhiranSah = regexp.MustCompile(`^\.[a-z0-9]{1,10}$`)

// NamaBerkasAntreLampiranTCO menyusun nama berkas antrean di `UNGGAHAN_DIR`.
//
// ⛔ Nama unggahan TIDAK PERNAH menentukan jalur: yang dipakai hanya kunci
// berkas (`IMAGEID`, heksa) dan akhiran yang lolos saringan. Nama
// `..\..\x.exe` tidak keluar dari folder antrean.
func NamaBerkasAntreLampiranTCO(imageID, namaAsli string) string {
	dasar := path.Base(strings.ReplaceAll(namaAsli, `\`, "/"))
	akhir := strings.ToLower(path.Ext(dasar))
	if !akhiranSah.MatchString(akhir) {
		return imageID
	}
	return imageID + akhir
}

// NamaEntriZipLampiranTCO menyusun nama entri zip "Download All".
//
// ⛔ Berawalan ID lampiran supaya dua berkas senama tidak saling timpa, dan
// hanya nama dasarnya - entri zip berjalur `../` adalah cara klasik menulis
// ke luar folder tujuan saat diekstrak.
func NamaEntriZipLampiranTCO(id, namaAsli string) string {
	dasar := path.Base(strings.ReplaceAll(strings.TrimSpace(namaAsli), `\`, "/"))
	if dasar == "." || dasar == "/" || dasar == ".." || dasar == "" {
		dasar = "berkas"
	}
	return id + "_" + dasar
}
