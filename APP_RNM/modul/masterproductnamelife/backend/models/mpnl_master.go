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
