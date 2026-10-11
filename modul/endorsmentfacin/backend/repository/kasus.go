package repository

// Kasus Endorsment Fac In - irisan I3. Tiga tulisan dalam SATU transaksi milik pemanggil (services):
//
//  1. nomor dari SEQ_WORK_POLIS_EDM (migrasi 260), dirakit `models.RakitIDKasusEDM`;
//  2. baris T_WORK_POLIS - tabel milik premiumlistlife (050/057/059/063), dipakai bersama (K-064), LINI 'FACEDM'
//     (keputusan work owner 11-10-2026; persetujuan tulis T_WORK_POLIS 11-10-2026);
//  3. baris T_EDM_FACIN_KASUS (migrasi 262) - berbagi PK dengan T_WORK_POLIS, tanpa FK.
//
// Pembuat = PENGENAL AKUN pelaku di CREATE_OP dan CREATE_OP_NAME, bukan nama orang (ADR-U-0030; pola nbfacin
// casenb.go). Seluruh nilai diikat.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/endorsmentfacin/backend/models"
)

// Nama objek I3.
const (
	// TabelWorkPolis - tabel kerja polis lintas lini (premiumlistlife 050), ejaan migrasinya.
	TabelWorkPolis = "T_WORK_POLIS"
	// TabelKasus - kepala kasus milik modul ini (262).
	TabelKasus = "T_EDM_FACIN_KASUS"
	// SequenceKasus - penghitung nomor kasus (260).
	SequenceKasus = "SEQ_WORK_POLIS_EDM"
)

// DaftarTabelDitulis - tabel yang ditulis modul ini.
var DaftarTabelDitulis = []string{TabelWorkPolis, TabelKasus}

type penjalan interface {
	QueryRowContext(ctx context.Context, q string, args ...any) *sql.Row
	ExecContext(ctx context.Context, q string, args ...any) (sql.Result, error)
}

func (g *Gudang) dari(tx *db.Tx) penjalan {
	if tx.Terisi() {
		return tx
	}
	return g.db
}

// NomorKasusBerikut - NEXTVAL SEQ_WORK_POLIS_EDM → "EDM-<5 digit>". Wajib di dalam transaksi.
func (g *Gudang) NomorKasusBerikut(ctx context.Context, tx *db.Tx) (string, error) {
	if !tx.Terisi() {
		return "", errors.New("repository: nomor kasus EDM wajib diambil di dalam transaksi")
	}
	s, err := g.db.NomorBerikut(ctx, tx, SequenceKasus)
	if err != nil {
		return "", bungkus(err, "mengambil nomor kasus ("+SequenceKasus+")")
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return "", fmt.Errorf("repository: nomor %s bukan bilangan bulat: %q", SequenceKasus, s)
	}
	return models.RakitIDKasusEDM(n), nil
}

// SqlSisipWorkPolis - baris T_WORK_POLIS berkeadaan awal kasus EDM.
func SqlSisipWorkPolis(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, LINI, POSITION, FLAG_ONGOING_POLICY, STATUS_WORK, CREATE_OP, CREATE_OP_NAME,
	  TGL_CREATE, TGL_UPDATE) VALUES (:1, :2, :3, :4, :5, :6, :7, SYSDATE, SYSDATE)`, t)
}

// SisipWorkPolis - satu baris T_WORK_POLIS.
func (g *Gudang) SisipWorkPolis(ctx context.Context, tx *db.Tx, id, akun string) error {
	t, err := g.db.Qualify(TabelWorkPolis)
	if err != nil {
		return err
	}
	return g.sisipSatu(ctx, tx, SqlSisipWorkPolis(t), "menyisipkan "+TabelWorkPolis,
		id, models.LiniEDMFacIn, models.PosisiAwalEDM, models.FlagAwalEDM, models.StatusAwalEDM, akun, akun)
}

// KolomKasus - urutan kolom T_EDM_FACIN_KASUS = urutan bind ArgSisipKasus.
var KolomKasus = []string{"ID", "OLD_POLICY_NO", "EDM_TYPE", "EDM_TYPE_NEW", "DETAIL_TYPE", "ENDORSEMENT_DATE",
	"ENDORSEMENT_REASON", "ENDORSEMENT_SOURCE", "ENDORSEMENT_SOURCE_NOTE", "ENDORSEMENT_STATUS", "ADMIN_CHARGE", "SURVEY",
	"INTERNAL_RETRO", "NB_STATUS", "IS_CEDING_CONFIRM"}

// SqlSisipKasus - baris kepala kasus; ENDORSEMENT_DATE = TO_DATE tanggal kalender.
func SqlSisipKasus(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, OLD_POLICY_NO, EDM_TYPE, EDM_TYPE_NEW, DETAIL_TYPE, ENDORSEMENT_DATE,
	  ENDORSEMENT_REASON, ENDORSEMENT_SOURCE, ENDORSEMENT_SOURCE_NOTE, ENDORSEMENT_STATUS, ADMIN_CHARGE, SURVEY,
	  INTERNAL_RETRO, NB_STATUS, IS_CEDING_CONFIRM)
	  VALUES (:1, :2, :3, :4, :5, TO_DATE(:6, 'YYYY-MM-DD'), :7, :8, :9, :10, :11, :12, :13, :14, :15)`, t)
}

