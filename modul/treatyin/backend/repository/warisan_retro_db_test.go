//go:build db

package repository_test

// Penjaga penanda Retro — MENGUKUR ULANG dari Oracle, bukan mempercayai
// daftar yang ditulis tangan.
//
// ⛔ Daftar yang hanya ditulis sekali akan basi tanpa suara. Uji ini merah
// pada hari kontrak ketiga memperoleh `RetroList`, atau salah satu dari
// kedua kontrak itu kehilangannya — dan hari itulah keputusan menunda tab
// Retro harus ditinjau ulang, bukan sesudah rekonsiliasi tiket 44 selesai.
//
// ⛔ BACA SAJA. Nol transaksi, nol tulisan.

import (
	"database/sql"
	"sort"
	"testing"

	"nusantarare/inti/backend/config"
	"nusantarare/modul/treatyin/backend/repository"
)

func TestKontrakRetroTertundaMasihDuaItu(t *testing.T) {
	g, ctx := gudangBaca(t)
	cfg, _ := config.Load()

	// ⚠️ DUA LANGKAH, dan yang kedua tidak dapat dilewati. Penyaringan teks
	// di Oracle hanya mempersempit calon; kunci `RetroList` dapat muncul
	// BERSARANG, dan sapuan substring atas dokumen bersarang sudah pernah
	// menipu ronde ini sekali (`EGNPI` terbaca 155 padahal 846). Yang
	// memutuskan pengurai JSON, atas kunci PUNCAK saja.
	sqlDB, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sqlDB.Close() }()

	q := `SELECT ID FROM ` + cfg.OracleSchema + `.M_TREATY_IN
	       WHERE DBMS_LOB.INSTR(JSONDATA, '"` + repository.LarikRetroTertunda + `":[') > 0
	       ORDER BY ID`
	baris, err := sqlDB.QueryContext(ctx, q)
	if err != nil {
		t.Fatalf("menyaring calon: %v", err)
	}
	var calon []string
	for baris.Next() {
		var id string
		if err := baris.Scan(&id); err != nil {
			_ = baris.Close()
			t.Fatal(err)
		}
		calon = append(calon, id)
	}
	_ = baris.Close()
	if err := baris.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("calon dari penyaringan teks: %d", len(calon))

	var nyata []string
	for _, id := range calon {
		teks, err := g.BacaDokumenMentah(ctx, id)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		doc, err := repository.UraiDokumen(teks)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		if n := len(repository.BarisLarik(doc, repository.LarikRetroTertunda)); n > 0 {
			nyata = append(nyata, id)
			t.Logf("  %s — %d elemen %s", id, n, repository.LarikRetroTertunda)
		}
	}
	sort.Strings(nyata)

	mau := append([]string(nil), repository.KontrakRetroTertunda...)
	sort.Strings(mau)
	if len(nyata) != len(mau) {
		t.Fatalf("kontrak ber-%s kini %d (%v), penanda menyebut %d (%v).\n"+
			"Keputusan §14 menunda tab Retro ATAS DASAR hanya dua kontrak yang punya. "+
			"Dasar itu berubah — tinjau ulang keputusannya, jangan sekadar perbarui daftarnya.",
			repository.LarikRetroTertunda, len(nyata), nyata, len(mau), mau)
	}
	for i := range mau {
		if nyata[i] != mau[i] {
			t.Errorf("penanda ke-%d: %s, Oracle punya %s", i, mau[i], nyata[i])
		}
	}

	// ⛔ Dan datanya memang BELUM dimuat ke mana pun — itu seluruh sebab
	// penandanya ada. Nol tabel pendaratan boleh memuat baris retro.
	for _, p := range repository.PetaPendaratan {
		if p.Larik == repository.LarikRetroTertunda {
			t.Errorf("%s memuat %s; penanda ini menyatakan sebaliknya, dan salah satunya basi",
				p.Tabel, repository.LarikRetroTertunda)
		}
	}
}
