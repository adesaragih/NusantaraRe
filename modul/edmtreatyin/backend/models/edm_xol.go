package models

// Untuk apa berkas ini: DAFTAR XOL ENDORSEMEN dari master - port korpus `EDM Treaty In\Activity`:
//
//	InsertToTreatyXOLListEDM          (Work, 01-01-77; 36 langkah aktif, nol `//`)  -> InsertToTreatyXOLListEDM
//	InsertToTreatyXOLListEDMOldData   (Work, 01-01-77; 35 aktif, nol `//`)          -> InsertToTreatyXOLListEDMOldData
//
// Versi treaty keluar (XOL Retro) di edm_xol_keluar.go; selisih XOL di edm_xol_selisih.go. Pemanggil:
// `InputPolicyTreatyEDMDetail_NP_AdjPremi` 8 dan `TreatyRealizationCheckXOLListEDM` 6 (edm_pilih.go).
//
// Semantik Pega yang ditiru sama dengan nonprop.go (port NB `InsertToTreatyXOLList`): `Page-New InputXOL`
// mengganti halaman CARI; lokal Decimal bernilai awal 0; properti angka kosong = 0 di aritmetika (`AngkaTeks`);
// perulangan bersyarat `.Currency==local.tempcurrency` per baris; ID mata uang hasil RDB `GetDataCurrencyByName_SQL`
// (`CurrencySearch.pxResults(1).ID`) datang sebagai `IDMataUang` - fungsi di sini murni.

import "github.com/cockroachdb/apd/v3"

// xeKonfig - beda dua aktivitas yang rumusnya sama (lihat tabel di `InsertToTreatyXOLListEDMOldData`).
type xeKonfig struct {
	sumber     *Halaman // halaman tempat daftar angsuran dan layer master dibaca
	angsuran   string   // jalur daftar angsuran master (`.Installment`)
	lapisan    string   // jalur daftar layer (`.Share`)
	fakultatif string   // jalur FacultativeShareList (selalu di pyWorkPage)
	tujuan     string   // jalur daftar XOL yang ditulis
	flagPPH    bool     // syarat langkah pajak induk dan layer
	ttInduk    string   // TypeTax rumus BrokerageFeeSebenarnya induk
	ttLapis    string   // TypeTax rumus CARI45 layer
	c45Mentah  bool     // CARI45 cabang bukan-Inclusive = `.Deduction` mentah (EDM) / `@toDecimal(CARI44)` (OldData)
	nolLapis   bool     // CARI44/32/47/13 = 0 setiap layer (EDM 3.3.3.1; OldData tidak punya langkah ini)
	nolPajak   bool     // lokal BrokerageFeeSebenarnya/PPH/PPN = 0 setiap angsuran (EDM 3.3.1; OldData tidak)
}

