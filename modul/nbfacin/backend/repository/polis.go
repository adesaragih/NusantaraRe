package repository

// Menu View Policy (Facultative) - tiket 71 (docs/issues/71-view-policy.md). Cari polis = pencari tiket 70 (FollowingOracle,
// BrowseDataViewPolis1 atas FACINPRODUCTION - SQL BrowseDataViewPolis yang dipanggil GetListViewPolicy_act L2934 TIDAK ada di
// korpus; BrowseDataViewPolis1 = salinan Save-As-nya `[dugaan kuat]`). Versi = kolom relasional JSON_POLIS (GetProdKeOldData_SQL:
// NOPOLIS / PRODKE / IDPEGA, urut TGL_INPUT turun); ringkasan polis lama = FACINPRODUCTION + CEDING_FACINPRODUCTION (A71-3).
// Nol DATA_JSON (K70-1 / K71-1).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// CariPolisQQ - opsi 7 View Policy: `UPPER(QQNAME) LIKE UPPER(filter)%` (GetListViewPolicy_act 8.9 + A71-2).
const CariPolisQQ = "qqupper"

// ErrPolisTidakAda - nomor polis / versi tidak ada di JSON_POLIS maupun FACINPRODUCTION.
var ErrPolisTidakAda = errors.New("repository: polis tidak ditemukan")

// PembacaPolis - data menu View Policy.
type PembacaPolis interface {
	// VersiPolis - versi polis urut TGL_INPUT turun; tanpa baris JSON_POLIS -> satu versi "0" dari FACINPRODUCTION
	// (MIN(IDPEGA), NOENDORS kosong). ErrPolisTidakAda bila keduanya kosong.
	VersiPolis(ctx context.Context, noPolis string) ([]models.VersiPolis, error)
	// RingkasanPolis - baris FACINPRODUCTION + ceding versi itu. ErrPolisTidakAda bila versi / barisnya tidak ada.
	RingkasanPolis(ctx context.Context, noPolis, prodKe string) (models.RingkasanPolis, error)
}

// PolisOracle - PembacaPolis atas Oracle.
type PolisOracle struct{ db *db.DB }

// NewPolisOracle merakit pembaca menu View Policy.
func NewPolisOracle(d *db.DB) *PolisOracle { return &PolisOracle{db: d} }

type tabelPolis struct{ jp, prod, ceding, work string }

func (r *PolisOracle) tabel() (tabelPolis, error) {
	var t tabelPolis
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelJSONPolis, &t.jp}, {TabelFacInProduction, &t.prod}, {TabelCedingProduksi, &t.ceding}, {TabelWorkPolis, &t.work}} {
		q, err := r.db.Qualify(x.nama)
		if err != nil {
			return t, err
		}
		*x.ke = q
	}
	return t, nil
}

// ---- SQL (bind posisional menurut kemunculan, go-ora) -------------------------------------------------------------

// sqlSistemBaru - ekspresi: case ber-IDPEGA `kolom` ada di T_WORK_POLIS LINI FAC (bind :n = LINI).
func sqlSistemBaru(work, kolom, bind string) string {
	return "(SELECT COUNT(*) FROM " + work + " w WHERE w.ID = " + sqlCaseIDDariIDPega(kolom) + " AND w.LINI = " + bind + ")"
}

// sqlVersiPolis - GetProdKeOldData_SQL + IDPEGA (label versi = substring(IDPEGA, 18) Java = SUBSTR(IDPEGA, 19)). :1 LINI, :2 NOPOLIS.
func sqlVersiPolis(jp, work string) string {
	return "SELECT j.PRODKE, " + sqlCaseIDDariIDPega("j.IDPEGA") + ", " + sqlSistemBaru(work, "j.IDPEGA", ":1") + " FROM " + jp +
		" j WHERE j.NOPOLIS = :2 ORDER BY j.TGL_INPUT DESC, j.PRODKE DESC"
}

// sqlVersiProduksi - polis tanpa baris JSON_POLIS: satu versi "0" (IDPEGA terkecil ber-NOENDORS kosong). :1 LINI, :2 NOPOLIS.
func sqlVersiProduksi(prod, work string) string {
	return "SELECT " + sqlCaseIDDariIDPega("x.IDPEGA") + ", " + sqlSistemBaru(work, "x.IDPEGA", ":1") + " FROM (SELECT MIN(p.IDPEGA) AS IDPEGA FROM " +
		prod + " p WHERE p.NOPOLIS = :2 AND p.NOENDORS IS NULL) x WHERE x.IDPEGA IS NOT NULL"
}

// sqlIDPegaVersi - IDPEGA versi (NOPOLIS, PRODKE) dari JSON_POLIS; ganda -> TGL_INPUT terbaru. :1 NOPOLIS, :2 PRODKE.
func sqlIDPegaVersi(jp string) string {
	return "SELECT IDPEGA FROM " + jp + " WHERE NOPOLIS = :1 AND PRODKE = :2 ORDER BY TGL_INPUT DESC FETCH FIRST 1 ROWS ONLY"
}

