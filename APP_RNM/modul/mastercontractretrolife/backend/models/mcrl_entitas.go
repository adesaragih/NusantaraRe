// Package models memuat bentuk data modul Master Contract Retro Life
// (`mastercontractretrolife`).
//
// Untuk apa berkas ini: lima entitas tabel WARISAN `POOLDATA.*_LIFE` dan
// master rujukan yang dibaca layar, persis kolom DDL `[data DBA]`
// (`docs/ddl-tables-from-dba.md`). Nol tabel baru (K1,
// `docs/RALAT-DEV-30-09-2026.md`).
//
// ⛔ Uang dan persen tidak pernah float (ADR-U-0003): `*apd.Decimal`, nil =
// kolom kosong (ADR-U-0027), dan JSON membawanya sebagai TEKS desimal.
// ⛔ `RIRATE` TEKS apa adanya - ia nama tabel rate (R7), bukan angka.
//
// Dibaca sesudah: docs/PARITAS-LAYAR-DAN-AKSI.md.
package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// TahunTreaty adalah satu baris `TREATYYEAR_LIFE` (grid `BrowseTreatyYear_Life_RD`).
type TahunTreaty struct {
	ID               string
	TreatyYear       string // label layar: TRANSACTION YEAR
	UnderwritingYear string
	StartDate        time.Time
	EndDate          time.Time
	UserID           string
	TglUpdate        time.Time
}

// Kontrak adalah satu baris `TREATYCONTRACT_LIFE` - layer proteksi satu jenis
// reasuransi di dalam satu tahun treaty.
type Kontrak struct {
	ID              string
	IDTreatyYear    string
	ReinsTypeID     string
	ReinsTypeName   string // salinan `REINSURANCETYPE.NOTE` (TreatyLimit_TypeProtect 4.1)
	TreatyStartDate time.Time
	TreatyEndDate   time.Time
	UserID          string
	TglUpdate       time.Time
	IDR             *apd.Decimal // MAXIMUM LIMIT (IDR) - batas atas
	USD             *apd.Decimal // MAXIMUM LIMIT (USD) - boleh kosong
	BIDR            *apd.Decimal // MINIMUM LIMIT (IDR) - batas bawah
	BUSD            *apd.Decimal // MINIMUM LIMIT (USD)
	IDRSelisih      *apd.Decimal // dihitung Go (K5)
	USDSelisih      *apd.Decimal
}

// Reinsurer adalah satu baris `TREATYREINSURER_LIFE`.
type Reinsurer struct {
	ID               string
	TreatyYearID     string
	TreatyContractID string
	ReinsTypeID      string // salinan kontrak (K4)
	ReinsTypeName    string
	ReinsurerID      string // `AGENT.ID`
	ReinsurerName    string // `AGENT.CLIENTNAME`
	PctShare         *apd.Decimal
	Komisi           *apd.Decimal // kolom `COMMISION` (ejaan warisan), label (%) DISCOUNT
	OvrComm          *apd.Decimal
	UserID           string
	TglUpdate        time.Time
}

// SecurityReinsurer adalah satu baris `TREATYSECURITYREINSURER_LIFE` - retrosesi
// atas bagian seorang reinsurer. `PctShare` = persen DARI share induknya.
type SecurityReinsurer struct {
	ID                string
	TreatyYearID      string
	TreatyContractID  string
	TreatyReinsurerID string
	ReinsurerID       string
	ReinsurerName     string
	PctShare          *apd.Decimal
	UserID            string
	TglUpdate         time.Time
}

// Business adalah satu baris `TREATYBUSINESS_LIFE`.
type Business struct {
	ID               string
	TreatyYearID     string
	TreatyYear       string // salinan tahun (K4)
	TreatyContractID string
	ReinsTypeID      string // salinan kontrak (K4)
	ReinsTypeName    string
	BizCode          string
	BizName          string
	RIRateID         string // ID tabel rate
	RIRate           string // nama tabel rate - TEKS apa adanya (R7)
	UserID           string
	TglUpdate        time.Time
}

