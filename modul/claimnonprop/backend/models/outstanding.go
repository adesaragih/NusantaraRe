package models

// Untuk apa berkas ini: PORT TOMBOL PENULIS OS_AKSEPTASI_KLAIM Non Prop - "Save to issue RNM" (`SaveDataToOSAksep_Act`,
// STS 0), "Save To OS" (`SaveToOS`, STS 0), "Close Claim" (`CloseClaimTNonProp`, STS 4), "Save Previously Paid"
// (`SaveCNPLayerList_Act`, STS 5), dan JSON halaman Pega (`@ASM.GetPageJSONString()`).
//
// Bagian murni (halaman, pesan, rencana baris OS, DATA_JSON) ada di sini; nomor (penomor), selisih terhadap baris OS
// tersimpan, dan tulisan OS / JSON_KLAIM / outbox dijalankan services lewat repository dalam SATU transaksi.
//
// Procedure `PEGA_JSON_OS_AKSEP_KLAIMTNP` (ALL_SOURCE DEV 09-10-2026) ditulis ulang tanpa procedure: SELALU INSERT
// (cabang UPDATE-nya dikomentari) kolom CASEID, NOCLAIM, MASTERID, DATA_JSON, TANGGAL (hari ini), NOPOLIS, STS_REJECT,
// STS_KONVERSI; TGL_PROD diisi trigger `TRG_TLG_PROD_OS_AKSEPTASI` (BEFORE INSERT).

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"
)

// Pesan VERBATIM.
const (
	PesanBreakQSKosong   = "SpreadingBreakQS tidak boleh kosong"                  // SaveDataToOSAksep_Act 3
	PesanProvinsiKosong  = "Province tidak boleh Kosong"                          // SaveDataToOSAksep_Act 3
	PesanNoPolisKosong   = "Nomor Polis tidak boleh kosong"                       // SaveDataToOSAksep_Act 2 / SaveToOS 4
	PesanKomitePending   = "Can not close claim, there is adjustment in comitee!" // CloseClaimTNonProp 1.1
	PesanCWPAdaAkseptasi = "Claim tidak bisa di Close, Sudah ada akseptasi"       // CloseClaimNP_preAct 2
)

// Teks riwayat.
const (
	// TeksTutupKlaim - CloseClaimTNonProp langkah 4. PERBAIKAN OQ-CNP-05 butir 6 (PARITAS `[penyimpangan sadar]`): XML
	// menulisnya ke `DataChronology.CARI12` sedangkan InsertChronology_DT membaca `CARI1`, sehingga kronologi Pega
	// berisi teks sisa aksi sebelumnya; di sini teks ini yang tercatat.
	TeksTutupKlaim = "Finish Adjustment (Close Claim)"
	TeksKirimCWP   = "Send to Committe (Acceptation)" // CloseClaimNP_preAct 2
)

// Status OS_AKSEPTASI_KLAIM.STS_REJECT (= parameter CARI10 procedure, PINDAI §5).
const (
	StsOSOutstanding = "0" // Save to issue RNM, Save To OS
	StsOSFinal       = "4" // Close Claim
	StsOSDibayar     = "5" // Save Previously Paid
)

// KelasOSAkseptasi - pxObjClass halaman InputParamOs / InputParamOsCNP (Pages & Classes activity; DATA_JSON DEV).
const KelasOSAkseptasi = "ASM-FW-GCNMFW-Data-osAkseptasi"

// KelasAdjustment - pxObjClass anggota CNPLayerList / CNPCurrencyList (DATA_JSON DEV).
const KelasAdjustment = "ASM-FW-GCNMFW-Data-Adjustment"

// Awalan nomor.
const (
	JenisNomorKlaim         = "K" // SaveDataToOSAksep_Act 15.6.6 `ParamSeq.HASIL3+"K"`
	TipeKodeProduksiNonLife = "NONLIFE"
)

