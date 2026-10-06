package repository

// Penyimpanan peserta hasil unggahan - tiket 04. TANPA Oracle.

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/premiumlistlife/backend/models"
)

// TestKolomSisipPesertaAdaDiMigrasi052 - nol kolom hantu.
//
// ⛔ Kolom yang tidak ada di DDL baru gagal di Oracle sungguhan (ORA-00904),
// yaitu di mesin work owner. Penjaga yang sama sudah menyelamatkan grid
// peserta tiket 03.
func TestKolomSisipPesertaAdaDiMigrasi052(t *testing.T) {
	ada := kolomTabelPeserta(t)
	periksa := func(nama string) {
		if !ada[strings.ToUpper(nama)] {
			t.Errorf("kolom sisip %q tidak ada di migrasi 052", nama)
		}
	}
	periksa("ID")
	periksa("PREMIUM_LIST_ID")
	for _, c := range kolomSisipPeserta {
		periksa(c.Kolom)
	}
	for _, c := range kolomTeksTanggalPeserta {
		periksa(c.Kolom)
	}
	for _, c := range kolomTanggalPeserta {
		periksa(c.Kolom)
	}
	for _, c := range kolomBulatPeserta {
		periksa(c.Kolom)
	}
	// ⛔ Ketiga puluh dua kolom uang unggahan juga harus punya rumah.
	for _, c := range models.KolomUangTersimpan() {
		periksa(c)
	}
}

// TestNolKolomSisipKembar - satu kolom dua kali adalah INSERT yang gagal.
func TestNolKolomSisipKembar(t *testing.T) {
	kolom, _, cacah := ekspresiSisipPeserta()
	bagian := strings.Split(kolom, ", ")
	if len(bagian) != cacah {
		t.Fatalf("%d nama kolom tetapi %d parameter", len(bagian), cacah)
	}
	lihat := map[string]bool{}
	for _, k := range bagian {
		if lihat[k] {
			t.Errorf("kolom %q disisipkan dua kali", k)
		}
		lihat[k] = true
	}
}

// TestNilaiSisipSejajarDenganKolomnya - urutan argumen SAMA dengan kolomnya.
//
// ⛔ Ini kekeliruan yang tidak satu pun galat tunjukkan: nilai `SUM_INSURED`
// yang mendarat di kolom `SUM_REASURED` tetap angka yang sah, tetap
// tersimpan, dan baru terlihat saat seseorang membandingkan totalnya.
func TestNilaiSisipSejajarDenganKolomnya(t *testing.T) {
	_, _, cacah := ekspresiSisipPeserta()
	arg := nilaiSisipPeserta("ID1", "POLIS1", models.BarisUnggah{
		Nomor: 1, Nilai: map[string]string{},
	})
	if len(arg) != cacah {
		t.Errorf("%d argumen tetapi %d parameter di query", len(arg), cacah)
	}
	if arg[0] != "ID1" || arg[1] != "POLIS1" {
		t.Errorf("dua argumen pertama %v, mau ID lalu PREMIUM_LIST_ID", arg[:2])
	}
	// Kolom kosong menjadi NULL, bukan teks kosong - dan bukan "0".
	for i, a := range arg[2:] {
		if a != nil {
			t.Errorf("argumen ke-%d dari baris kosong bernilai %v, mau nil", i+2, a)
		}
	}
}

// TestUangDikirimSebagaiTeks - nol float di jalur simpan.
func TestUangDikirimSebagaiTeks(t *testing.T) {
	const asli = "12345678901234567890.12345678"
	arg := nilaiSisipPeserta("ID1", "POLIS1", models.BarisUnggah{
		Nomor: 1, Nilai: map[string]string{"SUM_INSURED": asli},
	})
	found := false
	for _, a := range arg {
		if s, ok := a.(string); ok && s == asli {
			found = true
		}
	}
	if !found {
		t.Errorf("nilai uang %q tidak sampai sebagai teks apa adanya", asli)
	}
}

