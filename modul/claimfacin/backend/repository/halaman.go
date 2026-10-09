package repository

// Untuk apa berkas ini: HALAMAN <-> TABEL. Satu halaman kerja klaim Fac In disimpan ke T_GENERAL_CLAIM (kolom katalog),
// pohon objek -> item -> estimasi / spreading / Break QS / adjustment -> cucu adjustment, retro tingkat klaim, dan
// kronologi T_VIEW_SUGGEST - SEMUA di dalam transaksi pemanggil, digerakkan `models.SemuaTabel` (katalog). Pola
// penulisan disalin dari `modul/claimnonprop/backend/repository/halaman.go` (bukan diimpor) dan dibuat rekursif:
//
//   - Simpul ber-ID STABIL (T_CLAIM_OBJECT, T_CLAIM_OBJECT_ITEM, T_CLAIM_ADJUSTMENT): baris membawa `ID`
//     (models.PropID); baris lama disunting di tempat, baris baru disisipkan, baris yang hilang dihapus (anak ikut ON
//     DELETE CASCADE). Adjustment ber-KOMITE_ID tidak pernah terhapus, langsung maupun lewat induknya.
//   - Simpul lain: seluruh baris kasus dihapus lalu ditulis ulang. Tabel ber-CLAIM_ID (`DenganKlaim`) dihapus per klaim
//     dan NOURUT-nya berurut SE-KLAIM per tabel (UNIQUE (CLAIM_ID, NOURUT) tabel Claim Prop); cucu adjustment
//     (ADJUSTMENT_ID) per adjustment, NOURUT 1..n.
//   - T_CLAIM_ADJUSTMENT: NOURUT ditulis dua fase (negatif dulu) supaya UNIQUE (CLAIM_ID, NOURUT) tidak bertabrakan saat
//     urutan bergeser.
//   - T_VIEW_SUGGEST (tabel bersama): HANYA BERTAMBAH - baris Chronology ber-penanda `Baru` disisipkan dengan NO berikutnya
//     per klaim.

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/claimfacin/backend/models"
)

// ErrAdjustmentBerkomite - baris adjustment yang sudah diserahkan ke komite hendak dihapus (langsung atau lewat item /
// objek induknya).
var ErrAdjustmentBerkomite = errors.New("repository: baris adjustment yang sudah diserahkan ke komite tidak dapat dihapus")

// idBerikut - ID baris T_CLAIM_* (SEQ_T_CLAIM, urutan tabel bersama T_CLAIM_*) dan T_VIEW_SUGGEST klaim.
func (g *Gudang) idBerikut(ctx context.Context, tx *db.Tx) (string, error) {
	return g.db.NomorBerikut(ctx, tx, "SEQ_T_CLAIM")
}

// saringKlaim - klausa WHERE seluruh baris simpul `t` milik klaim (bind :1 = ID klaim, JENIS bila ada); `tabelAdj` =
// T_CLAIM_ADJUSTMENT berskema (cucu adjustment tanpa CLAIM_ID).
func saringKlaim(t *models.Tabel, tabelAdj string) (string, []string) {
	switch {
	case t == &models.TabelAdjSpreading || t == &models.TabelAdjQuotaShare:
		return fmt.Sprintf("ADJUSTMENT_ID IN (SELECT ID FROM %s WHERE CLAIM_ID = :1)", tabelAdj), nil
	case t == &models.TabelRetroTreaty:
		return "CLAIM_ID = :1 AND ADJUSTMENT_ID IS NULL", nil
	case t == &models.TabelAdjRetro:
		return "CLAIM_ID = :1 AND ADJUSTMENT_ID IS NOT NULL", nil
	case t.Jenis != "":
		return "CLAIM_ID = :1 AND JENIS = :2", []string{t.Jenis}
	}
	return "CLAIM_ID = :1", nil
}

// sqlBacaSimpul - seluruh baris simpul `t` milik satu klaim: ID, induk langsung, kolom katalog (+ KOMITE_ID adjustment).
func sqlBacaSimpul(t *models.Tabel, tabel, tabelAdj string) (string, []string) {
	w, jenis := saringKlaim(t, tabelAdj)
	ekstra := ""
	if t == &models.TabelAdjustment {
		ekstra = ", KOMITE_ID"
	}
	return fmt.Sprintf(`SELECT ID, %s%s, %s FROM %s WHERE %s ORDER BY NOURUT`, t.KolomInduk, ekstra,
		daftarBaca(t.Kolom, ""), tabel, w), jenis
}

