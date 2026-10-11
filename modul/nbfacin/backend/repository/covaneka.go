package repository

// Coverage objek Aneka (tiket 73 P2): `.AnekaList(k).CoverageList` -> T_COVERAGELIST ber-PARENT_TABLE 'T_ANEKALIST' (rancangan
// berinduk jamak, migrasi 193 + kolom Aneka 214), SRC_PATH menurut jalur objek (okupasi / Golf); deductible -> T_DEDUCTIBLELIST
// (FK ke coverage, pola tiket 45). Ikut baca / ganti tab Object Aneka (aneka.go). Kolom DIBANGUN dari satu tabel kolom.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	jalurCoverageAnekaOkupasi = jalurAnekaOkupasi + "/CoverageList"
	jalurCoverageAnekaGolf    = jalurAnekaGolf + "/CoverageList"
)

type kolomCoverageAneka = kolomData[models.CoverageAneka]

// kolomCoverageAnekaData - kolom data T_COVERAGELIST coverage Aneka (193 / 195 / 198 / 214) sesudah kolom sistem; CURRENCY_CODE
// terpisah (boleh tidak disebut -> DEFAULT, pola A169). LOST_LIMIT VARCHAR2 rancangan berisi teks desimal (jenisTeksAngka, pola FIRE).
var kolomCoverageAnekaData = []kolomCoverageAneka{
	kt("COVERAGE", func(c *models.CoverageAneka) *string { return &c.Coverage }),
	kt("OLDID", func(c *models.CoverageAneka) *string { return &c.OldID }),
	kt("COVERAGE_NOTE", func(c *models.CoverageAneka) *string { return &c.CoverageNote }),
	kt("CALCULATE_METHOD_FAC_IN", func(c *models.CoverageAneka) *string { return &c.CalculateMethodFacIn }),
	kt("CALCULATE_METHOD", func(c *models.CoverageAneka) *string { return &c.CalculateMethod }),
	ka("PRO_RATE_PERCENT", jenisAngka, func(c *models.CoverageAneka) **apd.Decimal { return &c.ProRatePercent }),
	ka("PCT_SHORT_PERIOD", jenisAngka, func(c *models.CoverageAneka) **apd.Decimal { return &c.PctShortPeriod }),
	kt("DAY", func(c *models.CoverageAneka) *string { return &c.Day }),
	ka("TSI", jenisAngka, func(c *models.CoverageAneka) **apd.Decimal { return &c.TSI }),
	ka("LOADING", jenisAngka, func(c *models.CoverageAneka) **apd.Decimal { return &c.Loading }),
	ka("LOST_LIMIT", jenisTeksAngka, func(c *models.CoverageAneka) **apd.Decimal { return &c.LostLimit }),
	ka("RATE", jenisAngka, func(c *models.CoverageAneka) **apd.Decimal { return &c.Rate }),
	kt("INDEMNITY", func(c *models.CoverageAneka) *string { return &c.Indemnity }),
	kt("UNIT", func(c *models.CoverageAneka) *string { return &c.Unit }),
	ka("INDEMNITY_PERCENTAGE", jenisAngka, func(c *models.CoverageAneka) **apd.Decimal { return &c.IndemnityPercentage }),
	ka("PREMIUM", jenisAngka, func(c *models.CoverageAneka) **apd.Decimal { return &c.Premium }),
	ka("DISCOUNT_PERCENTAGE", jenisAngka, func(c *models.CoverageAneka) **apd.Decimal { return &c.DiscountPercentage }),
	ka("DISCOUNT", jenisAngka, func(c *models.CoverageAneka) **apd.Decimal { return &c.Discount }),
	ka("LIMITOF_LIABILITY", jenisAngka, func(c *models.CoverageAneka) **apd.Decimal { return &c.LimitOfLiability }),
	ka("PCT_LO_L", jenisAngka, func(c *models.CoverageAneka) **apd.Decimal { return &c.PctLoL }),
	kt("SUB_LIMIT_NOTE", func(c *models.CoverageAneka) *string { return &c.SubLimitNote }),
	ka("MIN_PREMIUM", jenisAngka, func(c *models.CoverageAneka) **apd.Decimal { return &c.MinPremium }),
	ka("TSI_LIABILITY", jenisAngka, func(c *models.CoverageAneka) **apd.Decimal { return &c.TSILiability }),
	ka("TSI_NUSANTARA_RE", jenisAngka, func(c *models.CoverageAneka) **apd.Decimal { return &c.TSINusantaraRe }),
	ka("PREMI_NUSANTARA_RE", jenisAngka, func(c *models.CoverageAneka) **apd.Decimal { return &c.PremiNusantaraRe }),
}

// sqlCoverageAnekaKasus - ID T_COVERAGELIST coverage Aneka case (:1 dan :2 = ID case).
func sqlCoverageAnekaKasus(t tabelAneka) string {
	return "SELECT c.ID FROM " + t.cov + " c WHERE c.PARENT_TABLE = '" + indukCoverageAnek + "' AND c.PARENT_ID IN (" +
		sqlAnekaKasus(t) + ")"
}

