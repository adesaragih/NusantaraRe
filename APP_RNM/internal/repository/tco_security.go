package repository

// Security reinsurer - tiket 06 Treaty Contract Out.
//
// Untuk apa berkas ini: baca/tulis `T_MTREATYSECURITY` dengan struktur bersih
// (penyimpangan sadar 5) - pengganti tiga SQL mentah warisan:
//
//	InsertToMTreatySecurity.xml b60   insert ... values (7 nilai POSISIONAL)
//	UpdateMTreatySecurity.xml b85-89  ... where REAS_ID = .. and trim(REAS_SECURITY) = trim(..)
//	DeleteSecurityReinsurer.xml b85   ... where REAS_ID = .. and trim(REAS_SECURITY) = trim(..)
//
// ⛔ Setiap penulisan berdaftar-kolom (AC 19); setiap pencocokan berkunci `ID`
// surrogate + `REAS_ID` - nama security tidak pernah menjadi kunci (AC 18) dan
// tidak ada fungsi pemangkas spasi di SQL mana pun (AC 20).
//
// ⛔ `TOP_ID`, `TP_TREATY`, `USER_ID` ditulis NULL EKSPLISIT saat sisip -
// warisan mengosongkannya (`''`) - dan tidak disentuh saat diperbarui, sama
// dengan UPDATE warisan. Artinya `[terbuka]` (blocker tiket 06).
//
// Dibaca sesudah: tco_reinsurer.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/internal/models"
)

// ErrSecurityTidakAda - baris security bukan milik reinsurer itu.
var ErrSecurityTidakAda = errors.New("repository: security tidak ditemukan pada reinsurer ini")

// kolomSecurityTCO - daftar kolom sisip, BERNAMA (AC 19).
var kolomSecurityTCO = []string{"ID", "THN_TREATY", "TOP_ID", "TP_TREATY", "REAS_ID", "PCT_SHARE", "USER_ID", "REAS_SECURITY"}

// SecurityTCO adalah satu baris security beserta nama tampilnya dari master.
type SecurityTCO struct {
	models.SecurityReinsurer
	// ClientName - `.CLIENTNAME` grid b16878, dari master `AGENT`; tidak disimpan.
	ClientName string
}

// MasterSecurityTCO membaca dan menulis `T_MTREATYSECURITY`.
type MasterSecurityTCO struct{ db *DB }

// NewMasterSecurityTCO menyusun gudangnya.
func NewMasterSecurityTCO(db *DB) *MasterSecurityTCO { return &MasterSecurityTCO{db: db} }

func pilihSecurityTCO() string {
	return "s.ID, s.THN_TREATY, s.TOP_ID, s.TP_TREATY, s.REAS_ID, " + fmt.Sprintf(fmtDesimal, "s.PCT_SHARE") +
		", s.USER_ID, s.REAS_SECURITY, a.CLIENTNAME"
}

// sqlDaftarSecurityTCO - `SelectSecurityReinsurer` saringan `A AND D` (b550):
// `.THN_TREATY = Param.THN_TREATY` dan `.REAS_ID = Param.REAS_ID`; nama dari
// master `AGENT` (LEFT JOIN - security yang agennya nonaktif tetap tampil).
func sqlDaftarSecurityTCO(tabel, agent string) string {
	return fmt.Sprintf(`SELECT %s
	  FROM %s s LEFT JOIN %s a ON a.ID = s.REAS_SECURITY
	 WHERE s.REAS_ID = :1 AND s.THN_TREATY = :2
	 ORDER BY s.ID ASC`, pilihSecurityTCO(), tabel, agent)
}

func sqlAmbilSecurityTCO(tabel, agent string) string {
	return fmt.Sprintf(`SELECT %s
	  FROM %s s LEFT JOIN %s a ON a.ID = s.REAS_SECURITY
	 WHERE s.ID = :1 AND s.REAS_ID = :2`, pilihSecurityTCO(), tabel, agent)
}

// sqlCariDobelSecurityTCO - security yang sama di bawah reinsurer yang sama,
// dicocokkan PERSIS (nilai dibersihkan sekali di batas masukan).
func sqlCariDobelSecurityTCO(tabel string) string {
	return fmt.Sprintf(`SELECT ID FROM %s WHERE REAS_ID = :1 AND REAS_SECURITY = :2 AND (:3 IS NULL OR ID <> :4)
	 ORDER BY ID FETCH FIRST 1 ROWS ONLY`, tabel)
}

func sqlSisipSecurityTCO(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s
	       (ID, THN_TREATY, TOP_ID, TP_TREATY, REAS_ID, PCT_SHARE, USER_ID, REAS_SECURITY)
	VALUES (:1, :2, NULL, NULL, :3, :4, NULL, :5)`, tabel)
}

// sqlPerbaruiSecurityTCO - tiga kolom yang UPDATE warisan tulis (b86-b88),
// berkunci `ID` + `REAS_ID`.
func sqlPerbaruiSecurityTCO(tabel string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET THN_TREATY = :1, PCT_SHARE = :2, REAS_SECURITY = :3
	 WHERE ID = :4 AND REAS_ID = :5`, tabel)
}

