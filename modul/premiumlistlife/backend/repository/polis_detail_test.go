package repository

// Grid Premium List Detail - tiket 03 PremiumList Life. TANPA Oracle.
//
// ⛔ SEBAB PENJAGA UTAMA DI BERKAS INI ADA. Kotak masuk tiket 01 hampir
// menampilkan `KetentuanUnderwriting`, yang tidak punya kolom di migrasi mana
// pun; `PL_Detail_Sec` menampilkan `REINSTYPENAME` dan `RetrocadedShare`, yang
// juga tidak. Kolom layar yang tidak punya kolom tabel tidak gagal saat
// dirakit, tidak gagal saat dikirim, dan hanya terlihat sebagai ORA-00904 di
// Oracle sungguhan - yaitu di mesin work owner, bukan di sini.

import (
	"fmt"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/migrasi"
	"nusantarare/modul/premiumlistlife/backend/models"
)

// kolomTabelPeserta membaca nama kolom migrasi 052 dengan pengurai PRODUKSI.
//
// ⚠️ `KolomCreateTable`, bukan regex kedua yang ditulis khusus untuk uji ini.
// Pengurai uji yang berbeda dari pengurai produksi menguji pengurai uji.
func kolomTabelPeserta(t *testing.T) map[string]bool {
	t.Helper()
	isi, err := berkasMigrasi.ReadFile("migrations/052_t_premium_list_detail.sql")
	if err != nil {
		t.Fatalf("membaca migrasi 052: %v", err)
	}
	ada := map[string]bool{}
	for _, pernyataan := range strings.Split(string(isi), "\n/") {
		nama, kolom := migrasi.KolomCreateTable(pernyataan)
		if !strings.HasSuffix(strings.ToUpper(nama), "T_PREMIUM_LIST_DETAIL") {
			continue
		}
		for _, k := range kolom {
			ada[strings.ToUpper(k)] = true
		}
	}
	if len(ada) == 0 {
		t.Fatal("nol kolom terbaca dari migrasi 052; pembacanya yang rusak")
	}
	return ada
}

// TestKolomGridPesertaAdaDiMigrasi052 - setiap kolom grid punya kolom tabel.
func TestKolomGridPesertaAdaDiMigrasi052(t *testing.T) {
	ada := kolomTabelPeserta(t)
	for _, k := range models.KolomGridTampil() {
		if !ada[strings.ToUpper(k.Nama)] {
			t.Errorf("kolom grid %q tidak ada di migrasi 052.\n"+
				"Kolom layar tanpa kolom tabel baru gagal di Oracle sungguhan "+
				"(ORA-00904) - daftarkan di models.MedanGridTanpaKolom beserta "+
				"alasannya, atau perbaiki namanya.", k.Nama)
		}
	}
	// ID dipakai sebagai kunci baris dan ikut dipindai.
	if !ada["ID"] {
		t.Error("kolom ID tidak ada di migrasi 052")
	}
}

// TestPenjagaKolomGridMasihMenggigit - penjaga yang tidak pernah cocok adalah
// penjaga yang selalu hijau.
func TestPenjagaKolomGridMasihMenggigit(t *testing.T) {
	ada := kolomTabelPeserta(t)
	for _, hantu := range []string{"REINSTYPENAME", "RETROCADEDSHARE", "KETENTUANUNDERWRITING"} {
		if ada[hantu] {
			t.Errorf("kolom %q ternyata ADA di migrasi 052 - "+
				"catatan di models/polis_detail.go harus diperbarui", hantu)
		}
	}
}

