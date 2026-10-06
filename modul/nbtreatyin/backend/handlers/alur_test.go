package handlers_test

// Uji seam 1 (spec §6.2) - daur hidup berkas, tangga, wewenang, medan wajib,
// dan medan terkunci diuji LEWAT HTTP handler, di atas gudang tiruan
// (`backend/tiruan`) yang membatalkan transaksi sungguh-sungguh. Nilai harapan
// dari XML (INVENTARIS-XML.md bab 4-7) dan AC spec; fixture berawalan UJI-.
//
// ⚠️ Keutuhan transaksi terhadap Oracle sungguhan diuji terpisah di
// `repository/polis_db_test.go` (`-tags db`, spec §6.3) - di sini hanya
// urutan dan pembatalan di lapisan layanan.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbtreatyin/backend/handlers"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/services"
	"nusantarare/modul/nbtreatyin/backend/tiruan"
)

type pelakuUji struct{ akun, peran string }

var (
	admin    = pelakuUji{"UJI-ADMIN", "UJI-LAIN," + models.PosisiAdmin}
	secHead  = pelakuUji{"UJI-SH", models.PosisiSecHead}
	deptHead = pelakuUji{"UJI-DH", models.PosisiDeptHead}
	orang    = pelakuUji{"UJI-ORANG", "UJI-LAIN"}
)

type uji struct {
	t *testing.T
	g *tiruan.Gudang
	s http.Handler
}

func baru(t *testing.T) *uji {
	t.Helper()
	g := tiruan.Baru()
	g.Nama["UJI-ADMIN"] = "Uji Admin"
	g.Nama["UJI-SH"] = "Uji Sec"
	g.Nama["UJI-DH"] = "Uji Dept"
	return &uji{t: t, g: g, s: handlers.Router(services.Baru(g, jam), true)}
}

func (u *uji) panggil(metode, jalur string, p pelakuUji, badan any) (int, string) {
	u.t.Helper()
	var b bytes.Buffer
	if badan != nil {
		if err := json.NewEncoder(&b).Encode(badan); err != nil {
			u.t.Fatal(err)
		}
	}
	r := httptest.NewRequest(metode, "/api/nb-treaty-in"+jalur, &b)
	if p.akun != "" {
		r.Header.Set("X-Pelaku", p.akun)
		r.Header.Set("X-Peran", p.peran)
	}
	w := httptest.NewRecorder()
	u.s.ServeHTTP(w, r)
	return w.Code, w.Body.String()
}

func (u *uji) buat() string {
	u.t.Helper()
	kode, isi := u.panggil("POST", "/kasus", admin, nil)
	if kode != http.StatusCreated {
		u.t.Fatalf("buat: %d %s", kode, isi)
	}
	var k models.Kasus
	if err := json.Unmarshal([]byte(isi), &k); err != nil {
		u.t.Fatal(err)
	}
	return k.ID
}

func (u *uji) kirim(id string, p pelakuUji, h *models.Halaman) (int, string) {
	u.t.Helper()
	return u.panggil("POST", "/kasus/"+id+"/kirim", p, map[string]any{"halaman": h})
}

// halamanLengkap - isian layar admin yang memenuhi seluruh medan wajib.
func halamanLengkap(isApproved string) *models.Halaman {
	h := models.HalamanBaru()
	set := func(m, v string) { h.Setel("PolicyTreatyIn."+m, v) }
	set("StartDate", "2026-10-01")
	set("EndDate", "2026-12-31")
	set("QuotationData.IsSurveyReport", "UJI-NO")
	set("Quartal", "4")
	set("YearOfQuartal", "2026")
	set("QuotationData.MOID", "UJI-MO-1")
	for _, m := range []string{"PremiOgp", "RiCommOgp", "OveriddingCommOgp", "PremiOnp", "RiCommOnp",
		"OveriddingCommOnp", "Claim", "OutstandingClaim", "SalvageValue", "ExcessLoss", "Deduction1", "Deduction2"} {
		set(m, "0")
	}
	set("PremiOgp", "1000")
	// ResultOnp1 = 0: nilai yang CountOGPONP langkah 2 tulis saat PremiOnp "0"
	// (refresh layar admin) - wajib di layar atasan (AC 47).
	set("ResultOnp1", "0")
	set("IsApproved", isApproved)
	set("Suggest", "UJI-catatan")
	return h
}

func putusan(isApproved string) *models.Halaman {
	h := models.HalamanBaru()
	h.Setel("PolicyTreatyIn.IsApproved", isApproved)
	h.Setel("PolicyTreatyIn.Suggest", "UJI-"+isApproved)
	return h
}

