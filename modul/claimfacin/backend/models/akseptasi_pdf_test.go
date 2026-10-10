package models_test

// Uji perakit dokumen akseptasi `PrintPDFAccep_MultiAksep` (tombol Acceptation): saringan S9.3.11 VERBATIM
// (IsPrintAccept kosong + status 1 + nomor sama), TemporaryUang bertambah selama iterasi (S9.3.1-S9.3.9), format angka
// DecimalFormat id_ID, ClaimSpreaded 4 desimal.

import (
	"strings"
	"testing"
	"time"

	"nusantarare/modul/claimfacin/backend/models"
)

func TestSusunAcceptanceNoteSaringDanMataUangKumulatif(t *testing.T) {
	no := "UJI-AKS.03.2026.00009"
	h := models.HalamanBaru()
	h.Setel(models.OQ+"BusinessType", "Fire")
	h.Setel(models.JalurNoPolis, "UJI-RNM-F.001")
	h.Setel(models.CD+"NoClaim", "UJI-K-0009")
	h.Setel(models.CD+"DateOfLoss", "20260228T170000.000 GMT") // DateTime Pega -> 01/03/2026 Jakarta
	h.SetelDaftar(models.DaftarObjek, []models.Baris{{"ObjectName": "UJI GEDUNG"}})
	h.SetelDaftar(models.DaftarItem(1), []models.Baris{{"ObjectItemName": "UJI ITEM 1"}, {"ObjectItemName": "UJI ITEM 2"}})
	adj := func(cur, pt, nilai, cetak string) models.Baris {
		return models.Baris{"AcceptanceStatus": "1", "AcceptedNo": no, "IsPrintAccept": cetak, "CurrencyID": "UJI-" + cur,
			"Currency": cur, "PaymentType": pt, "AdjustmentValue": nilai, "AdjusterFeeValue": nilai, "SalvageValue": nilai,
			"pxCreateOpName": "UJI Admin"}
	}
	// item 1: satu adjustment sudah tercetak (dibuang) -> item 1 tanpa adjustment lolos, tidak tercetak
	h.SetelDaftar(models.DaftarAdj(1, 1), []models.Baris{adj("IDR", "1", "700", "1")})
	h.SetelDaftar(models.DaftarAdj(1, 2), []models.Baris{adj("IDR", "1", "1000", ""), adj("IDR", "3", "2000", ""),
		adj("USD", "6", "500", "")})
	h.SetelDaftar(models.DaftarDiItem(1, 2, models.AnakSpreadKlaim), []models.Baris{
		{"TreatyType": "10003", "SharePercentage": "33.33333"}})
	d, err := models.SusunAcceptanceNote(h, 1, no, "UJI-CLM-9", time.Date(2026, 3, 5, 10, 0, 0, 0, models.Jakarta),
		"UJI Teknik", map[string]string{"10003": "UJI QS"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, b := range d.Baris {
		switch {
		case b.Label == "Object Item Name" || b.Label == "Claim No / Claim ID" || b.Label == "Date of Loss" ||
			strings.HasPrefix(b.Label, "Accepted "):
			got = append(got, b.Label+"="+strings.Join(b.Nilai, "|"))
		case b.Tabel != nil:
			got = append(got, b.Label+"="+strings.Join(b.Tabel.Baris[0], ";"))
		}
	}
	mau := []string{"Claim No / Claim ID=UJI-K-0009/UJI-CLM-9", "Object Item Name=UJI GEDUNG  - UJI ITEM 2",
		"Date of Loss=01/03/2026",
		"Accepted Claim=IDR 1.000", "Spreading Adjustment=UJI QS;33,3333;166,6667",
		"Accepted Salvage=IDR 3.000", "Spreading Adjustment=UJI QS;33,3333;166,6667",
		"Accepted Consultant Fee=USD 500", "Spreading Adjustment=UJI QS;33,3333;166,6667"}
	if strings.Join(got, "\n") != strings.Join(mau, "\n") {
		t.Fatalf("dokumen:\n%s", strings.Join(got, "\n"))
	}
	if _, err := models.SusunAcceptanceNote(h, 1, "UJI-LAIN", "UJI-CLM-9", time.Now(), "", nil); err == nil {
		t.Fatal("nomor tanpa adjustment lolos harus galat")
	}
	if s := models.StreamAkseptasiLini(h); s != models.StreamAkseptasi {
		t.Fatalf("stream Fire %q", s)
	}
}
