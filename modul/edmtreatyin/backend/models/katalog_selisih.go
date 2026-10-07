package models

// Untuk apa berkas ini: KATALOG TABEL PROYEKSI SELISIH - pemetaan properti `PolicyTreatyIn.TreatyDifference.*`
// dan `PolicyTreatyIn.TreatyXOLDifferenceList(c).ValueList(l)` ke empat tabel yang DIBUAT modul ini (migrasi
// 360-363, `docs/STRUKTUR-TABEL-EDM-TREATY-IN.md`; spec-penyimpanan ID-6, ID-7, ID-23..ID-27, ID-42). Dicocokkan
// dengan DDL oleh `backend/migrasi_test.go`.
//
// ⛔ Tabel PROYEKSI, bukan tabel sumber: ditulis Go di transaksi yang sama dengan generasinya; bangun ulang hanya
// baris SUMBER = 'GO' (AC 24-27). Sumber kebenaran tetap dua baris generasi (ID-4).
//
// ⛔ Tidak disimpan, dan sebabnya:
//   - TreatyDifference.TotalPremium / TotalClaim / TotalSharePercentagePremium / TotalSharePercentageClaim -
//     turunan baris spreading selisih dan total spreading generasi baru (`HitungTotalSelisih`; pola NB
//     `HitungTotalSpreading`).
//   - TreatyDifference.DueTo / Currency / GuaranteeFund / Layer* - DIBACA `InsetTreatyInProdAddendum_Act` blok 9,
//     tetapi tidak satu pun rule korpus EDM menulisnya (sensus `TreatyDifference.*`): selalu kosong.
//   - Induk TreatyXOLDifferenceList(c) - kuncinya salinan kunci lapisan, nilainya jumlah lapisan
//     (`CalculateDifferenceEDM_act` langkah 2); DueTo induknya (langkah 1.1, `local.duetovalue` iterasi
//     sebelumnya) nol pembaca. Dibangun ulang dari baris lapisan (`BangunIndukSelisihXOL`).

const (
	sd = HalamanPolis + ".TreatyDifference."
	// DaftarSelisihSpreading, DaftarSelisihAngsuran - daftar selisih di halaman kerja.
	DaftarSelisihSpreading = HalamanPolis + ".TreatyDifference.SpreadingRiskList"
	DaftarSelisihAngsuran  = HalamanPolis + ".TreatyDifference.ListInstallment"
	// DaftarSelisihXOL - induk per mata uang; anaknya `JalurAnak(DaftarSelisihXOL, c, "ValueList")`.
	DaftarSelisihXOL = HalamanPolis + ".TreatyXOLDifferenceList"
)

// TabelSelisih - T_POLIS_DIFFERENCE <- PolicyTreatyIn.TreatyDifference (`EDMTCalculateTreatyDifference`
// langkah 1). Kolom kunci (ID, POLIS_ID, NOPOLIS, PRODKE, EDM_NO, IDPEGA, SUMBER) ditulis repository di luar
// katalog.
var TabelSelisih = Tabel{Nama: "T_POLIS_DIFFERENCE", Kolom: []Kolom{
	kKode(sd+"Installment", "INSTALLMENT", 16),
	kUang(sd+"GrossPremium", "GROSS_PREMIUM"),
	kUang(sd+"PremiOgp", "PREMI_OGP"),
	kPersen(sd+"RiCommOgp", "RI_COMM_OGP"),
	kUang(sd+"ResultOgp1", "RESULT_OGP1"),
	kPersen(sd+"OveriddingCommOgp", "OVERIDDING_COMM_OGP"),
	kUang(sd+"ResultOgp2", "RESULT_OGP2"),
	kUang(sd+"Claim", "CLAIM"),
	kUang(sd+"SalvageValue", "SALVAGE_VALUE"),
	kUang(sd+"ExcessLoss", "EXCESS_LOSS"),
	kUang(sd+"NetPremium", "NET_PREMIUM"),
	kUang(sd+"BalanceDueTo", "BALANCE_DUE_TO"),
	kUang(sd+"PremiOnp", "PREMI_ONP"),
	kPersen(sd+"RiCommOnp", "RI_COMM_ONP"),
	kUang(sd+"ResultOnp1", "RESULT_ONP1"),
	kPersen(sd+"OveriddingCommOnp", "OVERIDDING_COMM_ONP"),
	kUang(sd+"ResultOnp2", "RESULT_ONP2"),
	kUang(sd+"Deduction1", "DEDUCTION1"),
	kUang(sd+"Deduction2", "DEDUCTION2"),
	kUang(sd+"PPNValue", "PPN_VALUE"),
	kUang(sd+"PPHValue", "PPH_VALUE"),
	kUang(sd+"BalanceBeforeTax", "BALANCE_BEFORE_TAX"),
	kUang(sd+"BalanceBeforePPH", "BALANCE_BEFORE_PPH"),
}}

