package models

// Untuk apa berkas ini: PILIH BISNIS dan PRA-PROSES BERDATA - langkah-langkah
// rule yang membaca tabel acuan, dipisah dari pembacaannya. Repository
// membaca barisnya; fungsi di sini menerapkan barisnya ke halaman persis
// seperti langkah aslinya (seam 3).
//
// Sumber:
//   - `Activity/SetValue_Act` -> `InputPolicyTreatyInDetail_preACT` (tombol
//     "Choose" popup `BusinessAndSOBList`, tiket 01)
//   - `Activity/InputPolicyTreatyInPre_Act` (pra-proses kedua flow action)
//   - `Activity/CheckDataMkt`, `SetCurrency_act`, `SetTreatyCurrencyID`,
//     `ProtectDate`, `RemoveTypeTax_ACT`
//
// ⛔ Bagian JSON `InputPolicyTreatyInDetail_preACT` TIDAK dibangun
// (`[keputusan work owner]` P29, PERTANYAAN-untuk-DBA "enam aktivitas yang
// membongkar JSON tidak dimigrasi"): langkah 9-10 (`adoptJSONObject`), 13
// (jadwal angsuran dari `TreatyIn.INSTALLMENT`), 16
// (`InputPolicyTreatyInDetail_NonProp`), 17 (`TreatyInputPctCommSpreading` -
// `FetchMasterTreatyIn`), 18 (`TreatyIn.LimitShareSummaryList`). Yang
// dibangun: langkah 3-8, 11, 14, 15 - seluruhnya membaca view dan tabel acuan.

import (
	"strings"
	"time"

	"nusantarare/inti/backend/utils"
)

// BarisKontrak adalah satu baris view `TREATYINDETAILJOINEDM` /
// `TREATYINDETAIL`: nilai teks per NAMA KOLOM view (ID, TREATYID, ...).
type BarisKontrak = Baris

// TerapkanDetailKontrak = `InputPolicyTreatyInDetail_preACT` langkah 3, 5, 6, 11
// atas baris pertama hasil RD `BrowseTreatyJoinEDM` (filter `.ID = Param.ID`).
//
// ⚠️ Langkah 3 membaca `pxResults(1).CURRENCYID` - kolom yang TIDAK ada di RD
// maupun view (39 kolom, PERTANYAAN-untuk-DBA). Nilainya selalu kosong, dan
// itulah sebabnya langkah 4 (`SetTreatyCurrencyID` menurut nama) selalu
// berjalan: services membaca ID mata uang lalu memanggil `SetelIDMataUang`.
func TerapkanDetailKontrak(h *Halaman, b BarisKontrak) {
	p := func(m, v string) { h.Setel(HalamanPolis+"."+m, v) }
	q := func(m, v string) { h.Setel(HalamanQuotation+"."+m, v) }
	// langkah 3
	p("IDCurrency", b["CURRENCYID"])
	p("Currency", b["LIMITCURRENCY"])
	p("TreatyYear", b["TREATYYEAR"])
	p("TreatyGroupID", b["TREATYGROUPID"])
	p("TreatyGroupName", b["TREATYGROUP"])
	p("BizName", b["CLASSOFBUSINESS"])
	p("BizCode", b["CLASSOFBUSINESSID"])
	q("BusinessName", b["CLASSOFBUSINESS"])
	p("SOBName", b["SOB"])
	p("SOB", b["SOBID"])
	p("CedingCoName", b["CEDING"])
	p("CedingCo", b["CEDINGID"])
	p("InsuredID", h.Ambil(HalamanQuotation+".InsuredID"))
	p("InsuredName", h.Ambil(HalamanQuotation+".InsuredName"))
	p("NoOffer", b["TREATYID"])
	p("PremiOgp", b["MDPVALUE"])
	p("Installment", b["INSTALLMENTNO"])
	p("Deduction1", b["DEDUCTION1"])
	p("Deduction2", b["DEDUCTION2"])
	p("TreatyType", b["TREATYTYPE"])
	p("LayerType", b["LAYERTYPE"])
	p("Layer", b["LAYER"])
	p("LayerPartType", b["LAYERPARTTYPE"])
	p("LayerPart", b["LAYERPART"])
	p("ShareCurrency", b["SHARECURRENCY"])
	p("ShareValue", b["SHAREVALUE"])
	q("ProportionalType", b["PROPORTIONTYPE"])
	q("SourceOfBusiness", b["SOBID"])
	q("SobName", b["SOB"])
	q("CedingCoName", b["CEDING"])
	q("CedingCo", b["CEDINGID"])
	// langkah 5: TREATYTYPE kosong -> "XOL"
	if b["TREATYTYPE"] == "" {
		p("TreatyType", "XOL")
	}
	// langkah 6: TreatyIn.ID = pxResults(1).ID
	h.Setel(HalamanMaster+".ID", b["ID"])
	// langkah 11: Property-Remove ListInstallment (beserta InstallmentList bersarang)
	hapusDaftarBeserta(h, DaftarAngsuran)
}

