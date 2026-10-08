package handlers_test

// Uji seam HTTP Komite Claim Prop (handlers -> services -> tiruan + kontrak palsu Claim Prop). Lima skenario wajib
// prompt §9: setuju -> setuju (nomor terbit), subjectivity di tingkat 1, tolak di tingkat 1, penyetuju bukan pemilik
// ditolak, kasus tertutup -> 409.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/komiteclaimprop/backend/handlers"
	"nusantarare/modul/komiteclaimprop/backend/models"
	"nusantarare/modul/komiteclaimprop/backend/services"
	"nusantarare/modul/komiteclaimprop/backend/tiruan"
)

type uji struct {
	t   *testing.T
	g   *tiruan.Gudang
	a   *tiruan.Acuan
	l   *services.Layanan
	srv http.Handler
}

func siap(t *testing.T, n int, produksi bool) *uji {
	t.Helper()
	g, a := tiruan.Baru(), tiruan.AcuanBaru()
	tiruan.SiapkanUji(g, a, tiruan.KomiteUji, tiruan.AdjUji, n)
	l := services.Baru(g, a, g.Klaim, func() time.Time { return tiruan.SaatUji }, produksi).
		DenganKasir(models.KonfigurasiKasir{CompanyName: "UJI-CO", LjtdID: "UJI-LJTD", LdcID: "UJI-LDC"}).
		DenganEmail(models.KonfigurasiEmail{Akun: "UJI-AKUN", AkunSyariah: "UJI-SYARIAH", CC: "uji.cc@contoh.invalid"})
	return &uji{t: t, g: g, a: a, l: l, srv: handlers.Router(l, true)}
}

func (u *uji) minta(metode, jalur, pelaku string, badan any) *httptest.ResponseRecorder {
	u.t.Helper()
	var b []byte
	if badan != nil {
		b, _ = json.Marshal(badan)
	}
	r := httptest.NewRequest(metode, handlers.Prefix+jalur, bytes.NewReader(b))
	r.Header.Set("X-Pelaku", pelaku)
	w := httptest.NewRecorder()
	u.srv.ServeHTTP(w, r)
	return w
}

func (u *uji) putus(pelaku string, kep models.Keputusan, mau int) services.HasilKeputusan {
	u.t.Helper()
	w := u.minta("POST", "/kasus/"+tiruan.KomiteUji+"/putuskan", pelaku, kep)
	if w.Code != mau {
		u.t.Fatalf("putuskan %s: %d (mau %d) %s", pelaku, w.Code, mau, w.Body.String())
	}
	var h services.HasilKeputusan
	_ = json.Unmarshal(w.Body.Bytes(), &h)
	return h
}

func (u *uji) adj() map[string]string {
	for _, b := range u.g.Klaim.Daftar(tiruan.KlaimUji, "ClaimData.AdjustmentList") {
		if b["ID"] == tiruan.AdjUji {
			return b
		}
	}
	u.t.Fatal("baris adjustment uji hilang")
	return nil
}

func (u *uji) daftarKerja(pelaku string) []models.BarisKerja {
	u.t.Helper()
	w := u.minta("GET", "/kasus", pelaku, nil)
	if w.Code != http.StatusOK {
		u.t.Fatalf("daftar kerja %s: %d %s", pelaku, w.Code, w.Body.String())
	}
	var out []models.BarisKerja
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		u.t.Fatal(err)
	}
	return out
}

func setuju(komentar string) models.Keputusan {
	return models.Keputusan{AcceptStatus: models.KeputusanSetuju, Comment: komentar}
}

