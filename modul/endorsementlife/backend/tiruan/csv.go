package tiruan

// Tiruan `SaveCSVEDMLife` (repository/edm_csv.go).

import (
	"context"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/endorsementlife/backend/models"
)

// AcuanCSV meniru `sqlAcuanCSV`: peserta pertama urutan grid, di luar `New`.
func (g *Gudang) AcuanCSV(_ context.Context, _ *db.Tx, kasusID string) (models.AcuanCSV, bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return models.AcuanCSV{}, false, g.Galat
	}
	for _, d := range g.pesertaKasus(kasusID) {
		if d.EdmStatus != models.StatusNew {
			return models.AcuanCSV{Plan: d.Nilai["PLAN"], PolicyHolder: d.Nilai["POLICY_HOLDER"]}, true, nil
		}
	}
	return models.AcuanCSV{}, false, nil
}

// HapusPesertaBaru meniru `sqlHapusPesertaBaru`.
func (g *Gudang) HapusPesertaBaru(_ context.Context, _ *db.Tx, kasusID string) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return 0, g.Galat
	}
	sisa, c := g.Peserta[:0], 0
	for _, d := range g.Peserta {
		if d.PolisID == kasusID && d.EdmStatus == models.StatusNew {
			c++
			continue
		}
		sisa = append(sisa, d)
	}
	g.Peserta = sisa
	return c, nil
}

// SisipPesertaCSV meniru `sqlSisipPesertaCSV`; `GagalSisipKe` > 0 menggagalkan
// penyisipan baris bernomor itu (uji pembatalan transaksi).
func (g *Gudang) SisipPesertaCSV(_ context.Context, _ *db.Tx, kasusID, plNumber string, baris []models.BarisCSV) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return 0, g.Galat
	}
	for _, b := range baris {
		if g.GagalSisipKe > 0 && b.Nomor == g.GagalSisipKe {
			return 0, fmt.Errorf("tiruan: sisip baris CSV %d gagal", b.Nomor)
		}
		nilai := map[string]string{}
		for k, v := range b.Nilai {
			nilai[k] = v
		}
		g.Peserta = append(g.Peserta, &Peserta{ID: fmt.Sprintf("UJI-CSV-%s-%d", kasusID, b.Nomor), PolisID: kasusID,
			PLNumber: plNumber, EdmStatus: models.StatusNew, Nilai: nilai})
	}
	return len(baris), nil
}
