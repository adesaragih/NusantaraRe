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
	nama, err := g.db.Qualify(TabelJenisReasuransiWarisan)
	if err != nil {
		return nil, err
	}
	// ⭐ Isi menu Reinsurance Type (`modul/reinsurancetype`, tabel yang sama):
	// Name = `NOTE`, SOA Name = `SOANOTE`. Saringan dan urutan RD Pega.
	q := fmt.Sprintf(`SELECT ID, NOTE, SOANOTE FROM %s WHERE FLAG = 'active'
		ORDER BY ID DESC FETCH FIRST 500 ROWS ONLY`, nama)
	rows, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca pilihan Treaty Type dari %s: %w", TabelJenisReasuransiWarisan, err)
	}
	defer func() { _ = rows.Close() }()
	keluar := []models.PilihanWarisan{}
	for rows.Next() {
		var id, nm, soa sql.NullString
		if err := rows.Scan(&id, &nm, &soa); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s: %w", TabelJenisReasuransiWarisan, err)
		}
		keluar = append(keluar, models.PilihanWarisan{ID: id.String, Nama: nm.String, NamaSOA: soa.String})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	tandaiKembar(keluar)
	return keluar, nil
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
	if err := baris.Err(); err != nil {
		return nil, err
	}
	tandaiKembar(keluar)
	return keluar, nil
}

// tandaiKembar menandai pilihan yang NAMANYA dipakai lebih dari satu pengenal.
//
// ---------------------------------------------------------------------
// ⛔ MENGAPA INI PERLU: pencarian BALIK pengenal dari nama
// ---------------------------------------------------------------------
// Dokumen warisan menyimpan NAMA jenis treaty hampir selalu, tetapi
// PENGENALNYA hampir tidak pernah — terukur `TreatyTypeID` terisi pada 19
// dari 1.360 elemen `Limits[]`. Dropdown mengikat pengenal, jadi 98,6% layer
// berbunyi `Choose` padahal namanya ada.
//
// ⭐ Pega sendiri menyelesaikannya dengan mencari BALIK, dan rulenya ada di
// korpus — `Claim Non Prop/RDBList/GetReinsuranceTypeBYName_SQL.xml`:
//
//	select ID, NOTE as "Note" from REINSURANCETYPE
//	 where type = '4' and FLAG = 'active' and note = {InputSpreading.CARI1}
//
// Jadi mencari pengenal dari nama BUKAN karangan; ia pola yang sistem lama
// pakai. Yang TIDAK disalin saringan `type = '4'`: rule itu melayani
// spreading, sementara dropdown Limits (`BrowseReinsuranceType_RD`) menyaring
// `Flag = "active"` SAJA. Mencari di populasi yang lebih sempit daripada yang
// ditawarkan dropdown akan gagal menemukan jenis yang dropdown-nya sendiri
// tampilkan.
//
// ⚠️ DAN INILAH PAGARNYA: pencarian balik hanya aman bila namanya TUNGGAL.
// Dua baris bernama sama membuat pencarian memilih salah satu, dan layer akan
// menunjuk jenis treaty yang KELIRU tanpa ada yang tahu. Penanda ini yang
// membuat layar menolak menebak pada nama kembar, alih-alih menebak diam-diam.
func tandaiKembar(daftar []models.PilihanWarisan) {
	cacah := map[string]int{}
	for _, p := range daftar {
		cacah[p.Nama]++
	}
	for i := range daftar {
		daftar[i].Kembar = cacah[daftar[i].Nama] > 1
	}
}

// TabelKelasBisnisWarisan - tabel kelas `ASM-FW-GISFW-Int-TREATYBUSINESS`.
const TabelKelasBisnisWarisan = "TREATYBUSINESS"

// BacaDaftarKelasBisnisTreaty — autocomplete `Class of Business`
// (`BrowseTreatyBusinessWOType_RD`): `TreatyGroupID = pTreatyGroupId`,
// `TreatyGroupName` dan `BIZNAME` tidak kosong, `pyGetDistinctRows=true`,
// maks. 500, TANPA urutan. Nilai `.BizCode` → label `.BIZNAME`.
//
// ⚠️ DISTINCT atas pasangan (BizCode, BIZNAME) — tabelnya mengulang tiap
// pasangan per tahun treaty dan jenis reasuransi (4.489 baris untuk 26
// grup), dan RD-nya sendiri meminta baris berbeda.
func (g *Gudang) BacaDaftarKelasBisnisTreaty(ctx context.Context, treatyGroupID string) ([]models.PilihanWarisan, error) {
	nama, err := g.db.Qualify(TabelKelasBisnisWarisan)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT DISTINCT BIZCODE, BIZNAME FROM %s
		WHERE TREATYGROUPID = :1 AND TREATYGROUPNAME IS NOT NULL AND BIZNAME IS NOT NULL
		FETCH FIRST 500 ROWS ONLY`, nama)
	baris, err := g.db.QueryContext(ctx, q, treatyGroupID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca pilihan Class of Business dari %s: %w", TabelKelasBisnisWarisan, err)
	}
	defer func() { _ = baris.Close() }()
	keluar := []models.PilihanWarisan{}
	for baris.Next() {
		var id, nm sql.NullString
		if err := baris.Scan(&id, &nm); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s: %w", TabelKelasBisnisWarisan, err)
		}
		keluar = append(keluar, models.PilihanWarisan{ID: id.String, Nama: nm.String})
	}
	if err := baris.Err(); err != nil {
		return nil, err
	}
	tandaiKembar(keluar)
	return keluar, nil
}
