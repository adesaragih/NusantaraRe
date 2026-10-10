package skemauji

// Tiruan objek Pega Cover Life - `M_COVER_LIFE` (ID VARCHAR2(10) NOT NULL PK + JSONDATA IS JSON), view `COVER_LIFE`
// (teks DEV `SELECT ID, a.JSONDATA.Cover, a.JSONDATA.Note`), dan `M_COVER_LIFE_SEQ` (keputusan work owner 08-10-2026 K0,
// modul `coverlife`). Dibuat SEBELUM migrasi: migrasi MODUL 085-086 berjalan sebelum 900 di skema baru dan MENGUBAH
// objek itu (kolom + buang JSONDATA + buang view) - tanpa tiruan, 085 mati ORA-00942. Dibongkar SESUDAH migrasi mundur
// (086_down mengembalikan JSONDATA dan view).
//
// ⛔ Nol baris: hanya bentuk. Fixture uji ditulis uji modulnya sendiri (nilai UJI-).

import "fmt"

// namaTiruanCover - urutan buang: view dulu, lalu tabel; sequence terpisah.
var (
	namaViewTiruanCover     = []string{"COVER_LIFE"}
	namaTabelTiruanCover    = []string{"M_COVER_LIFE"}
	namaSequenceTiruanCover = []string{"M_COVER_LIFE_SEQ"}
)

// ddlTiruanCover - DDL ketiga objek, urut dibuat.
func ddlTiruanCover(skema string) []string {
	return []string{
		fmt.Sprintf(`CREATE TABLE %s.M_COVER_LIFE (ID VARCHAR2(10) NOT NULL PRIMARY KEY, JSONDATA CLOB
			CONSTRAINT ENSURE_M_COVER_LIFE_JSON CHECK (JSONDATA IS JSON))`, skema),
		fmt.Sprintf(`CREATE VIEW %s.COVER_LIFE AS SELECT ID, a.JSONDATA.Cover, a.JSONDATA.Note FROM %s.M_COVER_LIFE a`, skema, skema),
		fmt.Sprintf(`CREATE SEQUENCE %s.M_COVER_LIFE_SEQ START WITH 5`, skema),
	}
}
