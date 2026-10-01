package services

// Porsi periode sebelum dan sesudah tanggal endorsement - tiket E06.
//
// Untuk apa berkas ini: `SetValueToEDMWork` langkah 15 ("Count StartProRate
// and EndProRate").
//
// Dibaca sesudah: beforeimage.go.
//
// ⚠️ K-048 §8.1: nilai yang disetel di sini SEMENTARA. `CountPremiEDM_DT`
// menghitung ulang dan menimpa `ProrateEDMEnd` sebelum jalur produksi
// membacanya. Urutan before-image → selisih WAJIB dijaga; nilai ini tidak
// boleh di-cache sebagai final.
//
// ⛔ Dua hal di Pega yang BELUM terverifikasi, dan karena itu tidak ditebak:
//
//  1. SATUAN selisih DateTime. `[terverifikasi]` Ketiga lokal (`EdmToStart`,
//     `EdmToEnd`, `TotalPeriod`) dideklarasikan `int` di `pyLocalParameters`
//     `SetValueToEDMWork.xml`. `[dugaan]` Satuannya HARI: rule lain membagi
//     selisih DateTime yang sama dengan 365
//     (`setProRatePercent_Act.xml`, `IsThereAnyObjectLocation_Act.xml`:
//     `@Math.divide((…EndDateTime-…StartDateTime),365,20)*100`). Cara Pega
//     memotong pecahan hari ke `int` tidak diketahui - selisih yang bukan
//     hari bulat DITOLAK.
//  2. Mode pembulatan `@Math.divide(…, …, 20)`. Hasil yang tidak eksak pada 20
//     desimal DITOLAK - sikap yang sama dengan keputusan work owner
//     30-09-2026 untuk premi NB.

import (
	"errors"
	"fmt"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/endorsmentfacin/backend/models"
)

// desimalPorsiPeriode - `@Math.divide(…, …, 20)`: 20 digit di belakang koma.
// Juga menjadi `uang.Ratio.Scale`, yang di proyek ini berarti jumlah desimal
// (dicatat `nbfacin`, keputusan work owner 30-09-2026).
const desimalPorsiPeriode = 20

const sehari = 24 * time.Hour

var (
	// ErrModePembulatanBelumTerverifikasi - membagi ke 20 desimal akan
	// membuang digit bukan-nol, padahal mode pembulatan `@Math.divide` Pega
	// (half-up, half-even, potong) belum terverifikasi.
	ErrModePembulatanBelumTerverifikasi = errors.New("endorsement: porsi periode tidak eksak pada 20 desimal, padahal mode pembulatan @Math.divide belum terverifikasi")
	// ErrSatuanSelisihWaktuBelumTerverifikasi - selisih tanggal bukan hari
	// bulat, sehingga hasil Pega bergantung pada satuan selisih DateTime dan
	// cara pemotongannya ke `int`, keduanya belum terverifikasi.
	ErrSatuanSelisihWaktuBelumTerverifikasi = errors.New("endorsement: selisih tanggal porsi periode bukan hari bulat; satuan selisih DateTime Pega dan pemotongannya ke int belum terverifikasi")
	// ErrTanggalPorsiPeriodeKosong - salah satu tanggal yang dikurangkan
	// tidak terisi. Perilaku Pega atas pengurangan tanggal kosong belum
	// terverifikasi.
	ErrTanggalPorsiPeriodeKosong = errors.New("endorsement: tanggal untuk porsi periode tidak terisi")
)

// hitungPorsiPeriode - langkah 15, bergerbang
// `@SizeOfPropertyList(newWorkPage.OfferFacIn.OldData.CurrencyList)>=1`
// (WhenFalse = lewati langkah).
//
//	Local.EdmToStart  = QuotationData.EdmDate − OldData.PolicyData.StartDateTime
//	Local.EdmToEnd    = OldData.PolicyData.EndDateTime − QuotationData.EdmDate
//	Local.TotalPeriod = OldData.PolicyData.EndDateTime − OldData.PolicyData.StartDateTime
//	Local.TotalPeriod = @if(Local.TotalPeriod = 0, 1, Local.TotalPeriod)
//	.OfferFacIn.ProrateStartEDM = @Math.divide(Local.EdmToStart, Local.TotalPeriod, 20)
//	.OfferFacIn.ProrateEDMEnd   = @Math.divide(Local.EdmToEnd,   Local.TotalPeriod, 20)
//
// Galat tidak mengubah agregat: kedua rasio tetap kosong.
func hitungPorsiPeriode(o *models.OfferFacIn) error {
	if len(o.OldData.CurrencyList) < 1 {
		return nil
	}
	edm, mulai, akhir := o.QuotationData.EdmDate, o.OldData.PolicyData.StartDateTime, o.OldData.PolicyData.EndDateTime
	if edm.IsZero() || mulai.IsZero() || akhir.IsZero() {
		return ErrTanggalPorsiPeriodeKosong
	}
	keMulai, err := hariBulat(edm.Sub(mulai))
	if err != nil {
		return err
	}
	keAkhir, err := hariBulat(akhir.Sub(edm))
	if err != nil {
		return err
	}
	total, err := hariBulat(akhir.Sub(mulai))
	if err != nil {
		return err
	}
	// Guard bagi-nol: periode nol diperlakukan sebagai satu (hari).
	if total == 0 {
		total = 1
	}

	awal, err := bagiPorsi(keMulai, total)
	if err != nil {
		return fmt.Errorf("ProrateStartEDM: %w", err)
	}
	sisa, err := bagiPorsi(keAkhir, total)
	if err != nil {
		return fmt.Errorf("ProrateEDMEnd: %w", err)
	}
	o.ProrateStartEDM, o.ProrateEDMEnd = awal, sisa
	return nil
}

// hariBulat - selisih waktu dalam hari, hanya bila tepat bulat.
func hariBulat(d time.Duration) (int64, error) {
	if d%sehari != 0 {
		return 0, ErrSatuanSelisihWaktuBelumTerverifikasi
	}
	return int64(d / sehari), nil
}

// bagiPorsi - pembilang ÷ penyebut tepat pada 20 desimal, atau galat.
func bagiPorsi(pembilang, penyebut int64) (uang.Ratio, error) {
	ctx := utils.DecimalContext()
	hasil := new(apd.Decimal)
	kondisi, err := ctx.Quo(hasil, apd.New(pembilang, 0), apd.New(penyebut, 0))
	if err != nil {
		return uang.Ratio{}, err
	}
	// Tak eksak pada presisi 38 berarti pecahan berulang atau terlalu
	// panjang - tidak muat di 20 desimal tanpa pembulatan.
	if kondisi.Inexact() {
		return uang.Ratio{}, ErrModePembulatanBelumTerverifikasi
	}
	kondisi, err = ctx.Quantize(hasil, hasil, -desimalPorsiPeriode)
	if err != nil {
		return uang.Ratio{}, err
	}
	if kondisi.Inexact() {
		return uang.Ratio{}, ErrModePembulatanBelumTerverifikasi
	}
	return uang.Ratio{Value: hasil, Scale: desimalPorsiPeriode}, nil
}
