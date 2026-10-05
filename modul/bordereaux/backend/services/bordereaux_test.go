package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/bordereaux/backend/models"
	"nusantarare/modul/bordereaux/backend/services"
	"nusantarare/modul/bordereaux/backend/tiruan"
)

var ctx = context.Background()

var jamUji = time.Date(2026, 10, 4, 14, 5, 7, 123*int(time.Millisecond), time.FixedZone("WIB", 7*3600))

func layanan(g *tiruan.Gudang) *services.Layanan {
	return services.BaruLayanan(g, tiruan.Transaksi).DenganJam(func() time.Time { return jamUji })
}

var (
	maker   = services.Aktor{AkunID: "UJI-MAKER", Penuh: true}
	checker = services.Aktor{AkunID: "UJI-CHK", Peran: []string{models.PeranChecker}}
	spv     = services.Aktor{AkunID: "UJI-SPV", Peran: []string{models.PeranSupervisor}}
	super   = services.Aktor{AkunID: "UJI-SUPER", Superadmin: true, Penuh: true}
	tamu    = services.Aktor{AkunID: "UJI-TAMU"}
)

func kombinasi(t *testing.T, tipe, bisnis string) models.KombinasiBdx {
	t.Helper()
	k, err := models.CariKombinasi(tipe, bisnis)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

// csvUji - header + baris; sel diisi menurut nama kolom, sisanya kosong.
func csvUji(k models.KombinasiBdx, baris ...map[string]string) string {
	var b strings.Builder
	for i, c := range k.Kolom {
		if i > 0 {
			b.WriteString(";")
		}
		b.WriteString(c.Judul)
	}
	b.WriteString("\r\n")
	for _, isi := range baris {
		for i, c := range k.Kolom {
			if i > 0 {
				b.WriteString(";")
			}
			b.WriteString(isi[c.Kolom])
		}
		b.WriteString("\r\n")
	}
	return b.String()
}

func TestKombinasi29CocokTemplatDanPunyaRingkasan(t *testing.T) {
	if len(models.Kombinasi) != 29 {
		t.Fatalf("kombinasi %d", len(models.Kombinasi))
	}
	if b := models.DaftarBusiness(models.TypeSubrogation); len(b) != 1 || b[0] != models.BusinessBonding {
		t.Errorf("subrogation %v", b)
	}
	if b := models.DaftarBusiness(models.TypePremium); len(b) != 14 {
		t.Errorf("premium %d business", len(b))
	}
	tanpaMataUang := []string{}
	for _, k := range models.Kombinasi {
		mataUang := false
		for _, c := range k.Kolom {
			mataUang = mataUang || c.Kolom == "CURRENCY"
		}
		if !mataUang {
			tanpaMataUang = append(tanpaMataUang, k.Kode)
		}
		// Setiap kombinasi punya kolom Reinsurer dan RNM yang dijumlah tab Summary.
		r := services.HitungRingkasan(k, []models.Baris{{"CURRENCY": "IDR"}})
		if len(r) != 1 {
			t.Errorf("%s: tanpa kolom ringkasan", k.Kode)
		}
	}
	// Templat Premi Credit memang tanpa CURRENCY (Pega: Summary satu baris bermata uang kosong).
	if strings.Join(tanpaMataUang, ",") != "premi.credit" {
		t.Errorf("tanpa CURRENCY: %v", tanpaMataUang)
	}
}

func TestAngkaDanTanggalCSV(t *testing.T) {
	for masuk, mau := range map[string]string{
		"1.234.567,89": "1234567.89", "25%": "25", "2,5 %": "2.5", "-": "", "": "", "007": "7", "-1.000": "-1000", "12,500": "12.5",
	} {
		if v, ok := services.AngkaIndonesia(masuk); !ok || v != mau {
			t.Errorf("angka %q = %q %v, mau %q", masuk, v, ok, mau)
		}
	}
	for _, salah := range []string{"N/A", "1,2,3", "abc", "1e5"} {
		if _, ok := services.AngkaIndonesia(salah); ok {
			t.Errorf("angka %q harus ditolak", salah)
		}
	}
	for masuk, mau := range map[string]string{"31/12/2024": "31-12-2024", "1/2/2024": "01-02-2024", "05-06-2025": "05-06-2025", "": ""} {
		if v, ok := services.TanggalCSV(masuk); !ok || v != mau {
			t.Errorf("tanggal %q = %q %v, mau %q", masuk, v, ok, mau)
		}
	}
	for _, salah := range []string{"31/02/2024", "2024-12-31", "13/13/2024", "12/2024"} {
		if _, ok := services.TanggalCSV(salah); ok {
			t.Errorf("tanggal %q harus ditolak", salah)
		}
	}
}

func TestUnggahCSVPetakanDanValidasi(t *testing.T) {
	k := kombinasi(t, "PREMIUM", "FIRE")
	l := layanan(tiruan.Contoh())
	isi := csvUji(k,
		map[string]string{"COB": "UJI", "POI_START": "01/01/2026", "POI_END": "31/12/2026", "CURRENCY": "IDR", "PREMIUM_REINSURER": "1.000,50", "PREMIUM_RNM": "250,25"},
		map[string]string{"COB": "UJI", "CURRENCY": "usd", "PREMIUM_REINSURER": "10", "PREMIUM_RNM": "2"},
		map[string]string{},
		map[string]string{"COB": "UJI", "CURRENCY": "IDR", "PREMIUM_REINSURER": "999,5", "PREMIUM_RNM": "0,75"},
	)
	p, err := l.UnggahCSV(ctx, "premium", "fire", string(rune(0xFEFF))+isi)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Baris) != 3 || p.Baris[0]["POI_START"] != "01-01-2026" || p.Baris[0]["PREMIUM_REINSURER"] != "1000.5" {
		t.Fatalf("baris %+v", p.Baris)
	}
	if len(p.Ringkasan) != 2 || p.Ringkasan[0] != (models.Ringkasan{Currency: "IDR", Reinsurer: "2000.0", RNM: "251.00"}) ||
		p.Ringkasan[1].Currency != "USD" {
		t.Errorf("ringkasan %+v", p.Ringkasan)
	}
	salah := csvUji(k,
		map[string]string{"POI_START": "31/02/2026", "PREMIUM_RNM": "N/A"},
		map[string]string{"POI_START": "01/12/2026", "POI_END": "01/01/2026"},
	)
	_, err = l.UnggahCSV(ctx, "PREMIUM", "FIRE", salah)
	if !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Fatalf("harus ditolak: %v", err)
	}
	for _, mau := range []string{"Row 2, column PERIOD OF INSURANCE START", "Row 2, column *Premium (RNM Share)", "Row 3: PERIOD OF INSURANCE END is before PERIOD OF INSURANCE START"} {
		if !strings.Contains(err.Error(), mau) {
			t.Errorf("pesan tanpa %q:\n%v", mau, err)
		}
	}
	_, err = l.UnggahCSV(ctx, "PREMIUM", "ENGINEERING", isi)
	if err == nil || !strings.Contains(err.Error(), "the PREMIUM ENGINEERING template needs 61 columns") {
		t.Errorf("berkas Fire ke Engineering: %v", err)
	}
	if _, err := l.UnggahCSV(ctx, "SUBROGATION", "FIRE", csvUji(kombinasi(t, "SUBROGATION", "BONDING"), map[string]string{"COB": "UJI"})); err != nil {
		t.Errorf("subrogation dipaksa BONDING: %v", err)
	}
}

