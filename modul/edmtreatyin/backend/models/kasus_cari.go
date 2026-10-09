package models

import "strings"

// MaksKataCari - kata kotak saring portal terbanyak yang dipakai; kata sesudahnya diabaikan (batas argumen SQL).
const MaksKataCari = 5

// KataCari - kotak saring portal dipecah per spasi, huruf besar, paling banyak `MaksKataCari` kata (perintah work
// owner 07-10-2026, `SaringanKasus.Cari`).
func KataCari(cari string) []string {
	k := strings.Fields(strings.ToUpper(cari))
	if len(k) > MaksKataCari {
		k = k[:MaksKataCari]
	}
	return k
}

// CocokCari - setiap kata `cari` dimuat (tanpa beda huruf) salah satu `nilai`; kotak kosong = cocok. Padanan SQL
// portal untuk gudang tiruan.
func CocokCari(cari string, nilai ...string) bool {
	for _, k := range KataCari(cari) {
		ada := false
		for _, v := range nilai {
			if strings.Contains(strings.ToUpper(v), k) {
				ada = true
				break
			}
		}
		if !ada {
			return false
		}
	}
	return true
}
