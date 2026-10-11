package repository

// Tab Object + Coverage MARINE CARGO (tiket 74 P1 + P2): `.OfferFacIn.CargoList` -> T_CARGOLIST (induk T_GENERAL_POLIS) dengan
// T_SHIP / T_LC / T_SURVEYAGENT (satu halaman per cargo), T_ADDITIONALSHIP (induk T_SHIP), coverage T_COVERAGELIST ber-PARENT_TABLE
// 'T_CARGOLIST' + T_DEDUCTIBLELIST (+ DESCRIPTIONS) + T_ADDITIONALCOVERAGE + T_SPREADINGLIST (migrasi 256; kolom lihat
// docs/STRUKTUR-TABEL-NB-FACIN.md). Ganti = hapus seluruh pohon cargo case lalu sisip ulang PER BARIS dengan bind skalar (pola
// spreading FIRE yang terbukti di DEV, bug 10-10-2026), setiap baris lewat periksaBindBaris lebih dulu. Kolom DIBANGUN dari tabel
// kolom (kolomdata.go).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// Tabel migrasi 256.
	TabelCargoList         = "T_CARGOLIST"
	TabelShip256           = "T_SHIP"
	TabelLC                = "T_LC"
	TabelSurveyAgent       = "T_SURVEYAGENT"
	TabelAdditionalShip    = "T_ADDITIONALSHIP"
	TabelAdditionalCoverag = "T_ADDITIONALCOVERAGE"

	sequenceCargoList          = "SEQ_T_CARGOLIST"
	sequenceShip               = "SEQ_T_SHIP"
	sequenceLC                 = "SEQ_T_LC"
	sequenceSurveyAgent        = "SEQ_T_SURVEYAGENT"
	sequenceAdditionalShip     = "SEQ_T_ADDITIONALSHIP"
	sequenceAdditionalCoverage = "SEQ_T_ADDITIONALCOVERAGE"

	// indukCoverageCargo / jalur* - PARENT_TABLE / SRC_PATH (jalur rancangan loader `CargoList/...`).
	indukCoverageCargo   = TabelCargoList
	jalurCoverageCargo   = "CargoList/CoverageList"
	jalurDeductibleCargo = jalurCoverageCargo + "/DeductibleList"
	jalurSpreadingCargo  = jalurCoverageCargo + "/SpreadingList"
)

// PenyimpanObjekCargo - baca / ganti tab Object + Coverage MARINE CARGO case.
type PenyimpanObjekCargo interface {
	// BacaObjekCargo - cargo case urut SEQ_NO (selalu larik); ErrKasusTidakAda bila case tidak ada / bukan LINI Fac In.
	BacaObjekCargo(ctx context.Context, id string) ([]models.ObjekCargo, error)
	// GantiObjekCargo - ganti seluruh pohon cargo case di transaksi pemanggil (sentuh case + pastikan General).
	GantiObjekCargo(ctx context.Context, tx *db.Tx, id string, cargo []models.ObjekCargo) error
	// PercentShareKasus - T_GENERAL_POLIS.PERCENT_SHARE (aneka.go); tidak ada -> nil.
	PercentShareKasus(ctx context.Context, id string) (*apd.Decimal, error)
	// MataUangCargo - SetCurrencyCoverage_Act: CURRENCY berdasarkan nama -> ID / OLDID (ada = false bila tidak ada).
	MataUangCargo(ctx context.Context, nama string) (models.MataUangCsv, bool, error)
	// UbahPercentShare - T_GENERAL_POLIS.PERCENT_SHARE di transaksi pemanggil (aneka.go; Save tab Spreading cargo, tiket 74 P3).
	UbahPercentShare(ctx context.Context, tx *db.Tx, id string, ps *apd.Decimal) error
}

type tabelCargo struct {
	work, general, cargo, ship, lc, agen, kapalTambahan, cov, ded, addCov, spread string
}

func (r *ObjekOracle) tabelCargo() (tabelCargo, error) { return tabelCargoDari(r.db.Qualify) }

// tabelCargoDari - nama tabel lewat q (Qualify skema; uji: nama apa adanya).
func tabelCargoDari(q func(string) (string, error)) (tabelCargo, error) {
	var t tabelCargo
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelWorkPolis, &t.work}, {TabelGeneralPolis, &t.general}, {TabelCargoList, &t.cargo}, {TabelShip256, &t.ship},
		{TabelLC, &t.lc}, {TabelSurveyAgent, &t.agen}, {TabelAdditionalShip, &t.kapalTambahan}, {TabelCoverageList, &t.cov},
		{TabelDeductibleList, &t.ded}, {TabelAdditionalCoverag, &t.addCov}, {TabelSpreadingList, &t.spread}} {
		n, err := q(x.nama)
		if err != nil {
			return t, err
		}
		*x.ke = n
	}
	return t, nil
}

