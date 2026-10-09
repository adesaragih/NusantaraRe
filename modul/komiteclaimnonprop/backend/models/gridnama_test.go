package models_test

import (
	"testing"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimnonprop/backend/models"
)

// layarUji - kasus tingkat 1 atas klaim berlayer (Section ShowTransfer).
func layarUji(t *testing.T, kl kontrak.KlaimTreaty) (models.Layar, func(judul string, baris int, label string) string) {
	t.Helper()
	ly := models.SusunLayar(kasusUji(1), kl, "UJI-K1", nil)
	sel := func(judul string, baris int, label string) string {
		t.Helper()
		for _, b := range ly.Bagian {
			for _, g := range b.Grid {
				if g.Judul != judul {
					continue
				}
				for _, k := range g.Kolom {
					if k.Label == label {
						return g.Baris[baris][k.Properti]
					}
				}
				t.Fatalf("grid %q tanpa kolom %q", judul, label)
			}
		}
		t.Fatalf("grid %q tidak ada", judul)
		return ""
	}
	return ly, sel
}

func medan(ly models.Layar, label string) (string, bool) {
	for _, b := range ly.Bagian {
		for _, m := range b.Medan {
			if m.Label == label {
				return m.Nilai, true
			}
		}
	}
	return "", false
}

func adaGrid(ly models.Layar, judul string) bool {
	for _, b := range ly.Bagian {
		for _, g := range b.Grid {
			if g.Judul == judul {
				return true
			}
		}
	}
	return false
}

// Grid klaim menampilkan nama mata uang / treaty; kode hanya bila namanya kosong.
func TestLayarShowTransferGridNamaBukanKode(t *testing.T) {
	kl := klaimAkhir()
	kl.Daftar["ClaimData.AdjustmentList"][0]["Type"] = "2"
	kl.Nilai["ClaimData.PolicyData.StartDateTime"], kl.Nilai["ClaimData.PolicyData.EndDateTime"] = "2026-01-01", "2026-12-31"
	kl.Daftar["ClaimData.AdjustmentList(1).ListClaimAcceptation"] = []map[string]string{{"CurrencyID": "UJI-9",
		"Value": "10"}}
	ly, sel := layarUji(t, kl)
	if ly.Judul[0] != "CLAIM COMMITTEE -" || len(ly.Judul) != 1 {
		t.Fatalf("judul S1: %v", ly.Judul)
	}
	cek := func(judul string, baris int, label, mau string) {
		t.Helper()
		if v := sel(judul, baris, label); v != mau {
			t.Errorf("%s baris %d kolom %s = %q, mau %q", judul, baris+1, label, v, mau)
		}
	}
	cek("", 0, "Currency", "UJI-9") // Claim Acceptation tanpa judul: nama kosong -> kode
	cek("", 0, "Claim Amount", "10")
	cek("XOL Allocation", 1, "Treaty Name", "UJI-L1")
	cek("XOL Allocation", 1, "Claim Amount RNM", "100")
	cek("XOL Allocation", 1, "Reinstatement Premium", "2")
	cek("Spreading In", 0, "Treaty Type", "UJI-QS")
	cek("Spreading In", 0, "Reinstatement Premium", "1.5")
	if v, _ := medan(ly, "Type"); v != "Adjuster Fee" {
		t.Errorf("Type akseptasi = %q, mau label", v)
	}
	if v, _ := medan(ly, "Claim No / Claim ID"); v != "UJI-K-0001 / CLMNP-UJI001" {
		t.Errorf("S9 Claim No / Claim ID = %q", v)
	}
	if v, _ := medan(ly, "Insurance Period"); v != "01/01/2026 - 31/12/2026" {
		t.Errorf("S7 Insurance Period = %q", v)
	}
	if v, _ := medan(ly, "RNM Share (%)"); v != "25" {
		t.Errorf("RNM Share (%%) tanpa PersenRNM akseptasi = RNMShare master: %q", v)
	}
}

