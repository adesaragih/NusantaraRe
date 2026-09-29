package repository

// Security reinsurer - tiket 06 Treaty Contract Out, tco4: tabel WARISAN
// `MTREATYSECURITY` (tanpa PK `[data DBA]`).
//
// Untuk apa berkas ini: baca/tulis `MTREATYSECURITY` persis seperti tiga SQL
// warisannya, yang tidak dipanggil tetapi ditiru:
//
//	InsertToMTreatySecurity.xml b60   insert ... values (7 nilai POSISIONAL)
//	UpdateMTreatySecurity.xml b85     ... where REAS_ID = .. and trim(REAS_SECURITY) = trim(..)
//	DeleteSecurityReinsurer.xml b85   ... where REAS_ID = .. and trim(REAS_SECURITY) = trim(..)
//
// ⛔ Kunci baris = (`REAS_ID`, `TRIM(REAS_SECURITY)`) - seperti Pega. Tabel
// warisan tidak punya identitas lain; `REAS_SECURITY` CHAR(10) berekor spasi,
// maka `TRIM` di kedua sisi. `models.SecurityReinsurer.ID` membawa
// `REAS_SECURITY` yang sudah dipangkas - satu security per reinsurer
// (OQ-TCO-17, keputusan work owner) menjamin kunci itu unik.
// RALAT tco4 atas tiket 06 AC 18/AC 20 (ID surrogate, nol TRIM di SQL):
// tabelnya tidak punya ID.
//
// Sisipan menulis DAFTAR KOLOM bernama (AC 19) dengan nilai yang sama dengan
// sisipan posisional warisan: `TOP_ID`, `TP_TREATY`, `USER_ID` = `''` (NULL).
//
// Dibaca sesudah: tco_reinsurer.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/internal/models"
	"nusantarare/inti/db"
)

// ErrSecurityTidakAda - baris security bukan milik reinsurer itu.
var ErrSecurityTidakAda = errors.New("repository: security tidak ditemukan pada reinsurer ini")

// kolomSecurityTCO - daftar kolom sisip, BERNAMA (AC 19), urutan posisi
// `InsertToMTreatySecurity` b60.
var kolomSecurityTCO = KolomWarisanTCO(warisanSecurityTCO)

// SecurityTCO adalah satu baris security beserta nama tampilnya dari master.
type SecurityTCO struct {
	models.SecurityReinsurer
	// ClientName - `.CLIENTNAME` grid b16878, dari master `AGENT`; tidak disimpan.
	ClientName string
}

// MasterSecurityTCO membaca dan menulis `MTREATYSECURITY`.
type MasterSecurityTCO struct{ db *db.DB }

// NewMasterSecurityTCO menyusun gudangnya.
func NewMasterSecurityTCO(db *db.DB) *MasterSecurityTCO { return &MasterSecurityTCO{db: db} }

// pilihSecurityTCO - `PCT_SHARE` VARCHAR2(99) dibaca APA ADANYA (diurai di Go,
// titik atau koma); kolom CHAR dipangkas di Go.
func pilihSecurityTCO() string {
	return "s.THN_TREATY, s.TOP_ID, s.TP_TREATY, s.REAS_ID, s.PCT_SHARE, s.USER_ID, s.REAS_SECURITY, a.CLIENTNAME"
}

// sqlDaftarSecurityTCO - `SelectSecurityReinsurer` saringan `A AND D` (b550):
// `.THN_TREATY = Param.THN_TREATY` dan `.REAS_ID = Param.REAS_ID`; nama dari
// master `AGENT` (LEFT JOIN - security yang agennya nonaktif tetap tampil).
// Urutan `TRIM(REAS_SECURITY)` `[keputusan kami]`: tabel tanpa identitas.
func sqlDaftarSecurityTCO(tabel, agent string) string {
	return fmt.Sprintf(`SELECT %s
	  FROM %s s LEFT JOIN %s a ON a.ID = TRIM(s.REAS_SECURITY)
	 WHERE s.REAS_ID = :1 AND s.THN_TREATY = :2
	 ORDER BY TRIM(s.REAS_SECURITY) ASC`, pilihSecurityTCO(), tabel, agent)
}

