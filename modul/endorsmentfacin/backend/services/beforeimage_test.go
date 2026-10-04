package services_test

// Seam 4 - before-image saat kasus endorsement lahir: lapis A, lapis C, dan
// porsi periode. Tiket E06 (tracer FIRE) dan E07 (54 salinan lapis A).
//
// Dibaca sesudah: bantu_test.go, services/beforeimage.go.

import (
	"errors"
	"testing"

	"nusantarare/inti/backend/uang"
	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/services"
)

func siapkanFire(t *testing.T, tanggalEdm string) services.HasilBeforeImage {
	t.Helper()
	lama := polisLamaFire(t)
	h, err := services.PrepareBeforeImage(services.MasukanBeforeImage{
		Kasus:     kasusBaru(t, tanggalEdm),
		PolisLama: &lama,
		Predikat:  services.Predikat{IsFire: true, IsEDM: true},
	})
	if err != nil {
		t.Fatalf("PrepareBeforeImage: %v", err)
	}
	return h
}

// TestTracerFireTigaLapis - E06: satu kasus kebakaran, ketiga lapis terisi dan
// tetap terpisah, porsi periode terhitung.
func TestTracerFireTigaLapis(t *testing.T) {
	h := siapkanFire(t, "2026-01-02 00:00")
	o := h.Kasus.OfferFacIn

	// Lapis A - dokumen polis lama tersimpan utuh di OldData.
	if o.OldData == nil || len(o.OldData.LocationList) != 1 {
		t.Fatalf("lapis A tidak termuat: %+v", o.OldData)
	}
	if o.OldData.QuotationData.BusinessCode != "UJI-BC" {
		t.Errorf("OldData.BusinessCode %q", o.OldData.QuotationData.BusinessCode)
	}

	// Lapis C - daftar lokasi DISALIN dari OldData lalu ditandai.
	if len(o.LocationList) != 1 {
		t.Fatalf("LocationList tidak tersalin: %d baris", len(o.LocationList))
	}
	lok := o.LocationList[0]
	item := lok.Property.PropertyItemList[0]
	cov := item.CoverageList[0]
	for label, nilai := range map[string]string{
		"lokasi": lok.IsOldData, "item": item.IsOldData, "coverage": cov.IsOldData,
	} {
		if nilai != models.NilaiOld {
			t.Errorf("penanda %s %q, mau %q", label, nilai, models.NilaiOld)
		}
	}
	if o.IsProRate != models.NilaiProrate {
		t.Errorf("IsProRate %q, mau %q", o.IsProRate, models.NilaiProrate)
	}
	if cov.EDM != "" {
		t.Errorf("coverage.EDM %q, mau kosong (FIRE mengosongkannya)", cov.EDM)
	}

	// Lapis B - nilai lama per baris dari baris OldData yang sama.
	samaUang(t, "TSIObjectItemOld", item.TSIObjectItemOld, duit(t, "1000"))
	samaUang(t, "TotalGrossPremiOld", item.TotalGrossPremiOld, duit(t, "12"))
	samaUang(t, "TotalPremiumNusantaraReOld", item.TotalPremiumNusantaraReOld, duit(t, "6"))
	samaUang(t, "TotalTSIList.TSIOld", lok.Property.TotalTSIList[0].TSIOld, duit(t, "1000"))
	samaUang(t, "cedant TSIOld", o.CedingCedantList[0].CurrencyList[0].TSIOld, duit(t, "400"))
	samaUang(t, "cedant PremiumOld", o.CedingCedantList[0].CurrencyList[0].PremiumOld, duit(t, "4"))

	// Porsi periode - periode 4 hari, endorsement di hari ke-1.
	if h.GalatPorsiPeriode != nil {
		t.Fatalf("GalatPorsiPeriode %v", h.GalatPorsiPeriode)
	}
	samaRasio(t, "ProrateStartEDM", o.ProrateStartEDM, uang.Ratio{Value: rasio(t, "0.25").Value, Scale: 20})
	samaRasio(t, "ProrateEDMEnd", o.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "0.75").Value, Scale: 20})
}

// TestLapisATidakIkutTertandai - salinan lapis C hidup terpisah: penandaan
// daftar kerja TIDAK boleh mengubah dokumen polis lama.
func TestLapisATidakIkutTertandai(t *testing.T) {
	h := siapkanFire(t, "2026-01-02 00:00")
	lama := h.Kasus.OfferFacIn.OldData
	if lama.LocationList[0].IsOldData != "" ||
		lama.LocationList[0].Property.PropertyItemList[0].IsOldData != "" ||
		lama.LocationList[0].Property.PropertyItemList[0].CoverageList[0].EDM != "UJI-EDM" {
		t.Errorf("lapis A ikut berubah oleh lapis C: %+v", lama.LocationList[0])
	}
	if !lama.LocationList[0].Property.PropertyItemList[0].TSIObjectItemOld.Kosong() {
		t.Errorf("lapis A ikut berubah oleh lapis B")
	}
}

