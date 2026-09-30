package services

// Kontrak dan batas proteksi (paket 3, tiket 02/04) - padanan
// `SaveTreatyLimit_Act` (`Save` b5773 `InputRetroLimitReinsurers.xml`) TANPA
// procedure `INSERTTREATYCONTRACT_LIFE` (keputusan o):
//
//	langkah 3 b609  wajib REINSTYPEID, TREATYSTARTDATE, TREATYENDDATE, B_IDR, IDR, B_USD -> "All value cannot be empty."
//	langkah 4 b674  CARI13/14 <- IDR_SELISIH/USD_SELISIH - nol penulisnya di korpus; di sini DIHITUNG (K5)
//	REINS TYPE      dropdown `BrowseReinsuranceTypeLimit_RD` (FLAG = 1); nama <- `.Note` (`TreatyLimit_TypeProtect` 4.1)
//	TREATY START/END read-only = tanggal tahun induk (`SetValueRetroLimit_TreatyYearLife` b485/b507, R5)
//
// ⛔ `USD` boleh kosong (`[data DBA]`; bukan wajib di b609); `USD_SELISIH` lalu NULL.
// ⛔ Batas bawah > batas atas DITOLAK (tiket 02, tanpa bukti XML - OQ-MCRL-09).
// ⛔ Gerbang tahun TIDAK ditegakkan (K7, OQ-MCRL-01).
// ⛔ K4: mengubah jenis memperbarui salinannya di reinsurer dan business.

import (
	"context"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

// KontrakMasuk adalah badan permintaan simpan kontrak (form `Reins Type`).
//
// ⚠️ Tanggal dan nama jenis TIDAK diterima dari klien: keduanya salinan
// (tahun induk, master jenis). Medan yang terlanjur dikirim diabaikan.
type KontrakMasuk struct {
	ID            string `json:"id"`
	ReinsTypeID   string `json:"reinsTypeId"`
	ReinsTypeName string `json:"reinsTypeName"` // diabaikan - diambil dari master
	IDR           string `json:"idr"`
	USD           string `json:"usd"`
	BIDR          string `json:"bIdr"`
	BUSD          string `json:"bUsd"`
}

func (m KontrakMasuk) keModel(th models.TahunTreaty) (models.Kontrak, error) {
	var w wajib
	k := models.Kontrak{ReinsTypeID: w.teks("REINS TYPE", m.ReinsTypeID), IDTreatyYear: th.ID,
		TreatyStartDate: th.StartDate, TreatyEndDate: th.EndDate}
	if k.TreatyStartDate.IsZero() {
		w.kosong = append(w.kosong, "TREATY START")
	}
	if k.TreatyEndDate.IsZero() {
		w.kosong = append(w.kosong, "TREATY END")
	}
	bidr := w.teks("MINIMUM LIMIT (IDR)", m.BIDR)
	idr := w.teks("MAXIMUM LIMIT (IDR)", m.IDR)
	busd := w.teks("MINIMUM LIMIT (USD)", m.BUSD)
	if err := w.galat(PesanKosongSemua); err != nil {
		return k, err
	}
	var err error
	if k.BIDR, err = desimal("MINIMUM LIMIT (IDR)", bidr); err != nil {
		return k, err
	}
	if k.IDR, err = desimal("MAXIMUM LIMIT (IDR)", idr); err != nil {
		return k, err
	}
	if k.BUSD, err = desimal("MINIMUM LIMIT (USD)", busd); err != nil {
		return k, err
	}
	if k.USD, err = desimal("MAXIMUM LIMIT (USD)", m.USD); err != nil {
		return k, err
	}
	if err := periksaLayer("IDR", k.BIDR.Text('f'), k.IDR.Text('f'), k.BIDR.Cmp(k.IDR)); err != nil {
		return k, err
	}
	if k.USD != nil {
		if err := periksaLayer("USD", k.BUSD.Text('f'), k.USD.Text('f'), k.BUSD.Cmp(k.USD)); err != nil {
			return k, err
		}
	}
	if k.IDRSelisih, err = kurang(k.IDR, k.BIDR); err != nil {
		return k, err
	}
	if k.USDSelisih, err = kurang(k.USD, k.BUSD); err != nil {
		return k, err
	}
	return k, nil
}

// periksaLayer - batas bawah tidak boleh melebihi batas atas; sama = sah
// (layer menaik: batas atas satu layer = batas bawah layer berikutnya).
func periksaLayer(mataUang, bawah, atas string, banding int) error {
	if banding > 0 {
		return fmt.Errorf("%w: MINIMUM LIMIT (%s) %s is greater than MAXIMUM LIMIT (%s) %s",
			ErrMasukanTidakSah, mataUang, bawah, mataUang, atas)
	}
	return nil
}

// namaJenis membaca `.Note` jenis dari master life; di luar master = 422.
func (l *Layanan) namaJenis(ctx context.Context, id string) (string, error) {
	daftar, err := l.gudang.JenisReasuransiLife(ctx)
	if err != nil {
		return "", err
	}
	for _, j := range daftar {
		if j.ID == id {
			return j.Note, nil
		}
	}
	return "", fmt.Errorf("%w: REINS TYPE %q is not a life reinsurance type (FLAG = 1)", ErrMasukanTidakSah, id)
}

// SimpanKontrak menyimpan kontrak baru (`m.ID` kosong) di bawah tahunID, atau
// mengubah kontrak m.ID. Untuk ubah, tahunID boleh kosong (dibaca dari
// kontraknya); bila terisi dan berbeda, kontraknya dianggap tidak ada.
func (l *Layanan) SimpanKontrak(ctx context.Context, p inti.Pelaku, tahunID string, m KontrakMasuk) (models.Kontrak, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.Kontrak{}, err
	}
	var hasil models.Kontrak
	err := l.tx(ctx, func(tx *db.Tx) error {
		ubah := m.ID != ""
		if ubah {
			lama, err := l.ambilKontrak(ctx, tx, m.ID)
			if err != nil {
				return err
			}
			if tahunID != "" && lama.IDTreatyYear != tahunID {
				return fmt.Errorf("%w: %s in treaty year %s", ErrKontrakTidakAda, m.ID, tahunID)
			}
			tahunID = lama.IDTreatyYear
		}
		th, err := l.ambilTahun(ctx, tx, tahunID)
		if err != nil {
			return err
		}
		k, err := m.keModel(th)
		if err != nil {
			return err
		}
		if k.ReinsTypeName, err = l.namaJenis(ctx, k.ReinsTypeID); err != nil {
			return err
		}
		k.UserID = p.AkunID
		if !ubah {
			if k.ID, err = l.gudang.SisipKontrak(ctx, tx, k); err != nil {
				return err
			}
		} else {
			k.ID = m.ID
			if err := l.gudang.PerbaruiKontrak(ctx, tx, k); err != nil {
				return tidakAda(err, ErrKontrakTidakAda, k.ID)
			}
			n, err := l.gudang.SalinKontrakKeAnak(ctx, tx, k)
			if err != nil {
				return err
			}
			if n > 0 {
				l.catat(fmt.Sprintf("master contract retro life: treaty contract %s copied to %d child rows", k.ID, n))
			}
		}
		hasil, err = l.ambilKontrak(ctx, tx, k.ID)
		return err
	})
	return hasil, err
}
