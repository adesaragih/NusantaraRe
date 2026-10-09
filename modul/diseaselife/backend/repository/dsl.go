package repository

// SQL modul Disease Life atas `DISEASE_LIFE` (ID, ICD_CODE, DISEASE) - kolom bernama, nol JSON. Nol DELETE (XML
// `InboxDisease` tanpa Delete: `pyGridDeleteActivityExists` false b460 / b754 / b1929 / b3262 / b3619 / b5297).
//
// ⛔ 97.586 baris DEV: grid SELALU berhalaman dan bersaring DI SERVER (OFFSET / FETCH NEXT terikat, COUNT bersaring) -
// tidak ada jalan yang memuat seluruh tabel (uji TestSqlDaftarSelaluBerhalaman).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/diseaselife/backend/models"
)

var (
	// ErrTidakAda - ID tidak ada.
	ErrTidakAda = errors.New("repository: Disease tidak ada")
	// ErrBelumAda - tabel, kolom, atau sequence tidak ada di skema ini (ORA-00942 / ORA-00904 / ORA-02289).
	ErrBelumAda = errors.New("repository: tabel atau sequence Disease tidak ada di skema ini")
	// ErrKembar - ORA-00001 (PK_DISEASE_LIFE): ID sudah dipakai.
	ErrKembar = errors.New("repository: ID sudah dipakai")
	// ErrBacaSaja - SQL tulis diarahkan ke objek di luar DaftarTabelDitulis.
	ErrBacaSaja = errors.New("repository: objek ini dibaca saja")
)

// Gudang - akses Oracle modul ini.
type Gudang struct{ db *db.DB }

// Baru membuat gudang.
func Baru(d *db.DB) *Gudang { return &Gudang{db: d} }

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

// Bungkus memetakan galat Oracle ke galat paket ini (ORA-00001 = ErrKembar; tabel / kolom / sequence tidak ada =
// ErrBelumAda).
func Bungkus(err error, apa string) error {
	if err == nil {
		return nil
	}
	s := err.Error()
	switch {
	case strings.Contains(s, "ORA-00001"):
		return fmt.Errorf("%w: %v", ErrKembar, err)
	case strings.Contains(s, "ORA-00942") || strings.Contains(s, "ORA-00904") || strings.Contains(s, "ORA-02289"):
		return fmt.Errorf("%w: %v", ErrBelumAda, err)
	}
	return fmt.Errorf("repository: %s: %w", apa, err)
}

// PeriksaTulis - lapis penjaga: pernyataan bukan SELECT hanya boleh atas objek DaftarTabelDitulis.
func PeriksaTulis(objek, q string) error {
	if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(q)), "SELECT ") {
		return nil
	}
	for _, t := range DaftarTabelDitulis {
		if t == objek {
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrBacaSaja, objek)
}

func siap(objek, q string) error {
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	return PeriksaTulis(objek, q)
}

// nama - nama berskema tabel dan sequence.
func (g *Gudang) nama() (tabel, seq string, err error) {
	if tabel, err = g.db.Qualify(Tabel); err != nil {
		return "", "", err
	}
	if seq, err = g.db.Qualify(Seq); err != nil {
		return "", "", err
	}
	return tabel, seq, nil
}

// PolaCari - pola LIKE ber-ESCAPE '\' untuk saring "memuat", tanpa beda huruf; kosong = nil (tanpa saring).
func PolaCari(kata string) any {
	kata = strings.ToUpper(models.RapikanSaring(kata))
	if kata == "" {
		return nil
	}
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return "%" + r.Replace(kata) + "%"
}

// Urutan - klausa ORDER BY grid, selalu diakhiri ID supaya halaman stabil. Kolom ID (bawaan MENURUN: `pySortType`
// DESC b5012, `pySortOrder` 1 b5017) diurut sebagai ANGKA, bukan teks; kolom ICD Code dapat dipilih (`pyColumnSorting`
// true b5036). Kolom Disease TIDAK (b5058). Nilai `urut` lain = ID.
func Urutan(urut string, naik bool) string {
	arah := "DESC"
	if naik {
		arah = "ASC"
	}
	if urut == models.UrutICD {
		return fmt.Sprintf("ICD_CODE %s NULLS LAST, TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) %s NULLS LAST, ID %s", arah, arah, arah)
	}
	return fmt.Sprintf("TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) %s NULLS LAST, ID %s", arah, arah)
}