// TestMasukanTidakDiubah - hasil adalah salinan. Jalur Travel dipakai karena
// lapis B menulis langsung ke PersonList kerja tanpa lapis C menyalinnya lebih
// dulu: salinan dangkal akan membocorkan tulisan itu ke masukan pemanggil.
func TestMasukanTidakDiubah(t *testing.T) {
	lama := models.OfferFacIn{
		CurrencyList: []models.BarisMataUang{barisMU(t, "1", "1", "1")},
		PolicyData:   polisLamaFire(t).PolicyData,
		PersonList:   []models.Orang{{ASMCoverage: []models.Coverage{{TSI: duit(t, "5")}}}},
	}
	k := kasusBaru(t, "2026-01-02 00:00")
	k.OfferFacIn.PersonList = []models.Orang{{ASMCoverage: []models.Coverage{{TSI: duit(t, "6")}}}}
	if _, err := services.PrepareBeforeImage(services.MasukanBeforeImage{
		Kasus: k, PolisLama: &lama, Predikat: services.Predikat{IsTravel: true, IsEDM: true},
	}); err != nil {
		t.Fatal(err)
	}
	if !k.OfferFacIn.PersonList[0].ASMCoverage[0].TSIOld.Kosong() {
		t.Errorf("masukan Kasus ikut terisi lapis B")
	}
	if k.OfferFacIn.OldData != nil || !lama.PersonList[0].ASMCoverage[0].TSIOld.Kosong() {
		t.Errorf("masukan berubah: OldData=%v", k.OfferFacIn.OldData)
	}
}

// TestDeltaTerhadapPolisLamaBukanKosong - E06: data kerja berangkat dari polis
// lama. Medan yang kosong di kasus baru terisi dari OldData.
func TestDeltaTerhadapPolisLamaBukanKosong(t *testing.T) {
	h := siapkanFire(t, "2026-01-02 00:00")
	o := h.Kasus.OfferFacIn
	if o.QuotationData.BusinessCode != "UJI-BC" || o.Currency != mataUangUji ||
		len(o.CurrencyList) != 1 || !o.PolicyData.StartDateTime.Equal(tanggal(t, "2026-01-01 00:00")) {
		t.Errorf("data kerja tidak berangkat dari polis lama: %+v", o.QuotationData)
	}
	// Medan yang diisi langkah 1-13 tetap: lapis A tidak menimpa EdmDate.
	if !o.QuotationData.EdmDate.Equal(tanggal(t, "2026-01-02 00:00")) {
		t.Errorf("EdmDate tertimpa: %v", o.QuotationData.EdmDate)
	}
}

