package models

// Untuk apa berkas ini: AKTIVITAS INDUK jalur NonProp - `InputPolicyTreatyInDetail_NonProp`
// (preACT 16), langkah 18 `InputPolicyTreatyInDetail_preACT` (PPN/PPH per layer),
// pemulihan limit master (`TreatySetReinstatement` / `SetReinstatementPct`), dan
// pra-proses `TreatyRealizationCheckXOLList` (`InputPolicyTreatyInPre_Act` 10).
// Semantik yang ditiru: lihat kepala `nonprop.go`.

import (
	"strconv"

	"github.com/cockroachdb/apd/v3"
)

// PesanGagalXOL - VERBATIM `TreatyRealizationCheckXOLList` langkah 3 (`local.err`).
const PesanGagalXOL = "Error fetching XolList"

// teksSubstring meniru `@substring(s, awal, akhir)` (substring Java, indeks
// berbasis 0, akhir eksklusif). `[tafsiran]` indeks di luar panjang teks -> ""
// (galat ekspresi Pega tidak memberi nilai).
func teksSubstring(s string, awal, akhir int) string {
	if awal < 0 || akhir > len(s) || awal > akhir {
		return ""
	}
	return s[awal:akhir]
}

// tampilanRetro = NonProp langkah 7-9 (`FlagRetroTreaty==true`): halaman
// master ditampilkan dari bagian fakultatif-retro.
//
//	7  TotalShareNetNP = TotalFacShareDeductionNP; TotalShareGrossNP = TotalShareNetNP;
//	   LimitShareSummaryList = LimitFacShareSummaryList
//	8  setiap LimitShareSummaryList: NetPremi = Deductible, NetPremi2 = Deductible2,
//	   Deductible = 0, Deductible2 = 0
//	9  TotalSpreadedNetPremi(1).Value / TotalSpreadedNetPremiRI(1).Value dari
//	   TotalShareNetNP(1).Value x SpreadingListXOL(n).Pct / 100
//
// ⚠️ Ditiru apa adanya: `isOR = @substring(SpreadingListXOL(1).Pct,4,6)` -
// potongan teks PERSENTASE, dibandingkan "OR"; kedua cabang memakai `isOR`
// (`isRI` dihitung tetapi tidak dibaca - tidak diport).
func tampilanRetro(h *Halaman, k *kalkulator) {
	h.SetelDaftar(jMaster+"TotalShareNetNP", salinBaris(mDaftar(h, "TotalFacShareDeductionNP")))
	h.SetelDaftar(jMaster+"TotalShareGrossNP", salinBaris(mDaftar(h, "TotalShareNetNP")))
	h.SetelDaftar(jMaster+"LimitShareSummaryList", salinBaris(mDaftar(h, "LimitFacShareSummaryList")))
	for _, b := range mDaftar(h, "LimitShareSummaryList") { // 8.1
		b["NetPremi"], b["NetPremi2"] = b["Deductible"], b["Deductible2"]
		b["Deductible"], b["Deductible2"] = "0", "0"
	}
	// 9
	netValue := apd.New(0, 0)
	if tn := mDaftar(h, "TotalShareNetNP"); len(tn) > 0 {
		netValue = angkaBaris(k, tn[0], "Value")
	}
	slx := mAnak(h, "Share", 1, "SpreadingListXOL")
	pct := func(n int) *apd.Decimal {
		if n > len(slx) {
			return apd.New(0, 0)
		}
		return angkaBaris(k, slx[n-1], "Pct")
	}
	teksPct1 := ""
	if len(slx) > 0 {
		teksPct1 = slx[0]["Pct"]
	}
	isOR := teksSubstring(teksPct1, 4, 6)
	bagian := func(n int) string { return formatAngka(k.Bagi(k.Kali(netValue, pct(n)), seratus)) }
	v1, v2 := bagian(2), bagian(1)
	if isOR == "OR" {
		v1 = bagian(1)
	}
	if isOR == "R/I" {
		v2 = bagian(2)
	}
	setelBaris1(h, jMaster+"TotalSpreadedNetPremi", "Value", v1)
	setelBaris1(h, jMaster+"TotalSpreadedNetPremiRI", "Value", v2)
}