// TabelSelisihSpreading - T_POLIS_DIFFERENCE_SPREADING <- TreatyDifference.SpreadingRiskList (langkah 2.1).
var TabelSelisihSpreading = Tabel{Nama: "T_POLIS_DIFFERENCE_SPREADING", Daftar: DaftarSelisihSpreading, Kolom: []Kolom{
	kKode("TreatyType", "TREATY_TYPE", 64),
	kTeks("TreatyName", "TREATY_NAME", 255),
	kPersen("SharePercentage", "SHARE_PERCENTAGE"),
	kPersen("ClaimPercentage", "CLAIM_PERCENTAGE"),
	kUang("PremiumSpreaded", "PREMIUM_SPREADED"),
	kUang("ClaimSpreaded", "CLAIM_SPREADED"),
}}

// TabelSelisihAngsuran - T_POLIS_DIFFERENCE_INSTALMENT <- TreatyDifference.ListInstallment (langkah 3.1).
var TabelSelisihAngsuran = Tabel{Nama: "T_POLIS_DIFFERENCE_INSTALMENT", Daftar: DaftarSelisihAngsuran, Kolom: []Kolom{
	kCacah("InstallmentNo", "INSTALLMENT_NO"),
	kTgl("DueDate", "DUE_DATE"),
	kPersen("InstallmentPercentage", "INSTALLMENT_PERCENTAGE"),
	kUang("Premium", "PREMIUM"),
	kUang("PaymentTotal", "PAYMENT_TOTAL"),
}}

// TabelSelisihLapisan - T_POLIS_XOL_LAYER_DIFFERENCE <- TreatyXOLDifferenceList(c).ValueList(l), daftar DATAR
// urut (c, l) (`CalculateDifferenceEDM_act` langkah 1.2).
var TabelSelisihLapisan = Tabel{Nama: "T_POLIS_XOL_LAYER_DIFFERENCE", Daftar: "ValueList", Kolom: []Kolom{
	kKode("Layer", "LAYER", 64),
	kKode("LayerType", "LAYER_TYPE", 64),
	kKode("LayerPart", "LAYER_PART", 64),
	kKode("LayerPartType", "LAYER_PART_TYPE", 64),
	kKode("Currency", "CURRENCY", 16),
	kKode("IDCurrency", "ID_CURRENCY", 64),
	kKode("DueTo", "DUE_TO", 16),
	kUang("GrossPremi", "GROSS_PREMI"),
	kUang("Deduction", "DEDUCTION"),
	kUang("NetPremi", "NET_PREMI"),
	kUang("DueToValue", "DUE_TO_VALUE"),
	kUang("BrokerageFeeSebenarnya", "BROKERAGE_FEE_SEBENARNYA"),
	kUang("PPHValue", "PPH_VALUE"),
	kUang("PPNValue", "PPN_VALUE"),
	kUang("NetPremiAfterPPH", "NET_PREMI_AFTER_PPH"),
	kUang("NetPremiAfterPPN", "NET_PREMI_AFTER_PPN"),
	kUang("NetPremiAfterTax", "NET_PREMI_AFTER_TAX"),
}}

