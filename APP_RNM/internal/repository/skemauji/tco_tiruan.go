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
	repository.MasterJenisReasuransiTCO, repository.MasterGrupTreatyTCO,
	// Tiket 12: master kategori lampiran.
	repository.MasterKategoriLampiranTCO,
	// Tiket 05: master reinsurer.
	repository.MasterReinsurerAgentTCO,
	// Tiket 07: master bisnis.
	repository.MasterBusinessTCO,
	// Tiket 08: master jenis klausul dan pemilih ExclutionTreaty.
	repository.MasterJenisKlausulTCO, repository.MasterOccupationTCO, repository.MasterClauseTCO,
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
	// Master grup treaty: kolom yang SQL korpus sebut (FetchTreatyGroupOLDID).
	out = append(out, fmt.Sprintf(
		"CREATE TABLE %s.%s (ID VARCHAR2(1000), TREATYGROUPNAME VARCHAR2(1000), OLDID VARCHAR2(1000))",
		skema, repository.MasterGrupTreatyTCO))
	// Tiket 12: master kategori lampiran - hanya `NOTE` yang terbukti dipakai
	// (`SetCategoryAttachTreatyin.xml` b500); ID ditiru sebagai kolom bebas.
	out = append(out, fmt.Sprintf(
		"CREATE TABLE %s.%s (ID VARCHAR2(1000), NOTE VARCHAR2(1000))",
		skema, repository.MasterKategoriLampiranTCO))
	// Tiket 05: master reinsurer - kolom yang SQL korpus sebut (ID, CLIENTNAME,
	// CLIENTID) + STATUSACTIVE dari properti RD (OQ-TCO-12).
	out = append(out, fmt.Sprintf(
		"CREATE TABLE %s.%s (ID VARCHAR2(1000), CLIENTNAME VARCHAR2(1000), CLIENTID VARCHAR2(1000), STATUSACTIVE VARCHAR2(10))",
		skema, repository.MasterReinsurerAgentTCO))
	// Tiket 07: master bisnis - kolom yang SQL korpus sebut.
	out = append(out, fmt.Sprintf(
		"CREATE TABLE %s.%s (ID VARCHAR2(1000), OLDID VARCHAR2(1000), NOTE VARCHAR2(1000), BUSINESSGROUPID VARCHAR2(1000))",
		skema, repository.MasterBusinessTCO))
	// Tiket 08: kolom = nama properti RD / SQL korpus (OQ-TCO-16).
	out = append(out,
		fmt.Sprintf("CREATE TABLE %s.%s (ID VARCHAR2(1000), DESCNAME VARCHAR2(1000), ISXOL VARCHAR2(10), STATUSAKTIF VARCHAR2(10))",
			skema, repository.MasterJenisKlausulTCO),
		fmt.Sprintf("CREATE TABLE %s.%s (ID VARCHAR2(1000), NAME VARCHAR2(1000), TYPE VARCHAR2(100))",
			skema, repository.MasterOccupationTCO),
		fmt.Sprintf("CREATE TABLE %s.%s (ID VARCHAR2(1000), INFO VARCHAR2(1000), TYPE VARCHAR2(100))",
			skema, repository.MasterClauseTCO))
	// Tiket 11: master kurs - SELURUH kolom VARCHAR2 [data DBA]; master mata
	// uang - kolom yang SQL korpus sebut (`GetCurrencyIDByName`).
	out = append(out,
		fmt.Sprintf("CREATE TABLE %s.%s (TOIDR VARCHAR2(100), TOUSD VARCHAR2(100), IDCURRENCY VARCHAR2(100), "+
			"CURRENCY VARCHAR2(100), QUARTER VARCHAR2(10), STARTDATE VARCHAR2(100), ENDDATE VARCHAR2(100))",
			skema, repository.MasterKursTahunanTCO),
		fmt.Sprintf("CREATE TABLE %s.%s (ID VARCHAR2(100), OLDID VARCHAR2(100), CURRENCY VARCHAR2(100), "+
			"CURRENCYSYMBOL VARCHAR2(100))", skema, repository.MasterMataUangTCO))
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

// GrupTreatyUji adalah satu baris fixture master grup treaty.
type GrupTreatyUji struct {
	ID, TreatyGroupName string
}

// IsiGrupTreatyTCO mengisi tiruan TREATYGROUP.
func IsiGrupTreatyTCO(ctx context.Context, db *sql.DB, skema string, baris []GrupTreatyUji) error {
	q := fmt.Sprintf("INSERT INTO %s.%s (ID, TREATYGROUPNAME) VALUES (:1, :2)",
		skema, repository.MasterGrupTreatyTCO)
	for _, b := range baris {
		if _, err := db.ExecContext(ctx, q, b.ID, b.TreatyGroupName); err != nil {
			return fmt.Errorf("skemauji: mengisi tiruan master grup treaty: %w", err)
		}
	}
	return nil
}

// IsiKategoriLampiranTCO mengisi tiruan CATEGORY_ATTACH_REAS (tiket 12).
func IsiKategoriLampiranTCO(ctx context.Context, db *sql.DB, skema string, note []string) error {
	q := fmt.Sprintf("INSERT INTO %s.%s (ID, NOTE) VALUES (:1, :2)",
		skema, repository.MasterKategoriLampiranTCO)
	for i, n := range note {
		if _, err := db.ExecContext(ctx, q, fmt.Sprintf("UJI-%02d", i+1), n); err != nil {
			return fmt.Errorf("skemauji: mengisi tiruan master kategori lampiran: %w", err)
		}
	}
	return nil
}

