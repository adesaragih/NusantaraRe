//go:build db

package handlers_test

// Seam HTTP rute baca terhadap skema uji Oracle NYATA (paket 1).
// Tanpa ORACLE_DSN seluruhnya MELEWATI dengan pesan.

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func (u *ujiDB) get(t *testing.T, jalur string) (int, string) {
	t.Helper()
	req, err := http.NewRequest("GET", u.srv.URL+jalur, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Pelaku", "UJI-PELAKU")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(b)
}

func TestDBGridDibacaDenganSaringanDanUrutRD(t *testing.T) {
	u := pasangDB(t)
	s := u.skema
	u.exec(t, `INSERT INTO `+s+`.TREATYYEAR_LIFE (ID, TREATYYEAR, UNDERWRITINGYEAR, STARTDATE, ENDDATE)
		VALUES ('1000001', '2026', '2026', DATE '2026-01-01', DATE '2026-12-31')`)
	u.exec(t, `INSERT INTO `+s+`.TREATYCONTRACT_LIFE (ID, IDTREATYYEAR, REINSTYPEID, REINSTYPENAME, IDR, B_IDR, B_USD)
		VALUES ('1000002', '1000001', '10196', 'QS', 1500000000.123456789, 0, 0)`)
	u.exec(t, `INSERT INTO `+s+`.TREATYCONTRACT_LIFE (ID, IDTREATYYEAR, REINSTYPEID) VALUES ('1000001', '1000001', '10197')`)
	u.exec(t, `INSERT INTO `+s+`.TREATYREINSURER_LIFE (ID, TREATYYEARID, TREATYCONTRACTID, PCTSHARE)
		VALUES ('1000001', '1000001', '1000002', 40)`)
	u.exec(t, `INSERT INTO `+s+`.TREATYREINSURER_LIFE (ID, TREATYYEARID, TREATYCONTRACTID, PCTSHARE)
		VALUES ('1000009', '1000001', '1000002', 10.5)`)
	u.exec(t, `INSERT INTO `+s+`.TREATYREINSURER_LIFE (ID, TREATYYEARID, TREATYCONTRACTID, PCTSHARE)
		VALUES ('1000010', 'LAIN', '1000002', 1)`)

	kode, badan := u.get(t, "/api/master-contract-retro-life/tahun/1000001/kontrak")
	if kode != http.StatusOK {
		t.Fatalf("kontrak: %d %s", kode, badan)
	}
	if strings.Index(badan, `"id":"1000001"`) > strings.Index(badan, `"id":"1000002"`) {
		t.Errorf("kontrak tidak urut ID ASC: %s", badan)
	}
	if !strings.Contains(badan, `"idr":"1500000000.123456789"`) {
		t.Errorf("uang tidak identik menyeberang batas: %s", badan)
	}
	kode, badan = u.get(t, "/api/master-contract-retro-life/kontrak/1000002/reinsurer")
	if kode != http.StatusOK {
		t.Fatalf("reinsurer: %d %s", kode, badan)
	}
	if strings.Contains(badan, "1000010") {
		t.Errorf("baris bersalinan tahun lain ikut tampil: %s", badan)
	}
	if strings.Index(badan, `"id":"1000009"`) > strings.Index(badan, `"id":"1000001"`) {
		t.Errorf("reinsurer tidak urut ID DESC: %s", badan)
	}
	if !strings.Contains(badan, `"pctShare":"10.5"`) {
		t.Errorf("share tidak teks persis: %s", badan)
	}
}

func TestDBMasterJenisHanyaLife(t *testing.T) {
	u := pasangDB(t)
	s := u.skema
	u.exec(t, `INSERT INTO `+s+`.REINSURANCETYPE (ID, NOTE, FLAG) VALUES ('10196', 'QS', '1')`)
	u.exec(t, `INSERT INTO `+s+`.REINSURANCETYPE (ID, NOTE, FLAG) VALUES ('10200', 'OR', '1')`)
	u.exec(t, `INSERT INTO `+s+`.REINSURANCETYPE (ID, NOTE, FLAG) VALUES ('10004', 'UJI-NONLIFE', 'active')`)
	kode, badan := u.get(t, "/api/master-contract-retro-life/jenis-reasuransi")
	if kode != http.StatusOK || strings.Contains(badan, "UJI-NONLIFE") || !strings.Contains(badan, `"total":2`) {
		t.Errorf("jenis life: %d %s", kode, badan)
	}
}
