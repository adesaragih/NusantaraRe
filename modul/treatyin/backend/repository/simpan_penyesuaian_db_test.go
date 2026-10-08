//go:build db

package repository_test

// Jalur simpan layar ADJUSTMENT — `TREATY_IN_EDM` + pendaratan kedua sisi.
// ⛔ Seluruhnya di dalam transaksi yang DIBATALKAN: nol Commit di uji.
// ⛔ Nol akses ke dokumen JSON EDM — hanya `TREATY_IN_EDM` dan `T_TREATY_*`.

import (
	"errors"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/repository"
)

func TestSimpanDrafPenyesuaianLaluDibatalkan(t *testing.T) {
	g, ctx := gudangBaca(t)
	const id = "9999999/R01"
	baru := map[string]any{
		"OLDID": "9999999", "ProportionType": "NonProportional", "TreatyContractName": "UJI EDM — DIBATALKAN",
		"TreatyYear": "2026", "EDMState": "1", "EDMMaterialType": "2", "Position": "ReasTreatyInAdmin",
		"PositionUsername": "UJI",
		"Portfolio":        []any{map[string]any{"Description": "baru"}},
		"CommentList":      []any{map[string]any{"OperatorName": "UJI", "Suggest": "UJI  Had Created Internal Edit"}},
	}
	lama := map[string]any{"ProportionType": "NonProportional", "Portfolio": []any{map[string]any{"Description": "lama"}}}
	cBaru, cLama, kepala, err := g.SimpanPenyesuaianLaluBatalkanUntukUji(ctx, models.RencanaPenyesuaian{ID: id, Draf: true, Baru: baru, Lama: lama})
	if err != nil {
		t.Fatal(err)
	}
	if cBaru["T_TREATY_REVISION"] != 1 || cBaru["T_TREATY_PORTFOLIO"] != 1 || cBaru["T_VIEW_COMMENT"] != 1 || cLama["T_TREATY_PORTFOLIO"] != 1 {
		t.Errorf("cacah baru %v lama %v", cBaru, cLama)
	}
	if kepala["OLDID"] != "9999999" || kepala["EDMState"] != "1" || kepala["Position"] != "ReasTreatyInAdmin" || len(kepala["EDMDATE"].(string)) != 9 {
		t.Errorf("kepala %v", kepala)
	}
	// Dibatalkan — nol baris tersisa.
	if _, ada, _ := g.BacaKepalaPenyesuaian(ctx, id); ada {
		t.Errorf("kepala %s tertinggal", id)
	}
}

func TestDrafBerpengenalTersimpanDitolak(t *testing.T) {
	g, ctx := gudangBaca(t)
	k, ada, err := g.BacaKepalaPenyesuaian(ctx, "1000080/R02")
	if err != nil || !ada || k["OLDID"] != "1000080/R01" {
		t.Fatalf("kepala %v ada %v err %v", k, ada, err)
	}
	_, _, _, err = g.SimpanPenyesuaianLaluBatalkanUntukUji(ctx, models.RencanaPenyesuaian{ID: "1000080/R02", Draf: true, Baru: map[string]any{}})
	if !errors.Is(err, repository.ErrPenyesuaianSudahAda) {
		t.Errorf("galat %v", err)
	}
	// Bukan draf: UPDATE kepala yang ada.
	doc := map[string]any{}
	for a, b := range k {
		doc[a] = b
	}
	doc["Information"] = "UJI UBAH — DIBATALKAN"
	_, _, kepala, err := g.SimpanPenyesuaianLaluBatalkanUntukUji(ctx, models.RencanaPenyesuaian{ID: "1000080/R02", Baru: doc})
	if err != nil || kepala["Information"] != "UJI UBAH — DIBATALKAN" || kepala["OLDID"] != "1000080/R01" {
		t.Errorf("ubah %v %v", kepala, err)
	}
	sesudah, _, _ := g.BacaKepalaPenyesuaian(ctx, "1000080/R02")
	if sesudah["Information"] == "UJI UBAH — DIBATALKAN" {
		t.Error("UPDATE tertinggal sesudah rollback")
	}
}

func TestDeclineOfferEDMMenghapusLaluDibatalkan(t *testing.T) {
	g, ctx := gudangBaca(t)
	masih, err := g.HapusPenyesuaianLaluBatalkanUntukUji(ctx, "1000080/R02")
	if err != nil || masih {
		t.Errorf("masih %v err %v", masih, err)
	}
	if _, ada, _ := g.BacaKepalaPenyesuaian(ctx, "1000080/R02"); !ada {
		t.Error("baris hilang sesudah rollback")
	}
	// ⛔ Pengenal KONTRAK tidak pernah dihapus lewat jalur ini.
	if _, err := g.HapusPenyesuaianLaluBatalkanUntukUji(ctx, "1000080"); !errors.Is(err, repository.ErrBukanPengenalPenyesuaian) {
		t.Errorf("pengenal kontrak: %v", err)
	}
}
