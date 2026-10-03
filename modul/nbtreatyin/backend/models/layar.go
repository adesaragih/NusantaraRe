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
	return h.Ambil(HalamanQuotation+".ProportionalType") != JenisNonProporsional
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

// ---- wadah (`pyContainerVisibleWhen`): sel di wadah tersembunyi tidak
// ter-render, jadi `pyRequired`-nya tidak berlaku (dibaca ulang 2026-10-03).

// bukanXOLRetro = wadah FlagPPH/TypeTax/Choose Business `DetailPolicyTreatyIn`:
// `.ClaimType != 'XOL Retro'`.
func bukanXOLRetro(h *Halaman) bool { return h.Ambil(HalamanPolis+".ClaimType") != KlaimXOLRetro }

// bukanNonPropBaru = `.IsNewPolicyNonProp != 1` - wadah bagian uang,
// spreading, angsuran `DetailDeptHeadTreatyIn_UW`.
func bukanNonPropBaru(h *Halaman) bool { return h.Ambil(HalamanPolis+".IsNewPolicyNonProp") != "1" }

// wadahUangAdmin = `.IsNewPolicyNonProp != 1 && .IsNewPolicyListFormat != 1` -
// wadah bagian uang, spreading, angsuran `DetailPolicyTreatyIn`.
// ⚠️ `IsNewPolicyListFormat` tidak diisi rule mana pun di korpus (hanya dibaca
// dua section ini) - syaratnya ditiru apa adanya.
func wadahUangAdmin(h *Halaman) bool {
	return bukanNonPropBaru(h) && h.Ambil(HalamanPolis+".IsNewPolicyListFormat") != "1"
}

// dan menggabungkan syarat (nil = selalu).
func dan(fs ...func(*Halaman) bool) func(*Halaman) bool {
	return func(h *Halaman) bool {
		for _, f := range fs {
			if f != nil && !f(h) {
				return false
			}
		}
		return true
	}
}

func wajib(m, label string, syarat func(*Halaman) bool) MedanWajib {
	return MedanWajib{Jalur: HalamanPolis + "." + m, Label: label, Syarat: syarat}
}

// wajibUangAdmin - medan wajib bagian uang layar admin: `pyRequired`
// `bukanNonProp` DAN wadahnya tampil.
func wajibUangAdmin(m, label string) MedanWajib {
	return wajib(m, label, dan(bukanNonProp, wadahUangAdmin))
}

// wajibUangAtasan - medan wajib bagian uang layar atasan: `pyRequired` tanpa
// syarat, wadahnya `.IsNewPolicyNonProp != 1`.
func wajibUangAtasan(m, label string) MedanWajib { return wajib(m, label, bukanNonPropBaru) }

// medanWajibAdmin - `Section/DetailPolicyTreatyIn` + `Section/ListSuggest`.
//
// `ListSuggest.ProductionDate` wajib HANYA untuk dua identitas orang - tempat
// berperan tiket 05 (`TanggalProduksiWajib`), ditambahkan `MedanWajibBerlaku`;
// tertunda selama pemetaannya kosong (AC 81).
var medanWajibAdmin = []MedanWajib{
	wajib("StartDate", "Statement Period", nil),
	wajib("QuotationData.IsSurveyReport", "Survey Report", bukanNonProp),
	wajib("TypeTax", "Type Tax", dan(flagPPH, bukanXOLRetro)),
	wajib("StatementDate", "Statement Date", nil),
	wajib("EndDate", "To", nil),
	wajib("Quartal", ".Quartal", proporsionalQD),
	wajib("YearOfQuartal", "Text Input", proporsionalQD),
	wajib("QuotationData.MOID", "Marketing Officer", nil),
	wajib("ClaimType", "Claim Type", adaKlaim),
	wajib("ClaimPaymentType", "Payment Type", adaKlaim),
	wajibUangAdmin("PremiOgp", "Premi Ogp"),
	wajibUangAdmin("RiCommOgp", "(%) Deduction In A (OGP)"),
	wajibUangAdmin("OveriddingCommOgp", "(%) Deduction In B (OGP)"),
	wajibUangAdmin("PremiOnp", "Premi Onp"),
	wajibUangAdmin("RiCommOnp", "(%) Deduction In A (ONP)"),
	wajibUangAdmin("OveriddingCommOnp", "(%) Deduction In B (ONP)"),
	wajibUangAdmin("Claim", "Claim"),
	wajibUangAdmin("OutstandingClaim", "Outstanding Claim"),
	wajibUangAdmin("SalvageValue", "Salvage"),
	wajibUangAdmin("ExcessLoss", "Excess Loss"),
	wajibUangAdmin("Deduction1", "Deduction1"),
	wajibUangAdmin("Deduction2", "Deduction2"),
	wajib("IsApproved", "Approval", nil),
	wajib("Suggest", "Suggest", nil),
}

