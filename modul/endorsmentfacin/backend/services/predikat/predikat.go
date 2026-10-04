// Package predikat adalah registry rule `When` siklus Endorsement Fac In: SATU
// registry untuk seluruh `Endorsment Fac In\When`, ditanyai dengan nama lewat
// `Eval` (tiket E01). Sumber data adalah atribut PER-RULE: jalur properti di
// tiap baris kondisi (`OutData.pxResults(1).CARI2` untuk enam predikat lini,
// properti kasus untuk sisanya) diselesaikan `Kasus` pemanggil.
//
// ⚠️ ASAL SALINAN - mesin evaluator di berkas ini dan `logika/` DISALIN dari
// `modul/nbfacin/backend/services/rules/rules.go` + `rules/logika/` pada
// 01-10-2026 (sha256 rules.go 9f881801…, logika.go 6a9da6b4…, logika_test.go
// 7eabf9e7…; berkas sumber saat itu belum ter-commit di nbfacin). Dasarnya:
// jawaban work owner "EDM DAN NB MENU YANG TERPISAH" (register nbfacin butir
// 57), ditafsir sesi nusantarare-c3 dan nusantarare-0f sebagai: NB dan EDM
// modul terpisah penuh, mesin TIDAK dipindah ke inti, endorsmentfacin memegang
// mesinnya sendiri. ⚠️ Itu TAFSIRAN, bukan kutipan. Paket ini TIDAK mengimpor
// modul/nbfacin (penjaga impor lintas modul). Perubahan isi terhadap sumber:
// nama paket, jalur impor `logika`, awalan pesan galat/panic ("predikat:"),
// komentar, dan SATU penyesuaian EDM - `sikapKhusus` diperiksa sebelum registry
// (penilai.eval, kasus.go). Logika evaluasi lainnya sama.
//
// Isi registry (`registry_gen.go`) DIBANGKITKAN generator milik nbfacin
// (`rules/bangkit`, opsi (b) - keputusan work owner, butir 56) lalu disaring
// `docs/alat/saring_registry.py`; jangan disunting tangan. Sikap yang mengikat
// mesin: keputusan work owner nbfacin butir 19-25 dan 28 (dibawa bersama
// salinan; berlaku sampai work owner memutus lain untuk EDM).
package predikat

import (
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/endorsmentfacin/backend/services/predikat/logika"
)

// ErrTafsirBerbeda - perbandingan yang hasilnya bergantung pada tipe properti,
// padahal tipe itu tidak ada di korpus (nol rule Property): tafsir teks dan
// tafsir angka berbeda, atau nilai kosong dibandingkan dengan angka. Ditolak,
// bukan ditebak (keputusan work owner 01-10-2026, butir 20).
var ErrTafsirBerbeda = errors.New("predikat: tafsir teks dan angka berbeda, tipe properti belum terverifikasi")

// Kasus memberi nilai properti clipboard menurut JALUR PERSIS seperti tertulis
// di rule (`pyWorkPage.Quotation.BusinessCode`, `.Quotation.BusinessType`,
// `OperatorID.pyWorkBasketList(1).pyWorkBasketName`). Jalur relatif (berawal
// titik) diselesaikan Kasus terhadap halaman tempat predikat dievaluasi.
//
// [dugaan] Properti yang tidak ada (`ada == false`) diperlakukan sebagai teks
// kosong, sama dengan properti kosong.
type Kasus interface {
	Nilai(jalur string) (nilai string, ada bool)
}

// jenisKondisi - dua bentuk baris kondisi yang diport.
type jenisKondisi uint8

const (
	// banding - `@(Pega-RULES:ExpressionEvaluators).compareTwoValues(kiri, op, kanan)`.
	banding jenisKondisi = iota + 1
	// rujukWhen - `@(Pega-RULES:ExpressionEvaluators).evaluateWhen("Nama")`.
	rujukWhen
)

// jenisOperand - bentuk sisi kanan `compareTwoValues` yang diport. Kata tanpa
// kutip (`= ReasFacInDirector`, `= A1`) dibaca sebagai teks (keputusan work owner
// 01-10-2026, butir 28, mengganti butir 24). `= True` berhuruf besar dan login
// operator literal TIDAK diport - lihat pembangkit.
type jenisOperand uint8

const (
	// teksLiteral - `"10053"`, juga kata tanpa kutip `ReasFacInDirector`.
	teksLiteral jenisOperand = iota + 1
	// angkaLiteral - `1`.
	angkaLiteral
	// booleanLiteral - `true` / `false` huruf kecil, dibandingkan sebagai teksnya.
	booleanLiteral
)

type operand struct {
	jenis jenisOperand
	teks  string
}

type kondisi struct {
	jenis jenisKondisi
	// banding
	kiri      string // jalur properti
	tidakSama bool   // operator "!=" (selain itu "=")
	kanan     operand
	// rujukWhen - nama predikat, huruf besar
	rujukan string
}

type predikat struct {
	// asal - berkas korpus dan identitas `<pxInsName>`; lebih dari satu bila
	// identitas yang sama tersimpan di beberapa berkas berisi identik.
	asal []string
	// logika - `<pyLogic>`: label kondisi dirangkai AND/OR/kurung/negasi.
	logika  string
	kondisi map[string]kondisi
	// panik - bila terisi, Eval selalu panic dengan alasan ini.
	panik string
	// catatan - sikap khusus yang diputuskan untuk predikat ini (tiket 09),
	// mis. teks tampilan yang tidak dieksekusi atau arti nama yang menyesatkan.
	catatan string
}

