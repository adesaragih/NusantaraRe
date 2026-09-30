// Package skemauji menyiapkan skema uji Oracle untuk test bertag `db`.
//
// Untuk apa paket ini: membuat tabel, mengisinya dengan fixture buatan, lalu
// membongkarnya lagi sesudah test selesai.
//
// Dibaca sesudah: repository/migrasi.go.
//
// ⛔ BUKAN DDL produksi. Sejak tiket 14, tabel BARU tidak lagi dibuat di sini:
// Pasang memanggil pelari migrasi yang sama dengan yang dipakai aplikasi, dan
// Bongkar memanggil jalur mundurnya. Dengan begitu yang diuji adalah migrasi
// yang sebenarnya, bukan DDL tangan yang mirip.
//
// Yang masih dibuat di sini hanya satu: tabel TIRUAN OS_AKSEPTASI_KLAIM_LIFE,
// yaitu bentuk datar warisan yang dipakai untuk menguji pembongkaran data lama.
// Menjalankan migrasi terhadap tabel warisan SUNGGUHAN adalah tiket 13 dan
// tunduk pada persetujuan manusia.
package skemauji

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	_ "github.com/sijms/go-ora/v2"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/migrasi"
	"nusantarare/modul"
	"nusantarare/modul/claimlife/repository"
	"nusantarare/modul/premiumlistlife/services"
)

// ErrTanpaOracle dikembalikan bila ORACLE_DSN tidak dikonfigurasi.
var ErrTanpaOracle = errors.New("skemauji: ORACLE_DSN belum dikonfigurasi")

// ErrProduksi menolak pemasangan skema uji pada lingkungan produksi.
var ErrProduksi = errors.New("skemauji: menolak berjalan saat IS_PEGA_PROD=true")

// ErrBukanSkemaUji menolak berjalan di skema yang belum dinyatakan skema uji.
//
// ⛔ Kenapa galat, bukan SKIP: melewati diam-diam berarti pagarnya tidak
// pernah terlihat oleh orang yang salah menyetel env. Yang dipertaruhkan
// bukan test yang gagal, melainkan tabel warisan yang terhapus.
var ErrBukanSkemaUji = config.ErrBukanSkemaUji

// envSkemaUji adalah pengakuan sadar dari orang yang menjalankan test bahwa
// skema yang ditunjuk ORACLE_SCHEMA memang boleh dihapus isinya.
const envSkemaUji = config.EnvSkemaUji

// namaSkemaWarisan adalah skema Pega yang memuat tabel warisan sungguhan.
// Ia tidak pernah boleh menjadi skema uji, sekalipun di instance pengembangan.
const namaSkemaWarisan = config.NamaSkemaWarisan

// BolehDilewati menyatakan apakah sebuah galat dari Buka layak dijawab t.Skip.
//
// Hanya SATU yang layak: Oracle memang belum dikonfigurasi. Selebihnya -
// salah konfigurasi, menunjuk produksi, atau menunjuk skema yang bukan skema
// uji - harus menggagalkan test. Sebelum ronde 4 seluruhnya dilewati, sehingga
// IS_PEGA_PROD=true pun menghasilkan lari hijau yang tidak menguji apa pun.
func BolehDilewati(err error) bool { return errors.Is(err, ErrTanpaOracle) }

// SumberMigrasi mengembalikan folder migrasi SETIAP modul terdaftar - daftar
// yang sama dengan `-migrate`.
//
// Refactor bentuk B paket 8: paket ini SENGAJA mengenal semua modul - skema
// uji adalah skema UTUH - dan karena itu satu-satunya jalur lintas modul yang
// disahkan bagi berkas uji sebuah modul (`impor_lintas_modul_test.go`). Uji db
// modul memanggil ini, bukan daftar `nusantarare/modul` langsung.
func SumberMigrasi() []fs.FS { return modul.SumberMigrasi() }