// Skenario 1 - setuju -> setuju: tingkat akhir menerbitkan nomor akseptasi sekali, menulis klaim induk lewat kontrak,
// OS / JSON / log / riwayat, lalu kasus selesai.
func TestSetujuLaluSetujuNomorTerbit(t *testing.T) {
	u := siap(t, 2, false)
	if d := u.daftarKerja(tiruan.PenyetujuUji(1)); len(d) != 1 || d[0].KasusID != tiruan.KomiteUji || d[0].Tingkat != 1 {
		t.Fatalf("daftar kerja UJI-K1: %+v", d)
	}
	if d := u.daftarKerja(tiruan.PenyetujuUji(2)); len(d) != 0 {
		t.Fatalf("tingkat 2 belum boleh melihat kasus: %+v", d)
	}
	h := u.putus(tiruan.PenyetujuUji(1), setuju("UJI setuju satu"), http.StatusOK)
	if h.Selesai || h.KomiteCount != 2 || h.AcceptedNo != "" {
		t.Fatalf("tingkat 1: %+v", h)
	}
	if d := u.daftarKerja(tiruan.PenyetujuUji(2)); len(d) != 1 || d[0].Tingkat != 2 {
		t.Fatalf("giliran pindah ke UJI-K2: %+v", d)
	}
	if len(u.g.OS) != 0 || u.adj()["AcceptedNo"] != "" {
		t.Fatal("bukan tingkat akhir: nol OS, nol nomor")
	}
	h = u.putus(tiruan.PenyetujuUji(2), setuju("UJI setuju dua"), http.StatusOK)
	if !h.Selesai || h.KomiteCount != 3 {
		t.Fatalf("tingkat akhir: %+v", h)
	}
	if h.AcceptedNo != "UJIA12.10.2026.TP00001" {
		t.Fatalf("nomor akseptasi %q", h.AcceptedNo)
	}
	b := u.adj()
	if b["AcceptedNo"] != h.AcceptedNo || b["AcceptanceStatus"] != "1" || b["IsApproved"] != "1" ||
		b["Notes"] != "UJI setuju dua" || b["IsPrintAccept"] != "1" || b["IsFacRetro"] != "1" {
		t.Fatalf("adjustment induk: %+v", b)
	}
	n := u.g.Klaim.Nilai(tiruan.KlaimUji)
	if n["IsAnyAcceptation"] != "1" || n["AktifButton"] != "0" || n["IsOutstanding"] != "1" ||
		n["ClaimData.IsCloseFile"] != "false" {
		t.Fatalf("header induk: %+v", n)
	}
	if len(u.g.OS) != 1 || u.g.OS[0].AcceptedNo != h.AcceptedNo || u.g.OS[0].StsReject != "1" ||
		u.g.OS[0].CaseID != tiruan.KlaimUji {
		t.Fatalf("OS_AKSEPTASI_KLAIM: %+v", u.g.OS)
	}
	if dj := u.g.OS[0].DataJSON; !json.Valid([]byte(dj)) || !strings.Contains(dj, `"KomiteNo":"TKMT-UJI001"`) ||
		!strings.Contains(dj, `"AcceptedNo":"`+h.AcceptedNo+`"`) {
		t.Fatalf("DATA_JSON OS (halaman TempOSAkseptasi): %q", dj)
	}
	if _, ada := u.g.JSONKlaim[tiruan.KlaimUji]; !ada || len(u.g.Log) != 1 || u.g.Log[0].JenisService != "AKSEPTASI" ||
		u.g.Log[0].NoAkseptasi != h.AcceptedNo { // S30 Param.NoAkseptasi = OutputData.START_DATE (S16.8)
		t.Fatalf("JSON_KLAIM / log: %+v %+v", u.g.JSONKlaim, u.g.Log)
	}
	if len(u.g.Riwayat) != 2 || u.g.Riwayat[0].Status != "ACCEPT" || u.g.Riwayat[0].Username != "UJI Penyetuju Satu" ||
		u.g.Riwayat[1].IDKomite != tiruan.KomiteUji {
		t.Fatalf("HISTORYAKSEPTASIPEGA: %+v", u.g.Riwayat)
	}
	riw := u.g.Klaim.Daftar(tiruan.KlaimUji, "ClaimData.SuggestList")
	if len(riw) != 2 || riw[1]["CommentSuggest"] != "Accepted by UJI-JABATAN-2" {
		t.Fatalf("riwayat klaim: %+v", riw)
	}
	if retro := u.g.Klaim.Daftar(tiruan.KlaimUji, "ClaimData.FacRetroList"); len(retro) != 1 ||
		retro[0]["ReinsurerID"] != "UJI-R1" {
		t.Fatalf("FacRetroList: %+v", retro)
	}
	if k := u.g.Kasus[tiruan.KomiteUji]; k.StatusWork != models.StatusSelesai || k.AcceptStatus != "1" {
		t.Fatalf("kasus komite: %+v", k)
	}
	if len(u.g.Efek) != 0 {
		t.Fatalf("di luar produksi nol efek keluar: %v", u.g.Efek)
	}
	if d := u.daftarKerja(tiruan.PenyetujuUji(2)); len(d) != 0 {
		t.Fatalf("kasus selesai keluar dari daftar kerja: %+v", d)
	}
}

