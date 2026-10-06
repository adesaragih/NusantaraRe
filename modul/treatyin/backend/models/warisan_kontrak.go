package models

// Satu kontrak WARISAN selengkapnya — `TREATY_IN` + `M_TREATY_IN.JSONDATA`.
//
// ⛔ DUA SUMBER, dan pembagiannya diukur - bukan ditebak. Sapuan 3 Oktober
// 2026 atas seluruh 1.854 baris:
//
//	dari kolom TREATY_IN   TREATYCONTRACTNAME · TERITORIALSCOPE · COMMENCEMENT
//	                       TERMINATION · TREATYYEAR · CEDING(+ID)
//	                       LEADINGREINSSOURCE(+ID) · PROPORTIONTYPE
//	dari JSONDATA          Bordeaux 1.851/1.854 · AccountingMode 1.851
//	                       BordereauxNote 1.018 · ContractRefNo 742
//	                       TreatyLeader 659
//
// ⚠️ `TeritorialScope` ADA DI DUA TEMPAT, dan keduanya hampir sama: 1.844
// dari 1.854 identik begitu `\n` dan `\r` di JSON di-unescape; 10 sisanya
// benar-benar berbeda. Yang dipakai KOLOMNYA - ia sudah berbentuk teks yang
// dapat ditampilkan, sementara nilai JSON masih membawa escape-nya.
type KontrakWarisan struct {
	ID string `json:"id"`

	// --- dari kolom `TREATY_IN`
	NamaKontrak    string `json:"namaKontrak"`
	LingkupWilayah string `json:"lingkupWilayah"`
	TahunTreaty    string `json:"tahunTreaty"`
	Cedant         string `json:"cedant"`
	IDCedant       string `json:"idCedant"`
	AsalBisnis     string `json:"asalBisnis"`
	IDAsalBisnis   string `json:"idAsalBisnis"`

	SifatProporsiAsli string `json:"sifatProporsiAsli"`
	SifatProporsi     string `json:"sifatProporsi"`

	TanggalMulaiAsli    string `json:"tanggalMulaiAsli"`
	TanggalBerakhirAsli string `json:"tanggalBerakhirAsli"`
	TanggalMulai        string `json:"tanggalMulai"`
	TanggalBerakhir     string `json:"tanggalBerakhir"`

	// --- dari `M_TREATY_IN.JSONDATA`
	//
	// ⚠️ Nilainya APA ADANYA, huruf kecil seperti tersimpan. Sapuan
	// menemukan domainnya: `Bordeaux` adalah `reporting` (1.164) atau
	// `nonreporting` (687); `AccountingMode` adalah `underwriting` (1.110)
	// atau `accounting` (741). Layar lama menampilkan "Reporting" dan
	// "Accounting Year" - pemetaan dari yang tersimpan ke yang tampil TIDAK
	// ada di ekspor mana pun, jadi ia tidak dikarang di sini.
	// ⚠️ TIGA medan ini membawa LABEL; pasangan `…Asli`-nya membawa nilai
	// tersimpan. Layar memerlukan yang pertama, jalur tulis yang kedua.
	Bordereaux        string `json:"bordereaux"`
	BordereauxAsli    string `json:"bordereauxAsli"`
	BordereauxCatatan string `json:"bordereauxCatatan"`
	CaraPembukuan     string `json:"caraPembukuan"`
	CaraPembukuanAsli string `json:"caraPembukuanAsli"`
	// ⛔ Cara pembukuan cabang NON-PROPORSIONAL — properti yang berbeda,
	// bukan medan yang sama. `AccountingMode` bernilai
	// `underwriting`/`accounting`; `AccountingModeNonProp` bernilai
	// `loss`/`risk`. Layar non-prop membaca yang ini.
	CaraPembukuanNonProp     string `json:"caraPembukuanNonProp"`
	CaraPembukuanNonPropAsli string `json:"caraPembukuanNonPropAsli"`

	// Tiga nilai yang menentukan TAB MANA yang dirender. Kosong = kuncinya
	// tidak ada di dokumen ini, dan itu BUKAN sama dengan `false`.
	RetroBerganda    string `json:"retroBerganda"`
	EDMState         string `json:"edmState"`
	EDMJenisMaterial string `json:"edmJenisMaterial"`

	// `ContractRefNo` — ADA di 742 dari 1.854 baris.
	NomorRujukan string `json:"nomorRujukan"`
	// `TreatyLeader` — ADA di 659 baris; nilainya teks `"true"`/`"false"`.
	PemimpinTreaty string `json:"pemimpinTreaty"`

	// ⭐ TUJUH medan kepala tab Reporting Period — migrasi `444`,
	// 6 Oktober 2026.
	//
	// ⛔ Sampai hari itu ketujuhnya KOSONG di layar, dan sebabnya bukan
	// data yang tidak ada: terukur atas seluruh 1.855 dokumen
	// `M_TREATY_IN`, `ReportingPeriod` terisi di 1.851 dan kelima
	// saudaranya di 1.219 masing-masing — DUA PERTIGA korpus. Yang
	// kurang hanya kolom di `T_TREATY_REVISION`
	// (`PERTANYAAN-TERBUKA-LAYAR-PEGA.md` §17).
	//
	// ⚠️ Nilainya TERSIMPAN, bukan hasil `Apply`. Tombol `Apply`
	// (`services/periode_pelaporan.go`) MENGHITUNG daftar periode dari
	// ketujuh medan ini; medan ini sendiri adalah masukan yang pemakai
	// ketik dan Pega simpan.
	PeriodeMulai      string `json:"periodeMulai"`
	PeriodeAkhir      string `json:"periodeAkhir"`
	PeriodeJenis      string `json:"periodeJenis"`
	PeriodeInterval   string `json:"periodeInterval"`
	PeriodePenyerahan string `json:"periodePenyerahan"`
	PeriodeKonfirmasi string `json:"periodeKonfirmasi"`
	PeriodePelunasan  string `json:"periodePelunasan"`

	// AdaDiJSON menyatakan kunci MANA yang benar-benar ada di dokumen ini.
	//
	// ⛔ Perlu, dan sebabnya adalah seluruh pokok medan-medan ini: kunci
	// yang TIDAK ADA tidak dapat dibedakan dari kunci yang ada bernilai
	// kosong, dan keduanya berarti hal yang berbeda di layar. "Tidak ada di
	// sistem lama" bukan "belum diisi".
	AdaDiJSON map[string]bool `json:"adaDiJson"`

	// Empat larik isi - kosong berarti dokumennya memang tidak punya.
	Kurs             []BarisKursWarisan       `json:"kurs"`
	PeriodePelaporan []BarisPeriodeWarisan    `json:"periodePelaporan"`
	Portofolio       []BarisPortofolioWarisan `json:"portofolio"`
	Akumulasi        []BarisAkumulasiWarisan  `json:"akumulasi"`

	// Empat tab berikutnya, dari tabel pendaratan yang sama.
	Egnpi    []BarisEgnpiWarisan    `json:"egnpi"`
	Retensi  []BarisRetensiWarisan  `json:"retensi"`
	Angsuran []BarisAngsuranWarisan `json:"angsuran"`
	Catatan  []BarisCatatanWarisan  `json:"catatan"`

	// ⭐ Baris LAYER dari `M_TREATY_IN2` — melayani EMPAT tab sekaligus:
	// Limits, Share, Event Limits, dan RNM Share. Bukan empat medan, sebab
	// keempat tab itu memandang baris yang sama dari sisi yang berbeda.
	Layer []BarisLayerWarisan `json:"layer"`

	// ⭐ POHON tab Limits proporsional — `Limits[] → Detail[] → daftar`, tiap
	// simpul `map[string]any` berisi teks apa adanya atau larik simpul.
	// Bentuknya dari `Section/LimitProportional.xml` dan `DetailLimits.xml`.
	LimitsPohon []map[string]any `json:"limitsPohon"`

	// Pilihan dropdown kepala — pasangan NILAI TERSIMPAN ↔ LABEL TAMPIL,
	// diisi services. Dropdown memegang nilai tersimpan; labelnya untuk dibaca.
	OpsiKepala OpsiKepala `json:"opsiKepala"`

	// ⭐ Panel `Total Retention Amount` — DITURUNKAN dari `Retensi`, bukan
	// dibaca. Lihat `BarisTotalRetensiWarisan` untuk rumah rumusnya.
	TotalRetensi []BarisTotalRetensiWarisan `json:"totalRetensi"`

	// ⭐ Panel `Existing Policy for Master ID` — kanan atas, KEDUA cabang.
	PolisProduksi []BarisPolisProduksi `json:"polisProduksi"`

	// Tab Co-Ins Scale - tabel pendaratan kesembilan, migrasi 432.
	SkalaKoasuransi []BarisSkalaKoasuransiWarisan `json:"skalaKoasuransi"`

	// TeksMentah - kelima kunci tab teks APA ADANYA dari dokumen, berkunci
	// nama ejaannya. Diisi repository; yang MEMILIH di antaranya services,
	// sebab pemilihannya bergantung cabang. Tidak dikirim ke layar.
	TeksMentah map[string]string `json:"-"`

	// Dua tab teks, sudah dipilih menurut cabang.
	Pengecualian TabTeksWarisan `json:"pengecualian"`
	SyaratKhusus TabTeksWarisan `json:"syaratKhusus"`

	// Panel Attachment - dari `M_ATTACHMENTTREATY_2`, tabel WARISAN.
	Lampiran         []BarisLampiranWarisan  `json:"lampiran"`
	KategoriLampiran []BarisKategoriLampiran `json:"kategoriLampiran"`
}

