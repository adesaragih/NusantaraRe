package services

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/bordereaux/backend/models"
)

// Batas Upload CSV.
const (
	MaksBarisCSV = 5000
	// MaksPesan - pesan validasi terbanyak yang dikembalikan (seperti Aggregate).
	MaksPesan = 100
)

// GalatCSV membawa seluruh alasan penolakan berkas (baris dan kolom).
type GalatCSV struct{ Pesan []string }

func (g GalatCSV) Error() string { return strings.Join(g.Pesan, "\n") }

// Is - GalatCSV adalah masukan tidak sah.
func (g GalatCSV) Is(t error) bool { return t == ErrMasukanTidakSah }

// pemisah - `;` atau `,` menurut baris header (MODUploadCSVResults memakai pemisah daftar locale peminta; templat
// Bordereaux semuanya `;`).
func pemisah(header string) rune {
	if strings.Count(header, ",") > strings.Count(header, ";") {
		return ','
	}
	return ';'
}

// BacaCSV membaca berkas: BOM dibuang, baris pertama header, baris yang seluruh selnya kosong dilewati
// (`MODUploadCSVResults` langkah 8). Menjawab header dan baris data beserta nomor baris berkasnya (header = 1).
func BacaCSV(teks string) (header []string, baris [][]string, nomor []int, err error) {
	teks = strings.TrimPrefix(teks, string(rune(0xFEFF)))
	if strings.TrimSpace(teks) == "" {
		return nil, nil, nil, GalatCSV{[]string{"The file is empty"}}
	}
	pertama, _, _ := strings.Cut(teks, "\n")
	r := csv.NewReader(strings.NewReader(teks))
	r.Comma = pemisah(pertama)
	r.LazyQuotes = true
	r.FieldsPerRecord = -1
	r.ReuseRecord = false
	for {
		rec, e := r.Read()
		if errors.Is(e, io.EOF) {
			break
		}
		if e != nil {
			return nil, nil, nil, GalatCSV{[]string{fmt.Sprintf("The file cannot be read as CSV: %v", e)}}
		}
		if header == nil {
			header = rec
			continue
		}
		kosong := true
		for _, c := range rec {
			if strings.TrimSpace(c) != "" {
				kosong = false
				break
			}
		}
		if kosong {
			continue
		}
		baris_, _ := r.FieldPos(0)
		baris = append(baris, rec)
		nomor = append(nomor, baris_)
		if len(baris) > MaksBarisCSV {
			return nil, nil, nil, GalatCSV{[]string{fmt.Sprintf("The file has more than %d data rows", MaksBarisCSV)}}
		}
	}
	if len(baris) == 0 {
		return nil, nil, nil, GalatCSV{[]string{"The file has no data rows"}}
	}
	return header, baris, nomor, nil
}

var (
	polaAngka   = regexp.MustCompile(`^-?\d+(\.\d+)?$`)
	polaTanggal = regexp.MustCompile(`^(\d{1,2})[/.\-](\d{1,2})[/.\-](\d{4})$`)
	polaKanonik = regexp.MustCompile(`^\d{2}-\d{2}-\d{4}$`)
)

// AngkaIndonesia - sel angka CSV menjadi teks desimal kanonik: `1.234.567,89` -> `1234567.89` (titik ribuan dibuang,
// koma = desimal - rumus `*N` activity Mapping*), `%` dan spasi dibuang, `-` atau kosong = "" (NULL). ok false =
// bukan angka.
func AngkaIndonesia(s string) (string, bool) {
	s = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s, "%", ""), " ", ""))
	if s == "" || s == "-" {
		return "", true
	}
	s = strings.ReplaceAll(strings.ReplaceAll(s, ".", ""), ",", ".")
	if !polaAngka.MatchString(s) {
		return "", false
	}
	return kanonikAngka(s), true
}

// kanonikAngka - nol depan dibuang, nol belakang pecahan dibuang (Reduce), tanpa notasi ilmiah.
func kanonikAngka(s string) string {
	d, _, err := apd.NewFromString(s)
	if err != nil {
		return s
	}
	var r apd.Decimal
	r.Reduce(d)
	return r.Text('f')
}

// TanggalCSV - sel tanggal CSV `dd/MM/yyyy` (juga `d/M/yyyy`, pemisah `/` `-` `.`) menjadi `DD-MM-YYYY`; kosong = "".
// ok false = bukan tanggal yang ada (mis. 31/02/2024). Pega memotong substring tanpa memeriksa - diperbaiki.
func TanggalCSV(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", true
	}
	m := polaTanggal.FindStringSubmatch(s)
	if m == nil {
		return "", false
	}
	teks := fmt.Sprintf("%02s-%02s-%s", m[1], m[2], m[3])
	if _, err := time.Parse("02-01-2006", teks); err != nil {
		return "", false
	}
	return teks, true
}

// digitBulat - banyak digit bagian bulat teks desimal kanonik.
func digitBulat(s string) int {
	s = strings.TrimPrefix(s, "-")
	bulat, _, _ := strings.Cut(s, ".")
	bulat = strings.TrimLeft(bulat, "0")
	return len(bulat)
}

