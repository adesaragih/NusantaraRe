package repository

// Uji teks SQL tabel produksi lama (Utility1 `SaveJsonPolisTreatyIn_Act`, keputusan work owner 06-10-2026). SQL ini
// dicoba di DEV baca-saja 06-10-2026: SELECT dijalankan, INSERT lolos DBMS_SQL.PARSE dan ekspresi bind-nya.

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestSQLProduksiTanpaJSONTanpaProsedur(t *testing.T) {
	jp := sqlSisipJSONPolis("UJI_SKEMA.JSON_POLIS")
	if strings.Contains(jp, "DATA_JSON") {
		t.Errorf("json_polis tidak boleh menulis DATA_JSON (perintah WO 06-10-2026): %s", jp)
	}
	if !strings.Contains(jp, "TGL_INPUT") || !strings.Contains(jp, "SYSDATE") {
		t.Errorf("TGL_INPUT = SYSDATE seperti PEGA_JSON_POLIS_TREATYIN: %s", jp)
	}
	if c := sqlSisipCapaian("UJI_SKEMA.ACHIEVEMENT"); !strings.Contains(c, "TGL_PROD") || !strings.HasSuffix(c, "SYSDATE)") {
		t.Errorf("ACHIEVEMENT.TGL_PROD = sysdate (SaveAchievementSQL): %s", c)
	}
	isi, err := os.ReadFile("produksi.go")
	if err != nil {
		t.Fatal(err)
	}
	kode := strings.ToUpper(regexp.MustCompile(`(?m)^\s*//.*$`).ReplaceAllString(string(isi), ""))
	for _, kata := range []string{"COMMIT", "BEGIN ", "PEGA_JSON_POLIS", "INSERTUPDATEACHIEVMENT", "DATA_JSON"} {
		if strings.Contains(kode, kata) {
			t.Errorf("produksi.go memuat %q - nol prosedur, nol COMMIT, nol JSON", kata)
		}
	}
}

func TestSQLProduksiKolomDanPenampungSelaras(t *testing.T) {
	for _, c := range []struct {
		nama string
		q    string
		ks   []models.Kolom
	}{
		{"json_polis", sqlSisipJSONPolis("UJI_SKEMA.JSON_POLIS"), models.KolomJSONPolis},
		{"ACHIEVEMENT", sqlSisipCapaian("UJI_SKEMA.ACHIEVEMENT"), models.KolomCapaian},
		{"TREATYINPRODUCTION", sqlSisipProduksi("UJI_SKEMA.TREATYINPRODUCTION"), models.KolomProduksi},
	} {
		if !strings.HasPrefix(c.q, "INSERT INTO UJI_SKEMA.") {
			t.Errorf("%s tanpa skema eksplisit: %s", c.nama, c.q)
		}
		b := models.Baris{}
		for _, k := range c.ks {
			switch {
			case k.Golongan.Desimal():
				b[k.Properti] = "1.5"
			case k.Golongan.Tanggal():
				b[k.Properti] = "2026-10-06"
			default:
				b[k.Properti] = "UJI"
			}
		}
		arg, err := argumenKolom(c.ks, b)
		if err != nil {
			t.Fatalf("%s: %v", c.nama, err)
		}
		pen := regexp.MustCompile(`:(\d+)`).FindAllStringSubmatch(c.q, -1)
		if len(pen) != len(arg) {
			t.Errorf("%s: %d penampung, %d argumen", c.nama, len(pen), len(arg))
		}
		for i, p := range pen { // `:1, :2, ..` menurut urutan kemunculan = urutan argumen
			if p[1] != strconv.Itoa(i+1) {
				t.Fatalf("%s: penampung ke-%d = :%s", c.nama, i+1, p[1])
			}
		}
	}
}
