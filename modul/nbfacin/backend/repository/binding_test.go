package repository

import (
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestSQLBinding - tiket 63: INSERT FACINOFFER 50 kolom urutan InsertOfferProduction_Sql (diperiksa terhadap korpus saat
// ditulis: urutan sama, A\B dan B\A kosong), 4 kolom NULL (Pega "" / "NULL"; NOSERIAL ber-bind sejak tiket 73 P4, enam kolom
// Marine TRADING_NAME / CONVEYANCE_NAME / IDSHIP / FROMRUTE / TORUTE / SAILDATE ber-bind sejak tiket 74 P5; "" = NULL),
// tanggal To_date 24-jam, RATE Replace titik -> koma, SPREAD_DATE dari JSON_OFFER ber-NOOFFER;
// bind berurutan dan jumlahnya = argFacInOffer; SQL JSON_OFFER / flag / email.
func TestSQLBinding(t *testing.T) {
	q := sqlSisipFacInOffer("UJI.FO", "UJI.JO")
	if !strings.HasPrefix(q, "INSERT INTO UJI.FO ("+strings.Join(kolomFacInOffer, ", ")+") VALUES (:1, :2, :3, :4, :5, TO_DATE(:6, 'DD/MM/YYYY HH24:MI:SS'), TO_DATE(:7, 'DD/MM/YYYY HH24:MI:SS'),") ||
		len(kolomFacInOffer) != 50 || strings.Count(q, ", NULL") != 4 || !strings.Contains(q, "REPLACE(:21, '.', ',')") ||
		!strings.Contains(q, "(SELECT MAX(j.TGL_INPUT) FROM UJI.JO j WHERE j.NOOFFER = :30)") || strings.Contains(q, "TSITOP") {
		t.Errorf("%s", q)
	}
	b := models.BarisFacInOffer{IDPega: "UJI-NB-1", Rate: "0.25", TSI: desimalUji("100"), Premi: desimalUji("2"), GroupType: "2", TopRisk: "true"}
	arg := argFacInOffer(b, "RNM-OFR.10.2026.00007")
	if n := bindBerurutan(t, "facinoffer", q); n != 46 || len(arg) != n {
		t.Errorf("bind %d, arg %d", n, len(arg))
	}
	// IDPEGA, NOOFFER, RATE (bind 21), TSI_MENJADI = TSI_SELISIH (26, 28), PREMI (27, 29), SPREAD_DATE = nomor offer (30),
	// GROUP_TYPE (39), NOSERIAL (40, "" = NULL), TOP_RISK terakhir; enam kolom Marine "" = NULL.
	if arg[0] != "UJI-NB-1" || arg[1] != "RNM-OFR.10.2026.00007" || arg[20] != "0.25" || arg[25] != "100" || arg[27] != "100" ||
		arg[26] != "2" || arg[28] != "2" || arg[29] != "RNM-OFR.10.2026.00007" || arg[38] != "2" || arg[39] != nil || arg[45] != "true" ||
		arg[2] != nil || arg[34] != nil || arg[40] != nil || arg[43] != nil ||
		!strings.Contains(q, ", :38, NULL, NULL, NULL, :39, :40, NULL, :41, :42, :43, TO_DATE(:44, 'DD/MM/YYYY HH24:MI:SS'), :45, :46)") {
		t.Errorf("%#v", arg)
	}
	b.NoSerial = "SN-UJI"
	b.TradingName, b.ConveyanceName, b.IDShip, b.FromRute, b.ToRute, b.SailDate = "TR", "CV", "S1", "A", "B", "15/03/2026 12:00:00"
	if arg := argFacInOffer(b, ""); arg[39] != "SN-UJI" || arg[34] != "TR" || arg[35] != "CV" || arg[40] != "S1" || arg[41] != "A" ||
		arg[42] != "B" || arg[43] != "15/03/2026 12:00:00" {
		t.Errorf("NOSERIAL / Marine: %#v", arg[34:44])
	}
	for nama, c := range map[string]struct {
		q        string
		mau      string
		jmlhBind int
	}{
		"flag":  {sqlFlagBinding("UJI.W"), "UPDATE UJI.W SET FLAG_ONGOING_POLICY = :1 WHERE ID = :2 AND LINI = :3", 3},
		"cari":  {sqlCariJSONOffer("UJI.JO"), "SELECT COUNT(*), MAX(NOOFFER) FROM UJI.JO WHERE IDPEGA IN (:1, :2)", 2},
		"sisip": {sqlSisipJSONOffer("UJI.JO"), "INSERT INTO UJI.JO (ID, IDPEGA, JSONDATA, PRODKE, NOOFFER) VALUES (:1, :2, :3, 0, :4)", 4},
		"ubah":  {sqlUbahJSONOffer("UJI.JO"), "UPDATE UJI.JO SET JSONDATA = :1 WHERE IDPEGA IN (:2, :3)", 3},
		"hapus": {sqlHapusFacInOffer("UJI.FO"), "DELETE FROM UJI.FO WHERE IDPEGA IN (:1, :2)", 2},
		"email": {sqlEmailAkun("UJI.L"), "SELECT MAX(EMAIL) FROM UJI.L WHERE UPPER(TRIM(LOGIN_ID)) = UPPER(TRIM(:1))", 1},
	} {
		if c.q != c.mau || bindBerurutan(t, nama, c.q) != c.jmlhBind {
			t.Errorf("%s: %q", nama, c.q)
		}
	}
	if q := sqlEmailMarketing("UJI.M", "UJI.L"); !strings.Contains(q, " m WHERE m.ID = :1 OR m.ID = (SELECT MAX(m2.CLIENTID2) FROM UJI.M m2 WHERE m2.ID = :2)") ||
		bindBerurutan(t, "email marketing", q) != 2 {
		t.Errorf("%q", q)
	}
	if !strings.Contains(sqlKonteksSubmit(tbSubmit), "(SELECT MAX(m.CLIENTID) FROM "+tbSubmit.mo+" m WHERE m.ID = q.MOID), g.CONFIRM_BINDING, g.RECEIVED_RI_SLIP, w.STATUS_WORK\nFROM ") {
		t.Errorf("konteks tanpa CLIENTID")
	}
}

// TestNoOfferBaru - 'RNM-OFR.' || MM.YYYY || '.' || LPAD(n, 5, '0'): kurang dari 5 digit diberi nol depan, lebih dipotong
// seperti LPAD Oracle.
func TestNoOfferBaru(t *testing.T) {
	for n, mau := range map[string]string{"7": "RNM-OFR.10.2026.00007", "12345": "RNM-OFR.10.2026.12345", "123456": "RNM-OFR.10.2026.12345"} {
		if got := noOfferBaru("10.2026", n); got != mau {
			t.Errorf("%s: %s", n, got)
		}
	}
}

// TestGalatSequence - ORA-02289 -> ErrSequenceJSONOffer (galat jelas, D2); galat lain apa adanya.
func TestGalatSequence(t *testing.T) {
	if err := galatSequence(errors.New("repository: mengambil nomor dari RNM_JSON_OFFER_SEQ: ORA-02289: sequence does not exist")); !errors.Is(err, ErrSequenceJSONOffer) {
		t.Errorf("%v", err)
	}
	lain := errors.New("ORA-00001")
	if err := galatSequence(lain); err != lain || galatSequence(nil) != nil {
		t.Errorf("%v", err)
	}
}

// TestSQLFormBinding - tiket 64: CONFIRM_BINDING / RECEIVED_RI_SLIP (migrasi 213), penjaga posisi + flag, ConfirmBinding hanya
// bila belum, TGL_BIND ber-IDPEGA; bind berurutan.
func TestSQLFormBinding(t *testing.T) {
	for nama, c := range map[string]struct {
		q, mau string
		bind   int
	}{
		"status":  {sqlStatusBinding("UJI.G"), "UPDATE UJI.G SET CONFIRM_BINDING = :1, RECEIVED_RI_SLIP = :2 WHERE ID = :3", 3},
		"sentuh":  {sqlSentuhTahap("UJI.W"), "UPDATE UJI.W SET TGL_UPDATE = SYSDATE WHERE ID = :1 AND LINI = :2 AND POSITION = :3 AND FLAG_ONGOING_POLICY = :4", 4},
		"konfirm": {sqlKonfirmasiBinding("UJI.G"), "UPDATE UJI.G SET CONFIRM_BINDING = '1' WHERE ID = :1 AND (CONFIRM_BINDING IS NULL OR CONFIRM_BINDING <> '1')", 1},
		"tglBind": {sqlTglBind("UJI.JO"), "UPDATE UJI.JO SET TGL_BIND = SYSDATE, JSONDATA = :1 WHERE IDPEGA IN (:2, :3)", 3},
	} {
		if c.q != c.mau || bindBerurutan(t, nama, c.q) != c.bind {
			t.Errorf("%s: %q", nama, c.q)
		}
	}
	if q := sqlSisipPerformance("F"); !strings.HasSuffix(q, "VALUES (:1, :2, :3, SYSDATE, :4, :5, :6, :7, :8, :9)") {
		t.Errorf("%q", q)
	}
}
