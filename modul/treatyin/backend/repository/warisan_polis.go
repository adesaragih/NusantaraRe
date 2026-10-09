package repository

// Panel `Existing Policy for Master ID` — kanan atas layar Treaty In, kedua
// cabang (gambar 01 dan 26 dokumen desain 5 Oktober 2026).
//
// ⛔ SUMBERNYA DITELUSURI, bukan dikarang. Rantainya tiga langkah, dan
// ketiganya hidup (`pyRuleAvailable = Yes`, nol penjaga `1=2`/`Never`):
//
//  1. `Section/InputTreatyInOffer.xml` @89.949 — panel `pyTitle` "Existing
//     Policy for Master ID", gridnya `pyPageListProperty = PolisList.pxResults`
//     @105.310 berkelas `ASM-FW-GISFW-Data-Search`. Keempat judul kolomnya
//     `Policy No` @110.151 · `Pega ID` @114.008 · `Quarter` @117.127 ·
//     `Quarter Year` @121.484.
//
//  2. `Activity/FetchTreatyExistingProduction.xml` mengisinya:
//     `InputData.CARI1 := @if(TreatyIn.EDMState="", TreatyIn.ID, TreatyIn.OLDID)`
//     lalu `RDB-List` ke rule di bawah, lalu `.CARI2 := @substring(.CARI2,18)`.
//
//  3. `RDBList/FetchTreatyInProductionUsingNooffer.xml` — SQL-nya utuh:
//
//     SELECT DISTINCT NOPOLIS AS CARI1, IDPEGA AS CARI2,
//     QUARTER AS CARI3, QUARTER_YEAR AS CARI4
//     FROM POOLDATA.TREATYINPRODUCTION
//     WHERE SUBSTR(NOOFFER,1,7) = SUBSTR({InputData.CARI1},1,7)
//     ORDER BY CARI4, CARI3 ASC
//
// ⭐ DIADU DENGAN GAMBARNYA dan cocok persis. Kontrak `1001841` (gambar 26)
// memberi SATU baris, `NOPOLIS = RNM-QR.T02.05.2025.11987` dan `IDPEGA =
// ASM-FW-GISFW-WORK NB-147044` — dan layar menampilkan `NB-147044`, yaitu
// `IDPEGA` dipotong 18 aksara pertama, persis langkah `@substring(.CARI2,18)`.
// Kontrak `1001846` (gambar 01) memberi NOL baris, dan gambarnya berbunyi
// `No items`.
//
// ⛔ BACA SAJA. `TREATYINPRODUCTION` tabel WARISAN, 41.936 baris; modul ini
// tidak pernah menulisinya.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/modul/treatyin/backend/models"
)

// TabelWarisanProduksi - tabel warisan yang memuat polis produksi.
const TabelWarisanProduksi = "TREATYINPRODUCTION"

// panjangPrefiksPegaID - `@substring(.CARI2,18)` di Activity-nya.
//
// ⚠️ Angka 18 DISALIN dari Activity, dan ia kebetulan tepat sepanjang
// `"ASM-FW-GISFW-WORK "` yang terukur di POOLDATA. Yang menentukan
// Activity-nya; kecocokan itu memperkuatnya, bukan menggantikannya.
const panjangPrefiksPegaID = 18

// BacaPolisProduksi membaca baris panel `Existing Policy for Master ID`.
//
// ⚠️ Pencocokannya atas TUJUH aksara pertama `NOOFFER`, bukan kesamaan
// penuh — itu yang SQL rule-nya lakukan. Mempersempitnya menjadi `=` akan
// mengosongkan panel yang di Pega berisi.
func (g *Gudang) BacaPolisProduksi(ctx context.Context, id string) ([]models.BarisPolisProduksi, error) {
	t, err := g.db.Qualify(TabelWarisanProduksi)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT DISTINCT NOPOLIS, IDPEGA, QUARTER, QUARTER_YEAR
		FROM %s WHERE SUBSTR(NOOFFER,1,7) = SUBSTR(:1,1,7)
		ORDER BY QUARTER_YEAR, QUARTER`, t)
	baris, err := g.db.QueryContext(ctx, q, id)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca polis produksi %s: %w", id, err)
	}
	defer func() { _ = baris.Close() }()

	// ⛔ Irisan KOSONG, bukan nil - panelnya yang menyatakan "No items".
	hasil := []models.BarisPolisProduksi{}
	for baris.Next() {
		var nopolis, idpega, kuartal, tahunKuartal sql.NullString
		if err := baris.Scan(&nopolis, &idpega, &kuartal, &tahunKuartal); err != nil {
			return nil, fmt.Errorf("repository: memindai polis produksi %s: %w", id, err)
		}
		hasil = append(hasil, models.BarisPolisProduksi{
			NomorPolis:   nopolis.String,
			PegaID:       potongPrefiksPegaID(idpega.String),
			Kuartal:      kuartal.String,
			TahunKuartal: tahunKuartal.String,
		})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca polis produksi %s: %w", id, err)
	}
	return hasil, nil
}

// potongPrefiksPegaID menjalankan `@substring(.CARI2,18)`.
//
// ⛔ Nilai yang LEBIH PENDEK dari 18 aksara dikembalikan apa adanya. Pega
// `@substring` atas indeks di luar panjangnya mengembalikan kosong, dan
// kosong di kolom `Pega ID` tidak dapat dibedakan dari "polis ini tidak punya
// pengenal Pega" — sementara nilai pendek yang tampil utuh dapat dilihat dan
// dipertanyakan.
func potongPrefiksPegaID(s string) string {
	if len(s) <= panjangPrefiksPegaID {
		return strings.TrimSpace(s)
	}
	return strings.TrimSpace(s[panjangPrefiksPegaID:])
}