// TestTanggalDibungkusTODATEBerpola - bentuk DINYATAKAN, bukan dari sesi.
func TestTanggalDibungkusTODATEBerpola(t *testing.T) {
	_, penanda, _ := ekspresiSisipPeserta()
	if !strings.Contains(penanda, "TO_DATE(:") {
		t.Error("kolom tanggal tidak dibungkus TO_DATE")
	}
	if !strings.Contains(penanda, BentukTanggalOracle) {
		t.Errorf("pola tanggal %q tidak muncul di query", BentukTanggalOracle)
	}
	// ⛔ STNC dan WPC BUKAN kolom DATE di migrasi 052 - keduanya
	// VARCHAR2(255). Membungkusnya TO_DATE lalu menyimpan ke kolom teks
	// menyerahkan bentuk simpannya kepada sesi.
	n := strings.Count(penanda, "TO_DATE(:")
	if n != len(kolomTanggalPeserta) {
		t.Errorf("%d kolom dibungkus TO_DATE, mau %d", n, len(kolomTanggalPeserta))
	}
}

func TestQueryUnggahDibatasiSatuPolis(t *testing.T) {
	h := sqlHapusPesertaPolis("SKEMAUJI.T_PREMIUM_LIST_DETAIL")
	if !strings.Contains(h, "WHERE PREMIUM_LIST_ID = :1") {
		t.Errorf("hapus tidak dibatasi satu polis:\n%s", h)
	}
	if err := db.PeriksaSQL(h); err != nil {
		t.Errorf("%v\n%s", err, h)
	}
	s := sqlSisipPeserta("SKEMAUJI.T_PREMIUM_LIST_DETAIL")
	if err := db.PeriksaSQL(s); err != nil {
		t.Errorf("%v\n%s", err, s)
	}
	for _, q := range []string{h, s} {
		if strings.Contains(strings.ToUpper(q), "COMMIT") {
			t.Errorf("query memuat COMMIT:\n%s", q)
		}
	}
}

// TestPengenalPesertaTetapDanSepanjangKolomnya.
//
// ⛔ Kolomnya `VARCHAR2(32)`. Pengenal 33 karakter ditolak Oracle dengan
// ORA-12899 di baris yang tidak seorang pun tebak.
func TestPengenalPesertaTetapDanSepanjangKolomnya(t *testing.T) {
	a := PengenalPesertaUnggah("POLIS-1", 1)
	if len(a) != 32 {
		t.Errorf("pengenal %q berpanjang %d, mau 32", a, len(a))
	}
	// DETERMINISTIK: unggah ulang menghasilkan pengenal yang sama.
	if b := PengenalPesertaUnggah("POLIS-1", 1); b != a {
		t.Errorf("pengenal tidak tetap: %q lalu %q", a, b)
	}
	// Baris berbeda, polis berbeda: pengenal berbeda.
	if PengenalPesertaUnggah("POLIS-1", 2) == a {
		t.Error("dua baris berbeda berpengenal sama")
	}
	if PengenalPesertaUnggah("POLIS-2", 1) == a {
		t.Error("dua polis berbeda berpengenal sama")
	}
	// ⛔ Penyambungnya tidak boleh membuat (polis, baris) bertabrakan:
	// "POLIS-1" baris 11 dan "POLIS-11" baris 1 adalah dua hal berbeda.
	if PengenalPesertaUnggah("POLIS-1", 11) == PengenalPesertaUnggah("POLIS-11", 1) {
		t.Error("penyambung pengenal membuat dua pasangan berbeda bertabrakan")
	}
	// Heksa huruf besar, sejajar dengan IMAGEID.
	if strings.ToUpper(a) != a {
		t.Errorf("pengenal %q bukan huruf besar", a)
	}
}

// TestPengenalPesertaBerurutBarisCSV - keputusan work owner 03-10-2026:
// `ORDER BY ID` dalam satu polis = urutan baris berkas CSV.
func TestPengenalPesertaBerurutBarisCSV(t *testing.T) {
	var lalu string
	for _, n := range []int{1, 2, 9, 10, 11, 99, 100, 1000, 20000} {
		id := PengenalPesertaUnggah("NBLF-1", n)
		if id <= lalu {
			t.Errorf("baris %d berpengenal %q, tidak sesudah %q", n, id, lalu)
		}
		lalu = id
	}
	if PengenalPesertaUnggah("NBLF-1", 5)[:24] != PengenalPesertaUnggah("NBLF-1", 6)[:24] {
		t.Error("awalan pengenal berbeda di dalam satu polis")
	}
}
