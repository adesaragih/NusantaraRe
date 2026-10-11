package models

import "github.com/cockroachdb/apd/v3"

// Following Policy Number (tiket 70).

// BarisFollowing - satu baris grid GridViewFollowingNB (alias BrowseDataViewPolis1). Tanggal = teks DD-MM-YYYY.
type BarisFollowing struct {
	CaseID, PolicyNo, ReffNumber, InsuredName, SourceOfBusiness, QQName, StartPeriod, EndPeriod string
}

// NilaiFollowing - Copy All (CopyLocationFollowing_Act langkah 11): `.OfferFacIn.PercentShare` dan `.Policy.Payment.RICommision`
// case sumber; nil = NULL.
type NilaiFollowing struct {
	PercentShare, RICommision *apd.Decimal
}
