package repository

import (
	"context"
	"testing"

	"nusantarare/inti/db"
)

func TestBacaTxTCO(t *testing.T) {
	var d db.DB
	if DenganBacaTxTCO(context.Background(), nil) != context.Background() {
		t.Error("tx nil harus mengembalikan ctx apa adanya")
	}
	// Tx tanpa *sql.Tx (uji tanpa Oracle) tidak dipakai membaca.
	if bacaTCO(DenganBacaTxTCO(context.Background(), &db.Tx{}), &d) != kuerierTCO(&d) {
		t.Error("tx kosong dipakai membaca")
	}
	if bacaTCO(context.Background(), &d) != kuerierTCO(&d) {
		t.Error("tanpa penanda harus membaca pool")
	}
}
