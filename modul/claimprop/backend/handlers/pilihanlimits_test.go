package handlers_test

// Dropdown Treaty Type grid Loss Allocation (`pilihan/limits` = TreatyInMaster.Limits): setiap Treaty Type master
// menjadi pilihan - master QUOTA SHARE + SURPLUS menawarkan keduanya (laporan work owner 10-10-2026 "kenapa ga bisa
// pilih surplus"; master 1002308 tidak ada di JSON lama sehingga daftar dulu kosong, kini dari TREATYINDETAILJOINEDM).

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"nusantarare/modul/claimprop/backend/handlers"
	"nusantarare/modul/claimprop/backend/models"
)

func TestPilihanLimitsMenawarkanSemuaTreatyTypeMaster(t *testing.T) {
	u := baruUji(t)
	m := u.a.Master[masterUji]
	m.Limits = []models.LimitMaster{
		{TreatyType: "QUOTA SHARE", Detail: []models.DetailLimit{{TreatyGroupID: "UJI-TG", RNMShare: "5"}}},
		{TreatyType: "SURPLUS", Detail: []models.DetailLimit{{TreatyGroupID: "UJI-TG", RNMShare: "5"}}},
	}
	u.a.Master[masterUji] = m
	id := u.buat()
	kode, out := u.aksi(id, admin, "", "SetValueToClaim", 0, masterUji+"|UJI-TG|UJI-COB", nil)
	u.wajib(kode, http.StatusOK, out, "SetValueToClaim")

	r := httptest.NewRecorder()
	q := httptest.NewRequest(http.MethodGet, handlers.Prefix+"/kasus/"+id+"/pilihan/"+models.SumberLimits, nil)
	q.Header.Set("X-Pelaku", admin)
	u.srv.ServeHTTP(r, q)
	var opsi []models.Pilihan
	if r.Code != http.StatusOK || json.Unmarshal(r.Body.Bytes(), &opsi) != nil {
		t.Fatalf("pilihan limits: HTTP %d %s", r.Code, r.Body.String())
	}
	var nilai []string
	for _, o := range opsi {
		nilai = append(nilai, o.Nilai)
	}
	if len(nilai) != 2 || nilai[0] != "QUOTA SHARE" || nilai[1] != "SURPLUS" {
		t.Fatalf("pilihan Treaty Type Loss Allocation = %v, harap [QUOTA SHARE SURPLUS]", nilai)
	}
}

// Memilih SURPLUS pada baris Loss Allocation yang SUDAH QUOTA SHARE: nama yang dipilih menentukan ID (dulu
// SetNameTreaty menamai ulang baris dari ID lamanya, sehingga pilihan SURPLUS kembali menjadi QUOTA SHARE -
// laporan work owner 10-10-2026 "masih ga bisa yang surplus").
func TestLossAllocationGantiQuotaShareKeSurplus(t *testing.T) {
	u := baruUji(t)
	u.a.JenisReas["UJI-SURPID"] = "SURPLUS"
	u.a.TipeReas["UJI-SURPID"] = "4"
	id := u.buat()
	langkah := func(aksi string, n int, param string, isi map[string]string) {
		t.Helper()
		kode, out := u.aksi(id, admin, "", aksi, n, param, isi)
		u.wajib(kode, http.StatusOK, out, aksi)
	}
	langkah("SetValueToClaim", 0, masterUji+"|UJI-TG|UJI-COB", nil)
	langkah("CheckNoPolicy", 0, polisUji, nil)
	// kelengkapan sebelum grid Loss Allocation tersedia - urutan TestAlur (alur_test.go)
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
	jalur := models.JalurAnak(models.DaftarLossAlloc, 1, "TreatyName")
	langkah("SetNameTreaty", 1, "", map[string]string{jalur: "QUOTA SHARE"})
	if b := u.g.Halaman(id).AmbilDaftar(models.DaftarLossAlloc)[0]; b["TreatyType"] != jenisQSUji {
		t.Fatalf("QUOTA SHARE: TreatyType %q, harap %q", b["TreatyType"], jenisQSUji)
	}
	langkah("SetNameTreaty", 1, "", map[string]string{jalur: "SURPLUS"})
	b := u.g.Halaman(id).AmbilDaftar(models.DaftarLossAlloc)[0]
	if b["TreatyName"] != "SURPLUS" || b["TreatyType"] != "UJI-SURPID" {
		t.Fatalf("ganti ke SURPLUS: TreatyName %q TreatyType %q, harap SURPLUS / UJI-SURPID", b["TreatyName"], b["TreatyType"])
	}
}
