package models

// Untuk apa berkas ini: DAFTAR XOL dari master TREATY KELUAR (ClaimType "XOL Retro", master M_TREATY_OUT dibaca
// SetValueEDM_Act 7) - port korpus `EDM Treaty In\Activity`:
//
//	InsertToTreatyOutXOLList            (Work, 01-01-83; 36 aktif, 2 `//`: 4.3, 4.8)  -> InsertToTreatyOutXOLList
//	InsertToTreatyOutXOLListEDMOldData  (Work, 01-01-83; 35 aktif, 1 `//`: 4.7)       -> InsertToTreatyOutXOLListEDMOldData
//
// `InsertToTreatyOutXOLList` juga ada di korpus NB (identik) tetapi tidak dibangun modul NB (K8 butir 4).
// Pemanggil: `InputPolicyTreatyEDMDetail_NP` 9 dan `TreatyRealizationCheckXOLListEDM` 7. KEDUANYA membaca
// `pyWorkPage.TreatyIn` (master BARU), termasuk versi OldData.
//
// Langkah 3 kedua aktivitas MENGHAPUS baris master `ShareReins().ReinsuranceListTONP()` yang `ReinsName` berbeda dari
// `PolicyTreatyIn.SOBName` (`xeBuangReinsBerbeda`) - halaman master di halaman kerja ikut berubah.

import (
	"sort"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"
)

// xeTONP - jalur `TreatyIn.ShareReins(ri).ReinsuranceListTONP` (ri berbasis 1).
func xeTONP(ri int) string { return JalurAnak(jMaster+"ShareReins", ri, "ReinsuranceListTONP") }

// xeBuangReinsBerbeda = langkah 3 kedua aktivitas: per `pyWorkPage.TreatyIn.ShareReins` (Local.ShareIndex), per
// `.ReinsuranceListTONP` (Local.index = .pxListSubscript): PRE `.ReinsName==pyWorkPage.PolicyTreatyIn.SOBName` benar ->
// lewati (3), salah -> `Page-Remove ShareReins(Local.ShareIndex).ReinsuranceListTONP(Local.index)`.
//
// ⚠️ `[dugaan]` Perulangan berjalan menurut SUBSKRIP atas daftar yang sedang dihapus: sesudah baris ke-n dibuang,
// baris berikutnya bergeser ke n dan perulangan maju ke n+1 - baris yang bergeser itu TIDAK diperiksa. Dua baris
// berbeda berurutan -> baris kedua bertahan. Perilaku iterator Pega atas PageList yang diubah tidak ada di korpus.
func xeBuangReinsBerbeda(h *Halaman, sob string) {
	for ri := range mDaftar(h, "ShareReins") {
		jalur := xeTONP(ri + 1)
		for idx := 1; idx <= len(h.AmbilDaftar(jalur)); idx++ { // 3.1
			if h.AmbilDaftar(jalur)[idx-1]["ReinsName"] == sob { // 3.1.2 PRE benar -> lewati
				continue
			}
			xeHapusBaris(h, jalur, idx) // 3.1.2 Page-Remove
		}
	}
}

// xeHapusBaris meniru `Page-Remove <daftar>(n)`: baris n dibuang beserta daftar bersarangnya, daftar bersarang baris
// sesudahnya digeser satu subskrip ke bawah.
func xeHapusBaris(h *Halaman, jalur string, n int) {
	d := h.AmbilDaftar(jalur)
	if n < 1 || n > len(d) {
		return
	}
	baru := append(append([]Baris{}, d[:n-1]...), d[n:]...)
	awal := jalur + "("
	geser := map[string][]Baris{}
	var kunci []string
	for k := range h.Daftar {
		if strings.HasPrefix(k, awal) {
			kunci = append(kunci, k)
		}
	}
	sort.Strings(kunci)
	for _, k := range kunci {
		sisa := k[len(awal):]
		tutup := strings.Index(sisa, ")")
		if tutup < 0 {
			continue
		}
		i, err := strconv.Atoi(sisa[:tutup])
		if err != nil {
			continue
		}
		v := h.Daftar[k]
		delete(h.Daftar, k)
		switch {
		case i < n:
			geser[k] = v
		case i > n:
			geser[awal+strconv.Itoa(i-1)+sisa[tutup:]] = v
		}
	}
	for k, v := range geser {
		h.Daftar[k] = v
	}
	h.SetelDaftar(jalur, baru)
}