// Skenario 2 - subjectivity di tingkat 1 (OQ-KCP-01 "a", migrasi 682): isian tingkat 1 disimpan di header kasus komite
// dan dipakai tingkat akhir - tanpa nomor, tanpa OS, IsKomite 0. Tangga satu dan dua tingkat.
func TestSubjectivityTingkatSatu(t *testing.T) {
	u := siap(t, 2, false)
	kep := models.Keputusan{AcceptStatus: "1", Comment: "UJI", IsSubjectivity: true, SubjectivityNote: "3"}
	if h := u.putus(tiruan.PenyetujuUji(1), kep, http.StatusOK); h.Selesai {
		t.Fatalf("tingkat 1 dari 2: %+v", h)
	}
	if k := u.g.Kasus[tiruan.KomiteUji]; k.Subjectivity != "1" || k.SubjectivityNote != "3" {
		t.Fatalf("isian tingkat 1 tersimpan di header: %+v", k)
	}
	// tingkat 2: isian Subjectivity nonaktif (yang dikirim diabaikan), yang tersimpan dipakai.
	if h := u.putus(tiruan.PenyetujuUji(2), models.Keputusan{AcceptStatus: "1", Comment: "UJI 2"}, http.StatusOK); !h.Selesai ||
		h.AcceptedNo != "" {
		t.Fatalf("tingkat akhir subjectivity: %+v", h)
	}
	if b := u.adj(); b["IsSubjectivity"] != "true" || b["SubjectivityNote"] != "3" || b["IsKomite"] != "0" {
		t.Fatalf("adjustment subjectivity tangga dua tingkat: %+v", b)
	}
	if len(u.g.OS) != 0 {
		t.Fatal("subjectivity: nol OS")
	}
	u = siap(t, 1, false)
	h := u.putus(tiruan.PenyetujuUji(1), kep, http.StatusOK)
	if !h.Selesai || h.AcceptedNo != "" {
		t.Fatalf("subjectivity: %+v", h)
	}
	b := u.adj()
	if b["IsSubjectivity"] != "true" || b["SubjectivityNote"] != "3" || b["IsKomite"] != "0" ||
		b["AcceptedNo"] != "" {
		t.Fatalf("adjustment subjectivity: %+v", b)
	}
	if len(u.g.OS) != 0 || u.g.Klaim.Nilai(tiruan.KlaimUji)["ClaimData.IsSubjectivity"] != "true" {
		t.Fatal("subjectivity: nol OS, header IsSubjectivity true")
	}
	kep.SubjectivityNote = ""
	u = siap(t, 1, false)
	if w := u.minta("POST", "/kasus/"+tiruan.KomiteUji+"/putuskan", tiruan.PenyetujuUji(1), kep); w.Code !=
		http.StatusUnprocessableEntity || !strings.Contains(w.Body.String(), "Subjectivity Note") {
		t.Fatalf("catatan subjectivity wajib: %d %s", w.Code, w.Body.String())
	}
}