// TabTeksWarisan - satu tab yang isinya SATU medan teks panjang.
//
// ⛔ Membawa TIGA hal, dan ketiganya diperlukan di layar:
//
//	Isi      teks yang dipilih menurut cabang; kosong bila ejaan cabangnya
//	         tidak ada - dan kosong itu BUKAN diisi dari ejaan lain.
//	Ejaan    ejaan mana yang dipakai, supaya yang memeriksa tahu dari mana
//	         teks itu datang tanpa membuka dokumen.
//	EjaanLain ejaan LAIN yang juga berisi. Pembacanya berhak tahu ada teks
//	         lain yang tidak ia lihat - terukur, isinya BERBEDA di seluruh
//	         303 dokumen yang punya lebih dari satu.
type TabTeksWarisan struct {
	Isi       string   `json:"isi"`
	Ejaan     string   `json:"ejaan"`
	EjaanLain []string `json:"ejaanLain"`
}

// BarisSkalaKoasuransiWarisan - satu baris tab Co-Ins Scale.
//
// ⚠️ `BagianKoasuransi` BUKAN angka melainkan PITA: `>=30% up to < 50%`.
// Kolom angka akan menolak seluruh 702 barisnya.
//
// ⭐ `Penyusun` dan `DisusunPada` dibawa dari medan jejak Pega
// (`pxCreateOpName`, `pxCreateDateTime`), yang ada pada 701 dari 702
// elemen. Rancangan ronde ini menyebut "dua medan saja"; keduanya dibawa
// sebab tab inilah yang angkanya dinegosiasikan, dan "siapa yang menyusun
// skala ini" adalah pertanyaan yang akan ditanyakan.
type BarisSkalaKoasuransiWarisan struct {
	BagianKoasuransi string `json:"bagianKoasuransi"`
	PersenLimit      string `json:"persenLimit"`
	Penyusun         string `json:"penyusun"`
	DisusunPada      string `json:"disusunPada"`
}