func TestBuatKasusDiAntreanAdmin(t *testing.T) {
	u := baru(t)
	if kode, _ := u.panggil("POST", "/kasus", pelakuUji{}, nil); kode != http.StatusUnauthorized {
		t.Fatalf("tanpa identitas: %d", kode)
	}
	id := u.buat()
	if !strings.HasPrefix(id, "NB-") {
		t.Fatalf("pengenal %q, harap berawalan NB-", id)
	}
	k := u.g.Kasus[id]
	if k.PositionNote != models.PosisiAdmin || k.Position != "4" || k.StatusWork != "Input Realitation" {
		t.Fatalf("Start1 -> Assignment2: %+v", k)
	}
	if u.g.Halaman[id].Ambil("Quotation.BusinessFac") != "T" {
		t.Fatal("kasus treaty ber-BusinessFac T (filter E GetListOpportunity)")
	}
}

func TestBukanAnggotaAntreanDitolak(t *testing.T) { // AC 11, 14, 92
	u := baru(t)
	id := u.buat()
	if kode, _ := u.kirim(id, orang, halamanLengkap("1")); kode != http.StatusForbidden {
		t.Fatalf("bukan anggota: %d", kode)
	}
	if _, isi := u.panggil("GET", "/kasus/"+id, orang, nil); !strings.Contains(isi, `"bolehKerja":false`) {
		t.Fatalf("bukan anggota hanya-baca: %s", isi)
	}
	// Keanggotaan menurut NAMA: urutan daftar peran tidak berpengaruh.
	terbalik := pelakuUji{"UJI-ADMIN2", models.PosisiAdmin + ",UJI-LAIN"}
	if _, isi := u.panggil("GET", "/kasus/"+id, terbalik, nil); !strings.Contains(isi, `"bolehKerja":true`) {
		t.Fatalf("anggota dengan urutan peran berbeda: %s", isi)
	}
}

func TestAdminMenolakDiselesaikanDitolak(t *testing.T) { // AC 1, 5, 39-41, 43, 71
	u := baru(t)
	id := u.buat()
	kode, isi := u.kirim(id, admin, halamanLengkap("0"))
	if kode != http.StatusOK || !strings.Contains(isi, models.StatusDitolak) {
		t.Fatalf("admin menolak: %d %s", kode, isi)
	}
	if got := u.g.Halaman[id].Ambil("NBStatus"); got != "NB WAS DECLINED BY  UJI ADMIN" {
		t.Fatalf("NBStatus %q", got)
	}
	if len(u.g.Riwayat) != 1 {
		t.Fatalf("satu baris riwayat per perpindahan, dapat %d", len(u.g.Riwayat))
	}
	r := u.g.Riwayat[0]
	if r.Status != "REJECT" || r.OperatorID != "UJI-ADMIN" || r.Username != "Uji Admin" || r.Workbasket != models.PosisiAdmin {
		t.Fatalf("riwayat %+v", r)
	}
	if r.IDPega != "ASM-FW-GISFW-WORK-NB "+id {
		t.Fatalf("ID_PEGA %q", r.IDPega)
	}
	// Catatan usulan -> POOLDATA.HISTORYAKSEPTASIPRODUCTION (K4; SaveViewSuggest
	// langkah 2, InsertViewSuggest_SQL; spec-penyimpanan ID-31, AC 39-44).
	if len(u.g.Usulan) != 1 {
		t.Fatalf("satu baris riwayat produksi per catatan, dapat %d", len(u.g.Usulan))
	}
	c := u.g.Usulan[0]
	if c.IDPega != "ASM-FW-GISFW-WORK-NB "+id || c.NoUrut != 1 || c.TypePolis != "NB" || c.Posisi != "Policy" ||
		c.PIC != "Uji Admin" || c.AksesLogin != "UJI-ADMIN" || c.Approval != "Reject" || c.Keterangan != "UJI-catatan" ||
		c.Type != "T" || c.Putaran != "2" || c.TglInp != "2026-10-03 09:00:00" || c.Div != "" {
		t.Fatalf("baris riwayat produksi %+v", c)
	}
	// layar membaca catatan kembali dari tabel itu (riwayat catatan, AC 71)
	_, isi = u.panggil("GET", "/kasus/"+id, admin, nil)
	var ly services.Layar
	if err := json.Unmarshal([]byte(isi), &ly); err != nil {
		t.Fatal(err)
	}
	d := ly.Halaman.AmbilDaftar(models.DaftarUsulan)
	if len(d) != 1 || d[0]["OperatorID"] != "UJI-ADMIN" || d[0]["OperatorName"] != "Uji Admin" || d[0]["IsApproved"] != "0" ||
		d[0]["Suggest"] != "UJI-catatan" || d[0]["Date"] != "2026-10-03 09:00:00" {
		t.Fatalf("catatan dibaca kembali %+v", d)
	}
	if kode, _ := u.kirim(id, admin, halamanLengkap("1")); kode != http.StatusConflict {
		t.Fatalf("kasus tertutup harus ditolak: %d", kode)
	}
}

