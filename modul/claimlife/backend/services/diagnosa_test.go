package services

// Uji gerbang layanan diagnosa - butir bd.
//
// Tanpa Oracle: yang diperiksa URUTAN gerbangnya, dan itu justru bagian yang
// paling mudah salah. Gerbang yang benar di tempat yang salah membuat galat
// LAIN berbohong.

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/claimlife/backend/models"
)

// pelakuUjiDiagnosa - identitas ada, peran apa pun tidak relevan sebelum
// gerbang bentuk permintaan lewat.
var pelakuUjiDiagnosa = inti.Pelaku{AkunID: "UJI-1", Peran: []string{models.PeranAdminLife}}

func TestDiagnosaMenolakTanpaIdentitasLebihDulu(t *testing.T) {
	d := New(nil).Diagnosa()
	ctx := context.Background()
	// ⛔ Ketiga rute, bukan satu. Gerbang yang disalin ke tiga pintu adalah
	// tiga kesempatan untuk berbeda - dan uji satu pintu tidak akan tahu.
	if _, err := d.Tambah(ctx, inti.Pelaku{}, "K-1", "P-1"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("Tambah: %v, mau ErrTanpaIdentitas", err)
	}
	if err := d.Ubah(ctx, inti.Pelaku{}, "K-1", "P-1", 1, "", "", ""); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("Ubah: %v, mau ErrTanpaIdentitas", err)
	}
	if err := d.Hapus(ctx, inti.Pelaku{}, "K-1", "P-1", 1); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("Hapus: %v, mau ErrTanpaIdentitas", err)
	}
}

func TestDiagnosaMenjagaUrutanPagarnya(t *testing.T) {
	// ⛔ PENGENAL KOSONG dijawab "wajib diisi", BUKAN "ORACLE_DSN belum
	// dikonfigurasi". Cacat yang persis begitu pernah nyata di
	// `statusbaris.go`: penjaga yang benar, diletakkan di tempat yang membuat
	// galat lain berbohong. `TestUbahStatusMenjagaPagarnya` yang menangkapnya
	// di sana; uji ini menjaga pintu yang sama di sini.
	d := New(nil).Diagnosa()
	ctx := context.Background()
	for _, u := range []struct{ klaim, peserta string }{
		{"", "P-1"}, {"K-1", ""}, {"", ""}, {"   ", "P-1"},
	} {
		_, err := d.Tambah(ctx, pelakuUjiDiagnosa, u.klaim, u.peserta)
		if !errors.Is(err, galat.ErrPermintaanTidakSah) {
			t.Errorf("klaim=%q peserta=%q: %v, mau ErrPermintaanTidakSah",
				u.klaim, u.peserta, err)
		}
	}
}

func TestDiagnosaTanpaOracleSesudahBentuknyaSah(t *testing.T) {
	// Bentuk permintaan sah, identitas ada - barulah ketiadaan Oracle yang
	// menjadi galatnya. Urutan ini yang membuat pesan galat dapat dipercaya.
	_, err := New(nil).Diagnosa().Tambah(context.Background(),
		pelakuUjiDiagnosa, "K-1", "P-1")
	if err == nil || !strings.Contains(err.Error(), "ORACLE_DSN") {
		t.Errorf("galat = %v, mau menyebut ORACLE_DSN", err)
	}
}

