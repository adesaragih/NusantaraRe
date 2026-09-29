package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/utils"
)

// fmtTanggalOracle adalah bentuk tanggal yang diminta dari Oracle. Ia cocok
// dengan utils.TanggalWaktu, dan ParseTanggal juga menerima bentuk pendeknya.
const FmtTanggalOracle = `TO_CHAR(%s, 'YYYY-MM-DD HH24:MI:SS')`

// fmtDesimal memaksa Oracle menyerahkan angka sebagai TEKS, bukan sebagai
// bilangan pecahan biner.
//
// ⛔ Uang tidak pernah melewati float (ADR-U-0003, ADR-U-0016). Membiarkan
// driver menyerahkan NUMBER sebagai float64 melanggar aturan itu di tempat yang
// paling sulit terlihat - karena itu konversinya dilakukan di dalam SQL.
//
// 'TM9' memberi bentuk desimal terpendek tanpa notasi ilmiah. Argumen NLS
// memaksa titik sebagai pemisah desimal: tanpa itu, sesi ber-NLS Indonesia
// mengembalikan koma dan pembacaannya gagal senyap.
const FmtDesimal = `TO_CHAR(%s, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')`

// pastikanSatuBaris menuntut sebuah pernyataan menyentuh tepat satu baris.
func PastikanSatuBaris(hasil sql.Result, nama string) error {
	n, err := hasil.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: mencacah baris %s: %w", nama, err)
	}
	if n != 1 {
		return fmt.Errorf("repository: %s menyentuh %d baris, mau 1", nama, n)
	}
	return nil
}

// uraiDesimal mengubah teks yang datang dari Oracle menjadi desimal.
//
// ⛔ Galat urai DIKEMBALIKAN, tidak ditelan. Menelannya - yang dikerjakan ronde
// 1 lewat `if d, err := ...; err == nil` - membuat nilai yang tidak terbaca
// pulang sebagai kosong, dan kosong tidak dapat dibedakan dari nol. Angka yang
// rusak harus terdengar, bukan hilang diam-diam.
//
// NULL dan teks kosong BUKAN galat: keduanya berarti "tidak ada nilai"
// (ADR-U-0027), dan menghasilkan desimal nil.
func UraiDesimal(idBaris, kolom string, v sql.NullString) (*apd.Decimal, error) {
	if !v.Valid || strings.TrimSpace(v.String) == "" {
		return nil, nil
	}
	d, err := utils.ParseDecimal(v.String)
	if err != nil {
		return nil, fmt.Errorf("repository: %s kolom %s bernilai %q yang tidak terurai: %w",
			idBaris, kolom, v.String, err)
	}
	return d, nil
}

func KosongJadiNil(s string) any {
	if s == "" {
		return nil
	}
	return s
}
