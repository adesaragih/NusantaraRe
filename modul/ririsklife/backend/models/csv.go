package models

// Upload CSV ("Upload CSV" b3156 -> local action; "View Upload" b3676 -> harness `ViewCSVResult_RIRisk` b3695;
// "Simpan Upload" b4690 -> `SubmitRIRisk_Act` b4714 - `InboxSummaryRIRisk`). Rule ketiganya TIDAK ada di XML; yang ada
// hanya label `Format excel : USEDBY, CONTRACT, YEAR, MONTH, RISK` b5495 (di Pega labelnya tersembunyi `1=2` b5263 -
// urutan kolomnya tetap patokan). Aturan di bawah = gaya ricommlife / riratelife - ASUMSI, MODUL.md bab Asumsi.
//
// Pemisah dibaca dari baris kepala: ada `;` = pemisah `;` (RISK boleh berkoma desimal tanpa kutip), selain itu `,` (RISK
// berkoma desimal WAJIB dikutip `"0,5"`; tanpa kutip baris itu berkolom lebih dan ditolak, tidak ditebak).

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
)

// KolomCSV - kepala berkas, urutan label b5495. Urutan kolom di berkas bebas; nama wajib persis (tanpa beda huruf).
var KolomCSV = []string{"USEDBY", "CONTRACT", "YEAR", "MONTH", "RISK"}

// Batas unggahan (ASUMSI, sama dengan riratelife).
const (
	// MaksBarisCSV - baris data paling banyak sekali unggah.
	MaksBarisCSV = 10000
	// MaksBytesCSV - ukuran teks CSV paling besar.
	MaksBytesCSV = 4 << 20
	// MaksRingkasanCSV - R/I RISK NAME berbeda paling banyak sekali unggah (satu daftar IN Oracle <= 1000).
	MaksRingkasanCSV = 100
)

// ErrCSV - berkas tidak dapat dipakai seluruhnya (kepala, kosong, terlalu besar); pesannya untuk pengguna.
var ErrCSV = errors.New("models: berkas CSV ditolak")

// BarisCSV - satu baris data yang sudah dirapikan. Baris = nomor baris di berkas (kepala = 1).
type BarisCSV struct {
	Baris    int    `json:"baris"`
	UsedBy   string `json:"usedby"`
	Contract string `json:"contract"`
	Year     string `json:"year"`
	Month    string `json:"month"`
	Risk     string `json:"risk"`
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

// pemisahDari - `;` bila baris kepala memuatnya, selain itu `,`.
func pemisahDari(teks string) rune {
	kepala := strings.TrimLeft(teks, "\r\n")
	if i := strings.IndexAny(kepala, "\r\n"); i >= 0 {
		kepala = kepala[:i]
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
		return nil, nil, tolakCSV("the header must be USEDBY, CONTRACT, YEAR, MONTH, RISK (separated by ; or ,), got %q",
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
				p += " (with , as the separator a RISK with a decimal comma must be quoted, e.g. \"0,5\"; or use ; as the separator)"
			}
			galat = append(galat, GalatBaris{Baris: nomor, Pesan: p})
			continue
		}
		b, pesan := barisDari(nomor, rek, letak)
		if len(pesan) > 0 {
			galat = append(galat, GalatBaris{Baris: nomor, Pesan: strings.Join(pesan, "; ")})
			continue
		}
		k := KunciNama(b.UsedBy) + "\x01" + KunciRincian(b.Contract, b.Year, b.Month)
		if n, ada := pertama[k]; ada {
			galat = append(galat, GalatBaris{Baris: nomor, Pesan: fmt.Sprintf(
				"USEDBY %s, CONTRACT %s, YEAR %s, MONTH %s duplicates row %d", b.UsedBy, b.Contract, kosongTampil(b.Year),
				kosongTampil(b.Month), n)})
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

// letakKolom - indeks setiap KolomCSV di kepala; false bila kolomnya tidak tepat empat nama itu.
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
	b := BarisCSV{Baris: nomor, UsedBy: sel("USEDBY")}
	switch {
	case b.UsedBy == "":
		pesan = append(pesan, "USEDBY is required")
	case len(b.UsedBy) > BatasNama:
		pesan = append(pesan, fmt.Sprintf("USEDBY is longer than %d characters", BatasNama))
	}
	var err error
	if b.Contract, err = NormalContract(sel("CONTRACT")); err != nil {
		pesan = append(pesan, err.Error())
	}
	if b.Year, err = NormalYear(sel("YEAR")); err != nil {
		pesan = append(pesan, err.Error())
	}
	if b.Month, err = NormalMonth(sel("MONTH")); err != nil {
		pesan = append(pesan, err.Error())
	}
	if b.Risk, err = NormalRisk(sel("RISK")); err != nil {
		pesan = append(pesan, err.Error())
	}
	return b, pesan
}

// kosongTampil - nilai kosong ditulis "(empty)" di kalimat galat.
func kosongTampil(s string) string {
	if s == "" {
		return "(empty)"
	}
	return s
}
