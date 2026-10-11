package services_test

// Lampiran klaim untuk layar Komite Claim Prop - bagian `kontrak.KlaimTreatyKomite` sisi penyedia (perintah work owner
// 11-10-2026): kategori master + cacah (termasuk 0), dokumen klaim, isi berkas, dan View Office Online.

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/services"
	"nusantarare/modul/claimprop/backend/tiruan"
)

func TestKontrakKomiteLampiranKlaim(t *testing.T) {
	ctx := context.Background()
	g, a := tiruan.Baru(), tiruan.AcuanBaru()
	saat := time.Date(2026, 10, 11, 9, 30, 0, 0, models.Jakarta)
	if err := g.SisipKasus(ctx, nil, klaimUjiKomite, "UJI-ADMIN", "UJI Admin", saat); err != nil {
		t.Fatal(err)
	}
	g.SetelHalaman(klaimUjiKomite, models.HalamanBaru())
	g.KategoriDok = map[string]string{"ADU": "Approval", "File": "File", "LOD": "LOD"}
	kunci := models.KunciInstans(klaimUjiKomite)
	g.Dokumen = []models.BarisDokumenKlaim{
		{ID: "UJI-D1", Tanggal: saat, IDPega: kunci, NamaFile: "uji-a.pdf", MIME: "pdf", Kategori1: "ADU",
			StorageID: "UJI-IMG-1", Operator: "UJI-ADMIN"},
		{ID: "UJI-D2", Tanggal: saat, IDPega: kunci, NamaFile: "uji-b.xlsx", MIME: "xlsx", Kategori1: "ADU",
			StorageID: "UJI-IMG-2", Operator: "UJI-ADMIN"},
		{ID: "UJI-D3", Tanggal: saat, IDPega: kunci, NamaFile: "uji-c.pdf", MIME: "pdf", Kategori1: "File"},
		{ID: "UJI-D9", Tanggal: saat, IDPega: "CLMP-UJI-LAIN", NamaFile: "uji-lain.pdf", MIME: "pdf", Kategori1: "File",
			StorageID: "UJI-IMG-1"},
	}
	b := g.BerkasBaru()
	b.Unggahan = []penyimpanan.MasukUnggah{{Isi: []byte("%PDF-UJI")}, {Isi: []byte("UJI-XL")}}
	k := services.Baru(g, a, func() time.Time { return saat }, false).DenganPenyimpanan(b).KlaimUntukKomite()

	d, err := k.LampiranKlaimTreaty(ctx, klaimUjiKomite)
	if err != nil {
		t.Fatal(err)
	}
	cacah := map[string]int{}
	for _, x := range d.Kategori {
		cacah[x.ID] = x.Jumlah
	}
	if len(d.Kategori) != 3 || cacah["ADU"] != 2 || cacah["File"] != 1 || cacah["LOD"] != 0 {
		t.Fatalf("kategori %+v (kategori 0 tetap dikirim - penyaringan urusan Komite)", d.Kategori)
	}
	if len(d.Lampiran) != 3 || d.Lampiran[0].ID != "UJI-D3" || !d.Lampiran[2].AdaObjek || d.Lampiran[0].AdaObjek ||
		d.Lampiran[2].Tanggal != "11/10/2026 09:30" {
		t.Fatalf("dokumen %+v", d.Lampiran)
	}

	f, err := k.BukaLampiranKlaimTreaty(ctx, klaimUjiKomite, "UJI-D1", "UJI-KOMITE")
	if err != nil {
		t.Fatal(err)
	}
	isi, _ := io.ReadAll(f.Isi)
	_ = f.Isi.Close()
	if string(isi) != "%PDF-UJI" || f.Nama != "uji-a.pdf" || f.Mime != "application/pdf" {
		t.Fatalf("isi %q nama %q mime %q", isi, f.Nama, f.Mime)
	}
	for _, c := range []struct {
		klaim, lid string
		mau        error
	}{
		{klaimUjiKomite, "UJI-D3", penyimpanan.ErrObjekTidakAda},
		{klaimUjiKomite, "UJI-D9", kontrak.ErrLampiranKlaimTidakAda}, // dokumen klaim lain
		{"CLMP-UJI404", "UJI-D1", kontrak.ErrKlaimTreatyTidakAda},
	} {
		if _, err := k.BukaLampiranKlaimTreaty(ctx, c.klaim, c.lid, "UJI-KOMITE"); !errors.Is(err, c.mau) {
			t.Errorf("buka %s/%s: %v, mau %v", c.klaim, c.lid, err, c.mau)
		}
	}

	if url, err := k.TautanOfficeLampiranKlaimTreaty(ctx, klaimUjiKomite, "UJI-D2", "UJI-KOMITE"); err != nil ||
		url != "UJI-URL-UJI-IMG-2" {
		t.Fatalf("office xlsx: %q %v", url, err)
	}
	if _, err := k.TautanOfficeLampiranKlaimTreaty(ctx, klaimUjiKomite, "UJI-D1", "UJI-KOMITE"); !errors.Is(err,
		kontrak.ErrLampiranKlaimBukanOffice) {
		t.Fatalf("office pdf: %v", err)
	}
}
