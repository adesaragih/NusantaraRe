package handlers_test

// Uji seam HTTP - daur hidup endorsemen lewat handler di atas gudang tiruan (`backend/tiruan`) yang membatalkan
// transaksi sungguh-sungguh: Create dari polis terbit (CreateEDMT), Choose Business (EDMChooseBusiness_Act),
// selisih (EDMTCalculateTreatyDifference), tangga tiga jenjang, Utility1 produksi, endorsemen BERLAPIS sampai
// generasi ketiga, penolakan admin (generasi dilepas), pembatalan (EDMType 4), dan wewenang. Fixture UJI-.

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/edmtreatyin/backend/handlers"
	"nusantarare/modul/edmtreatyin/backend/models"
	"nusantarare/modul/edmtreatyin/backend/services"
	"nusantarare/modul/edmtreatyin/backend/tiruan"
)

type pelakuUji struct{ akun, peran string }

var (
	admin    = pelakuUji{"UJI-ADMIN", "UJI-LAIN," + models.PosisiAdmin}
	secHead  = pelakuUji{"UJI-SH", models.PosisiSecHead}
	deptHead = pelakuUji{"UJI-DH", models.PosisiDeptHead}
	orang    = pelakuUji{"UJI-ORANG", "UJI-LAIN"}
)

func jam() time.Time { return time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC) }

const (
	nopolUji   = "UJI-POL-0001"
	kontrakUji = "UJI-KONTRAK-01"
	nbUji      = "NB-1"
)

type uji struct {
	t *testing.T
	g *tiruan.Gudang
	s http.Handler
}

// polisNBUji - generasi NB proporsional terbit (PRODKE 0): dua baris spreading, satu angsuran.
func polisNBUji() *models.Halaman {
	h := models.HalamanBaru()
	p := func(m, v string) { h.Setel(models.HalamanPolis+"."+m, v) }
	q := func(m, v string) {
		h.Setel(models.HalamanQuotation+"."+m, v)
		h.Setel(models.HalamanPolis+".QuotationData."+m, v)
	}
	p("NoOffer", kontrakUji)
	p("BizCode", "UJI-BIZ")
	p("BizName", "UJI BISNIS")
	p("SOB", "UJI-SOB")
	p("SOBName", "UJI SOB")
	p("CedingCo", "UJI-CED")
	p("CedingCoName", "UJI CEDING")
	p("Currency", "IDR")
	p("IDCurrency", "UJI-IDR")
	p("TreatyYear", "2026")
	p("TreatyType", "QS")
	p("TreatyGroupID", "UJI-TG")
	p("TreatyGroupName", "UJI GRUP")
	p("Quartal", "4")
	p("YearOfQuartal", "2026")
	p("StartDate", "2026-01-01")
	p("EndDate", "2026-12-31")
	p("MarketingOfficer", "UJI MO")
	p("IsNewPolicyNonProp", "0")
	p("Installment", "1")
	p("PremiOgp", "1000")
	p("RiCommOgp", "10")
	p("ResultOgp1", "100")
	p("NetPremium", "900")
	p("BalanceDueTo", "900")
	p("Claim", "0")
	q("ProportionalType", models.JenisProporsional)
	q("BusinessFac", "T") // BusinessFac Treaty
	q("MarketingName", "UJI MO")
	q("IsSurveyReport", "No")
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{
		{"TreatyType": "UJI-T1", "TreatyName": "UJI SATU", "SharePercentage": "60", "PremiumSpreaded": "540"},
		{"TreatyType": "UJI-T2", "TreatyName": "UJI DUA", "SharePercentage": "40", "PremiumSpreaded": "360"},
	})
	h.SetelDaftar(models.DaftarAngsuran, []models.Baris{
		{"InstallmentNo": "1", "DueDate": "2026-01-01", "InstallmentPercentage": "100", "Premium": "900", "PaymentTotal": "900"},
	})
	return h
}

func baru(t *testing.T) *uji {
	t.Helper()
	g := tiruan.Baru()
	g.Nama["UJI-ADMIN"] = "Uji Admin"
	g.Nama["UJI-SH"] = "Uji Sec"
	g.Nama["UJI-DH"] = "Uji Dept"
	g.MataUang["IDR"] = "UJI-IDR"
	g.TanamPolis(nbUji, nopolUji, 0, "", polisNBUji(), kontrakUji)
	// master kontrak proporsional (RDB BrowseTreatyIn): tanpa Share / Installment XOL
	g.Master["ID:"+kontrakUji] = []models.MasterXOL{{Nilai: map[string]string{"ID": kontrakUji}, Daftar: map[string][]models.Baris{}}}
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
	r := httptest.NewRequest(metode, handlers.Prefix+jalur, &b)
	if p.akun != "" {
		r.Header.Set("X-Pelaku", p.akun)
		r.Header.Set("X-Peran", p.peran)
	}
	w := httptest.NewRecorder()
	u.s.ServeHTTP(w, r)
	return w.Code, w.Body.String()
}

