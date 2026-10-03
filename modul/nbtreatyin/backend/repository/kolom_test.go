package repository

// Uji murni lapisan repository - tanpa Oracle: konversi kolom (ID-14..18),
// bentuk SQL, dan kesepakatan katalog <-> DDL migrasi.

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestPecahAngkaEksak(t *testing.T) {
	for _, tt := range []struct {
		masuk, koef string
		skala       int64
	}{
		{"0", "0", 0},
		{"12.5", "125", 1},
		{"-0.43052837564", "-43052837564", 11},
		{"830.82191780804", "83082191780804", 11},
		{"1000", "1000", 0},
	} {
		d, err := utils.ParseDecimal(tt.masuk)
		if err != nil {
			t.Fatal(err)
		}
		k, s := pecahAngka(d)
		if k != tt.koef || s != tt.skala {
			t.Errorf("%s -> (%s, %d), harap (%s, %d)", tt.masuk, k, s, tt.koef, tt.skala)
		}
	}
}

func TestNilaiTulisKosongJadiNULL(t *testing.T) { // spec-penyimpanan AC 17, 18
	uang := models.Kolom{Properti: "X", Kolom: "X", Golongan: models.GolUang}
	if v, err := nilaiTulis(uang, ""); err != nil || v[0] != nil || v[1] != nil {
		t.Fatalf("uang kosong harus NULL, dapat %v %v", v, err)
	}
	tgl := models.Kolom{Properti: "T", Kolom: "T", Golongan: models.GolTanggal}
	if v, _ := nilaiTulis(tgl, ""); v[0] != nil {
		t.Fatal("tanggal kosong harus NULL")
	}
	kode := models.Kolom{Properti: "K", Kolom: "K", Golongan: models.GolKode, Panjang: 16}
	if v, _ := nilaiTulis(kode, "006"); v[0] != "006" {
		t.Fatal("kode tetap teks: nol di depan bermakna (ID-16)")
	}
	penanda := models.Kolom{Properti: "P", Kolom: "P", Golongan: models.GolPenanda, Panjang: 16}
	if v, _ := nilaiTulis(penanda, "0"); v[0] != "0" {
		t.Fatal("penanda \"0\" berbeda dari kosong (ID-17)")
	}
}

func TestCacahBilanganBulat(t *testing.T) { // ID-14
	c := models.Kolom{Properti: "InstallmentNo", Kolom: "INSTALLMENT_NO", Golongan: models.GolCacah}
	if v, err := nilaiTulis(c, "12"); err != nil || v[0] != "12" {
		t.Fatalf("12: %v %v", v, err)
	}
	if _, err := nilaiTulis(c, "1.5"); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Fatalf("1.5 bukan cacah: %v", err)
	}
	if eks, n := ekspresiTulis(c, 3); eks != "TO_NUMBER(:3)" || n != 1 {
		t.Fatalf("ekspresi %q %d", eks, n)
	}
}

func TestNilaiTulisMenolakMasukanRusak(t *testing.T) {
	uang := models.Kolom{Properti: "PremiOgp", Kolom: "P", Golongan: models.GolUang}
	if _, err := nilaiTulis(uang, "12,5"); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Fatalf("angka berkoma harus ditolak sebagai permintaan tidak sah: %v", err)
	}
	tgl := models.Kolom{Properti: "StartDate", Kolom: "T", Golongan: models.GolTanggal}
	if _, err := nilaiTulis(tgl, "31/12/2026"); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Fatalf("format tanggal kedua ditolak (AC 33): %v", err)
	}
	if v, err := nilaiTulis(tgl, "2026-12-31"); err != nil || v[0] != "2026-12-31 00:00:00" {
		t.Fatalf("tanggal satu format: %v %v", v, err)
	}
	teks := models.Kolom{Properti: "Remark", Kolom: "R", Golongan: models.GolTeks, Panjang: 3}
	if _, err := nilaiTulis(teks, "abcd"); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Fatalf("melebihi panjang kolom: %v", err)
	}
}

