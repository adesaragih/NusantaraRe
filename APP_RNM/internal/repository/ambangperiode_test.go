package repository

// Nol konstanta ambang tutup buku di kode produksi - AC tiket 02.
//
// ⛔ SEBAB PENJAGA INI ADA. Pega punya DUA versi hidup berdampingan:
// `SubmitPremiumList_Act` (20260122) membaca ambang dari tabel, sedangkan
// `InsertJsonPolisLife_Act` (20260728 - LEBIH BARU) masih menanam `25` di
// preconditionnya b1170. `[keputusan work owner]` ikuti yang dari DB.
//
// Artinya menanam angkanya kembali adalah kekeliruan yang sudah pernah
// terjadi di sistem yang ditiru, dan yang paling mungkin terulang: ia
// terlihat seperti "bawaan yang aman", dan ia membukukan transaksi ke
// periode yang salah tanpa satu pun galat.

import (
	"regexp"
	"strings"
	"testing"
)

// polaAmbangTertanam mencocokkan HARI INI yang dibandingkan dengan angka.
//
// ⛔ DIPERSEMPIT sesudah ronde pertamanya menuduh hal yang BENAR.
// Ronde itu mencocokkan `(tanggal|tutupBuku|...) [<>]= N`, dan ia menyalakan
// `tglTutupBuku <= 0` - pemeriksaan KEKOSONGAN yang justru menjadi inti
// tiket 02. Penjaga yang menuduh hal yang benar akan dilonggarkan orang,
// bukan dipatuhi.
//
// Yang membedakan cacat dari pemeriksaan sah adalah SISI KIRINYA:
//
//	cacat  : HARI INI  dibandingkan dengan angka   `lokal.Day() > 25`
//	sah    : AMBANGNYA dibandingkan dengan angka   `tglTutupBuku <= 0`
//
// Yang pertama membukukan transaksi ke periode yang salah; yang kedua
// menolak tabel yang kosong.
var polaAmbangTertanam = regexp.MustCompile(
	`(?i)(\.Day\(\)|currentdate|hariIni|tanggalKini)\s*>=?\s*[0-9]{1,2}\b`)

// polaKomentarAmbang membuang komentar sebelum pencocokan - prosa yang
// MENERANGKAN ambangnya tidak boleh dituduh sebagai ambangnya.
var polaKomentarAmbang = regexp.MustCompile(`(?m)^\s*//.*$`)

func TestNolAmbangTutupBukuTertanam(t *testing.T) {
	diperiksa := 0
	for nama, isi := range berkasGoSelainTest(t) {
		// ⛔ KETIGA lapisan, dan  yang PALING penting: aturan
		// periodenya tinggal di sana. Ronde pertama penjaga ini hanya
		// memindai repository dan services - lalu ambang tertanam yang
		// SENGAJA ditanam di  untuk membuktikannya
		// TIDAK menyalakannya. Penjaga yang tidak memandang tempat aturannya
		// tinggal adalah penjaga yang menjaga tempat yang salah.
		if !strings.Contains(nama, "/internal/repository/") &&
			!strings.Contains(nama, "/internal/services/") &&
			!strings.Contains(nama, "/internal/models/") {
			continue
		}
		diperiksa++
		bersih := polaKomentarAmbang.ReplaceAllString(isi, "")
		for _, m := range polaAmbangTertanam.FindAllString(bersih, -1) {
			t.Errorf("%s: ambang tutup buku tertanam %q.\n"+
				"AC tiket 02: tanggalnya dibaca dari POOLDATA.TANGGAL_CLOSING "+
				"setiap kali dibutuhkan - bukan dari konstanta. Pega sendiri "+
				"sudah terpeleset begitu di InsertJsonPolisLife_Act b1170.",
				nama, strings.TrimSpace(m))
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol berkas terbaca; pembacanya yang rusak, bukan kodenya")
	}
	t.Logf("%d berkas diperiksa", diperiksa)
}

func TestPenjagaAmbangMasihMenggigit(t *testing.T) {
	// ⛔ Penjaga yang polanya tidak pernah cocok adalah penjaga yang selalu
	// hijau. Diberi kode buatan yang jelas salah, tanpa menyentuh berkas
	// produksi mana pun.
	buruk := "if lokal.Day() > 25 {\n\tbulan++\n}"
	if !polaAmbangTertanam.MatchString(buruk) {
		t.Error("pola tidak menemukan ambang tertanam yang jelas")
	}
	// Dan yang AMAN tidak tertuduh - termasuk pemeriksaan kekosongan, yang
	// ronde pertama penjaga ini tuduh dengan keliru.
	for _, aman := range []string{
		"NUMBER(25)", "VARCHAR2(255)", "baris[25]",
		"if lokal.Day() > tglTutupBuku {",
		"if tglTutupBuku <= 0 {",
		"if tglTutupBuku > 31 {",
	} {
		if polaAmbangTertanam.MatchString(aman) {
			t.Errorf("kode aman tertuduh: %q", aman)
		}
	}
}

func TestQueryTanggalTutupBukuBerbatasSatuBaris(t *testing.T) {
	q := sqlTanggalTutupBuku("SKEMAUJI.TANGGAL_CLOSING")
	// ⚠️ Rule aslinya `SELECT *`; yang dipakai hanya kolom `TANGGAL`
	// (`SubmitPremiumList_Act` b918). Membawa seluruh kolom berarti
	// perubahan tabel di hulu mengubah bentuk baris kami tanpa diminta.
	if !strings.Contains(q, "SELECT TANGGAL FROM") {
		t.Errorf("query bukan SELECT TANGGAL:\n%s", q)
	}
	if strings.Contains(q, "SELECT *") {
		t.Errorf("query membawa seluruh kolom:\n%s", q)
	}
	if !strings.Contains(q, "FETCH FIRST 1 ROWS ONLY") {
		t.Errorf("query tidak berbatas satu baris:\n%s", q)
	}
	if err := PeriksaSQL(q); err != nil {
		t.Errorf("%v\n%s", err, q)
	}
}
