// Package repository adalah SATU-SATUNYA lapisan modul Marketing Officer yang berbicara ke Oracle.
//
// Tabel warisan `MARKETINGOFFICER` ditulis dan dibaca apa adanya (keputusan work owner 03-10-2026: nol DDL, nol
// tabel baru). `M_LOGIN_GO` (milik inti) dan `BRANCH` hanya dibaca.
//
// ⛔ Prosedur `PEGA_MARKETINGOFFICER` TIDAK dipanggil: INSERT dan UPDATE di sini menulis kolom yang SAMA dengan
// prosedur itu. `ID` = `1` + 7 digit `CURRENCY_SEQ` seperti prosedur (`vID := 1||lpad(currency_seq.nextval,7,'0')`).
// ⛔ Setiap tabel lewat `Qualify` (ADR-U-0033), nilai lewat bind, nol COMMIT di SQL (ADR-U-0029): transaksinya
// milik services.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/marketingofficer/backend/models"
)

// Objek Oracle modul ini.
const (
	TabelMO = "MARKETINGOFFICER"
	// TabelLog - log perubahan warisan; trigger `TRG_MARKETINGOFFICER_LOG` menyisipkan baris LAMA setiap UPDATE.
	// Migrasi modul 760 menambah `LOG_TIME` (default SYSTIMESTAMP) dan `AKSES_LOGIN`.
	TabelLog    = "MARKETINGOFFICER_LOG"
	TabelLogin  = "M_LOGIN_GO"
	TabelCabang = "BRANCH"
	// SeqMO - sequence yang dipakai `PEGA_MARKETINGOFFICER` (bukan `MARKETINGOFFICER_SEQ`).
	SeqMO = "CURRENCY_SEQ"
)

// DaftarTabelDitulis dan DaftarTabelDibacaSaja - penjaga modul.
var (
	DaftarTabelDitulis    = []string{TabelMO, TabelLog}
	DaftarTabelDibacaSaja = []string{TabelLogin, TabelCabang}
)

var (
	// ErrTidakAda - baris yang dicari tidak ada.
	ErrTidakAda = errors.New("repository: row not found")
	// ErrNomorMelampaui - nomor CURRENCY_SEQ tidak muat di 7 digit; ID tidak boleh dipotong seperti LPAD Pega.
	ErrNomorMelampaui = errors.New("repository: CURRENCY_SEQ exceeds 7 digits; the marketing officer ID cannot be formed")
	// ErrLogBelumDimigrasi - kolom LOG_TIME/AKSES_LOGIN MARKETINGOFFICER_LOG belum ada: migrasi 760 belum dijalankan.
	ErrLogBelumDimigrasi = errors.New("repository: MARKETINGOFFICER_LOG has not been fixed yet - migration 760 has not run")
)

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

// kolomMO - urutan kolom setiap SELECT baris MO (`pindaiMO`).
const kolomMO = `ID, CLIENTID, CLIENTNAME, CLIENTID2, MOLEADER, MOSTATUS, BRANCHPARENT, BRANCHDETAILID,
	BRANCHDETAILNAME, TEAMGROUP, BRANCHSTATUS, AKSES_LOGIN, USERUPDATE, TO_CHAR(TANGGAL, 'YYYY-MM-DD HH24:MI')`

func sqlDaftarMO(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s ORDER BY MOSTATUS, UPPER(CLIENTNAME), ID`, kolomMO, t)
}

func sqlAmbilMO(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, kolomMO, t)
}

func sqlNomorMO(seq string) string { return fmt.Sprintf(`SELECT %s.NEXTVAL FROM DUAL`, seq) }

// sqlSisipMO - kolom dan urutan `INSERT` `PEGA_MARKETINGOFFICER`; `TANGGAL` = SYSDATE. BRANCHSTATUS tidak ditulis.
func sqlSisipMO(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, CLIENTID, BRANCHDETAILID, MOLEADER, MOSTATUS, CLIENTID2, TEAMGROUP,
	    CLIENTNAME, TANGGAL, USERUPDATE, BRANCHPARENT, BRANCHDETAILNAME, AKSES_LOGIN)
	  VALUES (:1, :2, :3, :4, :5, :6, :7, :8, SYSDATE, :9, :10, :11, :12)`, t)
}

