package services

import "nusantarare/modul/treatyin/backend/models"

// NILAI PROMPT rule Property — satu tempat, sebagai DATA.
//
// ---------------------------------------------------------------------
// ⛔ MASALAH YANG BERKAS INI PECAHKAN
// ---------------------------------------------------------------------
// Rule Property Pega menyimpan DUA kolom untuk tiap pilihan:
//
//	Standard value   yang TERSIMPAN di basis data   `1`, `risk`, `OR`
//	Prompt value     yang DIBACA pemakai            `Of Cession to R/I`, …
//
// Pemilik proses 7 Oktober 2026: *"seharusnya yg diambil itu nilai dari
// prompt value bukan standar value, krn di aplikasi ada yg salah"*.
//
// ⛔ AKAR SEBABNYA: `Rule-Obj-Property` NOL diekspor. Sapuan keempat korpus
// `D:\XML_NURE` menemukan nol berkas berjenis itu; yang ada hanya dua
// Property yang dikirim terpisah ke
// `_migration-docs/treaty-in-adjustment/ekspor-tambahan/`. Jadi label tampil
// TIDAK DAPAT dibaca dari ekspor, dan selama ini nilai tersimpannya yang
// ditampilkan apa adanya.
//
// ---------------------------------------------------------------------
// ⛔ MENGAPA LABELNYA TIDAK DITEBAK
// ---------------------------------------------------------------------
// `OptionLimit` membuktikan menebak akan meleset. Pasangan yang masuk akal
// bagi `Of Cession to R/I` adalah sesuatu tentang cession; yang sebenarnya
// `Of 100% Limit`. Pilihan berlabel karangan terbaca BENAR sampai seseorang
// memilihnya, lalu angkanya salah tanpa satu pun tanda.
//
// ⭐ Jadi yang belum diketahui JATUH KE NILAINYA SENDIRI — terlihat ganjil,
// dan ganjilnya disengaja. `TestPromptValueYangBelumDiekspor` mencetak
// daftarnya supaya ia menyusut terukur, bukan terlupakan.

// promptValue[nama properti][standard value] = prompt value.
//
// ⛔ Tambahkan HANYA dari ekspor rule Property atau tangkapan layarnya —
// bukan dari dugaan, bukan dari nama kolomnya.
var promptValue = map[string]map[string]string{
	// Dikirim pemilik proses 7 Oktober 2026 (tangkapan layar rule).
	// `ASM-FW-GISFW-Int-TREATY_IN` · UI `pxTextInput` · Table `Prompt List`.
	"OptionLimit": {
		"1": "Of Cession to R/I",
		"2": "Of 100% Limit",
	},
	// `_migration-docs/treaty-in-adjustment/ekspor-tambahan/EDMState.xml`.
	"EDMState": {
		"1": "Internal",
		"2": "External",
	},
	// `…/ekspor-tambahan/EDMMaterialType.xml`.
	"EDMMaterialType": {
		"1": "Material",
		"2": "Non Material",
	},
}

// ⭐ TINGKAT KEYAKINAN tiap label — dan mengapa ia dicatat.
//
// Sampai 7 Oktober 2026 berkas ini hanya memuat label yang DIEKSPOR, dan
// sisanya jatuh ke nilainya sendiri. Pemilik proses lalu memerintahkan
// seluruh dropdown diperbaiki, dengan satu contoh: *"riskcat itu seharusnya
// yang diambil risk & cat"*.
//
// ⛔ Satu contoh itu MENGUNCI `Cover`, bukan yang lain. Jadi sisanya
// disimpulkan — dan kesimpulan disimpan TERPISAH dari yang terbukti, supaya
// yang membacanya tahu mana yang dapat dipercaya dan mana yang menunggu
// ralat. Menggabungkannya ke satu peta akan menghapus bedanya selamanya.
const (
	AsalEkspor     = "ekspor"     // berkas rule atau tangkapan layarnya
	AsalPemilik    = "pemilik"    // dinyatakan pemilik proses
	AsalKesimpulan = "kesimpulan" // DISIMPULKAN — menunggu ralat
)

