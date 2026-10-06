package models

// Untuk apa berkas ini: MUATAN KONVERSI ke Arasapas - Utility2
// `serviceInsertArasapas_act` sesudah realisasi selesai.
//
// `[keputusan work owner]` KEPUTUSAN-RONDE-12 butir 7 (23-09-2026): "Ikuti
// yang terbaca, cabut P8." Bentuk ditiru dari saudara kelas induk
// `ASM-FW-GISFW-WORK` (versi 19) - salinan kelas `Data-PolicyTreatyIn` di korpus
// hanya stub. `[penyimpangan sadar]` bentuk dari kelas yang berbeda.
//
//	InputData.CARI1  = pyWorkPage.pzInsKey
//	InputData.CARI2  = PolicyNo
//	InputData.CARI3  = PolicyNo
//	InputData.CARI21 = ProdDateTime, "dd/MM/yyyy HH:mm:ss"
//
// Kegagalan TIDAK membatalkan penyimpanan polis; layar menerima penanda
// `FlagErrorKonversi` (verbatim di bawah).

import (
	"strings"
	"time"
)

// PesanGagalKonversi - VERBATIM `FlagErrorKonversi`.
const PesanGagalKonversi = "Gagal Konversi, Silahkan Coba Lagi atau Hub IT !"

// MuatanKonversi adalah empat medan yang dikirim.
type MuatanKonversi struct {
	CARI1, CARI2, CARI3, CARI21 string
}

// RakitMuatanKonversi menyusun muatan dari halaman kasus yang selesai.
func RakitMuatanKonversi(id string, h *Halaman) MuatanKonversi {
	nopol := h.Ambil(HalamanPolis + ".PolicyNo")
	prod := strings.TrimSpace(h.Ambil(HalamanPolis + ".ProductionDate"))
	if t, err := time.Parse("2006-01-02 15:04:05", prod); err == nil {
		prod = t.Format("02/01/2006 15:04:05")
	}
	return MuatanKonversi{CARI1: KunciInstans(id), CARI2: nopol, CARI3: nopol, CARI21: prod}
}