// sqlPerbaruiMO - `UPDATE` `PEGA_MARKETINGOFFICER` TANPA `CLIENTID` dan `CLIENTNAME`: Marketing Code sudah
// tersalin ke data transaksi, dan nama adalah identitas baris (keputusan rancangan 03-10-2026). Trigger warisan
// `TRG_MARKETINGOFFICER_LOG` mencatat baris lama setiap UPDATE.
func sqlPerbaruiMO(t string) string {
	return fmt.Sprintf(`UPDATE %s SET BRANCHDETAILID = :1, MOLEADER = :2, MOSTATUS = :3, CLIENTID2 = :4,
	    TEAMGROUP = :5, TANGGAL = SYSDATE, USERUPDATE = :6, BRANCHPARENT = :7, BRANCHDETAILNAME = :8,
	    AKSES_LOGIN = :9
	  WHERE ID = :10`, t)
}

// sqlAktifDenganClientID dan sqlAktifDenganAkses - baris AKTIF lain dengan Marketing Code / akun yang sama.
// :3 = ID baris yang dikecualikan (NULL = tidak ada).
func sqlAktifDenganClientID(t string) string {
	return fmt.Sprintf(`SELECT ID FROM %s WHERE MOSTATUS = :1 AND CLIENTID = :2 AND ID <> NVL(:3, '-') ORDER BY ID`, t)
}

func sqlAktifDenganAkses(t string) string {
	return fmt.Sprintf(`SELECT ID FROM %s
	  WHERE MOSTATUS = :1 AND UPPER(AKSES_LOGIN) = UPPER(:2) AND ID <> NVL(:3, '-') ORDER BY ID`, t)
}

// sqlClientIDDariAkses - Marketing Code yang sudah dipakai akun itu di baris MO lama: baris aktif lebih dulu,
// lalu yang terbaru.
func sqlClientIDDariAkses(t string) string {
	return fmt.Sprintf(`SELECT CLIENTID FROM %s
	  WHERE UPPER(AKSES_LOGIN) = UPPER(:1) AND CLIENTID IS NOT NULL
	  ORDER BY MOSTATUS, TANGGAL DESC NULLS LAST, ID DESC`, t)
}

func sqlDaftarAkun(t string) string {
	return fmt.Sprintf(`SELECT LOGIN_ID, NAME, CONTACT_ID, EMAIL, IS_ACTIVE FROM %s ORDER BY UPPER(NAME), LOGIN_ID`, t)
}

func sqlAmbilAkun(t string) string {
	return fmt.Sprintf(`SELECT LOGIN_ID, NAME, CONTACT_ID, EMAIL, IS_ACTIVE FROM %s WHERE LOGIN_ID = :1`, t)
}

// sqlDaftarCabang - `BrowseBranchDetail_RD` filter `.Status = 1`.
func sqlDaftarCabang(t string) string {
	return fmt.Sprintf(`SELECT ID, NAME, BRANCHPARENTID, KANWILGROUP, STATUS FROM %s WHERE STATUS = :1
	  ORDER BY UPPER(NAME), ID`, t)
}

func sqlAmbilCabang(t string) string {
	return fmt.Sprintf(`SELECT ID, NAME, BRANCHPARENTID, KANWILGROUP, STATUS FROM %s WHERE ID = :1`, t)
}

// pemindai - Scan milik *sql.Row dan *sql.Rows.
type pemindai interface{ Scan(...any) error }

func pindaiMO(p pemindai) (models.MarketingOfficer, error) {
	var v [14]sql.NullString
	tujuan := make([]any, len(v))
	for i := range v {
		tujuan[i] = &v[i]
	}
	if err := p.Scan(tujuan...); err != nil {
		return models.MarketingOfficer{}, err
	}
	return models.MarketingOfficer{ID: v[0].String, ClientID: v[1].String, ClientName: v[2].String,
		ClientID2: v[3].String, MOLeader: v[4].String, MOStatus: v[5].String, BranchParent: v[6].String,
		BranchDetailID: v[7].String, BranchDetailName: v[8].String, TeamGroup: v[9].String,
		BranchStatus: v[10].String, AksesLogin: v[11].String, UserUpdate: v[12].String, Tanggal: v[13].String}, nil
}

func pindaiAkun(p pemindai) (models.Akun, error) {
	var id, nama, kontak, email, aktif sql.NullString
	if err := p.Scan(&id, &nama, &kontak, &email, &aktif); err != nil {
		return models.Akun{}, err
	}
	return models.Akun{LoginID: id.String, Nama: nama.String, ContactID: kontak.String, Email: email.String,
		Aktif: aktif.String == "1"}, nil
}

