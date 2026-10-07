package models

import "github.com/cockroachdb/apd/v3"

// Tab Spreading kasus FIRE (tiket 48). Uang / persen = *apd.Decimal (ADR-0034), nil = kosong.

// BarisSpreading - satu `.CoverageList(k).SpreadingList(s)` (kelas Data-SpreadingRisk; T_SPREADINGLIST).
type BarisSpreading struct {
	TreatyType, TreatyName                                                           string
	SharePercentage, TSIGrossSpreaded, TSISpreaded, ClaimEstimation, PremiumSpreaded *apd.Decimal
}

// CoverageSpreading - satu coverage dengan medan yang dipakai spreading. ID = T_COVERAGELIST.ID (internal, tidak dikirim).
type CoverageSpreading struct {
	ID, OldID, CoverageBasis, CoverageNote                                  string
	Rate, Premium, Discount, TSI, TSILiability, LimitOfLiability, FirstLoss *apd.Decimal
	TSINusantaraRe, PremiNusantaraRe                                        *apd.Decimal
	Spreading                                                               []BarisSpreading
}

// ItemSpreading - satu `.Property.PropertyItemList(m)`.
type ItemSpreading struct {
	ItemType, Currency      string
	TSI, TotalGrossPremi    *apd.Decimal
	TotalPremiumNusantaraRe *apd.Decimal // .TotalPremiumNusantaraRe = Σ PremiNusantaraRe coverage
	Coverages               []CoverageSpreading
}

// TotalSpreading - satu baris total per lokasi (`Property.TotalTSIPremiSpreadRNM`) atau ringkasan per mata uang × treaty
// (`TotalSpreadingCurrency`). Currency = mata uang objek (Pega menyimpannya di TreatyName); TreatyName = nama treaty.
type TotalSpreading struct {
	Currency, TreatyType, TreatyName                                                              string
	SharePercentage, ClaimSpreaded, TSISpreaded, ClaimEstimation, ClaimAmountIDR, PremiumSpreaded *apd.Decimal
}

// LokasiSpreading - satu `.LocationList(n)` dengan total per lokasi.
type LokasiSpreading struct {
	ObjectNo, ObjectName, Location string
	IsTopRisk                      bool
	Items                          []ItemSpreading
	Total                          []TotalSpreading
}

// RingkasanMataUangSpreading - satu baris `TotalSpreadAll` (per mata uang).
type RingkasanMataUangSpreading struct {
	Currency                      string
	TSITopRisk, TSI, LoL, Premium *apd.Decimal
}

// KasusSpreading - data tab Spreading satu case: % Share RNM, Begin date (teks Pega), lokasi.
type KasusSpreading struct {
	PercentShare  *apd.Decimal
	StartDateTime string
	Lokasi        []LokasiSpreading
}

// TreatySpreading - satu baris GetTreatyName (PROPORTIONALARRG TREATY LIMIT): ID = REINSTYPEID (CARI1), Name =
// REINSTYPENAME (CARI2), Limit = RP (CARI3, teks apa adanya).
type TreatySpreading struct {
	ID, Name, Limit string
}

// KapasitasTreaty - baris pertama KAPASITAS_TREATY yang memuat TSI (GetKapasitasTreaty Obj-Browse).
type KapasitasTreaty struct {
	TreatyNameQS, IDTreatyQS, TreatyNameSPL, IDTreatySPL string
	MaxLimitQSIDR, MaxLimitSPLIDR                        *apd.Decimal
}

// TemplateSpreading - satu baris template Copy Spreading (SpreadingList.pxResults).
type TemplateSpreading struct {
	TreatyType, TreatyName string
	SharePercentage        *apd.Decimal
}
