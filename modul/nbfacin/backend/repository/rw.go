package repository

// Saran Zip Code popup Add alamat risiko (tiket 37) dan simpan alamat baru ke RISKADDRESS.
//
// Saran `[terverifikasi]`: `NB FacIn\Section\InputRiskAddress.xml` sel 13 = pxAutoComplete atas RD
// `BrowseRW_RD` (kelas ASM-FW-GISFW-Int-RW = `pooldata.rw`, `RDBList\BrowseRW_SQL.xml`), medan cari
// `.ZipCode`; memilih baris mengisi ZipCode, NATIONNAME, PROVINCENAME, CITYNAME, DistrictName, Note
// (-> TerritoryName). RD: `.STS_AKTIF = "1"`, DISTINCT, pyMaxRecords 10000. Kolom DDL `RW.txt`:
// ZIPCODE, NOTE, DISTRICTNAME, CITYNAME, PROVINCENAME, NATION - TIDAK ada NATIONNAME; alias
// `NATION as "NATIONNAME"` `[terverifikasi]` di `RDBList\BrowseRW2_SQL.xml` (kelas RW yang sama).
//
// Simpan = PORT Go dari prosedur `pooldata.InsertUpdateRISKADDRESS` (DDL INSERTUPDATERISKADDRESS.txt),
// keputusan work owner butir 81 (ADR-0043: nol CALL prosedur): cabang INSERT saja (alamat baru,
// Pega P_ID "UNKNOWNID"); ID = `getcurrentsite || lpad(to_char(RISKADDRESS_SEQ.nextval),12,'0')`
// dihitung dengan ekspresi yang SAMA (isi fungsi getcurrentsite tidak ada di korpus - dipanggil apa
// adanya, bukan dikarang). Transaksi milik pemanggil; tanpa COMMIT/ROLLBACK di SQL.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// BatasSaranRW - A124: paling banyak 50 saran (RD 10000).
	BatasSaranRW = 50
	// StatusRWAktif - filter RD F `.STS_AKTIF = "1"`.
	StatusRWAktif = "1"
	// fungsiSitus, sequenceRiskAddress - objek POOLDATA yang dipakai prosedur.
	fungsiSitus         = "GETCURRENTSITE"
	sequenceRiskAddress = "RISKADDRESS_SEQ"
)

// PenyimpanRW - saran Zip Code (RW) dan simpan alamat baru (RISKADDRESS), tiket 37.
type PenyimpanRW interface {
	PembacaRW
	PenulisRiskAddress
}

// PembacaRW - saran Zip Code.
type PembacaRW interface {
	CariRW(ctx context.Context, awalanZip string) ([]models.BarisRW, error)
}

// PenulisRiskAddress - simpan alamat baru di transaksi pemanggil; mengembalikan ID baru.
type PenulisRiskAddress interface {
	SisipRiskAddress(ctx context.Context, tx *db.Tx, a models.AlamatBaru) (string, error)
}

// RWOracle - PembacaRW + PenulisRiskAddress atas Oracle.
type RWOracle struct{ db *db.DB }

// NewRWOracle merakit pembaca RW / penulis RISKADDRESS.
func NewRWOracle(d *db.DB) *RWOracle { return &RWOracle{db: d} }

// sqlCariRW - A123: ZIPCODE DIAWALI kata (`LIKE pola%`, ESCAPE); urut ZIPCODE lalu kolom lain
// (pyCategorizeSortType ASC sel 13; deterministik). :1 status, :2 pola, :3 batas.
func sqlCariRW(rw string) string {
	return "SELECT ZIPCODE, NOTE, DISTRICTNAME, CITYNAME, PROVINCENAME, NATION FROM (SELECT DISTINCT ZIPCODE, NOTE, DISTRICTNAME," +
		" CITYNAME, PROVINCENAME, NATION FROM " + rw + ` WHERE STS_AKTIF = :1 AND ZIPCODE LIKE :2 ESCAPE '\')` +
		" ORDER BY ZIPCODE, NOTE, DISTRICTNAME, CITYNAME, PROVINCENAME, NATION FETCH FIRST :3 ROWS ONLY"
}

