package repository

import (
	"errors"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

// TestMasterCargoSQL - view DDL 11-10-2026 (BrowseTrading_RD dst.: hanya ID / Description, tanpa TYPE); bind berurutan; daftar urut
// DESCRIPTION, batas 500; jenis tak dikenal / MaxAge / agen = galat jelas.
func TestMasterCargoSQL(t *testing.T) {
	for q, mau := range map[string]int{sqlDeskripsiMaster("S.TRADING"): 1, sqlIDMaster("S.TRADING"): 1, sqlDaftarMaster("S.TRADING"): 1} {
		if n := bindTerbesar(t, q); n != mau || strings.Contains(q, "TYPE") {
			t.Errorf("%d bind: %s", n, q)
		}
		if err := db.PeriksaSQL(q); err != nil {
			t.Error(err)
		}
	}
	if q := sqlDaftarMaster("S.X"); !strings.HasSuffix(q, "ORDER BY DESCRIPTION, ID FETCH FIRST :1 ROWS ONLY") {
		t.Error(q)
	}
	if q := sqlIDMaster("S.X"); q != "SELECT ID FROM S.X WHERE DESCRIPTION = :1 ORDER BY ID FETCH FIRST 1 ROWS ONLY" {
		t.Error(q)
	}
	for jenis, v := range map[string]string{"Trading": "TRADING", "Goods": "GOODSTYPE", "Conveyance": "CONVEYANCE", "Packing": "PACKING"} {
		if viewMasterCargo[jenis] != v {
			t.Errorf("%s -> %s", jenis, viewMasterCargo[jenis])
		}
	}
	r := &MasterCargoOracle{}
	if _, err := r.view("MaxAge"); !errors.Is(err, ErrJenisMasterTakDikenal) {
		t.Errorf("jenis: %v", err)
	}
}

// TestMasterMarineSQL2 - DDL kedua: bind berurutan (cari / tanpa cari), agen Name memuat PEKA huruf (BrowseLLOYDAGENT_RD tanpa
// pyCaseInsensitive), MaxAge ID / WARRANTY tak peka huruf, kondisi TYPE 'MARINE' tanpa saringan CONDITIONADD.
func TestMasterMarineSQL2(t *testing.T) {
	for q, mau := range map[string]int{sqlWarrantyMaxAge("S.MAXAGE"): 1, sqlAgen("S.L"): 1, sqlCariAgen("S.L", false): 1,
		sqlCariAgen("S.L", true): 2, sqlCariMaxAge("S.M", false): 1, sqlCariMaxAge("S.M", true): 3, sqlCariKondisi("S.C", false): 2,
		sqlCariKondisi("S.C", true): 3} {
		if n := bindTerbesar(t, q); n != mau {
			t.Errorf("%d bind, mau %d: %s", n, mau, q)
		}
		if err := db.PeriksaSQL(q); err != nil {
			t.Error(err)
		}
	}
	if q := sqlCariAgen("S.L", true); !strings.Contains(q, "WHERE NAME LIKE :1") || strings.Contains(q, "UPPER") {
		t.Error(q)
	}
	if q := sqlCariKondisi("S.C", true); !strings.Contains(q, "WHERE TYPE = :1 AND UPPER(CONDITIONNAME) LIKE :2") ||
		strings.Contains(q, "CONDITIONADD") {
		t.Error(q)
	}
	if polaPekaHuruf(" Ab_% ") != `%Ab\_\%%` || polaPekaHuruf("  ") != "" {
		t.Error(polaPekaHuruf(" Ab_% "))
	}
}
