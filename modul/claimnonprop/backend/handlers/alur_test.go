package handlers_test

// Uji seam HTTP: handlers -> services -> gudang tiruan. Alur Flow_TreatyIn Claim Non Prop dari pembuatan sampai kasus
// komite KMTNP- dan Resolved-Completed, dengan fixture `UJI-*`; gerbang wewenang dan kunci diuji lewat endpoint.

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/claimnonprop/backend/handlers"
	"nusantarare/modul/claimnonprop/backend/models"
	"nusantarare/modul/claimnonprop/backend/services"
	"nusantarare/modul/claimnonprop/backend/tiruan"
)

const (
	admin     = "UJI-ADMIN"
	lain      = "UJI-LAIN"
	teknik    = "UJI-TEKNIK"
	masterUji = "UJI-M000001"
	polisUji  = "UJI-RNM-QR.T1"
	bisnisUji = "UJI-PROPERTY"
	idr       = "10026"
)

func acuanUji() *tiruan.Acuan {
	a := tiruan.AcuanBaru()
	grup := []string{"UJI-GRUP"}
	layer := func(n, limit, ded string) models.LimitXOL {
		return models.LimitXOL{Layer: n, LayerType: "layer", LayerPart: n, LayerPartType: "layer", Currency: "IDR",
			Limit: limit, Limit2: "0", Deductible: ded, Deductible2: "0", ReinstatementPct: "100", IsCombineMDP: "false",
			NoRIPCalculation: "false", TreatyGroups: grup, MDPList: []models.NilaiMataUang{{Currency: "IDR", Value: "1000"}}}
	}
	a.Master[masterUji] = models.MasterTreaty{ID: masterUji, ProportionType: models.ProporsiMaster, Ceding: "UJI-CEDING",
		CedingID: "UJI-CED", LeadingReinsSource: "UJI-SOB", LeadingReinsSourceID: "UJI-SOBID", AccountingModeNonProp: "loss",
		Commencement: "2026-01-01", Termination: "2026-12-31", TreatyYear: "2026", RNMShare: "30",
		Limits:       []models.LimitXOL{layer("1", "3100000000", "2400000000"), layer("2", "6500000000", "5500000000")},
		CurrencyList: []models.KursMaster{{Currency: "IDR", Conversion: "1"}},
		Share: []models.ShareXOL{{SpreadingTypeIDXOL: "UJI-TRT", SpreadingTypeXOL: "UJI TRT", SpreadingTotalPctXOL: "100",
			SpreadingListXOL: []models.SpreadingMaster{{ReinsTypeID: "10004", ReinsTypeName: "QS (R/I)", Pct: "60"},
				{ReinsTypeID: "10028", ReinsTypeName: "QS (OR)", Pct: "40"}}}}}
	a.BarisMaster = []models.BarisMaster{{TreatyID: masterUji, TreatyContractName: "UJI-KONTRAK",
		ProportionType: models.ProporsiMaster, ClassOfBusinessID: "UJI-COB", ClassOfBusiness: bisnisUji, TreatyGroup: "UJI-GRUP",
		TreatyYear: "2026"}}
	a.Polis[masterUji] = []models.BarisPolis{{PolicyNo: polisUji, TreatyGroup: "UJI-COBPOLIS"}}
	a.PolisMaster[masterUji] = true
	a.NamaMU[idr] = "IDR"
	a.OldID[bisnisUji] = "22"
	a.JenisXOL["XL 1ST LAYER"] = models.BarisReinsType{ID: "UJI-XL1", Note: "XL 1ST LAYER"}
	a.JenisXOL["XL 2ND LAYER"] = models.BarisReinsType{ID: "UJI-XL2", Note: "XL 2ND LAYER"}
	a.Adjuster["UJI-ADJ"] = "UJI-ADJUSTER"
	a.Roster = []models.AnggotaKomite{
		{ID: "1", OperatorID: "ReasClaimDeptHead", Jabatan: "Claim Dept. Head", Degree: "1"},
		{ID: "2", OperatorID: "ReasClaimTechDivHead", Jabatan: "Technic Div. Head", Degree: "2"},
	}
	a.Rekening = []tiruan.Rekening{{ClientID: "UJI-CED", RekeningBank: models.RekeningBank{ClientName: "UJI-CEDING",
		NameOfBank: "UJI-BANK", BranchOfBank: "UJI-CABANG", AccountNo: "123", IDOfBank: "UJI-IDB", CurrencyID: idr,
		Currency: "IDR"}}}
	a.Email[teknik] = "UJI-EMAIL"
	return a
}

