package models

// Untuk apa berkas ini: SIMPAN POLIS KE TABEL PRODUKSI LAMA sesudah realisasi selesai - Flow
// `InputRealizationTreatyIn` `Utility1` ("Save json policy") -> `Activity/SaveJsonPolisTreatyIn_Act`.
//
// Patokannya ekspor yang diserahkan work owner 06-10-2026 (`D:\XML\RNM_BRD\NB Treaty In\`, akar folder):
// `SaveJsonPolisTreatyIn_Act` (pxUpdateDateTime 2026-07-09), `InsetTreatyInProd_Act` (2025-10-16), RDBList
// `InsertTreatyInProd_SQL` (2026-02-05), `SetAchivementValue` (2026-03-17), RDBList `SaveAchievementSQL`
// (2023-08-15). Isi prosedur `POOLDATA.PEGA_JSON_POLIS_TREATYIN` dan `POOLDATA.InsertUpdateAchievment` dibaca
// dari kamus data DEV (ALL_SOURCE) - keduanya satu INSERT.
//
// `[keputusan work owner 06-10-2026]` RALAT sebagian AC 16: "JSON-nya tidak disimpan, tapi tetap insert kolom
// lainnya" - baris JSON_POLIS ditulis TANPA DATA_JSON; TREATYINPRODUCTION dan ACHIEVEMENT ditulis seperti Pega.
// Nol prosedur (spec-penyimpanan AC 48): ketiganya INSERT langsung di transaksi submit (AC 29, 83).
//
// SaveJsonPolisTreatyIn_Act (2026-07-09) - kotak When langkah 4, 4.1, 6, 8 TIDAK dicentang, jadi selalu jalan:
//
//	1      Property-Remove SuggestDate/OperatorName/Suggest  ⛔ tidak ditiru: hanya membentuk dokumen JSON
//	2-4    tanggal closing, ProductionDate                   PrasimpanPolis (= GeneratePolicyNoTreaty_Act 5)
//	5      Policy.CaseID / Policy.OperatorID                 ⛔ halaman Policy hanya isi dokumen JSON
//	6      InputData.CARI1/2/21                              JSONPolis
//	7      IDNewBisnis = pyID                                PrasimpanPolis
//	8      ClaimType / ClaimPaymentType                      PrasimpanPolis
//	9      // RealizationPopulateCedingSOB                    berlabel //, tak pernah jalan
//	10     CARI3 = @ASM.GetPageJSONString()                  ⛔ perintah WO: JSON tidak disimpan
//	11     SavePolisTreatyIn_SQL                             JSONPolis
//	12     ERRMSG = IDPEGAOUT                                tanpa padanan (keluaran prosedur)
//	13     Call SetAchivementValue                           BarisCapaian
//	14-16  // INSERT KE JSON ERROR                           berlabel //
//	17     Call InsetTreatyInProd_Act                        BarisProduksi
//
// `[penyimpangan sadar]` CARI21 (JSON_POLIS.TGL_PROD) dan CARI55 (TREATYINPRODUCTION.PROD_DATE) diformat XML
// `dd/MM/yyyy hh:mm:ss` (jam 12) lalu dibaca `TO_DATE(.., 'DD/MM/YYYY HH24:MI:SS')` - ProductionDate sore
// tersimpan sebagai pagi. Di sini jam 24 apa adanya, sama dengan TGL_INP (models/usulan.go butir 3, disetujui WO
// 04-10-2026).