// promptDisimpulkan - label yang BELUM terbukti, beserta alasannya.
//
// ⚠️ Tiap baris menyebut DASAR kesimpulannya. Kesimpulan tanpa dasar tidak
// dapat diperiksa, dan yang tidak dapat diperiksa akan dianggap fakta.
var promptDisimpulkan = map[string]map[string]string{
	// ⭐ `riskcat` DINYATAKAN pemilik proses; dua lainnya dibiarkan apa
	// adanya sebab ia tidak menyebutnya — dan keduanya memang sudah kata
	// yang terbaca.
	// ⛔ HANYA `riskcat` yang terdaftar. `risk` dan `cat` labelnya SAMA
	// dengan nilainya, dan label yang sama dengan nilainya bukan kesimpulan
	// melainkan cadangan — mendaftarkannya membuat cermin di
	// `TestBerapaLabelMasihDisimpulkan` menghitung yang tidak perlu
	// diralat siapa pun.
	"Cover": {
		"riskcat": "risk & cat",
	},
	// Dasar: pola pemisahan yang SAMA dengan `riskcat` → `risk & cat`,
	// yaitu kode gandeng dipecah menjadi kata. `layer` sudah kata.
	"LayerType": {
		"sublayer": "sub layer",
	},
	// ⭐ `CurrencyRelation` SENGAJA TIDAK ADA di sini: `OR` dan `AND` sudah
	// kata utuh, dan Activity membacanya apa adanya
	// (`DetailCalculationROL` membandingkan `"AND"` / `"OR"`). Nol yang
	// perlu dipecah, jadi nol yang perlu diralat.
	// Dasar: pola pemisahan yang sama — `asamount` dan `astime` kode
	// gandeng dari frasa `as amount` dan `as time`.
	"ReinstatementNote": {
		"asamount": "as amount",
		"astime":   "as time",
	},
	// Dasar: pasangannya `loss` → `Loss Occuring` sudah ada di penerjemah,
	// dan lawan baku `Losses Occurring` di reasuransi adalah
	// `Risks Attaching`.
	// ⚠️ INI YANG PALING BERANI di berkas ini: bukan pemecahan kode,
	// melainkan istilah yang dibawa dari luar ekspor.
	"AccountingModeNonProp": {
		"risk": "Risk Attaching",
	},
	// Dasar: pasangannya `reporting` → `Reporting` sudah ada.
	"Bordeaux": {
		"nonreporting": "Non Reporting",
	},
}

// AsalLabel - dari mana label satu nilai berasal.
func AsalLabel(properti, nilai string) string {
	if m, ada := promptValue[properti]; ada {
		if _, ada := m[nilai]; ada {
			return AsalEkspor
		}
	}
	if properti == "Cover" {
		return AsalPemilik
	}
	switch properti {
	case "ReportingPeriod", "AccountingMode", "AccountingModeNonProp":
		if CaraPembukuanTampil(nilai) != nilai || PeriodePelaporanTampil(nilai) != nilai {
			return AsalEkspor
		}
	case "Bordeaux":
		if nilai == "reporting" {
			return AsalEkspor
		}
	}
	if m, ada := promptDisimpulkan[properti]; ada {
		if _, ada := m[nilai]; ada {
			return AsalKesimpulan
		}
	}
	return ""
}

