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

import (
	"strings"
	"time"
)

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

// pemicu - action set sel isian admin yang dijalankan server bila kiriman sel
// itu berubah: yang menulis DAFTAR diputar ulang `TerapkanPemicu`; yang menulis
// medan halaman dijalankan saat digabung.
type pemicu int

const (
	tanpaPemicu pemicu = iota
	// pemicuUang - change -> refresh `CountOGPONP_Act` (langkah 9 CountSpreading_Act
	// per baris spreading, 10 SetValidateInstallment_Act).
	pemicuUang
	// pemicuAngsuran - change -> refresh `FillPaymentInstallment(Installment=.Installment)`.
	pemicuAngsuran
	// pemicuHapusTypeTax - sel `.FlagPPH`: change -> runActivity `RemoveTypeTax_ACT`
	// (menulis medan halaman, bukan daftar - dijalankan saat digabung).
	pemicuHapusTypeTax
)

// medanIsianAdmin - satu sel isian TERBUKA layar admin: jalur relatif
// `PolicyTreatyIn`, syarat tampil (`pyVisible` sel DAN `pyContainerVisibleWhen`
// wadahnya; nil = selalu), dan action set yang menulis daftar.
type medanIsianAdmin struct {
	medan  string
	tampil func(*Halaman) bool
	pemicu pemicu
}

// medanAdmin - DAFTAR IZIN layar admin (`DetailPolicyTreatyIn` + `ListSuggest`),
// medan yang dapat diketik/dipilih (kolom "Kunci" kosong atau bersyarat yang tidak
// pernah benar bagi kasus treaty - IsUW), BERURUTAN: syarat tampil sebuah medan
// hanya membaca medan di atasnya (ClaimType -> FlagPPH -> TypeTax), dinilai atas
// halaman yang sedang digabung - persis sel yang tampil di layar saat itu.
//
// W4 audit silang P3: medan di sel / wadah TERSEMBUNYI tidak diterima (sel
// tersembunyi tidak ter-render dan tidak pernah mengirim nilai):
//
//	S19  `.IsNewPolicyNonProp != 1 && .IsNewPolicyListFormat != 1`  medan uang, .Installment
//	     (polis NonProp baru: PremiOgp, Deduction1/2 milik
//	     InputPolicyTreatyInDetail_NonProp 18-19, bukan layar)
//	S7   `.ClaimType != 'XOL Retro'`  .FlagPPH, .TypeTax (sel: `.FlagPPH = true`)
//	sel  `.ClaimType != 'XOL Retro'`  .FlagRetroTreaty
//	S14  `.QuotationData.ProportionalType = 'Proportional'`  .Quartal, .YearOfQuartal
//	sel  `.IsNewPolicyNonProp != 1`  .IDCurrency
//	sel  `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`  .QuotationData.IsSurveyReport
//
// Medan TURUNAN (NetPremium, Balance*, PPN/PPH, BrokerageFee*, DueTo,
// Currency) TIDAK diterima dari layar - services menghitungnya ulang dari
// isian (`services.turunkan`). Hasil tombol Enable / Disable Input Type
// (`IsNewPolicyNonProp`, `QuotationData.ProportionalType`): `KirimanEnableDisable`.
//
// ⛔ DAFTAR IZIN, bukan daftar larangan (temuan tinjauan 2026-10-03): medan
// milik server - hasil pilih bisnis (NoOffer, TreatyGroupID, OJKBusinessID,
// BizCode, SOB, CedingCo, ...), StatementDate (`ALWAYS`), PolicyNo - tidak
// pernah dapat ditimpa layar, termasuk yang tidak tampil.
var medanAdmin = func() []medanIsianAdmin {
	m := []medanIsianAdmin{}
	for _, n := range []string{"StartDate", "EndDate", "StatementType", "QuotationData.NoOfferSlip",
		"QuotationData.MOID", "ClaimType", "ClaimPaymentType", "Remark", "IsApproved", "Suggest"} {
		m = append(m, medanIsianAdmin{medan: n})
	}
	m = append(m,
		medanIsianAdmin{medan: "QuotationData.IsSurveyReport", tampil: bukanNonProp},
		medanIsianAdmin{medan: "FlagRetroTreaty", tampil: bukanXOLRetro},
		medanIsianAdmin{medan: "FlagPPH", tampil: bukanXOLRetro, pemicu: pemicuHapusTypeTax},
		medanIsianAdmin{medan: "TypeTax", tampil: dan(bukanXOLRetro, flagPPH)},
		medanIsianAdmin{medan: "Quartal", tampil: proporsionalQD},
		medanIsianAdmin{medan: "YearOfQuartal", tampil: proporsionalQD},
		medanIsianAdmin{medan: "IDCurrency", tampil: bukanNonPropBaru},
		// GrossPremium/GrossClaim: CalculatePremi_Act -> CountNetPremi_act (tanpa daftar)
		medanIsianAdmin{medan: "GrossPremium", tampil: wadahUangAdmin},
		medanIsianAdmin{medan: "GrossClaim", tampil: wadahUangAdmin},
		// OutstandingClaim: sel tanpa action set
		medanIsianAdmin{medan: "OutstandingClaim", tampil: wadahUangAdmin},
		medanIsianAdmin{medan: "Installment", tampil: wadahUangAdmin, pemicu: pemicuAngsuran},
	)
	for _, n := range []string{"PremiOgp", "RiCommOgp", "ResultOgp1", "OveriddingCommOgp", "ResultOgp2",
		"PremiOnp", "RiCommOnp", "ResultOnp1", "OveriddingCommOnp", "ResultOnp2",
		"Claim", "SalvageValue", "ExcessLoss", "Deduction1", "Deduction2"} {
		m = append(m, medanIsianAdmin{medan: n, tampil: wadahUangAdmin, pemicu: pemicuUang})
	}
	return m
}()

