package repository

// Klausul - SATU tabel `PROPORTIONALARRG` untuk 25 jenis - tiket 08.
//
// Untuk apa berkas ini: baca/tulis klausul (meniru logika
// `PEGA_PROPORTIONALARRG` 35 parameter dan `PEGA_M_PROPORTIONALARRG_CHILD` 26
// parameter - KEDUANYA menulis tabel yang sama `[data DBA]` - tanpa memanggil
// prosedur mana pun) dan tiga master yang dibaca saja: `TREATYDESC` (jenis),
// `OCCUPATION`, `CLAUSE` (pemilih ExclutionTreaty).
//
// ⛔ Lingkup daftar = (TREATYYEARID, TREATYDESCID, PARENTREINSTYPEID):
// `BrowseTreatyArrangement_EPI_RD` / `_Limit_RD` saringan `.TreatyDescID`,
// `.ParentReinsTypeID` ("00" untuk induk), `.TreatyYearID` (`C AND D AND F1`).
//
// ⛔ AC 25: baris ANAK menulis NULL pada sembilan kolom khusus induk - ditegakkan
// DI SINI, bukan diserahkan ke pemanggil.
//
// Dibaca sesudah: tco_reinsurer.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/internal/models"
)

// Master yang dibaca saja.
const (
	// MasterJenisKlausulTCO - `BrowseTreatyDesc_RD` (kelas
	// `ASM-FW-GISFW-Int-TREATYDESC`; `.ID` ASC, `.DescName`, `.IsXOL`,
	// `.StatusAktif`; saringan `.IsXOL = Param.IsXOL`). Kolom fisik = nama
	// properti, dikonfirmasi [keputusan work owner 29-09-2026] - nol SQL korpus
	// menyebutnya (OQ-TCO-16, ditutup).
	MasterJenisKlausulTCO = "TREATYDESC"
	// MasterOccupationTCO - `BrowseOccupationFIRE_RD` (`.Type = "FIRE"`,
	// `.ID` -> ID_Occupation, `.Name` -> Occupation). Kolom `ID`, `NAME`, `TYPE`
	// terbukti di SQL korpus modul lain.
	MasterOccupationTCO = "OCCUPATION"
	// MasterClauseTCO - `BrowseFireClauseFacIn_RD` (`.Type = "FIRE"`, urut
	// `.Info` ASC; `.ID` -> ID_Clause, `.Info` -> Clause). Kolom fisik = nama
	// properti, dikonfirmasi [keputusan work owner 29-09-2026] (OQ-TCO-16, ditutup).
	MasterClauseTCO = "CLAUSE"
	// TipeFireTCO - saringan kedua pemilih exclusion.
	TipeFireTCO = "FIRE"
)

const batasPilihanKlausulTCO = 100

// ErrKlausulTidakAda - klausul bukan milik tahun treaty itu.
var ErrKlausulTidakAda = errors.New("repository: klausul tidak ditemukan pada tahun treaty ini")

// ErrPilihanMasterTidakAda - ID occupation/clause tidak ada di master FIRE.
var ErrPilihanMasterTidakAda = errors.New("repository: pilihan tidak ada di master")

// JenisKlausulMasterTCO adalah satu baris master `TREATYDESC`.
type JenisKlausulMasterTCO struct{ ID, DescName, IsXOL, StatusAktif string }

// PilihanMasterTCO adalah satu pilihan occupation/clause.
type PilihanMasterTCO struct{ ID, Nama string }

// MasterKlausulPilihan membaca ketiga master klausul.
type MasterKlausulPilihan struct{ db *DB }

// NewMasterKlausulPilihan menyusun pembacanya.
func NewMasterKlausulPilihan(db *DB) *MasterKlausulPilihan { return &MasterKlausulPilihan{db: db} }

func sqlJenisKlausulTCO(tabel string) string {
	return fmt.Sprintf(`SELECT ID, DESCNAME, ISXOL, STATUSAKTIF FROM %s
	 WHERE (:1 IS NULL OR ISXOL = :2) ORDER BY ID ASC`, tabel)
}

