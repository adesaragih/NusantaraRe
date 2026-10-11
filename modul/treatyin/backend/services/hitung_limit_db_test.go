//go:build db

package services_test

// Ukur ulang `LimitCalculation` terhadap SELURUH nilai yang Pega simpan.
//
// ⛔ BACA SAJA. Nol transaksi, nol tulisan.
//
// Tiap `Detail` kontrak Treaty In (sisi `New`, tanpa `#LAMA`, tanpa
// pengenal penyesuaian `…/Rnn`) dihitung ulang dari MASUKANNYA sendiri —
// QS% dan 100% Limit untuk Quota Share; Lines dan 100% Limit Quota Share
// sepohon untuk Surplus — lalu hasilnya diadu dengan yang tersimpan, nilai
// demi nilai, secara BILANGAN (bukan teks: `400000000` dan `400000000.0000`
// sama).

import (
	"context"
	"database/sql"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/services"
)

type detailUkur struct {
	id, master    string
	d             services.DetailLimit
	ret, ces, ioo []services.NilaiMataUang // yang TERSIMPAN
}

func bacaPohonLimit(t *testing.T) (map[string][]*detailUkur, int) {
	t.Helper()
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
	saring := ` WHERE x.MASTERID NOT LIKE '%#LAMA' AND x.MASTERID NOT LIKE '%/%' `

	q := `SELECT x.ID, x.MASTERID, x.TREATYTYPE, x.TREATYGROUP, x.QSPCT,
	        x.SURPLUS, x.RETENTIONPCT, x.CESSIONPCT, x.IOOPCT
	      FROM ` + sk + `.T_TREATY_LIMIT_DETAIL x JOIN ` + sk + `.T_TREATY_LIMITS l ON l.ID = x.IDINDUK` +
		saring + ` ORDER BY x.MASTERID, l.URUTAN, x.URUTAN`
	rows, err := d.QueryContext(ctx, q)
	if err != nil {
		t.Fatal(err)
	}
	pohon := map[string][]*detailUkur{}
	perID := map[string]*detailUkur{}
	var n int
	for rows.Next() {
		var u detailUkur
		var tt, tg, qs, sp, rp, cp, ip sql.NullString
		if err := rows.Scan(&u.id, &u.master, &tt, &tg, &qs, &sp, &rp, &cp, &ip); err != nil {
			t.Fatal(err)
		}
		u.d = services.DetailLimit{TreatyType: tt.String, TreatyGroup: tg.String, QSPct: qs.String,
			Surplus: sp.String, RetentionPct: rp.String, CessionPct: cp.String, IOOPct: ip.String}
		pohon[u.master] = append(pohon[u.master], &u)
		perID[u.id] = &u
		n++
	}
	_ = rows.Close()

	// ⚠️ `NVL(x,'')` TIDAK menolong di Oracle: teks kosong ADALAH `NULL`.
	q = `SELECT x.IDINDUK, x.JENIS, x.CURRENCY, x.CURRENCYID, x.VALUE
	     FROM ` + sk + `.T_TREATY_LIMIT_AMOUNT x` + saring + ` ORDER BY x.IDINDUK, x.URUTAN`
	rows, err = d.QueryContext(ctx, q)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var induk, jenis string
		var mu, muid, v sql.NullString
		if err := rows.Scan(&induk, &jenis, &mu, &muid, &v); err != nil {
			t.Fatal(err)
		}
		u := perID[induk]
		if u == nil {
			continue
		}
		b := services.NilaiMataUang{Currency: mu.String, CurrencyID: muid.String, Value: v.String}
		switch jenis {
		case "IOOLimitList":
			u.ioo = append(u.ioo, b)
		case "RetentionList":
			u.ret = append(u.ret, b)
		case "CessionList":
			u.ces = append(u.ces, b)
		}
	}
	_ = rows.Close()
	for _, u := range perID {
		u.d.IOOLimitList, u.d.RetentionList, u.d.CessionList = u.ioo, u.ret, u.ces
	}
	return pohon, n
}

func samaBilangan(a, b string) bool {
	x, _, e1 := apd.NewFromString(a)
	y, _, e2 := apd.NewFromString(b)
	if e1 != nil || e2 != nil {
		return a == b
	}
	return x.Cmp(y) == 0
}

func samaLarik(a, b []services.NilaiMataUang) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Currency != b[i].Currency || !samaBilangan(a[i].Value, b[i].Value) {
			return false
		}
	}
	return true
}

