package models

// Untuk apa berkas ini: SIMPAN ENDORSEMEN KE TABEL PRODUKSI LAMA sesudah Dept Head menyetujui - Flow
// `InputAddendumTreatyIn` `Utility1` -> `Activity/SaveJsonPolisTreatyInEDM_Act`.
//
// Patokan (korpus `D:\XML\RNM_BRD\EDM Treaty In\`): Activity `SaveJsonPolisTreatyInEDM_Act` (ruleset 01-01-93,
// pxUpdateDateTime 2026-04-30), `SetEDMAchivementValue` (01-01-78, 2024-05-24), `InsetTreatyInProdAddendum_Act`
// (01-01-96, 2026-09-23); RDBList `SavePolisTreatyInEDM_SQL` (2025-11-20), `TreatyInSearchProdKe`,
// `GETTanggalClosing_SQL`, `SaveAchievementSQL` (2023-08-15, kata per kata = NB), `GetPolicyNoByCaseId`,
// `GetDataTreatyInProd_SQL`, `GetReinstypeIDbyName_SQL`, `InsertTreatyInProdEDMT_SQL` (2026-02-06). Isi prosedur
// `PEGA_JSON_POLIS_TREATYIN` / `InsertUpdateAchievment` tidak ada di korpus - dibaca tim NB dari ALL_SOURCE DEV:
// masing-masing SATU INSERT.
//
// `[keputusan work owner 06-10-2026]` (NB, berlaku sama): "JSON-nya tidak disimpan, tapi tetap insert kolom
// lainnya" - json_polis TANPA DATA_JSON; ACHIEVEMENT dan TREATYINPRODUCTION seperti Pega. Nol prosedur: INSERT
// langsung di transaksi Dept Head menyetujui.
//
// SaveJsonPolisTreatyInEDM_Act (18 aktif, 4 `//`; When TIDAK dicentang: 4, 4.1, 6):
//
//	1      Property-Remove SuggestDate/OperatorName/Suggest  ⛔ tidak ditiru (sama dengan NB): masukan layar usulan
//	2-3    TglProd = TANGGAL_CLOSING.TANGGAL, kosong -> 25    services (`HariClosing`) -> PrasimpanPolis
//	4      ProductionDate (4.1 sekarang, 4.2, 4.3)           PrasimpanPolis (`TanggalProduksiNomor`)
//	5      Policy.CaseID / Policy.OperatorID                 ⛔ halaman Policy hanya isi dokumen JSON (sama NB)
//	6      InputData.CARI1/21/2/4                            JSONPolis (CARI4 = EDMNo tak dibaca SQL mana pun)
//	7      IDNewBisnis = pyID                                PrasimpanMedan
//	8      ClaimType / ClaimPaymentType - When AKTIF (NB: tak dicentang)          PrasimpanMedan
//	9      // RealizationPopulateCedingSOB                    berlabel //, tak pernah jalan
//	10     CARI3 = @ASM.GetPageJSONString()                  ⛔ perintah WO: JSON tidak disimpan
//	11-12  TreatyInSearchProdKe -> CARI5 (PRODKE)            penulis (repository), di transaksi yang sama
//	13     SavePolisTreatyInEDM_SQL                          JSONPolis
//	14     ERRMSG = IDPEGAOUT                                tanpa padanan (keluaran prosedur)
//	15     Call SetEDMAchivementValue                        BarisCapaian
//	16-18  // INSERT KE JSON ERROR                           berlabel //
//	19     Call InsetTreatyInProdAddendum_Act                BarisProduksi
//
// `[penyimpangan sadar NB, disetujui WO 04-10-2026]` CARI21 (JSON_POLIS.TGL_PROD, langkah 6) dan CARI55
// (TREATYINPRODUCTION.PROD_DATE, InsetTreatyInProdAddendum_Act 4) diformat `dd/MM/yyyy hh:mm:ss` (jam 12) lalu
// dibaca `TO_DATE(.., 'DD/MM/YYYY HH24:MI:SS')` - ProductionDate sore tersimpan sebagai pagi. Di sini jam 24 apa
// adanya (models/usulan.go butir 3). ProductionDate digeser tanggal closing persis langkah 4.2-4.3.
//
// ⭐ Nilai uang produksi EDM = SELISIH (`PolicyTreatyIn.TreatyDifference`, `TreatyXOLDifferenceList`), bukan total
// generasi baru; pembatalan (`SetEDMTCancel`) menghasilkan selisih negatif dan ditulis apa adanya.

