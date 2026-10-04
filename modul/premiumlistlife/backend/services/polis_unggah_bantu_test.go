package services

// Pembantu uji unggahan - kolom wajib bergantung pada Type sejak 03-10-2026;
// uji lama memakai QR (Gross Valuation).

import (
	"io"

	"nusantarare/modul/premiumlistlife/backend/models"
)

func bacaQR(r io.Reader) ([]models.BarisUnggah, error) { return BacaCSVUnggah(r, "QR") }

func validasiQR(baris []models.BarisUnggah) models.HasilUnggah {
	h, err := models.ValidasiUnggah("QR", baris)
	if err != nil {
		panic(err)
	}
	return h
}
