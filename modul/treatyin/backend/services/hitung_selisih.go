package services

// Rumus tombol `Update Value` tab VALUE DIFFERENCE (Non-Prop, panel New
// layar Adjustment) — `Activity/TreatyEDMCalculateDifference.xml` dan
// kesembilan Activity yang dipanggilnya. Kesepuluhnya IDENTIK di korpus
// `Treaty In` dan `Treaty In Adjustment` (diadu 7 Oktober 2026).
//
//	 1  Property-Remove TreatyIn.ValueDifference
//	 2  TreatyEDMDifferencePremium      selisih EGNPI + total + proporsi
//	 3  TreatyEDMDifferenceLimits       selisih layer, ringkasan, total, ROL
//	 4  TreatyEDMDifferenceShare        selisih skalar + baris Share
//	 5  TreatyEDMDifferenceDeduction    Brokerage fee + CalculateDeduction
//	 6  TreatyInDifferenceShareSumary   ringkasan (= TreatyInSummaryLimitShare)
//	 7  TreatyInDifferenceShareTotal    total per mata uang
//	 8  TreatyInSummaryLimitShareActual ringkasan ActualValue.Share
//	 9-10 TreatyInSetValueDifferenceInstallment
//	11  TreatyEDMProRateCalculation     bila TreatyIn.IsProRate == true
//
// Semua selisih = New − Old (`TreatyIn.X − TreatyIn.OLDDATA.X`), per INDEKS
// larik. Nilai kosong / halaman yang tidak ada = 0, seperti aritmetika Pega
// atas properti kosong.
//
// ⚠️ DISALIN APA ADANYA walau janggal — rumus, bukan alamat:
//   - Limits [1.3.1] menulis `EgnpiTotalList(<CURRENT>).Currency` dari baris
//     Premium Earned; `PremiumEarnedList` selisih hanya membawa `Value`.
//   - Limits [3.5] `TotalLimitsROL` = `ActualValue.TotalLimitsROL − OLDDATA`
//     (halaman langkahnya `TreatyIn.ActualValue`).
//   - Share [3.3] menulis `RnmLimitListDisplay(2)`/`RnmGrossPremiDisplay(2)`
//     walau sumbernya bermata uang satu — baris kedua bermata uang kosong.
//   - Share [4] menimpa `ActualValue.LimitShareSummaryList` dengan
//     selisihnya sendiri terhadap OLDDATA; langkah 8 lalu menyusunnya ulang.
//   - Deduction: DeductionList dibuang, "Brokerage fee" ditambahkan dengan
//     PERSEN Old tanpa nilai, `CalculateDeduction` dipanggil TANPA parameter
//     (nilai tidak dihitung; total dan Net disusun dari Gross selisih), lalu
//     langkah 6 membuang baris bernilai < 1 — baris Brokerage itu sendiri.
//     Langkah 5 menghitung ulang `TreatyIn.FacultativeShareList` (akar).
//   - ProRate: langkah berblok `//` (Deduction, DeductionTotalList, Deduct
//     OR/RI, TotalShareDeductionNP) dikomentari di Pega — TIDAK dijalankan.
//
// ⛔ MURNI: nol baca, nol tulis basis data. Bentuk halaman = pohon generik
// (`HalamanPohon`): Value Difference memuat puluhan larik bersarang, dan
// layar Adjustment menyimpannya apa adanya.

