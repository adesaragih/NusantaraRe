package repository

// Kontrak treaty di dalam tahun treaty - tiket 04 Treaty Contract Out.
//
// Untuk apa berkas ini: baca/tulis tabel WARISAN `TREATYCONTRACT` (tco4), meniru logika
// `POOLDATA.PEGA_TREATYCONTRACT` (`RDBList/SaveMasterTreatyContract_SQL.xml`,
// 8 parameter + 2 keluaran) TANPA memanggil prosedurnya (prinsip o).
//
// `[data DBA]` upsert dikunci `ID`; identitas `'1' || lpad(seq, 6, '0')`;
// tanggal `to_date(..., 'DD/MM/YYYY')`; tidak COMMIT sendiri.
//
// ⛔ Anak-anak kontrak (reinsurer, business, klausul) menggantung pada
// KOMBINASI (TreatyYear, TreatyGroupID, ReinsTypeID), bukan pada `ID` kontrak
// [fakta bisnis work owner]. Kontrak membuka kombinasi itu; berkas ini tidak
// menulis anak apa pun.
//
// Dibaca sesudah: tco_tahun.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/internal/models"
)

// ErrKontrakTidakAda - kontrak tidak ada, atau bukan milik tahun treaty itu.
var ErrKontrakTidakAda = errors.New("repository: kontrak treaty tidak ditemukan pada tahun treaty ini")

// MasterKontrakTCO membaca dan menulis `TREATYCONTRACT`.
//
// ⚠️ tco4: `TREATYSTARTDATE`/`TREATYENDDATE` DATE (procedure `to_date(…,
// 'DD/MM/YYYY')`); `TGLUPDATE` VARCHAR2(1000) berisi stempel Pega.
type MasterKontrakTCO struct{ db *DB }

// NewMasterKontrakTCO menyusun gudangnya.
func NewMasterKontrakTCO(db *DB) *MasterKontrakTCO { return &MasterKontrakTCO{db: db} }

const pilihKontrakTCO = `ID, IDTREATYYEAR, REINSTYPEID, REINSTYPENAME,
	       TO_CHAR(TREATYSTARTDATE, 'YYYY-MM-DD HH24:MI:SS'), TO_CHAR(TREATYENDDATE, 'YYYY-MM-DD HH24:MI:SS'),
	       USERID, TGLUPDATE`

// sqlDaftarKontrakTCO - grid `BrowseTreatyContract_RD` (rule-nya TIDAK diekspor
// korpus) disaring `InputData.HASIL13 = Param.IDTreatyYear`
// (`BrowseReinsTypeYear.xml` b273). Urutan `ID DESC` `[keputusan kami]`,
// sejajar `BrowseTreatyYear_RD` b672.
func sqlDaftarKontrakTCO(tabel string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE IDTREATYYEAR = :1 ORDER BY ID DESC`, pilihKontrakTCO, tabel)
}

func sqlAmbilKontrakTCO(tabel string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE IDTREATYYEAR = :1 AND ID = :2`, pilihKontrakTCO, tabel)
}

func sqlSisipKontrakTCO(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s
	       (ID, IDTREATYYEAR, REINSTYPEID, REINSTYPENAME, TREATYSTARTDATE, TREATYENDDATE, USERID, TGLUPDATE)
	VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`, tabel)
}

// sqlPerbaruiKontrakTCO menimpa SELURUH medan (AC 8), dibatasi tahunnya.
func sqlPerbaruiKontrakTCO(tabel string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET REINSTYPEID = :1, REINSTYPENAME = :2, TREATYSTARTDATE = :3, TREATYENDDATE = :4,
	       USERID = :5, TGLUPDATE = :6
	 WHERE ID = :7 AND IDTREATYYEAR = :8`, tabel)
}

// sqlCariDobelKontrakTCO - satu jenis reasuransi satu kontrak per tahun.
//
// ⚠️ Placeholder berbeda untuk nilai yang sama (`:3`, `:4`) - lihat
// `sqlCariDobelTahunTreaty`.
func sqlCariDobelKontrakTCO(tabel string) string {
	return fmt.Sprintf(`SELECT ID FROM %s
	 WHERE IDTREATYYEAR = :1
	   AND REINSTYPEID = :2
	   AND (:3 IS NULL OR ID <> :4)
	 ORDER BY ID
	 FETCH FIRST 1 ROWS ONLY`, tabel)
}

