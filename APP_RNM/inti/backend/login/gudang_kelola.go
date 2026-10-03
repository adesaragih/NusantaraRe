package login

// Gudang Oracle - perubahan milik Kelola User (keputusan work owner 01-10-2026).
//
// ⛔ Setiap tabel lewat `Qualify` (ADR-U-0033), nilai lewat bind, nol COMMIT
// di teks SQL (ADR-U-0029): transaksi dibuka dan ditutup di sini.
//
// ⛔ ADMIN TERAKHIR: UbahAkun, SetelAktif, dan HapusAkun lebih dulu MENGUNCI
// baris setiap admin aktif (`SELECT ... FOR UPDATE`), lalu mengubah, lalu
// menghitung admin aktif yang tersisa - nol = ROLLBACK, `ErrAdminTerakhir`.
// Kunci itu membuat dua admin yang saling mencabut berjalan BERGILIR: yang
// kedua menunggu yang pertama selesai dan menghitung hasilnya.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/menu"
)

var _ GudangKelola = (*GudangOracle)(nil)

func sqlDaftarAkun(t string, satu bool) string {
	q := fmt.Sprintf(`SELECT LOGIN_ID, CONTACT_ID, NAME, ORGANIZATION_CODE, DIVISION_CODE, UNIT_CODE, IS_ACTIVE,
	        MUST_CHANGE_PASSWORD, CASE WHEN LOCKED_UNTIL > SYSDATE THEN 1 ELSE 0 END, LAST_LOGIN,
	        EMAIL, PHONE_NUMBER, EMPLOYEE_ID, JOB_POSITION
	   FROM %s`, t)
	if satu {
		return q + ` WHERE LOGIN_ID = :1`
	}
	return q + ` ORDER BY LOGIN_ID`
}

// sqlKunciAdmin dan sqlHitungAdmin: :1 = bendera aktif, :2 = KODE Kelola User.
// Urutan kunci TETAP (LOGIN_ID): dua transaksi admin serentak menunggu
// bergiliran, bukan saling mengunci (ORA-00060).
func sqlKunciAdmin(login, mn string) string {
	return fmt.Sprintf(`SELECT l.LOGIN_ID FROM %s l
	  WHERE l.IS_ACTIVE = :1 AND EXISTS (SELECT 1 FROM %s m WHERE m.LOGIN_ID = l.LOGIN_ID AND m.MENU_KODE = :2)
	  ORDER BY l.LOGIN_ID
	  FOR UPDATE`, login, mn)
}

func sqlHitungAdmin(login, mn string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s l
	  WHERE l.IS_ACTIVE = :1 AND EXISTS (SELECT 1 FROM %s m WHERE m.LOGIN_ID = l.LOGIN_ID AND m.MENU_KODE = :2)`, login, mn)
}

func sqlUbahProfil(t string) string {
	return fmt.Sprintf(`UPDATE %s SET NAME = :1, ORGANIZATION_CODE = :2, DIVISION_CODE = :3, UNIT_CODE = :4,
	    EMAIL = :5, PHONE_NUMBER = :6, EMPLOYEE_ID = :7, JOB_POSITION = :8, TGL_UPDATE = SYSDATE
	  WHERE LOGIN_ID = :9`, t)
}

// sqlSetelAktif - :1 = bendera, :2 = LOGIN_ID. Versi sesi selalu naik.
func sqlSetelAktif(t string) string {
	return fmt.Sprintf(`UPDATE %s SET IS_ACTIVE = :1, SESSION_VERSION = SESSION_VERSION + 1, TGL_UPDATE = SYSDATE
	  WHERE LOGIN_ID = :2`, t)
}

func sqlBukaKunci(t string) string {
	return fmt.Sprintf(`UPDATE %s SET FAILED_COUNT = 0, LOCKED_UNTIL = NULL, TGL_UPDATE = SYSDATE
	  WHERE LOGIN_ID = :1`, t)
}

// sqlAturSandi - password baru: :1 hash, :2 wajib ganti, :3 LOGIN_ID; sesi
// dicabut dan kunci dibuka. sqlSetelWajibGanti - centangnya saja.
func sqlAturSandi(t string) string {
	return fmt.Sprintf(`UPDATE %s SET PASSWORD_HASH = :1, MUST_CHANGE_PASSWORD = :2,
	    SESSION_VERSION = SESSION_VERSION + 1, FAILED_COUNT = 0, LOCKED_UNTIL = NULL, TGL_UPDATE = SYSDATE
	  WHERE LOGIN_ID = :3`, t)
}

func sqlSetelWajibGanti(t string) string {
	return fmt.Sprintf(`UPDATE %s SET MUST_CHANGE_PASSWORD = :1, TGL_UPDATE = SYSDATE WHERE LOGIN_ID = :2`, t)
}

func sqlHapusMilik(t string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE LOGIN_ID = :1`, t)
}

