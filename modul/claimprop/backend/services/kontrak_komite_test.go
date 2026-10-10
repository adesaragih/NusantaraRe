package services_test

// Uji kontrak `kontrak.KlaimTreatyKomite` sisi penyedia (Claim Prop) di atas tiruan Claim Prop. Sisi pemakai (Komite
// Claim Prop) diuji dengan kontrak palsu di modulnya sendiri.

import (
	"context"
	"errors"
	"testing"
	"time"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/services"
	"nusantarare/modul/claimprop/backend/tiruan"
)

const (
	klaimUjiKomite = "CLMP-UJI001"
	adjUjiKomite   = "UJI-ADJ-1"
)

func siapKontrakKomite(t *testing.T) (*tiruan.Gudang, services.KlaimUntukKomite) {
	t.Helper()
	g, a := tiruan.Baru(), tiruan.AcuanBaru()
	a.OldID["UJI-COB"] = "12"
	saat := time.Date(2026, 10, 8, 10, 0, 0, 0, models.Jakarta)
	if err := g.SisipKasus(context.Background(), nil, klaimUjiKomite, "UJI-ADMIN", "UJI Admin", saat); err != nil {
		t.Fatal(err)
	}
	h := models.HalamanBaru()
	h.Setel(models.CD+"NoClaim", "UJI-K-0001")
	h.Setel(models.CD+"QuotationData.BusinessCode", "UJI-COB")
	h.SetelDaftar(models.DaftarAdjustment, []models.Baris{
		{models.PropID: "UJI-ADJ-0", "Type": "1", "AcceptanceStatus": "1"},
		{models.PropID: adjUjiKomite, "Type": "2", "IsKomite": "1"},
	})
	g.SetelHalaman(klaimUjiKomite, h)
	l := services.Baru(g, a, func() time.Time { return saat }, false)
	return g, l.KlaimUntukKomite()
}

func TestKontrakKomiteBacaKlaimTreaty(t *testing.T) {
	_, k := siapKontrakKomite(t)
	kt, err := k.BacaKlaimTreaty(context.Background(), nil, klaimUjiKomite, adjUjiKomite)
	if err != nil {
		t.Fatal(err)
	}
	if kt.Adjustment != 2 {
		t.Fatalf("posisi adjustment %d, mau 2", kt.Adjustment)
	}
	if kt.Nilai[models.CD+"NoClaim"] != "UJI-K-0001" || kt.Nilai[models.OQ+"BusinessOldId"] != "12" {
		t.Fatalf("nilai klaim: NoClaim %q, BusinessOldId %q (mau 12 dari BUSINESS.OLDID)",
			kt.Nilai[models.CD+"NoClaim"], kt.Nilai[models.OQ+"BusinessOldId"])
	}
	if rows := kt.Daftar[models.DaftarAdjustment]; len(rows) != 2 || rows[1]["Type"] != "2" {
		t.Fatalf("daftar adjustment %+v", rows)
	}
	if kt.Tertutup {
		t.Fatal("kasus terbuka terbaca tertutup")
	}
	if _, err := k.BacaKlaimTreaty(context.Background(), nil, klaimUjiKomite, "UJI-ADJ-X"); !errors.Is(err,
		kontrak.ErrAdjustmentTreatyTidakAda) {
		t.Fatalf("adjustment tak dikenal: %v", err)
	}
	if _, err := k.BacaKlaimTreaty(context.Background(), nil, "CLMP-TIDAKADA", adjUjiKomite); !errors.Is(err,
		kontrak.ErrKlaimTreatyTidakAda) {
		t.Fatalf("klaim tak dikenal: %v", err)
	}
}

func TestKontrakKomiteKunciMenolakKlaimTertutup(t *testing.T) {
	g, k := siapKontrakKomite(t)
	if err := k.KunciKlaimTreaty(context.Background(), nil, klaimUjiKomite); err != nil {
		t.Fatalf("klaim terbuka: %v", err)
	}
	kasus := g.Kasus[klaimUjiKomite]
	kasus.StatusWork = models.StatusSelesai
	g.Kasus[klaimUjiKomite] = kasus
	if err := k.KunciKlaimTreaty(context.Background(), nil, klaimUjiKomite); !errors.Is(err,
		kontrak.ErrKlaimTreatyTertutup) {
		t.Fatalf("klaim tertutup: %v", err)
	}
}

