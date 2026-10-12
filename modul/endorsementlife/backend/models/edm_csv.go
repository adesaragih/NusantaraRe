package models

// Unggahan CSV endorsement (tiket 07) - `Activity/SaveCSVEDMLife.xml`
// langkah 4.1 b840: 72 penetapan dari baris unggahan
// (`TempWorkPage.ListLifePremiumDetailUpload`, b733) ke peserta kasus.
//
// ⛔ PEMISAH UANG. Pega menormalkan koma menjadi titik
// (`@divide(@toDecimal(@replaceAll(.X,",",".")),1,20)`, b866 …) - satu koma
// adalah titik desimal, dan itu ditiru. Nilai berpemisah ribuan menjadi
// `1.234.567.89` di Pega (bukan angka); di sini ia DITOLAK dengan pesan
// berbaris dan berkolom, bukan diterima diam-diam (AC 38).
//
// ⛔ Tanggal `DD/MM/YYYY` - bentuk `To_date(…, 'DD/MM/YYYY')` yang
// `RDBList/SaveMasterLPDet.xml` pakai untuk properti yang sama (b233 …),
// sama dengan unggahan PremiumList.

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// JenisCSV - cara satu nilai CSV diurai sebelum disimpan.
type JenisCSV int

const (
	// CSVTeks - apa adanya (`.X`), paling banyak 255 byte (`VARCHAR2(255)`).
	CSVTeks JenisCSV = iota
	// CSVUang - `@divide(@toDecimal(@replaceAll(.X,",",".")),1,20)`; kosong = 0 (`@toDecimal("")`).
	CSVUang
	// CSVDesimal - `.FACTOR` apa adanya: desimal bertitik, tanpa normalisasi koma.
	CSVDesimal
	// CSVBulat - kolom `NUMBER(5)`.
	CSVBulat
	// CSVTanggal - kolom `DATE`, teks `DD/MM/YYYY`.
	CSVTanggal
	// CSVTeksTanggal - `STNC`/`WPC`: `VARCHAR2` di tabel peserta, `To_date` di
	// `SaveMasterLPDet` - bentuknya diperiksa, disimpan sebagai teks.
	CSVTeksTanggal
)

// KolomCSV - satu penetapan langkah 4.1.
type KolomCSV struct {
	// Properti - judul kolom CSV = properti baris unggahan (`.X`).
	Properti string
	// Kolom - kolom `T_PREMIUM_LIST_DETAIL`; kosong = tabel tidak punya
	// kolomnya (juga tidak di `SaveMasterLPDet`), nilai tidak disimpan.
	Kolom string
	Jenis JenisCSV
}

func k(p string, j JenisCSV) KolomCSV { return KolomCSV{Properti: p, Kolom: p, Jenis: j} }

