package login

// Data kepegawaian akun - email, nomor HP, NIK, dan jabatan (Kelola User, permintaan work owner 03-10-2026:
// "tambahkan email, no hp, nik dan jabatan; buat dalam bahasa inggris"). Kolom `M_LOGIN_GO` migrasi 904 berbahasa
// Inggris (`EMAIL`, `PHONE_NUMBER`, `EMPLOYEE_ID`, `JOB_POSITION`), begitu pula label layar dan pesan galatnya.
//
// ⛔ SEMUANYA OPSIONAL: akun yang sudah ada tidak punya isinya, dan kosong = NULL. Yang diisi harus berformat sah.
// NIK di sini = Nomor Induk Karyawan (employee ID), bukan NIK kependudukan.

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"

	"nusantarare/inti/backend/db"
)

// Kontak adalah data kepegawaian satu akun. Disisipkan ke isian dan ringkasan akun; JSON-nya rata.
type Kontak struct {
	Email   string `json:"email"`
	Telepon string `json:"telepon"`
	NIK     string `json:"nik"`
	Jabatan string `json:"jabatan"`
}

// Batas panjang = lebar kolom migrasi 904.
const (
	PanjangMaksEmail   = 254
	PanjangMaksTelepon = 30
	PanjangMaksNIK     = 30
	PanjangMaksJabatan = 150
)

var (
	// ErrEmailTidakSah - bukan bentuk alamat email.
	ErrEmailTidakSah = errors.New("login: Email is not valid; use a format like name@company.com (max. 254 characters)")
	// ErrTeleponTidakSah - bukan nomor HP 8-15 digit.
	ErrTeleponTidakSah = errors.New("login: Phone Number must contain 8 to 15 digits, may start with +, and may use spaces or hyphens")
	// ErrNIKTidakSah - NIK berisi tanda selain huruf, angka, titik, garis miring, atau tanda hubung.
	ErrNIKTidakSah = errors.New("login: Employee ID (NIK) may only contain letters, digits, dots, slashes, or hyphens (max. 30 characters)")
	// ErrJabatanTidakSah - jabatan terlalu panjang.
	ErrJabatanTidakSah = errors.New("login: Position is too long (max. 150 characters)")
)

var (
	polaEmail   = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	polaTelepon = regexp.MustCompile(`^\+?[0-9][0-9 -]*[0-9]$`)
	polaNIK     = regexp.MustCompile(`^[A-Za-z0-9./-]+$`)
)

// RapikanKontak membuang spasi tepi setiap isian.
func RapikanKontak(k Kontak) Kontak {
	return Kontak{
		Email:   strings.TrimSpace(k.Email),
		Telepon: strings.TrimSpace(k.Telepon),
		NIK:     strings.TrimSpace(k.NIK),
		Jabatan: strings.TrimSpace(k.Jabatan),
	}
}

// PeriksaKontak menolak isian yang terisi tetapi tidak berformat sah - galat PERTAMA, urut medan di form.
// Kosong selalu sah. Panjang dihitung BYTE seperti kolom VARCHAR2.
func PeriksaKontak(k Kontak) error {
	if k.Email != "" && (len(k.Email) > PanjangMaksEmail || !polaEmail.MatchString(k.Email)) {
		return ErrEmailTidakSah
	}
	if k.Telepon != "" {
		digit := 0
		for _, c := range k.Telepon {
			if c >= '0' && c <= '9' {
				digit++
			}
		}
		if len(k.Telepon) > PanjangMaksTelepon || !polaTelepon.MatchString(k.Telepon) || digit < 8 || digit > 15 {
			return ErrTeleponTidakSah
		}
	}
	if k.NIK != "" && (len(k.NIK) > PanjangMaksNIK || !polaNIK.MatchString(k.NIK)) {
		return ErrNIKTidakSah
	}
	if len(k.Jabatan) > PanjangMaksJabatan || !utf8.ValidString(k.Jabatan) {
		return ErrJabatanTidakSah
	}
	return nil
}

// nilaiKontak - empat nilai bind berurutan EMAIL, PHONE_NUMBER, EMPLOYEE_ID, JOB_POSITION; kosong = NULL.
func nilaiKontak(k Kontak) []any {
	return []any{db.KosongJadiNil(k.Email), db.KosongJadiNil(k.Telepon), db.KosongJadiNil(k.NIK), db.KosongJadiNil(k.Jabatan)}
}