// TerapkanMasterKontrak mengisi halaman `TreatyIn` dari baris view yang sama
// - pengganti halaman master JSON untuk medan yang VIEW punya (P29):
// `Commencement`, `Termination` (layar atasan dan admin menampilkannya).
// Medan master lain (RNMShareP, RNMShare, BrokeragePercentP, CurrencyList,
// INSTALLMENT, Limits...) tidak punya kolom padanan - lihat `MasterTersedia`.
func TerapkanMasterKontrak(h *Halaman, b BarisKontrak) {
	h.Setel(HalamanMaster+".Commencement", b["COMMENCEMENT"])
	h.Setel(HalamanMaster+".Termination", b["TERMINATION"])
}

// SetelIDMataUang = `SetTreatyCurrencyID` langkah 3 (RDB `GetCurrencyIDByName`).
func SetelIDMataUang(h *Halaman, id string) { h.Setel(HalamanPolis+".IDCurrency", id) }

// SetelNamaMataUang = `SetCurrency_act` langkah 3 (RDB `GetCurrency`,
// `CURRENCY as HASIL1 where id = Param.CURR`).
func SetelNamaMataUang(h *Halaman, nama string) { h.Setel(HalamanPolis+".Currency", nama) }

// SetelOJK = `FetchTreatyGroupOJK` langkah 4.
func SetelOJK(h *Halaman, ojkBusinessID string) {
	h.Setel(HalamanPolis+".OJKBusinessID", ojkBusinessID)
}

// SetelGrupLama = `FetchTreatyGroupOldID` langkah 6 (dipanggil
// `GeneratePolicyNoTreaty_Act`; langkah 7 preACT berlabel `//`).
func SetelGrupLama(h *Halaman, oldID string) {
	h.Setel(HalamanPolis+".TreatyGroupOldID", oldID)
}

// TerapkanKlien = `InputPolicyTreatyInDetail_preACT` langkah 14.3: InsuredID
// diisi hasil `GetClientID_SQL` HANYA bila Quotation.InsuredID kosong.
func TerapkanKlien(h *Halaman, clientID string) {
	if h.Ambil(HalamanQuotation+".InsuredID") != "" {
		return
	}
	h.Setel(HalamanQuotation+".InsuredID", clientID)
	h.Setel(HalamanPolis+".InsuredID", clientID)
}

