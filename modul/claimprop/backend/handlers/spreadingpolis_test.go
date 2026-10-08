package handlers_test

// Spreading klaim lewat HTTP (keputusan work owner 08-10-2026): terisi dari polis saat polis dipilih, tabel bawah dari
// SpreadingList master, Add / Delete aktif, Treaty Type terkunci sesudah terisi. Kasus baru dengan spreading lolos
// ProteksiData langkah 5; semua baris dihapus = ditolak lagi. Fixture UJI-.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/modul/claimprop/backend/handlers"
	"nusantarare/modul/claimprop/backend/models"
)

func TestSpreadingDariPolisTambahHapusLolosSaveToIssueRNM(t *testing.T) {
	u := baruUji(t)
	u.a.JenisReas["UJI-INDUK"] = "UJI QS INDUK TRT"
	u.a.SpreadingPolisMap[polisUji] = []models.SpreadingPolis{
		{TreatyType: "UJI-INDUK", SharePercentage: "100", CurrencyID: matauang, Currency: "UJA"},
	}
	m := u.a.Master[masterUji]
	m.Limits[0].Detail[0].SpreadingList = []models.SpreadingMaster{
		{ReinsTypeID: "UJI-RI", ReinsTypeName: "UJI QS RI", Pct: "60"}, {ReinsTypeID: "UJI-OR", ReinsTypeName: "UJI QS OR", Pct: "40"},
	}
	u.a.Master[masterUji] = m
	id := u.buat()
	langkah := func(aksi string, n int, param string, isi map[string]string) map[string]any {
		t.Helper()
		kode, out := u.aksi(id, admin, "", aksi, n, param, isi)
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
		t.Fatalf("spreading bawah dari SpreadingList master: %v", bawah)
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

	// baris dari polis: Treaty Type terkunci
	kode, out := u.aksi(id, admin, "", "SetTreatyNameSpreading", 1, "",
		map[string]string{models.JalurAnak(models.DaftarSpreading, 1, "TreatyType"): "UJI-LAIN"})
	u.wajib(kode, http.StatusConflict, out, "ganti Treaty Type baris dari polis") // aksi tidak tersedia di layar

	// kelengkapan lain seperti alur penuh (Save to issue RNM baru aktif sesudah ada estimasi)
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
	// Delete semua baris: Save to issue RNM ditolak ProteksiData langkah 5
	langkah("DeleteSpreading", 1, "", nil)
	if n, m := len(u.g.Halaman(id).AmbilDaftar(models.DaftarSpreading)), len(u.g.Halaman(id).AmbilDaftar(models.DaftarBreakQS)); n != 0 || m != 0 {
		t.Fatalf("Delete: atas %d bawah %d, mau 0 0", n, m)
	}
	kode, out = u.aksi(id, admin, "", "SaveOutstanding", 0, "", nil)
	if kode == http.StatusOK || !strings.Contains(strings.Join(teks(out["pesan"]), ";"), models.PesanIsiSpreading) {
		t.Fatalf("tanpa spreading mau ditolak ProteksiData 5: HTTP %d %v", kode, out["pesan"])
	}

	// Add: baris kosong, pilih Treaty Type dari dropdown, isi Share
	langkah("AddSpreading", 0, "", nil)
	langkah("SetTreatyNameSpreading", 1, "", map[string]string{models.JalurAnak(models.DaftarSpreading, 1, "TreatyType"): "UJI-INDUK"})
	langkah("CountSpreading", 1, "", map[string]string{models.JalurAnak(models.DaftarSpreading, 1, "SharePercentage"): "100"})
	h = u.g.Halaman(id)
	if a := h.AmbilDaftar(models.DaftarSpreading); len(a) != 1 || a[0]["TreatyName"] != "UJI QS INDUK TRT" || a[0]["SharePercentage"] != "100" {
		t.Fatalf("Add + pilih Treaty Type: %v", a)
	}
	if len(h.AmbilDaftar(models.DaftarBreakQS)) != 2 {
		t.Fatalf("tabel bawah sesudah Add: %v", h.AmbilDaftar(models.DaftarBreakQS))
	}

	// Save to issue RNM lolos
	kode, out = u.aksi(id, admin, "", "SaveOutstanding", 0, "", nil)
	if strings.Contains(strings.Join(teks(out["pesan"]), ";"), models.PesanIsiSpreading) {
		t.Fatalf("spreading terisi masih ditolak ProteksiData 5: %v", out["pesan"])
	}
	u.wajib(kode, http.StatusOK, out, "Save to issue RNM dengan spreading")
}

// Dropdown Treaty Type = spreading polis + pasangan TreatyType -> TreatyName baris yang sudah ada (data lama di luar
// spreading polis tidak tampil sebagai ID mentah); label dari polis tidak ditimpa label baris.
func TestPilihanSpreadingMemuatBarisLama(t *testing.T) {
	u := baruUji(t)
	u.a.JenisReas["UJI-INDUK"] = "UJI QS INDUK TRT"
	u.a.SpreadingPolisMap[polisUji] = []models.SpreadingPolis{
		{TreatyType: "UJI-INDUK", SharePercentage: "100", CurrencyID: matauang, Currency: "UJA"},
	}
	id := u.buat()
	for _, a := range []struct{ aksi, param string }{{"SetValueToClaim", masterUji + "|UJI-TG|UJI-COB"}, {"CheckNoPolicy", polisUji}} {
		kode, out := u.aksi(id, admin, "", a.aksi, 0, a.param, nil)
		u.wajib(kode, http.StatusOK, out, a.aksi)
	}
	h := u.g.Halaman(id)
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{
		{"TreatyType": "UJI-INDUK", "TreatyName": "UJI NAMA BARIS LAIN"},
		{"TreatyType": "UJI-LAMA", "TreatyName": "UJI LAMA TRT", "IsOldData": "Yes"},
	})
	u.g.SetelHalaman(id, h)
	r := httptest.NewRecorder()
	q := httptest.NewRequest(http.MethodGet, handlers.Prefix+"/kasus/"+id+"/pilihan/"+models.SumberSpreading, nil)
	q.Header.Set("X-Pelaku", admin)
	u.srv.ServeHTTP(r, q)
	var opsi []models.Pilihan
	if r.Code != http.StatusOK || json.Unmarshal(r.Body.Bytes(), &opsi) != nil {
		t.Fatalf("pilihan spreading: HTTP %d %s", r.Code, r.Body.String())
	}
	mau := []models.Pilihan{{Nilai: "UJI-INDUK", Label: "UJI QS INDUK TRT"}, {Nilai: "UJI-LAMA", Label: "UJI LAMA TRT"}}
	if len(opsi) != len(mau) || opsi[0].Nilai != mau[0].Nilai || opsi[0].Label != mau[0].Label ||
		opsi[1].Nilai != mau[1].Nilai || opsi[1].Label != mau[1].Label {
		t.Fatalf("pilihan spreading = %+v, mau %+v", opsi, mau)
	}
}
