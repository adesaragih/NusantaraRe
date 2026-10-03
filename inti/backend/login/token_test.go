package login

import (
	"errors"
	"strings"
	"testing"
	"time"
)

var rahasiaUji = []byte("rahasia-uji-yang-panjangnya-32-byte!")

func TestTokenPulangPergi(t *testing.T) {
	t0 := time.Date(2026, 10, 1, 8, 0, 0, 0, time.UTC)
	tok := TokenBaru("UJI-ADMIN", 7, t0)
	if !tok.Habis.Equal(t0.Add(BatasDiam)) || !tok.Mulai.Equal(t0) {
		t.Fatalf("token baru: %+v", tok)
	}
	s := Tandatangani(rahasiaUji, tok)
	got, err := BacaToken(rahasiaUji, s, t0.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if got.Akun != "UJI-ADMIN" || got.Versi != 7 || !got.Mulai.Equal(t0) || !got.Habis.Equal(tok.Habis) {
		t.Errorf("token terbaca %+v", got)
	}
}

func TestTokenDitolak(t *testing.T) {
	t0 := time.Date(2026, 10, 1, 8, 0, 0, 0, time.UTC)
	s := Tandatangani(rahasiaUji, TokenBaru("UJI-ADMIN", 1, t0))
	isi, tanda, _ := strings.Cut(s, ".")
	palsu := Tandatangani(rahasiaUji, TokenBaru("UJI-LAIN", 1, t0))
	isiPalsu, _, _ := strings.Cut(palsu, ".")
	for nama, k := range map[string]struct {
		nilai   string
		rahasia []byte
		saat    time.Time
	}{
		"rahasia lain":             {s, []byte("rahasia-lain-yang-panjangnya-32-byte"), t0},
		"isi ditukar":              {isiPalsu + "." + tanda, rahasiaUji, t0},
		"tanda dibuang":            {isi, rahasiaUji, t0},
		"kosong":                   {"", rahasiaUji, t0},
		"sampah":                   {"a.b", rahasiaUji, t0},
		"lewat 30 menit diam":      {s, rahasiaUji, t0.Add(BatasDiam)},
		"rahasia kosong":           {s, nil, t0},
		"jam mundur sebelum lahir": {s, rahasiaUji, t0.Add(-time.Hour)},
	} {
		if _, err := BacaToken(k.rahasia, k.nilai, k.saat); !errors.Is(err, ErrSesiTidakSah) {
			t.Errorf("%s: %v, mau ErrSesiTidakSah", nama, err)
		}
	}
}

// Diperpanjang tiap dipakai (30 menit diam), tetapi TIDAK PERNAH melewati
// 10 jam sejak login - keputusan work owner 01-10-2026.
func TestTokenDiperpanjangSampaiBatasSesi(t *testing.T) {
	t0 := time.Date(2026, 10, 1, 8, 0, 0, 0, time.UTC)
	tok := TokenBaru("UJI-ADMIN", 1, t0)
	saat := t0
	for i := 0; i < 40; i++ { // 40 x 20 menit = 13 jam 20 menit
		saat = saat.Add(20 * time.Minute)
		got, err := BacaToken(rahasiaUji, Tandatangani(rahasiaUji, tok), saat)
		if saat.Before(t0.Add(BatasSesi)) {
			if err != nil {
				t.Fatalf("%v sesudah login: %v", saat.Sub(t0), err)
			}
			tok = got.Perpanjang(saat)
			if tok.Habis.After(t0.Add(BatasSesi)) {
				t.Fatalf("diperpanjang melewati 10 jam: %v", tok.Habis.Sub(t0))
			}
			continue
		}
		if !errors.Is(err, ErrSesiTidakSah) {
			t.Fatalf("%v sesudah login masih diterima", saat.Sub(t0))
		}
		return
	}
	t.Fatal("sesi tidak pernah berakhir")
}

func TestTokenMenolakAkunTidakSah(t *testing.T) {
	t0 := time.Now()
	s := Tandatangani(rahasiaUji, Token{Akun: "a|b", Versi: 1, Mulai: t0, Habis: t0.Add(time.Minute)})
	if _, err := BacaToken(rahasiaUji, s, t0); !errors.Is(err, ErrSesiTidakSah) {
		t.Errorf("akun berpemisah diterima: %v", err)
	}
}
