//go:build db

// Seam `repository` terhadap skema uji Oracle NYATA (brief bab 5).
// Tabel dibuat oleh skema uji, diisi fixture, lalu dibaca kembali.
//
// Jalankan: make test-db   (perlu ORACLE_DSN dan ORACLE_SCHEMA)
package repository_test

import (
	"context"
	"testing"
	"time"

	"nusantarare/internal/config"
	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/repository/skemauji"
	"nusantarare/pkg/utils"
)

func siapkan(t *testing.T) (*repository.KlaimLife, func()) {
	t.Helper()
	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		t.Skipf("lewati: %v", err)
	}
	ctx := context.Background()
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	if err := skemauji.Pasang(ctx, sqlDB, skema); err != nil {
		t.Fatalf("memasang skema uji: %v", err)
	}
	if err := skemauji.IsiContoh(ctx, sqlDB, skema); err != nil {
		t.Fatalf("mengisi fixture: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	db, err := repository.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return repository.NewKlaimLife(db), func() {
		_ = skemauji.Bongkar(ctx, sqlDB, skema)
		_ = db.Close()
		_ = sqlDB.Close()
	}
}

func TestAmbilHeader(t *testing.T) {
	repo, bersihkan := siapkan(t)
	defer bersihkan()
	ctx := context.Background()

	k, err := repo.AmbilHeader(ctx, "UJI-KLAIM-1")
	if err != nil {
		t.Fatal(err)
	}
	if k == nil {
		t.Fatal("klaim tidak ditemukan")
	}
	if k.NomorKlaim != "UJI-CLM-0001" || k.NomorPolis != "UJI-POL-0001" {
		t.Errorf("header salah: %+v", k)
	}

	hilang, err := repo.AmbilHeader(ctx, "TIDAK-ADA")
	if err != nil {
		t.Fatal(err)
	}
	if hilang != nil {
		t.Errorf("klaim yang tidak ada mengembalikan %+v", hilang)
	}
}

// ADR-U-0022: kode dan penanda tetap teks. Nilai kolom diperiksa LANGSUNG lewat
// pembacaan publik seam repository - "006" yang kembali sebagai "6" memecahkan
// penggolong, dan pulang-pergi saja tidak akan menangkapnya.
func TestNomorSertifikatBerawalanNolTetapTeks(t *testing.T) {
	repo, bersihkan := siapkan(t)
	defer bersihkan()

	peserta, err := repo.AmbilPeserta(context.Background(), "UJI-KLAIM-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(peserta) != 2 {
		t.Fatalf("cacah peserta = %d, mau 2", len(peserta))
	}
	if peserta[0].NomorSertifikat != "006" {
		t.Errorf("nomor sertifikat = %q, mau %q - awalan nol hilang",
			peserta[0].NomorSertifikat, "006")
	}
	if peserta[1].NomorSertifikat != "010" {
		t.Errorf("nomor sertifikat = %q, mau %q", peserta[1].NomorSertifikat, "010")
	}
}

// Tiket 01 AC-1: SELURUH baris AdjustmentList, bukan hanya yang terakhir.
func TestAmbilSeluruhBarisBukanHanyaTerakhir(t *testing.T) {
	repo, bersihkan := siapkan(t)
	defer bersihkan()

	baris, err := repo.AmbilBaris(context.Background(), "UJI-KLAIM-1")
	if err != nil {
		t.Fatal(err)
	}
	if n := len(baris["UJI-PESERTA-1"]); n != 3 {
		t.Errorf("peserta 1 punya %d baris, mau 3", n)
	}
	if n := len(baris["UJI-PESERTA-2"]); n != 2 {
		t.Errorf("peserta 2 punya %d baris, mau 2", n)
	}
	total := 0
	for _, b := range baris {
		total += len(b)
	}
	if total != 5 {
		t.Errorf("total baris = %d, mau 5", total)
	}
}

