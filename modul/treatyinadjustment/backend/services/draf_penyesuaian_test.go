package services_test

// Uji tombol Add Revision / Add Adjustment Premium dan `Choose` —
// `TreatyInEDMSetValue` [1]–[6] sebagai draf di layar.

import (
	"context"
	"errors"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyinadjustment/backend/models"
	"nusantarare/modul/treatyinadjustment/backend/services"
)

func (g *gudangTiruan) DaftarMasterPilihan(_ context.Context, hanyaNonProp bool) ([]models.BarisMasterPilihan, error) {
	g.hanyaNonProp = &hanyaNonProp
	return g.master, nil
}

func (g *gudangTiruan) BacaDokumenMaster(_ context.Context, id string) (models.SisiPenyesuaian, bool, error) {
	d, ada := g.dokumenMaster[id]
	return d, ada, nil
}

func (g *gudangTiruan) AdaRevisi(_ context.Context, id string) (bool, error) {
	return g.adaRevisi[id], nil
}

var pelakuUji = inti.Pelaku{AkunID: "kelvin"}

var kiniUji = time.Date(2026, 10, 7, 8, 30, 15, 0, time.UTC)

func dokumenUji() models.SisiPenyesuaian {
	return models.SisiPenyesuaian{
		Medan: map[string]string{
			"ID": "1000506", "ProportionType": "NonProportional", "Commencement": "20250101",
			"StatusAkseptasi": "Resolve Complete", "Position": "ReasTreatyInSecHead", "Comment": "lama",
		},
		Larik: map[string][]map[string]string{
			"CommentList": {{"Suggest": "komentar lama"}},
			"EGNPI":       {{"Currency": "IDR", "Amount": "1000"}},
		},
		Pohon: map[string][]map[string]any{},
	}
}

func TestDraftRevisiInternalDariMasterTanpaRevisi(t *testing.T) {
	p := services.SusunDraf(dokumenUji(), models.MasukanDraf{ID: "1000506", InternalType: "1", MaterialType: "2"}, false, "kelvin", kiniUji)

	if p.ID != "1000506/R01" || p.IDAsal != "1000506" {
		t.Fatalf("pengenal %q asal %q, mau 1000506/R01 asal 1000506", p.ID, p.IDAsal)
	}
	m := p.Baru.Medan
	harap := map[string]string{
		"ID": "1000506/R01", "OLDID": "1000506", "EDMState": "1", "EDMMaterialType": "2",
		"RevisionDate": "20261007T083015.000 GMT", "EDMEffective": "20250101",
		"ViewState": "0", "Position": "ReasTreatyInAdmin", "IsEditData": "0",
		"PositionUsername": "kelvin", "StatusAkseptasi": "", "Comment": "",
	}
	for k, v := range harap {
		if m[k] != v {
			t.Errorf("Baru.%s = %q, mau %q", k, m[k], v)
		}
	}
	// [3] mengosongkan CommentList, [4] menambah SATU baris.
	k := p.Baru.Larik["CommentList"]
	if len(k) != 1 || k[0]["Suggest"] != "kelvin  Had Created Internal Edit" || k[0]["Date"] != "20261007T083015.000 GMT" || k[0]["OperatorName"] != "kelvin" {
		t.Errorf("CommentList %v", k)
	}
	if _, ada := m["AddendumPremi"]; ada {
		t.Error("AddendumPremi hanya untuk EDMState 3")
	}
}

func TestDraftOldDataSalinanSebelumPengenalBerganti(t *testing.T) {
	p := services.SusunDraf(dokumenUji(), models.MasukanDraf{ID: "1000506", InternalType: "2", MaterialType: "1"}, false, "kelvin", kiniUji)
	// `TreatyInRevisi_post` [1] menyalin TreatyIn → OLDDATA SEBELUM [4]/[5].
	if p.Lama.Medan["ID"] != "1000506" || p.Lama.Medan["EDMState"] != "2" || p.Lama.Medan["Position"] != "ReasTreatyInAdmin" {
		t.Errorf("Lama %v", p.Lama.Medan)
	}
	if got := p.Lama.Larik["CommentList"][0]["Suggest"]; got != "kelvin  Had Created External Addendum" {
		t.Errorf("komentar Lama %q", got)
	}
	// EDMState 2 TIDAK menyetel EDMEffective.
	if _, ada := p.Baru.Medan["EDMEffective"]; ada {
		t.Error("EDMEffective hanya untuk EDMState 1")
	}
	// Salinan, bukan rujukan bersama.
	p.Baru.Larik["CommentList"][0]["Suggest"] = "x"
	if p.Lama.Larik["CommentList"][0]["Suggest"] == "x" {
		t.Error("Lama berbagi peta dengan Baru")
	}
}

