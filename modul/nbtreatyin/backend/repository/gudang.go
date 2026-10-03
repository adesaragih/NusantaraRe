package repository

// Untuk apa berkas ini: GUDANG ORACLE modul NB Treaty In - satu tipe yang
// memegang koneksi dan penomor bersama, beserta nama tabel dan bantuan kecil
// yang dipakai seluruh berkas paket ini.
//
// Tabel MILIK modul ini (migrasi 320-330): T_GENERAL_POLIS dan anak-anaknya
// T_POLIS_*, serta M_NBTRIN_PERAN_TEMPAT. Tabel MILIK modul lain yang ditulis:
// T_WORK_POLIS (premiumlistlife 050/059 - tabel kasus lintas-lini, dipakai
// bersama sesuai rancangan §2). Tabel WARISAN Pega yang dibaca: view
// TREATYINDETAILJOINEDM dan TREATYINDETAIL, CURRENCY, MARKETINGOFFICER,
// REINSURANCETYPE, TREATYGROUP, BUSINESS, CLIENT, AGENT, TREATYINPRODUCTION;
// yang ditulis: HISTORYAKSEPTASIPEGA (riwayat, sama dengan Pega). Seluruhnya
// lewat `db.Qualify` - skema eksplisit (AC 30).
//
// ⛔ Nol COMMIT di SQL (ADR-U-0029): transaksi milik services.

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penomor"
)

// Nama tabel - tanpa skema; skema dipasang `Qualify`.
const (
	tabelKerja         = "T_WORK_POLIS"
	tabelPeranTempat   = "M_NBTRIN_PERAN_TEMPAT"
	viewDetailGabung   = "TREATYINDETAILJOINEDM"
	tabelDetail        = "TREATYINDETAIL"
	tabelMataUang      = "CURRENCY"
	tabelMO            = "MARKETINGOFFICER"
	tabelJenisReas     = "REINSURANCETYPE"
	tabelGrupTreaty    = "TREATYGROUP"
	tabelBisnis        = "BUSINESS"
	tabelKlien         = "CLIENT"
	tabelAgen          = "AGENT"
	tabelProduksi      = "TREATYINPRODUCTION"
	tabelRiwayatPega   = "HISTORYAKSEPTASIPEGA"
	tabelAkunLogin     = "M_LOGIN_GO"
	sequenceKerjaPolis = "SEQ_WORK_POLIS"
)

var (
	// ErrKasusTidakAda - tidak ada kasus NB Treaty In dengan ID itu.
	ErrKasusTidakAda = errors.New("repository: kasus NB Treaty In tidak ada")
	// ErrGenerasiTertutup - baris generasi sudah ditutup (TGL_TUTUP terisi);
	// tidak boleh disunting (spec-penyimpanan ID-10, AC 6).
	ErrGenerasiTertutup = errors.New("repository: generasi polis sudah ditutup dan tidak boleh disunting")
	// ErrTahapBerubah - kasus berpindah tahap di antara baca dan tulis.
	ErrTahapBerubah = errors.New("repository: kasus sudah berpindah tahap; muat ulang")
	// ErrNomorPolisSudahAda - nomor polis dibentuk SEKALI per berkas (AC 74).
	ErrNomorPolisSudahAda = errors.New("repository: berkas ini sudah bernomor polis")
	// ErrDataKontrakTidakAda - baris view kontrak yang dipilih tidak ada
	// (AC 36-38: pembacaan gagal menghentikan proses).
	ErrDataKontrakTidakAda = errors.New("repository: data kontrak treaty yang dipilih tidak ditemukan di view")
	// ErrTipeNomorKosong - DueTo bukan 1/0 dan ClaimType bukan "XOL Retro":
	// nomor polis tanpa huruf tipe tidak diterbitkan.
	ErrTipeNomorKosong = errors.New("repository: tipe nomor polis kosong (Due To belum terisi)")
	// ErrOJKKosong - OJKBusinessID grup treaty kosong; nomor polis memuatnya.
	ErrOJKKosong = errors.New("repository: OJKBusinessID kosong; nomor polis memuatnya")
)

// Gudang adalah pintu Oracle modul NB Treaty In.
type Gudang struct {
	db    *db.DB
	nomor *penomor.Penomor
}

// Baru menyusun gudang di atas satu basis data.
func Baru(d *db.DB) *Gudang { return &Gudang{db: d, nomor: penomor.NewPenomor(d)} }

// nama memasang skema pada nama tabel.
func (g *Gudang) nama(t string) (string, error) { return g.db.Qualify(t) }

// pemeriksa membaca baris, di dalam transaksi bila ada.
type pemeriksa interface {
	QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, q string, args ...any) *sql.Row
}

func (g *Gudang) pembaca(tx *db.Tx) pemeriksa {
	if tx.Terisi() {
		return tx
	}
	return g.db
}

// jalankan menjalankan satu pernyataan tulis di dalam transaksi.
func jalankan(ctx context.Context, tx *db.Tx, apa, q string, args ...any) (sql.Result, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	hasil, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: %s: %w", apa, err)
	}
	return hasil, nil
}

// idBaru membangkitkan kunci baris anak: 32 karakter heksadesimal acak, pas
// VARCHAR2(32). Kunci ini bukan kunci dagang dan bukan kunci pasangan antar
// generasi (itu NOURUT, ID-11..13), sehingga tidak perlu sequence.
func idBaru() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("repository: membangkitkan kunci baris: %w", err)
	}
	return hex.EncodeToString(b[:]), nil
}

// teks membaca NullString sebagai teks ("" bila NULL).
func teks(v sql.NullString) string { return v.String }
