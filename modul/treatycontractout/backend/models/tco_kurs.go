package models

// Kurs USD -> IDR - tiket 11 Treaty Contract Out.
//
// Untuk apa berkas ini: kurs yang BERLAKU pada tanggal mulai tahun treaty dan
// konversi Rp <-> Usd dengan desimal persis (ADR-0003).
//
// Jalur yang benar-benar dipakai Pega adalah aktivitas `testingKurs` (nol
// langkah di-remark): `CARI1 = Param.StartDate` (b273) -> `GetMasterKursList`
// (b438) -> `InputTreatyArrangement.Kurs = .HASIL1` (b579). Kuerinya
// (`RDBList/GetMasterKursList.xml` b85-b86):
//
//	select TOIDR from treatyexchange where Quarter='0'
//	   and to_date(CARI1,'YYYYMMDD') BETWEEN trunc(TO_TIMESTAMP_TZ(STARTDATE, ...))
//	                                    AND trunc(TO_TIMESTAMP_TZ(ENDDATE, ...))
//	   and IDCURRENCY = <USD>
//
// ⛔ Nama komponen tanpa kata "testing" (penyimpangan sadar 8); identitas USD
// dibaca dari master mata uang, tidak ditanam (penyimpangan sadar 7).
//
// ⛔ Lanjutan 6 (29-09-2026): tanggal master diurai ORACLE, seperti Pega - Go
// hanya menerima `DATE` dan baris yang Oracle tolak. Pengurai teks Go dibuang:
// ia menolak jam lima angka (`20190801T00000.000 GMT`) yang Oracle terima,
// sehingga satu baris mematikan kurs seluruh tahun.
//
// Pemakai kurs di Pega:
//   - `HitungRpUsd_depan` (onchange Rp tujuh form induk): `Usd = Rp / Kurs`
//     (`@divide(..., 8)` b383 untuk TreatyLimit) - `Usd` form hanya dibaca;
//   - `CalculateTSIExcludeTreaty` (exclusion Occupation): IDR -> `Usd = Rp / Kurs`
//     4 desimal (b248), USD -> `Rp = Usd * Kurs` (b394);
//   - 14 aktivitas `NewTreatyArr*`: kurs kosong -> "Tidak ada Nilai Kurs di
//     Tahun : " + TreatyYear (mis. `NewTreatyArrEpi.xml` b870) dan form tidak tampil.
//
// Dibaca sesudah: tco_klausul.go.

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// QuarterKursTahunanTCO - `Quarter='0'` VERBATIM (AC 48); artinya `[terbuka]`.
const QuarterKursTahunanTCO = "0"

// FormatTanggalKursTCO - topeng `TO_TIMESTAMP_TZ` VERBATIM
// `RDBList/GetMasterKursList.xml` b85-b86; dipakai Oracle, bukan Go.
const FormatTanggalKursTCO = `YYYYMMDD"T"HH24MISS.FF3 TZR`

// KodeMataUangAsalKursTCO - kode mata uang asal konversi (aturan bisnis USD ->
// IDR). IDENTITASNYA dibaca dari master mata uang (`CURRENCY.CURRENCY` ->
// `ID`, `Claim Life/RDBList/GetCurrencyID.xml`), bukan ditanam (AC 47).
const KodeMataUangAsalKursTCO = "USD"

// SkalaUsdDariRpTCO - `HitungRpUsd_depan.xml` b383 `@divide(Rp, Kurs, 8)`; skala 8
// untuk ketujuh form induk [keputusan work owner 29-09-2026] (OQ-TCO-18).
const SkalaUsdDariRpTCO = 8

// SkalaUsdExclusionTCO - `CalculateTSIExcludeTreaty.xml` b248 `@divide(Rp, Kurs, 4)`.
const SkalaUsdExclusionTCO = 4

// skalaRpDariUsdTCO - NUMBER(38,8) kolom `RP`.
const skalaRpDariUsdTCO = 8

