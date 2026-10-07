package services

// Tab Coverage FIRE tahap C3 - tiket 45: pemeriksaan deductible coverage (tanpa basis data). Tidak ada medan wajib
// (`Section\addDeductible.xml` tanpa Required), termasuk mata uang (A169). Kode pilihan tidak diperiksa
// keanggotaannya (pola butir 85) - hanya bentuknya: kode bulat pyStandardValue (`DDL\TypeDeductible.xml`), kolom NUMBER(5).

import (
	"fmt"
	"regexp"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
)

// polaKodeDeductible - TYPE_DEDUCTIBLE / TYPE_DEDUCTIBLE2 NUMBER(5) (migrasi 194, A170).
var polaKodeDeductible = regexp.MustCompile(`^[0-9]{1,5}$`)

// lebarTimeExcess - T_DEDUCTIBLELIST.TIME_EXCESS VARCHAR2(30) rancangan (teks desimal, A171).
const lebarTimeExcess = 30

// lebarDeductible - lebar kolom (BYTE) migrasi 194 medan teks deductible.
var lebarDeductible = []struct {
	nama  string
	nilai func(models.Deductible) string
	n     int
}{
	{"minMax", func(d models.Deductible) string { return d.MinMax }, 50},
	{"currency", func(d models.Deductible) string { return d.Currency }, 50},
	{"condition", func(d models.Deductible) string { return d.Condition }, 500},
	{"inputCondition", func(d models.Deductible) string { return d.InputCondition }, 500},
}

// kodeDeductible - medan kode bulat.
var kodeDeductible = []struct {
	nama  string
	nilai func(models.Deductible) string
}{
	{"typeDeductible", func(d models.Deductible) string { return d.TypeDeductible }},
	{"typeDeductible2", func(d models.Deductible) string { return d.TypeDeductible2 }},
}

// periksaDeductible - deductible satu coverage; `awalCoverage` = awalan jalur pesan ("…coverages[k].").
func periksaDeductible(awalCoverage string, ded []models.Deductible) []string {
	if len(ded) > batasItem {
		return []string{fmt.Sprintf("%sdeductibles paling banyak %d", awalCoverage, batasItem)}
	}
	var masalah []string
	for i, d := range ded {
		awal := fmt.Sprintf("%sdeductibles[%d].", awalCoverage, i)
		for _, l := range lebarDeductible {
			if len(l.nilai(d)) > l.n {
				masalah = append(masalah, fmt.Sprintf("%s%s paling banyak %d byte", awal, l.nama, l.n))
			}
		}
		for _, k := range kodeDeductible {
			if v := k.nilai(d); v != "" && !polaKodeDeductible.MatchString(v) {
				masalah = append(masalah, awal+k.nama+" harus kode bilangan bulat (paling banyak 5 digit)")
			}
		}
		if d.TimeExcess != nil && len(utils.FormatDecimal(d.TimeExcess)) > lebarTimeExcess {
			masalah = append(masalah, fmt.Sprintf("%stimeExcess paling banyak %d karakter", awal, lebarTimeExcess))
		}
	}
	return masalah
}
