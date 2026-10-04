// Package repository adalah SATU-SATUNYA lapisan modul Company Detail yang berbicara ke Oracle.
//
// Tabel warisan `CLIENT`, `CLIENT_PICLIST`, dan `CLIENT_ADDRESS` ditulis langsung (kolom tambahannya dari migrasi
// 805-807); `M_ENUMERASI` dan `NATION` (migrasi 800-804) dibaca saja. Nomor ORG organisasi baru dari sequence
// `SEQ_CLIENT_ORG` (migrasi 810 memulainya sesudah nomor ORG tertinggi CLIENT dan M_CLIENT; perintah work owner
// 04-10-2026: "buat seq aja, start-nya dari id max+1").
//
// ⛔ Prosedur `RDBINSERTCLIENT` TIDAK dipanggil. Setiap tabel lewat `Qualify` (ADR-U-0033), nilai lewat bind, nol
// COMMIT di SQL (ADR-U-0029): transaksinya milik services.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/companydetail/backend/models"
)

// Objek Oracle modul ini.
const (
	TabelClient = "CLIENT"
	TabelPIC    = "CLIENT_PICLIST"
	TabelAlamat = "CLIENT_ADDRESS"
	TabelEnum   = "M_ENUMERASI"
	TabelNegara = "NATION"
	// TabelDokumen - dokumen organisasi Pega (Copy Old dan alat pindah).
	TabelDokumen = "M_CLIENT"
	// SequenceOrg - nomor ORG organisasi baru (migrasi 810).
	SequenceOrg = "SEQ_CLIENT_ORG"
)

// DaftarTabelDitulis dan DaftarTabelDibacaSaja - penjaga modul.
var (
	DaftarTabelDitulis    = []string{TabelClient, TabelPIC, TabelAlamat}
	DaftarTabelDibacaSaja = []string{TabelEnum, TabelNegara, TabelDokumen}
)

// BatasCariInduk - pilihan Parent organization paling banyak sekali cari.
const BatasCariInduk = 20

// polaID - `ID` organisasi `ASM-SFAGIS-WORK-ORG ORG-n` di CLIENT dan M_CLIENT (Copy Old dan alat pindah). Pola yang
// sama dipakai migrasi 810 untuk nilai awal SEQ_CLIENT_ORG.
const polaID = `^ASM-SFAGIS-WORK-ORG ORG-[0-9]+$`

var (
	// ErrTidakAda - baris yang dicari tidak ada.
	ErrTidakAda = errors.New("repository: row not found")
	// ErrBelumDimigrasi - kolom, tabel, atau sequence migrasi 800-810 belum ada.
	ErrBelumDimigrasi = errors.New("repository: the Company Detail tables are not migrated yet - run -migrate (migrations 800-810)")
)

// Gudang membaca dan menulis Oracle.
type Gudang struct{ db *db.DB }

// Baru menyusun Gudang.
func Baru(d *db.DB) *Gudang { return &Gudang{db: d} }

// penjalan - *db.DB dan *db.Tx.
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

func (g *Gudang) nama(objek string) (string, error) { return g.db.Qualify(objek) }

// belumDimigrasi - ORA-00904 (kolom belum ada), ORA-00942 (tabel belum ada), atau ORA-02289 (sequence belum ada).
func belumDimigrasi(err error) bool {
	if err == nil {
		return false
	}
	p := err.Error()
	return strings.Contains(p, "ORA-00904") || strings.Contains(p, "ORA-00942") || strings.Contains(p, "ORA-02289")
}

func bungkus(err error, apa string) error {
	if belumDimigrasi(err) {
		return ErrBelumDimigrasi
	}
	return fmt.Errorf("repository: %s: %w", apa, err)
}

// PolaCari membentuk pola LIKE huruf besar yang aman: `\`, `%`, dan `_` ketikan pemakai bukan wildcard.
func PolaCari(kueri string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return "%" + r.Replace(strings.ToUpper(strings.TrimSpace(kueri))) + "%"
}

