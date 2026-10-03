package services_test

// Seam 4 - lapis B, nilai lama per baris. Tiket E08 (enam lini + cedant), E09
// (guard K-046 dan A.5), E11 (tiga gerbang keluar).
//
// Dibaca sesudah: beforeimage_test.go, services/nilailama.go.

import (
	"testing"

	"nusantarare/inti/backend/uang"
	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/services"
)

// coverageLama - coverage polis lama dengan seluruh sumber terisi, nilai
// berbeda-beda supaya salah pasang ketahuan.
func coverageLama(t *testing.T) models.Coverage {
	return models.Coverage{
		TSI: duit(t, "101"), Premium: duit(t, "102"), PremiNusantaraRe: duit(t, "103"),
		PremiRp: duit(t, "104"), PremiumGrossDiscountFleet: duit(t, "105"),
	}
}

func cedantLama(t *testing.T) []models.Cedant {
	return []models.Cedant{{CurrencyList: []models.BarisMataUang{barisMU(t, "401", "402", "1")}}}
}

// kasusLapisB - kasus endorsement yang SUDAH lahir: lapis A di OldData, daftar
// kerja hasil lapis C. Lapis B diuji langsung lewat IsiNilaiLama.
func kasusLapisB(lama models.OfferFacIn, kerja models.OfferFacIn) models.KasusEndorsement {
	kerja.OldData = &lama
	return models.KasusEndorsement{OfferFacIn: kerja}
}

