package repository

// Bahan hitung kolom peserta Type QR - Calculate CSV (keputusan work owner
// 05-10-2026). Aturannya murni di models/polis_hitungqr.go.
//
// Untuk apa berkas ini: membaca SEKALI per polis - kepala polis, parameter
// produk Master Product Name Life, plan produk, master rate `M_RATE_LIFE` dan
// master risk `RIRISK_LIFE` - untuk dihitung.
//
// ⛔ BACA-SAJA, kolom disebut satu per satu, berkunci ID / IDUSEDBY. Master
// rate/risk dibaca UTUH per IDUSEDBY - TANPA BatasRateProduk/BatasRiskProduk
// (batas itu milik tampilan): memotong diam-diam akan membuat rate yang ada
// terbaca "tidak ditemukan".
//
// Dibaca sesudah: polis_rincianproduk.go (tabel dan view yang sama).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/premiumlistlife/backend/models"
)

// KepalaUnggah - kolom kepala polis yang dipakai Validate / Calculate CSV:
// hitung Type QR dan batas produk (keputusan work owner 05-10-2026).
type KepalaUnggah struct {
	Type, ProductNameID, BusinessName, ProRateType, RISlipRNM string
}

func sqlKepalaUnggah(polis string) string {
	return fmt.Sprintf(`SELECT TYPE, PRODUCT_NAME_ID, BUSINESS_NAME, PRO_RATE_TYPE, RI_SLIP_RNM FROM %s WHERE ID = :1`, polis)
}

func sqlBatasProdukMPNL(tabel string) string {
	return fmt.Sprintf(`SELECT %s, %s, %s, %s FROM %s WHERE ID = :1`,
		fmt.Sprintf(db.FmtDesimal, "MINAGE"), fmt.Sprintf(db.FmtDesimal, "MAXAGE"),
		fmt.Sprintf(db.FmtDesimal, "MINSUMINSURED"), fmt.Sprintf(db.FmtDesimal, "MAXSUMINSURED"), tabel)
}

// BatasProduk membaca batas umur dan sum insured satu produk dari tabel flat
// Master Product Name Life (`SavePremiumList_Act` 6-7, keputusan work owner
// 05-10-2026: BUKAN `PRODUCTINWARD_LIFE`). ada=false bila produknya tidak ada.
func (r *Rujukan) BatasProduk(ctx context.Context, produkID string) (models.BatasProduk, bool, error) {
	tabel, err := r.db.Qualify(TabelProdukMPNL)
	if err != nil {
		return models.BatasProduk{}, false, err
	}
	q := sqlBatasProdukMPNL(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return models.BatasProduk{}, false, err
	}
	var n [4]sql.NullString
	err = r.db.QueryRowContext(ctx, q, produkID).Scan(&n[0], &n[1], &n[2], &n[3])
	if errors.Is(err, sql.ErrNoRows) {
		return models.BatasProduk{}, false, nil
	}
	if err != nil {
		return models.BatasProduk{}, false, fmt.Errorf("repository: membaca %s: %w", TabelProdukMPNL, err)
	}
	var b models.BatasProduk
	for i, t := range []**apd.Decimal{&b.MinAge, &b.MaxAge, &b.MinSumInsured, &b.MaxSumInsured} {
		v := sql.NullString{String: angkaUtuh(strings.TrimSpace(n[i].String)), Valid: n[i].Valid}
		if *t, err = desimalAtauNil(v, TabelProdukMPNL); err != nil {
			return models.BatasProduk{}, false, err
		}
	}
	return b, true, nil
}

func sqlParamProdukQR(tabel string) string {
	return fmt.Sprintf(`SELECT PRODUCTNAME, %s, %s, %s, %s, %s, RIRISKID FROM %s WHERE ID = :1`,
		fmt.Sprintf(db.FmtDesimal, "CEDINGRETENTIONNUM"), fmt.Sprintf(db.FmtDesimal, "CEDINGLIMIT"),
		fmt.Sprintf(db.FmtDesimal, "RNMSHARE"), fmt.Sprintf(db.FmtDesimal, "RICOMM"),
		fmt.Sprintf(db.FmtDesimal, "BROKERAGE"), tabel)
}

func sqlPlanProdukQR(tabel string) string {
	return fmt.Sprintf(`SELECT NAME, RIRATEID FROM %s WHERE PRODUCTID = :1 ORDER BY URUT`, tabel)
}

func sqlRateHitungQR(tabel string) string {
	return fmt.Sprintf(`SELECT ID, GENDER, AGE, CONTRACT, RATE FROM %s WHERE IDUSEDBY = :1`, tabel)
}

// sqlRiskHitungQR - RISK kini NUMBER (migrasi inti 940, keputusan work owner 08-10-2026 K2): dibaca sebagai teks
// TANPA bergantung NLS sesi maupun format driver (db.FmtDesimal, `TM9` ber-NLS eksplisit - `921.9`, `.5`), lalu diurai
// angkaMaster yang menerima koma maupun titik desimal - hasil hitung SAMA dengan teks lama "921,9".
func sqlRiskHitungQR(tabel string) string {
	return fmt.Sprintf(`SELECT ID, YEAR, CONTRACT, %s FROM %s WHERE IDUSEDBY = :1`, fmt.Sprintf(db.FmtDesimal, "RISK"), tabel)
}