// PemetaanCSV - 72 penetapan 4.1, urut korpus (b865 … b2395).
var PemetaanCSV = []KolomCSV{
	k("GROSS_PREMIUM", CSVUang), k("NET_PREMIUM", CSVUang), k("SHARE_NUSANTARA_RE", CSVUang), k("SUM_INSURED", CSVUang),
	k("AGE", CSVBulat), k("BEGIN_DATE", CSVTanggal), k("CERTIFICATE_NO", CSVTeks), k("DOB", CSVTanggal),
	k("EXPIRED_DATE", CSVTanggal), k("EFFECTIVE_DATE", CSVTanggal), k("LAPSE_DATE", CSVTanggal), k("PASSED_PERIOD", CSVBulat),
	k("DESCRIPTION", CSVTeks), k("MEDICAL_STATUS", CSVTeks), k("SUM_REASURED", CSVUang), k("EM_PERCENT", CSVUang),
	k("NAME_OF_INSURED", CSVTeks), k("POLICY_NO", CSVTeks), k("PERIOD_MM", CSVBulat), k("PERIOD_YY", CSVBulat),
	k("PLAN", CSVTeks), k("POLICY_HOLDER", CSVTeks), k("COMM", CSVUang), k("GROSS_PREMIUM_REFUND", CSVUang),
	k("NET_PREMIUM_REFUND", CSVUang), k("COMM_REFUND", CSVUang), k("BROKERAGE_FEE_REFUND", CSVUang), k("OVR_COMM_REFUND", CSVUang),
	k("TAX_REFUND", CSVUang), k("SEX", CSVTeks), k("CEDING_RETENTION", CSVUang), k("CLAIM", CSVUang),
	k("TAX", CSVUang), k("BROKERAGE_FEE", CSVUang), k("OVR_COMM", CSVUang), k("CURRENCY", CSVTeks),
	k("PROF_COMM", CSVUang), k("SHARE_RETRO", CSVUang), k("GROSS_PREMIUM_RETRO", CSVUang), k("DISCOUNT_PREMIUM_RETRO", CSVUang),
	k("OVR_COMM_RETRO", CSVUang), k("BROKERAGE_FEE_RETRO", CSVUang), k("NET_PREMIUM_RETRO", CSVUang),
	k("GROSS_PREMIUM_REFUND_RETRO", CSVUang), k("DISCOUNT_PREMIUM_REFUND_RETRO", CSVUang), k("OVR_COMM_REFUND_RETRO", CSVUang),
	k("BROKERAGE_FEE_REFUND_RETRO", CSVUang), k("NET_PREMIUM_REFUND_RETRO", CSVUang), k("CLAIM_AMOUNT", CSVUang),
	k("WPC", CSVTeksTanggal), k("STNC", CSVTeksTanggal), k("ENTRY_AGE", CSVBulat), k("CURRENT_AGE", CSVBulat),
	k("GROSS_VALUATION_BEGIN_DATE", CSVTanggal), k("GROSS_VALUATION_EXPIRED_DATE", CSVTanggal),
	{Properti: "UW_STATUS", Jenis: CSVTeks}, {Properti: "SUM_AT_RISK", Jenis: CSVUang},
	k("RATE", CSVUang), k("DEDUCTION", CSVUang), k("RI_ADMIN_FEE", CSVUang), {Properti: "REMAINING_PERIOD", Jenis: CSVBulat},
	k("FACTOR", CSVDesimal), k("DEDUCTION_REFUND", CSVUang), k("RI_ADMIN_FEE_REFUND", CSVUang),
	// `SaveMasterLPDet` b233/b234: properti `RETROCESSION_*` → kolom `RETRO_VALUATION_*`.
	{Properti: "RETROCESSION_VALUATION_BEGIN_DATE", Kolom: "RETRO_VALUATION_BEGIN_DATE", Jenis: CSVTanggal},
	{Properti: "RETROCESSION_VALUATION_EXPIRED_DATE", Kolom: "RETRO_VALUATION_EXPIRED_DATE", Jenis: CSVTanggal},
	k("SHARE_NUSANTARA_RE_GROSS", CSVUang), k("RETROCEDED_SHARE", CSVUang), k("RI_ADMIN_FEE_RETRO", CSVUang),
	k("RI_ADMIN_FEE_REFUND_RETRO", CSVUang), k("SUM_AT_RISK_GROSS", CSVUang), k("SUM_AT_RISK_RETRO", CSVUang),
}

// KolomCSVTersimpan - penetapan yang punya kolom tabel, urut `PemetaanCSV`.
func KolomCSVTersimpan() []KolomCSV {
	var hasil []KolomCSV
	for _, c := range PemetaanCSV {
		if c.Kolom != "" {
			hasil = append(hasil, c)
		}
	}
	return hasil
}

// AcuanCSV - `PremiumListDetail(1)` (prakondisi 4.1 b2466/b2495): peserta
// pertama kasus di urutan grid, di luar baris `New` (dibuang langkah 2.1).
type AcuanCSV struct {
	Plan         string
	PolicyHolder string
	// Validasi wajib isi (keputusan work owner 10-10-2026, `edm_csvwajib.go`):
	// Tipe - `.Type` kasus (kolom wajib per Type); MataUang - mata uang peserta Old kasus (CURRENCY harus salah
	// satunya; kosong = tidak diperiksa); Sertifikat - CERTIFICATE_NO peserta Old kasus (`KunciSertifikat`).
	Tipe       string
	MataUang   []string
	Sertifikat map[string]bool
}

// PesanCSV - satu penolakan berbaris dan berkolom (AC 37). `Baris` 1 = baris
// data pertama di bawah judul.
type PesanCSV struct {
	Baris int    `json:"baris"`
	Kolom string `json:"kolom"`
	Pesan string `json:"pesan"`
}

