package repository

// Peserta yang MASIH AKTIF di klaim Claim Life lain - aturan IN-APP Register (keputusan work owner 10-10-2026:
// "peserta berada di case lain, dan di case itu peserta masih outstanding atau sudah di accept. tapi jika peserta
// di case lain di reject. peserta itu bisa di input lagi di case yang sama atau di case lain").
//
// ⛔ BUKAN gandawarisan.go. Pemeriksa itu (Save to RNM langkah 11, tidak berubah) mencocokkan tabel datar WARISAN
// OS_AKSEPTASI_KLAIM_LIFE lewat CEDINGCO + nama + DOB dan sengaja MELEWATI baris berstatus NULL - sehingga dua
// pendaftaran aplikasi atas peserta yang sama sebelum Save to RNM tidak pernah tertangkap. Yang ini membaca tabel
// peserta klaim APLIKASI, T_CLAIMLF_PREMIUMLIST_DETAIL.
//
// Kunci peserta = PL_NUMBER + CERTIFICATE_NO (TRIM), BUKAN SOURCE_ID: SOURCE_ID menunjuk SATU versi baris sumber
// (AmbilUntukKlaim memilih versi hidup TERAKHIR, OQ-N14), jadi peserta yang sama sesudah endorsement punya SOURCE_ID
// lain dan akan lolos; baris pindahan warisan pun dapat tanpa SOURCE_ID. PL + sertifikat adalah identitas yang
// Register sendiri pakai untuk membaca peserta - dan setiap baris ber-SOURCE_ID sama pasti ber-PL + sertifikat sama.
//
// Status yang menentukan = status tingkat PESERTA, bukan header. Kode (kontrak.*) di model:
//
//	NULL  - terdaftar, belum diputuskan (Save to RNM menulis '0' ke baris ADJUSTMENT saja, TandaiBarisOutstanding)
//	'0'   - kontrak.KodeOutstanding
//	'1'   - kontrak.KodeAksep    (PerbaruiStatusBaris, SaveAdjustment)
//	'2'   - kontrak.KodeDitolak  (PerbaruiStatusBaris, RejectOSClaimLife)
//
// ⛔ MODEL BERBEDA DARI DUGAAN "satu status per baris peserta": keputusan ditulis ke baris ADJUSTMENT dan dicerminkan
// ke baris peserta, tetapi tombol `Add` (`Putaran.Tambah`, putaran sesudah ditolak) menyisip baris adjustment BARU
// berstatus kosong TANPA mengubah baris peserta - pesertanya tetap '2' padahal aktif lagi di putaran baru. Karena itu
// peserta dihitung AKTIF bila status baris pesertanya bukan Ditolak ATAU ia punya baris adjustment yang bukan
// Ditolak; hanya peserta yang ditolak di SEMUA tingkat yang boleh didaftarkan lagi. Label "Accepted" bila peserta atau
// salah satu barisnya Aksep; selain itu "Outstanding".
//
// Baris yang dicabut (STS_HAPUS, migrasi 022) tidak dihitung. Tidak ada "case dihapus" lunak di model (hapus
// klaim = hapus fisik berantai), jadi tidak ada saringan lain.
//
// ⚠️ KINERJA: tabel ini hanya ber-index CLAIM_ID (migrasi 003); kueri per PL memindai tabelnya. Index
// (PL_NUMBER, CERTIFICATE_NO) butuh migrasi - menunggu persetujuan work owner, tidak dibuat di sini.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
)

// PesertaAktifInApp - satu peserta yang sudah aktif di klaim lain.
type PesertaAktifInApp struct {
	NomorPremiList string
	Sertifikat     string
	KlaimID        string
	NomorKlaim     string
	// KodeStatus - kode terkuat peserta itu di klaim tersebut: kontrak.KodeAksep bila peserta atau salah satu baris
	// adjustment-nya Aksep; selain itu kode baris peserta ("" = belum diputuskan).
	KodeStatus string
}

