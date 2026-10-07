package repository

// Untuk apa berkas ini: RANTAI GENERASI POLIS - pembacaan dan penulisan kunci generasi yang khas endorsemen
// (spec-penyimpanan ID-8..ID-15; diagram sheet EDM Treaty In Prop F9-F17). Pengganti SQL lama:
//
//	RDBList/FetchNopolisCount        Select DATA_JSON ... from POOLDATA.JSON_POLIS where NOPOLIS = {CARI1}
//	RDBList/FetchNoOfferFromNoPolis  select nooffer ... from POOLDATA.treatyinproduction where nopolis = .. and rownum = 1
//	RDBList/FetchPolisJsonPolis      ... JSON_POLIS where NOPOLIS = .. order by PRODKE desc fetch first 1 row only
//	RDBList/SelectProdKe             select PRODKE ... from json_polis where substr(nopolis,1,24) = .. order by prodke desc
//	RD GetListEdmTreaty              EndorsementTreaty: .pyStatusWork != "Resolved-Completed" AND
//	                                 .PolicyTreatyIn.PolicyNo = Param.Nopolis
//
// ⛔ Generasi lama dibaca dari T_GENERAL_POLIS_TREATY (penunjuk OLD_POLIS_ID), BUKAN dari JSON_POLIS.DATA_JSON
// (keputusan work owner "aku tidak mau ada json lagi"; NB menulis json_polis tanpa DATA_JSON). Polis lama Pega
// wajib sudah dimuat pemuat dokumen lama NB (`modul/nbtreatyin/backend/alat/pemuatlama`) dan endorsemennya oleh
// pemuat EDM (`backend/alat/pemuatlama`).
//
// ⛔ BUG LAMA DIPERBAIKI: `SelectProdKe` mengurutkan PRODKE TEKS (json_polis VARCHAR2(5) di DEV - "9" > "10") dan
// memotong nomor polis 24 karakter; di sini PRODKE bilangan (kolom NUMBER 320) dan NOPOLIS utuh.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/edmtreatyin/backend/models"
)

const tabelJSONPolis = "JSON_POLIS"

// CacahGenerasiJSONPolis = cacah generasi BERBEDA satu nomor polis di POOLDATA.JSON_POLIS (PRODKE kosong = generasi
// NB 0). Dua pemakai: `CheckNopolisAvailability` (RDB FetchNopolisCount: ada / tidak = cacah > 0) dan penjaga Create
// (tinjauan kode 06-10-2026): json_polis yang memuat generasi lebih banyak dari rantai relasional berarti dokumen
// lama belum dimuat - Create ditolak, bukan melahirkan EDMNo kembar. DISTINCT, bukan COUNT(*): DEV 06-10-2026
// 1 dari 149 polis json_polis berbaris ganda untuk generasi yang sama.
func (g *Gudang) CacahGenerasiJSONPolis(ctx context.Context, nopolis string) (int, error) {
	t, err := g.nama(tabelJSONPolis)
	if err != nil {
		return 0, err
	}
	q := fmt.Sprintf(`SELECT COUNT(DISTINCT NVL(TRIM(TO_CHAR(PRODKE)), '0')) FROM %s WHERE NOPOLIS = :1`, t)
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	var n int
	if err := g.db.QueryRowContext(ctx, q, nopolis).Scan(&n); err != nil {
		return 0, fmt.Errorf("repository: memeriksa nomor polis json_polis: %w", err)
	}
	return n, nil
}

// NoMasterDariNoPolis = `TrtEdmCheckPolicyError` langkah 6-7 (RDB FetchNoOfferFromNoPolis): NOOFFER baris
// TREATYINPRODUCTION pertama ber-NOPOLIS itu ("No Master Treaty" layar Create). Tanpa baris = "".
func (g *Gudang) NoMasterDariNoPolis(ctx context.Context, nopolis string) (string, error) {
	t, err := g.nama(tabelProduksi)
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf(`SELECT NOOFFER FROM %s WHERE NOPOLIS = :1 AND ROWNUM = 1`, t)
	return g.satuTeks(ctx, "membaca nomor master dari produksi", q, nopolis)
}

// AdaEDMBerjalan = `TrtEdmCheckPolicyError` langkah 2-4 (RD GetListEdmTreaty): ada kasus endorsemen polis itu yang
// belum selesai. Nomor polis kasus berjalan = `Quotation.OldPolicyNo` (CreateEDMT langkah 12; NOPOLIS generasi
// kosong sampai selesai).
func (g *Gudang) AdaEDMBerjalan(ctx context.Context, nopolis string) (bool, error) {
	kerja, err := g.nama(tabelKerja)
	if err != nil {
		return false, err
	}
	gen, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return false, err
	}
	quot, err := g.nama(models.TabelQuotation.Nama)
	if err != nil {
		return false, err
	}
	q := fmt.Sprintf(`SELECT COUNT(*)
	   FROM %s w JOIN %s g ON g.ID = w.ID JOIN %s q ON q.POLIS_ID = g.ID
	  WHERE g.PRODKE >= 1 AND w.ID LIKE '%s%%' AND q.OLD_POLICY_NO = :1
	    AND (w.STATUS_WORK IS NULL OR w.STATUS_WORK NOT IN (:2, :3))`, kerja, gen, quot, models.AwalanKasus)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var n int
	if err := g.db.QueryRowContext(ctx, q, nopolis, models.StatusDitolak, models.StatusSelesai).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: memeriksa endorsemen berjalan: %w", err)
	}
	return n > 0, nil
}