// TestLapisALimaPuluhEmpatSalinan - E07: seluruh 54 penugasan langkah 14.3,
// termasuk enam cermin ke halaman lain dan tiga yang bukan dari OldData.
func TestLapisALimaPuluhEmpatSalinan(t *testing.T) {
	lama := models.OfferFacIn{
		QuotationData: models.Quotation{
			BusinessOldId: "a01", BusinessCode: "a02", BusinessFac: "a03", BusinessName: "a04",
			BusinessType: "a05", CedingCo: "a06", CedingCoName: "a07", SourceOfBusiness: "a08",
			SobName: "a09", MarketingCode: "a10", MarketingName: "a11", MOID: "a12",
			TeamGroup: "a13", InsuredName: "a14", InsuredID: "a15", NoOfferSlip: "a16",
			IsGroup: "a17", QQName: "a18", PolicyType: "a19", EDMDay: "a20",
		},
		PolicyData: models.PolicyData{
			StartDateTime: tanggal(t, "2026-01-01 00:00"), EndDateTime: tanggal(t, "2026-01-05 00:00"),
			OfferingDate: "p03", ProdDateTime: "p04",
			Payment: models.Pembayaran{Installment: "b01", RICommision: "b02", PctBrokerageFee: "b03"},
		},
		Currency: "s01", CurrencyList: []models.BarisMataUang{barisMU(t, "1", "1", "1")},
		Parameters: []byte(`{"uji":"s03"}`), InwardScale: "s04", PercentShare: "s05",
		AdditionalCapital: "s06", MaxPctTreatyCapacity: "s07", MaxTreatyCapacity: "s08",
		CurrentYear: "s09", ProRatePercent: "s10", ProRateType: "s11",
		CedingCedantList: []models.Cedant{{}}, ShareCedantType: "s13", IsSpecialAcceptance: "s14",
		BinderRNM: "s15", PPnCheck: "s16", PolicyMasterNumber: "s17", IsB2B: "s18",
	}
	h, err := services.PrepareBeforeImage(services.MasukanBeforeImage{
		Kasus: kasusBaru(t, "2026-01-02 00:00"), PolisLama: &lama,
	})
	if err != nil {
		t.Fatal(err)
	}
	o, q := h.Kasus.OfferFacIn, h.Kasus.OfferFacIn.QuotationData

	// 20 identitas bisnis.
	identitas := []string{q.BusinessOldId, q.BusinessCode, q.BusinessFac, q.BusinessName, q.BusinessType,
		q.CedingCo, q.CedingCoName, q.SourceOfBusiness, q.SobName, q.MarketingCode, q.MarketingName,
		q.MOID, q.TeamGroup, q.InsuredName, q.InsuredID, q.NoOfferSlip, q.IsGroup, q.QQName,
		q.PolicyType, q.EDMDay}
	for i, v := range identitas {
		if want := "a" + dua(i+1); v != want {
			t.Errorf("identitas ke-%d %q, mau %q", i+1, v, want)
		}
	}
	// 4 periode.
	if !o.PolicyData.StartDateTime.Equal(lama.PolicyData.StartDateTime) || !o.PolicyData.EndDateTime.Equal(lama.PolicyData.EndDateTime) ||
		o.PolicyData.OfferingDate != "p03" || o.PolicyData.ProdDateTime != "p04" {
		t.Errorf("periode: %+v", o.PolicyData)
	}
	// 18 struktur share & kapasitas.
	share := []string{o.Currency, "", "", o.InwardScale, o.PercentShare, o.AdditionalCapital,
		o.MaxPctTreatyCapacity, o.MaxTreatyCapacity, o.CurrentYear, o.ProRatePercent, o.ProRateType,
		"", o.ShareCedantType, o.IsSpecialAcceptance, o.BinderRNM, o.PPnCheck, o.PolicyMasterNumber, o.IsB2B}
	for i, v := range share {
		if v == "" {
			continue // CurrencyList, Parameters, CedingCedantList diperiksa di bawah
		}
		if want := "s" + dua(i+1); v != want {
			t.Errorf("share ke-%d %q, mau %q", i+1, v, want)
		}
	}
	if len(o.CurrencyList) != 1 || string(o.Parameters) != `{"uji":"s03"}` || len(o.CedingCedantList) != 1 {
		t.Errorf("halaman share tidak tersalin: CurrencyList=%d Parameters=%s Cedant=%d",
			len(o.CurrencyList), o.Parameters, len(o.CedingCedantList))
	}
	// 3 pembayaran - di agregat penawaran ...
	if o.PolicyData.Payment != lama.PolicyData.Payment {
		t.Errorf("OfferFacIn.PolicyData.Payment %+v", o.PolicyData.Payment)
	}

	// ... dan enam cermin ke halaman lain. ⛔ Kriteria kunci E07: keempat
	// halaman ini TIDAK BOLEH kosong.
	if h.Kasus.Quotation.BusinessType != "a05" {
		t.Errorf("newWorkPage.Quotation.BusinessType %q", h.Kasus.Quotation.BusinessType)
	}
	if h.Kasus.Policy.Payment != lama.PolicyData.Payment {
		t.Errorf("newWorkPage.Policy.Payment %+v - tiga field pembayaran ditulis DUA kali", h.Kasus.Policy.Payment)
	}
	if h.Kasus.IsB2B != "s18" {
		t.Errorf("newWorkPage.IsB2B %q", h.Kasus.IsB2B)
	}
	if h.InputDataCreditCARI6 != "a11" {
		t.Errorf("InputDataCredit.CARI6 %q, mau nama marketing polis lama", h.InputDataCreditCARI6)
	}

	// Tiga yang bukan dari OldData: penyelarasan antar-halaman.
	if h.PortalQuotationBusinessType != "a05" {
		t.Errorf("pyWorkPage.Quotation.BusinessType %q", h.PortalQuotationBusinessType)
	}
	if h.Kasus.Quotation.BusinessCode != "a02" || !h.Kasus.Quotation.EdmDate.Equal(q.EdmDate) {
		t.Errorf("newWorkPage.Quotation bukan salinan OfferFacIn.QuotationData: %+v", h.Kasus.Quotation)
	}
	if h.Kasus.PPnCheck != "s16" {
		t.Errorf("newWorkPage.PPnCheck %q", h.Kasus.PPnCheck)
	}
}

