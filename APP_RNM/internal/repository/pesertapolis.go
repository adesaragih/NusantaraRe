package repository

// Pembaca peserta polis dari tabel warisan - tiket 02.
//
// Untuk apa berkas ini: layar Register memilih peserta dari
// POOLDATA.M_LIFE_PREMIUM_DETAIL, tabel warisan milik modul PremiumList Life.
//
// ⛔ TABEL ITU BERISI 66,8 JUTA BARIS. Katalog instance pengembangan
// (26-09-2026) menyebut index pada PL_NUMBER, CERTIFICATE_NO, dan POLICY_NO -
// dan TIDAK ada index yang berawalan EDMSTATUS. Query tanpa penyaring
// ber-index karena itu bukan "agak lambat", melainkan pemindaian penuh atas
// 66 juta baris yang menahan basis data produksi.
//
// Dua aturan yang lahir dari itu, dan dijaga test statik:
//  1. setiap query ke tabel ini menyaring dengan PL_NUMBER atau CERTIFICATE_NO
//  2. setiap query berbatas hasil (FETCH FIRST :n ROWS ONLY)
//
// Dibaca sesudah: pohonklaim.go.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// namaTabelPeserta adalah tabel warisan peserta polis.
const namaTabelPeserta = "M_LIFE_PREMIUM_DETAIL"

// batasHasilBawaan dipakai bila pemanggil tidak menyebut batas.
const batasHasilBawaan = 200

// batasHasilTertinggi menahan permintaan yang batasnya kelewat besar.
const batasHasilTertinggi = 1000

// PesertaPolis membaca calon peserta klaim dari tabel warisan.
type PesertaPolis struct{ db *DB }

// NewPesertaPolis menyusun pembacanya.
func NewPesertaPolis(db *DB) *PesertaPolis { return &PesertaPolis{db: db} }

// CalonPeserta adalah satu baris hasil pencarian.
//
// Sengaja TIDAK memuat seluruh 85 kolom tabel itu: yang dibawa hanya yang
// dipakai layar untuk memilih, dan ⛔ kolom KTP tidak pernah ikut.
type CalonPeserta struct {
	NomorPremiList  string
	NomorPolis      string
	NomorSertifikat string
	NamaTertanggung string
	MataUang        string
	EDMStatus       string
}

// penyaringHidup menyingkirkan peserta yang sudah batal atau dihapus lunak.
//
// `[terverifikasi]` verdict V14 grilling Endorsement Life: baris negatif hasil
// jurnal balik ditulis ke tabel yang SAMA dan hidup berdampingan dengan baris
// positifnya - akuntansi melihat keduanya, klaim hanya melihat yang masih
// hidup. Kontraknya berbunyi "peserta yang sudah EDM Batal atau soft-delete
// TIDAK BOLEH MUNCUL di Claim Life".
//
// ⚠️ Penyaring naif `EDMSTATUS = ”` KELIRU dan itu bukan kehati-hatian
// berlebihan: agregat instance pengembangan menghitung 59,1 juta baris
// ber-EDMSTATUS NULL dan NOL baris bernilai teks kosong. Peserta new business
// justru yang NULL, sehingga penyaring naif membuang hampir seluruh tabel.
//
// ⚠️ TRIM di kedua sisi. Tanpa itu, SQL dan PesertaHidup menjadi DUA aturan
// yang berbeda: "Batal " berspasi lolos SQL tetapi ditolak Go, dan AC 29 -
// "penyaringan terjadi di satu tempat" - dilanggar tanpa satu pun test gagal.
const penyaringHidup = `(EDMSTATUS IS NULL OR TRIM(EDMSTATUS) NOT IN ('Batal','Delete'))`

// Cari mengembalikan calon peserta satu premium list.
//
// Batas hasil WAJIB: layar tidak pernah memerlukan 66 juta baris, dan query
// tanpa batas adalah cara paling mudah menahan basis data tanpa sengaja.
func (r *PesertaPolis) Cari(ctx context.Context, nomorPremiList string, batas int) (
	[]CalonPeserta, error) {
	if strings.TrimSpace(nomorPremiList) == "" {
		return nil, fmt.Errorf("repository: mencari peserta tanpa nomor premium list; " +
			"tabel peserta hanya ber-index pada PL_NUMBER, CERTIFICATE_NO, dan POLICY_NO")
	}
	if batas <= 0 {
		batas = batasHasilBawaan
	}
	if batas > batasHasilTertinggi {
		batas = batasHasilTertinggi
	}
	tabel, err := r.db.Qualify(namaTabelPeserta)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT PL_NUMBER, POLICY_NO, CERTIFICATE_NO, NAME_OF_INSURED, CURRENCY, EDMSTATUS
	        FROM %s
	        WHERE PL_NUMBER = :1 AND %s
	        ORDER BY CERTIFICATE_NO
	        FETCH FIRST %d ROWS ONLY`, tabel, penyaringHidup, batas)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	baris, err := r.db.sql.QueryContext(ctx, q, nomorPremiList)
	if err != nil {
		return nil, fmt.Errorf("repository: mencari peserta %s: %w", nomorPremiList, err)
	}
	defer func() { _ = baris.Close() }()

	var out []CalonPeserta
	for baris.Next() {
		var pl, polis, sertifikat, nama, mataUang, edm sql.NullString
		if err := baris.Scan(&pl, &polis, &sertifikat, &nama, &mataUang, &edm); err != nil {
			return nil, fmt.Errorf("repository: membaca peserta %s: %w", nomorPremiList, err)
		}
		// ⛔ Disaring DUA KALI dengan sengaja, dan itu bukan pemborosan.
		// Penyaring SQL ada supaya Oracle tidak mengirim 66 juta baris;
		// PesertaHidup adalah ATURANNYA, dan ia yang menentukan. Tanpa
		// pemanggil produksi, aturan di Go dan aturan di SQL dapat berselisih
		// diam-diam - dan yang terlihat hanya salah satunya.
		if !PesertaHidup(edm.String) {
			continue
		}
		out = append(out, CalonPeserta{
			NomorPremiList:  pl.String,
			NomorPolis:      polis.String,
			NomorSertifikat: sertifikat.String,
			NamaTertanggung: nama.String,
			MataUang:        mataUang.String,
			EDMStatus:       edm.String,
		})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca peserta %s: %w", nomorPremiList, err)
	}
	return out, nil
}

// PesertaHidup memutuskan apakah satu nilai EDMSTATUS masih boleh tampil.
//
// Dipisah dari SQL supaya aturannya dapat diuji tanpa Oracle, dan supaya
// penyaring itu punya SATU rumah - AC 29 tiket 02 menuntut penyaringan terjadi
// di satu tempat, dan test yang menemukan jalur kedua gagal.
func PesertaHidup(edmStatus string) bool {
	switch strings.TrimSpace(edmStatus) {
	case "Batal", "Delete":
		return false
	default:
		// NULL, teks kosong, Old, New, dan nilai tak dikenal tetap tampil.
		// ⛔ Nilai tak dikenal sengaja LOLOS: menyembunyikan peserta karena
		// status yang belum kita pahami jauh lebih berbahaya daripada
		// menampilkan satu baris berlebih, sebab yang hilang tidak terlihat.
		return true
	}
}
