package models

import "testing"

func terbuka(h *Halaman, defs []Unsur, aksi string, indeks int) bool {
	return AksiTerbuka(Evaluasi(h, defs, false), aksi, indeks)
}

// Section OutstandingClaim S67: Save to issue RNM NA[IsOutstanding = 1 || ProtectDOL || ProtectStartDate], Print CFS
// NA[IsOutstanding != '1' || IsCFS == 1], Submit NA[IsCFS != '1' || IsOutstanding != '1'].
func TestTombolOutstanding(t *testing.T) {
	h := HalamanBaru()
	if !terbuka(h, LayarOutstanding(), "SaveDataToOSAksep", 0) {
		t.Error("Save to issue RNM harus terbuka pada kasus baru")
	}
	for _, a := range []string{"GenerateCFS", "Submit"} {
		if terbuka(h, LayarOutstanding(), a, 0) {
			t.Errorf("%s harus tertutup sebelum Save to issue RNM", a)
		}
	}
	h.Setel("IsOutstanding", "1")
	if terbuka(h, LayarOutstanding(), "SaveDataToOSAksep", 0) || !terbuka(h, LayarOutstanding(), "GenerateCFS", 0) {
		t.Error("sesudah issue RNM: Save to issue RNM tertutup, Print CFS terbuka")
	}
	h.Setel("IsCFS", "1")
	if !terbuka(h, LayarOutstanding(), "Submit", 0) || terbuka(h, LayarOutstanding(), "GenerateCFS", 0) {
		t.Error("sesudah CFS: Submit terbuka, Print CFS tertutup")
	}
	// DOL melewati akhir treaty pada mode loss (ProtectDOL) mengunci Save to issue RNM.
	h2 := HalamanBaru()
	h2.Setel(TM+"AccountingModeNonProp", ModeLoss)
	h2.Setel(CD+"EndDateTreaty", "2025-12-31")
	h2.Setel(CD+"DateOfLoss", "2026-01-02")
	if terbuka(h2, LayarOutstanding(), "SaveDataToOSAksep", 0) {
		t.Error("ProtectDOL harus menutup Save to issue RNM")
	}
}

// Tombol tanpa rule ekspor / jalur konfigurasi: tampil, nonaktif, ber-OQ.
func TestTombolOQNonaktif(t *testing.T) {
	ts := Evaluasi(HalamanBaru(), LayarOutstanding(), false)
	for id, oq := range map[string]string{"EditXOLAllocation": OQSandiEditXOL, "ViewMaster": OQViewMaster,
		"AddListClaimNPSpreading": OQGridSpreading} {
		u, ada := cariTata(ts, id)
		if !ada {
			t.Errorf("%s tidak tampil", id)
			continue
		}
		if !u.Nonaktif || u.Catatan != oq {
			t.Errorf("%s nonaktif=%v catatan=%q", id, u.Nonaktif, u.Catatan)
		}
	}
}

func cariTata(ts []Tata, id string) (Tata, bool) {
	for _, t := range ts {
		if t.ID == id {
			return t, true
		}
		if t.Tambah != nil && t.Tambah.ID == id {
			return *t.Tambah, true
		}
		if u, ada := cariTata(t.Anak, id); ada {
			return u, true
		}
	}
	return Tata{}, false
}

// Baris terkunci (CNPFlagOuts = 1, sesudah Save to issue RNM) tidak bisa diubah layar.
func TestBarisTerkunciTidakTerbuka(t *testing.T) {
	h := HalamanBaru()
	h.SetelDaftar(DaftarClaimAmount, []Baris{{"Value": "1"}, {"Value": "2", "CNPFlagOuts": "1"}})
	m := MedanTerbuka(Evaluasi(h, LayarOutstanding(), false))
	if !m[JalurAnak(DaftarClaimAmount, 1, "Value")] || m[JalurAnak(DaftarClaimAmount, 2, "Value")] {
		t.Errorf("Value baris 1 terbuka / baris 2 terkunci: %v %v", m[JalurAnak(DaftarClaimAmount, 1, "Value")],
			m[JalurAnak(DaftarClaimAmount, 2, "Value")])
	}
}

