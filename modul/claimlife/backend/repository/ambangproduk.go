package repository

// Ambang batas hari per produk - butir ba, dibuka keputusan bh.
//
// Untuk apa berkas ini: membaca kedua ambang yang dipakai
// `models.PenandaBatasHari`, berkunci `ProductNameID`.
//
// `[terverifikasi]` `Claim Life/RDBList/GetProductName.xml`:
//
//	SELECT * FROM POOLDATA.PRODUCTINWARD_LIFE
//	 WHERE ID = {pyWorkPage.PolicyDataLife.ProductNameID}
//
// ⛔ RIWAYATNYA DICATAT, sebab ia menjelaskan bentuk berkas ini. Pembaca ini
// sempat ditulis 28-09-2026 lalu DIBUANG ketika `TestMasterViewTidakDisentuh`
// berbunyi: view ini ada di daftar master yang tidak disentuh. Penjaga itu
// bekerja sebagaimana mestinya, dan executor tidak melonggarkannya sendiri.
//
// `[DIPUTUSKAN - butir bh, veto work owner 28-09-2026]` Ia dibuka dengan
// batas yang sempit, dan batas itulah yang membuat pembukaannya aman:
//
//	AC 38 melarang MENGURAI `JSONDATA` di aplikasi - bukan membaca kolom
//	BERTIPE dari view milik basis data. Pega sendiri membaca view ini.
//
// Karena itu yang diizinkan TEPAT: satu pembaca, READ-ONLY, DUA kolom,
// berkunci `ID`. Nol `JSONDATA` di kode, nol `m_product_life`, nol tulisan.
// Penjaganya dipersempit untuk mengizinkan persis itu - nama fungsinya dan
// kedua kolomnya dikunci - sehingga pembaca ketiga atau kolom ketiga akan
// tetap berbunyi.
//
// ⛔ `SELECT *` TIDAK DITIRU, dan di sini itu bukan soal kerapian melainkan
// soal batas yang baru saja dinyatakan: view-nya 40 kolom `[data DBA]`, dan
// membawa seluruhnya berarti membawa kolom yang izin ini TIDAK mencakup.
//
// ⚠️ OQ-001 tetap terbuka (pertanyaan DDL produksi); berkas ini tidak
// menunggunya.
//
// Dibaca sesudah: polis_ringkas.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// ErrAmbangProdukTakDitemukan - produk itu tidak ada di view.
//
// ⛔ GAGAL TERANG, bukan ambang bawaan. Ambang yang ditebak MELARANG atau
// MELOLOSKAN klaim dengan angka karangan, dan keduanya keputusan yang tidak
// seorang pun buat.
var ErrAmbangProdukTakDitemukan = errors.New(
	"repository: produk tidak ditemukan di view produk; ambang batas hari " +
		"tidak dapat ditentukan, dan TIDAK ada nilai pengganti yang dipakai")

// KolomAmbangProduk adalah DUA kolom yang izin butir bh cakup.
//
// ⛔ DIDAFTAR SEBAGAI DATA, bukan hanya ditulis di dalam query. Penjaga
// `TestMasterViewTidakDisentuh` membacanya untuk memastikan yang dibaca tepat
// kedua kolom ini - kolom ketiga akan berbunyi.
var KolomAmbangProduk = []string{"MAXEXPIREDCLAIM", "MAXDATARECEIVE"}

// NamaViewProdukLife adalah view yang dibaca.
//
// ⛔ DITULIS UTUH, DAN ITU DISENGAJA. Ronde pertama berkas ini merakitnya dari
// potongan (`"PRODUCTINWARD" + "_LIFE"`) supaya penjaga master tidak
// menemukannya - dan itu MENGELABUI penjaga, bukan mempersempitnya. Penjaga
// yang dikelabui satu kali menjadi penjaga yang buta selamanya, sebab
// pembaca berikutnya akan menyalin caranya.
//
// Yang benar: nama ditulis apa adanya, dan penjaganya dipersempit di tempat
// penjaga itu tinggal - dengan berkas, fungsi, dan kedua kolomnya dikunci.
const NamaViewProdukLife = "PRODUCTINWARD_LIFE"

// AmbangProduk adalah kedua ambang yang dipakai validasi tanggal.
type AmbangProduk struct {
	// MaxExpiredClaim membatasi jarak DATE_OF_LOSS -> CLAIM_RECEIVED_DATE.
	MaxExpiredClaim string
	// MaxDataReceive membatasi jarak EFFECTIVE_DATE -> RECEIVED_DATE.
	MaxDataReceive string
}

