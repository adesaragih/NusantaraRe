package models

// Untuk apa berkas ini: PDF DOKUMEN AKSEPTASI - `PrintPDFAccep_MultiAksep_KMT` S24-S26 (KomitePost_Adjustment S9):
// stream `AcceptanceNotePDF` (Rule-Obj-HTML kelas Work-PNC, korpus `Claim Fac In/AcceptanceNotePDF.xml`, diberikan work
// owner 10-10-2026) atas halaman `TempAcceptedNo` dijadikan PDF A4 portrait (S16-S22 pyPDFPageSize /
// pyPDFPageOrientation) dengan kepala dan kaki halaman S23 (pyPDFHeaderHTMLTemplate / pyPDFFooterHTMLTemplate).
//
// `[penyimpangan sadar]` (prompt tahap 2 §6 butir 7 "gambar dengan go-pdf/fpdf", pola Komite Claim Prop):
// `HTMLToPDF` memakai mesin HTML->PDF platform Pega; padanannya pustaka Go murni `github.com/go-pdf/fpdf` yang
// MENGGAMBAR susunan stream (judul, baris label : nilai, tabel di sel nilai, penutup) - markup HTML-nya tidak
// ditafsirkan. Isi dan urutan baris
// disusun `SusunAcceptanceNote` (syarat `pega:when` stream); teks label VERBATIM dari stream; ukuran dari CSS stream
// (Times New Roman 18px, line-height 2rem, margin kiri 60px / kanan 40px). Huruf di luar cp1252 tidak dapat dicetak fon
// dasar PDF.

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"
	"github.com/go-pdf/fpdf"
)

// Ukuran CSS stream dalam milimeter (1px = 25,4/96 mm) dan poin (1px = 0,75pt).
const (
	pxMM           = 25.4 / 96
	ptBadan        = 18 * 0.75       // body font-size 18px
	ptJudul        = 18 * 1.5 * .75  // <h2> = 1.5em
	ptNomor        = 18 * 1.17 * .75 // <h3> = 1.17em
	ptKepala       = 17 * 0.75       // header / footer font-size 17px
	tinggiBaris    = 32 * pxMM       // line-height 2rem
	marginKiriPDF  = 60 * pxMM
	marginKananPDF = 40 * pxMM
	marginAtasPDF  = 25.0 // kepala halaman (padding-bottom 30px) di atas badan
	marginBawahPDF = 18.0
	fonDokumen     = "Times" // font-family: Times New Roman
)

// KepalaPDFAkseptasi - teks kepala halaman S23 `pyPDFHeaderHTMLTemplate` (VERBATIM).
const KepalaPDFAkseptasi = "PT. REASURANSI NUSANTARA MAKMUR"

// BarisPDF - satu <tr> tabel utama stream: label 30%, ":", nilai (setiap paragraf boleh membungkus) atau tabel di sel
// nilai (Spreading Adjustment / BreakDown Spreading (QS) / TSI RNM).
type BarisPDF struct {
	Label string
	Nilai []string
	Tabel *TabelPDF
}

// TabelPDF - tabel tanpa garis di sel nilai (lebar kolom menurut isi, `tabel`).
type TabelPDF struct {
	Kepala []string
	Baris  [][]string
}

// DokumenPDF - isi dokumen yang sudah disusun (urut stream).
type DokumenPDF struct {
	Judul   []string // <h2> (pega:when per lini)
	Nomor   string   // <h3>Accepted No : ...
	Baris   []BarisPDF
	Penutup string // "Jakarta, dd MMMM yyyy"
}

// kertas - satu dokumen yang sedang digambar.
type kertas struct {
	pdf   *fpdf.Fpdf
	tr    func(string) string
	lebar float64 // lebar badan (100% tabel stream)
}

// spasiHTML - pemisah baris / tab di nilai halaman tampil sebagai satu spasi (HTML meruntuhkan spasi putih).
var spasiHTML = strings.NewReplacer("\r\n", " ", "\n", " ", "\r", " ", "\t", " ")

