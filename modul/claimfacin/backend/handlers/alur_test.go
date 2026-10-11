package handlers_test

// Uji seam HTTP: handlers -> services -> gudang tiruan. Alur Register_Flow Claim Fac In (lini Fire) dari pembuatan,
// Choose Polis, Submit Input Register, item + estimasi, Download Claim Face Sheet (nomor klaim K + baris OS STS 0),
// Print PLA, sampai Send to PIC Claim (Choose Surveyor); gerbang wewenang dan kunci diuji lewat endpoint. Fixture
// `UJI-*`.

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/claimfacin/backend/handlers"
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/repository"
	"nusantarare/modul/claimfacin/backend/services"
	"nusantarare/modul/claimfacin/backend/tiruan"
)

const (
	admin    = "UJI-ADMIN"
	lain     = "UJI-LAIN"
	teknik   = "UJI-TEKNIK"
	polisUji = "UJI-RNM-F.001"
	idr      = "10026"
)

// polisFire - dokumen JSON_POLIS fixture lini Fire: satu lokasi, satu property item, satu coverage ber-spreading QS
// (10003) dan fac retro (10015).
const polisFire = `{
 "PercentShare": "25",
 "IsB2B": "",
 "PolicyData": {"PolicyNo": "UJI-RNM-F.001", "StartDateTime": "20260101T000000.000 GMT",
  "EndDateTime": "20261231T000000.000 GMT"},
 "QuotationData": {"BusinessType": "Fire", "BusinessCode": "UJI-B1", "BusinessOldId": "77", "BusinessName": "UJI FIRE",
  "InsuredName": "UJI TERTANGGUNG", "CedingCoName": "UJI CEDING", "SobName": "UJI SOB",
  "CedingCoList": [{"CedingCoName": "UJI CEDING"}]},
 "LocationList": [{"Property": {"ObjectNo": "1", "ObjectName": "UJI GEDUNG", "RiskLocation": {"ASMAddress": "UJI JALAN"},
  "OccupationList": [{"OccupationId": "UJI-OC1", "OccupationName": "UJI OKUPASI"}],
  "PropertyItemList": [{"ItemType": "UJI BANGUNAN", "PropertyItemNo": "1", "Currency": "IDR", "CurrencyID": "10026",
   "TSIObjectItem": "1000000000",
   "CoverageList": [{"OLDID": "UJI-C1", "Coverage": "UJI-COV", "CoverageNote": "UJI FLEXAS", "IndexCoverage": "1",
    "TSINusantaraRe": "250000000", "CoverageBasis": "1",
    "SpreadingList": [{"TreatyType": "10003", "TreatyName": "UJI QS", "SharePercentage": "60", "TSISpreaded": "150000000"},
     {"TreatyType": "10015", "TreatyName": "UJI FAC RETRO", "SharePercentage": "40", "TSISpreaded": "100000000"}]}]}]}}],
 "FacRetroList": [{"ReinsurerID": "UJI-R1", "ReinsurerName": "UJI REAS", "PctShareAllObj": "100"}]
}`

func acuanUji() *tiruan.Acuan {
	a := tiruan.AcuanBaru()
	a.Polis = []models.BarisPolisCari{{PolicyNo: polisUji, CustomerName: "UJI TERTANGGUNG", CedingCoName: "UJI CEDING",
		Prodke: "1", BusinessName: "UJI FIRE"}}
	a.Dokumen[tiruan.KunciPolis(polisUji, "1")] = []byte(polisFire)
	a.NamaMU[idr] = "IDR"
	a.Kurs[idr] = "1"
	a.JenisReas["10003"] = "UJI QS"
	a.JenisReas["10015"] = "UJI FAC RETRO"
	a.Sebab = []repository.BarisSebab{{ID: "UJI-S1", Description: "UJI KEBAKARAN"}}
	a.Email[teknik] = "UJI-EMAIL"
	return a
}

type uji struct {
	t   *testing.T
	srv http.Handler
	g   *tiruan.Gudang
	a   *tiruan.Acuan
	b   *tiruan.Berkas
}

func baruUji(t *testing.T) *uji {
	g, a := tiruan.Baru(), acuanUji()
	for _, k := range []string{"ADU", "CloseClaim", "DLA", "Invoice", "LOD", "SPGR", "Salvage"} {
		g.KategoriDok[k] = "UJI " + k
	}
	a.Tangga = g.TanggaKomite
	jam := func() time.Time { return time.Date(2026, 3, 5, 10, 0, 0, 0, models.Jakarta) }
	b := g.BerkasBaru()
	l := services.Baru(g, a, jam, false).DenganPenyimpanan(b)
	return &uji{t: t, srv: handlers.Router(l, true), g: g, a: a, b: b}
}

