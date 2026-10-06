package repository

// Isi dropdown tab Limits (proporsional) — tiga Report Definition Pega.
//
// ⭐ Sumbernya dibaca dari ekspor, bukan ditebak:
//
//	`Treaty Type`  LimitProportional `.TreatyTypeID` → `BrowseReinsuranceType_RD`
//	               kelas `ASM-FW-GISFW-Int-REINSURANCETYPE`, nilai `.ID`, label
//	               `.Note`, `Flag="active"`, urut `.ID` DESC, maks. 500.
//	`Treaty Group` DetailLimits `.TreatyGroupID` → `BrowseTreatyGroup_RD`
//	               kelas `ASM-FW-GISFW-Int-TREATYGROUP`, nilai `ID`, label
//	               `.TreatyGroupName`, urut `.ID` DESC, maks. 500.
//	Mata uang      DetailLimits `.Currency` (grid 100% Limit, Retention,
//	               Cession, Reserve, PLA, Cash Loss, Claim Coop, EPI, Event
//	               Limits) → `BrowseCurrencyTreatyIn_RD`; `.CurrencyID` grid
//	               Deduction → `BrowseCurrency_RD`. Keduanya kelas
//	               `ASM-FW-GISFW-Int-CURRENCY`, `.Currency != "ITL"`, TANPA
//	               urutan, maks. 500.
//
// Kelas → tabel: SQL mentah kelas yang sama di `Claim Non Prop/RDBList/
// GetReinsuranceTypeBYName_SQL.xml` (`select ID, NOTE as "Note" from
// REINSURANCETYPE …`), dan kolom ketiga tabel persis medan RD-nya.
// Diukur 6 Oktober 2026: REINSURANCETYPE 130 baris (84 `active`),
// TREATYGROUP 33, CURRENCY 38 (satu `ITL`).
//
// ⛔ BACA SAJA — ketiganya dijaga `TestWarisanHanyaDibaca`.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/modul/treatyin/backend/models"
)

// Tabel kelas `ASM-FW-GISFW-Int-…` yang dibaca dropdown tab Limits.
const (
	TabelJenisReasuransiWarisan = "REINSURANCETYPE"
	TabelKelompokTreatyWarisan  = "TREATYGROUP"
	TabelMataUangWarisan        = "CURRENCY"
)

// BacaDaftarJenisTreaty — `BrowseReinsuranceType_RD`.
func (g *Gudang) BacaDaftarJenisTreaty(ctx context.Context) ([]models.PilihanWarisan, error) {
	return g.bacaPilihanRD(ctx, TabelJenisReasuransiWarisan, "ID", "NOTE", "FLAG = 'active'", "ID DESC")
}

// BacaDaftarKelompokTreaty — `BrowseTreatyGroup_RD`.
func (g *Gudang) BacaDaftarKelompokTreaty(ctx context.Context) ([]models.PilihanWarisan, error) {
	return g.bacaPilihanRD(ctx, TabelKelompokTreatyWarisan, "ID", "TREATYGROUPNAME", "", "ID DESC")
}

// BacaDaftarMataUangLimit — `BrowseCurrencyTreatyIn_RD`/`BrowseCurrency_RD`.
//
// ⚠️ RD-nya tidak menyatakan urutan; urutan Oracle-lah yang dilihat Pega,
// dan itu pula yang dikembalikan di sini.
func (g *Gudang) BacaDaftarMataUangLimit(ctx context.Context) ([]models.PilihanWarisan, error) {
	return g.bacaPilihanRD(ctx, TabelMataUangWarisan, "ID", "CURRENCY", "CURRENCY != 'ITL'", "")
}

// bacaPilihanRD menarik pasangan nilai+label satu RD, paling banyak 500
// (`pyMaxRecords`). Nama tabel, kolom, dan saringan DITANAM di berkas ini.
func (g *Gudang) bacaPilihanRD(ctx context.Context, tabel, kolomID, kolomNama, saring, urut string) ([]models.PilihanWarisan, error) {
	nama, err := g.db.Qualify(tabel)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf("SELECT %s, %s FROM %s", kolomID, kolomNama, nama)
	if saring != "" {
		q += " WHERE " + saring
	}
	if urut != "" {
		q += " ORDER BY " + urut
	}
	q += " FETCH FIRST 500 ROWS ONLY"
	baris, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca pilihan dari %s: %w", tabel, err)
	}
	defer func() { _ = baris.Close() }()

	keluar := []models.PilihanWarisan{}
	for baris.Next() {
		var id, nm sql.NullString
		if err := baris.Scan(&id, &nm); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s: %w", tabel, err)
		}
		keluar = append(keluar, models.PilihanWarisan{ID: id.String, Nama: nm.String})
	}
	return keluar, baris.Err()
}
