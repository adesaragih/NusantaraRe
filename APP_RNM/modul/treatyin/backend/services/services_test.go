package services_test

// Uji seam services tiket 15: identitas, himpunan tertutup, dan daftar kosong
// sebagai JAWABAN.
//
// Wewenang dan penolakan diuji DI SINI dengan `inti.Pelaku` langsung, bukan
// lewat jalur HTTP - aturannya tidak boleh bergantung pada `X-Pelaku`.

import (
	"context"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// gudangTiruan menjawab dari peta, dan MENCATAT himpunan yang diterimanya -
// supaya uji dapat membuktikan services tidak meneruskan yang tidak sah.
type gudangTiruan struct {
	isi     map[models.Himpunan][]models.Acuan
	diminta []models.Himpunan
	galat   error
}

func (g *gudangTiruan) DaftarAcuan(_ context.Context, h models.Himpunan) ([]models.Acuan, error) {
	g.diminta = append(g.diminta, h)
	if g.galat != nil {
		return nil, g.galat
	}
	return g.isi[h], nil
}

var pelakuAda = inti.Pelaku{AkunID: "AKUN-UJI"}

func TestDaftarAcuanMenolakPermintaanTanpaIdentitas(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	_, err := l.DaftarAcuan(context.Background(), inti.Pelaku{}, models.HimpunanMataUang)

	if !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Fatalf("mau ErrTanpaIdentitas, dapat %v", err)
	}
	// ⛔ Yang penting BUKAN hanya galatnya: gudang tidak boleh tersentuh sama
	// sekali. Penolakan yang tetap membaca basis data bukan penolakan.
	if len(g.diminta) != 0 {
		t.Errorf("gudang tersentuh %d kali walau identitas tidak ada", len(g.diminta))
	}
}

func TestDaftarAcuanMenolakHimpunanDiLuarEnam(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	for _, h := range []models.Himpunan{"", "KONTRAK", "mata_uang", "VERSI_KONTRAK", "../bahaya"} {
		_, err := l.DaftarAcuan(context.Background(), pelakuAda, h)
		if !errors.Is(err, services.ErrHimpunanTidakAda) {
			t.Errorf("himpunan %q: mau ErrHimpunanTidakAda, dapat %v", h, err)
		}
	}
	// Nama tabel tidak pernah datang dari teks permintaan: yang tidak sah
	// berhenti SEBELUM gudang.
	if len(g.diminta) != 0 {
		t.Errorf("gudang menerima %v; himpunan tak sah tidak boleh diteruskan", g.diminta)
	}
}

// Uji POSITIF - dan ia yang menangkap penyaring yang terlalu ketat. Penyaring
// yang menolak salah satu dari enam lulus setiap uji negatif di atas.
func TestKeenamHimpunanDiterima(t *testing.T) {
	enam := []models.Himpunan{
		models.HimpunanMataUang, models.HimpunanJenisPotongan, models.HimpunanKelasBisnis,
		models.HimpunanKelompokTreaty, models.HimpunanBahaya, models.HimpunanJenisReasuransi,
	}
	g := &gudangTiruan{isi: map[models.Himpunan][]models.Acuan{}}
	for _, h := range enam {
		g.isi[h] = []models.Acuan{{ID: 1, Kode: "K1", Nama: "N1", Aktif: "1"}}
	}
	l := services.LayananDengan(g)

	for _, h := range enam {
		baris, err := l.DaftarAcuan(context.Background(), pelakuAda, h)
		if err != nil {
			t.Errorf("himpunan %q ditolak: %v", h, err)
			continue
		}
		if len(baris) != 1 {
			t.Errorf("himpunan %q: mau 1 baris, dapat %d", h, len(baris))
		}
	}
	if len(g.diminta) != len(enam) {
		t.Errorf("gudang dipanggil %d kali, mau %d", len(g.diminta), len(enam))
	}
}

// Daftar kosong adalah JAWABAN, bukan galat: keenam tabel berdiri kosong
// sampai tiket 44 memindahkan isinya dari sistem lama.
func TestTabelAcuanKosongBukanGalat(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{isi: map[models.Himpunan][]models.Acuan{}})

	baris, err := l.DaftarAcuan(context.Background(), pelakuAda, models.HimpunanBahaya)

	if err != nil {
		t.Fatalf("tabel acuan kosong menghasilkan galat: %v", err)
	}
	if len(baris) != 0 {
		t.Fatalf("mau nol baris, dapat %d", len(baris))
	}
}

func TestGalatGudangDiteruskan(t *testing.T) {
	bocor := errors.New("oracle mati")
	l := services.LayananDengan(&gudangTiruan{galat: bocor})

	_, err := l.DaftarAcuan(context.Background(), pelakuAda, models.HimpunanMataUang)

	if !errors.Is(err, bocor) {
		t.Fatalf("galat gudang tidak diteruskan; dapat %v", err)
	}
}

func TestBersusunHanyaJenisReasuransi(t *testing.T) {
	if !models.HimpunanJenisReasuransi.Bersusun() {
		t.Error("JENIS_REASURANSI bersusun (§10.6) tetapi Bersusun() false")
	}
	for _, h := range []models.Himpunan{
		models.HimpunanMataUang, models.HimpunanJenisPotongan, models.HimpunanKelasBisnis,
		models.HimpunanKelompokTreaty, models.HimpunanBahaya,
	} {
		if h.Bersusun() {
			t.Errorf("%q tidak bersusun tetapi Bersusun() true", h)
		}
	}
}
