package repository

// Master Marine Cargo (tiket 74 P1 / P6; DDL dikirim work owner 11-10-2026, `D:\migrasi\RNM\DDL\`): view POOLDATA.TRADING /
// GOODSTYPE / PACKING (ID, DESCRIPTION, TYPE) dan CONVEYANCE (ID, DESCRIPTION, TYPECONVEYANCE, TYPE) atas JSONDATA m_trading /
// m_goodstype / m_packing / m_conveyance `[terverifikasi DDL]`. Baca saja; skema lewat db.Qualify (pola SHIP / OCCUPATION).
//
// `[terverifikasi]` RD NB FacIn BrowseTrading_RD / BrowseGoodsType / BrowsePacking_RD / BrowseConveyance_RD: kelas Int-TRADING /
// GOODSTYPE / PACKING / CONVEYANCE, filter HANYA `.ID = Param.ID` dan `.Description = Param.Description` (TANPA saringan TYPE;
// filter kosong diabaikan), tanpa urutan (pySortOrder 99999), pyMaxRecords 500.
//   - Deskripsi (catatan ProtectObjectMarine, PUT / GET objek-cargo): `ID =`.
//   - IDMaster (unggah CSV, ProtectUploadMarineCargo_act pxResults(1)): `DESCRIPTION =` persis; RD tanpa urutan - bila deskripsi
//     ganda, ID terkecil `[keputusan agent]` (deterministik).
//   - Daftar (dropdown form cargo): tanpa filter, urut DESCRIPTION lalu ID (RD tanpa urutan - arahan sesi 0f), paling banyak 500.
// DDL kedua (work owner 11-10-2026): view POOLDATA.LLOYDAGENT (ID, NAME, ADDRESS, CITY, CONTACT, COUNTRY, TYPEAGENTID, TYPE),
// MAXAGE (ID, WARRANTY, TYPE), CONDITION (ID, CONDITIONNAME, BUSINESSCODE, BUSINESSNAME, CONDITIONADD, TYPE; TYPE 'MARINE' ->
// CONDITIONNAME = JSON Description) `[terverifikasi DDL]`.
//   - Agen (CSV P6, ProtectUploadMarineCargo_act 5.2 BrowseLLOYDAGENT_RD Param.ID): `ID =`.
//   - CariAgen (popup Select Agent, Harness / Section SelectAgent -> BrowseLLOYDAGENT_RD NB): filter A `.ID = Param.ID` (kosong =
//     diabaikan), B `.Name Contains Param.Name` TANPA pyCaseInsensitive -> peka huruf `[terverifikasi]`; tanpa urutan -> NAME, ID;
//     maks 500.
//   - WarrantyMaxAge (PUT objek-cargo MaxAgeDetail, MarineCargoDtl autocomplete .PolicyData.MaxAgeDetail = .Warranty): `ID =`.
//   - CariMaxAge (autocomplete MarineCargoDtl L14861: BrowseMaxAge_RD NB tanpa filter, kolom pencarian .ID dan .Warranty; mode
//     cocok autocomplete tidak tertulis di ekspor -> "memuat", tak peka huruf `[keputusan agent]`); urut ID; maks 500.
//   - CariKondisiMarine (Select Coverage / Warranty Extension, BrowseMarineCondition_RD NB = kelas Int-MARINECONDITION, field .ID /
//     .Description / .Type / .ConditionAdd, TANPA filter `[terverifikasi]`). Varian Endorsment Fac In (Int-CONDITION, filter Type
//     "MARINE" + ConditionName / ConditionAdd / ID) identitas lain - TIDAK dipinjam (CLAUDE.md 6). Sumber = view CONDITION dari DDL
//     "MARINE CONDITION.txt" kiriman work owner; Int-MARINECONDITION = CONDITION ber-TYPE 'MARINE' (keputusan work owner
//     11-10-2026; kelas -> tabel tidak ada di ekspor) - TANPA saringan CONDITIONADD (RD NB tidak menyaring). q: memuat, tak peka huruf, pada
//     CONDITIONNAME `[keputusan agent]`; urut CONDITIONNAME, ID; maks 500.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// View master Marine Cargo (ejaan DDL).
	ViewLloydAgent = "LLOYDAGENT"
	ViewMaxAge     = "MAXAGE"
	ViewCondition  = "CONDITION"
	// tipeKondisiMarine - CONDITION.TYPE baris Marine (CASE view).
	tipeKondisiMarine = "MARINE"
	ViewTrading       = "TRADING"
	ViewGoodsType     = "GOODSTYPE"
	ViewConveyance    = "CONVEYANCE"
	ViewPacking       = "PACKING"
	// BatasMasterCargo - pyMaxRecords RD master Marine.
	BatasMasterCargo = 500
)