// ---- tabel kolom ----------------------------------------------------------------------------------------------------------

// kolomCargoData - kolom data T_CARGOLIST sesudah ID / PARENT_ID / SEQ_NO / ROW_UID / IS_TOP_RISK (bool, terpisah).
var kolomCargoData = []kolomData[models.ObjekCargo]{
	kt("BL_NUMBER", func(c *models.ObjekCargo) *string { return &c.BLNumber }),
	kt("CONVEYANCE_ID", func(c *models.ObjekCargo) *string { return &c.ConveyanceID }),
	kt("CONVEYANCE_NAME", func(c *models.ObjekCargo) *string { return &c.ConveyanceName }),
	kt("CONVEYANCE_NOTE", func(c *models.ObjekCargo) *string { return &c.ConveyanceNote }),
	kt("FLAG_OLD_DATA", func(c *models.ObjekCargo) *string { return &c.FlagOldData }),
	kt("FROM_RUTE", func(c *models.ObjekCargo) *string { return &c.FromRute }),
	kt("GOOD_ID", func(c *models.ObjekCargo) *string { return &c.GoodID }),
	kt("GOOD_NOTE", func(c *models.ObjekCargo) *string { return &c.GoodNote }),
	kt("INVOICE_NUMBER", func(c *models.ObjekCargo) *string { return &c.InvoiceNumber }),
	kt("PACKING_ID", func(c *models.ObjekCargo) *string { return &c.PackingID }),
	kt("PACKING_NOTE", func(c *models.ObjekCargo) *string { return &c.PackingNote }),
	kt("SAIL_DATE", func(c *models.ObjekCargo) *string { return &c.SailDate }),
	kt("TO_RUTE", func(c *models.ObjekCargo) *string { return &c.ToRute }),
	kt("TRADING_ID", func(c *models.ObjekCargo) *string { return &c.TradingID }),
	kt("TRADING_NOTE", func(c *models.ObjekCargo) *string { return &c.TradingNote }),
	kt("INVOICE_DATE", func(c *models.ObjekCargo) *string { return &c.InvoiceDate }),
	kt("MAX_AGE", func(c *models.ObjekCargo) *string { return &c.MaxAge }),
	kt("MAX_AGE_DETAIL", func(c *models.ObjekCargo) *string { return &c.MaxAgeDetail }),
	kt("FLAG_DELETE", func(c *models.ObjekCargo) *string { return &c.FlagDelete }),
}

// kolomShip - kolom data T_SHIP (sesudah ID / PARENT_ID).
var kolomShip = []kolomData[models.KapalCargo]{
	kt("AGE", func(s *models.KapalCargo) *string { return &s.Age }),
	kt("CONST", func(s *models.KapalCargo) *string { return &s.Construction }),
	kt("DWT", func(s *models.KapalCargo) *string { return &s.DWT }),
	kt("FLAG_NAME", func(s *models.KapalCargo) *string { return &s.FlagName }),
	kt("GRT", func(s *models.KapalCargo) *string { return &s.GRT }),
	kt("IMO", func(s *models.KapalCargo) *string { return &s.Imo }),
	kt("NM_SHIP", func(s *models.KapalCargo) *string { return &s.NmShip }),
	kt("NRT", func(s *models.KapalCargo) *string { return &s.NRT }),
	kt("REGISTER", func(s *models.KapalCargo) *string { return &s.Register }),
	kt("REMARK", func(s *models.KapalCargo) *string { return &s.Remark }),
	kt("SHIPTYPE_ID", func(s *models.KapalCargo) *string { return &s.ShipTypeID }),
	kt("SHIPTYPE_NAME", func(s *models.KapalCargo) *string { return &s.ShipTypeName }),
	kt("SHIP_REF_ID", func(s *models.KapalCargo) *string { return &s.ID }),
	kt("Y_MAKE1", func(s *models.KapalCargo) *string { return &s.YMake1 }),
	kt("Y_MAKE2", func(s *models.KapalCargo) *string { return &s.YMake2 }),
	kt("FORMER_SHIP_NAME", func(s *models.KapalCargo) *string { return &s.FormerShipName }),
	kt("CLASSIFICATION", func(s *models.KapalCargo) *string { return &s.Classification }),
	kt("REBUILT", func(s *models.KapalCargo) *string { return &s.Rebuilt }),
}

// kolomLC - kolom data T_LC selain LC (bool, terpisah).
var kolomLC = []kolomData[models.LCCargo]{
	kt("LCBANK", func(l *models.LCCargo) *string { return &l.Bank }),
	kt("LC_NUMBER", func(l *models.LCCargo) *string { return &l.Number }),
	kt("LC_CONDITION", func(l *models.LCCargo) *string { return &l.Condition }),
	kt("LC_REMARK", func(l *models.LCCargo) *string { return &l.Remark }),
	kt("LC_END_DATE", func(l *models.LCCargo) *string { return &l.EndDate }),
}

