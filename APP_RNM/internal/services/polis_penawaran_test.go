package services

// Uji gerbang keputusan penawaran - tiket 01 PremiumList Life.
//
// Tanpa Oracle: yang diperiksa URUTAN gerbangnya dan janji "nol aturan
// otomatis" (AC 6).

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"nusantarare/internal/models"
)

var pelakuUjiPolis = Pelaku{AkunID: "UJI-POLIS"}

func TestPenawaranMenolakTanpaIdentitasLebihDulu(t *testing.T) {
	p := New(nil).Penawaran()
	ctx := context.Background()
	if _, err := p.Putuskan(ctx, Pelaku{}, "POL-1", models.KeputusanConfirm,
		time.Now()); !errors.Is(err, ErrTanpaIdentitas) {
		t.Errorf("Putuskan: %v, mau ErrTanpaIdentitas", err)
	}
	if _, err := p.Golongkan(ctx, Pelaku{}, "POL-1", models.LanjutOffer,
		time.Now()); !errors.Is(err, ErrTanpaIdentitas) {
		t.Errorf("Golongkan: %v, mau ErrTanpaIdentitas", err)
	}
}

func TestPenawaranMenjagaUrutanPagarnya(t *testing.T) {
	// ⛔ Pengenal kosong dijawab "wajib diisi", BUKAN "ORACLE_DSN belum
	// dikonfigurasi". Cacat yang persis begitu pernah nyata di Claim Life:
	// penjaga yang benar, diletakkan di tempat yang membuat galat lain
	// berbohong.
	p := New(nil).Penawaran()
	for _, id := range []string{"", "   "} {
		_, err := p.Putuskan(context.Background(), pelakuUjiPolis, id,
			models.KeputusanConfirm, time.Now())
		if !errors.Is(err, ErrPermintaanTidakSah) {
			t.Errorf("pengenal %q: %v, mau ErrPermintaanTidakSah", id, err)
		}
	}
	// Bentuk sah -> barulah ketiadaan Oracle yang menjadi galatnya.
	_, err := p.Putuskan(context.Background(), pelakuUjiPolis, "POL-1",
		models.KeputusanConfirm, time.Now())
	if err == nil || !strings.Contains(err.Error(), "ORACLE_DSN") {
		t.Errorf("galat = %v, mau menyebut ORACLE_DSN", err)
	}
}

// TestNolAturanOtomatisMenetapkanKeputusan adalah penjaga AC 6.
//
// ⛔ Kedua decision table Pega mengekspor NOL baris keputusan; yang ditiru
// AKIBAT keputusan, bukan formula yang memilihnya. Layanan ini karena itu
// TIDAK boleh menghitung keputusannya sendiri - ia selalu menerimanya sebagai
// parameter. Penjaga ini membaca berkasnya dan menolak literal keputusan di
// luar daftar konstanta models.
func TestNolAturanOtomatisMenetapkanKeputusan(t *testing.T) {
	isi, err := os.ReadFile("polis_penawaran.go")
	if err != nil {
		t.Fatal(err)
	}
	teks := string(isi)
	// Keputusan hanya boleh disebut lewat konstanta models - nol literal.
	for _, literal := range []string{`"Confirm"`, `"Reject"`, `"Decline"`} {
		if strings.Contains(teks, literal) {
			t.Errorf("polis_penawaran.go memuat literal %s; keputusan datang dari "+
				"pengguna lewat models.Keputusan*, tidak dihitung di sini", literal)
		}
	}
	// Dan ia memang MENERIMA keputusannya, bukan menurunkannya.
	if !strings.Contains(teks, "models.TransisiPenawaran(keadaan.Position, keputusan)") {
		t.Error("Putuskan tidak meneruskan keputusan pengguna ke tabel transisi")
	}
}

func TestConfirmDiPenawaranTidakMenulisApaPun(t *testing.T) {
	// ⛔ `Confirm` di tahap penawaran hanya memindahkan kendali ke
	// `Decision3`. Menuliskan tahap di sini berarti kasus berpindah SEBELUM
	// ada yang memutuskan ke mana - dan penggolongnya kemudian memindahkan
	// lagi dari tahap yang bukan tahapnya.
	akibat, err := models.TransisiPenawaran(models.TahapPolisPenawaran,
		models.KeputusanConfirm)
	if err != nil {
		t.Fatal(err)
	}
	if !akibat.MenungguPenggolong {
		t.Fatal("Confirm di penawaran tidak menunggu penggolong")
	}
	if akibat.TahapTujuan != "" || akibat.Ditutup() {
		t.Errorf("akibat = %+v, mau kosong selain penanda menunggu", akibat)
	}
}