// saring - ICD Code dan Disease "memuat", tanpa beda huruf (`pyGridFiltering` true b5148). Penampung UNIK (penjaga
// penampung_statik: penampung berulang + OFFSET / FETCH = ORA-01008).
const saring = `(:1 IS NULL OR UPPER(ICD_CODE) LIKE :2 ESCAPE '\') AND (:3 IS NULL OR UPPER(DISEASE) LIKE :4 ESCAPE '\')`

// SqlDaftar - satu halaman grid `BrowseDiseaseLife_RD` (b5094): bersaring, berurut, OFFSET / FETCH NEXT terikat.
func SqlDaftar(t, urut string, naik bool) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE %s ORDER BY %s OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY`,
		kolomBaca, t, saring, Urutan(urut, naik))
}

// SqlJumlah - jumlah baris bersaring.
func SqlJumlah(t string) string { return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, t, saring) }

// SqlAmbil - satu baris.
func SqlAmbil(t string) string { return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, kolomBaca, t) }

// SqlPemakaiICD - baris lain ber-ICD Code sama (tanpa beda huruf dan spasi tepi), selain kecualiID (D3, di luar XML).
func SqlPemakaiICD(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE UPPER(TRIM(ICD_CODE)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0)) ORDER BY ID FETCH FIRST 5 ROWS ONLY`,
		kolomBaca, t)
}

// SqlNomorBaru - satu nomor SEQ_DISEASE_LIFE sebagai teks.
func SqlNomorBaru(seq string) string { return fmt.Sprintf(`SELECT TO_CHAR(%s.NEXTVAL) FROM DUAL`, seq) }

// SqlAdaID - ID sudah terpakai?
func SqlAdaID(t string) string { return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ID = :1`, t) }

// SqlSisip - Save Add (`AddToList_Act` b2121).
func SqlSisip(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, ICD_CODE, DISEASE) VALUES (:1, :2, :3)`, t)
}

// SqlUbah - Edit (`EditList_DT` b4852) lalu Save.
func SqlUbah(t string) string {
	return fmt.Sprintf(`UPDATE %s SET ICD_CODE = :1, DISEASE = :2 WHERE ID = :3`, t)
}

func pindai(rows *sql.Rows) ([]models.Penyakit, error) {
	var out []models.Penyakit
	for rows.Next() {
		var id, icd, nama sql.NullString
		if err := rows.Scan(&id, &icd, &nama); err != nil {
			return nil, err
		}
		out = append(out, models.Penyakit{ID: id.String, ICDCode: icd.String, Disease: nama.String})
	}
	return out, rows.Err()
}

func (g *Gudang) baca(ctx context.Context, tx *db.Tx, q string, args ...any) ([]models.Penyakit, error) {
	if err := siap(Tabel, q); err != nil {
		return nil, err
	}
	rows, err := g.dari(tx).QueryContext(ctx, q, args...)
	if err != nil {
		return nil, Bungkus(err, "membaca")
	}
	defer func() { _ = rows.Close() }()
	out, err := pindai(rows)
	return out, Bungkus(err, "membaca")
}

func (g *Gudang) satuNilai(ctx context.Context, tx *db.Tx, objek, q string, args ...any) (string, error) {
	if err := siap(objek, q); err != nil {
		return "", err
	}
	var s sql.NullString
	if err := g.dari(tx).QueryRowContext(ctx, q, args...).Scan(&s); err != nil {
		return "", Bungkus(err, "membaca")
	}
	return s.String, nil
}

