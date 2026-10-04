package repository

// Sub-tab Occupation (tiket 40): .Property.OccupationList -> T_OCCUPATIONLIST (induk T_PROPERTY, dibedakan
// PARENT_TABLE / SRC_PATH karena rancangan memberi tabel ini tiga induk) + .TableOfLimit -> T_TABLEOFLIMIT (satu per
// okupasi), migrasi 189; ikut baca/ganti objek (objek.go). PctLimit TEKS apa adanya (butir 68.1).

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// indukOkupasi, jalurOkupasi - isi PARENT_TABLE / SRC_PATH, sama dengan yang ditulis loader (flatten.go: induk.tabel
	// dan jalur rancangan `LocationList/Property/OccupationList`).
	indukOkupasi = TabelProperty
	jalurOkupasi = "LocationList/Property/OccupationList"
)

// syaratIndukOkupasi - hanya okupasi berinduk T_PROPERTY (literal tetap, bukan masukan pengguna). ⚠️ Mengandaikan
// T_OCCUPATIONLIST beralias `o` di setiap pernyataan yang memakainya (sqlBacaOkupasi, sqlHapusObjek).
const syaratIndukOkupasi = "o.PARENT_TABLE = '" + indukOkupasi + "'"

// kolomBacaOkupasi - kolom sqlBacaOkupasi BERNAMA -> medan OkupasiObjek.
var kolomBacaOkupasi = []string{"TO_CHAR(o.PARENT_ID)", "o.OCCUPATION_ID", "o.OCCUPATION_NAME", "k.CATEGORY", "k.DESCRIPTION",
	"k.PCT_LIMIT"}

// sqlBacaOkupasi - seluruh okupasi case :1, urut property lalu SEQ_NO; TableOfLimit lewat LEFT JOIN (UNIQUE PARENT_ID).
func sqlBacaOkupasi(t tabelObjek) string {
	return "SELECT " + strings.Join(kolomBacaOkupasi, ", ") + `
FROM ` + t.okupasi + ` o
JOIN ` + t.prop + ` p ON p.ID = o.PARENT_ID
JOIN ` + t.loc + ` l ON l.ID = p.PARENT_ID
LEFT JOIN ` + t.tol + ` k ON k.PARENT_ID = o.ID
WHERE ` + syaratIndukOkupasi + ` AND l.PARENT_ID = :1
ORDER BY o.PARENT_ID, o.SEQ_NO`
}

func sqlSisipOkupasi(okupasi string) string {
	return "INSERT INTO " + okupasi + " (ID, PARENT_ID, PARENT_TABLE, SRC_PATH, SEQ_NO, ROW_UID, OCCUPATION_ID, OCCUPATION_NAME)" +
		" VALUES (:1, :2, :3, :4, :5, :6, :7, :8)"
}

func sqlSisipTableOfLimit(tol string) string {
	return "INSERT INTO " + tol + " (ID, PARENT_ID, CATEGORY, DESCRIPTION, PCT_LIMIT) VALUES (:1, :2, :3, :4, :5)"
}

// bacaOkupasi - okupasi case `id` per T_PROPERTY.ID (teks).
func (r *ObjekOracle) bacaOkupasi(ctx context.Context, t tabelObjek, id string) (map[string][]models.OkupasiObjek, error) {
	baris, err := r.db.QueryContext(ctx, sqlBacaOkupasi(t), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca okupasi objek: %w", err)
	}
	defer baris.Close()
	hasil := map[string][]models.OkupasiObjek{}
	for baris.Next() {
		teks := make(map[string]*sql.NullString, len(kolomBacaOkupasi))
		tujuan := make([]any, len(kolomBacaOkupasi))
		for i, k := range kolomBacaOkupasi {
			teks[k] = &sql.NullString{}
			tujuan[i] = teks[k]
		}
		if err := baris.Scan(tujuan...); err != nil {
			return nil, fmt.Errorf("repository: okupasi objek: %w", err)
		}
		v := func(k string) string { return teks[k].String }
		induk := v("TO_CHAR(o.PARENT_ID)")
		hasil[induk] = append(hasil[induk], models.OkupasiObjek{OccupationID: v("o.OCCUPATION_ID"),
			OccupationName: v("o.OCCUPATION_NAME"), Category: v("k.CATEGORY"), ConstructionClass: v("k.DESCRIPTION"),
			PctLimit: v("k.PCT_LIMIT")})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: okupasi objek: %w", err)
	}
	return hasil, nil
}

// sisipOkupasi - okupasi satu property (SEQ_NO 1..n), masing-masing dengan satu baris TableOfLimit.
func (r *ObjekOracle) sisipOkupasi(ctx context.Context, tx *db.Tx, t tabelObjek, idProperty string, okupasi []models.OkupasiObjek) error {
	k := db.KosongJadiNil
	for j, o := range okupasi {
		idOkupasi, err := r.db.NomorBerikut(ctx, tx, sequenceOccupationList)
		if err != nil {
			return err
		}
		idTol, err := r.db.NomorBerikut(ctx, tx, sequenceTableOfLimit)
		if err != nil {
			return err
		}
		uid, err := uidAcak()
		if err != nil {
			return err
		}
		for _, l := range []struct {
			q, tabel string
			arg      []any
		}{
			{sqlSisipOkupasi(t.okupasi), TabelOccupationList, []any{idOkupasi, idProperty, indukOkupasi, jalurOkupasi, j + 1, uid,
				k(o.OccupationID), k(o.OccupationName)}},
			{sqlSisipTableOfLimit(t.tol), TabelTableOfLimit, []any{idTol, idOkupasi, k(o.Category), k(o.ConstructionClass), k(o.PctLimit)}},
		} {
			h, err := jalankan(ctx, tx, l.q, "menyisipkan "+l.tabel, l.arg...)
			if err != nil {
				return err
			}
			if err := db.PastikanSatuBaris(h, l.tabel); err != nil {
				return err
			}
		}
	}
	return nil
}
