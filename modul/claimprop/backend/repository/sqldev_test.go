//go:build ujidev

// Uji manual SQL Claim Prop di DEV - BACA SAJA (prompt implementasi §9): SELECT bacaan dijalankan lewat metode
// aslinya dengan masukan UJI-*; setiap INSERT / UPDATE / DELETE, dan SELECT yang hanya berjalan di dalam transaksi
// tulis, DIURAI lewat DBMS_SQL.PARSE (tidak dieksekusi). Hasil dicetak sebagai cacah dan kode galat saja, nol isi baris.
//
//	go test -tags ujidev -run TestSQLDiDEV -v ./modul/claimprop/backend/repository/
//
// Konfigurasi = lingkungan `cmd/api` (ORACLE_DSN, ORACLE_SCHEMA). Galat "objek belum ada" (ORA-00942 / 00904 /
// 02289) pada T_CLAIM_*, kolom 520, atau SEQ_T_CLAIM berarti migrasi 520-533 belum dijalankan work owner - dicatat,
// tidak menggagalkan.
package repository

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimprop/backend/models"
)

var (
	polaDML      = regexp.MustCompile(`(?is)^\s*(INSERT|UPDATE|DELETE|SELECT)\b`)
	polaBelumAda = regexp.MustCompile(`ORA-(00942|00904|02289)`)
	polaORA      = regexp.MustCompile(`ORA-\d{5}`)
)

// sqlUrai - DBMS_SQL.PARSE atas satu pernyataan DML (DDL DIEKSEKUSI oleh PARSE, maka ditolak lebih dulu).
const sqlUrai = `DECLARE c INTEGER := DBMS_SQL.OPEN_CURSOR;
BEGIN
  DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
  DBMS_SQL.CLOSE_CURSOR(c);
EXCEPTION WHEN OTHERS THEN
  IF DBMS_SQL.IS_OPEN(c) THEN DBMS_SQL.CLOSE_CURSOR(c); END IF;
  RAISE;
END;`

type hasilDEV struct{ ok, belum, gagal int }

func (h *hasilDEV) catat(t *testing.T, nama string, err error) {
	t.Helper()
	switch {
	case err == nil:
		h.ok++
		t.Logf("ok     %s", nama)
	case polaBelumAda.MatchString(err.Error()):
		h.belum++
		t.Logf("BELUM  %s: %s (objek belum ada - migrasi 520-533)", nama, polaBelumAda.FindString(err.Error()))
	default:
		h.gagal++
		kode := polaORA.FindString(err.Error())
		if kode == "" {
			kode = err.Error()
		}
		t.Errorf("GAGAL  %s: %s", nama, kode)
	}
}