// KunciCariBisnis = isian `OldID.CARI2` sebelum RDB `GetOldIDBusiness_SQL`
// (`where ID = CARI2 OR NOTE = CARI2`).
//
//	preACT 14.4-14.5.7 / Pre_Act 2.1-2.7:
//	  CARI2 = BusinessName
//	  @contains(BusinessName,"MBU")                      -> "MOTOR VEHICLE"
//	  ADVANCE PAYMENT / BID OR TENDER / PAYMENT / PERFORMANCE BONDS
//	                                                     -> @replaceAll(nama,"S","")
//	  "CUSTOMS BOND"             -> "OTHERS CUSTOMS BOND"
//	  "ASURANSI KREDIT"          -> "ASURANSI KREDIT (CASH LOAN)"
//	  "BOILER & PRESSURE VESSEL" -> "BOILER & EXCAVATOR"
//	  "GOLF INSURANCE"           -> "HOLE IN ONE"
//	Pre_Act 2.8 SAJA: "BID OR TENDER BONDS" -> "BID BOND" (menimpa 2.3)
//
// ⚠️ `@replaceAll(...,"S","")` membuang SETIAP huruf S, bukan hanya akhiran -
// ditiru apa adanya.
func KunciCariBisnis(nama string, langkahPra bool) string {
	k := nama
	if strings.Contains(nama, "MBU") {
		k = "MOTOR VEHICLE"
	}
	switch nama {
	case "ADVANCE PAYMENT BONDS", "BID OR TENDER BONDS", "PAYMENT BONDS", "PERFORMANCE BONDS":
		k = strings.ReplaceAll(nama, "S", "")
	case "CUSTOMS BOND":
		k = "OTHERS CUSTOMS BOND"
	case "ASURANSI KREDIT":
		k = "ASURANSI KREDIT (CASH LOAN)"
	case "BOILER & PRESSURE VESSEL":
		k = "BOILER & EXCAVATOR"
	case "GOLF INSURANCE":
		k = "HOLE IN ONE"
	}
	if langkahPra && nama == "BID OR TENDER BONDS" {
		k = "BID BOND"
	}
	return k
}

// BarisBisnis adalah hasil pertama `GetOldIDBusiness_SQL` (CARI1 oldid,
// CARI2 GROUPPANEL, CARI3 ID); kosong semua bila tidak ada baris.
type BarisBisnis struct{ OldID, GroupPanel, ID string }

// TerapkanBisnisPilih = preACT langkah 14.7-14.9.
func TerapkanBisnisPilih(h *Halaman, b BarisBisnis) {
	h.Setel(HalamanQuotation+".BusinessOldId", b.OldID)
	h.Setel(HalamanQuotation+".GroupPanel", b.GroupPanel)
	h.Setel(HalamanQuotation+".BusinessCode", b.ID)
	h.Setel(HalamanPolis+".BizCode", b.ID)
	// 14.8 DT BusinessType_DeT -> Quotation.BusinessType
	h.Setel(HalamanQuotation+".BusinessType", GolongkanJenisUsaha(b.GroupPanel, b.OldID))
	// 14.9 QuotationData = Quotation
	SalinQuotation(h)
}

// TerapkanBisnisPra = `InputPolicyTreatyInPre_Act` langkah 2.10 (berjalan hanya
// bila `PolicyTreatyIn.BizCode == ""`).
func TerapkanBisnisPra(h *Halaman, b BarisBisnis) {
	h.Setel(HalamanQuotation+".BusinessOldId", b.OldID)
	h.Setel(HalamanQuotation+".GroupPanel", b.GroupPanel)
	h.Setel(HalamanQuotation+".BusinessCode", b.ID)
	h.Setel(HalamanPolis+".BizCode", b.ID)
	h.Setel(HalamanPolis+".QuotationData.BusinessCode", b.ID)
	h.Setel(HalamanPolis+".QuotationData.BusinessOldId", b.OldID)
	h.Setel(HalamanPolis+".QuotationData.GroupPanel", b.GroupPanel)
}

// PraprosesTanggal = `InputPolicyTreatyInPre_Act` langkah 3-4 dan 9.
//
//	3-4  StatementDate = ProductionDate = sysdate. ⚠️ Kotak When langkah 3
//	     ("PositionNote==ReasTreatyInAdmin && SuggestList kosong") TIDAK
//	     dicentang - langkahnya SELALU berjalan, di ketiga posisi.
//	9    hari StatementDate melewati hari tutup buku -> ProductionDate = tanggal 1
//	     bulan berikut (`GeserTanggalProduksi` - hari dari TANGGAL_CLOSING)
func PraprosesTanggal(h *Halaman, sekarang time.Time, hariClosing int) {
	h.Setel(HalamanPolis+".StatementDate", FormatTanggalWaktu(sekarang))
	h.Setel(HalamanPolis+".ProductionDate", FormatTanggalWaktu(GeserTanggalProduksi(sekarang, hariClosing)))
}

