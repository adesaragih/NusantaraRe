package models

// Untuk apa berkas ini: SELISIH XOL - port `Activity/CalculateDifferenceEDM_act` (kelas Int-treaty_in_edm, ruleset
// 01-01-83; 12 langkah aktif, 1 `//`: 1.2.2). Pemanggil: `EDMChooseBusiness_Act` 6. Hanya halaman.

import (
	"strings"

	"github.com/cockroachdb/apd/v3"
)

// xeDueToYou - VERBATIM CalculateDifferenceEDM_act 1.1 / 1.2.1 (`@If(local.duetovalue>0,"DUE TO US","DUE TO YOU")`).
const xeDueToYou = "DUE TO YOU"

// xeMedanSelisih - medan uang layer yang dihitung 1.2.1 / 1.2.3 / 1.2.5, urut XML.
var xeMedanSelisih = []string{"GrossPremi", "Deduction", "NetPremi", "DueToValue", "BrokerageFeeSebenarnya",
	"PPHValue", "PPNValue", "NetPremiAfterPPH", "NetPremiAfterPPN", "NetPremiAfterTax"}

// xeKunciLayer - medan kunci yang disalin dari baris baru (1.1 / 1.2.1 / 1.2.3 / 1.2.5).
var xeKunciLayer = []string{"LayerType", "Layer", "LayerPartType", "LayerPart", "Currency", "IDCurrency"}