import (
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

// AksiSelisihEDM - tombol `Update Value` tab Value Difference.
const AksiSelisihEDM = "edm"

// HalamanPohon - satu halaman Pega sebagai skalar + larik baris generik.
type HalamanPohon struct {
	Medan map[string]string           `json:"medan"`
	Larik map[string][]map[string]any `json:"larik"`
}

// MasukanSelisih - halaman yang Activity Value Difference baca.
type MasukanSelisih struct {
	Aksi string `json:"aksi"`
	// Akar - `TreatyIn` (panel New).
	Akar HalamanPohon `json:"akar"`
	// Lama - `TreatyIn.OLDDATA` (panel Old).
	Lama HalamanPohon `json:"lama"`
	// Actual - `TreatyIn.ActualValue` (cabang Adjust Premium; kosong di luar itu).
	Actual HalamanPohon `json:"actual"`
}

// HasilSelisih - halaman yang Activity tulis.
type HasilSelisih struct {
	// Selisih - `TreatyIn.ValueDifference`, disusun ulang utuh (langkah 1).
	Selisih HalamanPohon `json:"selisih"`
	// SebelumProrata - `TreatyIn.ValueBeforeProrate` (ProRate langkah 1);
	// `null` bila `IsProRate` bukan `true`.
	SebelumProrata *HalamanPohon `json:"sebelumProrata"`
	// Actual - `TreatyIn.ActualValue` sesudah Share [4] dan langkah 8.
	Actual HalamanPohon `json:"actual"`
	// FacultativeShareList - `TreatyIn.FacultativeShareList` sesudah
	// Deduction [5] (`CalculateDeduction` atas akar).
	FacultativeShareList []map[string]any `json:"FacultativeShareList"`
	Pesan                []string         `json:"pesan"`
}

// HitungSelisih - bentuk ber-pelaku untuk handler.
func (l *Layanan) HitungSelisih(p inti.Pelaku, m MasukanSelisih) (HasilSelisih, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilSelisih{}, err
	}
	return HitungSelisih(m)
}

// HitungSelisih menjalankan satu aksi Value Difference.
func HitungSelisih(m MasukanSelisih) (HasilSelisih, error) {
	switch m.Aksi {
	case AksiSelisihEDM:
		return selisihEDM(m), nil
	}
	return HasilSelisih{}, fmt.Errorf("%w: aksi selisih %q", ErrMasukanTidakSah, m.Aksi)
}

func selisihEDM(m MasukanSelisih) HasilSelisih {
	akar, lama, act := salinHalaman(m.Akar), salinHalaman(m.Lama), salinHalaman(m.Actual)
	// [1]
	vd := halamanKosong()
	selisihPremi(akar, lama, &vd)
	selisihLimits(akar, lama, act, &vd)
	selisihShare(akar, lama, &act, &vd)
	pesan := selisihDeduksi(&akar, lama, &vd)
	// [6] — struktur TreatyInSummaryLimitShare, atas ValueDifference.Share.
	vd.Larik["LimitShareSummaryList"] = ringkasPohon(vd.Larik["Share"])
	totalSelisihShare(&vd)
	// [8]
	act.Larik["LimitShareSummaryList"] = ringkasPohon(act.Larik["Share"])
	selisihAngsuran(akar, &vd)
	var sebelum *HalamanPohon
	// [11] `TreatyIn.IsProRate==true`.
	if akar.Medan["IsProRate"] == "true" {
		s := salinHalaman(vd)
		sebelum = &s
		prorata(akar.Medan["ProRatePercent"], &vd)
	}
	return HasilSelisih{
		Selisih:              vd,
		SebelumProrata:       sebelum,
		Actual:               act,
		FacultativeShareList: akar.Larik["FacultativeShareList"],
		Pesan:                pesanUnik(pesan),
	}
}

// --- [2] TreatyEDMDifferencePremium -----------------------------------------

func selisihPremi(akar, lama HalamanPohon, vd *HalamanPohon) {
	out := []map[string]any{}
	for i, e := range akar.Larik["EGNPI"] {
		o := baris(lama.Larik["EGNPI"], i)
		out = append(out, map[string]any{
			"Currency":       teksSimpul(e, "Currency"),
			"TreatyGroup":    teksSimpul(e, "TreatyGroup"),
			"Amount":         beda(e, o, "Amount"),
			"AmountIDR":      beda(e, o, "AmountIDR"),
			"Proportion":     beda(e, o, "Proportion"),
			"AddendumStatus": "1",
		})
	}
	vd.Larik["EGNPI"] = out
	if len(out) == 0 {
		return
	}
	// [2]
	total := apd.New(0, 0)
	for _, r := range out {
		total = tambah(angka(teksSimpul(r, "AmountIDR")), total)
	}
	vd.Medan["TotalEgnpiAmount"] = teks(total)
	// [3] — `TotalEgnpiAmount=="0"` → KELUAR (6) sebelum baris pertama.
	if vd.Medan["TotalEgnpiAmount"] == "0" {
		return
	}
	np := []map[string]any{}
	proporsi := apd.New(0, 0)
	for _, r := range out {
		cur, nilai := teksSimpul(r, "Currency"), angka(teksSimpul(r, "Amount"))
		p := kali(bagiBulatDes(angka(teksSimpul(r, "AmountIDR")), total, 20), seratus)
		r["Proportion"] = teks(p)
		proporsi = tambah(p, proporsi)
		vd.Medan["TotalEgnpiProportion"] = teks(proporsi)
		ada := false
		for _, t := range np {
			if teksSimpul(t, "Currency") == cur {
				t["Value"] = teks(tambah(nilai, angka(teksSimpul(t, "Value"))))
				ada = true
			}
		}
		if !ada {
			np = append(np, map[string]any{"ID": "", "Currency": cur, "Value": teks(nilai), "AltValue": ""})
		}
	}
	vd.Larik["TotalEgnpiAmountNP"] = np
}

