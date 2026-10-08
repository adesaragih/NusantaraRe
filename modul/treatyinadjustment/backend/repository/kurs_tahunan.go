package repository

// Grid Rate of Exchange layar Adjustment — dari `TREATYEXCHANGEYEARLY`.
//
// ---------------------------------------------------------------------
// ⛔ BUKAN `CurrencyList` dokumen, dan bukan `T_TREATY_CURRENCY`
// ---------------------------------------------------------------------
// Pemilik proses menegaskannya dua kali (4 dan 6 Oktober 2026): Rate of
// Exchange dibaca dari `TREATYEXCHANGEYEARLY`, tabel WARISAN yang hidup.
// `CurrencyList` di dokumen hanyalah SALINAN tabel itu
// (`KOREKSI-ERD-VERSUS-POOLDATA.md` §1.2), dan ia tidak pernah dimuat ke
// tabel pendaratan mana pun — itulah sebab grid kedua panel layar ini
// selalu "No items" sebelum berkas ini ada.
//
// ⚠️ Disaring menurut `TREATYYEAR` sisi itu sendiri, sama dengan pembaca
// Treaty In (`treatyin/backend/repository.BacaKursTahunan`) — modul ini
// tidak boleh mengimpornya, jadi aturannya disalin, bukan dipanggil. Panel
// Old memakai Treaty Year Old, panel New memakai Treaty Year New: revisi
// yang memindah tahun kontrak menampilkan kurs tahun masing-masing.
//
// Kolom → kunci Pega grid (`kerangka.gen.ts` GRID_KURS):
//
//	CURRENCY   -> Currency     IDCURRENCY -> CurrencyID
//	TOIDR      -> Conversion   STARTDATE  -> PeriodStart   ENDDATE -> PeriodEnd

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyinadjustment/backend/models"
)

// TabelKursTahunan - master kurs warisan. Dibaca, tidak pernah ditulis.
const TabelKursTahunan = "TREATYEXCHANGEYEARLY"

// LarikKurs - nama larik grid Rate of Exchange di halaman Pega.
const LarikKurs = "CurrencyList"

// bacaKursTahunan membaca baris kurs satu Treaty Year sebagai larik berkunci
// Pega. Tahun kosong → nol baris, bukan galat: kontrak tanpa tahun tidak
// punya apa pun untuk disaring.
func (g *Gudang) bacaKursTahunan(ctx context.Context, tahun string) ([]map[string]string, error) {
	tahun = strings.TrimSpace(tahun)
	if tahun == "" {
		return nil, nil
	}
	nama, err := g.db.Qualify(TabelKursTahunan)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT CURRENCY, TO_CHAR(IDCURRENCY), TO_CHAR(TOIDR), STARTDATE, ENDDATE
		FROM %s WHERE TREATYYEAR = :1 ORDER BY CURRENCY`, nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, tahun)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s tahun %s: %w", TabelKursTahunan, tahun, err)
	}
	defer func() { _ = rows.Close() }()
	var out []map[string]string
	for rows.Next() {
		var cur, idCur, toidr, mulai, akhir sql.NullString
		if err := rows.Scan(&cur, &idCur, &toidr, &mulai, &akhir); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s: %w", TabelKursTahunan, err)
		}
		out = append(out, map[string]string{
			"Currency":    cur.String,
			"CurrencyID":  idCur.String,
			"Conversion":  toidr.String,
			"PeriodStart": mulai.String,
			"PeriodEnd":   akhir.String,
		})
	}
	return out, rows.Err()
}

// isiKurs memasang grid kurs satu sisi dari Treaty Year sisi itu.
func (g *Gudang) isiKurs(ctx context.Context, sisi models.SisiPenyesuaian) error {
	baris, err := g.bacaKursTahunan(ctx, sisi.Medan["TreatyYear"])
	if err != nil {
		return err
	}
	if len(baris) > 0 {
		sisi.Larik[LarikKurs] = baris
	}
	return nil
}
