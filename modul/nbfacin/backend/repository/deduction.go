package repository

// Tab Premium Deduction / Brokerage Fee kasus FIRE (tiket 50): persen dan centang di T_GENERAL_POLIS (PAY_RI_COMMISION,
// PAY_PCT_BROKERAGE_FEE, PPN_CHECK), potret `.OfferFacIn.CurrencyList(n).Policy.Payment` di T_CURRENCYLIST (migrasi 200),
// dan tarif pajak `POOLDATA.TAX` (`Obj-Browse` kelas ASM-FW-GISFW-Int-TAX di `SetRIComm_Act` langkah 4 /
// `TotalBrokerageFee_Act` langkah 2: urut `.TANGGAL` turun, satu baris, tanpa saringan; DDL `D:\migrasi\RNM\DDL\TAX.txt`
// `TAX (TANGGAL DATE, PPH NUMBER, PPN NUMBER)`).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelCurrencyList - CurrencyList pembayaran per mata uang (migrasi 200).
	TabelCurrencyList    = "T_CURRENCYLIST"
	sequenceCurrencyList = "SEQ_T_CURRENCYLIST"
	// TabelPajak - tarif PPH / PPN warisan (dibaca saja).
	TabelPajak = "TAX"
	// lebarNamaMataUang - T_CURRENCYLIST.CURRENCY_CODE (50; NAME 500).
	lebarNamaMataUang = 50
)

// PenyimpanDeduction - data tab Premium Deduction.
type PenyimpanDeduction interface {
	// BacaDeduction - ErrKasusTidakAda bila case tidak ada / bukan LINI Fac In.
	BacaDeduction(ctx context.Context, id string) (models.KasusDeduction, error)
	// TulisDeduction - persen, centang, dan CurrencyList diganti utuh di transaksi pemanggil. ErrKasusTidakAda bila case
	// tidak ada.
	TulisDeduction(ctx context.Context, tx *db.Tx, id string, k models.KasusDeduction) error
	// TarifPajak - baris TAX ber-TANGGAL terbaru; ada=false bila tabel kosong.
	TarifPajak(ctx context.Context) (models.TarifPajak, bool, error)
}

// DeductionOracle - PenyimpanDeduction atas Oracle.
type DeductionOracle struct{ db *db.DB }

// NewDeductionOracle merakit penyimpan tab Premium Deduction.
func NewDeductionOracle(d *db.DB) *DeductionOracle { return &DeductionOracle{db: d} }

type tabelDeduction struct{ work, general, cur string }

func (r *DeductionOracle) tabel() (tabelDeduction, error) {
	var t tabelDeduction
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelWorkPolis, &t.work}, {TabelGeneralPolis, &t.general}, {TabelCurrencyList, &t.cur}} {
		q, err := r.db.Qualify(x.nama)
		if err != nil {
			return t, err
		}
		*x.ke = q
	}
	return t, nil
}

func sqlBacaKepalaDeduction(t tabelDeduction) string {
	return "SELECT " + angkaKeluar("g.PAY_RI_COMMISION") + ", " + angkaKeluar("g.PAY_PCT_BROKERAGE_FEE") + ", g.PPN_CHECK FROM " +
		t.work + " w LEFT JOIN " + t.general + " g ON g.ID = w.ID WHERE w.ID = :1 AND w.LINI = :2"
}

func sqlBacaCurrencyList(cur string) string {
	return "SELECT NAME, " + angkaKeluar("PAY_PREMIUM") + ", " + angkaKeluar("PAY_COMMISION") + ", " + angkaKeluar("PAY_DEDUCTION2") + ", " +
		angkaKeluar("PAY_BROKERAGE_FEE") + ", " + angkaKeluar("PAY_PPH") + ", " + angkaKeluar("PAY_PPN") + " FROM " + cur +
		" WHERE PARENT_ID = :1 ORDER BY SEQ_NO"
}

func sqlUbahKepalaDeduction(general string) string {
	return "UPDATE " + general + " SET PAY_RI_COMMISION = " + angkaMasuk(":1") + ", PAY_PCT_BROKERAGE_FEE = " + angkaMasuk(":2") +
		", PPN_CHECK = :3 WHERE ID = :4"
}

func sqlSisipCurrencyList(cur string) string {
	return "INSERT INTO " + cur + " (ID, PARENT_ID, SEQ_NO, ROW_UID, CURRENCY_CODE, NAME, PAY_PREMIUM, PAY_COMMISION, PAY_DEDUCTION2," +
		" PAY_BROKERAGE_FEE, PAY_PPH, PAY_PPN) VALUES (:1, :2, :3, :4, :5, :6, " + angkaMasuk(":7") + ", " + angkaMasuk(":8") + ", " +
		angkaMasuk(":9") + ", " + angkaMasuk(":10") + ", " + angkaMasuk(":11") + ", " + angkaMasuk(":12") + ")"
}

