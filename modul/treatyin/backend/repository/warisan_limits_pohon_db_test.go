//go:build db

package repository_test

// Bukti Oracle — pohon tab Limits proporsional dari dokumen SUNGGUHAN.
//
// ⛔ BACA SAJA, di dalam transaksi `siapkan` yang selalu di-rollback.

import (
	"testing"

	"nusantarare/modul/treatyin/backend/repository"
)

// ⭐ Cacah tingkat 1 dan 2 pohon = cacah `JSON_TABLE` pada dokumen yang sama.
func TestPohonLimitsCocokDenganOracle(t *testing.T) {
	tx, skema, ctx := siapkan(t)

	var id, dok string
	if err := tx.QueryRowContext(ctx, `SELECT m.ID, m.JSONDATA FROM `+skema+`.M_TREATY_IN m
		JOIN `+skema+`.TREATY_IN t ON t.ID = m.ID
		WHERE t.PROPORTIONTYPE = 'Proportional'
		  AND JSON_EXISTS(m.JSONDATA, '$.Limits[0].Detail[0].IOOLimitList[0]')
		ORDER BY m.ID FETCH FIRST 1 ROWS ONLY`).Scan(&id, &dok); err != nil {
		t.Skipf("lewati: nol dokumen proporsional berpohon: %v", err)
	}
	var nLimit, nDetail int
	if err := tx.QueryRowContext(ctx, `SELECT
		(SELECT COUNT(*) FROM `+skema+`.M_TREATY_IN m, JSON_TABLE(m.JSONDATA, '$.Limits[*]' COLUMNS (x NUMBER PATH '$.Layer')) j WHERE m.ID = :1),
		(SELECT COUNT(*) FROM `+skema+`.M_TREATY_IN m, JSON_TABLE(m.JSONDATA, '$.Limits[*].Detail[*]' COLUMNS (x NUMBER PATH '$.QSPct')) j WHERE m.ID = :1)
		FROM DUAL`, id).Scan(&nLimit, &nDetail); err != nil {
		t.Fatal(err)
	}

	p := repository.PohonLimitsDariDokumen([]byte(dok))
	jumlahDetail := 0
	for _, l := range p {
		d, _ := l["Detail"].([]map[string]any)
		jumlahDetail += len(d)
	}
	if len(p) != nLimit || jumlahDetail != nDetail {
		t.Errorf("kontrak %s: pohon %d/%d, Oracle %d/%d", id, len(p), jumlahDetail, nLimit, nDetail)
	}
	d, _ := p[0]["Detail"].([]map[string]any)
	if io, _ := d[0]["IOOLimitList"].([]map[string]any); len(io) == 0 {
		t.Errorf("kontrak %s: IOOLimitList kosong di pohon", id)
	}
}
