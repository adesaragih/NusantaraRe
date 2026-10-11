//go:build db

package services_test

// Ukur ulang rumus tab Limits Non-Prop terhadap SELURUH nilai yang Pega
// simpan — `DetailCalculation` (adj, mdp) dan `DetailCalculationROL`.
//
// ⛔ BACA SAJA. Nol transaksi, nol tulisan.
//
// Masukannya nilai TERSIMPAN layer itu sendiri (EgnpiTotalList + AdjRate
// untuk Premium Earned; Premium Earned + MDP% untuk MDP; Limit, EGNPI, PE
// dan kurs tahun treaty untuk ROL), hasilnya diadu dengan yang tersimpan
// secara BILANGAN.

import (
	"context"
	"database/sql"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/services"
)

type layerUkur struct {
	l     services.LayerNP
	tahun string
	pe    []services.NilaiMataUang // TERSIMPAN
	mdp   []services.NilaiMataUang
	rol   string
}

func samaBilanganNP(a, b string) bool {
	x, _, e1 := apd.NewFromString(a)
	y, _, e2 := apd.NewFromString(b)
	if e1 != nil || e2 != nil {
		return a == b
	}
	d := new(apd.Decimal)
	_, _ = apd.BaseContext.Sub(d, x, y)
	d.Abs(d)
	return d.Cmp(apd.New(1, -6)) <= 0
}

// bulat2 membulatkan setengah-ke-atas ke 2 desimal.
func bulat2(s string) string {
	x, _, err := apd.NewFromString(s)
	if err != nil {
		return s
	}
	c := apd.BaseContext.WithPrecision(34)
	c.Rounding = apd.RoundHalfUp
	_, _ = c.Quantize(x, x, -2)
	return x.Text('f')
}

func samaDaftarNP(a, b []services.NilaiMataUang) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Currency != b[i].Currency || !samaBilanganNP(a[i].Value, b[i].Value) {
			return false
		}
	}
	return true
}

