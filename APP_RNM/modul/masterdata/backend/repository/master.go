// Package repository membaca dan menulis tabel master (modul Master Data) di Oracle. SQL dibangun dari daftar kolom
// `models.DaftarMaster` - nama tabel / kolom hanya dari konstanta itu, nilai selalu lewat bind posisi.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterdata/backend/models"
)

// ErrBarisTidakAda - ID tidak ada di tabel master.
var ErrBarisTidakAda = errors.New("repository: baris master tidak ada")

const (
	// StatusAktif / StatusNonaktif - isi kolom status (STS_AKTIF migrasi 760/761; OBJECTITEMTYPE.ISACTIVE `'1'` =
	// aktif `[terverifikasi]` RDBList GetObjectItembyName_SQL / GetDataObjectItem; `'0'` nonaktif `[dugaan]`).
	StatusAktif    = "1"
	StatusNonaktif = "0"
	// sequenceAkumulasi - urutan ID ACCUMULATION (prosedur RDBMASTERACCUMULATION `accumulation_seq.nextval`; DDL
	// sequence-nya tidak ada di korpus `[dugaan]` ada di skema).
	sequenceAkumulasi = "ACCUMULATION_SEQ"
	// TabelStatus - status master bertabel warisan tanpa kolom status (migrasi 761, MD-2); tanpa baris = aktif.
	TabelStatus = "T_MASTER_STATUS"
)

// Penyimpan - baca / tulis tabel master.
type Penyimpan interface {
	Daftar(ctx context.Context, t models.TabelMaster, kata, status string, nomor, ukuran int) (models.Halaman, error)
	Ada(ctx context.Context, t models.TabelMaster, id string) (bool, error)
	AdaCatatanAkumulasi(ctx context.Context, note, zip string) (bool, error)
	NilaiRujukan(ctx context.Context, r models.Rujukan, nilai string) (string, bool, error)
	IDAkumulasi(ctx context.Context, tx *db.Tx, negara, zip string) (string, error)
	Sisip(ctx context.Context, tx *db.Tx, t models.TabelMaster, b models.Baris) error
	Ubah(ctx context.Context, tx *db.Tx, t models.TabelMaster, b models.Baris) error
	UbahStatus(ctx context.Context, tx *db.Tx, t models.TabelMaster, id string, aktif bool) error
}

// MasterOracle - Penyimpan atas Oracle.
type MasterOracle struct{ db *db.DB }

// NewMasterOracle merakit penyimpan master.
func NewMasterOracle(d *db.DB) *MasterOracle { return &MasterOracle{db: d} }

// PolaCari - `%KATA%` huruf besar dengan %, _, \ di-escape; kosong = "".
func PolaCari(kata string) string {
	k := strings.ToUpper(strings.TrimSpace(kata))
	if k == "" {
		return ""
	}
	return "%" + strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(k) + "%"
}

// ekspresiStatus - nilai status baris (alias tabel `m`): kolom status tabel itu, atau - bila tabelnya warisan tanpa
// kolom status - baris T_MASTER_STATUS-nya (tanpa baris = aktif).
func ekspresiStatus(t models.TabelMaster, tabelStatus string) string {
	if t.KolomStatus != "" {
		return "m." + t.KolomStatus
	}
	return "NVL((SELECT s.STS_AKTIF FROM " + tabelStatus + " s WHERE s.NAMA_TABEL = '" + t.Nama + "' AND s.ID_BARIS = m.ID), '" +
		StatusAktif + "')"
}

// saring - syarat WHERE daftar: kata (Contains tidak peka huruf atas KolomCari, OR) dan status ("aktif" /
// "nonaktif" / "" semua) atas `eks` (ekspresiStatus). Bind dibangun bersama syaratnya.
func saring(t models.TabelMaster, eks, kata, status string) (string, []any) {
	var syarat []string
	var arg []any
	if pola := PolaCari(kata); pola != "" {
		var atau []string
		for _, k := range t.KolomCari {
			arg = append(arg, pola)
			atau = append(atau, "UPPER(m."+k+") LIKE :"+strconv.Itoa(len(arg))+` ESCAPE '\'`)
		}
		syarat = append(syarat, "("+strings.Join(atau, " OR ")+")")
	}
	switch status {
	case "aktif":
		syarat = append(syarat, eks+" = '"+StatusAktif+"'")
	case "nonaktif":
		syarat = append(syarat, "("+eks+" IS NULL OR "+eks+" <> '"+StatusAktif+"')")
	}
	if len(syarat) == 0 {
		return "", arg
	}
	return " WHERE " + strings.Join(syarat, " AND "), arg
}

