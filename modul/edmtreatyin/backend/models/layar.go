package models

// Untuk apa berkas ini: ATURAN LAYAR - medan wajib, medan yang boleh diubah per posisi, data transform pra dan
// pasca kedua flow action, dan tombol submit. Pola dan sebagian kode: salinan
// `modul/nbtreatyin/backend/models/layar.go` (06-10-2026), isinya menurut section EDM.
//
// Sumber (korpus EDM Treaty In): `Section/GeneralPolicyTreatyInAddendum` -> `Section/DetailPolicyTreatyInAddendum`
// (SATU section untuk semua posisi; beda posisi lewat `pyWorkPage.PositionNote`), `Section/
// DetailPolicyTreatyInAddGeneral` -> `...AddGeneralEditable` (tab Old Data / New Data / Value Difference),
// `Section/DetailPolicyTreatyInPropNewData2` (tab New Data admin - satu-satunya yang dapat diisi),
// `Section/ListSuggestEDM`, `DataTransform/InputPolicyTreatyInAddendum_preAddDT`, `InputPolicyTreatyIn_preAddDT`,
// `DeptHeadTreatyInAddendum_PreDT`, `InboxPolicyTreatyInAddendum_postDT`, `InboxPolicyTreatyIn_UW_postDT`,
// `AddToListCommentsPolicyTreatyIn_DT`.

import (
	"strings"
	"time"

	"nusantarare/inti/backend/utils"
)

// ---------------------------------------------------------------- syarat tampil

// bukanNonProp = `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` (syarat wajib `*` PropNewData2).
func bukanNonProp(h *Halaman) bool {
	return h.Ambil(HalamanQuotation+".ProportionalType") != JenisNonProporsional
}

// flagPPH = `.FlagPPH = true` (sel Type Tax `DetailPolicyTreatyInAddendum`).
func flagPPH(h *Halaman) bool { return h.Ambil(HalamanPolis+".FlagPPH") == "true" }

// bukanNonPropBaru = `.IsNewPolicyNonProp != 1` - wadah S22 `DetailPolicyTreatyInAddGeneral` (tab Old / New /
// Value Difference) dan S17 `DetailPolicyTreatyInAddendum` (`.IsNewPolicyNonProp = 0 || ”`).
func bukanNonPropBaru(h *Halaman) bool { return h.Ambil(HalamanPolis+".IsNewPolicyNonProp") != "1" }

// nonPropBaru = `.IsNewPolicyNonProp = 1` - wadah S11 `DetailPolicyTreatyInAddendum` (AddPremi + Installment).
func nonPropBaru(h *Halaman) bool { return !bukanNonPropBaru(h) }

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

// ---------------------------------------------------------------- medan wajib

// MedanWajib adalah satu medan bertanda wajib di satu layar.
type MedanWajib struct {
	// Jalur halaman (relatif pyWorkPage).
	Jalur string
	// Label - VERBATIM label sel di section.
	Label string
	// Syarat - `pyRequired` bersyarat; nil = selalu wajib.
	Syarat func(h *Halaman) bool
	// NolOtomatis - medan UANG wajib: kosong diisi "0" saat Save / Submit, bukan ditolak (ketetapan NB
	// `[permintaan work owner 06-10-2026]`).
	NolOtomatis bool
	// HanyaSubmit - wajib hanya saat Submit, bukan halangan Save (ketetapan NB: Approval, Suggest).
	HanyaSubmit bool
}

func wajib(m, label string, syarat func(*Halaman) bool) MedanWajib {
	return MedanWajib{Jalur: HalamanPolis + "." + m, Label: label, Syarat: syarat}
}

// wajibUang - medan uang `PropNewData2` bertanda `*`: `pyRequiredWhen` `pyWorkPage.Quotation.ProportionalType !=
// 'NonProportional' && pyWorkPage.PositionNote = 'ReasTreatyInAdmin'`, dan wadahnya tampil (tab New Data =
// `.IsNewPolicyNonProp != 1`).
func wajibUang(m, label string) MedanWajib {
	w := wajib(m, label, dan(bukanNonProp, bukanNonPropBaru))
	w.NolOtomatis = true
	return w
}

