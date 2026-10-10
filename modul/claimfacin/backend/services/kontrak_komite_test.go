package services_test

// Uji kontrak `kontrak.KlaimFacInKomite` sisi penyedia (Claim Fac In) di atas tiruan Claim Fac In. Sisi pemakai (Komite
// Claim Fac In) diuji dengan kontrak palsu di modulnya sendiri; uji ini memastikan bentuk halaman yang dibaca pemakai
// (pohon objek -> item -> adjustment) sama dengan halaman Claim Fac In, dan daftar putih tulis baliknya ditegakkan.

import (
	"context"
	"errors"
	"testing"
	"time"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/repository"
	"nusantarare/modul/claimfacin/backend/services"
	"nusantarare/modul/claimfacin/backend/tiruan"
)

const (
	klaimUjiKomite = "CLM-UJI001"
	adjUjiKomite   = "UJI-ADJ-21"
)

func siapKontrakKomite(t *testing.T) (*tiruan.Gudang, services.KlaimUntukKomite, time.Time) {
	t.Helper()
	g, a := tiruan.Baru(), tiruan.AcuanBaru()
	saat := time.Date(2026, 10, 10, 10, 0, 0, 0, models.Jakarta)
	if err := g.SisipKasus(context.Background(), nil, klaimUjiKomite, "UJI-ADMIN", "UJI Admin", saat); err != nil {
		t.Fatal(err)
	}
	h := models.HalamanBaru()
	h.Setel(models.CD+"NoClaim", "UJI-K-0001")
	h.SetelDaftar(models.DaftarObjek, []models.Baris{{models.PropID: "UJI-OBJ-1", "ObjectName": "UJI OBJEK"}})
	h.SetelDaftar(models.DaftarItem(1), []models.Baris{
		{models.PropID: "UJI-IT-1", "ObjectItemName": "UJI ITEM 1"},
		{models.PropID: "UJI-IT-2", "ObjectItemName": "UJI ITEM 2"},
	})
	h.SetelDaftar(models.DaftarAdj(1, 1), []models.Baris{{models.PropID: "UJI-ADJ-11", "PaymentType": "1"}})
	h.SetelDaftar(models.DaftarAdj(1, 2), []models.Baris{
		{models.PropID: "UJI-ADJ-20", "PaymentType": "1", "AcceptanceStatus": "1", "AcceptedNo": "UJI-NO-LAMA"},
		{models.PropID: adjUjiKomite, "PaymentType": "2", "IsKomite": "1", models.PropKomiteID: "KMT-UJI001",
			"DataCommitteFacin.Remarks": "UJI CATATAN"},
	})
	g.SetelHalaman(klaimUjiKomite, h)
	l := services.Baru(g, a, func() time.Time { return saat }, false)
	return g, l.KlaimUntukKomite(), saat
}

func TestKontrakKomiteFacInBacaPosisiPohon(t *testing.T) {
	_, k, _ := siapKontrakKomite(t)
	kl, err := k.BacaKlaimFacIn(context.Background(), nil, klaimUjiKomite, adjUjiKomite)
	if err != nil {
		t.Fatal(err)
	}
	if kl.Objek != 1 || kl.Item != 2 || kl.Adjustment != 2 || kl.Tertutup {
		t.Fatalf("posisi (%d, %d, %d) mau (1, 2, 2), tertutup %v", kl.Objek, kl.Item, kl.Adjustment, kl.Tertutup)
	}
	if kl.Nilai[models.CD+"NoClaim"] != "UJI-K-0001" {
		t.Fatalf("NoClaim %q", kl.Nilai[models.CD+"NoClaim"])
	}
	// ShowTransfer LS50 `pyWorkCover.ClaimData.ClaimNo` = pyID klaim (CallActivityInputRegister 9)
	if kl.Nilai[models.CD+"ClaimNo"] != klaimUjiKomite {
		t.Fatalf("ClaimNo %q mau %q", kl.Nilai[models.CD+"ClaimNo"], klaimUjiKomite)
	}
	rows := kl.Daftar[models.DaftarAdj(1, 2)]
	if len(rows) != 2 || rows[1]["DataCommitteFacin.Remarks"] != "UJI CATATAN" {
		t.Fatalf("baris adjustment: %+v", rows)
	}
	// TT3 / TT4: tanpa adjustment
	kl, err = k.BacaKlaimFacIn(context.Background(), nil, klaimUjiKomite, "")
	if err != nil || kl.Adjustment != 0 {
		t.Fatalf("baca tanpa adjustment: %+v %v", kl.Adjustment, err)
	}
	if _, err := k.BacaKlaimFacIn(context.Background(), nil, klaimUjiKomite, "UJI-TIDAK-ADA"); !errors.Is(err,
		kontrak.ErrAdjustmentFacInTidakAda) {
		t.Fatalf("adjustment tidak ada: %v", err)
	}
	if _, err := k.BacaKlaimFacIn(context.Background(), nil, "CLM-UJI999", ""); !errors.Is(err,
		kontrak.ErrKlaimFacInTidakAda) {
		t.Fatalf("klaim tidak ada: %v", err)
	}
}

