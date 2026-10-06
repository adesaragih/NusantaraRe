package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/riratelife/backend/models"
	"nusantarare/modul/riratelife/backend/services"
	"nusantarare/modul/riratelife/backend/tiruan"
)

var ctx = context.Background()

var (
	penuh = services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
	lihat = services.Aktor{AkunID: "UJI-LIHAT"}
	jam   = time.Date(2026, 10, 5, 3, 4, 5, 600e6, time.UTC)
)

func layanan(g *tiruan.Gudang) *services.Layanan {
	return services.BaruLayanan(g, tiruan.Transaksi{G: g}.Jalankan).DenganJam(func() time.Time { return jam })
}

// Add: nama dipangkas; ID dari sequence, nomor terpakai (105) dilewati; OPERATORID akun; MODIFIEDDATE format Pega.
func TestAdd(t *testing.T) {
	g := tiruan.Contoh()
	r, err := layanan(g).Simpan(ctx, penuh, "", models.Isian{UsedBy: "  UJI RATE BARU "})
	if err != nil {
		t.Fatal(err)
	}
	mau := models.Ringkasan{ID: "106", UsedBy: "UJI RATE BARU", OperatorID: "UJI-ADMIN", ModifiedDate: "20261005T030405.600 GMT"}
	if g.Ringkasan["106"] != mau {
		t.Errorf("tersimpan %+v", g.Ringkasan["106"])
	}
	if r.ID != "106" || r.Diubah != "05-10-2026" {
		t.Errorf("jawab %+v", r)
	}
}