func (u *uji) buat(edmType string) models.Kasus {
	u.t.Helper()
	kode, isi := u.panggil("POST", "/kasus", admin, map[string]string{"noPolis": nopolUji, "edmType": edmType})
	if kode != http.StatusCreated {
		u.t.Fatalf("buat: %d %s", kode, isi)
	}
	var k models.Kasus
	if err := json.Unmarshal([]byte(isi), &k); err != nil {
		u.t.Fatal(err)
	}
	return k
}

func (u *uji) layar(id string, p pelakuUji) services.Layar {
	u.t.Helper()
	kode, isi := u.panggil("GET", "/kasus/"+id, p, nil)
	if kode != http.StatusOK {
		u.t.Fatalf("buka %s: %d %s", id, kode, isi)
	}
	var ly services.Layar
	if err := json.Unmarshal([]byte(isi), &ly); err != nil {
		u.t.Fatal(err)
	}
	return ly
}

func (u *uji) wajib(kode int, isi string, harap int, apa string) {
	u.t.Helper()
	if kode != harap {
		u.t.Fatalf("%s: %d (harap %d) %s", apa, kode, harap, isi)
	}
}

// pilihBisnis = tombol Choose popup -> EDMChooseBusiness_Act.
func (u *uji) pilihBisnis(id string) services.Layar {
	u.t.Helper()
	kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"halaman": models.HalamanBaru()})
	u.wajib(kode, isi, http.StatusOK, "pilih bisnis")
	var ly services.Layar
	if err := json.Unmarshal([]byte(isi), &ly); err != nil {
		u.t.Fatal(err)
	}
	return ly
}

// isianAdmin - kiriman layar admin: PremiOgp baru, spreading DUA baris (baris lama ke-2 ditambahkan lagi), putusan.
func isianAdmin(premi, isApproved string) *models.Halaman {
	h := models.HalamanBaru()
	h.Setel(models.HalamanPolis+".PremiOgp", premi)
	h.Setel(models.HalamanPolis+".IsApproved", isApproved)
	h.Setel(models.HalamanPolis+".Suggest", "UJI-catatan")
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{
		{"TreatyType": "UJI-T1", "SharePercentage": "60"},
		{"TreatyType": "UJI-T2", "SharePercentage": "40"},
	})
	return h
}

func putusan(isApproved string) *models.Halaman {
	h := models.HalamanBaru()
	h.Setel(models.HalamanPolis+".IsApproved", isApproved)
	h.Setel(models.HalamanPolis+".Suggest", "UJI-"+isApproved)
	return h
}

func (u *uji) kirim(id string, p pelakuUji, h *models.Halaman) (int, string) {
	u.t.Helper()
	return u.panggil("POST", "/kasus/"+id+"/kirim", p, map[string]any{"halaman": h})
}

// tuntaskan menjalankan satu endorsemen sampai Dept Head menyetujui (Utility1).
func (u *uji) tuntaskan(id, premi string) {
	u.t.Helper()
	u.pilihBisnis(id)
	kode, isi := u.kirim(id, admin, isianAdmin(premi, "1"))
	u.wajib(kode, isi, http.StatusOK, "admin setuju")
	kode, isi = u.kirim(id, secHead, putusan("1"))
	u.wajib(kode, isi, http.StatusOK, "sec head setuju")
	kode, isi = u.kirim(id, deptHead, putusan("1"))
	u.wajib(kode, isi, http.StatusOK, "dept head setuju")
}

