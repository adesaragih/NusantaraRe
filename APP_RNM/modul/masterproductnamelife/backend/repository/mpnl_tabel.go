// Package repository adalah SATU-SATUNYA lapisan modul Master Product Name
// Life yang berbicara ke Oracle, dan satu-satunya yang mengenal bentuk JSON
// Pega. Seluruh SQL dan seluruh kunci `JSONDATA` modul ini ada di sini.
//
// Untuk apa berkas ini: nama objek Oracle, kolom, dan kunci JSON yang dibaca
// ketiga view DEV - ditulis SEKALI.
//
// ⛔ P1 (`docs/RALAT-DEV-30-09-2026.md`): produk disimpan SEPERTI PEGA - dua
// tabel warisan `POOLDATA.M_PRODUCT_LIFE` dan `M_PRODUCTINWARD_LIFE`, kolom
// `JSONDATA` berkunci persis Pega (tiga view, dua prosedur, dan Claim Life
// membacanya). Nol tabel baru, nol DDL.
// ⛔ P2: modul ini PENULIS kedua tabel - boleh membaca dan menulis `JSONDATA`
// miliknya sendiri. Larangan AC 38 milik Claim Life sebagai pembaca.
// ⛔ Setiap query menyebut skemanya lewat `db.Qualify` dan diperiksa
// `db.PeriksaSQL` (nol COMMIT). Prosedur `PEGA_M_PRODUCT_LIFE` dan
// `PEGA_M_PRODUCT_INWARD_LIFE` TIDAK dipanggil (logikanya ditiru).
//
// Dibaca sesudah: models/mpnl_produk.go.
package repository

import (
	"context"
	"database/sql"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
)

// Dua tabel warisan produk (`dba-procedures-and-ddl.md` DDL) - ditulis dan dibaca.
const (
	TabelProduk = "M_PRODUCT_LIFE"
	TabelInward = "M_PRODUCTINWARD_LIFE"
)

// Kolom fisik kedua tabel - katalog DEV `ALL_TAB_COLUMNS` 01-10-2026
// (`testdata/katalog-dev.json`, uji `TestKolomDitulisAdaDiKatalogDEV`). `M_PRODUCT_LIFE`:
// `JSONDATA` CLOB (`IS JSON`) + DUA kolom datar `RIRISKID`, `RIRISK` - penulisnya di
// korpus hanya `SaveProductNameLIfeFlat` b84. ⛔ `PRODUCTNAME` dan `BEGIN_DATE` TIDAK ADA
// di DEV (ralat lanjutan 1 L1; `dba-procedures-and-ddl.md` keliru menyebut empat).
// `M_PRODUCTINWARD_LIFE`: `ID` + `JSONDATA`.
var (
	KolomProduk = []string{"ID", "JSONDATA", "RIRISKID", "RIRISK"}
	KolomInward = []string{"ID", "JSONDATA"}
)

// Lebar kolom datar - katalog DEV `ALL_TAB_COLUMNS` (`testdata/katalog-dev.json` `lebar`, uji
// `TestLebarKolomDatarSamaDenganKatalogDEV`). Nilai yang lebih panjang ditolak services dengan kalimat
// yang menyebut medannya, bukan dipotong Oracle. ⚠️ Dihitung dalam BYTE: semantik BYTE/CHAR kolom DEV
// tidak tercatat di katalog - byte adalah batas yang aman untuk keduanya.
const (
	LebarRIRiskID = 10
	LebarRIRisk   = 100
)

// LebarKunciView - kunci skalar yang dibaca view `PRODUCT_LIFE` / `PRODUCTINWARD_LIFE` keluar sebagai
// VARCHAR2(4000) (`claimlife/docs/KATALOG-TABEL-PESERTA-DAN-TREATY.md`); nilai yang lebih panjang
// dibaca view sebagai NULL (`NULL ON ERROR`) - Claim Life diam-diam kehilangan nilainya.
const LebarKunciView = 4000

// KunciViewTerlaluPanjang - kunci view skalar yang nilainya (bentuk Pega) melampaui LebarKunciView.
func KunciViewTerlaluPanjang(p models.Produk) []string {
	umum, inward := HalamanPega(p)
	var hasil []string
	for _, k := range KunciViewProduk {
		if len(umum[k]) > LebarKunciView {
			hasil = append(hasil, k)
		}
	}
	for _, k := range KunciViewInward {
		if len(inward[k]) > LebarKunciView {
			hasil = append(hasil, k)
		}
	}
	return hasil
}