func TestRumusLimitNPCocokDenganYangPegaSimpan(t *testing.T) {
	cfg, _ := config.Load()
	if !cfg.PunyaOracle() {
		t.Skip("lewati: ORACLE_DSN kosong")
	}
	d, err := db.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = d.Close() })
	ctx := context.Background()
	sk := cfg.OracleSchema

	rows, err := d.QueryContext(ctx, `SELECT l.ID, r.TREATYYEAR, l.ADJRATE, l.MDPPCT,
	      l.MDPMINPCT, l.LIMIT, l.LIMIT2, l.CURRENCYRELATION, l.ROLPCT
	    FROM `+sk+`.T_TREATY_LIMITS l JOIN `+sk+`.T_TREATY_REVISION r ON r.MASTERID = l.MASTERID
	    WHERE l.MASTERID NOT LIKE '%#LAMA' AND l.MASTERID NOT LIKE '%/%' AND r.PROPORTIONTYPE = 'NonProportional'`)
	if err != nil {
		t.Fatal(err)
	}
	layer := map[string]*layerUkur{}
	for rows.Next() {
		var id string
		var c [8]sql.NullString
		if err := rows.Scan(&id, &c[0], &c[1], &c[2], &c[3], &c[4], &c[5], &c[6], &c[7]); err != nil {
			t.Fatal(err)
		}
		u := &layerUkur{tahun: c[0].String, rol: c[7].String}
		u.l.AdjRate, u.l.MDPPct, u.l.MDPMinPct = c[1].String, c[2].String, c[3].String
		u.l.Limit, u.l.Limit2, u.l.CurrencyRelation = c[4].String, c[5].String, c[6].String
		layer[id] = u
	}
	_ = rows.Close()

	rows, err = d.QueryContext(ctx, `SELECT IDINDUK, JENIS, CURRENCY, VALUE FROM `+sk+
		`.T_TREATY_LIMIT_MEASURE ORDER BY IDINDUK, URUTAN`)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var id, jenis string
		var cs, vs sql.NullString
		if err := rows.Scan(&id, &jenis, &cs, &vs); err != nil {
			t.Fatal(err)
		}
		cur, v := cs.String, vs.String
		u := layer[id]
		if u == nil {
			continue
		}
		n := services.NilaiMataUang{Currency: cur, Value: v}
		switch jenis {
		case "EgnpiTotalList":
			u.l.EgnpiTotalList = append(u.l.EgnpiTotalList, n)
		case "PremiumEarnedList":
			u.pe = append(u.pe, n)
		case "MDPList":
			u.mdp = append(u.mdp, n)
		}
	}
	_ = rows.Close()

	kurs := map[string][]services.KursNP{}
	rows, err = d.QueryContext(ctx, `SELECT TREATYYEAR, CURRENCY, TOIDR FROM `+sk+`.TREATYEXCHANGEYEARLY`)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var ts, cs, vs sql.NullString
		if err := rows.Scan(&ts, &cs, &vs); err != nil {
			t.Fatal(err)
		}
		th, cur, v := ts.String, cs.String, vs.String
		kurs[th] = append(kurs[th], services.KursNP{Currency: cur, Conversion: v})
	}
	_ = rows.Close()

	var nPE, okPE, nMDP, okMDP, nROL, okROL, okROL2 int
	for _, u := range layer {
		if len(u.l.EgnpiTotalList) > 0 && len(u.pe) > 0 && u.l.AdjRate != "" {
			nPE++
			l := u.l
			services.DetailCalculation(&l, services.AksiNPAdj, nil)
			if samaDaftarNP(l.PremiumEarnedList, u.pe) {
				okPE++
			}
		}
		if len(u.pe) > 0 && len(u.mdp) > 0 && u.l.MDPPct != "" {
			nMDP++
			l := u.l
			l.PremiumEarnedList = u.pe
			services.DetailCalculation(&l, services.AksiNPMDP, nil)
			if samaDaftarNP(l.MDPList, u.mdp) {
				okMDP++
			}
		}
		if len(u.l.EgnpiTotalList) > 0 && len(u.pe) > 0 && u.rol != "" {
			nROL++
			l := u.l
			l.PremiumEarnedList = u.pe
			services.DetailCalculationROL(&l, kurs[u.tahun])
			if samaBilanganNP(l.ROLPct, u.rol) {
				okROL++
			}
			// ⚠️ Data tersimpan dihitung versi rule LAMA: rasio dibulatkan 4
			// desimal (ROL 2 desimal), bukan 8 seperti ekspor sekarang.
			// Dicocokkan pula pada 2 desimal untuk membuktikan bagian lain
			// rumusnya.
			if samaBilanganNP(bulat2(l.ROLPct), bulat2(u.rol)) {
				okROL2++
			}
		}
	}
	persen := func(a, b int) float64 {
		if b == 0 {
			return 0
		}
		return float64(a) * 100 / float64(b)
	}
	t.Logf("Premium Earned %d/%d (%.1f%%) · MDP %d/%d (%.1f%%) · ROL tepat %d/%d (%.1f%%) · ROL 2 desimal %d/%d (%.1f%%)",
		okPE, nPE, persen(okPE, nPE), okMDP, nMDP, persen(okMDP, nMDP), okROL, nROL, persen(okROL, nROL),
		okROL2, nROL, persen(okROL2, nROL))
	if nPE == 0 || nMDP == 0 {
		t.Fatal("nol layer terukur — pembacanya yang rusak")
	}
	// Terukur 6 Oktober 2026: 89,0% · 97,8% · ROL 2 desimal 93,5%. Lantai,
	// bukan angka tepat — data bertambah tidak boleh memecahkan uji.
	for _, c := range []struct {
		nama   string
		a, b   int
		lantai float64
	}{{"Premium Earned", okPE, nPE, 85}, {"MDP", okMDP, nMDP, 95}, {"ROL 2 desimal", okROL2, nROL, 90}} {
		if persen(c.a, c.b) < c.lantai {
			t.Errorf("%s %.1f%% di bawah lantai %.0f%%", c.nama, persen(c.a, c.b), c.lantai)
		}
	}
}
