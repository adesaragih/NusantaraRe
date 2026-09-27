package repository

import (
	"strings"
	"testing"

	"nusantarare/internal/models"
)

// Uji query pencarian diagnosa — kelompok Medis.
//
// ⛔ Tabelnya 97.586 baris. Uji di berkas ini menjaga satu janji di atas
// segalanya: TIDAK ADA jalan keluar dari sini yang tidak berbatas.

const tabelUjiPenyakit = "SKEMAUJI.DISEASE_LIFE"

func TestQueryPenyakitSelaluBerbatas(t *testing.T) {
	// ⛔ SELURUH kombinasi kriteria, termasuk yang kosong. Pencarian kosong
	// SAH di Pega (`Contains ""` cocok dengan semua) dan yang menahannya
	// adalah `pyMaxRecords` 500 - jadi justru kombinasi kosong itu yang
	// paling harus berbatas.
	for _, k := range []models.KriteriaPenyakit{
		{},
		{KodeICD: "A00"},
		{Nama: "DIABETES"},
		{KodeICD: "A00", Nama: "DIABETES"},
	} {
		q, _ := sqlCariPenyakit(tabelUjiPenyakit, k, models.BatasPenyakit(0))
		if !strings.Contains(q, "FETCH FIRST") {
			t.Errorf("kriteria %+v menghasilkan query tanpa batas:\n%s", k, q)
		}
		if !strings.Contains(q, "ORDER BY") {
			t.Errorf("kriteria %+v menghasilkan query tanpa urutan; "+
				"halaman tanpa urutan stabil menampilkan baris yang berbeda "+
				"setiap kali:\n%s", k, q)
		}
	}
}

func TestBatasPenyakitDijepit(t *testing.T) {
	// ⛔ Batas yang diminta pemanggil TIDAK pernah dipercaya apa adanya.
	// `?batas=1000000` atas tabel 97.586 baris adalah permintaan yang memuat
	// seluruh tabel ke memori satu proses.
	if got := models.BatasPenyakit(0); got != models.UkuranHalamanPenyakit {
		t.Errorf("BatasPenyakit(0) = %d, mau %d", got, models.UkuranHalamanPenyakit)
	}
	if got := models.BatasPenyakit(-5); got != models.UkuranHalamanPenyakit {
		t.Errorf("BatasPenyakit(-5) = %d, mau %d", got, models.UkuranHalamanPenyakit)
	}
	if got := models.BatasPenyakit(1000000); got != models.BatasBarisPenyakit {
		t.Errorf("BatasPenyakit(sejuta) = %d, mau %d", got, models.BatasBarisPenyakit)
	}
	if got := models.BatasPenyakit(25); got != 25 {
		t.Errorf("BatasPenyakit(25) = %d, mau 25", got)
	}
	// Dan angkanya dari rule, bukan dari selera kami.
	if models.BatasBarisPenyakit != 500 {
		t.Errorf("BatasBarisPenyakit = %d; `pyMaxRecords` b659 menyebut 500",
			models.BatasBarisPenyakit)
	}
	if models.UkuranHalamanPenyakit != 50 {
		t.Errorf("UkuranHalamanPenyakit = %d; `pyPageSize` b514 menyebut 50",
			models.UkuranHalamanPenyakit)
	}
}

func TestKriteriaPenyakitDisambungAND(t *testing.T) {
	// ⛔ `A AND B` (b535/b754), bukan OR. OR akan mengembalikan seluruh
	// penyakit yang KODEnya cocok ditambah seluruh yang NAMAnya cocok -
	// hasil yang lebih banyak dan lebih tidak berguna, dan bedanya tidak
	// terlihat sampai seseorang mengetik dua kata kunci sekaligus.
	q, arg := sqlCariPenyakit(tabelUjiPenyakit,
		models.KriteriaPenyakit{KodeICD: "A00", Nama: "DIABETES"}, 50)
	if strings.Contains(q, " OR ") {
		t.Errorf("query memakai OR:\n%s", q)
	}
	if !strings.Contains(q, " AND ") {
		t.Errorf("dua kriteria tidak disambung AND:\n%s", q)
	}
	if len(arg) != 2 {
		t.Errorf("argumen = %d, mau 2", len(arg))
	}
	// ⛔ Nilainya lewat BIND, tidak pernah ditempel ke teks query.
	if strings.Contains(q, "A00") || strings.Contains(q, "DIABETES") {
		t.Errorf("nilai kriteria tertempel ke teks query:\n%s", q)
	}
}

