package repository

// Tab Object FIRE (tiket 35): .LocationList -> T_LOCATIONLIST (induk T_GENERAL_POLIS) ->
// T_PROPERTY -> T_RISKLOCATION / T_BUILDINGCONSTRUCTION (tabel rancangan, migrasi 186) /
// T_SURROUNDINGRISK (tiket 38, migrasi 187) / T_PROPERTYITEMLIST (tiket 39, migrasi 188, banyak per property).
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
	TabelSurroundingRisk      = "T_SURROUNDINGRISK"
	TabelPropertyItemList     = "T_PROPERTYITEMLIST"
	TabelOccupationList       = "T_OCCUPATIONLIST"
	TabelTableOfLimit         = "T_TABLEOFLIMIT"
	TabelFEAList              = "T_FEALIST"
	TabelListCauseOfLoss      = "T_LISTCAUSEOFLOSS"
	TabelCoinsData            = "T_COINSDATA"
	sequenceLocationList      = "SEQ_T_LOCATIONLIST"
	sequenceProperty          = "SEQ_T_PROPERTY"
	sequenceRiskLocation      = "SEQ_T_RISKLOCATION"
	sequenceBuildingConstruct = "SEQ_T_BUILDINGCONSTRUCTION"
	sequenceSurroundingRisk   = "SEQ_T_SURROUNDINGRISK"
	sequencePropertyItemList  = "SEQ_T_PROPERTYITEMLIST"
	sequenceOccupationList    = "SEQ_T_OCCUPATIONLIST"
	sequenceTableOfLimit      = "SEQ_T_TABLEOFLIMIT"
	sequenceFEAList           = "SEQ_T_FEALIST"
	sequenceListCauseOfLoss   = "SEQ_T_LISTCAUSEOFLOSS"
	sequenceCoinsData         = "SEQ_T_COINSDATA"
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

type tabelObjek struct{ work, general, loc, prop, risk, bang, sekitar, item, okupasi, tol, fea, rugi, koas, cov string }

func (r *ObjekOracle) tabel() (tabelObjek, error) {
	var t tabelObjek
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelWorkPolis, &t.work}, {TabelGeneralPolis, &t.general}, {TabelLocationList, &t.loc}, {TabelProperty, &t.prop},
		{TabelRiskLocation, &t.risk}, {TabelBuildingConstruction, &t.bang}, {TabelSurroundingRisk, &t.sekitar},
		{TabelPropertyItemList, &t.item}, {TabelOccupationList, &t.okupasi}, {TabelTableOfLimit, &t.tol},
		{TabelFEAList, &t.fea}, {TabelListCauseOfLoss, &t.rugi}, {TabelCoinsData, &t.koas}, {TabelCoverageList, &t.cov}} {
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
	"b.PARTITION_TYPE", "b.SUPPORT_WALL_TYPE", "b.OTHERS_TYPE",
	// tiket 38
	"p.OWNERSHIP", "p.IS_PRODUCTION_PROCESS_FLAG", "p.IS_HOT_WORK_PROCESS_FLAG", "p.IS_FLAMMABLE_ITEM_FLAG",
	"s.FRONT_OCCUPATION", "s.FRONT_CONSTRUCTION", "s.FRONT_DISTANCE", "s.FRONT_NOTE",
	"s.LEFT_OCCUPATION", "s.LEFT_CONSTRUCTION", "s.LEFT_DISTANCE", "s.LEFT_NOTE",
	"s.BACK_OCCUPATION", "s.BACK_CONSTRUCTION", "s.BACK_DISTANCE", "s.BACK_NOTE",
	"s.RIGHT_OCCUPATION", "s.RIGHT_CONSTRUCTION", "s.RIGHT_DISTANCE", "s.RIGHT_NOTE",
	"s.HOUSEKEEPING_STATUS", "s.FLOOD_AREA_STATUS", "s.FLOOD_AREA", "s.HOUSEKEEPING_REMARK",
	// tiket 39: kunci item (T_PROPERTYITEMLIST.PARENT_ID); tiket 41: kunci FEA (T_FEALIST.PARENT_ID)
	"TO_CHAR(p.ID)", "TO_CHAR(l.ID)",
	// tiket 42: loss ratio baris lokasi (baca-saja); amount uang lewat TO_CHAR TM9 ber-NLS titik
	angkaKeluar("l.LOSS_RATIO1_YEAR_AMOUNT"), "l.LOSS_RATIO1_YEAR_PERCENT",
	angkaKeluar("l.LOSS_RATIO35_YEAR_AMOUNT"), "l.LOSS_RATIO35_YEAR_PERCENT"}

