package models

// Pembantu uji unggahan - kolom wajib bergantung pada Type sejak 03-10-2026;
// uji lama memakai QR (Gross Valuation).

func validasiQR(baris []BarisUnggah) HasilUnggah {
	h, err := ValidasiUnggah("QR", baris)
	if err != nil {
		panic(err)
	}
	return h
}

func kolomWajibQR() []string {
	k, err := KolomWajibUnggah("QR")
	if err != nil {
		panic(err)
	}
	return k
}
