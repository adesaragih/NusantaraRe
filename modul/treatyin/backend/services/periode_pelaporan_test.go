package services_test

// Uji tombol Apply tab Reporting Period — `Activity/TreatyInSetReport.xml`.

import (
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/services"
)

func masukan(mulai, akhir, periode, sub, con, set string) services.MasukanPeriodePelaporan {
	return services.MasukanPeriodePelaporan{Mulai: mulai, Akhir: akhir, Periode: periode, Penyerahan: sub, Konfirmasi: con, Pelunasan: set}
}

// ⭐ Lampiran pemakai 6 Oktober 2026 — hasil layar Pega:
// Start 06/10/2026, End 06/10/2026, Quarter Year, 12/12/12 →
// Q 1 · 06/10/2026 · 17/01/2027 · 29/01/2027 · 10/02/2027.
func TestApplySamaDenganLayarPega(t *testing.T) {
	h := services.HitungPeriodePelaporan(masukan("06-10-2026", "06-10-2026", "quarter", "12", "12", "12"))
	if len(h.Galat) != 0 || len(h.Baris) != 1 {
		t.Fatalf("galat %v, %d baris — mau 1 baris", h.Galat, len(h.Baris))
	}
	b := h.Baris[0]
	got := [5]string{b.Periode, b.TanggalAwalAsli, b.JatuhTempoKirimAsli, b.JatuhTempoKonfirAsli, b.JatuhTempoBayarAsli}
	mau := [5]string{"Q 1", "20261006", "20270117", "20270129", "20270210"}
	if got != mau {
		t.Errorf("dapat %v, mau %v", got, mau)
	}
	if b.HitungOtomatis != "" {
		t.Errorf("Auto Calculate %q — Activity tidak mengisinya", b.HitungOtomatis)
	}
}

// ⭐ Kontrak nyata 1000019 — baris yang Pega SIMPAN: Start 01/01/2019 WIB,
// End 31/12/2019, quarter, 45/15/15 → empat baris.
func TestApplySamaDenganBarisTersimpan1000019(t *testing.T) {
	h := services.HitungPeriodePelaporan(masukan("20190101", "20191231", "quarter", "45", "15", "15"))
	mau := [][4]string{
		{"Q 1", "20190515", "20190530", "20190614"},
		{"Q 2", "20190814", "20190829", "20190913"},
		{"Q 3", "20191114", "20191129", "20191214"},
		{"Q 4", "20200214", "20200229", "20200315"},
	}
	if len(h.Baris) != len(mau) {
		t.Fatalf("%d baris, mau %d", len(h.Baris), len(mau))
	}
	for i, m := range mau {
		b := h.Baris[i]
		if got := [4]string{b.Periode, b.JatuhTempoKirimAsli, b.JatuhTempoKonfirAsli, b.JatuhTempoBayarAsli}; got != m {
			t.Errorf("baris %d: %v, mau %v", i+1, got, m)
		}
	}
	// InitialDate baris berikutnya = TempDate baris sebelumnya.
	if h.Baris[1].TanggalAwalAsli != "20190401" {
		t.Errorf("InitialDate Q 2 %s, mau 20190401", h.Baris[1].TanggalAwalAsli)
	}
}

// ⛔ Langkah 4–7: pesan per medan, apa adanya, lalu KELUAR — nol baris.
// ⭐ Sel Initial Date — `TreatyInSetReport(startdate=.InitialDate, …)`:
// baris pertama mulai dari tanggal sel itu, JUMLAH baris tetap dari
// ReportingStart→End (langkah 8: `@DateTimeDifference(Start, End, 'M')`).
func TestStartdateSelMenggeserAwalTanpaMengubahJumlah(t *testing.T) {
	m := masukan("20190101", "20191231", "quarter", "45", "15", "15")
	m.Awal = "20190401"
	h := services.HitungPeriodePelaporan(m)
	if len(h.Baris) != 4 {
		t.Fatalf("%d baris, mau 4", len(h.Baris))
	}
	if h.Baris[0].TanggalAwalAsli != "20190401" || h.Baris[3].TanggalAwalAsli != "20200101" {
		t.Errorf("awal %s … %s, mau 20190401 … 20200101", h.Baris[0].TanggalAwalAsli, h.Baris[3].TanggalAwalAsli)
	}
}