func TestNilaiBaca(t *testing.T) {
	uang := models.Kolom{Golongan: models.GolPersen}
	for masuk, harap := range map[string]string{".5": "0.5", "-.25": "-0.25", "12.5": "12.5"} {
		if got := nilaiBaca(uang, sql.NullString{String: masuk, Valid: true}); got != harap {
			t.Errorf("%q -> %q, harap %q", masuk, got, harap)
		}
	}
	tgl := models.Kolom{Golongan: models.GolTanggal}
	if got := nilaiBaca(tgl, sql.NullString{String: "2026-10-03 00:00:00", Valid: true}); got != "2026-10-03" {
		t.Fatalf("tanggal saja: %q", got)
	}
	if got := nilaiBaca(uang, sql.NullString{}); got != "" {
		t.Fatal("NULL dibaca kosong")
	}
}

func TestSQLDaftarKasusPenampungUnik(t *testing.T) {
	q := sqlDaftarKasus("S.W", "S.G", "S.Q", true, true)
	pen := regexp.MustCompile(`:(\d+)`).FindAllStringSubmatch(q, -1)
	lihat := map[string]bool{}
	for _, p := range pen {
		if lihat[p[1]] {
			t.Fatalf("penampung :%s berulang di SQL berpembatas baris:\n%s", p[1], q)
		}
		lihat[p[1]] = true
	}
	if len(lihat) != 7 || !strings.Contains(q, "FETCH FIRST 200 ROWS ONLY") {
		t.Fatalf("SQL daftar kasus:\n%s", q)
	}
}

// Setiap kolom katalog ada di CREATE TABLE tabelnya dengan tipe golongannya,
// dan setiap kolom DDL selain kolom kunci/generasi berasal dari katalog.
func TestKatalogSepakatDenganDDL(t *testing.T) {
	berkas, err := filepath.Glob(filepath.Join("..", "migrations", "3*.sql"))
	if err != nil || len(berkas) == 0 {
		t.Fatalf("migrasi tidak terbaca: %v", err)
	}
	ddl := map[string]map[string]string{}
	pola := regexp.MustCompile(`(?s)CREATE TABLE \{skema\}\.(\w+) \((.*?)\n\)`)
	for _, b := range berkas {
		if strings.HasSuffix(b, "_down.sql") {
			continue
		}
		isi, err := os.ReadFile(b)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range pola.FindAllStringSubmatch(string(isi), -1) {
			kol := map[string]string{}
			for _, baris := range strings.Split(m[2], "\n") {
				f := strings.Fields(strings.TrimSuffix(strings.TrimSpace(baris), ","))
				if len(f) >= 2 && f[0] != "CONSTRAINT" {
					kol[f[0]] = f[1]
				}
			}
			ddl[m[1]] = kol
		}
	}
	tipe := func(g models.Golongan) string {
		switch {
		case g.Desimal():
			return "NUMBER(38,8)"
		case g.Tanggal():
			return "DATE"
		case g == models.GolCacah:
			return "NUMBER(10)"
		}
		return "VARCHAR2"
	}
	bukanKatalog := map[string]bool{"ID": true, "POLIS_ID": true, "NOURUT": true, "INSTALMENT_ID": true, "XOL_ID": true,
		"QUOTATION_ID": true, "NOPOLIS": true, "PRODKE": true, "NOENDORS": true, "OLD_POLIS_ID": true, "IDPEGA": true,
		"TGL_INPUT": true, "USERNAME": true}
	for _, tb := range models.SemuaTabel {
		kol, ada := ddl[tb.Nama]
		if !ada {
			t.Errorf("%s tidak dibuat migrasi mana pun", tb.Nama)
			continue
		}
		dariKatalog := map[string]bool{}
		for _, k := range tb.Kolom {
			dariKatalog[k.Kolom] = true
			got, ada := kol[k.Kolom]
			if !ada {
				t.Errorf("%s.%s ada di katalog, tidak di DDL", tb.Nama, k.Kolom)
				continue
			}
			if !strings.HasPrefix(got, tipe(k.Golongan)) {
				t.Errorf("%s.%s bertipe %s, golongan %s menuntut %s", tb.Nama, k.Kolom, got, k.Golongan, tipe(k.Golongan))
			}
		}
		for k := range kol {
			if !dariKatalog[k] && !bukanKatalog[k] {
				t.Errorf("%s.%s ada di DDL, tidak di katalog", tb.Nama, k)
			}
		}
	}
	// AC 64: kolom isApprovedtoDeptHead tidak dibuat; AC 65: nomor surat bukan penanda arah.
	for tb, kol := range ddl {
		for k := range kol {
			if strings.Contains(k, "DEPT_HEAD") || strings.Contains(k, "LETTER") {
				t.Errorf("%s.%s: penanda arah tangga tidak disimpan (AC 64, 65)", tb, k)
			}
		}
	}
}

