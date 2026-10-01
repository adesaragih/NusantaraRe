package repository

// Migrasi 050-056 lawan STRUKTUR PremiumList Life - tiket 09. TANPA Oracle.
//
// ⚠️ Sejak GILIRAN-13 (057, butir bn) pembacanya menyisir SELURUH rentang 05x,
// termasuk `ALTER ... ADD`. Nama uji di bawah tetap `050Sampai056` sebab ia
// dirujuk sebagai bukti AC di tiket 00 dan 09 - mengganti nama berarti
// memutus rujukan itu.
//
// ⛔ Kenapa ada, padahal `TestKolomDDLCocokDenganStruktur` sudah ada: penjaga
// itu mencocokkan NAMA kolom. Kolom uang yang bernama benar tetapi bertipe
// `VARCHAR2`, kolom wajib yang lupa `NOT NULL`, atau FK tanpa index lolos
// darinya. Yang diperiksa di sini TIPE, NULLABILITY, FK, dan INDEX tiap FK.
//
// ⚠️ DUA CARA (CLAUDE.md §4a): sensus pertama dijalankan dengan skrip Python
// 28-09-2026 (dilaporkan di tiket 09) - 7 tabel, 219 kolom, nol selisih
// sesudah kosakata `DATE` dikenali. Uji ini cara kedua, jalan berbeda
// (pengurai Go), dan menagih hasil yang sama.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// kolomStruktur adalah satu baris tabel kolom STRUKTUR.
type kolomStruktur struct{ Tipe, Null, Kunci string }

// bacaStrukturPolis membaca tabel kolom STRUKTUR PremiumList Life.
func bacaStrukturPolis(t *testing.T) map[string]map[string]kolomStruktur {
	t.Helper()
	isi, err := os.ReadFile(filepath.FromSlash("../../docs/STRUKTUR-TABEL-PREMIUMLIST-LIFE.md"))
	if err != nil {
		t.Fatalf("STRUKTUR tidak terbaca: %v", err)
	}
	judul := regexp.MustCompile(`^## (T_[A-Z_]+)\s*$`)
	// ⚠️ Sel kosong `| |` sah (kolom tanpa kunci) - pola `\s*` bukan satu spasi.
	barisPola := regexp.MustCompile("^\\|\\s*`([A-Z0-9_]+)`\\s*\\|\\s*([^|]*?)\\s*\\|\\s*([^|]*?)\\s*\\|\\s*([^|]*?)\\s*\\|")
	hasil := map[string]map[string]kolomStruktur{}
	var kini string
	for _, b := range strings.Split(string(isi), "\n") {
		b = strings.TrimRight(b, "\r")
		if m := judul.FindStringSubmatch(b); m != nil {
			kini = m[1]
			hasil[kini] = map[string]kolomStruktur{}
			continue
		}
		if strings.HasPrefix(b, "## ") {
			kini = ""
			continue
		}
		if kini == "" {
			continue
		}
		if m := barisPola.FindStringSubmatch(b); m != nil {
			hasil[kini][m[1]] = kolomStruktur{Tipe: m[2], Null: m[3], Kunci: m[4]}
		}
	}
	return hasil
}

// kolomDDL adalah satu kolom `CREATE TABLE`.
type kolomDDL struct {
	Tipe     string
	WajibIsi bool
}

