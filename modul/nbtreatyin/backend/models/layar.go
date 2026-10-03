package models

// Untuk apa berkas ini: ATURAN LAYAR - medan wajib, medan yang boleh diubah per
// posisi, dan data transform pasca-submit kedua flow action (tiket 10, 11, 12;
// spec §5.10, §5.11; AC 45-52, 71, 76).
//
// Sumber: `Section/DetailPolicyTreatyIn` (layar admin), `Section/
// DetailDeptHeadTreatyIn_UW` (layar Sec Head dan Dept Head), `Section/
// ListSuggest` (keduanya), `DataTransform/InboxPolicyTreatyIn_postDT`,
// `DataTransform/DeptHeadTreatyIn_UW_postDT`, `DataTransform/
// AddToListCommentsPolicyTreatyIn_DT` - INVENTARIS-XML.md bab 5 dan 7.

import "strings"

// ---------------------------------------------------------------- medan wajib

// MedanWajib adalah satu medan bertanda wajib di satu layar.
type MedanWajib struct {
	// Jalur halaman (relatif pyWorkPage).
	Jalur string
	// Label - VERBATIM label sel di section.
	Label string
	// Syarat - `pyRequired` bersyarat; nil = selalu wajib.
	Syarat func(h *Halaman) bool
}

// bukanNonProp = `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`.
func bukanNonProp(h *Halaman) bool {
	return h.Ambil("Quotation.ProportionalType") != JenisNonProporsional
}

// proporsionalQD = `.QuotationData.ProportionalType = 'Proportional'`.
func proporsionalQD(h *Halaman) bool {
	return h.Ambil(HalamanPolis+".QuotationData.ProportionalType") == JenisProporsional
}

// adaKlaim = `.Claim != ” && .Claim != 0`.
func adaKlaim(h *Halaman) bool {
	c := strings.TrimSpace(h.Ambil(HalamanPolis + ".Claim"))
	if c == "" {
		return false
	}
	d, err := AngkaTeks("Claim", c)
	return err != nil || d.Sign() != 0
}

// flagPPH = `.FlagPPH = true`.
func flagPPH(h *Halaman) bool { return h.Ambil(HalamanPolis+".FlagPPH") == "true" }

func wajib(m, label string, syarat func(*Halaman) bool) MedanWajib {
	return MedanWajib{Jalur: HalamanPolis + "." + m, Label: label, Syarat: syarat}
}

// medanWajibAdmin - `Section/DetailPolicyTreatyIn` + `Section/ListSuggest`.
//
// ⛔ `ListSuggest.ProductionDate` wajib HANYA untuk dua identitas orang
// (`OperatorID.pyUserIdentifier == <ID-operator-3> || <ID-operator-4>`) - salah
// satu dari 12 tempat tiket 05; TERTUNDA (AC 81), tidak diwajibkan.
var medanWajibAdmin = []MedanWajib{
	wajib("StartDate", "Statement Period", nil),
	wajib("QuotationData.IsSurveyReport", "Survey Report", bukanNonProp),
	wajib("TypeTax", "Type Tax", flagPPH),
	wajib("StatementDate", "Statement Date", nil),
	wajib("EndDate", "To", nil),
	wajib("Quartal", ".Quartal", proporsionalQD),
	wajib("YearOfQuartal", "Text Input", proporsionalQD),
	wajib("QuotationData.MOID", "Marketing Officer", nil),
	wajib("ClaimType", "Claim Type", adaKlaim),
	wajib("ClaimPaymentType", "Payment Type", adaKlaim),
	wajib("PremiOgp", "Premi Ogp", bukanNonProp),
	wajib("RiCommOgp", "(%) Deduction In A (OGP)", bukanNonProp),
	wajib("OveriddingCommOgp", "(%) Deduction In B (OGP)", bukanNonProp),
	wajib("PremiOnp", "Premi Onp", bukanNonProp),
	wajib("RiCommOnp", "(%) Deduction In A (ONP)", bukanNonProp),
	wajib("OveriddingCommOnp", "(%) Deduction In B (ONP)", bukanNonProp),
	wajib("Claim", "Claim", bukanNonProp),
	wajib("OutstandingClaim", "Outstanding Claim", bukanNonProp),
	wajib("SalvageValue", "Salvage", bukanNonProp),
	wajib("ExcessLoss", "Excess Loss", bukanNonProp),
	wajib("Deduction1", "Deduction1", bukanNonProp),
	wajib("Deduction2", "Deduction2", bukanNonProp),
	wajib("IsApproved", "Approval", nil),
	wajib("Suggest", "Suggest", nil),
}