// RakitNomorKlaim = SaveDataToOSAksep_Act 15.6.8: `ParamSeq.CARI2 + BusinessOldId + "." + HASIL1 + ".TX" + HASIL2`
// (CARI2 = kode NONLIFE + "K", HASIL1 = MM.YYYY, HASIL2 = urut 5 digit).
func RakitNomorKlaim(jenis, oldID, mmYYYY string, urut int) string {
	return fmt.Sprintf("%s%s.%s.TX%05d", jenis, oldID, mmYYYY, urut)
}

// HalamanJSON - satu halaman untuk `@ASM.GetPageJSONString()`: nilai teks + daftar halaman (PageList).
type HalamanJSON struct {
	Nilai  map[string]string
	Daftar []DaftarJSON
}

// DaftarJSON - satu PageList bernama.
type DaftarJSON struct {
	Nama string
	Isi  []HalamanJSON
}

// JSONHalamanPega = `@ASM.GetPageJSONString()`. Fungsinya tidak diekspor; bentuknya dibaca dari DATA_JSON
// OS_AKSEPTASI_KLAIM `CLMNP-` DEV 09-10-2026 (salinan pola Claim Prop `JSONHalamanPega`):
//   - "{" LF, pasangan pertama, setiap pasangan berikutnya diawali LF ",", ditutup LF "}"; halaman puncak diakhiri LF;
//   - nilai teks dulu, urut kunci tanpa membedakan huruf besar; properti kosong tidak ditulis;
//   - PageList sesudah nilai teks: `"Nama":[ ` LF anggota LF `] ` (anggota dipisah LF "," - `[inferensi]`, DEV hanya
//     memuat daftar beranggota satu).
func JSONHalamanPega(p HalamanJSON) string {
	return tulisHalamanJSON(p) + "\n"
}

func tulisHalamanJSON(p HalamanJSON) string {
	kunci := make([]string, 0, len(p.Nilai))
	for k, v := range p.Nilai {
		if v != "" {
			kunci = append(kunci, k)
		}
	}
	sort.Slice(kunci, func(i, j int) bool {
		a, b := strings.ToLower(kunci[i]), strings.ToLower(kunci[j])
		if a != b {
			return a < b
		}
		return kunci[i] < kunci[j]
	})
	var bagian []string
	for _, k := range kunci {
		bagian = append(bagian, teksJSON(k)+":"+teksJSON(p.Nilai[k]))
	}
	for _, d := range p.Daftar {
		if len(d.Isi) == 0 {
			continue
		}
		var isi []string
		for _, a := range d.Isi {
			isi = append(isi, tulisHalamanJSON(a))
		}
		bagian = append(bagian, teksJSON(d.Nama)+":[ \n"+strings.Join(isi, "\n,")+"\n] ")
	}
	return "{\n" + strings.Join(bagian, "\n,") + "\n}"
}

// teksJSON - literal teks JSON tanpa escape HTML.
func teksJSON(s string) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(s)
	return strings.TrimSuffix(buf.String(), "\n")
}

// BarisOS - satu baris OS_AKSEPTASI_KLAIM (kolom yang ditulis procedure). CASEID = `KunciInstans` (ID kasus apa adanya).
type BarisOS struct {
	CaseID, NoClaim, MasterID, NoPolis string
	StsReject, StsKonversi             string
	DataJSON                           string
	// TypeLoss, Currency - kunci baris di DATA_JSON (dipakai menghitung selisih, tidak ditulis ke kolom sendiri).
	TypeLoss, Currency string
}

// NilaiOS - Σ nilai baris OS STS 0 tersimpan per layer x mata uang (`GetDataOS` / `GetDataCNPOS`, alias SQL Value /
// GrossValue / Adjusterfee / CNPOthersFee / Salvage).
type NilaiOS struct {
	Value, GrossValue, Adjusterfee, CNPOthersFee, Salvage string
}