func namaKolom(t models.TabelMaster, alias string) []string {
	out := make([]string, 0, len(t.Kolom))
	for _, k := range t.Kolom {
		out = append(out, alias+k.Nama)
	}
	return out
}

// sqlDaftar - kolom data + status, urut ID, halaman OFFSET / FETCH (dua bind terakhir).
func sqlDaftar(tabel string, t models.TabelMaster, eks, where string, nBind int) string {
	return "SELECT " + strings.Join(namaKolom(t, "m."), ", ") + ", " + eks + " FROM " + tabel + " m" + where +
		" ORDER BY m.ID OFFSET :" + strconv.Itoa(nBind+1) + " ROWS FETCH NEXT :" + strconv.Itoa(nBind+2) + " ROWS ONLY"
}

func sqlHitung(tabel, where string) string { return "SELECT COUNT(*) FROM " + tabel + " m" + where }

// sqlSisip - seluruh kolom data (urut daftar kolom) + status aktif bila tabelnya berkolom status (tanpa: tanpa
// baris T_MASTER_STATUS = aktif).
func sqlSisip(tabel string, t models.TabelMaster) string {
	kolom := namaKolom(t, "")
	nilai := make([]string, 0, len(t.Kolom)+1)
	for i := range t.Kolom {
		nilai = append(nilai, ":"+strconv.Itoa(i+1))
	}
	if t.KolomStatus != "" {
		kolom = append(kolom, t.KolomStatus)
		nilai = append(nilai, "'"+StatusAktif+"'")
	}
	return "INSERT INTO " + tabel + " (" + strings.Join(kolom, ", ") + ") VALUES (" + strings.Join(nilai, ", ") + ")"
}

// sqlUbah - kolom data selain ID (urut daftar kolom), ID bind terakhir.
func sqlUbah(tabel string, t models.TabelMaster) string {
	var set []string
	for i, k := range t.Kolom[1:] {
		set = append(set, k.Nama+" = :"+strconv.Itoa(i+1))
	}
	return "UPDATE " + tabel + " SET " + strings.Join(set, ", ") + " WHERE ID = :" + strconv.Itoa(len(t.Kolom))
}

func sqlStatus(tabel string, t models.TabelMaster) string {
	return "UPDATE " + tabel + " SET " + t.KolomStatus + " = :1 WHERE ID = :2"
}

// sqlUbahStatusTerpisah / sqlSisipStatusTerpisah - status di T_MASTER_STATUS: ubah barisnya, atau - bila belum ada
// - sisipkan (dua pernyataan di transaksi yang sama; bukan MERGE supaya setiap nama tabel terbaca penjaga ADR-U-0033).
func sqlUbahStatusTerpisah(tabelStatus string) string {
	return "UPDATE " + tabelStatus + " SET STS_AKTIF = :1 WHERE NAMA_TABEL = :2 AND ID_BARIS = :3"
}

func sqlSisipStatusTerpisah(tabelStatus string) string {
	return "INSERT INTO " + tabelStatus + " (STS_AKTIF, NAMA_TABEL, ID_BARIS) VALUES (:1, :2, :3)"
}

func sqlRujukan(tabel string, r models.Rujukan) string {
	nilai := "NULL"
	if r.KolomNilai != "" {
		nilai = "MAX(" + r.KolomNilai + ")"
	}
	return "SELECT COUNT(*), " + nilai + " FROM " + tabel + " WHERE " + r.KolomKunci + " = :1"
}

