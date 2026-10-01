package models

// Bentuk kasus endorsement, versi polis lama, dan peserta.
//
// ⛔ Uang dan angka menyeberang sebagai TEKS desimal (ADR-0003) - tidak pernah
// `float64`. Tanggal teks `YYYY-MM-DD`; kosong = "".
//
// ⛔ Nol nama orang: `CREATE_OP_NAME` dan `PIC_SUGGEST` diisi pengenal AKUN
// pelaku (ADR-U-0030), bukan nama tampilan.

import (
	"regexp"
	"strings"
)

// JenisSumber membedakan dari mana versi polis lama dibaca (RALAT bab 4).
type JenisSumber string

const (
	// SumberAplikasi - versi sistem baru: `T_PREMIUM_LIST` + anak-anaknya.
	SumberAplikasi JenisSumber = "aplikasi"
	// SumberWarisan - versi sistem lama: `JSON_POLIS` + `M_LIFE_PREMIUM_DETAIL`.
	SumberWarisan JenisSumber = "warisan"
)

// Versi adalah versi BERJALAN sebuah polis - `PRODKE` terbesar di kedua sumber.
type Versi struct {
	Jenis JenisSumber
	// ID - `T_PREMIUM_LIST.ID` (aplikasi) atau `JSON_POLIS.IDPEGA` (warisan).
	ID string
	// ProdKe - `PROD_KE`/`PRODKE`; kosong dibaca 1 (E1, RALAT R23).
	ProdKe int
	// EdmType - maksud endorsement versi itu; kosong pada new business.
	// Gerbang 4: memuat `3` = polis sudah pernah Batal (`GetEdmTypeLife`).
	EdmType string
}

// SudahBatal - gerbang 4: `@contains(OutData1.pxResults(1).CARI1,"3")`
// (`SetErrorBatalEndorsement_Act` 3.8 b1995).
func (v Versi) SudahBatal() bool { return strings.Contains(v.EdmType, EdmTypeBatal) }

// LebihBaru memilih versi berjalan di antara dua kandidat: `PRODKE` terbesar;
// seri dimenangkan sumber aplikasi (penomoran sistem baru melanjutkan warisan).
func LebihBaru(a, b Versi) Versi {
	switch {
	case a.ID == "":
		return b
	case b.ID == "":
		return a
	case b.ProdKe > a.ProdKe:
		return b
	case b.ProdKe == a.ProdKe && b.Jenis == SumberAplikasi && a.Jenis != SumberAplikasi:
		return b
	}
	return a
}

// KolomKepala adalah satu kolom kepala `T_PREMIUM_LIST` yang disalin dari versi
// lama ke kasus - `MappingEDMLife` langkah 9 b1690 (PRE=false, selalu jalan).
type KolomKepala struct {
	Kolom string
	// Properti - nama properti work page Pega, kunci `JSON_POLIS.DATA_JSON`
	// untuk sumber warisan.
	Properti string
	Tanggal  bool
}

// KolomKepalaSalin - 24 properti langkah 9 (yang ke-25, `PL_NUMBER` b2184,
// adalah nomor polis itu sendiri). VERBATIM urutan korpus.
var KolomKepalaSalin = []KolomKepala{
	{"TYPE", "Type", false},                                 // b1717
	{"TYPE_CEDING", "TypeCeding", false},                    // b1764
	{"BUSINESS_CODE", "BusinessCode", false},                // b1785
	{"SOB_NAME", "SobName", false},                          // b1806
	{"POLICY_HOLDER_NAME", "PolicyHolderName", false},       // b1827
	{"CEDING_CO_NAME", "CedingCoName", false},               // b1848
	{"MARKETING_NAME", "MarketingName", false},              // b1869
	{"PRO_RATE_TYPE", "ProRateType", false},                 // b1890
	{"BUSINESS_NAME", "BusinessName", false},                // b1911
	{"POLICY_HOLDER", "PolicyHolder", false},                // b1932
	{"SOB", "SourceOfBusiness", false},                      // b1953
	{"CEDING_CO", "CedingCo", false},                        // b1974
	{"MO_ID", "MOID", false},                                // b1995
	{"MARKETING_CODE", "MarketingCode", false},              // b2016
	{"DATE_RECEIVED", "DateReceived", true},                 // b2037
	{"NO_OFFER", "NoOffer", false},                          // b2058
	{"RETRO_ID", "RetroID", false},                          // b2079
	{"RETRO_NAME", "RetroName", false},                      // b2100
	{"SECURITY_REINSURER_ID", "SecurityReinsurerID", false}, // b2121
	{"SECURITY_REINSURER", "SecurityReinsurer", false},      // b2142
	{"RI_SLIP_RNM", "PremiumListSummary.RISLIPRNM", false},  // b2163
	{"PRODUCT_NAME", "ProductName", false},                  // b2205 (R09)
	{"PRODUCT_NAME_ID", "ProductNameID", false},             // b2226 (R09)
	{"WPC", "WPC", true},                                    // b2247 (R09)
}