// wajibSubmit - medan `Section/ListSuggestEDM` yang wajib (Approval, Suggest) - saat Submit saja.
func wajibSubmit(m, label string) MedanWajib {
	w := wajib(m, label, nil)
	w.HanyaSubmit = true
	return w
}

// medanWajibAdmin - `DetailPolicyTreatyInAddendum` (Type Tax wajib bila `.FlagPPH = true`) + `PropNewData2`
// (medan bertanda `*`, label VERBATIM) + `ListSuggestEDM`.
var medanWajibAdmin = []MedanWajib{
	wajib("TypeTax", "Type Tax", flagPPH),
	wajibUang("PremiOgp", "Premi Ogp"),
	wajibUang("RiCommOgp", "(%) Deduction In A (OGP)"),
	wajibUang("OveriddingCommOgp", "(%) Deduction In B (OGP)"),
	wajibUang("PremiOnp", "Premi Onp"),
	wajibUang("RiCommOnp", "(%) Deduction In A (ONP)"),
	wajibUang("OveriddingCommOnp", "(%) Deduction In B (ONP)"),
	wajibUang("Claim", "Claim"),
	wajibUang("OutstandingClaim", "Outstanding Claim"),
	wajibUang("SalvageValue", "Salvage"),
	wajibUang("ExcessLoss", "Excess Loss"),
	wajibUang("Deduction1", "Deduction1"),
	wajibUang("Deduction2", "Deduction2"),
	wajibSubmit("IsApproved", "Approval"),
	wajibSubmit("Suggest", "Suggest"),
}

// medanWajibAtasan - layar atasan: section yang SAMA, tetapi tab New Data adalah `PropNewData` (semua RO, nol
// `*`); wajib hanya `ListSuggestEDM` - dan Type Tax (sel header yang sama, Auto di XML).
var medanWajibAtasan = []MedanWajib{
	wajib("TypeTax", "Type Tax", flagPPH),
	wajibSubmit("IsApproved", "Approval"),
	wajibSubmit("Suggest", "Suggest"),
}

// DaftarMedanWajib - medan wajib layar posisi itu.
func DaftarMedanWajib(posisi string) []MedanWajib {
	switch posisi {
	case PosisiAdmin:
		return medanWajibAdmin
	case PosisiSecHead, PosisiDeptHead:
		return medanWajibAtasan
	}
	return nil
}

// MedanWajibBerlaku - medan wajib layar posisi itu yang BERLAKU saat ini (syarat terpenuhi). Urutan = urutan
// layar. SATU sumber untuk `Layar.MedanWajib`, tombol Save, dan Submit.
func MedanWajibBerlaku(h *Halaman, posisi string) []MedanWajib {
	var out []MedanWajib
	for _, m := range DaftarMedanWajib(posisi) {
		if m.Syarat == nil || m.Syarat(h) {
			out = append(out, m)
		}
	}
	return out
}

// IsiNolWajibUang - medan UANG wajib yang berlaku dan kosong diisi "0" (ketetapan NB 06-10-2026). Jawab jalur
// yang diisi.
func IsiNolWajibUang(h *Halaman, posisi string) []string {
	var diisi []string
	for _, m := range MedanWajibBerlaku(h, posisi) {
		if m.NolOtomatis && strings.TrimSpace(h.Ambil(m.Jalur)) == "" {
			h.Setel(m.Jalur, "0")
			diisi = append(diisi, m.Jalur)
		}
	}
	return diisi
}

// MedanWajibKosongSimpan - medan wajib kosong untuk tombol Save: medan `HanyaSubmit` tidak menghalangi.
func MedanWajibKosongSimpan(h *Halaman, posisi string) []string {
	var kosong []string
	for _, m := range MedanWajibBerlaku(h, posisi) {
		if !m.HanyaSubmit && strings.TrimSpace(h.Ambil(m.Jalur)) == "" {
			kosong = append(kosong, m.Label)
		}
	}
	return kosong
}

