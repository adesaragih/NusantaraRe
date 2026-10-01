// Package repository adalah satu-satunya lapisan modul Endorsement Life yang
// menulis SQL.
//
// ⛔ Setiap nama tabel lewat `Qualify` (ADR-U-0033); nol `COMMIT` di teks SQL
// (ADR-U-0029) - transaksi dibuka services. Kode mengimpor `inti/...` dan
// modul ini saja: tabel PremiumList dan tabel warisan disebut namanya sendiri
// di sini, tidak dipinjam dari paket modul lain.
//
// ⛔ E4 (brief gelombang 2): `M_LIFE_PREMIUM_DETAIL` ±66,8 juta baris. Setiap
// SQL ke tabel itu berkunci kolom ber-index dan MENYEBUT indeksnya di komentar
// di atasnya; `TestKueriWarisanBerindex` menagihnya.
package repository

import (
	"errors"

	"nusantarare/inti/backend/db"
)

// Tabel aplikasi milik PremiumList Life (migrasi 051-056) yang dipakai
// bersama: endorsement adalah VERSI baru di tabel yang sama (spec §16).
const (
	tabelPolis          = "T_PREMIUM_LIST"
	tabelPeserta        = "T_PREMIUM_LIST_DETAIL"
	tabelSpreading      = "T_PREMIUM_LIST_SPREADING"
	tabelSpreadingRetro = "T_PREMIUM_LIST_SPREADING_RETRO"
	tabelRekap          = "T_PREMIUM_LIST_SUMMARY"
	tabelRiwayat        = "T_VIEW_SUGGEST"
	urutanKasus         = "SEQ_WORK_EDM_LIFE" // migrasi 481
)

// Tabel warisan `POOLDATA` - dibaca dan (E2) ditulis, tidak pernah dibuat.
const (
	// tabelPolisWarisan - rekam polis Pega; DIBACA saja (RALAT R18).
	tabelPolisWarisan = "JSON_POLIS"
	// tabelPesertaWarisanEDM - ⛔ sengaja BUKAN identifier `namaTabelPeserta`:
	// penjaga Claim Life memeriksa berkas yang memuat nama itu dengan aturan
	// kueri BACA berbatas hasil; penyalinan himpunan di sini berkunci
	// `PL_NUMBER` tetapi tidak berbatas baris (seluruh peserta satu polis).
	// Penjaganya sendiri: `TestKueriWarisanBerindex`.
	tabelPesertaWarisanEDM = "M_LIFE_PREMIUM_DETAIL"
	urutanPesertaWarisan   = "M_LIFE_PREMIUM_DETAIL_SEQ" // `SaveMasterLPDet` VALUES butir 1
	tabelRekapWarisan      = "M_LIFE_PREMIUM_SUMMARY"
	urutanRekapWarisan     = "M_LIFE_PREMIUM_SUMMARY_SEQ"
)

// tabelArasapas - gerbang ke-5, lintas skema (spec §13, RALAT R26): satu-satunya
// nama tabel berskema tetap di modul ini, hanya di `edm_arasapas.go`.
const tabelArasapas = "ARASAPAS.DETAIL_INVOICE"

// ErrTidakAda - baris yang dicari tidak ada.
var ErrTidakAda = errors.New("repository: baris tidak ada")

// Gudang adalah implementasi Oracle penyimpanan modul ini.
type Gudang struct{ db *db.DB }

// Baru menyusun Gudang di atas koneksi bersama.
func Baru(d *db.DB) *Gudang { return &Gudang{db: d} }

// nama mengkualifikasi beberapa tabel sekaligus, urut masukan.
func (g *Gudang) nama(objek ...string) ([]string, error) {
	hasil := make([]string, len(objek))
	for i, o := range objek {
		q, err := g.db.Qualify(o)
		if err != nil {
			return nil, err
		}
		hasil[i] = q
	}
	return hasil, nil
}
