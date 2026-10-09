package handlers_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/tiruan"
)

// Keputusan work owner 09-10-2026 "Occupation cannot empty itu masih ada protek, tapi popup komite masih muncul!":
// Occupation kosong menggagalkan proteksi Send to Committe (Protect.CARI1 = 0) - popup CommitteeTreaty tidak dibuka
// (frontend membukanya hanya bila CARI1 = CARI2 = 1) dan Send Claim to Committee ditolak server dengan pesan
// itu. Occupation dapat diisi di Input Acceptation selama belum ada akseptasi; sesudah diisi, proteksi lolos.
func TestOccupationKosongMenahanKomiteDanBisaDiisiDiAkseptasi(t *testing.T) {
	u, _ := baruUjiLampiran(t)
	id := sampaiAdjustment(u)
	h := u.g.Halaman(id)
	h.Setel(models.CD+"Occupation", "")
	u.g.SetelHalaman(id, h)
	u.a.Roster = []tiruan.AnggotaRoster{
		{AnggotaKomite: models.AnggotaKomite{ID: "1", OperatorID: "UJI-K1", Jabatan: "UJI-JABATAN-1", Degree: "1"},
			Batas: "-1", Sts: models.STSKlaimProp},
	}
	isian := map[string]string{models.CD + "Payable": "2",
		models.JalurAnak(models.DaftarAdjustment, 1, "DataCommitteeTreaty.CircumCauseOfLoss"): "UJI",
		models.JalurAnak(models.DaftarAdjustment, 1, "DataCommitteeTreaty.Remarks"):           "UJI",
		models.JalurAnak(models.DaftarAdjustment, 1, "NameOfBank"):                            "UJI BANK",
		models.JalurAnak(models.DaftarAdjustment, 1, "NoAccount"):                             "UJI-REK"}
	langkah := func(aksi string, m map[string]string) map[string]any {
		t.Helper()
		kode, out := u.aksi(id, teknik, models.WorkbasketAcceptation, aksi, 1, "", m)
		u.wajib(kode, http.StatusOK, out, aksi)
		return out
	}
	langkah("SetNameCurrency", map[string]string{models.JalurAnak(models.DaftarAdjustment, 1, "CurrencyID"): matauang})
	langkah("CountGrossAdjTreaty", map[string]string{models.JalurAnak(models.DaftarAdjustment, 1, "Type"): "1",
		models.JalurAnak(models.DaftarAdjustment, 1, "GrossAdjustment"): "100"})
	for _, k := range []string{"LOD", "DLA", "SPGR"} {
		kode, out := u.unggah(id, teknik, models.WorkbasketAcceptation, k, map[string]string{"UJI-" + k + ".pdf": "%PDF"})
		u.wajib(kode, http.StatusOK, out, "unggah "+k)
	}
	adaKirim := func(out map[string]any) bool {
		m, _ := json.Marshal(out["modal"].(map[string]any)["komite:1"])
		return strings.Contains(string(m), "SendClaimToCommittee")
	}

	protek := func(out map[string]any) any {
		return out["halaman"].(map[string]any)["nilai"].(map[string]any)["Protect.CARI1"]
	}
	out := langkah("BukaKomite", isian)
	if !strings.Contains(strings.Join(teks(out["pesan"]), ";"), models.PesanOccupation) {
		t.Fatalf("Occupation kosong harus berpesan %q: %v", models.PesanOccupation, out["pesan"])
	}
	if protek(out) != "0" {
		t.Fatalf("Occupation kosong: Protect.CARI1 = %v, mau 0 (popup komite tidak dibuka)", protek(out))
	}
	if adaKirim(out) {
		t.Fatal("Occupation kosong: tombol Send Claim to Committee tidak boleh tampil")
	}
	kode, out := u.aksi(id, teknik, models.WorkbasketAcceptation, "AddKomiteTreatyChild", 1, "", isian)
	u.wajib(kode, http.StatusUnprocessableEntity, out, "Send Claim to Committee dengan Occupation kosong")
	if !strings.Contains(strings.Join(teks(out["pesan"]), ";"), models.PesanOccupation) {
		t.Fatalf("penolakan harus menyebut %q: %v", models.PesanOccupation, out["pesan"])
	}

	langkah("Simpan", map[string]string{models.CD + "Occupation": "uji-okupasi"})
	if v := u.g.Halaman(id).Ambil(models.CD + "Occupation"); v != "uji-okupasi" {
		t.Fatalf("Occupation diisi di Input Acceptation tidak tersimpan: %q", v)
	}
	out = langkah("BukaKomite", isian)
	if strings.Contains(strings.Join(teks(out["pesan"]), ";"), models.PesanOccupation) || protek(out) != "1" ||
		!adaKirim(out) {
		t.Fatalf("sesudah Occupation diisi: pesan %v, CARI1 %v, tombol kirim tampil = %v", out["pesan"], protek(out),
			adaKirim(out))
	}
}
