package repository

// Halaman `ActualValue` satu penyesuaian — dibaca dari MASTERID
// `id + AkhiranSisiAktual` (keputusan pemakai 8 Oktober 2026) dan digabung ke
// sisi New sebagai kunci BERTITIK.
//
// ⭐ Bentuk bertitik itulah yang kerangka ekspor dan rumus layar ikat:
// medan `ActualValue.TotalEgnpiAmount`, larik `ActualValue.EGNPI`, pohon
// `ActualValue.Share` (baris bersarang `GrossPremiumList` …). Penulisnya
// (`treatyin/services.halamanAktual`) mengembalikan kunci bertitik itu ke
// halaman tertanam sebelum mendaratkannya — gabungan ini kebalikannya.

import "nusantarare/modul/treatyinadjustment/backend/models"

const awalanAktual = "ActualValue."

// gabungAktual menyalin halaman Actual ke sisi New berawalan `ActualValue.`.
//
// ⛔ Kunci bertitik yang SUDAH ada di sisi New tidak ditimpa: sisi New
// pendaratannya sendiri hanya membawa `ValueDifference.*` bertitik, jadi
// tabrakan berarti data yang tidak semestinya — dan yang tersimpan lebih
// dulu menang, bukan dikarang ulang.
func gabungAktual(baru, aktual models.SisiPenyesuaian) {
	for k, v := range aktual.Medan {
		if _, ada := baru.Medan[awalanAktual+k]; !ada {
			baru.Medan[awalanAktual+k] = v
		}
	}
	for k, v := range aktual.Larik {
		if _, ada := baru.Larik[awalanAktual+k]; !ada {
			baru.Larik[awalanAktual+k] = v
		}
	}
	if len(aktual.Pohon) == 0 {
		return
	}
	if baru.Pohon == nil {
		// Peta milik pemanggil: `sisiKosong`/`bacaSisi` selalu mengisinya,
		// jadi cabang ini hanya penjaga — penugasan ke salinan nilai struct
		// tidak akan sampai ke pemanggil.
		return
	}
	for k, v := range aktual.Pohon {
		if _, ada := baru.Pohon[awalanAktual+k]; !ada {
			baru.Pohon[awalanAktual+k] = v
		}
	}
}
