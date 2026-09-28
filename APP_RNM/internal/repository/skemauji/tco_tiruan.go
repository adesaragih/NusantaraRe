package skemauji

// Tiruan tabel WARISAN dan MASTER Treaty Contract Out untuk uji bertag db
// (tiket 01 tco2; tiket 02 master jenis reasuransi).
//
// ⛔ Tipe kolom warisan mengikuti DEKLARASI `[data DBA]` (dba-procedures.md
// bab DDL): yang disebut NUMBER, DATE, atau CHAR dibuat begitu; sisanya
// VARCHAR2(1000). Tiruan bertipe "lebih benar" akan membuat Oracle mengurai
// angka lebih dulu, dan jalur "teks -> urai -> laporkan yang gagal" tidak
// pernah teruji.
//
// Daftar kolom warisan TIDAK ditulis ulang di sini: ia datang dari repository,
// tempat pembaca migrasi mengambil daftar yang sama.
//
// PROPORTIONALARRG tiruan memuat pula dua kolom mati PROPORTIONALLIST dan
// OBJECT (AC 70) supaya pencacahnya benar-benar berjalan.
//
// Master REINSURANCETYPE ditiru dengan kolom yang SQL korpus sebut (ID, NOTE,
// TYPE, FLAG, NOURUT); tipe fisiknya tidak diketahui DBA - VARCHAR2 generik.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/internal/repository"
)

// namaTabelTiruanTCO adalah tabel yang ditiru dan dibongkar bersama skema uji.
var namaTabelTiruanTCO = []string{
	"TREATYYEAR", "TREATYCONTRACT", "TREATYREINSURER",
	"MTREATYSECURITY", "TREATYBUSINESS", "PROPORTIONALARRG",
	repository.MasterJenisReasuransiTCO,
}

// namaTabelWarisanTCO adalah enam tabel warisan yang dipindahkan migrasi data.
var namaTabelWarisanTCO = namaTabelTiruanTCO[:6]

// lebarCharTCO adalah lebar kolom CHAR `[data DBA]` MTREATYSECURITY.
var lebarCharTCO = map[string]int{
	"TP_TREATY": 2, "REAS_ID": 7, "USER_ID": 99, "REAS_SECURITY": 10,
}

// kolomMatiTiruanTCO adalah kolom warisan yang ada di DDL tetapi tidak dibawa.
var kolomMatiTiruanTCO = []string{"PROPORTIONALLIST VARCHAR2(1000)", "OBJECT VARCHAR2(50)"}

func tipeTiruanTCO(tabel, kolom string) string {
	switch repository.TipeWarisanTCO(tabel, kolom) {
	case repository.WarisanAngka:
		return "NUMBER"
	case repository.WarisanTanggal:
		return "DATE"
	case repository.WarisanChar:
		return fmt.Sprintf("CHAR(%d)", lebarCharTCO[kolom])
	}
	return "VARCHAR2(1000)"
}

// ddlTiruanTCO membuat DDL keenam tabel warisan dan master jenis reasuransi.
func ddlTiruanTCO(skema string) []string {
	var out []string
	for _, tabel := range namaTabelWarisanTCO {
		var kolom []string
		for _, k := range repository.KolomWarisanTCO(tabel) {
			kolom = append(kolom, k+" "+tipeTiruanTCO(tabel, k))
		}
		if tabel == "PROPORTIONALARRG" {
			kolom = append(kolom, kolomMatiTiruanTCO...)
		}
		out = append(out, fmt.Sprintf("CREATE TABLE %s.%s (%s)", skema, tabel, strings.Join(kolom, ", ")))
	}
	out = append(out, fmt.Sprintf(
		"CREATE TABLE %s.%s (ID VARCHAR2(1000), NOTE VARCHAR2(1000), TYPE VARCHAR2(1000), FLAG VARCHAR2(1000), NOURUT VARCHAR2(1000))",
		skema, repository.MasterJenisReasuransiTCO))
	return out
}

// IsiWarisanTCO mengisi satu tabel warisan tiruan dengan baris buatan.
//
// Nilai diberikan sebagai TEKS per kolom, seperti tersimpan di warisan.
// Kolom DATE diisi lewat TO_DATE 'YYYY-MM-DD HH24:MI:SS'; kolom NUMBER
// di-bind teks (NLS sesi disamakan lebih dulu). Teks kosong menjadi NULL.
// Kolom yang tidak dikenal pembaca migrasi (kolom mati) tetap boleh diisi.
//
// ⛔ Nol nama orang, nol nomor polis nyata: seluruh nilai fixture berawalan UJI.
func IsiWarisanTCO(ctx context.Context, db *sql.DB, skema, tabel string, baris []map[string]string) error {
	if err := samakanNLS(ctx, db); err != nil {
		return err
	}
	dikenal := map[string]bool{}
	for _, k := range repository.KolomWarisanTCO(tabel) {
		dikenal[k] = true
	}
	for _, b := range baris {
		var kolom, penampung []string
		var arg []any
		for _, k := range repository.KolomWarisanTCO(tabel) {
			v, ada := b[k]
			if !ada {
				continue
			}
			kolom = append(kolom, k)
			i := len(arg) + 1
			if repository.TipeWarisanTCO(tabel, k) == repository.WarisanTanggal {
				penampung = append(penampung, fmt.Sprintf("TO_DATE(:%d, 'YYYY-MM-DD HH24:MI:SS')", i))
			} else {
				penampung = append(penampung, fmt.Sprintf(":%d", i))
			}
			if v == "" {
				arg = append(arg, nil)
			} else {
				arg = append(arg, v)
			}
		}
		for k, v := range b {
			if dikenal[k] {
				continue
			}
			kolom = append(kolom, k)
			penampung = append(penampung, fmt.Sprintf(":%d", len(arg)+1))
			if v == "" {
				arg = append(arg, nil)
			} else {
				arg = append(arg, v)
			}
		}
		q := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s)", skema, tabel,
			strings.Join(kolom, ", "), strings.Join(penampung, ", "))
		if _, err := db.ExecContext(ctx, q, arg...); err != nil {
			return fmt.Errorf("skemauji: mengisi tiruan %s: %w", tabel, err)
		}
	}
	return nil
}

// JenisReasuransiUji adalah satu baris fixture master jenis reasuransi.
type JenisReasuransiUji struct {
	ID, Note, Tipe, Flag string
}

// IsiJenisReasuransiTCO mengisi tiruan REINSURANCETYPE.
//
// ⚠️ Fixture yang berguna memuat yang HARUS tersaring keluar: ID blacklist,
// ID yang hanya BERAWALAN blacklist, Flag bukan active, Type di luar 1-3.
func IsiJenisReasuransiTCO(ctx context.Context, db *sql.DB, skema string, baris []JenisReasuransiUji) error {
	q := fmt.Sprintf("INSERT INTO %s.%s (ID, NOTE, TYPE, FLAG) VALUES (:1, :2, :3, :4)",
		skema, repository.MasterJenisReasuransiTCO)
	for _, b := range baris {
		if _, err := db.ExecContext(ctx, q, b.ID, b.Note, b.Tipe, b.Flag); err != nil {
			return fmt.Errorf("skemauji: mengisi tiruan master jenis reasuransi: %w", err)
		}
	}
	return nil
}
