package repository

// Total yang TIDAK tersimpan dihitung ulang dari rinciannya saat sisi
// (New / Old / ActualValue / ValueBeforeProrate) dibaca.
//
// ⭐ Laporan pemakai 9 Oktober 2026: sesudah dimuat ulang, rincian Share Prop
// terisi tetapi `Total Value Spreading OR/R/I` "No items"; di mode View
// tombol Refresh/Update Total mati, jadi total itu tidak pernah muncul.
// Sebabnya data: berkas tersimpan ketika total di layar masih kosong.
//
// Rumusnya sama dengan tombolnya (SALINAN `total_cadangan.go` modul Treaty In
// — modul tidak boleh saling impor):
//
//	Share Prop   Σ RNMShareList / RNMSpreadedList(RI) seluruh Limits.Detail
//	             per mata uang (`TreatyInPropshare` [6])
//	EGNPI        Σ `.Amount` per mata uang (`TreatyInNPSetTotal` egnpi [7])
//	Retention    Σ `.Amount` per mata uang (`TreatyInNPSetTotal` retention)
//	Installment  `AmountTotal` tiap halaman per mata uang (`TreatyInNPSetTotal`
//	             [27]-[28]); halaman tanpa AmountTotal: Σ `.Amount` jadwalnya
//
// ⛔ HANYA larik yang KOSONG yang diisi — total tersimpan selalu menang.

import (
	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/treatyinadjustment/backend/models"
)

var ktxTotal = apd.BaseContext.WithPrecision(60)

// jumlahan - Σ per mata uang, urutan kemunculan dipertahankan.
type jumlahan struct {
	urut  []string
	nilai map[string]*apd.Decimal
}

func (j *jumlahan) tambah(cur, v string) {
	if v == "" {
		return
	}
	d, _, err := apd.NewFromString(v)
	if err != nil {
		return
	}
	if j.nilai == nil {
		j.nilai = map[string]*apd.Decimal{}
	}
	if _, ada := j.nilai[cur]; !ada {
		j.urut = append(j.urut, cur)
		j.nilai[cur] = apd.New(0, 0)
	}
	_, _ = ktxTotal.Add(j.nilai[cur], j.nilai[cur], d)
}

func (j *jumlahan) baris() []map[string]string {
	out := []map[string]string{}
	for _, c := range j.urut {
		out = append(out, map[string]string{"Currency": c, "Value": j.nilai[c].Text('f')})
	}
	return out
}

func teksSimpulAny(s map[string]any, k string) string {
	v, _ := s[k].(string)
	return v
}

func anakSimpul(s map[string]any, k string) []map[string]any {
	switch xs := s[k].(type) {
	case []map[string]any:
		return xs
	case []any:
		out := []map[string]any{}
		for _, x := range xs {
			if m, ok := x.(map[string]any); ok {
				out = append(out, m)
			}
		}
		return out
	}
	return nil
}

// lengkapiTotal mengisi larik total yang kosong pada satu sisi.
func lengkapiTotal(s *models.SisiPenyesuaian) {
	if s.Larik == nil {
		s.Larik = map[string][]map[string]string{}
	}
	isi := func(nama string, j *jumlahan) {
		if len(s.Larik[nama]) == 0 && len(j.urut) > 0 {
			s.Larik[nama] = j.baris()
		}
	}

	var share, or, ri jumlahan
	for _, l := range s.Pohon["Limits"] {
		for _, d := range anakSimpul(l, "Detail") {
			for _, x := range []struct {
				larik string
				ke    *jumlahan
			}{{"RNMShareList", &share}, {"RNMSpreadedList", &or}, {"RNMSpreadedListRI", &ri}} {
				for _, c := range anakSimpul(d, x.larik) {
					x.ke.tambah(teksSimpulAny(c, "Currency"), teksSimpulAny(c, "Value"))
				}
			}
		}
	}
	isi("TotalShareRnmProp", &share)
	isi("TotalSpreadedRnmProp", &or)
	isi("TotalSpreadedRnmRIProp", &ri)

	var egnpi, retensi, angsuran jumlahan
	for _, r := range s.Larik["EGNPI"] {
		egnpi.tambah(r["Currency"], r["Amount"])
	}
	isi("TotalEgnpiAmountNP", &egnpi)
	for _, r := range s.Larik["Retention"] {
		retensi.tambah(r["Currency"], r["Amount"])
	}
	isi("TotalRetentionAmountNP", &retensi)
	for _, h := range s.Pohon["Installment"] {
		cur := teksSimpulAny(h, "Currency")
		if t := teksSimpulAny(h, "AmountTotal"); t != "" {
			angsuran.tambah(cur, t)
			continue
		}
		for _, b := range anakSimpul(h, "InstallmentList") {
			c := teksSimpulAny(b, "Currency")
			if c == "" {
				c = cur
			}
			angsuran.tambah(c, teksSimpulAny(b, "Amount"))
		}
	}
	isi("TotalInstallmentNP", &angsuran)
}