// JenisKlausul membaca master jenis; `isXOL` kosong = seluruhnya.
func (m *MasterKlausulPilihan) JenisKlausul(ctx context.Context, isXOL string) ([]JenisKlausulMasterTCO, error) {
	tabel, err := m.db.Qualify(MasterJenisKlausulTCO)
	if err != nil {
		return nil, err
	}
	q := sqlJenisKlausulTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := m.db.bacaTCO(ctx).QueryContext(ctx, q, kosongJadiNil(isXOL), kosongJadiNil(isXOL))
	if err != nil {
		return nil, fmt.Errorf("repository: membaca master %s: %w", MasterJenisKlausulTCO, err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []JenisKlausulMasterTCO
	for rows.Next() {
		var n [4]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3]); err != nil {
			return nil, err
		}
		hasil = append(hasil, JenisKlausulMasterTCO{ID: n[0].String, DescName: n[1].String, IsXOL: n[2].String,
			StatusAktif: n[3].String})
	}
	return hasil, rows.Err()
}

func sqlCariPilihanTCO(tabel, kolomNama string) string {
	return fmt.Sprintf(`SELECT ID, %s FROM %s WHERE TYPE = :1 AND UPPER(%s) LIKE :2 ESCAPE '\'
	 ORDER BY %s, ID FETCH FIRST %d ROWS ONLY`, kolomNama, tabel, kolomNama, kolomNama, batasPilihanKlausulTCO)
}

func sqlAmbilPilihanTCO(tabel, kolomNama string) string {
	return fmt.Sprintf(`SELECT ID, %s FROM %s WHERE ID = :1 AND TYPE = :2`, kolomNama, tabel)
}

// kolomNamaPilihan - kolom tampil tiap master pemilih.
var kolomNamaPilihan = map[string]string{MasterOccupationTCO: "NAME", MasterClauseTCO: "INFO"}

