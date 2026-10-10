package handlers_test

// Uji seam HTTP Komite Claim Fac In (handlers -> services -> tiruan + kontrak palsu Claim Fac In): TT2 tingkat 1 SPV A /
// SPV B, pita SPV B, perluasan tangga KCF-02 (tampil tanpa tulis saat GET, tersimpan saat Submit), tingkat akhir (nomor
// akseptasi, OS, JSON_KLAIM, HISTORY, SUBPROGRESS, tulis balik tiga tingkat), tolak, TT3 Reject, TT4 Close Without
// Payment, wewenang (403 / 409 / 422), efek keluar hanya produksi. Fixture `UJI-`.

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimfacin/backend/handlers"
	"nusantarare/modul/komiteclaimfacin/backend/models"
	"nusantarare/modul/komiteclaimfacin/backend/services"
	"nusantarare/modul/komiteclaimfacin/backend/tiruan"
)

var saatUji = time.Date(2026, 10, 10, 9, 30, 0, 0, models.Jakarta)

type uji struct {
	t   *testing.T
	g   *tiruan.Gudang
	a   *tiruan.Acuan
	srv http.Handler
}

func siap(t *testing.T, produksi bool) *uji {
	t.Helper()
	g, a := tiruan.Baru(), tiruan.AcuanBaru()
	a.Roster = tiruan.RosterUji()
	a.JenisReas["10003"] = "UJI QUOTA SHARE"
	a.Nama["UJI-SPVA"], a.Nama["UJI-HEAD"] = "UJI Supervisor", "UJI Kepala"
	a.Email[tiruan.PembuatUji] = "uji.admin@contoh.invalid"
	a.Anggota[models.WorkbasketDeptHead] = []string{"uji.head@contoh.invalid"}
	l := services.Baru(g, a, g.Klaim, func() time.Time { return saatUji }, produksi).
		DenganKasir(models.KonfigurasiKasir{CompanyName: "UJI-CO", LjtdID: "UJI-LJTD", LdcID: "UJI-LDC"}).
		DenganEmail(models.KonfigurasiEmail{Akun: "UJI-AKUN", AkunSyariah: "UJI-SYARIAH", CC: "uji.cc@contoh.invalid"})
	return &uji{t: t, g: g, a: a, srv: handlers.Router(l, true)}
}

func (u *uji) minta(metode, jalur, pelaku, peran string, badan any) *httptest.ResponseRecorder {
	u.t.Helper()
	var b []byte
	if badan != nil {
		b, _ = json.Marshal(badan)
	}
	r := httptest.NewRequest(metode, handlers.Prefix+jalur, bytes.NewReader(b))
	r.Header.Set("X-Pelaku", pelaku)
	r.Header.Set("X-Peran", peran)
	w := httptest.NewRecorder()
	u.srv.ServeHTTP(w, r)
	return w
}

func (u *uji) putus(pelaku, peran string, kep models.Keputusan, mau int) services.HasilKeputusan {
	u.t.Helper()
	w := u.minta("POST", "/kasus/"+tiruan.KomiteUji+"/putuskan", pelaku, peran, kep)
	if w.Code != mau {
		u.t.Fatalf("putuskan %s: %d (mau %d) %s", pelaku, w.Code, mau, w.Body.String())
	}
	var h services.HasilKeputusan
	_ = json.Unmarshal(w.Body.Bytes(), &h)
	return h
}

func (u *uji) layar(pelaku, peran string) models.Layar {
	u.t.Helper()
	w := u.minta("GET", "/kasus/"+tiruan.KomiteUji, pelaku, peran, nil)
	if w.Code != http.StatusOK {
		u.t.Fatalf("buka kasus: %d %s", w.Code, w.Body.String())
	}
	var ly models.Layar
	if err := json.Unmarshal(w.Body.Bytes(), &ly); err != nil {
		u.t.Fatal(err)
	}
	return ly
}

func (u *uji) adj() map[string]string {
	rows := u.g.Klaim.Daftar(tiruan.KlaimUji, models.DaftarAdj(1, 2))
	return rows[0]
}

func bagian(ly models.Layar, kunci string) (models.Bagian, bool) {
	for _, b := range ly.Bagian {
		if b.Kunci == kunci {
			return b, true
		}
	}
	return models.Bagian{}, false
}

