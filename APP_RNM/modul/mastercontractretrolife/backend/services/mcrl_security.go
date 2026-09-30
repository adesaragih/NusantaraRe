package services

// Security reinsurer (paket 5, tiket 06) - padanan `SaveSecurityReinsurerLife_Act`
// (`Save` b5151 `InputSecurityReinsurerLife.xml`) TANPA procedure
// `INSERTSECURITYREINSURER_LIFE`:
//
//	langkah 3 b580  wajib REINSURERNAME, REINSURERID, PCTSHARE -> "All value cannot be empty." (b284)
//	langkah 4 b635  PRE=false (selalu jalan): TREATYYEARID/CONTRACTID/REINSURERID <- reinsurer induk (K4)
//	onChange        `SetErrorMessageReinsurer` (halaman non-life) - maksud 0..100 ditegakkan (R3)
//
// ⛔ `PCTSHARE` security = persen DARI share reinsurer induknya
// (`[fakta bisnis — work owner]`); eksposur terhadap treaty = anak × induk / 100,
// dihitung saat BACA sehingga mengubah share induk langsung mengubah eksposur
// seluruh anaknya tanpa menyentuh baris anak (tiket 06 AC).

import (
	"context"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

// SecurityMasuk adalah badan permintaan simpan security (form `Security Reinsurer`).
type SecurityMasuk struct {
	ID            string `json:"id"`
	ReinsurerName string `json:"reinsurerName"` // SECURITY REINSURER NAME - nama akhir dari master
	ReinsurerID   string `json:"reinsurerId"`   // REINS ID (diisi autocomplete)
	PctShare      string `json:"pctShare"`      // (%) SHARE - persen dari share induk
}

func (m SecurityMasuk) keModel() (models.SecurityReinsurer, error) {
	var w wajib
	w.teks("SECURITY REINSURER NAME", m.ReinsurerName)
	s := models.SecurityReinsurer{ReinsurerID: w.teks("REINS ID", m.ReinsurerID)}
	share := w.teks("(%) SHARE", m.PctShare)
	if err := w.galat(PesanKosongSemua); err != nil {
		return s, err
	}
	var err error
	s.PctShare, err = persen("(%) SHARE", share)
	return s, err
}

// SimpanSecurity menyimpan security baru (`m.ID` kosong) di bawah reinsurerID,
// atau mengubah security m.ID (reinsurerID boleh kosong; bila terisi harus cocok).
func (l *Layanan) SimpanSecurity(ctx context.Context, p inti.Pelaku, reinsurerID string, m SecurityMasuk) (models.SecurityReinsurer, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.SecurityReinsurer{}, err
	}
	var hasil models.SecurityReinsurer
	err := l.tx(ctx, func(tx *db.Tx) error {
		ubah := m.ID != ""
		if ubah {
			lama, err := l.gudang.AmbilSecurity(ctx, tx, m.ID)
			if err != nil {
				return tidakAda(err, ErrSecurityTidakAda, m.ID)
			}
			if reinsurerID != "" && lama.TreatyReinsurerID != reinsurerID {
				return fmt.Errorf("%w: %s under reinsurer %s", ErrSecurityTidakAda, m.ID, reinsurerID)
			}
			reinsurerID = lama.TreatyReinsurerID
		}
		induk, err := l.ambilReinsurer(ctx, tx, reinsurerID)
		if err != nil {
			return err
		}
		s, err := m.keModel()
		if err != nil {
			return err
		}
		if s.ReinsurerName, err = l.namaReinsurer(ctx, "REINS ID", s.ReinsurerID); err != nil {
			return err
		}
		s.TreatyYearID, s.TreatyContractID, s.TreatyReinsurerID = induk.TreatyYearID, induk.TreatyContractID, induk.ID
		s.UserID = p.AkunID
		if !ubah {
			if s.ID, err = l.gudang.SisipSecurity(ctx, tx, s); err != nil {
				return err
			}
		} else {
			s.ID = m.ID
			if err := l.gudang.PerbaruiSecurity(ctx, tx, s); err != nil {
				return tidakAda(err, ErrSecurityTidakAda, s.ID)
			}
		}
		hasil, err = l.gudang.AmbilSecurity(ctx, tx, s.ID)
		return tidakAda(err, ErrSecurityTidakAda, s.ID)
	})
	return hasil, err
}

// Eksposur - share anak × share induk / 100 (persen terhadap treaty); nil bila
// salah satunya kosong.
func Eksposur(anak, induk *apd.Decimal) (*apd.Decimal, error) {
	if anak == nil || induk == nil {
		return nil, nil
	}
	c := utils.DecimalContext()
	hasil := new(apd.Decimal)
	if _, err := c.Mul(hasil, anak, induk); err != nil {
		return nil, fmt.Errorf("services: computing exposure: %w", err)
	}
	if _, err := c.Quo(hasil, hasil, seratus); err != nil {
		return nil, fmt.Errorf("services: computing exposure: %w", err)
	}
	hasil.Reduce(hasil)
	return hasil, nil
}