// sqlBacaObjek - satu baris per lokasi; anak tunggal lewat LEFT JOIN (UNIQUE PARENT_ID, 186/187).
func sqlBacaObjek(t tabelObjek) string {
	return "SELECT " + strings.Join(kolomBacaObjek, ", ") + `
FROM ` + t.loc + ` l
LEFT JOIN ` + t.prop + ` p ON p.PARENT_ID = l.ID
LEFT JOIN ` + t.risk + ` r ON r.PARENT_ID = p.ID
LEFT JOIN ` + t.bang + ` b ON b.PARENT_ID = p.ID
LEFT JOIN ` + t.sekitar + ` s ON s.PARENT_ID = p.ID
WHERE l.PARENT_ID = :1
ORDER BY l.SEQ_NO`
}

// sqlHapusObjek - anak sebelum induk (FK tanpa ON DELETE).
func sqlHapusObjek(t tabelObjek) []string {
	prop := "SELECT p.ID FROM " + t.prop + " p JOIN " + t.loc + " l ON l.ID = p.PARENT_ID WHERE l.PARENT_ID = :1"
	return []string{
		"DELETE FROM " + t.koas + " WHERE PARENT_ID IN (SELECT c.ID FROM " + t.rugi + " c WHERE c.PARENT_ID IN (" + prop + "))",
		"DELETE FROM " + t.rugi + " WHERE PARENT_ID IN (" + prop + ")",
		"DELETE FROM " + t.tol + " WHERE PARENT_ID IN (SELECT o.ID FROM " + t.okupasi + " o WHERE " + syaratIndukOkupasi +
			" AND o.PARENT_ID IN (" + prop + "))",
		"DELETE FROM " + t.okupasi + " o WHERE " + syaratIndukOkupasi + " AND o.PARENT_ID IN (" + prop + ")",
		sqlHapusCoverage(t),
		"DELETE FROM " + t.item + " WHERE PARENT_ID IN (" + prop + ")",
		"DELETE FROM " + t.risk + " WHERE PARENT_ID IN (" + prop + ")",
		"DELETE FROM " + t.bang + " WHERE PARENT_ID IN (" + prop + ")",
		"DELETE FROM " + t.sekitar + " WHERE PARENT_ID IN (" + prop + ")",
		"DELETE FROM " + t.prop + " WHERE PARENT_ID IN (SELECT ID FROM " + t.loc + " WHERE PARENT_ID = :1)",
		"DELETE FROM " + t.fea + " WHERE PARENT_ID IN (SELECT ID FROM " + t.loc + " WHERE PARENT_ID = :1)",
		"DELETE FROM " + t.loc + " WHERE PARENT_ID = :1",
	}
}

// sqlPastikanGeneral - T_LOCATIONLIST ber-FK ke T_GENERAL_POLIS: baris General kosong dibuat
// bila General belum pernah disimpan.
func sqlPastikanGeneral(general string) string {
	return "INSERT INTO " + general + " (ID) SELECT :1 FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM " + general + " WHERE ID = :2)"
}

// sqlSisipLokasi - :5..:8 loss ratio (tiket 42) hasil hitung services; amount uang lewat TO_NUMBER, percent teks.
func sqlSisipLokasi(loc string) string {
	return "INSERT INTO " + loc + " (ID, PARENT_ID, SEQ_NO, ROW_UID, LOSS_RATIO1_YEAR_AMOUNT, LOSS_RATIO1_YEAR_PERCENT," +
		" LOSS_RATIO35_YEAR_AMOUNT, LOSS_RATIO35_YEAR_PERCENT) VALUES (:1, :2, :3, :4, " + angkaMasuk(":5") +
		", :6, " + angkaMasuk(":7") + ", :8)"
}

func sqlSisipProperty(prop string) string {
	return "INSERT INTO " + prop + " (ID, PARENT_ID, OBJECT_NO, OBJECT_TYPE, OBJECT_NAME, IS_MATERIAL_DAMAGE, IS_TOP_RISK," +
		" ROAD_TYPE, ROAD_NAME, BUILDING_NO, COUNTRY, PROVINCE, ALM_RISK_ID, OWNERSHIP, IS_PRODUCTION_PROCESS_FLAG," +
		" IS_HOT_WORK_PROCESS_FLAG, IS_FLAMMABLE_ITEM_FLAG)" +
		" VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16, :17)"
}

func sqlSisipRisk(risk string) string {
	return "INSERT INTO " + risk + " (ID, PARENT_ID, ASM_ADDRESS, ASM_CITY, ASM_DISTRICT, ASMRW, ASM_ZIP_CODE) VALUES (:1, :2, :3, :4, :5, :6, :7)"
}

