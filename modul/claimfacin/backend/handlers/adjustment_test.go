package handlers_test

// Uji seam HTTP layar Choose Surveyor / Input Adjustment: tambah adjustment (CountTotalEstimasi_Act), Payment Type,
// Choose Currency, Gross Adjustment, deductible, daftar komite calon, Acceptation sesudah komite (fixture tahap 2),
// Print DLA fac retro (nomor P + UPDATE OS), Back ke Input Estimasi. Fixture `UJI-*`.

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/claimfacin/backend/handlers"
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/services"
)

const perTeknik = models.WorkbasketSurveyor

// sampaiSurveyor - kasus sampai Choose Surveyor (langkah TestAlurRegisterSampaiSurveyor tanpa asersi antara).
func (u *uji) sampaiSurveyor() string {
	u.t.Helper()
	id := u.sampaiEstimasi()
	est := models.KunciPanel(models.PanelObjekEst, models.DaftarObjek, 1)
	item := models.KunciPanel(models.PanelItemEst, models.DaftarItem(1), 1)
	baris := models.JalurBaris(models.DaftarDiItem(1, 1, models.AnakEstimasi), 1)
	for _, p := range []services.PermintaanAksi{
		{Aksi: "SetObjectItem", Konteks: est, Indeks: 1, Param: "1"},
		{Aksi: "PilihCoverageFire", Konteks: est, Indeks: 1, Param: "UJI-C1"},
		{Aksi: "ValidateInputEstimate", Konteks: item},
		{Aksi: "SetConvertValueKurs_Estimation", Konteks: item, Indeks: 1,
			Masukan: map[string]string{baris + ".CurrencyID": idr}},
		{Aksi: "CheckEstimateValue", Konteks: item, Indeks: 1,
			Masukan: map[string]string{baris + ".CurrencyID": idr, baris + ".GrossEstimationPct": "40000000"}},
		{Aksi: "CLaimFaceSheet", Indeks: 1},
		{Aksi: "GeneratePLA", Konteks: "pla:1", Masukan: map[string]string{models.JalurObjek(1) + ".RemarksPLA": "UJI"}},
		{Aksi: "SetDisable"},
	} {
		kode, out := u.aksiP(id, admin, "", p)
		u.wajib(kode, http.StatusOK, out, p.Aksi)
	}
	return id
}

func (u *uji) aksiT(id string, p services.PermintaanAksi) (int, map[string]any) {
	u.t.Helper()
	return u.aksiP(id, teknik, perTeknik, p)
}

func panelAdj(a int) string       { return models.KunciPanel(models.PanelAdj, models.DaftarAdj(1, 1), a) }
func jAdj(a int, p string) string { return models.JalurAdj(1, 1, a) + "." + p }

// daftarLayar - baris daftar halaman jawaban layar.
func daftarLayar(out map[string]any, daftar string) []any {
	h, _ := out["halaman"].(map[string]any)
	d, _ := h["daftar"].(map[string]any)
	rows, _ := d[daftar].([]any)
	return rows
}

func pesanMemuat(out map[string]any, teks string) bool {
	rows, _ := out["pesan"].([]any)
	for _, p := range rows {
		if s, _ := p.(string); strings.Contains(s, teks) {
			return true
		}
	}
	return false
}