// Section InputAcceptation S14: Add akseptasi NA[IsSaveToOs == 0]; Delete NA[AcceptanceStatus != ”].
func TestTombolAkseptasi(t *testing.T) {
	h := HalamanBaru()
	h.Setel("IsSaveToOs", "0")
	if terbuka(h, LayarAkseptasi(), "AddAkseptasiCNP", 0) {
		t.Error("Add akseptasi tertutup sebelum Save To OS")
	}
	h.Setel("IsSaveToOs", "1")
	if !terbuka(h, LayarAkseptasi(), "AddAkseptasiCNP", 0) {
		t.Error("Add akseptasi terbuka sesudah Save To OS")
	}
	h.SetelDaftar(DaftarAdjustment, []Baris{{"AcceptanceStatus": ""}, {"AcceptanceStatus": "1"}})
	if !terbuka(h, LayarAkseptasi(), "DeleteAkseptasi", 1) || terbuka(h, LayarAkseptasi(), "DeleteAkseptasi", 2) {
		t.Error("Delete hanya untuk akseptasi tanpa status")
	}
}

// AdjustmentDetailNP S40: Send to Committe tampil[TotalKomite != ”] NA[IsKomite = 1 || IsSaveToOs == 0]; Generate DLA
// nonaktif (OQ-CNP-19); Add grid akseptasi nonaktif (OQ-CNP-33).
func TestTombolDetailAkseptasi(t *testing.T) {
	h := HalamanBaru()
	h.Setel("IsSaveToOs", "1")
	h.SetelDaftar(DaftarAdjustment, []Baris{{}})
	if terbuka(h, LayarDetailAkseptasi(1), "BukaKomite", 1) {
		t.Error("Send to Committe tersembunyi tanpa TotalKomite")
	}
	h.AmbilDaftar(DaftarAdjustment)[0]["TotalKomite"] = "1"
	if !terbuka(h, LayarDetailAkseptasi(1), "BukaKomite", 1) {
		t.Error("Send to Committe terbuka")
	}
	h.AmbilDaftar(DaftarAdjustment)[0]["IsKomite"] = "1"
	if terbuka(h, LayarDetailAkseptasi(1), "BukaKomite", 1) {
		t.Error("Send to Committe tertutup saat akseptasi di komite")
	}
	for _, a := range []string{"GenerateDLA", "AddListClaimNPAkseptasi", "AddLossAlocationAkseptasi"} {
		if terbuka(h, LayarDetailAkseptasi(1), a, 1) {
			t.Errorf("%s harus nonaktif", a)
		}
	}
}

// KomiteCLMNP: Send Claim to Committee hanya bila gerbang proteksi lolos, Payable dan Occupation terisi.
func TestGerbangKirimKomite(t *testing.T) {
	h := HalamanBaru()
	h.SetelDaftar(DaftarAdjustment, []Baris{{}})
	h.Setel(CD+"Payable", "1")
	h.Setel(CD+"Occupation", "UJI")
	if terbuka(h, LayarKomite(1), "CreateChildKomiteCNP", 1) {
		t.Error("tanpa CekError.Lolos tertutup")
	}
	h.Setel(JalurLolosKomite, "1")
	if !terbuka(h, LayarKomite(1), "CreateChildKomiteCNP", 1) {
		t.Error("lolos -> terbuka")
	}
	h.Setel(CD+"Occupation", "")
	if terbuka(h, LayarKomite(1), "CreateChildKomiteCNP", 1) {
		t.Error("Occupation kosong -> tertutup")
	}
}

// CloseClaimNP: Yes nonaktif (kasus komite tanpa akseptasi tidak dapat ditulis, OQ-CNP-36).
func TestCWPYesNonaktif(t *testing.T) {
	if terbuka(HalamanBaru(), LayarCWP(), "CwpYes", 0) {
		t.Error("CWP Yes harus nonaktif")
	}
}
