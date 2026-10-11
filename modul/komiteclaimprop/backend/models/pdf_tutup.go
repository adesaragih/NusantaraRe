package models

// Untuk apa berkas ini: PDF CLOSE / REJECT CLAIM - KomitePost_Close / KomitePost_Reject S12.2-S12.5 (TT 4 / TT 3,
// disetujui MAUPUN ditolak). Reject = stream `CommitteReject_CC` (ClaimTreaty, diekspor work owner 10-10-2026):
// susunan close + tabel ESTIMATION DETAIL (`SusunDataTolak`). Close = stream
// `CommitteCloseClaim` (Rule-Obj-HTML ASM-FW-GCNMFW-Work-ClaimTreaty, diekspor work owner 10-10-2026 sebagai
// `Claim Prop/CommitteCloseClaim.xml`) atas TempOpenPage (klaim induk), pyWorkPage.Komite (SendCloseClaimToKomite 5.2)
// dan ApprovalCommite (S6) dijadikan PDF A4 portrait; kepala / kaki halaman S12.3 = PrintFileAcceptance_TKMT.
//
// `[penyimpangan sadar]` sama dengan PDF akseptasi (OQ-KCP-07 "A"): susunan stream DIGAMBAR go-pdf/fpdf (label
// VERBATIM, ukuran dari CSS stream), markup HTML tidak ditafsirkan; tanggal `format="date"` = dd/MM/yyyy;
// ApprovalCommite.DateOfComitee XML `"dd/MM/YYYY"` (YYYY Java = tahun-minggu) dicetak dd/MM/yyyy.

import (
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/kontrak"
)

// Teks tetap S12.3 dan stream CommitteCloseClaim (VERBATIM).
const (
	// KategoriDokumenTutup - S12.3 `param.AttachmentCategory`.
	KategoriDokumenTutup = "CloseClaim"
	tandaPDFTutup        = "Close Claim"
	tandaPDFTolak        = "Reject Claim" // CommitteReject_CC (ClaimTreaty)
	statusSetujuPDF      = "Approve"      // S6: keputusan anggota "1" -> Approve, selain itu Reject
	statusTolakPDF       = "Reject"
	tinggiPenutupTutup   = 18 * pxMM // tabel penutup style="line-height:18px"
	jarakBingkaiTutup    = 1.5       // padding sel berbingkai td1 / table1 (1px + 1px), dibulatkan agar teks tidak menempel
)

// DataTutupKlaim - isi `pega:reference` stream CommitteCloseClaim / CommitteReject_CC (versi ClaimTreaty).
type DataTutupKlaim struct {
	// Tanda - teks miring kanan atas ("Close Claim" / "Reject Claim").
	Tanda string
	// TempOpenPage (klaim induk, KomitePost_Close S1)
	SOB, Ceding, InsuredName, PolicyNo, PeriodeAwal, PeriodeAkhir, NoClaim, DateOfLoss, CauseOfLoss, Location, PIC string
	// pyWorkPage.Komite: CircumtansesCouseOfLoss, Remarks, Initial (OperatorID.pyUserName peminta), DateOfComitee
	// (tanggal permintaan, @CurrentDate)
	Kronologi, Remarks, Peminta, TanggalMinta string
	// ApprovalCommite (S6): IDKomite, KomiteAproval, DateOfComitee (tanggal putusan), KomiteComment
	Komite, Status, TanggalPutus, KomentarKomite string
	// TabelEstimasi - tabel ESTIMATION DETAIL (hanya CommitteReject_CC; CommitteCloseClaim hanya judulnya).
	TabelEstimasi bool
	Estimasi      []BarisEstimasiPDF
}

// BarisEstimasiPDF - satu baris tabel ESTIMATION DETAIL CommitteReject_CC (`TempOpenPage.ClaimData.EstimationList`).
type BarisEstimasiPDF struct {
	No, TypeLoss, Tanggal, Nilai string
}

// NamaBerkasTutup - S12.3 `Param.PDFName = "Close Claim "+TempOpenPage.ClaimData.NoClaim+".PDF"` (VERBATIM).
func NamaBerkasTutup(noClaim string) string { return "Close Claim " + noClaim + ".PDF" }

// NamaBerkasTolak - KomitePost_Reject S12.3 `"Reject Claim "+TempOpenPage.ClaimData.NoClaim+".PDF"` (VERBATIM;
// kategori tetap `CloseClaim`, sama dengan close).
func NamaBerkasTolak(noClaim string) string { return "Reject Claim " + noClaim + ".PDF" }