// periksaNilai menilai SATU nilai KANONIK terhadap kolom tabel; "" = sah (NULL). Dipakai untuk hasil CSV dan untuk
// baris yang dikirim kembali saat Save.
func periksaNilai(k models.KolomCSV, v string) string {
	if v == "" {
		return ""
	}
	switch k.Jenis {
	case models.Angka:
		if !polaAngka.MatchString(v) {
			return "is not a number"
		}
		if k.Presisi > 0 && digitBulat(v) > k.Presisi-k.Skala {
			return fmt.Sprintf("is too large (max %d digits before the decimal point)", k.Presisi-k.Skala)
		}
	case models.Tanggal:
		if !polaKanonik.MatchString(v) {
			return "is not a date dd/MM/yyyy"
		}
		if _, err := time.Parse("02-01-2006", v); err != nil {
			return "is not a valid date"
		}
	default:
		if k.Panjang > 0 && utf8.RuneCountInString(v) > k.Panjang {
			return fmt.Sprintf("is longer than %d characters", k.Panjang)
		}
	}
	return ""
}

// pasanganTanggal - kolom akhir periode untuk kolom awal (POI_START/POI_END, START_POI/END_POI, PERIOD..._START/_END,
// START_EXTEND/END_EXTEND, ...).
func pasanganTanggal(k models.KombinasiBdx) map[string]string {
	ada := map[string]bool{}
	for _, c := range k.Kolom {
		if c.Jenis == models.Tanggal {
			ada[c.Kolom] = true
		}
	}
	hasil := map[string]string{}
	for awal := range ada {
		var akhir string
		switch {
		case strings.HasSuffix(awal, "_START"):
			akhir = strings.TrimSuffix(awal, "_START") + "_END"
		case strings.HasPrefix(awal, "START_"):
			akhir = "END_" + strings.TrimPrefix(awal, "START_")
		default:
			continue
		}
		if ada[akhir] {
			hasil[awal] = akhir
		}
	}
	return hasil
}

func waktuKanonik(s string) time.Time {
	t, _ := time.Parse("02-01-2006", s)
	return t
}

// pengumpul pesan - berhenti menambah sesudah MaksPesan, lalu satu baris "... and N more".
type pengumpul struct {
	pesan []string
	lebih int
}

func (p *pengumpul) tambah(s string) {
	if len(p.pesan) < MaksPesan {
		p.pesan = append(p.pesan, s)
		return
	}
	p.lebih++
}

func (p *pengumpul) galat() error {
	if len(p.pesan) == 0 {
		return nil
	}
	pesan := p.pesan
	if p.lebih > 0 {
		pesan = append(pesan, fmt.Sprintf("... and %d more", p.lebih))
	}
	return GalatCSV{pesan}
}

// PetakanCSV membaca isi berkas untuk satu kombinasi: jumlah kolom header harus sama dengan templat, setiap sel dibaca
// menurut tipe kolom tabelnya, dan seluruh kesalahan dikumpulkan sekaligus (Pega: nol validasi).
func PetakanCSV(k models.KombinasiBdx, teks string) ([]models.Baris, error) {
	header, rekaman, nomor, err := BacaCSV(teks)
	if err != nil {
		return nil, err
	}
	if len(header) != len(k.Kolom) {
		return nil, GalatCSV{[]string{fmt.Sprintf(
			"The file header has %d columns; the %s %s template needs %d columns. Check that the file matches the chosen Type and Business.",
			len(header), k.Type, k.Business, len(k.Kolom))}}
	}
	var p pengumpul
	pasangan := pasanganTanggal(k)
	hasil := make([]models.Baris, 0, len(rekaman))
	for i, rec := range rekaman {
		b := models.Baris{}
		for j, c := range k.Kolom {
			mentah := ""
			if j < len(rec) {
				mentah = rec[j]
			}
			var (
				v  string
				ok = true
			)
			switch c.Jenis {
			case models.Angka:
				v, ok = AngkaIndonesia(mentah)
			case models.Tanggal:
				v, ok = TanggalCSV(mentah)
			default:
				v = strings.TrimSpace(mentah)
			}
			masalah := ""
			if !ok {
				masalah = map[models.Jenis]string{models.Angka: "is not a number", models.Tanggal: "is not a date dd/MM/yyyy"}[c.Jenis]
			} else {
				masalah = periksaNilai(c, v)
			}
			if masalah != "" {
				p.tambah(fmt.Sprintf("Row %d, column %s: %q %s", nomor[i], c.Judul, strings.TrimSpace(mentah), masalah))
				continue
			}
			b[c.Kolom] = v
		}
		for awal, akhir := range pasangan {
			if b[awal] != "" && b[akhir] != "" && waktuKanonik(b[akhir]).Before(waktuKanonik(b[awal])) {
				p.tambah(fmt.Sprintf("Row %d: %s is before %s", nomor[i], judulKolom(k, akhir), judulKolom(k, awal)))
			}
		}
		hasil = append(hasil, b)
	}
	if err := p.galat(); err != nil {
		return nil, err
	}
	return hasil, nil
}

func judulKolom(k models.KombinasiBdx, kolom string) string {
	for _, c := range k.Kolom {
		if c.Kolom == kolom {
			return c.Judul
		}
	}
	return kolom
}

// PeriksaBaris menilai baris KANONIK yang dikirim kembali saat Save (kolom tak dikenal ditolak).
func PeriksaBaris(k models.KombinasiBdx, baris []models.Baris) error {
	kenal := map[string]models.KolomCSV{}
	for _, c := range k.Kolom {
		kenal[c.Kolom] = c
	}
	var p pengumpul
	for i, b := range baris {
		for kolom, v := range b {
			if kolom == models.KolomID {
				continue
			}
			c, ok := kenal[kolom]
			if !ok {
				p.tambah(fmt.Sprintf("Row %d has an unknown column %s", i+1, kolom))
				continue
			}
			if m := periksaNilai(c, v); m != "" {
				p.tambah(fmt.Sprintf("Row %d, column %s: %q %s", i+1, c.Judul, v, m))
			}
		}
	}
	return p.galat()
}
