// Package db adalah pintu Oracle bersama: koneksi, transaksi, dan pembantu
// kueri yang dipakai repository setiap modul.
//
// Refactor bentuk B (30-09-2026): isinya dulu kepala paket `repository`.
// Repository tiap modul tetap SATU-SATUNYA lapisan yang menyentuh Oracle;
// paket ini hanya pintunya. Arah ketergantungan: handlers -> services ->
// repository -> db. Paket ini tidak pernah mengimpor modul mana pun.
//
// Fase 0 - scaffold: nol query, nol aturan dagang. Yang ada di sini adalah
// bentuk yang harus diikuti setiap tiket penyimpanan.
package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	// [usulan] Driver Oracle murni Go; mendaftar dengan nama "oracle".
	_ "github.com/sijms/go-ora/v2"

	"nusantarare/inti/config"
)

// namaDriver adalah nama yang didaftarkan go-ora ke database/sql.
const namaDriver = "oracle"

var (
	// ErrTanpaOracle muncul bila lapisan ini dipakai tanpa ORACLE_DSN.
	ErrTanpaOracle = errors.New("repository: ORACLE_DSN belum dikonfigurasi")
	// ErrTransaksiDiSQL muncul bila teks SQL memuat COMMIT.
	ErrTransaksiDiSQL = errors.New("repository: COMMIT tidak boleh ada di teks SQL (ADR-U-0029)")
	// ErrObjekTakBernama muncul bila nama objek kosong saat dikualifikasi.
	ErrObjekTakBernama = errors.New("repository: nama objek kosong")
)

// DB membungkus koneksi beserta nama skemanya.
//
// Skema disimpan di sini supaya setiap query dapat menyebutnya secara
// eksplisit (ADR-U-0033) tanpa mengandalkan skema bawaan sesi.
type DB struct {
	sql   *sql.DB
	skema string
	// isPegaProd menandai lingkungan yang menunjuk data produksi Pega
	// (ADR-U-0005). Pembacanya adalah pelari migrasi, yang menolak berjalan
	// di sana - lihat migrasi.go.
	isPegaProd bool
}

// Open membuka koneksi. sql.Open tidak menghubungi server; pemeriksaan
// sesungguhnya dilakukan Ping.
func Open(cfg config.Config) (*DB, error) {
	if !cfg.PunyaOracle() {
		return nil, ErrTanpaOracle
	}
	h, err := sql.Open(namaDriver, cfg.OracleDSN)
	if err != nil {
		return nil, fmt.Errorf("repository: membuka koneksi: %w", err)
	}
	return &DB{sql: h, skema: cfg.OracleSchema, isPegaProd: cfg.IsPegaProd}, nil
}

// Ping memeriksa koneksi.
func (d *DB) Ping(ctx context.Context) error {
	if d == nil || d.sql == nil {
		return ErrTanpaOracle
	}
	return d.sql.PingContext(ctx)
}

// Close menutup koneksi.
func (d *DB) Close() error {
	if d == nil || d.sql == nil {
		return nil
	}
	return d.sql.Close()
}

// Skema mengembalikan nama skema yang dipakai.
func (d *DB) Skema() string { return d.skema }

// PegaProduksi menyatakan koneksi ini menunjuk data produksi Pega
// (ADR-U-0005).
//
// Penandanya sempat tidak disimpan di lapisan ini karena belum ada
// pembacanya; sejak tiket 14 pembacanya ada, yaitu pelari migrasi yang
// menolak berjalan di produksi.
func (d *DB) PegaProduksi() bool { return d != nil && d.isPegaProd }

// Qualify mengembalikan nama objek berkualifikasi skema, mis. POOLDATA.T_X.
//
// Setiap query menyebut nama skema secara eksplisit (ADR-U-0033). Tidak ada
// query yang boleh menulis nama tabel telanjang.
func (d *DB) Qualify(objek string) (string, error) {
	objek = strings.TrimSpace(objek)
	if objek == "" {
		return "", ErrObjekTakBernama
	}
	if strings.Contains(objek, ".") {
		return "", fmt.Errorf("%w: %q sudah berkualifikasi", ErrObjekTakBernama, objek)
	}
	return d.skema + "." + objek, nil
}

var polaTransaksi = regexp.MustCompile(`(?is)\bCOMMIT\b`)

// PeriksaSQL menolak teks SQL yang mengurus transaksinya sendiri.
//
// Transaksi dibuka dan ditutup oleh aplikasi di lapisan services
// (ADR-U-0029). Hanya COMMIT yang ditolak: stored procedure ber-ROLLBACK
// justru DIIZINKAN aturan yang sama - ia dipanggil TERAKHIR dan penandanya
// diperiksa - sehingga menolak kata ROLLBACK akan menolak pemanggilan yang sah.
func PeriksaSQL(teks string) error {
	if polaTransaksi.MatchString(teks) {
		return fmt.Errorf("%w: %q", ErrTransaksiDiSQL, ringkas(teks))
	}
	return nil
}

