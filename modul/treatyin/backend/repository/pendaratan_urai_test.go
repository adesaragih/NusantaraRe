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
	if cacah["T_TREATY_INSTALLMENT"] != 2 {
		t.Errorf("induk angsuran %d, mau 2", cacah["T_TREATY_INSTALLMENT"])
	}
	if cacah["T_TREATY_INSTALLMENT_ITEM"] != 3 {
		t.Errorf("butir angsuran %d, mau 3", cacah["T_TREATY_INSTALLMENT_ITEM"])
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
	daftar, ada := asing["T_TREATY_RETENTION"]
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
	// 22 + 3 (migrasi 438) + 4 (migrasi 439) = 29.
	if len(repository.PetaPendaratan) != 29 {
		t.Errorf("%d tabel pendaratan, mau 29", len(repository.PetaPendaratan))
	}
}

// Tetapan indeks induk-anak menunjuk tabel yang benar. Keduanya tetapan
// angka, dan angka yang menunjuk tempat yang salah melewati SELURUH butir
// angsuran tanpa satu pun galat.
func TestIndeksAngsuranMenunjukTabelYangBenar(t *testing.T) {
	p := repository.PetaPendaratan
	if p[repository.IndeksAngsuran].Tabel != "T_TREATY_INSTALLMENT" {
		t.Errorf("IndeksAngsuran menunjuk %s", p[repository.IndeksAngsuran].Tabel)
	}
	if p[repository.IndeksButirAngsuran].Tabel != "T_TREATY_INSTALLMENT_ITEM" {
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

// ⛔ `CacahLarik` HARUS mengenal KEEMPAT bentuk yang `MuatKontrak` tangani.
//
// Ronde 6 Oktober 2026 menemukannya buta terhadap dua di antaranya —
// `LarikGabung` dan `Akar`, yaitu sembilan tabel migrasi 438/439. Kebutaan
// itu TIDAK memunculkan galat: rekonsiliasi hanya melaporkan "dokumen
// memberi 0", dan selisihnya terbaca sebagai pemuat yang rusak.
//
// ⚠️ Dokumen ujinya DIBANGKITKAN DARI PETA, bukan ditulis tangan. Dokumen
// tulisan tangan berhenti menyebut tabel ke-31 pada hari tabel itu lahir,
// dan berhentinya tidak terlihat.
func TestCacahLarikMengenalKeempatBentuk(t *testing.T) {
	// Satu elemen berisi apa pun — isinya tidak dibaca, hanya dicacah.
	el := func() map[string]any { return map[string]any{"Currency": "IDR"} }

	// Dibangun dua lintasan: induk lebih dulu, supaya elemen anak dapat
	// disisipkan ke dalam elemen induk yang sudah ada.
	doc := map[string]any{}
	indukEl := map[string]map[string]any{}
	for _, p := range repository.PetaPendaratan {
		if p.Induk != "" || p.Akar {
			continue
		}
		e := el()
		indukEl[p.Tabel] = e
		switch {
		case len(p.LarikGabung) > 0:
			for _, nama := range p.LarikGabung {
				doc[nama] = []any{el()}
			}
			// Elemen pertama larik gabung pertama menjadi wakil induknya.
			indukEl[p.Tabel] = doc[p.LarikGabung[0]].([]any)[0].(map[string]any)
		default:
			doc[p.Larik] = []any{e}
		}
	}
	for _, p := range repository.PetaPendaratan {
		if p.Induk == "" {
			continue
		}
		ind, ada := indukEl[p.Induk]
		if !ada {
			t.Fatalf("%s berinduk %s yang tidak ada di peta", p.Tabel, p.Induk)
		}
		e := el()
		if len(p.LarikGabung) > 0 {
			for _, nama := range p.LarikGabung {
				ind[nama] = []any{el()}
			}
			e = ind[p.LarikGabung[0]].([]any)[0].(map[string]any)
		} else {
			ind[p.KunciAnak] = []any{e}
		}
		indukEl[p.Tabel] = e
	}

	cacah := repository.CacahLarik(doc)
	for _, p := range repository.PetaPendaratan {
		mau := 1
		if len(p.LarikGabung) > 0 {
			mau = len(p.LarikGabung)
		}
		if cacah[p.Tabel] != mau {
			t.Errorf("%s: CacahLarik %d, dokumen uji memberi %d — "+
				"bentuk %s tidak dikenali", p.Tabel, cacah[p.Tabel], mau, bentuk(p))
		}
	}
}

// bentuk menamai salah satu dari keempat bentuk pendaratan, untuk pesan galat.
func bentuk(p repository.Pendaratan) string {
	switch {
	case p.Akar:
		return "Akar"
	case p.Induk == "" && len(p.LarikGabung) > 0:
		return "LarikGabung akar"
	case p.Induk == "":
		return "Larik akar"
	case len(p.LarikGabung) > 0:
		return "LarikGabung anak"
	default:
		return "Larik anak"
	}
}

// ⛔ `KunciTakTerpetakan` TIDAK boleh menyebut kunci yang menampung tabel
// ANAK — itu larik, bukan kolom yang hilang.
//
// Sampai 6 Oktober 2026 ia mengecualikan SATU kunci anak
// (`InstallmentList`) lewat tetapan tulisan tangan, sehingga ketiga belas
// kunci anak lain dilaporkan asing pada hampir setiap dokumen. Daftar yang
// selalu memuat belasan nama palsu membuat nama sungguhan — properti baru
// dari Pega, yang justru menjadi alasan daftar ini ada — tidak terlihat.
func TestKunciAnakBukanKunciAsing(t *testing.T) {
	// Dokumen yang tiap tabel induknya memuat larik anaknya, dibangkitkan
	// dari peta supaya tabel ke-30 tidak terlewat pada hari ia lahir.
	doc := map[string]any{}
	el := map[string]map[string]any{}
	for _, p := range repository.PetaPendaratan {
		if p.Induk != "" || p.Akar {
			continue
		}
		e := map[string]any{}
		if len(p.LarikGabung) > 0 {
			for _, n := range p.LarikGabung {
				doc[n] = []any{map[string]any{}}
			}
			e = doc[p.LarikGabung[0]].([]any)[0].(map[string]any)
		} else {
			doc[p.Larik] = []any{e}
		}
		el[p.Tabel] = e
	}
	for _, p := range repository.PetaPendaratan {
		if p.Induk == "" {
			continue
		}
		ind := el[p.Induk]
		if ind == nil {
			t.Fatalf("%s berinduk %s yang tidak ada", p.Tabel, p.Induk)
		}
		e := map[string]any{}
		if len(p.LarikGabung) > 0 {
			for _, n := range p.LarikGabung {
				ind[n] = []any{map[string]any{}}
			}
			e = ind[p.LarikGabung[0]].([]any)[0].(map[string]any)
		} else {
			ind[p.KunciAnak] = []any{e}
		}
		el[p.Tabel] = e
	}

	asing := repository.KunciTakTerpetakan(doc)
	for tabel, k := range asing {
		t.Errorf("%s: %v dilaporkan asing, padahal dokumen uji HANYA memuat "+
			"kunci anak dan larik tingkat pertama", tabel, k)
	}
}

// ⭐ Dan penjaga yang SEBENARNYA tetap bekerja: kunci yang betul-betul baru
// tetap dilaporkan, termasuk di dalam tabel ANAK — yang dulu tidak pernah
// diperiksa sama sekali.
func TestKunciBaruDiTabelAnakTetapDilaporkan(t *testing.T) {
	doc, err := repository.UraiDokumen(
		`{"Installment":[{"AmountTotal":"1","InstallmentList":[{"Installment":"1","KunciAnakYangBelumPernahAda":"x"}]}]}`)
	if err != nil {
		t.Fatal(err)
	}
	asing := repository.KunciTakTerpetakan(doc)
	k := asing["T_TREATY_INSTALLMENT_ITEM"]
	if len(k) != 1 || k[0] != "KunciAnakYangBelumPernahAda" {
		t.Errorf("T_TREATY_INSTALLMENT_ITEM asing = %v, mau [KunciAnakYangBelumPernahAda]", k)
	}
	if len(asing["T_TREATY_INSTALLMENT"]) != 0 {
		t.Errorf("induk ikut terlapor asing: %v", asing["T_TREATY_INSTALLMENT"])
	}
}

// ⛔ Daftar kolom PEMBACA tidak boleh menyebut kolom yang peta tidak isi.
//
// `pendaratan_akar.go` memegang daftarnya SENDIRI (`kolomRevisi`,
// `kolomTeksRevisi`) karena ia membedakan `VARCHAR2` dari `CLOB` — pembedaan
// yang peta tidak punya. Dua daftar untuk satu tabel akan berselisih, dan
// selisihnya DIAM dalam dua arah:
//
//   - kolom di pembaca yang TIDAK ada di peta -> `NULL` selamanya, sebab
//     nol yang mengisinya; atau ORA-00904 bila kolomnya pun tak ada.
//   - kolom di peta yang tidak ada di pembaca -> terisi di tabel tetapi
//     tidak pernah sampai ke layar.
//
// ⚠️ Yang diuji hanya arah PERTAMA. Arah kedua SAH: peta boleh mendaratkan
// kolom yang layar belum perlukan, dan menuntut keduanya sama akan memaksa
// layar membaca kolom yang tidak ia pakai.
func TestKolomRevisiAdaDiPeta(t *testing.T) {
	var peta *repository.Pendaratan
	for i := range repository.PetaPendaratan {
		if repository.PetaPendaratan[i].Tabel == "T_TREATY_REVISION" {
			peta = &repository.PetaPendaratan[i]
			break
		}
	}
	if peta == nil {
		t.Fatal("T_TREATY_REVISION tidak ada di peta")
	}
	diPeta := map[string]string{} // KOLOM -> kunci dokumen
	for i, k := range peta.Kolom {
		diPeta[k] = peta.Kunci[i]
	}
	for _, p := range repository.KolomRevisiUntukUji() {
		kunci, ada := diPeta[p[0]]
		if !ada {
			t.Errorf("pembaca menyebut kolom %s yang peta TIDAK isi — "+
				"ia akan NULL selamanya", p[0])
			continue
		}
		// ⭐ Ejaan dokumennya juga diadu: pembaca memakai ejaan itu sebagai
		// KUNCI petanya, dan kunci yang meleset membuat medan hilang dari
		// layar tanpa satu pun galat.
		if kunci != p[1] {
			t.Errorf("kolom %s: pembaca memakai ejaan %q, peta mengisinya dari %q",
				p[0], p[1], kunci)
		}
	}
}
