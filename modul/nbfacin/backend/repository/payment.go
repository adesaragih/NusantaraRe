package repository

// Tab Payment kasus FIRE (tiket 51): jumlah cicilan (T_GENERAL_POLIS.PAY_INSTALLMENT), Total Payment / Net Premium per mata
// uang (T_CURRENCYLIST.SUM_TOTAL_PAYMENT / PAY_NET_PREMIUM) dan `.Policy.Payment.ListInstallment` (T_LISTINSTALLMENT) -
// migrasi 201. Baris CurrencyList ditulis bersama tab Premium Deduction (`simpanCurrencyList`).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
)

// lebarPersenCicilan - T_LISTINSTALLMENT.INSTALLMENT_PERCENTAGE VARCHAR2(50).
const lebarPersenCicilan = 50

// PenyimpanPayment - data tab Payment.
type PenyimpanPayment interface {
	// BacaPayment - ErrKasusTidakAda bila case tidak ada / bukan LINI Fac In.
	BacaPayment(ctx context.Context, id string) (models.KasusPayment, error)
	// TulisPayment - jumlah cicilan, CurrencyList (nilai Payment, Total Payment, Net Premium), dan cicilan tiap mata uang
	// diganti utuh di transaksi pemanggil. ErrKasusTidakAda bila case tidak ada.
	TulisPayment(ctx context.Context, tx *db.Tx, id string, k models.SimpanPayment) error
}

// PaymentOracle - PenyimpanPayment atas Oracle.
type PaymentOracle struct{ db *db.DB }

// NewPaymentOracle merakit penyimpan tab Payment.
func NewPaymentOracle(d *db.DB) *PaymentOracle { return &PaymentOracle{db: d} }

type tabelPayment struct{ work, general, cur, inst string }

func (r *PaymentOracle) tabel() (tabelPayment, error) {
	var t tabelPayment
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelWorkPolis, &t.work}, {TabelGeneralPolis, &t.general}, {TabelCurrencyList, &t.cur}, {TabelListInstallment, &t.inst}} {
		q, err := r.db.Qualify(x.nama)
		if err != nil {
			return t, err
		}
		*x.ke = q
	}
	return t, nil
}

func sqlBacaKepalaPayment(t tabelPayment) string {
	return "SELECT TO_CHAR(g.PAY_INSTALLMENT) FROM " + t.work + " w LEFT JOIN " + t.general +
		" g ON g.ID = w.ID WHERE w.ID = :1 AND w.LINI = :2"
}

func sqlBacaCicilan(t tabelPayment) string {
	return "SELECT c.NAME, TO_CHAR(i.INSTALLMENT_NO), i.DUE_DATE, i.INSTALLMENT_PERCENTAGE FROM " + t.inst + " i JOIN " + t.cur +
		" c ON c.ID = i.PARENT_ID WHERE c.PARENT_ID = :1 ORDER BY c.SEQ_NO, i.SEQ_NO"
}

func sqlUbahInstallment(general string) string {
	return "UPDATE " + general + " SET PAY_INSTALLMENT = :1 WHERE ID = :2"
}

func sqlUbahTotalPayment(cur string) string {
	return "UPDATE " + cur + " SET SUM_TOTAL_PAYMENT = " + angkaMasuk(":1") + ", PAY_NET_PREMIUM = " + angkaMasuk(":2") + " WHERE ID = :3"
}

func sqlSisipCicilan(inst string) string {
	return "INSERT INTO " + inst + " (ID, PARENT_ID, SEQ_NO, ROW_UID, CURRENCY_CODE, INSTALLMENT_NO, DUE_DATE, INSTALLMENT_PERCENTAGE," +
		" PREMIUM, RI_COMMISION, DEDUCTION2, BROKERAGE_FEE, PPH, PPN, DISCOUNT, PAYMENT_TOTAL) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, " +
		angkaMasuk(":9") + ", " + angkaMasuk(":10") + ", " + angkaMasuk(":11") + ", " + angkaMasuk(":12") + ", " + angkaMasuk(":13") + ", " +
		angkaMasuk(":14") + ", " + angkaMasuk(":15") + ", " + angkaMasuk(":16") + ")"
}

