package services

// Tiket 59 lanjutan (K59-B6, diminta sesi 9d 08-10-2026): Save tab Spreading di tahap uw ikut menulis ulang nilai turunan tab
// Premium Deduction / Payment yang tersimpan, di transaksi yang SAMA dengan TulisSpreading. Port `[terverifikasi]`
// `D:\migrasi\RNM\NB FacIn\Activity\CountPremiAndTSINusantaraRe_ACT.xml` langkah 11 "WHEN UW" (prekondisi IsUWSpread):
//   - CountPremiCedant_Act - `CedingCedantList(i).CurrencyList` TSI / Premium per cedant; sistem baru TIDAK menyimpannya
//     (T_CEDINGCEDANTLIST tanpa kolom per mata uang, migrasi 199; dasar per mata uang dihitung saat GET tab Cedant, K49-8) -
//     tidak ada yang ditulis. SHARE_OF_CEDING tidak disentuh (langkah 11 tidak mengubahnya; SetShareOfCeding hanya klik
//     ShareCedantType).
//   - SetRIComm_Act - hitungRIComm (rumus tiket 50) dari persen TERSIMPAN atas premi Nusantara Re baru; langkah 5
//     `PPnCheck = true` SELALU; Deduction2 tidak disentuh (nominal tersimpan per mata uang). Cabang EDM (1, 6, 8.2) di luar
//     cakupan.
//   - CountPaymentInstallment_Act - segarkanCicilan (rumus tiket 51 / K50-12): % dan Due Date tersimpan, nilai uang / Total
//     Payment / Net Premium dari nilai Payment baru. NetPremium mengurangkan Deduction2 (K51-1, penyimpangan sadar keputusan
//     work owner 06-10-2026 - CountPaymentInstallment_Act langkah 3.2 tidak mengurangkannya).
// Tabel yang ditulis = tabel yang SUDAH ditulis Save tab Premium Deduction jalur marketing (T_GENERAL_POLIS PAY_* / PPN_CHECK,
// T_CURRENCYLIST, T_LISTINSTALLMENT, PAY_INSTALLMENT); tanpa migrasi. Endpoint tab itu tetap 409 di uw.

import (
	"context"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// dasarLangkah11 - data tersimpan tab Premium Deduction / Payment + tarif TAX, dibaca SEBELUM transaksi (pola SimpanDeduction).
type dasarLangkah11 struct {
	deduction models.KasusDeduction
	bayar     models.KasusPayment
	tarif     models.TarifPajak
}

// adaDeductionTersimpan - tab Premium Deduction pernah di-Save (persen atau baris CurrencyList ada).
func adaDeductionTersimpan(k models.KasusDeduction) bool {
	return k.RICommision != nil || k.PctBrokerageFee != nil || len(k.MataUang) > 0
}

// bacaDasarLangkah11 - nil bila penyimpan Premium Deduction tidak terpasang atau tab itu belum pernah di-Save (tidak ada nilai
// tersimpan yang basi; GET tab menghitungnya dari spreading).
func (s *Service) bacaDasarLangkah11(ctx context.Context, id string) (*dasarLangkah11, error) {
	if s.deduction == nil {
		return nil, nil
	}
	kd, err := s.deduction.BacaDeduction(ctx, id)
	if err != nil {
		return nil, galatKasus(err)
	}
	if !adaDeductionTersimpan(kd) {
		return nil, nil
	}
	tarif, _, err := s.deduction.TarifPajak(ctx)
	if err != nil {
		return nil, err
	}
	d := &dasarLangkah11{deduction: kd, tarif: tarif}
	if s.payment != nil {
		if d.bayar, err = s.payment.BacaPayment(ctx, id); err != nil {
			return nil, galatKasus(err)
		}
	}
	return d, nil
}

// barisLangkah11 - CurrencyList yang akan ditulis: SetRIComm_Act atas spreading baru (`sp`, nilai 8 desimal yang sama dengan
// yang akan dibaca ulang). d nil = nil (tidak ada yang ditulis).
func barisLangkah11(sp models.KasusSpreading, d *dasarLangkah11) []models.BarisDeduction {
	if d == nil {
		return nil
	}
	baris := premiPerMataUang(sp)
	lama := map[string]models.BarisDeduction{}
	for _, b := range d.deduction.MataUang {
		lama[b.Currency] = b
	}
	for i := range baris {
		baris[i].Deduction2 = lama[baris[i].Currency].Deduction2
	}
	// Persen berpresisi penuh dibangun ulang dari baris tersimpan (K50-9), seperti Save tab Premium Deduction.
	ri := persenPenuh(d.deduction.RICommision, d.deduction.MataUang, nilaiCommision)
	pct := persenPenuh(d.deduction.PctBrokerageFee, d.deduction.MataUang, nilaiBrokerageFee)
	hitungRIComm(baris, ri, pct, d.tarif, true)
	bulatkanDeduction(baris)
	return baris
}

// tulisLangkah11UW - Premium Deduction (baris dari barisLangkah11) + CountPaymentInstallment_Act di transaksi `tx`. d nil =
// tidak ada yang ditulis.
func (s *Service) tulisLangkah11UW(ctx context.Context, tx *db.Tx, id string, d *dasarLangkah11, baris []models.BarisDeduction) error {
	if d == nil {
		return nil
	}
	if err := s.deduction.TulisDeduction(ctx, tx, id, models.KasusDeduction{RICommision: d.deduction.RICommision,
		PctBrokerageFee: d.deduction.PctBrokerageFee, PPnCheck: "true", MataUang: baris}); err != nil {
		return err
	}
	return s.segarkanCicilan(ctx, tx, id, d.bayar, baris)
}
