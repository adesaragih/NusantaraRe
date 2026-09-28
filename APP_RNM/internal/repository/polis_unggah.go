package repository

// Penyimpanan peserta hasil unggahan CSV - tiket 04 PremiumList Life.
//
// Untuk apa berkas ini: mengganti seluruh baris peserta satu polis dengan isi
// berkas yang baru diunggah, di dalam SATU transaksi.
//
// ⛔ MENGGANTI, BUKAN MENUMPUK (AC tiket 04). `DeleteTempUploadDataLife`
// membersihkan staging per case sebelum unggahan berikutnya; kami melakukan
// hal setara atas tabel peserta. Unggahan kedua yang menumpuk menghasilkan
// peserta ganda yang tidak seorang pun minta.
//
// ⛔ HANYA baris polis INI. Penyaringnya `PREMIUM_LIST_ID`, jadi tidak ada
// cara sebuah unggahan menyentuh baris polis lain - berbeda dari
// `CekDoubleInsured`, yang dua dari tiga pernyataannya LUPA menyaring
// `idpega` (lihat OQ-PL-05 di tiket 04).
//
// ⛔ SELURUH UANG DITULIS SEBAGAI TEKS DESIMAL lewat parameter, tidak pernah
// lewat `float64` (ADR-U-0003, ADR-U-0016).
//
// ⛔ Setiap query menyebut skemanya lewat `Qualify` (ADR-U-0033).
//
// Dibaca sesudah: polis_detail.go, polis_nomor.go.

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/internal/models"
)

