package repository

// Master jenis reasuransi - dibaca saja (tiket 02 Treaty Contract Out).
//
// `[terverifikasi]` Tabel fisik `POOLDATA.REINSURANCETYPE` - kolom `ID`,
// `NOTE`, `TYPE`, `FLAG` disebut SQL korpus (`Claim Non Prop/RDBList/
// GetReinsuranceTypeBYName_SQL.xml`: `select ID, NOTE … from REINSURANCETYPE
// where type ='4' and FLAG ='active'`; `EDM Treaty In/RDBList/
// GetReinstypeIDbyName_SQL.xml`). Master ini BERSAMA life dan non-life
// (spec §11); saringanlah yang memisahkan.
//
// Saringan `[terverifikasi]` `Treaty Contract Out/ReportDefinition/
// BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull.xml` - RD yang dipakai 11 grid
// klausul, ruleset 01-01-91 (2026-02-05):
//
//	b557  <pyFilterLogic>A AND B AND C</pyFilterLogic>
//	b567  A  .ID    <pyFilterOperation>NotStartsWith</pyFilterOperation> b581
//	      "10004","10011","10012","10021","10022","10025","10026","10028",
//	      "10248","10249","10018","10217"                              b573
//	b586  C  .Flag  "active"  (tanpa pyFilterOperation = Equal)
//	b603  B  .Type  "1","2","3" (tanpa pyFilterOperation = Equal atas daftar)
//	b667  sort .Note ASC; b742 pyMaxRecords 500
//
// ⛔ RALAT ATAS SPEC §11 / TIKET 02 (28-09-2026): keduanya menulis `.ID NOT
// IN (…)`. Operatornya `NotStartsWith` - "tidak BERAWALAN" salah satu dari dua
// belas nilai. Untuk ID lima karakter keduanya sama; untuk ID yang lebih
// panjang tidak. XML menang: `NOT LIKE '10004%'`.
//
// `[keputusan work owner]` Q5: dua belas ID ditiru persis, jangan
// digeneralkan. Penyimpangan sadar 8: nol kata "Old" di nama mana pun di sini
// - nama rule sumbernya disebut hanya sebagai bukti.
//
// ⛔ BACA SAJA (AC 12): nol INSERT/UPDATE/DELETE - dijaga
// TestTCOWarisanHanyaDibaca (daftar `masterDibacaSajaTCO`).

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/db"
)

// MasterJenisReasuransiTCO adalah nama tabel master jenis reasuransi.
const MasterJenisReasuransiTCO = "REINSURANCETYPE"

// masterDibacaSajaTCO adalah master yang modul ini BACA tetapi tidak miliki
// (spec b107). Penjaga TestTCOWarisanHanyaDibaca menolak fungsi yang menyebut
// namanya bersama kata kerja tulis.
var masterDibacaSajaTCO = []string{
	MasterJenisReasuransiTCO, "TREATYDESC", "TREATYGROUP", "TREATYEXCHANGEYEARLY",
	"CATEGORY_ATTACH_REAS",
	// Tiket 05: master reinsurer (pemilih `BrowseAgentReinsSOA_RD`).
	"AGENT",
	// Tiket 07: master bisnis (pemilih `BrowseFilterBusiness_RD`).
	"BUSINESS",
	// Tiket 11: master mata uang (`GetCurrencyID`) - pengenal USD untuk kurs.
	"CURRENCY",
	// OQ-TCO-08: APPNAME penyimpanan (`GetAppName_SQL`).
	"T_FOLDER_IMAGE",
	// Tiket 08: pemilih ExclutionTreaty (`BrowseOccupationFIRE_RD`,
	// `BrowseFireClauseFacIn_RD`). `TREATYDESC` sudah di atas.
	"OCCUPATION", "CLAUSE",
}

// BlacklistJenisReasuransiNonLife adalah dua belas awalan ID yang disingkirkan
// - VERBATIM b573, urutan RD dipertahankan.
var BlacklistJenisReasuransiNonLife = []string{
	"10004", "10011", "10012", "10021", "10022", "10025",
	"10026", "10028", "10248", "10249", "10018", "10217",
}

// FlagJenisReasuransiAktif - `.Flag = "active"` b586.
//
// ⚠️ Beda ejaan antar konteks: di Master Contract Retro Life `.Flag` bernilai
// `1`; di sini `"active"`. Jangan menyalin dari konteks Life.
const FlagJenisReasuransiAktif = "active"

