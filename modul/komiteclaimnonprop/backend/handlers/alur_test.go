package handlers_test

// Uji seam HTTP Komite Claim Non Prop (handlers -> services -> tiruan + kontrak palsu Claim Non Prop): setuju -> setuju
// (nomor terbit, OS + CLAIMXOL2 + OS subjectivity), subjectivity di tingkat 1, tolak di tingkat 1 (EXIT), penyetuju
// bukan pemilik ditolak, kasus tertutup -> 409, efek keluar di produksi, kirim ulang subjectivity.

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
	"nusantarare/modul/komiteclaimnonprop/backend/handlers"
	"nusantarare/modul/komiteclaimnonprop/backend/models"
	"nusantarare/modul/komiteclaimnonprop/backend/services"
	"nusantarare/modul/komiteclaimnonprop/backend/tiruan"
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
	u.t.Fatal("baris akseptasi uji hilang")
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

// Setuju -> setuju: tingkat akhir menerbitkan nomor akseptasi sekali (S14.8-S14.13), menulis klaim induk lewat kontrak,
// OS + CLAIMXOL2 (S14.18), OS subjectivity (S17), riwayat, JSON_KLAIM, lalu kasus selesai (S28).
func TestSetujuLaluSetujuNomorTerbit(t *testing.T) {
	u := siap(t, 2, false)
	if d := u.daftarKerja(tiruan.PenyetujuUji(1)); len(d) != 1 || d[0].KasusID != tiruan.KomiteUji || d[0].Tingkat != 1 ||
		d[0].NoKlaim != "UJI-K-0001" {
		t.Fatalf("daftar kerja UJI-K1: %+v", d)
	}
	if d := u.daftarKerja(tiruan.PenyetujuUji(2)); len(d) != 0 {
		t.Fatalf("tingkat 2 belum boleh melihat kasus: %+v", d)
	}
	h := u.putus(tiruan.PenyetujuUji(1), setuju("UJI setuju satu"), http.StatusOK)
	if h.Selesai || h.KomiteCount != 2 || h.AcceptedNo != "" || u.g.Posisi[tiruan.KomiteUji] != tiruan.PenyetujuUji(2) {
		t.Fatalf("tingkat 1: %+v posisi %q", h, u.g.Posisi[tiruan.KomiteUji])
	}
	if d := u.daftarKerja(tiruan.PenyetujuUji(2)); len(d) != 1 || d[0].Tingkat != 2 {
		t.Fatalf("giliran pindah ke UJI-K2: %+v", d)
	}
	if len(u.g.OS) != 0 || len(u.g.OSSubj) != 0 || u.adj()["AcceptedNo"] != "" {
		t.Fatal("bukan tingkat akhir: nol OS, nol nomor")
	}
	h = u.putus(tiruan.PenyetujuUji(2), setuju("UJI setuju dua"), http.StatusOK)
	if !h.Selesai || h.KomiteCount != 3 || h.AcceptedNo != "UJIA123.10.2026.TX00001" {
		t.Fatalf("tingkat akhir: %+v", h)
	}
	b := u.adj()
	if b["AcceptedNo"] != h.AcceptedNo || b["AcceptanceStatus"] != "1" || b["AcceptedDate"] != "2026-10-08 10:00:00" ||
		b["IsSubjectivity"] != "false" {
		t.Fatalf("akseptasi induk: %+v", b)
	}
	n := u.g.Klaim.Nilai(tiruan.KlaimUji)
	if n["CNPStatusCase"] != "CLAIM ACCEPTED" || n["IsCloseFile"] != "0" {
		t.Fatalf("header induk: %+v", n)
	}
	if len(u.g.OS) != 1 || u.g.OS[0].StsReject != "4" || u.g.OS[0].CaseID != tiruan.KlaimUji ||
		u.g.OS[0].MasterID != "UJI-MASTER" {
		t.Fatalf("OS_AKSEPTASI_KLAIM: %+v", u.g.OS)
	}
	if dj := u.g.OS[0].DataJSON; !json.Valid([]byte(dj)) || !strings.Contains(dj, `"AcceptedNo":"`+h.AcceptedNo+`"`) ||
		!strings.Contains(dj, `"XOL":"UJI XOL 1"`) || strings.Contains(dj, `"XOL":"UR"`) {
		t.Fatalf("DATA_JSON OS (TempOSAkseptasi + CNPLayerList): %q", dj)
	}
	if len(u.g.XOL2) != 1 || u.g.XOL2[0].XOL != "UJI XOL 1" || u.g.XOL2[0].GrossAdjustment != "300" ||
		u.g.XOL2[0].CNPReinstatement != "10" {
		t.Fatalf("CLAIMXOL2 tanpa layer UR: %+v", u.g.XOL2)
	}
	if s := u.g.OSSubj[tiruan.KlaimUji]; s.StsSubjectivity != "0" || s.DataJSON != u.g.OS[0].DataJSON {
		t.Fatalf("OS_AKSEPTASI_SUBJECTIVITY: %+v", s)
	}
	if _, ada := u.g.JSONKlaim[tiruan.KlaimUji]; !ada {
		t.Fatal("JSON_KLAIM S24")
	}
	if len(u.g.Riwayat) != 2 || u.g.Riwayat[0].Status != "ACCEPT" || u.g.Riwayat[0].Username != "UJI Penyetuju Satu" ||
		u.g.Riwayat[1].IDKomite != tiruan.KomiteUji || u.g.Riwayat[1].Workbasket != "KLAIM" {
		t.Fatalf("HISTORYAKSEPTASIPEGA: %+v", u.g.Riwayat)
	}
	riw := u.g.Klaim.Daftar(tiruan.KlaimUji, "ClaimData.SuggestList")
	if len(riw) != 2 || riw[0]["CommentSuggest"] != "Accepted by UJI Penyetuju Satu" ||
		riw[1]["CommentSuggest"] != "Accepted by "+tiruan.PenyetujuUji(2) {
		t.Fatalf("riwayat klaim S8-S10: %+v", riw)
	}
	if k := u.g.Kasus[tiruan.KomiteUji]; k.StatusWork != models.StatusSelesai || k.AcceptStatus != "1" ||
		u.g.Posisi[tiruan.KomiteUji] != "" {
		t.Fatalf("kasus komite: %+v", k)
	}
	if len(u.g.Efek) != 0 {
		t.Fatalf("di luar produksi nol efek keluar: %v", u.g.Efek)
	}
	if d := u.daftarKerja(tiruan.PenyetujuUji(2)); len(d) != 0 {
		t.Fatalf("kasus selesai keluar dari daftar kerja: %+v", d)
	}
}

