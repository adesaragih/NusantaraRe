package repository

// Teks SQL pembaca mengikuti RD Pega (PARITAS §2–§6) - tanpa Oracle.

import (
	"errors"
	"regexp"
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

// K1 keputusan work owner 01-10-2026 (OQ-MCRL-13 + OQ-MCRL-05): kedua view rate dibaca SAJA - kolom
// yang dibaca RD XML saja (`BrowseRateLifeSummary` b692–b717, `BrowseRateLife_RD` b747–b791), nol
// `SELECT *`, nol `JSONDATA`; Rate List berkunci `IDUSEDBY`, urut dan batas RD.
func TestRateDibacaKolomRDSaja(t *testing.T) {
	kolomRD := map[string][]string{
		MasterRingkasanRate: {"ID", "USEDBY", "OPERATORID", "MODIFIEDDATE", "TYPE"},
		MasterRate:          {"ID", "IDUSEDBY", "USEDBY", "GENDER", "CONTRACT", "AGE", "RATE", "TYPE"},
	}
	kasus := []struct {
		objek string
		q     string
		mau   []string
	}{
		{MasterRingkasanRate, sqlCariRingkasanRate("S.V"), []string{"UPPER(USEDBY) LIKE :1", "ORDER BY ID ASC", "FETCH FIRST 100 ROWS ONLY"}},
		{MasterRingkasanRate, sqlAmbilRingkasanRate("S.V"), []string{"WHERE ID = :1"}},
		{MasterRate, sqlDaftarRate("S.V"), []string{"WHERE IDUSEDBY = :1", "ORDER BY ID DESC, RATE ASC", "FETCH FIRST 501 ROWS ONLY"}},
	}
	for _, k := range kasus {
		q := strings.Join(strings.Fields(k.q), " ")
		if strings.Contains(q, "*") || strings.Contains(strings.ToUpper(q), "JSONDATA") {
			t.Errorf("%s: SELECT * / JSONDATA dilarang: %s", k.objek, q)
		}
		for _, m := range k.mau {
			if !strings.Contains(q, m) {
				t.Errorf("%s: %q tidak ada: %s", k.objek, m, q)
			}
		}
		pilih := q[len("SELECT "):strings.Index(q, " FROM ")]
		for _, kol := range strings.Split(pilih, ",") {
			kol = strings.TrimSpace(kol)
			ada := false
			for _, r := range kolomRD[k.objek] {
				ada = ada || r == kol
			}
			if !ada {
				t.Errorf("%s: kolom %s bukan kolom RD %v", k.objek, kol, kolomRD[k.objek])
			}
		}
		if err := periksaBacaSaja(k.objek, k.q); err != nil {
			t.Errorf("%s: SELECT ditolak penjaga baca-saja: %v", k.objek, err)
		}
	}
	if BatasRate != 500 {
		t.Errorf("BatasRate %d, RD pyMaxRecords 500", BatasRate)
	}
}

// Uji gigit K1: tulisan ke view rate (dan master lain) ditolak SEBELUM sampai ke Oracle; tabel warisan
// modul ini tetap boleh ditulis.
func TestPeriksaBacaSajaMenolakTulisanKeView(t *testing.T) {
	for _, objek := range []string{MasterRate, MasterRingkasanRate, MasterBusiness} {
		for _, q := range []string{
			"UPDATE S.V SET RATE = :1 WHERE ID = :2",
			"insert into S.V (ID) values (:1)",
			"  DELETE FROM S.V WHERE ID = :1",
			"MERGE INTO S.V USING DUAL ON (1 = 1) WHEN MATCHED THEN UPDATE SET RATE = :1",
		} {
			if err := periksaBacaSaja(objek, q); !errors.Is(err, ErrMasterBacaSaja) {
				t.Errorf("%s: %q harus ditolak: %v", objek, q, err)
			}
		}
	}
	if err := periksaBacaSaja(TabelBusiness, "UPDATE S.T SET RIRATE = :1 WHERE ID = :2"); err != nil {
		t.Errorf("tabel warisan modul ini boleh ditulis: %v", err)
	}
}

// polaPanggilMaster - pembaca/penyusun SQL yang menyebut objek master: `siapkan(MasterX, sqlY` dan
// `bacaMaster(ctx, MasterX, sqlY`.
var polaPanggilMaster = regexp.MustCompile(`(?:siapkan|bacaMaster)\((?:ctx,\s*)?(Master\w+),\s*(sql\w+)`)

// sqlMasterBukanSelect - fungsi SQL yang dipasangkan dengan objek master tetapi tidak berawal SELECT,
// atau memuat `SELECT *` / `JSONDATA`.
func sqlMasterBukanSelect(berkas map[string]string) []string {
	semua := strings.Join(func() []string {
		var d []string
		for _, isi := range berkas {
			d = append(d, buangKomentarGo(isi))
		}
		return d
	}(), "\n")
	var hasil []string
	for _, m := range polaPanggilMaster.FindAllStringSubmatch(semua, -1) {
		i := strings.Index(semua, "\nfunc "+m[2]+"(")
		if i < 0 {
			hasil = append(hasil, m[2]+" (tidak ditemukan)")
			continue
		}
		badan := semua[i+1:]
		if j := strings.Index(badan, "\n}"); j >= 0 {
			badan = badan[:j]
		}
		awal := strings.Index(badan, "`")
		if awal < 0 {
			hasil = append(hasil, m[2]+" (tanpa teks SQL)")
			continue
		}
		q := strings.ToUpper(strings.TrimSpace(badan[awal+1:]))
		if !strings.HasPrefix(q, "SELECT ") || strings.Contains(q, "SELECT *") || strings.Contains(q, "JSONDATA") {
			hasil = append(hasil, m[1]+"/"+m[2])
		}
	}
	return hasil
}

func TestMCRLMasterHanyaDibacaSelect(t *testing.T) {
	berkas := map[string]string{}
	for jalur, isi := range berkasModul(t, ".go") {
		if strings.Contains(jalur, "/repository/") && !strings.HasSuffix(jalur, "_test.go") {
			berkas[jalur] = isi
		}
	}
	if n := len(polaPanggilMaster.FindAllString(strings.Join(func() []string {
		var d []string
		for _, isi := range berkas {
			d = append(d, isi)
		}
		return d
	}(), "\n"), -1)); n < 8 {
		t.Fatalf("hanya %d pemanggilan master terbaca - pembacanya yang rusak", n)
	}
	if bad := sqlMasterBukanSelect(berkas); len(bad) > 0 {
		t.Errorf("SQL master/view rate yang bukan SELECT kolom: %v", bad)
	}
}

func TestMCRLAturanMasterSelectMenggigit(t *testing.T) {
	isi := "\nfunc (g *Gudang) x(ctx context.Context) {\n\tg.bacaMaster(ctx, MasterRate, sqlUbahRate, nil)\n}\n" +
		"\nfunc sqlUbahRate(t string) string {\n\treturn fmt.Sprintf(`UPDATE %s SET RATE = :1`, t)\n}\n" +
		"\nfunc (g *Gudang) y(ctx context.Context) {\n\tg.siapkan(MasterRingkasanRate, sqlSemuaRate)\n}\n" +
		"\nfunc sqlSemuaRate(t string) string {\n\treturn fmt.Sprintf(`SELECT * FROM %s`, t)\n}\n"
	if bad := sqlMasterBukanSelect(map[string]string{"x.go": isi}); len(bad) != 2 {
		t.Errorf("UPDATE ke view dan SELECT * harus tertangkap: %v", bad)
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