func TestKolomViewHilangAdalahGalat(t *testing.T) { // AC 89
	tipe := map[string]string{"ID": "VARCHAR2", "LIMITVALUE": "NUMBER", "COMMENCEMENT": "DATE"}
	if _, _, err := pilihKolom("UJI_VIEW", tipe, []string{"ID", "LIMITVALUE", "TREATYID"}); !errors.Is(err, ErrDataKontrakTidakAda) ||
		!strings.Contains(err.Error(), "TREATYID") {
		t.Fatalf("kolom hilang harus galat yang menyebutnya: %v", err)
	}
	eks, ada, err := pilihKolom("UJI_VIEW", tipe, []string{"ID", "LIMITVALUE", "COMMENCEMENT"})
	if err != nil || len(ada) != 3 || !strings.Contains(eks[1], "TM9") || !strings.Contains(eks[2], "TO_CHAR(COMMENCEMENT") || eks[0] != "ID" {
		t.Fatalf("ekspresi menurut tipe: %v %v %v", eks, ada, err)
	}
}

// InsertViewSuggest_SQL: 15 kolom, KETERANGAN dipotong 3990, TGL_INP DATE,
// penampung unik; pembacaan balik berurut NOURUT sebagai bilangan (AC 39-44).
func TestSQLRiwayatProduksiMengikutiInsertViewSuggest(t *testing.T) {
	q := sqlSisipUsulan("S.HISTORYAKSEPTASIPRODUCTION")
	for _, k := range []string{"IDPEGA", "TYPE_POLIS", "NOURUT", "POSISI", "PIC", "TGL_INP", "DIV", "TYPE", "PUTARAN",
		"APPROVAL", "KETERANGAN", "AKSES_LOGIN", "B2B", "BUSINESS_CODE", "PERCENT_RNM"} {
		if !regexp.MustCompile(`[(,\s]` + k + `[,)\s]`).MatchString(q) {
			t.Errorf("kolom %s tidak disisipkan:\n%s", k, q)
		}
	}
	pen := regexp.MustCompile(`:(\d+)`).FindAllStringSubmatch(q, -1)
	if len(pen) != 15 || !strings.Contains(q, "SUBSTR(:11, 1, 3990)") || !strings.Contains(q, "TO_DATE(:6, '"+fmtTanggal+"')") {
		t.Fatalf("SQL sisip riwayat produksi:\n%s", q)
	}
	if strings.Contains(strings.ToUpper(q), "COMMIT") {
		t.Fatal("nol COMMIT - transaksi milik services")
	}
	b := sqlBacaUsulan("S.HISTORYAKSEPTASIPRODUCTION")
	if !strings.Contains(b, "WHERE IDPEGA = :1") || !strings.Contains(b, "ORDER BY TO_NUMBER(NOURUT)") {
		t.Fatalf("SQL baca riwayat produksi:\n%s", b)
	}
	if n := sqlNourutUsulan("S.HISTORYAKSEPTASIPRODUCTION"); !strings.Contains(n, "MAX(TO_NUMBER(NOURUT))") || !strings.Contains(n, "IDPEGA = :1") {
		t.Fatalf("NOURUT berikutnya per IDPEGA:\n%s", n)
	}
}

// kolomDDL membaca CREATE TABLE seluruh migrasi maju modul ini: tabel -> kolom.
func kolomDDL(t *testing.T) map[string][]string {
	t.Helper()
	berkas, err := filepath.Glob(filepath.Join("..", "migrations", "*.sql"))
	if err != nil || len(berkas) == 0 {
		t.Fatalf("migrasi tidak terbaca: %v", err)
	}
	pola := regexp.MustCompile(`(?s)CREATE TABLE \{skema\}\.(\w+) \((.*?)\n\)`)
	hasil := map[string][]string{}
	for _, b := range berkas {
		if strings.HasSuffix(b, "_down.sql") {
			continue
		}
		isi, err := os.ReadFile(b)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range pola.FindAllStringSubmatch(string(isi), -1) {
			if _, ganda := hasil[m[1]]; ganda {
				t.Fatalf("%s dibuat dua kali", m[1])
			}
			var kol []string
			for _, baris := range strings.Split(m[2], "\n") {
				if f := strings.Fields(strings.TrimSpace(baris)); len(f) >= 2 && f[0] != "CONSTRAINT" {
					kol = append(kol, f[0])
				}
			}
			hasil[m[1]] = kol
		}
	}
	return hasil
}

