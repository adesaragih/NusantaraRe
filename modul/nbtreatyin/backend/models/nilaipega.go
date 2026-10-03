package models

// Untuk apa berkas ini: PENGURAI NILAI PEGA - SATU tempat untuk membaca nilai
// dokumen JSON Pega menjadi teks halaman: nilai skalar (P29 sifat 1: seluruh
// nilai bertipe teks; angka JSON tidak pernah lewat float), tanggal Pega
// (`YYYYMMDD` dan cap waktu ber-GMT, P29 sifat 2; diagram F21), dan zona
// waktu bisnis Asia/Jakarta. Dipakai pemecah dokumen lama (`dokumenlama.go`,
// tiket 22), pengurai master XOL (`repository.MasterXOLDariJSON`, K8), dan
// aturan tanggal (`TanggalProduksiNomor`).

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"nusantarare/inti/backend/utils"
)

// TeksSkalarJSON - nilai skalar JSON (didekode `UseNumber`) sebagai teks
// properti Pega; ok false untuk objek/larik. Angka (`json.Number`) lewat
// desimal eksak ke bentuk bertitik tanpa notasi ilmiah - digit tidak
// dibulatkan (`592629512.880000276` tetap utuh, ADR-0003); literal yang tak
// terurai dibawa apa adanya.
func TeksSkalarJSON(v any) (string, bool) {
	switch x := v.(type) {
	case nil:
		return "", true
	case string:
		return x, true
	case bool:
		return strconv.FormatBool(x), true
	case json.Number:
		d, err := utils.ParseDecimal(x.String())
		if err != nil {
			return x.String(), true
		}
		return d.Text('f'), true
	}
	return "", false
}

// Galat pembacaan tanggal lama.
var (
	// ErrTanggalAmbigu - susunan hari/bulan tidak dapat ditentukan dari
	// nilainya (`05/06/2017`). `[keputusan work owner]` K15: TIDAK ditebak;
	// dokumennya masuk laporan galat dan jumlahnya dilaporkan.
	ErrTanggalAmbigu = errors.New("models: tanggal ambigu (susunan hari/bulan tidak dapat ditentukan) - tidak ditebak (K15)")
	// ErrFormatTanggal - bukan salah satu dari dua format dokumen lama
	// (ID-19), atau tanggalnya tidak ada di kalender.
	ErrFormatTanggal = errors.New("models: format tanggal di luar YYYYMMDD / cap waktu Pega ber-GMT")
)

var (
	// polaYYYYMMDD - tanggal Pega (`"20171130"`, P29 sifat 2).
	polaYYYYMMDD = regexp.MustCompile(`^\d{8}$`)
	// polaCapWaktuGMT - cap waktu Pega (`"20170930T170000.000 GMT"`, P29 sifat 2).
	polaCapWaktuGMT = regexp.MustCompile(`^(\d{8}T\d{6})(\.\d{1,3})? GMT$`)
	// polaGarisMiring - susunan `dd/MM/yyyy` atau `MM/dd/yyyy …` (P32:
	// `InputPolicyTreatyIn_preDT` memakai keduanya).
	polaGarisMiring = regexp.MustCompile(`^(\d{1,2})/(\d{1,2})/(\d{4})(\s.*)?$`)
)

// BacaTanggalLama mengubah nilai tanggal dokumen lama menjadi bentuk
// pertukaran repository (`utils.TanggalSaja` / `utils.TanggalWaktu`).
//
//   - `YYYYMMDD` (AC 21) -> `YYYY-MM-DD`; jam tidak dikarang, juga untuk
//     kolom tanggal-waktu.
//   - `YYYYMMDDTHHMMSS.mmm GMT` (AC 22) -> jam dinding Asia/Jakarta. Pega
//     menyimpan cap waktu dalam GMT dan rule-nya sendiri membaca harinya di
//     Asia/Jakarta (`GeneratePolicyNoTreaty_Act` langkah 5.3, `TanggalProduksiNomor`):
//     `20170930T170000.000 GMT` adalah 1 Oktober 2017 pukul 00.00 WIB.
//     Kolom bertanggal saja menerima tanggal kalender Jakarta itu.
//   - `05/06/2017` -> ErrTanggalAmbigu (K15), tidak ditebak.
//   - selain itu -> ErrFormatTanggal, juga garis miring yang tidak ambigu:
//     cara menentukan susunan per baris belum diputuskan (P32 butir 1).
func BacaTanggalLama(teks string, g Golongan) (string, error) {
	teks = strings.TrimSpace(teks)
	switch {
	case teks == "":
		return "", nil
	case polaYYYYMMDD.MatchString(teks):
		t, err := time.Parse("20060102", teks)
		if err != nil {
			return "", fmt.Errorf("%w: %q", ErrFormatTanggal, teks)
		}
		return t.Format("2006-01-02"), nil
	}
	if m := polaCapWaktuGMT.FindStringSubmatch(teks); m != nil {
		t, err := time.ParseInLocation("20060102T150405", m[1], time.UTC)
		if err != nil {
			return "", fmt.Errorf("%w: %q", ErrFormatTanggal, teks)
		}
		lokal := t.In(zonaJakarta())
		if g == GolTanggal {
			return lokal.Format("2006-01-02"), nil
		}
		return lokal.Format("2006-01-02 15:04:05"), nil
	}
	if m := polaGarisMiring.FindStringSubmatch(teks); m != nil {
		a, _ := strconv.Atoi(m[1])
		b, _ := strconv.Atoi(m[2])
		if a >= 1 && a <= 12 && b >= 1 && b <= 12 && a != b {
			return "", fmt.Errorf("%w: %q", ErrTanggalAmbigu, teks)
		}
	}
	return "", fmt.Errorf("%w: %q", ErrFormatTanggal, teks)
}

// TanggalMasterPega menormalkan tanggal dokumen master XOL ke "2006-01-02"
// lewat pembaca yang SAMA dengan dokumen lama (`BacaTanggalLama`, kolom
// bertanggal saja). Teks yang sudah berbentuk `2006-01-02...` dipotong ke
// tanggalnya. Bentuk lain dibawa apa adanya - penyimpanan menolaknya dengan
// pesan jalur medannya (tidak ditebak, K15).
func TanggalMasterPega(s string) string {
	if v, err := BacaTanggalLama(s, GolTanggal); err == nil {
		return v
	}
	t := strings.TrimSpace(s)
	if len(t) >= 10 {
		if d, err := time.Parse("2006-01-02", t[:10]); err == nil {
			return d.Format("2006-01-02")
		}
	}
	return s
}

// zonaJakarta - zona hari bisnis (Asia/Jakarta, WIB).
func zonaJakarta() *time.Location {
	if l, err := time.LoadLocation("Asia/Jakarta"); err == nil {
		return l
	}
	return time.FixedZone("WIB", 7*3600)
}
