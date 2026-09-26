package services_test

// Tahap dan jalur balik - TANPA Oracle.
//
// Pemilik: tiket 08. Dibaca sesudah: tahap.go.

import (
	"context"
	"errors"
	"testing"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/services"
)

// TestPeranPemegangTahap - AC: tiap tahap dipegang perannya.
//
// `[terverifikasi work owner]` ADR-U-0002 tabel "Peran dan pemetaan tahap".
func TestPeranPemegangTahap(t *testing.T) {
	mau := map[models.Tahap]string{
		models.TahapInputRegister: models.PeranAdminLife,
		models.TahapOutstanding:   models.PeranAdminLife,
		models.TahapMedicalCheck:  models.PeranMedicalLife,
		// ⛔ Claim Analis dipegang SPV - ADR-U-0002, bukan Flow. Ronde pertama
		// menandainya `[terbuka]` dengan alasan "nol pyPosition di Flow";
		// diamnya korpus justru celah yang ADR itu TUTUP.
		models.TahapClaimAnalis: models.PeranSPVLife,
	}
	for tahap, peran := range mau {
		got, ada := models.PeranPemegangTahap(tahap)
		if !ada || got != peran {
			t.Errorf("tahap %v dipegang %q (ada=%v), mau %q", tahap, got, ada, peran)
		}
	}
	if _, ada := models.PeranPemegangTahap(models.TahapTidakDikenal); ada {
		t.Error("tahap tidak dikenal punya pemegang")
	}
}

// TestKolomPyPositionMenyimpanNamaPeran mengunci ralat terbesar tiket ini.
//
// `[terverifikasi]` `Flow/Register_Flow.xml` pecahan baris 582, 605, 628, 668,
// 731: `pyWorkPage.pyPosition` disetel ke nama PERAN, tidak pernah ke
// `"Assignment<n>"`. Ronde pertama menganggapnya pengenal shape - dan itu
// membuat setiap pembacaan baris nyata berakhir "tidak dikenal".
func TestKolomPyPositionMenyimpanNamaPeran(t *testing.T) {
	for peran, mau := range map[string]models.Tahap{
		models.PeranAdminLife:   models.TahapOutstanding,
		models.PeranMedicalLife: models.TahapMedicalCheck,
		models.PeranSPVLife:     models.TahapClaimAnalis,
	} {
		if got := models.TahapDariPeran(peran); got != mau {
			t.Errorf("TahapDariPeran(%q) = %v, mau %v", peran, got, mau)
		}
	}
	// ⛔ Pengenal shape BUKAN nilai kolom ini - dan harus tetap tak dikenal.
	for _, bukan := range []string{"Assignment1", "Assignment2", "", "Admin"} {
		if got := models.TahapDariPeran(bukan); got.Diketahui() {
			t.Errorf("TahapDariPeran(%q) = %v; itu bukan nama peran", bukan, got)
		}
	}
}

