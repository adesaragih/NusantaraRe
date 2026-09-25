// Package skemauji membangun SKEMA UJI PROVISIONAL untuk Claim Life.
//
// ⛔ BUKAN DDL produksi dan BUKAN migrasi. Tiket 01 menyatakannya terang:
// tipe dan presisi kolom sesungguhnya belum diketahui - tidak ada DDL di korpus
// dan OQ-001 (pemilik DBA) masih terbuka - sehingga di sini dipakai tipe
// numerik berpresisi longgar dan tipe teks berpanjang longgar. Itu sah untuk
// skema uji, tidak sah untuk produksi. DDL sebenarnya milik tiket 14.
//
// Paket ini hanya dipakai oleh test bertag `db`.
package skemauji

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/sijms/go-ora/v2"

	"nusantarare/internal/config"
)

// ErrTanpaOracle dikembalikan bila ORACLE_DSN tidak dikonfigurasi.
var ErrTanpaOracle = errors.New("skemauji: ORACLE_DSN belum dikonfigurasi")

// ErrProduksi menolak pemasangan skema uji pada lingkungan produksi.
var ErrProduksi = errors.New("skemauji: menolak berjalan saat IS_PEGA_PROD=true")

// Buka menyambung ke instance uji memakai env var yang sama dengan aplikasi.
//
// ⛔ Menolak berjalan bila IS_PEGA_PROD bernilai benar (ADR-U-0005). Skema uji
// membuat dan membuang tabel; ia tidak pernah boleh menyentuh produksi.
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

// Tabel yang dibuat, urut dari induk ke anak.
var tabel = []string{"T_GENERAL_CLAIM", "T_CLAIMLF_PREMIUMLIST_DETAIL", "T_CLAIMLF_ADJUSTMENT"}

var ddl = map[string]string{
	"T_GENERAL_CLAIM": `CREATE TABLE %[1]s.T_GENERAL_CLAIM (
		ID            VARCHAR2(64) NOT NULL,
		CLAIM_NO      VARCHAR2(64),
		POLICY_NO     VARCHAR2(64),
		BUSINESS_NAME VARCHAR2(255),
		STS_REJECT    VARCHAR2(8),
		CONSTRAINT PK_UJI_GENERAL_CLAIM PRIMARY KEY (ID))`,

	"T_CLAIMLF_PREMIUMLIST_DETAIL": `CREATE TABLE %[1]s.T_CLAIMLF_PREMIUMLIST_DETAIL (
		ID             VARCHAR2(64) NOT NULL,
		CLAIM_ID       VARCHAR2(64),
		PL_NUMBER      VARCHAR2(64),
		POLICY_NO      VARCHAR2(64),
		CERTIFICATE_NO VARCHAR2(64),
		CURRENCY       VARCHAR2(8),
		CONSTRAINT PK_UJI_PLD PRIMARY KEY (ID))`,

	// CLAIM_AMOUNT memakai NUMBER tanpa presisi: presisi sesungguhnya belum
	// diketahui (OQ-001), dan menebaknya akan menetapkan yang belum diputuskan.
	"T_CLAIMLF_ADJUSTMENT": `CREATE TABLE %[1]s.T_CLAIMLF_ADJUSTMENT (
		ID                     VARCHAR2(64) NOT NULL,
		PREMIUM_LIST_DETAIL_ID VARCHAR2(64),
		CLAIM_AMOUNT           NUMBER,
		CURRENCY               VARCHAR2(8),
		STS_REJECT             VARCHAR2(8),
		ACCEPTED_NO            VARCHAR2(64),
		ACCEPTATION_DATE       DATE,
		KOMITE_ID              VARCHAR2(64),
		CONSTRAINT PK_UJI_ADJ PRIMARY KEY (ID))`,
}