func sqlHapusSecurityTCO(tabel string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE ID = :1 AND REAS_ID = :2`, tabel)
}

func pindaiSecurityTCO(baca interface{ Scan(...any) error }) (SecurityTCO, error) {
	var n [9]sql.NullString
	tujuan := make([]any, len(n))
	for i := range n {
		tujuan[i] = &n[i]
	}
	if err := baca.Scan(tujuan...); err != nil {
		return SecurityTCO{}, err
	}
	s := SecurityTCO{SecurityReinsurer: models.SecurityReinsurer{ID: n[0].String, ThnTreaty: n[1].String,
		TopID: n[2].String, TpTreaty: n[3].String, ReasID: n[4].String, UserID: n[6].String,
		ReasSecurity: n[7].String}, ClientName: n[8].String}
	var err error
	s.PctShare, err = desimalDariTeks(n[5], "PCT_SHARE")
	return s, err
}

func (m *MasterSecurityTCO) tabelBaca() (string, string, error) {
	tabel, err := m.db.Qualify(TabelSecurityTCO)
	if err != nil {
		return "", "", err
	}
	agent, err := m.db.Qualify(MasterReinsurerAgentTCO)
	if err != nil {
		return "", "", err
	}
	return tabel, agent, nil
}

// Daftar membaca security seorang reinsurer, `ID ASC`.
func (m *MasterSecurityTCO) Daftar(ctx context.Context, reasID, thnTreaty string) ([]SecurityTCO, error) {
	tabel, agent, err := m.tabelBaca()
	if err != nil {
		return nil, err
	}
	q := sqlDaftarSecurityTCO(tabel, agent)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := m.db.sql.QueryContext(ctx, q, reasID, thnTreaty)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca security reinsurer %s: %w", reasID, err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []SecurityTCO
	for rows.Next() {
		s, err := pindaiSecurityTCO(rows)
		if err != nil {
			return nil, err
		}
		hasil = append(hasil, s)
	}
	return hasil, rows.Err()
}

// Ambil membaca satu security milik reinsurer itu.
func (m *MasterSecurityTCO) Ambil(ctx context.Context, reasID, id string) (SecurityTCO, error) {
	tabel, agent, err := m.tabelBaca()
	if err != nil {
		return SecurityTCO{}, err
	}
	q := sqlAmbilSecurityTCO(tabel, agent)
	if err := PeriksaSQL(q); err != nil {
		return SecurityTCO{}, err
	}
	s, err := pindaiSecurityTCO(m.db.sql.QueryRowContext(ctx, q, id, reasID))
	if errors.Is(err, sql.ErrNoRows) {
		return SecurityTCO{}, ErrSecurityTidakAda
	}
	if err != nil {
		return SecurityTCO{}, fmt.Errorf("repository: membaca security %s: %w", id, err)
	}
	return s, nil
}

// CariDobel mencari baris LAIN yang memegang security yang sama.
func (m *MasterSecurityTCO) CariDobel(ctx context.Context, tx *Tx, reasID, reasSecurity, kecualiID string) (string, error) {
	if tx == nil {
		return "", errors.New("repository: pencarian dobel security menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelSecurityTCO)
	if err != nil {
		return "", err
	}
	q := sqlCariDobelSecurityTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	var id string
	err = tx.tx.QueryRowContext(ctx, q, reasID, reasSecurity, kosongJadiNil(kecualiID), kosongJadiNil(kecualiID)).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: mencari security dobel: %w", err)
	}
	return id, nil
}

// Sisip menulis security baru; ID dari `SEQ_T_MTREATYSECURITY`.
func (m *MasterSecurityTCO) Sisip(ctx context.Context, tx *Tx, s models.SecurityReinsurer) (string, error) {
	if tx == nil {
		return "", errors.New("repository: menyisipkan security menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelSecurityTCO)
	if err != nil {
		return "", err
	}
	id, err := m.db.IdentitasBerikutTCO(ctx, tx, SeqSecurityTCO)
	if err != nil {
		return "", err
	}
	q := sqlSisipSecurityTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, id, kosongJadiNil(s.ThnTreaty), s.ReasID, desimalJadiNil(s.PctShare),
		s.ReasSecurity)
	if err != nil {
		return "", fmt.Errorf("repository: menyisipkan security: %w", err)
	}
	return id, pastikanSatuBaris(hasil, "penyisipan security")
}

// Perbarui menimpa tahun, share, dan security - berkunci ID + REAS_ID (AC 18/20).
func (m *MasterSecurityTCO) Perbarui(ctx context.Context, tx *Tx, s models.SecurityReinsurer) error {
	if tx == nil {
		return errors.New("repository: memperbarui security menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelSecurityTCO)
	if err != nil {
		return err
	}
	q := sqlPerbaruiSecurityTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, kosongJadiNil(s.ThnTreaty), desimalJadiNil(s.PctShare), s.ReasSecurity,
		s.ID, s.ReasID)
	if err != nil {
		return fmt.Errorf("repository: memperbarui security %s: %w", s.ID, err)
	}
	if n, err := hasil.RowsAffected(); err == nil && n == 0 {
		return ErrSecurityTidakAda
	}
	return pastikanSatuBaris(hasil, "pembaruan security")
}

// Hapus membuang SATU security - baris reinsurer induknya tidak disentuh.
func (m *MasterSecurityTCO) Hapus(ctx context.Context, tx *Tx, reasID, id string) error {
	if tx == nil {
		return errors.New("repository: menghapus security menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelSecurityTCO)
	if err != nil {
		return err
	}
	q := sqlHapusSecurityTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, id, reasID)
	if err != nil {
		return fmt.Errorf("repository: menghapus security %s: %w", id, err)
	}
	if n, err := hasil.RowsAffected(); err == nil && n == 0 {
		return ErrSecurityTidakAda
	}
	return pastikanSatuBaris(hasil, "penghapusan security")
}
