package models_test

// Kolom "Status" grid "Committe Accept Status" menampilkan teks prompt `.KomiteAproval` (laporan work owner 09-10-2026:
// kolom tampil "0"); nilai tersimpan tetap kode. Sumber: ekspor work owner `Komite Claim Prop/KomiteAproval.xml`
// (ASM-FW-GCNMFW-Data-Comitee pyLocalList: 1 Approved, 2 Reject, 0 Waiting).

import (
	"strings"
	"testing"

	"nusantarare/modul/claimprop/backend/models"
)

func TestLabelStatusKomite(t *testing.T) {
	h := models.HalamanBaru()
	h.SetelDaftar(models.DaftarAdjustment, []models.Baris{{"Type": "1"}})
	h.SetelDaftar(models.JalurAdj(1, "ComiteeClaim"), []models.Baris{{models.PropKeputusanAnggota: "0"}})
	k := kolomGrid(models.Evaluasi(h, models.LayarAdjustment(1), false), models.JalurAdj(1, "ComiteeClaim"),
		models.PropKeputusanAnggota)
	if k == nil {
		t.Fatal("kolom Status grid Committe Accept Status tidak ada")
	}
	kunci, ada := models.AdaKode(k.Sumber)
	if !ada || kunci != "KomiteAproval" {
		t.Fatalf("sumber kolom Status %q, mau kode KomiteAproval", k.Sumber)
	}
	mau := map[string]string{"0": "Waiting", "1": "Approved", "2": "Reject"}
	for kode, label := range mau {
		if got := models.LabelKode[kunci][kode]; got != label {
			t.Errorf("label %s = %q, mau %q", kode, got, label)
		}
	}
	if got := strings.Join(models.KodePilihan[kunci], ","); got != "1,2,0" {
		t.Errorf("kode tersimpan / urutan XML: %s, mau 1,2,0", got)
	}
}
