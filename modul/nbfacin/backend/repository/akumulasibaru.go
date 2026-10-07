package repository

// Form Add New akumulasi popup Choose Accumulation (tiket 46, permintaan work owner 04-10-2026 "tombol add
// accumulationnya mana?") - dua pembaca untuk form `NB FacIn\Section\InputAccumulationCov.xml`. Simpannya lewat
// mesin master inti (master "accumulation"), bukan di sini.
//
//   - CZone dari zip - A181 (diganti 04-10-2026, work owner: "pada saat add new, ada kolom Cresta Zone, itu otomatis
//     keisi dari tabel rw, diambil dari zcone zipcode yang dipilih"): cZone = RW.CZONE ber-ZIPCODE zip itu (MIN sesudah
//     NULL / spasi dibuang - pasti bila nilainya beberapa; HANYA RW AKTIF `STS_AKTIF = '1'`, work owner: "yang aktif saja,
//     dan di dalam master datanya tidak ada yg czone nya kosong"), cZoneId = MIN CZONE.ID ber-CODE
//     itu bila ada (tanpa syarat aktif - nilainya datang dari RW). MENYIMPANG dari `GetCzone_Act` +
//     `RDBList\GetAccumulationZipcode_SQL.xml` `[terverifikasi]` (`select id, code from czone where code in (select distinct czone
//     from rw where zipcode = ...)`, baris 1) yang KOSONG bila CZONE RW tidak ada di tabel CZONE.
//   - Zip Code `[terverifikasi]` `ReportDefinition\BrowseRiskAddressZipCode_RD.xml` (kelas Int-RISKADDRESS = tabel
//     RISKADDRESS, `risiko.go`): `.ProvinceName Contains Param.ProvinceName`, pyMaxRecords 100, tanpa urutan, tanpa
//     DISTINCT; nationInitial = `Activity\SetCountryID_Act.xml` (BrowseNation_RD `ID = CARI33 OR Note = CARI38`,
//     pxResults(1).NationInitial; CARI33 / CARI38 = IDNation / NationName baris terpilih `[dugaan kuat]`
//     `Section\ChooseZipCodeDtl.xml`). Keputusan A182: `q` = Contains atas zip (tambahan agent, pencarian layar); baris
//     dibuat DISTINCT atas lima kolom yang tampil (RISKADDRESS = satu baris per ALAMAT - tanpa DISTINCT 100 baris dapat
//     berisi satu zip berulang); urut zip, kota, provinsi, negara; NationInitial = MAX (pxResults(1) tanpa urutan).
//     A185 (permintaan work owner 04-10-2026 "ada page nya, 15 list perpage"): halaman di server, total = COUNT atas
//     kueri yang sama; batas pyMaxRecords 100 DICABUT (Pega DEV menampilkan "Page 1 of 2421" - RD dipakai berhalaman).
//     A187 (work owner 04-10-2026 "yang aktif saja, dan di dalam master datanya tidak ada yg czone nya kosong"): HANYA zip
//     yang punya baris RW AKTIF ber-CZONE terisi (`RW.ZIPCODE = RISKADDRESS.POSTALCODE`, DDL keduanya) - zip terpilih
//     selalu ber-Cresta Zone. MENYIMPANG dari RD (yang tidak menyaring RW).

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// PembacaAddAkumulasi - pembaca form Add New akumulasi.
type PembacaAddAkumulasi interface {
	// CZoneZip - CZone pertama zip itu; ada=false bila tidak ada.
	CZoneZip(ctx context.Context, zip string) (models.CZoneZip, bool, error)
	// CariZip - satu halaman saran Zip Code (nomor >= 1) dan jumlah seluruhnya; provinsi / kata kosong = tanpa saringan.
	CariZip(ctx context.Context, provinsi, kata string, nomor, ukuran int) ([]models.ZipAkumulasi, int, error)
}

// AddAkumulasiOracle - PembacaAddAkumulasi atas Oracle.
type AddAkumulasiOracle struct{ db *db.DB }

// NewAddAkumulasiOracle merakit pembaca form Add New akumulasi.
func NewAddAkumulasiOracle(d *db.DB) *AddAkumulasiOracle { return &AddAkumulasiOracle{db: d} }

// syaratRWBerCZone - RW aktif ber-CZONE terisi (A181 / A187).
const syaratRWBerCZone = "r.STS_AKTIF = '" + StatusRWAktif + "' AND TRIM(r.CZONE) IS NOT NULL"

// sqlCZoneZip - A181: CZONE RW AKTIF zip itu (MIN, tanpa NULL / spasi) dan ID CZONE ber-CODE itu (MIN; NULL bila tidak ada).
// Agregat tanpa GROUP BY selalu satu baris - `WHERE w.CZ IS NOT NULL` membuatnya nol baris bila zip tanpa CZONE di RW.
func sqlCZoneZip(czone, rw string) string {
	return "SELECT w.CZ, (SELECT MIN(z.ID) FROM " + czone + " z WHERE z.CODE = w.CZ) FROM (SELECT MIN(r.CZONE) AS CZ FROM " +
		rw + " r WHERE r.ZIPCODE = :1 AND " + syaratRWBerCZone + ") w WHERE w.CZ IS NOT NULL"
}

