//go:build db

// Seam HTTP klausul (tiket 08) terhadap skema uji Oracle NYATA - satu tabel
// yang memuat banyak jenis dengan pola NULL berbeda hanya terbukti di Oracle.
package handlers_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/treatycontractout/backend/repository"
	"nusantarare/uji/skemauji"
)

func TestKlausulLingkaranPenuh(t *testing.T) {
	u, bersihkan := serverTCO(t)
	defer bersihkan()
	// Nama induk berkata " QS " -> pilihan anak QS (OR), QS (R/I), ORS (`TreatyContractSetReinsTypeList`).
	u.isiJenis([]skemauji.JenisReasuransiUji{
		{ID: "10003", Note: "UJI QS TREATY", Tipe: "1", Flag: "active"},
		{ID: "10005", Note: "UJI SURPLUS", Tipe: "2", Flag: "active"},
	})
	if err := skemauji.IsiJenisKlausulTCO(u.ctx, u.sqlDBMentah(), u.skema, []skemauji.JenisKlausulUji{
		{ID: "10001", DescName: "UJI TREATY LIMIT", IsXOL: "0"}, {ID: "10009", DescName: "UJI EPI", IsXOL: "0"},
		{ID: "10013", DescName: "UJI EXCLUSION", IsXOL: "0"}, {ID: "10017", DescName: "UJI LIMIT MB", IsXOL: "1"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := skemauji.IsiPilihanKlausulTCO(u.ctx, u.sqlDBMentah(), u.skema, repository.MasterOccupationTCO,
		map[string]string{"UJI-O1": "UJI OKUPASI"}); err != nil {
		t.Fatal(err)
	}
	// Tiket 11: induk EPI dan anaknya berkurs.
	isiKursUji(t, u)
	_, badan := u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun", badanTahun("2026-01-01", "2026-12-31"), true)
	var tahun tahunJSON
	_ = json.Unmarshal([]byte(badan), &tahun)
	dasar := "/api/treaty-contract-out/tahun/" + tahun.ID + "/klausul"

	kode, badan := u.minta(t, http.MethodGet, "/api/treaty-contract-out/jenis-klausul?isXol=0", nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"total":3`) || !strings.Contains(badan, `"jenis":"EpiList"`) {
		t.Fatalf("jenis: %d %s", kode, badan)
	}
	// Induk EPI: identitas '1' + 7 digit; sentinel "00"; Usd = Rp / Kurs (tiket 11).
	kode, badan = u.minta(t, http.MethodPost, dasar, map[string]any{"descId": "10009",
		"medan": map[string]string{"ReinsTypeID": "10003", "Rp": "1000000,5"}}, true)
	if kode != http.StatusOK {
		t.Fatalf("induk: %d %s", kode, badan)
	}
	var h struct {
		Klausul struct {
			ID, ParentReinsTypeID, Kurs string
			Medan                       map[string]string
		}
	}
	_ = json.Unmarshal([]byte(badan), &h)
	if len(h.Klausul.ID) != 8 || h.Klausul.ParentReinsTypeID != "00" || h.Klausul.Medan["Rp"] != "1000000.5" ||
		h.Klausul.Medan["Usd"] != "64.51512072" || h.Klausul.Kurs != "15500.25" {
		t.Fatalf("induk: %+v", h)
	}
	// Dobel, medan wajib, jenis ditahan.
	if kode, badan := u.minta(t, http.MethodPost, dasar, map[string]any{"descId": "10009",
		"medan": map[string]string{"ReinsTypeID": "10003", "Rp": "1"}}, true); kode != http.StatusConflict ||
		!strings.Contains(badan, "Data has already been entered") {
		t.Errorf("dobel: %d %s", kode, badan)
	}
	if kode, badan := u.minta(t, http.MethodPost, dasar, map[string]any{"descId": "10009",
		"medan": map[string]string{"ReinsTypeID": "10005"}}, true); kode != http.StatusUnprocessableEntity ||
		!strings.Contains(badan, "Rp") {
		t.Errorf("wajib: %d %s", kode, badan)
	}
	// MB Capacity mengikuti XML [keputusan work owner 02-10-2026]: tanpa wajib-isi; MORERP/MOREUSD NUMBER.
	if kode, badan := u.minta(t, http.MethodPost, dasar, map[string]any{"descId": "10017", "medan": map[string]string{
		"ID_Occupation": "01", "MoreRp": "5000000.5", "MoreUsd": "325.25", "TerritorialLimit": "UJI GRUP"}}, true); kode != http.StatusOK ||
		!strings.Contains(badan, `"Occupation":"RESIDENTIAL RISK"`) || !strings.Contains(badan, `"MoreRp":"5000000.5"`) {
		t.Errorf("LimitMB: %d %s", kode, badan)
	}
	if kode, _ := u.minta(t, http.MethodPost, dasar, map[string]any{"descId": "10002", "medan": map[string]string{}}, true); kode != http.StatusUnprocessableEntity {
		t.Errorf("Portfolio tetap ditahan: %d", kode)
	}
	// Anak EpiList: Rp/Usd turunan; sembilan kolom khusus induk NULL di Oracle.
	kode, badan = u.minta(t, http.MethodPost, dasar, map[string]any{"descId": "10009", "anak": true,
		"parentReinsTypeId": "10003", "medan": map[string]string{"ReinsTypeID": "10028", "Pct": "25"}}, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"Rp":"250000.12500000"`) || !strings.Contains(badan, `"Usd":"16.12878018"`) ||
		!strings.Contains(badan, `"reinsTypeName":"QS (OR)"`) {
		t.Fatalf("anak: %d %s", kode, badan)
	}
	var anak struct{ Klausul struct{ ID string } }
	_ = json.Unmarshal([]byte(badan), &anak)
	var nullInduk int
	if err := u.sqlDBMentah().QueryRowContext(u.ctx, `SELECT COUNT(*) FROM `+u.skema+`.PROPORTIONALARRG
		 WHERE ID = :1 AND ID_OCCUPATION IS NULL AND OCCUPATION IS NULL AND ID_CLAUSE IS NULL AND CLAUSE IS NULL
		   AND TREATYLIMIT IS NULL AND COINS_MIN IS NULL AND COINS_MAX IS NULL AND MORERP IS NULL AND MOREUSD IS NULL`,
		anak.Klausul.ID).Scan(&nullInduk); err != nil || nullInduk != 1 {
		t.Errorf("anak tidak menulis NULL pada sembilan kolom induk: %d %v", nullInduk, err)
	}
	if kode, _ := u.minta(t, http.MethodPost, dasar, map[string]any{"descId": "10009", "anak": true,
		"parentReinsTypeId": "10003", "medan": map[string]string{"ReinsTypeID": "10004", "Pct": "75.1"}}, true); kode != http.StatusUnprocessableEntity {
		t.Errorf("total anak > 100: %d", kode)
	}
	kode, badan = u.minta(t, http.MethodGet, dasar+"?descId=10009&induk=10003", nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"totalPct":"25"`) {
		t.Errorf("daftar anak: %d %s", kode, badan)
	}
	// Exclusion Occupation: nama dari master.
	kode, badan = u.minta(t, http.MethodPost, dasar, map[string]any{"descId": "10013", "subjenis": "Occupation",
		"medan": map[string]string{"ID_Occupation": "UJI-O1", "Occupation": "KARANGAN", "Line": "A", "Usd": "1", "Rp": "15000"}}, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"Occupation":"UJI OKUPASI"`) {
		t.Errorf("exclusion: %d %s", kode, badan)
	}
}
