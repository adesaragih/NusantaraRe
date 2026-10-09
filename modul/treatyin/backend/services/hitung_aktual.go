package services

// Rumus cabang ADJUST PREMIUM layar Adjustment — Section
// `TreatyInTabsNonProportionalAdjustPremi` (`TreatyInNONProportional` @566980:
// `ProportionType='NonProportional' && EDMState=3`). Halaman yang dihitung =
// `TreatyIn.ActualValue`; selisihnya = `TreatyIn.ValueDifference`
// (Actual − Treaty In). Disalin langkah demi langkah dari korpus
// `Treaty In Adjustment` (7 Oktober 2026):
//
//	nilai      Actual GNPI · Update Total — TreatyInActualUpdateValue:
//	           PremiumIncome, ActualUpdateValueLimits, LimitsActualSetTotal,
//	           SummaryLimitActual, ActualShare, SetBrokerageActual,
//	           ShareActualSetTotal, Property-Remove ValueDifference,
//	           DifferencePremium/Limits/LimitsSumary/LimitsSetTotal/Share/
//	           Deduction/ShareSumary/ShareTotal, SummaryLimitShareActual
//	share      Actual Share · Update Summary — TreatyInActualUpdateValueShare
//	           (+ DifferenceFacShare/SummaryFacShare/FacShareTotal)
//	limits     Actual Limits / Premium Adjustment · Update Total —
//	           TreatyInNPSetTotalActual(limits) = NPSetTotalActualLimits +
//	           NPSetTotalActualShare; lalu SummaryLimitActual dan
//	           SummaryLimitShareActual (langkah tombol berikutnya)
//	ringkas-limit / ringkas-share  kedua Summary itu sendiri
//	rnm-baris  rincian Share `% RNM Share` — TreatyInXOLAddSpreadingDetailActual(idx)
//
// ⚠️ DISALIN APA ADANYA walau janggal (semua tercatat di dokumen modul):
//   - ActualUpdateValueLimits [2] MENIMPA `ActualValue.Limits` dengan salinan
//     `TreatyIn.Limits`, lalu menghitung EGNPI/PE/MDP-nya dari Actual GNPI;
//     MDP% dipaksa 100.
//   - PremiumIncome [3] keluar bila `TreatyIn.TotalEgnpiAmount` (Treaty In,
//     bukan Actual) nol — pesannya `Param.errmsg1`, yang tak pernah dikirim.
//   - ActualShare [4] memakai `TreatyIn.FacultativeShare` untuk limit pertama
//     dan `ActualValue.FacultativeShare` untuk limit kedua/premi; [3] menulis
//     `TreatyIn.RnmShareDeducted` (akar).
//   - SetBrokerageActual [6] membuang `TreatyIn.LimitFacShareSummaryList`
//     (akar) dan menulis ringkasannya ke `ActualValue` per indeks.
//   - NPSetTotalActualLimits [10] menambahkan baris bermata uang (tanpa
//     nilai) ke `TreatyIn.TotalLimitPremiEarnNP/MDPNP` (akar).
//   - NPSetTotalActualShare [3] membuang `TreatyIn.TotalFac*` (akar), padahal
//     yang dijumlahkan [4] adalah `ActualValue.TotalFac*` — yang tidak
//     dibuang, sehingga terakumulasi tiap klik.
//   - XOLAddSpreadingDetailActual [6] mencari total di `ActualValue` (yang
//     baru dibuang [1]) dan menambahkan baris ke total AKAR
//     `TreatyIn.TotalSpreaded*` — tiap klik menambah baris.
//   - Difference*: selisih premi/MDP/Gross/Deduction/Net/spreading bernilai
//     NOL bila Actual < Treaty In (`@if(.Value < estimasi, 0, …)`); selisih
//     limit dan RnmLimit tidak.
//   - DifferenceLimits [3] (total & ROL) berblok `//`; DifferenceLimitsSetTotal
//     tidak mengisi total limit/deductible (parameter `type` tidak dikirim).
//
// ⛔ MURNI: nol baca, nol tulis basis data (`FetchQSfromMasterXOL` dipanggil
// TANPA induk — tidak membaca RD).

import (
	"fmt"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

// Aksi cabang Adjust Premium.
const (
	AksiAktualNilai        = "nilai"
	AksiAktualShare        = "share"
	AksiAktualLimits       = "limits"
	AksiAktualRingkasLimit = "ringkas-limit"
	AksiAktualRingkasShare = "ringkas-share"
	AksiAktualBaris        = "rnm-baris"
)

// MasukanAktual - halaman yang Activity cabang Adjust Premium baca.
type MasukanAktual struct {
	Aksi string `json:"aksi"`
	// Akar - `TreatyIn` (panel New).
	Akar HalamanPohon `json:"akar"`
	// Actual - `TreatyIn.ActualValue`.
	Actual HalamanPohon `json:"actual"`
	// Indeks - `Param.idx` (mulai 0) untuk `rnm-baris`.
	Indeks int `json:"indeks"`
}

// HasilAktual - halaman sesudah aksi.
type HasilAktual struct {
	// Actual - `TreatyIn.ActualValue` utuh.
	Actual HalamanPohon `json:"actual"`
	// Selisih - `TreatyIn.ValueDifference` disusun ulang; `null` bila aksinya
	// tidak menyentuhnya.
	Selisih *HalamanPohon `json:"selisih"`
	// Akar - HANYA kunci akar yang Activity tulis.
	Akar  HalamanPohon `json:"akar"`
	Pesan []string     `json:"pesan"`
}

// HitungAktual - bentuk ber-pelaku untuk handler.
func (l *Layanan) HitungAktual(p inti.Pelaku, m MasukanAktual) (HasilAktual, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilAktual{}, err
	}
	return HitungAktual(m)
}