// InsertToTreatyOutXOLList = `Activity/InsertToTreatyOutXOLList` (dipanggil `InputPolicyTreatyEDMDetail_NP` 9,
// `ClaimType=="XOL Retro"`).
//
//	1        Property-Remove PolicyTreatyIn.TreatyXOLList
//	2        CARI1 = "DUE TO US" - dihapus `Page-New` 4.2.2.1 / 4.7.2.1 sebelum dibaca (DueTo selalu kosong)
//	3        buang baris ReinsuranceListTONP ber-ReinsName lain (`xeBuangReinsBerbeda`)
//	4        per TreatyIn.Installment (mata uang c):
//	4.1      lokal dinolkan (BrokerageFeeSebenarnya TIDAK - nol pembaca)
//	4.2      per ShareReins(Local.Index) per ReinsuranceListTONP(Local.ShareIndex): Page-New; GrossPremiumList baris
//	         ini -> CARI32/31; di dalamnya NetPremiumList -> CARI47/13/31 dan DeductionList -> CARI44 (keduanya
//	         `ShareReins(Local.Index).ReinsuranceListTONP(Local.ShareIndex)`)
//	4.2.2.4  jumlah Gross/Net/DueToValue/Deduction; BrokerageFeeSebenarnya/PPH/PPN dihitung ke lokal tetapi TIDAK
//	         pernah ditulis (nol pembaca - tidak diport)
//	4.4-4.6  ID mata uang; baris induk TreatyXOLList(<CURRENT>)
//	4.7      per ShareReins/ReinsuranceListTONP: satu ValueList (Layer* baris TONP, CARI31/32/13/47/44, DueTo CARI1);
//	         CARI45 4.7.2.3.3.1 dihitung tetapi tidak pernah ditulis (tidak diport)
//
// ⛔ 4.3 dan 4.8 berlabel `//` (layer dari SpreadingTONP) - tidak diport. 4.9 langkah tanpa metode.
// ⚠️ Ditiru apa adanya: 4.2.2.4 menambah CARI13 ke `local.netpremi` DUA KALI (`local.netpremi = local.netpremi +
// @toDecimal(InputXOL.CARI13)` tertulis di awal dan di akhir langkah) - NetPremi induk = 2 x jumlah net layer.
// Tanpa langkah FlagPPH: kolom pajak induk dan layer kosong.
func InsertToTreatyOutXOLList(h *Halaman, idMU IDMataUang) error {
	return xeIsiXOLKeluar(h, idMU, DaftarXOL, false)
}

// InsertToTreatyOutXOLListEDMOldData = `Activity/InsertToTreatyOutXOLListEDMOldData` (dipanggil
// `TreatyRealizationCheckXOLListEDM` 7, `ClaimType=="XOL Retro"`). Sama dengan `InsertToTreatyOutXOLList`, kecuali:
//
//	tujuan              PolicyTreatyIn.OldData.TreatyXOLList (1, 4.5, 4.6.2.4)
//	potongan            `pyWorkPage.TreatyIn.SpreadingTONP(1).ReinsuranceListTONP(1).LayerList(Local.ShareIndex)
//	                    .DeductionList` (4.2.2.3.3, 4.6.2.3.3) - Local.ShareIndex = subskrip baris TONP
//	net induk           ditambah SEKALI (4.2.2.4)
//
// ⛔ 4.7 berlabel `//` (beserta 10 anak) - tidak diport. Sumbernya master BARU `pyWorkPage.TreatyIn`, bukan
// halaman `TreatyIn` tingkat atas hasil SetTreatyIn_Act - ditiru apa adanya.
func InsertToTreatyOutXOLListEDMOldData(h *Halaman, idMU IDMataUang) error {
	return xeIsiXOLKeluar(h, idMU, od+"TreatyXOLList", true)
}

