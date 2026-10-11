//go:build db

package repository_test

// Jalur Save — bagian yang dapat diuji TANPA mengikat apa pun (diukur
// 7 Oktober 2026). `SimpanKontrak` sendiri mengikat transaksinya, jadi tidak
// dijalankan di sini: nol Commit di uji.

import (
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
)

func TestSaveSeluruhTulisanSahDiOracleLaluDibatalkan(t *testing.T) {
	g, ctx := gudangBaca(t)
	doc := map[string]any{
		"ProportionType": "Proportional", "TreatyContractName": "UJI SAVE — DIBATALKAN",
		"Commencement": "20260101", "Termination": "20261231", "TreatyYear": "2026",
		"Ceding": "ASURANSI UJI", "LeadingReinsSource": "DIRECT", "Position": "ReasTreatyInAdmin",
		"Portfolio":   []any{map[string]any{"TypePortfolio": "premium", "Description": "satu"}},
		"CommentList": []any{map[string]any{"Date": "20261007T000000.000 GMT", "OperatorName": "UJI", "IsApproved": "Accept", "Suggest": "uji"}},
		"Limits": []any{map[string]any{"TreatyType": "QUOTA SHARE", "Detail": []any{
			map[string]any{"TreatyGroup": "PROPERTY", "QSPct": "40", "IOOLimitList": []any{map[string]any{"Currency": "IDR", "Value": "1000"}}},
		}}},
	}
	kurs := &models.KursSimpan{Tahun: "2026", Baris: []models.BarisKursWarisan{
		{MataUang: "UJI", MataUangID: "99999", NilaiKeIDR: "1", BerlakuDari: "20260101T000000.000 GMT", BerlakuSampai: "20261231T000000.000 GMT"},
	}}
	id, cacah, err := g.SimpanLaluBatalkanUntukUji(ctx, models.RencanaSimpan{Dokumen: doc, Kurs: kurs, Operator: "UJI", Stempel: "20261007T000000.000 GMT"})
	if err != nil {
		t.Fatal(err)
	}
	if id == "" || cacah["T_TREATY_REVISION"] != 1 || cacah["T_TREATY_PORTFOLIO"] != 1 || cacah["T_VIEW_COMMENT"] != 1 ||
		cacah["T_TREATY_LIMITS"] != 1 || cacah["T_TREATY_LIMIT_DETAIL"] != 1 || cacah["T_TREATY_LIMIT_AMOUNT"] != 1 {
		t.Errorf("pengenal %q, cacah %v", id, cacah)
	}
	// Dibatalkan — nol baris tersisa.
	if k, ada, _ := g.BacaKepalaTreatyIn(ctx, id); ada {
		t.Errorf("kepala %s tertinggal: %v", id, k)
	}
}

// ⭐ Migrasi 448 boleh belum terpasang: Save tetap berhasil, dan properti
// yang kolom/tabelnya belum ada DILAPORKAN. Sesudah terpasang, keduanya
// tersimpan dan tidak lagi dilaporkan — uji ini benar di kedua keadaan.
func TestSaveTolerantTerhadapMigrasi448(t *testing.T) {
	g, ctx := gudangBaca(t)
	doc := map[string]any{
		"ProportionType": "NonProportional", "TreatyContractName": "UJI 448 — DIBATALKAN",
		"RNMShareP": "10", "RevisionState": "1",
		"TotalInstallmentNP":    []any{map[string]any{"Currency": "IDR", "CurrencyID": "10026", "Value": "5"}},
		"LimitShareSummaryList": []any{map[string]any{"Layer": "1", "NetPremi": "7"}},
	}
	belum, err := g.KunciBelumTerpasang(ctx, doc)
	if err != nil {
		t.Fatal(err)
	}
	_, cacah, err := g.SimpanLaluBatalkanUntukUji(ctx, models.RencanaSimpan{Dokumen: doc})
	if err != nil {
		t.Fatalf("Save gagal sebelum/ sesudah 448: %v", err)
	}
	terpasang := !strings.Contains(strings.Join(belum, ","), "RNMShareP")
	if terpasang {
		if cacah["T_TREATY_SHARE_SUMMARY"] != 1 {
			t.Errorf("448 terpasang, ringkasan tidak tersisip: %v", cacah)
		}
	} else {
		// Belum terpasang — dilaporkan, tidak ditelan.
		for _, k := range []string{"RNMShareP", "RevisionState", "LimitShareSummaryList"} {
			if !strings.Contains(strings.Join(belum, ","), k) {
				t.Errorf("%s tidak dilaporkan: %v", k, belum)
			}
		}
	}
	// Total tidak butuh DDL — selalu tersisip ke T_TREATY_TOTAL.
	if cacah["T_TREATY_TOTAL"] != 1 {
		t.Errorf("TotalInstallmentNP tidak tersisip: %v", cacah)
	}
	t.Logf("448 terpasang: %v; dilaporkan: %v", terpasang, belum)
}