// Nama wajib, tidak kembar tanpa beda huruf (Add maupun Edit ke nama baris lain); View only ditolak.
func TestValidasiNama(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	for isi, mau := range map[string]string{
		"  ":                                    "R/I RATE NAME is required",
		" uji rate a ":                          "R/I RATE NAME uji rate a is already used by ID 101",
		strings.Repeat("A", models.BatasNama+1): "longer than 200",
	} {
		if _, err := l.Simpan(ctx, penuh, "", models.Isian{UsedBy: isi}); !errors.Is(err, services.ErrMasukanTidakSah) ||
			!strings.Contains(err.Error(), mau) {
			t.Errorf("%q: %v", isi, err)
		}
	}
	if _, err := l.Simpan(ctx, penuh, "102", models.Isian{UsedBy: "Uji Rate A"}); err == nil {
		t.Error("edit ke nama baris lain harus ditolak")
	}
	if _, err := l.Simpan(ctx, penuh, "999", models.Isian{UsedBy: "UJI X"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("edit ID tak ada: %v", err)
	}
	for _, a := range []services.Aktor{lihat, {Penuh: true}} {
		if _, err := l.Simpan(ctx, a, "", models.Isian{UsedBy: "UJI Y"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("%+v: %v", a, err)
		}
		if _, err := l.Hapus(ctx, a, "101"); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("hapus %+v: %v", a, err)
		}
		if _, err := l.Pratinjau(ctx, a, services.PermintaanUnggah{CSV: "x"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("pratinjau %+v: %v", a, err)
		}
		if _, err := l.SimpanUnggah(ctx, a, services.PermintaanUnggah{CSV: "x"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("simpan unggah %+v: %v", a, err)
		}
	}
	if len(g.Ringkasan) != 3 {
		t.Error("ditolak tetapi tertulis")
	}
}

// Edit: nama sendiri boleh (beda huruf); salinan nama di baris rate ikut berganti; baris rate ringkasan lain tetap.
func TestEdit(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	if _, err := l.Simpan(ctx, penuh, "101", models.Isian{UsedBy: "Uji Rate A2"}); err != nil {
		t.Fatal(err)
	}
	if r := g.Ringkasan["101"]; r.UsedBy != "Uji Rate A2" || r.OperatorID != "UJI-ADMIN" || r.ModifiedDate != "20261005T030405.600 GMT" {
		t.Errorf("ringkasan %+v", r)
	}
	if g.Rate["9000"].UsedBy != "Uji Rate A2" || g.Rate["9100"].UsedBy != "UJI RATE B" {
		t.Errorf("nama rate %+v %+v", g.Rate["9000"], g.Rate["9100"])
	}
	if _, err := l.Simpan(ctx, penuh, "101", models.Isian{UsedBy: "UJI RATE A2"}); err != nil {
		t.Errorf("nama sendiri beda huruf: %v", err)
	}
}

// Delete: ringkasan BESERTA rate ber-IDUSEDBY sama; rate ringkasan lain tetap.
func TestHapus(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	r, err := l.Buka(ctx, "101")
	if err != nil || r.JumlahRate == nil || *r.JumlahRate != 3 || r.Diubah != "02-01-2024" {
		t.Fatalf("buka %+v %v", r, err)
	}
	h, err := l.Hapus(ctx, penuh, "101")
	if err != nil || h != (models.HasilHapus{ID: "101", RateTerhapus: 3}) {
		t.Fatalf("hapus %+v %v", h, err)
	}
	if _, ada := g.Ringkasan["101"]; ada || len(g.Rate) != 1 || g.Rate["9100"].IDUsedBy != "102" {
		t.Errorf("sisa %v %v", g.Ringkasan, g.Rate)
	}
	if _, err := l.Hapus(ctx, penuh, "101"); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("hapus lagi %v", err)
	}
}

func TestDaftarDanDetail(t *testing.T) {
	l := layanan(tiruan.Contoh())
	d, err := l.Daftar(ctx, models.Saringan{})
	if err != nil || d.Total != 3 || d.Halaman != 1 || d.Ukuran != 50 || d.Daftar[0].ID != "105" {
		t.Errorf("bawaan ID menurun %+v %v", d, err)
	}
	if d, _ := l.Daftar(ctx, models.Saringan{UsedBy: "rate b"}); d.Total != 1 || d.Daftar[0].ID != "102" || d.Daftar[0].Diubah != "03-01-2024" {
		t.Errorf("saring nama %+v", d)
	}
	if d, _ := l.Daftar(ctx, models.Saringan{Urut: "USEDBY"}); d.Daftar[0].UsedBy != "UJI RATE A" {
		t.Errorf("urut nama %+v", d)
	}
	if d, _ := l.Daftar(ctx, models.Saringan{Halaman: 2}); len(d.Daftar) != 0 || d.Total != 3 {
		t.Errorf("halaman 2 %+v", d)
	}
	r, err := l.DaftarRate(ctx, "101", 0)
	if err != nil || r.Total != 3 || r.Daftar[0].ID != "9003" || r.Daftar[1].ID != "9002" || r.Daftar[2].ID != "9000" {
		t.Errorf("detail urut GENDER, CONTRACT, AGE %+v %v", r, err)
	}
	if _, err := l.DaftarRate(ctx, "999", 1); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("detail tak ada %v", err)
	}
}

const kepala = "USEDBY;CONTRACT;GENDER;AGE;RATE\n"

// View Upload: tanpa menulis; nama dicocokkan tanpa beda huruf ke ringkasan; kembar terhadap baris lama (AGE `05`
// warisan setara `5`, CONTRACT kosong) ditolak per baris.
func TestPratinjau(t *testing.T) {
	g := tiruan.Contoh()
	teks := kepala +
		"uji rate a;;U;5;9\n" + // 2 kembar rate 9002
		"uji rate a;2;U;30;0,9\n" + // 3 sah, ringkasan 101
		"UJI RATE BARU;1;F;1;1.5\n" // 4 sah, ringkasan baru
	h, err := layanan(g).Pratinjau(ctx, penuh, services.PermintaanUnggah{CSV: teks})
	if err != nil {
		t.Fatal(err)
	}
	if h.Sah || len(h.Galat) != 1 || h.Galat[0].Baris != 2 || !strings.Contains(h.Galat[0].Pesan, "already exists in R/I RATE NAME UJI RATE A (rate ID 9002)") {
		t.Errorf("galat %+v", h.Galat)
	}
	if len(h.Baris) != 2 || h.Baris[0].UsedBy != "UJI RATE A" || h.Baris[1].Rate != "1,5" {
		t.Errorf("baris %+v", h.Baris)
	}
	if len(h.Ringkasan) != 2 || h.Ringkasan[0] != (services.RingkasanUnggah{UsedBy: "UJI RATE A", ID: "101", Jumlah: 1}) ||
		h.Ringkasan[1] != (services.RingkasanUnggah{UsedBy: "UJI RATE BARU", Baru: true, Jumlah: 1}) {
		t.Errorf("ringkasan %+v", h.Ringkasan)
	}
	if len(g.Ringkasan) != 3 || len(g.Rate) != 4 {
		t.Error("pratinjau menulis")
	}
	if _, err := layanan(g).Pratinjau(ctx, penuh, services.PermintaanUnggah{CSV: "A;B\n"}); !errors.Is(err, services.ErrMasukanTidakSah) ||
		!strings.Contains(err.Error(), "header must be") {
		t.Errorf("kepala salah %v", err)
	}
}

// Simpan Upload: satu galat = nol tertulis (422 berkalimat per baris); sah = satu transaksi, ringkasan baru / lama,
// ID rate dari sequence melewati nomor terpakai.
func TestSimpanUnggah(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	_, err := l.SimpanUnggah(ctx, penuh, services.PermintaanUnggah{CSV: kepala + "UJI RATE A;1;U;30;1\nUJI BARU;1;X;1;1\n"})
	if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(),
		"The upload was not saved. Row 2: GENDER U, AGE 30, CONTRACT 1 already exists") || !strings.Contains(err.Error(), "Row 3: GENDER must be U, M, or F") {
		t.Errorf("tolak %v", err)
	}
	if len(g.Ringkasan) != 3 || len(g.Rate) != 4 {
		t.Fatal("ditolak tetapi tertulis")
	}
	h, err := l.SimpanUnggah(ctx, penuh, services.PermintaanUnggah{CSV: kepala +
		"Uji Rate A;3;U;30;\"0,25\"\nUJI BARU;;F;1;2.5\nUJI BARU;;M;1;3\n"})
	if err != nil {
		t.Fatal(err)
	}
	if h.Disimpan != 3 || h.RingkasanBaru != 1 || len(h.Ringkasan) != 2 || h.Ringkasan[1].ID != "106" {
		t.Errorf("hasil %+v", h)
	}
	if r := g.Ringkasan["106"]; r.UsedBy != "UJI BARU" || r.OperatorID != "UJI-ADMIN" {
		t.Errorf("ringkasan baru %+v", r)
	}
	if r := g.Ringkasan["101"]; r.UsedBy != "UJI RATE A" || r.ModifiedDate != "20261005T030405.600 GMT" {
		t.Errorf("ringkasan lama diperbarui %+v", r)
	}
	mau := map[string]models.Rate{
		"9001": {ID: "9001", IDUsedBy: "101", UsedBy: "UJI RATE A", Gender: "U", Contract: "3", Age: "30", Rate: "0,25"},
		"9004": {ID: "9004", IDUsedBy: "106", UsedBy: "UJI BARU", Gender: "F", Contract: "", Age: "1", Rate: "2,5"},
		"9005": {ID: "9005", IDUsedBy: "106", UsedBy: "UJI BARU", Gender: "M", Contract: "", Age: "1", Rate: "3"},
	}
	for id, r := range mau {
		if g.Rate[id] != r {
			t.Errorf("rate %s %+v, mau %+v", id, g.Rate[id], r)
		}
	}
	if len(g.Rate) != 7 {
		t.Errorf("jumlah rate %d", len(g.Rate))
	}
}

// Nama yang cocok ke lebih dari satu ringkasan warisan ditolak (tidak ditebak).
func TestUnggahNamaGanda(t *testing.T) {
	g := tiruan.Contoh()
	g.Ringkasan["107"] = models.Ringkasan{ID: "107", UsedBy: "uji rate b "}
	h, err := layanan(g).Pratinjau(ctx, penuh, services.PermintaanUnggah{CSV: kepala + "UJI RATE B;1;U;1;1\n"})
	if err != nil || h.Sah || len(h.Galat) != 1 || !strings.Contains(h.Galat[0].Pesan, "matches more than one R/I rate summary (IDs 102, 107)") {
		t.Errorf("%+v %v", h, err)
	}
}

func TestKalimatGalat(t *testing.T) {
	var g []models.GalatBaris
	for i := 0; i < services.MaksKalimatGalat+3; i++ {
		g = append(g, models.GalatBaris{Baris: i + 2, Pesan: "X"})
	}
	if s := services.KalimatGalat(g); !strings.HasSuffix(s, "Row 51: X. And 3 more rows.") {
		t.Errorf("%s", s[len(s)-60:])
	}
}
