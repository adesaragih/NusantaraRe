package login

// Token sesi - isi cookie `rnm_sesi`, ditandatangani HMAC-SHA256.
//
// Tanpa tabel sesi (keputusan work owner 01-10-2026). Isinya akun, versi sesi
// (`M_LOGIN_GO.SESSION_VERSION` saat diterbitkan), waktu login, dan waktu
// habis. Server mencocokkan versinya dengan kolom itu tiap permintaan, jadi
// menaikkan kolomnya mencabut cookie lama di mana pun ia berada.
//
// Bentuk: base64url("v1|akun|versi|mulai|habis") + "." + base64url(HMAC).

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ErrSesiTidakSah - cookie rusak, palsu, kedaluwarsa, atau sudah dicabut.
var ErrSesiTidakSah = errors.New("login: sesi tidak sah atau sudah berakhir")

// Token adalah isi cookie sesi.
type Token struct {
	Akun  string
	Versi int64
	// Mulai - waktu login; batas mutlak sesi dihitung darinya (BatasSesi).
	Mulai time.Time
	// Habis - batas diam; diperpanjang tiap sesi dipakai.
	Habis time.Time
}

// TokenBaru menerbitkan token saat login berhasil.
func TokenBaru(akun string, versi int64, sekarang time.Time) Token {
	t := Token{Akun: akun, Versi: versi, Mulai: sekarang.Truncate(time.Second)}
	return t.Perpanjang(sekarang)
}

// Perpanjang menggeser batas diam - tidak pernah melewati BatasSesi sejak login.
func (t Token) Perpanjang(sekarang time.Time) Token {
	habis := sekarang.Add(BatasDiam).Truncate(time.Second)
	if mutlak := t.Mulai.Add(BatasSesi); habis.After(mutlak) {
		habis = mutlak
	}
	t.Habis = habis
	return t
}

func tanda(rahasia []byte, isi string) string {
	m := hmac.New(sha256.New, rahasia)
	m.Write([]byte(isi))
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}

// Tandatangani merakit nilai cookie.
func Tandatangani(rahasia []byte, t Token) string {
	isi := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("v1|%s|%d|%d|%d",
		t.Akun, t.Versi, t.Mulai.Unix(), t.Habis.Unix())))
	return isi + "." + tanda(rahasia, isi)
}

// BacaToken memeriksa tanda tangan dan waktunya. Versi dan status akun
// diperiksa pemanggil terhadap M_LOGIN_GO.
func BacaToken(rahasia []byte, nilai string, sekarang time.Time) (Token, error) {
	isi, ttd, ok := strings.Cut(nilai, ".")
	if !ok || len(rahasia) == 0 || !hmac.Equal([]byte(ttd), []byte(tanda(rahasia, isi))) {
		return Token{}, ErrSesiTidakSah
	}
	mentah, err := base64.RawURLEncoding.DecodeString(isi)
	if err != nil {
		return Token{}, ErrSesiTidakSah
	}
	b := strings.Split(string(mentah), "|")
	if len(b) != 5 || b[0] != "v1" || !AkunSah(b[1]) {
		return Token{}, ErrSesiTidakSah
	}
	versi, e1 := strconv.ParseInt(b[2], 10, 64)
	mulai, e2 := strconv.ParseInt(b[3], 10, 64)
	habis, e3 := strconv.ParseInt(b[4], 10, 64)
	if e1 != nil || e2 != nil || e3 != nil {
		return Token{}, ErrSesiTidakSah
	}
	t := Token{Akun: b[1], Versi: versi, Mulai: time.Unix(mulai, 0), Habis: time.Unix(habis, 0)}
	if sekarang.Before(t.Mulai) || !sekarang.Before(t.Habis) || !sekarang.Before(t.Mulai.Add(BatasSesi)) {
		return Token{}, ErrSesiTidakSah
	}
	return t, nil
}
