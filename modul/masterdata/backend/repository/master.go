// Package repository membaca dan menulis tabel master (modul Master Data) di Oracle. SQL dibangun dari daftar kolom
// `models.DaftarMaster` - nama tabel / kolom hanya dari konstanta itu, nilai selalu lewat bind posisi.
//
// Setiap kueri membaca tabel master beralias `m`; tabel JejakTerpisah (NATION, OBJECTITEMTYPE - warisan, MD-2)
// di-LEFT JOIN ke T_MASTER_STATUS beralias `s` untuk jejak ubah (MD-7) dan - bila tanpa kolom status - statusnya.
// Tanggal ditulis SYSDATE di SQL dan dibaca TO_CHAR berformat tetap: tanpa zona / NLS di Go.
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
	// TabelStatus - status dan jejak ubah master bertabel warisan (migrasi 761 / 762, MD-2); tanpa baris = aktif.
	TabelStatus = "T_MASTER_STATUS"
	// formatTanggal - TO_CHAR kolom jejak ubah.
	formatTanggal = "YYYY-MM-DD HH24:MI:SS"
)

// Penyimpan - baca / tulis tabel master. `akun` = pelaku (jejak ubah, MD-7).
type Penyimpan interface {
	Daftar(ctx context.Context, t models.TabelMaster, kata, status string, nomor, ukuran int) (models.Halaman, error)
	Ada(ctx context.Context, t models.TabelMaster, id string) (bool, error)
	AdaCatatanAkumulasi(ctx context.Context, note, zip string) (bool, error)
	NilaiRujukan(ctx context.Context, r models.Rujukan, nilai string) (string, bool, error)
	IDAkumulasi(ctx context.Context, tx *db.Tx, negara, zip string) (string, error)
	Sisip(ctx context.Context, tx *db.Tx, t models.TabelMaster, b models.Baris, akun string) error
	Ubah(ctx context.Context, tx *db.Tx, t models.TabelMaster, b models.Baris, akun string) error
	UbahStatus(ctx context.Context, tx *db.Tx, t models.TabelMaster, id string, aktif bool, akun string) error
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

// sumber - klausa FROM: master `m` (+ T_MASTER_STATUS `s` bila JejakTerpisah).
func sumber(t models.TabelMaster, tabel, tabelStatus string) string {
	if !t.JejakTerpisah {
		return tabel + " m"
	}
	return tabel + " m LEFT JOIN " + tabelStatus + " s ON s.NAMA_TABEL = '" + t.Nama + "' AND s.ID_BARIS = m.ID"
}

// ekspresiStatus - nilai status baris: kolom status tabel itu, atau STS_AKTIF T_MASTER_STATUS (tanpa baris = aktif).
func ekspresiStatus(t models.TabelMaster) string {
	if t.KolomStatus != "" {
		return "m." + t.KolomStatus
	}
	return "NVL(s.STS_AKTIF, '" + StatusAktif + "')"
}

// kolomJejak - ekspresi terpilih keempat kolom jejak ubah.
func kolomJejak(t models.TabelMaster) []string {
	a := "m."
	if t.JejakTerpisah {
		a = "s."
	}
	out := make([]string, 0, len(models.KolomAudit))
	for _, k := range models.KolomAudit {
		if k.Tanggal {
			out = append(out, "TO_CHAR("+a+k.Nama+", '"+formatTanggal+"')")
			continue
		}
		out = append(out, a+k.Nama)
	}
	return out
}

// saring - syarat WHERE daftar: kata (Contains tidak peka huruf atas KolomCari, OR) dan status ("aktif" /
// "nonaktif" / "" semua). Bind dibangun bersama syaratnya.
func saring(t models.TabelMaster, kata, status string) (string, []any) {
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
	eks := ekspresiStatus(t)
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

// sqlDaftar - kolom data, jejak ubah, status; urut ID, halaman OFFSET / FETCH (dua bind terakhir).
func sqlDaftar(t models.TabelMaster, dari, where string, nBind int) string {
	pilih := append(append(namaKolom(t, "m."), kolomJejak(t)...), ekspresiStatus(t))
	return "SELECT " + strings.Join(pilih, ", ") + " FROM " + dari + where + " ORDER BY m.ID OFFSET :" + strconv.Itoa(nBind+1) +
		" ROWS FETCH NEXT :" + strconv.Itoa(nBind+2) + " ROWS ONLY"
}

func sqlHitung(dari, where string) string { return "SELECT COUNT(*) FROM " + dari + where }

// sqlSisip - seluruh kolom data (urut daftar kolom) + status aktif bila tabelnya berkolom status + pembuat / tanggal
// buat bila jejaknya di tabel itu (bind terakhir = pelaku).
func sqlSisip(tabel string, t models.TabelMaster) string {
	kolom := namaKolom(t, "")
	nilai := make([]string, 0, len(t.Kolom)+3)
	for i := range t.Kolom {
		nilai = append(nilai, ":"+strconv.Itoa(i+1))
	}
	if t.KolomStatus != "" {
		kolom = append(kolom, t.KolomStatus)
		nilai = append(nilai, "'"+StatusAktif+"'")
	}
	if !t.JejakTerpisah {
		kolom = append(kolom, "CREATE_OP", "TGL_CREATE")
		nilai = append(nilai, ":"+strconv.Itoa(len(t.Kolom)+1), "SYSDATE")
	}
	return "INSERT INTO " + tabel + " (" + strings.Join(kolom, ", ") + ") VALUES (" + strings.Join(nilai, ", ") + ")"
}

// sqlUbah - kolom data selain ID (urut daftar kolom) [+ pengubah / tanggal ubah], ID bind terakhir.
func sqlUbah(tabel string, t models.TabelMaster) string {
	var set []string
	for i, k := range t.Kolom[1:] {
		set = append(set, k.Nama+" = :"+strconv.Itoa(i+1))
	}
	n := len(t.Kolom)
	if !t.JejakTerpisah {
		set = append(set, "UPDATE_OP = :"+strconv.Itoa(n), "TGL_UPDATE = SYSDATE")
		n++
	}
	return "UPDATE " + tabel + " SET " + strings.Join(set, ", ") + " WHERE ID = :" + strconv.Itoa(n)
}

// sqlStatus - kolom status [+ pengubah / tanggal ubah]: :1 status, [:2 pelaku], ID bind terakhir.
func sqlStatus(tabel string, t models.TabelMaster) string {
	if t.JejakTerpisah {
		return "UPDATE " + tabel + " SET " + t.KolomStatus + " = :1 WHERE ID = :2"
	}
	return "UPDATE " + tabel + " SET " + t.KolomStatus + " = :1, UPDATE_OP = :2, TGL_UPDATE = SYSDATE WHERE ID = :3"
}

// sqlUbahJejakTerpisah / sqlSisipJejakTerpisah - jejak (dan status bila `denganStatus`) di T_MASTER_STATUS: ubah
// barisnya, atau - bila belum ada - sisipkan (dua pernyataan, bukan MERGE: penjaga ADR-U-0033). Bind: [status],
// pelaku, nama tabel, ID.
func sqlUbahJejakTerpisah(tabelStatus string, denganStatus bool) string {
	set, n := "", 1
	if denganStatus {
		set, n = "STS_AKTIF = :1, ", 2
	}
	return "UPDATE " + tabelStatus + " SET " + set + "UPDATE_OP = :" + strconv.Itoa(n) + ", TGL_UPDATE = SYSDATE WHERE NAMA_TABEL = :" +
		strconv.Itoa(n+1) + " AND ID_BARIS = :" + strconv.Itoa(n+2)
}

func sqlSisipJejakTerpisah(tabelStatus string, denganStatus bool) string {
	if denganStatus {
		return "INSERT INTO " + tabelStatus + " (STS_AKTIF, UPDATE_OP, TGL_UPDATE, NAMA_TABEL, ID_BARIS) VALUES (:1, :2, SYSDATE, :3, :4)"
	}
	return "INSERT INTO " + tabelStatus + " (UPDATE_OP, TGL_UPDATE, NAMA_TABEL, ID_BARIS) VALUES (:1, SYSDATE, :2, :3)"
}

// sqlUbahBuatJejakTerpisah / sqlBuatJejakTerpisah - baris T_MASTER_STATUS saat master warisan ditambah: pembuat /
// tanggal buat, status aktif. Baris sisa ber-(tabel, ID) sama - baris masternya dihapus di luar menu lalu ID-nya
// ditambah lagi - DIATUR ULANG, tidak mewarisi status / pengubah lama dan tidak menabrak PK. Bind: pelaku, nama
// tabel, ID.
func sqlUbahBuatJejakTerpisah(tabelStatus string) string {
	return "UPDATE " + tabelStatus + " SET STS_AKTIF = '" + StatusAktif + "', CREATE_OP = :1, TGL_CREATE = SYSDATE, " +
		"UPDATE_OP = NULL, TGL_UPDATE = NULL WHERE NAMA_TABEL = :2 AND ID_BARIS = :3"
}

func sqlBuatJejakTerpisah(tabelStatus string) string {
	return "INSERT INTO " + tabelStatus + " (CREATE_OP, TGL_CREATE, NAMA_TABEL, ID_BARIS) VALUES (:1, SYSDATE, :2, :3)"
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

// namaBerskema - nama berskema tabel master dan T_MASTER_STATUS.
func (r *MasterOracle) namaBerskema(t models.TabelMaster) (tabel, tabelStatus string, err error) {
	if tabel, err = r.db.Qualify(t.Nama); err != nil {
		return "", "", err
	}
	tabelStatus, err = r.db.Qualify(TabelStatus)
	return tabel, tabelStatus, err
}

// Daftar - satu halaman master.
func (r *MasterOracle) Daftar(ctx context.Context, t models.TabelMaster, kata, status string, nomor, ukuran int) (models.Halaman, error) {
	tabel, tabelStatus, err := r.namaBerskema(t)
	if err != nil {
		return models.Halaman{}, err
	}
	dari := sumber(t, tabel, tabelStatus)
	where, arg := saring(t, kata, status)
	h := models.Halaman{Nomor: nomor, Ukuran: ukuran, Baris: []models.Baris{}, Aktif: []bool{}}
	if err := r.db.QueryRowContext(ctx, sqlHitung(dari, where), arg...).Scan(&h.Total); err != nil {
		return h, fmt.Errorf("repository: hitung %s: %w", t.Nama, err)
	}
	baris, err := r.db.QueryContext(ctx, sqlDaftar(t, dari, where, len(arg)), append(arg, (nomor-1)*ukuran, ukuran)...)
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

// pindai - kolom data, keempat kolom jejak, lalu status.
func pindai(baris *sql.Rows, t models.TabelMaster) (models.Baris, bool, error) {
	kolom := t.SeluruhKolom()
	nilai := make([]sql.NullString, len(kolom)+1)
	tujuan := make([]any, len(nilai))
	for i := range nilai {
		tujuan[i] = &nilai[i]
	}
	if err := baris.Scan(tujuan...); err != nil {
		return nil, false, fmt.Errorf("repository: %s: %w", t.Nama, err)
	}
	b := models.Baris{}
	for i, k := range kolom {
		b[k.JSON] = nilai[i].String
	}
	return b, nilai[len(kolom)].String == StatusAktif, nil
}

// Ada - ID sudah ada.
func (r *MasterOracle) Ada(ctx context.Context, t models.TabelMaster, id string) (bool, error) {
	tb, err := r.db.Qualify(t.Nama)
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

// Sisip - baris baru berstatus aktif, pembuat = `akun`.
func (r *MasterOracle) Sisip(ctx context.Context, tx *db.Tx, t models.TabelMaster, b models.Baris, akun string) error {
	tabel, tabelStatus, err := r.namaBerskema(t)
	if err != nil {
		return err
	}
	arg := argKolom(b, t.Kolom)
	if !t.JejakTerpisah {
		arg = append(arg, akun)
	}
	h, err := tx.ExecContext(ctx, sqlSisip(tabel, t), arg...)
	if err != nil {
		return fmt.Errorf("repository: sisip %s: %w", t.Nama, err)
	}
	if err := db.PastikanSatuBaris(h, t.Nama); err != nil {
		return err
	}
	if !t.JejakTerpisah {
		return nil
	}
	return ubahAtauSisip(ctx, tx, t.Nama, sqlUbahBuatJejakTerpisah(tabelStatus), sqlBuatJejakTerpisah(tabelStatus),
		akun, t.Nama, b["id"])
}

// Ubah - kolom data selain ID, pengubah = `akun`; ErrBarisTidakAda bila nol baris.
func (r *MasterOracle) Ubah(ctx context.Context, tx *db.Tx, t models.TabelMaster, b models.Baris, akun string) error {
	tabel, tabelStatus, err := r.namaBerskema(t)
	if err != nil {
		return err
	}
	arg := argKolom(b, t.Kolom[1:])
	if !t.JejakTerpisah {
		arg = append(arg, akun)
	}
	h, err := tx.ExecContext(ctx, sqlUbah(tabel, t), append(arg, b["id"])...)
	if err != nil {
		return fmt.Errorf("repository: ubah %s: %w", t.Nama, err)
	}
	if err := satuBaris(h, t.Nama); err != nil {
		return err
	}
	if !t.JejakTerpisah {
		return nil
	}
	return r.jejakTerpisah(ctx, tx, tabelStatus, t, b["id"], nil, akun)
}

// UbahStatus - aktif / nonaktif, pengubah = `akun`; ErrBarisTidakAda bila ID tidak ada.
func (r *MasterOracle) UbahStatus(ctx context.Context, tx *db.Tx, t models.TabelMaster, id string, aktif bool, akun string) error {
	tabel, tabelStatus, err := r.namaBerskema(t)
	if err != nil {
		return err
	}
	nilai := StatusNonaktif
	if aktif {
		nilai = StatusAktif
	}
	switch {
	case t.KolomStatus == "":
		ada, err := r.Ada(ctx, t, id)
		if err != nil {
			return err
		}
		if !ada {
			return ErrBarisTidakAda
		}
		return r.jejakTerpisah(ctx, tx, tabelStatus, t, id, &nilai, akun)
	case t.JejakTerpisah:
		h, err := tx.ExecContext(ctx, sqlStatus(tabel, t), nilai, id)
		if err != nil {
			return fmt.Errorf("repository: status %s: %w", t.Nama, err)
		}
		if err := satuBaris(h, t.Nama); err != nil {
			return err
		}
		return r.jejakTerpisah(ctx, tx, tabelStatus, t, id, nil, akun)
	}
	h, err := tx.ExecContext(ctx, sqlStatus(tabel, t), nilai, akun, id)
	if err != nil {
		return fmt.Errorf("repository: status %s: %w", t.Nama, err)
	}
	return satuBaris(h, t.Nama)
}

// jejakTerpisah - pengubah (dan status bila `status` terisi) baris T_MASTER_STATUS master warisan.
func (r *MasterOracle) jejakTerpisah(ctx context.Context, tx *db.Tx, tabelStatus string, t models.TabelMaster, id string, status *string, akun string) error {
	var arg []any
	if status != nil {
		arg = append(arg, *status)
	}
	arg = append(arg, akun, t.Nama, id)
	return ubahAtauSisip(ctx, tx, t.Nama, sqlUbahJejakTerpisah(tabelStatus, status != nil),
		sqlSisipJejakTerpisah(tabelStatus, status != nil), arg...)
}

// ubahAtauSisip - jalankan `ubah`; bila nol baris, `sisip` dengan bind yang sama (pengganti MERGE).
func ubahAtauSisip(ctx context.Context, tx *db.Tx, tabel, ubah, sisip string, arg ...any) error {
	h, err := tx.ExecContext(ctx, ubah, arg...)
	if err != nil {
		return fmt.Errorf("repository: jejak %s: %w", tabel, err)
	}
	n, err := h.RowsAffected()
	if err != nil || n > 0 {
		return err
	}
	if _, err := tx.ExecContext(ctx, sisip, arg...); err != nil {
		return fmt.Errorf("repository: jejak %s: %w", tabel, err)
	}
	return nil
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
