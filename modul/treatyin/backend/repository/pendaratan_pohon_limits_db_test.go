//go:build db

package repository_test

// ⭐ Pohon tab Limits dari tabel pendaratan memuat larik yang dulu tidak
// dibaca — 100% Limit/Retention/Cession/EPI (Prop) dan Treaty Group per
// layer + MDP/PE/EGNPI (Non-Prop) — beserta larik akar Summary dan Total.

import "testing"

func TestPohonLimitsMemuatSeluruhLarik(t *testing.T) {
	g, ctx := gudangBaca(t)
	sdb, sk := bukaOracle(t)
	// Kontrak yang punya baris di setiap tabel yang diuji — dipilih dari
	// data, bukan ditanam.
	var prop, np string
	if err := sdb.QueryRowContext(ctx, `SELECT MIN(MASTERID) FROM `+sk+`.T_TREATY_LIMIT_AMOUNT WHERE JENIS = 'IOOLimitList' AND MASTERID NOT LIKE '%#LAMA'`).Scan(&prop); err != nil {
		t.Fatal(err)
	}
	if err := sdb.QueryRowContext(ctx, `SELECT MIN(MASTERID) FROM `+sk+`.T_TREATY_LIMIT_GROUP WHERE MASTERID NOT LIKE '%#LAMA' AND MASTERID IN (SELECT MASTERID FROM `+sk+`.T_TREATY_LIMIT_SUMMARY)`).Scan(&np); err != nil {
		t.Fatal(err)
	}
	pp, err := g.BacaPohonLimitsPendaratan(ctx, prop)
	if err != nil {
		t.Fatal(err)
	}
	adaIOO := false
	for _, l := range pp {
		for _, d := range l["Detail"].([]map[string]any) {
			if len(d["IOOLimitList"].([]map[string]any)) > 0 {
				adaIOO = true
			}
			for _, k := range []string{"RetentionList", "CessionList", "EPIList", "DeductionList", "ReserveList"} {
				if _, ok := d[k].([]map[string]any); !ok {
					t.Errorf("%s bukan larik di kontrak %s", k, prop)
				}
			}
		}
	}
	if !adaIOO {
		t.Errorf("kontrak %s: IOOLimitList kosong", prop)
	}
	pn, err := g.BacaPohonLimitsPendaratan(ctx, np)
	if err != nil {
		t.Fatal(err)
	}
	adaGrup := false
	for _, l := range pn {
		if len(l["TreatyGroupList"].([]map[string]any)) > 0 {
			adaGrup = true
		}
		for _, k := range []string{"MDPList", "PremiumEarnedList", "EgnpiTotalList", "MDPMinList", "Reinstatement_List"} {
			if _, ok := l[k].([]map[string]any); !ok {
				t.Errorf("%s bukan larik di kontrak %s", k, np)
			}
		}
	}
	if !adaGrup {
		t.Errorf("kontrak %s: TreatyGroupList kosong", np)
	}
	akar, err := g.BacaLimitsAkarPendaratan(ctx, np)
	if err != nil {
		t.Fatal(err)
	}
	if len(akar.LimitSummaryList) == 0 {
		t.Errorf("kontrak %s: LimitSummaryList kosong", np)
	}
	t.Logf("prop %s · non-prop %s · ringkasan %d · total IOO %d", prop, np, len(akar.LimitSummaryList), len(akar.Total["TotalLimitIOONP"]))
}

// ⭐ Sumber Achievement: baris produksi kontrak dan kurs ke IDR terbaru.
func TestAchievementDanKursKeIDR(t *testing.T) {
	g, ctx := gudangBaca(t)
	sdb, sk := bukaOracle(t)
	var id string
	if err := sdb.QueryRowContext(ctx, `SELECT MIN(SUBSTR(NOOFFER,1,7)) FROM `+sk+`.ACHIEVEMENT`).Scan(&id); err != nil {
		t.Fatal(err)
	}
	b, err := g.BacaAchievement(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatalf("kontrak %s: nol baris ACHIEVEMENT", id)
	}
	k, err := g.BacaKursKeIDR(ctx, []string{"10001"})
	if err != nil {
		t.Fatal(err)
	}
	if k["10001"] == "" {
		t.Error("kurs USD (10001) kosong")
	}
	t.Logf("kontrak %s: %d baris; USD→IDR %s", id, len(b), k["10001"])
}