// KelompokLayerOS = SaveDataToOSAksep_Act 15.7 / SaveToOS 6.2: SpreadingRisk dijumlah per (Currency, TreatyType)
// (AdjusterFee, ClaimSpreaded, Salvage, CNPOthersFee, ClaimEstimation, AdjClaimValue); baris UR dibuang (15.8.1 F:4).
func KelompokLayerOS(h *Halaman) []Baris {
	var kal Kalkulator
	var out []Baris
	for _, s := range h.AmbilDaftar(DaftarXOL) {
		var t Baris
		for _, o := range out {
			if o["Currency"] == s["Currency"] && o["TreatyType"] == s["TreatyType"] {
				t = o
			}
		}
		if t == nil {
			out = append(out, s.Salin())
			continue
		}
		for _, p := range []string{"AdjusterFee", "ClaimSpreaded", "Salvage", "CNPOthersFee", "ClaimEstimation", "AdjClaimValue"} {
			t[p] = Teks(kal.Tambah(kal.B(t, p), kal.B(s, p)))
		}
	}
	var tanpaUR []Baris
	for _, b := range out {
		if b["TreatyName"] != TreatyUR {
			tanpaUR = append(tanpaUR, b)
		}
	}
	return tanpaUR
}

// ParamOSLayer = SaveDataToOSAksep_Act 15.8.1 / SaveToOS 6.3.1: halaman InputParamOs satu layer x mata uang.
func ParamOSLayer(h *Halaman, l Baris) map[string]string {
	var kal Kalkulator
	return map[string]string{
		"CauseOfLoss": h.Ambil(CD + "CauseOfLoss"), "CauseOfLossID": h.Ambil(CD + "CauseOfLossID"),
		"Currency": l["Currency"], "CurrencyID": l["CurrencyID"], "EstimationDate": "",
		"GrossValue": Teks(kal.Tambah(kal.B(l, "ClaimEstimation"), kal.B(l, "AdjClaimValue"))),
		"KursValue":  l["Kurs"], "NoClaim": h.Ambil(CD + "NoClaim"), "TypeLoss": l["TreatyName"],
		"TypeLossID": l["TreatyType"], "Value": l["ClaimSpreaded"], "Adjusterfee": l["AdjusterFee"],
		"Salvage": l["Salvage"], "CNPOthersFee": l["CNPOthersFee"], "IDMasterTreaty": h.Ambil(TM + "ID"),
	}
}

// KurangiOS = SaveDataToOSAksep_Act 15.8.3 / SaveToOS 6.3.4: selisih terhadap Σ baris OS STS 0 tersimpan; Type 0,
// PersenRNM = RNM Share.
func KurangiOS(h *Halaman, p map[string]string, lama NilaiOS) error {
	var kal Kalkulator
	kurang := func(k, v string) { p[k] = Teks(kal.Kurang(kal.Teks(k, p[k]), kal.Teks(k, v))) }
	kurang("CNPOthersFee", lama.CNPOthersFee)
	kurang("Adjusterfee", lama.Adjusterfee)
	kurang("Value", lama.Value)
	kurang("Salvage", lama.Salvage)
	kurang("GrossValue", lama.GrossValue)
	p["Type"] = "0"
	p["PersenRNM"] = h.Ambil(TM + "RNMShare")
	return kal.Galat()
}

// SusunBarisOS - baris OS dari halaman InputParamOs (`InputData.CARI1 = @ASM.GetPageJSONString()`).
func SusunBarisOS(h *Halaman, id, sts string, p map[string]string) BarisOS {
	nilai := map[string]string{"pxObjClass": KelasOSAkseptasi}
	for k, v := range p {
		nilai[k] = v
	}
	return BarisOS{CaseID: KunciInstans(id), NoClaim: h.Ambil(CD + "NoClaim"), MasterID: h.Ambil(CD + "IDMaster"),
		NoPolis: h.Ambil(CD + "PolicyData.PolicyNo"), StsReject: sts, DataJSON: JSONHalamanPega(HalamanJSON{Nilai: nilai}),
		TypeLoss: p["TypeLoss"], Currency: p["Currency"]}
}