// sqlCatatanAkumulasi - prosedur RDBMASTERACCUMULATION `countid`: Note (huruf besar) + PostalCode sudah ada.
func sqlCatatanAkumulasi(tabel string) string {
	return "SELECT COUNT(*) FROM " + tabel + " WHERE UPPER(NOTE) = UPPER(:1) AND ZIPCODE = :2"
}

// sqlIDAkumulasi - prosedur RDBMASTERACCUMULATION, persis: `p_COUNTRY || '-' || p_ZIPCODE || '-' ||
// lpad(to_Char(accumulation_seq.nextval),6,'0')`.
func sqlIDAkumulasi(seq string) string {
	return "SELECT :1 || '-' || :2 || '-' || LPAD(TO_CHAR(" + seq + ".NEXTVAL), 6, '0') FROM DUAL"
}

func (r *MasterOracle) tabel(t models.TabelMaster) (string, error) { return r.db.Qualify(t.Nama) }

// Daftar - satu halaman master.
func (r *MasterOracle) Daftar(ctx context.Context, t models.TabelMaster, kata, status string, nomor, ukuran int) (models.Halaman, error) {
	tb, err := r.tabel(t)
	if err != nil {
		return models.Halaman{}, err
	}
	st, err := r.db.Qualify(TabelStatus)
	if err != nil {
		return models.Halaman{}, err
	}
	eks := ekspresiStatus(t, st)
	where, arg := saring(t, eks, kata, status)
	h := models.Halaman{Nomor: nomor, Ukuran: ukuran, Baris: []models.Baris{}, Aktif: []bool{}}
	if err := r.db.QueryRowContext(ctx, sqlHitung(tb, where), arg...).Scan(&h.Total); err != nil {
		return h, fmt.Errorf("repository: hitung %s: %w", t.Nama, err)
	}
	baris, err := r.db.QueryContext(ctx, sqlDaftar(tb, t, eks, where, len(arg)), append(arg, (nomor-1)*ukuran, ukuran)...)
	if err != nil {
		return h, fmt.Errorf("repository: baca %s: %w", t.Nama, err)
	}
	defer baris.Close()
	for baris.Next() {
		b, aktif, err := pindai(baris, t)
		if err != nil {
			return h, err
		}
		h.Baris = append(h.Baris, b)
		h.Aktif = append(h.Aktif, aktif)
	}
	return h, baris.Err()
}

func pindai(baris *sql.Rows, t models.TabelMaster) (models.Baris, bool, error) {
	nilai := make([]sql.NullString, len(t.Kolom)+1)
	tujuan := make([]any, len(nilai))
	for i := range nilai {
		tujuan[i] = &nilai[i]
	}
	if err := baris.Scan(tujuan...); err != nil {
		return nil, false, fmt.Errorf("repository: %s: %w", t.Nama, err)
	}
	b := models.Baris{}
	for i, k := range t.Kolom {
		b[k.JSON] = nilai[i].String
	}
	return b, nilai[len(t.Kolom)].String == StatusAktif, nil
}

// Ada - ID sudah ada.
func (r *MasterOracle) Ada(ctx context.Context, t models.TabelMaster, id string) (bool, error) {
	tb, err := r.tabel(t)
	if err != nil {
		return false, err
	}
	var n int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+tb+" WHERE ID = :1", id).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: %s: %w", t.Nama, err)
	}
	return n > 0, nil
}

// AdaCatatanAkumulasi - akumulasi ber-Note dan zip sama sudah ada (pemeriksaan prosedur sebelum INSERT).
func (r *MasterOracle) AdaCatatanAkumulasi(ctx context.Context, note, zip string) (bool, error) {
	tb, err := r.db.Qualify("ACCUMULATION")
	if err != nil {
		return false, err
	}
	var n int
	if err := r.db.QueryRowContext(ctx, sqlCatatanAkumulasi(tb), note, zip).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: ACCUMULATION: %w", err)
	}
	return n > 0, nil
}

