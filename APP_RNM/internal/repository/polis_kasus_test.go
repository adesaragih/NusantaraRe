package repository

// Uji kasus polis baru - GILIRAN-13 paket 1 (pl3, bn). TANPA Oracle.

import (
	"strings"
	"testing"

	"nusantarare/internal/models"
)

// TestPengenalWorkPolisTanpaNolDepan - bentuk pengenal warisan, dari data.
//
// ⛔ `[terverifikasi]` sampel baca-saja DEV (29-09-2026, `ROWNUM <= 200`,
// bentuk saja): `JSON_OFFER_LIFE` dan `M_LIFE_PREMIUM_SUMMARY` 400/400
// berbentuk `ASM-FW-GISFW-WORK NBLF-<1..5 digit>`, NOL berawalan nol. Berbeda
// dengan klaim (`CLM-` dipad enam digit, butir aa): polis mengikuti datanya.
func TestPengenalWorkPolisTanpaNolDepan(t *testing.T) {
	for urut, mau := range map[string]string{
		"7": "NBLF-7", "12345": "NBLF-12345", " 42 ": "NBLF-42", "1000000": "NBLF-1000000",
	} {
		if got := RakitPengenalWorkPolis(urut); got != mau {
			t.Errorf("urut %q: %q, mau %q", urut, got, mau)
		}
	}
	if len(RakitPengenalWorkPolis("999999999999")) > 32 {
		t.Error("pengenal melampaui T_WORK_POLIS.ID VARCHAR2(32)")
	}
}

func TestSisipKasusPolisMenulisKelimaKolomKerja(t *testing.T) {
	q := sqlSisipKasusPolis(tabelUjiWorkPolis)
	if err := PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"ID", "LINI", "POSITION", "STATUS", "FLAG_ONGOING_POLICY"} {
		if !strings.Contains(q, k) {
			t.Errorf("kolom %s tidak ditulis:\n%s", k, q)
		}
	}
	if got := urutanPenampung(q); got != "12345" {
		t.Errorf("penampung %q, mau 12345", got)
	}
}

func TestSisipPremiumListKosongBerbagiPengenal(t *testing.T) {
	q := sqlSisipPremiumListKosong("SKEMAUJI.T_PREMIUM_LIST")
	if err := PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	// ⛔ ID = ID_PEGA = pengenal work: shared PK (050/051), dan kotak masuk
	// menggabung `p.ID_PEGA = w.ID` - tanpa ID_PEGA kasus baru tampil tanpa
	// tanggal dan tanpa urutan.
	// Penampung UNIK (pengenal dikirim dua kali): penampung berulang terbukti
	// aman hanya untuk SELECT tanpa pembatas baris - INSERT tidak pernah diuji.
	for _, k := range []string{"(ID, ID_PEGA, TGL_INPUT)", "VALUES (:1, :2, :3)"} {
		if !strings.Contains(q, k) {
			t.Errorf("bentuk %q tidak ada:\n%s", k, q)
		}
	}
	// Tidak ada kolom isian lain: layar Input Offer yang mengisinya.
	for _, k := range []string{"TYPE", "NO_POLIS", "BUSINESS_CODE", "PL_NUMBER"} {
		if strings.Contains(q, k) {
			t.Errorf("kolom isian %s ditulis saat kasus lahir:\n%s", k, q)
		}
	}
}

func TestArgumenKasusPolisBerurutanSepertiPenampung(t *testing.T) {
	k, err := models.SusunKasusPolisBaru(models.FlagPolisPremium)
	if err != nil {
		t.Fatal(err)
	}
	arg := argSisipKasusPolis("NBLF-1", k)
	mau := []any{"NBLF-1", models.LiniLife, models.PosisiOffer, models.TahapPolisPenawaran, "1"}
	if len(arg) != len(mau) {
		t.Fatalf("argumen %v", arg)
	}
	for i := range mau {
		if arg[i] != mau[i] {
			t.Errorf("argumen %d = %v, mau %v", i+1, arg[i], mau[i])
		}
	}
}

// Migrasi 057 - butir bn.
func TestMigrasi057SequenceDanBendera(t *testing.T) {
	var maju, mundur string
	for nama, teks := range seluruhSQL(t, false) {
		if strings.Contains(nama, "057_") {
			maju = strings.ToUpper(teks)
		}
	}
	for nama, teks := range seluruhSQL(t, true) {
		if strings.Contains(nama, "057_") {
			mundur = strings.ToUpper(teks)
		}
	}
	for _, mau := range []string{
		"CREATE SEQUENCE {SKEMA}.SEQ_WORK_POLIS",
		"ALTER TABLE {SKEMA}.T_WORK_POLIS ADD (",
		"FLAG_ONGOING_POLICY VARCHAR2(1)",
	} {
		if !strings.Contains(maju, mau) {
			t.Errorf("057 tidak memuat %q", mau)
		}
	}
	for _, mau := range []string{
		"ALTER TABLE {SKEMA}.T_WORK_POLIS DROP COLUMN FLAG_ONGOING_POLICY",
		"DROP SEQUENCE {SKEMA}.SEQ_WORK_POLIS",
	} {
		if !strings.Contains(mundur, mau) {
			t.Errorf("057 mundur tidak memuat %q", mau)
		}
	}
}
