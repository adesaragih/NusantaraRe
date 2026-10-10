package models

// Form R/I RISK DETAIL - section Pega `InboxRIRisk` (kelas `ASM-FW-GISFW-Int-RI_RISK_LIFE`, judul b382): tambah dan
// ubah satu baris `RIRISK_LIFE` milik ringkasan yang sedang dilihat. USEDBY = `TempIDUsedBy.USEDBY` b1727 (disabled,
// wajib) dan IDUSEDBY = `TempIDUsedBy.ID` b1524 (tersembunyi `1=2` b1643): keduanya dari ringkasan, bukan isian.
//
// BEDA dengan R/I Comm Life (XML): CONTRACT pxTextInput (Comm: pxNumber); YEAR TIDAK wajib (Comm: wajib); medan MONTH
// (pyMax 4, tidak wajib) tidak ada di Comm; RISK pxNumber wajib (Comm COMM pxTextInput).
//
// Tipe kolom (STRUKTUR-TABEL-RIRISKLIFE.md): CONTRACT, YEAR, MONTH VARCHAR2(10) warisan (teks angka); RISK NUMBER
// warisan TANPA skala. Nilai yang tidak muat DITOLAK berkalimat - tidak pernah dipotong atau dibulatkan (ADR-U-0003).

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// IsianRincian - isian form R/I RISK DETAIL (Save b3132 -> `AddToList_Act` b3155; EDIT b9784 -> `EditRIRiskLife_Act`
// b9808 dengan ID, Contract, Year, risk, UsedTo, Month).
type IsianRincian struct {
	// Contract - `.CONTRACT` b1920 (pxTextInput b1923, wajib b1883/b1933, placeholder 0 b1936, ChangeDotToPoint_DT b1963).
	Contract string `json:"contract"`
	// Year - `.YEAR` b2202 (pxTextInput b2205, TIDAK wajib b2216, placeholder 0 b2219, pyMax 4 b2220).
	Year string `json:"year"`
	// Month - `.MONTH` b2389 (pxTextInput b2392, TIDAK wajib b2403, placeholder 0 b2406, pyMax 4 b2407).
	Month string `json:"month"`
	// Risk - `.RISK` b2579 (pxNumber b2582, wajib b2542/b2595), label "RISK (PERMIL)" b2548.
	Risk string `json:"risk"`
}

// Batas angka.
const (
	// DigitTeksAngka - CONTRACT: lebar kolom VARCHAR2(10) warisan.
	DigitTeksAngka = 10
	// DigitYearMonth - `pyMax` 4 YEAR b2220 / MONTH b2407.
	DigitYearMonth = 4
	// DigitRisk - RISK NUMBER tanpa presisi/skala = presisi Oracle 38 angka (bulat + pecahan).
	DigitRisk = 38
)

var (
	polaBulat   = regexp.MustCompile(`^[0-9]+$`)
	polaDesimal = regexp.MustCompile(`^([0-9]+)(\.([0-9]+))?$`)
)

// BulatKanonik - bilangan bulat tak bertanda tanpa nol depan (`007` -> `7`, `0` -> `0`); ok=false bila bukan bulat.
func BulatKanonik(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if !polaBulat.MatchString(s) {
		return "", false
	}
	s = strings.TrimLeft(s, "0")
	if s == "" {
		s = "0"
	}
	return s, true
}

// DesimalKanonik - desimal tak bertanda, pemisah desimal koma ATAU titik (satu, tanpa pemisah ribuan), menjadi bentuk
// kanonik bertitik tanpa nol depan / nol ekor (`921,9` -> `921.9`, `0,0` -> `0`) - sama dengan yang dibaca kembali dari
// Oracle. Ekspresi SETARA dengan konversi migrasi 940 (koma dan titik sama-sama pemisah desimal). Mengembalikan juga
// cacah angka bulat dan pecahan sesudah dikanonikkan.
func DesimalKanonik(s string) (kanonik string, bulat, pecahan int, ok bool) {
	t := strings.Replace(strings.TrimSpace(s), ",", ".", 1)
	m := polaDesimal.FindStringSubmatch(t)
	if m == nil {
		return "", 0, 0, false
	}
	b := strings.TrimLeft(m[1], "0")
	p := strings.TrimRight(m[3], "0")
	if b == "" {
		b = "0"
	}
	kanonik = b
	if p != "" {
		kanonik += "." + p
	}
	bulat = len(b)
	if b == "0" {
		bulat = 0
	}
	return kanonik, bulat, len(p), true
}

