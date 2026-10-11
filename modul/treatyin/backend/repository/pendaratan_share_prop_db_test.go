//go:build db

package repository_test

import "testing"

// TestAnakSpreadingPropLimaFilter - RD `Limit_MstTrt_RD` cabang Prop atas
// `PROPORTIONALARRG` (master, baca saja): parameter kosong dilewati, dan
// `TREATYYEAR` ikut terbaca (dipakai `FetchQSfromMaster` langkah 6.1/7).
func TestAnakSpreadingPropLimaFilter(t *testing.T) {
	g, ctx := gudangBaca(t)
	induk, err := g.BacaIndukSpreading(ctx, "", "20250301", "")
	if err != nil {
		t.Fatalf("induk: %v", err)
	}
	if len(induk) == 0 {
		t.Skip("PROPORTIONALARRG tanpa induk untuk 2025")
	}
	p := induk[0]
	if p.TreatyYear == "" {
		t.Errorf("TreatyYear induk kosong: %+v", p)
	}
	semua, err := g.BacaAnakSpreadingProp(ctx, "", "", "", p.ReinsTypeID, "")
	if err != nil {
		t.Fatalf("anak: %v", err)
	}
	bertahun, err := g.BacaAnakSpreadingProp(ctx, p.TreatyYear, "", "10001", p.ReinsTypeID, p.TreatyYearID)
	if err != nil {
		t.Fatalf("anak bertahun: %v", err)
	}
	if len(bertahun) > len(semua) {
		t.Errorf("filter tambahan menambah baris: %d > %d", len(bertahun), len(semua))
	}
	for _, a := range bertahun {
		if a.ParentReinsTypeID != p.ReinsTypeID || a.TreatyYearID != p.TreatyYearID {
			t.Errorf("baris di luar filter: %+v", a)
		}
	}
	t.Logf("induk %s (%s, tahun %s): anak %d tanpa filter tahun, %d bertahun", p.ReinsTypeID, p.ReinsTypeName, p.TreatyYear, len(semua), len(bertahun))
	if kosong, err := g.BacaAnakSpreadingProp(ctx, "", "", "", "", ""); err != nil || len(kosong) != 0 {
		t.Errorf("tanpa induk harus nol baris: %d %v", len(kosong), err)
	}
}
