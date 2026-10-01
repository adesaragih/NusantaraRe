package repository

// Teks SQL pembaca mengikuti RD Pega (PARITAS §2–§6) - tanpa Oracle.

import (
	"strings"
	"testing"

	"nusantarare/modul/mastercontractretrolife/backend/models"
)

func TestSQLGridMengikutiRD(t *testing.T) {
	kasus := []struct {
		nama  string
		sql   string
		wajib []string
	}{
		{"tahun BrowseTreatyYear_Life_RD b765", sqlPilih(KolomTahun, "", "ID ASC")("S.TREATYYEAR_LIFE"),
			[]string{"FROM S.TREATYYEAR_LIFE", "ORDER BY ID ASC"}},
		{"kontrak BrowseTreatyContract_Life_RD b628/b940", sqlPilih(KolomKontrak, "IDTREATYYEAR = :1", "ID ASC")("S.TREATYCONTRACT_LIFE"),
			[]string{"WHERE IDTREATYYEAR = :1", "ORDER BY ID ASC"}},
		{"reinsurer BrowseDetailTreatyReisurerLife_RD b613/b631/b896",
			sqlPilih(KolomReinsurer, "TREATYYEARID = :1 AND TREATYCONTRACTID = :2", "ID DESC")("S.TREATYREINSURER_LIFE"),
			[]string{"TREATYYEARID = :1 AND TREATYCONTRACTID = :2", "ORDER BY ID DESC"}},
		{"jenis BrowseReinsuranceTypeLimit_RD b578/b765", sqlJenisReasuransi("S.REINSURANCETYPE"),
			[]string{"WHERE FLAG = :1", "ORDER BY ID DESC"}},
		{"reinsurer master BrowseCedingCoLife_RD b565/b584/b601/b745", sqlCariReinsurer("S.AGENT"),
			[]string{"ID LIKE :1", "UPPER(CLIENTNAME) LIKE :2", "STATUSACTIVE = :3", "ORDER BY CLIENTNAME ASC"}},
		{"business master BrowseBusinessLife_RD b651/b1057", sqlCariBusiness("S.BUSINESS"),
			[]string{"OLDID LIKE :1", "UPPER(NOTE) LIKE :2", "ORDER BY ID ASC"}},
	}
	for _, k := range kasus {
		rata := strings.Join(strings.Fields(k.sql), " ")
		for _, w := range k.wajib {
			if !strings.Contains(rata, w) {
				t.Errorf("%s: SQL tanpa %q:\n%s", k.nama, w, rata)
			}
		}
	}
}

func TestSaringanMasterVerbatimRD(t *testing.T) {
	if models.FlagJenisReasuransiLife != "1" || idReinsurerLife != "%L0%" ||
		models.StatusMasterReinsurerAktif != "1" || awalanBizLife != "L%" {
		t.Errorf("nilai saringan RD berubah: flag %q, id %q, aktif %q, oldid %q", models.FlagJenisReasuransiLife,
			idReinsurerLife, models.StatusMasterReinsurerAktif, awalanBizLife)
	}
}

// ⛔ OQ-MCRL-13: nol pembaca sumber tabel rate di repository sampai disetujui.
func TestNolPembacaSumberRate(t *testing.T) {
	for jalur, isi := range berkasModul(t, ".go") {
		if strings.HasSuffix(jalur, "_test.go") {
			continue
		}
		for _, kata := range []string{"USEDBY", "M_RATE", "_SUMMARY"} {
			if strings.Contains(buangKomentarGo(isi), kata) {
				t.Errorf("%s menyebut %q - sumber tabel rate menunggu persetujuan (OQ-MCRL-13)", jalur, kata)
			}
		}
	}
}

func TestPilihMemintaTanggalDanUangSebagaiTeks(t *testing.T) {
	q := daftarPilih(KolomKontrak)
	for _, k := range []string{"IDR", "USD", "B_IDR", "B_USD", "IDR_SELISIH", "USD_SELISIH"} {
		if !strings.Contains(q, "TO_CHAR("+k+", 'TM9'") {
			t.Errorf("kolom uang %s tidak diminta sebagai teks TM9: %s", k, q)
		}
	}
	for _, k := range []string{"TGLUPDATE", "TREATYSTARTDATE", "TREATYENDDATE"} {
		if !strings.Contains(q, "TO_CHAR("+k+", 'YYYY-MM-DD HH24:MI:SS')") {
			t.Errorf("kolom tanggal %s tidak diminta sebagai teks: %s", k, q)
		}
	}
	if strings.Contains(daftarPilih([]string{"RIRATE"}), "TO_CHAR") {
		t.Error("RIRATE teks (R7) tidak boleh dikonversi")
	}
}

func TestPolaCariMeloloskanWildcard(t *testing.T) {
	if got := PolaCari(" a%b_c\\ "); got != `%A\%B\_C\\%` {
		t.Errorf("PolaCari = %q", got)
	}
}