// SimpanHalaman menulis seluruh halaman klaim `id`.
func (g *Gudang) SimpanHalaman(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error {
	if err := g.simpanHeader(ctx, tx, id, h); err != nil {
		return err
	}
	if err := g.jagaKomite(ctx, tx, id, h); err != nil {
		return err
	}
	if err := g.kosongkanTulisUlang(ctx, tx, id); err != nil {
		return err
	}
	if err := g.geserNourut(ctx, tx, id); err != nil {
		return err
	}
	urut := map[string]int{}
	if err := g.simpanSimpul(ctx, tx, id, &models.TabelObjek, id, models.DaftarObjek, h, urut); err != nil {
		return err
	}
	if err := g.simpanSimpul(ctx, tx, id, &models.TabelRetroTreaty, id, models.DaftarFacRetroTreaty, h, urut); err != nil {
		return err
	}
	return g.sisipKronologi(ctx, tx, id, h)
}

func (g *Gudang) simpanHeader(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error {
	t := models.TabelHeaderKlaim
	tabel, err := g.db.Qualify(t.Nama)
	if err != nil {
		return err
	}
	var set []string
	var args []any
	n := 1
	for _, k := range t.Kolom {
		e, jml := ekspresiTulis(k, n)
		set = append(set, k.Kolom+" = "+e)
		n += jml
		v, err := nilaiTulis(k, h.Ambil(k.Properti))
		if err != nil {
			return err
		}
		args = append(args, v...)
	}
	q := fmt.Sprintf(`UPDATE %s SET %s WHERE ID = :%d`, tabel, strings.Join(set, ", "), n)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, append(args, id)...)
	if err != nil {
		return fmt.Errorf("repository: menyimpan induk klaim: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "penyimpanan induk klaim")
}

// jagaKomite - adjustment ber-KOMITE_ID tersimpan wajib masih ada di halaman (penghapusan objek / item induknya
// meng-CASCADE baris itu).
func (g *Gudang) jagaKomite(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error {
	tabel, err := g.db.Qualify(models.TabelAdjustment.Nama)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`SELECT ID FROM %s WHERE CLAIM_ID = :1 AND KOMITE_ID IS NOT NULL`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	rows, err := tx.QueryContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("repository: membaca adjustment berkomite: %w", err)
	}
	defer func() { _ = rows.Close() }()
	ada := map[string]bool{}
	for o := range h.AmbilDaftar(models.DaftarObjek) {
		for i := range h.AmbilDaftar(models.DaftarItem(o + 1)) {
			for _, a := range h.AmbilDaftar(models.DaftarAdj(o+1, i+1)) {
				ada[a[models.PropID]] = true
			}
		}
	}
	for rows.Next() {
		var rid sql.NullString
		if err := rows.Scan(&rid); err != nil {
			return fmt.Errorf("repository: memindai adjustment berkomite: %w", err)
		}
		if !ada[rid.String] {
			return ErrAdjustmentBerkomite
		}
	}
	return rows.Err()
}

// kosongkanTulisUlang menghapus seluruh baris simpul tak-stabil ber-CLAIM_ID milik klaim (sekali per tabel).
func (g *Gudang) kosongkanTulisUlang(ctx context.Context, tx *db.Tx, id string) error {
	sudah := map[string]bool{}
	for _, t := range models.SemuaTabel() {
		if t.Stabil || t == &models.TabelHeaderKlaim || (!t.DenganKlaim && t.KolomInduk != "CLAIM_ID") || sudah[t.Nama] {
			continue
		}
		sudah[t.Nama] = true
		tabel, err := g.db.Qualify(t.Nama)
		if err != nil {
			return err
		}
		q := fmt.Sprintf(`DELETE FROM %s WHERE CLAIM_ID = :1`, tabel)
		if err := db.PeriksaSQL(q); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, q, id); err != nil {
			return fmt.Errorf("repository: mengosongkan %s: %w", t.Nama, err)
		}
	}
	return nil
}

// tabelGeser - tabel stabil (baris disunting di tempat, UNIQUE (CLAIM_ID, NOURUT)): T_CLAIM_OBJECT, T_CLAIM_OBJECT_ITEM,
// T_CLAIM_ADJUSTMENT. Tabel lain dikosongkan lalu ditulis ulang (`kosongkanTulisUlang`).
func tabelGeser() []string {
	var out []string
	for _, t := range models.SemuaTabel() {
		if t.Stabil {
			out = append(out, t.Nama)
		}
	}
	return out
}

func sqlGeserNourut(tabel string) string {
	return fmt.Sprintf(`UPDATE %s SET NOURUT = -NOURUT WHERE CLAIM_ID = :1 AND NOURUT > 0`, tabel)
}

// geserNourut - fase pertama NOURUT dua fase setiap tabel stabil: NOURUT tersimpan dibalik negatif sebelum baris
// ditulis menurut urutan halaman. Tanpa itu item baru di objek pertama (NOURUT se-klaim) bertabrakan dengan item objek
// sesudahnya yang belum digeser (temuan review 10-10-2026).
func (g *Gudang) geserNourut(ctx context.Context, tx *db.Tx, id string) error {
	for _, nama := range tabelGeser() {
		tabel, err := g.db.Qualify(nama)
		if err != nil {
			return err
		}
		q := sqlGeserNourut(tabel)
		if err := db.PeriksaSQL(q); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, q, id); err != nil {
			return fmt.Errorf("repository: menggeser urutan %s: %w", nama, err)
		}
	}
	return nil
}

// kolomSisip - kolom INSERT simpul (ID, induk, [CLAIM_ID], NOURUT, [JENIS], katalog) beserta ekspresi penampungnya.
func kolomSisip(t *models.Tabel) (string, string) {
	kol := []string{"ID", t.KolomInduk}
	nilai := []string{":1", ":2"}
	n := 3
	if t.DenganKlaim && t.KolomInduk != "CLAIM_ID" {
		kol, nilai = append(kol, "CLAIM_ID"), append(nilai, fmt.Sprintf(":%d", n))
		n++
	}
	kol, nilai = append(kol, "NOURUT"), append(nilai, fmt.Sprintf(":%d", n))
	n++
	if t.Jenis != "" {
		kol, nilai = append(kol, "JENIS"), append(nilai, fmt.Sprintf(":%d", n))
		n++
	}
	for _, k := range t.Kolom {
		e, jml := ekspresiTulis(k, n)
		kol, nilai = append(kol, k.Kolom), append(nilai, e)
		n += jml
	}
	return strings.Join(kol, ", "), strings.Join(nilai, ", ")
}

// sqlSisipSimpul - INSERT satu baris simpul.
func sqlSisipSimpul(tabel string, t *models.Tabel) string {
	kol, nilai := kolomSisip(t)
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, tabel, kol, nilai)
}