// pastikanAman menjalankan seluruh pemeriksaan pintu masuk, dalam satu urutan.
//
// Urutannya disengaja: produksi ditolak lebih dulu, lalu "Oracle memang belum
// ada" - satu-satunya yang layak dilewati - lalu pagar skema uji. Kedua pintu
// masuk paket ini memanggil fungsi yang sama, supaya tidak mungkin ada pintu
// yang tertinggal saat pemeriksaannya bertambah.
func pastikanAman(cfg config.Config) error {
	if cfg.IsPegaProd {
		return ErrProduksi
	}
	if !cfg.PunyaOracle() {
		return ErrTanpaOracle
	}
	return periksaPagarSkemaUji(cfg.OracleSchema)
}

// periksaPagarSkemaUji membaca pengakuan itu dari environment.
func periksaPagarSkemaUji(skema string) error {
	return pagarSkemaUji(skema, os.Getenv(envSkemaUji))
}

// pagarSkemaUji adalah isi keputusannya, dipisah dari environment supaya dapat
// diuji tanpa Oracle dan tanpa env var.
//
// Dua syarat, KEDUANYA harus benar:
//  1. env ORACLE_SKEMA_UJI bernilai "true"
//  2. skema yang ditunjuk bukan POOLDATA
//
// Syarat kedua tidak dapat ditutupi oleh syarat pertama: menyetel env tidak
// membuat skema warisan menjadi skema uji.
func pagarSkemaUji(skema, diakui string) error {
	return config.PagarSkemaUji(skema, diakui)
}

// namaTabelLama adalah tabel datar warisan yang ditiru untuk uji migrasi.
const namaTabelLama = "OS_AKSEPTASI_KLAIM_LIFE"

// Buka menyambung ke instance uji memakai env var yang sama dengan aplikasi.
//
// ⛔ Menolak berjalan bila IS_PEGA_PROD bernilai benar (ADR-U-0005).
func Buka() (*sql.DB, string, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, "", err
	}
	if err := pastikanAman(cfg); err != nil {
		return nil, "", err
	}
	db, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		return nil, "", err
	}
	return db, cfg.OracleSchema, nil
}

// BukaRepositori membuka koneksi lapisan repository dari env var yang sama.
func BukaRepositori() (*db.DB, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	if err := pastikanAman(cfg); err != nil {
		return nil, err
	}
	return db.Open(cfg)
}

