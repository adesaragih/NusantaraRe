package models

// Untuk apa berkas ini: PILIH BISNIS ENDORSEMEN - port rantai tombol Choose popup `Section/BusinessAndSOBListEDM`,
// `Activity/EDMChooseBusiness_Act` (kelas Int-treaty_in_edm, ruleset 01-01-83; 14 langkah aktif, nol `//`) beserta
// aktivitas yang dipanggilnya (korpus `EDM Treaty In\Activity`, `DataTransform\SetInstallmentValue`):
//
//	EDMChooseBusiness_Act 1   SetValueEDM_Act                       -> SetValueEDM (+ InputPolicyTreatyEDMDetail_NP*)
//	                      2   SetValueOldTax (FlagPPH=="true")      -> SetValueOldTax
//	                      3   TreatyRealizationCheckXOLListEDM      -> TreatyRealizationCheckXOLListEDM
//	                      4   CopyGeneralDataEDM_act                -> CopyGeneralDataEDM
//	                      5   SetEDMTCancel (EDMType=='4')          -> SetEDMTCancel
//	                      6   CalculateDifferenceEDM_act            -> CalculateDifferenceEDM (edm_xol_selisih.go)
//	                      7   FillMasterInstallment                 -> FillMasterInstallment (+ SetInstallmentValue)
//	                      8   DT ExpandAllExpandables               -> ⛔ keadaan layar (pyExpanded), tidak diport
//	                      9   FillSpreading                         -> FillSpreading
//	                      10  OfferFacIn.QuotationData; QuotationData.OldPolicyNo
//	                      11  Installment = jumlah TreatyIn.Installment(1).InstallmentList
//	                      12  FillPaymentInstallmentEDMT            -> FillPaymentInstallmentEDMT (edm_angsuran.go)
//	                      13-14 Obj-Save, Commit                    -> services (penyimpanan)
//
// DUA halaman `TreatyIn`: `pyWorkPage.TreatyIn` (jalur `TreatyIn.` di halaman kerja `h`, master BARU hasil
// SetValueEDM_Act 8) dan halaman clipboard TINGKAT ATAS `TreatyIn` (kelas Int-TREATY_IN; Pages & Classes setiap
// aktivitas di atas) yang diisi SetValueEDM_Act 2-4, SetTreatyIn_Act (master LAMA) dan SetTreatyInEDM_Act, lalu
// dibuang FillMasterInstallment 5. Yang kedua dibawa sebagai halaman terpisah `atas` (kunci juga `TreatyIn.`);
// tidak pernah disimpan.
//
// Master dibaca lewat `PembacaMasterEDM` (RDB-List + Java `adoptJSONObject`) - fungsi di sini murni.