func TestApplyMedanKosongMemberiPesanNolBaris(t *testing.T) {
	h := services.HitungPeriodePelaporan(masukan("", "", "quarter", "", "", ""))
	if len(h.Baris) != 0 {
		t.Errorf("%d baris walau isian kosong", len(h.Baris))
	}
	for k, v := range map[string]string{
		"mulai":      "Start date and End date must not be empty",
		"akhir":      "Start date and End date must not be empty",
		"penyerahan": "Submission days must not be empty",
		"konfirmasi": "Confirmation days must not be empty",
		"pelunasan":  "Settlement days must not be empty",
	} {
		if h.Galat[k] != v {
			t.Errorf("%s: %q, mau %q", k, h.Galat[k], v)
		}
	}
}

func TestApplyPeriodeLainDanInterval(t *testing.T) {
	// month: 12 bulan → 13 baris berlabel "M n".
	h := services.HitungPeriodePelaporan(masukan("01-01-2020", "01-01-2021", "month", "1", "1", "1"))
	if len(h.Baris) != 13 || h.Baris[12].Periode != "M 13" {
		t.Errorf("month: %d baris, terakhir %v", len(h.Baris), h.Baris[len(h.Baris)-1].Periode)
	}
	// other tanpa Interval → pesan, nol baris.
	o := services.HitungPeriodePelaporan(masukan("01-01-2020", "01-01-2021", "other", "1", "1", "1"))
	if len(o.Baris) != 0 || o.Galat["interval"] == "" {
		t.Errorf("other tanpa interval: %v", o)
	}
	m := masukan("01-01-2020", "01-01-2021", "other", "1", "1", "1")
	m.Interval = "4"
	if r := services.HitungPeriodePelaporan(m); len(r.Baris) != 4 || r.Baris[0].Periode != "T 1" {
		t.Errorf("other interval 4: %v", r.Baris)
	}
}

// ⚠️ `@addCalendar` memotong ke akhir bulan, dan pemotongannya MENUMPUK
// karena tiap periode ditambah dari TempDate sebelumnya.
func TestApplyAkhirBulanMenumpuk(t *testing.T) {
	h := services.HitungPeriodePelaporan(masukan("20191130", "20200630", "quarter", "1", "1", "1"))
	got := []string{}
	for _, b := range h.Baris {
		got = append(got, b.TanggalAwalAsli)
	}
	mau := []string{"20191130", "20200229", "20200529"}
	if len(got) != len(mau) {
		t.Fatalf("%v, mau %v", got, mau)
	}
	for i := range mau {
		if got[i] != mau[i] {
			t.Errorf("%v, mau %v", got, mau)
		}
	}
}

func TestApplyMenolakTanpaIdentitas(t *testing.T) {
	l := services.LayananDengan(nil)
	if _, err := l.HitungPeriodePelaporan(inti.Pelaku{}, masukan("", "", "", "", "", "")); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("mau ErrTanpaIdentitas, dapat %v", err)
	}
}

func TestTanggalWIBDariStempelPega(t *testing.T) {
	for in, mau := range map[string]string{
		"20181231T170000.000 GMT": "20190101",
		"20190515":                "20190515",
		"bukan":                   "bukan",
	} {
		if got := services.TanggalWIB(in); got != mau {
			t.Errorf("%q -> %q, mau %q", in, got, mau)
		}
	}
}

// ⭐ Prompt List Property ReportingPeriod (rule Pega, lampiran pemakai).
func TestLabelPeriodeDariPromptListPega(t *testing.T) {
	for in, mau := range map[string]string{
		"quarter": "Quarter Year", "half": "Half Year", "month": "Monthly", "other": "Others", "lain": "lain",
	} {
		if got := services.PeriodePelaporanTampil(in); got != mau {
			t.Errorf("%q -> %q, mau %q", in, got, mau)
		}
	}
}