// Skenario 3 - tolak di tingkat 1: tangga berhenti, sisa penyetuju ditolak otomatis, adjustment ditolak.
func TestTolakTingkatSatu(t *testing.T) {
	u := siap(t, 2, false)
	h := u.putus(tiruan.PenyetujuUji(1), models.Keputusan{AcceptStatus: "2", Comment: "UJI tolak",
		UsulTutup: true}, http.StatusOK)
	if !h.Selesai || h.AcceptedNo != "" || h.KomiteCount != 3 {
		t.Fatalf("tolak: %+v", h)
	}
	t2 := u.g.Tangga[tiruan.KomiteUji]
	if t2[0].Keputusan != "2" || t2[0].Komentar != "UJI tolak" || t2[1].Keputusan != "2" || t2[1].Komentar != "" ||
		t2[1].Tanggal == "" {
		t.Fatalf("tangga sesudah tolak: %+v", t2)
	}
	b := u.adj()
	if b["AcceptanceStatus"] != "2" || b["Notes"] != "UJI tolak" || b["AcceptedNo"] != "" {
		t.Fatalf("adjustment ditolak: %+v", b)
	}
	n := u.g.Klaim.Nilai(tiruan.KlaimUji)
	if n["AktifButton"] != "0" || n["ClaimData.IsCloseFile"] != "true" || n["IsAnyAcceptation"] == "1" {
		t.Fatalf("header ditolak: %+v", n)
	}
	if len(u.g.Riwayat) != 1 || u.g.Riwayat[0].Status != "REJECT" || len(u.g.OS) != 0 {
		t.Fatalf("riwayat tolak: %+v OS %d", u.g.Riwayat, len(u.g.OS))
	}
	if k := u.g.Kasus[tiruan.KomiteUji]; k.UsulTutup != "1" || k.StatusWork != models.StatusSelesai {
		t.Fatalf("kepala kasus: %+v", k)
	}
}

// Skenario 4 - penyetuju bukan pemilik tingkat berjalan ditolak (403) tanpa satu tulisan pun.
func TestBukanPemilikDitolak(t *testing.T) {
	u := siap(t, 2, false)
	u.putus(tiruan.PenyetujuUji(2), setuju("UJI"), http.StatusForbidden)
	u.putus("UJI-LAIN", setuju("UJI"), http.StatusForbidden)
	if k := u.g.Kasus[tiruan.KomiteUji]; k.Count != 1 || u.g.Tangga[tiruan.KomiteUji][0].Keputusan != "0" ||
		len(u.g.Riwayat) != 0 {
		t.Fatalf("nol tulisan: %+v", k)
	}
	if w := u.minta("POST", "/kasus/"+tiruan.KomiteUji+"/putuskan", "", setuju("UJI")); w.Code != http.StatusUnauthorized {
		t.Fatalf("tanpa identitas: %d", w.Code)
	}
}

// Skenario 5 - kasus tertutup -> 409 (juga bila klaim induknya yang tertutup).
func TestKasusTertutup409(t *testing.T) {
	u := siap(t, 1, false)
	u.putus(tiruan.PenyetujuUji(1), setuju("UJI"), http.StatusOK)
	u.putus(tiruan.PenyetujuUji(1), setuju("UJI"), http.StatusConflict)
	u = siap(t, 2, false)
	u.g.Klaim.Tutup(tiruan.KlaimUji)
	u.putus(tiruan.PenyetujuUji(1), setuju("UJI"), http.StatusConflict)
	if u.g.Tangga[tiruan.KomiteUji][0].Keputusan != "0" {
		t.Fatal("transaksi batal utuh")
	}
}