// Tiket 01 AC-4 / ADR-U-0003 - ADR-U-0016: uang tidak pernah lewat float.
// Nilai dibandingkan terhadap angka yang ditulis fixture, digit demi digit.
func TestJumlahKlaimUtuhTanpaFloat(t *testing.T) {
	repo, bersihkan := siapkan(t)
	defer bersihkan()

	baris, err := repo.AmbilBaris(context.Background(), "UJI-KLAIM-1")
	if err != nil {
		t.Fatal(err)
	}
	p1 := baris["UJI-PESERTA-1"]

	mau := map[string]string{
		"UJI-ADJ-1": "1234567890.12345678",
		"UJI-ADJ-2": "0.00000001",
		"UJI-ADJ-3": "250000",
	}
	for _, b := range p1 {
		harap, ada := mau[b.ID]
		if !ada {
			t.Fatalf("baris tak terduga: %s", b.ID)
		}
		// utils.FormatDecimal adalah satu-satunya jalan keluar desimal ke teks
		// (ADR-U-0034), dan ia aman terhadap nilai kosong.
		if got := utils.FormatDecimal(b.JumlahKlaim.Amount); got != harap {
			t.Errorf("%s: jumlah = %q, mau %q", b.ID, got, harap)
		}
		if b.JumlahKlaim.Currency != "IDR" {
			t.Errorf("%s: mata uang = %q", b.ID, b.JumlahKlaim.Currency)
		}
	}
}

// ADR-U-0027: kolom kosong tetap kosong, tidak menjadi nol.
func TestKolomKosongTidakMenjadiNol(t *testing.T) {
	repo, bersihkan := siapkan(t)
	defer bersihkan()

	baris, err := repo.AmbilBaris(context.Background(), "UJI-KLAIM-1")
	if err != nil {
		t.Fatal(err)
	}
	var kosong *models.BarisAdjustment
	for i, b := range baris["UJI-PESERTA-2"] {
		if b.ID == "UJI-ADJ-5" {
			kosong = &baris["UJI-PESERTA-2"][i]
		}
	}
	if kosong == nil {
		t.Fatal("baris UJI-ADJ-5 tidak terbaca")
	}
	if !kosong.JumlahKlaim.Kosong() {
		t.Errorf("jumlah kosong menjadi %q", kosong.JumlahKlaim.String())
	}
	if kosong.KodeStatus != "" {
		t.Errorf("kode status kosong menjadi %q", kosong.KodeStatus)
	}
	if kosong.Status().Diketahui() {
		t.Error("status kosong tidak boleh dianggap diketahui")
	}
	if !kosong.TanggalAkseptasi.IsZero() {
		t.Errorf("tanggal kosong menjadi %v", kosong.TanggalAkseptasi)
	}
}

// Kode di luar ketiga nilai spec dibawa apa adanya, tidak ditebak.
func TestKodeTidakDikenalDibawaApaAdanya(t *testing.T) {
	repo, bersihkan := siapkan(t)
	defer bersihkan()

	baris, err := repo.AmbilBaris(context.Background(), "UJI-KLAIM-1")
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range baris["UJI-PESERTA-2"] {
		if b.ID != "UJI-ADJ-4" {
			continue
		}
		if b.KodeStatus != "9" {
			t.Errorf("kode mentah = %q, mau %q", b.KodeStatus, "9")
		}
		if b.Status() != models.StatusTidakDiketahui {
			t.Errorf("kode 9 ditebak menjadi %v", b.Status())
		}
		return
	}
	t.Fatal("baris UJI-ADJ-4 tidak terbaca")
}

// Nama skema disebut eksplisit di setiap query (ADR-U-0033): pembacaan tetap
// berhasil walau skema bawaan sesi bukan skema sasaran.
func TestSkemaDisebutEksplisit(t *testing.T) {
	repo, bersihkan := siapkan(t)
	defer bersihkan()

	ctx, batal := context.WithTimeout(context.Background(), 20*time.Second)
	defer batal()
	if _, err := repo.AmbilHeader(ctx, "UJI-KLAIM-1"); err != nil {
		t.Fatalf("pembacaan berkualifikasi skema gagal: %v", err)
	}
}
