package services_test

// Seam 4 - lapis C: salinan daftar dari polis lama, penanda baris warisan,
// penggabungan spreading. Tiket E10, ditambah temuan korpus 01-10-2026 yang
// dicatat di services/lapisc.go.
//
// Dibaca sesudah: beforeimage_test.go, services/lapisc.go.

import (
	"errors"
	"testing"

	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/services"
)

func covSpread(t *testing.T, baris ...models.Spreading) models.Coverage {
	return models.Coverage{SpreadingList: baris}
}

func siapkanC(t *testing.T, lama models.OfferFacIn, jenis models.JenisPenyesuaian, p services.Predikat) models.OfferFacIn {
	t.Helper()
	k := kasusBaru(t, "2026-01-02 00:00")
	k.OfferFacIn.QuotationData.Type = jenis
	h, err := services.PrepareBeforeImage(services.MasukanBeforeImage{Kasus: k, PolisLama: &lama, Predikat: p})
	if err != nil {
		t.Fatalf("PrepareBeforeImage: %v", err)
	}
	return h.Kasus.OfferFacIn
}

func harusOld(t *testing.T, label, nilai string) {
	t.Helper()
	if nilai != models.NilaiOld {
		t.Errorf("%s = %q, mau %q", label, nilai, models.NilaiOld)
	}
}

