package models

// Gerbang tahun treaty - tiket 03 Treaty Contract Out.
//
// `[terverifikasi]` Validasi existing sangat tipis:
//
//	Activity/SaveTreatyYear_Act.xml  b388 prasyarat `InputTreatyYear.TreatyGroupID==""`
//	                                 b411 prasyarat `InputTreatyYear.TreatyYear==""`
//	                                 b434 prasyarat `@Default.isNumber(InputTreatyYear.TreatyYear)`
//	Activity/CheckYear.xml           b335 `@Default.isNumber(InputTreatyYear.TreatyYear)`
//	                                 b268 `pyMessageLabel CheckYearly` - teks pesannya
//	                                 TIDAK diekspor `[terbuka]`
//
// Dua gerbang lain adalah TAMBAHAN SADAR `[keputusan work owner]`, bukan tiruan:
// masa berlaku yang berakhir sebelum dimulai DITOLAK (AC 9), dan kombinasi
// (StartDate, EndDate, TreatyGroupID) yang sudah ada DITOLAK (AC 73) - yang
// kedua menuntut basis data dan hidup di services.
//
// Pesan menyebut MEDANNYA (nama kolom, dan label layar bila ada).

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/inti/utils"
)

var (
	// ErrTahunTreatyGrupKosong - b388.
	ErrTahunTreatyGrupKosong = errors.New("models: Treaty Group (TREATYGROUPID) wajib diisi")
	// ErrTahunTreatyTahunKosong - b411.
	ErrTahunTreatyTahunKosong = errors.New("models: Underwriting Year (TREATYYEAR) wajib diisi")
	// ErrTahunTreatyBukanAngka - b434 / CheckYear b335.
	ErrTahunTreatyBukanAngka = errors.New("models: Underwriting Year (TREATYYEAR) harus angka")
	// ErrPeriodeTerbalik - AC 9: masa berlaku berakhir sebelum dimulai.
	ErrPeriodeTerbalik = errors.New("models: masa berlaku berakhir sebelum dimulai")
)

// PeriksaPeriodeTCO menolak akhir yang mendahului mulai (AC 9).
//
// Salah satu KOSONG bukan pelanggaran - kolomnya nullable (ADR-U-0027) dan
// wajib-isinya gerbang lain. Pesan menyebut kedua medan dan nilainya.
func PeriksaPeriodeTCO(mulai, akhir time.Time, namaMulai, namaAkhir string) error {
	if mulai.IsZero() || akhir.IsZero() {
		return nil
	}
	if akhir.Before(mulai) {
		return fmt.Errorf("%w: %s %s mendahului %s %s", ErrPeriodeTerbalik,
			namaAkhir, utils.FormatTanggal(akhir), namaMulai, utils.FormatTanggal(mulai))
	}
	return nil
}

// `@Default.isNumber` ditiru dengan `angkaSaja` (polis_nomor.go): hanya digit.
//
// ⚠️ `isNumber` Pega menerima pula desimal; kode tahun tidak pernah
// berkoma, dan menerimanya berarti menyimpan "2026.5" sebagai tahun.

// PeriksaTahunTreaty menjalankan seluruh gerbang murni tahun treaty.
//
// Urutannya urutan langkah activity: grup, tahun kosong, tahun angka, lalu
// gerbang tambahan periode.
func PeriksaTahunTreaty(t TahunTreaty) error {
	if strings.TrimSpace(t.TreatyGroupID) == "" {
		return ErrTahunTreatyGrupKosong
	}
	tahun := strings.TrimSpace(t.TreatyYear)
	if tahun == "" {
		return ErrTahunTreatyTahunKosong
	}
	if !utils.AngkaSaja(tahun) {
		return fmt.Errorf("%w: %q", ErrTahunTreatyBukanAngka, t.TreatyYear)
	}
	return PeriksaPeriodeTCO(t.StartDate, t.EndDate, "StartDate", "EndDate")
}
