package repository

// Jalur pemulihan limit dan jalur baca warisan - tiket 32, 40.
//
// Tiket 41 TIDAK punya jalur repository: identitas bentuk lama diturunkan di
// services dari `BacaKontrak` yang sudah ada. Menambah kueri tersendiri
// untuknya berarti dua jalur membaca kolom yang sama, dan `INV-60` menuntut
// satu fakta satu penulis - termasuk satu pembaca kanonik.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/modul/treatyin/backend/models"
)

// CatatPemulihanLimit mengganti SELURUH pemulihan sebuah layer - tiket 32.
//
// ⛔ Hapus-lalu-sisip di dalam SATU transaksi, bukan sisip bertambah. Daftar
// pemulihan adalah satu pernyataan utuh tentang layernya: mengirim tiga baris
// berarti "pemulihan layer ini ada tiga", bukan "tambahkan tiga lagi". Tanpa
// penghapusannya, memperbaiki daftar dari tiga baris menjadi dua meninggalkan
// baris ketiga yang tidak seorang pun kirimkan.
func (g *Gudang) CatatPemulihanLimit(ctx context.Context, idLayer int64, baris []models.PemulihanLimit) error {
	tabel, err := g.db.Qualify("PEMULIHAN_LIMIT")
	if err != nil {
		return err
	}

	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return fmt.Errorf("repository: membuka transaksi: %w", err)
	}
	selesai := false
	defer func() {
		if !selesai {
			_ = tx.Rollback()
		}
	}()

	if _, err := tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE ID_LAYER = :1", tabel), idLayer); err != nil {
		return fmt.Errorf("repository: membersihkan PEMULIHAN_LIMIT layer %d: %w", idLayer, err)
	}

	q := fmt.Sprintf("INSERT INTO %s (ID_PEMULIHAN_LIMIT,ID_LAYER,NOMOR_URUT_PEMULIHAN,"+
		"PERSEN_PEMULIHAN,PERSEN_TAMBAHAN,CATATAN) VALUES "+
		"(:1,:2,:3,(TO_NUMBER(:4) / POWER(10, :5)),(TO_NUMBER(:6) / POWER(10, :7)),:8)", tabel)
	for _, b := range baris {
		id, err := nomor(ctx, g, tx, "SEQ_TRIN_PEMULIHAN_LIMIT")
		if err != nil {
			return err
		}
		koefP, skalaP := pecahAngka(b.PersenPemulihan)
		// ⚠️ Kosong tetap kosong. `TO_NUMBER(NULL)` menghasilkan NULL, jadi
		// pembagiannya pun NULL - tarif yang belum dinyatakan TIDAK menjadi
		// nol, dan itu beda yang tiket 32 ada untuk menjaganya.
		koefT, skalaT := pecahAngka(b.PersenTambahan)
		var catatan any
		if b.Catatan != "" {
			catatan = b.Catatan
		}
		if _, err := tx.ExecContext(ctx, q, id, idLayer, b.NomorUrut, koefP, skalaP, koefT, skalaT, catatan); err != nil {
			return fmt.Errorf("repository: menyisipkan PEMULIHAN_LIMIT urut %d: %w", b.NomorUrut, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("repository: menutup transaksi: %w", err)
	}
	selesai = true
	return nil
}

// BacaVersiDasar membaca versi yang menjadi dasar `idVersi` - tiket 40.
//
// ⛔ Nil, nil berarti "versi ini tidak punya dasar" - versi pertama. Itu
// BUKAN galat, dan membedakannya dari galat adalah pokok tiket 40.
//
// ⛔ Kolomnya diambil dari baris DASAR (alias `d`), bukan dari salinan mana
// pun. Satu-satunya hal yang dibaca dari baris anaknya adalah penunjuknya
// sendiri. Itulah bedanya dengan pohon `OLDDATA` sistem lama, dan ia terlihat
// hanya ketika versi dasarnya berubah sesudah versi anaknya dibuat.
func (g *Gudang) BacaVersiDasar(ctx context.Context, idVersi int64) (*models.VersiKontrak, error) {
	tabel, err := g.db.Qualify("VERSI_KONTRAK")
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT d.ID_VERSI_KONTRAK, d.ID_KONTRAK, d.NOMOR_URUT_VERSI,
		d.KEADAAN_SIKLUS_HIDUP, d.NAMA_KONTRAK, d.KODE_MATA_UANG_KONTRAK,
		d.BAGIAN_NURE_SERAGAM, d.MEMAKAI_BORDEREAUX, d.CARA_PEMBUKUAN,
		d.MEMAKAI_PRORATA, d.RETRO_BERGANDA
		FROM %s v JOIN %s d ON d.ID_VERSI_KONTRAK = v.ID_VERSI_KONTRAK_DASAR
		WHERE v.ID_VERSI_KONTRAK = :1`, tabel, tabel)

	var v models.VersiKontrak
	var urut sql.NullInt64
	err = g.db.QueryRowContext(ctx, q, idVersi).Scan(&v.ID, &v.IDKontrak, &urut,
		&v.KeadaanSiklusHidup, &v.NamaKontrak, &v.KodeMataUangKontak,
		&v.BagianNureSeragam, &v.MemakaiBordereaux, &v.CaraPembukuan,
		&v.MemakaiProrata, &v.RetroBerganda)
	if errors.Is(err, sql.ErrNoRows) {
		// Dua keadaan bermuara di sini dan keduanya sama artinya bagi
		// pemanggil: versinya tidak punya dasar, atau versinya sendiri tidak
		// ada. Yang kedua sudah dijawab jalur `BacaKontrak`.
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository: membaca versi dasar dari versi %d: %w", idVersi, err)
	}
	if urut.Valid {
		n := urut.Int64
		v.NomorUrutVersi = &n
	}
	return &v, nil
}