// PDFDokumen = S24-S26: PDF dokumen akseptasi. `saat` = tanggal dokumen PDF (S3 / S14 NameInput.CARI24).
func PDFDokumen(d DokumenPDF, saat time.Time) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginKiriPDF, marginAtasPDF, marginKananPDF)
	pdf.SetAutoPageBreak(false, marginBawahPDF) // pindah halaman diatur `ruang` per baris
	pdf.SetCreationDate(saat)
	pdf.SetModificationDate(saat)
	pdf.SetCatalogSort(true) // urutan katalog fon tetap: berkas yang sama untuk masukan yang sama
	pdf.AliasNbPages("")
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	k := &kertas{pdf: pdf, tr: func(s string) string { return tr(spasiHTML.Replace(s)) }}
	lebarHalaman, _ := pdf.GetPageSize()
	k.lebar = lebarHalaman - marginKiriPDF - marginKananPDF
	pdf.SetHeaderFuncMode(func() { // S23 pyPDFHeaderHTMLTemplate; sesudahnya badan mulai di margin atas
		pdf.SetFont(fonDokumen, "", ptKepala)
		pdf.SetXY(marginKiriPDF, 10)
		pdf.CellFormat(k.lebar, 17*pxMM*1.2, k.tr(KepalaPDFAkseptasi), "", 0, "L", false, 0, "")
	}, true)
	pdf.SetFooterFunc(func() { // S23 pyPDFFooterHTMLTemplate "Page $[page] of $[total] pages"
		pdf.SetFont(fonDokumen, "", ptKepala)
		pdf.SetY(-marginBawahPDF + 4)
		pdf.CellFormat(k.lebar*0.9, 17*pxMM*1.2, fmt.Sprintf("Page %d of {nb} pages", pdf.PageNo()), "", 0, "L",
			false, 0, "")
	})
	pdf.AddPage()
	k.isi(d)
	var b bytes.Buffer
	if err := pdf.Output(&b); err != nil {
		return nil, fmt.Errorf("models: PDF AcceptanceNotePDF: %w", err)
	}
	return b.Bytes(), nil
}

func (k *kertas) isi(d DokumenPDF) {
	p := k.pdf
	p.SetFont(fonDokumen, "B", ptJudul) // <center><h2>
	for _, j := range d.Judul {
		p.CellFormat(k.lebar, tinggiBaris*1.2, k.tr(j), "", 1, "C", false, 0, "")
	}
	p.SetFont(fonDokumen, "B", ptNomor) // <h3>Accepted No :
	p.CellFormat(k.lebar, tinggiBaris, k.tr("Accepted No : "+d.Nomor), "", 1, "C", false, 0, "")
	p.SetFont(fonDokumen, "", ptBadan)
	p.Ln(tinggiBaris / 2) // <br>
	for _, b := range d.Baris {
		if b.Tabel != nil { // tabel di sel nilai: baris kepalanya sejajar label
			k.ruang(tinggiBaris * 2)
			if b.Label != "" {
				k.label(b.Label)
			}
			k.tabel(marginKiriPDF+k.lebar*0.33, k.lebar*0.67, *b.Tabel)
			continue
		}
		k.baris(b.Label, b.Nilai...)
	}
	k.ruang(tinggiBaris * 2)
	p.Ln(tinggiBaris / 2)
	// <table style="text-align:right; margin-right:50px; width:100%"><td style="text-align:right">
	p.SetX(marginKiriPDF)
	p.CellFormat(k.lebar, tinggiBaris, k.tr(d.Penutup), "", 1, "R", false, 0, "")
}

// label - sel label 30% dan ":" sebuah <tr> di baris berjalan, tanpa berpindah baris.
func (k *kertas) label(label string) {
	p := k.pdf
	y := p.GetY()
	p.SetXY(marginKiriPDF, y)
	p.CellFormat(k.lebar*0.30, tinggiBaris, k.tr(label), "", 0, "L", false, 0, "")
	p.CellFormat(k.lebar*0.03, tinggiBaris, ":", "", 0, "L", false, 0, "")
	p.SetXY(marginKiriPDF, y)
}

// baris - satu <tr> tabel stream: label 30%, ":", nilai (setiap paragraf membungkus; nilai yang lebih panjang dari
// satu halaman berlanjut di halaman berikut).
func (k *kertas) baris(label string, nilai ...string) {
	p := k.pdf
	xNilai, wNilai := marginKiriPDF+k.lebar*0.33, k.lebar*0.67
	var baris []string
	for _, n := range nilai {
		baris = append(baris, k.pecah(k.tr(n), wNilai)...)
	}
	if len(baris) == 0 {
		baris = []string{""}
	}
	k.utuh(tinggiBaris * float64(len(baris)))
	k.label(label)
	for _, b := range baris {
		k.ruang(tinggiBaris)
		y := p.GetY()
		p.SetXY(xNilai, y)
		p.CellFormat(wNilai, tinggiBaris, b, "", 0, "L", false, 0, "")
		p.SetXY(marginKiriPDF, y+tinggiBaris)
	}
}

// selaTabel - jarak antarkolom tabel di luar margin sel fpdf (bawaan `cellspacing` HTML, mm).
const selaTabel = 1

// tabel - tabel tanpa garis (seperti stream) mulai `x`, paling lebar `lebar`. Stream tidak memberi lebar kolom
// (`<table>` tanpa width): lebar kolom = isi terpanjangnya, diperkecil berimbang bila jumlahnya melebihi `lebar`.
func (k *kertas) tabel(x, lebar float64, t TabelPDF) {
	n := len(t.Kepala)
	if n == 0 {
		return
	}
	w := make([]float64, n)
	total := 0.0
	for i := range w {
		w[i] = k.pdf.GetStringWidth(k.tr(t.Kepala[i]))
		for _, r := range t.Baris {
			if i < len(r) {
				w[i] = max(w[i], k.pdf.GetStringWidth(k.tr(r[i])))
			}
		}
		w[i] += 2*k.pdf.GetCellMargin() + selaTabel
		total += w[i]
	}
	if total > lebar {
		for i := range w {
			w[i] *= lebar / total
		}
	}
	k.barisTabel(x, w, t.Kepala)
	for _, r := range t.Baris {
		k.barisTabel(x, w, r)
	}
}