// DaftarTabelWarisan - tabel yang ditulis modul ini (penjaga modul).
var DaftarTabelWarisan = []string{TabelProduk, TabelInward}

// Kunci JSON yang dibaca ketiga view DEV (`docs/dba-view-produk-life.md`) -
// peka huruf besar-kecil. Uji murni menuntut setiap kunci ini ada di JSON
// hasil simpan (brief bab 4); kunci larik ditulis lewat barisnya.
var (
	// KunciViewProduk - view `PRODUCT_LIFE` atas `M_PRODUCT_LIFE`.
	KunciViewProduk = []string{"TYPE", "TYPE_CEDING", "CEDING", "CEDINGID", "SOBNAME", "SOBID", "CAUSEID", "GRUP",
		"PRODUCTNAME", "PRODUCTCODE", "PRODUCTTYPEID", "PRODUCTTYPE", "RIRISKID", "RIRISK", "RIRATEID", "RIRATE",
		"RICOMMID", "RICOMM", "INWARDNAME", "UnderwritingLimitList", "OUTWARDNAMEID", "OUTWARDNAME", "OUTWARDRATEID",
		"OUTWARDRATE", "OUTWARDCOMMID", "OUTWARDCOMM", "BENEFITID", "BENEFIT", "CAUSE", "OutwardList", "POLICYHODER",
		"POLICYHODERNAME", "TREATYNUMBER", "CREATEOP", "UPDATEOP"}
	// KunciViewProdukBarisOutward - `a.JSONDATA.OutwardList[0].OVR_COMM`.
	KunciViewProdukBarisOutward = []string{"OVR_COMM"}
	// KunciViewDokumen - view `DOCUMENTCLAIM_LIFE` (`$.DocumentClaim[*].Document`).
	KunciViewDokumen      = []string{"DocumentClaim"}
	KunciViewDokumenBaris = []string{"Document"}
	// KunciViewInward - view `PRODUCTINWARD_LIFE` atas `M_PRODUCTINWARD_LIFE`.
	KunciViewInward = []string{"PRODUCTID", "INSURED", "CEDING", "TREATYNUMBER", "ADDENDUMWORD", "AMANDEMENTSCHD",
		"INWARDTREATYNM", "BEGIN", "MATURE", "CEDINGRETENTIONNUM", "CEDINGRETENTIONPCT", "CEDINGLIMIT", "CEDINGLIMITXPN",
		"MINAGE", "MAXAGE", "BIRTHDAY", "EXTRAPREMI", "CURRENCY", "RNMSHARE", "EXTRAMORTALITY", "RNMLIMITNUM",
		"RNMLIMITPCT", "LIENCLAUSE", "MONTHS", "MINSUMINSURED", "MAXSUMINSURED", "MAXSUMREASURED", "MAXCONTRACT",
		"PAYMENT", "PROPORTIONALTABLE", "SUBJECTTO", "POLICYHODER", "POLICYHODERNAME", "BROKERAGE", "ADDENDUMNO",
		"AMANDEMENTNO", "MAXDATARECEIVE", "MAXEXPIREDCLAIM", "STNC"}
)

// kuerier - baca lewat transaksi bila ada (anti-basi di dalam penyimpanan),
// selain itu lewat koneksi.
type kuerier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Gudang membaca dan menulis kedua tabel produk (dan, paket 2+, membaca master).
type Gudang struct{ db *db.DB }

// Baru menyusun Gudang di atas koneksi bersama.
func Baru(d *db.DB) *Gudang { return &Gudang{db: d} }

func (g *Gudang) kueri(tx *db.Tx) kuerier {
	if tx.Terisi() {
		return tx
	}
	return g.db
}

// siapkan mengualifikasi nama objek dan memeriksa teks SQL-nya.
func (g *Gudang) siapkan(objek string, susun func(tabel string) string) (string, error) {
	tabel, err := g.db.Qualify(objek)
	if err != nil {
		return "", err
	}
	q := susun(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	return q, nil
}
