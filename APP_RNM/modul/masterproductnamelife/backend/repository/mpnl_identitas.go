package repository

// Identitas produk baru - tiruan `PEGA_M_PRODUCT_LIFE` (prosedur TIDAK
// dipanggil, P6): `concat('1', lpad(M_PRODUCT_LIFE_SEQ.nextval, 5, '0'))`
// (`docs/dba-procedures-and-ddl.md` §1; dipicu `SaveProductNameLIfe` b58 dengan
// `IDPEGA` kosong). Baris inward memakai ID yang SAMA (R14; OQ-MPNL-02 ditutup data DEV 01-10-2026: 196/196 produk) -
// `M_PRODUCT_INWARD_LIFE_SEQ` tidak dipakai.
//
// ⛔ PENYIMPANGAN SADAR KECIL dari `LPAD` (preseden Retro Life / Treaty
// Contract Out): nomor urut lebih dari 5 digit dipotong Oracle diam-diam dan
// melahirkan identitas bertabrakan; di sini ia GAGAL TERANG. ID yang sudah
// dipakai induk flat tidak pernah digandakan atau ditimpa.
//
// ⭐ Keputusan work owner 02-10-2026 ("simpan ke table flat semua", "semua simpan dan baca dari table flat"): yang
// diperiksa HANYA induk flat `M_PRODUCTNAME_LIFE`. Kedua tabel JSON warisan tidak lagi dibaca aplikasi - ID produk
// lama sampai di induk flat lewat alat pindah, jadi alat pindah dijalankan SEBELUM aplikasi dipakai.
//
// ⛔ Audit 02-10-2026 (RALAT 02-10-2026): di DEV `M_PRODUCT_LIFE_SEQ` tertinggal dari data
// (nilai berikut 200, ID `100202` sudah ada). `MERGE` prosedur Pega akan MENIMPA produk itu
// diam-diam; dulu modul ini gagal 500 pada produk baru ketiga. Kini nomor yang ID-nya sudah
// dipakai DILEWATI ke nomor berikut (paling banyak `MaksLewatiIdentitas` berturut-turut, setiap
// lompatan dicatat di log supaya DBA melihatnya); lebih dari itu tetap gagal terang.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"nusantarare/inti/backend/db"
)

// MaksLewatiIdentitas - batas ID terpakai yang dilewati berturut-turut dalam satu penerbitan.
const MaksLewatiIdentitas = 100

// SeqProduk - sequence identitas produk.
const SeqProduk = "M_PRODUCT_LIFE_SEQ"

// LebarNomorIdentitas - `lpad(…, 5, '0')`.
const LebarNomorIdentitas = 5

var (
	// ErrIdentitasMelampauiLebar - nomor urut tidak muat 5 digit.
	ErrIdentitasMelampauiLebar = errors.New(
		"repository: sequence number exceeds the 5-digit product identity width; the sequence must be reviewed, not truncated")
	// ErrIdentitasBentrok - ID baru sudah dipakai baris lain (sequence tertinggal dari data).
	ErrIdentitasBentrok = errors.New(
		"repository: the new product ID is already used; the sequence is behind the existing data and must be reviewed by the DBA")
)

// FormatIdentitas merakit '1' + nomor ber-padding nol 5 digit.
func FormatIdentitas(n int64) (string, error) {
	ekor := strconv.FormatInt(n, 10)
	if n < 0 || len(ekor) > LebarNomorIdentitas {
		return "", fmt.Errorf("%w: %d", ErrIdentitasMelampauiLebar, n)
	}
	return "1" + fmt.Sprintf("%0*d", LebarNomorIdentitas, n), nil
}

func sqlCacahID(tabel string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ID = :1`, tabel)
}

// PilihIdentitasBebas mengambil nomor sequence berikut sampai ID-nya belum dipakai: ID terpakai
// dilewati (dikembalikan di `dilewati`), paling banyak `MaksLewatiIdentitas` berturut-turut.
// Nomor yang tidak muat 5 digit tetap gagal terang (`ErrIdentitasMelampauiLebar`).
func PilihIdentitasBebas(berikut func() (int64, error), terpakai func(id string) (bool, error)) (id string, dilewati []string, err error) {
	for range MaksLewatiIdentitas + 1 {
		n, err := berikut()
		if err != nil {
			return "", dilewati, err
		}
		calon, err := FormatIdentitas(n)
		if err != nil {
			return "", dilewati, err
		}
		pakai, err := terpakai(calon)
		if err != nil {
			return "", dilewati, err
		}
		if !pakai {
			return calon, dilewati, nil
		}
		dilewati = append(dilewati, calon)
	}
	return "", dilewati, fmt.Errorf("%w: %d consecutive IDs from %s are already used (%s to %s)", ErrIdentitasBentrok,
		len(dilewati), SeqProduk, dilewati[0], dilewati[len(dilewati)-1])
}

// identitasBaru menerbitkan ID produk baru di dalam transaksi pemanggil dan
// memastikan ID itu belum dipakai induk flat (ID terpakai dilewati, lihat kepala berkas).
func (g *Gudang) identitasBaru(ctx context.Context, tx *db.Tx) (string, error) {
	if !tx.Terisi() {
		return "", errors.New("repository: a new product identity requires a transaction")
	}
	berikut := func() (int64, error) {
		nomor, err := g.db.NomorBerikut(ctx, tx, SeqProduk)
		if err != nil {
			return 0, err
		}
		n, err := strconv.ParseInt(nomor, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("repository: %s returned %q: %w", SeqProduk, nomor, err)
		}
		return n, nil
	}
	q, err := g.siapkan(TabelFlatInduk, sqlCacahID)
	if err != nil {
		return "", err
	}
	terpakai := func(id string) (bool, error) {
		var c int
		if err := tx.QueryRowContext(ctx, q, id).Scan(&c); err != nil {
			return false, fmt.Errorf("repository: checking %s %s: %w", TabelFlatInduk, id, err)
		}
		return c > 0, nil
	}
	id, dilewati, err := PilihIdentitasBebas(berikut, terpakai)
	if len(dilewati) > 0 {
		log.Printf("master product name life: %s issued %d already used ID(s) %v; skipped - the sequence is behind the data and should be reviewed by the DBA",
			SeqProduk, len(dilewati), dilewati)
	}
	return id, err
}