// Dampak - baris ANAK yang ikut terhapus bersama satu induk (K2, tiket 09);
// baris induknya sendiri tidak dihitung.
type Dampak struct {
	Security  int64 `json:"security"`
	Reinsurer int64 `json:"reinsurer"`
	Business  int64 `json:"business"`
}

// TotalShareKontrak - total share satu kontrak untuk laporan tiket 11.
type TotalShareKontrak struct {
	KontrakID     string
	TahunID       string
	TreatyYear    string
	ReinsTypeName string
	Total         *apd.Decimal // nil = tanpa reinsurer
}

// JenisReasuransi - satu pilihan dropdown `REINS TYPE` (`BrowseReinsuranceTypeLimit_RD`).
type JenisReasuransi struct {
	ID   string `json:"id"`
	Note string `json:"note"`
}

// MasterReinsurer - satu pilihan autocomplete `REINSURER NAME` (`BrowseCedingCoLife_RD`).
type MasterReinsurer struct {
	ID         string `json:"id"`
	ClientName string `json:"clientName"`
	// StatusActive - dibaca saat simpan untuk saringan `.StatusActive = 1`; tidak dikirim ke layar.
	StatusActive string `json:"-"`
}

// Life - saringan `BrowseCedingCoLife_RD` b565 `.ID Contains "L0"`.
func (m MasterReinsurer) Life() bool { return strings.Contains(m.ID, PenandaReinsurerLife) }

// Aktif - saringan `BrowseCedingCoLife_RD` b601 `.StatusActive = 1`.
func (m MasterReinsurer) Aktif() bool { return m.StatusActive == StatusMasterReinsurerAktif }

// MasterBusiness - satu pilihan autocomplete `BUSINESS NAME` (`BrowseBusinessLife_RD`).
type MasterBusiness struct {
	ID    string `json:"id"`
	Note  string `json:"note"`
	OldID string `json:"oldId"`
}

// Life - saringan `BrowseBusinessLife_RD` b651 `.OLDID StartsWith "L"`.
func (m MasterBusiness) Life() bool { return strings.HasPrefix(m.OldID, AwalanBusinessLife) }

// Nilai saringan master VERBATIM RD - di models karena bernama kode/status
// (penjaga Claim Life `TestKodeStatusLiteralHanyaDiModels`).
const (
	// FlagJenisReasuransiLife - `BrowseReinsuranceTypeLimit_RD` b578 `.Flag = 1` ("1 for life").
	FlagJenisReasuransiLife = "1"
	// StatusMasterReinsurerAktif - `BrowseCedingCoLife_RD` b601 `.StatusActive = 1`.
	StatusMasterReinsurerAktif = "1"
	// PenandaReinsurerLife - `BrowseCedingCoLife_RD` b565 `.ID Contains "L0"`.
	PenandaReinsurerLife = "L0"
	// AwalanBusinessLife - `BrowseBusinessLife_RD` b651 `.OLDID StartsWith "L"`.
	AwalanBusinessLife = "L"
)

// waktuJSON menulis stempel waktu; nol menjadi teks kosong (ADR-U-0027).
func waktuJSON(t time.Time) string { return utils.FormatTanggalWaktu(t) }

// MarshalJSON - nama medan JSON SAMA dengan `frontend/api.ts`.
func (t TahunTreaty) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID               string `json:"id"`
		TreatyYear       string `json:"treatyYear"`
		UnderwritingYear string `json:"underwritingYear"`
		StartDate        string `json:"startDate"`
		EndDate          string `json:"endDate"`
		UserID           string `json:"userId"`
		TglUpdate        string `json:"tglUpdate"`
	}{t.ID, t.TreatyYear, t.UnderwritingYear, utils.FormatTanggal(t.StartDate), utils.FormatTanggal(t.EndDate),
		t.UserID, waktuJSON(t.TglUpdate)})
}