// ⭐ EMPAT LARIK ISI — grid dan tab, dibaca dari dokumen warisan.
//
// ⛔ Keempatnya BUKAN model baru. Mereka potret apa adanya dari
// `M_TREATY_IN.JSONDATA`, dan medannya sengaja bertipe string: nilai di
// dokumen seluruhnya string (terukur atas 150 dokumen), dan mengubahnya
// menjadi angka atau tanggal di sini berarti menafsirkan sebelum ada yang
// meminta tafsiran itu. Yang membaca selisih pemindahan membaca yang
// tersimpan.

// BarisKursWarisan - satu baris grid Rate of Exchange.
// Dari `CurrencyList`; berisi di 297 dari 300 dokumen yang disapu.
type BarisKursWarisan struct {
	MataUang      string `json:"mataUang"`
	NilaiKeIDR    string `json:"nilaiKeIDR"`
	BerlakuDari   string `json:"berlakuDari"`
	BerlakuSampai string `json:"berlakuSampai"`
}

// BarisPeriodeWarisan - satu baris grid tab Reporting Period.
// Dari `ReportingPeriodList`; berisi di 225 dari 300.
type BarisPeriodeWarisan struct {
	Periode          string `json:"periode"`
	HitungOtomatis   string `json:"hitungOtomatis"`
	TanggalAwal      string `json:"tanggalAwal"`
	JatuhTempoKirim  string `json:"jatuhTempoKirim"`
	JatuhTempoKonfir string `json:"jatuhTempoKonfirmasi"`
	JatuhTempoBayar  string `json:"jatuhTempoBayar"`
	// Bentuk TERSIMPAN (`YYYYMMDD`, tanggal WIB) — untuk KOTAK tanggal.
	// Medan di atas bentuk tampil, untuk DIBACA. Terjemahan tampil yang
	// dimasukkan ke kotak tanggal membuat kotaknya kosong.
	TanggalAwalAsli      string `json:"tanggalAwalAsli"`
	JatuhTempoKirimAsli  string `json:"jatuhTempoKirimAsli"`
	JatuhTempoKonfirAsli string `json:"jatuhTempoKonfirmasiAsli"`
	JatuhTempoBayarAsli  string `json:"jatuhTempoBayarAsli"`
}