func TestTanggaPenuhDanNomorPolisSekali(t *testing.T) { // AC 6-9, 31, 73, 74
	u := baru(t)
	id := u.buat()
	if kode, isi := u.kirim(id, admin, halamanLengkap("1")); kode != http.StatusOK {
		t.Fatalf("admin menyetujui: %d %s", kode, isi)
	}
	if u.g.Kasus[id].PositionNote != models.PosisiSecHead {
		t.Fatalf("admin menyetujui -> %q", u.g.Kasus[id].PositionNote)
	}
	if got := u.g.Halaman[id].Ambil("NBStatus"); got != "NB IS IN REASTREATYINSECHEAD'S INBOX" {
		t.Fatalf("NBStatus dari data posisi, dapat %q", got)
	}
	// Sec Head menolak -> kembali ke admin, NBStatus kosong
	if kode, isi := u.kirim(id, secHead, putusan("0")); kode != http.StatusOK {
		t.Fatalf("Sec Head menolak: %d %s", kode, isi)
	}
	if u.g.Kasus[id].PositionNote != models.PosisiAdmin || u.g.Halaman[id].Ambil("NBStatus") != "" {
		t.Fatalf("Sec Head menolak: %+v / %q", u.g.Kasus[id], u.g.Halaman[id].Ambil("NBStatus"))
	}
	// naik lagi; Sec Head menyetujui -> Dept Head SELALU (AC 8, WO)
	u.kirim(id, admin, halamanLengkap("1"))
	if kode, isi := u.kirim(id, secHead, putusan("1")); kode != http.StatusOK || u.g.Kasus[id].PositionNote != models.PosisiDeptHead {
		t.Fatalf("Sec Head menyetujui -> %q (%d %s)", u.g.Kasus[id].PositionNote, kode, isi)
	}
	// Sec Head bukan anggota antrean Dept Head
	if kode, _ := u.panggil("POST", "/kasus/"+id+"/nomor-polis", secHead, map[string]any{"halaman": putusan("1")}); kode != http.StatusForbidden {
		t.Fatalf("Sec Head menerbitkan nomor: %d", kode)
	}
	var n1, n2 services.NomorPolis
	for _, n := range []*services.NomorPolis{&n1, &n2} {
		kode, isi := u.panggil("POST", "/kasus/"+id+"/nomor-polis", deptHead, map[string]any{"halaman": putusan("1")})
		if kode != http.StatusOK || json.Unmarshal([]byte(isi), n) != nil {
			t.Fatalf("nomor polis: %d %s", kode, isi)
		}
	}
	if !strings.HasPrefix(n1.PolicyNo, "UJI-QR.TUJI-OJK.10.2026.") || len(n1.PolicyNo) != len("UJI-QR.TUJI-OJK.10.2026.00001") {
		t.Fatalf("bentuk nomor %q", n1.PolicyNo)
	}
	if n2.PolicyNo != n1.PolicyNo {
		t.Fatalf("nomor dibentuk SEKALI: %q lalu %q", n1.PolicyNo, n2.PolicyNo)
	}
	kode, isi := u.kirim(id, deptHead, putusan("1"))
	if kode != http.StatusOK || u.g.Kasus[id].StatusWork != models.StatusSelesai {
		t.Fatalf("Dept Head menyetujui bernomor: %d %s", kode, isi)
	}
	if strings.Contains(isi, "pesanKonversi") {
		t.Fatal("di luar produksi konversi dilewati - bukan gagal")
	}
	if len(u.g.Riwayat) != 5 {
		t.Fatalf("lima submit, lima riwayat - dapat %d", len(u.g.Riwayat))
	}
	// [penyimpangan sadar] K4: catatan SETIAP jenjang ditulis di submit-nya,
	// NOURUT berurutan per kasus (XML: hanya pasca-submit admin).
	var akses []string
	for i, c := range u.g.Usulan {
		if c.NoUrut != i+1 || c.IDPega != "ASM-FW-GISFW-WORK-NB "+id {
			t.Fatalf("NOURUT %d: %+v", i+1, c)
		}
		akses = append(akses, c.AksesLogin+":"+c.Approval)
	}
	if got := strings.Join(akses, " "); got != "UJI-ADMIN:Accept UJI-SH:Reject UJI-ADMIN:Accept UJI-SH:Accept UJI-DH:Accept" {
		t.Fatalf("baris riwayat produksi per jenjang: %s", got)
	}
}