// MedanWajibKosong - label medan wajib berlaku yang kosong (Submit).
func MedanWajibKosong(h *Halaman, posisi string) []string {
	var kosong []string
	for _, m := range MedanWajibBerlaku(h, posisi) {
		if strings.TrimSpace(h.Ambil(m.Jalur)) == "" {
			kosong = append(kosong, m.Label)
		}
	}
	return kosong
}

// ---------------------------------------------------------------- medan yang boleh diubah

// pemicu - action set sel isian yang dijalankan server bila kiriman sel itu berubah.
type pemicu int

const (
	// tanpaPemicu - nilai nol: sel tanpa action set (kepala iota, dipakai sebagai bawaan medanIsian.pemicu).
	tanpaPemicu pemicu = iota
	// pemicuUang - sel uang PropNewData2: change -> `CountOGPONP_Act` (langkah 9 CountSpreading_Act, 10
	// SetValidateInstallment_Act) lalu refresh otherSection `DetailPolicyTreatyInPropValueDifference` - section
	// itu ber-defer-load `EDMTCalculateTreatyDifference` (selisih dihitung ulang).
	pemicuUang
	// pemicuAngsuran - sel `.Installment` PropNewData2: change -> `FillPaymentInstallment`.
	pemicuAngsuran
	// pemicuAngsuranEDMT - sel `.Installment` S12 `DetailPolicyTreatyInAddendum` (NonProp baru): change -> refresh
	// `FillPaymentInstallmentEDMT`.
	pemicuAngsuranEDMT
	// pemicuHapusTypeTax - sel `.FlagPPH` ("With Tax"): click -> runActivity `RemoveTypeTax_ACT`.
	pemicuHapusTypeTax
	// pemicuPajak - sel `.TypeTax`: XML hanya postValue. Ketetapan NB `[keputusan work owner 06-10-2026]` pajak
	// dihitung ulang saat diubah.
	pemicuPajak
)

// medanIsian - satu sel isian TERBUKA: jalur relatif `PolicyTreatyIn`, syarat tampil (nil = selalu), action set.
type medanIsian struct {
	medan  string
	tampil func(*Halaman) bool
	pemicu pemicu
}

// medanHeader - sel header `DetailPolicyTreatyInAddendum` yang `pyEditOptions=Auto`. XML tidak menguncinya per
// posisi (section yang sama dipakai flow action atasan); ⛔ keputusan work owner 07-10-2026: atasan TIDAK boleh
// mengubahnya - hanya diterima dari layar Admin (penyimpangan sadar). Overiding Commision (`.FlagRetroTreaty`), With Tax
// (`.FlagPPH` click -> RemoveTypeTax_ACT), Type Tax (`.TypeTax`, tampil `.FlagPPH = true`), Marketing Officer
// (`.QuotationData.MOID`, change -> CheckDataMkt - dijalankan lewat aksi hitung).
var medanHeader = []medanIsian{
	{medan: "FlagRetroTreaty"},
	{medan: "FlagPPH", pemicu: pemicuHapusTypeTax},
	{medan: "TypeTax", tampil: flagPPH, pemicu: pemicuPajak},
	{medan: "QuotationData.MOID"},
}

// medanSaran - `Section/ListSuggestEDM`: Approval (`.IsApproved`, radio) dan Suggest (`.Suggest`). Production
// Date ListSuggest NB TIDAK ada di ListSuggestEDM.
var medanSaran = []medanIsian{{medan: "IsApproved"}, {medan: "Suggest"}}