// TestLapisBEnamLiniSesuaiPeta - E08: tiap lini mengisi properti nilai
// lamanya SENDIRI, dan cedant identik di semua lini.
func TestLapisBEnamLiniSesuaiPeta(t *testing.T) {
	cov := coverageLama(t)
	mu := barisMU(t, "201", "202", "2")

	t.Run("Golf", func(t *testing.T) {
		lama := models.OfferFacIn{CedingCedantList: cedantLama(t), LocationList: []models.Lokasi{{Property: models.PropertiLokasi{
			RiskLocation:           models.RisikoLokasi{AnekaList: []models.Aneka{{TSI: duit(t, "301"), CoverageList: []models.Coverage{cov}}}},
			TotalTSIPremiGrossList: []models.BarisMataUang{mu},
		}}}}
		k := kasusLapisB(lama, lama.Klon())
		jalankanB(t, &k, services.Predikat{IsGolfInsurance: true, IsEDM: true})
		an := k.OfferFacIn.LocationList[0].Property.RiskLocation.AnekaList[0]
		samaUang(t, "Aneka.TSIOld", an.TSIOld, duit(t, "301"))
		periksaCoverageTSIPremi(t, an.CoverageList[0])
		periksaBarisMU(t, k.OfferFacIn.LocationList[0].Property.TotalTSIPremiGrossList[0])
		periksaCedant(t, k)
	})
	t.Run("Aneka", func(t *testing.T) {
		lama := models.OfferFacIn{CedingCedantList: cedantLama(t), LocationList: []models.Lokasi{{Property: models.PropertiLokasi{
			RiskLocation: models.RisikoLokasi{OccupationList: []models.Okupasi{{AnekaList: []models.Aneka{{
				TSI: duit(t, "301"), CoverageList: []models.Coverage{cov},
			}}}}},
			TotalTSIPremiGrossList: []models.BarisMataUang{mu},
		}}}}
		k := kasusLapisB(lama, lama.Klon())
		jalankanB(t, &k, services.Predikat{IsAneka: true, IsEDM: true})
		an := k.OfferFacIn.LocationList[0].Property.RiskLocation.OccupationList[0].AnekaList[0]
		samaUang(t, "Aneka.TSIOld", an.TSIOld, duit(t, "301"))
		// Aneka TIDAK mengisi coverage di bawah AnekaList.
		if c := an.CoverageList[0]; !c.TSIOld.Kosong() || !c.PremiumOld.Kosong() {
			t.Errorf("Aneka mengisi coverage: %+v", c)
		}
		periksaBarisMU(t, k.OfferFacIn.LocationList[0].Property.TotalTSIPremiGrossList[0])
		periksaCedant(t, k)
	})
	t.Run("PA", func(t *testing.T) {
		lama := models.OfferFacIn{CedingCedantList: cedantLama(t), PersonList: []models.Orang{{
			ASMCoverage: []models.Coverage{cov}, TotalTSIPremiGrossList: []models.BarisMataUang{mu},
		}}}
		k := kasusLapisB(lama, lama.Klon())
		jalankanB(t, &k, services.Predikat{IsPA: true, IsEDM: true})
		periksaCoverageTSIPremi(t, k.OfferFacIn.PersonList[0].ASMCoverage[0])
		periksaBarisMU(t, k.OfferFacIn.PersonList[0].TotalTSIPremiGrossList[0])
		periksaCedant(t, k)
	})
	t.Run("MarineCargo", func(t *testing.T) {
		lama := models.OfferFacIn{CedingCedantList: cedantLama(t), CargoList: []models.Kargo{{
			CoverageList: []models.Coverage{cov}, TotalTSIPremiGrossList: []models.BarisMataUang{mu},
		}}}
		k := kasusLapisB(lama, lama.Klon())
		jalankanB(t, &k, services.Predikat{IsMarineCargo: true, IsEDM: true})
		c := k.OfferFacIn.CargoList[0].CoverageList[0]
		periksaCoverageTSIPremi(t, c)
		samaUang(t, "PremiNusantaraReOld", c.PremiNusantaraReOld, duit(t, "103"))
		if !c.PremiRpOld.Kosong() || !c.PremiumGrossDiscountFleetOld.Kosong() {
			t.Errorf("MC mengisi properti milik MBU: %+v", c)
		}
		periksaBarisMU(t, k.OfferFacIn.CargoList[0].TotalTSIPremiGrossList[0])
		periksaCedant(t, k)
	})
	t.Run("MBU", func(t *testing.T) {
		lama := models.OfferFacIn{CedingCedantList: cedantLama(t), VehicleList: []models.Kendaraan{{
			CoverageList: []models.Coverage{cov}, TotalTSIPremiGrossList: []models.BarisMataUang{mu},
		}}}
		k := kasusLapisB(lama, lama.Klon())
		jalankanB(t, &k, services.Predikat{IsMBU: true, IsEDM: true})
		c := k.OfferFacIn.VehicleList[0].CoverageList[0]
		periksaCoverageTSIPremi(t, c)
		samaUang(t, "PremiNusantaraReOld", c.PremiNusantaraReOld, duit(t, "103"))
		samaUang(t, "PremiRpOld", c.PremiRpOld, duit(t, "104"))
		samaUang(t, "PremiumGrossDiscountFleetOld", c.PremiumGrossDiscountFleetOld, duit(t, "105"))
		periksaBarisMU(t, k.OfferFacIn.VehicleList[0].TotalTSIPremiGrossList[0])
		periksaCedant(t, k)
	})
	t.Run("Travel", func(t *testing.T) {
		lama := models.OfferFacIn{CedingCedantList: cedantLama(t), PersonList: []models.Orang{{
			ASMCoverage: []models.Coverage{cov}, TotalTSIPremiGrossList: []models.BarisMataUang{mu},
		}}}
		k := kasusLapisB(lama, lama.Klon())
		jalankanB(t, &k, services.Predikat{IsTravel: true, IsEDM: true})
		c := k.OfferFacIn.PersonList[0].ASMCoverage[0]
		periksaCoverageTSIPremi(t, c)
		if !c.PremiNusantaraReOld.Kosong() {
			t.Errorf("Travel mengisi PremiNusantaraReOld")
		}
		// Travel TIDAK mengisi TotalTSIPremiGrossList (berbeda dari PA).
		if b := k.OfferFacIn.PersonList[0].TotalTSIPremiGrossList[0]; !b.TSIOld.Kosong() {
			t.Errorf("Travel mengisi TotalTSIPremiGrossList: %+v", b)
		}
		periksaCedant(t, k)
	})
}