// ErrMasterBelumAda - master yang DDL-nya belum dikirim (MaxAge, LLOYDAGENT).
var ErrMasterBelumAda = errors.New("repository: master belum tersedia (DDL menunggu)")

// ErrJenisMasterTakDikenal - jenis master selain Trading / Goods / Conveyance / Packing.
var ErrJenisMasterTakDikenal = errors.New("repository: jenis master Marine tidak dikenal")

// viewMasterCargo - jenis (services.MasterTrading dst.) -> view.
var viewMasterCargo = map[string]string{"Trading": ViewTrading, "Goods": ViewGoodsType, "Conveyance": ViewConveyance,
	"Packing": ViewPacking}

// MasterCargoOracle - master Marine Cargo atas Oracle (services.PembacaMasterCargo).
type MasterCargoOracle struct{ db *db.DB }

// NewMasterCargoOracle merakit pembaca master Marine Cargo.
func NewMasterCargoOracle(d *db.DB) *MasterCargoOracle { return &MasterCargoOracle{db: d} }

func (r *MasterCargoOracle) view(jenis string) (string, error) {
	v, ok := viewMasterCargo[jenis]
	if !ok {
		return "", fmt.Errorf("%w: %q", ErrJenisMasterTakDikenal, jenis)
	}
	return r.db.Qualify(v)
}

func sqlDeskripsiMaster(v string) string {
	return "SELECT DESCRIPTION FROM " + v + " WHERE ID = :1 FETCH FIRST 1 ROWS ONLY"
}

func sqlIDMaster(v string) string {
	return "SELECT ID FROM " + v + " WHERE DESCRIPTION = :1 ORDER BY ID FETCH FIRST 1 ROWS ONLY"
}

func sqlDaftarMaster(v string) string {
	return "SELECT ID, DESCRIPTION FROM " + v + " ORDER BY DESCRIPTION, ID FETCH FIRST :1 ROWS ONLY"
}

// satuTeks - satu kolom teks satu baris; tidak ada = ada false.
func (r *MasterCargoOracle) satuTeks(ctx context.Context, q, nama, nilai string) (string, bool, error) {
	var s sql.NullString
	err := r.db.QueryRowContext(ctx, q, nilai).Scan(&s)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("repository: baca %s: %w", nama, err)
	}
	return s.String, true, nil
}

// Deskripsi - DESCRIPTION master `jenis` ber-ID `id`.
func (r *MasterCargoOracle) Deskripsi(ctx context.Context, jenis, id string) (string, bool, error) {
	v, err := r.view(jenis)
	if err != nil {
		return "", false, err
	}
	return r.satuTeks(ctx, sqlDeskripsiMaster(v), viewMasterCargo[jenis], id)
}

// IDMaster - ID master `jenis` ber-DESCRIPTION persis `deskripsi`.
func (r *MasterCargoOracle) IDMaster(ctx context.Context, jenis, deskripsi string) (string, bool, error) {
	v, err := r.view(jenis)
	if err != nil {
		return "", false, err
	}
	return r.satuTeks(ctx, sqlIDMaster(v), viewMasterCargo[jenis], deskripsi)
}