var setuju = models.Keputusan{AcceptStatus: models.KeputusanSetuju, Comment: "UJI SETUJU"}

func TestTT2SPVATingkatAkhirTanpaPerluasan(t *testing.T) {
	u := siap(t, false)
	u.g.SiapkanTT2("20000000", saatUji)
	ly := u.layar("UJI-SPVA", models.WorkbasketSPVA)
	if !ly.BolehKerja || strings.Join(ly.Judul, " ") != "CLAIM COMMITTEE - ADJUSTMENT" {
		t.Fatalf("layar: boleh %v judul %v", ly.BolehKerja, ly.Judul)
	}
	tg, _ := bagian(ly, "tangga")
	if len(tg.Grid) != 1 || len(tg.Grid[0].Baris) != 1 {
		t.Fatalf("tangga 20 jt mau satu tingkat: %+v", tg)
	}
	// ubin ringkasan (pola Komite Non Prop): "Claim No" = NoClaim, ID klaim di "Claim ID" (sejajar tabel Committee inbox)
	var ubin []string
	for _, m := range ly.Ubin {
		ubin = append(ubin, m.Label+"="+m.Nilai)
	}
	if strings.Join(ubin, "|") != "Committee No="+tiruan.KomiteUji+"|Claim ID="+tiruan.KlaimUji+
		"|Claim No=UJI-K-0001|Level=1 / 1" {
		t.Fatalf("ubin: %v", ubin)
	}
	h := u.putus("UJI-SPVA", models.WorkbasketSPVA, setuju, http.StatusOK)
	// prompt §5 butir 4: PDF akseptasi diterbitkan bila nomor terisi (berkasnya OQ-KCFI-01 -> pesan info)
	if !h.Selesai || h.AcceptedNo != "UJI-A12.10.2026.00001" || h.KomiteCount != 2 || h.Info != services.OQDokumenPDF {
		t.Fatalf("hasil: %+v", h)
	}
	if !models.PanjangNoAksepCLM(h.AcceptedNo) {
		t.Fatalf("nomor %q bukan 21/22 karakter (gerbang kasir CLM)", h.AcceptedNo)
	}
	b := u.adj()
	if b["AcceptanceStatus"] != "1" || b["AcceptedNo"] != h.AcceptedNo || b["Notes"] != "UJI SETUJU" ||
		b["IsApproved"] != "1" || b["IsPrintAccept"] != "1" || b["IsFacRetro"] != "0" {
		t.Fatalf("adjustment sesudah: %+v", b)
	}
	if ob := u.g.Klaim.Daftar(tiruan.KlaimUji, models.DaftarObjek)[0]; ob["DLAStatus"] != "0" || ob["IsFacretro"] != "0" {
		t.Fatalf("objek sesudah SaveAcceptation_KMT: %+v", ob)
	}
	if v := u.g.Klaim.Nilai(tiruan.KlaimUji); v["AktifButton"] != "0" || v["ClaimData.IsCloseFile"] != "0" {
		t.Fatalf("kepala klaim: AktifButton %q IsCloseFile %q", v["AktifButton"], v["ClaimData.IsCloseFile"])
	}
	kr := u.g.Klaim.Riwayat(tiruan.KlaimUji)
	if len(kr) != 1 || kr[0].Teks != "Accepted by Claim Supervisor - "+tiruan.KomiteUji || kr[0].Tingkat != "Claim Supervisor" {
		t.Fatalf("kronologi: %+v", kr)
	}
	if len(u.g.OS) != 1 || u.g.OS[0].StsReject != "4" || u.g.OS[0].StsDLA != "7" || u.g.OS[0].CaseID != tiruan.KlaimUji {
		t.Fatalf("OS: %+v", u.g.OS)
	}
	for _, s := range []string{`"KomiteNo":"` + tiruan.KomiteUji + `"`, `"CurrencyList":[ `, `"AcceptationList":[ `,
		`"AcceptedNo":"` + h.AcceptedNo + `"`, `"ObjectItemName":"UJI ITEM 2"`} {
		if !strings.Contains(u.g.OS[0].DataJSON, s) {
			t.Fatalf("DATA_JSON tanpa %s:\n%s", s, u.g.OS[0].DataJSON)
		}
	}
	if _, ada := u.g.JSONKlaim[tiruan.KlaimUji]; !ada {
		t.Fatal("JSON_KLAIM tidak ditulis")
	}
	if len(u.g.Riwayat) != 1 || u.g.Riwayat[0].Status != "ACCEPT" || u.g.Riwayat[0].Username != "UJI Supervisor" ||
		u.g.Riwayat[0].IDKomite != tiruan.KomiteUji || u.g.Riwayat[0].Workbasket != "KLAIM" {
		t.Fatalf("HISTORYAKSEPTASIPEGA: %+v", u.g.Riwayat)
	}
	if u.g.SubProgres[tiruan.KomiteUji] != "Accepted" {
		t.Fatalf("SUBPROGRESSCLAIM: %q", u.g.SubProgres[tiruan.KomiteUji])
	}
	if len(u.g.Log) != 1 || u.g.Log[0].JenisService != "AKSEPATSI" || u.g.Log[0].Parameter != tiruan.KlaimUji+" / UJI-RNM-F.001 / 1" {
		t.Fatalf("MONITORING_KLAIM_LOG: %+v", u.g.Log)
	}
	if k := u.g.Kasus[tiruan.KomiteUji]; !k.Tertutup() || u.g.Posisi[tiruan.KomiteUji] != "" {
		t.Fatalf("kasus komite: %+v posisi %q", k, u.g.Posisi[tiruan.KomiteUji])
	}
	if len(u.g.Efek) != 0 {
		t.Fatalf("efek keluar di luar produksi: %+v", u.g.Efek)
	}
	// Submit kedua -> 409
	u.putus("UJI-SPVA", models.WorkbasketSPVA, setuju, http.StatusConflict)
}

