package services

import (
	"testing"

	"nusantarare/modul/treatyin/backend/models"
)

// sumberGrupUji — RD induk BERSARING Treaty Group (kosong = semua grup) dan
// RD anak BERSARING TreatyYearID (kosong = SEMUA tahun, seperti Pega/RD).
type sumberGrupUji struct {
	induk map[string][]models.SusunanSpreading // grup → induk
	anak  []models.SusunanSpreading            // seluruh anak, semua tahun
}

func (s *sumberGrupUji) Induk(grup, _ string) []models.SusunanSpreading {
	if grup != "" {
		return s.induk[grup]
	}
	out := []models.SusunanSpreading{}
	for _, xs := range s.induk {
		out = append(out, xs...)
	}
	return out
}

func (s *sumberGrupUji) AnakProp(_, _, _, induk, tahunID string) []models.SusunanSpreading {
	out := []models.SusunanSpreading{}
	for _, a := range s.anak {
		if a.ParentReinsTypeID == induk && (tahunID == "" || a.TreatyYearID == tahunID) {
			out = append(out, a)
		}
	}
	return out
}

// Master DEV terukur 9 Oktober 2026: `ORS` (10007) punya 24 anak, satu per
// Treaty Group per tahun 2017–2025; di sini tiga tahun cukup.
func sumberORS() *sumberGrupUji {
	return &sumberGrupUji{
		induk: map[string][]models.SusunanSpreading{
			"10016": {{ReinsTypeID: "10007", ReinsTypeName: "ORS", ParentReinsTypeID: "00", TreatyYearID: "1000675"}},
		},
		anak: []models.SusunanSpreading{
			{ReinsTypeID: "10007", ReinsTypeName: "ORS", ParentReinsTypeID: "10007", TreatyYearID: "1000598", Pct: "100"},
			{ReinsTypeID: "10007", ReinsTypeName: "ORS", ParentReinsTypeID: "10007", TreatyYearID: "1000599", Pct: "100"},
			{ReinsTypeID: "10007", ReinsTypeName: "ORS", ParentReinsTypeID: "10007", TreatyYearID: "1000675", Pct: "100"},
		},
	}
}

func detailORS(pecahan []any) map[string]any {
	r := map[string]any{"ReinsTypeID": "10007", "Pct": "25"}
	if pecahan != nil {
		r["BreakDownSprdList"] = pecahan
	}
	return map[string]any{
		"TreatyGroupID":   "10055", // grup Detail ≠ grup induk ORS di master
		"RNMShare":        "25",
		"SpreadingTypeID": "",
		"RNMShareList":    []any{map[string]any{"Currency": "IDR", "Value": "56250000"}},
		"SpreadingList":   []any{r},
	}
}

// ⛔ Laporan pemakai: ▾ spread `ORS` berisi puluhan baris `ORS`. Induk tidak
// ada di grup Detail → dicari tanpa grup → SATU tahun, satu pecahan.
func TestPecahanORSTidakMengambilSemuaTahun(t *testing.T) {
	d := detailORS(nil)
	r := larikSimpul(d, "SpreadingList")[0]
	pecahan, tahunID := pecahanSebar(sumberORS(), "20260401", d, r, "")
	if tahunID != "1000675" {
		t.Fatalf("TreatyYearID = %q, mau 1000675 (induk tanpa saringan grup)", tahunID)
	}
	if len(pecahan) != 1 {
		t.Fatalf("pecahan %d baris, mau 1: %v", len(pecahan), pecahan)
	}
	if pecahan[0]["ReinsName"] != "ORS" || pecahan[0]["Amount"] != "56250000" {
		t.Errorf("pecahan = %v", pecahan[0])
	}
}

// Induk tidak ditemukan sama sekali → NOL pecahan, bukan anak semua tahun.
func TestPecahanTanpaIndukKosong(t *testing.T) {
	src := sumberORS()
	src.induk = map[string][]models.SusunanSpreading{}
	d := detailORS(nil)
	pecahan, _ := pecahanSebar(src, "20260401", d, larikSimpul(d, "SpreadingList")[0], "")
	if len(pecahan) != 0 {
		t.Fatalf("pecahan %d baris, mau 0", len(pecahan))
	}
}

// Kontrak dibuka: pecahan manual yang KOSONG dilengkapi; yang tersimpan dan
// spreading lama (Spreading Type terisi) tidak disentuh.
func TestLengkapiPecahanHanyaYangKosong(t *testing.T) {
	kosong := detailORS(nil)
	tersimpan := detailORS([]any{map[string]any{"ReinsName": "TERSIMPAN"}})
	lama := detailORS(nil)
	lama["SpreadingTypeID"] = "10227"
	limits := []map[string]any{{"Detail": []any{kosong, tersimpan, lama}}}
	LengkapiPecahanSpreading(sumberORS(), "20260401", limits)

	if n := len(larikSimpul(larikSimpul(kosong, "SpreadingList")[0], "BreakDownSprdList")); n != 1 {
		t.Errorf("Detail kosong: %d pecahan, mau 1", n)
	}
	if p := larikSimpul(larikSimpul(tersimpan, "SpreadingList")[0], "BreakDownSprdList"); len(p) != 1 || p[0]["ReinsName"] != "TERSIMPAN" {
		t.Errorf("pecahan tersimpan berubah: %v", p)
	}
	if _, ada := larikSimpul(lama, "SpreadingList")[0]["BreakDownSprdList"]; ada {
		t.Error("spreading lama ikut dilengkapi")
	}
}

// ⛔ Penjaga keras: master yang (salah) memberi anak kembar dalam SATU tahun
// tetap menghasilkan satu baris per tipe per mata uang — ORS maksimal 1.
func TestPecahanSatuBarisPerTipe(t *testing.T) {
	src := sumberORS()
	kembar := src.anak[2]
	src.anak = append(src.anak, kembar, kembar)
	d := detailORS(nil)
	pecahan, _ := pecahanSebar(src, "20260401", d, larikSimpul(d, "SpreadingList")[0], "")
	if len(pecahan) != 1 {
		t.Fatalf("pecahan %d baris, mau 1 ORS", len(pecahan))
	}
}