var kolomAgen = []kolomData[models.AgenSurvei]{
	kt("ASM_AGENT_ID", func(a *models.AgenSurvei) *string { return &a.AgentID }),
	kt("FULL_NAME", func(a *models.AgenSurvei) *string { return &a.FullName }),
	kt("ADDRESS", func(a *models.AgenSurvei) *string { return &a.Address }),
	kt("CITY", func(a *models.AgenSurvei) *string { return &a.City }),
	kt("COUNTRY", func(a *models.AgenSurvei) *string { return &a.Country }),
}

var kolomKapalTambahan = []kolomData[models.KapalTambahan]{
	kt("SHIP_REF_ID", func(s *models.KapalTambahan) *string { return &s.ID }),
	kt("REGISTER", func(s *models.KapalTambahan) *string { return &s.Register }),
	kt("NM_SHIP", func(s *models.KapalTambahan) *string { return &s.NmShip }),
	kt("VOYAGE_NO", func(s *models.KapalTambahan) *string { return &s.VoyageNo }),
	kt("SHIPTYPE_ID", func(s *models.KapalTambahan) *string { return &s.ShipTypeID }),
	kt("FLAG_NAME", func(s *models.KapalTambahan) *string { return &s.FlagName }),
	kt("NRT", func(s *models.KapalTambahan) *string { return &s.NRT }),
	kt("DWT", func(s *models.KapalTambahan) *string { return &s.DWT }),
	kt("GRT", func(s *models.KapalTambahan) *string { return &s.GRT }),
	kt("Y_MAKE1", func(s *models.KapalTambahan) *string { return &s.YMake1 }),
	kt("CONST", func(s *models.KapalTambahan) *string { return &s.Construction }),
}

// kolomCoverageCargoData - kolom data T_COVERAGELIST coverage cargo (193 / 198 / 214 / 256); CURRENCY_CODE terpisah (A169).
var kolomCoverageCargoData = []kolomData[models.CoverageCargo]{
	kt("COVERAGE", func(c *models.CoverageCargo) *string { return &c.Coverage }),
	kt("COVERAGE_NOTE", func(c *models.CoverageCargo) *string { return &c.CoverageNote }),
	kt("COVERAGE_INITIAL", func(c *models.CoverageCargo) *string { return &c.CoverageInitial }),
	ka("RATE", jenisAngka, func(c *models.CoverageCargo) **apd.Decimal { return &c.Rate }),
	ka("LIMITOF_LIABILITY", jenisAngka, func(c *models.CoverageCargo) **apd.Decimal { return &c.LimitOfLiability }),
	ka("DISCOUNT_PERCENTAGE", jenisAngka, func(c *models.CoverageCargo) **apd.Decimal { return &c.DiscountPercentage }),
	ka("TSI", jenisAngka, func(c *models.CoverageCargo) **apd.Decimal { return &c.TSI }),
	ka("PREMIUM", jenisAngka, func(c *models.CoverageCargo) **apd.Decimal { return &c.Premium }),
	ka("MIN_PREMIUM", jenisAngka, func(c *models.CoverageCargo) **apd.Decimal { return &c.MinPremium }),
	ka("DISCOUNT", jenisAngka, func(c *models.CoverageCargo) **apd.Decimal { return &c.Discount }),
	ka("TSI_LIABILITY", jenisAngka, func(c *models.CoverageCargo) **apd.Decimal { return &c.TSILiability }),
	ka("TSI_NUSANTARA_RE", jenisAngka, func(c *models.CoverageCargo) **apd.Decimal { return &c.TSINusantaraRe }),
	ka("PREMI_NUSANTARA_RE", jenisAngka, func(c *models.CoverageCargo) **apd.Decimal { return &c.PremiNusantaraRe }),
	ka("PRO_RATE_PERCENT", jenisAngka, func(c *models.CoverageCargo) **apd.Decimal { return &c.ProRatePercent }),
}

// kolomDeductibleCargoData - kolomDeductibleData (deductible.go) + DESCRIPTIONS (256); CURRENCY terpisah.
var kolomDeductibleCargoData = append(angkatDeductible(kolomDeductibleData),
	kt("DESCRIPTIONS", func(d *models.DeductibleCargo) *string { return &d.Descriptions }))

// angkatDeductible - kolom Deductible -> kolom DeductibleCargo (medan tertanam).
func angkatDeductible(kolom []kolomData[models.Deductible]) []kolomData[models.DeductibleCargo] {
	hasil := make([]kolomData[models.DeductibleCargo], 0, len(kolom))
	for _, k := range kolom {
		out := kolomData[models.DeductibleCargo]{nama: k.nama, jenis: k.jenis}
		if f := k.teks; f != nil {
			out.teks = func(d *models.DeductibleCargo) *string { return f(&d.Deductible) }
		}
		if f := k.des; f != nil {
			out.des = func(d *models.DeductibleCargo) **apd.Decimal { return f(&d.Deductible) }
		}
		hasil = append(hasil, out)
	}
	return hasil
}