func TestCreateMelahirkanGenerasiBerikut(t *testing.T) {
	u := baru(t)
	k := u.buat("1")
	if !strings.HasPrefix(k.ID, models.AwalanKasus) || k.ProdKe != 1 || k.OldPolisID != nbUji || k.PositionNote != models.PosisiAdmin {
		t.Fatalf("kasus lahir = %+v", k)
	}
	ly := u.layar(k.ID, admin)
	h := ly.Halaman
	cek := map[string]string{
		models.HalamanPolis + ".EDMNo":            nopolUji + "/E01", // SetEDMTNoPolis 3
		models.HalamanPolis + ".PolicyNo":         nopolUji,
		models.HalamanPolis + ".EDMType":          "1",
		models.HalamanPolis + ".OldData.PolicyNo": nopolUji,
		models.HalamanPolis + ".OldData.PremiOgp": "1000",
		models.HalamanPolis + ".OldData.NoOffer":  kontrakUji,
		models.HalamanQuotation + ".OldPolicyNo":  nopolUji,  // CreateEDMT 12
		models.HalamanPolis + ".IDCurrency":       "UJI-IDR", // CreateEDMT 14
		"NBStatus":                                "EDM IS IN UJI-ADMIN'S INBOX",
	}
	for j, v := range cek {
		if h.Ambil(j) != v {
			t.Errorf("%s = %q, harap %q", j, h.Ambil(j), v)
		}
	}
	// data baru KOSONG sampai Choose Business (MergePage baru)
	if h.Ambil(models.HalamanPolis+".PremiOgp") != "" || len(h.AmbilDaftar(models.DaftarSpreading)) != 0 {
		t.Fatalf("data baru harus kosong sebelum Choose Business")
	}
	if len(h.AmbilDaftar(models.HalamanPolis+".OldData.SpreadingRiskList")) != 2 {
		t.Fatal("OldData.SpreadingRiskList = generasi NB (2 baris)")
	}
}

func TestEDMKeduaSelamaBerjalanDitolak(t *testing.T) {
	u := baru(t)
	u.buat("1")
	kode, isi := u.panggil("POST", "/kasus", admin, map[string]string{"noPolis": nopolUji, "edmType": "1"})
	u.wajib(kode, isi, http.StatusUnprocessableEntity, "EDM kedua")
	if !strings.Contains(isi, models.PesanEDMBelumSelesai) {
		t.Fatalf("pesan TrtEdmCheckPolicyError: %s", isi)
	}
	kode, isi = u.panggil("GET", "/periksa-polis?nopolis="+nopolUji, admin, nil)
	u.wajib(kode, isi, http.StatusOK, "periksa polis")
	if !strings.Contains(isi, kontrakUji) || !strings.Contains(isi, models.PesanEDMBelumSelesai) {
		t.Fatalf("periksa polis: %s", isi)
	}
	kode, isi = u.panggil("GET", "/periksa-polis?nopolis=UJI-TIDAK-ADA", admin, nil)
	u.wajib(kode, isi, http.StatusOK, "periksa polis salah")
	if !strings.Contains(isi, models.PesanNopolisSalah) {
		t.Fatalf("CheckNopolisAvailability: %s", isi)
	}
}

