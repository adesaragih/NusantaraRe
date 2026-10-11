package handlers_test

// Uji seam HTTP penyerahan komite TT2 (SendPICProtect_Act -> pop-up Comittee -> CreateKMTNo_Act), Close Claim
// (ValidationAdjustmentKomite + CloseClaim), dan kelahiran kasus komite TT3 Reject Claim / TT4 Close Without Payment
// (KCF-03). Fixture `UJI-*`.

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/repository"
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

func TestRejectClaimTT3MelahirkanKomiteTanpaAdjustment(t *testing.T) {
	u := baruUji(t)
	id := u.buat()
	kode, out := u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "BukaRejectClaim"})
	u.wajib(kode, http.StatusOK, out, "buka reject")
	if out["bukaModal"] != services.ModalTolak {
		t.Fatalf("modal %v", out["bukaModal"])
	}
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "SendRejectClaimToKomite2",
		Konteks: services.ModalTolak, Masukan: map[string]string{models.JalurTKRemarks: "UJI ALASAN",
			models.JalurTKKronologi: "UJI KRONOLOGI", models.JalurTKExtent: "UJI EXTENT",
			models.JalurTKLiability: "UJI LIABILITAS"}})
	u.wajib(kode, http.StatusOK, out, "reject TT3 (KCF-03)")
	// ClaimComiteeReject LS4 tampil selalu: isian pop-up (halaman requestor TempCommiteClaim) tetap terbaca sesudah Yes
	if nl, _ := out["halaman"].(map[string]any)["nilai"].(map[string]any); nl[models.JalurTKRemarks] != "UJI ALASAN" ||
		nl[models.JalurTKKronologi] != "UJI KRONOLOGI" {
		t.Fatalf("isian pop-up hilang sesudah Yes: Remarks %v Chronology %v", nl[models.JalurTKRemarks],
			nl[models.JalurTKKronologi])
	}
	kmt := modeLayar(out)[models.JalurKomiteBaru]
	k, ada := u.g.Komite[kmt]
	// jawaban work owner 10-10-2026 (OQ-KCFI-03): teks pop-up disimpan di kepala kasus komite (migrasi 643)
	if k.Teks != (repository.TeksKomite{Kronologi: "UJI KRONOLOGI", Extent: "UJI EXTENT", Liability: "UJI LIABILITAS"}) {
		t.Fatalf("teks pop-up kasus komite TT3: %+v", k.Teks)
	}
	if !ada || k.Transfer != models.TransferTolak || k.AdjustmentID != "" || k.KlaimID != id || len(k.Anggota) != 1 ||
		k.Anggota[0].OperatorID != models.WorkbasketTutupKomite || k.Anggota[0].Jabatan != models.JabatanTutupKomite {
		t.Fatalf("kasus komite TT3 %q: %+v", kmt, k)
	}
	h := u.g.Halaman(id)
	if h.Ambil(models.CD+"Remark") != "UJI ALASAN" || h.Ambil(models.CD+"Remark_Close") != "UJI ALASAN" {
		t.Fatalf("Remark: %q / %q", h.Ambil(models.CD+"Remark"), h.Ambil(models.CD+"Remark_Close"))
	}
	kr := h.AmbilDaftar(models.DaftarKronologi)
	if len(kr) == 0 || kr[len(kr)-1]["pyNote"] != models.AwalanTolakKomite+kmt {
		t.Fatalf("kronologi: %+v", kr)
	}
	// permintaan kedua selagi menunggu komite ditolak (penjaga ganda)
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "SendRejectClaimToKomite2",
		Konteks: services.ModalTolak, Masukan: map[string]string{models.JalurTKRemarks: "UJI"}})
	if kode != http.StatusUnprocessableEntity || !strings.Contains(fmt.Sprint(out), models.PesanTutupKomiteGanda) {
		t.Fatalf("permintaan ganda: %d %v", kode, out)
	}
}

func TestCloseWithoutPaymentTT4(t *testing.T) {
	u := baruUji(t)
	id := u.adjustmentFinal()
	kode, out := u.aksiT(id, services.PermintaanAksi{Aksi: "PreventRejectClaim"})
	u.wajib(kode, http.StatusOK, out, "buka close claim")
	tutup := map[string]string{models.JalurTKRemarks: "UJI TANPA BAYAR", models.JalurAlokasiSalvage: "true"}
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "SetelAlokasiSalvage", Konteks: services.ModalTutup,
		Mode: map[string]string{models.JalurAlokasiSalvage: "true"}, Masukan: tutup})
	u.wajib(kode, http.StatusOK, out, "close without payment")
	mode := map[string]string{models.JalurAlokasiSalvage: "true"}
	// adjustment ber-AcceptanceStatus kosong menahan TT4 (SendCloseClaimToKomite 5.1.1.1, VERBATIM)
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "SendCloseClaimToKomite", Konteks: services.ModalTutup,
		Mode: mode, Masukan: tutup})
	if kode != http.StatusUnprocessableEntity || !pesanMemuat(out, models.PesanAdjustmentDiKomite(1)) {
		t.Fatalf("TT4 dengan adjustment belum diputus: %d %v", kode, out["pesan"])
	}
	if len(u.g.Komite) != 0 {
		t.Fatalf("kasus komite lahir walau ditolak: %+v", u.g.Komite)
	}
	h := u.g.Halaman(id)
	h.AmbilDaftar(models.DaftarAdj(1, 1))[0]["AcceptanceStatus"] = "2"
	u.g.SetelHalaman(id, h)
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "SendCloseClaimToKomite", Konteks: services.ModalTutup,
		Mode: mode, Masukan: tutup})
	u.wajib(kode, http.StatusOK, out, "close without payment ke komite")
	kmt := modeLayar(out)[models.JalurKomiteBaru]
	if k := u.g.Komite[kmt]; k.Transfer != models.TransferTutup || k.AdjustmentID != "" || k.KlaimID != id {
		t.Fatalf("kasus komite TT4 %q: %+v", kmt, k)
	}
	if k := u.g.Kasus[id]; k.Tertutup() {
		t.Fatal("klaim ditutup sebelum komite memutus")
	}
	kr := u.g.Halaman(id).AmbilDaftar(models.DaftarKronologi)
	if kr[len(kr)-1]["pyNote"] != models.AwalanTutupKomite+kmt {
		t.Fatalf("kronologi: %+v", kr[len(kr)-1])
	}
}

