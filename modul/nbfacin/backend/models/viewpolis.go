package models

// View Policy tab Facultative (List Offer) / Treaty / Life (tiket 71 perluasan, permintaan work owner 11-10-2026). Teks apa adanya
// dari basis data; tanggal DD-MM-YYYY (Oracle TO_CHAR); NULL = "".

// BarisOfferFac - BrowseDataViewOfferF (POOLDATA.FACINOFFER, DISTINCT). QQ selalu kosong (SQL Pega memilih literal teks kosong AS "QQ").
type BarisOfferFac struct {
	CaseID, NoOffer, InsuredName, SourceOfBusiness, StartPeriod, EndPeriod string
}

// BarisPolisTreaty - BrowsePolisTreaty (POOLDATA.TREATYINPRODUCTION, DISTINCT).
type BarisPolisTreaty struct {
	CaseID, PolicyNo, InsuredName, SourceOfBusiness, StartPeriod, EndPeriod string
}

// BarisPolisLife - BrowsePolisLife (POOLDATA.LIFEINPRODUCTION).
type BarisPolisLife struct {
	CaseID, PolicyNo, EndorsNo, TypePolicy, SourceOfBusiness, CedingCoName, QQ, RISlipCeding string
}

// BarisOfferLife - BrowseOfferLife (POOLDATA.JSON_OFFER_LIFE versi terakhir per IDPEGA, NoOffer = 'OFR-' || ID versi pertama).
type BarisOfferLife struct {
	CaseID, NoOffer, TypePolicy, BusinessName, CedingCoName, QQ, LastStatus string
}

// BerkasPolisTreaty - tombol View tab Treaty (work owner 11-10-2026 "tombol view treaty nya mana"): berkas NB / EDM Treaty In sistem
// baru yang memuat nomor polis itu. Pola Claim Non Prop / Claim Prop (OQ-CNP-13) DISALIN, tanpa impor modul lain.
type BerkasPolisTreaty struct{ Modul, Kasus string }

// Nama modul frontend tujuan tombol View Treaty.
const (
	ModulNBTreatyIn  = "nbtreatyin"
	ModulEDMTreatyIn = "edmtreatyin"
)

// ModulBerkasPolisTreaty - T_GENERAL_POLIS_TREATY.PRODKE 0 = NB Treaty In, >= 1 = EDM Treaty In.
func ModulBerkasPolisTreaty(prodke int64) string {
	if prodke >= 1 {
		return ModulEDMTreatyIn
	}
	return ModulNBTreatyIn
}