func TestTT2PerluasanTanggaTampilLaluTersimpan(t *testing.T) {
	u := siap(t, false)
	u.g.SiapkanTT2("100000000", saatUji)
	ly := u.layar("UJI-SPVA", models.WorkbasketSPVA)
	tg, _ := bagian(ly, "tangga")
	if n := len(tg.Grid[0].Baris); n != 2 || tg.Grid[0].Baris[1]["jabatan"] != "Claim Dept. Head" {
		t.Fatalf("tangga tampil: %+v", tg.Grid[0].Baris)
	}
	if n := len(u.g.Tangga[tiruan.KomiteUji]); n != 1 {
		t.Fatalf("GET menulis tangga: %d baris", n)
	}
	h := u.putus("UJI-SPVA", models.WorkbasketSPVA, setuju, http.StatusOK)
	if h.Selesai || h.AcceptedNo != "" || h.KomiteCount != 2 {
		t.Fatalf("tingkat 1 dari 2: %+v", h)
	}
	k := u.g.Kasus[tiruan.KomiteUji]
	if k.Loop != 2 || len(u.g.Tangga[tiruan.KomiteUji]) != 2 || u.g.Posisi[tiruan.KomiteUji] != models.WorkbasketDeptHead {
		t.Fatalf("perluasan tersimpan: loop %d tangga %+v posisi %q", k.Loop, u.g.Tangga[tiruan.KomiteUji],
			u.g.Posisi[tiruan.KomiteUji])
	}
	if b := u.adj(); b["AcceptanceStatus"] != "" || b["AcceptedNo"] != "" {
		t.Fatalf("adjustment diputus di tingkat 1: %+v", b)
	}
	// SPV A tidak memegang tingkat 2; Dept Head memegangnya
	u.putus("UJI-SPVA", models.WorkbasketSPVA, setuju, http.StatusForbidden)
	if d := u.daftarKerja("UJI-HEAD", models.WorkbasketDeptHead); len(d) != 1 || d[0].Tingkat != 2 || d[0].Jenis != "ADJUSTMENT" {
		t.Fatalf("daftar kerja Dept Head: %+v", d)
	}
	h = u.putus("UJI-HEAD", models.WorkbasketDeptHead, setuju, http.StatusOK)
	if !h.Selesai || h.AcceptedNo == "" {
		t.Fatalf("tingkat akhir: %+v", h)
	}
	kr := u.g.Klaim.Riwayat(tiruan.KlaimUji)
	if len(kr) != 2 || kr[1].Teks != "Accepted by Claim Dept. Head - "+tiruan.KomiteUji {
		t.Fatalf("kronologi: %+v", kr)
	}
	if t2 := u.g.Tangga[tiruan.KomiteUji]; t2[0].OperatorID != "UJI-SPVA" || t2[1].OperatorID != "UJI-HEAD" {
		t.Fatalf("pemutus tercatat di tangga: %+v", t2)
	}
}

