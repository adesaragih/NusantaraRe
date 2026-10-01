package login

// Gudang Oracle - M_LOGIN_GO, M_LOGIN_GO_WORKBASKET, dan master organisasi
// serta workbasket (dibaca, tidak dibuat).
//
// ⛔ Setiap tabel lewat `Qualify` (ADR-U-0033), nilai lewat bind, nol COMMIT
// di teks SQL (ADR-U-0029).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// Bendera '1'/'0' M_LOGIN_GO; IS_ACTIVE master berupa NUMBER(1).
const (
	benderaYa    = "1"
	benderaTidak = "0"
	masterAktif  = 1
	tabelLogin   = "M_LOGIN_GO"
	tabelLoginWB = "M_LOGIN_GO_WORKBASKET"
	tabelWB      = "M_WORKBASKET"
	tabelUnit    = "M_UNIT"
	tabelDivisi  = "M_DIVISION"
	tabelOrganis = "M_ORGANIZATION"
)

// GudangOracle menyimpan akun di Oracle.
type GudangOracle struct{ db *db.DB }

// NewGudangOracle menyusunnya.
func NewGudangOracle(d *db.DB) *GudangOracle { return &GudangOracle{db: d} }

// ⛔ `TERKUNCI` dihitung Oracle (`LOCKED_UNTIL > SYSDATE`), pasangan
// `sqlCatatGagal` yang mengunci dengan SYSDATE juga: satu jam saja.
func sqlAmbilAkun(t string) string {
	return fmt.Sprintf(`SELECT LOGIN_ID, NAME, PASSWORD_HASH, ORGANIZATION_CODE, DIVISION_CODE, UNIT_CODE,
	        IS_ACTIVE, MUST_CHANGE_PASSWORD, SESSION_VERSION,
	        CASE WHEN LOCKED_UNTIL > SYSDATE THEN 1 ELSE 0 END
	   FROM %s WHERE LOGIN_ID = :1`, t)
}

func sqlWorkbasket(login, wb string) string {
	return fmt.Sprintf(`SELECT l.WORKBASKET_ID FROM %s l JOIN %s w ON w.WORKBASKET_ID = l.WORKBASKET_ID
	  WHERE l.LOGIN_ID = :1 AND w.IS_ACTIVE = :2 ORDER BY l.WORKBASKET_ID`, login, wb)
}

// sqlCatatGagal - atomik di satu pernyataan, jadi dua percobaan serentak
// tetap terhitung dua. Kunci yang SUDAH lewat waktunya memulai hitungan dari
// satu lagi. Oracle membaca nilai LAMA di setiap ekspresi SET.
//
// :1 = BatasGagal, :2 = LamaKunci dalam menit, :3 = LOGIN_ID.
func sqlCatatGagal(t string) string {
	return fmt.Sprintf(`UPDATE %s SET
	    FAILED_COUNT = CASE WHEN LOCKED_UNTIL <= SYSDATE THEN 1 ELSE FAILED_COUNT + 1 END,
	    LOCKED_UNTIL = CASE
	        WHEN (CASE WHEN LOCKED_UNTIL <= SYSDATE THEN 1 ELSE FAILED_COUNT + 1 END) >= :1
	          THEN SYSDATE + :2 / 1440
	        WHEN LOCKED_UNTIL <= SYSDATE THEN NULL
	        ELSE LOCKED_UNTIL END,
	    TGL_UPDATE = SYSDATE
	  WHERE LOGIN_ID = :3`, t)
}

func sqlCatatBerhasil(t string) string {
	return fmt.Sprintf(`UPDATE %s SET FAILED_COUNT = 0, LOCKED_UNTIL = NULL, LAST_LOGIN = SYSDATE
	  WHERE LOGIN_ID = :1`, t)
}

// sqlGantiSandi - hanya bila versinya masih yang dibaca (`:4`): sesi yang
// sudah dicabut tidak dapat mengganti sandi.
func sqlGantiSandi(t string) string {
	return fmt.Sprintf(`UPDATE %s SET PASSWORD_HASH = :1, MUST_CHANGE_PASSWORD = :2,
	    SESSION_VERSION = SESSION_VERSION + 1, TGL_UPDATE = SYSDATE
	  WHERE LOGIN_ID = :3 AND SESSION_VERSION = :4`, t)
}