// medanWajibAtasan - `Section/DetailDeptHeadTreatyIn_UW` + `Section/ListSuggest`.
// Seluruh medan uangnya wajib TANPA syarat, termasuk `ResultOnp1` (AC 47), dan
// tidak satu pun dari ClaimPaymentType, ClaimType, IDCurrency, Quartal,
// TypeTax, YearOfQuartal (AC 46).
var medanWajibAtasan = []MedanWajib{
	wajib("StartDate", "Statement Period", nil),
	wajib("StatementDate", "Statement Date", nil),
	wajib("EndDate", "To", nil),
	wajib("PremiOgp", "Premi Ogp", nil),
	wajib("RiCommOgp", "Deduction In A (OGP)", nil),
	wajib("OveriddingCommOgp", "Deduction In B (OGP)", nil),
	wajib("Claim", "Claim", nil),
	wajib("OutstandingClaim", "Outstanding Claim", nil),
	wajib("SalvageValue", "Salvage", nil),
	wajib("ExcessLoss", "Excess Loss", nil),
	wajib("PremiOnp", "Premi Onp", nil),
	wajib("RiCommOnp", "Deduction In A (ONP)", nil),
	wajib("ResultOnp1", "ResultOnp1", nil),
	wajib("OveriddingCommOnp", "Deduction In B (ONP)", nil),
	wajib("Deduction1", "Deduction1", nil),
	wajib("Deduction2", "Deduction2", nil),
	wajib("IsApproved", "Approval", nil),
	wajib("Suggest", "Suggest", nil),
}

// DaftarMedanWajib - medan wajib layar posisi itu.
func DaftarMedanWajib(posisi string) []MedanWajib {
	if posisi == PosisiAdmin {
		return medanWajibAdmin
	}
	if posisi == PosisiSecHead || posisi == PosisiDeptHead {
		return medanWajibAtasan
	}
	return nil
}

// MedanWajibKosong - label medan wajib yang kosong (AC 45, 48). Urutan = urutan
// layar.
func MedanWajibKosong(h *Halaman, posisi string) []string {
	var kosong []string
	for _, m := range DaftarMedanWajib(posisi) {
		if m.Syarat != nil && !m.Syarat(h) {
			continue
		}
		if strings.TrimSpace(h.Ambil(m.Jalur)) == "" {
			kosong = append(kosong, m.Label)
		}
	}
	return kosong
}

// ---------------------------------------------------------------- medan yang boleh diubah

// medanAtasan - medan yang TIDAK terkunci di layar atasan
// (`DetailDeptHeadTreatyIn_UW` + `ListSuggest`): kolom "Kunci" kosong atau
// hanya `nonaktif` bersyarat kosong. Selain ini terkunci permanen (AC 49-52):
// nilai kiriman layar untuk medan lain DIABAIKAN, nilai tersimpan dipakai.
var medanAtasan = []string{
	HalamanPolis + ".DueTo",
	HalamanPolis + ".FlagPPH",
	HalamanPolis + ".QuotationData.NoOfferSlip",
	HalamanPolis + ".IsApproved",
	HalamanPolis + ".Suggest",
	HalamanPolis + ".ProductionDate",
}

// medanAdmin - DAFTAR IZIN layar admin (`DetailPolicyTreatyIn` + `ListSuggest`):
//
//	isian   medan yang dapat diketik/dipilih di section (kolom "Kunci" kosong
//	        atau bersyarat yang tidak pernah benar bagi kasus treaty - IsUW)
//	tombol  hasil `TreatyEnableDisableInput`
//
// Medan TURUNAN (NetPremium, Balance*, PPN/PPH, BrokerageFee*, DueTo,
// Currency) TIDAK diterima dari layar - services menghitungnya ulang dari
// isian (`services.turunkan`).
//
// ⛔ DAFTAR IZIN, bukan daftar larangan (temuan tinjauan 2026-10-03): medan
// milik server - hasil pilih bisnis (NoOffer, TreatyGroupID, OJKBusinessID,
// BizCode, SOB, CedingCo, ...), StatementDate (`ALWAYS`), PolicyNo - tidak
// pernah dapat ditimpa layar, termasuk yang tidak tampil.
var medanAdmin = func() map[string]bool {
	m := map[string]bool{}
	for _, n := range []string{
		// isian
		"StartDate", "EndDate", "QuotationData.IsSurveyReport", "StatementType", "QuotationData.NoOfferSlip",
		"FlagRetroTreaty", "FlagPPH", "TypeTax", "Quartal", "YearOfQuartal", "QuotationData.MOID",
		"IDCurrency", "ClaimType", "ClaimPaymentType", "Remark", "GrossPremium", "GrossClaim",
		"PremiOgp", "RiCommOgp", "ResultOgp1", "OveriddingCommOgp", "ResultOgp2",
		"PremiOnp", "RiCommOnp", "ResultOnp1", "OveriddingCommOnp", "ResultOnp2",
		"Claim", "OutstandingClaim", "SalvageValue", "ExcessLoss", "Deduction1", "Deduction2",
		"Installment", "IsApproved", "Suggest", "ProductionDate",
		// hasil tombol "Enable / Disable Input Type" (TreatyEnableDisableInput)
		"IsNewPolicyNonProp", "QuotationData.ProportionalType",
	} {
		m[HalamanPolis+"."+n] = true
	}
	return m
}()