// BarisPortofolioWarisan - satu baris tab Portfolio.
// Dari `Portfolio`. ⭐ RALAT 6 Oktober 2026: "152 dari 300" adalah sisa
// sapuan CONTOH yang lama. Terukur atas SELURUH korpus - 1.925 elemen di
// 845 dari 1.855 kontrak, nol gagal urai.
//
// ⚠️ Nama ruasnya MENYESATKAN dan itu dinyatakan, bukan diam-diam
// dibiarkan: `Jenis` membawa kunci `Type` (bernilai `Premium`/`Loss`) dan
// `JenisPortfolio` membawa kunci `TypePortfolio` (bernilai
// `Withdrawal`/`Assumption`) - jadi yang bernama "jenis" justru arahnya,
// dan sebaliknya. Yang BENAR adalah padanan kuncinya; penggantian namanya
// menyentuh `api.ts` dan uji milik orang lain, jadi ia dicatat di sini
// dan tidak dikerjakan menyelinap di ronde ini.
type BarisPortofolioWarisan struct {
	Jenis          string `json:"jenis"`
	JenisPortfolio string `json:"jenisPortfolio"`
	Keterangan     string `json:"keterangan"`
}

// BarisAkumulasiWarisan - satu baris tab Accumulation.
// Dari `AccumulationList`; berisi di 12 dari 300 - jarang, dan itu BUKAN
// alasan melewatkannya: dua belas kontrak yang punya akumulasi adalah dua
// belas kontrak yang datanya hilang bila tabnya tidak dibaca.
type BarisAkumulasiWarisan struct {
	Periode         string `json:"periode"`
	TanggalLapor    string `json:"tanggalLapor"`
	HariKirim       string `json:"hariKirim"`
	JatuhTempoKirim string `json:"jatuhTempoKirim"`
}

// BarisEgnpiWarisan - satu baris tab EGNPI. Dari `T_TREATY_EGNPI`;
// larik `EGNPI` berisi di 846 dari 1.854 kontrak (terukur, bukan 177 yang
// sempat dilaporkan - sapuan substring lama terkecoh kunci bersarang).
type BarisEgnpiWarisan struct {
	Jumlah         string `json:"jumlah"`
	JumlahIDR      string `json:"jumlahIDR"`
	PerTanggal     string `json:"perTanggal"`
	KelasBisnis    string `json:"kelasBisnis"`
	MataUang       string `json:"mataUang"`
	Keterangan     string `json:"keterangan"`
	Proporsi       string `json:"proporsi"`
	KelompokTreaty string `json:"kelompokTreaty"`
}

// BarisRetensiWarisan - satu baris tab Maximum Retention.
type BarisRetensiWarisan struct {
	Jumlah         string `json:"jumlah"`
	KelasBisnis    string `json:"kelasBisnis"`
	MataUang       string `json:"mataUang"`
	Keterangan     string `json:"keterangan"`
	KelompokTreaty string `json:"kelompokTreaty"`
}