// K4 [penyimpangan sadar — disetujui WO 04-10-2026] (F2): NOURUT diberikan
// gudang (MAX+1 per IDPEGA), XML `InsertViewSuggest_SQL` memakai
// `.pxListSubscript` (SaveViewSuggest 2.1.2 CARI2). Bukti keduanya sama
// sepanjang tangga (tolak-naik-setuju, tiga jenjang): sesudah SETIAP submit,
// SuggestList yang dibangun ulang dari tabel memuat baris NOURUT j tepat di
// pxListSubscript j, dan baris yang baru ditulis = ujung daftar.
func TestNourutUsulanSamaDenganSubskripSuggestList(t *testing.T) { // AC 39-44, K4
	u := baru(t)
	id := u.buat()
	langkah := []struct {
		p pelakuUji
		h *models.Halaman
	}{
		{admin, halamanLengkap("1")}, {secHead, putusan("0")}, {admin, halamanLengkap("1")},
		{secHead, putusan("1")}, {deptHead, putusan("1")},
	}
	for i, l := range langkah {
		l.h.Setel("PolicyTreatyIn.Suggest", fmt.Sprintf("UJI-catatan-%d", i+1))
		if l.p == deptHead {
			if kode, isi := u.panggil("POST", "/kasus/"+id+"/nomor-polis", deptHead, map[string]any{"halaman": l.h}); kode != http.StatusOK {
				t.Fatalf("nomor polis: %d %s", kode, isi)
			}
		}
		if kode, isi := u.kirim(id, l.p, l.h); kode != http.StatusOK {
			t.Fatalf("submit %d: %d %s", i+1, kode, isi)
		}
		_, isi := u.panggil("GET", "/kasus/"+id, admin, nil)
		var ly services.Layar
		if err := json.Unmarshal([]byte(isi), &ly); err != nil {
			t.Fatal(err)
		}
		daftar := ly.Halaman.AmbilDaftar(models.DaftarUsulan)
		if len(daftar) != i+1 || len(u.g.Usulan) != i+1 {
			t.Fatalf("submit %d: %d baris SuggestList, %d baris tabel", i+1, len(daftar), len(u.g.Usulan))
		}
		for j, b := range daftar { // posisi baris (berbasis 1) = j+1
			c := u.g.Usulan[j]
			if c.NoUrut != j+1 || c.Keterangan != b["Suggest"] || c.AksesLogin != b["OperatorID"] {
				t.Fatalf("submit %d: posisi baris (berbasis 1) %d = %+v, baris tabel NOURUT %d = %+v", i+1, j+1, b, c.NoUrut, c)
			}
		}
		if daftar[i]["Suggest"] != fmt.Sprintf("UJI-catatan-%d", i+1) {
			t.Fatalf("submit %d: catatan baru bukan di ujung daftar: %+v", i+1, daftar)
		}
	}
}

func TestMedanWajibMenahanKirimDanSimpan(t *testing.T) { // AC 45, 48
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("1")
	h.Setel("PolicyTreatyIn.QuotationData.MOID", "")
	h.Setel("PolicyTreatyIn.Suggest", "")
	kode, isi := u.kirim(id, admin, h)
	if kode != http.StatusUnprocessableEntity || !strings.Contains(isi, "Marketing Officer") || !strings.Contains(isi, "Suggest") {
		t.Fatalf("kirim: %d %s", kode, isi)
	}
	if kode, _ := u.panggil("PUT", "/kasus/"+id, admin, map[string]any{"halaman": h}); kode != http.StatusUnprocessableEntity {
		t.Fatalf("Save dengan medan wajib kosong (AC 48): %d", kode)
	}
	if len(u.g.Riwayat) != 0 || u.g.Kasus[id].PositionNote != models.PosisiAdmin {
		t.Fatal("validasi gagal tidak boleh menyisakan riwayat atau perpindahan")
	}
	// Approval belum dipilih: tidak ada tombol Submit
	if kode, isi := u.kirim(id, admin, halamanLengkap("")); kode != http.StatusUnprocessableEntity || !strings.Contains(isi, "Approval") {
		t.Fatalf("tanpa Approval: %d %s", kode, isi)
	}
}

func TestPolisSerupaMenahan(t *testing.T) { // TreatyRealizationCheckDuplicate, AC 83
	u := baru(t)
	id := u.buat()
	u.g.Serupa = []string{"UJI-NOPOL-1"}
	kode, isi := u.kirim(id, admin, halamanLengkap("1"))
	if kode != http.StatusUnprocessableEntity || !strings.Contains(isi, "Protect Duplicate Policy; data is similar to UJI-NOPOL-1 ") {
		t.Fatalf("dapat %d %s", kode, isi)
	}
	if len(u.g.Riwayat) != 0 {
		t.Fatal("submit yang tertahan tidak boleh meninggalkan riwayat")
	}
	h := halamanLengkap("1")
	h.Setel("PolicyTreatyIn.ClaimType", "XOL")
	h.Setel("PolicyTreatyIn.ClaimPaymentType", "UJI-Claim")
	if kode, isi := u.kirim(id, admin, h); kode != http.StatusOK {
		t.Fatalf("ClaimType XOL melewati pemeriksaan: %d %s", kode, isi)
	}
}

