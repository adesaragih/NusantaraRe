package repository

// Jalur tulis dan baca KONTRAK + VERSI_KONTRAK - tiket 14.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
)

// ErrKontrakTidakAda - pengenal tidak menunjuk kontrak mana pun.
var ErrKontrakTidakAda = errors.New("repository: kontrak tidak ada")

// ErrNomorUrutVersiGanda - INV-04 ditolak Oracle; pesannya menyebut nomornya.
var ErrNomorUrutVersiGanda = errors.New("repository: nomor urut versi sudah dipakai")

// pecahAngka memecah desimal menjadi koefisien BULAT (teks) dan skala.
//
// ⛔ Kenapa tidak membind desimalnya langsung: jalur tulis yang membind teks
// apa adanya bergantung pada NLS_NUMERIC_CHARACTERS sesi, dan sesi ber-NLS
// Indonesia membaca "1.5" sebagai galat. Bilangan bulat tidak punya pemisah
// desimal, dan pembagian dengan pangkat sepuluh pada NUMBER itu eksak - nol
// float di jalur mana pun (ADR-0003).
//
// Pola ini sama dengan mastercontractretrolife/backend/repository; ia DISALIN,
// bukan dipinjam, sebab modul tidak pernah mengimpor modul lain (bab 7).
func pecahAngka(d *apd.Decimal) (any, int64) {
	if d == nil {
		return nil, 0
	}
	s := d.Text('f')
	tanda := ""
	if strings.HasPrefix(s, "-") {
		tanda, s = "-", s[1:]
	}
	bulat, pecahan, _ := strings.Cut(s, ".")
	koef := strings.TrimLeft(bulat+pecahan, "0")
	if koef == "" {
		return "0", 0
	}
	return tanda + koef, int64(len(pecahan))
}

// terjemahkanGalatOracle memetakan pelanggaran constraint menjadi galat yang
// berkata-kata.
//
// ⛔ Memeriksa NAMA constraint, bukan sekadar kode ORA: ORA-00001 dapat datang
// dari UQ_VERSI_KONTRAK maupun dari kunci utama, dan keduanya berarti hal yang
// berbeda bagi pengisi.
func terjemahkanGalatOracle(err error, nomorUrut *int64) error {
	if err == nil {
		return nil
	}
	pesan := strings.ToUpper(err.Error())
	if strings.Contains(pesan, "UQ_VERSI_KONTRAK") {
		n := "kosong"
		if nomorUrut != nil {
			n = strconv.FormatInt(*nomorUrut, 10)
		}
		// Daftar periksa tiket 14: "pesannya menyebut nomor yang bentrok".
		return fmt.Errorf("%w: nomor urut versi %s sudah ada pada kontrak itu", ErrNomorUrutVersiGanda, n)
	}
	return err
}

// BuatKontrakDenganVersiPertama menulis keduanya dalam SATU transaksi.
//
// Kontrak tanpa versi pertamanya bukan keadaan yang sah menurut tiket 14, dan
// dua pernyataan terpisah dapat meninggalkannya begitu bila yang kedua gagal.
//
// INV-02: kedua pengenal datang dari SEQUENCE. Tidak ada jalur lain di modul
// ini yang dapat memberi pengenal - pemanggil tidak dapat mengusulkannya.
func (g *Gudang) BuatKontrakDenganVersiPertama(ctx context.Context, k models.Kontrak, v models.VersiKontrak) (int64, int64, error) {
	tKontrak, err := g.db.Qualify("KONTRAK")
	if err != nil {
		return 0, 0, err
	}
	tVersi, err := g.db.Qualify("VERSI_KONTRAK")
	if err != nil {
		return 0, 0, err
	}

	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("repository: membuka transaksi: %w", err)
	}
	selesai := false
	defer func() {
		if !selesai {
			_ = tx.Rollback()
		}
	}()

	idK, err := nomor(ctx, g, tx, "SEQ_TRIN_KONTRAK")
	if err != nil {
		return 0, 0, err
	}
	idV, err := nomor(ctx, g, tx, "SEQ_TRIN_VERSI_KONTRAK")
	if err != nil {
		return 0, 0, err
	}

	qK := fmt.Sprintf("INSERT INTO %s (ID_KONTRAK,NOMOR_KONTRAK_WARISAN,ID_CEDANT,ID_ASAL_BISNIS,"+
		"SIFAT_PROPORSI,TANGGAL_MULAI,TANGGAL_BERAKHIR) VALUES (:1,:2,:3,:4,:5,TO_DATE(:6,'YYYY-MM-DD'),"+
		"TO_DATE(:7,'YYYY-MM-DD'))", tKontrak)
	var warisan any
	if k.NomorKontrakWarisan != "" {
		warisan = k.NomorKontrakWarisan
	}
	if _, err := tx.ExecContext(ctx, qK, idK, warisan, k.IDCedant, k.IDAsalBisnis,
		k.SifatProporsi, k.TanggalMulai, k.TanggalBerakhir); err != nil {
		return 0, 0, fmt.Errorf("repository: menyisipkan KONTRAK: %w", err)
	}

	koef, skala := pecahAngka(v.PersenBagianNure)
	qV := fmt.Sprintf("INSERT INTO %s (ID_VERSI_KONTRAK,ID_KONTRAK,NOMOR_URUT_VERSI,KEADAAN_SIKLUS_HIDUP,"+
		"NAMA_KONTRAK,KODE_MATA_UANG_KONTRAK,PERSEN_BAGIAN_NURE,BAGIAN_NURE_SERAGAM,MEMAKAI_BORDEREAUX,"+
		"CARA_PEMBUKUAN,MEMAKAI_PRORATA,RETRO_BERGANDA) VALUES "+
		"(:1,:2,:3,:4,:5,:6,(TO_NUMBER(:7) / POWER(10, :8)),:9,:10,:11,:12,:13)", tVersi)
	if _, err := tx.ExecContext(ctx, qV, idV, idK, v.NomorUrutVersi, v.KeadaanSiklusHidup,
		v.NamaKontrak, v.KodeMataUangKontak, koef, skala, v.BagianNureSeragam,
		v.MemakaiBordereaux, v.CaraPembukuan, v.MemakaiProrata, v.RetroBerganda); err != nil {
		return 0, 0, terjemahkanGalatOracle(fmt.Errorf("repository: menyisipkan VERSI_KONTRAK: %w", err), v.NomorUrutVersi)
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("repository: menutup transaksi: %w", err)
	}
	selesai = true
	return idK, idV, nil
}

