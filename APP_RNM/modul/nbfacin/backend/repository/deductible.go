package repository

// Tab Coverage FIRE tahap C3 (tiket 45): .CoverageList(k).DeductibleList -> T_DEDUCTIBLELIST (rancangan, migrasi 194;
// induk T_COVERAGELIST ber-FK), banyak baris per coverage urut SEQ_NO; ikut baca/ganti objek lewat coverage.go.
// Mata uang kosong: kolom CURRENCY tidak disisipkan sehingga DEFAULT 'UNKNOWN' rancangan yang mengisinya (aplikasi tidak
// menulis 'UNKNOWN', K-012); dibaca kembali sebagai kosong (A169).

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
	// TabelDeductibleList - tabel rancangan T_DEDUCTIBLELIST.
	TabelDeductibleList    = "T_DEDUCTIBLELIST"
	sequenceDeductibleList = "SEQ_T_DEDUCTIBLELIST"
	// indukDeductible, jalurDeductible - isi PARENT_TABLE / SRC_PATH (sama dengan loader).
	indukDeductible = TabelCoverageList
	jalurDeductible = "LocationList/Property/PropertyItemList/CoverageList/DeductibleList"
	// mataUangTakDikenal - DEFAULT kolom CURRENCY rancangan (K-069); dibaca kembali sebagai kosong.
	mataUangTakDikenal = "UNKNOWN"
)

// kolomDeductibleData - sembilan kolom data (migrasi 194) selain CURRENCY; CURRENCY ditambahkan terpisah karena boleh
// tidak disisipkan.
var kolomDeductibleData = []kolomData[models.Deductible]{
	ka("AMOUNT", jenisAngka, func(d *models.Deductible) **apd.Decimal { return &d.Amount }),
	kt("CONDITION", func(d *models.Deductible) *string { return &d.Condition }),
	kt("INPUT_CONDITION", func(d *models.Deductible) *string { return &d.InputCondition }),
	kt("MIN_MAX", func(d *models.Deductible) *string { return &d.MinMax }),
	ka("PCT_DEDUCTIBLE", jenisAngka, func(d *models.Deductible) **apd.Decimal { return &d.PctDeductible }),
	ka("PCT_DEDUCTIBLE2", jenisAngka, func(d *models.Deductible) **apd.Decimal { return &d.PctDeductible2 }),
	ka("TIME_EXCESS", jenisTeksAngka, func(d *models.Deductible) **apd.Decimal { return &d.TimeExcess }),
	kk("TYPE_DEDUCTIBLE", func(d *models.Deductible) *string { return &d.TypeDeductible }),
	kk("TYPE_DEDUCTIBLE2", func(d *models.Deductible) *string { return &d.TypeDeductible2 }),
}

// kolomBacaDeductible - kunci coverage induk, CURRENCY, lalu kolom data (alias e).
func kolomBacaDeductible() []string {
	return append([]string{"TO_CHAR(e.PARENT_ID)", "e.CURRENCY"}, kolomBacaData(kolomDeductibleData, "e")...)
}

// sqlBacaDeductible - seluruh deductible case :1, urut coverage lalu SEQ_NO.
func sqlBacaDeductible(t tabelObjek) string {
	return "SELECT " + strings.Join(kolomBacaDeductible(), ", ") + `
FROM ` + t.ded + ` e
JOIN ` + t.cov + ` v ON v.ID = e.PARENT_ID
JOIN ` + t.item + ` i ON i.ID = v.PARENT_ID
JOIN ` + t.prop + ` p ON p.ID = i.PARENT_ID
JOIN ` + t.loc + ` l ON l.ID = p.PARENT_ID
WHERE ` + syaratIndukCoverage + ` AND l.PARENT_ID = :1
ORDER BY e.PARENT_ID, e.SEQ_NO`
}

// sqlHapusDeductible - deductible coverage berinduk item case :1 (sebelum coverage dihapus).
func sqlHapusDeductible(t tabelObjek) string {
	return "DELETE FROM " + t.ded + " WHERE PARENT_ID IN (SELECT v.ID FROM " + t.cov + " v JOIN " + t.item +
		" i ON i.ID = v.PARENT_ID JOIN " + t.prop + " p ON p.ID = i.PARENT_ID JOIN " + t.loc +
		" l ON l.ID = p.PARENT_ID WHERE " + syaratIndukCoverage + " AND l.PARENT_ID = :1)"
}