func sqlMasterOrganisasi(o string) string {
	return fmt.Sprintf(`SELECT CODE, NAME, NULL FROM %s WHERE IS_ACTIVE = :1 ORDER BY NAME, CODE`, o)
}

func sqlMasterDivisi(d, o string) string {
	return fmt.Sprintf(`SELECT d.CODE, d.NAME, o.CODE FROM %s d JOIN %s o ON o.ORGANIZATION_ID = d.ORGANIZATION_ID
	  WHERE d.IS_ACTIVE = :1 ORDER BY d.SORT_ORDER, d.NAME, d.CODE`, d, o)
}

func sqlMasterUnit(u, d string) string {
	return fmt.Sprintf(`SELECT u.CODE, u.NAME, d.CODE FROM %s u JOIN %s d ON d.DIVISION_ID = u.DIVISION_ID
	  WHERE u.IS_ACTIVE = :1 ORDER BY u.SORT_ORDER, u.NAME, u.CODE`, u, d)
}

func sqlMasterWorkbasket(w string) string {
	return fmt.Sprintf(`SELECT WORKBASKET_ID, NAME, NULL FROM %s WHERE IS_ACTIVE = :1 ORDER BY WORKBASKET_ID`, w)
}

// pemindai - Scan milik *sql.Row dan *sql.Rows.
type pemindai interface{ Scan(...any) error }

func pindaiRingkas(p pemindai) (RingkasAkun, error) {
	var r RingkasAkun
	var idKontak, org, div, unit, nama, email, telepon, nik, jabatan sql.NullString
	var aktif, wajib string
	var kunci int
	var terakhir sql.NullTime
	if err := p.Scan(&r.AkunID, &idKontak, &nama, &org, &div, &unit, &aktif, &wajib, &kunci, &terakhir,
		&email, &telepon, &nik, &jabatan); err != nil {
		return RingkasAkun{}, err
	}
	r.IDKontak, r.Nama, r.Organisasi, r.Divisi, r.Unit = idKontak.String, nama.String, org.String, div.String, unit.String
	r.Kontak = Kontak{Email: email.String, Telepon: telepon.String, NIK: nik.String, Jabatan: jabatan.String}
	r.Aktif, r.WajibGantiSandi, r.Terkunci = aktif == benderaYa, wajib == benderaYa, kunci == 1
	if terakhir.Valid {
		r.LoginTerakhir = terakhir.Time.Format("2006-01-02 15:04")
	}
	return r, nil
}