// HasilEnableDisable - dua medan yang hanya ditulis `TreatyEnableDisableInput`.
type HasilEnableDisable struct {
	IsNewPolicyNonProp string
	ProportionalType   string
}

const (
	jalurIsNewNonProp    = pt + "IsNewPolicyNonProp"
	jalurJenisProporsiQD = pt + "QuotationData.ProportionalType"
)

func bacaEnableDisable(h *Halaman) HasilEnableDisable {
	return HasilEnableDisable{IsNewPolicyNonProp: h.Ambil(jalurIsNewNonProp), ProportionalType: h.Ambil(jalurJenisProporsiQD)}
}

func tulisEnableDisable(h *Halaman, x HasilEnableDisable) {
	h.Setel(jalurIsNewNonProp, x.IsNewPolicyNonProp)
	h.Setel(jalurJenisProporsiQD, x.ProportionalType)
}

// PesanEnableDisableTidakCocok - pesan 422 W3.
func PesanEnableDisableTidakCocok(x HasilEnableDisable) string {
	return "Nilai " + jalurIsNewNonProp + " \"" + x.IsNewPolicyNonProp + "\" / " + jalurJenisProporsiQD + " \"" +
		x.ProportionalType + "\" bukan hasil tombol Enable / Disable Input Type " +
		"(TreatyEnableDisableInput) - medan ini terkunci di layar"
}

// KirimanEnableDisable = W3 audit silang P3, pola F4 (`TerimaKirimanTerkunci`):
// `.QuotationData.ProportionalType` (sel `pyReadOnly=true`,
// `pyEditOptions=Read-only`) dan `.IsNewPolicyNonProp` (bukan sel layar) hanya
// diubah DataTransform `TreatyEnableDisableInput` - tombol "Enable / Disable
// Input Type", pyVisible `.TreatyType='XOL'`, click -> refresh (TANPA simpan;
// hasilnya dipegang layar sampai Save/Submit). Hasil hitung ulang = DT atas
// halaman server (langkah 1 "NonProportional", langkah 2-4 selalu "0").
func KirimanEnableDisable(h *Halaman) KirimanTerkunci[HasilEnableDisable] {
	return KirimanTerkunci[HasilEnableDisable]{
		Jalur:  []string{jalurIsNewNonProp, jalurJenisProporsiQD},
		Tampil: h.Ambil(pt+"TreatyType") == "XOL",
		Baca:   bacaEnableDisable,
		Tulis:  tulisEnableDisable,
		Pesan:  PesanEnableDisableTidakCocok,
		HitungUlang: func() ([]HasilEnableDisable, error) {
			dt := h.Salin()
			TreatyEnableDisableInput(dt)
			return []HasilEnableDisable{bacaEnableDisable(dt)}, nil
		},
	}
}

