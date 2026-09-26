package services_test

// Mesin status per baris - TANPA Oracle.
//
// Pemilik: tiket 04. Dibaca sesudah: statusbaris.go.

import (
	"context"
	"errors"
	"testing"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/services"
)

// saatUji adalah jam yang disuntikkan, bukan jam sungguhan: fungsi yang
// membaca jamnya sendiri tidak dapat diuji tanpa menunggu waktu berlalu.
var saatUji = time.Date(2026, 9, 26, 21, 30, 0, 0, time.UTC)

func baris(id, kode string) models.BarisAdjustment {
	return models.BarisAdjustment{ID: id, KodeStatus: kode}
}

func klaimDengan(b ...models.BarisAdjustment) models.Klaim {
	return models.Klaim{Peserta: []models.Peserta{{ID: "UJI-P-1", Baris: b}}}
}

// TestBarisFinalTidakDapatBerubah - AC 2: Aksep dan Ditolak keduanya final,
// lewat jalur mana pun. `[terverifikasi]` sensus OQ-061: nol rule menulis "0"
// sesudah "1" atau "2".
func TestBarisFinalTidakDapatBerubah(t *testing.T) {
	for _, dari := range []string{models.KodeAksep, models.KodeDitolak} {
		for _, ke := range []models.StatusBaris{
			models.StatusAksep, models.StatusDitolak, models.StatusOutstanding,
		} {
			_, err := services.Transisi(baris("UJI-A", dari), ke, saatUji)
			if !errors.Is(err, services.ErrBarisSudahFinal) {
				t.Errorf("dari %q ke %v: galat = %v, mau ErrBarisSudahFinal", dari, ke, err)
			}
		}
	}
}

// TestTransisiHanyaDariOutstanding - baris tanpa status belum pernah disimpan
// ke Outstanding, jadi ia belum punya apa pun untuk ditransisikan.
func TestTransisiHanyaDariOutstanding(t *testing.T) {
	_, err := services.Transisi(baris("UJI-A", ""), models.StatusAksep, saatUji)
	if !errors.Is(err, services.ErrTransisiTidakSah) {
		t.Fatalf("galat = %v, mau ErrTransisiTidakSah", err)
	}
	hasil, err := services.Transisi(baris("UJI-A", models.KodeOutstanding), models.StatusAksep, saatUji)
	if err != nil {
		t.Fatalf("Outstanding -> Aksep: %v", err)
	}
	if hasil.KodeStatus != models.KodeAksep {
		t.Errorf("kode = %q, mau %q", hasil.KodeStatus, models.KodeAksep)
	}
}

// TestTujuanTransisiHanyaAksepAtauDitolak - Outstanding bukan tujuan (ia
// keadaan awal), dan TidakDiketahui tidak pernah ditulis. Kode "4" ada di data
// warisan, artinya belum diputuskan work owner, dan sistem baru tidak
// menulisnya.
func TestTujuanTransisiHanyaAksepAtauDitolak(t *testing.T) {
	for _, ke := range []models.StatusBaris{
		models.StatusOutstanding, models.StatusTidakDiketahui,
	} {
		_, err := services.Transisi(baris("UJI-A", models.KodeOutstanding), ke, saatUji)
		if !errors.Is(err, services.ErrTransisiTidakSah) {
			t.Errorf("tujuan %v: galat = %v, mau ErrTransisiTidakSah", ke, err)
		}
	}
}

// TestStatusKlaimTurunan - AC 8 dan kerabatnya.
func TestStatusKlaimTurunan(t *testing.T) {
	kasus := []struct {
		nama string
		k    models.Klaim
		mau  models.StatusKlaim
	}{
		{"nol baris", klaimDengan(), models.KlaimBelumBerbaris},
		{"satu Outstanding", klaimDengan(baris("A", models.KodeOutstanding)),
			models.KlaimBerjalan},
		{"Aksep dan Outstanding", klaimDengan(
			baris("A", models.KodeAksep), baris("B", models.KodeOutstanding)),
			models.KlaimBerjalan},
		{"satu Aksep", klaimDengan(baris("A", models.KodeAksep)), models.KlaimSelesai},
		{"Aksep dan Ditolak", klaimDengan(
			baris("A", models.KodeDitolak), baris("B", models.KodeAksep)),
			models.KlaimSelesai},
		{"seluruhnya Ditolak", klaimDengan(
			baris("A", models.KodeDitolak), baris("B", models.KodeDitolak)),
			models.KlaimDitolakSeluruhnya},
		// ⛔ Kode yang tidak dikenal bukan Outstanding, bukan Aksep, dan bukan
		// pula Ditolak. Melaporkannya "ditolak seluruhnya" membuat klaim
		// tampak SUDAH diputus padahal keputusannya justru yang tidak terbaca.
		{"hanya kode tak dikenal", klaimDengan(baris("A", "4")),
			models.KlaimTidakDapatDipastikan},
		{"tak dikenal bersama Outstanding", klaimDengan(
			baris("A", "4"), baris("B", models.KodeOutstanding)),
			models.KlaimBerjalan},
		{"tak dikenal bersama Aksep", klaimDengan(
			baris("A", "4"), baris("B", models.KodeAksep)),
			models.KlaimSelesai},
	}
	for _, k := range kasus {
		if got := k.k.StatusTurunan(); got != k.mau {
			t.Errorf("%s: StatusTurunan = %v, mau %v", k.nama, got, k.mau)
		}
	}
}

