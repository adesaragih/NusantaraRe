package models_test

import (
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimprop/backend/models"
)

// DecimalFormat("#,###.####") locale id_ID - nilai harapan dihitung tangan dari aturan Java (HALF_EVEN, tanpa '0').
func TestAngkaPegaDecimalFormatIndonesia(t *testing.T) {
	for masuk, mau := range map[string]string{
		"0": "0", "": "0", "100": "100", "1234567.891": "1.234.567,891", "0.5": ",5", "-1234.5": "-1.234,5",
		"1.23456": "1,2346", "0.00015": ",0002", "0.00025": ",0002", "25000000.00000000": "25.000.000",
	} {
		if got, err := models.AngkaPega("uji", masuk); err != nil || got != mau {
			t.Errorf("AngkaPega(%q) = %q, %v; mau %q", masuk, got, err, mau)
		}
	}
	if _, err := models.AngkaPega("uji", "UJI"); err == nil {
		t.Fatal("bukan angka harus galat")
	}
}

func klaimSurat() kontrak.KlaimTreaty {
	return kontrak.KlaimTreaty{Adjustment: 1, Nilai: map[string]string{
		"pyID": "CLMP-UJI001", "ClaimData.NoClaim": "UJI/CLM/001", "ClaimData.DateOfLoss": "20260912",
		"ClaimData.PolicyData.PolicyNo": "UJI-POL-1", "ClaimData.PolicyData.StartDateTime": "20260101T000000.000 GMT",
		"ClaimData.PolicyData.EndDateTime": "20261231", "ClaimData.InsuredName": "UJI <Tertanggung>",
		"ClaimData.CauseOfLoss": "UJI KEBAKARAN", "ClaimData.Location": "UJI LOKASI",
		"TreatyInMaster.Ceding": "UJI CEDING", "TreatyInMaster.LeadingReinsSource": "UJI SOB",
		"OfferFacIn.QuotationData.BusinessName": "UJI COB",
	}, Daftar: map[string][]map[string]string{
		"ClaimData.AdjustmentList": {{"KomiteNo": "UJI-KNO"}},
		"ClaimData.AdjustmentList(1).SpreadingAdjustment": {
			{"TreatyName": "UJI QS", "TreatyType": "10196", "Currency": "IDR", "SharePercentage": "60.5", "ClaimSpreaded": "605000.25"},
			{"TreatyName": "UJI SURPLUS", "TreatyType": "UJI-SP", "Currency": "IDR", "SharePercentage": "39.5", "ClaimSpreaded": "395000"},
		},
		"ClaimData.AdjustmentList(1).SpreadingQuotaShare": {
			{"TreatyName": "UJI RETRO QS", "SharePercentage": "50", "ClaimSpreaded": "302500.125"},
		},
		"ClaimData.InterestList": {
			{"ObjectName": "UJI GEDUNG", "Currency": "IDR", "TSIPerObject": "1000000000"},
			{"ObjectName": "UJI MESIN", "Currency": "USD", "TSIPerObject": "2500.5"},
		},
	}}
}