func TestSelisihTerhadapGenerasiTepatSebelumnyaSampaiLapisKetiga(t *testing.T) {
	u := baru(t)
	k1 := u.buat("1")
	u.tuntaskan(k1.ID, "1500")
	s1 := u.g.Selisih[k1.ID]
	if s1 == nil || s1.Halaman.Ambil(models.HalamanPolis+".TreatyDifference.PremiOgp") != "500" {
		t.Fatalf("selisih gen 1 = %+v", s1)
	}
	if s1.Kunci.Sumber != models.SumberGo || s1.Kunci.EDMNo != nopolUji+"/E01" || s1.Kunci.ProdKe != 1 || s1.Kunci.NoPolis != nopolUji {
		t.Fatalf("kunci saring T_POLIS_DIFFERENCE = %+v (AC 28)", s1.Kunci)
	}
	if u.g.Generasi[k1.ID].NoPolis != nopolUji || u.g.Kasus[k1.ID].StatusWork != models.StatusSelesai {
		t.Fatalf("gen 1 selesai: %+v %+v", u.g.Generasi[k1.ID], u.g.Kasus[k1.ID])
	}
	// generasi ke-2: OldData = generasi 1 (bukan NB)
	k2 := u.buat("1")
	if k2.ProdKe != 2 || k2.OldPolisID != k1.ID {
		t.Fatalf("gen 2 = %+v", k2)
	}
	h2 := u.layar(k2.ID, admin).Halaman
	if h2.Ambil(models.HalamanPolis+".EDMNo") != nopolUji+"/E02" || h2.Ambil(models.HalamanPolis+".OldData.EDMNo") != nopolUji+"/E01" ||
		h2.Ambil(models.HalamanPolis+".OldData.PremiOgp") != "1500" {
		t.Fatalf("gen 2 OldData: EDMNo %q old %q premi %q", h2.Ambil(models.HalamanPolis+".EDMNo"),
			h2.Ambil(models.HalamanPolis+".OldData.EDMNo"), h2.Ambil(models.HalamanPolis+".OldData.PremiOgp"))
	}
	// tab Old Data PropOldData2: selisih generasi sebelumnya
	if h2.Ambil(models.HalamanPolis+".OldData.TreatyDifference.PremiOgp") != "500" {
		t.Fatal("OldData.TreatyDifference = selisih generasi 1")
	}
	u.tuntaskan(k2.ID, "1700")
	// ⛔ SATU RUMUS (ID-28, AC 7, 16): 1700 - 1500 = 200, BUKAN 1700 - 500 (varian berlapis) atau 1700 - 1000 (NB)
	if v := u.g.Selisih[k2.ID].Halaman.Ambil(models.HalamanPolis + ".TreatyDifference.PremiOgp"); v != "200" {
		t.Fatalf("selisih gen 2 = %q, harap 200", v)
	}
	k3 := u.buat("1")
	if k3.ProdKe != 3 || k3.OldPolisID != k2.ID {
		t.Fatalf("gen 3 = %+v", k3)
	}
	u.tuntaskan(k3.ID, "1650")
	if v := u.g.Selisih[k3.ID].Halaman.Ambil(models.HalamanPolis + ".TreatyDifference.PremiOgp"); v != "-50" {
		t.Fatalf("selisih gen 3 = %q, harap -50", v)
	}
	if n := len(u.g.Simpanan); n != 3 {
		t.Fatalf("Utility1 = %d kali, harap 3", n)
	}
}

// AC 8 / ID-15 (keutuhan generasi): setiap baris spreading generasi lama ada di generasi baru. Sejak keputusan WO
// 07-10-2026 (Add dibuang) Choose Business membawa SEMUA baris lama, dan kiriman layar yang lebih pendek tidak dapat
// menghapusnya (ID-16) - penjaga `BarisSpreadingHilang` tetap ada untuk data yang tersimpan tidak utuh.
func TestBarisSpreadingGenerasiLamaWajibAdaSaatSubmit(t *testing.T) {
	u := baru(t)
	k := u.buat("1")
	ly := u.pilihBisnis(k.ID)
	if n := len(ly.Halaman.AmbilDaftar(models.DaftarSpreading)); n != 2 {
		t.Fatalf("spreading sesudah Choose Business = %d baris, harap 2 (semua baris generasi lama)", n)
	}
	h := isianAdmin("1500", "1")
	h.SetelDaftar(models.DaftarSpreading, h.AmbilDaftar(models.DaftarSpreading)[:1])
	kode, isi := u.kirim(k.ID, admin, h)
	u.wajib(kode, isi, http.StatusOK, "kiriman satu baris")
	if n := len(u.g.Generasi[k.ID].Halaman.AmbilDaftar(models.DaftarSpreading)); n != 2 {
		t.Fatalf("baris lama terhapus kiriman layar: %d baris (ID-16)", n)
	}
}

func TestAdminMenolakMelepasGenerasi(t *testing.T) {
	u := baru(t)
	k := u.buat("1")
	u.pilihBisnis(k.ID)
	kode, isi := u.kirim(k.ID, admin, isianAdmin("1500", "0"))
	u.wajib(kode, isi, http.StatusOK, "admin tolak")
	if u.g.Kasus[k.ID].StatusWork != models.StatusDitolak || u.g.Generasi[k.ID].OldPolisID != "" {
		t.Fatalf("tolak admin: %+v %+v", u.g.Kasus[k.ID], u.g.Generasi[k.ID])
	}
	if len(u.g.Simpanan) != 0 {
		t.Fatal("Decision3 No -> End3 TANPA Utility1")
	}
	if !strings.HasPrefix(u.g.Generasi[k.ID].Halaman.Ambil("NBStatus"), "EDM WAS DECLINED BY  UJI ADMIN") {
		t.Fatalf("NBStatus tolak = %q", u.g.Generasi[k.ID].Halaman.Ambil("NBStatus"))
	}
	// polis boleh di-endorse lagi, nomor generasi sama (Pega: json_polis tidak ditulis)
	k2 := u.buat("1")
	if k2.ProdKe != 1 || k2.OldPolisID != nbUji {
		t.Fatalf("endorse ulang sesudah tolak = %+v", k2)
	}
	// generasi tolak tetap terbaca OldData-nya lewat (OldPolicyNo, PRODKE - 1)
	if h := u.layar(k.ID, admin).Halaman; h.Ambil(models.HalamanPolis+".OldData.PremiOgp") != "1000" {
		t.Fatal("OldData kasus tolak")
	}
}