// DaftarAkun membaca seluruh akun, urut LOGIN_ID.
func (g *GudangOracle) DaftarAkun(ctx context.Context) ([]RingkasAkun, error) {
	t, err := g.nama(tabelLogin)
	if err != nil {
		return nil, err
	}
	q := sqlDaftarAkun(t, false)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("login: membaca daftar akun: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := []RingkasAkun{}
	for rows.Next() {
		r, err := pindaiRingkas(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// RingkasAkun membaca satu baris daftar.
func (g *GudangOracle) RingkasAkun(ctx context.Context, id string) (RingkasAkun, error) {
	t, err := g.nama(tabelLogin)
	if err != nil {
		return RingkasAkun{}, err
	}
	q := sqlDaftarAkun(t, true)
	if err := db.PeriksaSQL(q); err != nil {
		return RingkasAkun{}, err
	}
	r, err := pindaiRingkas(g.db.QueryRowContext(ctx, q, id))
	if errors.Is(err, sql.ErrNoRows) {
		return RingkasAkun{}, ErrAkunTidakAda
	}
	if err != nil {
		return RingkasAkun{}, fmt.Errorf("login: membaca akun: %w", err)
	}
	return r, nil
}

// tabelKelola - nama lengkap ketiga tabel login.
type tabelKelola struct{ login, wb, menu string }

func (g *GudangOracle) tabelKelola() (tabelKelola, error) {
	var t tabelKelola
	var err error
	if t.login, err = g.nama(tabelLogin); err != nil {
		return t, err
	}
	if t.wb, err = g.nama(tabelLoginWB); err != nil {
		return t, err
	}
	t.menu, err = g.nama(tabelLoginMn)
	return t, err
}

// dalamTransaksiAdmin mengunci admin aktif, menjalankan `ubah`, lalu menolak
// hasil yang menyisakan nol admin aktif. Nol COMMIT bila apa pun gagal.
func (g *GudangOracle) dalamTransaksiAdmin(ctx context.Context, t tabelKelola, ubah func(*db.Tx) error) (err error) {
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	kunci := sqlKunciAdmin(t.login, t.menu)
	if err = db.PeriksaSQL(kunci); err != nil {
		return err
	}
	rows, err := tx.QueryContext(ctx, kunci, benderaYa, menu.KodeKelolaUser)
	if err != nil {
		return fmt.Errorf("login: mengunci admin: %w", err)
	}
	for rows.Next() {
		var terkunci string
		if err = rows.Scan(&terkunci); err != nil {
			_ = rows.Close()
			return fmt.Errorf("login: mengunci admin: %w", err)
		}
	}
	if err = rows.Err(); err != nil {
		_ = rows.Close()
		return fmt.Errorf("login: mengunci admin: %w", err)
	}
	_ = rows.Close()
	if err = ubah(tx); err != nil {
		return err
	}
	hitung := sqlHitungAdmin(t.login, t.menu)
	if err = db.PeriksaSQL(hitung); err != nil {
		return err
	}
	var n int
	if err = tx.QueryRowContext(ctx, hitung, benderaYa, menu.KodeKelolaUser).Scan(&n); err != nil {
		return fmt.Errorf("login: menghitung admin: %w", err)
	}
	if n == 0 {
		err = ErrAdminTerakhir
		return err
	}
	return tx.Commit()
}

// tulisTx menjalankan satu pernyataan tulis; `harusSatu` = nol baris berarti
// akunnya tidak ada.
func tulisTx(ctx context.Context, tx *db.Tx, q string, harusSatu bool, args ...any) error {
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		if ganda := galatGanda(err); ganda != nil {
			return ganda
		}
		return fmt.Errorf("login: menulis akun: %w", err)
	}
	if harusSatu {
		if n, err := hasil.RowsAffected(); err != nil {
			return err
		} else if n == 0 {
			return ErrAkunTidakAda
		}
	}
	return nil
}

// UbahAkun menulis profil dan mengganti seluruh workbasket dan menunya.
func (g *GudangOracle) UbahAkun(ctx context.Context, id string, a IsianAkun) error {
	t, err := g.tabelKelola()
	if err != nil {
		return err
	}
	return g.dalamTransaksiAdmin(ctx, t, func(tx *db.Tx) error {
		args := append([]any{a.Nama, db.KosongJadiNil(a.Organisasi), db.KosongJadiNil(a.Divisi), db.KosongJadiNil(a.Unit)},
			nilaiKontak(a.Kontak)...)
		if err := tulisTx(ctx, tx, sqlUbahProfil(t.login), true, append(args, id)...); err != nil {
			return err
		}
		if err := tulisTx(ctx, tx, sqlHapusMilik(t.wb), false, id); err != nil {
			return err
		}
		for _, w := range a.Workbasket {
			if err := tulisTx(ctx, tx, sqlSisipWorkbasket(t.wb), false, id, w); err != nil {
				return err
			}
		}
		if err := tulisTx(ctx, tx, sqlHapusMilik(t.menu), false, id); err != nil {
			return err
		}
		for _, m := range a.Menu {
			if err := tulisTx(ctx, tx, sqlSisipMenu(t.menu), false, id, m); err != nil {
				return err
			}
		}
		return nil
	})
}

// SetelAktif mengubah IS_ACTIVE dan mencabut seluruh sesi akun itu.
func (g *GudangOracle) SetelAktif(ctx context.Context, id string, aktif bool) error {
	t, err := g.tabelKelola()
	if err != nil {
		return err
	}
	bendera := benderaTidak
	if aktif {
		bendera = benderaYa
	}
	return g.dalamTransaksiAdmin(ctx, t, func(tx *db.Tx) error {
		return tulisTx(ctx, tx, sqlSetelAktif(t.login), true, bendera, id)
	})
}

// BukaKunci menolkan hitungan gagal dan mencabut kunci.
func (g *GudangOracle) BukaKunci(ctx context.Context, id string) error {
	t, err := g.nama(tabelLogin)
	if err != nil {
		return err
	}
	n, err := g.ubah(ctx, sqlBukaKunci(t), id)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrAkunTidakAda
	}
	return nil
}

// AturSandi menulis password baru (bila `hash` terisi) dan centang wajib ganti.
func (g *GudangOracle) AturSandi(ctx context.Context, id, hash string, wajibGanti bool) error {
	t, err := g.nama(tabelLogin)
	if err != nil {
		return err
	}
	wajib := benderaTidak
	if wajibGanti {
		wajib = benderaYa
	}
	var n int64
	if hash != "" {
		n, err = g.ubah(ctx, sqlAturSandi(t), hash, wajib, id)
	} else {
		n, err = g.ubah(ctx, sqlSetelWajibGanti(t), wajib, id)
	}
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrAkunTidakAda
	}
	return nil
}

// HapusAkun membuang menu, workbasket, lalu akunnya - PERMANEN, satu transaksi.
func (g *GudangOracle) HapusAkun(ctx context.Context, id string) error {
	t, err := g.tabelKelola()
	if err != nil {
		return err
	}
	return g.dalamTransaksiAdmin(ctx, t, func(tx *db.Tx) error {
		if err := tulisTx(ctx, tx, sqlHapusMilik(t.menu), false, id); err != nil {
			return err
		}
		if err := tulisTx(ctx, tx, sqlHapusMilik(t.wb), false, id); err != nil {
			return err
		}
		return tulisTx(ctx, tx, sqlHapusMilik(t.login), true, id)
	})
}

// opsiMaster membaca satu daftar master aktif.
func (g *GudangOracle) opsiMaster(ctx context.Context, apa, q string) ([]OpsiMaster, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, masterAktif)
	if err != nil {
		return nil, fmt.Errorf("login: membaca master %s: %w", apa, err)
	}
	defer func() { _ = rows.Close() }()
	out := []OpsiMaster{}
	for rows.Next() {
		var kode string
		var nama, induk sql.NullString
		if err := rows.Scan(&kode, &nama, &induk); err != nil {
			return nil, err
		}
		out = append(out, OpsiMaster{Kode: kode, Nama: nama.String, Induk: induk.String})
	}
	return out, rows.Err()
}

// Master membaca organisasi, divisi, unit, dan workbasket yang AKTIF.
func (g *GudangOracle) Master(ctx context.Context) (PilihanKelola, error) {
	var p PilihanKelola
	nama := map[string]string{}
	for _, n := range []string{tabelOrganis, tabelDivisi, tabelUnit, tabelWB} {
		q, err := g.nama(n)
		if err != nil {
			return p, err
		}
		nama[n] = q
	}
	var err error
	if p.Organisasi, err = g.opsiMaster(ctx, "organisasi", sqlMasterOrganisasi(nama[tabelOrganis])); err != nil {
		return p, err
	}
	if p.Divisi, err = g.opsiMaster(ctx, "divisi", sqlMasterDivisi(nama[tabelDivisi], nama[tabelOrganis])); err != nil {
		return p, err
	}
	if p.Unit, err = g.opsiMaster(ctx, "unit", sqlMasterUnit(nama[tabelUnit], nama[tabelDivisi])); err != nil {
		return p, err
	}
	p.Workbasket, err = g.opsiMaster(ctx, "workbasket", sqlMasterWorkbasket(nama[tabelWB]))
	return p, err
}
