package models

// Untuk apa berkas ini: JALUR NB NON-PROPORSIONAL (XOL) - port aktivitas yang
// mengisi polis XOL dari master kontrak (`[keputusan work owner]` K8):
//
//	InputPolicyTreatyInDetail_NonProp   InputDetailNonProp   (preACT 16)
//	InsertToTreatyXOLList               InsertToTreatyXOLList
//	InsertToTreatyXOLListRetroShare     InsertToTreatyXOLListRetroShare
//	TreatyNonPropSetSpreading           TreatyNonPropSetSpreading
//	TreatySetReinstatement              TreatySetReinstatement
//	SetReinstatementPct                 SetReinstatementPct
//	TreatyRealizationCheckXOLList       LengkapiDaftarXOL     (Pre_Act 10)
//	InputPolicyTreatyInDetail_preACT 18 PPNPPHLapisanXOL
//
// Setiap rumus berasal dari `PropertiesValue` satu langkah (AC 79); nomor
// langkahnya ditulis di samping baris port. Halaman master `TreatyIn` sudah
// diisi `TerapkanMasterXOL` sebelum fungsi di sini berjalan (langkah RDB-List +
// Java `adoptJSONObject` milik repository); ID mata uang hasil RDB
// `GetDataCurrencyByName_SQL` datang sebagai `IDMataUang` - fungsi di sini murni.
//
// Hasilnya tersimpan lewat katalog: `TreatyXOLList` -> T_POLIS_XOL,
// `TreatyXOLList().ValueList` -> T_POLIS_XOL_LAYER, `ListInstallment` ->
// T_POLIS_INSTALMENT, `ListInstallment().InstallmentList` ->
// T_POLIS_INSTALMENT_DETAIL, `SpreadingRiskList` -> T_POLIS_SPREADING, dan medan
// polis -> T_GENERAL_POLIS. Nol tabel baru (bab 0 butir 11).
//
// Semantik Pega yang ditiru:
//   - variabel lokal ber-tipe Decimal (pyLocalParameters) bernilai awal 0;
//   - `@toDecimal("")` dan properti angka kosong = 0 di aritmetika (`AngkaTeks`);
//   - `Page-New InputXOL` MENGGANTI halaman - nilai CARI sebelumnya hilang;
//   - langkah berlabel `//` dinonaktifkan (tidak diport);
//   - langkah perulangan bersyarat `.Currency==...` dievaluasi per baris.

import (
	"strings"

	"github.com/cockroachdb/apd/v3"
)

// Jalur daftar XOL di halaman.
const (
	// DaftarXOL - `pyWorkPage.PolicyTreatyIn.TreatyXOLList` (T_POLIS_XOL).
	DaftarXOL = HalamanPolis + ".TreatyXOLList"
	// AnakLayerXOL - `TreatyXOLList().ValueList` (T_POLIS_XOL_LAYER).
	AnakLayerXOL = "ValueList"
	// AnakRinciAngsuran - `ListInstallment().InstallmentList` (T_POLIS_INSTALMENT_DETAIL).
	AnakRinciAngsuran = "InstallmentList"
)

// teksDueToUs - VERBATIM `InsertToTreatyXOLList` langkah 2 / RetroShare 2.1.
const teksDueToUs = "DUE TO US"

// hapusDaftarBeserta meniru `Property-Remove` sebuah PageList: daftarnya DAN
// setiap daftar bersarang di bawah barisnya (`<jalur>(n).<anak>`).
func hapusDaftarBeserta(h *Halaman, jalur string) {
	h.pastikan()
	delete(h.Daftar, jalur)
	for k := range h.Daftar {
		if strings.HasPrefix(k, jalur+"(") {
			delete(h.Daftar, k)
		}
	}
}

// halamanXOL meniru halaman sementara `InputXOL` (kelas Int, medan CARIn teks).
type halamanXOL map[string]string

// dec membaca CARI sebagai `@toDecimal` (kosong = 0).
func (x halamanXOL) dec(k *kalkulator, cari string) *apd.Decimal {
	d, err := AngkaTeks("InputXOL."+cari, x[cari])
	if err != nil {
		k.catat(err)
		return apd.New(0, 0)
	}
	return d
}

// pajakXOL = `@toDecimal(x)* @divide(2,100,8)` dan `* @divide(2.2,100,8)`.
func pajakXOL(k *kalkulator, x *apd.Decimal) (pph, ppn *apd.Decimal) {
	return k.Kali(x, k.BagiBulat(apd.New(2, 0), seratus, 8)), k.Kali(x, k.BagiBulat(k.d("2.2"), seratus, 8))
}

