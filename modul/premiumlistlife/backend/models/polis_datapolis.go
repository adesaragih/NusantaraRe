package models

// Data polis layar Input Premium Detail - tiket 03 bagian 2 (bagian 1-3 layar).
// MURNI.
//
// Untuk apa berkas ini: isian `Section/ShowLifePremiumDetail.xml` yang DAPAT
// DIISI di tahap Input Premium Detail, di luar unggah CSV dan perhitungan.
//
// `[terverifikasi]` dibaca sebagai pohon 01-10-2026 (salinan korpus
// `kelvin\PremiumListLife (Done)\`):
//
//	Choose Product Name -> ChooseProdName (BrowseProductForNB_Life, param
//	                       Ceding = pyWorkPage.CedingCo) -> SetProdNametoPolis:
//	                       SobName=.SOBNAME, SourceOfBusiness=.SOBID,
//	                       ProductName=.INWARDNAME, ProductNameID=.ID,
//	                       CedingCo=.CEDINGID, CedingCoName=.CEDING,
//	                       PolicyHolder=.POLICYHODER, PolicyHolderName=.POLICYHODERNAME
//	.Type             pxDropdown  "Type"                    pyRequired
//	.PremiumListSummary.RISLIPRNM "R/I SLIP RNM No."        wajib & tampil bila TR/TP
//	                  (GetPLNumber_Act -> GetPLandNopolis_sql: NOPOLIS RNML-Q/RNML-F)
//	.ProRateType      pxDropdown  "Premium Payment Method"  pyRequired
//	.MarketingName    pxAutoComplete "Marketing Officer"    pyRequired
//	                  (BrowseMarketingOfficer_RD: .ID->MOID, .ClientID->MarketingCode)
//	.AnnuityInterest  pxNumber "Annuity Interest"           pyRequired
//	.PremiumRefundFactor pxNumber "Premium Refund Factor"   pyRequired
//	.RetroName        pxAutoComplete "Billing Name" (+RetroID) - label layar
//	                  "BILLING NAME IS MANDATORY FOR TYPE TP & TR"
//	.SecurityReinsurer pxAutoComplete "Retrocessionaire" (+SecurityReinsurerID)
//	                  setSecurityReinsurer_act saat Billing Name berubah
//
// ⚠️ Sel lain di layar itu (Ceding, Policy Holder, tanggal-tanggal, dan
// seterusnya) adalah data penawaran - DITAMPILKAN, tidak diisi di sini.
//
// ⛔ `Save Data` hanya menyimpan + memeriksa medan wajib (keputusan work owner
// 05-10-2026); `Calculate1_Act` TIDAK dijalankan - keputusan work owner
// 01-10-2026. Batas produk `SavePremiumList_Act` 6-8.2 kini penolakan Validate
// CSV - lihat PeriksaBatasProduk di bawah.
//
// Dibaca sesudah: polis_isianpenawaran.go.

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"
)

// PilihanTypePolis - dropdown `Type`.
//
// `[terverifikasi]` keempat kode yang korpus bandingkan dengan `.Type`
// (InsertJsonPolisLife_Act, AppendCurrencySummary_DT, tombol View Upload per
// Type). ⚠️ Teks tampilannya kode itu sendiri: pilihan dropdown milik rule
// properti yang tidak diekspor, dan kepanjangannya tidak dikarang.
var PilihanTypePolis = []Pilihan{
	{Kode: "QR", Nama: "QR"},
	{Kode: "QP", Nama: "QP"},
	{Kode: "TP", Nama: "TP"},
	{Kode: "TR", Nama: "TR"},
}

