package repository

// Uji teks SQL khas EDM (tanpa Oracle): saringan daftar EDM lawan NB, penampung urut = argumen, kunci generasi
// endorsemen, dan kolom popup Retro tanpa OLDID. Bentuk yang sama dicoba baca-saja di DEV 06-10-2026
// (`scratchpad/edmx/dev-sql-baca.txt`).

import (
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/edmtreatyin/backend/models"
)

var rePenampung = regexp.MustCompile(`:(\d+)`)

// penampungUrut - penampung `:n` muncul berurut 1..n dan jumlahnya = jumlah argumen.
func penampungUrut(t *testing.T, q string, args []any) {
	t.Helper()
	terbesar := 0
	for _, m := range rePenampung.FindAllStringSubmatch(q, -1) {
		n := 0
		for _, c := range m[1] {
			n = n*10 + int(c-'0')
		}
		if n > terbesar+1 {
			t.Fatalf("penampung :%d melompat (terbesar sebelumnya :%d)\n%s", n, terbesar, q)
		}
		if n > terbesar {
			terbesar = n
		}
	}
	if terbesar != len(args) {
		t.Fatalf("penampung terbesar :%d, argumen %d\n%s", terbesar, len(args), q)
	}
}

func TestDaftarEDMHanyaGenerasiEndorsemen(t *testing.T) {
	q, args := sqlDaftarKasus("P.W", "P.G", "P.Q", models.SaringanKasus{Cari: "UJI-POL", Pembuat: "UJI-AKUN"})
	for _, w := range []string{"g.PRODKE >= 1", "w.ID LIKE 'EDMT-%'", "UPPER(q.OLD_POLICY_NO) LIKE", "w.CREATE_OP = ", "ORDER BY w.TGL_UPDATE DESC", "FETCH FIRST 500 ROWS ONLY", "NVL(g.NOPOLIS, q.OLD_POLICY_NO)", "TO_CHAR(g.TGL_PROD"} {
		if !strings.Contains(q, w) {
			t.Errorf("daftar portal tanpa %q\n%s", w, q)
		}
	}
	if strings.Contains(q, "PRODKE = 0") {
		t.Fatal("daftar EDM tidak boleh menyaring PRODKE = 0 (milik NB)")
	}
	penampungUrut(t, q, args)
	if args[len(args)-2] != "%UJI-POL%" || args[len(args)-1] != "UJI-AKUN" {
		t.Fatalf("argumen cari / pembuat = %v", args)
	}
}

// Perintah work owner 07-10-2026 ("pencarian ... buat bisa mencari nomor nb/edm, insured name dll"): setiap kata kotak
// saring wajib cocok dengan SALAH SATU kolom portal (AND antar-kata, OR antar-kolom), tanpa beda huruf, % _ \ harfiah.
func TestCariPortalBanyakKolomPerKata(t *testing.T) {
	q, args := sqlDaftarKasus("P.W", "P.G", "P.Q", models.SaringanKasus{Cari: `  uji-pol   50%_x\  `})
	penampungUrut(t, q, args)
	kolom := []string{"w.ID", "g.NO_OFFER", "g.NOPOLIS", "q.OLD_POLICY_NO", "g.NOENDORS", "g.INSURED_NAME", "q.INSURED_NAME",
		"q.BUSINESS_NAME", "g.SOB_NAME", "g.CEDING_CO_NAME", "q.MARKETING_NAME", "g.TREATY_GROUP_NAME", "g.BIZ_NAME",
		"w.CREATE_OP_NAME"}
	for _, k := range kolom {
		if n := strings.Count(q, "UPPER("+k+") LIKE :"); n != 2 {
			t.Errorf("kolom %s dicari %d kali, harap 2 (satu per kata)\n%s", k, n, q)
		}
	}
	if n := strings.Count(q, `ESCAPE '\'`); n != 2*len(kolom) {
		t.Errorf("ESCAPE %d, harap %d", n, 2*len(kolom))
	}
	if strings.Count(q, "\n\t    AND (UPPER(") != 2 {
		t.Errorf("kata kedua wajib AND, kolom OR\n%s", q)
	}
	// argumen: 3 tetap + satu per kolom per kata; huruf besar, wildcard di-escape
	if len(args) != 3+2*len(kolom) || args[3] != "%UJI-POL%" || args[3+len(kolom)] != `%50\%\_X\\%` {
		t.Fatalf("argumen cari %v", args)
	}
	q, args = sqlDaftarKasus("P.W", "P.G", "P.Q", models.SaringanKasus{Cari: "a b c d e f g"})
	penampungUrut(t, q, args)
	if len(args) != 3+models.MaksKataCari*len(kolom) {
		t.Fatalf("kata terbanyak %d, argumen %d", models.MaksKataCari, len(args))
	}
}

func TestDaftarEDMKotakMasukDanSelesai(t *testing.T) {
	q, args := sqlDaftarKasus("P.W", "P.G", "P.Q", models.SaringanKasus{Pembuat: "UJI", PembuatPosisi: models.PosisiAdmin,
		Antrean: []string{models.PosisiSecHead}})
	penampungUrut(t, q, args)
	if !strings.Contains(q, "(w.CREATE_OP = :4 AND g.POSITION_NOTE = :5) OR g.POSITION_NOTE IN (:6)") {
		t.Fatalf("kotak masuk: pembuat ATAU antrean\n%s", q)
	}
	q, args = sqlDaftarKasus("P.W", "P.G", "P.Q", models.SaringanKasus{Selesai: true})
	penampungUrut(t, q, args)
	if strings.Contains(q, "NOT IN") || !strings.Contains(q, "w.STATUS_WORK IN (:2, :3)") {
		t.Fatalf("switch Resolved = HANYA berkas selesai\n%s", q)
	}
}