var (
	// ErrKursTidakAda - tidak ada baris kurs berlaku pada tanggal itu (ADR-0015).
	ErrKursTidakAda = errors.New("models: no exchange rate in effect")
	// ErrKursTakTerurai - teks tanggal/nilai kurs master tidak dapat diurai.
	ErrKursTakTerurai = errors.New("models: exchange-rate master value cannot be parsed")
)

// GalatKursTidakAda - pesan `NewTreatyArrEpi.xml` b870 "Tidak ada Nilai Kurs di
// Tahun : " + TreatyYear, diterjemahkan [keputusan work owner 30-09-2026: bahasa Inggris].
type GalatKursTidakAda struct {
	TreatyYear string
	Tanggal    time.Time
}

func (g GalatKursTidakAda) Error() string {
	return "No exchange rate for Treaty Year : " + g.TreatyYear
}

// Is membuat `errors.Is(err, ErrKursTidakAda)` benar.
func (GalatKursTidakAda) Is(target error) bool { return target == ErrKursTidakAda }

// KursTCO adalah satu baris master kurs yang sudah diurai.
type KursTCO struct {
	// ToIDR - diurai hanya untuk baris yang BERLAKU (temuan /code-review):
	// seperti Pega, `TOIDR` baris lain yang rusak tidak menggagalkan pencarian.
	ToIDR *apd.Decimal
	// TeksToIDR - `TOIDR` apa adanya dari master (VARCHAR2).
	TeksToIDR string
	// Mulai, Akhir - `trunc(TO_TIMESTAMP_TZ(..))` hasil Oracle (AC 53).
	Mulai      time.Time
	Akhir      time.Time
	IDCurrency string
	Currency   string
	Quarter    string
	// ID - `ID` master (VARCHAR2, tidak unik di seluruh tabel): pemecah seri waktu input yang sama.
	ID string
	// DiinputPada - `DATEIU` diurai ORACLE (topeng `FormatTanggalKursTCO`) - HANYA untuk urutan antarbaris (driver
	// memberi jam dinding zona server, semua baris bergeser sama); nol = kosong / ditolak Oracle (paling lama).
	DiinputPada time.Time
	// BarisDitolak - cacah baris master yang tanggalnya ditolak Oracle pada
	// baca yang sama; diisi `PilihKursBerlakuTCO`, dilaporkan ke layar.
	BarisDitolak int
	// BarisKembar - cacah baris berlaku LAIN yang identik dengan yang dipakai
	// (master memuat baris kembar); diisi `PilihKursBerlakuTCO`, disebut di layar.
	BarisKembar int
	// BarisTertimpa - cacah baris berlaku LAIN yang TIDAK identik dan kalah oleh input yang lebih baru
	// (`PilihKursBerlakuTCO`); dilaporkan ke layar.
	BarisTertimpa int
}

// BarisKursDitolakTCO - satu baris master yang tanggalnya ditolak Oracle
// (`TO_TIMESTAMP_TZ` gagal, atau teksnya kosong): kolom pertama yang ditolak
// dan teksnya apa adanya.
type BarisKursDitolakTCO struct{ Kolom, Teks string }

// HasilMasterKursTCO - satu baca master kurs (satu mata uang, satu QUARTER):
// baris yang BERLAKU pada tanggal itu menurut Oracle, dan baris yang tanggalnya
// Oracle tolak - disaring per baris, tidak mematikan baris lain.
type HasilMasterKursTCO struct {
	Berlaku []KursTCO
	Ditolak []BarisKursDitolakTCO
}

// polaRibuanIndonesiaTCO - TOIDR berformat Indonesia yang ditulis menu Treaty Exchange Yearly: titik ribuan dengan
// koma desimal ("16.000,00"), atau beberapa kelompok ribuan tanpa desimal ("1.234.567").
var (
	polaRibuanKomaTCO = regexp.MustCompile(`^\d{1,3}(\.\d{3})+,\d+$`)
	polaRibuanSajaTCO = regexp.MustCompile(`^\d{1,3}(\.\d{3}){2,}$`)
)

