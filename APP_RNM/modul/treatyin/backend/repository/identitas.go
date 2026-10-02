package repository

// Jalur identitas kontrak - tiket 16, 17, 18, 19.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/modul/treatyin/backend/models"
)

// CariKontrakSerupa mengembalikan pengenal kontrak yang KUNCI ALAMI-nya sama -
// cedant + asal bisnis + tanggal mulai + sifat proporsi (tiket 16).
//
// ⛔ Tanggal BERAKHIR sengaja TIDAK ikut. `ADR-0040` §2 dan `SPEC-MODEL-DATA.md`
// §10.1 menyebut empat ruas, bukan lima - dua kontrak dengan cedant dan tanggal
// mulai yang sama tetapi berakhir berbeda tetap layak diperingatkan.
func (g *Gudang) CariKontrakSerupa(ctx context.Context, k models.Kontrak) ([]int64, error) {
	nama, err := g.db.Qualify("KONTRAK")
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT ID_KONTRAK FROM %s
		WHERE ID_CEDANT = :1 AND ID_ASAL_BISNIS = :2 AND SIFAT_PROPORSI = :3
		  AND TANGGAL_MULAI = TO_DATE(:4, 'YYYY-MM-DD')
		ORDER BY ID_KONTRAK ASC`, nama)
	baris, err := g.db.QueryContext(ctx, q, k.IDCedant, k.IDAsalBisnis, k.SifatProporsi, k.TanggalMulai)
	if err != nil {
		return nil, fmt.Errorf("repository: mencari kontrak berkunci alami sama: %w", err)
	}
	defer func() { _ = baris.Close() }()
	keluar := []int64{}
	for baris.Next() {
		var id int64
		if err := baris.Scan(&id); err != nil {
			return nil, fmt.Errorf("repository: membaca kontrak serupa: %w", err)
		}
		keluar = append(keluar, id)
	}
	return keluar, baris.Err()
}

// CariKontrakLewatNomorWarisan - tiket 17.
//
// ⛔ `NOMOR_KONTRAK_WARISAN IS NOT NULL` disebut EKSPLISIT walau pembandingnya
// sudah menyingkirkan NULL dengan sendirinya. Ia menyatakan maksudnya: baris
// non-warisan TIDAK PERNAH ikut terambil, dan pembaca berikutnya tidak perlu
// menyimpulkannya dari semantik NULL.
func (g *Gudang) CariKontrakLewatNomorWarisan(ctx context.Context, nomor string) ([]models.Kontrak, error) {
	nama, err := g.db.Qualify("KONTRAK")
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT ID_KONTRAK, NOMOR_KONTRAK_WARISAN, ID_CEDANT, ID_ASAL_BISNIS,
		SIFAT_PROPORSI, TO_CHAR(TANGGAL_MULAI,'YYYY-MM-DD'), TO_CHAR(TANGGAL_BERAKHIR,'YYYY-MM-DD')
		FROM %s WHERE NOMOR_KONTRAK_WARISAN IS NOT NULL AND UPPER(NOMOR_KONTRAK_WARISAN) = UPPER(:1)
		ORDER BY ID_KONTRAK ASC`, nama)
	baris, err := g.db.QueryContext(ctx, q, nomor)
	if err != nil {
		return nil, fmt.Errorf("repository: mencari lewat nomor warisan: %w", err)
	}
	defer func() { _ = baris.Close() }()
	// Daftar kosong adalah JAWABAN (tiket 17), bukan galat.
	keluar := []models.Kontrak{}
	for baris.Next() {
		var k models.Kontrak
		var warisan sql.NullString
		if err := baris.Scan(&k.ID, &warisan, &k.IDCedant, &k.IDAsalBisnis,
			&k.SifatProporsi, &k.TanggalMulai, &k.TanggalBerakhir); err != nil {
			return nil, fmt.Errorf("repository: membaca kontrak warisan: %w", err)
		}
		k.NomorKontrakWarisan = warisan.String
		keluar = append(keluar, k)
	}
	return keluar, baris.Err()
}

