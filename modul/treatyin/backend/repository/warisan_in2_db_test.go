//go:build db

package repository_test

// Bukti Oracle untuk pembacaan `M_TREATY_IN2` - empat tab, satu tabel.
//
// ⛔ BACA SAJA. Nol transaksi, nol tulisan; tidak ada yang perlu dibatalkan.

import (
	"context"
	"database/sql"
	"testing"

	"nusantarare/inti/backend/config"
)

// kontrakBerlayer memilih kontrak dengan baris `M_TREATY_IN2` TERBANYAK -
// kontrak berlayer satu tidak membuktikan apa pun tentang pengurutan.
func kontrakBerlayer(t *testing.T, ctx context.Context, n int) []string {
	t.Helper()
	cfg, _ := config.Load()
	sqlDB, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sqlDB.Close() }()

	q := `SELECT MASTERID FROM ` + cfg.OracleSchema + `.M_TREATY_IN2
	       GROUP BY MASTERID ORDER BY COUNT(*) DESC, MASTERID
	       FETCH FIRST :1 ROWS ONLY`
	baris, err := sqlDB.QueryContext(ctx, q, n)
	if err != nil {
		t.Fatalf("memilih kontrak berlayer: %v", err)
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
		t.Skip("lewati: M_TREATY_IN2 kosong")
	}
	return out
}

// ⭐ Pembacaan mengembalikan baris, dan ke-41 medannya terisi dari kolom
// yang BENAR - dibuktikan dengan mengadu dua medan yang nilainya diketahui
// berbeda bentuk, bukan sekadar "tidak kosong".
func TestBacaLayerMengisiMedanDariKolomYangBenar(t *testing.T) {
	g, ctx := gudangBaca(t)
	cfg, _ := config.Load()

	for _, id := range kontrakBerlayer(t, ctx, 5) {
		baris, err := g.BacaLayerWarisan(ctx, id)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		if len(baris) == 0 {
			t.Fatalf("%s: nol baris padahal ia kontrak dengan layer terbanyak", id)
		}
		for i, b := range baris {
			// `MASTERID` harus kontrak yang DIMINTA — kueri tanpa saring
			// mengembalikan baris kontrak lain, dan grid tetap terisi rapi.
			if b.MasterID != id {
				t.Errorf("%s baris %d: MASTERID %q — kuerinya tidak menyaring", id, i, b.MasterID)
			}
			// `SifatProporsi` hanya punya dua nilai sah; medan yang tergeser
			// akan memuat nama cedant atau sebuah angka.
			if b.SifatProporsi != "Proportional" && b.SifatProporsi != "NonProportional" {
				t.Errorf("%s baris %d: PROPORTIONTYPE %q — pemindainya bergeser", id, i, b.SifatProporsi)
			}
		}

		// ⛔ Cacah baris diadu dengan Oracle sendiri. Pemindai yang diam-diam
		// membuang baris lulus seluruh pemeriksaan di atas.
		var n int
		q := `SELECT COUNT(*) FROM ` + cfg.OracleSchema + `.M_TREATY_IN2 WHERE MASTERID = :1`
		sqlDB, err := sql.Open("oracle", cfg.OracleDSN)
		if err != nil {
			t.Fatal(err)
		}
		err = sqlDB.QueryRowContext(ctx, q, id).Scan(&n)
		_ = sqlDB.Close()
		if err != nil {
			t.Fatal(err)
		}
		if len(baris) != n {
			t.Errorf("%s: dibaca %d baris, Oracle punya %d", id, len(baris), n)
		}
	}
}

