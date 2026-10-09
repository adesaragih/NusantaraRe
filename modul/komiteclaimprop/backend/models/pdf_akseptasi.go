package models

// Untuk apa berkas ini: PDF DOKUMEN AKSEPTASI - PrintFileAcceptance_TKMT S9-S11: stream `FILEAcceptanceNote` atas
// halaman TempAcceptedNo dijadikan PDF A4 portrait (S9 pyPDFPageSize / pyPDFPageOrientation) dengan kepala dan kaki
// halaman S9 (pyPDFHeaderHTMLTemplate / pyPDFFooterHTMLTemplate).
//
// `[penyimpangan sadar]` (keputusan work owner 08-10-2026, OQ-KCP-07 "A"): `HTMLToPDF` memakai mesin HTML->PDF platform
// Pega; padanannya pustaka Go murni `github.com/go-pdf/fpdf` yang MENGGAMBAR susunan stream (judul, baris label : nilai,
// tabel spreading, penutup) - markup HTML-nya tidak ditafsirkan. Teks label VERBATIM dari stream; ukuran dari CSS
// stream (Times New Roman 18px, line-height 2rem, margin kiri 60px / kanan 40px). Huruf di luar cp1252 tidak dapat
// dicetak fon dasar PDF.

import (
	"bytes"
	"fmt"
	"time"

	"github.com/go-pdf/fpdf"
)

// Ukuran CSS stream dalam milimeter (1px = 25,4/96 mm) dan poin (1px = 0,75pt).
const (
	pxMM           = 25.4 / 96
	ptBadan        = 18 * 0.75       // body font-size 18px
	ptJudul        = 18 * 1.17 * .75 // <h3> = 1.17em
	ptKepala       = 17 * 0.75       // header / footer font-size 17px
	tinggiBaris    = 32 * pxMM       // line-height 2rem
	marginKiriPDF  = 60 * pxMM
	marginKananPDF = 40 * pxMM
	marginAtasPDF  = 25.0 // kepala halaman (padding-bottom 30px) di atas badan
	marginBawahPDF = 18.0
	fonDokumen     = "Times" // font-family: Times New Roman
)

// Teks tetap S9 (kepala / kaki halaman) dan stream (VERBATIM).
const (
	KepalaPDFAkseptasi = "PT. REASURANSI NUSANTARA MAKMUR"
	JudulPDFAkseptasi  = "ACCEPTED CLAIM INSURANCE"
	penutupPerusahaan  = "PT. Reasuransi Nusantara Makmur"
)

// kertasAkseptasi - satu dokumen yang sedang digambar.
type kertasAkseptasi struct {
	pdf   *fpdf.Fpdf
	tr    func(string) string
	lebar float64 // lebar badan (100% tabel stream)
}

// PDFAcceptanceNote = S10-S11: PDF dokumen ACCEPTED CLAIM INSURANCE. `saat` = tanggal dokumen PDF (S4 CARI24).
func PDFAcceptanceNote(d DataAcceptanceNote, saat time.Time) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginKiriPDF, marginAtasPDF, marginKananPDF)
	pdf.SetAutoPageBreak(true, marginBawahPDF)
	pdf.SetCreationDate(saat)
	pdf.SetModificationDate(saat)
	pdf.SetCatalogSort(true) // urutan katalog fon tetap: berkas yang sama untuk masukan yang sama
	pdf.AliasNbPages("")
	k := &kertasAkseptasi{pdf: pdf, tr: pdf.UnicodeTranslatorFromDescriptor("")}
	lebarHalaman, _ := pdf.GetPageSize()
	k.lebar = lebarHalaman - marginKiriPDF - marginKananPDF
	pdf.SetHeaderFuncMode(func() { // pyPDFHeaderHTMLTemplate; sesudahnya badan mulai di margin atas
		pdf.SetFont(fonDokumen, "", ptKepala)
		pdf.SetXY(marginKiriPDF, 10)
		pdf.CellFormat(k.lebar, 17*pxMM*1.2, k.tr(KepalaPDFAkseptasi), "", 0, "L", false, 0, "")
	}, true)
	pdf.SetFooterFunc(func() { // pyPDFFooterHTMLTemplate "Page $[page] of $[total] pages"
		pdf.SetFont(fonDokumen, "", ptKepala)
		pdf.SetY(-marginBawahPDF + 4)
		pdf.CellFormat(k.lebar*0.9, 17*pxMM*1.2, fmt.Sprintf("Page %d of {nb} pages", pdf.PageNo()), "", 0, "L",
			false, 0, "")
	})
	pdf.AddPage()
	k.isi(d)
	var b bytes.Buffer
	if err := pdf.Output(&b); err != nil {
		return nil, fmt.Errorf("models: PDF FILEAcceptanceNote: %w", err)
	}
	return b.Bytes(), nil
}

