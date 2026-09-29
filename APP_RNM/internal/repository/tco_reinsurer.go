package repository

// Reinsurer pada kombinasi (tahun, grup, jenis) - tiket 05 Treaty Contract Out.
//
// Untuk apa berkas ini: baca/tulis `T_TREATYREINSURER` (meniru logika
// `POOLDATA.PEGA_TREATYREINSURER`, 19 parameter, TANPA memanggilnya) dan
// pembaca master reinsurer `AGENT` (dibaca saja).
//
// `[terverifikasi]` daftar per KOMBINASI - `GetMasterReinsurerList.xml` dan
// `BrowseDetailTreatyReisurer_RD.xml` (saringan b651/b668/b685, urut `.ID ASC`
// b727/b730) - bukan per `ID` kontrak.
//
// Dibaca sesudah: tco_kontrak.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/internal/models"
	"nusantarare/pkg/utils"
)

// MasterReinsurerAgentTCO - master reinsurer, dibaca saja.
//
// `[terverifikasi]` pemilih `Reinsurer` `ViewDetailTreatyReinsurerGrid1.xml`
// b8285 -> `BrowseAgentReinsSOA_RD` (kelas `ASM-FW-GISFW-Int-AGENT`): `.ID` ->
// ReinsurerID, `.ClientName` -> NAME, `.ClientID` -> CLIENTID; saringan
// `.ClientName Contains` + `.StatusActive = 1`. Kolom fisik `ID`, `CLIENTNAME`,
// `CLIENTID` terbukti di SQL korpus (`Claim Fac In/RDBList/GetLeaderReport.xml`,
// `GetAddressCeding.xml`); `STATUSACTIVE` dari nama properti RD `[dugaan kuat]`
// (OQ-TCO-12).
const MasterReinsurerAgentTCO = "AGENT"

// NilaiAgentAktifTCO - `BrowseAgentReinsSOA_RD` b579 `pyFilterValue 1`.
const NilaiAgentAktifTCO = "1"

// batasCariReinsurerTCO membatasi hasil pemilih.
const batasCariReinsurerTCO = 100

// ErrReinsurerTidakAda - reinsurer bukan milik kombinasi itu.
var ErrReinsurerTidakAda = errors.New("repository: reinsurer tidak ditemukan pada kombinasi ini")

// ErrReinsurerMasterTidakAda - ID reinsurer tidak ada di master aktif.
var ErrReinsurerMasterTidakAda = errors.New("repository: reinsurer tidak ada di master AGENT yang aktif")

// ReinsurerMasterTCO adalah satu baris pemilih reinsurer.
type ReinsurerMasterTCO struct {
	ID, ClientName, ClientID string
}

// MasterReinsurerAgent membaca master `AGENT`.
type MasterReinsurerAgent struct{ db *DB }

// NewMasterReinsurerAgent menyusun pembacanya.
func NewMasterReinsurerAgent(db *DB) *MasterReinsurerAgent { return &MasterReinsurerAgent{db: db} }

// sqlCariReinsurerMasterTCO - `Contains` tanpa beda huruf (`pyCaseInsensitive`
// b569) dengan ESCAPE supaya `%` dan `_` ketikan pemakai tidak menjadi wildcard.
func sqlCariReinsurerMasterTCO(tabel string) string {
	return fmt.Sprintf(`SELECT ID, CLIENTNAME, CLIENTID FROM %s
	 WHERE STATUSACTIVE = :1 AND UPPER(CLIENTNAME) LIKE :2 ESCAPE '\'
	 ORDER BY CLIENTNAME, ID
	 FETCH FIRST %d ROWS ONLY`, tabel, batasCariReinsurerTCO)
}

func sqlAmbilReinsurerMasterTCO(tabel string) string {
	return fmt.Sprintf(`SELECT ID, CLIENTNAME, CLIENTID FROM %s WHERE ID = :1 AND STATUSACTIVE = :2`, tabel)
}

// polaLikeTCO membungkus teks ketikan menjadi pola LIKE yang aman.
func polaLikeTCO(teks string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return "%" + r.Replace(strings.ToUpper(strings.TrimSpace(teks))) + "%"
}