func (u *uji) minta(metode, jalur, pelaku, peran string, badan any) (int, map[string]any) {
	u.t.Helper()
	var buf bytes.Buffer
	if badan != nil {
		if err := json.NewEncoder(&buf).Encode(badan); err != nil {
			u.t.Fatal(err)
		}
	}
	r := httptest.NewRequest(metode, jalur, &buf)
	r.Header.Set("X-Pelaku", pelaku)
	if peran != "" {
		r.Header.Set("X-Peran", peran)
	}
	w := httptest.NewRecorder()
	u.srv.ServeHTTP(w, r)
	var out map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &out)
	return w.Code, out
}

func (u *uji) aksiP(id, pelaku, peran string, p services.PermintaanAksi) (int, map[string]any) {
	u.t.Helper()
	return u.minta(http.MethodPost, handlers.Prefix+"/kasus/"+id+"/aksi", pelaku, peran, p)
}

func (u *uji) wajib(kode, mau int, out map[string]any, apa string) {
	u.t.Helper()
	if kode != mau {
		u.t.Fatalf("%s: HTTP %d, mau %d: %v", apa, kode, mau, out)
	}
}

func (u *uji) buat() string {
	u.t.Helper()
	kode, out := u.minta(http.MethodPost, handlers.Prefix+"/kasus", admin, "", nil)
	u.wajib(kode, http.StatusCreated, out, "buat kasus")
	return out["kasus"].(map[string]any)["id"].(string)
}

// sampaiEstimasi - kasus baru, polis dipilih, data register lengkap, satu objek dipilih, Submit -> Input Estimasi.
func (u *uji) sampaiEstimasi() string {
	u.t.Helper()
	id := u.buat()
	kode, out := u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "CopyNB", Konteks: services.ModalPilihPolis,
		Param: polisUji + "|1"})
	u.wajib(kode, http.StatusOK, out, "pilih polis")
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "SetInputParam", Konteks: services.ModalPilihPolis})
	u.wajib(kode, http.StatusOK, out, "submit pilih polis")
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "GetNameCauseofLoss", Param: "UJI-S1"})
	u.wajib(kode, http.StatusOK, out, "pilih cause of loss")
	isian := map[string]string{
		models.CD + "DateOfLoss": "2026-03-01", models.CD + "ReportDate": "2026-03-02",
		models.CD + "DateReceived": "2026-03-02", models.CD + "ReporterName": "UJI PELAPOR",
		models.CD + "ReporterTelp": "0800", models.CD + "Currency": idr, models.CD + "Location": "UJI LOKASI",
		models.CD + "Country": "INDONESIA",
	}
	isian[models.JalurBaris(models.DaftarCalon, 1)+"."+models.PropDipilih] = "true"
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "CheckListEstimasi", Indeks: 1, Masukan: isian})
	u.wajib(kode, http.StatusOK, out, "pilih objek")
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "Submit", Masukan: isian})
	u.wajib(kode, http.StatusOK, out, "submit register")
	if k := u.g.Kasus[id]; k.Tahap != models.TahapEstimasi {
		u.t.Fatalf("sesudah Submit register tahap %q", k.Tahap)
	}
	return id
}

func TestBuatKasusMasukWorklistPembuat(t *testing.T) {
	u := baruUji(t)
	id := u.buat()
	if !strings.HasPrefix(id, models.AwalanKlaim) {
		t.Fatalf("ID %q tanpa awalan CLM-", id)
	}
	if k := u.g.Kasus[id]; k.Tahap != models.TahapRegister || k.PembuatID != admin {
		t.Fatalf("kasus baru %+v", k)
	}
	kode, out := u.minta(http.MethodGet, handlers.Prefix+"/kasus/"+id, lain, "", nil)
	u.wajib(kode, http.StatusOK, out, "buka oleh pelaku lain")
	if out["bolehKerja"] == true {
		t.Fatal("pelaku lain boleh kerja di worklist pembuat")
	}
	kode, out = u.aksiP(id, lain, "", services.PermintaanAksi{Aksi: "Simpan"})
	u.wajib(kode, http.StatusForbidden, out, "aksi pelaku lain")
	kode, out = u.minta(http.MethodGet, handlers.Prefix+"/kasus?daftar=saya", admin, "", nil)
	u.wajib(kode, http.StatusOK, out, "daftar saya")
}

