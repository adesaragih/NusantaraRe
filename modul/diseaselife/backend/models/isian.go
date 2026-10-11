package models

// Aturan isian form `InboxDisease` (keputusan work owner 08-10-2026 D3):
//
//   - Disease WAJIB: `pyRequired` true b1465 / b1515, `pyRequiredNew` always b1511.
//   - HURUF BESAR pada ICD Code DAN Disease: perubahan KEDUA medan menjalankan refresh dengan data transform
//     `SetUpperCase_DT` (ICD Code: behaviors b1246-b1283 / ActionSets b1353-b1406 = b1265 / b1386; Disease: b1519-b1556 /
//     b1622-b1675 = b1538 / b1654). Medan "Number / ID" tidak (disabled, nol behaviors b1071). Isi data transform itu TIDAK
//     ada di XML; yang diterapkan = arti namanya pada medan yang memanggilnya (preseden benefitlife).
//
// ICD Code WAJIB = keputusan work owner D3, BUKAN XML: XML menandai ICD Code `pyRequired` false b1195 / b1244 dan
// `pyRequiredNew` false b1237 (PARITAS "Di luar XML" + pertanyaan terbuka; prompt menyebutnya "wajib (XML)" - temuan).
//
// Di luar XML, mengikuti pola modul MASTER TREATY (PARITAS bab "Di luar XML" - WO boleh menolaknya):
//
//   - spasi tepi dipangkas;
//   - ICD Code TIDAK BOLEH KEMBAR tanpa beda huruf dan spasi tepi (D3; data DEV 97.586 ICD unik) - diperiksa services;
//   - panjang paling banyak lebar kolom (tanpa ini Oracle menjawab ORA-12899 = galat 500 mentah).

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// HurufBesar - `SetUpperCase_DT` (b1265 / b1538): nilai medan dijadikan huruf besar.
func HurufBesar(s string) string { return strings.ToUpper(s) }

// NormalICDCode merapikan ICD Code: pangkas, huruf besar, wajib (D3), paling panjang BatasICDCode byte. Galat = kalimat
// untuk pengguna (bahasa layar: Inggris).
func NormalICDCode(s string) (string, error) {
	t := HurufBesar(strings.TrimSpace(s))
	switch {
	case t == "":
		return "", errors.New("ICD Code is required")
	case len(t) > BatasICDCode:
		return "", fmt.Errorf("ICD Code is longer than %d characters", BatasICDCode)
	}
	return t, nil
}

// NormalDisease merapikan Disease: pangkas, huruf besar, wajib (b1465), paling panjang BatasDisease byte.
func NormalDisease(s string) (string, error) {
	t := HurufBesar(strings.TrimSpace(s))
	switch {
	case t == "":
		return "", errors.New("Disease is required")
	case len(t) > BatasDisease:
		return "", fmt.Errorf("Disease is longer than %d characters", BatasDisease)
	}
	return t, nil
}

// KunciTeks - pencocokan ICD Code kembar tanpa beda huruf besar-kecil dan spasi tepi (sama dengan
// `UPPER(TRIM(ICD_CODE))` di repository).
func KunciTeks(s string) string { return strings.ToUpper(strings.TrimSpace(s)) }

// RapikanSaring - kata saring: pangkas, dipotong paling panjang BatasSaring byte pada batas huruf (kata yang lebih
// panjang dari kolom terlebar tidak pernah cocok).
func RapikanSaring(s string) string {
	t := strings.TrimSpace(s)
	for len(t) > BatasSaring {
		_, n := utf8.DecodeLastRuneInString(t)
		t = t[:len(t)-n]
	}
	return t
}