// nomor mengambil satu pengenal dari sequence dan mengubahnya ke int64.
func nomor(ctx context.Context, g *Gudang, tx *db.Tx, sequence string) (int64, error) {
	s, err := g.db.NomorBerikut(ctx, tx, sequence)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("repository: nomor dari %s bukan bilangan: %w", sequence, err)
	}
	return n, nil
}

// BacaKontrak membaca kepala kontrak beserta SELURUH versinya.
//
// Lapisan beku dibaca SEKALI; ia tidak disalin ke tiap versi (ADR-0040).
func (g *Gudang) BacaKontrak(ctx context.Context, id int64) (models.KontrakDenganVersi, error) {
	var out models.KontrakDenganVersi
	tKontrak, err := g.db.Qualify("KONTRAK")
	if err != nil {
		return out, err
	}
	tVersi, err := g.db.Qualify("VERSI_KONTRAK")
	if err != nil {
		return out, err
	}

	qK := fmt.Sprintf("SELECT ID_KONTRAK,NOMOR_KONTRAK_WARISAN,ID_KONTRAK_DISALIN_DARI,ID_CEDANT,"+
		"ID_ASAL_BISNIS,SIFAT_PROPORSI,TO_CHAR(TANGGAL_MULAI,'YYYY-MM-DD'),"+
		"TO_CHAR(TANGGAL_BERAKHIR,'YYYY-MM-DD') FROM %s WHERE ID_KONTRAK = :1", tKontrak)
	var warisan sql.NullString
	var disalin sql.NullInt64
	k := &out.Kontrak
	switch err := g.db.QueryRowContext(ctx, qK, id).Scan(&k.ID, &warisan, &disalin, &k.IDCedant,
		&k.IDAsalBisnis, &k.SifatProporsi, &k.TanggalMulai, &k.TanggalBerakhir); {
	case errors.Is(err, sql.ErrNoRows):
		return out, fmt.Errorf("%w: %d", ErrKontrakTidakAda, id)
	case err != nil:
		return out, fmt.Errorf("repository: membaca KONTRAK: %w", err)
	}
	k.NomorKontrakWarisan = warisan.String
	if disalin.Valid {
		d := disalin.Int64
		k.IDKontrakDisalinDari = &d
	}

	qV := fmt.Sprintf("SELECT ID_VERSI_KONTRAK,ID_KONTRAK,NOMOR_URUT_VERSI,KEADAAN_SIKLUS_HIDUP,"+
		"NAMA_KONTRAK,KODE_MATA_UANG_KONTRAK,BAGIAN_NURE_SERAGAM,MEMAKAI_BORDEREAUX,CARA_PEMBUKUAN,"+
		"MEMAKAI_PRORATA,RETRO_BERGANDA FROM %s WHERE ID_KONTRAK = :1 "+
		"ORDER BY NOMOR_URUT_VERSI ASC NULLS LAST, ID_VERSI_KONTRAK ASC", tVersi)
	baris, err := g.db.QueryContext(ctx, qV, id)
	if err != nil {
		return out, fmt.Errorf("repository: membaca VERSI_KONTRAK: %w", err)
	}
	defer func() { _ = baris.Close() }()
	out.Versi = []models.VersiKontrak{}
	for baris.Next() {
		var v models.VersiKontrak
		var nomorUrut sql.NullInt64
		if err := baris.Scan(&v.ID, &v.IDKontrak, &nomorUrut, &v.KeadaanSiklusHidup, &v.NamaKontrak,
			&v.KodeMataUangKontak, &v.BagianNureSeragam, &v.MemakaiBordereaux, &v.CaraPembukuan,
			&v.MemakaiProrata, &v.RetroBerganda); err != nil {
			return out, fmt.Errorf("repository: membaca baris VERSI_KONTRAK: %w", err)
		}
		if nomorUrut.Valid {
			n := nomorUrut.Int64
			v.NomorUrutVersi = &n
		}
		out.Versi = append(out.Versi, v)
	}
	if err := baris.Err(); err != nil {
		return out, fmt.Errorf("repository: membaca VERSI_KONTRAK: %w", err)
	}
	return out, nil
}
