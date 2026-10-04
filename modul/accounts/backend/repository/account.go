// Package repository adalah SATU-SATUNYA lapisan modul Accounts yang berbicara ke Oracle.
//
// Tabel warisan `POOLDATA.T_M_ACCOUNT` (Pega SFAGIS Account) dibaca dan ditulis apa adanya, ditambah tiga kolom
// migrasi 840 (`CREATEDATE`, `CREATEOP`, `DESCRIPTION`) dan kunci utama 841. `CLIENT` (organisasi, milik Company
// Detail) dan `BUSINESSGROUP` hanya DIBACA. Nomor akun baru dari `SEQ_T_M_ACCOUNT` (migrasi 842).
// ⛔ Setiap tabel lewat `Qualify` (ADR-U-0033), nilai lewat bind, nol COMMIT di SQL (ADR-U-0029): transaksinya
// milik services. Nol DELETE (keputusan work owner 04-10-2026: "Delete tidak").
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/accounts/backend/models"
)

// Objek Oracle modul ini.
const (
	TabelAccount = "T_M_ACCOUNT"
	// TabelClient - organisasi (`FLAG` `Org`) untuk "Insured Name"; dibaca saja.
	TabelClient = "CLIENT"
	// TabelBusinessGroup - pilihan "Group Business" (`ID`, `NOTE`); dibaca saja.
	TabelBusinessGroup = "BUSINESSGROUP"
	// SeqAccount - nomor ACC akun baru (migrasi 842).
	SeqAccount = "SEQ_T_M_ACCOUNT"
	// FlagOrg - `CLIENT.FLAG` organisasi (katalog DEV 04-10-2026: 18.655 dari 18.658 INSUREDID).
	FlagOrg = "Org"
)

// DaftarTabelDitulis dan DaftarTabelDibacaSaja - penjaga modul.
var (
	DaftarTabelDitulis    = []string{TabelAccount}
	DaftarTabelDibacaSaja = []string{TabelClient, TabelBusinessGroup}
)

// ErrTidakAda - baris yang dicari tidak ada.
var ErrTidakAda = errors.New("repository: row not found")

// Gudang membaca dan menulis Oracle.
type Gudang struct{ db *db.DB }

// Baru menyusun Gudang.
func Baru(d *db.DB) *Gudang { return &Gudang{db: d} }

// penjalan - *db.DB dan *db.Tx.
type penjalan interface {
	QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, q string, args ...any) *sql.Row
	ExecContext(ctx context.Context, q string, args ...any) (sql.Result, error)
}

func (g *Gudang) dari(tx *db.Tx) penjalan {
	if tx.Terisi() {
		return tx
	}
	return g.db
}

func (g *Gudang) nama(objek string) (string, error) { return g.db.Qualify(objek) }

// siapkan - nama objek ter-Qualify lalu SQL yang sudah diperiksa `db.PeriksaSQL`.
func (g *Gudang) siapkan(objek string, susun func(string) string) (string, error) {
	t, err := g.nama(objek)
	if err != nil {
		return "", err
	}
	q := susun(t)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	return q, nil
}

// kolomAccount - urutan kolom setiap SELECT baris akun (`pindaiAccount`).
const kolomAccount = `ID, GROUPBUSINESSID, GROUPBUSINESS, INSUREDID, INSUREDNAME, DESCRIPTION, CREATEOP,
	TO_CHAR(CREATEDATE, 'YYYY-MM-DD HH24:MI')`

// saringCari - :1 kosong (NULL) = tanpa saringan; :2-:5 pola yang sama untuk Insured Name, Group Business, ID,
// dan Org ID (tanpa beda huruf).
const saringCari = `(:1 IS NULL OR UPPER(INSUREDNAME) LIKE :2 ESCAPE '\' OR UPPER(GROUPBUSINESS) LIKE :3 ESCAPE '\'
	OR UPPER(ID) LIKE :4 ESCAPE '\' OR UPPER(INSUREDID) LIKE :5 ESCAPE '\')`

// sqlDaftar - satu halaman daftar, terbaru (nomor ACC terbesar) lebih dulu.
func sqlDaftar(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE %s
	  ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '[0-9]+$')) DESC NULLS LAST, ID DESC
	  OFFSET :6 ROWS FETCH NEXT :7 ROWS ONLY`, kolomAccount, t, saringCari)
}

func sqlHitung(t string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, t, saringCari)
}

func sqlAmbil(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, kolomAccount, t)
}

func sqlAdaID(t string) string { return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ID = :1`, t) }

func sqlNomor(seq string) string { return fmt.Sprintf(`SELECT %s.NEXTVAL FROM DUAL`, seq) }

