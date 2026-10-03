package repository

import (
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/nbfacin/backend/models"
)

// bacaMigrasi - isi berkas migrasi modul ini (tanpa CR).
func bacaMigrasi(t *testing.T, nama string) string {
	t.Helper()
	b, err := os.ReadFile("../migrations/" + nama)
	if err != nil {
		t.Fatal(err)
	}
	return strings.ReplaceAll(string(b), "\r", "")
}

// TestSQLCaseNB - teks SQL tiket 29: baris T_WORK_POLIS hanya menulis kolom yang ADA di
// tabel premiumlistlife (050/059), LINI terikat, waktu dari jam basis data; nol nilai
// disambung ke teks.
func TestSQLCaseNB(t *testing.T) {
	mau := "INSERT INTO UJI.T_WORK_POLIS (ID, LINI, POSITION, FLAG_ONGOING_POLICY, STATUS_WORK, CREATE_OP, CREATE_OP_NAME, TGL_CREATE, TGL_UPDATE)" +
		" VALUES (:1, :2, :3, :4, :5, :6, :7, SYSDATE, SYSDATE)"
	if got := sqlSisipCase("UJI.T_WORK_POLIS"); got != mau {
		t.Errorf("SQL\n%q\nmau\n%q", got, mau)
	}
	// Butir 77: nilai bind keadaan awal VERBATIM teks work owner; pembuat di dua kolom.
	arg := argSisipCase("NB-1", "UJI-USER")
	want := []any{"NB-1", "FAC", "Offer", "0", "Pending-Policy", "UJI-USER", "UJI-USER"}
	if len(arg) != len(want) {
		t.Fatalf("argumen %v, mau %v", arg, want)
	}
	for i := range want {
		if arg[i] != want[i] {
			t.Errorf("bind :%d = %v, mau %v", i+1, arg[i], want[i])
		}
	}
	if LiniFacIn != "FAC" || AwalanCaseNB != "NB-" || SequenceCaseNB != "SEQ_WORK_POLIS_NB" || TabelWorkPolis != "T_WORK_POLIS" {
		t.Error("konstanta butir 76 berubah")
	}
	for urut, mau := range map[string]string{"184352": "NB-184352", " 7 ": "NB-7"} {
		if got := RakitPengenalCaseNB(urut); got != mau {
			t.Errorf("RakitPengenalCaseNB(%q) = %q, mau %q", urut, got, mau)
		}
	}
}

// TestKolomOpportunityCocokMigrasi - kolom INSERT T_NB_OPPORTUNITY = kolom CREATE TABLE
// migrasi 180, dua arah dan urutan; bind :1..:n sebanyak kolom; argumen sebanyak kolom.
func TestKolomOpportunityCocokMigrasi(t *testing.T) {
	sql := bacaMigrasi(t, "180_t_nb_opportunity.sql")
	blok := regexp.MustCompile(`(?s)CREATE TABLE \{skema\}\.T_NB_OPPORTUNITY \((.*?)\n\)`).FindStringSubmatch(sql)
	if blok == nil {
		t.Fatal("CREATE TABLE T_NB_OPPORTUNITY tidak terbaca")
	}
	var ddl []string
	for _, baris := range strings.Split(blok[1], "\n") {
		f := strings.Fields(baris)
		if len(f) == 0 || f[0] == "CONSTRAINT" {
			continue
		}
		ddl = append(ddl, f[0])
	}
	if strings.Join(ddl, ",") != strings.Join(kolomOpportunity, ",") {
		t.Errorf("kolom DDL %v\nkolom INSERT %v", ddl, kolomOpportunity)
	}
	q := sqlSisipOpportunity("UJI.T_NB_OPPORTUNITY")
	if !strings.HasPrefix(q, "INSERT INTO UJI.T_NB_OPPORTUNITY (ID, ESTIMATED_CLOSING_DATE, ") ||
		!strings.HasSuffix(q, ":14, :15)") || strings.Contains(q, ":16") {
		t.Errorf("SQL %q", q)
	}
	tgl := time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC)
	arg := argSisipOpportunity("NB-1", models.Opportunity{EstimatedClosingDate: tgl, BusinessProspectName: " UJI Prospek ", Phase: "UJI"})
	if len(arg) != len(kolomOpportunity) || arg[0] != "NB-1" || arg[1] != tgl || arg[2] != " UJI Prospek " || arg[10] != "UJI" {
		t.Errorf("argumen %v", arg)
	}
	if arg[3] != nil || arg[14] != nil {
		t.Errorf("teks kosong harus NULL: %v", arg)
	}
}

// TestMigrasiSequenceWajibDiisi - butir 76.5: angka awal SEQ_WORK_POLIS_NB milik work
// owner/DBA. Berkas memuat SATU pernyataan dengan penanda {NB_MULAI} (belum diisi) ATAU
// bilangan bulat positif tanpa nol di depan (sudah diisi work owner - 03-10-2026 di DEV).
// Angka itu tidak disalin ke uji ini.
func TestMigrasiSequenceWajibDiisi(t *testing.T) {
	sql := bacaMigrasi(t, "181_seq_work_polis_nb.sql")
	pola := regexp.MustCompile(`(?m)^CREATE SEQUENCE \{skema\}\.SEQ_WORK_POLIS_NB START WITH (\{NB_MULAI\}|[1-9][0-9]*) INCREMENT BY 1 NOCACHE NOCYCLE$`)
	if n := len(pola.FindAllString(sql, -1)); n != 1 || strings.Count(sql, "CREATE SEQUENCE") != 1 {
		t.Errorf("181: %d pernyataan sah, mau tepat 1 (penanda {NB_MULAI} atau angka work owner)", n)
	}
}