func pindaiCabang(p pemindai) (models.Cabang, error) {
	var id, nama, induk, kanwil, status sql.NullString
	if err := p.Scan(&id, &nama, &induk, &kanwil, &status); err != nil {
		return models.Cabang{}, err
	}
	return models.Cabang{ID: id.String, Nama: nama.String, Induk: induk.String, KanwilGroup: kanwil.String,
		Aktif: status.String == "1"}, nil
}

// daftar menjalankan satu SELECT banyak baris.
func daftar[T any](ctx context.Context, j penjalan, q string, pindai func(pemindai) (T, error), args ...any) ([]T, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := j.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := []T{}
	for rows.Next() {
		v, err := pindai(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

// satu menjalankan satu SELECT satu baris; nol baris = ErrTidakAda.
func satu[T any](ctx context.Context, j penjalan, q string, pindai func(pemindai) (T, error), args ...any) (T, error) {
	var nol T
	if err := db.PeriksaSQL(q); err != nil {
		return nol, err
	}
	v, err := pindai(j.QueryRowContext(ctx, q, args...))
	if errors.Is(err, sql.ErrNoRows) {
		return nol, ErrTidakAda
	}
	if err != nil {
		return nol, fmt.Errorf("repository: %w", err)
	}
	return v, nil
}

func pindaiTeks(p pemindai) (string, error) {
	var s sql.NullString
	err := p.Scan(&s)
	return s.String, err
}

// DaftarMO membaca seluruh baris MO.
func (g *Gudang) DaftarMO(ctx context.Context) ([]models.MarketingOfficer, error) {
	t, err := g.nama(TabelMO)
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.db, sqlDaftarMO(t), pindaiMO)
}

// AmbilMO membaca satu baris MO.
func (g *Gudang) AmbilMO(ctx context.Context, tx *db.Tx, id string) (models.MarketingOfficer, error) {
	t, err := g.nama(TabelMO)
	if err != nil {
		return models.MarketingOfficer{}, err
	}
	return satu(ctx, g.dari(tx), sqlAmbilMO(t), pindaiMO, id)
}

// FormatIDMO membentuk ID MO dari nomor sequence: `1` + 7 digit.
func FormatIDMO(n int64) (string, error) {
	if n < 0 || n > 9999999 {
		return "", ErrNomorMelampaui
	}
	return fmt.Sprintf("1%07d", n), nil
}

// SisipMO menyisipkan baris baru dan mengembalikan ID-nya.
func (g *Gudang) SisipMO(ctx context.Context, tx *db.Tx, m models.MarketingOfficer) (string, error) {
	t, err := g.nama(TabelMO)
	if err != nil {
		return "", err
	}
	seq, err := g.nama(SeqMO)
	if err != nil {
		return "", err
	}
	j := g.dari(tx)
	qNomor := sqlNomorMO(seq)
	if err := db.PeriksaSQL(qNomor); err != nil {
		return "", err
	}
	var n int64
	if err := j.QueryRowContext(ctx, qNomor).Scan(&n); err != nil {
		return "", fmt.Errorf("repository: nomor MO: %w", err)
	}
	id, err := FormatIDMO(n)
	if err != nil {
		return "", err
	}
	m.ID = id
	q := sqlSisipMO(t)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	if _, err := j.ExecContext(ctx, q, NilaiSisip(m)...); err != nil {
		return "", fmt.Errorf("repository: menyisipkan MO: %w", err)
	}
	return id, nil
}

// NilaiSisip - dua belas nilai bind `sqlSisipMO`, berurutan; kosong = NULL.
func NilaiSisip(m models.MarketingOfficer) []any {
	n := db.KosongJadiNil
	return []any{m.ID, n(m.ClientID), n(m.BranchDetailID), n(m.MOLeader), n(m.MOStatus), n(m.ClientID2),
		n(m.TeamGroup), n(m.ClientName), n(m.UserUpdate), n(m.BranchParent), n(m.BranchDetailName), n(m.AksesLogin)}
}

// NilaiPerbarui - sepuluh nilai bind `sqlPerbaruiMO`, berurutan; kosong = NULL.
func NilaiPerbarui(m models.MarketingOfficer) []any {
	n := db.KosongJadiNil
	return []any{n(m.BranchDetailID), n(m.MOLeader), n(m.MOStatus), n(m.ClientID2), n(m.TeamGroup),
		n(m.UserUpdate), n(m.BranchParent), n(m.BranchDetailName), n(m.AksesLogin), m.ID}
}

// PerbaruiMO menulis ulang kolom yang boleh berubah; nol baris = ErrTidakAda.
func (g *Gudang) PerbaruiMO(ctx context.Context, tx *db.Tx, m models.MarketingOfficer) error {
	t, err := g.nama(TabelMO)
	if err != nil {
		return err
	}
	q := sqlPerbaruiMO(t)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := g.dari(tx).ExecContext(ctx, q, NilaiPerbarui(m)...)
	if err != nil {
		return fmt.Errorf("repository: memperbarui MO: %w", err)
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrTidakAda
	}
	return nil
}

// AktifDenganClientID - ID baris AKTIF lain ber-Marketing Code itu; `kecuali` kosong = tidak ada.
func (g *Gudang) AktifDenganClientID(ctx context.Context, tx *db.Tx, clientID, kecuali string) ([]string, error) {
	t, err := g.nama(TabelMO)
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.dari(tx), sqlAktifDenganClientID(t), pindaiTeks, models.StatusAktif, clientID,
		db.KosongJadiNil(kecuali))
}

// AktifDenganAkses - ID baris AKTIF lain ber-akun itu (tanpa beda huruf).
func (g *Gudang) AktifDenganAkses(ctx context.Context, tx *db.Tx, login, kecuali string) ([]string, error) {
	t, err := g.nama(TabelMO)
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.dari(tx), sqlAktifDenganAkses(t), pindaiTeks, models.StatusAktif, login,
		db.KosongJadiNil(kecuali))
}

