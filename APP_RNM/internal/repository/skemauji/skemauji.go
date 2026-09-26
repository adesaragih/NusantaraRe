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
	"os"
	"strings"

	_ "github.com/sijms/go-ora/v2"

	"nusantarare/internal/config"
	"nusantarare/internal/repository"
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
func BukaRepositori() (*repository.DB, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	if err := pastikanAman(cfg); err != nil {
		return nil, err
	}
	return repository.Open(cfg)
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
	if _, err := repo.JalankanMigrasi(ctx); err != nil {
		return fmt.Errorf("skemauji: menjalankan migrasi: %w", err)
	}

	if _, err := db.ExecContext(ctx, ddlTabelLama(skema)); err != nil {
		if !strings.Contains(err.Error(), "ORA-00955") { // nama sudah dipakai
			return fmt.Errorf("skemauji: membuat tiruan tabel warisan: %w", err)
		}
	}
	return nil
}

// Bongkar membuang seluruh objek skema uji lewat jalur mundur migrasi,
// ditambah tabel tiruan warisan.
func Bongkar(ctx context.Context, db *sql.DB, skema string) error {
	q := fmt.Sprintf(`DROP TABLE %s.%s CASCADE CONSTRAINTS`, skema, namaTabelLama)
	if _, err := db.ExecContext(ctx, q); err != nil {
		if !strings.Contains(err.Error(), "ORA-00942") { // tabel tidak ada
			return fmt.Errorf("skemauji: membongkar tiruan warisan: %w", err)
		}
	}

	repo, err := BukaRepositori()
	if err != nil {
		return err
	}
	defer func() { _ = repo.Close() }()
	if _, err := repo.BongkarMigrasi(ctx); err != nil {
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
