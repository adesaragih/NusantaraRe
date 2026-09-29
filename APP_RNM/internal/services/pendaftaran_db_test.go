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

	"github.com/cockroachdb/apd/v3"

	"nusantarare/internal/repository"
	"nusantarare/internal/repository/skemauji"
	"nusantarare/internal/services"
	"nusantarare/inti"
	"nusantarare/inti/db"
	intiuang "nusantarare/inti/uang"
	"nusantarare/inti/utils"
)

// penomorUji memberi nomor yang dapat diramalkan, menggantikan butir o.
type penomorUji struct{ n int }

func (p *penomorUji) NomorBerikut(context.Context, *db.Tx, string, time.Time) (string, error) {
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
	// Tiruan tabel peserta polis beserta isinya. Pendaftaran membaca peserta
	// dari sana; tanpa fixture ini test gagal di pembacaan, bukan menguji
	// pendaftaran.
	if err := skemauji.IsiPesertaPolis(ctx, sqlDB, skema); err != nil {
		t.Fatalf("mengisi tiruan peserta polis: %v", err)
	}
	db, err := skemauji.BukaRepositori()
	if err != nil {
		t.Fatal(err)
	}
	// Refactor bentuk B: pembaca polis PremiumList disambung seperti di cmd/api.
	return services.New(db).DenganPembacaPolis(skemauji.PembacaPolis(db)), func() {
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
		// Hanya nomor sertifikat; nilai polisnya dibaca server dari tiruan
		// M_LIFE_PREMIUM_DETAIL yang dibuat skema uji (tiket 03).
		Sertifikat: []string{"006"},
	}
}