import (
	"strings"
	"time"

	"nusantarare/inti/backend/utils"
)

// Tabel produksi lama yang ditulis Utility1.
const (
	TabelJSONPolis = "JSON_POLIS"
	TabelProduksi  = "TREATYINPRODUCTION"
	TabelCapaian   = "ACHIEVEMENT"
)

// Teks DUE_TO - `InsetTreatyInProdAddendum_Act` 9.5, VERBATIM.
const teksDueToYou = "DUE TO YOU"

// ProporsionalPenuh - nilai `QuotationData.ProportionalType` jalur proporsional (`SetEDMAchivementValue` 3-5).
const ProporsionalPenuh = "Proportional"

// KolomJSONPolis - argumen `PEGA_JSON_POLIS_TREATYIN` (`SavePolisTreatyInEDM_SQL`, urutan argumen) selain
// DATA_JSON (perintah WO: tidak disimpan) dan TGL_INPUT (SYSDATE, ditulis SQL). Panjang: IDPEGA, NOPOLIS, PRODKE,
// TGL_PROD, USERNAME = kolom DEV (NB, ALL_TAB_COLUMNS 06-10-2026); NOENDORS VARCHAR2(100) dan PRODKE VARCHAR2(5)
// = DDL `JSON_POLIS.txt` (dikutip modul/nbfacin/docs/_DAFTAR-ISSUE-TERBUKA.md J-1 dan issues/24).
var KolomJSONPolis = []Kolom{
	kKode("IDPEGA", "IDPEGA", 50),
	kKode("NOPOLIS", "NOPOLIS", 100),
	kKode("NOENDORS", "NOENDORS", 100),
	kKode("PRODKE", "PRODKE", 5),
	kTglWaktu("TGL_PROD", "TGL_PROD"),
	kTeks("USERNAME", "USERNAME", 200),
}

