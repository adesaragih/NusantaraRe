package repository

// Business pada kombinasi (tahun, grup, jenis) - tiket 07 Treaty Contract Out.
//
// Untuk apa berkas ini: baca/tulis/hapus `TREATYBUSINESS` (meniru logika
// `POOLDATA.PEGA_TREATYBUSINESS`, 12 parameter, TANPA memanggilnya) dan pembaca
// master `BUSINESS` (dibaca saja).
//
// ⛔ AC 23 - PERBAIKAN SADAR atas prosedur warisan: pembaruan menulis SELURUH
// medan yang dikirim, bukan hanya ISACTIVE/BIZCODE/BIZNAME/USERID/TGLUPDATE
// seperti `PEGA_TREATYBUSINESS` `[data DBA]`.
//
// ⛔ AC 63/64 - hapus menyentuh SATU tabel. `RDBList/DeleteRowBusinessList.xml`
// memegang dua SQL (`treatybusiness` di `pyBrowseSQL`, tabel dokumen kembarnya di
// `pyDeleteSQL`); yang kedua tidak dibawa (penyimpangan sadar 1).
//
// Dibaca sesudah: tco_reinsurer.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatycontractout/models"
)

// MasterBusinessTCO - master bisnis, dibaca saja.
//
// `[terverifikasi]` pemilih `Business Name` `ViewDetailTreatyBusinessGrid.xml`
// b6294-b6328 -> `BrowseFilterBusiness_RD` (kelas `ASM-FW-GISFW-Int-BUSINESS`,
// INNER JOIN `BUSINESSGROUP` pada `.BusinessGroupID = A.ID`, urut `.Note ASC`
// b861-b867): `.ID` -> BizCode, `.Note` -> BIZNAME. Kolom fisik `ID`, `NOTE`,
// `BUSINESSGROUPID` terbukti di SQL korpus (`from business where ID = ...`).
const MasterBusinessTCO = "BUSINESS"

// batasBusinessMasterTCO - `pyMaxRecords` bawaan RD 500 (sejajar BrowseTreatyYear_RD b757).
const batasBusinessMasterTCO = 500

// ErrBusinessTidakAda - baris bisnis bukan milik kombinasi itu.
var ErrBusinessTidakAda = errors.New("repository: business row not found in this combination")

// ErrBusinessMasterTidakAda - kode bisnis tidak ada di master.
var ErrBusinessMasterTidakAda = errors.New("repository: business code not in the BUSINESS master")

// BusinessMasterTCO adalah satu pilihan bisnis.
type BusinessMasterTCO struct{ ID, Note string }

// MasterBusiness membaca master `BUSINESS`.
type MasterBusiness struct{ db *db.DB }

// NewMasterBusiness menyusun pembacanya.
func NewMasterBusiness(db *db.DB) *MasterBusiness { return &MasterBusiness{db: db} }

// sqlDaftarBusinessMasterTCO - saringan `Note`/`ID`/`OLDID`/`GRUPBIS` RD tidak
// berparameter dari layar (dilewati bila kosong), dan parameter `TREATYGROUP`
// yang dikirim layar (b6363) TIDAK dipakai saringan mana pun - jadi pilihan
// TIDAK tersaring grup treaty.
//
// ⚠️ `INNER JOIN BUSINESSGROUP` diganti `BUSINESSGROUPID IS NOT NULL`
// `[keputusan kami]`: tabel `BUSINESSGROUP` tidak terbukti di SQL korpus;
// bedanya hanya baris berkelompok yatim (OQ-TCO-13).
func sqlDaftarBusinessMasterTCO(tabel string) string {
	return fmt.Sprintf(`SELECT ID, NOTE FROM %s WHERE BUSINESSGROUPID IS NOT NULL
	 ORDER BY NOTE ASC, ID ASC FETCH FIRST %d ROWS ONLY`, tabel, batasBusinessMasterTCO)
}

func sqlAmbilBusinessMasterTCO(tabel string) string {
	return fmt.Sprintf(`SELECT ID, NOTE FROM %s WHERE ID = :1 AND BUSINESSGROUPID IS NOT NULL`, tabel)
}