// CalculateDifferenceEDM = `Activity/CalculateDifferenceEDM_act`: `PolicyTreatyIn.TreatyXOLDifferenceList` dari
// `TreatyXOLList` (baru) lawan `OldData.TreatyXOLList` (lama), berpasangan menurut SUBSKRIP induk (c) dan layer (l).
//
//	1.1    induk (c): Layer*, Currency, IDCurrency = baris baru; DueTo = @If(local.duetovalue>0, ...) - lokal BERISI
//	       nilai layer terakhir induk SEBELUMNYA (0 pada induk pertama -> "DUE TO YOU")
//	1.2.1  layer (l), tanpa syarat: X = @if(.X < Old.X, 0, .X - Old.X) untuk sepuluh medan; local.deduc = Deduction
//	       selisih; local.duetovalue = DueToValue selisih; DueTo = duetovalue > 0 ? "DUE TO US" : "DUE TO YOU"
//	1.2.2  `//` (salinan 1.2.4) - tidak diport
//	1.2.3  PRE `pyWorkPage.TreatyIn.IsProRate`: X = @if(.X < Old.X, 0, (.X - Old.X) * @Math.divide(ProRatePercent,100,8))
//	1.2.4  PRE FlagPPH=="true" DAN EDMType=="3": BrokerageFeeSebenarnya = @if(TypeTax="Inclusive",
//	       @divide(local.deduc,@divide(102.2,100,8),8), local.deduc); PPH = BFS*@divide(2,100,8); PPN = BFS*
//	       @divide(2.2,100,8); AfterPPH = NetPremi+PPH; AfterPPN = NetPremi+PPN; AfterTax = NetPremi+PPN+PPH
//	1.2.5  PRE EDMType=="4" || EDMType=="2" ("EDM BATAL"): X = .X - Old.X mentah (boleh negatif), menimpa
//	2      setiap induk: jumlah sepuluh medan layernya (`JumlahIndukSelisihXOL`)
//
// ⚠️ Ditiru apa adanya: daftar selisih TIDAK dihapus dulu - baris ditulis per subskrip, baris lama di luar panjang
// daftar baru tertinggal (dan ikut dijumlah langkah 2). Baris lama yang tidak ada = 0 (properti tak bernilai).
// `[dugaan]` Pembandingan `<` atas properti Decimal kosong = 0 (semantik aritmetika `AngkaTeks`).
// Pembandingan DUA nilai uang (`.X < Old.X`) ADA di rule ini (beda dari NB, yang tidak punya - guard
// pembandingan_uang_test.go); ditiru EKSAK seperti BigDecimal Pega, ditulis `.X - Old.X < 0` (setara, tanpa
// toleransi).
// `[dugaan]` `IsProRate` benar bila teksnya "true" (tanpa beda huruf; properti TrueFalse, diisi `= true` / `= True`
// di korpus Treaty In Adjustment).
func CalculateDifferenceEDM(h *Halaman) error {
	k := &kalkulator{}
	edm := h.Ambil(pt + "EDMType")
	pajak := h.Ambil(pt+"FlagPPH") == "true" && edm == "3" // 1.2.4 (dua baris PRE berurutan)
	batal := edm == "4" || edm == "2"                      // 1.2.5
	prorata := strings.EqualFold(strings.TrimSpace(h.Ambil(jMaster+"IsProRate")), "true")
	tt := h.Ambil(pt + "TypeTax")
	duetovalue := apd.New(0, 0) // local.duetovalue (Decimal)
	dueTo := func() string {
		if duetovalue.Sign() > 0 {
			return teksDueToUs
		}
		return xeDueToYou
	}
	faktor := func() *apd.Decimal { return k.BagiBulat(angkaMaster(k, h, "ProRatePercent"), seratus, 8) }
	for c, x := range h.AmbilDaftar(DaftarXOL) { // 1
		ci := c + 1
		d := xePastikanBaris(h, DaftarSelisihXOL, ci) // 1.1
		for _, m := range xeKunciLayer {
			d[m] = x[m]
		}
		d["DueTo"] = dueTo()
		lama := h.AmbilDaftar(JalurAnak(od+"TreatyXOLList", ci, AnakLayerXOL))
		jalurSelisih := JalurAnak(DaftarSelisihXOL, ci, AnakLayerXOL)
		for l, v := range h.AmbilDaftar(JalurAnak(DaftarXOL, ci, AnakLayerXOL)) { // 1.2
			o := barisKe(lama, l)
			r := xePastikanBaris(h, jalurSelisih, l+1)
			isi := func(rumus func(baru, lawas *apd.Decimal) *apd.Decimal) {
				for _, m := range xeKunciLayer {
					r[m] = v[m]
				}
				for _, m := range xeMedanSelisih {
					r[m] = formatAngka(rumus(angkaBaris(k, v, m), angkaBaris(k, o, m)))
					if m == "DueToValue" { // SET local.duetovalue = (rumus yang sama); DueTo = @If(...)
						duetovalue = angkaBaris(k, r, "DueToValue")
						r["DueTo"] = dueTo()
					}
				}
			}
			// 1.2.1 @if(.X<Old.X, 0, .X-Old.X)
			isi(func(baru, lawas *apd.Decimal) *apd.Decimal {
				if s := k.Kurang(baru, lawas); s.Sign() >= 0 {
					return s
				}
				return apd.New(0, 0)
			})
			if prorata { // 1.2.3 @if(.X<Old.X, 0, (.X-Old.X)*@Math.divide(ProRatePercent,100,8))
				isi(func(baru, lawas *apd.Decimal) *apd.Decimal {
					if s := k.Kurang(baru, lawas); s.Sign() >= 0 {
						return k.Kali(s, faktor())
					}
					return apd.New(0, 0)
				})
			}
			if pajak { // 1.2.4 - local.deduc = Deduction selisih (1.2.1 / 1.2.3)
				deduc := angkaBaris(k, r, "Deduction")
				r["BrokerageFeeSebenarnya"] = formatAngka(BrokerageSebenarnya(k, tt, deduc))
				lPPH, lPPN := pajakXOL(k, angkaBaris(k, r, "BrokerageFeeSebenarnya"))
				r["PPHValue"] = formatAngka(lPPH)
				r["PPNValue"] = formatAngka(lPPN)
				net := angkaBaris(k, r, "NetPremi")
				r["NetPremiAfterPPH"] = formatAngka(k.Tambah(net, angkaBaris(k, r, "PPHValue")))
				r["NetPremiAfterPPN"] = formatAngka(k.Tambah(net, angkaBaris(k, r, "PPNValue")))
				r["NetPremiAfterTax"] = formatAngka(k.Tambah(k.Tambah(net, angkaBaris(k, r, "PPNValue")), angkaBaris(k, r, "PPHValue")))
			}
			if batal { // 1.2.5
				isi(func(baru, lawas *apd.Decimal) *apd.Decimal { return k.Kurang(baru, lawas) })
			}
		}
	}
	if k.err != nil {
		return k.err
	}
	return JumlahIndukSelisihXOL(h) // 2
}

// xePastikanBaris - baris ke-n (berbasis 1) daftar `jalur`, dibuat kosong bila belum ada (penulisan per subskrip
// `List(n).Prop = ...` membuat anggota sampai n).
func xePastikanBaris(h *Halaman, jalur string, n int) Baris {
	d := h.AmbilDaftar(jalur)
	for len(d) < n {
		d = append(d, Baris{})
	}
	if d[n-1] == nil {
		d[n-1] = Baris{}
	}
	h.SetelDaftar(jalur, d)
	return d[n-1]
}