func kolomBacaCoverageAneka() []string {
	return append([]string{"TO_CHAR(v.PARENT_ID)", "TO_CHAR(v.ID)", "v.CURRENCY_CODE"}, kolomBacaData(kolomCoverageAnekaData, "v")...)
}

func sqlBacaCoverageAneka(t tabelAneka) string {
	return "SELECT " + strings.Join(kolomBacaCoverageAneka(), ", ") + " FROM " + t.cov + " v WHERE v.ID IN (" +
		sqlCoverageAnekaKasus(t) + ") ORDER BY v.PARENT_ID, v.SEQ_NO"
}

// sqlBacaSpreadingAneka - baris spreading coverage Aneka case (:1 dan :2 = ID case), urut coverage lalu SEQ_NO.
func sqlBacaSpreadingAneka(t tabelAneka) string {
	return "SELECT TO_CHAR(s.PARENT_ID), s.TREATY_TYPE, s.TREATY_NAME, " + angkaKeluar("s.SHARE_PERCENTAGE") + ", " +
		angkaKeluar("s.TSI_SPREADED") + ", " + angkaKeluar("s.TSI_GROSS_SPREADED") + ", " + angkaKeluar("s.PREMIUM_SPREADED") + ", " +
		angkaKeluar("s.CLAIM_ESTIMATION") + " FROM " + t.spread + " s WHERE s.PARENT_ID IN (" + sqlCoverageAnekaKasus(t) +
		") ORDER BY s.PARENT_ID, s.SEQ_NO"
}

func sqlBacaDeductibleAneka(t tabelAneka) string {
	return "SELECT " + strings.Join(kolomBacaDeductible(), ", ") + " FROM " + t.ded + " e WHERE e.PARENT_ID IN (" +
		sqlCoverageAnekaKasus(t) + ") ORDER BY e.PARENT_ID, e.SEQ_NO"
}

// sqlSisipCoverageAneka - :1 ID, :2 PARENT_ID, :3 PARENT_TABLE, :4 SRC_PATH, :5 SEQ_NO, :6 ROW_UID, [:7 CURRENCY_CODE], lalu data.
func sqlSisipCoverageAneka(cov string, denganMataUang bool) string {
	sistem := []string{"ID", "PARENT_ID", "PARENT_TABLE", "SRC_PATH", "SEQ_NO", "ROW_UID"}
	if denganMataUang {
		sistem = append(sistem, "CURRENCY_CODE")
	}
	nama, nilai := kolomSisipData(kolomCoverageAnekaData, len(sistem)+1)
	return sqlSisipData(cov, sistem, nama, nilai)
}

// bacaCoverageAneka - coverage Aneka case per T_ANEKALIST.ID (teks), deductible terpasang; selalu larik.
func (r *ObjekOracle) bacaCoverageAneka(ctx context.Context, t tabelAneka, id string) (map[string][]models.CoverageAneka, error) {
	ded := map[string][]models.Deductible{}
	bd, err := r.db.QueryContext(ctx, sqlBacaDeductibleAneka(t), id, id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca deductible Aneka: %w", err)
	}
	defer bd.Close()
	kd := kolomBacaDeductible()
	for bd.Next() {
		teks, err := pindaiBaris(bd, kd)
		if err != nil {
			return nil, fmt.Errorf("repository: deductible Aneka: %w", err)
		}
		induk := teks["TO_CHAR(e.PARENT_ID)"].String
		d, err := isiData(kolomDeductibleData, "e", TabelDeductibleList+" induk "+induk, teks)
		if err != nil {
			return nil, err
		}
		d.Currency = mataUangBaca(teks["e.CURRENCY"].String)
		ded[induk] = append(ded[induk], d)
	}
	if err := bd.Err(); err != nil {
		return nil, fmt.Errorf("repository: deductible Aneka: %w", err)
	}
	bd.Close()
	spread, err := r.bacaSpreadingAneka(ctx, t, id)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlBacaCoverageAneka(t), id, id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca coverage Aneka: %w", err)
	}
	defer baris.Close()
	kunci := kolomBacaCoverageAneka()
	hasil := map[string][]models.CoverageAneka{}
	for baris.Next() {
		teks, err := pindaiBaris(baris, kunci)
		if err != nil {
			return nil, fmt.Errorf("repository: coverage Aneka: %w", err)
		}
		induk := teks["TO_CHAR(v.PARENT_ID)"].String
		c, err := isiData(kolomCoverageAnekaData, "v", TabelCoverageList+" induk "+induk, teks)
		if err != nil {
			return nil, err
		}
		c.Currency = mataUangBaca(teks["v.CURRENCY_CODE"].String)
		c.Deductibles = ded[teks["TO_CHAR(v.ID)"].String]
		if c.Deductibles == nil {
			c.Deductibles = []models.Deductible{}
		}
		if c.Spreading = spread[teks["TO_CHAR(v.ID)"].String]; c.Spreading == nil {
			c.Spreading = []models.BarisSpreading{}
		}
		hasil[induk] = append(hasil[induk], c)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: coverage Aneka: %w", err)
	}
	return hasil, nil
}