import (
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// Tabel produksi lama yang ditulis Utility1 (MODUL.md "Tabel warisan").
const (
	TabelJSONPolis = "JSON_POLIS"
	TabelProduksi  = "TREATYINPRODUCTION"
	TabelCapaian   = "ACHIEVEMENT"
)

// ProdKeNB - `SavePolisTreatyIn_SQL` argumen ke-4: generasi NB = '0'.
const ProdKeNB = "0"

// Teks DUE_TO - `InsetTreatyInProd_Act` 9.4 / 9.5, VERBATIM.
const teksDueToYou = "DUE TO YOU"

// ProporsionalPenuh - nilai `QuotationData.ProportionalType` jalur proporsional (`SetAchivementValue` 4-6).
const ProporsionalPenuh = "Proportional"

// KolomJSONPolis - kolom `POOLDATA.json_polis` yang ditulis `PEGA_JSON_POLIS_TREATYIN` selain DATA_JSON
// (perintah WO: tidak disimpan) dan TGL_INPUT (SYSDATE, ditulis SQL). NOENDORS = NULL (argumen ke-3).
// Panjang = kolom DEV (ALL_TAB_COLUMNS 06-10-2026).
var KolomJSONPolis = []Kolom{
	kKode("IDPEGA", "IDPEGA", 50),
	kKode("NOPOLIS", "NOPOLIS", 100),
	kKode("PRODKE", "PRODKE", 5),
	kTglWaktu("TGL_PROD", "TGL_PROD"),
	kTeks("USERNAME", "USERNAME", 200),
}

// KolomProduksi - `InsertTreatyInProd_SQL`, urutan kolom INSERT apa adanya (58 kolom). Baris produksi memakai NAMA
// KOLOM sebagai kunci. Angka (`TO_NUMBER(REPLACE(x, ',', '.'))`) bergolongan uang/persen, tanggal
// `TO_DATE(x, 'YYYYMMDD')` bergolongan tanggal, PROD_DATE tanggal-waktu. Panjang = kolom DEV.
var KolomProduksi = []Kolom{
	kKode("IDPEGA", "IDPEGA", 100),
	kKode("NOPOLIS", "NOPOLIS", 100),
	kKode("BUSINESSCODE", "BUSINESSCODE", 100),
	kTeks("CEDINGCO", "CEDINGCO", 1000),
	kTeks("SOB", "SOB", 1000),
	kTeks("INSUREDNAME", "INSUREDNAME", 1000),
	kKode("TREATYGROUPID", "TREATYGROUPID", 100),
	kTeks("TREATYGROUP", "TREATYGROUP", 100),
	kTeks("DUE_TO", "DUE_TO", 100),
	kTgl("STATEMENT_DATE", "STATEMENT_DATE"),
	kTgl("BEGINDATE", "BEGINDATE"),
	kTgl("ENDDATE", "ENDDATE"),
	kKode("CURR_ID", "CURR_ID", 100),
	kKode("UW_YEAR", "UW_YEAR", 100),
	kUang("PREMI_OGP", "PREMI_OGP"),
	kPersen("PERCENT_RI_COMM_OGP", "PERCENT_RI_COMM_OGP"),
	kUang("RI_COMM_OGP", "RI_COMM_OGP"),
	kPersen("PERCENT_OVERRIDING_COMM_OGP", "PERCENT_OVERRIDING_COMM_OGP"),
	kUang("OVERRIDING_COMM_OGP", "OVERRIDING_COMM_OGP"),
	kUang("CLAIM", "CLAIM"),
	kUang("EXCESS_LOSS", "EXCESS_LOSS"),
	kUang("PREMI_ONP", "PREMI_ONP"),
	kPersen("PERCENT_RI_COMM_ONP", "PERCENT_RI_COMM_ONP"),
	kUang("RI_COMM_ONP", "RI_COMM_ONP"),
	kPersen("PERCENT_OVERRIDING_COMM_ONP", "PERCENT_OVERRIDING_COMM_ONP"),
	kUang("OVERRIDING_COMM_ONP", "OVERRIDING_COMM_ONP"),
	kUang("DEDUCTION1", "DEDUCTION1"),
	kUang("DEDUCTION2", "DEDUCTION2"),
	kUang("NET_PREMIUM", "NET_PREMIUM"),
	kUang("BALANCE_DUE_TO", "BALANCE_DUE_TO"),
	kKode("MARKETINGOFFICERCODE", "MARKETINGOFFICERCODE", 100),
	kKode("NOOFFER", "NOOFFER", 100),
	kKode("JN_REAS", "JN_REAS", 100),
	kPersen("PCT_SHARE_PREMI", "PCT_SHARE_PREMI"),
	kPersen("PCT_SHARE_CLAIM", "PCT_SHARE_CLAIM"),
	kTeks("TREATYTYPE", "TREATYTYPE", 100),
	kKode("LAYERTYPE", "LAYERTYPE", 20),
	kKode("LAYER", "LAYER", 20),
	kKode("LAYERPARTTYPE", "LAYERPARTTYPE", 20),
	kKode("LAYERPART", "LAYERPART", 20),
	kKode("CLAIMTYPE", "CLAIMTYPE", 20),
	kKode("CLAIMPAYMENTTYPE", "CLAIMPAYMENTTYPE", 20),
	kKode("QUARTER", "QUARTER", 5),
	kUang("PPHVALUE", "PPHVALUE"),
	kUang("PPNVALUE", "PPNVALUE"),
	kUang("BALANCE_BEFORE_PPH", "BALANCE_BEFORE_PPH"),
	kUang("BALANCE_BEFORE_TAX", "BALANCE_BEFORE_TAX"),
	kKode("CEDINGCOID", "CEDINGCOID", 50),
	kKode("SOBID", "SOBID", 50),
	kKode("INSUREDID", "INSUREDID", 50),
	kKode("JENIS_TREATY", "JENIS_TREATY", 50),
	kKode("QUARTER_YEAR", "QUARTER_YEAR", 10),
	kUang("GUARANTEE_FUND", "GUARANTEE_FUND"),
	kUang("SHARE_VALUE", "SHARE_VALUE"),
	kKode("PROPORTIONALTYPE", "PROPORTIONALTYPE", 50),
	kUang("OUTSTANDING_CLAIM", "OUTSTANDING_CLAIM"),
	kTglWaktu("PROD_DATE", "PROD_DATE"),
	kKode("STATEMENT_TYPE", "STATEMENT_TYPE", 10),
}

// KolomCapaian - `POOLDATA.InsertUpdateAchievment` (INSERT ACHIEVEMENT), urutan parameter prosedur; TGL_PROD =
// SYSDATE (argumen terakhir `SaveAchievementSQL`), ditulis SQL. Panjang = kolom DEV.
var KolomCapaian = []Kolom{
	kKode("IDPEGA", "IDPEGA", 50),
	kKode("NOPOLIS", "NOPOLIS", 100),
	kTeks("NOOFFER", "NOOFFER", 4000),
	kTeks("SOBNAME", "SOBNAME", 4000),
	kTeks("TREATYGROUPNAME", "TREATYGROUPNAME", 4000),
	kTeks("TREATYTYPE", "TREATYTYPE", 4000),
	kTeks("QUARTER", "QUARTER", 4000),
	kTeks("QUARTERYEAR", "QUARTERYEAR", 4000),
	kTeks("IDCURRENCY", "IDCURRENCY", 4000),
	kTeks("CURRENCY", "CURRENCY", 4000),
	kUang("PREMIUM", "PREMIUM"),
	kUang("RICOMM", "RICOMM"),
	kUang("BROKERAGE", "BROKERAGE"),
	kUang("NETPREMIUM", "NETPREMIUM"),
	kUang("PAIDCLAIM", "PAIDCLAIM"),
	kUang("OUTSTANDINGCLAIM", "OUTSTANDINGCLAIM"),
	kTeks("PROPORTIONALTYPE", "PROPORTIONALTYPE", 4000),
}

// SimpananPolis - seluruh baris yang ditulis Utility1 untuk satu kasus.
type SimpananPolis struct {
	// IDPega - `pyWorkPage.pzInsKey`.
	IDPega string
	// JSONPolis - satu baris json_polis (kunci = KolomJSONPolis). Tidak ditulis bila IDPEGA sudah punya baris:
	// prosedur Pega gagal pada kunci utama, menelan galatnya (ROLLBACK + pesan), dan aktivitas berlanjut.
	JSONPolis Baris
	// Produksi - baris TREATYINPRODUCTION urutan sisip; ditulis HANYA bila IDPEGA belum punya baris
	// (`InsetTreatyInProd_Act` 6-8 "when there is data already, exit activity").
	Produksi []Baris
	// Capaian - baris ACHIEVEMENT. NOPOLIS diisi penulis dari json_polis kasus ini (`GetPolicyNoByCaseId`).
	Capaian []Baris
}

// CacahProduksi - jumlah baris tabel produksi lama yang sudah dimiliki satu IDPEGA.
type CacahProduksi struct {
	JSONPolis, Capaian, Produksi int
}

// PrasimpanPolis = medan halaman yang Utility1 tulis sebelum halaman disimpan pada End:
// `SaveJsonPolisTreatyIn_Act` 4 (ProductionDate), 7 (IDNewBisnis), 8 (ClaimType / ClaimPaymentType), dan
// `InsetTreatyInProd_Act` 9.1 (BizCode, jalur bukan NonProp baru).
func PrasimpanPolis(h *Halaman, id string, sekarang time.Time, hariClosing int) {
	// 4 blok REPEAT (kotak When `OfferFacIn.PolicyData.PolicyNo==""` tak dicentang): 4.1 (When tak dicentang)
	// ProductionDate = sekarang; 4.2-4.3 = GeneratePolicyNoTreaty_Act 5.2-5.3. StatementDate dibaca seperti
	// `repository.TerbitkanNomorPolis`.
	var statement time.Time
	if s := strings.TrimSpace(h.Ambil(pt + "StatementDate")); s != "" {
		if t, err := time.ParseInLocation(utils.TanggalWaktu, s, sekarang.Location()); err == nil {
			statement = t
		}
	}
	h.Setel(pt+"ProductionDate", utils.FormatTanggalWaktu(TanggalProduksiNomor(sekarang, statement, hariClosing)))
	PrasimpanMedan(h, id)
}

// PrasimpanMedan = `PrasimpanPolis` TANPA langkah 4: ProductionDate yang tersimpan dipakai apa adanya. Dipakai alat
// `simpanproduksi` untuk kasus yang selesai sebelum penulisan produksi ada (tanggal produksinya sudah tetap).
func PrasimpanMedan(h *Halaman, id string) {
	// 7
	h.Setel(pt+"IDNewBisnis", id)
	// 8 (kotak When `.ClaimType = "" || .ClaimPaymentType = ""` tak dicentang - SELALU dihitung ulang):
	// @if((.Claim!=0&&.Claim!="")||(.SalvageValue!=0&&.SalvageValue!=""), "SOA" / "Claim", "")
	ada := bukanNolIsi(h.Ambil(pt+"Claim")) || bukanNolIsi(h.Ambil(pt+"SalvageValue"))
	h.Setel(pt+"ClaimType", jika(ada, "SOA", ""))
	h.Setel(pt+"ClaimPaymentType", jika(ada, "Claim", ""))
	// InsetTreatyInProd_Act 9.1 - hanya di dalam blok 9 (`IsNewPolicyNonProp != 1`)
	if !PolisNonPropBaru(h) && h.Ambil(pt+"BizCode") == "" {
		h.Setel(pt+"BizCode", NilaiQuotation(h, "BusinessCode"))
	}
}

// SusunSimpananPolis menyusun baris Utility1 dari halaman sesudah `PrasimpanPolis`. `akun` =
// `OperatorID.pyUserIdentifier` (USERNAME json_polis).
func SusunSimpananPolis(h *Halaman, id, akun string) (SimpananPolis, error) {
	prod, err := BarisProduksi(h, id)
	if err != nil {
		return SimpananPolis{}, err
	}
	return SimpananPolis{
		IDPega:    KunciInstans(id),
		JSONPolis: JSONPolis(h, id, akun),
		Produksi:  prod,
		Capaian:   BarisCapaian(h, id),
	}, nil
}

// JSONPolis = `SavePolisTreatyIn_SQL` -> `PEGA_JSON_POLIS_TREATYIN(pzInsKey, PolicyTreatyIn.PolicyNo, NULL, '0',
// CARI21, OperatorID.pyUserIdentifier, CARI3)`; CARI21 = ProductionDate (langkah 6).
func JSONPolis(h *Halaman, id, akun string) Baris {
	return Baris{
		"IDPEGA":   KunciInstans(id),
		"NOPOLIS":  h.Ambil(pt + "PolicyNo"),
		"PRODKE":   ProdKeNB,
		"TGL_PROD": h.Ambil(pt + "ProductionDate"),
		"USERNAME": akun,
	}
}

// BarisProduksi = `InsetTreatyInProd_Act` 4-11 + `InsertTreatyInProd_SQL`. `in` meniru halaman `InputTreaty`
// (CARIn teks), diisi berurutan seperti langkah aslinya; satu baris per panggilan RDB.
//
//	9   IsNewPolicyNonProp != 1: nilai polis (9.1), layer "0" bila TreatyType bukan XOL (9.2) atau layer polis
//	    (9.3), DUE TO US / DUE TO YOU (9.4-9.5), lalu SATU baris per baris SpreadingRiskList (9.6)
//	10  // berlabel, tak pernah jalan
//	11  IsNewPolicyNonProp == 1: per baris SpreadingRiskList x per TreatyXOLList x per ValueList, uang layer
//	    dikali `@divide(SharePercentage, 100, 20)`
func BarisProduksi(h *Halaman, id string) ([]Baris, error) {
	k := &kalkulator{}
	in := map[string]string{
		"CARI55": h.Ambil(pt + "ProductionDate"),         // 4
		"CARI11": tanggalSaja(h.Ambil(pt + "StartDate")), // 5; TO_DATE(.., 'YYYYMMDD')
		"CARI12": tanggalSaja(h.Ambil(pt + "EndDate")),
	}
	var out []Baris
	sisip := func() { out = append(out, barisProduksi(h, id, in)) }
	spreading := h.AmbilDaftar(DaftarSpreading)

	if !PolisNonPropBaru(h) { // 9
		for cari, prop := range map[string]string{ // 9.1
			"CARI31": "Currency", "CARI32": "PremiOgp", "CARI33": "RiCommOgp", "CARI34": "ResultOgp1",
			"CARI35": "OveriddingCommOgp", "CARI36": "ResultOgp2", "CARI38": "ExcessLoss", "CARI39": "PremiOnp",
			"CARI40": "RiCommOnp", "CARI41": "ResultOnp1", "CARI42": "OveriddingCommOnp", "CARI43": "ResultOnp2",
			"CARI44": "Deduction1", "CARI45": "Deduction2", "CARI47": "BalanceDueTo", "CARI48": "BalanceBeforePPH",
			"CARI49": "BalanceBeforeTax", "CARI28": "PPHValue", "CARI29": "PPNValue", "CARI54": "GuaranteeFund",
		} {
			in[cari] = h.Ambil(pt + prop)
		}
		if h.Ambil(pt+"TreatyType") != "XOL" { // 9.2
			in["CARI50"], in["CARI51"], in["CARI52"], in["CARI53"] = "0", "0", "0", "0"
		} else { // 9.3
			in["CARI50"], in["CARI51"] = h.Ambil(pt+"LayerType"), h.Ambil(pt+"Layer")
			in["CARI52"], in["CARI53"] = h.Ambil(pt+"LayerPartType"), h.Ambil(pt+"LayerPart")
		}
		switch dueTo := h.Ambil(pt + "DueTo"); { // 9.4, 9.5 (`.DueTo==1` / `.DueTo==0`)
		case samaDenganSatu(dueTo):
			in["CARI1"] = teksDueToUs
		case samaDenganNol(dueTo):
			in["CARI1"] = teksDueToYou
		}
		for _, s := range spreading { // 9.6
			in["CARI25"], in["CARI13"], in["CARI14"] = s["ClaimSpreaded"], s["PremiumSpreaded"], s["TreatyType"]
			in["CARI15"], in["CARI16"] = s["SharePercentage"], s["ClaimPercentage"]
			sisip() // 9.6.2 (kotak When `Local.year==0` tak dicentang; langkah 8 sudah keluar bila ada data)
		}
	}

	if PolisNonPropBaru(h) { // 11
		xol := h.AmbilDaftar(DaftarXOL)
		for _, s := range spreading {
			pct := s["SharePercentage"] // 11.1 local.SprdPct
			in["CARI14"] = s["TreatyType"]
			for _, c := range []string{"CARI33", "CARI34", "CARI35", "CARI36", "CARI38", "CARI40", "CARI41",
				"CARI42", "CARI43", "CARI44", "CARI45", "CARI16", "CARI25"} { // 11.2
				in[c] = "0"
			}
			in["CARI1"], in["CARI15"] = teksDueToUs, pct
			bagian := k.BagiBulat(angkaSel(k, "SharePercentage", pct), k.d("100"), 20)
			kali := func(prop, teks string) string { return utils.FormatDecimal(k.Kali(angkaSel(k, prop, teks), bagian)) }
			for i := range xol { // 11.3
				for _, v := range h.AmbilDaftar(JalurAnak(DaftarXOL, i+1, AnakLayerXOL)) { // 11.3.1
					in["CARI50"], in["CARI51"], in["CARI52"], in["CARI53"] = v["LayerType"], v["Layer"], v["LayerPartType"], v["LayerPart"]
					in["CARI32"], in["CARI31"] = kali("GrossPremi", v["GrossPremi"]), v["Currency"]
					in["CARI47"], in["CARI13"] = kali("DueToValue", v["DueToValue"]), kali("NetPremi", v["NetPremi"])
					in["CARI44"] = kali("Deduction", v["Deduction"])
					in["CARI28"], in["CARI29"] = v["PPHValue"], v["PPNValue"]
					sisip() // 11.3.1.2
				}
			}
		}
	}
	if k.err != nil {
		return nil, k.err
	}
	return out, nil
}

// barisProduksi = VALUES `InsertTreatyInProd_SQL` atas halaman dan `InputTreaty` saat itu.
func barisProduksi(h *Halaman, id string, in map[string]string) Baris {
	p := func(prop string) string { return h.Ambil(pt + prop) }
	return Baris{
		"IDPEGA":                      KunciInstans(id),
		"NOPOLIS":                     p("PolicyNo"),
		"BUSINESSCODE":                p("BizCode"),
		"CEDINGCO":                    strings.ReplaceAll(p("CedingCoName"), ";", ""),
		"SOB":                         p("SOBName"),
		"INSUREDNAME":                 p("InsuredName"),
		"TREATYGROUPID":               p("TreatyGroupID"),
		"TREATYGROUP":                 p("TreatyGroupName"),
		"DUE_TO":                      in["CARI1"],
		"STATEMENT_DATE":              in["CARI11"],
		"BEGINDATE":                   in["CARI11"],
		"ENDDATE":                     in["CARI12"],
		"CURR_ID":                     in["CARI31"],
		"UW_YEAR":                     p("TreatyYear"),
		"PREMI_OGP":                   in["CARI32"],
		"PERCENT_RI_COMM_OGP":         in["CARI33"],
		"RI_COMM_OGP":                 in["CARI34"],
		"PERCENT_OVERRIDING_COMM_OGP": in["CARI35"],
		"OVERRIDING_COMM_OGP":         in["CARI36"],
		"CLAIM":                       in["CARI25"],
		"EXCESS_LOSS":                 in["CARI38"],
		"PREMI_ONP":                   in["CARI39"],
		"PERCENT_RI_COMM_ONP":         in["CARI40"],
		"RI_COMM_ONP":                 in["CARI41"],
		"PERCENT_OVERRIDING_COMM_ONP": in["CARI42"],
		"OVERRIDING_COMM_ONP":         in["CARI43"],
		"DEDUCTION1":                  in["CARI44"],
		"DEDUCTION2":                  in["CARI45"],
		"NET_PREMIUM":                 in["CARI13"],
		"BALANCE_DUE_TO":              in["CARI47"],
		"MARKETINGOFFICERCODE":        NilaiQuotation(h, "MOID"),
		"NOOFFER":                     p("NoOffer"),
		"JN_REAS":                     in["CARI14"],
		"PCT_SHARE_PREMI":             in["CARI15"],
		"PCT_SHARE_CLAIM":             in["CARI16"],
		"TREATYTYPE":                  p("TreatyType"),
		"LAYERTYPE":                   in["CARI50"],
		"LAYER":                       in["CARI51"],
		"LAYERPARTTYPE":               in["CARI52"],
		"LAYERPART":                   in["CARI53"],
		"CLAIMTYPE":                   p("ClaimType"),
		"CLAIMPAYMENTTYPE":            p("ClaimPaymentType"),
		"QUARTER":                     p("Quartal"),
		"PPHVALUE":                    in["CARI28"],
		"PPNVALUE":                    in["CARI29"],
		"BALANCE_BEFORE_PPH":          in["CARI48"],
		"BALANCE_BEFORE_TAX":          in["CARI49"],
		"CEDINGCOID":                  p("CedingCo"),
		"SOBID":                       p("SOB"),
		"INSUREDID":                   p("InsuredID"),
		"JENIS_TREATY":                p("IsNewPolicyNonProp"),
		"QUARTER_YEAR":                p("YearOfQuartal"),
		"GUARANTEE_FUND":              in["CARI54"],
		"SHARE_VALUE":                 p("ShareValue"),
		"PROPORTIONALTYPE":            NilaiQuotation(h, "ProportionalType"),
		"OUTSTANDING_CLAIM":           p("OutstandingClaim"),
		"PROD_DATE":                   in["CARI55"],
		"STATEMENT_TYPE":              p("StatementType"),
	}
}

// BarisCapaian = `SetAchivementValue` + `SaveAchievementSQL` (InputData CARI1, CARI10..CARI25 -> parameter
// prosedur berurutan). Proporsional (`QuotationData.ProportionalType == "Proportional"`, langkah 4 + 6): SATU
// baris nilai polis. Selain itu (langkah 5): satu baris per TreatyXOLList, RICOMM kosong. NOPOLIS kosong di sini -
// diisi penulis dari json_polis (`GetPolicyNoByCaseId`).
func BarisCapaian(h *Halaman, id string) []Baris {
	p := func(prop string) string { return h.Ambil(pt + prop) }
	jenis := NilaiQuotation(h, "ProportionalType")
	dasar := func() Baris {
		return Baris{
			"IDPEGA": KunciInstans(id), "NOPOLIS": "", "NOOFFER": p("NoOffer"), "SOBNAME": p("SOBName"),
			"TREATYGROUPNAME": p("TreatyGroupName"), "TREATYTYPE": p("TreatyType"), "QUARTER": p("Quartal"),
			"QUARTERYEAR": p("YearOfQuartal"), "PAIDCLAIM": p("Claim"), "OUTSTANDINGCLAIM": p("OutstandingClaim"),
			"PROPORTIONALTYPE": jenis,
		}
	}
	if jenis == ProporsionalPenuh { // 4, 6
		b := dasar()
		b["IDCURRENCY"], b["CURRENCY"] = p("IDCurrency"), p("Currency")
		b["PREMIUM"], b["RICOMM"], b["BROKERAGE"], b["NETPREMIUM"] = p("PremiOgp"), p("ResultOgp1"), p("Deduction1"), p("NetPremium")
		return []Baris{b}
	}
	var out []Baris
	for _, x := range h.AmbilDaftar(DaftarXOL) { // 5 (5.1 When tak dicentang; 5.2 selalu di jalur ini)
		b := dasar()
		b["IDCURRENCY"], b["CURRENCY"] = x["IDCurrency"], x["Currency"]
		b["PREMIUM"], b["RICOMM"], b["BROKERAGE"], b["NETPREMIUM"] = x["GrossPremi"], "", x["Deduction"], x["NetPremi"]
		out = append(out, b)
	}
	return out
}

// bukanNolIsi = `(x != 0 && x != "")`: terisi dan bukan bilangan nol. Teks bukan bilangan terhitung terisi.
func bukanNolIsi(teks string) bool {
	s := strings.TrimSpace(teks)
	if s == "" {
		return false
	}
	d, err := utils.ParseDecimal(s)
	return err != nil || !d.IsZero()
}

// samaDenganNol = `.x == 0` atas teks bilangan (kosong bukan nol).
func samaDenganNol(teks string) bool {
	if !AdalahDesimal(teks) {
		return false
	}
	d, err := AngkaTeks("DueTo", teks)
	return err == nil && d.IsZero()
}

// tanggalSaja = `TO_DATE(x, 'YYYYMMDD')`: bagian tanggal saja; teks tak terbaca dibiarkan (penulis menolaknya).
func tanggalSaja(teks string) string {
	s := strings.TrimSpace(teks)
	if s == "" {
		return ""
	}
	t, err := utils.ParseTanggal(s)
	if err != nil {
		return s
	}
	return utils.FormatTanggal(t)
}

// angkaSel = `@toDecimal` sel kalkulasi (kosong = 0); galat dicatat sekali di kalkulator.
func angkaSel(k *kalkulator, prop, teks string) *apd.Decimal {
	d, err := AngkaTeks(prop, teks)
	if err != nil {
		k.catat(err)
		return apd.New(0, 0)
	}
	return d
}

func jika(syarat bool, ya, tidak string) string {
	if syarat {
		return ya
	}
	return tidak
}