// sqlUbahSimpul - UPDATE satu baris simpul stabil (kolom katalog, induk, NOURUT) menurut ID + CLAIM_ID. ⛔ Urutan
// penampung di teks = urutan argumen (driver mengikat menurut POSISI kemunculan).
func sqlUbahSimpul(tabel string, t *models.Tabel) string {
	var set []string
	n := 1
	for _, k := range t.Kolom {
		if k.MilikKomite {
			continue
		}
		e, jml := ekspresiTulis(k, n)
		set = append(set, k.Kolom+" = "+e)
		n += jml
	}
	induk := ""
	if t.KolomInduk != "CLAIM_ID" {
		induk = fmt.Sprintf(", %s = :%d", t.KolomInduk, n)
		n++
	}
	return fmt.Sprintf(`UPDATE %s SET %s%s, NOURUT = :%d WHERE ID = :%d AND CLAIM_ID = :%d`, tabel,
		strings.Join(set, ", "), induk, n, n+1, n+2)
}

// argsKatalog - nilai kolom katalog urut katalog; `ubah` = untuk sqlUbahSimpul (tanpa kolom milik komite).
func argsKatalog(t *models.Tabel, b models.Baris, ubah bool) ([]any, error) {
	var args []any
	for _, k := range t.Kolom {
		if ubah && k.MilikKomite {
			continue
		}
		v, err := nilaiTulis(k, b[k.Properti])
		if err != nil {
			return nil, err
		}
		args = append(args, v...)
	}
	return args, nil
}

