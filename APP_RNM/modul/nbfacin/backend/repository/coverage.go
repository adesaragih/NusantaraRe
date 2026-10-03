package repository

// Tab Coverage FIRE tahap C1 (tiket 43): .Property.PropertyItemList(m).CoverageList -> T_COVERAGELIST (rancangan berinduk
// jamak, migrasi 193; NB menulis PARENT_TABLE 'T_PROPERTYITEMLIST'), banyak baris per item urut SEQ_NO; ikut baca/ganti
// objek lewat item.go. SQL sisip, bind, dan kolom baca DIBANGUN dari SATU tabel kolom (kolomCoverage) - tidak ada
// pemasangan menurut urutan yang ditulis dua kali.

import (
	"context"
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

// kolomCoverage - satu kolom data T_COVERAGELIST (kolomdata.go).
type kolomCoverage = kolomData[models.CoverageObjek]

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

// kolomBacaCoverage - kunci induk, kunci coverage (untuk deductible, tiket 45), lalu kolom data.
func kolomBacaCoverage() []string {
	return append([]string{"TO_CHAR(v.PARENT_ID)", "TO_CHAR(v.ID)"}, kolomBacaData(kolomCoverageData, "v")...)
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
	nama, nilai := kolomSisipData(kolomCoverageData, 7)
	kolom := append(append([]string{"ID", "PARENT_ID", "PARENT_TABLE", "SRC_PATH", "SEQ_NO", "ROW_UID"}, nama...), "CURRENCY_CODE")
	isi := append(append([]string{":1", ":2", ":3", ":4", ":5", ":6"}, nilai...), ":"+strconv.Itoa(7+len(kolomCoverageData)))
	return "INSERT INTO " + cov + " (" + strings.Join(kolom, ", ") + ") VALUES (" + strings.Join(isi, ", ") + ")"
}

// argCoverage - bind kolomCoverageData (urutan tabel yang sama) + CURRENCY_CODE.
func argCoverage(c models.CoverageObjek, mataUang string) []any {
	return append(argData(kolomCoverageData, c), mataUang)
}

// bacaCoverage - coverage case `id` per T_PROPERTYITEMLIST.ID (teks), deductible terpasang (tiket 45).
func (r *ObjekOracle) bacaCoverage(ctx context.Context, t tabelObjek, id string) (map[string][]models.CoverageObjek, error) {
	baris, err := r.db.QueryContext(ctx, sqlBacaCoverage(t), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca coverage: %w", err)
	}
	defer baris.Close()
	kunci := kolomBacaCoverage()
	hasil := map[string][]models.CoverageObjek{}
	// letak coverage di `hasil` dan T_COVERAGELIST.ID-nya, untuk memasang deductible.
	type letakCoverage struct {
		induk  string
		indeks int
		id     string
	}
	var letak []letakCoverage
	for baris.Next() {
		teks, err := pindaiBaris(baris, kunci)
		if err != nil {
			return nil, fmt.Errorf("repository: coverage: %w", err)
		}
		induk := teks["TO_CHAR(v.PARENT_ID)"].String
		c, err := isiData(kolomCoverageData, "v", TabelCoverageList+" induk "+induk, teks)
		if err != nil {
			return nil, err
		}
		c.Deductibles = []models.Deductible{}
		hasil[induk] = append(hasil[induk], c)
		letak = append(letak, letakCoverage{induk, len(hasil[induk]) - 1, teks["TO_CHAR(v.ID)"].String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: coverage: %w", err)
	}
	baris.Close()
	ded, err := r.bacaDeductible(ctx, t, id)
	if err != nil {
		return nil, err
	}
	for _, l := range letak {
		if d, ada := ded[l.id]; ada && l.id != "" {
			hasil[l.induk][l.indeks].Deductibles = d
		}
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
		if err := r.sisipDeductible(ctx, tx, t, idCov, c.Deductibles); err != nil {
			return err
		}
	}
	return nil
}
