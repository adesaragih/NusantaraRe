package services

import (
	"context"
	"errors"
	"strings"

	"nusantarare/modul/ririsklife/backend/models"
	"nusantarare/modul/ririsklife/backend/repository"
)

// SimpanRincian - Save form R/I RISK DETAIL (`InboxRIRisk` Save b3132): idRincian kosong = tambah baris (`AddToList_Act`
// b3155), selain itu ubah baris yang diisi EDIT b9784 (`EditRIRiskLife_Act` b9808). USEDBY (`TempIDUsedBy.USEDBY`
// b1727, disabled) dan IDUSEDBY (`TempIDUsedBy.ID` b1524) SELALU dari ringkasan idRingkasan - bukan isian. Rule
// `AddToList_Act` tidak tersedia; aturan di bawahnya ASUMSI (MODUL.md):
//   - kombinasi (CONTRACT, YEAR, MONTH) tidak boleh kembar di ringkasan yang sama (seperti Upload);
//   - MODIFY OPERATOR / MODIFY DATE ringkasan ikut diperbarui (seperti Upload).
//
// Satu transaksi: rincian + ringkasan.
func (l *Layanan) SimpanRincian(ctx context.Context, a Aktor, idRingkasan, idRincian string, isi models.IsianRincian) (models.Rincian, error) {
	if err := wajibPenuh(a); err != nil {
		return models.Rincian{}, err
	}
	idRingkasan, idRincian = strings.TrimSpace(idRingkasan), strings.TrimSpace(idRincian)
	rapi, err := models.PeriksaIsianRincian(isi)
	if err != nil {
		return models.Rincian{}, tolak("%s", err.Error())
	}
	var hasil models.Rincian
	err = l.tx(ctx, func(tx *dbTx) error {
		ring, err := l.gudang.Ambil(ctx, tx, idRingkasan)
		if err != nil {
			return err
		}
		lama, err := l.gudang.RincianDari(ctx, tx, []string{idRingkasan})
		if err != nil {
			return err
		}
		kunci := models.KunciRincian(rapi.Contract, rapi.Year, rapi.Month)
		milik := false
		for _, k := range lama {
			if k.ID == idRincian {
				milik = true
				continue
			}
			if models.KunciRincian(k.Contract, k.Year, k.Month) == kunci {
				return tolak("CONTRACT %s, YEAR %s, MONTH %s already exists in R/I RISK NAME %s (ID %s)",
					rapi.Contract, kosongTampil(rapi.Year), kosongTampil(rapi.Month), ring.UsedBy, k.ID)
			}
		}
		if idRincian != "" && !milik {
			return ErrRincianTidakAda
		}
		hasil = models.Rincian{ID: idRincian, IDUsedBy: ring.ID, UsedBy: ring.UsedBy, Contract: rapi.Contract,
			Year: rapi.Year, Month: rapi.Month, Risk: rapi.Risk}
		if idRincian == "" {
			if hasil.ID, err = (&pemberiID{l: l}).baru(ctx, tx); err != nil {
				return err
			}
			if err := l.gudang.SisipRincian(ctx, tx, hasil); err != nil {
				return err
			}
		} else if err := l.gudang.UbahRincian(ctx, tx, hasil); err != nil {
			if errors.Is(err, repository.ErrTidakAda) {
				return ErrRincianTidakAda
			}
			return err
		}
		ring.OperatorID, ring.ModifiedDate = potongByte(a.AkunID, models.BatasNama), FormatWaktuPega(l.jam())
		return l.gudang.UbahRingkasan(ctx, tx, ring)
	})
	if err != nil {
		return models.Rincian{}, petaTidakAda(err)
	}
	return hasil, nil
}

// kosongTampil - nilai kosong ditulis "(empty)" di kalimat galat.
func kosongTampil(s string) string {
	if s == "" {
		return "(empty)"
	}
	return s
}