// idTersimpan - ID baris simpul stabil milik klaim.
func (g *Gudang) idTersimpan(ctx context.Context, tx *db.Tx, tabel, id string) (map[string]bool, error) {
	q := fmt.Sprintf(`SELECT ID FROM %s WHERE CLAIM_ID = :1`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := tx.QueryContext(ctx, q, id)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca ID tersimpan: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := map[string]bool{}
	for rows.Next() {
		var rid sql.NullString
		if err := rows.Scan(&rid); err != nil {
			return nil, fmt.Errorf("repository: memindai ID tersimpan: %w", err)
		}
		out[rid.String] = true
	}
	return out, rows.Err()
}

// simpanSimpul menulis baris daftar `daftar` simpul `t` (induk langsung `indukID`) lalu anak-anaknya.
func (g *Gudang) simpanSimpul(ctx context.Context, tx *db.Tx, klaim string, t *models.Tabel, indukID, daftar string,
	h *models.Halaman, urut map[string]int) error {
	tabel, err := g.db.Qualify(t.Nama)
	if err != nil {
		return err
	}
	rows := h.AmbilDaftar(daftar)
	qSisip := sqlSisipSimpul(tabel, t)
	if err := db.PeriksaSQL(qSisip); err != nil {
		return err
	}
	var tersimpan map[string]bool
	tetap := map[string]bool{}
	if t.Stabil {
		if tersimpan, err = g.idTersimpan(ctx, tx, tabel, klaim); err != nil {
			return err
		}
		for _, b := range rows {
			if rid := b[models.PropID]; rid != "" && tersimpan[rid] {
				tetap[rid] = true
			}
		}
		if err := g.hapusYatim(ctx, tx, t, tabel, klaim, indukID, tetap); err != nil {
			return err
		}
	} else if !t.DenganKlaim && t.KolomInduk != "CLAIM_ID" { // cucu adjustment: per induk
		q := fmt.Sprintf(`DELETE FROM %s WHERE %s = :1`, tabel, t.KolomInduk)
		if err := db.PeriksaSQL(q); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, q, indukID); err != nil {
			return fmt.Errorf("repository: mengosongkan %s: %w", t.Nama, err)
		}
	}
	qUbah := ""
	if t.Stabil {
		qUbah = sqlUbahSimpul(tabel, t)
		if err := db.PeriksaSQL(qUbah); err != nil {
			return err
		}
	}
	for i, b := range rows {
		nourut := i + 1
		if t.DenganKlaim || t.KolomInduk == "CLAIM_ID" {
			urut[t.Nama]++
			nourut = urut[t.Nama]
		}
		rid := b[models.PropID]
		ubah := t.Stabil && rid != "" && tetap[rid]
		v, err := argsKatalog(t, b, ubah)
		if err != nil {
			return err
		}
		if ubah {
			args := v
			if t.KolomInduk != "CLAIM_ID" {
				args = append(args, indukID)
			}
			hasil, err := tx.ExecContext(ctx, qUbah, append(args, nourut, rid, klaim)...)
			if err != nil {
				return fmt.Errorf("repository: menyunting %s: %w", t.Nama, err)
			}
			if err := db.PastikanSatuBaris(hasil, "penyuntingan "+t.Nama); err != nil {
				return err
			}
		} else {
			if rid, err = g.idBerikut(ctx, tx); err != nil {
				return err
			}
			args := []any{rid, indukID}
			if t.DenganKlaim && t.KolomInduk != "CLAIM_ID" {
				args = append(args, klaim)
			}
			args = append(args, nourut)
			if t.Jenis != "" {
				args = append(args, t.Jenis)
			}
			hasil, err := tx.ExecContext(ctx, qSisip, append(args, v...)...)
			if err != nil {
				return fmt.Errorf("repository: menulis %s: %w", t.Nama, err)
			}
			if err := db.PastikanSatuBaris(hasil, "penulisan "+t.Nama); err != nil {
				return err
			}
			if t.Stabil {
				b[models.PropID] = rid
			}
		}
		for _, a := range t.Anak {
			if err := g.simpanSimpul(ctx, tx, klaim, a, rid, models.JalurAnak(daftar, i+1, a.Daftar), h, urut); err != nil {
				return err
			}
		}
	}
	return nil
}

