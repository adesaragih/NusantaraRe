package repository

// Uji pemindah data Treaty Contract Out - bagian yang TANPA Oracle:
// bentuk SQL dan rekonsiliasi murni (tiket 01).

import (
	"strings"
	"testing"
)

func TestSQLBacaWarisanTCOMembungkusMenurutTipe(t *testing.T) {
	q := sqlBacaWarisanTCO("S.PROPORTIONALARRG", warisanKlausulTCO)
	for _, mau := range []string{
		"TO_CHAR(TREATYLIMIT, 'TM9'", "TO_CHAR(MOREUSD, 'TM9'",
		"TO_CHAR(TGLUPDATE, 'YYYY-MM-DD HH24:MI:SS')", " RP,", " USD,", "ORDER BY ID",
	} {
		if !strings.Contains(q, mau) {
			t.Errorf("SQL klausul tanpa %q:\n%s", mau, q)
		}
	}
	// RP/USD/PCT/PCTME VARCHAR2 di warisan: TIDAK dibungkus TO_CHAR angka.
	if strings.Contains(q, "TO_CHAR(RP") || strings.Contains(q, "TO_CHAR(PCT") {
		t.Errorf("kolom teks warisan dibungkus TO_CHAR angka:\n%s", q)
	}
	// Kolom mati TIDAK dibaca (AC 70).
	for _, mati := range kolomMatiKlausulTCO {
		if strings.Contains(q, mati) {
			t.Errorf("kolom mati %s ikut dibaca:\n%s", mati, q)
		}
	}

	s := sqlBacaWarisanTCO("S.MTREATYSECURITY", warisanSecurityTCO)
	for _, mau := range []string{"RTRIM(REAS_ID)", "RTRIM(REAS_SECURITY)", "ORDER BY ROWID"} {
		if !strings.Contains(s, mau) {
			t.Errorf("SQL security tanpa %q:\n%s", mau, s)
		}
	}
}

func TestSQLSisipBaruTCOBernamaDanIDPertama(t *testing.T) {
	q := sqlSisipBaruTCO("S.T_MTREATYSECURITY", KolomWarisanTCO(warisanSecurityTCO))
	if !strings.HasPrefix(q, "INSERT INTO S.T_MTREATYSECURITY (ID, THN_TREATY, TOP_ID, TP_TREATY, REAS_ID, PCT_SHARE, USER_ID, REAS_SECURITY) VALUES (:1, :2, :3, :4, :5, :6, :7, :8)") {
		t.Errorf("INSERT security tidak bernama/berurut:\n%s", q)
	}
	k := sqlSisipBaruTCO("S.T_PROPORTIONALARRG", KolomWarisanTCO(warisanKlausulTCO))
	if strings.Count(k, ":") != 35 || strings.Count(k, "ID,") < 1 {
		t.Errorf("INSERT klausul harus 35 penampung dengan ID sekali:\n%s", k)
	}
}

func TestSQLBacaBaruTCOMembungkusTujuan(t *testing.T) {
	q := sqlBacaBaruTCO("S.T_TREATYREINSURER", warisanReinsurerTCO)
	for _, mau := range []string{"SELECT ID, TREATYYEAR", "TO_CHAR(PCTSHARE, 'TM9'",
		"TO_CHAR(STARTDATE, 'YYYY-MM-DD HH24:MI:SS')", " IUDATE,"} {
		if !strings.Contains(q, mau) {
			t.Errorf("SQL baca kembali tanpa %q:\n%s", mau, q)
		}
	}
}

func kembaliDari(siap []BarisSiapTCO, ubah func(tabel string, baris map[string]string)) map[string][]BarisWarisanTCO {
	kembali := map[string][]BarisWarisanTCO{}
	for i, s := range siap {
		baris := map[string]string{"ID": s.ID}
		if s.Warisan == warisanSecurityTCO {
			baris["ID"] = "1000" + string(rune('0'+i))
		}
		for k, v := range s.Kanonik {
			baris[k] = v
		}
		if ubah != nil {
			ubah(s.Tabel, baris)
		}
		kembali[s.Tabel] = append(kembali[s.Tabel], BarisWarisanTCO{Tabel: s.Tabel, Nilai: baris})
	}
	return kembali
}