// TestSerahTerimaDanJalurBalik mengunci tangga kerja dan ADR-U-0002.
//
// ⛔ Ronde pertama menguji `JalurBalik` sebagai bendera BEBAS, dan bahkan
// menegaskan `{KeAdmin, KeMedical}` keduanya menyala - keadaan yang tidak
// mungkin, sebab keduanya menunjuk tujuan yang berbeda. Kini ia diturunkan
// dari pasangan perannya, sehingga keadaan mustahil itu tidak dapat dibentuk.
func TestSerahTerimaDanJalurBalik(t *testing.T) {
	kasus := []struct {
		dari, ke       string
		sah            bool
		mauAdm, mauMed bool
		apa            string
	}{
		{models.PeranAdminLife, models.PeranMedicalLife, true, false, false,
			"maju: Admin -> Medical"},
		{models.PeranMedicalLife, models.PeranSPVLife, true, false, false,
			"maju: Medical -> SPV"},
		{models.PeranMedicalLife, models.PeranAdminLife, true, true, false,
			"balik: Medical -> Admin"},
		{models.PeranSPVLife, models.PeranAdminLife, true, true, false,
			"balik: SPV -> Admin"},
		{models.PeranSPVLife, models.PeranMedicalLife, true, false, true,
			"balik: SPV -> Medical"},
		// ⛔ Lompatan yang tidak ada di tangga.
		{models.PeranAdminLife, models.PeranSPVLife, false, false, false,
			"lompat: Admin -> SPV"},
		{models.PeranAdminLife, models.PeranAdminLife, false, false, false,
			"ke dirinya sendiri"},
	}
	for _, k := range kasus {
		if got := models.SerahTerimaSah(k.dari, k.ke); got != k.sah {
			t.Errorf("%s: SerahTerimaSah = %v, mau %v", k.apa, got, k.sah)
		}
		adm, med := models.JalurBalikPeran(k.dari, k.ke)
		if adm != k.mauAdm || med != k.mauMed {
			t.Errorf("%s: jalur balik = (%v, %v), mau (%v, %v)",
				k.apa, adm, med, k.mauAdm, k.mauMed)
		}
		if adm && med {
			t.Errorf("%s: kedua penanda menyala sekaligus - mustahil", k.apa)
		}
	}
}

// TestPenandaSendtoBerbentukTeks - `[terverifikasi]` `IsSendtoAdmin`
// `pyConditionString` = `pyWorkPage.SendtoAdmin = 1` (pecahan baris 154);
// `IsSendtoMedical` = `pyWorkPage.SendtoMedical = 1` (baris 165).
func TestPenandaSendtoBerbentukTeks(t *testing.T) {
	adm, med := services.JalurBalik{KeAdmin: true}.NilaiSendto()
	if adm != "1" || med != "" {
		t.Errorf("ke Admin -> (%q, %q), mau (\"1\", \"\")", adm, med)
	}
	adm, med = services.JalurBalik{KeMedical: true}.NilaiSendto()
	if adm != "" || med != "1" {
		t.Errorf("ke Medical -> (%q, %q), mau (\"\", \"1\")", adm, med)
	}
	adm, med = services.JalurBalik{}.NilaiSendto()
	if adm != "" || med != "" {
		t.Errorf("bukan jalur balik -> (%q, %q), mau kosong keduanya", adm, med)
	}
}

// TestPindahTahapMenjagaPagarnya - identitas, tujuan, pengenal, Oracle.
func TestPindahTahapMenjagaPagarnya(t *testing.T) {
	svc := services.New(nil)
	ctx := context.Background()
	if err := svc.Tahap().Pindah(ctx, services.Pelaku{}, "CLM-1",
		models.TahapOutstanding, saatUji); !errors.Is(
		err, services.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: galat = %v, mau ErrTanpaIdentitas", err)
	}
	if err := svc.Tahap().Pindah(ctx, pelakuBerperan(services.PeranAdmin),
		"CLM-1", models.TahapTidakDikenal, saatUji); !errors.Is(
		err, services.ErrTahapTidakDikenal) {
		t.Errorf("tujuan asing: galat = %v, mau ErrTahapTidakDikenal", err)
	}
	if err := svc.Tahap().Pindah(ctx, pelakuBerperan(services.PeranAdmin),
		"  ", models.TahapOutstanding, saatUji); !errors.Is(
		err, services.ErrPermintaanTidakSah) {
		t.Errorf("pengenal kosong: galat = %v, mau ErrPermintaanTidakSah", err)
	}
	// ⚠️ Gerbang perannya menuntut pembacaan PY_POSITION, jadi tanpa Oracle ia
	// berhenti di sana. Peran per tahap sendiri diuji langsung di atas.
	if err := svc.Tahap().Pindah(ctx, pelakuBerperan(services.PeranAdmin),
		"CLM-1", models.TahapMedicalCheck, saatUji); !errors.Is(
		err, repository.ErrTanpaOracle) {
		t.Errorf("tanpa Oracle: galat = %v, mau ErrTanpaOracle", err)
	}
}
