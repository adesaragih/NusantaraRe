package models_test

// Consultant / Adjuster dipilih lewat NAMA (perintah work owner 08-10-2026: "yang dropdown hanya dari namanya aja; untuk
// ID dihapus dari tampilan, tapi tetap simpan ID"): medan tetap berjalur ID (yang disimpan), berlabel nama, dan
// menampilkan teks dari jalur nama. Baris nama hanya-baca hanya muncul saat dropdown tersembunyi (IsAnyAcceptation = 1).

import (
	"testing"

	"nusantarare/modul/claimprop/backend/models"
)

func medanBerjalur(ts []models.Tata, jalur string) []models.Tata {
	var out []models.Tata
	for _, t := range ts {
		if t.Jenis == models.JenisMedan && t.Jalur == jalur {
			out = append(out, t)
		}
		out = append(out, medanBerjalur(t.Anak, jalur)...)
	}
	return out
}

func TestAdjusterDipilihLewatNama(t *testing.T) {
	pasangan := []struct{ id, nama, label string }{
		{models.CD + "ConsultantID", models.CD + "ConsultantName", "Consultant Name"},
		{models.CD + "AppointedADJID", models.CD + "AppointedADJ", "Adjuster / Professional Name"},
	}
	for nama, layar := range map[string]func() []models.Unsur{
		"Outstanding": models.LayarOutstanding, "Akseptasi": models.LayarAkseptasi,
	} {
		for _, sudahAcc := range []string{"", "1"} {
			h := models.HalamanBaru()
			h.Setel("IsAnyAcceptation", sudahAcc)
			for _, p := range pasangan {
				h.Setel(p.id, "UJI-ADJ-1")
				h.Setel(p.nama, "UJI NAMA ADJUSTER")
			}
			ts := models.Evaluasi(h, layar(), false)
			for _, p := range pasangan {
				pilih, baca := medanBerjalur(ts, p.id), medanBerjalur(ts, p.nama)
				if sudahAcc == "1" {
					if len(pilih) != 0 || len(baca) != 1 {
						t.Fatalf("%s %s sudah akseptasi: dropdown %d, baris nama %d", nama, p.id, len(pilih), len(baca))
					}
					continue
				}
				if len(pilih) != 1 || len(baca) != 0 {
					t.Fatalf("%s %s: dropdown %d, baris nama %d (nama ditampilkan dropdown)", nama, p.id, len(pilih), len(baca))
				}
				if pilih[0].Label != p.label || pilih[0].Tampilan != p.nama || pilih[0].Sumber != models.SumberAdjuster {
					t.Fatalf("%s %s: %+v", nama, p.id, pilih[0])
				}
			}
		}
	}
}
