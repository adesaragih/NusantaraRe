package models

import (
	"sort"
	"strings"
)

// UrutanJenisReasuransi - urutan dropdown `REINS TYPE` menurut teks tampilnya
// (`Note`), keputusan work owner 04-10-2026. Jenis lain di master yang tidak
// tercantum TETAP ditampilkan, sesudahnya, dengan urutan aslinya - tidak
// dibuang diam-diam.
var UrutanJenisReasuransi = []string{"QS", "2ND QS", "SURPLUS", "2ND SURPLUS", "OR"}

// UrutkanJenisReasuransi mengurutkan pilihan menurut `UrutanJenisReasuransi`
// (Note dibandingkan tanpa beda huruf besar/kecil dan spasi tepi). Stabil.
func UrutkanJenisReasuransi(d []JenisReasuransi) []JenisReasuransi {
	nilai := func(j JenisReasuransi) int { return peringkatJenis(j.Note) }
	keluar := append([]JenisReasuransi(nil), d...)
	sort.SliceStable(keluar, func(a, b int) bool { return nilai(keluar[a]) < nilai(keluar[b]) })
	return keluar
}

// JenisKontrakGanda menjawab apakah `reinsTypeID` sudah dipakai kontrak LAIN
// (ID berbeda dari `idSendiri`) di daftar kontrak satu tahun treaty -
// keputusan work owner 04-10-2026: satu REINS TYPE satu kontrak per tahun.
func JenisKontrakGanda(daftar []Kontrak, idSendiri, reinsTypeID string) bool {
	jenis := strings.TrimSpace(reinsTypeID)
	for _, k := range daftar {
		if k.ID != idSendiri && strings.TrimSpace(k.ReinsTypeID) == jenis {
			return true
		}
	}
	return false
}

// peringkatJenis - posisi nama jenis di `UrutanJenisReasuransi` (tanpa beda
// huruf besar/kecil dan spasi tepi); jenis yang tidak tercantum di belakang.
func peringkatJenis(nama string) int {
	n := strings.ToUpper(strings.TrimSpace(nama))
	for i, u := range UrutanJenisReasuransi {
		if u == n {
			return i
		}
	}
	return len(UrutanJenisReasuransi)
}

// UrutkanKontrakMenurutJenis mengurutkan grid kontrak satu tahun treaty
// menurut REINS TYPE-nya (`ReinsTypeName`), urutan yang SAMA dengan dropdown
// - keputusan work owner 04-10-2026. Stabil: jenis yang sama atau tidak
// tercantum tetap pada urutan aslinya (ID).
func UrutkanKontrakMenurutJenis(d []Kontrak) []Kontrak {
	keluar := append([]Kontrak(nil), d...)
	sort.SliceStable(keluar, func(a, b int) bool {
		return peringkatJenis(keluar[a].ReinsTypeName) < peringkatJenis(keluar[b].ReinsTypeName)
	})
	return keluar
}

// BusinessGanda menjawab apakah `bizCode` sudah dipakai business LAIN (ID
// berbeda dari `idSendiri`) di daftar business satu kontrak - keputusan work
// owner 04-10-2026: satu Business Name sekali per kontrak.
func BusinessGanda(daftar []Business, idSendiri, bizCode string) bool {
	kode := strings.TrimSpace(bizCode)
	for _, b := range daftar {
		if b.ID != idSendiri && strings.TrimSpace(b.BizCode) == kode {
			return true
		}
	}
	return false
}

// ReinsurerGanda menjawab apakah `reinsurerID` sudah dipakai reinsurer LAIN
// (ID berbeda dari `idSendiri`) di daftar reinsurer satu kontrak - keputusan
// work owner 04-10-2026: satu Reinsurer Name sekali per kontrak.
func ReinsurerGanda(daftar []Reinsurer, idSendiri, reinsurerID string) bool {
	id := strings.TrimSpace(reinsurerID)
	for _, r := range daftar {
		if r.ID != idSendiri && strings.TrimSpace(r.ReinsurerID) == id {
			return true
		}
	}
	return false
}
