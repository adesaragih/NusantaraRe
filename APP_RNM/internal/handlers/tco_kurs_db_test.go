//go:build db

// Seam HTTP kurs (tiket 11) terhadap skema uji Oracle NYATA - pemilihan baris
// kurs menurut periode dan presisi konversinya hanya terbukti di basis data.
package handlers_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/internal/repository/skemauji"
)

// isiKursUji - satu pengenal USD dan tiga baris kurs: yang berlaku 2026,
// quarter lain (diabaikan), mata uang lain (diabaikan).
func isiKursUji(t *testing.T, u *ujiTCO) {
	t.Helper()
	if err := skemauji.IsiMataUangTCO(u.ctx, u.sqlDBMentah(), u.skema, map[string]string{"USD": "UJI-USD", "EUR": "UJI-EUR"}); err != nil {
		t.Fatal(err)
	}
	if err := skemauji.IsiKursTCO(u.ctx, u.sqlDBMentah(), u.skema, []skemauji.KursUji{
		{ToIDR: "15500,25", IDCurrency: "UJI-USD", Currency: "USD", Quarter: "0",
			StartDate: "20260101T000000.000 GMT", EndDate: "20261231T000000.000 GMT"},
		{ToIDR: "1", IDCurrency: "UJI-USD", Currency: "USD", Quarter: "1",
			StartDate: "20260101T000000.000 GMT", EndDate: "20261231T000000.000 GMT"},
		{ToIDR: "2", IDCurrency: "UJI-EUR", Currency: "EUR", Quarter: "0",
			StartDate: "20260101T000000.000 GMT", EndDate: "20261231T000000.000 GMT"},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestKursLingkaranPenuh(t *testing.T) {
	u, bersihkan := serverTCO(t)
	defer bersihkan()
	isiKursUji(t, u)
	_, badan := u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun", badanTahun("2026-01-01", "2026-12-31"), true)
	var tahun tahunJSON
	_ = json.Unmarshal([]byte(badan), &tahun)
	dasar := "/api/treaty-contract-out/tahun/" + tahun.ID + "/kurs"

	kode, badan := u.minta(t, http.MethodGet, dasar, nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"kurs":"15500.25"`) || !strings.Contains(badan, `"idCurrency":"UJI-USD"`) {
		t.Fatalf("kurs: %d %s", kode, badan)
	}
	kode, badan = u.minta(t, http.MethodGet, dasar+"/konversi?dari=Rp&nilai=31000500", nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"usd":"2000.00000000"`) {
		t.Errorf("konversi Rp: %d %s", kode, badan)
	}
	kode, badan = u.minta(t, http.MethodGet, dasar+"/konversi?dari=Usd&nilai=2000", nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"rp":"31000500.00000000"`) {
		t.Errorf("konversi Usd: %d %s", kode, badan)
	}
	// ADR-0015: tahun tanpa kurs -> 422 + pesan VERBATIM.
	_, badan = u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun", badanTahun("2027-01-01", "2027-12-31"), true)
	var tahun27 tahunJSON
	_ = json.Unmarshal([]byte(badan), &tahun27)
	kode, badan = u.minta(t, http.MethodGet, "/api/treaty-contract-out/tahun/"+tahun27.ID+"/kurs", nil, true)
	if kode != http.StatusUnprocessableEntity || !strings.Contains(badan, "Tidak ada Nilai Kurs di Tahun : 2027") {
		t.Errorf("tanpa kurs: %d %s", kode, badan)
	}
	// Master kurs tidak ditulis: jumlah barisnya tetap.
	var n int
	if err := u.sqlDBMentah().QueryRowContext(u.ctx, `SELECT COUNT(*) FROM `+u.skema+`.TREATYEXCHANGEYEARLY`).Scan(&n); err != nil || n != 3 {
		t.Errorf("master kurs berubah: %d %v", n, err)
	}
}