// sqlSisipDeductible - :1 ID, :2 PARENT_ID, :3 PARENT_TABLE, :4 SRC_PATH, :5 SEQ_NO, :6 ROW_UID, lalu
// kolomDeductibleData, terakhir CURRENCY bila `denganMataUang`.
func sqlSisipDeductible(ded string, denganMataUang bool) string {
	nama, nilai := kolomSisipData(kolomDeductibleData, 7)
	kolom := append([]string{"ID", "PARENT_ID", "PARENT_TABLE", "SRC_PATH", "SEQ_NO", "ROW_UID"}, nama...)
	isi := append([]string{":1", ":2", ":3", ":4", ":5", ":6"}, nilai...)
	if denganMataUang {
		kolom = append(kolom, "CURRENCY")
		isi = append(isi, ":"+strconv.Itoa(7+len(kolomDeductibleData)))
	}
	return "INSERT INTO " + ded + " (" + strings.Join(kolom, ", ") + ") VALUES (" + strings.Join(isi, ", ") + ")"
}

// adaMataUang - CURRENCY ikut disisipkan? SATU syarat untuk SQL (sqlSisipDeductible) dan bind (argDeductible) - jumlah
// bind selalu cocok.
func adaMataUang(d models.Deductible) bool { return d.Currency != "" }

// argDeductible - bind kolomDeductibleData (+ CURRENCY bila terisi).
func argDeductible(d models.Deductible) []any {
	arg := argData(kolomDeductibleData, d)
	if adaMataUang(d) {
		arg = append(arg, d.Currency)
	}
	return arg
}

// mataUangBaca - CURRENCY tersimpan -> medan: DEFAULT 'UNKNOWN' (mata uang tidak diisi) dibaca kosong (A169).
func mataUangBaca(s string) string {
	if s == mataUangTakDikenal {
		return ""
	}
	return s
}

// bacaDeductible - deductible case `id` per T_COVERAGELIST.ID (teks).
func (r *ObjekOracle) bacaDeductible(ctx context.Context, t tabelObjek, id string) (map[string][]models.Deductible, error) {
	baris, err := r.db.QueryContext(ctx, sqlBacaDeductible(t), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca deductible: %w", err)
	}
	defer baris.Close()
	kunci := kolomBacaDeductible()
	hasil := map[string][]models.Deductible{}
	for baris.Next() {
		teks, err := pindaiBaris(baris, kunci)
		if err != nil {
			return nil, fmt.Errorf("repository: deductible: %w", err)
		}
		induk := teks["TO_CHAR(e.PARENT_ID)"].String
		d, err := isiData(kolomDeductibleData, "e", TabelDeductibleList+" induk "+induk, teks)
		if err != nil {
			return nil, err
		}
		d.Currency = mataUangBaca(teks["e.CURRENCY"].String)
		hasil[induk] = append(hasil[induk], d)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: deductible: %w", err)
	}
	return hasil, nil
}

// sisipDeductible - deductible satu coverage (SEQ_NO 1..n).
func (r *ObjekOracle) sisipDeductible(ctx context.Context, tx *db.Tx, t tabelObjek, idCov string, ded []models.Deductible) error {
	for k, d := range ded {
		idDed, err := r.db.NomorBerikut(ctx, tx, sequenceDeductibleList)
		if err != nil {
			return err
		}
		uid, err := uidAcak()
		if err != nil {
			return err
		}
		arg := append([]any{idDed, idCov, indukDeductible, jalurDeductible, k + 1, uid}, argDeductible(d)...)
		h, err := jalankan(ctx, tx, sqlSisipDeductible(t.ded, adaMataUang(d)), "menyisipkan "+TabelDeductibleList, arg...)
		if err != nil {
			return err
		}
		if err := db.PastikanSatuBaris(h, TabelDeductibleList); err != nil {
			return err
		}
	}
	return nil
}