// saringCari - :2 kosong (NULL) = tanpa saringan; :3-:5 pola yang sama untuk NAME, IDVIEW, NPWP.
const saringCari = `FLAG = :1 AND (:2 IS NULL OR UPPER(NAME) LIKE :3 ESCAPE '\' OR UPPER(IDVIEW) LIKE :4 ESCAPE '\'
	    OR UPPER(NPWP) LIKE :5 ESCAPE '\')`

func sqlCariOrg(t string) string {
	return fmt.Sprintf(`SELECT ID, IDVIEW, NAME, TITLE, NPWP, COUNTRYNAME, BU_ID, GROUPNAME FROM %s
	  WHERE %s
	  ORDER BY UPPER(NAME), ID OFFSET :6 ROWS FETCH NEXT :7 ROWS ONLY`, t, saringCari)
}

func sqlHitungOrg(t string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, t, saringCari)
}

// kolomOrg - urutan kolom SELECT satu organisasi (`pindaiOrg`).
const kolomOrg = `ID, IDVIEW, NAME, TITLE, NPWP, COUNTRY, COUNTRYNAME, BU_ID, PARENT_ID, GROUPNAME, NOTE,
	CREATED_BY, TO_CHAR(CREATED_AT, 'YYYY-MM-DD HH24:MI'), UPDATED_BY, TO_CHAR(UPDATED_AT, 'YYYY-MM-DD HH24:MI')`

func sqlAmbilOrg(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1 AND FLAG = :2`, kolomOrg, t)
}

// sqlNomorOrgBerikut - satu nomor ORG dari sequence; dua Create bersamaan tidak pernah mendapat nomor yang sama,
// jadi CLIENT tidak perlu dikunci.
func sqlNomorOrgBerikut(seq string) string { return fmt.Sprintf(`SELECT %s.NEXTVAL FROM DUAL`, seq) }

// sqlSisipOrg - kolom baris organisasi: yang ditulis `RDBINSERTCLIENT` (tanpa IDNUMBER, BU_NOTE) + kolom 805.
func sqlSisipOrg(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, NAME, FLAG, BU_ID, IDVIEW, NPWP, GROUPNAME, TITLE, COUNTRY, COUNTRYNAME,
	    PARENT_ID, NOTE, CREATED_BY, CREATED_AT, UPDATED_BY, UPDATED_AT)
	  VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, SYSTIMESTAMP, :14, SYSTIMESTAMP)`, t)
}

// sqlPerbaruiOrg - ID, IDVIEW, FLAG, IDNUMBER, BU_NOTE, dan jejak pembuatan TIDAK pernah diubah.
func sqlPerbaruiOrg(t string) string {
	return fmt.Sprintf(`UPDATE %s SET NAME = :1, TITLE = :2, NPWP = :3, COUNTRY = :4, COUNTRYNAME = :5, BU_ID = :6,
	    PARENT_ID = :7, GROUPNAME = :8, NOTE = :9, UPDATED_BY = :10, UPDATED_AT = SYSTIMESTAMP
	  WHERE ID = :11 AND FLAG = :12`, t)
}

func sqlDaftarPIC(t string) string {
	return fmt.Sprintf(`SELECT USERIDENTIFIER, NICKNAME, POSITION, GENDER, EMAIL, DATEOFBIRTH, PHONENUMBER FROM %s
	  WHERE CLIENTID = :1 ORDER BY ROWID`, t)
}

func sqlHapusPIC(t string) string { return fmt.Sprintf(`DELETE FROM %s WHERE CLIENTID = :1`, t) }

