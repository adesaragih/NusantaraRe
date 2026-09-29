package repository

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"nusantarare/internal/models"
)

// AC 44 + AC 63/64: kaskade tidak MENGHAPUS klausul, dan hanya menyentuh
// keluarga tabel T_ (tidak ada tabel kembar JSON / warisan).
func TestKaskadeTidakMenghapusKlausul(t *testing.T) {
	isi, err := os.ReadFile("tco_kaskade.go")
	if err != nil {
		t.Fatal(err)
	}
	kode := buangKomentarSumber("tco_kaskade.go", string(isi))
	for _, m := range regexp.MustCompile(`DELETE FROM %s WHERE[^`+"`"+`]*`).FindAllString(kode, -1) {
		if strings.Contains(m, "PARENTREINSTYPEID") || strings.Contains(m, "TREATYDESCID") {
			t.Errorf("DELETE bersaringan klausul: %s", m)
		}
	}
	// Tabel yang dihapus kaskade: tepat empat keluarga T_, tanpa klausul.
	tb := tabelKaskadeTCO{kontrak: "S.K", reinsurer: "S.R", security: "S.S", business: "S.B", klausul: "S.P", tahun: "S.Y"}
	hapus := []string{sqlHapusKontrakKaskadeTCO(tb.kontrak), sqlHapusKaskadeTCO(tb.business, saringBusinessKaskadeTCO),
		sqlHapusKaskadeTCO(tb.security, fmt.Sprintf(saringSecurityKaskadeTCO, tb.reinsurer)),
		sqlHapusKaskadeTCO(tb.reinsurer, saringReinsurerKaskadeTCO), sqlHapusSecurityReinsurerTCO(tb.security),
		sqlHapusReinsurerSatuTCO(tb.reinsurer)}
	for _, q := range hapus {
		if strings.Contains(q, "S.P") {
			t.Errorf("kaskade menghapus klausul: %s", q)
		}
		if err := PeriksaSQL(q); err != nil {
			t.Errorf("%v: %s", err, q)
		}
	}
	for _, larang := range []string{"M_TREATY", "JSON", "PROPORTIONALARRG"} {
		for _, q := range hapus {
			if strings.Contains(strings.ToUpper(q), larang) {
				t.Errorf("kaskade menyebut %s: %s", larang, q)
			}
		}
	}
}

// Saringan hitung = saringan hapus (angka popup = angka terhapus), dan bisnis
// tahan NULL pada TREATYYEARID - VERBATIM `DeleteFromTREATYCONTRACT_SQL` b85.
func TestKaskadeSaringanSamaDanTahanNull(t *testing.T) {
	if !strings.Contains(saringBusinessKaskadeTCO, "(TREATYYEARID = :2 OR TREATYYEARID IS NULL)") {
		t.Error("bisnis tidak tahan NULL TREATYYEARID")
	}
	for _, s := range []string{saringBusinessKaskadeTCO, saringReinsurerKaskadeTCO, saringSecurityKaskadeTCO} {
		h, d := sqlHitungTCO("S.X", s), sqlHapusKaskadeTCO("S.X", s)
		if h[strings.Index(h, "WHERE"):] != d[strings.Index(d, "WHERE"):] {
			t.Errorf("saringan berbeda:\n%s\n%s", h, d)
		}
	}
	if err := PeriksaSQL(sqlHitungTCO("S.P", saringKlausulTetapTCO)); err != nil {
		t.Error(err)
	}
}

// AC 44 pada langkah yang BENAR-BENAR dijalankan HapusKontrak: tidak satu pun
// menyasar tabel klausul, dan urutannya urutan Pega.
func TestLangkahHapusKontrakTanpaKlausul(t *testing.T) {
	tb := tabelKaskadeTCO{kontrak: "S.K", reinsurer: "S.R", security: "S.S", business: "S.B", klausul: "S.P", tahun: "S.Y"}
	var d DampakHapusTCO
	var urut []string
	for _, l := range langkahHapusKontrakTCO(tb, models.KombinasiTCO{TreatyYear: "2026", TreatyGroupID: "10001", ReinsTypeID: "10003"}, "1000001", "1000003", &d) {
		if l.tabel == tb.klausul || strings.Contains(l.q, tb.klausul) {
			t.Errorf("langkah kaskade menyasar klausul: %s", l.q)
		}
		if l.isi == nil {
			t.Errorf("langkah tanpa penghitung: %s", l.q)
		}
		urut = append(urut, l.tabel)
	}
	if strings.Join(urut, ",") != "S.K,S.B,S.S,S.R" {
		t.Errorf("urutan kaskade: %v", urut)
	}
}

// Temuan /code-review: kombinasi dipakai bersama kontrak lain -> reinsurer,
// security, dan bisnis tanpa TREATYYEARID TIDAK ikut terhapus.
func TestLangkahHapusKontrakBersama(t *testing.T) {
	tb := tabelKaskadeTCO{kontrak: "S.K", reinsurer: "S.R", security: "S.S", business: "S.B", klausul: "S.P", tahun: "S.Y"}
	d := DampakHapusTCO{Bersama: 1}
	var tabel []string
	for _, l := range langkahHapusKontrakTCO(tb, models.KombinasiTCO{TreatyYear: "2026", TreatyGroupID: "10001", ReinsTypeID: "10003"}, "1000001", "1000003", &d) {
		tabel = append(tabel, l.tabel)
		if l.tabel == tb.business && strings.Contains(l.q, "IS NULL") {
			t.Errorf("bisnis milik bersama ikut terhapus: %s", l.q)
		}
	}
	if strings.Join(tabel, ",") != "S.K,S.B" {
		t.Errorf("langkah bersama: %v", tabel)
	}
	if err := PeriksaSQL(sqlKontrakBersamaTCO("S.K", "S.Y")); err != nil {
		t.Error(err)
	}
}

// Temuan /code-review: anti-dobel tahun dan kontrak dikunci.
func TestAntiDobelTahunDanKontrakDikunci(t *testing.T) {
	if q := sqlKunciTabelTahunTCO("S.Y"); q != "LOCK TABLE S.Y IN EXCLUSIVE MODE" || PeriksaSQL(q) != nil {
		t.Errorf("kunci tabel tahun: %s", q)
	}
	isi, _ := os.ReadFile("tco_kontrak.go")
	if !strings.Contains(string(isi), "kunci := sqlKunciTahunTCO(tahun)") {
		t.Error("CariDobel kontrak tidak mengunci tahun induk")
	}
	if err := PeriksaSQL(sqlJumlahAnakTahunTCO("S.K", "S.P")); err != nil {
		t.Error(err)
	}
}
