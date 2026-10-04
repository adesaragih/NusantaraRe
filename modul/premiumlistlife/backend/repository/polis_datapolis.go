package repository

// Data polis layar Input Premium Detail - tiket 03 bagian 2.
//
// Untuk apa berkas ini: membaca dan menulis kolom `T_PREMIUM_LIST` milik
// tahap Input Premium Detail (Type, produk, marketing, R/I SLIP, retro), dan
// ketiga pencarian master layar itu. Seluruh kolomnya sudah ada sejak 051 -
// nol migrasi.
//
// ⛔ Hanya BARIS UTAMA (`ID` = nomor kasus): baris status penawaran (062)
// tidak memuat kolom tahap ini.
//
// ⛔ `Qualify` di setiap query (ADR-U-0033), nol `COMMIT` (ADR-U-0029).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/premiumlistlife/backend/models"
)

// DataPolisTersimpan adalah isian data polis yang sudah ada di header.
type DataPolisTersimpan struct {
	models.IsianDataPolis
	// WPC - DIBACA saja di sini; ditulis saat Confirm (`TulisWPC`, rumus
	// `models.WPCPolis` pengganti `WPCLife_Act` / `POOLDATA.GETQUARTER`,
	// 03-10-2026). Kosong sampai polisnya jadi.
	WPC *time.Time
}

// kolomDataPolis - urutan SELECT dan pengurai `BacaDataPolis`; teks kecuali
// yang disebut.
var kolomDataPolis = []string{
	"TYPE", "PRODUCT_NAME_ID", "PRODUCT_NAME", "SOB", "SOB_NAME", "CEDING_CO", "CEDING_CO_NAME",
	"POLICY_HOLDER", "POLICY_HOLDER_NAME", "RI_SLIP_RNM", "PRO_RATE_TYPE", "MO_ID", "MARKETING_CODE",
	"MARKETING_NAME", "ANNUITY_INTEREST", "PREMIUM_REFUND_FACTOR", "RETRO_ID", "RETRO_NAME",
	"SECURITY_REINSURER_ID", "SECURITY_REINSURER", "WPC", "DATE_RECEIVED",
}

func sqlBacaDataPolis(polis string) string {
	bagian := make([]string, len(kolomDataPolis))
	for i, k := range kolomDataPolis {
		switch k {
		case "ANNUITY_INTEREST", "PREMIUM_REFUND_FACTOR":
			bagian[i] = fmt.Sprintf(db.FmtDesimal, "p."+k)
		case "WPC", "DATE_RECEIVED":
			bagian[i] = fmt.Sprintf(db.FmtTanggalOracle, "p."+k)
		default:
			bagian[i] = "p." + k
		}
	}
	return fmt.Sprintf(`SELECT %s FROM %s p WHERE p.ID = :1`, strings.Join(bagian, ", "), polis)
}

// BacaDataPolis membaca data polis baris utama satu kasus.
func (r *Penawaran) BacaDataPolis(ctx context.Context, id string) (DataPolisTersimpan, error) {
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return DataPolisTersimpan{}, err
	}
	q := sqlBacaDataPolis(polis)
	if err := db.PeriksaSQL(q); err != nil {
		return DataPolisTersimpan{}, err
	}
	n := make([]sql.NullString, len(kolomDataPolis))
	tujuan := make([]any, len(n))
	for i := range n {
		tujuan[i] = &n[i]
	}
	err = r.db.QueryRowContext(ctx, q, id).Scan(tujuan...)
	if errors.Is(err, sql.ErrNoRows) {
		return DataPolisTersimpan{}, ErrHeaderPolisTidakAda
	}
	if err != nil {
		return DataPolisTersimpan{}, fmt.Errorf("repository: membaca data polis: %w", err)
	}
	return uraiDataPolis(n)
}

