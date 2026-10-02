//go:build db

package repository_test

// Pembuatan kasus dan penyalinan versi terhadap skema uji Oracle - tiket 01-02.
// Tanpa ORACLE_DSN MELEWATI.
//
// ⛔ Tiruan tabel peserta warisan dibuat test ini sendiri dari bentuk 052
// (nama dan tipe kolom yang disalin `sqlSalinPesertaWarisan`), lalu dibuang:
// tiruan bersama skema uji hanya memuat kolom yang dibaca Claim Life.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
)

// ddlTiruanWarisan - kolom yang disalin dari tabel peserta warisan, bertipe
// seperti katalog (`STNC`/`WPC` DATE, angka NUMBER tanpa presisi).
func ddlTiruanWarisan(t *testing.T, skema string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "..", "premiumlistlife", "backend", "migrations", "052_t_premium_list_detail.sql"))
	if err != nil {
		t.Fatal(err)
	}
	// Kolom EDM `SaveMasterLPDet` (K5): nomor endorsement dan ketiga penanda - nama warisan, bukan nama 052.
	kolom := []string{"ID VARCHAR2(50)", "PL_NUMBER VARCHAR2(255)", "IDPEGA VARCHAR2(255)", "PL_NUMBER_EDM VARCHAR2(255)",
		"EDMSTATUS VARCHAR2(255)", "STATUSOLD VARCHAR2(255)", "STATUS VARCHAR2(255)"}
	lewati := map[string]bool{"ID": true, "PREMIUM_LIST_ID": true, "PARENT_ID": true, "ID_PEGA": true, "PL_NUMBER": true,
		"PL_NUMBER_EDM": true, "EDM_STATUS": true, "STATUS_OLD": true, "STATUS": true}
	for _, baris := range strings.Split(string(b), "\n") {
		m := regexp.MustCompile(`^\s*([A-Z][A-Z0-9_]*)\s+(VARCHAR2|NUMBER|DATE)`).FindStringSubmatch(baris)
		if m == nil || lewati[m[1]] {
			continue
		}
		nama, tipe := m[1], m[2]
		switch nama {
		case "PRO_RATE_TYPE":
			nama = "PRORATETYPE"
		case "STNC", "WPC":
			tipe = "DATE"
		case "RISK":
			tipe = "NUMBER"
		}
		if tipe == "VARCHAR2" {
			tipe = "VARCHAR2(255)"
		}
		kolom = append(kolom, nama+" "+tipe)
	}
	return fmt.Sprintf(`CREATE TABLE %s.M_LIFE_PREMIUM_DETAIL (%s)`, skema, strings.Join(kolom, ", "))
}

