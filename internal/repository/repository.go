// Package repository adalah SATU-SATUNYA lapisan yang menyentuh Oracle.
//
// Arah ketergantungan: handlers -> services -> repository. Tidak terbalik,
// tidak memotong. Paket ini tidak pernah mengimpor handlers atau services.
//
// Fase 0 - scaffold: nol query, nol aturan dagang. Yang ada di sini adalah
// bentuk yang harus diikuti setiap tiket penyimpanan.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	// [usulan] Driver Oracle murni Go; mendaftar dengan nama "oracle".
	_ "github.com/sijms/go-ora/v2"

	"nusantarare/internal/config"
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
	return &DB{sql: h, skema: cfg.OracleSchema}, nil
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

// Penanda IS_PEGA_PROD (ADR-U-0005) sengaja TIDAK disimpan di sini: ia sudah
// ada di config, dan menyalinnya ke lapisan ini hanya menambah keadaan yang
// belum ada pembacanya. Tiket yang perilakunya memang berbeda di produksi
// membawanya sendiri lewat services.

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