func TestKontrakKomiteFacInTulisBalikTigaTingkat(t *testing.T) {
	g, k, saat := siapKontrakKomite(t)
	u := kontrak.UbahanKlaimFacIn{
		Header: map[string]string{"AktifButton": "0", models.CD + "IsCloseFile": "1", models.CD + "IsReservedClaim": "0"},
		Adjustment: map[string]string{"AcceptanceStatus": "1", "AcceptedNo": "UJI-NO-BARU", "Notes": "UJI SETUJU",
			"IsApproved": "1", "IsPrintAccept": "1", "IsFacRetro": "0"},
		Objek:   map[string]string{"IsPrintAccept": "", "DLAStatus": "0", "IsFacretro": "0"},
		Item:    map[string]string{"IsFacretro": "0"},
		Riwayat: []kontrak.RiwayatKlaimFacIn{{Teks: "Accepted by UJI SPV - KMT-UJI001", Pelaku: "UJI-SPV", Tingkat: "UJI SPV", Saat: saat}},
	}
	if err := k.TulisBalikKlaimFacIn(context.Background(), nil, klaimUjiKomite, adjUjiKomite, u); err != nil {
		t.Fatal(err)
	}
	h := g.Halaman(klaimUjiKomite)
	if h.Ambil("AktifButton") != "0" || h.Ambil(models.CD+"IsCloseFile") != "1" {
		t.Fatalf("header: AktifButton %q IsCloseFile %q", h.Ambil("AktifButton"), h.Ambil(models.CD+"IsCloseFile"))
	}
	b := h.AmbilDaftar(models.DaftarAdj(1, 2))[1]
	if b["AcceptanceStatus"] != "1" || b["AcceptedNo"] != "UJI-NO-BARU" || b["Notes"] != "UJI SETUJU" || b["IsApproved"] != "1" ||
		b["IsPrintAccept"] != "1" {
		t.Fatalf("adjustment: %+v", b)
	}
	if lain := h.AmbilDaftar(models.DaftarAdj(1, 2))[0]; lain["AcceptedNo"] != "UJI-NO-LAMA" {
		t.Fatalf("adjustment lain tersentuh: %+v", lain)
	}
	if ob := h.AmbilDaftar(models.DaftarObjek)[0]; ob["DLAStatus"] != "0" || ob["IsFacretro"] != "0" {
		t.Fatalf("objek: %+v", ob)
	}
	if it := h.AmbilDaftar(models.DaftarItem(1))[1]; it["IsFacretro"] != "0" {
		t.Fatalf("item: %+v", it)
	}
	kr := h.AmbilDaftar(models.DaftarKronologi)
	if len(kr) != 1 || kr[0]["pyNote"] != "Accepted by UJI SPV - KMT-UJI001" || kr[0]["ASMNoteType"] != models.JenisCatatanKomite ||
		kr[0]["ASMUserID"] != "UJI SPV" || kr[0]["ASMUser"] != "UJI-SPV" {
		t.Fatalf("kronologi: %+v", kr)
	}

	// simpan halaman klaim berikutnya (halaman lama tanpa keputusan) tidak menimpa kolom milik komite
	basi := h.Salin()
	delete(basi.AmbilDaftar(models.DaftarAdj(1, 2))[1], "AcceptanceStatus")
	delete(basi.AmbilDaftar(models.DaftarAdj(1, 2))[1], "Notes")
	if err := g.SimpanHalaman(context.Background(), nil, klaimUjiKomite, basi); err != nil {
		t.Fatal(err)
	}
	if b := g.Halaman(klaimUjiKomite).AmbilDaftar(models.DaftarAdj(1, 2))[1]; b["AcceptanceStatus"] != "1" ||
		b["Notes"] != "UJI SETUJU" {
		t.Fatalf("kolom milik komite tertimpa simpan klaim: %+v", b)
	}
}