func TestAtasanMenolakKembaliKeAdmin(t *testing.T) {
	u := baru(t)
	k := u.buat("1")
	u.pilihBisnis(k.ID)
	kode, isi := u.kirim(k.ID, admin, isianAdmin("1500", "1"))
	u.wajib(kode, isi, http.StatusOK, "admin setuju")
	if kk := u.g.Kasus[k.ID]; kk.PositionNote != models.PosisiSecHead || kk.Position != models.PositionAtasan {
		t.Fatalf("sesudah admin: %+v", kk)
	}
	kode, isi = u.kirim(k.ID, secHead, putusan("0"))
	u.wajib(kode, isi, http.StatusOK, "sec head tolak")
	if kk := u.g.Kasus[k.ID]; kk.PositionNote != models.PosisiAdmin || kk.Position != models.PositionAdmin {
		t.Fatalf("Transition8: %+v", kk)
	}
	kode, isi = u.kirim(k.ID, admin, isianAdmin("1500", "1"))
	u.wajib(kode, isi, http.StatusOK, "admin setuju lagi")
	kode, isi = u.kirim(k.ID, secHead, putusan("1"))
	u.wajib(kode, isi, http.StatusOK, "sec head setuju")
	if kk := u.g.Kasus[k.ID]; kk.PositionNote != models.PosisiDeptHead {
		t.Fatalf("tangga tiga jenjang: %+v", kk)
	}
	kode, isi = u.kirim(k.ID, deptHead, putusan("0"))
	u.wajib(kode, isi, http.StatusOK, "dept head tolak")
	if kk := u.g.Kasus[k.ID]; kk.PositionNote != models.PosisiAdmin || kk.Position != models.PositionAtasan {
		t.Fatalf("Transition4: %+v", kk)
	}
}

func TestPembatalanNolkanDataBaru(t *testing.T) {
	u := baru(t)
	k := u.buat("4")
	ly := u.pilihBisnis(k.ID)
	h := ly.Halaman
	if h.Ambil(models.HalamanPolis+".PremiOgp") != "0" || len(h.AmbilDaftar(models.DaftarSpreading)) != 2 {
		t.Fatalf("SetEDMTCancel: PremiOgp %q, spreading %d", h.Ambil(models.HalamanPolis+".PremiOgp"), len(h.AmbilDaftar(models.DaftarSpreading)))
	}
	// tombol "Calculate Value Difference" (PropNewData2 S24)
	kode, isi := u.panggil("POST", "/kasus/"+k.ID+"/hitung", admin, map[string]any{
		"urutan": []map[string]string{{"aksi": "EDMTCalculateTreatyDifference"}}, "halaman": models.HalamanBaru()})
	u.wajib(kode, isi, http.StatusOK, "Calculate Value Difference")
	var hasil services.Layar
	if err := json.Unmarshal([]byte(isi), &hasil); err != nil {
		t.Fatal(err)
	}
	if v := hasil.Halaman.Ambil(models.HalamanPolis + ".TreatyDifference.PremiOgp"); v != "-1000" {
		t.Fatalf("selisih pembatalan = %q, harap -1000 (AC 14)", v)
	}
}

func TestWewenangDanPosisi(t *testing.T) {
	u := baru(t)
	k := u.buat("1")
	kode, isi := u.kirim(k.ID, secHead, putusan("1"))
	u.wajib(kode, isi, http.StatusForbidden, "bukan antrean")
	kode, isi = u.panggil("POST", "/kasus/"+k.ID+"/pilih-bisnis", orang, map[string]any{"halaman": models.HalamanBaru()})
	u.wajib(kode, isi, http.StatusForbidden, "orang lain")
	kode, isi = u.panggil("GET", "/kasus", admin, nil)
	u.wajib(kode, isi, http.StatusOK, "portal")
	if !strings.Contains(isi, k.ID) {
		t.Fatalf("portal pembuat: %s", isi)
	}
	kode, isi = u.panggil("GET", "/kasus", secHead, nil)
	u.wajib(kode, isi, http.StatusOK, "portal atasan")
	if strings.Contains(isi, k.ID) {
		t.Fatal("portal = filter A pembuat")
	}
	kode, isi = u.panggil("GET", "/kotak-masuk", admin, nil)
	u.wajib(kode, isi, http.StatusOK, "kotak masuk")
	if !strings.Contains(isi, `"jumlah":1`) {
		t.Fatalf("kotak masuk admin: %s", isi)
	}
}

