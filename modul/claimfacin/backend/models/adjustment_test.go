package models_test

// Uji tabel perhitungan adjustment (prompt §6 butir 7): SetGrossAdjustment_act + SetNilaiResikoSendiri (deductible
// tipe 1 / 2 / 3), SetValueAdjusterFee (VAT), SetSalvageValue, SetSpreadingAjsutement_Act, CountSpreadingAdjustment_ACT,
// CountTotalEstimasi_Act (pesan 20 / 21.2 / 21.3), DisableSendComite_Act. Angka buatan `UJI`.

import (
	"context"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/tiruan"
)

const idr = "10026"

func konteksUji() *models.Konteks {
	a := tiruan.AcuanBaru()
	a.NamaMU[idr] = "IDR"
	a.Kurs[idr] = "1"
	a.Kurs["UJI-USD"] = "15000"
	a.JenisReas["10003"] = "UJI QS FAC"
	a.JenisReas["10015"] = "UJI FAC RETRO"
	return &models.Konteks{Ctx: context.Background(), Acuan: a, Pelaku: "UJI-TEKNIK",
		Sekarang: time.Date(2026, 3, 5, 10, 0, 0, 0, models.Jakarta)}
}

// halamanAdj - satu objek, satu item (TSI RNM 250 jt IDR), estimasi 10 jt tercetak CFS, spreading klaim QS 60% / FAC
// 40%, satu baris adjustment bermata uang IDR.
func halamanAdj(adj models.Baris) *models.Halaman {
	h := models.HalamanBaru()
	h.Setel(models.AwalanPolis+".PercentShare", "25")
	h.Setel(models.OQ+"SobName", "UJI SOB")
	h.SetelDaftar(models.DaftarObjek, []models.Baris{{"ObjectName": "UJI"}})
	h.SetelDaftar(models.DaftarItem(1), []models.Baris{{"TSIPerObject": "1000000000", "KursObjectItem": "1",
		"TSINusare": "250000000", "TotalEstimasi": "10000000", "TotalEstimationValueinIDR": "10000000"}})
	h.SetelDaftar(models.DaftarDiItem(1, 1, models.AnakEstimasi), []models.Baris{{"CurrencyID": idr, "Currency": "IDR",
		"KursValue": "1", "EstimationValue": "10000000", "PrintFaceClaim": "1"}})
	h.SetelDaftar(models.DaftarDiItem(1, 1, models.AnakSpreadKlaim), []models.Baris{
		{"TreatyType": "10003", "TreatyName": "UJI QS FAC", "SharePercentage": "60", "CurrencyID": idr, "ClaimSpreaded": "6000000"},
		{"TreatyType": "10015", "TreatyName": "UJI FAC RETRO", "SharePercentage": "40", "CurrencyID": idr, "ClaimSpreaded": "4000000"},
	})
	base := models.Baris{"CurrencyID": idr, "UploadLOD": idr, "CurrencyDol": "1", "EstimationValue": "10000000",
		"TotalEstimasiValue": "10000000", "PaymentType": models.BayarFinal}
	for k, v := range adj {
		base[k] = v
	}
	h.SetelDaftar(models.DaftarAdj(1, 1), []models.Baris{base})
	h.SetelDaftar(models.DaftarDiAdj(1, 1, 1, models.AnakAdjSpread), []models.Baris{
		{"TreatyType": "10003", "TreatyName": "UJI QS FAC", "SharePercentage": "60"},
		{"TreatyType": "10015", "TreatyName": "UJI FAC RETRO", "SharePercentage": "40"},
	})
	return h
}