// Daftar membaca seluruh pilihan bisnis.
func (m *MasterBusiness) Daftar(ctx context.Context) ([]BusinessMasterTCO, error) {
	tabel, err := m.db.Qualify(MasterBusinessTCO)
	if err != nil {
		return nil, err
	}
	q := sqlDaftarBusinessMasterTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := bacaTCO(ctx, m.db).QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: reading master %s: %w", MasterBusinessTCO, err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []BusinessMasterTCO
	for rows.Next() {
		var id, note sql.NullString
		if err := rows.Scan(&id, &note); err != nil {
			return nil, err
		}
		hasil = append(hasil, BusinessMasterTCO{ID: id.String, Note: note.String})
	}
	return hasil, rows.Err()
}

// Ambil membaca satu pilihan bisnis.
func (m *MasterBusiness) Ambil(ctx context.Context, id string) (BusinessMasterTCO, error) {
	tabel, err := m.db.Qualify(MasterBusinessTCO)
	if err != nil {
		return BusinessMasterTCO{}, err
	}
	q := sqlAmbilBusinessMasterTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return BusinessMasterTCO{}, err
	}
	var gotID, note sql.NullString
	err = bacaTCO(ctx, m.db).QueryRowContext(ctx, q, id).Scan(&gotID, &note)
	if errors.Is(err, sql.ErrNoRows) {
		return BusinessMasterTCO{}, ErrBusinessMasterTidakAda
	}
	if err != nil {
		return BusinessMasterTCO{}, fmt.Errorf("repository: reading master %s %s: %w", MasterBusinessTCO, id, err)
	}
	return BusinessMasterTCO{ID: gotID.String, Note: note.String}, nil
}

// MasterBusinessKombinasiTCO membaca dan menulis `TREATYBUSINESS`.
type MasterBusinessKombinasiTCO struct{ db *db.DB }

// NewMasterBusinessKombinasiTCO menyusun gudangnya.
func NewMasterBusinessKombinasiTCO(db *db.DB) *MasterBusinessKombinasiTCO {
	return &MasterBusinessKombinasiTCO{db: db}
}

// pilihBusinessTCO - tco4: seluruh kolom `TREATYBUSINESS` VARCHAR2 `[data DBA]`.
const pilihBusinessTCO = `ID, ISACTIVE, TREATYYEAR, TREATYYEARID, TREATYGROUPID, TREATYGROUPNAME, REINSTYPEID,
	       REINSTYPENAME, BIZCODE, BIZNAME, USERID, TGLUPDATE`

// sqlDaftarBusinessTCO - SELURUH baris kombinasi, aktif maupun nonaktif (AC 22).
//
// ⚠️ Grid Pega (`BrowseTreatyBusiness_RD` b670-b675) menyaring `.IsActive = 1`,
// sehingga baris nonaktif lenyap dari layar. AC 22 menuntut baris nonaktif
// TETAP terbaca - ralat bertanggal tiket 07.
func sqlDaftarBusinessTCO(tabel string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE %s ORDER BY ID ASC`, pilihBusinessTCO, tabel, saringKombinasiTCO)
}

func sqlAmbilBusinessTCO(tabel string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE %s AND ID = :4`, pilihBusinessTCO, tabel, saringKombinasiTCO)
}

func sqlSisipBusinessTCO(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s
	       (ID, ISACTIVE, TREATYYEAR, TREATYYEARID, TREATYGROUPID, TREATYGROUPNAME, REINSTYPEID, REINSTYPENAME,
	        BIZCODE, BIZNAME, USERID, TGLUPDATE)
	VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12)`, tabel)
}

// sqlPerbaruiBusinessTCO - PERSIS UPDATE `PEGA_TREATYBUSINESS` `[data DBA]`:
// hanya `ISACTIVE, BIZCODE, BIZNAME, USERID, TGLUPDATE`, dibatasi kombinasi.
// RALAT tco4 atas AC 23 ("seluruh medan non-kunci"): procedure tidak menyentuh
// `TREATYYEARID`/`TREATYGROUPNAME`/`REINSTYPENAME` saat memperbarui.
func sqlPerbaruiBusinessTCO(tabel string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET ISACTIVE = :1, BIZCODE = :2, BIZNAME = :3, USERID = :4, TGLUPDATE = :5
	 WHERE ID = :6 AND TREATYYEAR = :7 AND TREATYGROUPID = :8 AND REINSTYPEID = :9`, tabel)
}

