package handlers_test

// Uji seam 2 (HTTP): handlers -> services -> gudang tiruan. Alur Flow_TreatyIn dari pembuatan sampai Resolved-Completed
// dengan fixture `UJI-*`; gerbang wewenang dan kunci diuji lewat endpoint, bukan lewat layar (AC 57, 58).

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/claimprop/backend/handlers"
	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/repository"
	"nusantarare/modul/claimprop/backend/services"
	"nusantarare/modul/claimprop/backend/tiruan"
)

var (
	_ services.Gudang = (*tiruan.Gudang)(nil)
	_ services.Acuan  = (*tiruan.Acuan)(nil)
)

const (
	admin      = "UJI-ADMIN"
	lain       = "UJI-LAIN"
	teknik     = "UJI-TEKNIK"
	masterUji  = "UJI-M0000001"
	polisUji   = "UJI-RNM-Q.T1"
	matauang   = "UJI-A"
	jenisQSUji = "UJI-QSID"
)

func acuanUji() *tiruan.Acuan {
	a := tiruan.AcuanBaru()
	a.Master[masterUji] = models.MasterTreaty{ID: masterUji, TreatyContractName: "UJI-KONTRAK", ProportionType: "Proportional",
		Ceding: "UJI-CEDING", CedingID: "UJI-CED", LeadingReinsSource: "UJI-SOB", LeadingReinsSourceID: "UJI-SOBID",
		Commencement: "2026-01-01", Termination: "2026-12-31", TreatyYear: "2026", RNMShareP: "10",
		StatusAkseptasi: "Resolve Complete", Limits: []models.LimitMaster{{TreatyType: "QUOTA SHARE",
			Detail: []models.DetailLimit{{TreatyGroupID: "UJI-TG", RNMShare: "10",
				SpreadingList: []models.SpreadingMaster{{ReinsTypeID: "UJI-R1", ReinsTypeName: "UJI-QS", Pct: "100"}}}}}}}
	a.BarisMaster = []models.BarisMaster{{TreatyID: masterUji, ClassOfBusiness: "UJI-COBNAMA", ClassOfBusinessID: "UJI-COB",
		TreatyContractName: "UJI-KONTRAK", ProportionType: "Proportional", TreatyGroup: "UJI-GRUP", TreatyGroupID: "UJI-TG",
		TreatyYear: "2026"}}
	a.Polis = []models.BarisPolis{{PolicyNo: polisUji, NoOffer: models.AwalanMaster(masterUji), TreatyGroup: "UJI-GRUP",
		TreatyYear: "2026", Quarter: "1", Prodke: "0"}}
	a.PolisRealisasi[polisUji] = true
	a.PolisMaster[masterUji] = true
	a.Kurs[matauang] = "2"
	a.NamaMU[matauang] = "UJA"
	a.OldID["UJI-COB"] = "12"
	a.TreatyGroup["UJI-COB|2026"] = "UJI-TG"
	a.Tahun["UJI-TG|*"] = "2026"
	a.JenisReas[jenisQSUji] = "QUOTA SHARE"
	a.TipeReas[jenisQSUji] = "4"
	a.Sebab = append(a.Sebab, repository.BarisSebab{ID: "UJI-S1", Description: "UJI-SEBAB"})
	a.Adjuster["UJI-ADJ"] = "UJI-ADJUSTER"
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

func (u *uji) aksi(id, pelaku, peran, aksi string, indeks int, param string, masukan map[string]string) (int, map[string]any) {
	u.t.Helper()
	return u.minta(http.MethodPost, handlers.Prefix+"/kasus/"+id+"/aksi", pelaku, peran,
		services.PermintaanAksi{Aksi: aksi, Indeks: indeks, Param: param, Masukan: masukan})
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
		t.Fatalf("ID %q tanpa awalan CLMP-", id)
	}
	k := u.g.Kasus[id]
	if k.Tahap != models.TahapOutstanding || k.PembuatID != admin {
		t.Fatalf("kasus baru %+v, mau tahap OutstandingClaim milik pembuat", k)
	}
	kode, out := u.minta(http.MethodGet, handlers.Prefix+"/kasus?daftar=saya", admin, "", nil)
	if kode != http.StatusOK {
		t.Fatalf("daftar saya: %d", kode)
	}
	_ = out
	kode, out = u.minta(http.MethodGet, handlers.Prefix+"/kasus/"+id, lain, "", nil)
	u.wajib(kode, http.StatusOK, out, "buka oleh pelaku lain")
	if out["bolehKerja"] == true {
		t.Fatalf("pelaku lain boleh kerja di worklist pembuat")
	}
}

