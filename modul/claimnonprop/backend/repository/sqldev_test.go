//go:build ujidev

// Uji manual SQL Claim Non Prop di DEV - BACA SAJA (prompt §8): SELECT bacaan dijalankan lewat metode aslinya dengan
// masukan UJI-*; setiap INSERT / UPDATE / DELETE, dan SELECT yang hanya berjalan di dalam transaksi tulis, DIURAI lewat
// DBMS_SQL.PARSE (tidak dieksekusi). Hasil dicetak sebagai cacah dan kode galat saja, nol isi baris.
//
//	go test -tags ujidev -run 'TestSQLDiDEV|TestHitungMasterDEV' -v ./modul/claimnonprop/backend/repository/
//
// Konfigurasi = lingkungan `cmd/api` (ORACLE_DSN, ORACLE_SCHEMA). Galat "objek belum ada" (ORA-00942 / 00904 /
// 02289) pada kolom 600-607, tabel T_CLAIM_NP_*, atau PLATNP_SEQ berarti migrasi 600-611 belum dijalankan work owner
// (atau objek DBA belum ada, OQ-CNP-39) - dicatat, tidak menggagalkan.
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
	"nusantarare/modul/claimnonprop/backend/models"
)

var (
	polaDML = regexp.MustCompile(`(?is)^\s*(INSERT|UPDATE|DELETE|SELECT)\b`)
	// ORA-04063: view CLAIMXOL INVALID di DEV (OQ-CNP-28) - objek DBA, bukan galat SQL modul.
	polaBelumAda = regexp.MustCompile(`ORA-(00942|00904|02289|04063)`)
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
		t.Logf("BELUM  %s: %s (objek belum ada - migrasi 600-611 / OQ)", nama, polaBelumAda.FindString(err.Error()))
	default:
		h.gagal++
		kode := polaORA.FindString(err.Error())
		if kode == "" {
			kode = err.Error()
		}
		t.Errorf("GAGAL  %s: %s", nama, kode)
	}
}

func bukaDEV(t *testing.T) (*db.DB, config.Config) {
	t.Helper()
	cfg, err := config.Load()
	if err != nil || !cfg.PunyaOracle() {
		t.Skip("ORACLE_DSN belum dikonfigurasi")
	}
	d, err := db.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return d, cfg
}

