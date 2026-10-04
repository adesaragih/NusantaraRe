package backend

// Migrasi 840-842 dan slot menu 994 - TANPA Oracle: bentuk SQL-nya.

import (
	"reflect"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
)

func langkah(t *testing.T, nama string) []string {
	t.Helper()
	p, err := migrasi.PernyataanLangkah(berkasMigrasi, nama)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// 840: tiga kolom tanpa default (baris lama tetap kosong); 841: PK atas ID; mundur membuang keduanya.
func TestMigrasiKolomDanPK(t *testing.T) {
	p := langkah(t, "840_t_m_account_kolom.sql")
	if len(p) != 1 {
		t.Fatalf("840: %q", p)
	}
	if tabel, kolom := migrasi.KolomAlterTambah(p[0]); tabel != "T_M_ACCOUNT" || !reflect.DeepEqual(kolom, []string{"CREATEDATE", "CREATEOP", "DESCRIPTION"}) {
		t.Errorf("840 kolom: %s %v", tabel, kolom)
	}
	if strings.Contains(strings.ToUpper(p[0]), "DEFAULT") {
		t.Error("840 memberi default - baris lama akan mendapat nilai migrasi")
	}
	if m := langkah(t, "840_t_m_account_kolom_down.sql"); !reflect.DeepEqual(m, []string{"ALTER TABLE {skema}.T_M_ACCOUNT DROP (CREATEDATE, CREATEOP, DESCRIPTION)"}) {
		t.Errorf("840 mundur: %q", m)
	}
	if p := langkah(t, "841_t_m_account_pk.sql"); !reflect.DeepEqual(p, []string{"ALTER TABLE {skema}.T_M_ACCOUNT ADD CONSTRAINT PK_T_M_ACCOUNT PRIMARY KEY (ID)"}) {
		t.Errorf("841: %q", p)
	}
	if m := langkah(t, "841_t_m_account_pk_down.sql"); !reflect.DeepEqual(m, []string{"ALTER TABLE {skema}.T_M_ACCOUNT DROP CONSTRAINT PK_T_M_ACCOUNT"}) {
		t.Errorf("841 mundur: %q", m)
	}
}

// 842: nilai awal = nomor ACC terbesar + 1, dihitung di basis data tempat migrasi berjalan (keputusan work owner
// 04-10-2026); bentuk blok sequence-dari-kueri.
func TestMigrasiSequenceDariNomorTerbesar(t *testing.T) {
	p := langkah(t, "842_seq_t_m_account.sql")
	if len(p) != 1 {
		t.Fatalf("842: %q", p)
	}
	s, ok := migrasi.BacaSequenceDariKueri(p[0])
	if !ok || s.Nama != "SEQ_T_M_ACCOUNT" || s.Opsi != "INCREMENT BY 1 NOCACHE NOCYCLE" {
		t.Fatalf("842 bukan blok sequence-dari-kueri: %+v %v", s, ok)
	}
	if !strings.Contains(s.Kueri, `NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^ASM-SFAGIS-WORK-ACCOUNT ACC-([0-9]+)$', 1, 1, NULL, 1))), 0) + 1`) ||
		!strings.Contains(s.Kueri, "FROM {skema}.T_M_ACCOUNT") {
		t.Errorf("842 kueri nilai awal: %s", s.Kueri)
	}
	if m := langkah(t, "842_seq_t_m_account_down.sql"); !reflect.DeepEqual(m, []string{"DROP SEQUENCE {skema}.SEQ_T_M_ACCOUNT"}) {
		t.Errorf("842 mundur: %q", m)
	}
}

func TestMigrasiSlotMenu(t *testing.T) {
	if p := langkah(t, "994_menu_accounts.sql"); !reflect.DeepEqual(p, []string{"UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE\nWHERE KODE = 'accounts'"}) {
		t.Errorf("994: %q", p)
	}
}
