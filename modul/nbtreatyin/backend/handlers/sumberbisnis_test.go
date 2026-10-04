package handlers_test

// Uji seam 1 pemilih Source Of Business (XOL Retro) - F4 putaran 3, IKUTI XML:
//
//	tombol `Select Source Of Business` (`Section/DetailPolicyTreatyIn`, pyVisible
//	  `.ClaimType = 'XOL Retro'`) -> showHarness `SOB` -> TreeGrid `SourceHierarki`
//	klik baris -> flow action `AgentSourceBizDetails` -> pra-proses
//	  `SearchHierarkiSourceBizAgent_PostDT`: menulis `pyWorkPage.Quotation.*` di
//	  clipboard SAJA (nol Obj-Save) - hasilnya dipegang layar
//	Save / Submit layar admin -> hasil itu disimpan bersama halaman, diterima
//	  server HANYA bila cocok dengan RD `BrowseAgentHierarkiList_RD` yang
//	  dijalankan ulang
//
// Helper (`baru`, `halamanLengkap`, pelaku) milik alur_test.go; `layarDari`,
// `angkaSama` milik logika_test.go.

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func agenUji() []models.BarisAgen {
	return []models.BarisAgen{
		{ID: "UJI-AG-1", ClientName: "UJI SUMBER SATU", ChildCount: "0", ClientID: "UJI-K1"},
		{ID: "UJI-AG-2", ClientName: "UJI SUMBER DUA", ChildCount: "2", ClientID: "UJI-K2"},
	}
}

func TestDaftarSumberBisnis(t *testing.T) { // RD BrowseAgentHierarkiList_RD
	u := baru(t)
	u.g.Agen = agenUji()
	if kode, _ := u.panggil("GET", "/sumber-bisnis", pelakuUji{}, nil); kode != http.StatusUnauthorized {
		t.Fatalf("tanpa identitas: %d", kode)
	}
	kode, isi := u.panggil("GET", "/sumber-bisnis", admin, nil)
	if kode != http.StatusOK {
		t.Fatalf("%d %s", kode, isi)
	}
	var b []models.BarisAgen
	if err := json.Unmarshal([]byte(isi), &b); err != nil {
		t.Fatal(err)
	}
	if len(b) != 2 || b[0].ID != "UJI-AG-1" || b[0].ClientName != "UJI SUMBER SATU" || b[1].ChildCount != "2" {
		t.Fatalf("daftar %+v", b)
	}
}

// halamanXOLRetro - isian layar admin lengkap (medan wajib terisi, Approval 1)
// dengan ClaimType yang memunculkan tombol.
func halamanXOLRetro() *models.Halaman {
	h := halamanLengkap("1")
	h.Setel("PolicyTreatyIn.ClaimType", "XOL Retro")
	return h
}

// hasilKlik - jawaban klik baris: nilai yang ditulis PostDT, per jalur halaman.
type hasilKlik struct {
	Nilai map[string]string `json:"nilai"`
}

func (u *uji) klikSOB(id string, p pelakuUji, idAgen string, h *models.Halaman) (int, string) {
	u.t.Helper()
	return u.panggil("POST", "/kasus/"+id+"/pilih-sumber-bisnis", p, map[string]any{"idAgen": idAgen, "halaman": h})
}

func (u *uji) klikSOBNilai(id, idAgen string, h *models.Halaman) map[string]string {
	u.t.Helper()
	kode, isi := u.klikSOB(id, admin, idAgen, h)
	if kode != http.StatusOK {
		u.t.Fatalf("klik %s: %d %s", idAgen, kode, isi)
	}
	var k hasilKlik
	if err := json.Unmarshal([]byte(isi), &k); err != nil {
		u.t.Fatalf("%v: %s", err, isi)
	}
	return k.Nilai
}

// pegang - layar memegang hasil klik di halamannya (state layar).
func pegang(h *models.Halaman, nilai map[string]string) *models.Halaman {
	for j, v := range nilai {
		h.Setel(j, v)
	}
	return h
}