func uraiDataPolis(n []sql.NullString) (DataPolisTersimpan, error) {
	v := map[string]sql.NullString{}
	for i, k := range kolomDataPolis {
		v[k] = n[i]
	}
	t := func(k string) string { return v[k].String }
	h := DataPolisTersimpan{IsianDataPolis: models.IsianDataPolis{
		Type: t("TYPE"), ProductNameID: t("PRODUCT_NAME_ID"), ProductName: t("PRODUCT_NAME"),
		SourceOfBusiness: t("SOB"), SobName: t("SOB_NAME"), CedingCo: t("CEDING_CO"),
		CedingCoName: t("CEDING_CO_NAME"), PolicyHolder: t("POLICY_HOLDER"),
		PolicyHolderName: t("POLICY_HOLDER_NAME"), RISlipRNM: t("RI_SLIP_RNM"),
		ProRateType: t("PRO_RATE_TYPE"), MoID: t("MO_ID"), MarketingCode: t("MARKETING_CODE"),
		MarketingName: t("MARKETING_NAME"), RetroID: t("RETRO_ID"), RetroName: t("RETRO_NAME"),
		SecurityReinsurerID: t("SECURITY_REINSURER_ID"), SecurityReinsurer: t("SECURITY_REINSURER"),
	}}
	var err error
	for kolom, tujuan := range map[string]**apd.Decimal{
		"ANNUITY_INTEREST": &h.AnnuityInterest, "PREMIUM_REFUND_FACTOR": &h.PremiumRefundFactor,
	} {
		if s := v[kolom]; s.Valid && strings.TrimSpace(s.String) != "" {
			if *tujuan, err = utils.ParseDecimal(strings.TrimSpace(s.String)); err != nil {
				return h, fmt.Errorf("repository: kolom %s bernilai %q: %w", kolom, s.String, err)
			}
		}
	}
	if h.WPC, err = uraiTanggalOracle(v["WPC"], "WPC"); err != nil {
		return h, err
	}
	if h.DateReceived, err = uraiTanggalOracle(v["DATE_RECEIVED"], "DATE_RECEIVED"); err != nil {
		return h, err
	}
	return h, nil
}

// sqlSimpanDataPolis menimpa kolom data polis tahap Input Premium Detail.
//
// ⛔ Hanya kolom yang layar ini isi; data penawaran DITULIS ULANG hanya untuk
// pasangan Ceding / Policy Holder yang `SetProdNametoPolis` timpa dari produk.
func sqlSimpanDataPolis(polis string) string {
	return fmt.Sprintf(`UPDATE %s
	    SET TYPE = :1, PRODUCT_NAME_ID = :2, PRODUCT_NAME = :3, SOB = :4, SOB_NAME = :5,
	        CEDING_CO = :6, CEDING_CO_NAME = :7, POLICY_HOLDER = :8, POLICY_HOLDER_NAME = :9,
	        RI_SLIP_RNM = :10, PRO_RATE_TYPE = :11, MO_ID = :12, MARKETING_CODE = :13,
	        MARKETING_NAME = :14, ANNUITY_INTEREST = :15, PREMIUM_REFUND_FACTOR = :16,
	        RETRO_ID = :17, RETRO_NAME = :18, SECURITY_REINSURER_ID = :19, SECURITY_REINSURER = :20,
	        DATE_RECEIVED = :21
	  WHERE ID = :22`, polis)
}

func argSimpanDataPolis(id string, d models.IsianDataPolis) []any {
	desimal := func(x *apd.Decimal) any {
		if x == nil {
			return nil
		}
		return utils.FormatDecimal(x)
	}
	k := db.KosongJadiNil
	return []any{
		k(d.Type), k(d.ProductNameID), k(d.ProductName), k(d.SourceOfBusiness), k(d.SobName),
		k(d.CedingCo), k(d.CedingCoName), k(d.PolicyHolder), k(d.PolicyHolderName),
		k(d.RISlipRNM), k(d.ProRateType), k(d.MoID), k(d.MarketingCode), k(d.MarketingName),
		desimal(d.AnnuityInterest), desimal(d.PremiumRefundFactor),
		k(d.RetroID), k(d.RetroName), k(d.SecurityReinsurerID), k(d.SecurityReinsurer),
		tanggalAtauNil(d.DateReceived), id,
	}
}

// SimpanDataPolis menulis data polis ke baris utama - di transaksi pemanggil.
func (r *Penawaran) SimpanDataPolis(ctx context.Context, tx *db.Tx, id string, d models.IsianDataPolis) error {
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return err
	}
	q := sqlSimpanDataPolis(polis)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, argSimpanDataPolis(id, d)...)
	if err != nil {
		return fmt.Errorf("repository: menyimpan data polis: %w", err)
	}
	if n, err := hasil.RowsAffected(); err == nil && n == 0 {
		return ErrHeaderPolisTidakAda
	}
	return db.PastikanSatuBaris(hasil, "data polis")
}

// ——— Save Data: bahan SavePremiumList_Act langkah 6-8 ———

// MasterProdukInwardLife - `GetRateProductLife`: `SELECT * FROM
// POOLDATA.PRODUCTINWARD_LIFE WHERE ID = {ParamData.CARI1}` (CARI1 =
// ProductNameID). `[terverifikasi]` MINAGE/MAXAGE (GetProductDtlPL `b.MINAGE`,
// `b.MAXAGE`). ⚠️ `[belum terverifikasi]` MINSUMINSURED/MAXSUMINSURED - nama
// properti yang SavePremiumList_Act baca dari `ProductNameInward.pxResults(1)`.
const MasterProdukInwardLife = "PRODUCTINWARD_LIFE"

