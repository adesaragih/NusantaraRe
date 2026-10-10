//go:build db

package repository_test

// Tab Share Non-Prop atas POOLDATA — BACA SAJA, nol transaksi.

import (
	"fmt"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/repository"
)

// Kontrak 1001270: enam layer, RNM Share 40, Spreading Type `2022 QS 145M
// TRT`, mulai 20230401 (terukur 7 Oktober 2026).
const kontrakShare = "1001270"

func TestPohonShareDariPendaratan(t *testing.T) {
	g, ctx := gudangBaca(t)
	sp, err := g.BacaSharePendaratan(ctx, kontrakShare)
	if err != nil {
		t.Fatal(err)
	}
	if len(sp.Share) != 6 {
		t.Fatalf("%d baris Share, mau 6", len(sp.Share))
	}
	for i, b := range sp.Share {
		if b["SpreadingTypeXOL"] != "2022 QS 145M TRT" || b["RNMShare"] != "40" {
			t.Errorf("baris %d: %v / %v", i, b["SpreadingTypeXOL"], b["RNMShare"])
		}
		for _, k := range []string{"SpreadingListXOL", "DeductionList", "GrossPremiumList", "NetPremiumList"} {
			if _, ok := b[k].([]map[string]any); !ok {
				t.Errorf("baris %d: %s bukan larik (%T)", i, k, b[k])
			}
		}
		if n := len(b["SpreadingListXOL"].([]map[string]any)); n != 2 {
			t.Errorf("baris %d: %d baris spreading, mau 2 (QS OR / QS R/I)", i, n)
		}
	}
	if sp.ShareReins == nil || sp.ShareFacultativeReinsurers == nil || sp.FacultativeShareList == nil {
		t.Error("larik kosong harus larik, bukan nil")
	}
}

func TestIndukDanAnakSpreadingDariSusunanMaster(t *testing.T) {
	g, ctx := gudangBaca(t)
	pohon, err := g.BacaPohonLimitsPendaratan(ctx, kontrakShare)
	if err != nil || len(pohon) == 0 {
		t.Fatalf("pohon %d, %v", len(pohon), err)
	}
	grup := pohon[0]["TreatyGroupList"].([]map[string]any)
	if len(grup) == 0 {
		t.Fatal("layer 1 tanpa Treaty Group")
	}
	idGrup := fmt.Sprint(grup[0]["TreatyGroupID"])
	induk, err := g.BacaIndukSpreading(ctx, idGrup, "20230401", "")
	if err != nil {
		t.Fatal(err)
	}
	var idInduk, tahun string
	for _, p := range induk {
		if p.ParentReinsTypeID != "00" {
			t.Errorf("induk ber-ParentReinsTypeID %q", p.ParentReinsTypeID)
		}
		if p.ReinsTypeName == "2022 QS 145M TRT" {
			idInduk, tahun = p.ReinsTypeID, p.TreatyYearID
		}
	}
	if idInduk != "10227" {
		t.Fatalf("induk `2022 QS 145M TRT` tidak ada / ID %q di %+v", idInduk, induk)
	}
	anak, err := g.BacaAnakSpreading(ctx, idInduk, tahun)
	if err != nil {
		t.Fatal(err)
	}
	pct := map[string]string{}
	for _, a := range anak {
		pct[a.ReinsTypeName] = a.Pct
	}
	if pct["QS (OR)"] != "40" || pct["QS (R/I)"] != "60" {
		t.Errorf("anak %+v", anak)
	}
}

func TestReasuradurShareTanpaSaringanChildCount(t *testing.T) {
	g, ctx := gudangBaca(t)
	d, err := g.BacaDaftarReasuradurShare(ctx)
	if err != nil {
		t.Fatal(err)
	}
	cedant, err := g.BacaDaftarCedant(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// Tanpa `CHILDCOUNT = '0'` daftarnya MELIPUTI pemilih Ceding.
	if len(d) < len(cedant) {
		t.Errorf("reasuradur %d < cedant %d", len(d), len(cedant))
	}
	for i := 1; i < len(d); i++ {
		if d[i-1].ID < d[i].ID {
			t.Fatalf("urutan bukan ID menurun: %s lalu %s", d[i-1].ID, d[i].ID)
		}
	}
}

// ⭐ Filter tanpa `pyUseNullIfEmpty` DILEWATI saat parameternya kosong:
// grup kosong → induk SEMUA grup; 10246 (`2025 XOL TRT`) hanya dikeluarkan
// bila pemanggil mengirimnya (kedua dropdown).
func TestIndukSpreadingFilterKosongDilewati(t *testing.T) {
	g, ctx := gudangBaca(t)
	satu, err := g.BacaIndukSpreading(ctx, "10002", "20250301", "")
	if err != nil {
		t.Fatal(err)
	}
	semua, err := g.BacaIndukSpreading(ctx, "", "20250301", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(semua) <= len(satu) {
		t.Errorf("grup kosong %d induk, grup 10002 %d — mau lebih banyak", len(semua), len(satu))
	}
	ada := func(d []models.SusunanSpreading, id string) bool {
		for _, p := range d {
			if p.ReinsTypeID == id {
				return true
			}
		}
		return false
	}
	if !ada(semua, "10246") {
		t.Fatal("10246 (2025 XOL TRT) mestinya ada tanpa filter H")
	}
	dropdown, err := g.BacaIndukSpreading(ctx, "", "20250301", repository.IndukDikecualikanDropdown)
	if err != nil {
		t.Fatal(err)
	}
	if ada(dropdown, "10246") || len(dropdown) != len(semua)-1 {
		t.Errorf("dropdown %d, semua %d", len(dropdown), len(semua))
	}
}

// ⭐ `TREATYINDETAIL` — nilai akar yang `SaveTreatyInDetail_Act` tulis.
// Kontrak 1000402: RNM Share 10 (terukur 7 Oktober 2026).
func TestShareDetailWarisanDariTreatyInDetail(t *testing.T) {
	g, ctx := gudangBaca(t)
	d, err := g.BacaShareDetailWarisan(ctx, "1000402")
	if err != nil {
		t.Fatal(err)
	}
	if d.RNMShare != "10" {
		t.Errorf("RNM_SHARE %q, mau 10", d.RNMShare)
	}
	kosong, err := g.BacaShareDetailWarisan(ctx, "TIDAK-ADA")
	if err != nil || kosong.RNMShare != "" {
		t.Errorf("kontrak tak dikenal: %+v %v", kosong, err)
	}
}

// ⭐ Kolom akar Share (migrasi 445) dibaca TOLERAN: sebelum migrasinya
// dijalankan, jawabannya peta kosong — bukan galat yang menggagalkan
// pembukaan kontrak.
func TestShareAkarRevisiToleranSebelumMigrasi(t *testing.T) {
	g, ctx := gudangBaca(t)
	akar, err := g.BacaShareAkarRevisi(ctx, kontrakShare)
	if err != nil {
		t.Fatalf("mestinya toleran: %v", err)
	}
	t.Logf("kolom akar terbaca: %v", akar)
}
