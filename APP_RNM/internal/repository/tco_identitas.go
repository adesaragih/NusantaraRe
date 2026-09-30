package repository

// Identitas baris Treaty Contract Out dari sequence - tiket 01 (ADR-0006).
//
// Bentuk warisan dipertahankan (spec §7): `'1' || lpad(seq.nextval, N, '0')`,
// N = 6 untuk tahun, kontrak, reinsurer, security, business; N = 7 untuk
// klausul. Pengguna TIDAK PERNAH mengetik identitas (AC 5).
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
)

// ErrIdentitasMelampauiLebar - nomor urut tidak muat di lebar identitas.
var ErrIdentitasMelampauiLebar = errors.New(
	"repository: nomor urut melampaui lebar identitas; sequence harus ditinjau, bukan dipotong")

// FormatIdentitasTCO merakit '1' + nomor ber-padding nol selebar `lebar`.
func FormatIdentitasTCO(n int64, lebar int) (string, error) {
	if n < 0 {
		return "", fmt.Errorf("%w: %d negatif", ErrIdentitasMelampauiLebar, n)
	}
	ekor := strconv.FormatInt(n, 10)
	if len(ekor) > lebar {
		return "", fmt.Errorf("%w: %d tidak muat di %d digit", ErrIdentitasMelampauiLebar, n, lebar)
	}
	return "1" + fmt.Sprintf("%0*d", lebar, n), nil
}

// sequenceDikenalTCO membatasi nama yang boleh disebut di SQL.
var sequenceDikenalTCO = map[string]int{
	SeqTahunTCO: LebarIdentitasTCO, SeqKontrakTCO: LebarIdentitasTCO,
	SeqReinsurerTCO: LebarIdentitasTCO, SeqSecurityTCO: LebarIdentitasTCO,
	SeqBusinessTCO: LebarIdentitasTCO, SeqKlausulTCO: LebarIdentitasKlausulTCO,
	SeqJejakTCO: LebarIdentitasTCO,
}

// ErrSequenceTakDikenal - nama sequence di luar daftar modul.
var ErrSequenceTakDikenal = errors.New("repository: sequence bukan milik Treaty Contract Out")

// IdentitasBerikutTCO menerbitkan identitas baru dari sequence, di dalam
// transaksi pemanggil.
func (d *DB) IdentitasBerikutTCO(ctx context.Context, tx *Tx, sequence string) (string, error) {
	lebar, dikenal := sequenceDikenalTCO[sequence]
	if !dikenal {
		return "", fmt.Errorf("%w: %q", ErrSequenceTakDikenal, sequence)
	}
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
	return FormatIdentitasTCO(n, lebar)
}