// DomainProperty - nilai TERSIMPAN tiap dropdown, terukur di data, beserta
// nama rule Property-nya.
//
// ⭐ Namanya dipakai dua arah: mencari prompt value di atas, DAN menyusun
// daftar permintaan ekspor ke pemilik proses.
var DomainProperty = []struct {
	Properti string
	Nilai    []string
}{
	{"LayerType", domainJenisLayer},
	{"Cover", domainCover},
	{"CurrencyRelation", domainRelasiMataUang},
	{"ReinstatementNote", domainCatatanReinstatement},
	{"Bordeaux", domainBordereaux},
	{"AccountingMode", domainCaraPembukuan},
	{"AccountingModeNonProp", domainCaraPembukuanNonProp},
	{"ReportingPeriod", domainPeriodePelaporan},
	{"OptionLimit", []string{"1", "2"}},
	{"EDMState", []string{"1", "2"}},
	{"EDMMaterialType", []string{"1", "2"}},
}

// LabelPrompt - prompt value satu nilai tersimpan, atau nilainya sendiri
// bila rule Property-nya belum diekspor.
//
// ⚠️ Jatuh ke nilai sendiri BUKAN bawaan yang netral — ia tanda bahwa
// labelnya belum diketahui. Yang membacanya di layar melihat `riskcat`
// alih-alih kalimat, dan itu memang yang harus terlihat sampai ekspornya
// datang.
func LabelPrompt(properti, nilai string) string {
	if m, ada := promptValue[properti]; ada {
		if l, ada := m[nilai]; ada {
			return l
		}
	}
	// ⭐ TIGA PROPERTY SUDAH PUNYA PENERJEMAHNYA SENDIRI, dan ia TIDAK
	// disalin ke peta di atas: dua sumber untuk satu label adalah cara
	// termudah keduanya berbeda diam-diam. Ketiganya sudah diuji terpisah
	// (`TestLabelPeriodeDariPromptListPega` menyebut Prompt List Pega
	// sebagai sumbernya).
	//
	// ⚠️ Ketiganya juga JATUH ke nilainya sendiri untuk nilai yang
	// penerjemahnya tidak kenal — `risk` dan `nonreporting` termasuk, dan
	// keduanya ikut terdaftar sebagai belum diekspor.
	switch properti {
	case "ReportingPeriod":
		return PeriodePelaporanTampil(nilai)
	case "AccountingMode", "AccountingModeNonProp":
		// ⛔ Hanya kembali bila penerjemahnya MENGENALNYA. `risk` tidak ia
		// kenal dan dikembalikan apa adanya; kembali di sini akan melewati
		// peta kesimpulan di bawah.
		if l := CaraPembukuanTampil(nilai); l != nilai {
			return l
		}
	case "Bordeaux":
		if l := BordereauxTampil(nilai); l != nilai {
			return l
		}
	}
	// ⭐ Terakhir: label yang DISIMPULKAN. Ia sengaja dicari paling akhir —
	// apa pun yang terbukti selalu menang atasnya.
	if m, ada := promptDisimpulkan[properti]; ada {
		if l, ada := m[nilai]; ada {
			return l
		}
	}
	return nilai
}

// PromptBelumDiekspor - pasangan (properti, nilai) yang masih nol prompt
// value. Dipakai uji dan laporan permintaan ekspor.
func PromptBelumDiekspor() []string {
	var kurang []string
	for _, d := range DomainProperty {
		for _, v := range d.Nilai {
			if LabelPrompt(d.Properti, v) == v {
				kurang = append(kurang, d.Properti+"."+v)
			}
		}
	}
	return kurang
}

// opsiDariProperty - daftar dropdown: NILAI tersimpan berpasangan PROMPT
// VALUE, dengan nilainya sendiri sebagai cadangan.
//
// ⛔ Menggantikan `opsiApaAdanya`, yang menyalin nilai ke label tanpa
// perantara. Perbedaannya bukan kosmetik: lewat fungsi ini, label yang
// datang kemudian cukup ditambahkan ke `promptValue` — nol perubahan kode
// di tujuh titik panggil.
func opsiDariProperty(properti string, domain []string) []models.Opsi {
	out := make([]models.Opsi, 0, len(domain))
	for _, v := range domain {
		out = append(out, models.Opsi{Nilai: v, Label: LabelPrompt(properti, v)})
	}
	return out
}