// PilihanProRateType - dropdown `Premium Payment Method`.
//
// Kode 1/2/3 `[terverifikasi]` `Calculate1_Act`: `@if(ProRateType==1,"AP",
// @if(ProRateType==2,"PY","PM"))`. Teks tampilnya `[keputusan work owner
// 01-10-2026]`: 1 Single, 2 Annually, 3 Others - pilihan radio/dropdown aslinya
// milik rule properti yang tidak diekspor. ⚠️ Singkatan AP/PY/PM di
// Calculate1_Act adalah kode turunan untuk perhitungan (tidak dibawa), BUKAN
// teks dropdown.
var PilihanProRateType = []Pilihan{
	{Kode: "1", Nama: "Single"},
	{Kode: "2", Nama: "Annually"},
	{Kode: "3", Nama: "Others"},
}

// Pesan wajib-isi layar Input Premium Detail.
//
// ⚠️ Lima yang pertama VERBATIM `ProtectAccept` langkah 5 (PesanTypeKosong dan
// kawan-kawan di polis_validasi.go dipakai ulang). Sisanya karangan layar baru
// berpola sama: sel ber-pyRequired tanpa kalimat korpus.
const (
	PesanRISlipKosong  = "R/I SLIP RNM No. can't null"
	PesanBillingKosong = "Billing Name can't null"
	// PesanRetroKosong - Retrocessionaire wajib untuk TP/TR (keputusan work owner
	// 01-10-2026; di Pega selnya tidak ber-pyRequired).
	PesanRetroKosong           = "Retrocessionaire can't null"
	PesanAnnuityInterestKosong = "Annuity Interest can't null"
	PesanPremiumRefundKosong   = "Premium Refund Factor can't null"
	// PesanDateReceivedKosong - Email Received Date wajib di Input Premium
	// Detail (keputusan work owner 02-10-2026).
	PesanDateReceivedKosong = "Email Received Date can't null"
)

// ErrDataPolisTidakSah - isian di luar pilihan tertutupnya.
var ErrDataPolisTidakSah = errors.New("models: data polis tidak sah")

// ErrDataPolisBelumLengkap - isian wajib yang kosong.
var ErrDataPolisBelumLengkap = errors.New("models: data polis belum lengkap")

// IsianDataPolis adalah yang pemakai kirim dari layar Input Premium Detail.
type IsianDataPolis struct {
	Type string
	// Produk dan turunannya (SetProdNametoPolis) - dipilih dari popup.
	ProductNameID    string
	ProductName      string
	SourceOfBusiness string
	SobName          string
	CedingCo         string
	CedingCoName     string
	PolicyHolder     string
	PolicyHolderName string

	RISlipRNM           string
	ProRateType         string
	MoID                string
	MarketingCode       string
	MarketingName       string
	AnnuityInterest     *apd.Decimal
	PremiumRefundFactor *apd.Decimal
	RetroID             string
	RetroName           string
	SecurityReinsurerID string
	SecurityReinsurer   string
	// DateReceived - `Email Received Date` (`DATE_RECEIVED`, sel Input Offer).
	// ⛔ [keputusan work owner 02-10-2026] DAPAT DIISI dan WAJIB di layar
	// Input Premium Detail; kolomnya sama dengan yang ditulis Save Offer.
	DateReceived *time.Time
}

// TypeRetro - Type yang menuntut R/I SLIP dan Billing Name.
func TypeRetro(t string) bool { return t == "TP" || t == "TR" }

// SecurityReinsurerOtomatis meniru `setSecurityReinsurer_act` langkah 1-3:
// Retrocessionaire dikosongkan, lalu diisi untuk dua Billing Name tertentu.
//
// `[terverifikasi]` RetroID L0000141 -> L0000134; L0000135 / L0000137 ->
// B0000020. Nama perusahaannya VERBATIM activity itu (badan usaha, bukan
// nama orang). Selain itu - kosong, dipilih pemakai.
func SecurityReinsurerOtomatis(retroID string) (id, nama string, ada bool) {
	switch retroID {
	case "L0000141":
		return "L0000134", "AON REINSURANCE BROKERS INDONESIA", true
	case "L0000135", "L0000137":
		return "B0000020", "IBS REINSURANCE BROKERS", true
	}
	return "", "", false
}