// TestLapisCTujuhVarianMenyalinDanMenandai - tiap varian menyalin daftarnya
// dari OldData lalu menandai setiap tingkat sesuai kedalaman model datanya.
func TestLapisCTujuhVarianMenyalinDanMenandai(t *testing.T) {
	sp := func() []models.Spreading { return []models.Spreading{spread(t, "UJI-T1", "1", "1", "100")} }

	t.Run("FIRE", func(t *testing.T) {
		lama := models.OfferFacIn{LocationList: []models.Lokasi{{Property: models.PropertiLokasi{
			PropertyItemList: []models.ItemProperti{{CoverageList: []models.Coverage{{SpreadingList: sp()}}}},
		}}}}
		o := siapkanC(t, lama, "1", services.Predikat{IsFire: true})
		l := o.LocationList[0]
		harusOld(t, "Lokasi.IsOldData", l.IsOldData)
		harusOld(t, "PropertyItem.IsOldData", l.Property.PropertyItemList[0].IsOldData)
		harusOld(t, "Coverage.IsOldData", l.Property.PropertyItemList[0].CoverageList[0].IsOldData)
		harusOld(t, "Spreading.IsOldData", l.Property.PropertyItemList[0].CoverageList[0].SpreadingList[0].IsOldData)
	})
	t.Run("MC", func(t *testing.T) {
		lama := models.OfferFacIn{CargoList: []models.Kargo{{CoverageList: []models.Coverage{{SpreadingList: sp()}}}}}
		o := siapkanC(t, lama, "1", services.Predikat{IsMarineCargo: true})
		harusOld(t, "Kargo.FlagOldData", o.CargoList[0].FlagOldData)
		harusOld(t, "Coverage.IsOldData", o.CargoList[0].CoverageList[0].IsOldData)
		harusOld(t, "Spreading.IsOldData", o.CargoList[0].CoverageList[0].SpreadingList[0].IsOldData)
	})
	t.Run("Aneka", func(t *testing.T) {
		lama := models.OfferFacIn{LocationList: []models.Lokasi{{Property: models.PropertiLokasi{RiskLocation: models.RisikoLokasi{
			OccupationList: []models.Okupasi{{AnekaList: []models.Aneka{{CoverageList: []models.Coverage{{SpreadingList: sp()}}}}}},
		}}}}}
		o := siapkanC(t, lama, "1", services.Predikat{IsAneka: true})
		okp := o.LocationList[0].Property.RiskLocation.OccupationList[0]
		harusOld(t, "Lokasi.IsOldData", o.LocationList[0].IsOldData)
		harusOld(t, "Occupation.IsOldData", okp.IsOldData)
		harusOld(t, "Aneka.IsOldData", okp.AnekaList[0].IsOldData)
		harusOld(t, "Coverage.IsOldData", okp.AnekaList[0].CoverageList[0].IsOldData)
		harusOld(t, "Spreading.IsOldData", okp.AnekaList[0].CoverageList[0].SpreadingList[0].IsOldData)
	})
	t.Run("MBU", func(t *testing.T) {
		lama := models.OfferFacIn{VehicleList: []models.Kendaraan{{CoverageList: []models.Coverage{{
			SpreadingList: sp(), AdditionalCoverage: []models.Coverage{{}},
		}}}}}
		o := siapkanC(t, lama, "1", services.Predikat{IsMBU: true})
		c := o.VehicleList[0].CoverageList[0]
		harusOld(t, "Kendaraan.FlagOldData", o.VehicleList[0].FlagOldData)
		harusOld(t, "Coverage.IsOldData", c.IsOldData)
		harusOld(t, "AdditionalCoverage.IsOldData", c.AdditionalCoverage[0].IsOldData)
		harusOld(t, "Spreading.IsOldData", c.SpreadingList[0].IsOldData)
	})
	t.Run("LIFE", func(t *testing.T) {
		lama := models.OfferFacIn{PersonList: []models.Orang{{CoverageList: []models.Coverage{{SpreadingList: sp()}}}}}
		o := siapkanC(t, lama, "1", services.Predikat{IsLife: true})
		harusOld(t, "Orang.FlagOldData", o.PersonList[0].FlagOldData)
		harusOld(t, "Coverage.IsOldData", o.PersonList[0].CoverageList[0].IsOldData)
		harusOld(t, "Spreading.IsOldData", o.PersonList[0].CoverageList[0].SpreadingList[0].IsOldData)
	})
	t.Run("PA", func(t *testing.T) {
		lama := models.OfferFacIn{PersonList: []models.Orang{{ASMCoverage: []models.Coverage{{SpreadingList: sp()}}}}}
		o := siapkanC(t, lama, "1", services.Predikat{IsPA: true})
		harusOld(t, "Orang.FlagOldData", o.PersonList[0].FlagOldData)
		harusOld(t, "ASMCoverage.IsOldData", o.PersonList[0].ASMCoverage[0].IsOldData)
		harusOld(t, "Spreading.IsOldData", o.PersonList[0].ASMCoverage[0].SpreadingList[0].IsOldData)
	})
	t.Run("Golf", func(t *testing.T) {
		lama := models.OfferFacIn{LocationList: []models.Lokasi{{Property: models.PropertiLokasi{RiskLocation: models.RisikoLokasi{
			AnekaList: []models.Aneka{{CoverageList: []models.Coverage{{SpreadingList: sp()}}}},
		}}}}}
		o := siapkanC(t, lama, "1", services.Predikat{IsGolfInsurance: true})
		an := o.LocationList[0].Property.RiskLocation.AnekaList[0]
		harusOld(t, "Lokasi.IsOldData", o.LocationList[0].IsOldData)
		harusOld(t, "Aneka.IsOldData", an.IsOldData)
		harusOld(t, "Coverage.IsOldData", an.CoverageList[0].IsOldData)
		harusOld(t, "Spreading.IsOldData", an.CoverageList[0].SpreadingList[0].IsOldData)
	})
}

// TestLapisCPenandaProrataHanyaFire - E10: hanya varian kebakaran menyetel
// IsProRate; enam lainnya tidak.
func TestLapisCPenandaProrataHanyaFire(t *testing.T) {
	lama := models.OfferFacIn{
		LocationList: []models.Lokasi{{}}, CargoList: []models.Kargo{{}},
		VehicleList: []models.Kendaraan{{}}, PersonList: []models.Orang{{}},
	}
	for nama, p := range map[string]services.Predikat{
		"MC": {IsMarineCargo: true}, "Aneka": {IsAneka: true}, "MBU": {IsMBU: true},
		"LIFE": {IsLife: true}, "PA": {IsPA: true}, "Golf": {IsGolfInsurance: true},
	} {
		if o := siapkanC(t, lama, "1", p); o.IsProRate != "" {
			t.Errorf("%s menyetel IsProRate %q", nama, o.IsProRate)
		}
	}
	if o := siapkanC(t, lama, "1", services.Predikat{IsFire: true}); o.IsProRate != models.NilaiProrate {
		t.Errorf("FIRE: IsProRate %q", o.IsProRate)
	}
}