// Daftar - pilihan dropdown master `jenis`.
func (r *MasterCargoOracle) Daftar(ctx context.Context, jenis string) ([]models.BarisMasterCargo, error) {
	v, err := r.view(jenis)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlDaftarMaster(v), BatasMasterCargo)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", viewMasterCargo[jenis], err)
	}
	defer baris.Close()
	hasil := []models.BarisMasterCargo{}
	for baris.Next() {
		var id, d sql.NullString
		if err := baris.Scan(&id, &d); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", viewMasterCargo[jenis], err)
		}
		hasil = append(hasil, models.BarisMasterCargo{ID: id.String, Description: d.String})
	}
	return hasil, baris.Err()
}

// polaPekaHuruf - pola LIKE "memuat" tanpa mengubah huruf (`\ % _` di-escape); kosong = "".
func polaPekaHuruf(kata string) string {
	k := strings.TrimSpace(kata)
	if k == "" {
		return ""
	}
	return "%" + strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(k) + "%"
}

func sqlWarrantyMaxAge(v string) string {
	return "SELECT WARRANTY FROM " + v + " WHERE ID = :1 FETCH FIRST 1 ROWS ONLY"
}

// sqlCariMaxAge - [:1 dan :2 pola huruf besar], batas.
func sqlCariMaxAge(v string, cari bool) string {
	syarat, batas := "", ":1"
	if cari {
		syarat, batas = ` WHERE (UPPER(ID) LIKE :1 ESCAPE '\' OR UPPER(WARRANTY) LIKE :2 ESCAPE '\')`, ":3"
	}
	return "SELECT ID, WARRANTY FROM " + v + syarat + " ORDER BY ID FETCH FIRST " + batas + " ROWS ONLY"
}

const kolomAgenMaster = "ID, NAME, ADDRESS, CITY, COUNTRY, CONTACT"

func sqlAgen(v string) string {
	return "SELECT " + kolomAgenMaster + " FROM " + v + " WHERE ID = :1 FETCH FIRST 1 ROWS ONLY"
}

// sqlCariAgen - [:1 pola peka huruf], batas.
func sqlCariAgen(v string, cari bool) string {
	syarat, batas := "", ":1"
	if cari {
		syarat, batas = ` WHERE NAME LIKE :1 ESCAPE '\'`, ":2"
	}
	return "SELECT " + kolomAgenMaster + " FROM " + v + syarat + " ORDER BY NAME, ID FETCH FIRST " + batas + " ROWS ONLY"
}

// sqlCariKondisi - :1 TYPE, [:2 pola huruf besar], batas.
func sqlCariKondisi(v string, cari bool) string {
	syarat, batas := "", ":2"
	if cari {
		syarat, batas = ` AND UPPER(CONDITIONNAME) LIKE :2 ESCAPE '\'`, ":3"
	}
	return "SELECT ID, CONDITIONNAME FROM " + v + " WHERE TYPE = :1" + syarat + " ORDER BY CONDITIONNAME, ID FETCH FIRST " + batas +
		" ROWS ONLY"
}

// WarrantyMaxAge - WARRANTY MaxAge ber-ID `id`.
func (r *MasterCargoOracle) WarrantyMaxAge(ctx context.Context, id string) (string, bool, error) {
	v, err := r.db.Qualify(ViewMaxAge)
	if err != nil {
		return "", false, err
	}
	return r.satuTeks(ctx, sqlWarrantyMaxAge(v), ViewMaxAge, id)
}

// bacaBaris - semua baris kueri sebagai teks (kolom = n).
func (r *MasterCargoOracle) bacaBaris(ctx context.Context, q, nama string, n int, arg ...any) ([][]string, error) {
	baris, err := r.db.QueryContext(ctx, q, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", nama, err)
	}
	defer baris.Close()
	var hasil [][]string
	for baris.Next() {
		v := make([]sql.NullString, n)
		tuju := make([]any, n)
		for i := range v {
			tuju[i] = &v[i]
		}
		if err := baris.Scan(tuju...); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", nama, err)
		}
		s := make([]string, n)
		for i := range v {
			s[i] = v[i].String
		}
		hasil = append(hasil, s)
	}
	return hasil, baris.Err()
}