import (
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// PembacaMasterEDM - pembacaan master kontrak treaty (JSONDATA, baca-saja, pengecualian K8) yang dipanggil rantai
// pilih bisnis EDM. Setiap metode = SATU RDB-List; hasilnya SATU `MasterXOL` per baris hasil, BERURUTAN seperti
// `pxResults` (Java langkah sesudahnya meng-`adoptJSONObject` setiap baris berurutan - `AdopsiMasterEDM`).
// Nol baris = irisan kosong TANPA galat (Pega: halaman tidak berubah, tidak ada pesan). Galat hanya untuk kegagalan
// baca / dokumen rusak.
type PembacaMasterEDM interface {
	// MasterMenurutOldID = RDB `BrowseTreatyInEDM` kelas Data-PolicyTreatyIn (SetValueEDM_Act 5):
	// `select JSONDATA as CLASSOFBUSINESS from pooldata.m_treaty_in_edm where OLDID = {OldData.NoOffer}`.
	MasterMenurutOldID(oldID string) ([]MasterXOL, error)
	// MasterMenurutID = RDB `BrowseTreatyIn` (SetValueEDM_Act 6, SetTreatyIn_Act 4): M_TREATY_IN WHERE ID = x
	// UNION ALL M_TREATY_IN_edm WHERE ID = x.
	MasterMenurutID(id string) ([]MasterXOL, error)
	// MasterEDMMenurutID = RDB `BrowseTreatyInEDM` kelas Int-treaty_in_edm (SetTreatyInEDM_Act 3):
	// `select JSONDATA as ClassofBusiness from pooldata.M_TREATY_IN_EDM where ID={TreatyIn.ID}`.
	MasterEDMMenurutID(id string) ([]MasterXOL, error)
	// MasterOutMenurutID = RDB `BrowseTreatyOutDetailEDM` (SetValueEDM_Act 7, SetTreatyInEDM_Act 4; XOL Retro):
	// `select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={OldData.NoOffer}`.
	MasterOutMenurutID(id string) ([]MasterXOL, error)
	// IDMataUang = RDB `GetDataCurrencyByName_SQL`: `select ID, ... from CURRENCY where CURRENCY = {nama}` ->
	// `CurrencySearch.pxResults(1).ID`; nama tak dikenal = "" tanpa galat.
	IDMataUang(nama string) (string, error)
}

// PesanSpreadingAsalKosong - VERBATIM FillSpreading langkah 2 (`local.errsprdempty`).
const PesanSpreadingAsalKosong = "Spreading for original policy (No Polis) is empty."

// PilihBisnisEDM = `Activity/EDMChooseBusiness_Act` langkah 1-12 atas halaman kerja `h`. `sekarang` menggantikan
// `@CurrentDateTime()` (CopyGeneralDataEDM_act 1).
//
// `[penyesuaian]` Langkah 11-12 berprakondisi `pyWorkPage.pyWorkIDPrefix=="EDMT-"`: setiap kasus modul ini lahir
// ber-awalan `AwalanKasus` (kasus.go `RakitIDKasus`), jadi kedua langkah selalu jalan.
// ⚠️ Halaman `pyWorkPage.TreatyIn` di Pega tersimpan di blob kasus, sehingga `adoptJSONObject` pilihan kedua
// menimpa di ATAS master pilihan pertama. Di sini isinya = apa pun yang ada di `h` saat dipanggil (halaman master
// tidak disimpan - katalog hanya `TreatyIn.ID`).
func PilihBisnisEDM(h *Halaman, baca PembacaMasterEDM, sekarang time.Time) error {
	atas := HalamanBaru() // clipboard tingkat atas: halaman `TreatyIn` (Int-TREATY_IN)
	// keputusan WO 07-10-2026: With Tax / Type Tax kosong diwarisi generasi lama SEBELUM rantai (kasus lama ikut)
	WarisPajakLama(h)
	if err := SetValueEDM(h, atas, baca); err != nil { // 1
		return err
	}
	if h.Ambil(pt+"FlagPPH") == "true" { // 2
		SetValueOldTax(h)
	}
	if err := TreatyRealizationCheckXOLListEDM(h, atas, baca); err != nil { // 3
		return err
	}
	CopyGeneralDataEDM(h, sekarang)   // 4
	if h.Ambil(pt+"EDMType") == "4" { // 5 PRE `EDMType=='4'`
		SetEDMTCancel(h)
	}
	if err := CalculateDifferenceEDM(h); err != nil { // 6
		return err
	}
	if err := FillMasterInstallment(h, atas, baca); err != nil { // 7
		return err
	}
	// 8 DT ExpandAllExpandables - keadaan layar, tidak diport
	if err := FillSpreading(h); err != nil { // 9
		return err
	}
	// 10 OfferFacIn.QuotationData = PolicyTreatyIn.QuotationData (DULU), lalu QuotationData.OldPolicyNo
	pbSalinHalaman(h, pt+"QuotationData", "OfferFacIn.QuotationData")
	h.Setel(pt+"QuotationData.OldPolicyNo", h.Ambil(od+"PolicyNo"))
	// 11 Local.Installment (String) = @LengthOfPageList(pyWorkPage.TreatyIn.Installment(1).InstallmentList)
	installment := strconv.Itoa(len(h.AmbilDaftar(JalurAnak(jMaster+"Installment", 1, AnakRinciAngsuran))))
	h.Setel(pt+"Installment", installment)
	return FillPaymentInstallmentEDMT(h, installment) // 12 (Param.Installment = Local.Installment)
}

// HitungUlangPajakNonPropEDM - ⛔ keputusan work owner 07-10-2026 ("UNTUK EDMT NON PROP KENAPA INI TU KOSONG??", pola
// NB `services.pajakNonProp` 06-10-2026): With Tax / Type Tax berubah SESUDAH Choose Business pada polis NonProp baru
// -> rantai `PilihBisnisEDM` dijalankan lagi, satu-satunya tempat XML menghitung pajak grid XOL
// (`InsertToTreatyXOLListEDM` 3.3.7 / 3.3.8.5) dan selisihnya (`CalculateDifferenceEDM` butuh halaman master
// IsProRate yang tidak disimpan). Sel `.TypeTax` XML hanya postValue, `.FlagPPH` hanya RemoveTypeTax_ACT.
//
//   - belum pernah Choose Business (TreatyXOLList kosong) -> tidak ada yang dihitung;
//   - master SAMA dengan pilihan terakhir: NoOffer = OldData.NoOffer berarti pilihan terakhir lewat jalur pertama
//     (SetValueEDM 4-6, master lama ber-ID) - NoOffer dikosongkan supaya jalurnya sama, bukan maju ke master EDM;
//   - isian yang ikut ditulis ulang rantai dikembalikan: `.Installment` (sel S12), StatementDate (DueDate angsuran),
//     MarketingOfficer - lalu angsuran disusun ulang (`FillPaymentInstallmentEDMT`, Param.Installment = .Installment).
func HitungUlangPajakNonPropEDM(h *Halaman, baca PembacaMasterEDM, sekarang time.Time) error {
	if !PolisNonPropBaru(h) || len(h.AmbilDaftar(DaftarXOL)) == 0 {
		return nil
	}
	simpan := map[string]string{}
	for _, m := range []string{"Installment", "StatementDate", "MarketingOfficer"} {
		simpan[m] = h.Ambil(pt + m)
	}
	if h.Ambil(pt+"NoOffer") == h.Ambil(od+"NoOffer") {
		h.Setel(pt+"NoOffer", "")
	}
	if err := PilihBisnisEDM(h, baca, sekarang); err != nil {
		return err
	}
	for m, v := range simpan {
		h.Setel(pt+m, v)
	}
	return FillPaymentInstallmentEDMT(h, simpan["Installment"])
}

// SetValueEDM = `Activity/SetValueEDM_Act` (kelas Int-treaty_in_edm, 01-01-83; 12 langkah aktif, nol `//`).
//
//	1    Page-New TempResult
//	2-3  Page-Remove / Page-New `TreatyIn` (tingkat atas)
//	4    PRE PolicyTreatyIn.NoOffer=="" -> TreatyIn.ID = OldData.NoOffer
//	5    PRE bukan XOL Retro -> RDB BrowseTreatyInEDM (OLDID = OldData.NoOffer)       `MasterMenurutOldID`
//	6    PRE NoOffer=="" -> RDB BrowseTreatyIn (ID = TreatyIn.ID tingkat atas)          `MasterMenurutID`
//	7    PRE XOL Retro -> RDB BrowseTreatyOutDetailEDM (M_TREATY_OUT ID = OldData.NoOffer) `MasterOutMenurutID`
//	8    Java: setiap TempResult.pxResults -> pyWorkPage.TreatyIn.adoptJSONObject(CLASSOFBUSINESS)
//	9    PRE EDMType=="3" benar -> lewati; selain itu InputPolicyTreatyEDMDetail_NP
//	10   PRE EDMType=="3" -> InputPolicyTreatyEDMDetail_NP_AdjPremi
//	11-12 Obj-Save, Commit - services
//
// ⚠️ `[dugaan]` Ketiga RDB-List menulis ke `TempResult` yang SAMA dan RDB-List MENGGANTI `pxResults`: hanya hasil
// RDB terakhir yang berjalan yang di-adopt. Akibatnya (ditiru): pilihan pertama (NoOffer kosong, bukan Retro) hasil
// langkah 5 (master EDM ber-OLDID) ditimpa langkah 6 (master lama ber-ID); XOL Retro selalu master M_TREATY_OUT.
// Baris popup yang dipilih TIDAK dibaca langkah mana pun - master ditentukan `OldData.NoOffer`.
// ⚠️ Langkah 9 tidak memeriksa jenis proporsi: polis Proporsional (EDMType != "3") ikut menjalankan
// InputPolicyTreatyEDMDetail_NP (medan uang lalu ditimpa CopyGeneralDataEDM_act 2) - ditiru apa adanya.
func SetValueEDM(h, atas *Halaman, baca PembacaMasterEDM) error {
	var hasil []MasterXOL                        // 1
	pbKosongkan(atas, HalamanMaster)             // 2-3
	noOfferKosong := h.Ambil(pt+"NoOffer") == "" // PRE 4, 6
	retro := h.Ambil(pt+"ClaimType") == KlaimXOLRetro
	if noOfferKosong { // 4
		atas.Setel(jMaster+"ID", h.Ambil(od+"NoOffer"))
	}
	var err error
	if !retro { // 5
		if hasil, err = baca.MasterMenurutOldID(h.Ambil(od + "NoOffer")); err != nil {
			return err
		}
	}
	if noOfferKosong { // 6
		if hasil, err = baca.MasterMenurutID(atas.Ambil(jMaster + "ID")); err != nil {
			return err
		}
	}
	if retro { // 7
		if hasil, err = baca.MasterOutMenurutID(h.Ambil(od + "NoOffer")); err != nil {
			return err
		}
	}
	for _, m := range hasil { // 8
		AdopsiMasterEDM(h, HalamanMaster, m)
	}
	if h.Ambil(pt+"EDMType") != "3" { // 9
		return InputPolicyTreatyEDMDetailNP(h, baca)
	}
	return InputPolicyTreatyEDMDetailNPAdjPremi(h, baca) // 10
}

// AdopsiMasterEDM meniru `ClipboardPage.adoptJSONObject` satu baris master ke halaman `halaman` (`TreatyIn`) di `h`:
// skalar dokumen menimpa; daftar dokumen MENGGANTI daftar halaman beserta daftar bersarangnya; medan yang tidak ada
// di dokumen tetap. `[dugaan]` sama dengan tafsiran pengurai `repository.uraiMasterXOL` (baris berikut menimpa kunci
// tingkat atas, daftar diganti utuh); halaman tersemat (`ValueDifference`) digabung per daftar.
func AdopsiMasterEDM(h *Halaman, halaman string, m MasterXOL) {
	awal := halaman + "."
	for k, v := range m.Nilai {
		h.Setel(awal+k, v)
	}
	for k := range m.Daftar {
		if !strings.Contains(k, "(") {
			hapusDaftarBeserta(h, awal+k)
		}
	}
	for k, v := range m.Daftar {
		h.SetelDaftar(awal+k, salinBaris(v))
	}
}

// InputPolicyTreatyEDMDetailNP = `Activity/InputPolicyTreatyEDMDetail_NP` (Work, 01-01-83; 13 aktif, 3 `//`).
// Jumlah atas `pyWorkPage.TreatyIn.Share`; lalu 8 bukan XOL Retro -> `InsertToTreatyXOLList` (port NB,
// nonprop.go - identik di korpus EDM), 9 XOL Retro -> `InsertToTreatyOutXOLList`.
func InputPolicyTreatyEDMDetailNP(h *Halaman, baca PembacaMasterEDM) error {
	if err := pbInputDetailNP(h, "Share"); err != nil {
		return err
	}
	idMU, err := pbPetaMataUang(baca, h)
	if err != nil {
		return err
	}
	if h.Ambil(pt+"ClaimType") != KlaimXOLRetro { // 8
		return InsertToTreatyXOLList(h, idMU)
	}
	return InsertToTreatyOutXOLList(h, idMU) // 9
}

// InputPolicyTreatyEDMDetailNPAdjPremi = `Activity/InputPolicyTreatyEDMDetail_NP_AdjPremi` (Work, 01-01-56; 12
// aktif, 3 `//`). Sama dengan `_NP` tetapi jumlah atas `pyWorkPage.TreatyIn.ValueDifference.Share`, dan 8
// `InsertToTreatyXOLListEDM` TANPA syarat (tidak ada cabang XOL Retro).
func InputPolicyTreatyEDMDetailNPAdjPremi(h *Halaman, baca PembacaMasterEDM) error {
	if err := pbInputDetailNP(h, "ValueDifference.Share"); err != nil {
		return err
	}
	idMU, err := pbPetaMataUang(baca, h)
	if err != nil {
		return err
	}
	return InsertToTreatyXOLListEDM(h, idMU) // 8
}

// pbInputDetailNP - langkah 1-7 kedua aktivitas (`daftar` = `Share` / `ValueDifference.Share`).
//
//	1    `//` pyExpanded - tidak diport
//	2    DueTo = "1"; StartDate = TreatyIn.Commencement; EndDate = TreatyIn.Termination
//	3    PRE TreatyIn.FacultativeShare>0 -> pyWorkPage.OfferFacIn.IsFacRetro = "1"
//	4    local.currency = PolicyTreatyIn.Currency
//	5    per baris daftar: 5.1-5.3 jumlah .Value GrossPremiumList / DeductionTotalList / NetPremiumList ber-Currency
//	     = local.currency; 5.4 RnmLimitList `//` - tidak diport
//	6    PremiOgp = gross; ShareValue = local.share; Deduction1 = deduction; Deduction2 = 0; NetPremium =
//	     BalanceDueTo = net; ListInstallment = "" (dikosongkan)
//	7    `//` susun installment (beserta 7 anak) - tidak diport
//
// ⚠️ `local.share` (Decimal) hanya diisi 5.4 yang `//`: ShareValue = 0 (nilai awal lokal Decimal, konvensi nonprop.go).
func pbInputDetailNP(h *Halaman, daftar string) error {
	k := &kalkulator{}
	p := polis{h, k}
	h.Setel(pt+"DueTo", "1") // 2
	h.Setel(pt+"StartDate", h.Ambil(jMaster+"Commencement"))
	h.Setel(pt+"EndDate", h.Ambil(jMaster+"Termination"))
	if angkaMaster(k, h, "FacultativeShare").Sign() > 0 { // 3
		h.Setel("OfferFacIn.IsFacRetro", "1")
	}
	mu := p.teks("Currency") // 4
	gross, ded, net, share := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
	for i := range mDaftar(h, daftar) { // 5
		gross = jumlahMataUang(k, mAnak(h, daftar, i+1, "GrossPremiumList"), mu, gross) // 5.1
		ded = jumlahMataUang(k, mAnak(h, daftar, i+1, "DeductionTotalList"), mu, ded)   // 5.2
		net = jumlahMataUang(k, mAnak(h, daftar, i+1, "NetPremiumList"), mu, net)       // 5.3
	}
	if k.err != nil {
		return k.err
	}
	p.setel("PremiOgp", gross) // 6
	p.setel("ShareValue", share)
	p.setel("Deduction1", ded)
	p.setel("Deduction2", apd.New(0, 0))
	p.setel("NetPremium", net)
	p.setel("BalanceDueTo", net)
	hapusDaftarBeserta(h, DaftarAngsuran)
	return nil
}

// pbPetaMataUang - ID mata uang (`IDMataUang`) setiap `TreatyIn.Installment().Currency` halaman `src` (RDB
// `GetDataCurrencyByName_SQL` yang dijalankan per baris angsuran oleh aktivitas XOL).
func pbPetaMataUang(baca PembacaMasterEDM, src *Halaman) (IDMataUang, error) {
	out := IDMataUang{}
	for _, mu := range MataUangAngsuranMaster(src) {
		id, err := baca.IDMataUang(mu)
		if err != nil {
			return nil, err
		}
		out[mu] = id
	}
	return out, nil
}

// SetValueOldTax = `Activity/SetValueOldTax` (Int-treaty_in_edm, 01-01-76; 4 aktif): pajak lama yang belum ada.
//
//	1    per OldData.TreatyXOLList: 1.1 PRE .PPNValue=="" -> PPHValue = PPNValue = 0; NetPremiAfterPPH/PPN/Tax =
//	     .NetPremi (`tes.CARI1 = .NetPremiAfterPPH` - nol pembaca, tidak diport)
//	1.2  per .ValueList: PRE .PPNValue=="" -> sama (tanpa tes.CARI1)
func SetValueOldTax(h *Halaman) {
	isi := func(b Baris) {
		if b["PPNValue"] != "" {
			return
		}
		b["PPHValue"], b["PPNValue"] = "0", "0"
		b["NetPremiAfterPPH"], b["NetPremiAfterPPN"], b["NetPremiAfterTax"] = b["NetPremi"], b["NetPremi"], b["NetPremi"]
	}
	for c, x := range h.AmbilDaftar(od + "TreatyXOLList") {
		isi(x) // 1.1
		for _, v := range h.AmbilDaftar(JalurAnak(od+"TreatyXOLList", c+1, AnakLayerXOL)) {
			isi(v) // 1.2.1
		}
	}
}

// TreatyRealizationCheckXOLListEDM = `Activity/TreatyRealizationCheckXOLListEDM` (Work, 01-01-83; 7 aktif).
//
//	1    Page-Clear-Messages pyWorkPage
//	2    PRE FlagRetroTreaty==true -> lompat AA (langkah 4)
//	3    PRE @SizeOfPropertyList(OldData.TreatyXOLList)<1 salah -> keluar aktivitas
//	4 AA Param.ID = OldData.NoOffer (`local.err` diisi, nol pembaca)
//	5    Call SetTreatyIn_Act (pass-current-param-page: Param.ID langkah 4; ID=3000004 / viewstate=1 di larik
//	     parameter diabaikan) - `pbSetTreatyIn`, master LAMA ke `TreatyIn` tingkat atas
//	6    bukan XOL Retro -> InsertToTreatyXOLListEDMOldData
//	7    XOL Retro -> InsertToTreatyOutXOLListEDMOldData
func TreatyRealizationCheckXOLListEDM(h, atas *Halaman, baca PembacaMasterEDM) error {
	h.BersihkanPesan()                           // 1
	if h.Ambil(pt+"FlagRetroTreaty") != "true" { // 2
		if len(h.AmbilDaftar(od+"TreatyXOLList")) >= 1 { // 3
			return nil
		}
	}
	id := h.Ambil(od + "NoOffer")                         // 4
	if err := pbSetTreatyIn(atas, baca, id); err != nil { // 5
		return err
	}
	if h.Ambil(pt+"ClaimType") != KlaimXOLRetro { // 6
		idMU, err := pbPetaMataUang(baca, atas)
		if err != nil {
			return err
		}
		return InsertToTreatyXOLListEDMOldData(h, atas, idMU)
	}
	idMU, err := pbPetaMataUang(baca, h) // 7 (sumbernya pyWorkPage.TreatyIn)
	if err != nil {
		return err
	}
	return InsertToTreatyOutXOLListEDMOldData(h, idMU)
}

// pbSetTreatyIn = `Activity/SetTreatyIn_Act` (Data-Portal; 15 langkah, identik korpus NB) atas `TreatyIn` tingkat
// atas, Param.ID = `id`:
//
//	3    TreatyIn.ID = Param.ID
//	4    RDB BrowseTreatyIn (M_TREATY_IN ∪ M_TREATY_IN_EDM WHERE ID = TreatyIn.ID)   `MasterMenurutID`
//	5    Java: adoptJSONObject setiap baris ke `TreatyIn`
//
// ⛔ Tidak diport, nol pembaca di rantai ini: 1 & 15 FlagExcel; 2 TreatyInInputVis (Add=1: TreatyIn.ID = "UnknownId"
// lalu keluar - ditimpa 3); 6 (Param.viewstate tidak ada di halaman parameter yang diteruskan); 7-11
// (revisionstate - SaveTreatyIn tidak pernah dari sini); 12 ConvertHistoryDate (TreatyIn.CommentList); 13
// TreatySetReinstatement (TreatyIn.Limits tingkat atas); 14 CheckDuplicateOffer (RDB GetCountClaim, pesan di
// halaman TreatyIn tingkat atas yang tidak disimpan maupun ditampilkan).
func pbSetTreatyIn(atas *Halaman, baca PembacaMasterEDM, id string) error {
	atas.Setel(jMaster+"ID", id) // 3
	hasil, err := baca.MasterMenurutID(id)
	if err != nil {
		return err
	}
	for _, m := range hasil { // 5
		AdopsiMasterEDM(atas, HalamanMaster, m)
	}
	return nil
}

// CopyGeneralDataEDM = `Activity/CopyGeneralDataEDM_act` (Int-treaty_in_edm, 01-01-77; 2 aktif; step page
// `PolicyTreatyIn.OldData`).
//
// ⚠️ Ditiru apa adanya: `IDCurrency = .Currency` (NAMA mata uang, menimpa ID CreateEDMT 14); TreatyGroupName diisi
// dua kali - `TreatyIn.Limits(1).TreatyGroupList(1).TreatyGroup` lalu ditimpa `OldData.TreatyGroupName`.
func CopyGeneralDataEDM(h *Halaman, sekarang time.Time) {
	old := func(m string) string { return h.Ambil(od + m) }
	set := func(m, v string) { h.Setel(pt+m, v) }
	// 1
	for _, m := range []string{"BizCode", "BizName", "CedingCo", "CedingCoName", "Currency", "EndDate"} {
		set(m, old(m))
	}
	set("IDCurrency", old("Currency"))
	for _, m := range []string{"Installment", "InsuredName", "InsuredID", "IsNewPolicyNonProp", "Layer", "LayerPart",
		"LayerPartType", "LayerType", "MarketingOfficer"} {
		set(m, old(m))
	}
	set("NoOffer", h.Ambil(jMaster+"ID"))
	for _, m := range []string{"Show", "SOB", "SOBName", "StartDate"} {
		set(m, old(m))
	}
	set("StatementDate", utils.FormatTanggalWaktu(sekarang)) // @CurrentDateTime()
	set("TreatyGroupName", barisKe(h.AmbilDaftar(JalurAnak(jMaster+"Limits", 1, "TreatyGroupList")), 0)["TreatyGroup"])
	for _, m := range []string{"TreatyYear", "TreatyType", "TreatyGroupID", "TreatyGroupName", "Quartal", "YearOfQuartal"} {
		set(m, old(m))
	}
	// 2 PRE OldData.IsNewPolicyNonProp != 1
	if samaDenganSatu(old("IsNewPolicyNonProp")) {
		return
	}
	for _, m := range []string{"Installment", "GrossPremium", "PremiOgp", "RiCommOgp", "ResultOgp1", "OveriddingCommOgp",
		"ResultOgp2", "Claim", "SalvageValue", "ExcessLoss", "NetPremium", "BalanceDueTo", "PremiOnp", "RiCommOnp",
		"ResultOnp1", "OveriddingCommOnp", "ResultOnp2", "Deduction1", "Deduction2"} {
		set(m, old(m))
	}
	// OldData.BalanceBeforePPH / BalanceBeforeTax = @if(.X==""||.X==0, .NetPremium, .X)
	for _, m := range []string{"BalanceBeforePPH", "BalanceBeforeTax"} {
		if pbKosongAtauNol(old(m)) {
			h.Setel(od+m, old("NetPremium"))
		}
	}
}

// SetEDMTCancel = `Activity/SetEDMTCancel` (Work, 01-01-77; 10 aktif) - EDM pembatalan (EDMType 4).
//
//	1    PRE .PolicyTreatyIn.IsNewPolicyNonProp==1 benar -> lewati (blok Proporsional, REPEAT sekali)
//	1.1  "0": PremiOgp, RiCommOgp, ResultOgp1, OveriddingCommOgp, ResultOgp2, PremiOnp, RiCommOnp, ResultOnp1,
//	     OveriddingCommOnp, ResultOnp2, Claim, OutstandingClaim, SalvageValue, ExcessLoss, Deduction1, Deduction2
//	1.2  Installment, ListInstallment, SpreadingRiskList, TreatyXOLList = OldData.* (salinan daftar utuh)
//	1.3  ListInstallment(): PaymentTotal = Premium = 0;  1.4 SpreadingRiskList(): PremiumSpreaded = SharePercentage = 0
//	2    PRE IsNewPolicyNonProp==1 (NonProp): per TreatyXOLList (c): Currency/IDCurrency/DueTo = OldData.TreatyXOLList(c);
//	     Deduction, DueToValue, GrossPremi, NetPremi, PPH, PPN, NetPremiAfter* = 0; per ValueList (l): Currency =
//	     Currency induk, IDCurrency = .IDCurrency (tetap), DueTo = DueTo induk, Layer* = OldData...ValueList(l), uang = 0
//
// ⚠️ Ditiru apa adanya: BrokerageFeeSebenarnya tidak dinolkan (2); GrossPremium, NetPremium, BalanceDueTo, PPN/PPH,
// BalanceBefore* tidak dinolkan (1.1).
func SetEDMTCancel(h *Halaman) {
	if !samaDenganSatu(h.Ambil(pt + "IsNewPolicyNonProp")) { // 1
		for _, m := range []string{"PremiOgp", "RiCommOgp", "ResultOgp1", "OveriddingCommOgp", "ResultOgp2", "PremiOnp",
			"RiCommOnp", "ResultOnp1", "OveriddingCommOnp", "ResultOnp2", "Claim", "OutstandingClaim", "SalvageValue",
			"ExcessLoss", "Deduction1", "Deduction2"} { // 1.1
			h.Setel(pt+m, "0")
		}
		h.Setel(pt+"Installment", h.Ambil(od+"Installment")) // 1.2
		pbSalinDaftar(h, h, od+"ListInstallment", DaftarAngsuran)
		pbSalinDaftar(h, h, od+"SpreadingRiskList", DaftarSpreading)
		pbSalinDaftar(h, h, od+"TreatyXOLList", DaftarXOL)
		for _, b := range h.AmbilDaftar(DaftarAngsuran) { // 1.3
			b["PaymentTotal"], b["Premium"] = "0", "0"
		}
		for _, b := range h.AmbilDaftar(DaftarSpreading) { // 1.4
			b["PremiumSpreaded"], b["SharePercentage"] = "0", "0"
		}
		return
	}
	nolUang := []string{"Deduction", "DueToValue", "GrossPremi", "NetPremi", "PPHValue", "PPNValue",
		"NetPremiAfterPPH", "NetPremiAfterPPN", "NetPremiAfterTax"}
	lamaInduk := h.AmbilDaftar(od + "TreatyXOLList")
	for c, x := range h.AmbilDaftar(DaftarXOL) { // 2.1
		o := barisKe(lamaInduk, c)
		x["Currency"], x["IDCurrency"], x["DueTo"] = o["Currency"], o["IDCurrency"], o["DueTo"]
		for _, m := range nolUang {
			x[m] = "0"
		}
		lamaLapis := h.AmbilDaftar(JalurAnak(od+"TreatyXOLList", c+1, AnakLayerXOL))
		for l, v := range h.AmbilDaftar(JalurAnak(DaftarXOL, c+1, AnakLayerXOL)) { // 2.1.1
			ol := barisKe(lamaLapis, l)
			v["Currency"], v["DueTo"] = x["Currency"], x["DueTo"] // local.XOLCurrency / XOLDueTo
			for _, m := range []string{"LayerType", "Layer", "LayerPartType", "LayerPart"} {
				v[m] = ol[m]
			}
			for _, m := range nolUang {
				v[m] = "0"
			}
		}
	}
}

// FillMasterInstallment = `Activity/FillMasterInstallment` (Int-treaty_in_edm, 01-01-64; 5 aktif).
//
//	1    Page-New TreatyIn (tingkat atas)
//	2    Param.ID = PolicyTreatyIn.NoOffer
//	3    Call SetTreatyInEDM_Act (pass-current-param-page; Data-Portal, 01-01-90):
//	       1 TreatyInInputVis (Add=1: TreatyIn.ID "UnknownId" lalu keluar - ditimpa 2); 2 TreatyIn.ID = Param.ID;
//	       3 RDB BrowseTreatyInEDM (M_TREATY_IN_EDM ID = TreatyIn.ID)            `MasterEDMMenurutID`
//	       4 PRE XOL Retro -> RDB BrowseTreatyOutDetailEDM (ID = OldData.NoOffer) `MasterOutMenurutID` (mengganti)
//	       5 Java adoptJSONObject ke TreatyIn;  6 viewstate - Param.viewstate tidak ada, tidak jalan
//	4    DataTransform SetInstallmentValue
//	5    Page-Remove TreatyIn (tingkat atas)
func FillMasterInstallment(h, atas *Halaman, baca PembacaMasterEDM) error {
	pbKosongkan(atas, HalamanMaster) // 1
	id := h.Ambil(pt + "NoOffer")    // 2
	atas.Setel(jMaster+"ID", id)     // 3.2
	hasil, err := baca.MasterEDMMenurutID(id)
	if err != nil {
		return err
	}
	if h.Ambil(pt+"ClaimType") == KlaimXOLRetro { // 3.4
		if hasil, err = baca.MasterOutMenurutID(h.Ambil(od + "NoOffer")); err != nil {
			return err
		}
	}
	for _, m := range hasil { // 3.5
		AdopsiMasterEDM(atas, HalamanMaster, m)
	}
	SetInstallmentValue(h, atas)     // 4
	pbKosongkan(atas, HalamanMaster) // 5
	return nil
}

// SetInstallmentValue = `DataTransform/SetInstallmentValue` (kelas Int-treaty_in_edm, nol baris nonaktif):
//
//	SET  PolicyTreatyIn.Installment = pyWorkPage.TreatyIn.InstallmentNo
//	FOR EACH TreatyIn.ValueDifference.Installment (tingkat atas, subskrip n):
//	  ListInstallment(n).Currency = .Currency; ListInstallment(n).Premium = .AmountTotal;
//	  TreatyDifference.TotalPremium = .AmountTotal; Param.Installment = n
//	  FOR EACH .InstallmentList (m): ListInstallment(n).InstallmentList(m).Currency = .Currency, .Premium = .Amount,
//	    .DueDate = .DueDate, .InstallmentNo = .Installment, .InstallmentPercentage = .InstallmentPct,
//	    .PaymentDate = .PaymentDate
//
// Ditulis per subskrip (anggota lain baris yang sudah ada tetap). Di rantai Choose `ListInstallment` lalu dibangun
// ulang FillPaymentInstallmentEDMT (EDMChooseBusiness_Act 12) - ditiru apa adanya.
func SetInstallmentValue(h, atas *Halaman) {
	h.Setel(pt+"Installment", h.Ambil(jMaster+"InstallmentNo"))
	induk := jMaster + "ValueDifference.Installment"
	for n, inst := range atas.AmbilDaftar(induk) {
		b := xePastikanBaris(h, DaftarAngsuran, n+1)
		b["Currency"], b["Premium"] = inst["Currency"], inst["AmountTotal"]
		h.Setel(sd+"TotalPremium", inst["AmountTotal"])
		rinci := JalurAnak(DaftarAngsuran, n+1, AnakRinciAngsuran)
		for m, il := range atas.AmbilDaftar(JalurAnak(induk, n+1, AnakRinciAngsuran)) {
			r := xePastikanBaris(h, rinci, m+1)
			r["Currency"], r["Premium"], r["DueDate"] = il["Currency"], il["Amount"], il["DueDate"]
			r["InstallmentNo"], r["InstallmentPercentage"], r["PaymentDate"] = il["Installment"], il["InstallmentPct"], il["PaymentDate"]
		}
	}
}

// FillSpreading = `Activity/FillSpreading` (Int-treaty_in_edm, 01-01-90; 7 aktif).
//
//	1    Page-Clear-Messages halaman PRIMER (baris popup Int-treaty_in_edm, bukan pyWorkPage) - tidak diport
//	2    local.errsprdempty = PesanSpreadingAsalKosong
//	3    SpreadingRiskList(1) = OldData.SpreadingRiskList(1) (salinan halaman); .TreatyType = Old(1).TreatyType;
//	     .PremiumSpreaded = TreatyXOLDifferenceList(1).NetPremi
//	4    jumlah ClaimPercentage, ClaimSpreaded, SharePercentage, PremiumSpreaded setiap baris
//	5    TotalSharePercentagePremium / TotalPremium / TotalSharePercentageClaim / TotalClaim
//	6    PRE OldData.SpreadingRiskList(1).TreatyType = "" -> Page-Set-Messages pyWorkPage
//
// `[dugaan]` Salinan halaman 3 MENGGANTI isi baris 1 (baris lain tetap); baris lama yang tidak ada = halaman kosong.
func FillSpreading(h *Halaman) error {
	k := &kalkulator{}
	p := polis{h, k}
	lama1 := barisKe(h.AmbilDaftar(od+"SpreadingRiskList"), 0)
	d := h.AmbilDaftar(DaftarSpreading)
	if len(d) == 0 {
		d = []Baris{{}}
	}
	baru := Baris{} // 3
	for m, v := range lama1 {
		baru[m] = v
	}
	baru["TreatyType"] = lama1["TreatyType"]
	baru["PremiumSpreaded"] = barisKe(h.AmbilDaftar(DaftarSelisihXOL), 0)["NetPremi"]
	d[0] = baru
	// ⚠️ Keputusan work owner 07-10-2026 (grid spreading hanya-baca, tombol Add dibuang - sama dengan NB): baris generasi
	// lama ke-2 dst. yang belum ada di data baru DIBAWA apa adanya, berurutan (di Pega admin menambahkannya lewat Add).
	// Tanpa ini generasi berbaris lebih dari satu tidak pernah lolos keutuhan ID-15 (`BarisSpreadingHilang`).
	lama := h.AmbilDaftar(od + "SpreadingRiskList")
	for i := len(d); i < len(lama); i++ {
		b := Baris{}
		for m, v := range lama[i] {
			b[m] = v
		}
		d = append(d, b)
	}
	h.SetelDaftar(DaftarSpreading, d)
	shareKlaim, klaim, sharePremi, premi := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
	for _, b := range d { // 4
		shareKlaim = k.Tambah(angkaBaris(k, b, "ClaimPercentage"), shareKlaim)
		klaim = k.Tambah(angkaBaris(k, b, "ClaimSpreaded"), klaim)
		sharePremi = k.Tambah(angkaBaris(k, b, "SharePercentage"), sharePremi)
		premi = k.Tambah(angkaBaris(k, b, "PremiumSpreaded"), premi)
	}
	p.setel("TotalSharePercentagePremium", sharePremi) // 5
	p.setel("TotalPremium", premi)
	p.setel("TotalSharePercentageClaim", shareKlaim)
	p.setel("TotalClaim", klaim)
	if k.err != nil {
		return k.err
	}
	if lama1["TreatyType"] == "" { // 6
		h.TambahPesan("", PesanSpreadingAsalKosong)
	}
	return nil
}

// ---------------------------------------------------------------- pembantu halaman

// pbKosongkan meniru `Page-New` / `Page-Remove` halaman `awalan` di `h`: setiap nilai dan daftar berawalan
// `awalan.` dibuang.
func pbKosongkan(h *Halaman, awalan string) {
	h.pastikan()
	a := awalan + "."
	for k := range h.Nilai {
		if strings.HasPrefix(k, a) {
			delete(h.Nilai, k)
		}
	}
	for k := range h.Daftar {
		if strings.HasPrefix(k, a) {
			delete(h.Daftar, k)
		}
	}
}

// pbSalinDaftar meniru `Property-Set <daftar> = <daftar lain>` (PageList): daftar `ke` di `dst` diganti salinan
// daftar `dari` di `src` beserta daftar bersarangnya (`dari(n).anak` -> `ke(n).anak`).
func pbSalinDaftar(src, dst *Halaman, dari, ke string) {
	hapusDaftarBeserta(dst, ke)
	isi := salinBaris(src.AmbilDaftar(dari))
	bersarang := map[string][]Baris{}
	for k, v := range src.Daftar {
		if strings.HasPrefix(k, dari+"(") {
			bersarang[ke+k[len(dari):]] = salinBaris(v)
		}
	}
	dst.SetelDaftar(ke, buangKosong(isi))
	for k, v := range bersarang {
		dst.SetelDaftar(k, v)
	}
}

// pbSalinHalaman meniru `Property-Set <halaman> = <halaman lain>` (halaman tersemat): isi `ke` diganti salinan
// nilai dan daftar berawalan `dari.`.
func pbSalinHalaman(h *Halaman, dari, ke string) {
	pbKosongkan(h, ke)
	a := dari + "."
	nilai := map[string]string{}
	for k, v := range h.Nilai {
		if strings.HasPrefix(k, a) {
			nilai[ke+k[len(dari):]] = v
		}
	}
	daftar := map[string][]Baris{}
	for k, v := range h.Daftar {
		if strings.HasPrefix(k, a) {
			daftar[ke+k[len(dari):]] = salinBaris(v)
		}
	}
	for k, v := range nilai {
		h.Setel(k, v)
	}
	for k, v := range daftar {
		h.SetelDaftar(k, v)
	}
}

// pbKosongAtauNol = `.X=="" || .X==0` (teks kosong, atau bilangan bernilai nol).
func pbKosongAtauNol(s string) bool {
	if strings.TrimSpace(s) == "" {
		return true
	}
	if !AdalahDesimal(s) {
		return false
	}
	d, err := AngkaTeks("", s)
	return err == nil && d.IsZero()
}