func (k *kertas) barisTabel(x float64, w []float64, sel []string) {
	p := k.pdf
	pecah := make([][]string, len(w))
	n := 1
	for i := range w {
		s := ""
		if i < len(sel) {
			s = sel[i]
		}
		pecah[i] = k.pecah(k.tr(s), w[i])
		n = max(n, len(pecah[i]))
	}
	k.utuh(tinggiBaris * float64(n))
	for j := 0; j < n; j++ {
		k.ruang(tinggiBaris)
		y := p.GetY()
		xi := x
		for i := range w {
			if j < len(pecah[i]) {
				p.SetXY(xi, y)
				p.CellFormat(w[i], tinggiBaris, pecah[i][j], "", 0, "L", false, 0, "")
			}
			xi += w[i]
		}
		p.SetXY(marginKiriPDF, y+tinggiBaris)
	}
}

// pecah - teks `s` (SUDAH cp1252, `k.tr`) dipecah menjadi baris selebar paling banyak `w`: di spasi terakhir yang
// muat; kata yang lebih lebar dari `w` dipotong per huruf. Pengganti `fpdf.SplitText`, yang membaca teks sebagai rune
// UTF-8 dan panik pada huruf cp1252 bukan ASCII (é, –, ’ menjadi U+FFFD di luar tabel lebar fon).
func (k *kertas) pecah(s string, w float64) []string {
	maks := w - 2*k.pdf.GetCellMargin()
	var out []string
	for k.pdf.GetStringWidth(s) > maks {
		n := 1 // s[:n] terpanjang yang muat, paling sedikit satu huruf
		for n < len(s) && k.pdf.GetStringWidth(s[:n+1]) <= maks {
			n++
		}
		potong, sisa := n, n
		if i := strings.LastIndexByte(s[:n+1], ' '); i > 0 {
			potong, sisa = i, i+1
		}
		out = append(out, s[:potong])
		s = s[sisa:]
	}
	return append(out, s)
}

// utuh - satu <tr> setinggi `h` tidak dipotong antarhalaman bila muat satu halaman; yang lebih tinggi mulai di sini dan
// berlanjut di halaman berikut (`ruang` per baris).
func (k *kertas) utuh(h float64) {
	_, tinggiHalaman := k.pdf.GetPageSize()
	if h <= tinggiHalaman-marginAtasPDF-marginBawahPDF {
		k.ruang(h)
	}
}

// ruang - pindah halaman bila tinggi `h` tidak muat.
func (k *kertas) ruang(h float64) {
	_, tinggiHalaman := k.pdf.GetPageSize()
	if k.pdf.GetY()+h > tinggiHalaman-marginBawahPDF {
		k.pdf.AddPage()
	}
}

// AngkaPega = `DecimalFormat("#,###.####")` locale id_ID (formatdesimal1 / formatdesimal2 stream): ribuan ".", desimal
// ",", paling banyak empat angka desimal dibulatkan HALF_EVEN (bawaan DecimalFormat), tanpa nol ekor. Pola tanpa digit
// '0': bagian bulat nol tidak dicetak (0,5 -> ",5"); nol -> "0". Kosong = 0. Disalin dari Komite Claim Prop.
// `[penyimpangan sadar]` stream memformat `double` (`getDouble`; Java 8+ membulatkan nilai biner tepatnya): di sini
// nilai DESIMAL tersimpan dibulatkan (nol float untuk uang, ADR-0003) - beda hanya pada seri tepat ...5 di desimal
// ke-5 (0,00015 -> Java ",0001") dan nilai di atas 15 digit bermakna.
func AngkaPega(jalur, s string) (string, error) {
	d, err := desimal(jalur, s)
	if err != nil {
		return "", err
	}
	c := *konteksUang
	c.Rounding = apd.RoundHalfEven
	r := new(apd.Decimal)
	if _, err := c.Quantize(r, d, -4); err != nil {
		return "", fmt.Errorf("models: %s: %w", jalur, err)
	}
	r.Reduce(r)
	t := r.Text('f')
	tanda := ""
	if strings.HasPrefix(t, "-") {
		tanda, t = "-", t[1:]
	}
	bulat, pecahan, _ := strings.Cut(t, ".")
	if bulat == "0" && pecahan != "" {
		bulat = ""
	}
	var b strings.Builder
	for i, c := range bulat {
		if i > 0 && (len(bulat)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(c)
	}
	if pecahan != "" {
		b.WriteString("," + pecahan)
	}
	return tanda + b.String(), nil
}
