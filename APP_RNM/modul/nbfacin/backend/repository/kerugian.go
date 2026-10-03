package repository

// Sub-tab Loss Record (tiket 42): .Property.ListCauseOfLoss -> T_LISTCAUSEOFLOSS (rancangan + lima kolom baru, migrasi
// 191, induk T_PROPERTY) + halaman .CoinsData -> T_COINSDATA (satu per catatan); ikut baca/ganti objek (objek.go).
// Uang CLAIM / AMOUNT / PREVENTION_OF_LOSS NUMBER(38,8): ditulis fmtAngkaMasuk, dibaca db.FmtDesimal (item.go).

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// kolomBacaKerugian - kolom sqlBacaKerugian BERNAMA -> medan CatatanKerugian.
var kolomBacaKerugian = []string{"TO_CHAR(c.PARENT_ID)", "c.DATE_OF_LOSS", "d.COINS_NAME", "c.LOSS_OBJECT", "c.CURRENCY",
	fmt.Sprintf(db.FmtDesimal, "c.AMOUNT"), fmt.Sprintf(db.FmtDesimal, "c.CLAIM"), fmt.Sprintf(db.FmtDesimal, "c.PREVENTION_OF_LOSS"),
	"c.CAUSE_OF_LOSS", "c.REMARKS", "c.DETAIL"}

// sqlBacaKerugian - seluruh catatan kerugian case :1, urut property lalu SEQ_NO; CoinsData lewat LEFT JOIN.
func sqlBacaKerugian(t tabelObjek) string {
	return "SELECT " + strings.Join(kolomBacaKerugian, ", ") + `
FROM ` + t.rugi + ` c
JOIN ` + t.prop + ` p ON p.ID = c.PARENT_ID
JOIN ` + t.loc + ` l ON l.ID = p.PARENT_ID
LEFT JOIN ` + t.koas + ` d ON d.PARENT_ID = c.ID
WHERE l.PARENT_ID = :1
ORDER BY c.PARENT_ID, c.SEQ_NO`
}

// sqlSisipKerugian - :1 ID, :2 PARENT_ID, :3 SEQ_NO, :4 ROW_UID, :5..:13 medan (urut kolom; uang lewat TO_NUMBER).
func sqlSisipKerugian(rugi string) string {
	angka := func(n int) string { return fmt.Sprintf(fmtAngkaMasuk, ":"+strconv.Itoa(n)) }
	return "INSERT INTO " + rugi + " (ID, PARENT_ID, SEQ_NO, ROW_UID, DATE_OF_LOSS, LOSS_OBJECT, CURRENCY, AMOUNT, CLAIM," +
		" PREVENTION_OF_LOSS, CAUSE_OF_LOSS, REMARKS, DETAIL) VALUES (:1, :2, :3, :4, :5, :6, :7, " + angka(8) + ", " + angka(9) +
		", " + angka(10) + ", :11, :12, :13)"
}

func sqlSisipKoas(koas string) string {
	return "INSERT INTO " + koas + " (ID, PARENT_ID, COINS_NAME) VALUES (:1, :2, :3)"
}

// argKerugian - bind :5..:13 sqlSisipKerugian, urutan sama dengan daftar kolomnya. Kosong -> NULL kecuali CURRENCY
// (wajib, diperiksa services). DateOfLoss di sini sudah teks Pega (services).
func argKerugian(c models.CatatanKerugian) []any {
	k := db.KosongJadiNil
	return []any{k(c.DateOfLoss), k(c.LossObject), c.Currency, k(c.Amount), k(c.Claim), k(c.PreventionOfLoss),
		k(c.CauseOfLoss), k(c.Remarks), k(c.Detail)}
}

// bacaKerugian - catatan kerugian case `id` per T_PROPERTY.ID (teks).
func (r *ObjekOracle) bacaKerugian(ctx context.Context, t tabelObjek, id string) (map[string][]models.CatatanKerugian, error) {
	baris, err := r.db.QueryContext(ctx, sqlBacaKerugian(t), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca catatan kerugian: %w", err)
	}
	defer baris.Close()
	hasil := map[string][]models.CatatanKerugian{}
	for baris.Next() {
		teks := make(map[string]*sql.NullString, len(kolomBacaKerugian))
		tujuan := make([]any, len(kolomBacaKerugian))
		for i, k := range kolomBacaKerugian {
			teks[k] = &sql.NullString{}
			tujuan[i] = teks[k]
		}
		if err := baris.Scan(tujuan...); err != nil {
			return nil, fmt.Errorf("repository: catatan kerugian: %w", err)
		}
		v := func(k string) string { return teks[k].String }
		induk := v("TO_CHAR(c.PARENT_ID)")
		c := models.CatatanKerugian{DateOfLoss: v("c.DATE_OF_LOSS"), CoinsName: v("d.COINS_NAME"), LossObject: v("c.LOSS_OBJECT"),
			Currency: v("c.CURRENCY"), CauseOfLoss: v("c.CAUSE_OF_LOSS"), Remarks: v("c.REMARKS"), Detail: v("c.DETAIL")}
		for _, d := range []struct {
			kolom string
			ke    *string
		}{{"AMOUNT", &c.Amount}, {"CLAIM", &c.Claim}, {"PREVENTION_OF_LOSS", &c.PreventionOfLoss}} {
			if *d.ke, err = desimalTeks(TabelListCauseOfLoss+" induk "+induk, d.kolom,
				teks[fmt.Sprintf(db.FmtDesimal, "c."+d.kolom)]); err != nil {
				return nil, err
			}
		}
		hasil[induk] = append(hasil[induk], c)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: catatan kerugian: %w", err)
	}
	return hasil, nil
}

// sisipKerugian - catatan kerugian satu property (SEQ_NO 1..n), masing-masing dengan satu baris CoinsData.
func (r *ObjekOracle) sisipKerugian(ctx context.Context, tx *db.Tx, t tabelObjek, idProperty string, rugi []models.CatatanKerugian) error {
	for j, c := range rugi {
		idRugi, err := r.db.NomorBerikut(ctx, tx, sequenceListCauseOfLoss)
		if err != nil {
			return err
		}
		idKoas, err := r.db.NomorBerikut(ctx, tx, sequenceCoinsData)
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
			{sqlSisipKerugian(t.rugi), TabelListCauseOfLoss, append([]any{idRugi, idProperty, j + 1, uid}, argKerugian(c)...)},
			{sqlSisipKoas(t.koas), TabelCoinsData, []any{idKoas, idRugi, db.KosongJadiNil(c.CoinsName)}},
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
