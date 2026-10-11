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
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
	intiuang "nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/claimlife/backend/repository"
	"nusantarare/modul/claimlife/backend/services"
	"nusantarare/uji/skemauji"
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
	lengkapiKolomVersi(t, ctx, sqlDB, skema)
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

// lengkapiKolomVersi - salinan `repository/versiterakhir_db_test.go` (paket uji berbeda): AmbilUntukKlaim
// memeringkat versi lewat `PL_NUMBER_EDM` dan `TGL_INPUT` (OQ-N14), yang belum dimuat tiruan bersama
// `uji/skemauji` (di luar folder modul ini). Kolomnya ditambahkan ke TIRUAN skema uji - bukan DDL tabel
// warisan; ORA-01430 = kolom sudah ada.
func lengkapiKolomVersi(t *testing.T, ctx context.Context, sqlDB *sql.DB, skema string) {
	t.Helper()
	for _, kolom := range []string{"PL_NUMBER_EDM VARCHAR2(100)", "TGL_INPUT DATE"} {
		if _, err := sqlDB.ExecContext(ctx, `ALTER TABLE `+skema+`.M_LIFE_PREMIUM_DETAIL ADD (`+kolom+`)`); err != nil &&
			!strings.Contains(err.Error(), "ORA-01430") {
			t.Fatalf("melengkapi tiruan peserta (%s): %v", kolom, err)
		}
	}
}

