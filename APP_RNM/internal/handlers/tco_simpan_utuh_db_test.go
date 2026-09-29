//go:build db

// Seam HTTP simpan utuh (tiket 09) terhadap skema uji Oracle NYATA -
// atomisitas lintas enam tabel adalah satu-satunya hal yang diuji tiket ini,
// dan ia hanya berperilaku benar pada basis data sungguhan.
package handlers_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/internal/repository"
	"nusantarare/internal/repository/skemauji"
)

// hitungUtuh - jumlah baris lima tabel + jejak, dan LAST_NUMBER enam sequence
// (NOCACHE: LAST_NUMBER = nomor berikut, sehingga nomor terpakai terlihat).
func hitungUtuh(t *testing.T, u *ujiTCO) map[string]int64 {
	t.Helper()
	hasil := map[string]int64{}
	for _, tabel := range []string{"T_TREATYCONTRACT", "T_TREATYREINSURER", "T_MTREATYSECURITY", "T_TREATYBUSINESS",
		"T_PROPORTIONALARRG", "T_TREATYCO_JEJAK"} {
		var n int64
		if err := u.sqlDBMentah().QueryRowContext(u.ctx, `SELECT COUNT(*) FROM `+u.skema+`.`+tabel).Scan(&n); err != nil {
			t.Fatal(err)
		}
		hasil[tabel] = n
	}
	rows, err := u.sqlDBMentah().QueryContext(u.ctx, `SELECT SEQUENCE_NAME, LAST_NUMBER FROM ALL_SEQUENCES
		 WHERE SEQUENCE_OWNER = UPPER(:1) AND SEQUENCE_NAME IN ('SEQ_T_TREATYCONTRACT', 'SEQ_T_TREATYREINSURER',
		 'SEQ_T_MTREATYSECURITY', 'SEQ_T_TREATYBUSINESS', 'SEQ_T_PROPORTIONALARRG', 'SEQ_T_TREATYCO_JEJAK')`, u.skema)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var nama string
		var n int64
		if err := rows.Scan(&nama, &n); err != nil {
			t.Fatal(err)
		}
		hasil[nama] = n
	}
	if len(hasil) != 12 {
		t.Fatalf("hitungan tidak lengkap: %v", hasil)
	}
	return hasil
}