// TestPencerminanDiagnosaDiTransaksiKeputusan adalah penjaga STATIK.
//
// ⛔ Ia membaca `statusbaris.go` dan menuntut `cerminkanKeDiagnosa` dipanggil
// DI DALAM `DalamTransaksi`. Kenapa statik dan bukan uji jalur: jalurnya
// menuntut Oracle, dan yang paling mudah hilang di sini bukan kebenaran
// query-nya melainkan LETAKNYA. Pencerminan yang tergeser ke luar transaksi
// tetap hijau di setiap uji db - sampai suatu hari ia gagal sendirian dan
// peserta yang sudah ditolak berdiri di samping diagnosa yang masih tampak
// Outstanding.
func TestPencerminanDiagnosaDiTransaksiKeputusan(t *testing.T) {
	isi, err := os.ReadFile("statusbaris.go")
	if err != nil {
		t.Fatal(err)
	}
	teks := string(isi)
	if !strings.Contains(teks, "cerminkanKeDiagnosa(ctx, st.svc, tx,") {
		t.Fatal("statusbaris.go tidak mencerminkan keputusan ke diagnosa; " +
			"SetSTS_Reject b241/b257 tidak tertiru")
	}
	// Letaknya: sesudah pembukaan transaksi, sebelum penutupnya.
	mulai := strings.Index(teks, "return st.svc.DalamTransaksi(ctx,")
	if mulai < 0 {
		t.Fatal("blok transaksi tidak ditemukan; penjaga ini yang usang")
	}
	if strings.Index(teks, "cerminkanKeDiagnosa(ctx, st.svc, tx,") < mulai {
		t.Error("pencerminan diagnosa berada DI LUAR transaksi keputusan")
	}
	// ⛔ Dan ia memakai kode BARU, bukan kode lama. Mencerminkan
	// `sasaranKodeLama` akan menuliskan keadaan sebelum keputusan - hijau di
	// setiap uji bentuk, dan salah pada setiap baris nyata.
	potong := teks[strings.Index(teks, "cerminkanKeDiagnosa(ctx, st.svc, tx,"):]
	if len(potong) > 200 {
		potong = potong[:200]
	}
	if !strings.Contains(potong, "baru.KodeStatus") {
		t.Errorf("pencerminan tidak memakai kode BARU:\n%s", potong)
	}
}

// TestPagariDiagnosaMemeriksaKetigaGerbang membaca `diagnosa.go` sendiri.
//
// ⚠️ Statik, dan sengaja sempit: yang dijaga adalah bahwa ketiga gerbang
// masih DIPANGGIL, bukan bahwa keduanya benar - kebenarannya diuji di
// `models` dan kelak di uji `db`. Gerbang yang hilang dari `pagari` tidak
// menggagalkan satu pun uji lain, sebab ketiga rute memanggil `pagari` dan
// `pagari`-lah yang menjadi kosong.
func TestPagariDiagnosaMemeriksaKetigaGerbang(t *testing.T) {
	isi, err := os.ReadFile("diagnosa.go")
	if err != nil {
		t.Fatal(err)
	}
	teks := string(isi)
	for _, panggilan := range []string{
		"WajibIdentitas(pelaku)",
		"PastikanKasusTerbuka(ctx, klaimID)",
		"models.TahapBergridPeserta(tahap)",
		"WajibPeran(pelaku, peranTahap)",
		"models.DiagnosaTerkunci(k.peserta.KodeStatus)",
	} {
		if !strings.Contains(teks, panggilan) {
			t.Errorf("pagari tidak lagi memanggil %s", panggilan)
		}
	}
	// ⛔ Dan ketiga rute lewat `pagari`. Rute yang memanggil repository
	// langsung adalah rute yang melewati seluruh gerbang sekaligus.
	for _, rute := range []string{
		"func (d *DiagnosaPeserta) Tambah(",
		"func (d *DiagnosaPeserta) Ubah(",
		"func (d *DiagnosaPeserta) Hapus(",
	} {
		i := strings.Index(teks, rute)
		if i < 0 {
			t.Errorf("rute %s hilang", rute)
			continue
		}
		badan := teks[i:]
		if j := strings.Index(badan[1:], "\nfunc "); j > 0 {
			badan = badan[:j]
		}
		if !strings.Contains(badan, "d.pagari(ctx, pelaku, klaimID, pesertaID)") {
			t.Errorf("%s tidak melewati pagari", rute)
		}
	}
}