// konteksAktual - halaman yang dibaca/ditulis satu rantai.
type konteksAktual struct {
	akar, av HalamanPohon
	// ubah - kunci AKAR yang ditulis (dikirim balik).
	ubah  HalamanPohon
	pesan []string
}

// HitungAktual menjalankan satu aksi cabang Adjust Premium.
func HitungAktual(m MasukanAktual) (HasilAktual, error) {
	k := &konteksAktual{akar: salinHalaman(m.Akar), av: salinHalaman(m.Actual), ubah: halamanKosong(), pesan: []string{}}
	var vd *HalamanPohon
	switch m.Aksi {
	case AksiAktualNilai:
		k.premiIncome()
		k.aktualLimits()
		k.totalLimitsAktual(false)
		k.ringkasLimitAktual()
		k.aktualShare()
		k.brokerageAktual()
		k.totalShareAktual(false)
		v := k.selisih(false)
		vd = &v
		k.ringkasShareAktual()
	case AksiAktualShare:
		k.aktualShare()
		k.brokerageAktual()
		k.totalShareAktual(false)
		v := k.selisih(true)
		vd = &v
		k.ringkasShareAktual()
	case AksiAktualLimits:
		k.totalLimitsAktual(true)
		k.totalShareAktual(true)
		k.ringkasLimitAktual()
		k.ringkasShareAktual()
	case AksiAktualRingkasLimit:
		k.ringkasLimitAktual()
	case AksiAktualRingkasShare:
		k.ringkasShareAktual()
	case AksiAktualBaris:
		k.xolDetailAktual(m.Indeks)
	default:
		return HasilAktual{}, fmt.Errorf("%w: aksi adjust premium %q", ErrMasukanTidakSah, m.Aksi)
	}
	return HasilAktual{Actual: k.av, Selisih: vd, Akar: k.ubah, Pesan: pesanUnik(k.pesan)}, nil
}

// tulisAkar - tulis satu larik akar (juga ke salinan kerja).
func (k *konteksAktual) tulisAkar(n string, v []map[string]any) {
	k.akar.Larik[n] = v
	k.ubah.Larik[n] = v
}

// --- TreatyInActualUpdateValuePremiumIncome -----------------------------------

