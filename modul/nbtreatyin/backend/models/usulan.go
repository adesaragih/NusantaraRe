package models

// Untuk apa berkas ini: CATATAN USULAN PER TAHAP (`PolicyTreatyIn.SuggestList`)
// <-> tabel lama `POOLDATA.HISTORYAKSEPTASIPRODUCTION` (spec-penyimpanan ID-5,
// ID-31; AC 39-44; diagram grilling sheet NB Treaty In Prop/NonProp:
// "POOLDATA.HISTORYAKSEPTASIPRODUCTION 1:N IDPEGA · 15 kolom · SUDAH datar ·
// ditulis oleh SaveViewSuggest -> InsertViewSuggest_SQL").
//
// Sumber XML:
//
//	Activity/SaveViewSuggest langkah 2 "UNTUK TREATY", kalang
//	pyWorkPage.PolicyTreatyIn.SuggestList, per baris (2.1) bila .IsSave == "":
//	  2.1.2 CARI1  = @replaceAll(pyWorkPage.pyWorkIDPrefix,"-","")   -> TYPE_POLIS
//	        CARI2  = .pxListSubscript                                -> NOURUT
//	        CARI3  = "Policy"                                        -> POSISI
//	        CARI4  = .OperatorName                                   -> PIC
//	        CARI5  = @FormatDateTime(.Date,"dd/MM/yyyy hh:mm:ss",..) -> TGL_INP
//	        CARI6  = OperatorID.pyOrgDivision                        -> DIV
//	        CARI7  = pyWorkPage.Quotation.BusinessFac                -> TYPE
//	        CARI8  = 2                                               -> PUTARAN
//	        CARI9  = @if(.IsApproved="1","Accept",@if(.IsApproved="0","Reject","")) -> APPROVAL
//	        CARI10 = @substring(.Suggest,0,3990)                     -> KETERANGAN
//	  2.1.3 RDB-List InsertViewSuggest_SQL (+ IDPEGA pyWorkPage.pzInsKey,
//	        AKSES_LOGIN OperatorID.pyUserIdentifier, B2B OfferFacIn.IsB2B,
//	        BUSINESS_CODE Quotation.BusinessCode, PERCENT_RNM OfferFacIn.PercentShare)
//	  2.1.4 .IsSave = "Yes"
//
// `[penyimpangan sadar]` (keputusan work owner K4 03-10-2026, grilling ID-31 /
// AC 39) terhadap XML, ditulis juga di tiket 10:
//  1. Syarat langkah 2 `Quotation.BusinessFac == "F"` TIDAK ditiru - kasus
//     treaty bernilai "T", sehingga di Pega kalang ini tidak pernah menulis.
//  2. XML memanggilnya HANYA dari `InputPolicyTreatyInPost_Act` langkah 4
//     (pasca-submit admin); sistem baru menulis baris yang ditambahkan SETIAP
//     submit (admin, Sec Head, Dept Head) di transaksi submit itu - catatan
//     jenjang atasan yang menyelesaikan kasus tidak pernah hilang (AC 71).
//  3. TGL_INP: XML memformat `hh` (jam 12) lalu `To_date(..,'HH24')` - catatan
//     sore tersimpan sebagai pagi. Di sini jam 24 apa adanya: layar (riwayat
//     catatan) membaca balik tabel ini, dan Pega menampilkan `.Date` halaman.
//  4. NOURUT diberikan repository (berikutnya di bawah kunci kasus), bukan
//     `.pxListSubscript` halaman - dua submit serentak tidak berbagi nomor.
//
// ⛔ DIV (`OperatorID.pyOrgDivision`) tidak punya sumber di inti (`inti.Pelaku`
// hanya AkunID dan Peran) - ditulis NULL, butir terbuka. B2B dan PERCENT_RNM
// dari halaman `OfferFacIn` (kasus fakultatif) - NULL di NB.

import "strings"

