package tiruan

// Tiruan `BrowseReinstypeOR_SQL` b84 (paket 9) - saringan yang sama dengan SQL-nya.

import (
	"context"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/repository"
)

// KontrakOR - satu baris kontrak treaty bergabung tahun treaty-nya.
type KontrakOR struct {
	ReinsTypeID, ReinsTypeName   string
	TreatyYear, UnderwritingYear string
	Mulai, Akhir                 time.Time
}

// DaftarReinstypeOR - `TO_DATE(BEGIN) >= TREATYSTARTDATE AND TO_DATE(MATURE) <= TREATYENDDATE AND REINSTYPEID = '10200'`.
func (g *Gudang) DaftarReinstypeOR(_ context.Context, _ *db.Tx, begin, mature string) ([]models.BarisOutward, error) {
	g.MintaOR = [2]string{begin, mature}
	if g.GagalMaster != nil {
		return nil, g.GagalMaster
	}
	hasil := []models.BarisOutward{}
	b, errB := time.Parse(repository.BentukTanggalPega, begin)
	m, errM := time.Parse(repository.BentukTanggalPega, mature)
	if errB != nil || errM != nil {
		return hasil, nil
	}
	for _, k := range g.KontrakOR {
		if k.ReinsTypeID == repository.ReinsTypeOR && !b.Before(k.Mulai) && !m.After(k.Akhir) {
			hasil = append(hasil, models.BarisOutward{ReinsTypeID: k.ReinsTypeID, ReinsTypeName: k.ReinsTypeName,
				TransactionYear: k.TreatyYear, UnderwritingYear: k.UnderwritingYear})
		}
	}
	return hasil, nil
}
