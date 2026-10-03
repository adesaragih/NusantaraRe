package repository

// Pembaca ketujuh master pendukung (paket 2, tiket 04) - DIBACA SAJA.
//
// Setiap pemilih: section `*_Section` → RD (PARITAS §4) berparam `CARI1`
// (`SearchPolicyHolder.CARI1`, dihurufbesarkan `SearchPolicyHolder_act` 1 b236).
// Objek fisik = nama kelas `ASM-FW-GISFW-Int-<X>` - TERBUKTI di katalog DEV
// `ALL_OBJECTS` 01-10-2026 (OQ-MPNL-04 ditutup, lanjutan 1 L3): tabel `AGENT`,
// `CLIENT`; view `CURRENCY`, `CAUSEOFLOSS_LIFE`, `PRODUCT_TYPE_LIFE`,
// `RIRISK_LIFE_SUMMARY` (uji `TestObjekMasterAdaDiKatalogDEV`). Objek yang tidak
// terbaca di skema yang dikonfigurasi dijawab 503 yang MENYEBUT objeknya.
//
// ⛔ R/I Rate (`BrowseRateLifeSummary`, kelas `RATE_LIFE_SUMMARY`) dan `View Rate`
// (`BrowseRateLife_RD`, kelas `M_RATE_LIFE` = view `RATE_LIFE`): K1 keputusan work owner 01-10-2026
// (OQ-MPNL-03) - kedua view rate dibaca SAJA, kolom RD saja, nol `SELECT *`, nol
// `JSONDATA`, nol tulisan (`periksaBacaSaja`).
// ⛔ "Contains" Pega = `LIKE '%…%'`; kata cari kosong = semua baris. Batas
// baris = `pyMaxRecords` RD.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/modul/masterproductnamelife/backend/models"
)

// Objek master yang dibaca (kelas `ASM-FW-GISFW-Int-<X>`).
const (
	MasterAgent     = "AGENT"               // BrowseCedingCoLife_RD
	MasterClient    = "CLIENT"              // BrowseClientNusaRe_RD
	MasterCurrency  = "CURRENCY"            // BrowseCurrencyLIFE_RD / BrowseCurrencyFacIn_RD
	MasterRIRisk    = "RIRISK_LIFE_SUMMARY" // BrowseRIRiskSummary
	MasterCause     = "CAUSEOFLOSS_LIFE"    // BrowseCauseofLossLife_RD
	MasterJenisPlan = "PRODUCT_TYPE_LIFE"   // BrowseProductTypeLife_RD (paket 6)
	MasterRIRate    = "RATE_LIFE_SUMMARY"   // BrowseRateLifeSummary b40 (pemilih `Choose R/I Rate`)
	// MasterRate - view `RATE_LIFE` atas `M_RATE_LIFE.JSONDATA` (dialog `View Rate`, `BrowseRateLife_RD`
	// b39); nama fisik kelasnya: `NB FacIn/RDBList/BrowseLifeRate_SQL.xml` b85 `… FROM RATE_LIFE …`.
	MasterRate = "RATE_LIFE"
)

// DaftarMasterDibacaSaja - objek yang dibaca tetapi tidak pernah ditulis.
var DaftarMasterDibacaSaja = []string{MasterAgent, MasterClient, MasterCurrency, MasterRIRisk, MasterCause, MasterJenisPlan,
	MasterKontrakTreaty, MasterTahunTreaty, MasterRIRate, MasterRate}

// Nilai saringan VERBATIM RD.
const (
	saringIDCeding  = "%L0%" // `BrowseCedingCoLife_RD` b570 `.ID Contains "L0"`
	saringAktif     = "1"    // b607 `.StatusActive = 1`
	saringNamaStrip = "-"    // `BrowseClientNusaRe_RD` b570 `.Name != "-"`
	saringBukanITL  = "ITL"  // `BrowseCurrencyLIFE_RD` b541 `.Currency != "ITL"`
)

