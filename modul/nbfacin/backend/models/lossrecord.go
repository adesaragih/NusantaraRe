package models

import "github.com/cockroachdb/apd/v3"

// Tab utama Loss Record kasus FIRE (tiket 55) - daftar kerugian TINGKAT CASE `.OfferFacIn.ListCauseOfLoss` (kelas
// Data-CauseOfLoss), Premi PAR / MB / PL / Other, dan loss ratio 1 / 3 / 5 tahun. BEDA dari tiket 42 (`.Property
// .ListCauseOfLoss` per lokasi, T_LISTCAUSEOFLOSS). Uang / kurs / persen = *apd.Decimal (ADR-0034), nil = kosong.

// BarisLossRecord - satu `.OfferFacIn.ListCauseOfLoss(n)`. DateOfLoss: DD-MM-YYYY di batas services <-> handler; teks
// Pega DateTime di batas services <-> repository (pola tiket 42). UwYear / TahunKe / LossAmountIDR / KursIDR dihitung
// services (`SaveLossRecord`); TahunKe 0 = di luar kelima jendela.
type BarisLossRecord struct {
	UwYear                                   string
	TahunKe                                  int
	ClassOfBusiness, DateOfLoss, CauseOfLoss string
	Currency                                 string
	LossAmount, LossAmountIDR, KursIDR       *apd.Decimal
}

// TotalLossRecord - satu baris "Total Loss" per Class Of Business + Currency (`TotalLossCOB`, SaveLossRecord langkah 8).
type TotalLossRecord struct {
	ClassOfBusiness, Currency string
	LossAmount, LossAmountIDR *apd.Decimal
}

// RasioLossRecord - `.OfferFacIn.LossRatio{1,3,5}YearPercent` (Gross*) / `.NetLossRatio{1,3,5}Year` (Net*), persen.
type RasioLossRecord struct {
	Gross1, Gross3, Gross5, Net1, Net3, Net5 *apd.Decimal
}

// LossRecord - isi tab Loss Record satu kasus (yang disimpan / ditampilkan).
type LossRecord struct {
	PremiPar, PremiMb, PremiPl, PremiOther *apd.Decimal
	Baris                                  []BarisLossRecord
	Rasio                                  RasioLossRecord
}
