package services

// Pemilihan ejaan untuk dua tab teks — Exclusions dan Special Conditions.
//
// ⛔ Jalan B, keputusan pemilik proses 4 Oktober 2026
// (`docs/KEPUTUSAN-PENYELARASAN-REPO.md` §15): isinya tetap di `JSONDATA`,
// nol tabel pendaratan. Yang dikerjakan di sini hanya MEMILIH di antara
// ejaan yang ada — dan pemilihan itulah bagian yang paling mudah salah.
//
// ⛔ EJAANNYA BUKAN SINONIM. Terukur atas 1.854 dokumen: dari **303**
// dokumen yang punya lebih dari satu ejaan `SpecialConditions*`, **nol**
// yang isinya identik. Karena itu "jatuh ke ejaan satunya" DILARANG: ia
// menampilkan teks yang salah, bukan salinan yang basi, dan pembacanya
// tidak punya cara tahu.

import (
	"sort"

	"nusantarare/modul/treatyin/backend/models"
)

// Ejaan per cabang. `…P` proporsional, telanjang non-proporsional.
const (
	EjaanPengecualianProp    = "ExclusionsP"
	EjaanPengecualianNonProp = "Exclusions"
	EjaanSyaratProp          = "SpecialConditionsP"
	EjaanSyaratNonProp       = "SpecialConditions"

	// ⚠️ Ejaan KETIGA, berhuruf kecil di akhir. Bukan salah ketik: 292
	// dokumen memakainya — 170 proporsional, 122 non-proporsional. Karena
	// ia dipakai KEDUA cabang, ia tidak terikat cabang mana pun, dan
	// karenanya SELALU muncul sebagai "ejaan lain" — tidak pernah terpilih.
	EjaanSyaratKetiga = "SpecialConditionsp"
)

// SifatProporsional menjawab apakah kontrak ini cabang proporsional.
//
// ⛔ Dibaca dari `TREATY_IN.PROPORTIONTYPE`, KOLOM — bukan dari kunci
// `ProportionType` di dalam dokumen. Terukur: kolomnya terisi pada seluruh
// 1.854 baris, sedangkan kunci dokumennya TIDAK ADA pada tiga (`1001854`,
// `1001855`, `1001856`). Memilih ejaan dengan kunci yang kadang hilang
// berarti ketiga kontrak itu diam-diam diperlakukan sebagai non-proporsional.
func SifatProporsional(sifatAsli string) bool {
	return sifatAsli == "Proportional"
}

// PilihTabTeks memilih satu ejaan menurut cabang, dan menyebut ejaan lain
// yang juga berisi.
//
// Aturannya, apa adanya dari keputusan §15:
//
//   - ejaan yang SESUAI cabang dipakai;
//   - bila ia tidak ada, hasilnya KOSONG — bukan ejaan lain;
//   - ejaan lain yang berisi tetap DISEBUT, supaya pembacanya tahu ada teks
//     yang tidak ia lihat.
func PilihTabTeks(mentah map[string]string, proporsional bool, ejaanProp, ejaanNonProp string, ejaanLain ...string) models.TabTeksWarisan {
	dipilih := ejaanNonProp
	if proporsional {
		dipilih = ejaanProp
	}

	hasil := models.TabTeksWarisan{EjaanLain: []string{}}
	if isi, ada := mentah[dipilih]; ada {
		hasil.Isi = isi
		hasil.Ejaan = dipilih
	}

	// ⛔ Seluruh ejaan lain diperiksa, termasuk yang MILIK cabang seberang.
	// Itu pokoknya: 175 kontrak proporsional punya `Exclusions` telanjang
	// berisi DI SAMPING `ExclusionsP`-nya, dan isinya berbeda.
	semua := append([]string{ejaanProp, ejaanNonProp}, ejaanLain...)
	for _, e := range semua {
		if e == dipilih {
			continue
		}
		if isi, ada := mentah[e]; ada && isi != "" {
			hasil.EjaanLain = append(hasil.EjaanLain, e)
		}
	}
	sort.Strings(hasil.EjaanLain)
	return hasil
}

// isiTabTeks mengisi kedua tab teks pada satu kontrak yang sudah dibaca.
func isiTabTeks(k *models.KontrakWarisan) {
	prop := SifatProporsional(k.SifatProporsiAsli)
	k.Pengecualian = PilihTabTeks(k.TeksMentah, prop,
		EjaanPengecualianProp, EjaanPengecualianNonProp)
	k.SyaratKhusus = PilihTabTeks(k.TeksMentah, prop,
		EjaanSyaratProp, EjaanSyaratNonProp, EjaanSyaratKetiga)
}