func sqlAmbilSecurityTCO(tabel, agent string) string {
	return fmt.Sprintf(`SELECT %s
	  FROM %s s LEFT JOIN %s a ON a.ID = TRIM(s.REAS_SECURITY)
	 WHERE s.REAS_ID = :1 AND TRIM(s.REAS_SECURITY) = TRIM(:2)`, pilihSecurityTCO(), tabel, agent)
}

// sqlCariDobelSecurityTCO - security yang sama di bawah reinsurer yang sama;
// baris yang sedang diubah (kunci lamanya) dikecualikan.
func sqlCariDobelSecurityTCO(tabel string) string {
	return fmt.Sprintf(`SELECT TRIM(REAS_SECURITY) FROM %s
	 WHERE REAS_ID = :1 AND TRIM(REAS_SECURITY) = TRIM(:2) AND (:3 IS NULL OR TRIM(REAS_SECURITY) <> TRIM(:4))
	 FETCH FIRST 1 ROWS ONLY`, tabel)
}

func sqlSisipSecurityTCO(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s
	       (THN_TREATY, TOP_ID, TP_TREATY, REAS_ID, PCT_SHARE, USER_ID, REAS_SECURITY)
	VALUES (:1, NULL, NULL, :2, :3, NULL, :4)`, tabel)
}

// sqlPerbaruiSecurityTCO - tiga kolom yang UPDATE warisan tulis (b86-b88),
// berkunci `REAS_ID` + `TRIM(REAS_SECURITY)` LAMA (b89).
func sqlPerbaruiSecurityTCO(tabel string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET THN_TREATY = :1, PCT_SHARE = :2, REAS_SECURITY = :3
	 WHERE REAS_ID = :4 AND TRIM(REAS_SECURITY) = TRIM(:5)`, tabel)
}