// medanAdminSaja - DAFTAR IZIN tambahan layar admin (berurutan; syarat tampil membaca medan di atasnya):
//
//	S10   `.Remark` pxTextArea - disabled bila `PositionNote != ADM`
//	S12   `.Installment` (S11 `.IsNewPolicyNonProp = 1`) - change -> FillPaymentInstallmentEDMT
//	NewData2 (tab New Data admin; wadah `.IsNewPolicyNonProp != 1`): GrossPremium (tanpa action set),
//	      OutstandingClaim (tanpa action set), .Installment (change -> FillPaymentInstallment), lima belas sel uang
//	      ber-CountOGPONP_Act (+ refresh Value Difference)
//
// Medan TURUNAN (NetPremium, Balance*, PPN/PPH, BrokerageFee*, DueTo, Currency) TIDAK diterima dari layar.
var medanAdminSaja = func() []medanIsian {
	m := []medanIsian{
		{medan: "Remark"},
		{medan: "Installment", tampil: nonPropBaru, pemicu: pemicuAngsuranEDMT},
		{medan: "GrossPremium", tampil: bukanNonPropBaru},
		{medan: "OutstandingClaim", tampil: bukanNonPropBaru},
	}
	// .Installment tab New Data (wadah berjudul "hidden" tanpa syarat di PropNewData2) - sel yang sama dengan S12
	// tetapi di wadah berlawanan: hanya satu yang tampil.
	m = append(m, medanIsian{medan: "Installment", tampil: bukanNonPropBaru, pemicu: pemicuAngsuran})
	for _, n := range []string{"PremiOgp", "RiCommOgp", "ResultOgp1", "OveriddingCommOgp", "ResultOgp2",
		"PremiOnp", "RiCommOnp", "ResultOnp1", "OveriddingCommOnp", "ResultOnp2",
		"Claim", "SalvageValue", "ExcessLoss", "Deduction1", "Deduction2"} {
		m = append(m, medanIsian{medan: n, tampil: bukanNonPropBaru, pemicu: pemicuUang})
	}
	return m
}()

// PemicuLayar - action set sel layar yang DIPICU perubahan kiriman, diputar ulang di server oleh
// `TerapkanPemicu`.
type PemicuLayar struct {
	// Spreading - %Share / %Share Claim spreading berubah, atau baris ditambah: refresh `CountSpreading_Act`.
	Spreading bool
	// Uang - sel uang ber-action set `CountOGPONP_Act` berubah.
	Uang bool
	// Angsuran - `.Installment` tab New Data berubah: `FillPaymentInstallment`.
	Angsuran bool
	// AngsuranEDMT - `.Installment` S12 (NonProp baru) berubah: `FillPaymentInstallmentEDMT`.
	AngsuranEDMT bool
	// Pajak - `.FlagPPH` / `.TypeTax` berubah.
	Pajak bool
	// Selisih - tab Value Difference disegarkan (defer-load `EDMTCalculateTreatyDifference`): setiap sel uang /
	// spreading / angsuran tab New Data yang berubah.
	Selisih bool
}

// TerapkanNilaiBawaanSel = `pyDefaultValue` sel terbuka `DetailPolicyTreatyInAddendum`: Type Tax "Inclusive"
// (pola NB; sel `.TypeTax` salinan sel NB yang sama).
func TerapkanNilaiBawaanSel(h *Halaman) {
	if flagPPH(h) && h.Ambil(pt+"TypeTax") == "" {
		h.Setel(pt+"TypeTax", TypeTaxInclusive)
	}
}

// GabungMasukanLayar menyalin nilai kiriman layar `masuk` ke halaman tersimpan `h`, HANYA untuk medan yang
// boleh diubah di posisi itu, dan menjawab action set yang dipicu kiriman.
//
//	Semua posisi  `medanSaran`
//	Admin         + `medanHeader` (keputusan WO 07-10-2026: atasan tidak boleh mengubah) + `medanAdminSaja` (Remark,
//	              Installment, tab New Data PropNewData2) + grid spreading tab New Data
//
// Halaman Quotation, OldData, TreatyDifference, TreatyXOLList / TreatyXOLDifferenceList (grid AddPremi: panel
// rinciannya `DetailPolicyAddPremiDetail` seluruhnya RO), TreatyIn tidak pernah diterima dari layar.
func GabungMasukanLayar(h, masuk *Halaman, posisi string) PemicuLayar {
	var p PemicuLayar
	if masuk == nil {
		return p
	}
	terima := func(daftar []medanIsian) {
		for _, m := range daftar {
			j := pt + m.medan
			v, ada := masuk.Nilai[j]
			if !ada || (m.tampil != nil && !m.tampil(h)) {
				continue
			}
			berubah := nilaiBerubah(h.Ambil(j), v)
			h.Setel(j, v)
			if !berubah {
				continue
			}
			switch m.pemicu {
			case pemicuUang:
				p.Uang, p.Selisih = true, true
			case pemicuAngsuran:
				p.Angsuran, p.Selisih = true, true
			case pemicuAngsuranEDMT:
				p.AngsuranEDMT = true
			case pemicuHapusTypeTax:
				RemoveTypeTax(h)
				p.Pajak = true
			case pemicuPajak:
				p.Pajak = true
			}
		}
	}
	if posisi == PosisiAdmin {
		terima(medanHeader)
		terima(medanAdminSaja)
		if bukanNonPropBaru(h) {
			if b, ada := masuk.Daftar[DaftarSpreading]; ada && gabungSpreading(h, b) {
				p.Spreading, p.Selisih = true, true
			}
		}
	}
	terima(medanSaran)
	// pajak proporsional: CountNetPremi_act dijalankan services.turunkan; daftar yang memakai BalanceDueTo
	// diputar ulang seperti sel uang berubah.
	if p.Pajak && bukanNonPropBaru(h) {
		p.Uang, p.Selisih = true, true
	}
	TerapkanNilaiBawaanSel(h)
	return p
}

