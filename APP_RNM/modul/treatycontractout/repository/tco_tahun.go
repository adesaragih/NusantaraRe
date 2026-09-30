package repository

// Tahun treaty - tabel WARISAN `TREATYYEAR` (tiket 03 Treaty Contract Out; tco4).
//
// Padanan `RDBList/SaveMasterTreatyYear_SQL.xml` -> `PEGA_TREATYYEAR`
// (`[data DBA]` UPSERT dikunci ID; ID baru `'1' || lpad(TreatyYear_seq, 6)`;
// tidak COMMIT sendiri; tco4: tabel warisan `TREATYYEAR`, TGLUPDATE stempel
// Pega) dan
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
// ⛔ tco4: seluruh kolom `TREATYYEAR` VARCHAR2 `[data DBA]`. `STARTDATE`/
// `ENDDATE` ditulis dan dibaca `YYYYMMDD` (OQ-TCO-01 ditutup dari data DEV,
// 182/182 baris), `TGLUPDATE` stempel `@getCurrentTimeStamp()`.
// Setiap query menyebut skemanya lewat Qualify (ADR-U-0033), nol COMMIT
// (ADR-U-0029).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/treatycontractout/models"
)

// ErrTahunTreatyTidakAda - tidak ada baris dengan ID itu.
var ErrTahunTreatyTidakAda = errors.New("repository: treaty year not found")

// HalamanTahunTreaty adalah satu halaman daftar tahun treaty.
type HalamanTahunTreaty struct {
	Baris   []models.TahunTreaty
	Total   int
	Halaman int
	Ukuran  int
}

// MasterTahunTreaty membaca dan menulis `TREATYYEAR`.
type MasterTahunTreaty struct{ db *db.DB }

// NewMasterTahunTreaty menyusunnya.
func NewMasterTahunTreaty(db *db.DB) *MasterTahunTreaty { return &MasterTahunTreaty{db: db} }

// daftarPilihTahunTreaty adalah sepuluh kolom warisan, apa adanya (VARCHAR2).
const daftarPilihTahunTreaty = `ID, TREATYYEAR, UNDERWRITINGYEAR, TREATYGROUPID, TREATYGROUPNAME, PROPORTION,
	       STARTDATE, ENDDATE, USERID, TGLUPDATE`

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