// BacaPayment - lihat PenyimpanPayment.
func (r *PaymentOracle) BacaPayment(ctx context.Context, id string) (models.KasusPayment, error) {
	k := models.KasusPayment{Cicilan: map[string][]models.CicilanPayment{}}
	t, err := r.tabel()
	if err != nil {
		return k, err
	}
	var n sql.NullString
	err = r.db.QueryRowContext(ctx, sqlBacaKepalaPayment(t), id, LiniFacIn).Scan(&n)
	if errors.Is(err, sql.ErrNoRows) {
		return k, ErrKasusTidakAda
	}
	if err != nil {
		return k, fmt.Errorf("repository: baca kepala tab Payment: %w", err)
	}
	k.Installment = n.String
	baris, err := r.db.QueryContext(ctx, sqlBacaCicilan(t), id)
	if err != nil {
		return k, fmt.Errorf("repository: baca %s: %w", TabelListInstallment, err)
	}
	defer baris.Close()
	for baris.Next() {
		var nama, no, jatuh, persen sql.NullString
		if err := baris.Scan(&nama, &no, &jatuh, &persen); err != nil {
			return k, fmt.Errorf("repository: %s: %w", TabelListInstallment, err)
		}
		c := models.CicilanPayment{DueDate: jatuh.String}
		if no.String != "" {
			if c.InstallmentNo, err = strconv.Atoi(no.String); err != nil {
				return k, fmt.Errorf("repository: %s.INSTALLMENT_NO %q: %w", TabelListInstallment, no.String, err)
			}
		}
		if c.Persen, err = bacaDesimal(id, TabelListInstallment+".INSTALLMENT_PERCENTAGE", &persen); err != nil {
			return k, err
		}
		k.Cicilan[nama.String] = append(k.Cicilan[nama.String], c)
	}
	if err := baris.Err(); err != nil {
		return k, fmt.Errorf("repository: %s: %w", TabelListInstallment, err)
	}
	return k, nil
}

// TulisPayment - lihat PenyimpanPayment. Urutan: sentuh T_WORK_POLIS (kunci + case ada), pastikan General, jumlah cicilan,
// CurrencyList menurut nama (`simpanCurrencyList`), lalu per mata uang Total Payment / Net Premium dan cicilan diganti utuh.
func (r *PaymentOracle) TulisPayment(ctx context.Context, tx *db.Tx, id string, k models.SimpanPayment) error {
	t, err := r.tabel()
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, sqlSentuhCase(t.work), "menyentuh "+TabelWorkPolis, id, LiniFacIn)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrKasusTidakAda
	}
	if _, err := jalankan(ctx, tx, sqlPastikanGeneral(t.general), "memastikan "+TabelGeneralPolis, id, id); err != nil {
		return err
	}
	if _, err := jalankan(ctx, tx, sqlUbahInstallment(t.general), "mengubah Installment", k.Installment, id); err != nil {
		return err
	}
	baris := make([]models.BarisDeduction, 0, len(k.MataUang))
	for _, m := range k.MataUang {
		baris = append(baris, m.Baris)
	}
	ids, err := simpanCurrencyList(ctx, r.db, tx, t.cur, t.inst, id, baris, false)
	if err != nil {
		return err
	}
	for _, m := range k.MataUang {
		idCur := ids[m.Baris.Currency]
		if _, err := jalankan(ctx, tx, sqlUbahTotalPayment(t.cur), "mengubah Total Payment", ikatDesimal(m.SumTotalPayment),
			ikatDesimal(m.NetPremium), idCur); err != nil {
			return err
		}
		if _, err := jalankan(ctx, tx, sqlHapusCicilan(t.inst), "menghapus "+TabelListInstallment, idCur); err != nil {
			return err
		}
		kode := m.Baris.Currency
		if kode == "" {
			kode = mataUangKosong
		}
		for i, c := range m.Cicilan {
			persen := ""
			if c.Persen != nil {
				persen = utils.FormatDecimal(c.Persen)
			}
			if len(persen) > lebarPersenCicilan {
				return fmt.Errorf("repository: persen cicilan %s ke-%d melebihi %d bita", m.Baris.Currency, i+1, lebarPersenCicilan)
			}
			urut, err := r.db.NomorBerikut(ctx, tx, sequenceListInstallment)
			if err != nil {
				return err
			}
			uid, err := uidAcak()
			if err != nil {
				return err
			}
			nilai := []*apd.Decimal{c.Premium, c.RICommision, c.Deduction2, c.BrokerageFee, c.PPh, c.PPN, c.Discount, c.PaymentTotal}
			arg := []any{urut, idCur, i + 1, uid, kode, c.InstallmentNo, db.KosongJadiNil(c.DueDate), db.KosongJadiNil(persen)}
			for _, v := range nilai {
				arg = append(arg, ikatDesimal(v))
			}
			if _, err := jalankan(ctx, tx, sqlSisipCicilan(t.inst), "menyisip "+TabelListInstallment, arg...); err != nil {
				return err
			}
		}
	}
	return nil
}
