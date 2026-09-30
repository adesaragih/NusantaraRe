package services

// Gerbang masukan bersama kelima jalur simpan.
//
// ⛔ Wajib-isi yang ditegakkan HANYA yang hidup di Pega (langkah 3
// `Page-Set-Messages` ber-PRE=true di kelima activity simpan, R9), dengan
// pesan VERBATIM. Gerbang mati (`PRE=false`) tidak dibangun - termasuk
// gerbang tahun (K7, OQ-MCRL-01).
// ⛔ Identitas baris baru tidak pernah datang dari klien (ADR-0006, K3).

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

var (
	// ErrWajibIsi - medan wajib kosong (422); pesan VERBATIM korpus.
	ErrWajibIsi = errors.New("services: required value is empty")
	// ErrIDDariKlien - baris baru membawa ID (400).
	ErrIDDariKlien = errors.New("services: a new row must not carry an id; the server assigns it")
	// ErrMasukanTidakSah - nilai tidak dapat dibaca (tanggal, angka) (422).
	ErrMasukanTidakSah = errors.New("services: value is not valid")
)

// GalatWajibIsi membawa pesan VERBATIM dan medan yang kosong (untuk log).
type GalatWajibIsi struct {
	Pesan string
	Medan []string
}

func (g GalatWajibIsi) Error() string { return g.Pesan }

// Is membuat errors.Is(err, ErrWajibIsi) benar.
func (GalatWajibIsi) Is(target error) bool { return target == ErrWajibIsi }

// wajib mengumpulkan medan kosong (`== ""` Pega; spasi dianggap kosong).
type wajib struct {
	kosong []string
}

func (w *wajib) teks(label, v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		w.kosong = append(w.kosong, label)
	}
	return v
}

func (w *wajib) galat(pesan string) error {
	if len(w.kosong) == 0 {
		return nil
	}
	return GalatWajibIsi{Pesan: pesan, Medan: w.kosong}
}

// tanggal membaca `YYYY-MM-DD`; kosong = waktu nol (wajib-isi diperiksa terpisah).
func tanggal(label, v string) (time.Time, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(utils.TanggalSaja, v)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s %q is not a date (YYYY-MM-DD)", ErrMasukanTidakSah, label, v)
	}
	return t, nil
}

// desimal membaca angka; koma diterima sebagai titik desimal (Pega
// `@replaceAll(…, ",", ".")`, `SetErrorMessageReinsurer` langkah 1), koma DAN
// titik sekaligus ditolak. Kosong = nil.
func desimal(label, v string) (*apd.Decimal, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil, nil
	}
	if strings.Contains(v, ",") {
		if strings.Contains(v, ".") {
			return nil, fmt.Errorf("%w: %s %q contains both a comma and a dot", ErrMasukanTidakSah, label, v)
		}
		v = strings.ReplaceAll(v, ",", ".")
	}
	d, err := utils.ParseDecimal(v)
	if err != nil {
		return nil, fmt.Errorf("%w: %s %q is not a number", ErrMasukanTidakSah, label, v)
	}
	return d, nil
}
