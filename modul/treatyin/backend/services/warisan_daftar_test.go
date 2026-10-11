package services_test

// Uji layar daftar WARISAN — terjemahan dan penomoran halaman.
//
// ⛔ Nol koneksi Oracle. Yang diuji di sini ATURAN TAMPILnya, dan aturan yang
// teruji tanpa basis data adalah aturan yang masih teruji ketika basis
// datanya tidak terjangkau.
//
// Bentuk data ujinya diambil dari sapuan 3 Oktober 2026 atas 1.854 baris
// `POOLDATA.TREATY_IN` — bukan dikarang:
//
//	ID                 tujuh digit, 1000001-1001856, nol ganda, nol bukan-angka
//	PROPORTIONTYPE     `Proportional` 1.079 · `NonProportional` 775 · nol NULL
//	STATUSAKSEPTASI    `Resolve Complete` 1.820 · `Accept` 12 · `Decline` 11 · NULL 11
//	COMMENCEMENT       delapan angka 1.853 · NULL 1
//	POSITIONUSERNAME   NULL 1.822 · terisi 32

import (
	"context"
	"errors"
	"fmt"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

func TestTanggalWarisanDiterjemahkan(t *testing.T) {
	// ⭐ Bentuk rujukan: `20190101` tampil `01/01/19`.
	for _, k := range []struct{ masuk, mau string }{
		{"20190101", "01/01/19"},
		{"20191231", "31/12/19"},
		{"20260917", "17/09/26"},
	} {
		if got := services.TanggalTampil(k.masuk); got != k.mau {
			t.Errorf("TanggalTampil(%q) = %q, mau %q", k.masuk, got, k.mau)
		}
	}
}

// ⛔ Uji NEGATIF yang paling penting di berkas ini: yang BUKAN delapan angka
// dikembalikan APA ADANYA. Menebak tanggal untuk nilai yang tidak berbentuk
// tanggal berarti menampilkan hari yang tidak pernah ada di baris mana pun,
// dan pembacanya tidak punya cara tahu ia karangan.
func TestTanggalBukanDelapanAngkaTidakDikarang(t *testing.T) {
	for _, masuk := range []string{"", "   ", "bukan tanggal", "2019", "201901011", "2019-01-01", "abcdefgh"} {
		if got := services.TanggalTampil(masuk); got != masuk {
			t.Errorf("TanggalTampil(%q) = %q — tanggal dikarang; mau apa adanya", masuk, got)
		}
	}
}

func TestSifatProporsiDiberiSpasi(t *testing.T) {
	if got := services.SifatProporsiTampil("NonProportional"); got != "Non Proportional" {
		t.Errorf("NonProportional -> %q, mau \"Non Proportional\"", got)
	}
	// `Proportional` TIDAK berubah — dan itu perlu diuji terpisah: aturan
	// yang menyisipkan spasi di sembarang tempat lulus uji di atas.
	if got := services.SifatProporsiTampil("Proportional"); got != "Proportional" {
		t.Errorf("Proportional berubah menjadi %q", got)
	}
	for _, lain := range []string{"", "   ", "Facultative", "NONPROPORTIONAL"} {
		if got := services.SifatProporsiTampil(lain); got != lain {
			t.Errorf("SifatProporsiTampil(%q) = %q, mau apa adanya", lain, got)
		}
	}
}

// barisWarisanUji membuat n baris dengan pengenal MENURUN, seperti yang
// kueri kembalikan.
func barisWarisanUji(n int) []models.BarisDaftarWarisan {
	out := make([]models.BarisDaftarWarisan, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, models.BarisDaftarWarisan{
			ID:                  fmt.Sprintf("%d", 1001856-i),
			NamaKontrak:         fmt.Sprintf("KONTRAK %d", i),
			SifatProporsiAsli:   "NonProportional",
			TanggalMulaiAsli:    "20190101",
			TanggalBerakhirAsli: "20191231",
			StatusAkseptasi:     "Resolve Complete",
		})
	}
	return out
}

func TestDaftarWarisanMenolakTanpaIdentitasSebelumGudang(t *testing.T) {
	g := &gudangTiruan{barisWarisan: barisWarisanUji(3), cacahWarisan: 3}
	l := services.LayananDengan(g)

	if _, err := l.DaftarKontrakWarisan(context.Background(), inti.Pelaku{}, 1); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Fatalf("mau ErrTanpaIdentitas, dapat %v", err)
	}
}

