package skemauji

// Tiruan dua tabel warisan Marketing Officer - `MARKETINGOFFICER` dan `MARKETINGOFFICER_LOG` (kolom katalog DEV
// 03-10-2026, tanpa trigger). Berbeda dengan tiruan lain, keduanya dibuat SEBELUM migrasi: migrasi modul
// `marketingofficer` 760 menambah kolom ke `MARKETINGOFFICER_LOG` (perbaikan log, izin work owner 03-10-2026), jadi
// tabel itu harus sudah ada; `MARKETINGOFFICER` ikut supaya pasangannya lengkap.
//
// ⛔ Nol nama orang, nol data produksi - hanya bentuk tabel.

import "fmt"

// namaTabelTiruanMO adalah tabel yang ditiru dan dibongkar bersama skema uji.
var namaTabelTiruanMO = []string{"MARKETINGOFFICER_LOG", "MARKETINGOFFICER"}

const kolomWarisanMO = `ID VARCHAR2(100), CLIENTID VARCHAR2(100), BRANCHDETAILID VARCHAR2(100), MOLEADER VARCHAR2(100),
  MOSTATUS VARCHAR2(100), CLIENTID2 VARCHAR2(100), BRANCHSTATUS VARCHAR2(100), TEAMGROUP VARCHAR2(100),
  CLIENTNAME VARCHAR2(100), TANGGAL DATE, USERUPDATE VARCHAR2(100), BRANCHPARENT VARCHAR2(10),
  BRANCHDETAILNAME VARCHAR2(100)`

// ddlTiruanMarketingOfficer - DDL kedua tiruan, urut dibuat.
func ddlTiruanMarketingOfficer(skema string) []string {
	return []string{
		fmt.Sprintf(`CREATE TABLE %s.MARKETINGOFFICER (%s, AKSES_LOGIN VARCHAR2(50))`, skema, kolomWarisanMO),
		fmt.Sprintf(`CREATE TABLE %s.MARKETINGOFFICER_LOG (ACTION VARCHAR2(10), %s)`, skema, kolomWarisanMO),
	}
}
