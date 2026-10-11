package models

// Tab Object + Coverage MARINE CARGO (tiket 74 P1 + P2). `.OfferFacIn.CargoList` -> T_CARGOLIST (migrasi 256) dengan halaman anak
// satu-per-cargo T_SHIP / T_LC / T_SURVEYAGENT, grid T_ADDITIONALSHIP (induk T_SHIP), coverage T_COVERAGELIST ber-PARENT_TABLE
// 'T_CARGOLIST' + T_DEDUCTIBLELIST + T_ADDITIONALCOVERAGE (Warranty Extension) + T_SPREADINGLIST (internal, P3). Teks apa adanya;
// angka *apd.Decimal; tanggal di model = teks Pega (layanan mengubah ke / dari kabel DD-MM-YYYY).

import "github.com/cockroachdb/apd/v3"

// ObjekCargo - satu baris CargoList. ID = T_CARGOLIST.ID (teks; kosong = baris baru). Medan "server" diisi layanan saat PUT.
type ObjekCargo struct {
	ID                           string
	IsTopRisk                    bool
	TradingID, TradingNote       string // TradingNote server (ProtectObjectMarine)
	GoodID, GoodNote             string // GoodNote server
	ConveyanceID, ConveyanceName string // ConveyanceName server
	ConveyanceNote               string // server
	PackingID, PackingNote       string // PackingNote server
	FromRute, ToRute             string
	SailDate, BLNumber           string
	InvoiceNumber, InvoiceDate   string
	MaxAge, MaxAgeDetail         string // MaxAgeDetail server (`.Warranty` master MaxAge)
	FlagOldData, FlagDelete      string // tidak diubah layar NB; dibawa apa adanya
	SurveyAgent                  AgenSurvei
	LC                           LCCargo
	Ship                         KapalCargo
	// Coverages - nil = tidak dikirim (PUT: coverage lama cargo ber-ID sama dipertahankan); larik = diganti utuh.
	Coverages []CoverageCargo
}

// AgenSurvei - `.PolicyData.SurveyAgent` (K74-6 per cargo).
type AgenSurvei struct {
	AgentID, FullName, Address, City, Country string
}

// Kosong - halaman tanpa isi (tidak disisipkan).
func (a AgenSurvei) Kosong() bool { return a == AgenSurvei{} }

// LCCargo - `.PolicyData.LC`; LC = cekbox "LC Status".
type LCCargo struct {
	LC                                       bool
	Bank, Number, Condition, Remark, EndDate string
}

// KapalCargo - `.PolicyData.Ship` (popup Select Ship, SelectShip_act L470-873); Age server (CountAgeShip). ID = `.Ship.ID`
// (POOLDATA.SHIP.ID).
type KapalCargo struct {
	ID, NmShip, ShipTypeID, ShipTypeName string
	GRT, NRT, DWT, Age                   string
	FlagName, Construction               string
	YMake1, YMake2, Remark, Register     string
	Imo, FormerShipName, Classification  string
	Rebuilt                              string
	AdditionalShip                       []KapalTambahan
}

// KapalTambahan - satu baris grid `.PolicyData.Ship.AdditionalShip` (autocomplete BrowseShip_RD, MarineCargoDtl L16928-17001).
type KapalTambahan struct {
	ID, Register, NmShip, VoyageNo, ShipTypeID, FlagName string
	NRT, DWT, GRT, YMake1, Construction                  string
}

// CoverageCargo - satu coverage cargo (form `InputCoverageCargo_FacIn`). Currency = nama (CURRENCY_CODE); CurrencyID / OldID,
// Discount, TSILiability, TSINusantaraRe, PremiNusantaraRe, ProRatePercent = server.
type CoverageCargo struct {
	Coverage, CoverageNote, CoverageInitial string
	Currency, CurrencyID, CurrencyOldID     string
	Rate, LimitOfLiability                  *apd.Decimal
	DiscountPercentage, TSI, Premium        *apd.Decimal
	MinPremium, Discount, TSILiability      *apd.Decimal
	TSINusantaraRe, PremiNusantaraRe        *apd.Decimal
	ProRatePercent                          *apd.Decimal
	WarrantyExtension                       []WarrantyExtension
	Deductibles                             []DeductibleCargo
	// Spreading - baris T_SPREADINGLIST coverage ini (tiket 74 P3); internal - dibaca dan ditulis ulang apa adanya.
	Spreading []BarisSpreading
}

// WarrantyExtension - satu baris grid `.AdditionalCoverage` (BrowseMarineCondition_RD: Coverage = .ID, CoverageNote =
// .ConditionName).
type WarrantyExtension struct {
	Coverage, CoverageNote string
}

// DeductibleCargo - deductible FIRE / Aneka + `.Descriptions` (InputCoverageCargo_FacIn L17652; kolom DESCRIPTIONS migrasi 256).
type DeductibleCargo struct {
	Deductible
	Descriptions string
}

// TotalGrossCargo - satu baris `TotalTSIPremiGrossList` (SumTotalTSIPremiGross_Act langkah 4): per mata uang.
type TotalGrossCargo struct {
	Currency           string
	TSI, Premium, Rate *apd.Decimal
}

// BarisMasterCargo - satu pilihan master Trading / Goods / Conveyance / Packing (view POOLDATA, tiket 74): .ID / .Description.
type BarisMasterCargo struct {
	ID, Description string
}

// BarisAgenSurvei - satu agen survei (view POOLDATA.LLOYDAGENT, popup Select Agent).
type BarisAgenSurvei struct {
	ID, Name, Address, City, Country, Contact string
}

// BarisMaxAge - satu Max Age (view POOLDATA.MAXAGE): ID -> .MaxAge, Warranty -> .MaxAgeDetail.
type BarisMaxAge struct {
	ID, Warranty string
}

// BarisKondisiMarine - satu kondisi marine (view POOLDATA.CONDITION TYPE 'MARINE'): Select Coverage / Warranty Extension.
type BarisKondisiMarine struct {
	ID, ConditionName string
}
