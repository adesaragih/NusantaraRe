package db

// Koneksi - SATU koneksi terkunci dari pool (`*sql.Conn`). Untuk pekerjaan yang menyetel sesi (`ALTER SESSION`) lalu
// membuka transaksi: pada `*sql.DB` keduanya dapat jatuh ke koneksi pool yang BERBEDA, sehingga setelan sesi diam-diam
// tidak berlaku. Di sini pernyataan sesi dan transaksi berjalan di koneksi yang sama sampai Close.
// Dipakai alat pindahflat modul ricommlife (keputusan work owner 06-10-2026).

import (
	"context"
	"database/sql"
	"fmt"
)

// Koneksi membungkus *sql.Conn.
type Koneksi struct{ c *sql.Conn }

// Koneksi mengambil satu koneksi dari pool; pemanggil WAJIB Close.
func (d *DB) Koneksi(ctx context.Context) (*Koneksi, error) {
	if d == nil || d.sql == nil {
		return nil, ErrTanpaOracle
	}
	c, err := d.sql.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository: mengambil koneksi: %w", err)
	}
	return &Koneksi{c: c}, nil
}

// Close mengembalikan koneksi ke pool.
func (k *Koneksi) Close() error { return k.c.Close() }

// Mulai membuka transaksi DI KONEKSI INI.
func (k *Koneksi) Mulai(ctx context.Context) (*Tx, error) {
	t, err := k.c.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("repository: membuka transaksi: %w", err)
	}
	return &Tx{tx: t}, nil
}

// ExecContext menjalankan pernyataan di koneksi ini (di luar transaksi).
func (k *Koneksi) ExecContext(ctx context.Context, q string, args ...any) (sql.Result, error) {
	return k.c.ExecContext(ctx, q, args...)
}

// QueryContext menjalankan kueri banyak baris di koneksi ini.
func (k *Koneksi) QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error) {
	return k.c.QueryContext(ctx, q, args...)
}

// QueryRowContext menjalankan kueri satu baris di koneksi ini.
func (k *Koneksi) QueryRowContext(ctx context.Context, q string, args ...any) *sql.Row {
	return k.c.QueryRowContext(ctx, q, args...)
}
