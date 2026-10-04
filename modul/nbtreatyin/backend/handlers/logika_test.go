package handlers_test

// Uji seam 1 - logika dan perhitungan yang dicocokkan ulang ke XML (putaran 2,
// paket P4): komisi OGP saat pilih bisnis, hari tutup buku dari TANGGAL_CLOSING,
// SetPPNPPH membaca status PKP agen, submit tanpa Approval (K6), dan syarat cek
// polis serupa. Nilai harapan dari XML; fixture berawalan UJI-.

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbtreatyin/backend/handlers"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/services"
)

func layarDari(t *testing.T, isi string) services.Layar {
	t.Helper()
	var ly services.Layar
	if err := json.Unmarshal([]byte(isi), &ly); err != nil {
		t.Fatalf("%v: %s", err, isi)
	}
	return ly
}

func angkaSama(t *testing.T, h *models.Halaman, jalur, harap string) {
	t.Helper()
	got, err := models.AngkaTeks(jalur, h.Ambil(jalur))
	if err != nil {
		t.Fatal(err)
	}
	want, _, _ := apd.NewFromString(harap)
	if got.Cmp(want) != 0 {
		t.Errorf("%s = %q, harap %s", jalur, h.Ambil(jalur), harap)
	}
}

// `InputPolicyTreatyInDetail_preACT` langkah 17 -> `TreatyInputPctCommSpreading`
// langkah 2.1.1.1: RiCommOgp = RIONR baris view ber-TREATYID = NoOffer yang
// TREATYTYPE dan TREATYGROUP-nya sama dengan pilihan (RIOGR ditimpa).
func TestPilihBisnisKomisiOgpDariRIONRView(t *testing.T) {
	u := baru(t)
	id := u.buat()
	u.g.Kontrak["UJI-D1"] = models.BarisKontrak{"ID": "UJI-D1", "TREATYID": "UJI-T1", "LIMITCURRENCY": "IDR",
		"CLASSOFBUSINESS": "UJI BISNIS", "PROPORTIONTYPE": models.JenisProporsional,
		"TREATYTYPE": "UJI-QS", "TREATYGROUP": "UJI-GRUP-1", "RIOGR": "30", "RIONR": "27.5"}
	u.g.Kontrak["UJI-D2"] = models.BarisKontrak{"ID": "UJI-D2", "TREATYID": "UJI-T1",
		"TREATYTYPE": "UJI-QS", "TREATYGROUP": "UJI-GRUP-2", "RIOGR": "31", "RIONR": "99"} // grup lain
	u.g.Kontrak["UJI-D3"] = models.BarisKontrak{"ID": "UJI-D3", "TREATYID": "UJI-T2",
		"TREATYTYPE": "UJI-QS", "TREATYGROUP": "UJI-GRUP-1", "RIOGR": "32", "RIONR": "77"} // kontrak lain
	kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D1"})
	if kode != http.StatusOK {
		t.Fatalf("pilih bisnis: %d %s", kode, isi)
	}
	if got := u.g.Halaman[id].Ambil("PolicyTreatyIn.RiCommOgp"); got != "27.5" {
		t.Fatalf("RiCommOgp tersimpan %q, harap 27.5 (RIONR UJI-D1)", got)
	}
}

// Hari tutup buku dibaca dari TANGGAL_CLOSING (gudang), bukan 25 tertanam;
// hari = batas tidak digeser (`InputPolicyTreatyInPre_Act` langkah 9, `>`).
func TestTanggalProduksiDariTanggalClosing(t *testing.T) {
	for _, tt := range []struct {
		closing, hari int
		harap         string
	}{
		{25, 24, "2026-10-24 09:00:00"},
		{25, 25, "2026-10-25 09:00:00"},
		{25, 26, "2026-11-01 09:00:00"},
		{27, 26, "2026-10-26 09:00:00"},
	} {
		u := baru(t)
		u.g.Closing = tt.closing
		saat := time.Date(2026, 10, tt.hari, 9, 0, 0, 0, time.UTC)
		u.s = handlers.Router(services.Baru(u.g, func() time.Time { return saat }), true)
		id := u.buat()
		kode, isi := u.panggil("GET", "/kasus/"+id, admin, nil)
		if kode != http.StatusOK {
			t.Fatalf("buka: %d %s", kode, isi)
		}
		if got := layarDari(t, isi).Halaman.Ambil("PolicyTreatyIn.ProductionDate"); got != tt.harap {
			t.Errorf("closing %d hari %d: ProductionDate %q, harap %q", tt.closing, tt.hari, got, tt.harap)
		}
	}
}

