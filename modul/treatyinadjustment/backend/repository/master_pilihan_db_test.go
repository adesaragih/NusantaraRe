//go:build db

package repository_test

// Picker Add Revision / Add Adjustment Premium dan dokumen master `Choose`,
// diadu dengan POOLDATA — BACA SAJA (diukur 7 Oktober 2026).

import "testing"

func TestDaftarMasterPilihanSamaDenganSQLEkspor(t *testing.T) {
	g, ctx := bacaSaja(t)
	// `TreatyLoadMasterJoinEdm`: TREATY_IN (1.856) ∪ TREATY_IN_EDM (280).
	semua, err := g.DaftarMasterPilihan(ctx, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(semua) != 2136 {
		t.Errorf("revisi %d baris, mau 2.136", len(semua))
	}
	// `TreatyLoadMasterJoinEdmXOL`: idem, NonProportional saja.
	np, err := g.DaftarMasterPilihan(ctx, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(np) != 1010 {
		t.Errorf("premi %d baris, mau 1.010", len(np))
	}
	for _, b := range np {
		if b.SifatProporsi != "NonProportional" {
			t.Fatalf("baris %s bercabang %q di picker premi", b.ID, b.SifatProporsi)
		}
	}
}

func TestAdaRevisiMengikutiTREATY_IN_EDM(t *testing.T) {
	g, ctx := bacaSaja(t)
	for id, mau := range map[string]bool{
		"1000506":     true,  // R01–R03 ada
		"1000506/R02": true,  // dirinya sendiri cocok `ID%`
		"1000218":     false, // master tanpa revisi
	} {
		ada, err := g.AdaRevisi(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		if ada != mau {
			t.Errorf("AdaRevisi(%s) = %v, mau %v", id, ada, mau)
		}
	}
}

func TestDokumenMasterDilengkapiKolomKepala(t *testing.T) {
	g, ctx := bacaSaja(t)
	d, ada, err := g.BacaDokumenMaster(ctx, "1000506")
	if err != nil {
		t.Fatal(err)
	}
	if !ada || d.Medan["ProportionType"] != "NonProportional" || d.Medan["TreatyYear"] != "2021" {
		t.Errorf("ada %v, kepala %v", ada, d.Medan)
	}
	// Revisi dibaca dari TREATY_IN_EDM, termasuk kepala EDM-nya.
	r, ada, err := g.BacaDokumenMaster(ctx, "1000506/R02")
	if err != nil {
		t.Fatal(err)
	}
	if !ada || r.Medan["OLDID"] != "1000506/R01" || r.Medan["EDMState"] == "" {
		t.Errorf("revisi: ada %v, kepala %v", ada, r.Medan)
	}
	if _, ada, err := g.BacaDokumenMaster(ctx, "9999999"); err != nil || ada {
		t.Errorf("master tak ada: ada %v err %v", ada, err)
	}
}