// --- [3] TreatyEDMDifferenceLimits ------------------------------------------

func selisihLimits(akar, lama, act HalamanPohon, vd *HalamanPohon) {
	out := []map[string]any{}
	for i, l := range akar.Larik["Limits"] {
		o := baris(lama.Larik["Limits"], i)
		v := map[string]any{"TreatyGroupList": salinPohon(larikSimpul(l, "TreatyGroupList"))}
		for _, k := range []string{"LayerType", "Layer", "LayerPartType", "LayerPart", "Cover", "ReinstatementPct", "Currency", "Currency2"} {
			v[k] = teksSimpul(l, k)
		}
		for _, k := range []string{"Limit", "Limit2", "Deductible", "Deductible2", "AdjRate", "MDPPct", "ROLPct"} {
			v[k] = beda(l, o, k)
		}
		egnpi := selisihNilai(larikSimpul(l, "EgnpiTotalList"), larikSimpul(o, "EgnpiTotalList"), true)
		// [1.3.1] ⚠️ Currency Premium Earned ke `EgnpiTotalList(<CURRENT>)`.
		pe := []map[string]any{}
		peLama := larikSimpul(o, "PremiumEarnedList")
		for j, x := range larikSimpul(l, "PremiumEarnedList") {
			var e map[string]any
			egnpi, e = pasang(egnpi, j)
			e["Currency"] = teksSimpul(x, "Currency")
			pe = append(pe, map[string]any{"Value": beda(x, baris(peLama, j), "Value")})
		}
		v["EgnpiTotalList"] = egnpi
		v["PremiumEarnedList"] = pe
		v["MDPList"] = selisihNilai(larikSimpul(l, "MDPList"), larikSimpul(o, "MDPList"), true)
		out = append(out, v)
	}
	vd.Larik["Limits"] = out
	// [2]
	ringkas := []map[string]any{}
	for i, r := range akar.Larik["LimitSummaryList"] {
		o := baris(lama.Larik["LimitSummaryList"], i)
		v := map[string]any{"Note": teksSimpul(r, "Note")}
		for _, k := range []string{"Limit", "Limit2", "Deductible", "Deductible2", "MDP", "MDP2", "AggregateLimit", "AggregateLimit2"} {
			v[k] = beda(r, o, k)
		}
		ringkas = append(ringkas, v)
	}
	vd.Larik["LimitSummaryList"] = ringkas
	// [3.1]–[3.4]
	for _, n := range []string{"TotalLimitIOONP", "TotalLimitDeductblNP", "TotalLimitPremiEarnNP", "TotalLimitMDPNP"} {
		vd.Larik[n] = selisihNilai(akar.Larik[n], lama.Larik[n], true)
	}
	// [3.5] ⚠️ halaman langkah `TreatyIn.ActualValue`.
	vd.Medan["TotalLimitsROL"] = teks(kurang(angka(act.Medan["TotalLimitsROL"]), angka(lama.Medan["TotalLimitsROL"])))
}

// --- [4] TreatyEDMDifferenceShare -------------------------------------------

// larikShareSelisih - larik bernilai per mata uang baris Share yang diselisihkan
// ([3.4]–[3.15]; DeductionList tersendiri).
var larikShareSelisih = []string{
	"RnmLimitList", "GrossPremiumList", "NetPremiumList",
	"RNMSpreadedListXOL", "RNMSpreadedListRIXOL", "RNMSpreadedListGrossXOL", "RNMSpreadedListGrossRIXOL",
	"RNMSpreadedListDeductXOL", "RNMSpreadedListDeductRIXOL", "RNMSpreadedListNetXOL", "RNMSpreadedListNetRIXOL",
}

