package repository

// Tahun treaty - `T_TREATYYEAR` (tiket 03 Treaty Contract Out).
//
// Padanan `RDBList/SaveMasterTreatyYear_SQL.xml` -> `PEGA_TREATYYEAR`
// (`[data DBA]` UPSERT dikunci ID; ID baru `'1' || lpad(TreatyYear_seq, 6)`;
// TGLUPDATE = SYSDATE; tidak COMMIT sendiri) dan
// `ReportDefinition/BrowseTreatyYear_RD.xml` (sort `.ID DESC` b672, maks 500
// b757). Prosedurnya TIDAK dipanggil (keputusan o): logikanya di sini.
//
// ⛔ Upsert dipecah dua: `Sisip` (ID dari sequence, ADR-0006) dan `Perbarui`
// (SELURUH kolom yang dikirim, WHERE ID; nol baris = tidak ada). Pemanggil
// yang memutuskan mana - ID kosong berarti baru (AC 8).
//
// ⛔ Anti-dobel (AC 73) adalah pertanyaan keunikan di basis data: `CariDobel`
// dijalankan DI DALAM transaksi penyimpanan, bukan dari cache layar.
//
// Tanggal dibaca lewat TO_CHAR berformat dan ditulis sebagai time.Time;
// nol pergeseran zona. Setiap query menyebut skemanya lewat Qualify
// (ADR-U-0033), nol COMMIT (ADR-U-0029).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"nusantarare/internal/models"
	"nusantarare/pkg/utils"
)

// ErrTahunTreatyTidakAda - tidak ada baris dengan ID itu.
var ErrTahunTreatyTidakAda = errors.New("repository: tahun treaty tidak ditemukan")

// HalamanTahunTreaty adalah satu halaman daftar tahun treaty.
type HalamanTahunTreaty struct {
	Baris   []models.TahunTreaty
	Total   int
	Halaman int
	Ukuran  int
}

// MasterTahunTreaty membaca dan menulis `T_TREATYYEAR`.
type MasterTahunTreaty struct{ db *DB }

// NewMasterTahunTreaty menyusunnya.
func NewMasterTahunTreaty(db *DB) *MasterTahunTreaty { return &MasterTahunTreaty{db: db} }

// daftarPilihTahunTreaty adalah sepuluh kolom warisan, tanggal sebagai teks.
const daftarPilihTahunTreaty = `ID, TREATYYEAR, UNDERWRITINGYEAR, TREATYGROUPID, TREATYGROUPNAME, PROPORTION,
	       TO_CHAR(STARTDATE, 'YYYY-MM-DD HH24:MI:SS'), TO_CHAR(ENDDATE, 'YYYY-MM-DD HH24:MI:SS'),
	       USERID, TO_CHAR(TGLUPDATE, 'YYYY-MM-DD HH24:MI:SS')`

// sqlDaftarTahunTreaty - urutan `.ID DESC` (BrowseTreatyYear_RD b672).
func sqlDaftarTahunTreaty(tabel string) string {
	return fmt.Sprintf(`SELECT %s
	  FROM %s
	 ORDER BY ID DESC
	 OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY`, daftarPilihTahunTreaty, tabel)
}

// sqlCacahTahunTreaty - query TERSENDIRI (pola sqlCacahInboxPolis).
func sqlCacahTahunTreaty(tabel string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s`, tabel)
}

func sqlAmbilTahunTreaty(tabel string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, daftarPilihTahunTreaty, tabel)
}

// sqlSisipTahunTreaty - INSERT kolom BERNAMA, urutan parameter procedure.
func sqlSisipTahunTreaty(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s
	       (ID, TREATYYEAR, UNDERWRITINGYEAR, TREATYGROUPID, TREATYGROUPNAME,
	        USERID, TGLUPDATE, PROPORTION, STARTDATE, ENDDATE)
	VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10)`, tabel)
}