func dua(n int) string { return string(rune('0'+n/10)) + string(rune('0'+n%10)) }

// TestTeamGroupDariMasterMarketing - langkah 14.5: TeamGroup master marketing
// menimpa salinan lapis A, hanya bila terisi.
func TestTeamGroupDariMasterMarketing(t *testing.T) {
	lama := polisLamaFire(t)
	lama.QuotationData.TeamGroup = "LAMA"
	for _, u := range []struct{ master, mau string }{{"", "LAMA"}, {"BARU", "BARU"}} {
		h, err := services.PrepareBeforeImage(services.MasukanBeforeImage{
			Kasus: kasusBaru(t, "2026-01-02 00:00"), PolisLama: &lama, TeamGroupMarketing: u.master,
		})
		if err != nil {
			t.Fatal(err)
		}
		if h.Kasus.OfferFacIn.QuotationData.TeamGroup != u.mau || h.Kasus.Quotation.TeamGroup != u.mau {
			t.Errorf("master %q: TeamGroup %q / %q, mau %q", u.master,
				h.Kasus.OfferFacIn.QuotationData.TeamGroup, h.Kasus.Quotation.TeamGroup, u.mau)
		}
	}
}

// TestPolisLamaTakDitemukanMenyalinKosong - query lapis A tanpa baris:
// OldData halaman kosong, dan 54 salinan tetap berjalan (diport apa adanya).
func TestPolisLamaTakDitemukanMenyalinKosong(t *testing.T) {
	k := kasusBaru(t, "2026-01-02 00:00")
	k.OfferFacIn.QuotationData.BusinessCode = "TERISI-SEBELUMNYA"
	h, err := services.PrepareBeforeImage(services.MasukanBeforeImage{Kasus: k})
	if err != nil {
		t.Fatal(err)
	}
	if h.Kasus.OfferFacIn.OldData == nil {
		t.Fatal("OldData nil; mau halaman kosong")
	}
	if h.Kasus.OfferFacIn.QuotationData.BusinessCode != "" {
		t.Errorf("BusinessCode %q, mau tertimpa kosong dari OldData", h.Kasus.OfferFacIn.QuotationData.BusinessCode)
	}
	// Langkah 15 bergerbang CurrencyList OldData >= 1: tidak dihitung.
	if !h.Kasus.OfferFacIn.ProrateEDMEnd.Kosong() {
		t.Errorf("ProrateEDMEnd terisi %s padahal OldData.CurrencyList kosong", h.Kasus.OfferFacIn.ProrateEDMEnd)
	}
}

// porsiPeriode - jalankan seam dengan polis lama FIRE yang periodenya diubah.
func porsiPeriode(t *testing.T, mulai, akhir, tanggalEdm string) services.HasilBeforeImage {
	t.Helper()
	lama := polisLamaFire(t)
	lama.PolicyData.StartDateTime = tanggal(t, mulai)
	lama.PolicyData.EndDateTime = tanggal(t, akhir)
	h, err := services.PrepareBeforeImage(services.MasukanBeforeImage{
		Kasus: kasusBaru(t, tanggalEdm), PolisLama: &lama,
		Predikat: services.Predikat{IsFire: true, IsEDM: true},
	})
	if err != nil {
		t.Fatalf("PrepareBeforeImage: %v", err)
	}
	return h
}

// TestPorsiPeriodeSetengahKeAtas - `@Math.divide(…,20)` = HALF_UP (A37,
// dikonfirmasi work owner 01-10-2026): 1/3 → …333 (bawah setengah), 2/3 →
// …667 (atas setengah) pada 20 desimal.
func TestPorsiPeriodeSetengahKeAtas(t *testing.T) {
	h := porsiPeriode(t, "2026-01-01 00:00", "2026-01-04 00:00", "2026-01-02 00:00") // 1/3 dan 2/3
	if h.GalatPorsiPeriode != nil {
		t.Fatalf("GalatPorsiPeriode %v", h.GalatPorsiPeriode)
	}
	o := h.Kasus.OfferFacIn
	samaRasio(t, "ProrateStartEDM", o.ProrateStartEDM, uang.Ratio{Value: rasio(t, "0.33333333333333333333").Value, Scale: 20})
	samaRasio(t, "ProrateEDMEnd", o.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "0.66666666666666666667").Value, Scale: 20})
	if o.ProrateEDMEnd.String() != "0.66666666666666666667" {
		t.Errorf("teks %q, mau tepat 20 desimal", o.ProrateEDMEnd.String())
	}
}

