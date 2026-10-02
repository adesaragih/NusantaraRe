package services_test

// Uji seam services tiket 21, 35, 36, 43.

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/treatyin/backend/services"
)

func desimal(t *testing.T, s string) *apd.Decimal {
	t.Helper()
	d, _, err := apd.NewFromString(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func lines(n int64) *int64 { return &n }

// ── Tiket 21 ──────────────────────────────────────────────────────────────

// Setiap kegagalan menghasilkan KETERANGAN yang menyebut apa yang gagal —
// bukan nol, dan bukan hasil yang diam-diam memakai kurs satu.
func TestHitungGagalMenghasilkanKeterangan(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{})
	kasus := []struct {
		nama        string
		nilai, kurs *apd.Decimal
		kode, sebut string
	}{
		{"kurs tidak ditemukan", desimal(t, "100"), nil, "USD", "tidak ditemukan"},
		{"kurs nol", desimal(t, "100"), desimal(t, "0"), "USD", "nol"},
		{"nilai sumber kosong", nil, desimal(t, "15000"), "USD", "kosong"},
		{"kode mata uang kosong", desimal(t, "100"), desimal(t, "15000"), "  ", "kode mata uang kosong"},
	}
	for _, k := range kasus {
		hasil, err := l.KonversiMataUang(context.Background(), pelakuAda, k.nilai, k.kurs, k.kode)
		if !errors.Is(err, services.ErrHitungGagal) {
			t.Errorf("%s: mau ErrHitungGagal, dapat %v", k.nama, err)
			continue
		}
		// ⛔ Yang paling penting: NIL, bukan nol. Nol terbaca sebagai hasil
		// yang sah di laporan, dan tidak ada yang tahu ia karangan.
		if hasil != nil {
			t.Errorf("%s: mengembalikan angka %v, mau nil", k.nama, hasil)
		}
		if !strings.Contains(err.Error(), k.sebut) {
			t.Errorf("%s: pesan tidak menyebut %q: %v", k.nama, k.sebut, err)
		}
	}
}

// Uji POSITIF: konversi yang dapat diselesaikan menghasilkan angkanya — dan
// nilai nol yang SAH tetap nol, bukan ditolak.
func TestKonversiSahMenghasilkanAngka(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{})

	h, err := l.KonversiMataUang(context.Background(), pelakuAda, desimal(t, "100"), desimal(t, "15000"), "USD")
	if err != nil {
		t.Fatalf("konversi sah ditolak: %v", err)
	}
	if h.String() != "1500000" {
		t.Errorf("hasil %s, mau 1500000", h)
	}

	nol, err := l.KonversiMataUang(context.Background(), pelakuAda, desimal(t, "0"), desimal(t, "15000"), "USD")
	if err != nil {
		t.Fatalf("nilai nol yang SAH ditolak: %v", err)
	}
	if !nol.IsZero() {
		t.Errorf("nol sah menjadi %s", nol)
	}
}

// ── Tiket 35 ──────────────────────────────────────────────────────────────

func TestKetentuanProporsionalTepatSatuTerisi(t *testing.T) {
	kasus := []struct {
		nama  string
		m     services.MasukanKetentuanProporsional
		sebut string
	}{
		{"keduanya terisi", services.MasukanKetentuanProporsional{
			IDLayer: 1, IDKelompokTreaty: 1, JenisTreaty: services.JenisQuotaShare,
			PersenQuotaShare: "25", JumlahLinesSurplus: lines(3)}, "KEDUANYA"},
		{"keduanya kosong", services.MasukanKetentuanProporsional{
			IDLayer: 1, IDKelompokTreaty: 1, JenisTreaty: services.JenisQuotaShare}, "KOSONG keduanya"},
		{"quota share tetapi lines yang terisi", services.MasukanKetentuanProporsional{
			IDLayer: 1, IDKelompokTreaty: 1, JenisTreaty: services.JenisQuotaShare,
			JumlahLinesSurplus: lines(3)}, "menuntut persenQuotaShare"},
		{"surplus tetapi persen yang terisi", services.MasukanKetentuanProporsional{
			IDLayer: 1, IDKelompokTreaty: 1, JenisTreaty: services.JenisSurplus,
			PersenQuotaShare: "25"}, "menuntut jumlahLinesSurplus"},
		{"jenis treaty di luar dua nilai", services.MasukanKetentuanProporsional{
			IDLayer: 1, IDKelompokTreaty: 1, JenisTreaty: "XOL", PersenQuotaShare: "25"}, "INV-30"},
	}
	for _, k := range kasus {
		g := &gudangTiruan{adaQS: true}
		l := services.LayananDengan(g)
		err := l.CatatKetentuanProporsional(context.Background(), pelakuAda, 7, k.m)
		if !errors.Is(err, services.ErrKetentuanProporsional) {
			t.Errorf("%s: mau ErrKetentuanProporsional, dapat %v", k.nama, err)
			continue
		}
		if !strings.Contains(err.Error(), k.sebut) {
			t.Errorf("%s: pesan tidak menyebut %q: %v", k.nama, k.sebut, err)
		}
		if len(g.dicatat) != 0 {
			t.Errorf("%s: baris tidak sah TETAP dicatat", k.nama)
		}
	}
}