func selisihShare(akar, lama HalamanPohon, act *HalamanPohon, vd *HalamanPohon) {
	// [1] EDMState 3 → dari ActualValue; [2] selain itu dari akar.
	sumber := akar.Medan
	if akar.Medan["EDMState"] == "3" {
		sumber = act.Medan
	}
	vd.Medan["RNMShare"] = teks(kurang(angka(sumber["RNMShare"]), angka(lama.Medan["RNMShare"])))
	vd.Medan["BrokeragePercent"] = teks(kurang(angka(sumber["BrokeragePercent"]), angka(lama.Medan["BrokeragePercent"])))
	out := []map[string]any{}
	for i, s := range akar.Larik["Share"] {
		o := baris(lama.Larik["Share"], i)
		v := map[string]any{
			"TreatyGroupList":      salinPohon(larikSimpul(s, "TreatyGroupList")),
			"SpreadingTypeXOL":     teksSimpul(o, "SpreadingTypeXOL"),
			"SpreadingTotalPctXOL": teksSimpul(o, "SpreadingTotalPctXOL"),
			"AddendumStatus":       "1",
		}
		for _, k := range []string{"LayerType", "Layer", "LayerPartType", "LayerPart", "Cover"} {
			v[k] = teksSimpul(s, k)
		}
		// [3.2]
		spread := []map[string]any{}
		for _, x := range larikSimpul(s, "SpreadingListXOL") {
			spread = append(spread, map[string]any{"Pct": teksSimpul(x, "Pct"), "ReinsTypeName": teksSimpul(x, "ReinsTypeName")})
		}
		v["SpreadingListXOL"] = spread
		// [3.3] kolom TAMPIL (1) dan (2). Larik tampil tidak didaratkan; ia
		// sama dengan sumbernya (`TreatyInNonAddItem` 14.3).
		for _, pas := range [][2]string{{"RnmLimitListDisplay", "RnmLimitList"}, {"RnmGrossPremiDisplay", "GrossPremiumList"}} {
			kini, dulu := larikTampil(s, pas[0], pas[1]), larikTampil(o, pas[0], pas[1])
			t := make([]map[string]any, 2)
			for j := range t {
				a := baris(kini, j)
				t[j] = map[string]any{"Currency": teksSimpul(a, "Currency"), "Value": beda(a, baris(dulu, j), "Value")}
			}
			v[pas[0]] = t
		}
		for _, n := range larikShareSelisih {
			v[n] = selisihNilai(larikSimpul(s, n), larikSimpul(o, n), true)
		}
		// [3.6]
		ded := []map[string]any{}
		dedLama := larikSimpul(o, "DeductionList")
		for j, d := range larikSimpul(s, "DeductionList") {
			x := baris(dedLama, j)
			ded = append(ded, map[string]any{
				"Currency":     teksSimpul(d, "Currency"),
				"Comment":      teksSimpul(d, "Comment"),
				"Deduction":    beda(d, x, "Deduction"),
				"DeductionPct": beda(d, x, "DeductionPct"),
			})
		}
		v["DeductionList"] = ded
		out = append(out, v)
	}
	vd.Larik["Share"] = out
	// [4] ⚠️ menimpa ActualValue dengan selisihnya terhadap OLDDATA.
	for i, r := range act.Larik["LimitShareSummaryList"] {
		o := baris(lama.Larik["LimitShareSummaryList"], i)
		for _, k := range []string{"Limit", "Limit2", "Deductible", "Deductible2", "MDP", "MDP2", "AggregateLimit", "AggregateLimit2"} {
			r[k] = beda(r, o, k)
		}
	}
}

// --- [5] TreatyEDMDifferenceDeduction ---------------------------------------

