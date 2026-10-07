package models

// Untuk apa berkas ini: LAHIRNYA KASUS ENDORSEMEN - port `Activity/CreateEDMT` (kelas Data-Portal, ruleset 01-01-91;
// 18 langkah aktif, 3 `//`), `Activity/SetEDMTNoPolis` (01-01-56), dan pesan pemeriksaan layar `TreatyCreateEdm`
// (`Activity/TrtEdmCheckPolicyError` 01-01-76, `Activity/CheckNopolisAvailability` 01-01-77).
//
//	CreateEDMT 1        InputData.CARI1 = DisplayData.CARI1 (nomor polis diketik); TrtERR.CARI1 = ""
//	           2-3      Page-New TempPage, MergePage
//	           4-6      `//` CheckNopolisAvailability, CheckDupeAdjPremi, CheckOngoingTreatyEDM - TIDAK diport
//	                    (pemeriksaan yang hidup hanya di layar: services.PeriksaPolis)
//	           7-8      RDB FetchPolisJsonPolis (json_polis PRODKE terbesar) + adoptJSONObject -> TempPage.PolicyTreatyIn
//	                    ⛔ di sistem baru: generasi terakhir dari T_GENERAL_POLIS_TREATY (repository.GenerasiTerakhir)
//	           9        TempPage.PolicyTreatyIn.PolicyNo = InputData.CARI1
//	           10       Page-Copy TempPage.PolicyTreatyIn -> MergePage.PolicyTreatyIn.OldData
//	           11       Page-Copy OldData.QuotationData -> MergePage.Quotation
//	           12       OfferFacIn.PolicyData.PolicyNo = OldData.PolicyNo; Quotation.OldPolicyNo = OldData.PolicyNo;
//	                    OfferFacIn.QuotationData = Quotation; PolicyTreatyIn.QuotationData = Quotation
//	           13       Page-Remove TempPage
//	           14       EDMType = Param.edmtype; IDCurrency, Currency, ClaimPaymentType, ClaimType = OldData.*;
//	                    NBStatus = "EDM IS IN " + OperatorID.pxInsName + "'S INBOX"
//	           15       Call SetEDMTNoPolis
//	           16       Call SetEDMTCancel bila Param.edmtype==4 - atas pyWorkPage PORTAL, sebelum kasus lahir dari
//	                    MergePage (18): [dugaan] tanpa efek pada kasus baru; yang berlaku EDMChooseBusiness_Act 5
//	           17       keluar bila TrtERR.CARI1 != "" (mati: langkah 1 mengosongkannya, 4-6 `//`)
//	           18-21    CreateWorkPage (ASM-FW-GISFW-Work-EndorsementTreaty, MergePage, pyStartCase), AddWork,
//	                    commit, DisplayData.CARI9 = pzInsKey (buka kasus) - services / repository
//
// Halaman `OfferFacIn` (kelas FacIn) tidak disimpan: `OfferFacIn.PolicyData.PolicyNo` dibaca
// `serviceInsertArasapas_act` (REST noPolis) - diturunkan dari PolicyNo saat dipakai.

import (
	"strconv"
	"strings"
)

// Pesan VERBATIM layar Create (`TrtERR.CARI1`).
const (
	// PesanEDMBelumSelesai = `TrtEdmCheckPolicyError` langkah 4.
	PesanEDMBelumSelesai = "There's EDM with this policy no that haven't finish yet!"
	// PesanNopolisSalah = `CheckNopolisAvailability` langkah 3.
	PesanNopolisSalah = "Error: Nopolis is incorrect / doesn't exists"
	// PesanGenerasiBelumDimuat - ⛔ BUKAN teks XML: nomor polis ada di json_polis tetapi generasinya belum ada di
	// tabel relasional (dokumen lama belum dipindah pemuat). Pega membaca OldData dari DATA_JSON; sistem baru dari
	// T_GENERAL_POLIS_TREATY, jadi keadaan ini tidak dikenal XML.
	PesanGenerasiBelumDimuat = "Error: policy data has not been migrated to the new system yet"
)