func permintaan() services.PermintaanDaftar {
	return services.PermintaanDaftar{
		NomorPremiList: "UJI-PL-1",
		NomorPolis:     "UJI-POL-0001",
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

	semua, err := svc.Pendaftaran().Daftar(ctx, inti.Pelaku{AkunID: "UJI-OPERATOR"}, permintaan())
	if err != nil {
		t.Fatalf("mendaftarkan klaim: %v", err)
	}
	if len(semua) != 1 {
		t.Fatalf("case terbentuk %d, mau 1 (satu PL, satu mata uang)", len(semua))
	}
	pohon := semua[0]
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
	// RALAT 10-10-2026: Register tidak menomori - CLAIM_NO kosong sampai Save to RNM.
	if klaim.NomorKlaim != "" {
		t.Errorf("nomor klaim = %q, mau kosong sampai Save to RNM", klaim.NomorKlaim)
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
	// RALAT 10-10-2026 (keputusan work owner): Register TIDAK menulis baris datar OS_AKSEPTASI_KLAIM_LIFE - ia lahir di
	// Save to RNM bersama nomor klaim. Dulu AC 32 menuntut satu baris di sini.
	datar, err := repository.NewPohonKlaim(db).CacahBarisLama(ctx, pohon.Work.ID)
	if err != nil {
		t.Fatal(err)
	}
	if datar != 0 {
		t.Errorf("Register menulis %d baris datar OS_AKSEPTASI_KLAIM_LIFE, mau 0 (RALAT 10-10-2026)", datar)
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
	_, err := svc.Pendaftaran().Daftar(context.Background(), inti.Pelaku{AkunID: "UJI-OPERATOR"}, minta)
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
	daftar := svc.Pendaftaran()

	satu, err := daftar.Daftar(ctx, inti.Pelaku{AkunID: "UJI-OPERATOR"}, permintaan())
	if err != nil {
		t.Fatalf("pendaftaran pertama: %v", err)
	}
	dua, err := daftar.Daftar(ctx, inti.Pelaku{AkunID: "UJI-OPERATOR"}, permintaan())
	if err != nil {
		t.Fatalf("pendaftaran kedua: %v", err)
	}
	if len(satu) != 1 || len(dua) != 1 || satu[0].Work.ID == dua[0].Work.ID {
		t.Errorf("dua pendaftaran memakai pengenal yang sama: %+v %+v", satu, dua)
	}
}

// RALAT 10-10-2026: pendaftaran TIDAK bergantung penomor (dulu: penomor yang gagal membatalkan seluruh transaksi) -
// Register berhasil tanpa penomor apa pun, CLAIM_NO kosong, pengenal berawalan CLMLF-.
func TestRegisterTanpaPenomor(t *testing.T) {
	svc, bersihkan := siapkanPendaftaran(t)
	defer bersihkan()
	ctx := context.Background()
	semua, err := svc.Pendaftaran().Daftar(ctx, inti.Pelaku{AkunID: "UJI-OPERATOR"}, permintaan())
	if err != nil || len(semua) != 1 {
		t.Fatalf("pendaftaran: %v", err)
	}
	if semua[0].Klaim.NomorKlaim != "" || !strings.HasPrefix(semua[0].Work.ID, "CLMLF-") {
		t.Errorf("nomor %q, pengenal %q", semua[0].Klaim.NomorKlaim, semua[0].Work.ID)
	}
}

// Aturan in-app Register (keputusan work owner 10-10-2026) terhadap Oracle uji: peserta 006 sudah terdaftar di klaim
// pertama -> pendaftaran kedua ditolak 409 (Outstanding, baris peserta NULL), lalu Accepted ('1') tetap ditolak,
// Rejected ('2') boleh, baris tercabut (STS_HAPUS) boleh.
func TestDaftarPesertaAktifDiKlaimLain(t *testing.T) {
	svc, bersihkan := siapkanPendaftaran(t)
	defer bersihkan()
	ctx := context.Background()
	daftar := svc.Pendaftaran()
	pelaku := inti.Pelaku{AkunID: "UJI-OPERATOR", Peran: []string{services.PeranInputRegister}}

	pertama, err := daftar.Daftar(ctx, pelaku, permintaan())
	if err != nil || len(pertama) != 1 {
		t.Fatalf("pendaftaran pertama: %v", err)
	}
	klaimID := pertama[0].Work.ID
	// ubahStatus menulis kode ke baris peserta DAN baris adjustment-nya (seperti PerbaruiStatusBaris); kolom STS_HAPUS
	// hanya ke peserta.
	ubahStatus := func(kolom, nilai string) {
		t.Helper()
		pes, err := svc.DB().Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := svc.DB().ExecContext(ctx, "UPDATE "+pes+" SET "+kolom+" = :1 WHERE CLAIM_ID = :2", nilai, klaimID); err != nil {
			t.Fatal(err)
		}
		if kolom != "STS_REJECT" {
			return
		}
		adj, err := svc.DB().Qualify("T_CLAIMLF_ADJUSTMENT")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := svc.DB().ExecContext(ctx, "UPDATE "+adj+" SET STS_REJECT = :1 WHERE PREMIUM_LIST_DETAIL_ID IN "+
			"(SELECT ID FROM "+pes+" WHERE CLAIM_ID = :2)", nilai, klaimID); err != nil {
			t.Fatal(err)
		}
	}
	harusDitolak := func(label string) {
		t.Helper()
		_, err := daftar.Daftar(ctx, pelaku, permintaan())
		if !errors.Is(err, services.ErrPesertaSudahDiklaim) || !strings.Contains(err.Error(), "(PL UJI-PL-1)") ||
			!strings.Contains(err.Error(), "("+label+")") {
			t.Fatalf("%s: galat %v", label, err)
		}
	}
	harusDitolak("Outstanding")
	ubahStatus("STS_REJECT", kontrak.KodeOutstanding)
	harusDitolak("Outstanding")
	ubahStatus("STS_REJECT", kontrak.KodeAksep)
	harusDitolak("Accepted")
	ubahStatus("STS_REJECT", kontrak.KodeDitolak)
	kedua, err := daftar.Daftar(ctx, pelaku, permintaan())
	if err != nil || len(kedua) != 1 {
		t.Fatalf("peserta Rejected harus boleh didaftarkan lagi: %v", err)
	}
	// Klaim kedua kini aktif; cabut pesertanya (STS_HAPUS) -> boleh lagi.
	klaimID = kedua[0].Work.ID
	ubahStatus("STS_HAPUS", "1")
	if _, err := daftar.Daftar(ctx, pelaku, permintaan()); err != nil {
		t.Fatalf("baris tercabut harus tidak dihitung: %v", err)
	}
}
