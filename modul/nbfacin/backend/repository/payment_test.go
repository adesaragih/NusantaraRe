package repository

import "testing"

// TestSQLPayment - tiket 51: kepala per case Fac In (PAY_INSTALLMENT), cicilan urut mata uang lalu cicilan, Total Payment /
// Net Premium lewat TO_NUMBER, sisip cicilan 16 kolom (persen teks apa adanya, uang lewat TO_NUMBER).
func TestSQLPayment(t *testing.T) {
	tb := tabelPayment{work: "U.W", general: "U.G", cur: "U.C", inst: "U.I"}
	for nama, c := range map[string]struct{ got, mau string }{
		"kepala": {sqlBacaKepalaPayment(tb), "SELECT TO_CHAR(g.PAY_INSTALLMENT) FROM U.W w LEFT JOIN U.G g ON g.ID = w.ID WHERE w.ID = :1 AND w.LINI = :2"},
		"cicilan": {sqlBacaCicilan(tb), "SELECT c.NAME, TO_CHAR(i.INSTALLMENT_NO), i.DUE_DATE, i.INSTALLMENT_PERCENTAGE FROM U.I i JOIN U.C c" +
			" ON c.ID = i.PARENT_ID WHERE c.PARENT_ID = :1 ORDER BY c.SEQ_NO, i.SEQ_NO"},
		"urut":        {sqlUrutCurrencyList("U.C"), "UPDATE U.C SET SEQ_NO = :1 WHERE ID = :2"},
		"installment": {sqlUbahInstallment("U.G"), "UPDATE U.G SET PAY_INSTALLMENT = :1 WHERE ID = :2"},
		"total": {sqlUbahTotalPayment("U.C"), "UPDATE U.C SET SUM_TOTAL_PAYMENT = " + angkaMasuk(":1") + ", PAY_NET_PREMIUM = " + angkaMasuk(":2") +
			" WHERE ID = :3"},
		"sisip": {sqlSisipCicilan("U.I"), "INSERT INTO U.I (ID, PARENT_ID, SEQ_NO, ROW_UID, CURRENCY_CODE, INSTALLMENT_NO, DUE_DATE," +
			" INSTALLMENT_PERCENTAGE, PREMIUM, RI_COMMISION, DEDUCTION2, BROKERAGE_FEE, PPH, PPN, DISCOUNT, PAYMENT_TOTAL) VALUES (:1, :2, :3, :4," +
			" :5, :6, :7, :8, " + angkaMasuk(":9") + ", " + angkaMasuk(":10") + ", " + angkaMasuk(":11") + ", " + angkaMasuk(":12") + ", " +
			angkaMasuk(":13") + ", " + angkaMasuk(":14") + ", " + angkaMasuk(":15") + ", " + angkaMasuk(":16") + ")"},
	} {
		if c.got != c.mau {
			t.Errorf("%s:\n got %s\nmau %s", nama, c.got, c.mau)
		}
	}
}