// NilaiNol - teks bernilai nol (kosong = 0).
func NilaiNol(s string) bool {
	var kal Kalkulator
	d := kal.Teks("x", s)
	return kal.Galat() == nil && d.IsZero()
}

// PeriksaSimpanOS = SaveDataToOSAksep_Act langkah 3-13 (Save to issue RNM). Mengembalikan false bila activity keluar.
// ⚠️ Kelainan XML (PARITAS, OQ): langkah 8 menolak bila `ClaimData.QuotationData.BusinessOldId` kosong, padahal
// pengisinya (langkah 15.2-15.3, GetDataBusiness_SQL menurut nama bisnis) baru jalan SESUDAHNYA - Save to issue RNM
// tidak pernah lolos untuk kasus baru. Di sini OLDID dibaca lebih dulu (`oldID`, dari `OfferFacIn.QuotationData
// .BusinessName` - `ClaimData.QuotationData.BusinessName` tidak punya penulis).
func PeriksaSimpanOS(h *Halaman, oldID, noClaimOS string) bool {
	if len(h.AmbilDaftar(DaftarBreakQS)) == 0 { // 4 (trans T:6)
		h.TambahPesan("", PesanBreakQSKosong)
		return false
	}
	if h.Ambil(CD+"Province") == "" { // 5
		h.TambahPesan("", PesanProvinsiKosong)
	}
	if h.Ambil(CD+"QuotationData.BusinessOldId") == "" && oldID != "" {
		h.Setel(CD+"QuotationData.BusinessOldId", oldID)
	}
	if h.Ambil(OQ+"BusinessName") == "" || h.Ambil(CD+"QuotationData.BusinessOldId") == "" { // 8 (trans T:6)
		h.TambahPesan("", PesanBisnisKosong)
		return false
	}
	if h.Ambil(CD+"NoClaim") == "" && noClaimOS != "" { // 10, 12
		h.Setel(CD+"NoClaim", noClaimOS)
	}
	return !h.AdaPesan() // 13
}

// TandaiOutstanding = SaveDataToOSAksep_Act langkah 16-22: IsOutstanding = 1, baris ListClaimAmount / Loss Allocation /
// XOL Allocation / Spreading List terkunci (CNPFlagOuts = 1).
func TandaiOutstanding(h *Halaman) {
	h.Setel("IsOutstanding", "1")
	for _, d := range []string{DaftarClaimAmount, DaftarLossAlloc, DaftarXOL, DaftarSpreading} {
		for _, b := range h.AmbilDaftar(d) {
			b["CNPFlagOuts"] = "1"
		}
	}
}

// ParamOSSaveToOS = SaveToOS langkah 6.3.1-6.3.7 untuk satu layer: `trtSebelum` = `local.TrtName` (diisi langkah 6.3.6
// SESUDAH dipakai langkah 6.3.2 - jadi nama layer iterasi SEBELUMNYA; ditiru apa adanya, PARITAS). Mengembalikan false
// bila langkah 6.3.2.1.1 melompat ke label B (Value 0 pada jalur premi aktual) - baris OS tidak ditulis.
func ParamOSSaveToOS(h *Halaman, l Baris, trtSebelum string, lama NilaiOS) (map[string]string, bool, error) {
	var kal Kalkulator
	p := ParamOSLayer(h, l)
	aktual := h.Ambil("FlagActualPremium") == "true"
	if aktual { // 6.3.2
		for i := range h.AmbilDaftar(DaftarAdjustment) {
			for _, s := range h.AmbilDaftar(JalurAdj(i+1, AnakXOL)) {
				if s["TreatyName"] == TreatyUR || trtSebelum != s["TreatyName"] {
					continue
				}
				p["CNPOthersFee"] = Teks(kal.Kurang(kal.Teks("o", p["CNPOthersFee"]), kal.B(s, "CNPOthersFee")))
				p["Adjusterfee"] = Teks(kal.Kurang(kal.Teks("a", p["Adjusterfee"]), kal.B(s, "AdjusterFee")))
				p["Value"] = Teks(kal.Kurang(kal.Teks("v", p["Value"]), kal.B(s, "ClaimSpreaded")))
				p["Salvage"] = Teks(kal.Kurang(kal.Teks("s", p["Salvage"]), kal.B(s, "Salvage")))
				p["GrossValue"] = Teks(kal.Kurang(kal.Teks("g", p["GrossValue"]), kal.B(s, "TotalClaim")))
				p["Type"], p["PersenRNM"] = "0", h.Ambil(TM+"RNMShare")
				if NilaiNol(p["Value"]) { // trans: Value 0 -> label B
					return p, false, kal.Galat()
				}
			}
		}
	}
	if err := KurangiOS(h, p, lama); err != nil { // 6.3.3-6.3.4
		return nil, false, err
	}
	if aktual { // 6.3.7
		p["Value"] = Teks(kal.Kurang(kal.B(l, "ClaimSpreaded"), kal.Teks("v", p["Value"])))
		p["GrossValue"] = Teks(kal.Kurang(kal.B(l, "TotalClaim"), kal.Teks("g", p["GrossValue"])))
	}
	return p, true, kal.Galat()
}