// hapusYatim - baris stabil induk `indukID` yang tidak lagi ada di halaman dihapus (anak ikut CASCADE).
func (g *Gudang) hapusYatim(ctx context.Context, tx *db.Tx, t *models.Tabel, tabel, klaim, indukID string,
	tetap map[string]bool) error {
	q := fmt.Sprintf(`SELECT ID FROM %s WHERE CLAIM_ID = :1 AND %s = :2`, tabel, t.KolomInduk)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	rows, err := tx.QueryContext(ctx, q, klaim, indukID)
	if err != nil {
		return fmt.Errorf("repository: membaca baris %s: %w", t.Nama, err)
	}
	var buang []string
	for rows.Next() {
		var rid sql.NullString
		if err := rows.Scan(&rid); err != nil {
			_ = rows.Close()
			return fmt.Errorf("repository: memindai baris %s: %w", t.Nama, err)
		}
		if !tetap[rid.String] {
			buang = append(buang, rid.String)
		}
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return err
	}
	_ = rows.Close()
	qHapus := fmt.Sprintf(`DELETE FROM %s WHERE ID = :1 AND CLAIM_ID = :2`, tabel)
	if err := db.PeriksaSQL(qHapus); err != nil {
		return err
	}
	for _, rid := range buang {
		if _, err := tx.ExecContext(ctx, qHapus, rid, klaim); err != nil {
			return fmt.Errorf("repository: menghapus baris %s: %w", t.Nama, err)
		}
	}
	return nil
}

// ---------------------------------------------------------------- kronologi (T_VIEW_SUGGEST)

func sqlNomorKronologi(tabel string) string {
	return fmt.Sprintf(`SELECT NVL(MAX(NO), 0) + 1 FROM %s WHERE CLAIM_ID = :1`, tabel)
}

