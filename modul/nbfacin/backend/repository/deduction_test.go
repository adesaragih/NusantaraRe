package repository

import "testing"

// TestSQLDeduction - tiket 50: kepala per case Fac In (General LEFT JOIN), CurrencyList urut SEQ_NO, persen lewat
// TO_NUMBER, CurrencyList diperbarui menurut nama (ubah / sisip 12 kolom / buang beserta cicilan), tarif TAX = Obj-Browse Sort D TANGGAL, satu baris.
func TestSQLDeduction(t *testing.T) {
	tb := tabelDeduction{work: "U.W", general: "U.G", cur: "U.C"}
	for nama, c := range map[string]struct{ got, mau string }{
		"kepala": {sqlBacaKepalaDeduction(tb), "SELECT " + angkaKeluar("g.PAY_RI_COMMISION") + ", " + angkaKeluar("g.PAY_PCT_BROKERAGE_FEE") +
			", g.PPN_CHECK FROM U.W w LEFT JOIN U.G g ON g.ID = w.ID WHERE w.ID = :1 AND w.LINI = :2"},
		"baris": {sqlBacaCurrencyList("U.C"), "SELECT NAME, " + angkaKeluar("PAY_PREMIUM") + ", " + angkaKeluar("PAY_COMMISION") + ", " +
			angkaKeluar("PAY_DEDUCTION2") + ", " + angkaKeluar("PAY_BROKERAGE_FEE") + ", " + angkaKeluar("PAY_PPH") + ", " + angkaKeluar("PAY_PPN") +
			" FROM U.C WHERE PARENT_ID = :1 ORDER BY SEQ_NO"},
		"ubah": {sqlUbahKepalaDeduction("U.G"), "UPDATE U.G SET PAY_RI_COMMISION = " + angkaMasuk(":1") + ", PAY_PCT_BROKERAGE_FEE = " +
			angkaMasuk(":2") + ", PPN_CHECK = :3 WHERE ID = :4"},
		"id": {sqlIDCurrencyList("U.C"), "SELECT TO_CHAR(ID), NAME FROM U.C WHERE PARENT_ID = :1"},
		"ubah baris": {sqlUbahCurrencyList("U.C"), "UPDATE U.C SET SEQ_NO = :1, PAY_PREMIUM = " + angkaMasuk(":2") + ", PAY_COMMISION = " +
			angkaMasuk(":3") + ", PAY_DEDUCTION2 = " + angkaMasuk(":4") + ", PAY_BROKERAGE_FEE = " + angkaMasuk(":5") + ", PAY_PPH = " +
			angkaMasuk(":6") + ", PAY_PPN = " + angkaMasuk(":7") + " WHERE ID = :8"},
		"buang cicilan": {sqlHapusCicilan("U.I"), "DELETE FROM U.I WHERE PARENT_ID = :1"},
		"buang baris":   {sqlHapusSatuCurrency("U.C"), "DELETE FROM U.C WHERE ID = :1"},
		"sisip": {sqlSisipCurrencyList("U.C"), "INSERT INTO U.C (ID, PARENT_ID, SEQ_NO, ROW_UID, CURRENCY_CODE, NAME, PAY_PREMIUM, PAY_COMMISION," +
			" PAY_DEDUCTION2, PAY_BROKERAGE_FEE, PAY_PPH, PAY_PPN) VALUES (:1, :2, :3, :4, :5, :6, " + angkaMasuk(":7") + ", " + angkaMasuk(":8") +
			", " + angkaMasuk(":9") + ", " + angkaMasuk(":10") + ", " + angkaMasuk(":11") + ", " + angkaMasuk(":12") + ")"},
		"tarif": {sqlTarifPajak("U.T"), "SELECT " + angkaKeluar("PPH") + ", " + angkaKeluar("PPN") + " FROM U.T ORDER BY TANGGAL DESC FETCH FIRST 1 ROWS ONLY"},
	} {
		if c.got != c.mau {
			t.Errorf("%s:\n got %s\nmau %s", nama, c.got, c.mau)
		}
	}
}