// InsertToTreatyXOLListEDM = `Activity/InsertToTreatyXOLListEDM` (dipanggil `InputPolicyTreatyEDMDetail_NP_AdjPremi`
// 8 tanpa syarat). Sumber: master BARU `pyWorkPage.TreatyIn.ActualValue.Share` / `.ActualValue.FacultativeShareList`,
// angsuran `pyWorkPage.TreatyIn.Installment`; tujuan `PolicyTreatyIn.TreatyXOLList`.
//
//	1        @SizeOfPropertyList(TreatyIn.ActualValue.Share) = 0 -> ActualValue.Share = TreatyIn.Share
//	2        Property-Remove PolicyTreatyIn.TreatyXOLList
//	3        blok "For Non Prop" - PRE `IsNewPolicyNonProp==1` TIDAK dicentang (PRE-NONAKTIF): SELALU jalan, satu kali
//	3.1-3.2  CARI1 = "DUE TO US"; CARI14 = ValueDifference.Share(1).SpreadingTypeIDXOL - keduanya dihapus `Page-New`
//	         3.3.2 sebelum dibaca (nol pembaca; CARI14 tidak diport)
//	3.3      per TreatyIn.Installment (mata uang c): induk jumlah semua layer (3.3.3), ID mata uang (3.3.4-3.3.5),
//	         baris induk (3.3.6), pajak induk bila FlagPPH (3.3.7), satu ValueList per layer (3.3.8)
//
// ⚠️ Ditiru apa adanya, walau tampak salah:
//   - DueTo induk dan layer SELALU kosong (CARI1 hilang oleh `Page-New`, sama dengan NB).
//   - CARI31 tidak dinolkan per layer (3.3.3.1 hanya CARI44/32/47/13): Currency induk = CARI31 terakhir yang terisi.
//   - FlagRetroTreaty: induk Gross/Net/DueToValue = CARI44 layer TERAKHIR (bukan jumlah), Deduction = 0 (3.3.3.4);
//     layer Gross/Net/DueToValue = CARI19 = potongan fakultatif + CARI44 (3.3.8.3.4.1), "" bila tak ada baris cocok.
//   - BrokerageFeeSebenarnya induk dihitung dari potongan KUMULATIF; nilai yang tertinggal = sesudah layer terakhir.
func InsertToTreatyXOLListEDM(h *Halaman, idMU IDMataUang) error {
	if len(mDaftar(h, "ActualValue.Share")) == 0 { // 1
		pbSalinDaftar(h, h, jMaster+"Share", jMaster+"ActualValue.Share")
	}
	tt := h.Ambil(pt + "TypeTax")
	return xeIsiXOL(h, idMU, xeKonfig{
		sumber: h, angsuran: jMaster + "Installment", lapisan: jMaster + "ActualValue.Share",
		fakultatif: jMaster + "ActualValue.FacultativeShareList", tujuan: DaftarXOL,
		flagPPH: h.Ambil(pt+"FlagPPH") == "true", ttInduk: tt, ttLapis: tt,
		c45Mentah: true, nolLapis: true, nolPajak: true,
	})
}

// InsertToTreatyXOLListEDMOldData = `Activity/InsertToTreatyXOLListEDMOldData` (dipanggil
// `TreatyRealizationCheckXOLListEDM` 6, bukan XOL Retro). Langkah 1-2.3.8.5 = 2-3.3.8.5 `InsertToTreatyXOLListEDM`
// tanpa langkah 1 (salin ActualValue), dengan beda yang ditiru apa adanya:
//
//	sumber layer/angsuran   `TreatyIn.Share` / `TreatyIn.Installment` - halaman TreatyIn TINGKAT ATAS = master LAMA
//	                        hasil SetTreatyIn_Act (`atas`); FacultativeShareList tetap `pyWorkPage.TreatyIn` (2.3.3.3)
//	tujuan                  PolicyTreatyIn.OldData.TreatyXOLList (1, 2.3.6, 2.3.8.4)
//	syarat pajak            OldData.FlagPPH=="true" (2.3.7, 2.3.8.5)
//	TypeTax induk           PolicyTreatyIn.TypeTax (BARU, 2.3.3.4); TypeTax CARI45 = OldData.TypeTax (2.3.8.3.3.1)
//	CARI45 bukan-Inclusive  @toDecimal(CARI44)
//	tanpa 2.3.3.1-nol       CARI44/32/47/13 TIDAK dinolkan per layer: nilai layer sebelumnya terbawa dan DIJUMLAH
//	                        lagi bila layer berikut tak punya baris mata uang itu
//	tanpa nol pajak         lokal BrokerageFeeSebenarnya/PPH/PPN tidak dinolkan per angsuran (2.3.1)
func InsertToTreatyXOLListEDMOldData(h, atas *Halaman, idMU IDMataUang) error {
	return xeIsiXOL(h, idMU, xeKonfig{
		sumber: atas, angsuran: jMaster + "Installment", lapisan: jMaster + "Share",
		fakultatif: jMaster + "FacultativeShareList", tujuan: od + "TreatyXOLList",
		flagPPH: h.Ambil(od+"FlagPPH") == "true", ttInduk: h.Ambil(pt + "TypeTax"), ttLapis: h.Ambil(od + "TypeTax"),
	})
}