func TestSatuTransaksiPembatalanUtuh(t *testing.T) { // AC 29, 83 (urutan layanan)
	u := baru(t)
	id := u.buat()
	for _, op := range []string{"PindahPosisi", "CatatRiwayat", "CatatUsulan"} {
		u.g.GagalDi = op
		if kode, _ := u.kirim(id, admin, halamanLengkap("1")); kode != http.StatusInternalServerError {
			t.Fatalf("%s gagal: %d", op, kode)
		}
		if len(u.g.Riwayat) != 0 || len(u.g.Usulan) != 0 || u.g.Halaman[id].Ambil("PolicyTreatyIn.Suggest") != "" || u.g.Kasus[id].PositionNote != models.PosisiAdmin {
			t.Fatalf("%s: kegagalan di tengah harus membatalkan riwayat, halaman, dan perpindahan", op)
		}
	}
}

func TestMedanTerkunciAtasanDanTurunanAdmin(t *testing.T) { // AC 49-52
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("1")
	h.Setel("PolicyTreatyIn.NetPremium", "999999")          // turunan: dihitung ulang server
	h.Setel("PolicyTreatyIn.OJKBusinessID", "UJI-KARANGAN") // milik server: diabaikan
	if kode, isi := u.kirim(id, admin, h); kode != http.StatusOK {
		t.Fatalf("%d %s", kode, isi)
	}
	if got := u.g.Halaman[id].Ambil("PolicyTreatyIn.NetPremium"); got == "999999" {
		t.Fatal("NetPremium kiriman layar tidak boleh tersimpan")
	}
	if u.g.Halaman[id].Ambil("PolicyTreatyIn.OJKBusinessID") != "" {
		t.Fatal("OJKBusinessID milik server")
	}
	a := putusan("1")
	a.Setel("PolicyTreatyIn.PremiOgp", "999999")
	// AC 52: `DetailDeptHeadTreatyIn_UW` - DueTo, FlagPPH, No Offer Slip
	// `pyDisabled=always`; `ListSuggest.ProductionDate` tidak tampil tanpa
	// tempat berperan (tiket 05) -> keempatnya tidak dapat diisi atasan.
	a.Setel("PolicyTreatyIn.QuotationData.NoOfferSlip", "UJI-SLIP")
	a.Setel("PolicyTreatyIn.FlagPPH", "true")
	a.Setel("PolicyTreatyIn.DueTo", "0")
	a.Setel("PolicyTreatyIn.ProductionDate", "2099-01-01 00:00:00")
	if kode, isi := u.kirim(id, secHead, a); kode != http.StatusOK {
		t.Fatalf("%d %s", kode, isi)
	}
	g := u.g.Halaman[id]
	if got := g.Ambil("PolicyTreatyIn.PremiOgp"); got != "1000" {
		t.Fatalf("PremiOgp terkunci di layar atasan, tersimpan %q", got)
	}
	for j, tolak := range map[string]string{
		"PolicyTreatyIn.QuotationData.NoOfferSlip": "UJI-SLIP",
		"PolicyTreatyIn.FlagPPH":                   "true",
		"PolicyTreatyIn.DueTo":                     "0",
		"PolicyTreatyIn.ProductionDate":            "2099-01-01 00:00:00",
	} {
		if got := g.Ambil(j); got == tolak {
			t.Errorf("%s tidak dapat diisi atasan (AC 52), tersimpan %q", j, got)
		}
	}
	// Catatan atasan tersimpan di POOLDATA.HISTORYAKSEPTASIPRODUCTION (K4).
	usulan := u.g.Usulan
	if len(usulan) == 0 || usulan[len(usulan)-1].Keterangan != "UJI-1" || usulan[len(usulan)-1].Approval != "Accept" ||
		usulan[len(usulan)-1].AksesLogin != "UJI-SH" {
		t.Fatalf("Approval dan Suggest dapat diisi atasan: %+v", usulan)
	}
}

func TestSimpanDrafHanyaAdmin(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("1")
	h.Setel("PolicyTreatyIn.Remark", "UJI-remark")
	if kode, isi := u.panggil("PUT", "/kasus/"+id, admin, map[string]any{"halaman": h}); kode != http.StatusOK {
		t.Fatalf("simpan: %d %s", kode, isi)
	}
	if u.g.Halaman[id].Ambil("PolicyTreatyIn.Remark") != "UJI-remark" {
		t.Fatal("draf tersimpan")
	}
	u.kirim(id, admin, halamanLengkap("1"))
	if kode, _ := u.panggil("PUT", "/kasus/"+id, secHead, map[string]any{"halaman": putusan("1")}); kode != http.StatusConflict {
		t.Fatalf("layar atasan tidak punya Save: %d", kode)
	}
}

