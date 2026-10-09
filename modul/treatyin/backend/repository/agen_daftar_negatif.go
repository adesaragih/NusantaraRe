package repository

// Jalur baca status agen untuk `TreatyInCheckCedingBlacklist` — "This name is
// on Agent Negative List".
//
// ⛔ BACA SAJA atas `AGENT`. Nol `INSERT`, nol `UPDATE`, nol `DELETE`, nol DDL.
//
// ⭐ SUMBER, dibaca dari ekspor (korpus Treaty In dan Adjustment identik):
//
//	Section/TreatyInNONProportional.xml sel 16 / 17 (autocomplete Ceding /
//	Business Source) -> RD `BrowseAgentNusaRe_RD`, kelas
//	ASM-FW-GISFW-Int-AGENT, "additional fields":
//	  .StatusActive -> TreatyIn.CedingStatusActive / TreatyIn.SourceStatusActive
//
// Kelas `ASM-FW-GISFW-Int-AGENT` = tabel `AGENT` (lihat `warisan_pemilih.go`),
// dan `.StatusActive` = kolom `STATUSACTIVE`. Activity pemeriksanya sendiri
// tidak membaca basis data — ia membandingkan properti yang autocomplete
// salin; di sini nilai itu dibaca langsung dari kolom asalnya.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// BacaStatusAktifAgen mengembalikan `STATUSACTIVE` apa adanya per `AGENT.ID`.
//
// Pengenal yang tidak ada di `AGENT` TIDAK muncul di peta — pemanggil yang
// memutuskan artinya. Pengenal kosong dilewati; nol pengenal = nol kueri.
//
// ⚠️ Terukur 8 Oktober 2026 atas seluruh `POOLDATA.AGENT` (429 baris):
// `STATUSACTIVE` hanya bernilai `'1'` (378) dan `'0'` (51).
func (g *Gudang) BacaStatusAktifAgen(ctx context.Context, ids []string) (map[string]string, error) {
	keluar := map[string]string{}
	unik := []any{}
	tanda := []string{}
	lihat := map[string]bool{}
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" || lihat[id] {
			continue
		}
		lihat[id] = true
		unik = append(unik, id)
		tanda = append(tanda, fmt.Sprintf(":%d", len(unik)))
	}
	if len(unik) == 0 {
		return keluar, nil
	}
	nama, err := g.db.Qualify(TabelAgen)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT TO_CHAR(ID), STATUSACTIVE FROM %s WHERE ID IN (%s)`, nama, strings.Join(tanda, ", "))
	baris, err := g.db.QueryContext(ctx, q, unik...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca status agen dari %s: %w", TabelAgen, err)
	}
	defer func() { _ = baris.Close() }()
	for baris.Next() {
		var id, status sql.NullString
		if err := baris.Scan(&id, &status); err != nil {
			return nil, fmt.Errorf("repository: membaca baris status agen: %w", err)
		}
		keluar[id.String] = status.String
	}
	return keluar, baris.Err()
}