// xeIsiXOLKeluar - badan bersama; `lama` memilih varian OldData.
func xeIsiXOLKeluar(h *Halaman, idMU IDMataUang, tujuan string, lama bool) error {
	k := &kalkulator{}
	hapusDaftarBeserta(h, tujuan) // 1
	// 2 CARI1 = "DUE TO US": hilang oleh Page-New 4.2.2.1 / 4.7.2.1 sebelum dibaca - nol pembaca.
	xeBuangReinsBerbeda(h, h.Ambil(pt+"SOBName")) // 3
	potongan := func(ri, ti int) string {
		if lama { // lama 4.2.2.3.3 / 4.6.2.3.3
			lapis := JalurAnak(JalurAnak(jMaster+"SpreadingTONP", 1, "ReinsuranceListTONP"), 1, "LayerList")
			return JalurAnak(lapis, ti, "DeductionList")
		}
		return JalurAnak(xeTONP(ri), ti, "DeductionList") // 4.2.2.3.3 / 4.7.2.3.3
	}
	// baca = perulangan `.GrossPremiumList` baris TONP (ri, ti) beserta perulangan bersarang net dan potongan
	// (4.2.2.3 / 4.7.2.3; lama 4.2.2.3 / 4.6.2.3). mu = local.tempcurrency.
	baca := func(in halamanXOL, mu string, ri, ti int) {
		for _, g := range h.AmbilDaftar(JalurAnak(xeTONP(ri), ti, "GrossPremiumList")) {
			if g["Currency"] == mu {
				in["CARI32"] = g["Value"]
				in["CARI31"] = g["Currency"]
			}
			for _, n := range h.AmbilDaftar(JalurAnak(xeTONP(ri), ti, "NetPremiumList")) {
				if n["Currency"] == mu {
					in["CARI47"] = n["Value"]
					in["CARI13"] = n["Value"]
					in["CARI31"] = n["Currency"]
				}
			}
			for _, d := range h.AmbilDaftar(potongan(ri, ti)) {
				if d["Currency"] == mu {
					in["CARI44"] = d["Deduction"]
				}
			}
		}
	}
	nol := func() *apd.Decimal { return apd.New(0, 0) }
	var induk []Baris
	for i, inst := range mDaftar(h, "Installment") { // 4
		mu := inst["Currency"] // 4.1
		cur, dueto := "", ""
		gross, net, dtv, ded := nol(), nol(), nol(), nol()
		for ri := range mDaftar(h, "ShareReins") { // 4.2 (Local.Index)
			for ti := range h.AmbilDaftar(xeTONP(ri + 1)) { // 4.2.2 (Local.ShareIndex)
				in := halamanXOL{}       // 4.2.2.1 Page-New
				baca(in, mu, ri+1, ti+1) // 4.2.2.3
				// 4.2.2.4
				cur = in["CARI31"]
				gross = k.Tambah(gross, in.dec(k, "CARI32"))
				net = k.Tambah(net, in.dec(k, "CARI13"))
				dtv = k.Tambah(dtv, in.dec(k, "CARI47"))
				ded = k.Tambah(ded, in.dec(k, "CARI44"))
				dueto = in["CARI1"]
				if !lama { // `local.netpremi = local.netpremi+ @toDecimal(InputXOL.CARI13)` yang kedua
					net = k.Tambah(net, in.dec(k, "CARI13"))
				}
			}
		}
		// 4.4-4.6 (lama 4.3-4.5)
		induk = append(induk, Baris{
			"Currency": cur, "IDCurrency": idMU[mu],
			"GrossPremi": formatAngka(gross), "NetPremi": formatAngka(net), "DueToValue": formatAngka(dtv),
			"Deduction": formatAngka(ded), "DueTo": dueto,
		})
		var lapis []Baris
		for ri := range mDaftar(h, "ShareReins") { // 4.7 (lama 4.6)
			for ti, s := range h.AmbilDaftar(xeTONP(ri + 1)) {
				in := awalLapisan(s) // Page-New + CARI50..53 (4.7.2.1-4.7.2.2)
				baca(in, mu, ri+1, ti+1)
				lapis = append(lapis, barisLapisan(in, idMU[mu])) // 4.7.2.4 (lama 4.6.2.4)
			}
		}
		h.SetelDaftar(JalurAnak(tujuan, i+1, AnakLayerXOL), lapis)
	}
	if k.err != nil {
		return k.err
	}
	h.SetelDaftar(tujuan, induk)
	return nil
}
