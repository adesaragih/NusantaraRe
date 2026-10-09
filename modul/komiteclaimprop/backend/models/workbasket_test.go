package models_test

import (
	"testing"

	"nusantarare/modul/komiteclaimprop/backend/models"
)

// Keputusan work owner 09-10-2026: KomiteID roster komite PROP = nama workbasket (migrasi claimprop 537); penyetuju =
// anggota workbasket tingkat berjalan; satu orang tidak menyetujui dua tingkat kasus yang sama; T_WORK_CLAIM.POSITION =
// workbasket tingkat berjalan.

// kasusWB - dua tingkat ber-KomiteID workbasket UJI-WB-1 / UJI-WB-2.
func kasusWB(count int, putusan ...string) models.Kasus {
	k := kasusUji(count, putusan...)
	k.Tangga[0].OperatorID, k.Tangga[1].OperatorID = "UJI-WB-1", "UJI-WB-2"
	return k
}

func TestPemegangAnggotaWorkbasketTingkatBerjalan(t *testing.T) {
	k := kasusWB(1)
	if !k.Pemegang("UJI-A", []string{"UJI-LAIN", "UJI-WB-1"}) {
		t.Fatal("anggota workbasket tingkat 1 memegang assignment")
	}
	if k.Pemegang("UJI-A", []string{"UJI-WB-2"}) || k.Pemegang("UJI-A", nil) {
		t.Fatal("bukan anggota workbasket tingkat berjalan: bukan pemegang")
	}
	if !kasusUji(1).Pemegang("UJI-K1", nil) {
		t.Fatal("tangga ber-akun (sebelum 537) tetap dipegang akunnya")
	}
}

func TestSatuOrangTidakMenyetujuiDuaTingkat(t *testing.T) {
	k := kasusWB(2, models.KeputusanSetuju)
	k.Tangga[0].OperatorID = "UJI-A" // baris yang diputuskan menyimpan akun pemutusnya
	if k.Pemegang("UJI-A", []string{"UJI-WB-1", "UJI-WB-2"}) {
		t.Fatal("UJI-A sudah memutus tingkat 1: tidak boleh memegang tingkat 2 kasus yang sama")
	}
	if !k.Pemegang("UJI-B", []string{"UJI-WB-2"}) {
		t.Fatal("anggota lain workbasket tingkat 2 tetap memegang")
	}
}

func TestRencanaMenimpaPemutusDanMemajukanPosisi(t *testing.T) {
	k := kasusWB(1)
	r := models.Rencanakan(k, klaimUji(""), models.Keputusan{AcceptStatus: models.KeputusanSetuju, Comment: "UJI"},
		"UJI-A", saatUji)
	if len(r.Tangga) == 0 || r.Tangga[0].ID != "L1" || r.Tangga[0].Pemutus != "UJI-A" {
		t.Fatalf("baris yang diputuskan menyimpan akun pemutus: %+v", r.Tangga)
	}
	if r.Selesai || r.Posisi != "UJI-WB-2" {
		t.Fatalf("setuju tingkat 1: POSITION %q (selesai %v), mau UJI-WB-2", r.Posisi, r.Selesai)
	}

	r = models.Rencanakan(kasusWB(2, models.KeputusanSetuju), klaimUji(""),
		models.Keputusan{AcceptStatus: models.KeputusanSetuju, Comment: "UJI"}, "UJI-B", saatUji)
	if !r.Selesai || r.Posisi != "" {
		t.Fatalf("setuju tingkat akhir: POSITION %q (selesai %v), mau kosong", r.Posisi, r.Selesai)
	}

	r = models.Rencanakan(kasusWB(1), klaimUji(""), models.Keputusan{AcceptStatus: models.KeputusanTolak,
		Comment: "UJI"}, "UJI-A", saatUji)
	if !r.Selesai || r.Posisi != "" {
		t.Fatalf("tolak: POSITION %q (selesai %v), mau kosong", r.Posisi, r.Selesai)
	}
	for _, u := range r.Tangga[1:] {
		if u.Pemutus != "" {
			t.Fatalf("baris yang ditolak otomatis (S26.1) bukan diputus pelaku: %+v", u)
		}
	}
}
