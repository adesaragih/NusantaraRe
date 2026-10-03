package models

// `DataTransform/SystemSetOneYear_DT`: `@DateTime.addCalendar(.StartDate,1,...)`.
// Fungsi Pega itu menambah tahun lewat `java.util.Calendar.add(YEAR, 1)`, yang
// MENJEPIT tanggal 29 Februari ke 28 Februari pada tahun bukan kabisat - bukan
// menggulirkannya ke 1 Maret seperti `time.AddDate`.

import "testing"

func TestSetahunDari29FebruariJatuhPada28Februari(t *testing.T) {
	h := HalamanBaru()
	h.Setel(HalamanPolis+".StartDate", "2028-02-29")
	SystemSetOneYear(h)
	if got := h.Ambil(HalamanPolis + ".EndDate"); got != "2029-02-28" {
		t.Fatalf("mau 2029-02-28, dapat %q", got)
	}
}