func TestKontrakKomiteTulisBalik(t *testing.T) {
	g, k := siapKontrakKomite(t)
	saat := time.Date(2026, 10, 8, 11, 0, 0, 0, models.Jakarta)
	u := kontrak.UbahanKlaimTreaty{
		Header:     map[string]string{"IsAnyAcceptation": "1", "ClaimData.IsCloseFile": "1"},
		Adjustment: map[string]string{"AcceptedNo": "UJI-A12.10.2026.TP00001", "AcceptanceStatus": "1"},
		Riwayat:    []kontrak.RiwayatKlaimTreaty{{Teks: "Accepted by UJI-JABATAN", Pelaku: "UJI-K1", Tingkat: "UJI-JABATAN", Saat: saat}},
		FacRetro:   []map[string]string{{"ReinsurerID": "UJI-RE1", "ReinsurerName": "UJI-REAS", "PctShareAllObj": "40"}},
	}
	if err := k.TulisBalikKlaimTreaty(context.Background(), nil, klaimUjiKomite, adjUjiKomite, u); err != nil {
		t.Fatal(err)
	}
	h := g.Halaman(klaimUjiKomite)
	if h.Ambil("IsAnyAcceptation") != "1" || h.Ambil(models.CD+"IsCloseFile") != "1" {
		t.Fatalf("header tidak ditulis: %+v", h.Nilai)
	}
	adj := h.AmbilDaftar(models.DaftarAdjustment)
	if adj[1]["AcceptedNo"] != "UJI-A12.10.2026.TP00001" || adj[0]["AcceptedNo"] != "" {
		t.Fatalf("adjustment: %+v", adj)
	}
	riw := h.AmbilDaftar(models.DaftarRiwayat)
	if len(riw) != 1 || riw[0]["CommentSuggest"] != "Accepted by UJI-JABATAN" || riw[0]["IsCedingConfirm"] != "UJI-JABATAN" ||
		riw[0]["PICSuggest"] != "UJI-K1" {
		t.Fatalf("riwayat: %+v", riw)
	}
	if fr := h.AmbilDaftar(models.DaftarFacRetro); len(fr) != 1 || fr[0]["ReinsurerID"] != "UJI-RE1" {
		t.Fatalf("FacRetroList: %+v", fr)
	}
	// FacRetroList yang sudah berisi tidak ditimpa (SaveAcceptationTreaty_TKMT S6.3 hanya bila SizeRetro < 1)
	u2 := kontrak.UbahanKlaimTreaty{FacRetro: []map[string]string{{"ReinsurerID": "UJI-RE2"}}}
	if err := k.TulisBalikKlaimTreaty(context.Background(), nil, klaimUjiKomite, adjUjiKomite, u2); err != nil {
		t.Fatal(err)
	}
	if fr := g.Halaman(klaimUjiKomite).AmbilDaftar(models.DaftarFacRetro); len(fr) != 1 || fr[0]["ReinsurerID"] != "UJI-RE1" {
		t.Fatalf("FacRetroList tertimpa: %+v", fr)
	}
}

func TestKontrakKomiteTulisBalikDaftarPutih(t *testing.T) {
	g, k := siapKontrakKomite(t)
	for _, u := range []kontrak.UbahanKlaimTreaty{
		{Header: map[string]string{"ClaimData.NoClaim": "UJI-LAIN"}},
		{Adjustment: map[string]string{"AdjustmentValue": "1"}},
		{FacRetro: []map[string]string{{"BukanKolom": "1"}}},
	} {
		if err := k.TulisBalikKlaimTreaty(context.Background(), nil, klaimUjiKomite, adjUjiKomite, u); !errors.Is(err,
			kontrak.ErrUbahanKlaimTreatyTidakSah) {
			t.Fatalf("ubahan %+v: %v", u, err)
		}
	}
	if g.Halaman(klaimUjiKomite).Ambil(models.CD+"NoClaim") != "UJI-K-0001" {
		t.Fatal("ubahan di luar daftar putih tertulis")
	}
	kasus := g.Kasus[klaimUjiKomite]
	kasus.StatusWork = models.StatusSelesai
	g.Kasus[klaimUjiKomite] = kasus
	if err := k.TulisBalikKlaimTreaty(context.Background(), nil, klaimUjiKomite, adjUjiKomite,
		kontrak.UbahanKlaimTreaty{Header: map[string]string{"IsAnyAcceptation": "1"}}); !errors.Is(err,
		kontrak.ErrKlaimTreatyTertutup) {
		t.Fatalf("klaim tertutup: %v", err)
	}
}