func (k *konteksAktual) premiIncome() {
	av := k.av
	// [1]
	delete(av.Larik, "TotalEgnpiAmountNP")
	delete(av.Medan, "TotalEgnpiAmount")
	delete(av.Medan, "TotalEgnpiProportion")
	// [2]
	egnpi := av.Larik["EGNPI"]
	total := apd.New(0, 0)
	for _, e := range egnpi {
		total = tambah(angka(teksSimpul(e, "AmountIDR")), total)
	}
	if len(egnpi) > 0 {
		av.Medan["TotalEgnpiAmount"] = teks(total)
	}
	// [3] `TreatyIn.TotalEgnpiAmount==0` → pesan `Param.errmsg1` (tak dikirim) dan KELUAR.
	if angka(k.akar.Medan["TotalEgnpiAmount"]).IsZero() {
		return
	}
	// [4]
	np := []map[string]any{}
	proporsi := apd.New(0, 0)
	for _, r := range egnpi {
		cur, nilai := teksSimpul(r, "Currency"), angka(teksSimpul(r, "Amount"))
		p := kali(bagiBulatDes(angka(teksSimpul(r, "AmountIDR")), total, 20), seratus)
		r["Proportion"] = teks(p)
		proporsi = tambah(p, proporsi)
		av.Medan["TotalEgnpiProportion"] = teks(proporsi)
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
	av.Larik["TotalEgnpiAmountNP"] = np
}

// --- TreatyInActualUpdateValueLimits ----------------------------------------

func (k *konteksAktual) aktualLimits() {
	// [2] ⚠️ salinan UTUH `TreatyIn.Limits`.
	lim := keJenis[[]map[string]any](k.akar.Larik["Limits"])
	egnpiAktual := k.av.Larik["EGNPI"]
	for _, l := range lim {
		// [3.2]–[3.3]
		l["AddendumStatus"] = "1"
		// [3.4] baris Actual GNPI ber-Treaty Group sama.
		cocok := []map[string]any{}
		for _, g := range larikSimpul(l, "TreatyGroupList") {
			grup := teksSimpul(g, "TreatyGroup")
			for _, e := range egnpiAktual {
				if teksSimpul(e, "TreatyGroup") == grup {
					cocok = append(cocok, e)
				}
			}
		}
		// [3.5] per mata uang, nilai = Amount.
		total := []map[string]any{}
		for _, e := range cocok {
			cur, nilai := teksSimpul(e, "Currency"), angka(teksSimpul(e, "Amount"))
			ada := false
			for _, t := range total {
				if teksSimpul(t, "Currency") == cur {
					t["Value"] = teks(tambah(nilai, angka(teksSimpul(t, "Value"))))
					ada = true
				}
			}
			if !ada {
				total = append(total, map[string]any{"Currency": cur, "Value": teks(nilai), "CurrencyID": teksSimpul(e, "CurrencyID")})
			}
		}
		l["EgnpiTotalList"] = total
		// [5]–[6] Premium Earned = EGNPI × @divide(AdjRate,100,15).
		rate := bagiBulat(angka(teksSimpul(l, "AdjRate")), 100, 15)
		pe := []map[string]any{}
		for _, e := range total {
			pe = append(pe, map[string]any{"Currency": teksSimpul(e, "Currency"), "Value": teks(kali(angka(teksSimpul(e, "Value")), rate))})
		}
		l["PremiumEarnedList"] = pe
		// [8] ⚠️ MDP% = 100.
		l["MDPPct"] = "100"
		satu := bagiBulat(seratus, 100, 15)
		mdp := []map[string]any{}
		for _, p := range pe {
			mdp = append(mdp, map[string]any{"Value": teks(kali(angka(teksSimpul(p, "Value")), satu)), "Currency": teksSimpul(p, "Currency")})
		}
		l["MDPList"] = mdp
	}
	k.av.Larik["Limits"] = lim
}

// --- TreatyInLimitsActualSetTotal / TreatyInNPSetTotalActualLimits ----------

// totalLimitsAktual - `np` = NPSetTotalActualLimits (langkah [10]/[11]
// berbeda: mata uang premi ditambahkan ke total AKAR).
func (k *konteksAktual) totalLimitsAktual(np bool) {
	av := k.av
	// [1]
	for _, n := range []string{"TotalLimitIOONP", "TotalLimitDeductblNP", "TotalLimitPremiEarnNP", "TotalLimitMDPNP"} {
		delete(av.Larik, n)
	}
	delete(av.Medan, "TotalLimitsROL")
	lim := av.Larik["Limits"]
	pertama := baris(lim, 0)
	jumlah := func(k1, k2 string) (string, string) {
		if len(lim) == 0 {
			return "", ""
		}
		a, b := apd.New(0, 0), apd.New(0, 0)
		for _, l := range lim {
			a = tambah(a, angka(teksSimpul(l, k1)))
			b = tambah(b, angka(teksSimpul(l, k2)))
		}
		return teks(a), teks(b)
	}
	// [4]–[5] / [3]–[4]
	v1, v2 := jumlah("Limit", "Limit2")
	av.Larik["TotalLimitIOONP"] = []map[string]any{
		{"Currency": teksSimpul(pertama, "Currency"), "Value": v1},
		{"Currency": teksSimpul(pertama, "Currency2"), "Value": v2},
	}
	// [7]–[8] / [6]–[7]
	d1, d2 := jumlah("Deductible", "Deductible2")
	av.Larik["TotalLimitDeductblNP"] = []map[string]any{
		{"Currency": teksSimpul(pertama, "Currency"), "Value": d1},
		{"Currency": teksSimpul(pertama, "Currency2"), "Value": d2},
	}
	// [9] / [8] mata uang Premium Earned, urut kemunculan.
	mataUang := []string{}
	for _, l := range lim {
		for _, p := range larikSimpul(l, "PremiumEarnedList") {
			cur := teksSimpul(p, "Currency")
			ada := false
			for _, c := range mataUang {
				ada = ada || c == cur
			}
			if !ada {
				mataUang = append(mataUang, cur)
			}
		}
	}
	pe, mdp := []map[string]any{}, []map[string]any{}
	if np {
		// [10] ⚠️ baris mata uang ke total AKAR.
		ap, am := k.akar.Larik["TotalLimitPremiEarnNP"], k.akar.Larik["TotalLimitMDPNP"]
		for _, c := range mataUang {
			ap = append(ap, map[string]any{"Currency": c})
			am = append(am, map[string]any{"Currency": c})
		}
		if len(mataUang) > 0 {
			k.tulisAkar("TotalLimitPremiEarnNP", ap)
			k.tulisAkar("TotalLimitMDPNP", am)
		}
	}
	// [11]/[13] — Σ nilai bermata uang sama.
	for _, c := range mataUang {
		sp, sm := apd.New(0, 0), apd.New(0, 0)
		for _, l := range lim {
			for _, p := range larikSimpul(l, "PremiumEarnedList") {
				if teksSimpul(p, "Currency") == c {
					sp = tambah(sp, angka(teksSimpul(p, "Value")))
				}
			}
			for _, x := range larikSimpul(l, "MDPList") {
				if teksSimpul(x, "Currency") == c {
					sm = tambah(sm, angka(teksSimpul(x, "Value")))
				}
			}
		}
		pe = append(pe, map[string]any{"Currency": c, "Value": teks(sp)})
		mdp = append(mdp, map[string]any{"Currency": c, "Value": teks(sm)})
	}
	av.Larik["TotalLimitPremiEarnNP"] = pe
	av.Larik["TotalLimitMDPNP"] = mdp
	// [14] / [12] ROL = Σ ROL% layer.
	if len(lim) > 0 {
		rol := apd.New(0, 0)
		for _, l := range lim {
			rol = tambah(rol, angka(teksSimpul(l, "ROLPct")))
		}
		av.Medan["TotalLimitsROL"] = teks(rol)
	}
}

// TreatyInSummaryLimitActual — struktur TreatyInSummaryLimit atas ActualValue.Limits.
func (k *konteksAktual) ringkasLimitAktual() {
	k.av.Larik["LimitSummaryList"] = keJenis[[]map[string]any](SummaryLimit(keJenis[[]LayerNP](k.av.Larik["Limits"])))
}

// TreatyInSummaryLimitShareActual — struktur TreatyInSummaryLimitShare atas ActualValue.Share.
func (k *konteksAktual) ringkasShareAktual() {
	k.av.Larik["LimitShareSummaryList"] = ringkasPohon(k.av.Larik["Share"])
}

// --- TreatyInActualShare ----------------------------------------------------

func (k *konteksAktual) aktualShare() {
	av, akar := k.av, k.akar
	// [1]
	share, fac := []map[string]any{}, []map[string]any{}
	// [2]
	for _, n := range []string{"RNMShare", "BrokeragePercent", "FacultativeShare", "FacultativeShareBrokerage"} {
		av.Medan[n] = akar.Medan[n]
	}
	hitung := angka(av.Medan["RNMShare"])
	lim := av.Larik["Limits"]
	// [3] FacultativeShare > 0 → [4]–[5]; selain itu lompat ke `shr` [6].
	if angka(av.Medan["FacultativeShare"]).Sign() > 0 {
		hitung = kurang(angka(av.Medan["RNMShare"]), angka(av.Medan["FacultativeShare"]))
		k.ubah.Medan["RnmShareDeducted"] = teks(hitung)
		fAkar := bagiBulat(angka(akar.Medan["FacultativeShare"]), 100, 4)
		fAv := bagiBulat(angka(av.Medan["FacultativeShare"]), 100, 4)
		for _, l := range lim {
			b := barisShareDariLimit(l)
			rnm := []map[string]any{{"Currency": teksSimpul(l, "Currency"), "Value": teks(kali(angka(teksSimpul(l, "Limit")), fAkar))}}
			b["Limit"], b["Limit2"] = teksSimpul(l, "Limit"), teksSimpul(l, "Limit2")
			// [4.2]
			if angka(teksSimpul(l, "Limit2")).Cmp(apd.New(1, 0)) >= 0 {
				rnm = append(rnm, map[string]any{"Currency": teksSimpul(l, "Currency2"), "Value": teks(kali(angka(teksSimpul(l, "Limit2")), fAv))})
			}
			b["RnmLimitList"] = rnm
			b["GrossPremiumList"] = kaliNilai(larikSimpul(l, "MDPList"), fAv)
			b["NetPremiumList"] = keJenis[[]map[string]any](b["GrossPremiumList"])
			fac = append(fac, b)
		}
		// [5]
		reins := av.Larik["ShareFacultativeReinsurers"]
		for _, b := range fac {
			r := keJenis[[]map[string]any](reins)
			for _, x := range r {
				pct := bagiBulat(angka(teksSimpul(x, "SharePct")), 100, 2)
				x["Amount"] = teks(kali(angka(teksSimpul(b, "Limit")), pct))
				x["Amount2"] = teks(kali(angka(teksSimpul(b, "Limit2")), pct))
			}
			b["ShareFacultativeReinsurers"] = r
			b["RnmLimitListDisplay"] = keJenis[[]map[string]any](b["RnmLimitList"])
			b["RnmGrossPremiDisplay"] = keJenis[[]map[string]any](b["GrossPremiumList"])
		}
	}
	// [6] baris Share per layer Actual.
	r := bagiBulat(hitung, 100, 4)
	shareAkar := akar.Larik["Share"]
	for i, l := range lim {
		s := baris(shareAkar, i)
		b := barisShareDariLimit(l)
		b["SpreadingTypeXOL"] = teksSimpul(s, "SpreadingTypeXOL")
		b["SpreadingTotalPctXOL"] = teksSimpul(s, "SpreadingTotalPctXOL")
		b["AddendumStatus"] = "1"
		spread := []map[string]any{}
		for _, x := range larikSimpul(s, "SpreadingListXOL") {
			spread = append(spread, map[string]any{"Pct": teksSimpul(x, "Pct"), "ReinsTypeName": teksSimpul(x, "ReinsTypeName")})
		}
		b["SpreadingListXOL"] = spread
		rnm := []map[string]any{{"Currency": teksSimpul(l, "Currency"), "Value": teks(kali(angka(teksSimpul(l, "Limit")), r))}}
		if angka(teksSimpul(l, "Limit2")).Cmp(apd.New(1, 0)) >= 0 {
			rnm = append(rnm, map[string]any{"Currency": teksSimpul(l, "Currency2"), "Value": teks(kali(angka(teksSimpul(l, "Limit2")), r))})
		}
		b["RnmLimitList"] = rnm
		b["GrossPremiumList"] = kaliNilai(larikSimpul(l, "MDPList"), r)
		// [7]
		b["RNMShare"] = akar.Medan["RNMShare"]
		b["NetPremiumList"] = keJenis[[]map[string]any](b["GrossPremiumList"])
		b["RnmLimitListDisplay"] = keJenis[[]map[string]any](b["RnmLimitList"])
		b["RnmGrossPremiDisplay"] = keJenis[[]map[string]any](b["GrossPremiumList"])
		share = append(share, b)
	}
	// [8] sebaran OR/RI Limit dan Gross; QS% TERBAWA antarbaris.
	q := apd.New(0, 0)
	for _, b := range share {
		q = qsOR(b, q)
		orP, riP := bagiPolos(q, seratus), bagiPolos(kurang(seratus, q), seratus)
		b["RNMSpreadedListXOL"], b["RNMSpreadedListRIXOL"] = sebarDua(larikSimpul(b, "RnmLimitList"), orP, riP)
		b["RNMSpreadedListGrossXOL"], b["RNMSpreadedListGrossRIXOL"] = sebarDua(larikSimpul(b, "GrossPremiumList"), orP, riP)
	}
	av.Larik["Share"] = share
	av.Larik["FacultativeShareList"] = fac
}

// barisShareDariLimit - kunci layer yang ActualShare salin ke baris Share/Fac.
func barisShareDariLimit(l map[string]any) map[string]any {
	b := map[string]any{
		"ClassofBusinessList": salinPohon(larikSimpul(l, "LayerList")),
		"TreatyGroupList":     salinPohon(larikSimpul(l, "TreatyGroupList")),
	}
	for _, kunci := range []string{"LayerType", "Layer", "LayerPartType", "LayerPart", "Cover"} {
		b[kunci] = teksSimpul(l, kunci)
	}
	return b
}

func kaliNilai(xs []map[string]any, f *apd.Decimal) []map[string]any {
	out := []map[string]any{}
	for _, x := range xs {
		out = append(out, map[string]any{"Currency": teksSimpul(x, "Currency"), "Value": teks(kali(angka(teksSimpul(x, "Value")), f))})
	}
	return out
}

// qsOR - `local.QSPCT` dari baris `QS (OR)`; tanpa baris itu nilai lama terbawa.
func qsOR(b map[string]any, q *apd.Decimal) *apd.Decimal {
	for _, s := range larikSimpul(b, "SpreadingListXOL") {
		if teksSimpul(s, "ReinsTypeName") == namaQSOR {
			q = angka(teksSimpul(s, "Pct"))
		}
	}
	return q
}

// sebarDua - (QS%/100)×nilai dan ((100−QS%)/100)×nilai per baris.
func sebarDua(xs []map[string]any, orP, riP *apd.Decimal) ([]map[string]any, []map[string]any) {
	or, ri := []map[string]any{}, []map[string]any{}
	for _, x := range xs {
		v := angka(teksSimpul(x, "Value"))
		or = append(or, map[string]any{"Value": teks(kali(orP, v)), "Currency": teksSimpul(x, "Currency")})
		ri = append(ri, map[string]any{"Value": teks(kali(riP, v)), "Currency": teksSimpul(x, "Currency")})
	}
	return or, ri
}

// --- TreatyInSetBrokerageActual ---------------------------------------------

func (k *konteksAktual) brokerageAktual() {
	av := k.av
	share := av.Larik["Share"]
	// [1]–[3]
	for _, b := range share {
		b["DeductionList"] = []map[string]any{{"Comment": "Brokerage fee", "DeductionPct": av.Medan["BrokeragePercent"], "DeductionPctCalculate": "true"}}
		k.pesan = append(k.pesan, deduksiPohon(b)...)
	}
	// [4]–[5]
	for _, b := range av.Larik["FacultativeShareList"] {
		b["DeductionList"] = append(larikSimpul(b, "DeductionList"), map[string]any{
			"Comment": "Facultative Brokerage fee", "DeductionPct": av.Medan["FacultativeShareBrokerage"], "DeductionPctCalculate": "true",
		})
		k.pesan = append(k.pesan, deduksiPohon(b)...)
	}
	// [6] TreatyInSummaryLimitActualFacShare — ⚠️ akar dibuang, Actual per indeks.
	k.tulisAkar("LimitFacShareSummaryList", []map[string]any{})
	ringkas := ringkasPohon(av.Larik["FacultativeShareList"])
	tujuan := av.Larik["LimitFacShareSummaryList"]
	for i, r := range ringkas {
		tujuan, _ = pasang(tujuan, i)
		tujuan[i] = r
	}
	if tujuan != nil {
		av.Larik["LimitFacShareSummaryList"] = tujuan
	}
	// [7] buang SATU deduksi < 1 — yang terakhir.
	hapusDeduksiKecil(share)
	// [8] sebaran Deduction dan Net (QS% terbawa).
	q := apd.New(0, 0)
	for _, b := range share {
		q = qsOR(b, q)
		orP, riP := bagiPolos(q, seratus), bagiPolos(kurang(seratus, q), seratus)
		b["RNMSpreadedListDeductXOL"], b["RNMSpreadedListDeductRIXOL"] = sebarDua(larikSimpul(b, "DeductionTotalList"), orP, riP)
		b["RNMSpreadedListNetXOL"], b["RNMSpreadedListNetRIXOL"] = sebarDua(larikSimpul(b, "NetPremiumList"), orP, riP)
	}
}

func hapusDeduksiKecil(share []map[string]any) {
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
}

// --- TreatyInShareActualSetTotal / TreatyInNPSetTotalActualShare ------------

func (k *konteksAktual) totalShareAktual(np bool) {
	av := k.av
	// [1]–[2]
	for _, n := range []string{"TotalShareRnmNP", "TotalShareGrossNP", "TotalShareNetNP", "TotalShareDeductionNP"} {
		av.Larik[n] = []map[string]any{}
	}
	for _, b := range av.Larik["Share"] {
		for _, p := range [][2]string{
			{"GrossPremiumList", "TotalShareGrossNP"}, {"NetPremiumList", "TotalShareNetNP"},
			{"DeductionTotalList", "TotalShareDeductionNP"}, {"RnmLimitList", "TotalShareRnmNP"},
		} {
			for _, x := range larikSimpul(b, p[0]) {
				av.Larik[p[1]] = jumlahPohon(av.Larik[p[1]], x, true)
			}
		}
	}
	// [3] ShareActualSetTotal membuang ActualValue.TotalFac*; NPSetTotalActualShare
	// ⚠️ membuang total AKAR dan MENGAKUMULASI ke ActualValue.
	fac := []string{"TotalFacShareRnmNP", "TotalFacShareGrossNP", "TotalFacShareNetNP", "TotalFacShareDeductionNP"}
	for _, n := range fac {
		if np {
			k.tulisAkar(n, []map[string]any{})
		} else {
			av.Larik[n] = []map[string]any{}
		}
	}
	// [4]
	for _, b := range av.Larik["FacultativeShareList"] {
		for _, p := range [][2]string{
			{"GrossPremiumList", "TotalFacShareGrossNP"}, {"NetPremiumList", "TotalFacShareNetNP"},
			{"DeductionTotalList", "TotalFacShareDeductionNP"}, {"RnmLimitList", "TotalFacShareRnmNP"},
		} {
			for _, x := range larikSimpul(b, p[0]) {
				av.Larik[p[1]] = jumlahPohon(av.Larik[p[1]], x, true)
			}
		}
	}
	if np {
		return
	}
	// [5]–[6]
	for _, n := range []string{"TotalSpreadedRnmProp", "TotalSpreadedRnmRIProp", "TotalSpreadedNetPremi", "TotalSpreadedNetPremiRI"} {
		av.Larik[n] = []map[string]any{}
	}
	for _, b := range av.Larik["Share"] {
		for _, p := range [][2]string{
			{"RNMSpreadedListXOL", "TotalSpreadedRnmProp"}, {"RNMSpreadedListRIXOL", "TotalSpreadedRnmRIProp"},
			{"RNMSpreadedListNetXOL", "TotalSpreadedNetPremi"}, {"RNMSpreadedListNetRIXOL", "TotalSpreadedNetPremiRI"},
		} {
			for _, x := range larikSimpul(b, p[0]) {
				av.Larik[p[1]] = jumlahPohon(av.Larik[p[1]], x, false)
			}
		}
	}
}

// --- Difference* (Actual − Treaty In) ---------------------------------------

// nolBilaKurang - `@if(.X < estimasi, 0, .X − estimasi)`.
func nolBilaKurang(a, b map[string]any, kunci string) string {
	x, y := angka(teksSimpul(a, kunci)), angka(teksSimpul(b, kunci))
	if x.Cmp(y) < 0 {
		return "0"
	}
	return teks(kurang(x, y))
}

func selisihNilaiNol(kini, dulu []map[string]any) []map[string]any {
	out := []map[string]any{}
	for j, x := range kini {
		out = append(out, map[string]any{"Currency": teksSimpul(x, "Currency"), "Value": nolBilaKurang(x, baris(dulu, j), "Value")})
	}
	return out
}

var kolomRingkasShare = []string{"Limit", "Limit2", "Deductible", "Deductible2", "MDP", "MDP2", "AggregateLimit", "AggregateLimit2"}

// selisih - ValueDifference disusun ulang (`Property-Remove` lalu Difference*);
// `fac` = TreatyInActualUpdateValueShare (dengan DifferenceFacShare*).
func (k *konteksAktual) selisih(fac bool) HalamanPohon {
	av, akar := k.av, k.akar
	vd := halamanKosong()
	// DifferencePremium — struktur TreatyEDMDifferencePremium, Actual vs akar.
	selisihPremi(av, akar, &vd)
	// DifferenceLimits [1]–[2]; [3] berblok `//`.
	lim := []map[string]any{}
	for i, l := range av.Larik["Limits"] {
		o := baris(akar.Larik["Limits"], i)
		v := map[string]any{"TreatyGroupList": salinPohon(larikSimpul(l, "TreatyGroupList"))}
		for _, kunci := range []string{"LayerType", "Layer", "LayerPartType", "LayerPart", "Cover", "ReinstatementPct", "Currency", "Currency2"} {
			v[kunci] = teksSimpul(l, kunci)
		}
		for _, kunci := range []string{"Limit", "Limit2", "Deductible", "Deductible2", "AdjRate", "MDPPct", "ROLPct"} {
			v[kunci] = beda(l, o, kunci)
		}
		for _, n := range []string{"EgnpiTotalList", "PremiumEarnedList", "MDPList"} {
			v[n] = selisihNilaiNol(larikSimpul(l, n), larikSimpul(o, n))
		}
		lim = append(lim, v)
	}
	vd.Larik["Limits"] = lim
	// DifferenceLimitsSumary — struktur TreatyInSummaryLimit atas ValueDifference.Limits.
	vd.Larik["LimitSummaryList"] = keJenis[[]map[string]any](SummaryLimit(keJenis[[]LayerNP](lim)))
	// DifferenceLimitsSetTotal — tanpa total limit/deductible (`type` tak dikirim).
	k2 := &konteksAktual{av: vd, akar: halamanKosong(), ubah: halamanKosong()}
	k2.totalLimitsAktual(false)
	delete(vd.Larik, "TotalLimitIOONP")
	delete(vd.Larik, "TotalLimitDeductblNP")
	// DifferenceShare
	vd.Medan["RNMShare"] = teks(kurang(angka(av.Medan["RNMShare"]), angka(akar.Medan["RNMShare"])))
	vd.Medan["BrokeragePercent"] = teks(kurang(angka(av.Medan["BrokeragePercent"]), angka(akar.Medan["BrokeragePercent"])))
	vd.Larik["Share"] = selisihBarisShare(av.Larik["Share"], akar.Larik["Share"], true)
	// [3] ⚠️ ActualValue.LimitShareSummaryList ditimpa selisihnya.
	for i, r := range av.Larik["LimitShareSummaryList"] {
		o := baris(akar.Larik["LimitShareSummaryList"], i)
		for _, kunci := range kolomRingkasShare {
			r[kunci] = nolBilaKurang(r, o, kunci)
		}
	}
	// DifferenceDeduction ([4]–[5] berblok `//`).
	for _, b := range vd.Larik["Share"] {
		b["DeductionList"] = []map[string]any{{"Comment": "Brokerage fee", "DeductionPct": akar.Medan["BrokeragePercent"], "DeductionPctCalculate": "true"}}
		k.pesan = append(k.pesan, deduksiPohon(b)...)
	}
	hapusDeduksiKecil(vd.Larik["Share"])
	// DifferenceShareSumary, DifferenceShareTotal.
	vd.Larik["LimitShareSummaryList"] = ringkasPohon(vd.Larik["Share"])
	totalSelisihShare(&vd)
	if !fac {
		return vd
	}
	// DifferenceFacShare
	vd.Medan["FacultativeShare"] = teks(kurang(angka(av.Medan["FacultativeShare"]), angka(akar.Medan["FacultativeShare"])))
	vd.Medan["FacultativeShareBrokerage"] = teks(kurang(angka(av.Medan["FacultativeShareBrokerage"]), angka(akar.Medan["FacultativeShareBrokerage"])))
	vd.Larik["FacultativeShareList"] = selisihBarisShare(av.Larik["FacultativeShareList"], akar.Larik["FacultativeShareList"], false)
	for i, r := range av.Larik["LimitShareSummaryList"] {
		o := baris(akar.Larik["LimitShareSummaryList"], i)
		for _, kunci := range kolomRingkasShare {
			r[kunci] = nolBilaKurang(r, o, kunci)
		}
	}
	// DifferenceSummaryFacShare
	vd.Larik["LimitFacShareSummaryList"] = ringkasPohon(vd.Larik["FacultativeShareList"])
	// DifferenceFacShareTotal — CalculateDeduction di tengah ([2.3]).
	for _, n := range []string{"TotalFacShareRnmNP", "TotalFacShareGrossNP", "TotalFacShareNetNP", "TotalFacShareDeductionNP"} {
		vd.Larik[n] = []map[string]any{}
	}
	for _, b := range vd.Larik["FacultativeShareList"] {
		for _, x := range larikSimpul(b, "GrossPremiumList") {
			vd.Larik["TotalFacShareGrossNP"] = jumlahPohon(vd.Larik["TotalFacShareGrossNP"], x, true)
		}
		for _, x := range larikSimpul(b, "NetPremiumList") {
			vd.Larik["TotalFacShareNetNP"] = jumlahPohon(vd.Larik["TotalFacShareNetNP"], x, true)
		}
		k.pesan = append(k.pesan, deduksiPohon(b)...)
		for _, x := range larikSimpul(b, "DeductionTotalList") {
			vd.Larik["TotalFacShareDeductionNP"] = jumlahPohon(vd.Larik["TotalFacShareDeductionNP"], x, true)
		}
		for _, x := range larikSimpul(b, "RnmLimitList") {
			vd.Larik["TotalFacShareRnmNP"] = jumlahPohon(vd.Larik["TotalFacShareRnmNP"], x, true)
		}
	}
	return vd
}

// selisihBarisShare - DifferenceShare [2] (`share` = true: spreading dari
// Share akar dan larik sebaran) / DifferenceFacShare [2] (Limit/Limit2 disalin).
func selisihBarisShare(kini, dulu []map[string]any, share bool) []map[string]any {
	out := []map[string]any{}
	for i, s := range kini {
		o := baris(dulu, i)
		v := map[string]any{"AddendumStatus": "1"}
		for _, kunci := range []string{"LayerType", "Layer", "LayerPartType", "LayerPart", "Cover"} {
			v[kunci] = teksSimpul(s, kunci)
		}
		if share {
			v["TreatyGroupList"] = salinPohon(larikSimpul(s, "TreatyGroupList"))
			for _, kunci := range []string{"SpreadingTypeXOL", "SpreadingTypeIDXOL", "SpreadingTotalPctXOL"} {
				v[kunci] = teksSimpul(o, kunci)
			}
			spread := []map[string]any{}
			for _, x := range larikSimpul(o, "SpreadingListXOL") {
				spread = append(spread, map[string]any{"Pct": teksSimpul(x, "Pct"), "ReinsTypeName": teksSimpul(x, "ReinsTypeName"), "ReinsTypeID": teksSimpul(x, "ReinsTypeID")})
			}
			v["SpreadingListXOL"] = spread
		} else {
			delete(v, "AddendumStatus")
			v["Limit"], v["Limit2"] = teksSimpul(s, "Limit"), teksSimpul(s, "Limit2")
		}
		for _, pas := range [][2]string{{"RnmLimitListDisplay", "RnmLimitList"}, {"RnmGrossPremiDisplay", "GrossPremiumList"}} {
			a, b := larikTampil(s, pas[0], pas[1]), larikTampil(o, pas[0], pas[1])
			t := make([]map[string]any, 2)
			for j := range t {
				x := baris(a, j)
				t[j] = map[string]any{"Currency": teksSimpul(x, "Currency"), "Value": beda(x, baris(b, j), "Value")}
			}
			v[pas[0]] = t
		}
		v["RnmLimitList"] = selisihNilai(larikSimpul(s, "RnmLimitList"), larikSimpul(o, "RnmLimitList"), true)
		v["GrossPremiumList"] = selisihNilaiNol(larikSimpul(s, "GrossPremiumList"), larikSimpul(o, "GrossPremiumList"))
		ded := []map[string]any{}
		dl := larikSimpul(o, "DeductionList")
		for j, d := range larikSimpul(s, "DeductionList") {
			x := baris(dl, j)
			ded = append(ded, map[string]any{
				"Currency": teksSimpul(d, "Currency"), "Comment": teksSimpul(d, "Comment"),
				"Deduction": nolBilaKurang(d, x, "Deduction"), "DeductionPct": nolBilaKurang(d, x, "DeductionPct"),
			})
		}
		v["DeductionList"] = ded
		v["NetPremiumList"] = selisihNilaiNol(larikSimpul(s, "NetPremiumList"), larikSimpul(o, "NetPremiumList"))
		if share {
			for _, n := range []string{
				"RNMSpreadedListXOL", "RNMSpreadedListRIXOL", "RNMSpreadedListGrossXOL", "RNMSpreadedListGrossRIXOL",
				"RNMSpreadedListDeductXOL", "RNMSpreadedListDeductRIXOL", "RNMSpreadedListNetXOL", "RNMSpreadedListNetRIXOL",
			} {
				v[n] = selisihNilaiNol(larikSimpul(s, n), larikSimpul(o, n))
			}
		}
		out = append(out, v)
	}
	return out
}

// --- TreatyInXOLAddSpreadingDetailActual ------------------------------------

func (k *konteksAktual) xolDetailAktual(idx int) {
	av := k.av
	share := av.Larik["Share"]
	if idx < 0 || idx >= len(share) {
		return
	}
	b := share[idx]
	// [1]
	for _, n := range []string{"TotalSpreadedRnmProp", "TotalSpreadedRnmRIProp", "TotalSpreadedNetPremi", "TotalSpreadedNetPremiRI"} {
		delete(av.Larik, n)
	}
	for _, n := range []string{"RNMSpreadedListDeductXOL", "RnmLimitList", "RNMSpreadedListNetXOL", "RNMSpreadedListNetRIXOL", "GrossPremiumList"} {
		delete(b, n)
	}
	// [2] dari ActualValue.Limits(idx) × @divide(RNMShare baris,100,4).
	l := baris(av.Larik["Limits"], idx)
	r := bagiBulat(angka(teksSimpul(b, "RNMShare")), 100, 4)
	rnm := []map[string]any{{"Currency": teksSimpul(l, "Currency"), "Value": teks(kali(angka(teksSimpul(l, "Limit")), r))}}
	if angka(teksSimpul(l, "Limit2")).Cmp(apd.New(1, 0)) >= 0 {
		rnm = append(rnm, map[string]any{"Currency": teksSimpul(l, "Currency2"), "Value": teks(kali(angka(teksSimpul(l, "Limit2")), r))})
	}
	b["RnmLimitList"] = rnm
	b["GrossPremiumList"] = kaliNilai(larikSimpul(l, "MDPList"), r)
	b["NetPremiumList"] = keJenis[[]map[string]any](b["GrossPremiumList"])
	// [3]
	b["RNMSpreadedListXOL"], b["RNMSpreadedListRIXOL"] = []map[string]any{}, []map[string]any{}
	// [4.1] FetchQSfromMasterXOL TANPA induk (struktur rute share-np), [4.2] CalculateDeduction.
	ks := konteksShare{idKontrak: k.akar.Medan["ID"]}
	t := keJenis[models.BarisShareNP](b)
	ks.fetchQS(&t, "", false)
	for kunci, v := range keJenis[map[string]any](t) {
		if _, adaDulu := b[kunci]; adaDulu || !kosongPohon(v) {
			b[kunci] = v
		}
	}
	k.pesan = append(k.pesan, deduksiPohon(b)...)
	// [5] (5.2/5.3/5.4.1 berblok `//`).
	q := qsOR(b, apd.New(0, 0))
	orP, riP := bagiPolos(q, seratus), bagiPolos(kurang(seratus, q), seratus)
	b["RNMSpreadedListNetXOL"], b["RNMSpreadedListNetRIXOL"] = sebarDua(larikSimpul(b, "NetPremiumList"), orP, riP)
	b["RnmLimitListDisplay"] = keJenis[[]map[string]any](b["RnmLimitList"])
	b["RnmGrossPremiDisplay"] = keJenis[[]map[string]any](b["GrossPremiumList"])
	b["RNMSpreadedListDeductXOL"], b["RNMSpreadedListDeductRIXOL"] = sebarDua(larikSimpul(b, "DeductionTotalList"), orP, riP)
	// [6] ⚠️ total ActualValue (kosong sejak [1]) tak pernah cocok — tiap
	// baris sebaran DITAMBAHKAN ke total AKAR.
	for _, p := range [][2]string{
		{"RNMSpreadedListXOL", "TotalSpreadedRnmProp"}, {"RNMSpreadedListRIXOL", "TotalSpreadedRnmRIProp"},
		{"RNMSpreadedListNetXOL", "TotalSpreadedNetPremi"}, {"RNMSpreadedListNetRIXOL", "TotalSpreadedNetPremiRI"},
	} {
		total := keJenis[[]map[string]any](k.akar.Larik[p[1]])
		berubah := false
		for _, s := range share {
			for _, x := range larikSimpul(s, p[0]) {
				total = append(total, map[string]any{"Currency": teksSimpul(x, "Currency"), "Value": teks(angka(teksSimpul(x, "Value")))})
				berubah = true
			}
		}
		if berubah {
			k.tulisAkar(p[1], total)
		}
	}
}

// kosongPohon - nilai Go kosong (`""`, `[]`, nil).
func kosongPohon(v any) bool {
	switch x := v.(type) {
	case nil:
		return true
	case string:
		return x == ""
	case []any:
		return len(x) == 0
	}
	return false
}