// samakanNLS memaksa titik sebagai pemisah desimal untuk sesi ini.
//
// Jalur BACA dijaga di dalam SQL lewat argumen NLS pada TO_CHAR. Jalur TULIS
// fixture mem-bind teks desimal ke kolom NUMBER, dan konversi implisit itu
// mengikuti NLS sesi: pada sesi ber-NLS koma, "0.00000001" gagal ORA-01722.
func samakanNLS(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER SESSION SET NLS_NUMERIC_CHARACTERS = '.,'`); err != nil {
		return fmt.Errorf("skemauji: menyamakan NLS: %w", err)
	}
	return nil
}

// ddlTabelLama membuat tiruan tabel datar warisan.
//
// Kelima puluh lima nama kolomnya `[terverifikasi]` dari rule
// ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL
// bertipe Rule-Connect-SQL. ✅ [data DBA] TIPE kolomnya SUDAH DIKETAHUI sejak
// 26-09-2026: dibaca dari ALL_TAB_COLUMNS instance pengembangan dan disimpan di
// .scratch/claim-life/TIPE-KOLOM-OS-AKSEPTASI-KLAIM-LIFE.md. Tiruan ini tidak
// lagi menebak, dan tabelnya 62 kolom - bukan 55. ⚠️ Produksi belum dibaca;
// DBA yang dapat memastikan bentuknya sama.
//
// Daftarnya TIDAK ditulis ulang di sini. Ia datang dari repository, tempat
// pembaca dan penulis mengambil daftar yang sama, sehingga ketiganya tidak
// mungkin berselisih nama maupun urutan.
//
// ⭐ 26-09-2026: tiruan memakai KEENAM PULUH DUA kolom katalog, bukan lagi
// kelima puluh lima yang ditulis rule. Tujuh sisanya tidak pernah ditulis
// maupun dibaca kode ini, tetapi ada di tabel sungguhan - dan tiruan yang
// kekurangan kolom tidak membuktikan bahwa INSERT 55 kolom berhasil pada
// tabel 62 kolom.
func ddlTabelLama(skema string) string {
	var kolom []string
	for _, n := range repository.NamaKolomTabelWarisan() {
		kolom = append(kolom, n+" "+repository.TipeKolomBarisLama(n))
	}
	return fmt.Sprintf("CREATE TABLE %s.%s (%s)", skema, namaTabelLama,
		strings.Join(kolom, ", "))
}

// Pasang membangun skema uji dari keadaan bersih.
//
// Urutannya: bongkar dulu, lalu jalankan migrasi yang sebenarnya, lalu buat
// tabel tiruan warisan.
func Pasang(ctx context.Context, db *sql.DB, skema string) error {
	if err := samakanNLS(ctx, db); err != nil {
		return err
	}
	if err := Bongkar(ctx, db, skema); err != nil {
		return err
	}

	repo, err := BukaRepositori()
	if err != nil {
		return err
	}
	defer func() { _ = repo.Close() }()
	if _, err := migrasi.Jalankan(ctx, repo, modul.SumberMigrasi()...); err != nil {
		return fmt.Errorf("skemauji: menjalankan migrasi: %w", err)
	}

	if _, err := db.ExecContext(ctx, ddlTabelLama(skema)); err != nil {
		if !strings.Contains(err.Error(), "ORA-00955") { // nama sudah dipakai
			return fmt.Errorf("skemauji: membuat tiruan tabel warisan: %w", err)
		}
	}

	// Tiruan tabel peserta polis. Tanpa ini, pendaftaran klaim tidak dapat
	// membaca peserta yang dipilih, dan test db tiket 02 gagal di pembacaan -
	// bukan menguji pendaftarannya.
	if _, err := db.ExecContext(ctx, ddlTiruanPesertaPolis(skema)); err != nil {
		if !strings.Contains(err.Error(), "ORA-00955") {
			return fmt.Errorf("skemauji: membuat tiruan peserta polis: %w", err)
		}
	}

	// Tiruan rekap warisan PremiumList + sequence-nya - PL-09 (GILIRAN-18).
	for _, q := range ddlTiruanSummaryPolis(skema) {
		if _, err := db.ExecContext(ctx, q); err != nil {
			if !strings.Contains(err.Error(), "ORA-00955") {
				return fmt.Errorf("skemauji: membuat tiruan summary polis: %w", err)
			}
		}
	}

	// Tiruan tabel treaty untuk perhitungan spreading (tiket 03).
	for _, q := range ddlTiruanTreaty(skema) {
		if _, err := db.ExecContext(ctx, q); err != nil {
			if !strings.Contains(err.Error(), "ORA-00955") {
				return fmt.Errorf("skemauji: membuat tiruan treaty: %w", err)
			}
		}
	}

	// Tiruan enam tabel warisan Treaty Contract Out + sequence warisannya -
	// tco4: tabel yang modul itu tulis dan baca langsung. ADITIF 28/29-09-2026.
	for _, q := range ddlTiruanTCO(skema) {
		if _, err := db.ExecContext(ctx, q); err != nil {
			if !strings.Contains(err.Error(), "ORA-00955") {
				return fmt.Errorf("skemauji: membuat tiruan warisan Treaty Contract Out: %w", err)
			}
		}
	}
	return nil
}

// Bongkar membuang seluruh objek skema uji lewat jalur mundur migrasi,
// ditambah tabel tiruan warisan.
func Bongkar(ctx context.Context, db *sql.DB, skema string) error {
	tiruan := []string{namaTabelLama, namaTabelPesertaPolis,
		namaTabelRetrosesi, namaTabelTahunTreaty, namaTabelSummaryPolis}
	// Enam tiruan warisan Treaty Contract Out ikut dibongkar (aditif 28-09-2026).
	tiruan = append(tiruan, namaTabelTiruanTCO...)
	for _, nama := range tiruan {
		q := fmt.Sprintf(`DROP TABLE %s.%s CASCADE CONSTRAINTS`, skema, nama)
		if _, err := db.ExecContext(ctx, q); err != nil {
			if !strings.Contains(err.Error(), "ORA-00942") { // tabel tidak ada
				return fmt.Errorf("skemauji: membongkar tiruan %s: %w", nama, err)
			}
		}
	}
	// Treaty Contract Out tco4 (aditif 29-09-2026) dan PL-09 (GILIRAN-18):
	// sequence warisan tiruan.
	for _, nama := range append([]string{namaSequenceSummaryPolis}, namaSequenceTiruanTCO...) {
		q := fmt.Sprintf(`DROP SEQUENCE %s.%s`, skema, nama)
		if _, err := db.ExecContext(ctx, q); err != nil {
			if !strings.Contains(err.Error(), "ORA-02289") { // sequence tidak ada
				return fmt.Errorf("skemauji: membongkar sequence tiruan %s: %w", nama, err)
			}
		}
	}

	repo, err := BukaRepositori()
	if err != nil {
		return err
	}
	defer func() { _ = repo.Close() }()
	if _, err := migrasi.Bongkar(ctx, repo, modul.SumberMigrasi()...); err != nil {
		return fmt.Errorf("skemauji: membongkar migrasi: %w", err)
	}

	// Tabel pencatat migrasi ikut dibuang supaya skema uji benar-benar bersih.
	if _, err := db.ExecContext(ctx,
		fmt.Sprintf(`DROP TABLE %s.T_MIGRASI CASCADE CONSTRAINTS`, skema)); err != nil {
		if !strings.Contains(err.Error(), "ORA-00942") {
			return fmt.Errorf("skemauji: membongkar catatan migrasi: %w", err)
		}
	}
	return nil
}

// IsiContoh mengisi fixture untuk tiket 01 - satu klaim, dua peserta, lima
// baris adjustment.
//
// ⛔ Nol nama orang, nol nomor polis nyata, nol potongan data produksi.
//
// Fixture ini sengaja memuat kasus yang harus GAGAL bila kode menebak:
//   - baris berkode "9" - di luar ketiga nilai yang tertulis spec
//   - baris ber-STS_REJECT NULL - kolom nullable (ADR-U-0027)
//   - nomor sertifikat "006" - berawalan nol, harus tetap teks (ADR-U-0022)
//   - jumlah berdesimal delapan angka - harus utuh, tidak lewat float
//   - dua peserta, masing-masing lebih dari satu baris
//   - satu baris komite yang sungguh ada, ditunjuk KOMITE_ID lewat FK
func IsiContoh(ctx context.Context, db *sql.DB, skema string) error {
	if err := samakanNLS(ctx, db); err != nil {
		return err
	}
	perintah := []struct {
		q    string
		args []any
	}{
		{fmt.Sprintf(`INSERT INTO %s.T_WORK_CLAIM (ID, LINI, TYPE, CASE_ID)
			VALUES (:1, :2, :3, :4)`, skema),
			[]any{"CLM-UJI001", "LIFE", "UJI-TYPE", "UJI-CASE-1"}},

		// Baris komite yang ditunjuk KOMITE_ID di bawah. Sejak keputusan work
		// owner d (26-09-2026) KOMITE_ID ber-REFERENCES ke T_WORK_CLAIM(ID),
		// jadi penunjuk yatim tidak lagi diterima Oracle - fixture harus
		// menyediakan barisnya, persis seperti data sungguhan harus.
		//
		// COVER_KEY-nya sengaja dibiarkan NULL. Mengisinya dengan CLM-UJI001
		// akan membuat penghapusan klaim (AC 38) ditolak ORA-02292, sebab FK
		// COVER_KEY tanpa ON DELETE menolak menghapus induk yang masih
		// ditunjuk. Itu pertanyaan jalur hapus tiket 15, bukan tiket ini.
		{fmt.Sprintf(`INSERT INTO %s.T_WORK_CLAIM (ID, LINI, TYPE, CASE_ID)
			VALUES (:1, :2, :3, :4)`, skema),
			[]any{"KMT-UJI001", "LIFE", "UJI-KOMITE", "UJI-CASE-KMT"}},

		{fmt.Sprintf(`INSERT INTO %s.T_GENERAL_CLAIM
			(ID, CLAIM_NO, POLICY_NO, BUSINESS_NAME, STS_REJECT)
			VALUES (:1, :2, :3, :4, :5)`, skema),
			[]any{"CLM-UJI001", "UJI-CLM-0001", "UJI-POL-0001", "UJI BISNIS", "0"}},

		{fmt.Sprintf(`INSERT INTO %s.T_CLAIMLF_PREMIUMLIST_DETAIL
			(ID, CLAIM_ID, PL_NUMBER, POLICY_NO, CERTIFICATE_NO, CURRENCY)
			VALUES (:1, :2, :3, :4, :5, :6)`, skema),
			[]any{"UJI-PESERTA-1", "CLM-UJI001", "UJI-PL-1", "UJI-POL-0001", "006", "IDR"}},
		{fmt.Sprintf(`INSERT INTO %s.T_CLAIMLF_PREMIUMLIST_DETAIL
			(ID, CLAIM_ID, PL_NUMBER, POLICY_NO, CERTIFICATE_NO, CURRENCY)
			VALUES (:1, :2, :3, :4, :5, :6)`, skema),
			[]any{"UJI-PESERTA-2", "CLM-UJI001", "UJI-PL-2", "UJI-POL-0001", "010", "IDR"}},

		{adjSQL(skema), []any{"UJI-ADJ-1", "UJI-PESERTA-1", "1234567890.12345678", "IDR", "0", nil, nil}},
		{adjSQL(skema), []any{"UJI-ADJ-2", "UJI-PESERTA-1", "0.00000001", "IDR", "1", "UJI-AKSEP-1", "KMT-UJI001"}},
		{adjSQL(skema), []any{"UJI-ADJ-3", "UJI-PESERTA-1", "250000", "IDR", "2", nil, nil}},
		{adjSQL(skema), []any{"UJI-ADJ-4", "UJI-PESERTA-2", "7500.5", "IDR", "9", nil, nil}},
		{adjSQL(skema), []any{"UJI-ADJ-5", "UJI-PESERTA-2", nil, nil, nil, nil, nil}},
	}
	for _, p := range perintah {
		if _, err := db.ExecContext(ctx, p.q, p.args...); err != nil {
			return fmt.Errorf("skemauji: mengisi contoh: %w", err)
		}
	}
	return nil
}

func adjSQL(skema string) string {
	return fmt.Sprintf(`INSERT INTO %s.T_CLAIMLF_ADJUSTMENT
		(ID, PREMIUM_LIST_DETAIL_ID, CLAIM_AMOUNT, CURRENCY, STS_REJECT, ACCEPTED_NO, KOMITE_ID)
		VALUES (:1, :2, :3, :4, :5, :6, :7)`, skema)
}

// IsiBarisLama mengisi tiruan tabel datar warisan dengan baris buatan, untuk
// menguji pembongkaran data lama.
func IsiBarisLama(ctx context.Context, db *sql.DB, skema string, baris []repository.BarisLama) error {
	if err := samakanNLS(ctx, db); err != nil {
		return err
	}
	kolom := repository.NamaKolomBarisLama()
	var penampung []string
	for i, n := range kolom {
		penampung = append(penampung, repository.PenampungTulisLama(n, i+1))
	}
	q := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s)", skema, namaTabelLama,
		strings.Join(kolom, ", "), strings.Join(penampung, ", "))
	for _, b := range baris {
		// ⛔ Dipagari sama seperti jalur tulis aplikasi (butir s1). Inilah
		// jalur yang melahirkan ORA-01722 pada ronde 5: fixture mengisi
		// CLAIM_RETRO dengan teks, dan tiruan yang berbentuk benar menolaknya.
		// Fixture yang salah harus berteriak di sini, bukan di Oracle.
		if err := repository.PeriksaNilaiWarisan(b); err != nil {
			return err
		}
		if _, err := db.ExecContext(ctx, q, repository.NilaiBarisLama(b)...); err != nil {
			return fmt.Errorf("skemauji: mengisi baris lama %s: %w", b.ID, err)
		}
	}
	return nil
}

// namaTabelPesertaPolis adalah tabel warisan peserta polis yang ditiru.
const namaTabelPesertaPolis = "M_LIFE_PREMIUM_DETAIL"

// ddlTiruanPesertaPolis membuat tiruan M_LIFE_PREMIUM_DETAIL.
//
// ⛔ Kenapa ini ada: tiket 02 membaca peserta dari tabel itu saat mendaftar,
// dan skema uji tidak pernah membuatnya - sehingga ketiga test db tiket 02
// akan GAGAL di pembacaan peserta begitu Oracle ada, bukan menguji pendaftaran.
// Hari ini semuanya SKIP, jadi cacatnya tidak terlihat.
//
// ⚠️ Yang ditiru hanya kolom yang DIBACA kode ini - dua puluh empat kolom
// kolomSalin ditambah EDMSTATUS. Tabel sungguhannya 85 kolom; menirunya utuh
// tidak menambah satu pun bukti, dan kolom yang tidak pernah dibaca hanya
// menambah tempat untuk salah.
//
// ⛔ Kolom KTP TIDAK ikut ditiru, dan itu disengaja: ia tidak pernah dibaca,
// dan tabel uji yang memuat tempat untuk nomor identitas adalah undangan.
//
// Tipe kolomnya `[data DBA]` dari KATALOG-TABEL-PESERTA-DAN-TREATY.md, bukan
// diturunkan dari nama. ⭐ Dua di antaranya mengejutkan: STNC bertipe DATE
// meski namanya tidak berbunyi begitu, dan seluruh kolom uang NUMBER tanpa
// presisi.
func ddlTiruanPesertaPolis(skema string) string {
	return fmt.Sprintf(`CREATE TABLE %s.%s (
		ID VARCHAR2(50),
		PL_NUMBER VARCHAR2(255),
		POLICY_NO VARCHAR2(255),
		CERTIFICATE_NO VARCHAR2(255),
		CURRENCY VARCHAR2(255),
		STNC DATE,
		GROSS_VALUATION_BEGIN_DATE DATE,
		GROSS_VALUATION_EXPIRED_DATE DATE,
		RETRO_VALUATION_BEGIN_DATE DATE,
		RETRO_VALUATION_EXPIRED_DATE DATE,
		WPC DATE,
		BEGIN_DATE DATE,
		EFFECTIVE_DATE DATE,
		LAPSE_DATE DATE,
		EXPIRED_DATE DATE,
		SUM_INSURED NUMBER,
		SUM_REASURED NUMBER,
		GROSS_PREMIUM NUMBER,
		NET_PREMIUM NUMBER,
		CEDING_RETENTION NUMBER,
		SHARE_NUSANTARA_RE NUMBER,
		SHARE_RETRO NUMBER,
		RETROCEDED_SHARE NUMBER,
		EM_PERCENT NUMBER,
		EDMSTATUS VARCHAR2(25),
		SHARE_NUSANTARA_RE_GROSS NUMBER,
		AGE NUMBER,
		ENTRY_AGE NUMBER,
		CURRENT_AGE NUMBER,
		CLAIM_AMOUNT NUMBER,
		NAME_OF_INSURED VARCHAR2(255),
		DOB DATE
	)`, skema, namaTabelPesertaPolis)
}

// IsiPesertaPolis mengisi tiruan itu dengan dua peserta buatan.
//
// ⛔ Nol nama orang, nol nomor polis nyata, nol potongan data produksi -
// seluruhnya berawalan UJI-.
//
// ⭐ Satu peserta ber-EDMSTATUS NULL (new business, yang HARUS muncul) dan satu
// ber-'Batal' (yang HARUS disaring keluar). Tanpa keduanya, penyaring hidup
// tidak pernah benar-benar diuji terhadap Oracle - hanya terhadap dirinya
// sendiri di test murni.
func IsiPesertaPolis(ctx context.Context, db *sql.DB, skema string) error {
	if err := samakanNLS(ctx, db); err != nil {
		return err
	}
	q := fmt.Sprintf(`INSERT INTO %s.%s
		(ID, PL_NUMBER, POLICY_NO, CERTIFICATE_NO, CURRENCY, STNC,
		 GROSS_VALUATION_BEGIN_DATE, GROSS_VALUATION_EXPIRED_DATE,
		 RETRO_VALUATION_BEGIN_DATE, RETRO_VALUATION_EXPIRED_DATE,
		 WPC, BEGIN_DATE, EFFECTIVE_DATE, LAPSE_DATE, EXPIRED_DATE,
		 SUM_INSURED, SUM_REASURED, GROSS_PREMIUM, NET_PREMIUM,
		 CEDING_RETENTION, SHARE_NUSANTARA_RE, SHARE_RETRO,
		 RETROCEDED_SHARE, EM_PERCENT, EDMSTATUS, CLAIM_AMOUNT,
		 NAME_OF_INSURED, DOB)
		VALUES (:1,:2,:3,:4,:5,
		 TO_DATE(:6,'YYYY-MM-DD'),
		 TO_DATE(:7,'YYYY-MM-DD'), TO_DATE(:8,'YYYY-MM-DD'),
		 TO_DATE(:9,'YYYY-MM-DD'), TO_DATE(:10,'YYYY-MM-DD'),
		 TO_DATE(:11,'YYYY-MM-DD'), TO_DATE(:12,'YYYY-MM-DD'),
		 TO_DATE(:13,'YYYY-MM-DD'), TO_DATE(:14,'YYYY-MM-DD'),
		 TO_DATE(:15,'YYYY-MM-DD'),
		 :16,:17,:18,:19,:20,:21,:22,:23,:24,:25,:26,
		 :27, TO_DATE(:28,'YYYY-MM-DD'))`, skema, namaTabelPesertaPolis)

	baris := [][]any{
		// Peserta hidup: EDMSTATUS NULL, seperti seluruh baris new business.
		{"UJI-SRC-1", "UJI-PL-1", "UJI-POL-0001", "006", "IDR", "2026-01-01",
			"2026-01-01", "2026-12-31", "2026-02-01", "2026-11-30",
			"2026-03-01", "2026-01-01", "2026-01-15", "2027-01-01", "2026-12-31",
			"1000000", "900000", "50000", "45000", "100000", "800000", "200000",
			"150000", "0.1", nil, "25000.123449", "UJI-TERTANGGUNG-1", "1980-01-01"},
		// Peserta batal: HARUS disaring keluar.
		{"UJI-SRC-2", "UJI-PL-1", "UJI-POL-0001", "010", "IDR", "2026-01-01",
			"2026-01-01", "2026-12-31", "2026-02-01", "2026-11-30",
			"2026-03-01", "2026-01-01", "2026-01-15", "2027-01-01", "2026-12-31",
			"2000000", "1800000", "60000", "55000", "200000", "1600000", "400000",
			"300000", "0.2", "Batal", nil, "UJI-TERTANGGUNG-2", "1981-02-02"},
	}
	for _, b := range baris {
		if _, err := db.ExecContext(ctx, q, b...); err != nil {
			return fmt.Errorf("skemauji: mengisi tiruan peserta polis: %w", err)
		}
	}
	return nil
}

// Tabel treaty yang ditiru untuk perhitungan spreading (tiket 03).
const (
	namaTabelRetrosesi   = "RETROCESSIONLIFE"
	namaTabelTahunTreaty = "TREATYYEAR_LIFE"
)

// ddlTiruanTreaty membuat tiruan tabel treaty.
//
// ⭐ RETROCESSIONLIFE di instance pengembangan adalah **VIEW** ber-13 kolom
// yang SELURUHNYA VARCHAR2(4000) - termasuk PERCENTSHARE, RATE, COMMISION, dan
// OVR_COMM. Tiruannya dibuat dengan tipe yang sama persis, bukan dengan NUMBER
// yang "lebih benar": justru jalur "teks -> ParseDecimal -> laporkan yang
// gagal" itulah yang perlu diuji. Tiruan bertipe NUMBER akan membuat Oracle
// mengurai angkanya lebih dulu, dan pembacanya tidak pernah menemui teks.
//
// ⛔ RATE_LIFE TIDAK ditiru. Katalog baru memuat enam kolom pertamanya
// (ID, IDUSEDBY, USEDBY, TYPE, GENDER, CONTRACT) dan tidak satu pun di
// antaranya kolom rate. Menirunya berarti mengarang bentuk, dan membaca rate
// dari tabel yang bentuknya dikarang berarti mengarang angkanya.
// [data DBA] - lihat bab Implementasi tiket 03.
func ddlTiruanTreaty(skema string) []string {
	return []string{
		fmt.Sprintf(`CREATE TABLE %s.%s (
			ID VARCHAR2(4000),
			IDTREATYYEAR_LIFE VARCHAR2(4000),
			PERCENTSHARE VARCHAR2(4000),
			RATE VARCHAR2(4000),
			COMMISION VARCHAR2(4000),
			OVR_COMM VARCHAR2(4000),
			TREATYTYPEID VARCHAR2(4000),
			TREATYTYPENAME VARCHAR2(4000),
			TREATYSTARTDATE VARCHAR2(4000),
			TREATYENDDATE VARCHAR2(4000),
			USERID VARCHAR2(4000),
			TGLUPDATE VARCHAR2(4000),
			REINSURERNAME VARCHAR2(4000)
		)`, skema, namaTabelRetrosesi),
		fmt.Sprintf(`CREATE TABLE %s.%s (
			ID VARCHAR2(100),
			TREATYYEAR VARCHAR2(100),
			UNDERWRITINGYEAR VARCHAR2(100),
			USERID VARCHAR2(100),
			TGLUPDATE DATE,
			STARTDATE DATE,
			ENDDATE DATE
		)`, skema, namaTabelTahunTreaty),
	}
}

// PembacaPolis menyambung pembaca polis PremiumList Life untuk test bertag
// `db` yang menguji Claim Life (butir pl4/av).
//
// Refactor bentuk B (30-09-2026): di produksi `cmd/api` yang menyambungnya
// lewat `inti/backend/kontrak`. Test Claim Life tidak boleh mengimpor modul
// PremiumList; penunjang uji ini - yang memang mengenal semua modul - yang
// menyerahkannya, sehingga yang dibaca test sama dengan yang dibaca produksi.
func PembacaPolis(repo *db.DB) kontrak.PembacaPolis {
	return services.PembacaPolis(services.New(repo))
}