// sqlCariDobelTahunTreaty - AC 73: kandidat baris LAIN segrup; tanggalnya
// dibandingkan di Go (tco4: `STARTDATE`/`ENDDATE` VARCHAR2 `YYYYMMDD`, `TRUNC`
// tidak berlaku atas teks; baris berbentuk lain dilewati, bukan ditebak).
//
// ⛔ `(:2 IS NULL OR ID <> :3)`, bukan `ID <> :2` telanjang: teks kosong
// adalah NULL di Oracle, dan `ID <> NULL` tidak pernah benar. Placeholder
// berbeda untuk nilai yang sama (tiket 04).
func sqlCariDobelTahunTreaty(tabel string) string {
	return fmt.Sprintf(`SELECT ID, STARTDATE, ENDDATE FROM %s
	 WHERE TREATYGROUPID = :1
	   AND (:2 IS NULL OR ID <> :3)
	 ORDER BY ID`, tabel)
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
		return time.Time{}, fmt.Errorf("repository: column %s has value %q: %w", kolom, v.String, err)
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
	if t.StartDate, err = tanggalTahunWarisanTeks(n[6], "STARTDATE"); err != nil {
		return t, err
	}
	if t.EndDate, err = tanggalTahunWarisanTeks(n[7], "ENDDATE"); err != nil {
		return t, err
	}
	if t.TglUpdate, err = waktuWarisanTeks(n[9], "TGLUPDATE"); err != nil {
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
	if err := db.PeriksaSQL(qCacah); err != nil {
		return hasil, err
	}
	if err := bacaTCO(ctx, m.db).QueryRowContext(ctx, qCacah).Scan(&hasil.Total); err != nil {
		return hasil, fmt.Errorf("repository: counting treaty years: %w", err)
	}
	q := sqlDaftarTahunTreaty(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return hasil, err
	}
	rows, err := bacaTCO(ctx, m.db).QueryContext(ctx, q, (halaman-1)*ukuran, ukuran)
	if err != nil {
		return hasil, fmt.Errorf("repository: reading treaty year list: %w", err)
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
	if err := db.PeriksaSQL(q); err != nil {
		return models.TahunTreaty{}, err
	}
	t, err := pindaiTahunTreaty(bacaTCO(ctx, m.db).QueryRowContext(ctx, q, id))
	if errors.Is(err, sql.ErrNoRows) {
		return models.TahunTreaty{}, ErrTahunTreatyTidakAda
	}
	if err != nil {
		return models.TahunTreaty{}, fmt.Errorf("repository: reading treaty year %s: %w", id, err)
	}
	return t, nil
}

// Sisip menulis tahun treaty BARU; identitas dari sequence (ADR-0006).
func (m *MasterTahunTreaty) Sisip(ctx context.Context, tx *db.Tx, t models.TahunTreaty) (string, error) {
	if tx == nil {
		return "", errors.New("repository: inserting treaty year requires a transaction")
	}
	tabel, err := m.db.Qualify(TabelTahunTCO)
	if err != nil {
		return "", err
	}
	id, err := IdentitasBerikutTCO(ctx, m.db, tx, SeqTahunTCO)
	if err != nil {
		return "", err
	}
	q := sqlSisipTahunTreaty(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	hasil, err := tx.ExecContext(ctx, q, id,
		db.KosongJadiNil(t.TreatyYear), db.KosongJadiNil(t.UnderwritingYear),
		db.KosongJadiNil(t.TreatyGroupID), db.KosongJadiNil(t.TreatyGroupName),
		db.KosongJadiNil(t.UserID), db.KosongJadiNil(StempelPegaTCO(t.TglUpdate)), db.KosongJadiNil(t.Proportion),
		db.KosongJadiNil(TanggalYYYYMMDDTCO(t.StartDate)), db.KosongJadiNil(TanggalYYYYMMDDTCO(t.EndDate)))
	if err != nil {
		return "", fmt.Errorf("repository: inserting treaty year: %w", err)
	}
	return id, db.PastikanSatuBaris(hasil, "treaty year insert")
}

// Perbarui menimpa SELURUH kolom tahun treaty ber-ID itu.
func (m *MasterTahunTreaty) Perbarui(ctx context.Context, tx *db.Tx, t models.TahunTreaty) error {
	if tx == nil {
		return errors.New("repository: updating treaty year requires a transaction")
	}
	tabel, err := m.db.Qualify(TabelTahunTCO)
	if err != nil {
		return err
	}
	q := sqlPerbaruiTahunTreaty(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q,
		db.KosongJadiNil(t.TreatyYear), db.KosongJadiNil(t.UnderwritingYear),
		db.KosongJadiNil(t.TreatyGroupID), db.KosongJadiNil(t.TreatyGroupName),
		db.KosongJadiNil(t.UserID), db.KosongJadiNil(StempelPegaTCO(t.TglUpdate)), db.KosongJadiNil(t.Proportion),
		db.KosongJadiNil(TanggalYYYYMMDDTCO(t.StartDate)), db.KosongJadiNil(TanggalYYYYMMDDTCO(t.EndDate)), t.ID)
	if err != nil {
		return fmt.Errorf("repository: updating treaty year %s: %w", t.ID, err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: counting updated treaty year rows: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("%w: %s", ErrTahunTreatyTidakAda, t.ID)
	}
	return nil
}

// CariDobel mengembalikan ID tahun treaty LAIN yang memakai kombinasi
// (grup, mulai, akhir) yang sama; teks kosong bila tidak ada (AC 73).
func (m *MasterTahunTreaty) CariDobel(ctx context.Context, tx *db.Tx, grupID string,
	mulai, akhir time.Time, kecualiID string) (string, error) {

	if tx == nil {
		return "", errors.New("repository: treaty year duplicate check requires a transaction")
	}
	// Temuan /code-review: periksa-lalu-sisip tanpa kunci membiarkan dua penulis
	// serentak sama-sama lolos. Tahun treaty tidak punya baris induk untuk
	// dikunci - tabelnya dikunci EXCLUSIVE sampai transaksi selesai (penulis
	// tahun jarang; pembaca tidak terhalang).
	if err := m.kunciTabelTahunTCO(ctx, tx); err != nil {
		return "", err
	}
	if mulai.IsZero() || akhir.IsZero() {
		return "", nil
	}
	tabel, err := m.db.Qualify(TabelTahunTCO)
	if err != nil {
		return "", err
	}
	q := sqlCariDobelTahunTreaty(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	rows, err := tx.QueryContext(ctx, q, grupID, db.KosongJadiNil(kecualiID), db.KosongJadiNil(kecualiID))
	if err != nil {
		return "", fmt.Errorf("repository: checking duplicate treaty year: %w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var id, awal, akhirTeks sql.NullString
		if err := rows.Scan(&id, &awal, &akhirTeks); err != nil {
			return "", err
		}
		// ⚠️ Baris warisan bertanggal tak terurai tidak dapat "sama": dilewati,
		// bukan menggagalkan penyimpanan tahun lain.
		a, okA := uraiYYYYMMDDTCO(awal.String)
		b, okB := uraiYYYYMMDDTCO(akhirTeks.String)
		if okA && okB && tanggalSamaTCO(a, mulai) && tanggalSamaTCO(b, akhir) {
			return id.String, nil
		}
	}
	return "", rows.Err()
}

// tanggalSamaTCO - tanggal kalender yang sama (jam diabaikan).
func tanggalSamaTCO(a, b time.Time) bool {
	ya, ma, da := a.Date()
	yb, mb, db := b.Date()
	return ya == yb && ma == mb && da == db
}

func sqlKunciTabelTahunTCO(tabel string) string {
	return fmt.Sprintf(`LOCK TABLE %s IN EXCLUSIVE MODE`, tabel)
}

func (m *MasterTahunTreaty) kunciTabelTahunTCO(ctx context.Context, tx *db.Tx) error {
	tabel, err := m.db.Qualify(TabelTahunTCO)
	if err != nil {
		return err
	}
	q := sqlKunciTabelTahunTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("repository: locking treaty year table: %w", err)
	}
	return nil
}

// sqlJumlahAnakTahunTCO - kontrak + klausul + LAMPIRAN (tco4: kunci lampiran
// `TREATYID` = teks `TREATYYEAR` + ID - mengubah TREATYYEAR memutusnya;
// temuan /code-review).
func sqlJumlahAnakTahunTCO(kontrak, klausul, lampiran, tahun string) string {
	return fmt.Sprintf(`SELECT (SELECT COUNT(*) FROM %s WHERE IDTREATYYEAR = :1) + (SELECT COUNT(*) FROM %s WHERE TREATYYEARID = :2)
	     + (SELECT COUNT(*) FROM %s a, %s y WHERE y.ID = :3 AND a.TREATYID = y.TREATYYEAR || y.ID) FROM DUAL`,
		kontrak, klausul, lampiran, tahun)
}

// JumlahAnak menghitung kontrak + klausul tahun itu - baris yang kombinasinya
// (TREATYYEAR teks, TREATYGROUPID) ikut ditulis dari tahun.
func (m *MasterTahunTreaty) JumlahAnak(ctx context.Context, tx *db.Tx, tahunID string) (int64, error) {
	if tx == nil {
		return 0, errors.New("repository: counting treaty year children requires a transaction")
	}
	kontrak, err := m.db.Qualify(TabelKontrakTCO)
	if err != nil {
		return 0, err
	}
	klausul, err := m.db.Qualify(TabelKlausulTCO)
	if err != nil {
		return 0, err
	}
	lampiran, err := m.db.Qualify(TabelLampiranTCO)
	if err != nil {
		return 0, err
	}
	tahun, err := m.db.Qualify(TabelTahunTCO)
	if err != nil {
		return 0, err
	}
	q := sqlJumlahAnakTahunTCO(kontrak, klausul, lampiran, tahun)
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	var n int64
	if err := tx.QueryRowContext(ctx, q, tahunID, tahunID, tahunID).Scan(&n); err != nil {
		return 0, fmt.Errorf("repository: counting children of treaty year %s: %w", tahunID, err)
	}
	return n, nil
}