// sqlZipAkumulasi - BrowseRiskAddressZipCode_RD (lihat kepala berkas), tanpa urutan / halaman: dasar bersama hitungan
// dan halaman. Saringan dan bind dibangun BERSAMA.
func sqlZipAkumulasi(alamat, nation, rw, polaProvinsi, polaZip string) (string, []any) {
	syarat := []string{"EXISTS (SELECT 1 FROM " + rw + " r WHERE r.ZIPCODE = a.POSTALCODE AND " + syaratRWBerCZone + ")"}
	var arg []any
	bind := func(v any) string { arg = append(arg, v); return ":" + strconv.Itoa(len(arg)) }
	if polaProvinsi != "" {
		syarat = append(syarat, "UPPER(a.PROVINCENAME) LIKE "+bind(polaProvinsi)+` ESCAPE '\'`)
	}
	if polaZip != "" {
		syarat = append(syarat, "UPPER(a.POSTALCODE) LIKE "+bind(polaZip)+` ESCAPE '\'`)
	}
	where := " WHERE " + strings.Join(syarat, " AND ")
	return "SELECT DISTINCT a.POSTALCODE, a.CITYNAME, a.PROVINCENAME, a.NATIONNAME, (SELECT MAX(n.NATIONINITIAL) FROM " +
		nation + " n WHERE n.ID = a.IDNATION OR n.NOTE = a.NATIONNAME) FROM " + alamat + " a" + where, arg
}

// sqlHitungZip - jumlah baris kueri dasar.
func sqlHitungZip(dasar string) string { return "SELECT COUNT(*) FROM (" + dasar + ")" }

// sqlHalamanZip - satu halaman kueri dasar: urut zip, kota, provinsi, negara, inisial (kelima kolom - halaman stabil),
// OFFSET / FETCH = dua bind sesudah `nBind`.
func sqlHalamanZip(dasar string, nBind int) string {
	return dasar + " ORDER BY 1, 2, 3, 4, 5 OFFSET :" + strconv.Itoa(nBind+1) + " ROWS FETCH NEXT :" + strconv.Itoa(nBind+2) + " ROWS ONLY"
}

// CZoneZip - lihat PembacaAddAkumulasi.
func (r *AddAkumulasiOracle) CZoneZip(ctx context.Context, zip string) (models.CZoneZip, bool, error) {
	czone, err := r.db.Qualify(TabelCZone)
	if err != nil {
		return models.CZoneZip{}, false, err
	}
	rw, err := r.db.Qualify(TabelRW)
	if err != nil {
		return models.CZoneZip{}, false, err
	}
	var id, kode sql.NullString
	err = r.db.QueryRowContext(ctx, sqlCZoneZip(czone, rw), zip).Scan(&kode, &id)
	if err == sql.ErrNoRows {
		return models.CZoneZip{}, false, nil
	}
	if err != nil {
		return models.CZoneZip{}, false, fmt.Errorf("repository: CZone zip: %w", err)
	}
	return models.CZoneZip{ID: id.String, Code: kode.String}, true, nil
}

// CariZip - lihat PembacaAddAkumulasi.
func (r *AddAkumulasiOracle) CariZip(ctx context.Context, provinsi, kata string, nomor, ukuran int) ([]models.ZipAkumulasi, int, error) {
	alamat, err := r.db.Qualify(TabelRiskAddress)
	if err != nil {
		return nil, 0, err
	}
	nation, err := r.db.Qualify(TabelNation)
	if err != nil {
		return nil, 0, err
	}
	rw, err := r.db.Qualify(TabelRW)
	if err != nil {
		return nil, 0, err
	}
	dasar, arg := sqlZipAkumulasi(alamat, nation, rw, PolaCari(provinsi), PolaCari(kata))
	var total int
	if err := r.db.QueryRowContext(ctx, sqlHitungZip(dasar), arg...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("repository: hitung %s: %w", TabelRiskAddress, err)
	}
	baris, err := r.db.QueryContext(ctx, sqlHalamanZip(dasar, len(arg)), append(arg, (nomor-1)*ukuran, ukuran)...)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: baca %s: %w", TabelRiskAddress, err)
	}
	defer baris.Close()
	hasil := []models.ZipAkumulasi{}
	for baris.Next() {
		var zip, kota, prov, negara, inisial sql.NullString
		if err := baris.Scan(&zip, &kota, &prov, &negara, &inisial); err != nil {
			return nil, 0, fmt.Errorf("repository: %s: %w", TabelRiskAddress, err)
		}
		hasil = append(hasil, models.ZipAkumulasi{ZipCode: zip.String, City: kota.String, Province: prov.String,
			Nation: negara.String, NationInitial: inisial.String})
	}
	return hasil, total, baris.Err()
}