func TestLimitCalculationCocokDenganYangPegaSimpan(t *testing.T) {
	pohon, n := bacaPohonLimit(t)
	t.Logf("detail dibaca: %d di %d kontrak", n, len(pohon))

	var qsN, qsPct, qsRet, qsCes, qsTanpaIOO int
	var spN, spPct, spAuto, spMan, spTanpaSumber int
	for _, daftar := range pohon {
		semua := make([]services.DetailLimit, 0, len(daftar))
		for _, u := range daftar {
			semua = append(semua, u.d)
		}
		for _, u := range daftar {
			switch u.d.TreatyType {
			case "QUOTA SHARE":
				qsN++
				if len(u.ioo) == 0 {
					qsTanpaIOO++
				}
				h := services.HitungLimit(services.MasukanLimit{Jenis: "qs", Otomatis: true, Detail: u.d})
				if samaBilangan(h.RetentionPct, u.d.RetentionPct) && samaBilangan(h.CessionPct, u.d.CessionPct) {
					qsPct++
				}
				if samaLarik(h.RetentionList, u.ret) {
					qsRet++
				}
				if samaLarik(h.CessionList, u.ces) {
					qsCes++
				}
			case "SURPLUS":
				spN++
				a := services.HitungLimit(services.MasukanLimit{Jenis: "surplus", Otomatis: true, Detail: u.d, Pohon: semua})
				m := services.HitungLimit(services.MasukanLimit{Jenis: "surplus", Tambah: "man", Otomatis: true, Detail: u.d, Pohon: semua})
				if samaBilangan(a.IOOPct, u.d.IOOPct) && samaBilangan(a.CessionPct, u.d.CessionPct) {
					spPct++
				}
				adaSumber := false
				for _, q := range semua {
					if q.TreatyType == "QUOTA SHARE" && q.TreatyGroup == u.d.TreatyGroup {
						adaSumber = true
					}
				}
				if !adaSumber {
					spTanpaSumber++
				}
				if samaLarik(a.IOOLimitList, u.ioo) && samaLarik(a.RetentionList, u.ret) && samaLarik(a.CessionList, u.ces) {
					spAuto++
				} else if samaLarik(m.IOOLimitList, u.ioo) && samaLarik(m.CessionList, u.ces) {
					spMan++
				}
			}
		}
	}
	persen := func(a, b int) float64 {
		if b == 0 {
			return 0
		}
		return 100 * float64(a) / float64(b)
	}
	t.Logf("QUOTA SHARE  %d detail (%d tanpa 100%% Limit)", qsN, qsTanpaIOO)
	t.Logf("   persen Retention/Cession cocok  %d (%.1f%%)", qsPct, persen(qsPct, qsN))
	t.Logf("   RetentionList cocok             %d (%.1f%%)", qsRet, persen(qsRet, qsN))
	t.Logf("   CessionList cocok               %d (%.1f%%)", qsCes, persen(qsCes, qsN))
	t.Logf("SURPLUS      %d detail (%d tanpa Quota Share bergrup sama)", spN, spTanpaSumber)
	t.Logf("   persen IOO/Cession cocok        %d (%.1f%%)", spPct, persen(spPct, spN))
	t.Logf("   ketiga larik cocok, mode OTOMATIS %d (%.1f%%)", spAuto, persen(spAuto, spN))
	t.Logf("   ketiga larik cocok, mode MANUAL   %d (%.1f%%) — dari sisanya", spMan, persen(spMan, spN))

	// ⛔ LANTAI, terukur 6 Oktober 2026 atas seluruh 2.868 detail di 1.080
	// kontrak Treaty In:
	//
	//   QS 1.468      persen 1.468 (100%) · RetentionList 1.461 · CessionList 1.461
	//   Surplus 1.345 persen 1.345 (100%) · otomatis 519 + manual 821 = 1.340
	//
	// Lantai, bukan angka persis: kontrak baru boleh bertambah. Yang menjadi
	// regresi adalah rumus yang berhenti cocok dengan yang Pega simpan.
	//
	// ⚠️ YANG UJI INI TIDAK BUKTIKAN, dan itu dinyatakan: pembulatan empat
	// tempat `@divide(…,100,4)`. Hanya 6 dari 1.468 detail QS yang persennya
	// lewat empat desimal sesudah dibagi 100 — satu-satunya kasus di mana
	// pembulatan berpengaruh — dan datanya TERBELAH: 1 cocok hanya dengan
	// pembulatan, 1 hanya tanpa, 4 tidak cocok keduanya (nilai bulat yang
	// tampaknya diketik tangan, mis. 200.000 dari 910.000 × 21,978…%).
	// Pembulatannya berdiri di atas TEKS ekspor, bukan di atas data; mode
	// setengah-ke-atasnya belum terverifikasi sama sekali.
	if qsPct < qsN || qsRet < 1450 || qsCes < 1450 {
		t.Errorf("QS menyimpang: persen %d/%d, retensi %d, cession %d (terukur 1.468 · 1.461 · 1.461)",
			qsPct, qsN, qsRet, qsCes)
	}
	if spPct < spN || spAuto+spMan < 1330 {
		t.Errorf("Surplus menyimpang: persen %d/%d, terjelaskan %d (terukur 1.345 · 1.340)",
			spPct, spN, spAuto+spMan)
	}
}