// TestMedanTanpaKolomDinyatakan - selisihnya dijawab, bukan dilupakan.
//
// `PL_Detail_Sec.xml` menampilkan empat puluh medan; grid kami tiga puluh
// delapan kolom. Selisih dua itu harus punya nama dan alasan.
func TestMedanTanpaKolomDinyatakan(t *testing.T) {
	// `[terverifikasi]` 40. Jendela: elemen `<pyValue>` ber-titik yang UNIK di
	// SATU berkas, `PremiumList Life/Section/PL_Detail_Sec.xml`. Dihitung dua
	// cara, keduanya menjawab 40:
	//
	//	grep -o "<pyValue>\.[A-Za-z0-9_.]*</pyValue>" Section/PL_Detail_Sec.xml \
	//	  | sed 's/<[^>]*>//g' | sort -u | wc -l
	//	grep -o "<pyValue>\.[A-Za-z0-9_.]*</pyValue>" Section/PL_Detail_Sec.xml \
	//	  | sed 's/<[^>]*>//g' | awk '!s[$0]++' | wc -l
	//
	// ⚠️ Kemunculan MENTAHnya 109 - satu medan muncul berkali-kali. Sensus
	// yang mencacah kemunculan akan menjawab 109 dan menyimpulkan gridnya
	// hampir tiga kali lebih lebar daripada yang ada.
	const mauTampil = 40
	if n := len(models.KolomGridPeserta) + len(models.MedanGridTanpaKolom); n != mauTampil {
		t.Errorf("kolom grid + medan tanpa kolom = %d, mau %d "+
			"(cacah <pyValue> UNIK di PL_Detail_Sec.xml).\n"+
			"Kalau korpusnya memang berubah, ubah angka ini bersama alasannya "+
			"DAN hitung ulang dua cara - lihat perintah audit di atas.",
			n, mauTampil)
	}
	for medan, alasan := range models.MedanGridTanpaKolom {
		if strings.TrimSpace(alasan) == "" {
			t.Errorf("medan %q didaftar tanpa alasan tertulis", medan)
		}
	}
}

// TestNamaTertanggungTidakAdaDiGrid - nol nama orang.
//
// ⛔ `NAME_OF_INSURED` ADA di migrasi 052 dan TIDAK ditampilkan Pega di grid
// ini - ia hidup di `ShowLifePremiumDetail`. Ia nama orang: tidak pernah masuk
// fixture, tiket, atau log, dan grid yang membawanya membawanya ke ekspor
// xlsx juga.
func TestNamaTertanggungTidakAdaDiGrid(t *testing.T) {
	// ⚠️ Yang dijaga TIRUAN PL_Detail_Sec saja. Sejak 02-10-2026 NAME_OF_INSURED
	// tampil lewat models.KolomGridTambahan (keputusan work owner) - lihat
	// TestKolomGridTambahanTerpisahDariTiruan.
	for _, k := range models.KolomGridPeserta {
		if strings.EqualFold(k.Nama, "NAME_OF_INSURED") {
			t.Error("NAME_OF_INSURED ada di grid peserta; Pega pun tidak " +
				"menampilkannya di PL_Detail_Sec")
		}
	}
}

// TestUrutanKolomGridMengikutiPLDetailSec - urutan bagian dari tiruan.
func TestUrutanKolomGridMengikutiPLDetailSec(t *testing.T) {
	if len(models.KolomGridPeserta) == 0 {
		t.Fatal("daftar kolom grid kosong")
	}
	nama := models.NamaKolomGridPeserta()
	if nama[0] != "BEGIN_DATE" {
		t.Errorf("kolom pertama %q, mau BEGIN_DATE (pyValue pertama di PL_Detail_Sec)", nama[0])
	}
	if akhir := nama[len(nama)-1]; akhir != "NET_PREMIUM_RETRO" {
		t.Errorf("kolom terakhir %q, mau NET_PREMIUM_RETRO", akhir)
	}
	// Nol kolom kembar: satu kolom dua kali berarti satu kolom lain hilang.
	lihat := map[string]bool{}
	for _, n := range nama {
		if lihat[n] {
			t.Errorf("kolom %q muncul dua kali", n)
		}
		lihat[n] = true
	}
}