// xeIsiXOL - badan bersama kedua aktivitas (nomor langkah: EDM / OldData).
func xeIsiXOL(h *Halaman, idMU IDMataUang, c xeKonfig) error {
	k := &kalkulator{}
	hapusDaftarBeserta(h, c.tujuan) // 2 / 1
	retro := h.Ambil(pt+"FlagRetroTreaty") == "true"
	lapisan := c.sumber.AmbilDaftar(c.lapisan)
	nol := func() *apd.Decimal { return apd.New(0, 0) }
	bfs, pph, ppn := nol(), nol(), nol()
	var induk []Baris
	for i, inst := range c.sumber.AmbilDaftar(c.angsuran) { // 3.3 / 2.3
		// 3.3.1 / 2.3.1
		mu := inst["Currency"]
		cur, dueto := "", ""
		gross, net, dtv, ded := nol(), nol(), nol(), nol()
		if c.nolPajak {
			bfs, pph, ppn = nol(), nol(), nol()
		}
		in := halamanXOL{}        // 3.3.2 / 2.3.2 Page-New (CARI1 "DUE TO US" ikut hilang)
		for si := range lapisan { // 3.3.3 / 2.3.3
			if c.nolLapis { // 3.3.3.1
				in["CARI44"], in["CARI32"], in["CARI47"], in["CARI13"] = "0", "0", "0", "0"
			}
			xeBacaLapisan(c.sumber, c.lapisan, si+1, mu, in, nil, nil) // 3.3.3.2 / 2.3.3.2
			if retro {                                                 // 3.3.3.3 / 2.3.3.3
				for _, d := range h.AmbilDaftar(JalurAnak(c.fakultatif, si+1, "DeductionList")) {
					if d["Currency"] == mu { // CARI44 = .Deduction+ @toDecimal(InputXOL.CARI44)
						in["CARI44"] = formatAngka(k.Tambah(angkaBaris(k, d, "Deduction"), in.dec(k, "CARI44")))
						in["CARI31"] = d["Currency"]
					}
				}
			}
			// 3.3.3.4 / 2.3.3.4
			cur = in["CARI31"]
			if retro {
				gross, net, dtv, ded = in.dec(k, "CARI44"), in.dec(k, "CARI44"), in.dec(k, "CARI44"), nol()
			} else {
				gross = k.Tambah(gross, in.dec(k, "CARI32"))
				net = k.Tambah(net, in.dec(k, "CARI13"))
				dtv = k.Tambah(dtv, in.dec(k, "CARI47"))
				ded = k.Tambah(ded, in.dec(k, "CARI44"))
			}
			dueto = in["CARI1"]
			bfs = BrokerageSebenarnya(k, c.ttInduk, ded)
			pph, ppn = pajakXOL(k, bfs)
		}
		// 3.3.4-3.3.6 / 2.3.4-2.3.6
		b := Baris{
			"Currency": cur, "IDCurrency": idMU[mu],
			"GrossPremi": formatAngka(gross), "NetPremi": formatAngka(net), "DueToValue": formatAngka(dtv),
			"Deduction": formatAngka(ded), "DueTo": dueto,
		}
		if c.flagPPH { // 3.3.7 / 2.3.7
			b["BrokerageFeeSebenarnya"] = formatAngka(bfs)
			b["PPHValue"] = formatAngka(pph)
			b["PPNValue"] = formatAngka(ppn)
			b["NetPremiAfterPPH"] = formatAngka(k.Tambah(net, pph))
			b["NetPremiAfterPPN"] = formatAngka(k.Tambah(net, ppn))
			b["NetPremiAfterTax"] = formatAngka(k.Tambah(k.Tambah(net, pph), ppn))
		}
		induk = append(induk, b)
		var lapis []Baris
		for si, s := range lapisan { // 3.3.8 / 2.3.8
			in = awalLapisan(s)    // 3.3.8.1-3.3.8.2 / 2.3.8.1-2.3.8.2
			c45 := func(d Baris) { // 3.3.8.3.3.1 / 2.3.8.3.3.1
				if c.c45Mentah && c.ttLapis != TypeTaxInclusive {
					in["CARI45"] = d["Deduction"] // @if(..., .Deduction)
					return
				}
				in["CARI45"] = formatAngka(BrokerageSebenarnya(k, c.ttLapis, in.dec(k, "CARI44")))
			}
			retroLapis := func() { // 3.3.8.3.4 / 2.3.8.3.4 (di dalam perulangan GrossPremiumList)
				if !retro {
					return
				}
				for _, d := range h.AmbilDaftar(JalurAnak(c.fakultatif, si+1, "DeductionList")) {
					if d["Currency"] == mu { // CARI19 = .Deduction + @toDecimal(InputXOL.CARI44)
						in["CARI19"] = formatAngka(k.Tambah(angkaBaris(k, d, "Deduction"), in.dec(k, "CARI44")))
					}
				}
			}
			xeBacaLapisan(c.sumber, c.lapisan, si+1, mu, in, c45, retroLapis) // 3.3.8.3 / 2.3.8.3
			v := barisLapisan(in, idMU[mu])                                   // 3.3.8.4 / 2.3.8.4
			if retro {
				v["GrossPremi"], v["NetPremi"], v["DueToValue"] = in["CARI19"], in["CARI19"], in["CARI19"]
			}
			if c.flagPPH { // 3.3.8.5 / 2.3.8.5
				c45d, c13 := in.dec(k, "CARI45"), in.dec(k, "CARI13")
				lPPH, lPPN := pajakXOL(k, c45d)
				v["BrokerageFeeSebenarnya"] = in["CARI45"]
				v["PPHValue"] = formatAngka(lPPH)
				v["PPNValue"] = formatAngka(lPPN)
				v["NetPremiAfterPPH"] = formatAngka(k.Tambah(c13, lPPH))
				v["NetPremiAfterPPN"] = formatAngka(k.Tambah(c13, lPPN))
				v["NetPremiAfterTax"] = formatAngka(k.Tambah(c13, k.Tambah(lPPH, lPPN)))
			}
			lapis = append(lapis, v)
		}
		h.SetelDaftar(JalurAnak(c.tujuan, i+1, AnakLayerXOL), lapis)
	}
	if k.err != nil {
		return k.err
	}
	h.SetelDaftar(c.tujuan, induk)
	return nil
}

