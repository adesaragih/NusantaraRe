package models

// Aturan isian form `InboxProductType`. XML TIDAK menandai satu medan pun wajib (`pyRequired` true nol kali; Plan Name
// `pyRequired` false b859, Business / Benefit `pyRequiredNew` false). Keputusan work owner 08-10-2026:
//
//   - K5 Plan Name WAJIB dan UNIK tanpa beda huruf - di luar XML, pola modul MASTER TREATY (data DEV: 32 nama unik);
//   - K5 Business dan Benefit WAJIB: XML / aktivitas TIDAK membuktikan boleh kosong (isi `SaveProductTypeLife_Act`
//     tidak ada; data DEV keenam kunci terisi di 32/32 baris), jadi diwajibkan;
//   - K4 teks Business / Benefit dicocokkan ulang server ke master (services) - tidak di sini.
//
// Spasi tepi dipangkas; panjang paling banyak BatasNama byte (lebar kolom 947; tanpa ini ORA-12899 = 500 mentah).
// Huruf TIDAK diubah (XML tanpa data transform huruf besar - refresh b864-b875 tanpa DT).

import (
	"errors"
	"fmt"
	"strings"
)

// teks - pangkas, wajib, paling panjang BatasNama byte. Galat = kalimat pengguna (bahasa layar: Inggris).
func teks(label, s string) (string, error) {
	t := strings.TrimSpace(s)
	switch {
	case t == "":
		return "", fmt.Errorf("%s is required", label)
	case len(t) > BatasNama:
		return "", fmt.Errorf("%s is longer than %d characters", label, BatasNama)
	}
	return t, nil
}

// PeriksaIsian merapikan isian; galat = SEMUA masalah sekaligus, dipisah "; ".
func PeriksaIsian(isi Isian) (Isian, error) {
	var out Isian
	var pesan []string
	var err error
	if out.CoverName, err = teks("Plan Name", isi.CoverName); err != nil {
		pesan = append(pesan, err.Error())
	}
	if out.Business, err = teks("Business", isi.Business); err != nil {
		pesan = append(pesan, err.Error())
	}
	if out.Benefit, err = teks("Benefit", isi.Benefit); err != nil {
		pesan = append(pesan, err.Error())
	}
	if len(pesan) > 0 {
		return Isian{}, errors.New(strings.Join(pesan, "; "))
	}
	return out, nil
}

// KunciTeks - pencocokan tanpa beda huruf besar-kecil dan spasi tepi (K4 master, K5 Plan Name unik).
func KunciTeks(s string) string { return strings.ToUpper(strings.TrimSpace(s)) }
