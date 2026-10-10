package services

// RUMUS TAB ACCUMULATION (cabang PROPORSIONAL) - dari ekspor Pega.
//
// ---------------------------------------------------------------------
// ⭐ KETERGANTUNGAN ANTARTAB
// ---------------------------------------------------------------------
// `TreatyInSetAccountReport` membaca `TreatyIn.ReportingStart` dan
// `TreatyIn.ReportingEnd` - medan tab REPORTING PERIOD, bukan tab ini. Di
// Pega keduanya hidup di halaman `TreatyIn` yang sama; di layar ini
// keduanya dibaca dari penampung halaman (`frontend/halaman.tsx`) dan
// dikirim bersama isian tab Accumulation.
//
// ---------------------------------------------------------------------
// Activity/TreatyInSetAccountReport.xml - isian `Period` (perilaku change)
// ---------------------------------------------------------------------
//
//	 1  startdate = ReportingStart; duration = @DateTimeDifference(Start, End, 'M')
//	 2  AccumulationPeriod == "none"   → KELUAR (daftar TIDAK disentuh)
//	3-5 quarter 3 "Q " · half 6 "H " · month 1 "M ";
//	    duration = duration / interval + 1
//	 6  "other" → Property-Remove AccumulationList, KELUAR
//	 7  Loop = 0
//	 8  Property-Remove AccumulationList
//	 9  (prasyarat `ReportingPeriod=="other"` NONAKTIF) baris 1:
//	    Period = simbol + 1, ReportDate = startdate, TempDate = start + interval bulan
//	10  REPEAT 0..duration-2 (prasyarat sama, NONAKTIF): Period = simbol + n,
//	    ReportDate = TempDate, TempDate += interval bulan
//
// Pembantu tanggalnya SAMA dengan `TreatyInSetReport` (`periode_pelaporan.go`,
// diukur 99,5% terhadap data Pega): `selisihBulan` = @DateTimeDifference 'M'
// bulan penuh, pembagian bilangan bulat, `tambahBulan` = @addCalendar.
//
// ---------------------------------------------------------------------
// Activity/TreatyInAccumulationSetSubDue.xml - sel Reporting Date / Submission Days
// ---------------------------------------------------------------------
//
//	per baris, kecuali `@length(.SubDays) < 1`:
//	    SubDueDate = @addToDate(.ReportDate, .SubDays, '0','0','0')
//
// ⚠️ YANG TIDAK BERBUKTI, DIPUTUSKAN DI SINI:
//   - Start/End kosong atau tak terbaca: @DateTimeDifference atas nilai itu
//     tidak terekspor; di sini selisihnya 0 bulan → SATU baris, ReportDate
//     = Start apa adanya.
//   - AccumulationPeriod di luar kelima nilai (kosong): langkah 3–6 tidak
//     berlaku, jadi `interval`, `periodsymbol` kosong dan `duration` TIDAK
//     dibagi - disalin apa adanya: baris sebanyak selisih bulan, semua
//     bertanggal Start, Period = urutan tanpa simbol.
//   - Submission Days bukan bilangan bulat: SubDueDate baris itu tidak
//     disentuh.
//   - Tanggal ditulis `YYYYMMDD` (tanggal Jakarta). Pega menyimpan
//     ReportDate/SubDueDate sebagai stempel DateTime GMT (panjang 23 di
//     data Pega, `docs/STRUKTUR-TABEL-TREATY-IN.md`) - tanggal yang
//     DITAMPILKAN sama; bentuk simpannya urusan tombol Save.

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
)

// Aksi tab Accumulation.
const (
	// AksiAkumulasiPeriode - isian `Period` (`TreatyInSetAccountReport`).
	AksiAkumulasiPeriode = "periode"
	// AksiAkumulasiJatuhTempo - sel Reporting Date / Submission Days
	// (`TreatyInAccumulationSetSubDue`).
	AksiAkumulasiJatuhTempo = "jatuh-tempo"
)

// BarisAkumulasi - satu baris `TreatyIn.AccumulationList` (ejaan Pega).
type BarisAkumulasi struct {
	Period     string `json:"Period"`
	ReportDate string `json:"ReportDate"`
	SubDays    string `json:"SubDays"`
	SubDueDate string `json:"SubDueDate"`
}