// PostDT langkah 1.1, 1.2, 4, 5 atas baris ChildCount 0 (btnSOB_DT sudah
// menulis btnQuotation "SOB"): ID dan ClientName baris; Leader0 baris (kosong -
// RD efektif NB `LEADER0 IS NULL`); `.Leader1` bukan kolom RD -> "".
// Klik TIDAK menyimpan apa pun: PostDT hanya menulis clipboard.
func TestKlikBarisSumberBisnisTidakMenyimpan(t *testing.T) {
	u := baru(t)
	u.g.Agen = agenUji()
	id := u.buat()
	sebelum := u.g.Halaman[id].Salin()
	panggilSebelum := len(u.g.Panggil)
	nilai := u.klikSOBNilai(id, "UJI-AG-1", halamanXOLRetro())
	harap := map[string]string{
		"Quotation.SourceOfBusiness": "UJI-AG-1",
		"Quotation.SobName":          "UJI SUMBER SATU",
		"Quotation.SobLeader0":       "",
		"Quotation.SobLeader1":       "",
	}
	if len(nilai) != len(harap) {
		t.Fatalf("PostDT menulis tepat empat medan Quotation, dapat %v", nilai)
	}
	for j, v := range harap {
		if got, ada := nilai[j]; !ada || got != v {
			t.Errorf("%s = %q (ada %v), harap %q", j, got, ada, v)
		}
	}
	if slices.Contains(u.g.Panggil[panggilSebelum:], "SimpanHalaman") {
		t.Fatal("klik baris tidak boleh menyimpan halaman")
	}
	s := u.g.Halaman[id]
	for _, j := range []string{"Quotation.SourceOfBusiness", "PolicyTreatyIn.QuotationData.SourceOfBusiness", "PolicyTreatyIn.ClaimType"} {
		if s.Ambil(j) != sebelum.Ambil(j) {
			t.Errorf("tersimpan %s berubah oleh klik: %q", j, s.Ambil(j))
		}
	}
}

// PostDT `@if(.ChildCount > 0, "", ...)`: klik simpul beranak (ChildCount 2)
// mengembalikan keempat medan KOSONG - pilihan yang dipegang layar ikut
// terhapus bila layar menerapkannya. Apa adanya.
func TestKlikSimpulBeranakMengosongkan(t *testing.T) {
	u := baru(t)
	u.g.Agen = agenUji()
	id := u.buat()
	nilai := u.klikSOBNilai(id, "UJI-AG-2", halamanXOLRetro())
	for _, m := range []string{"SourceOfBusiness", "SobName", "SobLeader0", "SobLeader1"} {
		got, ada := nilai["Quotation."+m]
		if !ada || got != "" {
			t.Errorf("ChildCount 2: Quotation.%s = %q (ada %v), harap kosong", m, got, ada)
		}
	}
}

// Tombol hanya di layar admin dan hanya bila ClaimType 'XOL Retro' (isian
// layar); baris yang tidak ada di daftar RD ditolak.
func TestKlikSumberBisnisDitolak(t *testing.T) {
	u := baru(t)
	u.g.Agen = agenUji()
	id := u.buat()
	bukanRetro := halamanXOLRetro()
	bukanRetro.Setel("PolicyTreatyIn.ClaimType", "XOL")
	if kode, isi := u.klikSOB(id, admin, "UJI-AG-1", bukanRetro); kode != http.StatusConflict {
		t.Fatalf("ClaimType bukan XOL Retro: %d %s", kode, isi)
	}
	if kode, isi := u.klikSOB(id, admin, "UJI-TIDAK-ADA", halamanXOLRetro()); kode != http.StatusBadRequest {
		t.Fatalf("agen di luar daftar: %d %s", kode, isi)
	}
	if kode, _ := u.klikSOB(id, admin, "", halamanXOLRetro()); kode != http.StatusBadRequest {
		t.Fatalf("ID kosong: %d", kode)
	}
	if kode, _ := u.klikSOB(id, secHead, "UJI-AG-1", halamanXOLRetro()); kode != http.StatusForbidden {
		t.Fatalf("bukan anggota antrean admin: %d", kode)
	}
	// Berkas di layar Sec Head - layar atasan tidak memuat tombolnya.
	naik := halamanXOLRetro()
	naik.Setel("PolicyTreatyIn.IsApproved", "1")
	if kode, isi := u.kirim(id, admin, naik); kode != http.StatusOK {
		t.Fatalf("admin menyetujui: %d %s", kode, isi)
	}
	if kode, isi := u.klikSOB(id, pelakuUji{"UJI-SH2", models.PosisiSecHead}, "UJI-AG-1", halamanXOLRetro()); kode != http.StatusConflict {
		t.Fatalf("layar atasan: %d %s", kode, isi)
	}
}

