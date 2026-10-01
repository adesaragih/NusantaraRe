//go:build db

package handlers_test

// Seam HTTP ketujuh pemilih master terhadap skema uji Oracle NYATA (paket 2).

import (
	"net/http"
	"strings"
	"testing"
)

func TestDBPemilihMasterSaringanRD(t *testing.T) {
	u := pasangDB(t)
	for _, q := range []string{
		`INSERT INTO {s}.AGENT VALUES ('L0UJI2', 'UJI ZETA', '1')`,
		`INSERT INTO {s}.AGENT VALUES ('L0UJI1', 'UJI ALFA', '1')`,
		`INSERT INTO {s}.AGENT VALUES ('L0UJI3', 'UJI NONAKTIF', '0')`,
		`INSERT INTO {s}.AGENT VALUES ('X9UJI4', 'UJI BUKAN LIFE', '1')`,
		`INSERT INTO {s}.CLIENT VALUES ('UJI-ORG-1', 'UJI PEMEGANG', 'LIFE')`,
		`INSERT INTO {s}.CLIENT VALUES ('UJI-ORG-2', '-', 'LIFE')`,
		`INSERT INTO {s}.CURRENCY VALUES ('1', 'IDR')`,
		`INSERT INTO {s}.CURRENCY VALUES ('9', 'ITL')`,
		`INSERT INTO {s}.RIRISK_LIFE_SUMMARY VALUES ('1000117', 'UJI RISK')`,
		`INSERT INTO {s}.CAUSEOFLOSS_LIFE VALUES ('100004', 'ANY CAUSE')`,
	} {
		u.exec(t, q)
	}
	kode, badan := u.kirim(t, "GET", pre+"/master/ceding?cari=uji", "")
	if kode != http.StatusOK || !strings.Contains(badan, `"total":2`) ||
		strings.Index(badan, "UJI ALFA") > strings.Index(badan, "UJI ZETA") {
		t.Errorf("AGENT: L0 + aktif, urut ClientName: %d %s", kode, badan)
	}
	if kode, badan := u.kirim(t, "GET", pre+"/master/pemegang-polis", ""); kode != http.StatusOK ||
		!strings.Contains(badan, `"total":1`) {
		t.Errorf("CLIENT Name != '-': %d %s", kode, badan)
	}
	if kode, badan := u.kirim(t, "GET", pre+"/master/mata-uang", ""); kode != http.StatusOK || strings.Contains(badan, "ITL") {
		t.Errorf("CURRENCY != ITL: %d %s", kode, badan)
	}
	if kode, badan := u.kirim(t, "GET", pre+"/master/ri-risk?cari=risk", ""); kode != http.StatusOK ||
		!strings.Contains(badan, `"id":"1000117"`) {
		t.Errorf("RIRISK: %d %s", kode, badan)
	}
	if kode, badan := u.kirim(t, "GET", pre+"/master/penyebab?cari=a%25", ""); kode != http.StatusOK ||
		!strings.Contains(badan, `"total":0`) {
		t.Errorf("wildcard diloloskan: %d %s", kode, badan)
	}
}