func TestPilihBisnis(t *testing.T) { // tiket 01; AC 36-38
	u := baru(t)
	id := u.buat()
	if kode, _ := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-TIDAK-ADA"}); kode != http.StatusUnprocessableEntity {
		t.Fatalf("kontrak tak ada: %d", kode)
	}
	if u.g.Halaman[id].Ambil("PolicyTreatyIn.NoOffer") != "" {
		t.Fatal("pembacaan gagal tidak boleh menyimpan apa pun")
	}
	u.g.Kontrak["UJI-D1"] = models.BarisKontrak{"ID": "UJI-D1", "TREATYID": "UJI-T1", "LIMITCURRENCY": "IDR",
		"CLASSOFBUSINESS": "UJI BISNIS", "TREATYTYPE": "", "PROPORTIONTYPE": "Proportional", "MDPVALUE": "500",
		"LAYERTYPE": "UJI-LT", "LAYER": "1", "LAYERPARTTYPE": "UJI-LPT", "LAYERPART": "2"}
	u.g.Bisnis["UJI BISNIS"] = models.BarisBisnis{OldID: "01", GroupPanel: "006", ID: "UJI-B1"}
	kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D1"})
	if kode != http.StatusOK {
		t.Fatalf("pilih bisnis: %d %s", kode, isi)
	}
	var ly services.Layar
	if err := json.Unmarshal([]byte(isi), &ly); err != nil {
		t.Fatal(err)
	}
	for j, harap := range map[string]string{
		"Quotation.BusinessType":                    "FireStyle2", // BusinessType_DeT 006/01
		"PolicyTreatyIn.QuotationData.BusinessType": "FireStyle2", // langkah 14.9
	} {
		if got := ly.Halaman.Ambil(j); got != harap {
			t.Errorf("layar %s = %q, harap %q", j, got, harap)
		}
	}
	// Tersimpan: medan berkolom saja. BusinessType TURUNAN masukannya yang
	// tersimpan (GroupPanel "006" + BusinessOldId "01" - diagram J38), tanpa kolom.
	h := u.g.Halaman[id]
	for j, harap := range map[string]string{
		"PolicyTreatyIn.NoOffer":       "UJI-T1",
		"PolicyTreatyIn.TreatyType":    "XOL", // langkah 5
		"PolicyTreatyIn.IDCurrency":    "UJI-ID-IDR",
		"PolicyTreatyIn.PremiOgp":      "500",
		"PolicyTreatyIn.BizCode":       "UJI-B1",
		"Quotation.GroupPanel":         "006",
		"Quotation.BusinessOldId":      "01",
		"Quotation.BusinessType":       "",
		"TreatyIn.ID":                  "UJI-D1",
		"PolicyTreatyIn.OJKBusinessID": "UJI-OJK",
		"PolicyTreatyIn.LayerType":     "", // diagram F26: dicoret dari T_GENERAL_POLIS_TREATY
	} {
		if got := h.Ambil(j); got != harap {
			t.Errorf("tersimpan %s = %q, harap %q", j, got, harap)
		}
	}
	if got := models.GolongkanJenisUsaha(h.Ambil("Quotation.GroupPanel"), h.Ambil("Quotation.BusinessOldId")); got != "FireStyle2" {
		t.Errorf("penggolong atas kolom tersimpan = %q (AC 14)", got)
	}
	// LAYER* tingkat polis = PANTULAN baris view (preACT langkah 3,
	// pxResults(1)) - dibaca balik lewat TreatyIn.ID saat layar dibuka (ID-22).
	_, isi = u.panggil("GET", "/kasus/"+id, admin, nil)
	ly = services.Layar{}
	if err := json.Unmarshal([]byte(isi), &ly); err != nil {
		t.Fatal(err)
	}
	for j, harap := range map[string]string{"PolicyTreatyIn.LayerType": "UJI-LT", "PolicyTreatyIn.Layer": "1",
		"PolicyTreatyIn.LayerPartType": "UJI-LPT", "PolicyTreatyIn.LayerPart": "2"} {
		if got := ly.Halaman.Ambil(j); got != harap {
			t.Errorf("dibuka ulang %s = %q, harap %q", j, got, harap)
		}
	}
}

// Pemetaan tempat -> peran KOSONG sampai IAM menjawab (K12, K16): setiap
// tempat TERTUNDA bagi siapa pun, dan Production Date tidak diwajibkan.
// Mekanisme arahnya diuji murni di models (TestTempatTampilMenurutArah).
func TestTempatBerperanTidakDitebak(t *testing.T) { // AC 81, 82, 91
	u := baru(t)
	id := u.buat()
	for _, p := range []pelakuUji{admin, {"UJI-A", models.PosisiAdmin + ",UJI-PERAN-A"}} {
		_, isi := u.panggil("GET", "/kasus/"+id, p, nil)
		var ly services.Layar
		if err := json.Unmarshal([]byte(isi), &ly); err != nil {
			t.Fatal(err)
		}
		if len(ly.Tempat) != len(models.SemuaTempat) {
			t.Fatalf("%s: kedua belas tempat tiket 05 dilaporkan layar (%v)", p.akun, ly.Tempat)
		}
		for kode, tampil := range ly.Tempat {
			if tampil {
				t.Fatalf("%s: tanpa pemetaan, tempat %s TERTUNDA - tidak tampil", p.akun, kode)
			}
		}
		for _, w := range ly.MedanWajib {
			if w == models.HalamanPolis+".ProductionDate" {
				t.Fatal("tempat tertunda tidak mewajibkan Production Date")
			}
		}
	}
}