func sqlSisipBangunan(bang string) string {
	return "INSERT INTO " + bang + " (ID, PARENT_ID, FLOOR_TYPE, NUMBER_OF_FLOOR, ROOF_TYPE, WALL_TYPE, PARTITION_TYPE," +
		" SUPPORT_WALL_TYPE, OTHERS_TYPE) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9)"
}

// sqlSisipSekitar - T_SURROUNDINGRISK (tiket 38): ID, PARENT_ID, empat sisi x 4 (urut FRONT, LEFT,
// BACK, RIGHT), lalu HOUSEKEEPING_STATUS, FLOOD_AREA_STATUS, FLOOD_AREA, HOUSEKEEPING_REMARK.
func sqlSisipSekitar(sekitar string) string {
	return "INSERT INTO " + sekitar + " (ID, PARENT_ID," +
		" FRONT_OCCUPATION, FRONT_CONSTRUCTION, FRONT_DISTANCE, FRONT_NOTE," +
		" LEFT_OCCUPATION, LEFT_CONSTRUCTION, LEFT_DISTANCE, LEFT_NOTE," +
		" BACK_OCCUPATION, BACK_CONSTRUCTION, BACK_DISTANCE, BACK_NOTE," +
		" RIGHT_OCCUPATION, RIGHT_CONSTRUCTION, RIGHT_DISTANCE, RIGHT_NOTE," +
		" HOUSEKEEPING_STATUS, FLOOD_AREA_STATUS, FLOOD_AREA, HOUSEKEEPING_REMARK)" +
		" VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16, :17, :18, :19, :20, :21, :22)"
}

