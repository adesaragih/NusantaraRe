package models

import "github.com/cockroachdb/apd/v3"

// Menu View Policy (Facultative) - tiket 71 (K71-1..3, A71-1..5). Semua dari tabel; nol DATA_JSON (K70-1).

// VersiPolis - satu versi polis ("View Data"): PRODKE + ID case (IDPEGA tanpa awalan kelas, GetPolicyData_ACT L1262);
// SistemBaru = case itu ada di T_WORK_POLIS (LINI FAC) -> detail = layar case (K71-1).
type VersiPolis struct {
	ProdKe, CaseID string
	SistemBaru     bool
}

// KepalaPolis - kepala ringkasan polis lama (A71-3), dari baris FACINPRODUCTION versi itu. Tanggal = teks DD-MM-YYYY; "" = NULL.
type KepalaPolis struct {
	CaseID, PolicyNo, NoEndors, InsuredName, QQName, SourceOfBusiness, CedingCo, StartPeriod, EndPeriod string
	PolicyType, ReffNumber, GroupName, StatusBusiness, NoMasterPolis                                    string
}

// BarisPolis - satu baris FACINPRODUCTION (A71-3). Rate = kolom VARCHAR2 apa adanya (dinormalkan services); angka nil = NULL.
type BarisPolis struct {
	ObjectNo, RiskLocation, ObjectItem, Coverage, AccumulationCode, Currency, Rate, Treaty string
	TSI100, TSIRNM, PremiRNM, PctShare, RIComm                                             *apd.Decimal
}

// CedingPolis - satu baris CEDING_FACINPRODUCTION.
type CedingPolis struct {
	CedingName string
	PctShare   *apd.Decimal
}

// RingkasanPolis - ringkasan satu versi polis lama.
type RingkasanPolis struct {
	Kepala KepalaPolis
	Baris  []BarisPolis
	Ceding []CedingPolis
}
