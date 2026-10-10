package models

import (
	"strings"
	"testing"
	"time"
)

// Perbaikan OQ-CNP-05 butir 1 - GenerateCFS_act: agregasi CFS mulai dari daftar kosong (XML: daftar salinan tidak
// dibuang -> jumlah ganda) dan Fee / Claim Amount Loss Allocation tidak tertukar.
func TestCFSAgregasiTanpaGandaDanFeeTidakTertukar(t *testing.T) {
	h := HalamanBaru()
	h.SetelDaftar(DaftarClaimAmount, []Baris{
		{"Currency": "IDR", "ClaimAmountCedant": "1000", "TPL": "0", "CNPDeductible": "0", "AdjusterFee": "10",
			"CNPOthersFee": "5", "Salvage": "1", "AltValue": "1"},
		{"Currency": "IDR", "ClaimAmountCedant": "500", "TPL": "0", "CNPDeductible": "0", "AdjusterFee": "20",
			"CNPOthersFee": "5", "Salvage": "1", "AltValue": "1"},
	})
	h.SetelDaftar(DaftarLossAlloc, []Baris{
		{"TreatyName": "OR", "Currency": "IDR", "ClaimAmountAdjust": "1200", "CNPOthersFee": "7", "AdjusterFee": "3", "Salvage": "1"},
		{"TreatyName": "OR", "Currency": "IDR", "ClaimAmountAdjust": "300", "CNPOthersFee": "3", "AdjusterFee": "2", "Salvage": "1"},
	})
	d, err := SusunDataCFS(h)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.KlaimMataUang) != 1 || len(d.LossAllocation) != 1 {
		t.Fatalf("CFS = %+v", d)
	}
	cekAngka(t, "Claim Amount Cedant IDR", d.KlaimMataUang[0]["ClaimAmountCedant"], "1500")
	cekAngka(t, "Adjuster Fee IDR", d.KlaimMataUang[0]["AdjusterFee"], "30")
	l := d.LossAllocation[0]
	cekAngka(t, "Loss Allocation Claim Amount", l["ClaimAmountAdjust"], "1500")
	cekAngka(t, "Loss Allocation Fee", l["CNPOthersFee"], "10")
}

// Perbaikan butir 2 - GetSelisihActual_Act membaca alias SQL GetDataOS (Value / GrossValue), bukan CARI1 / CARI2.
func TestSelisihAktualMembacaAliasGetDataOS(t *testing.T) {
	h := halamanDEV3998()
	if err := HitungKlaim(konteksUji(), h, masterDEV1001789()); err != nil {
		t.Fatal(err)
	}
	os := []NilaiOS{{Value: "930000000", GrossValue: "3100000000"}, {Value: "1350000000", GrossValue: "4500000000"}}
	s, err := HitungSelisihAktual(h, os)
	if err != nil {
		t.Fatal(err)
	}
	cekAngka(t, "EstimasiSpread (layer terakhir)", s.EstimasiSpread, "1350000000")
	cekAngka(t, "EstimasiTotal (layer terakhir)", s.EstimasiTotal, "4500000000")
	cekAngka(t, "SelisihSpread", s.SelisihSpread, "0")
	cekAngka(t, "SelisihTotal", s.SelisihTotal, "0")
}

// Perbaikan butir 3 - ProteksiSendKomiteCNP_Act: "Error No Account" aktif bila rekening Spreading In bermata uang
// rekening 1 bukan rekening 1 / 2 (XML: syarat saling meniadakan, tidak pernah aktif); IsError direset setiap jalan.
func TestProteksiKomiteRekeningSalahDanIsErrorDireset(t *testing.T) {
	h := HalamanBaru()
	h.Setel(CD+"PolicyData.PolicyNo", "UJI-POLIS")
	h.Setel(CD+"Payable", "1")
	h.Setel(CD+"Occupation", "UJI")
	h.Setel("IsError", "2")
	h.SetelDaftar(DaftarAdjustment, []Baris{{"PaymentType": "1", "Currency": "IDR", "CurrencyID": "10026",
		"NoAccount": "111", "NameOfBank": "UJI-BANK", "BranchOfBank": "UJI", "IDOfBank": "1"}})
	h.SetelDaftar(JalurAdj(1, AnakSpreadIn), []Baris{{"Currency": "IDR", "NoAccount": "999"}})
	c, err := ProteksiKirimKomite(konteksUji(), h, 1, "UJI-EMAIL")
	if err != nil {
		t.Fatal(err)
	}
	if !c.RekeningSalah {
		t.Error("rekening Spreading In bukan rekening akseptasi harus ditolak")
	}
	if h.Ambil("IsError") != "" {
		t.Errorf("IsError = %q, harus direset", h.Ambil("IsError"))
	}
	if c.Lolos(h) {
		t.Error("gerbang Send Claim to Committee harus tertutup")
	}
	h.SetelDaftar(JalurAdj(1, AnakSpreadIn), []Baris{{"Currency": "IDR", "NoAccount": "111"}})
	c, err = ProteksiKirimKomite(konteksUji(), h, 1, "UJI-EMAIL")
	if err != nil {
		t.Fatal(err)
	}
	if c.RekeningSalah {
		t.Error("rekening sama tidak boleh ditolak")
	}
}

