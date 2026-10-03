package repository

// Pencarian alamat risiko - popup Choose Risk Address (tiket 36). `[terverifikasi]`
// `NB FacIn\ReportDefinition\BrowseRisksAddress_RD.xml` (kelas ASM-FW-GISFW-Int-RISKADDRESS):
//   - INNER JOIN kelas ASM-FW-GISFW-Int-RW (prefix RW) pada `.PostalCode = RW.ZipCode`
//     -> tabel POOLDATA.RW kolom ZIPCODE (`D:\migrasi\RNM\DDL\RW.txt`; kelas->tabel `[dugaan]` -
//     RDB-List kelas yang sama `RDBList\BrowseRW_SQL.xml` membaca `FROM pooldata.rw`); hanya alamat
//     yang kode posnya ada di RW yang tampil;
//   - pyGetDistinctRows true -> SELECT DISTINCT;
//   - logika `A AND B AND C AND D AND E AND F AND G`: A `.Address Contains Param.JALAN`
//     (pyCaseInsensitive true), B `.PostalCode Contains KODEPOS`, C `.NationName Contains COUNTRY`,
//     D `.ProvinceName Contains PROVINCE`, E `.CityName Contains CITY`, F `.DistrictName Contains
//     DISTRICT`, G `.TerritoryName Contains TERRITORY` (B-G tanpa pyCaseInsensitive;
//     SearchRiskAddressAct meng-UPPER nilai A dan C-G, Zip Code / B TIDAK - baris 600);
//   - tanpa urutan; pyMaxRecords 100.
// Saringan kosong dilewati (Pega: parameter kosong tanpa "use null if empty" mengabaikan
// kondisinya `[dugaan]`). Kolom `[terverifikasi]` DDL `RISKADDRESS.txt` (16 kolom VARCHAR2(4000)).

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelRiskAddress, TabelRW - tabel warisan POOLDATA, baca saja.
	TabelRiskAddress = "RISKADDRESS"
	TabelRW          = "RW"
	// BatasRiskAddress - pyMaxRecords RD: paling banyak 100 baris hasil (A117).
	BatasRiskAddress = 100
)

// kolomSaringRisk - kolom RISKADDRESS tiap saringan (filter RD A..G).
func kolomSaringRisk(s models.SaringRisk) []struct{ kolom, nilai string } {
	return []struct{ kolom, nilai string }{
		{"a.ADDRESS", s.Address}, {"a.POSTALCODE", s.ZipCode}, {"a.NATIONNAME", s.Country}, {"a.PROVINCENAME", s.Province},
		{"a.CITYNAME", s.City}, {"a.DISTRICTNAME", s.District}, {"a.TERRITORYNAME", s.Territory},
	}
}

// kolomRisk - kolom hasil (RD) = medan models.RiskAddress, urut.
const kolomRisk = "a.ID, a.TITLE, a.ADDRESS, a.NATIONNAME, a.PROVINCENAME, a.CITYNAME, a.DISTRICTNAME, a.TERRITORYNAME, a.POSTALCODE"

// PembacaRisk - satu halaman alamat risiko + cacah (dibatasi BatasRiskAddress).
type PembacaRisk interface {
	CariRisk(ctx context.Context, s models.SaringRisk, offset, ukuran int) ([]models.RiskAddress, int, error)
}

// RiskOracle - PembacaRisk atas Oracle.
type RiskOracle struct{ db *db.DB }

// NewRiskOracle merakit pembaca alamat risiko.
func NewRiskOracle(d *db.DB) *RiskOracle { return &RiskOracle{db: d} }

// dasarRisk - SELECT DISTINCT ... JOIN RW ... WHERE (saringan terisi), argumen terikat.
// Tiap saringan: UPPER(kolom) LIKE pola huruf besar ESCAPE (A116: tidak peka huruf di semua
// kolom; RD hanya A yang pyCaseInsensitive).
func dasarRisk(alamat, rw string, s models.SaringRisk) (string, []any) {
	var syarat []string
	var arg []any
	for _, k := range kolomSaringRisk(s) {
		pola := PolaCari(k.nilai)
		if pola == "" {
			continue
		}
		arg = append(arg, pola)
		syarat = append(syarat, fmt.Sprintf("UPPER(%s) LIKE :%d ESCAPE '\\'", k.kolom, len(arg)))
	}
	q := "SELECT DISTINCT " + kolomRisk + " FROM " + alamat + " a JOIN " + rw + " w ON w.ZIPCODE = a.POSTALCODE"
	if len(syarat) > 0 {
		q += " WHERE " + strings.Join(syarat, " AND ")
	}
	return q, arg
}

// sqlCariRisk - 100 baris pertama, lalu satu halaman. A118: RD tanpa urutan; urut SEMBILAN
// kolom hasil (ID tidak ber-PK di DDL, DISTINCT atas sembilan kolom) supaya deterministik.
func sqlCariRisk(dasar string, nArg int) string {
	luar := strings.ReplaceAll(kolomRisk, "a.", "")
	// Kueri luar ber-ORDER BY sendiri: Oracle tidak menjamin urutan inline view terbawa.
	return fmt.Sprintf("SELECT %s FROM (%s ORDER BY %s FETCH FIRST %d ROWS ONLY) ORDER BY %s OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY",
		luar, dasar, kolomRisk, BatasRiskAddress, luar, nArg+1, nArg+2)
}

// sqlCacahRisk - cacah baris berbeda, dibatasi 100.
func sqlCacahRisk(dasar string) string {
	return fmt.Sprintf("SELECT COUNT(*) FROM (%s FETCH FIRST %d ROWS ONLY)", dasar, BatasRiskAddress)
}

// CariRisk - lihat PembacaRisk.
func (r *RiskOracle) CariRisk(ctx context.Context, s models.SaringRisk, offset, ukuran int) ([]models.RiskAddress, int, error) {
	alamat, err := r.db.Qualify(TabelRiskAddress)
	if err != nil {
		return nil, 0, err
	}
	rw, err := r.db.Qualify(TabelRW)
	if err != nil {
		return nil, 0, err
	}
	dasar, arg := dasarRisk(alamat, rw, s)
	var total int
	if err := r.db.QueryRowContext(ctx, sqlCacahRisk(dasar), arg...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("repository: cacah %s: %w", TabelRiskAddress, err)
	}
	baris, err := r.db.QueryContext(ctx, sqlCariRisk(dasar, len(arg)), append(arg, offset, ukuran)...)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: baca %s: %w", TabelRiskAddress, err)
	}
	defer baris.Close()
	hasil := []models.RiskAddress{}
	for baris.Next() {
		var id, judul, alamatBaris, negara, provinsi, kota, kecamatan, wilayah, kodePos sql.NullString
		if err := baris.Scan(&id, &judul, &alamatBaris, &negara, &provinsi, &kota, &kecamatan, &wilayah, &kodePos); err != nil {
			return nil, 0, fmt.Errorf("repository: %s: %w", TabelRiskAddress, err)
		}
		hasil = append(hasil, models.RiskAddress{ID: id.String, Title: judul.String, Address: alamatBaris.String,
			NationName: negara.String, ProvinceName: provinsi.String, CityName: kota.String, DistrictName: kecamatan.String,
			TerritoryName: wilayah.String, PostalCode: kodePos.String})
	}
	if err := baris.Err(); err != nil {
		return nil, 0, fmt.Errorf("repository: %s: %w", TabelRiskAddress, err)
	}
	return hasil, total, nil
}
