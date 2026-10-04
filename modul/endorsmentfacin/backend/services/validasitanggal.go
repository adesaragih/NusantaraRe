package services

// Validasi tanggal endorsement - tiket E05.
//
// Untuk apa berkas ini: tanggal endorsement wajib berada di dalam periode polis
// yang di-endors. Asal (Pega): `Endorsment Fac In/Activity/CheckEDMPolisDate.xml`
// (6 langkah) · `RDBList/GetStartDate.xml` (`Rule-Connect-SQL`:
// `BEGINDATE AS CARI1, ENDDATE AS CARI2 from pooldata.facinproduction`, baris
// SPREAD_DATE terakhir, `FETCH FIRST 1 ROW ONLY`).
//
// Dibaca sesudah: bukakasus.go.
//
// Alur korpus `[terverifikasi]`:
//
//	1  CARI1 = @FormatDateTime(EndorsementDate,"dd/MM/yyyy","Asia/Jakarta"); EdmDate = CARI1
//	2  PolicyNo == <literal> || @contains(PolicyNo,"RNML")  → transisi 1, lompat ke label END
//	3  RDB-List GetStartDate
//	4  per baris: StartDate = @replaceAll(CARI1, @substring(CARI1,8,11), "T05"); idem EndDate
//	   (4.1 `pyStepsPreCondition=false`: tetap jalan, P-11)
//	   CARI12 = @FormatDateTime(StartDate,"dd/MM/yyyy","Asia/Jakarta"); CARI20 = CARI1==CARI12
//	5  CARI20=="true" → lewati (3); selain itu:
//	   @CompareDates(@toDate(StartDate),@toDate(EdmDate))=="true" ||
//	   @CompareDates(@toDate(EdmDate),@toDate(EndDate))=="true" → pesan galat
//
// ⚠️ `[dugaan]` yang dipakai, semuanya dari bentuk aturannya sendiri:
//   - BEGINDATE/ENDDATE terbaca sebagai teks DateTime Pega
//     `yyyyMMddTHHmmss.SSS GMT` - satu-satunya bentuk yang membuat
//     `substring(8,11)` = `"Thh"`. Akibatnya StartDate/EndDate = tanggal
//     KALENDER GMT polis pukul 05:00 GMT (= 12:00 WIB, tanggal sama).
//   - `@CompareDates(a,b)=="true"` = a SESUDAH b - satu-satunya bacaan yang
//     cocok dengan pesan "EDM date cannot be outside the period".
//   - `@toDate` membaca tanggal kalender kedua bentuk itu.
//   - "Asia/Jakarta" = UTC+7 tetap (zona itu tanpa DST sejak 1964).
//
// ⛔ K-046 - diport apa adanya: polis yang tersimpan 17:00 GMT (00:00 WIB hari
// berikutnya) dibandingkan dengan tanggal GMT-nya, yaitu SEHARI LEBIH AWAL
// dari tanggal mulai WIB; tanggal EDM pada hari itu lolos lewat langkah 5
// baris 1. ⛔ Premis kasus uji `K046_ValidasiTanggal_PerbandinganSubstring8Karakter`
// tidak sepenuhnya didukung korpus: perbandingan `@substring(…,0,8)` hanya ada
// di langkah 4.2 "temporary utk monitoring" (CARI17) dan di deskripsi 4.1 -
// keputusannya memakai CARI20 dan `@CompareDates`.

import (
	"errors"
	"strings"
	"time"
)

// pesanTanggalDiLuarPeriode - `Local.ErrMsg` langkah 1.
const pesanTanggalDiLuarPeriode = "EDM date cannot be outside the period"

// polaPolisDilewati - langkah 2 `@contains(Primary.PolicyNo,"RNML")`: bypass
// validasi tanggal. Ini POLA awalan, bukan nomor polis produksi - aman sebagai
// literal.
const polaPolisDilewati = "RNML"

var jakarta = time.FixedZone("Asia/Jakarta", 7*3600)

// ErrPeriodeProduksiTidakAda - `GetStartDate` tanpa baris: StartDate kosong,
// dan perilaku `@toDate("")` Pega belum terverifikasi (`[pertanyaan terbuka]`,
// laporan "Menunggu work owner": arti `@toDate` atas "" dan atas `dd/MM/yyyy`).
var ErrPeriodeProduksiTidakAda = errors.New("endorsement: periode polis di tabel produksi tidak ditemukan")

// PeriodeProduksi - satu baris `GetStartDate`.
type PeriodeProduksi struct {
	BeginDate, EndDate time.Time
}

// MasukanValidasiTanggal - masukan `CheckEDMPolisDate`.
type MasukanValidasiTanggal struct {
	PolicyNo           string
	TanggalEndorsement time.Time
	// Periode - hasil `GetStartDate`; nil = tanpa baris.
	Periode *PeriodeProduksi
	// Pengecualian - konfigurasi pengganti literal nomor polis langkah 2.
	Pengecualian DaftarPolis
}

// HasilValidasiTanggal - keluaran bisnis validasi.
type HasilValidasiTanggal struct {
	// Dilewati - langkah 2 melompat ke END (pengecualian).
	Dilewati bool
	// Ditolak - langkah 5 memasang pesan galat (`CARIGROUP = "Error"`).
	Ditolak bool
	Pesan   string
	// Medan - medan tempat pesan 5.2 dipasang: `Primary.EndorsementDate`.
	Medan string
}

// ValidasiTanggalEndorsemen - `CheckEDMPolisDate`.
func ValidasiTanggalEndorsemen(m MasukanValidasiTanggal) (HasilValidasiTanggal, error) {
	// 2
	if m.Pengecualian.Memuat(m.PolicyNo) || strings.Contains(m.PolicyNo, polaPolisDilewati) {
		return HasilValidasiTanggal{Dilewati: true}, nil
	}
	// 3
	if m.Periode == nil {
		return HasilValidasiTanggal{}, ErrPeriodeProduksiTidakAda
	}
	edm := tanggalKalender(m.TanggalEndorsement.In(jakarta))
	mulai := tanggalKalender(m.Periode.BeginDate.UTC())
	akhir := tanggalKalender(m.Periode.EndDate.UTC())
	// 5 baris 1 - CARI20: tanggal EDM (dd/MM/yyyy, Jakarta) = tanggal mulai.
	if edm.Equal(mulai) {
		return HasilValidasiTanggal{}, nil
	}
	// 5 baris 2.
	if mulai.After(edm) || edm.After(akhir) {
		return HasilValidasiTanggal{Ditolak: true, Pesan: pesanTanggalDiLuarPeriode, Medan: "EndorsementDate"}, nil
	}
	return HasilValidasiTanggal{}, nil
}

// tanggalKalender - tanggal pada zona waktu nilai itu, tanpa jam.
func tanggalKalender(w time.Time) time.Time {
	y, b, d := w.Date()
	return time.Date(y, b, d, 0, 0, 0, 0, time.UTC)
}
