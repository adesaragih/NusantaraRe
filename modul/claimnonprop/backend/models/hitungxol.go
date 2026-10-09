package models

// Untuk apa berkas ini: MESIN HITUNG XoL tingkat klaim - port `Activity/CountClaimTNP_Act.xml` (klaim 100% per mata
// uang), `Activity/CountLossAllocation_act.xml` (Loss Allocation + waterfall layer XoL + spreading), dan
// `Activity/AdjClaimAmount_Act.xml` (penyesuaian alokasi satu layer), serta `Activity/SetActualPremium_ACT.xml`.
//
// Nomor langkah di komentar = nomor langkah Pega. Langkah ber-REMARK (`//`) tidak ditulis. Hardcode per kasus / master
// (CLMNP-367/382, master 1000393) DIBUANG (bawaan OQ-CNP-03). Rumus dibangun persis XML termasuk keanehannya yang tidak
// masuk daftar perbaikan OQ-CNP-05 (mis. konversi tak simetris Limit2 / Deductible2, retensi `.Deductible*0` bila mata
// uang klaim = mata uang layer - OQ-CNP-18); setiap keanehan dicatat di docs/PARITAS.md.
//
// Medan TURUNAN (tidak disimpan): `ListTotalEstimation` (Summary XOL) disusun ulang dari SpreadingRisk dengan rumus
// CountLossAllocation_act langkah 19 (`SusunSummaryXOL`).

import (
	"strings"

	"github.com/cockroachdb/apd/v3"
)

// Pesan VERBATIM CountLossAllocation_act langkah 2.
const PesanNamaTreatyXOL = "Error TreatyName, Hubungi IT"

// Nilai baris UR (retensi) - CountLossAllocation_act langkah 17.2.5.
const TreatyUR = "UR"

// IDMataUangIDR - CurrencyID IDR yang tertulis di rumus (`=="10026"`, CountClaimTNP_Act langkah 13.2.3,
// SetCurrency_Act langkah 7.2 / 8.3).
const IDMataUangIDR = "10026"

// HitungKlaim = `CountClaimTNP_Act` lalu `CountLossAllocation_act` (Calculation "CountXOL", langkah 15).
func HitungKlaim(k *Konteks, h *Halaman, m MasterTreaty) error {
	var kal Kalkulator
	// 3-5: alokasi lama dibuang (TempAlokasiOld tidak dibaca lagi di activity ini).
	h.SetelDaftar(DaftarXOL, nil)
	h.SetelDaftar(DaftarSpreading, nil)
	h.SetelDaftar(DaftarBreakQS, nil)
	limits := limitGrup(h, m)                    // 6
	kursLayer, kursExc := kursLapisan(limits, m) // 7-8
	share := kal.H(h, CD+"ShareCeding")
	deductOn := h.Ambil(CD+"DeductibleType") == "true"
	baris := h.AmbilDaftar(DaftarClaimAmount)

	// 10: total nilai klaim per mata uang + prorata.
	type totalMU struct {
		cur        string
		totalClaim *apd.Decimal // CNPTotalClaim
	}
	var temp []*totalMU
	nilaiKonversi := apd.New(0, 0)
	totalAlokasi := apd.New(0, 0)
	tplAkhir := ""
	for i, r := range baris {
		deductFin := apd.New(0, 0)
		if deductOn { // 10.1
			deduct := kal.H(h, CD+"DeductibleValue")
			dasar := kal.B(r, "CNPTSI")
			if h.Ambil(CD+"TypeDeductible") == "1" {
				dasar = kal.B(r, "Value")
			}
			pctValue := kal.BagiPega(kal.Kali(dasar, kal.H(h, CD+"Amount")), k100())
			minMax := h.Ambil(CD+"CNPDeducMinMax") == "1"
			switch {
			case minMax:
				deductFin = maks(deduct, pctValue)
			case h.Ambil(CD+"FormType") == "1":
				deductFin = pctValue
			default:
				deductFin = deduct
			}
		}
		tplAkhir = r["TPL"] // 10.2 Local.TPLAmTC
		ded := apd.New(0, 0)
		if i == 0 {
			ded = deductFin
		}
		r["CNPDeductible"] = Teks(ded)
		nilai := kal.Tambah(kal.Kurang(kal.B(r, "Value"), ded), kal.B(r, "TPL"))
		kurs := kal.B(r, "AltValue")
		bersihPenuh := kal.Kurang(kal.Tambah(kal.Tambah(nilai, kal.B(r, "CNPOthersFee")), kal.B(r, "AdjusterFee")),
			kal.B(r, "Salvage"))
		totalAll := kal.Kali(bersihPenuh, kurs)
		nilaiKonversi = kal.Tambah(nilaiKonversi, totalAll)
		r["USD"] = Teks(kal.Kali(kal.B(r, "Value"), kurs))
		r["ClaimAmountCedant"] = Teks(kal.Kali(nilai, kal.BagiPega(share, k100())))
		// 10.3 TotalAllocation (gerbang langkah 12)
		var porsi *apd.Decimal
		switch {
		case kursLayer == "IDR":
			porsi = kal.B(r, "USD")
		case kursLayer == r["Currency"]:
			porsi = nilai
		default:
			porsi = kal.BagiPega(nilai, kursExc)
		}
		totalAlokasi = kal.Tambah(totalAlokasi, porsi)
		// 10.4-10.5 TempTotalClaim per mata uang
		ketemu := false
		for _, t := range temp {
			if t.cur == r["Currency"] {
				t.totalClaim = kal.Tambah(t.totalClaim, totalAll)
				ketemu = true
			}
		}
		if !ketemu {
			temp = append(temp, &totalMU{cur: r["Currency"], totalClaim: totalAll})
		}
	}
	if err := kal.Galat(); err != nil {
		return err
	}
	// 12: keluar bila total nol (SpreadingRisk dkk. sudah dibuang langkah 5).
	if !Lebih(nilaiKonversi, apd.New(0, 0)) && !Lebih(totalAlokasi, apd.New(0, 0)) {
		return nil
	}
	// 13: prorata per mata uang (10.4.1 AltValue = ValueKonversi akhir).
	for _, t := range temp {
		pr := kal.BagiPega(kal.Kali(t.totalClaim, k100()), nilaiKonversi)
		for _, r := range baris {
			if r["Currency"] == t.cur {
				r["PctProrateClaim"] = Teks(pr)
			}
		}
	}
	// 13.2.2-13.2.4 lalu 14: dua putaran rumus yang sama; putaran kedua memakai CNPDeductible hasil putaran pertama.
	for putaran := 0; putaran < 2; putaran++ {
		for _, r := range baris {
			if err := hitungBarisKlaim(&kal, h, r, tplAkhir, share); err != nil {
				return err
			}
		}
	}
	if err := kal.Galat(); err != nil {
		return err
	}
	return HitungXOL(k, h, m, "CountXOL", 0, "") // 15
}

