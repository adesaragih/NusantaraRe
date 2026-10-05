package skemauji

// Tiruan tabel warisan Adjuster Consultant - `ADJUSTERCONSULTANT` (kolom katalog DEV 05-10-2026, PK `ID`). Dibuat
// SEBELUM migrasi: migrasi modul `adjusterconsultant` 870 menambah kolom IS_ACTIVE (keputusan work owner 05-10-2026).
// Dibongkar bersama tiruan lain sebelum migrasi mundur.
//
// ⛔ Nol nama orang, nol data produksi - hanya bentuk tabel.

import "fmt"

// namaTabelTiruanAdjuster adalah tabel yang ditiru dan dibongkar bersama skema uji.
var namaTabelTiruanAdjuster = []string{"ADJUSTERCONSULTANT"}

// ddlTiruanAdjuster - DDL tiruan, urut dibuat.
func ddlTiruanAdjuster(skema string) []string {
	return []string{
		fmt.Sprintf(`CREATE TABLE %s.ADJUSTERCONSULTANT (ID VARCHAR2(100) NOT NULL, NAME VARCHAR2(500),
  ADDRESS VARCHAR2(1000), TELPNO VARCHAR2(50), USERNAME VARCHAR2(50), EDITDATE DATE, PRIMARY KEY (ID))`, skema),
	}
}
