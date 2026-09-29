package repository

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/internal/models"
	"nusantarare/inti/db"
)

// AC 24: SATU tabel untuk seluruh jenis - setiap SQL klausul memakai tabel yang
// sama, dan migrasi modul tidak membuat tabel klausul per jenis.
func TestKlausulSatuTabel(t *testing.T) {
	for nama, q := range map[string]string{
		"daftar": sqlDaftarKlausulTCO("S.PROPORTIONALARRG"), "ambil": sqlAmbilKlausulTCO("S.PROPORTIONALARRG"),
		"induk": sqlIndukKlausulTCO("S.PROPORTIONALARRG"), "sisip": sqlSisipKlausulTCO("S.PROPORTIONALARRG"),
		"perbarui": sqlPerbaruiKlausulTCO("S.PROPORTIONALARRG"), "pct": sqlPctAnakLainTCO("S.PROPORTIONALARRG"),
	} {
		if strings.Count(q, "S.PROPORTIONALARRG") != 1 {
			t.Errorf("%s tidak memakai tepat satu tabel klausul: %s", nama, q)
		}
	}
	berkas, _ := filepath.Glob(filepath.Join("migrations", "3*.sql"))
	pola := regexp.MustCompile(`CREATE TABLE \{skema\}\.(\w+)`)
	for _, b := range berkas {
		isi, _ := os.ReadFile(b)
		for _, m := range pola.FindAllStringSubmatch(string(isi), -1) {
			nama := strings.ToUpper(m[1])
			for _, jenis := range []string{"EPI", "PLA", "RICOMM", "PROFIT", "BORDER", "EXGRATIA", "FACIN", "COINS", "LOL", "PORTFOLIO"} {
				if strings.Contains(nama, jenis) {
					t.Errorf("tabel %s di %s tampak tabel per jenis klausul - AC 24 menuntut satu tabel", nama, b)
				}
			}
		}
	}
}

// AC 25: baris ANAK menyimpan NULL pada sembilan kolom khusus induk, meski
// struct-nya (keliru) berisi.
func TestKlausulAnakMenulisNullKolomInduk(t *testing.T) {
	isi := models.KlausulTreaty{ID: "10000001", ParentReinsTypeID: "10003", IDOccupation: "X", Occupation: "X",
		IDClause: "X", Clause: "X", TreatyLimit: apd.New(1, 0), CoinsMin: apd.New(1, 0), CoinsMax: apd.New(1, 0),
		MoreRp: apd.New(1, 0), MoreUsd: apd.New(1, 0), Pct: apd.New(5, 0)}
	v := nilaiKolomKlausul(isi)
	for _, k := range KolomKhususIndukTCO {
		if v[k] != nil {
			t.Errorf("anak menulis %s = %v", k, v[k])
		}
	}
	if len(KolomKhususIndukTCO) != 9 || v["PCT"] != "5" {
		t.Errorf("kolom khusus induk %d, PCT %v", len(KolomKhususIndukTCO), v["PCT"])
	}
	isi.ParentReinsTypeID = models.ParentReinsTypeTanpaInduk
	if nilaiKolomKlausul(isi)["TREATYLIMIT"] == nil {
		t.Error("induk kehilangan TREATYLIMIT")
	}
}

// 35 kolom VERBATIM urutan prosedur; sisip ber-35 placeholder; perbarui tanpa
// mengubah kunci; daftar ID ASC (AC 31) dengan lingkup RD.
func TestSQLKlausulTCO(t *testing.T) {
	if len(kolomKlausulTCO) != 35 {
		t.Fatalf("kolom %d, mau 35", len(kolomKlausulTCO))
	}
	if !strings.Contains(sqlSisipKlausulTCO("S.T"), ":35)") {
		t.Error("sisip tanpa 35 placeholder")
	}
	p := sqlPerbaruiKlausulTCO("S.T")
	for _, tetap := range []string{"TREATYYEARID =", "TREATYDESCID =", "PARENTREINSTYPEID =", " ID ="} {
		if strings.Contains(p[:strings.Index(p, "WHERE")], tetap) {
			t.Errorf("perbarui mengubah kunci %q", tetap)
		}
	}
	if !strings.Contains(p, "WHERE ID = :32 AND TREATYYEARID = :33") {
		t.Errorf("perbarui: %s", p)
	}
	d := sqlDaftarKlausulTCO("S.T")
	if !strings.Contains(d, "TREATYYEARID = :1 AND TREATYDESCID = :2 AND PARENTREINSTYPEID = :3 ORDER BY ID ASC") {
		t.Errorf("daftar: %s", d)
	}
	q := sqlCariDobelKlausulTCO("S.T")
	if !strings.Contains(q, "TREATYYEARID = :1 AND TREATYDESCID = :2 AND PARENTREINSTYPEID = :3") ||
		!strings.Contains(q, "(:4 IS NULL OR ID <> :5)") {
		t.Errorf("dobel: %s", q)
	}
	// tco4: kolom desimal TEKS warisan tidak pernah dibungkus TO_CHAR berformat angka.
	pilih := pilihKlausulTCO()
	for _, teks := range []string{"RP", "USD", "PCT", "PCTME", "KURS"} {
		if strings.Contains(pilih, "TO_CHAR("+teks+",") {
			t.Errorf("TO_CHAR atas kolom VARCHAR2 %s", teks)
		}
	}
	for _, angka := range []string{"TREATYLIMIT", "COINS_MIN", "COINS_MAX", "MORERP", "MOREUSD"} {
		if !strings.Contains(pilih, "TO_CHAR("+angka+",") {
			t.Errorf("kolom NUMBER %s tidak dibaca TO_CHAR ber-NLS", angka)
		}
	}
	for _, s := range []string{d, p, q, sqlSisipKlausulTCO("S.T"), sqlJenisKlausulTCO("S.D"),
		sqlCariPilihanTCO("S.O", "NAME"), sqlAmbilPilihanTCO("S.C", "INFO"), sqlKunciTahunTCO("S.Y")} {
		if err := db.PeriksaSQL(s); err != nil {
			t.Errorf("PeriksaSQL: %v", err)
		}
	}
}