// setelBaris1 menulis `<daftar>(1).<m>` - baris 1 dibuat bila belum ada.
func setelBaris1(h *Halaman, daftar, m, v string) {
	d := h.AmbilDaftar(daftar)
	if len(d) == 0 {
		d = []Baris{{}}
		h.SetelDaftar(daftar, d)
	}
	d[0][m] = v
}

// InputDetailNonProp = `Activity/InputPolicyTreatyInDetail_NonProp`, langkah
// SESUDAH master dimuat (2-6: RDB `BrowseTreatyInJoinEDM` + `adoptJSONObject` ke
// `pyWorkPage.TreatyIn` - `TerapkanMasterXOL`).
//
//	1      Page-Clear-Messages
//	7-9    FlagRetroTreaty -> tampilan master retro (`tampilanRetro`)
//	10     TreatyMasterInEDM -> IsEDMInputOnNB = true
//	13     DueTo "1"; StartDate/EndDate = TreatyIn.Commencement/Termination
//	15-17  jumlah per mata uang polis atas TreatyIn.Share (+ potongan retro)
//	18-19  medan uang polis; ListInstallment dikosongkan
//	20     ListInstallment / InstallmentList dari TreatyIn.Installment
//	21-22  FacultativeShare != 0 -> RetroShare; == 0 -> InsertToTreatyXOLList
//	23     TreatyNonPropSetSpreading
//
// ⛔ Tidak diport, dan sebabnya:
//   - 11 `Call InputPolicyTreatyEDMDetail_NP`: syaratnya (`TreatyMasterInEDM` DAN
//     `IsEDMInputOnNB != true`) mustahil sesudah langkah 10 mengisi
//     `IsEDMInputOnNB = true` tepat ketika `TreatyMasterInEDM` - tidak terjangkau.
//   - 12 `pyExpanded` (keadaan tampilan grid).
//   - 14 `pyWorkPage.OfferFacIn.IsFacRetro`: halaman kasus FAKULTATIF, tidak ada di
//     kasus treaty dan tidak ada kolomnya di diagram.
//   - 24 berlabel `//`.
//
// ⚠️ Ditiru apa adanya: 16.2.2 menghitung PPN/PPH dengan presisi 4 dari potongan
// KUMULATIF pada setiap baris DeductionTotalList (mata uang apa pun); 19
// memakai `local.netpremi`, bukan nilai retro; 17 menambah potongan retro
// SESUDAH PPN/PPH dihitung; 20.4 "When" tidak dicentang -> setiap baris
// InstallmentList diambil, apa pun mata uangnya.
func InputDetailNonProp(h *Halaman, idMU IDMataUang) error {
	k := &kalkulator{}
	p := polis{h, k}
	h.BersihkanPesan() // 1
	retro := p.teks("FlagRetroTreaty") == "true"
	if retro { // 7-9
		tampilanRetro(h, k)
	}
	if TreatyMasterInEDM(h) { // 10
		h.Setel(pt+"IsEDMInputOnNB", "true")
	}
	// 13 "Set always DUE TO US"
	h.Setel(pt+"DueTo", "1")
	h.Setel(pt+"StartDate", h.Ambil(jMaster+"Commencement"))
	h.Setel(pt+"EndDate", h.Ambil(jMaster+"Termination"))
	// 15
	mu := p.teks("Currency")
	flagPPH := p.teks("FlagPPH") == "true"
	typeTax := p.teks("TypeTax")
	gross, ded, net, share := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
	ppn, pph := apd.New(0, 0), apd.New(0, 0)
	for i := range mDaftar(h, "Share") { // 16
		si := i + 1
		gross = jumlahMataUang(k, mAnak(h, "Share", si, "GrossPremiumList"), mu, gross) // 16.1
		for _, d := range mAnak(h, "Share", si, "DeductionTotalList") {                 // 16.2
			if d["Currency"] == mu { // 16.2.1
				ded = k.Tambah(ded, angkaBaris(k, d, "Value"))
			}
			if flagPPH { // 16.2.2
				// @if(TypeTax="Inclusive",@divide(local.deduction,@divide(102.2,100,4),4),local.deduction)
				bfs := ded
				if typeTax == TypeTaxInclusive {
					bfs = k.BagiBulat(ded, k.BagiBulat(k.d("102.2"), seratus, 4), 4)
				}
				ppn = k.Kali(bfs, k.BagiBulat(k.d("2.2"), seratus, 4))
				pph = k.Kali(bfs, k.BagiBulat(apd.New(2, 0), seratus, 4))
			}
		}
		net = jumlahMataUang(k, mAnak(h, "Share", si, "NetPremiumList"), mu, net) // 16.3
		share = jumlahMataUang(k, mAnak(h, "Share", si, "RnmLimitList"), mu, share)
	}
	if retro { // 17
		for i := range mDaftar(h, "FacultativeShareList") {
			ded = jumlahMataUang(k, mAnak(h, "FacultativeShareList", i+1, "DeductionTotalList"), mu, ded)
		}
	}
	// 18
	nilaiRetro := func(biasa *apd.Decimal) *apd.Decimal {
		if retro {
			return ded
		}
		return biasa
	}
	p.setel("PremiOgp", nilaiRetro(gross))
	p.setel("ShareValue", share)
	p.setel("Deduction1", ded)
	p.setel("Deduction2", apd.New(0, 0))
	p.setel("NetPremium", nilaiRetro(net))
	p.setel("BalanceDueTo", nilaiRetro(net))
	hapusDaftarBeserta(h, DaftarAngsuran)
	if flagPPH { // 19
		p.setel("PPNValue", ppn)
		p.setel("PPHValue", pph)
		p.setel("BalanceDueTo", k.Tambah(k.Tambah(net, ppn), pph))
		p.setel("BalanceBeforePPH", k.Tambah(net, ppn))
		p.setel("BalanceBeforeTax", net)
	}
	// 20
	var angsuran []Baris
	for i, inst := range mDaftar(h, "Installment") {
		id := idMU[inst["Currency"]]       // 20.1-20.2
		angsuran = append(angsuran, Baris{ // 20.3
			"PaymentTotal": inst["AmountTotal"], "InstallmentPercentage": inst["PctTotal"],
			"Currency": inst["Currency"], "IDCurrency": id,
		})
		var rinci []Baris
		for _, il := range mAnak(h, "Installment", i+1, "InstallmentList") { // 20.4.1
			premi := il["Amount"]
			if retro { // local.deduction*@divide(.InstallmentPct,100,10)
				premi = formatAngka(k.Kali(ded, k.BagiBulat(angkaBaris(k, il, "InstallmentPct"), seratus, 10)))
			}
			rinci = append(rinci, Baris{
				"Premium": premi, "Currency": il["Currency"], "IDCurrency": id,
				"InstallmentPercentage": il["InstallmentPct"], "InstallmentNo": il["Installment"],
				"DueDate": il["PaymentDate"],
			})
		}
		h.SetelDaftar(JalurAnak(DaftarAngsuran, i+1, AnakRinciAngsuran), rinci)
	}
	if k.err != nil {
		return k.err
	}
	h.SetelDaftar(DaftarAngsuran, angsuran)
	// 21-22
	fac := angkaMaster(k, h, "FacultativeShare")
	if k.err != nil {
		return k.err
	}
	var err error
	if fac.IsZero() {
		err = InsertToTreatyXOLList(h, idMU)
	} else {
		err = InsertToTreatyXOLListRetroShare(h, idMU)
	}
	if err != nil {
		return err
	}
	return TreatyNonPropSetSpreading(h) // 23
}

