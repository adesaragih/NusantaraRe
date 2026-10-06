//go:build db

package repository_test

// AMBANG PEMBALIKAN §17 — DIUKUR ULANG dari Oracle, bukan dipercaya dari
// daftar yang ditulis tangan.
//
// ⛔ Angka yang hanya ditulis sekali akan basi tanpa suara. Uji ini merah
// pada hari kontrak KETIGA memperoleh `RetroList` — dan hari itulah
// KEPUTUSAN §17 harus ditinjau ulang, bukan sekadar angkanya diperbarui.
//
// ⛔ BACA SAJA. Nol transaksi, nol tulisan.

import (
	"database/sql"
	"sort"
	"testing"

	"nusantarare/inti/backend/config"
	"nusantarare/modul/treatyin/backend/repository"
)

func TestKontrakRetroJarangMasihDuaItu(t *testing.T) {
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
	       WHERE DBMS_LOB.INSTR(JSONDATA, '"` + repository.LarikRetroJarang + `":[') > 0
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
		if n := len(repository.BarisLarik(doc, repository.LarikRetroJarang)); n > 0 {
			nyata = append(nyata, id)
			t.Logf("  %s — %d elemen %s", id, n, repository.LarikRetroJarang)
		}
	}
	sort.Strings(nyata)

	mau := append([]string(nil), repository.KontrakRetroJarang...)
	sort.Strings(mau)
	if len(nyata) != len(mau) {
		t.Fatalf("kontrak ber-%s kini %d (%v), penanda menyebut %d (%v).\n"+
			"KEPUTUSAN §17 tidak membangun tab Retro ATAS DASAR hanya dua kontrak yang punya. "+
			"Dasar itu berubah - TINJAU ULANG keputusannya, jangan sekadar perbarui daftarnya.",
			repository.LarikRetroJarang, len(nyata), nyata, len(mau), mau)
	}
	for i := range mau {
		if nyata[i] != mau[i] {
			t.Errorf("penanda ke-%d: %s, Oracle punya %s", i, mau[i], nyata[i])
		}
	}

	// ⛔ Dan datanya memang BELUM dimuat ke mana pun — itu seluruh sebab
	// penandanya ada. Nol tabel pendaratan boleh memuat baris retro.
	for _, p := range repository.PetaPendaratan {
		if p.Larik == repository.LarikRetroJarang {
			t.Errorf("%s memuat %s; penanda ini menyatakan sebaliknya, dan salah satunya basi",
				p.Tabel, repository.LarikRetroJarang)
		}
	}
}

// AMBANG KEDUA §17 — `IsMultipleRetro`, ambangnya LIMA.
//
// ⛔ §17 berdiri di atas DUA angka, dan sampai 4 Oktober 2026 baru satu yang
// dijaga. Keputusan yang setengah angkanya dihafal adalah keputusan yang
// setengahnya dapat basi tanpa ada yang tahu.
//
// ⚠️ Angka inilah yang paling mungkin bergerak lebih dulu: `IsMultipleRetro`
// disetel saat kontrak disusun, sementara `RetroList` baru terisi sesudah
// retro benar-benar dibagi. Kontrak keenam ber-`IsMultipleRetro` adalah
// tanda paling awal bahwa Retro mulai dipakai.
//
// ⛔ BACA SAJA, dan pengurai JSON yang memutuskan - bukan penyaringan teks.
// Kunci dapat muncul bersarang, dan sapuan substring atas dokumen bersarang
// sudah pernah menipu pekerjaan ini sekali (`EGNPI` terbaca 155 padahal 846).
func TestAmbangRetroMultipleMasihLima(t *testing.T) {
	g, ctx := gudangBaca(t)
	cfg, _ := config.Load()

	sqlDB, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sqlDB.Close() }()

	// ⛔ DUA LANGKAH PENYARINGAN, dan yang kedua ada supaya yang pertama
	// boleh sempit.
	//
	// Menyaring hanya pada nama kuncinya mengembalikan 1.536 calon, dan
	// membaca 1.536 CLOB membuat uji ini berjalan setengah menit. Menyaring
	// langsung pada `"kunci":"true"` mengembalikan lima - tetapi penyaringan
	// sesempit itu akan DIAM bila Pega suatu hari menulis spasi sesudah
	// titik dua, atau nilai ketiga selain true/false.
	//
	// Jadi bentuk yang tidak terduga dihitung TERPISAH, dengan satu kueri
	// yang tidak membaca satu CLOB pun. Selama cacahnya nol, penyaringan
	// sempit aman; begitu ia bukan nol, uji ini merah dan menyebut sebabnya.
	var takTerduga int
	qCek := `SELECT COUNT(*) FROM ` + cfg.OracleSchema + `.M_TREATY_IN
	          WHERE DBMS_LOB.INSTR(JSONDATA, '"` + repository.KunciRetroMultiple + `"') > 0
	            AND DBMS_LOB.INSTR(JSONDATA, '"` + repository.KunciRetroMultiple + `":"true"') = 0
	            AND DBMS_LOB.INSTR(JSONDATA, '"` + repository.KunciRetroMultiple + `":"false"') = 0`
	if err := sqlDB.QueryRowContext(ctx, qCek).Scan(&takTerduga); err != nil {
		t.Fatal(err)
	}
	if takTerduga != 0 {
		t.Fatalf("%d dokumen memuat kunci %s dalam bentuk selain `:\"true\"` atau `:\"false\"`; "+
			"penyaringan sempit di bawah akan MELEWATKANNYA. Lebarkan penyaringannya "+
			"sebelum mempercayai angkanya.", takTerduga, repository.KunciRetroMultiple)
	}

	q := `SELECT ID FROM ` + cfg.OracleSchema + `.M_TREATY_IN
	       WHERE DBMS_LOB.INSTR(JSONDATA, '"` + repository.KunciRetroMultiple + `":"true"') > 0
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

	var benar []string
	for _, id := range calon {
		teks, err := g.BacaDokumenMentah(ctx, id)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		doc, err := repository.UraiDokumen(teks)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		v, ada := doc[repository.KunciRetroMultiple]
		if !ada {
			continue
		}
		// ⚠️ Nilainya TEKS `"true"`/`"false"`, bukan boolean JSON - pola yang
		// sama dengan `TreatyLeader`.
		if s, _ := repository.NilaiTeks(v); s == "true" {
			benar = append(benar, id)
		}
	}

	t.Logf("%s = \"true\" pada %d dokumen: %v", repository.KunciRetroMultiple, len(benar), benar)
	if len(benar) != repository.AmbangRetroMultiple {
		t.Errorf("%s = \"true\" pada %d dokumen, ambang §17 adalah %d.\n"+
			"KEPUTUSAN §17 tidak membangun tab Retro ATAS DASAR angka itu. "+
			"Dasar itu berubah - TINJAU ULANG keputusannya, jangan sekadar perbarui ambangnya.",
			repository.KunciRetroMultiple, len(benar), repository.AmbangRetroMultiple)
	}
}
