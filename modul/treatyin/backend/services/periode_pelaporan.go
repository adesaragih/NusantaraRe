package services

// Tombol `Apply` tab Reporting Period — `Activity/TreatyInSetReport.xml`.
//
// ⛔ RUMUSNYA DISALIN dari Activity, langkah demi langkah, dan DIUKUR ULANG
// terhadap seluruh 4.548 baris `ReportingPeriodList` yang Pega simpan di
// 1.137 dokumen (6 Oktober 2026): 4.524 baris cocok (99,5%), jumlah baris
// cocok pada 1.135 kontrak. Sisanya kontrak yang Start/Period-nya diubah
// SESUDAH Apply — label `T 1` pada periode `quarter`, misalnya.
//
// Tombolnya (`Section/TreatyInTabsProportional.xml` @254615) menjalankan
// `TreatyInSetReport` dengan `startdate = TreatyIn.ReportingStart` dan
// `autocalculate = true`:
//
//	[4–7]  Property-Set-Messages bila Start/End, Submission, Confirmation,
//	       atau Settlement kosong — lalu KELUAR (nol baris)
//	[8]    duration = @DateTimeDifference(Start, End, 'M')
//	[9–12] quarter 3 · half 6 · month 1 · other = ReportingInterval;
//	       duration = duration / interval + 1; simbol "Q " "H " "M " "T "
//	[13]   baris 1:  InitialDate = Start; TempDate = Start + interval bulan
//	[14]   REPEAT 0…duration-2: InitialDate = TempDate; TempDate += interval
//	       Due = TempDate + Submission; + Confirmation; + Settlement (kumulatif)
//
// ⚠️ Syarat `ReportingPeriod=="other"` pada langkah 13–14 TIDAK AKTIF
// (`pyStepsPreCondition = false`) — keduanya berjalan untuk SEMUA periode.
//
// # Kenapa jatuh tempo "kurang satu hari"
//
// Pega menghitung di zona waktu lokal (WIB) lalu menyimpan jatuh tempo
// sebagai Date dari jam GMT-nya: tengah malam WIB = 17:00 GMT hari
// SEBELUMNYA. Lampiran pemakai: Start 06/10/2026 + 3 bulan + 12 hari =
// 18/01/2027 WIB → tersimpan dan tampil **17/01/2027**. Rumus di bawah
// menirunya (`- 1` hari) — itu yang membuat 4.524 baris cocok.

import (
	"strconv"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
)

// MasukanPeriodePelaporan - isian kepala tab Reporting Period.
type MasukanPeriodePelaporan struct {
	Mulai      string `json:"mulai"` // DD-MM-YYYY atau YYYYMMDD
	Akhir      string `json:"akhir"`
	Periode    string `json:"periode"` // quarter · half · month · other
	Interval   string `json:"interval"`
	Penyerahan string `json:"penyerahan"`
	Konfirmasi string `json:"konfirmasi"`
	Pelunasan  string `json:"pelunasan"`
}

// HasilPeriodePelaporan - baris hasil Apply, atau pesan per medan.
type HasilPeriodePelaporan struct {
	Baris []models.BarisPeriodeWarisan `json:"baris"`
	// Pesan per medan — `Property-Set-Messages` langkah 4–7, teksnya apa adanya.
	Galat map[string]string `json:"galat"`
}

// Pesan langkah 3 Activity — DISALIN apa adanya.
const (
	galatMulaiAkhir = "Start date and End date must not be empty"
	galatPenyerahan = "Submission days must not be empty"
	galatKonfirmasi = "Confirmation days must not be empty"
	galatPelunasan  = "Settlement days must not be empty"
	// ⚠️ Activity tidak memeriksa Interval; tanpanya pembagian langkah 12
	// gagal di Pega. Kalimatnya dirangkai dari label layar sendiri
	// (`, and Interval` @276424 + `. Must Not Be Empty` @281471).
	galatInterval = "Interval must not be empty"
)

// intervalPeriode - langkah 9–11.
var intervalPeriode = map[string]int{"quarter": 3, "half": 6, "month": 1}

// simbolPeriode - `local.periodsymbol`.
var simbolPeriode = map[string]string{"quarter": "Q ", "half": "H ", "month": "M ", "other": "T "}