// jumlahMataUang menjumlah `.Value` baris bermata uang `mu` (perulangan
// bersyarat `.Currency==local.currency` per baris).
func jumlahMataUang(k *kalkulator, rows []Baris, mu string, awal *apd.Decimal) *apd.Decimal {
	for _, r := range rows {
		if r["Currency"] == mu {
			awal = k.Tambah(awal, angkaBaris(k, r, "Value"))
		}
	}
	return awal
}

// PPNPPHLapisanXOL = `InputPolicyTreatyInDetail_preACT` langkah 18 (FlagPPH ==
// "true"): PPN/PPH per layer `TreatyIn.LimitShareSummaryList`, total per mata
// uang di `TotalShareNetNP` / `TotalShareDeductionNP`, lalu PPN/PPh angsuran
// polis (`ListInstallment` dan `InstallmentList`).
//
// ⚠️ Ditiru apa adanya: 18.2-18.3 berjalan DI DALAM perulangan layer, memakai
// total berjalan dan `Local.curr` (IDR bila `.Limit>0`) / `Local.curr2` (USD bila
// `.Limit2>0`) baris itu; 18.3.4.2 memakai `Local.PPNins` dkk. terakhir yang
// terisi - termasuk dari mata uang lain.
func PPNPPHLapisanXOL(h *Halaman) error { return ppnPphLapisan(h, true) }