func sqlNaikkanVersi(t string) string {
	return fmt.Sprintf(`UPDATE %s SET SESSION_VERSION = SESSION_VERSION + 1, TGL_UPDATE = SYSDATE
	  WHERE LOGIN_ID = :1`, t)
}

func sqlInfoUnit(unit, divisi string) string {
	return fmt.Sprintf(`SELECT d.CODE, u.IS_ACTIVE FROM %s u JOIN %s d ON d.DIVISION_ID = u.DIVISION_ID
	  WHERE u.CODE = :1`, unit, divisi)
}

func sqlInfoDivisi(divisi, org string) string {
	return fmt.Sprintf(`SELECT o.CODE, d.IS_ACTIVE FROM %s d JOIN %s o ON o.ORGANIZATION_ID = d.ORGANIZATION_ID
	  WHERE d.CODE = :1`, divisi, org)
}

func sqlInfoOrganisasi(org string) string {
	return fmt.Sprintf(`SELECT IS_ACTIVE FROM %s WHERE CODE = :1`, org)
}

func sqlWorkbasketAktif(wb string) string {
	return fmt.Sprintf(`SELECT IS_ACTIVE FROM %s WHERE WORKBASKET_ID = :1`, wb)
}

func sqlSisipAkun(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (LOGIN_ID, NAME, PASSWORD_HASH, ORGANIZATION_CODE, DIVISION_CODE,
	    UNIT_CODE, IS_ACTIVE, MUST_CHANGE_PASSWORD) VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`, t)
}

func sqlSisipWorkbasket(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (LOGIN_ID, WORKBASKET_ID) VALUES (:1, :2)`, t)
}

func (g *GudangOracle) nama(logis string) (string, error) { return g.db.Qualify(logis) }

func (g *GudangOracle) baris(ctx context.Context, q string, args []any, tujuan ...any) error {
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	return g.db.QueryRowContext(ctx, q, args...).Scan(tujuan...)
}

// AmbilAkun membaca satu akun.
func (g *GudangOracle) AmbilAkun(ctx context.Context, id string) (Akun, error) {
	t, err := g.nama(tabelLogin)
	if err != nil {
		return Akun{}, err
	}
	var a Akun
	var org, div, unit sql.NullString
	var aktif, wajib string
	var kunci int
	err = g.baris(ctx, sqlAmbilAkun(t), []any{id}, &a.ID, &a.Nama, &a.HashSandi, &org, &div, &unit,
		&aktif, &wajib, &a.VersiSesi, &kunci)
	if errors.Is(err, sql.ErrNoRows) {
		return Akun{}, ErrAkunTidakAda
	}
	if err != nil {
		return Akun{}, fmt.Errorf("login: membaca akun: %w", err)
	}
	a.Organisasi, a.Divisi, a.Unit = org.String, div.String, unit.String
	a.Aktif, a.WajibGantiSandi, a.Terkunci = aktif == benderaYa, wajib == benderaYa, kunci == 1
	return a, nil
}

// Workbasket membaca WORKBASKET_ID aktif akun itu.
func (g *GudangOracle) Workbasket(ctx context.Context, id string) ([]string, error) {
	l, err := g.nama(tabelLoginWB)
	if err != nil {
		return nil, err
	}
	w, err := g.nama(tabelWB)
	if err != nil {
		return nil, err
	}
	q := sqlWorkbasket(l, w)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, id, masterAktif)
	if err != nil {
		return nil, fmt.Errorf("login: membaca workbasket: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := []string{}
	for rows.Next() {
		var wb string
		if err := rows.Scan(&wb); err != nil {
			return nil, err
		}
		out = append(out, wb)
	}
	return out, rows.Err()
}

func (g *GudangOracle) ubah(ctx context.Context, q string, args ...any) (int64, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	hasil, err := g.db.ExecContext(ctx, q, args...)
	if err != nil {
		return 0, fmt.Errorf("login: menulis M_LOGIN_GO: %w", err)
	}
	return hasil.RowsAffected()
}

// CatatGagal menaikkan hitungan gagal dan mengunci sesudah BatasGagal.
func (g *GudangOracle) CatatGagal(ctx context.Context, id string) error {
	t, err := g.nama(tabelLogin)
	if err != nil {
		return err
	}
	_, err = g.ubah(ctx, sqlCatatGagal(t), BatasGagal, int(LamaKunci.Minutes()), id)
	return err
}

// CatatBerhasil menolkan hitungan gagal dan mengisi LAST_LOGIN.
func (g *GudangOracle) CatatBerhasil(ctx context.Context, id string) error {
	t, err := g.nama(tabelLogin)
	if err != nil {
		return err
	}
	_, err = g.ubah(ctx, sqlCatatBerhasil(t), id)
	return err
}

// GantiSandi menulis hash baru dan menaikkan versi sesi.
func (g *GudangOracle) GantiSandi(ctx context.Context, id, hash string, versiLama int64) error {
	t, err := g.nama(tabelLogin)
	if err != nil {
		return err
	}
	n, err := g.ubah(ctx, sqlGantiSandi(t), hash, benderaTidak, id, versiLama)
	if err != nil {
		return err
	}
	if n != 1 {
		return ErrSesiTidakSah
	}
	return nil
}

// NaikkanVersi mencabut seluruh sesi akun itu.
func (g *GudangOracle) NaikkanVersi(ctx context.Context, id string) error {
	t, err := g.nama(tabelLogin)
	if err != nil {
		return err
	}
	_, err = g.ubah(ctx, sqlNaikkanVersi(t), id)
	return err
}

func (g *GudangOracle) info(ctx context.Context, q string, code string, tujuan ...any) error {
	err := g.baris(ctx, q, []any{code}, tujuan...)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrMasterTidakAda
	}
	return err
}

