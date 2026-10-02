package models

// Tujuh pemilih master (tiket 04, PARITAS §4). Setiap section pemilih
// menampilkan dua kolom - `ID` dan `Name` (`RIRate Name` untuk R/I Rate) - dan
// tombol `Choose` menyalin keduanya ke produk lewat `set*_DT`.

// JenisMaster - segmen jalur `GET …/master/{jenis}`.
type JenisMaster string

// Ketujuh master. Ceding dan SOB berbagi RD `BrowseCedingCoLife_RD` (kelas
// `AGENT`); R/I Risk dan R/I Rate DUA master berbeda (`[fakta bisnis — work owner]`).
const (
	MasterCeding        JenisMaster = "ceding"         // `ChooseCeding` → `setCeding_DT`
	MasterSOB           JenisMaster = "sob"            // `ChooseSOB` → `setSOB_DT`
	MasterPemegangPolis JenisMaster = "pemegang-polis" // `ChoosePolicyHolder` → `setPolicyHolder_DT`
	MasterMataUang      JenisMaster = "mata-uang"      // `ChooseCurrency` → `setCurrency_DT`
	MasterRIRisk        JenisMaster = "ri-risk"        // `ChooseRIRisk` → `setRIRISK_DT`
	MasterRIRate        JenisMaster = "ri-rate"        // `ChooseRIRate` → `SetRIRate` (baris plan)
	MasterPenyebab      JenisMaster = "penyebab"       // `ChooseCauseOfLoss` → `setCauseOfLoss_DT`
)

// SemuaJenisMaster - urutan tombol `Choose*` di layar.
var SemuaJenisMaster = []JenisMaster{MasterCeding, MasterSOB, MasterRIRisk, MasterPenyebab, MasterPemegangPolis,
	MasterMataUang, MasterRIRate}

// NilaiMaster - satu baris grid pemilih: kolom `ID` dan `Name`.
type NilaiMaster struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}

// JenisPlan - satu baris autocomplete `Plan Name` (`.Plan` b33121, RD
// `BrowseProductTypeLife_RD` kelas `PRODUCT_TYPE_LIFE`): `.CoverName` → `.Plan`,
// `.ID` → `.PlanID`, `.Business` → `.Name`, `.Benefit` → `.Benefit`.
type JenisPlan struct {
	ID        string `json:"id"`
	CoverName string `json:"coverName"`
	Business  string `json:"business"`
	Benefit   string `json:"benefit"`
}

// BarisRate - satu baris dialog `View Rate` (section `ViewRate`, RD `BrowseRateLife_RD`): enam kolom
// grid (`ViewRate.xml` b2035–b2809). ⛔ TEKS apa adanya - `RATE` di view berdesimal koma maupun titik.
type BarisRate struct {
	ID       string `json:"id"`
	UsedBy   string `json:"usedBy"`
	Gender   string `json:"gender"`
	Contract string `json:"contract"`
	Age      string `json:"age"`
	Rate     string `json:"rate"`
}
