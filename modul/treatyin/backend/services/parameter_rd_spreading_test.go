package services_test

// ⛔ PARAMETER RD SUSUNAN — dropdown dan Activity TIDAK SAMA, dan menyamakannya
// mengosongkan dropdown.
//
// `BrowseTreatyArrangement_ParentReinsMasterTrt` punya delapan filter dan NOL
// di antaranya ber-`pyUseNullIfEmpty`: parameter kosong membuat filternya
// DILEWATI. Jadi yang menentukan isinya bukan RD-nya, melainkan parameter apa
// yang tiap pemanggil kirim — dan itu dibaca satu per satu dari ekspor
// 8 Oktober 2026:
//
//	Section/Share.xml                 grup —   desc —        (dropdown XOL)
//	Section/DetailShare.xml           grup —   desc "10001"  (dropdown Prop)
//	Activity/FetchQSfromMaster[2]     grup ✓   desc "10001"
//	Activity/FetchQSfromMasterXOL[4]  grup ✓   desc "10001"
//
// ⚠️ Uji ini menjaga BARIS KEDUA tabel itu, yang sempat salah dua kali dalam
// satu giliran: mula-mula `TreatyDescID` dipakukan untuk SEMUA pemanggil
// (dropdown XOL ikut tersaring, dan itu laporan *"kenapa tidak bisa milih"*),
// lalu ralatnya menghapusnya dari SEMUA pemanggil — termasuk Activity, yang
// di Pega memang mengirimnya.

import (
	"context"
	"testing"

	"nusantarare/modul/treatyin/backend/services"
)

func TestActivityShareTetapMengirimTreatyDescID(t *testing.T) {
	for _, u := range []struct {
		nama string
		muat func(*gudangTiruan) error
	}{
		{"Prop (FetchQSfromMaster)", func(g *gudangTiruan) error {
			// Satu Detail ber-`SpreadingTypeID` — itulah yang memanggil RD induk.
			limits := []map[string]any{{
				"TreatyType": "QUOTA SHARE",
				"Detail": []any{map[string]any{
					"TreatyGroup": "PROPERTY", "TreatyGroupID": "10007",
					"RNMShare": "25", "SpreadingTypeID": "10252",
				}},
			}}
			_, err := services.LayananDengan(g).HitungShareProp(context.Background(), pelakuAda, services.MasukanShareProp{
				Aksi:         services.AksiSharePropSpreading,
				Limits:       limits,
				Commencement: "20250101",
			})
			return err
		}},
	} {
		g := &gudangTiruan{}
		if err := u.muat(g); err != nil {
			t.Fatalf("%s: %v", u.nama, err)
		}
		// ⚠️ Tanpa ini uji lulus ketika RD-nya tidak dipanggil sama sekali —
		// loop atas daftar kosong selalu hijau.
		if len(g.indukSpreading) == 0 {
			t.Fatalf("%s: RD induk tidak dipanggil sama sekali", u.nama)
		}
		for _, p := range g.indukSpreading {
			if p[1] != services.DescSpreading {
				t.Errorf("%s: TreatyDescID %q, mau %q — `FetchQSfromMaster[2]` mengirimnya",
					u.nama, p[1], services.DescSpreading)
			}
		}
	}
}

// ⭐ Dan nilainya memang yang ada di ekspor, bukan angka yang hanyut.
func TestDescSpreadingSamaDenganEkspor(t *testing.T) {
	if services.DescSpreading != "10001" {
		t.Fatalf("DescSpreading = %q, mau 10001 (Param.TreatyDescID seluruh pemanggil)", services.DescSpreading)
	}
}
