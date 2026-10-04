package repository

// Nol konstanta ambang tutup buku di kode produksi - AC tiket 02.
//
// ⛔ SEBAB PENJAGA INI ADA. Pega punya DUA aturan hidup berdampingan:
// `PROC_GENERATE_SEQUENCE_NUMBER` menggulir periode nomor PL dari tabel
// `TANGGAL_CLOSING`, sedangkan `InsertJsonPolisLife_Act` (20260728) masih
// menanam `25` di precondition langkah 4 b1170. `[keputusan work owner]`
// ikuti yang dari DB. (Ralat sensus 28-09-2026: yang dulu disebut pembaca
// tabel, `SubmitPremiumList_Act`, hanya mengalirkan nilainya ke langkah 15
// yang ter-remark - OQ-PL-13.)
//
// Artinya menanam angkanya kembali adalah kekeliruan yang sudah pernah
// terjadi di sistem yang ditiru, dan yang paling mungkin terulang: ia
// terlihat seperti "bawaan yang aman", dan ia membukukan transaksi ke
// periode yang salah tanpa satu pun galat.

import (
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
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
		// Refactor bentuk B (30-09-2026): lapisan yang sama di SETIAP modul,
		// ditambah `inti/` - periode produksi kini tinggal di
		// `inti/backend/penomor/periode.go`. Dulu hanya `/internal/...`: saringan itu
		// diam-diam menyempit (cacah log dasar 158; tanpa perbaikan ini, 107).
		if !strings.Contains(nama, "/repository/") &&
			!strings.Contains(nama, "/services/") &&
			!strings.Contains(nama, "/models/") &&
			!strings.Contains(nama, "/inti/") &&
			// Paket 5: skema uji pindah dari internal/repository/skemauji.
			!strings.Contains(nama, "/uji/skemauji/") {
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
	if err := db.PeriksaSQL(q); err != nil {
		t.Errorf("%v\n%s", err, q)
	}
}