// angkaSama - dua teks angka bernilai sama (format penulis boleh berbeda).
func angkaSama(t *testing.T, apa, v, harap string) {
	t.Helper()
	a, _, err := apd.NewFromString(v)
	b, _, _ := apd.NewFromString(harap)
	if err != nil || a.Cmp(b) != 0 {
		t.Fatalf("%s = %q, harap %s", apa, v, harap)
	}
}

// Tinjauan kode 06-10-2026 (1): selisih dihitung ulang server SETIAP generasi ditulis - tidak bergantung pada tombol
// "Calculate Value Difference" atau sel uang yang berubah (diagram EDM Prop J98: "bila beda dari hitung ulang,
// TABELNYA yang salah"). Pembatalan, Submit tanpa sentuhan: selisih = 0 - lama.
func TestSelisihDihitungUlangSaatSubmitTanpaTombolHitung(t *testing.T) {
	u := baru(t)
	k := u.buat("4")
	u.pilihBisnis(k.ID)
	kode, isi := u.kirim(k.ID, admin, putusan("1"))
	u.wajib(kode, isi, http.StatusOK, "admin setuju tanpa sentuhan")
	s := u.g.Selisih[k.ID]
	if s == nil {
		t.Fatal("proyeksi selisih tidak ditulis")
	}
	angkaSama(t, "TreatyDifference.PremiOgp", s.Halaman.Ambil(models.HalamanPolis+".TreatyDifference.PremiOgp"), "-1000")
	if sp := s.Halaman.AmbilDaftar(models.DaftarSelisihSpreading); len(sp) != 2 {
		t.Fatalf("selisih spreading %d baris, harap 2", len(sp))
	} else {
		angkaSama(t, "spreading(1).PremiumSpreaded", sp[0]["PremiumSpreaded"], "-540")
	}
	// Tinjauan (3): kunci NOPOLIS proyeksi mengikuti generasinya - kosong selama berjalan
	if s.Kunci.NoPolis != "" {
		t.Fatalf("NOPOLIS selisih selama berjalan = %q, harap kosong", s.Kunci.NoPolis)
	}
	kode, isi = u.kirim(k.ID, secHead, putusan("1"))
	u.wajib(kode, isi, http.StatusOK, "sec head setuju")
	kode, isi = u.kirim(k.ID, deptHead, putusan("1"))
	u.wajib(kode, isi, http.StatusOK, "dept head setuju")
	if s := u.g.Selisih[k.ID]; s.Kunci.NoPolis != nopolUji {
		t.Fatalf("NOPOLIS selisih sesudah selesai = %q", s.Kunci.NoPolis)
	}
	if len(u.g.Simpanan) != 1 || len(u.g.Simpanan[0].Produksi) == 0 {
		t.Fatalf("Utility1 tanpa baris produksi: %+v", u.g.Simpanan)
	}
	angkaSama(t, "TREATYINPRODUCTION.PREMI_OGP", u.g.Simpanan[0].Produksi[0]["PREMI_OGP"], "-1000")
}

// Tinjauan kode 06-10-2026 (2): json_polis memuat generasi yang belum ada di tabel relasional (dokumen lama belum
// dimuat pemuat EDM) - Create DITOLAK, bukan melahirkan EDMNo kembar di atas generasi yang tertinggal.
func TestCreateDitolakBilaJSONPolisLebihMaju(t *testing.T) {
	u := baru(t)
	u.g.JSONPolis[nopolUji]++ // E01 Pega di json_polis, generasinya belum dimuat
	kode, isi := u.panggil("POST", "/kasus", admin, map[string]string{"noPolis": nopolUji, "edmType": "1"})
	if kode != http.StatusUnprocessableEntity || !strings.Contains(isi, models.PesanGenerasiBelumDimuat) {
		t.Fatalf("create di atas generasi tertinggal: %d %s", kode, isi)
	}
	if len(u.g.Kasus) != 0 {
		t.Fatal("kasus tetap lahir")
	}
}

