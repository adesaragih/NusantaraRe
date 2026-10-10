//go:build db

package services_test

// Ukur ulang `TreatyInNonAddItem(share)` lewat KODE terhadap Gross Premium
// yang Pega simpan — BACA SAJA, nol transaksi.
//
// Masukan: layer Limits Non-Prop tersimpan (MDPList) + `RNMShare` salinan
// baris + `FacultativeShare` revisi. Keluaran Gross tiap baris/mata uang
// diadu dengan `T_TREATY_SHARE_AMOUNT` (GrossPremiumList).

import (
	"context"
	"testing"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/repository"
	"nusantarare/modul/treatyin/backend/services"
)

func TestUpdateSummaryShareCocokDenganGrossPega(t *testing.T) {
	cfg, _ := config.Load()
	if !cfg.PunyaOracle() {
		t.Skip("lewati: ORACLE_DSN kosong")
	}
	d, err := db.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = d.Close() })
	ctx := context.Background()
	g := repository.Baru(d)
	sk := cfg.OracleSchema

	rows, err := d.QueryContext(ctx, `SELECT DISTINCT MASTERID FROM `+sk+`.T_TREATY_SHARE
	    WHERE MASTERID NOT LIKE '%#LAMA' AND TRIM(RNMSHARE) IS NOT NULL
	    ORDER BY MASTERID FETCH FIRST 120 ROWS ONLY`)
	if err != nil {
		t.Fatal(err)
	}
	var kontrak []string
	for rows.Next() {
		var m string
		if err := rows.Scan(&m); err != nil {
			t.Fatal(err)
		}
		kontrak = append(kontrak, m)
	}
	_ = rows.Close()

	total, cocok := 0, 0
	for _, m := range kontrak {
		pohon, err := g.BacaPohonLimitsPendaratan(ctx, m)
		if err != nil {
			t.Fatal(err)
		}
		sp, err := g.BacaSharePendaratan(ctx, m)
		if err != nil {
			t.Fatal(err)
		}
		rev, err := g.BacaRevisiPendaratan(ctx, m)
		if err != nil {
			t.Fatal(err)
		}
		simpan := services.ShareDariPendaratan(sp, rev.Medan)
		hitung := simpan
		services.NonAddItemShare(&hitung, services.LayerDariPohon(pohon))
		for i, b := range simpan.Share {
			if i >= len(hitung.Share) {
				break
			}
			for _, v := range b.GrossPremiumList {
				total++
				for _, w := range hitung.Share[i].GrossPremiumList {
					if w.Currency == v.Currency && samaBilanganNP(bulat2(w.Value), bulat2(v.Value)) {
						cocok++
						break
					}
				}
			}
		}
	}
	if total == 0 {
		t.Skip("nol baris Gross untuk diukur")
	}
	persen := cocok * 100 / total
	t.Logf("Gross Update Summary cocok %d/%d (%d%%) atas %d kontrak", cocok, total, persen, len(kontrak))
	// Lantai — sisa = nilai yang tidak dihitung ulang sesudah MDP berubah.
	if persen < 90 {
		t.Errorf("hanya %d%% cocok", persen)
	}
}