// hitungBarisKlaim = CountClaimTNP_Act langkah 13.2.2-13.2.4 / 14.1-14.3 untuk satu baris ListClaimAmount.
func hitungBarisKlaim(kal *Kalkulator, h *Halaman, r Baris, tplAkhir string, share *apd.Decimal) error {
	isTPL := h.Ambil(CD+"IsTPL") == "true"
	nilai := kal.B(r, "Value")
	if isTPL { // 13.2.3 baris 2
		tpl, tplAm := kal.B(r, "TPL"), kal.Teks(".TPLAmTC", tplAkhir)
		nilai = kal.Tambah(nilai, maks2Lebih(tpl, tplAm))
	}
	ded := kal.B(r, "CNPDeductible")
	nilaiAdjust := kal.Kurang(nilai, ded)
	fee := kal.Kurang(kal.Tambah(kal.B(r, "CNPOthersFee"), kal.B(r, "AdjusterFee")), kal.B(r, "Salvage"))
	if isTPL {
		nilai = kal.Tambah(kal.Kurang(nilai, ded), fee)
	} else {
		nilai = kal.Tambah(kal.Tambah(kal.Kurang(nilai, ded), kal.B(r, "TPL")), fee)
	}
	r["CNPDeductible"] = Teks(deductibleBaris(kal, h, r))
	kurs := kal.B(r, "AltValue")
	pctShare := kal.BagiPega(share, k100())
	r["USD"] = Teks(kal.Kali(nilai, kurs))
	r["ClaimAmountCedant"] = Teks(kal.Kali(nilai, pctShare))
	claimValue := kal.Kali(nilaiAdjust, pctShare)
	r["ClaimAmountAdjust"] = Teks(claimValue)
	adjFee := kal.Kali(kal.B(r, "AdjusterFee"), pctShare)
	salvage := kal.Kali(kal.B(r, "Salvage"), pctShare)
	others := kal.Kali(kal.B(r, "CNPOthersFee"), pctShare)
	for _, l := range h.AmbilDaftar(DaftarLossAlloc) { // 13.2.4 / 14.3
		if l["Currency"] != r["Currency"] || l["CNPFlagOuts"] == "1" {
			continue
		}
		p := kal.BagiPega(kal.B(l, "ClaimPercentage"), k100())
		l["ClaimAmountAdjust"] = Teks(kal.Kali(claimValue, p))
		l["AdjusterFee"] = Teks(kal.Kali(adjFee, p))
		l["CNPOthersFee"] = Teks(kal.Kali(others, p))
		l["Salvage"] = Teks(kal.Kali(salvage, p))
		l["Currency"] = r["Currency"]
		l["CurrencyID"] = r["CurrencyID"]
	}
	return kal.Galat()
}