// GabungMasukanLayar menyalin nilai kiriman layar `masuk` ke halaman
// tersimpan `h`, HANYA untuk medan yang boleh diubah di posisi itu.
//
//	Admin   `medanAdmin` (daftar izin), beserta daftar SpreadingRiskList,
//	        ListInstallment, TreatyXOLList (baris boleh ditambah/dihapus -
//	        tombol Add/Delete layar admin)
//	Atasan  hanya `medanAtasan`; seluruh daftar terkunci
//
// Halaman Quotation dan TreatyIn tidak pernah diterima dari layar: Quotation
// diisi pilih bisnis dan CheckDataMkt; TreatyIn dibaca dari view.
func GabungMasukanLayar(h, masuk *Halaman, posisi string) {
	if masuk == nil {
		return
	}
	if posisi != PosisiAdmin {
		// Hanya medan yang DIKIRIM layar: medan yang tidak ada di kiriman
		// tetap bernilai tersimpan (bukan dikosongkan).
		for _, j := range medanAtasan {
			if v, ada := masuk.Nilai[j]; ada {
				h.Setel(j, v)
			}
		}
		return
	}
	for j, v := range masuk.Nilai {
		if medanAdmin[j] {
			h.Setel(j, v)
		}
	}
	for _, d := range []string{DaftarSpreading, DaftarAngsuran, HalamanPolis + ".TreatyXOLList"} {
		if b, ada := masuk.Daftar[d]; ada {
			h.SetelDaftar(d, salinBaris(b))
		}
	}
}

func salinBaris(b []Baris) []Baris {
	out := make([]Baris, len(b))
	for i, x := range b {
		nb := Baris{}
		for k, v := range x {
			nb[k] = v
		}
		out[i] = nb
	}
	return out
}

// ---------------------------------------------------------------- pasca submit

// TambahCatatan = `DataTransform/AddToListCommentsPolicyTreatyIn_DT`, dipanggil
// kedua DT pasca dengan Comment=.Suggest, Operator=.OperatorName,
// Approved=.IsApproved, Date=.SuggestDate (parameter APPLY_MODEL).
//
//	APPEND .PolicyTreatyIn.SuggestList:
//	  .Suggest = Param.Comment
//	  WHEN Param.Approved != "" -> .IsApproved = Param.Approved
//	  .Date = Param.Date ; .OperatorName = Param.Operator
//
// `operatorID` (identitas akses login) DITAMBAHKAN sistem baru supaya catatan
// tersimpan bersama operatornya (AC 71; P4 - OPERATORID = identitas login).
// ⛔ Parameter ApprovedtoDeptHead tidak dibangun (P36, AC 64).
func TambahCatatan(h *Halaman, operatorID string) {
	b := Baris{
		"Suggest":      h.Ambil(HalamanPolis + ".Suggest"),
		"Date":         h.Ambil(HalamanPolis + ".SuggestDate"),
		"OperatorName": h.Ambil(HalamanPolis + ".OperatorName"),
		"OperatorID":   operatorID,
	}
	if a := h.Ambil(HalamanPolis + ".IsApproved"); a != "" {
		b["IsApproved"] = a
	}
	h.SetelDaftar(DaftarUsulan, append(h.AmbilDaftar(DaftarUsulan), b))
}

