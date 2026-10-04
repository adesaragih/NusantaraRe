package services

// Sub-tab FEA - tiket 41: pemeriksaan baris FEA. Tiga dropdown (Private Team Fire Brigade, Team & SOP Safety / Risk
// Management) tanpa enumerasi (M-4); hanya jumlah unit (M-2) dan lebar kolom migrasi 190.

import (
	"fmt"

	"nusantarare/modul/nbfacin/backend/models"
)

// batasFEA - T_FEALIST.SEQ_NO NUMBER(5).
const batasFEA = 99999

// lebarFEA - lebar kolom (BYTE) migrasi 190 tiap medan; unit = jumlah unit (pxNumber `Section\InputFEA_IsUW.xml`):
// kosong atau bilangan bulat >= 0 (M-2, polaUnit).
var lebarFEA = []struct {
	nama  string
	nilai func(models.BarisFEA) string
	n     int
	unit  bool
}{
	{"apar", func(f models.BarisFEA) string { return f.APAR }, 50, true},
	{"sprinkler", func(f models.BarisFEA) string { return f.Sprinkler }, 50, true},
	{"smokeDetector", func(f models.BarisFEA) string { return f.SmokeDetector }, 50, true},
	{"hydrant", func(f models.BarisFEA) string { return f.Hydrant }, 50, true},
	{"privateTruckBrigade", func(f models.BarisFEA) string { return f.PrivateTruckBrigade }, 50, true},
	{"privateFireBrigade", func(f models.BarisFEA) string { return f.PrivateFireBrigade }, 50, false},
	{"teamSopSafety", func(f models.BarisFEA) string { return f.TeamSOPSafety }, 50, false},
	{"teamSopRiskManagement", func(f models.BarisFEA) string { return f.TeamSOPRiskManagement }, 50, false},
	{"info", func(f models.BarisFEA) string { return f.Info }, 500, false},
}

// periksaFEA - FEA baris objek ke-`n` (tanpa basis data).
func periksaFEA(n int, fea []models.BarisFEA) []string {
	if len(fea) > batasFEA {
		return []string{fmt.Sprintf("baris[%d].fea paling banyak %d", n, batasFEA)}
	}
	var masalah []string
	for m, f := range fea {
		for _, l := range lebarFEA {
			v := l.nilai(f)
			if l.unit && v != "" && !polaUnit.MatchString(v) {
				masalah = append(masalah, fmt.Sprintf("baris[%d].fea[%d].%s harus bilangan bulat >= 0", n, m, l.nama))
			}
			if len(v) > l.n {
				masalah = append(masalah, fmt.Sprintf("baris[%d].fea[%d].%s paling banyak %d byte", n, m, l.nama, l.n))
			}
		}
	}
	return masalah
}
