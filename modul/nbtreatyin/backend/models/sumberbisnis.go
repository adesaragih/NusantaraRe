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
// `KirimanSumberBisnis` (pola `TerimaKirimanTerkunci`, Save/Submit/refresh,
// `services.terimaSumberBisnis`).

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

// SumberBisnisPostDT - keempat medan `Quotation.*` yang ditulis
// `SearchHierarkiSourceBizAgent_PostDT` (langkah 1.1, 1.2, 4, 5). Inilah yang
// DIPEGANG layar sampai Save/Submit (F4) - PostDT hanya menulis clipboard,
// nol Obj-Save.
type SumberBisnisPostDT struct {
	SourceOfBusiness string `json:"sourceOfBusiness"`
	SobName          string `json:"sobName"`
	SobLeader0       string `json:"sobLeader0"`
	SobLeader1       string `json:"sobLeader1"`
}

// MedanSumberBisnis - medan halaman Quotation yang ditulis PostDT.
var MedanSumberBisnis = []string{"SourceOfBusiness", "SobName", "SobLeader0", "SobLeader1"}

const jq = HalamanQuotation + "."

// BacaSumberBisnis - keempat medan PostDT dari halaman.
func BacaSumberBisnis(h *Halaman) SumberBisnisPostDT {
	return SumberBisnisPostDT{
		SourceOfBusiness: h.Ambil(jq + "SourceOfBusiness"),
		SobName:          h.Ambil(jq + "SobName"),
		SobLeader0:       h.Ambil(jq + "SobLeader0"),
		SobLeader1:       h.Ambil(jq + "SobLeader1"),
	}
}

// TerapkanPilihanSumberBisnis menulis keempat medan ke `Quotation.<medan>`
// (persis hasil PostDT) lalu ke salinannya `PolicyTreatyIn.QuotationData.<medan>`
// (`SalinKeQuotationData`).
func TerapkanPilihanSumberBisnis(h *Halaman, x SumberBisnisPostDT) {
	h.Setel(jq+"SourceOfBusiness", x.SourceOfBusiness)
	h.Setel(jq+"SobName", x.SobName)
	h.Setel(jq+"SobLeader0", x.SobLeader0)
	h.Setel(jq+"SobLeader1", x.SobLeader1)
	SalinKeQuotationData(h, MedanSumberBisnis...)
}

// HasilPostDT = tombol `Select Source Of Business` (`btnSOB_DT` langkah 1) lalu
// klik baris `b` (`SearchHierarkiSourceBizAgent_PostDT`).
func HasilPostDT(b BarisAgen) (SumberBisnisPostDT, error) {
	h := HalamanBaru()
	TombolSOB(h)
	if err := TerapkanSumberBisnis(h, b); err != nil {
		return SumberBisnisPostDT{}, err
	}
	return BacaSumberBisnis(h), nil
}

// ---------------------------------------------------------------- F4: Save/Submit

// KirimanSumberBisnis = F4 dengan pola `TerimaKirimanTerkunci`: pilihan Source
// Of Business yang dipegang layar admin ikut kiriman Save/Submit/refresh.
//
//   - pembeda = `Quotation.SourceOfBusiness`, satu-satunya medan PostDT yang
//     berkolom (T_POLIS_QUOTATION.SOURCE_OF_BUSINESS) dan yang dibaca rule NB
//     (`SetPPNPPH` langkah 1, lewat salinan QuotationData); tiga medan lainnya
//     (tanpa kolom, nol pembaca NB) bukan pemicu;
//   - ClaimType (isian layar) bukan 'XOL Retro' -> tombol tidak tampil, medan
//     TERKUNCI (pola AC 49-51);
//   - hasil hitung ulang = PostDT atas setiap baris `daftar` (RD
//     `BrowseAgentHierarkiList_RD` yang dijalankan ulang): pilihan kosong
//     semua cocok dengan simpul beranak mana pun (`ChildCount > 0`); baris yang
//     ChildCount-nya bukan angka tidak dapat diklik (PostDT gagal) - dilewati;
//   - diterima -> keempat medan ditulis ke Quotation dan salinan QuotationData
//     (`TerapkanPilihanSumberBisnis`).
func KirimanSumberBisnis(h *Halaman, daftar func() ([]BarisAgen, error)) KirimanTerkunci[SumberBisnisPostDT] {
	return KirimanTerkunci[SumberBisnisPostDT]{
		Jalur:   jalurSumberBisnis(),
		Tampil:  TampilPilihSumberBisnis(h),
		Baca:    BacaSumberBisnis,
		Tulis:   TerapkanPilihanSumberBisnis,
		Berubah: func(k, s SumberBisnisPostDT) bool { return k.SourceOfBusiness != s.SourceOfBusiness },
		Pesan:   func(k SumberBisnisPostDT) string { return PesanSumberBisnisTidakCocok(k.SourceOfBusiness) },
		HitungUlang: func() ([]SumberBisnisPostDT, error) {
			d, err := daftar()
			if err != nil {
				return nil, err
			}
			out := make([]SumberBisnisPostDT, 0, len(d))
			for _, b := range d {
				if x, err := HasilPostDT(b); err == nil {
					out = append(out, x)
				}
			}
			return out, nil
		},
	}
}

func jalurSumberBisnis() []string {
	out := make([]string, len(MedanSumberBisnis))
	for i, m := range MedanSumberBisnis {
		out[i] = jq + m
	}
	return out
}

// PesanSumberBisnisTidakCocok - pesan validasi (422) bila pilihan yang dipegang
// layar tidak cocok dengan hasil pencarian hierarki yang dijalankan ulang.
func PesanSumberBisnisTidakCocok(id string) string {
	return fmt.Sprintf("Source Of Business %q tidak cocok dengan hasil pencarian hierarki sumber bisnis "+
		"(BrowseAgentHierarkiList_RD) - pilih ulang lewat tombol Select Source Of Business", id)
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