// Relasi dan index dari STRUKTUR-TABEL-CLAIM-LIFE.md.
//
// Lisensi "skema uji provisional" tiket 01 melonggarkan **tipe dan presisi**
// kolom - bukan relasinya. Kaskade dan keunikan di bawah ditulis apa adanya
// sebagaimana dokumen struktur menetapkannya.
var relasi = []string{
	// "anaknya T_CLAIMLF_PREMIUMLIST_DETAIL lewat CLAIM_ID - 1:N - ON DELETE CASCADE"
	`ALTER TABLE %[1]s.T_CLAIMLF_PREMIUMLIST_DETAIL ADD CONSTRAINT FK_UJI_PLD_CLAIM
		FOREIGN KEY (CLAIM_ID) REFERENCES %[1]s.T_GENERAL_CLAIM (ID) ON DELETE CASCADE`,
	// "induknya T_CLAIMLF_PREMIUMLIST_DETAIL lewat PREMIUM_LIST_DETAIL_ID - 1:N - ON DELETE CASCADE"
	`ALTER TABLE %[1]s.T_CLAIMLF_ADJUSTMENT ADD CONSTRAINT FK_UJI_ADJ_PLD
		FOREIGN KEY (PREMIUM_LIST_DETAIL_ID) REFERENCES %[1]s.T_CLAIMLF_PREMIUMLIST_DETAIL (ID) ON DELETE CASCADE`,
	// "Index: PREMIUM_LIST_DETAIL_ID - KOMITE_ID (UNIK)"
	`CREATE INDEX %[1]s.IX_UJI_ADJ_PLD ON %[1]s.T_CLAIMLF_ADJUSTMENT (PREMIUM_LIST_DETAIL_ID)`,
	// Satu baris AdjustmentList = TEPAT satu kasus komite. Oracle mengizinkan
	// banyak NULL di index unik, jadi baris yang belum pernah dikirim ke Komite
	// tetap boleh banyak - persis yang diminta dokumen struktur.
	`CREATE UNIQUE INDEX %[1]s.UX_UJI_ADJ_KOMITE ON %[1]s.T_CLAIMLF_ADJUSTMENT (KOMITE_ID)`,
}

// samakanNLS memaksa titik sebagai pemisah desimal untuk sesi ini.
//
// Jalur BACA dijaga di dalam SQL lewat argumen NLS pada TO_CHAR. Jalur TULIS
// fixture mem-bind teks desimal ke kolom NUMBER, dan konversi implisit itu
// mengikuti NLS sesi: pada sesi ber-NLS koma, "0.00000001" gagal dengan
// ORA-01722. Karena itu sesinya disamakan di sini.
func samakanNLS(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `ALTER SESSION SET NLS_NUMERIC_CHARACTERS = '.,'`)
	if err != nil {
		return fmt.Errorf("skemauji: menyamakan NLS: %w", err)
	}
	return nil
}

func jalankan(ctx context.Context, db *sql.DB, q string) error {
	if _, err := db.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("skemauji: %q: %w", ringkas(q), err)
	}
	return nil
}

func ringkas(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > 70 {
		return s[:70] + "..."
	}
	return s
}

// Bongkar membuang tabel uji. Galat "tabel tidak ada" diabaikan.
func Bongkar(ctx context.Context, db *sql.DB, skema string) error {
	for i := len(tabel) - 1; i >= 0; i-- {
		q := fmt.Sprintf(`DROP TABLE %s.%s CASCADE CONSTRAINTS`, skema, tabel[i])
		if _, err := db.ExecContext(ctx, q); err != nil {
			if strings.Contains(err.Error(), "ORA-00942") { // tabel atau view tidak ada
				continue
			}
			return fmt.Errorf("skemauji: membongkar %s: %w", tabel[i], err)
		}
	}
	return nil
}

// Pasang membuat ulang tabel uji dari keadaan bersih, beserta relasinya.
func Pasang(ctx context.Context, db *sql.DB, skema string) error {
	if err := samakanNLS(ctx, db); err != nil {
		return err
	}
	if err := Bongkar(ctx, db, skema); err != nil {
		return err
	}
	for _, t := range tabel {
		if err := jalankan(ctx, db, fmt.Sprintf(ddl[t], skema)); err != nil {
			return err
		}
	}
	for _, r := range relasi {
		if err := jalankan(ctx, db, fmt.Sprintf(r, skema)); err != nil {
			return err
		}
	}
	return nil
}