// ArgSisipKasus - nilai bind, urut KolomKasus. Teks kosong = NULL.
func ArgSisipKasus(k models.KepalaKasus) []any {
	n := db.KosongJadiNil
	tanggal := any(nil)
	if !k.EndorsementDate.IsZero() {
		tanggal = k.EndorsementDate.Format("2006-01-02")
	}
	return []any{k.ID, n(k.OldPolicyNo), n(k.EdmType), n(k.EdmTypeNew), n(k.DetailType), tanggal,
		n(k.EndorsementReason), n(k.EndorsementSource), n(k.EndorsementSourceNote), n(k.EndorsementStatus),
		n(k.AdminCharge), n(k.Survey), n(k.InternalRetro), n(k.NBStatus), n(k.IsCedingConfirm)}
}

// SisipKasus - satu baris T_EDM_FACIN_KASUS.
func (g *Gudang) SisipKasus(ctx context.Context, tx *db.Tx, k models.KepalaKasus) error {
	t, err := g.db.Qualify(TabelKasus)
	if err != nil {
		return err
	}
	return g.sisipSatu(ctx, tx, SqlSisipKasus(t), "menyisipkan "+TabelKasus, ArgSisipKasus(k)...)
}

func (g *Gudang) sisipSatu(ctx context.Context, tx *db.Tx, q, apa string, args ...any) error {
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	h, err := g.dari(tx).ExecContext(ctx, q, args...)
	if err != nil {
		return bungkus(err, apa)
	}
	return db.PastikanSatuBaris(h, apa)
}

// SqlEDMBelumSelesai - gerbang 2: kasus EDM sistem baru ber-polis lama sama yang STATUS_WORK-nya belum "Resolved-…".
// Pengganti `GetListEdm` (report Pega tanpa body di korpus) - kasus era Pega TIDAK terbaca (keputusan work owner
// 11-10-2026, W-EDM-12).
func SqlEDMBelumSelesai(work, kasus string) string {
	return fmt.Sprintf(`SELECT w.ID FROM %s w JOIN %s k ON k.ID = w.ID
	  WHERE w.LINI = :1 AND k.OLD_POLICY_NO = :2 AND (w.STATUS_WORK IS NULL OR w.STATUS_WORK NOT LIKE :3)
	  ORDER BY w.ID FETCH FIRST 1 ROWS ONLY`, work, kasus)
}

// EDMBelumSelesai - ID kasus EDM terbuka pertama untuk polis itu; ("", false) bila tidak ada. Tabel 262 belum ada →
// ErrBelumAda.
func (g *Gudang) EDMBelumSelesai(ctx context.Context, polis string) (string, bool, error) {
	w, err := g.db.Qualify(TabelWorkPolis)
	if err != nil {
		return "", false, err
	}
	k, err := g.db.Qualify(TabelKasus)
	if err != nil {
		return "", false, err
	}
	return g.satuTeks(ctx, "membaca kasus EDM terbuka", SqlEDMBelumSelesai(w, k),
		models.LiniEDMFacIn, polis, models.PolaStatusSelesai)
}
