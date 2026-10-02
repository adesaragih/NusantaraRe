package tiruan

// Tiruan `SetPremi_EDM`: penandaan + jurnal balik + rekap mata uang - desimal
// `apd`, bukan float (ADR-0003).

import (
	"context"
	"fmt"
	"sort"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/endorsementlife/backend/models"
)

// tm9 - teks desimal tanpa nol buntut, seperti `TO_CHAR(x, 'TM9')` Oracle.
func tm9(d *apd.Decimal) string {
	var r apd.Decimal
	r.Reduce(d)
	return utils.FormatDecimal(&r)
}

// balik membalik tanda satu nilai teks desimal; kosong tetap kosong.
func balik(v string) (string, error) {
	if v == "" {
		return "", nil
	}
	d, err := utils.ParseDecimal(v)
	if err != nil {
		return "", err
	}
	var n apd.Decimal
	n.Neg(d)
	return tm9(&n), nil
}

// Tandai meniru `sqlTandai`: hanya peserta `Old`, status dan tanda berubah bersama.
func (g *Gudang) Tandai(_ context.Context, _ *db.Tx, kasusID, statusBaru string, p models.PilihanHapus) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return 0, g.Galat
	}
	pilih := map[string]bool{}
	for _, id := range p.Pilih {
		pilih[id] = true
	}
	kecuali := map[string]bool{}
	for _, id := range p.Kecuali {
		kecuali[id] = true
	}
	c := 0
	for _, d := range g.Peserta {
		if d.PolisID != kasusID || d.EdmStatus != models.StatusOld {
			continue
		}
		if p.Semua && kecuali[d.ID] || !p.Semua && !pilih[d.ID] {
			continue
		}
		for _, k := range models.KolomJurnalBalik {
			v, err := balik(d.Nilai[k])
			if err != nil {
				return 0, fmt.Errorf("tiruan: %s %s: %w", d.ID, k, err)
			}
			if d.Nilai[k] != "" {
				d.Nilai[k] = v
			}
		}
		d.EdmStatus = statusBaru
		c++
	}
	return c, nil
}

// jumlah menjumlahkan beberapa kolom satu baris dengan tanda (+1/-1).
func jumlah(n map[string]string, kolom []string, tanda []int) (*apd.Decimal, error) {
	hasil := apd.New(0, 0)
	for i, k := range kolom {
		if n[k] == "" {
			continue
		}
		d, err := utils.ParseDecimal(n[k])
		if err != nil {
			return nil, err
		}
		if tanda[i] < 0 {
			_, err = utils.DecimalContext().Sub(hasil, hasil, d)
		} else {
			_, err = utils.DecimalContext().Add(hasil, hasil, d)
		}
		if err != nil {
			return nil, err
		}
	}
	return hasil, nil
}

// rumus - kolom dan tanda PREMIUM/BALANCE per Type (`rumusPremiBalance`).
func rumus(tipe string) (premi []string, balance []string, tandaBalance []int) {
	switch tipe {
	case models.TypeQR:
		return []string{"GROSS_PREMIUM"}, []string{"GROSS_PREMIUM", "DEDUCTION", "RI_ADMIN_FEE", "BROKERAGE_FEE", "TAX", "PROF_COMM", "CLAIM"},
			[]int{1, -1, -1, -1, -1, -1, -1}
	case models.TypeQP:
		return []string{"GROSS_PREMIUM_REFUND"}, []string{"GROSS_PREMIUM_REFUND", "CLAIM_AMOUNT", "DEDUCTION_REFUND", "BROKERAGE_FEE_REFUND",
			"RI_ADMIN_FEE_REFUND", "TAX", "PROF_COMM", "CLAIM"}, []int{1, 1, -1, -1, -1, -1, -1, -1}
	case models.TypeTP:
		return []string{"GROSS_PREMIUM_RETRO"}, []string{"GROSS_PREMIUM_RETRO", "DISCOUNT_PREMIUM_RETRO", "RI_ADMIN_FEE_RETRO", "BROKERAGE_FEE_RETRO"},
			[]int{1, -1, -1, 1}
	case models.TypeTR:
		return []string{"GROSS_PREMIUM_REFUND_RETRO"}, []string{"GROSS_PREMIUM_REFUND_RETRO", "DISCOUNT_PREMIUM_REFUND_RETRO",
			"RI_ADMIN_FEE_REFUND_RETRO", "BROKERAGE_FEE_REFUND_RETRO"}, []int{1, -1, -1, 1}
	}
	return nil, nil, nil
}

// pasanganRekap - kolom rekap → kolom peserta (`kolomJumlahRekap` repository).
var pasanganRekap = map[string]string{"COMMISSION": "COMM"}

// HitungRekap meniru `sqlHapusRekap` + `sqlSisipRekap`.
func (g *Gudang) HitungRekap(_ context.Context, _ *db.Tx, kasusID, tipe string) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return 0, g.Galat
	}
	perUang := map[string][]map[string]string{}
	for _, d := range g.Peserta {
		if d.PolisID == kasusID {
			perUang[d.Nilai["CURRENCY"]] = append(perUang[d.Nilai["CURRENCY"]], d.Nilai)
		}
	}
	var hasil []map[string]string
	premi, balance, tanda := rumus(tipe)
	for uang, baris := range perUang {
		r := map[string]string{"CURRENCY": uang}
		for _, k := range models.KolomRekapKasus[1:] {
			var kolom []string
			var tk []int
			switch k {
			case "PREMIUM":
				kolom, tk = premi, []int{1}
			case "BALANCE":
				kolom, tk = balance, tanda
			default:
				sumber := k
				if s, ada := pasanganRekap[k]; ada {
					sumber = s
				}
				kolom, tk = []string{sumber}, []int{1}
			}
			total := apd.New(0, 0)
			for _, n := range baris {
				v, err := jumlah(n, kolom, tk)
				if err != nil {
					return 0, err
				}
				if _, err := utils.DecimalContext().Add(total, total, v); err != nil {
					return 0, err
				}
			}
			r[k] = tm9(total)
		}
		hasil = append(hasil, r)
	}
	sort.Slice(hasil, func(i, j int) bool { return hasil[i]["CURRENCY"] < hasil[j]["CURRENCY"] })
	if g.RekapData == nil {
		g.RekapData = map[string][]map[string]string{}
	}
	g.RekapData[kasusID] = hasil
	g.Rekap[kasusID] = len(hasil)
	return len(hasil), nil
}

// RekapKasus meniru `sqlRekapKasus`.
func (g *Gudang) RekapKasus(_ context.Context, _ *db.Tx, kasusID string) ([]map[string]string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return nil, g.Galat
	}
	return g.RekapData[kasusID], nil
}
