package models_test

import (
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/ricommlife/backend/models"
)

// Rumus prosedur PEGA_M_RICOMM_LIFE b21 dengan nilai DEV: site 1, nomor 44 -> 1000044.
func TestBentukIDRumusProsedur(t *testing.T) {
	for _, k := range []struct{ site, nomor, mau string }{
		{"1", "44", "1000044"},
		{"1", "5", "1000005"},
		{" 2 ", "123456", "2123456"},
		{"1234", "1", "1234000001"},
	} {
		if got, err := models.BentukID(k.site, k.nomor); err != nil || got != k.mau {
			t.Errorf("BentukID(%q, %q) = %q %v, mau %q", k.site, k.nomor, got, err, k.mau)
		}
	}
	for _, k := range []struct{ site, nomor, mau string }{
		{"", "1", "site"},
		{"A", "1", "site"},
		{"1", "", "nomor sequence"},
		{"1", "-1", "nomor sequence"},
		{"1", "1234567", "lebih dari 6 angka"},
		{"12345", "1", "lebih dari 10 karakter"},
	} {
		if _, err := models.BentukID(k.site, k.nomor); !errors.Is(err, models.ErrIDTidakSah) || !strings.Contains(err.Error(), k.mau) {
			t.Errorf("BentukID(%q, %q): %v, mau %q", k.site, k.nomor, err, k.mau)
		}
	}
}

// InboxRIComm: CONTRACT b1923, YEAR b2201, COMM b2412 wajib; YEAR pyMax 4 b2206; COMM muat NUMBER(38,8).
func TestPeriksaIsianKomisi(t *testing.T) {
	got, err := models.PeriksaIsianKomisi(models.IsianKomisi{Contract: " 05 ", Year: "0002", Comm: "007,500"})
	if err != nil || got != (models.IsianKomisi{Contract: "5", Year: "2", Comm: "7.5"}) {
		t.Fatalf("sah %+v %v", got, err)
	}
	if got, err := models.PeriksaIsianKomisi(models.IsianKomisi{Contract: "0", Year: "2026", Comm: "0.12345678"}); err != nil ||
		got.Comm != "0.12345678" || got.Year != "2026" {
		t.Errorf("batas pecahan %+v %v", got, err)
	}
	if _, err := models.PeriksaIsianKomisi(models.IsianKomisi{}); err == nil ||
		err.Error() != "CONTRACT is required; YEAR is required; COMM is required" {
		t.Errorf("kosong: %v", err)
	}
	for _, k := range []struct {
		isi models.IsianKomisi
		mau string
	}{
		{models.IsianKomisi{Contract: "1.5", Year: "1", Comm: "1"}, "CONTRACT must be a whole number from 0 to 99999"},
		{models.IsianKomisi{Contract: "100000", Year: "1", Comm: "1"}, "CONTRACT must be a whole number"},
		{models.IsianKomisi{Contract: "1", Year: "12345", Comm: "1"}, "YEAR must be a whole number of at most 4 digits"},
		{models.IsianKomisi{Contract: "1", Year: "-1", Comm: "1"}, "YEAR must be a whole number"},
		{models.IsianKomisi{Contract: "1", Year: "1", Comm: "1,234.5"}, "COMM must be a number"},
		{models.IsianKomisi{Contract: "1", Year: "1", Comm: "-1"}, "COMM must be a number"},
		{models.IsianKomisi{Contract: "1", Year: "1", Comm: "0.123456789"}, "at most 30 digits before and 8 after"},
		{models.IsianKomisi{Contract: "1", Year: "1", Comm: strings.Repeat("9", 31)}, "at most 30 digits before"},
	} {
		if _, err := models.PeriksaIsianKomisi(k.isi); err == nil || !strings.Contains(err.Error(), k.mau) {
			t.Errorf("%+v: galat %v, mau %q", k.isi, err, k.mau)
		}
	}
}

func TestDesimalKanonik(t *testing.T) {
	for masuk, mau := range map[string]string{"0,75": "0.75", " 1.50 ": "1.5", "007,5": "7.5", "0": "0", "0,000": "0",
		"12": "12", "0.000123": "0.000123", "100": "100"} {
		if got, _, _, ok := models.DesimalKanonik(masuk); !ok || got != mau {
			t.Errorf("%q -> %q %v, mau %q", masuk, got, ok, mau)
		}
	}
	for _, salah := range []string{"", "1e3", "NaN", "1,2,3", ",5", "5,", "+1", "1 000", "1.234,5"} {
		if got, _, _, ok := models.DesimalKanonik(salah); ok {
			t.Errorf("%q diterima: %q", salah, got)
		}
	}
	if models.KunciKomisi("05", " 1") != models.KunciKomisi("5", "1") || models.KunciKomisi("1", "2") == models.KunciKomisi("2", "1") {
		t.Error("KunciKomisi")
	}
}

func pesanGalat(g []models.GalatBaris) string {
	var s []string
	for _, x := range g {
		s = append(s, x.Kalimat())
	}
	return strings.Join(s, " | ")
}