// BarisPolisProduksi - satu baris panel `Existing Policy for Master ID`.
//
// ⛔ Dari `TREATYINPRODUCTION`, tabel WARISAN yang BELUM pernah dibaca modul
// ini sampai 5 Oktober 2026 — 41.936 baris. Keempat medannya berpadanan
// satu-ke-satu dengan judul kolom di gambar 01 dan 26.
type BarisPolisProduksi struct {
	// `NOPOLIS` -> kolom `Policy No`.
	NomorPolis string `json:"nomorPolis"`
	// `IDPEGA` dipotong 18 aksara -> kolom `Pega ID`.
	PegaID string `json:"pegaID"`
	// `QUARTER` -> kolom `Quarter`. Kosong pada contoh desain, dan itu sah.
	Kuartal string `json:"kuartal"`
	// `QUARTER_YEAR` -> kolom `Quarter Year`.
	TahunKuartal string `json:"tahunKuartal"`
}

// BarisTotalRetensiWarisan - satu baris panel `Total Retention Amount`.
//
// ⭐ TURUNAN, bukan simpanan. Nol kolom, nol kunci JSON, nol tabel
// pendaratan memuatnya; ia dihitung dari `Retensi` oleh services.
//
// ⛔ RUMUSNYA DIBACA, bukan ditebak dari nama panelnya —
// `Activity/TreatyInNPSetTotal.xml` (`pyRuleAvailable = Yes`, nol penjaga
// `1=2`), cabang `param.type=retention`:
//
//	local.Currency  := .Currency
//	local.Value     := .Amount
//	bila mata uangnya SUDAH ADA di daftar -> .Value := local.Value + .Value
//	bila belum                             -> APPEND {Currency, Value}
//
// Jadi: JUMLAH `Amount` tab Maximum Retention, DIKELOMPOKKAN per mata uang.
// Panelnya sendiri membaca `TreatyIn.TotalRetentionAmountNP`
// (`pyPageListProperty` @147731 di `Section/TreatyInTabsNonProportional.xml`)
// dengan dua kolom: `.Currency` @161105 dan `.Value` @166072.
type BarisTotalRetensiWarisan struct {
	MataUang string `json:"mataUang"`
	Nilai    string `json:"nilai"`
}

// BarisAngsuranWarisan - satu baris jadwal tab Installment, dari tabel ANAK
// (`T_TREATY_INSTALLMENT_ITEM`). Induknya memuat total; jadwalnya yang
// dibaca orang.
type BarisAngsuranWarisan struct {
	Angsuran     string `json:"angsuran"`
	MataUang     string `json:"mataUang"`
	Jumlah       string `json:"jumlah"`
	Persen       string `json:"persen"`
	JatuhTempo   string `json:"jatuhTempo"`
	TanggalBayar string `json:"tanggalBayar"`
	WPC          string `json:"wpc"`
}

// BarisCatatanWarisan - satu baris tab Information & Submit.
type BarisCatatanWarisan struct {
	Tanggal   string `json:"tanggal"`
	Operator  string `json:"operator"`
	Disetujui string `json:"disetujui"`
	Catatan   string `json:"catatan"`
}

// Opsi - satu pilihan dropdown: nilai TERSIMPAN dan label TAMPILNYA.
type Opsi struct {
	Nilai string `json:"value"`
	Label string `json:"label"`
}

// OpsiKepala - pilihan ketiga dropdown kepala form.
//
// ⛔ Daftar pilihannya TIDAK ada di ekspor maupun di database: dropdown
// `AccountingMode`/`AccountingModeNonProp`/`Bordeaux` ber-`pyListSource =
// associated` (rule Property yang tidak diekspor), dan tabel enumerasi Pega
// `DATAPEGA.PR_ASM_FW_GISFW_DATA_ENUMERATI` (1.217 baris) tidak memuatnya —
// diukur 6 Oktober 2026. Nilainya domain TERSIMPAN di 1.854 dokumen; labelnya
// dari tangkapan layar (services `CaraPembukuanTampil`, `BordereauxTampil`).
type OpsiKepala struct {
	Bordereaux           []Opsi `json:"bordereaux"`
	CaraPembukuan        []Opsi `json:"caraPembukuan"`
	CaraPembukuanNonProp []Opsi `json:"caraPembukuanNonProp"`
	PeriodePelaporan     []Opsi `json:"periodePelaporan"`
}
