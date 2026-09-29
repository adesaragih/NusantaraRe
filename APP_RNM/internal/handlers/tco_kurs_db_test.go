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

// Lanjutan 6: tanggal master kurs diurai ORACLE seperti `GetMasterKursList` -
// jam lima angka (`T00000.000`, 11 baris STARTDATE di DEV) diterima; bentuk
// yang Oracle tolak (huruf, bulan 13, panjang salah, kosong) disaring per baris
// dan dicacah, tidak mematikan kurs yang berlaku.
func TestKursTanggalDiuraiOracle(t *testing.T) {
	u, bersihkan := serverTCO(t)
	defer bersihkan()
	if err := skemauji.IsiMataUangTCO(u.ctx, u.sqlDBMentah(), u.skema, map[string]string{"USD": "UJI-USD"}); err != nil {
		t.Fatal(err)
	}
	baris := []skemauji.KursUji{{ToIDR: "15500", StartDate: "20260101T00000.000 GMT", EndDate: "20261231T000000.000 GMT"}}
	for _, buruk := range []string{"2019A801T000000.000 GMT", "20191301T000000.000 GMT", "2019T000000.000 GMT", ""} {
		baris = append(baris, skemauji.KursUji{ToIDR: "1", StartDate: buruk, EndDate: "20191231T000000.000 GMT"})
	}
	for i := range baris {
		baris[i].IDCurrency, baris[i].Currency, baris[i].Quarter = "UJI-USD", "USD", "0"
	}
	if err := skemauji.IsiKursTCO(u.ctx, u.sqlDBMentah(), u.skema, baris); err != nil {
		t.Fatal(err)
	}
	_, badan := u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun", badanTahun("2026-01-01", "2026-12-31"), true)
	var tahun tahunJSON
	_ = json.Unmarshal([]byte(badan), &tahun)
	kode, badan := u.minta(t, http.MethodGet, "/api/treaty-contract-out/tahun/"+tahun.ID+"/kurs", nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"kurs":"15500"`) || !strings.Contains(badan, `"mulai":"2026-01-01"`) ||
		!strings.Contains(badan, `"barisMasterDitolak":4`) {
		t.Errorf("jam lima angka + empat ditolak: %d %s", kode, badan)
	}
	// Tanpa baris berlaku, penolakan itu 503 berkata-kata - bukan "tidak ada kurs".
	_, badan = u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun", badanTahun("2019-08-01", "2019-12-31"), true)
	var tahun19 tahunJSON
	_ = json.Unmarshal([]byte(badan), &tahun19)
	kode, badan = u.minta(t, http.MethodGet, "/api/treaty-contract-out/tahun/"+tahun19.ID+"/kurs", nil, true)
	if kode != http.StatusServiceUnavailable || !strings.Contains(badan, "ditolak Oracle") || !strings.Contains(badan, "2019A801T000000.000 GMT") {
		t.Errorf("hanya ditolak: %d %s", kode, badan)
	}
}
