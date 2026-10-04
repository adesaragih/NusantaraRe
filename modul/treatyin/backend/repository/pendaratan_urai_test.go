package repository_test

// Uji pengurai pendaratan - NOL koneksi Oracle.
//
// Bentuk data ujinya dari sapuan 3 Oktober 2026 atas SELURUH 1.854 dokumen
// `POOLDATA.M_TREATY_IN.JSONDATA`, bukan dikarang.

import (
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/repository"
)

// ⛔ UJI YANG PALING PENTING DI BERKAS INI.
//
// `AmountTotal` terpanjang di korpus 61 aksara. `encoding/json` tanpa
// `UseNumber()` mengubahnya menjadi `float64`, yang membawa 15-17 digit
// berarti - sisanya hilang, DIAM-DIAM, dan hasilnya tetap terlihat seperti
// angka yang masuk akal. Uji ini membandingkan literalnya aksara demi
// aksara; tidak ada cara lulus kecuali nilainya tidak pernah disentuh.
func TestAngkaPanjangTidakDibulatkan(t *testing.T) {
	const panjang = "1323411750.0000008394305684816601000000000000000000000000000"
	doc, err := repository.UraiDokumen(`{"Installment":[{"AmountTotal":` + panjang + `}]}`)
	if err != nil {
		t.Fatal(err)
	}
	baris := repository.BarisLarik(doc, "Installment")
	if len(baris) != 1 {
		t.Fatalf("%d baris, mau 1", len(baris))
	}
	got, ada := repository.NilaiTeks(baris[0]["AmountTotal"])
	if !ada {
		t.Fatal("nilainya hilang")
	}
	if got != panjang {
		t.Errorf("angka berubah saat diurai:\n  dapat %s\n  mau   %s", got, panjang)
	}
}

// Kunci yang TIDAK ADA dan kunci yang ADA tetapi kosong harus berbeda:
// yang pertama mendarat NULL, yang kedua mendarat string kosong. Di tabel
// pendaratan bedanya satu-satunya bukti apakah Pega pernah menulis medan
// itu.
func TestKunciHilangBedaDariNilaiKosong(t *testing.T) {
	doc, err := repository.UraiDokumen(`{"Retention":[{"Currency":"","Note":null}]}`)
	if err != nil {
		t.Fatal(err)
	}
	el := repository.BarisLarik(doc, "Retention")[0]

	if teks, ada := repository.NilaiTeks(el["Currency"]); !ada || teks != "" {
		t.Errorf("kunci ada bernilai kosong: teks=%q ada=%v, mau \"\" dan true", teks, ada)
	}
	if _, ada := repository.NilaiTeks(el["Note"]); ada {
		t.Error("null terbaca sebagai nilai; ia harus mendarat NULL")
	}
	if _, ada := repository.NilaiTeks(el["TreatyGroup"]); ada {
		t.Error("kunci yang tidak ada terbaca sebagai nilai")
	}
}

// Larik KOSONG di korpus ditulis `[ ]`, dengan spasi - dan itu pernah
// membuat sapuan berbasis teks mengembalikan nol di mana-mana. Pengurai
// JSON tidak peduli, dan uji ini mengunci bahwa ia memang tidak peduli.
func TestLarikKosongDanKunciHilangSamaSamaNolBaris(t *testing.T) {
	for _, dok := range []string{
		`{"Portfolio":[ ]}`,
		`{"Portfolio":[]}`,
		`{}`,
		`{"Portfolio":null}`,
		`{"Portfolio":"bukan larik"}`,
	} {
		doc, err := repository.UraiDokumen(dok)
		if err != nil {
			t.Fatalf("%s: %v", dok, err)
		}
		if n := len(repository.BarisLarik(doc, "Portfolio")); n != 0 {
			t.Errorf("%s menghasilkan %d baris, mau 0", dok, n)
		}
	}
}

// Dokumen rusak -> galat yang dapat dikenali, bukan panik.
func TestDokumenRusakMenghasilkanGalatYangDikenali(t *testing.T) {
	_, err := repository.UraiDokumen(`{"Portfolio":[`)
	if err == nil {
		t.Fatal("dokumen rusak diterima")
	}
	if !errors.Is(err, repository.ErrJSONWarisanRusak) {
		t.Errorf("galatnya %v; mau yang dikenali sebagai JSON rusak", err)
	}
}

const dokAngsuran = `{"Installment":[
  {"AmountTotal":"100","Currency":"IDR","InstallmentList":[
     {"Installment":"1","Amount":"25"},
     {"Installment":"2","Amount":"25"}]},
  {"AmountTotal":"50","Currency":"USD","InstallmentList":[
     {"Installment":"1","Amount":"50"}]}]}`

// Butir angsuran diratakan DALAM URUTAN BACA, dan cacahnya menyeberangi
// kedua induk. Larik Pega berurut; urutan yang hilang tidak terlihat sampai
// seseorang membandingkan layar dengan sistem lama.
func TestButirAngsuranDiratakanBerurut(t *testing.T) {
	doc, err := repository.UraiDokumen(dokAngsuran)
	if err != nil {
		t.Fatal(err)
	}
	cacah := repository.CacahLarik(doc)
	if cacah["M_TREATYIN_INSTALLMENT"] != 2 {
		t.Errorf("induk angsuran %d, mau 2", cacah["M_TREATYIN_INSTALLMENT"])
	}
	if cacah["M_TREATYIN_INSTALLMENTITEM"] != 3 {
		t.Errorf("butir angsuran %d, mau 3", cacah["M_TREATYIN_INSTALLMENTITEM"])
	}
}

