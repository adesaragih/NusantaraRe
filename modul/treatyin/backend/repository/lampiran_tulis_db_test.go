//go:build db

package repository_test

// Jalur tulis panel Attachment — `T_STORAGE_IMAGE` + `M_ATTACHMENTTREATY_2`.
// ⛔ Seluruh tulisan di dalam transaksi yang DIBATALKAN: nol Commit di uji.
// ⛔ Kolom JSON tabel lampiran tidak dibaca dan tidak ditulis.

import (
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/treatyin/backend/models"
)

// ⭐ Pertanyaan terbuka kode kategori TERJAWAB oleh katalog RD Pega.
func TestKatalogKategoriDariMasterSebelasPasangan(t *testing.T) {
	g, ctx := gudangBaca(t)
	k, err := g.BacaKatalogKategoriLampiran(ctx)
	if err != nil {
		t.Fatal(err)
	}
	mau := map[string]string{
		"00003": "Binding, signed share Email", "00004": "Info Pack",
		"00008": "Letter of Acknowledgment / LOA", "00009": "Claim Data",
		"00002": "Approval Email", "00007": "Pega Proportional Calculation /Perhitungan Pega Proportional",
	}
	for kode, nama := range mau {
		if k[kode] != nama {
			t.Errorf("%s = %q, mau %q", kode, k[kode], nama)
		}
	}
	if len(k) < 11 {
		t.Errorf("katalog %d pasangan, mau >= 11: %v", len(k), k)
	}
}

func TestNamaAplikasiSimpananDariFolderImage(t *testing.T) {
	g, ctx := gudangBaca(t)
	app, err := g.NamaAplikasiSimpanan(ctx)
	if err != nil || app == "" {
		t.Fatalf("APPNAME %q %v", app, err)
	}
	t.Logf("APPNAME %s", app)
}

func TestCatatLampiranSahDiOracleLaluDibatalkan(t *testing.T) {
	g, ctx := gudangBaca(t)
	imageID, err := unggah.ImageIDBaru(time.Now())
	if err != nil {
		t.Fatal(err)
	}
	id, nObjek, nLampiran, err := g.CatatLampiranLaluBatalkanUntukUji(ctx, models.LampiranBaru{
		IDKontrak: "1002305", KodeKategori: "00002", NamaKategori: "Approval Email",
		NamaBerkas: "UJI — DIBATALKAN.pdf", Ekstensi: "pdf", Pengguna: "UJI",
		ImageID: imageID, URLPublik: "https://storage.googleapis.com/rnmtest/uji.pdf",
		AppFolder: "gs://rnmtest/Contract/Doc/2026/10/uji.pdf", Exp: "08/10/2026 12:30:00",
		NamaObjek: "20261008-123000-1 - UJI.pdf", App: "rnmtest",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(id) != 17 || nObjek != 1 || nLampiran != 1 {
		t.Errorf("id %q objek %d lampiran %d", id, nObjek, nLampiran)
	}
	// Dibatalkan — nol baris tersisa.
	l, err := g.BacaLampiranKontrak(ctx, "1002305")
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range l {
		if b.ID == id {
			t.Errorf("lampiran uji tertinggal: %v", b)
		}
	}
}

// GetLinkStorage_SQL atas objek nyata (lampiran kontrak 1001513, 12 Agustus 2025).
func TestBacaObjekSimpananNyata(t *testing.T) {
	g, ctx := gudangBaca(t)
	o, ada, err := g.BacaObjekSimpanan(ctx, "3E626C5BA54EE379A2A562260EC01CC7")
	if err != nil || !ada {
		t.Fatalf("ada %v err %v", ada, err)
	}
	if o.App != "rnmtest" || !strings.HasPrefix(o.AppFolder, "gs://rnmtest/Contract/Doc/") || o.NamaObjek == "" || len(o.Exp) != 19 {
		t.Errorf("objek %+v", o)
	}
}

// Delete_act: DeleteStorage_SQL + DeleteAttachment2_Sql — DIBATALKAN.
func TestHapusLampiranSahDiOracleLaluDibatalkan(t *testing.T) {
	g, ctx := gudangBaca(t)
	nObjek, nLampiran, err := g.HapusLampiranLaluBatalkanUntukUji(ctx, "1001513", "20250812172620369", "3E626C5BA54EE379A2A562260EC01CC7")
	if err != nil {
		t.Fatal(err)
	}
	if nObjek != 1 || nLampiran != 1 {
		t.Errorf("terhapus objek %d lampiran %d", nObjek, nLampiran)
	}
	if _, ada, _ := g.BacaObjekSimpanan(ctx, "3E626C5BA54EE379A2A562260EC01CC7"); !ada {
		t.Error("objek hilang sesudah rollback")
	}
}