// UraiNilaiKursTCO mengurai `TOIDR` (VARCHAR2) menjadi desimal persis > 0.
//
// Format Indonesia ("16.000,00", "1.234.567") diterima [keputusan work owner 11-10-2026, data master dari menu
// Treaty Exchange Yearly]. Satu titik tanpa koma ("16500.00") tetap titik DESIMAL; koma-lalu-titik ("15,500.25")
// dan bentuk campuran lain tetap ditolak - menebak salah pada kurs berarti salah pada setiap nilai Rupiah.
func UraiNilaiKursTCO(teks string) (*apd.Decimal, error) {
	t := strings.TrimSpace(teks)
	switch {
	case polaRibuanKomaTCO.MatchString(t):
		t = strings.ReplaceAll(strings.ReplaceAll(t, ".", ""), ",", ".")
	case polaRibuanSajaTCO.MatchString(t):
		t = strings.ReplaceAll(t, ".", "")
	case strings.Contains(t, ","):
		if strings.Contains(t, ".") {
			return nil, fmt.Errorf("%w: TOIDR %q contains both a comma and a dot", ErrKursTakTerurai, teks)
		}
		t = strings.ReplaceAll(t, ",", ".")
	}
	d, err := utils.ParseDecimal(t)
	if err != nil || t == "" {
		return nil, fmt.Errorf("%w: TOIDR %q", ErrKursTakTerurai, teks)
	}
	if d.Sign() <= 0 {
		return nil, fmt.Errorf("%w: TOIDR %q must be greater than zero", ErrKursTakTerurai, teks)
	}
	return d, nil
}

// PilihKursBerlakuTCO memilih SATU baris dari yang Oracle nyatakan berlaku
// (`to_date(CARI1,'YYYYMMDD') BETWEEN trunc(..STARTDATE..) AND trunc(..ENDDATE..)`).
//
// ⭐ INPUT TERBARU MENANG [keputusan work owner 11-10-2026, menggantikan "tolak, jangan menebak" 29-09-2026
// (OQ-TCO-18)]: master kurs dari menu Treaty Exchange Yearly membolehkan periode yang tumpang tindih, sehingga setiap
// input kurs baru membuat TCO menolak. Dari baris yang berlaku dipakai yang `DATEIU`-nya paling baru (diurai
// Oracle); waktu sama -> `ID` terbesar. Pega memutar seluruh hasil dan menyimpan yang TERAKHIR (`testingKurs`
// langkah 3) dalam urutan tak tentu - di sini urutannya tentu. Baris kembar identik (TOIDR, mulai, akhir sama)
// dicacah `BarisKembar`; baris lain yang kalah dicacah `BarisTertimpa` - keduanya dilaporkan ke layar.
//
// Tanpa baris berlaku tetapi ada baris yang tanggalnya ditolak Oracle: salah
// satunya mungkin baris yang dicari - master rusak, bukan "tidak ada kurs".
func PilihKursBerlakuTCO(h HasilMasterKursTCO, tanggal time.Time) (KursTCO, error) {
	tgl := utils.FormatTanggal(time.Date(tanggal.Year(), tanggal.Month(), tanggal.Day(), 0, 0, 0, 0, time.UTC))
	if len(h.Berlaku) == 0 {
		if len(h.Ditolak) > 0 {
			d := h.Ditolak[0]
			return KursTCO{}, fmt.Errorf("%w: no exchange rate in effect on %s, and %d master rows have dates rejected by Oracle "+
				"(format %s), e.g. %s %q", ErrKursTakTerurai, tgl, len(h.Ditolak), FormatTanggalKursTCO, d.Kolom, d.Teks)
		}
		return KursTCO{}, fmt.Errorf("%w on %s", ErrKursTidakAda, tgl)
	}
	pilih := 0
	for i := 1; i < len(h.Berlaku); i++ {
		if lebihBaruTCO(h.Berlaku[i], h.Berlaku[pilih]) {
			pilih = i
		}
	}
	k := h.Berlaku[pilih]
	for i, b := range h.Berlaku {
		switch {
		case i == pilih:
		case kembarIdentikTCO(k, b):
			k.BarisKembar++
		default:
			k.BarisTertimpa++
		}
	}
	k.BarisDitolak = len(h.Ditolak)
	return k, nil
}

