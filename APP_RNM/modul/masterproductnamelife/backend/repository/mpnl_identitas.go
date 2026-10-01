package repository

// Identitas produk baru - tiruan `PEGA_M_PRODUCT_LIFE` (prosedur TIDAK
// dipanggil, P6): `concat('1', lpad(M_PRODUCT_LIFE_SEQ.nextval, 5, '0'))`
// (`docs/dba-procedures-and-ddl.md` §1; dipicu `SaveProductNameLIfe` b58 dengan
// `IDPEGA` kosong). Baris inward memakai ID yang SAMA (R14, OQ-MPNL-02) -
// `M_PRODUCT_INWARD_LIFE_SEQ` tidak dipakai.
//
// ⛔ PENYIMPANGAN SADAR KECIL dari `LPAD` (preseden Retro Life / Treaty
// Contract Out): nomor urut lebih dari 5 digit dipotong Oracle diam-diam dan
// melahirkan identitas bertabrakan; di sini ia GAGAL TERANG. ID yang sudah
// dipakai salah satu tabel (nol PK di DEV) juga ditolak, bukan digandakan.

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"nusantarare/inti/backend/db"
)

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

// identitasBaru menerbitkan ID produk baru di dalam transaksi pemanggil dan
// memastikan ID itu belum dipakai kedua tabel.
func (g *Gudang) identitasBaru(ctx context.Context, tx *db.Tx) (string, error) {
	if !tx.Terisi() {
		return "", errors.New("repository: a new product identity requires a transaction")
	}
	nomor, err := g.db.NomorBerikut(ctx, tx, SeqProduk)
	if err != nil {
		return "", err
	}
	n, err := strconv.ParseInt(nomor, 10, 64)
	if err != nil {
		return "", fmt.Errorf("repository: %s returned %q: %w", SeqProduk, nomor, err)
	}
	id, err := FormatIdentitas(n)
	if err != nil {
		return "", err
	}
	for _, tabel := range []string{TabelProduk, TabelInward} {
		q, err := g.siapkan(tabel, sqlCacahID)
		if err != nil {
			return "", err
		}
		var c int
		if err := tx.QueryRowContext(ctx, q, id).Scan(&c); err != nil {
			return "", fmt.Errorf("repository: checking %s %s: %w", tabel, id, err)
		}
		if c > 0 {
			return "", fmt.Errorf("%w: %s in %s", ErrIdentitasBentrok, id, tabel)
		}
	}
	return id, nil
}