// Nilai tetap InsertViewSuggest_SQL untuk kasus treaty.
const (
	// PosisiUsulanProduksi - CARI3 kalang "UNTUK TREATY".
	PosisiUsulanProduksi = "Policy"
	// PutaranUsulanProduksi - CARI8 kalang "UNTUK TREATY".
	PutaranUsulanProduksi = "2"
	// PanjangKeterangan - `substr(CARI10, 0, 3990)` (AC 40).
	PanjangKeterangan = 3990
	// TandaTersimpan - `.IsSave = "Yes"` (langkah 2.1.4).
	TandaTersimpan = "Yes"
)

// UsulanProduksi adalah satu baris `POOLDATA.HISTORYAKSEPTASIPRODUCTION`.
// IDPEGA (`pyWorkPage.pzInsKey`, `KunciInstans`) dibawa repository.
type UsulanProduksi struct {
	// NoUrut - NOURUT; diberikan repository saat menulis.
	NoUrut       int
	TypePolis    string
	Posisi       string
	PIC          string
	TglInp       string // "2006-01-02 15:04:05"
	Div          string
	Type         string
	Putaran      string
	Approval     string
	Keterangan   string
	AksesLogin   string
	B2B          string
	BusinessCode string
	PercentRNM   string
}

// approvalUsulan = CARI9.
func approvalUsulan(isApproved string) string {
	switch isApproved {
	case "1":
		return "Accept"
	case NilaiDitolak:
		return "Reject"
	}
	return ""
}

// potongKarakter = `substr(x, 0, n)` Oracle: n KARAKTER pertama.
func potongKarakter(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

// UsulanBelumTersimpan = `SaveViewSuggest` langkah 2.1-2.1.4: setiap baris
// SuggestList yang `IsSave`-nya kosong dipetakan ke baris riwayat produksi,
// lalu ditandai `IsSave = "Yes"` di halaman.
//
// `OperatorID` baris (identitas akses login penulis catatan,
// `TambahCatatan`) menjadi AKSES_LOGIN - di XML `OperatorID.pyUserIdentifier`
// operator yang menjalankan SaveViewSuggest; karena baris ditulis pada submit
// yang menambahkannya, keduanya orang yang sama (P4, AC 43).
func UsulanBelumTersimpan(h *Halaman) []UsulanProduksi {
	daftar := h.AmbilDaftar(DaftarUsulan)
	var out []UsulanProduksi
	for _, b := range daftar {
		if b["IsSave"] != "" {
			continue
		}
		out = append(out, UsulanProduksi{
			TypePolis:    strings.ReplaceAll(AwalanKasus, "-", ""),
			Posisi:       PosisiUsulanProduksi,
			PIC:          b["OperatorName"],
			TglInp:       b["Date"],
			Type:         h.Ambil(HalamanQuotation + ".BusinessFac"), // CARI7; "T" sejak kasus lahir
			Putaran:      PutaranUsulanProduksi,
			Approval:     approvalUsulan(b["IsApproved"]),
			Keterangan:   potongKarakter(b["Suggest"], PanjangKeterangan),
			AksesLogin:   b["OperatorID"],
			BusinessCode: h.Ambil(HalamanQuotation + ".BusinessCode"),
		})
		b["IsSave"] = TandaTersimpan
	}
	return out
}

// BarisCatatan memetakan balik satu baris riwayat produksi menjadi baris
// SuggestList layar (`Section/ListSuggest`: Date, OperatorName, IsApproved,
// Suggest). APPROVAL kosong tidak memasang IsApproved - sama dengan
// `AddToListCommentsPolicyTreatyIn_DT` langkah 1.2 (`Param.Approved != ""`).
func BarisCatatan(u UsulanProduksi) Baris {
	b := Baris{
		"Suggest":      u.Keterangan,
		"Date":         u.TglInp,
		"OperatorName": u.PIC,
		"OperatorID":   u.AksesLogin,
		"IsSave":       TandaTersimpan,
	}
	switch u.Approval {
	case "Accept":
		b["IsApproved"] = "1"
	case "Reject":
		b["IsApproved"] = NilaiDitolak
	}
	return b
}