func selisihDeduksi(akar *HalamanPohon, lama HalamanPohon, vd *HalamanPohon) []string {
	pesan := []string{}
	share := vd.Larik["Share"]
	for _, b := range share {
		// [1] Property-Remove .DeductionList, [2] Brokerage fee (persen Old).
		b["DeductionList"] = []map[string]any{{
			"Comment": "Brokerage fee", "DeductionPct": lama.Medan["BrokeragePercent"], "DeductionPctCalculate": "true",
		}}
	}
	// [3] CalculateDeduction TANPA parameter.
	for _, b := range share {
		pesan = append(pesan, deduksiPohon(b)...)
	}
	// [4] — `ValueDifference.FacultativeShareList` tidak pernah diisi rantai ini.
	for _, b := range vd.Larik["FacultativeShareList"] {
		b["DeductionList"] = append(larikSimpul(b, "DeductionList"), map[string]any{
			"Comment": "Facultative Brokerage fee", "DeductionPct": akar.Medan["FacultativeShareBrokerage"], "DeductionPctCalculate": "true",
		})
	}
	// [5] ⚠️ atas `TreatyIn.FacultativeShareList` (akar).
	for _, b := range akar.Larik["FacultativeShareList"] {
		pesan = append(pesan, deduksiPohon(b)...)
	}
	// [6] buang SATU baris bernilai < 1 — yang TERAKHIR ditemukan
	// (`subscriptdelete != 0` selalu benar: subscript mulai 1).
	satu := apd.New(1, 0)
	for _, b := range share {
		dl := larikSimpul(b, "DeductionList")
		hapus := -1
		for j, d := range dl {
			if angka(teksSimpul(d, "Deduction")).Cmp(satu) < 0 {
				hapus = j
			}
		}
		if hapus >= 0 {
			b["DeductionList"] = append(append([]map[string]any{}, dl[:hapus]...), dl[hapus+1:]...)
		}
	}
	return pesan
}

// deduksiPohon - `CalculateDeduction` tanpa parameter atas satu baris pohon.
// Kunci baris deduksi yang tidak dikenal rute ditimpakan, bukan dibuang.
func deduksiPohon(b map[string]any) []string {
	dl := larikSimpul(b, "DeductionList")
	h := HitungDeduksi(MasukanDeduksi{
		Sts:              "",
		Indeks:           -1,
		DeductionList:    keJenis[[]BarisDeduksi](dl),
		GrossPremiumList: keJenis[[]NilaiMataUang](larikSimpul(b, "GrossPremiumList")),
		BenderaPct:       true,
	})
	baru := keJenis[[]map[string]any](h.DeductionList)
	for i := range baru {
		if i < len(dl) {
			gabung := map[string]any{}
			for k, v := range dl[i] {
				gabung[k] = v
			}
			for k, v := range baru[i] {
				if s, ok := v.(string); ok && s == "" {
					if _, ada := gabung[k]; ada {
						continue
					}
				}
				gabung[k] = v
			}
			baru[i] = gabung
		}
	}
	b["DeductionList"] = baru
	b["DeductionTotalList"] = keJenis[[]map[string]any](h.DeductionTotalList)
	b["NetPremiumList"] = keJenis[[]map[string]any](h.NetPremiumList)
	return h.Pesan
}

// --- [6]/[8] ringkasan — struktur TreatyInSummaryLimitShare --------------------

func ringkasPohon(rows []map[string]any) []map[string]any {
	return keJenis[[]map[string]any](SummaryLimitShare(keJenis[[]models.BarisShareNP](rows)))
}

// --- [7] TreatyInDifferenceShareTotal ---------------------------------------

func totalSelisihShare(vd *HalamanPohon) {
	// [1]
	for _, n := range []string{"TotalShareRnmNP", "TotalShareGrossNP", "TotalShareNetNP", "TotalShareDeductionNP"} {
		vd.Larik[n] = []map[string]any{}
	}
	for _, n := range []string{"TotalSpreadedRnmProp", "TotalSpreadedRnmRIProp", "TotalSpreadedNetPremi", "TotalSpreadedNetPremiRI"} {
		if vd.Larik[n] == nil {
			vd.Larik[n] = []map[string]any{}
		}
	}
	// [2] — baris tambahan HANYA bila mata uangnya tidak kosong.
	for _, b := range vd.Larik["Share"] {
		for _, p := range [][2]string{
			{"GrossPremiumList", "TotalShareGrossNP"}, {"NetPremiumList", "TotalShareNetNP"},
			{"DeductionTotalList", "TotalShareDeductionNP"}, {"RnmLimitList", "TotalShareRnmNP"},
		} {
			for _, x := range larikSimpul(b, p[0]) {
				vd.Larik[p[1]] = jumlahPohon(vd.Larik[p[1]], x, true)
			}
		}
	}
	// [3] — tanpa syarat mata uang.
	for _, b := range vd.Larik["Share"] {
		for _, p := range [][2]string{
			{"RNMSpreadedListXOL", "TotalSpreadedRnmProp"}, {"RNMSpreadedListRIXOL", "TotalSpreadedRnmRIProp"},
			{"RNMSpreadedListNetXOL", "TotalSpreadedNetPremi"}, {"RNMSpreadedListNetRIXOL", "TotalSpreadedNetPremiRI"},
		} {
			for _, x := range larikSimpul(b, p[0]) {
				vd.Larik[p[1]] = jumlahPohon(vd.Larik[p[1]], x, false)
			}
		}
	}
}