// SelesaiSaveToOS = SaveToOS langkah 8-9: IsSaveToOs = 1, ListClaimAmount terkunci.
func SelesaiSaveToOS(h *Halaman) {
	h.Setel("IsSaveToOs", "1")
	for _, b := range h.AmbilDaftar(DaftarClaimAmount) {
		b["CNPFlagOuts"] = "1"
	}
}

// PeriksaTutupKlaim = CloseClaimTNonProp langkah 1-3: akseptasi masih di komite (AcceptanceStatus "0") -> pesan, keluar.
func PeriksaTutupKlaim(h *Halaman) bool {
	for _, b := range h.AmbilDaftar(DaftarAdjustment) {
		if b["AcceptanceStatus"] == "0" {
			h.TambahPesan("", PesanKomitePending)
			return false
		}
	}
	return true
}

// TutupKlaim = CloseClaimTNonProp langkah 1-4 + 6 (tombol Yes local action Close Claim): proteksi akseptasi di komite,
// lalu kronologi `TeksTutupKlaim` (perbaikan butir 6). Baris OS STS 4 (5), JSON_KLAIM (7), outbox
// insertClaimFinalOrClosed_NP (8), dan penutupan kasus Resolved-Completed (9) dijalankan services dalam transaksi yang
// sama. Mengembalikan false bila keluar karena proteksi (pesan terpasang).
func TutupKlaim(k *Konteks, h *Halaman) bool {
	if !PeriksaTutupKlaim(h) {
		return false
	}
	k.Riwayat(h, TeksTutupKlaim)
	return true
}

// PraCWP = CloseClaimNP_preAct (pra-proses local action CloseClaimNP): `IsAcceptation` = AcceptanceStatus akseptasi
// PERTAMA (langkah 2, ditiru apa adanya - menimpa penanda tersimpan), inisial = nama pelaku, tanggal = hari ini,
// kronologi `TeksKirimCWP` (3), pesan bila IsAcceptation 1 (4). `nama` = `OperatorID.pyLabel`.
func PraCWP(k *Konteks, h *Halaman, nama string) {
	status := ""
	if d := h.AmbilDaftar(DaftarAdjustment); len(d) > 0 {
		status = d[0]["AcceptanceStatus"]
	}
	h.Setel("IsAcceptation", status)
	h.Setel(JalurPICKomite, nama)
	h.Setel(JalurTanggalKomite, k.Hari())
	k.Riwayat(h, TeksKirimCWP)
	if status == "1" {
		h.TambahPesan("", PesanCWPAdaAkseptasi)
	}
}