// SusunDataTolak = KomitePost_Reject S5-S6 + stream CommitteReject_CC (ClaimTreaty, diekspor work owner 10-10-2026):
// isian sama dengan CommitteCloseClaim ditambah tabel ESTIMATION DETAIL dari `ClaimData.EstimationList` (TypeLoss,
// EstimationDate dd/MM/yyyy, Currency + EstimationValue `DecimalFormat(",###.0000")`). `[penyimpangan sadar]` kolom
// Estimation No: stream mengulang `a = 1` di dalam loop (setiap baris bernomor 1) - di sini 1..n.
func SusunDataTolak(kl kontrak.KlaimTreaty, k Kasus, a Anggota, saat time.Time) (DataTutupKlaim, error) {
	d := SusunDataTutup(kl, k, a, saat)
	d.Tanda, d.TabelEstimasi = tandaPDFTolak, true
	for i, e := range kl.Daftar["ClaimData.EstimationList"] {
		nilai, err := angkaPDFTolak(fmt.Sprintf("ClaimData.EstimationList(%d).EstimationValue", i+1), e["EstimationValue"])
		if err != nil {
			return DataTutupKlaim{}, err
		}
		d.Estimasi = append(d.Estimasi, BarisEstimasiPDF{No: fmt.Sprint(i + 1), TypeLoss: e["TypeLoss"],
			Tanggal: TanggalSurat(e["EstimationDate"]), Nilai: strings.TrimSpace(e["Currency"] + " " + nilai)})
	}
	return d, nil
}