func (u *uji) simpan(id string, h *models.Halaman) (int, string) {
	u.t.Helper()
	return u.panggil("PUT", "/kasus/"+id, admin, map[string]any{"halaman": h})
}

// Save membawa hasil klik yang dipegang layar; server menjalankan ulang RD,
// menemukan baris UJI-AG-1 yang PostDT-nya sama persis, lalu menyimpannya.
// Hanya `SourceOfBusiness` berkolom (T_POLIS_QUOTATION.SOURCE_OF_BUSINESS);
// SobName/SobLeader0/SobLeader1 tanpa kolom (katalog: dibuang) - tidak bertahan.
// `PolicyTreatyIn.SOBName` (medan layar, dari view saat pilih bisnis) tidak
// disentuh PostDT.
func TestSaveMenyimpanSumberBisnisYangCocok(t *testing.T) {
	u := baru(t)
	u.g.Agen = agenUji()
	id := u.buat()
	h := halamanXOLRetro()
	pegang(h, u.klikSOBNilai(id, "UJI-AG-1", h))
	if kode, isi := u.simpan(id, h); kode != http.StatusOK {
		t.Fatalf("simpan: %d %s", kode, isi)
	}
	s := u.g.Halaman[id]
	for j, harap := range map[string]string{
		"Quotation.SourceOfBusiness":                    "UJI-AG-1",
		"PolicyTreatyIn.QuotationData.SourceOfBusiness": "UJI-AG-1",
		"Quotation.SobName":                             "", // tanpa kolom
		"PolicyTreatyIn.SOBName":                        "",
	} {
		if got := s.Ambil(j); got != harap {
			t.Errorf("tersimpan %s = %q, harap %q", j, got, harap)
		}
	}
}

// Nilai yang bukan hasil PostDT baris RD yang dijalankan ulang ditolak 422 dan
// tidak satu pun tersimpan - Save maupun Submit:
//   - ID di luar daftar;
//   - ID benar, ClientName palsu (PostDT 1.2 menulis ClientName baris);
//   - ID simpul beranak UJI-AG-2 (PostDT-nya mengosongkan, tidak pernah ID);
//   - SobLeader0 palsu (RD efektif NB `LEADER0 IS NULL` -> selalu "").
func TestNilaiSumberBisnisPalsuDitolak(t *testing.T) {
	for _, palsu := range []map[string]string{
		{"Quotation.SourceOfBusiness": "UJI-KARANGAN", "Quotation.SobName": "UJI KARANGAN"},
		{"Quotation.SourceOfBusiness": "UJI-AG-1", "Quotation.SobName": "UJI NAMA LAIN"},
		{"Quotation.SourceOfBusiness": "UJI-AG-2", "Quotation.SobName": "UJI SUMBER DUA"},
		{"Quotation.SourceOfBusiness": "UJI-AG-1", "Quotation.SobName": "UJI SUMBER SATU", "Quotation.SobLeader0": "UJI-AG-0"},
	} {
		u := baru(t)
		u.g.Agen = agenUji()
		id := u.buat()
		for _, kirim := range []func(*models.Halaman) (int, string){
			func(h *models.Halaman) (int, string) { return u.simpan(id, h) },
			func(h *models.Halaman) (int, string) { return u.kirim(id, admin, h) },
		} {
			kode, isi := kirim(pegang(halamanXOLRetro(), palsu))
			if kode != http.StatusUnprocessableEntity || !strings.Contains(isi, "Source Of Business") ||
				!strings.Contains(isi, "tidak cocok") {
				t.Fatalf("%v: %d %s", palsu, kode, isi)
			}
		}
		if got := u.g.Halaman[id].Ambil("Quotation.SourceOfBusiness"); got != "" {
			t.Fatalf("%v: penolakan tidak boleh menyimpan, tersimpan %q", palsu, got)
		}
		if k := u.g.Kasus[id]; k.PositionNote != models.PosisiAdmin || len(u.g.Riwayat) != 0 {
			t.Fatalf("%v: submit ditolak tidak boleh memindah berkas atau mencatat riwayat", palsu)
		}
	}
}