func TestValidasiDanBukaKasus(t *testing.T) {
	u := siap(t, 2, false)
	if w := u.minta("POST", "/kasus/"+tiruan.KomiteUji+"/putuskan", tiruan.PenyetujuUji(1),
		models.Keputusan{AcceptStatus: "1"}); w.Code != http.StatusUnprocessableEntity ||
		!strings.Contains(w.Body.String(), "Note: Value cannot be blank") {
		t.Fatalf("Note wajib: %d %s", w.Code, w.Body.String())
	}
	if w := u.minta("GET", "/kasus/TKMT-TIDAKADA", tiruan.PenyetujuUji(1), nil); w.Code != http.StatusNotFound {
		t.Fatalf("kasus tak dikenal: %d", w.Code)
	}
	if w := u.minta("GET", "/kasus/CLMP-UJI001", tiruan.PenyetujuUji(1), nil); w.Code != http.StatusNotFound {
		t.Fatalf("bukan TKMT-: %d", w.Code)
	}
	w := u.minta("GET", "/kasus/"+tiruan.KomiteUji, tiruan.PenyetujuUji(1), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("buka kasus: %d %s", w.Code, w.Body.String())
	}
	var ly models.Layar
	if err := json.Unmarshal(w.Body.Bytes(), &ly); err != nil {
		t.Fatal(err)
	}
	if !ly.BolehKerja || ly.Judul[0] != "CLAIM COMMITTEE -" || ly.Judul[1] != "ADJUSTMENT" || !ly.Isian.Terbuka {
		t.Fatalf("layar: %+v", ly)
	}
	if !strings.Contains(w.Body.String(), `"Committe Accept Status"`) ||
		!strings.Contains(w.Body.String(), `"Total in IDR"`) || !strings.Contains(w.Body.String(), `"Dedutible Type"`) {
		t.Fatal("layar memuat grid tangga, total per mata uang, dan blok deductible (VERBATIM)")
	}
	// prompt values (diekspor 08-10-2026): KomiteAproval, AcceptanceStatus, SubjectivityNote
	if !strings.Contains(w.Body.String(), `"keputusan":"Waiting"`) ||
		!strings.Contains(w.Body.String(), `"AcceptanceStatus":"Approve"`) ||
		!strings.Contains(w.Body.String(), `{"nilai":"1","label":"Treaty Leader Approval"}`) {
		t.Fatalf("label prompt values: %s", w.Body.String())
	}
	w = u.minta("GET", "/kasus/"+tiruan.KomiteUji, tiruan.PenyetujuUji(2), nil)
	_ = json.Unmarshal(w.Body.Bytes(), &ly)
	for _, tb := range ly.Tombol {
		if tb.Label == "Submit" && tb.Aktif {
			t.Fatal("Submit nonaktif bagi bukan pemegang")
		}
		if tb.Label == "View more details" && !tb.Aktif {
			t.Fatal("View more details aktif juga bagi bukan pemegang (hanya-baca, ViewClaimFormKomite)")
		}
	}
}