// sqlHapusBusinessTCO - SATU tabel (AC 63/64), dibatasi kombinasi.
func sqlHapusBusinessTCO(tabel string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE ID = :1 AND TREATYYEAR = :2 AND TREATYGROUPID = :3 AND REINSTYPEID = :4`, tabel)
}

// sqlCariDobelBusinessTCO - kode bisnis yang sama pada kombinasi yang sama.
func sqlCariDobelBusinessTCO(tabel string) string {
	return fmt.Sprintf(`SELECT ID FROM %s WHERE %s AND BIZCODE = :4 AND (:5 IS NULL OR ID <> :6)
	 ORDER BY ID FETCH FIRST 1 ROWS ONLY`, tabel, saringKombinasiTCO)
}

func pindaiBusinessTCO(baca interface{ Scan(...any) error }) (models.BusinessTreaty, error) {
	var n [12]sql.NullString
	tujuan := make([]any, len(n))
	for i := range n {
		tujuan[i] = &n[i]
	}
	if err := baca.Scan(tujuan...); err != nil {
		return models.BusinessTreaty{}, err
	}
	b := models.BusinessTreaty{ID: n[0].String, IsActive: n[1].String, TreatyYear: n[2].String,
		TreatyYearID: n[3].String, TreatyGroupID: n[4].String, TreatyGroupName: n[5].String,
		ReinsTypeID: n[6].String, ReinsTypeName: n[7].String, BizCode: n[8].String, BizName: n[9].String,
		UserID: n[10].String}
	var err error
	b.TglUpdate, err = waktuWarisanTeks(n[11], "TGLUPDATE")
	return b, err
}

// Daftar membaca seluruh baris bisnis satu kombinasi.
func (m *MasterBusinessKombinasiTCO) Daftar(ctx context.Context, k models.KombinasiTCO) ([]models.BusinessTreaty, error) {
	tabel, err := m.db.Qualify(TabelBusinessTCO)
	if err != nil {
		return nil, err
	}
	q := sqlDaftarBusinessTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := bacaTCO(ctx, m.db).QueryContext(ctx, q, argKombinasi(k)...)
	if err != nil {
		return nil, fmt.Errorf("repository: reading combination business rows: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []models.BusinessTreaty
	for rows.Next() {
		b, err := pindaiBusinessTCO(rows)
		if err != nil {
			return nil, err
		}
		hasil = append(hasil, b)
	}
	return hasil, rows.Err()
}

// Ambil membaca satu baris bisnis milik kombinasi itu.
func (m *MasterBusinessKombinasiTCO) Ambil(ctx context.Context, k models.KombinasiTCO, id string) (models.BusinessTreaty, error) {
	tabel, err := m.db.Qualify(TabelBusinessTCO)
	if err != nil {
		return models.BusinessTreaty{}, err
	}
	q := sqlAmbilBusinessTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return models.BusinessTreaty{}, err
	}
	b, err := pindaiBusinessTCO(bacaTCO(ctx, m.db).QueryRowContext(ctx, q, append(argKombinasi(k), id)...))
	if errors.Is(err, sql.ErrNoRows) {
		return models.BusinessTreaty{}, ErrBusinessTidakAda
	}
	if err != nil {
		return models.BusinessTreaty{}, fmt.Errorf("repository: reading business %s: %w", id, err)
	}
	return b, nil
}

// CariDobel mencari baris LAIN berkode bisnis sama pada kombinasi.
func (m *MasterBusinessKombinasiTCO) CariDobel(ctx context.Context, tx *db.Tx, k models.KombinasiTCO, bizCode, kecualiID string) (string, error) {
	if tx == nil {
		return "", errors.New("repository: business duplicate search requires a transaction")
	}
	tabel, err := m.db.Qualify(TabelBusinessTCO)
	if err != nil {
		return "", err
	}
	q := sqlCariDobelBusinessTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var id string
	err = tx.QueryRowContext(ctx, q, append(argKombinasi(k), bizCode, db.KosongJadiNil(kecualiID),
		db.KosongJadiNil(kecualiID))...).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: searching duplicate business: %w", err)
	}
	return id, nil
}

// Sisip menulis baris bisnis baru; ID dari `TREATY_BUSINESS_SEQ`.
func (m *MasterBusinessKombinasiTCO) Sisip(ctx context.Context, tx *db.Tx, b models.BusinessTreaty) (string, error) {
	if tx == nil {
		return "", errors.New("repository: inserting business requires a transaction")
	}
	tabel, err := m.db.Qualify(TabelBusinessTCO)
	if err != nil {
		return "", err
	}
	id, err := IdentitasBerikutTCO(ctx, m.db, tx, SeqBusinessTCO)
	if err != nil {
		return "", err
	}
	q := sqlSisipBusinessTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	hasil, err := tx.ExecContext(ctx, q, id, db.KosongJadiNil(b.IsActive), db.KosongJadiNil(b.TreatyYear),
		db.KosongJadiNil(b.TreatyYearID), db.KosongJadiNil(b.TreatyGroupID), db.KosongJadiNil(b.TreatyGroupName),
		db.KosongJadiNil(b.ReinsTypeID), db.KosongJadiNil(b.ReinsTypeName), db.KosongJadiNil(b.BizCode),
		db.KosongJadiNil(b.BizName), db.KosongJadiNil(b.UserID), db.KosongJadiNil(StempelPegaTCO(b.TglUpdate)))
	if err != nil {
		return "", fmt.Errorf("repository: inserting business: %w", err)
	}
	return id, db.PastikanSatuBaris(hasil, "business insert")
}

// Perbarui menimpa SELURUH medan non-kunci baris bisnis (AC 23).
func (m *MasterBusinessKombinasiTCO) Perbarui(ctx context.Context, tx *db.Tx, b models.BusinessTreaty) error {
	if tx == nil {
		return errors.New("repository: updating business requires a transaction")
	}
	tabel, err := m.db.Qualify(TabelBusinessTCO)
	if err != nil {
		return err
	}
	q := sqlPerbaruiBusinessTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, db.KosongJadiNil(b.IsActive), db.KosongJadiNil(b.BizCode),
		db.KosongJadiNil(b.BizName), db.KosongJadiNil(b.UserID), db.KosongJadiNil(StempelPegaTCO(b.TglUpdate)),
		b.ID, b.TreatyYear, b.TreatyGroupID, b.ReinsTypeID)
	if err != nil {
		return fmt.Errorf("repository: updating business %s: %w", b.ID, err)
	}
	if n, err := hasil.RowsAffected(); err == nil && n == 0 {
		return ErrBusinessTidakAda
	}
	return db.PastikanSatuBaris(hasil, "business update")
}

// Hapus membuang SATU baris bisnis dari SATU tabel.
func (m *MasterBusinessKombinasiTCO) Hapus(ctx context.Context, tx *db.Tx, k models.KombinasiTCO, id string) error {
	if tx == nil {
		return errors.New("repository: deleting business requires a transaction")
	}
	tabel, err := m.db.Qualify(TabelBusinessTCO)
	if err != nil {
		return err
	}
	q := sqlHapusBusinessTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, id, k.TreatyYear, k.TreatyGroupID, k.ReinsTypeID)
	if err != nil {
		return fmt.Errorf("repository: deleting business %s: %w", id, err)
	}
	if n, err := hasil.RowsAffected(); err == nil && n == 0 {
		return ErrBusinessTidakAda
	}
	return db.PastikanSatuBaris(hasil, "business delete")
}