// InfoUnit membaca divisi pemilik unit.
func (g *GudangOracle) InfoUnit(ctx context.Context, code string) (string, bool, error) {
	u, err := g.nama(tabelUnit)
	if err != nil {
		return "", false, err
	}
	d, err := g.nama(tabelDivisi)
	if err != nil {
		return "", false, err
	}
	var div string
	var aktif int
	err = g.info(ctx, sqlInfoUnit(u, d), code, &div, &aktif)
	return div, aktif == masterAktif, err
}

// InfoDivisi membaca organisasi pemilik divisi.
func (g *GudangOracle) InfoDivisi(ctx context.Context, code string) (string, bool, error) {
	d, err := g.nama(tabelDivisi)
	if err != nil {
		return "", false, err
	}
	o, err := g.nama(tabelOrganis)
	if err != nil {
		return "", false, err
	}
	var org string
	var aktif int
	err = g.info(ctx, sqlInfoDivisi(d, o), code, &org, &aktif)
	return org, aktif == masterAktif, err
}

// InfoOrganisasi membaca status organisasi.
func (g *GudangOracle) InfoOrganisasi(ctx context.Context, code string) (bool, error) {
	o, err := g.nama(tabelOrganis)
	if err != nil {
		return false, err
	}
	var aktif int
	err = g.info(ctx, sqlInfoOrganisasi(o), code, &aktif)
	return aktif == masterAktif, err
}

// WorkbasketAktif membaca status workbasket.
func (g *GudangOracle) WorkbasketAktif(ctx context.Context, id string) (bool, error) {
	w, err := g.nama(tabelWB)
	if err != nil {
		return false, err
	}
	var aktif int
	err = g.info(ctx, sqlWorkbasketAktif(w), id, &aktif)
	return aktif == masterAktif, err
}

// BuatAkun menulis akun dan workbasket-nya dalam satu transaksi.
func (g *GudangOracle) BuatAkun(ctx context.Context, a AkunBaru, hash string) (err error) {
	t, err := g.nama(tabelLogin)
	if err != nil {
		return err
	}
	l, err := g.nama(tabelLoginWB)
	if err != nil {
		return err
	}
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	langkah := []struct {
		q    string
		args []any
	}{{sqlSisipAkun(t), []any{a.ID, a.Nama, hash, db.KosongJadiNil(a.Organisasi), db.KosongJadiNil(a.Divisi),
		db.KosongJadiNil(a.Unit), benderaYa, benderaYa}}}
	for _, w := range a.Workbasket {
		langkah = append(langkah, struct {
			q    string
			args []any
		}{sqlSisipWorkbasket(l), []any{a.ID, w}})
	}
	for _, s := range langkah {
		if err = db.PeriksaSQL(s.q); err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, s.q, s.args...); err != nil {
			if strings.Contains(err.Error(), "ORA-00001") {
				return fmt.Errorf("%w: %v", ErrAkunSudahAda, err)
			}
			return fmt.Errorf("login: membuat akun: %w", err)
		}
	}
	return tx.Commit()
}