// sqlPesertaAktifDiKlaimLain - peserta PL bersertifikat (daftar IN, n butir) di klaim mana pun, beserta dua cacah
// baris adjustment-nya: yang BUKAN ditolak dan yang Aksep. Kode diikat (tidak pernah literal - penjaga kode status).
//
// ⚠️ go-ora mengikat menurut URUTAN kemunculan: :1 = kode ditolak (subkueri 1), :2 = kode aksep (subkueri 2),
// :3 = PL, :4..:n+3 = sertifikat.
func sqlPesertaAktifDiKlaimLain(pes, hdr, adj string, n int) string {
	ikat := make([]string, n)
	for i := range ikat {
		ikat[i] = fmt.Sprintf(":%d", i+4)
	}
	return fmt.Sprintf(`SELECT TRIM(p.CERTIFICATE_NO), p.CLAIM_ID, h.CLAIM_NO, TRIM(p.STS_REJECT),
	        (SELECT COUNT(*) FROM %s a WHERE a.PREMIUM_LIST_DETAIL_ID = p.ID
	            AND (a.STS_REJECT IS NULL OR TRIM(a.STS_REJECT) <> :1)),
	        (SELECT COUNT(*) FROM %s a WHERE a.PREMIUM_LIST_DETAIL_ID = p.ID
	            AND TRIM(a.STS_REJECT) = :2)
	   FROM %s p JOIN %s h ON h.ID = p.CLAIM_ID
	  WHERE TRIM(p.PL_NUMBER) = :3
	    AND TRIM(p.CERTIFICATE_NO) IN (%s)
	    AND p.STS_HAPUS IS NULL
	  ORDER BY TRIM(p.CERTIFICATE_NO), p.CLAIM_ID`, adj, adj, pes, hdr, strings.Join(ikat, ", "))
}

// aktifInApp - aturan atas satu baris hasil: aktif bila peserta bukan Ditolak ATAU ada baris adjustment bukan
// Ditolak; kode terkuat Aksep bila peserta atau barisnya Aksep. Murni, diuji tanpa Oracle.
func aktifInApp(kodePeserta string, barisAktif, barisAksep int) (aktif bool, kode string) {
	kodePeserta = strings.TrimSpace(kodePeserta)
	aktif = kodePeserta != kontrak.KodeDitolak || barisAktif > 0
	kode = kodePeserta
	if kodePeserta == kontrak.KodeAksep || barisAksep > 0 {
		kode = kontrak.KodeAksep
	} else if kodePeserta == kontrak.KodeDitolak {
		kode = "" // peserta ditolak tetapi putaran baru masih berjalan = Outstanding
	}
	return aktif, kode
}

// PesertaAktifDiKlaimLain - SATU kueri untuk satu PL (daftar IN sertifikat; Register membatasi 500 peserta).
// Sertifikat dan PL dipangkas seperti Register membacanya.
func (r *KlaimLife) PesertaAktifDiKlaimLain(ctx context.Context, nomorPremiList string, sertifikat []string) (
	[]PesertaAktifInApp, error) {
	pl := strings.TrimSpace(nomorPremiList)
	if pl == "" || len(sertifikat) == 0 {
		return nil, nil
	}
	pes, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return nil, err
	}
	hdr, err := r.db.Qualify("T_GENERAL_CLAIM")
	if err != nil {
		return nil, err
	}
	adj, err := r.db.Qualify("T_CLAIMLF_ADJUSTMENT")
	if err != nil {
		return nil, err
	}
	q := sqlPesertaAktifDiKlaimLain(pes, hdr, adj, len(sertifikat))
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	args := []any{kontrak.KodeDitolak, kontrak.KodeAksep, pl}
	for _, s := range sertifikat {
		args = append(args, strings.TrimSpace(s))
	}
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca peserta aktif di klaim lain: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []PesertaAktifInApp
	for rows.Next() {
		var sert, klaim, nomor, kode sql.NullString
		var barisAktif, barisAksep int
		if err := rows.Scan(&sert, &klaim, &nomor, &kode, &barisAktif, &barisAksep); err != nil {
			return nil, fmt.Errorf("repository: memindai peserta aktif di klaim lain: %w", err)
		}
		aktif, kodeKuat := aktifInApp(kode.String, barisAktif, barisAksep)
		if !aktif {
			continue
		}
		out = append(out, PesertaAktifInApp{NomorPremiList: pl, Sertifikat: sert.String, KlaimID: klaim.String,
			NomorKlaim: nomor.String, KodeStatus: kodeKuat})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca peserta aktif di klaim lain: %w", err)
	}
	return out, nil
}