// TerapkanNilaiBawaanSel = `pyDefaultValue` sel terbuka `Section/DetailPolicyTreatyIn`
// (juga salinannya `GeneralPolicyTreatyIn`): nilai yang dipakai bila medan
// kosong saat sel DIRENDER - hanya bila selnya tampil. Dua satu-satunya di layar
// NB (sisanya label mati `1=2` dan grid master baca-saja):
//
//	.QuotationData.IsSurveyReport  "No"         pyCondition `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`
//	.TypeTax                       "Inclusive"  wadah S7 `.ClaimType != 'XOL Retro'`, pyCondition `.FlagPPH = true`
//
// Layar atasan (`DetailDeptHeadTreatyIn_UW`) tidak memuat `pyDefaultValue` sel
// terbuka. Dipanggil saat layar admin dirender (`services.BukaKasus`) dan
// sesudah kiriman admin digabung (sel baru tampil, mis. FlagPPH dicentang).
func TerapkanNilaiBawaanSel(h *Halaman) {
	if bukanNonProp(h) && h.Ambil(pt+"QuotationData.IsSurveyReport") == "" {
		h.Setel(pt+"QuotationData.IsSurveyReport", "No")
	}
	if bukanXOLRetro(h) && flagPPH(h) && h.Ambil(pt+"TypeTax") == "" {
		h.Setel(pt+"TypeTax", TypeTaxInclusive)
	}
}

// GabungMasukanLayar menyalin nilai kiriman layar `masuk` ke halaman
// tersimpan `h`, HANYA untuk medan yang boleh diubah di posisi itu, dan
// menjawab action set yang dipicu kiriman (`PemicuLayar`).
//
//	Admin   `KirimanEnableDisable` (W3, pola F4), lalu `medanAdmin` (daftar izin
//	        berurutan, hanya sel yang TAMPIL - W4); `.FlagPPH` yang berubah
//	        menjalankan `RemoveTypeTax_ACT` (change -> runActivity); nilai bawaan
//	        sel (`TerapkanNilaiBawaanSel`).
//	Atasan  hanya `medanAtasan`
//	Daftar  `.SpreadingRiskList` bila gridnya terbuka (`SpreadingDariLayar`,
//	        kedua posisi - W2), TANPA kolom hanya-baca (W5). `.ListInstallment`
//	        (grid `readOnly`) TIDAK PERNAH dari layar - barisnya ditulis action
//	        set server (`TerapkanPemicu`). ⛔ RALAT K8: TreatyXOLList tidak tampil
//	        di section mana pun - tidak diterima dari layar.
//
// Kedua layar: `.ProductionDate` (`ListSuggest`) diterima hanya bila tampil -
// `TanggalProduksiTampil` atas IsApproved (sesudah digabung) dan `tempat`
// berperan pelaku (tiket 05).
//
// Halaman Quotation dan TreatyIn tidak pernah diterima dari layar DI SINI:
// Quotation diisi pilih bisnis dan CheckDataMkt; TreatyIn dibaca dari view.
// Satu-satunya kekecualian - keempat medan `SearchHierarkiSourceBizAgent_PostDT`
// (pemilih Source Of Business, F4) - diterima terpisah oleh
// `services.terimaSumberBisnis`, hanya bila ClaimType 'XOL Retro' dan cocok
// dengan RD `BrowseAgentHierarkiList_RD` yang dijalankan ulang.
func GabungMasukanLayar(h, masuk *Halaman, posisi string, tempat map[string]bool) (PemicuLayar, error) {
	var p PemicuLayar
	if masuk == nil {
		return p, nil
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
		if _, err := TerimaKirimanTerkunci(h, masuk, KirimanEnableDisable(h)); err != nil {
			return p, err
		}
		for _, m := range medanAdmin {
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
				p.Uang = true
			case pemicuAngsuran:
				p.Angsuran = true
			case pemicuHapusTypeTax:
				RemoveTypeTax(h) // change -> runActivity RemoveTypeTax_ACT
			}
		}
		TerapkanNilaiBawaanSel(h)
	}
	if SpreadingDariLayar(h, posisi) { // nonprop_layar.go
		if b, ada := masuk.Daftar[DaftarSpreading]; ada {
			p.Spreading = gabungSpreading(h, b)
		}
	}
	if v, ada := masuk.Nilai[jalurTanggalProduksi]; ada && TanggalProduksiTampil(h, tempat) {
		h.Setel(jalurTanggalProduksi, v)
	}
	return p, nil
}