// Uji POSITIF: terjemahannya sampai ke baris yang dikembalikan, dan nilai
// aslinya TIDAK hilang.
func TestDaftarWarisanMenerjemahkanTanpaMembuangAsli(t *testing.T) {
	g := &gudangTiruan{barisWarisan: barisWarisanUji(1), cacahWarisan: 1}
	l := services.LayananDengan(g)

	h, err := l.DaftarKontrakWarisan(context.Background(), pelakuAda, 1)
	if err != nil {
		t.Fatalf("mau diterima, dapat %v", err)
	}
	b := h.Baris[0]
	if b.SifatProporsi != "Non Proportional" || b.SifatProporsiAsli != "NonProportional" {
		t.Errorf("sifat proporsi: tampil %q asli %q", b.SifatProporsi, b.SifatProporsiAsli)
	}
	if b.TanggalMulai != "01/01/19" || b.TanggalMulaiAsli != "20190101" {
		t.Errorf("tanggal mulai: tampil %q asli %q", b.TanggalMulai, b.TanggalMulaiAsli)
	}
	if b.TanggalBerakhir != "31/12/19" {
		t.Errorf("tanggal berakhir: %q", b.TanggalBerakhir)
	}
}

// ⭐ Penomoran halaman: halaman kedua mengembalikan baris yang BERBEDA dari
// halaman pertama. Penomoran yang mengembalikan baris yang sama di tiap
// halaman lulus setiap uji yang hanya memeriksa cacahnya.
func TestHalamanKeduaBerbedaDariPertama(t *testing.T) {
	n := services.UkuranHalamanWarisan*2 + 5
	g := &gudangTiruan{barisWarisan: barisWarisanUji(n), cacahWarisan: n}
	l := services.LayananDengan(g)

	h1, err := l.DaftarKontrakWarisan(context.Background(), pelakuAda, 1)
	if err != nil {
		t.Fatal(err)
	}
	h2, err := l.DaftarKontrakWarisan(context.Background(), pelakuAda, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(h1.Baris) != services.UkuranHalamanWarisan || len(h2.Baris) != services.UkuranHalamanWarisan {
		t.Fatalf("cacah halaman: %d dan %d, mau %d", len(h1.Baris), len(h2.Baris), services.UkuranHalamanWarisan)
	}
	if h1.Baris[0].ID == h2.Baris[0].ID {
		t.Errorf("halaman 1 dan 2 mengembalikan baris yang SAMA (%s)", h1.Baris[0].ID)
	}
	if h1.Total != n || h2.Total != n {
		t.Errorf("total: %d dan %d, mau %d", h1.Total, h2.Total, n)
	}
	// Halaman terakhir membawa SISANYA, bukan satu halaman penuh.
	h3, err := l.DaftarKontrakWarisan(context.Background(), pelakuAda, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(h3.Baris) != 5 {
		t.Errorf("halaman terakhir: %d baris, mau 5", len(h3.Baris))
	}
}

// Halaman nol, negatif, dan di luar jangkauan adalah PERMINTAAN yang wajar,
// bukan galat: pemakai yang menekan "berikutnya" sekali terlalu banyak tidak
// sedang melakukan kesalahan.
func TestHalamanDiLuarJangkauanMenjawabKosongBukanGalat(t *testing.T) {
	g := &gudangTiruan{barisWarisan: barisWarisanUji(3), cacahWarisan: 3}
	l := services.LayananDengan(g)

	for _, h := range []int{0, -5} {
		hasil, err := l.DaftarKontrakWarisan(context.Background(), pelakuAda, h)
		if err != nil {
			t.Fatalf("halaman %d: %v", h, err)
		}
		if hasil.Halaman != 1 || len(hasil.Baris) != 3 {
			t.Errorf("halaman %d dibulatkan ke %d dengan %d baris, mau 1 dan 3", h, hasil.Halaman, len(hasil.Baris))
		}
	}
	jauh, err := l.DaftarKontrakWarisan(context.Background(), pelakuAda, 99)
	if err != nil {
		t.Fatalf("halaman 99: %v", err)
	}
	if len(jauh.Baris) != 0 || jauh.Total != 3 {
		t.Errorf("halaman 99: %d baris, total %d — mau 0 dan 3", len(jauh.Baris), jauh.Total)
	}
	// Daftar KOSONG, bukan nil: pemanggil JSON tidak perlu membedakan
	// `null` dari `[]`.
	if jauh.Baris == nil {
		t.Error("Baris nil, mau daftar kosong")
	}
}

// ⛔ Aturan tombol aksi SENGAJA TIDAK diuji di sini, dan tidak ada
// padanannya di Go. Ia sudah hidup di `frontend/pages/DaftarKontrakTreatyIn.tsx`
// sebagai `aksiUntuk()`, lengkap dengan ujinya di `frontend/layar.test.ts`.
// Menuliskannya kembali di services akan membuat DUA tempat memutuskan satu
// hal, dan salah satunya akan basi — yang briefing ronde ini sebut
// "sambungkan, jangan tulis ulang aturannya".
//
// Yang dibawa dari sapuan ke sana: `STATUSAKSEPTASI` punya EMPAT nilai
// nyata, bukan satu — `Resolve Complete` (1.820), `Accept` (12),
// `Decline` (11), dan NULL (11). Ketiga yang terakhir masuk cabang "masih
// dapat disunting", dan itu 34 baris nyata.