// lebihBaruTCO - `a` diinput lebih baru dari `b`: `DATEIU` lebih akhir (nol = paling lama); waktu sama -> `ID`
// lebih besar (angka bila keduanya angka, selain itu teks).
func lebihBaruTCO(a, b KursTCO) bool {
	if !a.DiinputPada.Equal(b.DiinputPada) {
		return a.DiinputPada.After(b.DiinputPada)
	}
	x, ex := strconv.ParseInt(strings.TrimSpace(a.ID), 10, 64)
	y, ey := strconv.ParseInt(strings.TrimSpace(b.ID), 10, 64)
	if ex == nil && ey == nil {
		return x > y
	}
	return strings.TrimSpace(a.ID) > strings.TrimSpace(b.ID)
}

// kembarIdentikTCO - dua baris berlaku identik: teks TOIDR sama persis, hari
// mulai dan hari akhir (hasil `trunc(TO_TIMESTAMP_TZ(..))` Oracle) sama.
func kembarIdentikTCO(a, b KursTCO) bool {
	return strings.TrimSpace(a.TeksToIDR) == strings.TrimSpace(b.TeksToIDR) && a.Mulai.Equal(b.Mulai) && a.Akhir.Equal(b.Akhir)
}

func kuantisasiKurs(d *apd.Decimal, skala int32) (*apd.Decimal, error) {
	k := utils.DecimalContext()
	k.Rounding = apd.RoundHalfUp
	hasil := new(apd.Decimal)
	if _, err := k.Quantize(hasil, d, -skala); err != nil {
		return nil, err
	}
	return hasil, nil
}

// UsdDariRpTCO - `Usd = Rp / Kurs`, dibulatkan setengah-ke-atas pada skala.
// Rp kosong -> Usd kosong (medan wajib yang memeriksanya).
func UsdDariRpTCO(rp, kurs *apd.Decimal, skala int32) (*apd.Decimal, error) {
	if kurs == nil || kurs.Sign() <= 0 {
		return nil, ErrKursTidakAda
	}
	if rp == nil {
		return nil, nil
	}
	hasil := new(apd.Decimal)
	if _, err := utils.DecimalContext().Quo(hasil, rp, kurs); err != nil {
		return nil, fmt.Errorf("models: computing Usd from Rp: %w", err)
	}
	return kuantisasiKurs(hasil, skala)
}

// RpDariUsdTCO - `Rp = Usd * Kurs` (`CalculateTSIExcludeTreaty` b394), 8 desimal.
func RpDariUsdTCO(usd, kurs *apd.Decimal) (*apd.Decimal, error) {
	if kurs == nil || kurs.Sign() <= 0 {
		return nil, ErrKursTidakAda
	}
	if usd == nil {
		return nil, nil
	}
	hasil := new(apd.Decimal)
	if _, err := utils.DecimalContext().Mul(hasil, usd, kurs); err != nil {
		return nil, fmt.Errorf("models: computing Rp from Usd: %w", err)
	}
	return kuantisasiKurs(hasil, skalaRpDariUsdTCO)
}

// LengkapiKursTCO mengurai `TOIDR` baris yang terpilih bila belum diurai.
func LengkapiKursTCO(k KursTCO) (KursTCO, error) {
	if k.ToIDR != nil {
		return k, nil
	}
	d, err := UraiNilaiKursTCO(k.TeksToIDR)
	if err != nil {
		return KursTCO{}, fmt.Errorf("%w (row in effect %s - %s)", err, utils.FormatTanggal(k.Mulai), utils.FormatTanggal(k.Akhir))
	}
	k.ToIDR = d
	return k, nil
}