// S23 / S36: Previously Calculated, Payable To, dan rekening hanya bila IsPrevious (AlokasiXOLPaid terisi); S7 Occupation
// NOTBLANK; tombol "View Claim" / Cancel / Submit.
func TestLayarShowTransferSyaratTampil(t *testing.T) {
	kl := klaimAkhir()
	ly, _ := layarUji(t, kl)
	if adaGrid(ly, "Previously Calculated") {
		t.Fatal("tanpa AlokasiXOLPaid: Previously Calculated tidak tampil")
	}
	if _, ada := medan(ly, "Payable To"); ada {
		t.Fatal("S36 tampil[IsPrevious==1]: Payable To tidak tampil")
	}
	if _, ada := medan(ly, "Occupation"); ada {
		t.Fatal("Occupation kosong tidak tampil")
	}
	kl.Nilai["ClaimData.Occupation"] = "UJI OKUPASI"
	kl.Daftar["ClaimData.AdjustmentList(1).AlokasiXOLPaid"] = []map[string]string{{"TreatyName": "UJI-L1",
		"ClaimSpreaded": "40"}}
	kl.Daftar["ClaimData.AdjustmentList"][0]["NameOfBank"] = "UJI BANK"
	kl.Daftar["ClaimData.AdjustmentList"][0]["Payable"] = "2"
	kl.Nilai["ClaimData.Payable"] = "1" // header klaim bukan sumber S37
	kl.Daftar["ClaimData.AdjustmentList"][0]["FlagCurrency"] = "1"
	ly, sel := layarUji(t, kl)
	if v := sel("Previously Calculated", 0, "Claim Amount RNM"); v != "40" {
		t.Fatalf("Previously Calculated: %q", v)
	}
	if v, _ := medan(ly, "Payable To"); v != "Broker Name" {
		t.Fatalf("S37 Payable To = .Adjustment.Payable: %q", v)
	}
	if v, _ := medan(ly, "Specify"); v != "UJI-PENERIMA" {
		t.Fatalf("S37 Specify = .Adjustment.PayableTo: %q", v)
	}
	if v, ada := medan(ly, "Name of Bank"); !ada || v != "UJI BANK" {
		t.Fatalf("rekening S38: %q %v", v, ada)
	}
	n := 0
	for _, b := range ly.Bagian {
		if b.Kunci == "bank2" {
			n++
		}
	}
	if n != 1 {
		t.Fatal("S40 rekening kedua bila FlagCurrency 1")
	}
	if v, ada := medan(ly, "Occupation"); !ada || v != "UJI OKUPASI" {
		t.Fatalf("Occupation: %q", v)
	}
	if len(ly.Tombol) != 3 || ly.Tombol[0].Label != "View Claim" || ly.Tombol[1].Label != "Cancel" ||
		ly.Tombol[2].Label != "Submit" || !ly.Tombol[2].Aktif {
		t.Fatalf("tombol: %+v", ly.Tombol)
	}
	if ly.Isian.Label["isSubjectivity"] != "Subjectivity" {
		t.Fatalf("label isian: %+v", ly.Isian.Label)
	}
}

func TestLayarTanggaCommitteAcceptStatus(t *testing.T) {
	k := kasusUji(2, models.KeputusanSetuju)
	k.Tangga[0].Komentar, k.Tangga[0].Tanggal = "UJI-OK", "2026-10-08 10:00:00"
	ly := models.SusunLayar(k, klaimAkhir(), "UJI-K2", nil)
	for _, b := range ly.Bagian {
		for _, g := range b.Grid {
			if g.Judul != "Committe Accept Status" {
				continue
			}
			if len(g.Baris) != 2 || g.Baris[0]["jabatan"] != "UJI-JABATAN-1" || g.Baris[0]["keputusan"] != "Approved" ||
				g.Baris[0]["komentar"] != "UJI-OK" || g.Baris[1]["keputusan"] != "Waiting" {
				t.Fatalf("tangga: %+v", g.Baris)
			}
			if ly.Isian.Nilai.Comment != "UJI-OK" {
				t.Fatalf("Note tingkat 2 = komentar tingkat sebelumnya (pyWorkPage.Comment): %q", ly.Isian.Nilai.Comment)
			}
			return
		}
	}
	t.Fatal("grid Committe Accept Status tidak ada")
}