// NomorEDM = `SetEDMTNoPolis` langkah 3: EDMNo = OldData.PolicyNo + "/E" + (ProdKe < 10 ? "0" + ProdKe : ProdKe).
// Generasi >= 100 menjadi tiga digit (tanpa batas - KEPUTUSAN 23-09-2026 butir 1).
func NomorEDM(nopolis string, prodKe int) string {
	n := strconv.Itoa(prodKe)
	if prodKe < 10 {
		n = "0" + n
	}
	return nopolis + "/E" + n
}

// PindahAwalan menyalin nilai dan daftar `src` berjalur awalan `dari` ke `dst` dengan awalan `ke` (Page-Copy
// halaman bersarang).
func PindahAwalan(src, dst *Halaman, dari, ke string) {
	if src == nil {
		return
	}
	for j, v := range src.Nilai {
		if strings.HasPrefix(j, dari) {
			dst.Setel(ke+strings.TrimPrefix(j, dari), v)
		}
	}
	for j, d := range src.Daftar {
		if strings.HasPrefix(j, dari) {
			dst.SetelDaftar(ke+strings.TrimPrefix(j, dari), salinBaris(d))
		}
	}
}

// PasangOldData = `CreateEDMT` langkah 10 (Page-Copy halaman generasi terakhir -> `PolicyTreatyIn.OldData`):
// seluruh `PolicyTreatyIn.*` generasi lama - termasuk EDMNo, QuotationData, daftar, dan TreatyDifference miliknya
// (tab `PropOldData2`) - menjadi `PolicyTreatyIn.OldData.*`. Satu tingkat saja (spec ID-8).
func PasangOldData(h, lama *Halaman) {
	PindahAwalan(lama, h, HalamanPolis+".", od)
}

// RakitHalamanBaru = `CreateEDMT` langkah 9-15 atas generasi terakhir `lama` (halaman satu generasi, jalur
// `PolicyTreatyIn.*` / `Quotation.*`): halaman kasus endorsemen yang baru lahir. `edmType` = Param.edmtype,
// `akunPembuat` = OperatorID.pxInsName, `prodKe` = PRODKE generasi terakhir + 1.
//
// Data baru (premi, komisi, spreading, angsuran, XOL) KOSONG: MergePage baru - diisi tombol Choose Business
// (`EDMChooseBusiness_Act`).
func RakitHalamanBaru(lama *Halaman, nopolis, edmType, akunPembuat string, prodKe int) *Halaman {
	h := HalamanBaru()
	// 9-10
	lama.Setel(pt+"PolicyNo", nopolis)
	PasangOldData(h, lama)
	// 11-12: Quotation = OldData.QuotationData (Page-Copy, termasuk daftar CedingCoList / SurveyReportList)
	PindahAwalan(h, h, od+"QuotationData.", HalamanQuotation+".")
	h.Setel(HalamanQuotation+".OldPolicyNo", h.Ambil(od+"PolicyNo"))
	SalinQuotation(h)
	for j, d := range h.Daftar {
		if strings.HasPrefix(j, od+"QuotationData.") {
			h.SetelDaftar(pt+"QuotationData."+strings.TrimPrefix(j, od+"QuotationData."), salinBaris(d))
		}
	}
	// 14
	h.Setel(pt+"EDMType", edmType)
	for _, m := range []string{"IDCurrency", "Currency", "ClaimPaymentType", "ClaimType"} {
		h.Setel(pt+m, h.Ambil(od+m))
	}
	h.Setel("NBStatus", TeksNBStatusBaru(akunPembuat))
	// 15 SetEDMTNoPolis
	h.Setel(pt+"ProdKe", strconv.Itoa(prodKe))
	h.Setel(pt+"EDMNo", NomorEDM(h.Ambil(od+"PolicyNo"), prodKe))
	h.Setel(pt+"PolicyNo", h.Ambil(od+"PolicyNo"))
	// connector Start1 -> Assignment2
	h.Setel("Position", PositionAdmin)
	h.Setel("PositionNote", PosisiAdmin)
	h.Setel("FlagOnGoingPolicy", "1")
	WarisPajakLama(h) // keputusan WO 07-10-2026 - header With Tax / Type Tax sejak lahir
	return h
}