// MasukanAkumulasi - satu aksi tab Accumulation. `ReportingStart` dan
// `ReportingEnd` milik tab Reporting Period.
type MasukanAkumulasi struct {
	Aksi               string           `json:"aksi"`
	AccumulationPeriod string           `json:"AccumulationPeriod"`
	ReportingStart     string           `json:"ReportingStart"`
	ReportingEnd       string           `json:"ReportingEnd"`
	AccumulationList   []BarisAkumulasi `json:"AccumulationList"`
}

// HasilAkumulasi - daftar sesudah aksi.
type HasilAkumulasi struct {
	AccumulationList []BarisAkumulasi `json:"AccumulationList"`
}

// intervalAkumulasi / simbolAkumulasi - langkah 3–5 (TANPA "other": langkah 6
// keluar sebelum interval dipakai).
var (
	intervalAkumulasi = map[string]int{"quarter": 3, "half": 6, "month": 1}
	simbolAkumulasi   = map[string]string{"quarter": "Q ", "half": "H ", "month": "M "}
)

// HitungAkumulasi - bentuk ber-pelaku untuk handler.
func (l *Layanan) HitungAkumulasi(p inti.Pelaku, m MasukanAkumulasi) (HasilAkumulasi, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilAkumulasi{}, err
	}
	return HitungAkumulasi(m)
}

// HitungAkumulasi menjalankan satu aksi. MURNI: nol basis data.
func HitungAkumulasi(m MasukanAkumulasi) (HasilAkumulasi, error) {
	switch m.Aksi {
	case AksiAkumulasiPeriode:
		return HasilAkumulasi{AccumulationList: setAccountReport(m)}, nil
	case AksiAkumulasiJatuhTempo:
		return HasilAkumulasi{AccumulationList: setSubDue(m.AccumulationList)}, nil
	}
	return HasilAkumulasi{}, fmt.Errorf("%w: aksi accumulation %q", ErrMasukanTidakSah, m.Aksi)
}

// setAccountReport = `TreatyInSetAccountReport` langkah 1–10.
func setAccountReport(m MasukanAkumulasi) []BarisAkumulasi {
	periode := strings.TrimSpace(m.AccumulationPeriod)
	switch periode {
	case "none":
		// Langkah 2 - keluar; daftar apa adanya.
		return salinAkumulasi(m.AccumulationList)
	case "other":
		// Langkah 6 - dibuang, lalu keluar.
		return []BarisAkumulasi{}
	}
	mulai, okMulai := tanggalMasukan(m.ReportingStart)
	akhir, okAkhir := tanggalMasukan(m.ReportingEnd)
	durasi := 0
	if okMulai && okAkhir {
		durasi = selisihBulan(mulai, akhir)
	}
	interval, dikenal := intervalAkumulasi[periode]
	if dikenal {
		durasi = durasi/interval + 1
	}
	simbol := simbolAkumulasi[periode]

	tanggal := func(t time.Time) string { return t.Format("20060102") }
	awal := strings.TrimSpace(m.ReportingStart)
	if okMulai {
		awal = tanggal(mulai)
	}
	// Langkah 8-9: daftar dibuang, baris 1 bertanggal Start.
	out := []BarisAkumulasi{{Period: simbol + "1", ReportDate: awal}}
	if !okMulai {
		return out
	}
	temp := tambahBulan(mulai, interval)
	// Langkah 10 - REPEAT 0..duration-2 (batas negatif: nol putaran).
	for i := 0; i <= durasi-2; i++ {
		out = append(out, BarisAkumulasi{Period: simbol + strconv.Itoa(len(out)+1), ReportDate: tanggal(temp)})
		temp = tambahBulan(temp, interval)
	}
	return out
}

// setSubDue = `TreatyInAccumulationSetSubDue`.
func setSubDue(daftar []BarisAkumulasi) []BarisAkumulasi {
	out := salinAkumulasi(daftar)
	for i := range out {
		b := &out[i]
		if len(b.SubDays) < 1 {
			continue
		}
		hari, err := strconv.Atoi(strings.TrimSpace(b.SubDays))
		if err != nil {
			continue
		}
		t, ok := tanggalMasukan(b.ReportDate)
		if !ok {
			continue
		}
		b.SubDueDate = t.AddDate(0, 0, hari).Format("20060102")
	}
	return out
}

func salinAkumulasi(xs []BarisAkumulasi) []BarisAkumulasi {
	out := make([]BarisAkumulasi, len(xs))
	copy(out, xs)
	return out
}