// ParamOSTutup = CloseClaimTNonProp langkah 5.1: halaman InputParamOs penutupan (STS 4).
func ParamOSTutup(h *Halaman) map[string]string {
	return map[string]string{"CauseOfLoss": h.Ambil(CD + "CauseOfLoss"), "CauseOfLossID": h.Ambil(CD + "CauseOfLossID"),
		"NoClaim": h.Ambil(CD + "NoClaim"), "IDMasterTreaty": h.Ambil(TM + "ID")}
}

// LayerListAkseptasi = SaveCNPLayerList_Act langkah 2-3 (dan CreateChildKomiteCNP_Act 20-23): `CNPLayerList` akseptasi
// `n` - satu layer per nama XoL (non-UR), satu `CNPCurrencyList` per baris mata uang; Previously Calculated dikurangkan
// (langkah 3).
func LayerListAkseptasi(h *Halaman, n int) ([]HalamanJSON, error) {
	var kal Kalkulator
	rnm := kal.BagiPega(kal.H(h, TM+"RNMShare"), k100())
	type layer struct {
		xol, xolID string
		mu         []map[string]string
	}
	var daftar []*layer
	for _, s := range h.AmbilDaftar(JalurAdj(n, AnakXOL)) {
		if s["TreatyName"] == TreatyUR {
			continue
		}
		pct := kal.BagiPega(kal.B(s, "ClaimPercentage"), k100())
		var l *layer
		for _, d := range daftar {
			if d.xol == s["TreatyName"] {
				l = d
			}
		}
		mu := map[string]string{"pxObjClass": KelasAdjustment, "AdjusterFeeRNM": s["AdjusterFee"],
			"CNPLimit": s["CNPLimit"], "CNPMindep": s["CNPMDP"], "CNPPctReinstate": s["CNPPctReinstate"],
			"CNPReinstatement": s["CNPReinstatement"], "CNPReinstatementRNM": s["CNPReinstatementRNM"],
			"Currency": s["Currency"], "CurrencyID": s["CurrencyID"], "GrossAdjustment": s["ClaimEstimation"],
			"GrossValue": s["ClaimSpreaded"], "KursIDR": s["Kurs"], "PersenRNM": s["ClaimPercentage"],
			"SalvageRNM": s["Salvage"], "TotalXOL": s["TotalClaim"], "TotalXOLGross": s["TotalClaim"],
			"CNPOthersFeeRNM": s["CNPOthersFee"]}
		total := kal.Kali(kal.B(s, "TotalClaim"), pct)
		mu["TotalXOLLayerRNM"], mu["TotalXOLRNM"] = Teks(total), Teks(total)
		if l == nil { // 2.3: layer baru - pembagi RNM Share master
			mu["AdjusterFeeValue"] = Teks(kal.BagiPega(kal.B(s, "AdjusterFee"), rnm))
			mu["SalvageValue"] = Teks(kal.BagiPega(kal.B(s, "Salvage"), rnm))
			mu["CNPOthersFee"] = Teks(kal.BagiPega(kal.B(s, "CNPOthersFee"), rnm))
			daftar = append(daftar, &layer{xol: s["TreatyName"], xolID: s["TreatyType"], mu: []map[string]string{mu}})
			continue
		}
		mu["AdjusterFeeValue"] = Teks(kal.BagiPega(kal.B(s, "AdjusterFee"), pct)) // 2.2.2: pembagi persen baris
		mu["SalvageValue"] = Teks(kal.BagiPega(kal.B(s, "Salvage"), pct))
		mu["CNPOthersFee"] = Teks(kal.BagiPega(kal.B(s, "CNPOthersFee"), pct))
		l.mu = append(l.mu, mu)
	}
	dibayar := h.AmbilDaftar(JalurAdj(n, AnakXOLDibayar))
	if len(dibayar) > 0 { // 3
		for _, l := range daftar {
			for _, mu := range l.mu {
				var p Baris
				for _, d := range dibayar {
					if d["Currency"] == mu["Currency"] && d["TreatyName"] == l.xol {
						p = d
					}
				}
				if p == nil {
					continue
				}
				kurang := func(k string, v *apd.Decimal) { mu[k] = Teks(kal.Kurang(kal.Teks(k, mu[k]), v)) }
				kurang("AdjusterFeeRNM", kal.B(p, "AdjusterFee"))
				kurang("AdjusterFeeValue", kal.BagiPega(kal.B(p, "AdjusterFee"), rnm))
				kurang("CNPOthersFee", kal.BagiPega(kal.B(p, "CNPOthersFee"), rnm))
				kurang("CNPOthersFeeRNM", kal.B(p, "CNPOthersFee"))
				kurang("CNPReinstatement", kal.B(p, "CNPReinstatement"))
				kurang("CNPReinstatementRNM", kal.B(p, "CNPReinstatementRNM"))
				kurang("GrossAdjustment", kal.B(p, "ClaimAmountAdjust"))
				kurang("GrossValue", kal.B(p, "ClaimSpreaded"))
				kurang("SalvageRNM", kal.B(p, "Salvage"))
				kurang("SalvageValue", kal.BagiPega(kal.B(p, "Salvage"), rnm))
				g := kal.Tambah(kal.Kurang(kal.Tambah(kal.Teks("g", mu["GrossAdjustment"]), kal.Teks("a", mu["AdjusterFeeValue"])),
					kal.Teks("s", mu["SalvageValue"])), kal.Teks("o", mu["CNPOthersFee"]))
				mu["TotalXOL"], mu["TotalXOLGross"] = Teks(g), Teks(g)
				r := kal.Kurang(kal.Tambah(kal.Tambah(kal.Teks("g", mu["GrossValue"]), kal.Teks("a", mu["AdjusterFeeRNM"])),
					kal.Teks("o", mu["CNPOthersFeeRNM"])), kal.Teks("s", mu["SalvageRNM"]))
				mu["TotalXOLLayerRNM"], mu["TotalXOLRNM"] = Teks(r), Teks(r)
			}
		}
	}
	var out []HalamanJSON
	for _, l := range daftar {
		var mu []HalamanJSON
		for _, m := range l.mu {
			mu = append(mu, HalamanJSON{Nilai: m})
		}
		out = append(out, HalamanJSON{Nilai: map[string]string{"pxObjClass": KelasAdjustment, "XOL": l.xol, "XOLID": l.xolID},
			Daftar: []DaftarJSON{{Nama: "CNPCurrencyList", Isi: mu}}})
	}
	return out, kal.Galat()
}

