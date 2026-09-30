//go:build db

// Seam `repository` terhadap skema uji Oracle NYATA - tiket 14.
//
// Untuk apa berkas ini: membuktikan bahwa migrasi benar-benar membangun tujuh
// tabel dengan relasinya, bahwa satu pohon klaim dapat ditulis dan dibaca
// kembali sampai tingkat terdalam, dan bahwa menghapus klaim mengkaskade sampai
// cicit.
//
// Jalankan: make test-db   (perlu ORACLE_DSN dan ORACLE_SCHEMA)
//
// Tanpa instance Oracle, seluruh test di sini MELEWATI dengan pesan - bukan
// lulus diam-diam.
package repository_test

import (
	"context"
	"testing"

	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/migrasi"
	intiuang "nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
	claimlife "nusantarare/modul/claimlife/backend"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
	"nusantarare/uji/skemauji"
)

func siapkanPohon(t *testing.T) (*db.DB, *repository.PohonKlaim, func()) {
	t.Helper()
	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		// Hanya "Oracle belum dikonfigurasi" yang dilewati. Salah
		// konfigurasi, menunjuk produksi, atau menunjuk skema yang bukan
		// skema uji harus MENGGAGALKAN - test ini menghapus tabel.
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
	return db, repository.NewPohonKlaim(db), func() {
		_ = skemauji.Bongkar(ctx, sqlDB, skema)
		_ = db.Close()
		_ = sqlDB.Close()
	}
}

// contohPohon membuat satu klaim lengkap sampai tingkat terdalam.
// ⛔ Nol nama orang, nol nomor polis nyata.
func contohPohon(t *testing.T) models.PohonKlaim {
	t.Helper()
	uang := func(s string) intiuang.Money {
		m, err := intiuang.NewMoney(s, "IDR")
		if err != nil {
			t.Fatal(err)
		}
		return m
	}
	rasio := func(s string) intiuang.Ratio {
		r, err := intiuang.NewRatio(s, 8)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}
	return models.PohonKlaim{
		Work: models.WorkClaim{
			ID: "CLM-UJI900", Lini: inti.LiniLife, Type: "UJI-TYPE", CaseID: "UJI-CASE-900",
		},
		Klaim: models.Klaim{
			ID: "CLM-UJI900", NomorKlaim: "UJI-CLM-9", NomorPolis: "UJI-POL-9",
			NamaBisnis: "UJI BISNIS", KodeStatus: "0",
			Peserta: []models.Peserta{{
				ID: "UJI-P-1", NomorSertifikat: "006", MataUang: "IDR",
				Baris: []models.BarisAdjustment{{
					ID: "UJI-A-1", KodeStatus: "1", JumlahKlaim: uang("1234567890.12345678"),
					Spreading: []models.Spreading{{
						ID: "UJI-S-1", TreatyTypeName: "UJI-TREATY", TreatyYearLife: "2026",
						IDR: uang("500000.5"), Currency: "IDR",
						// RetrocadedShare UANG (diralat 26-09-2026 menurut
						// SpreadingClaimLife_Act langkah 8.2.1.4-7); Rate rasio.
						// Keduanya sengaja bertipe berbeda.
						RetrocadedShare: uang("0.12345678"),
						Rate:            rasio("0.075"),
						Retro: []models.SpreadingRetro{
							{ID: "UJI-RT-1", ReinsurerName: "UJI-REINSURER",
								Amount: uang("250000.25"), PercentShare: rasio("0.6"),
								Rate: rasio("0.0125"), PremiumSpreadedGross: uang("99999.99999999"),
								PremiumSpreadedNet: uang("88888.88888888"),
								// COMMISION dan OVR_COMM adalah PERSEN, bukan uang.
								Commision: rasio("1234.5"), OvrComm: rasio("0.00000001")},
							{ID: "UJI-RT-2", ReinsurerName: "UJI-REINSURER-2",
								Amount: uang("0.00000001"), PercentShare: rasio("0.4")},
						},
					}},
				}},
			}},
		},
	}
}

