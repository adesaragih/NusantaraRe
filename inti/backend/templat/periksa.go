package templat

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
)

// BatasUkuran - berkas templat terbesar yang diterima (templat Bordereaux terbesar sekitar 14 KB).
const BatasUkuran = 1 << 20

// Kepala membaca baris header CSV: BOM dibuang, judul berisi baris baru di dalam kutip tetap SATU kolom.
func Kepala(isi []byte, pemisah rune) ([]string, error) {
	isi = bytes.TrimPrefix(isi, []byte("\xef\xbb\xbf"))
	r := csv.NewReader(bytes.NewReader(isi))
	r.Comma = pemisah
	r.LazyQuotes = true
	r.FieldsPerRecord = -1
	baris, err := r.Read()
	if errors.Is(err, io.EOF) {
		return nil, errors.New("file is empty")
	}
	if err != nil {
		return nil, fmt.Errorf("header cannot be read: %v", err)
	}
	return baris, nil
}

// rapikan - judul kolom untuk dibandingkan dan ditampilkan: spasi dan baris baru diringkas.
func rapikan(s string) string { return strings.Join(strings.Fields(s), " ") }

// Perbedaan - satu judul kolom yang berbeda dari versi aktif.
type Perbedaan struct {
	Kolom int    `json:"kolom"`
	Lama  string `json:"lama"`
	Baru  string `json:"baru"`
}

// HasilPeriksa - jawaban pemeriksaan berkas sebelum disimpan.
type HasilPeriksa struct {
	// Galat - alasan berkas DITOLAK; kosong = boleh disimpan.
	Galat []string `json:"galat"`
	// JumlahKolom - kolom header berkas (0 untuk slot bukan CSV).
	JumlahKolom int `json:"jumlahKolom"`
	// Perbedaan - judul kolom yang berbeda dari versi aktif; peringatan, bukan penolakan.
	Perbedaan []Perbedaan `json:"perbedaan"`
}

// periksaBerkas menilai berkas unggahan untuk slot `s` terhadap isi versi aktif `aktif`.
func periksaBerkas(s Slot, nama string, isi, aktif []byte) HasilPeriksa {
	h := HasilPeriksa{Galat: []string{}, Perbedaan: []Perbedaan{}}
	if !strings.HasSuffix(strings.ToLower(nama), s.Ekstensi) {
		h.Galat = append(h.Galat, fmt.Sprintf("File must be %s", s.Ekstensi))
	}
	switch {
	case len(isi) == 0:
		h.Galat = append(h.Galat, "File is empty")
	case len(isi) > BatasUkuran:
		h.Galat = append(h.Galat, fmt.Sprintf("File is larger than %d KB", BatasUkuran>>10))
	}
	if s.Pemisah == 0 || len(h.Galat) > 0 {
		return h
	}
	kolom, err := Kepala(isi, s.Pemisah)
	if err != nil {
		h.Galat = append(h.Galat, "File "+err.Error())
		return h
	}
	h.JumlahKolom = len(kolom)
	if s.JumlahKolom > 0 && len(kolom) != s.JumlahKolom {
		h.Galat = append(h.Galat, fmt.Sprintf(
			"Header has %d columns, this template needs %d columns separated by %q", len(kolom), s.JumlahKolom, string(s.Pemisah)))
		return h
	}
	lama, err := Kepala(aktif, s.Pemisah)
	if err != nil {
		return h
	}
	for i := range kolom {
		var l string
		if i < len(lama) {
			l = rapikan(lama[i])
		}
		if b := rapikan(kolom[i]); b != l {
			h.Perbedaan = append(h.Perbedaan, Perbedaan{Kolom: i + 1, Lama: l, Baru: b})
		}
	}
	return h
}