func TestBuatKasusTerhadapOracle(t *testing.T) {
	repo, skema, tutup := pasangSkema(t)
	defer tutup()
	ctx := context.Background()
	_, _ = repo.ExecContext(ctx, `DROP TABLE `+skema+`.M_LIFE_PREMIUM_DETAIL PURGE`)
	jalankan(t, repo, ddlTiruanWarisan(t, skema))
	defer func() { _, _ = repo.ExecContext(ctx, `DROP TABLE `+skema+`.M_LIFE_PREMIUM_DETAIL PURGE`) }()
	jalankan(t, repo,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST (ID, ID_PEGA, TGL_INPUT, TYPE, PRODUCT_NAME, DATE_RECEIVED)
		   VALUES ('UJI-NB-3', 'UJI-NB-3', SYSDATE, 'QR', 'UJI-PRODUK', DATE '2026-01-02')`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PL_NUMBER, CERTIFICATE_NO, SUM_INSURED, STNC)
		   VALUES ('UJI-D3A', 'UJI-NB-3', 'UJI-PL-3', 'UJI-C1', -0.00000001, '02/01/2026')`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PL_NUMBER, CERTIFICATE_NO, EDM_STATUS)
		   VALUES ('UJI-D3B', 'UJI-NB-3', 'UJI-PL-3', 'UJI-C2', 'Delete')`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_SPREADING (ID, DETAIL_ID, TREATY_TYPE_NAME, RETROCADED_SHARE) VALUES ('UJI-S3', 'UJI-D3A', 'UJI-QS', 50)`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_SPREADING_RETRO (ID, SPREADING_ID, REINSURER_NAME, PERCENT_SHARE) VALUES ('UJI-R3', 'UJI-S3', 'UJI-REAS', 100)`,
		`INSERT INTO `+skema+`.JSON_POLIS (IDPEGA, NOPOLIS, PRODKE, TGL_INPUT, DATA_JSON)
		   VALUES ('UJI-IDPEGA-W', 'UJI-PL-W', NULL, SYSDATE, '{"Type":"TR","DateReceived":"20240115","PremiumListSummary":{"RISLIPRNM":"UJI-SLIP"}}')`,
		`INSERT INTO `+skema+`.M_LIFE_PREMIUM_DETAIL (ID, PL_NUMBER, IDPEGA, CERTIFICATE_NO, PRORATETYPE, STNC, SUM_INSURED)
		   VALUES ('UJI-M1', 'UJI-PL-W', 'UJI-IDPEGA-W', 'UJI-WC1', 'UJI-PRO', DATE '2025-03-04', 12.5)`,
		`INSERT INTO `+skema+`.M_LIFE_PREMIUM_DETAIL (ID, PL_NUMBER, IDPEGA, CERTIFICATE_NO) VALUES ('UJI-M2', 'UJI-PL-W', 'UJI-IDPEGA-LAIN', 'UJI-WC2')`,
	)
	g := repository.Baru(repo)

	buat := func(polis string) (string, repository.Salinan, error) {
		t.Helper()
		tx, err := repo.Mulai(ctx)
		if err != nil {
			t.Fatal(err)
		}
		v, ada, err := g.VersiBerjalan(ctx, tx, polis, 0)
		if err != nil || !ada {
			t.Fatalf("versi %s: %v %v", polis, ada, err)
		}
		kepala, _, err := g.KepalaSumber(ctx, tx, v, polis)
		if err != nil {
			t.Fatal(err)
		}
		id, err := g.PengenalKasusBaru(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if err := g.SisipKasus(ctx, tx, repository.KasusTulis{ID: id, NomorPolis: polis, EdmType: "1", EdmDate: "2026-10-01",
			ProdKe: v.ProdKe + 1, Pembuat: "UJI-AKUN-1", Kepala: kepala}); err != nil {
			_ = tx.Rollback()
			return "", repository.Salinan{}, err
		}
		s, err := g.SalinVersi(ctx, tx, id, v, polis)
		if err != nil {
			_ = tx.Rollback()
			return "", repository.Salinan{}, err
		}
		if err := tx.Commit(); err != nil {
			t.Fatal(err)
		}
		return id, s, nil
	}

	id, s, err := buat("UJI-PL-3")
	if err != nil {
		t.Fatal(err)
	}
	if s != (repository.Salinan{Peserta: 1, Spreading: 1, SpreadingRetro: 1}) {
		t.Fatalf("salinan aplikasi = %+v (baris Delete tidak tersalin)", s)
	}
	k, err := g.AmbilKasus(ctx, nil, id, false)
	if err != nil || k.ProdKe != 2 || k.Kepala["PRODUCT_NAME"] != "UJI-PRODUK" || k.Kepala["DATE_RECEIVED"] != "2026-01-02" || k.EdmDate != "2026-10-01" {
		t.Fatalf("kasus = %+v, %v", k, err)
	}
	p, err := g.DaftarPeserta(ctx, id, 1, 20)
	if err != nil || len(p) != 1 || p[0].ParentID != "UJI-D3A" || p[0].EdmStatus != models.StatusOld || len(p[0].ID) != 32 {
		t.Fatalf("peserta salinan = %+v, %v", p, err)
	}
	r, err := g.RincianPeserta(ctx, id, p[0].ID)
	if err != nil || r.Peserta.Nilai["SUM_INSURED"] != "-0.00000001" || r.Peserta.Nilai["STNC"] != "02/01/2026" ||
		len(r.Spreading) != 1 || len(r.Spreading[0].Retro) != 1 || r.Spreading[0].Retro[0].PercentShare != "100" {
		t.Fatalf("rincian = %+v, %v (uang negatif kecil utuh; spreading ikut)", r, err)
	}

	// Index unik 482: kasus terbuka kedua atas polis yang sama ditolak basis data.
	if _, _, err := buat("UJI-PL-3"); !errors.Is(err, repository.ErrKasusTerbukaGanda) {
		t.Fatalf("kasus terbuka kedua: %v", err)
	}
	// Kasus yang ditolak tidak memblokir (OQ-EDM-002).
	jalankan(t, repo, `UPDATE `+skema+`.T_PREMIUM_LIST SET STATUSS = 'Resolved-Rejected' WHERE ID = '`+id+`'`)
	if _, _, err := buat("UJI-PL-3"); err != nil {
		t.Fatalf("sesudah Decline: %v", err)
	}

	// Versi NB warisan: peserta IDPEGA itu saja, kepala dari DATA_JSON.
	idW, sW, err := buat("UJI-PL-W")
	if err != nil || sW.Peserta != 1 {
		t.Fatalf("salinan warisan = %+v, %v", sW, err)
	}
	kW, _ := g.AmbilKasus(ctx, nil, idW, false)
	if kW.Kepala["TYPE"] != "TR" || kW.Kepala["DATE_RECEIVED"] != "2024-01-15" || kW.Kepala["RI_SLIP_RNM"] != "UJI-SLIP" || kW.ProdKe != 2 {
		t.Fatalf("kepala warisan = %+v", kW)
	}
	pW, _ := g.DaftarPeserta(ctx, idW, 1, 20)
	rW, err := g.RincianPeserta(ctx, idW, pW[0].ID)
	if err != nil || rW.Peserta.Nilai["STNC"] != "04/03/2025" || rW.Peserta.Nilai["SUM_INSURED"] != "12.5" || rW.Peserta.ParentID != "" {
		t.Fatalf("peserta warisan = %+v, %v", rW, err)
	}
	pv, total, err := g.PesertaVersi(ctx, models.Versi{Jenis: models.SumberWarisan, ID: "UJI-IDPEGA-W"}, "UJI-PL-W", 1, 20)
	if err != nil || total != 1 || len(pv) != 1 || pv[0].Nilai["CERTIFICATE_NO"] != "UJI-WC1" {
		t.Fatalf("polis lama warisan = %+v %d %v", pv, total, err)
	}
}