func TestIDMengikutiXML(t *testing.T) {
	if got := services.IDBerkas(jamUji); got != "BDX-2026.10.04.07123" {
		t.Errorf("BDX %s", got)
	}
	if got := services.IDDetail(jamUji); got != "BDX_DTL-20261004020507123" {
		t.Errorf("detail %s (hh 12 jam seperti XML)", got)
	}
}

func TestHakAtas(t *testing.T) {
	diPembuat := models.Header{UserInput: "UJI-MAKER", Position: "UJI-MAKER"}
	if h := services.HakAtas(maker, diPembuat); !h.Ubah || !h.Hapus || !h.Submit || h.Putuskan {
		t.Errorf("pembuat %+v", h)
	}
	if h := services.HakAtas(super, diPembuat); !h.Ubah {
		t.Errorf("superadmin mengambil alih %+v", h)
	}
	// View only berlaku juga untuk superadmin (keputusan work owner 05-10-2026).
	superLihat := services.Aktor{AkunID: "UJI-SUPER", Superadmin: true}
	if h := services.HakAtas(superLihat, diPembuat); h != (services.Hak{}) || superLihat.BolehBuat() {
		t.Errorf("superadmin View only %+v", h)
	}
	if services.DaftarPilihan(superLihat).CopyOld {
		t.Error("superadmin View only tanpa tombol Copy Old Data")
	}
	if h := services.HakAtas(tamu, diPembuat); h != (services.Hak{}) {
		t.Errorf("orang lain %+v", h)
	}
	// Pembuat yang menunya kini View only: tidak lagi Edit/Delete/Submit (hak menu PENUH, 04-10-2026).
	if h := services.HakAtas(services.Aktor{AkunID: "UJI-MAKER"}, diPembuat); h != (services.Hak{}) {
		t.Errorf("pembuat View only %+v", h)
	}
	diChecker := models.Header{UserInput: "UJI-MAKER", Position: models.PosisiChecker, StatusAksep: models.StatusAccept}
	if !services.HakAtas(checker, diChecker).Putuskan || services.HakAtas(maker, diChecker).Ubah {
		t.Error("checker memutuskan, pembuat tidak lagi mengubah")
	}
	pembuatChecker := services.Aktor{AkunID: "UJI-MAKER", Peran: []string{models.PeranChecker}}
	if services.HakAtas(pembuatChecker, diChecker).Putuskan {
		t.Error("pembuat tidak boleh menjadi Checker berkasnya sendiri")
	}
	diSpv := models.Header{UserInput: "UJI-SPV", Position: models.PosisiSupervisor}
	if !services.HakAtas(spv, diSpv).Putuskan {
		t.Error("supervisor boleh siapa saja pemegang workbasket, juga pembuatnya")
	}
	if services.HakAtas(super, models.Header{UserInput: "X", Position: "X", StatusAksep: models.StatusResolveComplete}) != (services.Hak{}) {
		t.Error("Resolve-Complete hanya View")
	}
}