// ⛔ Kunci BARU dari Pega tidak boleh hilang tanpa suara. Pega menambah
// properti tanpa memberi tahu siapa pun; kunci tanpa kolom mendarat di
// mana-mana sebagai ketiadaan, dan cacah barisnya tetap cocok.
func TestKunciBaruDilaporkanBukanDitelan(t *testing.T) {
	doc, err := repository.UraiDokumen(
		`{"Retention":[{"Amount":"1","KunciYangBelumPernahAda":"x"}]}`)
	if err != nil {
		t.Fatal(err)
	}
	asing := repository.KunciTakTerpetakan(doc)
	daftar, ada := asing["M_TREATYIN_RETENTION"]
	if !ada {
		t.Fatalf("kunci baru tidak dilaporkan; yang dilaporkan: %v", asing)
	}
	if len(daftar) != 1 || daftar[0] != "KunciYangBelumPernahAda" {
		t.Errorf("dilaporkan %v, mau [KunciYangBelumPernahAda]", daftar)
	}
}

// Pasangannya: dokumen yang SELURUH kuncinya terpetakan melaporkan nol.
// Tanpa uji ini, pelapor yang selalu mengeluh lulus uji di atas.
func TestDokumenTerpetakanPenuhMelaporkanNol(t *testing.T) {
	doc, err := repository.UraiDokumen(dokAngsuran)
	if err != nil {
		t.Fatal(err)
	}
	if asing := repository.KunciTakTerpetakan(doc); len(asing) != 0 {
		t.Errorf("kunci asing dilaporkan padahal seluruhnya terpetakan: %v", asing)
	}
}

// `InstallmentList` adalah larik ANAK, bukan kolom - ia tidak boleh
// dilaporkan sebagai kunci yang belum punya rumah.
func TestLarikAnakBukanKunciAsing(t *testing.T) {
	doc, err := repository.UraiDokumen(dokAngsuran)
	if err != nil {
		t.Fatal(err)
	}
	for tabel, daftar := range repository.KunciTakTerpetakan(doc) {
		for _, k := range daftar {
			if k == repository.LarikAnakAngsuran {
				t.Errorf("%s melaporkan %s sebagai kunci tanpa kolom; ia larik anak", tabel, k)
			}
		}
	}
}

// Peta kunci dan peta kolom WAJIB sejajar, dan nama tabelnya unik.
// Keduanya diperiksa di sini sebab seluruh SQL dibangkitkan darinya: satu
// kunci yang tergeser memindahkan nilai ke kolom tetangganya, dan hasilnya
// tetap berupa tabel yang terisi rapi.
func TestPetaPendaratanSejajarDanUnik(t *testing.T) {
	lihat := map[string]bool{}
	for _, p := range repository.PetaPendaratan {
		if len(p.Kunci) != len(p.Kolom) {
			t.Errorf("%s: %d kunci tetapi %d kolom", p.Tabel, len(p.Kunci), len(p.Kolom))
		}
		if lihat[p.Tabel] {
			t.Errorf("%s terdaftar dua kali", p.Tabel)
		}
		lihat[p.Tabel] = true
		if p.Seq == "" || p.Tabel == "" {
			t.Errorf("%+v: tabel atau sequence kosong", p)
		}
		for _, k := range p.Kolom {
			if k != strings.ToUpper(k) {
				t.Errorf("%s: kolom %q tidak huruf besar; DDL menulisnya huruf besar", p.Tabel, k)
			}
		}
	}
	if len(repository.PetaPendaratan) != 9 {
		t.Errorf("%d tabel pendaratan, mau 9", len(repository.PetaPendaratan))
	}
}

// Tetapan indeks induk-anak menunjuk tabel yang benar. Keduanya tetapan
// angka, dan angka yang menunjuk tempat yang salah melewati SELURUH butir
// angsuran tanpa satu pun galat.
func TestIndeksAngsuranMenunjukTabelYangBenar(t *testing.T) {
	p := repository.PetaPendaratan
	if p[repository.IndeksAngsuran].Tabel != "M_TREATYIN_INSTALLMENT" {
		t.Errorf("IndeksAngsuran menunjuk %s", p[repository.IndeksAngsuran].Tabel)
	}
	if p[repository.IndeksButirAngsuran].Tabel != "M_TREATYIN_INSTALLMENTITEM" {
		t.Errorf("IndeksButirAngsuran menunjuk %s", p[repository.IndeksButirAngsuran].Tabel)
	}
	// Anaknya SESUDAH induknya; pemuat bersandar pada urutan itu.
	if repository.IndeksButirAngsuran <= repository.IndeksAngsuran {
		t.Error("butir angsuran dimuat sebelum induknya; IDINDUK-nya belum lahir")
	}
}

// `Date` -> `TANGGAL`, dan tidak ada kolom bernama `DATE` di peta mana pun.
func TestKunciDateDipetakanKeTanggal(t *testing.T) {
	for _, p := range repository.PetaPendaratan {
		for j, k := range p.Kunci {
			if k == "Date" && p.Kolom[j] != "TANGGAL" {
				t.Errorf("%s: kunci Date dipetakan ke %q, mau TANGGAL", p.Tabel, p.Kolom[j])
			}
			if p.Kolom[j] == "DATE" {
				t.Errorf("%s: kolom DATE kata cadangan Oracle (ORA-00923)", p.Tabel)
			}
		}
	}
}
