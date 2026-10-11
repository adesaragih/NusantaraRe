package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/ririsklife/backend/models"
	"nusantarare/modul/ririsklife/backend/repository"
	"nusantarare/modul/ririsklife/backend/services"
	"nusantarare/modul/ririsklife/backend/tiruan"
)

var ctx = context.Background()

var (
	penuh = services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
	lihat = services.Aktor{AkunID: "UJI-LIHAT"}
	jam   = time.Date(2026, 10, 8, 3, 4, 5, 600e6, time.UTC)
)

func layanan(g *tiruan.Gudang) *services.Layanan {
	return services.BaruLayanan(g, tiruan.Transaksi{G: g}.Jalankan).DenganJam(func() time.Time { return jam })
}

// Situs tidak tepat satu baris = galat, nol tulisan (ringkasan baru lewat upload membutuhkan site); nomor sequence
// rincian 6 angka = galat (LPAD 5 akan memotongnya).
func TestSitusDanNomorRincian(t *testing.T) {
	for _, situs := range [][]string{nil, {"1", "2"}} {
		g := tiruan.Contoh()
		g.SitusAktif = situs
		if _, err := layanan(g).SimpanUnggah(ctx, penuh, services.PermintaanUnggah{CSV: kepala + "UJI X;1;1;;1\n"}); !errors.Is(err, repository.ErrSitus) {
			t.Errorf("%v: %v", situs, err)
		}
		if len(g.Ringkasan) != 2 || len(g.Rincian) != 3 {
			t.Error("tertulis tanpa situs")
		}
	}
	g := tiruan.Contoh()
	g.SeqRincian = 100000
	if _, err := layanan(g).SimpanRincian(ctx, penuh, "1000003", "", models.IsianRincian{Contract: "9", Year: "9", Risk: "1"}); !errors.Is(err, models.ErrIDTidakSah) {
		t.Errorf("nomor 6 angka: %v", err)
	}
	// Rincian TIDAK membaca site: site rusak tidak menghalangi Save detail.
	g = tiruan.Contoh()
	g.SitusAktif = nil
	if k, err := layanan(g).SimpanRincian(ctx, penuh, "1000003", "", models.IsianRincian{Contract: "9", Year: "9", Risk: "1"}); err != nil || k.ID != "131723" {
		t.Errorf("rincian tanpa site %+v %v", k, err)
	}
}

// Add (keputusan work owner 08-10-2026): nama dipangkas; ID = site || LPAD(M_RIRISK_LIFE_SUMMARY_SEQ, 6) - 1000004
// terpakai dilewati; OPERATORID akun; MODIFIEDDATE format Pega.
func TestAdd(t *testing.T) {
	g := tiruan.Contoh()
	r, err := layanan(g).Simpan(ctx, penuh, "", models.Isian{UsedBy: "  UJI RISK BARU "})
	if err != nil {
		t.Fatal(err)
	}
	mau := models.Ringkasan{ID: "1000005", UsedBy: "UJI RISK BARU", OperatorID: "UJI-ADMIN", ModifiedDate: "20261008T030405.600 GMT"}
	if g.Ringkasan["1000005"] != mau || r.ID != "1000005" || r.Diubah != "08-10-2026" {
		t.Errorf("tersimpan %+v / jawab %+v", g.Ringkasan["1000005"], r)
	}
}