// Migrasi idempoten: menjalankannya dua kali tidak menambah apa pun.
func TestMigrasiIdempoten(t *testing.T) {
	db, _, bersihkan := siapkanPohon(t)
	defer bersihkan()

	lap, err := migrasi.Jalankan(context.Background(), db, skemauji.SumberMigrasi()...)
	if err != nil {
		t.Fatalf("migrasi kedua gagal: %v", err)
	}
	if len(lap.Dijalankan) != 0 {
		t.Errorf("migrasi kedua menjalankan ulang %v - tidak idempoten", lap.Dijalankan)
	}
	if len(lap.Dilewati) == 0 {
		t.Error("tidak satu pun langkah dilaporkan dilewati")
	}
}

// Langkah yang dulu gagal SEPARUH JALAN tetap dapat dituntaskan.
//
// DDL Oracle menutup transaksinya sendiri. Bila sebuah langkah gagal di
// pernyataan kedua, pernyataan pertamanya sudah terlanjur jadi sementara
// T_MIGRASI tidak mencatat apa pun - sehingga percobaan berikutnya mengulang
// langkah itu dari awal dan mati di ORA-00955. Test ini menirukan keadaan itu:
// skema dibongkar habis, satu objek dibuat sendirian, lalu migrasi harus LULUS
// dan melaporkan objek yang dilewatinya.
func TestLangkahGagalSeparuhJalanTetapSelesai(t *testing.T) {
	db, _, bersihkan := siapkanPohon(t)
	defer bersihkan()
	ctx := context.Background()

	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		// Hanya "Oracle belum dikonfigurasi" yang dilewati. Salah
		// konfigurasi, menunjuk produksi, atau menunjuk skema yang bukan
		// skema uji harus MENGGAGALKAN - test ini menghapus tabel.
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	if _, err := migrasi.Bongkar(ctx, db, skemauji.SumberMigrasi()...); err != nil {
		t.Fatalf("membongkar: %v", err)
	}
	// ⛔ Pernyataan yang dijalankan adalah pernyataan PERTAMA langkah 001 yang
	// sesungguhnya - bukan tabel bikinan sendiri yang bentuknya berbeda.
	// Menirukan kegagalan separuh jalan dengan objek berbentuk lain akan
	// menguji hal lain, dan lebih buruk: ia akan membuat skema cacat tercatat
	// sebagai migrasi yang sukses.
	pernyataan, err := migrasi.PernyataanLangkah(claimlife.SumberMigrasi(), "001_t_work_claim.sql")
	if err != nil {
		t.Fatal(err)
	}
	if len(pernyataan) == 0 {
		t.Fatal("langkah 001 tidak memuat pernyataan")
	}
	if _, err := sqlDB.ExecContext(ctx,
		strings.ReplaceAll(pernyataan[0], "{skema}", skema)); err != nil {
		t.Fatalf("membuat sisa objek: %v", err)
	}

	lap, err := migrasi.Jalankan(ctx, db, skemauji.SumberMigrasi()...)
	if err != nil {
		t.Fatalf("migrasi menolak meneruskan langkah yang separuh jadi: %v", err)
	}
	if len(lap.ObjekSudahAda) == 0 {
		t.Error("objek yang sudah berdiri tidak dilaporkan; ia dilewati diam-diam")
	}
	if len(lap.Dijalankan) == 0 {
		t.Error("tidak satu pun langkah dituntaskan")
	}
}

// Jalur mundur benar-benar membongkar, dan migrasi dapat dijalankan lagi
// sesudahnya.
func TestJalurMundurDiuji(t *testing.T) {
	db, _, bersihkan := siapkanPohon(t)
	defer bersihkan()
	ctx := context.Background()

	if _, err := migrasi.Bongkar(ctx, db, skemauji.SumberMigrasi()...); err != nil {
		t.Fatalf("jalur mundur gagal: %v", err)
	}
	lap, err := migrasi.Jalankan(ctx, db, skemauji.SumberMigrasi()...)
	if err != nil {
		t.Fatalf("migrasi ulang sesudah mundur gagal: %v", err)
	}
	if len(lap.Dijalankan) == 0 {
		t.Error("sesudah dibongkar, migrasi ulang tidak menjalankan apa pun")
	}
}

