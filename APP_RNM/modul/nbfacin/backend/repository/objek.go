package repository

// Tab Object FIRE (tiket 35): .LocationList -> T_LOCATIONLIST (induk T_GENERAL_POLIS) ->
// T_PROPERTY -> T_RISKLOCATION / T_BUILDINGCONSTRUCTION (tabel rancangan, migrasi 186).
// Daftar dibaca urut SEQ_NO dan DIGANTI UTUH di dalam transaksi pemanggil (anak lebih dulu
// dihapus, lalu disisip ulang urut). Seluruh nilai parameter terikat.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	TabelLocationList         = "T_LOCATIONLIST"
	TabelProperty             = "T_PROPERTY"
	TabelRiskLocation         = "T_RISKLOCATION"
	TabelBuildingConstruction = "T_BUILDINGCONSTRUCTION"
	sequenceLocationList      = "SEQ_T_LOCATIONLIST"
	sequenceProperty          = "SEQ_T_PROPERTY"
	sequenceRiskLocation      = "SEQ_T_RISKLOCATION"
	sequenceBuildingConstruct = "SEQ_T_BUILDINGCONSTRUCTION"
	teksBenar, teksSalah      = "true", "false" // boolean Pega di IS_* (fixture)
)

// PenyimpanObjek - baca/ganti daftar objek case.
type PenyimpanObjek interface {
	// BacaObjek - ErrKasusTidakAda bila case tidak ada / bukan LINI Fac In.
	BacaObjek(ctx context.Context, id string) ([]models.ObjekFire, error)
	// GantiObjek - ganti seluruh daftar di dalam transaksi pemanggil.
	GantiObjek(ctx context.Context, tx *db.Tx, id string, baris []models.ObjekFire) error
}

// ObjekOracle - PenyimpanObjek atas Oracle.
type ObjekOracle struct{ db *db.DB }

// NewObjekOracle merakit penyimpan objek.
func NewObjekOracle(d *db.DB) *ObjekOracle { return &ObjekOracle{db: d} }

type tabelObjek struct{ work, general, loc, prop, risk, bang string }

func (r *ObjekOracle) tabel() (tabelObjek, error) {
	var t tabelObjek
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelWorkPolis, &t.work}, {TabelGeneralPolis, &t.general}, {TabelLocationList, &t.loc}, {TabelProperty, &t.prop},
		{TabelRiskLocation, &t.risk}, {TabelBuildingConstruction, &t.bang}} {
		q, err := r.db.Qualify(x.nama)
		if err != nil {
			return t, err
		}
		*x.ke = q
	}
	return t, nil
}

func sqlAdaKasus(work string) string {
	return "SELECT COUNT(*) FROM " + work + " WHERE ID = :1 AND LINI = :2"
}

// kolomBacaObjek - kolom sqlBacaObjek BERNAMA (dibaca lewat kunci, bukan urutan) -> medan.
var kolomBacaObjek = []string{"p.OBJECT_NO", "p.OBJECT_TYPE", "p.OBJECT_NAME", "p.IS_MATERIAL_DAMAGE", "p.IS_TOP_RISK",
	"p.ROAD_TYPE", "p.ROAD_NAME", "p.BUILDING_NO", "r.ASM_ZIP_CODE", "p.COUNTRY", "r.ASM_ADDRESS", "r.ASMRW", "r.ASM_CITY",
	"r.ASM_DISTRICT", "p.PROVINCE", "p.ALM_RISK_ID", "b.NUMBER_OF_FLOOR", "b.ROOF_TYPE", "b.WALL_TYPE", "b.FLOOR_TYPE",
	"b.PARTITION_TYPE", "b.SUPPORT_WALL_TYPE", "b.OTHERS_TYPE"}

// sqlBacaObjek - satu baris per lokasi; anak tunggal lewat LEFT JOIN (UNIQUE PARENT_ID, 186).
func sqlBacaObjek(t tabelObjek) string {
	return "SELECT " + strings.Join(kolomBacaObjek, ", ") + `
FROM ` + t.loc + ` l
LEFT JOIN ` + t.prop + ` p ON p.PARENT_ID = l.ID
LEFT JOIN ` + t.risk + ` r ON r.PARENT_ID = p.ID
LEFT JOIN ` + t.bang + ` b ON b.PARENT_ID = p.ID
WHERE l.PARENT_ID = :1
ORDER BY l.SEQ_NO`
}