// SendEmailKlaim_KMT S5-S16 jalur CLMP + EmailKlaim_HTML_KMT.
func TestEmailKlaimKMT(t *testing.T) {
	d, err := models.SusunDataEmail(klaimSurat(), models.BahanEmail{Jenis: models.EmailPenyetujuBerikut,
		KomiteID: "TKMT-UJI001", Keputusan: "1", Komentar: "UJI <catatan>", Sapaan: "UJI-JABATAN-2",
		Pengirim: "UJI Penyetuju Satu", AcceptedNo: "UJI-NO"})
	if err != nil {
		t.Fatal(err)
	}
	// S11.1 @divide(.SharePercentage,1,0) HALF_UP: 60.5 -> 61, 39.5 -> 40; S10 total 100; ClaimSpreaded #,###.####
	if d.KomiteNo != "UJI-KNO" || d.Status != "Accepted" || d.AcceptedNo != "" || d.TotalShare != "100" ||
		d.TotalClaim != "1.000.000,25" || len(d.Baris) != 2 || d.Baris[0].Share != "61" || d.Baris[1].Share != "40" ||
		d.Baris[0].ClaimSpreaded != "605.000,25" || d.Baris[1].No != "2" || d.DateOfLoss != "12/09/2026" {
		t.Fatalf("data email: %+v", d)
	}
	h, err := models.RenderEmailKomite(d)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"Dear <strong>UJI-JABATAN-2</strong>,", ": </span>CLMP-UJI001 / UJI-KNO</td>",
		"Total Share Percentage (%): 100", "Total Claim Spreaded: 1.000.000,25", "UJI &lt;Tertanggung&gt;",
		"UJI &lt;catatan&gt;", "<strong>UJI Penyetuju Satu</strong>", "@media only screen and (max-width: 600px)"} {
		if !strings.Contains(h, s) {
			t.Errorf("badan email tanpa %q", s)
		}
	}
	if strings.Contains(h, "Accepted No:") || strings.Contains(h, "pega:") || strings.Contains(h, "<%") {
		t.Fatal("S12 tanpa Accepted No; tanpa sisa JSP / pega:reference")
	}
	d, _ = models.SusunDataEmail(klaimSurat(), models.BahanEmail{Jenis: models.EmailPembuatTolak, Keputusan: "2",
		AcceptedNo: "UJI-NO"})
	if h, _ = models.RenderEmailKomite(d); d.Status != "Rejected" || !strings.Contains(h, "Accepted No: <span") {
		t.Fatalf("S15 Rejected + CARI16: %+v", d)
	}
	cfg := models.KonfigurasiEmail{Akun: "UJI-AKUN", AkunSyariah: "UJI-SYARIAH"}
	if cfg.AkunUntuk("uji.syariah@contoh.invalid") != "UJI-SYARIAH" || cfg.AkunUntuk("uji@contoh.invalid") != "UJI-AKUN" {
		t.Fatal("S17-S18 akun notifikasi")
	}
}

// PrintFileAcceptance_TKMT S5-S10 + FILEAcceptanceNote.
func TestFileAcceptanceNote(t *testing.T) {
	adj := map[string]string{"AcceptedNo": "UJIA12.10.2026.TP00001", "Currency": "IDR", "AdjustmentValue": "1000000.25",
		"PayableTo": "UJI PENERIMA", "NameOfBank": "UJI BANK", "BranchOfBank": "UJI CABANG", "NoAccount": "UJI-REK",
		"pxCreateOpName": "UJI Admin"}
	saat := time.Date(2026, 10, 7, 20, 30, 0, 0, time.UTC) // 08-10-2026 03:30 Jakarta
	d, err := models.SusunAcceptanceNote(klaimSurat(), adj, saat, "UJI Penyetuju Akhir")
	if err != nil {
		t.Fatal(err)
	}
	if d.Location != "UJI KEBAKARAN" || !d.TampilLini || !d.TampilQS || len(d.QS) != 1 || d.TanggalCetak != "08 October 2026" ||
		d.PeriodeAwal != "01/01/2026" || d.InterestInsured != "UJI GEDUNG" || d.TSI[1].Nilai != "2.500,5" {
		t.Fatalf("data dokumen: %+v", d)
	}
	h, err := models.RenderAcceptanceNote(d)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"ACCEPTED CLAIM INSURANCE", "Accepted No :  \tUJIA12.10.2026.TP00001", "UJI COB",
		"UJI/CLM/001 / CLMP-UJI001", "UJI GEDUNG&nbsp;IDR&nbsp;1.000.000.000", "IDR&nbsp;1.000.000,25",
		"BreakDown Spreading (QS)", "<td>UJI RETRO QS</td>", "<td> 60,5</td>", "Jakarta, 08 October 2026",
		"Create by : UJI Penyetuju Akhir", "UJI &lt;Tertanggung&gt;"} {
		if !strings.Contains(h, s) {
			t.Errorf("dokumen tanpa %q", s)
		}
	}
	if strings.Contains(h, "Swift Code") {
		t.Fatal("Swift Code hanya bila terisi")
	}
	if n := models.NamaBerkasAkseptasi("UJI-NO"); n != "Persetujuan Klaim   AcceptNo UJI-NO.pdf" {
		t.Fatalf("S9 PDFName: %q", n)
	}
}
