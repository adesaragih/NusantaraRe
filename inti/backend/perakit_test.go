package backend_test

// Uji perakit modul - struktur tim satu folder per modul, R4 (30-09-2026).
//
// Yang dijaga: setiap `backend/modul.go` menyatakan kontrak yang DISEDIAKAN
// dan yang DIBUTUHKAN, dan perakit menyambungnya menurut jenis antarmuka -
// tanpa satu baris tangan per modul di berkas bersama. Modul di sini tiruan;
// sambungan modul sungguhan diuji di `inti/backend/daftar`.

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/config"
)

// pembacaUji dan penulisUji - dua kontrak tiruan (antarmuka, seperti
// `inti/backend/kontrak`).
type pembacaUji interface{ Baca() string }
type penulisUji interface{ Tulis(string) }

type pembacaTetap string

func (p pembacaTetap) Baca() string { return string(p) }

// modulUji - `inti.Modul` tiruan yang mengingat kontrak yang diterimanya.
type modulUji struct {
	nama    string
	pembaca pembacaUji
}

func (m modulUji) Nama() string                               { return m.nama }
func (modulUji) DaftarkanRute(*http.ServeMux)                 {}
func (modulUji) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }

func rakitUji(t *testing.T, daftar ...inti.Pendaftaran) (inti.Rakitan, error) {
	t.Helper()
	return inti.Rakit(inti.NewDasar(nil), config.Config{}, func(string) {}, daftar)
}

func TestPemakaiMenerimaKontrakYangDisediakanPenyedia(t *testing.T) {
	disediakan := pembacaTetap("polis dari penyedia")
	penyedia := inti.Pendaftaran{
		Nama:        "penyedia",
		Menyediakan: []inti.Kontrak{inti.KontrakDari[pembacaUji]()},
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			inti.Sediakan[pembacaUji](p, disediakan)
			return modulUji{nama: "penyedia"}, nil
		},
	}
	pemakai := inti.Pendaftaran{
		Nama:        "pemakai",
		Membutuhkan: []inti.Kontrak{inti.KontrakDari[pembacaUji]()},
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			return modulUji{nama: "pemakai", pembaca: inti.Ambil[pembacaUji](p)}, nil
		},
	}
	// Pemakai didaftar LEBIH DULU: urutan bangun ditentukan kontrak, bukan
	// urutan daftar.
	r, err := rakitUji(t, pemakai, penyedia)
	if err != nil {
		t.Fatal(err)
	}
	var nama []string
	for _, m := range r.Modul {
		nama = append(nama, m.Nama())
	}
	// Urutan hasil = urutan nama modul, apa pun urutan bangunnya.
	if mau := []string{"pemakai", "penyedia"}; !reflect.DeepEqual(nama, mau) {
		t.Fatalf("modul %v, mau %v", nama, mau)
	}
	if got := r.Modul[0].(modulUji).pembaca; got != disediakan {
		t.Errorf("pemakai menerima %v, mau nilai yang disediakan penyedia %v", got, disediakan)
	}
	mau := []inti.Sambungan{{Kontrak: "backend_test.pembacaUji", Penyedia: "penyedia", Pemakai: "pemakai"}}
	if !reflect.DeepEqual(r.Sambungan, mau) {
		t.Errorf("sambungan %+v, mau %+v", r.Sambungan, mau)
	}
}

// ⛔ Kontrak dibutuhkan tanpa penyedia: galat saat menyala yang menyebut
// kontrak DAN modulnya - bukan nil diam-diam.
func TestKontrakTanpaPenyediaGagalMenyebutKontrakDanModul(t *testing.T) {
	pemakai := inti.Pendaftaran{
		Nama:        "pemakai",
		Membutuhkan: []inti.Kontrak{inti.KontrakDari[pembacaUji]()},
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			t.Error("modul dibangun walau kontraknya tanpa penyedia")
			return modulUji{nama: "pemakai"}, nil
		},
	}
	_, err := rakitUji(t, pemakai)
	if err == nil {
		t.Fatal("perakit tidak menolak kontrak tanpa penyedia")
	}
	for _, mau := range []string{"backend_test.pembacaUji", "pemakai"} {
		if !strings.Contains(err.Error(), mau) {
			t.Errorf("galat %q tidak menyebut %q", err, mau)
		}
	}
}

