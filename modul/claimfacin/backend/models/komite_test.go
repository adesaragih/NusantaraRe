package models_test

// Uji tabel batas komite (SetListKomite_act, SendPICProtect_Act, CreateKMTNo_Act) dan penutupan / penolakan
// (ValidationAdjustmentKomite, CloseClaim, SendRejectClaimToKomite2, SendCloseClaimToKomite). Angka `UJI`.

import (
	"strings"
	"testing"

	"nusantarare/modul/claimfacin/backend/models"
)

var rosterUji = []models.AnggotaKomite{
	{OperatorID: "UJI-K1", Jabatan: "UJI SPV", LimitBottom: "-9999999999999", LimitTop: "57750000"},
	{OperatorID: "UJI-K2", Jabatan: "UJI HEAD", LimitBottom: "57750000", LimitTop: "189750000"},
	{OperatorID: "UJI-K3", Jabatan: models.JabatanBatasFacOut, LimitBottom: "189750000", LimitTop: "495000000"},
	{OperatorID: "UJI-K4", Jabatan: "UJI TANPA BATAS"},
}

func TestRosterKomiteCalon(t *testing.T) {
	kasus := []struct {
		nama, retro, nilai string
		mau                []string
	}{
		{"non retro: batas 1", "", "900000000", []string{"UJI-K1"}},
		{"retro di bawah batas kedua", "1", "50000000", []string{"UJI-K1"}},
		{"retro di atas batas kedua", "1", "200000000", []string{"UJI-K1", "UJI-K2", "UJI-K3"}},
		// melampaui LIMIT_TOP Technic Div. Head: saringan batas diabaikan (LIMIT_BOTTOM NULL ikut)
		{"retro melampaui LIMIT_TOP", "1", "600000000", []string{"UJI-K1", "UJI-K2", "UJI-K3", "UJI-K4"}},
	}
	for _, c := range kasus {
		t.Run(c.nama, func(t *testing.T) {
			h := halamanAdj(models.Baris{"ValueAdjustment": c.nilai})
			h.SetelDaftar(models.DaftarObjek, []models.Baris{{"IsFacretro": c.retro}})
			calon, err := models.RosterKomiteCalon(h, 1, 1, 1, rosterUji)
			if err != nil {
				t.Fatal(err)
			}
			var got []string
			for _, a := range calon {
				got = append(got, a.OperatorID)
			}
			if strings.Join(got, ",") != strings.Join(c.mau, ",") {
				t.Fatalf("calon %v, mau %v", got, c.mau)
			}
		})
	}
}

func TestProteksiKomite(t *testing.T) {
	bank := models.Baris{"NameOfBank": "UJI", "NoAccount": "1", "IDOfBank": "UJI-B", "ValueAdjustment": "1000"}
	kasus := []struct {
		nama     string
		pt       string
		lampiran []string
		limit    string
		lolos    bool
	}{
		{"final lengkap", models.BayarFinal, []string{"LOD", "DLA", "SPGR"}, "", true},
		{"final tanpa SPGR", models.BayarFinal, []string{"LOD", "DLA"}, "", false},
		{"final lewat batas dirut tanpa ADU", models.BayarFinal, []string{"LOD", "DLA", "SPGR"}, "500", false},
		{"final lewat batas dirut dengan ADU", models.BayarFinal, []string{"LOD", "DLA", "SPGR", "ADU"}, "500", true},
		{"salvage", models.BayarSalvage, []string{"Salvage", "DLA"}, "", true},
		{"fee tanpa invoice", models.BayarFee, []string{"DLA"}, "", false},
		{"adjustment PT 5", models.BayarAdjust, nil, "", true},
	}
	for _, c := range kasus {
		t.Run(c.nama, func(t *testing.T) {
			adj := models.Baris{"PaymentType": c.pt}
			for k, v := range bank {
				adj[k] = v
			}
			h := halamanAdj(adj)
			lamp, cacah := map[string]bool{}, map[string]int{}
			for _, k := range c.lampiran {
				lamp[k], cacah[k] = true, 1
			}
			lolos, err := models.ProteksiKomite(h, 1, 1, 1, cacah, c.limit)
			if err != nil {
				t.Fatal(err)
			}
			if lolos != c.lolos {
				t.Fatalf("lolos %v pesan %v", lolos, h.SemuaPesan())
			}
			adaADU := strings.Contains(strings.Join(h.SemuaPesan(), "|"), "Approval Direktur Utama")
			if c.limit != "" && !lamp["ADU"] && !adaADU {
				t.Fatalf("pesan ADU tidak muncul: %v", h.SemuaPesan())
			}
		})
	}
	// tanpa rekening: pesan dan gagal walau lampiran lengkap (18)
	h := halamanAdj(models.Baris{"PaymentType": models.BayarAdjust})
	if lolos, _ := models.ProteksiKomite(h, 1, 1, 1, nil, ""); lolos ||
		!strings.Contains(strings.Join(h.SemuaPesan(), "|"), models.PesanBankKosong) {
		t.Fatalf("tanpa rekening lolos=%v pesan %v", lolos, h.SemuaPesan())
	}
}