// PengenalPesertaUnggah menyusun `ID` satu baris peserta.
//
// ⛔ NOL SEQUENCE, dan itu keputusan yang TERCATAT - tiket 00 PremiumList:
// *"Nol sequence. Pengenalnya dirakit di `repository` mengikuti pola
// `PengenalWorkBerikut`, bukan `DEFAULT seq.NEXTVAL` di DDL."* Migrasi
// 050-056 memang tidak membuat satu pun sequence untuk keluarga tabel ini,
// dan membuatnya sekarang berarti migrasi baru dari keputusan yang tidak ada.
//
// ⛔ DETERMINISTIK, dan itu justru yang diinginkan di sini - bukan kompromi.
// Unggahan MENGGANTI isi (hapus lalu sisip), jadi baris ke-N polis yang sama
// selalu memperoleh pengenal yang sama. Akibatnya unggah ulang tidak
// menghasilkan banjir pengenal baru yang menggantung, dan dua unggahan berkas
// yang sama menghasilkan baris yang sama persis - dapat dibandingkan.
//
// ⚠️ 32 heksa TEPAT, sebab kolomnya `VARCHAR2(32)`. Menyambung
// `polisID + "-" + nomor` akan MELEBIHI batas itu untuk polis yang
// pengenalnya sudah panjang - dan kelebihan satu karakter ditolak Oracle
// dengan ORA-12899 di baris yang tidak seorang pun tebak.
//
// ⚠️ MD5 dipakai sebagai PEMADAT, bukan sebagai pengaman. Isinya bukan
// rahasia dan tidak perlu tidak-dapat-ditebak: ia pengenal baris, bukan kunci
// penyimpanan. (Bandingkan `models.ImageIDBaru`, yang justru HARUS tidak
// dapat ditebak dan karena itu memakai GUID.)
func PengenalPesertaUnggah(polisID string, nomorBaris int) string {
	sum := md5.Sum([]byte(polisID + "\x00peserta\x00" + strconv.Itoa(nomorBaris)))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

// PesertaUnggah menyimpan peserta hasil unggahan.
type PesertaUnggah struct{ db *DB }

// NewPesertaUnggah menyusunnya.
func NewPesertaUnggah(db *DB) *PesertaUnggah { return &PesertaUnggah{db: db} }

// kolomSisipPeserta adalah kolom yang diisi unggahan, urut.
//
// ⛔ SATU DAFTAR merakit nama kolom, penanda parameter, DAN pengambilan
// nilainya. Tiga daftar terpisah akan berselisih, dan selisihnya menyimpan
// nilai sebuah kolom ke kolom lain - kekeliruan yang tidak satu pun galat
// tunjukkan sebab tipenya sama-sama teks atau sama-sama angka.
//
// ⚠️ Nama di kiri adalah kolom migrasi 052; nama di kanan kolom CSV. Keduanya
// SAMA untuk hampir semuanya, dan yang berbeda disebut.
var kolomSisipPeserta = []struct{ Kolom, DariCSV string }{
	{"POLICY_NO", "POLICY_NO"},
	{"CERTIFICATE_NO", "CERTIFICATE_NO"},
	{"NAME_OF_INSURED", "NAME_OF_INSURED"},
	{"POLICY_HOLDER", "POLICY_HOLDER"},
	{"PLAN", "PLAN"},
	{"CURRENCY", "CURRENCY"},
	{"MEDICAL_STATUS", "MEDICAL_STATUS"},
	{"RISK", "RISK"},
	{"SEX", "SEX"},
	{"DESCRIPTION", "DESCRIPTION"},
	{"PRO_RATE_TYPE", "PRO_RATE_TYPE"},
}

// kolomTanggalPeserta adalah kolom tanggal yang diisi unggahan.
//
// ⛔ Dikirim sebagai TEKS `dd/mm/yyyy` dan diubah Oracle lewat `TO_DATE`
// dengan pola yang DINYATAKAN. Menyerahkan bentuknya kepada `NLS_DATE_FORMAT`
// sesi membuat tanggal yang sama terbaca berbeda di dua mesin yang sama-sama
// benar - dan `03/04/2026` adalah tanggal yang sah dalam dua pembacaan.
//
// ⚠️ `STNC` dan `WPC` IKUT di sini walau namanya tidak berbunyi seperti
// tanggal: korpus memvalidasi keduanya dengan `@toDate(...)!=0`. Di migrasi
// 052 keduanya `VARCHAR2(255)` pada tabel detail - jadi keduanya TIDAK
// dibungkus `TO_DATE`, dan itu disebut di `kolomTeksTanggalPeserta`.
var kolomTanggalPeserta = []struct{ Kolom, DariCSV string }{
	{"DOB", "DOB"},
	{"BEGIN_DATE", "BEGIN_DATE"},
	{"EXPIRED_DATE", "EXPIRED_DATE"},
	{"EFFECTIVE_DATE", "EFFECTIVE_DATE"},
	{"LAPSE_DATE", "LAPSE_DATE"},
	{"GROSS_VALUATION_BEGIN_DATE", "GROSS_VALUATION_BEGIN_DATE"},
	{"GROSS_VALUATION_EXPIRED_DATE", "GROSS_VALUATION_EXPIRED_DATE"},
	{"RETRO_VALUATION_BEGIN_DATE", "RETROCESSION_VALUATION_BEGIN_DATE"},
	{"RETRO_VALUATION_EXPIRED_DATE", "RETROCESSION_VALUATION_EXPIRED_DATE"},
}

// kolomTeksTanggalPeserta adalah kolom yang divalidasi sebagai tanggal tetapi
// disimpan sebagai TEKS.
//
// ⛔ `STNC` dan `WPC` bertipe `VARCHAR2(255)` di migrasi 052 - lihat
// catatannya di `kolomTanggalPeserta`. Membungkusnya `TO_DATE` lalu menyimpan
// ke kolom teks akan menyerahkan bentuk simpannya kepada sesi.
var kolomTeksTanggalPeserta = []struct{ Kolom, DariCSV string }{
	{"STNC", "STNC"},
	{"WPC", "WPC"},
}

// kolomBulatPeserta adalah kolom bilangan bulat.
var kolomBulatPeserta = []struct{ Kolom, DariCSV string }{
	{"AGE", "AGE"},
	{"ENTRY_AGE", "ENTRY_AGE"},
	{"CURRENT_AGE", "CURRENT_AGE"},
	{"PERIOD_YY", "PERIOD_YY"},
	{"PERIOD_MM", "PERIOD_MM"},
	{"PASSED_PERIOD", "PASSED_PERIOD"},
}

// BentukTanggalOracle adalah pola `TO_DATE` untuk tanggal CSV.
//
// ⛔ DINYATAKAN, tidak diserahkan `NLS_DATE_FORMAT` sesi.
const BentukTanggalOracle = "DD/MM/YYYY"

// sqlHapusPesertaPolis merakit pembersihan baris peserta satu polis.
func sqlHapusPesertaPolis(detail string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE PREMIUM_LIST_ID = :1`, detail)
}

// HapusPesertaPolis membuang seluruh baris peserta polis ini.
//
// Mengembalikan cacah baris terhapus. NOL bukan galat: polis yang belum
// pernah diunggahi memang belum punya peserta.
func (r *PesertaUnggah) HapusPesertaPolis(ctx context.Context, tx *Tx,
	polisID string) (int, error) {

	detail, err := r.db.Qualify("T_PREMIUM_LIST_DETAIL")
	if err != nil {
		return 0, err
	}
	q := sqlHapusPesertaPolis(detail)
	if err := PeriksaSQL(q); err != nil {
		return 0, err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, polisID)
	if err != nil {
		return 0, fmt.Errorf("repository: menghapus peserta polis: %w", err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("repository: membaca cacah baris terhapus: %w", err)
	}
	return int(n), nil
}

// ekspresiSisipPeserta merakit daftar kolom dan penanda parameternya.
//
// Mengembalikan daftar kolom, daftar penanda, dan cacah parameter per baris.
func ekspresiSisipPeserta() (kolom, penanda string, cacah int) {
	var k, p []string
	tambah := func(nama, bentuk string) {
		cacah++
		k = append(k, nama)
		p = append(p, fmt.Sprintf(bentuk, cacah))
	}
	// ID dan PREMIUM_LIST_ID di depan - keduanya diisi kode, bukan CSV.
	tambah("ID", ":%d")
	tambah("PREMIUM_LIST_ID", ":%d")
	for _, c := range kolomSisipPeserta {
		tambah(c.Kolom, ":%d")
	}
	for _, c := range kolomTeksTanggalPeserta {
		tambah(c.Kolom, ":%d")
	}
	for _, c := range kolomTanggalPeserta {
		tambah(c.Kolom, "TO_DATE(:%d,'"+BentukTanggalOracle+"')")
	}
	for _, c := range kolomBulatPeserta {
		tambah(c.Kolom, ":%d")
	}
	for _, c := range models.KolomUangUnggah {
		tambah(c, ":%d")
	}
	return strings.Join(k, ", "), strings.Join(p, ", "), cacah
}

// sqlSisipPeserta merakit satu penyisipan baris peserta.
func sqlSisipPeserta(detail string) string {
	kolom, penanda, _ := ekspresiSisipPeserta()
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, detail, kolom, penanda)
}

// nilaiSisipPeserta menyusun argumen satu baris, urut sama dengan kolomnya.
//
// ⛔ Kosong menjadi NULL, bukan teks kosong maupun "0". Kolom yang belum diisi
// dan kolom bernilai nol adalah dua keadaan berbeda, dan hanya satu di
// antaranya perlu dikerjakan orang.
func nilaiSisipPeserta(id, polisID string, b models.BarisUnggah) []any {
	ambil := func(k string) any {
		v := strings.TrimSpace(b.Nilai[k])
		if v == "" {
			return nil
		}
		return v
	}
	arg := []any{id, polisID}
	for _, c := range kolomSisipPeserta {
		arg = append(arg, ambil(c.DariCSV))
	}
	for _, c := range kolomTeksTanggalPeserta {
		arg = append(arg, ambil(c.DariCSV))
	}
	for _, c := range kolomTanggalPeserta {
		arg = append(arg, ambil(c.DariCSV))
	}
	for _, c := range kolomBulatPeserta {
		arg = append(arg, ambil(c.DariCSV))
	}
	// ⛔ Uang dikirim sebagai TEKS DESIMAL apa adanya. Mengubahnya menjadi
	// `float64` di sini membuang angka di belakang koma tepat sebelum
	// tersimpan - sesudah seluruh validasi menyatakannya utuh.
	for _, c := range models.KolomUangUnggah {
		arg = append(arg, ambil(c))
	}
	return arg
}

// SisipPeserta menyisipkan seluruh baris unggahan.
//
// Mengembalikan cacah baris tersimpan.
func (r *PesertaUnggah) SisipPeserta(ctx context.Context, tx *Tx, polisID string,
	baris []models.BarisUnggah) (int, error) {

	detail, err := r.db.Qualify("T_PREMIUM_LIST_DETAIL")
	if err != nil {
		return 0, err
	}
	q := sqlSisipPeserta(detail)
	if err := PeriksaSQL(q); err != nil {
		return 0, err
	}
	// ⚠️ Satu pernyataan disiapkan SEKALI lalu dipakai berulang. Merakit
	// pernyataan baru per baris membuat Oracle menyusun rencana baru untuk
	// tiap baris, dan berkas seribu peserta menjadi seribu parse.
	stmt, err := tx.tx.PrepareContext(ctx, q)
	if err != nil {
		return 0, fmt.Errorf("repository: menyiapkan penyisipan peserta: %w", err)
	}
	defer stmt.Close()

	for _, b := range baris {
		id := PengenalPesertaUnggah(polisID, b.Nomor)
		if _, err := stmt.ExecContext(ctx, nilaiSisipPeserta(id, polisID, b)...); err != nil {
			// ⛔ Nomor barisnya IKUT di galat. "ORA-12899" tanpa nomor baris
			// menyuruh orang mencari sendiri baris mana dari seribu.
			return 0, fmt.Errorf("repository: menyisipkan peserta baris %d: %w",
				b.Nomor, err)
		}
	}
	return len(baris), nil
}
