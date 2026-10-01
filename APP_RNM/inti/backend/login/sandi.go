package login

// Sandi - aturan, hash bcrypt, dan sandi sementara.
//
// ⛔ Sandi asli tidak pernah disimpan, dicatat ke log, atau dikirim balik.

import (
	"crypto/rand"
	"errors"
	"math/big"
	"regexp"
	"sync"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrSandiTerlaluPendek - kurang dari PanjangMinSandi karakter.
	ErrSandiTerlaluPendek = errors.New("login: sandi minimal 10 karakter")
	// ErrSandiTerlaluPanjang - lebih dari 72 byte (batas bcrypt).
	ErrSandiTerlaluPanjang = errors.New("login: sandi maksimal 72 byte")
)

// polaAkun - huruf, angka, titik, garis bawah, tanda hubung, @; 1-64.
// Tanpa spasi dan tanpa `|` (pemisah isi cookie sesi).
var polaAkun = regexp.MustCompile(`^[A-Za-z0-9._@-]{1,64}$`)

// AkunSah menjawab apakah bentuk LOGIN_ID dapat dipakai.
func AkunSah(akun string) bool { return polaAkun.MatchString(akun) }

// PeriksaSandiBaru menegakkan aturan sandi.
func PeriksaSandiBaru(sandi string) error {
	if utf8.RuneCountInString(sandi) < PanjangMinSandi {
		return ErrSandiTerlaluPendek
	}
	if len(sandi) > PanjangMaksSandi {
		return ErrSandiTerlaluPanjang
	}
	return nil
}

// HashSandi menghasilkan hash bcrypt untuk PASSWORD_HASH.
func HashSandi(sandi string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(sandi), BiayaBcrypt)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

// CocokSandi membandingkan sandi dengan hash tersimpan (waktu tetap).
func CocokSandi(hash, sandi string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(sandi)) == nil
}

// hashTiruan dipakai saat akun tidak ada atau nonaktif, supaya waktu jawab
// sama dengan akun yang ada - waktu tidak boleh membocorkan akun mana yang ada.
var hashTiruan = sync.OnceValue(func() string {
	h, _ := HashSandi("tiruan-penyama-waktu-jawab")
	return h
})

// abjadSementara - tanpa huruf yang mudah tertukar saat dibacakan (0/O, 1/l/I).
const abjadSementara = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789"

// SandiSementara membuat sandi 16 karakter dari crypto/rand - untuk akun baru
// dan reset; akunnya WAJIB ganti sandi saat login pertama.
func SandiSementara() (string, error) {
	b := make([]byte, 16)
	batas := big.NewInt(int64(len(abjadSementara)))
	for i := range b {
		n, err := rand.Int(rand.Reader, batas)
		if err != nil {
			return "", err
		}
		b[i] = abjadSementara[n.Int64()]
	}
	return string(b), nil
}