func jalankanB(t *testing.T, k *models.KasusEndorsement, p services.Predikat) {
	t.Helper()
	if err := services.IsiNilaiLama(k, p); err != nil {
		t.Fatalf("IsiNilaiLama: %v", err)
	}
}

func periksaCoverageTSIPremi(t *testing.T, c models.Coverage) {
	t.Helper()
	samaUang(t, "Coverage.TSIOld", c.TSIOld, duit(t, "101"))
	samaUang(t, "Coverage.PremiumOld", c.PremiumOld, duit(t, "102"))
}

func periksaBarisMU(t *testing.T, b models.BarisMataUang) {
	t.Helper()
	samaUang(t, "TotalTSIPremiGross.TSIOld", b.TSIOld, duit(t, "201"))
	samaUang(t, "TotalTSIPremiGross.PremiumOld", b.PremiumOld, duit(t, "202"))
}

func periksaCedant(t *testing.T, k models.KasusEndorsement) {
	t.Helper()
	c := k.OfferFacIn.CedingCedantList[0].CurrencyList[0]
	samaUang(t, "cedant TSIOld", c.TSIOld, duit(t, "401"))
	samaUang(t, "cedant PremiumOld", c.PremiumOld, duit(t, "402"))
}

// TestPremiNusantaraReOldHanyaDuaLini - E08: hanya MC dan MBU. Lini lain
// dengan coverage berisi PremiNusantaraRe tidak menyentuhnya.
func TestPremiNusantaraReOldHanyaDuaLini(t *testing.T) {
	cov := coverageLama(t)
	lama := models.OfferFacIn{
		LocationList: []models.Lokasi{{Property: models.PropertiLokasi{RiskLocation: models.RisikoLokasi{
			AnekaList: []models.Aneka{{CoverageList: []models.Coverage{cov}}},
		}}}},
		PersonList: []models.Orang{{ASMCoverage: []models.Coverage{cov}}},
	}
	k := kasusLapisB(lama, lama.Klon())
	jalankanB(t, &k, services.Predikat{IsGolfInsurance: true, IsPA: true, IsTravel: true, IsEDM: true})
	if !k.OfferFacIn.LocationList[0].Property.RiskLocation.AnekaList[0].CoverageList[0].PremiNusantaraReOld.Kosong() ||
		!k.OfferFacIn.PersonList[0].ASMCoverage[0].PremiNusantaraReOld.Kosong() {
		t.Error("PremiNusantaraReOld terisi di luar MC dan MBU")
	}
}

