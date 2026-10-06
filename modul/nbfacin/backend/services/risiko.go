package services

// Pencarian alamat risiko - popup Choose Risk Address (tiket 36).

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// UkuranHalamanRisk - `[terverifikasi]` `Section\ChooseRiskAddress_ResultList.xml` pyPageSize 10.
const UkuranHalamanRisk = 10

// batasSaringRisk - A119: tiap saringan paling banyak 255 karakter (pola A73).
const batasSaringRisk = 255

// halamanMaksRisk - pagar luapan offset int32 (bind Oracle). Hasil RD dibatasi 100 baris (=
// 10 halaman); halaman 11 sampai batas ini sah tetapi kosong.
const halamanMaksRisk = (1<<31-1)/UkuranHalamanRisk + 1

// ErrMasukanRisk - tanpa satu pun saringan, halaman tak sah, atau saringan terlalu panjang. 400.
var ErrMasukanRisk = errors.New("services: isian pencarian alamat risiko tidak sah")

// ErrRiskTanpaDatabase - tabel RISKADDRESS tidak terbaca. 503.
var ErrRiskTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, tabel RISKADDRESS tidak terbaca")

// HasilRisk - satu halaman alamat risiko.
type HasilRisk struct {
	Baris          []models.RiskAddress
	Total, Halaman int
	Ukuran         int
}

// DenganRisk memasang pembaca alamat risiko (tiket 36).
func (s *Service) DenganRisk(r repository.PembacaRisk) *Service {
	s.risk = r
	return s
}

// CariRisk - halaman ke-`halaman` alamat risiko yang lolos saringan terisi. Minimal satu
// saringan - A121: semua kosong = 400 (Pega SearchRiskAddressAct mencari dengan kunci palsu
// RoadName "z" / ZipCode "123456" yang memberi hasil kosong). Tanpa identitas (A120, pola
// lookup tiket 27/28/33).
func (s *Service) CariRisk(ctx context.Context, saring models.SaringRisk, halaman int) (HasilRisk, error) {
	isian := []struct{ nama, nilai string }{{"address", saring.Address}, {"zipCode", saring.ZipCode},
		{"country", saring.Country}, {"province", saring.Province}, {"city", saring.City}, {"district", saring.District},
		{"territory", saring.Territory}}
	var masalah []string
	terisi := 0
	for _, i := range isian {
		if strings.TrimSpace(i.nilai) != "" {
			terisi++
		}
		if utf8.RuneCountInString(i.nilai) > batasSaringRisk {
			masalah = append(masalah, fmt.Sprintf("%s paling banyak %d karakter", i.nama, batasSaringRisk))
		}
	}
	if terisi == 0 {
		masalah = append(masalah, "minimal satu saringan wajib diisi")
	}
	if halaman < 1 || halaman > halamanMaksRisk {
		masalah = append(masalah, fmt.Sprintf("halaman harus bilangan bulat 1..%d", halamanMaksRisk))
	}
	if len(masalah) > 0 {
		return HasilRisk{}, fmt.Errorf("%w: %s", ErrMasukanRisk, strings.Join(masalah, "; "))
	}
	if s.risk == nil {
		return HasilRisk{}, ErrRiskTanpaDatabase
	}
	baris, total, err := s.risk.CariRisk(ctx, saring, (halaman-1)*UkuranHalamanRisk, UkuranHalamanRisk)
	if err != nil {
		return HasilRisk{}, err
	}
	return HasilRisk{Baris: baris, Total: total, Halaman: halaman, Ukuran: UkuranHalamanRisk}, nil
}
