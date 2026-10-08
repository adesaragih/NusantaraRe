package handlers_test

// Spreading klaim terisi dari polis saat polis dipilih (keputusan work owner 08-10-2026), tabel bawah dari
// PROPORTIONALARRG; kasus baru dengan spreading itu tidak lagi ditolak ProteksiData langkah 5 saat Save to issue RNM.
// Dropdown Treaty Type spreading (Input Acceptation) menawarkan spreading polis - koreksi atas `[dugaan]` SpreadingList
// master yang dibantah data DEV. Fixture UJI-.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/modul/claimprop/backend/handlers"
	"nusantarare/modul/claimprop/backend/models"
)

func TestSpreadingTerisiDariPolisLolosSaveToIssueRNM(t *testing.T) {
	u := baruUji(t)
	u.a.JenisReas["UJI-INDUK"] = "UJI QS INDUK TRT"
	u.a.SpreadingPolisMap[polisUji] = []models.SpreadingPolis{
		{TreatyType: "UJI-INDUK", SharePercentage: "100", CurrencyID: matauang, Currency: "UJA"},
	}
	u.a.AnakSpreadingMap["UJI-INDUK|2026|UJI-TG"] = []models.AnakSpreading{
		{ReinsTypeID: "UJI-RI", ReinsTypeName: "QS (R/I)", Pct: "60"}, {ReinsTypeID: "UJI-OR", ReinsTypeName: "QS (OR)", Pct: "40"},
	}
	id := u.buat()
	langkah := func(aksi string, n int, param string, m map[string]string) map[string]any {
		t.Helper()
		kode, out := u.aksi(id, admin, "", aksi, n, param, m)
		u.wajib(kode, http.StatusOK, out, aksi)
		return out
	}
	langkah("SetValueToClaim", 0, masterUji+"|UJI-TG|UJI-COB", nil)
	langkah("CheckNoPolicy", 0, polisUji, nil)

	h := u.g.Halaman(id)
	atas, bawah := h.AmbilDaftar(models.DaftarSpreading), h.AmbilDaftar(models.DaftarBreakQS)
	if len(atas) != 1 || atas[0]["TreatyType"] != "UJI-INDUK" || atas[0]["TreatyName"] != "UJI QS INDUK TRT" ||
		atas[0]["SharePercentage"] != "100" {
		t.Fatalf("spreading atas dari polis: %v", atas)
	}
	if len(bawah) != 2 || bawah[0]["TreatyType"] != "UJI-RI" || bawah[1]["SharePercentage"] != "40" {
		t.Fatalf("spreading bawah dari PROPORTIONALARRG: %v", bawah)
	}

	// dropdown Treaty Type spreading = spreading polis berlabel nama treaty
	r := httptest.NewRecorder()
	q := httptest.NewRequest(http.MethodGet, handlers.Prefix+"/kasus/"+id+"/pilihan/"+models.SumberSpreading, nil)
	q.Header.Set("X-Pelaku", admin)
	u.srv.ServeHTTP(r, q)
	var opsi []models.Pilihan
	if r.Code != http.StatusOK || json.Unmarshal(r.Body.Bytes(), &opsi) != nil || len(opsi) != 1 ||
		opsi[0].Nilai != "UJI-INDUK" || opsi[0].Label != "UJI QS INDUK TRT" {
		t.Fatalf("pilihan spreading: HTTP %d %s", r.Code, r.Body.String())
	}

	// kelengkapan lain seperti alur penuh, lalu Save to issue RNM: tidak ada pesan "please Fill SpreadingList"
	langkah("GetNameCauseofLoss", 0, "UJI-S1", nil)
	langkah("SetOutstanding", 0, "", map[string]string{
		models.CD + "PolicyData.StartDateTime": "2026-01-01", models.CD + "PolicyData.EndDateTime": "2026-12-31",
		models.CD + "DateOfLoss": "2026-05-01 08:00:00", models.CD + "ReportDate": "2026-05-02",
		models.CD + "DateReceived": "2026-05-03 09:00:00", models.CD + "ReporterName": "UJI-PELAPOR",
		models.CD + "ReporterTelp": "1", models.CD + "ReportAddress": "UJI-ALAMAT", models.CD + "Location": "uji-lokasi",
		models.CD + "Province": "UJI-PROVINSI", models.CD + "ConsultantID": "UJI-ADJ", models.CD + "AppointedADJID": "UJI-ADJ",
		models.CD + "ShareCeding": "100", models.CD + "Occupation": "uji-okupasi"})
	langkah("AddInterest", 0, "", nil)
	langkah("SetCurencyInterest", 1, "", map[string]string{models.JalurAnak(models.DaftarInterest, 1, "CurrencyID"): matauang,
		models.JalurAnak(models.DaftarInterest, 1, "ObjectName"): "UJI-OBJ"})
	langkah("CountTotalInsterest", 1, "", map[string]string{models.JalurAnak(models.DaftarInterest, 1, "TSIPerObject"): "1000"})
	langkah("AddListClaimAmount", 0, "", nil)
	langkah("AddLossAllocation", 0, "", nil)
	langkah("SetNameTreaty", 1, "", map[string]string{models.JalurAnak(models.DaftarLossAlloc, 1, "TreatyName"): "QUOTA SHARE"})
	langkah("CountPersen", 1, "", map[string]string{models.JalurAnak(models.DaftarLossAlloc, 1, "SharePercentage"): "100"})
	langkah("AddEstimation", 0, "", nil)
	langkah("CountEstimation", 1, "", map[string]string{models.JalurAnak(models.DaftarEstimasi, 1, "Type"): "1"})
	kode, out := u.aksi(id, admin, "", "SaveOutstanding", 0, "", nil)
	if strings.Contains(strings.Join(teks(out["pesan"]), ";"), models.PesanIsiSpreading) {
		t.Fatalf("spreading dari polis masih ditolak ProteksiData 5: %v", out["pesan"])
	}
	u.wajib(kode, http.StatusOK, out, "Save to issue RNM dengan spreading dari polis")
}