// bacaDDLPolis membaca 050-056: kolom, FK, dan kolom pertama tiap index.
func bacaDDLPolis(t *testing.T) (map[string]map[string]kolomDDL, map[string]map[string]bool,
	map[string]map[string]bool) {
	t.Helper()
	kolom := map[string]map[string]kolomDDL{}
	fk := map[string]map[string]bool{}
	idx := map[string]map[string]bool{}
	entri, err := berkasMigrasi.ReadDir("migrations")
	if err != nil {
		t.Fatal(err)
	}
	tabelPola := regexp.MustCompile(`(?s)CREATE TABLE \{skema\}\.(\w+) \((.*?)\n\)`)
	// ⛔ Kolom yang lahir di ALTER ... ADD ikut dibaca (GILIRAN-13, 057 butir
	// bn): tanpa itu kolom lanjutan tidak pernah diperiksa tipenya.
	tambahPola := regexp.MustCompile(`(?s)ALTER TABLE \{skema\}\.(\w+) ADD \((.*?)\n\)`)
	// ⛔ 059 (seragam T_WORK_CLAIM): kolom yang DIBUANG dan FK yang lahir di
	// `ADD CONSTRAINT` ikut dibaca - keduanya di dalam blok berpelindung
	// katalog. Tanpa itu STATUS yang sudah diganti STATUS_WORK tetap "ada", dan
	// FK COVER_KEY tidak pernah ditagih index-nya.
	buangPola := regexp.MustCompile(`ALTER TABLE \{skema\}\.(\w+) DROP COLUMN (\w+)`)
	fkAlterPola := regexp.MustCompile(`ALTER TABLE \{skema\}\.(\w+) ADD CONSTRAINT \w+ FOREIGN KEY \((\w+)\)`)
	kolomPola := regexp.MustCompile(`^([A-Z0-9_]+)\s+([A-Z0-9_]+(?:\(\d+(?:,\d+)?\))?)(.*)$`)
	fkPola := regexp.MustCompile(`FOREIGN KEY \((\w+)\)`)
	idxPola := regexp.MustCompile(`CREATE (?:UNIQUE )?INDEX \{skema\}\.\w+ ON \{skema\}\.(\w+) \(\s*(\w+)`)
	for _, e := range entri {
		n := e.Name()
		// Seluruh rentang migrasi polis 05x - bukan hanya 050-056 pembuat tabel.
		if !regexp.MustCompile(`^05[0-9]_`).MatchString(n) || strings.HasSuffix(n, "_down.sql") {
			continue
		}
		b, err := berkasMigrasi.ReadFile("migrations/" + n)
		if err != nil {
			t.Fatal(err)
		}
		var bersih []string
		for _, l := range strings.Split(string(b), "\n") {
			if !strings.HasPrefix(strings.TrimSpace(l), "--") {
				bersih = append(bersih, strings.TrimRight(l, "\r"))
			}
		}
		teks := strings.Join(bersih, "\n")
		for _, m := range tabelPola.FindAllStringSubmatch(teks, -1) {
			tab := m[1]
			kolom[tab] = map[string]kolomDDL{}
			fk[tab] = map[string]bool{}
			for _, l := range strings.Split(m[2], "\n") {
				l = strings.TrimSuffix(strings.TrimSpace(l), ",")
				if f := fkPola.FindStringSubmatch(l); f != nil {
					fk[tab][f[1]] = true
				}
				if c := kolomPola.FindStringSubmatch(l); c != nil && c[1] != "CONSTRAINT" {
					kolom[tab][c[1]] = kolomDDL{Tipe: c[2], WajibIsi: strings.Contains(c[3], "NOT NULL")}
				}
			}
		}
		for _, m := range tambahPola.FindAllStringSubmatch(teks, -1) {
			tab := m[1]
			if kolom[tab] == nil {
				kolom[tab] = map[string]kolomDDL{}
			}
			for _, l := range strings.Split(m[2], "\n") {
				l = strings.TrimSuffix(strings.TrimSpace(l), ",")
				if c := kolomPola.FindStringSubmatch(l); c != nil {
					kolom[tab][c[1]] = kolomDDL{Tipe: c[2], WajibIsi: strings.Contains(c[3], "NOT NULL")}
				}
			}
		}
		for _, m := range buangPola.FindAllStringSubmatch(teks, -1) {
			delete(kolom[m[1]], m[2])
		}
		for _, m := range fkAlterPola.FindAllStringSubmatch(teks, -1) {
			if fk[m[1]] == nil {
				fk[m[1]] = map[string]bool{}
			}
			fk[m[1]][m[2]] = true
		}
		for _, m := range idxPola.FindAllStringSubmatch(teks, -1) {
			if idx[m[1]] == nil {
				idx[m[1]] = map[string]bool{}
			}
			idx[m[1]][m[2]] = true
		}
	}
	return kolom, fk, idx
}

