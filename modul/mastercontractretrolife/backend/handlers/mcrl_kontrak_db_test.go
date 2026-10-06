//go:build db

package handlers_test

// Seam HTTP simpan kontrak terhadap Oracle NYATA (paket 3): uang ditulis
// kebal NLS (`TO_NUMBER(koef) / POWER(10, skala)`) dan dibaca identik.

import (
	"net/http"
	"strings"
	"testing"
)

func TestDBKontrakUangIdentikSelisihDitulis(t *testing.T) {
	u := pasangDB(t)
	s := u.skema
	u.exec(t, `INSERT INTO `+s+`.TREATYYEAR_LIFE (ID, TREATYYEAR, STARTDATE, ENDDATE)
		VALUES ('1000001', '2026', DATE '2026-01-01', DATE '2026-12-31')`)
	u.exec(t, `INSERT INTO `+s+`.REINSURANCETYPE (ID, NOTE, FLAG) VALUES ('10196', 'QS', '1')`)
	kode, badan := u.kirim(t, "POST", "/api/master-contract-retro-life/tahun/1000001/kontrak",
		`{"reinsTypeId":"10196","bIdr":"1000000000.25","idr":"1500000000.123456789","bUsd":"0.5","usd":"2.5"}`)
	if kode != http.StatusOK || !strings.Contains(badan, `"id":"1000044"`) {
		t.Fatalf("POST kontrak: %d %s", kode, badan)
	}
	idr := u.teks(t, `SELECT TO_CHAR(IDR, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''') FROM {s}.TREATYCONTRACT_LIFE WHERE ID = '1000044'`)
	sel := u.teks(t, `SELECT TO_CHAR(IDR_SELISIH, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''') FROM {s}.TREATYCONTRACT_LIFE WHERE ID = '1000044'`)
	if idr != "1500000000.123456789" || sel != "500000000.123456789" {
		t.Errorf("IDR %q IDR_SELISIH %q", idr, sel)
	}
	if n := u.cacah(t, "TREATYCONTRACT_LIFE", "ID = '1000044' AND USD IS NULL AND USD_SELISIH IS NULL "+
		"AND REINSTYPENAME = 'QS' AND TREATYSTARTDATE = DATE '2026-01-01'"); n != 1 {
		t.Errorf("USD kosong / nama / tanggal salinan: %d", n)
	}
}
