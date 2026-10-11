package services_test

// Uji panel Attachment - NOL koneksi Oracle.
//
// Angka dalam berkas ini dari pengukuran 4 Oktober 2026 atas
// `POOLDATA.M_ATTACHMENTTREATY_2` (43 baris) dan sapuan korpus kedua modul.

import (
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// Katalog yang TERBACA dari data - ketujuh pasangan yang terpakai.
func katalogUji() map[string]string {
	return map[string]string{
		"00000": "Others",
		"00001": "Analysed Email",
		"00002": "Approval Email",
		"00005": "Summary Treaty Leader",
		"00006": "Assessment Inward Treaty Form / Format Analisa Treaty",
		"00007": "Pega Proportional Calculation /Perhitungan Pega Proportional",
		"00010": "Offer Email",
	}
}

// ⛔ UJI PALING PENTING DI BERKAS INI.
//
// Keempat nama yang kodenya TIDAK diketahui tidak boleh diberi kode. Keempat
// kode tanpa nama dan keempat nama tanpa kode sama-sama empat dan berurutan,
// dan memasangkannya menurut abjad adalah tebakan yang terlihat benar sampai
// seseorang mengunduh berkas dari kategori yang salah bertahun kemudian.
//
// ⭐ RALAT 6 Oktober 2026: panel kini dirender dari DAFTAR DESIGN (kesebelas
// nama, dipastikan pemilik proses lewat tangkapan layar), bukan dari katalog.
// Yang berubah arah penyusunannya; yang TIDAK berubah larangan menebak
// pasangan kode-nama — nama tanpa kode tetap bernilai kode KOSONG.
func TestEmpatNamaTanpaKodeTIDAKDitebakKodenya(t *testing.T) {
	baris := services.SusunKategoriLampiran(models.NamaKategoriLampiranProp, katalogUji(), nil)

	tanpaKode := map[string]models.BarisKategoriLampiran{}
	for _, b := range baris {
		if !b.Dipastikan {
			tanpaKode[b.Nama] = b
		}
	}
	if len(tanpaKode) != 4 {
		t.Fatalf("%d kategori tanpa kode, mau 4: %v", len(tanpaKode), tanpaKode)
	}
	for _, nama := range models.NamaKategoriBelumDipastikan {
		b, ada := tanpaKode[nama]
		if !ada {
			t.Errorf("%q tidak ditandai tanpa kode", nama)
			continue
		}
		if b.Kode != "" {
			t.Errorf("%q diberi kode %q — pasangannya TIDAK diketahui, dan "+
				"menebaknya menaruh berkas di kategori yang salah", nama, b.Kode)
		}
	}
	// Dan sebaliknya: keempat kode yatim tidak menempel ke nama mana pun.
	for _, b := range baris {
		for _, kode := range models.KodeKategoriBelumDipastikan {
			if b.Kode == kode {
				t.Errorf("%q diberi kode %s yang namanya belum diketahui", b.Nama, kode)
			}
		}
	}
}

// Ketujuh pasangan yang terbukti tetap dibaca DARI DATA, bukan dihafal.
func TestTujuhKategoriTerbuktiDibacaDariKatalog(t *testing.T) {
	baris := services.SusunKategoriLampiran(models.NamaKategoriLampiranProp, katalogUji(), nil)

	n := 0
	for _, b := range baris {
		if b.Dipastikan {
			n++
			if b.Kode == "" {
				t.Errorf("%q dipastikan tetapi tanpa kode", b.Nama)
			}
		}
	}
	if n != 7 {
		t.Errorf("%d kategori berkode, mau 7", n)
	}
}

// ⭐ KESEBELAS BARIS TAMPIL, berurut sesuai design — termasuk yang nol berkas.
//
// Panel lama memperlihatkan kesebelasnya beserta `Count 0`-nya. Katalog hanya
// memuat TUJUH yang pernah dipakai; panel yang dirender dari katalog
// kehilangan empat baris, dan kehilangan itu DIAM.
func TestKesebelasKategoriTampilBerurutSesuaiDesign(t *testing.T) {
	baris := services.SusunKategoriLampiran(models.NamaKategoriLampiranProp, katalogUji(), nil)
	if len(baris) != 11 {
		t.Fatalf("%d baris kategori, mau 11", len(baris))
	}
	for i, nama := range models.NamaKategoriLampiranProp {
		if baris[i].Nama != nama {
			t.Errorf("baris %d berbunyi %q, mau %q — urutannya urutan layar lama",
				i, baris[i].Nama, nama)
		}
		if baris[i].Nama == "" {
			t.Errorf("baris %d tanpa nama; panel menampilkan nama, bukan kode", i)
		}
	}
}

// Cabang non-proporsional berbeda PADA SATU butir saja, dan butir itu nyata.
func TestCabangNonPropBerbedaSatuButir(t *testing.T) {
	prop := models.NamaKategoriLampiranProp
	non := models.NamaKategoriLampiranNonProp
	if len(prop) != len(non) {
		t.Fatalf("cacah cabang berbeda: %d lawan %d", len(prop), len(non))
	}
	beda := 0
	for i := range prop {
		if prop[i] != non[i] {
			beda++
		}
	}
	if beda != 1 {
		t.Errorf("%d butir berbeda antar cabang, mau tepat 1", beda)
	}
}

// Kategori yang ADA DI DATA tetapi di luar daftar design tetap TAMPIL —
// menyembunyikannya membuat berkasnya tidak terjangkau siapa pun.
func TestKategoriAsingTetapTampil(t *testing.T) {
	lampiran := []models.BarisLampiranWarisan{
		{KodeKategori: "00099", NamaKategori: "Kategori Yang Tidak Digambar", NamaBerkas: "a.pdf"},
	}
	baris := services.SusunKategoriLampiran(models.NamaKategoriLampiranProp, katalogUji(), lampiran)
	if len(baris) != 12 {
		t.Fatalf("%d baris, mau 12 (11 design + 1 asing)", len(baris))
	}
	akhir := baris[len(baris)-1]
	if akhir.Nama != "Kategori Yang Tidak Digambar" || akhir.Cacah != 1 {
		t.Errorf("kategori asing tidak tampil apa adanya: %+v", akhir)
	}
}

// Cacah dihitung menurut NAMA, sehingga kategori yang kodenya belum
// dipastikan tetap memperlihatkan berkasnya.
func TestCacahMenurutNamaBukanKode(t *testing.T) {
	lampiran := []models.BarisLampiranWarisan{
		{KodeKategori: "", NamaKategori: "Claim Data", NamaBerkas: "a.pdf"},
		{KodeKategori: "", NamaKategori: "Claim Data", NamaBerkas: "b.pdf"},
	}
	baris := services.SusunKategoriLampiran(models.NamaKategoriLampiranProp, katalogUji(), lampiran)
	for _, b := range baris {
		if b.Nama == "Claim Data" {
			if b.Cacah != 2 {
				t.Errorf("Claim Data bercacah %d, mau 2 — kodenya belum dipastikan, "+
					"dan mencacah menurut kode akan memberi nol yang salah", b.Cacah)
			}
			return
		}
	}
	t.Fatal("Claim Data tidak ada di daftar")
}

// ⛔ Kategori BERNOL berkas tetap tampil. Panel lama punya kolom `Count`,
// dan kolom itu tidak akan pernah berbunyi `0` kalau barisnya disembunyikan
// saat kosong.
func TestKategoriNolBerkasTetapTampil(t *testing.T) {
	lampiran := []models.BarisLampiranWarisan{
		{KodeKategori: "00001", NamaKategori: "Analysed Email"},
		{KodeKategori: "00001", NamaKategori: "Analysed Email"},
		{KodeKategori: "00010", NamaKategori: "Offer Email"},
	}
	baris := services.SusunKategoriLampiran(models.NamaKategoriLampiranProp, katalogUji(), lampiran)
	if len(baris) != 11 {
		t.Fatalf("%d baris, mau 11", len(baris))
	}
	cacah := map[string]int{}
	for _, b := range baris {
		cacah[b.Nama] = b.Cacah
	}
	if cacah["Analysed Email"] != 2 || cacah["Offer Email"] != 1 {
		t.Errorf("cacah salah: %v", cacah)
	}
	for _, nama := range []string{"Others", "Approval Email", "Summary Treaty Leader", "Claim Data"} {
		if cacah[nama] != 0 {
			t.Errorf("%q cacahnya %d, mau 0", nama, cacah[nama])
		}
	}
}

// ⭐ RALAT 6 Oktober 2026: urutannya kini urutan DESIGN (alfabetis menurut
// nama), bukan urutan kode.
//
// Sebabnya berubah bersama arah penyusunannya: panel dirender dari daftar
// nama, dan keempat nama yang kodenya belum dipastikan tidak punya kode untuk
// diurutkan. Tangkapan layar pemilik proses memperlihatkan urutan alfabetis.
func TestKategoriBerurutMenurutNamaDesign(t *testing.T) {
	baris := services.SusunKategoriLampiran(models.NamaKategoriLampiranProp, katalogUji(), nil)
	for i := 1; i < len(baris); i++ {
		if baris[i-1].Nama >= baris[i].Nama {
			t.Fatalf("urutan rusak di %d: %q lalu %q", i, baris[i-1].Nama, baris[i].Nama)
		}
	}
	if baris[0].Nama != "Analysed Email" || baris[len(baris)-1].Nama != "Summary Treaty Leader" {
		t.Errorf("ujungnya %q dan %q", baris[0].Nama, baris[len(baris)-1].Nama)
	}
}

// Bila kode yang tadinya tidak diketahui ternyata MUNCUL di data lengkap
// dengan namanya, ia berhenti ditandai - pertanyaannya terjawab sendiri.
func TestKodeYangTernyataTerbacaBerhentiDitandai(t *testing.T) {
	k := katalogUji()
	k["00003"] = "Claim Data"
	baris := services.SusunKategoriLampiran(models.NamaKategoriLampiranProp, k, nil)
	for _, b := range baris {
		if b.Kode == "00003" {
			if !b.Dipastikan || b.Nama != "Claim Data" {
				t.Errorf("00003 = %+v; ia terbaca dari data dan harus dipastikan", b)
			}
			return
		}
	}
	t.Error("00003 hilang")
}

// ⭐ Spanduk biru itu ATURAN. Diuji penolakannya, bukan hanya penerimaannya.
func TestNamaBerkasAmanMenolakYangTidakAman(t *testing.T) {
	for _, n := range []string{
		"laporan treaty.pdf",  // spasi
		"laporan(1).pdf",      // kurung
		"a/b.pdf", "a\\b.pdf", // pemisah folder
		"../rahasia.pdf", // keluar folder
		"laporan#2.pdf", "surat&lampiran.pdf", "data%20.pdf",
		"", "   ", ".", "..",
	} {
		if services.NamaBerkasAman(n) {
			t.Errorf("%q DITERIMA; spanduk menyebut titik atau garis bawah sebagai "+
				"pengganti yang aman, dan nama ini bukan keduanya", n)
		}
	}
}

func TestNamaBerkasAmanMenerimaYangAman(t *testing.T) {
	// ⛔ Pasangannya wajib: pemeriksa yang menolak SEGALANYA lulus uji di
	// atas dengan sempurna, dan menutup pintu yang sistem lama buka.
	for _, n := range []string{
		"laporan.pdf", "laporan_treaty.pdf", "laporan-treaty.pdf",
		"laporan.2026.pdf", "LAPORAN.PDF", "a1.docx", "Treaty_In_2026-01.xlsx",
	} {
		if !services.NamaBerkasAman(n) {
			t.Errorf("%q DITOLAK; ia hanya memuat huruf, angka, titik, garis bawah, "+
				"dan tanda hubung", n)
		}
	}
}
