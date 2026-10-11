package models

// Validasi wajib isi CSV peserta baru + templat unduhan (keputusan work owner 10-10-2026: "iya lakukan a dan b").
// ⚠️ MENYIMPANG dari Pega: `SaveCSVEDMLife` tidak mewajibkan kolom apa pun selain acuan PLAN/POLICY_HOLDER.
// CSV tetap TIDAK menghitung premi (keputusan terpisah, ditunda).

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// KolomWajibUmumCSV - wajib ada di judul DAN terisi di setiap baris, semua Type. `NET_PREMIUM` ditambahkan untuk
// SEMUA Type (keputusan work owner 10-10-2026: "buat net_premium wajib untuk semua type") walau di luar rumus rekap;
// varian NET_PREMIUM_REFUND / _RETRO / _REFUND_RETRO tidak diwajibkan.
// PLAN dan POLICY_HOLDER TIDAK wajib (keputusan work owner 10-10-2026): nilainya diambil dari polis (`NilaiPolisCSV`),
// nilai CSV diabaikan - seperti PremiumList Life.
var KolomWajibUmumCSV = []string{
	"NAME_OF_INSURED", "CERTIFICATE_NO", "SEX", "DOB", "CURRENCY", "BEGIN_DATE", "EXPIRED_DATE", "SUM_INSURED",
	"NET_PREMIUM",
}

// ErrDataPolisCSVKosong - PLAN / POLICY_HOLDER kosong di polis DAN di peserta Old kasus: unggahan ditolak seluruhnya
// (pola `DataPolisPeserta.Periksa` PremiumList Life, disalin - tidak diimpor).
var ErrDataPolisCSVKosong = errors.New("models: policy data for participants is empty")

// NilaiPolisCSV - PLAN / POLICY_HOLDER peserta baru (keputusan work owner 10-10-2026, "seperti premium list life"):
// kepala kasus (`BUSINESS_NAME` / `POLICY_HOLDER_NAME`, disalin dari versi polis saat Create Addendum) bila terisi,
// selain itu acuan yang sudah ada (peserta Old pertama: `acuan.Plan` / `acuan.PolicyHolder`) - per kolom. Keduanya
// kosong → `ErrDataPolisCSVKosong` berpesan.
func NilaiPolisCSV(acuan AcuanCSV, businessName, policyHolderName string) (AcuanCSV, error) {
	if v := strings.TrimSpace(businessName); v != "" {
		acuan.Plan = v
	}
	if v := strings.TrimSpace(policyHolderName); v != "" {
		acuan.PolicyHolder = v
	}
	var kosong []string
	if strings.TrimSpace(acuan.PolicyHolder) == "" {
		kosong = append(kosong, "POLICY_HOLDER")
	}
	if strings.TrimSpace(acuan.Plan) == "" {
		kosong = append(kosong, "PLAN")
	}
	switch len(kosong) {
	case 0:
		return acuan, nil
	case 1:
		return acuan, galatCSV(ErrDataPolisCSVKosong,
			"%s is empty on the policy and on its existing participants, so it cannot be filled in for the new participants.", kosong[0])
	}
	return acuan, galatCSV(ErrDataPolisCSVKosong,
		"%s are empty on the policy and on its existing participants, so they cannot be filled in for the new participants.",
		strings.Join(kosong, " and "))
}

// KolomUsiaCSV - minimal satu ada di judul dan minimal satu terisi per baris.
var KolomUsiaCSV = []string{"AGE", "ENTRY_AGE"}

// kolomWajibTipeCSV - wajib per Type = kolom yang dijumlah rumus rekap PREMIUM/BALANCE per Type
// (`repository.rumusPremiBalance`; uji repository memastikan keduanya sama). TP = *_RETRO, TR = *_REFUND_RETRO
// seperti rekap dan grid mata uang layar. TAX, PROF_COMM, CLAIM (rumus rekap QR/QP) SENGAJA tidak wajib (keputusan
// work owner 10-10-2026: "kolom tax, prof_comm, claim tidak perlu") - tidak ada/kosong tersimpan 0 (UangCSVEDM), rekap
// memakai NVL; bila diisi tetap divalidasi formatnya.
var kolomWajibTipeCSV = map[string][]string{
	TypeQR: {"GROSS_PREMIUM", "DEDUCTION", "RI_ADMIN_FEE", "BROKERAGE_FEE"},
	TypeQP: {"GROSS_PREMIUM_REFUND", "CLAIM_AMOUNT", "DEDUCTION_REFUND", "BROKERAGE_FEE_REFUND", "RI_ADMIN_FEE_REFUND"},
	TypeTP: {"GROSS_PREMIUM_RETRO", "DISCOUNT_PREMIUM_RETRO", "RI_ADMIN_FEE_RETRO", "BROKERAGE_FEE_RETRO"},
	TypeTR: {"GROSS_PREMIUM_REFUND_RETRO", "DISCOUNT_PREMIUM_REFUND_RETRO", "RI_ADMIN_FEE_REFUND_RETRO",
		"BROKERAGE_FEE_REFUND_RETRO"},
}

