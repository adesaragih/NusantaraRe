//go:build ujidev

package repository

// Uji DEV cadangan tabel bawah spreading (keputusan work owner 09-10-2026): AnakSpreading dijalankan (SELECT) dengan
// induk / treaty group / tahun yang ADA di PROPORTIONALARRG, dan cacah master treaty ber-SpreadingList kosong.

import (
	"context"
	"testing"
	"time"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
)

func TestAnakSpreadingDiDEV(t *testing.T) {
	cfg, err := config.Load()
	if err != nil || !cfg.PunyaOracle() {
		t.Skip("ORACLE_DSN belum dikonfigurasi")
	}
	d, err := db.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = d.Close() }()
	ctx, batal := context.WithTimeout(context.Background(), 2*time.Minute)
	defer batal()
	a := AcuanDari(Baru(d), false)
	pa, err := d.Qualify("PROPORTIONALARRG")
	if err != nil {
		t.Fatal(err)
	}
	var induk, grup, tahun string
	if err := d.QueryRowContext(ctx, `SELECT PARENTREINSTYPEID, TREATYGROUPID, TO_CHAR(TREATYYEAR) FROM `+pa+
		` WHERE PARENTREINSTYPEID IS NOT NULL AND TREATYGROUPID IS NOT NULL AND ROWNUM = 1`).Scan(&induk, &grup, &tahun); err != nil {
		t.Fatal(err)
	}
	anak, err := a.AnakSpreading(ctx, induk, tahun, grup)
	if err != nil {
		t.Fatalf("AnakSpreading: %v", err)
	}
	t.Logf("AnakSpreading contoh: %d anak (tahun %s)", len(anak), tahun)
	if len(anak) == 0 {
		t.Error("contoh induk yang ada di PROPORTIONALARRG harus punya anak")
	}
	for _, m := range []string{"M_TREATY_IN", "M_TREATY_IN_EDM"} {
		tb, err := d.Qualify(m)
		if err != nil {
			t.Fatal(err)
		}
		var semua, kosong int
		if err := d.QueryRowContext(ctx, `SELECT COUNT(*), SUM(CASE WHEN JSON_EXISTS(JSONDATA,
			'$.Limits[0].Detail[0].SpreadingList[0].ReinsTypeID') THEN 0 ELSE 1 END) FROM `+tb).Scan(&semua, &kosong); err != nil {
			t.Fatalf("%s: %v", m, err)
		}
		t.Logf("%s: %d master, %d tanpa SpreadingList (kini memakai PROPORTIONALARRG)", m, semua, kosong)
	}
}