var kolomWarranty = []kolomData[models.WarrantyExtension]{
	kt("COVERAGE", func(w *models.WarrantyExtension) *string { return &w.Coverage }),
	kt("COVERAGE_NOTE", func(w *models.WarrantyExtension) *string { return &w.CoverageNote }),
}

// ---- SQL ------------------------------------------------------------------------------------------------------------------

// sqlCargoKasus - ID T_CARGOLIST case (:1 = ID case).
func sqlCargoKasus(t tabelCargo) string {
	return "SELECT k.ID FROM " + t.cargo + " k WHERE k.PARENT_ID = :1"
}

// sqlCoverageCargoKasus - ID T_COVERAGELIST coverage cargo case (:1 = ID case).
func sqlCoverageCargoKasus(t tabelCargo) string {
	return "SELECT v.ID FROM " + t.cov + " v WHERE v.PARENT_TABLE = '" + indukCoverageCargo + "' AND v.PARENT_ID IN (" +
		sqlCargoKasus(t) + ")"
}

// kolomBacaCargo - cargo (alias c) lalu tiga halaman anak satu-per-cargo (LEFT JOIN, UNIQUE PARENT_ID).
func kolomBacaCargo() []string {
	k := []string{"TO_CHAR(c.ID)", "c.IS_TOP_RISK", "TO_CHAR(s.ID)", "l.LC"}
	k = append(k, kolomBacaData(kolomCargoData, "c")...)
	k = append(k, kolomBacaData(kolomShip, "s")...)
	k = append(k, kolomBacaData(kolomLC, "l")...)
	return append(k, kolomBacaData(kolomAgen, "a")...)
}

func sqlBacaCargo(t tabelCargo) string {
	return "SELECT " + strings.Join(kolomBacaCargo(), ", ") + " FROM " + t.cargo + " c LEFT JOIN " + t.ship +
		" s ON s.PARENT_ID = c.ID LEFT JOIN " + t.lc + " l ON l.PARENT_ID = c.ID LEFT JOIN " + t.agen +
		" a ON a.PARENT_ID = c.ID WHERE c.PARENT_ID = :1 ORDER BY c.SEQ_NO"
}

func sqlBacaKapalTambahan(t tabelCargo) string {
	return "SELECT " + strings.Join(kolomBacaPerInduk("x", nil, kolomKapalTambahan), ", ") +
		" FROM " + t.kapalTambahan + " x WHERE x.PARENT_ID IN (SELECT s.ID FROM " + t.ship + " s WHERE s.PARENT_ID IN (" +
		sqlCargoKasus(t) + ")) ORDER BY x.PARENT_ID, x.SEQ_NO"
}

func kolomBacaCoverageCargo() []string {
	return append([]string{"TO_CHAR(v.PARENT_ID)", "TO_CHAR(v.ID)", "v.CURRENCY_CODE"}, kolomBacaData(kolomCoverageCargoData, "v")...)
}

func sqlBacaCoverageCargo(t tabelCargo) string {
	return "SELECT " + strings.Join(kolomBacaCoverageCargo(), ", ") + " FROM " + t.cov + " v WHERE v.ID IN (" +
		sqlCoverageCargoKasus(t) + ") ORDER BY v.PARENT_ID, v.SEQ_NO"
}

// kolomBacaPerInduk - TO_CHAR(alias.PARENT_ID), ekstra, lalu kolom data (SATU urutan untuk SQL dan pemindaian).
func kolomBacaPerInduk[T any](alias string, ekstra []string, kolom []kolomData[T]) []string {
	k := append([]string{"TO_CHAR(" + alias + ".PARENT_ID)"}, ekstra...)
	return append(k, kolomBacaData(kolom, alias)...)
}

func kolomBacaDeductibleCargo() []string {
	return kolomBacaPerInduk("e", []string{"e.CURRENCY"}, kolomDeductibleCargoData)
}

func sqlBacaDeductibleCargo(t tabelCargo) string {
	return "SELECT " + strings.Join(kolomBacaDeductibleCargo(), ", ") + " FROM " + t.ded + " e WHERE e.PARENT_ID IN (" +
		sqlCoverageCargoKasus(t) + ") ORDER BY e.PARENT_ID, e.SEQ_NO"
}

func sqlBacaWarranty(t tabelCargo) string {
	return "SELECT " + strings.Join(kolomBacaPerInduk("w", nil, kolomWarranty), ", ") +
		" FROM " + t.addCov + " w WHERE w.PARENT_ID IN (" + sqlCoverageCargoKasus(t) + ") ORDER BY w.PARENT_ID, w.SEQ_NO"
}

