//go:build db

// Seam HTTP kaskade hapus (tiket 10) terhadap skema uji Oracle NYATA - kaskade
// yang MENGECUALIKAN satu tabel adalah pernyataan tentang basis data. Uji dua
// arah: tiga anak hilang, klausul masih ada.
package handlers_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/uji/skemauji"
)

func TestKaskadeHapusKontrakKlausulTetapHidup(t *testing.T) {
	u, bersihkan := serverTCO(t)
	defer bersihkan()
	u.isiJenis([]skemauji.JenisReasuransiUji{
		{ID: "10003", Note: "UJI QUOTA SHARE", Tipe: "1", Flag: "active"},
		{ID: "10005", Note: "UJI SURPLUS", Tipe: "2", Flag: "active"},
	})
	u.isiAgent([]skemauji.AgentUji{
		{ID: "UJI-R1", ClientName: "UJI REAS SATU", ClientID: "UJI-C1", StatusActive: "1"},
		{ID: "UJI-R2", ClientName: "UJI REAS DUA", ClientID: "UJI-C2", StatusActive: "1"},
	})
	if err := skemauji.IsiBusinessTCO(u.ctx, u.sqlDBMentah(), u.skema, []skemauji.BusinessUji{{ID: "UJI-B1", Note: "UJI BISNIS", BusinessGroupID: "UJI-G1"}}); err != nil {
		t.Fatal(err)
	}
	if err := skemauji.IsiJenisKlausulTCO(u.ctx, u.sqlDBMentah(), u.skema, []skemauji.JenisKlausulUji{{ID: "10009", DescName: "UJI EPI", IsXOL: "0"}}); err != nil {
		t.Fatal(err)
	}
	isiKursUji(t, u)
	_, badan := u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun", badanTahun("2026-01-01", "2026-12-31"), true)
	var tahun tahunJSON
	_ = json.Unmarshal([]byte(badan), &tahun)
	dasarTahun := "/api/treaty-contract-out/tahun/" + tahun.ID
	kontrak := func(jenis string) string {
		_, b := u.minta(t, http.MethodPost, dasarTahun+"/kontrak",
			map[string]string{"reinsTypeId": jenis, "treatyStartDate": "2026-01-01", "treatyEndDate": "2027-01-01"}, true)
		var k kontrakJSON
		_ = json.Unmarshal([]byte(b), &k)
		return k.ID
	}
	a, b := kontrak("10003"), kontrak("10005")
	// OQ-TCO-21: tahun LAIN berteks tahun + grup sama, kontrak berjenis sama ->
	// kombinasi (2026, 10001, 10003) dipakai bersama; popup menyebut 1 kontrak lain.
	_, badanLain := u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun", badanTahun("2026-07-01", "2027-06-30"), true)
	var tahunLain tahunJSON
	_ = json.Unmarshal([]byte(badanLain), &tahunLain)
	if kode, bd := u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun/"+tahunLain.ID+"/kontrak",
		map[string]string{"reinsTypeId": "10003", "treatyStartDate": "2026-07-01", "treatyEndDate": "2027-07-01"}, true); kode != http.StatusOK {
		t.Fatalf("kontrak tahun lain: %d %s", kode, bd)
	}
	reas := func(kid, agen, share string) string {
		_, bd := u.minta(t, http.MethodPost, dasarTahun+"/kontrak/"+kid+"/reinsurer",
			map[string]string{"reinsurerId": agen, "pctShare": share, "ricomm": "0"}, true)
		var r struct{ Reinsurer struct{ ID string } }
		_ = json.Unmarshal([]byte(bd), &r)
		return r.Reinsurer.ID
	}
	ra := reas(a, "UJI-R1", "40")
	reas(a, "UJI-R2", "60")
	rb := reas(b, "UJI-R1", "100")
	if kode, bd := u.minta(t, http.MethodPost, dasarTahun+"/kontrak/"+a+"/reinsurer/"+ra+"/security",
		map[string]string{"reasSecurity": "UJI-R2", "pctShare": "10"}, true); kode != http.StatusOK {
		t.Fatalf("security: %d %s", kode, bd)
	}
	if kode, bd := u.minta(t, http.MethodPost, dasarTahun+"/kontrak/"+a+"/business",
		map[string]string{"bizCode": "UJI-B1", "isActive": "1"}, true); kode != http.StatusOK {
		t.Fatalf("business: %d %s", kode, bd)
	}
	// Baris bisnis lama tanpa TREATYYEARID - kaskade wajib tahan NULL.
	if _, err := u.sqlDBMentah().ExecContext(u.ctx, `INSERT INTO `+u.skema+`.TREATYBUSINESS
		 (ID, ISACTIVE, TREATYYEAR, TREATYYEARID, TREATYGROUPID, REINSTYPEID, BIZCODE) VALUES ('UJI-BN1', '1', '2026', NULL, '10001', '10003', 'UJI-B9')`); err != nil {
		t.Fatal(err)
	}
	if kode, bd := u.minta(t, http.MethodPost, dasarTahun+"/klausul", map[string]any{"descId": "10009",
		"medan": map[string]string{"ReinsTypeID": "10003", "Rp": "1000"}}, true); kode != http.StatusOK {
		t.Fatalf("klausul: %d %s", kode, bd)
	}

	// Popup: jumlah tiap jenis + klausul yang tetap hidup.
	kode, badan := u.minta(t, http.MethodGet, dasarTahun+"/kontrak/"+a+"/dampak-hapus", nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"reinsurer":2`) || !strings.Contains(badan, `"security":1`) ||
		!strings.Contains(badan, `"business":2`) || !strings.Contains(badan, `"klausulTetap":1`) || !strings.Contains(badan, `"bersama":1`) {
		t.Fatalf("dampak: %d %s", kode, badan)
	}
	// Angka lain dari popup -> 409, tidak ada yang terhapus.
	if kode, _ := u.minta(t, http.MethodDelete, dasarTahun+"/kontrak/"+a+"?reinsurer=2&security=1&business=1&bersama=1", nil, true); kode != http.StatusConflict {
		t.Errorf("angka berubah: %d", kode)
	}
	kode, badan = u.minta(t, http.MethodDelete, dasarTahun+"/kontrak/"+a+"?reinsurer=2&security=1&business=2&bersama=1", nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, "Data Berhasil di Hapus") {
		t.Fatalf("hapus: %d %s", kode, badan)
	}
	hitung := func(q string, args ...any) int {
		var n int
		if err := u.sqlDBMentah().QueryRowContext(u.ctx, q, args...).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}
	s := u.skema
	if n := hitung(`SELECT COUNT(*) FROM `+s+`.TREATYREINSURER WHERE REINSTYPEID = '10003'`) +
		hitung(`SELECT COUNT(*) FROM `+s+`.MTREATYSECURITY`) +
		hitung(`SELECT COUNT(*) FROM `+s+`.TREATYBUSINESS WHERE REINSTYPEID = '10003'`) +
		hitung(`SELECT COUNT(*) FROM `+s+`.TREATYCONTRACT WHERE ID = :1`, a); n != 0 {
		t.Errorf("anak tersisa: %d", n)
	}
	// Klausul TETAP HIDUP dan terbaca (AC 44); kontrak lain tidak tersentuh.
	kode, badan = u.minta(t, http.MethodGet, dasarTahun+"/klausul?descId=10009&induk=00", nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"total":1`) {
		t.Errorf("klausul ikut terhapus: %d %s", kode, badan)
	}
	if n := hitung(`SELECT COUNT(*) FROM `+s+`.TREATYREINSURER WHERE ID = :1`, rb); n != 1 {
		t.Errorf("reinsurer kontrak lain terhapus")
	}
	// Hapus reinsurer kontrak lain: popup 0 security; Ya -> hilang.
	kode, badan = u.minta(t, http.MethodGet, dasarTahun+"/kontrak/"+b+"/reinsurer/"+rb+"/dampak-hapus", nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"security":0`) {
		t.Errorf("dampak reinsurer: %d %s", kode, badan)
	}
	if kode, _ := u.minta(t, http.MethodDelete, dasarTahun+"/kontrak/"+b+"/reinsurer/"+rb+"?security=0", nil, true); kode != http.StatusOK ||
		hitung(`SELECT COUNT(*) FROM `+s+`.TREATYREINSURER WHERE ID = :1`, rb) != 0 {
		t.Errorf("hapus reinsurer: %d", kode)
	}
}
