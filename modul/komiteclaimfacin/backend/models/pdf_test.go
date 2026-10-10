package models_test

// Uji PDF dokumen akseptasi (stream `AcceptanceNotePDF`, PrintPDFAccep_MultiAksep_KMT): format angka Java
// `DecimalFormat("#,###.####")` id_ID dan penggambar fpdf (teks dibaca dari stream konten PDF yang dimampatkan).

import (
	"bytes"
	"compress/zlib"
	"io"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimfacin/backend/models"
	"nusantarare/modul/komiteclaimfacin/backend/tiruan"
)

// DecimalFormat("#,###.####") locale id_ID - nilai harapan dihitung tangan dari pola Java (HALF_EVEN atas nilai
// DESIMAL, tanpa '0'; Java atas double berbeda di seri tepat ...5, lihat AngkaPega).
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

func TestPDFDokumenMenggambarSusunanStream(t *testing.T) {
	saat := time.Date(2026, 10, 10, 9, 0, 0, 0, models.Jakarta)
	d := models.DokumenPDF{
		Judul: []string{"ACCEPTED CLAIM INSURANCE"}, Nomor: "UJI-A12.10.2026.00001",
		Baris: []models.BarisPDF{
			{Label: "Policy No", Nilai: []string{"UJI-POL-1"}},
			{Label: "TSI RNM", Nilai: []string{"UJI ITEM 1 IDR 250.000.000", "UJI ITEM 2 USD 2.500,5"}},
			{Label: "Payable To", Nilai: []string{"UJI SOB", "Name of Bank : UJI BANK"}},
			{Label: "Spreading Adjustment", Tabel: &models.TabelPDF{Kepala: []string{"Treaty Name", "Share (%)", "Claim Spread"},
				Baris: [][]string{{"UJI QS", "60", "3.000.000"}}}},
		},
		Penutup: "Jakarta, 10 October 2026",
	}
	b, err := models.PDFDokumen(d, saat)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(b, []byte("%PDF-")) {
		t.Fatal("bukan PDF")
	}
	h := teksPDF(t, b)
	for _, s := range []string{"(" + models.KepalaPDFAkseptasi + ")", "(ACCEPTED CLAIM INSURANCE)",
		"(Accepted No : UJI-A12.10.2026.00001)", "(Policy No)", "(UJI-POL-1)", "(UJI ITEM 2 USD 2.500,5)",
		"(Name of Bank : UJI BANK)", `(Share \(%\))`, "(3.000.000)", "(Jakarta, 10 October 2026)", "(Page 1 of 1 pages)"} {
		if !strings.Contains(h, s) {
			t.Errorf("PDF tanpa %q", s)
		}
	}
	if ulang, _ := models.PDFDokumen(d, saat); !bytes.Equal(b, ulang) {
		t.Fatal("PDF tidak deterministik")
	}
}

// awalStream - penanda awal isi stream PDF yang ditulis fpdf.
const awalStream = "stream\n"

// teksPDF - isi semua stream FlateDecode PDF (konten halaman) sebagai teks.
func teksPDF(t *testing.T, b []byte) string {
	t.Helper()
	var out strings.Builder
	for {
		i := bytes.Index(b, []byte(awalStream))
		if i < 0 {
			return out.String()
		}
		b = b[i+len(awalStream):]
		j := bytes.Index(b, []byte("endstream"))
		if j < 0 {
			t.Fatal("stream tanpa endstream")
		}
		if r, err := zlib.NewReader(bytes.NewReader(b[:j])); err == nil {
			isi, _ := io.ReadAll(r)
			out.Write(isi)
		}
		b = b[j+len("endstream"):]
	}
}