func TestTT2PitaSPVB(t *testing.T) {
	// 40 jt: SPV A memutus sendiri; SPV B (cadangan) memicu satu jenjang atas (ApprovalKomite_Act S5, KCF-01).
	u := siap(t, false)
	u.g.SiapkanTT2("40000000", saatUji)
	if d := u.daftarKerja("UJI-SPVB", models.WorkbasketSPVB); len(d) != 1 {
		t.Fatalf("SPV B tidak melihat tingkat 1: %+v", d)
	}
	ly := u.layar("UJI-SPVB", models.WorkbasketSPVB)
	if tg, _ := bagian(ly, "tangga"); len(tg.Grid[0].Baris) != 2 {
		t.Fatalf("pita SPV B mau dua tingkat: %+v", tg.Grid[0].Baris)
	}
	if ly := u.layar("UJI-SPVA", models.WorkbasketSPVA); len(mustBagian(t, ly, "tangga").Grid[0].Baris) != 1 {
		t.Fatal("SPV A pada 40 jt tidak boleh diperluas")
	}
	h := u.putus("UJI-SPVB", models.WorkbasketSPVB, setuju, http.StatusOK)
	if h.Selesai || u.g.Posisi[tiruan.KomiteUji] != models.WorkbasketDeptHead {
		t.Fatalf("SPV B: %+v posisi %q", h, u.g.Posisi[tiruan.KomiteUji])
	}
	// di luar pita (25 jt) SPV B tidak memperluas
	u2 := siap(t, false)
	u2.g.SiapkanTT2("25000000", saatUji)
	if h := u2.putus("UJI-SPVB", models.WorkbasketSPVB, setuju, http.StatusOK); !h.Selesai {
		t.Fatalf("SPV B 25 jt mau selesai di tingkat 1: %+v", h)
	}
}

func mustBagian(t *testing.T, ly models.Layar, kunci string) models.Bagian {
	t.Helper()
	b, ada := bagian(ly, kunci)
	if !ada {
		t.Fatalf("bagian %q tidak ada", kunci)
	}
	return b
}

func (u *uji) daftarKerja(pelaku, peran string) []models.BarisKerja {
	u.t.Helper()
	w := u.minta("GET", "/kasus", pelaku, peran, nil)
	if w.Code != http.StatusOK {
		u.t.Fatalf("daftar kerja %s: %d %s", pelaku, w.Code, w.Body.String())
	}
	var out []models.BarisKerja
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		u.t.Fatal(err)
	}
	return out
}

func TestTT2TolakMenutupSisaTangga(t *testing.T) {
	u := siap(t, false)
	u.g.SiapkanTT2("100000000", saatUji)
	h := u.putus("UJI-SPVA", models.WorkbasketSPVA, models.Keputusan{AcceptStatus: models.KeputusanTolak,
		Comment: "UJI TOLAK"}, http.StatusOK)
	if !h.Selesai || h.AcceptedNo != "" {
		t.Fatalf("tolak: %+v", h)
	}
	tg := u.g.Tangga[tiruan.KomiteUji]
	if len(tg) != 2 || tg[0].Keputusan != "2" || tg[0].Komentar != "UJI TOLAK" || tg[1].Keputusan != "2" || tg[1].Komentar != "" {
		t.Fatalf("tangga sesudah tolak: %+v", tg)
	}
	if b := u.adj(); b["AcceptanceStatus"] != "2" || b["IsApproved"] != "" || b["AcceptedNo"] != "" {
		t.Fatalf("adjustment ditolak: %+v", b)
	}
	if len(u.g.OS) != 0 || len(u.g.Log) != 0 {
		t.Fatalf("tolak menulis OS / log: %+v %+v", u.g.OS, u.g.Log)
	}
	if len(u.g.Riwayat) != 1 || u.g.Riwayat[0].Status != "REJECT" || u.g.SubProgres[tiruan.KomiteUji] != "Rejected" {
		t.Fatalf("riwayat / sub-progres: %+v %q", u.g.Riwayat, u.g.SubProgres[tiruan.KomiteUji])
	}
	if kr := u.g.Klaim.Riwayat(tiruan.KlaimUji); len(kr) != 1 || kr[0].Teks != "Rejected by Claim Supervisor - "+tiruan.KomiteUji {
		t.Fatalf("kronologi: %+v", kr)
	}
}