func (g *Gudang) tulis(ctx context.Context, tx *db.Tx, q, apa string, args ...any) (int64, error) {
	if err := siap(Tabel, q); err != nil {
		return 0, err
	}
	h, err := g.dari(tx).ExecContext(ctx, q, args...)
	if err != nil {
		return 0, Bungkus(err, apa)
	}
	n, _ := h.RowsAffected()
	return n, nil
}

func angka(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}

// ArgSaring - nilai keempat penampung saring (:1-:4), urutan sama dengan `saring`.
func ArgSaring(s models.Saringan) []any {
	icd, nama := PolaCari(s.ICDCode), PolaCari(s.Disease)
	return []any{icd, icd, nama, nama}
}

// Daftar - satu halaman dan total baris bersaring.
func (g *Gudang) Daftar(ctx context.Context, s models.Saringan) ([]models.Penyakit, int, error) {
	t, _, err := g.nama()
	if err != nil {
		return nil, 0, err
	}
	arg := ArgSaring(s)
	total, err := g.satuNilai(ctx, nil, Tabel, SqlJumlah(t), arg...)
	if err != nil {
		return nil, 0, err
	}
	d, err := g.baca(ctx, nil, SqlDaftar(t, s.Urut, s.Naik), append(arg, (s.Halaman-1)*models.UkuranHalaman, models.UkuranHalaman)...)
	return d, angka(total), err
}

// Ambil - satu baris; ErrTidakAda bila ID tidak ada.
func (g *Gudang) Ambil(ctx context.Context, tx *db.Tx, id string) (models.Penyakit, error) {
	t, _, err := g.nama()
	if err != nil {
		return models.Penyakit{}, err
	}
	d, err := g.baca(ctx, tx, SqlAmbil(t), id)
	if err != nil {
		return models.Penyakit{}, err
	}
	if len(d) == 0 {
		return models.Penyakit{}, ErrTidakAda
	}
	return d[0], nil
}

// PemakaiICD - baris ber-ICD Code sama selain kecualiID ("" = semua).
func (g *Gudang) PemakaiICD(ctx context.Context, tx *db.Tx, icd, kecualiID string) ([]models.Penyakit, error) {
	t, _, err := g.nama()
	if err != nil {
		return nil, err
	}
	return g.baca(ctx, tx, SqlPemakaiICD(t), icd, db.KosongJadiNil(kecualiID))
}

// NomorBaru - satu nomor `SEQ_DISEASE_LIFE` (teks).
func (g *Gudang) NomorBaru(ctx context.Context, tx *db.Tx) (string, error) {
	_, seq, err := g.nama()
	if err != nil {
		return "", err
	}
	return g.satuNilai(ctx, tx, "DUAL", SqlNomorBaru(seq))
}

// AdaID - ID sudah terpakai di DISEASE_LIFE?
func (g *Gudang) AdaID(ctx context.Context, tx *db.Tx, id string) (bool, error) {
	t, _, err := g.nama()
	if err != nil {
		return false, err
	}
	s, err := g.satuNilai(ctx, tx, Tabel, SqlAdaID(t), id)
	return angka(s) > 0, err
}

// Sisip - baris baru; ErrKembar bila ID terpakai (ORA-00001, PK_DISEASE_LIFE).
func (g *Gudang) Sisip(ctx context.Context, tx *db.Tx, p models.Penyakit) error {
	t, _, err := g.nama()
	if err != nil {
		return err
	}
	_, err = g.tulis(ctx, tx, SqlSisip(t), "menyimpan", p.ID, p.ICDCode, p.Disease)
	return err
}

// Ubah - ICD_CODE dan DISEASE satu baris; ErrTidakAda bila ID tidak ada.
func (g *Gudang) Ubah(ctx context.Context, tx *db.Tx, p models.Penyakit) error {
	t, _, err := g.nama()
	if err != nil {
		return err
	}
	n, err := g.tulis(ctx, tx, SqlUbah(t), "mengubah", p.ICDCode, p.Disease, p.ID)
	if err == nil && n == 0 {
		return ErrTidakAda
	}
	return err
}