// sqlHapusObjek - anak sebelum induk (FK tanpa ON DELETE).
func sqlHapusObjek(t tabelObjek) []string {
	prop := "SELECT p.ID FROM " + t.prop + " p JOIN " + t.loc + " l ON l.ID = p.PARENT_ID WHERE l.PARENT_ID = :1"
	return []string{
		"DELETE FROM " + t.risk + " WHERE PARENT_ID IN (" + prop + ")",
		"DELETE FROM " + t.bang + " WHERE PARENT_ID IN (" + prop + ")",
		"DELETE FROM " + t.prop + " WHERE PARENT_ID IN (SELECT ID FROM " + t.loc + " WHERE PARENT_ID = :1)",
		"DELETE FROM " + t.loc + " WHERE PARENT_ID = :1",
	}
}

// sqlPastikanGeneral - T_LOCATIONLIST ber-FK ke T_GENERAL_POLIS: baris General kosong dibuat
// bila General belum pernah disimpan.
func sqlPastikanGeneral(general string) string {
	return "INSERT INTO " + general + " (ID) SELECT :1 FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM " + general + " WHERE ID = :2)"
}

func sqlSisipLokasi(loc string) string {
	return "INSERT INTO " + loc + " (ID, PARENT_ID, SEQ_NO, ROW_UID) VALUES (:1, :2, :3, :4)"
}

func sqlSisipProperty(prop string) string {
	return "INSERT INTO " + prop + " (ID, PARENT_ID, OBJECT_NO, OBJECT_TYPE, OBJECT_NAME, IS_MATERIAL_DAMAGE, IS_TOP_RISK," +
		" ROAD_TYPE, ROAD_NAME, BUILDING_NO, COUNTRY, PROVINCE, ALM_RISK_ID) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13)"
}

func sqlSisipRisk(risk string) string {
	return "INSERT INTO " + risk + " (ID, PARENT_ID, ASM_ADDRESS, ASM_CITY, ASM_DISTRICT, ASMRW, ASM_ZIP_CODE) VALUES (:1, :2, :3, :4, :5, :6, :7)"
}

func sqlSisipBangunan(bang string) string {
	return "INSERT INTO " + bang + " (ID, PARENT_ID, FLOOR_TYPE, NUMBER_OF_FLOOR, ROOF_TYPE, WALL_TYPE, PARTITION_TYPE," +
		" SUPPORT_WALL_TYPE, OTHERS_TYPE) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9)"
}

// teksBool - boolean -> teks Pega. Baca: hanya teks "true" = benar (kosong/NULL = salah).
// `[dugaan]` IsMaterialDamage tak dicentang: fixture tidak pernah berisi "false" (kuncinya
// tidak ada), aplikasi tetap menulis "false" (A115).
func teksBool(b bool) string {
	if b {
		return teksBenar
	}
	return teksSalah
}