// jumlahPohon - `.Value += local.Value` pada baris bermata uang sama, atau
// baris baru (bila `wajibMataUang`, hanya bila mata uangnya tidak kosong).
func jumlahPohon(total []map[string]any, x map[string]any, wajibMataUang bool) []map[string]any {
	cur, nilai := teksSimpul(x, "Currency"), angka(teksSimpul(x, "Value"))
	ada := false
	for _, t := range total {
		if teksSimpul(t, "Currency") == cur {
			t["Value"] = teks(tambah(angka(teksSimpul(t, "Value")), nilai))
			ada = true
		}
	}
	if !ada && (!wajibMataUang || cur != "") {
		total = append(total, map[string]any{"Currency": cur, "Value": teks(nilai)})
	}
	return total
}

// --- [9]–[10] TreatyInSetValueDifferenceInstallment -------------------------

func selisihAngsuran(akar HalamanPohon, vd *HalamanPohon) {
	// [1]–[2] satu halaman per mata uang Total Net Premium selisih.
	inst := []map[string]any{}
	for _, n := range vd.Larik["TotalShareNetNP"] {
		inst = append(inst, map[string]any{"Currency": teksSimpul(n, "Currency")})
	}
	// [3] jadwal disalin dari Installment akar bermata uang sama.
	for _, v := range inst {
		cur := teksSimpul(v, "Currency")
		rows := []map[string]any{}
		for _, a := range akar.Larik["Installment"] {
			if teksSimpul(a, "Currency") != cur {
				continue
			}
			for j, r := range larikSimpul(a, "InstallmentList") {
				var x map[string]any
				rows, x = pasang(rows, j)
				for _, k := range []string{"Currency", "CurrencyID", "Installment", "InstallmentPct", "DueDate", "PaymentDate", "WPC"} {
					x[k] = teksSimpul(r, k)
				}
			}
		}
		v["InstallmentList"] = rows
	}
	// [4] Amount = @divide(Net × %, 100, 4).
	for _, n := range vd.Larik["TotalShareNetNP"] {
		cur, nilai := teksSimpul(n, "Currency"), angka(teksSimpul(n, "Value"))
		for _, v := range inst {
			if teksSimpul(v, "Currency") != cur {
				continue
			}
			for _, r := range larikSimpul(v, "InstallmentList") {
				r["Amount"] = teks(bagiBulat(kali(nilai, angka(teksSimpul(r, "InstallmentPct"))), 100, 4))
			}
		}
	}
	// [5] SetTotalInstallment (tanpa status: hanya total), [6] total.
	total := []map[string]any{}
	for _, v := range inst {
		a := totalHalamanAngsuran(keJenis[Angsuran](v), "", nil)
		v["AmountTotal"], v["PctTotal"] = a.AmountTotal, a.PctTotal
		total = append(total, map[string]any{"Currency": teksSimpul(v, "Currency"), "Value": a.AmountTotal})
	}
	vd.Larik["Installment"] = inst
	vd.Larik["TotalInstallmentNP"] = total
}

// --- [11] TreatyEDMProRateCalculation ---------------------------------------

