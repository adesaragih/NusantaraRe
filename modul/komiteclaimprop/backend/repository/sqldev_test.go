//go:build ujidev

// Uji manual SQL Komite Claim Prop di DEV - BACA SAJA (prompt §9): SELECT dijalankan lewat metode aslinya dengan
// masukan UJI-*; setiap INSERT / UPDATE, dan SELECT yang hanya berjalan di dalam transaksi tulis, DIURAI lewat
// DBMS_SQL.PARSE (tidak dieksekusi). Hasil dicetak sebagai cacah dan kode galat saja, nol isi baris.
//
//	go test -tags ujidev -run TestSQLKomiteDiDEV -v ./modul/komiteclaimprop/backend/repository/
//
// Nomor akseptasi (penomor AwalanProduksi / HariClosing / UrutNomorBerikut) memakai fungsi `inti/backend/penomor`
// yang sama dengan nomor PLA / DLA Claim Prop - diuji di sana; urutannya MENULIS GENERATE_SEQUENCE_NUMBER, jadi tidak
// dijalankan di DEV.
//
// Galat "objek belum ada" (ORA-00942 / 00904) pada KOMITE_USUL_TUTUP / KOMITE_USUL_CADANG berarti migrasi 680 belum
// dijalankan work owner - dicatat, tidak menggagalkan.
package repository

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/komiteclaimprop/backend/models"
)

var (
	polaDML      = regexp.MustCompile(`(?is)^\s*(INSERT|UPDATE|DELETE|SELECT)\b`)
	polaBelumAda = regexp.MustCompile(`ORA-(00942|00904|02289)`)
	polaORA      = regexp.MustCompile(`ORA-\d{5}`)
)

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
		t.Logf("BELUM  %s: %s (objek belum ada - migrasi 680)", nama, polaBelumAda.FindString(err.Error()))
	default:
		h.gagal++
		kode := polaORA.FindString(err.Error())
		if kode == "" {
			kode = err.Error()
		}
		t.Errorf("GAGAL  %s: %s", nama, kode)
	}
}

