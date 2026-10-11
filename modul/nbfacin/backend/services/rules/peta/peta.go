// Package peta - pembetulan nama workbasket (ADR-0032, penyimpangan sadar) yang dipakai BERSAMA generator registry When
// (`services/rules/bangkit`) dan aturan tahap layanan (`services/bacasaja.go`): satu sumber pemetaan lama -> baru.
//
// Keputusan work owner 07-10-2026: "jika ada yang menggunakan ReasFacInJuniorUnderwriting ganti menjadi
// ReasFacInJuniorUnderwritingB, karena kedepannya akan menggunakan itu"; 08-10-2026: "di when jika ada workbasket
// ReasFacInJuniorUnderwriting diubah jadi ReasFacInJuniorUnderwritingB". Tabel pemetaan juga di docs/issues/58 bab
// "Pemetaan workbasket". Paket daun: tanpa impor paket modul, supaya generator tetap dapat dibangun walau registry rusak.
package peta

// Workbasket - nama workbasket korpus Pega -> nama sistem baru.
var Workbasket = map[string]string{"ReasFacInJuniorUnderwriting": "ReasFacInJuniorUnderwritingB"}

// Petakan - nama baru bila nama korpus dipetakan; selain itu apa adanya.
func Petakan(nama string) string {
	if baru, ada := Workbasket[nama]; ada {
		return baru
	}
	return nama
}