// Cari membaca reinsurer aktif yang namanya memuat teks.
func (m *MasterReinsurerAgent) Cari(ctx context.Context, teks string) ([]ReinsurerMasterTCO, error) {
	tabel, err := m.db.Qualify(MasterReinsurerAgentTCO)
	if err != nil {
		return nil, err
	}
	q := sqlCariReinsurerMasterTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := m.db.bacaTCO(ctx).QueryContext(ctx, q, NilaiAgentAktifTCO, polaLikeTCO(teks))
	if err != nil {
		return nil, fmt.Errorf("repository: membaca master %s: %w", MasterReinsurerAgentTCO, err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []ReinsurerMasterTCO
	for rows.Next() {
		var n [3]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2]); err != nil {
			return nil, err
		}
		hasil = append(hasil, ReinsurerMasterTCO{ID: n[0].String, ClientName: n[1].String, ClientID: n[2].String})
	}
	return hasil, rows.Err()
}

// Ambil membaca SATU reinsurer aktif.
func (m *MasterReinsurerAgent) Ambil(ctx context.Context, id string) (ReinsurerMasterTCO, error) {
	tabel, err := m.db.Qualify(MasterReinsurerAgentTCO)
	if err != nil {
		return ReinsurerMasterTCO{}, err
	}
	q := sqlAmbilReinsurerMasterTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return ReinsurerMasterTCO{}, err
	}
	var n [3]sql.NullString
	err = m.db.bacaTCO(ctx).QueryRowContext(ctx, q, id, NilaiAgentAktifTCO).Scan(&n[0], &n[1], &n[2])
	if errors.Is(err, sql.ErrNoRows) {
		return ReinsurerMasterTCO{}, ErrReinsurerMasterTidakAda
	}
	if err != nil {
		return ReinsurerMasterTCO{}, fmt.Errorf("repository: membaca master %s %s: %w", MasterReinsurerAgentTCO, id, err)
	}
	return ReinsurerMasterTCO{ID: n[0].String, ClientName: n[1].String, ClientID: n[2].String}, nil
}

// MasterReinsurerTCO membaca dan menulis `T_TREATYREINSURER`.
type MasterReinsurerTCO struct{ db *DB }

// NewMasterReinsurerTCO menyusun gudangnya.
func NewMasterReinsurerTCO(db *DB) *MasterReinsurerTCO { return &MasterReinsurerTCO{db: db} }

var pilihReinsurerTCO = `ID, TREATYYEAR, TREATYGROUPID, TREATYGROUPNAME, REINSTYPEID, REINSTYPENAME,
	       REINSURERID, CLIENTID, NAME, ` + fmt.Sprintf(fmtDesimal, "RICOMM") + `, ` +
	fmt.Sprintf(fmtDesimal, "PCTSHARE") + `, IUDATE, USERID,
	       TO_CHAR(STARTDATE, 'YYYY-MM-DD HH24:MI:SS'), TO_CHAR(ENDDATE, 'YYYY-MM-DD HH24:MI:SS'),
	       STATUSON, STDRATING, OPERATORNAME, TO_CHAR(TGLUPDATE, 'YYYY-MM-DD HH24:MI:SS')`

const saringKombinasiTCO = `TREATYYEAR = :1 AND TREATYGROUPID = :2 AND REINSTYPEID = :3`

func sqlDaftarReinsurerTCO(tabel string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE %s ORDER BY ID ASC`, pilihReinsurerTCO, tabel, saringKombinasiTCO)
}

func sqlAmbilReinsurerTCO(tabel string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE %s AND ID = :4`, pilihReinsurerTCO, tabel, saringKombinasiTCO)
}

// sqlShareLainTCO membaca share reinsurer LAIN, dikunci selama transaksi.
func sqlShareLainTCO(tabel string) string {
	return fmt.Sprintf(`SELECT ID, %s FROM %s WHERE %s AND (:4 IS NULL OR ID <> :5) FOR UPDATE`,
		fmt.Sprintf(fmtDesimal, "PCTSHARE"), tabel, saringKombinasiTCO)
}