func sqlBatasProduk(tabel string) string {
	return fmt.Sprintf(`SELECT %s, %s, %s, %s FROM %s WHERE ID = :1 FETCH FIRST 1 ROWS ONLY`,
		fmt.Sprintf(db.FmtDesimal, "MINAGE"), fmt.Sprintf(db.FmtDesimal, "MAXAGE"),
		fmt.Sprintf(db.FmtDesimal, "MINSUMINSURED"), fmt.Sprintf(db.FmtDesimal, "MAXSUMINSURED"), tabel)
}

// sqlPesertaBatas - urutan GRID peserta (`sqlGridPeserta`): nomor "at list"
// pesan galat menunjuk baris yang pemakai lihat.
func sqlPesertaBatas(detail string) string {
	return fmt.Sprintf(`SELECT d.NAME_OF_INSURED, %s, %s FROM %s d
	  WHERE d.PREMIUM_LIST_ID = :1 ORDER BY d.ID`,
		fmt.Sprintf(db.FmtDesimal, "d.ENTRY_AGE"), fmt.Sprintf(db.FmtDesimal, "d.SUM_INSURED"), detail)
}

func desimalAtauNil(v sql.NullString, kolom string) (*apd.Decimal, error) {
	if !v.Valid || strings.TrimSpace(v.String) == "" {
		return nil, nil
	}
	d, err := utils.ParseDecimal(strings.TrimSpace(v.String))
	if err != nil {
		return nil, fmt.Errorf("repository: kolom %s bernilai %q: %w", kolom, v.String, err)
	}
	return d, nil
}

// BatasProduk membaca batas umur dan sum insured satu produk; ada=false bila
// produknya tidak punya baris di PRODUCTINWARD_LIFE.
func (r *Penawaran) BatasProduk(ctx context.Context, produkID string) (models.BatasProduk, bool, error) {
	tabel, err := r.db.Qualify(MasterProdukInwardLife)
	if err != nil {
		return models.BatasProduk{}, false, err
	}
	q := sqlBatasProduk(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return models.BatasProduk{}, false, err
	}
	var n [4]sql.NullString
	err = r.db.QueryRowContext(ctx, q, produkID).Scan(&n[0], &n[1], &n[2], &n[3])
	if errors.Is(err, sql.ErrNoRows) {
		return models.BatasProduk{}, false, nil
	}
	if err != nil {
		return models.BatasProduk{}, false, fmt.Errorf("repository: membaca %s: %w", MasterProdukInwardLife, err)
	}
	var b models.BatasProduk
	for i, t := range []**apd.Decimal{&b.MinAge, &b.MaxAge, &b.MinSumInsured, &b.MaxSumInsured} {
		if *t, err = desimalAtauNil(n[i], MasterProdukInwardLife); err != nil {
			return models.BatasProduk{}, false, err
		}
	}
	return b, true, nil
}

// PesertaBatas membaca umur masuk dan sum insured seluruh peserta satu polis.
func (r *Penawaran) PesertaBatas(ctx context.Context, polisID string) ([]models.PesertaBatas, error) {
	detail, err := r.db.Qualify("T_PREMIUM_LIST_DETAIL")
	if err != nil {
		return nil, err
	}
	q := sqlPesertaBatas(detail)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, polisID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca peserta untuk batas produk: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []models.PesertaBatas
	for rows.Next() {
		var n [3]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2]); err != nil {
			return nil, err
		}
		p := models.PesertaBatas{NameOfInsured: n[0].String}
		if p.EntryAge, err = desimalAtauNil(n[1], "ENTRY_AGE"); err != nil {
			return nil, err
		}
		if p.SumInsured, err = desimalAtauNil(n[2], "SUM_INSURED"); err != nil {
			return nil, err
		}
		hasil = append(hasil, p)
	}
	return hasil, rows.Err()
}

// tanggalAtauNil - tanggal kosong menjadi NULL.
func tanggalAtauNil(t *time.Time) any {
	if t == nil {
		return nil
	}
	return *t
}

// sqlTulisWPC mengisi WPC baris utama header satu polis.
func sqlTulisWPC(polis string) string {
	return fmt.Sprintf(`UPDATE %s SET WPC = :1 WHERE ID = :2`, polis)
}

// TulisWPC menulis WPC polis (`models.WPCPolis`) - saat Confirm, sesudah PL
// Number terbit, di transaksi pemanggil (keputusan work owner 03-10-2026,
// pengganti `WPCLife_Act` / `POOLDATA.GETQUARTER`).
func (r *Penawaran) TulisWPC(ctx context.Context, tx *db.Tx, id string, wpc time.Time) error {
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return err
	}
	q := sqlTulisWPC(polis)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, wpc, id)
	if err != nil {
		return fmt.Errorf("repository: menulis WPC: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "WPC polis")
}
