package services

import (
	"fmt"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// periksaMask - setiap desimal yang akan diikat ke angkaMasuk wajib cocok repository.PolaAngkaMasuk (<= 30 digit bulat,
// <= 8 desimal, tanpa eksponen); selain itu Oracle menolak ORA-01722. Dipakai tiruan Tulis* supaya kelas bug "nilai skala
// 20 sampai ke penyimpan" (DEV 08-10-2026) tertangkap uji.
func periksaMask(nama string, nilai ...*apd.Decimal) error {
	for i, x := range nilai {
		if x != nil && !repository.PolaAngkaMasuk.MatchString(utils.FormatDecimal(x)) {
			return fmt.Errorf("UJI mask angkaMasuk %s[%d]: %q", nama, i, utils.FormatDecimal(x))
		}
	}
	return nil
}

// maskSpreading - nilai yang ditulis TulisSpreading lewat angkaMasuk.
func maskSpreading(k models.KasusSpreading) error {
	if err := periksaMask("PercentShare", k.PercentShare); err != nil {
		return err
	}
	for _, l := range k.Lokasi {
		for _, it := range l.Items {
			for _, c := range it.Coverages {
				if err := periksaMask("NR "+c.OldID, c.TSINusantaraRe, c.PremiNusantaraRe); err != nil {
					return err
				}
				for _, b := range c.Spreading {
					if err := periksaMask("baris "+c.OldID, b.SharePercentage, b.TSISpreaded, b.TSIGrossSpreaded, b.PremiumSpreaded,
						b.ClaimEstimation); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// maskBarisDeduction - CurrencyList yang ditulis TulisDeduction / TulisPayment.
func maskBarisDeduction(baris []models.BarisDeduction) error {
	for _, b := range baris {
		if err := periksaMask("CurrencyList "+b.Currency, b.Premium, b.Commision, b.Deduction2, b.BrokerageFee, b.PPh, b.PPN); err != nil {
			return err
		}
	}
	return nil
}
