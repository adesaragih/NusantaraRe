package handlers_test

// Uji seam HTTP - polis NonProp baru: With Tax dicentang SESUDAH Choose Business menghitung ulang pajak grid XOL
// Current Premium (keputusan work owner 07-10-2026, `models.HitungUlangPajakNonPropEDM`). Fixture UJI-.

import (
	"encoding/json"
	"net/http"
	"testing"

	"nusantarare/modul/edmtreatyin/backend/models"
	"nusantarare/modul/edmtreatyin/backend/services"
)

const (
	nopolNonProp   = "UJI-POL-0002"
	kontrakNonProp = "UJI-KONTRAK-02"
)

// masterNonPropUji - master XOL satu layer IDR (gross 1000 / net 800 / potongan 200), satu angsuran satu termin.
func masterNonPropUji(id string) models.MasterXOL {
	mv := func(v string) []models.Baris { return []models.Baris{{"Currency": "IDR", "Value": v}} }
	return models.MasterXOL{Nilai: map[string]string{"ID": id, "Commencement": "2026-01-01", "Termination": "2026-12-31",
		"FacultativeShare": "0", "InstallmentNo": "UJI-NO"}, Daftar: map[string][]models.Baris{
		"Share":                          {{"Layer": "1", "LayerType": "UJI-LT"}},
		"Share(1).GrossPremiumList":      mv("1000"),
		"Share(1).NetPremiumList":        mv("800"),
		"Share(1).DeductionTotalList":    mv("200"),
		"Share(1).DeductionList":         {{"Currency": "IDR", "Deduction": "200"}},
		"Installment":                    {{"Currency": "IDR"}},
		"Installment(1).InstallmentList": {{"Installment": "UJI"}},
	}}
}

// polis NB NonProp baru TANPA pajak (With Tax tidak dicentang) - warisan WarisPajakLama = "false".
func baruNonProp(t *testing.T) *uji {
	u := baru(t)
	h := polisNBUji()
	h.Setel(models.HalamanPolis+".NoOffer", kontrakNonProp)
	h.Setel(models.HalamanPolis+".IsNewPolicyNonProp", "1")
	h.Setel(models.HalamanPolis+".FlagPPH", "false")
	h.Setel(models.HalamanQuotation+".ProportionalType", models.JenisNonProporsional)
	h.Setel(models.HalamanPolis+".QuotationData.ProportionalType", models.JenisNonProporsional)
	u.g.TanamPolis("NB-2", nopolNonProp, 0, "", h, kontrakNonProp)
	u.g.Master["ID:"+kontrakNonProp] = []models.MasterXOL{masterNonPropUji(kontrakNonProp)}
	return u
}

func lapisXOL(h *models.Halaman) models.Baris {
	b := h.AmbilDaftar(models.JalurAnak(models.DaftarXOL, 1, models.AnakLayerXOL))
	if len(b) == 0 {
		return nil
	}
	return b[0]
}

func TestWithTaxSesudahChooseBusinessMenghitungPajakXOL(t *testing.T) {
	u := baruNonProp(t)
	kode, isi := u.panggil("POST", "/kasus", admin, map[string]string{"noPolis": nopolNonProp, "edmType": "1"})
	u.wajib(kode, isi, http.StatusCreated, "buat NonProp")
	var k models.Kasus
	if err := json.Unmarshal([]byte(isi), &k); err != nil {
		t.Fatal(err)
	}
	ly := u.pilihBisnis(k.ID)
	if !models.PolisNonPropBaru(ly.Halaman) || lapisXOL(ly.Halaman) == nil {
		t.Fatalf("prasyarat: NonProp baru ber-grid XOL, IsNewPolicyNonProp %q", ly.Halaman.Ambil(models.HalamanPolis+".IsNewPolicyNonProp"))
	}
	if v := lapisXOL(ly.Halaman)["PPHValue"]; v != "" {
		t.Fatalf("prasyarat: tanpa With Tax layer belum berpajak, PPHValue %q", v)
	}

	// admin mencentang With Tax (Type Tax bawaan Inclusive) lalu Save
	isian := models.HalamanBaru()
	isian.Setel(models.HalamanPolis+".FlagPPH", "true")
	kode, isi = u.panggil("PUT", "/kasus/"+k.ID, admin, map[string]any{"halaman": isian})
	u.wajib(kode, isi, http.StatusOK, "simpan With Tax")
	var hasil services.Layar
	if err := json.Unmarshal([]byte(isi), &hasil); err != nil {
		t.Fatal(err)
	}
	h := hasil.Halaman
	if h.Ambil(models.HalamanPolis+".TypeTax") != models.TypeTaxInclusive {
		t.Fatalf("Type Tax bawaan %q", h.Ambil(models.HalamanPolis+".TypeTax"))
	}
	if l := lapisXOL(h); l["PPHValue"] == "" || l["NetPremiAfterTax"] == "" {
		t.Fatalf("Current Premium tetap tanpa pajak sesudah With Tax dicentang: %+v", l)
	}
	if b := h.AmbilDaftar(models.DaftarXOL)[0]; b["PPNValue"] == "" {
		t.Fatalf("induk Current Premium tanpa PPN: %+v", b)
	}
	// tersimpan: dibuka ulang tetap berpajak
	if l := lapisXOL(u.layar(k.ID, admin).Halaman); l["PPHValue"] == "" {
		t.Fatalf("pajak hilang sesudah dibuka ulang: %+v", l)
	}
}