// bulatOpsional - kosong boleh (""), selain itu bulat tak bertanda paling banyak `digit` angka, kanonik.
func bulatOpsional(nama, s string, digit int) (string, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return "", nil
	}
	k, ok := BulatKanonik(t)
	if !ok || len(k) > digit {
		return "", fmt.Errorf("%s must be a whole number of at most %d digits (got %q)", nama, digit, t)
	}
	return k, nil
}

// NormalContract - CONTRACT wajib (b1933), bulat tak bertanda yang muat VARCHAR2(10). [penyimpangan sadar - menunggu
// WO]: pxTextInput tanpa batas di XML; data DEV seluruhnya angka (fakta WO 08-10-2026). `ChangeDotToPoint_DT` b1963
// (isinya tidak ada di XML) tidak ditiru - bilangan bulat tidak berpemisah.
func NormalContract(s string) (string, error) {
	if strings.TrimSpace(s) == "" {
		return "", errors.New("CONTRACT is required")
	}
	return bulatOpsional("CONTRACT", s, DigitTeksAngka)
}

// NormalYear - YEAR TIDAK wajib (b2216), bulat paling banyak 4 angka (pyMax 4 b2220). Bukan tahun kalender: tahun
// polis (premiumlistlife menghitung `tahun - tahun mulai`); data DEV YEAR NULL di 1.994 baris.
func NormalYear(s string) (string, error) { return bulatOpsional("YEAR", s, DigitYearMonth) }

// NormalMonth - MONTH TIDAK wajib (b2403), bulat paling banyak 4 angka (pyMax 4 b2407); data DEV 0-180.
func NormalMonth(s string) (string, error) { return bulatOpsional("MONTH", s, DigitYearMonth) }

// NormalRisk - RISK wajib (b2595), desimal TIDAK NEGATIF (tanda `-` ditolak) yang muat NUMBER (38 angka); koma ATAU
// titik desimal (data DEV: 11.063 baris berkoma, 89 bertitik). [penyimpangan sadar - menunggu WO]: pxNumber b2582
// tanpa presisi/rentang; "tidak negatif" seperti COMM ricommlife.
func NormalRisk(s string) (string, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return "", errors.New("RISK is required")
	}
	if strings.HasPrefix(t, "-") {
		return "", fmt.Errorf("RISK must not be negative (got %q)", t)
	}
	k, bulat, pecahan, ok := DesimalKanonik(t)
	if !ok {
		return "", fmt.Errorf("RISK must be a number with at most one decimal comma or point (got %q)", t)
	}
	if bulat+pecahan > DigitRisk {
		return "", fmt.Errorf("RISK must have at most %d digits (got %q)", DigitRisk, t)
	}
	return k, nil
}

// PeriksaIsianRincian merapikan isian R/I RISK DETAIL; galat = kalimat untuk pengguna (semua masalah sekaligus).
func PeriksaIsianRincian(isi IsianRincian) (IsianRincian, error) {
	var out IsianRincian
	var pesan []string
	var err error
	if out.Contract, err = NormalContract(isi.Contract); err != nil {
		pesan = append(pesan, err.Error())
	}
	if out.Year, err = NormalYear(isi.Year); err != nil {
		pesan = append(pesan, err.Error())
	}
	if out.Month, err = NormalMonth(isi.Month); err != nil {
		pesan = append(pesan, err.Error())
	}
	if out.Risk, err = NormalRisk(isi.Risk); err != nil {
		pesan = append(pesan, err.Error())
	}
	if len(pesan) > 0 {
		return IsianRincian{}, errors.New(strings.Join(pesan, "; "))
	}
	return out, nil
}

// KunciRincian - kunci kembar (CONTRACT, YEAR, MONTH) di satu ringkasan: angka kanonik (`05` setara `5`), kosong tetap
// kosong. ASUMSI MODUL.md A1 (Comm memakai CONTRACT + YEAR; Risk menambah MONTH karena baris DEV memakai YEAR ATAU MONTH).
func KunciRincian(contract, year, month string) string {
	kanon := func(s string) string {
		if k, ok := BulatKanonik(s); ok {
			return k
		}
		return strings.TrimSpace(s)
	}
	return kanon(contract) + "\x00" + kanon(year) + "\x00" + kanon(month)
}

// KunciNama - pencocokan R/I RISK NAME (upload mencocokkan ringkasan lewat nama): tanpa beda huruf dan spasi tepi.
func KunciNama(s string) string { return strings.ToUpper(strings.TrimSpace(s)) }
