package repository

// Untuk apa berkas ini: HALAMAN <-> TABEL. Satu halaman kerja klaim disimpan ke T_GENERAL_CLAIM (kolom katalog), tujuh
// tabel anak T_CLAIM_*, dan T_CLAIM_ADJUSTMENT beserta tiga tabel cucunya - SEMUA di dalam transaksi pemanggil, digerakkan `models.TabelHeaderKlaim` / `TabelAnakKlaim` / `TabelAdjustment` /
// `TabelCucuAdjustment`.
//
//   - Tabel anak (tanpa rujukan dari luar): hapus lalu tulis ulang, NOURUT = urutan baris (1..n).
//   - T_CLAIM_ADJUSTMENT: ID STABIL - baris membawa `ID` (models.PropID); baris lama disunting di tempat, baris baru
//     disisipkan, baris yang hilang dihapus (tabel cucu ikut ON DELETE CASCADE). Baris yang sudah punya KOMITE_ID tidak
//     pernah dihapus (T_GENERAL_KOMITE.ADJUSTMENT_ID menunjuknya). NOURUT ditulis dua fase (negatif dulu) supaya
//     UNIQUE (CLAIM_ID, NOURUT) tidak bertabrakan saat urutan bergeser.
//   - SuggestList ("Claim History") TIDAK disimpan: CLAIM_ID di T_VIEW_SUGGEST dibatalkan 07-10-2026 (OQ-CP-17).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/claimprop/backend/models"
)

// ErrAdjustmentBerkomite - baris adjustment yang sudah diserahkan ke komite hendak dihapus.
var ErrAdjustmentBerkomite = errors.New("repository: baris adjustment yang sudah diserahkan ke komite tidak dapat dihapus")

// idBerikut - ID baris T_CLAIM_* dan T_VIEW_SUGGEST klaim (SEQ_T_CLAIM, migrasi 533).
func (g *Gudang) idBerikut(ctx context.Context, tx *db.Tx) (string, error) {
	return g.db.NomorBerikut(ctx, tx, "SEQ_T_CLAIM")
}