func TestSQLDiDEV(t *testing.T) {
	cfg, err := config.Load()
	if err != nil || !cfg.PunyaOracle() {
		t.Skip("ORACLE_DSN belum dikonfigurasi")
	}
	d, err := db.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = d.Close() }()
	ctx, batal := context.WithTimeout(context.Background(), 5*time.Minute)
	defer batal()
	g := Baru(d)
	q := func(o string) string {
		s, err := d.Qualify(o)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	var h hasilDEV

	// ---- tulisan: diurai saja
	sqlOS, osArgs, err := sqlSisipOS(q("OS_AKSEPTASI_KLAIM"), models.BarisOS{CaseID: "UJI-C", Value: "1", GrossValue: "1",
		KursValue: "1", PersenRNM: "1", EstimationDate: time.Now()}, time.Now())
	if err != nil || len(osArgs) == 0 {
		t.Fatal(err)
	}
	tulis := map[string]string{
		"SisipKasus T_WORK_CLAIM":         sqlSisipKasus(q("T_WORK_CLAIM")),
		"SisipInduk T_GENERAL_CLAIM":      sqlSisipInduk(q("T_GENERAL_CLAIM")),
		"SisipKasusLama T_WORK_CLAIM":     sqlSisipKasusLama(q("T_WORK_CLAIM")),
		"PindahTahap":                     sqlPindahTahap(q("T_WORK_CLAIM")),
		"TutupKasus":                      sqlTutupKasus(q("T_WORK_CLAIM")),
		"SentuhKasus":                     sqlSentuhKasus(q("T_WORK_CLAIM")),
		"UbahHeader T_GENERAL_CLAIM":      sqlUbahHeader(q(models.TabelHeaderKlaim.Nama)),
		"HapusAdjustment":                 sqlHapusAdjustment(q(models.TabelAdjustment.Nama)),
		"GeserNourut":                     sqlGeserNourut(q(models.TabelAdjustment.Nama)),
		"UbahAdjustment":                  sqlUbahAdjustment(q(models.TabelAdjustment.Nama)),
		"SisipBaris T_CLAIM_ADJUSTMENT":   sqlSisipBaris(q(models.TabelAdjustment.Nama), models.TabelAdjustment),
		"NomorSementara OUTTREATYINTEMP":  sqlNomorSementara(q("OUTTREATYINTEMP")),
		"NomorCFS CFSTREATYIN":            sqlNomorCFS(q("CFSTREATYIN")),
		"SisipOS OS_AKSEPTASI_KLAIM":      sqlOS,
		"SisipJSONKlaim JSON_KLAIM":       sqlSisipJSONKlaim(q("JSON_KLAIM")),
		"LogLayanan MONITORING_KLAIM_LOG": sqlLogLayanan(q("MONITORING_KLAIM_LOG")),
		"SisipKatastrofe CATASTROPHE":     sqlSisipKatastrofe(q("CATASTROPHE")),
	}
	for _, tb := range append(append([]models.Tabel{}, models.TabelAnakKlaim...), models.TabelCucuAdjustment...) {
		tulis["HapusAnak "+tb.Nama] = sqlHapusAnak(q(tb.Nama), tb.KolomInduk)
		tulis["SisipBaris "+tb.Nama] = sqlSisipBaris(q(tb.Nama), tb)
	}
	tulis["AdaJSONKlaim JSON_KLAIM (SELECT)"] = sqlAdaJSONKlaim(q("JSON_KLAIM"))
	tulis["IDAdjustment (SELECT)"] = sqlIDAdjustment(q(models.TabelAdjustment.Nama))
	tulis["NomorRiwayat T_VIEW_SUGGEST (SELECT)"] = sqlNomorRiwayat(q("T_VIEW_SUGGEST"))
	tulis["SisipRiwayat T_VIEW_SUGGEST"] = sqlSisipRiwayat(q("T_VIEW_SUGGEST"))
	tulis["BacaRiwayat T_VIEW_SUGGEST (SELECT)"] = sqlBacaRiwayat(q("T_VIEW_SUGGEST"))
	for nama, s := range tulis {
		if !polaDML.MatchString(s) {
			t.Fatalf("%s bukan DML - DBMS_SQL.PARSE mengeksekusi DDL, ditolak", nama)
		}
		_, err := d.ExecContext(ctx, sqlUrai, s)
		h.catat(t, "urai "+nama, err)
	}

	// ---- bacaan: dijalankan (SELECT), masukan UJI-*
	a := AcuanDari(g, false)
	ap := AcuanDari(g, true)
	baca := []struct {
		nama string
		f    func() error
	}{
		{"KursStandar", func() error { _, err := a.KursStandar(ctx, "UJI"); return err }},
		{"NamaMataUang", func() error { _, err := a.NamaMataUang(ctx, "UJI"); return err }},
		{"NamaJenisReasuransi", func() error { _, err := a.NamaJenisReasuransi(ctx, "UJI"); return err }},
		{"IDJenisReasuransi", func() error { _, err := a.IDJenisReasuransi(ctx, "UJI", "4"); return err }},
		{"MasterTreaty", func() error { _, _, err := a.MasterTreaty(ctx, "UJI"); return err }},
		{"DaftarMaster", func() error { _, err := a.DaftarMaster(ctx, "UJI"); return err }},
		{"BarisMasterDari", func() error { _, _, err := a.BarisMasterDari(ctx, "UJI", "UJI", "UJI"); return err }},
		{"AdaPolisMaster", func() error { _, err := a.AdaPolisMaster(ctx, "UJI"); return err }},
		{"AdaPolisRealisasi", func() error { _, err := a.AdaPolisRealisasi(ctx, "UJI"); return err }},
		{"DaftarPolis", func() error { _, err := a.DaftarPolis(ctx, "UJI", "UJI"); return err }},
		{"TreatyGroupBisnis", func() error { _, err := a.TreatyGroupBisnis(ctx, "UJI", "2026"); return err }},
		{"YearOfQuartal", func() error { _, err := a.YearOfQuartal(ctx, "UJI", "0"); return err }},
		{"Wilayah", func() error { _, err := a.Wilayah(ctx, "UJI"); return err }},
		{"AlamatKlien", func() error { _, _, err := a.AlamatKlien(ctx, "UJI"); return err }},
		{"NamaAdjuster", func() error { _, err := a.NamaAdjuster(ctx, "UJI"); return err }},
		{"DaftarAdjuster", func() error { _, err := a.DaftarAdjuster(ctx, "UJI"); return err }},
		{"MarketingPolis", func() error { _, _, err := a.MarketingPolis(ctx, "UJI"); return err }},
		{"RiwayatKlaimPolis", func() error { _, err := a.RiwayatKlaimPolis(ctx, "UJI"); return err }},
		{"KodeLamaBisnis", func() error { _, err := a.KodeLamaBisnis(ctx, "UJI"); return err }},
		{"TahunTreaty", func() error { _, err := a.TahunTreaty(ctx, "UJI", "20260101"); return err }},
		{"LimitPLA", func() error { _, _, err := a.LimitPLA(ctx, "2026", "UJI", "UJI"); return err }},
		{"DaftarRetro", func() error { _, err := a.DaftarRetro(ctx, "UJI", "2026", "UJI"); return err }},
		{"RosterKomite", func() error { _, err := a.RosterKomite(ctx, "0", "UJI"); return err }},
		{"TingkatPelaku", func() error { _, err := a.TingkatPelaku(ctx, "UJI"); return err }},
		{"LimitDirekturUtama", func() error { _, _, err := a.LimitDirekturUtama(ctx); return err }},
		{"RekeningBank", func() error { _, err := a.RekeningBank(ctx, "UJI", "UJI"); return err }},
		{"RekeningBankMataUang", func() error { _, err := a.RekeningBankMataUang(ctx, "UJI"); return err }},
		{"RekeningBankKlien", func() error { _, err := a.RekeningBankKlien(ctx, "UJI"); return err }},
		{"IDBankRekening", func() error { _, err := a.IDBankRekening(ctx, "UJI", "UJI", "UJI"); return err }},
		{"SaldoPremi", func() error { _, err := a.SaldoPremi(ctx, "UJI", "UJI"); return err }},
		{"AdaProteksiPremi", func() error { _, err := a.AdaProteksiPremi(ctx, "UJI"); return err }},
		{"StatusKasir", func() error { _, _, err := a.StatusKasir(ctx, "UJI"); return err }},
		{"StatusKonversi (produksi)", func() error { _, err := ap.StatusKonversi(ctx, "UJI"); return err }},
		{"EmailCeding (produksi)", func() error { _, err := ap.EmailCeding(ctx, "UJI"); return err }},
		{"DaftarMataUang", func() error { _, err := a.DaftarMataUang(ctx); return err }},
		{"DaftarJenisReas", func() error { _, err := a.DaftarJenisReas(ctx, "4"); return err }},
		{"DaftarProvinsi", func() error { _, err := a.DaftarProvinsi(ctx, "UJI"); return err }},
		{"DaftarSebab", func() error { _, err := a.DaftarSebab(ctx, "UJI"); return err }},
		{"DaftarKatastrofe", func() error { _, err := a.DaftarKatastrofe(ctx, "UJI"); return err }},
		{"BarisKatastrofeID", func() error { _, _, err := a.BarisKatastrofeID(ctx, "UJI"); return err }},
		{"BarisSebabID", func() error { _, _, err := a.BarisSebabID(ctx, "UJI"); return err }},
		{"RingkasanOS", func() error { _, err := a.RingkasanOS(ctx, "UJI"); return err }},
		{"NamaPelaku", func() error { _, err := a.NamaPelaku(ctx, "UJI"); return err }},
		{"Keadaan", func() error {
			_, err := g.Keadaan(ctx, nil, "UJI-CLMP")
			if errors.Is(err, ErrKasusTidakAda) {
				return nil
			}
			return err
		}},
		{"DaftarKasus", func() error {
			_, err := g.DaftarKasus(ctx, SaringanKasus{Tahap: models.TahapOutstanding, Cari: "UJI", Batas: 5})
			return err
		}},
		{"BacaHalaman", func() error {
			_, err := g.BacaHalaman(ctx, nil, "UJI-CLMP")
			if errors.Is(err, ErrKasusTidakAda) {
				return nil
			}
			return err
		}},
		{"BacaJSONKlaimLama", func() error { _, err := g.BacaJSONKlaimLama(ctx); return err }},
		{"BacaOSLama", func() error {
			rows, err := g.BacaOSLama(ctx)
			t.Logf("       BacaOSLama: %d baris", len(rows))
			return err
		}},
	}
	for _, b := range baca {
		h.catat(t, "baca "+b.nama, b.f())
	}
	t.Logf("RINGKASAN: %d ok, %d menunggu migrasi, %d gagal", h.ok, h.belum, h.gagal)
	if strings.TrimSpace(cfg.OracleSchema) == "" {
		t.Logf("ORACLE_SCHEMA kosong")
	}
}
