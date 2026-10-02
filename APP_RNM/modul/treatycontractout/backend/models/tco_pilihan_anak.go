package models

// Pilihan ReinsType baris anak SEMUA grid `Show Child` Treaty Desc - isi
// activity `TreatyContractSetReinsTypeList` (`@baseclass`, 4 langkah, nol `//`)
// [keputusan work owner 02-10-2026: "untuk child-nya ikuti XML-nya
// TreatyContractSetReinsTypeList, untuk semua child pada Treaty Desc"].
//
// `[terverifikasi]` `Activity/TreatyContractSetReinsTypeList.xml`:
//
//	1   b267  Page-Remove ReinsTypeList
//	2   b418  satu putaran; setiap anak langkah menimpa ReinsTypeList.pxResults(n)
//	2.1 b418  @contains(InputData.CARIDESCFACIN," QS ")  b628 -> 10028 QS (OR), 10004 QS (R/I), 10007 ORS
//	2.2 b668  @contains(InputData.CARIDESCFACIN," SPL ") b878 -> 10248 SPL (OR), 10249 SPL (RI), 10007 ORS
//	2.3 b918  @contains(InputData.CARIDESCFACIN," XOL ") b1120 -> 10028 QS (OR), 10004 QS (R/I), 10217 XL
//	2.4 b1160 @contains(InputData.CARIDESCFACIN,"ORS")   b1284 -> 10007 ORS
//
// `InputData.CARIDESCFACIN` = `.ReinsTypeName` baris INDUK yang di-`Show Child`
// (`BrowseTreatyArrLimitParentList` langkah 1 b368/b369). Pilihan tampil
// `.CARI2` (nama, dicari) dan `.CARI1` (ID -> `ReinsTypeID`, b3019); nama yang
// tersimpan = `.CARI2` (sel terikat `.ReinsTypeName` b2935).
//
// ⚠️ `@contains` peka huruf besar-kecil dan spasinya bagian dari kata:
// induk bernama "QS 2020" (tanpa spasi di depan) TIDAK cocok " QS ". Tanpa
// kata yang cocok, daftar KOSONG - seperti Pega.

import "strings"

// PilihanReinsAnak - satu baris `ReinsTypeList.pxResults` (`.CARI1`, `.CARI2`).
type PilihanReinsAnak struct {
	ID   string
	Nama string
}

// langkahPilihanReinsAnak - langkah 2.1-2.4, urut XML.
var langkahPilihanReinsAnak = []struct {
	kata string
	isi  []PilihanReinsAnak
}{
	{" QS ", []PilihanReinsAnak{{"10028", "QS (OR)"}, {"10004", "QS (R/I)"}, {"10007", "ORS"}}},
	{" SPL ", []PilihanReinsAnak{{"10248", "SPL (OR)"}, {"10249", "SPL (RI)"}, {"10007", "ORS"}}},
	{" XOL ", []PilihanReinsAnak{{"10028", "QS (OR)"}, {"10004", "QS (R/I)"}, {"10217", "XL"}}},
	{"ORS", []PilihanReinsAnak{{"10007", "ORS"}}},
}

// PilihanReinsAnakDari menjalankan `TreatyContractSetReinsTypeList` atas nama
// ReinsType induk. Setiap langkah yang cocok MENIMPA `pxResults(1..n)` (langkah
// berikutnya menimpa yang sebelumnya pada indeks yang sama); ID yang berulang
// sesudahnya tampil sekali.
func PilihanReinsAnakDari(namaInduk string) []PilihanReinsAnak {
	var hasil []PilihanReinsAnak
	for _, l := range langkahPilihanReinsAnak {
		if !strings.Contains(namaInduk, l.kata) {
			continue
		}
		for i, p := range l.isi {
			if i < len(hasil) {
				hasil[i] = p
			} else {
				hasil = append(hasil, p)
			}
		}
	}
	lihat := map[string]bool{}
	unik := []PilihanReinsAnak{}
	for _, p := range hasil {
		if !lihat[p.ID] {
			lihat[p.ID] = true
			unik = append(unik, p)
		}
	}
	return unik
}
