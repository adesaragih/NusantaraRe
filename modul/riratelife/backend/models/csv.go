package models

// Upload CSV (container "Upload CSV" b2949 -> local action `UploadCSV_RIRATE` b3019; "View Upload" b3467 -> popup
// `ViewCSVResult_RIRate` b3486; "Simpan Upload" b4511 -> `SubmitRIRate_Act` b4535). Rule ketiganya TIDAK ada di XML
// (K3 keputusan work owner 05-10-2026): yang ada hanya label `Format excel : USEDBY, CONTRACT, GENDER, AGE, RATE`
// b5305. Seluruh aturan di bawah dirancang dari label itu dan dari isi `RATE_LIFE` DEV - ASUMSI, tabel "Asumsi
// terbuka" MODUL.md.
//
// Aturan pemisah yang tidak ambigu: pemisah dibaca dari baris kepala - ada `;` = pemisah `;` (RATE boleh berkoma
// desimal tanpa kutip), selain itu `,` (RATE berkoma desimal WAJIB dikutip `"0,5"`; tanpa kutip baris itu berkolom
// lebih dan ditolak, tidak ditebak).

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"nusantarare/inti/backend/utils"
)

// KolomCSV - kepala berkas, urutan label b5305. Urutan kolom di berkas bebas; nama wajib persis (tanpa beda huruf).
var KolomCSV = []string{"USEDBY", "CONTRACT", "GENDER", "AGE", "RATE"}

// Batas unggahan (ASUMSI). `M_RATE_LIFE` DEV: 98.305 baris untuk 348 ringkasan.
const (
	// MaksBarisCSV - baris data paling banyak sekali unggah.
	MaksBarisCSV = 10000
	// MaksBytesCSV - ukuran teks CSV paling besar.
	MaksBytesCSV = 4 << 20
	// MaksRingkasanCSV - R/I RATE NAME berbeda paling banyak sekali unggah (satu daftar IN Oracle <= 1000).
	MaksRingkasanCSV = 100
	// UmurMaks - AGE dan CONTRACT 0-120 (`RATE_LIFE` DEV).
	UmurMaks = 120
	// batasRate - panjang teks RATE (DEV: 1-20).
	batasRate = 30
)

// Gender sah (`RATE_LIFE` DEV: U 98.004, M 201, F 99).
var Gender = []string{"U", "M", "F"}

// ErrCSV - berkas tidak dapat dipakai seluruhnya (kepala, kosong, terlalu besar); pesannya untuk pengguna.
var ErrCSV = errors.New("models: berkas CSV ditolak")

// BarisCSV - satu baris data yang sudah dirapikan. Baris = nomor baris di berkas (kepala = 1).
type BarisCSV struct {
	Baris    int    `json:"baris"`
	UsedBy   string `json:"usedby"`
	Contract string `json:"contract"`
	Gender   string `json:"gender"`
	Age      string `json:"age"`
	Rate     string `json:"rate"`
}

// GalatBaris - galat satu baris berkas.
type GalatBaris struct {
	Baris int    `json:"baris"`
	Pesan string `json:"pesan"`
}

// Kalimat - `Row n: pesan.`
func (g GalatBaris) Kalimat() string { return fmt.Sprintf("Row %d: %s.", g.Baris, g.Pesan) }

func tolakCSV(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrCSV, fmt.Sprintf(format, a...))
}

var (
	polaBulat = regexp.MustCompile(`^[0-9]+$`)
	polaRate  = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
)

// NormalRate - RATE desimal tak bertanda; pemisah desimal koma ATAU titik (satu, tanpa pemisah ribuan). Disimpan
// berkoma desimal seperti mayoritas data DEV (`0.50` -> `0,50`, `007,5` -> `7,5`). Diurai `utils.ParseDecimal`.
func NormalRate(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", errors.New("RATE is required")
	}
	t := strings.Replace(s, ",", ".", 1)
	if strings.Contains(t, ",") || strings.Count(t, ".") > 1 || !polaRate.MatchString(t) || len(t) > batasRate {
		return "", fmt.Errorf("RATE must be a number with at most one decimal comma or point (got %q)", s)
	}
	d, err := utils.ParseDecimal(t)
	if err != nil {
		return "", fmt.Errorf("RATE must be a number (got %q)", s)
	}
	return strings.Replace(utils.FormatDecimal(d), ".", ",", 1), nil
}