func TestSQLDiDEV(t *testing.T) {
	d, cfg := bukaDEV(t)
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
	tulis := map[string]string{
		"SisipKasus T_WORK_CLAIM":              sqlSisipKasus(q("T_WORK_CLAIM")),
		"SisipInduk T_GENERAL_CLAIM":           sqlSisipInduk(q("T_GENERAL_CLAIM")),
		"PindahTahap":                          sqlPindahTahap(q("T_WORK_CLAIM")),
		"TutupKasus":                           sqlTutupKasus(q("T_WORK_CLAIM")),
		"SentuhKasus":                          sqlSentuhKasus(q("T_WORK_CLAIM")),
		"UbahHeader T_GENERAL_CLAIM":           sqlUbahHeader(q(models.TabelHeaderKlaim.Nama)),
		"HapusAdjustment":                      sqlHapusAdjustment(q(models.TabelAdjustment.Nama)),
		"GeserNourut":                          sqlGeserNourut(q(models.TabelAdjustment.Nama)),
		"UbahAdjustment":                       sqlUbahAdjustment(q(models.TabelAdjustment.Nama)),
		"SisipBaris T_CLAIM_ADJUSTMENT":        sqlSisipBaris(q(models.TabelAdjustment.Nama), models.TabelAdjustment),
		"SisipOS OS_AKSEPTASI_KLAIM":           sqlSisipOS(q("OS_AKSEPTASI_KLAIM")),
		"JumlahOS (SELECT)":                    sqlJumlahOS(q("OS_AKSEPTASI_KLAIM")),
		"NomorKlaimOS (SELECT)":                sqlNomorKlaimOS(q("OS_AKSEPTASI_KLAIM")),
		"SisipJSONKlaim JSON_KLAIM":            sqlSisipJSONKlaim(q("JSON_KLAIM")),
		"AdaJSONKlaim JSON_KLAIM (SELECT)":     sqlAdaJSONKlaim(q("JSON_KLAIM")),
		"SisipKatastrofe CATASTROPHE":          sqlSisipKatastrofe(q("CATASTROPHE")),
		"IDAdjustment (SELECT)":                sqlIDAdjustment(q(models.TabelAdjustment.Nama)),
		"NomorRiwayat T_VIEW_SUGGEST (SELECT)": sqlNomorRiwayat(q("T_VIEW_SUGGEST")),
		"SisipRiwayat T_VIEW_SUGGEST":          sqlSisipRiwayat(q("T_VIEW_SUGGEST")),
		"BacaRiwayat T_VIEW_SUGGEST (SELECT)":  sqlBacaRiwayat(q("T_VIEW_SUGGEST")),
		"SisipKasusKomite T_WORK_CLAIM":        sqlSisipKasusKomite(q("T_WORK_CLAIM")),
		"SisipKepalaKomite T_GENERAL_KOMITE":   sqlSisipKepalaKomite(q("T_GENERAL_KOMITE")),
		"SisipAnggotaKomite tangga":            sqlSisipAnggotaKomite(q(tabelTanggaKomite)),
		"SetelKomiteAdjustment":                sqlSetelKomiteAdjustment(q(models.TabelAdjustment.Nama)),
		"TanggaKomite (SELECT)":                sqlTanggaKomite(q(tabelTanggaKomite), q("T_GENERAL_KOMITE"), q("T_WORK_CLAIM")),
	}
	for _, tb := range append(append([]models.Tabel{}, models.TabelAnakKlaim...), models.TabelCucuAdjustment...) {
		kunci := tb.Nama + "/" + tb.Daftar
		tulis["HapusAnak "+kunci] = sqlHapusAnak(q(tb.Nama), tb.KolomInduk, tb.Jenis != "")
		tulis["SisipBaris "+kunci] = sqlSisipBaris(q(tb.Nama), tb)
	}
	for nama, s := range tulis {
		if !polaDML.MatchString(s) {
			t.Fatalf("%s bukan DML - DBMS_SQL.PARSE mengeksekusi DDL, ditolak", nama)
		}
		_, err := d.ExecContext(ctx, sqlUrai, s)
		h.catat(t, "urai "+nama, err)
	}

	// ---- DATA_JSON OS: bentuk JSONHalamanPega lolos CHECK `DATA_JSON IS JSON`.
	var sah int
	err := d.QueryRowContext(ctx, `SELECT COUNT(*) FROM DUAL WHERE :1 IS JSON`, models.JSONHalamanPega(models.HalamanJSON{
		Nilai:  map[string]string{"CauseOfLoss": `UJI "A"/B`, "Value": "1.5", "pxObjClass": models.KelasOSAkseptasi},
		Daftar: []models.DaftarJSON{{Nama: "CNPLayerList", Isi: []models.HalamanJSON{{Nilai: map[string]string{"XOL": "UJI"}}}}},
	})).Scan(&sah)
	if err == nil && sah != 1 {
		err = errors.New("JSONHalamanPega ditolak IS JSON")
	}
	h.catat(t, "baca DATA_JSON IS JSON", err)

	// ---- bacaan: dijalankan (SELECT), masukan UJI-*
	a := AcuanDari(g, false)
	ap := AcuanDari(g, true)
	baca := []struct {
		nama string
		f    func() error
	}{
		{"NamaMataUang", func() error { _, err := a.NamaMataUang(ctx, "UJI"); return err }},
		{"JenisReasXOL", func() error { _, err := a.JenisReasXOL(ctx, "XL 1ST LAYER"); return err }},
		{"MasterTreaty", func() error { _, _, err := a.MasterTreaty(ctx, "UJI"); return err }},
		{"DaftarMaster IN", func() error { _, err := a.DaftarMaster(ctx, models.SaringanMaster{TreatyID: "UJI"}); return err }},
		{"DaftarMaster INEDM", func() error {
			_, err := a.DaftarMaster(ctx, models.SaringanMaster{Sumber: models.SumberMasterINEDM, Ceding: "UJI"})
			return err
		}},
		{"BarisMasterDari", func() error { _, _, err := a.BarisMasterDari(ctx, "IN", "UJI", "UJI", "UJI"); return err }},
		{"AdaPolisMaster", func() error { _, err := a.AdaPolisMaster(ctx, "UJI"); return err }},
		{"DaftarPolis", func() error { _, err := a.DaftarPolis(ctx, "UJI-0000001"); return err }},
		{"BerkasPolis", func() error { _, _, err := a.BerkasPolis(ctx, "UJI"); return err }},
		{"NamaTreatySpreading", func() error { _, err := a.NamaTreatySpreading(ctx, "UJI", "2026-01-01"); return err }},
		{"Wilayah", func() error { _, err := a.Wilayah(ctx, "UJI"); return err }},
		{"AlamatAgen", func() error { _, err := a.AlamatAgen(ctx, "UJI"); return err }},
		{"NamaAdjuster", func() error { _, err := a.NamaAdjuster(ctx, "UJI"); return err }},
		{"DaftarAdjuster", func() error { _, err := a.DaftarAdjuster(ctx, "UJI"); return err }},
		{"MarketingPolis", func() error { _, _, err := a.MarketingPolis(ctx, "UJI"); return err }},
		{"RiwayatKlaimPolis", func() error { _, err := a.RiwayatKlaimPolis(ctx, "UJI"); return err }},
		{"KodeLamaBisnis", func() error { _, err := a.KodeLamaBisnis(ctx, "UJI"); return err }},
		{"RosterKomite tingkat 1", func() error { _, err := a.RosterKomite(ctx, true); return err }},
		{"RosterKomite semua", func() error { _, err := a.RosterKomite(ctx, false); return err }},
		{"TingkatPelaku", func() error { _, err := a.TingkatPelaku(ctx, "UJI"); return err }},
		{"EmailPelaku", func() error { _, err := a.EmailPelaku(ctx, "UJI"); return err }},
		{"NamaPelaku", func() error { _, err := a.NamaPelaku(ctx, "UJI"); return err }},
		{"RekeningBank", func() error { _, err := a.RekeningBank(ctx, "UJI", "UJI"); return err }},
		{"RekeningBankKlien", func() error { _, err := a.RekeningBankKlien(ctx, "UJI"); return err }},
		{"RekeningBankNama", func() error { _, err := a.RekeningBankNama(ctx, "UJI", ""); return err }},
		{"RekeningBankNama mata uang", func() error { _, err := a.RekeningBankNama(ctx, "UJI", "UJI"); return err }},
		{"DaftarKlien", func() error { _, err := a.DaftarKlien(ctx, "UJI"); return err }},
		{"IDBankRekening", func() error { _, err := a.IDBankRekening(ctx, "UJI", "UJI", "UJI"); return err }},
		{"StatusKasir", func() error { _, _, err := a.StatusKasir(ctx, "UJI"); return err }},
		{"StatusKonversi (produksi)", func() error { _, err := ap.StatusKonversi(ctx, "UJI"); return err }},
		{"EmailCeding (produksi)", func() error { _, err := ap.EmailCeding(ctx, "UJI"); return err }},
		{"RiwayatMaster (CLAIMXOL INVALID, OQ-CNP-28)", func() error {
			_, err := a.RiwayatMaster(ctx, [3]string{"UJI", "UJI/R01", "UJI/R02"})
			return err
		}},
		{"LampiranBayar", func() error { _, err := a.LampiranBayar(ctx, "UJI-CLMNP"); return err }},
		{"DaftarMataUang", func() error { _, err := a.DaftarMataUang(ctx); return err }},
		{"DaftarProvinsi", func() error { _, err := a.DaftarProvinsi(ctx, "UJI"); return err }},
		{"DaftarSebab", func() error { _, err := a.DaftarSebab(ctx, "UJI"); return err }},
		{"BarisSebabID", func() error { _, _, err := a.BarisSebabID(ctx, "UJI"); return err }},
		{"DaftarKatastrofe", func() error { _, err := a.DaftarKatastrofe(ctx, "UJI"); return err }},
		{"BarisKatastrofeID", func() error { _, _, err := a.BarisKatastrofeID(ctx, "UJI"); return err }},
		{"TanggaKomite", func() error { _, err := a.TanggaKomite(ctx, "UJI"); return err }},
		{"Keadaan", func() error {
			_, err := g.Keadaan(ctx, nil, "UJI-CLMNP")
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
			_, err := g.BacaHalaman(ctx, nil, "UJI-CLMNP")
			if errors.Is(err, ErrKasusTidakAda) {
				return nil
			}
			return err
		}},
		{"JumlahOS", func() error { _, err := g.JumlahOS(ctx, nil, "UJI-CLMNP", "XL 1ST LAYER", "IDR"); return err }},
	}
	for _, b := range baca {
		h.catat(t, "baca "+b.nama, b.f())
	}
	t.Logf("RINGKASAN: %d ok, %d menunggu migrasi / objek DBA, %d gagal", h.ok, h.belum, h.gagal)
	if strings.TrimSpace(cfg.OracleSchema) == "" {
		t.Logf("ORACLE_SCHEMA kosong")
	}
}

