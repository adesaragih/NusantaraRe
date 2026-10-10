package handlers_test

import (
	"context"

	"nusantarare/modul/treatyin/backend/repository"
)

// Gudang tiruan rute: Submit Achievement — nol tulisan.
func (g gudangTiruan) CatatLogAchievement(_ context.Context, _, _ string, _ []repository.BarisLogSiap) error {
	return nil
}
