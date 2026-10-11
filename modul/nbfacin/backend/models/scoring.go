package models

import "github.com/cockroachdb/apd/v3"

// Tab Scoring Risk kasus FIRE (tiket 54) - halaman `.OfferFacIn.ScoringRisk` (kelas Data-ScoringRisk). Skor / persen =
// *apd.Decimal (ADR-0034), nil = kosong.

// FaktorScoring - satu sub-faktor: Kunci = kunci kontrak (`pilihan` / `skor`), Grup / Nama = nama halaman Pega di bawah
// `.DataScoringRiskList(n)` (Grup kosong = halaman faktor langsung) - juga FACTOR_GROUP / FACTOR_NAME T_SCORING_FACTOR
// (aturan loader V-30), JumlahOpsi = banyak ChechBoxN.
type FaktorScoring struct {
	Kunci, Grup, Nama string
	JumlahOpsi        int
}

// Kunci sub-faktor yang diisi otomatis / dibaca aturan selain penjumlahan (tiket 54).
const (
	KunciOccupation      = "occupation"
	KunciEarthquake      = "earthquake"
	KunciLossRatio       = "lossRatio"
	KunciTotalSumInsured = "totalSumInsured"
)

// DaftarFaktorScoring - 24 sub-faktor, urut layar (kontrak); `.Score` Pega (`Activity\SetScore_Act.xml` langkah 5)
// menjumlah faktor yang sama dengan Flood sebelum Riot - jumlahnya sama `[terverifikasi]`. Jumlah opsi = ChechBox1..n
// `ValueScoringRisk_act` / `SetScore_Act`.
var DaftarFaktorScoring = []FaktorScoring{
	{"occupation", "", "Occupation", 5},
	{"riot", "RiskLocation", "Riot", 4},
	{"flood", "RiskLocation", "Flood", 4},
	{"earthquake", "RiskLocation", "Earthquake", 7},
	{"tsunami", "RiskLocation", "Tsunami", 4},
	{"classConstruction", "ObjectConditions", "ClassConstruction", 6},
	{"building", "ObjectConditions", "Building", 5},
	{"machinery", "ObjectConditions", "Machinery", 5},
	{"stock", "ObjectConditions", "Stockinthepolicycover", 4},
	{"typeOfStock", "ObjectConditions", "TypeofStock", 4},
	{"portableFireExtinguisher", "FEA", "PortableFireExtinguisher", 4},
	{"hydrant", "FEA", "HydrantSystem", 3},
	{"sprinkler", "FEA", "SprinklerSystem", 3},
	{"fireAlarm", "FEA", "FireAlarmSystem", 2},
	{"securityGuard", "FEA", "SecurityGuard", 3},
	{"fireBrigade", "FEA", "Firebrigade", 5},
	{"lossRatio", "LossRatio", "LossRatio", 10},
	{"riskImprovement", "", "RiskImprovement", 3},
	{"surveyReport", "", "SurveyReport", 6},
	{"rate", "Others", "Rate", 3},
	{"deductible", "Others", "Deductible", 3},
	{"totalSumInsured", "Others", "TotalSumInsured", 2},
	{"otherInformation", "Others", "OtherInformation", 3},
	{"clauses", "Others", "Clauses", 3},
}

// KolomScoring - satu `.DataScoringRiskList(n)`: kolom 1 Top Risk, kolom 2 Dominant Risk. Masukan: IndexLocation
// (ObjectNo lokasi dominan; kolom 1 tidak memakainya), Pilihan (kunci faktor -> nomor ChechBox 1..n, 0 / absen = tidak
// dicentang), PersenLossRatio (`.LossRatio.LossRatio.PercentLossRatio`). Hasil hitung: OccupationNote / OccupationCode /
// TopRiskLocation (dari lokasinya), Skor (ScorePerFactor), Total (`.Score`). Ada = kolom dihitung dan disimpan (kolom 2
// hanya bila Dominant dipilih).
type KolomScoring struct {
	IndexLocation   string
	Pilihan         map[string]int
	PersenLossRatio *apd.Decimal

	Ada                                             bool
	OccupationNote, OccupationCode, TopRiskLocation string
	Skor                                            map[string]*apd.Decimal
	Total                                           *apd.Decimal
}

// HasilScoring - `.ScoringResult` (`Activity\ScoringResult.xml`): kelas skor, "Due to", "Approval From", rekomendasi.
type HasilScoring struct {
	Good, AboveAverage, Average, BelowAverage, Poor              bool
	AcceptableNotApproval, AcceptableWithApproval, NotAcceptable bool
	Occupation, ObjectConditions, LossRatio, Others              bool
	OperationalDivHead, OperationalDirector, TechnicalDirector   bool
	PresidentDirector                                            bool
}

// Scoring - halaman ScoringRisk hasil hitung (yang disimpan).
type Scoring struct {
	Status         string
	Kolom          [2]KolomScoring
	FinalScore     *apd.Decimal
	NoteFinalScore string
	Hasil          HasilScoring
}

// RingkasanScoring - `.ScoringRisk.FinalScore` / `.NoteFinalScore` tersimpan (`Section\Periode.xml` sel "Risk Scoring",
// baca-saja). FinalScore nil = belum ada.
type RingkasanScoring struct {
	FinalScore     *apd.Decimal
	NoteFinalScore string
}

// KasusScoring - yang dibaca repository: posisi kasus dan FLAG_ONGOING_POLICY (T_WORK_POLIS) serta masukan kolom
// tersimpan (Pilihan / PersenLossRatio / IndexLocation; Ada = baris kolom itu tersimpan).
type KasusScoring struct {
	Position, FlagOnGoingPolicy string
	Kolom                       [2]KolomScoring
}