// PemicuLayar - action set sel layar yang DIPICU perubahan kiriman (nilai
// kiriman yang diterima berbeda dari halaman server), diputar ulang di server
// oleh `TerapkanPemicu` sesudah medan turunan dihitung (W5 audit silang P3:
// kolom hanya-baca daftar tidak pernah diterima dari layar).
type PemicuLayar struct {
	// Spreading - sel `.SharePercentage` / `.ClaimPercentage` grid spreading
	// berubah, baris baru membawa %Share, atau baris dihapus: refresh
	// `CountSpreading_Act`.
	Spreading bool
	// Uang - sel uang ber-action set `CountOGPONP_Act` berubah (admin).
	Uang bool
	// Angsuran - sel `.Installment` berubah: refresh `FillPaymentInstallment`.
	Angsuran bool
}

// kolomHitungSpreading - kolom `Read-only` grid spreading yang hanya ditulis
// `CountSpreading_Act` langkah 4.1 (sel C[2.3]/C[2.5] `DetailPolicyTreatyIn` S30,
// `SpreadingRiskList`, `DetailDeptHeadTreatyIn_UW` S96).
var kolomHitungSpreading = []string{"PremiumSpreaded", "ClaimSpreaded"}

// selSpreading - satu-satunya anggota baris spreading yang berupa SEL terbuka
// grid (`DetailPolicyTreatyIn` S30, `SpreadingRiskList`): `.TreatyType`
// (dropdown, `pyValue .TreatyType`; `.TreatyName` hanya `pyPrompt` daftar
// pilihan - tidak ditulis sel), `.SharePercentage`, `.ClaimPercentage`.
var selSpreading = []string{"TreatyType", "SharePercentage", "ClaimPercentage"}

// anggotaServerSpreading - anggota baris spreading yang BUKAN sel grid dan
// tidak ditulis action set sel mana pun (`CountSpreading_Act` hanya menulis
// SharePercentage/ClaimPercentage/PremiumSpreaded/ClaimSpreaded). Penulisnya
// hanya aktivitas server saat pilih bisnis:
//
//	TreatyName        `TreatyInputPctCommSpreading` (baris 1 = `.SpreadingType` master)
//	Currency, CurrencyID, SplitRNMSharePct
//	                  `TreatyNonPropSetSpreading` langkah 3, 4.1, 6
//	                  (`InputPolicyTreatyInDetail_NonProp` 23)
//
// Pengerasan tercatat (PERMINTAAN H2): nilainya TIDAK diterima dari layar.
// Baris kiriman membawa nilai baris server apa adanya (layar menyalin baris
// utuh); nilai itu dipakai hanya bila SAMA dengan baris server yang belum
// terpakai, berurutan (`pasangAnggotaServer`) - sehingga Delete tetap membawa
// nilai barisnya sendiri. Baris baru (Add membuat baris kosong) atau nilai
// yang tidak ada di server: kosong.
var anggotaServerSpreading = []string{"TreatyName", "Currency", "CurrencyID", "SplitRNMSharePct"}