// ClaimType (isian layar) bukan 'XOL Retro': tombol tidak tampil, medan itu
// terkunci - kiriman diabaikan (pola AC 49-51), nilai tersimpan dipakai.
func TestSumberBisnisTerkunciBilaBukanXOLRetro(t *testing.T) {
	u := baru(t)
	u.g.Agen = agenUji()
	id := u.buat()
	h := halamanXOLRetro()
	pegang(h, u.klikSOBNilai(id, "UJI-AG-1", h))
	if kode, isi := u.simpan(id, h); kode != http.StatusOK {
		t.Fatalf("simpan: %d %s", kode, isi)
	}
	for _, claim := range []string{"XOL", ""} {
		lain := halamanXOLRetro()
		lain.Setel("PolicyTreatyIn.ClaimType", claim)
		pegang(lain, map[string]string{"Quotation.SourceOfBusiness": "UJI-KARANGAN", "Quotation.SobName": "UJI KARANGAN"})
		kode, isi := u.simpan(id, lain)
		if kode != http.StatusOK {
			t.Fatalf("ClaimType %q: %d %s", claim, kode, isi)
		}
		if got := layarDari(t, isi).Halaman.Ambil("Quotation.SourceOfBusiness"); got != "UJI-AG-1" {
			t.Errorf("ClaimType %q: layar %q, harap nilai tersimpan", claim, got)
		}
		if got := u.g.Halaman[id].Ambil("Quotation.SourceOfBusiness"); got != "UJI-AG-1" {
			t.Errorf("ClaimType %q: tersimpan %q, harap tetap UJI-AG-1", claim, got)
		}
	}
}

// Klik simpul beranak lalu Save: pilihan kosong-semua adalah hasil PostDT yang
// sah (ChildCount > 0) - sumber bisnis tersimpan TERHAPUS, apa adanya. Bila RD
// yang dijalankan ulang tidak punya simpul beranak, kosong-semua bukan hasil
// klik mana pun -> 422.
func TestSaveSimpulBeranakMenghapusSumberBisnis(t *testing.T) {
	u := baru(t)
	u.g.Agen = agenUji()
	id := u.buat()
	h := halamanXOLRetro()
	pegang(h, u.klikSOBNilai(id, "UJI-AG-1", h))
	if kode, isi := u.simpan(id, h); kode != http.StatusOK {
		t.Fatalf("simpan: %d %s", kode, isi)
	}
	pegang(h, u.klikSOBNilai(id, "UJI-AG-2", h))
	u.g.Agen = agenUji()[:1] // hierarki tanpa simpul beranak
	if kode, isi := u.simpan(id, h); kode != http.StatusUnprocessableEntity {
		t.Fatalf("kosong tanpa simpul beranak: %d %s", kode, isi)
	}
	u.g.Agen = agenUji()
	u.g.GagalDi = "SimpanHalaman"
	if kode, _ := u.simpan(id, h); kode != http.StatusInternalServerError {
		t.Fatalf("simpan gagal: %d", kode)
	}
	if got := u.g.Halaman[id].Ambil("Quotation.SourceOfBusiness"); got != "UJI-AG-1" {
		t.Fatalf("kegagalan simpan membatalkan semuanya, tersimpan %q", got)
	}
	u.g.GagalDi = ""
	if kode, isi := u.simpan(id, h); kode != http.StatusOK {
		t.Fatalf("simpan: %d %s", kode, isi)
	}
	s := u.g.Halaman[id]
	for _, j := range []string{"Quotation.SourceOfBusiness", "PolicyTreatyIn.QuotationData.SourceOfBusiness"} {
		if got := s.Ambil(j); got != "" {
			t.Errorf("ChildCount 2: %s = %q, harap kosong", j, got)
		}
	}
}

