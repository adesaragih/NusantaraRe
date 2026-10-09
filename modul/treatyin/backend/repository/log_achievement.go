package repository

// Tombol `Submit` sub-tab Achievement (Treaty In dan Adjustment) —
// `InsertToLogAchievement` → RDB `InsertToLogAchievement_SQL`.
//
// ⭐ Keputusan pemakai 8 Oktober 2026: sasarannya `LOG_ACHIEVEMENT` apa adanya
// (tabel sudah ada, 17 kolom sama persis dengan SQL ekspor), setiap klik
// MENYISIPKAN (log = riwayat; Pega tidak mencegah duplikat), baris ber-Quarter
// kosong DILEWATI (services).
//
//	INSERT INTO POOLDATA.LOG_ACHIEVEMENT (MASTERID, QUARTER, QUARTERYEAR,
//	  IDCURRENCY, CURRENCY, PREMIUM, RICOMM, BROKERAGE, NETPREMIUM, PAIDCLAIM,
//	  CASHCALLCLAIM, OUTSTANDINGCLAIM, INCUREDCLAIM, TOTAL, LOSSRATIO,
//	  PXCREATEOPNAME, INSERTDATE) VALUES (…, sysdate)
//
// ⛔ Kolom angka `NUMBER` diikat `TO_NUMBER(koef) / POWER(10, skala)` —
// koefisien bilangan bulat, jadi tidak bergantung `NLS_NUMERIC_CHARACTERS`
// sesi. Angka kosong = `NULL`.

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
)

// TabelLogAchievement - log Submit Achievement (dibaca nol menu di ekspor).
const TabelLogAchievement = "LOG_ACHIEVEMENT"

// BarisLogSiap - satu baris log yang SUDAH diurai services: teks apa adanya,
// angka sebagai desimal (`nil` = kosong).
type BarisLogSiap struct {
	Quarter, QuarterYear, CurrencyID, Currency string
	// Urut kolom: PREMIUM, RICOMM, BROKERAGE, NETPREMIUM, PAIDCLAIM,
	// CASHCALLCLAIM, OUTSTANDINGCLAIM, INCUREDCLAIM, TOTAL, LOSSRATIO.
	Angka [10]*apd.Decimal
}

var kolomAngkaLog = []string{"PREMIUM", "RICOMM", "BROKERAGE", "NETPREMIUM", "PAIDCLAIM",
	"CASHCALLCLAIM", "OUTSTANDINGCLAIM", "INCUREDCLAIM", "TOTAL", "LOSSRATIO"}

// CatatLogAchievement menyisipkan seluruh baris dalam SATU transaksi.
func (g *Gudang) CatatLogAchievement(ctx context.Context, masterID, operator string, baris []BarisLogSiap) (err error) {
	if len(baris) == 0 {
		return nil
	}
	nama, err := g.db.Qualify(TabelLogAchievement)
	if err != nil {
		return err
	}
	kolom := append(append([]string{"MASTERID", "QUARTER", "QUARTERYEAR", "IDCURRENCY", "CURRENCY"}, kolomAngkaLog...),
		"PXCREATEOPNAME", "INSERTDATE")
	nilai := []string{":1", ":2", ":3", ":4", ":5"}
	n := 6
	for range kolomAngkaLog {
		nilai = append(nilai, fmt.Sprintf("(TO_NUMBER(:%d) / POWER(10, :%d))", n, n+1))
		n += 2
	}
	nilai = append(nilai, fmt.Sprintf(":%d", n), "SYSDATE")
	q := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", nama, strings.Join(kolom, ", "), strings.Join(nilai, ", "))
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	for _, b := range baris {
		args := []any{masterID, kosongNil(b.Quarter), kosongNil(b.QuarterYear), kosongNil(b.CurrencyID), kosongNil(b.Currency)}
		for _, d := range b.Angka {
			koef, skala := pecahAngka(d)
			args = append(args, koef, skala)
		}
		args = append(args, kosongNil(operator))
		if _, err = tx.ExecContext(ctx, q, args...); err != nil {
			return fmt.Errorf("repository: menyisipkan %s %s: %w", TabelLogAchievement, masterID, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("repository: mengikat %s %s: %w", TabelLogAchievement, masterID, err)
	}
	return nil
}

func kosongNil(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}
