package models

// Uji murni pengurai nilai Pega SATU tempat (temuan tinjauan P9: semula
// tersalin di repository/masterxol.go dan models/dokumenlama.go).

import (
	"encoding/json"
	"testing"
)

func TestTeksSkalarJSONTanpaFloat(t *testing.T) { // ADR-0003, P29 sifat 1
	for masuk, harap := range map[any]string{
		nil:                                     "",
		"UJI-teks":                              "UJI-teks",
		true:                                    "true",
		json.Number("592629512.880000276"):      "592629512.880000276",
		json.Number("1000000000000000000000.1"): "1000000000000000000000.1",
		json.Number("1E+2"):                     "100",
		json.Number("12.50"):                    "12.50",
	} {
		got, ok := TeksSkalarJSON(masuk)
		if !ok || got != harap {
			t.Errorf("TeksSkalarJSON(%v) = %q %v, harap %q", masuk, got, ok, harap)
		}
	}
	for _, bukan := range []any{map[string]any{}, []any{}} {
		if _, ok := TeksSkalarJSON(bukan); ok {
			t.Errorf("%T bukan skalar", bukan)
		}
	}
}

func TestTanggalMasterPegaSamaDenganDokumenLama(t *testing.T) { // diagram F21, K15
	for masuk, harap := range map[string]string{
		"20261101":                "2026-11-01",
		"20270131T170000.000 GMT": "2027-02-01", // 17.00 GMT = 00.00 WIB hari berikut
		"20270131T165959 GMT":     "2027-01-31",
		"2026-11-01 00:00:00":     "2026-11-01",
		"05/06/2017":              "05/06/2017", // ambigu: tidak ditebak, penyimpanan menolak
		"UJI-bukan-tanggal":       "UJI-bukan-tanggal",
		"":                        "",
	} {
		if got := TanggalMasterPega(masuk); got != harap {
			t.Errorf("TanggalMasterPega(%q) = %q, harap %q", masuk, got, harap)
		}
	}
}
