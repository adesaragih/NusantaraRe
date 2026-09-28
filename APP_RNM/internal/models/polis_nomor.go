package models

// Bentuk `PL_NUMBER` - tiket 03 PremiumList Life.
//
// Untuk apa berkas ini: merakit nomor premium list dari potongan-potongannya.
// Seluruhnya MURNI - tidak satu pun bagian di sini menyentuh Oracle, supaya
// bentuk nomornya dapat diuji tanpa basis data.
//
// Sumbernya `Activity/SubmitPremiumList_Act.xml`, dibaca sebagai pohon
// 28-09-2026.
//
// `[terverifikasi]` 17 langkah tingkat pertama. Jendela: elemen
// `<pyStepPageReference>` berbentuk `RH_1.pySteps(N)` - tingkat SATU saja;
// sub-langkah `RH_1.pySteps(N).pySteps(M)` tidak ikut. Dua cara:
//
//	cara A - cacah rujukan tingkat-1:
//	  grep -c "<pyStepPageReference>RH_1\.pySteps([0-9]*)</pyStepPageReference>" \
//	    Activity/SubmitPremiumList_Act.xml
//	cara B - indeks TERTINGGI, yang menangkap langkah hilang di tengah:
//	  grep -o "<pyStepPageReference>RH_1\.pySteps([0-9]*)</pyStepPageReference>" \
//	    Activity/SubmitPremiumList_Act.xml | grep -o "[0-9]*" | sort -n | tail -1
//
// Keduanya menjawab 17, jadi nol langkah hilang di tengah. Enam di antaranya
// menyusun nomor:
//
//	langkah 6  b1783  `.Type=="QR"` -> `InputData.CARI20 = "QR"`
//	langkah 7  b1959  `.Type=="QP"` -> `InputData.CARI20 = "QP"`
//	langkah 8  b2127  `.Type=="TP"` -> `InputData.CARI20 = "TP"`
//	langkah 9  b2295  `.Type=="TR"` -> `InputData.CARI20 = "TR"`
//	langkah 10 b2463  `GetKodeProdLife_SQL`  -> `ParamSeq.HASIL3` (awalan)
//	langkah 11 b2669  `ParamSeq.CARI1 = pyWorkPage.pxObjClass`
//	           b2717  `ParamSeq.CARI2 = ParamSeq.HASIL3+"QR/QP/TP/TR"`
//	langkah 12 b2810  `GetSequenceNumber_SQL` -> `HASIL1` (MM.YYYY), `HASIL2`
//	langkah 13 b3021  `local.HASIL1 = @substring(ParamSeq.HASIL1,0,2)`
//	           b3075  `local.HASIL2 = @substring(ParamSeq.HASIL1,5,7)`
//	           b3104  `ParamSeq.HASIL1 = local.HASIL1+"."+local.HASIL2`
//	           b3126  `.PremiumListSummary.PL_NUMBER =
//	                    ParamSeq.HASIL3 + InputData.CARI20 +
//	                    pyWorkPage.BusinessCode + "." + ParamSeq.HASIL1 +
//	                    "." + ParamSeq.HASIL2`
//
// ⛔ SATU PENGHITUNG UNTUK EMPAT TIPE, dan inilah yang paling mudah dirusak
// "perbaikan". Baris b2717 memakai teks HARFIAH `"QR/QP/TP/TR"` sebagai
// bagian kunci penghitung - BUKAN tipe yang sedang berjalan. Tipe hanya masuk
// ke NOMORNYA (b3126), tidak ke kuncinya. Siapa pun yang membacanya sekilas
// akan mengira empat tipe berarti empat penghitung, memecahnya, lalu
// menerbitkan empat deret yang masing-masing mulai dari 1 - dan setiap nomor
// baru bertabrakan dengan nomor lama.
//
// ⛔ TAHUNNYA DUA ANGKA DI SINI, EMPAT ANGKA DI NOMOR KLAIM. Procedure yang
// sama mengembalikan `MM.YYYY`; nomor klaim memakainya apa adanya
// (`RakitNomorKlaim`), sedangkan nomor PL memotongnya menjadi `MM.YY` lewat
// kedua `@substring` di b3021/b3075. Dua modul, satu procedure, dua bentuk -
// dan menyeragamkannya membuat salah satunya tidak lagi cocok dengan nomor
// yang sudah beredar.
//
// Dibaca sesudah: polis_periode.go (periodenya), polis_validasi.go.