func TestWewenangDanValidasi(t *testing.T) {
	u := siap(t, false)
	u.g.SiapkanTT2("20000000", saatUji)
	u.putus("UJI-LAIN", "", setuju, http.StatusForbidden)
	u.putus("UJI-HEAD", models.WorkbasketDeptHead, setuju, http.StatusForbidden) // bukan tingkat berjalan
	w := u.minta("POST", "/kasus/"+tiruan.KomiteUji+"/putuskan", "UJI-SPVA", models.WorkbasketSPVA,
		models.Keputusan{AcceptStatus: models.KeputusanSetuju})
	if w.Code != http.StatusUnprocessableEntity || !strings.Contains(w.Body.String(), "Note: Value cannot be blank") {
		t.Fatalf("Note kosong: %d %s", w.Code, w.Body.String())
	}
	w = u.minta("POST", "/kasus/"+tiruan.KomiteUji+"/putuskan", "UJI-SPVA", models.WorkbasketSPVA,
		models.Keputusan{Comment: "UJI"})
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("AcceptStatus kosong: %d %s", w.Code, w.Body.String())
	}
	if w := u.minta("GET", "/kasus/KMTLF-000001", "UJI-SPVA", "", nil); w.Code != http.StatusNotFound {
		t.Fatalf("kasus lini lain: %d", w.Code)
	}
	if len(u.g.Riwayat) != 0 || u.g.Kasus[tiruan.KomiteUji].Count != 1 {
		t.Fatal("penolakan wewenang / validasi menulis")
	}
}

func TestTT3RejectDisetujuiMenutupKlaim(t *testing.T) {
	u := siap(t, false)
	u.g.SiapkanTutup(models.TransferReject, saatUji)
	ly := u.layar("UJI-HEAD", models.WorkbasketDeptHead)
	if strings.Join(ly.Judul, " ") != "CLAIM COMMITTEE - REJECT" || ly.Isian.TampilUsul {
		t.Fatalf("layar TT3: %+v", ly.Judul)
	}
	if _, ada := bagian(ly, "tangga"); ada {
		t.Fatal("List of Committee hanya TT2 (LS41)")
	}
	// Remarks = ClaimData.Remark klaim; tiga teks pop-up lain dari kepala kasus (urutan: TestTT3TT4TeksPopUpTampil)
	if tk := mustBagian(t, ly, "teksKomite"); len(tk.Medan) != 4 || tk.Medan[3].Nilai != "UJI ALASAN" {
		t.Fatalf("teks komite TT3: %+v", tk)
	}
	u.putus("UJI-SPVA", models.WorkbasketSPVA, setuju, http.StatusForbidden)
	h := u.putus("UJI-HEAD", models.WorkbasketDeptHead, setuju, http.StatusOK)
	if !h.Selesai || h.Info == "" {
		t.Fatalf("hasil TT3: %+v", h)
	}
	if st := u.g.Klaim.Status(tiruan.KlaimUji); st != kontrak.StatusKlaimDitolak {
		t.Fatalf("klaim induk %q", st)
	}
	if len(u.g.OS) != 1 || u.g.OS[0].StsReject != "2" || u.g.OS[0].StsDLA != "8" || !strings.Contains(u.g.OS[0].DataJSON,
		`"Value":"1000000"`) {
		t.Fatalf("OS STS 2 (estimasi PrintFaceClaim 1 saja): %+v", u.g.OS)
	}
	if len(u.g.KlaimTolak) != 1 || u.g.KlaimTolak[0].ID != tiruan.KlaimUji || u.g.KlaimTolak[0].Remark != "UJI ALASAN" ||
		u.g.KlaimTolak[0].PembuatID != tiruan.PembuatUji {
		t.Fatalf("CLAIMREJECTED: %+v", u.g.KlaimTolak)
	}
	if kr := u.g.Klaim.Riwayat(tiruan.KlaimUji); len(kr) != 1 || kr[0].Teks != "Accepted by Claim Dept. Head - "+tiruan.KomiteUji {
		t.Fatalf("kronologi TT3: %+v", kr)
	}
	if len(u.g.Riwayat) != 0 {
		t.Fatalf("TT3 tanpa HISTORYAKSEPTASIPEGA: %+v", u.g.Riwayat)
	}
}

