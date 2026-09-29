package repository

// Teks SQL pelaksana penyimpanan nyata (OQ-TCO-08) - tanpa Oracle.

import (
	"context"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/db"
)

func TestSQLStorageTCO(t *testing.T) {
	app := sqlAppStorageTCO("SKEMA_UJI." + MasterFolderImageTCO)
	tok := sqlTokenStorageBerlakuTCO("SKEMA_UJI.GCP_IMAGE")
	for nama, q := range map[string]string{"app": app, "token": tok} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Errorf("%s: %v", nama, err)
		}
		if !strings.HasPrefix(strings.TrimSpace(q), "SELECT") {
			t.Errorf("%s bukan bacaan: %s", nama, q)
		}
	}
	// `GET_TOKEN_STORAGE`: token TERBARU yang belum kedaluwarsa, nilai lewat bind;
	// sisa umur dihitung DI ORACLE (tanpa membaca DATE ke jam aplikasi).
	for _, wajib := range []string{"(CAST(INPUTDATE AS DATE) - CAST(:1 AS DATE)) * 86400", "APPNAME = :2",
		"INPUTDATE > :3", "ORDER BY INPUTDATE DESC", "FETCH FIRST 1 ROWS ONLY"} {
		if !strings.Contains(tok, wajib) {
			t.Errorf("kueri token tanpa %q", wajib)
		}
	}
}

func TestTokenStorageBerlakuTCOMenuntutTransaksi(t *testing.T) {
	var d *db.DB
	if _, _, err := TokenStorageBerlakuTCO(context.Background(), d, nil, "UJI-APP", time.Now(), time.Second); err == nil {
		t.Error("tanpa transaksi diterima")
	}
}
