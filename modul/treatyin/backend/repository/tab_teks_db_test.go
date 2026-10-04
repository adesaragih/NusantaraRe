//go:build db

package repository_test

// Bukti Oracle untuk dua tab TEKS — Jalan B, keputusan §15.
//
// ⛔ BACA SAJA. Nol transaksi, nol tulisan.
//
// ⭐ Yang dibuktikan di sini bukan "kuncinya terbaca", melainkan **ejaan
// yang dipilih benar-benar milik cabangnya** dan **isinya benar-benar
// berbeda antar ejaan**. Yang pertama menjaga teks yang salah tidak tampil;
// yang kedua membuktikan aturannya memang perlu.

import (
	"database/sql"
	"testing"

	"nusantarare/inti/backend/config"
	"nusantarare/modul/treatyin/backend/repository"
)

// kontrakBerEjaan memilih kontrak yang dokumennya memuat kunci tertentu.
func kontrakBerEjaan(t *testing.T, kunci string, n int) []string {
	t.Helper()
	cfg, _ := config.Load()
	sqlDB, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sqlDB.Close() }()

	q := `SELECT ID FROM ` + cfg.OracleSchema + `.M_TREATY_IN
	       WHERE DBMS_LOB.INSTR(JSONDATA, '"` + kunci + `":') > 0
	       ORDER BY ID FETCH FIRST :1 ROWS ONLY`
	baris, err := sqlDB.Query(q, n)
	if err != nil {
		t.Fatalf("memilih kontrak ber-%s: %v", kunci, err)
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
	return out
}

// ⭐ Ketiga ejaan `SpecialConditions*` membawa isi yang BERBEDA — dibuktikan
// pada kontrak NYATA, bukan pada data buatan.
//
// Inilah yang membuat aturan "jangan jatuh ke ejaan lain" perlu. Bila
// isinya ternyata sama, aturannya berlebihan; terukur, ia tidak pernah sama.
func TestKetigaEjaanSyaratKhususBerisiTeksYangBERBEDA(t *testing.T) {
	g, ctx := gudangBaca(t)

	diperiksa, identik := 0, 0
	for _, id := range kontrakBerEjaan(t, "SpecialConditionsp", 40) {
		teks, err := g.BacaDokumenMentah(ctx, id)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		doc, err := repository.UraiDokumen(teks)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		isi := map[string]string{}
		for _, e := range []string{"SpecialConditionsP", "SpecialConditions", "SpecialConditionsp"} {
			if v, ada := doc[e]; ada {
				s, _ := repository.NilaiTeks(v)
				isi[e] = s
			}
		}
		if len(isi) < 2 {
			continue
		}
		diperiksa++
		unik := map[string]bool{}
		for _, v := range isi {
			unik[v] = true
		}
		if len(unik) < len(isi) {
			identik++
			t.Errorf("%s: dua ejaan berisi teks yang SAMA — bila ini meluas, "+
				"aturan 'jangan jatuh ke ejaan lain' perlu ditinjau", id)
		}
	}
	if diperiksa == 0 {
		t.Skip("lewati: nol kontrak dengan lebih dari satu ejaan")
	}
	t.Logf("%d kontrak punya >1 ejaan; %d di antaranya isinya identik "+
		"(sapuan 4 Okt 2026 atas 1.854: 303 dan 0)", diperiksa, identik)
}

// ⭐ Panjang teksnya NYATA, dan itu sebab tampilannya harus dapat digulir.
func TestTeksPanjangBenarBenarPanjang(t *testing.T) {
	g, ctx := gudangBaca(t)

	maks := 0
	var dariKunci, dariKontrak string
	for _, kunci := range []string{"Exclusions", "ExclusionsP", "SpecialConditionsP"} {
		for _, id := range kontrakBerEjaan(t, kunci, 60) {
			teks, err := g.BacaDokumenMentah(ctx, id)
			if err != nil {
				t.Fatalf("%s: %v", id, err)
			}
			doc, err := repository.UraiDokumen(teks)
			if err != nil {
				t.Fatalf("%s: %v", id, err)
			}
			v, ada := doc[kunci]
			if !ada {
				continue
			}
			s, _ := repository.NilaiTeks(v)
			if len(s) > maks {
				maks, dariKunci, dariKontrak = len(s), kunci, id
			}
		}
	}
	if maks == 0 {
		t.Skip("lewati: nol teks terbaca")
	}
	t.Logf("teks terpanjang pada contoh ini: %d aksara (%s, kontrak %s)", maks, dariKunci, dariKontrak)
	// ⛔ Jauh melampaui satu baris sel. Bila suatu hari tidak, tampilan
	// bergulir berhenti punya sebab dan pernyataannya di CSS perlu ditinjau.
	if maks < 1000 {
		t.Errorf("terpanjang hanya %d aksara; sapuan menemukan sampai 23.453, "+
			"dan tampilan bergulir dipasang atas dasar itu", maks)
	}
}

// ⭐ Kunci yang ADA tidak pernah bernilai KOSONG — sifat yang dipakai
// `PilihTabTeks` untuk memutuskan "ejaan lain berisi".
func TestKunciTeksYangAdaTidakPernahKosong(t *testing.T) {
	g, ctx := gudangBaca(t)

	diperiksa := 0
	for _, kunci := range []string{"Exclusions", "ExclusionsP", "SpecialConditions", "SpecialConditionsP", "SpecialConditionsp"} {
		for _, id := range kontrakBerEjaan(t, kunci, 20) {
			teks, err := g.BacaDokumenMentah(ctx, id)
			if err != nil {
				t.Fatalf("%s: %v", id, err)
			}
			doc, err := repository.UraiDokumen(teks)
			if err != nil {
				t.Fatalf("%s: %v", id, err)
			}
			v, ada := doc[kunci]
			if !ada {
				continue
			}
			diperiksa++
			s, _ := repository.NilaiTeks(v)
			if s == "" {
				t.Errorf("%s: kunci %s ada tetapi kosong — `PilihTabTeks` menganggap "+
					"kunci-ada berarti teks-ada, dan andaian itu patah di sini", id, kunci)
			}
		}
	}
	if diperiksa == 0 {
		t.Skip("lewati: nol kunci teks terbaca")
	}
	t.Logf("%d kunci teks diperiksa, seluruhnya berisi", diperiksa)
}