// Subjectivity di tingkat 1: isian disimpan di header kasus komite dan dipakai tingkat akhir - tanpa nomor, tanpa OS /
// CLAIMXOL2, IsKomite 0, OS subjectivity STS 1. Tangga satu dan dua tingkat.
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
	if h := u.putus(tiruan.PenyetujuUji(2), models.Keputusan{AcceptStatus: "1", Comment: "UJI 2"},
		http.StatusOK); !h.Selesai || h.AcceptedNo != "" {
		t.Fatalf("tingkat akhir subjectivity: %+v", h)
	}
	if b := u.adj(); b["IsSubjectivity"] != "true" || b["SubjectivityNote"] != "3" || b["IsKomite"] != "0" {
		t.Fatalf("akseptasi subjectivity tangga dua tingkat: %+v", b)
	}
	if len(u.g.OS) != 0 || len(u.g.XOL2) != 0 || u.g.OSSubj[tiruan.KlaimUji].StsSubjectivity != "1" {
		t.Fatalf("subjectivity: nol OS / CLAIMXOL2, OS subjectivity 1: %+v", u.g.OSSubj)
	}
	if _, ada := u.g.Klaim.Nilai(tiruan.KlaimUji)["CNPStatusCase"]; ada {
		t.Fatal("subjectivity melewati S14: CNPStatusCase tidak disentuh")
	}
	u = siap(t, 1, false)
	h := u.putus(tiruan.PenyetujuUji(1), kep, http.StatusOK)
	if !h.Selesai || h.AcceptedNo != "" {
		t.Fatalf("subjectivity: %+v", h)
	}
	b := u.adj()
	if b["IsSubjectivity"] != "true" || b["SubjectivityNote"] != "3" || b["IsKomite"] != "0" || b["AcceptedNo"] != "" {
		t.Fatalf("akseptasi subjectivity: %+v", b)
	}
	if len(u.g.OS) != 0 || u.g.OSSubj[tiruan.KlaimUji].StsSubjectivity != "1" {
		t.Fatal("subjectivity: nol OS, OS subjectivity STS 1")
	}
	kep.SubjectivityNote = ""
	u = siap(t, 1, false)
	if w := u.minta("POST", "/kasus/"+tiruan.KomiteUji+"/putuskan", tiruan.PenyetujuUji(1), kep); w.Code !=
		http.StatusUnprocessableEntity || !strings.Contains(w.Body.String(), "Subjectivity Note") {
		t.Fatalf("catatan subjectivity wajib: %d %s", w.Code, w.Body.String())
	}
}

