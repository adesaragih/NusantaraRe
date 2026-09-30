package repository

// Identitas baris Treaty Contract Out dari sequence WARISAN - tiket 01
// (ADR-0006), tco4: sequence yang procedure `PEGA_*` sendiri pakai.
//
// Bentuk warisan (spec §7): `'1' || lpad(seq.nextval, N, '0')`, N = 6 untuk
// tahun, kontrak, reinsurer, business; N = 7 untuk klausul. Security TIDAK
// punya identitas. Pengguna TIDAK PERNAH mengetik identitas (AC 5).
//
// ⛔ PENYIMPANGAN SADAR KECIL dari `lpad`: bila nomor urut punya LEBIH banyak
// digit daripada lebarnya, Oracle `LPAD` memotongnya DIAM-DIAM dari kanan
// ('1234567' lebar 6 -> '123456') dan melahirkan identitas yang dapat
// bertabrakan. Di sini ia GAGAL TERANG (ADR-0015).

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"nusantarare/inti/backend/db"
)

// ErrIdentitasMelampauiLebar - nomor urut tidak muat di lebar identitas.
var ErrIdentitasMelampauiLebar = errors.New(
	"repository: sequence number exceeds the identity width; the sequence must be reviewed, not truncated")

// FormatIdentitasTCO merakit '1' + nomor ber-padding nol selebar `lebar`.
func FormatIdentitasTCO(n int64, lebar int) (string, error) {
	if n < 0 {
		return "", fmt.Errorf("%w: %d is negative", ErrIdentitasMelampauiLebar, n)
	}
	ekor := strconv.FormatInt(n, 10)
	if len(ekor) > lebar {
		return "", fmt.Errorf("%w: %d does not fit in %d digits", ErrIdentitasMelampauiLebar, n, lebar)
	}
	return "1" + fmt.Sprintf("%0*d", lebar, n), nil
}

// sequenceDikenalTCO membatasi nama yang boleh disebut di SQL.
var sequenceDikenalTCO = map[string]int{
	SeqTahunTCO: LebarIdentitasTCO, SeqKontrakTCO: LebarIdentitasTCO,
	SeqReinsurerTCO: LebarIdentitasTCO, SeqBusinessTCO: LebarIdentitasTCO,
	SeqKlausulTCO: LebarIdentitasKlausulTCO,
}

// ErrSequenceTakDikenal - nama sequence di luar daftar modul.
var ErrSequenceTakDikenal = errors.New("repository: sequence does not belong to Treaty Contract Out")

// IdentitasBerikutTCO menerbitkan identitas baru dari sequence, di dalam
// transaksi pemanggil.
func IdentitasBerikutTCO(ctx context.Context, d *db.DB, tx *db.Tx, sequence string) (string, error) {
	lebar, dikenal := sequenceDikenalTCO[sequence]
	if !dikenal {
		return "", fmt.Errorf("%w: %q", ErrSequenceTakDikenal, sequence)
	}
	nama, err := d.Qualify(sequence)
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf(`SELECT %s.NEXTVAL FROM DUAL`, nama)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var n int64
	if err := tx.QueryRowContext(ctx, q).Scan(&n); err != nil {
		return "", fmt.Errorf("repository: taking a number from %s: %w", sequence, err)
	}
	return FormatIdentitasTCO(n, lebar)
}