// BarisOSDibayar = SaveCNPLayerList_Act langkah 4-6 (Save Previously Paid, akseptasi `n`): satu baris OS STS 5 berisi
// halaman InputParamOsCNP (CNPLayerList, CauseOfLoss, CauseOfLossID, NoClaim, PersenRNM, Type "1").
func BarisOSDibayar(h *Halaman, id string, n int) (BarisOS, error) {
	layer, err := LayerListAkseptasi(h, n)
	if err != nil {
		return BarisOS{}, err
	}
	p := HalamanJSON{Nilai: map[string]string{"pxObjClass": KelasOSAkseptasi, "CauseOfLoss": h.Ambil(CD + "CauseOfLoss"),
		"CauseOfLossID": h.Ambil(CD + "CauseOfLossID"), "NoClaim": h.Ambil(CD + "NoClaim"),
		"PersenRNM": h.Ambil(TM + "RNMShare"), "Type": "1"}, Daftar: []DaftarJSON{{Nama: "CNPLayerList", Isi: layer}}}
	return BarisOS{CaseID: KunciInstans(id), NoClaim: h.Ambil(CD + "NoClaim"), MasterID: h.Ambil(CD + "IDMaster"),
		NoPolis: h.Ambil(CD + "PolicyData.PolicyNo"), StsReject: StsOSDibayar, DataJSON: JSONHalamanPega(p)}, nil
}

// TanggalStempel - `@getCurrentDateStamp()` / TANGGAL OS (yyyyMMdd hari ini, Jakarta).
func TanggalStempel(t time.Time) string { return t.In(Jakarta).Format("20060102") }
