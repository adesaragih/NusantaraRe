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
	ErrKontrakJenisReasuransiKosong = errors.New("models: ReinsTypeID is required - choose a reinsurance type from the list")
	// ErrKontrakMulaiKosong - `SaveTreatyContract_Act` langkah 3
	// (`ASMMessageStartDateKosong`, juga dikomentari) dan AC 53.
	ErrKontrakMulaiKosong = errors.New("models: TreatyStartDate is required")
	// ErrKontrakAkhirKosong - tanggal akhir wajib (AC 53).
	ErrKontrakAkhirKosong = errors.New("models: TreatyEndDate is required")
	// ErrKontrakTahunMulaiBeda - `SetTanggalTreatyContract.xml` langkah 4 b847
	// (`ASMMessageStartDate`, HIDUP): tahun tanggal mulai harus sama dengan
	// tahun treaty induknya. Teks pesan Pega tidak diekspor.
	ErrKontrakTahunMulaiBeda = errors.New("models: the TreatyStartDate year differs from the treaty year")
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
		return fmt.Errorf("%w: TreatyStartDate %s is in year %s, treaty year %s", ErrKontrakTahunMulaiBeda,
			k.TreatyStartDate.Format("2006-01-02"), mulai, strings.TrimSpace(tahunTreaty))
	}
	return nil
}

// AkhirKontrakBawaanTCO menghitung tanggal akhir yang diisikan layar saat
// tanggal mulai dipilih: MULAI + 1 TAHUN KALENDER [keputusan work owner
// 29-09-2026, OQ-TCO-10 ditutup].
//
// ⛔ PENYIMPANGAN SADAR dari `SetTanggalTreatyContract.xml` (dipanggil
// `InputTreatyContractReinsType.xml` b3007): langkah 3 `JumlahHari = 365`
// (b620-b621), langkah 5 `JumlahHari = 366` bila prasyaratnya lolos
// (b918-b1191), langkah 6 `TreatyEndDate = TreatyStartDate + JumlahHari`
// hari (b1273-b1274). Prasyarat langkah 5 memotong cap waktu `yyyyMMdd...`
// seolah `dd/MM/yyyy`, sehingga yang berjalan adalah "mulai Oktober-Desember
// tahun kabisat = 366 hari" - anomali yang TIDAK lagi ditiru.
//
// Semantik `ADD_MONTHS(mulai, 12)`: hari terakhir bulan tetap hari terakhir
// bulan (28 Februari tahun biasa -> 29 Februari tahun kabisat), dan tanggal
// yang tidak ada di bulan tujuan dijepit ke akhir bulan (29 Februari -> 28
// Februari). Tanggal akhir tetap dapat diubah pemakai.
func AkhirKontrakBawaanTCO(mulai time.Time) time.Time {
	y, m, d := mulai.Date()
	akhirBulan := func(tahun int) int { return time.Date(tahun, m+1, 0, 0, 0, 0, 0, time.UTC).Day() }
	hari := d
	if d == akhirBulan(y) || d > akhirBulan(y+1) {
		hari = akhirBulan(y + 1)
	}
	return time.Date(y+1, m, hari, mulai.Hour(), mulai.Minute(), mulai.Second(), mulai.Nanosecond(), mulai.Location())
}