// TestLapisBKosongMenjadiNol - E08: seluruh properti nilai lama memakai pola
// kosong → NOL (angka), bukan kosong. Mata uang ikut sumber.
func TestLapisBKosongMenjadiNol(t *testing.T) {
	covKosong := models.Coverage{TSI: kosong(), Premium: kosong(), PremiNusantaraRe: kosong(), PremiRp: kosong(), PremiumGrossDiscountFleet: kosong()}
	muKosong := models.BarisMataUang{TSI: kosong(), Premium: kosong(), Rate: uang.Ratio{Scale: 6}}
	lama := models.OfferFacIn{
		CedingCedantList: []models.Cedant{{CurrencyList: []models.BarisMataUang{muKosong}}},
		VehicleList:      []models.Kendaraan{{CoverageList: []models.Coverage{covKosong}, TotalTSIPremiGrossList: []models.BarisMataUang{muKosong}}},
	}
	k := kasusLapisB(lama, lama.Klon())
	jalankanB(t, &k, services.Predikat{IsMBU: true, IsEDM: true})
	nol := duit(t, "0")
	c := k.OfferFacIn.VehicleList[0].CoverageList[0]
	for label, v := range map[string]uang.Money{
		"TSIOld": c.TSIOld, "PremiumOld": c.PremiumOld, "PremiNusantaraReOld": c.PremiNusantaraReOld,
		"PremiRpOld": c.PremiRpOld, "PremiumGrossDiscountFleetOld": c.PremiumGrossDiscountFleetOld,
		"TotalTSIPremiGross.TSIOld":     k.OfferFacIn.VehicleList[0].TotalTSIPremiGrossList[0].TSIOld,
		"TotalTSIPremiGross.PremiumOld": k.OfferFacIn.VehicleList[0].TotalTSIPremiGrossList[0].PremiumOld,
		"cedant.TSIOld":                 k.OfferFacIn.CedingCedantList[0].CurrencyList[0].TSIOld,
	} {
		samaUang(t, label, v, nol)
	}
	samaRasio(t, "RateOld", k.OfferFacIn.VehicleList[0].TotalTSIPremiGrossList[0].RateOld, uang.Ratio{Value: duit(t, "0").Amount, Scale: 6})
}

// TestA5_PremiNusantaraReOldKosongMenjadiAngkaNol - E09 A.5: korpus memakai
// cadangan STRING "0"; di sini ANGKA nol bertipe Money. Ini PERBAIKAN, bukan
// port - satu-satunya di modul before-image.
func TestA5_PremiNusantaraReOldKosongMenjadiAngkaNol(t *testing.T) {
	cov := models.Coverage{PremiNusantaraRe: kosong()}
	for nama, p := range map[string]services.Predikat{
		"MC":  {IsMarineCargo: true, IsEDM: true},
		"MBU": {IsMBU: true, IsEDM: true},
	} {
		lama := models.OfferFacIn{
			CargoList:   []models.Kargo{{CoverageList: []models.Coverage{cov}}},
			VehicleList: []models.Kendaraan{{CoverageList: []models.Coverage{cov}}},
		}
		k := kasusLapisB(lama, lama.Klon())
		jalankanB(t, &k, p)
		var got uang.Money
		if p.IsMarineCargo {
			got = k.OfferFacIn.CargoList[0].CoverageList[0].PremiNusantaraReOld
		} else {
			got = k.OfferFacIn.VehicleList[0].CoverageList[0].PremiNusantaraReOld
		}
		if got.Kosong() || got.Amount.Sign() != 0 {
			t.Errorf("%s: PremiNusantaraReOld %q, mau angka nol", nama, got.String())
		}
	}
}

