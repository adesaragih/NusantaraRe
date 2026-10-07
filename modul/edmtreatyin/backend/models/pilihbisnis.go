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
// membongkar JSON tidak dimigrasi"): langkah 9-10 (`adoptJSONObject`
// `M_TREATY_IN_DETAIL_EDM`), 13 (jadwal angsuran dari `TreatyIn.INSTALLMENT`).
// Yang dibangun: langkah 3-8, 11, 14, 15 - seluruhnya membaca view dan tabel
// acuan - dan langkah 17 (`TreatyInputPctCommSpreading`) SEBAGIAN: RiCommOgp
// dari kolom view RIOGR/RIONR (`models/komisi.go`, RALAT putaran 2); baris
// spreading-nya tidak (bukan kolom view, alasan c).
// ⭐ RALAT K8 (03-10-2026): langkah 16 (`InputPolicyTreatyInDetail_NonProp`) dan
// 18 (`TreatyIn.LimitShareSummaryList`) DIBANGUN - master XOL dibaca baca-saja
// (`nonprop.go`, `nonprop_detail.go`, `services/nonprop.go`).

import (
	"strings"
	"time"

	"nusantarare/inti/backend/utils"
)

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
	h.Setel(HalamanPolis+".StatementDate", utils.FormatTanggalWaktu(sekarang))
	h.Setel(HalamanPolis+".ProductionDate", utils.FormatTanggalWaktu(GeserTanggalProduksi(sekarang, hariClosing)))
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

// RemoveTypeTax = `RemoveTypeTax_ACT` langkah 2: FlagPPH==false -> hapus
// TypeTax. Langkah 1 (`TempWorkPage`) halaman sementara Pega, tidak ada di
// sistem baru. `FlagPPH==false` benar untuk "false" dan kosong (properti
// TrueFalse Pega).
func RemoveTypeTax(h *Halaman) {
	if h.Ambil(HalamanPolis+".FlagPPH") != "true" {
		h.Hapus(HalamanPolis + ".TypeTax")
	}
}
