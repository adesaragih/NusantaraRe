package models

// Lampiran kontrak — panel Attachment, dari `POOLDATA.M_ATTACHMENTTREATY_2`.
//
// ⛔ Tabel WARISAN yang hidup, 43 baris. Nol tabel baru, nol migrasi: aturan
// "tabel baru hanya untuk struktur tab ber-JSON" tidak berlaku di sini sebab
// lampirannya sudah relasional.

// BarisLampiranWarisan - satu berkas yang terlampir pada satu kontrak.
//
// Nama medan mengikuti alias di `GetAttachment2_Sql`, yang menamai ulang
// kolomnya: `filename as "pyFileName"`, `CATEGORY_ID as "pyCategory"`,
// `FILEMIMETYPE as "pyFileMimeType"`, `T_STORAGE_ID as "type"`.
type BarisLampiranWarisan struct {
	ID           string `json:"id"`
	KodeKategori string `json:"kodeKategori"`
	NamaKategori string `json:"namaKategori"`
	NamaBerkas   string `json:"namaBerkas"`
	JenisMime    string `json:"jenisMime"`
	IDSimpanan   string `json:"idSimpanan"`
	Diunggah     string `json:"diunggah"`
	Pengunggah   string `json:"pengunggah"`
}

// BarisKategoriLampiran - satu baris panel Attachment: Category + Count.
//
// ⚠️ Panel lama menampilkan SELURUH kategori, termasuk yang nol berkas —
// kolom `Count` tidak akan pernah berbunyi `0` kalau barisnya disembunyikan
// saat kosong.
type BarisKategoriLampiran struct {
	Kode  string `json:"kode"`
	Nama  string `json:"nama"`
	Cacah int    `json:"cacah"`
	// ⛔ Dipastikan menyatakan apakah pasangan KODE↔NAMA ini terbukti.
	//
	// Tujuh terbukti dari data itu sendiri: `M_ATTACHMENTTREATY_2` menyimpan
	// `CATEGORY_ID` DAN `CATEGORY` berdampingan, jadi pasangannya dibaca,
	// bukan dihafal. Empat sisanya — `00003` `00004` `00008` `00009` — nol
	// baris, dan namanya TIDAK ADA di korpus kedua modul (disapu 4 Oktober
	// 2026). Layar menandainya, sebab berkas yang mendarat di kategori yang
	// salah baru ketahuan bertahun kemudian.
	Dipastikan bool `json:"dipastikan"`
}

// NamaKategoriBelumDipastikan adalah keempat nama yang layar lama
// tampilkan tetapi kodenya TIDAK diketahui.
//
// ⛔ URUTAN DI SINI BUKAN PASANGAN. Keempatnya ditulis menurut abjad supaya
// tidak ada yang membacanya sebagai urutan kode — menebak `00003` =
// "Binding" hanya karena keduanya pertama adalah persis kesalahan yang
// penandaan ini ada untuk mencegahnya.
//
// Disapu dan NIHIL di `D:\XML_NURE\Treaty In` dan `Treaty In Adjustment`:
// keempat kodenya nol kemunculan (satu-satunya `00009` ternyata potongan
// cap waktu `20210609T100009_887_GMT`), dan keempat namanya nol berkas.
// Pertanyaannya di `docs/PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md`.
var NamaKategoriBelumDipastikan = []string{
	"Binding, signed share Email",
	"Claim Data",
	"Info Pack",
	"Letter of Acknowledgment / LOA",
}

// KodeKategoriBelumDipastikan adalah keempat kode yang namanya tidak
// diketahui. Sejajar dengan daftar di atas HANYA pada cacahnya, BUKAN pada
// urutannya.
var KodeKategoriBelumDipastikan = []string{"00003", "00004", "00008", "00009"}

// ===========================================================================
// KESEBELAS NAMA KATEGORI LAMPIRAN — yang panel tampilkan, apa adanya
// ===========================================================================
//
// ⛔ Daftar INILAH yang merender panel, bukan katalog `M_ATTACHMENTTREATY_2`.
// Sebabnya: panel lama menampilkan kesebelas kategori beserta `Count 0`-nya,
// sementara katalog hanya memuat kategori yang PERNAH dipakai — tujuh. Panel
// yang dirender dari katalog kehilangan empat baris pada kontrak mana pun
// yang belum memakainya, dan kehilangan itu diam.
//
// Dipastikan pemilik proses dengan tangkapan layar, 6 Oktober 2026.
//
// ⚠️ Urutannya ALFABETIS, persis seperti layar lama — dan itu BUKAN urutan
// kode. Menebak `00003` = "Binding" hanya karena keduanya pertama adalah
// persis kesalahan yang `Dipastikan` ada untuk mencegahnya.
//
// ⭐ Daftarnya BERCABANG pada butir kesepuluh saja.
var NamaKategoriLampiranProp = []string{
	"Analysed Email",
	"Approval Email",
	"Assessment Inward Treaty Form / Format Analisa Treaty",
	"Binding, signed share Email",
	"Claim Data",
	"Info Pack",
	"Letter of Acknowledgment / LOA",
	"Offer Email",
	"Others",
	"Pega Proportional Calculation /Perhitungan Pega Proportional",
	"Summary Treaty Leader",
}

// NamaKategoriLampiranNonProp - sama, kecuali butir kesepuluh.
var NamaKategoriLampiranNonProp = []string{
	"Analysed Email",
	"Approval Email",
	"Assessment Inward Treaty Form / Format Analisa Treaty",
	"Binding, signed share Email",
	"Claim Data",
	"Info Pack",
	"Letter of Acknowledgment / LOA",
	"Offer Email",
	"Others",
	"Pega Non Proportional Calculation /Perhitungan Pega Non Proportional",
	"Summary Treaty Leader",
}
