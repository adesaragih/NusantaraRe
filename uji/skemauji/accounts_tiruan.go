package skemauji

// Tiruan tabel warisan Accounts - `T_M_ACCOUNT` (kolom katalog DEV 04-10-2026, tanpa PK). Dibuat SEBELUM migrasi:
// migrasi modul `accounts` 840 menambah kolom, 841 memasang PK, dan 842 membaca nomor ACC terbesarnya (keputusan
// work owner 04-10-2026). Dibongkar bersama tiruan lain sebelum migrasi mundur - jalur mundur 840-842 menoleransi
// ORA-00942 / ORA-02289.
//
// ⛔ Nol nama orang, nol data produksi - hanya bentuk tabel.

import "fmt"

// namaTabelTiruanAccounts adalah tabel yang ditiru dan dibongkar bersama skema uji.
var namaTabelTiruanAccounts = []string{"T_M_ACCOUNT"}

// ddlTiruanAccounts - DDL tiruan, urut dibuat.
func ddlTiruanAccounts(skema string) []string {
	return []string{
		fmt.Sprintf(`CREATE TABLE %s.T_M_ACCOUNT (ID VARCHAR2(1020) NOT NULL, GROUPBUSINESSID VARCHAR2(128),
  GROUPBUSINESS VARCHAR2(256), INSUREDID VARCHAR2(1020) NOT NULL, INSUREDNAME VARCHAR2(256))`, skema),
	}
}