// deductibleBaris = CountClaimTNP_Act langkah 13.2.3 `.CNPDeductible` (berlaku untuk SETIAP baris - syarat
// `.pxListSubscript=="1"||.pxListSubscript!="1"` selalu benar). Cabang terakhir memakai DeductibleValue (bukan
// TSIDeductible) - ditiru apa adanya.
func deductibleBaris(kal *Kalkulator, h *Halaman, r Baris) *apd.Decimal {
	if h.Ambil(CD+"DeductibleType") != "true" {
		return apd.New(0, 0)
	}
	cur, kurs := h.Ambil(CD+"CurrencyDeductible"), kal.B(r, "AltValue")
	nilaiDed := kal.H(h, CD+"DeductibleValue")
	form := h.Ambil(CD + "FormType")
	if form == "1" || (form == "2" && Lebih(nilaiDed, apd.New(0, 0))) {
		switch {
		case cur == r["CurrencyID"]:
			return nilaiDed
		case cur == IDMataUangIDR:
			return kal.BagiPega(nilaiDed, kurs)
		default:
			return kal.Kali(nilaiDed, kurs)
		}
	}
	pct := kal.BagiPega(kal.H(h, CD+"Amount"), k100())
	switch {
	case h.Ambil(CD+"TypeDeductible") == "1":
		return kal.Kali(kal.B(r, "Value"), pct)
	case cur == r["CurrencyID"]:
		return kal.Kali(kal.H(h, CD+"TSIDeductible"), pct)
	case cur == IDMataUangIDR:
		return kal.Kali(kal.BagiPega(kal.H(h, CD+"TSIDeductible"), kurs), pct)
	default:
		return kal.Kali(kal.Kali(nilaiDed, kurs), pct)
	}
}

// limitGrup = CountLossAllocation_act langkah 11-12 (CountClaimTNP_Act langkah 6): layer master yang TreatyGroupList-nya
// memuat `TreatyInMaster.TreatyGroup`.
func limitGrup(h *Halaman, m MasterTreaty) []LimitXOL {
	grup := h.Ambil(TM + "TreatyGroup")
	var out []LimitXOL
	for _, l := range m.Limits {
		for _, g := range l.TreatyGroups {
			if g == grup {
				out = append(out, l)
				break
			}
		}
	}
	return out
}

// kursLapisan = langkah 14-15: mata uang layer pertama dan kursnya (`CurrencyList.Conversion`, baris terakhir yang
// cocok; tanpa layer / tanpa kurs = 1).
func kursLapisan(limits []LimitXOL, m MasterTreaty) (string, *apd.Decimal) {
	cur := ""
	if len(limits) > 0 {
		cur = limits[0].Currency
	}
	kurs := apd.New(1, 0)
	var kal Kalkulator
	for _, c := range m.CurrencyList {
		if c.Currency == cur {
			kurs = kal.Teks(".Conversion", c.Conversion)
		}
	}
	return cur, kurs
}

// Mode `Calculation` CountLossAllocation_act.
const (
	HitungPersen = "pct"
	HitungJumlah = "amount"
	HitungXOLMod = "CountXOL"
)

