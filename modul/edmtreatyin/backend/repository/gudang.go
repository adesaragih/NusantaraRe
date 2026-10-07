package repository

// Asal: salinan modul/nbtreatyin/backend/repository/gudang.go (06-10-2026).
//
// Untuk apa berkas ini: GUDANG ORACLE modul EDM Treaty In - satu tipe yang
// memegang koneksi dan penomor bersama, beserta nama tabel dan bantuan kecil
// yang dipakai seluruh berkas paket ini.
//
// Tabel MILIK modul ini (migrasi 360-363): empat tabel proyeksi selisih T_POLIS_DIFFERENCE* /
// T_POLIS_XOL_LAYER_DIFFERENCE (selisih.go). Tabel MILIK modul lain yang dibaca DAN ditulis: generasi
// T_GENERAL_POLIS_TREATY + anak T_POLIS_* (nbtreatyin 320-328; generasi endorsemen PRODKE >= 1) dan
// T_WORK_POLIS / SEQ_WORK_POLIS (premiumlistlife - tabel kasus lintas-lini). Tabel WARISAN Pega yang dibaca:
// CURRENCY, MARKETINGOFFICER, REINSURANCETYPE, BUSINESS, AGENT, M_LOGIN_GO, master JSON M_TREATY_IN /
// M_TREATY_IN_EDM / M_TREATY_OUT (master_edm.go), popup TREATY_IN / TREATY_IN_EDM / TREATY_OUT2
// (bisnis_edm.go); yang ditulis: HISTORYAKSEPTASIPEGA, HISTORYAKSEPTASIPRODUCTION, JSON_POLIS (tanpa
// DATA_JSON), ACHIEVEMENT, TREATYINPRODUCTION. Seluruhnya lewat `db.Qualify` - skema eksplisit (AC 30).
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
	tabelMataUang      = "CURRENCY"
	tabelMO            = "MARKETINGOFFICER"
	tabelJenisReas     = "REINSURANCETYPE"
	tabelBisnis        = "BUSINESS"
	tabelAgen          = "AGENT"
	tabelProduksi      = "TREATYINPRODUCTION"
	tabelRiwayatPega   = "HISTORYAKSEPTASIPEGA"
	tabelAkunLogin     = "M_LOGIN_GO"
	sequenceKerjaPolis = "SEQ_WORK_POLIS"
)

var (
	// ErrKasusTidakAda - tidak ada kasus EDM Treaty In dengan ID itu.
	ErrKasusTidakAda = errors.New("repository: kasus EDM Treaty In tidak ada")
	// ErrGenerasiTertutup - baris generasi sudah ditutup (punya penerus);
	// tidak boleh disunting (spec-penyimpanan ID-10, AC 6).
	ErrGenerasiTertutup = errors.New("repository: generasi polis sudah ditutup dan tidak boleh disunting")
	// ErrTahapBerubah - kasus berpindah tahap di antara baca dan tulis.
	ErrTahapBerubah = errors.New("repository: kasus sudah berpindah tahap; muat ulang")
	// ErrNomorPolisSudahAda - nomor polis dibentuk SEKALI per berkas (AC 74).
	ErrNomorPolisSudahAda = errors.New("repository: berkas ini sudah bernomor polis")

	// ErrGenerasiPolisTidakAda - nomor polis tidak punya generasi bernomor di T_GENERAL_POLIS_TREATY (polis lama
	// belum dimuat, atau nomor salah).
	ErrGenerasiPolisTidakAda = errors.New("repository: generasi polis tidak ada di tabel relasional")
	// ErrGenerasiSudahDiendorse - generasi itu sudah punya penerus (UNIQUE OLD_POLIS_ID, ID-10, AC 3-4): endorsemen
	// serentak yang kedua ditolak basis data.
	ErrGenerasiSudahDiendorse = errors.New("repository: generasi polis sudah punya endorsemen berjalan atau penerus")
	// ErrSelisihBeku - proyeksi selisih SUMBER 'PEGA' (hasil migrasi) tidak boleh dibangun ulang aplikasi (ID-38).
	ErrSelisihBeku = errors.New("repository: selisih hasil migrasi (SUMBER PEGA) beku")
)

// Gudang adalah pintu Oracle modul EDM Treaty In.
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
