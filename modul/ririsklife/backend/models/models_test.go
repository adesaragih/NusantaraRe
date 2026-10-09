package models_test

import (
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/ririsklife/backend/models"
)

// Rumus ringkasan PEGA_M_RIRISK_LIFE_SUMMARY: site || LPAD(seq, 6, '0') - nilai DEV: site 1, last_number 166.
func TestBentukIDRingkasanRumusProsedur(t *testing.T) {
	for _, k := range []struct{ site, nomor, mau string }{
		{"1", "166", "1000166"},
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

// Rumus rincian PEGA_M_RIRISK_LIFE: '1' || LPAD(seq, 5, '0') - huruf '1' TETAP (bukan site), lebar 6. Nilai DEV:
// M_RIRISK_LIFE_SEQ last_number 31723 -> 131723. Nomor 6 angka DITOLAK (LPAD Oracle memotongnya -> ID kembar).
func TestBentukIDRincianRumusProsedur(t *testing.T) {
	for nomor, mau := range map[string]string{"31723": "131723", "1": "100001", "99999": "199999", " 44 ": "100044"} {
		if got, err := models.BentukIDRincian(nomor); err != nil || got != mau {
			t.Errorf("BentukIDRincian(%q) = %q %v, mau %q", nomor, got, err, mau)
		}
	}
	for nomor, mau := range map[string]string{"100000": "lebih dari 5 angka", "": "nomor sequence", "1A": "nomor sequence"} {
		if _, err := models.BentukIDRincian(nomor); !errors.Is(err, models.ErrIDTidakSah) || !strings.Contains(err.Error(), mau) {
			t.Errorf("BentukIDRincian(%q): %v, mau %q", nomor, err, mau)
		}
	}
	if len("1"+strings.Repeat("9", models.LebarNomorIDRincian)) != models.BatasIDRincian {
		t.Error("lebar ID rincian harus 6 = RIRISK_LIFE.ID VARCHAR2(6)")
	}
}

// InboxRIRisk: CONTRACT b1933 dan RISK b2595 wajib; YEAR b2216 dan MONTH b2403 TIDAK wajib (pyMax 4); RISK tidak negatif,
// koma atau titik desimal, muat NUMBER (38 angka) - tanpa skala tetap.
func TestPeriksaIsianRincian(t *testing.T) {
	got, err := models.PeriksaIsianRincian(models.IsianRincian{Contract: " 05 ", Year: " 0001 ", Month: "", Risk: "921,9"})
	if err != nil || got != (models.IsianRincian{Contract: "5", Year: "1", Month: "", Risk: "921.9"}) {
		t.Fatalf("sah %+v %v", got, err)
	}
	if got, err := models.PeriksaIsianRincian(models.IsianRincian{Contract: "1", Month: "180", Risk: "580.894351210924"}); err != nil ||
		got.Risk != "580.894351210924" || got.Month != "180" || got.Year != "" {
		t.Errorf("MONTH tanpa YEAR, RISK 12 desimal %+v %v", got, err)
	}
	if _, err := models.PeriksaIsianRincian(models.IsianRincian{}); err == nil || err.Error() != "CONTRACT is required; RISK is required" {
		t.Errorf("kosong: %v", err)
	}
	for _, k := range []struct {
		isi models.IsianRincian
		mau string
	}{
		{models.IsianRincian{Contract: "1.5", Risk: "1"}, "CONTRACT must be a whole number"},
		{models.IsianRincian{Contract: "12345678901", Risk: "1"}, "CONTRACT must be a whole number of at most 10 digits"},
		{models.IsianRincian{Contract: "1", Year: "12345", Risk: "1"}, "YEAR must be a whole number of at most 4 digits"},
		{models.IsianRincian{Contract: "1", Month: "-1", Risk: "1"}, "MONTH must be a whole number"},
		{models.IsianRincian{Contract: "1", Risk: "1,234.5"}, "RISK must be a number"},
		{models.IsianRincian{Contract: "1", Risk: "-0,5"}, "RISK must not be negative"},
		{models.IsianRincian{Contract: "1", Risk: strings.Repeat("9", 39)}, "RISK must have at most 38 digits"},
	} {
		if _, err := models.PeriksaIsianRincian(k.isi); err == nil || !strings.Contains(err.Error(), k.mau) {
			t.Errorf("%+v: galat %v, mau %q", k.isi, err, k.mau)
		}
	}
}

// Konversi RISK = setara migrasi 940 (koma DAN titik sama-sama pemisah desimal): "921,9" -> 921.9,
// "580.894351210924" -> 580.894351210924, bilangan bulat apa adanya.
func TestDesimalKanonikSetaraKonversi940(t *testing.T) {
	for masuk, mau := range map[string]string{"921,9": "921.9", "580.894351210924": "580.894351210924", "12": "12",
		"0,75": "0.75", " 1.50 ": "1.5", "007,5": "7.5", "0": "0", "0,000": "0", "0.000123": "0.000123"} {
		if got, _, _, ok := models.DesimalKanonik(masuk); !ok || got != mau {
			t.Errorf("%q -> %q %v, mau %q", masuk, got, ok, mau)
		}
	}
	for _, salah := range []string{"", "1e3", "NaN", "1,2,3", ",5", "5,", "+1", "1 000", "1.234,5"} {
		if got, _, _, ok := models.DesimalKanonik(salah); ok {
			t.Errorf("%q diterima: %q", salah, got)
		}
	}
	if models.KunciRincian("05", " 1", "") != models.KunciRincian("5", "1", "") ||
		models.KunciRincian("1", "2", "") == models.KunciRincian("1", "", "2") {
		t.Error("KunciRincian: YEAR dan MONTH tidak boleh tertukar")
	}
}

func pesanGalat(g []models.GalatBaris) string {
	var s []string
	for _, x := range g {
		s = append(s, x.Kalimat())
	}
	return strings.Join(s, " | ")
}

// Format excel b5495: USEDBY, CONTRACT, YEAR, MONTH, RISK. Pemisah `;` - RISK berkoma desimal boleh tanpa kutip.
func TestUraiCSVTitikKoma(t *testing.T) {
	teks := string(rune(0xFEFF)) + "USEDBY;CONTRACT;YEAR;MONTH;RISK\r\n" +
		"UJI RISK A; 1 ; 1 ;;921,9\r\n" +
		"\r\n" +
		"UJI RISK A;1;;12;580.894351210924\r\n" +
		";;;;\r\n" +
		"uji risk b;10;2;;20\r\n"
	b, g, err := models.UraiCSV(teks)
	if err != nil || len(g) != 0 {
		t.Fatalf("galat %v %s", err, pesanGalat(g))
	}
	mau := []models.BarisCSV{
		{Baris: 2, UsedBy: "UJI RISK A", Contract: "1", Year: "1", Risk: "921.9"},
		{Baris: 4, UsedBy: "UJI RISK A", Contract: "1", Month: "12", Risk: "580.894351210924"},
		{Baris: 6, UsedBy: "uji risk b", Contract: "10", Year: "2", Risk: "20"},
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

// Pemisah `,`: kolom kepala dalam urutan apa pun; RISK berkoma desimal WAJIB dikutip.
func TestUraiCSVKomaDanKutip(t *testing.T) {
	teks := "risk,month,year,contract,usedby\n" +
		"\"0,125\",,1,2,UJI RISK A\n" +
		"0,125,,2,2,UJI RISK A\n" +
		"0.125,,3,2,UJI RISK A\n"
	b, g, err := models.UraiCSV(teks)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 2 || b[0].Risk != "0.125" || b[1].Risk != "0.125" || b[1].Year != "3" {
		t.Errorf("baris %+v", b)
	}
	if len(g) != 1 || g[0].Baris != 3 || !strings.Contains(g[0].Pesan, "has 6 columns, expected 5") ||
		!strings.Contains(g[0].Pesan, "decimal comma must be quoted") {
		t.Errorf("galat %s", pesanGalat(g))
	}
}

func TestUraiCSVGalatPerBaris(t *testing.T) {
	teks := "USEDBY;CONTRACT;YEAR;MONTH;RISK\n" +
		";1;1;;1\n" + // 2 nama kosong
		"UJI;;1;;1\n" + // 3 contract kosong
		"UJI;1;12345;;1\n" + // 4 year lebih dari 4 angka
		"UJI;1;;1,5;1\n" + // 5 month pecahan
		"UJI;1,5;1;;1\n" + // 6 contract pecahan
		"UJI;1;2;;abc\n" + // 7 risk bukan angka
		"UJI;1;3;;\n" + // 8 risk kosong
		"UJI;1;4;;1\n" + // 9 sah
		"uji ;01;04;;2\n" + // 10 kembar baris 9 (tanpa beda huruf, nol depan)
		"UJI;1;;4;2\n" + // 11 sah: MONTH 4, bukan YEAR 4
		strings.Repeat("A", models.BatasNama+1) + ";1;1;;1\n" + // 12 nama terlalu panjang
		"UJI;3;1;;-1\n" // 13 risk negatif
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
		4:  "YEAR must be a whole number of at most 4 digits",
		5:  "MONTH must be a whole number",
		6:  "CONTRACT must be a whole number",
		7:  "RISK must be a number",
		8:  "RISK is required",
		10: "duplicates row 9",
		12: "USEDBY is longer than 200 characters",
		13: "RISK must not be negative",
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
		"":                                    "the CSV file is empty",
		"\n\n":                                "the CSV file is empty",
		"USEDBY;CONTRACT;YEAR;RISK\nA;1;1;1":  "header must be USEDBY, CONTRACT, YEAR, MONTH, RISK",
		"USEDBY;CONTRACT;YEAR;MONTH;RISK;X\n": "header must be USEDBY, CONTRACT, YEAR, MONTH, RISK",
		"USEDBY;USEDBY;YEAR;MONTH;RISK\n":     "header must be USEDBY, CONTRACT, YEAR, MONTH, RISK",
		"USEDBY;CONTRACT;YEAR;MONTH;RISK\n":   "the CSV file has no data rows",
		"USEDBY,CONTRACT,YEAR,MONTH,RISK\n\"A,1,1,,1\n": "cannot be read",
	}
	for teks, mau := range kasus {
		_, _, err := models.UraiCSV(teks)
		if !errors.Is(err, models.ErrCSV) || !strings.Contains(err.Error(), mau) {
			t.Errorf("%q: %v, mau %q", teks, err, mau)
		}
	}
	banyak := "USEDBY;CONTRACT;YEAR;MONTH;RISK\n" + strings.Repeat("A;1;1;;1\n", models.MaksBarisCSV+1)
	if _, _, err := models.UraiCSV(banyak); err == nil || !strings.Contains(err.Error(), "at most 10000 rows") {
		t.Errorf("terlalu banyak: %v", err)
	}
}
