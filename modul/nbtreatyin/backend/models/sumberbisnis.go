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
//
// F4 (keputusan WO 04-10-2026: IKUTI XML) - hasil PostDT DIPEGANG layar dan
// disimpan bersama Save/Submit admin: `HasilPostDT` (klik), lalu
// `SumberBisnisKiriman` + `CocokHasilPostDT` + `TerapkanPilihanSumberBisnis`
// (Save/Submit/refresh, `services.terimaSumberBisnis`).

import "fmt"

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

// HasilPostDT = tombol `Select Source Of Business` (`btnSOB_DT` langkah 1) lalu
// klik baris `b` (`SearchHierarkiSourceBizAgent_PostDT`): nilai keempat
// `MedanSumberBisnis` yang ditulis PostDT, per jalur halaman
// (`Quotation.<medan>`). Inilah yang DIPEGANG layar sampai Save/Submit (F4) -
// PostDT hanya menulis clipboard, nol Obj-Save.
func HasilPostDT(b BarisAgen) (map[string]string, error) {
	h := HalamanBaru()
	TombolSOB(h)
	if err := TerapkanSumberBisnis(h, b); err != nil {
		return nil, err
	}
	out := make(map[string]string, len(MedanSumberBisnis))
	for _, m := range MedanSumberBisnis {
		out[HalamanQuotation+"."+m] = h.Ambil(HalamanQuotation + "." + m)
	}
	return out, nil
}

// ---------------------------------------------------------------- F4: Save/Submit

// SumberBisnisKiriman - pilihan Source Of Business yang dipegang layar
// (`masuk`, kiriman Save/Submit/refresh admin) BILA berbeda dari halaman server
// `h` (tersimpan + pra-proses). Pembanding = `Quotation.SourceOfBusiness`,
// satu-satunya medan PostDT yang berkolom (T_POLIS_QUOTATION.SOURCE_OF_BUSINESS)
// dan yang dibaca rule NB (`SetPPNPPH` langkah 1, lewat salinan QuotationData).
// Jalur itu tidak dikirim, atau sama = tidak ada pilihan baru: tiga medan
// lainnya (tanpa kolom, nol pembaca NB) tidak diambil dari layar.
//
// `pilihan` memuat keempat `MedanSumberBisnis` per jalur halaman; medan yang
// tidak dikirim bernilai "".
func SumberBisnisKiriman(h, masuk *Halaman) (pilihan map[string]string, berubah bool) {
	j := HalamanQuotation + ".SourceOfBusiness"
	if masuk == nil {
		return nil, false
	}
	v, ada := masuk.Nilai[j]
	if !ada || v == h.Ambil(j) {
		return nil, false
	}
	pilihan = make(map[string]string, len(MedanSumberBisnis))
	for _, m := range MedanSumberBisnis {
		pilihan[HalamanQuotation+"."+m] = masuk.Ambil(HalamanQuotation + "." + m)
	}
	return pilihan, true
}

// CocokHasilPostDT - adakah baris `daftar` (RD `BrowseAgentHierarkiList_RD`
// yang dijalankan ulang) yang hasil `HasilPostDT`-nya SAMA PERSIS dengan
// `pilihan` di keempat medan. Pilihan kosong semua cocok dengan simpul
// beranak mana pun (`ChildCount > 0`); pilihan berisi hanya dengan baris
// ChildCount 0 ber-ID, ClientName, dan Leader0 itu. Baris yang ChildCount-nya
// bukan angka tidak dapat diklik (PostDT gagal) - dilewati.
func CocokHasilPostDT(daftar []BarisAgen, pilihan map[string]string) bool {
	for _, b := range daftar {
		hasil, err := HasilPostDT(b)
		if err != nil {
			continue
		}
		sama := true
		for _, m := range MedanSumberBisnis {
			j := HalamanQuotation + "." + m
			if hasil[j] != pilihan[j] {
				sama = false
				break
			}
		}
		if sama {
			return true
		}
	}
	return false
}

// PesanSumberBisnisTidakCocok - pesan validasi (422) bila pilihan yang dipegang
// layar tidak cocok dengan hasil pencarian hierarki yang dijalankan ulang.
func PesanSumberBisnisTidakCocok(id string) string {
	return fmt.Sprintf("Source Of Business %q tidak cocok dengan hasil pencarian hierarki sumber bisnis "+
		"(BrowseAgentHierarkiList_RD) - pilih ulang lewat tombol Select Source Of Business", id)
}

// TerapkanPilihanSumberBisnis menulis pilihan yang sudah dicocokkan ke
// `Quotation.<medan>` (persis hasil PostDT) lalu ke salinannya
// `PolicyTreatyIn.QuotationData.<medan>` (`SalinKeQuotationData`).
func TerapkanPilihanSumberBisnis(h *Halaman, pilihan map[string]string) {
	for _, m := range MedanSumberBisnis {
		h.Setel(HalamanQuotation+"."+m, pilihan[HalamanQuotation+"."+m])
	}
	SalinKeQuotationData(h, MedanSumberBisnis...)
}

// SalinKeQuotationData menyalin medan `Quotation.<m>` ke salinannya
// `PolicyTreatyIn.QuotationData.<m>` = `InputPolicyTreatyIn_preDT` langkah 14
// (`PolicyTreatyIn.QuotationData = pyWorkPage.Quotation`, `SalinQuotation`)
// diterapkan atas pilihan yang dipegang layar.
//
// ⚠️ `[penyesuaian sadar]`, bukan langkah PostDT. PostDT menulis `Quotation.*`
// saja; `SetPPNPPH` langkah 1 membaca `PolicyTreatyIn.QuotationData.
// SourceOfBusiness`, dan satu-satunya penyalin di korpus adalah pra-proses
// flow action (langkah 14 di atas; juga preACT 14.9 dan
// `GeneratePolicyNoTreaty_Act` langkah 10). XML tidak menjalankan refresh apa
// pun sesudah klik/Choose; apakah `opener.location.reload` (tombol Choose)
// membuat mesin Pega menjalankan ulang pra-proses itu TIDAK ada di korpus.
// Di sini pra-proses dijalankan ulang di SETIAP permintaan atas data tersimpan
// (`services.siapkan`), sehingga pilihan yang dipegang layar - yang belum
// tersimpan - harus disalin sesudah digabung; tanpa itu refresh, Save, dan
// Submit menghitung status PKP dari nilai tersimpan, bukan nilai layar. Kedua
// halaman juga disimpan SATU baris T_POLIS_QUOTATION dengan QuotationData
// didahulukan (`NilaiQuotation`, ID-23) - tanpa salinan ini pilihan tertimpa
// salinan lama saat disimpan.
func SalinKeQuotationData(h *Halaman, medan ...string) {
	for _, m := range medan {
		h.Setel(HalamanPolis+".QuotationData."+m, h.Ambil(HalamanQuotation+"."+m))
	}
}
