package templat

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// Versi adalah satu berkas yang pernah diunggah untuk sebuah slot (tanpa isinya).
type Versi struct {
	Versi        int    `json:"versi"`
	NamaBerkas   string `json:"namaBerkas"`
	Ukuran       int    `json:"ukuran"`
	JumlahKolom  int    `json:"jumlahKolom"`
	Catatan      string `json:"catatan"`
	Aktif        bool   `json:"aktif"`
	DiunggahOleh string `json:"diunggahOleh"`
	// TglUnggah - `DD-MM-YYYY HH24:MI`.
	TglUnggah string `json:"tglUnggah"`
}

// ErrVersiTidakAda - versi yang diminta tidak ada untuk slot itu.
var ErrVersiTidakAda = errors.New("templat: versi tidak ada")

// ErrBelumDimigrasi - tabel M_TEMPLATE_FILE belum ada (migrasi inti 912 belum dijalankan).
var ErrBelumDimigrasi = errors.New("templat: tabel M_TEMPLATE_FILE belum ada - migrasi inti 912 belum dijalankan (-migrate)")

// Gudang menyimpan versi berkas templat.
type Gudang interface {
	// Aktif - versi aktif setiap kode yang punya versi aktif.
	Aktif(ctx context.Context) (map[string]Versi, error)
	// Riwayat - seluruh versi satu kode, terbaru dulu.
	Riwayat(ctx context.Context, kode string) ([]Versi, error)
	// Isi - isi berkas satu versi.
	Isi(ctx context.Context, kode string, versi int) ([]byte, error)
	// Sisip menyimpan versi baru dan menjadikannya SATU-SATUNYA versi aktif kode itu; menjawab nomor versinya.
	Sisip(ctx context.Context, kode string, v Versi, isi []byte) (int, error)
	// Aktifkan menjadikan versi itu satu-satunya yang aktif; 0 = kembali ke berkas bawaan (semua nonaktif).
	Aktifkan(ctx context.Context, kode string, versi int) error
}

// TabelTemplat dan SeqTemplat - migrasi inti 912.
const (
	TabelTemplat = "M_TEMPLATE_FILE"
	SeqTemplat   = "SEQ_M_TEMPLATE_FILE"
	benderaYa    = "1"
	benderaTidak = "0"
)

const kolomVersi = `VERSI, NAMA_BERKAS, UKURAN, NVL(JUMLAH_KOLOM, 0), NVL(CATATAN, ' '), AKTIF, DIUNGGAH_OLEH,
	  TO_CHAR(TGL_UNGGAH, 'DD-MM-YYYY HH24:MI')`

func sqlAktif(t string) string {
	return fmt.Sprintf(`SELECT KODE, %s FROM %s WHERE AKTIF = :1`, kolomVersi, t)
}

func sqlRiwayat(t string) string {
	return fmt.Sprintf(`SELECT KODE, %s FROM %s WHERE KODE = :1 ORDER BY VERSI DESC`, kolomVersi, t)
}

func sqlIsi(t string) string {
	return fmt.Sprintf(`SELECT ISI FROM %s WHERE KODE = :1 AND VERSI = :2`, t)
}

func sqlVersiBerikut(t string) string {
	return fmt.Sprintf(`SELECT NVL(MAX(VERSI), 0) + 1 FROM %s WHERE KODE = :1`, t)
}

func sqlNonaktif(t string) string {
	return fmt.Sprintf(`UPDATE %s SET AKTIF = :1 WHERE KODE = :2 AND AKTIF = :3`, t)
}

func sqlAktifkan(t string) string {
	return fmt.Sprintf(`UPDATE %s SET AKTIF = :1 WHERE KODE = :2 AND VERSI = :3`, t)
}