func TestAlurAdjustmentSampaiDLA(t *testing.T) {
	u := baruUji(t)
	u.a.Adjuster["UJI-ADJ1"] = "UJI ADJUSTER"
	u.a.Roster = []models.AnggotaKomite{
		{ID: "1", OperatorID: "UJI-K1", Jabatan: "UJI SPV", Degree: "1", LimitBottom: "-9999999999999", LimitTop: "57750000"},
		{ID: "2", OperatorID: "UJI-K2", Jabatan: "UJI HEAD", Degree: "2", LimitBottom: "57750000", LimitTop: "189750000"},
	}
	id := u.sampaiSurveyor()
	itemAdj := models.KunciPanel(models.PanelItemAdj, models.DaftarItem(1), 1)

	kode, out := u.minta(http.MethodGet, handlers.Prefix+"/kasus/"+id, teknik, perTeknik, nil)
	u.wajib(kode, http.StatusOK, out, "buka oleh anggota workbasket")
	if out["bolehKerja"] != true {
		t.Fatal("anggota workbasket tidak boleh kerja")
	}
	if p, _ := out["panel"].(map[string]any); p[itemAdj] == nil {
		t.Fatalf("panel item adjustment %q tidak ada", itemAdj)
	}

	// "+" tanpa adjuster / consultant: pesan CountTotalEstimasi_Act 20, baris tidak ditambahkan
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "CountTotalEstimasi", Konteks: itemAdj})
	u.wajib(kode, http.StatusUnprocessableEntity, out, "tambah adjustment tanpa adjuster")
	if !pesanMemuat(out, strings.TrimSpace(models.PesanIsiAdjuster)) {
		t.Fatalf("pesan %v", out["pesan"])
	}
	isi := map[string]string{models.CD + "ConsultantID": "UJI-ADJ1", models.CD + "AppointedADJID": "UJI-ADJ1"}
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "CountTotalEstimasi", Konteks: itemAdj, Masukan: isi})
	u.wajib(kode, http.StatusOK, out, "tambah adjustment")
	h := u.g.Halaman(id)
	adj := h.AmbilDaftar(models.DaftarAdj(1, 1))
	if len(adj) != 1 {
		t.Fatalf("adjustment tersimpan %d", len(adj))
	}
	b := adj[0]
	if b["IsFacRetro"] != "1" || b["DirectToKasir"] != "false" || b["Payable"] != "2" || b["PayableTo"] != "UJI SOB" ||
		b["EstimationValue"] != "10000000" || b["pxCreateOperator"] != teknik {
		t.Fatalf("baris adjustment baru %v", b)
	}
	if n := len(h.AmbilDaftar(models.DaftarDiAdj(1, 1, 1, models.AnakAdjSpread))); n != 2 {
		t.Fatalf("spreading adjustment %d baris", n)
	}

	langkah := []services.PermintaanAksi{
		{Aksi: "SetAdjTypePayment", Konteks: panelAdj(1), Masukan: map[string]string{jAdj(1, "PaymentType"): "1"}},
		{Aksi: "CheckCurrency", Konteks: panelAdj(1), Masukan: map[string]string{jAdj(1, "UploadLOD"): idr}},
		{Aksi: "SetGrossAdjustment", Konteks: panelAdj(1), Masukan: map[string]string{jAdj(1, "GrossAdjustment"): "20000000"}},
	}
	for _, p := range langkah {
		kode, out = u.aksiT(id, p)
		u.wajib(kode, http.StatusOK, out, p.Aksi)
	}
	b = u.g.Halaman(id).AmbilDaftar(models.DaftarAdj(1, 1))[0]
	if b["CurrencyID"] != idr || b["PersenRNM"] != "25" || b["AdjustmentValue"] != "5000000" ||
		b["ValueAdjustment"] != "5000000" || b["GrossValue"] != "5000000" {
		t.Fatalf("gross adjustment %v", b)
	}
	spread := u.g.Halaman(id).AmbilDaftar(models.DaftarDiAdj(1, 1, 1, models.AnakAdjSpread))
	if spread[0]["ClaimSpreaded"] != "3000000" || spread[1]["ClaimSpreaded"] != "2000000" {
		t.Fatalf("spreading adjustment %v", spread)
	}

	// deductible tipe 1 (10% gross)
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "SetNilaiResikoSendiri", Konteks: panelAdj(1),
		Masukan: map[string]string{jAdj(1, "IndividualRiskType"): "1", jAdj(1, "IndividualRiskPercentage"): "10"}})
	u.wajib(kode, http.StatusOK, out, "deductible")
	b = u.g.Halaman(id).AmbilDaftar(models.DaftarAdj(1, 1))[0]
	if b["IndividualRiskValue"] != "2000000" || b["ProposeAdjustmentValue"] != "18000000" || b["AdjustmentValue"] != "4500000" {
		t.Fatalf("deductible %v", b)
	}
	// komite calon (SetListKomite_act): objek retro -> batas ValueAdjustment, hanya LIMIT_BOTTOM <= batas
	komite := daftarLayar(out, models.DaftarDiAdj(1, 1, 1, models.AnakKomiteAdj))
	if len(komite) != 1 || komite[0].(map[string]any)["IDKomite"] != "UJI SPV" {
		t.Fatalf("daftar komite calon %v", komite)
	}

	// adjustment melebihi estimasi: pesan, aksi batal
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "SetGrossAdjustment", Konteks: panelAdj(1),
		Masukan: map[string]string{jAdj(1, "GrossAdjustment"): "80000000"}})
	u.wajib(kode, http.StatusUnprocessableEntity, out, "gross melebihi estimasi")
	if !pesanMemuat(out, models.PesanAdjLebihEstimasi) {
		t.Fatalf("pesan %v", out["pesan"])
	}
	if v := u.g.Halaman(id).AmbilDaftar(models.DaftarAdj(1, 1))[0]["AdjustmentValue"]; v != "4500000" {
		t.Fatalf("aksi batal masih menulis %q", v)
	}

	// Acceptation tertutup sebelum komite
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "Acceptation", Konteks: panelAdj(1)})
	u.wajib(kode, http.StatusConflict, out, "acceptation sebelum komite")

	// komite tahap 2 menerima adjustment (fixture)
	h = u.g.Halaman(id)
	b = h.AmbilDaftar(models.DaftarAdj(1, 1))[0]
	b["AcceptanceStatus"], b["AcceptedNo"], b["IsApproved"], b["IsKomite"] = "1", "UJI-AKS.03.2026.00001", "1", "1"
	b["AcceptedDate"] = "2026-03-05 10:00:00"
	u.g.SetelHalaman(id, h)
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "Acceptation", Konteks: panelAdj(1)})
	u.wajib(kode, http.StatusOK, out, "acceptation")
	h = u.g.Halaman(id)
	b, ob := h.AmbilDaftar(models.DaftarAdj(1, 1))[0], h.AmbilDaftar(models.DaftarObjek)[0]
	if b["IsPrintAccept"] != "1" || ob["IsPrintAccept"] != "1" || ob["DLAStatus"] != "0" {
		t.Fatalf("sesudah acceptation adj %v objek %v", b, ob)
	}
	if len(u.g.Efek) != 0 {
		t.Fatalf("efek keluar di luar produksi %v", u.g.Efek)
	}

	// Print DLA fac retro
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "BukaDLA", Indeks: 1})
	u.wajib(kode, http.StatusOK, out, "buka DLA")
	if out["bukaModal"] != "dla:1" {
		t.Fatalf("modal %v", out["bukaModal"])
	}
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "ChooseDla", Konteks: "dla:1"})
	u.wajib(kode, http.StatusOK, out, "submit DLA")
	h = u.g.Halaman(id)
	b, ob = h.AmbilDaftar(models.DaftarAdj(1, 1))[0], h.AmbilDaftar(models.DaftarObjek)[0]
	if ob["NoDLA"] != "UJI-P77.03.2026.00001" || b["DLA_No"] != ob["NoDLA"] || ob["DLAStatus"] != "1" {
		t.Fatalf("DLA objek %v adj %v", ob, b)
	}
	if len(u.g.DLAOS) != 1 || u.g.DLAOS[0] != "UJI-AKS.03.2026.00001=UJI-P77.03.2026.00001" {
		t.Fatalf("UPDATE OS DLA %v", u.g.DLAOS)
	}
	if v := h.AmbilDaftar(models.DaftarItem(1))[0]["TotalEstimasiReas"]; v != "1800000" {
		t.Fatalf("TotalEstimasiReas %q (TSI RNM 250 jt x 40%% / adjustment 4,5 jt x 40%% x 100%%)", v)
	}
	jenis := map[string]bool{}
	for _, l := range u.g.Log {
		jenis[l.JenisService] = true
	}
	if !jenis[models.JenisServiceAkseptasi] || !jenis[models.JenisServiceDLA] {
		t.Fatalf("log layanan %v", u.g.Log)
	}
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "ChooseDla", Konteks: "dla:1"})
	u.wajib(kode, http.StatusConflict, out, "DLA kedua (DLAStatus 1)")

	// adjustment berkomite tidak dapat dihapus (tombol tersembunyi IsKomite 1)
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "DisableSendComite", Konteks: itemAdj, Indeks: 1})
	u.wajib(kode, http.StatusConflict, out, "hapus adjustment berkomite")
}

