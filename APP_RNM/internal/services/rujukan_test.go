package services_test

// Uji gerbang dropdown Register - A3 kelompok Register.

import (
	"context"
	"errors"
	"testing"

	"nusantarare/internal/repository"
	"nusantarare/internal/services"
)

func TestRujukanMenolakSebelumMenyentuhOracle(t *testing.T) {
	rj := services.New(nil).SumberRujukan()
	ctx := context.Background()
	pelaku := services.Pelaku{AkunID: "UJI-AKUN", Peran: []string{"ReasLifeAdmin"}}

	// ⛔ Anonim ditolak: isinya memuat NAMA - mitra dan petugas pemasaran.
	if _, err := rj.Cari(ctx, services.Pelaku{}, services.RujukanCeding, "abc"); !errors.Is(
		err, services.ErrTanpaIdentitas) {
		t.Errorf("anonim: galat = %v, mau ErrTanpaIdentitas", err)
	}

	// ⚠️ Ketikan terlalu pendek menjawab daftar KOSONG, bukan galat: pemakai
	// yang baru mengetik satu huruf belum salah melakukan apa pun.
	for _, pendek := range []string{"", " ", "a", " x "} {
		hasil, err := rj.Cari(ctx, pelaku, services.RujukanCeding, pendek)
		if err != nil {
			t.Errorf("cari %q: galat %v, mau nil", pendek, err)
		}
		if len(hasil) != 0 {
			t.Errorf("cari %q: %d baris, mau 0", pendek, len(hasil))
		}
	}

	// Jenis di luar ketiga yang ada.
	if _, err := rj.Cari(ctx, pelaku, "penyakit", "abc"); !errors.Is(
		err, services.ErrJenisRujukanTidakDikenal) {
		t.Errorf("jenis asing: galat = %v, mau ErrJenisRujukanTidakDikenal", err)
	}

	// Baru sesudah gerbangnya lolos, ketiadaan Oracle yang terasa.
	for _, jenis := range []string{
		services.RujukanCeding, services.RujukanBisnis, services.RujukanMarketing,
	} {
		if _, err := rj.Cari(ctx, pelaku, jenis, "abc"); !errors.Is(
			err, repository.ErrTanpaOracle) {
			t.Errorf("%s tanpa Oracle: galat = %v, mau ErrTanpaOracle", jenis, err)
		}
	}
}

// Ketiga jenis, dan TIDAK LEBIH.
//
// ⛔ Himpunan tertutup. Nama objek yang datang dari pemakai adalah injeksi
// lewat pintu yang tidak dijaga penanda `:1`; petanya hidup di kode.
func TestJenisRujukanHimpunanTertutup(t *testing.T) {
	semua := []string{
		services.RujukanCeding, services.RujukanBisnis, services.RujukanMarketing,
	}
	if len(semua) != 3 {
		t.Fatalf("%d jenis, mau 3", len(semua))
	}
	lihat := map[string]bool{}
	for _, j := range semua {
		if lihat[j] {
			t.Errorf("jenis %q ganda", j)
		}
		lihat[j] = true
	}
}