func sqlSisip(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, KODE, VERSI, NAMA_BERKAS, UKURAN, JUMLAH_KOLOM, ISI, CATATAN, AKTIF,
	  DIUNGGAH_OLEH, TGL_UNGGAH) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, SYSDATE)`, t)
}

// GudangOracle - Gudang di atas M_TEMPLATE_FILE.
type GudangOracle struct{ db *db.DB }

// NewGudangOracle membuat gudang Oracle.
func NewGudangOracle(d *db.DB) *GudangOracle { return &GudangOracle{db: d} }

func (g *GudangOracle) tabel() (string, error) { return g.db.Qualify(TabelTemplat) }

// bungkus menerjemahkan ORA-00942 menjadi ErrBelumDimigrasi.
func bungkus(err error, apa string) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "ORA-00942") || strings.Contains(err.Error(), "ORA-02289") {
		return fmt.Errorf("%w: %v", ErrBelumDimigrasi, err)
	}
	return fmt.Errorf("templat: %s: %w", apa, err)
}

func pindaiVersi(r *sql.Rows) (string, Versi, error) {
	var (
		kode, aktif string
		v           Versi
	)
	err := r.Scan(&kode, &v.Versi, &v.NamaBerkas, &v.Ukuran, &v.JumlahKolom, &v.Catatan, &aktif, &v.DiunggahOleh, &v.TglUnggah)
	v.Catatan = strings.TrimSpace(v.Catatan)
	v.Aktif = aktif == benderaYa
	return kode, v, err
}

func (g *GudangOracle) daftar(ctx context.Context, q string, args ...any) ([]string, []Versi, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, nil, bungkus(err, "membaca versi")
	}
	defer func() { _ = rows.Close() }()
	var (
		kode []string
		vs   []Versi
	)
	for rows.Next() {
		k, v, err := pindaiVersi(rows)
		if err != nil {
			return nil, nil, bungkus(err, "membaca versi")
		}
		kode, vs = append(kode, k), append(vs, v)
	}
	return kode, vs, bungkus(rows.Err(), "membaca versi")
}

// Aktif memenuhi Gudang.
func (g *GudangOracle) Aktif(ctx context.Context) (map[string]Versi, error) {
	t, err := g.tabel()
	if err != nil {
		return nil, err
	}
	kode, vs, err := g.daftar(ctx, sqlAktif(t), benderaYa)
	if err != nil {
		return nil, err
	}
	hasil := map[string]Versi{}
	for i, k := range kode {
		hasil[k] = vs[i]
	}
	return hasil, nil
}

// Riwayat memenuhi Gudang.
func (g *GudangOracle) Riwayat(ctx context.Context, kode string) ([]Versi, error) {
	t, err := g.tabel()
	if err != nil {
		return nil, err
	}
	_, vs, err := g.daftar(ctx, sqlRiwayat(t), kode)
	return vs, err
}

// Isi memenuhi Gudang.
func (g *GudangOracle) Isi(ctx context.Context, kode string, versi int) ([]byte, error) {
	t, err := g.tabel()
	if err != nil {
		return nil, err
	}
	q := sqlIsi(t)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	var p db.PindaiBiner
	if err := g.db.QueryRowContext(ctx, q, kode, versi).Scan(&p); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrVersiTidakAda
		}
		return nil, bungkus(err, "membaca isi berkas")
	}
	return p.Isi(), nil
}

// Sisip memenuhi Gudang - satu transaksi: nomor versi, nonaktifkan versi lain, sisip versi baru aktif.
func (g *GudangOracle) Sisip(ctx context.Context, kode string, v Versi, isi []byte) (versi int, err error) {
	t, err := g.tabel()
	if err != nil {
		return 0, err
	}
	for _, q := range []string{sqlVersiBerikut(t), sqlNonaktif(t), sqlSisip(t)} {
		if err := db.PeriksaSQL(q); err != nil {
			return 0, err
		}
	}
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return 0, bungkus(err, "memulai transaksi")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	if err = tx.QueryRowContext(ctx, sqlVersiBerikut(t), kode).Scan(&versi); err != nil {
		return 0, bungkus(err, "menghitung nomor versi")
	}
	id, err := g.db.NomorBerikut(ctx, tx, SeqTemplat)
	if err != nil {
		return 0, bungkus(err, "mengambil ID")
	}
	if _, err = tx.ExecContext(ctx, sqlNonaktif(t), benderaTidak, kode, benderaYa); err != nil {
		return 0, bungkus(err, "menonaktifkan versi lama")
	}
	var kolom any
	if v.JumlahKolom > 0 {
		kolom = v.JumlahKolom
	}
	if _, err = tx.ExecContext(ctx, sqlSisip(t), id, kode, versi, v.NamaBerkas, len(isi), kolom, db.Biner(isi),
		db.KosongJadiNil(v.Catatan), benderaYa, v.DiunggahOleh); err != nil {
		return 0, bungkus(err, "menyimpan versi baru")
	}
	if err = tx.Commit(); err != nil {
		return 0, bungkus(err, "commit")
	}
	return versi, nil
}

// Aktifkan memenuhi Gudang.
func (g *GudangOracle) Aktifkan(ctx context.Context, kode string, versi int) (err error) {
	t, err := g.tabel()
	if err != nil {
		return err
	}
	for _, q := range []string{sqlNonaktif(t), sqlAktifkan(t)} {
		if err := db.PeriksaSQL(q); err != nil {
			return err
		}
	}
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return bungkus(err, "memulai transaksi")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	if _, err = tx.ExecContext(ctx, sqlNonaktif(t), benderaTidak, kode, benderaYa); err != nil {
		return bungkus(err, "menonaktifkan versi")
	}
	if versi > 0 {
		hasil, e := tx.ExecContext(ctx, sqlAktifkan(t), benderaYa, kode, versi)
		if e != nil {
			err = bungkus(e, "mengaktifkan versi")
			return err
		}
		if n, _ := hasil.RowsAffected(); n != 1 {
			err = ErrVersiTidakAda
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return bungkus(err, "commit")
	}
	return nil
}
