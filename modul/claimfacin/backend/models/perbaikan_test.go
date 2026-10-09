package models_test

// Uji perbaikan kelainan XML (prompt §5, OQ-CFI-03) yang berada di models: setiap uji menunjukkan beda dengan perilaku
// Pega lama yang tercatat di docs/PARITAS.md. Butir 1 (bind) dan 11 (tanggal penomoran) = uji repository; butir 4, 5,
// 8 = uji seam HTTP (handlers); butir 10 tidak berlaku (CheckTotalSpreadingPct_Act tak terjangkau, PARITAS).

import (
	"testing"

	"nusantarare/modul/claimfacin/backend/models"
)

// §5 butir 2: Pega `SetIndexObj_Act` 1 menyetel PlaStatus 1 saat pop-up dibuka; di sini hanya sesudah GeneratePLA.
func TestPlaStatusHanyaSesudahGeneratePLA(t *testing.T) {
	h := halamanDLA()
	h.SetelDaftar(models.DaftarObjek, []models.Baris{{"IsFacretro": "1"}})
	if err := models.BukaPLA(h, 1); err != nil {
		t.Fatal(err)
	}
	if ob, _ := models.Objek(h, 1); ob["PlaStatus"] != "" {
		t.Fatalf("PlaStatus %q sesudah pop-up dibuka", ob["PlaStatus"])
	}
	if err := models.SelesaiPLA(h, 1); err != nil {
		t.Fatal(err)
	}
	if ob, _ := models.Objek(h, 1); ob["PlaStatus"] != "1" {
		t.Fatalf("PlaStatus %q sesudah GeneratePLA", ob["PlaStatus"])
	}
}

// §5 butir 3: Pega `BackFromRegister` 8-11 mengisi medan kosong dengan nilai uji; di sini medan kosong tetap kosong.
func TestBackFromRegisterTanpaNilaiUji(t *testing.T) {
	h := halamanAdj(models.Baris{})
	k := konteksUji()
	k.Langkah = "Input Estimasi"
	models.BackFromRegister(k, h)
	for _, j := range []string{models.CD + "DateOfLoss", models.CD + "ReportDate", models.CD + "DateReceived",
		models.CD + "ReporterName", models.CD + "ReporterTelp", models.CD + "Currency"} {
		if v := h.Ambil(j); v != "" {
			t.Errorf("%s diisi %q", j, v)
		}
	}
	if h.Ambil(models.JalurCatatan) != models.CatatanBack {
		t.Fatalf("pyNote %q", h.Ambil(models.JalurCatatan))
	}
}

// §5 butir 9: Pega `InsertProgressClaim` cabang REGISTER memakai literal REGISTER tanpa kutip / "RESIGTER"; di sini
// posisi REGISTER.
func TestProgresPosisiRegister(t *testing.T) {
	h := models.HalamanBaru()
	h.Setel(models.JalurIsRegister, "1")
	pr, sub, ada := models.RencanaProgres(konteksUji(), h, "CLM-000001", models.ParamProgres{})
	if !ada || pr.Posisi != models.PosisiRegister || sub.Jenis != models.PosisiRegister || sub.Status != "Auto Create Register" {
		t.Fatalf("%+v %+v", pr, sub)
	}
}

// Nomor klaim / PLA / DLA: KODE_PRODUKSI NONLIFE + huruf + BusinessOldId + "." + MM.YYYY + "." + urut 5 digit.
func TestRakitNomor(t *testing.T) {
	kasus := map[string]string{
		models.RakitNomorKlaim("UJI-K", "77", "03.2026", 1):  "UJI-K77.03.2026.00001",
		models.RakitNomorKlaim("UJI-P", "77", "12.2026", 42): "UJI-P77.12.2026.00042",
	}
	for got, mau := range kasus {
		if got != mau {
			t.Errorf("%q, mau %q", got, mau)
		}
	}
}