// sqlSisip - CREATEDATE = SYSDATE; CREATEOP = akun pelaku.
func sqlSisip(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, GROUPBUSINESSID, GROUPBUSINESS, INSUREDID, INSUREDNAME, DESCRIPTION,
	    CREATEOP, CREATEDATE)
	  VALUES (:1, :2, :3, :4, :5, :6, :7, SYSDATE)`, t)
}

// sqlPasanganLain - akun yang sudah ada dengan Insured + Group Business yang sama.
func sqlPasanganLain(t string) string {
	return fmt.Sprintf(`SELECT ID FROM %s WHERE INSUREDID = :1 AND GROUPBUSINESSID = :2
	  ORDER BY ID FETCH FIRST 1 ROWS ONLY`, t)
}

// sqlCariOrganisasi - organisasi untuk "Insured Name": nama atau ORG-n mengandung kata cari.
func sqlCariOrganisasi(t string) string {
	return fmt.Sprintf(`SELECT ID, IDVIEW, NAME FROM %s
	  WHERE FLAG = :1 AND NAME IS NOT NULL
	    AND (:2 IS NULL OR UPPER(NAME) LIKE :3 ESCAPE '\' OR UPPER(IDVIEW) LIKE :4 ESCAPE '\')
	  ORDER BY UPPER(NAME), ID FETCH FIRST :5 ROWS ONLY`, t)
}

func sqlAmbilOrganisasi(t string) string {
	return fmt.Sprintf(`SELECT ID, IDVIEW, NAME FROM %s WHERE ID = :1 AND FLAG = :2`, t)
}

func sqlDaftarGroupBusiness(t string) string {
	return fmt.Sprintf(`SELECT ID, NOTE FROM %s WHERE NOTE IS NOT NULL ORDER BY UPPER(NOTE), ID`, t)
}

func sqlAmbilGroupBusiness(t string) string {
	return fmt.Sprintf(`SELECT ID, NOTE FROM %s WHERE ID = :1`, t)
}

// pola - `%KATA%` tanpa beda huruf, dengan `%`, `_`, dan `\` di-escape.
func pola(kata string) string {
	k := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(strings.ToUpper(strings.TrimSpace(kata)))
	return "%" + k + "%"
}

// argCari - :1-:5 saringCari; kata kosong = NULL.
func argCari(kata string) []any {
	if strings.TrimSpace(kata) == "" {
		return []any{nil, nil, nil, nil, nil}
	}
	p := pola(kata)
	return []any{p, p, p, p, p}
}

func pindaiAccount(r interface{ Scan(...any) error }) (models.Account, error) {
	var id, gbID, gb, insID, insNama, desk, op, tgl sql.NullString
	if err := r.Scan(&id, &gbID, &gb, &insID, &insNama, &desk, &op, &tgl); err != nil {
		return models.Account{}, err
	}
	return models.Account{ID: id.String, IDView: models.BagianAkhir(id.String, models.AwalanID),
		GroupBusinessID: gbID.String, GroupBusiness: gb.String, InsuredID: insID.String,
		OrgID: models.BagianAkhir(insID.String, models.AwalanOrg), InsuredName: insNama.String,
		Description: desk.String, CreateOp: op.String, CreateDate: tgl.String}, nil
}

// Daftar - satu halaman akun dan jumlah seluruh baris yang cocok.
func (g *Gudang) Daftar(ctx context.Context, s models.Saringan) ([]models.Account, int, error) {
	qHitung, err := g.siapkan(TabelAccount, sqlHitung)
	if err != nil {
		return nil, 0, err
	}
	var total int
	if err := g.db.QueryRowContext(ctx, qHitung, argCari(s.Cari)...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("repository: counting %s: %w", TabelAccount, err)
	}
	q, err := g.siapkan(TabelAccount, sqlDaftar)
	if err != nil {
		return nil, 0, err
	}
	args := append(argCari(s.Cari), (s.Halaman-1)*s.Ukuran, s.Ukuran)
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: reading %s: %w", TabelAccount, err)
	}
	defer func() { _ = rows.Close() }()
	out := []models.Account{}
	for rows.Next() {
		a, err := pindaiAccount(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("repository: reading %s: %w", TabelAccount, err)
		}
		out = append(out, a)
	}
	return out, total, rows.Err()
}

// Ambil - satu akun menurut ID penuh.
func (g *Gudang) Ambil(ctx context.Context, tx *db.Tx, id string) (models.Account, error) {
	q, err := g.siapkan(TabelAccount, sqlAmbil)
	if err != nil {
		return models.Account{}, err
	}
	a, err := pindaiAccount(g.dari(tx).QueryRowContext(ctx, q, id))
	if errors.Is(err, sql.ErrNoRows) {
		return models.Account{}, ErrTidakAda
	}
	if err != nil {
		return models.Account{}, fmt.Errorf("repository: reading %s: %w", TabelAccount, err)
	}
	return a, nil
}

// AdaID - ID itu sudah dipakai baris lain.
func (g *Gudang) AdaID(ctx context.Context, tx *db.Tx, id string) (bool, error) {
	q, err := g.siapkan(TabelAccount, sqlAdaID)
	if err != nil {
		return false, err
	}
	var n int
	if err := g.dari(tx).QueryRowContext(ctx, q, id).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: reading %s: %w", TabelAccount, err)
	}
	return n > 0, nil
}

// NomorBerikut - `SEQ_T_M_ACCOUNT.NEXTVAL`.
func (g *Gudang) NomorBerikut(ctx context.Context, tx *db.Tx) (int64, error) {
	seq, err := g.nama(SeqAccount)
	if err != nil {
		return 0, err
	}
	var n int64
	if err := g.dari(tx).QueryRowContext(ctx, sqlNomor(seq)).Scan(&n); err != nil {
		return 0, fmt.Errorf("repository: %s (migration 842): %w", SeqAccount, err)
	}
	return n, nil
}

// Sisip - satu baris baru; CREATEDATE = SYSDATE.
func (g *Gudang) Sisip(ctx context.Context, tx *db.Tx, a models.Account) error {
	q, err := g.siapkan(TabelAccount, sqlSisip)
	if err != nil {
		return err
	}
	_, err = g.dari(tx).ExecContext(ctx, q, a.ID, a.GroupBusinessID, a.GroupBusiness, a.InsuredID, a.InsuredName,
		db.KosongJadiNil(a.Description), a.CreateOp)
	if err != nil {
		return fmt.Errorf("repository: inserting %s: %w", TabelAccount, err)
	}
	return nil
}

// PasanganLain - ID akun yang sudah ada dengan Insured + Group Business yang sama; kosong bila tidak ada.
func (g *Gudang) PasanganLain(ctx context.Context, tx *db.Tx, insuredID, groupBusinessID string) (string, error) {
	q, err := g.siapkan(TabelAccount, sqlPasanganLain)
	if err != nil {
		return "", err
	}
	var id string
	err = g.dari(tx).QueryRowContext(ctx, q, insuredID, groupBusinessID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: reading %s: %w", TabelAccount, err)
	}
	return id, nil
}

// CariOrganisasi - paling banyak `batas` organisasi yang nama atau ORG-n-nya mengandung kata cari.
func (g *Gudang) CariOrganisasi(ctx context.Context, kata string, batas int) ([]models.Organisasi, error) {
	q, err := g.siapkan(TabelClient, sqlCariOrganisasi)
	if err != nil {
		return nil, err
	}
	var p any
	if strings.TrimSpace(kata) != "" {
		p = pola(kata)
	}
	rows, err := g.db.QueryContext(ctx, q, FlagOrg, p, p, p, batas)
	if err != nil {
		return nil, fmt.Errorf("repository: reading %s: %w", TabelClient, err)
	}
	defer func() { _ = rows.Close() }()
	out := []models.Organisasi{}
	for rows.Next() {
		var id, view, nama sql.NullString
		if err := rows.Scan(&id, &view, &nama); err != nil {
			return nil, fmt.Errorf("repository: reading %s: %w", TabelClient, err)
		}
		out = append(out, models.Organisasi{ID: id.String, IDView: view.String, Nama: nama.String})
	}
	return out, rows.Err()
}

// AmbilOrganisasi - satu organisasi menurut `CLIENT.ID`.
func (g *Gudang) AmbilOrganisasi(ctx context.Context, tx *db.Tx, id string) (models.Organisasi, error) {
	q, err := g.siapkan(TabelClient, sqlAmbilOrganisasi)
	if err != nil {
		return models.Organisasi{}, err
	}
	var cid, view, nama sql.NullString
	err = g.dari(tx).QueryRowContext(ctx, q, id, FlagOrg).Scan(&cid, &view, &nama)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Organisasi{}, ErrTidakAda
	}
	if err != nil {
		return models.Organisasi{}, fmt.Errorf("repository: reading %s: %w", TabelClient, err)
	}
	return models.Organisasi{ID: cid.String, IDView: view.String, Nama: nama.String}, nil
}

// DaftarGroupBusiness - seluruh `BUSINESSGROUP` ber-NOTE, urut NOTE.
func (g *Gudang) DaftarGroupBusiness(ctx context.Context) ([]models.GroupBusiness, error) {
	q, err := g.siapkan(TabelBusinessGroup, sqlDaftarGroupBusiness)
	if err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: reading %s: %w", TabelBusinessGroup, err)
	}
	defer func() { _ = rows.Close() }()
	out := []models.GroupBusiness{}
	for rows.Next() {
		var id, note sql.NullString
		if err := rows.Scan(&id, &note); err != nil {
			return nil, fmt.Errorf("repository: reading %s: %w", TabelBusinessGroup, err)
		}
		out = append(out, models.GroupBusiness{ID: id.String, Note: note.String})
	}
	return out, rows.Err()
}

// AmbilGroupBusiness - satu `BUSINESSGROUP` menurut ID.
func (g *Gudang) AmbilGroupBusiness(ctx context.Context, tx *db.Tx, id string) (models.GroupBusiness, error) {
	q, err := g.siapkan(TabelBusinessGroup, sqlAmbilGroupBusiness)
	if err != nil {
		return models.GroupBusiness{}, err
	}
	var gid, note sql.NullString
	err = g.dari(tx).QueryRowContext(ctx, q, id).Scan(&gid, &note)
	if errors.Is(err, sql.ErrNoRows) {
		return models.GroupBusiness{}, ErrTidakAda
	}
	if err != nil {
		return models.GroupBusiness{}, fmt.Errorf("repository: reading %s: %w", TabelBusinessGroup, err)
	}
	return models.GroupBusiness{ID: gid.String, Note: note.String}, nil
}
