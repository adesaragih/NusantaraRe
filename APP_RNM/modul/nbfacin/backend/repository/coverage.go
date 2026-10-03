package repository

// Tab Coverage FIRE tahap C1 (tiket 43): .Property.PropertyItemList(m).CoverageList -> T_COVERAGELIST (rancangan berinduk
// jamak, migrasi 193; NB menulis PARENT_TABLE 'T_PROPERTYITEMLIST'), banyak baris per item urut SEQ_NO; ikut baca/ganti
// objek lewat item.go. SQL sisip, bind, dan kolom baca DIBANGUN dari SATU tabel kolom (kolomCoverage) - tidak ada
// pemasangan menurut urutan yang ditulis dua kali.

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelCoverageList - tabel rancangan T_COVERAGELIST.
	TabelCoverageList    = "T_COVERAGELIST"
	sequenceCoverageList = "SEQ_T_COVERAGELIST"
	// indukCoverage, jalurCoverage - isi PARENT_TABLE / SRC_PATH (sama dengan loader).
	indukCoverage = TabelPropertyItemList
	jalurCoverage = "LocationList/Property/PropertyItemList/CoverageList"
)

// syaratIndukCoverage - hanya coverage berinduk item (literal tetap). Mengandaikan alias `v`.
const syaratIndukCoverage = "v.PARENT_TABLE = '" + indukCoverage + "'"

// Jenis kolom coverage.
const (
	jenisTeks      = iota // VARCHAR2 teks
	jenisAngka            // NUMBER(38,8) - angkaMasuk / angkaKeluar
	jenisTeksAngka        // VARCHAR2 berisi teks desimal (rancangan) - ikatDesimal / bacaDesimal tanpa TO_NUMBER
)

// kolomCoverage - satu kolom data T_COVERAGELIST dan medan CoverageObjek-nya (teks ATAU desimal).
type kolomCoverage struct {
	nama  string
	jenis int
	teks  func(*models.CoverageObjek) *string
	des   func(*models.CoverageObjek) **apd.Decimal
}

func kt(nama string, f func(*models.CoverageObjek) *string) kolomCoverage {
	return kolomCoverage{nama: nama, jenis: jenisTeks, teks: f}
}

func ka(nama string, jenis int, f func(*models.CoverageObjek) **apd.Decimal) kolomCoverage {
	return kolomCoverage{nama: nama, jenis: jenis, des: f}
}

// kolomCoverageData - 25 kolom data (migrasi 193) sesudah kolom sistem; CURRENCY_CODE ditambahkan terpisah.
var kolomCoverageData = []kolomCoverage{
	kt("COVERAGE", func(c *models.CoverageObjek) *string { return &c.Coverage }),
	kt("OLDID", func(c *models.CoverageObjek) *string { return &c.OldID }),
	kt("COVERAGE_NOTE", func(c *models.CoverageObjek) *string { return &c.CoverageNote }),
	kt("COVERAGE_BASIS", func(c *models.CoverageObjek) *string { return &c.CoverageBasis }),
	kt("DAY", func(c *models.CoverageObjek) *string { return &c.Day }),
	ka("TSI", jenisAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.TSI }),
	kt("INDEMNITY", func(c *models.CoverageObjek) *string { return &c.Indemnity }),
	ka("RATE", jenisAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.Rate }),
	ka("RATE_OJK", jenisAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.RateOJK }),
	ka("FIRST_LOSS", jenisTeksAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.FirstLoss }),
	ka("DISCOUNT_PERCENTAGE", jenisAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.DiscountPercentage }),
	ka("TSI_LIABILITY", jenisAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.TSILiability }),
	ka("NET_RATE", jenisAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.NetRate }),
	ka("LIMITOF_LIABILITY", jenisAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.LimitOfLiability }),
	ka("PCT_LO_L", jenisAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.PctLoL }),
	ka("PRO_RATE_PERCENT", jenisAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.ProRatePercent }),
	ka("INDEMNITY_PERCENTAGE", jenisAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.IndemnityPercentage }),
	ka("FIRST_SCALE", jenisTeksAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.FirstScale }),
	ka("SUBLIMIT", jenisAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.Sublimit }),
	ka("LOST_LIMIT", jenisTeksAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.LostLimit }),
	ka("EML_PML", jenisTeksAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.EmlPml }),
	ka("DISCOUNT", jenisAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.Discount }),
	ka("PREMIUM", jenisAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.Premium }),
	kt("CONDITIONS", func(c *models.CoverageObjek) *string { return &c.Conditions }),
	ka("PCT_ADJUSTMENT", jenisAngka, func(c *models.CoverageObjek) **apd.Decimal { return &c.PctAdjustment }),
}

// kunciBaca - ekspresi terpilih kolom (alias v).
func (k kolomCoverage) kunciBaca() string {
	if k.jenis == jenisAngka {
		return angkaKeluar("v." + k.nama)
	}
	return "v." + k.nama
}

// kolomBacaCoverage - kunci induk + kolom data.
func kolomBacaCoverage() []string {
	hasil := []string{"TO_CHAR(v.PARENT_ID)"}
	for _, k := range kolomCoverageData {
		hasil = append(hasil, k.kunciBaca())
	}
	return hasil
}

