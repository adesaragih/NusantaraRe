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
	a.Setel("PolicyTreatyIn.QuotationData.NoOfferSlip", "UJI-SLIP")
	if kode, isi := u.kirim(id, secHead, a); kode != http.StatusOK {
		t.Fatalf("%d %s", kode, isi)
	}
	if got := u.g.Halaman[id].Ambil("PolicyTreatyIn.PremiOgp"); got != "1000" {
		t.Fatalf("PremiOgp terkunci di layar atasan, tersimpan %q", got)
	}
	if got := u.g.Halaman[id].Ambil("PolicyTreatyIn.QuotationData.NoOfferSlip"); got != "UJI-SLIP" {
		t.Fatalf("No Offer Slip dapat diisi atasan, tersimpan %q", got)
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
		"CLASSOFBUSINESS": "UJI BISNIS", "TREATYTYPE": "", "PROPORTIONTYPE": "Proportional", "MDPVALUE": "500"}
	u.g.Bisnis["UJI BISNIS"] = models.BarisBisnis{OldID: "01", GroupPanel: "006", ID: "UJI-B1"}
	kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D1"})
	if kode != http.StatusOK {
		t.Fatalf("pilih bisnis: %d %s", kode, isi)
	}
	h := u.g.Halaman[id]
	for j, harap := range map[string]string{
		"PolicyTreatyIn.NoOffer":                    "UJI-T1",
		"PolicyTreatyIn.TreatyType":                 "XOL", // langkah 5
		"PolicyTreatyIn.IDCurrency":                 "UJI-ID-IDR",
		"PolicyTreatyIn.PremiOgp":                   "500",
		"PolicyTreatyIn.BizCode":                    "UJI-B1",
		"Quotation.BusinessType":                    "FireStyle2", // BusinessType_DeT 006/01
		"PolicyTreatyIn.QuotationData.BusinessType": "FireStyle2", // langkah 14.9
		"TreatyIn.ID":                               "UJI-D1",
		"PolicyTreatyIn.OJKBusinessID":              "UJI-OJK",
	} {
		if got := h.Ambil(j); got != harap {
			t.Errorf("%s = %q, harap %q", j, got, harap)
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
		if tampil, ada := ly.Tempat[services.TempatTanggalProduksi]; !ada || tampil {
			t.Fatalf("%s: tanpa pemetaan, tempat TERTUNDA - tidak tampil (%v)", p.akun, ly.Tempat)
		}
		for _, w := range ly.MedanWajib {
			if w == models.HalamanPolis+".ProductionDate" {
				t.Fatal("tempat tertunda tidak mewajibkan Production Date")
			}
		}
	}
}

func TestHitungTidakMenyimpan(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("")
	h.Setel("PolicyTreatyIn.RiCommOgp", "10")
	kode, isi := u.panggil("POST", "/kasus/"+id+"/hitung", admin, map[string]any{"aksi": "CountResult1", "param": "Pct", "halaman": h})
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