// HitungXOL = `CountLossAllocation_act` (params Calculation, Index, Currency).
func HitungXOL(k *Konteks, h *Halaman, m MasterTreaty, calc string, idx int, cur string) error {
	var kal Kalkulator
	shareCedant := kal.BagiPega(kal.H(h, CD+"ShareCeding"), k100()) // 2
	h.Setel("ProtectTreatyName.CARI1", "0")
	loss := h.AmbilDaftar(DaftarLossAlloc)
	if calc != HitungXOLMod { // 3-5 satu baris Loss Allocation
		var claim, salv, adj, oth *apd.Decimal
		claim, salv, adj, oth = apd.New(0, 0), apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
		for _, r := range h.AmbilDaftar(DaftarClaimAmount) {
			if r["Currency"] == cur {
				claim, salv, adj, oth = kal.B(r, "ClaimAmountAdjust"), kal.B(r, "Salvage"), kal.B(r, "AdjusterFee"),
					kal.B(r, "CNPOthersFee")
			}
		}
		if idx >= 1 && idx <= len(loss) {
			l := loss[idx-1]
			p := kal.BagiPega(kal.B(l, "ClaimPercentage"), k100())
			switch calc {
			case HitungPersen: // 4
				l["ClaimAmountAdjust"] = Teks(kal.Kali(p, claim))
				l["Salvage"] = Teks(kal.Kali(kal.Kali(p, salv), shareCedant))
				l["AdjusterFee"] = Teks(kal.Kali(kal.Kali(p, adj), shareCedant))
				l["CNPOthersFee"] = Teks(kal.Kali(kal.Kali(p, oth), shareCedant))
			case HitungJumlah: // 5
				l["ClaimPercentage"] = Teks(kal.Kali(kal.BagiPega(kal.B(l, "ClaimAmountAdjust"), claim), k100()))
			}
		}
	}
	// 7-10 XOL dimulai: daftar hasil dibuang.
	h.SetelDaftar(DaftarXOL, nil)
	h.SetelDaftar(DaftarSpreading, nil)
	h.SetelDaftar(DaftarBreakQS, nil)
	limits := limitGrup(h, m) // 11-12
	noRIP := ""
	if len(m.Limits) > 0 { // `Local.NoRIPCalculation` di-set pada SETIAP Limits -> nilai Limits terakhir
		noRIP = m.Limits[len(m.Limits)-1].NoRIPCalculation
	}
	// 13 total klaim per mata uang dari Loss Allocation ber-To XOL.
	type grupMU struct {
		cur, curID            string
		value, adj, salv, oth *apd.Decimal
		kurs, prorate         *apd.Decimal
	}
	var temp []*grupMU
	for _, l := range loss {
		if l["CNPFlagXOL"] != "true" {
			continue // 13.1 F:4
		}
		g := &grupMU{cur: l["Currency"], curID: l["CurrencyID"], value: kal.B(l, "ClaimAmountAdjust"),
			adj: kal.B(l, "AdjusterFee"), salv: kal.B(l, "Salvage"), oth: kal.B(l, "CNPOthersFee"),
			kurs: apd.New(0, 0), prorate: apd.New(0, 0)}
		for _, r := range h.AmbilDaftar(DaftarClaimAmount) { // 13.2
			if r["Currency"] == g.cur {
				g.curID, g.prorate, g.kurs = r["CurrencyID"], kal.B(r, "PctProrateClaim"), kal.B(r, "AltValue")
			}
		}
		ada := false
		for _, t := range temp { // 13.3
			if t.cur == g.cur {
				t.value = kal.Tambah(t.value, g.value)
				t.adj = kal.Tambah(t.adj, g.adj)
				t.salv = kal.Tambah(t.salv, g.salv)
				t.oth = kal.Tambah(t.oth, g.oth)
				ada = true
			}
		}
		if !ada { // 13.4
			temp = append(temp, g)
		}
	}
	kursLayer, kursExc := kursLapisan(limits, m) // 14-15
	rnm := kal.H(h, TM+"RNMShare")
	pctRNM := kal.BagiPega(rnm, k100())
	var xol []Baris
	var kursMDP *apd.Decimal = apd.New(0, 0) // `Local.KursMDP` bertahan antar-iterasi
	for _, t := range temp {                 // 17
		value := t.value
		kurs := kursExc
		if kursLayer == "IDR" {
			kurs = t.kurs
		}
		claimValue := kal.Kurang(kal.Tambah(kal.Tambah(value, t.adj), t.oth), t.salv)
		for li, L := range limits { // 17.2
			pr := t.prorate
			limit, limit2 := kal.Teks(".Limit", L.Limit), kal.Teks(".Limit2", L.Limit2)
			ded := kal.Teks(".Deductible", L.Deductible)
			ded2 := kal.Teks(".Deductible2", strings.ReplaceAll(L.Deductible2, ",", "."))
			persenProrata := func(x *apd.Decimal) *apd.Decimal { return kal.BagiPega(kal.Kali(x, pr), k100()) }
			var layerLimit *apd.Decimal
			switch {
			case t.cur == "IDR":
				layerLimit = persenProrata(kal.BagiPega(limit, kurs))
			case t.cur == L.Currency:
				layerLimit = persenProrata(kal.Kali(limit, kurs))
			case limit2.IsZero():
				layerLimit = persenProrata(kal.BagiPega(limit, kurs))
			default:
				layerLimit = persenProrata(limit2)
			}
			var urLimit *apd.Decimal
			switch {
			case L.Currency == "IDR":
				urLimit = persenProrata(kal.BagiPega(ded, kurs))
			case t.cur == L.Currency:
				urLimit = persenProrata(ded)
			default:
				urLimit = persenProrata(kal.Kali(ded2, kurs))
			}
			lewatUR := Lebih(claimValue, urLimit)
			ur, totalUR := value, value
			if lewatUR {
				switch {
				case L.Currency == "IDR":
					ur = persenProrata(kal.BagiPega(ded, kurs))
				case t.cur == L.Currency:
					ur = persenProrata(ded)
				default:
					ur = persenProrata(kal.Kali(ded, kurs))
				}
				switch {
				case t.cur == "IDR":
					totalUR = persenProrata(kal.BagiPega(ded, kurs))
				case t.cur == L.Currency:
					totalUR = apd.New(0, 0) // `.Deductible*0*ProrateClaim/100` (OQ-CNP-18, ditiru)
				case ded2.IsZero():
					totalUR = persenProrata(kal.BagiPega(ded, kurs))
				default:
					totalUR = persenProrata(ded2)
				}
			}
			var limitPenuh *apd.Decimal
			switch {
			case t.cur == "IDR":
				limitPenuh = kal.BagiPega(limit, kurs)
			case t.cur == L.Currency2:
				if limit2.IsZero() {
					limitPenuh = kal.BagiPega(limit, kurs)
				} else {
					limitPenuh = limit2
				}
			default:
				limitPenuh = kal.Kali(limit2, kurs)
			}
			nama := NamaLayerXOL(L.Layer, L.LayerType)
			if !strings.Contains(nama, "ST") && !strings.Contains(nama, "ND") && !strings.Contains(nama, "RD") &&
				!strings.Contains(nama, "TH") { // 17.2.2-17.2.3
				h.Setel("ProtectTreatyName.CARI1", "1")
				h.TambahPesan("", PesanNamaTreatyXOL)
			}
			jenis, err := k.Acuan.JenisReasXOL(k.Ctxt(), nama) // 17.2.4
			if err != nil {
				return err
			}
			if li == 0 { // 17.2.5 baris UR
				xol = append(xol, Baris{"TreatyType": TreatyUR, "TreatyName": TreatyUR, "ClaimEstimation": Teks(ur),
					"ClaimAmountAdjust": Teks(ur), "ClaimPercentage": "0", "ClaimSpreaded": "0", "CurrencyID": t.curID,
					"Currency": t.cur, "Kurs": Teks(kurs), "AdjClaimValue": "0", "AdjusterFee": "0", "Salvage": "0",
					"TotalClaim": Teks(totalUR), "IsEditClaim": "0"})
			}
			totalAlokasi := apd.New(0, 0) // 17.2.6
			for _, x := range xol {
				if x["Currency"] == t.cur {
					totalAlokasi = kal.Tambah(totalAlokasi, kal.B(x, "TotalClaim"))
				}
			}
			mdp, curMDP := mdpLayer(&kal, L, t.cur) // 17.2.7-17.2.9
			if L.IsCombineMDP == "true" {           // 17.2.10-17.2.11
				for _, s := range m.LimitSummaryList {
					if awal1(L.Layer) == awal1(s.Layer) && L.LayerType == s.LayerType {
						if Lebih(kal.Teks(".MDP", s.MDP), apd.New(0, 0)) {
							mdp = kal.Tambah(mdp, kal.Teks(".MDP", s.MDP))
							curMDP = s.Currency
						} else {
							mdp = kal.Tambah(mdp, kal.Teks(".MDP2", s.MDP2))
							curMDP = s.Currency2
						}
					}
				}
			}
			for _, c := range m.CurrencyList { // 17.2.12
				if c.Currency == curMDP {
					kursMDP = kal.Teks(".Conversion", c.Conversion)
				}
			}
			if curMDP != t.cur { // 17.2.13
				adaKurs := Lebih(kursMDP, apd.New(0, 0))
				switch {
				case curMDP == "IDR" && adaKurs:
					mdp = kal.BagiPega(mdp, kursMDP)
				case curMDP == "IDR":
					mdp = kal.BagiPega(mdp, kurs)
				case adaKurs:
					mdp = kal.Kali(mdp, kursMDP)
				default:
					mdp = kal.Kali(mdp, kurs)
				}
			}
			if Banding(totalAlokasi, claimValue) == 0 { // 17.2.14 T:1 -> TO: layer berikut dilewati
				break
			}
			est := kal.Kurang(claimValue, totalAlokasi)
			if !Lebih(layerLimit, est) {
				est = layerLimit
			}
			penuh := Banding(limitPenuh, est) == 0
			feeLayer := func(totalFee *apd.Decimal) *apd.Decimal {
				switch {
				case penuh:
					return apd.New(0, 0)
				case value.IsZero() && !totalFee.IsZero():
					return kal.BagiPega(kal.Kali(est, rnm), k100())
				default:
					return kal.Kali(totalFee, pctRNM)
				}
			}
			adjF, salv, oth := feeLayer(t.adj), feeLayer(t.salv), feeLayer(t.oth)
			var spread *apd.Decimal
			switch {
			case value.IsZero():
				spread = apd.New(0, 0)
			case Banding(est, layerLimit) == 0:
				spread = kal.Kurang(kal.BagiPega(kal.Kali(est, rnm), k100()), kal.Tambah(kal.Kurang(adjF, salv), oth))
			default:
				spread = kal.BagiPega(kal.Kali(est, rnm), k100())
			}
			mdpTulis := mdp
			if noRIP == "true" {
				mdpTulis = apd.New(0, 0)
			}
			idr := est
			if t.cur != "IDR" {
				idr = kal.Kali(est, kurs)
			}
			b := Baris{"TreatyName": "", "TreatyType": "", "ClaimEstimation": Teks(est), "TotalClaim": Teks(est),
				"ClaimPercentage": h.Ambil(TM + "RNMShare"), "CurrencyID": t.curID, "Currency": t.cur, "Kurs": Teks(kurs),
				"AdjClaimValue": "0", "AdjusterFee": Teks(adjF), "Salvage": Teks(salv), "CNPOthersFee": Teks(oth),
				"ClaimSpreaded": Teks(spread), "CNPLimit": Teks(limitPenuh), "CNPMDP": Teks(mdpTulis),
				"CNPPctReinstate": L.ReinstatementPct, "CNPProrateClaim": Teks(pr), "ClaimAmountIDR": Teks(idr),
				"Layer": L.Layer, "LayerType": L.LayerType, "LayerPart": L.LayerPart, "LayerPartType": L.LayerPartType,
				"KursIDR": Teks(t.kurs), "ClaimAmountAdjust": Teks(est),
				"TotalSpread":  Teks(kal.Tambah(kal.Kurang(kal.Tambah(spread, adjF), salv), oth)),
				"CNPLimitFull": L.Limit, "CurrLayerOri": L.Currency, "IsEditClaim": "0"}
			if len(jenis) > 0 {
				b["TreatyType"], b["TreatyName"] = jenis[0].ID, jenis[0].Note
			}
			xol = append(xol, b)
		}
	}
	if err := kal.Galat(); err != nil {
		return err
	}
	h.SetelDaftar(DaftarXOL, xol)
	return SusunSpreadingXOL(h, m) // 19, 21 (22 CountReinstatement_Act: grid NEVER, OQ-CNP-12)
}