// SimpanHalaman menulis seluruh halaman klaim `id`.
func (g *Gudang) SimpanHalaman(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error {
	if err := g.simpanHeader(ctx, tx, id, h); err != nil {
		return err
	}
	for _, t := range models.TabelAnakKlaim {
		if err := g.tulisUlangAnak(ctx, tx, t, "CLAIM_ID", id, h.AmbilDaftar(t.Daftar)); err != nil {
			return err
		}
	}
	if err := g.simpanAdjustment(ctx, tx, id, h); err != nil {
		return err
	}
	// SuggestList ("Claim History") tidak disimpan: CLAIM_ID di T_VIEW_SUGGEST dibatalkan 07-10-2026 (OQ-CP-17).
	return nil
}

func (g *Gudang) simpanHeader(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error {
	tabel, err := g.db.Qualify(models.TabelHeaderKlaim.Nama)
	if err != nil {
		return err
	}
	var args []any
	for _, k := range models.TabelHeaderKlaim.Kolom {
		v, err := nilaiTulis(k, h.Ambil(k.Properti))
		if err != nil {
			return err
		}
		args = append(args, v...)
	}
	q := sqlUbahHeader(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	args = append(args, id)
	hasil, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("repository: menyimpan induk klaim: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "penyimpanan induk klaim")
}

// sqlUbahHeader - UPDATE seluruh kolom katalog header satu klaim (penampung ID terakhir).
func sqlUbahHeader(tabel string) string {
	var set []string
	n := 1
	for _, k := range models.TabelHeaderKlaim.Kolom {
		e, jml := ekspresiTulis(k, n)
		set = append(set, k.Kolom+" = "+e)
		n += jml
	}
	return fmt.Sprintf(`UPDATE %s SET %s WHERE ID = :%d`, tabel, strings.Join(set, ", "), n)
}

// sqlHapusAnak - seluruh baris anak/cucu satu induk.
func sqlHapusAnak(tabel, kolInduk string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE %s = :1`, tabel, kolInduk)
}

// sqlHapusAdjustment / sqlGeserNourut / sqlUbahAdjustment - penyuntingan baris adjustment ber-ID stabil.
func sqlHapusAdjustment(tabel string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE ID = :1 AND CLAIM_ID = :2`, tabel)
}

func sqlGeserNourut(tabel string) string {
	return fmt.Sprintf(`UPDATE %s SET NOURUT = -NOURUT WHERE CLAIM_ID = :1 AND NOURUT > 0`, tabel)
}

// ⛔ Urutan penampung di teks = urutan argumen: driver mengikat menurut POSISI kemunculan, bukan nomornya.
func sqlUbahAdjustment(tabel string) string {
	var set []string
	n := 1
	for _, k := range models.TabelAdjustment.Kolom {
		e, jml := ekspresiTulis(k, n)
		set = append(set, k.Kolom+" = "+e)
		n += jml
	}
	return fmt.Sprintf(`UPDATE %s SET %s, NOURUT = :%d WHERE ID = :%d AND CLAIM_ID = :%d`, tabel,
		strings.Join(set, ", "), n, n+1, n+2)
}

// sqlSisipBaris - INSERT satu baris tabel anak/cucu (ID, induk, NOURUT, kolom katalog).
func sqlSisipBaris(tabel string, t models.Tabel) string {
	kol := []string{"ID", t.KolomInduk, "NOURUT"}
	nilai := []string{":1", ":2", ":3"}
	n := 4
	for _, k := range t.Kolom {
		e, jml := ekspresiTulis(k, n)
		kol = append(kol, k.Kolom)
		nilai = append(nilai, e)
		n += jml
	}
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, tabel, strings.Join(kol, ", "), strings.Join(nilai, ", "))
}

func argsBaris(t models.Tabel, b models.Baris) ([]any, error) {
	var args []any
	for _, k := range t.Kolom {
		v, err := nilaiTulis(k, b[k.Properti])
		if err != nil {
			return nil, err
		}
		args = append(args, v...)
	}
	return args, nil
}

// tulisUlangAnak menghapus seluruh baris induk `indukID` lalu menyisipkan `rows`.
func (g *Gudang) tulisUlangAnak(ctx context.Context, tx *db.Tx, t models.Tabel, kolInduk, indukID string,
	rows []models.Baris) error {
	tabel, err := g.db.Qualify(t.Nama)
	if err != nil {
		return err
	}
	qHapus := sqlHapusAnak(tabel, kolInduk)
	if err := db.PeriksaSQL(qHapus); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, qHapus, indukID); err != nil {
		return fmt.Errorf("repository: mengosongkan %s: %w", t.Nama, err)
	}
	q := sqlSisipBaris(tabel, t)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	for i, b := range rows {
		rid, err := g.idBerikut(ctx, tx)
		if err != nil {
			return err
		}
		v, err := argsBaris(t, b)
		if err != nil {
			return err
		}
		args := append([]any{rid, indukID, i + 1}, v...)
		hasil, err := tx.ExecContext(ctx, q, args...)
		if err != nil {
			return fmt.Errorf("repository: menulis %s: %w", t.Nama, err)
		}
		if err := db.PastikanSatuBaris(hasil, "penulisan "+t.Nama); err != nil {
			return err
		}
	}
	return nil
}

// simpanAdjustment - ID stabil, NOURUT dua fase, cucu ditulis ulang.
func (g *Gudang) simpanAdjustment(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error {
	t := models.TabelAdjustment
	tabel, err := g.db.Qualify(t.Nama)
	if err != nil {
		return err
	}
	ada, err := g.idAdjustment(ctx, tx, id)
	if err != nil {
		return err
	}
	rows := h.AmbilDaftar(models.DaftarAdjustment)
	tetap := map[string]bool{}
	for _, b := range rows {
		if rid := b[models.PropID]; rid != "" && ada[rid] != nil {
			tetap[rid] = true
		}
	}
	qHapus := sqlHapusAdjustment(tabel)
	if err := db.PeriksaSQL(qHapus); err != nil {
		return err
	}
	for rid, komite := range ada {
		if tetap[rid] {
			continue
		}
		if *komite != "" {
			return ErrAdjustmentBerkomite
		}
		if _, err := tx.ExecContext(ctx, qHapus, rid, id); err != nil {
			return fmt.Errorf("repository: menghapus baris adjustment: %w", err)
		}
	}
	qNeg := sqlGeserNourut(tabel)
	if err := db.PeriksaSQL(qNeg); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, qNeg, id); err != nil {
		return fmt.Errorf("repository: menggeser urutan adjustment: %w", err)
	}
	qUbah := sqlUbahAdjustment(tabel)
	qSisip := sqlSisipBaris(tabel, t)
	for _, q := range []string{qUbah, qSisip} {
		if err := db.PeriksaSQL(q); err != nil {
			return err
		}
	}
	for i, b := range rows {
		v, err := argsBaris(t, b)
		if err != nil {
			return err
		}
		rid := b[models.PropID]
		if rid != "" && tetap[rid] {
			hasil, err := tx.ExecContext(ctx, qUbah, append(v, i+1, rid, id)...)
			if err != nil {
				return fmt.Errorf("repository: menyunting baris adjustment: %w", err)
			}
			if err := db.PastikanSatuBaris(hasil, "penyuntingan baris adjustment"); err != nil {
				return err
			}
		} else {
			if rid, err = g.idBerikut(ctx, tx); err != nil {
				return err
			}
			hasil, err := tx.ExecContext(ctx, qSisip, append([]any{rid, id, i + 1}, v...)...)
			if err != nil {
				return fmt.Errorf("repository: menyisipkan baris adjustment: %w", err)
			}
			if err := db.PastikanSatuBaris(hasil, "penyisipan baris adjustment"); err != nil {
				return err
			}
			b[models.PropID] = rid
		}
		for _, c := range models.TabelCucuAdjustment {
			j := models.JalurAdj(i+1, c.Daftar)
			if err := g.tulisUlangAnak(ctx, tx, c, c.KolomInduk, rid, h.AmbilDaftar(j)); err != nil {
				return err
			}
		}
	}
	return nil
}

// sqlIDAdjustment - ID baris adjustment satu klaim beserta KOMITE_ID-nya.
func sqlIDAdjustment(tabel string) string {
	return fmt.Sprintf(`SELECT ID, KOMITE_ID FROM %s WHERE CLAIM_ID = :1`, tabel)
}

// idAdjustment - ID baris adjustment tersimpan beserta KOMITE_ID-nya.
func (g *Gudang) idAdjustment(ctx context.Context, tx *db.Tx, id string) (map[string]*string, error) {
	tabel, err := g.db.Qualify(models.TabelAdjustment.Nama)
	if err != nil {
		return nil, err
	}
	q := sqlIDAdjustment(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := tx.QueryContext(ctx, q, id)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca ID adjustment: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := map[string]*string{}
	for rows.Next() {
		var rid, komite sql.NullString
		if err := rows.Scan(&rid, &komite); err != nil {
			return nil, fmt.Errorf("repository: memindai ID adjustment: %w", err)
		}
		k := komite.String
		out[rid.String] = &k
	}
	return out, rows.Err()
}

// potong - teks dipangkas ke lebar kolom tujuannya (rune, bukan byte).
func potong(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}

// ---------------------------------------------------------------- baca

// BacaHalaman membaca halaman klaim `id` (tx boleh nil).
func (g *Gudang) BacaHalaman(ctx context.Context, tx *db.Tx, id string) (*models.Halaman, error) {
	h := models.HalamanBaru()
	q1 := func(q string, args ...any) *sql.Row {
		if tx != nil {
			return tx.QueryRowContext(ctx, q, args...)
		}
		return g.db.QueryRowContext(ctx, q, args...)
	}
	qn := func(q string, args ...any) (*sql.Rows, error) {
		if tx != nil {
			return tx.QueryContext(ctx, q, args...)
		}
		return g.db.QueryContext(ctx, q, args...)
	}
	t := models.TabelHeaderKlaim
	tabel, err := g.db.Qualify(t.Nama)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, daftarBaca(t.Kolom, ""), tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	nilai := make([]sql.NullString, len(t.Kolom))
	tujuan := make([]any, len(t.Kolom))
	for i := range nilai {
		tujuan[i] = &nilai[i]
	}
	if err := q1(q, id).Scan(tujuan...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrKasusTidakAda
		}
		return nil, fmt.Errorf("repository: membaca induk klaim: %w", err)
	}
	for i, k := range t.Kolom {
		if v := nilaiBaca(k, nilai[i]); v != "" {
			h.Setel(k.Properti, v)
		}
	}
	for _, a := range models.TabelAnakKlaim {
		rows, err := g.bacaAnak(qn, a, "CLAIM_ID", id, false)
		if err != nil {
			return nil, err
		}
		h.SetelDaftar(a.Daftar, rows)
	}
	adj, err := g.bacaAnak(qn, models.TabelAdjustment, "CLAIM_ID", id, true)
	if err != nil {
		return nil, err
	}
	h.SetelDaftar(models.DaftarAdjustment, adj)
	for i, b := range adj {
		for _, c := range models.TabelCucuAdjustment {
			rows, err := g.bacaAnak(qn, c, c.KolomInduk, b[models.PropID], false)
			if err != nil {
				return nil, err
			}
			h.SetelDaftar(models.JalurAdj(i+1, c.Daftar), rows)
		}
	}
	return h, nil
}

func (g *Gudang) bacaAnak(qn func(string, ...any) (*sql.Rows, error), t models.Tabel, kolInduk, indukID string,
	denganID bool) ([]models.Baris, error) {
	tabel, err := g.db.Qualify(t.Nama)
	if err != nil {
		return nil, err
	}
	extra := ""
	if denganID {
		extra = "ID, KOMITE_ID, "
	}
	q := fmt.Sprintf(`SELECT %s%s FROM %s WHERE %s = :1 ORDER BY NOURUT`, extra, daftarBaca(t.Kolom, ""), tabel, kolInduk)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := qn(q, indukID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s: %w", t.Nama, err)
	}
	defer func() { _ = rows.Close() }()
	var out []models.Baris
	for rows.Next() {
		lebar := len(t.Kolom)
		if denganID {
			lebar += 2
		}
		nilai := make([]sql.NullString, lebar)
		tujuan := make([]any, lebar)
		for i := range nilai {
			tujuan[i] = &nilai[i]
		}
		if err := rows.Scan(tujuan...); err != nil {
			return nil, fmt.Errorf("repository: memindai %s: %w", t.Nama, err)
		}
		b := models.Baris{}
		awal := 0
		if denganID {
			b[models.PropID] = nilai[0].String
			if nilai[1].String != "" {
				b[models.PropKomiteID] = nilai[1].String
			}
			awal = 2
		}
		for i, k := range t.Kolom {
			if v := nilaiBaca(k, nilai[awal+i]); v != "" {
				b[k.Properti] = v
			}
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// angkaBulat - teks bilangan bulat ("" = 0).
func angkaBulat(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}

// waktuOracle - waktu sebagai teks pertukaran (`utils.TanggalWaktu`).
func waktuOracle(t time.Time) string { return t.Format(utils.TanggalWaktu) }