func TestSQLKomiteDiDEV(t *testing.T) {
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

	// ---- tulisan dan SELECT bertransaksi: diurai saja
	dataJSON := models.JSONHalamanPega(map[string]string{"KomiteNo": "TKMT-UJI", "CauseOfLoss": `UJI "A"/B`,
		"Value": "1.5", "pxObjClass": models.KelasOSAkseptasi})
	sqlOS, _, err := sqlSisipOS(q("OS_AKSEPTASI_KLAIM"), models.BarisOSAkseptasi{CaseID: "UJI-C", Value: "1",
		GrossValue: "1", KursValue: "1", PersenRNM: "1", EstimationDate: time.Now(), DataJSON: dataJSON}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	tulis := map[string]string{
		"KepalaKasus FOR UPDATE (SELECT)":       sqlKepalaKasus(q("T_GENERAL_KOMITE"), q("T_WORK_CLAIM"), true),
		"SimpanKepala T_GENERAL_KOMITE":         sqlSimpanKepala(q("T_GENERAL_KOMITE")),
		"TutupKasus T_WORK_CLAIM":               sqlTutupKasus(q("T_WORK_CLAIM")),
		"SentuhKasus T_WORK_CLAIM":              sqlSentuhKasus(q("T_WORK_CLAIM")),
		"TulisAnggota +komentar":                sqlTulisAnggota(q(tabelTangga), true),
		"TulisAnggota tanpa komentar":           sqlTulisAnggota(q(tabelTangga), false),
		"SisipOS OS_AKSEPTASI_KLAIM":            sqlOS,
		"AdaJSONKlaim (SELECT)":                 sqlAdaJSONKlaim(q("JSON_KLAIM")),
		"SisipJSONKlaim JSON_KLAIM":             sqlSisipJSONKlaim(q("JSON_KLAIM")),
		"LogLayanan MONITORING_KLAIM_LOG":       sqlLogLayanan(q("MONITORING_KLAIM_LOG")),
		"RiwayatAkseptasi HISTORYAKSEPTASIPEGA": sqlRiwayatAkseptasi(q("HISTORYAKSEPTASIPEGA")),
		"SisipDokumenKlaim (tabel warisan)":     sqlSisipDokumenKlaim(q(TabelDokumenKlaim)),
	}
	for nama, s := range tulis {
		if !polaDML.MatchString(s) {
			t.Fatalf("%s bukan DML - DBMS_SQL.PARSE mengeksekusi DDL, ditolak", nama)
		}
		_, err := d.ExecContext(ctx, sqlUrai, s)
		h.catat(t, "urai "+nama, err)
	}

	// ---- DATA_JSON OS (keputusan work owner 08-10-2026): bentuk JSONHalamanPega lolos CHECK `DATA_JSON IS JSON`.
	var sah int
	err = d.QueryRowContext(ctx, `SELECT COUNT(*) FROM DUAL WHERE :1 IS JSON`, dataJSON).Scan(&sah)
	if err == nil && sah != 1 {
		err = errors.New("JSONHalamanPega ditolak IS JSON")
	}
	h.catat(t, "baca DATA_JSON IS JSON", err)

	// ---- bacaan: dijalankan (SELECT), masukan UJI-*
	a := AcuanDari(g, false)
	baca := []struct {
		nama string
		f    func() error
	}{
		{"BacaKasus", func() error {
			_, err := g.BacaKasus(ctx, nil, "TKMT-UJI000", false)
			if errors.Is(err, ErrKasusTidakAda) {
				return nil
			}
			return err
		}},
		{"BacaTangga", func() error { _, err := g.BacaTangga(ctx, nil, "TKMT-UJI000"); return err }},
		{"DaftarKerja", func() error {
			rows, err := g.DaftarKerja(ctx, "UJI-K1", []string{"UJI-WB-1", "UJI-WB-2"})
			if err == nil {
				t.Logf("       daftar kerja UJI-K1: %d baris", len(rows))
			}
			return err
		}},
		{"CacahTanpaGiliran", func() error {
			n, err := cacahTanpaGiliran(g, ctx)
			if err == nil {
				t.Logf("       kasus komite PROP terbuka tanpa baris menunggu: %d", n)
			}
			return err
		}},
		{"TahunTreaty", func() error { _, err := a.TahunTreaty(ctx, "UJI", "20260101"); return err }},
		{"LimitPLA", func() error { _, err := a.LimitPLA(ctx, "2026", "UJI", "UJI"); return err }},
		{"DaftarRetro", func() error { _, err := a.DaftarRetro(ctx, "2026", "UJI", "UJI"); return err }},
		{"IDBankRekening", func() error { _, err := a.IDBankRekening(ctx, "UJI", "UJI", "UJI"); return err }},
		{"NamaPelaku", func() error { _, err := a.NamaPelaku(ctx, "UJI"); return err }},
		{"EmailPelaku", func() error { _, err := a.EmailPelaku(ctx, "UJI"); return err }},
		{"EmailAnggotaWorkbasket", func() error { _, err := a.EmailAnggotaWorkbasket(ctx, "UJI-WB"); return err }},
		{"KomentarAwal", func() error { _, err := g.KomentarAwal(ctx, "CLMP-UJI", "UJI", "TKMT-UJI"); return err }},
		{"StatusKonversi (produksi; DEV tanpa hak = BELUM)", func() error {
			_, err := AcuanDari(g, true).StatusKonversi(ctx, "UJI")
			return err
		}},
	}
	for _, b := range baca {
		h.catat(t, "baca "+b.nama, b.f())
	}
	t.Logf("ringkasan: ok %d, belum-migrasi %d, gagal %d", h.ok, h.belum, h.gagal)
}

// sqlCacahTanpaGiliran - kasus komite PROP terbuka tanpa satu pun baris tangga menunggu: tidak tampil di daftar kerja
// siapa pun (KomiteRouter tanpa cadangan, prompt §6 butir 6 - dilaporkan, bukan ditampilkan).
func sqlCacahTanpaGiliran(gen, work, list string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s g JOIN %s w ON w.ID = g.ID
		 WHERE w.LINI = :1 AND w.ID LIKE :2 AND w.STATUS_WORK IS NULL
		   AND NOT EXISTS (SELECT 1 FROM %s l WHERE l.DATA_KOMITE_ID = g.ID AND l.KOMITE_APPROVAL = :3)`, gen, work, list)
}

// cacahTanpaGiliran (uji DEV) mencacah kasus komite PROP terbuka yang tidak menunggu siapa pun.
func cacahTanpaGiliran(g *Gudang, ctx context.Context) (int, error) {
	gen, work, list, _, _, err := g.tabelKerja()
	if err != nil {
		return 0, err
	}
	q := sqlCacahTanpaGiliran(gen, work, list)
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	var n int
	if err := g.db.QueryRowContext(ctx, q, models.LiniProp, awalanLike(), models.KeputusanMenunggu).Scan(&n); err != nil {
		return 0, fmt.Errorf("repository: mencacah kasus komite tanpa giliran: %w", err)
	}
	return n, nil
}