// Mekanisme yang akan dipakai begitu IAM mengisi pemetaan: `ListSuggest.
// ProductionDate` tampil dan wajib menurut DUA tempat tersendiri (pyVisible,
// pyRequiredWhen) - satu sumber `models.MedanWajibBerlaku` untuk layar, Save,
// dan submit. Peran fiktif UJI- (AC 91); pemetaan dipulihkan sesudah uji.
func TestTanggalProduksiMengikutiPemetaanTempat(t *testing.T) { // AC 81, tiket 05
	lama := models.PemetaanPeranTempat
	t.Cleanup(func() { models.PemetaanPeranTempat = lama })
	models.PemetaanPeranTempat = []models.PeranTempat{
		{KodeTempat: models.TempatProduksiTampilOperator3, Peran: "UJI-PERAN-PROD", Arah: models.ArahMuncul},
		{KodeTempat: models.TempatProduksiWajibOperator3, Peran: "UJI-PERAN-PROD", Arah: models.ArahMuncul},
	}
	u := baru(t)
	berperan := pelakuUji{"UJI-ADMIN", models.PosisiAdmin + ",UJI-PERAN-PROD"}
	id := u.buat()
	// Pra-proses (InputPolicyTreatyInPre_Act 3-4, 9) selalu mengisi
	// ProductionDate; layar yang menampilkannya dapat mengosongkannya.
	h := halamanLengkap("1")
	h.Setel("PolicyTreatyIn.ProductionDate", "")
	kode, isi := u.kirim(id, berperan, h)
	if kode != http.StatusUnprocessableEntity || !strings.Contains(isi, "Production Date") {
		t.Fatalf("tempat wajib terbuka, Production Date kosong: %d %s", kode, isi)
	}
	if kode, isi := u.panggil("PUT", "/kasus/"+id, berperan, map[string]any{"halaman": h}); kode != http.StatusUnprocessableEntity ||
		!strings.Contains(isi, "Production Date") {
		t.Fatalf("Save juga menahan (AC 48): %d %s", kode, isi)
	}
	_, isi = u.panggil("GET", "/kasus/"+id, berperan, nil)
	var ly services.Layar
	if err := json.Unmarshal([]byte(isi), &ly); err != nil {
		t.Fatal(err)
	}
	if !ly.Tempat[models.TempatProduksiTampilOperator3] || ly.Tempat[models.TempatProduksiTampilOperator4] {
		t.Fatalf("tempat pelaku %v", ly.Tempat)
	}
	if w := strings.Join(ly.MedanWajib, " "); strings.Contains(w, "ProductionDate") {
		t.Fatalf("IsApproved tersimpan kosong: Production Date belum tampil/wajib (%s)", w)
	}
	h.Setel("PolicyTreatyIn.ProductionDate", "2026-10-31")
	if kode, isi := u.kirim(id, berperan, h); kode != http.StatusOK {
		t.Fatalf("Production Date terisi: %d %s", kode, isi)
	}
	if got := u.g.Halaman[id].Ambil("PolicyTreatyIn.ProductionDate"); !strings.HasPrefix(got, "2026-10-31") {
		t.Fatalf("Production Date tampil diterima dari layar, tersimpan %q", got)
	}
	// pelaku tanpa peran itu: tidak tampil (kiriman diabaikan), tidak wajib
	id2 := u.buat()
	if kode, isi := u.kirim(id2, admin, h); kode != http.StatusOK {
		t.Fatalf("tanpa peran: %d %s", kode, isi)
	}
	if got := u.g.Halaman[id2].Ambil("PolicyTreatyIn.ProductionDate"); strings.HasPrefix(got, "2026-10-31") {
		t.Fatalf("tanpa peran: Production Date tak tampil tidak diterima, tersimpan %q", got)
	}
}

func TestHitungTidakMenyimpan(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("")
	h.Setel("PolicyTreatyIn.RiCommOgp", "10")
	kode, isi := u.panggil("POST", "/kasus/"+id+"/hitung", admin, map[string]any{"urutan": []map[string]string{{"aksi": "CountResult1", "param": "Pct"}}, "halaman": h})
	if kode != http.StatusOK {
		t.Fatalf("%d %s", kode, isi)
	}
	var ly services.Layar
	if err := json.Unmarshal([]byte(isi), &ly); err != nil {
		t.Fatal(err)
	}
	d, err := models.AngkaTeks("ResultOgp1", ly.Halaman.Ambil("PolicyTreatyIn.ResultOgp1"))
	if err != nil || d.Cmp(apd.New(100, 0)) != 0 {
		t.Fatalf("ResultOgp1 = 1000 x 10/100: %v %v", d, err)
	}
	if u.g.Halaman[id].Ambil("PolicyTreatyIn.ResultOgp1") != "" {
		t.Fatal("refresh berhitung tidak menyimpan")
	}
}