// ClientIDDariAkses - Marketing Code yang sudah dipakai akun itu; `ada` false = belum pernah.
func (g *Gudang) ClientIDDariAkses(ctx context.Context, tx *db.Tx, login string) (string, bool, error) {
	t, err := g.nama(TabelMO)
	if err != nil {
		return "", false, err
	}
	d, err := daftar(ctx, g.dari(tx), sqlClientIDDariAkses(t), pindaiTeks, login)
	if err != nil || len(d) == 0 {
		return "", false, err
	}
	return d[0], true, nil
}

// DaftarAkun membaca seluruh akun M_LOGIN_GO (aktif dan nonaktif).
func (g *Gudang) DaftarAkun(ctx context.Context) ([]models.Akun, error) {
	t, err := g.nama(TabelLogin)
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.db, sqlDaftarAkun(t), pindaiAkun)
}

// AmbilAkun membaca satu akun M_LOGIN_GO.
func (g *Gudang) AmbilAkun(ctx context.Context, tx *db.Tx, login string) (models.Akun, error) {
	t, err := g.nama(TabelLogin)
	if err != nil {
		return models.Akun{}, err
	}
	return satu(ctx, g.dari(tx), sqlAmbilAkun(t), pindaiAkun, login)
}

// DaftarCabang membaca cabang aktif.
func (g *Gudang) DaftarCabang(ctx context.Context) ([]models.Cabang, error) {
	t, err := g.nama(TabelCabang)
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.db, sqlDaftarCabang(t), pindaiCabang, models.StatusAktif)
}

// AmbilCabang membaca satu cabang (aktif atau tidak).
func (g *Gudang) AmbilCabang(ctx context.Context, tx *db.Tx, id string) (models.Cabang, error) {
	t, err := g.nama(TabelCabang)
	if err != nil {
		return models.Cabang{}, err
	}
	return satu(ctx, g.dari(tx), sqlAmbilCabang(t), pindaiCabang, id)
}

// kolomLog - urutan kolom SELECT log (`pindaiLog`): 14 kolom baris MO, lalu ACTION dan LOG_TIME.
const kolomLog = `ID, CLIENTID, CLIENTNAME, CLIENTID2, MOLEADER, MOSTATUS, BRANCHPARENT, BRANCHDETAILID,
	BRANCHDETAILNAME, TEAMGROUP, BRANCHSTATUS, AKSES_LOGIN, USERUPDATE, TO_CHAR(TANGGAL, 'YYYY-MM-DD HH24:MI'),
	ACTION, TO_CHAR(LOG_TIME, 'YYYY-MM-DD HH24:MI:SS')`

