package models_test

import (
	"testing"

	"nusantarare/modul/premiumlistlife/backend/models"
)

// TestAkibatStatusOffer - Save Offer Bind/Closed menutup kasus bendera "0"
// (keputusan work owner 08-10-2026); status lain dan bendera "1" tidak.
func TestAkibatStatusOffer(t *testing.T) {
	a, keputusan, sebab, tutup, err := models.AkibatStatusOffer(models.FlagPolisPenawaran, models.StatusOfferBind)
	if err != nil || !tutup || keputusan != models.KeputusanConfirm || a.StatusWork != models.StatusPolisSelesai ||
		a.SimpanPolis || a.TahapTujuan != "" || a.KeDecision3 {
		t.Fatalf("Bind: %+v %q %v %v", a, keputusan, tutup, err)
	}
	if sebab != "Save Offer Bind -> Confirm -> Offer" {
		t.Errorf("sebab Bind %q", sebab)
	}
	// Sama persis dengan Confirm tahap penawaran yang dirutekan Decision3 bendera "0".
	lanjut, _ := models.PenggolongOtomatis(models.FlagPolisPenawaran)
	if a != lanjut {
		t.Errorf("Bind %+v != jalur Offer %+v", a, lanjut)
	}

	a, keputusan, sebab, tutup, err = models.AkibatStatusOffer(models.FlagPolisPenawaran, models.StatusOfferClosed)
	if err != nil || !tutup || keputusan != models.KeputusanDecline || a.StatusWork != models.StatusPolisDitolak || a.SimpanPolis {
		t.Fatalf("Closed: %+v %q %v %v", a, keputusan, tutup, err)
	}
	if decline, _ := models.TransisiPenawaran(models.TahapPolisPenawaran, models.KeputusanDecline); a != decline {
		t.Errorf("Closed %+v != Decline %+v", a, decline)
	}
	if sebab != "Save Offer Closed -> Decline" {
		t.Errorf("sebab Closed %q", sebab)
	}

	for _, s := range []string{"Accept", "Pending", "Decline", ""} {
		if _, _, _, tutup, err := models.AkibatStatusOffer(models.FlagPolisPenawaran, s); tutup || err != nil {
			t.Errorf("%q menutup (%v)", s, err)
		}
	}
	// Bendera "1" (EmailTypePL) tidak pernah menutup, apa pun kodenya.
	for _, s := range []string{"1", "2", "3", models.StatusOfferBind, models.StatusOfferClosed} {
		if _, _, _, tutup, err := models.AkibatStatusOffer(models.FlagPolisPremium, s); tutup || err != nil {
			t.Errorf("bendera 1 %q menutup (%v)", s, err)
		}
	}
	// Kode Bind/Closed ada di radio Status bendera "0".
	ada := map[string]bool{}
	for _, p := range models.PilihanStatusOffer {
		ada[p.Kode] = true
	}
	if !ada[models.StatusOfferBind] || !ada[models.StatusOfferClosed] {
		t.Error("Bind/Closed tidak ada di PilihanStatusOffer")
	}
}

// TestConfirmPenawaranTersedia - Confirm tahap penawaran hanya bendera "1"
// (keputusan work owner 08-10-2026); tahap detail dan keputusan lain tetap.
func TestConfirmPenawaranTersedia(t *testing.T) {
	for _, k := range []struct {
		tahap, keputusan, flag string
		mau                    bool
	}{
		{models.TahapPolisPenawaran, models.KeputusanConfirm, models.FlagPolisPenawaran, false},
		{models.TahapPolisPenawaran, models.KeputusanConfirm, models.FlagPolisPremium, true},
		{models.TahapPolisPenawaran, models.KeputusanDecline, models.FlagPolisPenawaran, true},
		{models.TahapPolisDetail, models.KeputusanConfirm, models.FlagPolisPenawaran, true},
		{models.TahapPolisDetail, models.KeputusanConfirm, models.FlagPolisPremium, true},
	} {
		if got := models.ConfirmPenawaranTersedia(k.tahap, k.keputusan, k.flag); got != k.mau {
			t.Errorf("%q %q bendera %q = %v", k.tahap, k.keputusan, k.flag, got)
		}
	}
	// Save Offer Bind tetap menutup Resolved-Completed (jalur internal, bukan Confirm klien).
	a, _, _, tutup, err := models.AkibatStatusOffer(models.FlagPolisPenawaran, models.StatusOfferBind)
	if err != nil || !tutup || a.StatusWork != models.StatusPolisSelesai {
		t.Errorf("Bind %+v %v %v", a, tutup, err)
	}
	if models.PesanConfirmPenawaranTidakTersedia != "Confirm is not available for offer cases. Choose Status Bind and press Save Offer." {
		t.Error("pesan berubah")
	}
}