func TestSaveKontrakLamaMemperbaruiTreatyInLaluDibatalkan(t *testing.T) {
	g, ctx := gudangBaca(t)
	kepala, ada, err := g.BacaKepalaTreatyIn(ctx, "1000506")
	if err != nil || !ada {
		t.Fatalf("kepala %v %v", ada, err)
	}
	asli := kepala["TreatyContractName"]
	doc := map[string]any{}
	for k, v := range kepala {
		doc[k] = v
	}
	doc["TreatyContractName"] = "UJI UBAH — DIBATALKAN"
	id, cacah, err := g.SimpanLaluBatalkanUntukUji(ctx, models.RencanaSimpan{ID: "1000506", Dokumen: doc})
	if err != nil {
		t.Fatal(err)
	}
	if id != "1000506" || cacah["T_TREATY_REVISION"] != 1 {
		t.Errorf("pengenal %q cacah %v", id, cacah)
	}
	sesudah, _, _ := g.BacaKepalaTreatyIn(ctx, "1000506")
	if sesudah["TreatyContractName"] != asli {
		t.Errorf("TREATY_IN berubah sesudah rollback: %v → %v", asli, sesudah["TreatyContractName"])
	}
}

func TestPengenalKontrakBaruSesudahYangTertinggi(t *testing.T) {
	g, ctx := gudangBaca(t)
	id, err := g.IDKontrakBaruUntukUji(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// Situs `1` + 6 digit; tertinggi terukur 1002304 (T_TREATY_REVISION kosong).
	if len(id) != 7 || !strings.HasPrefix(id, "1") || id <= "1002304" {
		t.Errorf("pengenal baru %q", id)
	}
}

func TestPemegangPosisiDariKelolaUser(t *testing.T) {
	g, ctx := gudangBaca(t)
	admin, err := g.PemegangPosisi(ctx, "ReasTreatyInAdmin")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(admin, ","), "ADESAMUEL") {
		t.Errorf("pemegang Admin %v", admin)
	}
	kosong, err := g.PemegangPosisi(ctx, "BUKAN_WORKBASKET")
	if err != nil || len(kosong) != 0 {
		t.Errorf("workbasket tak ada: %v %v", kosong, err)
	}
}

func TestKepalaDanDokumenKontrakTersimpan(t *testing.T) {
	g, ctx := gudangBaca(t)
	k, ada, err := g.BacaKepalaTreatyIn(ctx, "1000506")
	if err != nil || !ada || k["ProportionType"] != "NonProportional" || k["TreatyYear"] != "2021" {
		t.Errorf("kepala %v ada %v err %v", k, ada, err)
	}
	// Tabel pendaratan sengaja dikosongkan sampai Save jadi: dokumennya kosong,
	// BUKAN galat.
	d, err := g.BacaDokumenPendaratan(ctx, "1000506")
	if err != nil {
		t.Fatal(err)
	}
	_ = d
}