// Produksi: konversi, Kasir (DirectToKasir), dan email diantre di transaksi yang sama.
func TestEfekKeluarDiProduksi(t *testing.T) {
	u := siap(t, 1, true)
	nilai := u.g.Klaim.Nilai(tiruan.KlaimUji)
	daftar := map[string][]map[string]string{}
	for _, j := range []string{"ClaimData.AdjustmentList", "ClaimData.AdjustmentList(2).SpreadingAdjustment"} {
		daftar[j] = u.g.Klaim.Daftar(tiruan.KlaimUji, j)
	}
	daftar["ClaimData.AdjustmentList"][1]["DirectToKasir"] = "true"
	daftar["ClaimData.AdjustmentList"][1]["NameOfBank"] = "UJI BANK"
	u.g.Klaim.Setel(tiruan.KlaimUji, nilai, daftar)
	u.a.Bank["UJI BANK||"] = "UJI-BANK-1"
	u.a.Email["UJI-CED"] = "uji@contoh.invalid"
	u.g.KodeProduksi = "UJIX"                  // S14.1: @length(AcceptedNo) 23 / 24
	u.a.Konversi["UJIXA12102026TP00001"] = "1" // S3 getStatusKonversi_Act (IsPEGAPROD): sudah dikonversi
	h := u.putus(tiruan.PenyetujuUji(1), setuju("UJI"), http.StatusOK)
	jenis := map[string]bool{}
	for _, e := range u.g.Efek {
		jenis[strings.SplitN(e, ":", 2)[0]] = true
	}
	if !jenis[services.JenisEfekKonversi] || !jenis[services.JenisEfekKasir] || !jenis[services.JenisEfekEmailKomite] {
		t.Fatalf("efek produksi: %v", u.g.Efek)
	}
	if u.adj()["IDOfBank"] != "UJI-BANK-1" || len(h.AcceptedNo) != 23 {
		t.Fatalf("IDOfBank S12-S13: %+v", u.adj())
	}
	if !jenis[services.JenisEfekDokumen] {
		t.Fatalf("S21 -> PrintFileAcceptance_TKMT: dokumen akseptasi diantre: %v", u.g.Efek)
	}
	// MUATAN hanya pengenal (claimlife/015): tanpa alamat email, subjek, maupun nama.
	muatan := map[string]services.MuatanOutbox{}
	for _, e := range u.g.Efek {
		p := strings.SplitN(e, ":", 3)
		if p[0] == services.JenisEfekEmailKomite || p[0] == services.JenisEfekDokumen {
			if strings.Contains(p[2], "@") || strings.Contains(p[2], "UJI Penyetuju") || strings.Contains(p[2], "subjek") {
				t.Fatalf("MUATAN %s memuat alamat / nama: %s", p[0], p[2])
			}
			var m services.MuatanOutbox
			if err := json.Unmarshal([]byte(p[2]), &m); err != nil {
				t.Fatal(err)
			}
			muatan[p[0]] = m
		}
	}
	isi := func(j string) map[string]string {
		out := map[string]string{}
		for k, v := range muatan[j].Isi.(map[string]any) {
			out[k], _ = v.(string)
		}
		return out
	}
	ctx := context.Background()
	u.a.Nama[tiruan.PembuatUji] = "UJI Pembuat"
	u.a.Surel[tiruan.PembuatUji] = "uji.pembuat.syariah@contoh.invalid"
	sr, err := u.l.SusunEmailKomite(ctx, tiruan.KomiteUji, isi(services.JenisEfekEmailKomite))
	if err != nil {
		t.Fatal(err)
	}
	// S13-S14 pembuat; S18 akun syariah; S4 CC; S16 badan EmailKlaim_HTML_KMT
	if sr.Kepada != "uji.pembuat.syariah@contoh.invalid" || sr.Akun != "UJI-SYARIAH" || sr.CC != "uji.cc@contoh.invalid" ||
		!strings.HasPrefix(sr.Subjek, "(Approval) Pengajuan Akseptasi : ") ||
		!strings.Contains(sr.HTML, "Dear <strong>UJI Pembuat</strong>") ||
		!strings.Contains(sr.HTML, "Accepted No: <span style=\"font-weight: normal;\">"+h.AcceptedNo) ||
		!strings.Contains(sr.HTML, "<strong>UJI Penyetuju Satu</strong>") || !strings.Contains(sr.HTML, "Status: <span") {
		t.Fatalf("email komite: %+v", sr)
	}
	dk, err := u.l.SusunDokumenAkseptasi(ctx, tiruan.KomiteUji, isi(services.JenisEfekDokumen))
	if err != nil {
		t.Fatal(err)
	}
	if dk.NamaBerkas != "Persetujuan Klaim   AcceptNo "+h.AcceptedNo+".pdf" || dk.Kategori != "AcceptanceNote" ||
		dk.Folder != "Claim" || !strings.Contains(dk.HTML, "Create by : UJI Penyetuju Satu") ||
		!strings.Contains(dk.HTML, "Jakarta, 08 October 2026") || !strings.Contains(dk.HTML, h.AcceptedNo) {
		t.Fatalf("dokumen akseptasi: %+v", dk)
	}
	// Pelaksana produksi: isi dirakit, lalu panggilan nyata ditahan sampai disetujui manusia.
	pl := services.PelaksanaKomiteClaimProp{Lingkungan: inti.Produksi, Penyusun: u.l}
	for j, mau := range map[string]error{services.JenisEfekEmailKomite: outbox.ErrEmailBelumDisetujui,
		services.JenisEfekDokumen: outbox.ErrPenyimpananBelumDisetujui} {
		b, _ := json.Marshal(muatan[j])
		err := pl.Laksanakan(ctx, nil, outbox.BarisEfekKeluar{Modul: "komiteclaimprop", Jenis: j, Muatan: string(b)})
		if !errors.Is(err, mau) {
			t.Fatalf("pelaksana %s: %v", j, err)
		}
	}
}