// TestLapisCFlagOldDataBukanHanyaLife - temuan korpus: penanda tingkat atas
// FlagOldData dipasang MC, MBU, PA, dan LIFE. Premis kasus uji
// `K046_Life_DuaPenandaOldData` ("hanya LIFE dua penanda") TIDAK didukung
// korpus; test ini mengikuti korpus.
func TestLapisCFlagOldDataBukanHanyaLife(t *testing.T) {
	lama := models.OfferFacIn{CargoList: []models.Kargo{{}}, VehicleList: []models.Kendaraan{{}}, PersonList: []models.Orang{{}}}
	periksa := map[string]struct {
		p     services.Predikat
		ambil func(models.OfferFacIn) string
	}{
		"MC":   {services.Predikat{IsMarineCargo: true}, func(o models.OfferFacIn) string { return o.CargoList[0].FlagOldData }},
		"MBU":  {services.Predikat{IsMBU: true}, func(o models.OfferFacIn) string { return o.VehicleList[0].FlagOldData }},
		"PA":   {services.Predikat{IsPA: true}, func(o models.OfferFacIn) string { return o.PersonList[0].FlagOldData }},
		"LIFE": {services.Predikat{IsLife: true}, func(o models.OfferFacIn) string { return o.PersonList[0].FlagOldData }},
	}
	for nama, u := range periksa {
		harusOld(t, nama+" FlagOldData", u.ambil(siapkanC(t, lama, "1", u.p)))
	}
}

// TestLapisCGabungSpreadingPerTreatyType - baris spreading berjenis treaty
// sama DIGABUNG (TSI, premi, bagian dijumlahkan), berlaku untuk SEMUA jenis
// endorsement karena prakondisi `Type=="4"`-nya nonaktif (P-11).
func TestLapisCGabungSpreadingPerTreatyType(t *testing.T) {
	lama := models.OfferFacIn{CargoList: []models.Kargo{{CoverageList: []models.Coverage{covSpread(t,
		spread(t, "UJI-T1", "100", "10", "30"),
		spread(t, "UJI-T2", "50", "5", "20"),
		spread(t, "UJI-T1", "200", "20", "50"),
	)}}}}
	for _, jenis := range []models.JenisPenyesuaian{"1", services.TypeAdjSpreading} {
		o := siapkanC(t, lama, jenis, services.Predikat{IsMarineCargo: true})
		got := o.CargoList[0].CoverageList[0].SpreadingList
		if len(got) != 2 || got[0].TreatyType != "UJI-T1" || got[1].TreatyType != "UJI-T2" {
			t.Fatalf("jenis %s: hasil gabung %+v", jenis, got)
		}
		samaUang(t, "T1 TSISpreaded", got[0].TSISpreaded, duit(t, "300"))
		samaUang(t, "T1 PremiumSpreaded", got[0].PremiumSpreaded, duit(t, "30"))
		samaRasio(t, "T1 SharePercentage", got[0].SharePercentage, rasio(t, "80"))
		samaUang(t, "T2 TSISpreaded", got[1].TSISpreaded, duit(t, "50"))
		// Polis lama tidak ikut tergabung.
		if n := len(o.OldData.CargoList[0].CoverageList[0].SpreadingList); n != 3 {
			t.Errorf("jenis %s: spreading lapis A %d baris, mau 3", jenis, n)
		}
	}
}

