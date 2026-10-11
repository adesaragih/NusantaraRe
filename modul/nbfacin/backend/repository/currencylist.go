package repository

// `.OfferFacIn.CurrencyList` (T_CURRENCYLIST, migrasi 200) ditulis dua tab: Premium Deduction (tiket 50) dan Payment
// (tiket 51, anak T_LISTINSTALLMENT migrasi 201). Baris diperbarui menurut NAMA mata uang (bukan hapus-sisip) supaya
// cicilan mata uang yang tetap ada tidak ikut terhapus oleh Save tab Premium Deduction (Pega: Save menghitung ulang cicilan,
// tidak membuangnya); mata uang yang hilang dibuang beserta cicilannya.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelListInstallment - cicilan per mata uang (migrasi 201).
	TabelListInstallment    = "T_LISTINSTALLMENT"
	sequenceListInstallment = "SEQ_T_LISTINSTALLMENT"
)

func sqlIDCurrencyList(cur string) string {
	return "SELECT TO_CHAR(ID), NAME FROM " + cur + " WHERE PARENT_ID = :1"
}

func sqlHapusCicilan(inst string) string { return "DELETE FROM " + inst + " WHERE PARENT_ID = :1" }

func sqlHapusSatuCurrency(cur string) string { return "DELETE FROM " + cur + " WHERE ID = :1" }

func sqlUrutCurrencyList(cur string) string {
	return "UPDATE " + cur + " SET SEQ_NO = :1 WHERE ID = :2"
}

func sqlUbahCurrencyList(cur string) string {
	return "UPDATE " + cur + " SET SEQ_NO = :1, PAY_PREMIUM = " + angkaMasuk(":2") + ", PAY_COMMISION = " + angkaMasuk(":3") +
		", PAY_DEDUCTION2 = " + angkaMasuk(":4") + ", PAY_BROKERAGE_FEE = " + angkaMasuk(":5") + ", PAY_PPH = " + angkaMasuk(":6") +
		", PAY_PPN = " + angkaMasuk(":7") + " WHERE ID = :8"
}

// simpanCurrencyList - perbarui-atau-sisip baris CurrencyList case `id` menurut NAMA; mata uang yang tidak lagi ada dibuang
// beserta cicilannya. tulisNilai=false (Save tab Payment): baris yang ada hanya dinomori ulang, PAY_* milik tab Premium
// Deduction tidak ditimpa potret Payment (lost update). Mengembalikan ID baris per nama.
func simpanCurrencyList(ctx context.Context, d *db.DB, tx *db.Tx, cur, inst, id string, baris []models.BarisDeduction, tulisNilai bool) (map[string]string, error) {
	q := sqlIDCurrencyList(cur)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := tx.QueryContext(ctx, q, id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelCurrencyList, err)
	}
	ada := map[string]string{}
	for rows.Next() {
		var idBaris, nama sql.NullString
		if err := rows.Scan(&idBaris, &nama); err != nil {
			rows.Close()
			return nil, fmt.Errorf("repository: %s: %w", TabelCurrencyList, err)
		}
		ada[nama.String] = idBaris.String
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("repository: %s: %w", TabelCurrencyList, err)
	}
	rows.Close()
	baru := map[string]bool{}
	for _, b := range baris {
		baru[b.Currency] = true
	}
	for nama, idBaris := range ada {
		if baru[nama] {
			continue
		}
		if _, err := jalankan(ctx, tx, sqlHapusCicilan(inst), "menghapus "+TabelListInstallment, idBaris); err != nil {
			return nil, err
		}
		if _, err := jalankan(ctx, tx, sqlHapusSatuCurrency(cur), "menghapus "+TabelCurrencyList, idBaris); err != nil {
			return nil, err
		}
	}
	hasil := map[string]string{}
	for i, b := range baris {
		if len(b.Currency) > lebarNamaMataUang {
			return nil, fmt.Errorf("repository: mata uang %q melebihi %d bita", b.Currency, lebarNamaMataUang)
		}
		nilai := []any{ikatDesimal(b.Premium), ikatDesimal(b.Commision), ikatDesimal(b.Deduction2), ikatDesimal(b.BrokerageFee),
			ikatDesimal(b.PPh), ikatDesimal(b.PPN)}
		if idBaris, lama := ada[b.Currency]; lama && !tulisNilai {
			if _, err := jalankan(ctx, tx, sqlUrutCurrencyList(cur), "mengurutkan "+TabelCurrencyList, i+1, idBaris); err != nil {
				return nil, err
			}
			hasil[b.Currency] = idBaris
			continue
		}
		if idBaris, lama := ada[b.Currency]; lama {
			if _, err := jalankan(ctx, tx, sqlUbahCurrencyList(cur), "mengubah "+TabelCurrencyList,
				append(append([]any{i + 1}, nilai...), idBaris)...); err != nil {
				return nil, err
			}
			hasil[b.Currency] = idBaris
			continue
		}
		urut, err := d.NomorBerikut(ctx, tx, sequenceCurrencyList)
		if err != nil {
			return nil, err
		}
		uid, err := uidAcak()
		if err != nil {
			return nil, err
		}
		kode := b.Currency
		if kode == "" {
			kode = mataUangKosong
		}
		if _, err := jalankan(ctx, tx, sqlSisipCurrencyList(cur), "menyisip "+TabelCurrencyList,
			append([]any{urut, id, i + 1, uid, kode, db.KosongJadiNil(b.Currency)}, nilai...)...); err != nil {
			return nil, err
		}
		hasil[b.Currency] = urut
	}
	return hasil, nil
}