// Tinjauan kode 06-10-2026 (3): admin menolak -> generasi dilepas DAN proyeksi selisihnya dibuang (tidak menyisakan
// baris berkunci saring sama dengan endorsemen berikutnya).
func TestAdminMenolakMembuangSelisih(t *testing.T) {
	u := baru(t)
	k := u.buat("1")
	u.pilihBisnis(k.ID)
	kode, isi := u.panggil("PUT", "/kasus/"+k.ID, admin, map[string]any{"halaman": isianAdmin("1500", "")})
	u.wajib(kode, isi, http.StatusOK, "simpan draf")
	if u.g.Selisih[k.ID] == nil {
		t.Fatal("draf tanpa proyeksi selisih")
	}
	kode, isi = u.kirim(k.ID, admin, isianAdmin("1500", "0"))
	u.wajib(kode, isi, http.StatusOK, "admin tolak")
	if s := u.g.Selisih[k.ID]; s != nil {
		t.Fatalf("selisih generasi tolak tersisa: %+v", s.Kunci)
	}
}

// Keputusan work owner 07-10-2026 (aturan portal NB berlaku untuk EDM): In Progress = buatan akun ini yang masih
// proses; Resolved (`?status=selesai`) = SEMUA berkas selesai siapa pun pembuatnya; berkas selesai hanya-baca.
func TestPortalInProgressDanResolved(t *testing.T) {
	u := baru(t)
	k := u.buat("1")
	u.tuntaskan(k.ID, "1500")
	adminLain := pelakuUji{"UJI-ADMIN2", models.PosisiAdmin}
	daftar := func(p pelakuUji, kueri string) []models.RingkasanKasus {
		t.Helper()
		kode, isi := u.panggil("GET", "/kasus"+kueri, p, nil)
		u.wajib(kode, isi, http.StatusOK, "portal "+kueri)
		var out []models.RingkasanKasus
		if err := json.Unmarshal([]byte(isi), &out); err != nil {
			t.Fatal(err)
		}
		return out
	}
	if d := daftar(admin, ""); len(d) != 0 {
		t.Fatalf("In Progress memuat berkas selesai: %+v", d)
	}
	if d := daftar(adminLain, "?status=selesai"); len(d) != 1 || d[0].ID != k.ID {
		t.Fatalf("Resolved harus memuat berkas selesai buatan akun lain: %+v", d)
	} else if d[0].NoPolis != nopolUji || d[0].EDMNo != nopolUji+"/E01" || d[0].TglProd == "" {
		// WO 07-10-2026 "kalo dah resolve tambahin kolom nopolisnya": kolom XML "Policy Number" = NOPOLIS generasi selesai
		t.Fatalf("Resolved: Policy Number %q / EDM Number %q / Production Date %q", d[0].NoPolis, d[0].EDMNo, d[0].TglProd)
	}
	if ly := u.layar(k.ID, adminLain); ly.BolehKerja {
		t.Fatal("berkas Resolved harus hanya-baca")
	}
}

// Keputusan work owner 07-10-2026: atasan TIDAK boleh mengubah With Tax, Type Tax, Overiding Commision, Marketing
// Officer (XML tidak menguncinya per posisi - penyimpangan sadar).
func TestAtasanTidakBolehMengubahHeader(t *testing.T) {
	u := baru(t)
	k := u.buat("1")
	u.pilihBisnis(k.ID)
	kode, isi := u.kirim(k.ID, admin, isianAdmin("1500", "1"))
	u.wajib(kode, isi, http.StatusOK, "admin setuju")
	sebelum := u.layar(k.ID, secHead).Halaman
	h := putusan("1")
	h.Setel(models.HalamanPolis+".FlagPPH", "true")
	h.Setel(models.HalamanPolis+".TypeTax", "UJI-PAJAK")
	h.Setel(models.HalamanPolis+".FlagRetroTreaty", "true")
	h.Setel(models.HalamanPolis+".QuotationData.MOID", "UJI-MO-LAIN")
	kode, isi = u.kirim(k.ID, secHead, h)
	u.wajib(kode, isi, http.StatusOK, "sec head setuju")
	sesudah := u.g.Generasi[k.ID].Halaman
	for _, m := range []string{"FlagPPH", "TypeTax", "FlagRetroTreaty", "QuotationData.MOID"} {
		j := models.HalamanPolis + "." + m
		if sesudah.Ambil(j) != sebelum.Ambil(j) {
			t.Errorf("%s diubah atasan: %q -> %q", m, sebelum.Ambil(j), sesudah.Ambil(j))
		}
	}
	for _, aksi := range []string{"RemoveTypeTax", "HitungPajak", "CheckDataMkt"} {
		kode, _ := u.panggil("POST", "/kasus/"+k.ID+"/hitung", deptHead, map[string]any{
			"urutan": []map[string]string{{"aksi": aksi}}, "halaman": models.HalamanBaru()})
		if kode == http.StatusOK {
			t.Errorf("aksi header %s terbuka bagi atasan", aksi)
		}
	}
}