func TestTT3DitolakKlaimTetapTerbuka(t *testing.T) {
	u := siap(t, false)
	u.g.SiapkanTutup(models.TransferReject, saatUji)
	h := u.putus("UJI-HEAD", models.WorkbasketDeptHead, models.Keputusan{AcceptStatus: models.KeputusanTolak,
		Comment: "UJI"}, http.StatusOK)
	if !h.Selesai || u.g.Klaim.Status(tiruan.KlaimUji) != "" || len(u.g.OS) != 0 || len(u.g.KlaimTolak) != 0 {
		t.Fatalf("TT3 ditolak: %+v status %q OS %d", h, u.g.Klaim.Status(tiruan.KlaimUji), len(u.g.OS))
	}
	if _, ada := u.g.JSONKlaim[tiruan.KlaimUji]; !ada { // KomitePost_Reject S11 tanpa syarat
		t.Fatal("JSON_KLAIM TT3 tanpa syarat")
	}
}

func TestTT4CloseWithoutPayment(t *testing.T) {
	u := siap(t, false)
	u.g.SiapkanTutup(models.TransferClose, saatUji)
	if ly := u.layar("UJI-HEAD", models.WorkbasketDeptHead); strings.Join(ly.Judul, " ") != "CLAIM COMMITTEE - CLOSE" {
		t.Fatalf("judul TT4: %v", ly.Judul)
	}
	u.putus("UJI-HEAD", models.WorkbasketDeptHead, setuju, http.StatusOK)
	if st := u.g.Klaim.Status(tiruan.KlaimUji); st != kontrak.StatusKlaimSelesai {
		t.Fatalf("klaim induk %q", st)
	}
	if len(u.g.OS) != 1 || u.g.OS[0].StsReject != "4" || u.g.OS[0].StsDLA != "" ||
		u.g.OS[0].DataJSON != "{\n\"CauseOfLoss\":\"UJI KEBAKARAN\"\n,\"CauseOfLossID\":\"UJI-COL\"\n,\"NoClaim\":\"UJI-K-0001\"\n,\"pxObjClass\":\"ASM-FW-GCNMFW-Data-osAkseptasi\"\n}\n" {
		t.Fatalf("OS TT4: %+v", u.g.OS)
	}
	if len(u.g.KlaimTolak) != 0 {
		t.Fatal("TT4 tidak menulis CLAIMREJECTED")
	}
}

func TestKlaimTertutupMenolakSubmitUtuh(t *testing.T) {
	u := siap(t, false)
	u.g.SiapkanTutup(models.TransferReject, saatUji)
	if err := u.g.Klaim.TutupKlaimFacIn(context.Background(), nil, tiruan.KlaimUji, "", kontrak.StatusKlaimSelesai, saatUji); err != nil {
		t.Fatal(err)
	}
	u.putus("UJI-HEAD", models.WorkbasketDeptHead, setuju, http.StatusConflict)
	if u.g.Tangga[tiruan.KomiteUji][0].Keputusan != models.KeputusanMenunggu || len(u.g.OS) != 0 {
		t.Fatal("Submit gagal meninggalkan tulisan")
	}
}

func TestEfekKeluarHanyaProduksi(t *testing.T) {
	u := siap(t, true)
	u.g.SiapkanTT2("20000000", saatUji)
	u.a.Bank["UJI BANK|UJI CABANG|12-34"] = "UJI-BANK-01"
	u.a.Konversi["UJI-A1210202600001"] = "1" // getStatusKonversi_Act: nomor tanpa titik sudah di TRLOSS_DETAIL_T
	u.putus("UJI-SPVA", models.WorkbasketSPVA, setuju, http.StatusOK)
	jenis := map[string]int{}
	for _, e := range u.g.Efek {
		jenis[e.Jenis]++
		if strings.Contains(e.Muatan, "@") {
			t.Fatalf("MUATAN memuat alamat email: %s", e.Muatan)
		}
	}
	if jenis[services.JenisEfekKonversi] != 1 || jenis[services.JenisEfekKasir] != 1 || jenis[services.JenisEfekEmailKomite] != 1 {
		t.Fatalf("efek produksi: %+v", u.g.Efek)
	}
	if b := u.adj(); b["IDOfBank"] != "UJI-BANK-01" {
		t.Fatalf("IDOfBank S12-S13: %+v", b)
	}
	for _, e := range u.g.Efek {
		if e.Jenis == services.JenisEfekKasir && !strings.Contains(e.Muatan, `"AccountNo":"1234"`) {
			t.Fatalf("muatan kasir: %s", e.Muatan)
		}
	}
}

