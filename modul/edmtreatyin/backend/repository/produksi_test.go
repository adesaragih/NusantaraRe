package repository

// Uji teks SQL dan susunan tulisan tabel produksi lama EDM (Utility1 `SaveJsonPolisTreatyInEDM_Act`, keputusan work
// owner 06-10-2026: json_polis tanpa DATA_JSON). Tanpa Oracle: penjaga IDPEGA, PRODKE, NOPOLIS ACHIEVEMENT, dan
// JN_REAS diuji lewat `prSusunTulisan` atas hasil baca tiruan. Fixture UJI-.

import (
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"nusantarare/modul/edmtreatyin/backend/models"
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

func TestSQLProduksiEDMNoEndorsProdKeJenisReas(t *testing.T) {
	jp := sqlSisipJSONPolis("UJI_SKEMA.JSON_POLIS")
	for _, k := range []string{"NOENDORS", "PRODKE"} {
		if !strings.Contains(jp, k) {
			t.Errorf("json_polis EDM menulis %s (argumen ke-3/ke-4 PEGA_JSON_POLIS_TREATYIN): %s", k, jp)
		}
	}
	pr := sqlSisipProduksi("UJI_SKEMA.TREATYINPRODUCTION")
	if !strings.Contains(pr, "NOENDORS") || strings.Contains(pr, "TREATYGROUPID") {
		t.Errorf("InsertTreatyInProdEDMT_SQL: + NOENDORS, - TREATYGROUPID: %s", pr)
	}
	if q := prSQLCacahNoPolis("UJI_SKEMA.JSON_POLIS"); q != "SELECT COUNT(*) FROM UJI_SKEMA.JSON_POLIS WHERE NOPOLIS = :1" {
		t.Errorf("TreatyInSearchProdKe = %s", q)
	}
	if q := prSQLIDJenisReas("UJI_SKEMA.REINSURANCETYPE"); !strings.HasPrefix(q, "SELECT TO_CHAR(ID) FROM UJI_SKEMA.REINSURANCETYPE WHERE NOTE = :1") {
		t.Errorf("GetReinstypeIDbyName_SQL = %s", q)
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
				b[k.Properti] = "-1.5" // selisih pembatalan
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

func TestSQLProduksiSelisihNegatifUtuh(t *testing.T) {
	for _, k := range models.KolomProduksi {
		if k.Kolom != "PREMI_OGP" {
			continue
		}
		v, err := nilaiTulis(k, "-1000.5")
		if err != nil || !reflect.DeepEqual(v, []any{"-10005", int64(1)}) {
			t.Errorf("PREMI_OGP -1000.5 -> %v %v", v, err)
		}
		return
	}
	t.Fatal("PREMI_OGP tidak ada di katalog")
}

// simpananUji - satu kasus NonProp: dua baris capaian, dua baris produksi yang JN_REAS-nya menunggu pencarian.
func simpananUji(nota string) models.SimpananPolis {
	prod := func(layer string) models.Baris {
		return models.Baris{"IDPEGA": "UJI-EDMT-1", "NOPOLIS": "UJI-POL-1", "NOENDORS": "UJI-POL-1/E01", "LAYER": layer,
			"PREMI_OGP": "-1000", "JN_REAS": "", "PROD_DATE": "2026-10-06 17:07:03"}
	}
	return models.SimpananPolis{
		IDPega: "UJI-EDMT-1",
		JSONPolis: models.Baris{"IDPEGA": "UJI-EDMT-1", "NOPOLIS": "UJI-POL-1", "NOENDORS": "UJI-POL-1/E01",
			"PRODKE": "", "TGL_PROD": "2026-10-06 17:07:03", "USERNAME": "UJI-DH"},
		Capaian: []models.Baris{
			{"IDPEGA": "UJI-EDMT-1", "NOPOLIS": "", "CURRENCY": "USD", "PREMIUM": "-1500"},
			{"IDPEGA": "UJI-EDMT-1", "NOPOLIS": "", "CURRENCY": "IDR", "PREMIUM": "200"},
		},
		Produksi:         []models.Baris{prod("1"), prod("2")},
		NotaJenisReasXOL: nota,
	}
}

func harapArg(t *testing.T, ks []models.Kolom, b models.Baris, ubah map[string]string) []any {
	t.Helper()
	salin := models.Baris{}
	for k, v := range b {
		salin[k] = v
	}
	for k, v := range ubah {
		salin[k] = v
	}
	arg, err := argumenKolom(ks, salin)
	if err != nil {
		t.Fatal(err)
	}
	return arg
}

func TestSusunTulisanProduksiPenjagaIDPega(t *testing.T) {
	tb := prTabelProduksi{"UJI_SKEMA.JSON_POLIS", "UJI_SKEMA.ACHIEVEMENT", "UJI_SKEMA.TREATYINPRODUCTION"}
	for _, c := range []struct {
		nama          string
		nota          string
		baca          prBacaanProduksi
		apa           []string
		prodKe, nopol string
		jnReas        string
	}{
		{"baru", "UJI NOTA XOL", prBacaanProduksi{prodKe: 2, idJenisReas: "7"},
			[]string{"menulis json_polis", "menulis ACHIEVEMENT", "menulis ACHIEVEMENT", "menulis TREATYINPRODUCTION",
				"menulis TREATYINPRODUCTION"}, "2", "UJI-POL-1", "7"},
		// json_polis sudah ada: tidak disisip ulang; NOPOLIS ACHIEVEMENT dari baris lama
		{"json ada", "UJI NOTA XOL", prBacaanProduksi{adaJSON: 1, noPolisJSON: "UJI-POL-LAMA", idJenisReas: "7"},
			[]string{"menulis ACHIEVEMENT", "menulis ACHIEVEMENT", "menulis TREATYINPRODUCTION", "menulis TREATYINPRODUCTION"},
			"", "UJI-POL-LAMA", "7"},
		// InsetTreatyInProdAddendum_Act 6-8: produksi ber-IDPEGA sudah ada -> keluar
		{"produksi ada", "UJI NOTA XOL", prBacaanProduksi{prodKe: 0, adaProduksi: 3},
			[]string{"menulis json_polis", "menulis ACHIEVEMENT", "menulis ACHIEVEMENT"}, "0", "UJI-POL-1", ""},
		// tanpa nota: JN_REAS baris dibiarkan
		{"tanpa nota", "", prBacaanProduksi{prodKe: 1, idJenisReas: "7"},
			[]string{"menulis json_polis", "menulis ACHIEVEMENT", "menulis ACHIEVEMENT", "menulis TREATYINPRODUCTION",
				"menulis TREATYINPRODUCTION"}, "1", "UJI-POL-1", ""},
	} {
		s := simpananUji(c.nota)
		w, err := prSusunTulisan(tb, s, c.baca)
		if err != nil {
			t.Fatalf("%s: %v", c.nama, err)
		}
		var apa []string
		for _, x := range w {
			apa = append(apa, x.apa)
		}
		if !reflect.DeepEqual(apa, c.apa) {
			t.Fatalf("%s: urutan tulisan %v, harap %v", c.nama, apa, c.apa)
		}
		i := 0
		if c.baca.adaJSON == 0 {
			if !reflect.DeepEqual(w[0].arg, harapArg(t, models.KolomJSONPolis, s.JSONPolis, map[string]string{"PRODKE": c.prodKe})) ||
				w[0].q != sqlSisipJSONPolis(tb.jsonPolis) {
				t.Errorf("%s: json_polis %s %v", c.nama, w[0].q, w[0].arg)
			}
			i = 1
		}
		for j, cp := range s.Capaian {
			if !reflect.DeepEqual(w[i+j].arg, harapArg(t, models.KolomCapaian, cp, map[string]string{"NOPOLIS": c.nopol})) {
				t.Errorf("%s: ACHIEVEMENT %d %v", c.nama, j, w[i+j].arg)
			}
		}
		i += len(s.Capaian)
		for j, p := range s.Produksi[:len(w)-i] {
			if !reflect.DeepEqual(w[i+j].arg, harapArg(t, models.KolomProduksi, p, map[string]string{"JN_REAS": c.jnReas})) ||
				w[i+j].q != sqlSisipProduksi(tb.produksi) {
				t.Errorf("%s: TREATYINPRODUCTION %d %v", c.nama, j, w[i+j].arg)
			}
		}
		if s.Produksi[0]["JN_REAS"] != "" || s.JSONPolis["PRODKE"] != "" || s.Capaian[0]["NOPOLIS"] != "" {
			t.Errorf("%s: simpanan pemanggil berubah", c.nama)
		}
	}
}
