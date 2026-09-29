package repository

import (
	"context"
	"testing"
)

func TestBacaTxTCO(t *testing.T) {
	var d DB
	if DenganBacaTxTCO(context.Background(), nil) != context.Background() {
		t.Error("tx nil harus mengembalikan ctx apa adanya")
	}
	// Tx tanpa *sql.Tx (uji tanpa Oracle) tidak dipakai membaca.
	if d.bacaTCO(DenganBacaTxTCO(context.Background(), &Tx{})) != d.sql {
		t.Error("tx kosong dipakai membaca")
	}
	if d.bacaTCO(context.Background()) != d.sql {
		t.Error("tanpa penanda harus membaca pool")
	}
}