func adjSatu(t *testing.T, h *models.Halaman) models.Baris {
	t.Helper()
	b, err := models.Adj(h, 1, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestGrossAdjustment(t *testing.T) {
	// SetGrossAdjustment_act 6 ber-pre=false: gerbang `.IndividualRiskType=="3"` mati, persen deductible SELALU dinolkan
	// sebelum SetNilaiResikoSendiri (tipe 1 / 2 lalu menghitung deductible 0) - ditiru apa adanya.
	kasus := []struct {
		nama                              string
		adj                               models.Baris
		propose, adjRNM, valueIDR, persen string
	}{
		{"tanpa deductible", models.Baris{"GrossAdjustment": "20000000"}, "20000000", "5000000", "5000000", "0"},
		{"tipe 1 persen dinolkan langkah 6", models.Baris{"GrossAdjustment": "20000000", "IndividualRiskType": "1",
			"IndividualRiskPercentage": "10"}, "20000000", "5000000", "5000000", "0"},
		{"tipe 3 nilai tetap", models.Baris{"GrossAdjustment": "20000000", "IndividualRiskType": "3",
			"IndividualRiskValue": "3000000"}, "17000000", "4250000", "4250000", "0"},
	}
	for _, c := range kasus {
		t.Run(c.nama, func(t *testing.T) {
			h := halamanAdj(c.adj)
			if err := models.SetGrossAdjustment(konteksUji(), h, 1, 1, 1); err != nil {
				t.Fatal(err)
			}
			b := adjSatu(t, h)
			if b["PersenRNM"] != "25" || b["ProposeAdjustmentValue"] != c.propose || b["AdjustmentValue"] != c.adjRNM ||
				b["ValueAdjustment"] != c.valueIDR || b["IndividualRiskPercentage"] != c.persen || b["GrossValue"] != "5000000" {
				t.Fatalf("%v", b)
			}
			if h.AdaPesan() {
				t.Fatalf("pesan %v", h.SemuaPesan())
			}
		})
	}
}

func TestDeductibleLangsung(t *testing.T) {
	kasus := []struct {
		nama                           string
		adj                            models.Baris
		nilai, propose, adjRNM, dedRNM string
	}{
		{"tipe 1", models.Baris{"GrossAdjustment": "20000000", "IndividualRiskType": "1", "IndividualRiskPercentage": "10",
			"PersenRNM": "25"}, "2000000", "18000000", "4500000", "500000"},
		{"tipe 2", models.Baris{"GrossAdjustment": "20000000", "IndividualRiskType": "2", "IndividualRiskPercentage": "1",
			"PersenRNM": "25"}, "10000000", "10000000", "2500000", "2500000"},
		{"tipe 3", models.Baris{"GrossAdjustment": "20000000", "IndividualRiskType": "3", "IndividualRiskValue": "3000000",
			"IndividualRiskPercentage": "7", "PersenRNM": "25"}, "3000000", "17000000", "4250000", "750000"},
	}
	for _, c := range kasus {
		t.Run(c.nama, func(t *testing.T) {
			h := halamanAdj(c.adj)
			if err := models.SetNilaiResikoSendiri(konteksUji(), h, 1, 1, 1); err != nil {
				t.Fatal(err)
			}
			b := adjSatu(t, h)
			if b["IndividualRiskValue"] != c.nilai || b["ProposeAdjustmentValue"] != c.propose ||
				b["AdjustmentValue"] != c.adjRNM || b["IndividualRiskRNM"] != c.dedRNM {
				t.Fatalf("%v", b)
			}
			if c.adj["IndividualRiskType"] == "3" && b["IndividualRiskPercentage"] != "0" {
				t.Fatalf("tipe 3 persen %q (langkah 11 menolkan)", b["IndividualRiskPercentage"])
			}
			it, _ := models.Item(h, 1, 1)
			if it["AdjustmentVal"] != "1" || it[models.PropIndeksAdj] != "1" {
				t.Fatalf("penanda item %v", it)
			}
		})
	}
}

func TestAdjustmentMelebihiEstimasi(t *testing.T) {
	h := halamanAdj(models.Baris{"GrossAdjustment": "80000000", "PersenRNM": "25"})
	if err := models.SetGrossAdjustment(konteksUji(), h, 1, 1, 1); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(h.SemuaPesan(), "|"), models.PesanAdjLebihEstimasi) ||
		h.Ambil(models.JalurIsError) != "2" {
		t.Fatalf("pesan %v IsError %q", h.SemuaPesan(), h.Ambil(models.JalurIsError))
	}
	it, _ := models.Item(h, 1, 1)
	if it["AdjustmentVal"] != "0" {
		t.Fatalf("tombol komite tetap aktif %v", it)
	}
}

func TestFeeDenganVAT(t *testing.T) {
	h := halamanAdj(models.Baris{"PaymentType": models.BayarFee, "GrossAdjustment": "1000000", "VAT": "11",
		"PersenRNM": "25"})
	if err := models.SetValueAdjusterFee(konteksUji(), h, 1, 1, 1); err != nil {
		t.Fatal(err)
	}
	b := adjSatu(t, h)
	// VAT 11% x 1 jt = 110 rb; fee 100% = 1,11 jt; RNM 25% = 277,5 rb
	if b["VATValue"] != "110000" || b["ProfessionalFee"] != "1110000" || b["AdjusterFeeValue"] != "277500" ||
		b["IndividualRiskRNM"] != "27500" {
		t.Fatalf("%v", b)
	}
	sp := h.AmbilDaftar(models.DaftarDiAdj(1, 1, 1, models.AnakAdjSpread))
	if sp[0]["ClaimSpreaded"] != "166500" || sp[1]["ClaimSpreaded"] != "111000" {
		t.Fatalf("spreading fee %v", sp)
	}
}

func TestSalvage(t *testing.T) {
	h := halamanAdj(models.Baris{"PaymentType": models.BayarSalvage, "GrossAdjustment": "-4000000", "PersenRNM": "25"})
	if err := models.SetSalvageValue(konteksUji(), h, 1, 1, 1); err != nil {
		t.Fatal(err)
	}
	b := adjSatu(t, h)
	if b["SalvageValue"] != "-1000000" || b["ValueAdjustment"] != "-1000000" || h.AdaPesan() {
		t.Fatalf("%v pesan %v", b, h.SemuaPesan())
	}
	h = halamanAdj(models.Baris{"PaymentType": models.BayarSalvage, "GrossAdjustment": "4000000", "PersenRNM": "25"})
	if err := models.SetSalvageValue(konteksUji(), h, 1, 1, 1); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(h.SemuaPesan(), "|"), models.PesanSalvageNegatif) {
		t.Fatalf("salvage positif tanpa pesan: %v", h.SemuaPesan())
	}
}