// xeBacaLapisan = perulangan `.GrossPremiumList` satu layer beserta perulangan BERSARANG di dalamnya (net, potongan,
// lalu `setelah` - langkah tambahan di akhir badan perulangan gross): untuk SETIAP baris gross, daftar net dan
// potongan layer yang sama dibaca ulang; layer tanpa baris gross tidak pernah membaca net maupun potongannya.
// `c45` (boleh nil) dijalankan pada baris potongan yang cocok sesudah CARI44 diisi.
func xeBacaLapisan(src *Halaman, lapisan string, si int, mu string, in halamanXOL, c45 func(Baris), setelah func()) {
	for _, g := range src.AmbilDaftar(JalurAnak(lapisan, si, "GrossPremiumList")) {
		if g["Currency"] == mu { // SET CARI32 = .Value; CARI31 = .Currency
			in["CARI32"] = g["Value"]
			in["CARI31"] = g["Currency"]
		}
		for _, n := range src.AmbilDaftar(JalurAnak(lapisan, si, "NetPremiumList")) {
			if n["Currency"] == mu { // CARI47 = CARI13 = .Value; CARI31 = .Currency
				in["CARI47"] = n["Value"]
				in["CARI13"] = n["Value"]
				in["CARI31"] = n["Currency"]
			}
		}
		for _, d := range src.AmbilDaftar(JalurAnak(lapisan, si, "DeductionList")) {
			if d["Currency"] == mu { // CARI44 = .Deduction
				in["CARI44"] = d["Deduction"]
				if c45 != nil {
					c45(d)
				}
			}
		}
		if setelah != nil {
			setelah()
		}
	}
}