func (k *kertasAkseptasi) isi(d DataAcceptanceNote) {
	p := k.pdf
	p.SetFont(fonDokumen, "BU", ptJudul) // <center><u><h3>
	p.CellFormat(k.lebar, tinggiBaris, k.tr(JudulPDFAkseptasi), "", 1, "C", false, 0, "")
	p.SetFont(fonDokumen, "B", ptBadan) // <th class="h3"><center>
	p.CellFormat(k.lebar, tinggiBaris, k.tr("Accepted No :  "+d.AcceptedNo), "", 1, "C", false, 0, "")
	p.SetFont(fonDokumen, "", ptBadan)
	p.Ln(tinggiBaris / 2) // <br>

	k.baris("Policy No", d.PolicyNo)
	if d.TampilLini {
		k.baris("Line of Business", d.LineOfBusiness)
	}
	k.baris("Claim No / Claim ID", d.NoClaim+" / "+d.ClaimID)
	k.baris("SOB Name", d.SOB)
	k.baris("Ceding Co Name", d.Ceding)
	k.baris("Name of Insured", d.InsuredName)
	k.baris("Interest Insured", d.InterestInsured)
	tsi := make([]string, 0, len(d.TSI))
	for _, t := range d.TSI {
		tsi = append(tsi, t.ObjectName+" "+t.Currency+" "+t.Nilai)
	}
	k.baris("TSI RNM", tsi...)
	k.baris("Period of Policy", d.PeriodeAwal+"- "+d.PeriodeAkhir)
	k.baris("Date of Loss", d.DateOfLoss)
	k.baris("Cause of Loss", d.CauseOfLoss)
	k.baris("Location of Loss", d.Location)
	k.baris("Accepted Claim", d.Currency+" "+d.AcceptedClaim)
	bayar := []string{d.PayableTo, "Name of Bank : " + d.NameOfBank}
	if d.SwiftCode != "" {
		bayar = append(bayar, "Swift Code : "+d.SwiftCode)
	}
	bayar = append(bayar, "Branch of Bank : "+d.BranchOfBank, "Account No : "+d.NoAccount)
	k.baris("Payable To", bayar...)
	k.baris("PIC Name", d.PIC)
	k.baris("Spreading Adjustment")
	k.tabelSpread(marginKiriPDF+k.lebar*0.33, k.lebar*0.67, d.Spreading)
	if d.TampilQS {
		// stream baris 242-275: tabel QS ditutup di luar sel nilai - tergambar dari tepi kiri.
		k.baris("BreakDown Spreading (QS)")
		k.tabelSpread(marginKiriPDF, k.lebar*0.67, d.QS)
	}

	k.ruang(tinggiBaris/2 + tinggiBaris*6) // tabel penutup (tempat, tanggal, <br> x3, perusahaan, pembuat) tidak dipotong
	p.Ln(tinggiBaris / 2)
	k.teks("Jakarta, " + d.TanggalCetak)
	p.Ln(tinggiBaris * 3) // <br><br><br>
	k.teks(penutupPerusahaan)
	k.teks("Create by : " + d.Pembuat)
}

// baris - satu <tr> tabel stream: label 30%, ":", nilai (setiap paragraf boleh membungkus).
func (k *kertasAkseptasi) baris(label string, nilai ...string) {
	p := k.pdf
	wLabel, wTitik := k.lebar*0.30, k.lebar*0.03
	wNilai := k.lebar - wLabel - wTitik
	var baris []string
	for _, n := range nilai {
		baris = append(baris, p.SplitText(k.tr(n), wNilai)...)
	}
	if len(baris) == 0 {
		baris = []string{""}
	}
	k.ruang(tinggiBaris * float64(len(baris)))
	y := p.GetY()
	p.SetXY(marginKiriPDF, y)
	p.CellFormat(wLabel, tinggiBaris, k.tr(label), "", 0, "L", false, 0, "")
	p.CellFormat(wTitik, tinggiBaris, ":", "", 0, "L", false, 0, "")
	for i, b := range baris {
		p.SetXY(marginKiriPDF+wLabel+wTitik, y+float64(i)*tinggiBaris)
		p.CellFormat(wNilai, tinggiBaris, b, "", 0, "L", false, 0, "")
	}
	p.SetXY(marginKiriPDF, y+float64(len(baris))*tinggiBaris)
}

// tabelSpread - tabel "Treaty Name / Share (%) / Claim Spread" (tanpa garis, seperti stream).
func (k *kertasAkseptasi) tabelSpread(x, lebar float64, rows []BarisSpread) {
	w := []float64{lebar * 0.5, lebar * 0.2, lebar * 0.3}
	k.barisTabel(x, w, []string{"Treaty Name", "Share (%)", "Claim Spread"})
	for _, r := range rows {
		k.barisTabel(x, w, []string{r.TreatyName, r.Share, r.ClaimSpread})
	}
}

func (k *kertasAkseptasi) barisTabel(x float64, w []float64, sel []string) {
	p := k.pdf
	pecah := make([][]string, len(sel))
	n := 1
	for i, s := range sel {
		pecah[i] = p.SplitText(k.tr(s), w[i])
		if len(pecah[i]) > n {
			n = len(pecah[i])
		}
	}
	k.ruang(tinggiBaris * float64(n))
	y := p.GetY()
	xi := x
	for i := range sel {
		for j, b := range pecah[i] {
			p.SetXY(xi, y+float64(j)*tinggiBaris)
			p.CellFormat(w[i], tinggiBaris, b, "", 0, "L", false, 0, "")
		}
		xi += w[i]
	}
	p.SetXY(marginKiriPDF, y+float64(n)*tinggiBaris)
}

func (k *kertasAkseptasi) teks(s string) {
	k.ruang(tinggiBaris)
	k.pdf.SetX(marginKiriPDF)
	k.pdf.CellFormat(k.lebar, tinggiBaris, k.tr(s), "", 1, "L", false, 0, "")
}

// ruang - pindah halaman bila tinggi `h` tidak muat (satu baris tabel tidak dipotong antarhalaman).
func (k *kertasAkseptasi) ruang(h float64) {
	_, tinggiHalaman := k.pdf.GetPageSize()
	if k.pdf.GetY()+h > tinggiHalaman-marginBawahPDF {
		k.pdf.AddPage()
	}
}
