package models

// Aturan isian form `InboxBenefit` - HANYA dari XML (keputusan work owner 08-10-2026 K6):
//
//   - Benefit WAJIB: `pyRequired` true b1137 / b1188, `pyRequiredNew` always b1185.
//   - Benefit HURUF BESAR: perubahan textarea Benefit menjalankan refresh dengan data transform `SetUpperCase_DT`
//     (b1192-b1211, ActionSets b1310-b1327). Isi data transform itu TIDAK ada di XML; yang diterapkan = arti namanya
//     (nilai Benefit dijadikan huruf besar) - PARITAS-LAYAR-DAN-AKSI.md, pertanyaan terbuka T1.
//
// Di luar XML, mengikuti pola modul MASTER TREATY (PARITAS bab "Di luar XML" - WO boleh menolaknya):
//
//   - spasi tepi dipangkas (riratelife / ricommlife / ririsklife memangkas nama);
//   - panjang paling banyak BatasBenefit byte = lebar kolom (tanpa ini Oracle menjawab ORA-12899 = galat 500 mentah).
//
// TIDAK ada penolakan Benefit kembar: XML tidak memuatnya dan isi `AddToList_Act` tidak ada (pola ririsklife
// menolak nama kembar, tetapi K6 melarang menambahkannya tanpa bukti).

import (
	"errors"
	"fmt"
	"strings"
)

// NormalBenefit merapikan Benefit: pangkas, huruf besar (`SetUpperCase_DT`), wajib, paling panjang BatasBenefit byte.
// Galat = kalimat untuk pengguna (bahasa layar: Inggris).
func NormalBenefit(s string) (string, error) {
	t := strings.ToUpper(strings.TrimSpace(s))
	switch {
	case t == "":
		return "", errors.New("Benefit is required")
	case len(t) > BatasBenefit:
		return "", fmt.Errorf("Benefit is longer than %d characters", BatasBenefit)
	}
	return t, nil
}