// PascaAdmin = `DataTransform/InboxPolicyTreatyIn_postDT`.
//
//	1   AddToListComments          -> TambahCatatan
//	2-3 HasFacOut "0", TestTreatyToFacStatus -> TetapkanHasFacOut (AC 76)
//	4   PositionNote==Admin && IsApproved=="0" ->
//	    NBStatus = "NB WAS DECLINED BY  " + @toUpperCase(OperatorID.pyUserName)
func PascaAdmin(h *Halaman, operatorID, namaTampilan string) {
	TambahCatatan(h, operatorID)
	TetapkanHasFacOut(h)
	if h.Ambil("PositionNote") == PosisiAdmin && h.Ambil(HalamanPolis+".IsApproved") == NilaiDitolak {
		h.Setel("NBStatus", TeksNBStatusDitolakOleh(namaTampilan))
	}
}

// PascaAtasan = `DataTransform/DeptHeadTreatyIn_UW_postDT`.
//
//	1   AddToListComments -> TambahCatatan
//	2   PositionNote==Admin && "0"   -> tidak pernah benar di layar atasan; ditiru
//	5   PositionNote==DeptHead && "0" -> "NB IS IN "+@toUpperCase(pxCreateOpName)+"'S INBOX"
//	3,4,6,7  posisi Director / GroupLeader - dibuang P13 (AC 10), tidak dibangun
//
// `namaPembuat` = nama tampilan pembuat kasus (`pxCreateOpName`), dari data.
// ⚠️ Connector penolakan sesudahnya menulis `NBStatus = ""` (Decision2/11 No),
// sehingga nilai langkah 5 tertimpa - urutan itu ditiru services.
func PascaAtasan(h *Halaman, operatorID, namaTampilan, namaPembuat string) {
	TambahCatatan(h, operatorID)
	ditolak := h.Ambil(HalamanPolis+".IsApproved") == NilaiDitolak
	switch {
	case h.Ambil("PositionNote") == PosisiAdmin && ditolak:
		h.Setel("NBStatus", TeksNBStatusDitolakOleh(namaTampilan))
	case h.Ambil("PositionNote") == PosisiDeptHead && ditolak:
		h.Setel("NBStatus", TeksNBStatusKotakMasuk(namaPembuat))
	}
}

// ---------------------------------------------------------------- tombol layar

// TombolKirim adalah tombol submit yang tampil di layar, menurut
// `Section/DetailPolicyTreatyIn` dan `Section/DetailDeptHeadTreatyIn_UW`.
type TombolKirim string

const (
	// TombolTidakAda - IsApproved belum bernilai yang memunculkan tombol.
	TombolTidakAda TombolKirim = ""
	// TombolKirimLangsung - `finishAssignment` (admin IsApproved==1 didahului
	// SetDueTo_act; atasan IsApproved='0'; Sec Head IsApproved==1).
	TombolKirimLangsung TombolKirim = "kirim"
	// TombolKonfirmasiTolak - admin IsApproved==0: `PolicyTreatyInDeclineConfirm`
	// (Yes = finishAssignment).
	TombolKonfirmasiTolak TombolKirim = "konfirmasi-tolak"
	// TombolNomorPolis - Dept Head IsApproved==1: `GeneratePolicyNoTreaty_Act`
	// lalu `ShowPolicyNoTreaty` (OK = finishAssignment).
	TombolNomorPolis TombolKirim = "nomor-polis"
)

// TombolUntuk menentukan tombol submit yang tampil.
//
// ⛔ Syarat identitas orang tombol atasan (`OperatorID.pyUserIdentifier ==
// <ID-operator-1>`) DIGANTI posisi kasus: Dept Head - jenjang terakhir tangga
// P13 - satu-satunya yang menerbitkan nomor polis; Sec Head selalu
// menaikkan (AC 8, WO). Catatan di tiket 05.
func TombolUntuk(h *Halaman, posisi string) TombolKirim {
	ia := h.Ambil(HalamanPolis + ".IsApproved")
	switch posisi {
	case PosisiAdmin:
		switch ia {
		case "1":
			return TombolKirimLangsung
		case NilaiDitolak:
			return TombolKonfirmasiTolak
		}
	case PosisiSecHead:
		if ia == NilaiDitolak || ia == "1" {
			return TombolKirimLangsung
		}
	case PosisiDeptHead:
		switch ia {
		case NilaiDitolak:
			return TombolKirimLangsung
		case "1":
			return TombolNomorPolis
		}
	}
	return TombolTidakAda
}