// TestK046_RateOld_GuardSelfReferential_6dari6 - E09: guard menguji RateOld
// di baris polis lama, lalu mengambil Rate. ⛔ Perilaku BENAR (K-046), bukan
// cacat. Polis lama tanpa RateOld → RateOld = 0 walau Rate terisi.
func TestK046_RateOld_GuardSelfReferential_6dari6(t *testing.T) {
	tanpaRateOld := barisMU(t, "1", "1", "7")
	denganRateOld := barisMU(t, "1", "1", "7")
	denganRateOld.RateOld = rasio(t, "3")

	lokasi := func(b models.BarisMataUang) models.OfferFacIn {
		return models.OfferFacIn{LocationList: []models.Lokasi{{Property: models.PropertiLokasi{TotalTSIPremiGrossList: []models.BarisMataUang{b}}}}}
	}
	orang := func(b models.BarisMataUang) models.OfferFacIn {
		return models.OfferFacIn{PersonList: []models.Orang{{TotalTSIPremiGrossList: []models.BarisMataUang{b}}}}
	}
	kargo := func(b models.BarisMataUang) models.OfferFacIn {
		return models.OfferFacIn{CargoList: []models.Kargo{{TotalTSIPremiGrossList: []models.BarisMataUang{b}}}}
	}
	kendaraan := func(b models.BarisMataUang) models.OfferFacIn {
		return models.OfferFacIn{VehicleList: []models.Kendaraan{{TotalTSIPremiGrossList: []models.BarisMataUang{b}}}}
	}
	lini := []struct {
		nama   string
		p      services.Predikat
		bangun func(models.BarisMataUang) models.OfferFacIn
		ambil  func(models.OfferFacIn) models.BarisMataUang
	}{
		{"FIRE 4.1.2", services.Predikat{IsFire: true, IsEDM: true}, lokasi, func(o models.OfferFacIn) models.BarisMataUang {
			return o.LocationList[0].Property.TotalTSIPremiGrossList[0]
		}},
		{"Golf 4.3.2", services.Predikat{IsGolfInsurance: true, IsEDM: true}, lokasi, func(o models.OfferFacIn) models.BarisMataUang {
			return o.LocationList[0].Property.TotalTSIPremiGrossList[0]
		}},
		{"Aneka 4.5.2", services.Predikat{IsAneka: true, IsEDM: true}, lokasi, func(o models.OfferFacIn) models.BarisMataUang {
			return o.LocationList[0].Property.TotalTSIPremiGrossList[0]
		}},
		{"PA 4.7.2", services.Predikat{IsPA: true, IsEDM: true}, orang, func(o models.OfferFacIn) models.BarisMataUang { return o.PersonList[0].TotalTSIPremiGrossList[0] }},
		{"MC 4.9.2", services.Predikat{IsMarineCargo: true, IsEDM: true}, kargo, func(o models.OfferFacIn) models.BarisMataUang { return o.CargoList[0].TotalTSIPremiGrossList[0] }},
		{"MBU 4.11.2", services.Predikat{IsMBU: true, IsEDM: true}, kendaraan, func(o models.OfferFacIn) models.BarisMataUang { return o.VehicleList[0].TotalTSIPremiGrossList[0] }},
	}
	for _, l := range lini {
		for _, u := range []struct {
			ket   string
			baris models.BarisMataUang
			mau   string
		}{{"tanpa RateOld lama", tanpaRateOld, "0"}, {"dengan RateOld lama", denganRateOld, "7"}} {
			lama := l.bangun(u.baris)
			k := kasusLapisB(lama, lama.Klon())
			jalankanB(t, &k, l.p)
			got := l.ambil(k.OfferFacIn).RateOld
			if got.Kosong() || got.Value.Cmp(rasio(t, u.mau).Value) != 0 {
				t.Errorf("%s, %s: RateOld %s, mau %s", l.nama, u.ket, got, u.mau)
			}
		}
	}
}

// TestK046_PremiumOld_SelfReferential_HanyaFire - E09: hanya FIRE 4.1.2 yang
// menguji PremiumOld; lini lain menguji Premium.
func TestK046_PremiumOld_SelfReferential_HanyaFire(t *testing.T) {
	b := barisMU(t, "1", "9", "1") // Premium terisi, PremiumOld lama kosong
	lama := models.OfferFacIn{LocationList: []models.Lokasi{{Property: models.PropertiLokasi{TotalTSIPremiGrossList: []models.BarisMataUang{b}}}}}

	fire := kasusLapisB(lama, lama.Klon())
	jalankanB(t, &fire, services.Predikat{IsFire: true, IsEDM: true})
	samaUang(t, "FIRE PremiumOld", fire.OfferFacIn.LocationList[0].Property.TotalTSIPremiGrossList[0].PremiumOld, duit(t, "0"))

	golf := kasusLapisB(lama, lama.Klon())
	jalankanB(t, &golf, services.Predikat{IsGolfInsurance: true, IsEDM: true})
	samaUang(t, "Golf PremiumOld", golf.OfferFacIn.LocationList[0].Property.TotalTSIPremiGrossList[0].PremiumOld, duit(t, "9"))
}

