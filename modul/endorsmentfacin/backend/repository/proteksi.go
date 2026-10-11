package repository

// Klep pembatal gerbang 1/3/4/5 - `BrowseOpenProteksiEdm_RD` atas tabel warisan POOLDATA.OPENPROTEKSI_EDM.
//
// `[terverifikasi]` `D:\migrasi\RNM\DDL\BrowseOpenProteksiEdm_RD.xml` (ASM-FW-GISFW-INT-OPENPROTEKSI_EDM!
// BROWSEOPENPROTEKSIEDM_RD, kiriman work owner 11-10-2026), `pyUIFilters`:
//
//	pyFilterLogic = A AND B AND (C OR D) AND F1
//	A  .POLICY_NO  =        Param.Nopolis
//	B  .TYPE       =        Param.Type
//	C  .STS_AKSEP  !=       "2"
//	D  .STS_AKSEP  IS NULL
//	F1 .IDPEGA     =        Param.Idpega
//
// Pemanggil satu-satunya, `Endorsment Fac In\Activity\SetErrorBatalEndorsement_Act.xml` langkah 5.3/5.12/5.16/5.20,
// menyetel Param.Nopolis dan Param.Type (1 Batal, 2 Klaim, 3 Pembayaran, 4 Renewal) dan TIDAK menyetel Param.Idpega
// (`grep -rl BrowseOpenProteksiEdm_RD D:\migrasi\RNM` → 2 berkas: RD itu sendiri dan activity itu).
//
// ⚠️ `[dugaan]` F1 TIDAK diport. Dasarnya: F1 dibuat dan rule terakhir diubah di server DEV (`pxCreateSystemID` /
// `pxUpdateSystemID` pegadevnusare2, 03-12-2025), sedangkan C dan D lahir di produksi; dan parameternya selalu kosong.
// Dua bacaan - versi produksi tanpa F1, atau Pega mengabaikan saringan berparameter kosong - memberi hasil yang sama.
// Hanya bacaan "parameter kosong = IDPEGA IS NULL" yang berbeda; perilaku Pega itu `belum terverifikasi`.
// C/D: `!=` Oracle tidak meloloskan NULL, karena itu D. Arti kode STS_AKSEP "2" `belum terverifikasi` (OQ-020).
// Tipe kolom dari `D:\migrasi\RNM\DDL\OPENPROTEKSI_EDM.txt`: POLICY_NO VARCHAR2(100), TYPE VARCHAR2(10),
// STS_AKSEP VARCHAR2(5) - seluruhnya teks, diikat sebagai teks.

import (
	"context"
	"fmt"
	"strconv"
)

// TabelProteksi - tabel warisan klep; tidak pernah dibuat modul ini.
const TabelProteksi = "OPENPROTEKSI_EDM"

// StsAksepDikecualikan - saringan C `.STS_AKSEP != "2"`.
const StsAksepDikecualikan = "2"

// SqlAdaProteksi - saringan A AND B AND (C OR D); penampung unik.
func SqlAdaProteksi(t string) string {
	return fmt.Sprintf(`SELECT 1 FROM %s o WHERE o.POLICY_NO = :1 AND o.TYPE = :2
	  AND (o.STS_AKSEP <> :3 OR o.STS_AKSEP IS NULL) FETCH FIRST 1 ROWS ONLY`, t)
}

// AdaProteksi - klep `tipe` (1-4) untuk polis itu terisi (`@LengthOfPageList(OpenProteksiEDM….pxResults)>0`).
// Tabel tidak ada → ErrBelumAda (pemanggil memutuskan).
func (g *Gudang) AdaProteksi(ctx context.Context, polis string, tipe int) (bool, error) {
	t, err := g.db.Qualify(TabelProteksi)
	if err != nil {
		return false, err
	}
	_, ada, err := g.satuTeks(ctx, "membaca "+TabelProteksi+" (BrowseOpenProteksiEdm_RD)", SqlAdaProteksi(t),
		polis, strconv.Itoa(tipe), StsAksepDikecualikan)
	return ada, err
}
