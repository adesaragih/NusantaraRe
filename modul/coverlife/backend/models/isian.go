package models

// Aturan isian form `InboxCoverLife` - HANYA dari XML (keputusan work owner 08-10-2026 C3):
//
//   - Cover WAJIB: `pyRequired` true b843 / b895, `pyRequiredNew` always b888.
//   - Note TIDAK wajib: `pyRequired` false b1020 / b1070, `pyRequiredNew` false b1063.
//   - Huruf TIDAK diubah: XML tidak memuat pengubah huruf (nol data transform `SetUpperCase*`, `pyBehaviors` medan Cover
//     b897 dan Note b1072 kosong).
//
// Di luar XML, mengikuti pola modul MASTER TREATY (PARITAS bab "Di luar XML" - WO boleh menolaknya; C3):
//
//   - spasi tepi Cover dipangkas;
//   - Cover TIDAK BOLEH KEMBAR tanpa beda huruf dan spasi tepi - diperiksa services;
//   - panjang paling banyak lebar kolom (tanpa ini Oracle menjawab ORA-12899 = galat 500 mentah).
//
// Note disimpan apa adanya; isian yang hanya spasi = kosong (NULL, seperti 4/4 baris DEV).

import (
	"errors"
	"fmt"
	"strings"
)

// NormalCover merapikan Cover: pangkas, wajib, paling panjang BatasCover byte. Huruf tidak diubah. Galat = kalimat
// untuk pengguna (bahasa layar: Inggris).
func NormalCover(s string) (string, error) {
	t := strings.TrimSpace(s)
	switch {
	case t == "":
		return "", errors.New("Cover is required")
	case len(t) > BatasCover:
		return "", fmt.Errorf("Cover is longer than %d characters", BatasCover)
	}
	return t, nil
}

// NormalNote - Note opsional: hanya spasi = kosong; paling panjang BatasNote byte.
func NormalNote(s string) (string, error) {
	if strings.TrimSpace(s) == "" {
		return "", nil
	}
	if len(s) > BatasNote {
		return "", fmt.Errorf("Note is longer than %d characters", BatasNote)
	}
	return s, nil
}

// KunciTeks - pencocokan Cover kembar tanpa beda huruf besar-kecil dan spasi tepi (sama dengan `UPPER(TRIM(COVER))` di
// repository).
func KunciTeks(s string) string { return strings.ToUpper(strings.TrimSpace(s)) }