import (
	"errors"
	"fmt"
	"strings"
)

// Tipe premium list yang punya cabang penomoran.
//
// `[terverifikasi]` Keempatnya, dan HANYA keempatnya, punya langkah sendiri
// di `SubmitPremiumList_Act` (b1783/b1959/b2127/b2295).
const (
	TipePLQuotationRealisasi = "QR"
	TipePLQuotationProposal  = "QP"
	TipePLTreatyProposal     = "TP"
	TipePLTreatyRealisasi    = "TR"
)

// ErrTipePLTanpaCabang - `Type` polis tidak punya cabang penomoran.
//
// ⛔ PENYIMPANGAN SADAR DARI PEGA, dan disengaja. Di Pega `InputData.CARI20`
// tidak direset langkah 5 (yang mereset `TempInputData`, bukan `InputData`),
// jadi tipe di luar keempatnya menghasilkan nomor tanpa huruf tipe - atau,
// lebih buruk, dengan huruf tipe polis SEBELUMNYA yang masih tertinggal di
// halaman. Keduanya nomor yang terlihat sah. Kami menolaknya.
var ErrTipePLTanpaCabang = errors.New(
	"models: Type polis tidak punya cabang penomoran (bukan QR/QP/TP/TR)")

// ErrPeriodeNomorPLTakBerbentuk - periode dari penghitung bukan `MM.YYYY`.
var ErrPeriodeNomorPLTakBerbentuk = errors.New(
	"models: periode penghitung bukan berbentuk MM.YYYY")

// ErrKodeBisnisKosong - `BusinessCode` kosong, padahal nomor memuatnya.
var ErrKodeBisnisKosong = errors.New(
	"models: BusinessCode kosong; nomor premium list memuatnya")

// ErrAwalanProduksiKosong - awalan dari `KODE_PRODUKSI` kosong.
var ErrAwalanProduksiKosong = errors.New(
	"models: awalan produksi kosong; nomor premium list memuatnya")

// KodeTipePL memilih huruf tipe yang masuk ke nomor.
//
// ⚠️ Cabangnya dipilih DARI DATA - satu `switch` atas nilai `Type`, bukan
// empat langkah berurutan yang masing-masing menyalakan dirinya sendiri
// (AC tiket 03). Urutan langkah di Pega adalah cara Pega menulis `switch`.
func KodeTipePL(tipe string) (string, error) {
	switch strings.TrimSpace(tipe) {
	case TipePLQuotationRealisasi:
		return TipePLQuotationRealisasi, nil
	case TipePLQuotationProposal:
		return TipePLQuotationProposal, nil
	case TipePLTreatyProposal:
		return TipePLTreatyProposal, nil
	case TipePLTreatyRealisasi:
		return TipePLTreatyRealisasi, nil
	}
	return "", fmt.Errorf("%w: %q", ErrTipePLTanpaCabang, tipe)
}

// jenisPenghitungPLHarfiah adalah teks b2717, disalin apa adanya.
//
// ⛔ Ini BUKAN daftar pilihan. Ia satu nilai teks yang menjadi bagian kunci
// `GENERATE_SEQUENCE_NUMBER.JENIS`, dan mengubahnya - termasuk mengurutkan
// ulang keempat hurufnya - memindahkan seluruh penomoran ke baris penghitung
// yang lain, yang urutannya mulai dari nol.
const jenisPenghitungPLHarfiah = "QR/QP/TP/TR"

