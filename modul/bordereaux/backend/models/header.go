package models

// Status dan posisi berkas - `AkseptasiBdx_DT` (kolom BORDEREAUX.STATUSAKSEP dan POSITION).
const (
	StatusAccept          = "Accept"
	StatusRejected        = "Rejected"
	StatusResolveComplete = "Resolve-Complete"
	PosisiChecker         = "Checker"
	PosisiSupervisor      = "Supervisor"
)

// Peran (workbasket) Bordereaux - migrasi 891. `ReasBordereauxAdmin` dibuang migrasi 893: Input Data kini hak menu
// PENUH (`M_LOGIN_GO_MENU.HAK`, keputusan work owner 04-10-2026).
const (
	PeranChecker    = "ReasBordereauxChecker"
	PeranSupervisor = "ReasBordereauxSupervisor"
)

// Header adalah satu baris BORDEREAUX.
type Header struct {
	BdxID string `json:"bdxId"`
	// Tanggal - waktu buat, `DD-MM-YYYY HH24:MI`.
	Tanggal      string `json:"tanggal"`
	UserInput    string `json:"userInput"`
	Type         string `json:"type"`
	TypeBusiness string `json:"typeBusiness"`
	MasterID     string `json:"masterId"`
	CedingID     string `json:"cedingId"`
	CedingName   string `json:"cedingName"`
	SobID        string `json:"sobId"`
	SobName      string `json:"sobName"`
	TreatyName   string `json:"treatyName"`
	// ReportStart dan ReportEnd - `DD-MM-YYYY`.
	ReportStart string `json:"reportStart"`
	ReportEnd   string `json:"reportEnd"`
	ReffNoSOA   string `json:"reffNoSoa"`
	ReffNoBDX   string `json:"reffNoBdx"`
	Position    string `json:"position"`
	StatusAksep string `json:"statusAksep"`
}

// Baris adalah satu baris detail: kolom tabel -> teks (angka bertitik desimal tanpa pemisah ribuan, tanggal
// `DD-MM-YYYY`, kosong = NULL). `ID` hanya pada baris yang dibaca dari tabel.
type Baris map[string]string

// KolomID - kunci ID baris detail di Baris.
const KolomID = "ID"

// Filter daftar - parameter `InboxBordereaux_RD` (CARI1..CARI11). Kosong = tidak disaring.
type Filter struct {
	BdxID, Type, Business, ReffNoSOA, ReffNoBDX, Ceding, Treaty string
	// ReportStart/ReportEnd `DD-MM-YYYY`: rentang (Pega mencari tanggal PERSIS sama - diperbaiki).
	ReportStart, ReportEnd string
	Position, Status       string
}

// Ringkasan - satu baris tab Summary (`CountSummaryBdx`): total per mata uang.
type Ringkasan struct {
	Currency  string `json:"currency"`
	Reinsurer string `json:"reinsurer"`
	RNM       string `json:"rnm"`
}

// Riwayat - satu baris BORDEREAUX_HISTORY.
type Riwayat struct {
	// Tanggal - `DD-MM-YYYY HH24:MI`.
	Tanggal    string `json:"tanggal"`
	PIC        string `json:"pic"`
	IsApproved bool   `json:"isApproved"`
	Komentar   string `json:"komentar"`
}

// Cedant - satu baris AGENT untuk popup Choose Master Treaty (`BrowseAgentNonLife_RD`).
type Cedant struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}

// MasterTreaty - satu baris TREATY_IN (`BrowseTREATY_IN`).
type MasterTreaty struct {
	ID           string `json:"id"`
	ContractName string `json:"contractName"`
	ReinsType    string `json:"reinsType"`
	SobID        string `json:"sobId"`
	SobName      string `json:"sobName"`
	CedingID     string `json:"cedingId"`
	CedingName   string `json:"cedingName"`
}

// IrisanChart - jumlah berkas satu Business x Type x Ceding (chart daftar: kelompok Business, dibuka per Ceding).
type IrisanChart struct {
	Business   string `json:"business"`
	Type       string `json:"type"`
	CedingID   string `json:"cedingId"`
	CedingName string `json:"cedingName"`
	Jumlah     int    `json:"jumlah"`
}