// sqlTarifPajak - Obj-Browse TAX: Sort D pada TANGGAL, MaxRecords 1 (urutan NULL bawaan Oracle, seperti SQL Pega).
func sqlTarifPajak(tax string) string {
	return "SELECT " + angkaKeluar("PPH") + ", " + angkaKeluar("PPN") + " FROM " + tax + " ORDER BY TANGGAL DESC FETCH FIRST 1 ROWS ONLY"
}

// BacaDeduction - lihat PenyimpanDeduction.
func (r *DeductionOracle) BacaDeduction(ctx context.Context, id string) (models.KasusDeduction, error) {
	var k models.KasusDeduction
	t, err := r.tabel()
	if err != nil {
		return k, err
	}
	var ri, pct, ppn sql.NullString
	err = r.db.QueryRowContext(ctx, sqlBacaKepalaDeduction(t), id, LiniFacIn).Scan(&ri, &pct, &ppn)
	if errors.Is(err, sql.ErrNoRows) {
		return k, ErrKasusTidakAda
	}
	if err != nil {
		return k, fmt.Errorf("repository: baca kepala tab Premium Deduction: %w", err)
	}
	if k.RICommision, err = bacaDesimal(id, TabelGeneralPolis+".PAY_RI_COMMISION", &ri); err != nil {
		return k, err
	}
	if k.PctBrokerageFee, err = bacaDesimal(id, TabelGeneralPolis+".PAY_PCT_BROKERAGE_FEE", &pct); err != nil {
		return k, err
	}
	k.PPnCheck = ppn.String
	baris, err := r.db.QueryContext(ctx, sqlBacaCurrencyList(t.cur), id)
	if err != nil {
		return k, fmt.Errorf("repository: baca %s: %w", TabelCurrencyList, err)
	}
	defer baris.Close()
	k.MataUang = []models.BarisDeduction{}
	for baris.Next() {
		var nama sql.NullString
		var v [6]sql.NullString
		if err := baris.Scan(&nama, &v[0], &v[1], &v[2], &v[3], &v[4], &v[5]); err != nil {
			return k, fmt.Errorf("repository: %s: %w", TabelCurrencyList, err)
		}
		b := models.BarisDeduction{Currency: nama.String}
		for i, kol := range []struct {
			nama string
			ke   **apd.Decimal
		}{{"PAY_PREMIUM", &b.Premium}, {"PAY_COMMISION", &b.Commision}, {"PAY_DEDUCTION2", &b.Deduction2},
			{"PAY_BROKERAGE_FEE", &b.BrokerageFee}, {"PAY_PPH", &b.PPh}, {"PAY_PPN", &b.PPN}} {
			if *kol.ke, err = bacaDesimal(id, TabelCurrencyList+"."+kol.nama, &v[i]); err != nil {
				return k, err
			}
		}
		k.MataUang = append(k.MataUang, b)
	}
	if err := baris.Err(); err != nil {
		return k, fmt.Errorf("repository: %s: %w", TabelCurrencyList, err)
	}
	return k, nil
}

// TulisDeduction - lihat PenyimpanDeduction. Urutan: sentuh T_WORK_POLIS (kunci + case ada), pastikan General, persen +
// centang, CurrencyList diperbarui menurut nama mata uang (cicilan tab Payment dipertahankan, `simpanCurrencyList`).
func (r *DeductionOracle) TulisDeduction(ctx context.Context, tx *db.Tx, id string, k models.KasusDeduction) error {
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
	if _, err := jalankan(ctx, tx, sqlUbahKepalaDeduction(t.general), "mengubah persen Premium Deduction",
		ikatDesimal(k.RICommision), ikatDesimal(k.PctBrokerageFee), db.KosongJadiNil(k.PPnCheck), id); err != nil {
		return err
	}
	inst, err := r.db.Qualify(TabelListInstallment)
	if err != nil {
		return err
	}
	_, err = simpanCurrencyList(ctx, r.db, tx, t.cur, inst, id, k.MataUang, true)
	return err
}

// TarifPajak - lihat PenyimpanDeduction.
func (r *DeductionOracle) TarifPajak(ctx context.Context) (models.TarifPajak, bool, error) {
	var p models.TarifPajak
	tax, err := r.db.Qualify(TabelPajak)
	if err != nil {
		return p, false, err
	}
	var pph, ppn sql.NullString
	err = r.db.QueryRowContext(ctx, sqlTarifPajak(tax)).Scan(&pph, &ppn)
	if errors.Is(err, sql.ErrNoRows) {
		return p, false, nil
	}
	if err != nil {
		return p, false, fmt.Errorf("repository: baca %s: %w", TabelPajak, err)
	}
	if p.PPH, err = bacaDesimal("baris TANGGAL terbaru", TabelPajak+".PPH", &pph); err != nil {
		return p, false, err
	}
	if p.PPN, err = bacaDesimal("baris TANGGAL terbaru", TabelPajak+".PPN", &ppn); err != nil {
		return p, false, err
	}
	return p, true, nil
}