// Bab 0 butir 11-12 PROMPT putaran 2: TEPAT delapan tabel diagram grilling
// (Diagram-Skema-Tabel-NusantaraRe.xlsx sheet NB Treaty In Prop/NonProp), dan
// kolomnya mengikuti diagram + rancangan-tabel-datar; kolom di luar keduanya
// hanya yang XML buktikan DIBACA rule terjangkau (RALAT rancangan,
// docs/PERBANDINGAN-KOLOM-DIAGRAM.md).
func TestTabelDanKolomMengikutiDiagramGrilling(t *testing.T) {
	daftar := func(s string) []string { return strings.Fields(s) }
	angsuran := "INSTALLMENT_NO DUE_DATE INSTALLMENT_PERCENTAGE PREMIUM PAYMENT_TOTAL PREMIUM_AFTER_PPH PREMIUM_AFTER_PPN PREMIUM_AFTER_TAX CURRENCY ID_CURRENCY "
	xol := "CURRENCY ID_CURRENCY GROSS_PREMI NET_PREMI DEDUCTION DUE_TO DUE_TO_VALUE BROKERAGE_FEE_SEBENARNYA PPH_VALUE PPN_VALUE NET_PREMI_AFTER_PPH NET_PREMI_AFTER_PPN NET_PREMI_AFTER_TAX"
	harap := map[string][]string{
		"T_GENERAL_POLIS": daftar(
			// kunci + json_polis (diagram F11-F16)
			"ID NOPOLIS PRODKE NOENDORS OLD_POLIS_ID IDPEGA TGL_INPUT USERNAME TGL_PROD " +
				// PolicyTreatyIn - rancangan §4.1
				"NO_OFFER MASTER_ID IS_APPROVED SUGGEST SUGGEST_DATE OPERATOR_NAME IS_NEW_POLICY_NON_PROP IS_EDM_INPUT_ON_NB " +
				"HAS_FAC_OUT FLAG_PPH FLAG_RETRO_TREATY DUE_TO TYPE_TAX STATEMENT_TYPE TREATY_GROUP_ID TREATY_GROUP_NAME " +
				"TREATY_GROUP_OLD_ID OJK_BUSINESS_ID ID_NEW_BISNIS BIZ_CODE BIZ_NAME SOB SOB_NAME CEDING_CO CEDING_CO_NAME " +
				"INSURED_ID INSURED_NAME MARKETING_OFFICER TREATY_TYPE TREATY_YEAR CURRENCY ID_CURRENCY QUARTAL " +
				"YEAR_OF_QUARTAL CLAIM_TYPE CLAIM_PAYMENT_TYPE INSTALLMENT REMARK START_DATE END_DATE STATEMENT_DATE " +
				"GROSS_PREMIUM PREMI_OGP RESULT_OGP1 RESULT_OGP2 PREMI_ONP RESULT_ONP1 RESULT_ONP2 CLAIM OUTSTANDING_CLAIM " +
				"SALVAGE_VALUE EXCESS_LOSS NET_PREMIUM BALANCE_DUE_TO BALANCE_BEFORE_TAX BALANCE_BEFORE_PPH DEDUCTION1 " +
				"DEDUCTION2 PPH_VALUE PPN_VALUE SHARE_VALUE RI_COMM_OGP OVERIDDING_COMM_OGP RI_COMM_ONP OVERIDDING_COMM_ONP " +
				// RALAT - dibaca rule terjangkau
				"POSITION_NOTE NB_STATUS TREATY_IN_ID SHARE_CURRENCY GROSS_CLAIM BROKERAGE_FEE_SEBENARNYA"),
		"T_POLIS_QUOTATION": daftar("POLIS_ID PROPORTIONAL_TYPE MO_ID BUSINESS_CODE BUSINESS_OLD_ID GROUP_PANEL " +
			"SOURCE_OF_BUSINESS TYPE EDM_TYPE OLD_POLICY_NO MARKETING_NAME " +
			// RALAT - dibaca rule terjangkau / tampil di Section NB
			"BUSINESS_NAME BUSINESS_FAC INSURED_ID INSURED_NAME NO_OFFER_SLIP IS_SURVEY_REPORT"),
		"T_POLIS_CEDING":            daftar("ID QUOTATION_ID NOURUT CEDING_CO_ID CEDING_CO_NAME"),
		"T_POLIS_INSTALMENT":        daftar("ID POLIS_ID NOURUT " + angsuran + "PPN PPH PAYMENT_TOTAL_AFTER_PPN PAYMENT_TOTAL_AFTER_TAX"),
		"T_POLIS_INSTALMENT_DETAIL": daftar("ID INSTALMENT_ID NOURUT " + angsuran + "PAYMENT_DATE"),
		"T_POLIS_SPREADING": daftar("ID POLIS_ID NOURUT TREATY_TYPE TREATY_NAME CURRENCY CURRENCY_ID SHARE_PERCENTAGE " +
			"SPLIT_RNM_SHARE_PCT CLAIM_PERCENTAGE PREMIUM_SPREADED CLAIM_SPREADED"),
		"T_POLIS_XOL":       daftar("ID POLIS_ID NOURUT " + xol),
		"T_POLIS_XOL_LAYER": daftar("ID XOL_ID NOURUT LAYER LAYER_TYPE LAYER_PART LAYER_PART_TYPE " + xol),
	}
	ddl := kolomDDL(t)
	if len(ddl) != 8 {
		var nama []string
		for n := range ddl {
			nama = append(nama, n)
		}
		t.Fatalf("TEPAT delapan CREATE TABLE (diagram grilling), dapat %d: %v", len(ddl), nama)
	}
	for tabel, mau := range harap {
		ada, dibuat := ddl[tabel]
		if !dibuat {
			t.Errorf("%s tidak dibuat", tabel)
			continue
		}
		punya := map[string]bool{}
		for _, k := range ada {
			punya[k] = true
		}
		diminta := map[string]bool{}
		for _, k := range mau {
			diminta[k] = true
			if !punya[k] {
				t.Errorf("%s: kolom %s (diagram/rancangan/RALAT) tidak ada", tabel, k)
			}
		}
		for _, k := range ada {
			if !diminta[k] {
				t.Errorf("%s: kolom %s di luar diagram, rancangan, dan RALAT", tabel, k)
			}
		}
	}
	// Ceding di bawah QUOTATION (diagram O39 "1:N QUOTATION_ID").
	isi, err := os.ReadFile(filepath.Join("..", "migrations", "322_t_polis_ceding.sql"))
	if err != nil || !strings.Contains(string(isi), "FOREIGN KEY (QUOTATION_ID) REFERENCES {skema}.T_POLIS_QUOTATION (POLIS_ID)") {
		t.Fatalf("T_POLIS_CEDING menunjuk T_POLIS_QUOTATION: %v\n%s", err, isi)
	}
}

