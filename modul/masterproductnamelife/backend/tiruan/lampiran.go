package tiruan

// Lampiran di memori: `M_ATTACHMENTPRODUCTNAME`, `T_STORAGE_IMAGE`, `T_FOLDER_IMAGE`,
// dan baris outbox `T_LOG_SERVICE_RNM` modul ini - semantik sama dengan repository
// (status terunggah = objek tercatat; selain itu dari efek terakhir).

import (
	"context"
	"fmt"
	"sort"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/repository"
)

// Efek - satu baris outbox tiruan.
type Efek struct {
	ID, Rujukan, Status, Galat string
	Percobaan                  int
	Jadwal                     time.Time
}

func (g *Gudang) siapLampiran() {
	if g.Lampiran == nil {
		g.Lampiran = map[string]models.Lampiran{}
	}
	if g.Objek == nil {
		g.Objek = map[string]models.ObjekPenyimpanan{}
	}
}

func (g *Gudang) status(l models.Lampiran) models.Lampiran {
	l.Status, l.Galat = models.StatusBelum, ""
	if _, ada := g.Objek[l.StorageID]; ada {
		l.Status = models.StatusTerunggah
		return l
	}
	for i := len(g.Outbox) - 1; i >= 0; i-- {
		e := g.Outbox[i]
		if e.Rujukan != l.ID {
			continue
		}
		if e.Galat != "" || e.Status == outbox.StatusEfekGagalPermanen {
			l.Status, l.Galat = models.StatusGagal, e.Galat
		}
		break
	}
	return l
}

// DaftarLampiran - urut ID.
func (g *Gudang) DaftarLampiran(_ context.Context, produkID string) ([]models.Lampiran, error) {
	g.siapLampiran()
	hasil := []models.Lampiran{}
	for _, l := range g.Lampiran {
		if l.ProdukID == produkID {
			hasil = append(hasil, g.status(l))
		}
	}
	sort.Slice(hasil, func(i, j int) bool { return hasil[i].ID < hasil[j].ID })
	return hasil, nil
}

// AmbilLampiran - satu lampiran produk.
func (g *Gudang) AmbilLampiran(_ context.Context, _ *db.Tx, produkID, id string) (models.Lampiran, error) {
	g.siapLampiran()
	l, ada := g.Lampiran[id]
	if !ada || l.ProdukID != produkID {
		return models.Lampiran{}, fmt.Errorf("%w: %s", repository.ErrLampiranTidakAda, id)
	}
	return g.status(l), nil
}

// SisipLampiran - ID berurut bergaya stempel FF3.
func (g *Gudang) SisipLampiran(_ context.Context, _ *db.Tx, l models.Lampiran) (string, error) {
	g.siapLampiran()
	g.nomorLampiran++
	l.ID = fmt.Sprintf("20261001000000%03d", g.nomorLampiran)
	g.Lampiran[l.ID] = l
	return l.ID, nil
}

// HapusLampiran - satu baris.
func (g *Gudang) HapusLampiran(_ context.Context, _ *db.Tx, produkID, id string) error {
	if l, ada := g.Lampiran[id]; !ada || l.ProdukID != produkID {
		return fmt.Errorf("%w: %s", repository.ErrLampiranTidakAda, id)
	}
	delete(g.Lampiran, id)
	return nil
}

// CatatObjek - idempoten.
func (g *Gudang) CatatObjek(_ context.Context, _ *db.Tx, o models.ObjekPenyimpanan) error {
	g.siapLampiran()
	if _, ada := g.Objek[o.ImageID]; !ada {
		g.Objek[o.ImageID] = o
	}
	return nil
}

// HapusObjek - yang tidak ada bukan galat.
func (g *Gudang) HapusObjek(_ context.Context, _ *db.Tx, imageID string) error {
	delete(g.Objek, imageID)
	return nil
}

// NamaAplikasi - `T_FOLDER_IMAGE.APPNAME`.
func (g *Gudang) NamaAplikasi(context.Context, *db.Tx) (string, error) { return g.AppName, nil }

// AntreUnggah - satu efek antre.
func (g *Gudang) AntreUnggah(_ context.Context, _ *db.Tx, lampiranID, _ string, saat time.Time) error {
	g.Outbox = append(g.Outbox, Efek{ID: fmt.Sprint(len(g.Outbox) + 1), Rujukan: lampiranID, Status: outbox.StatusEfekAntre,
		Jadwal: saat})
	return nil
}

// PungutUnggah - efek antre pertama lampiran itu, ditandai jalan.
func (g *Gudang) PungutUnggah(_ context.Context, _ *db.Tx, lampiranID string, _ time.Time) (string, int, bool, error) {
	for i := range g.Outbox {
		e := &g.Outbox[i]
		if e.Rujukan == lampiranID && e.Status == outbox.StatusEfekAntre {
			e.Status = outbox.StatusEfekJalan
			e.Percobaan++
			return e.ID, e.Percobaan, true, nil
		}
	}
	return "", 0, false, nil
}

// AdaUnggahAntre - efek antre lampiran ini ada (tanpa menandai jalan).
func (g *Gudang) AdaUnggahAntre(_ context.Context, _ *db.Tx, lampiranID string) (bool, error) {
	for _, e := range g.Outbox {
		if e.Rujukan == lampiranID && e.Status == outbox.StatusEfekAntre {
			return true, nil
		}
	}
	return false, nil
}

// TuntaskanUnggah - hasil satu percobaan.
func (g *Gudang) TuntaskanUnggah(_ context.Context, _ *db.Tx, efekID, status string, jadwal time.Time, galat string,
	_ time.Time) error {
	for i := range g.Outbox {
		if g.Outbox[i].ID == efekID {
			g.Outbox[i].Status, g.Outbox[i].Jadwal, g.Outbox[i].Galat = status, jadwal, galat
			return nil
		}
	}
	return fmt.Errorf("tiruan: no effect %s", efekID)
}
