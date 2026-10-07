package models

import "github.com/cockroachdb/apd/v3"

// Tab Inw Fac Cedant Panels kasus FIRE (tiket 49).

// BarisCedant - satu `.OfferFacIn.CedingCedantList(n)` (T_CEDINGCEDANTLIST) atau satu Ceding Co tab General (sumber
// Add pertama, ShareCeding nil).
type BarisCedant struct {
	CedingCo, CedingCoName string
	ShareCeding            *apd.Decimal // % Share cedant
}

// KasusCedant - data tab Cedant yang tersimpan: Source of Business, Share Cedant Type ("" / "0" / "1"), FlagOnGoingPolicy
// (T_WORK_POLIS, gerbang proteksi), grid cedant, dan Ceding Co tab General (`QuotationData.CedingCoList`, atau CedingCo /
// CedingCoName bila daftar kosong).
type KasusCedant struct {
	SobName, ShareCedantType, FlagOnGoingPolicy string
	Cedant, CedingUmum                          []BarisCedant
}

// SimpanCedant - isi Save tab Cedant. ShareOfCeding nil = tidak diubah (Share Cedant Type kosong, SetShareOfCeding
// tidak menjalankan langkah mana pun).
type SimpanCedant struct {
	ShareCedantType string
	ShareOfCeding   *string
	Cedant          []BarisCedant
}