// AgentUji adalah satu baris fixture master reinsurer.
type AgentUji struct {
	ID, ClientName, ClientID, StatusActive string
}

// IsiAgentTCO mengisi tiruan AGENT (tiket 05).
func IsiAgentTCO(ctx context.Context, db *sql.DB, skema string, baris []AgentUji) error {
	q := fmt.Sprintf("INSERT INTO %s.%s (ID, CLIENTNAME, CLIENTID, STATUSACTIVE) VALUES (:1, :2, :3, :4)",
		skema, repository.MasterReinsurerAgentTCO)
	for _, b := range baris {
		if _, err := db.ExecContext(ctx, q, b.ID, b.ClientName, b.ClientID, b.StatusActive); err != nil {
			return fmt.Errorf("skemauji: mengisi tiruan master reinsurer: %w", err)
		}
	}
	return nil
}

// BusinessUji adalah satu baris fixture master bisnis.
type BusinessUji struct{ ID, Note, BusinessGroupID string }

// IsiBusinessTCO mengisi tiruan BUSINESS (tiket 07).
func IsiBusinessTCO(ctx context.Context, db *sql.DB, skema string, baris []BusinessUji) error {
	q := fmt.Sprintf("INSERT INTO %s.%s (ID, NOTE, BUSINESSGROUPID) VALUES (:1, :2, :3)",
		skema, repository.MasterBusinessTCO)
	for _, b := range baris {
		if _, err := db.ExecContext(ctx, q, b.ID, b.Note, kosongJadiNilUji(b.BusinessGroupID)); err != nil {
			return fmt.Errorf("skemauji: mengisi tiruan master bisnis: %w", err)
		}
	}
	return nil
}

func kosongJadiNilUji(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// JenisKlausulUji adalah satu baris fixture master TREATYDESC.
type JenisKlausulUji struct{ ID, DescName, IsXOL string }

// IsiJenisKlausulTCO mengisi tiruan TREATYDESC (tiket 08).
func IsiJenisKlausulTCO(ctx context.Context, db *sql.DB, skema string, baris []JenisKlausulUji) error {
	q := fmt.Sprintf("INSERT INTO %s.%s (ID, DESCNAME, ISXOL, STATUSAKTIF) VALUES (:1, :2, :3, '1')",
		skema, repository.MasterJenisKlausulTCO)
	for _, b := range baris {
		if _, err := db.ExecContext(ctx, q, b.ID, b.DescName, b.IsXOL); err != nil {
			return fmt.Errorf("skemauji: mengisi tiruan master jenis klausul: %w", err)
		}
	}
	return nil
}

// IsiPilihanKlausulTCO mengisi tiruan OCCUPATION atau CLAUSE (tiket 08).
func IsiPilihanKlausulTCO(ctx context.Context, db *sql.DB, skema, master string, pilihan map[string]string) error {
	kolom := "NAME"
	if master == repository.MasterClauseTCO {
		kolom = "INFO"
	}
	q := fmt.Sprintf("INSERT INTO %s.%s (ID, %s, TYPE) VALUES (:1, :2, :3)", skema, master, kolom)
	for id, nama := range pilihan {
		if _, err := db.ExecContext(ctx, q, id, nama, repository.TipeFireTCO); err != nil {
			return fmt.Errorf("skemauji: mengisi tiruan %s: %w", master, err)
		}
	}
	return nil
}

// KursUji adalah satu baris tiruan master kurs (teks, seperti warisannya).
type KursUji struct{ ToIDR, IDCurrency, Currency, Quarter, StartDate, EndDate string }

// IsiKursTCO mengisi tiruan TREATYEXCHANGEYEARLY (tiket 11).
func IsiKursTCO(ctx context.Context, db *sql.DB, skema string, baris []KursUji) error {
	q := fmt.Sprintf("INSERT INTO %s.%s (TOIDR, IDCURRENCY, CURRENCY, QUARTER, STARTDATE, ENDDATE) VALUES (:1, :2, :3, :4, :5, :6)",
		skema, repository.MasterKursTahunanTCO)
	for _, b := range baris {
		if _, err := db.ExecContext(ctx, q, b.ToIDR, b.IDCurrency, b.Currency, b.Quarter, b.StartDate, b.EndDate); err != nil {
			return fmt.Errorf("skemauji: mengisi tiruan master kurs: %w", err)
		}
	}
	return nil
}

// IsiMataUangTCO mengisi tiruan CURRENCY: kode -> pengenal (tiket 11).
func IsiMataUangTCO(ctx context.Context, db *sql.DB, skema string, kodeKeID map[string]string) error {
	q := fmt.Sprintf("INSERT INTO %s.%s (ID, CURRENCY) VALUES (:1, :2)", skema, repository.MasterMataUangTCO)
	for kode, id := range kodeKeID {
		if _, err := db.ExecContext(ctx, q, id, kode); err != nil {
			return fmt.Errorf("skemauji: mengisi tiruan mata uang: %w", err)
		}
	}
	return nil
}