// sumber - satu master: SQL cari, SQL ambil-satu, dan argumennya.
type sumber struct {
	objek    string
	sqlCari  func(tabel string) string
	argCari  func(pola string) []any
	sqlAmbil func(tabel string) string
	argAmbil func(id string) []any
	kolom    []string
}

var sumberAgent = sumber{
	objek: MasterAgent,
	sqlCari: func(t string) string {
		return fmt.Sprintf(`SELECT ID, CLIENTNAME FROM %s
			WHERE ID LIKE :1 AND UPPER(CLIENTNAME) LIKE :2 ESCAPE '\' AND STATUSACTIVE = :3
			ORDER BY CLIENTNAME ASC FETCH FIRST 10000 ROWS ONLY`, t)
	},
	argCari: func(pola string) []any { return []any{saringIDCeding, pola, saringAktif} },
	sqlAmbil: func(t string) string {
		return fmt.Sprintf(`SELECT ID, CLIENTNAME FROM %s WHERE ID = :1 AND ID LIKE :2 AND STATUSACTIVE = :3`, t)
	},
	argAmbil: func(id string) []any { return []any{id, saringIDCeding, saringAktif} },
	kolom:    []string{"ID", "CLIENTNAME", "STATUSACTIVE"},
}

// sumberMaster - pemilih → sumbernya.
var sumberMaster = map[models.JenisMaster]sumber{
	models.MasterCeding: sumberAgent,
	models.MasterSOB:    sumberAgent,
	models.MasterPemegangPolis: {
		objek: MasterClient,
		sqlCari: func(t string) string {
			return fmt.Sprintf(`SELECT ID, NAME FROM %s
				WHERE UPPER(NAME) LIKE :1 ESCAPE '\' AND NAME <> :2 AND NAME IS NOT NULL
				ORDER BY NAME ASC, BU_NOTE ASC FETCH FIRST 100000 ROWS ONLY`, t)
		},
		argCari: func(pola string) []any { return []any{pola, saringNamaStrip} },
		sqlAmbil: func(t string) string {
			return fmt.Sprintf(`SELECT ID, NAME FROM %s WHERE ID = :1 AND NAME <> :2 AND NAME IS NOT NULL`, t)
		},
		argAmbil: func(id string) []any { return []any{id, saringNamaStrip} },
		kolom:    []string{"ID", "NAME", "BU_NOTE"},
	},
	models.MasterMataUang: {
		objek: MasterCurrency,
		sqlCari: func(t string) string {
			return fmt.Sprintf(`SELECT ID, CURRENCY FROM %s
				WHERE CURRENCY <> :1 AND UPPER(CURRENCY) LIKE :2 ESCAPE '\'
				ORDER BY CURRENCY ASC FETCH FIRST 500 ROWS ONLY`, t)
		},
		argCari: func(pola string) []any { return []any{saringBukanITL, pola} },
		sqlAmbil: func(t string) string {
			return fmt.Sprintf(`SELECT ID, CURRENCY FROM %s WHERE ID = :1 AND CURRENCY <> :2`, t)
		},
		argAmbil: func(id string) []any { return []any{id, saringBukanITL} },
		kolom:    []string{"ID", "CURRENCY"},
	},
	models.MasterRIRisk: {
		objek: MasterRIRisk,
		sqlCari: func(t string) string {
			return fmt.Sprintf(`SELECT ID, USEDBY FROM %s
				WHERE UPPER(USEDBY) LIKE :1 ESCAPE '\'
				ORDER BY ID ASC FETCH FIRST 500 ROWS ONLY`, t)
		},
		argCari:  func(pola string) []any { return []any{pola} },
		sqlAmbil: func(t string) string { return fmt.Sprintf(`SELECT ID, USEDBY FROM %s WHERE ID = :1`, t) },
		argAmbil: func(id string) []any { return []any{id} },
		kolom:    []string{"ID", "USEDBY"},
	},
	// `RIRate_Section` grid RD `BrowseRateLifeSummary` b2707: `id` kosong b1461 → saringan `.ID = param.id`
	// gugur; `idusedby` = `SearchPolicyHolder.CARI1` b1466 → `.USEDBY Contains` b807; urut `.ID ASC` b694;
	// `pyMaxRecords` 500 b674. `Choose` → `SetRIRate` b2448 (`.RIRATEID ← .ID`, `.RIRATE ← .USEDBY`).
	models.MasterRIRate: {
		objek: MasterRIRate,
		sqlCari: func(t string) string {
			return fmt.Sprintf(`SELECT ID, USEDBY FROM %s
				WHERE UPPER(USEDBY) LIKE :1 ESCAPE '\'
				ORDER BY ID ASC FETCH FIRST 500 ROWS ONLY`, t)
		},
		argCari:  func(pola string) []any { return []any{pola} },
		sqlAmbil: func(t string) string { return fmt.Sprintf(`SELECT ID, USEDBY FROM %s WHERE ID = :1`, t) },
		argAmbil: func(id string) []any { return []any{id} },
		kolom:    []string{"ID", "USEDBY"},
	},
	models.MasterPenyebab: {
		objek: MasterCause,
		sqlCari: func(t string) string {
			return fmt.Sprintf(`SELECT ID, CAUSEOFLOSS FROM %s
				WHERE UPPER(CAUSEOFLOSS) LIKE :1 ESCAPE '\'
				ORDER BY ID ASC FETCH FIRST 500 ROWS ONLY`, t)
		},
		argCari:  func(pola string) []any { return []any{pola} },
		sqlAmbil: func(t string) string { return fmt.Sprintf(`SELECT ID, CAUSEOFLOSS FROM %s WHERE ID = :1`, t) },
		argAmbil: func(id string) []any { return []any{id} },
		kolom:    []string{"ID", "CAUSEOFLOSS"},
	},
}