// MarshalJSON - uang sebagai teks desimal, nil = "".
func (k Kontrak) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID              string `json:"id"`
		IDTreatyYear    string `json:"idTreatyYear"`
		ReinsTypeID     string `json:"reinsTypeId"`
		ReinsTypeName   string `json:"reinsTypeName"`
		TreatyStartDate string `json:"treatyStartDate"`
		TreatyEndDate   string `json:"treatyEndDate"`
		UserID          string `json:"userId"`
		TglUpdate       string `json:"tglUpdate"`
		IDR             string `json:"idr"`
		USD             string `json:"usd"`
		BIDR            string `json:"bIdr"`
		BUSD            string `json:"bUsd"`
		IDRSelisih      string `json:"idrSelisih"`
		USDSelisih      string `json:"usdSelisih"`
	}{k.ID, k.IDTreatyYear, k.ReinsTypeID, k.ReinsTypeName, utils.FormatTanggal(k.TreatyStartDate),
		utils.FormatTanggal(k.TreatyEndDate), k.UserID, waktuJSON(k.TglUpdate),
		utils.FormatDecimal(k.IDR), utils.FormatDecimal(k.USD), utils.FormatDecimal(k.BIDR), utils.FormatDecimal(k.BUSD),
		utils.FormatDecimal(k.IDRSelisih), utils.FormatDecimal(k.USDSelisih)})
}

// MarshalJSON - `komisi` = kolom `COMMISION`.
func (r Reinsurer) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID               string `json:"id"`
		TreatyYearID     string `json:"treatyYearId"`
		TreatyContractID string `json:"treatyContractId"`
		ReinsTypeID      string `json:"reinsTypeId"`
		ReinsTypeName    string `json:"reinsTypeName"`
		ReinsurerID      string `json:"reinsurerId"`
		ReinsurerName    string `json:"reinsurerName"`
		PctShare         string `json:"pctShare"`
		Komisi           string `json:"komisi"`
		OvrComm          string `json:"ovrComm"`
		UserID           string `json:"userId"`
		TglUpdate        string `json:"tglUpdate"`
	}{r.ID, r.TreatyYearID, r.TreatyContractID, r.ReinsTypeID, r.ReinsTypeName, r.ReinsurerID, r.ReinsurerName,
		utils.FormatDecimal(r.PctShare), utils.FormatDecimal(r.Komisi), utils.FormatDecimal(r.OvrComm),
		r.UserID, waktuJSON(r.TglUpdate)})
}

// MarshalJSON - share mentah; eksposur dihitung di jawaban daftar (services).
func (s SecurityReinsurer) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID                string `json:"id"`
		TreatyYearID      string `json:"treatyYearId"`
		TreatyContractID  string `json:"treatyContractId"`
		TreatyReinsurerID string `json:"treatyReinsurerId"`
		ReinsurerID       string `json:"reinsurerId"`
		ReinsurerName     string `json:"reinsurerName"`
		PctShare          string `json:"pctShare"`
		UserID            string `json:"userId"`
		TglUpdate         string `json:"tglUpdate"`
	}{s.ID, s.TreatyYearID, s.TreatyContractID, s.TreatyReinsurerID, s.ReinsurerID, s.ReinsurerName,
		utils.FormatDecimal(s.PctShare), s.UserID, waktuJSON(s.TglUpdate)})
}

// MarshalJSON - `riRate` teks apa adanya.
func (b Business) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID               string `json:"id"`
		TreatyYearID     string `json:"treatyYearId"`
		TreatyYear       string `json:"treatyYear"`
		TreatyContractID string `json:"treatyContractId"`
		ReinsTypeID      string `json:"reinsTypeId"`
		ReinsTypeName    string `json:"reinsTypeName"`
		BizCode          string `json:"bizCode"`
		BizName          string `json:"bizName"`
		RIRateID         string `json:"riRateId"`
		RIRate           string `json:"riRate"`
		UserID           string `json:"userId"`
		TglUpdate        string `json:"tglUpdate"`
	}{b.ID, b.TreatyYearID, b.TreatyYear, b.TreatyContractID, b.ReinsTypeID, b.ReinsTypeName, b.BizCode, b.BizName,
		b.RIRateID, b.RIRate, b.UserID, waktuJSON(b.TglUpdate)})
}