func TestSpreadingAdjustmentDanQS(t *testing.T) {
	h := halamanAdj(models.Baris{"AdjustmentValue": "5000000"})
	h.SetelDaftar(models.DaftarDiAdj(1, 1, 1, models.AnakAdjQS), []models.Baris{
		{"TreatyType": "UJI-A", "SharePercentage": "50"}, {"TreatyType": "UJI-B", "SharePercentage": "50"}})
	if err := models.SetSpreadingAdjustment(h, 1, 1, 1); err != nil {
		t.Fatal(err)
	}
	b := adjSatu(t, h)
	sp := h.AmbilDaftar(models.DaftarDiAdj(1, 1, 1, models.AnakAdjSpread))
	qs := h.AmbilDaftar(models.DaftarDiAdj(1, 1, 1, models.AnakAdjQS))
	// QS+FAC 60% x 5 jt = 3 jt dipecah QS 50/50
	if sp[0]["ClaimSpreaded"] != "3000000" || sp[1]["ClaimSpreaded"] != "2000000" || sp[1]["TotalSpread"] != "5000000" ||
		qs[0]["ClaimSpreaded"] != "1500000" || qs[1]["TotalSpread"] != "3000000" ||
		b["TotalSharePersen"] != "100" || b["TotalSpreadAdjustment"] != "5000000" || b["TotalSpreadBreakQs"] != "3000000" {
		t.Fatalf("spread %v qs %v adj %v", sp, qs, b)
	}
}

