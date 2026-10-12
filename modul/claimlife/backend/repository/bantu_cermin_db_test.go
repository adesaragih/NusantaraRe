//go:build db

package repository_test

import (
	"context"
	"fmt"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

// tulisCerminSepertiSaveRNM - cobaTulisCermin yang wajib sukses.
func tulisCerminSepertiSaveRNM(t *testing.T, ctx context.Context, d *db.DB, klaimID string) int {
	t.Helper()
	n, err := cobaTulisCermin(ctx, d, klaimID)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

// cobaTulisCermin menulis baris datar OS_AKSEPTASI_KLAIM_LIFE klaim yang sudah tersimpan, persis jalur
// Save to RNM (RALAT 10-10-2026: Register tidak lagi menulisnya). Mengembalikan cacah baris datar yang ditulis.
func cobaTulisCermin(ctx context.Context, d *db.DB, klaimID string) (int, error) {
	baca := repository.NewKlaimLife(d)
	hdr, err := baca.AmbilHeader(ctx, klaimID)
	if err != nil {
		return 0, err
	}
	if hdr == nil {
		return 0, fmt.Errorf("header %s tidak ada", klaimID)
	}
	peserta, err := baca.AmbilPeserta(ctx, klaimID)
	if err != nil {
		return 0, err
	}
	baris, err := baca.AmbilBaris(ctx, klaimID)
	if err != nil {
		return 0, err
	}
	for i := range peserta {
		peserta[i].Baris = baris[peserta[i].ID]
	}
	klaim := *hdr
	klaim.Peserta = peserta
	datar := map[string]repository.BarisLama{}
	for _, b := range repository.BarisLamaDari(models.PohonKlaim{Work: models.WorkClaim{ID: klaimID}, Klaim: klaim}) {
		b.STS_REJECT = ""
		datar[b.ID] = b
	}
	tx, err := d.Mulai(ctx)
	if err != nil {
		return 0, err
	}
	pohon := repository.NewPohonKlaim(d)
	ditulis := 0
	for _, p := range peserta {
		k := repository.KunciPesertaSumber{PLNumber: p.NomorPremiList, Sertifikat: p.NomorSertifikat, SumberID: p.SumberID}
		for _, b := range p.Baris {
			ada, err := pohon.TulisCerminSaveRNM(ctx, tx, datar[b.ID], k, repository.CerminBolehDisisip(b.KodeStatus))
			if err != nil {
				_ = tx.Rollback()
				return 0, fmt.Errorf("baris datar %s: %w", b.ID, err)
			}
			if ada {
				ditulis++
			}
		}
	}
	return ditulis, tx.Commit()
}
