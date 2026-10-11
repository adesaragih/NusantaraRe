// Package repository adalah SATU-SATUNYA lapisan modul Master Contract Retro
// Life yang berbicara ke Oracle. Seluruh SQL modul ini ada di sini.
//
// Untuk apa berkas ini: nama objek Oracle dan daftar kolom modul ini, ditulis
// SEKALI. Lima tabel WARISAN `POOLDATA.*_LIFE` ditulis dan dibaca apa adanya
// (K1: nol DDL, nol tabel baru); master rujukan hanya dibaca.
//
// ⛔ Setiap query menyebut skemanya lewat `db.Qualify` (ADR-U-0033) dan
// diperiksa `db.PeriksaSQL` (nol COMMIT, ADR-U-0029).
// ⛔ Uang, persen, dan tanggal diminta sebagai TEKS lewat `db.FmtDesimal` /
// `db.FmtTanggalOracle`, lalu diurai sekali (ADR-U-0003).
//
// Dibaca sesudah: models/mcrl_entitas.go.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// Lima tabel warisan yang modul ini tulis dan baca (`ddl-tables-from-dba.md`).
const (
	TabelTahun     = "TREATYYEAR_LIFE"
	TabelKontrak   = "TREATYCONTRACT_LIFE"
	TabelReinsurer = "TREATYREINSURER_LIFE"
	TabelSecurity  = "TREATYSECURITYREINSURER_LIFE"
	TabelBusiness  = "TREATYBUSINESS_LIFE"
)

// Master rujukan - dibaca saja.
const (
	// MasterJenisReasuransi - kelas `ASM-FW-GISFW-Int-REINSURANCETYPE`
	// (tabel fisik yang sama dipakai Treaty Contract Out).
	MasterJenisReasuransi = "REINSURANCETYPE"
	// MasterReinsurer - kelas `ASM-FW-GISFW-Int-AGENT` (`BrowseCedingCoLife_RD`).
	MasterReinsurer = "AGENT"
	// MasterBusiness - kelas `ASM-FW-GISFW-Int-BUSINESS` (`BrowseBusinessLife_RD`).
	MasterBusiness = "BUSINESS"
	// MasterRingkasanRate - kelas `ASM-FW-GISFW-Int-RATE_LIFE_SUMMARY` (`BrowseRateLifeSummary` b40),
	// sumber autocomplete `R/I RATE`; tabel `M_RATE_LIFE_SUMMARY` (kolom ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID - keputusan work owner 07-10-2026).
	// MasterRate - kelas `ASM-FW-GISFW-Int-M_RATE_LIFE` (`BrowseRateLife_RD` b39), section `Rate List`;
	// dulu view DEV `RATE_LIFE` (8 kolom, atas `M_RATE_LIFE.JSONDATA`) - nama fisik kelas itu terbukti
	// `NB FacIn/RDBList/BrowseLifeRate_SQL.xml` b85 `… FROM RATE_LIFE WHERE IDUSEDBY= …`. RALAT 07-10-2026
	// (keputusan work owner, `modul/riratelife/MODUL.md` RALAT R7): kini TABEL flat `M_RATE_LIFE` berkolom sama
	// (view dibuang migrasi inti 930); tetap baca-saja.
	//
	// ⛔ K1 keputusan work owner 01-10-2026 (OQ-MCRL-13 + OQ-MCRL-05): kedua view dibaca SAJA - kolom
	// yang dibaca RD XML saja, nol `SELECT *`, nol `JSONDATA`, nol tulisan (`periksaBacaSaja`).
	MasterRingkasanRate = "M_RATE_LIFE_SUMMARY"
	MasterRate          = "M_RATE_LIFE"
)

// DaftarTabelWarisan - lima tabel yang ditulis modul ini (penjaga modul).
var DaftarTabelWarisan = []string{TabelTahun, TabelKontrak, TabelReinsurer, TabelSecurity, TabelBusiness}

// DaftarMasterDibacaSaja - objek yang dibaca tetapi tidak pernah ditulis.
var DaftarMasterDibacaSaja = []string{MasterJenisReasuransi, MasterReinsurer, MasterBusiness,
	MasterRingkasanRate, MasterRate}