func agenDari(s []string) models.BarisAgenSurvei {
	return models.BarisAgenSurvei{ID: s[0], Name: s[1], Address: s[2], City: s[3], Country: s[4], Contact: s[5]}
}

// Agen - agen survei ber-ID `id` (ASMAgentID / pyFullName / pyAddress / pyCity / pyCountry).
func (r *MasterCargoOracle) Agen(ctx context.Context, id string) (models.AgenSurvei, bool, error) {
	v, err := r.db.Qualify(ViewLloydAgent)
	if err != nil {
		return models.AgenSurvei{}, false, err
	}
	b, err := r.bacaBaris(ctx, sqlAgen(v), ViewLloydAgent, 6, id)
	if err != nil || len(b) == 0 {
		return models.AgenSurvei{}, false, err
	}
	a := agenDari(b[0])
	return models.AgenSurvei{AgentID: a.ID, FullName: a.Name, Address: a.Address, City: a.City, Country: a.Country}, true, nil
}

// CariAgen - popup Select Agent (lihat kepala berkas).
func (r *MasterCargoOracle) CariAgen(ctx context.Context, kata string) ([]models.BarisAgenSurvei, error) {
	v, err := r.db.Qualify(ViewLloydAgent)
	if err != nil {
		return nil, err
	}
	q, arg := sqlCariAgen(v, false), []any{BatasMasterCargo}
	if p := polaPekaHuruf(kata); p != "" {
		q, arg = sqlCariAgen(v, true), []any{p, BatasMasterCargo}
	}
	b, err := r.bacaBaris(ctx, q, ViewLloydAgent, 6, arg...)
	hasil := []models.BarisAgenSurvei{}
	for _, s := range b {
		hasil = append(hasil, agenDari(s))
	}
	return hasil, err
}

// CariMaxAge - autocomplete Max Age (lihat kepala berkas).
func (r *MasterCargoOracle) CariMaxAge(ctx context.Context, kata string) ([]models.BarisMaxAge, error) {
	v, err := r.db.Qualify(ViewMaxAge)
	if err != nil {
		return nil, err
	}
	q, arg := sqlCariMaxAge(v, false), []any{BatasMasterCargo}
	if p := PolaCari(kata); p != "" {
		q, arg = sqlCariMaxAge(v, true), []any{p, p, BatasMasterCargo}
	}
	b, err := r.bacaBaris(ctx, q, ViewMaxAge, 2, arg...)
	hasil := []models.BarisMaxAge{}
	for _, s := range b {
		hasil = append(hasil, models.BarisMaxAge{ID: s[0], Warranty: s[1]})
	}
	return hasil, err
}

// CariKondisiMarine - Select Coverage / Warranty Extension (lihat kepala berkas).
func (r *MasterCargoOracle) CariKondisiMarine(ctx context.Context, kata string) ([]models.BarisKondisiMarine, error) {
	v, err := r.db.Qualify(ViewCondition)
	if err != nil {
		return nil, err
	}
	q, arg := sqlCariKondisi(v, false), []any{tipeKondisiMarine, BatasMasterCargo}
	if p := PolaCari(kata); p != "" {
		q, arg = sqlCariKondisi(v, true), []any{tipeKondisiMarine, p, BatasMasterCargo}
	}
	b, err := r.bacaBaris(ctx, q, ViewCondition, 2, arg...)
	hasil := []models.BarisKondisiMarine{}
	for _, s := range b {
		hasil = append(hasil, models.BarisKondisiMarine{ID: s[0], ConditionName: s[1]})
	}
	return hasil, err
}