func prorata(persen string, vd *HalamanPohon) {
	f := bagiBulat(angka(persen), 100, 12)
	skala := func(r map[string]any, k string) { r[k] = teks(kali(angka(teksSimpul(r, k)), f)) }
	for _, b := range vd.Larik["Share"] {
		// [2.1] Deduction, [2.8]–[2.10] berblok `//` — dikomentari.
		for _, n := range []string{
			"RnmLimitList", "RNMSpreadedListXOL", "RNMSpreadedListRIXOL", "GrossPremiumList",
			"RNMSpreadedListGrossXOL", "RNMSpreadedListGrossRIXOL", "NetPremiumList",
			"RNMSpreadedListNetXOL", "RNMSpreadedListNetRIXOL",
		} {
			for _, x := range larikSimpul(b, n) {
				skala(x, "Value")
			}
		}
		// [2.14] tampil (1) dan (2).
		for _, n := range []string{"RnmLimitListDisplay", "RnmGrossPremiDisplay"} {
			xs := larikSimpul(b, n)
			for j := 0; j < 2; j++ {
				var x map[string]any
				xs, x = pasang(xs, j)
				skala(x, "Value")
			}
			b[n] = xs
		}
	}
	// [3]
	for _, r := range vd.Larik["LimitShareSummaryList"] {
		for _, k := range []string{"Limit", "Limit2", "MDP", "MDP2", "Deductible", "Deductible2", "NetPremi", "NetPremi2"} {
			skala(r, k)
		}
	}
	// [4]–[11] ([8] TotalShareDeductionNP berblok `//`).
	for _, n := range []string{
		"TotalShareRnmNP", "TotalSpreadedRnmProp", "TotalSpreadedRnmRIProp", "TotalShareGrossNP",
		"TotalShareNetNP", "TotalSpreadedNetPremi", "TotalSpreadedNetPremiRI", "TotalInstallmentNP",
	} {
		for _, x := range vd.Larik[n] {
			skala(x, "Value")
		}
	}
	// [12]
	for _, a := range vd.Larik["Installment"] {
		skala(a, "AmountTotal")
		for _, r := range larikSimpul(a, "InstallmentList") {
			skala(r, "Amount")
		}
	}
}

// --- pembantu pohon -----------------------------------------------------------

func halamanKosong() HalamanPohon {
	return HalamanPohon{Medan: map[string]string{}, Larik: map[string][]map[string]any{}}
}

// salinHalaman - salinan DALAM (larik bersarang ikut disalin).
func salinHalaman(h HalamanPohon) HalamanPohon {
	out := halamanKosong()
	for k, v := range h.Medan {
		out.Medan[k] = v
	}
	for k, v := range h.Larik {
		out.Larik[k] = keJenis[[]map[string]any](v)
	}
	return out
}

// baris - baris ke-i, atau halaman kosong (properti halaman yang tidak ada = kosong).
func baris(xs []map[string]any, i int) map[string]any {
	if i >= 0 && i < len(xs) {
		return xs[i]
	}
	return map[string]any{}
}

// pasang - `X(n)` Pega: halaman ke-n DIBUAT bila belum ada.
func pasang(xs []map[string]any, i int) ([]map[string]any, map[string]any) {
	for len(xs) <= i {
		xs = append(xs, map[string]any{})
	}
	return xs, xs[i]
}

// beda - `.k − OLDDATA.k` (kosong = 0).
func beda(a, b map[string]any, k string) string {
	return teks(kurang(angka(teksSimpul(a, k)), angka(teksSimpul(b, k))))
}

// selisihNilai - larik {Currency, Value} per indeks: Value = kini − dulu.
func selisihNilai(kini, dulu []map[string]any, denganMataUang bool) []map[string]any {
	out := []map[string]any{}
	for j, x := range kini {
		v := map[string]any{"Value": beda(x, baris(dulu, j), "Value")}
		if denganMataUang {
			v["Currency"] = teksSimpul(x, "Currency")
		}
		out = append(out, v)
	}
	return out
}

// larikTampil - larik tampil, atau sumbernya bila tampilnya tidak ada.
func larikTampil(s map[string]any, tampil, sumber string) []map[string]any {
	if xs := larikSimpul(s, tampil); len(xs) > 0 {
		return xs
	}
	return larikSimpul(s, sumber)
}

// keJenis - konversi lewat JSON antara pohon generik dan tipe rute.
func keJenis[T any](v any) T {
	var t T
	b, err := json.Marshal(v)
	if err == nil {
		_ = json.Unmarshal(b, &t)
	}
	return t
}

// pesanUnik - pesan Activity tanpa ulangan, urut kemunculan.
func pesanUnik(xs []string) []string {
	out := []string{}
	lihat := map[string]bool{}
	for _, x := range xs {
		if !lihat[x] {
			lihat[x] = true
			out = append(out, x)
		}
	}
	return out
}
