package models

// Aturan isian form `InboxCauseofLossLife` - HANYA dari XML (keputusan work owner 08-10-2026 K4):
//
//   - Cause of Loss WAJIB: `pyRequired` true b996 / b1044, `pyRequiredNew` always b1037.
//   - Huruf TIDAK diubah: XML tidak memuat pengubah huruf (nol data transform `SetUpperCase*`, nol refresh pada medan
//     Cause of Loss - b1027-b1131); data DEV memang huruf besar, tetapi itu isi, bukan aturan.
//
// Di luar XML, mengikuti pola modul MASTER TREATY (PARITAS bab "Di luar XML" - WO boleh menolaknya):
//
//   - spasi tepi dipangkas (riratelife / ricommlife / ririsklife / benefitlife / planlife memangkas nama);
//   - nama TIDAK BOLEH KEMBAR tanpa beda huruf dan spasi tepi (pola planlife K5, ririsklife) - diperiksa services;
//   - panjang paling banyak BatasCauseOfLoss byte = lebar kolom (tanpa ini Oracle menjawab ORA-12899 = galat 500).

import (
	"errors"
	"fmt"
	"strings"
)

// NormalCauseOfLoss merapikan Cause of Loss: pangkas, wajib, paling panjang BatasCauseOfLoss byte. Huruf tidak diubah.
// Galat = kalimat untuk pengguna (bahasa layar: Inggris).
func NormalCauseOfLoss(s string) (string, error) {
	t := strings.TrimSpace(s)
	switch {
	case t == "":
		return "", errors.New("Cause of Loss is required")
	case len(t) > BatasCauseOfLoss:
		return "", fmt.Errorf("Cause of Loss is longer than %d characters", BatasCauseOfLoss)
	}
	return t, nil
}

// KunciTeks - pencocokan nama kembar tanpa beda huruf besar-kecil dan spasi tepi (sama dengan
// `UPPER(TRIM(CAUSEOFLOSS))` di repository).
func KunciTeks(s string) string { return strings.ToUpper(strings.TrimSpace(s)) }