// Satu pohon utuh ditulis dalam SATU transaksi, ke tabel relasional DAN ke
// baris datar warisan.
func TestSimpanPohonMenulisDuaTempatDalamSatuTransaksi(t *testing.T) {
	db, repo, bersihkan := siapkanPohon(t)
	defer bersihkan()
	ctx := context.Background()
	p := contohPohon(t)

	tx, err := db.Mulai(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Simpan(ctx, tx, p); err != nil {
		_ = tx.Rollback()
		t.Fatalf("menyimpan pohon: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	// Penulisan kedua benar-benar terjadi: satu baris datar per baris adjustment.
	n, err := repo.CacahBarisLama(ctx, p.Work.CaseID)
	if err != nil {
		t.Fatal(err)
	}
	if n != p.Klaim.CacahBaris() {
		t.Errorf("baris datar %d, mau %d", n, p.Klaim.CacahBaris())
	}
}

// Pembacaan turun sampai CICIT - tingkat keenam pohon.
func TestBacaSampaiCicit(t *testing.T) {
	db, repo, bersihkan := siapkanPohon(t)
	defer bersihkan()
	ctx := context.Background()
	p := contohPohon(t)

	tx, _ := db.Mulai(ctx)
	if err := repo.Simpan(ctx, tx, p); err != nil {
		_ = tx.Rollback()
		t.Fatal(err)
	}
	_ = tx.Commit()

	// Sesudah hapus, SETIAP tingkat wajib nol - diperiksa satu per satu,
	// bukan lewat satu angka total yang dapat menutupi satu tingkat yang
	// tertinggal.
	sesudah, err := repo.Dampak(ctx, p.Work.ID, p.Work.CaseID)
	if err != nil {
		t.Fatalf("mencacah dampak sesudah hapus: %v", err)
	}
	for _, k := range []struct {
		nama string
		got  int
	}{
		{"header", sesudah.Header}, {"peserta", sesudah.Peserta},
		{"adjustment", sesudah.Adjustment}, {"spreading", sesudah.Spreading},
		{"spreading retro", sesudah.SpreadingRetro},
		{"dokumen", sesudah.Dokumen}, {"baris work", sesudah.WorkClaim},
		{"baris datar warisan", sesudah.BarisDatarWarisan},
	} {
		if k.got != 0 {
			t.Errorf("sesudah hapus masih ada %d %s", k.got, k.nama)
		}
	}

	spr, err := repo.AmbilSpreading(ctx, p.Work.ID)
	if err != nil {
		t.Fatal(err)
	}
	daftar := spr["UJI-A-1"]
	if len(daftar) != 1 {
		t.Fatalf("spreading %d, mau 1", len(daftar))
	}
	if got := utils.FormatDecimal(daftar[0].IDR.Amount); got != "500000.5" {
		t.Errorf("IDR = %q, mau 500000.5", got)
	}
	if len(daftar[0].Retro) != 2 {
		t.Fatalf("spreading retro %d, mau 2 - cicit tidak terbaca", len(daftar[0].Retro))
	}
	// Uang di tingkat terdalam tetap utuh, digit demi digit.
	if got := utils.FormatDecimal(daftar[0].Retro[1].Amount.Amount); got != "0.00000001" {
		t.Errorf("jumlah retro = %q, mau 0.00000001", got)
	}

	// ⛔ SELURUH kolom uang dan rasio diperiksa pulang-pergi, satu per satu.
	// Ronde 1 hanya membaca IDR/USD dan AMOUNT, sehingga tujuh kolom lain
	// pulang kosong tanpa ada satu pun test yang menyadarinya.
	sprd := daftar[0]
	rt := sprd.Retro[0]
	for _, k := range []struct {
		nama, mau, dapat string
	}{
		{"spreading RETROCADED_SHARE", "0.12345678", utils.FormatDecimal(sprd.RetrocadedShare.Amount)},
		{"spreading RATE", "0.075", utils.FormatDecimal(sprd.Rate.Value)},
		{"retro PERCENT_SHARE", "0.6", utils.FormatDecimal(rt.PercentShare.Value)},
		{"retro RATE", "0.0125", utils.FormatDecimal(rt.Rate.Value)},
		{"retro PREMIUM_SPREADED_GROSS", "99999.99999999", utils.FormatDecimal(rt.PremiumSpreadedGross.Amount)},
		{"retro PREMIUM_SPREADED_NET", "88888.88888888", utils.FormatDecimal(rt.PremiumSpreadedNet.Amount)},
		{"retro COMMISION", "1234.5", utils.FormatDecimal(rt.Commision.Value)},
		{"retro OVR_COMM", "0.00000001", utils.FormatDecimal(rt.OvrComm.Value)},
	} {
		if k.dapat != k.mau {
			t.Errorf("%s = %q, mau %q", k.nama, k.dapat, k.mau)
		}
	}

	// Kolom yang memang KOSONG di fixture pulang kosong - bukan menjadi nol.
	if !sprd.Retro[1].Commision.Kosong() {
		t.Errorf("COMMISION yang tidak diisi pulang sebagai %v, seharusnya kosong",
			sprd.Retro[1].Commision)
	}
}

// Menghapus klaim mengkaskade sampai tingkat terdalam. Test memeriksa CICIT
// ikut hilang, bukan hanya anak langsungnya.
func TestHapusMengkaskadeSampaiCicit(t *testing.T) {
	db, repo, bersihkan := siapkanPohon(t)
	defer bersihkan()
	ctx := context.Background()
	p := contohPohon(t)

	tx, _ := db.Mulai(ctx)
	if err := repo.Simpan(ctx, tx, p); err != nil {
		_ = tx.Rollback()
		t.Fatal(err)
	}
	_ = tx.Commit()

	// ⭐ Dampak dihitung SEBELUM menghapus, lalu dibandingkan dengan yang
	// benar-benar hilang. Itulah satu-satunya cara membuktikan AC "jumlah yang
	// ditampilkan popup sama persis dengan yang benar-benar terhapus" -
	// membandingkan angka dengan dirinya sendiri tidak membuktikan apa pun.
	sebelum, err := repo.Dampak(ctx, p.Work.ID, p.Work.CaseID)
	if err != nil {
		t.Fatalf("mencacah dampak: %v", err)
	}
	for _, k := range []struct {
		nama string
		got  int
	}{
		{"header", sebelum.Header}, {"peserta", sebelum.Peserta},
		{"adjustment", sebelum.Adjustment}, {"spreading", sebelum.Spreading},
		{"spreading retro", sebelum.SpreadingRetro}, {"baris work", sebelum.WorkClaim},
	} {
		if k.got == 0 {
			t.Errorf("dampak %s = 0 sebelum hapus; fixture-nya yang kosong, "+
				"sehingga perbandingan sesudah hapus tidak membuktikan apa pun", k.nama)
		}
	}

	tx2, _ := db.Mulai(ctx)
	if err := repo.HapusFisik(ctx, tx2, p.Work.ID, p.Work.CaseID); err != nil {
		_ = tx2.Rollback()
		t.Fatalf("menghapus pohon: %v", err)
	}
	if err := tx2.Commit(); err != nil {
		t.Fatal(err)
	}

	spr, err := repo.AmbilSpreading(ctx, p.Work.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(spr) != 0 {
		t.Errorf("sesudah hapus masih ada %d kelompok spreading - kaskade tidak sampai cicit", len(spr))
	}
	n, err := repo.CacahBarisLama(ctx, p.Work.CaseID)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("baris datar warisan tersisa %d", n)
	}
}

// Pembongkaran data lama berjalan di atas tabel tiruan warisan, dan jumlah
// baris sesudah migrasi sama dengan sebelumnya.
func TestBongkarDataLamaDariTabelTiruan(t *testing.T) {
	_, _, bersihkan := siapkanPohon(t)
	defer bersihkan()
	ctx := context.Background()

	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		// Hanya "Oracle belum dikonfigurasi" yang dilewati. Salah
		// konfigurasi, menunjuk produksi, atau menunjuk skema yang bukan
		// skema uji harus MENGGAGALKAN - test ini menghapus tabel.
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	masuk := []repository.BarisLama{
		{ID: "L1", CASEID: "UJI-CASE-800", NO_CLAIM: "UJI-CLM-8", CERTIFICATE_NO: "006",
			CURRENCY: "IDR", CLAIM_AMOUNT: "100.00000001", STS_REJECT: "0"},
		{ID: "L2", CASEID: "UJI-CASE-800", NO_CLAIM: "UJI-CLM-8", CERTIFICATE_NO: "006",
			CURRENCY: "IDR", CLAIM_AMOUNT: "200", STS_REJECT: "1"},
	}
	if err := skemauji.IsiBarisLama(ctx, sqlDB, skema, masuk); err != nil {
		t.Fatalf("mengisi tabel tiruan: %v", err)
	}

	// ⛔ Dibaca KEMBALI dari tabel, bukan dari slice di memori. Ronde 1
	// menulis ke tabel lalu mengabaikannya dan membongkar slice yang sama -
	// sehingga pembacaan dari Oracle tidak pernah teruji sama sekali.
	repo, err := skemauji.BukaRepositori()
	if err != nil {
		t.Fatalf("membuka repositori: %v", err)
	}
	defer func() { _ = repo.Close() }()

	kembali, err := repository.NewPohonKlaim(repo).AmbilBarisLama(ctx, "UJI-CASE-800")
	if err != nil {
		t.Fatalf("membaca baris lama: %v", err)
	}
	if len(kembali) != len(masuk) {
		t.Fatalf("baris terbaca %d, mau %d", len(kembali), len(masuk))
	}
	// Digit demi digit: uang pulang persis seperti yang ditulis.
	for i, b := range kembali {
		if b.CLAIM_AMOUNT != masuk[i].CLAIM_AMOUNT {
			t.Errorf("baris %s CLAIM_AMOUNT = %q, mau %q - pemisah desimal NLS?",
				b.ID, b.CLAIM_AMOUNT, masuk[i].CLAIM_AMOUNT)
		}
		if b.CERTIFICATE_NO != masuk[i].CERTIFICATE_NO {
			t.Errorf("baris %s CERTIFICATE_NO = %q, mau %q - nol di depan hilang?",
				b.ID, b.CERTIFICATE_NO, masuk[i].CERTIFICATE_NO)
		}
		if b.STS_REJECT != masuk[i].STS_REJECT {
			t.Errorf("baris %s STS_REJECT = %q, mau %q", b.ID, b.STS_REJECT, masuk[i].STS_REJECT)
		}
	}

	pohon, lap := repository.BongkarBarisLama(kembali)
	if len(pohon) != 1 {
		t.Fatalf("klaim terbentuk %d, mau 1", len(pohon))
	}
	if lap.AdjustmentTerbentuk != len(masuk) {
		t.Errorf("adjustment %d, mau %d", lap.AdjustmentTerbentuk, len(masuk))
	}
	if lap.BarisHardcode != 1 {
		t.Errorf("baris hardcode %d, mau 1", lap.BarisHardcode)
	}
}

// ⛔ Nama constraint yang bertabrakan harus MENGGAGALKAN migrasi, bukan dilewati.
//
// ORA-02264 berarti nama constraint sudah dipakai objek lain, dan Oracle baru
// memeriksanya ketika tabelnya belum ada. Jadi ORA-02264 pada sebuah
// CREATE TABLE berarti tabel itu JUSTRU TIDAK terbuat. Ronde 2 menelannya dan
// mencatat langkahnya sukses - skema tanpa tabel tercatat sebagai migrasi yang
// berhasil. Test ini menirukan tabrakan itu dan menuntut migrasi berhenti.
func TestNamaConstraintBertabrakanMenggagalkanMigrasi(t *testing.T) {
	db, _, bersihkan := siapkanPohon(t)
	defer bersihkan()
	ctx := context.Background()

	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		// Hanya "Oracle belum dikonfigurasi" yang dilewati. Salah
		// konfigurasi, menunjuk produksi, atau menunjuk skema yang bukan
		// skema uji harus MENGGAGALKAN - test ini menghapus tabel.
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	if _, err := migrasi.Bongkar(ctx, db, skemauji.SumberMigrasi()...); err != nil {
		t.Fatalf("membongkar: %v", err)
	}
	// Tabel lain yang sudah memakai nama constraint milik langkah 001.
	_, err = sqlDB.ExecContext(ctx, "CREATE TABLE "+skema+
		".UJI_TABRAKAN (ID VARCHAR2(1), CONSTRAINT PK_T_WORK_CLAIM PRIMARY KEY (ID))")
	if err != nil {
		t.Fatalf("membuat tabel penabrak: %v", err)
	}
	defer func() {
		_, _ = sqlDB.ExecContext(ctx, "DROP TABLE "+skema+".UJI_TABRAKAN CASCADE CONSTRAINTS")
	}()

	_, err = migrasi.Jalankan(ctx, db, skemauji.SumberMigrasi()...)
	if err == nil {
		t.Fatal("migrasi LULUS padahal nama constraint bertabrakan - tabel 001 tidak terbuat")
	}
	for _, wajib := range []string{"001", "ORA-02264"} {
		if !strings.Contains(err.Error(), wajib) {
			t.Errorf("galat %q tidak menyebut %q", err, wajib)
		}
	}

	// Langkah yang gagal TIDAK boleh tercatat selesai.
	var n int
	if err := sqlDB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM "+skema+".T_MIGRASI WHERE NAMA = :1",
		"001_t_work_claim").Scan(&n); err != nil {
		t.Fatalf("membaca T_MIGRASI: %v", err)
	}
	if n != 0 {
		t.Errorf("langkah 001 tercatat %d kali di T_MIGRASI padahal gagal", n)
	}
}

// Tabel yang sudah ada dengan BENTUK berbeda menggagalkan migrasi.
//
// `[keputusan work owner 26-09-2026, butir x]`. Ini kelemahan ronde 3 yang
// sengaja dibiarkan separuh: keberadaan objek dibuktikan, bentuknya tidak.
// Pada 26-09-2026 ia terbukti nyata - tabel warisan DOCUMENT_CLAIM sudah ada
// di POOLDATA dengan empat belas kolom yang bukan milik 007.
//
// Yang dituntut: migrasi GAGAL, galatnya menyebut nama tabel dan kolom yang
// berselisih, dan TIDAK SATU PUN langkah tercatat di T_MIGRASI - sebab
// pra-terbang berhenti sebelum satu pernyataan pun dikirim.
func TestBentukTabelBerbedaMenggagalkanMigrasi(t *testing.T) {
	db, _, bersihkan := siapkanPohon(t)
	defer bersihkan()
	ctx := context.Background()

	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	if _, err := migrasi.Bongkar(ctx, db, skemauji.SumberMigrasi()...); err != nil {
		t.Fatalf("membongkar: %v", err)
	}
	// Tabel bernama sama dengan yang dibuat 007, tetapi kolomnya sengaja
	// BUKAN kolom 007 - meniru keadaan DOCUMENT_CLAIM warisan di POOLDATA.
	//
	// ⚠️ Sesudah keputusan v1 tabel baru bernama T_CLAIMLF_DOCUMENT, sehingga
	// DOCUMENT_CLAIM warisan tidak lagi bertabrakan dengan migrasi mana pun.
	// Yang diuji di sini karena itu bukan tabrakan nama itu, melainkan PAGAR
	// BENTUK-nya: tabel apa pun yang sudah ada dengan kolom yang berbeda harus
	// menghentikan migrasi.
	_, err = sqlDB.ExecContext(ctx, "CREATE TABLE "+skema+
		".T_CLAIMLF_DOCUMENT (ID VARCHAR2(100) NOT NULL, IDPEGA VARCHAR2(100), NAMAFILE VARCHAR2(255))")
	if err != nil {
		t.Fatalf("membuat tiruan tabel warisan: %v", err)
	}
	defer func() {
		_, _ = sqlDB.ExecContext(ctx, "DROP TABLE "+skema+".T_CLAIMLF_DOCUMENT CASCADE CONSTRAINTS")
	}()

	_, err = migrasi.Jalankan(ctx, db, skemauji.SumberMigrasi()...)
	if err == nil {
		t.Fatal("migrasi LULUS padahal T_CLAIMLF_DOCUMENT berbentuk lain - " +
			"aplikasi akan berjalan di atas tabel yang kolomnya bukan miliknya")
	}
	// Pesannya harus menyebut tabelnya, salah satu kolom yang tidak diminta,
	// dan salah satu kolom yang diminta tetapi tidak ada.
	for _, wajib := range []string{"T_CLAIMLF_DOCUMENT", "IDPEGA", "BENTUKNYA BERBEDA"} {
		if !strings.Contains(err.Error(), wajib) {
			t.Errorf("galat %q tidak menyebut %q", err, wajib)
		}
	}

	// ⛔ NOL langkah tercatat: pra-terbang berhenti sebelum eksekusi.
	var n int
	if err := sqlDB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM "+skema+".T_MIGRASI").Scan(&n); err != nil {
		t.Fatalf("membaca T_MIGRASI: %v", err)
	}
	if n != 0 {
		t.Errorf("%d langkah tercatat di T_MIGRASI padahal pra-terbang menolak", n)
	}
}
