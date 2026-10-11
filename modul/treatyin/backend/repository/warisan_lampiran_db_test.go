//go:build db

package repository_test

// Bukti Oracle untuk panel Attachment — `M_ATTACHMENTTREATY_2`, 43 baris.
//
// ⛔ BACA SAJA. Nol transaksi, nol tulisan, nol `Commit`. Tabel ini WARISAN
// yang hidup, dan ke-43 barisnya tidak boleh berubah oleh pekerjaan ini.

import (
	"database/sql"
	"testing"

	"nusantarare/inti/backend/config"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// ⭐ Ke-43 baris UTUH. Dijalankan di luar transaksi mana pun, jadi ia
// melihat apa yang orang lain lihat.
func TestEmpatPuluhTigaBarisLampiranUtuh(t *testing.T) {
	_, ctx := gudangBaca(t)
	cfg, _ := config.Load()
	sqlDB, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sqlDB.Close() }()

	var n int
	if err := sqlDB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM `+cfg.OracleSchema+`.M_ATTACHMENTTREATY_2`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	// ⛔ LANTAI, bukan angka PERSIS — ralat 6 Oktober 2026.
	//
	// Angka persis membuat gerbang ini merah setiap kali bisnis berjalan
	// normal: Pega menambah baris ke tabel warisan ini, dan penambahan itu
	// bukan regresi. Terukur hari itu — `TREATY_IN` dan `M_TREATY_IN` sudah
	// 1.855 dari 1.854 (kontrak `1002059`), `M_TREATY_IN_DETAIL` 27.618 dari
	// 27.617.
	//
	// ⚠️ YANG HILANG karenanya, dan itu dinyatakan: penambahan baris oleh
	// MODUL INI tidak lagi tertangkap cacahnya. Yang menangkapnya
	// `TestWarisanHanyaDibaca`, yang menyapu naskah SQL dan menolak
	// `INSERT`/`UPDATE`/`DELETE` terhadap tabel ini — penjaga yang
	// sesungguhnya, dan yang 6 Oktober 2026 diperluas supaya mencakupnya.
	if n < 43 {
		t.Errorf("M_ATTACHMENTTREATY_2 %d baris, terukur 43 pada 4 Oktober 2026 — "+
			"tabel WARISAN ini tidak boleh ditulis oleh modul ini", n)
	}
}

// ⭐ Katalog dibaca DARI DATA, dan ketujuh pasangannya persis yang terukur.
func TestKatalogKategoriTerbacaDariData(t *testing.T) {
	g, ctx := gudangBaca(t)

	katalog, err := g.BacaKatalogKategoriLampiran(ctx)
	if err != nil {
		t.Fatal(err)
	}
	mau := map[string]string{
		"00000": "Others",
		"00001": "Analysed Email",
		"00002": "Approval Email",
		"00005": "Summary Treaty Leader",
		"00006": "Assessment Inward Treaty Form / Format Analisa Treaty",
		"00007": "Pega Proportional Calculation /Perhitungan Pega Proportional",
		"00010": "Offer Email",
	}
	if len(katalog) != len(mau) {
		t.Errorf("%d kategori terbaca, terukur %d: %v", len(katalog), len(mau), katalog)
	}
	for kode, nama := range mau {
		if katalog[kode] != nama {
			t.Errorf("kode %s = %q, terukur %q", kode, katalog[kode], nama)
		}
	}
	// ⛔ Dan keempat kode yang belum diketahui memang TIDAK ADA di data —
	// itulah sebab pasangannya tidak dapat dibaca. Bila salah satunya
	// muncul, pertanyaan terbukanya terjawab dan penandanya harus dicabut.
	for _, kode := range models.KodeKategoriBelumDipastikan {
		if nama, ada := katalog[kode]; ada {
			t.Errorf("kode %s kini ADA di data bernama %q — "+
				"PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md terjawab untuk kode ini, "+
				"cabut penandanya", kode, nama)
		}
	}
}

// ⛔ Satu kode, satu nama. Bila satu `CATEGORY_ID` punya dua `CATEGORY`
// berbeda, katalog yang memakai `MAX()` memilih satu dan menyembunyikan
// yang lain — dan berkas di bawah nama yang tersembunyi menjadi sulit
// ditemukan.
func TestSatuKodeSatuNama(t *testing.T) {
	_, ctx := gudangBaca(t)
	cfg, _ := config.Load()
	sqlDB, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sqlDB.Close() }()

	baris, err := sqlDB.QueryContext(ctx,
		`SELECT CATEGORY_ID, COUNT(DISTINCT CATEGORY) FROM `+cfg.OracleSchema+`.M_ATTACHMENTTREATY_2
		  WHERE CATEGORY_ID IS NOT NULL GROUP BY CATEGORY_ID HAVING COUNT(DISTINCT CATEGORY) > 1`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = baris.Close() }()
	for baris.Next() {
		var kode string
		var n int
		if err := baris.Scan(&kode, &n); err != nil {
			t.Fatal(err)
		}
		t.Errorf("kode %s punya %d nama berbeda; katalog hanya menampilkan satu", kode, n)
	}
}

// ⭐ Pembacaan per kontrak menyaring, dan isinya sampai ke modelnya.
func TestBacaLampiranKontrakMenyaringDanMengisi(t *testing.T) {
	g, ctx := gudangBaca(t)
	cfg, _ := config.Load()
	sqlDB, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sqlDB.Close() }()

	var id string
	if err := sqlDB.QueryRowContext(ctx,
		`SELECT TREATYID FROM `+cfg.OracleSchema+`.M_ATTACHMENTTREATY_2
		  WHERE TREATYID IS NOT NULL GROUP BY TREATYID ORDER BY COUNT(*) DESC, TREATYID
		  FETCH FIRST 1 ROWS ONLY`).Scan(&id); err != nil {
		t.Skipf("lewati: nol lampiran berpengenal kontrak: %v", err)
	}

	baris, err := g.BacaLampiranKontrak(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if len(baris) == 0 {
		t.Fatalf("%s: nol lampiran padahal ia kontrak dengan lampiran terbanyak", id)
	}
	var n int
	if err := sqlDB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM `+cfg.OracleSchema+`.M_ATTACHMENTTREATY_2 WHERE TREATYID = :1`,
		id).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if len(baris) != n {
		t.Errorf("%s: dibaca %d lampiran, Oracle punya %d", id, len(baris), n)
	}
	for i, b := range baris {
		if b.ID == "" {
			t.Errorf("%s baris %d: ID kosong", id, i)
		}
		if b.NamaBerkas == "" {
			t.Errorf("%s baris %d: nama berkas kosong", id, i)
		}
	}
	t.Logf("kontrak %s: %d lampiran, kategori pertama %q (%s)",
		id, len(baris), baris[0].NamaKategori, baris[0].KodeKategori)

	// ⛔ Kuerinya BENAR-BENAR menyaring — pengenal yang tidak ada
	// mengembalikan nol, bukan ke-43 barisnya.
	kosong, err := g.BacaLampiranKontrak(ctx, "ZZ-TIDAK-ADA")
	if err != nil {
		t.Fatal(err)
	}
	if len(kosong) != 0 {
		t.Errorf("pengenal tak dikenal mengembalikan %d lampiran", len(kosong))
	}
	if kosong == nil {
		t.Error("nil, mau irisan kosong")
	}
}

