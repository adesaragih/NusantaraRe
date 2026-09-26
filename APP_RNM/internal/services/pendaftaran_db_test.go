//go:build db

package services_test

// Pendaftaran klaim terhadap skema uji Oracle - tiket 02.
//
// ⛔ Seluruh test di sini MELEWATI dengan pesan bila ORACLE_DSN belum
// dikonfigurasi. Melewati bukan lulus, dan laporan sesi menyebutnya apa adanya.
//
// Dibaca sesudah: services/pendaftaran.go.

import (
	"context"
	"strings"
	"testing"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/repository/skemauji"
	"nusantarare/internal/services"
)

// penomorUji memberi nomor yang dapat diramalkan, menggantikan butir o.
type penomorUji struct{ n int }

func (p *penomorUji) NomorBerikut(context.Context, *repository.Tx, string, time.Time) (string, error) {
	p.n++
	return "UJI-CLM-NOMOR", nil
}

func siapkanPendaftaran(t *testing.T) (*services.Service, func()) {
	t.Helper()
	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	ctx := context.Background()
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	if err := skemauji.Pasang(ctx, sqlDB, skema); err != nil {
		t.Fatalf("memasang skema uji: %v", err)
	}
	db, err := skemauji.BukaRepositori()
	if err != nil {
		t.Fatal(err)
	}
	return services.New(db), func() {
		_ = skemauji.Bongkar(ctx, sqlDB, skema)
		_ = db.Close()
		_ = sqlDB.Close()
	}
}

func permintaan() services.PermintaanDaftar {
	return services.PermintaanDaftar{
		NomorPremiList: "UJI-PL-1",
		NomorPolis:     "UJI-POL-0001",
		Type:           "QP",
		KodeBisnis:     "L1",
		MataUang:       "IDR",
		Peserta: []models.Peserta{
			{NomorPremiList: "UJI-PL-1", NomorPolis: "UJI-POL-0001",
				NomorSertifikat: "006", MataUang: "IDR"},
		},
	}
}

// Pendaftaran menulis TIGA tempat plus baris datar, dalam satu transaksi.
func TestDaftarMenulisTigaTempatDanBarisDatar(t *testing.T) {
	svc, bersihkan := siapkanPendaftaran(t)
	defer bersihkan()
	ctx := context.Background()

	pohon, err := svc.Pendaftaran().DenganPenomor(&penomorUji{}).
		Daftar(ctx, services.Pelaku{AkunID: "UJI-OPERATOR"}, permintaan())
	if err != nil {
		t.Fatalf("mendaftarkan klaim: %v", err)
	}
	if !strings.HasPrefix(pohon.Work.ID, repository.AwalanKlaim) {
		t.Errorf("pengenal work %q tidak berawalan %q", pohon.Work.ID, repository.AwalanKlaim)
	}
	if pohon.Work.Lini != models.LiniLife {
		t.Errorf("LINI = %q, mau %q", pohon.Work.Lini, models.LiniLife)
	}

	// Header terbaca kembali, beserta mata uangnya (butir z1).
	klaim, err := svc.KlaimLife().Ambil(ctx, pohon.Work.ID)
	if err != nil {
		t.Fatalf("membaca klaim: %v", err)
	}
	if klaim.NomorKlaim != "UJI-CLM-NOMOR" {
		t.Errorf("nomor klaim = %q", klaim.NomorKlaim)
	}
	if klaim.ClaimRetro.Currency != "IDR" {
		t.Errorf("mata uang header = %q, mau IDR; butir z1 belum sampai ke jalur baca",
			klaim.ClaimRetro.Currency)
	}
	if len(klaim.Peserta) != 1 {
		t.Fatalf("peserta tersimpan %d, mau 1", len(klaim.Peserta))
	}
	// ⛔ Nomor sertifikat berawalan nol tetap utuh (ADR-U-0022).
	if klaim.Peserta[0].NomorSertifikat != "006" {
		t.Errorf("nomor sertifikat = %q, mau %q", klaim.Peserta[0].NomorSertifikat, "006")
	}
}

// Dua pendaftaran menghasilkan DUA pengenal berbeda.
func TestDuaPendaftaranDuaPengenal(t *testing.T) {
	svc, bersihkan := siapkanPendaftaran(t)
	defer bersihkan()
	ctx := context.Background()
	daftar := svc.Pendaftaran().DenganPenomor(&penomorUji{})

	satu, err := daftar.Daftar(ctx, services.Pelaku{AkunID: "UJI-OPERATOR"}, permintaan())
	if err != nil {
		t.Fatalf("pendaftaran pertama: %v", err)
	}
	dua, err := daftar.Daftar(ctx, services.Pelaku{AkunID: "UJI-OPERATOR"}, permintaan())
	if err != nil {
		t.Fatalf("pendaftaran kedua: %v", err)
	}
	if satu.Work.ID == dua.Work.ID {
		t.Errorf("dua pendaftaran memakai pengenal yang sama: %q", satu.Work.ID)
	}
}

// ⛔ Penomoran yang gagal MEMBATALKAN seluruh transaksi.
//
// Tanpa ini, klaim setengah jadi tertinggal di basis data: baris work object
// ada, headernya tidak, dan tidak ada yang melaporkannya.
func TestPenomorGagalMembatalkanSeluruhTransaksi(t *testing.T) {
	svc, bersihkan := siapkanPendaftaran(t)
	defer bersihkan()
	ctx := context.Background()

	// Penomor bawaan selalu gagal selama butir o belum diputuskan.
	_, err := svc.Pendaftaran().Daftar(ctx, services.Pelaku{AkunID: "UJI-OPERATOR"}, permintaan())
	if err == nil {
		t.Fatal("pendaftaran berhasil padahal penomoran belum diputuskan")
	}

	// Nol baris work object tertinggal.
	pohon, err := svc.KlaimLife().Ambil(ctx, repository.AwalanKlaim+"000001")
	if err == nil && pohon != nil {
		t.Errorf("klaim setengah jadi tertinggal: %+v", pohon.ID)
	}
}
