package backend

// Migrasi 800-808 dan slot menu 992 - TANPA Oracle: bentuk SQL-nya.

import (
	"reflect"
	"regexp"
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

// Kolom tambahan tiga tabel client terbaca penjaga STRUKTUR (`KolomAlterTambah`) dan jalur mundurnya membuang
// kolom yang sama.
func TestKolomTambahanTigaTabelClient(t *testing.T) {
	for _, k := range []struct {
		maju, mundur, tabel string
		kolom               []string
	}{
		{"805_client_kolom.sql", "805_client_kolom_down.sql", "CLIENT",
			[]string{"PARENT_ID", "NOTE", "CREATED_BY", "CREATED_AT", "UPDATED_BY", "UPDATED_AT"}},
		{"806_client_piclist_kolom.sql", "806_client_piclist_kolom_down.sql", "CLIENT_PICLIST", []string{"GENDER"}},
		{"807_client_address_kolom.sql", "807_client_address_kolom_down.sql", "CLIENT_ADDRESS",
			[]string{"TELFAX_TYPE", "TELFAX_CODE", "TELFAX_NO"}},
	} {
		p := langkah(t, k.maju)
		if len(p) != 1 {
			t.Fatalf("%s: %d pernyataan, mau 1", k.maju, len(p))
		}
		tabel, kolom := migrasi.KolomAlterTambah(p[0])
		if tabel != k.tabel || !reflect.DeepEqual(kolom, k.kolom) {
			t.Errorf("%s: terbaca %s %v, mau %s %v", k.maju, tabel, kolom, k.tabel, k.kolom)
		}
		mundur := langkah(t, k.mundur)
		mau := "ALTER TABLE {skema}." + k.tabel + " DROP (" + strings.Join(k.kolom, ", ") + ")"
		if !reflect.DeepEqual(mundur, []string{mau}) {
			t.Errorf("%s: %q, mau %q", k.mundur, mundur, mau)
		}
	}
}

// M_ENUMERASI berisi tujuh jenis (213 baris) dalam SATU pernyataan INSERT ALL - atomik; COUNTRY tidak di sana.
func TestIsiEnumerasiSatuPernyataan(t *testing.T) {
	p := langkah(t, "801_m_enumerasi_isi.sql")
	if len(p) != 1 || !strings.HasPrefix(p[0], "INSERT ALL") || !strings.HasSuffix(p[0], "SELECT 1 FROM DUAL") {
		t.Fatalf("801 bukan satu INSERT ALL: %d pernyataan", len(p))
	}
	baris := regexp.MustCompile(`INTO \{skema\}\.M_ENUMERASI \(JENIS, KODE, LABEL, AKTIF, URUTAN\) VALUES \('(\w+)', '([^']*)', (?:'[^']*'|NULL), '([01])', \d+\)`).
		FindAllStringSubmatch(p[0], -1)
	if len(baris) != 213 {
		t.Fatalf("801: %d baris terbaca, mau 213", len(baris))
	}
	per, aktif := map[string]int{}, map[string][]string{}
	for _, b := range baris {
		per[b[1]]++
		if b[3] == "1" && (b[1] == "telfax" || b[1] == "jenisalamat" || b[1] == "title" || b[1] == "gender") {
			aktif[b[1]] = append(aktif[b[1]], b[2])
		}
	}
	mau := map[string]int{"title": 7, "bidangusaha": 119, "posisi": 20, "jenisalamat": 8, "telfax": 6, "kodehp": 51, "gender": 2}
	if !reflect.DeepEqual(per, mau) {
		t.Errorf("baris per jenis %v, mau %v", per, mau)
	}
	mauAktif := map[string][]string{"telfax": {"3", "5", "6"}, "jenisalamat": {"1", "2", "7", "8"},
		"title": {"04", "05", "06", "07"}, "gender": {"1", "2"}}
	if !reflect.DeepEqual(aktif, mauAktif) {
		t.Errorf("kode aktif %v, mau %v", aktif, mauAktif)
	}
	if strings.Contains(p[0], "negara") {
		t.Error("801 memuat negara; COUNTRY bersumber NATION")
	}
}

// NATION: buang view, buat tabel bernama dan berkolom SAMA, salin isi dengan ekspresi view lamanya - tiga langkah
// terpisah; jalur mundur 802 membuat ulang view-nya.
func TestNationDariViewMenjadiTabel(t *testing.T) {
	if p := langkah(t, "802_nation_lepas_view.sql"); !reflect.DeepEqual(p, []string{"DROP VIEW {skema}.NATION"}) {
		t.Errorf("802: %q", p)
	}
	p := langkah(t, "803_nation.sql")
	nama, kolom := migrasi.KolomCreateTable(p[0])
	if nama != "NATION" || !reflect.DeepEqual(kolom, []string{"ID", "OLDID", "NOTE", "NATIONINITIAL"}) {
		t.Errorf("803: %s %v", nama, kolom)
	}
	isi := strings.Join(strings.Fields(langkah(t, "804_nation_isi.sql")[0]), " ")
	if isi != "INSERT INTO {skema}.NATION (ID, OLDID, NOTE, NATIONINITIAL) SELECT a.JSONDATA.ID, a.OLDID, "+
		"a.JSONDATA.Note, a.JSONDATA.NationInitial FROM {skema}.M_NATION a" {
		t.Errorf("804: %s", isi)
	}
	view := strings.Join(strings.Fields(langkah(t, "802_nation_lepas_view_down.sql")[0]), " ")
	if !strings.HasPrefix(view, "CREATE OR REPLACE VIEW {skema}.NATION AS SELECT a.JSONDATA.ID, a.OLDID,") {
		t.Errorf("802_down: %s", view)
	}
}

// 808 hanya MENONAKTIFKAN kedua trigger M_CLIENT (tidak membuang); mundurnya menghidupkan lagi.
func TestTriggerMClientDinonaktifkan(t *testing.T) {
	if p := langkah(t, "808_m_client_trigger.sql"); !reflect.DeepEqual(p, []string{
		"ALTER TRIGGER {skema}.TRG_M_CLIENT DISABLE", "ALTER TRIGGER {skema}.TRG_M_CLIENT_PIC DISABLE"}) {
		t.Errorf("808: %q", p)
	}
	if p := langkah(t, "808_m_client_trigger_down.sql"); !reflect.DeepEqual(p, []string{
		"ALTER TRIGGER {skema}.TRG_M_CLIENT_PIC ENABLE", "ALTER TRIGGER {skema}.TRG_M_CLIENT ENABLE"}) {
		t.Errorf("808_down: %q", p)
	}
}

// 809 menonaktifkan EMAIL di Phone and Fax tanpa membuang barisnya; mundurnya mengaktifkan lagi.
func TestTelfaxTanpaEmail(t *testing.T) {
	if p := langkah(t, "809_m_enumerasi_telfax_tanpa_email.sql"); !reflect.DeepEqual(p, []string{
		"UPDATE {skema}.M_ENUMERASI SET AKTIF = '0'\nWHERE JENIS = 'telfax' AND KODE = '6'"}) {
		t.Errorf("809: %q", p)
	}
	if p := langkah(t, "809_m_enumerasi_telfax_tanpa_email_down.sql"); !reflect.DeepEqual(p, []string{
		"UPDATE {skema}.M_ENUMERASI SET AKTIF = '1'\nWHERE JENIS = 'telfax' AND KODE = '6'"}) {
		t.Errorf("809_down: %q", p)
	}
}

// Nol COMMIT di setiap langkah modul ini.
func TestNolCommit(t *testing.T) {
	entri, err := berkasMigrasi.ReadDir("migrations")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entri {
		for _, p := range langkah(t, e.Name()) {
			if strings.Contains(strings.ToUpper(p), "COMMIT") {
				t.Errorf("%s memuat COMMIT", e.Name())
			}
		}
	}
}