// bacaSpreadingAneka - baris spreading coverage Aneka case per T_COVERAGELIST.ID (teks).
func (r *ObjekOracle) bacaSpreadingAneka(ctx context.Context, t tabelAneka, id string) (map[string][]models.BarisSpreading, error) {
	baris, err := r.db.QueryContext(ctx, sqlBacaSpreadingAneka(t), id, id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca spreading Aneka: %w", err)
	}
	defer baris.Close()
	hasil := map[string][]models.BarisSpreading{}
	for baris.Next() {
		var induk, jenis, nama sql.NullString
		var share, tsi, gross, premi, klaim sql.NullString
		if err := baris.Scan(&induk, &jenis, &nama, &share, &tsi, &gross, &premi, &klaim); err != nil {
			return nil, fmt.Errorf("repository: spreading Aneka: %w", err)
		}
		b := models.BarisSpreading{TreatyType: jenis.String, TreatyName: nama.String}
		label := TabelSpreadingList + " induk " + induk.String
		for _, d := range []struct {
			kolom string
			v     *sql.NullString
			ke    **apd.Decimal
		}{{"SHARE_PERCENTAGE", &share, &b.SharePercentage}, {"TSI_SPREADED", &tsi, &b.TSISpreaded},
			{"TSI_GROSS_SPREADED", &gross, &b.TSIGrossSpreaded}, {"PREMIUM_SPREADED", &premi, &b.PremiumSpreaded},
			{"CLAIM_ESTIMATION", &klaim, &b.ClaimEstimation}} {
			if *d.ke, err = bacaDesimal(label, d.kolom, d.v); err != nil {
				return nil, err
			}
		}
		hasil[induk.String] = append(hasil[induk.String], b)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: spreading Aneka: %w", err)
	}
	return hasil, nil
}

// sisipCoverageAneka - coverage satu objek (SEQ_NO 1..n), deductible, dan spreading-nya; jalur = SRC_PATH coverage; mataUangObjek =
// CURRENCY_CODE baris spreading bila coverage tanpa mata uang (kosong keduanya -> mataUangKosong, pola FIRE).
func (w *tulisAneka) sisipCoverageAneka(ctx context.Context, idAneka, jalur, mataUangObjek string, cov []models.CoverageAneka) error {
	for k, c := range cov {
		idCov, err := w.p.nomor(ctx, sequenceCoverageList)
		if err != nil {
			return err
		}
		uid, err := w.uid()
		if err != nil {
			return err
		}
		arg := []any{idCov, idAneka, indukCoverageAnek, jalur, k + 1, uid}
		if c.Currency != "" {
			arg = append(arg, c.Currency)
		}
		arg = append(arg, argData(kolomCoverageAnekaData, c)...)
		if err := w.p.sisip(ctx, sqlSisipCoverageAneka(w.t.cov, c.Currency != ""), TabelCoverageList, arg); err != nil {
			return err
		}
		for d, de := range c.Deductibles {
			idDed, err := w.p.nomor(ctx, sequenceDeductibleList)
			if err != nil {
				return err
			}
			uid, err := w.uid()
			if err != nil {
				return err
			}
			arg := append([]any{idDed, idCov, indukDeductible, jalur + "/DeductibleList", d + 1, uid}, argDeductible(de)...)
			if err := w.p.sisip(ctx, sqlSisipDeductible(w.t.ded, adaMataUang(de)), TabelDeductibleList, arg); err != nil {
				return err
			}
		}
		mataUang := c.Currency
		if mataUang == "" {
			mataUang = mataUangObjek
		}
		if mataUang == "" {
			mataUang = mataUangKosong
		}
		for n, s := range c.Spreading {
			idBaris, err := w.p.nomor(ctx, sequenceSpreadingList)
			if err != nil {
				return err
			}
			uid, err := w.uid()
			if err != nil {
				return err
			}
			if err := w.p.sisip(ctx, sqlSisipSpreading(w.t.spread), TabelSpreadingList, []any{idBaris, idCov, indukSpreading,
				jalur + "/SpreadingList", n + 1, uid, db.KosongJadiNil(s.TreatyType), db.KosongJadiNil(s.TreatyName),
				ikatDesimal(s.SharePercentage), ikatDesimal(s.TSISpreaded), ikatDesimal(s.TSIGrossSpreaded), ikatDesimal(s.PremiumSpreaded),
				ikatDesimal(s.ClaimEstimation), mataUang}); err != nil {
				return err
			}
		}
	}
	return nil
}
