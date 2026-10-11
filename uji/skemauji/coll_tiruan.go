package skemauji

// Tiruan objek Pega Cause Of Loss Life - `M_CAUSEOFLOSS_LIFE` (ID VARCHAR2(10) NOT NULL PK + JSONDATA IS JSON), view
// `CAUSEOFLOSS_LIFE` (teks DEV `a.JSONDATA.CauseofLoss`), dan `M_CAUSEOFLOSS_LIFE_SEQ` (keputusan work owner 08-10-2026
// K0, modul `causeoflosslife`). Dibuat SEBELUM migrasi: migrasi MODUL 090-092 berjalan sebelum 900 di skema baru dan
// MENGUBAH objek itu (RENAME + kolom + buang JSONDATA) - tanpa tiruan, 091 mati ORA-00942. Dibongkar SESUDAH migrasi
// mundur (090_down mengembalikan nama dan view).
//
// ⛔ Nol baris: hanya bentuk. Fixture uji ditulis uji modulnya sendiri (nilai UJI-).

import "fmt"

// namaTiruanCOL - urutan buang: view dulu, lalu tabel; sequence terpisah.
var (
	namaViewTiruanCOL     = []string{"CAUSEOFLOSS_LIFE"}
	namaTabelTiruanCOL    = []string{"CAUSEOFLOSS_LIFE", "M_CAUSEOFLOSS_LIFE"}
	namaSequenceTiruanCOL = []string{"M_CAUSEOFLOSS_LIFE_SEQ"}
)

// ddlTiruanCOL - DDL ketiga objek, urut dibuat.
func ddlTiruanCOL(skema string) []string {
	return []string{
		fmt.Sprintf(`CREATE TABLE %s.M_CAUSEOFLOSS_LIFE (ID VARCHAR2(10) NOT NULL PRIMARY KEY, JSONDATA CLOB
			CONSTRAINT ENSURE_M_CAUSEOFLOSS_LIFE_JSON CHECK (JSONDATA IS JSON))`, skema),
		fmt.Sprintf(`CREATE VIEW %s.CAUSEOFLOSS_LIFE AS SELECT a.ID, a.JSONDATA.CauseofLoss FROM %s.M_CAUSEOFLOSS_LIFE a`, skema, skema),
		fmt.Sprintf(`CREATE SEQUENCE %s.M_CAUSEOFLOSS_LIFE_SEQ START WITH 5`, skema),
	}
}