func simpanBaru(t *testing.T, l *services.Layanan, a services.Aktor, baris ...models.Baris) string {
	t.Helper()
	id, err := l.Simpan(ctx, a, services.PermintaanSimpan{Type: "PREMIUM", Business: "FIRE", MasterID: "UJI-TI-1",
		ReportStart: "01-07-2026", ReportEnd: "30-09-2026", ReffNoSOA: "UJI-SOA", Baris: baris})
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func TestSimpanTanpaJSONDanHak(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	if _, err := l.Simpan(ctx, tamu, services.PermintaanSimpan{Type: "PREMIUM", Business: "FIRE", MasterID: "UJI-TI-1",
		ReportStart: "01-07-2026", ReportEnd: "30-09-2026"}); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("menu View only: %v", err)
	}
	baris := models.Baris{"COB": "UJI", "CURRENCY": "IDR", "PREMIUM_RNM": "10.5"}
	id := simpanBaru(t, l, maker, baris, baris)
	if id != "BDX-2026.10.04.07123" {
		t.Errorf("id %s", id)
	}
	h := g.Header[id]
	if h.UserInput != "UJI-MAKER" || h.Position != "UJI-MAKER" || h.CedingName != "UJI CEDANT SATU" || h.TreatyName != "UJI KONTRAK SATU" || h.SobName != "UJI SOB" {
		t.Errorf("header disalin dari TREATY_IN: %+v", h)
	}
	d := g.Rinci["BORDEREAUX_PREMI_FIRE"][id]
	if len(d) != 2 || d[0][models.KolomID] != "BDX_DTL-20261004020507123" || d[1][models.KolomID] != "BDX_DTL-20261004020507124" {
		t.Errorf("detail langsung ke tabel (bukan menunggu Resolve-Complete): %+v", d)
	}
	// Save kedua di detik yang sama: ID berkas bentrok -> milidetik berikutnya.
	if id2 := simpanBaru(t, l, maker); id2 != "BDX-2026.10.04.07124" {
		t.Errorf("bentrok %s", id2)
	}
	// Ubah Type: detail lama di tabel lain ikut terhapus.
	_, err := l.Simpan(ctx, maker, services.PermintaanSimpan{BdxID: id, Type: "CLAIM", Business: "FIRE", MasterID: "UJI-TI-2",
		ReportStart: "01-07-2026", ReportEnd: "30-09-2026", Baris: []models.Baris{{"COB": "UJI"}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Rinci["BORDEREAUX_PREMI_FIRE"][id]) != 0 || len(g.Rinci["BORDEREAUX_CLAIM_FIRE"][id]) != 1 || g.Header[id].MasterID != "UJI-TI-2" {
		t.Errorf("ganti Type: %+v", g.Header[id])
	}
	_, err = l.Simpan(ctx, maker, services.PermintaanSimpan{Type: "PREMIUM", Business: "FIRE", MasterID: "UJI-TI-X",
		ReportStart: "30-09-2026", ReportEnd: "01-07-2026", Baris: []models.Baris{{"PREMIUM_RNM": "abc", "TIDAK_ADA": "1"}}})
	if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), "Master Treaty UJI-TI-X is not found") ||
		!strings.Contains(err.Error(), "End is before") {
		t.Errorf("validasi header: %v", err)
	}
}