func sqlSisipReinsurerTCO(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s
	       (ID, TREATYYEAR, TREATYGROUPID, TREATYGROUPNAME, REINSTYPEID, REINSTYPENAME, REINSURERID, CLIENTID,
	        NAME, RICOMM, PCTSHARE, IUDATE, USERID, STARTDATE, ENDDATE, STATUSON, STDRATING, OPERATORNAME, TGLUPDATE)
	VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16, :17, :18, :19)`, tabel)
}

// sqlPerbaruiReinsurerTCO menimpa seluruh medan milik baris (AC 8), dibatasi
// kombinasinya. Kunci kombinasi sendiri tidak berubah.
func sqlPerbaruiReinsurerTCO(tabel string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET TREATYGROUPNAME = :1, REINSTYPENAME = :2, REINSURERID = :3, CLIENTID = :4, NAME = :5,
	       RICOMM = :6, PCTSHARE = :7, IUDATE = :8, USERID = :9, STARTDATE = :10, ENDDATE = :11,
	       STATUSON = :12, STDRATING = :13, OPERATORNAME = :14, TGLUPDATE = :15
	 WHERE ID = :16 AND TREATYYEAR = :17 AND TREATYGROUPID = :18 AND REINSTYPEID = :19`, tabel)
}

func desimalDariTeks(v sql.NullString, kolom string) (*apd.Decimal, error) {
	if !v.Valid || strings.TrimSpace(v.String) == "" {
		return nil, nil
	}
	d, err := utils.ParseDecimal(strings.TrimSpace(v.String))
	if err != nil {
		return nil, fmt.Errorf("repository: kolom %s bernilai %q: %w", kolom, v.String, err)
	}
	return d, nil
}

func pindaiReinsurerTCO(baca interface{ Scan(...any) error }) (models.ReinsurerTreaty, error) {
	var n [19]sql.NullString
	tujuan := make([]any, len(n))
	for i := range n {
		tujuan[i] = &n[i]
	}
	if err := baca.Scan(tujuan...); err != nil {
		return models.ReinsurerTreaty{}, err
	}
	r := models.ReinsurerTreaty{ID: n[0].String, TreatyYear: n[1].String, TreatyGroupID: n[2].String,
		TreatyGroupName: n[3].String, ReinsTypeID: n[4].String, ReinsTypeName: n[5].String,
		ReinsurerID: n[6].String, ClientID: n[7].String, Name: n[8].String, IUDate: n[11].String,
		UserID: n[12].String, StatusOn: n[15].String, StdRating: n[16].String, OperatorName: n[17].String}
	var err error
	if r.Ricomm, err = desimalDariTeks(n[9], "RICOMM"); err != nil {
		return r, err
	}
	if r.PctShare, err = desimalDariTeks(n[10], "PCTSHARE"); err != nil {
		return r, err
	}
	if r.StartDate, err = uraiTanggalTeks(n[13], "STARTDATE"); err != nil {
		return r, err
	}
	if r.EndDate, err = uraiTanggalTeks(n[14], "ENDDATE"); err != nil {
		return r, err
	}
	r.TglUpdate, err = uraiTanggalTeks(n[18], "TGLUPDATE")
	return r, err
}

func argKombinasi(k models.KombinasiTCO) []any {
	return []any{k.TreatyYear, k.TreatyGroupID, k.ReinsTypeID}
}

func desimalJadiNil(d *apd.Decimal) any {
	if d == nil {
		return nil
	}
	return d.Text('f')
}

// Daftar membaca reinsurer satu kombinasi, `ID ASC` (RD b730).
func (m *MasterReinsurerTCO) Daftar(ctx context.Context, k models.KombinasiTCO) ([]models.ReinsurerTreaty, error) {
	tabel, err := m.db.Qualify(TabelReinsurerTCO)
	if err != nil {
		return nil, err
	}
	q := sqlDaftarReinsurerTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := m.db.bacaTCO(ctx).QueryContext(ctx, q, argKombinasi(k)...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca reinsurer kombinasi: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []models.ReinsurerTreaty
	for rows.Next() {
		r, err := pindaiReinsurerTCO(rows)
		if err != nil {
			return nil, err
		}
		hasil = append(hasil, r)
	}
	return hasil, rows.Err()
}

// Ambil membaca SATU reinsurer milik kombinasi itu.
func (m *MasterReinsurerTCO) Ambil(ctx context.Context, k models.KombinasiTCO, id string) (models.ReinsurerTreaty, error) {
	tabel, err := m.db.Qualify(TabelReinsurerTCO)
	if err != nil {
		return models.ReinsurerTreaty{}, err
	}
	q := sqlAmbilReinsurerTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return models.ReinsurerTreaty{}, err
	}
	r, err := pindaiReinsurerTCO(m.db.bacaTCO(ctx).QueryRowContext(ctx, q, append(argKombinasi(k), id)...))
	if errors.Is(err, sql.ErrNoRows) {
		return models.ReinsurerTreaty{}, ErrReinsurerTidakAda
	}
	if err != nil {
		return models.ReinsurerTreaty{}, fmt.Errorf("repository: membaca reinsurer %s: %w", id, err)
	}
	return r, nil
}