func TestSelisihRekonsiliasiTCOBersihBilaSama(t *testing.T) {
	siap, _ := KonversiWarisanTCO(contohWarisanTCO())
	kembali := kembaliDari(siap, func(tabel string, baris map[string]string) {
		// TO_CHAR TM9 menulis "12.5" dan ".5" - bentuk berbeda, nilai sama.
		if tabel == TabelKlausulTCO {
			baris["USD"] = ".5"
			baris["RP"] = "1000000000.12345678"
		}
		if tabel == TabelReinsurerTCO {
			baris["PCTSHARE"] = "33.333000"
		}
	})
	if selisih := SelisihRekonsiliasiTCO(siap, kembali); len(selisih) != 0 {
		t.Errorf("selisih palsu:\n%s", strings.Join(selisih, "\n"))
	}
}

func TestSelisihRekonsiliasiTCOMenemukanPerubahanSatuDigit(t *testing.T) {
	siap, _ := KonversiWarisanTCO(contohWarisanTCO())
	kembali := kembaliDari(siap, func(tabel string, baris map[string]string) {
		switch tabel {
		case TabelKlausulTCO:
			baris["RP"] = "1000000000.12345679" // satu digit terakhir berubah
		case TabelTahunTCO:
			baris["STARTDATE"] = "2026-01-02 00:00:00" // bergeser satu hari
		case TabelSecurityTCO:
			baris["REAS_SECURITY"] = "SEC-B"
		}
	})
	selisih := SelisihRekonsiliasiTCO(siap, kembali)
	if len(selisih) != 4 {
		t.Fatalf("selisih %d, mau 4 (RP, STARTDATE, security hilang, security asing):\n%s",
			len(selisih), strings.Join(selisih, "\n"))
	}
	gabung := strings.Join(selisih, "\n")
	for _, mau := range []string{"T_PROPORTIONALARRG/10000009.RP", "T_TREATYYEAR/1000001.STARTDATE",
		"tidak terbaca kembali", "tidak dikonversi"} {
		if !strings.Contains(gabung, mau) {
			t.Errorf("selisih tanpa %q:\n%s", mau, gabung)
		}
	}
}

func TestSelisihRekonsiliasiTCOMenemukanBarisHilangDanAsing(t *testing.T) {
	siap, _ := KonversiWarisanTCO(contohWarisanTCO())
	kembali := kembaliDari(siap, nil)
	kembali[TabelBusinessTCO] = nil
	kembali[TabelTahunTCO] = append(kembali[TabelTahunTCO],
		BarisWarisanTCO{Tabel: TabelTahunTCO, Nilai: map[string]string{"ID": "1000777"}})
	selisih := strings.Join(SelisihRekonsiliasiTCO(siap, kembali), "\n")
	for _, mau := range []string{"T_TREATYBUSINESS: ditulis 1 baris, terbaca kembali 0",
		"T_TREATYBUSINESS/1000002: tidak terbaca kembali", "T_TREATYYEAR/1000777: ada di tabel baru"} {
		if !strings.Contains(selisih, mau) {
			t.Errorf("selisih tanpa %q:\n%s", mau, selisih)
		}
	}
}

func TestNilaiBindTCO(t *testing.T) {
	if nilaiBindTCO(nil) != nil {
		t.Error("nil harus nil")
	}
	siap, _ := KonversiWarisanTCO(contohWarisanTCO())
	klausul := cari(siap, warisanKlausulTCO)
	if v := nilaiBindTCO(klausul.Nilai["RP"]); v != "1000000000.12345678" {
		t.Errorf("desimal harus menyeberang sebagai teks, dapat %v", v)
	}
	if v := nilaiBindTCO(klausul.Nilai["TREATYLIMIT"]); v != nil {
		t.Errorf("desimal kosong harus nil, dapat %v", v)
	}
}
