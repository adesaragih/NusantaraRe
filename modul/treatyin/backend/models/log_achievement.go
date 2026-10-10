package models

// BarisLogAchievement - satu baris `AchievementLists` yang tombol `Submit`
// sub-tab Achievement catat (`InsertToLogAchievement` langkah [2.1]), ejaan
// properti Pega apa adanya. Angka sebagai TEKS — services mengurainya.
type BarisLogAchievement struct {
	Quarter          string `json:"Quarter"`
	QUARTERYEAR      string `json:"QUARTERYEAR"`
	CurrencyID       string `json:"CurrencyID"`
	Currency         string `json:"Currency"`
	PREMIUM          string `json:"PREMIUM"`
	RICOMM           string `json:"RICOMM"`
	BROKERAGE        string `json:"BROKERAGE"`
	NETPREMIUM       string `json:"NETPREMIUM"`
	PaidClaim        string `json:"PaidClaim"`
	CASHCALL         string `json:"CASHCALL"`
	OutstandingClaim string `json:"OutstandingClaim"`
	IncuredClaim     string `json:"IncuredClaim"`
	Total            string `json:"Total"`
	LossRatio        string `json:"LossRatio"`
}
