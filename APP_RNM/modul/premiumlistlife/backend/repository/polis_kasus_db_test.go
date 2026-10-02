//go:build db

package repository_test

// Kasus polis menulis pembuat dan waktu - brief seragam kolom 01-10-2026 §4,
// terhadap skema uji Oracle: lahir dengan CREATE_OP/CREATE_OP_NAME = akun
// pelaku, TGL_CREATE = TGL_UPDATE = SYSDATE, COVER_KEY kosong; setiap ubah
// baris kasus (pindah tahap, tutup) memajukan TGL_UPDATE, TGL_CREATE tetap.
//
// Tanpa instance Oracle, seluruh test di sini MELEWATI dengan pesan.

import (
	"testing"
	"time"

	intidb "nusantarare/inti/backend/db"
	"nusantarare/modul/premiumlistlife/backend/models"
	"nusantarare/modul/premiumlistlife/backend/repository"
	"nusantarare/uji/skemauji"
)

func TestKasusPolisMenulisPembuatDanWaktu(t *testing.T) {
	// Pintu skema uji bersama - bantu_skemauji_db_test.go.
	_, _, ctx := pasangSkemaUji(t)
	db, err := skemauji.BukaRepositori()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	kerja := repository.NewWorkPolis(db)

	dalamTx := func(langkah func(tx *intidb.Tx) error) {
		t.Helper()
		tx, err := db.Mulai(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = tx.Rollback() }()
		if err := langkah(tx); err != nil {
			t.Fatal(err)
		}
		if err := tx.Commit(); err != nil {
			t.Fatal(err)
		}
	}

	const id = "UJI-NBLF-059"
	awal, err := models.SusunKasusPolisBaru(models.FlagPolisPenawaran)
	if err != nil {
		t.Fatal(err)
	}
	dalamTx(func(tx *intidb.Tx) error {
		return kerja.SisipKasusBaru(ctx, tx, id, awal.DenganPembuat("UJI-AKUN"), time.Now())
	})
	lahir, err := kerja.Keadaan(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if lahir.Status != models.TahapPolisPenawaran || lahir.CreateOp != "UJI-AKUN" ||
		lahir.CreateOpName != "UJI-AKUN" || lahir.CoverKey != "" ||
		lahir.TglCreate.IsZero() || lahir.TglUpdate.IsZero() {
		t.Fatalf("kasus lahir: %+v", lahir)
	}

	// Jeda melewati satu detik: DATE Oracle berpresisi detik.
	time.Sleep(1100 * time.Millisecond)
	dalamTx(func(tx *intidb.Tx) error {
		return kerja.PindahTahap(ctx, tx, id, models.TahapPolisPenawaran, models.TahapPolisDetail)
	})
	pindah, err := kerja.Keadaan(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if pindah.Status != models.TahapPolisDetail || !pindah.TglUpdate.After(lahir.TglUpdate) ||
		!pindah.TglCreate.Equal(lahir.TglCreate) {
		t.Errorf("pindah tahap: %+v (lahir %+v)", pindah, lahir)
	}

	time.Sleep(1100 * time.Millisecond)
	dalamTx(func(tx *intidb.Tx) error {
		return kerja.TutupKasus(ctx, tx, id, models.TahapPolisDetail, models.StatusPolisSelesai)
	})
	tutup, err := kerja.Keadaan(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if tutup.Status != models.StatusPolisSelesai || tutup.Position != "" ||
		!tutup.TglUpdate.After(pindah.TglUpdate) || !tutup.TglCreate.Equal(lahir.TglCreate) ||
		tutup.CreateOp != "UJI-AKUN" {
		t.Errorf("tutup kasus: %+v", tutup)
	}
}
