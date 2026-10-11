package services

// Panel Attachment — menyusun baris Category + Count.
//
// ⛔ Yang diputuskan di sini: kategori mana yang TAMPIL, dengan nama apa,
// dan mana yang ditandai belum dipastikan. Pembacaan dan penandaannya
// dipisah dengan sengaja — repository membaca fakta, services menyatakan
// apa yang belum diketahui.

import (
	"sort"
	"strings"

	"nusantarare/modul/treatyin/backend/models"
)

// namaKategoriSetara - nama Non-Prop → nama Prop untuk kode yang sama
// (`GetMasterTreatyCategory_Act` [2.2], `TreatyIn.ProportionType ==
// "NonProportional"`).
var namaKategoriSetara = map[string]string{
	"Pega Non Proportional Calculation /Perhitungan Pega Non Proportional": "Pega Proportional Calculation /Perhitungan Pega Proportional",
}

// SusunKategoriLampiran menggabungkan katalog yang TERBACA dengan empat
// kode yang namanya belum dipastikan.
//
// ⛔ ATURAN YANG PALING PENTING DI BERKAS INI: keempat kode tanpa nama
// TIDAK dipasangkan dengan keempat nama tanpa kode. Keduanya sama-sama
// empat, berurutan, dan menggoda — dan memasangkannya menurut abjad adalah
// tebakan yang akan terlihat benar sampai seseorang mengunduh berkas dari
// kategori yang salah bertahun kemudian.
//
// Yang ditampilkan: kodenya, dengan `Dipastikan=false`, dan keempat nama
// yang belum berumah disebut terpisah di layar.
func SusunKategoriLampiran(namaDesain []string, katalog map[string]string,
	lampiran []models.BarisLampiranWarisan) []models.BarisKategoriLampiran {
	// Kode yang diketahui per NAMA, dibalik dari katalog. Katalog menyimpan
	// `CATEGORY_ID` dan `CATEGORY` berdampingan, jadi pasangan yang terbukti
	// dibaca dari data — bukan dihafal.
	kodeDariNama := map[string]string{}
	for kode, nama := range katalog {
		if nama != "" {
			kodeDariNama[nama] = kode
		}
	}
	// ⭐ `GetMasterTreatyCategory_Act` [2.2]: kontrak Non-Prop menampilkan
	// kode `00007` dengan nama Non-Prop — kodenya SAMA, hanya namanya.
	for nonProp, prop := range namaKategoriSetara {
		if kode, ada := kodeDariNama[prop]; ada {
			if _, sudah := kodeDariNama[nonProp]; !sudah {
				kodeDariNama[nonProp] = kode
			}
		}
	}
	// ⛔ Dicacah menurut NAMA, bukan kode. Nama itulah yang kedua sisi
	// punya: panel menampilkan nama, dan baris lampiran membawa `CATEGORY`
	// berdampingan dengan kodenya. Mencacah menurut kode akan memberi NOL
	// pada keempat kategori yang kodenya belum dipastikan, padahal berkasnya
	// mungkin ada — dan nol yang salah terbaca persis seperti nol yang benar.
	cacahNama := map[string]int{}
	cacahKode := map[string]int{}
	for _, l := range lampiran {
		cacahNama[l.NamaKategori]++
		cacahKode[l.KodeKategori]++
	}

	out := make([]models.BarisKategoriLampiran, 0, len(namaDesain)+4)
	dipakai := map[string]bool{}
	for _, nama := range namaDesain {
		kode, ada := kodeDariNama[nama]
		cacah := cacahNama[nama]
		if cacah == 0 && ada {
			// Baris lama yang `CATEGORY`-nya kosong tetapi kodenya terisi.
			cacah = cacahKode[kode]
		}
		out = append(out, models.BarisKategoriLampiran{
			Kode: kode, Nama: nama, Cacah: cacah, Dipastikan: ada,
		})
		dipakai[nama] = true
	}

	// ⛔ Kategori yang ADA DI DATA tetapi tidak di daftar design tetap
	// TAMPIL. Menyembunyikannya membuat berkas yang tersimpan di sana tidak
	// terjangkau siapa pun — dan daftar design adalah gambar layar, bukan
	// jaminan tentang isi tabel warisan.
	var asing []string
	for nama := range cacahNama {
		if nama != "" && !dipakai[nama] {
			asing = append(asing, nama)
		}
	}
	sort.Strings(asing)
	for _, nama := range asing {
		out = append(out, models.BarisKategoriLampiran{
			Kode: kodeDariNama[nama], Nama: nama, Cacah: cacahNama[nama],
			Dipastikan: kodeDariNama[nama] != "",
		})
	}
	return out
}

// NamaBerkasAman memeriksa aturan spanduk biru layar lama:
// *"Recommended safe substitute should be . or _"*.
//
// ⛔ Spanduk itu ATURAN, bukan hiasan — dan ia menyebut penggantinya, bukan
// yang diganti. Yang ditegakkan di sini: nama berkas hanya boleh memuat
// huruf, angka, titik, garis bawah, dan tanda hubung. Apa pun di luar itu —
// spasi, tanda kurung, garis miring, tanda baca — ditolak dengan pesan yang
// MENYEBUT pengganti amannya, persis seperti spanduknya.
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