// TestLapisBSumbernyaLapisA - nilai lama diambil dari baris OldData, bukan dari
// baris kerja. Diisi ulang setiap layar dibuka: perubahan pengguna pada baris
// kerja tidak menggeser nilai lamanya.
func TestLapisBSumbernyaLapisA(t *testing.T) {
	lama := models.OfferFacIn{CargoList: []models.Kargo{{CoverageList: []models.Coverage{coverageLama(t)}}}}
	k := kasusLapisB(lama, lama.Klon())
	k.OfferFacIn.CargoList[0].CoverageList[0].TSI = duit(t, "999") // diubah pengguna
	p := services.Predikat{IsMarineCargo: true, IsEDM: true}
	jalankanB(t, &k, p)
	jalankanB(t, &k, p) // layar dibuka lagi
	c := k.OfferFacIn.CargoList[0].CoverageList[0]
	samaUang(t, "TSIOld", c.TSIOld, duit(t, "101"))
	samaUang(t, "TSI kerja", c.TSI, duit(t, "999"))
}

// TestLapisBTravelMembuatBarisTujuan - `[dugaan]`: Travel tidak punya varian
// lapis C, sehingga PersonList kerja bisa kosong; baris tujuan dibuat dan
// berisi nilai lama saja.
func TestLapisBTravelMembuatBarisTujuan(t *testing.T) {
	lama := models.OfferFacIn{PersonList: []models.Orang{{ASMCoverage: []models.Coverage{coverageLama(t)}}}}
	k := kasusLapisB(lama, models.OfferFacIn{})
	jalankanB(t, &k, services.Predikat{IsTravel: true, IsEDM: true})
	if len(k.OfferFacIn.PersonList) != 1 || len(k.OfferFacIn.PersonList[0].ASMCoverage) != 1 {
		t.Fatalf("baris tujuan tidak dibuat: %+v", k.OfferFacIn.PersonList)
	}
	c := k.OfferFacIn.PersonList[0].ASMCoverage[0]
	samaUang(t, "TSIOld", c.TSIOld, duit(t, "101"))
	if !c.TSI.Kosong() {
		t.Errorf("baris baru membawa TSI %q; mau nilai lama saja", c.TSI.String())
	}
}

// TestGerbangKeluarLapisB - E11: Life, bukan endorsement, dan lebih dari 100
// lokasi menghasilkan lapis B TIDAK TERSENTUH, tanpa galat.
func TestGerbangKeluarLapisB(t *testing.T) {
	lokasiN := func(n int) models.OfferFacIn {
		o := models.OfferFacIn{CedingCedantList: cedantLama(t)}
		for i := 0; i < n; i++ {
			o.LocationList = append(o.LocationList, models.Lokasi{Property: models.PropertiLokasi{
				TotalTSIList: []models.BarisMataUang{barisMU(t, "1", "1", "1")},
			}})
		}
		return o
	}
	for _, u := range []struct {
		nama   string
		lokasi int
		p      services.Predikat
		keluar bool
	}{
		{"jiwa", 1, services.Predikat{IsFire: true, IsLife: true, IsEDM: true}, true},
		{"bukan endorsement", 1, services.Predikat{IsFire: true}, true},
		{"101 lokasi", 101, services.Predikat{IsFire: true, IsEDM: true}, true},
		{"tepat 100 lokasi", 100, services.Predikat{IsFire: true, IsEDM: true}, false},
	} {
		lama := lokasiN(u.lokasi)
		k := kasusLapisB(lama, lama.Klon())
		jalankanB(t, &k, u.p)
		terisi := !k.OfferFacIn.CedingCedantList[0].CurrencyList[0].TSIOld.Kosong() ||
			!k.OfferFacIn.LocationList[0].Property.TotalTSIList[0].TSIOld.Kosong()
		if terisi == u.keluar {
			t.Errorf("%s: lapis B terisi=%v, mau keluar=%v", u.nama, terisi, u.keluar)
		}
	}
}