// angkaPDFTolak = `new DecimalFormat(",###.0000")` CommitteReject_CC tanpa locale: locale bawaan JVM Pega - `[dugaan]`
// en_US (ribuan ",", desimal "."), tepat empat desimal HALF_EVEN, bagian bulat nol tidak dicetak (0.5 -> ".5000").
// Kosong tetap kosong (Pega: `Double.parseDouble("")` menggagalkan PDF).
func angkaPDFTolak(jalur, s string) (string, error) {
	if strings.TrimSpace(s) == "" {
		return "", nil
	}
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
	t := r.Text('f')
	tanda := ""
	if strings.HasPrefix(t, "-") {
		tanda, t = "-", t[1:]
	}
	bulat, pecahan, _ := strings.Cut(t, ".")
	if bulat == "0" {
		bulat = ""
	}
	var b strings.Builder
	for i, c := range bulat {
		if i > 0 && (len(bulat)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	return tanda + b.String() + "." + pecahan, nil
}

// SusunDataTutup = KomitePost_Close S5-S6 + referensi stream: klaim `kl`, kasus komite TT 4 `k`, baris tangga yang
// baru memutus `a` (`KomiteList(Local.Count)`), `saat` = DateApprove (S5 @CurrentDateTime).
func SusunDataTutup(kl kontrak.KlaimTreaty, k Kasus, a Anggota, saat time.Time) DataTutupKlaim {
	v := kl.Nilai
	status := statusTolakPDF
	if a.Keputusan == KeputusanSetuju {
		status = statusSetujuPDF
	}
	minta := ""
	if !k.TglCreate.IsZero() {
		minta = TanggalSurat(k.TglCreate.In(Jakarta).Format("20060102"))
	}
	return DataTutupKlaim{Tanda: tandaPDFTutup,
		SOB: v["TreatyInMaster.LeadingReinsSource"], Ceding: v["TreatyInMaster.Ceding"],
		InsuredName: v["ClaimData.InsuredName"], PolicyNo: v["ClaimData.PolicyData.PolicyNo"],
		PeriodeAwal:  TanggalSurat(v["ClaimData.PolicyData.StartDateTime"]),
		PeriodeAkhir: TanggalSurat(v["ClaimData.PolicyData.EndDateTime"]),
		NoClaim:      v["ClaimData.NoClaim"], DateOfLoss: TanggalSurat(v["ClaimData.DateOfLoss"]),
		CauseOfLoss: v["ClaimData.CauseOfLoss"], Location: v["ClaimData.Location"],
		PIC: v[kontrak.JalurNamaPembuatKlaimTreaty], // TempOpenPage.pxCreateOpName
		// Komite.Remarks = TempCommiteClaim.Remarks saat permintaan = ClaimData.Remark_Close (SendCloseClaimToKomite 2)
		Kronologi: k.Kronologi, Remarks: v["ClaimData.Remark_Close"], Peminta: k.PembuatNama, TanggalMinta: minta,
		Komite: a.Jabatan, Status: status, TanggalPutus: saat.In(Jakarta).Format("02/01/2006"), KomentarKomite: a.Komentar,
	}
}

// PDFTutupKlaim = S12.4-S12.5: PDF stream CommitteCloseClaim / CommitteReject_CC (`d.TabelEstimasi`). `saat` =
// tanggal dokumen PDF.
func PDFTutupKlaim(d DataTutupKlaim, saat time.Time) ([]byte, error) {
	k := kertasPDF(saat)
	k.isiTutup(d)
	stream := "CommitteCloseClaim"
	if d.TabelEstimasi {
		stream = "CommitteReject_CC"
	}
	return k.keluaran(stream)
}

func (k *kertasAkseptasi) isiTutup(d DataTutupKlaim) {
	p := k.pdf
	p.SetFont(fonDokumen, "I", ptBadan) // <td align=right><i>Close Claim</i>
	p.CellFormat(k.lebar, tinggiBaris, k.tr(d.Tanda), "", 1, "R", false, 0, "")
	k.judulTengah("DETAILS CLAIM")
	k.kotakBaris([][]string{
		{"SOB", d.SOB},
		{"Ceding", d.Ceding},
		{"Name of Insured", d.InsuredName},
		{"Policy No.", d.PolicyNo},
		{"Period of Policy", d.PeriodeAwal + " - " + d.PeriodeAkhir},
		{"Claim No.", d.NoClaim},
		{"Date of Loss", d.DateOfLoss},
		{"Cause of Loss", d.CauseOfLoss},
		{"Location of Loss", d.Location},
		{"PIC Name", d.PIC},
	})
	p.Ln(tinggiBaris / 2) // </table><br>

	k.judulTengah("COMMENT AND ANALYSIS")
	k.baris("• Chronology", d.Kronologi)
	k.baris("• Remarks", d.Remarks)
	k.baris("• Request Reject", d.Peminta)
	k.baris("• Request Date", d.TanggalMinta)
	p.Ln(tinggiBaris / 2)

	k.judulTengah("COMMITTEE COMMENT")
	k.baris("• Committee", d.Komite)
	k.baris("• Status", d.Status)
	k.baris("• Date Approve", d.TanggalPutus)
	k.baris("• Remarks", d.KomentarKomite)
	p.Ln(tinggiBaris / 2)

	k.judulTengah("ESTIMATION DETAIL") // CommitteCloseClaim: judul tanpa isi
	if d.TabelEstimasi {
		k.tabelEstimasi(d.Estimasi)
	}

	h := tinggiPenutupTutup
	k.ruang(h * 7) // <br><br> Jakarta, tanggal <br><br><br> perusahaan - tidak dipotong antarhalaman
	p.Ln(h * 2)
	p.SetX(marginKiriPDF)
	p.CellFormat(k.lebar, h, k.tr("Jakarta, "+d.TanggalPutus), "", 1, "L", false, 0, "")
	p.Ln(h * 2)
	p.SetX(marginKiriPDF)
	p.CellFormat(k.lebar, h, k.tr(penutupPerusahaan), "", 1, "L", false, 0, "")
}

// tabelEstimasi - CommitteReject_CC `table.tablet01`: kepala th02 (latar #eee) "Estimation No" / (kosong) / "Estimation
// Date" / "Estimation Value" (kolom Type ber-komentar HTML, tidak tampil), lalu satu baris per estimasi.
func (k *kertasAkseptasi) tabelEstimasi(rows []BarisEstimasiPDF) {
	p := k.pdf
	w := []float64{k.lebar * 0.18, k.lebar * 0.32, k.lebar * 0.22, k.lebar * 0.28}
	k.ruang(tinggiBaris * 2)
	p.SetFont(fonDokumen, "B", ptBadan)
	p.SetFillColor(0xee, 0xee, 0xee)
	p.SetX(marginKiriPDF)
	for i, h := range []string{"Estimation No", "", "Estimation Date", "Estimation Value"} {
		p.CellFormat(w[i], tinggiBaris, k.tr(h), "", 0, "L", true, 0, "")
	}
	p.Ln(tinggiBaris)
	p.SetFont(fonDokumen, "", ptBadan)
	for _, r := range rows {
		k.barisTabel(marginKiriPDF, w, []string{r.No, r.TypeLoss, r.Tanggal, r.Nilai})
	}
}

// judulTengah - `<td colspan="3" align=center><b>...</b></td>`.
func (k *kertasAkseptasi) judulTengah(s string) {
	k.ruang(tinggiBaris * 2) // judul tidak terpisah dari baris pertamanya
	k.pdf.SetFont(fonDokumen, "B", ptBadan)
	k.pdf.SetX(marginKiriPDF)
	k.pdf.CellFormat(k.lebar, tinggiBaris, k.tr(s), "", 1, "C", false, 0, "")
	k.pdf.SetFont(fonDokumen, "", ptBadan)
}

// kotakBaris - tabel `table1` berbingkai 1px #999 berisi baris label : nilai (td2 30%); seluruh kotak di satu halaman.
func (k *kertasAkseptasi) kotakBaris(rows [][]string) {
	p := k.pdf
	x, lebar := marginKiriPDF+jarakBingkaiTutup, k.lebar-2*jarakBingkaiTutup
	n := 0
	for _, r := range rows {
		n += len(k.pecahNilai(lebar, r[1:]...))
	}
	tinggi := float64(n)*tinggiBaris + 2*jarakBingkaiTutup
	k.ruang(tinggi)
	y0 := p.GetY()
	p.SetY(y0 + jarakBingkaiTutup)
	for _, r := range rows {
		k.barisDi(x, lebar, r[0], r[1:]...)
	}
	p.SetDrawColor(0x99, 0x99, 0x99)
	p.SetLineWidth(pxMM)
	p.Rect(marginKiriPDF, y0, k.lebar, tinggi, "D")
	p.SetDrawColor(0, 0, 0)
	p.SetXY(marginKiriPDF, y0+tinggi)
}