// pohonLogika - `<pyLogic>` setiap predikat yang tidak panic, diurai SEKALI saat
// paket dimuat. Logika yang tidak terurai berarti registry rusak: panic saat
// dimuat, bukan saat dievaluasi.
var pohonLogika = uraiSemuaLogika()

func uraiSemuaLogika() map[string]logika.Simpul {
	hasil := map[string]logika.Simpul{}
	for nama, pr := range registry {
		if pr.panik != "" {
			continue
		}
		s, err := logika.Urai(pr.logika)
		if err != nil {
			panic(fmt.Sprintf("predikat: registry rusak, logika %s %q: %v", nama, pr.logika, err))
		}
		hasil[nama] = s
	}
	return hasil
}

// Eval mengevaluasi predikat bernama `nama` atas `k`. Nama tidak peka huruf
// besar-kecil: [terverifikasi] korpus sendiri merujuk dengan ejaan berbeda dari
// identitasnya (mis. `evaluateWhen("isGrowingTrees")`), dan identitas
// `<pxInsName>` tersimpan berhuruf besar.
//
// Panic - kesalahan program: nama tidak dikenal, predikat yang sikapnya panic
// (IsPKSASM), bentuk kondisi yang belum diport, rujukan melingkar.
// Galat - keraguan yang bergantung data: ErrTafsirBerbeda.
//
// SEMUA kondisi yang dirujuk logika dinilai, urut label (A, B, …, AA), tanpa
// hubung-singkat: satu perbandingan yang ditolak membuat seluruh predikat
// ditolak, sehingga hasil tidak pernah bergantung pada urutan evaluasi Pega
// yang belum terverifikasi, dan galat yang dilaporkan selalu sama (butir 25).
func Eval(nama string, k Kasus) (bool, error) {
	return (&penilai{kasus: k, sedang: map[string]bool{}}).eval(strings.ToUpper(nama))
}

type penilai struct {
	kasus Kasus
	// sedang - predikat di rantai evaluasi saat ini, untuk mengenali rujukan
	// melingkar.
	sedang map[string]bool
}

func (p *penilai) eval(nama string) (bool, error) {
	// Penyesuaian EDM (bukan salinan nbfacin): sikap khusus menggantikan entri
	// registry yang bentuk ekspresinya tidak diport generator (kasus.go).
	if f, ada := sikapKhusus[nama]; ada {
		return f(p.kasus)
	}
	pr, ada := registry[nama]
	if !ada {
		panic(fmt.Sprintf("predikat: predikat %q tidak ada di registry", nama))
	}
	if pr.panik != "" {
		panic(fmt.Sprintf("predikat: %s - %s", nama, pr.panik))
	}
	if p.sedang[nama] {
		panic(fmt.Sprintf("predikat: rujukan melingkar lewat %s", nama))
	}
	p.sedang[nama] = true
	defer delete(p.sedang, nama)

	pohon, ada := pohonLogika[nama]
	if !ada {
		// Predikat yang ditambahkan setelah paket dimuat (hanya uji).
		var err error
		if pohon, err = logika.Urai(pr.logika); err != nil {
			panic(fmt.Sprintf("predikat: logika %s %q: %v", nama, pr.logika, err))
		}
	}
	hasil := map[string]bool{}
	for _, label := range pohon.Label() {
		b, err := p.kondisi(pr.kondisi[label])
		if err != nil {
			return false, fmt.Errorf("%s baris %s: %w", nama, label, err)
		}
		hasil[label] = b
	}
	return pohon.Nilai(hasil), nil
}

func (p *penilai) kondisi(kd kondisi) (bool, error) {
	switch kd.jenis {
	case rujukWhen:
		return p.eval(kd.rujukan)
	case banding:
		kiri, _ := p.kasus.Nilai(kd.kiri)
		sama, err := samaDengan(kiri, kd.kanan.teks, kd.kanan.jenis == teksLiteral)
		if err != nil {
			return false, err
		}
		if kd.tidakSama {
			return !sama, nil
		}
		return sama, nil
	}
	panic(fmt.Sprintf("predikat: jenis kondisi %d tidak dikenal", kd.jenis))
}

// samaDengan - `compareTwoValues` dengan operator `=`. Tafsir teks selalu
// dihitung; tafsir angka bila kedua sisi terbaca angka. ErrTafsirBerbeda bila
// keduanya ada dan berbeda, atau bila tepat satu sisi kosong sementara sisi
// lain ANGKA. Literal TEKS berkutip (`= "10053"`) bertipe teks di rule-nya
// sendiri, sehingga kosong dibandingkan dengannya sebagai teks (butir 20).
func samaDengan(kiri, kanan string, kananTeks bool) (bool, error) {
	teks := kiri == kanan
	dKiri, errKiri := utils.ParseDecimal(kiri)
	dKanan, errKanan := utils.ParseDecimal(kanan)
	kosongLawanAngka := (kiri == "" && errKanan == nil) || (kanan == "" && errKiri == nil)
	switch {
	case errKiri == nil && errKanan == nil:
		if angka := dKiri.Cmp(dKanan) == 0; angka != teks {
			return false, fmt.Errorf("%w: %q dan %q", ErrTafsirBerbeda, kiri, kanan)
		}
	case !kananTeks && kosongLawanAngka:
		return false, fmt.Errorf("%w: nilai kosong dibandingkan dengan angka (%q dan %q)", ErrTafsirBerbeda, kiri, kanan)
	}
	return teks, nil
}
