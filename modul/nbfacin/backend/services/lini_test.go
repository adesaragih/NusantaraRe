package services

import (
	"context"
	"reflect"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// Tiket 73 P0 - gerbang Class of Business.

// TestBusinessTypeDari - BusinessType_DeT first match (tabel Go = urai XML, 36 baris, sensus dua arah sesi c3). Instrumen: dua
// sampel nyata sesi 9d (003/30 -> Aneka, 006/22 -> FireStyle1).
func TestBusinessTypeDari(t *testing.T) {
	for _, c := range []struct{ gp, old, mau string }{
		{"003", "30", "Aneka"}, {"006", "22", "FireStyle1"}, // sampel nyata
		{"003", "15", "Liability"}, {"003", "E9", "Liability"}, {"003", "28", "FireStyle1"}, {"003", "B3", "LandRig"},
		{"003", "36", "Bonding"}, {"003", "C1", "BondingKBG"}, {"003", "", "Aneka"},
		{"006", "78", "OilGas"}, {"006", "59", "FireStyle2"}, {"006", "99", "Fire"}, {"006", "", "Fire"},
		{"002", "03", "PA"}, {"002", "04", "UNKNOWN"}, {"004", "x", "MarineCargo"}, {"009", "L21", "Life"}, {"009", "L22", "UNKNOWN"},
		{"008", "", "Medicare"}, {"001", "", "Medicare"}, {"", "", "UNKNOWN"}, {"3", "15", "UNKNOWN"},
	} {
		if got := businessTypeDari(c.gp, c.old); got != c.mau {
			t.Errorf("%s/%s: %q mau %q", c.gp, c.old, got, c.mau)
		}
	}
	if len(tabelBusinessType) != 36 {
		t.Errorf("baris %d", len(tabelBusinessType))
	}
}

// layananLini - kasus "K" ber-Class of Business `nama` (kosong = "UJI LIABILITY"), baris BUSINESS NOTE + group case `bisnis`
// (GroupPanel "" = NULL).
func layananLini(bisnis []models.BarisBisnisLini, nama ...string) *Service {
	cob := "UJI LIABILITY"
	if len(nama) > 0 {
		cob = nama[0]
	}
	k := &models.Kasus{CaseID: "K", Opportunity: models.Opportunity{ClassOfBusiness: cob, GroupBusinessID: "3"}}
	return Baru(nil).DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"K": k}}).
		DenganTableOfLimit(tolTiruan{bisnis: bisnis})
}

// satuBisnis - satu baris BUSINESS (ID 10999 kecuali disebut).
func satuBisnis(gp, old string, id ...string) []models.BarisBisnisLini {
	b := models.BarisBisnisLini{ID: "10999", OldID: old, GroupPanel: gp}
	if len(id) > 0 {
		b.ID = id[0]
	}
	return []models.BarisBisnisLini{b}
}

// tanpaTampilan - LiniKasus tanpa TampilanAneka (diuji terpisah di TestTampilanAneka).
func tanpaTampilan(l LiniKasus) LiniKasus {
	l.TampilanAneka, l.JenisCsv = TampilanAneka{}, "" // JenisCsv diuji TestJenisCsvKasus
	return l
}

