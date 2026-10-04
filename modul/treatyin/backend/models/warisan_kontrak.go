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
	Bordereaux        string `json:"bordereaux"`
	BordereauxCatatan string `json:"bordereauxCatatan"`
	CaraPembukuan     string `json:"caraPembukuan"`

	// `ContractRefNo` — ADA di 742 dari 1.854 baris.
	NomorRujukan string `json:"nomorRujukan"`
	// `TreatyLeader` — ADA di 659 baris; nilainya teks `"true"`/`"false"`.
	PemimpinTreaty string `json:"pemimpinTreaty"`

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

	// Tab Co-Ins Scale - tabel pendaratan kesembilan, migrasi 432.
	SkalaKoasuransi []BarisSkalaKoasuransiWarisan `json:"skalaKoasuransi"`

	// TeksMentah - kelima kunci tab teks APA ADANYA dari dokumen, berkunci
	// nama ejaannya. Diisi repository; yang MEMILIH di antaranya services,
	// sebab pemilihannya bergantung cabang. Tidak dikirim ke layar.
	TeksMentah map[string]string `json:"-"`

	// Dua tab teks, sudah dipilih menurut cabang.
	Pengecualian TabTeksWarisan `json:"pengecualian"`
	SyaratKhusus TabTeksWarisan `json:"syaratKhusus"`
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
}

// BarisPortofolioWarisan - satu baris tab Portfolio.
// Dari `Portfolio`; berisi di 152 dari 300.
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

// BarisEgnpiWarisan - satu baris tab EGNPI. Dari `M_TREATYIN_EGNPI`;
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

// BarisAngsuranWarisan - satu baris jadwal tab Installment, dari tabel ANAK
// (`M_TREATYIN_INSTALLMENTITEM`). Induknya memuat total; jadwalnya yang
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
