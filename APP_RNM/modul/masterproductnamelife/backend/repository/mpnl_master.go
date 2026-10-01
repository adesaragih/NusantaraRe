package repository

// Pembaca ketujuh master pendukung (paket 2, tiket 04) - DIBACA SAJA.
//
// Setiap pemilih: section `*_Section` → RD (PARITAS §4) berparam `CARI1`
// (`SearchPolicyHolder.CARI1`, dihurufbesarkan `SearchPolicyHolder_act` 1 b234).
// Tabel fisik = nama kelas `ASM-FW-GISFW-Int-<X>` (konvensi; preseden Retro
// Life `AGENT`) - tidak terbukti di katalog (OQ-MPNL-04): objek yang tidak ada
// dijawab 503 yang MENYEBUT objeknya, bukan daftar kosong.
//
// ⛔ R/I Rate (`BrowseRateLifeSummary`, kelas `RATE_LIFE_SUMMARY`) TIDAK
// dibaca: view atas JSON rate - membacanya menunggu persetujuan work owner
// (OQ-MPNL-03, preseden OQ-MCRL-13). Namanya sengaja tidak ada di sini.
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
)

// DaftarMasterDibacaSaja - objek yang dibaca tetapi tidak pernah ditulis.
var DaftarMasterDibacaSaja = []string{MasterAgent, MasterClient, MasterCurrency, MasterRIRisk, MasterCause, MasterJenisPlan,
	MasterKontrakTreaty, MasterTahunTreaty}

// Nilai saringan VERBATIM RD.
const (
	saringIDCeding  = "%L0%" // `BrowseCedingCoLife_RD` b565 `.ID Contains "L0"`
	saringAktif     = "1"    // b601 `.StatusActive = 1`
	saringNamaStrip = "-"    // `BrowseClientNusaRe_RD` b565 `.Name != "-"`
	saringBukanITL  = "ITL"  // `BrowseCurrencyLIFE_RD` b534 `.Currency != "ITL"`
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

// sumberMaster - pemilih → sumbernya. R/I Rate sengaja tidak ada (OQ-MPNL-03).
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
	return fmt.Sprintf("master data %s cannot be read; check that the object exists in the configured schema (OQ-MPNL-04)", g.Objek)
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

func (g *Gudang) bacaMaster(ctx context.Context, objek, q string, args ...any) ([]models.NilaiMaster, error) {
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, galatMaster(objek, err)
	}
	defer func() { _ = rows.Close() }()
	hasil := []models.NilaiMaster{}
	for rows.Next() {
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

// CariMaster - grid section pemilih / autocomplete medan form.
func (g *Gudang) CariMaster(ctx context.Context, jenis models.JenisMaster, kata string) ([]models.NilaiMaster, error) {
	s, ada := sumberMaster[jenis]
	if !ada {
		return nil, fmt.Errorf("%w: %q", ErrJenisMasterTidakDikenal, jenis)
	}
	q, err := g.siapkan(s.objek, s.sqlCari)
	if err != nil {
		return nil, err
	}
	return g.bacaMaster(ctx, s.objek, q, s.argCari(PolaCari(kata))...)
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
	d, err := g.bacaMaster(ctx, s.objek, q, s.argAmbil(id)...)
	if err != nil || len(d) == 0 {
		return models.NilaiMaster{}, false, err
	}
	return d[0], true, nil
}

// KolomJenisPlan - kolom `PRODUCT_TYPE_LIFE` yang dibaca (`BrowseProductTypeLife_RD` b680–b710).
var KolomJenisPlan = []string{"ID", "COVERNAME", "BUSINESS", "BENEFIT"}

// sqlCariPlan - autocomplete `Plan Name`: dicari pada `.CoverName` dan `.Business`
// (`pyUseForSearch` true), tanpa saringan RD (b530), maks 500 b664. RD tidak
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