// Salinan dokumen Pega berkunci IDPEGA UTUH (WO 07-10-2026 "IDPEGA BAWAAN PEGA JANGAN DI POTONG"): daftar portal, kotak
// masuk, dan cek EDM terbuka mengenali `EDMT-<n>` dan `<kelas> EDMT-<n>`.
func TestIDKasusEDMMencakupPzInsKeyUtuh(t *testing.T) {
	harap := "(w.ID LIKE 'EDMT-%' OR w.ID LIKE '% EDMT-%')"
	if sqlIDKasusEDM("w.ID") != harap {
		t.Fatalf("syarat ID kasus %q", sqlIDKasusEDM("w.ID"))
	}
	q, _ := sqlDaftarKasus("P.W", "P.G", "P.Q", models.SaringanKasus{})
	k, _ := sqlHitungKotakMasuk("P.W", "P.G", "UJI", true, nil)
	for _, x := range []string{q, k} {
		if !strings.Contains(x, harap) {
			t.Errorf("tanpa syarat ID kasus utuh\n%s", x)
		}
	}
}

func TestKotakMasukEDMSaringPRODKE(t *testing.T) {
	q, args := sqlHitungKotakMasuk("P.W", "P.G", "UJI", true, []string{models.PosisiSecHead, models.PosisiDeptHead})
	penampungUrut(t, q, args)
	if !strings.Contains(q, "g.PRODKE >= 1") || !strings.Contains(q, "w.ID LIKE 'EDMT-%'") {
		t.Fatalf("kotak masuk EDM tanpa saringan generasi endorsemen\n%s", q)
	}
}

func TestSisipGenerasiMembawaKunciEndorsemen(t *testing.T) {
	q := sqlSisipGenerasi("P.G")
	for _, k := range []string{"PRODKE", "NOENDORS", "OLD_POLIS_ID", "EDM_TYPE", "POSITION_NOTE"} {
		if !strings.Contains(q, k) {
			t.Errorf("sisip generasi tanpa %s (AC 1, 2, 15)\n%s", k, q)
		}
	}
	if strings.Contains(q, "NOPOLIS") {
		t.Fatal("NOPOLIS generasi endorsemen KOSONG sampai selesai (SetelNomorPolisSelesai)")
	}
	penampungUrut(t, q, make([]any, 7))
}

func TestKeadaanHanyaGenerasiEndorsemen(t *testing.T) {
	q := sqlKeadaan("P.W", "P.G")
	if !strings.Contains(q, "g.PRODKE >= 1") || !strings.Contains(q, "g.OLD_POLIS_ID") {
		t.Fatalf("keadaan EDM\n%s", q)
	}
}

func TestPopupRetroTanpaKolomOldID(t *testing.T) {
	p := pilihTanpaOldID("A")
	if strings.Contains(p, "A.OLDID") || !strings.HasPrefix(p, "TO_CHAR(A.ID), NULL") {
		t.Fatalf("RDB TreatyLoadMasterJoinEdmChooseBusinessRetro: OLDID = (select null from dual); TREATY_OUT2 tanpa OLDID\n%s", p)
	}
	if awalanKarakter("1234567890123", 7) != "1234567" || awalanKarakter("12", 7) != "12" {
		t.Fatal("@substring(.., 0, n)")
	}
}

// Copy Old (perintah work owner 07-10-2026): data lama = JSON_POLIS x tabel kerja Pega (PZINSKEY = IDPEGA) x
// TREATYINPRODUCTION (IDPEGA), lewat EXISTS; generasi endorsemen saja; baca saja.
func TestKunciCopyOldGabungKerjaPegaDanProduksi(t *testing.T) {
	q := sqlPmKunciJSONPolisEDMCopyOld("P.JSON_POLIS", "P.TREATYINPRODUCTION")
	for _, w := range []string{"FROM P.JSON_POLIS b", "EXISTS (SELECT 1 FROM DATAPEGA.PC_ASM_FW_GISFW_WORK a WHERE a.PZINSKEY = b.IDPEGA)",
		"EXISTS (SELECT 1 FROM P.TREATYINPRODUCTION c WHERE c.IDPEGA = b.IDPEGA)", "TRIM(TO_CHAR(b.PRODKE)) <> '0'"} {
		if !strings.Contains(q, w) {
			t.Errorf("tanpa %q", w)
		}
	}
	penampungUrut(t, q, nil)
}

// WO 07-10-2026 "PXCREATEOPERATOR,PXCREATEOPNAME": pembuat berkas salinan dibaca dari tabel kerja Pega menurut
// pzInsKey = IDPEGA; baca saja, satu penampung.
func TestPembuatPegaDariTabelKerjaPega(t *testing.T) {
	q := sqlPmPembuatPega()
	if q != `SELECT PXCREATEOPERATOR, PXCREATEOPNAME FROM DATAPEGA.PC_ASM_FW_GISFW_WORK WHERE PZINSKEY = :1` {
		t.Fatalf("pembuat Pega: %s", q)
	}
}