func TestSimpanUtuhAtomikLintasEnamTabel(t *testing.T) {
	u, bersihkan := serverTCO(t)
	defer bersihkan()
	u.isiJenis([]skemauji.JenisReasuransiUji{{ID: "10003", Note: "UJI QUOTA SHARE", Tipe: "1", Flag: "active"}})
	u.isiAgent([]skemauji.AgentUji{
		{ID: "UJI-R1", ClientName: "UJI REAS SATU", ClientID: "UJI-C1", StatusActive: "1"},
		{ID: "UJI-R2", ClientName: "UJI REAS DUA", ClientID: "UJI-C2", StatusActive: "1"},
		{ID: "UJI-R3", ClientName: "UJI REAS TIGA", ClientID: "UJI-C3", StatusActive: "1"},
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
	_, badan = u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun/"+tahun.ID+"/kontrak",
		map[string]string{"reinsTypeId": "10003", "treatyStartDate": "2026-01-01", "treatyEndDate": "2027-01-01"}, true)
	var k kontrakJSON
	_ = json.Unmarshal([]byte(badan), &k)
	jalur := "/api/treaty-contract-out/tahun/" + tahun.ID + "/kontrak/" + k.ID + "/utuh"
	epiInduk := map[string]any{"descId": "10009", "medan": map[string]string{"ReinsTypeID": "10003", "Rp": "1000000"}}
	bundel := func(klausul ...map[string]any) map[string]any {
		return map[string]any{
			"kontrak": map[string]string{"reinsTypeId": "10003", "treatyStartDate": "2026-01-01", "treatyEndDate": "2026-12-31"},
			"reinsurer": []map[string]any{
				{"reinsurerId": "UJI-R1", "pctShare": "40", "ricomm": "10",
					"security": []map[string]string{{"reasSecurity": "UJI-R2", "pctShare": "20"}, {"reasSecurity": "UJI-R3", "pctShare": "20"}}},
				{"reinsurerId": "UJI-R2", "pctShare": "60", "ricomm": "10"},
			},
			"business": []map[string]string{{"bizCode": "UJI-B1", "isActive": "1"}},
			"klausul":  klausul,
		}
	}
	sebelum := hitungUtuh(t, u)

	// AC 37: gagal pada klausul ke-3 (dobel klausul ke-1) -> NOL perubahan di
	// keenam tabel, NOL nomor sequence terpakai, kontrak tidak berubah.
	kode, badan := u.minta(t, http.MethodPut, jalur, bundel(epiInduk,
		map[string]any{"descId": "10009", "anak": true, "parentReinsTypeId": "10003", "medan": map[string]string{"ReinsTypeID": "10003", "Pct": "50"}},
		epiInduk), true)
	if kode != http.StatusConflict || !strings.Contains(badan, "klausul ke-3") || !strings.Contains(badan, "Data sudah pernah di Input") ||
		repository.PolaIdentitasSementaraTCO().MatchString(badan) {
		t.Fatalf("gagal ke-3: %d %s", kode, badan)
	}
	sesudah := hitungUtuh(t, u)
	for kunci, n := range sebelum {
		if sesudah[kunci] != n {
			t.Errorf("%s berubah %d -> %d sesudah kegagalan", kunci, n, sesudah[kunci])
		}
	}
	_, badan = u.minta(t, http.MethodGet, "/api/treaty-contract-out/tahun/"+tahun.ID+"/kontrak", nil, true)
	if !strings.Contains(badan, "01-01-2027") && !strings.Contains(badan, "2027-01-01") {
		t.Errorf("kontrak berubah sesudah kegagalan: %s", badan)
	}

	// Sukses: satu transaksi, identitas tetap, security menunjuk reinsurer tetap.
	kode, badan = u.minta(t, http.MethodPut, jalur, bundel(epiInduk,
		map[string]any{"descId": "10009", "anak": true, "parentReinsTypeId": "10003", "medan": map[string]string{"ReinsTypeID": "10003", "Pct": "50"}}), true)
	if kode != http.StatusOK || !strings.Contains(badan, `"status":"1"`) || repository.PolaIdentitasSementaraTCO().MatchString(badan) {
		t.Fatalf("sukses: %d %s", kode, badan)
	}
	var h struct {
		Reinsurer []struct {
			Reinsurer struct{ ID string }
			Security  []struct{ ID, ReasID string }
		}
		Klausul []struct{ ID string }
	}
	_ = json.Unmarshal([]byte(badan), &h)
	if len(h.Reinsurer) != 2 || len(h.Reinsurer[0].Security) != 2 || h.Reinsurer[0].Security[0].ReasID != h.Reinsurer[0].Reinsurer.ID ||
		len(h.Reinsurer[0].Reinsurer.ID) != 7 || len(h.Klausul) != 2 || len(h.Klausul[0].ID) != 8 {
		t.Fatalf("identitas: %+v", h)
	}
	akhir := hitungUtuh(t, u)
	for tabel, tambah := range map[string]int64{"T_TREATYREINSURER": 2, "T_MTREATYSECURITY": 2, "T_TREATYBUSINESS": 1, "T_PROPORTIONALARRG": 2} {
		if akhir[tabel]-sebelum[tabel] != tambah {
			t.Errorf("%s bertambah %d, mau %d", tabel, akhir[tabel]-sebelum[tabel], tambah)
		}
	}
	if akhir["SEQ_T_TREATYREINSURER"]-sebelum["SEQ_T_TREATYREINSURER"] != 2 ||
		akhir["SEQ_T_PROPORTIONALARRG"]-sebelum["SEQ_T_PROPORTIONALARRG"] != 2 {
		t.Errorf("nomor sequence terpakai tidak sama dengan baris baru: %v -> %v", sebelum, akhir)
	}
	var sementara int
	if err := u.sqlDBMentah().QueryRowContext(u.ctx, `SELECT
		 (SELECT COUNT(*) FROM `+u.skema+`.T_TREATYCO_JEJAK WHERE ID LIKE 'S%T' OR BARIS_ID LIKE 'S%T') +
		 (SELECT COUNT(*) FROM `+u.skema+`.T_TREATYREINSURER WHERE ID LIKE 'S%T') +
		 (SELECT COUNT(*) FROM `+u.skema+`.T_MTREATYSECURITY WHERE ID LIKE 'S%T' OR REAS_ID LIKE 'S%T') +
		 (SELECT COUNT(*) FROM `+u.skema+`.T_PROPORTIONALARRG WHERE ID LIKE 'S%T') +
		 (SELECT COUNT(*) FROM `+u.skema+`.T_TREATYCO_JEJAK WHERE REGEXP_LIKE(KETERANGAN, 'S[0-9]{16}T')) FROM DUAL`).Scan(&sementara); err != nil || sementara != 0 {
		t.Errorf("identitas sementara tersisa: %d %v", sementara, err)
	}
}