// JenisPenghitungPL menyusun kunci `JENIS` penghitung.
//
// `[terverifikasi]` b2717 `ParamSeq.CARI2 = ParamSeq.HASIL3+"QR/QP/TP/TR"` -
// awalan produksi DITEMPEL di depan teks harfiahnya.
func JenisPenghitungPL(awalan string) string {
	return awalan + jenisPenghitungPLHarfiah
}

// ClassPenghitungPL adalah kunci `CLASS` penghitung.
//
// `[terverifikasi]` b2670 `ParamSeq.CARI1 = pyWorkPage.pxObjClass`.
//
// ⚠️ `[terbuka - pemilik kerja / DBA]` NILAINYA BELUM PASTI, dan sebabnya
// dinyatakan alih-alih disembunyikan di balik konstanta yang kelihatan yakin.
// `pxObjClass` adalah kelas KONKRET saat berjalan, sedangkan korpus ini
// lapisan kerangkanya. Dua kelas hidup berdampingan di sana:
//
//	ASM-FW-GISFW-Work-LIFE    kelas tempat SELURUH rule korpus ini
//	                          didefinisikan (Activity, Harness, Section)
//	RNM-FW-LIFEFW-Work-LIFE   kelas yang merujuk balik ke rule itu
//
// `[terverifikasi]` Jendela pertama: elemen `<pyActivityClass>` di SATU
// berkas, `Section/ShowLifePremiumDetail.xml`.
//
//	grep -o "<pyActivityClass>[^<]*" Section/ShowLifePremiumDetail.xml \
//	  | sed 's/<[^>]*>//' | sort | uniq -c
//
// menghasilkan 107 `ASM-FW-GISFW-Work-LIFE`, 8 `RNM-FW-LIFEFW-Work-LIFE`,
// 8 `ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL`, 2
// `RNM-FW-LIFEFW-Int-LIFE_PREMIUM_DETAIL`. Jendela kedua: SELURUH korpus
// PremiumList Life.
//
//	grep -rho "RNM-FW-LIFEFW-Work-LIFE" . | wc -l    -> 394
//
// ⚠️ KEDUA ANGKA MENGUKUR HAL YANG BERBEDA - rujukan activity di satu section
// versus kemunculan teks di seluruh korpus - jadi keduanya TIDAK saling
// meneguhkan, dan tidak boleh dibaca begitu. Yang mereka tunjukkan bersama
// hanya ini: kedua kelas hidup, dan korpus tidak memberi tahu mana yang
// konkret saat berjalan.
//
// Yang dipakai di bawah adalah yang pertama, sebab itulah kelas yang rule-nya
// sendiri sebut. Bila baris `GENERATE_SEQUENCE_NUMBER` di basis data nyata
// ternyata bertuliskan yang kedua, penomoran akan MULAI DARI SATU dan setiap
// nomor baru bertabrakan dengan nomor lama. Karena itu nilainya dikunci uji,
// bukan hanya ditulis: menggantinya menuntut alasan.
const ClassPenghitungPL = "ASM-FW-GISFW-Work-LIFE"

// PeriodeNomorPL memotong `MM.YYYY` dari penghitung menjadi `MM.YY`.
//
// `[terverifikasi]` b3021 `@substring(ParamSeq.HASIL1,0,2)` dan b3075
// `@substring(ParamSeq.HASIL1,5,7)` - keduanya indeks Java, awal disertakan,
// akhir tidak. Atas `"09.2026"` hasilnya `"09"` dan `"26"`.
//
// ⛔ BENTUKNYA DIPERIKSA, tidak dipotong membabi buta. `@substring` Pega atas
// teks yang lebih pendek melempar; atas teks yang lebih panjang ia memotong
// diam-diam dan menghasilkan nomor yang salah periodenya. Keduanya bukan yang
// kami inginkan: periode yang tidak berbentuk `MM.YYYY` berarti penghitungnya
// mengembalikan sesuatu yang tidak kami mengerti, dan nomor yang lahir
// darinya tidak dapat dipercaya.
func PeriodeNomorPL(mmYYYY string) (string, error) {
	s := strings.TrimSpace(mmYYYY)
	if len(s) != 7 || s[2] != '.' || !angkaSaja(s[0:2]) || !angkaSaja(s[3:7]) {
		return "", fmt.Errorf("%w: %q", ErrPeriodeNomorPLTakBerbentuk, mmYYYY)
	}
	return s[0:2] + "." + s[5:7], nil
}