// kolomHitungSpreading - kolom `Read-only` grid spreading (`.PremiumSpreaded`, `.ClaimSpreaded`) yang hanya
// ditulis `CountSpreading_Act`.
var kolomHitungSpreading = []string{"PremiumSpreaded", "ClaimSpreaded"}

// selSpreading - sel grid spreading tab New Data (`PropNewData2`) yang diterima dari layar. Keputusan work owner
// 07-10-2026 (screenshot grid, "INI READ ONLY JUGA" - sama dengan NB 06-10-2026): tombol Add DIBUANG dan `.TreatyType`
// hanya-baca; yang tersisa `.SharePercentage` dan `.ClaimPercentage` (change -> `CountSpreading_Act`). Di XML
// `PropNewData2` juga `.TreatyType` pxDropdown + tombol Add - penyimpangan sadar.
var selSpreading = []string{"SharePercentage", "ClaimPercentage"}

// gabungSpreading menerima kiriman grid spreading tab New Data.
//
//   - baris kiriman ke-i = baris server ke-i (pasangan menurut posisi, `.pxListSubscript`); hanya `selSpreading`;
//   - ⛔ `[keputusan work owner]` spec-penyimpanan ID-16: di endorsemen baris TIDAK DAPAT DIHAPUS - baris yang
//     hilang dari kiriman TETAP ada (tombol Delete `PropNewData2` tidak dibangun);
//   - baris kiriman di belakang baris server DIABAIKAN (Add dibuang, keputusan WO 07-10-2026).
//
// Jawab apakah `CountSpreading_Act` terpicu (%Share / %Claim berubah): kolom hitung dikosongkan, `TerapkanPemicu`
// menghitungnya ulang.
func gabungSpreading(h *Halaman, kiriman []Baris) bool {
	lama := h.AmbilDaftar(DaftarSpreading)
	baru := make([]Baris, 0, len(lama))
	terpicu := false
	for i, l := range lama {
		r := Baris{}
		for k, v := range l {
			r[k] = v
		}
		if i < len(kiriman) {
			for _, m := range selSpreading {
				if v, ada := kiriman[i][m]; ada {
					terpicu = terpicu || nilaiBerubah(l[m], v)
					r[m] = v
				}
			}
		}
		baru = append(baru, r)
	}
	if terpicu {
		for _, r := range baru {
			for _, k := range kolomHitungSpreading {
				delete(r, k)
			}
		}
	}
	h.SetelDaftar(DaftarSpreading, baru)
	return terpicu
}

// nilaiBerubah - isian sel berubah = event `change` sel Pega (teks berbeda).
func nilaiBerubah(lama, baru string) bool {
	return strings.TrimSpace(lama) != strings.TrimSpace(baru)
}