// KolomMaster - kolom setiap master yang disebut SQL (penjaga kata cadangan).
func KolomMaster() [][]string {
	hasil := [][]string{KolomJenisPlan}
	for _, s := range sumberMaster {
		hasil = append(hasil, s.kolom)
	}
	return hasil
}

var (
	// ErrMasterTidakTerbaca - master rujukan tidak dapat dibaca (503, pesan menyebut objeknya).
	ErrMasterTidakTerbaca = errors.New("repository: master data cannot be read")
	// ErrJenisMasterTidakDikenal - segmen jalur bukan salah satu pemilih (404).
	ErrJenisMasterTidakDikenal = errors.New("repository: unknown master picker")
)

// GalatMaster - master tidak terbaca; sebab Oracle hanya di log.
type GalatMaster struct {
	Objek string
	Sebab error
}

func (g GalatMaster) Error() string {
	return fmt.Sprintf("repository: master data %s cannot be read: %v", g.Objek, g.Sebab)
}

// PesanLayar - kalimat layar tanpa sebab Oracle.
func (g GalatMaster) PesanLayar() string {
	return fmt.Sprintf("master data %s cannot be read; check that the object exists in the configured schema", g.Objek)
}

// Is membuat errors.Is(err, ErrMasterTidakTerbaca) benar.
func (GalatMaster) Is(target error) bool { return target == ErrMasterTidakTerbaca }

func (g GalatMaster) Unwrap() error { return g.Sebab }

func galatMaster(objek string, sebab error) error { return GalatMaster{Objek: objek, Sebab: sebab} }

// PolaCari - `Contains` Pega: huruf besar, wildcard LIKE diloloskan.
func PolaCari(kata string) string {
	k := strings.ToUpper(strings.TrimSpace(kata))
	k = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(k)
	return "%" + k + "%"
}