func TestKontrakKomiteFacInDaftarPutih(t *testing.T) {
	_, k, _ := siapKontrakKomite(t)
	for nama, u := range map[string]kontrak.UbahanKlaimFacIn{
		"header":     {Header: map[string]string{models.CD + "NoClaim": "UJI"}},
		"adjustment": {Adjustment: map[string]string{"PaymentType": "1"}},
		"objek":      {Objek: map[string]string{"ObjectName": "UJI"}},
		"item":       {Item: map[string]string{"ObjectItemName": "UJI"}},
	} {
		err := k.TulisBalikKlaimFacIn(context.Background(), nil, klaimUjiKomite, adjUjiKomite, u)
		if !errors.Is(err, kontrak.ErrUbahanKlaimFacInTidakSah) {
			t.Fatalf("%s di luar daftar putih diterima: %v", nama, err)
		}
	}
	// ubahan pohon tanpa adjustment
	err := k.TulisBalikKlaimFacIn(context.Background(), nil, klaimUjiKomite, "",
		kontrak.UbahanKlaimFacIn{Objek: map[string]string{"DLAStatus": "0"}})
	if !errors.Is(err, kontrak.ErrAdjustmentFacInTidakAda) {
		t.Fatalf("ubahan objek tanpa adjustment: %v", err)
	}
}

func TestKontrakKomiteFacInTutupKlaim(t *testing.T) {
	g, k, saat := siapKontrakKomite(t)
	if err := k.TutupKlaimFacIn(context.Background(), nil, klaimUjiKomite, "", "Resolved-Withdrawn", saat); !errors.Is(err,
		kontrak.ErrUbahanKlaimFacInTidakSah) {
		t.Fatalf("status di luar KomitePost: %v", err)
	}
	if err := k.KunciKlaimFacIn(context.Background(), nil, klaimUjiKomite); err != nil {
		t.Fatal(err)
	}
	// dua kasus komite anak klaim: yang memutus (TT3) dikecualikan, yang lain (TT2 menunggu) ikut ditutup
	// (pxForceCaseClose `CloseAllSubCases=true`, KomitePost_Reject S17 / KomitePost_CloseClaim S14)
	ang := []repository.AnggotaTangga{{Urut: 1, OperatorID: "ReasClaimDeptHead", Jabatan: "Claim Dept. Head"}}
	kmtTutup, err := g.BuatKasusKomite(context.Background(), nil, klaimUjiKomite, "", models.TransferTolak, "UJI-ADMIN",
		"UJI Admin", ang, saat)
	if err != nil {
		t.Fatal(err)
	}
	kmtLain, err := g.BuatKasusKomite(context.Background(), nil, klaimUjiKomite, adjUjiKomite, models.TransferAdjustment,
		"UJI-ADMIN", "UJI Admin", ang, saat)
	if err != nil {
		t.Fatal(err)
	}
	if err := k.TutupKlaimFacIn(context.Background(), nil, klaimUjiKomite, kmtTutup, kontrak.StatusKlaimDitolak,
		saat); err != nil {
		t.Fatal(err)
	}
	if st := g.Kasus[klaimUjiKomite].StatusWork; st != kontrak.StatusKlaimDitolak {
		t.Fatalf("status %q", st)
	}
	if st := g.Kasus[kmtLain].StatusWork; st != kontrak.StatusKlaimDitolak {
		t.Fatalf("kasus komite lain %s tidak ikut ditutup: %q", kmtLain, st)
	}
	if g.Kasus[kmtTutup].Tertutup() {
		t.Fatalf("kasus komite yang memutus %s ditutup kontrak (ditutup modul komite sendiri)", kmtTutup)
	}
	if err := k.TutupKlaimFacIn(context.Background(), nil, klaimUjiKomite, "", kontrak.StatusKlaimSelesai, saat); !errors.Is(err,
		kontrak.ErrKlaimFacInTertutup) {
		t.Fatalf("tutup dua kali: %v", err)
	}
	if err := k.KunciKlaimFacIn(context.Background(), nil, klaimUjiKomite); !errors.Is(err, kontrak.ErrKlaimFacInTertutup) {
		t.Fatalf("kunci klaim tertutup: %v", err)
	}
	err = k.TulisBalikKlaimFacIn(context.Background(), nil, klaimUjiKomite, adjUjiKomite,
		kontrak.UbahanKlaimFacIn{Header: map[string]string{"AktifButton": "0"}})
	if !errors.Is(err, kontrak.ErrKlaimFacInTertutup) {
		t.Fatalf("tulis balik klaim tertutup: %v", err)
	}
}
