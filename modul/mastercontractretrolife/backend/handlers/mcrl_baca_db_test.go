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

// K1 (01-10-2026, OQ-MCRL-13): kedua view rate dibaca dari Oracle - autocomplete `Contains` USEDBY urut
// `ID ASC`; Rate List satu IDUSEDBY urut `ID DESC, RATE ASC`, RATE teks apa adanya; nol tulisan.
func TestDBRateDibacaSaja(t *testing.T) {
	u := pasangDB(t)
	s := u.skema
	for _, r := range [][2]string{{"UJI-2", "UJI RATE DUA"}, {"UJI-1", "uji rate satu"}, {"UJI-3", "LAIN"}} {
		u.exec(t, `INSERT INTO `+s+`.M_RATE_LIFE_SUMMARY (ID, USEDBY) VALUES (:1, :2)`, r[0], r[1])
	}
	for _, r := range [][3]string{{"UJI-A", "UJI-1", "0,5"}, {"UJI-B", "UJI-1", "1.25"}, {"UJI-C", "UJI-2", "9"}} {
		u.exec(t, `INSERT INTO `+s+`.RATE_LIFE (ID, IDUSEDBY, USEDBY, GENDER, CONTRACT, AGE, RATE)
			VALUES (:1, :2, 'UJI', 'U', '10', '30', :3)`, r[0], r[1], r[2])
	}
	// Sidik isi kedua view sebelum dan sesudah (code review #10: cacah saja tidak melihat UPDATE).
	sidik := func() string {
		return u.teks(t, `SELECT (SELECT COUNT(*) || '/' || SUM(LENGTH(ID || IDUSEDBY || USEDBY || GENDER || CONTRACT || AGE || RATE))
		    FROM {s}.RATE_LIFE) || '|' || (SELECT COUNT(*) || '/' || SUM(LENGTH(ID || USEDBY)) FROM {s}.M_RATE_LIFE_SUMMARY) FROM DUAL`)
	}
	awal := sidik()
	kode, badan := u.get(t, "/api/master-contract-retro-life/ringkasan-rate?cari=RATE")
	if kode != http.StatusOK || !strings.Contains(badan, `"total":2`) || strings.Index(badan, "UJI-1") > strings.Index(badan, "UJI-2") {
		t.Errorf("ringkasan rate: %d %s", kode, badan)
	}
	kode, badan = u.get(t, "/api/master-contract-retro-life/rate?idusedby=UJI-1")
	if kode != http.StatusOK || !strings.Contains(badan, `"rate":"0,5"`) || strings.Contains(badan, "UJI-C") ||
		strings.Index(badan, "UJI-B") > strings.Index(badan, "UJI-A") || !strings.Contains(badan, `"terpotong":false`) {
		t.Errorf("rate list: %d %s", kode, badan)
	}
	if n := u.cacah(t, "RATE_LIFE", ""); n != 3 || sidik() != awal {
		t.Errorf("view rate tersentuh: %d baris, sidik %s → %s", n, awal, sidik())
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
