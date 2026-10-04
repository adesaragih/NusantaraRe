package repository

// Uang dan persen ke/dari Oracle (tiket 39/42) - ADR-0034 (butir 94): SATU pasangan fungsi ikat/baca dan SATU tempat
// untuk konversi TO_NUMBER / TO_CHAR ber-NLS titik. Di dalam aplikasi nilai berbentuk *apd.Decimal; teks hanya di bind
// dan di hasil kueri. Tidak pernah bergantung NLS sesi; nol float.

import (
	"database/sql"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
)

// fmtAngkaMasuk - teks desimal bertitik -> NUMBER(38,8): 30 digit bulat, 8 desimal; NLS titik eksplisit.
const fmtAngkaMasuk = `TO_NUMBER(%s, 'FM999999999999999999999999999999D99999999', 'NLS_NUMERIC_CHARACTERS=''.,''')`

// angkaMasuk - ekspresi SQL penerima bind desimal (mis. ":13") untuk kolom NUMBER(38,8).
func angkaMasuk(bind string) string { return fmt.Sprintf(fmtAngkaMasuk, bind) }

// angkaKeluar - ekspresi SQL pembaca kolom NUMBER sebagai teks TM9 ber-NLS titik (db.FmtDesimal).
func angkaKeluar(kolom string) string { return fmt.Sprintf(db.FmtDesimal, kolom) }

// ikatDesimal - *apd.Decimal -> nilai bind: nil -> NULL, selain itu teks desimal bertitik (utils.FormatDecimal).
func ikatDesimal(d *apd.Decimal) any {
	if d == nil {
		return nil
	}
	return utils.FormatDecimal(d)
}

// bacaDesimal - teks hasil kueri (TO_CHAR TM9, atau teks kolom VARCHAR2 berisi desimal) -> *apd.Decimal; NULL / kosong
// -> nil; tak terurai -> galat (bukan diam-diam kosong, db.UraiDesimal).
func bacaDesimal(idBaris, kolom string, v *sql.NullString) (*apd.Decimal, error) {
	return db.UraiDesimal(idBaris, kolom, *v)
}

// kolomDesimal - satu kolom desimal sebaris hasil kueri: kunci = ekspresi terpilih (angkaKeluar(...) atau nama kolom
// VARCHAR2), kolom = nama untuk pesan galat, ke = tujuan.
type kolomDesimal struct {
	kunci, kolom string
	ke           **apd.Decimal
}

// bacaDesimalKe - baca beberapa kolom desimal sebaris lewat bacaDesimal; galat pertama diteruskan.
func bacaDesimalKe(idBaris string, teks map[string]*sql.NullString, daftar ...kolomDesimal) error {
	for _, d := range daftar {
		v, err := bacaDesimal(idBaris, d.kolom, teks[d.kunci])
		if err != nil {
			return err
		}
		*d.ke = v
	}
	return nil
}
