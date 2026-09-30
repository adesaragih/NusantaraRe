package services

// Tahun treaty (paket 2, tiket 01) - padanan `SaveTreatyYearLife_Act`
// (`Save` b4944 `InputDtlRetrocessionLife.xml`) TANPA procedure
// `INSERTTREATYYEAR_LIFE` (keputusan o):
//
//	langkah 3 b609  wajib UNDERWRITINGYEAR, TREATYYEAR, STARTDATE, ENDDATE -> "Value cannot be empty." (b313)
//	langkah 4 b664  CARI4 USERID = operator; TGLUPDATE = SYSDATE (procedure `[data DBA]`)
//	langkah 6 b1124 upsert dikunci ID; ID baru '1' || LPAD(TREATYYEAR_LIFE_SEQ, 6, '0') (K3)
//
// ⛔ Baris baru lewat tombol `End Period`, ubah lewat `Edit` - dua rute
// (POST / PUT), supaya ID baris baru tidak pernah datang dari klien.
// ⛔ K4: mengubah tahun memperbarui salinannya di anak dalam transaksi yang
// sama - `TREATYBUSINESS_LIFE.TREATYYEAR` dan tanggal kontrak (R5, OQ-MCRL-10).
// ⛔ Tahun treaty ABADI - tidak ada jalur hapus di sini maupun di mana pun.
// ⛔ Gerbang tahun (K7) dan periode terbalik tidak diperiksa: nol di XML.

import (
	"context"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

// TahunMasuk adalah badan permintaan simpan tahun treaty (form `Input New Data`).
type TahunMasuk struct {
	ID               string `json:"id"`
	TreatyYear       string `json:"treatyYear"` // TRANSACTION YEAR
	UnderwritingYear string `json:"underwritingYear"`
	StartDate        string `json:"startDate"` // YYYY-MM-DD
	EndDate          string `json:"endDate"`
}

func (m TahunMasuk) keModel() (models.TahunTreaty, error) {
	var w wajib
	t := models.TahunTreaty{
		UnderwritingYear: w.teks("UNDERWRITING YEAR", m.UnderwritingYear),
		TreatyYear:       w.teks("TRANSACTION YEAR", m.TreatyYear),
	}
	mulai := w.teks("START DATE", m.StartDate)
	akhir := w.teks("END DATE", m.EndDate)
	if err := w.galat(PesanKosongTahun); err != nil {
		return t, err
	}
	var err error
	if t.StartDate, err = tanggal("START DATE", mulai); err != nil {
		return t, err
	}
	if t.EndDate, err = tanggal("END DATE", akhir); err != nil {
		return t, err
	}
	return t, nil
}

// SimpanTahun menyimpan tahun treaty baru (`baru`) atau mengubah yang ada.
func (l *Layanan) SimpanTahun(ctx context.Context, p inti.Pelaku, m TahunMasuk, baru bool) (models.TahunTreaty, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.TahunTreaty{}, err
	}
	if baru && m.ID != "" {
		return models.TahunTreaty{}, ErrIDDariKlien
	}
	t, err := m.keModel()
	if err != nil {
		return models.TahunTreaty{}, err
	}
	t.UserID = p.AkunID
	var hasil models.TahunTreaty
	err = l.tx(ctx, func(tx *db.Tx) error {
		if baru {
			id, err := l.gudang.SisipTahun(ctx, tx, t)
			if err != nil {
				return err
			}
			t.ID = id
		} else {
			if _, err := l.ambilTahun(ctx, tx, m.ID); err != nil {
				return err
			}
			t.ID = m.ID
			if err := l.gudang.PerbaruiTahun(ctx, tx, t); err != nil {
				return tidakAda(err, ErrTahunTidakAda, t.ID)
			}
			n, err := l.gudang.SalinTahunKeAnak(ctx, tx, t)
			if err != nil {
				return err
			}
			if n > 0 {
				l.catat(fmt.Sprintf("master contract retro life: treaty year %s copied to %d child rows", t.ID, n))
			}
		}
		var err error
		hasil, err = l.ambilTahun(ctx, tx, t.ID)
		return err
	})
	return hasil, err
}
