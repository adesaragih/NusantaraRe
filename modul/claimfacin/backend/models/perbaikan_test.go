package models_test

// Uji perbaikan kelainan XML (prompt §5, OQ-CFI-03) yang berada di models: setiap uji menunjukkan beda dengan perilaku
// Pega lama yang tercatat di docs/PARITAS.md. Butir 1 (bind) dan 11 (tanggal penomoran) = uji repository; butir 4, 5,
// 8 = uji seam HTTP (handlers); butir 10 tidak berlaku (CheckTotalSpreadingPct_Act tak terjangkau, PARITAS).

import (
	"testing"

	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/tiruan"
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

func TestInputKodePosLokasiTanpaPemisahMenggantung(t *testing.T) {
	// InputKodePos1_Act 5 `TempLocaion.CARI1 + ", " + RW + ...`: lokasi sementara kosong (Zip Code dipilih sebelum
	// Location diketik) menghasilkan ", KEL, ..." - bagian kosong dilewati (work owner 11-10-2026 "RAPIHKAN!").
	k := konteksUji()
	a := k.Acuan.(*tiruan.Acuan)
	a.WilayahPos["11111"] = []models.BarisWilayah{{RW: "UJI KEL", District: "UJI KEC", City: "UJI KOTA", Province: "UJI PROV"}}
	for _, c := range []struct{ sementara, mau string }{
		{"", "UJI KEL, UJI KEC, UJI KOTA, UJI PROV"},
		{"UJI JALAN 1", "UJI JALAN 1, UJI KEL, UJI KEC, UJI KOTA, UJI PROV"},
	} {
		h := models.HalamanBaru()
		h.Setel(models.CD+"PostalCode", "11111")
		h.Setel(models.JalurLokasiSementara, c.sementara)
		if err := models.InputKodePos(k, h); err != nil {
			t.Fatal(err)
		}
		if got := h.Ambil(models.CD + "Location"); got != c.mau {
			t.Errorf("lokasi sementara %q: Location %q, mau %q", c.sementara, got, c.mau)
		}
	}
}

func TestAlamatPelaporTanpaPemisahMenggantung(t *testing.T) {
	// GetCeding_act 6 / GetReportStatus_Act 13 `CARI10 + ", " + CARI11 + ", " + CARI12 + ", " + CARI13` apa adanya: klien
	// ceding tanpa alamat (DEV: CLIENT_ADDRESS kosong) membuat Reporter Address ", , , " (work owner 11-10-2026
	// "PERBAIKI CHOOSENYA INI!"). Bagian kosong dilewati; tanpa alamat sama sekali Reporter Address kosong.
	k := konteksUji()
	a := k.Acuan.(*tiruan.Acuan)
	a.Agen["UJI CEDING TANPA ALAMAT"] = "UJI-AG1"
	a.Agen["UJI CEDING SEBAGIAN"] = "UJI-AG2"
	a.Alamat["UJI-AG2"] = [4]string{"UJI JALAN 1", "", "", "UJI KOTA"}
	for _, c := range []struct{ ceding, mau string }{
		{"UJI CEDING TANPA ALAMAT", ""},
		{"UJI CEDING SEBAGIAN", "UJI JALAN 1, UJI KOTA"},
	} {
		h := models.HalamanBaru()
		h.Setel(models.CD+"ReportAddress", "UJI ALAMAT LAMA")
		if err := models.GetCeding(k, h, c.ceding); err != nil {
			t.Fatal(err)
		}
		if got := h.Ambil(models.CD + "ReportAddress"); got != c.mau {
			t.Errorf("Choose %q: Reporter Address %q, mau %q", c.ceding, got, c.mau)
		}
		if got := h.Ambil(models.CD + "InsuredRelationshipOthers"); got != c.ceding {
			t.Errorf("Choose %q: Specify %q", c.ceding, got)
		}
		// Reporter Status 1 (satu ceding di CedingCoList) memakai rangkaian yang sama
		h = models.HalamanBaru()
		h.Setel(models.CD+"InsuredRelationship", "1")
		h.SetelDaftar(models.DaftarCedingCo, []models.Baris{{"CedingCoName": c.ceding}})
		if err := models.GetReportStatus(k, h); err != nil {
			t.Fatal(err)
		}
		if got := h.Ambil(models.CD + "ReportAddress"); got != c.mau {
			t.Errorf("Reporter Status ceding %q: Reporter Address %q, mau %q", c.ceding, got, c.mau)
		}
	}
}