// angkaSaja menjawab apakah seluruh isinya angka desimal.
func angkaSaja(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// PanjangUrutNomorPL adalah lebar bagian urut nomor.
//
// `[data DBA]` `LPAD(v_seq, 5, '0')` di dalam
// `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`.
const PanjangUrutNomorPL = 5

// RakitNomorPL menyusun `PL_NUMBER` dari potongan-potongannya.
//
// Bentuknya awalan + tipe + kodeBisnis + "." + MM.YY + "." + lima angka -
// b3126.
//
// ⚠️ `awalan` DITERIMA, tidak diambil sendiri: fungsi ini murni supaya bentuk
// nomornya dapat diuji tanpa Oracle. Yang mengambilnya penomornya - sejajar
// dengan `RakitNomorKlaim`.
//
// ⚠️ Urut yang sudah lebih panjang dari lima angka TIDAK dipotong.
// Memotongnya menerbitkan nomor yang bertabrakan - persis yang `LPAD` sendiri
// tidak lakukan.
func RakitNomorPL(awalan, kodeTipe, kodeBisnis, mmYY string, urut int) string {
	u := fmt.Sprintf("%d", urut)
	for len(u) < PanjangUrutNomorPL {
		u = "0" + u
	}
	return awalan + kodeTipe + kodeBisnis + "." + mmYY + "." + u
}

// BahanNomorPL adalah seluruh bahan yang dibutuhkan satu nomor.
//
// ⚠️ Satu tipe, bukan lima parameter yang selalu berjalan bersama. Kelimanya
// harus konsisten satu sama lain, dan tipe yang memuatnya membuat pemanggil
// tidak dapat menukar `Tipe` dengan `KodeBisnis` - keduanya teks pendek, dan
// tertukar keduanya menghasilkan nomor yang tetap terlihat sah.
type BahanNomorPL struct {
	// Awalan dari `KODE_PRODUKSI` - TIDAK pernah konstanta di kode.
	Awalan string
	// Tipe adalah `T_PREMIUM_LIST.TYPE` apa adanya; hurufnya dipilih
	// `KodeTipePL`.
	Tipe string
	// KodeBisnis adalah `T_PREMIUM_LIST.BUSINESS_CODE`.
	KodeBisnis string
	// PeriodeMMYYYY adalah keluaran penghitung, berbentuk `MM.YYYY`.
	PeriodeMMYYYY string
	// Urut adalah nilai penghitung sesudah dinaikkan.
	Urut int
}

// NomorPL merakit nomor dari bahannya, memeriksa setiap bahan lebih dahulu.
//
// ⛔ Nol bahan kosong yang dibiarkan lewat. Nomor yang kehilangan satu
// bagiannya tetap berbentuk nomor - satu yang awalannya hilang terbaca sah
// oleh mata, tersimpan, dan baru ketahuan salah ketika seseorang mencarinya
// dan tidak menemukannya.
func NomorPL(b BahanNomorPL) (string, error) {
	if strings.TrimSpace(b.Awalan) == "" {
		return "", ErrAwalanProduksiKosong
	}
	kodeTipe, err := KodeTipePL(b.Tipe)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(b.KodeBisnis) == "" {
		return "", ErrKodeBisnisKosong
	}
	mmYY, err := PeriodeNomorPL(b.PeriodeMMYYYY)
	if err != nil {
		return "", err
	}
	if b.Urut < 1 {
		return "", fmt.Errorf("models: urut penghitung %d tidak masuk akal", b.Urut)
	}
	return RakitNomorPL(b.Awalan, kodeTipe, strings.TrimSpace(b.KodeBisnis),
		mmYY, b.Urut), nil
}