func TestKriteriaKosongTidakMenghasilkanKlausa(t *testing.T) {
	// ⚠️ Bukan `LIKE '%%'`. Keduanya bermakna sama di SQL, tetapi yang
	// pertama membiarkan Oracle memakai indeks pada kolom yang satunya.
	q, arg := sqlCariPenyakit(tabelUjiPenyakit,
		models.KriteriaPenyakit{Nama: "DIABETES"}, 50)
	if len(arg) != 1 {
		t.Fatalf("argumen = %d, mau 1", len(arg))
	}
	if strings.Contains(q, kolomICDPenyakit+") LIKE") {
		t.Errorf("kriteria kode yang kosong tetap menghasilkan klausa:\n%s", q)
	}
	kosong, argKosong := sqlCariPenyakit(tabelUjiPenyakit, models.KriteriaPenyakit{}, 50)
	if len(argKosong) != 0 {
		t.Errorf("kriteria kosong menghasilkan %d argumen", len(argKosong))
	}
	if strings.Contains(kosong, "WHERE") {
		t.Errorf("kriteria kosong tetap menghasilkan WHERE:\n%s", kosong)
	}
}

func TestNormalkanKriteriaPenyakit(t *testing.T) {
	// ⛔ CARI1 -> ICD_Code, CARI2 -> Disease (b1657-1658). Dibaca dari
	// pemetaannya, bukan ditebak dari namanya - "CARI1/CARI2" tidak
	// menyebutkan apa pun, dan menukarnya membuat setiap pencarian gagal
	// dengan cara yang terlihat seperti "datanya tidak ada".
	k := models.NormalkanKriteriaPenyakit(" a00 ", " diabetes ")
	if k.KodeICD != "A00" {
		t.Errorf("KodeICD = %q, mau %q", k.KodeICD, "A00")
	}
	if k.Nama != "DIABETES" {
		t.Errorf("Nama = %q, mau %q", k.Nama, "DIABETES")
	}
	if !models.NormalkanKriteriaPenyakit("", "  ").Kosong() {
		t.Error("kriteria spasi belaka tidak terbaca kosong")
	}
}

// TestKolomPenyakitBelumDipastikan menagih OQ-K.
//
// ⛔ Ekspor yang kami terima TIDAK memuat pemetaan kelas-ke-tabel untuk
// `Int-DISEASE_LIFE`; yang ada hanya nama properti Pega. Ketiga nama kolom
// di `penyakit.go` karena itu belum dipastikan, dan `NUMBER_` paling lemah -
// `NUMBER` kata cadangan Oracle, jadi kolomnya pasti bernama lain.
//
// ⚠️ Uji ini tidak memeriksa kebenaran nama - ia tidak bisa. Yang ia jaga
// adalah supaya pertanyaannya TIDAK HILANG: ketiganya harus tetap terkumpul
// di satu tempat, sehingga koreksi DBA adalah satu suntingan dan bukan
// perburuan. Bila kelak jawabannya datang, ubah konstanta-konstanta itu dan
// hapus uji ini beserta OQ-K.
func TestKolomPenyakitBelumDipastikan(t *testing.T) {
	q, _ := sqlCariPenyakit(tabelUjiPenyakit, models.KriteriaPenyakit{Nama: "X"}, 50)
	for _, kolom := range []string{kolomNomorPenyakit, kolomNamaPenyakit, kolomICDPenyakit} {
		if kolom == "" {
			t.Fatal("nama kolom kosong")
		}
		if !strings.Contains(q, kolom) {
			t.Errorf("kolom %q tidak lagi dipakai query; bila pemetaannya sudah "+
				"dipastikan DBA, ubah konstantanya dan tutup OQ-K", kolom)
		}
	}
	// Nama kolom tidak pernah tertempel dari luar.
	if strings.Contains(q, ";") || strings.Contains(q, "--") {
		t.Errorf("query memuat pemisah pernyataan atau komentar:\n%s", q)
	}
}