func TestAlurRegisterSampaiSurveyor(t *testing.T) {
	u := baruUji(t)
	id := u.sampaiEstimasi()
	h := u.g.Halaman(id)
	if n := len(h.AmbilDaftar(models.DaftarObjek)); n != 1 {
		t.Fatalf("objek tersimpan %d", n)
	}
	if v := h.Ambil(models.JalurNoPolis); v != polisUji {
		t.Fatalf("polis tersimpan %q", v)
	}
	est := models.KunciPanel(models.PanelObjekEst, models.DaftarObjek, 1)
	kode, out := u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "SetObjectItem", Konteks: est, Indeks: 1, Param: "1"})
	u.wajib(kode, http.StatusOK, out, "pilih property item")
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "PilihCoverageFire", Konteks: est, Indeks: 1,
		Param: "UJI-C1"})
	u.wajib(kode, http.StatusOK, out, "pilih coverage")
	h = u.g.Halaman(id)
	if n := len(h.AmbilDaftar(models.DaftarDiItem(1, 1, models.AnakSpreadPolis))); n != 2 {
		t.Fatalf("spreading polis item %d baris", n)
	}
	item := models.KunciPanel(models.PanelItemEst, models.DaftarItem(1), 1)
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "ValidateInputEstimate", Konteks: item})
	u.wajib(kode, http.StatusOK, out, "tambah estimasi")
	baris := models.JalurBaris(models.DaftarDiItem(1, 1, models.AnakEstimasi), 1)
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "CheckEstimateValue", Konteks: item, Indeks: 1,
		Masukan: map[string]string{baris + ".CurrencyID": idr, baris + ".GrossEstimationPct": "40000000"}})
	u.wajib(kode, http.StatusOK, out, "isi estimasi")
	h = u.g.Halaman(id)
	if v := models.AmbilJalur(h, baris+".EstimationValue"); v != "10000000" {
		t.Fatalf("EstimationValue %q (40.000.000 x 25%%)", v)
	}
	if v := h.AmbilDaftar(models.DaftarObjek)[0]["CFS"]; v != "1" {
		t.Fatalf("objek CFS %q sesudah estimasi", v)
	}
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "CLaimFaceSheet", Indeks: 1})
	u.wajib(kode, http.StatusOK, out, "download claim face sheet")
	h = u.g.Halaman(id)
	if no := h.Ambil(models.CD + "NoClaim"); no != "UJI-K77.03.2026.00001" {
		t.Fatalf("nomor klaim %q", no)
	}
	if len(u.g.OS) != 1 || u.g.OS[0].StsReject != models.StsOSOutstanding || u.g.OS[0].StsDLA != models.StsDLARetro ||
		u.g.OS[0].CaseID != id {
		t.Fatalf("baris OS %+v", u.g.OS)
	}
	if !strings.Contains(u.g.OS[0].DataJSON, `"Value":"10000000"`) {
		t.Fatalf("DATA_JSON OS %q", u.g.OS[0].DataJSON)
	}
	if _, ada := u.g.JSONKlaim[id]; !ada {
		t.Fatal("JSON_KLAIM tidak ditulis")
	}
	ob := h.AmbilDaftar(models.DaftarObjek)[0]
	if ob["IsFacretro"] != "1" || ob["PrintFaceClaim"] != "1" || ob["CFS"] != "0" || h.Ambil(models.JalurIsCFS) != "1" {
		t.Fatalf("penanda CFS objek %v IsCFS %q", ob, h.Ambil(models.JalurIsCFS))
	}
	// CFS kedua tertutup (PrintFaceClaim 1)
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "CLaimFaceSheet", Indeks: 1})
	u.wajib(kode, http.StatusConflict, out, "CFS kedua")
	// Send to PIC sebelum PLA: objek retro tanpa PLA -> pesan
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "SetDisable"})
	u.wajib(kode, http.StatusConflict, out, "send to PIC tertutup (IsPicTransfer)")
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "BukaPLA", Indeks: 1})
	u.wajib(kode, http.StatusOK, out, "buka PLA")
	if out["bukaModal"] != "pla:1" {
		t.Fatalf("modal %v", out["bukaModal"])
	}
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "GeneratePLA", Konteks: "pla:1"})
	u.wajib(kode, http.StatusUnprocessableEntity, out, "PLA tanpa remarks")
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "GeneratePLA", Konteks: "pla:1",
		Masukan: map[string]string{models.JalurObjek(1) + ".RemarksPLA": "UJI CATATAN"}})
	u.wajib(kode, http.StatusOK, out, "generate PLA")
	if v := u.g.Halaman(id).AmbilDaftar(models.DaftarObjek)[0]["PlaStatus"]; v != "1" {
		t.Fatalf("PlaStatus %q", v)
	}
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "SetDisable"})
	u.wajib(kode, http.StatusOK, out, "send to PIC")
	if k := u.g.Kasus[id]; k.Tahap != models.TahapSurveyor || k.Posisi != models.WorkbasketSurveyor {
		t.Fatalf("sesudah Send to PIC %+v", k)
	}
	kode, out = u.aksiP(id, admin, "", services.PermintaanAksi{Aksi: "Simpan"})
	u.wajib(kode, http.StatusForbidden, out, "pembuat bukan pemegang workbasket")
	if len(u.g.Progres) < 2 {
		t.Fatalf("PROGRESSCLAIM %+v", u.g.Progres)
	}
}
