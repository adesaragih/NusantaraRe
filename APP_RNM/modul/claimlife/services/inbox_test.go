package services_test

// Uji gerbang kotak masuk - F0.4.
//
// Yang diuji: siapa boleh melihat antrian mana, antrian mana yang pribadi,
// dan batas halamannya. Pembacaan Oracle-nya diuji terpisah oleh uji berskema.

import (
	"context"
	"errors"
	"testing"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/modul/claimlife/models"
	"nusantarare/modul/claimlife/repository"
	"nusantarare/modul/claimlife/services"
)

func pelakuUji(peran ...string) inti.Pelaku {
	return inti.Pelaku{AkunID: "UJI-AKUN", Peran: peran}
}

// `[terverifikasi]` `Register_Flow.xml`: Input Register `WorkList` 1511 +
// `Current operator` 1508; Outstanding `WorkList` 1300 + `ToWorklist` 1351;
// Medical Check `WorkBasket` 993; Claim Analis `WorkBasket` 1119.
func TestAntrianPribadiHanyaKeduaTahapAdmin(t *testing.T) {
	kasus := []struct {
		tahap   models.Tahap
		pribadi bool
	}{
		{models.TahapInputRegister, true},
		{models.TahapOutstanding, true},
		{models.TahapMedicalCheck, false},
		{models.TahapClaimAnalis, false},
		{models.TahapTidakDikenal, false},
	}
	for _, k := range kasus {
		if got := services.AntrianPribadi(k.tahap); got != k.pribadi {
			t.Errorf("AntrianPribadi(%v) = %v, mau %v", k.tahap, got, k.pribadi)
		}
	}
}

func TestTahapTerlihatPerPeran(t *testing.T) {
	kasus := []struct {
		apa    string
		pelaku inti.Pelaku
		mau    []models.Tahap
	}{
		{
			"Admin melihat KEDUA antrian Admin",
			pelakuUji(models.PeranAdminLife),
			[]models.Tahap{models.TahapInputRegister, models.TahapOutstanding},
		},
		{
			"Medical Advisor hanya Medical Check",
			pelakuUji(models.PeranMedicalLife),
			[]models.Tahap{models.TahapMedicalCheck},
		},
		{
			"SPV hanya Claim Analis",
			pelakuUji(models.PeranSPVLife),
			[]models.Tahap{models.TahapClaimAnalis},
		},
		{
			// ⛔ Peran ganda melihat GABUNGAN. `pelakuDari` memecah `X-Peran`
			// pada koma justru supaya ini mungkin.
			"peran ganda melihat gabungan, berurut tangga",
			pelakuUji(models.PeranSPVLife, models.PeranAdminLife),
			[]models.Tahap{
				models.TahapInputRegister, models.TahapOutstanding,
				models.TahapClaimAnalis,
			},
		},
		{
			"tanpa peran tidak melihat satu pun tab",
			pelakuUji(),
			nil,
		},
	}
	for _, k := range kasus {
		got := services.TahapTerlihat(k.pelaku)
		if len(got) != len(k.mau) {
			t.Errorf("%s: %d tahap, mau %d (%v)", k.apa, len(got), len(k.mau), got)
			continue
		}
		for i := range got {
			if got[i] != k.mau[i] {
				t.Errorf("%s: tahap[%d] = %v, mau %v", k.apa, i, got[i], k.mau[i])
			}
		}
	}
}

// Batas halaman dari `InboxPremiumList.xml` 592 dan 942.
func TestBatasUkuranHalaman(t *testing.T) {
	kasus := []struct{ minta, mau int }{
		{0, repository.UkuranHalamanBawaan},
		{-5, repository.UkuranHalamanBawaan},
		{1, 1},
		{50, 50},
		{500, 500},
		// ⛔ DIPANGKAS, bukan ditolak: pemanggil yang meminta 10.000 hampir
		// selalu salah tulis, dan menolaknya membuat layar kosong.
		{501, repository.UkuranHalamanMaksimum},
		{10000, repository.UkuranHalamanMaksimum},
	}
	for _, k := range kasus {
		if got := services.BatasUkuran(k.minta); got != k.mau {
			t.Errorf("BatasUkuran(%d) = %d, mau %d", k.minta, got, k.mau)
		}
	}
}

func TestInboxMenolakSebelumMenyentuhOracle(t *testing.T) {
	in := services.New(nil).KotakMasuk()
	ctx := context.Background()

	// ⛔ Pelaku KOSONG ditolak, bukan dijawab daftar kosong. Daftar kosong
	// terbaca "tidak ada pekerjaan" - kalimat yang berbeda artinya dari
	// "saya tidak tahu siapa Anda".
	if _, err := in.Ambil(ctx, inti.Pelaku{}, models.TahapOutstanding, 0, 50); !errors.Is(
		err, inti.ErrTanpaIdentitas) {
		t.Errorf("pelaku kosong: galat = %v, mau ErrTanpaIdentitas", err)
	}

	// Tahap di luar keempat yang ada.
	if _, err := in.Ambil(ctx, pelakuUji(models.PeranAdminLife),
		models.TahapTidakDikenal, 0, 50); !errors.Is(err, services.ErrTahapTidakSah) {
		t.Errorf("tahap tidak dikenal: galat = %v, mau ErrTahapTidakSah", err)
	}

	// ⛔ Peran yang tidak memegang tahap itu DITOLAK - bukan diberi daftar
	// kosong. Admin tidak boleh mengintip antrian Medical Check.
	if _, err := in.Ambil(ctx, pelakuUji(models.PeranAdminLife),
		models.TahapMedicalCheck, 0, 50); !errors.Is(err, inti.ErrTanpaWewenang) {
		t.Errorf("peran salah: galat = %v, mau ErrTanpaWewenang", err)
	}

	// Baru sesudah semua gerbang lolos, ketiadaan Oracle yang terasa.
	if _, err := in.Ambil(ctx, pelakuUji(models.PeranAdminLife),
		models.TahapOutstanding, 0, 50); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("tanpa Oracle: galat = %v, mau ErrTanpaOracle", err)
	}
}
