package repository

// Halaman `ValueBeforeProrate` satu penyesuaian — dibaca dari MASTERID
// `id + AkhiranSisiSebelumProrata` dan digabung ke sisi New sebagai kunci
// BERTITIK.
//
// ---------------------------------------------------------------------
// ⛔ CELAH YANG DITUTUP 8 Oktober 2026
// ---------------------------------------------------------------------
// `TreatyEDMProRateCalculation` langkah *"Copy value from ValueDifference to
// ValueBeforeProrate"* menyetel `TreatyIn.ValueBeforeProrate :=
// TreatyIn.ValueDifference` — selisih SEBELUM pro-rate dikenakan — dan tab
// `TreatyInTabsNPValueDifference_NoProRate` membacanya.
//
// Sampai hari ini NOL tabel menampungnya: Save melaporkannya tidak tersimpan,
// dan tab itu kosong sesudah kontrak dimuat ulang. Sekarang ia mendarat di
// tabel `T_TREATY_*` yang SAMA lewat peta yang sama, ber-MASTERID sendiri —
// pola `ActualValue`, nol DDL.

import "nusantarare/modul/treatyinadjustment/backend/models"

const awalanSebelumProrata = "ValueBeforeProrate."

// gabungSebelumProrata menyalin halaman itu ke sisi New berawalan
// `ValueBeforeProrate.`.
//
// ⛔ Kunci bertitik yang SUDAH ada di sisi New tidak ditimpa — aturan yang
// sama dengan `gabungAktual`, dan sebabnya sama: yang tersimpan lebih dulu
// menang, bukan dikarang ulang.
func gabungSebelumProrata(baru, prorata models.SisiPenyesuaian) {
	for k, v := range prorata.Medan {
		if _, ada := baru.Medan[awalanSebelumProrata+k]; !ada {
			baru.Medan[awalanSebelumProrata+k] = v
		}
	}
	for k, v := range prorata.Larik {
		if _, ada := baru.Larik[awalanSebelumProrata+k]; !ada {
			baru.Larik[awalanSebelumProrata+k] = v
		}
	}
	if len(prorata.Pohon) == 0 || baru.Pohon == nil {
		return
	}
	for k, v := range prorata.Pohon {
		if _, ada := baru.Pohon[awalanSebelumProrata+k]; !ada {
			baru.Pohon[awalanSebelumProrata+k] = v
		}
	}
}