// bacaLapisan menjalankan perulangan `.GrossPremiumList` satu layer beserta dua
// perulangan BERSARANG di dalamnya (InsertToTreatyXOLList 3.2.3 / 3.7.3,
// RetroShare 2.3.3.2 / 2.3.7.3): untuk SETIAP baris gross, daftar net dan
// potongan layer yang sama dibaca ulang. Layer tanpa baris gross tidak pernah
// membaca net maupun potongannya.
func bacaLapisan(h *Halaman, k *kalkulator, in halamanXOL, si int, mu, typeTax string, hitung45 bool) {
	for _, g := range mAnak(h, "Share", si, "GrossPremiumList") {
		if g["Currency"] == mu {
			in["CARI32"] = g["Value"]
			in["CARI31"] = g["Currency"]
		}
		for _, n := range mAnak(h, "Share", si, "NetPremiumList") {
			if n["Currency"] == mu {
				in["CARI47"] = n["Value"]
				in["CARI13"] = n["Value"]
				in["CARI31"] = n["Currency"]
			}
		}
		for _, d := range mAnak(h, "Share", si, "DeductionList") {
			if d["Currency"] == mu {
				in["CARI44"] = d["Deduction"]
				if hitung45 { // 3.7.3.3.1
					// CARI45 = @if(TypeTax="Inclusive",@divide(@toDecimal(CARI44),@divide(102.2,100,8),8),@toDecimal(CARI44))
					in["CARI45"] = formatAngka(BrokerageSebenarnya(k, typeTax, in.dec(k, "CARI44")))
				}
			}
		}
	}
}

// barisLapisan = InsertToTreatyXOLList 3.7.4 / RetroShare 2.3.7.5 (bagian yang sama).
func barisLapisan(in halamanXOL, idMU string) Baris {
	return Baris{
		"Currency": in["CARI31"], "IDCurrency": idMU,
		"GrossPremi": in["CARI32"], "NetPremi": in["CARI13"], "DueToValue": in["CARI47"],
		"Deduction": in["CARI44"],
		"LayerType": in["CARI50"], "Layer": in["CARI51"], "LayerPartType": in["CARI52"], "LayerPart": in["CARI53"],
		"DueTo": in["CARI1"],
	}
}

// awalLapisan = 3.7.1-3.7.2 / 2.3.7.1-2.3.7.2: `Page-New InputXOL`, lalu penanda layer.
func awalLapisan(s Baris) halamanXOL {
	return halamanXOL{"CARI50": s["LayerType"], "CARI51": s["Layer"], "CARI52": s["LayerPartType"], "CARI53": s["LayerPart"]}
}