// medanWajibAtasan - `Section/DetailDeptHeadTreatyIn_UW` + `Section/ListSuggest`.
// Seluruh medan uangnya `pyRequired` TANPA syarat, termasuk `ResultOnp1` (AC 47),
// dan tidak satu pun dari ClaimPaymentType, ClaimType, IDCurrency, Quartal,
// TypeTax, YearOfQuartal (AC 46). Satu-satunya syarat: wadah bagian uang
// tampil (`.IsNewPolicyNonProp != 1`).
var medanWajibAtasan = []MedanWajib{
	wajib("StartDate", "Statement Period", nil),
	wajib("StatementDate", "Statement Date", nil),
	wajib("EndDate", "To", nil),
	wajibUangAtasan("PremiOgp", "Premi Ogp"),
	wajibUangAtasan("RiCommOgp", "Deduction In A (OGP)"),
	wajibUangAtasan("OveriddingCommOgp", "Deduction In B (OGP)"),
	wajibUangAtasan("Claim", "Claim"),
	wajibUangAtasan("OutstandingClaim", "Outstanding Claim"),
	wajibUangAtasan("SalvageValue", "Salvage"),
	wajibUangAtasan("ExcessLoss", "Excess Loss"),
	wajibUangAtasan("PremiOnp", "Premi Onp"),
	wajibUangAtasan("RiCommOnp", "Deduction In A (ONP)"),
	wajibUangAtasan("ResultOnp1", "ResultOnp1"),
	wajibUangAtasan("OveriddingCommOnp", "Deduction In B (ONP)"),
	wajibUangAtasan("Deduction1", "Deduction1"),
	wajibUangAtasan("Deduction2", "Deduction2"),
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

// medanTanggalProduksi - `Section/ListSuggest` `.ProductionDate` (label sel).
var medanTanggalProduksi = MedanWajib{Jalur: jalurTanggalProduksi, Label: "Production Date"}

// MedanWajibBerlaku - medan wajib layar posisi itu yang BERLAKU saat ini:
// syarat `pyRequired` (dan wadahnya) terpenuhi, ditambah `.ProductionDate`
// bila tempat berperannya terbuka bagi pelaku (`TanggalProduksiWajib`, tiket
// 05; `tempat` = `TempatTampil` pelaku). Urutan = urutan layar. SATU sumber
// untuk daftar wajib layar (`Layar.MedanWajib`), tombol Save (AC 48), dan
// submit (AC 45).
func MedanWajibBerlaku(h *Halaman, posisi string, tempat map[string]bool) []MedanWajib {
	daftar := DaftarMedanWajib(posisi)
	if daftar == nil {
		return nil
	}
	var out []MedanWajib
	for _, m := range daftar {
		if m.Syarat == nil || m.Syarat(h) {
			out = append(out, m)
		}
	}
	if TanggalProduksiWajib(h, tempat) {
		out = append(out, medanTanggalProduksi)
	}
	return out
}

// MedanWajibKosong - label medan wajib berlaku yang kosong (AC 45, 48).
func MedanWajibKosong(h *Halaman, posisi string, tempat map[string]bool) []string {
	var kosong []string
	for _, m := range MedanWajibBerlaku(h, posisi, tempat) {
		if strings.TrimSpace(h.Ambil(m.Jalur)) == "" {
			kosong = append(kosong, m.Label)
		}
	}
	return kosong
}

// ---------------------------------------------------------------- medan yang boleh diubah

// medanAtasan - medan yang DAPAT DIISI di layar atasan (flow action
// `DeptHeadTreatyIn_UW` -> `Section/GeneralDeptHeadTreatyIn_UW` (nol sel
// sendiri; hanya menyertakan `DetailDeptHeadTreatyIn_UW` atas `.PolicyTreatyIn`)
// -> `DetailDeptHeadTreatyIn_UW` + `ListSuggest`). AC 52, dibaca ulang
// 2026-10-03:
//
//	DetailDeptHeadTreatyIn_UW  sel ber-`pyReadOnly=false` hanya `.DueTo`,
//	                           `.FlagPPH`, `.QuotationData.NoOfferSlip` - ketiganya
//	                           mode sunting `pyDisabled=true`/`pyDisabledNew=always`
//	                           -> TIDAK dapat diisi; sisanya lima pxButton
//	ListSuggest                `.IsApproved` (wajib), `.Suggest` (wajib) - selalu;
//	                           `.ProductionDate` - hanya bila tampil (lihat
//	                           `TanggalProduksiTampil`, peran_tempat.go)
//
// Selain ini terkunci (AC 49-51): nilai kiriman layar untuk medan lain
// DIABAIKAN, nilai tersimpan dipakai.
var medanAtasan = []string{
	HalamanPolis + ".IsApproved",
	HalamanPolis + ".Suggest",
}

// jalurTanggalProduksi - `Section/ListSuggest` `.ProductionDate`.
const jalurTanggalProduksi = HalamanPolis + ".ProductionDate"

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
		"Installment", "IsApproved", "Suggest",
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
//	Admin   `medanAdmin` (daftar izin), beserta daftar `DaftarDariLayar`
//	        (SpreadingRiskList, ListInstallment - baris boleh ditambah/dihapus,
//	        tombol Add/Delete layar admin). ⛔ RALAT K8: TreatyXOLList tidak
//	        tampil di section mana pun - tidak lagi diterima dari layar.
//	Atasan  hanya `medanAtasan`; seluruh daftar terkunci
//
// Kedua layar: `.ProductionDate` (`ListSuggest`) diterima hanya bila tampil -
// `TanggalProduksiTampil` atas IsApproved (sesudah digabung) dan `tempat`
// berperan pelaku (tiket 05).
//
// Halaman Quotation dan TreatyIn tidak pernah diterima dari layar: Quotation
// diisi pilih bisnis dan CheckDataMkt; TreatyIn dibaca dari view.
func GabungMasukanLayar(h, masuk *Halaman, posisi string, tempat map[string]bool) {
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
	} else {
		for j, v := range masuk.Nilai {
			if medanAdmin[j] {
				h.Setel(j, v)
			}
		}
		for _, d := range DaftarDariLayar(h) { // nonprop_layar.go
			if b, ada := masuk.Daftar[d]; ada {
				h.SetelDaftar(d, salinBaris(b))
			}
		}
	}
	if v, ada := masuk.Nilai[jalurTanggalProduksi]; ada && TanggalProduksiTampil(h, tempat) {
		h.Setel(jalurTanggalProduksi, v)
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