// NilaiRujukan - nilai `KolomNilai` baris rujukan ber-kunci `nilai`; ada=false bila tidak ada.
func (r *MasterOracle) NilaiRujukan(ctx context.Context, ru models.Rujukan, nilai string) (string, bool, error) {
	tb, err := r.db.Qualify(ru.Tabel)
	if err != nil {
		return "", false, err
	}
	var n int
	var isi sql.NullString
	if err := r.db.QueryRowContext(ctx, sqlRujukan(tb, ru), nilai).Scan(&n, &isi); err != nil {
		return "", false, fmt.Errorf("repository: rujukan %s: %w", ru.Tabel, err)
	}
	return isi.String, n > 0, nil
}

// IDAkumulasi - ID akumulasi baru (sequence dimajukan di transaksi pemanggil).
func (r *MasterOracle) IDAkumulasi(ctx context.Context, tx *db.Tx, negara, zip string) (string, error) {
	seq, err := r.db.Qualify(sequenceAkumulasi)
	if err != nil {
		return "", err
	}
	var id string
	if err := tx.QueryRowContext(ctx, sqlIDAkumulasi(seq), negara, zip).Scan(&id); err != nil {
		return "", fmt.Errorf("repository: ID akumulasi: %w", err)
	}
	return id, nil
}

func argKolom(b models.Baris, kolom []models.Kolom) []any {
	arg := make([]any, 0, len(kolom))
	for _, k := range kolom {
		arg = append(arg, db.KosongJadiNil(b[k.JSON]))
	}
	return arg
}

// Sisip - baris baru berstatus aktif.
func (r *MasterOracle) Sisip(ctx context.Context, tx *db.Tx, t models.TabelMaster, b models.Baris) error {
	tb, err := r.tabel(t)
	if err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, sqlSisip(tb, t), argKolom(b, t.Kolom)...)
	if err != nil {
		return fmt.Errorf("repository: sisip %s: %w", t.Nama, err)
	}
	return db.PastikanSatuBaris(h, t.Nama)
}

// Ubah - kolom data selain ID; ErrBarisTidakAda bila nol baris.
func (r *MasterOracle) Ubah(ctx context.Context, tx *db.Tx, t models.TabelMaster, b models.Baris) error {
	tb, err := r.tabel(t)
	if err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, sqlUbah(tb, t), append(argKolom(b, t.Kolom[1:]), b["id"])...)
	if err != nil {
		return fmt.Errorf("repository: ubah %s: %w", t.Nama, err)
	}
	return satuBaris(h, t.Nama)
}

// UbahStatus - aktif / nonaktif; ErrBarisTidakAda bila nol baris.
func (r *MasterOracle) UbahStatus(ctx context.Context, tx *db.Tx, t models.TabelMaster, id string, aktif bool) error {
	tb, err := r.tabel(t)
	if err != nil {
		return err
	}
	nilai := StatusNonaktif
	if aktif {
		nilai = StatusAktif
	}
	if t.KolomStatus == "" {
		ada, err := r.Ada(ctx, t, id)
		if err != nil {
			return err
		}
		if !ada {
			return ErrBarisTidakAda
		}
		st, err := r.db.Qualify(TabelStatus)
		if err != nil {
			return err
		}
		h, err := tx.ExecContext(ctx, sqlUbahStatusTerpisah(st), nilai, t.Nama, id)
		if err != nil {
			return fmt.Errorf("repository: status %s: %w", t.Nama, err)
		}
		if n, err := h.RowsAffected(); err != nil || n > 0 {
			return err
		}
		if _, err := tx.ExecContext(ctx, sqlSisipStatusTerpisah(st), nilai, t.Nama, id); err != nil {
			return fmt.Errorf("repository: status %s: %w", t.Nama, err)
		}
		return nil
	}
	h, err := tx.ExecContext(ctx, sqlStatus(tb, t), nilai, id)
	if err != nil {
		return fmt.Errorf("repository: status %s: %w", t.Nama, err)
	}
	return satuBaris(h, t.Nama)
}

func satuBaris(h sql.Result, tabel string) error {
	n, err := h.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrBarisTidakAda
	}
	if n > 1 {
		return fmt.Errorf("repository: %s: %d baris berubah, mau 1 (ID ganda)", tabel, n)
	}
	return nil
}
