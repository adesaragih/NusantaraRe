//go:build db

package repository_test

// Bukti Oracle untuk jalur BACA tab dari tabel pendaratan.
//
// ⛔ BACA SAJA. Nol transaksi, nol tulisan — berkas ini tidak membuat satu
// baris pun, jadi tidak ada yang perlu dibatalkan.
//
// ⭐ Yang dibuktikan di sini BUKAN "kuerinya jalan", melainkan **isinya
// sama dengan dokumennya**. Pembacaan yang mengembalikan baris yang rapi
// tetapi salah terlihat persis seperti pembacaan yang benar.

import (
	"context"
	"database/sql"
	"testing"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/repository"
)

func gudangBaca(t *testing.T) (*repository.Gudang, context.Context) {
	t.Helper()
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("konfigurasi: %v", err)
	}
	if !cfg.PunyaOracle() {
		t.Skip("lewati: ORACLE_DSN kosong")
	}
	d, err := db.Open(cfg)
	if err != nil {
		t.Fatalf("membuka oracle: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })
	ctx := context.Background()
	if err := d.Ping(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	return repository.Baru(d), ctx
}

// ⭐ Isi tabel = isi dokumen, untuk sepuluh kontrak nyata, kedua cabang.
//
// Dibandingkan NILAI demi NILAI, bukan sekadar cacah barisnya: pemuat yang
// menggeser satu kolom menghasilkan cacah yang cocok sempurna dan isi yang
// tertukar seluruhnya.
func TestIsiTabelSamaDenganIsiDokumen(t *testing.T) {
	g, ctx := gudangBaca(t)

	diperiksa := 0
	for _, sifat := range []string{"Proportional", "NonProportional"} {
		for _, id := range kontrakTerpadat(t, ctx, sifat, 15) {
			teks, err := g.BacaDokumenMentah(ctx, id)
			if err != nil {
				t.Fatalf("%s: %v", id, err)
			}
			doc, err := repository.UraiDokumen(teks)
			if err != nil {
				t.Fatalf("%s: %v", id, err)
			}

			periode, err := g.BacaPeriodePelaporan(ctx, id)
			if err != nil {
				t.Fatal(err)
			}
			dok := repository.BarisLarik(doc, "ReportingPeriodList")
			if len(periode) != len(dok) {
				t.Errorf("%s ReportingPeriod: tabel %d baris, dokumen %d", id, len(periode), len(dok))
				continue
			}
			for i := range dok {
				diperiksa++
				mau, _ := repository.NilaiTeks(dok[i]["Period"])
				if periode[i].Periode != mau {
					t.Errorf("%s ReportingPeriod[%d].Period: tabel %q, dokumen %q", id, i, periode[i].Periode, mau)
				}
				mau, _ = repository.NilaiTeks(dok[i]["SettlementDue"])
				if periode[i].JatuhTempoBayar != mau {
					t.Errorf("%s ReportingPeriod[%d].SettlementDue: tabel %q, dokumen %q", id, i, periode[i].JatuhTempoBayar, mau)
				}
			}

			porto, err := g.BacaPortofolio(ctx, id)
			if err != nil {
				t.Fatal(err)
			}
			dok = repository.BarisLarik(doc, "Portfolio")
			if len(porto) != len(dok) {
				t.Errorf("%s Portfolio: tabel %d baris, dokumen %d", id, len(porto), len(dok))
				continue
			}
			for i := range dok {
				diperiksa++
				// ⛔ `Description` — yang terpanjang di korpus 869 aksara.
				// Kalau ada yang terpotong, di sinilah terlihat.
				mau, _ := repository.NilaiTeks(dok[i]["Description"])
				if porto[i].Keterangan != mau {
					t.Errorf("%s Portfolio[%d].Description: tabel %d aksara, dokumen %d aksara",
						id, i, len(porto[i].Keterangan), len(mau))
				}
				mau, _ = repository.NilaiTeks(dok[i]["TypePortfolio"])
				if porto[i].JenisPortfolio != mau {
					t.Errorf("%s Portfolio[%d].TypePortfolio: tabel %q, dokumen %q", id, i, porto[i].JenisPortfolio, mau)
				}
			}

			akum, err := g.BacaAkumulasi(ctx, id)
			if err != nil {
				t.Fatal(err)
			}
			dok = repository.BarisLarik(doc, "AccumulationList")
			if len(akum) != len(dok) {
				t.Errorf("%s Accumulation: tabel %d baris, dokumen %d", id, len(akum), len(dok))
			}
		}
	}
	// ⛔ Perbandingan nol nilai lulus setiap pemeriksaan di atas.
	if diperiksa == 0 {
		t.Fatal("nol nilai dibandingkan; kesepuluh kontrak kosong, atau tabelnya belum dimuat")
	}
	t.Logf("%d nilai dibandingkan satu per satu, seluruhnya sama", diperiksa)
}

// ⭐ Kontrak yang tabnya KOSONG mengembalikan nol baris, bukan galat — dan
// irisannya tidak nil, sehingga pemanggil JSON tidak perlu membedakan
// `null` dari `[]`.
func TestTabKosongMengembalikanIrisanKosongBukanGalat(t *testing.T) {
	g, ctx := gudangBaca(t)

	// Pengenal yang tidak menunjuk kontrak mana pun — kasus paling kosong
	// yang ada.
	periode, err := g.BacaPeriodePelaporan(ctx, "ZZ-TIDAK-ADA")
	if err != nil {
		t.Fatalf("pengenal tak dikenal menghasilkan galat %v; ia pertanyaan yang sah", err)
	}
	if len(periode) != 0 {
		t.Errorf("%d baris untuk pengenal yang tidak ada", len(periode))
	}
}

// ⭐ URUTAN. Dibaca `ORDER BY URUTAN`, dan dibuktikan terhadap dokumen —
// bukan terhadap dirinya sendiri.
func TestUrutanBarisSamaDenganUrutanDokumen(t *testing.T) {
	g, ctx := gudangBaca(t)

	dibuktikan := 0
	for _, id := range kontrakTerpadat(t, ctx, "Proportional", 15) {
		teks, err := g.BacaDokumenMentah(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		doc, err := repository.UraiDokumen(teks)
		if err != nil {
			t.Fatal(err)
		}
		dok := repository.BarisLarik(doc, "ReportingPeriodList")
		if len(dok) < 2 {
			continue
		}
		baris, err := g.BacaPeriodePelaporan(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		if len(baris) != len(dok) {
			t.Fatalf("%s: %d lawan %d baris", id, len(baris), len(dok))
		}
		for i := range dok {
			mau, _ := repository.NilaiTeks(dok[i]["Period"])
			if baris[i].Periode != mau {
				t.Errorf("%s urutan ke-%d: tabel %q, dokumen %q", id, i, baris[i].Periode, mau)
			}
		}
		dibuktikan++
	}
	if dibuktikan == 0 {
		t.Skip("lewati: kelima kontrak tidak punya lebih dari satu periode")
	}
	t.Logf("urutan terbukti pada %d kontrak berperiode banyak", dibuktikan)
}

// kontrakTerpadat memilih kontrak dengan BARIS PENDARATAN terbanyak.
//
// ⛔ Berbeda dari `kontrakUji`, yang memilih DOKUMEN terbesar. Dokumen
// terbesar belum tentu yang tabnya paling penuh - sepuluh kontrak terbesar
// hanya menghasilkan 26 nilai untuk dibandingkan, dan 26 perbandingan
// bukan bukti apa pun tentang 26.536 baris.
func kontrakTerpadat(t *testing.T, ctx context.Context, sifat string, n int) []string {
	t.Helper()
	cfg, _ := config.Load()
	sqlDB, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sqlDB.Close() }()

	q := `SELECT x.MASTERID FROM (
	        SELECT MASTERID, COUNT(*) n FROM ` + cfg.OracleSchema + `.M_TREATYIN_REPORTINGPERIOD GROUP BY MASTERID
	        UNION ALL
	        SELECT MASTERID, COUNT(*) n FROM ` + cfg.OracleSchema + `.M_TREATYIN_PORTFOLIO GROUP BY MASTERID
	        UNION ALL
	        SELECT MASTERID, COUNT(*) n FROM ` + cfg.OracleSchema + `.M_TREATYIN_ACCUMULATION GROUP BY MASTERID) x
	        JOIN ` + cfg.OracleSchema + `.TREATY_IN t ON t.ID = x.MASTERID
	       WHERE t.PROPORTIONTYPE = :1
	       GROUP BY x.MASTERID
	       ORDER BY SUM(x.n) DESC, x.MASTERID
	       FETCH FIRST :2 ROWS ONLY`
	baris, err := sqlDB.QueryContext(ctx, q, sifat, n)
	if err != nil {
		t.Fatalf("memilih kontrak terpadat %s: %v", sifat, err)
	}
	defer func() { _ = baris.Close() }()
	var out []string
	for baris.Next() {
		var id string
		if err := baris.Scan(&id); err != nil {
			t.Fatal(err)
		}
		out = append(out, id)
	}
	if len(out) == 0 {
		t.Skipf("lewati: nol kontrak %s punya baris pendaratan; tabelnya belum dimuat", sifat)
	}
	return out
}