// TampilanMasterNonProp menyusun ulang tampilan halaman master sesudah master
// dimuat ulang saat berkas dibuka: NonProp 7-9 (FlagRetroTreaty) dan bagian
// `TreatyIn` langkah 18 preACT - TANPA menyentuh medan polis tersimpan.
//
// `[penyimpangan sadar]` Pega menyimpan `pyWorkPage.TreatyIn` di blob kasus; di
// sini halaman master tidak disimpan (nol tabel di diagram, K8 butir 3) dan
// dibaca ulang dari master, sehingga tampilannya dihitung dari penanda TERKINI.
func TampilanMasterNonProp(h *Halaman) error {
	k := &kalkulator{}
	if h.Ambil(pt+"FlagRetroTreaty") == "true" {
		tampilanRetro(h, k)
	}
	if k.err != nil {
		return k.err
	}
	return ppnPphLapisan(h, false)
}

func ppnPphLapisan(h *Halaman, angsuran bool) error {
	k := &kalkulator{}
	p := polis{h, k}
	if p.teks("FlagPPH") != "true" {
		return nil
	}
	tt := p.teks("TypeTax")
	nol := func() *apd.Decimal { return apd.New(0, 0) }
	tPPN, tPPN2, tPPH, tPPH2 := nol(), nol(), nol(), nol()
	tAPPN, tAPPN2, tAPPH, tAPPH2, tTax, tTax2 := nol(), nol(), nol(), nol(), nol(), nol()
	ppnIns, pphIns, ptAPPN, ptATax := nol(), nol(), nol(), nol()
	persen := func(r Baris) *apd.Decimal { return k.BagiBulat(angkaBaris(k, r, "InstallmentPercentage"), seratus, 8) }
	for _, r := range mDaftar(h, "LimitShareSummaryList") { // 18
		// 18.1
		bfs := BrokerageSebenarnya(k, tt, angkaBaris(k, r, "Deductible"))
		lPPH, lPPN := pajakXOL(k, bfs)
		r["PPNValue"] = formatAngka(lPPN)
		r["NetPremiAfterPPN"] = formatAngka(k.Tambah(angkaBaris(k, r, "NetPremi"), lPPN))
		r["PPHValue"] = formatAngka(lPPH)
		r["NetPremiAfterPPH"] = formatAngka(k.Tambah(angkaBaris(k, r, "NetPremi"), lPPH))
		bfs2 := BrokerageSebenarnya(k, tt, angkaBaris(k, r, "Deductible2"))
		lPPH2, lPPN2 := pajakXOL(k, bfs2)
		r["PPNValue2"] = formatAngka(lPPN2)
		r["NetPremiAfterPPN2"] = formatAngka(k.Tambah(angkaBaris(k, r, "NetPremi2"), lPPN2))
		r["PPHValue2"] = formatAngka(lPPH2)
		r["NetPremiAfterPPH2"] = formatAngka(k.Tambah(angkaBaris(k, r, "NetPremi2"), lPPH2))
		tPPN, tPPN2 = k.Tambah(tPPN, lPPN), k.Tambah(tPPN2, lPPN2)
		tPPH, tPPH2 = k.Tambah(tPPH, lPPH), k.Tambah(tPPH2, lPPH2)
		tAPPN = k.Tambah(tAPPN, angkaBaris(k, r, "NetPremiAfterPPN"))
		tAPPN2 = k.Tambah(tAPPN2, angkaBaris(k, r, "NetPremiAfterPPN2"))
		tAPPH = k.Tambah(tAPPH, angkaBaris(k, r, "NetPremiAfterPPH"))
		tAPPH2 = k.Tambah(tAPPH2, angkaBaris(k, r, "NetPremiAfterPPH2"))
		// Local.TotalNetPremiAfterTax+.NetPremiAfterPPN+ .PPHValue
		tTax = k.Tambah(k.Tambah(tTax, angkaBaris(k, r, "NetPremiAfterPPN")), lPPH)
		tTax2 = k.Tambah(k.Tambah(tTax2, angkaBaris(k, r, "NetPremiAfterPPN2")), lPPH2)
		curr, curr2 := "", "" // @if(.Limit>0,"IDR","") / @if(.Limit2>0,"USD","")
		if angkaBaris(k, r, "Limit").Sign() > 0 {
			curr = "IDR"
		}
		if angkaBaris(k, r, "Limit2").Sign() > 0 {
			curr2 = "USD"
		}
		for _, t := range mDaftar(h, "TotalShareNetNP") { // 18.2
			if t["Currency"] == curr {
				t["TotalNetPremiAfterPPN"], t["TotalNetPremiAfterPPH"], t["TotalNetPremiAfterTax"] = formatAngka(tAPPN), formatAngka(tAPPH), formatAngka(tTax)
			}
			if t["Currency"] == curr2 {
				t["TotalNetPremiAfterPPN"], t["TotalNetPremiAfterPPH"], t["TotalNetPremiAfterTax"] = formatAngka(tAPPN2), formatAngka(tAPPH2), formatAngka(tTax2)
			}
		}
		for _, t := range mDaftar(h, "TotalShareDeductionNP") { // 18.3
			if t["Currency"] == curr { // 18.3.1
				t["TotalPPNValue"], t["TotalPPHValue"] = formatAngka(tPPN), formatAngka(tPPH)
			}
			if t["Currency"] == curr2 { // 18.3.2
				t["TotalPPNValue"], t["TotalPPHValue"] = formatAngka(tPPN2), formatAngka(tPPH2)
			}
			cur := t["Currency"] // 18.3.3
			if !angsuran {
				continue
			}
			for ai, a := range h.AmbilDaftar(DaftarAngsuran) { // 18.3.4
				currList := a["Currency"]
				if cur == a["Currency"] { // 18.3.4.1
					idr := a["Currency"] == "IDR"
					pilih := func(x, y *apd.Decimal) *apd.Decimal {
						if idr {
							return x
						}
						return y
					}
					a["PPN"], a["PPh"] = formatAngka(pilih(tPPN, tPPN2)), formatAngka(pilih(tPPH, tPPH2))
					ppnIns, pphIns = pilih(tPPN, tPPN2), pilih(tPPH, tPPH2)
					ptAPPN = k.Tambah(angkaBaris(k, a, "PaymentTotal"), ppnIns)
					ptATax = pilih(tTax, tTax2)
					a["PaymentTotalAfterPPN"], a["PaymentTotalAfterTax"] = formatAngka(ptAPPN), formatAngka(ptATax)
				}
				for _, ri := range h.AmbilDaftar(JalurAnak(DaftarAngsuran, ai+1, AnakRinciAngsuran)) { // 18.3.4.2
					if currList == ri["Currency"] {
						f := persen(ri)
						ri["PPN"], ri["PPh"] = formatAngka(k.Kali(ppnIns, f)), formatAngka(k.Kali(pphIns, f))
						ri["PremiumAfterPPN"], ri["PremiumAfterTax"] = formatAngka(k.Kali(ptAPPN, f)), formatAngka(k.Kali(ptATax, f))
					}
				}
			}
		}
	}
	return k.err
}