// NamaLayerXOL = CountLossAllocation_act langkah 17.2.1 `InputSpreading.CARI1`: "XL " + digit pertama Layer + ST/ND/RD/TH
// + " " + LayerType huruf besar ("SUBLAYER" -> "SUB LAYER").
func NamaLayerXOL(layer, layerType string) string {
	d := awal1(layer)
	akhiran := "TH "
	switch d {
	case "1":
		akhiran = "ST "
	case "2":
		akhiran = "ND "
	case "3":
		akhiran = "RD "
	}
	tipe := strings.ToUpper(layerType)
	if tipe == "SUBLAYER" {
		tipe = "SUB LAYER"
	}
	return "XL " + d + akhiran + tipe
}

// awal1 = `@substring(x,0,1)`.
func awal1(s string) string {
	if s == "" {
		return ""
	}
	return string([]rune(s)[:1])
}

// mdpLayer = CountLossAllocation_act langkah 17.2.7-17.2.9: MDP layer menurut mata uang klaim (baris terakhir yang
// cocok); MDP 0 / tidak ada -> MDPList(1); IsCombineMDP -> 0 (dijumlah ulang dari LimitSummaryList).
func mdpLayer(kal *Kalkulator, L LimitXOL, cur string) (*apd.Decimal, string) {
	mdp, curMDP := apd.New(0, 0), ""
	for _, x := range L.MDPList {
		if x.Currency == cur {
			mdp, curMDP = kal.Teks(".Value", x.Value), x.Currency
		}
	}
	if !Lebih(mdp, apd.New(0, 0)) && len(L.MDPList) > 0 {
		mdp, curMDP = kal.Teks(".Value", L.MDPList[0].Value), L.MDPList[0].Currency
	}
	if L.IsCombineMDP == "true" {
		mdp = apd.New(0, 0)
	}
	return mdp, curMDP
}