func TestAdjustmentHapusDanKembaliEstimasi(t *testing.T) {
	u := baruUji(t)
	u.a.Adjuster["UJI-ADJ1"] = "UJI ADJUSTER"
	id := u.sampaiSurveyor()
	itemAdj := models.KunciPanel(models.PanelItemAdj, models.DaftarItem(1), 1)
	isi := map[string]string{models.CD + "ConsultantID": "UJI-ADJ1", models.CD + "AppointedADJID": "UJI-ADJ1"}
	kode, out := u.aksiT(id, services.PermintaanAksi{Aksi: "CountTotalEstimasi", Konteks: itemAdj, Masukan: isi})
	u.wajib(kode, http.StatusOK, out, "tambah adjustment")
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "DisableSendComite", Konteks: itemAdj, Indeks: 1})
	u.wajib(kode, http.StatusOK, out, "hapus adjustment")
	h := u.g.Halaman(id)
	if n := len(h.AmbilDaftar(models.DaftarAdj(1, 1))); n != 0 {
		t.Fatalf("adjustment tersisa %d", n)
	}
	if !strings.Contains(h.AmbilDaftar(models.DaftarKronologi)[len(h.AmbilDaftar(models.DaftarKronologi))-1]["pyNote"],
		models.TeksBatalAdjustment) {
		t.Fatalf("kronologi %v", h.AmbilDaftar(models.DaftarKronologi))
	}
	kode, out = u.aksiT(id, services.PermintaanAksi{Aksi: "BackToEstimasi"})
	u.wajib(kode, http.StatusOK, out, "back")
	if k := u.g.Kasus[id]; k.Tahap != models.TahapEstimasi || k.Posisi != admin {
		t.Fatalf("sesudah Back %+v", k)
	}
	if v := u.g.Halaman(id).Ambil(models.JalurIsAdjustment); v != "" {
		t.Fatalf("IsAdjustment %q sesudah Back", v)
	}
}