// Action set sel `.ResultOgp1` `DetailPolicyTreatyIn`: event change ->
// `refresh CountResult1_Act(Data="Amount")` LALU `refresh CountOGPONP_Act` -
// dua activity berurutan atas clipboard yang sama. Hitung tangan (PremiOgp 1000,
// ResultOgp1 250, sisanya 0):
//
//	CountResult1_Act langkah 6   RiCommOgp = (250/1000) x 100 = 25
//	CountNetPremi_act langkah 4  NetPremium = (1000-250)+(0-0)-0-0-0-0 = 750
//	CountOGPONP_Act langkah 8    ClaimType = "" (Claim 0, Salvage 0) - menimpa isian
func TestHitungUrutanActionSet(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("")
	h.Setel("PolicyTreatyIn.ResultOgp1", "250")
	h.Setel("PolicyTreatyIn.ClaimType", "UJI-LAMA")
	kode, isi := u.panggil("POST", "/kasus/"+id+"/hitung", admin, map[string]any{
		"urutan":  []map[string]string{{"aksi": "CountResult1", "param": "Amount"}, {"aksi": "CountOGPONP"}},
		"halaman": h,
	})
	if kode != http.StatusOK {
		t.Fatalf("%d %s", kode, isi)
	}
	var ly services.Layar
	if err := json.Unmarshal([]byte(isi), &ly); err != nil {
		t.Fatal(err)
	}
	for m, harap := range map[string]int64{"RiCommOgp": 25, "NetPremium": 750} {
		d, err := models.AngkaTeks(m, ly.Halaman.Ambil("PolicyTreatyIn."+m))
		if err != nil || d.Cmp(apd.New(harap, 0)) != 0 {
			t.Errorf("%s = %v (%v), harap %d", m, d, err, harap)
		}
	}
	if got := ly.Halaman.Ambil("PolicyTreatyIn.ClaimType"); got != "" {
		t.Errorf("CountOGPONP_Act langkah 8 berjalan sesudah CountResult1_Act: ClaimType %q", got)
	}
	kode, _ = u.panggil("POST", "/kasus/"+id+"/hitung", admin, map[string]any{
		"urutan":  []map[string]string{{"aksi": "CountOGPONP"}, {"aksi": "UJI-Karangan"}},
		"halaman": h,
	})
	if kode != http.StatusBadRequest {
		t.Fatalf("aksi tak dikenal di urutan: %d", kode)
	}
}

// pengirimUji mencatat muatan dan memberi galat yang diminta.
type pengirimUji struct {
	galat  error
	muatan []models.MuatanKonversi
}

func (p *pengirimUji) Kirim(_ context.Context, m models.MuatanKonversi) error {
	p.muatan = append(p.muatan, m)
	return p.galat
}

// KEPUTUSAN-RONDE-12 butir 7: konversi sesudah selesai, empat medan, gagal
// tidak membatalkan penyimpanan.
func TestKonversiSesudahSelesai(t *testing.T) {
	for _, gagal := range []bool{false, true} {
		u := baru(t)
		p := &pengirimUji{}
		if gagal {
			p.galat = errors.New("UJI-jaringan")
		}
		u.s = handlers.Router(services.Baru(u.g, jam).DenganKonversi(p, true), true)
		id := u.buat()
		u.kirim(id, admin, halamanLengkap("1"))
		u.kirim(id, secHead, putusan("1"))
		_, isi := u.panggil("POST", "/kasus/"+id+"/nomor-polis", deptHead, map[string]any{"halaman": putusan("1")})
		var n services.NomorPolis
		_ = json.Unmarshal([]byte(isi), &n)
		kode, isi := u.kirim(id, deptHead, putusan("1"))
		if kode != http.StatusOK || u.g.Kasus[id].StatusWork != models.StatusSelesai {
			t.Fatalf("konversi gagal tidak boleh menggagalkan submit: %d %s", kode, isi)
		}
		if len(p.muatan) != 1 {
			t.Fatalf("satu konversi per penyelesaian, dapat %d", len(p.muatan))
		}
		m := p.muatan[0]
		if m.CARI1 != "ASM-FW-GISFW-WORK-NB "+id || m.CARI2 != n.PolicyNo || m.CARI3 != n.PolicyNo || m.CARI21 != "03/10/2026 09:00:00" {
			t.Fatalf("muatan %+v", m)
		}
		if gagal != strings.Contains(isi, models.PesanGagalKonversi) {
			t.Fatalf("gagal=%v, jawaban %s", gagal, isi)
		}
	}
}