// normalBulat - bilangan bulat 0..UmurMaks tanpa nol depan; kosong = "" bila boleh.
func normalBulat(nama, s string, bolehKosong bool) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		if bolehKosong {
			return "", nil
		}
		return "", fmt.Errorf("%s is required", nama)
	}
	n, err := strconv.Atoi(s)
	if !polaBulat.MatchString(s) || err != nil || n > UmurMaks {
		return "", fmt.Errorf("%s must be a whole number from 0 to %d (got %q)", nama, UmurMaks, s)
	}
	return strconv.Itoa(n), nil
}

// angkaKunci - teks angka warisan (`05`, ` 5`) setara `5`; teks lain apa adanya (dipangkas).
func angkaKunci(s string) string {
	s = strings.TrimSpace(s)
	if polaBulat.MatchString(s) {
		if n, err := strconv.Atoi(s); err == nil {
			return strconv.Itoa(n)
		}
	}
	return s
}

// KunciRate - kunci kembar (USEDBY, GENDER, AGE, CONTRACT): nama tanpa beda huruf dan spasi tepi, angka tanpa nol
// depan; CONTRACT kosong (warisan NULL) berbeda dari `0`.
func KunciRate(usedBy, gender, age, contract string) string {
	return strings.ToUpper(strings.TrimSpace(usedBy)) + "\x00" + strings.ToUpper(strings.TrimSpace(gender)) + "\x00" +
		angkaKunci(age) + "\x00" + angkaKunci(contract)
}

// KunciNama - pencocokan R/I RATE NAME (upload mencocokkan ringkasan lewat nama): tanpa beda huruf dan spasi tepi.
func KunciNama(s string) string { return strings.ToUpper(strings.TrimSpace(s)) }

// pemisahDari - `;` bila baris kepala memuatnya, selain itu `,`.
func pemisahDari(teks string) rune {
	kepala := strings.TrimLeft(teks, "\r\n")
	if i := strings.IndexAny(kepala, "\r\n"); i >= 0 {
		kepala = teks[:i]
	}
	if strings.Contains(kepala, ";") {
		return ';'
	}
	return ','
}

func barisKosong(rek []string) bool {
	for _, s := range rek {
		if strings.TrimSpace(s) != "" {
			return false
		}
	}
	return true
}