// PerbaruiKontrak menyunting ruas yang BUKAN kunci alami - tiket 18.
//
// ⛔ Pernyataan ini SENGAJA hanya menyentuh `NOMOR_KONTRAK_WARISAN`. Kelima
// ruas beku tidak ada di `SET`-nya sama sekali, sehingga sekalipun lapisan
// services kelak keliru, SQL ini tidak dapat mengubahnya. Dua lapis untuk satu
// invarian, dan yang kedua di sini sebab `INV-19` menjaga identitas.
func (g *Gudang) PerbaruiKontrak(ctx context.Context, k models.Kontrak) error {
	nama, err := g.db.Qualify("KONTRAK")
	if err != nil {
		return err
	}
	var warisan any
	if k.NomorKontrakWarisan != "" {
		warisan = k.NomorKontrakWarisan
	}
	q := fmt.Sprintf(`UPDATE %s SET NOMOR_KONTRAK_WARISAN = :1 WHERE ID_KONTRAK = :2`, nama)
	hasil, err := g.db.ExecContext(ctx, q, warisan, k.ID)
	if err != nil {
		return fmt.Errorf("repository: memperbarui KONTRAK: %w", err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: membaca cacah baris terubah: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("%w: %d", ErrKontrakTidakAda, k.ID)
	}
	return nil
}

// TambahVersi menambahkan versi pada kontrak yang sudah ada - tiket 19.
func (g *Gudang) TambahVersi(ctx context.Context, idKontrak int64, v models.VersiKontrak) (int64, error) {
	nama, err := g.db.Qualify("VERSI_KONTRAK")
	if err != nil {
		return 0, err
	}
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return 0, fmt.Errorf("repository: membuka transaksi: %w", err)
	}
	selesai := false
	defer func() {
		if !selesai {
			_ = tx.Rollback()
		}
	}()
	id, err := nomor(ctx, g, tx, "SEQ_TRIN_VERSI_KONTRAK")
	if err != nil {
		return 0, err
	}
	// Kolom wajib isi yang tiket 19 tidak bicarakan diisi dari versi yang ada
	// lewat sub-SELECT: menyalinnya ke Go lalu mengirimkannya kembali membuka
	// celah antara baca dan tulis, dan nilainya toh tidak diputuskan di sini.
	q := fmt.Sprintf(`INSERT INTO %s (ID_VERSI_KONTRAK, ID_KONTRAK, NOMOR_URUT_VERSI,
		KEADAAN_SIKLUS_HIDUP, NAMA_KONTRAK, KODE_MATA_UANG_KONTRAK, PERSEN_BAGIAN_NURE,
		BAGIAN_NURE_SERAGAM, MEMAKAI_BORDEREAUX, CARA_PEMBUKUAN, MEMAKAI_PRORATA, RETRO_BERGANDA)
		SELECT :1, :2, :3, :4, :5, v.KODE_MATA_UANG_KONTRAK, v.PERSEN_BAGIAN_NURE,
		       v.BAGIAN_NURE_SERAGAM, v.MEMAKAI_BORDEREAUX, v.CARA_PEMBUKUAN,
		       v.MEMAKAI_PRORATA, v.RETRO_BERGANDA
		  FROM %s v
		 WHERE v.ID_KONTRAK = :6
		   AND v.NOMOR_URUT_VERSI = (SELECT MAX(NOMOR_URUT_VERSI) FROM %s WHERE ID_KONTRAK = :7)`,
		nama, nama, nama)
	hasil, err := tx.ExecContext(ctx, q, id, idKontrak, v.NomorUrutVersi, v.KeadaanSiklusHidup,
		v.NamaKontrak, idKontrak, idKontrak)
	if err != nil {
		return 0, terjemahkanGalatOracle(fmt.Errorf("repository: menyisipkan versi: %w", err), v.NomorUrutVersi)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("repository: membaca cacah baris tersisip: %w", err)
	}
	if n == 0 {
		// Nol baris berarti kontraknya tidak punya versi mana pun - keadaan
		// yang tiket 14 larang, dan karena itu layak disebut apa adanya.
		return 0, fmt.Errorf("%w: %d tidak punya versi yang dapat disalin kepalanya", ErrKontrakTidakAda, idKontrak)
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("repository: menutup transaksi: %w", err)
	}
	selesai = true
	return id, nil
}