// Perbaikan butir 4 - CreateChildKomiteCNP_Act: validasi Gross Value vs Spreading In sebelum kasus komite dibuat.
func TestValidasiGrossSpreadingIn(t *testing.T) {
	h := HalamanBaru()
	h.SetelDaftar(DaftarAdjustment, []Baris{{}})
	h.SetelDaftar(JalurAdj(1, AnakXOL), []Baris{{"TreatyName": "UR", "Currency": "IDR", "ClaimSpreaded": "0"},
		{"TreatyName": "XL 1ST LAYER", "Currency": "IDR", "ClaimSpreaded": "930"}})
	h.SetelDaftar(JalurAdj(1, AnakSpreadIn), []Baris{{"Currency": "IDR", "ClaimSpreaded": "930"}})
	if ok, err := ValidasiGrossSpreadingIn(h, 1); err != nil || !ok {
		t.Fatalf("sama -> lolos, dapat %v %v", ok, err)
	}
	h.AmbilDaftar(JalurAdj(1, AnakSpreadIn))[0]["ClaimSpreaded"] = "900"
	if ok, err := ValidasiGrossSpreadingIn(h, 1); err != nil || ok {
		t.Fatalf("beda -> tolak, dapat %v %v", ok, err)
	}
}

// Perbaikan butir 5 - HitServiceToKasir_Act: tahun Tgl Boleh Bayar ikut bergulir (XML: hanya bila bulan SEKARANG 12).
func TestTanggalBolehBayarTahunBergulir(t *testing.T) {
	for _, c := range []struct {
		aksep time.Time
		ingin string
	}{
		{time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC), "10-04-2026"},
		{time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC), "01-05-2026"},
		{time.Date(2026, 11, 26, 0, 0, 0, 0, time.UTC), "01-01-2027"},
		{time.Date(2026, 12, 10, 0, 0, 0, 0, time.UTC), "10-01-2027"},
		{time.Date(2026, 12, 27, 0, 0, 0, 0, time.UTC), "01-02-2027"},
	} {
		if got := TanggalBolehBayar(c.aksep); got != c.ingin {
			t.Errorf("TanggalBolehBayar(%s) = %s, ingin %s", c.aksep.Format("2006-01-02"), got, c.ingin)
		}
	}
}

// Perbaikan butir 6 - CloseClaimTNonProp: kronologi "Finish Adjustment (Close Claim)" tercatat (XML menulis CARI12).
func TestTutupKlaimMencatatKronologi(t *testing.T) {
	h := HalamanBaru()
	h.SetelDaftar(DaftarAdjustment, []Baris{{"AcceptanceStatus": "1"}})
	if !TutupKlaim(konteksUji(), h) {
		t.Fatal("tutup ditolak")
	}
	r := h.AmbilDaftar(DaftarRiwayat)
	if len(r) != 1 || r[0]["CommentSuggest"] != TeksTutupKlaim {
		t.Fatalf("riwayat = %v", r)
	}
	h2 := HalamanBaru()
	h2.SetelDaftar(DaftarAdjustment, []Baris{{"AcceptanceStatus": "0"}})
	if TutupKlaim(konteksUji(), h2) || !h2.AdaPesan() || len(h2.AmbilDaftar(DaftarRiwayat)) != 0 {
		t.Fatal("akseptasi di komite harus menolak tanpa kronologi")
	}
}

// Tangga komite OQ-CNP-01 (ikut XML): RNM Share <= 30 DAN ValueAdjustment akseptasi terakhir <= 30 jt -> tingkat 1 saja.
func TestTanggaKomiteHanyaTingkat1(t *testing.T) {
	for _, c := range []struct {
		share, nilai, subj string
		ingin              bool
	}{
		{"30", "30000000", "", true},
		{"30", "30000001", "", false},
		{"30.5", "1", "", false},
		{"50", "100000000", "true", true},
	} {
		h := HalamanBaru()
		h.Setel(TM+"RNMShare", c.share)
		h.SetelDaftar(DaftarAdjustment, []Baris{{"ValueAdjustment": c.nilai, "IsSubjectivity": c.subj}})
		if got := HanyaTingkat1(h, 1); got != c.ingin {
			t.Errorf("share %s nilai %s subj %q -> %v, ingin %v", c.share, c.nilai, c.subj, got, c.ingin)
		}
	}
}

func TestNomorKlaimDanPLA(t *testing.T) {
	if got := RakitNomorKlaim("RNM-K", "22", "04.2026", 4); got != "RNM-K22.04.2026.TX00004" {
		t.Errorf("nomor klaim = %s (DEV CLMNP-3998: RNM-K22.04.2026.TX00004)", got)
	}
	pla := RakitNomorPLA("22", time.Date(2026, 10, 9, 10, 0, 0, 0, Jakarta), 1)
	if pla != "RNM-M22.10.2026.TX00001" {
		t.Errorf("nomor PLA = %s", pla)
	}
	if got := RevisiNomorPLA(pla); got != pla+"/1" {
		t.Errorf("revisi pertama = %s", got)
	}
	if got := RevisiNomorPLA(pla + "/1"); got != pla+"/1/2" {
		t.Errorf("revisi kedua = %s", got)
	}
}