// InsertToTreatyXOLList = `Activity/InsertToTreatyXOLList`: satu baris
// `TreatyXOLList` per `TreatyIn.Installment` (mata uang angsuran master),
// jumlah seluruh layer; di bawahnya satu baris `ValueList` per layer.
//
// ⚠️ Ditiru apa adanya, walau tampak salah:
//   - `DueTo` SELALU kosong: langkah 2 mengisi `InputXOL.CARI1 = "DUE TO US"`,
//     tetapi `Page-New InputXOL` (3.2.1, 3.7.1) mengganti halamannya sebelum
//     CARI1 dibaca (3.2.4, 3.7.4).
//   - `Currency` induk = CARI31 layer TERAKHIR (3.2.4 menimpa tiap layer): layer
//     terakhir tanpa baris mata uang itu menjadikannya kosong.
//   - BrokerageFeeSebenarnya/PPH/PPN induk dihitung dari potongan KUMULATIF
//     setiap layer; nilai yang tertinggal = sesudah layer terakhir.
func InsertToTreatyXOLList(h *Halaman, idMU IDMataUang) error {
	k := &kalkulator{}
	p := polis{h, k}
	hapusDaftarBeserta(h, DaftarXOL) // 1
	typeTax := p.teks("TypeTax")
	flagPPH := p.teks("FlagPPH") == "true"
	shares := mDaftar(h, "Share")
	in := halamanXOL{"CARI1": teksDueToUs} // 2
	var induk []Baris
	for i, inst := range mDaftar(h, "Installment") { // 3
		// 3.1
		mu := inst["Currency"]
		cur, dueto := "", ""
		gross, net, dtv, ded := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
		bfs, pph, ppn := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
		for si := range shares { // 3.2
			in = halamanXOL{}                               // 3.2.1 Page-New (CARI1 ikut hilang)
			bacaLapisan(h, k, in, si+1, mu, typeTax, false) // 3.2.3
			// 3.2.4
			cur = in["CARI31"]
			gross = k.Tambah(gross, in.dec(k, "CARI32"))
			net = k.Tambah(net, in.dec(k, "CARI13"))
			dtv = k.Tambah(dtv, in.dec(k, "CARI47"))
			ded = k.Tambah(ded, in.dec(k, "CARI44"))
			dueto = in["CARI1"]
			bfs = BrokerageSebenarnya(k, typeTax, ded)
			pph, ppn = pajakXOL(k, bfs)
		}
		// 3.3-3.5
		b := Baris{
			"Currency": cur, "IDCurrency": idMU[mu],
			"GrossPremi": formatAngka(gross), "NetPremi": formatAngka(net), "DueToValue": formatAngka(dtv),
			"Deduction": formatAngka(ded), "DueTo": dueto,
		}
		if flagPPH { // 3.6
			b["BrokerageFeeSebenarnya"] = formatAngka(bfs)
			b["PPHValue"] = formatAngka(pph)
			b["PPNValue"] = formatAngka(ppn)
			b["NetPremiAfterPPH"] = formatAngka(k.Tambah(net, pph))
			b["NetPremiAfterPPN"] = formatAngka(k.Tambah(net, ppn))
			b["NetPremiAfterTax"] = formatAngka(k.Tambah(k.Tambah(net, pph), ppn))
		}
		induk = append(induk, b)
		var lapis []Baris
		for si, s := range shares { // 3.7
			in = awalLapisan(s)                            // 3.7.1-3.7.2
			bacaLapisan(h, k, in, si+1, mu, typeTax, true) // 3.7.3
			v := barisLapisan(in, idMU[mu])                // 3.7.4
			if flagPPH {                                   // 3.7.5
				c45 := in.dec(k, "CARI45")
				lPPH, lPPN := pajakXOL(k, c45)
				c13 := in.dec(k, "CARI13")
				v["BrokerageFeeSebenarnya"] = in["CARI45"]
				v["PPHValue"] = formatAngka(lPPH)
				v["PPNValue"] = formatAngka(lPPN)
				v["NetPremiAfterPPH"] = formatAngka(k.Tambah(c13, lPPH))
				v["NetPremiAfterPPN"] = formatAngka(k.Tambah(c13, lPPN))
				v["NetPremiAfterTax"] = formatAngka(k.Tambah(c13, k.Tambah(lPPH, lPPN)))
			}
			lapis = append(lapis, v)
		}
		h.SetelDaftar(JalurAnak(DaftarXOL, i+1, AnakLayerXOL), lapis)
	}
	if k.err != nil {
		return k.err
	}
	h.SetelDaftar(DaftarXOL, induk)
	return nil
}

