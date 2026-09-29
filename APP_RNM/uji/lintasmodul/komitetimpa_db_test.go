//go:build db

package lintasmodul_test

// OQ-K-05 (GILIRAN-17) terhadap skema uji Oracle: langkah 5.1
// `KomitePostAdjustment` - keadaan tangga terbaca di dalam transaksi, lalu
// SATU UPDATE bersyarat menimpa seluruh tingkat yang memutus.
//
// Tanpa instance Oracle, seluruh test di sini MELEWATI dengan pesan.

import (
	"context"
	"testing"
	"time"

	"nusantarare/inti"
	"nusantarare/modul/claimlife/repository"
	"nusantarare/modul/komite/models"
	komiterepository "nusantarare/modul/komite/repository"
)

func TestTimpaTanggaTolakAkhirMenimpaSeluruhTingkatBerkeputusan(t *testing.T) {
	db, repo, bersihkan := siapkanPohon(t)
	defer bersihkan()
	ctx := context.Background()
	p := contohPohon(t)
	tx, err := db.Mulai(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Simpan(ctx, tx, p); err != nil {
		_ = tx.Rollback()
		t.Fatalf("menyimpan pohon: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	saat := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)
	tx, err = db.Mulai(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback() }()
	kasusID, err := repo.BuatKasusKomite(ctx, tx, p.Work.ID, "UJI-A-1", inti.LiniLife, "UJI-TYPE",
		[]repository.AnggotaTangga{
			{Urut: 1, OperatorID: "UJI-OP-1", Jabatan: "UJI-J-1", Email: "uji1@uji.invalid"},
			{Urut: 2, OperatorID: "UJI-OP-2", Jabatan: "UJI-J-2", Email: "uji2@uji.invalid"},
			{Urut: 3, OperatorID: "UJI-OP-3", Jabatan: "UJI-J-3", Email: "uji3@uji.invalid"},
		}, saat)
	if err != nil {
		t.Fatalf("melahirkan kasus komite: %v", err)
	}
	inbox := komiterepository.NewInboxKomite(db)
	if err := inbox.CatatKeputusan(ctx, tx, kasusID, 1, "UJI-OP-1", models.KeputusanKomiteSetuju, "UJI ok", saat); err != nil {
		t.Fatalf("keputusan tingkat 1: %v", err)
	}
	// Tingkat 2 DILEWATI eskalasi (NULL) - satu-satunya penyimpangan sadar
	// dari 5.1: tingkat yang tidak pernah memutus tidak ditimpa.
	if err := inbox.Eskalasi(ctx, tx, kasusID, 2); err != nil {
		t.Fatalf("eskalasi tingkat 2: %v", err)
	}
	if err := inbox.CatatKeputusan(ctx, tx, kasusID, 3, "UJI-OP-3", models.KeputusanKomiteTolak, "UJI tolak", saat); err != nil {
		t.Fatalf("keputusan tingkat 3: %v", err)
	}
	sebelum, err := inbox.TanggaSebelumDitimpa(ctx, tx, kasusID)
	if err != nil || len(sebelum) != 3 {
		t.Fatalf("tangga sebelum: %v (%d)", err, len(sebelum))
	}
	if sebelum[0].Approval != "1" || sebelum[0].Komentar != "UJI ok" || sebelum[1].Approval != "" ||
		sebelum[2].Approval != "2" || sebelum[2].Komentar != "UJI tolak" {
		t.Errorf("keadaan sebelum penimpaan tidak terbaca utuh: %+v", sebelum)
	}
	if err := inbox.TimpaTanggaTolakAkhir(ctx, tx, kasusID, saat); err != nil {
		t.Fatalf("TimpaTanggaTolakAkhir: %v", err)
	}
	sesudah, err := inbox.TanggaSebelumDitimpa(ctx, tx, kasusID)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range sesudah {
		if a.Urut == 2 {
			if a.Approval != "" {
				t.Errorf("tingkat dilewati eskalasi ikut ditimpa: %+v", a)
			}
			continue
		}
		if a.Approval != models.KeputusanKomiteTolak || a.Komentar != "" {
			t.Errorf("tingkat %d sesudah 5.1: %+v, mau 2 tanpa komentar", a.Urut, a)
		}
	}
}