// Daftar kontrak membawa `POSITION` — syarat tampil tombol Revision.
func TestDaftarKontrakMembawaPosisi(t *testing.T) {
	g, ctx := gudangBaca(t)
	baris, err := g.DaftarKontrakWarisan(ctx, 0, 2000)
	if err != nil {
		t.Fatal(err)
	}
	var tuntasTanpaPosisi int
	for _, b := range baris {
		if b.StatusAkseptasi == models.StatusTuntas && b.Posisi == "" {
			tuntasTanpaPosisi++
		}
	}
	if tuntasTanpaPosisi == 0 {
		t.Errorf("nol kontrak Resolve Complete tanpa posisi dari %d baris", len(baris))
	}
	t.Logf("%d baris, %d boleh direvisi", len(baris), tuntasTanpaPosisi)
}

// ⭐ Migrasi 449 — rincian Detail tab Share Prop. Sebelum terpasang: Save
// tetap berhasil, ketiga larik nilai SUDAH tersimpan (`T_TREATY_LIMIT_AMOUNT`,
// tanpa DDL), dan `Currency`/`NusaReLimit`/`SpreadingList` DILAPORKAN.
// Sesudah terpasang: semuanya tersimpan, nol laporan. Benar di kedua keadaan.
func TestSaveTolerantTerhadapMigrasi449(t *testing.T) {
	g, ctx := gudangBaca(t)
	doc := map[string]any{
		"ProportionType": "Proportional", "TreatyContractName": "UJI 449 — DIBATALKAN",
		"Limits": []any{map[string]any{"TreatyType": "QUOTA SHARE", "Detail": []any{map[string]any{
			"TreatyGroup": "PROPERTY", "TreatyGroupID": "10007", "RNMShare": "25", "SpreadingTypeID": "10252",
			"Currency": "IDR", "NusaReLimit": "15000000000",
			"RNMShareList":      []any{map[string]any{"Currency": "IDR", "Value": "15000000000"}},
			"RNMSpreadedList":   []any{map[string]any{"Currency": "IDR", "Value": "6000000000"}},
			"RNMSpreadedListRI": []any{map[string]any{"Currency": "IDR", "Value": "9000000000"}},
			"SpreadingList": []any{
				map[string]any{"ReinsTypeName": "QS (R/I)", "ReinsTypeID": "10004", "ParentReinsTypeID": "10252", "Pct": "60", "Value": "9000000000"},
				map[string]any{"ReinsTypeName": "QS (OR)", "ReinsTypeID": "10028", "ParentReinsTypeID": "10252", "Pct": "40", "Value": "6000000000"},
			},
		}}}},
	}
	belum, err := g.KunciBelumTerpasang(ctx, doc)
	if err != nil {
		t.Fatal(err)
	}
	_, cacah, err := g.SimpanLaluBatalkanUntukUji(ctx, models.RencanaSimpan{Dokumen: doc})
	if err != nil {
		t.Fatalf("Save gagal: %v", err)
	}
	// Tanpa DDL: tiga larik nilai selalu tersisip.
	if cacah["T_TREATY_LIMIT_AMOUNT"] != 3 {
		t.Errorf("RNMShareList/RNMSpreadedList(RI) tidak tersisip: %v", cacah)
	}
	terpasang := !strings.Contains(strings.Join(belum, ","), "SpreadingList")
	if terpasang {
		if cacah["T_TREATY_LIMIT_SPREADING"] != 2 || len(belum) != 0 {
			t.Errorf("449 terpasang: cacah %v laporan %v", cacah, belum)
		}
	} else {
		for _, k := range []string{"Currency", "NusaReLimit", "SpreadingList"} {
			if !strings.Contains(strings.Join(belum, ","), k) {
				t.Errorf("%s tidak dilaporkan: %v", k, belum)
			}
		}
	}
	t.Logf("449 terpasang: %v; dilaporkan: %v", terpasang, belum)
}