// IsiContoh mengisi fixture.
//
// ⛔ Nol nama orang, nol nomor polis nyata, nol potongan dokumen produksi.
// Seluruh nilai jelas buatan dan berawalan "UJI".
//
// Fixture ini sengaja memuat kasus yang harus GAGAL bila kode menebak:
//   - baris berkode "9" - di luar ketiga nilai yang tertulis spec
//   - baris ber-STS_REJECT NULL - kolom nullable (ADR-U-0027)
//   - nomor sertifikat "006" - berawalan nol, harus tetap teks (ADR-U-0022)
//   - jumlah berdesimal delapan angka - harus utuh, tidak boleh lewat float
//   - dua peserta, masing-masing lebih dari satu baris - membuktikan pembacaan
//     tidak berhenti di baris terakhir (tiket 01 AC-1)
func IsiContoh(ctx context.Context, db *sql.DB, skema string) error {
	perintah := []struct {
		q    string
		args []any
	}{
		{fmt.Sprintf(`INSERT INTO %s.T_GENERAL_CLAIM
			(ID, CLAIM_NO, POLICY_NO, BUSINESS_NAME, STS_REJECT)
			VALUES (:1, :2, :3, :4, :5)`, skema),
			[]any{"UJI-KLAIM-1", "UJI-CLM-0001", "UJI-POL-0001", "UJI BISNIS", "0"}},

		{fmt.Sprintf(`INSERT INTO %s.T_CLAIMLF_PREMIUMLIST_DETAIL
			(ID, CLAIM_ID, PL_NUMBER, POLICY_NO, CERTIFICATE_NO, CURRENCY)
			VALUES (:1, :2, :3, :4, :5, :6)`, skema),
			[]any{"UJI-PESERTA-1", "UJI-KLAIM-1", "UJI-PL-1", "UJI-POL-0001", "006", "IDR"}},
		{fmt.Sprintf(`INSERT INTO %s.T_CLAIMLF_PREMIUMLIST_DETAIL
			(ID, CLAIM_ID, PL_NUMBER, POLICY_NO, CERTIFICATE_NO, CURRENCY)
			VALUES (:1, :2, :3, :4, :5, :6)`, skema),
			[]any{"UJI-PESERTA-2", "UJI-KLAIM-1", "UJI-PL-2", "UJI-POL-0001", "010", "IDR"}},

		{fmt.Sprintf(`INSERT INTO %s.T_CLAIMLF_ADJUSTMENT
			(ID, PREMIUM_LIST_DETAIL_ID, CLAIM_AMOUNT, CURRENCY, STS_REJECT, ACCEPTED_NO, KOMITE_ID)
			VALUES (:1, :2, :3, :4, :5, :6, :7)`, skema),
			[]any{"UJI-ADJ-1", "UJI-PESERTA-1", "1234567890.12345678", "IDR", "0", nil, nil}},
		{fmt.Sprintf(`INSERT INTO %s.T_CLAIMLF_ADJUSTMENT
			(ID, PREMIUM_LIST_DETAIL_ID, CLAIM_AMOUNT, CURRENCY, STS_REJECT, ACCEPTED_NO, KOMITE_ID)
			VALUES (:1, :2, :3, :4, :5, :6, :7)`, skema),
			[]any{"UJI-ADJ-2", "UJI-PESERTA-1", "0.00000001", "IDR", "1", "UJI-AKSEP-1", "UJI-KOMITE-1"}},
		{fmt.Sprintf(`INSERT INTO %s.T_CLAIMLF_ADJUSTMENT
			(ID, PREMIUM_LIST_DETAIL_ID, CLAIM_AMOUNT, CURRENCY, STS_REJECT, ACCEPTED_NO, KOMITE_ID)
			VALUES (:1, :2, :3, :4, :5, :6, :7)`, skema),
			[]any{"UJI-ADJ-3", "UJI-PESERTA-1", "250000", "IDR", "2", nil, nil}},
		{fmt.Sprintf(`INSERT INTO %s.T_CLAIMLF_ADJUSTMENT
			(ID, PREMIUM_LIST_DETAIL_ID, CLAIM_AMOUNT, CURRENCY, STS_REJECT, ACCEPTED_NO, KOMITE_ID)
			VALUES (:1, :2, :3, :4, :5, :6, :7)`, skema),
			[]any{"UJI-ADJ-4", "UJI-PESERTA-2", "7500.5", "IDR", "9", nil, nil}},
		{fmt.Sprintf(`INSERT INTO %s.T_CLAIMLF_ADJUSTMENT
			(ID, PREMIUM_LIST_DETAIL_ID, CLAIM_AMOUNT, CURRENCY, STS_REJECT, ACCEPTED_NO, KOMITE_ID)
			VALUES (:1, :2, :3, :4, :5, :6, :7)`, skema),
			[]any{"UJI-ADJ-5", "UJI-PESERTA-2", nil, nil, nil, nil, nil}},
	}
	if err := samakanNLS(ctx, db); err != nil {
		return err
	}
	for _, p := range perintah {
		if _, err := db.ExecContext(ctx, p.q, p.args...); err != nil {
			return fmt.Errorf("skemauji: mengisi contoh: %w", err)
		}
	}
	return nil
}
