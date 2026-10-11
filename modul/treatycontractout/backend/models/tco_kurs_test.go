package models

import (
	"errors"
	"os"
	"path/filepath"
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

func TestUraiNilaiKursTCO(t *testing.T) {
	// format Indonesia dari menu Treaty Exchange Yearly diterima (keputusan work owner 11-10-2026)
	for teks, mau := range map[string]string{"15500": "15500", "15500.25": "15500.25", "15500,25": "15500.25",
		"15.500,25": "15500.25", "16.000,00": "16000.00", "1.234.567": "1234567", " 14.000,00 ": "14000.00"} {
		d, err := UraiNilaiKursTCO(teks)
		if err != nil || d.Text('f') != mau {
			t.Errorf("%q: %v %v", teks, d, err)
		}
	}
	for _, buruk := range []string{"", "0", "-1", "abc", "15,500.25", "1.5.5,0", "16.00,00"} {
		if _, err := UraiNilaiKursTCO(buruk); !errors.Is(err, ErrKursTakTerurai) {
			t.Errorf("%q: %v", buruk, err)
		}
	}
}

// AC 46: kurs yang BERLAKU menurut Oracle (`BETWEEN` di SQL). Lanjutan 6:
// baris yang tanggalnya ditolak Oracle dicacah dan TIDAK mematikan baris
// berlaku; tanpa baris berlaku, penolakan itu galat berkata-kata.
func TestPilihKursBerlakuTCO(t *testing.T) {
	berlaku := KursTCO{ToIDR: apd.New(15500, 0), Mulai: tglKurs("2026-01-01"), Akhir: tglKurs("2026-12-31")}
	ditolak := []BarisKursDitolakTCO{{Kolom: "STARTDATE", Teks: "2019A801T000000.000 GMT"}, {Kolom: "ENDDATE", Teks: ""}}
	k, err := PilihKursBerlakuTCO(HasilMasterKursTCO{Berlaku: []KursTCO{berlaku}, Ditolak: ditolak}, tglKurs("2026-06-30"))
	if err != nil || k.ToIDR.Text('f') != "15500" || k.BarisDitolak != 2 {
		t.Errorf("berlaku + dua ditolak: %+v %v", k, err)
	}
	// ADR-0015: periode tanpa kurs = kegagalan terlihat, bukan nol.
	if _, err := PilihKursBerlakuTCO(HasilMasterKursTCO{}, tglKurs("2027-01-01")); !errors.Is(err, ErrKursTidakAda) {
		t.Errorf("tanpa kurs: %v", err)
	}
	// Tanpa baris berlaku TETAPI ada penolakan: salah satunya mungkin baris
	// yang dicari - master rusak, bukan "tidak ada kurs".
	_, err = PilihKursBerlakuTCO(HasilMasterKursTCO{Ditolak: ditolak}, tglKurs("2019-08-01"))
	if !errors.Is(err, ErrKursTakTerurai) || errors.Is(err, ErrKursTidakAda) ||
		!strings.Contains(err.Error(), `STARTDATE "2019A801T000000.000 GMT"`) || !strings.Contains(err.Error(), "2 master rows") ||
		!strings.Contains(err.Error(), FormatTanggalKursTCO) {
		t.Errorf("hanya ditolak: %v", err)
	}
}

// Dua baris atau lebih berlaku [keputusan work owner 29-09-2026, "Kembar
// identik = satu kurs", mempersempit OQ-TCO-18]: baris IDENTIK - teks TOIDR,
// hari mulai, hari akhir sama - adalah satu kurs (Pega "terakhir menang"
// memberi nilai yang sama, tanpa menebak); selain itu tetap master rusak.
// Data DEV yang melahirkannya: dua pasang baris kembar persis.
func TestPilihKursBerlakuTCOBarisKembar(t *testing.T) {
	baris := func(teks, mulai, akhir string) KursTCO {
		return KursTCO{TeksToIDR: teks, Mulai: tglKurs(mulai), Akhir: tglKurs(akhir)}
	}
	a := baris("16500.00", "2025-07-01", "2026-06-30")

	k, err := PilihKursBerlakuTCO(HasilMasterKursTCO{Berlaku: []KursTCO{a, a}}, tglKurs("2026-06-01"))
	if err != nil || k.TeksToIDR != "16500.00" || k.BarisKembar != 1 {
		t.Errorf("kembar persis: %+v %v", k, err)
	}
	if k, err := PilihKursBerlakuTCO(HasilMasterKursTCO{Berlaku: []KursTCO{a, a, a}}, tglKurs("2026-06-01")); err != nil ||
		k.BarisKembar != 2 {
		t.Errorf("tiga kembar persis: %+v %v", k, err)
	}

	// TIDAK identik - input terbaru (DATEIU) menang, yang kalah dicacah (keputusan work owner 11-10-2026, menggantikan
	// "tolak, jangan menebak" 29-09-2026).
	input := func(k KursTCO, id, waktu string) KursTCO {
		k.ID = id
		if waktu != "" {
			t0, err := time.Parse("2006-01-02 15:04", waktu)
			if err != nil {
				t.Fatal(err)
			}
			k.DiinputPada = t0
		}
		return k
	}
	lama := input(a, "10112", "2025-11-02 20:24")
	for nama, b := range map[string]KursTCO{
		"TOIDR berbeda":               baris("16600", "2025-07-01", "2026-06-30"),
		"angka sama, teks lain":       baris("16500", "2025-07-01", "2026-06-30"),
		"TOIDR sama, periode berbeda": baris("16500.00", "2026-06-01", "2027-05-31"),
		"TOIDR sama, akhir berbeda":   baris("16500.00", "2025-07-01", "2026-12-31"),
	} {
		baru := input(b, "10140", "2026-10-09 08:27")
		for _, urut := range [][]KursTCO{{lama, baru}, {baru, lama}} {
			k, err := PilihKursBerlakuTCO(HasilMasterKursTCO{Berlaku: urut}, tglKurs("2026-06-15"))
			if err != nil || k.TeksToIDR != b.TeksToIDR || k.ID != "10140" || k.BarisTertimpa != 1 || k.BarisKembar != 0 {
				t.Errorf("%s: harap input terbaru (10140) dipakai, satu tertimpa: %+v %v", nama, k, err)
			}
		}
	}
	// waktu input sama -> ID terbesar (angka, bukan teks: 10140 > 9999); DATEIU kosong = paling lama
	x := input(baris("15000", "2024-04-01", "2025-03-31"), "9999", "2026-10-09 08:27")
	y := input(baris("15500", "2024-04-01", "2025-03-31"), "10140", "2026-10-09 08:27")
	if k, err := PilihKursBerlakuTCO(HasilMasterKursTCO{Berlaku: []KursTCO{x, y}}, tglKurs("2024-07-01")); err != nil || k.ID != "10140" {
		t.Errorf("seri waktu input -> ID terbesar: %+v %v", k, err)
	}
	tanpa := input(baris("17000", "2024-04-01", "2025-03-31"), "99999", "")
	if k, err := PilihKursBerlakuTCO(HasilMasterKursTCO{Berlaku: []KursTCO{tanpa, x}}, tglKurs("2024-07-01")); err != nil || k.ID != "9999" {
		t.Errorf("DATEIU kosong kalah dari yang berwaktu: %+v %v", k, err)
	}
	// dua baris lama kembar + satu berbeda yang lebih baru: yang baru dipakai, dua lama tertimpa
	tigaBaris := []KursTCO{lama, lama, input(baris("16600", "2025-07-01", "2026-06-30"), "10200", "2026-10-10 12:00")}
	if k, err := PilihKursBerlakuTCO(HasilMasterKursTCO{Berlaku: tigaBaris}, tglKurs("2026-06-15")); err != nil ||
		k.TeksToIDR != "16600" || k.BarisTertimpa != 2 {
		t.Errorf("dua lama kalah oleh satu baru: %+v %v", k, err)
	}
}

// Topeng VERBATIM `RDBList/GetMasterKursList.xml` b85-b86 - Oracle yang mengurai.
func TestFormatTanggalKursTCO(t *testing.T) {
	if FormatTanggalKursTCO != `YYYYMMDD"T"HH24MISS.FF3 TZR` {
		t.Errorf("%s", FormatTanggalKursTCO)
	}
}

// Pesan VERBATIM `NewTreatyArrEpi.xml` b870.
func TestGalatKursTidakAdaTCO(t *testing.T) {
	err := GalatKursTidakAda{TreatyYear: "2026", Tanggal: tglKurs("2026-01-01")}
	if err.Error() != "No exchange rate for Treaty Year : 2026" || !errors.Is(err, ErrKursTidakAda) {
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

// Temuan /code-review: modul tidak membuat konteks apd sendiri dan tidak
// menulis konteks bersama - `utils.DecimalContext()` memberi salinan.
func TestModulTanpaKonteksApdSendiri(t *testing.T) {
	for _, pola := range []string{"tco_*.go", "../services/tco_*.go", "../repository/tco_*.go"} {
		berkas, _ := filepath.Glob(pola)
		for _, b := range berkas {
			if strings.HasSuffix(b, "_test.go") {
				continue
			}
			isi, err := os.ReadFile(b)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(isi), "apd.BaseContext") || strings.Contains(string(isi), "WithPrecision(") {
				t.Errorf("%s membuat apd.Context sendiri", b)
			}
		}
	}
}
