package models

import "time"

// Opportunity - isian form Opportunity yang sudah diperiksa services (tiket 29): satu
// baris T_NB_OPPORTUNITY (butir 76.3). Teks apa adanya; tanggal tanpa jam.
type Opportunity struct {
	EstimatedClosingDate time.Time
	BusinessProspectName string
	AccountID            string
	InsuredID            string
	GroupBusinessID      string
	GroupBusiness        string
	ClassOfBusiness      string
	TypeOfInward         string
	TypeOfFacultative    string
	Phase                string
	Stage                string
	OpportunitySource    string
	BusinessStatus       string
	Description          string
}