// Nama wajib, dipangkas, tidak kembar tanpa beda huruf; ID tak ada 404.
func TestValidasiNama(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	for isi, mau := range map[string]string{
		"  ":                                    "R/I RISK NAME is required",
		" uji risk retro ":                      "R/I RISK NAME uji risk retro is already used by ID 1000003",
		strings.Repeat("A", models.BatasNama+1): "longer than 200",
	} {
		if _, err := l.Simpan(ctx, penuh, "", models.Isian{UsedBy: isi}); !errors.Is(err, services.ErrMasukanTidakSah) ||
			!strings.Contains(err.Error(), mau) {
			t.Errorf("%q: %v", isi, err)
		}
	}
	if _, err := l.Simpan(ctx, penuh, "1000004", models.Isian{UsedBy: "Uji Risk Retro"}); err == nil {
		t.Error("edit ke nama baris lain harus ditolak")
	}
	if _, err := l.Simpan(ctx, penuh, "999", models.Isian{UsedBy: "UJI X"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("edit ID tak ada: %v", err)
	}
	if len(g.Ringkasan) != 2 {
		t.Error("ditolak tetapi tertulis")
	}
}

// Edit nama ikut mengganti USEDBY rincian ringkasan itu dalam SATU transaksi; rincian ringkasan lain tetap; gagal di
// tengah (ringkasan tak ada) = nol perubahan.
func TestEditGantiNamaRincian(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	if _, err := l.Simpan(ctx, penuh, "1000003", models.Isian{UsedBy: "Uji Risk Retro 2"}); err != nil {
		t.Fatal(err)
	}
	if r := g.Ringkasan["1000003"]; r.UsedBy != "Uji Risk Retro 2" || r.OperatorID != "UJI-ADMIN" || r.ModifiedDate != "20261008T030405.600 GMT" {
		t.Errorf("ringkasan %+v", r)
	}
	if g.Rincian["131720"].UsedBy != "Uji Risk Retro 2" || g.Rincian["131721"].UsedBy != "Uji Risk Retro 2" || g.Rincian["131722"].UsedBy != "UJI RISK B" {
		t.Errorf("nama rincian %+v", g.Rincian)
	}
	if _, err := l.Simpan(ctx, penuh, "1000003", models.Isian{UsedBy: "UJI RISK RETRO 2"}); err != nil {
		t.Errorf("nama sendiri beda huruf: %v", err)
	}
}

// View only (dan tanpa akun) ditolak di setiap tulis; nol tulisan.
func TestHakTulis(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	for _, a := range []services.Aktor{lihat, {Penuh: true}} {
		if _, err := l.Simpan(ctx, a, "", models.Isian{UsedBy: "UJI Y"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("simpan %+v: %v", a, err)
		}
		if _, err := l.Hapus(ctx, a, "1000003"); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("hapus %+v: %v", a, err)
		}
		if _, err := l.SimpanRincian(ctx, a, "1000003", "", models.IsianRincian{Contract: "1", Year: "9", Risk: "1"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("detail %+v: %v", a, err)
		}
		if _, err := l.Pratinjau(ctx, a, services.PermintaanUnggah{CSV: "x"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("pratinjau %+v: %v", a, err)
		}
		if _, err := l.SimpanUnggah(ctx, a, services.PermintaanUnggah{CSV: "x"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("simpan unggah %+v: %v", a, err)
		}
	}
	if len(g.Ringkasan) != 2 || len(g.Rincian) != 3 {
		t.Error("ditolak tetapi tertulis")
	}
}

// Delete (`DeleteSummaryDetail`): ringkasan BESERTA rinciannya, satu transaksi.
func TestHapus(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	r, err := l.Buka(ctx, "1000003")
	if err != nil || r.JumlahRincian == nil || *r.JumlahRincian != 2 || r.Diubah != "13-11-2019" {
		t.Fatalf("buka %+v %v", r, err)
	}
	h, err := l.Hapus(ctx, penuh, "1000003")
	if err != nil || h != (models.HasilHapus{ID: "1000003", RincianTerhapus: 2}) {
		t.Fatalf("hapus %+v %v", h, err)
	}
	if _, ada := g.Ringkasan["1000003"]; ada || len(g.Rincian) != 1 || g.Rincian["131722"].IDUsedBy != "1000004" {
		t.Errorf("sisa %v %v", g.Ringkasan, g.Rincian)
	}
	if _, err := l.Hapus(ctx, penuh, "1000003"); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("hapus lagi %v", err)
	}
}

// Grid ringkasan 50 per halaman (b10030), grid detail 200 per halaman (b10169/b10243), keduanya ID menaik.
func TestDaftarDanDetail(t *testing.T) {
	l := layanan(tiruan.Contoh())
	d, err := l.Daftar(ctx, models.Saringan{})
	if err != nil || d.Total != 2 || d.Ukuran != 50 || d.Daftar[0].ID != "1000003" {
		t.Errorf("bawaan ID menaik %+v %v", d, err)
	}
	if d, _ := l.Daftar(ctx, models.Saringan{UsedBy: "risk b"}); d.Total != 1 || d.Daftar[0].Diubah != "03-01-2024" {
		t.Errorf("saring nama %+v", d)
	}
	if d, _ := l.Daftar(ctx, models.Saringan{Urut: "ID", Turun: true}); d.Daftar[0].ID != "1000004" {
		t.Errorf("urut ID turun %+v", d)
	}
	k, err := l.DaftarRincian(ctx, "1000003", 0)
	if err != nil || k.Total != 2 || k.Ukuran != 200 || k.Daftar[0].ID != "131720" || k.Daftar[1].ID != "131721" {
		t.Errorf("detail ID menaik %+v %v", k, err)
	}
	if _, err := l.DaftarRincian(ctx, "999", 1); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("detail tak ada %v", err)
	}
}

// R/I RISK DETAIL: tambah (ID 131723 = '1' || LPAD(31723, 5); 131722 terpakai dilewati), USEDBY/IDUSEDBY dari ringkasan,
// YEAR / MONTH boleh kosong, kembar (CONTRACT, YEAR, MONTH) ditolak, ubah baris sendiri (EDIT), baris ringkasan lain
// 404; ringkasan ikut diperbarui; gagal = nol tulisan.
func TestSimpanRincian(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	k, err := l.SimpanRincian(ctx, penuh, "1000003", "", models.IsianRincian{Contract: "02", Year: "1", Risk: "7,50"})
	if err != nil {
		t.Fatal(err)
	}
	mau := models.Rincian{ID: "131723", IDUsedBy: "1000003", UsedBy: "UJI RISK RETRO", Contract: "2", Year: "1", Risk: "7.5"}
	if k != mau || g.Rincian["131723"] != mau {
		t.Errorf("tambah %+v / %+v", k, g.Rincian["131723"])
	}
	if r := g.Ringkasan["1000003"]; r.OperatorID != "UJI-ADMIN" || r.ModifiedDate != "20261008T030405.600 GMT" {
		t.Errorf("ringkasan sesudah tambah %+v", r)
	}
	if _, err := l.SimpanRincian(ctx, penuh, "1000003", "", models.IsianRincian{Contract: "01", Month: "012", Risk: "1"}); !errors.Is(err, services.ErrMasukanTidakSah) ||
		!strings.Contains(err.Error(), "CONTRACT 1, YEAR (empty), MONTH 12 already exists in R/I RISK NAME UJI RISK RETRO (ID 131721)") {
		t.Errorf("kembar %v", err)
	}
	// Kunci sama di ringkasan lain boleh (131722 = 2/1 milik 1000004).
	if _, err := l.SimpanRincian(ctx, penuh, "1000004", "", models.IsianRincian{Contract: "1", Year: "1", Risk: "1"}); err != nil {
		t.Errorf("kunci ringkasan lain %v", err)
	}
	if _, err := l.SimpanRincian(ctx, penuh, "1000003", "131721", models.IsianRincian{Contract: "1", Month: "12", Risk: "11"}); err != nil ||
		g.Rincian["131721"].Risk != "11" || g.Rincian["131721"].UsedBy != "UJI RISK RETRO" || g.Rincian["131721"].Month != "12" {
		t.Errorf("ubah sendiri %v %+v", err, g.Rincian["131721"])
	}
	if _, err := l.SimpanRincian(ctx, penuh, "1000003", "131721", models.IsianRincian{Contract: "1", Year: "1", Risk: "1"}); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("ubah jadi kembar %v", err)
	}
	if _, err := l.SimpanRincian(ctx, penuh, "1000003", "131722", models.IsianRincian{Contract: "8", Year: "8", Risk: "1"}); !errors.Is(err, services.ErrRincianTidakAda) {
		t.Errorf("baris ringkasan lain %v", err)
	}
	if _, err := l.SimpanRincian(ctx, penuh, "999", "", models.IsianRincian{Contract: "8", Year: "8", Risk: "1"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("ringkasan tak ada %v", err)
	}
	if _, err := l.SimpanRincian(ctx, penuh, "1000003", "", models.IsianRincian{}); !errors.Is(err, services.ErrMasukanTidakSah) ||
		!strings.Contains(err.Error(), "CONTRACT is required; RISK is required") {
		t.Errorf("wajib %v", err)
	}
}

const kepala = "USEDBY;CONTRACT;YEAR;MONTH;RISK\n"

// View Upload: tanpa menulis; nama dicocokkan tanpa beda huruf; kembar terhadap rincian lama ditolak per baris.
func TestPratinjau(t *testing.T) {
	g := tiruan.Contoh()
	teks := kepala +
		"uji risk retro;01;1;;9\n" + // 2 kembar 131720 (CONTRACT nol depan)
		"uji risk retro;3;1;;0,9\n" + // 3 sah, ringkasan 1000003
		"UJI RISK BARU;1;1;;1.5\n" // 4 sah, ringkasan baru
	h, err := layanan(g).Pratinjau(ctx, penuh, services.PermintaanUnggah{CSV: teks})
	if err != nil {
		t.Fatal(err)
	}
	if h.Sah || len(h.Galat) != 1 || h.Galat[0].Baris != 2 || !strings.Contains(h.Galat[0].Pesan, "already exists in R/I RISK NAME UJI RISK RETRO (ID 131720)") {
		t.Errorf("galat %+v", h.Galat)
	}
	if len(h.Baris) != 2 || h.Baris[0].UsedBy != "UJI RISK RETRO" || h.Baris[1].Risk != "1.5" {
		t.Errorf("baris %+v", h.Baris)
	}
	if len(h.Ringkasan) != 2 || h.Ringkasan[0] != (services.RingkasanUnggah{UsedBy: "UJI RISK RETRO", ID: "1000003", Jumlah: 1}) ||
		h.Ringkasan[1] != (services.RingkasanUnggah{UsedBy: "UJI RISK BARU", Baru: true, Jumlah: 1}) {
		t.Errorf("ringkasan %+v", h.Ringkasan)
	}
	if len(g.Ringkasan) != 2 || len(g.Rincian) != 3 {
		t.Error("pratinjau menulis")
	}
}

// Simpan Upload: satu galat = nol tertulis; sah = satu transaksi, ringkasan baru (site || LPAD 6) / lama, rincian
// '1' || LPAD 5.
func TestSimpanUnggah(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	_, err := l.SimpanUnggah(ctx, penuh, services.PermintaanUnggah{CSV: kepala + "UJI RISK RETRO;1;1;;1\nUJI BARU;1;x;;1\n"})
	if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(),
		"The upload was not saved. Row 2: CONTRACT 1, YEAR 1, MONTH (empty) already exists") || !strings.Contains(err.Error(), "Row 3: YEAR must be") {
		t.Errorf("tolak %v", err)
	}
	if len(g.Ringkasan) != 2 || len(g.Rincian) != 3 {
		t.Fatal("ditolak tetapi tertulis")
	}
	h, err := l.SimpanUnggah(ctx, penuh, services.PermintaanUnggah{CSV: kepala +
		"Uji Risk Retro;3;1;;\"0,25\"\nUJI BARU;1;1;;921,9\nUJI BARU;1;;24;580.894351210924\n"})
	if err != nil {
		t.Fatal(err)
	}
	if h.Disimpan != 3 || h.RingkasanBaru != 1 || len(h.Ringkasan) != 2 || h.Ringkasan[1].ID != "1000005" {
		t.Errorf("hasil %+v", h)
	}
	if r := g.Ringkasan["1000005"]; r.UsedBy != "UJI BARU" || r.OperatorID != "UJI-ADMIN" {
		t.Errorf("ringkasan baru %+v", r)
	}
	if r := g.Ringkasan["1000003"]; r.UsedBy != "UJI RISK RETRO" || r.ModifiedDate != "20261008T030405.600 GMT" {
		t.Errorf("ringkasan lama diperbarui %+v", r)
	}
	mau := map[string]models.Rincian{
		"131723": {ID: "131723", IDUsedBy: "1000003", UsedBy: "UJI RISK RETRO", Contract: "3", Year: "1", Risk: "0.25"},
		"131724": {ID: "131724", IDUsedBy: "1000005", UsedBy: "UJI BARU", Contract: "1", Year: "1", Risk: "921.9"},
		"131725": {ID: "131725", IDUsedBy: "1000005", UsedBy: "UJI BARU", Contract: "1", Month: "24", Risk: "580.894351210924"},
	}
	for id, k := range mau {
		if g.Rincian[id] != k {
			t.Errorf("rincian %s %+v, mau %+v", id, g.Rincian[id], k)
		}
	}
}

func TestUnggahNamaGandaDanKalimat(t *testing.T) {
	g := tiruan.Contoh()
	g.Ringkasan["1000009"] = models.Ringkasan{ID: "1000009", UsedBy: "uji risk b "}
	h, err := layanan(g).Pratinjau(ctx, penuh, services.PermintaanUnggah{CSV: kepala + "UJI RISK B;1;1;;1\n"})
	if err != nil || h.Sah || len(h.Galat) != 1 || !strings.Contains(h.Galat[0].Pesan, "matches more than one R/I risk summary (IDs 1000004, 1000009)") {
		t.Errorf("%+v %v", h, err)
	}
	var gb []models.GalatBaris
	for i := 0; i < services.MaksKalimatGalat+3; i++ {
		gb = append(gb, models.GalatBaris{Baris: i + 2, Pesan: "X"})
	}
	if s := services.KalimatGalat(gb); !strings.HasSuffix(s, "Row 51: X. And 3 more rows.") {
		t.Errorf("%s", s[len(s)-60:])
	}
	if services.FormatWaktuPega(jam) != "20261008T030405.600 GMT" || services.TampilTanggal("20191113T025753.044 GMT") != "13-11-2019" {
		t.Error("format waktu Pega")
	}
}
