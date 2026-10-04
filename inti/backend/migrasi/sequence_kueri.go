package migrasi

// Blok sequence yang nilai awalnya dihitung dari data - `810_seq_client_org`
// Company Detail (perintah work owner 04-10-2026: "buat seq aja, start-nya dari
// id max+1"; "max dari org yang ada di client").
//
// Untuk apa berkas ini: Oracle hanya menerima ANGKA pada `CREATE SEQUENCE ...
// START WITH`, sehingga sequence yang harus mulai sesudah nomor tertinggi yang
// sudah terpakai tidak dapat ditulis sebagai SQL biasa tanpa membekukan angka
// satu basis data (DEV) ke basis data lain (PROD). Bentuk ini menghitungnya di
// basis data tempat migrasi berjalan.
//
// Ini bentuk PL/SQL KEDUA yang diterima penjaga, di samping blok berpelindung
// katalog (perintah_katalog.go) - sama sempitnya: tanya `SYS.ALL_SEQUENCES`
// dulu (pengulangan sesudah gagal di tengah tidak mati di ORA-00955), satu
// SELECT baca-saja mengisi `awal`, lalu satu CREATE SEQUENCE atas nama yang
// SAMA dengan yang ditanyakan.

import (
	"regexp"
	"strings"
)

// SequenceDariKueri adalah satu blok sequence-dari-kueri yang sudah diurai.
type SequenceDariKueri struct {
	// Nama sequence - yang ditanyakan katalog SAMA dengan yang dibuat.
	Nama string
	// Kueri - SELECT pengisi `awal`, tanpa `INTO awal`.
	Kueri string
	// Opsi - sisa CREATE SEQUENCE sesudah `START WITH <awal>`.
	Opsi string
}

var (
	// polaSequenceDariKueri - bentuk tunggal blok ini. `[^;]` menjaga kueri
	// tetap SATU pernyataan.
	polaSequenceDariKueri = regexp.MustCompile(`(?s)^DECLARE\s+n\s+NUMBER;\s+awal\s+NUMBER;\s+BEGIN\s+` +
		`SELECT COUNT\(\*\)\s+INTO\s+n\s+FROM\s+SYS\.ALL_SEQUENCES\s+` +
		`WHERE SEQUENCE_OWNER = UPPER\('\{skema\}'\) AND SEQUENCE_NAME = '(\w+)';\s+` +
		`IF n = 0 THEN\s+SELECT\s+([^;]+?)\s+INTO\s+awal\s+FROM\s+([^;]+);\s+` +
		`EXECUTE IMMEDIATE 'CREATE SEQUENCE \{skema\}\.(\w+) START WITH ' \|\| awal \|\| ' ([^']*)';\s+` +
		`END IF;\s+END;$`)
	// polaOpsiSequence - opsi yang dikenal saja; MAXVALUE dan semacamnya tidak.
	polaOpsiSequence = regexp.MustCompile(`^(INCREMENT BY [0-9]+|NOCACHE|NOCYCLE)( (NOCACHE|NOCYCLE))*$`)
	// polaKueriMenulis - kata yang membuat kueri bukan lagi baca-saja.
	polaKueriMenulis = regexp.MustCompile(
		`(?i)\b(INSERT|UPDATE|DELETE|MERGE|DROP|ALTER|CREATE|TRUNCATE|GRANT|REVOKE|EXECUTE|DECLARE|BEGIN|INTO)\b`)
)

// BacaSequenceDariKueri mengurai satu pernyataan sebagai blok sequence-dari-
// kueri. `ok` false berarti pernyataan itu BUKAN bentuk tersebut.
func BacaSequenceDariKueri(pernyataan string) (SequenceDariKueri, bool) {
	m := polaSequenceDariKueri.FindStringSubmatch(strings.TrimSpace(pernyataan))
	if m == nil || m[1] != m[4] || !polaOpsiSequence.MatchString(m[5]) {
		return SequenceDariKueri{}, false
	}
	kueri := "SELECT " + m[2] + " FROM " + m[3]
	if polaKueriMenulis.MatchString(kueri) {
		return SequenceDariKueri{}, false
	}
	return SequenceDariKueri{Nama: m[1], Kueri: kueri, Opsi: m[5]}, true
}
