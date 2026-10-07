package repository

// Sub-tab FEA / Fire Extinguisher Availability (tiket 41): .LocationList(n).FEAList -> T_FEALIST (tabel BARU, migrasi
// 190, induk T_LOCATIONLIST), banyak baris per lokasi urut SEQ_NO; ikut baca/ganti objek (objek.go).

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// kolomBacaFEA - kolom sqlBacaFEA BERNAMA -> medan BarisFEA.
var kolomBacaFEA = []string{"TO_CHAR(f.PARENT_ID)", "f.APAR", "f.SPRINKLER", "f.SMOKE_DETECTOR", "f.HYDRANT",
	"f.PRIVATE_TRUCK_BRIGADE", "f.PRIVATE_FIRE_BRIGADE", "f.TEAM_SOP_SAFETY", "f.TEAM_SOP_RISK_MANAGEMENT", "f.INFO_FEA"}

// sqlBacaFEA - seluruh FEA case :1, urut lokasi lalu SEQ_NO.
func sqlBacaFEA(t tabelObjek) string {
	return "SELECT " + strings.Join(kolomBacaFEA, ", ") + `
FROM ` + t.fea + ` f
JOIN ` + t.loc + ` l ON l.ID = f.PARENT_ID
WHERE l.PARENT_ID = :1
ORDER BY f.PARENT_ID, f.SEQ_NO`
}

// sqlSisipFEA - :1 ID, :2 PARENT_ID, :3 SEQ_NO, :4 ROW_UID, :5..:13 medan (urut kolom).
func sqlSisipFEA(fea string) string {
	return "INSERT INTO " + fea + " (ID, PARENT_ID, SEQ_NO, ROW_UID, APAR, SPRINKLER, SMOKE_DETECTOR, HYDRANT," +
		" PRIVATE_TRUCK_BRIGADE, PRIVATE_FIRE_BRIGADE, TEAM_SOP_SAFETY, TEAM_SOP_RISK_MANAGEMENT, INFO_FEA)" +
		" VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13)"
}

// argFEA - bind :5..:13 sqlSisipFEA, urutan sama dengan daftar kolomnya; kosong -> NULL.
func argFEA(f models.BarisFEA) []any {
	k := db.KosongJadiNil
	return []any{k(f.APAR), k(f.Sprinkler), k(f.SmokeDetector), k(f.Hydrant), k(f.PrivateTruckBrigade),
		k(f.PrivateFireBrigade), k(f.TeamSOPSafety), k(f.TeamSOPRiskManagement), k(f.Info)}
}

// bacaFEA - FEA case `id` per T_LOCATIONLIST.ID (teks).
func (r *ObjekOracle) bacaFEA(ctx context.Context, t tabelObjek, id string) (map[string][]models.BarisFEA, error) {
	baris, err := r.db.QueryContext(ctx, sqlBacaFEA(t), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca FEA objek: %w", err)
	}
	defer baris.Close()
	hasil := map[string][]models.BarisFEA{}
	for baris.Next() {
		teks := make(map[string]*sql.NullString, len(kolomBacaFEA))
		tujuan := make([]any, len(kolomBacaFEA))
		for i, k := range kolomBacaFEA {
			teks[k] = &sql.NullString{}
			tujuan[i] = teks[k]
		}
		if err := baris.Scan(tujuan...); err != nil {
			return nil, fmt.Errorf("repository: FEA objek: %w", err)
		}
		v := func(k string) string { return teks[k].String }
		induk := v("TO_CHAR(f.PARENT_ID)")
		hasil[induk] = append(hasil[induk], models.BarisFEA{APAR: v("f.APAR"), Sprinkler: v("f.SPRINKLER"),
			SmokeDetector: v("f.SMOKE_DETECTOR"), Hydrant: v("f.HYDRANT"), PrivateTruckBrigade: v("f.PRIVATE_TRUCK_BRIGADE"),
			PrivateFireBrigade: v("f.PRIVATE_FIRE_BRIGADE"), TeamSOPSafety: v("f.TEAM_SOP_SAFETY"),
			TeamSOPRiskManagement: v("f.TEAM_SOP_RISK_MANAGEMENT"), Info: v("f.INFO_FEA")})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: FEA objek: %w", err)
	}
	return hasil, nil
}

// sisipFEA - FEA satu lokasi, SEQ_NO 1..n.
func (r *ObjekOracle) sisipFEA(ctx context.Context, tx *db.Tx, t tabelObjek, idLokasi string, fea []models.BarisFEA) error {
	for j, f := range fea {
		idFEA, err := r.db.NomorBerikut(ctx, tx, sequenceFEAList)
		if err != nil {
			return err
		}
		uid, err := uidAcak()
		if err != nil {
			return err
		}
		arg := append([]any{idFEA, idLokasi, j + 1, uid}, argFEA(f)...)
		h, err := jalankan(ctx, tx, sqlSisipFEA(t.fea), "menyisipkan "+TabelFEAList, arg...)
		if err != nil {
			return err
		}
		if err := db.PastikanSatuBaris(h, TabelFEAList); err != nil {
			return err
		}
	}
	return nil
}