// TestLiniKasus - predikat When korpus atas BusinessType / BusinessCode / BusinessName case.
func TestLiniKasus(t *testing.T) {
	ctx := context.Background()
	ganda := append(satuBisnis("003", "15", "1"), satuBisnis("003", "16", "2")...)
	for nama, c := range map[string]struct {
		bisnis []models.BarisBisnisLini
		mau    LiniKasus
	}{
		"liability":      {satuBisnis("003", "15", "10021"), LiniKasus{BusinessType: "Liability", Aneka: true, AnekaNotBonding: true, Liability: true}},
		"aneka":          {satuBisnis("003", "30"), LiniKasus{BusinessType: "Aneka", Aneka: true, AnekaNotBonding: true}},
		"bonding":        {satuBisnis("003", "36"), LiniKasus{BusinessType: "Bonding", Aneka: true, Bonding: true}},
		"golf":           {satuBisnis("003", "27"), LiniKasus{BusinessType: "GolfInsurance", Golf: true}},
		"fire":           {satuBisnis("006", "99"), LiniKasus{BusinessType: "Fire", Fire: true}},
		"marine cargo":   {satuBisnis("004", "24"), LiniKasus{BusinessType: "MarineCargo", MarineCargo: true}}, // tiket 74 P0
		"fire style 1":   {satuBisnis("006", "22"), LiniKasus{BusinessType: "FireStyle1", Fire: true}},
		"BUSINESS ganda": {ganda, LiniKasus{BusinessType: "UNKNOWN"}},
		// Keputusan work owner 09-10-2026 (saringan `GROUPPANEL IS NOT NULL` Pega GetOldIDBusiness_SQL): baris ber-GROUPPANEL
		// NULL tidak ikut dihitung - dua baris NOTE + group sama, satu NULL -> baris ber-GROUPPANEL (bukan UNKNOWN).
		"ganda satu GROUPPANEL NULL": {append(satuBisnis("", "16", "1"), satuBisnis("003", "15", "10021")...),
			LiniKasus{BusinessType: "Liability", Aneka: true, AnekaNotBonding: true, Liability: true}},
		"ganda NULL lebih dulu, Fire": {append(satuBisnis("", "30", "1"), satuBisnis("006", "22", "2")...),
			LiniKasus{BusinessType: "FireStyle1", Fire: true}},
		"satu GROUPPANEL NULL saja": {satuBisnis("", "15"), LiniKasus{BusinessType: "UNKNOWN"}},
		"tanpa baris":               {nil, LiniKasus{BusinessType: "UNKNOWN"}},
	} {
		if got := tanpaTampilan(layananLini(c.bisnis).LiniKasus(ctx, "K")); !reflect.DeepEqual(got, c.mau) {
			t.Errorf("%s: %+v mau %+v", nama, got, c.mau)
		}
	}
	// BusinessName = Class of Business: IsBillboardNeon (`BusinessName = "BILLBOARD/NEON SIGN"`) -> IsLiability.
	if l := layananLini(satuBisnis("003", "30"), "BILLBOARD/NEON SIGN").LiniKasus(ctx, "K"); !l.Liability ||
		l.BusinessType != "Aneka" {
		t.Errorf("billboard: %+v", l)
	}
	// BusinessCode = BUSINESS.ID baris ber-GROUPPANEL: IsLiability E (`BusinessCode = "10048"`).
	if l := layananLini(append(satuBisnis("", "30", "10001"), satuBisnis("003", "30", "10048")...)).LiniKasus(ctx, "K"); !l.Liability {
		t.Errorf("kode dari baris ber-GROUPPANEL: %+v", l)
	}
	if l := Baru(nil).LiniKasus(ctx, "K"); !reflect.DeepEqual(l, LiniKasus{}) {
		t.Errorf("tanpa DB: %+v", l)
	}
	if l := layananLini(nil).LiniKasus(ctx, "K-404"); !reflect.DeepEqual(l, LiniKasus{}) {
		t.Errorf("kasus tidak ada: %+v", l)
	}
}

