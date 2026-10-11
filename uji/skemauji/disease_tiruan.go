package skemauji

// Tiruan objek Pega Disease Life - `DISEASE_LIFE` (ID VARCHAR2(100) NULLABLE, ICD_CODE VARCHAR2(100), DISEASE
// VARCHAR2(1000), TANPA PK dan TANPA indeks; fakta WO 08-10-2026, modul `diseaselife`). Dibuat SEBELUM migrasi: migrasi
// MODUL 080-081 berjalan sebelum 900 di skema baru dan MEMBACA / MENGUBAH tabel itu (sequence dari ID tertinggi + PK) -
// tanpa tiruan, 080 mati ORA-00942. Dibongkar SESUDAH migrasi mundur (081_down membuang PK, 080_down sequence).
//
// ⛔ Nol baris: hanya bentuk. Fixture uji ditulis uji modulnya sendiri (nilai UJI-).

import "fmt"

// namaTabelTiruanDisease - tabel yang dibongkar sesudah migrasi mundur.
var namaTabelTiruanDisease = []string{"DISEASE_LIFE"}

// ddlTiruanDisease - DDL tabel, urut dibuat.
func ddlTiruanDisease(skema string) []string {
	return []string{
		fmt.Sprintf(`CREATE TABLE %s.DISEASE_LIFE (ID VARCHAR2(100), ICD_CODE VARCHAR2(100), DISEASE VARCHAR2(1000))`, skema),
	}
}