// SemuaTabelSelisih - urutan tulis (induk lebih dulu).
var SemuaTabelSelisih = []Tabel{TabelSelisih, TabelSelisihSpreading, TabelSelisihAngsuran, TabelSelisihLapisan}

// Sumber baris proyeksi (kolom SUMBER, ID-25).
const (
	SumberPega = "PEGA"
	SumberGo   = "GO"
)

// DatarSelisihLapisan - baris lapisan TreatyXOLDifferenceList(c).ValueList(l) dalam urutan (c, l), untuk
// disimpan sebagai daftar datar (NOURUT 1..n).
func DatarSelisihLapisan(h *Halaman) []Baris {
	var out []Baris
	for c := range h.AmbilDaftar(DaftarSelisihXOL) {
		out = append(out, h.AmbilDaftar(JalurAnak(DaftarSelisihXOL, c+1, "ValueList"))...)
	}
	return out
}

// BangunIndukSelisihXOL menyusun ulang TreatyXOLDifferenceList dari baris lapisan datar: satu induk per
// ID_CURRENCY berurut kemunculan pertama (pengelompokan per mata uang `InsertToTreatyXOLList` = satu baris induk
// per mata uang angsuran master). Kunci induk = kunci lapisan pertamanya (`CalculateDifferenceEDM_act` langkah
// 1.1 menyalin kunci baris induk TreatyXOLList yang sama); nilai uang = jumlah lapisan (langkah 2.2-2.3).
func BangunIndukSelisihXOL(h *Halaman, lapisan []Baris) error {
	var induk []Baris
	var anak [][]Baris
	for _, b := range lapisan {
		n := len(induk) - 1
		if n < 0 || induk[n]["IDCurrency"] != b["IDCurrency"] {
			induk = append(induk, Baris{
				"Layer": b["Layer"], "LayerType": b["LayerType"], "LayerPart": b["LayerPart"],
				"LayerPartType": b["LayerPartType"], "Currency": b["Currency"], "IDCurrency": b["IDCurrency"],
			})
			anak = append(anak, nil)
			n++
		}
		anak[n] = append(anak[n], b)
	}
	h.SetelDaftar(DaftarSelisihXOL, induk)
	for c := range induk {
		h.SetelDaftar(JalurAnak(DaftarSelisihXOL, c+1, "ValueList"), anak[c])
	}
	return JumlahIndukSelisihXOL(h)
}

// medanUangLapisanXOL - sepuluh medan uang lapisan selisih XOL (`CalculateDifferenceEDM_act` langkah 2.2.1).
var medanUangLapisanXOL = []string{"GrossPremi", "Deduction", "NetPremi", "DueToValue", "BrokerageFeeSebenarnya",
	"PPNValue", "PPHValue", "NetPremiAfterPPN", "NetPremiAfterPPH", "NetPremiAfterTax"}

// JumlahIndukSelisihXOL = `CalculateDifferenceEDM_act` langkah 2: setiap induk TreatyXOLDifferenceList(c)
// menerima jumlah sepuluh medan uang lapisannya.
func JumlahIndukSelisihXOL(h *Halaman) error {
	k := &kalkulator{}
	induk := h.AmbilDaftar(DaftarSelisihXOL)
	for c := range induk {
		total := map[string]string{}
		for _, m := range medanUangLapisanXOL {
			jml := k.d("0")
			for _, l := range h.AmbilDaftar(JalurAnak(DaftarSelisihXOL, c+1, "ValueList")) {
				jml = k.Tambah(jml, angkaBaris(k, l, m))
			}
			total[m] = formatAngka(jml)
		}
		for m, v := range total {
			induk[c][m] = v
		}
	}
	if k.err != nil {
		return k.err
	}
	h.SetelDaftar(DaftarSelisihXOL, induk)
	return nil
}