// UraiCSV mengurai teks berkas: baris data sah (dirapikan) dan galat per baris. err (ErrCSV) = berkas ditolak
// seluruhnya. Baris kosong dilewati. Baris kepala wajib.
func UraiCSV(teks string) ([]BarisCSV, []GalatBaris, error) {
	if len(teks) > MaksBytesCSV {
		return nil, nil, tolakCSV("the CSV file is larger than %d MB", MaksBytesCSV>>20)
	}
	teks = strings.TrimPrefix(teks, string(rune(0xFEFF)))
	if strings.TrimSpace(teks) == "" {
		return nil, nil, tolakCSV("the CSV file is empty")
	}
	pemisah := pemisahDari(teks)
	r := csv.NewReader(strings.NewReader(teks))
	r.Comma = pemisah
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = true
	bacaGalat := func(err error) error {
		var pe *csv.ParseError
		if errors.As(err, &pe) {
			return tolakCSV("the CSV file cannot be read at line %d: %v", pe.Line, pe.Err)
		}
		return tolakCSV("the CSV file cannot be read")
	}
	kepala, err := r.Read()
	if err != nil {
		return nil, nil, bacaGalat(err)
	}
	letak, ok := letakKolom(kepala)
	if !ok {
		return nil, nil, tolakCSV("the header must be USEDBY, CONTRACT, GENDER, AGE, RATE (separated by ; or ,), got %q",
			strings.Join(kepala, string(pemisah)))
	}
	var (
		hasil   []BarisCSV
		galat   []GalatBaris
		pertama = map[string]int{}
		jumlah  int
	)
	for {
		rek, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, nil, bacaGalat(err)
		}
		if barisKosong(rek) {
			continue
		}
		jumlah++
		if jumlah > MaksBarisCSV {
			return nil, nil, tolakCSV("the CSV file has more than %d data rows; at most %d rows can be uploaded at once",
				MaksBarisCSV, MaksBarisCSV)
		}
		nomor, _ := r.FieldPos(0)
		if len(rek) != len(KolomCSV) {
			p := fmt.Sprintf("the row has %d columns, expected %d", len(rek), len(KolomCSV))
			if pemisah == ',' {
				p += " (with , as the separator a RATE with a decimal comma must be quoted, e.g. \"0,5\"; or use ; as the separator)"
			}
			galat = append(galat, GalatBaris{Baris: nomor, Pesan: p})
			continue
		}
		b, pesan := barisDari(nomor, rek, letak)
		if len(pesan) > 0 {
			galat = append(galat, GalatBaris{Baris: nomor, Pesan: strings.Join(pesan, "; ")})
			continue
		}
		k := KunciRate(b.UsedBy, b.Gender, b.Age, b.Contract)
		if n, ada := pertama[k]; ada {
			galat = append(galat, GalatBaris{Baris: nomor, Pesan: fmt.Sprintf(
				"USEDBY %s, GENDER %s, AGE %s, CONTRACT %s duplicates row %d", b.UsedBy, b.Gender, b.Age, tampilKosong(b.Contract), n)})
			continue
		}
		pertama[k] = nomor
		hasil = append(hasil, b)
	}
	if jumlah == 0 {
		return nil, nil, tolakCSV("the CSV file has no data rows")
	}
	return hasil, galat, nil
}

func tampilKosong(s string) string {
	if s == "" {
		return "(empty)"
	}
	return s
}

// letakKolom - indeks setiap KolomCSV di kepala; false bila kolomnya tidak tepat lima nama itu.
func letakKolom(kepala []string) (map[string]int, bool) {
	if len(kepala) != len(KolomCSV) {
		return nil, false
	}
	letak := map[string]int{}
	for i, k := range kepala {
		letak[strings.ToUpper(strings.TrimSpace(k))] = i
	}
	for _, k := range KolomCSV {
		if _, ada := letak[k]; !ada {
			return nil, false
		}
	}
	return letak, len(letak) == len(KolomCSV)
}

func barisDari(nomor int, rek []string, letak map[string]int) (BarisCSV, []string) {
	sel := func(k string) string { return strings.TrimSpace(rek[letak[k]]) }
	var pesan []string
	b := BarisCSV{Baris: nomor, UsedBy: sel("USEDBY"), Gender: strings.ToUpper(sel("GENDER"))}
	switch {
	case b.UsedBy == "":
		pesan = append(pesan, "USEDBY is required")
	case len(b.UsedBy) > BatasNama:
		pesan = append(pesan, fmt.Sprintf("USEDBY is longer than %d characters", BatasNama))
	}
	if !sahGender(b.Gender) {
		pesan = append(pesan, fmt.Sprintf("GENDER must be U, M, or F (got %q)", sel("GENDER")))
	}
	var err error
	if b.Age, err = normalBulat("AGE", sel("AGE"), false); err != nil {
		pesan = append(pesan, err.Error())
	}
	if b.Contract, err = normalBulat("CONTRACT", sel("CONTRACT"), true); err != nil {
		pesan = append(pesan, err.Error())
	}
	if b.Rate, err = NormalRate(sel("RATE")); err != nil {
		pesan = append(pesan, err.Error())
	}
	return b, pesan
}

func sahGender(g string) bool {
	for _, s := range Gender {
		if g == s {
			return true
		}
	}
	return false
}