// Kolom tiap tabel warisan, urutan DDL `[data DBA]`.
var (
	KolomTahun   = []string{"ID", "TREATYYEAR", "UNDERWRITINGYEAR", "USERID", "TGLUPDATE", "STARTDATE", "ENDDATE"}
	KolomKontrak = []string{"ID", "IDTREATYYEAR", "REINSTYPEID", "REINSTYPENAME", "USERID", "TGLUPDATE",
		"IDR", "USD", "B_IDR", "B_USD", "IDR_SELISIH", "USD_SELISIH", "TREATYENDDATE", "TREATYSTARTDATE"}
	KolomReinsurer = []string{"ID", "TREATYYEARID", "TREATYCONTRACTID", "REINSTYPEID", "REINSTYPENAME",
		"REINSURERNAME", "PCTSHARE", "COMMISION", "OVR_COMM", "USERID", "TGLUPDATE", "REINSURERID"}
	KolomSecurity = []string{"ID", "TREATYYEARID", "TREATYCONTRACTID", "TREATYREINSURERID", "REINSURERID",
		"REINSURERNAME", "PCTSHARE", "USERID", "TGLUPDATE"}
	KolomBusiness = []string{"ID", "TREATYYEARID", "TREATYYEAR", "REINSTYPEID", "REINSTYPENAME", "BIZCODE",
		"BIZNAME", "USERID", "TGLUPDATE", "TREATYCONTRACTID", "RIRATEID", "RIRATE"}
)

// kolomTanggal dan kolomDesimal - kolom bertipe DATE / NUMBER `[data DBA]`;
// sisanya VARCHAR2.
var (
	kolomTanggal = map[string]bool{"TGLUPDATE": true, "STARTDATE": true, "ENDDATE": true,
		"TREATYSTARTDATE": true, "TREATYENDDATE": true}
	kolomDesimal = map[string]bool{"IDR": true, "USD": true, "B_IDR": true, "B_USD": true, "IDR_SELISIH": true,
		"USD_SELISIH": true, "PCTSHARE": true, "COMMISION": true, "OVR_COMM": true}
)

// KolomTanggal / KolomDesimal menyatakan tipe fisik kolom untuk penjaga dan
// tiruan uji `db`.
func KolomTanggal(k string) bool { return kolomTanggal[k] }

// KolomDesimal - lihat KolomTanggal.
func KolomDesimal(k string) bool { return kolomDesimal[k] }

// daftarPilih merakit daftar SELECT: tanggal dan desimal diminta sebagai teks.
func daftarPilih(kolom []string) string {
	b := make([]string, len(kolom))
	for i, k := range kolom {
		switch {
		case kolomTanggal[k]:
			b[i] = fmt.Sprintf(db.FmtTanggalOracle, k)
		case kolomDesimal[k]:
			b[i] = fmt.Sprintf(db.FmtDesimal, k)
		default:
			b[i] = k
		}
	}
	return strings.Join(b, ", ")
}

// kuerier - baca lewat transaksi bila ada (anti-basi di dalam penyimpanan),
// selain itu lewat koneksi.
type kuerier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Gudang membaca dan menulis kelima tabel warisan serta membaca master.
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
	if err := periksaBacaSaja(objek, q); err != nil {
		return "", err
	}
	return q, nil
}

// ErrMasterBacaSaja - SQL yang bukan SELECT diarahkan ke objek master/view yang dibaca saja.
var ErrMasterBacaSaja = errors.New("repository: reference master is read-only")

// periksaBacaSaja - lapis kedua penjaga modul `TestMCRLMasterDibacaSaja`: objek di
// DaftarMasterDibacaSaja hanya menerima SELECT, sebelum teks SQL apa pun sampai ke Oracle.
func periksaBacaSaja(objek, q string) error {
	for _, m := range DaftarMasterDibacaSaja {
		if m == objek && !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(q)), "SELECT ") {
			return fmt.Errorf("%w: %s", ErrMasterBacaSaja, objek)
		}
	}
	return nil
}
