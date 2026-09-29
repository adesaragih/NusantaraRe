package handlers

// Uji pintu perpindahan tahap - butir aw.

import (
	"os"
	"strings"
	"testing"

	"nusantarare/modul/claimlife/models"
)

func TestTahapTujuanHimpunanTertutup(t *testing.T) {
	sah := map[string]models.Tahap{
		"input-register": models.TahapInputRegister,
		"outstanding":    models.TahapOutstanding,
		"medical-check":  models.TahapMedicalCheck,
		"claim-analis":   models.TahapClaimAnalis,
	}
	for potongan, mau := range sah {
		got, ok := tahapTujuan(potongan)
		if !ok || got != mau {
			t.Errorf("tahapTujuan(%q) = (%v, %v), mau (%v, true)", potongan, got, ok, mau)
		}
	}
	// ⛔ KATA, bukan angka. Jalur `/tahap/2` tidak dapat dibaca siapa pun, dan
	// angka yang bergeser bila urutan `models.Tahap` berubah akan memindahkan
	// kasus ke tempat yang salah tanpa satu pun galat.
	for _, asing := range []string{"1", "2", "", "Outstanding", "ReasLifeAdmin", "medical"} {
		if _, ok := tahapTujuan(asing); ok {
			t.Errorf("tahapTujuan(%q) diterima; mau ditolak", asing)
		}
	}
}

// Butir aw: kedua tombol Outstanding memakai rute ini, dan sebabnya dicatat.
func TestRuteTahapTerdaftarDenganBuktinya(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	teks := string(isi)
	if !strings.Contains(teks, `"POST /api/klaim-life/{id}/tahap/{tujuan}"`) {
		t.Error("rute perpindahan tahap tidak terdaftar")
	}
	// ⚠️ Komentar rutenya WAJIB menyebut kedua tombol beserta barisnya.
	// Ronde pertama komentarnya terpotong substitusi shell dan menjadi
	// "( 21404,  21839)" - kalimat yang tidak menjelaskan apa pun.
	for _, bukti := range []string{"Send Back to Register", "Send to Medical Check", "21404", "21839"} {
		if !strings.Contains(teks, bukti) {
			t.Errorf("komentar rute tahap tidak menyebut %q", bukti)
		}
	}
}
