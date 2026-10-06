package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/ricommlife/backend/models"
	"nusantarare/modul/ricommlife/backend/repository"
	"nusantarare/modul/ricommlife/backend/services"
	"nusantarare/modul/ricommlife/backend/tiruan"
)

var ctx = context.Background()

var (
	penuh = services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
	lihat = services.Aktor{AkunID: "UJI-LIHAT"}
	jam   = time.Date(2026, 10, 6, 3, 4, 5, 600e6, time.UTC)
)

func layanan(g *tiruan.Gudang) *services.Layanan {
	return services.BaruLayanan(g, tiruan.Transaksi{G: g}.Jalankan).DenganJam(func() time.Time { return jam })
}

// Add: nama dipangkas; ID = site || LPAD(seq, 6) - 1000004 terpakai dilewati; OPERATORID akun; MODIFIEDDATE Pega.
func TestAdd(t *testing.T) {
	g := tiruan.Contoh()
	r, err := layanan(g).Simpan(ctx, penuh, "", models.Isian{UsedBy: "  UJI COMM BARU "})
	if err != nil {
		t.Fatal(err)
	}
	mau := models.Ringkasan{ID: "1000005", UsedBy: "UJI COMM BARU", OperatorID: "UJI-ADMIN", ModifiedDate: "20261006T030405.600 GMT"}
	if g.Ringkasan["1000005"] != mau {
		t.Errorf("tersimpan %+v", g.Ringkasan["1000005"])
	}
	if r.ID != "1000005" || r.Diubah != "06-10-2026" {
		t.Errorf("jawab %+v", r)
	}
}

// Situs tidak tepat satu baris = galat, nol tulisan.
func TestSitusTidakTepatSatu(t *testing.T) {
	for _, situs := range [][]string{nil, {"1", "2"}} {
		g := tiruan.Contoh()
		g.SitusAktif = situs
		if _, err := layanan(g).Simpan(ctx, penuh, "", models.Isian{UsedBy: "UJI X"}); !errors.Is(err, repository.ErrSitus) {
			t.Errorf("%v: %v", situs, err)
		}
		if len(g.Ringkasan) != 2 {
			t.Error("tertulis tanpa situs")
		}
	}
	g := tiruan.Contoh()
	g.SeqKomisi = 1234567
	if _, err := layanan(g).SimpanKomisi(ctx, penuh, "1000003", "", models.IsianKomisi{Contract: "9", Year: "9", Comm: "1"}); !errors.Is(err, models.ErrIDTidakSah) {
		t.Errorf("nomor 7 angka: %v", err)
	}
}

