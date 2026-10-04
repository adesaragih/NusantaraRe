package migrasi

import "testing"

// blokSequence menyusun blok sequence-dari-kueri; setiap bagian dapat diganti uji penolakan.
func blokSequence(namaKatalog, kueri, namaBuat, opsi string) string {
	return "DECLARE\n  n    NUMBER;\n  awal NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_SEQUENCES\n" +
		"   WHERE SEQUENCE_OWNER = UPPER('{skema}') AND SEQUENCE_NAME = '" + namaKatalog + "';\n" +
		"  IF n = 0 THEN\n    " + kueri + ";\n" +
		"    EXECUTE IMMEDIATE 'CREATE SEQUENCE {skema}." + namaBuat + " START WITH ' || awal || ' " + opsi + "';\n" +
		"  END IF;\nEND;"
}

const kueriNomor = "SELECT NVL(MAX(N), 0) + 1 INTO awal FROM (\n" +
	"      SELECT TO_NUMBER(REGEXP_SUBSTR(IDVIEW, '[0-9]+$')) N FROM {skema}.CLIENT WHERE REGEXP_LIKE(IDVIEW, '^ORG-[0-9]+$'))"

// Sequence yang nilai awalnya dihitung dari data (810 Company Detail) dikenali persis - nama, kueri, opsi.
func TestBacaSequenceDariKueri(t *testing.T) {
	s, ok := BacaSequenceDariKueri(blokSequence("SEQ_A", kueriNomor, "SEQ_A", "INCREMENT BY 1 NOCACHE NOCYCLE"))
	if !ok {
		t.Fatal("blok sequence dari kueri tidak dikenali")
	}
	mau := SequenceDariKueri{Nama: "SEQ_A",
		Kueri: "SELECT NVL(MAX(N), 0) + 1 FROM (\n      SELECT TO_NUMBER(REGEXP_SUBSTR(IDVIEW, '[0-9]+$')) N FROM " +
			"{skema}.CLIENT WHERE REGEXP_LIKE(IDVIEW, '^ORG-[0-9]+$'))",
		Opsi: "INCREMENT BY 1 NOCACHE NOCYCLE"}
	if s != mau {
		t.Errorf("dapat %+v, mau %+v", s, mau)
	}
}

// Bentuk lain DITOLAK: blok ini satu-satunya PL/SQL yang boleh menghitung, jadi ia tidak boleh berbuat selain
// membaca lalu membuat sequence yang ditanyakan katalog.
func TestBacaSequenceDariKueriMenolak(t *testing.T) {
	for nama, salah := range map[string]string{
		"nama katalog lain dari yang dibuat": blokSequence("SEQ_A", kueriNomor, "SEQ_B", "INCREMENT BY 1 NOCACHE"),
		"kueri menulis":                      blokSequence("SEQ_A", "DELETE FROM {skema}.CLIENT", "SEQ_A", "NOCACHE"),
		"kueri menyunting di subkueri": blokSequence("SEQ_A",
			"SELECT 1 INTO awal FROM (SELECT 1 FROM DUAL WHERE 1 = (UPDATE {skema}.CLIENT SET X = 1))", "SEQ_A", "NOCACHE"),
		"kueri tanpa INTO awal": blokSequence("SEQ_A", "SELECT 1 FROM DUAL", "SEQ_A", "NOCACHE"),
		"opsi asing":            blokSequence("SEQ_A", kueriNomor, "SEQ_A", "INCREMENT BY 1 MAXVALUE 5"),
		"opsi kosong":           blokSequence("SEQ_A", kueriNomor, "SEQ_A", ""),
		"dua pernyataan di kueri": blokSequence("SEQ_A", kueriNomor+";\n    DROP TABLE {skema}.CLIENT", "SEQ_A",
			"NOCACHE"),
		"blok katalog 901": "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_INDEXES\n" +
			"   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T' AND INDEX_NAME = 'I';\n" +
			"  IF n > 0 THEN\n    EXECUTE IMMEDIATE 'DROP INDEX {skema}.I';\n  END IF;\nEND;",
		"SQL biasa": "CREATE SEQUENCE {skema}.SEQ_A START WITH 1 NOCACHE",
	} {
		if s, ok := BacaSequenceDariKueri(salah); ok {
			t.Errorf("%s: diterima sebagai sequence dari kueri: %+v", nama, s)
		}
	}
}