// sqlSisipKronologi - satu baris kronologi: DATE_SUGGEST = ASMDateTimeChronology, PIC_SUGGEST = ASMUser (akun),
// IS_CEDING_CONFIRM = ASMUserID (jabatan, kolom "User"), COMMENT_SUGGEST = pyNote, INITIAL_SUGGEST = ASMNoteType
// (kolom "Status") - pemetaan OQ-CFI-19.
func sqlSisipKronologi(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, CLAIM_ID, NO, DATE_SUGGEST, PIC_SUGGEST, IS_CEDING_CONFIRM, COMMENT_SUGGEST,
		INITIAL_SUGGEST) VALUES (:1, :2, :3, TO_DATE(:4, '%s'), :5, :6, :7, :8)`, tabel, fmtTanggal)
}

func sqlBacaKronologi(tabel string) string {
	return fmt.Sprintf(`SELECT TO_CHAR(NO), %s, PIC_SUGGEST, IS_CEDING_CONFIRM, COMMENT_SUGGEST, INITIAL_SUGGEST
		  FROM %s WHERE CLAIM_ID = :1 ORDER BY NO`, fmt.Sprintf(db.FmtTanggalOracle, "DATE_SUGGEST"), tabel)
}

// PengenalKronologi - `T_VIEW_SUGGEST.ID` kronologi klaim: 32 heksa deterministik dari (klaim, NO), pola
// `PengenalRiwayat` Claim Non Prop (disalin).
func PengenalKronologi(klaimID string, no int) string {
	sum := md5.Sum([]byte(klaimID + "\x00suggest\x00" + strconv.Itoa(no)))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

func (g *Gudang) sisipKronologi(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error {
	var baru []models.Baris
	for _, b := range h.AmbilDaftar(models.DaftarKronologi) {
		if b[models.PropRiwayatBaru] == "1" {
			baru = append(baru, b)
		}
	}
	if len(baru) == 0 {
		return nil
	}
	tabel, err := g.db.Qualify("T_VIEW_SUGGEST")
	if err != nil {
		return err
	}
	qNo, q := sqlNomorKronologi(tabel), sqlSisipKronologi(tabel)
	for _, x := range []string{qNo, q} {
		if err := db.PeriksaSQL(x); err != nil {
			return err
		}
	}
	var noTeks sql.NullString
	if err := tx.QueryRowContext(ctx, qNo, id).Scan(&noTeks); err != nil {
		return fmt.Errorf("repository: membaca nomor kronologi klaim: %w", err)
	}
	no, err := strconv.Atoi(strings.TrimSpace(noTeks.String))
	if err != nil {
		return fmt.Errorf("repository: nomor kronologi klaim %q: %w", noTeks.String, err)
	}
	for _, b := range baru {
		var tgl any
		if t, ok := models.UraiTanggal(b["ASMDateTimeChronology"]); ok {
			tgl = t.Format(utils.TanggalWaktu)
		}
		hasil, err := tx.ExecContext(ctx, q, PengenalKronologi(id, no), id, no, tgl, teksAtauNil(potong(b["ASMUser"], 255)),
			teksAtauNil(potong(b["ASMUserID"], 255)), teksAtauNil(potong(b["pyNote"], 255)),
			teksAtauNil(potong(b["ASMNoteType"], 255)))
		if err != nil {
			return fmt.Errorf("repository: menyisipkan kronologi klaim: %w", err)
		}
		if err := db.PastikanSatuBaris(hasil, "kronologi klaim"); err != nil {
			return err
		}
		delete(b, models.PropRiwayatBaru)
		no++
	}
	return nil
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

type barisSimpul struct {
	id, induk string
	b         models.Baris
}

// BacaHalaman membaca halaman klaim `id` (tx boleh nil): kepala, pohon (satu kueri per simpul, dikelompokkan menurut
// induk), retro tingkat klaim, kronologi.
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
	simpul := map[*models.Tabel]map[string][]barisSimpul{}
	for _, s := range models.SemuaTabel()[1:] {
		rows, err := g.bacaSimpul(qn, s, id)
		if err != nil {
			return nil, err
		}
		per := map[string][]barisSimpul{}
		for _, r := range rows {
			per[r.induk] = append(per[r.induk], r)
		}
		simpul[s] = per
	}
	var pasang func(t *models.Tabel, induk, daftar string)
	pasang = func(t *models.Tabel, induk, daftar string) {
		rows := simpul[t][induk]
		var bs []models.Baris
		for _, r := range rows {
			bs = append(bs, r.b)
		}
		h.SetelDaftar(daftar, bs)
		for i, r := range rows {
			for _, a := range t.Anak {
				pasang(a, r.id, models.JalurAnak(daftar, i+1, a.Daftar))
			}
		}
	}
	pasang(&models.TabelObjek, id, models.DaftarObjek)
	pasang(&models.TabelRetroTreaty, id, models.DaftarFacRetroTreaty)
	kr, err := g.bacaKronologi(qn, id)
	if err != nil {
		return nil, err
	}
	h.SetelDaftar(models.DaftarKronologi, kr)
	return h, nil
}

func (g *Gudang) bacaSimpul(qn func(string, ...any) (*sql.Rows, error), t *models.Tabel, id string) ([]barisSimpul, error) {
	tabel, err := g.db.Qualify(t.Nama)
	if err != nil {
		return nil, err
	}
	tabelAdj, err := g.db.Qualify(models.TabelAdjustment.Nama)
	if err != nil {
		return nil, err
	}
	q, jenis := sqlBacaSimpul(t, tabel, tabelAdj)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	args := []any{id}
	for _, j := range jenis {
		args = append(args, j)
	}
	rows, err := qn(q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s: %w", t.Nama, err)
	}
	defer func() { _ = rows.Close() }()
	awal := 2
	if t == &models.TabelAdjustment {
		awal = 3
	}
	var out []barisSimpul
	for rows.Next() {
		lebar := awal + len(t.Kolom)
		nilai := make([]sql.NullString, lebar)
		tujuan := make([]any, lebar)
		for i := range nilai {
			tujuan[i] = &nilai[i]
		}
		if err := rows.Scan(tujuan...); err != nil {
			return nil, fmt.Errorf("repository: memindai %s: %w", t.Nama, err)
		}
		b := models.Baris{}
		if t.Stabil {
			b[models.PropID] = nilai[0].String
		}
		if awal == 3 && nilai[2].String != "" {
			b[models.PropKomiteID] = nilai[2].String
		}
		for i, k := range t.Kolom {
			if v := nilaiBaca(k, nilai[awal+i]); v != "" {
				b[k.Properti] = v
			}
		}
		out = append(out, barisSimpul{id: nilai[0].String, induk: nilai[1].String, b: b})
	}
	return out, rows.Err()
}

func (g *Gudang) bacaKronologi(qn func(string, ...any) (*sql.Rows, error), id string) ([]models.Baris, error) {
	tabel, err := g.db.Qualify("T_VIEW_SUGGEST")
	if err != nil {
		return nil, err
	}
	q := sqlBacaKronologi(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := qn(q, id)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca kronologi klaim: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []models.Baris
	for rows.Next() {
		var n [6]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5]); err != nil {
			return nil, fmt.Errorf("repository: memindai kronologi klaim: %w", err)
		}
		out = append(out, models.Baris{"No": n[0].String, "ASMDateTimeChronology": n[1].String, "ASMUser": n[2].String,
			"ASMUserID": n[3].String, "pyNote": n[4].String, "ASMNoteType": n[5].String})
	}
	return out, rows.Err()
}