// DT `TreatyEDMListType` (screenshot Pega, work owner 07-10-2026): kode -> teks "Source of Change".
func TestAcuanJenisEDMBerlabel(t *testing.T) {
	u := baru(t)
	kode, isi := u.panggil("GET", "/acuan", admin, nil)
	u.wajib(kode, isi, http.StatusOK, "acuan")
	var a services.Acuan
	if err := json.Unmarshal([]byte(isi), &a); err != nil {
		t.Fatal(err)
	}
	harap := []models.Pilihan{{Nilai: "1", Label: "Internal"}, {Nilai: "2", Label: "External"},
		{Nilai: "3", Label: "Adjustment Premium"}, {Nilai: "4", Label: "Cancel Input"}}
	if len(a.JenisEDM) != len(harap) {
		t.Fatalf("jenisEdm = %+v", a.JenisEDM)
	}
	for i := range harap {
		if a.JenisEDM[i] != harap[i] {
			t.Errorf("jenisEdm[%d] = %+v, harap %+v", i, a.JenisEDM[i], harap[i])
		}
	}
}

// Keputusan work owner 07-10-2026 (grid spreading hanya-baca, Add dibuang - sama dengan NB): Choose Business membawa
// SETIAP baris generasi lama; baris yang dikirim layar di belakang baris server diabaikan; Type Treaty tidak dapat
// diubah. Baris lama ber-%Share kosong tetap ditolak dengan pesan di baris itu (rekomendasi b).
func TestSpreadingHanyaBacaBarisLamaDibawa(t *testing.T) {
	u := baru(t)
	k := u.buat("1")
	ly := u.pilihBisnis(k.ID)
	if d := ly.Halaman.AmbilDaftar(models.DaftarSpreading); len(d) != 2 || d[1]["TreatyType"] != "UJI-T2" {
		t.Fatalf("Choose Business harus membawa kedua baris generasi lama: %+v", d)
	}
	h := isianAdmin("1500", "1")
	sp := h.AmbilDaftar(models.DaftarSpreading)
	sp[0]["TreatyType"] = "UJI-GANTI"
	h.SetelDaftar(models.DaftarSpreading, append(sp, models.Baris{"TreatyType": "UJI-T3", "SharePercentage": "5"}))
	kode, isi := u.kirim(k.ID, admin, h)
	u.wajib(kode, isi, http.StatusOK, "admin setuju (baris tambahan diabaikan)")
	d := u.g.Generasi[k.ID].Halaman.AmbilDaftar(models.DaftarSpreading)
	if len(d) != 2 || d[0]["TreatyType"] != "UJI-T1" {
		t.Fatalf("Add / ubah Type Treaty dari layar harus diabaikan: %+v", d)
	}

	// baris lama ber-%Share kosong: pesan tepat sekali, Submit ditahan
	nb := polisNBUji()
	sp = nb.AmbilDaftar(models.DaftarSpreading)
	delete(sp[1], "SharePercentage")
	nb.SetelDaftar(models.DaftarSpreading, sp)
	u.g.TanamPolis("NB-2", "UJI-POL-0002", 0, "", nb, kontrakUji)
	kode, isi = u.panggil("POST", "/kasus", admin, map[string]string{"noPolis": "UJI-POL-0002", "edmType": "1"})
	u.wajib(kode, isi, http.StatusCreated, "buat kasus 2")
	var k2 models.Kasus
	if err := json.Unmarshal([]byte(isi), &k2); err != nil {
		t.Fatal(err)
	}
	u.pilihBisnis(k2.ID)
	kode, isi = u.kirim(k2.ID, admin, putusan("1"))
	var tolak struct {
		Pesan []string `json:"pesan"`
	}
	_ = json.Unmarshal([]byte(isi), &tolak)
	if kode != http.StatusUnprocessableEntity || len(tolak.Pesan) != 1 || tolak.Pesan[0] != models.PesanShareSpreadingKosong(2) {
		t.Fatalf("baris lama tanpa %%Share: %d %s", kode, isi)
	}
}
