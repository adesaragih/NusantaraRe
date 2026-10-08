package repository

// Panel `Existing Policy for Master ID` layar Adjustment —
// `Section/InputTreatyInAdjustment.xml` @111283, sel 31/32 (`Policy No` ·
// `Pega ID`), di antara kepala dan panel Old/New.
//
// Sumbernya Activity `FetchTreatyExistingProduction` korpus Adjustment:
//
//	InputData.CARI1 = @if(TreatyIn.EDMState="", TreatyIn.ID, TreatyIn.OLDID)
//	RDB-List  RDBList/FetchTreatyInProductionUsingNooffer.xml:
//	  SELECT DISTINCT NOPOLIS AS CARI1, IDPEGA AS CARI2, QUARTER AS CARI3,
//	         QUARTER_YEAR AS CARI4 FROM POOLDATA.TREATYINPRODUCTION
//	  WHERE SUBSTR(NOOFFER,1,7) = SUBSTR({InputData.CARI1},1,7)
//	  ORDER BY CARI4, CARI3 ASC
//	.CARI2 = @substring(.CARI2,18)
//
// ⛔ `TREATYINPRODUCTION` tabel WARISAN relasional — dibaca, tidak pernah
// ditulis; bukan `M_TREATY_IN`, bukan JSON. Modul Treaty In membaca SQL yang
// sama (`treatyin/backend/repository/warisan_polis.go`); modul ini tidak
// boleh mengimpornya, jadi pembacanya disalin.
//
// ⚠️ DISTINCT atas EMPAT kolom walau panel menampilkan DUA: satu polis yang
// muncul di beberapa kuartal tampil beberapa kali di Pega, dan di sini juga.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyinadjustment/backend/models"
)

// TabelWarisanProduksi - sumber panel polis. Dibaca, tidak pernah ditulis.
const TabelWarisanProduksi = "TREATYINPRODUCTION"

// panjangPrefiksPegaID - `@substring(.CARI2,18)` di Activity-nya.
const panjangPrefiksPegaID = 18

// BacaPolisMaster membaca baris panel `Existing Policy for Master ID`.
//
// `idMaster` sudah dipilih pemanggil menurut `@if(EDMState="", ID, OLDID)`.
func (g *Gudang) BacaPolisMaster(ctx context.Context, idMaster string) ([]models.BarisPolisMaster, error) {
	t, err := g.db.Qualify(TabelWarisanProduksi)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT DISTINCT NOPOLIS, IDPEGA, QUARTER, QUARTER_YEAR
		FROM %s WHERE SUBSTR(NOOFFER,1,7) = SUBSTR(:1,1,7)
		ORDER BY QUARTER_YEAR, QUARTER`, t)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	baris, err := g.db.QueryContext(ctx, q, idMaster)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca polis master %s: %w", idMaster, err)
	}
	defer func() { _ = baris.Close() }()

	// ⛔ Irisan KOSONG, bukan nil — panelnya yang menyatakan "No items".
	hasil := []models.BarisPolisMaster{}
	for baris.Next() {
		var nopolis, idpega, kuartal, tahunKuartal sql.NullString
		if err := baris.Scan(&nopolis, &idpega, &kuartal, &tahunKuartal); err != nil {
			return nil, fmt.Errorf("repository: memindai polis master %s: %w", idMaster, err)
		}
		hasil = append(hasil, models.BarisPolisMaster{
			NomorPolis: nopolis.String,
			PegaID:     potongPrefiksPegaID(idpega.String),
		})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca polis master %s: %w", idMaster, err)
	}
	return hasil, nil
}

// potongPrefiksPegaID menjalankan `@substring(.CARI2,18)`. Nilai yang lebih
// pendek dikembalikan apa adanya — kosong tidak dapat dibedakan dari "tanpa
// pengenal Pega", nilai pendek yang tampil utuh dapat dipertanyakan.
func potongPrefiksPegaID(s string) string {
	if len(s) <= panjangPrefiksPegaID {
		return strings.TrimSpace(s)
	}
	return strings.TrimSpace(s[panjangPrefiksPegaID:])
}