const (
	bentukTanggalCSV = "02/01/2006"
	batasTeksCSV     = 255
	batasDesimalCSV  = 8  // `NUMBER(38,8)`
	batasBulatCSV    = 30 // digit di depan titik `NUMBER(38,8)`
	batasNumber5     = 99999
)

var (
	// ErrUangRibuan - nilai uang berpemisah ribuan, atau berpemisah ganda.
	ErrUangRibuan = errors.New("models: thousands separators are not accepted; use one decimal separator (\".\" or \",\")")
	// ErrUangTakTerurai - bukan angka desimal.
	ErrUangTakTerurai = errors.New("models: not a decimal number")
	// ErrUangTerlaluPresisi - lebih dari 8 angka di belakang koma: Oracle akan membulatkannya diam-diam.
	ErrUangTerlaluPresisi = errors.New("models: more than 8 decimal places would be rounded silently")
	polaDesimal           = regexp.MustCompile(`^-?(\d+)(?:\.(\d+))?$`)
)

// desimalCSV memeriksa teks desimal bertitik.
func desimalCSV(s string) (string, error) {
	m := polaDesimal.FindStringSubmatch(s)
	if m == nil {
		return "", fmt.Errorf("%w: %q", ErrUangTakTerurai, s)
	}
	if len(m[2]) > batasDesimalCSV || len(strings.TrimLeft(m[1], "0")) > batasBulatCSV {
		return "", fmt.Errorf("%w: %q", ErrUangTerlaluPresisi, s)
	}
	return s, nil
}

// UangCSVEDM - satu nilai uang 4.1; kosong menjadi "0". ⚠️ Satu pemisah
// selalu titik desimal seperti Pega - `12.500` = 12,5 (ambigu bagi notasi
// Indonesia, OQ-EDM-020); dua pemisah atau lebih ditolak.
func UangCSVEDM(teks string) (string, error) {
	s := strings.TrimSpace(teks)
	if s == "" {
		return "0", nil
	}
	koma, titik := strings.Count(s, ","), strings.Count(s, ".")
	if koma+titik > 1 {
		return "", fmt.Errorf("%w: %q", ErrUangRibuan, teks)
	}
	return desimalCSV(strings.Replace(s, ",", ".", 1))
}

// urai mengurai satu nilai menurut jenisnya; kosong = "" (NULL) kecuali uang. Galatnya berkalimat layar
// (`edm_csvpesan.go`, 10-10-2026) yang menyebut nama kolom.
func urai(c KolomCSV, teks string) (string, error) {
	s := strings.TrimSpace(teks)
	switch c.Jenis {
	case CSVUang:
		v, err := UangCSVEDM(s)
		switch {
		case errors.Is(err, ErrUangRibuan):
			return "", fmt.Errorf(PesanCSVRibuan, c.Properti)
		case errors.Is(err, ErrUangTerlaluPresisi):
			return "", fmt.Errorf(PesanCSVDesimal, c.Properti)
		case err != nil:
			return "", fmt.Errorf(PesanCSVBukanAngka, c.Properti)
		}
		return v, nil
	case CSVDesimal:
		if s == "" {
			return "", nil
		}
		// FACTOR tidak dinormalkan koma di Pega: koma ditolak.
		v, err := desimalCSV(s)
		switch {
		case strings.Contains(s, ","), errors.Is(err, ErrUangTakTerurai):
			return "", fmt.Errorf(PesanCSVFactor, c.Properti)
		case err != nil:
			return "", fmt.Errorf(PesanCSVDesimal, c.Properti)
		}
		return v, nil
	case CSVBulat:
		if s == "" {
			return "", nil
		}
		n, err := strconv.Atoi(s)
		if err != nil || n < -batasNumber5 || n > batasNumber5 {
			return "", fmt.Errorf(PesanCSVBulat, c.Properti)
		}
		return strconv.Itoa(n), nil
	case CSVTanggal, CSVTeksTanggal:
		if s == "" {
			return "", nil
		}
		t, err := time.Parse(bentukTanggalCSV, s)
		if err != nil {
			return "", fmt.Errorf(PesanCSVTanggal, c.Properti)
		}
		if c.Jenis == CSVTeksTanggal {
			return t.Format(bentukTanggalCSV), nil
		}
		return t.Format(time.DateOnly), nil
	}
	if len(s) > batasTeksCSV {
		return "", fmt.Errorf(PesanCSVTeksPanjang, c.Properti)
	}
	return s, nil
}

