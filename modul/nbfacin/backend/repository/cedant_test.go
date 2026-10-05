package repository

import "testing"

// TestSQLCedant - tiket 49: kepala per case Fac In (General / Quotation LEFT JOIN), grid urut SEQ_NO, Share of Ceding
// UPDATE-atau-INSERT Quotation, grid diganti utuh (sisip 7 kolom, % lewat TO_NUMBER).
func TestSQLCedant(t *testing.T) {
	tb := tabelCedant{work: "U.W", general: "U.G", quo: "U.Q", ceding: "U.C", cedant: "U.D"}
	for nama, c := range map[string]struct{ got, mau string }{
		"kepala": {sqlBacaKepalaCedant(tb), "SELECT TO_CHAR(g.SHARE_CEDANT_TYPE), q.SOB_NAME, q.CEDING_CO, q.CEDING_CO_NAME, w.FLAG_ONGOING_POLICY FROM U.W w " +
			"LEFT JOIN U.G g ON g.ID = w.ID LEFT JOIN U.Q q ON q.PARENT_ID = w.ID WHERE w.ID = :1 AND w.LINI = :2"},
		"grid": {sqlBacaCedant("U.D"), "SELECT CEDING_CO, CEDING_CO_NAME, " + angkaKeluar("SHARE_CEDING") +
			" FROM U.D WHERE PARENT_ID = :1 ORDER BY SEQ_NO"},
		"tipe":        {sqlUbahShareCedantType("U.G"), "UPDATE U.G SET SHARE_CEDANT_TYPE = :1 WHERE ID = :2"},
		"share ubah":  {sqlUbahShareOfCeding("U.Q"), "UPDATE U.Q SET SHARE_OF_CEDING = :1 WHERE PARENT_ID = :2"},
		"share sisip": {sqlSisipQuotationShareOfCeding("U.Q"), "INSERT INTO U.Q (ID, PARENT_ID, SHARE_OF_CEDING) VALUES (:1, :2, :3)"},
		"hapus":       {sqlHapusCedant("U.D"), "DELETE FROM U.D WHERE PARENT_ID = :1"},
		"sisip": {sqlSisipCedant("U.D"), "INSERT INTO U.D (ID, PARENT_ID, SEQ_NO, ROW_UID, CEDING_CO, CEDING_CO_NAME, SHARE_CEDING) " +
			"VALUES (:1, :2, :3, :4, :5, :6, " + angkaMasuk(":7") + ")"},
	} {
		if c.got != c.mau {
			t.Errorf("%s:\n got %s\nmau %s", nama, c.got, c.mau)
		}
	}
}

// TestPecahGabunganCeding - gabungan `;` T_QUOTATIONDATA dipecah per kode; nama dipasang hanya bila jumlah potongan sama.
func TestPecahGabunganCeding(t *testing.T) {
	if b := pecahGabunganCeding("A1;B2", "UJI A;UJI B"); len(b) != 2 || b[1].CedingCo != "B2" || b[1].CedingCoName != "UJI B" {
		t.Errorf("pasangan: %+v", b)
	}
	if b := pecahGabunganCeding("A1;B2", "UJI; A;UJI B"); len(b) != 2 || b[0].CedingCoName != "" {
		t.Errorf("nama ber-; : %+v", b)
	}
	if b := pecahGabunganCeding("A1", "UJI A"); len(b) != 1 || b[0].CedingCoName != "UJI A" {
		t.Errorf("tunggal: %+v", b)
	}
}
