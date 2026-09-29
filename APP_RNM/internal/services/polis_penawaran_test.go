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
	"nusantarare/inti"
	"nusantarare/inti/galat"
)

var pelakuUjiPolis = inti.Pelaku{AkunID: "UJI-POLIS"}

func TestPenawaranMenolakTanpaIdentitasLebihDulu(t *testing.T) {
	p := New(nil).Penawaran()
	ctx := context.Background()
	if _, err := p.Putuskan(ctx, inti.Pelaku{}, "POL-1", models.KeputusanConfirm,
		time.Now()); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("Putuskan: %v, mau ErrTanpaIdentitas", err)
	}
}

// TestDecision3DirutekanDariBendera - GILIRAN-14 butir bq.
//
// ⛔ `Confirm` di tahap penawaran menerapkan hasil `IsFlagOnGoingPolicy` atas
// bendera KASUS, di pintu yang sama - bukan menunggu pengguna memilih
// `Offer`/`Premium` lewat rute kedua. Rute dan layanan penggolong manual
// (`Golongkan`) dibuang: kode mati.
func TestDecision3DirutekanDariBendera(t *testing.T) {
	isi, err := os.ReadFile("polis_penawaran.go")
	if err != nil {
		t.Fatal(err)
	}
	teks := string(isi)
	if !strings.Contains(teks, ".Bendera(ctx, keadaan.ID)") ||
		!strings.Contains(teks, "models.PenggolongOtomatis(flag)") {
		t.Error("Putuskan tidak merutekan Decision3 dari bendera kasus")
	}
	if strings.Contains(teks, "func (p *Penawaran) Golongkan") {
		t.Error("penggolong manual masih ada; Decision3 tidak ditanyakan ke pengguna")
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
		if !errors.Is(err, galat.ErrPermintaanTidakSah) {
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
// ⛔ Keputusan `Confirm`/`Reject`/`Decline` datang dari PENGGUNA - layanan ini
// tidak menghitungnya. Penjaga ini membaca berkasnya dan menolak literal
// keputusan di luar daftar konstanta models.
//
// ⚠️ RALAT 29-09-2026: premis lama "kedua decision table mengekspor NOL
// baris" keliru (GILIRAN-13). Untuk `IsLifeAccepted` (Decision1/2) keputusan
// manual tetap selaras; `IsFlagOnGoingPolicy` (Decision3) kini dirutekan dari
// bendera - lihat TestDecision3DirutekanDariBendera.
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
	//
	// ⚠️ Tahapnya dibaca dari Status, bukan Position -
	// ralat 28-09-2026; lihat KeadaanPolis di repository.
	if !strings.Contains(teks, "models.TransisiPenawaran(keadaan.Status, keputusan)") {
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
	if !akibat.KeDecision3 {
		t.Fatal("Confirm di penawaran tidak menunggu penggolong")
	}
	if akibat.TahapTujuan != "" || akibat.Ditutup() {
		t.Errorf("akibat = %+v, mau kosong selain penanda menunggu", akibat)
	}
}
