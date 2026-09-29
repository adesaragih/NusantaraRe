package models

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"
)

func tglKurs(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

func desKurs(t *testing.T, s string) *apd.Decimal {
	t.Helper()
	d, _, err := apd.NewFromString(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

// AC 53: tanggal kurs diurai SEKALI ke tanggal - bentuk warisan
// 'YYYYMMDD"T"HH24MISS.FF3 TZR', dipangkas ke tanggalnya (`trunc`).
func TestUraiTanggalKursTCO(t *testing.T) {
	for teks, mau := range map[string]string{
		"20260101T000000.000 GMT":          "2026-01-01",
		"20261231T235959.999 Asia/Jakarta": "2026-12-31",
		"20260315T120000.000":              "2026-03-15",
	} {
		got, err := UraiTanggalKursTCO("STARTDATE", teks)
		if err != nil || !got.Equal(tglKurs(mau)) {
			t.Errorf("%q: %v %v, mau %s", teks, got, err, mau)
		}
	}
	for _, buruk := range []string{"", "2026-01-01", "20261301T000000.000 GMT", "20260101", "20260101T0000"} {
		if _, err := UraiTanggalKursTCO("ENDDATE", buruk); !errors.Is(err, ErrKursTakTerurai) ||
			!strings.Contains(err.Error(), "ENDDATE") {
			t.Errorf("%q: %v", buruk, err)
		}
	}
}

func TestUraiNilaiKursTCO(t *testing.T) {
	for teks, mau := range map[string]string{"15500": "15500", "15500.25": "15500.25", "15500,25": "15500.25"} {
		d, err := UraiNilaiKursTCO(teks)
		if err != nil || d.Text('f') != mau {
			t.Errorf("%q: %v %v", teks, d, err)
		}
	}
	for _, buruk := range []string{"", "0", "-1", "15.500,25", "abc"} {
		if _, err := UraiNilaiKursTCO(buruk); !errors.Is(err, ErrKursTakTerurai) {
			t.Errorf("%q: %v", buruk, err)
		}
	}
}

// AC 46: kurs yang BERLAKU - tanggal di antara mulai dan akhir, inklusif.
func TestPilihKursBerlakuTCO(t *testing.T) {
	baris := []KursTCO{
		{ToIDR: apd.New(15000, 0), Mulai: tglKurs("2025-01-01"), Akhir: tglKurs("2025-12-31")},
		{ToIDR: apd.New(15500, 0), Mulai: tglKurs("2026-01-01"), Akhir: tglKurs("2026-12-31")},
	}
	for tgl, mau := range map[string]string{"2026-01-01": "15500", "2026-12-31": "15500", "2025-06-30": "15000"} {
		k, err := PilihKursBerlakuTCO(baris, tglKurs(tgl))
		if err != nil || k.ToIDR.Text('f') != mau {
			t.Errorf("%s: %+v %v", tgl, k, err)
		}
	}
	// ADR-0015: periode tanpa kurs = kegagalan terlihat, bukan nol.
	if _, err := PilihKursBerlakuTCO(baris, tglKurs("2027-01-01")); !errors.Is(err, ErrKursTidakAda) {
		t.Errorf("tanpa kurs: %v", err)
	}
	ganda := append(baris, KursTCO{ToIDR: apd.New(1, 0), Mulai: tglKurs("2026-06-01"), Akhir: tglKurs("2026-06-30")})
	if _, err := PilihKursBerlakuTCO(ganda, tglKurs("2026-06-15")); !errors.Is(err, ErrKursGanda) {
		t.Errorf("dua baris berlaku: %v", err)
	}
}

// Pesan VERBATIM `NewTreatyArrEpi.xml` b870.
func TestGalatKursTidakAdaTCO(t *testing.T) {
	err := GalatKursTidakAda{TreatyYear: "2026", Tanggal: tglKurs("2026-01-01")}
	if err.Error() != "Tidak ada Nilai Kurs di Tahun : 2026" || !errors.Is(err, ErrKursTidakAda) {
		t.Errorf("%q", err.Error())
	}
}

// ADR-0003: konversi desimal persis. Kasus besar menangkap `float64`.
func TestKonversiKursTCO(t *testing.T) {
	kurs := desKurs(t, "15500.25")
	usd, err := UsdDariRpTCO(desKurs(t, "31000500"), kurs, SkalaUsdDariRpTCO)
	if err != nil || usd.Text('f') != "2000.00000000" {
		t.Errorf("Usd = Rp / Kurs: %v %v", usd, err)
	}
	usd, _ = UsdDariRpTCO(desKurs(t, "1"), desKurs(t, "3"), SkalaUsdDariRpTCO)
	if usd.Text('f') != "0.33333333" {
		t.Errorf("pembulatan 8: %v", usd)
	}
	usd, _ = UsdDariRpTCO(desKurs(t, "2"), desKurs(t, "3"), SkalaUsdExclusionTCO)
	if usd.Text('f') != "0.6667" {
		t.Errorf("pembulatan 4 exclusion: %v", usd)
	}
	rp, err := RpDariUsdTCO(desKurs(t, "123456789012345678.12345678"), desKurs(t, "16234.56789012"))
	if err != nil || rp.Text('f') != "2004267622717146774324.82837777" {
		t.Errorf("Rp = Usd x Kurs besar: %v %v", rp, err)
	}
	usd, _ = UsdDariRpTCO(desKurs(t, "2004267622717146774324.82837777"), desKurs(t, "16234.56789012"), SkalaUsdDariRpTCO)
	if usd.Text('f') != "123456789012345678.12345678" {
		t.Errorf("Usd besar: %v", usd)
	}
	if _, err := UsdDariRpTCO(desKurs(t, "1"), nil, SkalaUsdDariRpTCO); !errors.Is(err, ErrKursTidakAda) {
		t.Errorf("tanpa kurs: %v", err)
	}
	if u, err := UsdDariRpTCO(nil, kurs, SkalaUsdDariRpTCO); u != nil || err != nil {
		t.Errorf("Rp kosong: %v %v", u, err)
	}
}

// AC 47: identitas mata uang USD bukan konstanta program; AC 50: nama jujur.
func TestKursTanpaLiteralMataUangDanNamaJujur(t *testing.T) {
	literal := regexp.MustCompile(`IDCURRENCY\s*=\s*'`)
	for _, berkas := range []string{"tco_kurs.go", "../repository/tco_kurs.go", "../services/tco_kurs.go", "../handlers/tco_kurs.go"} {
		isi, err := os.ReadFile(berkas)
		if err != nil {
			t.Fatal(err)
		}
		s := string(isi)
		if strings.Contains(s, "10001") || literal.MatchString(s) {
			t.Errorf("%s menanam identitas mata uang sebagai literal", berkas)
		}
		if strings.Contains(strings.ToLower(s), "testing") && !strings.Contains(s, "`testingKurs`") {
			t.Errorf("%s memakai kata testing di luar rujukan korpus", berkas)
		}
		for _, larang := range []string{"float32", "float64", "ParseFloat"} {
			if strings.Contains(s, larang) {
				t.Errorf("%s memuat %s", berkas, larang)
			}
		}
	}
}