// ⚠️ Lampiran yang kontraknya TIDAK ADA di `TREATY_IN` — terukur ada, dan
// itu temuan, bukan kegagalan. Dikunci supaya perubahannya terlihat.
func TestLampiranYatimTerukur(t *testing.T) {
	_, ctx := gudangBaca(t)
	cfg, _ := config.Load()
	sqlDB, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sqlDB.Close() }()

	var semua, cocok int
	if err := sqlDB.QueryRowContext(ctx,
		`SELECT COUNT(DISTINCT TREATYID) FROM `+cfg.OracleSchema+`.M_ATTACHMENTTREATY_2`).Scan(&semua); err != nil {
		t.Fatal(err)
	}
	if err := sqlDB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM (SELECT DISTINCT TREATYID t FROM `+cfg.OracleSchema+`.M_ATTACHMENTTREATY_2) x
		   JOIN `+cfg.OracleSchema+`.TREATY_IN t ON t.ID = x.t`).Scan(&cocok); err != nil {
		t.Fatal(err)
	}
	t.Logf("TREATYID berbeda=%d, cocok ke TREATY_IN=%d, yatim=%d", semua, cocok, semua-cocok)
	if semua != 7 || cocok != 6 {
		t.Errorf("terukur 4 Oktober 2026: 7 pengenal, 6 cocok (1 yatim); kini %d dan %d",
			semua, cocok)
	}
}

// ⭐ Aturan nama berkas ditegakkan terhadap nama berkas NYATA yang sudah
// tersimpan — bukan terhadap contoh buatan.
//
// ⚠️ Ia TIDAK menuntut seluruhnya lulus: ke-43 baris lahir sebelum
// aturannya ada. Yang dicatat cacahnya, supaya besarnya persoalan terbaca
// sebelum seseorang memutuskan menegakkan aturan itu surut.
func TestNamaBerkasNyataDiukurTerhadapAturannya(t *testing.T) {
	_, ctx := gudangBaca(t)
	cfg, _ := config.Load()
	sqlDB, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sqlDB.Close() }()

	baris, err := sqlDB.QueryContext(ctx,
		`SELECT FILENAME FROM `+cfg.OracleSchema+`.M_ATTACHMENTTREATY_2 WHERE FILENAME IS NOT NULL`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = baris.Close() }()

	var total, aman int
	var contoh []string
	for baris.Next() {
		var n string
		if err := baris.Scan(&n); err != nil {
			t.Fatal(err)
		}
		total++
		if services.NamaBerkasAman(n) {
			aman++
		} else if len(contoh) < 3 {
			contoh = append(contoh, n)
		}
	}
	if total == 0 {
		t.Skip("lewati: nol nama berkas")
	}
	t.Logf("%d nama berkas tersimpan, %d lolos aturan spanduk, %d tidak. Contoh yang tidak: %q",
		total, aman, total-aman, contoh)
	// ⛔ Pemeriksanya harus MENGGIGIT dan harus MELEWATKAN — bila ia
	// menolak seluruh 43, ia bukan aturan melainkan pintu tertutup.
	if aman == 0 {
		t.Error("nol dari nama berkas nyata lolos; aturannya menutup pintu yang sistem lama buka")
	}
}
