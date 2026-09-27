package services_test

// Penegakan peran - TANPA Oracle.
//
// Pemilik: tiket 07. Dibaca sesudah: wewenang.go.

import (
	"context"
	"errors"
	"testing"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
)

func pelakuBerperan(peran ...string) services.Pelaku {
	return services.Pelaku{AkunID: "UJI-AKUN", Peran: peran}
}

// TestMedicalAdvisorTidakDapatMengubahStatus - AC 9 spec.
//
// `[terverifikasi]` `Flow/Register_Flow.xml`: `ReasLifeMedicalAdvisor`
// memegang `Assignment3` (Medical Check) dan `Decision3`; nol gerbang status
// di `AdjustmentDetail_Section` menyebutnya. Tahap telaah medis tidak
// memutuskan akseptasi - ia menelaah.
func TestMedicalAdvisorTidakDapatMengubahStatus(t *testing.T) {
	err := services.WajibPeranPengubahStatus(
		pelakuBerperan(services.PeranMedicalAdvisor), models.StatusDitolak)
	if !errors.Is(err, services.ErrTanpaWewenang) {
		t.Fatalf("galat = %v, mau ErrTanpaWewenang", err)
	}
	// Hanya Admin yang boleh MENOLAK - gerbang XML baris 15399.
	if err := services.WajibPeranPengubahStatus(
		pelakuBerperan(services.PeranRejectOutstanding), models.StatusDitolak); err != nil {
		t.Errorf("Admin menolak baris ditolak: %v", err)
	}
	// ⛔ SPV TIDAK boleh menolak. Ronde pertama memakai satu daftar peran yang
	// datar untuk tujuan apa pun, sehingga SPV dapat menolak lewat Ubah -
	// melewati gerbang Admin di Tolak beserta syarat klaim-bernomornya.
	if err := services.WajibPeranPengubahStatus(
		pelakuBerperan(services.PeranSimpanOutstanding),
		models.StatusDitolak); !errors.Is(err, services.ErrTanpaWewenang) {
		t.Errorf("SPV menolak baris: galat = %v, mau ErrTanpaWewenang", err)
	}
	// ⛔ RALAT audit A0. Kasus ini semula menuntut Aksep DITOLAK bagi siapa
	// pun di modul ini, dengan alasan "jalurnya lewat Komite". Menurut XML itu
	// keliru: `SaveAdjustment_Act` (pecahan 1833, 1879, 1899) mengaksep DI
	// CLAIM LIFE SENDIRI.
	//
	// Kini Aksep lolos lapisan peran, dan yang menggerbanginya PEMEGANG TAHAP
	// beserta prasyarat XML - diperiksa `SimpanAdjustment`, sebab keduanya
	// menuntut keadaan yang hanya terbaca dari basis data.
	for _, peran := range []string{
		services.PeranSimpanOutstanding, services.PeranRejectOutstanding,
	} {
		if err := services.WajibPeranPengubahStatus(
			pelakuBerperan(peran), models.StatusAksep); err != nil {
			t.Errorf("peran %q menuju Aksep ditolak di lapisan peran: %v; "+
				"gerbangnya kini pemegang tahap, bukan daftar peran", peran, err)
		}
	}
}

// TestWewenangKirimKomitePerType mengunci gerbang XML
// `pyWorkPage.pyPosition =='ReasLifeSPV' || pyWorkPage.Type = 'TP' ||
// pyWorkPage.Type = 'TR'` (`AdjustmentDetail_Section.xml` baris 16091).
//
// ⚠️ Bacaan harfiahnya: untuk TP/TR gerbangnya terbuka TANPA memeriksa peran
// sama sekali. Tiket menuliskannya sebagai "ReasLifeAdmin dapat mengirim";
// keduanya sejalan selama Admin memang yang memegang tahapnya, dan selisihnya
// dicatat di tiket.
func TestWewenangKirimKomitePerType(t *testing.T) {
	kasus := []struct {
		tipe  string
		peran []string
		boleh bool
		apa   string
	}{
		{services.TypeQP, []string{services.PeranSimpanOutstanding}, true, "QP oleh SPV"},
		{services.TypeQR, []string{services.PeranSimpanOutstanding}, true, "QR oleh SPV"},
		{services.TypeQP, []string{services.PeranRejectOutstanding}, false, "QP oleh Admin"},
		{services.TypeQR, []string{services.PeranMedicalAdvisor}, false, "QR oleh Medical"},
		{services.TypeTP, []string{services.PeranRejectOutstanding}, true, "TP oleh Admin"},
		{services.TypeTR, []string{services.PeranRejectOutstanding}, true, "TR oleh Admin"},
		{services.TypeTP, []string{services.PeranSimpanOutstanding}, true, "TP oleh SPV"},
	}
	for _, k := range kasus {
		err := services.WajibWewenangKomite(pelakuBerperan(k.peran...), k.tipe)
		if k.boleh && err != nil {
			t.Errorf("%s ditolak: %v", k.apa, err)
		}
		if !k.boleh && !errors.Is(err, services.ErrTanpaWewenang) {
			t.Errorf("%s: galat = %v, mau ErrTanpaWewenang", k.apa, err)
		}
	}
}

// TestKirimKomiteMenuntutTypeDikenal - Type asing bukan izin lewat.
func TestKirimKomiteMenuntutTypeDikenal(t *testing.T) {
	for _, tipe := range []string{"", "XX", "tp"} {
		err := services.WajibWewenangKomite(
			pelakuBerperan(services.PeranSimpanOutstanding), tipe)
		if !errors.Is(err, services.ErrTypeTidakDikenal) {
			t.Errorf("Type %q: galat = %v, mau ErrTypeTidakDikenal", tipe, err)
		}
	}
}

// TestWewenangDitegakkanDiLayanan - AC: penolakan terjadi di lapisan layanan,
// dan tetap terjadi meskipun kontrol UI-nya disembunyikan. Test ini memanggil
// layanan LANGSUNG, tanpa melewati satu pun kontrol layar.
func TestWewenangDitegakkanDiLayanan(t *testing.T) {
	svc := services.New(nil)
	ctx := context.Background()
	// Medical Advisor menolak baris: ditolak di layanan, bukan di layar.
	err := svc.Status().Tolak(ctx, pelakuBerperan(services.PeranMedicalAdvisor),
		"CLM-1", "A-1", saatUji)
	if !errors.Is(err, services.ErrTanpaWewenang) {
		t.Errorf("Medical menolak baris: galat = %v, mau ErrTanpaWewenang", err)
	}
	// Medical Advisor mengubah status: sama.
	err = svc.Status().Ubah(ctx, pelakuBerperan(services.PeranMedicalAdvisor),
		"CLM-1", "P-1", "A-1", models.StatusDitolak, saatUji)
	if !errors.Is(err, services.ErrTanpaWewenang) {
		t.Errorf("Medical mengubah status: galat = %v, mau ErrTanpaWewenang", err)
	}
}