// SetPPNPPH langkah 1-3 (RD BrowseClientName_RD -> STS_PKP) mengendalikan
// langkah 4 lewat refresh layar (sel `.Deduction1` admin: refresh CountOGPONP_Act ->
// CountNetPremi_act langkah 1). FlagPPH kosong: hanya agen PKP yang berjalan.
// 0,025 x (1000 + 600) = 40; 51,1 / 1,022 = 50; PPH 50 x 0,02 = 1; PPN 50 x 0,022 = 1,1.
func TestSetPPNPPHMembacaStatusPKPAgen(t *testing.T) {
	for _, tt := range []struct {
		pkp           string
		pph, ppn, fee string
	}{
		{"1", "1", "1.1", "40"},
		{"0", "0", "0", "0"},
	} {
		u := baru(t)
		u.g.StsPKP = tt.pkp
		id := u.buat()
		h := halamanLengkap("")
		h.Setel("PolicyTreatyIn.PremiOnp", "600")
		h.Setel("PolicyTreatyIn.Deduction1", "51.1")
		h.Setel("PolicyTreatyIn.TypeTax", models.TypeTaxInclusive)
		h.Setel("PolicyTreatyIn.FlagPPH", "")
		kode, isi := u.panggil("POST", "/kasus/"+id+"/hitung", admin, map[string]any{"urutan": []map[string]string{{"aksi": "CountOGPONP"}}, "halaman": h})
		if kode != http.StatusOK {
			t.Fatalf("hitung: %d %s", kode, isi)
		}
		ly := layarDari(t, isi)
		angkaSama(t, ly.Halaman, "PolicyTreatyIn.PPHValue", tt.pph)
		angkaSama(t, ly.Halaman, "PolicyTreatyIn.PPNValue", tt.ppn)
		angkaSama(t, ly.Halaman, "PolicyTreatyIn.BrokerageFee", tt.fee)
	}
}

// K6 (03-10-2026): submit tanpa Approval ditolak 422 di KETIGA layar - tombol
// Submit hanya tampil untuk IsApproved 1/0 dan `ListSuggest` mewajibkan Approval.
func TestSubmitTanpaApprovalDitolakDiSetiapJenjang(t *testing.T) {
	u := baru(t)
	id := u.buat()
	tolak := func(p pelakuUji, h *models.Halaman, posisi string, riwayat int) {
		t.Helper()
		kode, isi := u.kirim(id, p, h)
		if kode != http.StatusUnprocessableEntity || !strings.Contains(isi, "Approval") {
			t.Fatalf("%s tanpa Approval: %d %s", p.akun, kode, isi)
		}
		if u.g.Kasus[id].PositionNote != posisi || len(u.g.Riwayat) != riwayat {
			t.Fatalf("%s: submit ditolak tidak boleh memindah berkas atau mencatat riwayat", p.akun)
		}
	}
	tolak(admin, halamanLengkap(""), models.PosisiAdmin, 0)
	if kode, isi := u.kirim(id, admin, halamanLengkap("1")); kode != http.StatusOK {
		t.Fatalf("admin menyetujui: %d %s", kode, isi)
	}
	tolak(secHead, putusan(""), models.PosisiSecHead, 1)
	if kode, isi := u.kirim(id, secHead, putusan("1")); kode != http.StatusOK {
		t.Fatalf("Sec Head menyetujui: %d %s", kode, isi)
	}
	tolak(deptHead, putusan(""), models.PosisiDeptHead, 2)
}

// `InputPolicyTreatyInPost_Act` langkah 3 bersyarat `.PolicyTreatyIn.IsApproved==1`,
// dan hanya pasca-proses flow action admin yang memanggilnya: admin menolak
// dan putusan atasan tidak diperiksa. Pesan langkah 5.1 merangkai SETIAP nopolis.
func TestCekPolisSerupaHanyaSaatAdminMenyetujui(t *testing.T) {
	u := baru(t)
	id := u.buat()
	u.g.Serupa = []string{"UJI-NOPOL-1", "UJI-NOPOL-2"}
	kode, isi := u.kirim(id, admin, halamanLengkap("1"))
	if kode != http.StatusUnprocessableEntity ||
		!strings.Contains(isi, "Protect Duplicate Policy; data is similar to UJI-NOPOL-1 UJI-NOPOL-2 ") {
		t.Fatalf("admin menyetujui: %d %s", kode, isi)
	}
	u.g.Serupa = nil
	if kode, isi := u.kirim(id, admin, halamanLengkap("1")); kode != http.StatusOK {
		t.Fatalf("tanpa polis serupa: %d %s", kode, isi)
	}
	u.g.Serupa = []string{"UJI-NOPOL-1"}
	if kode, isi := u.kirim(id, secHead, putusan("1")); kode != http.StatusOK || u.g.Kasus[id].PositionNote != models.PosisiDeptHead {
		t.Fatalf("putusan atasan tidak menjalankan cek: %d %s", kode, isi)
	}
	id2 := u.buat()
	if kode, isi := u.kirim(id2, admin, halamanLengkap("0")); kode != http.StatusOK || u.g.Kasus[id2].StatusWork != models.StatusDitolak {
		t.Fatalf("admin menolak (IsApproved 0) tidak menjalankan cek: %d %s", kode, isi)
	}
}