// KasusBaru adalah isian `EndorsmentLife_Section` yang membuat kasus.
type KasusBaru struct {
	NomorPolis string // `TempWork.PolicyNo` b1075
	EdmType    string // `TempWork.EdmType` b1347
	EdmDate    string // `TempWork.EdmDate` b3160, `YYYY-MM-DD`
	Deskripsi  string // `TempWork.DesBatal` b2140 → `EDM_NOTE` (R13)
}

// Kasus adalah kepala satu kasus endorsement, sebagaimana `InputEDMLife`
// membutuhkannya. Seluruh medan baca-saja di layar (`ro`).
type Kasus struct {
	ID          string `json:"id"`
	NomorPolis  string `json:"policyNo"`
	PLNumber    string `json:"plNumber"`
	PLNumberEDM string `json:"plNumberEdm"`
	ProdKe      int    `json:"prodKe"`
	EdmType     string `json:"edmType"`
	EdmNote     string `json:"description"`
	EdmDate     string `json:"edmDate"`
	Status      string `json:"status"`
	TglInput    string `json:"createDate"`
	Pembuat     string `json:"createOperator"`
	// Kepala - kolom `KolomKepalaSalin`, kunci nama kolom.
	Kepala map[string]string `json:"kepala"`
	// SudahSimpan - `IsJsonPolis` (`SetPremi_EDM` 8 b5600): rekap kasus ada.
	SudahSimpan bool `json:"sudahSimpan"`
	// CSVTerkunci - `EditInput1` (`SaveCSVEDMLife` 5 b2899): ada peserta `New`.
	CSVTerkunci bool `json:"csvTerkunci"`
	// Cacah - jumlah peserta per `EDM_STATUS`.
	Cacah map[string]int `json:"cacah"`
	// Sumber - jalur versi lama yang disalin (bab 4 RALAT).
	Sumber JenisSumber `json:"sumber"`
	// SumberID - `T_PREMIUM_LIST.ID` versi lama, atau `JSON_POLIS.IDPEGA`.
	SumberID string `json:"-"`
}

// Terbuka - kasus belum diputuskan.
func (k Kasus) Terbuka() bool { return k.Status == StatusKasusTerbuka }

// BarisInbox adalah satu baris kotak masuk `InboxEndorsementLife` (grid b8284).
type BarisInbox struct {
	CaseID         string `json:"caseId"`
	EndorsementNo  string `json:"endorsementNo"`
	Type           string `json:"type"`
	EdmType        string `json:"edmType"`
	PolicyNo       string `json:"policyNo"`
	Sob            string `json:"sob"`
	Ceding         string `json:"ceding"`
	PolicyHolder   string `json:"policyHolder"`
	MarketingName  string `json:"marketingName"`
	CreateDate     string `json:"createDate"`
	CreateOperator string `json:"createOperator"`
	Status         string `json:"status"`
}

// Peserta adalah satu baris `T_PREMIUM_LIST_DETAIL` kasus.
type Peserta struct {
	ID        string `json:"id"`
	ParentID  string `json:"parentId"`
	EdmStatus string `json:"edmStatus"`
	// Terkunci - kotak centang `.EdmBatal` mati (`InputEDMLife.xml` b15763
	// `.EditInput==1`): peserta yang sudah `Delete`/`Batal`.
	Terkunci bool `json:"terkunci"`
	// Nilai - kolom layar, kunci nama kolom, teks.
	Nilai map[string]string `json:"nilai"`
}

// JenisKolom - cara sebuah kolom peserta dibaca ke teks.
type JenisKolom int

const (
	KolomTeks JenisKolom = iota
	KolomTanggal
	KolomAngka
)