type uji struct {
	t   *testing.T
	srv http.Handler
	g   *tiruan.Gudang
	a   *tiruan.Acuan
}

func baruUji(t *testing.T) *uji {
	g, a := tiruan.Baru(), acuanUji()
	jam := func() time.Time { return time.Date(2026, 5, 20, 10, 0, 0, 0, models.Jakarta) }
	return &uji{t: t, srv: handlers.Router(services.Baru(g, a, jam, false), true), g: g, a: a}
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

func (u *uji) aksi(id, pelaku, peran, aksi string, indeks int, param string, masukan map[string]string) (int, map[string]any) {
	u.t.Helper()
	return u.aksiP(id, pelaku, peran, services.PermintaanAksi{Aksi: aksi, Indeks: indeks, Param: param, Masukan: masukan})
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

func TestBuatKasusMasukWorklistPembuat(t *testing.T) {
	u := baruUji(t)
	id := u.buat()
	if !strings.HasPrefix(id, models.AwalanKlaim) {
		t.Fatalf("ID %q tanpa awalan CLMNP-", id)
	}
	if k := u.g.Kasus[id]; k.Tahap != models.TahapOutstanding || k.PembuatID != admin {
		t.Fatalf("kasus baru %+v", k)
	}
	kode, out := u.minta(http.MethodGet, handlers.Prefix+"/kasus/"+id, lain, "", nil)
	u.wajib(kode, http.StatusOK, out, "buka oleh pelaku lain")
	if out["bolehKerja"] == true {
		t.Fatal("pelaku lain boleh kerja di worklist pembuat")
	}
	kode, out = u.aksi(id, lain, "", "AddInterestListCNP", 0, "", nil)
	u.wajib(kode, http.StatusForbidden, out, "aksi pelaku lain")
}

// isiOutstanding - Outstanding Claim sampai klaim 10 M IDR teralokasi XoL (pola kasus DEV CLMNP-3998).
func (u *uji) isiOutstanding(id string) {
	u.t.Helper()
	cek := func(kode int, out map[string]any, apa string) { u.t.Helper(); u.wajib(kode, http.StatusOK, out, apa) }
	kode, out := u.aksi(id, admin, "", "SetValueClaimTNP", 0, masterUji+"|UJI-COB|UJI-GRUP|IN", nil)
	cek(kode, out, "pilih master")
	kode, out = u.aksi(id, admin, "", "CheckNoPolicy", 0, "", map[string]string{models.CD + "PolicyData.PolicyNo": polisUji})
	cek(kode, out, "pilih polis")
	kode, out = u.aksi(id, admin, "", "AddInterestListCNP", 0, "", nil)
	cek(kode, out, "tambah interest")
	ji := func(p string) string { return models.JalurAnak(models.DaftarInterest, 1, p) }
	kode, out = u.aksi(id, admin, "", "SetCurrency:Interest", 1, "", map[string]string{ji("ObjectName"): "UJI-OBJEK",
		ji("CurrencyID"): idr, ji("TSIPerObject"): "10000000000"})
	cek(kode, out, "mata uang interest")
	kode, out = u.aksi(id, admin, "", "AddListClaimNP:Claim", 0, "", nil)
	cek(kode, out, "tambah klaim")
	kode, out = u.aksi(id, admin, "", "SetCurrency:Claim", 1, "",
		map[string]string{models.JalurAnak(models.DaftarClaimAmount, 1, "CurrencyID"): idr})
	cek(kode, out, "mata uang klaim")
	kode, out = u.aksi(id, admin, "", "AddLossAlocation", 0, "", nil)
	cek(kode, out, "tambah loss allocation")
	jl := func(p string) string { return models.JalurAnak(models.DaftarLossAlloc, 1, p) }
	kode, out = u.aksi(id, admin, "", "CountLossAllocation:CountXOL", 1, "", map[string]string{jl("TreatyName"): "OR",
		jl("ClaimPercentage"): "100", jl("CNPFlagXOL"): "true"})
	cek(kode, out, "to XOL")
	kode, out = u.aksi(id, admin, "", "CountClaimTNP", 1, "", map[string]string{
		models.JalurAnak(models.DaftarClaimAmount, 1, "AltValue"): "1", models.CD + "ShareCeding": "100",
		models.CD + "Province": "UJI-PROVINSI", models.CD + "CNPCircumtances": "UJI", models.CD + "Location": "UJI",
		models.CD + "DateOfLoss": "2026-04-01", models.CD + "ReportDate": "2026-04-14",
		models.CD + "DateReceived": "2026-04-14", models.CD + "ReporterName": "UJI", models.CD + "ReporterTelp": "1",
		models.CD + "PolicyData.StartDateTime": "2026-01-01", models.CD + "PolicyData.EndDateTime": "2026-12-31"})
	cek(kode, out, "hitung klaim")
}

func TestAlurOutstandingSampaiInputAcceptation(t *testing.T) {
	u := baruUji(t)
	id := u.buat()
	u.isiOutstanding(id)
	h := u.g.Halaman(id)
	xol := h.AmbilDaftar(models.DaftarXOL)
	if len(xol) != 3 || !models.SamaAngka(xol[1]["ClaimSpreaded"], "930000000") ||
		!models.SamaAngka(xol[2]["TotalClaim"], "4500000000") {
		t.Fatalf("SpreadingRisk tersimpan = %v", xol)
	}
	if len(h.AmbilDaftar(models.DaftarBreakQS)) != 2 {
		t.Fatalf("Break QS = %v", h.AmbilDaftar(models.DaftarBreakQS))
	}
	// Submit tertutup sebelum Save to issue RNM dan CFS.
	kode, out := u.aksi(id, admin, "", "Submit", 0, "", nil)
	u.wajib(kode, http.StatusConflict, out, "submit dini")

	kode, out = u.aksi(id, admin, "", "SaveDataToOSAksep", 0, "", nil)
	u.wajib(kode, http.StatusOK, out, "save to issue RNM")
	h = u.g.Halaman(id)
	if no := h.Ambil(models.CD + "NoClaim"); !strings.HasPrefix(no, tiruan.KodeProduksiUji+"K22.05.2026.TX") {
		t.Fatalf("nomor klaim = %q", no)
	}
	if h.Ambil("IsOutstanding") != "1" || h.AmbilDaftar(models.DaftarClaimAmount)[0]["CNPFlagOuts"] != "1" {
		t.Fatal("IsOutstanding / baris terkunci tidak tersimpan")
	}
	if len(u.g.OS) != 2 {
		t.Fatalf("baris OS = %d, mau 2 (XL 1ST + XL 2ND)", len(u.g.OS))
	}
	for _, b := range u.g.OS {
		if b.CaseID != id || b.StsReject != "0" || !strings.Contains(b.DataJSON, "\"TypeLoss\":\"XL ") {
			t.Fatalf("baris OS = %+v", b)
		}
	}
	if _, ada := u.g.JSONKlaim[id]; !ada {
		t.Fatal("JSON_KLAIM tidak ditulis")
	}
	if len(u.g.Efek) != 0 {
		t.Fatalf("efek luar di non-produksi: %v", u.g.Efek)
	}
	// Save to issue RNM kedua tertutup (IsOutstanding = 1).
	kode, out = u.aksi(id, admin, "", "SaveDataToOSAksep", 0, "", nil)
	u.wajib(kode, http.StatusConflict, out, "save to issue RNM kedua")

	kode, out = u.aksi(id, admin, "", "GenerateCFS", 0, "", nil)
	u.wajib(kode, http.StatusOK, out, "print CFS")
	kode, out = u.aksi(id, admin, "", "Submit", 0, "", nil)
	u.wajib(kode, http.StatusOK, out, "submit")
	k := u.g.Kasus[id]
	if k.Tahap != models.TahapAcceptation || k.Posisi != models.WorkbasketAcceptation {
		t.Fatalf("sesudah submit = %+v", k)
	}
	if u.g.Halaman(id).Ambil("CNPStatusCase") != models.StatusInputAcceptance {
		t.Fatal("CNPStatusCase tidak diisi post-activity")
	}
	// Assignment1 milik workbasket ReasKlaimTeknik.
	kode, out = u.aksi(id, admin, "", "SaveToOS", 0, "", nil)
	u.wajib(kode, http.StatusForbidden, out, "pembuat di Input Acceptation")
}

// kirimKeTeknik menyiapkan kasus di Input Acceptation.
func (u *uji) kirimKeTeknik() string {
	u.t.Helper()
	id := u.buat()
	u.isiOutstanding(id)
	for _, a := range []string{"SaveDataToOSAksep", "GenerateCFS", "Submit"} {
		kode, out := u.aksi(id, admin, "", a, 0, "", nil)
		u.wajib(kode, http.StatusOK, out, a)
	}
	return id
}

func TestAkseptasiSampaiKasusKomite(t *testing.T) {
	u := baruUji(t)
	id := u.kirimKeTeknik()
	peran := models.WorkbasketAcceptation
	kode, out := u.aksi(id, teknik, peran, "AddAkseptasiCNP", 0, "", nil)
	u.wajib(kode, http.StatusConflict, out, "add sebelum Save To OS")
	kode, out = u.aksi(id, teknik, peran, "SaveToOS", 0, "", nil)
	u.wajib(kode, http.StatusOK, out, "save to OS")
	if u.g.Halaman(id).Ambil("IsSaveToOs") != "1" {
		t.Fatal("IsSaveToOs tidak tersimpan")
	}
	// Tanpa Adjuster / Consultant akseptasi ditolak dengan pesan XML.
	kode, out = u.aksi(id, teknik, peran, "AddAkseptasiCNP", 0, "", nil)
	u.wajib(kode, http.StatusUnprocessableEntity, out, "add tanpa adjuster")
	h := u.g.Halaman(id)
	h.Setel(models.CD+"AppointedADJID", "UJI-ADJ")
	h.Setel(models.CD+"ConsultantID", "UJI-ADJ")
	h.Setel(models.CD+"QuotationData.BusinessOldId", "22")
	u.g.SetelHalaman(id, h)
	kode, out = u.aksi(id, teknik, peran, "AddAkseptasiCNP", 0, "", nil)
	u.wajib(kode, http.StatusOK, out, "add akseptasi")
	h = u.g.Halaman(id)
	if n := len(h.AmbilDaftar(models.DaftarAdjustment)); n != 1 {
		t.Fatalf("akseptasi = %d", n)
	}
	if len(h.AmbilDaftar(models.JalurAdj(1, models.AnakXOL))) != 3 || len(h.AmbilDaftar(models.JalurAdj(1, models.AnakClaimAccept))) != 1 {
		t.Fatal("salinan XOL / Claim Acceptation akseptasi tidak tersimpan")
	}
	// Payable 1 (ceding) -> rekening ceding.
	kode, out = u.aksiP(id, teknik, peran, services.PermintaanAksi{Aksi: "SetPayableTreatyNP:Acc", Indeks: 1,
		Masukan: map[string]string{models.CD + "Payable": "1", models.JalurAdj(1, "PaymentType"): "1",
			models.CD + "Occupation": "UJI"}})
	u.wajib(kode, http.StatusOK, out, "payable")
	h = u.g.Halaman(id)
	if b := h.AmbilDaftar(models.DaftarAdjustment)[0]; b["NoAccount"] != "123" || b["IDOfBank"] != "UJI-IDB" {
		t.Fatalf("rekening akseptasi = %v", b)
	}
	kode, out = u.aksiP(id, teknik, peran, services.PermintaanAksi{Aksi: "BukaKomite", Indeks: 1})
	u.wajib(kode, http.StatusOK, out, "buka komite")
	if out["bukaModal"] != "komite:1" {
		t.Fatalf("modal = %v", out["bukaModal"])
	}
	mode, _ := out["mode"].(map[string]any)
	modeTeks := map[string]string{}
	for k, v := range mode {
		modeTeks[k] = v.(string)
	}
	// Perbaikan OQ-CNP-05 butir 4: Gross Value layer != Spreading In -> ditolak SEBELUM kasus komite dibuat, nol tulisan.
	h = u.g.Halaman(id)
	asli := h.AmbilDaftar(models.JalurAdj(1, models.AnakSpreadIn))[0]["ClaimSpreaded"]
	h.AmbilDaftar(models.JalurAdj(1, models.AnakSpreadIn))[0]["ClaimSpreaded"] = "1"
	u.g.SetelHalaman(id, h)
	kode, out = u.aksiP(id, teknik, peran, services.PermintaanAksi{Aksi: "CreateChildKomiteCNP", Indeks: 1, Mode: modeTeks,
		Masukan: map[string]string{models.JalurAdj(1, "DataCommitteeTreaty.Remarks"): "UJI"}})
	u.wajib(kode, http.StatusUnprocessableEntity, out, "gross value beda")
	if p := strings.Join(teks(out["pesan"]), ";"); p != models.PesanGrossSpreadingIn {
		t.Fatalf("pesan = %q", p)
	}
	if len(u.g.Komite) != 0 || u.g.Halaman(id).AmbilDaftar(models.DaftarAdjustment)[0]["IsKomite"] == "1" {
		t.Fatal("kasus komite / penanda tertulis walau validasi gagal")
	}
	h = u.g.Halaman(id)
	h.AmbilDaftar(models.JalurAdj(1, models.AnakSpreadIn))[0]["ClaimSpreaded"] = asli
	u.g.SetelHalaman(id, h)
	kode, out = u.aksiP(id, teknik, peran, services.PermintaanAksi{Aksi: "CreateChildKomiteCNP", Indeks: 1, Mode: modeTeks,
		Masukan: map[string]string{models.JalurAdj(1, "DataCommitteeTreaty.Remarks"): "UJI"}})
	u.wajib(kode, http.StatusOK, out, "kirim komite")
	if len(u.g.Komite) != 1 {
		t.Fatalf("kasus komite = %v", u.g.Komite)
	}
	for kid, tangga := range u.g.Komite {
		if !strings.HasPrefix(kid, models.AwalanKomite) {
			t.Fatalf("ID komite %q", kid)
		}
		// RNM Share 30, ValueAdjustment tanpa penulis (0) -> hanya tingkat 1 (OQ-CNP-01 ikut XML).
		if len(tangga) != 1 || tangga[0].OperatorID != "ReasClaimDeptHead" {
			t.Fatalf("tangga = %+v", tangga)
		}
	}
	b := u.g.Halaman(id).AmbilDaftar(models.DaftarAdjustment)[0]
	if b["IsKomite"] != "1" || b["AcceptanceStatus"] != "0" || b[models.PropKomiteID] == "" {
		t.Fatalf("akseptasi sesudah kirim = %v", b)
	}
	// Akseptasi di komite: tulisan ditolak; Close Claim ditolak pesan XML.
	kode, out = u.aksiP(id, teknik, peran, services.PermintaanAksi{Aksi: "SetInterimXOL", Indeks: 1})
	u.wajib(kode, http.StatusConflict, out, "tulis akseptasi di komite")
	kode, out = u.aksi(id, teknik, peran, "CloseClaimTNonProp", 0, "", map[string]string{"Message": "UJI"})
	u.wajib(kode, http.StatusUnprocessableEntity, out, "close saat komite")
}

func TestCloseClaimMenutupKasusDanMenulisOS4(t *testing.T) {
	u := baruUji(t)
	id := u.kirimKeTeknik()
	peran := models.WorkbasketAcceptation
	kode, out := u.aksi(id, teknik, peran, "CloseClaimTNonProp", 0, "", nil)
	u.wajib(kode, http.StatusUnprocessableEntity, out, "close tanpa remarks")
	n := len(u.g.OS)
	kode, out = u.aksi(id, teknik, peran, "CloseClaimTNonProp", 0, "", map[string]string{"Message": "UJI-ALASAN"})
	u.wajib(kode, http.StatusOK, out, "close claim")
	if !u.g.Kasus[id].Tertutup() {
		t.Fatal("kasus tidak tertutup")
	}
	if len(u.g.OS) != n+1 || u.g.OS[n].StsReject != models.StsOSFinal {
		t.Fatalf("baris OS STS 4 = %+v", u.g.OS[n:])
	}
	h := u.g.Halaman(id)
	riw := h.AmbilDaftar(models.DaftarRiwayat)
	if riw[len(riw)-1]["CommentSuggest"] != models.TeksTutupKlaim {
		t.Fatalf("kronologi terakhir = %v", riw[len(riw)-1])
	}
	kode, out = u.aksi(id, teknik, peran, "SaveToOS", 0, "", nil)
	u.wajib(kode, http.StatusConflict, out, "tulis kasus tertutup")
}

func TestTombolOQTidakDapatDijalankan(t *testing.T) {
	u := baruUji(t)
	id := u.kirimKeTeknik()
	peran := models.WorkbasketAcceptation
	kode, out := u.aksi(id, teknik, peran, "BukaCWP", 0, "", nil)
	u.wajib(kode, http.StatusOK, out, "buka CWP")
	if out["bukaModal"] != "cwp" {
		t.Fatalf("modal = %v", out["bukaModal"])
	}
	kode, out = u.aksi(id, teknik, peran, "CwpYes", 0, "", nil)
	if kode != http.StatusBadRequest && kode != http.StatusConflict {
		t.Fatalf("CWP Yes = %d %v", kode, out)
	}
}

func teks(v any) []string {
	var out []string
	if l, ok := v.([]any); ok {
		for _, x := range l {
			out = append(out, x.(string))
		}
	}
	return out
}
