package models_test

import (
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/riratelife/backend/models"
)

func pesanGalat(g []models.GalatBaris) string {
	var s []string
	for _, x := range g {
		s = append(s, x.Kalimat())
	}
	return strings.Join(s, " | ")
}

// Format excel b5305: USEDBY, CONTRACT, GENDER, AGE, RATE. Pemisah `;` - RATE berkoma desimal boleh tanpa kutip.
func TestUraiCSVTitikKoma(t *testing.T) {
	teks := string(rune(0xFEFF)) + "USEDBY;CONTRACT;GENDER;AGE;RATE\r\n" +
		"UJI RATE A; 1 ;u;05;0,75\r\n" +
		"\r\n" +
		"UJI RATE A;;M;30;1.5\r\n" +
		";;;;\r\n" +
		"uji rate b;10;F;120;2\r\n"
	b, g, err := models.UraiCSV(teks)
	if err != nil || len(g) != 0 {
		t.Fatalf("galat %v %s", err, pesanGalat(g))
	}
	mau := []models.BarisCSV{
		{Baris: 2, UsedBy: "UJI RATE A", Contract: "1", Gender: "U", Age: "5", Rate: "0,75"},
		{Baris: 4, UsedBy: "UJI RATE A", Contract: "", Gender: "M", Age: "30", Rate: "1,5"},
		{Baris: 6, UsedBy: "uji rate b", Contract: "10", Gender: "F", Age: "120", Rate: "2"},
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

// Pemisah `,`: kolom kepala dalam urutan apa pun; RATE berkoma desimal WAJIB dikutip - tanpa kutip ditolak per baris.
func TestUraiCSVKomaDanKutip(t *testing.T) {
	teks := "rate,gender,age,contract,usedby\n" +
		"\"0,125\",U,1,2,UJI RATE A\n" +
		"0,125,U,2,2,UJI RATE A\n" +
		"0.125,U,3,2,UJI RATE A\n"
	b, g, err := models.UraiCSV(teks)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 2 || b[0].Rate != "0,125" || b[1].Rate != "0,125" || b[0].UsedBy != "UJI RATE A" || b[1].Age != "3" {
		t.Errorf("baris %+v", b)
	}
	if len(g) != 1 || g[0].Baris != 3 || !strings.Contains(g[0].Pesan, "has 6 columns, expected 5") ||
		!strings.Contains(g[0].Pesan, "decimal comma must be quoted") {
		t.Errorf("galat %s", pesanGalat(g))
	}
}

func TestUraiCSVGalatPerBaris(t *testing.T) {
	teks := "USEDBY;CONTRACT;GENDER;AGE;RATE\n" +
		";1;U;1;1\n" + // 2 nama kosong
		"UJI;1;X;1;1\n" + // 3 gender
		"UJI;1;U;;1\n" + // 4 age kosong
		"UJI;121;U;1;1\n" + // 5 contract > 120
		"UJI;1;U;1,5;1\n" + // 6 age pecahan
		"UJI;1;U;2;abc\n" + // 7 rate bukan angka
		"UJI;1;U;3;1,234.5\n" + // 8 rate koma dan titik
		"UJI;1;U;4;-1\n" + // 9 rate negatif
		"UJI;1;U;5;\n" + // 10 rate kosong
		"UJI;1;U;6;1\n" + // 11 sah
		"uji ;1;u;6;2\n" + // 12 kembar baris 11 (tanpa beda huruf)
		"UJI;;U;6;2\n" + // 13 sah: CONTRACT kosong beda kunci
		strings.Repeat("A", models.BatasNama+1) + ";1;U;7;1\n" // 14 nama terlalu panjang
	b, g, err := models.UraiCSV(teks)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 2 || b[0].Baris != 11 || b[1].Baris != 13 {
		t.Errorf("baris sah %+v", b)
	}
	mau := map[int]string{
		2:  "USEDBY is required",
		3:  "GENDER must be U, M, or F",
		4:  "AGE is required",
		5:  "CONTRACT must be a whole number from 0 to 120",
		6:  "AGE must be a whole number from 0 to 120",
		7:  "RATE must be a number",
		8:  "RATE must be a number",
		9:  "RATE must be a number",
		10: "RATE is required",
		12: "duplicates row 11",
		14: "USEDBY is longer than 200 characters",
	}
	if len(g) != len(mau) {
		t.Errorf("jumlah galat %d: %s", len(g), pesanGalat(g))
	}
	for _, x := range g {
		if m, ok := mau[x.Baris]; !ok || !strings.Contains(x.Pesan, m) {
			t.Errorf("baris %d: %q, mau %q", x.Baris, x.Pesan, m)
		}
	}
	if k := (models.GalatBaris{Baris: 3, Pesan: "GENDER must be U, M, or F"}).Kalimat(); k != "Row 3: GENDER must be U, M, or F." {
		t.Errorf("kalimat %q", k)
	}
}

func TestUraiCSVGalatBerkas(t *testing.T) {
	kasus := map[string]string{
		"":                                    "the CSV file is empty",
		"\n\n":                                "the CSV file is empty",
		"USEDBY;CONTRACT;GENDER;AGE\nA;1;U;1": "header must be USEDBY, CONTRACT, GENDER, AGE, RATE",
		"USEDBY;CONTRACT;GENDER;AGE;RATE;X\n": "header must be USEDBY, CONTRACT, GENDER, AGE, RATE",
		"USEDBY;USEDBY;GENDER;AGE;RATE\n":     "header must be USEDBY, CONTRACT, GENDER, AGE, RATE",
		"USEDBY;CONTRACT;GENDER;AGE;RATE\n":   "the CSV file has no data rows",
		"USEDBY,CONTRACT,GENDER,AGE,RATE\n\"A,1,U,1,1\n": "cannot be read",
	}
	for teks, mau := range kasus {
		_, _, err := models.UraiCSV(teks)
		if !errors.Is(err, models.ErrCSV) || !strings.Contains(err.Error(), mau) {
			t.Errorf("%q: %v, mau %q", teks, err, mau)
		}
	}
	banyak := "USEDBY;CONTRACT;GENDER;AGE;RATE\n" + strings.Repeat("A;1;U;1;1\n", models.MaksBarisCSV+1)
	if _, _, err := models.UraiCSV(banyak); err == nil || !strings.Contains(err.Error(), "at most 10000 rows") {
		t.Errorf("terlalu banyak: %v", err)
	}
}

func TestNormalRate(t *testing.T) {
	for masuk, mau := range map[string]string{"0,75": "0,75", " 1.50 ": "1,50", "007,5": "7,5", "0": "0", "12": "12", "0.000123": "0,000123"} {
		if got, err := models.NormalRate(masuk); err != nil || got != mau {
			t.Errorf("%q -> %q %v, mau %q", masuk, got, err, mau)
		}
	}
	for _, salah := range []string{"", "1e3", "NaN", "1,2,3", ",5", "5,", "+1", "1 000", "1.234,5"} {
		if got, err := models.NormalRate(salah); err == nil {
			t.Errorf("%q diterima: %q", salah, got)
		}
	}
}

// Kunci kembar (USEDBY, GENDER, AGE, CONTRACT): tanpa beda huruf/spasi tepi; angka warisan berawalan nol setara.
func TestKunciRate(t *testing.T) {
	a := models.KunciRate(" uji ", "u", "05", "")
	if a != models.KunciRate("UJI", "U", "5", " ") {
		t.Errorf("kunci %q", a)
	}
	if a == models.KunciRate("UJI", "U", "5", "0") {
		t.Error("CONTRACT kosong dan 0 harus beda")
	}
}