// gabungSpreading menerima baris grid spreading kiriman layar `kiriman`
// (Add/Delete, sel `selSpreading`) TANPA kolom hanya-baca dan tanpa
// `anggotaServerSpreading`, dan menjawab apakah `CountSpreading_Act` terpicu.
//
//   - terpicu (sel %Share / %Share Claim berubah dibanding baris server di
//     urutan yang sama, baris baru membawa %Share, atau baris dihapus - urutan
//     bergeser): kolom hanya-baca dikosongkan, `TerapkanPemicu` menghitungnya
//     (langkah 4.1 seluruh baris);
//   - tidak terpicu (hanya `.TreatyType` berubah, atau baris baru tanpa
//     %Share): kolom hanya-baca = nilai server di urutan yang sama; baris baru
//     kosong - di Pega pun baris itu kosong sampai refresh berikutnya.
//
// ⚠️ `[penyesuaian sadar]` Baris yang dihapus TANPA perubahan %Share dihitung
// ulang: baris tak berkunci, nilai server tidak dapat dipasangkan lagi
// menurut urutan. Hasilnya sama dengan nilai server selama NetPremium / klaim
// belum berubah sejak refresh terakhir.
func gabungSpreading(h *Halaman, kiriman []Baris) bool {
	lama := h.AmbilDaftar(DaftarSpreading)
	baru := make([]Baris, len(kiriman))
	for i, k := range kiriman {
		baru[i] = Baris{}
		for _, m := range selSpreading {
			if v, ada := k[m]; ada {
				baru[i][m] = v
			}
		}
	}
	pasangAnggotaServer(lama, kiriman, baru)
	terpicu := len(baru) < len(lama)
	for i, b := range baru {
		if i < len(lama) {
			terpicu = terpicu || nilaiBerubah(lama[i]["SharePercentage"], b["SharePercentage"]) ||
				nilaiBerubah(lama[i]["ClaimPercentage"], b["ClaimPercentage"])
		} else {
			terpicu = terpicu || b["SharePercentage"] != "" || b["ClaimPercentage"] != ""
		}
	}
	if !terpicu {
		for i := 0; i < len(baru) && i < len(lama); i++ {
			for _, k := range kolomHitungSpreading {
				if v, ada := lama[i][k]; ada {
					baru[i][k] = v
				}
			}
		}
	}
	h.SetelDaftar(DaftarSpreading, baru)
	return terpicu
}

// pasangAnggotaServer mengisi `anggotaServerSpreading` baris `baru` dari baris
// server `lama`: baris kiriman ke-i dipasangkan dengan baris server berikutnya
// (belum terpakai, berurutan) yang keempat anggotanya SAMA dengan yang dibawa
// kiriman; kiriman tanpa nilai keempatnya, atau tanpa pasangan = kosong.
func pasangAnggotaServer(lama, kiriman, baru []Baris) {
	sama := func(a, b Baris) bool {
		for _, m := range anggotaServerSpreading {
			if a[m] != b[m] {
				return false
			}
		}
		return true
	}
	j := 0
	for i, k := range kiriman {
		if sama(k, Baris{}) {
			continue // baris baru / tanpa nilai server: tidak memakai pasangan
		}
		for c := j; c < len(lama); c++ {
			if sama(lama[c], k) {
				for _, m := range anggotaServerSpreading {
					if v, ada := lama[c][m]; ada {
						baru[i][m] = v
					}
				}
				j = c + 1
				break
			}
		}
	}
}

// nilaiBerubah - isian sel berubah = event `change` sel Pega: teks yang
// dikirim layar berbeda dari teks yang server kirimkan (layar mengirim nilai
// mentah apa adanya, `frontend/components/InputAngka.tsx`). Bukan pembandingan
// dua nilai uang (spec-penyimpanan AC 25): tidak ada rumus yang bergantung
// pada selisihnya, hanya ada/tidaknya ketikan.
func nilaiBerubah(lama, baru string) bool {
	return strings.TrimSpace(lama) != strings.TrimSpace(baru)
}

// TerapkanPemicu memutar ulang di server action set sel yang dipicu kiriman
// layar (`PemicuLayar`), atas halaman yang medan turunannya sudah dihitung
// (`services.turunkan` - NetPremium, BalanceDueTo; CountNetPremi_act adalah
// langkah 7 CountOGPONP_Act):
//
//	Angsuran   `FillPaymentInstallment` (baris disusun ulang; DueDate = sekarang)
//	Uang       `CountOGPONP_Act` langkah 10 `SetValidateInstallment_Act` (3.1 Premium /
//	           PaymentTotal per baris) - bila Installment tidak ikut berubah - dan
//	           langkah 9 `CountSpreading_Act`
//	Spreading  `CountSpreading_Act` langkah 4-5 (langkah 1-3 berlabel `//`)
//
// ⚠️ `[penyesuaian sadar]` Urutan beberapa refresh dalam satu kiriman tidak
// terbaca; Installment yang berubah bersama sel uang dihitung ulang dengan
// BalanceDueTo akhir (FillPaymentInstallment sesudah sel uang).
// Total spreading selalu = jumlah baris (`HitungTotalSpreading`, langkah 4.2/5).
// `sekarang` menggantikan `@CurrentDateTime()` (FillPaymentInstallment 3.5).
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
	if p.Uang || p.Spreading {
		if err := CountSpreading(h, 0); err != nil {
			return err
		}
	}
	return HitungTotalSpreading(h)
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
