package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"
)

// Selesai Policy (tiket 69a-2, keputusan work owner K69a2-1..9 di docs/issues/69-jalur-policy.md): port
// Activity\SaveJsonPolicyFacIn_Act.xml jalur NB FIRE tanpa CALL / COMMIT, satu transaksi dengan Submit.

// TulisanPolicy - seluruh tulisan selesai Policy. Nomor polis, ProdDateTime, NOPOLIS dan SPREAD_DATE baris diisi repository
// di dalam transaksi (TANGGAL_CLOSING, KODE_PRODUKSI, GENERATE_SEQUENCE_NUMBER; JSON_POLIS case yang sudah ada dipakai ulang).
type TulisanPolicy struct {
	// BusinessOldID - QuotationData.BusinessOldId (GetOldIDBusiness_SQL), badan nomor polis.
	BusinessOldID string
	// Sekarang / MulaiPolis - dasar ProdDateTime (SaveJsonPolicy langkah 5); MulaiPolis nol = Begin date kosong.
	Sekarang, MulaiPolis time.Time
	// Username - JSON_POLIS.USERNAME (OperatorID.pyUserIdentifier). JSONData - DATA_JSON minimal (K69a2-1, ADR-0028).
	Username, JSONData string
	// Ceding - CEDING_FACINPRODUCTION (InsertCedingProduction langkah 2), satu per baris CedingCedantList.
	Ceding []BarisCedingProduksi
	// RisikoLoss - M_RISK_LOSS_PROFILE (SaveJsonPolicy langkah 19-21), sekali per polis (K69a2-8).
	RisikoLoss ProfilRisikoLoss
	// Produksi - FACINPRODUCTION (SaveFacinProdFireNB_Act langkah 13), ditulis bila IDPEGA belum ada (CekFacin, K69a2-7).
	Produksi []BarisFacInProduction
	// NomorPolis - KELUARAN: diisi repository di dalam transaksi (nomor baru atau NOPOLIS JSON_POLIS yang dipakai ulang);
	// dikembalikan Submit (work owner 09-10-2026, HasilSubmit.nomorPolis).
	NomorPolis string
}

// BarisCedingProduksi - satu baris CEDING_FACINPRODUCTION: CEDINGID / CEDINGNAME tanpa ";" dan PCT_SHARE (TempInput.CARI5).
type BarisCedingProduksi struct {
	CedingID, CedingName string
	Share                *apd.Decimal
}

// ProfilRisikoLoss - kolom M_RISK_LOSS_PROFILE (InsertRiskAndLossProfile_SQL); nil = NULL (variabel Local tak pernah diisi).
type ProfilRisikoLoss struct {
	TSINusantaraRe, PremiNusantaraRe, PremiumOR, PremiumQS, PremiumFacout *apd.Decimal
	COB                                                                   string
}

// BarisFacInProduction - satu baris FACINPRODUCTION (InsertTreatyProduction_Sql, 79 kolom; NOPOLIS / SPREAD_DATE dari
// repository). Tanggal = teks "DD/MM/YYYY HH24:MI:SS" jam dinding WIB ("" = NULL); angka *apd.Decimal (nil = NULL); teks "" =
// NULL. Kolom yang selalu sama di jalur FIRE (PRORATE 1, NOENDORS / SAILDATE / TYPE_EDM / DEDUCTION2 NULL, "NULL" teks) ada di
// repository. MarketingCode = QuotationData.MarketingCode; MARKETINGOFFICERCODE dibaca repository lewat GetMOID_Sql.
type BarisFacInProduction struct {
	IDPega, BusinessCode, CedingCo, SOB, BeginDate, EndDate                             string
	Province, Kabupaten, Kecamatan, Kelurahan, RoadName, BuildingNo, ZipCode, ObjectNo  string
	ObjectName, RiskLocation, ObjectItem, Coverage, CoverageNote, JnReas, CurrID        string
	AccumulationCode, OccupationCode, MarketingCode, InsuredName, RISlipCeding          string
	GroupType, QQName, ObjectNameID, LayerNo, PolicyType, CedingID, SOBID, InsuredID    string
	IdxObjectItem, StatusBusiness, GroupName, TopRisk                                   string
	Rate, PercentRNM, PercentShareSpreaded, TSI, Premi, PercentRIComm, PctRICommSelisih *apd.Decimal
	RIComm, LoL, PctBrokerageFee, BrokerageFee, TSI100                                  *apd.Decimal
	// Tiket 73 P5 (cabang Aneka / Golf SaveTreatyProduction_Act 19 / 21): NoSerial = `.SerialNo` objek ("" = teks "NULL" seperti
	// FIRE); MOIDKasus = QuotationData.MOID yang DIDAHULUKAN atas GetMOID (S11 L2492, K73-46; "" = GetMOID seperti FIRE).
	NoSerial, MOIDKasus string
}

// HariClosingBawaan - SaveJsonPolicy langkah 4 `@if(Local.TglProd=="",25,Local.TglProd)`.
const HariClosingBawaan = 25

// ProdDateTime - SaveJsonPolicy langkah 5 (ProdDateTime belum terisi):
//   - 5.1 sekarang;
//   - 5.2 Begin date lewat sehari penuh atau lebih dari sekarang -> Begin date (`@toInt(DateTimeDifference(Start, Current, D)) <
//     0`, deskripsi "jika start date lebih besar dari currentdate"; arah selisih `[dugaan kuat]` dari deskripsi itu, @toInt
//     memotong pecahan hari);
//   - 5.3 hari (Asia/Jakarta) > hari closing -> tanggal 01 bulan berikutnya 05:00 GMT (`yyyyMM` Jakarta + "01T050000.000 GMT",
//     @addCalendar +1 bulan).
func ProdDateTime(sekarang, mulai time.Time, hariClosing int, jakarta *time.Location) time.Time {
	prod := sekarang
	if !mulai.IsZero() && mulai.Sub(sekarang) >= 24*time.Hour {
		prod = mulai
	}
	lokal := prod.In(jakarta)
	if lokal.Day() > hariClosing {
		prod = time.Date(lokal.Year(), lokal.Month(), 1, 5, 0, 0, 0, time.UTC).AddDate(0, 1, 0)
	}
	return prod
}

// KelasNomorPolis - SaveJsonPolicy 9.3 `ParamSeq.CARI1` (kunci CLASS penghitung GENERATE_SEQUENCE_NUMBER).
const KelasNomorPolis = "ASM-FW-GISFW-Work-NB/RNW"

// JenisNomorPolis - `ParamSeq.CARI2 = ParamSeq.HASIL3 + "F"` (KODE_PRODUKSI NONLIFE + "F"; kunci JENIS penghitung).
func JenisNomorPolis(kodeProduksi string) string { return kodeProduksi + "F" }

// RakitNomorPolis - SaveJsonPolicy langkah 12: CARI2 + BusinessOldId + "." + MM.YYYY + "." + LPAD(urut, 5, '0'). Urut yang
// lebih dari lima digit TIDAK dipotong (`[keputusan agent]`, sama dengan RakitNomorKlaim: memotong menerbitkan nomor kembar).
func RakitNomorPolis(jenis, oldID, mmYYYY string, urut int) string {
	u := strconv.Itoa(urut)
	if len(u) < 5 {
		u = strings.Repeat("0", 5-len(u)) + u
	}
	return jenis + oldID + "." + mmYYYY + "." + u
}