func (u *uji) hitungNetPremi(id string, h *models.Halaman) *models.Halaman {
	u.t.Helper()
	kode, isi := u.panggil("POST", "/kasus/"+id+"/hitung", admin,
		map[string]any{"urutan": []map[string]string{{"aksi": "CountNetPremi"}}, "halaman": h})
	if kode != http.StatusOK {
		u.t.Fatalf("hitung: %d %s", kode, isi)
	}
	return layarDari(u.t, isi).Halaman
}

// `SetPPNPPH` (dipanggil `CountNetPremi_act`) langkah 1-3 membaca STS_PKP agen
// `PolicyTreatyIn.QuotationData.SourceOfBusiness`; langkah 4 berjalan bila
// `.FlagPPH=="true"` atau STS_PKP == 1. Refresh dan Submit memakai pilihan
// yang DIPEGANG layar saat itu (belum tersimpan). XML tidak menjalankan
// refresh sesudah klik: nilai PPH/PPN baru muncul pada refresh berikutnya yang
// dipicu pengguna. Hitung tangan, FlagPPH kosong, TypeTax kosong, Deduction1
// 100: BrokerageFeeSebenarnya = 100; PPHValue = 100 x 2/100 = 2;
// PPNValue = 100 x 2,2/100 = 2,2.
func TestPilihanDipegangMenggerakkanPPNPPH(t *testing.T) {
	u := baru(t)
	u.g.Agen = agenUji()
	u.g.PKPAgen = map[string]string{"UJI-AG-1": "1"}
	id := u.buat()
	h := halamanXOLRetro()
	h.Setel("PolicyTreatyIn.Deduction1", "100")
	angkaSama(t, u.hitungNetPremi(id, h), "PolicyTreatyIn.PPHValue", "0") // tanpa agen PKP
	pegang(h, u.klikSOBNilai(id, "UJI-AG-1", h))
	ly := u.hitungNetPremi(id, h)
	angkaSama(t, ly, "PolicyTreatyIn.BrokerageFeeSebenarnya", "100")
	angkaSama(t, ly, "PolicyTreatyIn.PPHValue", "2")
	angkaSama(t, ly, "PolicyTreatyIn.PPNValue", "2.2")
	if got := ly.Ambil("Quotation.SourceOfBusiness"); got != "UJI-AG-1" {
		t.Fatalf("layar tetap memegang pilihan sesudah refresh, dapat %q", got)
	}
	if got := u.g.Halaman[id].Ambil("Quotation.SourceOfBusiness"); got != "" {
		t.Fatalf("refresh tidak menyimpan, tersimpan %q", got)
	}
	if kode, isi := u.kirim(id, admin, h); kode != http.StatusOK {
		t.Fatalf("submit: %d %s", kode, isi)
	}
	s := u.g.Halaman[id]
	if got := s.Ambil("Quotation.SourceOfBusiness"); got != "UJI-AG-1" {
		t.Fatalf("submit menyimpan pilihan yang dipegang, tersimpan %q", got)
	}
	angkaSama(t, s, "PolicyTreatyIn.PPHValue", "2")
	angkaSama(t, s, "PolicyTreatyIn.PPNValue", "2.2")
	if u.g.Kasus[id].PositionNote != models.PosisiSecHead {
		t.Fatalf("submit admin Approval 1 -> Sec Head, posisi %q", u.g.Kasus[id].PositionNote)
	}
}
