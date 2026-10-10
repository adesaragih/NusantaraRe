package services

// Kalimat layar Treaty Contract Out BERBAHASA INGGRIS [keputusan work owner
// 30-09-2026: "untuk bahasa pake bahasa inggris, jangan indo"].
//
// Kalimat modul ini sudah Inggris. Kalimat sentinel milik `inti/` (bersama;
// tetap berbahasa Indonesia untuk modul lain) yang ikut di rantai galat Treaty
// diganti padanannya di SATU tempat ini - dipakai handler (`pesanTCO`).

import (
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
)

// pesanIntiInggrisTCO - kalimat asal (`Error()` sentinel `inti/`) -> padanan Inggris.
var pesanIntiInggrisTCO = []struct{ asal, inggris string }{
	{galat.ErrPermintaanTidakSah.Error(), "services: invalid request"},
	{db.ErrMataUangTidakDikenal.Error(), "repository: unknown currency code"},
	{db.ErrTanpaOracle.Error(), "repository: ORACLE_DSN is not configured"},
}

// TeksInggrisTCO mengganti kalimat sentinel `inti/` di dalam teks galat.
func TeksInggrisTCO(s string) string {
	for _, g := range pesanIntiInggrisTCO {
		s = strings.ReplaceAll(s, g.asal, g.inggris)
	}
	return s
}