// WarisPajakLama - keputusan work owner 07-10-2026 ("CURRENT PREMIUM ... GA ADA ANGKANYA"): With Tax (`.FlagPPH`)
// dan Type Tax (`.TypeTax`) generasi baru yang BELUM pernah diisi diwarisi dari generasi lama. Di korpus EDM nol rule
// menulis keduanya (hanya sel header), sehingga grid XOL Current Premium (`InsertToTreatyXOLListEDM` 3.3.7 /
// 3.3.8.5) tanpa pajak bila admin tidak mencentang With Tax SEBELUM Choose Business. Nilai "false" (centang dilepas
// admin) dihormati - penyimpangan sadar.
func WarisPajakLama(h *Halaman) {
	if h.Ambil(pt+"FlagPPH") != "" {
		return
	}
	h.Setel(pt+"FlagPPH", h.Ambil(od+"FlagPPH"))
	if h.Ambil(pt+"TypeTax") == "" {
		h.Setel(pt+"TypeTax", h.Ambil(od+"TypeTax"))
	}
}

// KodeJenisEDM - pilihan "Source of Change" (`Section/TreatyCreateEdm` `DisplayData.CARI3`, dropdown
// `TypeList.pxResults` dari DT `TreatyEDMListType`: nilai `.CARI1`, teks `.CARI2`). DT itu TIDAK ada di korpus -
// isinya dari screenshot Pega yang dikirim work owner 07-10-2026 (lihat `LabelJenisEDM`). Kode yang dibaca rule:
// "2" dan "4" (`CalculateDifferenceEDM_act` 1.2.5 "EDM BATAL"), "3" (`SetValueEDM_Act` 9-10 AdjPremi,
// `BusinessAndSOBListEDM` S1, `InputPolicyTreatyInPre_Act` 10), "4" (`EDMChooseBusiness_Act` 5 SetEDMTCancel).
var KodeJenisEDM = []string{"1", "2", "3", "4"}

// LabelJenisEDM - DT `TreatyEDMListType` (screenshot Pega, work owner 07-10-2026) VERBATIM:
// TypeList.pxResults(n).CARI1 -> .CARI2. Sel `.EDMType` portal (`SFAPortal_Endorsement_Treaty`) dan header
// (`DetailPolicyTreatyInAddendum`) juga pxDropdown (sumber `associated`, prompt values properti tidak ada di
// korpus) - teks yang sama dipakai di sana.
var LabelJenisEDM = map[string]string{"1": "Internal", "2": "External", "3": "Adjustment Premium", "4": "Cancel Input"}

// HitungTotalSpreadingOldData - total spreading generasi lama (`pyWorkPage.PolicyTreatyIn.OldData.Total*`, kaki
// grid tab Old Data `PropOldData`): turunan baris OldData.SpreadingRiskList, tidak disimpan (pola
// `HitungTotalSpreading`).
func HitungTotalSpreadingOldData(h *Halaman) error {
	if len(h.AmbilDaftar(od+"SpreadingRiskList")) == 0 {
		return nil
	}
	sem := HalamanBaru()
	sem.SetelDaftar(DaftarSpreading, h.AmbilDaftar(od+"SpreadingRiskList"))
	if err := HitungTotalSpreading(sem); err != nil {
		return err
	}
	for _, m := range []string{"TotalSharePercentagePremium", "TotalPremium", "TotalSharePercentageClaim", "TotalClaim"} {
		h.Setel(od+m, sem.Ambil(pt+m))
	}
	return nil
}
