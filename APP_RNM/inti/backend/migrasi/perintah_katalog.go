package migrasi

// Blok berpelindung katalog dan kolom yang dibuang - `901_m_nav_menu_datar`
// (brief menu datar 30-09-2026).
//
// Untuk apa berkas ini: pelari menjalankan tiap pernyataan tanpa transaksi dan
// hanya menoleransi ORA-00955 pada CREATE (maju) serta ORA-00942/ORA-02289
// (mundur). Langkah yang MEMBUANG atau MENAMBAH bagian tabel yang sudah ada -
// kolom, kunci tamu, indeks - tidak termasuk, jadi pengulangannya sesudah gagal
// di tengah mati di ORA-00904, ORA-02443, ORA-01430. Langkah semacam itu
// ditulis sebagai blok PL/SQL SATU BENTUK: tanya `SYS.ALL_*` dulu, jalankan
// perintahnya lewat EXECUTE IMMEDIATE hanya bila masih perlu.
//
// Pembaca di sini dipakai penjaga (`inti/backend/penjaga`), supaya satu-satunya
// bentuk PL/SQL yang diterima adalah bentuk itu - bukan blok yang berbuat lain
// di balik kata DECLARE.

import (
	"regexp"
	"strings"
)

// PerintahKatalog adalah satu blok berpelindung katalog yang sudah diurai.
type PerintahKatalog struct {
	// Katalog - ALL_TAB_COLUMNS, ALL_CONSTRAINTS, atau ALL_INDEXES.
	Katalog string
	// Tabel yang ditanyakan (TABLE_NAME).
	Tabel string
	// Objek - nama kolom, constraint, atau indeks yang ditanyakan.
	Objek string
	// BilaAda: true = perintah jalan bila objeknya ADA (`n > 0`, membuang);
	// false = bila TIDAK ada (`n = 0`, menambah).
	BilaAda bool
	// Perintah adalah isi EXECUTE IMMEDIATE, apa adanya.
	Perintah string
}

// polaPerintahKatalog - bentuk tunggal blok berpelindung katalog. Kolom yang
// ditanyakan harus sepadan dengan katalognya (COLUMN_NAME di ALL_TAB_COLUMNS,
// dst.); pemanggil memeriksa kesepadanan itu lewat `sepadan`.
var polaPerintahKatalog = regexp.MustCompile(`(?s)^DECLARE\s+n NUMBER;\s+BEGIN\s+` +
	`SELECT COUNT\(\*\)\s+INTO\s+n\s+FROM\s+SYS\.(ALL_TAB_COLUMNS|ALL_CONSTRAINTS|ALL_INDEXES)\s+` +
	`WHERE OWNER = UPPER\('\{skema\}'\) AND TABLE_NAME = '(\w+)' AND (COLUMN_NAME|CONSTRAINT_NAME|INDEX_NAME) = '(\w+)';\s+` +
	`IF n (>|=) 0 THEN\s+EXECUTE IMMEDIATE '([^']+)';\s+END IF;\s+END;$`)

// sepadan - kolom katalog yang boleh ditanyakan setiap katalog.
var sepadan = map[string]string{
	"ALL_TAB_COLUMNS": "COLUMN_NAME",
	"ALL_CONSTRAINTS": "CONSTRAINT_NAME",
	"ALL_INDEXES":     "INDEX_NAME",
}

// BacaPerintahKatalog mengurai satu pernyataan sebagai blok berpelindung
// katalog. `ok` false berarti pernyataan itu BUKAN bentuk tersebut.
func BacaPerintahKatalog(pernyataan string) (PerintahKatalog, bool) {
	m := polaPerintahKatalog.FindStringSubmatch(strings.TrimSpace(pernyataan))
	if m == nil || sepadan[m[1]] != m[3] {
		return PerintahKatalog{}, false
	}
	return PerintahKatalog{Katalog: m[1], Tabel: m[2], Objek: m[4], BilaAda: m[5] == ">", Perintah: m[6]}, true
}

// polaAlterBuang mengenali ALTER TABLE {skema}.X DROP COLUMN Y - juga di dalam
// teks EXECUTE IMMEDIATE, tempat 901 menulisnya.
var polaAlterBuang = regexp.MustCompile(`(?is)ALTER\s+TABLE\s+\{skema\}\.(\w+)\s+DROP\s+COLUMN\s+(\w+)`)

// KolomAlterBuang membaca nama tabel dan kolom yang DIBUANG satu pernyataan.
//
// Nama tabel kosong berarti pernyataan itu tidak membuang kolom. Padanan
// `KolomAlterTambah`: tanpa ini, kolom yang dibuang migrasi lanjutan tetap
// terlihat ada oleh pembanding STRUKTUR dan penjaga tipe.
func KolomAlterBuang(pernyataan string) (string, []string) {
	m := polaAlterBuang.FindStringSubmatch(pernyataan)
	if m == nil {
		return "", nil
	}
	return strings.ToUpper(m[1]), []string{strings.ToUpper(m[2])}
}
