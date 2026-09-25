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
	"strings"

	_ "github.com/sijms/go-ora/v2"

	"nusantarare/internal/config"
	"nusantarare/internal/repository"
)

// ErrTanpaOracle dikembalikan bila ORACLE_DSN tidak dikonfigurasi.
var ErrTanpaOracle = errors.New("skemauji: ORACLE_DSN belum dikonfigurasi")

// ErrProduksi menolak pemasangan skema uji pada lingkungan produksi.
var ErrProduksi = errors.New("skemauji: menolak berjalan saat IS_PEGA_PROD=true")

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
	if cfg.IsPegaProd {
		return nil, "", ErrProduksi
	}
	if !cfg.PunyaOracle() {
		return nil, "", ErrTanpaOracle
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
// Kelima puluh lima kolomnya `[terverifikasi]` dari rule
// ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL
// bertipe Rule-Connect-SQL. Tipe kolomnya TIDAK diketahui - korpus tidak
// memuat DDL - sehingga di sini seluruhnya teks longgar, cukup untuk menguji
// pembongkaran. Itu sah untuk skema uji dan tidak sah untuk produksi.
const ddlTabelLama = `CREATE TABLE %[1]s.` + namaTabelLama + ` (
	CASEID VARCHAR2(64), NO_CLAIM VARCHAR2(64), POLICY_NO VARCHAR2(64),
	POLICY_HOLDER VARCHAR2(255), CERTIFICATE_NO VARCHAR2(64),
	NAME_OF_INSURED VARCHAR2(255), SEX VARCHAR2(8), DOB VARCHAR2(32),
	AGE VARCHAR2(8), PLAN VARCHAR2(128),
	BEGIN_DATE VARCHAR2(32), LAPSE_DATE VARCHAR2(32), EXPIRED_DATE VARCHAR2(32),
	STATUS VARCHAR2(32), EM_PERCENT VARCHAR2(64),
	CURRENCY VARCHAR2(8), SUM_INSURED VARCHAR2(64), CEDING_RETENTION VARCHAR2(64),
	SUM_REASURED VARCHAR2(64), SHARE_NUSANTARA_RE VARCHAR2(64),
	CLAIM_AMOUNT VARCHAR2(64), WPC VARCHAR2(32), PL_NUMBER VARCHAR2(64),
	DISEASE VARCHAR2(255), ICD_CODE VARCHAR2(32),
	NOTES VARCHAR2(1000), CEDINGCO VARCHAR2(64), CEDINGCONAME VARCHAR2(255),
	SOB VARCHAR2(64), SOBNAME VARCHAR2(255),
	BUSINESSID VARCHAR2(64), BUSINESSNAME VARCHAR2(255), SHARE_RETRO VARCHAR2(64),
	KETERANGAN VARCHAR2(1000), ACCEPTATION_DATE VARCHAR2(32),
	NO_ACCEPTATION VARCHAR2(64), STS_REJECT VARCHAR2(8), CLAIM_RETRO VARCHAR2(64),
	RETROID VARCHAR2(64), RETRONAME VARCHAR2(255),
	SECURITYREINSURERID VARCHAR2(64), SECURITYREINSURER VARCHAR2(255),
	TYPECEDING VARCHAR2(64), TYPE VARCHAR2(32), CONFIRMATION_DATE VARCHAR2(32),
	CLAIM_RECEIVED_DATE VARCHAR2(32), COMPLETE_DATE VARCHAR2(32), ID VARCHAR2(64),
	NAME_OF_BANK VARCHAR2(255), IDBANK VARCHAR2(64),
	ACCOUNTNO VARCHAR2(64), CREATEOPNAME VARCHAR2(128), PRODUCTNAMEID VARCHAR2(64),
	PRODUCTNAME VARCHAR2(255), RETROCEDED_SHARE VARCHAR2(64))`

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

	if _, err := db.ExecContext(ctx, fmt.Sprintf(ddlTabelLama, skema)); err != nil {
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
	q := fmt.Sprintf(`INSERT INTO %s.%s
		(ID, CASEID, NO_CLAIM, POLICY_NO, CERTIFICATE_NO, PL_NUMBER, BUSINESSNAME,
		 CLAIM_RETRO, CURRENCY, CLAIM_AMOUNT, STS_REJECT, NO_ACCEPTATION,
		 ACCEPTATION_DATE, TYPE, CREATEOPNAME)
		VALUES (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15)`, skema, namaTabelLama)
	for _, b := range baris {
		_, err := db.ExecContext(ctx, q,
			b.ID, b.CASEID, b.NO_CLAIM, b.POLICY_NO, b.CERTIFICATE_NO, b.PL_NUMBER,
			b.BUSINESSNAME, b.CLAIM_RETRO, b.CURRENCY, b.CLAIM_AMOUNT, b.STS_REJECT,
			b.NO_ACCEPTATION, b.ACCEPTATION_DATE, b.TYPE, b.CREATEOPNAME)
		if err != nil {
			return fmt.Errorf("skemauji: mengisi baris lama %s: %w", b.ID, err)
		}
	}
	return nil
}
