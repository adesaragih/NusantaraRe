package services_test

// Uji tiket 32 — syarat berbeda tiap pemulihan limit.
//
// ⛔ Data uji di sini WAJIB memuat nilai selain 100, dan tiket 32 menyebut
// alasannya sendiri: uji yang seluruhnya bernilai 100 menghasilkan hasil yang
// sama pada rancangan lama maupun baru - ia tidak memisahkan apa pun.

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// persen membangun desimal dari teks.
//
// Ia PANIC bila teksnya bukan desimal, dan itu disengaja: teks yang salah di
// sini adalah uji yang salah tulis, bukan perilaku yang sedang diuji. Tanpa
// `*testing.T` supaya dapat dipakai juga saat tabel kasus dibangun.
func persen(s string) *apd.Decimal {
	d, _, err := apd.NewFromString(s)
	if err != nil {
		panic("data uji " + s + " bukan desimal: " + err.Error())
	}
	return d
}

func TestPemulihanMenolakTanpaIdentitasSebelumGudang(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	err := l.CatatPemulihanLimit(context.Background(), inti.Pelaku{}, 1,
		[]models.PemulihanLimit{{NomorUrut: 1, PersenPemulihan: persen("100")}})
	if !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Fatalf("mau ErrTanpaIdentitas, dapat %v", err)
	}
	// Penolakan yang tetap menulis ke basis data bukan penolakan.
	if len(g.pemulihan) != 0 {
		t.Errorf("gudang ditulis %d kali walau identitas tidak ada", len(g.pemulihan))
	}
}

func TestPemulihanMenolakMasukanTidakSah(t *testing.T) {
	kasus := []struct {
		nama   string
		baris  []models.PemulihanLimit
		petik  string
		sasar  error
		idLayr int64
	}{
		{
			nama:   "persen pemulihan 101 - di luar rentang INV-41",
			baris:  []models.PemulihanLimit{{NomorUrut: 1, PersenPemulihan: persen("101")}},
			petik:  "101",
			sasar:  services.ErrPemulihanTidakSah,
			idLayr: 7,
		},
		{
			nama:   "persen pemulihan negatif",
			baris:  []models.PemulihanLimit{{NomorUrut: 1, PersenPemulihan: persen("-0.5")}},
			petik:  "-0.5",
			sasar:  services.ErrPemulihanTidakSah,
			idLayr: 7,
		},
		{
			nama:   "persen tambahan di luar rentang",
			baris:  []models.PemulihanLimit{{NomorUrut: 1, PersenPemulihan: persen("100"), PersenTambahan: persen("150")}},
			petik:  "150",
			sasar:  services.ErrPemulihanTidakSah,
			idLayr: 7,
		},
		{
			nama: "nomor urut kembar pada satu layer",
			baris: []models.PemulihanLimit{
				{NomorUrut: 2, PersenPemulihan: persen("100")},
				{NomorUrut: 2, PersenPemulihan: persen("50")},
			},
			petik:  "sudah dipakai",
			sasar:  services.ErrPemulihanTidakSah,
			idLayr: 7,
		},
		{
			nama:   "nomor urut nol",
			baris:  []models.PemulihanLimit{{NomorUrut: 0, PersenPemulihan: persen("100")}},
			petik:  "nomor urut",
			sasar:  services.ErrPemulihanTidakSah,
			idLayr: 7,
		},
		{
			nama:   "persen pemulihan kosong - kolomnya wajib isi",
			baris:  []models.PemulihanLimit{{NomorUrut: 1}},
			petik:  "kosong",
			sasar:  services.ErrPemulihanTidakSah,
			idLayr: 7,
		},
		{
			nama:   "daftar kosong bukan jalur menghapus",
			baris:  []models.PemulihanLimit{},
			petik:  "nol baris",
			sasar:  services.ErrPemulihanTidakSah,
			idLayr: 7,
		},
		{
			nama:   "pengenal layer tidak sah",
			baris:  []models.PemulihanLimit{{NomorUrut: 1, PersenPemulihan: persen("100")}},
			petik:  "pengenal layer",
			sasar:  services.ErrMasukanTidakSah,
			idLayr: 0,
		},
	}

	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			g := &gudangTiruan{}
			l := services.LayananDengan(g)

			err := l.CatatPemulihanLimit(context.Background(), pelakuAda, k.idLayr, k.baris)
			if !errors.Is(err, k.sasar) {
				t.Fatalf("mau %v, dapat %v", k.sasar, err)
			}
			// Galat yang tidak menyebut APA yang ditolak membuat yang
			// memperbaikinya menebak.
			if !strings.Contains(err.Error(), k.petik) {
				t.Errorf("pesan tidak menyebut %q: %v", k.petik, err)
			}
			if len(g.pemulihan) != 0 {
				t.Errorf("masukan yang ditolak tetap sampai ke gudang")
			}
		})
	}
}

// Uji POSITIF, dan ia pokok tiket 32: dua pemulihan dengan tarif BERBEDA
// diterima, dan keduanya sampai ke gudang apa adanya.
//
// Nilainya sengaja 0 dan 100 - ketentuan pasar "pemulihan pertama gratis,
// kedua berbayar penuh" yang sistem lama tidak dapat nyatakan.
func TestPemulihanBertarifBerbedaDiterima(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	baris := []models.PemulihanLimit{
		{NomorUrut: 1, PersenPemulihan: persen("100"), PersenTambahan: persen("0")},
		{NomorUrut: 2, PersenPemulihan: persen("100"), PersenTambahan: persen("100")},
	}
	if err := l.CatatPemulihanLimit(context.Background(), pelakuAda, 42, baris); err != nil {
		t.Fatalf("mau diterima, dapat %v", err)
	}
	if len(g.pemulihan) != 1 || len(g.pemulihan[0]) != 2 {
		t.Fatalf("gudang menerima %v", g.pemulihan)
	}
	if g.layerDicatat[0] != 42 {
		t.Errorf("layer yang dicatat %d, mau 42", g.layerDicatat[0])
	}
	// Keduanya terbaca kembali BERBEDA - bila nilainya diseragamkan di
	// perjalanan, inilah yang menangkapnya.
	pertama := g.pemulihan[0][0].PersenTambahan.Text('f')
	kedua := g.pemulihan[0][1].PersenTambahan.Text('f')
	if pertama == kedua {
		t.Errorf("kedua tarif sampai sama (%s): nilainya diseragamkan di perjalanan", pertama)
	}
	if pertama != "0" || kedua != "100" {
		t.Errorf("tarif yang sampai: %q dan %q, mau 0 dan 100", pertama, kedua)
	}
}

// Positif kedua: tarif yang BELUM dinyatakan tetap kosong, tidak menjadi nol.
//
// Kosong dan nol adalah dua pernyataan yang berbeda - "belum ditentukan"
// versus "gratis" - dan sistem lama tidak dapat membedakannya sebab keduanya
// selalu 100.
func TestTarifPemulihanKosongTidakMenjadiNol(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	baris := []models.PemulihanLimit{{NomorUrut: 1, PersenPemulihan: persen("75.5")}}
	if err := l.CatatPemulihanLimit(context.Background(), pelakuAda, 9, baris); err != nil {
		t.Fatalf("mau diterima, dapat %v", err)
	}
	if g.pemulihan[0][0].PersenTambahan != nil {
		t.Errorf("tarif kosong berubah menjadi %v", g.pemulihan[0][0].PersenTambahan)
	}
	if got := g.pemulihan[0][0].PersenPemulihan.Text('f'); got != "75.5" {
		t.Errorf("persen pemulihan sampai sebagai %q, mau 75.5", got)
	}
}