// tipeCocok menjawab apakah kosakata STRUKTUR sepadan dengan tipe DDL.
func tipeCocok(struktur, ddl string) bool {
	switch struktur {
	case "teks":
		return strings.HasPrefix(ddl, "VARCHAR2(")
	case "angka desimal":
		// ADR-U-0003: uang/persen NUMBER(38,8), bukan NUMBER telanjang.
		return ddl == "NUMBER(38,8)"
	case "bilangan bulat":
		return regexp.MustCompile(`^NUMBER\(\d+\)$`).MatchString(ddl)
	case "DATE", "tanggal":
		return ddl == "DATE"
	}
	return false
}

func TestMigrasi050Sampai056TipeNullFKIndexSesuaiStruktur(t *testing.T) {
	struktur := bacaStrukturPolis(t)
	ddl, fk, idx := bacaDDLPolis(t)
	// ⚠️ Instrumen diuji atas jawaban yang diketahui: tujuh tabel, 219 kolom
	// (sensus Python 28-09-2026). Pengurai yang rusak akan meluluskan apa pun.
	// ⛔ 220 sejak GILIRAN-13: butir bn menambah T_WORK_POLIS.FLAG_ONGOING_POLICY
	// (057, ALTER) - diperbarui dengan sadar, bukan dilonggarkan.
	// ⛔ 225 sejak 059 (seragam T_WORK_CLAIM, 01-10-2026): STATUS menjadi
	// STATUS_WORK (nol bersih) ditambah COVER_KEY, CREATE_OP, CREATE_OP_NAME,
	// TGL_CREATE, TGL_UPDATE.
	total := 0
	for _, k := range struktur {
		total += len(k)
	}
	if len(struktur) != 7 || total != 225 {
		t.Fatalf("STRUKTUR terbaca %d tabel / %d kolom, mau 7 / 225; pengurainya rusak, "+
			"atau STRUKTUR berubah - perbarui angka ini dengan sadar", len(struktur), total)
	}
	for tab, kol := range struktur {
		d, ada := ddl[tab]
		if !ada {
			t.Errorf("%s ada di STRUKTUR, tidak dibuat migrasi 050-056", tab)
			continue
		}
		for c := range d {
			if _, ada := kol[c]; !ada {
				t.Errorf("%s.%s ada di DDL, tidak di STRUKTUR", tab, c)
			}
		}
		for c, s := range kol {
			k, ada := d[c]
			if !ada {
				t.Errorf("%s.%s ada di STRUKTUR, tidak di DDL", tab, c)
				continue
			}
			if !tipeCocok(s.Tipe, k.Tipe) {
				t.Errorf("%s.%s: STRUKTUR %q, DDL %s", tab, c, s.Tipe, k.Tipe)
			}
			if (s.Null == "tidak") != k.WajibIsi {
				t.Errorf("%s.%s: STRUKTUR nullable %q, DDL NOT NULL=%v", tab, c, s.Null, k.WajibIsi)
			}
			if strings.Contains(s.Kunci, "FK") && !fk[tab][c] {
				t.Errorf("%s.%s: STRUKTUR FK, DDL tanpa FOREIGN KEY", tab, c)
			}
		}
		for c := range fk[tab] {
			if !strings.Contains(kol[c].Kunci, "FK") {
				t.Errorf("%s.%s: FOREIGN KEY di DDL, tidak ditandai FK di STRUKTUR", tab, c)
			}
			// ⛔ FK tanpa index: setiap hapus induk memindai penuh tabel anak.
			if !idx[tab][c] {
				t.Errorf("%s.%s: FK tanpa index", tab, c)
			}
		}
	}
}
