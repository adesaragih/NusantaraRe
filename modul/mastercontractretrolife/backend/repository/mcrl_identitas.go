package repository

// Identitas baris baru - tiruan procedure `INSERT*_LIFE` (keputusan o, K3):
// `'1' || LPAD(<seq>.NEXTVAL, 6, '0')`.
//
// Sequence per tabel: `TREATYBUSINESS_LIFE_SEQ` DARI XML
// (`SaveTreatyBusinessAll_Life_SQL.xml` b84); empat lainnya dari body procedure
// yang dipanggil XML (`docs/procedure-bodies-from-dba.md`). DEV juga punya
// `TREATYSECURITY_LIFE_SEQ` dan `TREATY_REINSURER_LIFE_SEQ` - TIDAK dipakai.
//
// ⛔ PENYIMPANGAN SADAR KECIL dari `LPAD` (preseden `FormatIdentitasTCO`):
// nomor urut lebih dari 6 digit dipotong Oracle diam-diam dan melahirkan
// identitas bertabrakan; di sini ia GAGAL TERANG.

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"nusantarare/inti/backend/db"
)

// Sequence kelima tabel.
const (
	SeqTahun     = "TREATYYEAR_LIFE_SEQ"
	SeqKontrak   = "TREATYCONTRACT_LIFE_SEQ"
	SeqReinsurer = "TREATYREINSURER_LIFE_SEQ"
	SeqSecurity  = "TREATYSECURITYREINSURER_LIFE_SEQ"
	SeqBusiness  = "TREATYBUSINESS_LIFE_SEQ"
)

// LebarIdentitas - `lpad(…, 6, '0')`.
const LebarIdentitas = 6

var sequenceDikenal = map[string]bool{SeqTahun: true, SeqKontrak: true, SeqReinsurer: true, SeqSecurity: true,
	SeqBusiness: true}

var (
	// ErrIdentitasMelampauiLebar - nomor urut tidak muat di 6 digit.
	ErrIdentitasMelampauiLebar = errors.New(
		"repository: sequence number exceeds the 6-digit identity width; the sequence must be reviewed, not truncated")
	// ErrSequenceTakDikenal - nama sequence di luar lima milik modul ini.
	ErrSequenceTakDikenal = errors.New("repository: sequence does not belong to Master Contract Retro Life")
)

// FormatIdentitas merakit '1' + nomor ber-padding nol 6 digit.
func FormatIdentitas(n int64) (string, error) {
	ekor := strconv.FormatInt(n, 10)
	if n < 0 || len(ekor) > LebarIdentitas {
		return "", fmt.Errorf("%w: %d", ErrIdentitasMelampauiLebar, n)
	}
	return "1" + fmt.Sprintf("%0*d", LebarIdentitas, n), nil
}

func sqlNomorBerikut(seq string) string { return fmt.Sprintf(`SELECT %s.NEXTVAL FROM DUAL`, seq) }

// identitasBaru menerbitkan identitas baru di dalam transaksi pemanggil.
func (g *Gudang) identitasBaru(ctx context.Context, tx *db.Tx, sequence string) (string, error) {
	if !sequenceDikenal[sequence] {
		return "", fmt.Errorf("%w: %q", ErrSequenceTakDikenal, sequence)
	}
	if !tx.Terisi() {
		return "", errors.New("repository: a new identity requires a transaction")
	}
	nama, err := g.db.Qualify(sequence)
	if err != nil {
		return "", err
	}
	q := sqlNomorBerikut(nama)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var n int64
	if err := tx.QueryRowContext(ctx, q).Scan(&n); err != nil {
		return "", fmt.Errorf("repository: taking a number of %s: %w", sequence, err)
	}
	return FormatIdentitas(n)
}
