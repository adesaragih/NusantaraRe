package models_test

// Kolom Type grid Estimation List menampilkan teks prompt (screenshot work owner 08-10-2026, properti
// ASM-FW-GCNMFW-Data-Estimasi.Type: 1 Claim, 2 Adjuster Fee, 3 Salvage, 4 Consultant Fee); nilai tersimpan tetap kode.
// Kunci label harus SAMA dengan kunci sumber kolom itu - layar membaca `labelKode[<kunci sumber tanpa "kode:">]`.

import (
	"strings"
	"testing"

	"nusantarare/modul/claimprop/backend/models"
)

func kolomGrid(ts []models.Tata, daftar, jalur string) *models.Tata {
	for i := range ts {
		if ts[i].Jenis == models.JenisGrid && ts[i].Jalur == daftar {
			for j := range ts[i].Kolom {
				if ts[i].Kolom[j].Jalur == jalur {
					return &ts[i].Kolom[j]
				}
			}
		}
		if k := kolomGrid(ts[i].Anak, daftar, jalur); k != nil {
			return k
		}
	}
	return nil
}

func TestLabelTypeEstimasi(t *testing.T) {
	mau := map[string]string{"1": "Claim", "2": "Adjuster Fee", "3": "Salvage", "4": "Consultant Fee"}
	for nama, layar := range map[string]func() []models.Unsur{
		"Outstanding": models.LayarOutstanding, "Akseptasi": models.LayarAkseptasi,
	} {
		h := models.HalamanBaru()
		h.SetelDaftar(models.DaftarEstimasi, []models.Baris{{"Type": "1"}})
		k := kolomGrid(models.Evaluasi(h, layar(), false), models.DaftarEstimasi, "Type")
		if k == nil {
			t.Fatalf("%s: kolom Type Estimation List tidak ada", nama)
		}
		kunci, ada := models.AdaKode(k.Sumber)
		if !ada || kunci != "EstimationType" {
			t.Fatalf("%s: sumber kolom Type %q", nama, k.Sumber)
		}
		for _, kode := range models.KodePilihan[kunci] {
			if got := models.LabelKode[kunci][kode]; got != mau[kode] {
				t.Fatalf("%s: label %s = %q, mau %q", nama, kode, got, mau[kode])
			}
		}
		if got := strings.Join(models.KodePilihan[kunci], ","); got != "1,2,3,4" {
			t.Fatalf("kode tersimpan berubah: %s", got)
		}
	}
}