func pindaiKontrakTCO(baca interface{ Scan(...any) error }) (models.KontrakTreaty, error) {
	var n [8]sql.NullString
	if err := baca.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6], &n[7]); err != nil {
		return models.KontrakTreaty{}, err
	}
	k := models.KontrakTreaty{ID: n[0].String, IDTreatyYear: n[1].String, ReinsTypeID: n[2].String,
		ReinsTypeName: n[3].String, UserID: n[6].String}
	var err error
	if k.TreatyStartDate, err = uraiTanggalTeks(n[4], "TREATYSTARTDATE"); err != nil {
		return k, err
	}
	if k.TreatyEndDate, err = uraiTanggalTeks(n[5], "TREATYENDDATE"); err != nil {
		return k, err
	}
	k.TglUpdate, err = waktuWarisanTeks(n[7], "TGLUPDATE")
	return k, err
}

// Daftar membaca kontrak satu tahun treaty.
func (m *MasterKontrakTCO) Daftar(ctx context.Context, tahunID string) ([]models.KontrakTreaty, error) {
	tabel, err := m.db.Qualify(TabelKontrakTCO)
	if err != nil {
		return nil, err
	}
	q := sqlDaftarKontrakTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := m.db.bacaTCO(ctx).QueryContext(ctx, q, tahunID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca kontrak tahun treaty %s: %w", tahunID, err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []models.KontrakTreaty
	for rows.Next() {
		k, err := pindaiKontrakTCO(rows)
		if err != nil {
			return nil, err
		}
		hasil = append(hasil, k)
	}
	return hasil, rows.Err()
}

// Ambil membaca SATU kontrak milik tahun treaty itu.
func (m *MasterKontrakTCO) Ambil(ctx context.Context, tahunID, id string) (models.KontrakTreaty, error) {
	tabel, err := m.db.Qualify(TabelKontrakTCO)
	if err != nil {
		return models.KontrakTreaty{}, err
	}
	q := sqlAmbilKontrakTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return models.KontrakTreaty{}, err
	}
	k, err := pindaiKontrakTCO(m.db.bacaTCO(ctx).QueryRowContext(ctx, q, tahunID, id))
	if errors.Is(err, sql.ErrNoRows) {
		return models.KontrakTreaty{}, ErrKontrakTidakAda
	}
	if err != nil {
		return models.KontrakTreaty{}, fmt.Errorf("repository: membaca kontrak %s: %w", id, err)
	}
	return k, nil
}

// Sisip menulis kontrak baru; ID dari `TREATYCONTRACT_SEQ` warisan (ADR-0006).
func (m *MasterKontrakTCO) Sisip(ctx context.Context, tx *Tx, k models.KontrakTreaty) (string, error) {
	if tx == nil {
		return "", errors.New("repository: menyisipkan kontrak menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelKontrakTCO)
	if err != nil {
		return "", err
	}
	id, err := m.db.IdentitasBerikutTCO(ctx, tx, SeqKontrakTCO)
	if err != nil {
		return "", err
	}
	q := sqlSisipKontrakTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, id, k.IDTreatyYear, kosongJadiNil(k.ReinsTypeID),
		kosongJadiNil(k.ReinsTypeName), tanggalJadiNil(k.TreatyStartDate), tanggalJadiNil(k.TreatyEndDate),
		kosongJadiNil(k.UserID), kosongJadiNil(StempelPegaTCO(k.TglUpdate)))
	if err != nil {
		return "", fmt.Errorf("repository: menyisipkan kontrak: %w", err)
	}
	return id, pastikanSatuBaris(hasil, "penyisipan kontrak")
}

// Perbarui menimpa kontrak yang ada; nol baris = kontrak bukan milik tahun itu.
func (m *MasterKontrakTCO) Perbarui(ctx context.Context, tx *Tx, k models.KontrakTreaty) error {
	if tx == nil {
		return errors.New("repository: memperbarui kontrak menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelKontrakTCO)
	if err != nil {
		return err
	}
	q := sqlPerbaruiKontrakTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, kosongJadiNil(k.ReinsTypeID), kosongJadiNil(k.ReinsTypeName),
		tanggalJadiNil(k.TreatyStartDate), tanggalJadiNil(k.TreatyEndDate), kosongJadiNil(k.UserID),
		kosongJadiNil(StempelPegaTCO(k.TglUpdate)), k.ID, k.IDTreatyYear)
	if err != nil {
		return fmt.Errorf("repository: memperbarui kontrak %s: %w", k.ID, err)
	}
	if n, err := hasil.RowsAffected(); err == nil && n == 0 {
		return ErrKontrakTidakAda
	}
	return pastikanSatuBaris(hasil, "pembaruan kontrak")
}

// CariDobel mencari kontrak LAIN di tahun yang sama berjenis reasuransi sama.
//
// ⚠️ Di DALAM transaksi penulisnya, seperti `MasterTahunTreaty.CariDobel`.
// Bukan unique index: data warisan boleh sudah berduplikat.
func (m *MasterKontrakTCO) CariDobel(ctx context.Context, tx *Tx, tahunID, reinsTypeID, kecualiID string) (string, error) {
	if tx == nil {
		return "", errors.New("repository: pencarian dobel kontrak menuntut transaksi")
	}
	// Temuan /code-review: tahun induk DIKUNCI dulu - dua penulis kontrak
	// serentak pada tahun yang sama tidak boleh sama-sama lolos pemeriksaan.
	tahun, err := m.db.Qualify(TabelTahunTCO)
	if err != nil {
		return "", err
	}
	kunci := sqlKunciTahunTCO(tahun)
	if err := PeriksaSQL(kunci); err != nil {
		return "", err
	}
	var terkunci string
	if err := tx.tx.QueryRowContext(ctx, kunci, tahunID).Scan(&terkunci); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrTahunTreatyTidakAda
		}
		return "", fmt.Errorf("repository: mengunci tahun treaty %s: %w", tahunID, err)
	}
	tabel, err := m.db.Qualify(TabelKontrakTCO)
	if err != nil {
		return "", err
	}
	q := sqlCariDobelKontrakTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	var id string
	err = tx.tx.QueryRowContext(ctx, q, tahunID, reinsTypeID, kosongJadiNil(kecualiID),
		kosongJadiNil(kecualiID)).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: mencari kontrak dobel: %w", err)
	}
	return id, nil
}