func TestAlurPersetujuan(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	id := simpanBaru(t, l, maker)
	if err := l.Submit(ctx, maker, id, true, ""); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("submit tanpa detail: %v", err)
	}
	simpan := func() {
		if _, err := l.Simpan(ctx, maker, services.PermintaanSimpan{BdxID: id, Type: "PREMIUM", Business: "FIRE", MasterID: "UJI-TI-1",
			ReportStart: "01-07-2026", ReportEnd: "30-09-2026", Baris: []models.Baris{{"COB": "UJI"}}}); err != nil {
			t.Fatal(err)
		}
	}
	simpan()
	if err := l.Submit(ctx, maker, id, true, "siap"); err != nil {
		t.Fatal(err)
	}
	if h := g.Header[id]; h.Position != models.PosisiChecker || h.StatusAksep != models.StatusAccept {
		t.Fatalf("ke checker %+v", h)
	}
	if err := l.Submit(ctx, checker, id, false, ""); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("tolak tanpa komentar: %v", err)
	}
	if err := l.Submit(ctx, checker, id, false, "periksa ulang"); err != nil {
		t.Fatal(err)
	}
	if h := g.Header[id]; h.Position != "UJI-MAKER" || h.StatusAksep != models.StatusRejected {
		t.Fatalf("kembali ke pembuat %+v", h)
	}
	simpan()
	_ = l.Submit(ctx, maker, id, true, "")
	if err := l.Submit(ctx, spv, id, true, ""); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("supervisor sebelum checker: %v", err)
	}
	if err := l.Submit(ctx, checker, id, true, ""); err != nil {
		t.Fatal(err)
	}
	if err := l.Submit(ctx, spv, id, true, "ok"); err != nil {
		t.Fatal(err)
	}
	if h := g.Header[id]; h.Position != "" || h.StatusAksep != models.StatusResolveComplete {
		t.Fatalf("selesai %+v", h)
	}
	r, _ := l.Buka(ctx, super, id)
	if len(r.Riwayat) != 5 || r.Riwayat[1].IsApproved || r.Riwayat[1].Komentar != "periksa ulang" || r.Riwayat[4].PIC != "UJI-SPV" {
		t.Errorf("riwayat %+v", r.Riwayat)
	}
	if r.Hak != (services.Hak{}) || len(r.Baris) != 1 {
		t.Errorf("selesai hanya View: %+v", r)
	}
	if err := l.Hapus(ctx, maker, id); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("hapus berkas selesai: %v", err)
	}
}

func TestHapusDanSudahDiproses(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	id := simpanBaru(t, l, maker, models.Baris{"COB": "UJI"})
	if err := l.Hapus(ctx, tamu, id); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("orang lain: %v", err)
	}
	if err := l.Hapus(ctx, maker, id); err != nil {
		t.Fatal(err)
	}
	if _, ok := g.Header[id]; ok || len(g.Rinci["BORDEREAUX_PREMI_FIRE"][id]) != 0 {
		t.Error("berkas dan detail terhapus")
	}
	if _, err := l.Buka(ctx, maker, id); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("buka sesudah hapus: %v", err)
	}
	// Dua checker serentak: yang kedua mendapati posisi sudah berubah.
	h := g.Header["BDX-UJI.1"]
	h.Position = models.PosisiSupervisor
	g.Header["BDX-UJI.1"] = h
	if err := l.Submit(ctx, checker, "BDX-UJI.1", true, ""); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("posisi berubah: %v", err)
	}
}

func TestDaftarDanPilihan(t *testing.T) {
	l := layanan(tiruan.Contoh())
	h, err := l.Daftar(ctx, checker, models.Filter{Ceding: "satu"}, 1, 0)
	if err != nil || h.Total != 1 || !h.Daftar[0].Hak.Putuskan || h.Ukuran != services.UkuranHalaman {
		t.Errorf("daftar %+v %v", h, err)
	}
	if _, err := l.Daftar(ctx, checker, models.Filter{ReportStart: "2026-01-01"}, 1, 0); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("tanggal filter salah bentuk: %v", err)
	}
	p := services.DaftarPilihan(maker)
	if !p.BolehBuat || len(p.Kolom) != 29 || len(p.Kolom["PREMIUM|LIABILITY"]) != 41 || p.Kolom["PREMIUM|LIABILITY"][40].Kolom != "NOTE" {
		t.Errorf("pilihan %+v", p.Business)
	}
	if services.DaftarPilihan(checker).BolehBuat {
		t.Error("checker tidak boleh Input Data")
	}
}
