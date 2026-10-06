package services

// Panel Attachment modul Adjustment — menyusun baris Category + Count.
//
// ⛔ ATURANNYA SAMA PERSIS dengan modul Treaty In, dan kesamaan itu
// DINYATAKAN alih-alih dibagi lewat impor: `TestModulTidakMengimporModulLain`
// menolak impor silang, dan alasannya berdiri sendiri — dua modul yang
// berbagi kode services berbagi juga jadwal rilisnya.
//
// ⚠️ Yang TIDAK boleh menyimpang: perlakuan terhadap empat kode tanpa nama.
// Modul ini menampilkan `00003` `00004` `00008` `00009` dengan kodenya,
// tanpa nama, bertanda belum dipastikan — sama persis seperti Treaty In.
// Menebak pasangannya di SALAH SATU modul saja sudah cukup untuk menaruh
// berkas di kategori yang salah.

import (
	"sort"
	"strings"

	"nusantarare/modul/treatyinadjustment/backend/models"
)

// SusunKategoriLampiran menggabungkan katalog yang TERBACA dengan empat
// kode yang namanya belum dipastikan.
//
// ⛔ Keempat kode tanpa nama TIDAK dipasangkan dengan keempat nama tanpa
// kode. Keduanya sama-sama empat, berurutan, dan menggoda — dan
// memasangkannya menurut abjad adalah tebakan yang akan terlihat benar
// sampai seseorang mengunduh berkas dari kategori yang salah bertahun
// kemudian. Disapu 4 Oktober 2026: keempat kode dan keempat namanya NIHIL
// di korpus kedua modul.
func SusunKategoriLampiran(katalog map[string]string, lampiran []models.BarisLampiranWarisan) []models.BarisKategoriLampiran {
	cacah := map[string]int{}
	for _, l := range lampiran {
		cacah[l.KodeKategori]++
	}

	baris := map[string]models.BarisKategoriLampiran{}
	for kode, nama := range katalog {
		baris[kode] = models.BarisKategoriLampiran{
			Kode: kode, Nama: nama, Cacah: cacah[kode], Dipastikan: true,
		}
	}
	// ⛔ Keempat kode tanpa nama tetap TAMPIL. Menyembunyikannya membuat
	// kategori yang ada di sistem lama lenyap dari layar tanpa suara, dan
	// berkas yang tersimpan di sana menjadi tidak terjangkau.
	for _, kode := range models.KodeKategoriBelumDipastikan {
		if _, sudah := baris[kode]; sudah {
			continue
		}
		baris[kode] = models.BarisKategoriLampiran{
			Kode: kode, Cacah: cacah[kode], Dipastikan: false,
		}
	}
	for kode := range cacah {
		if _, sudah := baris[kode]; !sudah {
			baris[kode] = models.BarisKategoriLampiran{Kode: kode, Cacah: cacah[kode], Dipastikan: false}
		}
	}

	out := make([]models.BarisKategoriLampiran, 0, len(baris))
	for _, b := range baris {
		out = append(out, b)
	}
	// Urut KODE, bukan nama: kode itulah yang tetap, dan empat di antaranya
	// belum punya nama untuk diurutkan.
	sort.Slice(out, func(i, j int) bool { return out[i].Kode < out[j].Kode })
	return out
}

// NamaBerkasAman memeriksa aturan spanduk biru layar lama:
// *"Recommended safe substitute should be . or _"*.
//
// ⛔ Lingkup berlakunya ditetapkan KEPUTUSAN §18: **pintu masuk saja**.
// Ke-43 nama yang sudah tersimpan dibiarkan apa adanya — 37 di antaranya
// melanggar, dan menggantinya berarti menulis ke tabel warisan serta dapat
// memutus tautan `FILENAME` ke `T_STORAGE_IMAGE`.
//
// ⚠️ Garis miring ditolak lebih dulu dan terpisah: `../` di nama berkas
// bukan sekadar nama yang jelek, ia upaya menulis di luar folder unggahan.
func NamaBerkasAman(nama string) bool {
	n := strings.TrimSpace(nama)
	if n == "" || n == "." || n == ".." {
		return false
	}
	if strings.ContainsAny(n, `/\`) {
		return false
	}
	for _, r := range n {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		case r == '.' || r == '_' || r == '-':
		default:
			return false
		}
	}
	return true
}

// LampiranKontrak adalah jawaban panel Attachment untuk satu kontrak.
type LampiranKontrak struct {
	Kategori []models.BarisKategoriLampiran `json:"kategori"`
	Berkas   []models.BarisLampiranWarisan  `json:"berkas"`
}
