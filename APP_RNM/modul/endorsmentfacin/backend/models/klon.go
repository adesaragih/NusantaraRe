package models

// Salinan dalam (deep copy) agregat penawaran.
//
// Untuk apa berkas ini: penugasan halaman Pega seperti
// `newWorkPage.OfferFacIn.LocationList = newWorkPage.OfferFacIn.OldData.LocationList`
// MENYALIN isinya - sesudahnya kedua daftar hidup terpisah. Di Go, menyalin
// struct berisi slice hanya menyalin penunjuknya, sehingga penandaan lapis C
// atas salinan akan ikut mengubah lapis A. Fungsi di sini menutup lubang itu.
//
// Dibaca sesudah: penawaran.go.
//
// ⚠️ `uang.Money` dan `uang.Ratio` membawa penunjuk `*apd.Decimal`. Penunjuk
// itu SENGAJA dibagi, tidak disalin: tidak satu pun kode di modul ini
// mengubah desimal di tempat - aritmetika `uang` selalu menghasilkan nilai
// baru.

import "encoding/json"

// Klon menyalin agregat beserta seluruh daftarnya dan lapis A-nya.
func (o OfferFacIn) Klon() OfferFacIn {
	h := o
	h.CurrencyList = KlonBarisMataUang(o.CurrencyList)
	h.Parameters = KlonJSON(o.Parameters)
	h.CedingCedantList = KlonCedant(o.CedingCedantList)
	h.LocationList = KlonLokasi(o.LocationList)
	h.CargoList = KlonKargo(o.CargoList)
	h.VehicleList = KlonKendaraan(o.VehicleList)
	h.PersonList = KlonOrang(o.PersonList)
	if o.OldData != nil {
		lama := o.OldData.Klon()
		h.OldData = &lama
	}
	return h
}

// Klon menyalin objek kerja endorsement beserta agregatnya.
func (k KasusEndorsement) Klon() KasusEndorsement {
	h := k
	h.OfferFacIn = k.OfferFacIn.Klon()
	h.TradingList = KlonJSON(k.TradingList)
	h.GoodsList = KlonJSON(k.GoodsList)
	h.ConveyanceList = KlonJSON(k.ConveyanceList)
	return h
}

// KlonCedant menyalin `CedingCedantList`.
func KlonCedant(s []Cedant) []Cedant {
	if s == nil {
		return nil
	}
	h := make([]Cedant, len(s))
	for i, c := range s {
		h[i] = Cedant{CurrencyList: KlonBarisMataUang(c.CurrencyList)}
	}
	return h
}

// KlonBarisMataUang menyalin daftar baris mata uang.
func KlonBarisMataUang(s []BarisMataUang) []BarisMataUang {
	if s == nil {
		return nil
	}
	return append([]BarisMataUang{}, s...)
}

// KlonLokasi menyalin `LocationList` sampai baris spreading.
func KlonLokasi(s []Lokasi) []Lokasi {
	if s == nil {
		return nil
	}
	h := make([]Lokasi, len(s))
	for i, l := range s {
		h[i] = l
		p := &h[i].Property
		p.TotalTSIList = KlonBarisMataUang(l.Property.TotalTSIList)
		p.TotalTSIPremiGrossList = KlonBarisMataUang(l.Property.TotalTSIPremiGrossList)
		p.RiskLocation.AnekaList = klonAneka(l.Property.RiskLocation.AnekaList)
		if l.Property.RiskLocation.OccupationList != nil {
			p.RiskLocation.OccupationList = make([]Okupasi, len(l.Property.RiskLocation.OccupationList))
			for j, o := range l.Property.RiskLocation.OccupationList {
				p.RiskLocation.OccupationList[j] = Okupasi{IsOldData: o.IsOldData, AnekaList: klonAneka(o.AnekaList)}
			}
		}
		if l.Property.PropertyItemList != nil {
			p.PropertyItemList = make([]ItemProperti, len(l.Property.PropertyItemList))
			for j, it := range l.Property.PropertyItemList {
				it.CoverageList = KlonCoverage(it.CoverageList)
				p.PropertyItemList[j] = it
			}
		}
	}
	return h
}

func klonAneka(s []Aneka) []Aneka {
	if s == nil {
		return nil
	}
	h := make([]Aneka, len(s))
	for i, a := range s {
		a.CoverageList = KlonCoverage(a.CoverageList)
		h[i] = a
	}
	return h
}

// KlonKargo menyalin `CargoList`.
func KlonKargo(s []Kargo) []Kargo {
	if s == nil {
		return nil
	}
	h := make([]Kargo, len(s))
	for i, k := range s {
		k.CoverageList = KlonCoverage(k.CoverageList)
		k.TotalTSIPremiGrossList = KlonBarisMataUang(k.TotalTSIPremiGrossList)
		k.TradingList = KlonJSON(k.TradingList)
		k.GoodList = KlonJSON(k.GoodList)
		k.ConveyanceList = KlonJSON(k.ConveyanceList)
		h[i] = k
	}
	return h
}

// KlonKendaraan menyalin `VehicleList`.
func KlonKendaraan(s []Kendaraan) []Kendaraan {
	if s == nil {
		return nil
	}
	h := make([]Kendaraan, len(s))
	for i, k := range s {
		k.CoverageList = KlonCoverage(k.CoverageList)
		k.TotalTSIPremiGrossList = KlonBarisMataUang(k.TotalTSIPremiGrossList)
		h[i] = k
	}
	return h
}

// KlonOrang menyalin `PersonList`.
func KlonOrang(s []Orang) []Orang {
	if s == nil {
		return nil
	}
	h := make([]Orang, len(s))
	for i, o := range s {
		o.ASMCoverage = KlonCoverage(o.ASMCoverage)
		o.CoverageList = KlonCoverage(o.CoverageList)
		o.TotalTSIPremiGrossList = KlonBarisMataUang(o.TotalTSIPremiGrossList)
		h[i] = o
	}
	return h
}

// KlonCoverage menyalin daftar coverage beserta spreading dan coverage
// tambahannya.
func KlonCoverage(s []Coverage) []Coverage {
	if s == nil {
		return nil
	}
	h := make([]Coverage, len(s))
	for i, c := range s {
		if c.SpreadingList != nil {
			c.SpreadingList = append([]Spreading{}, c.SpreadingList...)
		}
		c.AdditionalCoverage = KlonCoverage(c.AdditionalCoverage)
		h[i] = c
	}
	return h
}

// KlonJSON menyalin halaman berbentuk JSON mentah; nil tetap nil.
func KlonJSON(m json.RawMessage) json.RawMessage {
	if m == nil {
		return nil
	}
	return append(json.RawMessage(nil), m...)
}