func TestBukanPemegangDitolakDiLayanan(t *testing.T) {
	u := baruUji(t)
	id := u.buat()
	kode, out := u.aksi(id, lain, "", "SetOutstanding", 0, "", nil)
	u.wajib(kode, http.StatusForbidden, out, "aksi bukan pemegang")
}

func TestMedanTerkunciTidakDitulisLayar(t *testing.T) {
	u := baruUji(t)
	id := u.buat()
	kode, out := u.aksi(id, admin, "", "SetOutstanding", 0, "", map[string]string{
		models.CD + "NoClaim": "UJI-PALSU", models.CD + "ReporterName": "UJI-PELAPOR"})
	u.wajib(kode, http.StatusOK, out, "save")
	h := u.g.Halaman(id)
	if h.Ambil(models.CD+"NoClaim") != "" {
		t.Errorf("Claim No (read-only selalu) tertulis dari layar: %q", h.Ambil(models.CD+"NoClaim"))
	}
	if h.Ambil(models.CD+"ReporterName") != "UJI-PELAPOR" {
		t.Errorf("Reporter Name tidak tersimpan")
	}
}

func TestAlurPenuhSampaiResolved(t *testing.T) {
	u := baruUji(t)
	id := u.buat()
	langkah := func(aksi string, n int, param string, m map[string]string) map[string]any {
		t.Helper()
		kode, out := u.aksi(id, admin, "", aksi, n, param, m)
		u.wajib(kode, http.StatusOK, out, aksi)
		return out
	}
	langkah("SetValueToClaim", 0, masterUji+"|UJI-TG|UJI-COB", nil)
	langkah("CheckNoPolicy", 0, polisUji, nil)
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
	// tambah / hapus baris Spreading Claim nonaktif (keputusan work owner 07-10-2026): dipanggil langsung pun ditolak
	kode, out := u.aksi(id, admin, "", models.DaftarSpreading+"#tambah", 0, "", nil)
	u.wajib(kode, http.StatusBadRequest, out, "tambah baris spreading")
	// tanpa baris spreading, Save to issue RNM ditolak ProteksiData_act langkah 5 - dan layar tidak dapat
	// melahirkannya. Kasus hasil pemuat data lama membawa barisnya: fixture setara.
	kode, out = u.aksi(id, admin, "", "SaveOutstanding", 0, "", nil)
	u.wajib(kode, http.StatusUnprocessableEntity, out, "Save to issue RNM tanpa baris spreading")
	if !strings.Contains(strings.Join(teks(out["pesan"]), ";"), models.PesanIsiSpreading) {
		t.Fatalf("pesan tanpa spreading: %v", out["pesan"])
	}
	h0 := u.g.Halaman(id)
	h0.SetelDaftar(models.DaftarSpreading, []models.Baris{{"CurrencyID": matauang, "Currency": "UJA"}})
	u.g.SetelHalaman(id, h0)
	langkah("CountSpreading", 1, "", map[string]string{models.JalurAnak(models.DaftarSpreading, 1, "TreatyName"): "UJI-QS",
		models.JalurAnak(models.DaftarSpreading, 1, "SharePercentage"): "100"})
	if s := u.g.Halaman(id).AmbilDaftar(models.DaftarSpreading); len(s) != 1 || s[0]["SharePercentage"] != "100" {
		t.Fatalf("baris Spreading Claim tersimpan %+v, mau satu baris share 100", s)
	}

	h := u.g.Halaman(id)
	if v := h.AmbilDaftar(models.DaftarLossAlloc)[0]["TreatyType"]; v != jenisQSUji {
		t.Fatalf("SetNameTreaty: TreatyType %q, mau %q", v, jenisQSUji)
	}
	if v := h.AmbilDaftar(models.DaftarEstimasi)[0]["GrossEstimationPct"]; v != "1000" {
		t.Fatalf("AddEstimation: gross %q, mau 1000 (100%% x claim amount 1000)", v)
	}
	if h.Ambil("IsCFS") != "1" {
		t.Fatalf("CountEstimation: IsCFS %q, mau 1", h.Ambil("IsCFS"))
	}

	out = langkah("SaveOutstanding", 0, "", nil)
	h = u.g.Halaman(id)
	if !riwayatMemuat(h, models.TeksSaveOutstanding) {
		t.Fatalf("Claim History tidak menyimpan %q: %+v", models.TeksSaveOutstanding, h.AmbilDaftar(models.DaftarRiwayat))
	}
	if no := h.Ambil(models.CD + "NoClaim"); no != "UJI-K12.05.2026.T00001" {
		t.Fatalf("nomor klaim %q", no)
	}
	if len(u.g.OS) != 1 || u.g.OS[0].StsReject != models.StsOSEstimasi || u.g.OS[0].CaseID != id {
		t.Fatalf("baris OS %+v", u.g.OS)
	}
	if _, ada := u.g.JSONKlaim[id]; !ada || len(u.g.Log) != 1 {
		t.Fatalf("JSON_KLAIM / log layanan tidak ditulis")
	}
	if h.Ambil("IsOutstanding") != "1" || h.Ambil("IsCFS") != "" || h.AmbilDaftar(models.DaftarEstimasi)[0]["PrintFaceClaim"] != "1" {
		t.Fatalf("penanda sesudah Save to issue RNM salah")
	}
	if out["info"] != models.PesanPrintPla {
		t.Errorf("IsPLA 1 (estimasi melewati batas PLA kosong) mau local action PrintFile: %v", out["info"])
	}
	if len(u.g.Efek) != 0 {
		t.Errorf("non-produksi tetap mengantre efek keluar: %v", u.g.Efek)
	}
	// sekali lagi tanpa estimasi baru: tombol nonaktif (IsCFS kosong)
	kode, out = u.aksi(id, admin, "", "SaveOutstanding", 0, "", nil)
	u.wajib(kode, http.StatusConflict, out, "Save to issue RNM tanpa estimasi baru")

	langkah("CheckNopolicy", 0, "", nil)
	if u.g.Halaman(id).Ambil("IsAcceptation") != "1" {
		t.Fatalf("Send to Acceptation tidak menyetel IsAcceptation")
	}
	langkah("SubmitOutstanding", 0, "", nil)
	if k := u.g.Kasus[id]; k.Tahap != models.TahapAcceptation || k.Posisi != models.WorkbasketAcceptation {
		t.Fatalf("Submit: %+v", k)
	}

	// pembuat tanpa peran workbasket tidak lagi memegang kasus
	kode, out = u.aksi(id, admin, "", "Simpan", 0, "", nil)
	u.wajib(kode, http.StatusForbidden, out, "pembuat sesudah Submit")
	kerja := func(aksi string, n int, param string, m map[string]string) (int, map[string]any) {
		return u.aksi(id, teknik, models.WorkbasketAcceptation, aksi, n, param, m)
	}
	kode, out = kerja("AddAdjustment", 0, "", nil)
	u.wajib(kode, http.StatusOK, out, "AddAdjustment")
	if n := len(u.g.Halaman(id).AmbilDaftar(models.DaftarAdjustment)); n != 1 {
		t.Fatalf("baris adjustment %d, mau 1", n)
	}
	// penyerahan komite tidak dibangun (OQ-CP-16): dipanggil langsung pun ditolak, nol kasus TKMT- lahir
	kode, out = kerja("AddKomiteTreatyChild", 1, "", nil)
	u.wajib(kode, http.StatusBadRequest, out, "penyerahan komite")
	for kid := range u.g.Kasus {
		if strings.HasPrefix(kid, "TKMT-") {
			t.Fatalf("kasus komite %s lahir", kid)
		}
	}
	// tutup klaim ditolak selama ada adjustment berstatus kosong
	kode, out = kerja("CloseClaimProp", 0, "", map[string]string{"TempCommiteClaim.Remarks": "UJI-TUTUP"})
	u.wajib(kode, http.StatusUnprocessableEntity, out, "tutup dengan adjustment tertunda")
	if !strings.Contains(strings.Join(teks(out["pesan"]), ";"), models.PesanKomiteMasihJalan) {
		t.Fatalf("pesan tutup klaim: %v", out["pesan"])
	}
	kode, out = kerja("DeleteAjsutment", 1, "", nil)
	u.wajib(kode, http.StatusOK, out, "hapus adjustment")
	kode, out = kerja("CloseClaimProp", 0, "", map[string]string{"TempCommiteClaim.Remarks": "UJI-TUTUP"})
	u.wajib(kode, http.StatusOK, out, "tutup klaim")
	if !u.g.Kasus[id].Tertutup() {
		t.Fatalf("kasus belum Resolved-Completed")
	}
	if last := u.g.OS[len(u.g.OS)-1]; last.StsReject != models.StsOSTutupBerkas {
		t.Fatalf("baris OS tutup %+v", last)
	}
	kode, out = kerja("Simpan", 0, "", nil)
	u.wajib(kode, http.StatusConflict, out, "tulis ke kasus tertutup")
}

// riwayatMemuat - baris SuggestList tersimpan berkomentar `isi` (tanpa penanda baris baru).
func riwayatMemuat(h *models.Halaman, isi string) bool {
	for _, b := range h.AmbilDaftar(models.DaftarRiwayat) {
		if b["CommentSuggest"] == isi && b[models.PropRiwayatBaru] == "" && b["No"] != "" {
			return true
		}
	}
	return false
}

func teks(v any) []string {
	var out []string
	if xs, ok := v.([]any); ok {
		for _, x := range xs {
			out = append(out, x.(string))
		}
	}
	return out
}