func TestDraftPremiMenyalinEGNPIDanAddendumPremi(t *testing.T) {
	asal := dokumenUji()
	p := services.SusunDraf(asal, models.MasukanDraf{ID: "1000506", InternalType: "3", MaterialType: "1"}, false, "kelvin", kiniUji)
	if p.Baru.Medan["AddendumPremi"] != "1" || p.Baru.Medan["EDMMaterialType"] != "1" {
		t.Errorf("medan premi %v", p.Baru.Medan)
	}
	if got := p.Baru.Larik["ActualValue.EGNPI"]; len(got) != 1 || got[0]["Amount"] != "1000" {
		t.Errorf("ActualValue.EGNPI %v", got)
	}
	if got := p.Baru.Larik["CommentList"][0]["Suggest"]; got != "kelvin  Had Created Addendum Premium" {
		t.Errorf("komentar %q", got)
	}
	// Dokumen asal tidak tersentuh.
	if asal.Medan["ID"] != "1000506" || len(asal.Larik["CommentList"]) != 1 || asal.Larik["CommentList"][0]["Suggest"] != "komentar lama" {
		t.Errorf("dokumen asal berubah: %v", asal)
	}
}

func TestDraftViewStateMengikutiRevisionState(t *testing.T) {
	d := dokumenUji()
	d.Medan["RevisionState"] = "1"
	p := services.SusunDraf(d, models.MasukanDraf{ID: "1000506", InternalType: "1"}, false, "kelvin", kiniUji)
	if p.Baru.Medan["ViewState"] != "1" {
		t.Errorf("ViewState %q, mau 1 (TreatyInSetEdit [2])", p.Baru.Medan["ViewState"])
	}
}

func TestIDRevisiBaru(t *testing.T) {
	for _, c := range []struct {
		id   string
		ada  bool
		mau  string
		alas string
	}{
		{"1000506", false, "1000506/R01", "[4] belum ada revisi"},
		{"1000506/R01", true, "1000506/R02", "[5] terukur: R01 → R02"},
		{"1000506/R02", true, "1000506/R03", "[5] terukur: R02 → R03"},
		{"1000506/R09", true, "1000506/R010", "[5] apa adanya: syaratnya n < 10"},
		{"1000506/R10", true, "1000506/R11", "[5] cabang kedua"},
		{"1000506", true, "1000506/R01", "[5] master 7 aksara: n = 0"},
	} {
		if got := services.IDRevisiBaru(c.id, c.ada); got != c.mau {
			t.Errorf("%s: IDRevisiBaru(%q, %v) = %q, mau %q", c.alas, c.id, c.ada, got, c.mau)
		}
	}
}

func TestDrafPenyesuaianLewatLayanan(t *testing.T) {
	g := &gudangTiruan{
		dokumenMaster: map[string]models.SisiPenyesuaian{"1000506/R01": dokumenUji()},
		adaRevisi:     map[string]bool{"1000506/R01": true},
	}
	l := services.LayananDengan(g)
	p, err := l.DrafPenyesuaian(context.Background(), pelakuUji, models.MasukanDraf{ID: " 1000506/R01 ", InternalType: "2", MaterialType: "1"})
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != "1000506/R02" || p.IDAsal != "1000506/R01" || p.Baru.Medan["OLDID"] != "1000506/R01" {
		t.Errorf("draf %s asal %s OLDID %s", p.ID, p.IDAsal, p.Baru.Medan["OLDID"])
	}
}

func TestDrafPenyesuaianMenolakMasukanSalah(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{dokumenMaster: map[string]models.SisiPenyesuaian{}})
	ctx := context.Background()
	if _, err := l.DrafPenyesuaian(ctx, inti.Pelaku{}, models.MasukanDraf{ID: "1"}); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
	if _, err := l.DrafPenyesuaian(ctx, pelakuUji, models.MasukanDraf{ID: "", InternalType: "1"}); !errors.Is(err, services.ErrIDTidakSah) {
		t.Errorf("pengenal kosong: %v", err)
	}
	if _, err := l.DrafPenyesuaian(ctx, pelakuUji, models.MasukanDraf{ID: "1000506", InternalType: "4"}); !errors.Is(err, services.ErrIDTidakSah) {
		t.Errorf("InternalType 4: %v", err)
	}
	if _, err := l.DrafPenyesuaian(ctx, pelakuUji, models.MasukanDraf{ID: "9999999", InternalType: "1"}); !errors.Is(err, services.ErrKontrakTidakAda) {
		t.Errorf("master tak ada: %v", err)
	}
}

func TestDaftarMasterPilihanPerJenis(t *testing.T) {
	g := &gudangTiruan{master: []models.BarisMasterPilihan{{ID: "1000506"}}}
	l := services.LayananDengan(g)
	ctx := context.Background()
	if _, err := l.DaftarMasterPilihan(ctx, pelakuUji, "revisi"); err != nil || g.hanyaNonProp == nil || *g.hanyaNonProp {
		t.Errorf("revisi: err %v, hanyaNonProp %v", err, g.hanyaNonProp)
	}
	if _, err := l.DaftarMasterPilihan(ctx, pelakuUji, "premi"); err != nil || !*g.hanyaNonProp {
		t.Errorf("premi harus NonProportional saja (TreatyLoadMasterJoinEdmXOL): err %v", err)
	}
	if _, err := l.DaftarMasterPilihan(ctx, pelakuUji, "lain"); !errors.Is(err, services.ErrIDTidakSah) {
		t.Errorf("jenis lain: %v", err)
	}
}
