package models

// Form R/I COMM DETAIL - section Pega `InboxRIComm` (kelas `ASM-FW-GISFW-Int-RI_COMM_LIFE`, judul b349): tambah dan
// ubah satu baris `M_RICOMM_LIFE` milik ringkasan yang sedang dilihat. USEDBY = `TempIDUsedBy.USEDBY` b1718 (disabled,
// wajib) dan IDUSEDBY = `TempIDUsedBy.ID` b1509 (tersembunyi `1=2` b1633): keduanya dari ringkasan, bukan isian.
//
// Tipe (STRUKTUR-TABEL-RICOMMLIFE.md): CONTRACT dan YEAR NUMBER(5) bilangan bulat, COMM NUMBER(38,8). Nilai yang tidak
// muat DITOLAK berkalimat - tidak pernah dipotong atau dibulatkan (ADR-U-0003, nol float).

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// IsianKomisi - isian form R/I COMM DETAIL (Save b2931 -> `AddToList_Act` b2955; Edit -> `EditList_DT` b9323).
type IsianKomisi struct {
	// Contract - `.CONTRACT` b1904 (pxNumber b1907, wajib b1917/b1923, placeholder 0 b1927, ChangeDotToPoint_DT b1953).
	Contract string `json:"contract"`
	// Year - `.YEAR` b2185 (pxTextInput b2188, wajib b2196/b2201, placeholder 0 b2205, pyMax 4 b2206).
	Year string `json:"year"`
	// Comm - `.COMM` b2392 (pxTextInput b2395, wajib b2406/b2412; grid pxNumber b9101).
	Comm string `json:"comm"`
}

// Batas angka (kolom flat).
const (
	// MaksContract - NUMBER(5).
	MaksContract = 99999
	// DigitYear - `pyMax` 4 b2206; YEAR TEPAT 4 angka (NormalYear).
	DigitYear = 4
	// DigitBulatComm / SkalaComm - NUMBER(38,8): 30 angka di depan koma, 8 di belakang.
	DigitBulatComm = 30
	SkalaComm      = 8
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
// kanonik bertitik tanpa nol depan / nol ekor (`007,50` -> `7.5`, `0,0` -> `0`) - sama dengan yang dibaca kembali dari
// Oracle. Mengembalikan juga cacah angka bulat dan pecahan sesudah dikanonikkan.
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

// NormalContract - CONTRACT wajib, bulat 0..MaksContract. [penyimpangan sadar - menunggu WO]: pxNumber b1907 tanpa
// presisi (`Precision` kosong b2089).
func NormalContract(s string) (string, error) {
	if strings.TrimSpace(s) == "" {
		return "", errors.New("CONTRACT is required")
	}
	k, ok := BulatKanonik(s)
	if !ok || len(k) > 5 {
		return "", fmt.Errorf("CONTRACT must be a whole number from 0 to %d (got %q)", MaksContract, strings.TrimSpace(s))
	}
	return k, nil
}

// NormalYear - YEAR wajib, TEPAT DigitYear angka tanpa nol depan (1000-9999). XML hanya memberi batas atas (pyMax 4
// b2206); "tepat 4" = [penyimpangan sadar - menunggu WO] (bawaan keputusan work owner 06-10-2026). Nol depan ditolak:
// kolom NUMBER(5) tidak dapat menyimpannya, jadi `0026` tidak akan pulang sebagai 4 angka.
func NormalYear(s string) (string, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return "", errors.New("YEAR is required")
	}
	if !polaBulat.MatchString(t) || len(t) != DigitYear || t[0] == '0' {
		return "", fmt.Errorf("YEAR must be exactly %d digits from 1000 to 9999 (got %q)", DigitYear, t)
	}
	return t, nil
}

// NormalComm - COMM wajib, desimal TIDAK NEGATIF (tanda `-` ditolak) yang muat NUMBER(38,8); koma atau titik desimal;
// tanpa batas 100. [penyimpangan sadar - menunggu WO]: XML tanpa presisi/rentang (`InboxRIComm` b2392).
func NormalComm(s string) (string, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return "", errors.New("COMM is required")
	}
	if strings.HasPrefix(t, "-") {
		return "", fmt.Errorf("COMM must not be negative (got %q)", t)
	}
	k, bulat, pecahan, ok := DesimalKanonik(t)
	if !ok {
		return "", fmt.Errorf("COMM must be a number with at most one decimal comma or point (got %q)", t)
	}
	if bulat > DigitBulatComm || pecahan > SkalaComm {
		return "", fmt.Errorf("COMM must have at most %d digits before and %d after the decimal separator (got %q)",
			DigitBulatComm, SkalaComm, t)
	}
	return k, nil
}

// PeriksaIsianKomisi merapikan isian R/I COMM DETAIL; galat = kalimat untuk pengguna (semua masalah sekaligus).
func PeriksaIsianKomisi(isi IsianKomisi) (IsianKomisi, error) {
	var out IsianKomisi
	var pesan []string
	var err error
	if out.Contract, err = NormalContract(isi.Contract); err != nil {
		pesan = append(pesan, err.Error())
	}
	if out.Year, err = NormalYear(isi.Year); err != nil {
		pesan = append(pesan, err.Error())
	}
	if out.Comm, err = NormalComm(isi.Comm); err != nil {
		pesan = append(pesan, err.Error())
	}
	if len(pesan) > 0 {
		return IsianKomisi{}, errors.New(strings.Join(pesan, "; "))
	}
	return out, nil
}

// KunciKomisi - kunci kembar (CONTRACT, YEAR) di satu ringkasan: angka kanonik (`05` setara `5`). ASUMSI MODUL.md.
func KunciKomisi(contract, year string) string {
	kanon := func(s string) string {
		if k, ok := BulatKanonik(s); ok {
			return k
		}
		return strings.TrimSpace(s)
	}
	return kanon(contract) + "\x00" + kanon(year)
}

// KunciNama - pencocokan R/I COMM NAME (upload mencocokkan ringkasan lewat nama): tanpa beda huruf dan spasi tepi.
func KunciNama(s string) string { return strings.ToUpper(strings.TrimSpace(s)) }