// Perbaikan prompt §5 butir 6 dan 7: penyerahan TIDAK memotong tanggal tersimpan dan TIDAK menimpa circumstances.
func TestTandaiKirimKomiteTanpaMengubahTanggalDanKronologi(t *testing.T) {
	h := halamanAdj(models.Baris{"DataCommitteFacin.CircumCauseOfLoss": "UJI KRONOLOGI"})
	h.Setel(models.CD+"DateOfLoss", "2026-03-01")
	if err := models.TandaiKirimKomite(konteksUji(), h, 1, 1, 1, "KMT-UJI1"); err != nil {
		t.Fatal(err)
	}
	b := adjSatu(t, h)
	it, _ := models.Item(h, 1, 1)
	if h.Ambil(models.CD+"DateOfLoss") != "2026-03-01" || b["DataCommitteFacin.CircumCauseOfLoss"] != "UJI KRONOLOGI" ||
		b[models.PropKomiteID] != "KMT-UJI1" || b["IsKomite"] != "1" || it["IsKomite"] != "0" {
		t.Fatalf("adj %v item %v DOL %q", b, it, h.Ambil(models.CD+"DateOfLoss"))
	}
	if kr := h.AmbilDaftar(models.DaftarKronologi); kr[len(kr)-1]["pyNote"] != models.AwalanKirimKomite+"KMT-UJI1" {
		t.Fatalf("kronologi %v", kr)
	}
	if models.PesanPembayaran(models.BayarInterim) != "Payment - Claim" || models.PesanPembayaran("6") != "Payment - Consultant Fee" {
		t.Fatal("teks pembayaran CreateKMTNo_Act 15")
	}
}

func TestValidasiTutup(t *testing.T) {
	kasus := []struct {
		nama     string
		adj      []models.Baris
		lampiran int
		estimasi string
		pesan    []string
	}{
		{"komite menunggu", []models.Baris{{"AcceptanceStatus": "0"}}, 0, "1", []string{models.PesanTutupDiKomite}},
		{"belum cetak akseptasi", []models.Baris{{"AcceptanceStatus": "1"}}, 0, "1", []string{models.PesanBelumCetakAksep}},
		{"kasir belum masuk", []models.Baris{{"AcceptanceStatus": "1", "IsPrintAccept": "1", "DirectToKasir": "true"}}, 0,
			"1", []string{models.PesanKasirBelumMasuk}},
		{"tanpa adjustment, tanpa lampiran, estimasi ada", nil, 0, "10000000",
			[]string{models.PesanTanpaPembayaran, models.PesanEstimasiBelumNol}},
		{"tanpa adjustment, lampiran Close Claim, estimasi nol", nil, 1, "0", nil},
		{"siap ditutup", []models.Baris{{"AcceptanceStatus": "1", "IsPrintAccept": "1", "DirectToKasir": "false"}}, 0, "1",
			nil},
	}
	for _, c := range kasus {
		t.Run(c.nama, func(t *testing.T) {
			h := halamanAdj(models.Baris{})
			h.SetelDaftar(models.DaftarAdj(1, 1), c.adj)
			it, _ := models.Item(h, 1, 1)
			it["TotalEstimationValueinIDR"] = c.estimasi
			if err := models.ValidasiTutup(h, c.lampiran); err != nil {
				t.Fatal(err)
			}
			got := strings.Join(h.SemuaPesan(), "|")
			if len(c.pesan) == 0 && got != "" {
				t.Fatalf("pesan tak terduga %q", got)
			}
			for _, p := range c.pesan {
				if !strings.Contains(got, p) {
					t.Fatalf("pesan %q tidak ada di %q", p, got)
				}
			}
		})
	}
}

// Perbaikan prompt §5 butir 8 (model): proteksi DLA menghasilkan pesan yang - di services - menghentikan penutupan.
func TestCekDLATutup(t *testing.T) {
	h := halamanAdj(models.Baris{})
	h.SetelDaftar(models.DaftarObjek, []models.Baris{{"IsFacretro": "1", "RemarksDLA": ""}})
	h.Setel(models.JalurTKRemarks, "UJI")
	models.CekDLATutup(h)
	if !strings.Contains(strings.Join(h.SemuaPesan(), "|"), models.PesanDLASebelumTutup) ||
		h.Ambil(models.CD+"Remark_Close") != "UJI" {
		t.Fatalf("pesan %v", h.SemuaPesan())
	}
}

func TestBarisOSTutup(t *testing.T) {
	h := models.HalamanBaru()
	h.Setel(models.CD+"NoClaim", "UJI-K77.03.2026.00001")
	h.Setel(models.CD+"CauseOfLoss", "UJI KEBAKARAN")
	h.Setel(models.JalurNoPolis, "UJI-RNM-F.001")
	b := models.BarisOSTutup(h, "CLM-000001")
	mau := "{\n\"CauseOfLoss\":\"UJI KEBAKARAN\"\n,\"NoClaim\":\"UJI-K77.03.2026.00001\"\n," +
		"\"pxObjClass\":\"ASM-FW-GCNMFW-Data-osAkseptasi\"\n}\n"
	if b.StsReject != "4" || b.CaseID != "CLM-000001" || b.StsDLA != "" || b.DataJSON != mau {
		t.Fatalf("%+v\n%q", b, b.DataJSON)
	}
}