// Uji POSITIF: kedua bentuk yang SAH diterima. Gerbang yang menolak segalanya
// lulus kelima uji negatif di atas.
func TestKeduaBentukKetentuanProporsionalSahDiterima(t *testing.T) {
	for _, m := range []services.MasukanKetentuanProporsional{
		{IDLayer: 1, IDKelompokTreaty: 1, JenisTreaty: services.JenisQuotaShare, PersenQuotaShare: "25"},
		{IDLayer: 1, IDKelompokTreaty: 1, JenisTreaty: services.JenisSurplus, JumlahLinesSurplus: lines(3)},
	} {
		g := &gudangTiruan{adaQS: true}
		l := services.LayananDengan(g)
		if err := l.CatatKetentuanProporsional(context.Background(), pelakuAda, 7, m); err != nil {
			t.Errorf("%s DITOLAK: %v", m.JenisTreaty, err)
			continue
		}
		if len(g.dicatat) != 1 {
			t.Errorf("%s: gudang menerima %d baris, mau 1", m.JenisTreaty, len(g.dicatat))
		}
	}
}

// ── Tiket 36 ──────────────────────────────────────────────────────────────

// Baris SURPLUS tanpa baris QUOTA_SHARE → kegagalan yang MENYEBUTKAN apa yang
// kurang; bukan galat tanpa isi, bukan hasil bernilai nol.
func TestSurplusTanpaQuotaShareDitolakBerketerangan(t *testing.T) {
	g := &gudangTiruan{adaQS: false}
	l := services.LayananDengan(g)

	err := l.CatatKetentuanProporsional(context.Background(), pelakuAda, 7,
		services.MasukanKetentuanProporsional{IDLayer: 1, IDKelompokTreaty: 1,
			JenisTreaty: services.JenisSurplus, JumlahLinesSurplus: lines(3)})

	if !errors.Is(err, services.ErrSurplusTanpaQuotaShare) {
		t.Fatalf("mau ErrSurplusTanpaQuotaShare, dapat %v", err)
	}
	for _, sebut := range []string{"QUOTA_SHARE", "retensi"} {
		if !strings.Contains(err.Error(), sebut) {
			t.Errorf("pesan tidak menyebut %q: %v", sebut, err)
		}
	}
	if len(g.dicatat) != 0 {
		t.Error("baris surplus TETAP dicatat")
	}
}

// Baris QUOTA_SHARE tidak menuntut apa pun — ia prasyaratnya, bukan yang
// bergantung.
func TestQuotaShareTidakMenuntutBarisLain(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{adaQS: false})

	err := l.CatatKetentuanProporsional(context.Background(), pelakuAda, 7,
		services.MasukanKetentuanProporsional{IDLayer: 1, IDKelompokTreaty: 1,
			JenisTreaty: services.JenisQuotaShare, PersenQuotaShare: "25"})

	if err != nil {
		t.Fatalf("baris quota share pertama DITOLAK: %v", err)
	}
}

// ── Tiket 43 ──────────────────────────────────────────────────────────────

func TestSakelarKeadaannyaTerlihatTanpaBasisData(t *testing.T) {
	var s services.SakelarPemindahan

	if k := s.Keadaan(); k.PenegakanMati {
		t.Error("sakelar lahir MATI; mau menyala")
	}
	if err := s.Matikan("pemindahan warisan batch 1"); err != nil {
		t.Fatal(err)
	}
	k := s.Keadaan()
	if !k.PenegakanMati || k.Alasan != "pemindahan warisan batch 1" {
		t.Errorf("keadaan tidak terbaca: %+v", k)
	}
}

// Alasan WAJIB: sakelar yang dapat dimatikan tanpa alasan akan ditemukan mati
// tanpa ada yang ingat kenapa.
func TestSakelarMenolakDimatikanTanpaAlasan(t *testing.T) {
	var s services.SakelarPemindahan
	if err := s.Matikan("   "); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Fatalf("mau ErrMasukanTidakSah, dapat %v", err)
	}
	if s.Keadaan().PenegakanMati {
		t.Error("sakelar mati walau alasannya ditolak")
	}
}

// Menyalakan kembali MEMERIKSA ULANG baris yang masuk selagi ia mati, dan
// menyebut yang gagal.
func TestNyalakanMemeriksaUlangBarisYangMasuk(t *testing.T) {
	var s services.SakelarPemindahan
	if err := s.Matikan("pemindahan"); err != nil {
		t.Fatal(err)
	}
	for _, id := range []int64{1, 2, 3} {
		s.Catat(id)
	}
	if n := s.Keadaan().BarisMasuk; n != 3 {
		t.Fatalf("mau 3 baris tercatat, dapat %d", n)
	}

	gagal, err := s.Nyalakan(func(id int64) error {
		if id == 2 {
			return errors.New("melanggar INV-53")
		}
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
	if len(gagal) != 1 || gagal[0] != 2 {
		t.Errorf("mau baris 2 gagal, dapat %v", gagal)
	}
	if k := s.Keadaan(); k.PenegakanMati || k.BarisMasuk != 0 {
		t.Errorf("sesudah dinyalakan: %+v", k)
	}
}

// Baris yang masuk selagi penegakan MENYALA tidak ikut tercatat.
func TestBarisSelagiMenyalaTidakTercatat(t *testing.T) {
	var s services.SakelarPemindahan
	s.Catat(1)
	if n := s.Keadaan().BarisMasuk; n != 0 {
		t.Errorf("mau nol, dapat %d", n)
	}
}

func TestNyalakanTanpaPemeriksaDitolak(t *testing.T) {
	var s services.SakelarPemindahan
	if _, err := s.Nyalakan(nil); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Fatalf("mau ErrMasukanTidakSah, dapat %v", err)
	}
}