// TestMenolakSatuBarisTidakMenutupKlaim - AC 4 dan AC "Ditolak selalu berarti
// baris itu ditolak". Sesudah satu baris ditolak, klaim masih menerima baris
// baru dan kembali berjalan.
func TestMenolakSatuBarisTidakMenutupKlaim(t *testing.T) {
	k := klaimDengan(baris("A", models.KodeDitolak))
	if got := k.StatusTurunan(); got == models.KlaimSelesai {
		t.Fatal("klaim yang barisnya ditolak dilaporkan SELESAI; " +
			"Ditolak berarti baris itu ditolak, tidak pernah berarti klaim selesai")
	}
	if err := services.TambahBaris(&k.Peserta[0], models.BarisAdjustment{}); err != nil {
		t.Fatalf("TambahBaris: %v", err)
	}
	if n := services.TandaiOutstandingKlaim(&k); n != 1 {
		t.Errorf("baris tersentuh = %d, mau 1; hanya baris BARU yang menjadi "+
			"Outstanding, yang sudah final tidak disentuh ulang", n)
	}
	if k.Peserta[0].Baris[0].KodeStatus != models.KodeDitolak {
		t.Errorf("baris yang sudah Ditolak dikembalikan ke Outstanding (%q); "+
			"kefinalan berlaku pada jalur ini juga", k.Peserta[0].Baris[0].KodeStatus)
	}
	if got := k.StatusTurunan(); got != models.KlaimBerjalan {
		t.Errorf("sesudah baris baru: %v, mau KlaimBerjalan", got)
	}
}

// TestBarisTerakhirAdalahSumberPencerminan - `[terverifikasi]`
// `serviceInsertArasapasClaimLife_act` langkah 1.1.1 menyalin `.STS_REJECT`
// dan `.ACCEPTEDNO` baris adjustment ke header di dalam putaran BERSARANG
// tanpa henti, sehingga baris yang terakhir diulang yang menang.
func TestBarisTerakhirAdalahSumberPencerminan(t *testing.T) {
	k := models.Klaim{Peserta: []models.Peserta{
		{ID: "UJI-P-1", Baris: []models.BarisAdjustment{
			baris("A", models.KodeAksep), baris("B", models.KodeDitolak)}},
		{ID: "UJI-P-2", Baris: []models.BarisAdjustment{
			baris("C", models.KodeOutstanding)}},
	}}
	akhir := services.BarisTerakhir(&k)
	if akhir == nil {
		t.Fatal("BarisTerakhir nil padahal ada tiga baris")
	}
	if akhir.ID != "C" {
		t.Errorf("baris terakhir = %q, mau C (peserta terakhir, baris terakhir)", akhir.ID)
	}
	if services.BarisTerakhir(&models.Klaim{}) != nil {
		t.Error("klaim tanpa baris mengembalikan baris")
	}
}

// TestUbahStatusMenjagaPagarnya - pelaku anonim, pengenal kosong, dan tanpa
// Oracle masing-masing gagal terang.
func TestUbahStatusMenjagaPagarnya(t *testing.T) {
	svc := services.New(nil)
	// Berperan sah: yang diuji di sini pagar PENGENAL dan Oracle, bukan
	// pagar peran - yang punya testnya sendiri di wewenang_test.go.
	pelaku := services.Pelaku{AkunID: "UJI-AKUN",
		Peran: []string{services.PeranRejectOutstanding}}
	if err := svc.Status().Ubah(context.Background(), services.Pelaku{},
		"CLM-1", "P-1", "A-1", models.StatusDitolak, saatUji); !errors.Is(
		err, services.ErrTanpaIdentitas) {
		t.Errorf("pelaku anonim: galat = %v, mau ErrTanpaIdentitas", err)
	}
	if err := svc.Status().Ubah(context.Background(), pelaku,
		"CLM-1", "", "A-1", models.StatusDitolak, saatUji); !errors.Is(
		err, services.ErrPermintaanTidakSah) {
		t.Errorf("peserta kosong: galat = %v, mau ErrPermintaanTidakSah", err)
	}
	if err := svc.Status().Ubah(context.Background(), pelaku,
		"CLM-1", "P-1", "A-1", models.StatusDitolak, saatUji); !errors.Is(
		err, repository.ErrTanpaOracle) {
		t.Errorf("tanpa Oracle: galat = %v, mau ErrTanpaOracle", err)
	}
}

// TestAksepMenstempelTanggalAkseptasi - `[terverifikasi]` SaveAdjustment_Act
// menulis STS_REJECT, ACCEPTEDNO, dan ACCEPTATION_DATE dalam SATU
// Property-Set. Ditolak tidak menstempelnya: tanggal akseptasi milik aksi
// akseptasi, bukan milik setiap keputusan.
func TestAksepMenstempelTanggalAkseptasi(t *testing.T) {
	aksep, err := services.Transisi(
		baris("A", models.KodeOutstanding), models.StatusAksep, saatUji)
	if err != nil {
		t.Fatalf("Aksep: %v", err)
	}
	if !aksep.TanggalAkseptasi.Equal(saatUji) {
		t.Errorf("TanggalAkseptasi = %v, mau %v", aksep.TanggalAkseptasi, saatUji)
	}
	tolak, err := services.Transisi(
		baris("B", models.KodeOutstanding), models.StatusDitolak, saatUji)
	if err != nil {
		t.Fatalf("Ditolak: %v", err)
	}
	if !tolak.TanggalAkseptasi.IsZero() {
		t.Errorf("baris yang DITOLAK ikut distempel tanggal akseptasi: %v",
			tolak.TanggalAkseptasi)
	}
}