// Tolak di tingkat 1: S11.3 baris berjalan dan sesudahnya ditolak berkomentar, S11.4 AcceptanceStatus 2, S12 -> EXIT
// (tanpa nomor / OS / email), S20 CLAIM REJECTED, kasus selesai.
func TestTolakTingkatSatu(t *testing.T) {
	u := siap(t, 2, false)
	h := u.putus(tiruan.PenyetujuUji(1), models.Keputusan{AcceptStatus: "2", Comment: "UJI tolak",
		UsulTutup: true}, http.StatusOK)
	if !h.Selesai || h.AcceptedNo != "" || h.KomiteCount != 3 {
		t.Fatalf("tolak: %+v", h)
	}
	t2 := u.g.Tangga[tiruan.KomiteUji]
	if t2[0].Keputusan != "2" || t2[0].Komentar != "UJI tolak" || t2[1].Keputusan != "2" ||
		t2[1].Komentar != "UJI tolak" || t2[1].Tanggal == "" || t2[1].OperatorID != tiruan.PenyetujuUji(2) {
		t.Fatalf("tangga sesudah tolak: %+v", t2)
	}
	b := u.adj()
	if b["AcceptanceStatus"] != "2" || b["AcceptedNo"] != "" {
		t.Fatalf("akseptasi ditolak: %+v", b)
	}
	n := u.g.Klaim.Nilai(tiruan.KlaimUji)
	if n["CNPStatusCase"] != "CLAIM REJECTED" {
		t.Fatalf("header ditolak: %+v", n)
	}
	if _, ada := n["IsCloseFile"]; ada {
		t.Fatal("S14.12 hanya di tingkat akhir disetujui")
	}
	if len(u.g.Riwayat) != 1 || u.g.Riwayat[0].Status != "REJECT" || len(u.g.OS) != 0 || len(u.g.OSSubj) != 0 {
		t.Fatalf("riwayat tolak: %+v OS %d", u.g.Riwayat, len(u.g.OS))
	}
	if riw := u.g.Klaim.Daftar(tiruan.KlaimUji, "ClaimData.SuggestList"); len(riw) != 1 ||
		riw[0]["CommentSuggest"] != "Rejected by UJI Penyetuju Satu" {
		t.Fatalf("riwayat klaim S9: %+v", riw)
	}
	if k := u.g.Kasus[tiruan.KomiteUji]; k.UsulTutup != "1" || k.StatusWork != models.StatusSelesai {
		t.Fatalf("kepala kasus: %+v", k)
	}
}

// Penyetuju bukan pemilik tingkat berjalan ditolak (403) tanpa satu tulisan pun.
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

// Kasus tertutup -> 409 (juga bila klaim induknya yang tertutup); transaksi batal utuh.
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
	if w := u.minta("GET", "/kasus/KMTNP-TIDAKADA", tiruan.PenyetujuUji(1), nil); w.Code != http.StatusNotFound {
		t.Fatalf("kasus tak dikenal: %d", w.Code)
	}
	if w := u.minta("GET", "/kasus/"+tiruan.KlaimUji, tiruan.PenyetujuUji(1), nil); w.Code != http.StatusNotFound {
		t.Fatalf("bukan KMTNP-: %d", w.Code)
	}
	w := u.minta("GET", "/kasus/"+tiruan.KomiteUji, tiruan.PenyetujuUji(1), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("buka kasus: %d %s", w.Code, w.Body.String())
	}
	var ly models.Layar
	if err := json.Unmarshal(w.Body.Bytes(), &ly); err != nil {
		t.Fatal(err)
	}
	if !ly.BolehKerja || len(ly.Judul) != 1 || ly.Judul[0] != "CLAIM COMMITTEE -" || !ly.Isian.Terbuka {
		t.Fatalf("layar: %+v", ly)
	}
	for _, s := range []string{`"Committe Accept Status"`, `"XOL Allocation"`, `"Loss Allocation"`, `"Spreading In"`,
		`"Spreading Out"`, `"Claim Analysis"`, `"keputusan":"Waiting"`, `{"nilai":"1","label":"Treaty Leader Approval"}`,
		`"UJI KRONOLOGI"`, `"UJI CATATAN KOMITE"`} {
		if !strings.Contains(w.Body.String(), s) {
			t.Errorf("layar ShowTransfer tanpa %s", s)
		}
	}
	w = u.minta("GET", "/kasus/"+tiruan.KomiteUji, tiruan.PenyetujuUji(2), nil)
	_ = json.Unmarshal(w.Body.Bytes(), &ly)
	for _, tb := range ly.Tombol {
		if tb.Label == "Submit" && tb.Aktif {
			t.Fatal("Submit nonaktif bagi bukan pemegang")
		}
		if tb.Label == "View Claim" && !tb.Aktif {
			t.Fatal("View Claim aktif juga bagi bukan pemegang (hanya-baca, ViewDtlClaimKmt)")
		}
	}
}