// ---------------------------------------------------------------- pemulihan limit master

// SetReinstatementPct = `Activity/SetReinstatementPct` atas baris
// `TreatyIn.Limits(idx)`: `Reinstatement_List` dibangun ulang, satu baris per
// pemulihan (`.ReinstatementValue` kali).
//
//	3.1  ReinstatementPct "100", ReinstatementValue n, ReinstatementNote, Amount1/2 =
//	     Limit/Limit2, ID = pxListSubscript limit, AdditionalPct "100"
//	3.2  MDPList 1 baris: AdditionalAmount1 = MDP IDR, AdditionalAmount2 = MDP USD
//	3.3  MDPList >= 2:    AdditionalAmount1 = MDPList(1) bila IDR, Amount2 = MDPList(2) bila USD
//
// Nilai master ini TIDAK dibaca rule NB lain dan tidak disimpan (halaman master
// baca-saja; pemulihan limit milik modul `treatyin`, tabel PEMULIHAN_LIMIT).
func SetReinstatementPct(h *Halaman, idx int) error {
	limits := mDaftar(h, "Limits")
	if idx < 1 || idx > len(limits) {
		return nil
	}
	l := limits[idx-1]
	jalur := JalurAnak(jMaster+"Limits", idx, "Reinstatement_List")
	h.SetelDaftar(jalur, nil) // 1
	k := &kalkulator{}
	maks := angkaBaris(k, l, "ReinstatementValue") // 2: Local.maxidx (int)
	if k.err != nil {
		return k.err
	}
	// bilangan bulat; pecahan dibuang (`int`)
	c := k.ctx()
	c.Rounding = apd.RoundDown
	var bulat apd.Decimal
	if _, err := c.RoundToIntegralValue(&bulat, maks); err != nil {
		return err
	}
	n, err := bulat.Int64()
	if err != nil {
		return err
	}
	mdp := mAnak(h, "Limits", idx, "MDPList")
	nilaiMU := func(b Baris, mu string) string {
		if b["Currency"] == mu {
			return b["Value"]
		}
		return "0"
	}
	var baris []Baris
	for r := int64(0); r < n; r++ { // 3
		b := Baris{ // 3.1
			"ReinstatementPct": "100", "ReinstatementValue": strconv.FormatInt(r+1, 10),
			"ReinstatementNote": l["ReinstatementNote"], "ReinstatementAmount1": l["Limit"],
			"ReinstatementAmount2": l["Limit2"], "ID": strconv.Itoa(idx), "AdditionalPct": "100",
		}
		if len(mdp) == 1 { // 3.2
			b["AdditionalAmount1"], b["AdditionalAmount2"] = nilaiMU(mdp[0], "IDR"), nilaiMU(mdp[0], "USD")
		}
		if len(mdp) >= 2 { // 3.3
			b["AdditionalAmount1"], b["AdditionalAmount2"] = nilaiMU(mdp[0], "IDR"), nilaiMU(mdp[1], "USD")
		}
		baris = append(baris, b)
	}
	h.SetelDaftar(jalur, baris)
	return nil
}