// sqlAdaJSONPolis - polis punya baris JSON_POLIS.
func sqlAdaJSONPolis(jp string) string { return "SELECT COUNT(*) FROM " + jp + " WHERE NOPOLIS = :1" }

// sqlIDPegaProduksi - IDPEGA versi "0" polis tanpa JSON_POLIS (sama dengan sqlVersiProduksi). :1 NOPOLIS.
func sqlIDPegaProduksi(prod string) string {
	return "SELECT MIN(IDPEGA) FROM " + prod + " WHERE NOPOLIS = :1 AND NOENDORS IS NULL"
}

// kolomBarisPolis - kepala (13) + baris (8 teks/tanggal) + angka (5); RiSlipCeding = BrowseDataViewPolis1 (RISLIPCEDING terakhir
// menurut SPREAD_DATE se-NOPOLIS).
func sqlBarisPolis(prod string) string {
	return "SELECT " + sqlCaseIDDariIDPega("p.IDPEGA") + ", p.NOPOLIS, p.NOENDORS, p.INSUREDNAME, p.QQNAME, p.SOB, p.CEDINGCO, " +
		"TO_CHAR(p.BEGINDATE, 'DD-MM-YYYY'), TO_CHAR(p.ENDDATE, 'DD-MM-YYYY'), p.POLICYTYPE, " +
		"(SELECT MAX(B.RISLIPCEDING) KEEP (DENSE_RANK LAST ORDER BY B.SPREAD_DATE) FROM " + prod + " B WHERE B.NOPOLIS = :1), " +
		"p.GROUPNAME, p.STATUS_BUSINESS, p.NO_MASTER_POLIS, p.OBJECTNO, p.RISK_LOCATION, p.OBJECT_ITEM, p.COVERAGENOTE, " +
		"p.ACCUMULATIONCODE, p.CURR_ID, p.RATE, p.JN_REAS, " + angkaKeluar("p.TSI100_MENJADI") + ", " + angkaKeluar("p.TSI_MENJADI") +
		", " + angkaKeluar("p.PREMI_MENJADI") + ", " + angkaKeluar("p.PERCENT_SHARE_SPREADED") + ", " + angkaKeluar("p.RICOMM") +
		" FROM " + prod + " p WHERE p.IDPEGA IN (:2, :3) AND p.NOPOLIS = :4 ORDER BY p.OBJECTNO, p.IDX_OBJECT_ITEM, p.COVERAGENOTE, p.JN_REAS"
}

// sqlCedingPolis - CEDING_FACINPRODUCTION versi itu, urut kode / nama `[keputusan agent]` (tanpa kolom urutan).
func sqlCedingPolis(ceding string) string {
	return "SELECT CEDINGNAME, " + angkaKeluar("PCT_SHARE") + " FROM " + ceding + " WHERE IDPEGA IN (:1, :2) ORDER BY CEDINGID, CEDINGNAME"
}

// ---- baca ----------------------------------------------------------------------------------------------------------

// VersiPolis - lihat PembacaPolis.
func (r *PolisOracle) VersiPolis(ctx context.Context, noPolis string) ([]models.VersiPolis, error) {
	t, err := r.tabel()
	if err != nil {
		return nil, err
	}
	q := sqlVersiPolis(t.jp, t.work)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, LiniFacIn, noPolis)
	if err != nil {
		return nil, fmt.Errorf("repository: versi polis: %w", err)
	}
	defer rows.Close()
	hasil := []models.VersiPolis{}
	for rows.Next() {
		var prodKe, caseID sql.NullString
		var n int
		if err := rows.Scan(&prodKe, &caseID, &n); err != nil {
			return nil, fmt.Errorf("repository: versi polis: %w", err)
		}
		hasil = append(hasil, models.VersiPolis{ProdKe: prodKe.String, CaseID: caseID.String, SistemBaru: n > 0})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: versi polis: %w", err)
	}
	if len(hasil) > 0 {
		return hasil, nil
	}
	q = sqlVersiProduksi(t.prod, t.work)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	var caseID sql.NullString
	var n int
	err = r.db.QueryRowContext(ctx, q, LiniFacIn, noPolis).Scan(&caseID, &n)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPolisTidakAda
	}
	if err != nil {
		return nil, fmt.Errorf("repository: versi polis FACINPRODUCTION: %w", err)
	}
	return []models.VersiPolis{{ProdKe: "0", CaseID: caseID.String, SistemBaru: n > 0}}, nil
}

