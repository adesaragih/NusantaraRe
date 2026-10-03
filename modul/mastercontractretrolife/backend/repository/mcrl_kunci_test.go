package repository

// Perbaikan /code-review (01-10-2026): kunci baris induk dan kolom master yang dibaca saat simpan.

import (
	"errors"
	"strings"
	"testing"
)

func TestSQLKunciBarisForUpdate(t *testing.T) {
	if got := sqlKunci("S.TREATYCONTRACT_LIFE"); got != "SELECT ID FROM S.TREATYCONTRACT_LIFE WHERE ID = :1 FOR UPDATE" {
		t.Errorf("sqlKunci = %q", got)
	}
	for jenis, mau := range map[string]string{"kontrak": TabelKontrak, "reinsurer": TabelReinsurer,
		"security": TabelSecurity, "business": TabelBusiness} {
		if got, err := tabelKunci(jenis); err != nil || got != mau {
			t.Errorf("tabelKunci(%q) = %q, %v; mau %q", jenis, got, err, mau)
		}
	}
	// ⛔ Tahun treaty abadi (nol penghapus) - tidak pernah perlu dikunci; nama asing ditolak.
	for _, jenis := range []string{"tahun", "AGENT", ""} {
		if _, err := tabelKunci(jenis); !errors.Is(err, errJenisKunci) {
			t.Errorf("tabelKunci(%q) seharusnya ditolak: %v", jenis, err)
		}
	}
}

func TestSQLAmbilMasterMembawaKolomSaringan(t *testing.T) {
	if q := sqlAmbilReinsurer("S.AGENT"); !strings.Contains(q, "STATUSACTIVE") {
		t.Errorf("master reinsurer saat simpan tanpa STATUSACTIVE: %q", q)
	}
	if q := sqlAmbilBusiness("S.BUSINESS"); !strings.Contains(q, "OLDID") {
		t.Errorf("master business saat simpan tanpa OLDID: %q", q)
	}
	if idReinsurerLife != "%L0%" || awalanBizLife != "L%" {
		t.Errorf("saringan RD bergeser: %q %q", idReinsurerLife, awalanBizLife)
	}
}
