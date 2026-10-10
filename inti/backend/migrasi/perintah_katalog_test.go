package migrasi

import (
	"reflect"
	"testing"
)

// Blok berpelindung katalog 901 dikenali persis - dan bentuk lain DITOLAK,
// supaya penjaga yang memakainya tidak menerima blok yang berbuat lain.
func TestBacaPerintahKatalog(t *testing.T) {
	blok := "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS\n" +
		"   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_NAV_MENU' AND COLUMN_NAME = 'PARENT_ID';\n" +
		"  IF n > 0 THEN\n    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_NAV_MENU DROP COLUMN PARENT_ID';\n  END IF;\nEND;"
	p, ok := BacaPerintahKatalog(blok)
	if !ok {
		t.Fatal("blok berpelindung katalog tidak dikenali")
	}
	mau := PerintahKatalog{Katalog: "ALL_TAB_COLUMNS", Tabel: "M_NAV_MENU", Objek: "PARENT_ID",
		BilaAda: true, Perintah: "ALTER TABLE {skema}.M_NAV_MENU DROP COLUMN PARENT_ID"}
	if p != mau {
		t.Errorf("dapat %+v, mau %+v", p, mau)
	}
	tambah := "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_INDEXES\n" +
		"   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_NAV_MENU' AND INDEX_NAME = 'IX_A';\n" +
		"  IF n = 0 THEN\n    EXECUTE IMMEDIATE 'CREATE INDEX {skema}.IX_A ON {skema}.M_NAV_MENU (A)';\n  END IF;\nEND;"
	if p, ok := BacaPerintahKatalog(tambah); !ok || p.BilaAda || p.Katalog != "ALL_INDEXES" || p.Objek != "IX_A" {
		t.Errorf("blok n = 0: %+v, %v", p, ok)
	}
	for nama, salah := range map[string]string{
		"tanpa pemeriksaan": "BEGIN\n  EXECUTE IMMEDIATE 'DROP INDEX {skema}.IX_A';\nEND;",
		"katalog lain": "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.DBA_OBJECTS\n" +
			"   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T' AND COLUMN_NAME = 'A';\n" +
			"  IF n > 0 THEN\n    EXECUTE IMMEDIATE 'X';\n  END IF;\nEND;",
		"dua perintah": "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_INDEXES\n" +
			"   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T' AND INDEX_NAME = 'I';\n" +
			"  IF n > 0 THEN\n    EXECUTE IMMEDIATE 'DROP INDEX {skema}.I';\n    EXECUTE IMMEDIATE 'DROP TABLE {skema}.T';\n  END IF;\nEND;",
		"SQL biasa": "DELETE FROM {skema}.M_NAV_MENU WHERE PARENT_ID IS NOT NULL",
		"perintah kosong": "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_INDEXES\n" +
			"   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T' AND INDEX_NAME = 'I';\n" +
			"  IF n > 0 THEN\n    EXECUTE IMMEDIATE '';\n  END IF;\nEND;",
	} {
		if _, ok := BacaPerintahKatalog(salah); ok {
			t.Errorf("%s: diterima sebagai blok berpelindung katalog", nama)
		}
	}
}

// Kolom yang dibuang ALTER ... DROP COLUMN terbaca - juga di dalam EXECUTE
// IMMEDIATE blok berpelindung katalog, tempat 901 menulisnya.
func TestKolomAlterBuang(t *testing.T) {
	for p, mau := range map[string][2]string{
		"ALTER TABLE {skema}.M_NAV_MENU DROP COLUMN PARENT_ID":                        {"M_NAV_MENU", "PARENT_ID"},
		"  EXECUTE IMMEDIATE 'alter table {skema}.m_nav_menu drop column parent_id';": {"M_NAV_MENU", "PARENT_ID"},
		"ALTER TABLE {skema}.M_NAV_MENU DROP CONSTRAINT FK_M_NAV_MENU_INDUK":          {"", ""},
		"ALTER TABLE {skema}.M_NAV_MENU ADD (PARENT_ID NUMBER(10))":                   {"", ""},
	} {
		tabel, kolom := KolomAlterBuang(p)
		var dapat [2]string
		dapat[0] = tabel
		if len(kolom) == 1 {
			dapat[1] = kolom[0]
		}
		if !reflect.DeepEqual(dapat, mau) {
			t.Errorf("%q: dapat %v, mau %v", p, dapat, mau)
		}
	}
}

// Bentuk kedua dan ketiga (R/I Risk, keputusan work owner 08-10-2026): ALL_VIEWS (DROP VIEW hanya bila objeknya VIEW)
// dan ALL_IND_COLUMNS (CREATE INDEX hanya bila belum ada indeks berkolom pertama itu). Varian lain DITOLAK.
func TestBacaPerintahKatalogViewDanKolomIndeks(t *testing.T) {
	view := "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_VIEWS\n" +
		"   WHERE OWNER = UPPER('{skema}') AND VIEW_NAME = 'V_A';\n" +
		"  IF n > 0 THEN\n    EXECUTE IMMEDIATE 'DROP VIEW {skema}.V_A';\n  END IF;\nEND;"
	if p, ok := BacaPerintahKatalog(view); !ok || p != (PerintahKatalog{Katalog: "ALL_VIEWS", Tabel: "V_A", Objek: "V_A",
		BilaAda: true, Perintah: "DROP VIEW {skema}.V_A"}) {
		t.Errorf("blok ALL_VIEWS: %+v %v", p, ok)
	}
	indeks := "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_IND_COLUMNS\n" +
		"   WHERE TABLE_OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_A' AND COLUMN_NAME = 'C' AND COLUMN_POSITION = 1;\n" +
		"  IF n = 0 THEN\n    EXECUTE IMMEDIATE 'CREATE INDEX {skema}.IX_T_A_C ON {skema}.T_A (C)';\n  END IF;\nEND;"
	if p, ok := BacaPerintahKatalog(indeks); !ok || p != (PerintahKatalog{Katalog: "ALL_IND_COLUMNS", Tabel: "T_A", Objek: "C",
		BilaAda: false, Perintah: "CREATE INDEX {skema}.IX_T_A_C ON {skema}.T_A (C)"}) {
		t.Errorf("blok ALL_IND_COLUMNS: %+v %v", p, ok)
	}
	for nama, salah := range map[string]string{
		"ALL_VIEWS ber-OWNER lain": "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_VIEWS\n" +
			"   WHERE OWNER = 'X' AND VIEW_NAME = 'V_A';\n  IF n > 0 THEN\n    EXECUTE IMMEDIATE 'DROP VIEW {skema}.V_A';\n  END IF;\nEND;",
		"ALL_IND_COLUMNS tanpa COLUMN_POSITION": "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_IND_COLUMNS\n" +
			"   WHERE TABLE_OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_A' AND COLUMN_NAME = 'C';\n" +
			"  IF n = 0 THEN\n    EXECUTE IMMEDIATE 'CREATE INDEX {skema}.I ON {skema}.T_A (C)';\n  END IF;\nEND;",
		"ALL_IND_COLUMNS ber-OWNER": "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_IND_COLUMNS\n" +
			"   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_A' AND COLUMN_NAME = 'C' AND COLUMN_POSITION = 1;\n" +
			"  IF n = 0 THEN\n    EXECUTE IMMEDIATE 'CREATE INDEX {skema}.I ON {skema}.T_A (C)';\n  END IF;\nEND;",
	} {
		if _, ok := BacaPerintahKatalog(salah); ok {
			t.Errorf("%s: diterima", nama)
		}
	}
}