// idPegaVersi - IDPEGA versi: JSON_POLIS (NOPOLIS, PRODKE); polis tanpa JSON_POLIS + PRODKE "0" -> FACINPRODUCTION.
func (r *PolisOracle) idPegaVersi(ctx context.Context, t tabelPolis, noPolis, prodKe string) (string, error) {
	q := sqlIDPegaVersi(t.jp)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var id sql.NullString
	err := r.db.QueryRowContext(ctx, q, noPolis, prodKe).Scan(&id)
	if err == nil && id.String != "" {
		return id.String, nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("repository: versi polis: %w", err)
	}
	if prodKe != "0" {
		return "", ErrPolisTidakAda
	}
	q = sqlAdaJSONPolis(t.jp)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var n int
	if err := r.db.QueryRowContext(ctx, q, noPolis).Scan(&n); err != nil {
		return "", fmt.Errorf("repository: versi polis: %w", err)
	}
	if n > 0 {
		return "", ErrPolisTidakAda // polis ber-JSON_POLIS tanpa PRODKE "0"
	}
	q = sqlIDPegaProduksi(t.prod)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	if err := r.db.QueryRowContext(ctx, q, noPolis).Scan(&id); err != nil {
		return "", fmt.Errorf("repository: versi polis FACINPRODUCTION: %w", err)
	}
	if id.String == "" {
		return "", ErrPolisTidakAda
	}
	return id.String, nil
}

// RingkasanPolis - lihat PembacaPolis.
func (r *PolisOracle) RingkasanPolis(ctx context.Context, noPolis, prodKe string) (models.RingkasanPolis, error) {
	var h models.RingkasanPolis
	t, err := r.tabel()
	if err != nil {
		return h, err
	}
	idPega, err := r.idPegaVersi(ctx, t, noPolis, prodKe)
	if err != nil {
		return h, err
	}
	q := sqlBarisPolis(t.prod)
	if err := db.PeriksaSQL(q); err != nil {
		return h, err
	}
	// IDPEGA versi dicocokkan dua bentuk: JSON_POLIS lama berawalan yang dipakai ulang bersama ceding / produksi polos (idpegaproduksi.go).
	polos, awalan := DuaBentukIDPega(idPega)
	rows, err := r.db.QueryContext(ctx, q, noPolis, polos, awalan, noPolis)
	if err != nil {
		return h, fmt.Errorf("repository: ringkasan polis: %w", err)
	}
	defer rows.Close()
	h.Baris = []models.BarisPolis{}
	for rows.Next() {
		var v [27]sql.NullString
		tujuan := make([]any, len(v))
		for i := range v {
			tujuan[i] = &v[i]
		}
		if err := rows.Scan(tujuan...); err != nil {
			return h, fmt.Errorf("repository: ringkasan polis: %w", err)
		}
		if len(h.Baris) == 0 {
			h.Kepala = models.KepalaPolis{CaseID: v[0].String, PolicyNo: v[1].String, NoEndors: v[2].String, InsuredName: v[3].String,
				QQName: v[4].String, SourceOfBusiness: v[5].String, CedingCo: v[6].String, StartPeriod: v[7].String, EndPeriod: v[8].String,
				PolicyType: v[9].String, ReffNumber: v[10].String, GroupName: v[11].String, StatusBusiness: v[12].String,
				NoMasterPolis: v[13].String}
		}
		b := models.BarisPolis{ObjectNo: v[14].String, RiskLocation: v[15].String, ObjectItem: v[16].String, Coverage: v[17].String,
			AccumulationCode: v[18].String, Currency: v[19].String, Rate: v[20].String, Treaty: v[21].String}
		teks := map[string]*sql.NullString{"tsi100": &v[22], "tsi": &v[23], "premi": &v[24], "pct": &v[25], "ricomm": &v[26]}
		if err := bacaDesimalKe(idPega, teks, kolomDesimal{"tsi100", "TSI100_MENJADI", &b.TSI100}, kolomDesimal{"tsi", "TSI_MENJADI", &b.TSIRNM},
			kolomDesimal{"premi", "PREMI_MENJADI", &b.PremiRNM}, kolomDesimal{"pct", "PERCENT_SHARE_SPREADED", &b.PctShare},
			kolomDesimal{"ricomm", "RICOMM", &b.RIComm}); err != nil {
			return h, err
		}
		h.Baris = append(h.Baris, b)
	}
	if err := rows.Err(); err != nil {
		return h, fmt.Errorf("repository: ringkasan polis: %w", err)
	}
	if len(h.Baris) == 0 {
		return h, ErrPolisTidakAda
	}
	q = sqlCedingPolis(t.ceding)
	if err := db.PeriksaSQL(q); err != nil {
		return h, err
	}
	cr, err := r.db.QueryContext(ctx, q, polos, awalan)
	if err != nil {
		return h, fmt.Errorf("repository: ceding polis: %w", err)
	}
	defer cr.Close()
	h.Ceding = []models.CedingPolis{}
	for cr.Next() {
		var nama, pct sql.NullString
		if err := cr.Scan(&nama, &pct); err != nil {
			return h, fmt.Errorf("repository: ceding polis: %w", err)
		}
		c := models.CedingPolis{CedingName: nama.String}
		if c.PctShare, err = bacaDesimal(idPega, "PCT_SHARE", &pct); err != nil {
			return h, err
		}
		h.Ceding = append(h.Ceding, c)
	}
	if err := cr.Err(); err != nil {
		return h, fmt.Errorf("repository: ceding polis: %w", err)
	}
	return h, nil
}