// bacaMaster membaca hasil RD; `batas` > 0 = berhenti sesudah sekian baris (autocomplete).
func (g *Gudang) bacaMaster(ctx context.Context, objek, q string, batas int, args ...any) ([]models.NilaiMaster, error) {
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, galatMaster(objek, err)
	}
	defer func() { _ = rows.Close() }()
	hasil := []models.NilaiMaster{}
	for (batas <= 0 || len(hasil) < batas) && rows.Next() {
		var id, nama sql.NullString
		if err := rows.Scan(&id, &nama); err != nil {
			return nil, galatMaster(objek, err)
		}
		hasil = append(hasil, models.NilaiMaster{ID: id.String, Nama: nama.String})
	}
	if err := rows.Err(); err != nil {
		return nil, galatMaster(objek, err)
	}
	return hasil, nil
}

// CariMaster - grid section pemilih (`batas` 0 = `pyMaxRecords` RD) / autocomplete medan form (`batas` kecil).
func (g *Gudang) CariMaster(ctx context.Context, jenis models.JenisMaster, kata string, batas int) ([]models.NilaiMaster, error) {
	s, ada := sumberMaster[jenis]
	if !ada {
		return nil, fmt.Errorf("%w: %q", ErrJenisMasterTidakDikenal, jenis)
	}
	q, err := g.siapkan(s.objek, s.sqlCari)
	if err != nil {
		return nil, err
	}
	return g.bacaMaster(ctx, s.objek, q, batas, s.argCari(PolaCari(kata))...)
}

// AmbilMaster - satu nilai master menurut ID dan saringan RD yang sama
// (verifikasi pilihan saat simpan: "pilihan yang tidak ada di master ditolak").
func (g *Gudang) AmbilMaster(ctx context.Context, jenis models.JenisMaster, id string) (models.NilaiMaster, bool, error) {
	s, ada := sumberMaster[jenis]
	if !ada {
		return models.NilaiMaster{}, false, fmt.Errorf("%w: %q", ErrJenisMasterTidakDikenal, jenis)
	}
	q, err := g.siapkan(s.objek, s.sqlAmbil)
	if err != nil {
		return models.NilaiMaster{}, false, err
	}
	d, err := g.bacaMaster(ctx, s.objek, q, 1, s.argAmbil(id)...)
	if err != nil || len(d) == 0 {
		return models.NilaiMaster{}, false, err
	}
	return d[0], true, nil
}

// KolomJenisPlan - kolom `PRODUCT_TYPE_LIFE` yang dibaca (`BrowseProductTypeLife_RD` b682–b712).
var KolomJenisPlan = []string{"ID", "COVERNAME", "BUSINESS", "BENEFIT"}

// sqlCariPlan - autocomplete `Plan Name`: dicari pada `.CoverName` dan `.Business`
// (`pyUseForSearch` true), tanpa saringan RD (b531), maks 500 b664. RD tidak
// mengurutkan (`pySortOrder` 99999) - di sini `ID ASC` supaya tetap.
func sqlCariPlan(tabel string) string {
	return fmt.Sprintf(`SELECT ID, COVERNAME, BUSINESS, BENEFIT FROM %s
		WHERE UPPER(COVERNAME) LIKE :1 ESCAPE '\' OR UPPER(BUSINESS) LIKE :2 ESCAPE '\'
		ORDER BY ID ASC FETCH FIRST 500 ROWS ONLY`, tabel)
}

func sqlAmbilPlan(tabel string) string {
	return fmt.Sprintf(`SELECT ID, COVERNAME, BUSINESS, BENEFIT FROM %s WHERE ID = :1`, tabel)
}