// TestPorsiPeriodeNegatifSetengahMenjauhiNol - HALF_UP Java membulatkan
// setengah MENJAUHI nol; tanggal endorsement sesudah akhir periode memberi
// porsi sisa negatif. −1/3 → −0.33333333333333333333; −2/3 →
// −0.66666666666666666667.
//
// ⚠️ Ralat 01-10-2026: kasus 151/150 di bawah memakai tanggal BUATAN, bukan
// tanggal kasus nyata - yang sama dengan fixture NB-15 hanya BENTUK angkanya.
// Dan 151/150 pada 20 desimal memberi hasil yang sama untuk HALF_UP dan
// HALF_EVEN (hanya pemotongan yang tersingkir), jadi ia BUKAN bukti HALF_UP.
// Tanggal kasus nyata: TestPorsiPeriodeKasusNyataEDMFire (rekonsiliasi_test.go).
func TestPorsiPeriodeNegatifSetengahMenjauhiNol(t *testing.T) {
	h := porsiPeriode(t, "2026-01-01 00:00", "2026-01-04 00:00", "2026-01-05 00:00") // 4/3 dan −1/3
	if h.GalatPorsiPeriode != nil {
		t.Fatalf("GalatPorsiPeriode %v", h.GalatPorsiPeriode)
	}
	samaRasio(t, "ProrateStartEDM", h.Kasus.OfferFacIn.ProrateStartEDM, uang.Ratio{Value: rasio(t, "1.33333333333333333333").Value, Scale: 20})
	samaRasio(t, "ProrateEDMEnd", h.Kasus.OfferFacIn.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "-0.33333333333333333333").Value, Scale: 20})
	h = porsiPeriode(t, "2026-01-01 00:00", "2026-01-04 00:00", "2026-01-06 00:00") // 5/3 dan −2/3
	samaRasio(t, "5/3", h.Kasus.OfferFacIn.ProrateStartEDM, uang.Ratio{Value: rasio(t, "1.66666666666666666667").Value, Scale: 20})
	samaRasio(t, "−2/3", h.Kasus.OfferFacIn.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "-0.66666666666666666667").Value, Scale: 20})
	h = porsiPeriode(t, "2025-10-01 00:00", "2026-02-28 00:00", "2026-03-01 00:00") // 151/150, tanggal buatan
	samaRasio(t, "151/150", h.Kasus.OfferFacIn.ProrateStartEDM, uang.Ratio{Value: rasio(t, "1.00666666666666666667").Value, Scale: 20})
}

// TestPorsiPeriodeBukanHariBulatDitolak - lokal Pega bertipe int dan satuannya
// [dugaan] hari; pemotongan pecahan hari belum terverifikasi.
func TestPorsiPeriodeBukanHariBulatDitolak(t *testing.T) {
	h := porsiPeriode(t, "2026-01-01 00:00", "2026-01-05 00:00", "2026-01-02 12:00")
	if !errors.Is(h.GalatPorsiPeriode, services.ErrSatuanSelisihWaktuBelumTerverifikasi) {
		t.Fatalf("GalatPorsiPeriode %v, mau ErrSatuanSelisihWaktuBelumTerverifikasi", h.GalatPorsiPeriode)
	}
}

// TestPorsiPeriodeNol - guard `TotalPeriod = 0 → 1` (satu hari).
func TestPorsiPeriodeNol(t *testing.T) {
	sama := porsiPeriode(t, "2026-01-01 00:00", "2026-01-01 00:00", "2026-01-01 00:00")
	if sama.GalatPorsiPeriode != nil {
		t.Fatalf("GalatPorsiPeriode %v", sama.GalatPorsiPeriode)
	}
	samaRasio(t, "ProrateEDMEnd", sama.Kasus.OfferFacIn.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "0").Value, Scale: 20})

	// Endorsement sehari sesudah periode nol: 1 ÷ 1 dan −1 ÷ 1, diport apa adanya.
	sesudah := porsiPeriode(t, "2026-01-01 00:00", "2026-01-01 00:00", "2026-01-02 00:00")
	if sesudah.GalatPorsiPeriode != nil {
		t.Fatalf("GalatPorsiPeriode %v", sesudah.GalatPorsiPeriode)
	}
	samaRasio(t, "ProrateStartEDM", sesudah.Kasus.OfferFacIn.ProrateStartEDM, uang.Ratio{Value: rasio(t, "1").Value, Scale: 20})
	samaRasio(t, "ProrateEDMEnd", sesudah.Kasus.OfferFacIn.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "-1").Value, Scale: 20})
}