// TreatySetReinstatement = `Activity/TreatySetReinstatement` (SetTreatyIn_Act 13):
// setiap `TreatyIn.Limits` yang `Reinstatement_List(1).ReinstatementValue` kosong
// -> `SetReinstatementPct`.
func TreatySetReinstatement(h *Halaman) error {
	for i := range mDaftar(h, "Limits") {
		rl := mAnak(h, "Limits", i+1, "Reinstatement_List")
		if len(rl) == 0 || rl[0]["ReinstatementValue"] == "" {
			if err := SetReinstatementPct(h, i+1); err != nil {
				return err
			}
		}
	}
	return nil
}

// ---------------------------------------------------------------- pra-proses

// PerluCekDaftarXOL = syarat `InputPolicyTreatyInPre_Act` langkah 10:
// `IsNewPolicyNonProp==1` dan `.PolicyTreatyIn.EDMType` bukan "3".
//
// `[penyimpangan sadar]` ditambah: jenis proporsi `QuotationData` bukan
// "Proportional". Penanda IsNewPolicyNonProp tidak pernah diturunkan kembali
// oleh pra-proses (InputPolicyTreatyIn_preDT 4-5), sehingga berkas yang dipilih
// ulang ke kontrak proporsional akan dibuatkan baris XOL - yang ditolak
// penyimpanan (spec-penyimpanan AC 33, `PeriksaBentukSimpan`).
func PerluCekDaftarXOL(h *Halaman) bool {
	return samaDenganSatu(h.Ambil(pt+"IsNewPolicyNonProp")) &&
		h.Ambil(pt+"EDMType") != "3" &&
		h.Ambil(pt+"QuotationData.ProportionalType") != JenisProporsional
}

