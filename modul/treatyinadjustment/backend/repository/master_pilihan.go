package repository

// Picker `Add Revision` / `Add Adjustment Premium` dan dokumen master yang
// tombol `Choose` muat — BACA SAJA.
//
// ⛔ Nol tulisan di berkas ini. `Choose` di Pega juga MENYIMPAN
// (`TreatyInEDMSetValue` [8]–[9]); di aplikasi ini ia menyusun draf di layar
// sampai jalur Save diputuskan (pemilik proses, 7 Oktober 2026).

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyinadjustment/backend/models"
)

// TabelMasterTreatyIn - kepala kontrak Treaty In sistem lama.
const TabelMasterTreatyIn = "TREATY_IN"

// kolomMasterPilihan - `CARI1`…`CARI7` RDB-List `TreatyLoadMasterJoinEdm`.
const kolomMasterPilihan = "ID, TREATYCONTRACTNAME, PROPORTIONTYPE, LEADINGREINSSOURCE, CEDING, COMMENCEMENT, TERMINATION"

// DaftarMasterPilihan membaca grid picker.
//
//	Revisi  `TreatyLoadMasterJoinEdm`     TREATY_IN ∪ TREATY_IN_EDM
//	Premi   `TreatyLoadMasterJoinEdmXOL`  idem, keduanya `PROPORTIONTYPE = 'NonProportional'`
//
// ⚠️ SQL ekspor tanpa `ORDER BY`; `UNION` Oracle mengurutkan lewat
// penghilangan duplikat, dan `ORDER BY 1` di sini menyatakan urutan itu
// alih-alih bergantung padanya.
func (g *Gudang) DaftarMasterPilihan(ctx context.Context, hanyaNonProp bool) ([]models.BarisMasterPilihan, error) {
	a, err := g.db.Qualify(TabelMasterTreatyIn)
	if err != nil {
		return nil, err
	}
	b, err := g.db.Qualify(TabelWarisanPenyesuaian)
	if err != nil {
		return nil, err
	}
	saringA, saringB, arg := "", "", []any{}
	if hanyaNonProp {
		saringA, saringB = " WHERE PROPORTIONTYPE = :1", " WHERE PROPORTIONTYPE = :2"
		arg = []any{"NonProportional", "NonProportional"}
	}
	q := fmt.Sprintf("SELECT %s FROM %s%s UNION SELECT %s FROM %s%s ORDER BY 1",
		kolomMasterPilihan, a, saringA, kolomMasterPilihan, b, saringB)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca daftar master: %w", err)
	}
	defer func() { _ = rows.Close() }()

	out := []models.BarisMasterPilihan{}
	for rows.Next() {
		var k [7]sql.NullString
		if err := rows.Scan(&k[0], &k[1], &k[2], &k[3], &k[4], &k[5], &k[6]); err != nil {
			return nil, fmt.Errorf("repository: membaca baris daftar master: %w", err)
		}
		out = append(out, models.BarisMasterPilihan{
			ID: k[0].String, NamaKontrak: k[1].String, SifatProporsi: k[2].String,
			AsalBisnis: k[3].String, Cedant: k[4].String,
			TanggalMulai: k[5].String, TanggalBerakhir: k[6].String,
		})
	}
	return out, rows.Err()
}

// AdaRevisi - `GetTreatyRevisionID` (`TreatyInRevisi_post` [3]): adakah
// penyesuaian yang pengenalnya diawali `id`?
//
// ⚠️ Ekspor membaca `M_TREATY_IN_EDM.ID`; di sini `TREATY_IN_EDM.ID` —
// pengenal yang sama (280 lawan 280, diukur 5 Oktober 2026), tanpa
// menyentuh tabel ber-JSON. Yang ekspor pakai hanya KOSONG-tidaknya
// hasilnya (`HASIL1 = ""`), jadi cacahnya cukup.
func (g *Gudang) AdaRevisi(ctx context.Context, id string) (bool, error) {
	b, err := g.db.Qualify(TabelWarisanPenyesuaian)
	if err != nil {
		return false, err
	}
	q := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ID LIKE :1 ESCAPE '\' AND ROWNUM = 1`, b)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var n int
	if err := g.db.QueryRowContext(ctx, q, lolosLike(id)+"%").Scan(&n); err != nil {
		return false, fmt.Errorf("repository: memeriksa revisi %s: %w", id, err)
	}
	return n > 0, nil
}

// lolosLike meloloskan `%`, `_`, dan `\` — pengenal dicocokkan APA ADANYA.
func lolosLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}

