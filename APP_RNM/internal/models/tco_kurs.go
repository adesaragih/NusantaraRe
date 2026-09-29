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
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/pkg/utils"
)

// QuarterKursTahunanTCO - `Quarter='0'` VERBATIM (AC 48); artinya `[terbuka]`.
const QuarterKursTahunanTCO = "0"

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
	ErrKursTidakAda = errors.New("models: tidak ada kurs berlaku")
	// ErrKursGanda - lebih dari satu baris kurs berlaku pada tanggal yang sama.
	ErrKursGanda = errors.New("models: lebih dari satu kurs berlaku pada tanggal yang sama")
	// ErrKursTakTerurai - teks tanggal/nilai kurs master tidak dapat diurai.
	ErrKursTakTerurai = errors.New("models: nilai master kurs tidak dapat diurai")
)

// GalatKursTidakAda - pesan VERBATIM `NewTreatyArrEpi.xml` b870.
type GalatKursTidakAda struct {
	TreatyYear string
	Tanggal    time.Time
}

func (g GalatKursTidakAda) Error() string { return "Tidak ada Nilai Kurs di Tahun : " + g.TreatyYear }

// Is membuat `errors.Is(err, ErrKursTidakAda)` benar.
func (GalatKursTidakAda) Is(target error) bool { return target == ErrKursTidakAda }

// KursTCO adalah satu baris master kurs yang sudah diurai.
type KursTCO struct {
	// ToIDR - diurai hanya untuk baris yang BERLAKU (temuan /code-review):
	// seperti Pega, `TOIDR` baris lain yang rusak tidak menggagalkan pencarian.
	ToIDR *apd.Decimal
	// TeksToIDR - `TOIDR` apa adanya dari master (VARCHAR2).
	TeksToIDR  string
	Mulai      time.Time
	Akhir      time.Time
	IDCurrency string
	Currency   string
	Quarter    string
}

// polaTanggalKursTCO - 'YYYYMMDD"T"HH24MISS.FF3 TZR'; TZR boleh kosong.
var polaTanggalKursTCO = regexp.MustCompile(`^(\d{8})T\d{6}\.\d{1,9}( \S+)?$`)

// UraiTanggalKursTCO mengurai tanggal teks master kurs SEKALI, di batas baca
// (AC 53), lalu memangkasnya ke tanggal seperti `trunc(TO_TIMESTAMP_TZ(..))`:
// tanggal yang tertulis, di zona yang tertulis.
func UraiTanggalKursTCO(kolom, teks string) (time.Time, error) {
	m := polaTanggalKursTCO.FindStringSubmatch(strings.TrimSpace(teks))
	if m == nil {
		return time.Time{}, fmt.Errorf("%w: %s %q bukan bentuk YYYYMMDD\"T\"HH24MISS.FF3 TZR", ErrKursTakTerurai, kolom, teks)
	}
	t, err := time.Parse("20060102", m[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s %q: %v", ErrKursTakTerurai, kolom, teks, err)
	}
	return t, nil
}

// UraiNilaiKursTCO mengurai `TOIDR` (VARCHAR2) menjadi desimal persis > 0.
func UraiNilaiKursTCO(teks string) (*apd.Decimal, error) {
	t := strings.TrimSpace(teks)
	if strings.Contains(t, ",") {
		if strings.Contains(t, ".") {
			return nil, fmt.Errorf("%w: TOIDR %q memuat koma dan titik sekaligus", ErrKursTakTerurai, teks)
		}
		t = strings.ReplaceAll(t, ",", ".")
	}
	d, err := utils.ParseDecimal(t)
	if err != nil || t == "" {
		return nil, fmt.Errorf("%w: TOIDR %q", ErrKursTakTerurai, teks)
	}
	if d.Sign() <= 0 {
		return nil, fmt.Errorf("%w: TOIDR %q harus lebih dari nol", ErrKursTakTerurai, teks)
	}
	return d, nil
}

// PilihKursBerlakuTCO memilih SATU baris dengan Mulai <= tanggal <= Akhir.
//
// ⚠️ Pega memutar seluruh hasil dan menyimpan yang TERAKHIR (`testingKurs`
// langkah 3) - urutan tak tentu. Dua baris berlaku dinyatakan sebagai galat
// (master rusak) [keputusan work owner 29-09-2026] (OQ-TCO-18, ditutup).
func PilihKursBerlakuTCO(baris []KursTCO, tanggal time.Time) (KursTCO, error) {
	tgl := time.Date(tanggal.Year(), tanggal.Month(), tanggal.Day(), 0, 0, 0, 0, time.UTC)
	var cocok []KursTCO
	for _, k := range baris {
		if !tgl.Before(k.Mulai) && !tgl.After(k.Akhir) {
			cocok = append(cocok, k)
		}
	}
	switch len(cocok) {
	case 0:
		return KursTCO{}, fmt.Errorf("%w pada %s", ErrKursTidakAda, utils.FormatTanggal(tgl))
	case 1:
		return cocok[0], nil
	}
	return KursTCO{}, fmt.Errorf("%w: %d baris pada %s", ErrKursGanda, len(cocok), utils.FormatTanggal(tgl))
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
		return nil, fmt.Errorf("models: menghitung Usd dari Rp: %w", err)
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
		return nil, fmt.Errorf("models: menghitung Rp dari Usd: %w", err)
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
		return KursTCO{}, fmt.Errorf("%w (baris berlaku %s - %s)", err, utils.FormatTanggal(k.Mulai), utils.FormatTanggal(k.Akhir))
	}
	k.ToIDR = d
	return k, nil
}