func sqlKunciKontrakTCO(tabel string) string {
	return fmt.Sprintf(`SELECT ID FROM %s WHERE IDTREATYYEAR = :1 AND ID = :2 FOR UPDATE`, tabel)
}

// Kunci mengunci baris kontrak selama transaksi penulis anak-anaknya.
//
// ⛔ Tiket 05: penulisan reinsurer MENJUMLAHKAN share kombinasi lalu menulis;
// dua penulis serentak yang sama-sama melihat total 60 akan sama-sama
// menambah 40. Mengunci kontrak pembuka kombinasinya menjadikan keduanya
// berurutan. `FOR UPDATE` atas baris reinsurer saja tidak cukup: baris BARU
// milik penulis lain tidak terkunci olehnya.
func (m *MasterKontrakTCO) Kunci(ctx context.Context, tx *Tx, tahunID, id string) error {
	if tx == nil {
		return errors.New("repository: mengunci kontrak menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelKontrakTCO)
	if err != nil {
		return err
	}
	q := sqlKunciKontrakTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	var got string
	err = tx.tx.QueryRowContext(ctx, q, tahunID, id).Scan(&got)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrKontrakTidakAda
	}
	if err != nil {
		return fmt.Errorf("repository: mengunci kontrak %s: %w", id, err)
	}
	return nil
}

// JumlahAnakKombinasi menghitung reinsurer + business yang menggantung pada
// kombinasi kontrak itu - saringan SAMA dengan kaskade hapus.
func (m *MasterKontrakTCO) JumlahAnakKombinasi(ctx context.Context, tx *Tx, kom models.KombinasiTCO, tahunID string) (int64, error) {
	if tx == nil {
		return 0, errors.New("repository: menghitung anak kombinasi menuntut transaksi")
	}
	reas, err := m.db.Qualify(TabelReinsurerTCO)
	if err != nil {
		return 0, err
	}
	biz, err := m.db.Qualify(TabelBusinessTCO)
	if err != nil {
		return 0, err
	}
	var jumlah int64
	for _, h := range []struct {
		q    string
		args []any
	}{
		{sqlHitungTCO(reas, saringReinsurerKaskadeTCO), argKombinasi(kom)},
		{sqlHitungTCO(biz, saringBusinessKaskadeTCO), []any{kom.TreatyYear, tahunID, kom.TreatyGroupID, kom.ReinsTypeID}},
	} {
		if err := PeriksaSQL(h.q); err != nil {
			return 0, err
		}
		var n int64
		if err := tx.tx.QueryRowContext(ctx, h.q, h.args...).Scan(&n); err != nil {
			return 0, fmt.Errorf("repository: menghitung anak kombinasi: %w", err)
		}
		jumlah += n
	}
	return jumlah, nil
}
