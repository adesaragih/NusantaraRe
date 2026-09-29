package models

// Kontrak treaty di dalam tahun treaty - tiket 04 Treaty Contract Out.
//
// Untuk apa berkas ini: gerbang simpan kontrak dan tanggal akhir bawaannya.
//
// Pohon yang ditiru, dibaca 29-09-2026 (nomor baris mentah):
//
//	`Activity/SaveTreatyContract_Act.xml`    langkah 7 (UserID, TglUpdate, tanggal)
//	                                        dan 8 (RDB `SaveMasterTreatyContract_SQL`);
//	                                        langkah 2-6 dan 12 DIKOMENTARI (`//`)
//	`Activity/SetTanggalTreatyContract.xml`  langkah 3-6 tanggal akhir bawaan,
//	                                        langkah 4 b847 tahun mulai = tahun treaty
//
// Dibaca sesudah: tco_tahun.go.

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrKontrakJenisReasuransiKosong - jenis reasuransi wajib (AC 10, 11).
	//
	// ⚠️ Di Pega gerbang ini DIKOMENTARI (`SaveTreatyContract_Act` langkah 2,
	// `ASMMessageReinstype`): kontrak baru dengan tanggal terisi tersimpan
	// tanpa jenis reasuransi. AC tiket 04 menuntutnya - ralat bertanggal.
	ErrKontrakJenisReasuransiKosong = errors.New("models: ReinsTypeID wajib diisi - jenis reasuransi dipilih dari daftar")
	// ErrKontrakMulaiKosong - `SaveTreatyContract_Act` langkah 3
	// (`ASMMessageStartDateKosong`, juga dikomentari) dan AC 53.
	ErrKontrakMulaiKosong = errors.New("models: TreatyStartDate wajib diisi")
	// ErrKontrakAkhirKosong - tanggal akhir wajib (AC 53).
	ErrKontrakAkhirKosong = errors.New("models: TreatyEndDate wajib diisi")
	// ErrKontrakTahunMulaiBeda - `SetTanggalTreatyContract.xml` langkah 4 b847
	// (`ASMMessageStartDate`, HIDUP): tahun tanggal mulai harus sama dengan
	// tahun treaty induknya. Teks pesan Pega tidak diekspor.
	ErrKontrakTahunMulaiBeda = errors.New("models: tahun TreatyStartDate tidak sama dengan tahun treaty")
)

// PeriksaKontrakTreaty menjalankan gerbang simpan kontrak.
//
// `tahunTreaty` adalah `TREATYYEAR` tahun induknya -
// `InputTreatyContractReinsType.TreatyYear` di Pega.
func PeriksaKontrakTreaty(k KontrakTreaty, tahunTreaty string) error {
	if strings.TrimSpace(k.ReinsTypeID) == "" {
		return ErrKontrakJenisReasuransiKosong
	}
	if k.TreatyStartDate.IsZero() {
		return ErrKontrakMulaiKosong
	}
	if k.TreatyEndDate.IsZero() {
		return ErrKontrakAkhirKosong
	}
	if err := PeriksaPeriodeTCO(k.TreatyStartDate, k.TreatyEndDate, "TreatyStartDate", "TreatyEndDate"); err != nil {
		return err
	}
	// ⛔ Perbandingan TEKS, seperti `@substring(TreatyStartDate,0,4) <>
	// TreatyYear` - tahun treaty warisan yang bukan angka tidak pernah sama.
	if mulai := strconv.Itoa(k.TreatyStartDate.Year()); mulai != strings.TrimSpace(tahunTreaty) {
		return fmt.Errorf("%w: TreatyStartDate %s bertahun %s, tahun treaty %s", ErrKontrakTahunMulaiBeda,
			k.TreatyStartDate.Format("2006-01-02"), mulai, strings.TrimSpace(tahunTreaty))
	}
	return nil
}

// AkhirKontrakBawaanTCO menghitung tanggal akhir yang diisikan layar saat
// tanggal mulai dipilih - `SetTanggalTreatyContract.xml` langkah 3-6.
//
// Langkah 3: `JumlahHari = 365`. Langkah 5: `JumlahHari = 366` bila SELURUH
// prasyaratnya lolos. Langkah 6: `TreatyEndDate = TreatyStartDate +
// JumlahHari` hari.
//
// ⛔ DITIRU APA ADANYA, TERMASUK ANOMALINYA (OQ-TCO-10). Prasyarat langkah 5
// memotong cap waktu Pega `yyyyMMdd...` seolah berformat `dd/MM/yyyy`
// (contoh uji penulisnya `01/02/2018`):
//
//	substring(0,2) di 1..31   -> yang terbaca dua digit ABAD, bukan tanggal
//	substring(4,5) di 1..2    -> yang terbaca digit PULUHAN bulan, jadi hanya
//	                             Oktober-Desember yang lolos, bukan Januari-Februari
//
// Maksud penulisnya tampak "mulai Januari-Februari tahun kabisat = 366 hari".
// Yang BERJALAN adalah "mulai Oktober-Desember tahun kabisat = 366 hari".
// Memperbaiki diam-diam berarti tanggal akhir bawaan berbeda dari sistem lama
// tanpa seorang pun memutuskannya; tanggal akhir tetap dapat diubah pemakai.
func AkhirKontrakBawaanTCO(mulai time.Time, tahunTreaty string) time.Time {
	hari := 365
	tahun := strings.TrimSpace(tahunTreaty)
	abad := mulai.Year() / 100
	puluhanBulan := int(mulai.Month()) / 10
	if strconv.Itoa(mulai.Year()) == tahun && abad >= 1 && abad <= 31 &&
		puluhanBulan >= 1 && puluhanBulan <= 2 && tahunKabisatTeks(tahun) {
		hari = 366
	}
	return mulai.AddDate(0, 0, hari)
}

// tahunKabisatTeks - `@if(Tahun%100==0, @if(Tahun%400==0,366,365),
// @if(Tahun%4==0,366,365)) == 366`; teks bukan angka tidak kabisat.
func tahunKabisatTeks(tahun string) bool {
	n, err := strconv.Atoi(tahun)
	if err != nil {
		return false
	}
	if n%100 == 0 {
		return n%400 == 0
	}
	return n%4 == 0
}
