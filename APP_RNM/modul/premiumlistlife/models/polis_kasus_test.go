package models_test

// Kasus polis baru - tombol `Input Offer` / `Input Premium` (GILIRAN-13 paket 1, butir bn).

import (
	"errors"
	"os"
	"strings"
	"testing"

	"nusantarare/inti"
	"nusantarare/modul/premiumlistlife/models"
)

func TestKeduaTombolMulaiDiTahapPenawaran(t *testing.T) {
	// ⛔ Flow `InputPolicyHolder.xml`: Start1 -> [Always] -> Assignment2 "Input
	// Offer" untuk KEDUA bendera. Benderanya baru bekerja di Decision3 sesudah
	// Confirm - bukan memilih tahap awal.
	for _, flag := range []string{models.FlagPolisPenawaran, models.FlagPolisPremium} {
		k, err := models.SusunKasusPolisBaru(flag)
		if err != nil {
			t.Fatalf("flag %q: %v", flag, err)
		}
		if k.Status != models.TahapPolisPenawaran || k.Posisi != models.PosisiOffer ||
			k.Lini != inti.LiniLife || k.Flag != flag {
			t.Errorf("flag %q: %+v", flag, k)
		}
	}
}

func TestBenderaDiLuarNolSatuDitolak(t *testing.T) {
	for _, flag := range []string{"", "2", " 0", "true", "Offer"} {
		if _, err := models.SusunKasusPolisBaru(flag); !errors.Is(err, models.ErrFlagPolisTidakSah) {
			t.Errorf("flag %q: %v, mau ErrFlagPolisTidakSah", flag, err)
		}
	}
}

// TestNilaiBenderaVERBATIMDariKorpus membaca section portal LANGSUNG.
func TestNilaiBenderaVERBATIMDariKorpus(t *testing.T) {
	const letak = `D:\XML\RNM_BRD\PremiumList Life\Section\PremiumList.xml`
	isi, err := os.ReadFile(letak)
	if err != nil {
		t.Skipf("korpus tidak terjangkau di mesin ini (%v)", err)
	}
	teks := string(isi)
	for _, v := range []string{models.FlagPolisPenawaran, models.FlagPolisPremium} {
		if !strings.Contains(teks, `<pyValue>"`+v+`"</pyValue>`) {
			t.Errorf("nilai FlagPolicy %q tidak ada di PremiumList.xml", v)
		}
	}
}