// InsertToTreatyXOLListRetroShare = `Activity/InsertToTreatyXOLListRetroShare`
// (NonProp 21, bila `TreatyIn.FacultativeShare != 0`).
//
// Langkah 1 SELALU menghapus `TreatyXOLList`; langkah 2 (seluruh pengisian)
// bersyarat `IsNewPolicyNonProp==1` - pada pilih bisnis pertama penanda itu
// masih "0" (pra-proses berjalan sebelum ProportionalType terisi), sehingga
// daftarnya kosong dan baru diisi `TreatyRealizationCheckXOLList` saat berkas
// dibuka berikutnya. Ditiru apa adanya.
//
// ⚠️ Ditiru apa adanya, walau tampak salah:
//   - `Page-New InputXOL` hanya per ANGSURAN (2.3.2), bukan per layer: nilai CARI
//     layer sebelumnya terbawa ke layer yang tak punya baris mata uang itu dan
//     DIJUMLAH lagi di 2.3.3.4.
//   - 2.3.7.4 (tanpa syarat) membaca ulang net dan potongan LAYER (bukan retro)
//     di dalam perulangan gross retro, lalu menambahkannya.
//   - CARI45 tidak pernah diisi: BrokerageFeeSebenarnya layer kosong, PPH/PPN 0
//     (2.3.7.6); induk tidak punya langkah FlagPPH sama sekali.
//   - Langkah 2.2 (`CARI14 = Share(1).SpreadingTypeIDXOL`) dihapus `Page-New`
//     2.3.2 sebelum dibaca - nol pembaca, tidak diport. 2.3.3.3.2-3 berlabel `//`.
func InsertToTreatyXOLListRetroShare(h *Halaman, idMU IDMataUang) error {
	k := &kalkulator{}
	p := polis{h, k}
	hapusDaftarBeserta(h, DaftarXOL)                   // 1
	if !samaDenganSatu(p.teks("IsNewPolicyNonProp")) { // 2
		return nil
	}
	retro := p.teks("FlagRetroTreaty") == "true"
	flagPPH := p.teks("FlagPPH") == "true"
	shares := mDaftar(h, "Share")
	in := halamanXOL{"CARI1": teksDueToUs} // 2.1
	var induk []Baris
	for i, inst := range mDaftar(h, "Installment") { // 2.3
		// 2.3.1
		mu := inst["Currency"]
		cur, dueto := "", ""
		gross, net, dtv, ded := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
		in = halamanXOL{}        // 2.3.2 Page-New (CARI1 "DUE TO US" 2.1 ikut hilang)
		for si := range shares { // 2.3.3
			bacaLapisan(h, k, in, si+1, mu, "", false) // 2.3.3.2
			if retro {                                 // 2.3.3.3
				for _, d := range mAnak(h, "FacultativeShareList", si+1, "DeductionList") {
					if d["Currency"] == mu { // CARI44 = .Deduction+ @toDecimal(CARI44)
						in["CARI44"] = formatAngka(k.Tambah(angkaBaris(k, d, "Deduction"), in.dec(k, "CARI44")))
						in["CARI31"] = d["Currency"]
					}
				}
			}
			// 2.3.3.4
			cur = in["CARI31"]
			if retro {
				gross, net, dtv, ded = in.dec(k, "CARI44"), in.dec(k, "CARI44"), in.dec(k, "CARI44"), apd.New(0, 0)
			} else {
				gross = k.Tambah(gross, in.dec(k, "CARI32"))
				net = k.Tambah(net, in.dec(k, "CARI13"))
				dtv = k.Tambah(dtv, in.dec(k, "CARI47"))
				ded = k.Tambah(ded, in.dec(k, "CARI44"))
			}
			dueto = in["CARI1"]
		}
		// 2.3.4-2.3.6
		induk = append(induk, Baris{
			"Currency": cur, "IDCurrency": idMU[mu],
			"GrossPremi": formatAngka(gross), "NetPremi": formatAngka(net), "DueToValue": formatAngka(dtv),
			"Deduction": formatAngka(ded), "DueTo": dueto,
		})
		var lapis []Baris
		for si, s := range shares { // 2.3.7
			in = awalLapisan(s)                                                            // 2.3.7.1-2.3.7.2
			bacaLapisan(h, k, in, si+1, mu, "", false)                                     // 2.3.7.3
			for _, g := range mAnak(h, "FacultativeShareList", si+1, "GrossPremiumList") { // 2.3.7.4
				if g["Currency"] == mu { // 2.3.7.4.1
					in["CARI32"] = formatAngka(k.Tambah(angkaBaris(k, g, "Value"), in.dec(k, "CARI32")))
					in["CARI31"] = g["Currency"]
				}
				for _, n := range mAnak(h, "Share", si+1, "NetPremiumList") { // 2.3.7.4.2
					if n["Currency"] == mu {
						in["CARI47"] = formatAngka(k.Tambah(angkaBaris(k, n, "Value"), in.dec(k, "CARI47")))
						in["CARI13"] = formatAngka(k.Tambah(angkaBaris(k, n, "Value"), in.dec(k, "CARI13")))
						in["CARI31"] = n["Currency"]
					}
				}
				for _, d := range mAnak(h, "Share", si+1, "DeductionList") { // 2.3.7.4.3
					if d["Currency"] == mu {
						in["CARI44"] = formatAngka(k.Tambah(angkaBaris(k, d, "Deduction"), in.dec(k, "CARI44")))
					}
				}
				if retro { // 2.3.7.4.4
					for _, d := range mAnak(h, "FacultativeShareList", si+1, "DeductionList") {
						if d["Currency"] == mu {
							in["CARI19"] = formatAngka(k.Tambah(angkaBaris(k, d, "Deduction"), in.dec(k, "CARI44")))
						}
					}
				}
			}
			v := barisLapisan(in, idMU[mu]) // 2.3.7.5
			if retro {
				v["GrossPremi"], v["NetPremi"], v["DueToValue"] = in["CARI19"], in["CARI19"], in["CARI19"]
			}
			if flagPPH { // 2.3.7.6
				lPPH, lPPN := pajakXOL(k, in.dec(k, "CARI45"))
				v["BrokerageFeeSebenarnya"] = in["CARI45"]
				v["PPHValue"] = formatAngka(lPPH)
				v["PPNValue"] = formatAngka(lPPN)
			}
			lapis = append(lapis, v)
		}
		h.SetelDaftar(JalurAnak(DaftarXOL, i+1, AnakLayerXOL), lapis)
	}
	if k.err != nil {
		return k.err
	}
	h.SetelDaftar(DaftarXOL, induk)
	return nil
}

