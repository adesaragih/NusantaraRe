package services

// Reinsurer (paket 4, tiket 05/11) - padanan `SaveSecurityLife_Act` (`Save`
// b5656 `InputSecurityLifeReinsurers.xml`) TANPA procedure `INSERTREINSURER_LIFE`:
//
//	langkah 3 b595  wajib NAME, REINSURERID_LIFE, PCTSHARE, COMMISION, OVR_COMM -> "All value cannot be empty." (b299)
//	langkah 4 b650  PRE=false (selalu jalan): TREATYYEARID/CONTRACTID/REINSTYPEID/NAME <- konteks kontrak (K4)
//	onChange        `SetErrorMessageReinsurer` - 0..100 PctShare/Ricomm, di Pega memeriksa halaman non-life
//	                (tidak pernah mengena); MAKSUDNYA ditegakkan untuk PCTSHARE dan COMMISION (R3, OQ-TCO-17).
//	                OVR_COMM tidak diperiksa Pega - tidak di sini.
//	REINSURER NAME  autocomplete `BrowseCedingCoLife_RD` (AGENT) mengisi REINS ID; nama = master CLIENTNAME
//
// ⛔ Total share (`CountingPercentShare_Act` 3.1: jumlah PCTSHARE) ditampilkan
// dan ditandai bila ≠ 100 - TIDAK memblokir (tiket 05, `[keputusan work owner]`).
// Nama jujur: `TotalShare`, bukan `STDRATING`.

import (
	"context"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

var seratus = apd.New(100, 0)

// ReinsurerMasuk adalah badan permintaan simpan reinsurer (form `Reinsurer List`).
type ReinsurerMasuk struct {
	ID            string `json:"id"`
	ReinsurerName string `json:"reinsurerName"` // REINSURER NAME - nama akhir dari master
	ReinsurerID   string `json:"reinsurerId"`   // REINS ID (diisi autocomplete)
	PctShare      string `json:"pctShare"`      // (%) SHARE
	Komisi        string `json:"komisi"`        // (%) DISCOUNT - kolom COMMISION
	OvrComm       string `json:"ovrComm"`       // (%) OVR COMM
}

// persen - desimal 0..100 (maksud `SetErrorMessageReinsurer`, pesan
// `SetErrorMessageBetween` tidak diekspor - OQ-MCRL-04).
func persen(label, v string) (*apd.Decimal, error) {
	d, err := desimal(label, v)
	if err != nil || d == nil {
		return d, err
	}
	if d.Negative && !d.IsZero() || d.Cmp(seratus) > 0 {
		return nil, fmt.Errorf("%w: %s must be between 0 and 100, got %s", ErrMasukanTidakSah, label, d.Text('f'))
	}
	return d, nil
}

func (m ReinsurerMasuk) keModel() (models.Reinsurer, error) {
	var w wajib
	w.teks("REINSURER NAME", m.ReinsurerName)
	r := models.Reinsurer{ReinsurerID: w.teks("REINS ID", m.ReinsurerID)}
	share := w.teks("(%) SHARE", m.PctShare)
	komisi := w.teks("(%) DISCOUNT", m.Komisi)
	ovr := w.teks("(%) OVR COMM", m.OvrComm)
	if err := w.galat(PesanKosongSemua); err != nil {
		return r, err
	}
	var err error
	if r.PctShare, err = persen("(%) SHARE", share); err != nil {
		return r, err
	}
	if r.Komisi, err = persen("(%) DISCOUNT", komisi); err != nil {
		return r, err
	}
	if r.OvrComm, err = desimal("(%) OVR COMM", ovr); err != nil {
		return r, err
	}
	return r, nil
}

// namaReinsurer membaca nama reinsurer dari master AGENT; di luar master = 422. Pilihan BARU (baris
// baru, atau reinsurer yang diganti) wajib lolos saringan autocomplete `BrowseCedingCoLife_RD` - life
// (b565) dan aktif (b601); baris lama yang reinsurernya kini nonaktif tetap dapat disunting.
func (l *Layanan) namaReinsurer(ctx context.Context, label, id string, pilihanBaru bool) (string, error) {
	m, ada, err := l.gudang.AmbilMasterReinsurer(ctx, id)
	if err != nil {
		return "", err
	}
	if !ada {
		return "", fmt.Errorf("%w: %s %q is not in the reinsurer master", ErrMasukanTidakSah, label, id)
	}
	if pilihanBaru && (!m.Life() || !m.Aktif()) {
		return "", fmt.Errorf("%w: %s %q is not an active life reinsurer", ErrMasukanTidakSah, label, id)
	}
	return m.ClientName, nil
}

// SimpanReinsurer menyimpan reinsurer baru (`m.ID` kosong) di bawah kontrakID,
// atau mengubah reinsurer m.ID (kontrakID boleh kosong; bila terisi harus cocok).
func (l *Layanan) SimpanReinsurer(ctx context.Context, p inti.Pelaku, kontrakID string, m ReinsurerMasuk) (models.Reinsurer, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.Reinsurer{}, err
	}
	if err := pelakuMuat(p, lebarTeks); err != nil {
		return models.Reinsurer{}, err
	}
	var hasil models.Reinsurer
	err := l.tx(ctx, func(tx *db.Tx) error {
		ubah := m.ID != ""
		var reinsurerLama string
		if ubah {
			lama, err := l.ambilReinsurer(ctx, tx, m.ID)
			if err != nil {
				return err
			}
			if kontrakID != "" && lama.TreatyContractID != kontrakID {
				return fmt.Errorf("%w: %s in treaty contract %s", ErrReinsurerTidakAda, m.ID, kontrakID)
			}
			kontrakID, reinsurerLama = lama.TreatyContractID, lama.ReinsurerID
		}
		// ⛔ Induk dikunci lebih dulu - kaskade hapus kontrak yang bersamaan menunggu (K2 tanpa FK).
		if err := l.kunci(ctx, tx, HapusKontrak, kontrakID); err != nil {
			return err
		}
		k, err := l.ambilKontrak(ctx, tx, kontrakID)
		if err != nil {
			return err
		}
		r, err := m.keModel()
		if err != nil {
			return err
		}
		if r.ReinsurerName, err = l.namaReinsurer(ctx, "REINS ID", r.ReinsurerID, !ubah || r.ReinsurerID != reinsurerLama); err != nil {
			return err
		}
		r.TreatyYearID, r.TreatyContractID = k.IDTreatyYear, k.ID
		r.ReinsTypeID, r.ReinsTypeName = k.ReinsTypeID, k.ReinsTypeName
		r.UserID = p.AkunID
		if !ubah {
			if r.ID, err = l.gudang.SisipReinsurer(ctx, tx, r); err != nil {
				return err
			}
		} else {
			r.ID = m.ID
			if err := l.gudang.PerbaruiReinsurer(ctx, tx, r); err != nil {
				return tidakAda(err, ErrReinsurerTidakAda, r.ID)
			}
		}
		hasil, err = l.ambilReinsurer(ctx, tx, r.ID)
		return err
	})
	return hasil, err
}