func (g *Gudang) bacaPlan(ctx context.Context, q string, args ...any) ([]models.JenisPlan, error) {
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, galatMaster(MasterJenisPlan, err)
	}
	defer func() { _ = rows.Close() }()
	hasil := []models.JenisPlan{}
	for rows.Next() {
		var id, cover, biz, benefit sql.NullString
		if err := rows.Scan(&id, &cover, &biz, &benefit); err != nil {
			return nil, galatMaster(MasterJenisPlan, err)
		}
		hasil = append(hasil, models.JenisPlan{ID: id.String, CoverName: cover.String, Business: biz.String, Benefit: benefit.String})
	}
	if err := rows.Err(); err != nil {
		return nil, galatMaster(MasterJenisPlan, err)
	}
	return hasil, nil
}

// CariPlan - autocomplete `Plan Name` grid `PLAN LIST`.
func (g *Gudang) CariPlan(ctx context.Context, kata string) ([]models.JenisPlan, error) {
	q, err := g.siapkan(MasterJenisPlan, sqlCariPlan)
	if err != nil {
		return nil, err
	}
	pola := PolaCari(kata)
	return g.bacaPlan(ctx, q, pola, pola)
}

// AmbilPlan - satu jenis plan menurut ID (verifikasi saat simpan).
func (g *Gudang) AmbilPlan(ctx context.Context, id string) (models.JenisPlan, bool, error) {
	q, err := g.siapkan(MasterJenisPlan, sqlAmbilPlan)
	if err != nil {
		return models.JenisPlan{}, false, err
	}
	d, err := g.bacaPlan(ctx, q, id)
	if err != nil || len(d) == 0 {
		return models.JenisPlan{}, false, err
	}
	return d[0], true, nil
}

// BatasRate - `BrowseRateLife_RD` b730 `pyMaxRecords` 500; satu baris lebih dibaca untuk mengetahui
// dialog terpotong (Pega memotong diam-diam).
const BatasRate = 500

// KolomRate - kolom `RATE_LIFE` yang dibaca: enam kolom grid `ViewRate` (subset kolom RD b747–b791).
var KolomRate = []string{"ID", "USEDBY", "GENDER", "CONTRACT", "AGE", "RATE"}

// sqlDaftarRate - `BrowseRateLife_RD`: `.IDUSEDBY = Param.idusedby` b859/b868, urut `.ID DESC` b748,
// `.RATE ASC` b786. Berkunci `IDUSEDBY` (view atas CLOB tanpa index - nol pembacaan tanpa kunci).
func sqlDaftarRate(t string) string {
	return fmt.Sprintf(`SELECT ID, USEDBY, GENDER, CONTRACT, AGE, RATE FROM %s
		WHERE IDUSEDBY = :1
		ORDER BY ID DESC, RATE ASC FETCH FIRST %d ROWS ONLY`, t, BatasRate+1)
}

// DaftarRate - dialog `View Rate` satu R/I Rate (`IDUSEDBY` = `RIRATEID` baris plan), paling banyak
// BatasRate baris; `terpotong` benar bila view memuat lebih.
func (g *Gudang) DaftarRate(ctx context.Context, idUsedBy string) ([]models.BarisRate, bool, error) {
	q, err := g.siapkan(MasterRate, sqlDaftarRate)
	if err != nil {
		return nil, false, err
	}
	rows, err := g.db.QueryContext(ctx, q, idUsedBy)
	if err != nil {
		return nil, false, galatMaster(MasterRate, err)
	}
	defer func() { _ = rows.Close() }()
	hasil := []models.BarisRate{}
	for rows.Next() {
		var n [6]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5]); err != nil {
			return nil, false, galatMaster(MasterRate, err)
		}
		hasil = append(hasil, models.BarisRate{ID: n[0].String, UsedBy: n[1].String, Gender: n[2].String,
			Contract: n[3].String, Age: n[4].String, Rate: n[5].String})
	}
	if err := rows.Err(); err != nil {
		return nil, false, galatMaster(MasterRate, err)
	}
	terpotong := len(hasil) > BatasRate
	if terpotong {
		hasil = hasil[:BatasRate]
	}
	return hasil, terpotong, nil
}