// ⭐ URUTAN LAYER menurut ANGKA, bukan teks. `LAYER` kolom teks: urutan teks
// menaruh 10 di antara 1 dan 2, dan grid limit yang layernya tertukar
// terbaca benar.
func TestLayerBerurutMenurutAngka(t *testing.T) {
	g, ctx := gudangBaca(t)

	diuji := 0
	for _, id := range kontrakBerlayer(t, ctx, 20) {
		baris, err := g.BacaLayerWarisan(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		var sebelum float64 = -1
		naik := true
		punyaBanyak := false
		for _, b := range baris {
			n, ok := angkaLayer(b.Layer)
			if !ok {
				continue
			}
			if sebelum >= 0 && n != sebelum {
				punyaBanyak = true
			}
			if n < sebelum {
				naik = false
			}
			sebelum = n
		}
		if !punyaBanyak {
			continue
		}
		diuji++
		if !naik {
			var urut []string
			for _, b := range baris {
				urut = append(urut, b.Layer)
			}
			t.Errorf("%s: layer tidak berurut naik: %v", id, urut)
		}
	}
	if diuji == 0 {
		t.Skip("lewati: nol kontrak berlayer lebih dari satu")
	}
	t.Logf("urutan layer terbukti pada %d kontrak berlayer banyak", diuji)
}

func angkaLayer(s string) (float64, bool) {
	var n float64
	ada := false
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + float64(c-'0')
		ada = true
	}
	return n, ada
}

// ⭐ Kontrak yang TIDAK ada di `M_TREATY_IN2` mengembalikan irisan KOSONG,
// bukan nil dan bukan galat.
//
// ⚠️ Keadaan ini NYATA dan sering: 510 kontrak punya `Limits[]` berisi di
// dokumennya tetapi nol baris di tabel ini. Menolaknya berarti menolak
// hampir sepertiga tabel.
func TestKontrakTanpaBarisLayerMengembalikanKosong(t *testing.T) {
	g, ctx := gudangBaca(t)

	baris, err := g.BacaLayerWarisan(ctx, "ZZ-TIDAK-ADA")
	if err != nil {
		t.Fatalf("pengenal tak dikenal menghasilkan galat %v; ia pertanyaan yang sah", err)
	}
	if len(baris) != 0 {
		t.Errorf("%d baris untuk pengenal yang tidak ada", len(baris))
	}
	if baris == nil {
		t.Error("nil, mau irisan kosong — pemanggil JSON tidak perlu membedakan null dari []")
	}
}

// ⭐ Jangkauan tabel ini TIDAK penuh, dan angkanya dikunci di sini supaya
// perubahannya terlihat.
func TestJangkauanTabelLayerTerukur(t *testing.T) {
	_, ctx := gudangBaca(t)
	cfg, _ := config.Load()
	sqlDB, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sqlDB.Close() }()

	var baris, kontrak int
	if err := sqlDB.QueryRowContext(ctx,
		`SELECT COUNT(*), COUNT(DISTINCT MASTERID) FROM `+cfg.OracleSchema+`.M_TREATY_IN2`).
		Scan(&baris, &kontrak); err != nil {
		t.Fatal(err)
	}
	if baris != 7281 {
		t.Errorf("M_TREATY_IN2 %d baris, sapuan 3 Okt 2026 menemukan 7.281 — "+
			"tabel WARISAN ini tidak boleh berubah oleh modul ini", baris)
	}
	if kontrak != 1340 {
		t.Errorf("M_TREATY_IN2 mencakup %d kontrak, terukur 1.340", kontrak)
	}
	// ⛔ Dan ia benar-benar TIDAK mencakup seluruh kontrak — bila suatu hari
	// mencakup, keempat tab berhenti punya kasus kosong dan petunjuk
	// kosongnya harus ditinjau.
	var semua int
	if err := sqlDB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM `+cfg.OracleSchema+`.TREATY_IN`).Scan(&semua); err != nil {
		t.Fatal(err)
	}
	if kontrak >= semua {
		t.Errorf("M_TREATY_IN2 kini mencakup %d dari %d kontrak; petunjuk kosong "+
			"keempat tab menyatakan sebaliknya dan perlu ditinjau", kontrak, semua)
	}
	t.Logf("M_TREATY_IN2: %d baris · %d dari %d kontrak", baris, kontrak, semua)
}