// sqlHapusSecurityTCO - `DeleteSecurityReinsurer` b85.
func sqlHapusSecurityTCO(tabel string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE REAS_ID = :1 AND TRIM(REAS_SECURITY) = TRIM(:2)`, tabel)
}

func pindaiSecurityTCO(baca interface{ Scan(...any) error }) (SecurityTCO, error) {
	var n [8]sql.NullString
	tujuan := make([]any, len(n))
	for i := range n {
		tujuan[i] = &n[i]
	}
	if err := baca.Scan(tujuan...); err != nil {
		return SecurityTCO{}, err
	}
	pangkas := func(v sql.NullString) string { return strings.TrimSpace(v.String) }
	s := SecurityTCO{SecurityReinsurer: models.SecurityReinsurer{ID: pangkas(n[6]), ThnTreaty: pangkas(n[0]),
		TopID: pangkas(n[1]), TpTreaty: pangkas(n[2]), ReasID: pangkas(n[3]), UserID: pangkas(n[5]),
		ReasSecurity: pangkas(n[6])}, ClientName: n[7].String}
	d, alasan, ok := UraiDesimalWarisanTCO(n[4].String)
	if !ok {
		return s, fmt.Errorf("repository: kolom PCT_SHARE bernilai %q: %s", n[4].String, alasan)
	}
	s.PctShare = d
	return s, nil
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

// Daftar membaca security seorang reinsurer.
func (m *MasterSecurityTCO) Daftar(ctx context.Context, reasID, thnTreaty string) ([]SecurityTCO, error) {
	tabel, agent, err := m.tabelBaca()
	if err != nil {
		return nil, err
	}
	q := sqlDaftarSecurityTCO(tabel, agent)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := bacaTCO(ctx, m.db).QueryContext(ctx, q, reasID, thnTreaty)
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

// Ambil membaca satu security milik reinsurer itu; `id` = `REAS_SECURITY`.
func (m *MasterSecurityTCO) Ambil(ctx context.Context, reasID, id string) (SecurityTCO, error) {
	tabel, agent, err := m.tabelBaca()
	if err != nil {
		return SecurityTCO{}, err
	}
	q := sqlAmbilSecurityTCO(tabel, agent)
	if err := db.PeriksaSQL(q); err != nil {
		return SecurityTCO{}, err
	}
	s, err := pindaiSecurityTCO(bacaTCO(ctx, m.db).QueryRowContext(ctx, q, reasID, id))
	if errors.Is(err, sql.ErrNoRows) {
		return SecurityTCO{}, ErrSecurityTidakAda
	}
	if err != nil {
		return SecurityTCO{}, fmt.Errorf("repository: membaca security %s: %w", id, err)
	}
	return s, nil
}

// CariDobel mencari baris LAIN yang memegang security yang sama; `kecualiID`
// = kunci lama baris yang sedang diubah.
func (m *MasterSecurityTCO) CariDobel(ctx context.Context, tx *db.Tx, reasID, reasSecurity, kecualiID string) (string, error) {
	if tx == nil {
		return "", errors.New("repository: pencarian dobel security menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelSecurityTCO)
	if err != nil {
		return "", err
	}
	q := sqlCariDobelSecurityTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var id string
	err = tx.QueryRowContext(ctx, q, reasID, reasSecurity, db.KosongJadiNil(kecualiID), db.KosongJadiNil(kecualiID)).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: mencari security dobel: %w", err)
	}
	return id, nil
}

// Sisip menulis security baru; "identitas"nya = `REAS_SECURITY` terpangkas.
func (m *MasterSecurityTCO) Sisip(ctx context.Context, tx *db.Tx, s models.SecurityReinsurer) (string, error) {
	if tx == nil {
		return "", errors.New("repository: menyisipkan security menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelSecurityTCO)
	if err != nil {
		return "", err
	}
	q := sqlSisipSecurityTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	hasil, err := tx.ExecContext(ctx, q, db.KosongJadiNil(s.ThnTreaty), s.ReasID, TulisDesimalWarisanTCO(s.PctShare),
		s.ReasSecurity)
	if err != nil {
		return "", fmt.Errorf("repository: menyisipkan security: %w", err)
	}
	return strings.TrimSpace(s.ReasSecurity), db.PastikanSatuBaris(hasil, "penyisipan security")
}

// Perbarui menimpa tahun, share, dan security - berkunci `REAS_ID` +
// `TRIM(REAS_SECURITY)` lama (`s.ID`).
func (m *MasterSecurityTCO) Perbarui(ctx context.Context, tx *db.Tx, s models.SecurityReinsurer) error {
	if tx == nil {
		return errors.New("repository: memperbarui security menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelSecurityTCO)
	if err != nil {
		return err
	}
	q := sqlPerbaruiSecurityTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, db.KosongJadiNil(s.ThnTreaty), TulisDesimalWarisanTCO(s.PctShare),
		s.ReasSecurity, s.ReasID, s.ID)
	if err != nil {
		return fmt.Errorf("repository: memperbarui security %s: %w", s.ID, err)
	}
	return palingSedikitSatuSecurityTCO(hasil)
}

// Hapus membuang SATU security - baris reinsurer induknya tidak disentuh.
func (m *MasterSecurityTCO) Hapus(ctx context.Context, tx *db.Tx, reasID, id string) error {
	if tx == nil {
		return errors.New("repository: menghapus security menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelSecurityTCO)
	if err != nil {
		return err
	}
	q := sqlHapusSecurityTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, reasID, id)
	if err != nil {
		return fmt.Errorf("repository: menghapus security %s: %w", id, err)
	}
	return palingSedikitSatuSecurityTCO(hasil)
}

// palingSedikitSatuSecurityTCO - UPDATE/DELETE berkunci nama mengenai SEMUA
// baris senama, seperti `UpdateMTreatySecurity`/`DeleteSecurityReinsurer`:
// warisan tanpa PK boleh sudah menyimpan duplikat, dan menolak lebih dari satu
// baris membuat duplikat itu tidak pernah dapat diperbaiki (temuan /code-review).
func palingSedikitSatuSecurityTCO(hasil sql.Result) error {
	n, err := hasil.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: mencacah baris security: %w", err)
	}
	if n == 0 {
		return ErrSecurityTidakAda
	}
	return nil
}

// desimalWarisanTCO - teks desimal warisan (VARCHAR2) -> desimal; galat
// menyebut kolomnya.
func desimalWarisanTCO(v sql.NullString, kolom string) (*apd.Decimal, error) {
	d, alasan, ok := UraiDesimalWarisanTCO(v.String)
	if !ok {
		return nil, fmt.Errorf("repository: kolom %s bernilai %q: %s", kolom, v.String, alasan)
	}
	return d, nil
}
