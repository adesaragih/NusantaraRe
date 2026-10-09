package models

import "github.com/cockroachdb/apd/v3"

// Baris DETAIL kontrak tuntas — `TREATYINDETAIL` (Treaty In) dan
// `TREATYINDETAILEDM` (Treaty In Adjustment).
//
// ---------------------------------------------------------------------
// ⭐ PERINTAH WO 9 Oktober 2026
// ---------------------------------------------------------------------
//
//	"implementasikan untuk simpan ke TREATYINDETAIL UNTUK MASTER TREATY IN
//	 DAN SIMPAN KE TREATYINDETAILEDM UNTUK MASTER TREATY IN ADJUSTMENT …
//	 inisert ke treatyindetail dan treatyindetailedm jika sudah resolve
//	 complete!"
//
// Mencabut keputusan 7 Oktober 2026 ("biarkan data nya ditarik dari table
// nya masing masing saja") untuk kedua tabel ini. Penulisnya
// `SaveTreatyInDetail_Act` / `SaveTreatyInDetailEdm_Act`, dipanggil
// `SaveTreatyIn_Act` [8] / `SaveTreatyIn_EDM_Act` [14] bila
// `StatusAkseptasi == "Resolve Complete"`.

// KolomDetail - satu kolom tabel detail, urut kolom `INSERT` prosedur
// `PEGA_M_TREATY_IN_DETAIL(_EDM)`. `Angka` = kolom `NUMBER` di DEV.
//
// ⛔ `ID`, `COMMENCEMENT`, dan `TERMINATION` tidak di sini: prosedur
// mengisinya sendiri (pengenal situs, tanggal dari kepala kontrak), dan
// repository menirunya.
type KolomDetail struct {
	Nama  string
	Angka bool
}

// kolomDetailBersama - kolom yang dimiliki KEDUA tabel, urut prosedur.
var kolomDetailBersama = []KolomDetail{
	{"TREATYID", false}, {"TREATYYEAR", false}, {"SOBID", false}, {"SOB", false},
	{"CEDINGID", false}, {"CEDING", false}, {"TREATYCONTRACTNAME", false},
	{"PROPORTIONTYPE", false}, {"TREATYTYPE", false}, {"TREATYGROUPID", false},
	{"TREATYGROUP", false}, {"CLASSOFBUSINESSID", false}, {"CLASSOFBUSINESS", false},
	{"LIMITCURRENCY", false}, {"LIMITVALUE", true},
	{"RETENTIONCURRENCY", false}, {"RETENTIONVALUE", true},
	{"EPICURRENCY", false}, {"EPIVALUE", true},
	{"MDPCURRENCY", false}, {"MDPVALUE", true},
	{"NETPREMICURRENCY", false}, {"NETPREMIVALUE", true},
	{"DEDUCTION1", true}, {"DEDUCTION2", true},
	{"LAYERTYPE", false}, {"LAYER", false}, {"LAYERPARTTYPE", false}, {"LAYERPART", false},
	{"INSTALLMENTNO", false}, {"SHARECURRENCY", false}, {"SHAREVALUE", true},
	{"RIOGR", true}, {"RIONR", true}, {"RNM_SHARE", true},
	{"SPREADINGTYPE", false}, {"SPREADINGTYPEID", false},
	{"QS_OR", true}, {"QS_RI", true}, {"BROKERAGE", true},
	{"CESSIONCURRENCY", false}, {"CESSIONVALUE", true},
	{"CURRENCY_RSMD", false}, {"CURRENCY_EQ", false},
	{"CURRENCY_FLOOD_JABO", false}, {"CURRENCY_FLOOD_NATION", false},
	{"LIMIT_RSMD", true}, {"LIMIT_EQ", true}, {"LIMIT_FLOOD_JABO", true}, {"LIMIT_FLOOD_NATION", true},
	{"SPL_LINE", true}, {"QS_PCT", true},
}

// KolomDetailTreatyIn - kolom `TREATYINDETAIL` yang `SaveTreatyInDetail`
// isi: kolom bersama ditambah sembilan kolom yang hanya dimiliki tabel ini.
//
// ⚠️ `ADJ_RATE` = `SaveData.ADJUSMENT_RATE` (parameter `P_ADJUSMENT_RATE`) —
// satu-satunya kolom yang namanya berbeda dari properti `SaveData`-nya.
var KolomDetailTreatyIn = append(append([]KolomDetail{}, kolomDetailBersama...),
	KolomDetail{"DEDUCTIBLE", true}, KolomDetail{"DEDUCTIBLE2", true}, KolomDetail{"ADJ_RATE", true},
	KolomDetail{"PREMIUM_EARNED", true}, KolomDetail{"MDP_PCT", true}, KolomDetail{"ROL_PCT", true},
	KolomDetail{"MDP", true}, KolomDetail{"SPREAD_RNM_SHARE_PCT", true}, KolomDetail{"SPREAD_RNM_SHARE_VALUE", true},
)

// KolomDetailTreatyInEDM - kolom `TREATYINDETAILEDM` yang
// `SaveTreatyInDetailEdm` isi: kolom bersama saja (prosedur EDM berhenti di
// `QS_PCT`).
var KolomDetailTreatyInEDM = append([]KolomDetail{}, kolomDetailBersama...)

// BarisDetailTreaty - satu baris detail siap tulis: teks apa adanya,
// angka sebagai desimal (`nil` = `NULL`).
type BarisDetailTreaty struct {
	Teks  map[string]string
	Angka map[string]*apd.Decimal
}

// RencanaDetail - tulisan detail satu kontrak tuntas.
//
// ⭐ `Baris` KOSONG tetap berarti "hapus lalu tulis nol baris" — persis
// `RemoveTreatyInDetail(Edm)` yang selalu jalan lebih dulu, apa pun
// hasil susunannya.
type RencanaDetail struct {
	Baris []BarisDetailTreaty
}