// PesanSpreadingMaster - VERBATIM `TreatyNonPropSetSpreading` langkah 3/4.1/6
// (`local.sprderror`).
const PesanSpreadingMaster = "Spreading in Master data is incomplete"

// TreatyNonPropSetSpreading = `Activity/TreatyNonPropSetSpreading` (NonProp 23).
//
//	1    hapus SpreadingRiskList
//	3    Share(1).SpreadingTypeXOL terisi -> baris 1: SpreadingTypeIDXOL, 100%
//	4    SpreadingTypeXOL kosong -> satu baris per Share(1).SpreadingListXOL:
//	     SplitRNMSharePct = .Pct; SharePercentage = @divide(.Pct,RNMShare,20)*100
//	6    FlagRetroTreaty -> baris 1: "ORS", 100%
//	7    SpreadingTypeIDXOL kosong && FacultativeShare==0 -> pesan halaman
//	10   setiap baris -> CountSpreading_Act
//
// ⚠️ `PremiumSpreaded = local.netpremi` (3, 4.1, 6): lokal aktivitas INI yang
// tidak pernah diisi (Decimal, 0); nilainya dihitung ulang langkah 10.
// ⛔ Tidak diport: langkah 2 (`InputData.CARI44`, nol pembaca); langkah 5
// (`ClaimType=="XOL Retro"` -> `ShareReins().ReinsuranceListTONP()` - data
// TREATY KELUAR, K8 butir 4 tetap ⛔); langkah 8-9 berlabel `//` (RetroList).
// `[tafsiran]` langkah 7 dengan `local.sprderror` kosong (langkah 3/4/6 tak
// berjalan) tidak memasang pesan - pesan tanpa teks tidak membawa isi.
func TreatyNonPropSetSpreading(h *Halaman) error {
	k := &kalkulator{}
	p := polis{h, k}
	var baris []Baris // 1
	share1 := Baris{}
	if s := mDaftar(h, "Share"); len(s) > 0 {
		share1 = s[0]
	}
	sprderror := ""
	nol := "0" // local.netpremi
	satu := func(isi Baris) {
		if len(baris) == 0 {
			baris = append(baris, Baris{})
		}
		for kk, v := range isi {
			baris[0][kk] = v
		}
	}
	if share1["SpreadingTypeXOL"] != "" { // 3
		satu(Baris{"TreatyType": share1["SpreadingTypeIDXOL"], "Currency": p.teks("Currency"),
			"CurrencyID": p.teks("IDCurrency"), "PremiumSpreaded": nol, "SharePercentage": "100"})
		sprderror = PesanSpreadingMaster
	}
	if share1["SpreadingTypeXOL"] == "" { // 4
		rnm := angkaMaster(k, h, "RNMShare")
		for _, x := range mAnak(h, "Share", 1, "SpreadingListXOL") { // 4.1
			baris = append(baris, Baris{"TreatyType": x["ReinsTypeID"], "Currency": p.teks("Currency"),
				"CurrencyID": p.teks("IDCurrency"), "PremiumSpreaded": nol, "SplitRNMSharePct": x["Pct"],
				"SharePercentage": formatAngka(k.Kali(k.BagiBulat(angkaBaris(k, x, "Pct"), rnm, 20), seratus))})
			sprderror = PesanSpreadingMaster
		}
	}
	if p.teks("FlagRetroTreaty") == "true" { // 6
		satu(Baris{"TreatyType": "ORS", "Currency": p.teks("Currency"), "CurrencyID": p.teks("IDCurrency"),
			"PremiumSpreaded": nol, "SharePercentage": "100"})
		sprderror = PesanSpreadingMaster
	}
	if k.err != nil {
		return k.err
	}
	h.SetelDaftar(DaftarSpreading, baris)
	// 7
	if share1["SpreadingTypeIDXOL"] == "" && angkaMaster(k, h, "FacultativeShare").IsZero() && sprderror != "" {
		h.TambahPesan("", sprderror)
	}
	if k.err != nil {
		return k.err
	}
	for i := range baris { // 10
		if err := CountSpreading(h, i+1); err != nil {
			return err
		}
	}
	return nil
}

// angkaMaster membaca `TreatyIn.<m>` sebagai angka (kosong = 0).
func angkaMaster(k *kalkulator, h *Halaman, m string) *apd.Decimal {
	d, err := h.Angka(jMaster + m)
	if err != nil {
		k.catat(err)
		return apd.New(0, 0)
	}
	return d
}
