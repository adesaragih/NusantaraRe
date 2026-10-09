package handlers_test

// Uji seam HTTP penyerahan komite TT2 (SendPICProtect_Act -> pop-up Comittee -> CreateKMTNo_Act), Close Claim
// (ValidationAdjustmentKomite + CloseClaim), dan pop-up Reject Claim (TT3 = OQ-CFI-27). Fixture `UJI-*`.

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/services"
)

// adjustmentFinal - kasus Choose Surveyor dengan satu adjustment PT 1 (IDR, gross 20 jt -> RNM 5 jt).
func (u *uji) adjustmentFinal() string {
	u.t.Helper()
	u.a.Adjuster["UJI-ADJ1"] = "UJI ADJUSTER"
	id := u.sampaiSurveyor()
	itemAdj := models.KunciPanel(models.PanelItemAdj, models.DaftarItem(1), 1)
	for _, p := range []services.PermintaanAksi{
		{Aksi: "CountTotalEstimasi", Konteks: itemAdj, Masukan: map[string]string{
			models.CD + "ConsultantID": "UJI-ADJ1", models.CD + "AppointedADJID": "UJI-ADJ1"}},
		{Aksi: "SetAdjTypePayment", Konteks: panelAdj(1), Masukan: map[string]string{jAdj(1, "PaymentType"): "1"}},
		{Aksi: "CheckCurrency", Konteks: panelAdj(1), Masukan: map[string]string{jAdj(1, "UploadLOD"): idr}},
		{Aksi: "SetGrossAdjustment", Konteks: panelAdj(1), Masukan: map[string]string{jAdj(1, "GrossAdjustment"): "20000000"}},
	} {
		kode, out := u.aksiT(id, p)
		u.wajib(kode, http.StatusOK, out, p.Aksi)
	}
	return id
}

func (u *uji) lampirkan(id string, kategori ...string) {
	for n, k := range kategori {
		u.g.Dokumen = append(u.g.Dokumen, models.BarisDokumenKlaim{ID: "UJI-DOK-" + k, IDPega: models.KunciInstans(id),
			NamaFile: "UJI-" + k + ".pdf", MIME: "pdf", Kategori1: k, Tanggal: time.Date(2026, 3, 5, 9, n, 0, 0, models.Jakarta)})
	}
}

func modeLayar(out map[string]any) map[string]string {
	m := map[string]string{}
	mode, _ := out["mode"].(map[string]any)
	for k, v := range mode {
		m[k], _ = v.(string)
	}
	return m
}

func TestKirimKomiteMelahirkanKMT(t *testing.T) {
	u := baruUji(t)
	u.a.Roster = []models.AnggotaKomite{
		{ID: "1", OperatorID: "UJI-K1", Jabatan: "UJI SPV", Degree: "1", LimitBottom: "-9999999999999", LimitTop: "57750000"},
		{ID: "2", OperatorID: "UJI-K2", Jabatan: "UJI HEAD", Degree: "2", LimitBottom: "57750000", LimitTop: "189750000"},
	}
	id := u.adjustmentFinal()
	modal := models.KunciPanel(models.ModalKomite, models.DaftarAdj(1, 1), 1)

	// tanpa lampiran / rekening: pesan SendPICProtect_Act, pop-up tidak dibuka
	kode, out := u.aksiT(id, services.PermintaanAksi{Aksi: "SendPICProtect", Konteks: panelAdj(1)})
	u.wajib(kode, http.StatusOK, out, "send to committe tanpa lampiran")
	if v, _ := out["bukaModal"].(string); v != "" {
		t.Fatalf("pop-up komite dibuka tanpa lampiran: %v", v)
	}
	if pm, _ := out["pesanMedan"].(map[string]any); pm[models.JalurIsError] == nil {
		t.Fatalf("pesan proteksi %v", out["pesanMedan"])
	}
	// penyerahan langsung ditolak server (gerbang dihitung ulang, bukan dari kiriman)
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "KirimKomite", Konteks: modal,
		Mode:    map[string]string{models.JalurProtect1: "1", models.JalurProtect2: "1"},
		Masukan: map[string]string{jAdj(1, "DataCommitteFacin.Remarks"): "UJI"}})
	if kode == http.StatusOK {
		t.Fatalf("penyerahan tanpa lampiran diterima: %v", out)
	}

	// lampiran PT 1 + rekening
	u.lampirkan(id, "LOD", "DLA", "SPGR")
	h := u.g.Halaman(id)
	b := h.AmbilDaftar(models.DaftarAdj(1, 1))[0]
	b["NameOfBank"], b["NoAccount"], b["IDOfBank"] = "UJI BANK", "0001", "UJI-B1"
	u.g.SetelHalaman(id, h)
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "SendPICProtect", Konteks: panelAdj(1)})
	u.wajib(kode, http.StatusOK, out, "send to committe")
	if out["bukaModal"] != modal {
		t.Fatalf("modal %v, mau %q (pesan %v)", out["bukaModal"], modal, out["pesanMedan"])
	}
	m := modeLayar(out)
	// Remarks wajib
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "KirimKomite", Konteks: modal, Mode: m})
	if kode == http.StatusOK {
		t.Fatalf("penyerahan tanpa Remarks diterima: %v", out)
	}
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "KirimKomite", Konteks: modal, Mode: m,
		Masukan: map[string]string{jAdj(1, "DataCommitteFacin.Remarks"): "UJI REMARKS"}})
	u.wajib(kode, http.StatusOK, out, "send claim to committee")
	h = u.g.Halaman(id)
	b, it := h.AmbilDaftar(models.DaftarAdj(1, 1))[0], h.AmbilDaftar(models.DaftarItem(1))[0]
	kmt := b[models.PropKomiteID]
	if !strings.HasPrefix(kmt, models.AwalanKomite) || b["IsKomite"] != "1" || it["IsKomite"] != "0" ||
		b["DataCommitteFacin.Remarks"] != "UJI REMARKS" {
		t.Fatalf("adjustment sesudah penyerahan %v item %v", b, it)
	}
	k := u.g.Komite[kmt]
	if k.KlaimID != id || k.AdjustmentID != b[models.PropID] || k.Posisi != "UJI-K1" || len(k.Anggota) != 1 {
		t.Fatalf("kasus komite %+v", k)
	}
	if kk := u.g.Kasus[kmt]; kk.Tahap != models.TahapKomite {
		t.Fatalf("T_WORK_CLAIM komite %+v", kk)
	}
	ada := false
	for _, s := range u.g.SubProgres {
		if s.IDPega == kmt && s.Posisi2 == models.ProgresMenungguKmt && s.Posisi1 == "Payment - Final" {
			ada = true
		}
	}
	if !ada {
		t.Fatalf("sub-progress Waiting Committee %+v", u.g.SubProgres)
	}
	// grid Committee Accept Status = tangga tersimpan; penyerahan kedua ditolak
	komite := daftarLayar(out, models.DaftarDiAdj(1, 1, 1, models.AnakKomiteAdj))
	if len(komite) != 1 || komite[0].(map[string]any)["IDKomite"] != "UJI SPV" {
		t.Fatalf("grid komite %v", komite)
	}
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "SendPICProtect", Konteks: panelAdj(1)})
	u.wajib(kode, http.StatusConflict, out, "send to committe kedua")
}

