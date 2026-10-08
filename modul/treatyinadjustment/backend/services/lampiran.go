package services

// Panel Attachment modul Adjustment — menyusun baris Category + Count.
//
// ⛔ ATURANNYA SAMA dengan modul Treaty In, dan kesamaan itu DINYATAKAN
// alih-alih dibagi lewat impor: `TestModulTidakMengimporModulLain` menolak
// impor silang, dan alasannya berdiri sendiri — dua modul yang berbagi kode
// services berbagi juga jadwal rilisnya.
//
// ⭐ 8 Oktober 2026 — penamaan kategori diperbaiki (permintaan pemakai,
// tangkapan layar Pega): kesebelas NAMA dari `M_KATEGORIMASTERTREATY`,
// urut abjad nama seperti `GetMasterTreatyCategory_SQL` (`order by note`).
// Empat kode yang dulu tampil "belum dipastikan" kini bernama dari katalog.

import (
	"sort"
	"strings"

	"nusantarare/modul/treatyinadjustment/backend/models"
)

// namaKategoriSetara - nama Prop → nama Non-Prop untuk kode yang sama
// (`GetMasterTreatyCategory_Act` [2.2], `TreatyIn.ProportionType ==
// "NonProportional"`).
var namaKategoriSetara = map[string]string{
	"Pega Proportional Calculation /Perhitungan Pega Proportional": "Pega Non Proportional Calculation /Perhitungan Pega Non Proportional",
}

// SusunKategoriLampiran menyusun panel dari katalog dan lampiran kontrak.
//
//   - urut NAMA katalog (`order by note`); nama Non-Prop untuk kontrak
//     NonProportional menggantikan namanya di tempat yang sama;
//   - dicacah menurut KODE, atau menurut NAMA bila baris lama tidak
//     berkode;
//   - ⛔ kode yang ADA DI DATA tetapi tidak di katalog tetap TAMPIL —
//     menyembunyikannya membuat berkasnya tidak terjangkau. Tanpa nama, ia
//     tampil dengan kodenya dan `Dipastikan=false`.
func SusunKategoriLampiran(katalog map[string]string, lampiran []models.BarisLampiranWarisan, proporsional bool) []models.BarisKategoriLampiran {
	kodeDariNama := map[string]string{}
	for kode, nama := range katalog {
		if nama != "" {
			kodeDariNama[nama] = kode
		}
	}
	cacah := map[string]int{}
	for _, l := range lampiran {
		kode := strings.TrimSpace(l.KodeKategori)
		if kode == "" {
			kode = kodeDariNama[strings.TrimSpace(l.NamaKategori)]
		}
		cacah[kode]++
	}

	out := make([]models.BarisKategoriLampiran, 0, len(katalog)+len(cacah))
	for kode, nama := range katalog {
		out = append(out, models.BarisKategoriLampiran{
			Kode: kode, Nama: nama, Cacah: cacah[kode], Dipastikan: nama != "",
		})
	}
	for kode, n := range cacah {
		if _, ada := katalog[kode]; !ada && kode != "" {
			out = append(out, models.BarisKategoriLampiran{Kode: kode, Cacah: n})
		}
	}
	// Bernama lebih dulu, urut nama; yang tanpa nama di ujung, urut kode.
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if (a.Nama == "") != (b.Nama == "") {
			return a.Nama != ""
		}
		if a.Nama != b.Nama {
			return a.Nama < b.Nama
		}
		return a.Kode < b.Kode
	})
	if !proporsional {
		for i := range out {
			if np, ada := namaKategoriSetara[out[i].Nama]; ada {
				out[i].Nama = np
			}
		}
	}
	return out
}

// SifatProporsional - `TreatyIn.ProportionType`; selain `NonProportional`
// (termasuk kosong) memakai nama katalog apa adanya.
func SifatProporsional(jenis string) bool {
	return strings.TrimSpace(jenis) != "NonProportional"
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