// jumlahShare - `Local.TotalShare += @toDecimal(.PCTSHARE)`; NULL = 0.
func jumlahShare(share []*apd.Decimal) (*apd.Decimal, error) {
	total := new(apd.Decimal)
	for _, s := range share {
		if s == nil {
			continue
		}
		if _, err := utils.DecimalContext().Add(total, total, s); err != nil {
			return nil, fmt.Errorf("services: summing shares: %w", err)
		}
	}
	return total, nil
}

func sharesReinsurer(d []models.Reinsurer) []*apd.Decimal {
	s := make([]*apd.Decimal, len(d))
	for i, r := range d {
		s[i] = r.PctShare
	}
	return s
}

// BarisLaporanShare - satu kontrak yang total sharenya ≠ 100 (tiket 11).
type BarisLaporanShare struct {
	KontrakID     string `json:"kontrakId"`
	TahunID       string `json:"tahunId"`
	TreatyYear    string `json:"treatyYear"`
	ReinsTypeName string `json:"reinsTypeName"`
	TotalShare    string `json:"totalShare"`
	Selisih       string `json:"selisih"` // 100 - total
}

// LaporanTotalShareBukan100 - kontrak dengan total share ≠ 100, termasuk yang
// tanpa reinsurer (total 0). ⛔ Kemampuan BARU (nol padanan Pega): rute baca
// saja, nol layar (OQ-MCRL-07). tahunID kosong = semua tahun.
func (l *Layanan) LaporanTotalShareBukan100(ctx context.Context, p inti.Pelaku, tahunID string) ([]BarisLaporanShare, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	totals, err := l.gudang.TotalSharePerKontrak(ctx, tahunID)
	if err != nil {
		return nil, err
	}
	hasil := []BarisLaporanShare{}
	for _, t := range totals {
		total := t.Total
		if total == nil {
			total = new(apd.Decimal)
		}
		if total.Cmp(seratus) == 0 {
			continue
		}
		sisa := new(apd.Decimal)
		if _, err := utils.DecimalContext().Sub(sisa, seratus, total); err != nil {
			return nil, fmt.Errorf("services: computing share gap: %w", err)
		}
		hasil = append(hasil, BarisLaporanShare{KontrakID: t.KontrakID, TahunID: t.TahunID, TreatyYear: t.TreatyYear,
			ReinsTypeName: t.ReinsTypeName, TotalShare: total.Text('f'), Selisih: sisa.Text('f')})
	}
	return hasil, nil
}
