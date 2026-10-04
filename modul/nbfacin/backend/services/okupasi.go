package services

// Sub-tab Occupation - tiket 40: pemeriksaan okupasi objek. Tanpa enumerasi (Category / Class of Construction / PctLimit
// disimpan apa adanya); hanya jumlah baris dan lebar kolom migrasi 189.

import (
	"fmt"

	"nusantarare/modul/nbfacin/backend/models"
)

// batasOkupasi - T_OCCUPATIONLIST.SEQ_NO NUMBER(5).
const batasOkupasi = 99999

// lebarOkupasi - lebar kolom (BYTE) migrasi 189 tiap medan okupasi.
var lebarOkupasi = []struct {
	nama  string
	nilai func(models.OkupasiObjek) string
	n     int
}{
	{"occupationId", func(o models.OkupasiObjek) string { return o.OccupationID }, 1000},
	{"occupationName", func(o models.OkupasiObjek) string { return o.OccupationName }, 1000},
	{"category", func(o models.OkupasiObjek) string { return o.Category }, 50},
	{"constructionClass", func(o models.OkupasiObjek) string { return o.ConstructionClass }, 500},
	{"pctLimit", func(o models.OkupasiObjek) string { return o.PctLimit }, 50},
}

// periksaOkupasi - okupasi baris objek ke-`n` (tanpa basis data). Class of Construction tidak diwajibkan (L-2: Save
// Pega tidak memvalidasinya).
func periksaOkupasi(n int, okupasi []models.OkupasiObjek) []string {
	if len(okupasi) > batasOkupasi {
		return []string{fmt.Sprintf("baris[%d].occupations paling banyak %d", n, batasOkupasi)}
	}
	var masalah []string
	for m, o := range okupasi {
		for _, l := range lebarOkupasi {
			if len(l.nilai(o)) > l.n {
				masalah = append(masalah, fmt.Sprintf("baris[%d].occupations[%d].%s paling banyak %d byte", n, m, l.nama, l.n))
			}
		}
	}
	return masalah
}
