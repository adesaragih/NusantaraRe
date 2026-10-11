package services

// Gerbang CSV objek template FIRE (bug DEV 11-10-2026: unggahan CSV FIRE ±1 juta baris masuk ke case MARINE CARGO). Pega
// `NB FacIn\Activity\ChooseActivityMappingCSV_Act.xml` menjalankan InsertUploadFire_act (L439) HANYA bila When IsFire
// (prakondisi L499); Marine Cargo -> InsertUploadMarineCargo_act (L539 / L599), Aneka -> InsertUploadAneka_act (L639 / L699),
// Briguna -> InsertUploadCredit_Act (L243 / L303) `[terverifikasi]` (dibaca sesi 0f). Jalur CSV FIRE sistem baru (tiket 66 JSON,
// tiket 75 berkas / simpan / ulangi + jalur latar) memakai predikat yang SAMA dengan `lini.fire` GET kasus (LiniKasus). Lini tak
// terbaca (BUSINESS bukan tepat satu baris) = bukan FIRE - ditolak, seperti Pega (IsFire salah).

import (
	"context"
	"errors"

	"nusantarare/modul/nbfacin/backend/models"
)

// ErrCsvBukanFire - CSV template FIRE untuk case bukan FIRE. 409.
var ErrCsvBukanFire = errors.New("Upload CSV template FIRE hanya untuk Class of Business FIRE (Pega ChooseActivityMappingCSV_Act). " +
	"CSV untuk Class of Business ini belum tersedia di sistem baru.")

// jenisCsvKasus - importer unggah CSV objek menurut Class of Business case, persis ChooseActivityMappingCSV_Act: IsFire -> fire,
// IsMarineCargo -> cargo (tiket 74 P6), IsAneka -> aneka (tiket 73 P6) HANYA bila templatenya Marine Hull (K73-60; template
// Aneka umum / CPM belum ada -> ErrCsvTemplateAnekaBelumAda 409); selain itu ErrCsvBukanFire (409).
func (s *Service) jenisCsvKasus(ctx context.Context, id string) (string, error) {
	l := s.LiniKasus(ctx, id)
	switch {
	case l.Fire:
		return models.JenisCsvFire, nil
	case l.MarineCargo:
		return models.JenisCsvCargo, nil
	case l.Aneka && l.JenisCsv == models.TemplateCsvMarineHull:
		return models.JenisCsvAneka, nil
	case l.Aneka:
		return "", ErrCsvTemplateAnekaBelumAda
	}
	return "", ErrCsvBukanFire
}

// wajibCsvFire - ErrCsvBukanFire bila case bukan IsFire (Save Data JSON lama tiket 66 = template FIRE saja).
func (s *Service) wajibCsvFire(ctx context.Context, id string) error {
	if !s.LiniKasus(ctx, id).Fire {
		return ErrCsvBukanFire
	}
	return nil
}