// TerapkanPemicu memutar ulang di server action set sel yang dipicu kiriman layar, atas halaman yang medan
// turunannya sudah dihitung (`services.turunkan`):
//
//	Angsuran      `FillPaymentInstallment` (baris disusun ulang)
//	AngsuranEDMT  `FillPaymentInstallmentEDMT` (dari TreatyXOLDifferenceList, Param.Installment = .Installment)
//	Uang          `CountOGPONP_Act` langkah 10 `SetValidateInstallment_Act` dan langkah 9 `CountSpreading_Act`
//	Spreading     `CountSpreading_Act`
//	Selisih       refresh `DetailPolicyTreatyInPropValueDifference` -> `EDMTCalculateTreatyDifference`
func TerapkanPemicu(h *Halaman, p PemicuLayar, sekarang time.Time) error {
	switch {
	case p.Angsuran:
		if err := FillPaymentInstallment(h, sekarang); err != nil {
			return err
		}
	case p.Uang:
		if err := SetValidateInstallment(h); err != nil {
			return err
		}
	}
	if p.AngsuranEDMT {
		if err := FillPaymentInstallmentEDMT(h, h.Ambil(pt+"Installment")); err != nil {
			return err
		}
	}
	if p.Uang || p.Spreading {
		if err := CountSpreading(h, 0); err != nil {
			return err
		}
	}
	if err := HitungTotalSpreading(h); err != nil {
		return err
	}
	if p.Selisih {
		return EDMTCalculateTreatyDifference(h)
	}
	return nil
}

// ---------------------------------------------------------------- pra flow action

// PraprosesAdmin = `DataTransform/InputPolicyTreatyInAddendum_preAddDT` (pra-DT flow action
// `InboxPolicyTreatyInAddendum`):
//
//	1    APPLY `InputPolicyTreatyIn_preAddDT`:
//	     1-3  StartDate / EndDate / StatementDate = hari ini bila kosong
//	     4-5  IsNewPolicyNonProp = "1" bila Quotation.ProportionalType == "NonProportional", selain itu "0"
//	     6    MasterID = ""
//	     7    MarketingOfficer = QuotationData.MarketingName
//	     8-9  IsApproved = "", Suggest = ""
//	     10   SuggestDate = sekarang ; 11 OperatorName = OperatorID.pyUserIdentifier ; 12 FlagOnGoingPolicy "1"
//	          (⛔ OperatorName = NAMA TAMPILAN - ketetapan NB P33: pengenal akun di sini bug, bukan maksud)
//	     13   TempEmail.CARI28 = pyWorkBasketList(2) - penentu tombol admin aktif: diganti keanggotaan antrean
//	2    ViewState = "0" - keadaan layar, tidak disimpan
//	3    PositionNote = ReasTreatyInAdmin ; 4 Show = True (keadaan layar) ; 5 Position = "5"
func PraprosesAdmin(h *Halaman, sekarang time.Time, namaTampilan string) {
	if h.Ambil(pt+"StartDate") == "" {
		h.Setel(pt+"StartDate", utils.FormatTanggal(sekarang))
	}
	if h.Ambil(pt+"EndDate") == "" {
		h.Setel(pt+"EndDate", utils.FormatTanggal(sekarang))
	}
	if h.Ambil(pt+"StatementDate") == "" {
		h.Setel(pt+"StatementDate", utils.FormatTanggalWaktu(sekarang))
	}
	if h.Ambil(HalamanQuotation+".ProportionalType") == JenisNonProporsional {
		h.Setel(pt+"IsNewPolicyNonProp", "1")
	} else {
		h.Setel(pt+"IsNewPolicyNonProp", "0")
	}
	h.Setel(pt+"MasterID", "")
	h.Setel(pt+"MarketingOfficer", h.Ambil(pt+"QuotationData.MarketingName"))
	h.Setel(pt+"IsApproved", "")
	h.Setel(pt+"Suggest", "")
	h.Setel(pt+"SuggestDate", utils.FormatTanggalWaktu(sekarang))
	h.Setel(pt+"OperatorName", namaTampilan)
	h.Setel("FlagOnGoingPolicy", "1")
	h.Setel("PositionNote", PosisiAdmin)
	h.Setel("Position", PositionAtasan)
}