// ShareLain membaca share reinsurer LAIN pada kombinasi, di dalam transaksi.
func (m *MasterReinsurerTCO) ShareLain(ctx context.Context, tx *Tx, k models.KombinasiTCO, kecualiID string) (
	[]*apd.Decimal, error) {
	if tx == nil {
		return nil, errors.New("repository: membaca share lain menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelReinsurerTCO)
	if err != nil {
		return nil, err
	}
	q := sqlShareLainTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := tx.tx.QueryContext(ctx, q, append(argKombinasi(k), kosongJadiNil(kecualiID), kosongJadiNil(kecualiID))...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca share reinsurer lain: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []*apd.Decimal
	for rows.Next() {
		var id, share sql.NullString
		if err := rows.Scan(&id, &share); err != nil {
			return nil, err
		}
		d, err := desimalDariTeks(share, "PCTSHARE")
		if err != nil {
			return nil, err
		}
		hasil = append(hasil, d)
	}
	return hasil, rows.Err()
}

// Sisip menulis reinsurer baru; ID dari `SEQ_T_TREATYREINSURER` (ADR-0006).
func (m *MasterReinsurerTCO) Sisip(ctx context.Context, tx *Tx, r models.ReinsurerTreaty) (string, error) {
	if tx == nil {
		return "", errors.New("repository: menyisipkan reinsurer menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelReinsurerTCO)
	if err != nil {
		return "", err
	}
	id, err := m.db.IdentitasBerikutTCO(ctx, tx, SeqReinsurerTCO)
	if err != nil {
		return "", err
	}
	q := sqlSisipReinsurerTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, id, kosongJadiNil(r.TreatyYear), kosongJadiNil(r.TreatyGroupID),
		kosongJadiNil(r.TreatyGroupName), kosongJadiNil(r.ReinsTypeID), kosongJadiNil(r.ReinsTypeName),
		kosongJadiNil(r.ReinsurerID), kosongJadiNil(r.ClientID), kosongJadiNil(r.Name),
		desimalJadiNil(r.Ricomm), desimalJadiNil(r.PctShare), kosongJadiNil(r.IUDate), kosongJadiNil(r.UserID),
		tanggalJadiNil(r.StartDate), tanggalJadiNil(r.EndDate), kosongJadiNil(r.StatusOn),
		kosongJadiNil(r.StdRating), kosongJadiNil(r.OperatorName), tanggalJadiNil(r.TglUpdate))
	if err != nil {
		return "", fmt.Errorf("repository: menyisipkan reinsurer: %w", err)
	}
	return id, pastikanSatuBaris(hasil, "penyisipan reinsurer")
}

// Perbarui menimpa reinsurer yang ada, dibatasi kombinasinya.
func (m *MasterReinsurerTCO) Perbarui(ctx context.Context, tx *Tx, r models.ReinsurerTreaty) error {
	if tx == nil {
		return errors.New("repository: memperbarui reinsurer menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelReinsurerTCO)
	if err != nil {
		return err
	}
	q := sqlPerbaruiReinsurerTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, kosongJadiNil(r.TreatyGroupName), kosongJadiNil(r.ReinsTypeName),
		kosongJadiNil(r.ReinsurerID), kosongJadiNil(r.ClientID), kosongJadiNil(r.Name),
		desimalJadiNil(r.Ricomm), desimalJadiNil(r.PctShare), kosongJadiNil(r.IUDate), kosongJadiNil(r.UserID),
		tanggalJadiNil(r.StartDate), tanggalJadiNil(r.EndDate), kosongJadiNil(r.StatusOn),
		kosongJadiNil(r.StdRating), kosongJadiNil(r.OperatorName), tanggalJadiNil(r.TglUpdate),
		r.ID, r.TreatyYear, r.TreatyGroupID, r.ReinsTypeID)
	if err != nil {
		return fmt.Errorf("repository: memperbarui reinsurer %s: %w", r.ID, err)
	}
	if n, err := hasil.RowsAffected(); err == nil && n == 0 {
		return ErrReinsurerTidakAda
	}
	return pastikanSatuBaris(hasil, "pembaruan reinsurer")
}