// BarisMO adalah hasil `Obj-Browse` marketing officer (CheckDataMkt langkah 3).
type BarisMO struct{ ID, ClientID, ClientName, TeamGroup, BranchDetailID, BranchDetailName string }

// TerapkanMO = `CheckDataMkt` langkah 2 dan 4 (langkah 3 dilewati bila
// QuotationData.MOID kosong - hasil kosong). Langkah 4 juga menulis
// `pyWorkPage.OfferFacIn.QuotationData.*` - halaman kasus FAKULTATIF, tidak
// ada di kasus treaty; tidak dibangun. Langkah 6-7 berlabel `//`.
func TerapkanMO(h *Halaman, m BarisMO) {
	h.Setel(HalamanPolis+".MarketingOfficer", "")
	q := func(n, v string) { h.Setel(HalamanQuotation+"."+n, v) }
	d := func(n, v string) { h.Setel(HalamanPolis+".QuotationData."+n, v) }
	q("MOID", m.ID)
	q("MarketingCode", m.ClientID)
	q("MarketingName", m.ClientName)
	q("TeamGroup", m.TeamGroup)
	q("BranchCode", m.BranchDetailID)
	q("BranchName", m.BranchDetailName)
	d("MOID", m.ID)
	d("MarketingCode", m.ClientID)
	d("MarketingName", m.ClientName)
	d("TeamGroup", m.TeamGroup)
	h.Setel(HalamanPolis+".MarketingOfficer", m.ClientName)
}

// PesanTanggalAkhir - VERBATIM `ProtectDate` langkah 1.
const PesanTanggalAkhir = "End Date Cannot be less than Start Date"

// ProtectDate = `Activity/ProtectDate`: `@CompareDates(.StartDate,.EndDate)`
// benar (StartDate SESUDAH EndDate) -> pesan pada `.EndDate`.
func ProtectDate(h *Halaman) {
	a, errA := utils.ParseTanggal(strings.TrimSpace(h.Ambil(HalamanPolis + ".StartDate")))
	b, errB := utils.ParseTanggal(strings.TrimSpace(h.Ambil(HalamanPolis + ".EndDate")))
	if errA == nil && errB == nil && a.After(b) {
		h.TambahPesan(HalamanPolis+".EndDate", PesanTanggalAkhir)
	}
}

// RemoveTypeTax = `RemoveTypeTax_ACT` langkah 2: FlagPPH==false -> hapus
// TypeTax. Langkah 1 (`TempWorkPage`) halaman sementara Pega, tidak ada di
// sistem baru. `FlagPPH==false` benar untuk "false" dan kosong (properti
// TrueFalse Pega).
func RemoveTypeTax(h *Halaman) {
	if h.Ambil(HalamanPolis+".FlagPPH") != "true" {
		h.Hapus(HalamanPolis + ".TypeTax")
	}
}

// TreatyEnableDisableInput = `DataTransform/TreatyEnableDisableInput` (tombol
// "Enable / Disable Input Type", tampil bila `.TreatyType='XOL'`).
//
//	1  .QuotationData.ProportionalType = "NonProportional"
//	2  WHEN IsNewPolicyNonProp == ""  -> "1"
//	3  WHEN IsNewPolicyNonProp == "0" -> "1"
//	4  WHEN IsNewPolicyNonProp == "1" -> "0"
//
// ⚠️ Langkah WHEN data transform dievaluasi BERURUTAN atas nilai terkini:
// "" dan "0" menjadi "1" lalu langkah 4 menjadikannya "0" - tombol ini SELALU
// berakhir di "0" untuk ketiga nilai itu. Ditiru apa adanya, bukan diperbaiki
// menjadi saklar.
func TreatyEnableDisableInput(h *Halaman) {
	j := HalamanPolis + ".IsNewPolicyNonProp"
	h.Setel(HalamanPolis+".QuotationData.ProportionalType", JenisNonProporsional)
	if h.Ambil(j) == "" {
		h.Setel(j, "1")
	}
	if h.Ambil(j) == "0" {
		h.Setel(j, "1")
	}
	if h.Ambil(j) == "1" {
		h.Setel(j, "0")
	}
}