// Pembaca pohon Limits TOLERAN — kontrak tersimpan tetap terbuka sebelum 449.
func TestPohonLimitsTerbacaSebelumDanSesudah449(t *testing.T) {
	g, ctx := gudangBaca(t)
	pohon, err := g.BacaPohonLimitsPendaratan(ctx, "1002305")
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range pohon {
		det, _ := l["Detail"].([]map[string]any)
		for _, d := range det {
			if _, ada := d["SpreadingList"]; !ada {
				t.Errorf("Detail tanpa larik SpreadingList: %v", d["TreatyGroup"])
			}
			t.Logf("Detail %v: RNMShareList %v · SpreadingTypeID %v", d["TreatyGroup"], d["RNMShareList"], d["SpreadingTypeID"])
		}
	}
}

// ⭐ Migrasi 450 — rincian Limits. Enam larik nilai TANPA DDL selalu tersisip;
// Deduction/Parameter Achievement/Reinstatement dilaporkan sampai 450 terpasang.
func TestSaveTolerantTerhadapMigrasi450(t *testing.T) {
	g, ctx := gudangBaca(t)
	nilai := func(v string) []any { return []any{map[string]any{"Currency": "IDR", "Value": v}} }
	doc := map[string]any{
		"ProportionType": "Proportional", "TreatyContractName": "UJI 450 — DIBATALKAN",
		"Limits": []any{map[string]any{
			"TreatyType": "QUOTA SHARE", "MDPMinList": nilai("5"),
			"Reinstatement_List": []any{map[string]any{"ID": "1", "ReinstatementPct": "100", "ReinstatementValue": "1"}},
			"Detail": []any{map[string]any{
				"TreatyGroup": "PROPERTY", "ReserveList": nilai("1"), "PLAList": nilai("2"), "CashLossList": nilai("3"),
				"ClaimCoopList": nilai("4"), "DeductionTotalList": nilai("6"),
				"DeductionList": []any{map[string]any{"Comment": "Brokerage fee", "CurrencyID": "10026", "Currency": "IDR", "DeductionPct": "2.5", "DeductionPctCalculate": "true"}},
				"CurrencyList":  []any{map[string]any{"Parameter": "Based on Nett", "AchievementPctGross": "1", "LossRatioGross": "0"}},
			}},
		}},
	}
	belum, err := g.KunciBelumTerpasang(ctx, doc)
	if err != nil {
		t.Fatal(err)
	}
	_, cacah, err := g.SimpanLaluBatalkanUntukUji(ctx, models.RencanaSimpan{Dokumen: doc})
	if err != nil {
		t.Fatalf("Save gagal: %v", err)
	}
	if cacah["T_TREATY_LIMIT_AMOUNT"] != 5 || cacah["T_TREATY_LIMIT_MEASURE"] != 1 {
		t.Errorf("larik nilai tanpa DDL tidak tersisip: %v", cacah)
	}
	terpasang := !strings.Contains(strings.Join(belum, ","), "DeductionList")
	if terpasang {
		if cacah["T_TREATY_LIMIT_DEDUCTION"] != 1 || cacah["T_TREATY_LIMIT_ACH_PARAM"] != 1 || cacah["T_TREATY_LIMIT_REINSTATEMENT"] != 1 || len(belum) != 0 {
			t.Errorf("450 terpasang: cacah %v laporan %v", cacah, belum)
		}
	} else {
		for _, k := range []string{"DeductionList", "CurrencyList", "Reinstatement_List"} {
			if !strings.Contains(strings.Join(belum, ","), k) {
				t.Errorf("%s tidak dilaporkan: %v", k, belum)
			}
		}
	}
	t.Logf("450 terpasang: %v; dilaporkan: %v", terpasang, belum)
}