// sqlBacaCoverage - seluruh coverage case :1 (berinduk item), urut item lalu SEQ_NO.
func sqlBacaCoverage(t tabelObjek) string {
	return "SELECT " + strings.Join(kolomBacaCoverage(), ", ") + `
FROM ` + t.cov + ` v
JOIN ` + t.item + ` i ON i.ID = v.PARENT_ID
JOIN ` + t.prop + ` p ON p.ID = i.PARENT_ID
JOIN ` + t.loc + ` l ON l.ID = p.PARENT_ID
WHERE ` + syaratIndukCoverage + ` AND l.PARENT_ID = :1
ORDER BY v.PARENT_ID, v.SEQ_NO`
}

// sqlHapusCoverage - coverage berinduk item case :1 (sebelum item dihapus).
func sqlHapusCoverage(t tabelObjek) string {
	return "DELETE FROM " + t.cov + " v WHERE " + syaratIndukCoverage + " AND v.PARENT_ID IN (SELECT i.ID FROM " + t.item +
		" i JOIN " + t.prop + " p ON p.ID = i.PARENT_ID JOIN " + t.loc + " l ON l.ID = p.PARENT_ID WHERE l.PARENT_ID = :1)"
}

// sqlSisipCoverage - :1 ID, :2 PARENT_ID, :3 PARENT_TABLE, :4 SRC_PATH, :5 SEQ_NO, :6 ROW_UID, lalu kolomCoverageData,
// terakhir CURRENCY_CODE.
func sqlSisipCoverage(cov string) string {
	kolom := []string{"ID", "PARENT_ID", "PARENT_TABLE", "SRC_PATH", "SEQ_NO", "ROW_UID"}
	nilai := []string{":1", ":2", ":3", ":4", ":5", ":6"}
	n := len(nilai)
	for _, k := range kolomCoverageData {
		n++
		kolom = append(kolom, k.nama)
		b := ":" + strconv.Itoa(n)
		if k.jenis == jenisAngka {
			b = angkaMasuk(b)
		}
		nilai = append(nilai, b)
	}
	kolom = append(kolom, "CURRENCY_CODE")
	nilai = append(nilai, ":"+strconv.Itoa(n+1))
	return "INSERT INTO " + cov + " (" + strings.Join(kolom, ", ") + ") VALUES (" + strings.Join(nilai, ", ") + ")"
}

// argCoverage - bind kolomCoverageData (urutan tabel yang sama) + CURRENCY_CODE.
func argCoverage(c models.CoverageObjek, mataUang string) []any {
	var arg []any
	for _, k := range kolomCoverageData {
		if k.jenis == jenisTeks {
			arg = append(arg, db.KosongJadiNil(*k.teks(&c)))
		} else {
			arg = append(arg, ikatDesimal(*k.des(&c)))
		}
	}
	return append(arg, mataUang)
}

// bacaCoverage - coverage case `id` per T_PROPERTYITEMLIST.ID (teks).
func (r *ObjekOracle) bacaCoverage(ctx context.Context, t tabelObjek, id string) (map[string][]models.CoverageObjek, error) {
	baris, err := r.db.QueryContext(ctx, sqlBacaCoverage(t), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca coverage: %w", err)
	}
	defer baris.Close()
	kunci := kolomBacaCoverage()
	hasil := map[string][]models.CoverageObjek{}
	for baris.Next() {
		teks := make(map[string]*sql.NullString, len(kunci))
		tujuan := make([]any, len(kunci))
		for i, k := range kunci {
			teks[k] = &sql.NullString{}
			tujuan[i] = teks[k]
		}
		if err := baris.Scan(tujuan...); err != nil {
			return nil, fmt.Errorf("repository: coverage: %w", err)
		}
		induk := teks["TO_CHAR(v.PARENT_ID)"].String
		var c models.CoverageObjek
		for _, k := range kolomCoverageData {
			v := teks[k.kunciBaca()]
			if k.jenis == jenisTeks {
				*k.teks(&c) = v.String
				continue
			}
			if *k.des(&c), err = bacaDesimal(TabelCoverageList+" induk "+induk, k.nama, v); err != nil {
				return nil, err
			}
		}
		hasil[induk] = append(hasil[induk], c)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: coverage: %w", err)
	}
	return hasil, nil
}

// sisipCoverage - coverage satu item (SEQ_NO 1..n); CURRENCY_CODE = mata uang item (K-069/K-012).
func (r *ObjekOracle) sisipCoverage(ctx context.Context, tx *db.Tx, t tabelObjek, idItem, mataUang string, cov []models.CoverageObjek) error {
	for k, c := range cov {
		idCov, err := r.db.NomorBerikut(ctx, tx, sequenceCoverageList)
		if err != nil {
			return err
		}
		uid, err := uidAcak()
		if err != nil {
			return err
		}
		arg := append([]any{idCov, idItem, indukCoverage, jalurCoverage, k + 1, uid}, argCoverage(c, mataUang)...)
		h, err := jalankan(ctx, tx, sqlSisipCoverage(t.cov), "menyisipkan "+TabelCoverageList, arg...)
		if err != nil {
			return err
		}
		if err := db.PastikanSatuBaris(h, TabelCoverageList); err != nil {
			return err
		}
	}
	return nil
}
