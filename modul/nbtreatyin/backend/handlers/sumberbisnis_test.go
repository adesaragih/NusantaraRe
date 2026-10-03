package handlers_test

// Uji seam 1 pemilih Source Of Business (XOL Retro) - tombol `Select Source Of
// Business` `Section/DetailPolicyTreatyIn` -> Harness `SOB` -> klik baris
// TreeGrid `SourceHierarki` (pra-proses `SearchHierarkiSourceBizAgent_PostDT`).
// Helper (`baru`, `halamanLengkap`, pelaku) milik alur_test.go.

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/services"
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

// halamanXOLRetro - isian layar admin dengan ClaimType yang memunculkan tombol.
func halamanXOLRetro() *models.Halaman {
	h := halamanLengkap("")
	h.Setel("PolicyTreatyIn.ClaimType", "XOL Retro")
	return h
}

func (u *uji) pilihSOB(id string, p pelakuUji, idAgen string, h *models.Halaman) (int, string) {
	u.t.Helper()
	return u.panggil("POST", "/kasus/"+id+"/pilih-sumber-bisnis", p, map[string]any{"idAgen": idAgen, "halaman": h})
}

func TestPilihSumberBisnisMenyimpanQuotation(t *testing.T) { // SearchHierarkiSourceBizAgent_PostDT
	u := baru(t)
	u.g.Agen = agenUji()
	id := u.buat()
	u.g.Halaman[id].Setel("PolicyTreatyIn.SOBName", "UJI SOB KONTRAK") // dari view saat pilih bisnis
	kode, isi := u.pilihSOB(id, admin, "UJI-AG-1", halamanXOLRetro())
	if kode != http.StatusOK {
		t.Fatalf("%d %s", kode, isi)
	}
	var ly services.Layar
	if err := json.Unmarshal([]byte(isi), &ly); err != nil {
		t.Fatal(err)
	}
	for j, harap := range map[string]string{
		"Quotation.SourceOfBusiness": "UJI-AG-1",
		"Quotation.SobName":          "UJI SUMBER SATU",
		"Quotation.SobLeader0":       "",
		"Quotation.SobLeader1":       "",
		"PolicyTreatyIn.SOBName":     "UJI SOB KONTRAK", // medan layar TIDAK disentuh PostDT
	} {
		if got := ly.Halaman.Ambil(j); got != harap {
			t.Errorf("layar %s = %q, harap %q", j, got, harap)
		}
	}
	// Tersimpan satu kali di T_POLIS_QUOTATION: nilai QuotationData didahulukan
	// repository, jadi keduanya harus membawa pilihan baru.
	s := u.g.Halaman[id]
	for _, j := range []string{"Quotation.SourceOfBusiness", "PolicyTreatyIn.QuotationData.SourceOfBusiness"} {
		if got := s.Ambil(j); got != "UJI-AG-1" {
			t.Errorf("tersimpan %s = %q", j, got)
		}
	}
	if s.Ambil("PolicyTreatyIn.ClaimType") != "XOL Retro" {
		t.Fatal("isian layar ikut terkirim (pySubmitData=Yes) dan tersimpan")
	}
}

// Tombol hanya di layar admin dan hanya bila ClaimType 'XOL Retro'; baris yang
// tidak ada di daftar RD ditolak; nol penyimpanan pada setiap penolakan.
func TestPilihSumberBisnisDitolak(t *testing.T) {
	u := baru(t)
	u.g.Agen = agenUji()
	id := u.buat()
	bukanRetro := halamanXOLRetro()
	bukanRetro.Setel("PolicyTreatyIn.ClaimType", "XOL")
	if kode, isi := u.pilihSOB(id, admin, "UJI-AG-1", bukanRetro); kode != http.StatusConflict {
		t.Fatalf("ClaimType bukan XOL Retro: %d %s", kode, isi)
	}
	if kode, isi := u.pilihSOB(id, admin, "UJI-TIDAK-ADA", halamanXOLRetro()); kode != http.StatusBadRequest {
		t.Fatalf("agen di luar daftar: %d %s", kode, isi)
	}
	if kode, _ := u.pilihSOB(id, admin, "", halamanXOLRetro()); kode != http.StatusBadRequest {
		t.Fatalf("ID kosong: %d", kode)
	}
	if kode, _ := u.pilihSOB(id, secHead, "UJI-AG-1", halamanXOLRetro()); kode != http.StatusForbidden {
		t.Fatalf("bukan anggota antrean admin: %d", kode)
	}
	if got := u.g.Halaman[id].Ambil("Quotation.SourceOfBusiness"); got != "" {
		t.Fatalf("penolakan tidak boleh menyimpan, tersimpan %q", got)
	}
	// Berkas di layar Sec Head - ClaimType tersimpan 'XOL Retro' pun tombolnya
	// tidak ada (layar atasan tidak memuatnya).
	naik := halamanXOLRetro()
	naik.Setel("PolicyTreatyIn.IsApproved", "1")
	if kode, isi := u.kirim(id, admin, naik); kode != http.StatusOK {
		t.Fatalf("admin menyetujui: %d %s", kode, isi)
	}
	if kode, isi := u.pilihSOB(id, pelakuUji{"UJI-SH2", models.PosisiSecHead}, "UJI-AG-1", halamanXOLRetro()); kode != http.StatusConflict {
		t.Fatalf("layar atasan: %d %s", kode, isi)
	}
}