func ringkas(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > 80 {
		return s[:80] + "..."
	}
	return s
}

// Tx adalah transaksi yang dibuka dan ditutup oleh services, bukan oleh SQL.
type Tx struct {
	tx *sql.Tx
}

// Mulai membuka transaksi. Pemanggilnya - services - yang wajib menutupnya.
func (d *DB) Mulai(ctx context.Context) (*Tx, error) {
	if d == nil || d.sql == nil {
		return nil, ErrTanpaOracle
	}
	t, err := d.sql.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("repository: membuka transaksi: %w", err)
	}
	return &Tx{tx: t}, nil
}

// Commit menutup transaksi dengan sukses.
func (t *Tx) Commit() error { return t.tx.Commit() }

// Rollback membatalkan transaksi. Aman dipanggil lewat defer.
func (t *Tx) Rollback() error {
	err := t.tx.Rollback()
	if errors.Is(err, sql.ErrTxDone) {
		return nil
	}
	return err
}

// Pintu pernyataan di dalam transaksi.
//
// Refactor bentuk B (30-09-2026): repository tiap modul kini tinggal di paket
// sendiri dan tidak dapat lagi menyentuh medan `tx`. Keempat metode ini
// meneruskan apa adanya ke `*sql.Tx` - nol perilaku baru. ⛔ Pemeriksaan
// `PeriksaSQL` TETAP kewajiban pemanggil, persis seperti sebelumnya.

// ExecContext menjalankan pernyataan tulis di dalam transaksi.
func (t *Tx) ExecContext(ctx context.Context, q string, args ...any) (sql.Result, error) {
	return t.tx.ExecContext(ctx, q, args...)
}

// QueryContext menjalankan kueri banyak baris di dalam transaksi.
func (t *Tx) QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, q, args...)
}

// QueryRowContext menjalankan kueri satu baris di dalam transaksi.
func (t *Tx) QueryRowContext(ctx context.Context, q string, args ...any) *sql.Row {
	return t.tx.QueryRowContext(ctx, q, args...)
}

// PrepareContext menyiapkan pernyataan di dalam transaksi.
func (t *Tx) PrepareContext(ctx context.Context, q string) (*sql.Stmt, error) {
	return t.tx.PrepareContext(ctx, q)
}

// Terisi menyatakan transaksi ini memegang *sql.Tx sungguhan. Tx kosong
// (`&Tx{}`) dipakai uji tanpa Oracle.
func (t *Tx) Terisi() bool { return t != nil && t.tx != nil }

// Pintu pernyataan di luar transaksi - pembacaan biasa.

// ExecContext menjalankan pernyataan di luar transaksi.
func (d *DB) ExecContext(ctx context.Context, q string, args ...any) (sql.Result, error) {
	return d.sql.ExecContext(ctx, q, args...)
}

// QueryContext menjalankan kueri banyak baris di luar transaksi.
func (d *DB) QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error) {
	return d.sql.QueryContext(ctx, q, args...)
}

// QueryRowContext menjalankan kueri satu baris di luar transaksi.
func (d *DB) QueryRowContext(ctx context.Context, q string, args ...any) *sql.Row {
	return d.sql.QueryRowContext(ctx, q, args...)
}

// Terhubung menyatakan DB ini memegang koneksi (bukan nil, bukan kosong).
func (d *DB) Terhubung() bool { return d != nil && d.sql != nil }

// NomorBerikut mengambil satu nomor dari sequence.
//
// "Sequence" adalah pembangkit angka berurut milik Oracle. ADR-U-0006
// menetapkan identitas seluruh tabel T_CLAIMLF_* berasal dari sini -
// bukan dari cap waktu, bukan dari teks yang disusun sendiri.
//
// Refactor bentuk B (30-09-2026): dipindah apa adanya dari
// `(*PohonKlaim).nomorBerikut`, sebab empat modul memakainya.
func (d *DB) NomorBerikut(ctx context.Context, tx *Tx, sequence string) (string, error) {
	nama, err := d.Qualify(sequence)
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf(`SELECT %s.NEXTVAL FROM DUAL`, nama)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	var n int64
	if err := tx.tx.QueryRowContext(ctx, q).Scan(&n); err != nil {
		return "", fmt.Errorf("repository: mengambil nomor dari %s: %w", sequence, err)
	}
	return fmt.Sprint(n), nil
}