// kolomKepalaMaster - kolom kepala `TREATY_IN` / `TREATY_IN_EDM` → kunci
// DOKUMEN (ejaan `PetaPendaratan` `T_TREATY_REVISION`).
//
// ⛔ Hanya yang dokumennya punya kunci akar sama. `CLASSOFBUSINESS`,
// `NUSARESHAREPCT`, `BROKERAGEPCT`, dan `EDMDATE` tidak dipetakan: kunci
// akarnya tidak ada di peta pendaratan, dan mengarang ejaannya berarti
// menebak.
var kolomKepalaMaster = [][2]string{
	{"ID", "ID"},
	{"PROPORTIONTYPE", "ProportionType"},
	{"TREATYCONTRACTNAME", "TreatyContractName"},
	{"TERITORIALSCOPE", "TeritorialScope"},
	{"COMMENCEMENT", "Commencement"},
	{"TERMINATION", "Termination"},
	{"LEADINGREINSSOURCE", "LeadingReinsSource"},
	{"LEADINGREINSSOURCEID", "LeadingReinsSourceID"},
	{"CEDING", "Ceding"},
	{"CEDINGID", "CedingID"},
	{"LEADINGREINSID", "LeadingReinsID"},
	{"INFORMATION", "Information"},
	{"POSITIONUSERNAME", "PositionUsername"},
	{"POSITION", "Position"},
	{"STATUSAKSEPTASI", "StatusAkseptasi"},
	{"CHOOSESTATUSAKSEPTASI", "ChooseStatusAkseptasi"},
	{"TREATYYEAR", "TreatyYear"},
}

// Tambahan kepala `TREATY_IN_EDM` — master yang dipilih bisa sebuah revisi.
var kolomKepalaEDM = [][2]string{
	{"OLDID", "OLDID"},
	{"EDMSTATE", "EDMState"},
	{"EDMMATERIALTYPE", "EDMMaterialType"},
}

// BacaDokumenMaster - dokumen yang `SetTreatyIn_Act` (pengenal 7 aksara)
// atau `SetTreatyInEDM_Act` (selainnya) muat ke `TreatyIn`.
//
// ⭐ Dokumennya dari TABEL PENDARATAN — `MASTERID` = pengenalnya, peta yang
// sama dengan layar Adjustment. Kunci kepala yang TIDAK ada di dokumen
// dilengkapi dari kolom `TREATY_IN` / `TREATY_IN_EDM`, sama seperti form
// Treaty In membaca kepalanya: tanpa itu, dokumen yang belum mendarat tidak
// punya `ProportionType`, dan panel Old/New tidak pernah terbuka.
//
// `ada` = barisnya ada di tabel kepala ATAU dokumennya ada di pendaratan.
func (g *Gudang) BacaDokumenMaster(ctx context.Context, id string) (models.SisiPenyesuaian, bool, error) {
	sisi, err := g.bacaSisi(ctx, id)
	if err != nil {
		return models.SisiPenyesuaian{}, false, err
	}
	tabel, kolom := TabelMasterTreatyIn, kolomKepalaMaster
	if len(id) != 7 {
		tabel, kolom = TabelWarisanPenyesuaian, append(append([][2]string{}, kolomKepalaMaster...), kolomKepalaEDM...)
	}
	kepala, adaBaris, err := g.bacaKepala(ctx, tabel, kolom, id)
	if err != nil {
		return models.SisiPenyesuaian{}, false, err
	}
	ada := adaBaris || len(sisi.Medan) > 0
	for k, v := range kepala {
		if _, sudah := sisi.Medan[k]; !sudah {
			sisi.Medan[k] = v
		}
	}
	if err := g.isiKurs(ctx, sisi); err != nil {
		return models.SisiPenyesuaian{}, false, err
	}
	return sisi, ada, nil
}

func (g *Gudang) bacaKepala(ctx context.Context, tabel string, kolom [][2]string, id string) (map[string]string, bool, error) {
	nama, err := g.db.Qualify(tabel)
	if err != nil {
		return nil, false, err
	}
	daftar := make([]string, len(kolom))
	for i, k := range kolom {
		daftar[i] = k[0]
	}
	q := fmt.Sprintf("SELECT %s FROM %s WHERE ID = :1", strings.Join(daftar, ", "), nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, false, err
	}
	sel := make([]sql.NullString, len(kolom))
	tuju := make([]any, len(kolom))
	for i := range sel {
		tuju[i] = &sel[i]
	}
	err = g.db.QueryRowContext(ctx, q, id).Scan(tuju...)
	if err == sql.ErrNoRows {
		return map[string]string{}, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("repository: membaca kepala %s %s: %w", tabel, id, err)
	}
	out := map[string]string{}
	for i, k := range kolom {
		// `NULL` dilewati — kunci yang tidak ada berbeda dari yang kosong.
		if sel[i].Valid {
			out[k[1]] = sel[i].String
		}
	}
	return out, true, nil
}