// TipeJenisReasuransiNonLife - `.Type` di antara "1","2","3" b603.
var TipeJenisReasuransiNonLife = []string{"1", "2", "3"}

// JenisReasuransiTCO adalah satu baris master seperti dibaca modul ini.
//
// Hanya tiga kolom yang dibaca: ID (kunci), NOTE (nama tampil `.Note`),
// TYPE. Kolom `.SOANote` dan `.Code` yang RD pilih tidak terbukti nama
// fisiknya di SQL korpus mana pun dan tidak dipakai layar - `[terbuka]`.
type JenisReasuransiTCO struct {
	ID   string
	Note string
	Tipe string
}

// LolosSaringanNonLifeTCO adalah tabel kebenaran saringan RD, di Go.
//
// Dipakai uji sebagai acuan SQL di bawah, dan dapat dipakai layar mana pun
// yang perlu memeriksa satu nilai. Keduanya WAJIB sepakat.
func LolosSaringanNonLifeTCO(id, flag, tipe string) bool {
	if flag != FlagJenisReasuransiAktif {
		return false
	}
	tipeOK := false
	for _, t := range TipeJenisReasuransiNonLife {
		if tipe == t {
			tipeOK = true
		}
	}
	if !tipeOK {
		return false
	}
	for _, awalan := range BlacklistJenisReasuransiNonLife {
		if strings.HasPrefix(id, awalan) {
			return false
		}
	}
	return true
}

// sqlJenisReasuransiNonLifeTCO merakit pembacaan tersaring.
//
// Bind: :1 flag; :2-:4 tipe; :5-:16 awalan blacklist berakhiran '%'.
// Urutan `NOTE ASC` mengikuti b667; `ID` sebagai pemecah seri supaya
// hasilnya dapat diulang.
func sqlJenisReasuransiNonLifeTCO(tabel string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "SELECT ID, NOTE, TYPE FROM %s\n WHERE FLAG = :1\n   AND TYPE IN (:2, :3, :4)", tabel)
	for i := range BlacklistJenisReasuransiNonLife {
		fmt.Fprintf(&b, "\n   AND ID NOT LIKE :%d", 5+i)
	}
	b.WriteString("\n ORDER BY NOTE ASC, ID ASC")
	return b.String()
}

// argJenisReasuransiNonLifeTCO menyusun argumen bind dalam urutan yang sama.
func argJenisReasuransiNonLifeTCO() []any {
	arg := []any{FlagJenisReasuransiAktif}
	for _, t := range TipeJenisReasuransiNonLife {
		arg = append(arg, t)
	}
	for _, awalan := range BlacklistJenisReasuransiNonLife {
		arg = append(arg, awalan+"%")
	}
	return arg
}

// MasterJenisReasuransi membaca master jenis reasuransi.
type MasterJenisReasuransi struct{ db *db.DB }

// NewMasterJenisReasuransi menyusunnya.
func NewMasterJenisReasuransi(db *db.DB) *MasterJenisReasuransi {
	return &MasterJenisReasuransi{db: db}
}

// DaftarNonLife membaca jenis reasuransi yang lolos saringan non-life.
//
// Daftar kosong dikembalikan APA ADANYA; yang memutuskan bahwa kosong adalah
// kegagalan (ADR-0015) adalah lapisan services.
func (m *MasterJenisReasuransi) DaftarNonLife(ctx context.Context) ([]JenisReasuransiTCO, error) {
	tabel, err := m.db.Qualify(MasterJenisReasuransiTCO)
	if err != nil {
		return nil, err
	}
	q := sqlJenisReasuransiNonLifeTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := bacaTCO(ctx, m.db).QueryContext(ctx, q, argJenisReasuransiNonLifeTCO()...)
	if err != nil {
		return nil, fmt.Errorf("repository: reading reinsurance type master: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []JenisReasuransiTCO
	for rows.Next() {
		var id, note, tipe sql.NullString
		if err := rows.Scan(&id, &note, &tipe); err != nil {
			return nil, fmt.Errorf("repository: reading reinsurance type master: %w", err)
		}
		out = append(out, JenisReasuransiTCO{ID: id.String, Note: note.String, Tipe: tipe.String})
	}
	return out, rows.Err()
}
