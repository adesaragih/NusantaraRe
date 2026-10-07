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
// polis -> T_GENERAL_POLIS_TREATY. Nol tabel baru (bab 0 butir 11).
//
// Semantik Pega yang ditiru:
//   - variabel lokal ber-tipe Decimal (pyLocalParameters) bernilai awal 0;
//   - `@toDecimal("")` dan properti angka kosong = 0 di aritmetika (`AngkaTeks`);
//   - `Page-New InputXOL` MENGGANTI halaman - nilai CARI sebelumnya hilang;
//   - langkah berlabel `//` dinonaktifkan (tidak diport);
//   - langkah perulangan bersyarat `.Currency==...` dievaluasi per baris.
//
// Keanehan rumus yang ditiru apa adanya (`⚠️ Ditiru apa adanya` di bawah dan di
// nonprop_detail.go) - `[keputusan work owner]` F5 04-10-2026 - didaftar beserta
// langkah XML-nya di docs/PERMINTAAN-TIM-INTI.md bagian G (G1-G17) untuk
// ditinjau Finance/Product; diubah hanya lewat keputusan tertulis per butir.

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

// angkaMaster membaca `TreatyIn.<m>` sebagai angka (kosong = 0).
func angkaMaster(k *kalkulator, h *Halaman, m string) *apd.Decimal {
	d, err := h.Angka(jMaster + m)
	if err != nil {
		k.catat(err)
		return apd.New(0, 0)
	}
	return d
}