// DATA_JSON OS STS 0 (Save to issue RNM) = baris DEV CLMNP-3998 (OS_AKSEPTASI_KLAIM, baca-saja 09-10-2026) pada kunci
// yang sama: urutan kunci, pemisah LF, dan nilai. Kunci DEV EstimationDate / PolicyNo / pzInsKey tidak ditulis XML
// (SaveDataToOSAksep_Act 15.8.1: EstimationDate = "") - OQ, PARITAS.
func TestDataJSONOSSamaDenganBarisDEV(t *testing.T) {
	h := halamanDEV3998()
	h.Setel(CD+"NoClaim", "RNM-K22.04.2026.TX00004")
	h.Setel(CD+"IDMaster", "1001789")
	h.Setel(TM+"ID", "1001789")
	if err := HitungKlaim(konteksUji(), h, masterDEV1001789()); err != nil {
		t.Fatal(err)
	}
	layer := KelompokLayerOS(h)
	if len(layer) != 2 {
		t.Fatalf("layer OS = %d, ingin 2 (UR dibuang)", len(layer))
	}
	p := ParamOSLayer(h, layer[0])
	if err := KurangiOS(h, p, NilaiOS{}); err != nil {
		t.Fatal(err)
	}
	b := SusunBarisOS(h, "CLMNP-000001", StsOSOutstanding, p)
	dev := "{\n\"Adjusterfee\":\"0\"\n,\"CNPOthersFee\":\"0\"\n,\"Currency\":\"IDR\"\n,\"CurrencyID\":\"10026\"\n" +
		",\"EstimationDate\":\"20260504\"\n,\"GrossValue\":\"3100000000.000000000000000\"\n,\"IDMasterTreaty\":\"1001789\"\n" +
		",\"KursValue\":\"1.0\"\n,\"NoClaim\":\"RNM-K22.04.2026.TX00004\"\n,\"PersenRNM\":\"30\"\n" +
		",\"PolicyNo\":\"RNM-QR.T14.04.2025.11346\"\n,\"pxObjClass\":\"ASM-FW-GCNMFW-Data-osAkseptasi\"\n" +
		",\"pzInsKey\":\"ASM-FW-GCNMFW-WORK CLMNP-3998\"\n,\"Salvage\":\"0\"\n,\"Type\":\"0\"\n,\"TypeLoss\":\"XL 1ST LAYER\"\n" +
		",\"TypeLossID\":\"10046\"\n,\"Value\":\"930000000.000000000000000\"\n}\n"
	got, ingin := pasanganJSON(t, b.DataJSON), pasanganJSON(t, dev)
	var kunciSama []string
	for _, k := range ingin.urut {
		if _, ada := got.nilai[k]; ada {
			kunciSama = append(kunciSama, k)
		}
	}
	var urutGot []string
	for _, k := range got.urut {
		if _, ada := ingin.nilai[k]; ada {
			urutGot = append(urutGot, k)
		}
	}
	if strings.Join(urutGot, ",") != strings.Join(kunciSama, ",") {
		t.Errorf("urutan kunci = %v, DEV %v", urutGot, kunciSama)
	}
	for _, k := range kunciSama {
		if !SamaAngka(got.nilai[k], ingin.nilai[k]) {
			t.Errorf("%s = %q, DEV %q", k, got.nilai[k], ingin.nilai[k])
		}
	}
	for k := range got.nilai {
		if _, ada := ingin.nilai[k]; !ada {
			t.Errorf("kunci %s tidak ada di baris DEV", k)
		}
	}
	if !strings.HasPrefix(b.DataJSON, "{\n") || !strings.HasSuffix(b.DataJSON, "\n}\n") {
		t.Errorf("bentuk GetPageJSONString salah: %q", b.DataJSON)
	}
	if b.CaseID != "CLMNP-000001" || b.StsReject != "0" || b.MasterID != "1001789" {
		t.Errorf("kolom OS = %+v", b)
	}
}

type jsonDatar struct {
	urut  []string
	nilai map[string]string
}

// pasanganJSON - pengurai JSON datar bentuk GetPageJSONString (satu pasangan per baris).
func pasanganJSON(t *testing.T, s string) jsonDatar {
	t.Helper()
	out := jsonDatar{nilai: map[string]string{}}
	for _, baris := range strings.Split(s, "\n") {
		baris = strings.TrimPrefix(strings.TrimSpace(baris), ",")
		if !strings.HasPrefix(baris, "\"") {
			continue
		}
		i := strings.Index(baris, "\":\"")
		if i < 0 {
			t.Fatalf("baris JSON tidak dikenal: %q", baris)
		}
		k, v := baris[1:i], strings.TrimSuffix(baris[i+3:], "\"")
		out.urut = append(out.urut, k)
		out.nilai[k] = v
	}
	return out
}
