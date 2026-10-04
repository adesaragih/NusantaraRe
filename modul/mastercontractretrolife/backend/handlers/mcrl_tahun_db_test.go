//go:build db

package handlers_test

// Seam HTTP simpan tahun treaty terhadap skema uji Oracle NYATA (paket 2).

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func (u *ujiDB) kirim(t *testing.T, metode, jalur, badan string) (int, string) {
	t.Helper()
	req, err := http.NewRequest(metode, u.srv.URL+jalur, bytes.NewBufferString(badan))
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

func TestDBTahunBaruDariSequenceLaluUpsertDikunciID(t *testing.T) {
	u := pasangDB(t)
	kode, badan := u.kirim(t, "POST", "/api/master-contract-retro-life/tahun",
		`{"treatyYear":"2026","underwritingYear":"2026","startDate":"2026-01-01","endDate":"2026-12-31"}`)
	if kode != http.StatusOK || !strings.Contains(badan, `"id":"1000044"`) {
		t.Fatalf("POST: %d %s", kode, badan)
	}
	if n := u.cacah(t, "TREATYYEAR_LIFE", "ID = '1000044' AND USERID = 'UJI-PELAKU' AND TGLUPDATE IS NOT NULL "+
		"AND STARTDATE = DATE '2026-01-01'"); n != 1 {
		t.Errorf("baris baru: %d", n)
	}
	u.exec(t, `INSERT INTO `+u.skema+`.TREATYBUSINESS_LIFE (ID, TREATYYEARID, TREATYYEAR) VALUES ('1000001', '1000044', '2026')`)
	kode, badan = u.kirim(t, "PUT", "/api/master-contract-retro-life/tahun/1000044",
		`{"treatyYear":"2027","underwritingYear":"2026","startDate":"2027-01-01","endDate":"2027-12-31"}`)
	if kode != http.StatusOK {
		t.Fatalf("PUT: %d %s", kode, badan)
	}
	if n := u.cacah(t, "TREATYYEAR_LIFE", ""); n != 1 {
		t.Errorf("upsert dikunci ID: %d baris, mau 1", n)
	}
	if got := u.teks(t, `SELECT TREATYYEAR FROM {s}.TREATYBUSINESS_LIFE WHERE ID = '1000001'`); got != "2027" {
		t.Errorf("salinan TREATYYEAR business (K4): %q", got)
	}
}