// argSekitar - bind :3..:22 sqlSisipSekitar, urutan sama dengan daftar kolomnya.
func argSekitar(s models.SurroundingRisk) []any {
	k := db.KosongJadiNil
	var arg []any
	for _, x := range []models.SisiRisiko{s.Front, s.Left, s.Back, s.Right} {
		arg = append(arg, k(x.Occupation), k(x.Construction), k(x.Distance), k(x.Note))
	}
	return append(arg, k(s.HousekeepingStatus), k(s.FloodAreaStatus), k(s.FloodArea), k(s.HousekeepingRemark))
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
	var idProperty, idLokasi []string
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
			PartitionType: v("b.PARTITION_TYPE"), SupportWallType: v("b.SUPPORT_WALL_TYPE"), OtherType: v("b.OTHERS_TYPE"),
			Ownership: v("p.OWNERSHIP"), IsProductionProcess: v("p.IS_PRODUCTION_PROCESS_FLAG") == teksBenar,
			IsHotWorkProcess: v("p.IS_HOT_WORK_PROCESS_FLAG") == teksBenar, IsFlammableItem: v("p.IS_FLAMMABLE_ITEM_FLAG") == teksBenar,
			SurroundingRisk: models.SurroundingRisk{
				Front: models.SisiRisiko{Occupation: v("s.FRONT_OCCUPATION"), Construction: v("s.FRONT_CONSTRUCTION"),
					Distance: v("s.FRONT_DISTANCE"), Note: v("s.FRONT_NOTE")},
				Left: models.SisiRisiko{Occupation: v("s.LEFT_OCCUPATION"), Construction: v("s.LEFT_CONSTRUCTION"),
					Distance: v("s.LEFT_DISTANCE"), Note: v("s.LEFT_NOTE")},
				Back: models.SisiRisiko{Occupation: v("s.BACK_OCCUPATION"), Construction: v("s.BACK_CONSTRUCTION"),
					Distance: v("s.BACK_DISTANCE"), Note: v("s.BACK_NOTE")},
				Right: models.SisiRisiko{Occupation: v("s.RIGHT_OCCUPATION"), Construction: v("s.RIGHT_CONSTRUCTION"),
					Distance: v("s.RIGHT_DISTANCE"), Note: v("s.RIGHT_NOTE")},
				HousekeepingStatus: v("s.HOUSEKEEPING_STATUS"), FloodAreaStatus: v("s.FLOOD_AREA_STATUS"),
				FloodArea: v("s.FLOOD_AREA"), HousekeepingRemark: v("s.HOUSEKEEPING_REMARK")},
			Items: []models.ItemObjek{}, Occupations: []models.OkupasiObjek{}, FEA: []models.BarisFEA{},
			LossRecords: []models.CatatanKerugian{}})
		// Loss ratio: amount NUMBER lewat angkaKeluar; percent kolom VARCHAR2 berisi teks desimal - keduanya bacaDesimal.
		lr := &hasil[len(hasil)-1].LossRatio
		// ⚠️ R-ADR34: percent lama yang bukan teks desimal polos (mis. koma / "%") kini GALAT (gagal keras, pola
		// db.UraiDesimal) - sebelumnya diteruskan apa adanya. Contoh data: semua "0".
		if err := bacaDesimalKe(TabelLocationList+" case "+id, teks,
			kolomDesimal{angkaKeluar("l.LOSS_RATIO1_YEAR_AMOUNT"), "LOSS_RATIO1_YEAR_AMOUNT", &lr.OneYearAmount},
			kolomDesimal{"l.LOSS_RATIO1_YEAR_PERCENT", "LOSS_RATIO1_YEAR_PERCENT", &lr.OneYearPercent},
			kolomDesimal{angkaKeluar("l.LOSS_RATIO35_YEAR_AMOUNT"), "LOSS_RATIO35_YEAR_AMOUNT", &lr.ThreeFiveYearAmount},
			kolomDesimal{"l.LOSS_RATIO35_YEAR_PERCENT", "LOSS_RATIO35_YEAR_PERCENT", &lr.ThreeFiveYearPercent}); err != nil {
			return nil, err
		}
		idProperty = append(idProperty, v("TO_CHAR(p.ID)"))
		idLokasi = append(idLokasi, v("TO_CHAR(l.ID)"))
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: objek: %w", err)
	}
	baris.Close()
	item, err := r.bacaItem(ctx, t, id)
	if err != nil {
		return nil, err
	}
	okupasi, err := r.bacaOkupasi(ctx, t, id)
	if err != nil {
		return nil, err
	}
	fea, err := r.bacaFEA(ctx, t, id)
	if err != nil {
		return nil, err
	}
	rugi, err := r.bacaKerugian(ctx, t, id)
	if err != nil {
		return nil, err
	}
	for i, l := range idLokasi {
		if d, ada := fea[l]; ada {
			hasil[i].FEA = d
		}
	}
	for i, p := range idProperty {
		if p == "" {
			continue
		}
		if d, ada := item[p]; ada {
			hasil[i].Items = d
		}
		if d, ada := okupasi[p]; ada {
			hasil[i].Occupations = d
		}
		if d, ada := rugi[p]; ada {
			hasil[i].LossRecords = d
		}
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
		ids := make([]string, 5)
		// ids: 0 lokasi, 1 property, 2 risk, 3 bangunan, 4 sekitar; item bernomor sendiri (sisipItem).
		for j, seq := range []string{sequenceLocationList, sequenceProperty, sequenceRiskLocation, sequenceBuildingConstruct,
			sequenceSurroundingRisk} {
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
			{sqlSisipLokasi(t.loc), TabelLocationList, []any{ids[0], id, i + 1, uid, ikatDesimal(o.LossRatio.OneYearAmount),
				ikatDesimal(o.LossRatio.OneYearPercent), ikatDesimal(o.LossRatio.ThreeFiveYearAmount),
				ikatDesimal(o.LossRatio.ThreeFiveYearPercent)}},
			{sqlSisipProperty(t.prop), TabelProperty, []any{ids[1], ids[0], k(o.ObjectNo), k(o.ObjectType), k(o.ObjectName),
				teksBool(o.IsMaterialDamage), teksBool(o.IsTopRisk), k(o.RoadType), k(o.RoadName), k(o.BuildingNo), k(o.Country),
				k(o.Province), k(o.RiskAddressID), k(o.Ownership), teksBool(o.IsProductionProcess), teksBool(o.IsHotWorkProcess),
				teksBool(o.IsFlammableItem)}},
			{sqlSisipRisk(t.risk), TabelRiskLocation, []any{ids[2], ids[1], k(o.RiskLocation), k(o.City), k(o.District),
				k(o.Territory), k(o.ZipCode)}},
			{sqlSisipBangunan(t.bang), TabelBuildingConstruction, []any{ids[3], ids[1], k(o.FloorType), k(o.NumberOfFloor),
				k(o.RoofType), k(o.WallType), k(o.PartitionType), k(o.SupportWallType), k(o.OtherType)}},
			{sqlSisipSekitar(t.sekitar), TabelSurroundingRisk, append([]any{ids[4], ids[1]}, argSekitar(o.SurroundingRisk)...)},
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
		if err := r.sisipItem(ctx, tx, t, ids[1], o.Items); err != nil {
			return err
		}
		if err := r.sisipOkupasi(ctx, tx, t, ids[1], o.Occupations); err != nil {
			return err
		}
		if err := r.sisipFEA(ctx, tx, t, ids[0], o.FEA); err != nil {
			return err
		}
		if err := r.sisipKerugian(ctx, tx, t, ids[1], o.LossRecords); err != nil {
			return err
		}
	}
	return nil
}