// AC 27: master TREATYDESC tidak pernah ditulis - di sumber produksi mana pun.
func TestTREATYDESCTidakDitulis(t *testing.T) {
	pola := regexp.MustCompile(`(?i)(INSERT\s+INTO|UPDATE|DELETE\s+FROM|MERGE\s+INTO)\s+[^\s]*TREATYDESC\b`)
	for nama, isi := range berkasSumberProduksi(t) {
		if pola.MatchString(isi) {
			t.Errorf("%s menulis TREATYDESC - master jenis klausul dibaca saja (AC 27)", nama)
		}
	}
	if !strings.Contains(sqlJenisKlausulTCO("S.D"), "SELECT ID, DESCNAME, ISXOL, STATUSAKTIF") {
		t.Error("master jenis tidak membaca kolom RD")
	}
}

func TestPindaiKlausulTCO(t *testing.T) {
	nilai := make([]any, 35)
	for i := range nilai {
		nilai[i] = ""
	}
	// tco4: PCT/RP teks warisan - koma diterima.
	nilai[0], nilai[5], nilai[17], nilai[22], nilai[24], nilai[14] = "10000001", "10009", "12,5", "00", "1000000.5", "2026-09-29 10:00:00"
	k, err := pindaiKlausulTCO(barisPalsu{nilai: nilai})
	if err != nil {
		t.Fatal(err)
	}
	if k.ID != "10000001" || k.TreatyDescID != "10009" || k.Pct.Text('f') != "12.5" || k.Rp.Text('f') != "1000000.5" ||
		k.ParentReinsTypeID != "00" || k.Usd != nil || k.TglUpdate.Hour() != 10 {
		t.Errorf("pindai: %+v", k)
	}
}

// Kontrak hilir: setiap kolom yang dibaca Claim Prop / Claim Fac In
// (`KolomKlausulHilir`) dan saringan induk `PARENTREINSTYPEID` ditulis oleh
// penulis klausul - hilir tidak membaca kolom yang tidak pernah terisi.
func TestKlausulMenulisKolomHilir(t *testing.T) {
	tulis := map[string]bool{}
	for _, k := range kolomKlausulTCO {
		tulis[k] = true
	}
	for _, k := range append(append([]string{}, KolomKlausulHilir...), "PARENTREINSTYPEID") {
		if !tulis[k] {
			t.Errorf("kolom hilir %s tidak ditulis penulis klausul", k)
		}
	}
}

// tco4: kunci dobel dibandingkan di Go - nilai desimal sama walau teksnya beda,
// kosong tidak pernah sama, medan di luar daftar putih ditolak.
func TestMedanSamaKlausulTCO(t *testing.T) {
	d := func(teks string) *apd.Decimal {
		v, _, _ := UraiDesimalWarisanTCO(teks)
		return v
	}
	a := models.KlausulTreaty{Pct: d("12.5"), Layer: "UJI-L1"}
	b := models.KlausulTreaty{Pct: d("12,50"), Layer: " UJI-L1 "}
	for medan, mau := range map[string]bool{models.MedanPct: true, models.MedanLayer: true, models.MedanRp: false,
		models.MedanMethod: false} {
		sama, err := medanSamaKlausulTCO(a, b, medan)
		if err != nil || sama != mau {
			t.Errorf("%s: %v %v, mau %v", medan, sama, err, mau)
		}
	}
	if _, err := medanSamaKlausulTCO(a, b, "JSONDATA"); err == nil {
		t.Error("kunci dobel di luar daftar putih diterima")
	}
}
