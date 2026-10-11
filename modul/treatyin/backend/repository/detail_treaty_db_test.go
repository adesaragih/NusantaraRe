//go:build db

package repository_test

// Tulisan detail kontrak tuntas di Oracle — SELURUHNYA di dalam transaksi
// yang SELALU dibatalkan (`…LaluBatalkanUntukUji`). Nol Commit, data UJI.
//
// ⛔ Menjalankan DELETE/INSERT di DEV (lalu ROLLBACK) — hanya dengan izin WO.

import (
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/repository"
)

func barisDetailUji(kolom []models.KolomDetail) []models.BarisDetailTreaty {
	out := []models.BarisDetailTreaty{}
	for _, mataUang := range []string{"IDR", "USD"} {
		b := models.BarisDetailTreaty{Teks: map[string]string{}, Angka: map[string]*apd.Decimal{}}
		for _, k := range kolom {
			if k.Angka {
				b.Angka[k.Nama] = apd.New(12345, -2) // 123.45
				continue
			}
			b.Teks[k.Nama] = "UJI"
		}
		b.Teks["LIMITCURRENCY"] = mataUang
		b.Teks["CESSIONCURRENCY"] = mataUang
		b.Teks["SPREADINGTYPEID"] = "10260"
		b.Angka["LIMITVALUE"] = nil
		out = append(out, b)
	}
	return out
}

func TestDetailTreatyInSahDiOracleLaluDibatalkan(t *testing.T) {
	g, ctx := gudangBaca(t)
	doc := map[string]any{
		"ProportionType": "Proportional", "TreatyContractName": "UJI DETAIL — DIBATALKAN",
		"Commencement": "20260101", "Termination": "20261231", "TreatyYear": "2026",
		"StatusAkseptasi": models.StatusTuntas,
	}
	id, cacah, err := g.SimpanLaluBatalkanUntukUji(ctx, models.RencanaSimpan{
		Dokumen: doc, Detail: &models.RencanaDetail{Baris: barisDetailUji(models.KolomDetailTreatyIn)},
	})
	if err != nil {
		t.Fatal(err)
	}
	if cacah[repository.TabelDetailTreatyIn] != 2 {
		t.Errorf("kontrak %s: %d baris detail, mau 2", id, cacah[repository.TabelDetailTreatyIn])
	}
}

func TestDetailEDMSahDiOracleLaluDibatalkan(t *testing.T) {
	g, ctx := gudangBaca(t)
	const id = "9999999/R01"
	baru := map[string]any{
		"OLDID": "9999999", "ProportionType": "NonProportional", "TreatyContractName": "UJI DETAIL EDM — DIBATALKAN",
		"TreatyYear": "2026", "StatusAkseptasi": models.StatusTuntas,
	}
	cBaru, _, _, err := g.SimpanPenyesuaianLaluBatalkanUntukUji(ctx, models.RencanaPenyesuaian{
		ID: id, Draf: true, Baru: baru, Detail: &models.RencanaDetail{Baris: barisDetailUji(models.KolomDetailTreatyInEDM)},
	})
	if err != nil {
		t.Fatal(err)
	}
	if cBaru[repository.TabelDetailTreatyInEDM] != 2 {
		t.Errorf("%d baris detail EDM, mau 2", cBaru[repository.TabelDetailTreatyInEDM])
	}
}
