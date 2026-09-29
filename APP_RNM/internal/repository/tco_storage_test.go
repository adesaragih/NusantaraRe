package repository

// Teks SQL pelaksana penyimpanan nyata (OQ-TCO-08) - tanpa Oracle.

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestSQLStorageTCO(t *testing.T) {
	app := sqlAppStorageTCO("SKEMA_UJI." + MasterFolderImageTCO)
	tok := sqlTokenStorageBerlakuTCO("SKEMA_UJI.GCP_IMAGE")
	for nama, q := range map[string]string{"app": app, "token": tok} {
		if err := PeriksaSQL(q); err != nil {
			t.Errorf("%s: %v", nama, err)
		}
		if !strings.HasPrefix(strings.TrimSpace(q), "SELECT") {
			t.Errorf("%s bukan bacaan: %s", nama, q)
		}
	}
	// `GET_TOKEN_STORAGE`: token TERBARU yang belum kedaluwarsa, nilai lewat bind.
	for _, wajib := range []string{"APPNAME = :1", "INPUTDATE > :2", "ORDER BY INPUTDATE DESC", "FETCH FIRST 1 ROWS ONLY"} {
		if !strings.Contains(tok, wajib) {
			t.Errorf("kueri token tanpa %q", wajib)
		}
	}
}

func TestTokenStorageBerlakuTCOMenuntutTransaksi(t *testing.T) {
	var d *DB
	if _, _, err := d.TokenStorageBerlakuTCO(context.Background(), nil, "UJI-APP", time.Now()); err == nil {
		t.Error("tanpa transaksi diterima")
	}
}
