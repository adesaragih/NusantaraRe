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

func TestDBMasterPlanCariCoverNameAtauBusiness(t *testing.T) {
	u := pasangDB(t)
	u.exec(t, `INSERT INTO {s}.PRODUCT_TYPE_LIFE VALUES ('P2', 'UJI COVER DUA', 'UJI KREDIT', 'UJI MANFAAT')`)
	u.exec(t, `INSERT INTO {s}.PRODUCT_TYPE_LIFE VALUES ('P1', 'UJI COVER SATU', 'UJI LAIN', 'UJI MANFAAT')`)
	kode, badan := u.kirim(t, "GET", pre+"/master-plan?cari=kredit", "")
	if kode != http.StatusOK || !strings.Contains(badan, `"total":1`) || !strings.Contains(badan, `"id":"P2"`) {
		t.Errorf("cari Business: %d %s", kode, badan)
	}
	kode, badan = u.kirim(t, "GET", pre+"/master-plan", "")
	if kode != http.StatusOK || strings.Index(badan, `"id":"P1"`) > strings.Index(badan, `"id":"P2"`) {
		t.Errorf("urut ID: %d %s", kode, badan)
	}
}

// K1 keputusan work owner 01-10-2026 (OQ-MPNL-03): `Choose R/I Rate` membaca view `RATE_LIFE_SUMMARY`
// (`Contains` USEDBY, urut `ID ASC`), `View Rate` membaca view `RATE_LIFE` satu RIRATEID (urut `ID DESC,
// RATE ASC`, RATE teks apa adanya); nol tulisan.
func TestDBRIRateDanViewRateDibacaSaja(t *testing.T) {
	u := pasangDB(t)
	for _, q := range []string{
		`INSERT INTO {s}.RATE_LIFE_SUMMARY (ID, USEDBY) VALUES ('UJI-2', 'UJI RATE DUA')`,
		`INSERT INTO {s}.RATE_LIFE_SUMMARY (ID, USEDBY) VALUES ('UJI-1', 'uji rate satu')`,
		`INSERT INTO {s}.RATE_LIFE_SUMMARY (ID, USEDBY) VALUES ('UJI-3', 'LAIN')`,
		`INSERT INTO {s}.RATE_LIFE (ID, IDUSEDBY, USEDBY, GENDER, CONTRACT, AGE, RATE) VALUES ('UJI-A', 'UJI-1', 'UJI', 'U', '10', '30', '0,5')`,
		`INSERT INTO {s}.RATE_LIFE (ID, IDUSEDBY, USEDBY, GENDER, CONTRACT, AGE, RATE) VALUES ('UJI-B', 'UJI-1', 'UJI', 'U', '10', '31', '1.25')`,
		`INSERT INTO {s}.RATE_LIFE (ID, IDUSEDBY, USEDBY, GENDER, CONTRACT, AGE, RATE) VALUES ('UJI-C', 'UJI-2', 'UJI', 'U', '10', '30', '9')`,
	} {
		u.exec(t, q)
	}
	// Sidik isi kedua view sebelum dan sesudah (code review #10: cacah saja tidak melihat UPDATE).
	sidik := func() string {
		return u.teks(t, `SELECT (SELECT COUNT(*) || '/' || SUM(LENGTH(ID || IDUSEDBY || USEDBY || GENDER || CONTRACT || AGE || RATE))
		    FROM {s}.RATE_LIFE) || '|' || (SELECT COUNT(*) || '/' || SUM(LENGTH(ID || USEDBY)) FROM {s}.RATE_LIFE_SUMMARY) FROM DUAL`)
	}
	awal := sidik()
	kode, badan := u.kirim(t, "GET", pre+"/master/ri-rate?cari=rate", "")
	if kode != http.StatusOK || !strings.Contains(badan, `"total":2`) || strings.Index(badan, "UJI-1") > strings.Index(badan, "UJI-2") {
		t.Errorf("R/I Rate: %d %s", kode, badan)
	}
	kode, badan = u.kirim(t, "GET", pre+"/rate?riRateId=UJI-1", "")
	if kode != http.StatusOK || !strings.Contains(badan, `"rate":"0,5"`) || strings.Contains(badan, "UJI-C") ||
		strings.Index(badan, "UJI-B") > strings.Index(badan, "UJI-A") || !strings.Contains(badan, `"terpotong":false`) {
		t.Errorf("View Rate: %d %s", kode, badan)
	}
	if akhir := sidik(); akhir != awal {
		t.Errorf("view rate tersentuh: sidik %s → %s", awal, akhir)
	}
}