func TestKasirMenungguStatusKonversi(t *testing.T) {
	// HitServiceToKasirKMT_Act S3: transisi pasca `.StatusKonversi=="1"` (getStatusKonversi_Act - COUNT TRLOSS_DETAIL_T,
	// hanya produksi), selainnya keluar SEBELUM S12-S13 IDOfBank: kasir tidak diantre, IDOfBank tidak ditulis balik.
	u := siap(t, true)
	u.g.SiapkanTT2("20000000", saatUji)
	u.a.Bank["UJI BANK|UJI CABANG|12-34"] = "UJI-BANK-01"
	u.putus("UJI-SPVA", models.WorkbasketSPVA, setuju, http.StatusOK)
	for _, e := range u.g.Efek {
		if e.Jenis == services.JenisEfekKasir {
			t.Fatalf("kasir diantre tanpa status konversi: %+v", e)
		}
	}
	if b := u.adj(); b["IDOfBank"] != "" {
		t.Fatalf("IDOfBank ditulis padahal S3 keluar: %q", b["IDOfBank"])
	}
}

func TestPitaSPVBTigaKasusPemutus(t *testing.T) {
	// Jawaban work owner 10-10-2026 (OQ-KCFI-06): pita SPV B (ApprovalKomite_Act S5) hanya untuk pemutus tingkat 1 yang
	// anggota ReasClaimSPVB SAJA; anggota ReasClaimSPVA + ReasClaimSPVB sekaligus = SPV A (tangga tidak diperluas).
	kasus := []struct {
		nama, peran string
		tingkat     int
		selesai     bool
	}{
		{"SPVA saja", models.WorkbasketSPVA, 1, true},
		{"SPVB saja", models.WorkbasketSPVB, 2, false},
		{"keduanya", models.WorkbasketSPVA + "," + models.WorkbasketSPVB, 1, true},
	}
	for _, c := range kasus {
		t.Run(c.nama, func(t *testing.T) {
			u := siap(t, false)
			u.g.SiapkanTT2("40000000", saatUji)
			if n := len(mustBagian(t, u.layar("UJI-SPV", c.peran), "tangga").Grid[0].Baris); n != c.tingkat {
				t.Fatalf("tangga tampil %d tingkat, mau %d", n, c.tingkat)
			}
			if h := u.putus("UJI-SPV", c.peran, setuju, http.StatusOK); h.Selesai != c.selesai {
				t.Fatalf("Submit: %+v, selesai mau %v", h, c.selesai)
			}
			if n := len(u.g.Tangga[tiruan.KomiteUji]); n != c.tingkat {
				t.Fatalf("tangga tersimpan %d tingkat, mau %d", n, c.tingkat)
			}
		})
	}
}

func TestTT3TT4TeksPopUpTampil(t *testing.T) {
	// SendRejectClaimToKomite2 / SendCloseClaimToKomite menyalin Chronology / Extent / Policy Liability pop-up ke
	// `Komite.*`; ShowTransfer LS39 (VIS NOTBLANK) menampilkannya - urutan LS39: Legal, Chronology, Extent, Remarks.
	for _, tt := range []string{models.TransferReject, models.TransferClose} {
		u := siap(t, false)
		u.g.SiapkanTutup(tt, saatUji)
		var label []string
		for _, m := range mustBagian(t, u.layar("UJI-HEAD", models.WorkbasketDeptHead), "teksKomite").Medan {
			label = append(label, m.Label+"="+m.Nilai)
		}
		mau := "Legal Liability / Policy Liability=UJI LIABILITAS|Chronology=UJI KRONOLOGI|Extent Of Loss=UJI EXTENT|" +
			"Remarks=UJI ALASAN"
		if strings.Join(label, "|") != mau {
			t.Fatalf("TT%s teks komite: %v", tt, label)
		}
	}
}
