package services_test

// Uji kontrak `kontrak.KlaimTreatyNonPropKomite` sisi penyedia (Claim Non Prop) di atas tiruan Claim Non Prop. Sisi
// pemakai (Komite Claim Non Prop) diuji dengan kontrak palsu di modulnya sendiri; uji ini memastikan bentuk halaman yang
// dibaca pemakai (baris akseptasi, anak `ClaimData.AdjustmentList(n).<anak>`) sama dengan halaman Claim Non Prop.

import (
	"context"
	"errors"
	"testing"
	"time"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimnonprop/backend/models"
	"nusantarare/modul/claimnonprop/backend/services"
	"nusantarare/modul/claimnonprop/backend/tiruan"
)

const (
	klaimUjiKomite = "CLMNP-UJI001"
	adjUjiKomite   = "UJI-ADJ-1"
)

func siapKontrakKomite(t *testing.T) (*tiruan.Gudang, services.KlaimUntukKomite) {
	t.Helper()
	g, a := tiruan.Baru(), tiruan.AcuanBaru()
	a.OldID["UJI COB"] = "123" // GetDataBusiness_SQL ada, tetapi S14.11 membaca OfferFacIn.QuotationData.BusinessOldId
	saat := time.Date(2026, 10, 9, 10, 0, 0, 0, models.Jakarta)
	if err := g.SisipKasus(context.Background(), nil, klaimUjiKomite, "UJI-ADMIN", "UJI Admin", saat); err != nil {
		t.Fatal(err)
	}
	h := models.HalamanBaru()
	h.Setel(models.CD+"NoClaim", "UJI-K-0001")
	h.Setel(models.OQ+"BusinessName", "UJI COB")
	h.SetelDaftar(models.DaftarAdjustment, []models.Baris{
		{models.PropID: "UJI-ADJ-0", "Type": "1", "AcceptanceStatus": "1"},
		{models.PropID: adjUjiKomite, "Type": "1", "IsKomite": "1", "DataCommitteeTreaty.Remarks": "UJI CATATAN"},
	})
	h.SetelDaftar(models.JalurAdj(2, models.AnakXOL), []models.Baris{{"TreatyName": "UJI XOL 1", "TotalClaim": "300"}})
	g.SetelHalaman(klaimUjiKomite, h)
	l := services.Baru(g, a, func() time.Time { return saat }, false)
	return g, l.KlaimUntukKomite()
}

func TestKontrakKomiteNonPropBacaKlaimTreaty(t *testing.T) {
	_, k := siapKontrakKomite(t)
	kt, err := k.BacaKlaimTreaty(context.Background(), nil, klaimUjiKomite, adjUjiKomite)
	if err != nil {
		t.Fatal(err)
	}
	if kt.Adjustment != 2 || kt.Tertutup {
		t.Fatalf("posisi akseptasi %d (mau 2), tertutup %v", kt.Adjustment, kt.Tertutup)
	}
	if kt.Nilai[models.CD+"NoClaim"] != "UJI-K-0001" || kt.Nilai[models.OQ+"BusinessOldId"] != "" {
		t.Fatalf("nilai klaim: NoClaim %q, BusinessOldId %q (mau kosong persis XML, OQ-CNP-37)", kt.Nilai[models.CD+"NoClaim"],
			kt.Nilai[models.OQ+"BusinessOldId"])
	}
	if rows := kt.Daftar[models.DaftarAdjustment]; len(rows) != 2 || rows[1]["DataCommitteeTreaty.Remarks"] != "UJI CATATAN" {
		t.Fatalf("baris akseptasi: %+v", rows)
	}
	if xol := kt.Daftar["ClaimData.AdjustmentList(2).SpreadingRisk"]; len(xol) != 1 || xol[0]["TreatyName"] != "UJI XOL 1" {
		t.Fatalf("anak XOL Allocation (jalur yang dibaca Komite): %+v", xol)
	}
	if _, err := k.BacaKlaimTreaty(context.Background(), nil, klaimUjiKomite, "UJI-ADJ-X"); !errors.Is(err,
		kontrak.ErrAdjustmentTreatyTidakAda) {
		t.Fatalf("akseptasi tak dikenal: %v", err)
	}
	if _, err := k.BacaKlaimTreaty(context.Background(), nil, "CLMNP-TIDAKADA", adjUjiKomite); !errors.Is(err,
		kontrak.ErrKlaimTreatyTidakAda) {
		t.Fatalf("klaim tak dikenal: %v", err)
	}
}

func TestKontrakKomiteNonPropKunciMenolakKlaimTertutup(t *testing.T) {
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

func TestKontrakKomiteNonPropTulisBalikDaftarPutih(t *testing.T) {
	g, k := siapKontrakKomite(t)
	saat := time.Date(2026, 10, 9, 11, 0, 0, 0, models.Jakarta)
	u := kontrak.UbahanKlaimTreaty{
		Header:     map[string]string{"CNPStatusCase": "CLAIM ACCEPTED", "IsCloseFile": "0"},
		Adjustment: map[string]string{"AcceptedNo": "UJIA123.10.2026.TX00001", "AcceptanceStatus": "1", "NoAccount": "123"},
		Riwayat:    []kontrak.RiwayatKlaimTreaty{{Teks: "Accepted by UJI", Pelaku: "UJI-K1", Tingkat: "UJI-JABATAN", Saat: saat}},
	}
	if err := k.TulisBalikKlaimTreaty(context.Background(), nil, klaimUjiKomite, adjUjiKomite, u); err != nil {
		t.Fatal(err)
	}
	h := g.Halaman(klaimUjiKomite)
	if h.Ambil("CNPStatusCase") != "CLAIM ACCEPTED" || h.Ambil("IsCloseFile") != "0" {
		t.Fatalf("header tidak ditulis: %+v", h.Nilai)
	}
	adj := h.AmbilDaftar(models.DaftarAdjustment)
	if adj[1]["AcceptedNo"] != "UJIA123.10.2026.TX00001" || adj[1]["NoAccount"] != "123" || adj[0]["AcceptedNo"] != "" {
		t.Fatalf("akseptasi: %+v", adj)
	}
	riw := h.AmbilDaftar(models.DaftarRiwayat)
	if len(riw) == 0 || riw[len(riw)-1]["CommentSuggest"] != "Accepted by UJI" || riw[len(riw)-1]["PICSuggest"] != "UJI-K1" {
		t.Fatalf("riwayat: %+v", riw)
	}
	for _, salah := range []kontrak.UbahanKlaimTreaty{
		{Header: map[string]string{"IsAnyAcceptation": "1"}},
		{Adjustment: map[string]string{"IsApproved": "1"}},
		{FacRetro: []map[string]string{{"ReinsurerID": "UJI"}}},
	} {
		if err := k.TulisBalikKlaimTreaty(context.Background(), nil, klaimUjiKomite, adjUjiKomite, salah); !errors.Is(err,
			kontrak.ErrUbahanKlaimTreatyTidakSah) {
			t.Errorf("ubahan di luar daftar putih Non Prop %+v: %v", salah, err)
		}
	}
}