// RapikanBarisCSV mengurai satu baris CSV (kunci = judul huruf besar) menjadi
// nilai siap simpan berkunci kolom tabel, beserta seluruh penolakannya.
// Tanggal `DATE` keluar sebagai `YYYY-MM-DD`.
func RapikanBarisCSV(nomor int, baris map[string]string, acuan AcuanCSV) (map[string]string, []PesanCSV) {
	nilai := map[string]string{}
	var pesan []PesanCSV
	for _, c := range PemetaanCSV {
		// PLAN / POLICY_HOLDER dari polis, seperti PremiumList Life `TerapkanDataPolisKePeserta` (keputusan work owner
		// 10-10-2026, menyimpang dari Pega 4.1 b2466/b2495): nilai CSV DIABAIKAN - tidak divalidasi, selalu ditimpa
		// nilai polis (`AcuanCSV.Plan` / `PolicyHolder`, lihat `NilaiPolisCSV`).
		switch c.Properti {
		case "PLAN":
			nilai[c.Kolom] = strings.TrimSpace(acuan.Plan)
			continue
		case "POLICY_HOLDER":
			nilai[c.Kolom] = strings.TrimSpace(acuan.PolicyHolder)
			continue
		}
		v, err := urai(c, baris[c.Properti])
		if err != nil {
			pesan = append(pesan, PesanCSV{Baris: nomor, Kolom: c.Properti, Pesan: err.Error()})
			continue
		}
		if c.Kolom != "" {
			nilai[c.Kolom] = v
		}
	}
	// Wajib isi, mata uang, sertifikat yang sudah ada (keputusan work owner 10-10-2026).
	pesan = append(pesan, periksaWajibBaris(nomor, baris, acuan)...)
	return nilai, pesan
}

// ErrJudulCSV - baris judul tidak dapat dipakai; seluruh berkas ditolak.
var ErrJudulCSV = errors.New("models: invalid CSV header")

// JudulCSV menormalkan judul (huruf besar, tanpa BOM) dan menolak judul
// ganda serta berkas tanpa kolom wajib Type ini (`kolomJudulKurang`, keputusan work owner 10-10-2026; dulu hanya
// `PLAN`/`POLICY_HOLDER` - acuan 4.1) - SELURUH kolom yang kurang disebut sekaligus. Judul yang bukan properti 4.1 DIABAIKAN dan dikembalikan agar
// tampil di tinjauan - unduhan `Generate Data Detail` korpus sendiri memuat
// `STATUS` dan `OVRR_COMM` (`GenerateDataDtlLife_act` b276) yang tidak
// dipetakan 4.1; menolaknya membuat berkas Pega tak dapat diunggah ulang.
func JudulCSV(judul []string, tipe string) ([]string, []string, error) {
	dikenal := map[string]bool{}
	for _, c := range PemetaanCSV {
		dikenal[c.Properti] = true
	}
	bom := string(rune(0xFEFF))
	hasil := make([]string, len(judul))
	lihat := map[string]bool{}
	var asing []string
	for i, j := range judul {
		n := strings.ToUpper(strings.TrimSpace(strings.TrimPrefix(j, bom)))
		if lihat[n] {
			// Kalimat layar (10-10-2026); aturannya tetap: judul ganda (termasuk dua judul kosong) menolak berkas.
			if n == "" {
				return nil, nil, galatCSV(ErrJudulCSV, PesanCSVJudulKosong, i+1)
			}
			return nil, nil, galatCSV(ErrJudulCSV, PesanCSVJudulGanda, n)
		}
		lihat[n] = true
		if !dikenal[n] {
			asing = append(asing, n)
		}
		hasil[i] = n
	}
	if kurang := kolomJudulKurang(lihat, tipe); len(kurang) > 0 {
		return nil, nil, galatCSV(ErrJudulCSV, PesanCSVJudulKurang, strings.Join(kurang, ", "))
	}
	return hasil, asing, nil
}

// BarisCSV - satu baris lolos validasi: nomor baris data dan nilai berkunci
// kolom tabel (`RapikanBarisCSV`).
type BarisCSV struct {
	Nomor int
	Nilai map[string]string
}

// AksiTambahCSV - komentar jejak `Add CSV Data` (label tombol b10405).
const AksiTambahCSV = "Add CSV Data"
