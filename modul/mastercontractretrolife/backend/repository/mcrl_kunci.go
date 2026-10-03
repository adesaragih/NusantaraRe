package repository

// Kunci baris induk - perbaikan /code-review 01-10-2026.
//
// ⛔ K2: kelima tabel tanpa FK. Tanpa kunci, anak yang ditulis BERSAMAAN dengan kaskade hapus induknya
// (READ COMMITTED) dapat menjadi yatim: penulis anak membaca induk yang masih ada, kaskade menghapusnya
// dan commit, lalu anak commit. Kedua pihak karena itu mengunci baris INDUK yang sama lebih dulu
// (`SELECT … FOR UPDATE`): yang kedua menunggu, lalu melihat keadaan sesudah yang pertama - kaskade
// menghitung anak baru itu (cacah berbeda → 409), atau penulis anak tidak menemukan induknya (404).

import (
	"context"
	"errors"
	"fmt"

	"nusantarare/inti/backend/db"
)

// errJenisKunci - jenis baris yang tidak dikenal kunci (termasuk tahun treaty: abadi, nol penghapus).
var errJenisKunci = errors.New("repository: unknown row kind to lock")

// tabelKunci memetakan jenis baris (nama yang sama dengan jenis hapus layanan) ke tabelnya.
func tabelKunci(jenis string) (string, error) {
	switch jenis {
	case "kontrak":
		return TabelKontrak, nil
	case "reinsurer":
		return TabelReinsurer, nil
	case "security":
		return TabelSecurity, nil
	case "business":
		return TabelBusiness, nil
	}
	return "", fmt.Errorf("%w: %q", errJenisKunci, jenis)
}

func sqlKunci(t string) string {
	return fmt.Sprintf(`SELECT ID FROM %s WHERE ID = :1 FOR UPDATE`, t)
}

// KunciBaris mengunci satu baris sampai transaksi tx selesai; baris tidak ada = ErrTidakAda.
func (g *Gudang) KunciBaris(ctx context.Context, tx *db.Tx, jenis, id string) error {
	if !tx.Terisi() {
		return fmt.Errorf("repository: locking %s %s requires a transaction", jenis, id)
	}
	tabel, err := tabelKunci(jenis)
	if err != nil {
		return err
	}
	q, err := g.siapkan(tabel, sqlKunci)
	if err != nil {
		return err
	}
	rows, err := tx.QueryContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("repository: locking %s %s: %w", jenis, id, err)
	}
	defer func() { _ = rows.Close() }()
	ada := rows.Next()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("repository: locking %s %s: %w", jenis, id, err)
	}
	if !ada {
		return fmt.Errorf("%w: %s %s", ErrTidakAda, jenis, id)
	}
	return nil
}
