package penomor

import (
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/utils"
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
	if len(s) != 7 || s[2] != '.' || !utils.AngkaSaja(s[0:2]) || !utils.AngkaSaja(s[3:7]) {
		return "", fmt.Errorf("%w: %q", ErrPeriodeNomorPLTakBerbentuk, mmYYYY)
	}
	return s[0:2] + "." + s[5:7], nil
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