// TestKolomUangGridDibungkusTM9 - nol float di jalur pulang.
func TestKolomUangGridDibungkusTM9(t *testing.T) {
	ekspresi := ekspresiKolomPeserta(models.KolomGridPeserta)
	for _, k := range models.KolomGridPeserta {
		switch k.Jenis {
		case models.KolomPesertaAngka:
			mau := fmt.Sprintf(db.FmtDesimal, "d."+k.Nama)
			if !strings.Contains(ekspresi, mau) {
				t.Errorf("kolom angka %q tidak dibungkus TM9", k.Nama)
			}
		case models.KolomPesertaTanggal:
			mau := fmt.Sprintf(db.FmtTanggalOracle, "d."+k.Nama)
			if !strings.Contains(ekspresi, mau) {
				t.Errorf("kolom tanggal %q tidak dibungkus TO_CHAR berpola", k.Nama)
			}
		}
	}
	// ⛔ FACTOR adalah DESIMAL tujuh angka (AC 39 spec), bukan bilangan bulat.
	if !strings.Contains(ekspresi, fmt.Sprintf(db.FmtDesimal, "d.FACTOR")) {
		t.Error("FACTOR tidak dibaca sebagai desimal")
	}
}

// TestQueryGridPesertaBerbatasDanTerurut - halaman yang dapat dipercaya.
func TestQueryGridPesertaBerbatasDanTerurut(t *testing.T) {
	q := sqlGridPeserta("SKEMAUJI.T_PREMIUM_LIST_DETAIL", models.KolomGridPeserta)
	for _, potong := range []string{
		"ORDER BY d.ID\n",
		"OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY",
		"WHERE d.PREMIUM_LIST_ID = :1",
	} {
		if !strings.Contains(q, potong) {
			t.Errorf("query grid tidak memuat %q:\n%s", potong, q)
		}
	}
	if err := db.PeriksaSQL(q); err != nil {
		t.Errorf("%v\n%s", err, q)
	}
	if err := db.PeriksaSQL(sqlCacahPeserta("SKEMAUJI.T_PREMIUM_LIST_DETAIL")); err != nil {
		t.Error(err)
	}
}

// TestBatasUkuranHalamanPesertaMenjepit - satu polis grup ribuan peserta.
func TestBatasUkuranHalamanPesertaMenjepit(t *testing.T) {
	for _, k := range []struct{ minta, mau int }{
		{0, 50}, {-5, 50}, {1, 1}, {50, 50}, {500, 500}, {501, 500}, {1000000, 500},
	} {
		if got := BatasUkuranHalamanPeserta(k.minta); got != k.mau {
			t.Errorf("ukuran %d -> %d, mau %d", k.minta, got, k.mau)
		}
	}
}

// TestKolomGridTambahanTerpisahDariTiruan - empat kolom tambahan work owner
// (02-10-2026) dibaca dan dikirim ke layar, TANPA mengubah tiruan
// PL_Detail_Sec, dan dengan pembungkus jenisnya.
func TestKolomGridTambahanTerpisahDariTiruan(t *testing.T) {
	tampil := models.NamaKolomGridTampil()
	for _, k := range models.KolomGridTambahan {
		for _, p := range models.KolomGridPeserta {
			if p.Nama == k.Nama {
				t.Errorf("kolom tambahan %q juga ada di tiruan PL_Detail_Sec", k.Nama)
			}
		}
		found := false
		for _, n := range tampil {
			found = found || n == k.Nama
		}
		if !found {
			t.Errorf("kolom tambahan %q tidak dikirim ke layar", k.Nama)
		}
	}
	if len(tampil) != len(models.KolomGridPeserta)+len(models.KolomGridTambahan) {
		t.Errorf("kolom tampil %d, mau tiruan + tambahan", len(tampil))
	}
	ekspresi := ekspresiKolomPeserta(models.KolomGridTampil())
	for _, kolom := range []string{"d.DOB", "d.GROSS_VALUATION_BEGIN_DATE", "d.GROSS_VALUATION_EXPIRED_DATE"} {
		if !strings.Contains(ekspresi, fmt.Sprintf(db.FmtTanggalOracle, kolom)) {
			t.Errorf("%s tidak dibungkus TO_CHAR berpola tanggal", kolom)
		}
	}
}