// AwalCekDaftarXOL = `TreatyRealizationCheckXOLList` langkah 1-2: pesan
// halaman dibersihkan; true = `TreatyXOLList` kosong, aktivitas berlanjut.
func AwalCekDaftarXOL(h *Halaman) bool {
	h.BersihkanPesan()
	return len(h.AmbilDaftar(DaftarXOL)) < 1
}

// LengkapiDaftarXOL = `TreatyRealizationCheckXOLList` langkah 4-8 sesudah master
// dimuat (`SetTreatyIn_Act` 4-5: RDB `BrowseTreatyIn` + `adoptJSONObject`; 5
// `Page-Copy` ke `pyWorkPage.TreatyIn` - `TerapkanMasterXOL`).
//
//	SetTreatyIn_Act 13  ProportionType "NonProportional" -> TreatySetReinstatement
//	6                   InsertToTreatyXOLList (BUKAN RetroShare)
//	8                   daftar tetap kosong -> pesan "Error fetching XolList"
//
// ⛔ Tidak diport: langkah 7 (`PolicyTreatyIn.OldData.TreatyXOLList`) - OldData
// adalah generasi sebelumnya, digantikan `OLD_POLIS_ID` (diagram F13-F14) yang
// di NB selalu kosong; SetTreatyIn_Act 1-3, 6-12, 14-15 (penanda tampilan dan
// komentar layar master `Data-Portal`, revisi master yang MENULIS JSON
// `SaveTreatyIn` - hanya bila `revisionstate==1`, tidak pernah dari sini;
// `CheckDuplicateOffer` langkah 1-4 `//`).
func LengkapiDaftarXOL(h *Halaman, idMU IDMataUang) error {
	if h.Ambil(jMaster+"ProportionType") == JenisNonProporsional {
		if err := TreatySetReinstatement(h); err != nil {
			return err
		}
	}
	if err := InsertToTreatyXOLList(h, idMU); err != nil {
		return err
	}
	if len(h.AmbilDaftar(DaftarXOL)) < 1 {
		h.TambahPesan("", PesanGagalXOL)
	}
	return nil
}