// SusunDataPolis memeriksa isian data polis.
//
// ⛔ SELURUH yang kosong dilaporkan sekaligus - pola `ProtectAccept`.
func SusunDataPolis(isi IsianDataPolis) (IsianDataPolis, error) {
	for _, s := range []*string{&isi.Type, &isi.ProductNameID, &isi.ProductName, &isi.SourceOfBusiness,
		&isi.SobName, &isi.CedingCo, &isi.CedingCoName, &isi.PolicyHolder, &isi.PolicyHolderName,
		&isi.RISlipRNM, &isi.ProRateType, &isi.MoID, &isi.MarketingCode, &isi.MarketingName,
		&isi.RetroID, &isi.RetroName, &isi.SecurityReinsurerID, &isi.SecurityReinsurer} {
		*s = strings.TrimSpace(*s)
	}
	if isi.Type != "" && namaPilihan(PilihanTypePolis, isi.Type) == "" {
		return isi, fmt.Errorf("%w: Type %q", ErrDataPolisTidakSah, isi.Type)
	}
	if isi.ProRateType != "" && namaPilihan(PilihanProRateType, isi.ProRateType) == "" {
		return isi, fmt.Errorf("%w: Premium Payment Method %q", ErrDataPolisTidakSah, isi.ProRateType)
	}
	for nama, pasangan := range map[string][2]string{
		"Product Name":      {isi.ProductNameID, isi.ProductName},
		"Marketing Officer": {isi.MarketingCode, isi.MarketingName},
		"Billing Name":      {isi.RetroID, isi.RetroName},
		"Retrocessionaire":  {isi.SecurityReinsurerID, isi.SecurityReinsurer},
	} {
		if (pasangan[0] == "") != (pasangan[1] == "") {
			return isi, fmt.Errorf("%w: %s tanpa pasangan kode-nama", ErrDataPolisTidakSah, nama)
		}
	}
	// R/I SLIP hanya bermakna untuk TP/TR (selnya pun hanya tampil di sana).
	if !TypeRetro(isi.Type) {
		isi.RISlipRNM = ""
		// Billing Name dan Retrocessionaire hanya untuk TP/TR - layar
		// menyembunyikannya di Type lain (keputusan work owner 01-10-2026);
		// nilai tersembunyi tidak boleh ikut tersimpan.
		isi.RetroID, isi.RetroName = "", ""
		isi.SecurityReinsurerID, isi.SecurityReinsurer = "", ""
	}
	if id, nama, ada := SecurityReinsurerOtomatis(isi.RetroID); ada {
		isi.SecurityReinsurerID, isi.SecurityReinsurer = id, nama
	}
	var kurang []string
	tambah := func(kosongkah bool, pesan string) {
		if kosongkah {
			kurang = append(kurang, pesan)
		}
	}
	tambah(isi.ProductName == "", PesanProductNameKosong)
	tambah(isi.Type == "", PesanTypeKosong)
	tambah(TypeRetro(isi.Type) && isi.RISlipRNM == "", PesanRISlipKosong)
	tambah(isi.ProRateType == "", PesanProRateTypeKosong)
	tambah(isi.MarketingName == "", PesanMarketingKosong)
	tambah(isi.AnnuityInterest == nil, PesanAnnuityInterestKosong)
	tambah(isi.PremiumRefundFactor == nil, PesanPremiumRefundKosong)
	tambah(isi.SourceOfBusiness == "", PesanSOBKosong)
	tambah(TypeRetro(isi.Type) && isi.RetroName == "", PesanBillingKosong)
	tambah(TypeRetro(isi.Type) && isi.SecurityReinsurer == "", PesanRetroKosong)
	tambah(isi.DateReceived == nil, PesanDateReceivedKosong)
	if len(kurang) > 0 {
		return isi, fmt.Errorf("%w: %s", ErrDataPolisBelumLengkap, GabungPesanPenawaran(kurang))
	}
	return isi, nil
}