func TestSpreadingAdjustmentGabungJenisKembar(t *testing.T) {
	h := halamanAdj(models.Baris{"AdjustmentValue": "1000000"})
	h.SetelDaftar(models.DaftarDiAdj(1, 1, 1, models.AnakAdjSpread), []models.Baris{
		{"TreatyType": "10003", "SharePercentage": "30"}, {"TreatyType": "10003", "SharePercentage": "30"},
		{"TreatyType": "10015", "SharePercentage": "40"}})
	if err := models.SetSpreadingAdjustment(h, 1, 1, 1); err != nil {
		t.Fatal(err)
	}
	sp := h.AmbilDaftar(models.DaftarDiAdj(1, 1, 1, models.AnakAdjSpread))
	if len(sp) != 2 || sp[0]["SharePercentage"] != "60" || sp[0]["ClaimSpreaded"] != "600000" {
		t.Fatalf("%v", sp)
	}
}

func TestCountSpreadingAdjustmentExGratia(t *testing.T) {
	h := halamanAdj(models.Baris{"AdjustmentValue": "1000000", "ExGratia": "1"})
	h.SetelDaftar(models.DaftarDiAdj(1, 1, 1, models.AnakAdjSpread), []models.Baris{
		{"TreatyType": "10003", "SharePercentage": "70"}, {"TreatyType": "10015", "SharePercentage": "40"}})
	if err := models.CountSpreadingAdjustment(h, 1, 1, 1); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(h.SemuaPesan(), "|"), models.PesanSharePersenLebih) {
		t.Fatalf("total 110%% tanpa pesan %v", h.SemuaPesan())
	}
	h = halamanAdj(models.Baris{"AdjustmentValue": "1000000", "ExGratia": "1"})
	h.SetelDaftar(models.DaftarDiAdj(1, 1, 1, models.AnakAdjSpread), []models.Baris{
		{"TreatyType": "10003", "SharePercentage": "50"}})
	if err := models.CountSpreadingAdjustment(h, 1, 1, 1); err != nil {
		t.Fatal(err)
	}
	if b := adjSatu(t, h); b["TotalSpreadAdjustment"] != "500000" ||
		!strings.Contains(strings.Join(h.SemuaPesan(), "|"), models.PesanSharePersenKurang) {
		t.Fatalf("%v pesan %v", b, h.SemuaPesan())
	}
}

func TestCheckCurrencyEstimasiMataUang(t *testing.T) {
	h := halamanAdj(models.Baris{"UploadLOD": idr, "CurrencyID": "", "Currency": ""})
	h.SetelDaftar(models.DaftarDiAdj(1, 1, 1, models.DaftarMataUangAdj), []models.Baris{
		{"CurrencyID": idr, "Currency": "IDR", "KursValue": "1"}})
	if err := models.CheckCurrency(konteksUji(), h, 1, 1, 1); err != nil {
		t.Fatal(err)
	}
	b := adjSatu(t, h)
	if b["CurrencyID"] != idr || b["CurrencyDol"] != "1" || b["EstimationValue"] != "10000000" {
		t.Fatalf("%v", b)
	}
	if n := len(h.AmbilDaftar(models.DaftarDiAdj(1, 1, 1, models.AnakAdjSpread))); n != 2 {
		t.Fatalf("spreading adjustment dari spreading klaim %d", n)
	}
}

