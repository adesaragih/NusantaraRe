package models

// Untuk apa berkas ini: PROTEKSI SUBMIT - port `Activity/Protection_Act` (kelas Data-PolicyTreatyIn; 20 langkah,
// 2 `//`) dan `Activity/ProtectionNonProp_Act` (7 langkah, khas EDM). Pemicu: radio Approval `.IsApproved`
// `Section/ListSuggestEDM` (change/click -> refresh + Protection_Act) - kedua layar. Tombol Submit admin
// (`IsApproved==1`) menyelesaikan assignment hanya bila `Protect.CARI1 = 0 && Protect.CARI2 = 0`; pesan halaman
// (Page-Set-Messages / Property-Set-Messages) menahan submit di posisi mana pun.
//
//	Protection_Act 1     Protect.CARI1..3 = 0; pesan
//	               2     Quotation.ProportionalType == "NonProportional" -> lompat SkipDate (14)
//	               4-5   SurveyReportList < 1 && IsSurveyReport == "Yes" -> CARI1 = 1, "Survey report can't be empty"
//	               6-7   `//` (PolicyTreatyInDetail kosong -> CARI3) - tidak diport
//	               8-10  kalang PolicyTreatyIn.PolicyTreatyInDetail (total %Share 100) -> CARI2 - ⚠️ daftar itu tidak
//	                     pernah ditulis rule EDM mana pun (sensus korpus): kalang nol putaran, CARI2 tetap 0
//	               11    NoOffer == "" -> "Please choose business"
//	               12-13 ProductionDate kosong -> tanggal SQL (GetSQLDate)
//	               15    NonProportional -> ProtectionNonProp_Act
//	ProtectionNonProp 2  NoOffer == "" -> "Please choose business"
//	                  3,5 MarketingOfficer == "" -> "Please choose Marketing Officer" (teks langkah 1 VERBATIM)
//	                  4  SurveyReportList < 1 && IsSurveyReport == "Yes" -> "Survey report can't be empty"
//	                  6-7 ProductionDate kosong -> tanggal SQL

import (
	"time"

	"nusantarare/inti/backend/utils"
)

// Pesan VERBATIM proteksi.
const (
	PesanPilihBisnis    = "Please choose business"
	PesanSurveiKosong   = "Survey report can't be empty"
	PesanPilihMarketing = "Please choose Marketing Officer"
)

// Proteksi adalah hasil `Protection_Act` (`Protect.CARI1` survei, `Protect.CARI2` share).
type Proteksi struct {
	CARI1, CARI2 int
}

// surveiKosong = `(@LengthOfPageList(.QuotationData.SurveyReportList)<1 && .QuotationData.IsSurveyReport=="Yes")`.
func surveiKosong(h *Halaman) bool {
	return len(h.AmbilDaftar(DaftarSurvei)) < 1 && h.Ambil(pt+"QuotationData.IsSurveyReport") == "Yes"
}

// ProtectionAct menjalankan `Protection_Act` atas halaman: pesan dipasang di halaman (`TambahPesan`), dan
// ProductionDate kosong diisi `sekarang` (GetSQLDate).
func ProtectionAct(h *Halaman, sekarang time.Time) Proteksi {
	var p Proteksi
	if h.Ambil(HalamanQuotation+".ProportionalType") == JenisNonProporsional {
		protectionNonProp(h, sekarang)
		return p
	}
	if surveiKosong(h) {
		p.CARI1 = 1
		h.TambahPesan("", PesanSurveiKosong)
	}
	if h.Ambil(pt+"NoOffer") == "" {
		h.TambahPesan("", PesanPilihBisnis)
	}
	if h.Ambil(pt+"ProductionDate") == "" {
		h.Setel(pt+"ProductionDate", utils.FormatTanggalWaktu(sekarang))
	}
	return p
}

// protectionNonProp = `ProtectionNonProp_Act` langkah 2-7.
func protectionNonProp(h *Halaman, sekarang time.Time) {
	if h.Ambil(pt+"NoOffer") == "" {
		h.TambahPesan("", PesanPilihBisnis)
	}
	if h.Ambil(pt+"MarketingOfficer") == "" {
		h.TambahPesan("", PesanPilihMarketing)
	}
	if surveiKosong(h) {
		h.TambahPesan("", PesanSurveiKosong)
	}
	if h.Ambil(pt+"ProductionDate") == "" {
		h.Setel(pt+"ProductionDate", utils.FormatTanggalWaktu(sekarang))
	}
}