func TestCloseClaim(t *testing.T) {
	u := baruUji(t)
	u.a.Roster = []models.AnggotaKomite{{ID: "1", OperatorID: "UJI-K1", Jabatan: "UJI SPV", LimitBottom: "-1"}}
	id := u.adjustmentFinal()
	// adjustment belum diputuskan komite (roster calon menunggu) -> penutupan ditolak
	kode, out := u.aksiT(id, services.PermintaanAksi{Aksi: "PreventRejectClaim"})
	u.wajib(kode, http.StatusOK, out, "buka close claim")
	if out["bukaModal"] != services.ModalTutup {
		t.Fatalf("modal %v", out["bukaModal"])
	}
	tutup := map[string]string{models.JalurTKRemarks: "UJI TUTUP"}
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "CloseClaim", Konteks: services.ModalTutup, Masukan: tutup})
	u.wajib(kode, http.StatusUnprocessableEntity, out, "close claim dengan adjustment di komite")
	if !pesanMemuat(out, models.PesanTutupDiKomite) {
		t.Fatalf("pesan %v", out["pesan"])
	}
	// komite menerima (fixture tahap 2), Acceptation dicetak -> penutupan berjalan
	h := u.g.Halaman(id)
	b := h.AmbilDaftar(models.DaftarAdj(1, 1))[0]
	b["AcceptanceStatus"], b["AcceptedNo"], b["IsApproved"], b["IsKomite"] = "1", "UJI-AKS.03.2026.00002", "1", "1"
	b[models.PropKomiteID] = "KMT-UJI1"
	u.g.SetelHalaman(id, h)
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "Acceptation", Konteks: panelAdj(1)})
	u.wajib(kode, http.StatusOK, out, "acceptation")
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "CloseClaim", Konteks: services.ModalTutup, Masukan: tutup})
	u.wajib(kode, http.StatusOK, out, "close claim")
	if k := u.g.Kasus[id]; !k.Tertutup() {
		t.Fatalf("kasus belum tertutup %+v", k)
	}
	akhir := u.g.OS[len(u.g.OS)-1]
	if akhir.StsReject != models.StsOSTutup || akhir.CaseID != id ||
		!strings.Contains(akhir.DataJSON, `"pxObjClass":"ASM-FW-GCNMFW-Data-osAkseptasi"`) {
		t.Fatalf("baris OS tutup %+v", akhir)
	}
	if v := u.g.Halaman(id).Ambil(models.CD + "Remark_Close"); v != "UJI TUTUP" {
		t.Fatalf("Remark_Close %q", v)
	}
}

func TestRejectClaimTT3MenungguKeputusan(t *testing.T) {
	u := baruUji(t)
	id := u.buat()
	kode, out := u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "BukaRejectClaim"})
	u.wajib(kode, http.StatusOK, out, "buka reject")
	if out["bukaModal"] != services.ModalTolak {
		t.Fatalf("modal %v", out["bukaModal"])
	}
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "SendRejectClaimToKomite2",
		Konteks: services.ModalTolak, Masukan: map[string]string{models.JalurTKRemarks: "UJI"}})
	if kode == http.StatusOK {
		t.Fatalf("reject TT3 diterima padahal OQ-CFI-27: %v", out)
	}
}