// TestHitungMasterDEV - mesin XoL atas master DEV 1001789 yang DIBACA repository (JSON_TABLE M_TREATY_IN, REINSURANCETYPE)
// dibandingkan dengan hasil Pega tersimpan kasus CLMNP-3998 (JSON_KLAIM; nilai dicatat di models/hitungxol_test.go).
func TestHitungMasterDEV(t *testing.T) {
	d, _ := bukaDEV(t)
	defer func() { _ = d.Close() }()
	ctx, batal := context.WithTimeout(context.Background(), 2*time.Minute)
	defer batal()
	a := AcuanDari(Baru(d), false)
	m, ada, err := a.MasterTreaty(ctx, "1001789")
	if err != nil || !ada {
		t.Fatalf("master 1001789: ada=%v err=%v", ada, err)
	}
	t.Logf("master: %d layer, %d ringkasan, %d kurs, %d share", len(m.Limits), len(m.LimitSummaryList),
		len(m.CurrencyList), len(m.Share))
	h := models.HalamanBaru()
	h.Setel(models.CD+"ShareCeding", "100")
	h.Setel(models.CD+"DeductibleType", "false")
	h.Setel(models.TM+"RNMShare", m.RNMShare)
	h.Setel(models.TM+"TreatyGroup", "PROPERTY")
	h.SetelDaftar(models.DaftarClaimAmount, []models.Baris{{"Currency": "IDR", "CurrencyID": "10026", "AltValue": "1.0",
		"Value": "10000000000", "TPL": "0", "AdjusterFee": "0", "Salvage": "0", "CNPOthersFee": "0", "CNPTSI": "10000000000"}})
	h.SetelDaftar(models.DaftarLossAlloc, []models.Baris{{"Currency": "IDR", "CurrencyID": "10026", "TreatyName": "OR",
		"ClaimPercentage": "100", "CNPFlagXOL": "true"}})
	k := &models.Konteks{Ctx: ctx, Acuan: a, Pelaku: "UJI", Sekarang: time.Now()}
	if err := models.HitungKlaim(k, h, m); err != nil {
		t.Fatal(err)
	}
	ingin := [][3]string{{"UR", "2400000000", "0"}, {"10046", "3100000000", "930000000"}, {"10036", "4500000000", "1350000000"}}
	xol := h.AmbilDaftar(models.DaftarXOL)
	if len(xol) != len(ingin) {
		t.Fatalf("SpreadingRisk %d baris, DEV 3", len(xol))
	}
	for i, w := range ingin {
		x := xol[i]
		if x["TreatyType"] != w[0] || !models.SamaAngka(x["TotalClaim"], w[1]) || !models.SamaAngka(x["ClaimSpreaded"], w[2]) {
			t.Errorf("baris %d = %s / %s / %s, DEV %v", i+1, x["TreatyType"], x["TotalClaim"], x["ClaimSpreaded"], w)
		}
	}
	qs := h.AmbilDaftar(models.DaftarBreakQS)
	if len(qs) != 2 || !models.SamaAngka(qs[0]["ClaimSpreaded"], "1368000000") || !models.SamaAngka(qs[1]["ClaimSpreaded"], "912000000") {
		t.Errorf("Break QS = %v", qs)
	}
	t.Logf("CLMNP-3998: SpreadingRisk %d baris, Break QS %d baris - cocok dengan DEV bila tanpa galat di atas", len(xol), len(qs))
	// GetDataOS atas baris OS STS 0 kasus Pega lama (kunci `ASM-FW-GCNMFW-WORK CLMNP-3998`).
	os, err := Baru(d).JumlahOS(ctx, nil, "CLMNP-3998", "XL 1ST LAYER", "IDR")
	if err != nil {
		t.Fatal(err)
	}
	if !models.SamaAngka(os.Value, "930000000") || !models.SamaAngka(os.GrossValue, "3100000000") {
		t.Errorf("JumlahOS CLMNP-3998 XL 1ST LAYER = %+v, DEV Value 930000000 / GrossValue 3100000000", os)
	}
}
