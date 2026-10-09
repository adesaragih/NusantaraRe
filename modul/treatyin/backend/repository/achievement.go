package repository

// Sumber sub-tab Achievement tab Limits Prop — `Activity/GetAchievement.xml`.
//
//	RDB `GetAchievement`        baris produksi per kuartal (`ACHIEVEMENT`)
//	RDB `GetCurrencyToIDR_SQL`  kurs ke IDR terbaru per `IDCURRENCY`
//
// ⛔ BACA SAJA — dijaga `TestWarisanHanyaDibaca`.
//
// ⚠️ Klaim (`GetHistoryClaim_act` / `GetEstimasiClaim_act`) TIDAK dibaca:
// sumbernya `OS_AKSEPTASI_KLAIM.DATA_JSON`, dan nilai yang ditarik dari JSON
// dilarang pemilik proses (6 Oktober 2026) sampai ada keputusan untuk tabel
// klaim. Cash Call dan Estimation karena itu nol.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
)

// Tabel sumber Achievement.
const (
	TabelAchievement = "ACHIEVEMENT"
)

// BacaAchievement — RDB `GetAchievement`, apa adanya:
//
//	... FROM POOLDATA.ACHIEVEMENT WHERE SUBSTR(NOOFFER,1,7) = SUBSTR(:id,1,7)
//	ORDER BY QUARTER, QUARTERYEAR ASC
func (g *Gudang) BacaAchievement(ctx context.Context, idKontrak string) ([]models.BarisAchievement, error) {
	nama, err := g.db.Qualify(TabelAchievement)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT IDPEGA, NOPOLIS, NOOFFER, SOBNAME, TREATYGROUPNAME, TREATYTYPE, QUARTER,
		QUARTERYEAR, IDCURRENCY, CURRENCY, PREMIUM, RICOMM, BROKERAGE, NETPREMIUM, PAIDCLAIM, OUTSTANDINGCLAIM
		FROM %s WHERE SUBSTR(NOOFFER,1,7) = SUBSTR(:1,1,7) ORDER BY QUARTER ASC, QUARTERYEAR ASC`, nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, idKontrak)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s kontrak %s: %w", TabelAchievement, idKontrak, err)
	}
	defer func() { _ = rows.Close() }()
	out := []models.BarisAchievement{}
	for rows.Next() {
		var c [16]sql.NullString
		tuju := make([]any, len(c))
		for i := range c {
			tuju[i] = &c[i]
		}
		if err := rows.Scan(tuju...); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s: %w", TabelAchievement, err)
		}
		out = append(out, models.BarisAchievement{
			IDPega: c[0].String, NoPolis: c[1].String, NoOffer: c[2].String, SoBName: c[3].String,
			TreatyGroupName: c[4].String, TreatyType: c[5].String, Quarter: c[6].String, QuarterYear: c[7].String,
			IDCurrency: c[8].String, Currency: c[9].String, Premium: c[10].String, RIComm: c[11].String,
			Brokerage: c[12].String, NetPremium: c[13].String, PaidClaim: c[14].String, OutstandingClaim: c[15].String,
		})
	}
	return out, rows.Err()
}

// BacaKursKeIDR — RDB `GetCurrencyToIDR_SQL` untuk setiap `IDCURRENCY`:
//
//	SELECT TOIDR FROM (SELECT * FROM treatyexchangeyearly
//	  WHERE IDCURRENCY = :cur ORDER BY startdate DESC) WHERE rownum = 1
//
// Pengenal tanpa baris tidak muncul di peta (Pega: `pxResults(1)` kosong).
func (g *Gudang) BacaKursKeIDR(ctx context.Context, idMataUang []string) (map[string]string, error) {
	nama, err := g.db.Qualify(TabelKursTahunan)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT TOIDR FROM (SELECT TOIDR FROM %s WHERE IDCURRENCY = :1
		ORDER BY STARTDATE DESC) WHERE ROWNUM = 1`, nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	out := map[string]string{}
	for _, id := range idMataUang {
		if _, sudah := out[id]; sudah || id == "" {
			continue
		}
		var v sql.NullString
		err := g.db.QueryRowContext(ctx, q, id).Scan(&v)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("repository: membaca kurs %s: %w", id, err)
		}
		out[id] = v.String
	}
	return out, nil
}