// CariPilihan membaca pilihan FIRE yang namanya memuat teks.
func (m *MasterKlausulPilihan) CariPilihan(ctx context.Context, master, teks string) ([]PilihanMasterTCO, error) {
	kolom, ok := kolomNamaPilihan[master]
	if !ok {
		return nil, fmt.Errorf("repository: master pilihan %q tidak dikenal", master)
	}
	tabel, err := m.db.Qualify(master)
	if err != nil {
		return nil, err
	}
	q := sqlCariPilihanTCO(tabel, kolom)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := m.db.bacaTCO(ctx).QueryContext(ctx, q, TipeFireTCO, polaLikeTCO(teks))
	if err != nil {
		return nil, fmt.Errorf("repository: membaca master %s: %w", master, err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []PilihanMasterTCO
	for rows.Next() {
		var id, nama sql.NullString
		if err := rows.Scan(&id, &nama); err != nil {
			return nil, err
		}
		hasil = append(hasil, PilihanMasterTCO{ID: id.String, Nama: nama.String})
	}
	return hasil, rows.Err()
}

// AmbilPilihan membaca satu pilihan FIRE.
func (m *MasterKlausulPilihan) AmbilPilihan(ctx context.Context, master, id string) (PilihanMasterTCO, error) {
	kolom, ok := kolomNamaPilihan[master]
	if !ok {
		return PilihanMasterTCO{}, fmt.Errorf("repository: master pilihan %q tidak dikenal", master)
	}
	tabel, err := m.db.Qualify(master)
	if err != nil {
		return PilihanMasterTCO{}, err
	}
	q := sqlAmbilPilihanTCO(tabel, kolom)
	if err := PeriksaSQL(q); err != nil {
		return PilihanMasterTCO{}, err
	}
	var gotID, nama sql.NullString
	err = m.db.bacaTCO(ctx).QueryRowContext(ctx, q, id, TipeFireTCO).Scan(&gotID, &nama)
	if errors.Is(err, sql.ErrNoRows) {
		return PilihanMasterTCO{}, fmt.Errorf("%w %s: %q", ErrPilihanMasterTidakAda, master, id)
	}
	if err != nil {
		return PilihanMasterTCO{}, fmt.Errorf("repository: membaca master %s %s: %w", master, id, err)
	}
	return PilihanMasterTCO{ID: gotID.String, Nama: nama.String}, nil
}

// MasterKlausulTCO membaca dan menulis tabel WARISAN `PROPORTIONALARRG` (tco4).
//
// ⚠️ Tipe CAMPUR `[data DBA]`: `TREATYLIMIT, COINS_MIN, COINS_MAX, MORERP,
// MOREUSD` NUMBER; `RP, USD, PCT, PCTME` VARCHAR2(1000); `KURS` tidak disebut
// DBA (diperlakukan teks); `TGLUPDATE` DATE. Kolom teks dibaca APA ADANYA dan
// diurai di Go (titik atau koma), ditulis sebagai teks desimal
// (`TulisDesimalWarisanTCO`) - `TO_CHAR` berformat angka atasnya akan gagal.
type MasterKlausulTCO struct{ db *DB }

// NewMasterKlausulTCO menyusun gudangnya.
func NewMasterKlausulTCO(db *DB) *MasterKlausulTCO { return &MasterKlausulTCO{db: db} }

// kolomKlausulTCO - 35 kolom, urutan parameter `PEGA_PROPORTIONALARRG`; SATU
// sumber dengan skema warisan (`KolomWarisanTCO`) - temuan /code-review.
var kolomKlausulTCO = KolomWarisanTCO(warisanKlausulTCO)

// KolomKhususIndukTCO - sembilan kolom yang NULL pada baris anak (AC 25).
var KolomKhususIndukTCO = []string{"ID_OCCUPATION", "OCCUPATION", "ID_CLAUSE", "CLAUSE", "TREATYLIMIT", "COINS_MIN",
	"COINS_MAX", "MORERP", "MOREUSD"}

// kolomDesimalKlausul - kolom bernilai desimal; kolomAngkaKlausul - yang NUMBER.
var kolomDesimalKlausul = map[string]bool{"KURS": true, "PCT": true, "PCTME": true, "RP": true, "USD": true,
	"TREATYLIMIT": true, "COINS_MIN": true, "COINS_MAX": true, "MORERP": true, "MOREUSD": true}

var kolomAngkaKlausul = func() map[string]bool {
	hasil := map[string]bool{}
	for _, k := range kolomKlausulTCO {
		if TipeWarisanTCO(warisanKlausulTCO, k) == WarisanAngka {
			hasil[k] = true
		}
	}
	return hasil
}()

func pilihKlausulTCO() string {
	var b strings.Builder
	for i, k := range kolomKlausulTCO {
		if i > 0 {
			b.WriteString(", ")
		}
		switch {
		case kolomAngkaKlausul[k]:
			b.WriteString(fmt.Sprintf(fmtDesimal, k))
		case k == "TGLUPDATE":
			b.WriteString("TO_CHAR(TGLUPDATE, 'YYYY-MM-DD HH24:MI:SS')")
		default:
			b.WriteString(k)
		}
	}
	return b.String()
}

// kolomDariMedan - medan klausul (VERBATIM Pega) -> kolom DDL.
var kolomDariMedan = map[string]string{
	models.MedanReinsTypeID: "REINSTYPEID", models.MedanLine: "LINE", models.MedanRp: "RP", models.MedanUsd: "USD",
	models.MedanPct: "PCT", models.MedanPctMe: "PCTME", models.MedanYdcf: "YDCF", models.MedanMethod: "METHOD",
	models.MedanTerritorialLimit: "TERRITORIALLIMIT", models.MedanCoInsMin: "COINS_MIN", models.MedanCoInsMax: "COINS_MAX",
	models.MedanTreatyLimit: "TREATYLIMIT", models.MedanIDOccupation: "ID_OCCUPATION", models.MedanOccupation: "OCCUPATION",
	models.MedanIDClause: "ID_CLAUSE", models.MedanClause: "CLAUSE", models.MedanLayer: "LAYER",
}

// sqlCariDobelKlausulTCO - kandidat baris LAIN pada lingkup yang sama; kunci
// jenis dibandingkan di Go (tco4: `PCT`/`RP`/`USD` teks warisan - `12.5` dan
// `12,50` sama nilainya, tidak sama teksnya).
func sqlCariDobelKlausulTCO(tabel string) string {
	return fmt.Sprintf(`SELECT %s FROM %s
	 WHERE TREATYYEARID = :1 AND TREATYDESCID = :2 AND PARENTREINSTYPEID = :3
	   AND (:4 IS NULL OR ID <> :5)
	 ORDER BY ID`, pilihKlausulTCO(), tabel)
}

// medanSamaKlausulTCO - satu medan kunci bernilai sama. Kosong tidak pernah
// sama (padanan `=` SQL atas NULL). Nama medan dari daftar putih.
func medanSamaKlausulTCO(a, b models.KlausulTreaty, medan string) (bool, error) {
	kolom, ok := kolomDariMedan[medan]
	if !ok {
		return false, fmt.Errorf("repository: medan kunci dobel %q tidak dikenal", medan)
	}
	if kolomDesimalKlausul[kolom] {
		da, _, okA := UraiDesimalWarisanTCO(models.NilaiMedanKlausul(a, medan))
		db, _, okB := UraiDesimalWarisanTCO(models.NilaiMedanKlausul(b, medan))
		return okA && okB && da != nil && db != nil && da.Cmp(db) == 0, nil
	}
	x, y := strings.TrimSpace(models.NilaiMedanKlausul(a, medan)), strings.TrimSpace(models.NilaiMedanKlausul(b, medan))
	return x != "" && x == y, nil
}

// sqlDaftarKlausulTCO - urut `ID ASC`: identitas dari sequence, jadi urutan
// masukan dipertahankan saat dibaca kembali (AC 31).
func sqlDaftarKlausulTCO(tabel string) string {
	return fmt.Sprintf(`SELECT %s FROM %s
	 WHERE TREATYYEARID = :1 AND TREATYDESCID = :2 AND PARENTREINSTYPEID = :3 ORDER BY ID ASC`, pilihKlausulTCO(), tabel)
}

func sqlAmbilKlausulTCO(tabel string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE TREATYYEARID = :1 AND ID = :2`, pilihKlausulTCO(), tabel)
}

func sqlIndukKlausulTCO(tabel string) string {
	return fmt.Sprintf(`SELECT %s FROM %s
	 WHERE TREATYYEARID = :1 AND TREATYDESCID = :2 AND PARENTREINSTYPEID = :3 AND REINSTYPEID = :4
	 ORDER BY ID ASC FETCH FIRST 1 ROWS ONLY`, pilihKlausulTCO(), tabel)
}

// sqlPctAnakLainTCO - `PCT` VARCHAR2 dibaca apa adanya (tco4).
func sqlPctAnakLainTCO(tabel string) string {
	return fmt.Sprintf(`SELECT ID, PCT FROM %s
	 WHERE TREATYYEARID = :1 AND TREATYDESCID = :2 AND PARENTREINSTYPEID = :3 AND (:4 IS NULL OR ID <> :5) FOR UPDATE`,
		tabel)
}

func sqlSisipKlausulTCO(tabel string) string {
	var ph []string
	for i := range kolomKlausulTCO {
		ph = append(ph, fmt.Sprintf(":%d", i+1))
	}
	return fmt.Sprintf("INSERT INTO %s\n\t       (%s)\n\tVALUES (%s)", tabel, strings.Join(kolomKlausulTCO, ", "),
		strings.Join(ph, ", "))
}

// kunciTetapKlausul - kolom yang tidak berubah saat diperbarui.
var kunciTetapKlausul = map[string]bool{"ID": true, "TREATYYEARID": true, "TREATYDESCID": true, "PARENTREINSTYPEID": true}

// sqlPerbaruiKlausulTCO - SELURUH kolom non-kunci, dibatasi tahunnya.
func sqlPerbaruiKlausulTCO(tabel string) string {
	var set []string
	n := 0
	for _, k := range kolomKlausulTCO {
		if kunciTetapKlausul[k] {
			continue
		}
		n++
		set = append(set, fmt.Sprintf("%s = :%d", k, n))
	}
	return fmt.Sprintf("UPDATE %s\n\t   SET %s\n\t WHERE ID = :%d AND TREATYYEARID = :%d", tabel, strings.Join(set, ", "), n+1, n+2)
}

func pindaiKlausulTCO(baca interface{ Scan(...any) error }) (models.KlausulTreaty, error) {
	n := make([]sql.NullString, len(kolomKlausulTCO))
	tujuan := make([]any, len(n))
	for i := range n {
		tujuan[i] = &n[i]
	}
	if err := baca.Scan(tujuan...); err != nil {
		return models.KlausulTreaty{}, err
	}
	nilai := map[string]sql.NullString{}
	for i, k := range kolomKlausulTCO {
		nilai[k] = n[i]
	}
	s := func(k string) string { return nilai[k].String }
	k := models.KlausulTreaty{ID: s("ID"), TreatyYear: s("TREATYYEAR"), TreatyYearID: s("TREATYYEARID"),
		TreatyGroupID: s("TREATYGROUPID"), TreatyGroupName: s("TREATYGROUPNAME"), TreatyDescID: s("TREATYDESCID"),
		TreatyDescName: s("TREATYDESCNAME"), ReinsTypeID: s("REINSTYPEID"), ReinsTypeName: s("REINSTYPENAME"),
		Layer: s("LAYER"), LayerPart: s("LAYERPART"), LayerPartType: s("LAYERPARTTYPE"), LayerType: s("LAYERTYPE"),
		UserID: s("USERID"), Line: s("LINE"), Ydcf: s("YDCF"), Method: s("METHOD"),
		TerritorialLimit: s("TERRITORIALLIMIT"), ParentReinsTypeID: s("PARENTREINSTYPEID"),
		SpreadingOrder: s("SPREADINGORDER"), IDOccupation: s("ID_OCCUPATION"), Occupation: s("OCCUPATION"),
		IDClause: s("ID_CLAUSE"), Clause: s("CLAUSE")}
	desimal := []struct {
		kolom string
		ke    **apd.Decimal
	}{{"KURS", &k.Kurs}, {"PCT", &k.Pct}, {"PCTME", &k.PctMe}, {"RP", &k.Rp}, {"USD", &k.Usd},
		{"TREATYLIMIT", &k.TreatyLimit}, {"COINS_MIN", &k.CoinsMin}, {"COINS_MAX", &k.CoinsMax},
		{"MORERP", &k.MoreRp}, {"MOREUSD", &k.MoreUsd}}
	for _, d := range desimal {
		baca := desimalWarisanTCO
		if kolomAngkaKlausul[d.kolom] {
			baca = desimalDariTeks
		}
		v, err := baca(nilai[d.kolom], d.kolom)
		if err != nil {
			return k, err
		}
		*d.ke = v
	}
	var err error
	k.TglUpdate, err = uraiTanggalTeks(nilai["TGLUPDATE"], "TGLUPDATE")
	return k, err
}

// nilaiKolomKlausul menyusun argumen untuk ke-35 kolom; baris ANAK menulis
// NULL pada sembilan kolom khusus induk (AC 25).
func nilaiKolomKlausul(k models.KlausulTreaty) map[string]any {
	d := TulisDesimalWarisanTCO // RP, USD, PCT, PCTME, KURS: teks warisan
	n := desimalJadiNil         // TREATYLIMIT, COINS_*, MORE*: NUMBER
	v := map[string]any{"ID": k.ID, "TREATYYEAR": kosongJadiNil(k.TreatyYear), "TREATYYEARID": kosongJadiNil(k.TreatyYearID),
		"TREATYGROUPID": kosongJadiNil(k.TreatyGroupID), "TREATYGROUPNAME": kosongJadiNil(k.TreatyGroupName),
		"TREATYDESCID": kosongJadiNil(k.TreatyDescID), "TREATYDESCNAME": kosongJadiNil(k.TreatyDescName),
		"REINSTYPEID": kosongJadiNil(k.ReinsTypeID), "REINSTYPENAME": kosongJadiNil(k.ReinsTypeName),
		"LAYER": kosongJadiNil(k.Layer), "LAYERPART": kosongJadiNil(k.LayerPart),
		"LAYERPARTTYPE": kosongJadiNil(k.LayerPartType), "LAYERTYPE": kosongJadiNil(k.LayerType), "KURS": d(k.Kurs),
		"TGLUPDATE": tanggalJadiNil(k.TglUpdate), "USERID": kosongJadiNil(k.UserID), "LINE": kosongJadiNil(k.Line),
		"PCT": d(k.Pct), "PCTME": d(k.PctMe), "YDCF": kosongJadiNil(k.Ydcf), "METHOD": kosongJadiNil(k.Method),
		"TERRITORIALLIMIT": kosongJadiNil(k.TerritorialLimit), "PARENTREINSTYPEID": kosongJadiNil(k.ParentReinsTypeID),
		"SPREADINGORDER": kosongJadiNil(k.SpreadingOrder), "RP": d(k.Rp), "USD": d(k.Usd),
		"ID_OCCUPATION": kosongJadiNil(k.IDOccupation), "OCCUPATION": kosongJadiNil(k.Occupation),
		"ID_CLAUSE": kosongJadiNil(k.IDClause), "CLAUSE": kosongJadiNil(k.Clause), "TREATYLIMIT": n(k.TreatyLimit),
		"COINS_MIN": n(k.CoinsMin), "COINS_MAX": n(k.CoinsMax), "MORERP": n(k.MoreRp), "MOREUSD": n(k.MoreUsd)}
	if models.KlausulAnakTCO(k) {
		for _, kolom := range KolomKhususIndukTCO {
			v[kolom] = nil
		}
	}
	return v
}

func (m *MasterKlausulTCO) bacaBanyak(ctx context.Context, q string, arg ...any) ([]models.KlausulTreaty, error) {
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := m.db.bacaTCO(ctx).QueryContext(ctx, q, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca klausul: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []models.KlausulTreaty
	for rows.Next() {
		k, err := pindaiKlausulTCO(rows)
		if err != nil {
			return nil, err
		}
		hasil = append(hasil, k)
	}
	return hasil, rows.Err()
}

// Daftar membaca klausul satu tahun, satu jenis, satu induk ("00" = induk).
func (m *MasterKlausulTCO) Daftar(ctx context.Context, tahunID, descID, parentReinsTypeID string) ([]models.KlausulTreaty, error) {
	tabel, err := m.db.Qualify(TabelKlausulTCO)
	if err != nil {
		return nil, err
	}
	return m.bacaBanyak(ctx, sqlDaftarKlausulTCO(tabel), tahunID, descID, parentReinsTypeID)
}

// Ambil membaca satu klausul milik tahun itu.
func (m *MasterKlausulTCO) Ambil(ctx context.Context, tahunID, id string) (models.KlausulTreaty, error) {
	tabel, err := m.db.Qualify(TabelKlausulTCO)
	if err != nil {
		return models.KlausulTreaty{}, err
	}
	baris, err := m.bacaBanyak(ctx, sqlAmbilKlausulTCO(tabel), tahunID, id)
	if err != nil {
		return models.KlausulTreaty{}, err
	}
	if len(baris) == 0 {
		return models.KlausulTreaty{}, ErrKlausulTidakAda
	}
	return baris[0], nil
}

// Induk membaca baris induk (sentinel "00") berjenis reasuransi itu.
func (m *MasterKlausulTCO) Induk(ctx context.Context, tahunID, descID, reinsTypeID string) (models.KlausulTreaty, error) {
	tabel, err := m.db.Qualify(TabelKlausulTCO)
	if err != nil {
		return models.KlausulTreaty{}, err
	}
	baris, err := m.bacaBanyak(ctx, sqlIndukKlausulTCO(tabel), tahunID, descID, models.ParentReinsTypeTanpaInduk, reinsTypeID)
	if err != nil {
		return models.KlausulTreaty{}, err
	}
	if len(baris) == 0 {
		return models.KlausulTreaty{}, fmt.Errorf("%w: induk jenis reasuransi %s", ErrKlausulTidakAda, reinsTypeID)
	}
	return baris[0], nil
}

// PctAnakLain membaca Pct anak LAIN satu induk, dikunci selama transaksi.
func (m *MasterKlausulTCO) PctAnakLain(ctx context.Context, tx *Tx, tahunID, descID, parentReinsTypeID, kecualiID string) (
	[]*apd.Decimal, error) {
	if tx == nil {
		return nil, errors.New("repository: membaca Pct anak menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelKlausulTCO)
	if err != nil {
		return nil, err
	}
	q := sqlPctAnakLainTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := tx.tx.QueryContext(ctx, q, tahunID, descID, parentReinsTypeID, kosongJadiNil(kecualiID), kosongJadiNil(kecualiID))
	if err != nil {
		return nil, fmt.Errorf("repository: membaca Pct anak: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []*apd.Decimal
	for rows.Next() {
		var id, pct sql.NullString
		if err := rows.Scan(&id, &pct); err != nil {
			return nil, err
		}
		d, err := desimalWarisanTCO(pct, "PCT")
		if err != nil {
			return nil, err
		}
		hasil = append(hasil, d)
	}
	return hasil, rows.Err()
}

// CariDobel mencari baris LAIN pada lingkup yang sama dengan kunci jenis sama.
func (m *MasterKlausulTCO) CariDobel(ctx context.Context, tx *Tx, k models.KlausulTreaty, kunci []string) (string, error) {
	if tx == nil {
		return "", errors.New("repository: pencarian dobel klausul menuntut transaksi")
	}
	if len(kunci) == 0 {
		return "", nil
	}
	tabel, err := m.db.Qualify(TabelKlausulTCO)
	if err != nil {
		return "", err
	}
	q := sqlCariDobelKlausulTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	rows, err := tx.tx.QueryContext(ctx, q, k.TreatyYearID, k.TreatyDescID, k.ParentReinsTypeID,
		kosongJadiNil(k.ID), kosongJadiNil(k.ID))
	if err != nil {
		return "", fmt.Errorf("repository: mencari klausul dobel: %w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		lain, err := pindaiKlausulTCO(rows)
		if err != nil {
			return "", err
		}
		sama := true
		for _, medan := range kunci {
			ok, err := medanSamaKlausulTCO(k, lain, medan)
			if err != nil {
				return "", err
			}
			sama = sama && ok
		}
		if sama {
			return lain.ID, nil
		}
	}
	return "", rows.Err()
}

// Sisip menulis klausul baru; ID '1' + 7 digit dari `PROPORTIONALARRG_SEQ` warisan.
func (m *MasterKlausulTCO) Sisip(ctx context.Context, tx *Tx, k models.KlausulTreaty) (string, error) {
	if tx == nil {
		return "", errors.New("repository: menyisipkan klausul menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelKlausulTCO)
	if err != nil {
		return "", err
	}
	if k.ID, err = m.db.IdentitasBerikutTCO(ctx, tx, SeqKlausulTCO); err != nil {
		return "", err
	}
	q := sqlSisipKlausulTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	v := nilaiKolomKlausul(k)
	arg := make([]any, 0, len(kolomKlausulTCO))
	for _, kolom := range kolomKlausulTCO {
		arg = append(arg, v[kolom])
	}
	hasil, err := tx.tx.ExecContext(ctx, q, arg...)
	if err != nil {
		return "", fmt.Errorf("repository: menyisipkan klausul: %w", err)
	}
	return k.ID, pastikanSatuBaris(hasil, "penyisipan klausul")
}

// Perbarui menimpa seluruh kolom non-kunci klausul milik tahun itu.
func (m *MasterKlausulTCO) Perbarui(ctx context.Context, tx *Tx, k models.KlausulTreaty) error {
	if tx == nil {
		return errors.New("repository: memperbarui klausul menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelKlausulTCO)
	if err != nil {
		return err
	}
	q := sqlPerbaruiKlausulTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	v := nilaiKolomKlausul(k)
	var arg []any
	for _, kolom := range kolomKlausulTCO {
		if !kunciTetapKlausul[kolom] {
			arg = append(arg, v[kolom])
		}
	}
	arg = append(arg, k.ID, k.TreatyYearID)
	hasil, err := tx.tx.ExecContext(ctx, q, arg...)
	if err != nil {
		return fmt.Errorf("repository: memperbarui klausul %s: %w", k.ID, err)
	}
	if n, err := hasil.RowsAffected(); err == nil && n == 0 {
		return ErrKlausulTidakAda
	}
	return pastikanSatuBaris(hasil, "pembaruan klausul")
}

func sqlKunciTahunTCO(tabel string) string {
	return fmt.Sprintf(`SELECT ID FROM %s WHERE ID = :1 FOR UPDATE`, tabel)
}

// Kunci mengunci baris tahun treaty selama penulis klausulnya bekerja - total
// Pct anak dan pencarian dobel tidak boleh dilewati penulis serentak.
func (m *MasterTahunTreaty) Kunci(ctx context.Context, tx *Tx, id string) error {
	if tx == nil {
		return errors.New("repository: mengunci tahun treaty menuntut transaksi")
	}
	tabel, err := m.db.Qualify(TabelTahunTCO)
	if err != nil {
		return err
	}
	q := sqlKunciTahunTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	var got string
	err = tx.tx.QueryRowContext(ctx, q, id).Scan(&got)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrTahunTreatyTidakAda
	}
	if err != nil {
		return fmt.Errorf("repository: mengunci tahun treaty %s: %w", id, err)
	}
	return nil
}