// KepalaUnggah membaca Type, Product Name ID, Class of Business, Premium
// Payment Method, dan R/I SLIP RNM polis - SEKALI per Validate / Calculate CSV.
func (r *Rujukan) KepalaUnggah(ctx context.Context, polisID string) (KepalaUnggah, error) {
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return KepalaUnggah{}, err
	}
	q := sqlKepalaUnggah(polis)
	if err := db.PeriksaSQL(q); err != nil {
		return KepalaUnggah{}, err
	}
	var n [5]sql.NullString
	err = r.db.QueryRowContext(ctx, q, polisID).Scan(&n[0], &n[1], &n[2], &n[3], &n[4])
	if errors.Is(err, sql.ErrNoRows) {
		return KepalaUnggah{}, fmt.Errorf("%w: %s", ErrHeaderPolisTidakAda, polisID)
	}
	if err != nil {
		return KepalaUnggah{}, fmt.Errorf("repository: membaca kepala polis: %w", err)
	}
	t := func(i int) string { return strings.TrimSpace(n[i].String) }
	return KepalaUnggah{Type: t(0), ProductNameID: t(1), BusinessName: t(2), ProRateType: t(3), RISlipRNM: t(4)}, nil
}

// BahanHitungQR membaca parameter produk dan seluruh plan-nya. Master rate
// plan yang cocok dan master risk dibaca terpisah (RateHitungQR, RiskHitungQR).
//
// `ErrRincianProdukTidakAda` bila produknya tidak ada di M_PRODUCTNAME_LIFE.
func (r *Rujukan) BahanHitungQR(ctx context.Context, produkID string) (models.ParamProdukQR, []models.PlanProdukQR, error) {
	tabel, err := r.db.Qualify(TabelProdukMPNL)
	if err != nil {
		return models.ParamProdukQR{}, nil, err
	}
	q := sqlParamProdukQR(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return models.ParamProdukQR{}, nil, err
	}
	var n [7]sql.NullString
	err = r.db.QueryRowContext(ctx, q, produkID).Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6])
	if errors.Is(err, sql.ErrNoRows) {
		return models.ParamProdukQR{}, nil, fmt.Errorf("%w: %s", ErrRincianProdukTidakAda, produkID)
	}
	if err != nil {
		return models.ParamProdukQR{}, nil, fmt.Errorf("repository: membaca %s: %w", TabelProdukMPNL, err)
	}
	p := models.ParamProdukQR{ProductName: n[0].String, CedingRetentionNum: angkaUtuh(n[1].String),
		CedingLimit: angkaUtuh(n[2].String), RNMShare: angkaUtuh(n[3].String), RIComm: angkaUtuh(n[4].String),
		Brokerage: angkaUtuh(n[5].String), RIRiskID: strings.TrimSpace(n[6].String)}

	tPlan, err := r.db.Qualify(TabelProdukMPNLPlan)
	if err != nil {
		return models.ParamProdukQR{}, nil, err
	}
	qp := sqlPlanProdukQR(tPlan)
	if err := db.PeriksaSQL(qp); err != nil {
		return models.ParamProdukQR{}, nil, err
	}
	rows, err := r.db.QueryContext(ctx, qp, produkID)
	if err != nil {
		return models.ParamProdukQR{}, nil, fmt.Errorf("repository: membaca %s: %w", TabelProdukMPNLPlan, err)
	}
	defer func() { _ = rows.Close() }()
	var plan []models.PlanProdukQR
	for rows.Next() {
		var a, b sql.NullString
		if err := rows.Scan(&a, &b); err != nil {
			return models.ParamProdukQR{}, nil, fmt.Errorf("repository: memindai %s: %w", TabelProdukMPNLPlan, err)
		}
		plan = append(plan, models.PlanProdukQR{Name: a.String, RIRateID: b.String})
	}
	return p, plan, rows.Err()
}

// RateHitungQR membaca SELURUH baris `M_RATE_LIFE` satu IDUSEDBY (tanpa batas tampilan).
func (r *Rujukan) RateHitungQR(ctx context.Context, riRateID string) ([]models.BarisRateQR, error) {
	tabel, err := r.db.Qualify(ViewRateProduk)
	if err != nil {
		return nil, err
	}
	q := sqlRateHitungQR(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, riRateID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s: %w", ViewRateProduk, err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []models.BarisRateQR
	for rows.Next() {
		var n [5]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4]); err != nil {
			return nil, fmt.Errorf("repository: memindai %s: %w", ViewRateProduk, err)
		}
		hasil = append(hasil, models.BarisRateQR{ID: n[0].String, Gender: n[1].String, Age: n[2].String,
			Contract: n[3].String, Rate: n[4].String})
	}
	return hasil, rows.Err()
}

// RiskHitungQR membaca SELURUH baris tabel `RIRISK_LIFE` satu IDUSEDBY (tanpa batas tampilan).
func (r *Rujukan) RiskHitungQR(ctx context.Context, riRiskID string) ([]models.BarisRiskQR, error) {
	tabel, err := r.db.Qualify(ViewRiskProduk)
	if err != nil {
		return nil, err
	}
	q := sqlRiskHitungQR(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, riRiskID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s: %w", ViewRiskProduk, err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []models.BarisRiskQR
	for rows.Next() {
		var n [4]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3]); err != nil {
			return nil, fmt.Errorf("repository: memindai %s: %w", ViewRiskProduk, err)
		}
		hasil = append(hasil, models.BarisRiskQR{ID: n[0].String, Year: n[1].String, Contract: n[2].String, Risk: n[3].String})
	}
	return hasil, rows.Err()
}