// TestLapisCBuangBagianFireBerbedaDariLainnya - FIRE membuang bagian TEPAT
// nol dan menyimpan bagian negatif; lima lini lain membuang nol DAN negatif.
func TestLapisCBuangBagianFireBerbedaDariLainnya(t *testing.T) {
	baris := func() []models.Spreading {
		return []models.Spreading{
			spread(t, "UJI-NOL", "1", "1", "0"),
			spread(t, "UJI-NEG", "1", "1", "-5"),
			spread(t, "UJI-POS", "1", "1", "5"),
		}
	}
	treaty := func(s []models.Spreading) []string {
		var h []string
		for _, b := range s {
			h = append(h, b.TreatyType)
		}
		return h
	}
	fire := siapkanC(t, models.OfferFacIn{LocationList: []models.Lokasi{{Property: models.PropertiLokasi{
		PropertyItemList: []models.ItemProperti{{CoverageList: []models.Coverage{covSpread(t, baris()...)}}},
	}}}}, "1", services.Predikat{IsFire: true})
	if got := treaty(fire.LocationList[0].Property.PropertyItemList[0].CoverageList[0].SpreadingList); len(got) != 2 || got[0] != "UJI-NEG" || got[1] != "UJI-POS" {
		t.Errorf("FIRE menyisakan %v, mau [UJI-NEG UJI-POS]", got)
	}
	mbu := siapkanC(t, models.OfferFacIn{VehicleList: []models.Kendaraan{{CoverageList: []models.Coverage{covSpread(t, baris()...)}}}},
		"1", services.Predikat{IsMBU: true})
	if got := treaty(mbu.VehicleList[0].CoverageList[0].SpreadingList); len(got) != 1 || got[0] != "UJI-POS" {
		t.Errorf("MBU menyisakan %v, mau [UJI-POS]", got)
	}
}

// TestLapisCLifeTanpaPenggabungan - varian jiwa tidak menggabung spreading.
func TestLapisCLifeTanpaPenggabungan(t *testing.T) {
	lama := models.OfferFacIn{PersonList: []models.Orang{{CoverageList: []models.Coverage{covSpread(t,
		spread(t, "UJI-T1", "1", "1", "0"), spread(t, "UJI-T1", "1", "1", "0"),
	)}}}}
	o := siapkanC(t, lama, "1", services.Predikat{IsLife: true})
	if n := len(o.PersonList[0].CoverageList[0].SpreadingList); n != 2 {
		t.Errorf("LIFE: %d baris spreading, mau 2 (tanpa penggabungan)", n)
	}
}

// TestLapisCSpreadingKosongDitolak - menjumlahkan nilai kosong tidak ditebak.
func TestLapisCSpreadingKosongDitolak(t *testing.T) {
	b := spread(t, "UJI-T1", "1", "1", "1")
	kosongTSI := b
	kosongTSI.TSISpreaded = kosong()
	lama := models.OfferFacIn{CargoList: []models.Kargo{{CoverageList: []models.Coverage{covSpread(t, b, kosongTSI)}}}}
	k := kasusBaru(t, "2026-01-02 00:00")
	_, err := services.PrepareBeforeImage(services.MasukanBeforeImage{Kasus: k, PolisLama: &lama, Predikat: services.Predikat{IsMarineCargo: true}})
	if !errors.Is(err, services.ErrSpreadingTakTerjumlahkan) {
		t.Fatalf("galat %v, mau ErrSpreadingTakTerjumlahkan", err)
	}
}

// TestLapisCMarineCargoMenyalinPolicyData - varian MC menyalin PolicyData utuh,
// mengosongkan nomor endorsement, dan membawa tiga daftar kargo PERTAMA ke
// akar objek kerja.
func TestLapisCMarineCargoMenyalinPolicyData(t *testing.T) {
	lama := models.OfferFacIn{
		PolicyData: models.PolicyData{EndorsementNo: "UJI-EDM-01", OfferingDate: "UJI-OD"},
		CargoList: []models.Kargo{
			{TradingList: []byte(`["t1"]`), GoodList: []byte(`["g1"]`), ConveyanceList: []byte(`["c1"]`)},
			{TradingList: []byte(`["t2"]`)},
		},
	}
	k := kasusBaru(t, "2026-01-02 00:00")
	h, err := services.PrepareBeforeImage(services.MasukanBeforeImage{Kasus: k, PolisLama: &lama, Predikat: services.Predikat{IsMarineCargo: true}})
	if err != nil {
		t.Fatal(err)
	}
	if pd := h.Kasus.OfferFacIn.PolicyData; pd.EndorsementNo != "" || pd.OfferingDate != "UJI-OD" {
		t.Errorf("PolicyData %+v", pd)
	}
	if string(h.Kasus.TradingList) != `["t1"]` || string(h.Kasus.GoodsList) != `["g1"]` || string(h.Kasus.ConveyanceList) != `["c1"]` {
		t.Errorf("daftar akar: %s %s %s", h.Kasus.TradingList, h.Kasus.GoodsList, h.Kasus.ConveyanceList)
	}
}