// sqlPerbaruiTahunTreaty - SELURUH kolom yang dikirim diperbarui (AC 8).
func sqlPerbaruiTahunTreaty(tabel string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET TREATYYEAR = :1, UNDERWRITINGYEAR = :2, TREATYGROUPID = :3, TREATYGROUPNAME = :4,
	       USERID = :5, TGLUPDATE = :6, PROPORTION = :7, STARTDATE = :8, ENDDATE = :9
	 WHERE ID = :10`, tabel)
}

// sqlCariDobelTahunTreaty - AC 73: (STARTDATE, ENDDATE, TREATYGROUPID) yang
// sudah dipakai baris LAIN.
//
// ⛔ `(:4 IS NULL OR ID <> :5)`, bukan `ID <> :4` telanjang: teks kosong
// adalah NULL di Oracle, dan `ID <> NULL` tidak pernah benar - baris baru
// tidak akan pernah menemukan dobelnya.
//
// ⚠️ 29-09-2026 (tiket 04): DUA placeholder berbeda, nilai yang sama diikat
// dua kali. Semula `:4` dipakai dua kali dengan satu argumen; benar tidaknya
// bergantung pada driver mengikat per nama atau per kemunculan, dan belum
// pernah dijalankan terhadap Oracle. Placeholder berbeda benar di keduanya.
//
// TRUNC di kedua sisi: baris warisan boleh membawa jam pada kolom DATE.
func sqlCariDobelTahunTreaty(tabel string) string {
	return fmt.Sprintf(`SELECT ID FROM %s
	 WHERE TREATYGROUPID = :1
	   AND TRUNC(STARTDATE) = TRUNC(:2)
	   AND TRUNC(ENDDATE) = TRUNC(:3)
	   AND (:4 IS NULL OR ID <> :5)
	 ORDER BY ID
	 FETCH FIRST 1 ROWS ONLY`, tabel)
}

// tanggalJadiNil membind waktu nol sebagai NULL (ADR-U-0027).
func tanggalJadiNil(t time.Time) any {
	if t.IsZero() {
		return nil
	}
	return t
}

// uraiTanggalTeks membaca TO_CHAR 'YYYY-MM-DD HH24:MI:SS'; NULL = kosong.
func uraiTanggalTeks(v sql.NullString, kolom string) (time.Time, error) {
	if !v.Valid || v.String == "" {
		return time.Time{}, nil
	}
	t, err := utils.ParseTanggal(v.String)
	if err != nil {
		return time.Time{}, fmt.Errorf("repository: kolom %s bernilai %q: %w", kolom, v.String, err)
	}
	return t, nil
}

// pindaiTahunTreaty membaca satu baris hasil daftarPilihTahunTreaty.
func pindaiTahunTreaty(baca interface{ Scan(...any) error }) (models.TahunTreaty, error) {
	var n [10]sql.NullString
	if err := baca.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6], &n[7], &n[8], &n[9]); err != nil {
		return models.TahunTreaty{}, err
	}
	t := models.TahunTreaty{
		ID: n[0].String, TreatyYear: n[1].String, UnderwritingYear: n[2].String,
		TreatyGroupID: n[3].String, TreatyGroupName: n[4].String, Proportion: n[5].String,
		UserID: n[8].String,
	}
	var err error
	if t.StartDate, err = uraiTanggalTeks(n[6], "STARTDATE"); err != nil {
		return t, err
	}
	if t.EndDate, err = uraiTanggalTeks(n[7], "ENDDATE"); err != nil {
		return t, err
	}
	if t.TglUpdate, err = uraiTanggalTeks(n[9], "TGLUPDATE"); err != nil {
		return t, err
	}
	return t, nil
}

// Daftar membaca satu halaman, terbaru dahulu.
func (m *MasterTahunTreaty) Daftar(ctx context.Context, halaman, ukuran int) (HalamanTahunTreaty, error) {
	tabel, err := m.db.Qualify(TabelTahunTCO)
	if err != nil {
		return HalamanTahunTreaty{}, err
	}
	if halaman < 1 {
		halaman = 1
	}
	if ukuran < 1 {
		ukuran = 20
	}
	hasil := HalamanTahunTreaty{Halaman: halaman, Ukuran: ukuran}
	qCacah := sqlCacahTahunTreaty(tabel)
	if err := PeriksaSQL(qCacah); err != nil {
		return hasil, err
	}
	if err := m.db.sql.QueryRowContext(ctx, qCacah).Scan(&hasil.Total); err != nil {
		return hasil, fmt.Errorf("repository: mencacah tahun treaty: %w", err)
	}
	q := sqlDaftarTahunTreaty(tabel)
	if err := PeriksaSQL(q); err != nil {
		return hasil, err
	}
	rows, err := m.db.sql.QueryContext(ctx, q, (halaman-1)*ukuran, ukuran)
	if err != nil {
		return hasil, fmt.Errorf("repository: membaca daftar tahun treaty: %w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		t, err := pindaiTahunTreaty(rows)
		if err != nil {
			return hasil, err
		}
		hasil.Baris = append(hasil.Baris, t)
	}
	return hasil, rows.Err()
}

// Ambil membaca satu tahun treaty.
func (m *MasterTahunTreaty) Ambil(ctx context.Context, id string) (models.TahunTreaty, error) {
	tabel, err := m.db.Qualify(TabelTahunTCO)
	if err != nil {
		return models.TahunTreaty{}, err
	}
	q := sqlAmbilTahunTreaty(tabel)
	if err := PeriksaSQL(q); err != nil {
		return models.TahunTreaty{}, err
	}
	t, err := pindaiTahunTreaty(m.db.sql.QueryRowContext(ctx, q, id))
	if errors.Is(err, sql.ErrNoRows) {
		return models.TahunTreaty{}, ErrTahunTreatyTidakAda
	}
	if err != nil {
		return models.TahunTreaty{}, fmt.Errorf("repository: membaca tahun treaty %s: %w", id, err)
	}
	return t, nil
}

// Sisip menulis tahun treaty BARU; identitas dari sequence (ADR-0006).
func (m *MasterTahunTreaty) Sisip(ctx context.Context, tx *Tx, t models.TahunTreaty) (string, error) {
	if tx == nil {
		return "", errors.New("repository: menyisipkan tahun treaty menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelTahunTCO)
	if err != nil {
		return "", err
	}
	id, err := m.db.IdentitasBerikutTCO(ctx, tx, SeqTahunTCO)
	if err != nil {
		return "", err
	}
	q := sqlSisipTahunTreaty(tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, id,
		kosongJadiNil(t.TreatyYear), kosongJadiNil(t.UnderwritingYear),
		kosongJadiNil(t.TreatyGroupID), kosongJadiNil(t.TreatyGroupName),
		kosongJadiNil(t.UserID), tanggalJadiNil(t.TglUpdate), kosongJadiNil(t.Proportion),
		tanggalJadiNil(t.StartDate), tanggalJadiNil(t.EndDate))
	if err != nil {
		return "", fmt.Errorf("repository: menyisipkan tahun treaty: %w", err)
	}
	return id, pastikanSatuBaris(hasil, "penyisipan tahun treaty")
}

// Perbarui menimpa SELURUH kolom tahun treaty ber-ID itu.
func (m *MasterTahunTreaty) Perbarui(ctx context.Context, tx *Tx, t models.TahunTreaty) error {
	if tx == nil {
		return errors.New("repository: memperbarui tahun treaty menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelTahunTCO)
	if err != nil {
		return err
	}
	q := sqlPerbaruiTahunTreaty(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q,
		kosongJadiNil(t.TreatyYear), kosongJadiNil(t.UnderwritingYear),
		kosongJadiNil(t.TreatyGroupID), kosongJadiNil(t.TreatyGroupName),
		kosongJadiNil(t.UserID), tanggalJadiNil(t.TglUpdate), kosongJadiNil(t.Proportion),
		tanggalJadiNil(t.StartDate), tanggalJadiNil(t.EndDate), t.ID)
	if err != nil {
		return fmt.Errorf("repository: memperbarui tahun treaty %s: %w", t.ID, err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: membaca cacah pembaruan tahun treaty: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("%w: %s", ErrTahunTreatyTidakAda, t.ID)
	}
	return nil
}

// CariDobel mengembalikan ID tahun treaty LAIN yang memakai kombinasi
// (grup, mulai, akhir) yang sama; teks kosong bila tidak ada (AC 73).
func (m *MasterTahunTreaty) CariDobel(ctx context.Context, tx *Tx, grupID string,
	mulai, akhir time.Time, kecualiID string) (string, error) {

	if tx == nil {
		return "", errors.New("repository: pemeriksaan dobel tahun treaty menuntut transaksi")
	}
	if mulai.IsZero() || akhir.IsZero() {
		return "", nil
	}
	tabel, err := m.db.Qualify(TabelTahunTCO)
	if err != nil {
		return "", err
	}
	q := sqlCariDobelTahunTreaty(tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	var id sql.NullString
	err = tx.tx.QueryRowContext(ctx, q, grupID, mulai, akhir, kosongJadiNil(kecualiID),
		kosongJadiNil(kecualiID)).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: memeriksa dobel tahun treaty: %w", err)
	}
	return id.String, nil
}