// tanggalMasukan membaca DD-MM-YYYY (kabel), YYYYMMDD, atau stempel Pega.
func tanggalMasukan(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	for _, pola := range []string{"02-01-2006", "20060102"} {
		if t, err := time.Parse(pola, s); err == nil {
			return t, true
		}
	}
	if len(s) >= 8 {
		if t, err := time.Parse("20060102", s[:8]); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// tambahBulan - `@addCalendar(…, months)`: hari dipotong ke akhir bulan.
func tambahBulan(t time.Time, n int) time.Time {
	awal := time.Date(t.Year(), t.Month()+time.Month(n), 1, 0, 0, 0, 0, time.UTC)
	akhir := awal.AddDate(0, 1, -1).Day()
	hari := t.Day()
	if hari > akhir {
		hari = akhir
	}
	return time.Date(awal.Year(), awal.Month(), hari, 0, 0, 0, 0, time.UTC)
}

// selisihBulan - `@DateTimeDifference(a, b, 'M')`: bulan PENUH.
func selisihBulan(a, b time.Time) int {
	m := (b.Year()-a.Year())*12 + int(b.Month()) - int(a.Month())
	if b.Day() < a.Day() {
		m--
	}
	return m
}

// HitungPeriodePelaporan - tombol Apply. Murni: nol basis data.
func HitungPeriodePelaporan(m MasukanPeriodePelaporan) HasilPeriodePelaporan {
	h := HasilPeriodePelaporan{Baris: []models.BarisPeriodeWarisan{}, Galat: map[string]string{}}
	mulai, okMulai := tanggalMasukan(m.Mulai)
	akhir, okAkhir := tanggalMasukan(m.Akhir)
	if !okMulai || !okAkhir {
		h.Galat["mulai"] = galatMulaiAkhir
		h.Galat["akhir"] = galatMulaiAkhir
	}
	hari := map[string]int{}
	for _, x := range []struct{ kunci, nilai, pesan string }{
		{"penyerahan", m.Penyerahan, galatPenyerahan},
		{"konfirmasi", m.Konfirmasi, galatKonfirmasi},
		{"pelunasan", m.Pelunasan, galatPelunasan},
	} {
		n, err := strconv.Atoi(strings.TrimSpace(x.nilai))
		if err != nil {
			h.Galat[x.kunci] = x.pesan
			continue
		}
		hari[x.kunci] = n
	}
	periode := strings.TrimSpace(m.Periode)
	interval, dikenal := intervalPeriode[periode]
	if !dikenal {
		n, err := strconv.Atoi(strings.TrimSpace(m.Interval))
		if err != nil || n <= 0 {
			h.Galat["interval"] = galatInterval
		}
		interval = n
	}
	// ⛔ Langkah 4–7: pesan lalu KELUAR — nol baris, bukan baris setengah.
	if len(h.Galat) > 0 {
		return h
	}
	jumlah := selisihBulan(mulai, akhir)/interval + 1
	simbol, ada := simbolPeriode[periode]
	if !ada {
		simbol = "T "
	}
	format := func(t time.Time) string { return t.Format("20060102") }
	temp := mulai
	for i := 1; i <= jumlah; i++ {
		awal := temp
		temp = tambahBulan(temp, interval)
		// `- 1` hari: tanggal GMT dari tengah malam WIB — lihat kepala berkas.
		jatuh := func(n int) string { return format(temp.AddDate(0, 0, n-1)) }
		b := models.BarisPeriodeWarisan{
			Periode:          simbol + strconv.Itoa(i),
			HitungOtomatis:   "",
			TanggalAwal:      format(awal),
			JatuhTempoKirim:  jatuh(hari["penyerahan"]),
			JatuhTempoKonfir: jatuh(hari["penyerahan"] + hari["konfirmasi"]),
			JatuhTempoBayar:  jatuh(hari["penyerahan"] + hari["konfirmasi"] + hari["pelunasan"]),
		}
		h.Baris = append(h.Baris, tampilkanPeriode(b))
	}
	return h
}

// tampilkanPeriode menyimpan bentuk TERSIMPAN di `…Asli` lalu menerjemahkan
// bentuk tampilnya — penerjemah yang SAMA dengan baris dokumen.
func tampilkanPeriode(p models.BarisPeriodeWarisan) models.BarisPeriodeWarisan {
	p.TanggalAwalAsli, p.JatuhTempoKirimAsli = p.TanggalAwal, p.JatuhTempoKirim
	p.JatuhTempoKonfirAsli, p.JatuhTempoBayarAsli = p.JatuhTempoKonfir, p.JatuhTempoBayar
	p.TanggalAwal = TanggalTampil(p.TanggalAwal)
	p.JatuhTempoKirim = TanggalTampil(p.JatuhTempoKirim)
	p.JatuhTempoKonfir = TanggalTampil(p.JatuhTempoKonfir)
	p.JatuhTempoBayar = TanggalTampil(p.JatuhTempoBayar)
	return p
}

// HitungPeriodePelaporan - rute `POST /hitung/periode-pelaporan`.
func (l *Layanan) HitungPeriodePelaporan(p inti.Pelaku, m MasukanPeriodePelaporan) (HasilPeriodePelaporan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilPeriodePelaporan{}, err
	}
	return HitungPeriodePelaporan(m), nil
}

// PeriodePelaporanTampil - label pilihan `Period`.
//
// ⭐ Prompt List Property `ASM-FW-GISFW-Int-TREATY_IN.ReportingPeriod`
// (tangkapan layar rule dari pemakai, 6 Oktober 2026): Standard value →
// Prompt value. Nilai di luar daftar tampil apa adanya.
var labelPeriodePelaporan = map[string]string{
	"quarter": "Quarter Year",
	"half":    "Half Year",
	"month":   "Monthly",
	"other":   "Others",
}

func PeriodePelaporanTampil(tersimpan string) string {
	if l, ada := labelPeriodePelaporan[strings.TrimSpace(tersimpan)]; ada {
		return l
	}
	return tersimpan
}

// TanggalWIB mengubah stempel DateTime Pega (`20181231T170000.000 GMT`)
// menjadi tanggal WIB `YYYYMMDD` — tanggal yang layar Pega tampilkan.
// Nilai tanpa jam dikembalikan delapan digit pertamanya; yang tak terbaca
// apa adanya.
func TanggalWIB(s string) string {
	t := strings.TrimSpace(s)
	if len(t) >= 15 && t[8] == 'T' {
		if x, err := time.Parse("20060102T150405", t[:15]); err == nil {
			return x.Add(7 * time.Hour).Format("20060102")
		}
	}
	if len(t) >= 8 {
		if _, err := time.Parse("20060102", t[:8]); err == nil {
			return t[:8]
		}
	}
	return s
}
