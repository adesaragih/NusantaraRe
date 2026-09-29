package models

// KombinasiTCO adalah kunci gabungan anak-anak kontrak treaty.
//
// `[fakta bisnis work owner]` reinsurer, business, dan klausul menggantung pada
// (TreatyYear, TreatyGroupID, ReinsTypeID) - BUKAN pada `ID` kontrak.
// `[terverifikasi]` `GetMasterReinsurerList.xml` (`TreatyYear={CARI4}`,
// `TreatyGroupID={CARI5}`, `ReinsTypeID={CARI6}`) dan
// `BrowseTreatyReinsurerList_Act` langkah 1: `CARI4 = Param.TreatyYear` - teks
// `TREATYYEAR` tahun treaty, bukan `ID`-nya.
//
// ⚠️ Akibat yang dibawa apa adanya: dua tahun treaty berteks tahun dan grup
// sama (periode berbeda - anti-dobel tiket 03 mengizinkannya) berbagi anak
// untuk jenis reasuransi yang sama.
type KombinasiTCO struct {
	TreatyYear      string
	TreatyGroupID   string
	TreatyGroupName string
	ReinsTypeID     string
	ReinsTypeName   string
}

// KombinasiDari menyusun kombinasi dari tahun induk dan kontraknya.
func KombinasiDari(t TahunTreaty, k KontrakTreaty) KombinasiTCO {
	return KombinasiTCO{TreatyYear: t.TreatyYear, TreatyGroupID: t.TreatyGroupID,
		TreatyGroupName: t.TreatyGroupName, ReinsTypeID: k.ReinsTypeID, ReinsTypeName: k.ReinsTypeName}
}