func sqlSisipPIC(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (CLIENTID, USERIDENTIFIER, NICKNAME, EMAIL, POSITION, DATEOFBIRTH, PHONENUMBER,
	    GENDER)
	  VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`, t)
}

// kolomAlamat - urutan kolom baris CLIENT_ADDRESS (`pindaiAlamat`, `NilaiSisipAlamat`).
const kolomAlamat = `CLIENTID, ASMADDRESSTYPE, ASMADDRESS, ASMCITY, CITYNAME, ASMZIPCODE, DISTRICTNAME, PROVINCENAME,
	RWNAME, PXCREATEOPERATOR, PXCREATEDATETIME, TELFAX_TYPE, TELFAX_CODE, TELFAX_NO`

func sqlDaftarAlamat(t string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE CLIENTID = :1 ORDER BY ROWID`, kolomAlamat, t)
}

func sqlHapusAlamat(t string) string { return fmt.Sprintf(`DELETE FROM %s WHERE CLIENTID = :1`, t) }

func sqlSisipAlamat(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (%s)
	  VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14)`, t, kolomAlamat)
}

func sqlDaftarPilihan(t string) string {
	return fmt.Sprintf(`SELECT JENIS, KODE, LABEL, AKTIF FROM %s ORDER BY JENIS, URUTAN, KODE`, t)
}

func sqlDaftarNegara(t string) string {
	return fmt.Sprintf(`SELECT ID, OLDID, NOTE, NATIONINITIAL FROM %s ORDER BY UPPER(NOTE), ID`, t)
}

// sqlNamaOrg - nama seluruh organisasi, untuk pemeriksaan nama sama/mirip sebelum Create (18 ribuan baris).
func sqlNamaOrg(t string) string {
	return fmt.Sprintf(`SELECT ID, IDVIEW, NAME, TITLE, COUNTRYNAME FROM %s WHERE FLAG = :1 AND NAME IS NOT NULL`, t)
}

// sqlCariInduk - :3 = ID organisasi yang sedang diubah (dikecualikan; NULL = tidak ada).
func sqlCariInduk(t string) string {
	return fmt.Sprintf(`SELECT ID, IDVIEW, NAME FROM %s
	  WHERE FLAG = :1 AND UPPER(NAME) LIKE :2 ESCAPE '\' AND ID <> NVL(:3, '-')
	  ORDER BY UPPER(NAME), ID FETCH FIRST %d ROWS ONLY`, t, BatasCariInduk)
}

// pemindai - Scan milik *sql.Row dan *sql.Rows.
type pemindai interface{ Scan(...any) error }

func pindaiTeks(p pemindai, n int) ([]string, error) {
	v := make([]sql.NullString, n)
	tujuan := make([]any, n)
	for i := range v {
		tujuan[i] = &v[i]
	}
	if err := p.Scan(tujuan...); err != nil {
		return nil, err
	}
	out := make([]string, n)
	for i := range v {
		out[i] = v[i].String
	}
	return out, nil
}

func pindaiBaris(p pemindai) (models.BarisDaftar, error) {
	v, err := pindaiTeks(p, 8)
	if err != nil {
		return models.BarisDaftar{}, err
	}
	return models.BarisDaftar{ID: v[0], IDView: v[1], Nama: v[2], Title: v[3], NPWP: v[4], CountryName: v[5],
		BusinessField: v[6], ParentName: v[7]}, nil
}

func pindaiOrg(p pemindai) (models.Organisasi, error) {
	v, err := pindaiTeks(p, 15)
	if err != nil {
		return models.Organisasi{}, err
	}
	return models.Organisasi{ID: v[0], IDView: v[1], Nama: v[2], Title: v[3], NPWP: v[4], Country: v[5],
		CountryName: v[6], BusinessField: v[7], ParentID: v[8], ParentName: v[9], Note: v[10], CreatedBy: v[11],
		CreatedAt: v[12], UpdatedBy: v[13], UpdatedAt: v[14]}, nil
}

func pindaiPIC(p pemindai) (models.PIC, error) {
	v, err := pindaiTeks(p, 7)
	if err != nil {
		return models.PIC{}, err
	}
	return models.PIC{UserIdentifier: v[0], Nama: v[1], Position: v[2], Gender: v[3], Email: v[4], DateOfBirth: v[5],
		Phone: v[6]}, nil
}

func pindaiAlamat(p pemindai) (models.BarisAlamat, error) {
	v, err := pindaiTeks(p, 14)
	if err != nil {
		return models.BarisAlamat{}, err
	}
	return models.BarisAlamat{ClientID: v[0], Type: v[1], Address: v[2], City: v[3], CityName: v[4], ZipCode: v[5],
		DistrictName: v[6], ProvinceName: v[7], RWName: v[8], PxCreateOperator: v[9], PxCreateDateTime: v[10],
		TelfaxType: v[11], TelfaxCode: v[12], TelfaxNo: v[13]}, nil
}

func pindaiPilihan(p pemindai) (models.Pilihan, error) {
	v, err := pindaiTeks(p, 4)
	if err != nil {
		return models.Pilihan{}, err
	}
	return models.Pilihan{Jenis: v[0], Kode: v[1], Label: v[2], Aktif: v[3] == "1"}, nil
}

func pindaiNegara(p pemindai) (models.Negara, error) {
	v, err := pindaiTeks(p, 4)
	if err != nil {
		return models.Negara{}, err
	}
	return models.Negara{ID: v[0], OldID: v[1], Nama: v[2], NationInitial: v[3]}, nil
}

func pindaiInduk(p pemindai) (models.BarisDaftar, error) {
	v, err := pindaiTeks(p, 3)
	if err != nil {
		return models.BarisDaftar{}, err
	}
	return models.BarisDaftar{ID: v[0], IDView: v[1], Nama: v[2]}, nil
}

// daftar menjalankan satu SELECT banyak baris.
func daftar[T any](ctx context.Context, j penjalan, q string, pindai func(pemindai) (T, error), args ...any) ([]T, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := j.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, bungkus(err, "membaca")
	}
	defer func() { _ = rows.Close() }()
	out := []T{}
	for rows.Next() {
		v, err := pindai(rows)
		if err != nil {
			return nil, bungkus(err, "memindai")
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, bungkus(err, "membaca")
	}
	return out, nil
}

// satu menjalankan satu SELECT satu baris; nol baris = ErrTidakAda.
func satu[T any](ctx context.Context, j penjalan, q string, pindai func(pemindai) (T, error), args ...any) (T, error) {
	var nol T
	if err := db.PeriksaSQL(q); err != nil {
		return nol, err
	}
	v, err := pindai(j.QueryRowContext(ctx, q, args...))
	if errors.Is(err, sql.ErrNoRows) {
		return nol, ErrTidakAda
	}
	if err != nil {
		return nol, bungkus(err, "membaca")
	}
	return v, nil
}

// jalankan menjalankan satu DML/DDL.
func jalankan(ctx context.Context, j penjalan, q, apa string, args ...any) (sql.Result, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	hasil, err := j.ExecContext(ctx, q, args...)
	if err != nil {
		return nil, bungkus(err, apa)
	}
	return hasil, nil
}

// NilaiCari - lima nilai bind `saringCari`; kueri kosong = tanpa saringan.
func NilaiCari(kueri string) []any {
	if strings.TrimSpace(kueri) == "" {
		return []any{models.FlagOrg, nil, nil, nil, nil}
	}
	p := PolaCari(kueri)
	return []any{models.FlagOrg, p, p, p, p}
}

// CariOrg membaca satu halaman organisasi dan jumlah seluruhnya.
func (g *Gudang) CariOrg(ctx context.Context, kueri string, offset, ukuran int) ([]models.BarisDaftar, int, error) {
	t, err := g.nama(TabelClient)
	if err != nil {
		return nil, 0, err
	}
	nilai := NilaiCari(kueri)
	jumlah, err := satu(ctx, g.db, sqlHitungOrg(t), func(p pemindai) (int, error) {
		var n int
		err := p.Scan(&n)
		return n, err
	}, nilai...)
	if err != nil {
		return nil, 0, err
	}
	baris, err := daftar(ctx, g.db, sqlCariOrg(t), pindaiBaris, append(nilai, offset, ukuran)...)
	return baris, jumlah, err
}

// AmbilOrg membaca satu organisasi.
func (g *Gudang) AmbilOrg(ctx context.Context, tx *db.Tx, id string) (models.Organisasi, error) {
	t, err := g.nama(TabelClient)
	if err != nil {
		return models.Organisasi{}, err
	}
	return satu(ctx, g.dari(tx), sqlAmbilOrg(t), pindaiOrg, id, models.FlagOrg)
}

// NomorOrgBerikut - SEQ_CLIENT_ORG.NEXTVAL, nomor ORG organisasi baru.
func (g *Gudang) NomorOrgBerikut(ctx context.Context, tx *db.Tx) (int64, error) {
	s, err := g.nama(SequenceOrg)
	if err != nil {
		return 0, err
	}
	return satu(ctx, g.dari(tx), sqlNomorOrgBerikut(s), func(p pemindai) (int64, error) {
		var n int64
		err := p.Scan(&n)
		return n, err
	})
}

// NilaiSisipOrg - empat belas nilai bind `sqlSisipOrg`, berurutan; kosong = NULL.
func NilaiSisipOrg(o models.Organisasi) []any {
	n := db.KosongJadiNil
	return []any{o.ID, n(o.Nama), models.FlagOrg, n(o.BusinessField), n(o.IDView), n(o.NPWP), n(o.ParentName),
		n(o.Title), n(o.Country), n(o.CountryName), n(o.ParentID), n(o.Note), n(o.CreatedBy), n(o.UpdatedBy)}
}

// NilaiPerbaruiOrg - dua belas nilai bind `sqlPerbaruiOrg`, berurutan; kosong = NULL.
func NilaiPerbaruiOrg(o models.Organisasi) []any {
	n := db.KosongJadiNil
	return []any{n(o.Nama), n(o.Title), n(o.NPWP), n(o.Country), n(o.CountryName), n(o.BusinessField), n(o.ParentID),
		n(o.ParentName), n(o.Note), n(o.UpdatedBy), o.ID, models.FlagOrg}
}

// SisipOrg menyisipkan baris organisasi baru (ID sudah dibentuk services).
func (g *Gudang) SisipOrg(ctx context.Context, tx *db.Tx, o models.Organisasi) error {
	t, err := g.nama(TabelClient)
	if err != nil {
		return err
	}
	_, err = jalankan(ctx, g.dari(tx), sqlSisipOrg(t), "menyisipkan organisasi", NilaiSisipOrg(o)...)
	return err
}

// PerbaruiOrg menulis ulang kolom layar; nol baris = ErrTidakAda.
func (g *Gudang) PerbaruiOrg(ctx context.Context, tx *db.Tx, o models.Organisasi) error {
	t, err := g.nama(TabelClient)
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, g.dari(tx), sqlPerbaruiOrg(t), "memperbarui organisasi", NilaiPerbaruiOrg(o)...)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrTidakAda
	}
	return nil
}

// DaftarPIC membaca PIC satu organisasi apa adanya (DATEOFBIRTH bentuk simpan).
func (g *Gudang) DaftarPIC(ctx context.Context, tx *db.Tx, clientID string) ([]models.PIC, error) {
	t, err := g.nama(TabelPIC)
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.dari(tx), sqlDaftarPIC(t), pindaiPIC, clientID)
}

// NilaiSisipPIC - delapan nilai bind `sqlSisipPIC`, berurutan; kosong = NULL.
func NilaiSisipPIC(clientID string, p models.PIC) []any {
	n := db.KosongJadiNil
	return []any{clientID, n(p.UserIdentifier), n(p.Nama), n(p.Email), n(p.Position), n(p.DateOfBirth), n(p.Phone),
		n(p.Gender)}
}

// GantiPIC mengganti SELURUH PIC satu organisasi: hapus lalu sisip (transaksi yang sama).
func (g *Gudang) GantiPIC(ctx context.Context, tx *db.Tx, clientID string, pic []models.PIC) error {
	t, err := g.nama(TabelPIC)
	if err != nil {
		return err
	}
	j := g.dari(tx)
	if _, err := jalankan(ctx, j, sqlHapusPIC(t), "menghapus PIC", clientID); err != nil {
		return err
	}
	for _, p := range pic {
		if _, err := jalankan(ctx, j, sqlSisipPIC(t), "menyisipkan PIC", NilaiSisipPIC(clientID, p)...); err != nil {
			return err
		}
	}
	return nil
}

// DaftarBarisAlamat membaca baris CLIENT_ADDRESS satu organisasi apa adanya.
func (g *Gudang) DaftarBarisAlamat(ctx context.Context, tx *db.Tx, clientID string) ([]models.BarisAlamat, error) {
	t, err := g.nama(TabelAlamat)
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.dari(tx), sqlDaftarAlamat(t), pindaiAlamat, clientID)
}

// NilaiSisipAlamat - empat belas nilai bind `sqlSisipAlamat`, berurutan; kosong = NULL.
func NilaiSisipAlamat(b models.BarisAlamat) []any {
	n := db.KosongJadiNil
	return []any{b.ClientID, n(b.Type), n(b.Address), n(b.City), n(b.CityName), n(b.ZipCode), n(b.DistrictName),
		n(b.ProvinceName), n(b.RWName), n(b.PxCreateOperator), n(b.PxCreateDateTime), n(b.TelfaxType),
		n(b.TelfaxCode), n(b.TelfaxNo)}
}

// GantiAlamat mengganti SELURUH baris alamat satu organisasi: hapus lalu sisip (transaksi yang sama).
func (g *Gudang) GantiAlamat(ctx context.Context, tx *db.Tx, clientID string, baris []models.BarisAlamat) error {
	t, err := g.nama(TabelAlamat)
	if err != nil {
		return err
	}
	j := g.dari(tx)
	if _, err := jalankan(ctx, j, sqlHapusAlamat(t), "menghapus alamat", clientID); err != nil {
		return err
	}
	for _, b := range baris {
		b.ClientID = clientID
		if _, err := jalankan(ctx, j, sqlSisipAlamat(t), "menyisipkan alamat", NilaiSisipAlamat(b)...); err != nil {
			return err
		}
	}
	return nil
}

// DaftarPilihan membaca seluruh M_ENUMERASI.
func (g *Gudang) DaftarPilihan(ctx context.Context) ([]models.Pilihan, error) {
	t, err := g.nama(TabelEnum)
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.db, sqlDaftarPilihan(t), pindaiPilihan)
}

// DaftarNegara membaca seluruh NATION.
func (g *Gudang) DaftarNegara(ctx context.Context) ([]models.Negara, error) {
	t, err := g.nama(TabelNegara)
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.db, sqlDaftarNegara(t), pindaiNegara)
}

// DaftarNamaOrg membaca ID, ORG ID, nama, title, dan negara seluruh organisasi.
func (g *Gudang) DaftarNamaOrg(ctx context.Context) ([]models.BarisDaftar, error) {
	t, err := g.nama(TabelClient)
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.db, sqlNamaOrg(t), func(p pemindai) (models.BarisDaftar, error) {
		v, err := pindaiTeks(p, 5)
		if err != nil {
			return models.BarisDaftar{}, err
		}
		return models.BarisDaftar{ID: v[0], IDView: v[1], Nama: v[2], Title: v[3], CountryName: v[4]}, nil
	}, models.FlagOrg)
}

// CariInduk - organisasi bernama mirip `kueri` (pilihan Parent organization); `kecuali` kosong = tidak ada.
func (g *Gudang) CariInduk(ctx context.Context, kueri, kecuali string) ([]models.BarisDaftar, error) {
	t, err := g.nama(TabelClient)
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.db, sqlCariInduk(t), pindaiInduk, models.FlagOrg, PolaCari(kueri), db.KosongJadiNil(kecuali))
}