// Format excel b5569: USEDBY, CONTRACT, YEAR, COMM. Pemisah `;` - COMM berkoma desimal boleh tanpa kutip.
func TestUraiCSVTitikKoma(t *testing.T) {
	teks := string(rune(0xFEFF)) + "USEDBY;CONTRACT;YEAR;COMM\r\n" +
		"UJI COMM A; 1 ;05;0,75\r\n" +
		"\r\n" +
		"UJI COMM A;2;1;12.5\r\n" +
		";;;\r\n" +
		"uji comm b;10;2026;20\r\n"
	b, g, err := models.UraiCSV(teks)
	if err != nil || len(g) != 0 {
		t.Fatalf("galat %v %s", err, pesanGalat(g))
	}
	mau := []models.BarisCSV{
		{Baris: 2, UsedBy: "UJI COMM A", Contract: "1", Year: "5", Comm: "0.75"},
		{Baris: 4, UsedBy: "UJI COMM A", Contract: "2", Year: "1", Comm: "12.5"},
		{Baris: 6, UsedBy: "uji comm b", Contract: "10", Year: "2026", Comm: "20"},
	}
	if len(b) != len(mau) {
		t.Fatalf("baris %+v", b)
	}
	for i := range mau {
		if b[i] != mau[i] {
			t.Errorf("baris %d\n%+v\nmau\n%+v", i, b[i], mau[i])
		}
	}
}

// Pemisah `,`: kolom kepala dalam urutan apa pun; COMM berkoma desimal WAJIB dikutip.
func TestUraiCSVKomaDanKutip(t *testing.T) {
	teks := "comm,year,contract,usedby\n" +
		"\"0,125\",1,2,UJI COMM A\n" +
		"0,125,2,2,UJI COMM A\n" +
		"0.125,3,2,UJI COMM A\n"
	b, g, err := models.UraiCSV(teks)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 2 || b[0].Comm != "0.125" || b[1].Comm != "0.125" || b[1].Year != "3" {
		t.Errorf("baris %+v", b)
	}
	if len(g) != 1 || g[0].Baris != 3 || !strings.Contains(g[0].Pesan, "has 5 columns, expected 4") ||
		!strings.Contains(g[0].Pesan, "decimal comma must be quoted") {
		t.Errorf("galat %s", pesanGalat(g))
	}
}

func TestUraiCSVGalatPerBaris(t *testing.T) {
	teks := "USEDBY;CONTRACT;YEAR;COMM\n" +
		";1;1;1\n" + // 2 nama kosong
		"UJI;;1;1\n" + // 3 contract kosong
		"UJI;1;;1\n" + // 4 year kosong
		"UJI;1;12345;1\n" + // 5 year > 4 angka
		"UJI;1,5;1;1\n" + // 6 contract pecahan
		"UJI;1;2;abc\n" + // 7 comm bukan angka
		"UJI;1;3;\n" + // 8 comm kosong
		"UJI;1;6;1\n" + // 9 sah
		"uji ;01;6;2\n" + // 10 kembar baris 9 (tanpa beda huruf, nol depan)
		"UJI;2;6;2\n" + // 11 sah: CONTRACT beda
		strings.Repeat("A", models.BatasNama+1) + ";1;7;1\n" // 12 nama terlalu panjang
	b, g, err := models.UraiCSV(teks)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 2 || b[0].Baris != 9 || b[1].Baris != 11 {
		t.Errorf("baris sah %+v", b)
	}
	mau := map[int]string{
		2:  "USEDBY is required",
		3:  "CONTRACT is required",
		4:  "YEAR is required",
		5:  "YEAR must be a whole number of at most 4 digits",
		6:  "CONTRACT must be a whole number",
		7:  "COMM must be a number",
		8:  "COMM is required",
		10: "duplicates row 9",
		12: "USEDBY is longer than 200 characters",
	}
	if len(g) != len(mau) {
		t.Errorf("jumlah galat %d: %s", len(g), pesanGalat(g))
	}
	for _, x := range g {
		if m, ok := mau[x.Baris]; !ok || !strings.Contains(x.Pesan, m) {
			t.Errorf("baris %d: %q, mau %q", x.Baris, x.Pesan, m)
		}
	}
}

func TestUraiCSVGalatBerkas(t *testing.T) {
	kasus := map[string]string{
		"":                                       "the CSV file is empty",
		"\n\n":                                   "the CSV file is empty",
		"USEDBY;CONTRACT;YEAR\nA;1;1":            "header must be USEDBY, CONTRACT, YEAR, COMM",
		"USEDBY;CONTRACT;YEAR;COMM;X\n":          "header must be USEDBY, CONTRACT, YEAR, COMM",
		"USEDBY;USEDBY;YEAR;COMM\n":              "header must be USEDBY, CONTRACT, YEAR, COMM",
		"USEDBY;CONTRACT;YEAR;COMM\n":            "the CSV file has no data rows",
		"USEDBY,CONTRACT,YEAR,COMM\n\"A,1,1,1\n": "cannot be read",
	}
	for teks, mau := range kasus {
		_, _, err := models.UraiCSV(teks)
		if !errors.Is(err, models.ErrCSV) || !strings.Contains(err.Error(), mau) {
			t.Errorf("%q: %v, mau %q", teks, err, mau)
		}
	}
	banyak := "USEDBY;CONTRACT;YEAR;COMM\n" + strings.Repeat("A;1;1;1\n", models.MaksBarisCSV+1)
	if _, _, err := models.UraiCSV(banyak); err == nil || !strings.Contains(err.Error(), "at most 10000 rows") {
		t.Errorf("terlalu banyak: %v", err)
	}
}