// SusunSummaryXOL = CountLossAllocation_act langkah 19 (Summary XOL Allocation, `ListTotalEstimation`): per mata uang
// SpreadingRisk - IDR = Σ(ClaimEstimation + AdjClaimValue), Value = Σ ClaimSpreaded, fee, ClaimAmount = Σ ClaimAmountIDR,
// CNPTotalClaim = Σ TotalClaim; AltValue = Kurs baris pertama mata uang itu, PctProrateClaim = baris terakhirnya.
// Medan turunan: dihitung ulang dari SpreadingRisk tersimpan setiap halaman dimuat.
func SusunSummaryXOL(h *Halaman) error {
	var kal Kalkulator
	var out []Baris
	for _, x := range h.AmbilDaftar(DaftarXOL) {
		var s Baris
		for _, o := range out {
			if o["Currency"] == x["Currency"] {
				s = o
			}
		}
		if s == nil {
			s = Baris{"Currency": x["Currency"], "CurrencyID": x["CurrencyID"], "IDR": "0", "Value": "0",
				"AdjusterFee": "0", "Salvage": "0", "CNPOthersFee": "0", "ClaimAmount": "0", "CNPTotalClaim": "0",
				"AltValue": x["Kurs"]}
			out = append(out, s)
		}
		tambah := func(p string, v *apd.Decimal) { s[p] = Teks(kal.Tambah(kal.B(s, p), v)) }
		tambah("IDR", kal.Tambah(kal.B(x, "ClaimEstimation"), kal.B(x, "AdjClaimValue")))
		tambah("Value", kal.B(x, "ClaimSpreaded"))
		tambah("AdjusterFee", kal.B(x, "AdjusterFee"))
		tambah("Salvage", kal.B(x, "Salvage"))
		tambah("CNPOthersFee", kal.B(x, "CNPOthersFee"))
		tambah("ClaimAmount", kal.B(x, "ClaimAmountIDR"))
		tambah("CNPTotalClaim", kal.B(x, "TotalClaim"))
		s["PctProrateClaim"] = x["CNPProrateClaim"]
	}
	if err := kal.Galat(); err != nil {
		return err
	}
	h.SetelDaftar(DaftarSummaryXOL, out)
	return nil
}

// SusunSpreadingXOL = CountLossAllocation_act langkah 19 + 21: Summary XOL, lalu Spreading List (`SpreadingClaim`,
// Share(1) XoL master) dan Break QS (`SpreadingBreakQS`, `Share(1).SpreadingListXOL`) per mata uang.
func SusunSpreadingXOL(h *Halaman, m MasterTreaty) error {
	if err := SusunSummaryXOL(h); err != nil {
		return err
	}
	var kal Kalkulator
	var sh ShareXOL
	if len(m.Share) > 0 {
		sh = m.Share[0]
	}
	pct := kal.BagiPega(kal.Teks(".SpreadingTotalPctXOL", sh.SpreadingTotalPctXOL), k100())
	var spr, qs []Baris
	for _, s := range h.AmbilDaftar(DaftarSummaryXOL) {
		claimValue := kal.Kali(kal.B(s, "Value"), pct)
		adj, salv := kal.Kali(kal.B(s, "AdjusterFee"), pct), kal.Kali(kal.B(s, "Salvage"), pct)
		oth, idr := kal.Kali(kal.B(s, "CNPOthersFee"), pct), kal.Kali(kal.B(s, "ClaimAmount"), pct)
		spr = append(spr, Baris{"TreatyType": sh.SpreadingTypeIDXOL, "TreatyName": sh.SpreadingTypeXOL,
			"SharePercentage": sh.SpreadingTotalPctXOL, "Currency": s["Currency"], "CurrencyID": s["CurrencyID"],
			"ClaimSpreaded": Teks(claimValue), "AdjusterFee": Teks(adj), "Salvage": Teks(salv), "CNPOthersFee": Teks(oth),
			"ClaimAmountIDR": Teks(idr)})
		for _, q := range sh.SpreadingListXOL { // 21.2
			tt := q.ReinsTypeID
			if tt == "" {
				switch q.ReinsTypeName {
				case "QS (OR)":
					tt = "10028"
				case "QS (R/I)":
					tt = "10004"
				}
			}
			p := kal.Teks(".Pct", q.Pct)
			bagi := func(x *apd.Decimal) string { return Teks(kal.BagiPega(kal.Kali(x, p), k100())) }
			qs = append(qs, Baris{"TreatyType": tt, "TreatyName": q.ReinsTypeName, "Currency": s["Currency"],
				"CurrencyID": s["CurrencyID"], "SharePercentage": q.Pct, "ClaimSpreaded": bagi(claimValue),
				"AdjusterFee": bagi(adj), "Salvage": bagi(salv), "CNPOthersFee": bagi(oth), "ClaimAmountIDR": bagi(idr)})
		}
	}
	if err := kal.Galat(); err != nil {
		return err
	}
	h.SetelDaftar(DaftarSpreading, spr)
	h.SetelDaftar(DaftarBreakQS, qs)
	return nil
}

