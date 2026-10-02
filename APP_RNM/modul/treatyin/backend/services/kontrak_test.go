package services_test

// Uji seam services tiket 14 - jalur buat dan baca kontrak.
//
// ⭐ Di sinilah INV-53 dan INV-29 akhirnya punya uji. Migrasi 401 menyatakan
// keduanya ditegakkan di services dan menulis "DITAGIH: tiket lapisan
// aplikasi"; berkas ini menagihnya.

import (
	"context"
	"errors"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// masukanSah adalah kontrak yang SELURUH gerbangnya lolos. Tiap uji negatif
// merusak SATU ruas darinya, supaya yang diuji benar-benar ruas itu.
func masukanSah() services.MasukanKontrak {
	return services.MasukanKontrak{
		IDCedant:            1,
		IDAsalBisnis:        2,
		SifatProporsi:       models.SifatNonProporsional,
		TanggalMulai:        "2026-01-01",
		TanggalBerakhir:     "2026-12-31",
		NamaKontrak:         "Kontrak Uji",
		KodeMataUangKontrak: 3,
		PersenBagianNure:    "12.5",
		BagianNureSeragam:   "0",
		MemakaiBordereaux:   "0",
		CaraPembukuan:       "X",
		MemakaiProrata:      "0",
		RetroBerganda:       "0",
		KeadaanSiklusHidup:  "DRAFT",
	}
}

// Uji POSITIF, dan ia yang menangkap gerbang yang TERLALU KETAT: masukan sah
// diterima, dan pengenal yang sistem berikan dikembalikan.
func TestBuatKontrakSahDiterima(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	h, err := l.BuatKontrak(context.Background(), pelakuAda, masukanSah())

	if err != nil {
		t.Fatalf("masukan sah DITOLAK: %v", err)
	}
	if h.IDKontrak == 0 || h.IDVersi == 0 {
		t.Errorf("pengenal tidak dikembalikan: %+v", h)
	}
	if len(g.dibuat) != 1 {
		t.Fatalf("gudang menerima %d kontrak, mau 1", len(g.dibuat))
	}
	// Versi PERTAMA selalu bernomor 1.
	v := g.versiDibuat[0]
	if v.NomorUrutVersi == nil || *v.NomorUrutVersi != 1 {
		t.Errorf("versi pertama bernomor %v, mau 1", v.NomorUrutVersi)
	}
}

func TestBuatKontrakMenolakTanpaIdentitas(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	_, err := l.BuatKontrak(context.Background(), inti.Pelaku{}, masukanSah())

	if !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Fatalf("mau ErrTanpaIdentitas, dapat %v", err)
	}
	if len(g.dibuat) != 0 {
		t.Error("gudang tersentuh walau identitas tidak ada")
	}
}

// INV-53 - batas INKLUSIF keduanya (ADR-0022).
func TestInv53MenolakTanggalTerbalik(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)
	m := masukanSah()
	m.TanggalMulai, m.TanggalBerakhir = "2026-12-31", "2026-01-01"

	_, err := l.BuatKontrak(context.Background(), pelakuAda, m)

	if !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Fatalf("mau ErrMasukanTidakSah, dapat %v", err)
	}
	if !strings.Contains(err.Error(), "INV-53") {
		t.Errorf("pesan tidak menyebut INV-53: %v", err)
	}
	if len(g.dibuat) != 0 {
		t.Error("masukan yang ditolak TETAP diteruskan ke gudang")
	}
}

// INV-53 uji POSITIF - batasnya INKLUSIF: mulai dan berakhir pada hari yang
// SAMA diterima. Gerbang yang memakai <= terbalik akan menolaknya, dan uji
// negatif saja tidak pernah menangkapnya.
func TestInv53MenerimaPeriodeSatuHari(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{})
	m := masukanSah()
	m.TanggalMulai, m.TanggalBerakhir = "2026-06-01", "2026-06-01"

	if _, err := l.BuatKontrak(context.Background(), pelakuAda, m); err != nil {
		t.Fatalf("periode satu hari DITOLAK, batas inklusif: %v", err)
	}
}

// INV-29 - SIFAT_PROPORSI dua nilai, dan hanya dua.
func TestInv29MenolakSifatProporsiDiLuarDuaNilai(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{})
	for _, sifat := range []string{"", "proporsional", "QUOTA_SHARE", "NONPROPORSIONAL"} {
		m := masukanSah()
		m.SifatProporsi = sifat
		_, err := l.BuatKontrak(context.Background(), pelakuAda, m)
		if !errors.Is(err, services.ErrMasukanTidakSah) {
			t.Errorf("sifatProporsi %q diterima, mau ditolak", sifat)
		}
	}
}

// Uji POSITIF INV-29: KEDUA nilai sah diterima. Gerbang yang hanya mengizinkan
// satu lulus setiap uji negatif di atas.
func TestInv29MenerimaKeduaNilaiSah(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{})
	for _, sifat := range []string{models.SifatProporsional, models.SifatNonProporsional} {
		m := masukanSah()
		m.SifatProporsi = sifat
		if _, err := l.BuatKontrak(context.Background(), pelakuAda, m); err != nil {
			t.Errorf("sifatProporsi %q DITOLAK: %v", sifat, err)
		}
	}
}

// Gerbang mengumpulkan SELURUH pelanggaran, bukan berhenti di yang pertama.
func TestMasukanRusakMelaporkanSeluruhPelanggaran(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{})
	m := masukanSah()
	m.SifatProporsi = "entah"
	m.NamaKontrak = "  "
	m.TanggalBerakhir = "bukan tanggal"

	_, err := l.BuatKontrak(context.Background(), pelakuAda, m)

	if err == nil {
		t.Fatal("masukan rusak diterima")
	}
	for _, potongan := range []string{"sifatProporsi", "namaKontrak", "tanggalBerakhir"} {
		if !strings.Contains(err.Error(), potongan) {
			t.Errorf("pesan tidak menyebut %s: %v", potongan, err)
		}
	}
}

func TestBacaKontrakMenolakPengenalTidakSah(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)
	for _, id := range []int64{0, -1} {
		if _, err := l.BacaKontrak(context.Background(), pelakuAda, id); !errors.Is(err, services.ErrMasukanTidakSah) {
			t.Errorf("pengenal %d: mau ErrMasukanTidakSah, dapat %v", id, err)
		}
	}
	if len(g.dibaca) != 0 {
		t.Error("gudang tersentuh untuk pengenal tak sah")
	}
}

// Tiket 14: kontrak ditemukan kembali dengan pengenalnya, dan lapisan bekunya
// dibaca SEKALI - tidak disalin ke tiap versi.
func TestBacaKontrakMengembalikanLapisanBekuSekali(t *testing.T) {
	satu, dua := int64(1), int64(2)
	g := &gudangTiruan{kontrak: models.KontrakDenganVersi{
		Kontrak: models.Kontrak{ID: 777, SifatProporsi: models.SifatNonProporsional},
		Versi: []models.VersiKontrak{
			{ID: 888, IDKontrak: 777, NomorUrutVersi: &satu},
			{ID: 889, IDKontrak: 777, NomorUrutVersi: &dua},
		},
	}}
	l := services.LayananDengan(g)

	k, err := l.BacaKontrak(context.Background(), pelakuAda, 777)

	if err != nil {
		t.Fatalf("baca gagal: %v", err)
	}
	if k.Kontrak.ID != 777 || len(k.Versi) != 2 {
		t.Errorf("bentuk tidak sesuai: %+v", k)
	}
	if k.Kontrak.SifatProporsi == "" {
		t.Error("lapisan beku kosong")
	}
}