// PraprosesAtasan = `DataTransform/DeptHeadTreatyInAddendum_PreDT`: ViewState "1" (keadaan layar),
// OperatorName = OperatorID.pyUserIdentifier (⛔ nama tampilan - ketetapan NB P33), Suggest = "", SuggestDate =
// sekarang, IsApproved = "".
func PraprosesAtasan(h *Halaman, sekarang time.Time, namaTampilan string) {
	h.Setel(pt+"OperatorName", namaTampilan)
	h.Setel(pt+"Suggest", "")
	h.Setel(pt+"SuggestDate", utils.FormatTanggalWaktu(sekarang))
	h.Setel(pt+"IsApproved", "")
}

// ---------------------------------------------------------------- pasca submit

// TambahCatatan = `DataTransform/AddToListCommentsPolicyTreatyIn_DT`, dipanggil kedua DT pasca dengan
// Comment=.Suggest, Operator=.OperatorName, Approved=.IsApproved, Date=.SuggestDate:
//
//	APPEND .PolicyTreatyIn.SuggestList: .Suggest, WHEN Param.Approved != "" -> .IsApproved, .Date, .OperatorName
//
// `operatorID` (identitas akses login) ditambahkan supaya catatan tersimpan bersama operatornya (P4).
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

// PascaAdmin = `DataTransform/InboxPolicyTreatyInAddendum_postDT`:
//
//	1   AddToListComments -> TambahCatatan
//	2   PositionNote == Admin && IsApproved == "0" -> NBStatus = "EDM WAS DECLINED BY  " + @toUpperCase(pyUserName)
//	3   PositionNote == Admin && IsApproved == "1" -> NBStatus = "EDM IS IN <nama tertanam>'S INBOX" - ditimpa
//	    connector ke Sec Head sesudahnya (services, nama dari data) - tidak ditulis di sini
func PascaAdmin(h *Halaman, operatorID, namaTampilan string) {
	TambahCatatan(h, operatorID)
	if h.Ambil("PositionNote") == PosisiAdmin && h.Ambil(HalamanPolis+".IsApproved") == NilaiDitolak {
		h.Setel("NBStatus", TeksNBStatusDitolakOleh(namaTampilan))
	}
}

// PascaAtasan = `DataTransform/InboxPolicyTreatyIn_UW_postDT` langkah 1 (AddToListComments). Langkah 2-5 hanya
// untuk PositionNote `ReasTreatyInDirector` / `ReasTreatyInGroupLeader` - posisi yang tidak ada di flow EDM, mati.
func PascaAtasan(h *Halaman, operatorID string) {
	TambahCatatan(h, operatorID)
}

// ---------------------------------------------------------------- tombol layar

// TombolKirim adalah tombol submit yang tampil di layar (`DetailPolicyTreatyInAddendum` S18 / S19).
type TombolKirim string

const (
	// TombolTidakAda - IsApproved belum bernilai yang memunculkan tombol.
	TombolTidakAda TombolKirim = ""
	// TombolKirimLangsung - `finishAssignment` (admin IsApproved==1 didahului SetDueTo_act; atasan
	// IsApproved='0'; Sec Head IsApproved==1).
	TombolKirimLangsung TombolKirim = "kirim"
	// TombolKonfirmasiTolak - admin IsApproved==0: localAction `PolicyTreatyInDeclineConfirm` (Yes = finish).
	TombolKonfirmasiTolak TombolKirim = "konfirmasi-tolak"
	// TombolNomorPolis - Dept Head IsApproved==1: localAction `ShowPolicyNoTreaty` (PolicyNo sudah terisi sejak
	// `SetEDMTNoPolis`, sehingga cabang `GeneratePolicyNoTreatyAddendum_Act` bila `PolicyNo = ""` tidak pernah
	// jalan); OK = closeContainer + finishAssignment.
	TombolNomorPolis TombolKirim = "nomor-polis"
)

// TombolUntuk menentukan tombol submit yang tampil.
//
// ⛔ Syarat identitas orang tombol atasan S18 (`OperatorID.pyUserIdentifier == <ID operator #1 / #2>`, dua ID
// operator tertanam) DIGANTI posisi kasus (pola NB): Dept Head - jenjang terakhir - menampilkan popup nomor polis;
// Sec Head menaikkan (`LetterNo == 'TREATYINDEPTHEAD'`, varian B - tangga tiga jenjang).
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