// Pesan validasi wajib isi - kalimat untuk user awam (disederhanakan 10-10-2026).
const (
	// PesanCSVWajib - %s = nama kolom.
	PesanCSVWajib     = "%s is empty. Please fill it in."
	PesanCSVUsiaWajib = "AGE or ENTRY_AGE must be filled in."
	PesanCSVMataUang  = "CURRENCY must be IDR or USD."
	// PesanCSVMataUangPolis - %s = mata uang baris, %s = mata uang polis ("IDR" / "IDR or USD").
	PesanCSVMataUangPolis = "CURRENCY %s is not used in this policy (use %s)."
	// Peserta dianggap sama HANYA bila CERTIFICATE_NO sama (keputusan work owner 10-10-2026); %s = nomor sertifikat
	// baris CSV sesudah TrimSpace.
	PesanCSVSertifikatAda   = `Certificate No "%s" is already registered in this policy.`
	PesanCSVSertifikatGanda = `Certificate No "%s" is used more than once in this file (also in row %d).`
)

// MataUangCSV - mata uang yang diterima CSV.
var MataUangCSV = []string{"IDR", "USD"}

// TipeCSV - Type kasus untuk aturan CSV; Type lain diperlakukan QR (seperti `tipePopup` layar).
func TipeCSV(tipe string) string {
	if _, ada := kolomWajibTipeCSV[strings.TrimSpace(tipe)]; ada {
		return strings.TrimSpace(tipe)
	}
	return TypeQR
}

// KolomWajibTipeCSV - kolom wajib khusus Type.
func KolomWajibTipeCSV(tipe string) []string {
	return append([]string{}, kolomWajibTipeCSV[TipeCSV(tipe)]...)
}

// KolomWajibCSV - seluruh kolom wajib berurutan (umum, AGE, ENTRY_AGE, Type) - kolom depan templat.
func KolomWajibCSV(tipe string) []string {
	h := append([]string{}, KolomWajibUmumCSV...)
	h = append(h, KolomUsiaCSV...)
	return append(h, KolomWajibTipeCSV(tipe)...)
}

// DaftarKolomWajibCSV - daftar layar (modal Upload CSV): umum, "AGE / ENTRY_AGE", Type. Satu sumber dengan validasi.
func DaftarKolomWajibCSV(tipe string) []string {
	h := append([]string{}, KolomWajibUmumCSV...)
	h = append(h, strings.Join(KolomUsiaCSV, " or ")) // "AGE or ENTRY_AGE" (disederhanakan 10-10-2026)
	return append(h, KolomWajibTipeCSV(tipe)...)
}

// KunciSertifikat - pembanding CERTIFICATE_NO: TrimSpace + huruf besar (tidak peka huruf).
func KunciSertifikat(s string) string { return strings.ToUpper(strings.TrimSpace(s)) }

// kolomJudulKurang - kolom wajib yang tidak ada di judul (urut KolomWajibCSV; AGE/ENTRY_AGE sebagai satu butir).
func kolomJudulKurang(lihat map[string]bool, tipe string) []string {
	var kurang []string
	for _, k := range KolomWajibUmumCSV {
		if !lihat[k] {
			kurang = append(kurang, k)
		}
	}
	if !slices.ContainsFunc(KolomUsiaCSV, func(k string) bool { return lihat[k] }) {
		kurang = append(kurang, strings.Join(KolomUsiaCSV, " or "))
	}
	for _, k := range KolomWajibTipeCSV(tipe) {
		if !lihat[k] {
			kurang = append(kurang, k)
		}
	}
	return kurang
}