func sqlBacaSpreadingCargo(t tabelCargo) string {
	return "SELECT TO_CHAR(s.PARENT_ID), s.TREATY_TYPE, s.TREATY_NAME, " + angkaKeluar("s.SHARE_PERCENTAGE") + ", " +
		angkaKeluar("s.TSI_SPREADED") + ", " + angkaKeluar("s.TSI_GROSS_SPREADED") + ", " + angkaKeluar("s.PREMIUM_SPREADED") + ", " +
		angkaKeluar("s.CLAIM_ESTIMATION") + " FROM " + t.spread + " s WHERE s.PARENT_ID IN (" + sqlCoverageCargoKasus(t) +
		") ORDER BY s.PARENT_ID, s.SEQ_NO"
}

// sqlHapusCargo - pohon cargo case, anak sebelum induk (FK); setiap pernyataan satu bind :1 = ID case.
func sqlHapusCargo(t tabelCargo) []string {
	cov := sqlCoverageCargoKasus(t)
	cargo := sqlCargoKasus(t)
	return []string{
		"DELETE FROM " + t.spread + " WHERE PARENT_ID IN (" + cov + ")",
		"DELETE FROM " + t.ded + " WHERE PARENT_ID IN (" + cov + ")",
		"DELETE FROM " + t.addCov + " WHERE PARENT_ID IN (" + cov + ")",
		"DELETE FROM " + t.cov + " WHERE PARENT_TABLE = '" + indukCoverageCargo + "' AND PARENT_ID IN (" + cargo + ")",
		"DELETE FROM " + t.kapalTambahan + " WHERE PARENT_ID IN (SELECT s.ID FROM " + t.ship + " s WHERE s.PARENT_ID IN (" + cargo + "))",
		"DELETE FROM " + t.ship + " WHERE PARENT_ID IN (" + cargo + ")",
		"DELETE FROM " + t.lc + " WHERE PARENT_ID IN (" + cargo + ")",
		"DELETE FROM " + t.agen + " WHERE PARENT_ID IN (" + cargo + ")",
		"DELETE FROM " + t.cargo + " WHERE PARENT_ID = :1",
	}
}

func sqlSisipCargo(cargo string) string {
	nama, nilai := kolomSisipData(kolomCargoData, 6)
	return sqlSisipData(cargo, []string{"ID", "PARENT_ID", "SEQ_NO", "ROW_UID", "IS_TOP_RISK"}, nama, nilai)
}

// sqlSisipAnakCargo - halaman satu-per-induk: :1 ID, :2 PARENT_ID, [sistem tambahan], lalu data.
func sqlSisipAnakCargo[T any](tabel string, sistem []string, kolom []kolomData[T]) string {
	s := append([]string{"ID", "PARENT_ID"}, sistem...)
	nama, nilai := kolomSisipData(kolom, len(s)+1)
	return sqlSisipData(tabel, s, nama, nilai)
}

// sqlSisipBarisCargo - baris grid: :1 ID, :2 PARENT_ID, [:3 PARENT_TABLE, :4 SRC_PATH], SEQ_NO, ROW_UID, [mata uang], lalu data.
func sqlSisipBarisCargo[T any](tabel string, denganInduk bool, mataUang string, kolom []kolomData[T]) string {
	s := []string{"ID", "PARENT_ID"}
	if denganInduk {
		s = append(s, "PARENT_TABLE", "SRC_PATH")
	}
	s = append(s, "SEQ_NO", "ROW_UID")
	if mataUang != "" {
		s = append(s, mataUang)
	}
	nama, nilai := kolomSisipData(kolom, len(s)+1)
	return sqlSisipData(tabel, s, nama, nilai)
}

// ---- baca -----------------------------------------------------------------------------------------------------------------

// BacaObjekCargo - lihat PenyimpanObjekCargo.
func (r *ObjekOracle) BacaObjekCargo(ctx context.Context, id string) ([]models.ObjekCargo, error) {
	t, err := r.tabelCargo()
	if err != nil {
		return nil, err
	}
	var n int
	if err := r.db.QueryRowContext(ctx, sqlAdaKasus(t.work), id, LiniFacIn).Scan(&n); err != nil {
		return nil, fmt.Errorf("repository: cek case NB: %w", err)
	}
	if n == 0 {
		return nil, ErrKasusTidakAda
	}
	tambahan, err := bacaPerInduk(ctx, r.db, sqlBacaKapalTambahan(t), id, "x", nil, kolomKapalTambahan, TabelAdditionalShip, nil)
	if err != nil {
		return nil, err
	}
	coverage, err := r.bacaCoverageCargo(ctx, t, id)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlBacaCargo(t), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelCargoList, err)
	}
	defer baris.Close()
	kunci := kolomBacaCargo()
	hasil := []models.ObjekCargo{}
	for baris.Next() {
		teks, err := pindaiBaris(baris, kunci)
		if err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelCargoList, err)
		}
		idCargo := teks["TO_CHAR(c.ID)"].String
		label := TabelCargoList + " " + idCargo
		c, err := isiData(kolomCargoData, "c", label, teks)
		if err != nil {
			return nil, err
		}
		c.ID, c.IsTopRisk = idCargo, teks["c.IS_TOP_RISK"].String == teksBenar
		if c.Ship, err = isiData(kolomShip, "s", label, teks); err != nil {
			return nil, err
		}
		if c.Ship.AdditionalShip = tambahan[teks["TO_CHAR(s.ID)"].String]; c.Ship.AdditionalShip == nil {
			c.Ship.AdditionalShip = []models.KapalTambahan{}
		}
		if c.LC, err = isiData(kolomLC, "l", label, teks); err != nil {
			return nil, err
		}
		c.LC.LC = teks["l.LC"].String == teksBenar
		if c.SurveyAgent, err = isiData(kolomAgen, "a", label, teks); err != nil {
			return nil, err
		}
		if c.Coverages = coverage[idCargo]; c.Coverages == nil {
			c.Coverages = []models.CoverageCargo{}
		}
		hasil = append(hasil, c)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelCargoList, err)
	}
	return hasil, nil
}

