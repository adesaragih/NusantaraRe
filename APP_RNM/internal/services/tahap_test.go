package services_test

// Tahap dan jalur balik - TANPA Oracle.
//
// Pemilik: tiket 08. Dibaca sesudah: tahap.go.

import (
	"context"
	"errors"
	"testing"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/galat"
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
// dari pasangan TAHAPnya, sehingga keadaan mustahil itu tidak dapat dibentuk.
//
// ⛔ BUTIR at: tangganya kini antar-TAHAP, bukan antar-PERAN. Peta peran
// tidak dapat menyatakan Input Register ⇄ Outstanding Claim - keduanya
// dipegang `ReasLifeAdmin` - sehingga satu perpindahan yang XML tunjukkan
// dengan terang (`Send Back to Register`, `InputOSClaimLife.xml:21404`)
// tidak punya tempat di dalamnya.
func TestSerahTerimaDanJalurBalik(t *testing.T) {
	kasus := []struct {
		dari, ke       models.Tahap
		sah            bool
		mauAdm, mauMed bool
		apa            string
	}{
		// ⛔ Admin → Admin: kemajuan Register → Outstanding saat borang
		// register disimpan, dan kembalinya lewat `Send Back to Register`.
		// KEDUANYA sah, dan KEDUANYA bukan jalur balik antarperan:
		// tidak ada yang "dikembalikan" kepada siapa pun.
		{models.TahapInputRegister, models.TahapOutstanding, true, false, false,
			"maju: Input Register -> Outstanding"},
		{models.TahapOutstanding, models.TahapInputRegister, true, false, false,
			"balik dalam tangan Admin: Outstanding -> Input Register"},
		{models.TahapOutstanding, models.TahapMedicalCheck, true, false, false,
			"maju: Outstanding -> Medical"},
		{models.TahapMedicalCheck, models.TahapClaimAnalis, true, false, false,
			"maju: Medical -> Analis"},
		{models.TahapMedicalCheck, models.TahapOutstanding, true, true, false,
			"balik: Medical -> Admin"},
		{models.TahapClaimAnalis, models.TahapOutstanding, true, true, false,
			"balik: Analis -> Admin"},
		{models.TahapClaimAnalis, models.TahapMedicalCheck, true, false, true,
			"balik: Analis -> Medical"},
		// ⛔ Lompatan yang tidak ada di tangga.
		{models.TahapOutstanding, models.TahapClaimAnalis, false, false, false,
			"lompat: Outstanding -> Analis"},
		{models.TahapInputRegister, models.TahapMedicalCheck, false, false, false,
			"lompat: Register -> Medical"},
		{models.TahapOutstanding, models.TahapOutstanding, false, false, false,
			"ke dirinya sendiri"},
		{models.TahapTidakDikenal, models.TahapOutstanding, false, false, false,
			"dari tahap yang tidak dikenal"},
	}
	for _, k := range kasus {
		if got := models.SerahTerimaSah(k.dari, k.ke); got != k.sah {
			t.Errorf("%s: SerahTerimaSah = %v, mau %v", k.apa, got, k.sah)
		}
		adm, med := models.JalurBalikTahap(k.dari, k.ke)
		if adm != k.mauAdm || med != k.mauMed {
			t.Errorf("%s: jalur balik = (%v, %v), mau (%v, %v)",
				k.apa, adm, med, k.mauAdm, k.mauMed)
		}
		if adm && med {
			t.Errorf("%s: kedua penanda menyala sekaligus - mustahil", k.apa)
		}
	}
}

// TestTahapDariNamaKebalikanString - butir at.
//
// ⛔ Perjalanan bolak-balik. Kolom `TAHAP` menyimpan nama yang `String()`
// hasilkan; bila `TahapDariNama` tidak dapat membacanya kembali, kasus yang
// sudah tersimpan menjadi tak terbaca - dan gerbang tangganya diam-diam
// jatuh ke cadangan `PY_POSITION` yang tidak dapat membedakan kedua tahap
// Admin.
func TestTahapDariNamaKebalikanString(t *testing.T) {
	semua := []models.Tahap{
		models.TahapInputRegister, models.TahapOutstanding,
		models.TahapMedicalCheck, models.TahapClaimAnalis,
	}
	for _, tahap := range semua {
		if balik := models.TahapDariNama(tahap.String()); balik != tahap {
			t.Errorf("TahapDariNama(%q) = %v, mau %v", tahap.String(), balik, tahap)
		}
	}
	// Nama asing menjadi TIDAK DIKENAL, bukan tebakan: kolom berisi nilai
	// asing adalah kolom yang ditulis di luar aplikasi ini.
	for _, asing := range []string{"", "Outstanding", "ReasLifeAdmin", "Assignment1"} {
		if got := models.TahapDariNama(asing); got.Diketahui() {
			t.Errorf("TahapDariNama(%q) = %v; itu bukan nama tahap", asing, got)
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
	if err := svc.Tahap().Pindah(ctx, inti.Pelaku{}, "CLM-1",
		models.TahapOutstanding, saatUji); !errors.Is(
		err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: galat = %v, mau ErrTanpaIdentitas", err)
	}
	if err := svc.Tahap().Pindah(ctx, pelakuBerperan(inti.PeranAdmin),
		"CLM-1", models.TahapTidakDikenal, saatUji); !errors.Is(
		err, services.ErrTahapTidakDikenal) {
		t.Errorf("tujuan asing: galat = %v, mau ErrTahapTidakDikenal", err)
	}
	if err := svc.Tahap().Pindah(ctx, pelakuBerperan(inti.PeranAdmin),
		"  ", models.TahapOutstanding, saatUji); !errors.Is(
		err, galat.ErrPermintaanTidakSah) {
		t.Errorf("pengenal kosong: galat = %v, mau ErrPermintaanTidakSah", err)
	}
	// ⚠️ Gerbang perannya menuntut pembacaan PY_POSITION, jadi tanpa Oracle ia
	// berhenti di sana. Peran per tahap sendiri diuji langsung di atas.
	if err := svc.Tahap().Pindah(ctx, pelakuBerperan(inti.PeranAdmin),
		"CLM-1", models.TahapMedicalCheck, saatUji); !errors.Is(
		err, db.ErrTanpaOracle) {
		t.Errorf("tanpa Oracle: galat = %v, mau ErrTanpaOracle", err)
	}
}
