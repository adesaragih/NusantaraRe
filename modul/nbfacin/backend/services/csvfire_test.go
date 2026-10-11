package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestCsvFireHanyaKasusFire - bug DEV 11-10-2026 (Pega ChooseActivityMappingCSV_Act: InsertUploadFire_act hanya IsFire): case
// Aneka / lini tak terbaca -> 409 di Save Data JSON lama, unggah berkas, simpan, ulangi (pekerjaan yang sudah ada) dan jalur latar;
// objek case TIDAK berubah. Case FIRE (fixtur) tetap lolos. MARINE CARGO (tiket 74 P6) memakai importer cargo: Save Data JSON lama
// (template FIRE) tetap 409, dan pekerjaan FIRE yang sudah ada di case itu tidak dapat disimpan / diulang (jenis berbeda). Tiket 73
// P6: case Aneka bertemplate Aneka umum (bukan Marine Hull) -> ErrCsvTemplateAnekaBelumAda (K73-60), Save Data JSON lama tetap
// ErrCsvBukanFire.
func TestCsvFireHanyaKasusFire(t *testing.T) {
	ctx := context.Background()
	for nama, gp := range map[string][]models.BarisBisnisLini{"Aneka": satuBisnis("003", "30"), "lini tak terbaca": nil,
		"Marine Cargo": satuBisnis("004", "24")} {
		marine := nama == "Marine Cargo"
		svc, u, _, ob, _ := layananUnggahUji()
		svc.tableOfLimit = tolTiruan{bisnis: gp}
		lama := len(ob.ada["UJI-NB-1"])
		if _, err := svc.SimpanCsvObjek(ctx, pelakuCsv, "UJI-NB-1", []map[string]string{barisCsvUji()}); !errors.Is(err, ErrCsvBukanFire) {
			t.Errorf("%s JSON lama: %v", nama, err)
		}
		tolak := ErrCsvBukanFire
		if nama == "Aneka" {
			tolak = ErrCsvTemplateAnekaBelumAda
		}
		if marine {
			tolak = ErrUnggahStatus
		} else {
			if _, err := svc.MulaiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", "uji.csv", strings.NewReader(csvUji(3, nil))); !errors.Is(err, tolak) {
				t.Errorf("%s berkas: %v", nama, err)
			}
			if len(u.kerja) != 0 {
				t.Errorf("%s: pekerjaan dibuat %d", nama, len(u.kerja))
			}
		}
		if marine { // dependensi importer cargo terpasang -> sampai ke pemeriksaan jenis pekerjaan
			svc = svc.DenganObjekCargo(cargoTiruan{ada: map[string][]models.ObjekCargo{"UJI-NB-1": nil}}).DenganMasterCargo(masterUji()).
				DenganUnggahCsvCargo(&cargoCsvTiruan{}, &kapalTiruan{})
		}
		valid, _ := u.BuatUnggah(ctx, nil, models.UnggahCsv{CaseID: "UJI-NB-1", Status: StatusUnggahValid, Fase: FaseUnggahValidasi, TotalBaris: 3})
		if _, err := svc.SimpanUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", valid); !errors.Is(err, tolak) {
			t.Errorf("%s simpan: %v", nama, err)
		}
		u.kerja[valid].Status = StatusUnggahGagal
		u.kerja[valid].Fase = FaseUnggahSimpan
		if _, err := svc.UlangiUnggahCsv(ctx, pelakuCsv, "UJI-NB-1", valid); !errors.Is(err, tolak) {
			t.Errorf("%s ulangi: %v", nama, err)
		}
		svc.simpanUnggah(ctx, pelakuCsv, *u.kerja[valid])
		mau := "hanya untuk Class of Business FIRE"
		if nama == "Aneka" {
			mau = "unggah CSV baru tersedia untuk Marine Hull"
		}
		if marine {
			mau = "pekerjaan unggah CSV fire tidak sesuai Class of Business case (cargo)"
		}
		if x := u.kerja[valid]; x.Status != StatusUnggahGagal || !strings.Contains(x.GalatProses, mau) {
			t.Errorf("%s latar: %+v", nama, x)
		}
		if len(ob.ada["UJI-NB-1"]) != lama {
			t.Errorf("%s: objek berubah %d -> %d", nama, lama, len(ob.ada["UJI-NB-1"]))
		}
	}
}
