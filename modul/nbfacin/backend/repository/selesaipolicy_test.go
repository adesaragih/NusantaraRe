package repository

import (
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestSQLSelesaiPolicy - tiket 69a-2: FACINPRODUCTION 79 kolom urut InsertTreatyProduction_Sql (sensus dua cara sesi c3: urai
// daftar kolom korpus = 79, token = 79, VALUES = 79; A\B = B\A = ∅, urutan sama), bind berurutan = jumlah argumen, teks tetap
// (K69a2-4 "NULL", PRORATE 1, NULL halaman tak dikenal), GetMOID_Sql; JSON_POLIS / CEDING / M_RISK_LOSS_PROFILE.
func TestSQLSelesaiPolicy(t *testing.T) {
	unik := map[string]bool{}
	for _, k := range kolomFacInProduction {
		unik[k] = true
	}
	if len(kolomFacInProduction) != 79 || len(unik) != 79 || kolomFacInProduction[0] != "IDPEGA" ||
		kolomFacInProduction[78] != "DEDUCTION2_SELISIH" {
		t.Fatalf("kolom %d / unik %d", len(kolomFacInProduction), len(unik))
	}
	q := sqlSisipFacInProduction("UJI.P", "UJI.M")
	n := bindBerurutan(t, "FACINPRODUCTION", q)
	arg := argFacInProduction(models.BarisFacInProduction{IDPega: "ID", LayerNo: "NULL", Rate: apd.New(15, -2)}, "NP", "01/11/2026 12:00:00")
	// 79 kolom - 10 NULL - 2 'NULL' - PRORATE 1 = 66 kolom ber-bind; MARKETINGOFFICERCODE dua bind (NVL MOID case, tiket 73 P5) = 67.
	// Tiket 74 P5: enam kolom Marine (TRADING_NAME / CONVEYANCE_NAME / IDSHIP / FROMRUTE / TORUTE / SAILDATE) ber-bind.
	if n != len(arg) || n != 67 {
		t.Errorf("bind %d argumen %d", n, len(arg))
	}
	for _, x := range []string{"INSERT INTO UJI.P (" + strings.Join(kolomFacInProduction, ", ") + ") VALUES (:1, :2, NULL, :3,",
		"TO_DATE(:6, 'DD/MM/YYYY HH24:MI:SS'), TO_DATE(:7, 'DD/MM/YYYY HH24:MI:SS')", "'NULL', NULL, NULL, :", ", :41, 'NULL', :42,",
		"NVL(:33, (SELECT MIN(m.ID) FROM UJI.M m WHERE m.CLIENTID = :34", " AND m.MOSTATUS = '1'))", ", 1, TO_NUMBER(:"} {
		if !strings.Contains(q, x) {
			t.Errorf("tanpa %q:\n%s", x, q)
		}
	}
	if strings.Count(q, "'NULL'") != 2 || strings.Count(q, "TO_NUMBER(") != 19 || strings.Count(q, "TO_DATE(") != 4 {
		t.Errorf("teks tetap 'NULL' %d / TO_NUMBER %d / TO_DATE %d", strings.Count(q, "'NULL'"), strings.Count(q, "TO_NUMBER("),
			strings.Count(q, "TO_DATE("))
	}
	if arg[0] != "ID" || arg[1] != "NP" {
		t.Errorf("argumen awal %v", arg[:2])
	}
	// FIRE: MOID case NULL (NVL -> GetMOID), NOSERIAL teks "NULL" (K69a2-4) - hasil sama dengan sebelum tiket 73 P5.
	if arg[32] != nil || arg[40] != "NULL" || arg[35] != nil {
		t.Errorf("MOID / NOSERIAL / TRADING_NAME FIRE: %v / %v / %v", arg[32], arg[40], arg[35])
	}
	arg = argFacInProduction(models.BarisFacInProduction{MOIDKasus: "MO-UJI", NoSerial: "SN-UJI", TradingName: "TR", IDShip: "S1",
		SailDate: "15/03/2026 12:00:00"}, "", "")
	if arg[32] != "MO-UJI" || arg[40] != "SN-UJI" || arg[35] != "TR" {
		t.Errorf("MOID / NOSERIAL / TRADING_NAME: %v / %v / %v", arg[32], arg[40], arg[35])
	}
	if !strings.Contains(q, "TO_DATE(:") || !strings.Contains(sqlSisipFacInProduction("U.P", "U.M"), "SAILDATE, POLICYTYPE") {
		t.Error("SAILDATE")
	}
	for nama, c := range map[string]struct {
		q    string
		mau  string
		bind int
	}{
		"cari polis":   {sqlCariJSONPolis("UJI.J"), "SELECT NOPOLIS, TO_CHAR(TGL_PROD, 'DD/MM/YYYY HH24:MI:SS') FROM UJI.J WHERE IDPEGA IN (:1, :2) ORDER BY TGL_INPUT DESC FETCH FIRST 1 ROWS ONLY", 2},
		"sisip polis":  {sqlSisipJSONPolis("UJI.J"), "INSERT INTO UJI.J (IDPEGA, DATA_JSON, TGL_INPUT, NOPOLIS, PRODKE, OLDNOPOLIS, TGL_PROD, USERNAME, NOENDORS) VALUES (:1, :2, SYSDATE, :3, '0', NULL, TO_DATE(:4, 'DD/MM/YYYY HH24:MI:SS'), :5, NULL)", 5},
		"hapus ceding": {sqlHapusCedingProduksi("UJI.C"), "DELETE FROM UJI.C WHERE IDPEGA IN (:1, :2)", 2},
		"ada risiko":   {sqlAdaProfilRisiko("UJI.R"), "SELECT COUNT(*) FROM UJI.R WHERE NOPOLIS = :1", 1},
		"cek facin":    {sqlAdaFacInProduction("UJI.P"), "SELECT COUNT(*) FROM UJI.P WHERE IDPEGA IN (:1, :2)", 2},
	} {
		if c.q != c.mau || bindBerurutan(t, nama, c.q) != c.bind {
			t.Errorf("%s: %q", nama, c.q)
		}
	}
	ced := sqlSisipCedingProduksi("UJI.C", "UJI.J")
	// NOENDORS_ARAS tidak ada di skema DEV (ORA-00904, 09-10-2026) - tidak ditulis.
	if bindBerurutan(t, "ceding", ced) != 8 || !strings.HasPrefix(ced, "INSERT INTO UJI.C (IDPEGA, NOPOLIS, NOENDORS, CEDINGID, CEDINGNAME, PCT_SHARE, PRODKE, NOPOLIS_ARAS) VALUES (:1, :2, NULL, :3, :4, TO_NUMBER(:5,") ||
		!strings.HasSuffix(ced, "(SELECT MAX(j.PRODKE) FROM UJI.J j WHERE j.IDPEGA IN (:6, :7)), REPLACE(:8, '.', ''))") || strings.Contains(ced, "UPDATE") ||
		strings.Contains(ced, "NOENDORS_ARAS") {
		t.Errorf("ceding %q", ced)
	}
	ris := sqlSisipProfilRisiko("UJI.R")
	if bindBerurutan(t, "risiko", ris) != 7 || !strings.HasPrefix(ris, "INSERT INTO UJI.R (NOPOLIS, TANGGALSPREADING, TSINUSANTARARE, PREMINUSANTARARE, PREMIUMOR, PREMIUMQS, PREMIUMFACOUT, COB) VALUES (:1, SYSDATE, TO_NUMBER(:2,") {
		t.Errorf("risiko %q", ris)
	}
	if !strings.Contains(sqlPosisiSubmit("W"), "STATUS_WORK NOT LIKE 'Resolved-%'") {
		t.Error("penjaga Resolved")
	}
}