// SebarUlangSpreading = AdjClaimAmount_Act langkah 14 / SetActualPremium_ACT langkah 5: fee dan Claim Spreaded baris
// Spreading List dan Break QS dihitung ulang dari Summary XOL per mata uang (`ClaimAmountIDR` tidak disentuh).
func SebarUlangSpreading(h *Halaman) error {
	if err := SusunSummaryXOL(h); err != nil {
		return err
	}
	var kal Kalkulator
	for _, s := range h.AmbilDaftar(DaftarSummaryXOL) {
		for _, d := range []string{DaftarSpreading, DaftarBreakQS} {
			for _, b := range h.AmbilDaftar(d) {
				if b["Currency"] != s["Currency"] {
					continue
				}
				p := kal.BagiPega(kal.B(b, "SharePercentage"), k100())
				b["AdjusterFee"] = Teks(kal.Kali(kal.B(s, "AdjusterFee"), p))
				b["ClaimSpreaded"] = Teks(kal.Kali(kal.B(s, "Value"), p))
				b["CNPOthersFee"] = Teks(kal.Kali(kal.B(s, "CNPOthersFee"), p))
				b["Salvage"] = Teks(kal.Kali(kal.B(s, "Salvage"), p))
			}
		}
	}
	return kal.Galat()
}

// bagiPersen - `@divide(x*100, pct, 20)` (pembagi nol = 0).
func bagiPersen(kal *Kalkulator, x, pct *apd.Decimal) *apd.Decimal {
	return kal.BagiPega(kal.Kali(x, k100()), pct)
}