// BacaObjek - lihat PenyimpanObjek.
func (r *ObjekOracle) BacaObjek(ctx context.Context, id string) ([]models.ObjekFire, error) {
	t, err := r.tabel()
	if err != nil {
		return nil, err
	}
	var n int
	if err := r.db.QueryRowContext(ctx, sqlAdaKasus(t.work), id, LiniFacIn).Scan(&n); err != nil {
		return nil, fmt.Errorf("repository: cek case NB: %w", err)
	}
	if n == 0 {
		return nil, ErrKasusTidakAda
	}
	baris, err := r.db.QueryContext(ctx, sqlBacaObjek(t), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca objek: %w", err)
	}
	defer baris.Close()
	hasil := []models.ObjekFire{}
	for baris.Next() {
		teks := make(map[string]*sql.NullString, len(kolomBacaObjek))
		tujuan := make([]any, len(kolomBacaObjek))
		for i, k := range kolomBacaObjek {
			teks[k] = &sql.NullString{}
			tujuan[i] = teks[k]
		}
		if err := baris.Scan(tujuan...); err != nil {
			return nil, fmt.Errorf("repository: objek: %w", err)
		}
		v := func(k string) string { return teks[k].String }
		hasil = append(hasil, models.ObjekFire{ObjectNo: v("p.OBJECT_NO"), ObjectType: v("p.OBJECT_TYPE"),
			ObjectName: v("p.OBJECT_NAME"), IsMaterialDamage: v("p.IS_MATERIAL_DAMAGE") == teksBenar,
			IsTopRisk: v("p.IS_TOP_RISK") == teksBenar, RoadType: v("p.ROAD_TYPE"), RoadName: v("p.ROAD_NAME"),
			BuildingNo: v("p.BUILDING_NO"), ZipCode: v("r.ASM_ZIP_CODE"), Country: v("p.COUNTRY"),
			RiskLocation: v("r.ASM_ADDRESS"), Territory: v("r.ASMRW"), City: v("r.ASM_CITY"), District: v("r.ASM_DISTRICT"),
			Province: v("p.PROVINCE"), RiskAddressID: v("p.ALM_RISK_ID"), NumberOfFloor: v("b.NUMBER_OF_FLOOR"),
			RoofType: v("b.ROOF_TYPE"), WallType: v("b.WALL_TYPE"), FloorType: v("b.FLOOR_TYPE"),
			PartitionType: v("b.PARTITION_TYPE"), SupportWallType: v("b.SUPPORT_WALL_TYPE"), OtherType: v("b.OTHERS_TYPE")})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: objek: %w", err)
	}
	return hasil, nil
}

// GantiObjek - lihat PenyimpanObjek. Urutan: sentuh case (kunci + 404), pastikan General,
// hapus anak -> induk, sisip ulang urut (SEQ_NO 1..n).
func (r *ObjekOracle) GantiObjek(ctx context.Context, tx *db.Tx, id string, baris []models.ObjekFire) error {
	t, err := r.tabel()
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, sqlSentuhCase(t.work), "menyentuh "+TabelWorkPolis, id, LiniFacIn)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrKasusTidakAda
	}
	if _, err := jalankan(ctx, tx, sqlPastikanGeneral(t.general), "memastikan "+TabelGeneralPolis, id, id); err != nil {
		return err
	}
	for _, q := range sqlHapusObjek(t) {
		if _, err := jalankan(ctx, tx, q, "menghapus objek", id); err != nil {
			return err
		}
	}
	k := db.KosongJadiNil
	for i, o := range baris {
		ids := make([]string, 4)
		for j, seq := range []string{sequenceLocationList, sequenceProperty, sequenceRiskLocation, sequenceBuildingConstruct} {
			if ids[j], err = r.db.NomorBerikut(ctx, tx, seq); err != nil {
				return err
			}
		}
		uid, err := uidAcak()
		if err != nil {
			return err
		}
		langkah := []struct {
			q, tabel string
			arg      []any
		}{
			{sqlSisipLokasi(t.loc), TabelLocationList, []any{ids[0], id, i + 1, uid}},
			{sqlSisipProperty(t.prop), TabelProperty, []any{ids[1], ids[0], k(o.ObjectNo), k(o.ObjectType), k(o.ObjectName),
				teksBool(o.IsMaterialDamage), teksBool(o.IsTopRisk), k(o.RoadType), k(o.RoadName), k(o.BuildingNo), k(o.Country),
				k(o.Province), k(o.RiskAddressID)}},
			{sqlSisipRisk(t.risk), TabelRiskLocation, []any{ids[2], ids[1], k(o.RiskLocation), k(o.City), k(o.District),
				k(o.Territory), k(o.ZipCode)}},
			{sqlSisipBangunan(t.bang), TabelBuildingConstruction, []any{ids[3], ids[1], k(o.FloorType), k(o.NumberOfFloor),
				k(o.RoofType), k(o.WallType), k(o.PartitionType), k(o.SupportWallType), k(o.OtherType)}},
		}
		for _, l := range langkah {
			h, err := jalankan(ctx, tx, l.q, "menyisipkan "+l.tabel, l.arg...)
			if err != nil {
				return err
			}
			if err := db.PastikanSatuBaris(h, l.tabel); err != nil {
				return err
			}
		}
	}
	return nil
}