// Nama wajib, tidak kembar tanpa beda huruf; View only ditolak di setiap tulis.
func TestValidasiNamaDanHak(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	for isi, mau := range map[string]string{
		"  ":                                    "R/I COMM NAME is required",
		" uji comm retro ":                      "R/I COMM NAME uji comm retro is already used by ID 1000003",
		strings.Repeat("A", models.BatasNama+1): "longer than 200",
	} {
		if _, err := l.Simpan(ctx, penuh, "", models.Isian{UsedBy: isi}); !errors.Is(err, services.ErrMasukanTidakSah) ||
			!strings.Contains(err.Error(), mau) {
			t.Errorf("%q: %v", isi, err)
		}
	}
	if _, err := l.Simpan(ctx, penuh, "1000004", models.Isian{UsedBy: "Uji Comm Retro"}); err == nil {
		t.Error("edit ke nama baris lain harus ditolak")
	}
	if _, err := l.Simpan(ctx, penuh, "999", models.Isian{UsedBy: "UJI X"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("edit ID tak ada: %v", err)
	}
	for _, a := range []services.Aktor{lihat, {Penuh: true}} {
		if _, err := l.Simpan(ctx, a, "", models.Isian{UsedBy: "UJI Y"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("%+v: %v", a, err)
		}
		if _, err := l.Hapus(ctx, a, "1000003"); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("hapus %+v: %v", a, err)
		}
		if _, err := l.SimpanKomisi(ctx, a, "1000003", "", models.IsianKomisi{Contract: "1", Year: "9", Comm: "1"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("detail %+v: %v", a, err)
		}
		if _, err := l.Pratinjau(ctx, a, services.PermintaanUnggah{CSV: "x"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("pratinjau %+v: %v", a, err)
		}
		if _, err := l.SimpanUnggah(ctx, a, services.PermintaanUnggah{CSV: "x"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("simpan unggah %+v: %v", a, err)
		}
	}
	if len(g.Ringkasan) != 2 || len(g.Komisi) != 3 {
		t.Error("ditolak tetapi tertulis")
	}
}

// Edit nama ikut mengganti USEDBY rincian ringkasan itu; rincian ringkasan lain tetap.
func TestEdit(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	if _, err := l.Simpan(ctx, penuh, "1000003", models.Isian{UsedBy: "Uji Comm Retro 2"}); err != nil {
		t.Fatal(err)
	}
	if r := g.Ringkasan["1000003"]; r.UsedBy != "Uji Comm Retro 2" || r.OperatorID != "UJI-ADMIN" || r.ModifiedDate != "20261006T030405.600 GMT" {
		t.Errorf("ringkasan %+v", r)
	}
	if g.Komisi["1000040"].UsedBy != "Uji Comm Retro 2" || g.Komisi["1000041"].UsedBy != "Uji Comm Retro 2" || g.Komisi["1000042"].UsedBy != "UJI COMM B" {
		t.Errorf("nama rincian %+v", g.Komisi)
	}
	if _, err := l.Simpan(ctx, penuh, "1000003", models.Isian{UsedBy: "UJI COMM RETRO 2"}); err != nil {
		t.Errorf("nama sendiri beda huruf: %v", err)
	}
}

// Delete: ringkasan BESERTA rinciannya, satu transaksi.
func TestHapus(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	r, err := l.Buka(ctx, "1000003")
	if err != nil || r.JumlahKomisi == nil || *r.JumlahKomisi != 2 || r.Diubah != "05-12-2018" {
		t.Fatalf("buka %+v %v", r, err)
	}
	h, err := l.Hapus(ctx, penuh, "1000003")
	if err != nil || h != (models.HasilHapus{ID: "1000003", KomisiTerhapus: 2}) {
		t.Fatalf("hapus %+v %v", h, err)
	}
	if _, ada := g.Ringkasan["1000003"]; ada || len(g.Komisi) != 1 || g.Komisi["1000042"].IDUsedBy != "1000004" {
		t.Errorf("sisa %v %v", g.Ringkasan, g.Komisi)
	}
	if _, err := l.Hapus(ctx, penuh, "1000003"); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("hapus lagi %v", err)
	}
}

func TestDaftarDanDetail(t *testing.T) {
	l := layanan(tiruan.Contoh())
	d, err := l.Daftar(ctx, models.Saringan{})
	if err != nil || d.Total != 2 || d.Ukuran != 50 || d.Daftar[0].ID != "1000003" {
		t.Errorf("bawaan ID menaik %+v %v", d, err)
	}
	if d, _ := l.Daftar(ctx, models.Saringan{UsedBy: "comm b"}); d.Total != 1 || d.Daftar[0].Diubah != "03-01-2024" {
		t.Errorf("saring nama %+v", d)
	}
	if d, _ := l.Daftar(ctx, models.Saringan{Urut: "ID", Turun: true}); d.Daftar[0].ID != "1000004" {
		t.Errorf("urut ID turun %+v", d)
	}
	k, err := l.DaftarKomisi(ctx, "1000003", 0)
	if err != nil || k.Total != 2 || k.Ukuran != 50 || k.Daftar[0].ID != "1000040" || k.Daftar[1].ID != "1000041" {
		t.Errorf("detail ID menaik %+v %v", k, err)
	}
	if _, err := l.DaftarKomisi(ctx, "999", 1); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("detail tak ada %v", err)
	}
}

// R/I COMM DETAIL: tambah (ID 1000044 = site 1 + seq 44), USEDBY/IDUSEDBY dari ringkasan, kembar (CONTRACT, YEAR)
// ditolak, ubah baris sendiri, baris ringkasan lain 404; ringkasan ikut diperbarui; gagal = nol tulisan.
func TestSimpanKomisi(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	k, err := l.SimpanKomisi(ctx, penuh, "1000003", "", models.IsianKomisi{Contract: "02", Year: "1", Comm: "7,50"})
	if err != nil {
		t.Fatal(err)
	}
	mau := models.Komisi{ID: "1000044", IDUsedBy: "1000003", UsedBy: "UJI COMM RETRO", Contract: "2", Year: "1", Comm: "7.5"}
	if k != mau || g.Komisi["1000044"] != mau {
		t.Errorf("tambah %+v / %+v", k, g.Komisi["1000044"])
	}
	if r := g.Ringkasan["1000003"]; r.OperatorID != "UJI-ADMIN" || r.ModifiedDate != "20261006T030405.600 GMT" {
		t.Errorf("ringkasan sesudah tambah %+v", r)
	}
	if _, err := l.SimpanKomisi(ctx, penuh, "1000003", "", models.IsianKomisi{Contract: "1", Year: "02", Comm: "1"}); !errors.Is(err, services.ErrMasukanTidakSah) ||
		!strings.Contains(err.Error(), "CONTRACT 1, YEAR 2 already exists in R/I COMM NAME UJI COMM RETRO (ID 1000041)") {
		t.Errorf("kembar %v", err)
	}
	// Kunci sama di ringkasan lain boleh (1000042 = 2/1 milik 1000004).
	if _, err := l.SimpanKomisi(ctx, penuh, "1000004", "", models.IsianKomisi{Contract: "1", Year: "1", Comm: "1"}); err != nil {
		t.Errorf("kunci ringkasan lain %v", err)
	}
	if _, err := l.SimpanKomisi(ctx, penuh, "1000003", "1000041", models.IsianKomisi{Contract: "1", Year: "2", Comm: "11"}); err != nil ||
		g.Komisi["1000041"].Comm != "11" || g.Komisi["1000041"].UsedBy != "UJI COMM RETRO" {
		t.Errorf("ubah sendiri %v %+v", err, g.Komisi["1000041"])
	}
	if _, err := l.SimpanKomisi(ctx, penuh, "1000003", "1000041", models.IsianKomisi{Contract: "1", Year: "1", Comm: "1"}); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("ubah jadi kembar %v", err)
	}
	if _, err := l.SimpanKomisi(ctx, penuh, "1000003", "1000042", models.IsianKomisi{Contract: "8", Year: "8", Comm: "1"}); !errors.Is(err, services.ErrKomisiTidakAda) {
		t.Errorf("baris ringkasan lain %v", err)
	}
	if _, err := l.SimpanKomisi(ctx, penuh, "999", "", models.IsianKomisi{Contract: "8", Year: "8", Comm: "1"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("ringkasan tak ada %v", err)
	}
	if _, err := l.SimpanKomisi(ctx, penuh, "1000003", "", models.IsianKomisi{}); !errors.Is(err, services.ErrMasukanTidakSah) ||
		!strings.Contains(err.Error(), "CONTRACT is required; YEAR is required; COMM is required") {
		t.Errorf("wajib %v", err)
	}
}

const kepala = "USEDBY;CONTRACT;YEAR;COMM\n"

// View Upload: tanpa menulis; nama dicocokkan tanpa beda huruf; kembar terhadap rincian lama ditolak per baris.
func TestPratinjau(t *testing.T) {
	g := tiruan.Contoh()
	teks := kepala +
		"uji comm retro;1;01;9\n" + // 2 kembar 1000040
		"uji comm retro;3;1;0,9\n" + // 3 sah, ringkasan 1000003
		"UJI COMM BARU;1;1;1.5\n" // 4 sah, ringkasan baru
	h, err := layanan(g).Pratinjau(ctx, penuh, services.PermintaanUnggah{CSV: teks})
	if err != nil {
		t.Fatal(err)
	}
	if h.Sah || len(h.Galat) != 1 || h.Galat[0].Baris != 2 || !strings.Contains(h.Galat[0].Pesan, "already exists in R/I COMM NAME UJI COMM RETRO (ID 1000040)") {
		t.Errorf("galat %+v", h.Galat)
	}
	if len(h.Baris) != 2 || h.Baris[0].UsedBy != "UJI COMM RETRO" || h.Baris[1].Comm != "1.5" {
		t.Errorf("baris %+v", h.Baris)
	}
	if len(h.Ringkasan) != 2 || h.Ringkasan[0] != (services.RingkasanUnggah{UsedBy: "UJI COMM RETRO", ID: "1000003", Jumlah: 1}) ||
		h.Ringkasan[1] != (services.RingkasanUnggah{UsedBy: "UJI COMM BARU", Baru: true, Jumlah: 1}) {
		t.Errorf("ringkasan %+v", h.Ringkasan)
	}
	if len(g.Ringkasan) != 2 || len(g.Komisi) != 3 {
		t.Error("pratinjau menulis")
	}
}

// Simpan Upload: satu galat = nol tertulis; sah = satu transaksi, ringkasan baru / lama, ID site + sequence.
func TestSimpanUnggah(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	_, err := l.SimpanUnggah(ctx, penuh, services.PermintaanUnggah{CSV: kepala + "UJI COMM RETRO;1;1;1\nUJI BARU;1;x;1\n"})
	if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(),
		"The upload was not saved. Row 2: CONTRACT 1, YEAR 1 already exists") || !strings.Contains(err.Error(), "Row 3: YEAR must be") {
		t.Errorf("tolak %v", err)
	}
	if len(g.Ringkasan) != 2 || len(g.Komisi) != 3 {
		t.Fatal("ditolak tetapi tertulis")
	}
	h, err := l.SimpanUnggah(ctx, penuh, services.PermintaanUnggah{CSV: kepala +
		"Uji Comm Retro;3;1;\"0,25\"\nUJI BARU;1;1;2.5\nUJI BARU;1;2;3\n"})
	if err != nil {
		t.Fatal(err)
	}
	if h.Disimpan != 3 || h.RingkasanBaru != 1 || len(h.Ringkasan) != 2 || h.Ringkasan[1].ID != "1000005" {
		t.Errorf("hasil %+v", h)
	}
	if r := g.Ringkasan["1000005"]; r.UsedBy != "UJI BARU" || r.OperatorID != "UJI-ADMIN" {
		t.Errorf("ringkasan baru %+v", r)
	}
	if r := g.Ringkasan["1000003"]; r.UsedBy != "UJI COMM RETRO" || r.ModifiedDate != "20261006T030405.600 GMT" {
		t.Errorf("ringkasan lama diperbarui %+v", r)
	}
	mau := map[string]models.Komisi{
		"1000044": {ID: "1000044", IDUsedBy: "1000003", UsedBy: "UJI COMM RETRO", Contract: "3", Year: "1", Comm: "0.25"},
		"1000045": {ID: "1000045", IDUsedBy: "1000005", UsedBy: "UJI BARU", Contract: "1", Year: "1", Comm: "2.5"},
		"1000046": {ID: "1000046", IDUsedBy: "1000005", UsedBy: "UJI BARU", Contract: "1", Year: "2", Comm: "3"},
	}
	for id, k := range mau {
		if g.Komisi[id] != k {
			t.Errorf("rincian %s %+v, mau %+v", id, g.Komisi[id], k)
		}
	}
}

func TestUnggahNamaGandaDanKalimat(t *testing.T) {
	g := tiruan.Contoh()
	g.Ringkasan["1000009"] = models.Ringkasan{ID: "1000009", UsedBy: "uji comm b "}
	h, err := layanan(g).Pratinjau(ctx, penuh, services.PermintaanUnggah{CSV: kepala + "UJI COMM B;1;1;1\n"})
	if err != nil || h.Sah || len(h.Galat) != 1 || !strings.Contains(h.Galat[0].Pesan, "matches more than one R/I comm summary (IDs 1000004, 1000009)") {
		t.Errorf("%+v %v", h, err)
	}
	var gb []models.GalatBaris
	for i := 0; i < services.MaksKalimatGalat+3; i++ {
		gb = append(gb, models.GalatBaris{Baris: i + 2, Pesan: "X"})
	}
	if s := services.KalimatGalat(gb); !strings.HasSuffix(s, "Row 51: X. And 3 more rows.") {
		t.Errorf("%s", s[len(s)-60:])
	}
	if services.FormatWaktuPega(jam) != "20261006T030405.600 GMT" || services.TampilTanggal("20181205T073755.559 GMT") != "05-12-2018" {
		t.Error("format waktu Pega")
	}
}