// TestTampilanAneka - tiket 73 P1: kontainer C1..C9 + sakelar, dari predikat When korpus (kontrak 9d: Liability C1, MarineHull
// C2, HE C5+C6, Boiler C1+C4+C5, Golf C1, FIRE []).
func TestTampilanAneka(t *testing.T) {
	ctx := context.Background()
	for nama, c := range map[string]struct {
		gp, old string
		mau     TampilanAneka
	}{
		"liability":   {"003", "15", TampilanAneka{Kontainer: []string{"C1"}, TopRisk: true, JumlahSatuan: true}},
		"marine hull": {"003", "05", TampilanAneka{Kontainer: []string{"C2"}, TopRisk: true, JumlahSatuan: true, CurrencyCoverage: true}},
		"HE": {"003", "20", TampilanAneka{Kontainer: []string{"C5", "C6"}, HE: true, TopRisk: true, JumlahSatuan: true,
			CurrencyCoverage: true}},
		// tiket 73 P2: MBD (OldId 10) - sel Currency dan blok indemnity form coverage.
		"MBD": {"003", "10", TampilanAneka{Kontainer: []string{"C4", "C5"}, KolomTahunJumlah: true, Engineering: true, TopRisk: true,
			JumlahSatuan: true, CurrencyCoverage: true, MBD: true}},
		"boiler": {"003", "11", TampilanAneka{Kontainer: []string{"C1", "C4", "C5"}, KolomTahunJumlah: true,
			JumlahSatuan: true, CurrencyCoverage: true}},
		"golf": {"003", "27", TampilanAneka{Kontainer: []string{"C1"}, Golf: true, JumlahSatuan: true}},
		"fire": {"006", "99", TampilanAneka{Kontainer: []string{}}},
	} {
		got := layananLini(satuBisnis(c.gp, c.old), "UJI KELAS").LiniKasus(ctx, "K").TampilanAneka
		if !reflect.DeepEqual(got, c.mau) {
			t.Errorf("%s: %+v mau %+v", nama, got, c.mau)
		}
	}
	// Sakelar ber-BusinessCode (BUSINESS.ID) di bawah gerbang Liability: IsCreditBriguna `10235`, IsEngineering E `10032`,
	// IsObjectSectionAneka B `15` (K73-4: kode 2 karakter ditiru apa adanya), BusinessCode 20 -> jumlahSatuan false.
	for kode, mau := range map[string]TampilanAneka{
		"10235": {Kontainer: []string{"C1"}, Briguna: true, TopRisk: true, JumlahSatuan: true},
		"10032": {Kontainer: []string{"C1"}, Engineering: true, TopRisk: true, JumlahSatuan: true},
		"15":    {Kontainer: []string{"C1"}, SectionLol: true, TopRisk: true, JumlahSatuan: true},
		"20":    {Kontainer: []string{"C1"}, TopRisk: true},
	} {
		got := layananLini(satuBisnis("003", "15", kode), "UJI KELAS").LiniKasus(ctx, "K").TampilanAneka
		if !reflect.DeepEqual(got, mau) {
			t.Errorf("kode %s: %+v mau %+v", kode, got, mau)
		}
	}
	for kode, mau := range map[string]bool{"20": false, " 20 ": false, "20.0": false, "10020": true, "": true} {
		if got := !samaDuaPuluh(kode); got != mau {
			t.Errorf("jumlahSatuan %q: %v", kode, got)
		}
	}
}

// TestJenisCsvKasus - tiket 73 P6: template CSV objek = SetDownloadMasterCSV langkah 2-6 berurutan (yang terakhir benar menang):
// Marine Hull dan CPM (keduanya IsAneka) menimpa "aneka"; Golf bukan IsAneka -> ""; tanpa BUSINESS -> "".
func TestJenisCsvKasus(t *testing.T) {
	ctx := context.Background()
	for nama, c := range map[string]struct{ gp, old, mau string }{
		"marine hull": {"003", "05", "marinehull"}, "marine hull 55": {"003", "55", "marinehull"}, "cpm": {"003", "12", "cpm"},
		"liability": {"003", "15", "aneka"}, "HE": {"003", "20", "aneka"}, "golf": {"003", "27", ""}, "fire": {"006", "24", "fire"},
		"cargo": {"004", "", "cargo"},
	} {
		if got := layananLini(satuBisnis(c.gp, c.old), "UJI KELAS").LiniKasus(ctx, "K").JenisCsv; got != c.mau {
			t.Errorf("%s: %q mau %q", nama, got, c.mau)
		}
	}
	if got := Baru(nil).LiniKasus(ctx, "K").JenisCsv; got != "" {
		t.Errorf("tanpa pembaca: %q", got)
	}
}