// Produksi: konversi (S14.22), Kasir (S19.3, DirectToKasir), dan email (S19.2) diantre di transaksi yang sama.
func TestEfekKeluarDiProduksi(t *testing.T) {
	u := siap(t, 1, true)
	u.a.Bank["UJI BANK|UJI CABANG|123456"] = "UJI-BANK-1"
	u.a.Email["UJI-CED"] = "uji@contoh.invalid"
	u.a.Konversi["UJIA123102026TX00001"] = "1" // S3 getStatusKonversi_Act (IsPEGAPROD): sudah dikonversi
	h := u.putus(tiruan.PenyetujuUji(1), setuju("UJI"), http.StatusOK)
	jenis := map[string]string{}
	for _, e := range u.g.Efek {
		p := strings.SplitN(e, ":", 3)
		jenis[p[0]] = p[2]
	}
	for _, j := range []string{services.JenisEfekKonversi, services.JenisEfekKasir, services.JenisEfekEmailKomite} {
		if _, ada := jenis[j]; !ada {
			t.Fatalf("efek produksi tanpa %s: %v", j, u.g.Efek)
		}
	}
	if b := u.adj(); b["IDOfBank"] != "UJI-BANK-1" || b["NoAccount"] != "123456" || len(h.AcceptedNo) != 23 {
		t.Fatalf("Kasir S7 / S12-S13: %+v", b)
	}
	var kasir services.MuatanOutbox
	if err := json.Unmarshal([]byte(jenis[services.JenisEfekKasir]), &kasir); err != nil {
		t.Fatal(err)
	}
	if kasir.Kategori1 != "Kasir" || !strings.Contains(jenis[services.JenisEfekKasir], `"Nett":"75"`) ||
		!strings.Contains(jenis[services.JenisEfekKasir], `"LkuId":"IDR"`) {
		t.Fatalf("muatan Kasir (Nett = TotalClaim - PremiumSpreaded): %s", jenis[services.JenisEfekKasir])
	}
	if !strings.Contains(jenis[services.JenisEfekKonversi], `"STS_REJECT":"1"`) {
		t.Fatalf("konversi STSREJECT 1: %s", jenis[services.JenisEfekKonversi])
	}
	// MUATAN email hanya pengenal (claimlife/015): tanpa alamat email, subjek, maupun nama.
	if m := jenis[services.JenisEfekEmailKomite]; strings.Contains(m, "@") || strings.Contains(m, "UJI Penyetuju") {
		t.Fatalf("MUATAN email memuat alamat / nama: %s", m)
	}
	var m services.MuatanOutbox
	if err := json.Unmarshal([]byte(jenis[services.JenisEfekEmailKomite]), &m); err != nil {
		t.Fatal(err)
	}
	isi := map[string]string{}
	for k, v := range m.Isi.(map[string]any) {
		isi[k], _ = v.(string)
	}
	ctx := context.Background()
	u.a.Nama[tiruan.PembuatUji] = "UJI Pembuat"
	u.a.Surel[tiruan.PembuatUji] = "uji.pembuat.syariah@contoh.invalid"
	sr, err := u.l.SusunEmailKomite(ctx, tiruan.KomiteUji, isi)
	if err != nil {
		t.Fatal(err)
	}
	if sr.Kepada != "uji.pembuat.syariah@contoh.invalid" || sr.Akun != "UJI-SYARIAH" || sr.CC != "uji.cc@contoh.invalid" ||
		!strings.HasPrefix(sr.Subjek, "(Approval) Pengajuan Akseptasi : ") ||
		!strings.Contains(sr.HTML, "Dear <strong>UJI Pembuat</strong>") ||
		!strings.Contains(sr.HTML, "Accepted No: <span style=\"font-weight: normal;\">"+h.AcceptedNo) ||
		!strings.Contains(sr.HTML, "<strong>UJI Penyetuju Satu</strong>") {
		t.Fatalf("email komite: %+v", sr)
	}
	// Pelaksana produksi: isi dirakit, lalu pengiriman SMTP ditahan sampai disetujui manusia.
	pl := services.PelaksanaKomiteClaimNonProp{Lingkungan: inti.Produksi, Penyusun: u.l}
	err = pl.Laksanakan(ctx, nil, outbox.BarisEfekKeluar{Modul: "komiteclaimnonprop", Jenis: services.JenisEfekEmailKomite,
		Muatan: jenis[services.JenisEfekEmailKomite]})
	if !errors.Is(err, outbox.ErrEmailBelumDisetujui) {
		t.Fatalf("pelaksana email: %v", err)
	}
}

