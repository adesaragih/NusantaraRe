package repository

import (
	"strings"
	"testing"
)

// TestSQLFormAddAkumulasi - CZone zip = RW.CZONE zip itu + ID CZONE (A181, menyimpang dari GetAccumulationZipcode_SQL);
// Zip Code = BrowseRiskAddressZipCode_RD (ProvinceName Contains; berhalaman, A185) + q atas zip, NationInitial
// SetCountryID_Act (ID = IDNATION OR NOTE = NATIONNAME); saringan kosong dibuang, bind sejalan.
func TestSQLFormAddAkumulasi(t *testing.T) {
	// A181: Cresta Zone dari RW zip itu (work owner 04-10-2026), ID dari CZONE bila ada, tanpa syarat aktif.
	if q := sqlCZoneZip("UJI.CZONE", "UJI.RW"); q != "SELECT w.CZ, (SELECT MIN(z.ID) FROM UJI.CZONE z WHERE z.CODE = w.CZ) FROM "+
		"(SELECT MIN(r.CZONE) AS CZ FROM UJI.RW r WHERE r.ZIPCODE = :1 AND r.STS_AKTIF = '1' AND TRIM(r.CZONE) IS NOT NULL) w WHERE w.CZ IS NOT NULL" {
		t.Errorf("czone: %q", q)
	}
	// A187: hanya zip ber-RW aktif ber-CZONE (berlaku untuk hitungan dan halaman - satu kueri dasar).
	q, arg := sqlZipAkumulasi("UJI.RA", "UJI.N", "UJI.RW", "%JAWA%", "%40%")
	if !strings.HasSuffix(q, " WHERE EXISTS (SELECT 1 FROM UJI.RW r WHERE r.ZIPCODE = a.POSTALCODE AND r.STS_AKTIF = '1' AND "+
		`TRIM(r.CZONE) IS NOT NULL) AND UPPER(a.PROVINCENAME) LIKE :1 ESCAPE '\' AND UPPER(a.POSTALCODE) LIKE :2 ESCAPE '\'`) ||
		!strings.Contains(q, "(SELECT MAX(n.NATIONINITIAL) FROM UJI.N n WHERE n.ID = a.IDNATION OR n.NOTE = a.NATIONNAME)") ||
		!strings.HasPrefix(q, "SELECT DISTINCT a.POSTALCODE, a.CITYNAME, a.PROVINCENAME, a.NATIONNAME,") ||
		len(arg) != 2 || arg[0] != "%JAWA%" || arg[1] != "%40%" {
		t.Errorf("zip: %q %v", q, arg)
	}
	// A185: halaman di server - hitungan atas kueri yang sama; halaman urut kelima kolom, OFFSET / FETCH sesudah bind saringan.
	if h := sqlHitungZip(q); h != "SELECT COUNT(*) FROM ("+q+")" {
		t.Errorf("hitung: %q", h)
	}
	if h := sqlHalamanZip(q, len(arg)); h != q+" ORDER BY 1, 2, 3, 4, 5 OFFSET :3 ROWS FETCH NEXT :4 ROWS ONLY" {
		t.Errorf("halaman: %q", h)
	}
	if q, arg := sqlZipAkumulasi("UJI.RA", "UJI.N", "UJI.RW", "", ""); strings.Contains(q, "UPPER(") ||
		!strings.HasSuffix(q, " a WHERE EXISTS (SELECT 1 FROM UJI.RW r WHERE r.ZIPCODE = a.POSTALCODE AND r.STS_AKTIF = '1' AND TRIM(r.CZONE) IS NOT NULL)") ||
		len(arg) != 0 ||
		sqlHalamanZip(q, 0) != q+" ORDER BY 1, 2, 3, 4, 5 OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY" {
		t.Errorf("zip tanpa saringan: %q %v", q, arg)
	}
	if daftarSaran["accumtype"].ekstra != "NOTE" {
		t.Error("ekstra accumtype .Note")
	}
}