// periksaWajibBaris - wajib isi, usia, mata uang, sertifikat yang sudah ada (pemeriksaan ganda DALAM berkas di
// layanan karena butuh keadaan antarbaris).
func periksaWajibBaris(nomor int, baris map[string]string, acuan AcuanCSV) []PesanCSV {
	var pesan []PesanCSV
	for _, k := range append(append([]string{}, KolomWajibUmumCSV...), KolomWajibTipeCSV(acuan.Tipe)...) {
		if strings.TrimSpace(baris[k]) == "" {
			pesan = append(pesan, PesanCSV{Baris: nomor, Kolom: k, Pesan: fmt.Sprintf(PesanCSVWajib, k)})
		}
	}
	if !slices.ContainsFunc(KolomUsiaCSV, func(k string) bool { return strings.TrimSpace(baris[k]) != "" }) {
		pesan = append(pesan, PesanCSV{Baris: nomor, Kolom: KolomUsiaCSV[0], Pesan: PesanCSVUsiaWajib})
	}
	if mu := strings.TrimSpace(baris["CURRENCY"]); mu != "" {
		switch {
		case !slices.Contains(MataUangCSV, mu):
			pesan = append(pesan, PesanCSV{Baris: nomor, Kolom: "CURRENCY", Pesan: PesanCSVMataUang})
		case len(acuan.MataUang) > 0 && !slices.Contains(acuan.MataUang, mu):
			pesan = append(pesan, PesanCSV{Baris: nomor, Kolom: "CURRENCY",
				Pesan: fmt.Sprintf(PesanCSVMataUangPolis, mu, strings.Join(acuan.MataUang, " or "))})
		}
	}
	if cn := KunciSertifikat(baris["CERTIFICATE_NO"]); cn != "" && acuan.Sertifikat[cn] {
		pesan = append(pesan, PesanCSV{Baris: nomor, Kolom: "CERTIFICATE_NO", Pesan: fmt.Sprintf(PesanCSVSertifikatAda, strings.TrimSpace(baris["CERTIFICATE_NO"]))})
	}
	return pesan
}

// --- templat ------------------------------------------------------------------

// JudulTemplatCSV - seluruh properti `PemetaanCSV` (72): kolom wajib Type ini paling depan, sisanya urut PemetaanCSV.
func JudulTemplatCSV(tipe string) []string {
	h := KolomWajibCSV(tipe)
	for _, c := range PemetaanCSV {
		if !slices.Contains(h, c.Properti) {
			h = append(h, c.Properti)
		}
	}
	return h
}

// BarisContohCSV - satu baris contoh templat dari nilai peserta (berkunci KOLOM tabel, bentuk basis data: tanggal
// DD/MM/YYYY, angka TM9 bertitik) dalam bentuk yang lolos `RapikanBarisCSV`: uang/desimal dirapikan apd (maks 8
// desimal, tanpa pemisah ribuan; uang kosong = "0" seperti `UangCSVEDM`), STNC/WPC dinormalkan ke DD/MM/YYYY (bentuk
// lain dikosongkan). CERTIFICATE_NO dan NAME_OF_INSURED DIKOSONGKAN (supaya tidak menjadi duplikat).
func BarisContohCSV(judul []string, nilai map[string]string) []string {
	jenis := map[string]KolomCSV{}
	for _, c := range PemetaanCSV {
		jenis[c.Properti] = c
	}
	baris := make([]string, len(judul))
	for i, p := range judul {
		c, ada := jenis[p]
		if !ada || c.Kolom == "" || p == "CERTIFICATE_NO" || p == "NAME_OF_INSURED" {
			continue
		}
		baris[i] = nilaiContohCSV(c, strings.TrimSpace(nilai[c.Kolom]))
	}
	return baris
}

func nilaiContohCSV(c KolomCSV, v string) string {
	switch c.Jenis {
	case CSVUang, CSVDesimal:
		if v == "" {
			if c.Jenis == CSVUang {
				return "0"
			}
			return ""
		}
		d, err := utils.ParseDecimal(v)
		if err != nil {
			return ""
		}
		var q, r apd.Decimal
		if _, err := utils.DecimalContext().Quantize(&q, d, -batasDesimalCSV); err != nil {
			return ""
		}
		r.Reduce(&q)
		return utils.FormatDecimal(&r)
	case CSVTanggal:
		if _, err := time.Parse(bentukTanggalCSV, v); err != nil {
			return ""
		}
		return v
	case CSVTeksTanggal:
		for _, bentuk := range []string{bentukTanggalCSV, time.DateOnly, "20060102"} {
			if len(v) >= len(bentuk) {
				if t, err := time.Parse(bentuk, v[:len(bentuk)]); err == nil {
					return t.Format(bentukTanggalCSV)
				}
			}
		}
		return ""
	}
	return v
}
