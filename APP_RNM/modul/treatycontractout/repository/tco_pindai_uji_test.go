package repository

// Penopang uji pemindai baris - TANPA Oracle. Mengisi *sql.NullString dari
// nilai teks; teks kosong menjadi NULL, meniru driver.

import (
	"database/sql"
	"fmt"
)

type sqlNullString = sql.NullString

func isiNullString(tujuan []any, nilai []any) error {
	if len(tujuan) != len(nilai) {
		return fmt.Errorf("tujuan %d, nilai %d", len(tujuan), len(nilai))
	}
	for i, d := range tujuan {
		ns, ok := d.(*sql.NullString)
		if !ok {
			return fmt.Errorf("tujuan ke-%d bukan *sql.NullString", i)
		}
		s, _ := nilai[i].(string)
		*ns = sql.NullString{String: s, Valid: s != ""}
	}
	return nil
}