// Kolom yang dibuang (bab 0 butir 12; PERBANDINGAN-KOLOM-DIAGRAM.md) tidak
// boleh tersisa di SQL mana pun paket ini - ia tidak ada di DDL.
func TestSQLTidakMenyebutKolomYangDibuang(t *testing.T) {
	dibuang := []string{"TGL_TUTUP", "IS_OJK_NOPOLIS", "BUSINESS_TYPE", "SOB_LEADER0", "SOB_LEADER1", "MARKETING_CODE",
		"TEAM_GROUP", "BRANCH_CODE", "BRANCH_NAME", "M_NBTRIN_PERAN_TEMPAT", "T_POLIS_SUGGEST", "T_POLIS_MEDAN_LAIN"}
	berkas, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	dibaca := 0
	for _, b := range berkas {
		if strings.HasSuffix(b, "_test.go") {
			continue
		}
		isi, err := os.ReadFile(b)
		if err != nil {
			t.Fatal(err)
		}
		dibaca++
		for i, baris := range strings.Split(string(isi), "\n") {
			if strings.HasPrefix(strings.TrimSpace(baris), "//") {
				continue
			}
			for _, k := range dibuang {
				if regexp.MustCompile(`\b` + k + `\b`).MatchString(baris) {
					t.Errorf("%s:%d menyebut kolom/tabel yang dibuang %s", b, i+1, k)
				}
			}
		}
	}
	if dibaca < 5 {
		t.Fatalf("hanya %d berkas terbaca", dibaca)
	}
}