// bacaPerInduk - baris grid per PARENT_ID (teks), urut kueri. Kunci pindai = TO_CHAR(alias.PARENT_ID), ekstra, lalu kolom data -
// urutan yang SAMA dengan SQL-nya (kolomBacaPerInduk); ubah (boleh nil) melengkapi baris dari kolom ekstra.
func bacaPerInduk[T any](ctx context.Context, d *db.DB, q, id, alias string, ekstra []string, kolom []kolomData[T], tabel string,
	ubah func(*T, map[string]*sql.NullString)) (map[string][]T, error) {
	baris, err := d.QueryContext(ctx, q, id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", tabel, err)
	}
	defer baris.Close()
	kunci := kolomBacaPerInduk(alias, ekstra, kolom)
	hasil := map[string][]T{}
	for baris.Next() {
		teks, err := pindaiBaris(baris, kunci)
		if err != nil {
			return nil, fmt.Errorf("repository: %s: %w", tabel, err)
		}
		induk := teks["TO_CHAR("+alias+".PARENT_ID)"].String
		v, err := isiData(kolom, alias, tabel+" induk "+induk, teks)
		if err != nil {
			return nil, err
		}
		if ubah != nil {
			ubah(&v, teks)
		}
		hasil[induk] = append(hasil[induk], v)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", tabel, err)
	}
	return hasil, nil
}

// bacaCoverageCargo - coverage cargo case per T_CARGOLIST.ID (teks) dengan deductible, Warranty Extension, dan spreading.
func (r *ObjekOracle) bacaCoverageCargo(ctx context.Context, t tabelCargo, id string) (map[string][]models.CoverageCargo, error) {
	ded, err := bacaPerInduk(ctx, r.db, sqlBacaDeductibleCargo(t), id, "e", []string{"e.CURRENCY"}, kolomDeductibleCargoData, TabelDeductibleList,
		func(d *models.DeductibleCargo, teks map[string]*sql.NullString) {
			d.Currency = mataUangBaca(teks["e.CURRENCY"].String)
		})
	if err != nil {
		return nil, err
	}
	war, err := bacaPerInduk(ctx, r.db, sqlBacaWarranty(t), id, "w", nil, kolomWarranty, TabelAdditionalCoverag, nil)
	if err != nil {
		return nil, err
	}
	spread, err := r.bacaSpreadingCargo(ctx, t, id)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlBacaCoverageCargo(t), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca coverage cargo: %w", err)
	}
	defer baris.Close()
	kunci := kolomBacaCoverageCargo()
	hasil := map[string][]models.CoverageCargo{}
	for baris.Next() {
		teks, err := pindaiBaris(baris, kunci)
		if err != nil {
			return nil, fmt.Errorf("repository: coverage cargo: %w", err)
		}
		induk, idCov := teks["TO_CHAR(v.PARENT_ID)"].String, teks["TO_CHAR(v.ID)"].String
		c, err := isiData(kolomCoverageCargoData, "v", TabelCoverageList+" induk "+induk, teks)
		if err != nil {
			return nil, err
		}
		c.Currency = mataUangBaca(teks["v.CURRENCY_CODE"].String)
		if c.Deductibles = ded[idCov]; c.Deductibles == nil {
			c.Deductibles = []models.DeductibleCargo{}
		}
		if c.WarrantyExtension = war[idCov]; c.WarrantyExtension == nil {
			c.WarrantyExtension = []models.WarrantyExtension{}
		}
		if c.Spreading = spread[idCov]; c.Spreading == nil {
			c.Spreading = []models.BarisSpreading{}
		}
		hasil[induk] = append(hasil[induk], c)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: coverage cargo: %w", err)
	}
	return hasil, nil
}

