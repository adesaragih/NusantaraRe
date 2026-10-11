package models_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimprop/backend/models"
)

// Stream CommitteCloseClaim (diekspor work owner 10-10-2026): setiap `pega:reference` dipetakan ke klaim induk
// (TempOpenPage), kasus komite TT 4 (pyWorkPage.Komite) dan baris tangga yang memutus (ApprovalCommite, S6).
func TestSusunDataTutupIkutStream(t *testing.T) {
	kl := kontrak.KlaimTreaty{Nilai: map[string]string{
		"TreatyInMaster.LeadingReinsSource": "UJI SOB", "TreatyInMaster.Ceding": "UJI CEDING",
		"ClaimData.InsuredName": "UJI TERTANGGUNG", "ClaimData.PolicyData.PolicyNo": "UJI-POL-1",
		"ClaimData.PolicyData.StartDateTime": "2026-01-01", "ClaimData.PolicyData.EndDateTime": "20261231",
		"ClaimData.NoClaim": "UJI-K-1", "ClaimData.DateOfLoss": "2026-09-15", "ClaimData.CauseOfLoss": "UJI SEBAB",
		"ClaimData.Location": "UJI LOKASI", "ClaimData.Remark_Close": "UJI REMARKS",
		kontrak.JalurNamaPembuatKlaimTreaty: "UJI PEMBUAT",
	}}
	// dibuat 20:00 UTC 9 Oktober = 10 Oktober di Jakarta; diputus 23:30 UTC 10 Oktober = 11 Oktober di Jakarta
	k := models.Kasus{Kronologi: "UJI KRONOLOGI", PembuatNama: "UJI PEMINTA",
		TglCreate: time.Date(2026, 10, 9, 20, 0, 0, 0, time.UTC)}
	a := models.Anggota{Jabatan: "Claim Dept. Head", Keputusan: models.KeputusanSetuju, Komentar: "UJI OK"}
	saat := time.Date(2026, 10, 10, 23, 30, 0, 0, time.UTC)

	d := models.SusunDataTutup(kl, k, a, saat)
	mau := models.DataTutupKlaim{Tanda: "Close Claim", SOB: "UJI SOB", Ceding: "UJI CEDING", InsuredName: "UJI TERTANGGUNG",
		PolicyNo: "UJI-POL-1", PeriodeAwal: "01/01/2026", PeriodeAkhir: "31/12/2026", NoClaim: "UJI-K-1",
		DateOfLoss: "15/09/2026", CauseOfLoss: "UJI SEBAB", Location: "UJI LOKASI", PIC: "UJI PEMBUAT",
		Kronologi: "UJI KRONOLOGI", Remarks: "UJI REMARKS", Peminta: "UJI PEMINTA", TanggalMinta: "10/10/2026",
		Komite: "Claim Dept. Head", Status: "Approve", TanggalPutus: "11/10/2026", KomentarKomite: "UJI OK"}
	if !reflect.DeepEqual(d, mau) {
		t.Fatalf("data stream:\n dapat %+v\n mau   %+v", d, mau)
	}
	a.Keputusan = models.KeputusanTolak
	if s := models.SusunDataTutup(kl, k, a, saat).Status; s != "Reject" {
		t.Fatalf("status ditolak %q", s)
	}
	if n := models.NamaBerkasTutup("UJI-K-1"); n != "Close Claim UJI-K-1.PDF" {
		t.Fatalf("nama berkas %q", n)
	}

	b1, err := models.PDFTutupKlaim(d, saat)
	if err != nil || !bytes.HasPrefix(b1, []byte("%PDF-")) {
		t.Fatalf("PDF: %v", err)
	}
	b2, _ := models.PDFTutupKlaim(d, saat)
	if !bytes.Equal(b1, b2) {
		t.Fatal("PDF tidak tetap untuk masukan yang sama")
	}
}

// Stream CommitteReject_CC versi ClaimTreaty (diekspor work owner 10-10-2026): isian close + tabel ESTIMATION DETAIL dari
// ClaimData.EstimationList; angka `DecimalFormat(",###.0000")` locale bawaan (en_US), nomor baris 1..n.
func TestSusunDataTolakTabelEstimasi(t *testing.T) {
	kl := kontrak.KlaimTreaty{Nilai: map[string]string{"ClaimData.NoClaim": "UJI-K-1"},
		Daftar: map[string][]map[string]string{"ClaimData.EstimationList": {
			{"TypeLoss": "QUOTA SHARE", "EstimationDate": "2026-09-15", "Currency": "IDR", "EstimationValue": "1234567.5"},
			{"TypeLoss": "SURPLUS", "EstimationDate": "20261001", "Currency": "USD", "EstimationValue": "0.123456"},
			{"TypeLoss": "UJI", "EstimationDate": "", "Currency": "USD", "EstimationValue": ""},
		}}}
	a := models.Anggota{Jabatan: "Claim Dept. Head", Keputusan: models.KeputusanSetuju}
	d, err := models.SusunDataTolak(kl, models.Kasus{}, a, time.Date(2026, 10, 10, 3, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	mau := []models.BarisEstimasiPDF{
		{No: "1", TypeLoss: "QUOTA SHARE", Tanggal: "15/09/2026", Nilai: "IDR 1,234,567.5000"},
		{No: "2", TypeLoss: "SURPLUS", Tanggal: "01/10/2026", Nilai: "USD .1235"},
		{No: "3", TypeLoss: "UJI", Tanggal: "", Nilai: "USD"},
	}
	if d.Tanda != "Reject Claim" || !d.TabelEstimasi || !reflect.DeepEqual(d.Estimasi, mau) {
		t.Fatalf("data reject: tanda %q tabel %v\n dapat %+v\n mau   %+v", d.Tanda, d.TabelEstimasi, d.Estimasi, mau)
	}
	if n := models.NamaBerkasTolak("UJI-K-1"); n != "Reject Claim UJI-K-1.PDF" {
		t.Fatalf("nama berkas %q", n)
	}
	if _, err := models.SusunDataTolak(kontrak.KlaimTreaty{Daftar: map[string][]map[string]string{
		"ClaimData.EstimationList": {{"EstimationValue": "bukan angka"}}}}, models.Kasus{}, a, time.Now()); err == nil {
		t.Fatal("nilai estimasi bukan angka harus galat")
	}
	b, err := models.PDFTutupKlaim(d, time.Date(2026, 10, 10, 3, 0, 0, 0, time.UTC))
	if err != nil || !bytes.HasPrefix(b, []byte("%PDF-")) {
		t.Fatalf("PDF reject: %v", err)
	}
}