// Simpul beranak (ChildCount > 0) MENGOSONGKAN sumber bisnis yang sudah
// terpilih - PostDT apa adanya - dan kegagalan simpan membatalkan semuanya.
func TestPilihSumberBisnisBeranakDanPembatalan(t *testing.T) {
	u := baru(t)
	u.g.Agen = agenUji()
	id := u.buat()
	if kode, isi := u.pilihSOB(id, admin, "UJI-AG-1", halamanXOLRetro()); kode != http.StatusOK {
		t.Fatalf("%d %s", kode, isi)
	}
	u.g.GagalDi = "SimpanHalaman"
	if kode, _ := u.pilihSOB(id, admin, "UJI-AG-2", halamanXOLRetro()); kode != http.StatusInternalServerError {
		t.Fatalf("simpan gagal: %d", kode)
	}
	if got := u.g.Halaman[id].Ambil("Quotation.SourceOfBusiness"); got != "UJI-AG-1" {
		t.Fatalf("kegagalan simpan harus membatalkan, tersimpan %q", got)
	}
	u.g.GagalDi = ""
	if kode, isi := u.pilihSOB(id, admin, "UJI-AG-2", halamanXOLRetro()); kode != http.StatusOK {
		t.Fatalf("%d %s", kode, isi)
	}
	s := u.g.Halaman[id]
	for _, j := range []string{"Quotation.SourceOfBusiness", "Quotation.SobName", "PolicyTreatyIn.QuotationData.SourceOfBusiness"} {
		if got := s.Ambil(j); got != "" {
			t.Errorf("ChildCount 2: %s = %q, harap kosong", j, got)
		}
	}
}

// Akibat pada data kasus: `SetPPNPPH` langkah 1-3 membaca STS_PKP agen
// `PolicyTreatyIn.QuotationData.SourceOfBusiness`. Sesudah memilih agen PKP,
// Save berikutnya menghitung PPH/PPN (langkah 4) - FlagPPH tidak dicentang.
// Hitung tangan, Deduction1 100, TypeTax kosong: BrokerageFeeSebenarnya = 100,
// PPHValue = 100 x 2/100 = 2, PPNValue = 100 x 2.2/100 = 2.2.
func TestSumberBisnisPKPMenggerakkanPPNPPH(t *testing.T) {
	u := baru(t)
	u.g.Agen = agenUji()
	u.g.PKPAgen = map[string]string{"UJI-AG-1": "1"}
	id := u.buat()
	h := halamanXOLRetro()
	h.Setel("PolicyTreatyIn.IsApproved", "1")
	h.Setel("PolicyTreatyIn.Deduction1", "100")
	angka := func(j string) *apd.Decimal {
		t.Helper()
		d, err := models.AngkaTeks(j, u.g.Halaman[id].Ambil(j))
		if err != nil {
			t.Fatal(err)
		}
		return d
	}
	if kode, isi := u.panggil("PUT", "/kasus/"+id, admin, map[string]any{"halaman": h}); kode != http.StatusOK {
		t.Fatalf("simpan: %d %s", kode, isi)
	}
	if angka("PolicyTreatyIn.PPHValue").Sign() != 0 {
		t.Fatal("tanpa sumber bisnis PKP dan FlagPPH, PPH 0")
	}
	if kode, isi := u.pilihSOB(id, admin, "UJI-AG-1", h); kode != http.StatusOK {
		t.Fatalf("pilih: %d %s", kode, isi)
	}
	if kode, isi := u.panggil("PUT", "/kasus/"+id, admin, map[string]any{"halaman": h}); kode != http.StatusOK {
		t.Fatalf("simpan: %d %s", kode, isi)
	}
	for j, harap := range map[string]*apd.Decimal{
		"PolicyTreatyIn.BrokerageFeeSebenarnya": apd.New(100, 0),
		"PolicyTreatyIn.PPHValue":               apd.New(2, 0),
		"PolicyTreatyIn.PPNValue":               apd.New(22, -1),
	} {
		if got := angka(j); got.Cmp(harap) != 0 {
			t.Errorf("%s = %s, harap %s", j, got, harap)
		}
	}
}