// AdjClaimAmount = `AdjClaimAmount_Act` (params idx, Curr): Claim Amount satu layer XOL Allocation diubah; layer
// terakhir mata uang itu menyerap selisih terhadap total Loss Allocation.
func AdjClaimAmount(h *Halaman, idx int) error {
	xol := h.AmbilDaftar(DaftarXOL)
	if idx < 1 || idx > len(xol) {
		return ErrBarisTidakAda
	}
	var kal Kalkulator
	r := xol[idx-1]
	pct := kal.B(r, "ClaimPercentage")
	// 2
	var adjV *apd.Decimal
	if r["TreatyType"] == TreatyUR {
		adjV = kal.Kurang(kal.B(r, "TotalClaim"), kal.B(r, "ClaimEstimation"))
	} else {
		adjV = kal.Kurang(kal.B(r, "TotalClaim"), kal.Kurang(kal.Tambah(kal.Tambah(kal.B(r, "ClaimEstimation"),
			bagiPersen(&kal, kal.B(r, "AdjusterFee"), pct)), bagiPersen(&kal, kal.B(r, "CNPOthersFee"), pct)),
			bagiPersen(&kal, kal.B(r, "Salvage"), pct)))
	}
	r["AdjClaimValue"] = Teks(adjV)
	r["ClaimAmountAdjust"] = Teks(kal.Tambah(adjV, kal.B(r, "ClaimEstimation")))
	// 3
	r["ClaimSpreaded"] = Teks(kal.BagiPega(kal.Kali(kal.Tambah(kal.B(r, "ClaimEstimation"), adjV), pct), k100()))
	cur, nilaiLama, adjustment := r["Currency"], r["ClaimEstimation"], adjV
	totalClaim, totalFee, claimAdjust := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
	for _, l := range h.AmbilDaftar(DaftarLossAlloc) { // 4
		if l["Currency"] != cur {
			continue
		}
		claimAdjust = kal.B(l, "ClaimAmountAdjust")
		totalClaim = kal.Tambah(totalClaim, claimAdjust)
		totalFee = kal.Tambah(totalFee, kal.Kurang(kal.Tambah(kal.B(l, "AdjusterFee"), kal.B(l, "CNPOthersFee")),
			kal.B(l, "Salvage")))
	}
	idxAkhir, totalAl := 0, apd.New(0, 0)
	for i, s := range xol { // 5
		if s["Currency"] == cur && SamaAngka(s["ClaimEstimation"], nilaiLama) {
			s["AdjClaimValue"] = Teks(adjustment)
			jumlah := kal.Tambah(kal.B(s, "ClaimEstimation"), adjustment)
			s["ClaimSpreaded"] = Teks(kal.BagiPega(kal.Kali(jumlah, kal.B(s, "ClaimPercentage")), k100()))
			s["ClaimAmountAdjust"] = Teks(jumlah)
		}
		if s["Currency"] == cur {
			idxAkhir = i + 1
			totalAl = kal.Tambah(totalAl, kal.Tambah(kal.B(s, "ClaimEstimation"), kal.B(s, "AdjClaimValue")))
		}
	}
	if idxAkhir == 0 {
		return kal.Galat()
	}
	L := xol[idxAkhir-1] // 6
	pctL := kal.B(L, "ClaimPercentage")
	L["AdjClaimValue"] = Teks(kal.Tambah(kal.B(L, "AdjClaimValue"), kal.Kurang(totalClaim, totalAl)))
	jumlahL := kal.Tambah(kal.B(L, "ClaimEstimation"), kal.B(L, "AdjClaimValue"))
	L["ClaimSpreaded"] = Teks(kal.BagiPega(kal.Kali(jumlahL, pctL), k100()))
	L["ClaimAmountAdjust"] = Teks(jumlahL)
	if !claimAdjust.IsZero() { // 7
		if L["TreatyType"] == TreatyUR {
			L["TotalClaim"] = L["ClaimAmountAdjust"]
		} else {
			L["TotalClaim"] = Teks(kal.Kurang(kal.Tambah(kal.Tambah(jumlahL, bagiPersen(&kal, kal.B(L, "AdjusterFee"), pctL)),
				bagiPersen(&kal, kal.B(L, "CNPOthersFee"), pctL)), bagiPersen(&kal, kal.B(L, "Salvage"), pctL)))
		}
	} else { // 8
		if L["TreatyType"] == TreatyUR {
			L["TotalClaim"] = L["ClaimAmountAdjust"]
		} else {
			L["TotalClaim"] = Teks(totalFee)
		}
		L["CNPOthersFee"] = Teks(kal.BagiPega(kal.Kali(kal.B(L, "TotalClaim"), pctL), k100()))
	}
	if err := kal.Galat(); err != nil {
		return err
	}
	return SebarUlangSpreading(h) // 9-14
}

// PremiAktual - nilai Claim Spreaded layer terakhir saat Waiting For Actual Premium (SetActualPremium_ACT langkah 2,
// tertulis mati 50.000; bawaan OQ-CNP-17 konstanta bernama).
const PremiAktual = "50000"

// SetActualPremium = `SetActualPremium_ACT`: FlagActualPremium "false" -> keluar (1). Untuk SETIAP baris non-UR, baris
// TERAKHIR SpreadingRisk ditimpa (langkah 3.1 menulis `SpreadingRisk(<LAST>)`, persen dari baris yang diiterasi):
// ClaimSpreaded = 50.000, TotalClaim = 50.000 / (persen / 100), lalu Summary XOL dan spreading dihitung ulang (4-5).
func SetActualPremium(h *Halaman) error {
	if h.Ambil("FlagActualPremium") == "false" {
		return nil
	}
	var kal Kalkulator
	xol := h.AmbilDaftar(DaftarXOL)
	if len(xol) == 0 {
		return nil
	}
	akhir := xol[len(xol)-1]
	premi := kal.D(PremiAktual)
	for _, x := range xol {
		if x["TreatyName"] == TreatyUR {
			continue
		}
		total := kal.BagiPega(premi, kal.BagiPega(kal.B(x, "ClaimPercentage"), k100()))
		akhir["ClaimSpreaded"] = PremiAktual
		akhir["TotalClaim"] = Teks(total)
		akhir["ClaimAmountAdjust"] = akhir["TotalClaim"]
		akhir["ClaimAmountIDR"] = akhir["TotalClaim"]
		akhir["ClaimEstimation"] = akhir["TotalClaim"]
		akhir["TotalSpread"] = akhir["ClaimSpreaded"]
	}
	if err := kal.Galat(); err != nil {
		return err
	}
	return SebarUlangSpreading(h)
}

// ---------------------------------------------------------------- pembantu

func k100() *apd.Decimal { return apd.New(100, 0) }

// maks = `@if(a>b, a, b)`.
func maks(a, b *apd.Decimal) *apd.Decimal {
	if Lebih(a, b) {
		return a
	}
	return b
}

// maks2Lebih = `@if(a>b, a, b)` untuk TPL (CountClaimTNP_Act 13.2.3).
func maks2Lebih(a, b *apd.Decimal) *apd.Decimal { return maks(a, b) }

// SamaAngka - dua teks bernilai sama sebagai bilangan (`==` Pega atas Decimal); teks bukan angka dibandingkan sebagai
// teks.
func SamaAngka(a, b string) bool {
	var kal Kalkulator
	x, y := kal.Teks("a", a), kal.Teks("b", b)
	if kal.Galat() != nil {
		return a == b
	}
	return Banding(x, y) == 0
}