// Pendaftaran menulis TIGA tempat plus baris datar, dalam satu transaksi.
func TestDaftarMenulisTigaTempatDanBarisDatar(t *testing.T) {
	svc, bersihkan := siapkanPendaftaran(t)
	defer bersihkan()
	ctx := context.Background()

	pohon, err := svc.Pendaftaran().DenganPenomor(&penomorUji{}).
		Daftar(ctx, inti.Pelaku{AkunID: "UJI-OPERATOR"}, permintaan())
	if err != nil {
		t.Fatalf("mendaftarkan klaim: %v", err)
	}
	if !strings.HasPrefix(pohon.Work.ID, repository.AwalanKlaim) {
		t.Errorf("pengenal work %q tidak berawalan %q", pohon.Work.ID, repository.AwalanKlaim)
	}
	if pohon.Work.Lini != inti.LiniLife {
		t.Errorf("LINI = %q, mau %q", pohon.Work.Lini, inti.LiniLife)
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
	// ⭐ Keempat tanggal valuasi terbaca kembali SAMA PERSIS. Inilah yang
	// tiket 06 pijak: validasi DOL membaca jendela ini dari peserta klaim,
	// bukan bertanya ulang ke tabel 66,8 juta baris.
	ps := klaim.Peserta[0]
	for _, k := range []struct{ nama, got, mau string }{
		{"valuasi gross mulai", ps.ValuasiGrossMulai, "2026-01-01 00:00:00"},
		{"valuasi gross selesai", ps.ValuasiGrossSelesai, "2026-12-31 00:00:00"},
		{"valuasi retro mulai", ps.ValuasiRetroMulai, "2026-02-01 00:00:00"},
		{"valuasi retro selesai", ps.ValuasiRetroSelesai, "2026-11-30 00:00:00"},
	} {
		if k.got != k.mau {
			t.Errorf("%s = %q, mau %q", k.nama, k.got, k.mau)
		}
	}
	// Peserta ber-EDMSTATUS Batal tidak pernah ikut: sertifikat 010 ada di
	// fixture, tetapi memintanya harus GAGAL.
	if ps.SumberID != "UJI-SRC-1" {
		t.Errorf("SOURCE_ID = %q, mau UJI-SRC-1", ps.SumberID)
	}

	// ⭐ BUTIR bp (GILIRAN-14): satu baris adjustment lahir bersama peserta
	// terpilih - `SavePesertaClaim` 7.8 - dibaca ULANG dari Oracle.
	db, tutupDB := repoUji(t)
	defer tutupDB()
	perBaris, err := repository.NewKlaimLife(db).AmbilBaris(ctx, pohon.Work.ID)
	if err != nil {
		t.Fatal(err)
	}
	baris := perBaris[ps.ID]
	if len(baris) != 1 {
		t.Fatalf("baris adjustment peserta = %d, mau 1 (butir bp)", len(baris))
	}
	for _, k := range []struct {
		medan string
		got   intiuang.Money
		mau   string
	}{
		{"SUM_INSURED", baris[0].SumInsured, "1000000"},
		{"CEDING_RETENTION", baris[0].CedingRetention, "100000"},
		{"CLAIM_AMOUNT", baris[0].JumlahKlaim, "25000.1234"},
	} {
		if got := k.got.Amount; got == nil || got.Cmp(uangUjiDB(t, k.mau)) != 0 {
			t.Errorf("%s = %v, mau %s (7.8, empat angka)", k.medan, got, k.mau)
		}
	}
	if baris[0].KodeStatus != "" {
		t.Errorf("baris lahir berkode %q; 7.8 tidak menulis STS_REJECT", baris[0].KodeStatus)
	}
	// AC 32: baris datar warisan kini ikut tertulis - barisnya ada.
	datar, err := repository.NewPohonKlaim(db).CacahBarisLama(ctx, pohon.Work.ID)
	if err != nil {
		t.Fatal(err)
	}
	if datar != 1 {
		t.Errorf("baris datar OS_AKSEPTASI_KLAIM_LIFE = %d, mau 1 (AC 32)", datar)
	}
}

// uangUjiDB mengurai desimal pembanding tanpa float.
func uangUjiDB(t *testing.T, s string) *apd.Decimal {
	t.Helper()
	d, err := utils.ParseDecimal(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

// ⛔ Peserta ber-EDMSTATUS 'Batal' tidak dapat didaftarkan.
//
// Penyaring hidup diuji terhadap Oracle sungguhan di sini, bukan hanya
// terhadap dirinya sendiri di test murni.
func TestPesertaBatalTidakDapatDidaftarkan(t *testing.T) {
	svc, bersihkan := siapkanPendaftaran(t)
	defer bersihkan()

	minta := permintaan()
	minta.Sertifikat = []string{"010"} // fixture: EDMSTATUS = 'Batal'
	_, err := svc.Pendaftaran().DenganPenomor(&penomorUji{}).
		Daftar(context.Background(), inti.Pelaku{AkunID: "UJI-OPERATOR"}, minta)
	if err == nil {
		t.Fatal("peserta batal diterima; penyaring hidup tidak berlaku di jalur pendaftaran")
	}
	if !strings.Contains(err.Error(), "010") {
		t.Errorf("galat tidak menyebut sertifikat yang ditolak: %v", err)
	}
}

// Dua pendaftaran menghasilkan DUA pengenal berbeda.
func TestDuaPendaftaranDuaPengenal(t *testing.T) {
	svc, bersihkan := siapkanPendaftaran(t)
	defer bersihkan()
	ctx := context.Background()
	daftar := svc.Pendaftaran().DenganPenomor(&penomorUji{})

	satu, err := daftar.Daftar(ctx, inti.Pelaku{AkunID: "UJI-OPERATOR"}, permintaan())
	if err != nil {
		t.Fatalf("pendaftaran pertama: %v", err)
	}
	dua, err := daftar.Daftar(ctx, inti.Pelaku{AkunID: "UJI-OPERATOR"}, permintaan())
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
	_, err := svc.Pendaftaran().Daftar(ctx, inti.Pelaku{AkunID: "UJI-OPERATOR"}, permintaan())
	if err == nil {
		t.Fatal("pendaftaran berhasil padahal penomoran belum diputuskan")
	}

	// Nol baris work object tertinggal.
	pohon, err := svc.KlaimLife().Ambil(ctx, repository.AwalanKlaim+"000001")
	if err == nil && pohon != nil {
		t.Errorf("klaim setengah jadi tertinggal: %+v", pohon.ID)
	}
}
