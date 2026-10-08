package repository

// Skalar AKAR tab Share Non-Prop — `TreatyIn.RNMShare`, `BrokeragePercent`,
// `RNMShareAcrossTheBoard`, `RnmShareDeducted`.
//
// ⭐ KEPUTUSAN PEMAKAI 6–7 Oktober 2026: Save/Submit menyimpan ke tabel
// masing-masing yang ditentukan (skema v2: akar `TreatyIn` → `T_TREATY_REVISION`),
// bukan mengikuti Pega. Keempat kolomnya ditambahkan migrasi `445`
// (diparkir di folder `treatyinadjustment`, seperti `444` — rentang
// `treatyin` 400-439 habis). Kolom itulah sumber PERTAMA.
//
// Untuk kontrak Pega yang belum pernah disimpan aplikasi, cadangannya nilai
// yang Pega TULIS SENDIRI (nol JSON, nol `M_TREATY_IN`):
//
//	salinan baris Share   `TreatyInNonAddItem` [14.1] `.RNMShare := TreatyIn.RNMShare`;
//	                      deduksi `Brokerage fee` [`TreatyInSetBrokerage` 2]
//	TREATYINDETAIL        `SaveTreatyInDetail_Act` [5.1] `RNM_SHARE := TreatyIn.RNMShare`,
//	                      `BROKERAGE := TreatyIn.BrokeragePercent`
//
// Terukur 7 Oktober 2026: pada 283 kontrak yang punya salinan baris DAN
// `TREATYINDETAIL`, nilainya sama 283/283; `TREATYINDETAIL` menutup 480 dari
// 765 kontrak yang salinan barisnya kosong.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
)

// KolomShareRevisi — kolom `T_TREATY_REVISION` (migrasi 445) ↔ ejaan dokumen.
var KolomShareRevisi = [][2]string{
	{"RNMSHARE", "RNMShare"},
	{"BROKERAGEPERCENT", "BrokeragePercent"},
	{"RNMSHAREACROSSTHEBOARD", "RNMShareAcrossTheBoard"},
	{"RNMSHAREDEDUCTED", "RnmShareDeducted"},
}

// kolomBelumAda — `ORA-00904` (identifier tidak sah): migrasi kolomnya belum
// dijalankan di basis data ini.
func kolomBelumAda(err error) bool {
	return err != nil && strings.Contains(err.Error(), "ORA-00904")
}

// BacaShareAkarRevisi membaca keempat kolom akar Share satu kontrak.
//
// ⚠️ TOLERAN terhadap migrasi yang belum dijalankan: bila kolomnya belum
// ada (`ORA-00904`), jawabannya peta KOSONG tanpa galat — layar tetap
// terbuka dengan sumber cadangan, alih-alih seluruh kontrak gagal dibuka
// hanya karena `-migrate` belum dijalankan. Galat lain tetap galat.
func (g *Gudang) BacaShareAkarRevisi(ctx context.Context, masterID string) (map[string]string, error) {
	nama, err := g.db.Qualify("T_TREATY_REVISION")
	if err != nil {
		return nil, err
	}
	kolom := make([]string, 0, len(KolomShareRevisi))
	for _, k := range KolomShareRevisi {
		kolom = append(kolom, k[0])
	}
	q := fmt.Sprintf("SELECT %s FROM %s WHERE MASTERID = :1", strings.Join(kolom, ", "), nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	sel := make([]sql.NullString, len(kolom))
	tuju := make([]any, len(kolom))
	for i := range sel {
		tuju[i] = &sel[i]
	}
	out := map[string]string{}
	err = g.db.QueryRowContext(ctx, q, masterID).Scan(tuju...)
	if err == sql.ErrNoRows || kolomBelumAda(err) {
		return out, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository: membaca akar Share T_TREATY_REVISION kontrak %s: %w", masterID, err)
	}
	for i, k := range KolomShareRevisi {
		if sel[i].Valid && strings.TrimSpace(sel[i].String) != "" {
			out[k[1]] = sel[i].String
		}
	}
	return out, nil
}

// TabelDetailWarisan — tabel datar yang `SaveTreatyInDetail_Act` tulis
// (kelas `ASM-FW-GISFW-Int-TREATYINDETAIL`). Warisan, BACA SAJA.
const TabelDetailWarisan = "TREATYINDETAIL"

// BacaShareDetailWarisan membaca `RNM_SHARE` / `BROKERAGE` kontrak Non-Prop
// dari `TREATYINDETAIL` — nilai akar yang Pega tulis per baris layer; terukur
// SATU nilai berbeda per kontrak (nol kontrak ganda).
func (g *Gudang) BacaShareDetailWarisan(ctx context.Context, masterID string) (models.ShareDetailWarisan, error) {
	nama, err := g.db.Qualify(TabelDetailWarisan)
	if err != nil {
		return models.ShareDetailWarisan{}, err
	}
	q := fmt.Sprintf(`SELECT TO_CHAR(MAX(RNM_SHARE), 'TM9'), TO_CHAR(MAX(BROKERAGE), 'TM9') FROM %s
		WHERE TREATYID = :1 AND PROPORTIONTYPE = 'NonProportional'`, nama)
	var rnm, brk sql.NullString
	if err := g.db.QueryRowContext(ctx, q, masterID).Scan(&rnm, &brk); err != nil && err != sql.ErrNoRows {
		return models.ShareDetailWarisan{}, fmt.Errorf("repository: membaca %s kontrak %s: %w", TabelDetailWarisan, masterID, err)
	}
	return models.ShareDetailWarisan{RNMShare: angkaOracle(rnm.String), BrokeragePercent: angkaOracle(brk.String)}, nil
}

// angkaOracle — `TO_CHAR(…, 'TM9')` menulis `0,5` sebagai `.5`; ejaan
// desimal aplikasi menuntut nol di depan.
func angkaOracle(s string) string {
	s = strings.TrimSpace(s)
	switch {
	case strings.HasPrefix(s, "."):
		return "0" + s
	case strings.HasPrefix(s, "-."):
		return "-0" + s[1:]
	}
	return s
}