// ⛔ Kesalahan pernyataan di `backend/modul.go` menggagalkan perakitan dengan
// kalimat yang menyebut modul dan kontraknya - pemakai tidak pernah menerima
// nil diam-diam.
func TestPernyataanKontrakYangTidakDitepatiGagal(t *testing.T) {
	pembaca := inti.KontrakDari[pembacaUji]()
	kasus := []struct {
		nama   string
		daftar []inti.Pendaftaran
		mau    []string
	}{
		{
			nama: "menyatakan Menyediakan tetapi tidak menyerahkannya",
			daftar: []inti.Pendaftaran{
				{Nama: "penyedia", Menyediakan: []inti.Kontrak{pembaca},
					Bangun: func(*inti.Perakitan) (inti.Modul, error) { return modulUji{nama: "penyedia"}, nil }},
				{Nama: "pemakai", Membutuhkan: []inti.Kontrak{pembaca},
					Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
						return modulUji{nama: "pemakai", pembaca: inti.Ambil[pembacaUji](p)}, nil
					}},
			},
			mau: []string{"penyedia", "backend_test.pembacaUji"},
		},
		{
			nama: "mengambil kontrak yang tidak ada di Membutuhkan",
			daftar: []inti.Pendaftaran{
				{Nama: "penyedia", Menyediakan: []inti.Kontrak{pembaca},
					Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
						inti.Sediakan[pembacaUji](p, pembacaTetap("x"))
						return modulUji{nama: "penyedia"}, nil
					}},
				{Nama: "pemakai",
					Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
						return modulUji{nama: "pemakai", pembaca: inti.Ambil[pembacaUji](p)}, nil
					}},
			},
			mau: []string{"pemakai", "backend_test.pembacaUji", "Membutuhkan"},
		},
		{
			nama: "menyerahkan kontrak yang tidak ada di Menyediakan",
			daftar: []inti.Pendaftaran{
				{Nama: "penyedia",
					Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
						inti.Sediakan[pembacaUji](p, pembacaTetap("x"))
						return modulUji{nama: "penyedia"}, nil
					}},
			},
			mau: []string{"penyedia", "backend_test.pembacaUji", "Menyediakan"},
		},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			_, err := rakitUji(t, k.daftar...)
			if err == nil {
				t.Fatal("perakitan tidak gagal")
			}
			for _, m := range k.mau {
				if !strings.Contains(err.Error(), m) {
					t.Errorf("galat %q tidak menyebut %q", err, m)
				}
			}
		})
	}
}

// ⛔ Daftar yang bentuknya salah ditolak SEBELUM satu modul pun dibangun (nama
// modul hasil `Bangun` yang berbeda dari pendaftarannya: sesudahnya), dengan
// kalimat yang menyebut modulnya - termasuk ketergantungan melingkar, yang
// tanpa penjagaan membuat proses berputar tanpa henti saat menyala.
func TestDaftarBerbentukSalahDitolak(t *testing.T) {
	pembaca := inti.KontrakDari[pembacaUji]()
	penulis := inti.KontrakDari[penulisUji]()
	bangun := func(nama string) func(*inti.Perakitan) (inti.Modul, error) {
		return func(*inti.Perakitan) (inti.Modul, error) { return modulUji{nama: nama}, nil }
	}
	kasus := []struct {
		nama   string
		daftar []inti.Pendaftaran
		mau    []string
	}{
		{"nama ganda", []inti.Pendaftaran{
			{Nama: "kembar", Bangun: bangun("kembar")}, {Nama: "kembar", Bangun: bangun("kembar")},
		}, []string{"kembar", "dua kali"}},
		{"dua penyedia satu kontrak", []inti.Pendaftaran{
			{Nama: "satu", Menyediakan: []inti.Kontrak{pembaca}, Bangun: bangun("satu")},
			{Nama: "dua", Menyediakan: []inti.Kontrak{pembaca}, Bangun: bangun("dua")},
		}, []string{"backend_test.pembacaUji", "satu", "dua"}},
		{"ketergantungan melingkar", []inti.Pendaftaran{
			{Nama: "ayam", Menyediakan: []inti.Kontrak{pembaca}, Membutuhkan: []inti.Kontrak{penulis}, Bangun: bangun("ayam")},
			{Nama: "telur", Menyediakan: []inti.Kontrak{penulis}, Membutuhkan: []inti.Kontrak{pembaca}, Bangun: bangun("telur")},
		}, []string{"melingkar", "ayam", "telur"}},
		{"kontrak bukan antarmuka", []inti.Pendaftaran{
			{Nama: "salah", Menyediakan: []inti.Kontrak{inti.KontrakDari[pembacaTetap]()}, Bangun: bangun("salah")},
		}, []string{"salah", "antarmuka"}},
		{"tanpa Bangun", []inti.Pendaftaran{{Nama: "hampa"}}, []string{"hampa", "Bangun"}},
		{"Nama modul berbeda dari pendaftarannya", []inti.Pendaftaran{
			{Nama: "didaftar", Bangun: bangun("lain")},
		}, []string{"didaftar", "lain"}},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			_, err := rakitUji(t, k.daftar...)
			if err == nil {
				t.Fatal("daftar berbentuk salah diterima")
			}
			for _, m := range k.mau {
				if !strings.Contains(err.Error(), m) {
					t.Errorf("galat %q tidak menyebut %q", err, m)
				}
			}
		})
	}
}
