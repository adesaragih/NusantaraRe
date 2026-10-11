package repository

import "testing"

// TestSQLCsvObjek - tiket 66: pencarian persis Pega (BrowseRW2_SQL zip + NOTE tanpa spasi, GetDataCurrencyByName_SQL,
// GetDataObjectItem TYPE FIRE + ISACTIVE, BrowseOccupationFacInFIRE_RD TYPE FIRE + OldID urut Name), Obj-Browse persis, dan
// penghapus CurrencyList beserta cicilannya; bind berurutan.
func TestSQLCsvObjek(t *testing.T) {
	for nama, c := range map[string]struct {
		q, mau string
		bind   int
	}{
		"rw": {sqlRWKelurahan("UJI.RW"), "SELECT NOTE, DISTRICTNAME, CITYNAME, PROVINCENAME, NATION FROM UJI.RW" +
			" WHERE ZIPCODE = :1 AND REPLACE(NOTE, ' ', '') = REPLACE(:2, ' ', '')", 2},
		"risk":      {sqlCariIDRiskAddress("UJI.RA"), "SELECT ID FROM UJI.RA WHERE ADDRESS = :1 AND POSTALCODE = :2 ORDER BY ID FETCH FIRST 1 ROWS ONLY", 2},
		"akumulasi": {sqlIDAkumulasi("UJI.AC"), "SELECT ID FROM UJI.AC WHERE NOTE = :1 ORDER BY ID FETCH FIRST 1 ROWS ONLY", 1},
		"negara":    {sqlInisialNegara("UJI.N"), "SELECT NATIONINITIAL FROM UJI.N WHERE NOTE = :1 ORDER BY ID FETCH FIRST 1 ROWS ONLY", 1},
		"provinsi":  {sqlIDProvinsi("UJI.P"), "SELECT ID FROM UJI.P WHERE NOTE = :1 ORDER BY ID FETCH FIRST 1 ROWS ONLY", 1},
		"mata uang": {sqlMataUangCsv("UJI.C"), "SELECT ID, OLDID, CURRENCY FROM UJI.C WHERE CURRENCY = :1 ORDER BY ID FETCH FIRST 1 ROWS ONLY", 1},
		"item": {sqlJenisItemFire("UJI.O"), "SELECT ID FROM UJI.O WHERE TYPE = 'FIRE' AND REPLACE(OBJECTITEMTYPE, ' ', '') = REPLACE(:1, ' ', '')" +
			" AND ISACTIVE = '1' ORDER BY ID FETCH FIRST 1 ROWS ONLY", 1},
		"okupasi":  {sqlNamaOkupasiFire("UJI.OC"), "SELECT NAME FROM UJI.OC WHERE TYPE = 'FIRE' AND OLDID = :1 ORDER BY NAME, OLDID FETCH FIRST 1 ROWS ONLY", 1},
		"coverage": {sqlCoverageAda("UJI.CV", 3), "SELECT ID FROM UJI.CV WHERE ID IN (:1, :2, :3)", 3},
		"cicilan":  {sqlHapusCicilanCase("UJI.I", "UJI.CL"), "DELETE FROM UJI.I WHERE PARENT_ID IN (SELECT ID FROM UJI.CL WHERE PARENT_ID = :1)", 1},
		"currency": {sqlHapusCurrencyCase("UJI.CL"), "DELETE FROM UJI.CL WHERE PARENT_ID = :1", 1},
	} {
		if c.q != c.mau || bindBerurutan(t, nama, c.q) != c.bind {
			t.Errorf("%s: %q", nama, c.q)
		}
	}
	if TabelObjectItemType != "OBJECTITEMTYPE" {
		t.Error(TabelObjectItemType)
	}
}