// OQ-KCP-06 "a": penyerahan ulang baris subjectivity - putaran kedua disetujui tanpa syarat menerbitkan nomor; `.Comment`
// awal = komentar anggota pertama putaran pertama (AddKomiteTreatyChild_ACT S16); keputusan di baris terakhir (S7).
func TestKirimUlangSubjectivity(t *testing.T) {
	u := siap(t, 1, false)
	u.putus(tiruan.PenyetujuUji(1), models.Keputusan{AcceptStatus: "1", Comment: "UJI putaran 1", IsSubjectivity: true,
		SubjectivityNote: "3"}, http.StatusOK)
	if b := u.adj(); b["IsSubjectivity"] != "true" || b["IsKomite"] != "0" || b["AcceptedNo"] != "" {
		t.Fatalf("putaran 1 subjectivity: %+v", b)
	}
	// Claim Prop menyerahkan ulang (kasus baru, KOMITE_ID ditimpa) ke jenjang terbawah.
	u.g.Lahirkan("TKMT-UJI002", tiruan.KlaimUji, tiruan.AdjUji, tiruan.PembuatUji, "UJI Admin",
		[]models.Anggota{{OperatorID: tiruan.PenyetujuUji(1), Jabatan: "UJI-JABATAN-1"}}, tiruan.SaatUji.Add(time.Hour))
	w := u.minta("GET", "/kasus/TKMT-UJI002", tiruan.PenyetujuUji(1), nil)
	var ly models.Layar
	if err := json.Unmarshal(w.Body.Bytes(), &ly); err != nil || ly.Isian.Nilai.Comment != "UJI putaran 1" {
		t.Fatalf("S16 komentar awal putaran kedua: %q (%d)", ly.Isian.Nilai.Comment, w.Code)
	}
	r := u.minta("POST", "/kasus/TKMT-UJI002/putuskan", tiruan.PenyetujuUji(1),
		models.Keputusan{AcceptStatus: "1", Comment: "UJI putaran 2"})
	var h services.HasilKeputusan
	_ = json.Unmarshal(r.Body.Bytes(), &h)
	if r.Code != http.StatusOK || !h.Selesai || h.AcceptedNo == "" {
		t.Fatalf("putaran 2 tanpa syarat menerbitkan nomor: %d %s", r.Code, r.Body.String())
	}
	if b := u.adj(); b["IsSubjectivity"] != "false" || b["AcceptedNo"] != h.AcceptedNo || b["AcceptanceStatus"] != "1" {
		t.Fatalf("adjustment sesudah putaran 2: %+v", b)
	}
	if tg := u.g.Tangga["TKMT-UJI002"]; tg[len(tg)-1].Keputusan != "1" || tg[len(tg)-1].Komentar != "UJI putaran 2" {
		t.Fatalf("S7 baris terakhir: %+v", tg)
	}
}
