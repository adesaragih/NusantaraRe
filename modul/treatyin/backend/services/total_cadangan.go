package services

import "nusantarare/modul/treatyin/backend/models"

// lengkapiTotalPenampung - total penampung yang TIDAK tersimpan dihitung
// ulang dari rinciannya saat kontrak dibuka.
//
// ⭐ Laporan pemakai 9 Oktober 2026: tab Share Prop sesudah dimuat ulang -
// rincian `Value Spreading OR/R/I` terisi, tetapi `Total Value Spreading
// OR/R/I` "No items"; di mode View semua tombol (Refresh, Update Total)
// mati, jadi total itu tidak pernah dapat dimunculkan. Sebabnya data: kontrak
// tersimpan ketika total di layar masih kosong (sebelum rumus dijalankan,
// atau sebelum `T_TREATY_TOTAL` menyimpan larik ini - migrasi 448).
//
// Rumusnya SAMA dengan tombolnya, bukan rumus baru:
//
//	Share Prop  `jumlahkanTotal` - Σ RNMShareList / RNMSpreadedList(RI)
//	            seluruh Limits.Detail per mata uang (`TreatyInPropshare` [6])
//	EGNPI       Σ `.Amount` per mata uang (`TreatyInNPSetTotal` egnpi [7])
//	Installment Σ `.Amount` jadwal per mata uang = `AmountTotal` tiap halaman
//	            (`SetTotalInstallment`, `TreatyInNPSetTotal` [27]-[28])
//
// ⛔ HANYA larik yang KOSONG yang diisi - total tersimpan selalu menang. Total
// yang hasilnya nol baris tetap tidak diisi (tab menyemai bawaannya).
func lengkapiTotalPenampung(k *models.KontrakWarisan) {
	if k.PenampungLarik == nil {
		k.PenampungLarik = map[string][]map[string]any{}
	}
	kosong := func(n string) bool { return len(k.PenampungLarik[n]) == 0 }
	isi := func(n string, xs []NilaiMataUang) {
		if kosong(n) && len(xs) > 0 {
			k.PenampungLarik[n] = keSimpul(xs)
		}
	}

	if kosong("TotalShareRnmProp") || kosong("TotalSpreadedRnmProp") || kosong("TotalSpreadedRnmRIProp") {
		kp := &konteksProp{limits: k.LimitsPohon}
		kp.jumlahkanTotal(false)
		isi("TotalShareRnmProp", kp.share)
		isi("TotalSpreadedRnmProp", kp.or)
		isi("TotalSpreadedRnmRIProp", kp.ri)
	}

	if kosong("TotalEgnpiAmountNP") {
		per := []NilaiMataUang{}
		for _, r := range k.Egnpi {
			per = tambahPerMataUang(per, r.MataUang, "", angka(r.Jumlah))
		}
		isi("TotalEgnpiAmountNP", per)
	}

	if kosong("TotalInstallmentNP") {
		per := []NilaiMataUang{}
		for _, r := range k.Angsuran {
			per = tambahPerMataUang(per, r.MataUang, "", angka(r.Jumlah))
		}
		isi("TotalInstallmentNP", per)
	}
}