// polaAwalan - awalan LIKE: kata dipangkas, `\` `%` `_` diloloskan, huruf APA ADANYA (kolom
// ZIPCODE juga dibandingkan apa adanya), lalu `%` di belakang. Kosong = "".
func polaAwalan(kata string) string {
	k := strings.TrimSpace(kata)
	if k == "" {
		return ""
	}
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(k) + "%"
}

// panjangIDRiskAddress - prosedur menampung ID di `vID varchar2(15)`: ID lebih panjang gagal
// di Pega; port Go menolak hal yang sama (bukan memotong).
const panjangIDRiskAddress = 15

// periksaIDRiskAddress - ID baru tidak kosong dan <= 15 karakter (pola prosedur).
func periksaIDRiskAddress(id string) error {
	if id == "" || len(id) > panjangIDRiskAddress {
		return fmt.Errorf("repository: ID %s %q tidak sah (kosong atau > %d karakter, vID prosedur)", TabelRiskAddress, id, panjangIDRiskAddress)
	}
	return nil
}

// CariRW - lihat PembacaRW.
func (r *RWOracle) CariRW(ctx context.Context, awalanZip string) ([]models.BarisRW, error) {
	q, err := r.db.Qualify(TabelRW)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlCariRW(q), StatusRWAktif, polaAwalan(awalanZip), BatasSaranRW)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelRW, err)
	}
	defer baris.Close()
	hasil := []models.BarisRW{}
	for baris.Next() {
		var zip, wilayah, kecamatan, kota, provinsi, negara sql.NullString
		if err := baris.Scan(&zip, &wilayah, &kecamatan, &kota, &provinsi, &negara); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelRW, err)
		}
		hasil = append(hasil, models.BarisRW{ZipCode: zip.String, TerritoryName: wilayah.String, DistrictName: kecamatan.String,
			CityName: kota.String, ProvinceName: provinsi.String, NationName: negara.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelRW, err)
	}
	return hasil, nil
}

// sqlIDRiskAddress - ekspresi ID prosedur, persis.
func sqlIDRiskAddress(fungsi, seq string) string {
	return "SELECT " + fungsi + " || LPAD(TO_CHAR(" + seq + ".NEXTVAL), 12, '0') FROM DUAL"
}

// sqlSisipRiskAddress - kolom = urutan parameter prosedur (P_NationName .. P_Postalcode).
func sqlSisipRiskAddress(tabel string) string {
	return "INSERT INTO " + tabel + " (ID, NATIONNAME, PROVINCENAME, DISTRICTNAME, CITYNAME, TERRITORYNAME, TITLE, ADDRESS, POSTALCODE)" +
		" VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9)"
}

func argSisipRiskAddress(id string, a models.AlamatBaru) []any {
	k := db.KosongJadiNil
	return []any{id, k(a.NationName), k(a.ProvinceName), k(a.DistrictName), k(a.CityName), k(a.TerritoryName), k(a.Title),
		k(a.Address), k(a.PostalCode)}
}

// SisipRiskAddress - lihat PenulisRiskAddress.
func (r *RWOracle) SisipRiskAddress(ctx context.Context, tx *db.Tx, a models.AlamatBaru) (string, error) {
	fungsi, err := r.db.Qualify(fungsiSitus)
	if err != nil {
		return "", err
	}
	seq, err := r.db.Qualify(sequenceRiskAddress)
	if err != nil {
		return "", err
	}
	tabel, err := r.db.Qualify(TabelRiskAddress)
	if err != nil {
		return "", err
	}
	qID := sqlIDRiskAddress(fungsi, seq)
	if err := db.PeriksaSQL(qID); err != nil {
		return "", err
	}
	var id sql.NullString
	if err := tx.QueryRowContext(ctx, qID).Scan(&id); err != nil {
		return "", fmt.Errorf("repository: ID %s: %w", TabelRiskAddress, err)
	}
	if err := periksaIDRiskAddress(id.String); err != nil {
		return "", err
	}
	h, err := jalankan(ctx, tx, sqlSisipRiskAddress(tabel), "menyisipkan "+TabelRiskAddress, argSisipRiskAddress(id.String, a)...)
	if err != nil {
		return "", err
	}
	if err := db.PastikanSatuBaris(h, TabelRiskAddress); err != nil {
		return "", err
	}
	return id.String, nil
}