func (r *ObjekOracle) bacaSpreadingCargo(ctx context.Context, t tabelCargo, id string) (map[string][]models.BarisSpreading, error) {
	baris, err := r.db.QueryContext(ctx, sqlBacaSpreadingCargo(t), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca spreading cargo: %w", err)
	}
	defer baris.Close()
	hasil := map[string][]models.BarisSpreading{}
	for baris.Next() {
		var induk, jenis, nama, share, tsi, gross, premi, klaim sql.NullString
		if err := baris.Scan(&induk, &jenis, &nama, &share, &tsi, &gross, &premi, &klaim); err != nil {
			return nil, fmt.Errorf("repository: spreading cargo: %w", err)
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
		return nil, fmt.Errorf("repository: spreading cargo: %w", err)
	}
	return hasil, nil
}

// ---- tulis ----------------------------------------------------------------------------------------------------------------

// GantiObjekCargo - lihat PenyimpanObjekCargo. Urutan: sentuh case (kunci + 404), pastikan General (FK T_CARGOLIST), hapus
// pohon cargo case, sisip ulang per baris.
func (r *ObjekOracle) GantiObjekCargo(ctx context.Context, tx *db.Tx, id string, cargo []models.ObjekCargo) error {
	t, err := r.tabelCargo()
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
	for _, q := range sqlHapusCargo(t) {
		if _, err := jalankan(ctx, tx, q, "menghapus cargo", id); err != nil {
			return err
		}
	}
	w := &tulisCargo{t: t, nomor: func(seq string) (string, error) { return r.db.NomorBerikut(ctx, tx, seq) }, uid: uidAcak,
		sisip: func(q, tabel string, arg []any) error { return sisipSatuBaris(ctx, tx, q, tabel, arg) }}
	return w.sisipCargo(id, cargo)
}

// sisipSatuBaris - periksaBindBaris lalu INSERT skalar; tepat satu baris.
func sisipSatuBaris(ctx context.Context, tx *db.Tx, q, tabel string, arg []any) error {
	if err := periksaBindBaris(q, tabel, [][]any{arg}, 0); err != nil {
		return err
	}
	h, err := jalankan(ctx, tx, q, "menyisipkan "+tabel, arg...)
	if err != nil {
		return err
	}
	if n, err := h.RowsAffected(); err != nil {
		return fmt.Errorf("repository: %s: %w", tabel, err)
	} else if n != 1 {
		return fmt.Errorf("repository: %s: %d baris tersisip, seharusnya 1", tabel, n)
	}
	return nil
}

// tulisCargo - traversal sisip pohon cargo (nomor = NEXTVAL, sisip = satu INSERT); dapat diuji tanpa basis data.
type tulisCargo struct {
	t     tabelCargo
	nomor func(seq string) (string, error)
	uid   func() (string, error)
	sisip func(q, tabel string, arg []any) error
	// awal - SEQ_NO cargo = awal + urutan (tiket 74 P6: kelompok berikutnya unggah CSV menambah sesudah cargo yang ada); 0 = PUT.
	awal int
}

// kapalKosong - halaman Ship tanpa isi dan tanpa Additional Ship (tidak disisipkan).
func kapalKosong(s models.KapalCargo) bool {
	for _, k := range kolomShip {
		if *k.teks(&s) != "" {
			return false
		}
	}
	return len(s.AdditionalShip) == 0
}

// lcKosong - halaman LC tanpa centang dan tanpa teks.
func lcKosong(l models.LCCargo) bool { return l == models.LCCargo{} }

func (w *tulisCargo) baru(seq string) (string, string, error) {
	id, err := w.nomor(seq)
	if err != nil {
		return "", "", err
	}
	uid, err := w.uid()
	return id, uid, err
}

func (w *tulisCargo) sisipCargo(idKasus string, cargo []models.ObjekCargo) error {
	t := w.t
	for n, c := range cargo {
		idC, uid, err := w.baru(sequenceCargoList)
		if err != nil {
			return err
		}
		arg := append([]any{idC, idKasus, w.awal + n + 1, uid, teksBool(c.IsTopRisk)}, argData(kolomCargoData, c)...)
		if err := w.sisip(sqlSisipCargo(t.cargo), TabelCargoList, arg); err != nil {
			return err
		}
		if !kapalKosong(c.Ship) {
			idS, err := w.nomor(sequenceShip)
			if err != nil {
				return err
			}
			if err := w.sisip(sqlSisipAnakCargo(t.ship, nil, kolomShip), TabelShip256,
				append([]any{idS, idC}, argData(kolomShip, c.Ship)...)); err != nil {
				return err
			}
			for m, x := range c.Ship.AdditionalShip {
				idX, uidX, err := w.baru(sequenceAdditionalShip)
				if err != nil {
					return err
				}
				if err := w.sisip(sqlSisipBarisCargo(t.kapalTambahan, false, "", kolomKapalTambahan), TabelAdditionalShip,
					append([]any{idX, idS, m + 1, uidX}, argData(kolomKapalTambahan, x)...)); err != nil {
					return err
				}
			}
		}
		if !lcKosong(c.LC) {
			idL, err := w.nomor(sequenceLC)
			if err != nil {
				return err
			}
			if err := w.sisip(sqlSisipAnakCargo(t.lc, []string{"LC"}, kolomLC), TabelLC,
				append([]any{idL, idC, teksBool(c.LC.LC)}, argData(kolomLC, c.LC)...)); err != nil {
				return err
			}
		}
		if !c.SurveyAgent.Kosong() {
			idA, err := w.nomor(sequenceSurveyAgent)
			if err != nil {
				return err
			}
			if err := w.sisip(sqlSisipAnakCargo(t.agen, nil, kolomAgen), TabelSurveyAgent,
				append([]any{idA, idC}, argData(kolomAgen, c.SurveyAgent)...)); err != nil {
				return err
			}
		}
		if err := w.sisipCoverage(idC, c.Coverages); err != nil {
			return err
		}
	}
	return nil
}

func (w *tulisCargo) sisipCoverage(idCargo string, cov []models.CoverageCargo) error {
	t := w.t
	for k, c := range cov {
		idV, uid, err := w.baru(sequenceCoverageList)
		if err != nil {
			return err
		}
		arg := []any{idV, idCargo, indukCoverageCargo, jalurCoverageCargo, k + 1, uid}
		mu := ""
		if c.Currency != "" {
			mu = "CURRENCY_CODE"
			arg = append(arg, c.Currency)
		}
		arg = append(arg, argData(kolomCoverageCargoData, c)...)
		if err := w.sisip(sqlSisipBarisCargo(t.cov, true, mu, kolomCoverageCargoData), TabelCoverageList, arg); err != nil {
			return err
		}
		for d, de := range c.Deductibles {
			idD, uidD, err := w.baru(sequenceDeductibleList)
			if err != nil {
				return err
			}
			arg := []any{idD, idV, indukDeductible, jalurDeductibleCargo, d + 1, uidD}
			mu := ""
			if de.Currency != "" {
				mu = "CURRENCY"
				arg = append(arg, de.Currency)
			}
			arg = append(arg, argData(kolomDeductibleCargoData, de)...)
			if err := w.sisip(sqlSisipBarisCargo(t.ded, true, mu, kolomDeductibleCargoData), TabelDeductibleList, arg); err != nil {
				return err
			}
		}
		for m, x := range c.WarrantyExtension {
			idW, uidW, err := w.baru(sequenceAdditionalCoverage)
			if err != nil {
				return err
			}
			if err := w.sisip(sqlSisipBarisCargo(t.addCov, false, "", kolomWarranty), TabelAdditionalCoverag,
				append([]any{idW, idV, m + 1, uidW}, argData(kolomWarranty, x)...)); err != nil {
				return err
			}
		}
		mataUang := c.Currency
		if mataUang == "" {
			mataUang = mataUangKosong
		}
		for n, s := range c.Spreading {
			idB, uidB, err := w.baru(sequenceSpreadingList)
			if err != nil {
				return err
			}
			if err := w.sisip(sqlSisipSpreading(t.spread), TabelSpreadingList, []any{idB, idV, indukSpreading, jalurSpreadingCargo,
				n + 1, uidB, db.KosongJadiNil(s.TreatyType), db.KosongJadiNil(s.TreatyName), ikatDesimal(s.SharePercentage),
				ikatDesimal(s.TSISpreaded), ikatDesimal(s.TSIGrossSpreaded), ikatDesimal(s.PremiumSpreaded),
				ikatDesimal(s.ClaimEstimation), mataUang}); err != nil {
				return err
			}
		}
	}
	return nil
}

// ---- mata uang ------------------------------------------------------------------------------------------------------------

// MataUangCargo - lihat PenyimpanObjekCargo (SQL sqlMataUangCsv - tabel dan kueri yang sama dengan Save Data CSV).
func (r *ObjekOracle) MataUangCargo(ctx context.Context, nama string) (models.MataUangCsv, bool, error) {
	q, err := r.db.Qualify(TabelCurrency)
	if err != nil {
		return models.MataUangCsv{}, false, err
	}
	var id, old, cur sql.NullString
	err = r.db.QueryRowContext(ctx, sqlMataUangCsv(q), nama).Scan(&id, &old, &cur)
	if errors.Is(err, sql.ErrNoRows) {
		return models.MataUangCsv{}, false, nil
	}
	if err != nil {
		return models.MataUangCsv{}, false, fmt.Errorf("repository: baca %s: %w", TabelCurrency, err)
	}
	return models.MataUangCsv{ID: id.String, OldID: old.String, Currency: cur.String}, true, nil
}
