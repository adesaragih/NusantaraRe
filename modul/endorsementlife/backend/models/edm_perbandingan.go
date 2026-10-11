package models

// Perbandingan summary sebelum/sesudah endorsement per mata uang - FITUR BARU,
// bukan korpus (permintaan work owner 08-10-2026). Baca-saja, nol migrasi.
//
// Sumber: peserta kasus (`T_PREMIUM_LIST_DETAIL`) dikelompokkan per mata uang
// dan kelas `EDM_STATUS`, dijumlah dengan rumus rekap yang SAMA dengan
// `T_PREMIUM_LIST_SUMMARY` (repository `ekspresiRekap`). Peserta Delete/Batal
// tersimpan dengan kolom uang SUDAH dibalik (negatif) sesudah Process
// Endorsement, sehingga:
//
//	Before     = Old + nilai ASLI Delete/Batal = Old - (Delete+Batal)
//	Deleted    = Delete + Batal (negatif, apa adanya)
//	Added      = New
//	After      = Before + Deleted + Added       (= Old + New = total rekap)
//	Difference = After - Before                 (= Deleted + Added)
//
// Sebelum Process Endorsement belum ada Delete/Batal - Deleted nol; layar
// menandainya "Calculated after Process Endorsement".

import (
	"sort"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// Kelas `EDM_STATUS` untuk perbandingan (kolom bentukan SQL).
const (
	KelompokLama  = "O" // Old (dan status lain/kosong - ikut rekap)
	KelompokHapus = "D" // Delete + Batal
	KelompokBaru  = "N" // New (CSV)
)

// KelompokRekap - jumlah satu (mata uang, kelas) dari basis data.
type KelompokRekap struct {
	Currency string
	Kelompok string
	Cacah    int
	// Nilai - kolom rekap (`KolomRekapKasus[1:]`) → teks desimal.
	Nilai map[string]string
}

// SelPerbandingan - satu baris perbandingan: cacah peserta + kolom rekap.
type SelPerbandingan struct {
	Cacah int               `json:"cacah"`
	Nilai map[string]string `json:"nilai"`
}

// BarisPerbandingan - perbandingan satu mata uang.
type BarisPerbandingan struct {
	Currency   string          `json:"currency"`
	Before     SelPerbandingan `json:"before"`
	Deleted    SelPerbandingan `json:"deleted"`
	Added      SelPerbandingan `json:"added"`
	After      SelPerbandingan `json:"after"`
	Difference SelPerbandingan `json:"difference"`
}

// SusunPerbandingan menghitung perbandingan per mata uang (urut mata uang)
// dengan presisi penuh (apd).
func SusunPerbandingan(kelompok []KelompokRekap, kolom []string) ([]BarisPerbandingan, error) {
	per := map[string]map[string]KelompokRekap{}
	for _, k := range kelompok {
		if per[k.Currency] == nil {
			per[k.Currency] = map[string]KelompokRekap{}
		}
		per[k.Currency][k.Kelompok] = k
	}
	uang := make([]string, 0, len(per))
	for u := range per {
		uang = append(uang, u)
	}
	sort.Strings(uang)
	hasil := make([]BarisPerbandingan, 0, len(uang))
	for _, u := range uang {
		o, d, n := per[u][KelompokLama], per[u][KelompokHapus], per[u][KelompokBaru]
		b := BarisPerbandingan{
			Currency:   u,
			Before:     SelPerbandingan{Cacah: o.Cacah + d.Cacah, Nilai: map[string]string{}},
			Deleted:    SelPerbandingan{Cacah: -d.Cacah, Nilai: map[string]string{}},
			Added:      SelPerbandingan{Cacah: n.Cacah, Nilai: map[string]string{}},
			After:      SelPerbandingan{Cacah: o.Cacah + n.Cacah, Nilai: map[string]string{}},
			Difference: SelPerbandingan{Cacah: n.Cacah - d.Cacah, Nilai: map[string]string{}},
		}
		for _, kol := range kolom {
			vo, err := desimal(o.Nilai[kol])
			if err != nil {
				return nil, err
			}
			vd, err := desimal(d.Nilai[kol])
			if err != nil {
				return nil, err
			}
			vn, err := desimal(n.Nilai[kol])
			if err != nil {
				return nil, err
			}
			ctx := utils.DecimalContext()
			var before, after, beda apd.Decimal
			if _, err := ctx.Sub(&before, vo, vd); err != nil {
				return nil, err
			}
			if _, err := ctx.Add(&after, vo, vn); err != nil {
				return nil, err
			}
			if _, err := ctx.Sub(&beda, &after, &before); err != nil {
				return nil, err
			}
			b.Before.Nilai[kol] = teksDesimal(&before)
			b.Deleted.Nilai[kol] = teksDesimal(vd)
			b.Added.Nilai[kol] = teksDesimal(vn)
			b.After.Nilai[kol] = teksDesimal(&after)
			b.Difference.Nilai[kol] = teksDesimal(&beda)
		}
		hasil = append(hasil, b)
	}
	return hasil, nil
}

// desimal - teks desimal; kosong = nol.
func desimal(v string) (*apd.Decimal, error) {
	if v == "" {
		return apd.New(0, 0), nil
	}
	return utils.ParseDecimal(v)
}

// teksDesimal - bentuk `TM9`: tanpa nol ekor, nol = "0".
func teksDesimal(d *apd.Decimal) string {
	var r apd.Decimal
	r.Reduce(d)
	return utils.FormatDecimal(&r)
}
