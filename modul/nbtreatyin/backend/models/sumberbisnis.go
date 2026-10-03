package models

// Untuk apa berkas ini: PEMILIH SOURCE OF BUSINESS (XOL Retro) layar admin -
// port murni `DataTransform/btnSOB_DT` dan
// `DataTransform/SearchHierarkiSourceBizAgent_PostDT`.
//
// Rantai Pega (korpus `NB Treaty In (Done)`):
//
//	Section/DetailPolicyTreatyIn  tombol "Select Source Of Business", pyVisible
//	                              `.ClaimType = 'XOL Retro'`, click -> showHarness
//	                              `SOB` (pyUsingPage pyWorkPage.Quotation, pySubmitData
//	                              Yes) dengan aktivitas `InputQuotation_PreAct(Acton=SOB)`
//	Activity/InputQuotation_PreAct langkah 1 (`Param.Acton=="SOB"`) -> `btnSOB_DT`
//	Harness/SOB -> Section/SourceHierarki  TreeGrid `TempBusinessSource.pxResults`,
//	                              pyDeferLoadActivity `AgentSourceBizTreatyIn_Act`
//	                              (RD `BrowseAgentHierarkiList_RD`), pyRowEditing
//	                              masterDetail, pyEditAction `AgentSourceBizDetails`
//	FlowAction/AgentSourceBizDetails  pyPreProcessingTransformRule =
//	                              `SearchHierarkiSourceBizAgent_PostDT` (klik baris)

// BarisAgen adalah satu baris RD `BrowseAgentHierarkiList_RD` (kelas
// `ASM-FW-GISFW-Int-AGENT`): kelima kolom RD, apa adanya sebagai teks.
type BarisAgen struct {
	ID         string `json:"id"`
	ClientName string `json:"clientName"`
	Leader0    string `json:"leader0"`
	ChildCount string `json:"childCount"`
	ClientID   string `json:"clientId"`
}

// KlaimXOLRetro - nilai `.ClaimType` yang memunculkan tombol.
const KlaimXOLRetro = "XOL Retro"

// TampilPilihSumberBisnis = pyVisible tombol `Select Source Of Business`
// (`Section/DetailPolicyTreatyIn`, juga salinannya di `GeneralPolicyTreatyIn`):
// `.ClaimType = 'XOL Retro'` atas halaman PolicyTreatyIn. Tombol ini hanya ada
// di layar admin (flow action `InboxPolicyTreatyIn`); layar atasan
// (`DetailDeptHeadTreatyIn_UW`) tidak memuatnya.
func TampilPilihSumberBisnis(h *Halaman) bool {
	return h.Ambil(HalamanPolis+".ClaimType") == KlaimXOLRetro
}

// TombolSOB = `DataTransform/btnSOB_DT` langkah 1:
// `pyWorkPage.Quotation.btnQuotation = "SOB"`.
func TombolSOB(h *Halaman) {
	h.Setel(HalamanQuotation+".btnQuotation", "SOB")
}

// TerapkanSumberBisnis = `DataTransform/SearchHierarkiSourceBizAgent_PostDT`
// atas baris agen yang diklik (Primary = baris itu):
//
//	1   WHEN pyWorkPage.Quotation.btnQuotation=="SOB"
//	1.1   Quotation.SourceOfBusiness = @if(.ChildCount > 0, "", .ID)
//	1.2   Quotation.SobName          = @if(.ChildCount > 0, "", .ClientName)
//	4   Quotation.SobLeader0 = @if(.ChildCount > 0, "", .Leader0)
//	5   Quotation.SobLeader1 = @if(.ChildCount > 0, "", .Leader1)
//
// ⚠️ `.Leader1` BUKAN kolom RD `BrowseAgentHierarkiList_RD` (ID, ClientName,
// Leader0, ChildCount, ClientID) dan tidak disalin `AgentSourceBizTreatyIn_Act`
// langkah 4 (Page-Copy baris RD) - nilainya selalu "". Ditiru apa adanya.
// ⚠️ Simpul yang punya anak (`ChildCount > 0`) MENGOSONGKAN keempat medan -
// sumber bisnis yang sudah terpilih ikut terhapus. Ditiru apa adanya.
// ChildCount kosong = 0 (properti numerik kosong di ekspresi Pega).
//
// ⛔ Cabang 2 (`btnQuotation=="CedingCo"` -> Quotation.CedingCo/CedingCoName) dan
// cabang 3 (`=="CedingCedant"` -> `OfferFacIn.CedingCedant*`) TIDAK diport:
// tidak terjangkau di NB. Penulis `btnQuotation` di korpus hanya `btnSOB_DT`
// ("SOB") dan `btnCedingCO_DT` ("CedingCo"); yang kedua hanya lewat
// `InputQuotation_PreAct` langkah 2 (`Param.Acton=="Ceding"`), padahal satu-
// satunya pengirim Acton di korpus (`DetailPolicyTreatyIn` /
// `GeneralPolicyTreatyIn`) mengirim `Acton=SOB`; nilai "CedingCedant" ditulis
// nol rule.
func TerapkanSumberBisnis(h *Halaman, b BarisAgen) error {
	anak, err := AngkaTeks("ChildCount", b.ChildCount)
	if err != nil {
		return err
	}
	beranak := anak.Sign() > 0
	jika := func(v string) string {
		if beranak {
			return ""
		}
		return v
	}
	q := func(m, v string) { h.Setel(HalamanQuotation+"."+m, v) }
	if h.Ambil(HalamanQuotation+".btnQuotation") == "SOB" { // 1
		q("SourceOfBusiness", jika(b.ID)) // 1.1
		q("SobName", jika(b.ClientName))  // 1.2
	}
	q("SobLeader0", jika(b.Leader0)) // 4
	q("SobLeader1", jika(""))        // 5 - .Leader1 tidak ada di baris RD
	return nil
}

// MedanSumberBisnis - medan halaman Quotation yang ditulis PostDT.
var MedanSumberBisnis = []string{"SourceOfBusiness", "SobName", "SobLeader0", "SobLeader1"}

// SalinKeQuotationData menyalin medan `Quotation.<m>` ke salinannya
// `PolicyTreatyIn.QuotationData.<m>`.
//
// ⚠️ PENYESUAIAN PENYIMPANAN, bukan langkah PostDT: Pega menyalin seluruh
// halaman pada pra-proses berikutnya (`InputPolicyTreatyIn_preDT` langkah 14 =
// `SalinQuotation`). Di sini kedua halaman disimpan SATU baris
// T_POLIS_QUOTATION dan repository mendahulukan nilai QuotationData
// (`repository.nilaiQuotation`, ID-23) - tanpa salinan ini pilihan baru
// tertimpa salinan lama saat disimpan. Hasil akhirnya sama dengan sesudah
// pra-proses berikut di Pega.
func SalinKeQuotationData(h *Halaman, medan ...string) {
	for _, m := range medan {
		h.Setel(HalamanPolis+".QuotationData."+m, h.Ambil(HalamanQuotation+"."+m))
	}
}
