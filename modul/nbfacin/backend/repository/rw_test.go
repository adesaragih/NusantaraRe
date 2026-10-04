package repository

import (
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestSQLRW - tiket 37: saran Zip Code - STS_AKTIF terikat (:1), ZIPCODE DIAWALI pola (:2,
// ESCAPE), DISTINCT, urut ZIPCODE, batas terikat (:3); NATION dibaca sebagai negara.
func TestSQLRW(t *testing.T) {
	mau := "SELECT ZIPCODE, NOTE, DISTRICTNAME, CITYNAME, PROVINCENAME, NATION FROM (SELECT DISTINCT ZIPCODE, NOTE, DISTRICTNAME," +
		" CITYNAME, PROVINCENAME, NATION FROM UJI.RW WHERE STS_AKTIF = :1 AND ZIPCODE LIKE :2 ESCAPE '\\')" +
		" ORDER BY ZIPCODE, NOTE, DISTRICTNAME, CITYNAME, PROVINCENAME, NATION FETCH FIRST :3 ROWS ONLY"
	if got := sqlCariRW("UJI.RW"); got != mau {
		t.Errorf("SQL\n%q\nmau\n%q", got, mau)
	}
	for masuk, mau := range map[string]string{"123": "123%", " 12_4 ": `12\_4%`, "1%": `1\%%`, "": "", "  ": "", "ab1": "ab1%"} {
		if got := polaAwalan(masuk); got != mau {
			t.Errorf("polaAwalan(%q) = %q, mau %q", masuk, got, mau)
		}
	}
	if StatusRWAktif != "1" || BatasSaranRW != 50 {
		t.Error("konstanta RD/A124 berubah")
	}
}

// TestSQLSimpanRiskAddress - butir 81: ID = ekspresi prosedur InsertUpdateRISKADDRESS PERSIS;
// INSERT = urutan parameter prosedur (P_NationName .. P_Postalcode), tanpa COMMIT.
func TestSQLSimpanRiskAddress(t *testing.T) {
	if got := sqlIDRiskAddress("UJI.GETCURRENTSITE", "UJI.RISKADDRESS_SEQ"); got != "SELECT UJI.GETCURRENTSITE || LPAD(TO_CHAR(UJI.RISKADDRESS_SEQ.NEXTVAL), 12, '0') FROM DUAL" {
		t.Errorf("ID %q", got)
	}
	q := sqlSisipRiskAddress("UJI.RISKADDRESS")
	if q != "INSERT INTO UJI.RISKADDRESS (ID, NATIONNAME, PROVINCENAME, DISTRICTNAME, CITYNAME, TERRITORYNAME, TITLE, ADDRESS, POSTALCODE) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9)" ||
		strings.Contains(strings.ToUpper(q), "COMMIT") {
		t.Errorf("INSERT %q", q)
	}
	arg := argSisipRiskAddress("UJI-ID", models.AlamatBaru{NationName: "N", ProvinceName: "P", DistrictName: "D", CityName: "C",
		TerritoryName: "T", Title: "JL.", Address: "A", PostalCode: "Z"})
	mau := []any{"UJI-ID", "N", "P", "D", "C", "T", "JL.", "A", "Z"}
	for i := range mau {
		if arg[i] != mau[i] {
			t.Errorf("bind :%d = %v, mau %v", i+1, arg[i], mau[i])
		}
	}
	if kosong := argSisipRiskAddress("X", models.AlamatBaru{}); kosong[1] != nil || kosong[8] != nil {
		t.Errorf("teks kosong harus NULL: %v", kosong)
	}
}

// TestPeriksaIDRiskAddress - ID baru <= 15 karakter (vID varchar2(15) prosedur), tidak kosong.
func TestPeriksaIDRiskAddress(t *testing.T) {
	for id, sah := range map[string]bool{"UJI000000000001": true, "U": true, "": false, "UJIX000000000001": false} {
		if err := periksaIDRiskAddress(id); (err == nil) != sah {
			t.Errorf("%q: %v, mau sah=%v", id, err, sah)
		}
	}
}