func TestTambahAdjustmentPemeriksaan(t *testing.T) {
	k := konteksUji()
	// 20: tanpa adjuster / consultant
	h := halamanAdj(models.Baris{})
	h.SetelDaftar(models.DaftarAdj(1, 1), nil)
	if err := models.TambahAdjustment(k, h, 1, 1); err != nil {
		t.Fatal(err)
	}
	if !h.AdaPesan() || len(h.AmbilDaftar(models.DaftarAdj(1, 1))) != 0 {
		t.Fatalf("tanpa adjuster: pesan %v baris %d", h.SemuaPesan(), len(h.AmbilDaftar(models.DaftarAdj(1, 1))))
	}
	// 21.3: adjustment terakhir final belum berkomite
	h = halamanAdj(models.Baris{"PaymentType": models.BayarFinal})
	h.Setel(models.CD+"ConsultantID", "UJI-A")
	h.Setel(models.CD+"AppointedADJID", "UJI-A")
	if err := models.TambahAdjustment(k, h, 1, 1); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(h.SemuaPesan(), "|"), models.PesanAdjSesudahFinal) {
		t.Fatalf("pesan %v", h.SemuaPesan())
	}
	// baris baru: estimasi item, payable SOB, retro -> DirectToKasir false
	h = halamanAdj(models.Baris{"PaymentType": models.BayarInterim, "IsKomite": "1"})
	h.Setel(models.CD+"ConsultantID", "UJI-A")
	h.Setel(models.CD+"AppointedADJID", "UJI-A")
	h.SetelDaftar(models.DaftarDiItem(1, 1, models.AnakSpreadPolis), []models.Baris{{"TreatyType": models.TreatyFacRetro}})
	if err := models.TambahAdjustment(k, h, 1, 1); err != nil {
		t.Fatal(err)
	}
	adj := h.AmbilDaftar(models.DaftarAdj(1, 1))
	b := adj[len(adj)-1]
	if len(adj) != 2 || b["EstimationValue"] != "10000000" || b["TotalEstimasiValue"] != "10000000" ||
		b["PayableTo"] != "UJI SOB" || b["IsFacRetro"] != "1" || b["DirectToKasir"] != "false" ||
		b["pxCreateOperator"] != "UJI-TEKNIK" {
		t.Fatalf("%v", b)
	}
	if sp := h.AmbilDaftar(models.DaftarDiAdj(1, 1, 2, models.AnakAdjSpread)); len(sp) != 2 || sp[0]["ClaimSpreaded"] != "" {
		t.Fatalf("spreading adjustment baru %v", sp)
	}
}

func TestHapusAdjustment(t *testing.T) {
	h := halamanAdj(models.Baris{})
	h.SetelDaftar(models.DaftarAdj(1, 1), []models.Baris{{"PaymentType": "1"}, {"PaymentType": "2"}})
	if err := models.HapusAdjustment(konteksUji(), h, 1, 1, 2); err != nil {
		t.Fatal(err)
	}
	it, _ := models.Item(h, 1, 1)
	if len(h.AmbilDaftar(models.DaftarAdj(1, 1))) != 1 || it[models.PropIndeksAdj] != "1" || it["IsKomite"] != "" ||
		it["AdjustmentVal"] != "0" {
		t.Fatalf("item %v", it)
	}
	kr := h.AmbilDaftar(models.DaftarKronologi)
	if len(kr) != 1 || kr[0]["pyNote"] != models.TeksBatalAdjustment+"2" {
		t.Fatalf("kronologi %v", kr)
	}
}

func TestDeductibleNolKeluarDiLangkah15(t *testing.T) {
	// SetNilaiResikoSendiri 15 (pesan "can not be filled by Zero") bertransisi pasca `true` -> 6 keluar: langkah 16+
	// tidak berjalan - pesan "lebih besar dari estimasi" (17) TIDAK ikut tampil walau jumlah adjustment melampaui.
	h := halamanAdj(models.Baris{"GrossAdjustment": "20000000", "IndividualRiskType": "1",
		"IndividualRiskPercentage": "100", "PersenRNM": "25"})
	rows := h.AmbilDaftar(models.DaftarAdj(1, 1))
	h.SetelDaftar(models.DaftarAdj(1, 1), append(rows, models.Baris{"CurrencyID": idr, "UploadLOD": idr,
		"PaymentType": models.BayarFinal, "ValueAdjustment": "50000000"}))
	if err := models.SetNilaiResikoSendiri(konteksUji(), h, 1, 1, 1); err != nil {
		t.Fatal(err)
	}
	pesan := strings.Join(h.SemuaPesan(), "|")
	if !strings.Contains(pesan, models.PesanAdjNol) || strings.Contains(pesan, models.PesanAdjLebihEstimasi) {
		t.Fatalf("pesan %q", pesan)
	}
}