func TestSusunAcceptanceNoteFire(t *testing.T) {
	// PrintPDFAccep_MultiAksep_KMT S1-S15 + stream AcceptanceNotePDF atas klaim Fire uji: item 1 (adjustment ber-nomor
	// lama) disaring, item 2 = item terakhir yang lolos; Location of Loss = CauseOfLoss (S14 VERBATIM).
	n, d := tiruan.KlaimFacInUji("20000000")
	kl := kontrak.KlaimFacIn{Nilai: n, Daftar: d, Objek: 1, Item: 2, Adjustment: 1}
	adj := d[models.DaftarAdj(1, 2)][0]
	adj["AcceptanceStatus"], adj["AcceptedNo"], adj["IsPrintAccept"] = "1", "UJI-A12.10.2026.00001", "1"
	saat := time.Date(2026, 10, 10, 9, 0, 0, 0, models.Jakarta)
	doc, err := models.SusunAcceptanceNote(kl, "UJI-A12.10.2026.00001", saat, "UJI Kepala",
		map[string]string{"10003": "UJI QS", "10015": "UJI FAC RETRO"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(doc.Judul, "|") != "ACCEPTED CLAIM INSURANCE" || doc.Nomor != "UJI-A12.10.2026.00001" ||
		doc.Penutup != "Jakarta, 10 October 2026" {
		t.Fatalf("kepala dokumen: %+v", doc)
	}
	var got []string
	for _, b := range doc.Baris {
		s := b.Label + "=" + strings.Join(b.Nilai, "|")
		if b.Tabel != nil {
			for _, r := range b.Tabel.Baris {
				s += "[" + strings.Join(r, ";") + "]"
			}
		}
		got = append(got, s)
	}
	mau := []string{
		"Policy No=UJI-RNM-F.001", "Line of Business=UJI FIRE", "Claim No / Claim ID=UJI-K-0001/" + tiruan.KlaimUji,
		"Name of Insured=UJI TERTANGGUNG", "SOB Name=UJI SOB", "Ceding Co Name=UJI CEDING",
		"Period of Policy=01/01/2026- 31/12/2026", "Object Item Name=UJI GEDUNG  - UJI ITEM 2", "TSI RNM=UJI ITEM 2 IDR 0",
		"Location of Loss=UJI KEBAKARAN", "Date of Loss=30/09/2026", "Cause of Loss=UJI KEBAKARAN", "Premium Paid On=",
		"PIC Name / Technical PIC=UJI Admin/UJI Kepala", "Accepted Claim=IDR 20.000.000",
		"Payable To=UJI PENERIMA|Name of Bank : UJI BANK|Branch of Bank : UJI CABANG|Account No : 12-34",
		"Spreading Adjustment=[UJI QS;60;12.000.000][UJI FAC RETRO;40;8.000.000]",
	}
	if strings.Join(got, "\n") != strings.Join(mau, "\n") {
		t.Fatalf("baris dokumen:\n%s", strings.Join(got, "\n"))
	}
	if s := models.StreamAkseptasiLini(n); s != models.StreamAkseptasi {
		t.Fatalf("lini Fire memakai %q", s)
	}
	if s := models.StreamAkseptasiLini(map[string]string{models.OQ + "BusinessType": "MBUCar"}); s != "AcceptanceNotePDFMBU" {
		t.Fatalf("lini MBU memakai %q", s)
	}
}

func TestSusunAcceptanceNoteMataUangKumulatif(t *testing.T) {
	// S10.3.1-S10.3.9: TemporaryUang bertambah SELAMA iterasi - adjustment ke-2 (IDR, U=1) menjumlah, ke-3 (USD, U=2)
	// mengganti; PT 4 (Adjuster Fee) memakai AdjusterFeeValue. Status bukan 1 tidak dicetak.
	n, d := tiruan.KlaimFacInUji("1000")
	kl := kontrak.KlaimFacIn{Nilai: n, Daftar: d, Objek: 1, Item: 2, Adjustment: 1}
	no := "UJI-A12.10.2026.00002"
	a := d[models.DaftarAdj(1, 2)][0]
	a["AcceptanceStatus"], a["AcceptedNo"] = "1", no
	salin := func(ubah map[string]string) {
		b := map[string]string{}
		for k, v := range a {
			b[k] = v
		}
		for k, v := range ubah {
			b[k] = v
		}
		d[models.DaftarAdj(1, 2)] = append(d[models.DaftarAdj(1, 2)], b)
	}
	salin(map[string]string{"AdjustmentValue": "2000"})
	salin(map[string]string{"AcceptanceStatus": "2", "AdjustmentValue": "999"})
	salin(map[string]string{"CurrencyID": "UJI-USD", "Currency": "USD", "PaymentType": "4", "AdjusterFeeValue": "500"})
	doc, err := models.SusunAcceptanceNote(kl, no, time.Date(2026, 10, 10, 9, 0, 0, 0, models.Jakarta), "UJI Kepala",
		nil)
	if err != nil {
		t.Fatal(err)
	}
	var nilai []string
	for _, b := range doc.Baris {
		if strings.HasPrefix(b.Label, "Accepted ") {
			nilai = append(nilai, b.Label+"="+strings.Join(b.Nilai, "|"))
		}
	}
	if got := strings.Join(nilai, ";"); got != "Accepted Claim=IDR 1.000;Accepted Claim=IDR 3.000;Accepted Adjuster Fee=USD 500" {
		t.Fatalf("nilai akseptasi: %s", got)
	}
}

func TestPDFDokumenHurufCp1252DanNilaiPanjang(t *testing.T) {
	// temuan review 10-10-2026: fpdf.SplitText panik pada teks cp1252 bukan ASCII; nilai lebih panjang dari satu halaman
	// dahulu menjadi satu baris per halaman.
	var tsi []string
	for i := 0; i < 60; i++ {
		tsi = append(tsi, "UJI ITEM IDR 1.000.000")
	}
	d := models.DokumenPDF{Judul: []string{"ACCEPTED CLAIM INSURANCE"}, Nomor: "UJI-A1", Baris: []models.BarisPDF{
		{Label: "Cause of Loss", Nilai: []string{"Kebakaran – gudang café ’UJI’ " + strings.Repeat("sangat panjang ", 20)}},
		{Label: "TSI RNM", Nilai: tsi},
		{Label: "Spreading Adjustment", Tabel: &models.TabelPDF{Kepala: []string{"Treaty Name", "Share (%)", "Claim Spread"},
			Baris: [][]string{{"UJI Réas – QS", "60", "1"}}}},
	}, Penutup: "Jakarta, 10 October 2026"}
	b, err := models.PDFDokumen(d, time.Date(2026, 10, 10, 9, 0, 0, 0, models.Jakarta))
	if err != nil {
		t.Fatal(err)
	}
	h := teksPDF(t, b)
	if !strings.Contains(h, "(Kebakaran \x96 gudang caf\xe9 \x92UJI\x92") || !strings.Contains(h, "(UJI R\xe9as \x96 QS)") {
		t.Fatal("huruf cp1252 tidak tercetak")
	}
	if n := strings.Count(h, "(Page "); n != 3 || !strings.Contains(h, "(Page 3 of 3 pages)") {
		t.Fatalf("60 baris TSI mau 3 halaman, dapat %d", n)
	}
}

func TestSusunAcceptanceNoteTanggalPegaGMT(t *testing.T) {
	// teks DateTime Pega ("yyyyMMddTHHmmss.SSS GMT") dicetak menurut Asia/Jakarta (@FormatDateTime ... "Asia/Jakarta").
	n, d := tiruan.KlaimFacInUji("1000")
	n["ClaimData.DateOfLoss"] = "20260930T170000.000 GMT"
	a := d[models.DaftarAdj(1, 2)][0]
	a["AcceptanceStatus"], a["AcceptedNo"] = "1", "UJI-A3"
	doc, err := models.SusunAcceptanceNote(kontrak.KlaimFacIn{Nilai: n, Daftar: d, Objek: 1}, "UJI-A3",
		time.Date(2026, 10, 10, 9, 0, 0, 0, models.Jakarta), "UJI", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range doc.Baris {
		if b.Label == "Date of Loss" && strings.Join(b.Nilai, "") != "01/10/2026" {
			t.Fatalf("Date of Loss %v", b.Nilai)
		}
	}
}