// kolomLogLama - sebelum migrasi 760: AKSES_LOGIN dan LOG_TIME belum ada (dibaca NULL).
const kolomLogLama = `ID, CLIENTID, CLIENTNAME, CLIENTID2, MOLEADER, MOSTATUS, BRANCHPARENT, BRANCHDETAILID,
	BRANCHDETAILNAME, TEAMGROUP, BRANCHSTATUS, NULL, USERUPDATE, TO_CHAR(TANGGAL, 'YYYY-MM-DD HH24:MI'), ACTION, NULL`

// sqlLogMO - log satu MO, TERTUA lebih dulu: LOG_TIME (kosong = baris lama, lebih dulu), lalu TANGGAL, lalu ROWID.
// Untuk baris lama tanpa LOG_TIME urutannya perkiraan.
func sqlLogMO(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1 ORDER BY LOG_TIME NULLS FIRST, TANGGAL NULLS FIRST, ROWID`, kolomLog, t)
}

func sqlLogMOLama(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1 ORDER BY TANGGAL NULLS FIRST, ROWID`, kolomLogLama, t)
}

// sqlTandaiLog - baris log yang BARU disisipkan trigger untuk UPDATE ini (LOG_TIME terbesar ID itu, belum bertanda)
// mendapat AKSES_LOGIN lama dan tanda ACTION. :3 dan :4 sama-sama ID (go-ora mengikat menurut urutan placeholder).
func sqlTandaiLog(t string) string {
	return fmt.Sprintf(`UPDATE %s SET ACTION = :1, AKSES_LOGIN = :2
	  WHERE ID = :3 AND ACTION IS NULL AND LOG_TIME = (SELECT MAX(LOG_TIME) FROM %s WHERE ID = :4)`, t, t)
}

func pindaiLog(p pemindai) (models.BarisLog, error) {
	var v [16]sql.NullString
	tujuan := make([]any, len(v))
	for i := range v {
		tujuan[i] = &v[i]
	}
	if err := p.Scan(tujuan...); err != nil {
		return models.BarisLog{}, err
	}
	return models.BarisLog{MarketingOfficer: models.MarketingOfficer{ID: v[0].String, ClientID: v[1].String,
		ClientName: v[2].String, ClientID2: v[3].String, MOLeader: v[4].String, MOStatus: v[5].String,
		BranchParent: v[6].String, BranchDetailID: v[7].String, BranchDetailName: v[8].String, TeamGroup: v[9].String,
		BranchStatus: v[10].String, AksesLogin: v[11].String, UserUpdate: v[12].String, Tanggal: v[13].String},
		Aksi: v[14].String, LogTime: v[15].String}, nil
}

// kolomBelumAda - ORA-00904: kolom migrasi 760 belum ada.
func kolomBelumAda(err error) bool { return err != nil && strings.Contains(err.Error(), "ORA-00904") }

// LogMO membaca log satu MO, tertua lebih dulu. Sebelum migrasi 760 dibaca tanpa LOG_TIME dan AKSES_LOGIN.
func (g *Gudang) LogMO(ctx context.Context, id string) ([]models.BarisLog, error) {
	t, err := g.nama(TabelLog)
	if err != nil {
		return nil, err
	}
	d, err := daftar(ctx, g.db, sqlLogMO(t), pindaiLog, id)
	if kolomBelumAda(err) {
		return daftar(ctx, g.db, sqlLogMOLama(t), pindaiLog, id)
	}
	return d, err
}

// TandaiLog mengisi AKSES_LOGIN lama dan tanda `UPDATE-GO` pada baris log UPDATE yang baru saja terjadi - di
// transaksi yang sama. Nol baris (trigger tidak menulis) bukan galat; sebelum migrasi 760 = ErrLogBelumDimigrasi.
func (g *Gudang) TandaiLog(ctx context.Context, tx *db.Tx, id, aksesLama string) error {
	t, err := g.nama(TabelLog)
	if err != nil {
		return err
	}
	q := sqlTandaiLog(t)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	_, err = g.dari(tx).ExecContext(ctx, q, models.AksiGo, db.KosongJadiNil(aksesLama), id, id)
	if kolomBelumAda(err) {
		return ErrLogBelumDimigrasi
	}
	if err != nil {
		return fmt.Errorf("repository: menandai log MO: %w", err)
	}
	return nil
}