// Lengkap menjawab apakah kedua ambang terbaca.
//
// ⚠️ Ambang yang KOSONG bukan ambang nol. Produk yang kolomnya belum diisi
// tidak berarti "tidak boleh terlambat satu hari pun".
func (a AmbangProduk) Lengkap() bool {
	return strings.TrimSpace(a.MaxExpiredClaim) != "" &&
		strings.TrimSpace(a.MaxDataReceive) != ""
}

// ProdukLife membaca ambang per produk. READ-ONLY, seluruhnya.
type ProdukLife struct{ db *db.DB }

// NewProdukLife menyusunnya.
func NewProdukLife(db *db.DB) *ProdukLife { return &ProdukLife{db: db} }

// TabelProdukNameLife - tabel DATAR Master Product Name Life (modul masterproductnamelife, migrasi 140) - sumber ambang
// SEJAK 11-10-2026: aplikasi Master menulis tabel ini, bukan JSON di balik view (view DEV "tidak dipakai dan tidak dibangun
// ulang", masterproductnamelife mpnl_tabel.go). `MAXEXPIREDCLAIM` = "Max Notification Claim Expired", `MAXDATARECEIVE` =
// "Max Production Data Receive", keduanya NUMBER(5) hari. BACA SAJA, dua kolom yang sama, berkunci `ID` (= ProductNameID
// polis). Salinan bentuk baca, tanpa impor modul itu.
const TabelProdukNameLife = "M_PRODUCTNAME_LIFE"

// sqlAmbangProdukDatar - kedua ambang dari tabel datar Master (sumber utama).
func sqlAmbangProdukDatar(tabel string) string {
	return fmt.Sprintf(`SELECT t.%s, t.%s FROM %s t WHERE t.ID = :1`,
		KolomAmbangProduk[0], KolomAmbangProduk[1], tabel)
}

// sqlAmbangProduk merakit pembacaan kedua ambang dari view warisan (cadangan: produk yang belum ada di tabel datar).
//
// ⚠️ Keduanya dibaca sebagai TEKS, tidak diubah menjadi bilangan di sini:
// `models.PenandaBatasHari` yang menguraikannya, dan konversi tipe dilakukan
// SEKALI saat masuk (ADR-U-0022).
func sqlAmbangProduk(produk string) string {
	return fmt.Sprintf(`SELECT v.%s, v.%s FROM %s v WHERE v.ID = :1`,
		KolomAmbangProduk[0], KolomAmbangProduk[1], produk)
}

// Ambang membaca kedua ambang sebuah produk.
//
// ⛔ SATU-SATUNYA fungsi yang menyentuh view itu. Penjaga master mengunci
// namanya; pembaca kedua akan berbunyi.
func (r *ProdukLife) Ambang(ctx context.Context, produkID string) (AmbangProduk, error) {
	if strings.TrimSpace(produkID) == "" {
		return AmbangProduk{}, fmt.Errorf("%w: id produk kosong",
			ErrAmbangProdukTakDitemukan)
	}
	// 11-10-2026: tabel DATAR Master Product Name Life lebih dulu (yang ditulis aplikasi Master hari ini); view warisan
	// hanya bila produknya tidak ada di tabel datar (belum dipindah alat pindahflat).
	datar, err := r.db.Qualify(TabelProdukNameLife)
	if err != nil {
		return AmbangProduk{}, err
	}
	qd := sqlAmbangProdukDatar(datar)
	if err := db.PeriksaSQL(qd); err != nil {
		return AmbangProduk{}, err
	}
	var dKlaim, dTerima sql.NullString
	err = r.db.QueryRowContext(ctx, qd, produkID).Scan(&dKlaim, &dTerima)
	if err == nil {
		return AmbangProduk{
			MaxExpiredClaim: strings.TrimSpace(dKlaim.String),
			MaxDataReceive:  strings.TrimSpace(dTerima.String),
		}, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return AmbangProduk{}, fmt.Errorf("repository: membaca ambang produk (tabel datar): %w", err)
	}
	produk, err := r.db.Qualify(NamaViewProdukLife)
	if err != nil {
		return AmbangProduk{}, err
	}
	q := sqlAmbangProduk(produk)
	if err := db.PeriksaSQL(q); err != nil {
		return AmbangProduk{}, err
	}
	var maksKlaim, maksTerima sql.NullString
	if err := r.db.QueryRowContext(ctx, q, produkID).Scan(
		&maksKlaim, &maksTerima); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return AmbangProduk{}, fmt.Errorf("%w: %q",
				ErrAmbangProdukTakDitemukan, produkID)
		}
		return AmbangProduk{}, fmt.Errorf("repository: membaca ambang produk: %w", err)
	}
	return AmbangProduk{
		MaxExpiredClaim: strings.TrimSpace(maksKlaim.String),
		MaxDataReceive:  strings.TrimSpace(maksTerima.String),
	}, nil
}