func TestKasirAcceptationMenungguStatusKonversi(t *testing.T) {
	// HitServiceToKasir_Act 3: transisi PASCA-langkah `.StatusKonversi=="1"` T=2 F=6 (getStatusKonversi_Act, COUNT
	// reinsurance.trloss_detail_t - hanya produksi); selainnya keluar sebelum IDOfBank / muatan kasir.
	u := baruUjiProduksi(t)
	id := u.adjustmentFinal()
	h := u.g.Halaman(id)
	b := h.AmbilDaftar(models.DaftarAdj(1, 1))[0]
	b["AcceptanceStatus"], b["AcceptedNo"], b["IsApproved"], b["IsKomite"] = "1", "UJI-A77.03.2026.00005", "1", "1"
	b["IsPrintAccept"], b["DirectToKasir"], b["StatusKasir"] = "1", "true", ""
	b[models.PropKomiteID] = "KMT-UJI8"
	u.g.SetelHalaman(id, h)
	kasir := func() int {
		n := 0
		for _, e := range u.g.Efek {
			if strings.HasPrefix(e, services.JenisEfekKasir+":") {
				n++
			}
		}
		return n
	}
	kode, out := u.aksiT(id, services.PermintaanAksi{Aksi: "Acceptation", Konteks: panelAdj(1)})
	u.wajib(kode, http.StatusOK, out, "acceptation tanpa status konversi")
	if kasir() != 0 {
		t.Fatalf("kasir diantre sebelum konversi tercatat: %v", u.g.Efek)
	}
	u.a.Konversi["UJI-A7703202600005"] = "1"
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "Acceptation", Konteks: panelAdj(1)})
	u.wajib(kode, http.StatusOK, out, "acceptation sesudah konversi")
	if kasir() != 1 {
		t.Fatalf("kasir tidak diantre sesudah konversi: %v", u.g.Efek)
	}
}

func TestCloseClaimMenutupKomiteAnak(t *testing.T) {
	// CloseClaim 12 `ASMForceCaseClose` CloseAllSubCases=true: kasus komite TT4 yang masih menunggu ikut ditutup
	// (tanpa itu barisnya tertinggal di daftar kerja komite dan setiap Submit-nya 409).
	u := baruUji(t)
	id := u.adjustmentFinal()
	h := u.g.Halaman(id)
	h.AmbilDaftar(models.DaftarAdj(1, 1))[0]["AcceptanceStatus"] = "2"
	u.g.SetelHalaman(id, h)
	kode, out := u.aksiT(id, services.PermintaanAksi{Aksi: "PreventRejectClaim"})
	u.wajib(kode, http.StatusOK, out, "buka close claim")
	mode := map[string]string{models.JalurAlokasiSalvage: "true"}
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "SendCloseClaimToKomite", Konteks: services.ModalTutup,
		Mode: mode, Masukan: map[string]string{models.JalurTKRemarks: "UJI TT4", models.JalurAlokasiSalvage: "true"}})
	u.wajib(kode, http.StatusOK, out, "TT4")
	kmt := modeLayar(out)[models.JalurKomiteBaru]
	if kmt == "" || u.g.Kasus[kmt].Tertutup() {
		t.Fatalf("kasus komite TT4 %q", kmt)
	}
	// adjustment kemudian diterima + dicetak, klaim ditutup lewat Close Claim biasa
	h = u.g.Halaman(id)
	b := h.AmbilDaftar(models.DaftarAdj(1, 1))[0]
	b["AcceptanceStatus"], b["AcceptedNo"], b["IsApproved"], b["IsKomite"] = "1", "UJI-AKS.03.2026.00006", "1", "1"
	b[models.PropKomiteID] = "KMT-UJI7"
	u.g.SetelHalaman(id, h)
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "Acceptation", Konteks: panelAdj(1)})
	u.wajib(kode, http.StatusOK, out, "acceptation")
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "CloseClaim", Konteks: services.ModalTutup,
		Masukan: map[string]string{models.JalurTKRemarks: "UJI TUTUP"}})
	u.wajib(kode, http.StatusOK, out, "close claim")
	if !u.g.Kasus[id].Tertutup() {
		t.Fatal("klaim belum tertutup")
	}
	if k := u.g.Kasus[kmt]; k.StatusWork != models.StatusSelesai {
		t.Fatalf("kasus komite TT4 %s tidak ikut ditutup: %+v", kmt, k)
	}
}