// KolomProduksi - `InsertTreatyInProdEDMT_SQL`, urutan kolom INSERT apa adanya (58 kolom). Beda dengan NB
// `InsertTreatyInProd_SQL`: + NOENDORS (sesudah QUARTER), - TREATYGROUPID. Baris produksi memakai NAMA KOLOM sebagai
// kunci. Angka (`TO_NUMBER(REPLACE(x, ',', '.'))`) bergolongan uang/persen, tanggal `TO_DATE(x, 'YYYYMMDD')`
// bergolongan tanggal, PROD_DATE tanggal-waktu. Panjang = kolom DEV (NB), kecuali NOENDORS.
var KolomProduksi = []Kolom{
	kKode("IDPEGA", "IDPEGA", 100),
	kKode("NOPOLIS", "NOPOLIS", 100),
	kKode("BUSINESSCODE", "BUSINESSCODE", 100),
	kTeks("CEDINGCO", "CEDINGCO", 1000),
	kTeks("SOB", "SOB", 1000),
	kTeks("INSUREDNAME", "INSUREDNAME", 1000),
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
	// `PolicyTreatyIn.MarketingOfficer` = NAMA MO (layar.go langkah 7), NB: QuotationData.MOID
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
	// TREATYINPRODUCTION.NOENDORS VARCHAR2(100) - DEV ALL_TAB_COLUMNS 06-10-2026 (asisten utama, baca-saja)
	kKode("NOENDORS", "NOENDORS", 100),
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

// KolomCapaian - `POOLDATA.InsertUpdateAchievment` (INSERT ACHIEVEMENT), urutan parameter prosedur
// (`SaveAchievementSQL`: CARI1, CARI10..CARI25); TGL_PROD = SYSDATE (argumen ke-18), ditulis SQL. Panjang = DEV (NB).
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

// SimpananPolis - seluruh baris yang ditulis Utility1 untuk satu kasus, beserta masukan acuan yang dicari penulis
// di transaksi yang sama (RDB-List baca di tengah aktivitas Pega; models tidak menyentuh basis data).
type SimpananPolis struct {
	// IDPega - `pyWorkPage.pzInsKey` (`KunciInstans`).
	IDPega string
	// JSONPolis - satu baris json_polis (kunci = KolomJSONPolis). PRODKE KOSONG di sini: diisi penulis dengan
	// `TreatyInSearchProdKe` (langkah 11-12). Tidak ditulis bila IDPEGA sudah punya baris (prosedur Pega gagal pada
	// kunci utama, menelan galatnya, dan aktivitas berlanjut).
	JSONPolis Baris
	// Produksi - baris TREATYINPRODUCTION urutan sisip; ditulis HANYA bila IDPEGA belum punya baris
	// (`InsetTreatyInProdAddendum_Act` 6-8 "when there is data already, exit activity").
	Produksi []Baris
	// Capaian - baris ACHIEVEMENT. NOPOLIS diisi penulis dari json_polis kasus ini (`GetPolicyNoByCaseId`).
	Capaian []Baris
	// NotaJenisReasXOL - `TreatyIn.Share(1).SpreadingTypeXOL` yang ID-nya dicari penulis lewat
	// `GetReinstypeIDbyName_SQL` (`REINSURANCETYPE.ID WHERE NOTE = ..`), lalu dipasang ke JN_REAS setiap baris
	// produksi (`TerapkanJenisReasXOL`). Terisi HANYA di jalur NonProp bila `Share(1).SpreadingTypeIDXOL` kosong
	// (`InsetTreatyInProdAddendum_Act` 10.3); dicari hanya bila baris produksi memang ditulis (sesudah langkah 8).
	NotaJenisReasXOL string
}

// CacahProduksi - jumlah baris tabel produksi lama yang sudah dimiliki satu IDPEGA.
type CacahProduksi struct {
	JSONPolis, Capaian, Produksi int
}

// PrasimpanPolis = medan halaman yang Utility1 tulis sebelum halaman disimpan pada End:
// `SaveJsonPolisTreatyInEDM_Act` 4 (ProductionDate), 7, 8, dan `InsetTreatyInProdAddendum_Act` 9.1
// (`PrasimpanMedan`). `hariClosing` = langkah 2-3 (`TANGGAL_CLOSING.TANGGAL`, kosong -> 25; dibaca services).
func PrasimpanPolis(h *Halaman, id string, sekarang time.Time, hariClosing int) {
	// 4 blok REPEAT (When `OfferFacIn.PolicyData.PolicyNo==""` tak dicentang): 4.1 (When
	// `@PropertyHasValue(ProductionDate)` tak dicentang) ProductionDate = sekarang; 4.2 StatementDate di depan;
	// 4.3 lewat closing -> tanggal 1 bulan berikut 05:00 GMT. Sama dengan NB `SaveJsonPolisTreatyIn_Act` 4 (beda
	// hanya teks kondisi 4.1 yang nonaktif).
	var statement time.Time
	if s := strings.TrimSpace(h.Ambil(pt + "StatementDate")); s != "" {
		if t, err := time.ParseInLocation(utils.TanggalWaktu, s, sekarang.Location()); err == nil {
			statement = t
		}
	}
	h.Setel(pt+"ProductionDate", utils.FormatTanggalWaktu(TanggalProduksiNomor(sekarang, statement, hariClosing)))
	PrasimpanMedan(h, id)
}

// PrasimpanMedan = `PrasimpanPolis` TANPA langkah 4: ProductionDate yang tersimpan dipakai apa adanya (alat
// `simpanproduksi`).
func PrasimpanMedan(h *Halaman, id string) {
	// 7 `IDNewBisnis = pyWorkPage.pyID`
	h.Setel(pt+"IDNewBisnis", id)
	// 8 (halaman PolicyTreatyIn; When AKTIF `.ClaimType = "" ||.ClaimPaymentType = ""`, benar -> jalan, salah ->
	// lewati; `=` tunggal di kondisi dibaca sebagai pembanding): keduanya terisi = dibiarkan.
	// @if((.Claim!=0&&.Claim!="")||(.SalvageValue!=0&&.SalvageValue!=""), "SOA" / "Claim", "")
	if h.Ambil(pt+"ClaimType") == "" || h.Ambil(pt+"ClaimPaymentType") == "" {
		ada := bukanNolIsi(h.Ambil(pt+"Claim")) || bukanNolIsi(h.Ambil(pt+"SalvageValue"))
		h.Setel(pt+"ClaimType", jika(ada, "SOA", ""))
		h.Setel(pt+"ClaimPaymentType", jika(ada, "Claim", ""))
	}
	// InsetTreatyInProdAddendum_Act 9.1 - hanya di dalam blok 9 (`IsNewPolicyNonProp != 1`):
	// BizCode = @if(@PropertyHasValue(BizCode), BizCode, QuotationData.BusinessCode)
	if !PolisNonPropBaru(h) && h.Ambil(pt+"BizCode") == "" {
		h.Setel(pt+"BizCode", NilaiQuotation(h, "BusinessCode"))
	}
}

// SusunSimpananPolis menyusun baris Utility1 dari halaman sesudah `PrasimpanPolis`. `akun` =
// `OperatorID.pyUserIdentifier` (USERNAME json_polis). Tanda tangan = NB (dipanggil services); galat selalu nil
// di EDM - tidak ada aritmetika (NB mengalikan spreading di jalur NonProp, EDM tidak).
func SusunSimpananPolis(h *Halaman, id, akun string) (SimpananPolis, error) {
	return SimpananPolis{
		IDPega:           KunciInstans(id),
		JSONPolis:        JSONPolis(h, id, akun),
		Produksi:         BarisProduksi(h, id),
		Capaian:          BarisCapaian(h, id),
		NotaJenisReasXOL: NotaJenisReasXOL(h),
	}, nil
}

// TerapkanJenisReasXOL = `InsetTreatyInProdAddendum_Act` 10.3.3 `InputTreaty.CARI14 = ReinstypeOut.pxResults(1).CARI1`:
// `idJenis` hasil `GetReinstypeIDbyName_SQL` ("" bila tak ada baris) menjadi JN_REAS seluruh baris produksi. Hanya
// berlaku bila `NotaJenisReasXOL` terisi (jalur NonProp, satu CARI14 untuk semua lapisan).
func (s *SimpananPolis) TerapkanJenisReasXOL(idJenis string) {
	if s.NotaJenisReasXOL == "" {
		return
	}
	baru := make([]Baris, len(s.Produksi)) // daftar baru: baris dan larik milik pemanggil tidak diubah
	for i, b := range s.Produksi {
		salin := Baris{}
		for k, v := range b {
			salin[k] = v
		}
		salin["JN_REAS"] = idJenis
		baru[i] = salin
	}
	s.Produksi = baru
}

// JSONPolis = `SavePolisTreatyInEDM_SQL` -> `PEGA_JSON_POLIS_TREATYIN(pzInsKey, PolicyTreatyIn.PolicyNo,
// PolicyTreatyIn.EDMNo, CARI5, CARI21, OperatorID.pyUserIdentifier, CARI3)`. NOPOLIS = nomor polis INDUK (EDM
// menyimpan nomor endorsemen di NOENDORS); CARI21 = ProductionDate (langkah 6); CARI5 = PRODKE, diisi penulis
// (`select count(*) from pooldata.json_polis where Nopolis = CARI2`, langkah 11 - SEBELUM sisip langkah 13).
func JSONPolis(h *Halaman, id, akun string) Baris {
	return Baris{
		"IDPEGA":   KunciInstans(id),
		"NOPOLIS":  h.Ambil(pt + "PolicyNo"),
		"NOENDORS": h.Ambil(pt + "EDMNo"),
		"PRODKE":   "",
		"TGL_PROD": h.Ambil(pt + "ProductionDate"),
		"USERNAME": akun,
	}
}

// NotaJenisReasXOL = `InsetTreatyInProdAddendum_Act` 10.2-10.3 (jalur NonProp): CARI14 =
// `TreatyIn.Share(1).SpreadingTypeIDXOL`; bila kosong (When 1 `@PropertyHasValue(CARI14)` benar -> lewati) dan
// `Share(1).SpreadingTypeXOL` terisi (When 2), ID dicari menurut NOTE = SpreadingTypeXOL. Mengembalikan NOTE yang
// harus dicari, atau "" (jalur Prop, ID sudah ada, atau keduanya kosong).
func NotaJenisReasXOL(h *Halaman) string {
	if !PolisNonPropBaru(h) {
		return ""
	}
	s := prShareSatu(h)
	if s["SpreadingTypeIDXOL"] != "" || s["SpreadingTypeXOL"] == "" {
		return ""
	}
	return s["SpreadingTypeXOL"]
}

// BarisProduksi = `InsetTreatyInProdAddendum_Act` 4-10 + `InsertTreatyInProdEDMT_SQL`. `in` meniru halaman
// `InputTreaty` (CARIn teks), diisi berurutan seperti langkah aslinya; satu baris per panggilan RDB.
//
//	9   IsNewPolicyNonProp != 1, step page `PolicyTreatyIn.TreatyDifference`: nilai SELISIH (9.1), layer "0" bila
//	    TreatyType bukan XOL (9.2) atau layer TreatyDifference (9.3), DUE TO US / DUE TO YOU dari
//	    TreatyDifference.DueTo (9.4-9.5), lalu SATU baris per `TreatyDifference.SpreadingRiskList` (9.6)
//	10  IsNewPolicyNonProp == 1: per `TreatyXOLDifferenceList(c).ValueList(l)` (When 10.4.2.2 tak dicentang = setiap
//	    lapisan), DUE TO US, % 100, nilai lapisan APA ADANYA (NB mengalikan SharePercentage/100 per spreading)
//
// ⛔ `TreatyDifference.DueTo / Currency / GuaranteeFund / LayerType / Layer / LayerPartType / LayerPart` DIBACA blok
// 9 tetapi tidak satu pun rule korpus EDM menulisnya (katalog_selisih.go): DUE_TO, CURR_ID, GUARANTEE_FUND (dan
// LAYER* bila TreatyType XOL) jalur Prop = NULL apa adanya. Jalur NonProp tidak mengisi CARI39/48/49/54:
// PREMI_ONP, BALANCE_BEFORE_PPH/TAX, GUARANTEE_FUND = NULL.
func BarisProduksi(h *Halaman, id string) []Baris {
	in := map[string]string{
		"CARI55": h.Ambil(pt + "ProductionDate"),         // 4
		"CARI12": tanggalSaja(h.Ambil(pt + "EndDate")),   // 5; TO_DATE(.., 'YYYYMMDD')
		"CARI11": tanggalSaja(h.Ambil(pt + "StartDate")), // 5
	}
	var out []Baris
	sisip := func() { out = append(out, barisProduksi(h, id, in)) }

	if !PolisNonPropBaru(h) { // 9
		for cari, prop := range map[string]string{ // 9.1 (BizCode: PrasimpanMedan)
			"CARI31": "Currency", "CARI32": "PremiOgp", "CARI33": "RiCommOgp", "CARI34": "ResultOgp1",
			"CARI35": "OveriddingCommOgp", "CARI36": "ResultOgp2", "CARI38": "ExcessLoss", "CARI39": "PremiOnp",
			"CARI40": "RiCommOnp", "CARI41": "ResultOnp1", "CARI42": "OveriddingCommOnp", "CARI43": "ResultOnp2",
			"CARI44": "Deduction1", "CARI45": "Deduction2", "CARI47": "BalanceDueTo", "CARI48": "BalanceBeforePPH",
			"CARI49": "BalanceBeforeTax", "CARI28": "PPHValue", "CARI29": "PPNValue", "CARI54": "GuaranteeFund",
		} {
			in[cari] = h.Ambil(sd + prop)
		}
		if h.Ambil(pt+"TreatyType") != "XOL" { // 9.2
			in["CARI50"], in["CARI51"], in["CARI52"], in["CARI53"] = "0", "0", "0", "0"
		} else { // 9.3
			in["CARI50"], in["CARI51"] = h.Ambil(sd+"LayerType"), h.Ambil(sd+"Layer")
			in["CARI52"], in["CARI53"] = h.Ambil(sd+"LayerPartType"), h.Ambil(sd+"LayerPart")
		}
		switch dueTo := h.Ambil(sd + "DueTo"); { // 9.4, 9.5 (`.DueTo==1` / `.DueTo==0`)
		case samaDenganSatu(dueTo):
			in["CARI1"] = teksDueToUs
		case samaDenganNol(dueTo):
			in["CARI1"] = teksDueToYou
		}
		for _, s := range h.AmbilDaftar(DaftarSelisihSpreading) { // 9.6
			in["CARI25"], in["CARI13"], in["CARI14"] = s["ClaimSpreaded"], s["PremiumSpreaded"], s["TreatyType"]
			in["CARI15"], in["CARI16"] = s["SharePercentage"], s["ClaimPercentage"]
			sisip() // 9.6.2 (When `Local.year==0` tak dicentang; langkah 8 sudah keluar bila ada data)
		}
	}

	if PolisNonPropBaru(h) { // 10
		for _, c := range []string{"CARI33", "CARI34", "CARI35", "CARI36", "CARI38", "CARI40", "CARI41", "CARI42",
			"CARI43", "CARI44", "CARI45", "CARI16", "CARI25"} { // 10.1
			in[c] = "0"
		}
		in["CARI1"], in["CARI15"] = teksDueToUs, "100"
		// 10.2 CARI14 = Share(1).SpreadingTypeIDXOL; 10.3 (cari ID menurut NOTE) lewat NotaJenisReasXOL + penulis
		in["CARI14"] = prShareSatu(h)["SpreadingTypeIDXOL"]
		// 10.4 (10.4.1 local.tempcurrency: nol pembaca)
		for c := range h.AmbilDaftar(DaftarSelisihXOL) {
			for _, v := range h.AmbilDaftar(JalurAnak(DaftarSelisihXOL, c+1, "ValueList")) { // 10.4.2
				in["CARI50"], in["CARI51"], in["CARI52"], in["CARI53"] = v["LayerType"], v["Layer"], v["LayerPartType"], v["LayerPart"]
				in["CARI32"], in["CARI31"], in["CARI47"] = v["GrossPremi"], v["Currency"], v["DueToValue"]
				in["CARI13"], in["CARI44"] = v["NetPremi"], v["Deduction"]
				in["CARI28"], in["CARI29"] = v["PPHValue"], v["PPNValue"]
				sisip() // 10.4.2.2 (When `.Currency==local.tempcurrency` tak dicentang)
			}
		}
	}
	return out
}

// barisProduksi = VALUES `InsertTreatyInProdEDMT_SQL` atas halaman dan `InputTreaty` saat itu.
func barisProduksi(h *Halaman, id string, in map[string]string) Baris {
	p := func(prop string) string { return h.Ambil(pt + prop) }
	return Baris{
		"IDPEGA":                      KunciInstans(id),
		"NOPOLIS":                     p("PolicyNo"),
		"BUSINESSCODE":                p("BizCode"),
		"CEDINGCO":                    strings.ReplaceAll(p("CedingCoName"), ";", ""),
		"SOB":                         p("SOBName"),
		"INSUREDNAME":                 p("InsuredName"),
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
		"MARKETINGOFFICERCODE":        p("MarketingOfficer"),
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
		"NOENDORS":                    p("EDMNo"),
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

// BarisCapaian = `SetEDMAchivementValue` + `SaveAchievementSQL` (InputData CARI1, CARI10..CARI25 -> parameter
// prosedur berurutan). Proporsional (`QuotationData.ProportionalType == "Proportional"`, langkah 3 + 5): SATU baris
// nilai `TreatyDifference`. Selain itu (langkah 4, When 4.1 tak dicentang, 4.2 benar di jalur ini): satu baris per
// induk `TreatyXOLDifferenceList` (per mata uang), RICOMM kosong, PAIDCLAIM = `PolicyTreatyIn.Claim` (BUKAN
// selisih). QUARTERYEAR = `TreatyYear` (NB: YearOfQuartal). NOPOLIS kosong di sini - diisi penulis dari json_polis
// (`GetPolicyNoByCaseId`, langkah 2). Tanpa Page-Remove InputData (NB langkah 1): CARI10..CARI25 ditimpa seluruhnya
// sebelum setiap panggilan, jadi sisa CARI SaveJson tidak terbawa.
func BarisCapaian(h *Halaman, id string) []Baris {
	p := func(prop string) string { return h.Ambil(pt + prop) }
	jenis := NilaiQuotation(h, "ProportionalType")
	dasar := func() Baris {
		return Baris{
			"IDPEGA": KunciInstans(id), "NOPOLIS": "", "NOOFFER": p("NoOffer"), "SOBNAME": p("SOBName"),
			"TREATYGROUPNAME": p("TreatyGroupName"), "TREATYTYPE": p("TreatyType"), "QUARTER": p("Quartal"),
			"QUARTERYEAR": p("TreatyYear"), "OUTSTANDINGCLAIM": p("OutstandingClaim"), "PROPORTIONALTYPE": jenis,
		}
	}
	if jenis == ProporsionalPenuh { // 3, 5
		b := dasar()
		b["IDCURRENCY"], b["CURRENCY"] = p("IDCurrency"), p("Currency")
		b["PREMIUM"], b["RICOMM"] = h.Ambil(sd+"PremiOgp"), h.Ambil(sd+"ResultOgp1")
		b["BROKERAGE"], b["NETPREMIUM"] = h.Ambil(sd+"Deduction1"), h.Ambil(sd+"NetPremium")
		b["PAIDCLAIM"] = h.Ambil(sd + "Claim")
		return []Baris{b}
	}
	var out []Baris
	for _, x := range h.AmbilDaftar(DaftarSelisihXOL) { // 4 -> 4.1, 4.2
		b := dasar()
		b["IDCURRENCY"], b["CURRENCY"] = x["IDCurrency"], x["Currency"]
		b["PREMIUM"], b["RICOMM"], b["BROKERAGE"], b["NETPREMIUM"] = x["GrossPremi"], "", x["Deduction"], x["NetPremi"]
		b["PAIDCLAIM"] = p("Claim")
		out = append(out, b)
	}
	return out
}

// prShareSatu - baris `TreatyIn.Share(1)` master (kosong bila daftar kosong; Pega membaca properti baris yang tak
// ada sebagai "").
func prShareSatu(h *Halaman) Baris {
	if s := mDaftar(h, "Share"); len(s) > 0 {
		return s[0]
	}
	return Baris{}
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

func jika(syarat bool, ya, tidak string) string {
	if syarat {
		return ya
	}
	return tidak
}