// Penyerahan ulang baris subjectivity: `.Comment` awal = komentar anggota pertama putaran pertama
// (CreateChildKomiteCNP_Act); keputusan di baris terakhir (S7); putaran kedua tanpa syarat menerbitkan nomor.
func TestKirimUlangSubjectivity(t *testing.T) {
	u := siap(t, 1, false)
	u.putus(tiruan.PenyetujuUji(1), models.Keputusan{AcceptStatus: "1", Comment: "UJI putaran 1", IsSubjectivity: true,
		SubjectivityNote: "3"}, http.StatusOK)
	if b := u.adj(); b["IsSubjectivity"] != "true" || b["IsKomite"] != "0" || b["AcceptedNo"] != "" {
		t.Fatalf("putaran 1 subjectivity: %+v", b)
	}
	// Claim Non Prop menyerahkan ulang (kasus baru) ke jenjang terbawah.
	u.g.Lahirkan("KMTNP-UJI002", tiruan.KlaimUji, tiruan.AdjUji, tiruan.PembuatUji, "UJI Admin",
		[]models.Anggota{{OperatorID: tiruan.PenyetujuUji(1), Jabatan: "UJI-JABATAN-1"}}, tiruan.SaatUji.Add(time.Hour))
	w := u.minta("GET", "/kasus/KMTNP-UJI002", tiruan.PenyetujuUji(1), nil)
	var ly models.Layar
	if err := json.Unmarshal(w.Body.Bytes(), &ly); err != nil || ly.Isian.Nilai.Comment != "UJI putaran 1" {
		t.Fatalf("komentar awal putaran kedua: %q (%d)", ly.Isian.Nilai.Comment, w.Code)
	}
	r := u.minta("POST", "/kasus/KMTNP-UJI002/putuskan", tiruan.PenyetujuUji(1),
		models.Keputusan{AcceptStatus: "1", Comment: "UJI putaran 2"})
	var h services.HasilKeputusan
	_ = json.Unmarshal(r.Body.Bytes(), &h)
	if r.Code != http.StatusOK || !h.Selesai || h.AcceptedNo == "" {
		t.Fatalf("putaran 2 tanpa syarat menerbitkan nomor: %d %s", r.Code, r.Body.String())
	}
	if b := u.adj(); b["IsSubjectivity"] != "false" || b["AcceptedNo"] != h.AcceptedNo || b["AcceptanceStatus"] != "1" {
		t.Fatalf("akseptasi sesudah putaran 2: %+v", b)
	}
	if tg := u.g.Tangga["KMTNP-UJI002"]; tg[len(tg)-1].Keputusan != "1" || tg[len(tg)-1].Komentar != "UJI putaran 2" {
		t.Fatalf("S7 baris terakhir: %+v", tg)
	}
}

// HitServiceToKasirKMT_Act S9 / S10: panjang AcceptedNo bukan 23 / 24 = Exit Activity - tanpa muatan Kasir dan tanpa
// S12-S13 IDOfBank (S7 NoAccount angka saja sudah berjalan).
func TestKasirKeluarBilaPanjangNomorBukan23Atau24(t *testing.T) {
	u := siap(t, 1, true)
	u.a.Bank["UJI BANK|UJI CABANG|123456"] = "UJI-BANK-1"
	u.g.KodeProduksi = "UJIXYZ" // nomor 26 karakter
	u.a.Konversi["UJIXYZA123102026TX00001"] = "1"
	h := u.putus(tiruan.PenyetujuUji(1), setuju("UJI"), http.StatusOK)
	if len(h.AcceptedNo) == 23 || len(h.AcceptedNo) == 24 {
		t.Fatalf("nomor uji harus di luar 23 / 24: %q", h.AcceptedNo)
	}
	for _, e := range u.g.Efek {
		if strings.HasPrefix(e, services.JenisEfekKasir+":") {
			t.Fatalf("Kasir tidak boleh diantre: %s", e)
		}
	}
	if b := u.adj(); b["IDOfBank"] != "" || b["NoAccount"] != "123456" {
		t.Fatalf("S12-S13 dilewati, S7 tetap: %+v", b)
	}
}
