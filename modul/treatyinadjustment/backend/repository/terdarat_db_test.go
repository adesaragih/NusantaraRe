//go:build db

package repository_test

// `Terdarat` diadu dengan Oracle — BACA SAJA terhadap `T_TREATY_REVISION`,
// `TREATY_IN`, `TREATY_IN_EDM`, dan tabel pendaratan. Nol tulisan.

import "testing"

func TestTerdaratSamaDenganBarisAkarPendaratan(t *testing.T) {
	g, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)
	ada := func(id string) bool {
		var n int
		if err := h.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+skema+`.T_TREATY_REVISION WHERE MASTERID = :1`, id).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n > 0
	}
	// Master yang tersimpan lewat aplikasi (1002305) dan penyesuaian yang
	// dilaporkan tanpa pendaratan (1001540/R01, §12 LAYAR-ADJUSTMENT).
	d, _, err := g.BacaDokumenMaster(ctx, "1002305")
	if err != nil {
		t.Fatal(err)
	}
	if d.Terdarat != ada("1002305") {
		t.Errorf("master 1002305: Terdarat %v, baris akar %v", d.Terdarat, ada("1002305"))
	}
	p, err := g.BacaPenyesuaianPendaratan(ctx, "1001540/R01")
	if err != nil {
		t.Fatal(err)
	}
	if p.Terdarat != ada("1001540/R01") {
		t.Errorf("1001540/R01: Terdarat %v, baris akar %v", p.Terdarat, ada("1001540/R01"))
	}
	t.Logf("1002305 terdarat=%v · 1001540/R01 terdarat=%v", d.Terdarat, p.Terdarat)
}