// KolomPeserta - satu kolom layar peserta.
type KolomPeserta struct {
	Nama  string
	Jenis JenisKolom
}

// KolomPesertaGrid - kolom grid `InputEDMLife` b11899 (sel b13805 … b15587),
// sama dengan grid Batal b17500 tanpa `POLICY_HOLDER`.
var KolomPesertaGrid = []KolomPeserta{
	{"POLICY_NO", KolomTeks}, {"POLICY_HOLDER", KolomTeks}, {"CERTIFICATE_NO", KolomTeks},
	{"NAME_OF_INSURED", KolomTeks}, {"SEX", KolomTeks}, {"DOB", KolomTanggal},
	{"ENTRY_AGE", KolomAngka}, {"PLAN", KolomTeks}, {"BEGIN_DATE", KolomTanggal},
	{"EFFECTIVE_DATE", KolomTanggal}, {"EXPIRED_DATE", KolomTanggal},
}

// KolomPesertaRinci - kolom `PL_Detail_Sec` keempat `Type` (rincian baris,
// `pyEditAction` `PL_DetailAction`) dan grid popup polis lama, gabungan.
var KolomPesertaRinci = []KolomPeserta{
	{"POLICY_NO", KolomTeks}, {"POLICY_HOLDER", KolomTeks}, {"CERTIFICATE_NO", KolomTeks},
	{"NAME_OF_INSURED", KolomTeks}, {"SEX", KolomTeks}, {"DOB", KolomTanggal},
	{"AGE", KolomAngka}, {"ENTRY_AGE", KolomAngka}, {"CURRENT_AGE", KolomAngka},
	{"PLAN", KolomTeks}, {"STNC", KolomTeks}, {"WPC", KolomTeks}, {"CURRENCY", KolomTeks},
	{"BEGIN_DATE", KolomTanggal}, {"EFFECTIVE_DATE", KolomTanggal}, {"EXPIRED_DATE", KolomTanggal},
	{"PERIOD_YY", KolomAngka}, {"PERIOD_MM", KolomAngka}, {"PASSED_PERIOD", KolomAngka},
	{"DESCRIPTION", KolomTeks}, {"RISK", KolomTeks},
	{"SUM_INSURED", KolomAngka}, {"CEDING_RETENTION", KolomAngka}, {"SUM_REASURED", KolomAngka},
	{"SHARE_NUSANTARA_RE", KolomAngka}, {"SHARE_NUSANTARA_RE_GROSS", KolomAngka},
	{"SUM_AT_RISK_GROSS", KolomAngka}, {"SUM_AT_RISK_RETRO", KolomAngka},
	{"RETROCEDED_SHARE", KolomAngka}, {"SHARE_RETRO", KolomAngka}, {"RATE", KolomAngka},
	{"FACTOR", KolomAngka}, {"EM_PERCENT", KolomAngka},
	{"GROSS_PREMIUM", KolomAngka}, {"DEDUCTION", KolomAngka}, {"NET_PREMIUM", KolomAngka},
	{"BROKERAGE_FEE", KolomAngka}, {"CLAIM_AMOUNT", KolomAngka},
	{"GROSS_PREMIUM_REFUND", KolomAngka}, {"DEDUCTION_REFUND", KolomAngka},
	{"NET_PREMIUM_REFUND", KolomAngka}, {"BROKERAGE_FEE_REFUND", KolomAngka},
	{"GROSS_PREMIUM_RETRO", KolomAngka}, {"DISCOUNT_PREMIUM_RETRO", KolomAngka},
	{"RI_ADMIN_FEE_RETRO", KolomAngka}, {"BROKERAGE_FEE_RETRO", KolomAngka},
	{"NET_PREMIUM_RETRO", KolomAngka}, {"GROSS_PREMIUM_REFUND_RETRO", KolomAngka},
	{"DISCOUNT_PREMIUM_REFUND_RETRO", KolomAngka}, {"RI_ADMIN_FEE_REFUND_RETRO", KolomAngka},
	{"NET_PREMIUM_REFUND_RETRO", KolomAngka},
}

// polaKolom - nama kolom yang sah dirakit ke SQL (huruf besar, angka, `_`).
var polaKolom = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// KolomSah menjawab apakah sebuah nama kolom boleh dirakit ke SQL.
func KolomSah(nama string) bool { return polaKolom.MatchString(nama) }
