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
	// Katalog - ALL_TAB_COLUMNS, ALL_CONSTRAINTS, ALL_INDEXES, ALL_VIEWS (Tabel = Objek = nama view), atau
	// ALL_IND_COLUMNS (Objek = kolom pertama sebuah indeks).
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

// polaPerintahView - bentuk KEDUA (R/I Risk 935/938, keputusan work owner 08-10-2026): tanya SYS.ALL_VIEWS -
// satu-satunya katalog yang membedakan VIEW dari TABLE bernama sama (ALL_TAB_COLUMNS memuat keduanya). Tabel = Objek =
// nama view; penjaga menuntut perintahnya menyebut `{skema}.<nama view>`.
var polaPerintahView = regexp.MustCompile(`(?s)^DECLARE\s+n NUMBER;\s+BEGIN\s+` +
	`SELECT COUNT\(\*\)\s+INTO\s+n\s+FROM\s+SYS\.ALL_VIEWS\s+` +
	`WHERE OWNER = UPPER\('\{skema\}'\) AND VIEW_NAME = '(\w+)';\s+` +
	`IF n (>|=) 0 THEN\s+EXECUTE IMMEDIATE '([^']+)';\s+END IF;\s+END;$`)

// polaPerintahKolomIndeks - bentuk KETIGA (R/I Risk 940): tanya SYS.ALL_IND_COLUMNS apakah SUDAH ada indeks yang
// kolom PERTAMA-nya kolom itu (indeks warisan yang namanya tidak diketahui, mis. `INDEX4`), supaya CREATE INDEX tidak
// membuat indeks kembar (ORA-01408). Katalog ini memakai TABLE_OWNER, bukan OWNER.
var polaPerintahKolomIndeks = regexp.MustCompile(`(?s)^DECLARE\s+n NUMBER;\s+BEGIN\s+` +
	`SELECT COUNT\(\*\)\s+INTO\s+n\s+FROM\s+SYS\.ALL_IND_COLUMNS\s+` +
	`WHERE TABLE_OWNER = UPPER\('\{skema\}'\) AND TABLE_NAME = '(\w+)' AND COLUMN_NAME = '(\w+)' AND COLUMN_POSITION = 1;\s+` +
	`IF n (>|=) 0 THEN\s+EXECUTE IMMEDIATE '([^']+)';\s+END IF;\s+END;$`)

// BacaPerintahKatalog mengurai satu pernyataan sebagai blok berpelindung
// katalog. `ok` false berarti pernyataan itu BUKAN bentuk tersebut.
func BacaPerintahKatalog(pernyataan string) (PerintahKatalog, bool) {
	t := strings.TrimSpace(pernyataan)
	if m := polaPerintahKatalog.FindStringSubmatch(t); m != nil && sepadan[m[1]] == m[3] {
		return PerintahKatalog{Katalog: m[1], Tabel: m[2], Objek: m[4], BilaAda: m[5] == ">", Perintah: m[6]}, true
	}
	if m := polaPerintahView.FindStringSubmatch(t); m != nil {
		return PerintahKatalog{Katalog: "ALL_VIEWS", Tabel: m[1], Objek: m[1], BilaAda: m[2] == ">", Perintah: m[3]}, true
	}
	if m := polaPerintahKolomIndeks.FindStringSubmatch(t); m != nil {
		return PerintahKatalog{Katalog: "ALL_IND_COLUMNS", Tabel: m[1], Objek: m[2], BilaAda: m[3] == ">", Perintah: m[4]}, true
	}
	return PerintahKatalog{}, false
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