// GenerasiPolis - satu generasi polis yang sudah SELESAI (NOPOLIS terisi).
type GenerasiPolis struct {
	ID     string
	ProdKe int
	EDMNo  string
}

// GenerasiTerakhir = `FetchPolisJsonPolis` + `SelectProdKe`: generasi bernomor polis itu dengan PRODKE terbesar -
// dasar OldData dan ProdKe berikutnya (`SetEDMTNoPolis` langkah 3: ProdKe = terbesar + 1). Tanpa generasi =
// `ErrGenerasiPolisTidakAda`.
func (g *Gudang) GenerasiTerakhir(ctx context.Context, tx *db.Tx, nopolis string) (GenerasiPolis, error) {
	gen, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return GenerasiPolis{}, err
	}
	q := fmt.Sprintf(`SELECT ID, PRODKE, NOENDORS FROM %s WHERE NOPOLIS = :1 ORDER BY PRODKE DESC FETCH FIRST 1 ROWS ONLY`, gen)
	if err := db.PeriksaSQL(q); err != nil {
		return GenerasiPolis{}, err
	}
	var r GenerasiPolis
	var prodke sql.NullInt64
	var edm sql.NullString
	err = g.pembaca(tx).QueryRowContext(ctx, q, nopolis).Scan(&r.ID, &prodke, &edm)
	if errors.Is(err, sql.ErrNoRows) {
		return GenerasiPolis{}, ErrGenerasiPolisTidakAda
	}
	if err != nil {
		return GenerasiPolis{}, fmt.Errorf("repository: membaca generasi terakhir polis: %w", err)
	}
	r.ProdKe, r.EDMNo = int(prodke.Int64), teks(edm)
	return r, nil
}

// GenerasiSebelumnya - ID generasi yang menjadi OldData kasus ini: OLD_POLIS_ID bila terisi; sesudah admin menolak
// (penunjuk dilepas, `LepasGenerasi`) = generasi bernomor polis `nopolis` ber-PRODKE `prodke - 1`. "" bila tidak ada.
func (g *Gudang) GenerasiSebelumnya(ctx context.Context, tx *db.Tx, k models.Kasus, nopolis string) (string, error) {
	if k.OldPolisID != "" {
		return k.OldPolisID, nil
	}
	if nopolis == "" || k.ProdKe < 1 {
		return "", nil
	}
	gen, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf(`SELECT ID FROM %s WHERE NOPOLIS = :1 AND PRODKE = :2`, gen)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var id string
	err = g.pembaca(tx).QueryRowContext(ctx, q, nopolis, k.ProdKe-1).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca generasi sebelumnya: %w", err)
	}
	return id, nil
}

// SetelNomorPolisSelesai menulis NOPOLIS generasi endorsemen saat Utility1 (`SaveJsonPolisTreatyInEDM_Act`) -
// generasi resmi menjadi bagian rantai polis; UQ_GP_TREATY_NOPOLIS (NOPOLIS, PRODKE) menolak nomor generasi kembar.
func (g *Gudang) SetelNomorPolisSelesai(ctx context.Context, tx *db.Tx, id, nopolis string) error {
	gen, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, "menyimpan nomor polis generasi endorsemen",
		fmt.Sprintf(`UPDATE %s g SET NOPOLIS = :1 WHERE g.ID = :2 AND g.NOPOLIS IS NULL AND g.PRODKE >= 1`, gen), nopolis, id)
	if err != nil {
		return err
	}
	if n, _ := hasil.RowsAffected(); n != 1 {
		return ErrNomorPolisSudahAda
	}
	return nil
}

// LepasGenerasi - admin menolak (Decision3 No -> End3, tanpa Utility1): generasi endorsemen itu bukan bagian rantai
// polis. ⛔ PENYIMPANGAN SADAR: OLD_POLIS_ID dikosongkan supaya generasi sebelumnya dapat di-endorse lagi (UNIQUE
// OLD_POLIS_ID, ID-10). Di Pega json_polis tidak pernah ditulis untuk endorsemen yang ditolak, sehingga endorsemen
// berikutnya memakai ProdKe / EDMNo yang sama (SelectProdKe membaca json_polis) - di sini sama: NOPOLIS generasi
// tolak tetap kosong dan tidak dihitung GenerasiTerakhir.
func (g *Gudang) LepasGenerasi(ctx context.Context, tx *db.Tx, id string) error {
	gen, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return err
	}
	_, err = jalankan(ctx, tx, "melepas generasi endorsemen yang ditolak",
		fmt.Sprintf(`UPDATE %s g SET OLD_POLIS_ID = NULL WHERE g.ID = :1 AND g.PRODKE >= 1 AND g.NOPOLIS IS NULL`, gen), id)
	return err
}
